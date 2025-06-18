import { defineConfig } from 'vite'
import solid from 'vite-plugin-solid'
import path from "path"

export default defineConfig({
  plugins: [solid()],
  resolve: {
    alias: {
        "@services": path.resolve(__dirname, "/bindings/changeme/index.js"),
    },
},
})
