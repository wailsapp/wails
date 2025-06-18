import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from "path"

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
        "@services": path.resolve(__dirname, "/bindings/changeme/index.ts"),
    },
},
})
