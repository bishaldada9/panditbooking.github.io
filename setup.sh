#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT_DIR"

echo "============================================================"
echo "  Bishal Puja Sewa - Docker Setup"
echo "============================================================"
echo

echo "[1/5] Checking Docker..."
if ! command -v docker >/dev/null 2>&1; then
  echo "[ERROR] Docker is not installed or not on PATH."
  echo "        Install Docker Desktop or Docker Engine with Compose."
  exit 1
fi
docker compose version >/dev/null
echo "  [OK] Docker Compose is available"

echo
echo "[2/5] Preparing environment files..."
if [ ! -f .env ]; then
  cp .env.example .env
  echo "  [OK] Created .env"
else
  echo "  [OK] .env already exists"
fi
if [ ! -f backend/.env ]; then
  cp backend/.env.example backend/.env
  echo "  [OK] Created backend/.env for local non-Docker runs"
else
  echo "  [OK] backend/.env already exists"
fi

set -a
# shellcheck disable=SC1091
. ./.env
set +a
FRONTEND_HOST_PORT="${FRONTEND_HOST_PORT:-3000}"
BACKEND_HOST_PORT="${BACKEND_HOST_PORT:-8081}"

echo
echo "[3/5] Building and starting containers..."
docker compose up -d --build postgres redis backend frontend

echo
echo "[4/5] Waiting for backend health..."
healthy=0
for _ in $(seq 1 40); do
  if curl -fsS "http://127.0.0.1:${BACKEND_HOST_PORT}/health" >/dev/null 2>&1; then
    healthy=1
    break
  fi
  sleep 2
done
if [ "$healthy" -ne 1 ]; then
  echo "[ERROR] Backend did not become healthy. Showing logs:"
  docker compose logs --tail=120 backend
  exit 1
fi
echo "  [OK] Backend is healthy"

echo
echo "[5/5] Seeding demo data..."
docker compose run --rm seed

echo
echo "============================================================"
echo "  Setup Complete"
echo "============================================================"
echo "Frontend: http://localhost:${FRONTEND_HOST_PORT}"
echo "Backend:  http://localhost:${BACKEND_HOST_PORT}"
echo
echo "Default accounts:"
echo "Admin:    admin@bishalpujasewa.com / AdminPass123!"
echo "Pandit:   pandit.atri@example.com / PanditPass123!"
echo "Customer: customer@example.com / CustomerPass123!"
echo
echo "Useful commands:"
echo "docker compose logs -f"
echo "docker compose down"
echo "docker compose down -v   # reset database too"
