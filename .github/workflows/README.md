# CI workflows

## `ci.yml`

Runs on every `push` to `main` and every `pull_request` targeting `main`. Two
jobs execute in parallel:

### `backend (go test)`
- `go vet ./...`
- `go build ./...`
- `go test ./... -race -coverprofile=coverage.out`
- Uploads `backend/coverage.out` as an artifact.

### `frontend (vitest)`
- `npm ci`
- `npm run typecheck` (`tsc --noEmit`)
- `npm run build` (`tsc && vite build`)
- `npm run test:coverage` (`vitest run --coverage`)
- Uploads `frontend/coverage/` as an artifact.

## Reproducing locally

```bash
# backend
make test-backend
cd backend && go test ./... -race -coverprofile=coverage.out

# frontend
cd frontend && npm install
npm run typecheck
npm run build
npm run test:run
npm run test:coverage
```

## Known local-only quirks

- **Backend on macOS**: the default `GOCACHE` lives under `~/Library/Caches/go-build/`.
  If you run inside a sandbox that doesn't allow writes there, export
  `GOCACHE="$TMPDIR/go-build"` before `go test`. CI is unaffected.
- **Frontend on Node 22+**: if you see `--localstorage-file` warnings, an
  experimental Node flag is set in `NODE_OPTIONS`. The test suite stubs
  `localStorage` per-test, so warnings are cosmetic.
