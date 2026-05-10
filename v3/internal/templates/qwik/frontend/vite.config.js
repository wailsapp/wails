import { defineConfig } from "vite";
import { qwikVite } from "@builder.io/qwik/optimizer";
import wails from "@wailsio/runtime/plugins/vite";

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    host: "127.0.0.1",
  },
  plugins: [
    qwikVite({
      csr: true,
    }),
    wails("./bindings"),
  ],
});
