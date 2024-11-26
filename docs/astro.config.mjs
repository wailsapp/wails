// @ts-check
import { defineConfig } from "astro/config";
import starlight from "@astrojs/starlight";
import sitemap from "@astrojs/sitemap";
import starlightLinksValidator from "starlight-links-validator";
import starlightImageZoom from "starlight-image-zoom";

// https://astro.build/config
export default defineConfig({
  site: "https://wails.io",
  trailingSlash: "ignore",
  compressHTML: true,
  build: { format: "directory" },
  devToolbar: { enabled: true },
  integrations: [
    sitemap(),
    starlight({
      title: "Wails",
      description: "Build desktop applications using Go & Web Technologies.",
      pagefind: true,
      lastUpdated: true, // Note, this needs git clone with fetch depth 0 to work
      pagination: true,
      editLink: {
        // TODO: update this
        baseUrl: "https://github.com/wailsapp/wails/edit/v3-alpha/docs",
      },
      social: {
        github: "https://github.com/wailsapp/wails",
        discord: "https://discord.gg/JDdSxwjhGf",
        "x.com": "https://x.com/wailsapp",
      },
      plugins: [
        // https://starlight-links-validator.vercel.app/configuration/
        starlightLinksValidator(),
        // https://starlight-image-zoom.vercel.app/configuration/
        starlightImageZoom(),
      ],
      sidebar: [
        { label: "Changelog", link: "/changelog" },
        {
          label: "What's New",
          link: "/whats-new",
          badge: { text: "New", variant: "tip" },
        },
        { label: "Status", link: "/status" },
        {
          label: "API",
          autogenerate: { directory: "api", collapsed: false },
        },
        {
          label: "Getting Started",
          autogenerate: { directory: "getting-started", collapsed: false },
        },
        {
          label: "Learn",
          autogenerate: { directory: "learn", collapsed: false },
        },
        {
          label: "Development",
          autogenerate: { directory: "development", collapsed: false },
        },
      ],
    }),
  ],
});
