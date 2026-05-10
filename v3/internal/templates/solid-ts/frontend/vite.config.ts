import { defineConfig } from "vite";
import solid from "vite-plugin-solid";
import wails from "@wailsio/runtime/plugins/vite";

export default defineConfig({
  server: {
    host: "127.0.0.1",
  },
  plugins: [solid(), wails("./bindings")],
});
