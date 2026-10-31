---
title: YAML builds
description: Configure and run the experimental Wails v3 native build system with one wails.yaml file.
sidebar:
  order: 2
  badge:
    text: Experimental
    variant: caution
---

@note{type="caution" title="Experimental preview"}
The accepted first release supports literal configuration only. This experiment
also prototypes bounded lifecycle hooks so their shape can be evaluated before
they are proposed for release. Custom stages, YAML expressions and typed Wails
tool calls remain deferred. Existing projects keep using Taskfiles until you
explicitly activate a reviewed migration.
@end

## Enable YAML builds

Set `WAILS_EXP_USE_WAKE` in every terminal that runs the experimental commands.
For Bash or Zsh:

```bash
export WAILS_EXP_USE_WAKE=1
```

For PowerShell:

```powershell
$env:WAILS_EXP_USE_WAKE = '1'
```

The variable is a presence-based opt-in: even an empty value or `false` enables
it. Unset it (`unset WAILS_EXP_USE_WAKE` in Bash/Zsh, or
`Remove-Item Env:WAILS_EXP_USE_WAKE` in PowerShell) to disable the experiment. Once `wails.yaml` is active, commands require
re-enabling Wake and never fall back to Taskfiles.
Without the variable, YAML build/dev routing and YAML project generation are
disabled, and `migrate`, `config` and `eject` commands are unavailable.
Projects that have not migrated continue through the separate embedded Task runtime.

## Configure a project

With the opt-in enabled, add `wails.yaml` at the project root. Its presence
selects the native build system and makes Wails ignore every Taskfile in that
project. Keep the opt-in set for all commands below. The schema modeline enables
completion and inline validation in editors with YAML language-server support.

```yaml
# yaml-language-server: $schema=https://v3.wails.io/schemas/wails.v3.json
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
frontend:
  directory: frontend
  install: [pnpm, install, --frozen-lockfile]
  build: [pnpm, run, build]
  dev: [pnpm, run, dev]
  output: frontend/dist
build:
  output: bin
  trim_path: true
  strip: true
profiles:
  release:
    targets:
      linux/amd64:
        formats: [deb, rpm, appimage]
```

Run the default host build or the complete named release request:

```bash
wails3 build
wails3 build release
```

For a one-off build, select comma-separated targets and formats without adding
a profile:

```bash
wails3 build --targets linux/amd64 --formats deb,rpm
```

## Build several executables

The project defines the primary program. Set `project.package` to its Go package
(default `.`) and `project.kind` to `gui`, `cli` or `server` (default `gui`).
Additional programs use named entries under `binaries`:

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
binaries:
  notes-tool:
    package: cmd/notes-tool
    kind: cli
    tags: [tools]
    ldflags: [-X, main.Mode=release]
    environment:
      CGO_ENABLED: "0"
profiles:
  release:
    targets:
      linux/amd64: {}
```

Each program gets its own compilation and cache inputs. Imported project
packages and embedded files are tracked; an edit to an independent program
does not rebuild the others. Additional executables are supported on desktop
and included in the default app bundles and installers. User-owned replacement
package files control their own contents.

```bash
wails3 build --binaries notes,notes-tool --jobs 4
wails3 build --binaries notes-tool --skip-frontend --skip-bindings
```

A named profile is a complete request; use these selection overrides with an
anonymous build. `--jobs` bounds execution workers. Selecting only CLI/server
programs omits frontend work.

To read a version without running application code, configure both
`project.version_file` and `project.version_symbol`. The file must contain a
static Go string declaration. The literal `project.version` can then be omitted.
Set `commit_symbol: main.Commit` on the project or a named binary to inject the
repository revision into that executable's string variable. Missing repository
metadata is an error when commit injection is requested.

`upx_level: 1` through `9` optionally compresses release Linux and Windows
executables; zero disables it. UPX must be installed. Development and macOS
builds omit compression. A compilation or compression failure preserves the
previous published program.

## Native package adapters

Windows MSIX archives are created directly in Go; creating an unsigned package
does not require the Windows SDK or Docker. Windows signing still requires a
Windows host and its signing tools. NSIS uses the external NSIS compiler and
accepts literal definitions through `packages.nsis.defines`.

Linux supports AppImage, DEB, RPM, Arch Linux and Void packages. AppImages retain
the external linuxdeploy and GTK plugin boundary. Set
`packages.appimage.update_info` to request a zsync file and
`packages.appimage.checksum: true` for a SHA256 file. Sidecars are restored from
the package cache and listed in the artifact receipt alongside the AppImage.

macOS app bundles include companion programs. Set
`darwin.microphone_usage_description` to add the usage text and, when no custom
entitlements file is supplied, the audio-input entitlement. Custom plist and
entitlements sources remain unchanged. Signing and notarization use the
existing Apple tools and credentials on macOS; notarization requires a real
notary profile. Signed outputs retain their platform extension, for example
`notes.signed.app` and `notes.signed.dmg`.

## Development controls

During `wails3 dev`, enter a command followed by Enter:

| Command | Action |
| --- | --- |
| `r` | Force a rebuild and restart, including an unchanged backend. |
| `b` | Regenerate bindings and rebuild. |
| `f` | Restart the frontend. |
| `c` | Clear the terminal. |
| `q` | Stop the session and its child processes. |
| `h` | Show controls. |

Failed builds keep the current application running. Controls require an
interactive terminal; piped input does not change session behaviour.

## Inspect work before running it

Plan inspection resolves the same immutable graph used by execution. Text is
for people; JSON is a versioned machine-readable contract.

```bash
wails3 build release --plan
wails3 build release --plan --json
```

## Validate configuration without building

Use the dedicated preflight command in CI, an editor task, or before committing
configuration. With no profile it validates the base manifest and every named
profile; with a profile it checks only that complete request:

```bash
wails3 config check
wails3 config check release
```

Syntax, schema, semantic, path, target, format, signing, planner and
configuration-owned host-capability failures use the same diagnostic renderer
as Build and Dev. Diagnostics identify the field and source range, print the
source line and caret, and include a correction hint where one is useful.
Environment values are redacted from diagnostic source excerpts. Host failures
caused only by command-line choices remain ordinary command errors rather than
being misreported as invalid YAML.

## Use project-owned assets

Manifest paths refer to files you own, such as icons, platform manifests,
package templates and maintainer scripts. Wails never edits those sources. It
copies them into a disposable `.wails/` workspace before a tool consumes them.
Final artifacts go under `build.output` and a verified receipt is written to
`.wails/artifacts/receipt.json`.

## Deploy to Android

Production Android profiles produce AABs. For local deployment, Wails builds an
ABI-specific debug APK internally, installs it with ADB, and launches the app:

```bash
wails3 android devices
wails3 android run --device SERIAL
wails3 android run --emulator AVD_NAME
```

Without an explicit selector, an interactive terminal offers connected devices
and configured Android Virtual Devices. Automation with multiple available
targets must provide `--device` or `--emulator`.

Attach an application-scoped log stream and monitor the target until the app
exits or you press Ctrl+C:

```bash
wails3 android run --emulator AVD_NAME --logs
wails3 android run --emulator AVD_NAME --logs --stop-emulator
```

Attached runs follow an application restart and report target disconnects or
offline transitions. Emulators remain running by default for quick iteration.
`--stop-emulator` is explicit, requires `--logs`, and stops only an emulator
target when the attached session ends. If Wails starts a new emulator but boot
readiness fails or is cancelled, it cleans up that process automatically.

## Run bounded lifecycle hooks

The experimental hook surface invokes one project-owned script at each declared
phase. A hook is not a shell command: Wails executes the file directly, from the
project root unless `directory` is set.

On Unix hosts the script must be executable and start with a shebang. Windows
hosts accept `.cmd`, `.bat`, and `.ps1` scripts through fixed interpreter
adapters; Wails does not parse an inline shell command on either platform.

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
hooks:
  before_build:
    script: scripts/generate-version.sh
  after_package:
    script: scripts/check-packages.sh
profiles:
  release:
    targets:
      linux/amd64: {}
```

