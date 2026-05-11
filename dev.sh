#!/usr/bin/env bash
set -euo pipefail

SESSION="dashboard"
ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"

if tmux has-session -t "$SESSION" 2>/dev/null; then
  echo "Session '$SESSION' already exists. Attaching..."
  exec tmux attach -t "$SESSION"
fi

tmux new-session -d -s "$SESSION" -c "$ROOT_DIR"

tmux send-keys -t "$SESSION" "cd backend && go run ./cmd/server" Enter

tmux split-window -h -t "$SESSION" -c "$ROOT_DIR"
tmux send-keys -t "$SESSION" "cd frontend && npm run dev" Enter

tmux attach -t "$SESSION"
