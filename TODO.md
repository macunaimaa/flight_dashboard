# TODO — Technical Debt & Easy Wins

Generated from tech-debt audit on 2026-05-11.

---

## Security (fix first)

- [ ] **[HIGH]** Fail-fast in `backend/internal/config/config.go:30` when `JWT_SECRET` is empty or equals the dev default — tokens signed with a known secret are forgeable
- [ ] **[HIGH]** Restrict WS `CheckOrigin` in `backend/internal/handler/ws.go:16` to the same `AllowedOrigins` list used by CORS — current `return true` is a CSWSH vulnerability
- [ ] **[HIGH]** Add `httprate` per-IP rate limit on `POST /api/v1/auth/login` (`backend/cmd/server/main.go:113`) — unlimited bcrypt attempts allow credential stuffing
- [ ] **[HIGH]** Decode JWT `exp` in `frontend/src/store/authStore.ts` `restoreSession` and drop expired tokens; also prefer `sessionStorage` over `localStorage` to limit XSS exposure
- [ ] **[MED]** `regexp.QuoteMeta` + anchor the `$regex` filters in `backend/internal/repository/mongo_aircraft.go:73-77` to prevent ReDoS
- [ ] **[MED]** Enable Mongo `--auth` with an init user in `docker-compose.yml` and bind Mongo/Redis to `127.0.0.1` only; update `.env.example`

---

## Correctness & Reliability

- [ ] **[MED]** Fix race in `backend/internal/ws/hub.go:68-80` — `clients` map is mutated under RLock; promote to write lock or send client to `unregister` channel instead of closing inline
- [ ] **[MED]** Add `Hub.Shutdown(ctx)` and call from graceful-shutdown block in `backend/cmd/server/main.go:139` — currently hub goroutines and per-client goroutines leak on SIGTERM
- [ ] **[MED]** Move `viewportCenterRef` write out of render body into a `useEffect` in `frontend/src/components/layout/AppShell.tsx:22` — mutating refs during render breaks React 18 concurrent mode
- [ ] **[MED]** Poll tenants in parallel with `errgroup` + per-tenant timeouts in `backend/internal/ingestion/pipeline.go` `pollAllTenants` — sequential polling lets one slow tenant starve all others
- [ ] **[MED]** Define a `ErrRateLimited` sentinel in `backend/internal/ingestion` and return it from clients instead of string-matching `"429"` in `pipeline.go:150`; use per-source backoff state
- [ ] **[MED]** Retry Mongo `BulkWrite`/`InsertMany` failures with backoff in `backend/internal/ingestion/pipeline.go:176` — transient errors currently silently drop a full poll cycle
- [ ] **[MED]** Add `select { case ...: case <-ctx.Done(): }` guards to `Hub.Register` and `BroadcastToTenant` in `backend/internal/ws/hub.go:85-91` to prevent deadlock if hub stalls
- [ ] **[LOW]** Wrap `mongoClient.Disconnect` with a 5s timeout context in `backend/cmd/server/main.go:44` — `context.Background()` hangs indefinitely on a hung Mongo
- [ ] **[LOW]** Capture unmarshal errors in `backend/internal/ingestion/opensky_client.go:118-134` `parseStateVector` — currently all `json.Unmarshal` errors are silently ignored
- [ ] **[LOW]** Fix literal `{"°"}` bug in `frontend/src/components/weather/WeatherPanel.tsx:62` — change to `°` so the degree sign actually renders

---

## API & Frontend Polish

- [ ] **[LOW]** Add a 401 interceptor in `frontend/src/api/client.ts:36` that calls `logout()` and surfaces a toast — expired JWT currently leaves the user silently stuck
- [ ] **[LOW]** Replace strict `response.json()` in `frontend/src/api/client.ts:34` with text-then-try-JSON so 502/HTML gateway errors produce readable messages instead of a parse exception
- [ ] **[MED]** Add empty/disconnected state to `frontend/src/components/aircraft/AircraftList.tsx` when `connectionStatus === "disconnected"`
- [ ] **[EASY]** Add keyboard shortcuts in `AppShell.tsx`: `Esc` deselects aircraft, `/` focuses the search field, `w`/`n` toggle Weather/News panels
- [ ] **[EASY]** Propagate `request-id` into WS upgrade/connect/disconnect log entries in `backend/internal/handler/ws.go`

---

## Performance

- [ ] **[MED]** Add Vite `build.rollupOptions.output.manualChunks` to split Cesium into its own chunk in `frontend/vite.config.ts` — Cesium is huge and hurts first-paint if bundled together
- [ ] **[MED]** Update `useCesiumSync` trail polyline in-place (or use `CallbackProperty`) instead of recreating it on every track-points append (`frontend/src/hooks/useCesiumSync.ts:106-156`)
- [ ] **[LOW]** Add chunked pagination to `GetTrackPoints` in `backend/internal/repository/mongo_aircraft.go:189` — hard-coded `SetLimit(10000)` will OOM on long-running aircraft
- [ ] **[LOW]** Sort aircraft lazily in a selector rather than re-sorting on every WS update in `frontend/src/store/aircraftStore.ts:16-35`
- [ ] **[LOW]** Add WS client-to-server `viewport.update` messages and filter broadcasts per-viewport in `backend/internal/ws/hub.go` — currently every client receives every tenant update

---

## Observability & Ops

- [ ] **[MED]** Add real `/api/v1/healthz` (liveness) + `/api/v1/readyz` (Mongo + Redis ping) endpoints — current `backend/internal/handler/health.go` returns a static 200
- [ ] **[MED]** Add a Prometheus `/metrics` endpoint and instrument poll duration, WS connection count, and broadcast queue depth
- [ ] **[LOW]** Convert `frontend/Dockerfile` to a multi-stage build: compile with Node → serve `dist/` from `nginx:alpine` (current image runs the Vite dev server in production)
- [ ] **[LOW]** Drop `version: "3.9"` from `docker-compose.yml` (Compose v2 warns on it)

---

## Developer Experience

- [ ] **[MED]** Add ESLint + `eslint-plugin-react-hooks` + `@typescript-eslint` to `frontend/` and wire `npm run lint` into CI
- [ ] **[EASY]** Add a `.github/workflows/ci.yml` running `go test ./...`, `go vet`, `npm run typecheck`, `npm run test:run` on every push/PR
- [ ] **[EASY]** Add a `lefthook` or `pre-commit` config: `gofmt -l`, `go vet`, `npm run typecheck`, `vitest --run` on staged files
- [ ] **[EASY]** Write a top-level `README.md` covering stack, `make docker-up`, env vars table, architecture diagram link, and demo login
- [ ] **[LOW]** Refactor inline style objects across `frontend/src/components/**/*.tsx` into a `theme.ts` token file or CSS modules

---

## Notes

Priority was computed as `(Impact + Risk) × (6 - Effort)` where each dimension is scored 1–5.
Top security items (JWT secret validation, WS CSRF, login rate limiting) should be addressed before any public deployment.
