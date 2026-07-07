// Package usecase orquesta el dominio contra los puertos (agente, store, eventos) — sin conocer transporte.
package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/alpacapurpura/dev-studio/internal/domain"
	"github.com/alpacapurpura/dev-studio/internal/ports"
)

// ErrBusy: la sesión ya tiene un turno en vuelo (un-turno-a-la-vez).
var ErrBusy = errors.New("session busy: turn already in flight")

// ErrNotFound: no existe una sesión con ese id.
var ErrNotFound = errors.New("session not found")

// EventPublisher es el sub-puerto mínimo que el caso de uso necesita del transporte realtime.
type EventPublisher interface {
	Publish(eventType string, data []byte)
}

// dockFrame es lo que viaja por el canal realtime hacia el frontend — un evento por sesión/turno.
type dockFrame struct {
	SessionID       string `json:"session_id"`
	Kind            string `json:"kind"` // status | init | delta | result | error
	Text            string `json:"text,omitempty"`
	Status          string `json:"status,omitempty"`
	ClaudeSessionID string `json:"claude_session_id,omitempty"`
	Model           string `json:"model,omitempty"`
}

type sessionRuntime struct {
	meta       *domain.Session
	live       ports.AgentSession
	assembling strings.Builder
}

// SessionService es el registro en memoria de sesiones + su ciclo de vida de conductor.
type SessionService struct {
	mu    sync.Mutex
	rt    map[string]*sessionRuntime
	order []string

	agent   ports.AgentPort
	store   ports.SessionStore
	pub     EventPublisher
	baseCtx context.Context
}

func NewSessionService(baseCtx context.Context, agent ports.AgentPort, store ports.SessionStore, pub EventPublisher) (*SessionService, error) {
	s := &SessionService{
		rt:      make(map[string]*sessionRuntime),
		agent:   agent,
		store:   store,
		pub:     pub,
		baseCtx: baseCtx,
	}
	sessions, err := store.Load(baseCtx)
	if err != nil {
		return nil, fmt.Errorf("session service: load registry: %w", err)
	}
	for i := range sessions {
		sess := sessions[i]
		sess.Status = domain.StatusIdle // ningún conductor se relanza al arrancar
		s.rt[sess.ID] = &sessionRuntime{meta: &sess}
		s.order = append(s.order, sess.ID)
	}
	return s, nil
}

func newID() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return "s" + hex.EncodeToString(b)
}

func (s *SessionService) List() []domain.Session {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]domain.Session, 0, len(s.order))
	for _, id := range s.order {
		out = append(out, *s.rt[id].meta)
	}
	return out
}

func (s *SessionService) Get(id string) (domain.Session, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.rt[id]
	if !ok {
		return domain.Session{}, false
	}
	return *r.meta, true
}

// NewSession parametriza la creación de una sesión (spec shell §4.1): desde el shell toda
// sesión nace ligada a un repo + historia + rol; los campos extra son opcionales para
// mantener compatibilidad con sesiones sueltas (tests, F1).
type NewSession struct {
	Nombre   string
	Cwd      string
	RepoID   string
	Historia *domain.Historia
	Rol      string
}

// Create registra una sesión nueva. Ningún proceso `claude` se lanza todavía —
// el conductor arranca recién en el primer Turn().
func (s *SessionService) Create(nombre, cwd string) (domain.Session, error) {
	return s.CreateSession(NewSession{Nombre: nombre, Cwd: cwd})
}

// CreateSession registra una sesión con su ligadura completa (repo/historia/rol).
func (s *SessionService) CreateSession(p NewSession) (domain.Session, error) {
	s.mu.Lock()
	sess := domain.Session{
		ID: newID(), Nombre: p.Nombre, Cwd: p.Cwd, Status: domain.StatusIdle,
		RepoID: p.RepoID, Historia: p.Historia, Rol: p.Rol, Conv: []domain.Turn{},
	}
	s.rt[sess.ID] = &sessionRuntime{meta: &sess}
	s.order = append(s.order, sess.ID)
	err := s.persistLocked()
	s.mu.Unlock()
	return sess, err
}

func (s *SessionService) Rename(id, nombre string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.rt[id]
	if !ok {
		return ErrNotFound
	}
	r.meta.Nombre = nombre
	return s.persistLocked()
}

