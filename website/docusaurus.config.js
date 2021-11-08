const lightCodeTheme = require('prism-react-renderer/themes/github');
const darkCodeTheme = require('prism-react-renderer/themes/palenight');
// With JSDoc @type annotations, IDEs can provide config autocompletion
/** @type {import('@docusaurus/types').DocusaurusConfig} */
(module.exports = {
  title: 'Wails',
  tagline: 'Build applications using Go + HTML + CSS + JS (BETA)',
  url: 'https://wails.io',
  baseUrl: '/',
  onBrokenLinks: 'warn',
  onBrokenMarkdownLinks: 'warn',
  favicon: 'img/favicon.ico',
  organizationName: 'wailsapp', // Usually your GitHub org/user name.
  projectName: 'wails', // Usually your repo name.
  // i18n
  i18n: {
    defaultLocale: 'en',
    locales: ['en', 'zh-Hans'],
    localeConfigs: {
      en: {
        label: 'English',
        direction: 'ltr',
      },
      "zh-Hans": {
        label: '简体中文',
        direction: 'ltr',
      },
    },
  },
  plugins: [
    [
      'docusaurus-plugin-plausible',
      {
        domain: 'wails.io',
      },
    ],
  ],
  presets: [
    [
      '@docusaurus/preset-classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          // Please change this to your repo.
          editUrl: 'https://github.com/wailsapp/wails/edit/master/website',
        },
        blog: {
          showReadingTime: true,
          // Please change this to your repo.
          editUrl:
              'https://github.com/wailsapp/wails/edit/master/website/blog/',
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      }),
    ],
  ],

  themeConfig:
  /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
      ({
        announcementBar: {
          id: 'beta-message',
          content: 'Wails v2 is currently Beta for Windows & Mac. Linux in progress.',
          backgroundColor: '#b00',
          textColor: '#FFF',
          isCloseable: false,
        },
        colorMode: {
          respectPrefersColorScheme: true,
          defaultMode: 'dark',
        },
        navbar: {
          title: 'Wails',
          logo: {
            alt: 'Wails Logo',
            src: 'img/wails-logo-horizontal.svg',
            srcDark: 'img/wails-logo-horizontal-dark.svg',
          },
          items: [
            {
              type: 'docsVersionDropdown', //version
              position: 'left',
              dropdownActiveClassDisabled: true,
              dropdownItemsAfter: [
                // { to: 'https://v1.wails.app', label: '1.13.1', },
                // { to: '/versions', label: 'All versions', }, //Can add custom pages
              ],
            },
            {
              to: 'https://github.com/sponsors/leaanthony',
              label: 'Sponsor',
              position: 'left',
            },
            {
              type: 'doc',
              docId: 'about',
              position: 'right',
              label: 'About',
            },
            {to: '/blog', label: 'Blog', position: 'right'},
            {type: 'localeDropdown', position: 'right',},
            {
              href: 'https://github.com/wailsapp/wails',
              label: 'GitHub',
              position: 'right',
            },
          ],
        },
        footer: {
          style: 'dark',
          links: [
            {
              title: 'Docs',
              items: [
                {
                  label: 'About',
                  to: '/docs/about',
                },
              ],
            },
            {
              title: 'Community',
              items: [
                {
                  label: 'Github',
                  href: 'https://github.com/wailsapp/wails',
                },
                {
                  label: 'Twitter',
                  href: 'https://twitter.com/wailsapp',
                },
                {
                  label: 'Slack',
                  href: 'https://gophers.slack.com/messages/CJ4P9F7MZ/',
                },
                {
                  label: 'Slack invite',
                  href: 'https://invite.slack.golangbridge.org/',
                },
              ],
            },
            {
              title: 'More',
              items: [
                {
                  label: 'Blog',
                  to: '/blog',
                },
              ],
            },
          ],
          copyright: `Copyright © ${new Date().getFullYear()} Lea Anthony. Built with Docusaurus.`,
        },
        prism: {
          theme: lightCodeTheme,
          darkTheme: darkCodeTheme,
        },
      }),
});
