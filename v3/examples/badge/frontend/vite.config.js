import { defineConfig } from "vite";

export default defineConfig({
  server: {
    port: Number(process.env.WAILS_VITE_PORT) || 9245,
    strictPort: true,
  },
});
