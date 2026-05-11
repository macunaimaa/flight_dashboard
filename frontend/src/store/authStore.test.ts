import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { http, HttpResponse } from "msw";
import { server } from "../test/server";
import { useAuthStore } from "./authStore";

// Some Node versions inject a broken localStorage global via --localstorage-file.
// Replace it with an in-memory Storage shim for these tests.
function makeMemStorage(): Storage {
  const map = new Map<string, string>();
  return {
    get length() {
      return map.size;
    },
    clear: () => map.clear(),
    getItem: (k) => (map.has(k) ? (map.get(k) as string) : null),
    setItem: (k, v) => {
      map.set(k, String(v));
    },
    removeItem: (k) => {
      map.delete(k);
    },
    key: (i) => Array.from(map.keys())[i] ?? null,
  };
}

function resetStore() {
  useAuthStore.setState({
    token: null,
    tenantId: null,
    role: null,
    isAuthenticated: false,
  });
}

describe("authStore", () => {
  beforeEach(() => {
    vi.stubGlobal("localStorage", makeMemStorage());
    resetStore();
  });
  afterEach(() => {
    resetStore();
    vi.unstubAllGlobals();
  });

  it("login() stores token, tenantId, role and flips isAuthenticated", async () => {
    server.use(
      http.post("/api/v1/auth/login", () =>
        HttpResponse.json({
          data: {
            token: "tok-abc",
            expiresAt: "2030-01-01T00:00:00Z",
            tenantId: "t1",
            role: "admin",
          },
        })
      )
    );

    await useAuthStore
      .getState()
      .login({ tenantId: "t1", email: "a@b.com", password: "pw" });

    const s = useAuthStore.getState();
    expect(s.token).toBe("tok-abc");
    expect(s.tenantId).toBe("t1");
    expect(s.role).toBe("admin");
    expect(s.isAuthenticated).toBe(true);
    expect(localStorage.getItem("token")).toBe("tok-abc");
    expect(localStorage.getItem("tenantId")).toBe("t1");
    expect(localStorage.getItem("role")).toBe("admin");
  });

  it("login() propagates errors from the API", async () => {
    server.use(
      http.post("/api/v1/auth/login", () =>
        HttpResponse.json({ error: "invalid credentials" }, { status: 401 })
      )
    );

    await expect(
      useAuthStore.getState().login({ tenantId: "t1", email: "a", password: "p" })
    ).rejects.toThrow(/invalid credentials/);

    expect(useAuthStore.getState().isAuthenticated).toBe(false);
  });

  it("logout() clears state and localStorage", () => {
    useAuthStore.setState({
      token: "x",
      tenantId: "t1",
      role: "admin",
      isAuthenticated: true,
    });
    localStorage.setItem("token", "x");
    localStorage.setItem("tenantId", "t1");
    localStorage.setItem("role", "admin");

    useAuthStore.getState().logout();

    const s = useAuthStore.getState();
    expect(s.token).toBeNull();
    expect(s.isAuthenticated).toBe(false);
    expect(localStorage.getItem("token")).toBeNull();
    expect(localStorage.getItem("role")).toBeNull();
    expect(localStorage.getItem("tenantId")).toBeNull();
  });

  it("restoreSession() reads from localStorage when token present", () => {
    localStorage.setItem("token", "stored");
    localStorage.setItem("tenantId", "t1");
    localStorage.setItem("role", "viewer");

    useAuthStore.getState().restoreSession();

    const s = useAuthStore.getState();
    expect(s.token).toBe("stored");
    expect(s.tenantId).toBe("t1");
    expect(s.role).toBe("viewer");
    expect(s.isAuthenticated).toBe(true);
  });

  it("restoreSession() is a no-op when no token stored", () => {
    useAuthStore.getState().restoreSession();
    expect(useAuthStore.getState().isAuthenticated).toBe(false);
  });
});
