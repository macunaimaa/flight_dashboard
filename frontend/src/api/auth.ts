import { apiClient } from "./client";
import type { LoginRequest, LoginResponse } from "./types";

export async function login(req: LoginRequest): Promise<LoginResponse> {
  const response = await apiClient.post<LoginResponse>(
    "/api/v1/auth/login",
    req
  );
  return response.data;
}
