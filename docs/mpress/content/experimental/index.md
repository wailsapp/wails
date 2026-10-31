---
title: "Experimental Features"
description: "A home for in-progress experiments in Wails v3 — what they are, why they exist, and how to give feedback."
slug: "experimental"
sourcePath: "experimental/index.md"
---

@note{type="caution" title="Here be experiments"}
Everything in this section is, by definition, an experiment. Features here are opt-in, off by default, and may change shape, get renamed, or be removed entirely between releases. Don't build anything load-bearing on them without being ready for churn.

@end

## What "Experimental" means

Wails ships experiments in the open. An experiment is an idea we think is promising enough to put in your hands early — but one we haven't fully committed to yet. We're publishing it *because* we want to learn from real use before deciding whether it becomes a permanent, supported part of Wails.

That means a few things are true of everything in this section:

- **It's opt-in.** Experiments never change the default behaviour of `wails3`. You turn them on deliberately (usually with an environment variable or a build flag), and with them off, nothing about your existing workflow changes.
- **It might not survive.** Some experiments graduate into stable features. Others get reworked beyond recognition, or dropped. We'd rather try things in public and learn fast than ship only what we're already certain about.
- **The API isn't frozen.** Names, flags, defaults and behaviour can shift between releases while an experiment is finding its shape. Release notes will call out the changes, but don't expect the stability guarantees that stable features carry.

## We want your feedback

This is the important part. Experiments live or die by what we hear back from people actually using them. If you try one of these, we genuinely want to know:

- Did it work for your project? Where did it fall down?
- Was it faster / clearer / nicer — or not worth the switch?
- What would have to be true for you to use it by default?

The most useful feedback is concrete: what you ran, what you expected, what actually happened. Each experiment has its own thread in the **Experiments** category on GitHub Discussions:

@container{display="grid" columns="2" gap="1rem"}
@linkcard{title="Experiments Discussions" href="https://github.com/wailsapp/wails/discussions/categories/experiments" description="Find the thread for the experiment you're using and tell us how it worked out, what broke, or what's missing."}
@end

## Current experiments

@container{display="grid" columns="2" gap="1rem"}
@linkcard{title="Wake" href="/experimental/wake/" description="An alternative, Wails-aware build runner for your existing Taskfiles. Faster incremental builds, structured output, parallel by default."}
@linkcard{title="CGO-free builds" href="/experimental/#cgo-free-builds-purego" description="Build macOS and Linux apps with CGO_ENABLED=0: no C toolchain, cross-compile without Docker. Enabled with WAILS_PUREGO=1 or the purego build tag."}
@linkcard{title="LLM Control (MCP)" href="/guides/mcp-service/" description="A built-in Model Context Protocol server that lets LLM agents inspect, test, and drive a running Wails app — no user code required, enabled with a build tag."}
@end

## CGO-free builds (purego)

@note{type="caution" title="Experimental Feature"}
The CGO-free backend is opt-in behind the `purego` build tag (or `WAILS_PUREGO=1`) and is **not** the default. Without it, macOS and Linux builds use the cgo backend exactly as before. Behaviour may change between releases.

@end

### What it is

On macOS and Linux, Wails normally talks to the system's UI libraries (AppKit and WKWebView on macOS, GTK and WebKitGTK on Linux) through cgo. The CGO-free backend makes the same calls from pure Go using [purego](https://github.com/ebitengine/purego), which loads the system libraries at run time.

The public application API is unchanged. Select the backend when building; application code does not need to change. Test your app on the destination system before distributing it.

What you gain:

- **No C toolchain for Go compilation.** The Go build does not need Xcode command line tools, `gcc`, or Linux development headers. Frontend builds still need their usual JavaScript tooling.
- **Cross-compiling without Docker.** Build a macOS or Linux binary from any OS with plain `go build`.

Windows already builds without cgo and is unaffected.

### Before enabling it

