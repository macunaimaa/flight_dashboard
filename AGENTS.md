# AGENTS.md

## Project overview

Real-time aircraft tracking dashboard. Go backend ingests ADS-B data (adsb.lol, OpenSky, or mock), stores in MongoDB, streams updates via WebSocket. React/Cesium frontend renders a 3D globe.

## Commands

```bash
# Full stack (requires Docker)
make docker-up          # builds and starts all services (mongo, redis, backend, frontend)
make docker-down

# Local development (requires local mongo + redis)
make dev-backend        # cd backend && go run ./cmd/server
make dev-frontend       # cd frontend && npm run dev

# Build
make build-backend      # go build -o server ./cmd/server  (output: backend/server)
make build-frontend     # tsc && vite build

# Test
make test-backend       # cd backend && go test ./...
# No frontend test runner configured (no vitest/jest)

# Seed demo data (requires mongosh + Go)
make seed               # creates demo tenant + admin user
```

## Architecture

```
backend/
  cmd/server/main.go        # entrypoint — wires config, DB, repos, services, handlers, router
  cmd/genhash/main.go       # CLI helper — generates bcrypt hash (used by seed script)
  internal/
    config/                  # env-based config (config.Load())
    domain/                  # data models (AircraftLiveState, TrackPoint, Tenant, User, BBox)
    handler/                 # Chi HTTP handlers (auth, aircraft, ws, health)
    middleware/               # request-id, logging, CORS, JWT auth
    ingestion/               # data pipeline: fetcher → normalizer → upsert → WS broadcast
    repository/              # interfaces.go + MongoDB/Redis implementations
    service/                 # business logic layer
    ws/                      # WebSocket hub + client
    pkg/mongo, pkg/geo, pkg/logger  # shared utilities

frontend/
  src/
    api/                     # fetch wrappers (auth, aircraft, client)
    cesium/                  # Cesium viewer, aircraft layer, utils
    components/              # React components (auth, aircraft, layout, common)
    hooks/                   # useWebSocket, useCesiumSync
    store/                   # Zustand stores (aircraftStore, authStore, uiStore)
    ws/                      # WebSocket connection + types
```

## Key facts

- **Backend module**: `github.com/macunaimaa/dashboard/backend`
- **Data source** controlled by `DATA_SOURCE` env var: `"adsblol"` (default), `"opensky"`, or `"mock"`. Legacy `USE_MOCK_DATA=true` still works.
- **Backend loads `.env`** from three locations relative to `cmd/server/`: `../../.env`, `../.env`, `.env` — so a root `.env` file works when running locally.
- **Frontend reads env from parent dir** (`envDir: ".."` in vite.config.ts). Vite env vars: `VITE_API_URL`, `VITE_WS_URL`, `VITE_CESIUM_ION_TOKEN`.
- **Frontend proxies** `/api` → `http://localhost:8080` and `/api/v1/ws` (WebSocket) via Vite dev server.
- **WebSocket auth**: token passed as query param `?token=...` (not header).
- **GeoJSON coordinates**: `[lon, lat]` order throughout (MongoDB standard).
- **Demo login**: `admin@demo.com` / `admin123` (seeded via `make seed`).
- **Weather data** fetched client-side from Open-Meteo API (free, no key required, CORS enabled). Weather panel shows conditions for the selected aircraft's position.
- **News search** proxied through backend (`/api/v1/news/search?q=...`) to GNews API. Requires `NEWS_API_KEY` env var (free at gnews.io).
- **Flight path trails** rendered as Cesium polylines when an aircraft is selected, colored by altitude (green→yellow→blue→purple).
- **News and Weather panels** are toggleable from the TopBar; panels overlay the right side of the globe.

## Constraints

- **Frontend TypeScript is strict**: `strict`, `noUnusedLocals`, `noUnusedParameters`. Unused imports/vars will fail `tsc`.
- **No linter configured** for either backend or frontend. Use `go test ./...` and `tsc && vite build` as verification.
- **No CI/CD workflows** in this repo.
- **MongoDB indexes** auto-created on startup via `pkgmongo.EnsureIndexes()`.
- **Seed script** (`scripts/seed_tenants.sh`) requires `mongosh` CLI and compiles `backend/cmd/genhash` on the fly.

## Environment variables

Copy `.env.example` to `.env`. Required for local dev:
- `MONGO_URI`, `MONGO_DB`, `REDIS_ADDR` — database connections
- `JWT_SECRET` — token signing
- `DATA_SOURCE` — which ingestion source to use
- `OPENSKY_USERNAME` / `OPENSKY_PASSWORD` — only if `DATA_SOURCE=opensky`
- `VITE_CESIUM_ION_TOKEN` — Cesium Ion access token for the frontend
- `NEWS_API_KEY` — GNews API key for the news search panel (optional, free at gnews.io)
