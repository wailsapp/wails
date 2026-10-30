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

The result is the same app with the same API. You don't change any code.

What you gain:

- **No C toolchain.** Builds need only Go. No Xcode command line tools, no `gcc`, no `-dev` packages.
- **Cross-compiling without Docker.** Build a macOS or Linux binary from any OS with plain `go build`.
- **Faster builds.** No C compilation step.

Windows already builds without cgo and is unaffected.

### Using it

Pick any one of these. They all select the same backend.

#### With the Wails CLI

Set the environment variable once, and `build`, `package` and `dev` all use it. No Taskfile edits are needed:

```sh
WAILS_PUREGO=1 wails3 dev
WAILS_PUREGO=1 wails3 build
WAILS_PUREGO=1 wails3 package
```

Or pass the tag to a single command:

```sh
wails3 build -tags purego
wails3 package -tags purego
```

With the tag set, the stock Taskfiles build with `CGO_ENABLED=0` and run natively on every host, so a Linux binary builds on a Mac (and the other way round) without the cross-compiling Docker image:

```sh
WAILS_PUREGO=1 wails3 build GOOS=linux ARCH=arm64
```

#### With plain Go

```sh
CGO_ENABLED=0 go build -tags purego .
```

`CGO_ENABLED=1 go build -tags purego` works too: the backend is the same, cgo just stays available to the rest of your program.

On Linux, combine it with the legacy stack as usual: `-tags purego,gtk3`.

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

### How it is tested

Every native function the backend calls is checked by a script, not by hand:

- A static check lists every C symbol, Objective-C class and selector the backend uses, and fails if any purego function is an empty stub where the cgo one is not.
- On clean macOS and Ubuntu (arm64 and amd64) machines, every symbol is resolved against the real system libraries, and any call newer than the supported floor must be guarded.
- The same parity test runs on both backends, and an end-to-end test drives menus, key bindings, dialogs, context menus and the clipboard with real input. Screenshots are kept as evidence for every visual check.

### Known differences

- **Menu accelerators on Linux are app-wide.** If two windows' menus bind the same shortcut, the first one registered can fire whichever window has focus. The cgo backend does the same.
- **macOS open and save panels** run in a separate system process, as they do with cgo; the tests check they open and cancel cleanly.

Found something else? Please tell us.

### Feedback

Tell us how it went in the **CGO-free builds** thread under [Experiments](https://github.com/wailsapp/wails/discussions/categories/experiments): your OS and version, how you built, and what worked or didn't. A failing case with the output of `wails3 doctor` is the most useful report we can get.
