// @ts-check
import { defineConfig } from "astro/config";
import starlight from "@astrojs/starlight";
import sitemap from "@astrojs/sitemap";
import starlightLinksValidator from "starlight-links-validator";
import starlightImageZoom from "starlight-image-zoom";
import starlightBlog from "starlight-blog";
import { authors } from "./src/content/authors";
import d2 from 'astro-d2';

// https://astro.build/config
export default defineConfig({
  site: "https://wails.io",
  trailingSlash: "ignore",
  compressHTML: true,
  output: "static",
  build: { format: "directory" },
  devToolbar: { enabled: true },
  integrations: [
    d2(),
    sitemap(),
    starlight({
      title: "",
      titleDelimiter: "",
      logo: {
        dark: "./src/assets/wails-logo-horizontal-dark.svg",
        light: "./src/assets/wails-logo-horizontal-light.svg",
      },
      favicon: "./public/favicon.svg",
      description: "Build beautiful desktop applications using Go and modern web technologies.",
      pagefind: true,
      customCss: ["./src/stylesheets/extra.css"],
      lastUpdated: true,
      pagination: true,
      editLink: {
        baseUrl: "https://github.com/wailsapp/wails/edit/v3-alpha/docs",
      },
      social: {
        github: "https://github.com/wailsapp/wails",
        discord: "https://discord.gg/JDdSxwjhGf",
        "x.com": "https://x.com/wailsapp",
      },
      head: [
        {
          tag: 'script',
          content: `
            document.addEventListener('DOMContentLoaded', () => {
              const socialLinks = document.querySelector('.social-icons');
              if (socialLinks) {
                const sponsorLink = document.createElement('a');
                sponsorLink.href = 'https://github.com/sponsors/leaanthony';
                sponsorLink.className = 'sl-flex';
                sponsorLink.title = 'Sponsor';
                sponsorLink.innerHTML = '<span class="sr-only">Sponsor</span><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="#ef4444" stroke="none"><path d="M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.3 1.5 4.05 3 5.5l7 7Z"/></svg>';
                socialLinks.appendChild(sponsorLink);
              }
            });
          `,
        },
      ],
      defaultLocale: "root",
      locales: {
        root: { label: "English", lang: "en", dir: "ltr" },
      },
      plugins: [
        starlightImageZoom(),
        starlightBlog({
          title: "Wails Blog",
          authors: authors,
        }),
      ],
      sidebar: [
        { label: "Home", link: "/" },

        // Progressive Onboarding - Netflix Principle: Start with the problem
        { label: "Why Wails?", link: "/quick-start/why-wails" },

        {
          label: "Quick Start",
          collapsed: false,
          items: [
            { label: "Installation", link: "/quick-start/installation" },
            { label: "Your First App", link: "/quick-start/first-app" },
            { label: "Next Steps", link: "/quick-start/next-steps" },
          ],
        },

        // Tutorials
        {
          label: "Tutorials",
          collapsed: true,
          autogenerate: { directory: "tutorials" },
        },

        // Core Concepts
        {
          label: "Core Concepts",
          collapsed: true,
          items: [
            { label: "How Wails Works", link: "/concepts/architecture" },
            { label: "Manager API", link: "/concepts/manager-api" },
            { label: "Application Lifecycle", link: "/concepts/lifecycle" },
            { label: "Go-Frontend Bridge", link: "/concepts/bridge" },
            { label: "Build System", link: "/concepts/build-system" },
          ],
        },

        {
          label: "Features",
          collapsed: true,
          items: [
            {
              label: "Windows",
              collapsed: true,
              items: [
                { label: "Window Basics", link: "/features/windows/basics" },
                { label: "Window Options", link: "/features/windows/options" },
                { label: "Multiple Windows", link: "/features/windows/multiple" },
                { label: "Frameless Windows", link: "/features/windows/frameless" },
                { label: "Window Events", link: "/features/windows/events" },
              ],
            },
            {
              label: "Menus",
              collapsed: true,
              items: [
                { label: "Application Menus", link: "/features/menus/application" },
                { label: "Context Menus", link: "/features/menus/context" },
                { label: "System Tray Menus", link: "/features/menus/systray" },
                { label: "Menu Reference", link: "/features/menus/reference" },
              ],
            },
            {
              label: "Bindings & Services",
              collapsed: true,
              items: [
                { label: "Method Binding", link: "/features/bindings/methods" },
                { label: "Services", link: "/features/bindings/services" },
                { label: "Advanced Binding", link: "/features/bindings/advanced" },
                { label: "Best Practices", link: "/features/bindings/best-practices" },
              ],
            },
            {
              label: "Events",
              collapsed: true,
              items: [
                { label: "Event System", link: "/features/events/system" },
                { label: "Application Events", link: "/features/events/application" },
                { label: "Window Events", link: "/features/events/window" },
                { label: "Custom Events", link: "/features/events/custom" },
              ],
            },
            {
              label: "Dialogs",
              collapsed: true,
              items: [
                { label: "File Dialogs", link: "/features/dialogs/file" },
                { label: "Message Dialogs", link: "/features/dialogs/message" },
                { label: "Custom Dialogs", link: "/features/dialogs/custom" },
              ],
            },
            {
              label: "Clipboard",
              collapsed: true,
              autogenerate: { directory: "features/clipboard" },
            },
            {
              label: "Browser",
              collapsed: true,
              autogenerate: { directory: "features/browser" },
            },
            { label: "Drag & Drop", link: "/features/drag-drop" },
            {
              label: "Keyboard",
              collapsed: true,
              autogenerate: { directory: "features/keyboard" },
            },
            {
              label: "Notifications",
              collapsed: true,
              autogenerate: { directory: "features/notifications" },
            },
            {
              label: "Screens",
              collapsed: true,
              autogenerate: { directory: "features/screens" },
            },
            {
              label: "Environment",
              collapsed: true,
              autogenerate: { directory: "features/environment" },
            },
            {
              label: "Platform-Specific",
              collapsed: true,
              autogenerate: { directory: "features/platform" },
            },
          ],
        },

        // Guides - Task-oriented patterns (Netflix: When to use it, when not to use it)
        {
          label: "Guides",
          collapsed: true,
          items: [
            {
              label: "Development",
              collapsed: true,
              items: [
                { label: "Project Structure", link: "/guides/dev/project-structure" },
                { label: "Development Workflow", link: "/guides/dev/workflow" },
                { label: "Debugging", link: "/guides/dev/debugging" },
                { label: "Testing", link: "/guides/dev/testing" },
              ],
            },
            {
              label: "Building & Packaging",
              collapsed: true,
              items: [
                { label: "Building Applications", link: "/guides/build/building" },
                { label: "Build Customization", link: "/guides/build/customization" },
                { label: "Cross-Platform Builds", link: "/guides/build/cross-platform" },
                { label: "Code Signing", link: "/guides/build/signing" },
                { label: "Windows Packaging", link: "/guides/build/windows" },
                { label: "macOS Packaging", link: "/guides/build/macos" },
                { label: "Linux Packaging", link: "/guides/build/linux" },
                { label: "MSIX Packaging", link: "/guides/build/msix" },
              ],
            },
            {
              label: "Distribution",
              collapsed: true,
              items: [
                { label: "Auto-Updates", link: "/guides/distribution/auto-updates" },
                { label: "File Associations", link: "/guides/distribution/file-associations" },
                { label: "Custom Protocols", link: "/guides/distribution/custom-protocols" },
                { label: "Single Instance", link: "/guides/distribution/single-instance" },
              ],
            },
            {
              label: "Integration Patterns",
              collapsed: true,
              items: [
                { label: "Using Gin Router", link: "/guides/patterns/gin-routing" },
                { label: "Gin Services", link: "/guides/patterns/gin-services" },
                { label: "Database Integration", link: "/guides/patterns/database" },
                { label: "REST APIs", link: "/guides/patterns/rest-api" },
              ],
            },
            {
              label: "Advanced Topics",
              collapsed: true,
              items: [
                { label: "Custom Templates", link: "/guides/advanced/custom-templates" },
                { label: "WML (Wails Markup)", link: "/guides/advanced/wml" },
                { label: "Panic Handling", link: "/guides/advanced/panic-handling" },
                { label: "Security Best Practices", link: "/guides/advanced/security" },
              ],
            },
          ],
        },

        // Reference - Comprehensive API docs (Netflix: Complete technical reference)
        {
          label: "API Reference",
          collapsed: true,
          items: [
            { label: "Overview", link: "/reference/overview" },
            {
              label: "Application",
              collapsed: true,
              autogenerate: { directory: "reference/application" },
            },
            {
              label: "Window",
              collapsed: true,
              autogenerate: { directory: "reference/window" },
            },
            {
              label: "Menu",
              collapsed: true,
              autogenerate: { directory: "reference/menu" },
            },
            {
              label: "Events",
              collapsed: true,
              autogenerate: { directory: "reference/events" },
            },
            {
              label: "Dialogs",
              collapsed: true,
              autogenerate: { directory: "reference/dialogs" },
            },
            {
              label: "Runtime",
              collapsed: true,
              autogenerate: { directory: "reference/runtime" },
            },
            {
              label: "CLI",
              collapsed: true,
              autogenerate: { directory: "reference/cli" },
            },
          ],
        },

        // Contributing - Onboarding Wails developers (Netflix: Internal documentation)
        {
          label: "Contributing",
          collapsed: true,
          items: [
            { label: "Getting Started", link: "/contributing/getting-started" },
            { label: "Development Setup", link: "/contributing/setup" },
            {
              label: "Architecture",
              collapsed: true,
              items: [
                { label: "Overview", link: "/contributing/architecture/overview" },
                { label: "CLI Layer", link: "/contributing/architecture/cli" },
                { label: "Runtime Layer", link: "/contributing/architecture/runtime" },
                { label: "Platform Layer", link: "/contributing/architecture/platform" },
                { label: "Build System", link: "/contributing/architecture/build" },
                { label: "Binding System", link: "/contributing/architecture/bindings" },
              ],
            },
            {
              label: "Codebase Guide",
              collapsed: true,
              items: [
                { label: "Repository Structure", link: "/contributing/codebase/structure" },
                { label: "Application Package", link: "/contributing/codebase/application" },
                { label: "Internal Packages", link: "/contributing/codebase/internal" },
                { label: "Platform Bindings", link: "/contributing/codebase/platform" },
                { label: "Testing", link: "/contributing/codebase/testing" },
              ],
            },
            {
              label: "Development Workflows",
              collapsed: true,
              items: [
                { label: "Building from Source", link: "/contributing/workflows/building" },
                { label: "Running Tests", link: "/contributing/workflows/testing" },
                { label: "Debugging", link: "/contributing/workflows/debugging" },
                { label: "Adding Features", link: "/contributing/workflows/features" },
                { label: "Fixing Bugs", link: "/contributing/workflows/bugs" },
              ],
            },
            { label: "Coding Standards", link: "/contributing/standards" },
            { label: "Pull Request Process", link: "/contributing/pull-requests" },
            { label: "Documentation", link: "/contributing/documentation" },
          ],
        },

        // Migration & Troubleshooting
        {
          label: "Migration",
          collapsed: true,
          items: [
            { label: "From v2 to v3", link: "/migration/v2-to-v3" },
            { label: "From Electron", link: "/migration/from-electron" },
          ],
        },

        {
          label: "Troubleshooting",
          collapsed: true,
          autogenerate: { directory: "troubleshooting" },
        },

        // Community & Resources
        {
          label: "Community",
          collapsed: true,
          items: [
            { label: "Links", link: "/community/links" },
            { label: "Templates", link: "/community/templates" },
            {
              label: "Showcase",
              collapsed: true,
              items: [
                { label: "Overview", link: "/community/showcase" },
                {
                  label: "Applications",
                  autogenerate: {
                    directory: "community/showcase",
                    collapsed: true,
                  },
                },
              ],
            },
          ],
        },

        { label: "What's New", link: "/whats-new", badge: { text: "v3", variant: "tip" } },
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
