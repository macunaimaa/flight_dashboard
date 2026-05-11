export function LoadingSpinner() {
  return (
    <div style={styles.container}>
      <div style={styles.spinner} />
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
  spinner: {
    width: 32,
    height: 32,
    border: "3px solid #1f2937",
    borderTop: "3px solid #2563eb",
    borderRadius: "50%",
    animation: "spin 0.8s linear infinite",
  },
};
