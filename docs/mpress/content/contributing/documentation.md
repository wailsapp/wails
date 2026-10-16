---
title: "Correct the documentation"
description: "Submit a correction PR for the Wails v3 documentation using M-Press."
sourcePath: "contributing/documentation.md"
---

Correction PRs are welcome. Fix typos, broken links, outdated examples,  
unclear explanations, or translations. You do not need an issue or a failing  
code test for a documentation-only correction.

## Preview locally

Fork [wailsapp/wails](https://github.com/wailsapp/wails/fork), clone your fork,  
and create a branch from `master`.

Install the pinned documentation generator:

```sh
go install github.com/leaanthony/mpress/cmd/mpress@v1.0.17
```

You can also download a verified binary from the  
[M-Press v1.0.17 release](https://github.com/leaanthony/mpress/releases/tag/v1.0.17).

From the Wails repository root:

```sh
mpress version
mpress dev
```

Edit `.md` source files in `docs/mpress/content/`. English is the default  
language and lives directly in that directory. Existing translations live in  
language folders such as `fr/` and `id/`. The preview rebuilds as you save.

Preserve the metadata block at the top of each page and paired `@...` / `@end`  
components. Ordinary paragraphs, headings, lists and fenced code are editable  
as text. Do not edit generated files in `docs/mpress/site/`.

## Check the correction

```sh
python3 docs/mpress/scripts/check_translations.py
mpress build --strict --no-purge-css
mpress check
# Use `python` on Windows if `python3` is unavailable.
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

Check the changed page in the browser and run any code examples you modify.  
For translation corrections, compare the complete changed passage with English.  
Keep commands, API names, links, code samples and diagram connections intact.  
Use natural technical language, preserve requirements and caveats, and translate  
visible diagram labels, navigation labels and image descriptions as well as prose.

Every published language must have a complete translation of every English page.  
Do not use English placeholders or fallback pages. A correction to one translation  
can change just that language. If you change English meaning, update the matching  
pages in the other published languages; unrelated pages do not need regeneration.

## Submit a pull request

Open a PR against `master`. Describe what was wrong, explain your correction,  
and list the checks you ran. Include screenshots for visible layout changes  
and platform/version details for changed code examples.

Cloudflare credentials and private services are not required. Public PR checks  
build and validate the static site without deployment credentials.

For code changes and feature proposals, see [Contributing to Wails](/contributing/).  
For the internals, see the [Technical Overview](/contributing/overview/).
