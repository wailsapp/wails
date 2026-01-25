// @ts-check
import { defineConfig } from "astro/config";
import starlight from "@astrojs/starlight";
import sitemap from "@astrojs/sitemap";
import starlightLinksValidator from "starlight-links-validator";
import starlightImageZoom from "starlight-image-zoom";
import starlightBlog from "starlight-blog";
import { authors } from "./src/content/authors";
import d2 from 'astro-d2';
import react from '@astrojs/react';

// https://astro.build/config
export default defineConfig({
  site: "https://wails.io",
  vite: {
    resolve: {
      alias: {
        '@components': '/src/components'
      }
    }
  },
  trailingSlash: "ignore",
  compressHTML: true,
  output: "static",
  build: { format: "directory" },
  devToolbar: { enabled: true },
  integrations: [
    react(),
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
      social: [
        { icon: 'github', label: 'GitHub', href: 'https://github.com/wailsapp/wails' },
        { icon: 'discord', label: 'Discord', href: 'https://discord.gg/JDdSxwjhGf' },
        { icon: 'x.com', label: 'X', href: 'https://x.com/wailsapp' },
      ],
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
            {
              label: "Drag & Drop",
              collapsed: true,
              autogenerate: { directory: "features/drag-and-drop" },
            },
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
                { label: "Server Build", link: "/guides/server-build", badge: { text: "Experimental", variant: "caution" } },
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
            { label: "Application", link: "/reference/application" },
            { label: "Window", link: "/reference/window" },
            { label: "Menu", link: "/reference/menu" },
            { label: "Events", link: "/reference/events" },
            { label: "Dialogs", link: "/reference/dialogs" },
            { label: "Frontend Runtime", link: "/reference/frontend-runtime" },
            { label: "CLI", link: "/reference/cli" },
          ],
        },

        // Contributing
        {
          label: "Contributing",
          collapsed: true,
          items: [
            { label: "Getting Started", link: "/contributing/getting-started" },
            { label: "Development Setup", link: "/contributing/setup" },
            { label: "Coding Standards", link: "/contributing/standards" },
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

        { label: "What's New", link: "/whats-new" },
        { label: "Status", link: "/status" },
        { label: "Changelog", link: "/changelog" },
        {
          label: "Sponsor",
          link: "https://github.com/sponsors/leaanthony",
          badge: { text: "❤️" },
        },
        { label: "Credits", link: "/credits" },
      ],
    }),
  ],
});
