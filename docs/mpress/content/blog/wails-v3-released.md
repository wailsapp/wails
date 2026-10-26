---
title: "Wails v3.0.0 Released"
description: "Wails v3.0.0 is here: multiple windows, system tray, native menus, services and generated bindings, an updater, mobile support and a new documentation site."
authors: ["leaanthony"]
tags: ["wails","v3","release"]
date: "2026-10-26"
slug: "blog/wails-v3-released"
image: "/assets/videos/wails-v3-release-clean-poster.webp"
sourcePath: "blog/wails-v3-released.md"
---

@video{src="/assets/videos/wails-v3-release-clean-av1.mp4" poster="/assets/videos/wails-v3-release-clean-poster.webp" title="Wails v3.0.0 release film — best viewed full screen" width="1920" height="1080"}

## It's here!

Today I'm releasing Wails v3.0.0.

I'll be honest: it has been a long road. The first v3 alpha tag went out in January 2023, and a great many people have built, tested, broken and fixed things along the way since then. I'm truly grateful to every one of you. If you want the candid version of that journey, I wrote about it in the [v3 beta announcement](/blog/wails-v3-beta/).

## What *is* Wails?

If you're new here: Wails lets Go programmers build desktop applications with the web frontend tools they already know, using the native WebView on each platform rather than an embedded browser. You write your backend in Go, your UI in HTML, CSS and JavaScript (or whichever framework you like), and Wails wraps it all up into a single native application.

## Why v3?

Wails v2 made it straightforward to build a Go application with a modern web frontend, and it has served a lot of applications well. But its single-window, context-driven runtime made some everyday desktop work harder than it should have been.

v3 starts from a different model. Applications, windows, services, events and platform capabilities are explicit objects. You create an application, register services, create windows, and talk to the object that owns the behaviour you need. It's easier to follow as an app grows, and things like multiple windows are a normal part of the model rather than a workaround.

## What's in v3.0.0

