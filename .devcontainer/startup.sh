set -euo pipefail

echo "== versions =="
command -v go >/dev/null && go version || echo "go: missing"
command -v docker >/dev/null && docker --version || echo "docker: missing"
command -v docker >/dev/null && docker compose version || echo "docker compose: missing"

echo "== env =="
echo "PORT=${PORT:-}"
echo "TARGET_DOMAIN=${TARGET_DOMAIN:-}"

echo "== next =="
echo "Run: go mod tidy && go run ."
echo "Or:  docker compose up --build"
