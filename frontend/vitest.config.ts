import { defineConfig } from "vitest/config";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  test: {
    globals: true,
    environment: "jsdom",
    setupFiles: ["./src/test/setup.ts"],
    css: false,
    coverage: {
      provider: "v8",
      include: ["src/api/**", "src/store/**", "src/utils/**", "src/ws/**"],
      exclude: ["src/**/*.test.{ts,tsx}", "src/test/**", "src/cesium/**"],
      reporter: ["text", "html", "lcov"],
    },
  },
});