- **Multiple windows.** Create, manage and close windows at runtime, each with its own configuration and lifecycle.
- **System tray.** Give your app a persistent presence, with menus, light and dark icons, template icons on macOS, and the option to attach a window to the tray icon.
- **Native menus and dialogs.** Application menus, context menus and native dialogs on each platform.
- **Services and generated bindings.** Your backend stays ordinary Go code. Bindings are generated with static source analysis, so comments and meaningful parameter names carry through to the frontend, with a single `wails3 generate bindings` command.
- **A build system you can see.** Projects use a [Taskfile](https://taskfile.dev)-based build you can inspect, extend and debug, instead of one opaque command.
- **Cross-platform.** Windows (amd64 and arm64), macOS (Intel and Apple Silicon) and Linux (amd64 and arm64), with GTK4 and WebKitGTK 6.0 as the Linux default. The legacy GTK3 stack is still available with `-tags gtk3` until v3.1.
- **iOS and Android.** The same `main.go` and frontend can run on mobile. Mobile support is still marked experimental in the docs, so treat it as something to explore and give feedback on rather than a stable promise.
- **Built with Go 1.26.** The v3 module builds with Go 1.26.
- **An in-app updater.** Pluggable providers, cryptographic verification, an atomic swap and a default UI you can theme or replace. See the [updater guide](/guides/updater/).
- **Packaging and signing.** Guides for [code signing](/guides/build/signing/) on every platform, [installers](/guides/installers/) and [MSIX packaging](/guides/build/msix/).
- **Server builds.** Run the same application and services as an HTTP server without a native window.

## Highlights since the beta

The beta post went out on 2 August, and the changelog has been busy since then. A few things I want to call out:

- **A new documentation site.** The docs now run on [M-Press](https://m-press.me), our own documentation tool, and the translations were completed and brought across from the community. The English documentation remains the source of truth.
- **A tougher Windows experience.** Wails now recovers from WebView2 process failures instead of leaving you with a blank window, and a number of WebView2 COM lifetime issues were fixed along the way.
- **Global shortcuts that respect your keyboard.** On macOS, letter accelerators now follow the labelled key on layouts such as AZERTY, QWERTZ and Dvorak.
- **Updater hardening.** Fixes for swapping binaries across filesystems on Linux and for AppImage updates, among others.
- **Setup and tooling.** The `wails3 setup` wizard was completed, and `wails3 dev` now stops its background processes when your app exits.
- **Native Liquid Glass on macOS.** Support for the Liquid Glass effect on macOS 26 and later, with a fallback for earlier versions.

## Moving from v2

v3 is a new major version, and migrating is a real port rather than a version bump. The big shifts are the application and window lifecycle, services in place of context-bound bindings, direct application and window APIs in place of the v2 runtime package, and regenerated frontend bindings.

The [v2-to-v3 migration guide](/migration/v2-to-v3/) walks through all of it, with before-and-after examples. Please port deliberately, test as you go, and keep v2 in place until your v3 application is ready.

**Wails v2 isn't going away.** It continues to receive critical updates. For new projects, use v3.0.0. The two CLIs are separate binaries (`wails` and `wails3`) with different module import paths, so you can have both on one machine.

## Install

Install the v3.0.0 CLI with:

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0
```

Then run the guided setup wizard, which checks your environment and helps with the dependencies Wails needs:

```sh
wails3 setup
```

![The Wails v3 setup wizard in dark mode](/assets/screenshots/wails3-setup-wizard-dark.png)

And create your first project:

```sh
wails3 init
```

- Documentation: <https://v3.wails.io/>
- Migration guide: <https://v3.wails.io/migration/v2-to-v3/>
- Release notes: [Wails v3.0.0 on GitHub](https://github.com/wailsapp/wails/releases/tag/v3.0.0)

## Support and releases

Wails has grown far beyond what I expected, and it has been hard to keep up with the volume of issues and pull requests. Releases have sometimes been ad hoc, and fixes and contributions have waited too long to reach users. We’ve put systems in place to automate more of the groundwork in triage and release preparation, so we can respond more quickly to problems and spend more time shipping improvements and new features.

After v3.0.0, we’re moving to **nightly pre-releases** with dated tags such as `v3.1.0-nightly.20261026`, alongside **weekly minor releases** such as `v3.1.0` and `v3.2.0`. Nightlies make recent fixes available sooner, while the weekly releases give you a regular cadence of tested updates.

If you find a reproducible bug, please report it with the output of `wails3 doctor` and a minimal example. If you want to propose a new capability or a change to public behaviour, start a draft WEP (Wails Enhancement Proposal) PR rather than a feature-request issue.

## Translations

To everyone who helped translate the Wails v2 documentation: thank you, from the bottom of our hearts. You gave your time so that developers could read the docs in their own language, and that mattered to a lot of people. We are deeply grateful.

We found Crowdin particularly difficult to manage, and we have stopped using it. From v3 onwards, the documentation is translated automatically, and the translations live in the repository next to the English pages.

If you would like a translation in a language we don't have yet, or you have spotted a mistake in one that we do, please open an issue and tell us. We would love to hear from you, and we are thankful for every report.

## Credit & Thanks

Wails v3 exists because of the people who tested incomplete builds, reported hard bugs, translated documentation, answered questions and contributed code. Thank you. The full list of contributors lives on the [credits page](/credits/), and I'd encourage you to go and look at it.

First, a huge thank you to [Atterpac](https://github.com/atterpac), a Wails maintainer, for the development tooling, Linux desktop integration and community support that have helped bring v3 to this point.

Special thanks to [Fabio Massaioli (fbbdev)](https://github.com/fbbdev) for creating such an amazing static source analyser and binding generator. It is a huge part of what makes v3 work: ordinary Go services become JavaScript and TypeScript APIs, with types and documentation carried through to your frontend.

Thank you also to [Travis McLane](https://github.com/tmclane) for his extensive early platform work; [Stffabi](https://github.com/stffabi) for the asset server and native platform foundations; [Zach Botterman (popaprozac)](https://github.com/popaprozac) for notifications, Dock integration and other desktop features; [Ian VanSchooten (IanVS)](https://github.com/IanVS) for typed events and runtime improvements; and [Lyimmi](https://github.com/Lyimmi) for Linux work. These contributions remain part of v3, alongside the work of everyone who has helped build and refine it.

A very heartfelt thank you to the sponsors who have carried Wails through this long transition. Your support gave us the time to work on the architecture, tooling, testing and documentation that got us here. Special thanks also to JetBrains and ZSA, and to [CodeRabbit](https://www.coderabbit.ai/oss) for supporting the project with code reviews.

Finally, a special thank you to Valerio Tesei for all of his support of the project over the years.

## Final Words

We already have a lot of work waiting in the wings. There are some exciting developments across the tech world, and we’re ready to explore new paradigms and bring whole new features to Wails. First, we need v3 to be a stable foundation that you can rely on. Getting that right comes before rolling out what’s next.

I'm really proud of v3. I hope it gives you as much pleasure to build with as it has given us to make. Try it, build something, and tell us what you think. Come and say hello on [Discord](https://discord.gg/JDdSxwjhGf), [GitHub Discussions](https://github.com/wailsapp/wails/discussions) or the [subreddit](https://www.reddit.com/r/wails/).

Enjoy!

&dash;&nbsp;Lea

PS: If you or your company find Wails useful, please consider [sponsoring the project](https://github.com/sponsors/leaanthony). Thanks!
