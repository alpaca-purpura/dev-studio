#!/usr/bin/env bash
# Instalador dogfooding de DevStudio (PB-26, spec workspace-aislado §3.1).
#   scripts/install.sh            → instalación completa (binario + lanzador + icono + app.json)
#   scripts/install.sh --update   → solo rebuild + binario + app.json (lo usa POST /api/update)
# RN-4: acá NO se muta el repo — git solo LEE la versión (rev-parse). Pull/push los hace Chris.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BIN_DIR="$HOME/.local/bin"
BIN="$BIN_DIR/dev-studio"
APP_DIR="$HOME/.dev-studio"
MODE="${1:-full}"

# versión = SHA corto del working tree (+dirty si hay cambios sin commitear) + fecha
SHA="$(git -C "$ROOT" rev-parse --short HEAD 2>/dev/null || echo local)"
DIRTY=""
if ! git -C "$ROOT" diff --quiet 2>/dev/null || ! git -C "$ROOT" diff --cached --quiet 2>/dev/null; then
  DIRTY="+dirty"
fi
VERSION="${SHA}${DIRTY}"
DATE="$(date -Iseconds)"

echo "── DevStudio install (${MODE}) · versión ${VERSION}"

# El proceso lanzado por .desktop hereda un PATH pelado (bug real DH-16.1/DH-18): npm vive
# en nvm (init en .bashrc — bash -lc NO lo trae) y go en ~/.local/go. Fallbacks explícitos.
if ! command -v npm > /dev/null 2>&1; then
  for d in "$HOME"/.nvm/versions/node/*/bin /usr/local/bin; do
    [ -x "$d/npm" ] && PATH="$d:$PATH" && break
  done
fi
if ! command -v go > /dev/null 2>&1; then
  for d in "$HOME/.local/go/bin" /usr/local/go/bin; do
    [ -x "$d/go" ] && PATH="$d:$PATH" && break
  done
fi
command -v npm > /dev/null 2>&1 || { echo "✗ npm no encontrado (ni en nvm) — no puedo compilar la SPA" >&2; exit 1; }
command -v go > /dev/null 2>&1 || { echo "✗ go no encontrado — no puedo compilar el binario" >&2; exit 1; }

echo "→ build SPA (vite)"
npm --prefix "$ROOT/web" run build --silent

echo "→ build binario (go, versión embebida)"
mkdir -p "$BIN_DIR"
go build -C "$ROOT" -ldflags "-X main.version=${VERSION} -X main.buildDate=${DATE}" -o "$BIN" ./cmd/dev-studio

mkdir -p "$APP_DIR"
printf '{"source": "%s"}\n' "$ROOT" > "$APP_DIR/app.json"

if [ "$MODE" != "--update" ]; then
  echo "→ icono + lanzador de escritorio"
  ICON_DIR="$HOME/.local/share/icons/hicolor/scalable/apps"
  mkdir -p "$ICON_DIR"
  cp "$ROOT/assets/dev-studio.svg" "$ICON_DIR/dev-studio.svg"

  # lanzador: levanta el server si no corre y abre la app en modo ventana
  cat > "$BIN_DIR/dev-studio-open" << 'LAUNCHER'
#!/usr/bin/env bash
set -u
# PATH completo aunque el entorno .desktop venga pelado (claude/go/npm viven en ~/.local/bin)
export PATH="$HOME/.local/bin:$HOME/bin:/usr/local/bin:$PATH"
ADDR="127.0.0.1:4173"
URL="http://$ADDR"
if ! curl -sf --max-time 1 "$URL/api/version" > /dev/null 2>&1; then
  nohup "$HOME/.local/bin/dev-studio" --addr "$ADDR" >> "$HOME/.dev-studio/app.log" 2>&1 &
  for _ in $(seq 1 40); do
    curl -sf --max-time 1 "$URL/api/version" > /dev/null 2>&1 && break
    sleep 0.25
  done
fi
if command -v google-chrome > /dev/null 2>&1; then
  # ventana propia, sin diálogos de primera vez del perfil dedicado.
  # --class (X11) + --wayland-app-id (Wayland): WM_CLASS/app-id = dev-studio → la barra de
  # tareas agrupa la ventana bajo NUESTRO .desktop (StartupWMClass), no bajo Google Chrome.
  # OutdatedBuildDetector: mata el globo «No se puede actualizar Chrome» dentro de la
  # ventana app (el updater de Chrome es asunto del sistema, no de DevStudio — DH-18.1).
  # --disable-gpu(+compositing): compositor por software. Sin esto, en máquinas con el driver
  # de GPU trabado Chrome no produce frames → VENTANA BLANCA (procesos gpu-process en estado D,
  # Page.captureScreenshot cuelga). Diagnosticado por CDP 2026-07-09: con GPU off la app
  # renderiza completa. El software compositing es de sobra para esta UI (no hay 3D/canvas pesado).
  exec google-chrome --app="$URL" --user-data-dir="$HOME/.dev-studio/chrome-profile" \
    --class=dev-studio --wayland-app-id=dev-studio \
    --no-first-run --no-default-browser-check \
    --disable-gpu --disable-gpu-compositing \
    --disable-features=DefaultBrowserPrompt,OutdatedBuildDetector
fi
exec xdg-open "$URL"
LAUNCHER
  chmod +x "$BIN_DIR/dev-studio-open"

  DESKTOP_DIR="$HOME/.local/share/applications"
  mkdir -p "$DESKTOP_DIR"
  cat > "$DESKTOP_DIR/dev-studio.desktop" << DESKTOP
[Desktop Entry]
Type=Application
Name=DevStudio
Comment=Construir y mantener software basado en proceso y arquitectura
Exec=$BIN_DIR/dev-studio-open
Icon=dev-studio
Terminal=false
Categories=Development;
StartupWMClass=dev-studio
DESKTOP
  command -v update-desktop-database > /dev/null 2>&1 && update-desktop-database "$DESKTOP_DIR" || true
fi

echo "✓ instalado: $BIN (${VERSION})"
