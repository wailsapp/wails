import { defineConfig } from 'vite'
import { qwikVite } from '@builder.io/qwik/optimizer'
import path from "path"

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    qwikVite({
      csr: true,
    }),
  ],
  resolve: {
    alias: {
        "@services": path.resolve(__dirname, "/bindings/changeme/index.ts"),
    },
},
})
