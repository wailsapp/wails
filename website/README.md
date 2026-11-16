# Website

This website is built using [Docusaurus 2](https://docusaurus.io/), a modern static website generator.

### Installation

```
$ npm
```

### Local Development

```
$ npm run start
```

Other languages:

```
npm run start -- --locale <language>
```

language - The language code configured in the i18n field in the docusaurus.config.js file.

This command starts a local development server and opens up a browser window. Most changes are reflected live without having to restart the server.

### Translations

The Wails v2 documentation retains its existing translations. Crowdin imports,
source uploads and build-time progress requests have been retired. A Crowdin
account or token is not required to build the site.

Submit corrections to the existing locale files under `website/i18n/` through
GitHub. For current Wails v3 documentation, use the
[M-Press contribution guide](https://v3.wails.io/contributing/documentation/).

### Hosting

The existing Cloudflare Pages project serves the same v2 build on both
`wails.io` and `v2.wails.io`. Keep `wails.io` attached until the separate v3
launch is ready. Adding the v2 hostname does not require copying or moving
any documentation.

### Build

```
$ yarn build
```

This command generates static content into the `build` directory and can be served using any static contents hosting service.
