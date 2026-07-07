// Package sse implementa un broker Server-Sent-Events multiplexado a mano sobre net/http.
package sse

import (
	"fmt"
	"net/http"
	"sync"
)

// Event es un frame SSE con tipo + payload ya serializado.
type Event struct {
	Type string
	Data []byte
}

const subBuffer = 256

// Broker multiplexa un único stream `/events` para todas las sesiones activas.
// Cada frame lleva su propio session_id/run_id embebido en Data — el filtrado es responsabilidad del cliente.
type Broker struct {
	mu   sync.Mutex
	subs map[chan Event]struct{}
}

func NewBroker() *Broker {
	return &Broker{subs: make(map[chan Event]struct{})}
}

// Publish emite un evento a todos los subscriptores. Un subscriptor lento se descarta (nunca
// bloquea al productor) — el cliente reconecta; no hay replay por Last-Event-ID todavía (TBD).
func (b *Broker) Publish(eventType string, data []byte) {
	ev := Event{Type: eventType, Data: data}
	b.mu.Lock()
	defer b.mu.Unlock()
	for ch := range b.subs {
		select {
		case ch <- ev:
		default:
			// subscriptor rezagado: se descarta este frame para él, no se bloquea el productor.
		}
	}
}

func (b *Broker) subscribe() chan Event {
	ch := make(chan Event, subBuffer)
	b.mu.Lock()
	b.subs[ch] = struct{}{}
	b.mu.Unlock()
	return ch
}

func (b *Broker) unsubscribe(ch chan Event) {
	b.mu.Lock()
	delete(b.subs, ch)
	b.mu.Unlock()
	close(ch)
}

func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ch := b.subscribe()
	defer b.unsubscribe(ch)

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case ev, ok := <-ch:
			if !ok {
				return
			}
			fmt.Fprintf(w, "event: %s\ndata: %s\n\n", ev.Type, ev.Data)
			flusher.Flush()
		}
	}
}
