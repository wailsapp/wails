// @ts-check
import { defineConfig } from "astro/config";
import starlight from "@astrojs/starlight";
import sitemap from "@astrojs/sitemap";
import starlightLinksValidator from "starlight-links-validator";
import starlightImageZoom from "starlight-image-zoom";
import starlightBlog from "starlight-blog";
import { authors } from "./src/content/authors";

// https://astro.build/config
export default defineConfig({
  // TODO: update this
  site: "https://wails.io",
  trailingSlash: "ignore",
  compressHTML: true,
  output: "static",
  build: { format: "directory" },
  devToolbar: { enabled: true },
  integrations: [
    sitemap(),
    starlight({
      title: "",
      // If a title is added, also update the delimiter.
      titleDelimiter: "",
      logo: {
        dark: "./src/assets/wails-logo-horizontal-dark.svg",
        light: "./src/assets/wails-logo-horizontal-light.svg",
      },
      favicon: "./public/favicon.svg",
      description: "Build desktop applications using Go & Web Technologies.",
      pagefind: true,
      customCss: ["./src/stylesheets/extra.css"],
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
        // Example of how a new language is added.
        // After this, you create a directory named after the language inside content/docs/
        // with the same structure as the root language
        // eg content/docs/gr/changelog.md or content/docs/gr/api/application.mdx
        // gr: { label: "Greek", lang: "el", dir: "ltr" },
      },
      plugins: [
        // https://starlight-links-validator.vercel.app/configuration/
        // starlightLinksValidator({
          // exclude: [
        //     // TODO: Fix these links in the blog/wails-v2-released file
        //     // "/docs/reference/options#theme",
        //     // "/docs/reference/options#customtheme",
        //     // "/docs/guides/application-development#application-menu",
        //     // "/docs/reference/runtime/dialog",
        //     // "/docs/reference/options#windowistranslucent",
        //     // "/docs/reference/options#windowistranslucent-1",
        //     // "/docs/guides/windows-installer",
        //     // "/docs/reference/runtime/intro",
        //     // "/docs/guides/obfuscated",
        //     // "/docs/howdoesitwork#calling-bound-go-methods",
        //   ],
        // }),
        // https://starlight-image-zoom.vercel.app/configuration/
        starlightImageZoom(),
        // https://starlight-blog-docs.vercel.app/configuration
        starlightBlog({
          title: "Wails Blog",
          authors: authors,
        }),
      ],
      sidebar: [
        { label: "Home", link: "/" },
        {
          label: "Getting Started",
          autogenerate: { directory: "getting-started", collapsed: false },
        },
        {
          label: "Tutorials",
          collapsed: true,
          autogenerate: { directory: "tutorials", collapsed: true },
        },
        {
          label: "What's New",
          link: "/whats-new",
          badge: { text: "New", variant: "tip" },
        },
        { label: "v3 Alpha Feedback", link: "/feedback" },
        {
          label: "Learn",
          collapsed: true,
          autogenerate: { directory: "learn", collapsed: true },
        },
        {
          label: "Guides",
          collapsed: true,
          autogenerate: { directory: "guides", collapsed: true },
        },
        // {
        //   label: "API",
        //   collapsed: true,
        //   autogenerate: { directory: "api", collapsed: true },
        // },
        {
          label: "Community",
          collapsed: true,
          items: [
            { label: "Links", link: "/community/links" },
            { label: "Templates", link: "/community/templates" },
            {
              label: "Showcase",
              autogenerate: {
                directory: "community/showcase",
                collapsed: true,
              },
            },
          ],
        },
        // {
        //   label: "Development",
        //   collapsed: true,
        //   autogenerate: { directory: "development", collapsed: true },
        // },
        { label: "Status", link: "/status" },
        { label: "Changelog", link: "/changelog" },
        {
          label: "Sponsor",
          link: "https://github.com/sponsors/leaanthony",
          badge: { text: "❤️" },
        },
        {
          label: "Credits",
          link: "/credits",
          badge: { text: "👑" },
        },
      ],
    }),
  ],
});
