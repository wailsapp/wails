import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig, searchForWorkspaceRoot } from 'vite';
import path from "path"

export default defineConfig({
    server: {
        fs: {
          allow: [
            // search up for workspace root
            searchForWorkspaceRoot(process.cwd()),
            // your custom rules
            './bindings/*',
          ],
        },
    },
	plugins: [sveltekit()],
  resolve: {
    alias: {
        "@services": path.resolve(__dirname, "/bindings/changeme/index.ts"),
    },
},
});