// Close cierra el conductor vivo (si lo hay) y borra la sesión del registro.
func (s *SessionService) Close(id string) error {
	s.mu.Lock()
	r, ok := s.rt[id]
	if !ok {
		s.mu.Unlock()
		return ErrNotFound
	}
	delete(s.rt, id)
	for i, oid := range s.order {
		if oid == id {
			s.order = append(s.order[:i], s.order[i+1:]...)
			break
		}
	}
	err := s.persistLocked()
	live := r.live
	s.mu.Unlock()

	if live != nil {
		_ = live.Close()
	}
	return err
}

// Turn envía un mensaje a la sesión. Si no hay conductor vivo, lo spawnea (o resume) primero.
// Responde de inmediato (no espera la respuesta del asistente) — el resultado llega por eventos.
func (s *SessionService) Turn(id, text string) error {
	s.mu.Lock()
	r, ok := s.rt[id]
	if !ok {
		s.mu.Unlock()
		return ErrNotFound
	}
	if r.meta.Status == domain.StatusStreaming {
		s.mu.Unlock()
		return ErrBusy
	}
	if r.live == nil {
		if err := s.spawnLocked(id, r); err != nil {
			s.mu.Unlock()
			return err
		}
	}
	r.meta.Status = domain.StatusStreaming
	r.meta.Conv = append(r.meta.Conv, domain.Turn{Role: "user", Text: text})
	live := r.live
	s.publish(dockFrame{SessionID: id, Kind: "status", Status: "streaming"})
	_ = s.persistLocked()
	s.mu.Unlock()

	return live.Send(s.baseCtx, text)
}

// spawnLocked asume s.mu ya tomado.
func (s *SessionService) spawnLocked(id string, r *sessionRuntime) error {
	live, err := s.agent.Spawn(s.baseCtx, ports.SpawnOpts{Resume: r.meta.ClaudeSessionID, Cwd: r.meta.Cwd})
	if err != nil {
		return fmt.Errorf("session %s: spawn: %w", id, err)
	}
	r.live = live
	go s.consume(id, live)
	return nil
}

// consume drena los eventos normalizados de un conductor y los traduce a frames + estado.
func (s *SessionService) consume(id string, live ports.AgentSession) {
	for ev := range live.Events() {
		s.mu.Lock()
		r, ok := s.rt[id]
		if !ok {
			s.mu.Unlock()
			continue
		}
		switch ev.Kind {
		case ports.EventInit:
			r.meta.ClaudeSessionID = ev.ClaudeSessionID
			r.meta.Model = ev.Model
			s.publish(dockFrame{SessionID: id, Kind: "init", ClaudeSessionID: ev.ClaudeSessionID, Model: ev.Model})
		case ports.EventDelta:
			r.assembling.WriteString(ev.Text)
			s.publish(dockFrame{SessionID: id, Kind: "delta", Text: ev.Text})
		case ports.EventResult:
			text := r.assembling.String()
			if text == "" {
				text = ev.Text
			}
			r.meta.Conv = append(r.meta.Conv, domain.Turn{Role: "assistant", Text: text})
			r.assembling.Reset()
			r.meta.Status = domain.StatusIdle
			_ = s.persistLocked()
			s.publish(dockFrame{SessionID: id, Kind: "result", Text: text, Status: "idle"})
		case ports.EventError:
			r.meta.Status = domain.StatusIdle
			r.live = nil
			s.publish(dockFrame{SessionID: id, Kind: "error", Text: ev.Err, Status: "idle"})
		}
		s.mu.Unlock()
	}
	// el canal se cerró: el proceso terminó. El próximo Turn() relanza (spawnLocked).
	s.mu.Lock()
	if r, ok := s.rt[id]; ok {
		r.live = nil
	}
	s.mu.Unlock()
}

func (s *SessionService) publish(f dockFrame) {
	if s.pub == nil {
		return
	}
	data, err := json.Marshal(f)
	if err != nil {
		return
	}
	s.pub.Publish("dock", data)
}

// persistLocked asume s.mu ya tomado.
func (s *SessionService) persistLocked() error {
	sessions := make([]domain.Session, 0, len(s.order))
	for _, id := range s.order {
		sessions = append(sessions, *s.rt[id].meta)
	}
	return s.store.Save(s.baseCtx, sessions)
}
