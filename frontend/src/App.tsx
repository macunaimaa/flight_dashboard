import { useEffect } from "react";
import { useAuthStore } from "./store/authStore";
import { LoginForm } from "./components/auth/LoginForm";
import { AppShell } from "./components/layout/AppShell";
import { ErrorBoundary } from "./components/common/ErrorBoundary";

export function App() {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);
  const restoreSession = useAuthStore((s) => s.restoreSession);

  useEffect(() => {
    restoreSession();
  }, [restoreSession]);

  return (
    <ErrorBoundary>
      {isAuthenticated ? <AppShell /> : <LoginForm />}
    </ErrorBoundary>
  );
}
