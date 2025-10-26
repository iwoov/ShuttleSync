#!/usr/bin/env bash
set -euo pipefail

# Build frontend (Vue), copy built assets into backend for Go embed,
# build backend (Go), and run the server.

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
FRONTEND_DIR="$ROOT_DIR/frontend"
BACKEND_DIR="$ROOT_DIR/backend"
TARGET_WEB_DIR="$BACKEND_DIR/web"
BIN_DIR="$BACKEND_DIR/bin"
BIN_PATH="$BIN_DIR/shuttlesync"

need() { command -v "$1" >/dev/null 2>&1 || { echo "[ERROR] Missing required tool: $1" >&2; exit 1; }; }

echo "[Check] tools"
need go
need npm

echo "[1/4] Building frontend"
pushd "$FRONTEND_DIR" >/dev/null
if [ -f package-lock.json ]; then
  echo "[npm] Detected lockfile, running npm ci"
  if ! npm ci; then
    echo "[npm] npm ci failed; falling back to npm install to sync lockfile"
    npm install
  fi
else
  npm install
fi
npm run build
popd >/dev/null

echo "[2/4] Syncing frontend dist -> backend/web"
mkdir -p "$TARGET_WEB_DIR"
if command -v rsync >/dev/null 2>&1; then
  rsync -a --delete "$FRONTEND_DIR/dist/" "$TARGET_WEB_DIR/"
else
  # Fallback without deletion (may leave stale files if names change)
  cp -R "$FRONTEND_DIR/dist/." "$TARGET_WEB_DIR/"
fi

echo "[3/4] Building backend"
pushd "$BACKEND_DIR" >/dev/null
mkdir -p "$BIN_DIR"
go mod tidy
go build -o "$BIN_PATH"
popd >/dev/null

echo "[4/4] Running server -> $BIN_PATH"
echo "Open: http://127.0.0.1:4050/"
exec "$BIN_PATH"
