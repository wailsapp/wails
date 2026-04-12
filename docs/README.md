# Wails v3 Documentation

[![Built with Starlight](https://astro.badg.es/v2/built-with-starlight/tiny.svg)](https://starlight.astro.build)

## Prerequisites

* [D2 Diagram Library](https://d2lang.com/). You can download the binary manually from [Github releases](https://github.com/terrastruct/d2/releases) or use the following

```
# Use install script:
curl -fsSL https://d2lang.com/install.sh | sh -s --

# Install from source with Go
go install oss.terrastruct.com/d2@latest
```

World-class documentation for Wails v3, redesigned following Netflix documentation principles.

## 📚 Documentation Redesign (2025-10-01)

This documentation has been completely redesigned to follow the **Netflix approach** to developer documentation:

- **Problem-first framing** - Start with why, not what
- **Progressive disclosure** - Multiple entry points for different skill levels
- **Real production examples** - No toy code
- **Story-Code-Context pattern** - Why → How → When
- **Scannable content** - Clear structure, visual aids

**Status:** Foundation complete (~20%), ready for content migration

See [IMPLEMENTATION_SUMMARY.md](./IMPLEMENTATION_SUMMARY.md) for full details.

## 🚀 Project Structure

Inside of your Astro + Starlight project, you'll see the following folders and
files:

```sh
.
├── public/
├── src/
│   ├── assets/
│   ├── content/
│   │   ├── docs/
│   │   └── config.ts
│   └── env.d.ts
├── astro.config.mjs
├── package.json
└── tsconfig.json
```

Starlight looks for `.md` or `.mdx` files in the `src/content/docs/` directory.
Each file is exposed as a route based on its file name.

Images can be added to `src/assets/` and embedded in Markdown with a relative
link.

Static assets, like favicons, can be placed in the `public/` directory.

## 🧞 Commands

All commands are run from the root of the project, from a terminal:

| Command                   | Action                                           |
| :------------------------ | :----------------------------------------------- |
| `npm install`             | Installs dependencies                            |
| `npm run dev`             | Starts local dev server at `localhost:4321`      |
| `npm run build`           | Build your production site to `./dist/`          |
| `npm run preview`         | Preview your build locally, before deploying     |
| `npm run astro ...`       | Run CLI commands like `astro add`, `astro check` |
| `npm run astro -- --help` | Get help using the Astro CLI                     |

## 👀 Want to learn more?

Check out [Starlight’s docs](https://starlight.astro.build/), read
[the Astro documentation](https://docs.astro.build), or jump into the
[Astro Discord server](https://astro.build/chat).