Use a CLI and framework version that both include this experiment. Updating only the CLI does not add the backend to an older framework dependency in your application's `go.mod`. To try the development version, run these commands from your application directory:

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@master
go get github.com/wailsapp/wails/v3@master
```

Both commands track the moving `master` branch. For reproducible builds, use the same tested commit for the CLI and framework. Follow the [installation guide](/getting-started/installation/) for the Go toolchain and CLI setup.

The experimental macOS and Linux backends support `amd64` and `arm64`. Cross-compilation creates an executable for the requested target; it does not run that executable on the build machine. Test it on the target operating system and architecture.

### Using it

Pick one of these ways to select the backend.

#### With the Wails CLI

In a project generated with this version of the CLI, `build`, `package` and `dev` use the experimental backend when `WAILS_PUREGO` is enabled. In a POSIX shell:

```sh
WAILS_PUREGO=1 wails3 dev
WAILS_PUREGO=1 wails3 build
WAILS_PUREGO=1 wails3 package
```

In PowerShell:

```powershell
$env:WAILS_PUREGO = "1"
wails3 build
```

The CLI recognises `1`, `true`, `on` and `yes`, ignoring case and surrounding whitespace. An empty value does not enable this experiment.

Or pass the tag to a single build or packaging command:

```sh
wails3 build -tags purego
wails3 package -tags purego
```

With the tag set, the updated stock Taskfiles build with `CGO_ENABLED=0` using the host's Go toolchain. For example, build a Linux ARM64 executable from macOS or Windows:

```sh
wails3 build -tags purego GOOS=linux ARCH=arm64
```

Older projects keep their existing Taskfiles. Compare custom build tasks with the updated templates, or use plain Go to test the backend without changing those tasks. The environment variable selects the tag in Wails CLI commands; plain `go build` needs the tag explicitly.

#### With plain Go

Build your frontend first so any assets referenced by `go:embed` exist, then compile the application:

```sh
CGO_ENABLED=0 go build -tags purego,production .
```

`CGO_ENABLED=1 go build -tags purego` works too: the backend is the same, cgo just stays available to the rest of your program.

For a cross-build, also set `GOOS` and `GOARCH`:

```sh
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -tags purego,production -o myapp .
```

On Linux, select the GTK3 stack with `-tags purego,gtk3,production`.

To return to the default backend, remove the `purego` tag and unset `WAILS_PUREGO` (`unset WAILS_PUREGO` in a POSIX shell, or `Remove-Item Env:WAILS_PUREGO` in PowerShell). Also remove an explicit `CGO_ENABLED=0` setting when returning to the macOS or Linux cgo backend.

If you build with `CGO_ENABLED=0` but forget the tag, the build stops with one error that names the fix:

```
undefined: requires_CGO_ENABLED_1_or_build_tag_purego
```

### Runtime requirements

The app loads the system libraries when it starts, so the target machine needs them installed, just as with a cgo build:

| Platform | Libraries | Oldest supported |
|---|---|---|
| macOS | AppKit, WebKit (part of the OS) | Same as the cgo backend |
| Linux (default) | GTK 4, WebKitGTK 6.0 | GTK 4.8, WebKitGTK 2.40 (Debian 12) |
| Linux (`gtk3`) | GTK 3, WebKit2GTK 4.1 | GTK 3.24, WebKit2GTK 2.36 (Ubuntu 22.04) |

Functions newer than these floors are looked up when the app starts and used only when the installed library has them, so one binary runs on old and new systems alike.

The app still uses native UI libraries and may dynamically link system libraries such as Linux's libc. `CGO_ENABLED=0` does not mean a completely static executable. When a required runtime library is missing, install the matching GTK/WebKit stack on the destination system; installing C headers on the build machine does not fix a missing runtime library.

### Packaging and signing

The backend removes the C toolchain requirement from compilation. Each package format keeps its own requirements. In particular, selecting `purego` does not make the existing AppImage deployment tools cross-platform. See [Linux packaging](/guides/build/linux/).

Cross-built macOS applications still need the appropriate signing and notarisation before distribution. Windows installers still need their installer tooling. See [macOS packaging](/guides/build/macos/) and [Windows packaging](/guides/build/windows/). For dependencies that use cgo, `CGO_ENABLED=1` also requires a C compiler and any headers or libraries for the target; selecting the Wails pure-Go backend does not remove those dependencies' requirements.

### Validation

Generated-binding checks and `puregocheck` inspect native symbols, optional-version guards and backend coverage. The Linux parity workflow exercises GTK4 and GTK3 with both backends, including asset requests, binary request bodies and native dispatch; it also cross-compiles Linux ARM64 targets. Cross-compilation checks establish that code builds, while runtime checks exercise the native libraries installed on the test machine.

The [backend test guide](https://github.com/wailsapp/wails/blob/master/v3/TESTING.md) describes the parity probes and end-to-end input tests. When reporting a problem, distinguish a compilation failure from a runtime or packaging failure.

### Known differences

- **Menu accelerators on Linux use application-level registration.** Avoid assigning conflicting accelerators to different menus. Explicitly removing a window's menu with `SetMenu(nil)` also prevents its removed menu shortcuts from activating in that window; other windows retain their menus and shortcuts.
- **macOS open and save panels** run in a separate system process, as they do with cgo; the tests check they open and cancel cleanly.

Found something else? Please tell us.

### Feedback

Tell us how it went in the **CGO-free builds** thread under [Experiments](https://github.com/wailsapp/wails/discussions/categories/experiments): your OS and version, how you built, and what worked or didn't. A failing case with the output of `wails3 doctor` is the most useful report we can get.