Supported phases are `before_build`, `after_build`, `before_package`,
`after_package`, `before_sign`, and `after_sign`. Build hooks run once for the
shared project or once per target as appropriate. Package and sign hooks form
barriers around all formats requested for a target while the format operations
inside the barrier remain concurrent.

Hooks receive one Wails-owned environment variable:
`WAILS_HOOK_CONTEXT_FILE`. It names a protected JSON file containing the
versioned invocation context:

```json
{
  "version": 1,
  "phase": "after_build",
  "command": "build",
  "scope": "target",
  "project_dir": "/work/notes",
  "working_directory": "/work/notes",
  "profile": "release",
  "target": { "os": "linux", "arch": "amd64" },
  "output": "/work/notes/bin/notes",
  "declared_outputs": ["/work/notes/generated/version.go"]
}
```

The context contains resolved paths but never copies the inherited process
environment or secrets. Wails creates it atomically with owner-only file
permissions under `.wails/hooks/context/`. It is removed after success or
cancellation. On failure it is retained for diagnosis and its path is included
in the error. A non-zero exit fails the normal Plan step, and cancellation
terminates the hook process group.

Hooks run on every invocation by default. Opt into caching only when the script
has complete declared inputs and outputs:

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
hooks:
  before_build:
    script: scripts/generate-version.sh
    cache: true
    inputs: [version.txt]
    outputs: [generated/version.go]
profiles:
  release:
    targets:
      linux/amd64: {}
```

The script and its executable mode are implicit inputs. Multiple outputs must
share one non-root directory, which is the bounded artifact Wails records and
restores. The output root cannot contain the script or a declared input, and
cached output trees can contain only regular files and directories, never
symbolic links or special files. An always-run hook makes its downstream
operations non-cacheable because undeclared side effects cannot be fingerprinted
safely. A cacheable hook preserves downstream caching through its declared
artifact identity.

## Migrate an existing project

Create and review an inactive draft:

```bash
wails3 migrate
wails3 migrate --dry-run
wails3 migrate --json
```

Conservative argument-free lifecycle tasks that invoke a project-owned script
are translated to hooks. Other reachable custom Taskfile behaviour is reported
as a blocker. After reviewing
`wails.migrated.yaml`, activate it explicitly:

```bash
wails3 migrate --activate
```

Activation creates `wails.yaml` atomically. It does not delete, rename or edit
Taskfiles, `build/config.yml`, scripts or assets. Invalid YAML fails as YAML; it
never falls back to the legacy route.

## Materialise all defaults

To inspect or copy the complete resolved configuration, write an inactive
snapshot:

```bash
wails3 eject
```

This creates `wails.ejected.yaml`. It never rewrites active configuration.

## Compatibility window

`wails3 package` and `wails3 sign` are deprecated aliases for the same planner;
new automation should use `wails3 build` with a profile or formats. Rollout is
release-gated: the first release makes YAML the new-project path, the following
release warns on legacy Taskfile routing, and the next removes Taskfile
execution together with the aliases. Each transition requires the native-host
acceptance matrix for that release to be green.

See the generated [`wails.yaml` field reference](/reference/wails-yaml-fields)
for every field, type, default and applicability rule.

The [`wails3 run` configuration reference](/reference/wails-yaml-run) describes
local launch defaults, platform overrides and the no-manifest `go run .` fallback.
