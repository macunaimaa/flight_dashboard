import { create } from "zustand";
import { apiClient } from "../api/client";
import { login as apiLogin } from "../api/auth";
import type { LoginRequest } from "../api/types";

interface AuthState {
  token: string | null;
  tenantId: string | null;
  role: string | null;
  isAuthenticated: boolean;

  login: (req: LoginRequest) => Promise<void>;
  logout: () => void;
  restoreSession: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  token: null,
  tenantId: null,
  role: null,
  isAuthenticated: false,

  login: async (req: LoginRequest) => {
    const resp = await apiLogin(req);
    apiClient.setToken(resp.token);
    localStorage.setItem("token", resp.token);
    localStorage.setItem("tenantId", resp.tenantId);
    localStorage.setItem("role", resp.role);
    set({
      token: resp.token,
      tenantId: resp.tenantId,
      role: resp.role,
      isAuthenticated: true,
    });
  },

  logout: () => {
    apiClient.setToken(null);
    localStorage.removeItem("token");
    localStorage.removeItem("tenantId");
    localStorage.removeItem("role");
    set({
      token: null,
      tenantId: null,
      role: null,
      isAuthenticated: false,
    });
  },

  restoreSession: () => {
    const token = localStorage.getItem("token");
    const tenantId = localStorage.getItem("tenantId");
    const role = localStorage.getItem("role");
    if (token) {
      apiClient.setToken(token);
      set({ token, tenantId, role, isAuthenticated: true });
    }
  },
}));
