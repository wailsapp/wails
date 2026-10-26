---
title: "Write documentation with M-Press"
description: "Download and install M-Press, preview the Wails docs, and submit a documentation contribution."
sourcePath: "contributing/documentation.md"
---

Wails v3 uses **M-Press** to turn Markdown into this documentation site. You can fix a typo, improve an explanation, add an example or guide, or correct a translation. Documentation-only corrections do not need a prior issue or a failing code test.

You need Git to prepare a contribution and a GitHub account to submit a pull request. The downloadable M-Press executable includes the documentation generator and diagram renderer; previewing the docs needs no Go, Node.js, Wails CLI, or desktop build dependencies. Python 3 is needed for the repository's additional validation scripts.

## Start from a documentation page

Open **Contributing → Documentation** in the sidebar to reach this guide. To work on a particular page, select the **edit icon** in that page's header, then **Edit this page locally**. Copy the command for your operating system and run it in a terminal.

The command uses M-Press if installed; otherwise its installer downloads an official release and verifies the checksum. It prepares a local contribution checkout and opens the selected page in M-Press. Check `mpress version` against the pinned version below before running the repository checks: the automatic installer can download a newer release.

You can edit the source in your own editor while the local preview rebuilds. The wizard also offers **Review changes**, **Commit changes**, **Push contribution branch**, and **Open draft pull request** as separate actions. Follow the checks and PR guidance below before submitting. Translation generation is optional and may need a configured provider; editing and previewing Markdown does not.

For a reproducible manual setup, follow the steps below.

## Download and install M-Press

The repository build script downloads the **latest stable M-Press release**. Use the same release locally: choose an installation method below, run `mpress version`, and compare it with the version shown on the release page.

### Download an executable

