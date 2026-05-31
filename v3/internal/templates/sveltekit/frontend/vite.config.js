import { sveltekit } from "@sveltejs/kit/vite";
import { defineConfig, searchForWorkspaceRoot } from "vite";
import wails from "@wailsio/runtime/plugins/vite";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [sveltekit(), wails("./bindings")],
  server: {
    host: "127.0.0.1",
    port: Number(process.env.WAILS_VITE_PORT) || 9245,
    strictPort: true,
    fs: {
      allow: [
        // search up for workspace root
        searchForWorkspaceRoot(process.cwd()),
        // your custom rules
        "./bindings/*",
      ],
    },
  },
});
