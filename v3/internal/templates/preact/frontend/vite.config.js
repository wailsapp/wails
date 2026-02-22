import { defineConfig } from "vite";
import preact from "@preact/preset-vite";
import wails from "@wailsio/runtime/plugins/vite";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [preact(), wails("./bindings")],
});
