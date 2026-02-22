// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require("prism-react-renderer/themes/github");
const darkCodeTheme = require("prism-react-renderer/themes/dracula");

const { getTranslationProgress } = require("./src/api/crowdin.js");

module.exports = async function configCreatorAsync() {
  const translationProgress = await getTranslationProgress();
  return {
    title: "Wails",
    tagline: "",
    url: "https://wails.io",
    baseUrl: "/",
    onBrokenLinks: "warn",
    onBrokenMarkdownLinks: "warn",
    favicon: "img/favicon.ico",
    organizationName: "wailsapp",
    projectName: "wails",

    webpack: {
      jsLoader: (isServer) => ({
        loader: require.resolve("swc-loader"),
        options: {
          jsc: {
            parser: {
              syntax: "typescript",
              tsx: true,
            },
            target: "es2017",
          },
          module: {
            type: isServer ? "commonjs" : "es6",
          },
        },
      }),
    },
    i18n: {
      defaultLocale: "en",
      locales: ["en", "zh-Hans", "ja", "ru", "ko", "fr", "pt"],
      localeConfigs: {
        en: {
          label: "English",
          direction: "ltr",
          htmlLang: "en-US",
        },
        "zh-Hans": {
          label: `简体中文 (${translationProgress["zh-CN"]}%)`,
          direction: "ltr",
          htmlLang: "zh-Hans",
        },
        ja: {
          label: `日本語 (${translationProgress["ja"]}%)`,
          direction: "ltr",
          htmlLang: "ja-JP",
        },
        ru: {
          label: `Русский (${translationProgress["ru"]}%)`,
          direction: "ltr",
          htmlLang: "ru-RU",
        },
        ko: {
          label: `한국어 (${translationProgress["ko"]}%)`,
          direction: "ltr",
          htmlLang: "ko-KR",
        },
        fr: {
          label: `Français (${translationProgress["fr"]}%)`,
          direction: "ltr",
          htmlLang: "fr",
        },
        pt: {
          label: `Português (${translationProgress["pt-PT"]}%)`,
          direction: "ltr",
          htmlLang: "pt-PT",
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
            editUrl:
              "https://github.com/wailsapp/wails/edit/master/website/blog",
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
                  label: "Discord",
                  href: "https://discord.gg/JDdSxwjhGf",
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
          additionalLanguages: ["json5"],
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
};
