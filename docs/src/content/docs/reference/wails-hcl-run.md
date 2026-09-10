---
title: HCL run configuration (proposed)
description: Proposed local run defaults, platform overrides and launch behaviour for wails3 run.
---

## Status

This page specifies a proposed extension. The HCL build-system branch does not
yet implement top-level `wails3 run`, `run` blocks or nested `target.run` blocks.
The examples below are specification examples and are not accepted by the current
parser. Existing supported fields are listed in the
[`wails.hcl` field reference](/reference/wails-hcl-fields).

While HCL is experimental, the command follows the existing
[`WAILS_EXP_USE_WAKE=1` opt-in](/experimental/hcl-builds#enable-hcl-builds).
Enabling that flag alone does not implement this extension.

## Command contract

`wails3 run` discovers the nearest `wails.hcl`, resolves configuration for one
target, builds through the existing cached pipeline and launches the result.
Discovery searches from the working directory upwards, checking for a manifest
before stopping at a `go.mod` boundary. The manifest directory is the project
root and the desktop application's working directory.

The default target is the host operating system and architecture. Mobile targets
require explicit selection. Cross-compilation does not imply that a desktop
executable can run on the host; incompatible desktop targets fail before building.
The exact mobile target and device-selection CLI syntax remains to be specified.

The build uses embedded assets and the production build path, which supplies
the `production` tag automatically. `wails3 dev` continues to own live reload;
`dev.tags` does not contribute to `wails3 run`.

```shell
# Proposed commands; not implemented yet.
wails3 run
wails3 run --tags extra_tag
wails3 run -- --example-mode other
wails3 run --plan
```

`--plan` resolves and displays the selected target, effective build tags,
executable or installable artefact, arguments and working directory without
building, installing or launching. Environment overrides must be represented
without exposing sensitive values.

## Fields

The optional `run` block supplies shared launch defaults. An optional `run` block
inside an existing `target "<os>"` block supplies platform-specific defaults.

| Field | Type | Default | Meaning |
| --- | --- | --- | --- |
| `run.tags` | list(string) | `[]` | Additional Go tags used only by run builds. |
| `run.args` | list(string) | `[]` | Arguments passed to the application. |
| `run.environment` | map(string) | `{}` | Runtime environment overrides, separate from compiler environment. |
| `target["<os>"].run.tags` | list(string) | `[]` | Additional run tags for the selected platform. |
| `target["<os>"].run.args` | list(string) | inherited | Replaces shared arguments when explicitly set; `[]` clears them. |
| `target["<os>"].run.environment` | map(string) | `{}` | Overrides shared runtime environment values by key. |

Arguments are passed as an argument vector, without shell expansion. Environment
values are literal strings under the existing HCL literal-only contract.

## Platform example

This fragment accompanies the normal project configuration. `example_private_api`
is a placeholder for the actual private API tag required by the example.

```hcl
build {
  tags = ["devtools"]
}

run {
  args = ["--example-mode"]
}

target "darwin" {
  tags = ["example_private_api"]

  run {
    environment = {
      EXAMPLE_WINDOW_MODE = "native"
    }
  }
}

target "linux" {
  run {
    environment = {
      GDK_BACKEND = "wayland"
    }
  }
}

target "windows" {
  run {
    args = ["--example-mode", "--windows-option"]
  }
}
```

The Linux environment setting is specific to this example and requires a Wayland
session; it is not a Wails default. The private API tag is selected only on macOS.
Existing `target.tags` applies to production builds, including run builds.
`target.run.tags` restricts an additional tag to run builds.

## Resolution rules

| Setting | Resolution, from lowest to highest precedence |
| --- | --- |
| Tags | Combine `build.tags`, selected `target.tags`, automatic `production`, `run.tags`, selected `target.run.tags` and CLI `--tags`; retain first occurrence of each tag. |
| Runtime environment | Inherit the launching shell, overlay `run.environment`, then overlay selected `target.run.environment`. An empty value sets an empty string. |
| Application arguments | Shared `run.args`, replaced by explicitly set target arguments, then replaced by explicit arguments after `--`. A bare `--` clears configured arguments. |

Binding generation and Go compilation receive the same resolved tags. Build
cache keys include those tags. Changing only application arguments or runtime
environment must not invalidate compilation; the next launch uses the new values.

Platform blocks customise a platform; their presence does not make other
platforms unsupported. A supported-platform declaration is also required for
platform-exclusive examples, with early rejection of unsupported targets. Its
field name and syntax remain unspecified and must be settled before implementation.

## Platform launch behaviour

| Platform | Required behaviour |
| --- | --- |
| Linux and Windows | Launch the executable with connected standard input/output, termination handling and application exit-status propagation. |
| macOS | Assemble and launch an app bundle where required, using the existing metadata, signing and entitlement configuration. Verify argument delivery, lifecycle and exit reporting through the bundle launcher. |
| Android | Build an APK, select a device or emulator, install and launch the application. Report build, installation and launch failures distinctly. |
| iOS | Build an app bundle, select a simulator or device, satisfy signing and provisioning requirements, install and launch. Device builds require the appropriate credentials and Apple toolchain. |

Mobile launchers validate whether configured arguments and environment overrides
are supported by the selected device or simulator transport. Unsupported settings
produce actionable errors before deployment; they must not be silently ignored.
Desktop process exit semantics must not be claimed for mobile platforms without
an explicit lifecycle contract and platform verification.

## Implementation acceptance

- Go-only examples require no frontend toolchain. An example-local manifest can
  build using the repository's parent `go.mod` without creating a nested module.
- Tests cover tag-dependent source selection, identical binding/compiler tags,
  platform isolation, cache invalidation and unchanged-build cache reuse.
- Tests cover argument replacement, empty overrides, environment precedence,
  working directory, failed builds, cancellation and desktop exit status.
- `--plan` exposes the effective selection without build or deployment effects.
- Native acceptance covers Windows execution, macOS bundle launching, Android
  device/emulator deployment and iOS simulator/signed-device deployment.
- Schema validation, generated field documentation and CLI help are updated
  together when the capability is implemented.
