import { defineConfig } from 'astro/config';
import d2 from 'astro-d2';
// https://astro.build/config
export default defineConfig({
  integrations: [d2()]
});