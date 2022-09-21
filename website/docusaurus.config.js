// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require("prism-react-renderer/themes/github");
const darkCodeTheme = require("prism-react-renderer/themes/dracula");

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: "Wails",
  tagline: "",
  url: "https://wails.io",
  baseUrl: "/",
  onBrokenLinks: "warn",
  onBrokenMarkdownLinks: "warn",
  favicon: "img/favicon.ico",
  organizationName: "wailsapp",
  projectName: "wails",

  i18n: {
    defaultLocale: "en",
    locales: ["en", "zh-Hans", "ja"],
    localeConfigs: {
      en: {
        label: "English",
        direction: "ltr",
        htmlLang: "en-US",
      },
      "zh-Hans": {
        label: "简体中文",
        direction: "ltr",
        htmlLang: "zh-Hans",
      },
      ja: {
        label: "日本語",
        direction: "ltr",
        htmlLang: "ja-JP",
      },
      ru: {
        label: "Русский",
        direction: "ltr",
        htmlLang: "ru-RU",
      },
      ko: {
        label: "한국어",
        direction: "ltr",
        htmlLang: "ko-KR",
      },
    },
  },
  plugins: [],
  presets: [
    [
      "classic",
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          sidebarPath: require.resolve("./sidebars.js"),
          // Please change this to your repo.
          editUrl: "https://github.com/wailsapp/wails/edit/master/website",
        },
        blog: {
          showReadingTime: true,
          // Please change this to your repo.
          editUrl: "https://github.com/wailsapp/wails/edit/master/website/blog",
        },
        theme: {
          customCss: [
            require.resolve("./src/css/custom.css"),
            require.resolve("./src/css/carousel.css"),
          ],
        },
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      navbar: {
        title: "",
        logo: {
          alt: "Wails Logo",
          src: "img/wails-logo-horizontal.svg",
          srcDark: "img/wails-logo-horizontal-dark.svg",
        },
        items: [
          {
            type: "docsVersionDropdown", //version
            position: "left",
            dropdownActiveClassDisabled: true,
            dropdownItemsAfter: [],
          },
          {
            to: "https://github.com/sponsors/leaanthony",
            label: "Sponsor",
            position: "left",
          },
          {
            type: "doc",
            docId: "introduction",
            position: "right",
            label: "Docs",
          },
          {
            to: "/blog",
            label: "Blog",
            position: "right",
          },
          {
            type: "dropdown",
            label: "About",
            position: "right",
            items: [
              {
                to: "/faq",
                label: "FAQ",
              },
              {
                to: "/changelog",
                label: "Changelog",
              },
              {
                to: "/community-guide",
                label: "Community Guide",
              },
              {
                to: "/coc",
                label: "Code of Conduct",
              },
              {
                to: "/credits",
                label: "Credits",
              },
            ],
          },
          {
            type: "localeDropdown",
            position: "right",
            dropdownItemsAfter: [
              {
                to: "/community-guide#documenting",
                label: "Help Us Translate ❤",
              },
            ],
          },
          {
            href: "https://github.com/wailsapp/wails",
            label: "GitHub",
            position: "right",
          },
          {
            type: "search",
            position: "right",
          },
        ],
      },
      footer: {
        style: "dark",
        logo: {
          alt: "Wails Logo",
          src: "img/wails-logo-horizontal.svg",
          srcDark: "img/wails-logo-horizontal-dark.svg",
          href: "https://github.com/wailsapp/wails",
          width: 160,
          height: 50,
        },
        links: [
          {
            title: "Docs",
            items: [
              {
                label: "Introduction",
                to: "/docs/introduction",
              },
              {
                label: "Getting Started",
                to: "/docs/gettingstarted/installation",
              },
              {
                label: "Changelog",
                to: "/changelog",
              },
            ],
          },
          {
            title: "Community",
            items: [
              {
                label: "Github",
                href: "https://github.com/wailsapp/wails",
              },
              {
                label: "Twitter",
                href: "https://twitter.com/wailsapp",
              },
              {
                label: "Slack",
                href: "https://gophers.slack.com/messages/CJ4P9F7MZ/",
              },
              {
                label: "Slack invite",
                href: "https://invite.slack.golangbridge.org/",
              },
            ],
          },
          {
            title: "More",
            items: [
              {
                label: "Blog",
                to: "/blog",
              },
              {
                label: "Awesome",
                href: "https://github.com/wailsapp/awesome-wails",
              },
            ],
          },
        ],
        copyright: `Copyright © ${new Date().getFullYear()} Lea Anthony`,
      },
      tableOfContents: {
        minHeadingLevel: 2,
        maxHeadingLevel: 5,
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
      },
      colorMode: {
        defaultMode: "light",
        disableSwitch: false,
        respectPrefersColorScheme: true,
      },
      algolia: {
        appId: "AWTCNFZ4FF",
        apiKey: "a6c0adbe611ee2535f3da5e8fd7b2200",
        indexName: "wails",
        contextualSearch: true,
      },
    }),
};

module.exports = config;
