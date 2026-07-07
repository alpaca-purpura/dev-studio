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
  exec google-chrome --app="$URL" --user-data-dir="$HOME/.dev-studio/chrome-profile"
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
DESKTOP
  command -v update-desktop-database > /dev/null 2>&1 && update-desktop-database "$DESKTOP_DIR" || true
fi

echo "✓ instalado: $BIN (${VERSION})"
