import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import cesium from "vite-plugin-cesium";

export default defineConfig({
  envDir: "..",
  plugins: [react(), cesium()],
  server: {
    port: 5173,
    host: true,
    proxy: {
      "/api/v1/ws": {
        target: "http://localhost:8080",
        ws: true,
      },
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
});
