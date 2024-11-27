// @ts-check
import { defineConfig } from "astro/config";
import starlight from "@astrojs/starlight";
import sitemap from "@astrojs/sitemap";
import starlightLinksValidator from "starlight-links-validator";
import starlightImageZoom from "starlight-image-zoom";
import starlightBlog from "starlight-blog";

// https://astro.build/config
export default defineConfig({
  site: "https://wails.io",
  trailingSlash: "ignore",
  compressHTML: true,
  output: "static",
  build: { format: "directory" },
  devToolbar: { enabled: true },
  integrations: [
    sitemap(),
    starlight({
      title: "Wails",
      logo: {
        dark: "./src/assets/wails-logo-horizontal-dark.svg",
        light: "./src/assets/wails-logo-horizontal-light.svg",
      },
      favicon: "./public/favicon.svg",
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
      defaultLocale: "root",
      locales: {
        root: { label: "English", lang: "en", dir: "ltr" },
        gr: { label: "Greek", lang: "el", dir: "ltr" },
      },
      plugins: [
        // https://starlight-links-validator.vercel.app/configuration/
        starlightLinksValidator(),
        // https://starlight-image-zoom.vercel.app/configuration/
        starlightImageZoom(),
        // https://starlight-blog-docs.vercel.app/configuration
        starlightBlog({
          title: "Wails Blog",
          authors: {
            leaanthony: {
              name: "Lea Anthony",
              title: "Maintainer of Wails",
              url: "https://github.com/leaanthony",
              picture: "https://github.com/leaanthony.png",
            },
            misitebao: {
              name: "Misite Bao",
              title: "Architect",
              url: "https://github.com/misitebao",
              picture: "https://github.com/misitebao.png",
            },
          },
        }),
      ],
      sidebar: [
        {
          label: "What's New",
          link: "/whats-new",
          badge: { text: "New", variant: "tip" },
        },
        {
          label: "Getting Started",
          autogenerate: { directory: "getting-started", collapsed: false },
        },
        {
          label: "Learn",
          collapsed: true,
          autogenerate: { directory: "learn", collapsed: true },
        },
        { label: "Status", link: "/status" },
        {
          label: "API",
          collapsed: true,
          autogenerate: { directory: "api", collapsed: true },
        },
        { label: "Changelog", link: "/changelog" },
        {
          label: "Development",
          collapsed: true,
          autogenerate: { directory: "development", collapsed: true },
        },
      ],
    }),
  ],
});
