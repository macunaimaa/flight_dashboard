import { setupServer } from "msw/node";

// Per-test handlers are added via `server.use(...)` inside individual tests.
// The base handler list is intentionally empty so unhandled requests fail fast
// (configured in setup.ts via onUnhandledRequest: "error").
export const server = setupServer();
