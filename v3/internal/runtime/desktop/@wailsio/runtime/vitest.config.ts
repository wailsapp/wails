import { defineConfig } from "vitest/config";

export default defineConfig({
  test: {
    environment: './scripts/protocol-test-environment.mjs',
    testTimeout: 200
  },
});
