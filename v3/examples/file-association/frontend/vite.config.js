import { defineConfig } from "vite";
import { fileURLToPath } from "node:url";

export default defineConfig({
  resolve: {
    alias: {
      "@wailsio/runtime": fileURLToPath(new URL("../../../internal/runtime/desktop/@wailsio/runtime/src/index.ts", import.meta.url)),
    },
  },
  server: {
    port: Number(process.env.WAILS_VITE_PORT) || 9245,
    strictPort: true,
  },
});
