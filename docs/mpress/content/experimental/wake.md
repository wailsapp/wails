---
title: "Wake"
description: "The experimental native Wails build system, configured with wails.yaml."
slug: "experimental/wake"
sourcePath: "experimental/wake.md"
---

Wake uses a declarative `wails.yaml` manifest for building, development,
packaging and signing. Enable it with `WAILS_EXP_USE_WAKE=1` and follow the
[YAML build guide](/experimental/yaml-builds/).

For desktop cross-compilation, `CGO_ENABLED=0` selects the experimental Go-only
backend without a C compiler, target SDK, Docker or emulation. Unsigned macOS
app/DMG, Windows MSIX and Linux package-manager formats can be assembled on any
supported desktop host. NSIS still needs a host-native compiler. AppImage
deployment currently requires a matching Linux host/architecture and does not
support dynamically loaded pure-Go GTK libraries. Signing has its own platform
requirements; Arch detached signatures use the native Go signer.

Migration from Taskfiles is **one-way**. `wails3 migrate` analyses legacy
files without executing their commands, writes an inactive proposal, and
activates reviewed YAML only when the migration is complete. Once active,
Wake ignores Taskfiles. Invalid YAML or an unset experimental flag cannot
silently return the project to Taskfile execution.

Unmigrated projects and explicit `wails3 task` commands use the separate
embedded Task runtime. Wake does not implement Taskfile execution, template
expansion or shell-variable evaluation. The former `WAILS_USE_WAKE` selector
is retired.

See [manifest fields](/reference/wails-yaml-fields/) for the generated schema
reference. `wails3 eject` writes resolved YAML; it does not generate Taskfiles.
