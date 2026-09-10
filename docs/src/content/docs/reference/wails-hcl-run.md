---
title: HCL run configuration
description: Local run defaults, platform overrides and launch behaviour for wails3 run.
---

## Availability

This reference describes the run extension under development on
`codex/hcl-build-system`. It is not a guarantee of availability in a released CLI.

`wails3 run` builds and launches the application using local manifest defaults.
Manifest-driven execution requires the
[`WAILS_EXP_USE_WAKE=1` opt-in](/experimental/hcl-builds#enable-hcl-builds).
When no manifest is found, the command runs `go run .` in the current directory,
without requiring that opt-in. An invalid or unreadable manifest is an error;
it never selects the fallback. If a manifest exists while HCL is disabled, the
command reports how to enable it.

The generated [`wails.hcl` field reference](/reference/wails-hcl-fields) lists
field types and defaults. Native mobile and macOS bundle acceptance must be
verified on the appropriate host before a release claims those platforms ready.

## Command contract

`wails3 run` discovers the nearest `wails.hcl`, resolves configuration for one
target, builds through the existing cached pipeline and launches the result.
Discovery searches from the working directory upwards, checking for a manifest
before stopping at a `go.mod` boundary. The manifest directory is the project
root and the desktop application's working directory.

The default target is the host operating system and architecture. Mobile targets
require explicit selection. Cross-compilation does not imply that a desktop
executable can run on the host; incompatible desktop targets fail before building.
Use `--target android/arm64 --device <adb-serial>` or `--target android/arm64 --emulator <avd-name>` for Android.
The selected architecture must match the device or emulator. Use `--target ios/arm64 --device <udid>` for an iOS simulator, adding `--destination device` for a physical device. iOS planning and execution require macOS and Xcode.

With HCL, the build uses embedded assets and the production build path, which supplies
the `production` tag automatically. `wails3 dev` continues to own live reload;
`dev.tags` does not contribute to `wails3 run`.

```shell
wails3 run
wails3 run --tags extra_tag
wails3 run -- --example-mode other
wails3 run --plan
```

`--plan` resolves and displays the selected target, effective build tags,
executable or installable artefact, arguments and working directory without
building, installing or launching. The plan lists environment override key names and omits their values.

## CLI options

| Option | Meaning |
| --- | --- |
| `--tags <tags>` | Additional comma-separated Go build tags. |
| `--target <os>/<arch>` | Select one target; defaults to the host OS and architecture. |
| `--plan` | Display resolved settings without building or launching. Cross-platform desktop plans are allowed. |
| `--device <id>` | Android adb serial or iOS simulator/device identifier. |
| `--emulator <name>` | Android Virtual Device to start or reuse; cannot be combined with `--device`. |
| `--destination <kind>` | iOS destination: `simulator` (default) or `device`. |
| `-- <args...>` | Replace configured application arguments. A bare `--` clears them. |

Mobile selectors require a manifest and a matching mobile target.

## Fields

The optional `run` block supplies shared launch defaults. An optional `run` block
inside `target "<os>"` or `target "<os>/<arch>"` supplies platform- or
architecture-specific defaults. Architecture-specific settings take precedence
over OS-wide settings, regardless of the order of blocks in the file.

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

An architecture-specific override can refine the OS-wide defaults:

```hcl
target "linux/arm64" {
  run {
    tags = ["example_arm64"]
    args = []
  }
}
```

Together with the preceding Linux block, this keeps `GDK_BACKEND`, adds the
architecture-specific tag and clears the shared application arguments.

## Resolution rules

| Setting | Resolution, from lowest to highest precedence |
| --- | --- |
| Tags | Combine `build.tags`, selected `target.tags`, automatic `production`, `run.tags`, selected `target.run.tags` and CLI `--tags`; retain first occurrence of each tag. |
| Runtime environment | Inherit the launching shell, overlay `run.environment`, then overlay selected `target.run.environment`. An empty value sets an empty string. |
| Application arguments | Shared `run.args`, replaced by explicitly set target arguments, then replaced by explicit arguments after `--`. A bare `--` clears configured arguments. |

Within the selected target, OS-wide settings are resolved before architecture-specific
settings: tags accumulate, environment values merge by key and explicitly supplied
arguments replace the less-specific list. Run builds also supply `android` or `ios`
for the corresponding mobile target.

Binding generation and Go compilation receive the same resolved tags. Build
cache keys include those tags. Changing only application arguments or runtime
environment does not invalidate compilation; the next launch uses the new values.

Platform blocks customise a platform; their presence does not make other
platforms unsupported. `project.supported_platforms` restricts the platforms accepted by builds and `run`; an omitted or empty list permits all platforms. For example, `supported_platforms = ["darwin"]` inside `project` makes a macOS-only example fail early on Linux.

## Platform launch behaviour

| Platform | Launch behaviour |
| --- | --- |
| Linux and Windows | Launch the executable with connected standard input/output, termination handling and application exit-status propagation. |
| macOS | Assemble an app bundle and execute its `Contents/MacOS` binary with connected standard input/output. Bundle registration, signing and native lifecycle acceptance require macOS verification. |
| Android | Build an APK, select a device or emulator, install and launch the application. Report build, installation and launch failures distinctly. |
| iOS | Build an app bundle, select a simulator or device, satisfy signing and provisioning requirements, install and launch. Device builds require the appropriate credentials and Apple toolchain. |

Android rejects application arguments and runtime environment overrides before deployment. Shared desktop-only settings belong in desktop target blocks when Android is also supported. Android runs attach application logs and stop the application when attachment ends. iOS passes arguments to the selected launcher and environment overrides through its child-environment mechanism. Mobile exit status describes the launcher or log-monitoring result, rather than a portable application exit code.

## Go-only examples and fallback

A Go-only example with pre-existing or embedded assets disables frontend stages:

```hcl
frontend {
  disabled = true
}
```

This skips frontend dependency installation, binding generation and frontend
building. Examples that need generated binding classes can instead set
`frontend.bindings.interfaces = false` using nested blocks.

The no-manifest fallback executes `go run [-tags <tags>] . [application arguments]`
without adding `production` or `devtools`. It follows Go's own exit-status behaviour.
Arguments after `--` go to the application, including flags such as `--help`.
Runtime values shown by `--plan` are limited to environment key names; values are
not printed. The fallback plan displays the Go command and working directory.
