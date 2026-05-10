import { defineConfig } from "vite";
import preact from "@preact/preset-vite";
import wails from "@wailsio/runtime/plugins/vite";

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    host: "127.0.0.1",
  },
  plugins: [preact(), wails("./bindings")],
});