Open the [latest stable M-Press release](https://github.com/leaanthony/mpress/releases/latest), or download the matching archive below. ARM64 includes Apple Silicon; AMD64 is the usual Intel/AMD PC architecture.

| Computer | AMD64 / Intel | ARM64 / Apple Silicon |
| --- | --- | --- |
| macOS | [mpress-darwin-amd64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-darwin-amd64.tar.gz) | [mpress-darwin-arm64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-darwin-arm64.tar.gz) |
| Linux | [mpress-linux-amd64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-linux-amd64.tar.gz) | [mpress-linux-arm64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-linux-arm64.tar.gz) |
| Windows | [mpress-windows-amd64.zip](https://github.com/leaanthony/mpress/releases/latest/download/mpress-windows-amd64.zip) | [mpress-windows-arm64.zip](https://github.com/leaanthony/mpress/releases/latest/download/mpress-windows-arm64.zip) |

Download [checksums.txt](https://github.com/leaanthony/mpress/releases/latest/download/checksums.txt) from the same release. In your download directory, calculate the archive's SHA-256 hash and compare it with the line for that filename in `checksums.txt`:

@tabs
[macOS]
```sh
shasum -a 256 mpress-darwin-arm64.tar.gz
```

[Linux]
```sh
sha256sum mpress-linux-amd64.tar.gz
```

[Windows PowerShell]
```powershell
Get-FileHash .\mpress-windows-amd64.zip -Algorithm SHA256
```
@end

Use your actual downloaded filename if you chose another architecture. After the hashes match, extract and install the executable:

@tabs
[macOS]
From the download directory, using the Apple Silicon archive:

```sh
mkdir -p mpress-download
tar -xzf mpress-darwin-arm64.tar.gz -C mpress-download
mkdir -p "$HOME/.local/bin"
install -m 755 mpress-download/mpress "$HOME/.local/bin/mpress"
export PATH="$HOME/.local/bin:$PATH"
mpress version
```

For Intel Macs, substitute `mpress-darwin-amd64.tar.gz`. Add the `export PATH` line to `~/.zshrc` to keep the command available in new terminals.

[Linux]
From the download directory:

```sh
mkdir -p mpress-download
tar -xzf mpress-linux-amd64.tar.gz -C mpress-download
mkdir -p "$HOME/.local/bin"
install -m 755 mpress-download/mpress "$HOME/.local/bin/mpress"
export PATH="$HOME/.local/bin:$PATH"
mpress version
```

For ARM64, substitute `mpress-linux-arm64.tar.gz`. Add the `export PATH` line to your shell's startup file, such as `~/.bashrc`, to keep the command available in new terminals.

[Windows PowerShell]
From the download directory:

```powershell
Expand-Archive .\mpress-windows-amd64.zip -DestinationPath .\mpress-download
New-Item -ItemType Directory -Force "$HOME\bin"
Copy-Item .\mpress-download\mpress.exe "$HOME\bin\mpress.exe"
$env:Path = "$HOME\bin;$env:Path"
mpress version
```

For ARM64, substitute `mpress-windows-arm64.zip`. Add your user `bin` directory to your user **Path** in Windows **Environment Variables**, then reopen the terminal, to keep the command available in future sessions.
@end

### Install with Go

If you already use Go, you can build the latest stable release instead. Check its `go.mod` for the required toolchain; Go’s automatic toolchain selection may download it. The release used for this documentation pass requires Go 1.26.10 or later.

```sh
go install github.com/leaanthony/mpress/cmd/mpress@latest
mpress version
```

Go installs the command in `GOBIN`, or in the `bin` directory under `go env GOPATH` when `GOBIN` is unset. Add that directory to your `PATH` and reopen the terminal if `mpress` is not found.

## Get the documentation source

[Fork wailsapp/wails](https://github.com/wailsapp/wails/fork), then run these commands, replacing `YOUR_USERNAME` with your GitHub username:

```sh
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails
git remote add upstream https://github.com/wailsapp/wails.git
git fetch upstream master
git switch -c docs/my-change upstream/master
```

This starts your documentation branch from the upstream **master** branch. Use your own descriptive branch name in place of `docs/my-change`.

## Run the live preview

From the Wails repository root, where `mpress.yaml` lives:

```sh
mpress dev
```

Open [http://127.0.0.1:3000](http://127.0.0.1:3000) in your browser. For this guide, visit [http://127.0.0.1:3000/contributing/documentation/](http://127.0.0.1:3000/contributing/documentation/). Leave the command running while you edit; saved source changes rebuild the preview. Press **Ctrl+C** to stop it before running the validation commands below.

The repository already contains the M-Press configuration. No new site needs to be initialised. To use another port, run `mpress dev --port 3001` and open the corresponding localhost address.

## Edit or add a page

English source pages live in `docs/mpress/content/`. Translations live under language folders such as `fr/` and `id/`. The edit icon's contribution instructions identify the source file for the current page.

For an existing page, edit its Markdown and preserve the YAML metadata at the top. Use headings, paragraphs, lists, links and fenced code blocks. Preserve paired M-Press components such as `@note` / `@end` and `@tabs` / `@end`. The [M-Press authoring documentation](https://m-press.me/authoring/) describes the supported syntax.

For a new guide, create a file such as `docs/mpress/content/guides/my-guide.md`:

```markdown
---
title: "My guide"
description: "The task this guide helps someone complete."
slug: "guides/my-guide"
sourcePath: "guides/my-guide.md"
---

Explain the task, give working steps, and show the expected result.
```

Add the page to the appropriate section of `docs/mpress/content/_nav.yaml`, unless that section already uses `autogenerate`. Check that a reader can reach it through navigation and links from related pages. New English pages also need complete locale copies before the coverage checks can pass; coordinate those with maintainers or translation contributors.

Generated output lives in `docs/mpress/site/`. Commit the source pages and any intended navigation or static asset changes; generated site files and local caches stay out of the PR.

## Check your contribution

From the repository root:

```sh
python3 docs/mpress/scripts/check_translations.py
mpress build --strict --no-purge-css
mpress check
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

Use `python` on Windows if that is how Python 3 is installed. If [Task](https://taskfile.dev/) is already installed, `task docs:build` runs the same checks with the `mpress` on your PATH.

The translation check fails on missing, empty or orphaned translation files. Stale or damaged translated passages are reported as warnings by default; `--strict` makes those findings fail too. An edit to existing English prose can be submitted while the translation sync catches up. A correction to one translation can change just that language. Keep commands, API names, links, code samples and diagram connections intact; do not use English placeholders as translated content.

Read the changed page in the browser and run any code examples you modify. If a check fails on an untouched page, include the output and relevant context in your PR rather than weakening the check.

## Submit a pull request

Review your diff, commit your source changes, and push your contribution branch to your fork. Open a PR against **master** in `wailsapp/wails`. Describe what changed, link the affected page, and list the checks you ran. Include before/after screenshots for layout changes and platform/version details for changed code examples.

The production site is static; local preview and validation need no deployment credentials or translation provider. Maintainers handle publication after review.

For code contributions and feature proposals, see [Contributing to Wails](/contributing/). For framework internals, see the [Technical Overview](/contributing/overview/).
