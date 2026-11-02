---
title: "Wails v3.1.0: More Go, fewer build hurdles"
description: "Experimental CGO-free desktop builds and Wake: a simpler way to build, package and sign Wails applications."
authors: ["leaanthony"]
tags: ["wails", "v3", "release"]
date: "2026-11-01"
slug: "blog/wails-v3-1-released"
sourcePath: "blog/wails-v3-1-released.md"
---

It’s been a busy week at Wails: CGO-free builds, easier project setup, and less friction around packaging and signing.

Wails v3.1.0 starts our weekly feature releases, getting new capabilities into your hands regularly. Bug fixes will continue through nightly patch updates as they’re ready.

## This week’s fixes and optimisations

This week brought **100 fixes, 21 performance improvements and four security hardening changes** across the framework, CLI, build tooling and CI. Highlights include more reliable mixed-DPI window positioning, Linux tray and dark-mode fixes, leaner asset serving, and lower memory use during MSIX packaging. [See a few of the changes below](#a-few-of-this-weeks-highlights).

## Experimental CGO-free builds

Purego lets Go call native system libraries without CGo, Go’s usual bridge to C code. We started exploring it when it first appeared, but maintaining two backends was a challenge. With purego maturing and generative AI making that maintenance more manageable, we’re ready to share experimental macOS and Linux backends.

You can now build and cross-compile desktop apps without a C compiler:

```sh
WAILS_PUREGO=1 wails3 build
```

After building your frontend assets, you can also build directly with Go:

```sh
CGO_ENABLED=0 go build -tags purego,production .
```

The macOS backend uses AppKit and WebKit; Linux supports GTK4/WebKitGTK 6.0 and GTK3. Users still need native runtime libraries, and packaging and signing have separate requirements. CGo remains the default while we gather feedback.

[Implementation and platform details](https://github.com/wailsapp/wails/pull/6521)

## Why Wake?

Wails v2 hid too much of the build process. In v3, Taskfiles opened everything up, but most people want a few customisations without maintaining a build system.

Wake puts common choices in `wails.yaml`: applications, companion CLI tools, servers, targets and packaging. Task remains our core build system and is available when you need its flexibility.

Wake now supports live reload, incremental builds, cached outputs and parallel packaging. Already have binaries? `--skip-build` packages them without Go or frontend tools installed.

Enable it with `WAILS_EXP_USE_WAKE=1`.

### A guided start

With Wake enabled, `wails3 init -ui` walks you through backend, target and package choices, checks your tools, and saves the project configuration.

@image{src="/assets/screenshots/wake/application.png" alt="The Wake project wizard showing desktop application and backend choices, including experimental Pure Go." expand}

Choose package formats and signing options for each target:

@image{src="/assets/screenshots/wake/packages.png" alt="Wake package selection for Linux, showing AppImage, DEB, RPM, Void and Arch Linux formats and the option to sign release outputs." expand}

The tool check below shows a standard Linux build:

@image{src="/assets/screenshots/wake/tools.png" alt="Wake checking the C compiler, Go, npm and pkg-config for a standard Linux build." expand}

### What can I build from my machine?

Compilation, packaging and signing have different requirements. Here’s what you can do from each host:

@table{header search sort column-separators align=["left", "centre", "centre", "centre", "left"]}
| Operation | Windows host | macOS host | Linux host | What you need |
| Compile desktop apps with purego | ✓ | ✓ | ✓ | Go and the experimental purego backend for macOS/Linux targets |
| Create MSIX packages | ✓ | ✓ | ✓ | Built application; no Windows SDK required |
| Create NSIS installers | ✓ | ✓ | ✓ | Built application and `makensis` |
| Create macOS apps and DMGs | ✓ | ✓ | ✓ | Built application; Go handles package creation |
| Create DEB, RPM, Arch and Void packages | ✓ | ✓ | ✓ | Built application |
| Create AppImages | ✓ | ✓ | ✓ | Target AppImage runtime; Linux GUI runtime tree for cross-packaging |
| Sign Windows packages with an exported identity | ✓ | ✓ | ✓ | Exported RSA PFX certificate and private key |
| Sign Apple apps and DMGs with an exported identity | ✓ | ✓ | ✓ | Exported RSA PKCS#12 identity; see validation note below |
| Sign DEB, RPM and Arch packages | ✓ | ✓ | ✓ | Exported OpenPGP key |
| Submit to Apple notarization and attach tickets | ✓ | ✓ | ✓ | Apple team API key and an appropriately signed application |
| Use Apple’s `codesign` or keychain profiles | – | ✓ | – | macOS and the appropriate signing identity |
@end

CGo builds still need a C toolchain. Custom hooks that modify DMGs require macOS and `hdiutil`.

### Signing and notarization

Wake supports Windows, Apple and Linux package signing, including ad-hoc Apple signatures. Cross-platform signing uses exported identities; certificate-store and keychain identities still need native tools. User-level credentials keep private keys out of project configuration. AppImage and Void signing aren’t included yet.

Apple notarization works through Go with a team API key. We’ve tested submission and ticket attachment from Linux, with Gatekeeper accepting the result. That app was signed using Apple’s [`codesign`](https://developer.apple.com/library/archive/documentation/Security/Conceptual/CodeSigningGuide/Procedures/Procedures.html); Apple acceptance of Wake’s Go-generated Developer ID signatures remains unverified.

### Moving to Wake

Migration supports all nine built-in templates. It’s one-way: activating `wails.yaml` replaces Taskfile routing. Unsupported custom commands and overrides are rejected for review.

[Wake features, requirements and migration details](https://github.com/wailsapp/wails/pull/6531)

## A few of this week’s highlights

- **Windows land where they should on mixed-DPI displays.** Screen bounds and positions are now handled correctly on macOS and Linux setups with different display scales. [#6346](https://github.com/wailsapp/wails/pull/6346)
- **Linux tray and theme detection work more reliably.** Trays created after the app starts are no longer dropped, tooltips are populated, and dark-mode detection reads the desktop’s appearance setting from startup.
- **Go event callbacks keep their order per listener.** Each listener runs serially in emit order; different listeners can run concurrently, and concurrent emits have no defined order. Event hooks are snapshotted under a lock to avoid a data race. [#6275](https://github.com/wailsapp/wails/pull/6275), [#6276](https://github.com/wailsapp/wails/pull/6276)
- **Build caching handles more real-world projects.** Quoted or escaped spaces no longer cause Go output cache misses, and concurrent cache saves are serialised to prevent corrupted files. [#6495](https://github.com/wailsapp/wails/pull/6495), [#6488](https://github.com/wailsapp/wails/pull/6488)
- **Asset serving does less unnecessary work.** Known file extensions skip MIME sniffing, and disabled debug logging no longer prepares messages that will be discarded. Both changes reduce work on the request path. [#6414](https://github.com/wailsapp/wails/pull/6414), [#6420](https://github.com/wailsapp/wails/pull/6420)
- **MSIX packaging avoids binary-sized buffers.** Executables are streamed into the package instead of being read entirely into memory. A representative 16 MiB copy previously allocated about 16.8 MB per operation. [#6472](https://github.com/wailsapp/wails/pull/6472)

## Give it a try

Both features are experimental. Try a cross-platform build, package your app, and tell us where you hit friction. Thanks to everyone testing, reporting issues and contributing.

A special shout-out to beesophie from our Discord community, who built [Cook](https://github.com/SimTXNG/cook), an incredible build system that showed us some of the gaps in our own. Seeing what Cook made possible helped us recognise where we could make building with Wails easier, and was a real inspiration for this work.

Thanks Beesophie!
