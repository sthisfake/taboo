#!/usr/bin/env bash
set -euo pipefail

cd /workspaces/taboo

echo "== waiting for docker =="

until docker info >/dev/null 2>&1; do
  sleep 2
done

echo "== starting compose =="

docker compose up -d --build

echo "== waiting for app =="

until curl -s http://localhost:3000 >/dev/null; do
  sleep 2
done

echo "== setting public port =="

gh codespace ports visibility 3000:public || true

echo "== keepalive =="

nohup bash -c '
while true; do
  echo "$(date) keepalive $(uuidgen)"
  sleep 20
done
' >/tmp/keepalive.log 2>&1 &