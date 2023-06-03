# Website

This website is built using [Docusaurus 2](https://docusaurus.io/), a modern
static website generator.

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

language - The language code configured in the i18n field in the
docusaurus.config.js file.

This command starts a local development server and opens up a browser window.
Most changes are reflected live without having to restart the server.

### Translate

After the English source file is updated, run the following command to submit
the source file to Crowdin:

```
npm run crowdin push -- -b <branch>
```

branch - Branch name in crowdin project

Run the following command to pull the translated files in crowdin to the local:

```
npm run crowdin pull -- -b <branch> -l <languageCode>
```

languageCode - **Note** that this refers to the language code in the crowdin
project.

The recommended practice is to update the English source file locally, then
translate the file in crowdin, and finally pull the translated file to the
local.

### Build

```
$ yarn build
```

This command generates static content into the `build` directory and can be
served using any static contents hosting service.
