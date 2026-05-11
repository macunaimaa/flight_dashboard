import { useState } from "react";
import { useAuthStore } from "../../store/authStore";

export function LoginForm() {
  const login = useAuthStore((s) => s.login);
  const [tenantId, setTenantId] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      await login({ tenantId, email, password });
    } catch (err) {
      setError(err instanceof Error ? err.message : "Login failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={styles.container}>
      <div style={styles.card}>
        <h1 style={styles.title}>Aircraft Dashboard</h1>
        <p style={styles.subtitle}>Situational Awareness Platform</p>
        <form onSubmit={handleSubmit} style={styles.form}>
          <input
            type="text"
            placeholder="Tenant ID"
            value={tenantId}
            onChange={(e) => setTenantId(e.target.value)}
            style={styles.input}
            required
          />
          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            style={styles.input}
            required
          />
          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            style={styles.input}
            required
          />
          {error && <p style={styles.error}>{error}</p>}
          <button type="submit" style={styles.button} disabled={loading}>
            {loading ? "Signing in..." : "Sign In"}
          </button>
        </form>
      </div>
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  container: {
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
    height: "100vh",
    background: "#0a0e17",
  },
  card: {
    background: "#111827",
    borderRadius: 12,
    padding: 40,
    width: 360,
    boxShadow: "0 4px 24px rgba(0,0,0,0.5)",
  },
  title: {
    color: "#e0e6ed",
    fontSize: 24,
    fontWeight: 700,
    marginBottom: 4,
    textAlign: "center" as const,
  },
  subtitle: {
    color: "#6b7280",
    fontSize: 13,
    textAlign: "center" as const,
    marginBottom: 32,
  },
  form: {
    display: "flex",
    flexDirection: "column" as const,
    gap: 12,
  },
  input: {
    background: "#1f2937",
    border: "1px solid #374151",
    borderRadius: 6,
    padding: "10px 14px",
    color: "#e0e6ed",
    fontSize: 14,
    outline: "none",
  },
  button: {
    background: "#2563eb",
    color: "white",
    border: "none",
    borderRadius: 6,
    padding: "10px 14px",
    fontSize: 14,
    fontWeight: 600,
    cursor: "pointer",
    marginTop: 8,
  },
  error: {
    color: "#ef4444",
    fontSize: 13,
    margin: 0,
  },
};
