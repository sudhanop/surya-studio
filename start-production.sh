#!/bin/bash

# Visual window title (supported by some Linux terminal emulators)
echo -ne "\033]0;Surya Photography - Production\007"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

# Ensure logs directory exists
mkdir -p "logs"

# Load environment variables from .env if present
if [ -f ".env" ]; then
  echo "Loading environment from .env"
  # Export variables while ignoring comments (#) and blank lines
  export $(grep -v '^#' .env | xargs)
else
  echo "No .env found. Using default local settings."
fi

# Fallback default values
export APP_ENV="${APP_ENV:-production}"
export API_PORT="${API_PORT:-8080}"
export WEB_PORT="${WEB_PORT:-3000}"

export NEXT_PUBLIC_API_URL="${NEXT_PUBLIC_API_URL:-http://localhost:$API_PORT}"
export NEXT_PUBLIC_SITE_URL="${NEXT_PUBLIC_SITE_URL:-http://localhost:$WEB_PORT}"
export NEXT_PUBLIC_MEDIA_URL="${NEXT_PUBLIC_MEDIA_URL:-http://localhost:$API_PORT}/uploads"

export FRONTEND_URL="${FRONTEND_URL:-$NEXT_PUBLIC_SITE_URL}"
export CORS_ORIGINS="${CORS_ORIGINS:-$NEXT_PUBLIC_SITE_URL}"
export PUBLIC_MEDIA_URL="${PUBLIC_MEDIA_URL:-$NEXT_PUBLIC_MEDIA_URL}"

echo ""
echo "[1/2] Starting API..."

# Go Pre-flight & Build Check
cd "$SCRIPT_DIR/backend"
if [ ! -f "server" ]; then
  if ! command -v go &> /dev/null; then
    echo "Go is not installed or not in PATH. Install Go and retry."
    read -p "Press [Enter] key to exit..."
    exit 1
  fi
  echo "Compiling Go production binary..."
  go build -o server ./cmd/server
fi

# Run API in background, passing variables explicitly and logging output
export PORT="$API_PORT"
./server > "$SCRIPT_DIR/logs/api.log" 2>&1 &
API_PID=$!

# Wait 3 seconds for the API to initialize
sleep 3

echo "[2/2] Starting website..."
cd "$SCRIPT_DIR/frontend"

# Node/NPM Pre-flight Check
if ! command -v npm &> /dev/null; then
  echo "Node.js/npm is not installed or not in PATH. Install Node.js LTS and retry."
  # Kill backend before exiting to avoid leaving orphaned background processes
  kill $API_PID 2>/dev/null
  read -p "Press [Enter] key to exit..."
  exit 1
fi

# Install dependencies if missing
if [ ! -d "node_modules" ]; then
  echo "Installing frontend dependencies..."
  npm install
fi

# Build Next.js production asset bundle if standalone server doesn't exist
if [ ! -d ".next/standalone" ]; then
  echo "Building production frontend (this may take a moment)..."
  npm run build
fi

# Run Next.js production standalone server in the background
export NODE_ENV=production
export HOSTNAME="0.0.0.0"
export PORT="$WEB_PORT"

node .next/standalone/server.js > "$SCRIPT_DIR/logs/web.log" 2>&1 &
WEB_PID=$!

echo ""
echo "API:   http://localhost:$API_PORT"
echo "Site:  http://localhost:$WEB_PORT"
echo "Admin: http://localhost:$WEB_PORT/admin/login"
echo ""
echo "Logs are being written to the 'logs/' folder."
echo "Keep this window open. Press Ctrl+C to shut down both servers safely."

# Trap termination signals to clean up background processes cleanly
trap "echo -e '\nShutting down servers...'; kill $API_PID $WEB_PID 2>/dev/null; exit" SIGINT SIGTERM

# Keep the script running to monitor background tasks
wait
