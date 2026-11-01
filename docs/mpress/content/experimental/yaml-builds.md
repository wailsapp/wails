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

New to Wake? Start with the [getting-started tutorial](/experimental/wake/).
It covers the experiment's limits, every project and build setup screen, a first
local build, and a reviewed migration. This page covers individual build,
packaging and configuration workflows in more detail.

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

To create and configure a project in the browser, run `wails3 init -ui` with the opt-in enabled and follow the [Wake tutorial](/experimental/wake/). Project choices live in `wails.yaml`; edit that file to change an existing project. `wails3 setup` configures shared preferences, tools and signing credentials.

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

## Cross-compile desktop applications

For the experimental pure-Go desktop backend, set `CGO_ENABLED: "0"` in the
build environment. Wake selects the `purego` tag for macOS and Linux GUI
programs and uses Go directly on Windows, macOS and Linux hosts. Building
desktop executables does not require Docker, emulation, a C compiler or a
target-platform SDK in this mode. Application dependencies must also support
`CGO_ENABLED=0`.

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
build:
  output: bin
  environment:
    CGO_ENABLED: "0"
```

```bash
wails3 build --targets windows/amd64,windows/arm64,linux/amd64,linux/arm64,darwin/amd64,darwin/arm64
wails3 build --targets darwin/universal --formats dmg
```

The second command combines both macOS architectures and creates an unsigned
DMG using Go. It can run on any supported desktop host.

The resolved mode is selected independently for each executable. A declared
`build.environment.CGO_ENABLED`, target environment or binary environment
overrides the invoking terminal's value. Otherwise, an explicitly selected
`purego` tag, a CLI/server program or a Windows target defaults to `0`; macOS
and Linux GUI programs default to `1`. `--plan --json` reports each executable's
resolved compiler, CGo mode and tags. Cross-compilation creates target files;
running an incompatible desktop target still requires the appropriate host.

With `CGO_ENABLED=1`, Wake uses a compatible native toolchain or the configured
`CC` cross-compiler first, then Zig for a Windows cross-build when available,
then its existing Docker/Podman fallback where supported. An explicit `CC`
must supply the target compiler, libraries and sysroot; Zig alone does not
provide a Linux GTK/WebKit sysroot. Packaging and signing may need additional
tools as described below. Mobile targets require CGo and their platform SDKs.

Go-only compilation does not bundle the destination's system libraries.
Linux GUI applications, including GUI companion executables, still require the
GTK/WebKitGTK runtime selected by their build tags. Package dependencies or an
appropriate deployment environment must provide it. CLI/server-only packages
do not add those GUI dependencies by default; explicit dependency overrides
remain available.

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

## Package and sign for distribution

### Use prebuilt desktop executables

Use `--skip-build` when another build system or CI job has already produced the
executables:

```bash
wails3 build --targets linux/amd64 --formats deb,rpm --skip-build
wails3 run --skip-build -- --example-mode
```

Place each executable in `build.output` (default `bin`), using
`project.binary_name` for the main executable and each `binaries` entry's name
for companions. Windows names end in `.exe`. A multi-target invocation reads
each target's files from an `<os>-<arch>` subdirectory. For a macOS GUI, Wake
reads the raw executable at `bin/<name>` if present, otherwise
`bin/<name>.app/Contents/MacOS/<name>`. A `darwin/universal` input must contain
both amd64 and arm64 executable slices.

Wake validates file type and target architecture, copies the executables into
its generated workspace and tracks their contents for package caching. Missing,
incompatible or symlinked executable inputs fail instead of rebuilding or
restoring an old input. Supplied executables are not overwritten or chmodded.
Changes to any executable invalidate packages that use it.

This mode skips frontend installation/builds, bindings, Go compilation, UPX and
the `before_build`/`after_build` hooks. It still generates package assets,
assembles macOS app bundles and runs configured packaging/signing hooks.
Package and signing tools remain subject to the table below; Go and frontend
build tools are not required. Package metadata comes from the current manifest,
including `version_file` when configured; no version is inferred from the binary.
Use matching build and package settings, including GUI versus CLI kind.

`wails3 package --skip-build` is also supported. With `wails3 run`, launch
arguments and environment still apply, and the executable must run on the host.
This mode requires an active YAML project and supports production desktop
targets only; it does not apply to development sessions or mobile packaging.

Unsigned package creation has a different host/tool requirement from signing:

IPA payloads and macOS notarization ZIPs are written in Go, retaining bundle
paths, executable permissions and relative symlinks. Neither operation needs
an external `zip` or `ditto` executable. These archives contain the application
files rather than Finder metadata or resource forks; macOS signing already
removes that metadata because it invalidates signed bundles.

| Output | Creation hosts and tools | Signing |
| --- | --- | --- |
| Windows MSIX | Windows, macOS or Linux; Go writer, no Windows SDK required. | Go signing with an exported RSA PFX on any desktop host; certificate-store/provider identities require Windows SignTool. |
| Windows NSIS installer | Windows, macOS or Linux; host-native `makensis` required. | Go PE signing with an exported RSA PFX on any desktop host; certificate-store/provider identities require Windows SignTool. |
| macOS app and DMG | Windows, macOS or Linux; Go assembly, universal-binary merge and portable DMG writer. | Go signing with an exported RSA PKCS#12 identity, or ad-hoc signing; keychain identities retain macOS tools. API-key notarization and ticket stapling run in Go. |
| Linux DEB, RPM, Arch and Void | Windows, macOS or Linux; native Go package writers. | Exported OpenPGP keys use the native DEB/RPM/Arch signer. Legacy DEB/RPM key IDs require external tools. Void output is unsigned. |
| Linux AppImage | Windows, macOS or Linux; Go deployment and SquashFS writer, plus a target AppImage runtime. GUI cross-packaging requires a supplied Linux runtime tree. | No built-in AppImage signing flow. |

For signed DMGs, Wake normally signs a private copy of the application before
the Go writer compresses the image, then signs the image. This avoids mounting,
expanding and recompressing the DMG with `hdiutil`. Signing runs on every request;
the signed package is not reused from the build cache. Projects with
`after_package` or `before_sign` hooks retain the original package-then-sign
flow, including `hdiutil`, so edits those hooks make to the unsigned image are
preserved. That hook-compatible image-editing path requires macOS and `hdiutil`.

For portable Apple signing, set `darwin.signing.identity: "-"` for an ad-hoc
signature, or use `darwin.signing.certificate` for an exported RSA PKCS#12
identity. `darwin.signing.credential` names its password environment variable.
Certificate signatures use the shared Go PFX/CMS implementation; no private key
is exported from a keychain automatically. Thin and universal executables,
nested app bundles, frameworks and DMGs have passed Apple's strict native
signature checks, including tamper rejection, using generated test identities.

For notarization, an exported Developer ID identity also needs a compatible
RFC 3161 `darwin.signing.timestamp_server`. The shared verifier accepts its
supported SHA-2 timestamp profiles and rejects unsupported legacy algorithms.
An ad-hoc signature cannot be notarized. Native certificate-format validation
is separate from real Developer ID trust and Apple's acceptance of a submission.

The Go stapler authenticates Apple's ticket signature and verifies that it
covers every architecture before attaching it to a private copy. It uses the
observed ticket-delivery protocol, which Apple does not publish as a supported
API. It does not replace Gatekeeper's current revocation checks. App and signed
DMG outputs have passed native `stapler validate` and strict signature checks.

Windows MSIX archives are created directly in Go; creating an unsigned package
does not require the Windows SDK or Docker. Exported RSA PFX identities sign
Windows executables and individual MSIX packages directly in Go on every desktop
host. Set `windows.signing.certificate` to the project-relative PFX and
`windows.signing.credential` to its password environment variable. AES/PBKDF2 and
legacy 3DES exports are supported; RC2 exports require re-exporting or the native
Windows SignTool fallback. Certificate-store/provider identities and MSI signing
require Windows SignTool. The default timestamp service is DigiCert; timestamps
are verified before publishing signed output. See the
[signing guide](/guides/build/signing/) for format limits.
NSIS uses the external NSIS compiler and
accepts literal definitions through `packages.nsis.defines`. On macOS, Wails
supplies a UTF-8 locale only to the NSIS compiler child when its inherited
locale is not UTF-8; terminal and application environments are unchanged.
The stock NSIS installer omits the WebView2 bootstrapper when every packaged
executable is a CLI or server. A GUI primary or companion retains WebView2.
Custom NSIS scripts retain the bootstrapper for compatibility with their own
runtime-installation logic.

Linux supports AppImage, DEB, RPM, Arch Linux and Void packages. Wake assembles
AppImages directly in Go, including GTK/WebKit deployment, SquashFS compression
and zsync generation. It does not invoke linuxdeploy, the GTK plugin, appimagetool
or mksquashfs. A pinned, checksum-verified AppImage type-2 runtime is downloaded
for the target architecture; `packages.appimage.runtime_path` selects a
project-relative runtime ELF instead.

On Linux, the default library source is the host root filesystem. To package a
GUI for another architecture or from Windows/macOS, set
`packages.appimage.runtime_root` to a project-relative target Linux root tree
containing its libraries, compiled schemas and module caches. This is the whole
root layout, including `usr/`, rather than the `usr/` directory itself. Supply a
runtime from the oldest Linux distribution you intend to support. Wake reads
ELF imports and explicit GTK/WebKit profiles without executing target programs.
It bundles the selected closure, not glibc or the dynamic loader. Library and
resource content hashes invalidate the package cache after runtime updates.

Pure-Go GUI builds are supported. Wake selects GTK 3 for GUI executables built
with `gtk3`, and GTK 4 otherwise, including GUI companions of a CLI/server.
All GUI executables in an image must use one GTK version. For prebuilt binaries,
`packages.appimage.gtk_version` can explicitly select `3` or `4`; imported CGo
libraries must agree with that choice.

Set `packages.appimage.update_info` to request a zsync file and
`packages.appimage.checksum: true` for a SHA256 file. Sidecars are restored from
the package cache and listed in the artifact receipt alongside the AppImage.
Windows and macOS packaging have been verified with a supplied GTK/WebKit tree;
the resulting pure-Go GUI images passed runtime checks on Linux. In one matched
CGo GUI build, the native path took 8.77 seconds versus 26.19 seconds for
linuxdeploy. The native gzip image was 133 MiB versus 114 MiB with the old zstd
path, so this measures the whole packaging path rather than compression alone.

Arch packages preserve prereleases in native package metadata, for example
`1.2.3-beta.2` becomes `1.2.3beta.2-1`, which sorts before the final release.
Void dependencies accept bare package names such as `gtk4`, normalised to
`gtk4>=0`; explicit version comparisons and complete `name-version_revision`
values are preserved. Use Void package names rather than library filenames.

DEB, RPM and Arch signing can use exported armored or binary OpenPGP private
keys on any supported desktop host, without GnuPG, `dpkg-sig` or `rpmsign`.
Configure `linux.signing.certificate` with the key file and `credential` with
the name of its passphrase environment variable. DEB and RPM embed signatures;
Arch publishes a detached signature. Private keys are read only by the signing
step and are not copied into package assets or reusable cache artifacts.

Existing DEB/RPM configurations containing a bare GnuPG key identifier keep
their external signer. To remove that dependency, export the secret key and
configure its path instead. Use a path such as `./signing/release.asc`; missing,
unreadable or invalid key files fail rather than falling back to a keyring.
DEB signatures retain the `dpkg-sig` package-signature format; this does not
replace signing an apt repository's release metadata.

To create an Arch package and detached signature, use an exported armored PGP
private-key file and optionally name a passphrase environment variable:

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
build:
  environment:
    CGO_ENABLED: "0"
linux:
  signing:
    certificate: signing/private-key.asc
    credential: ARCH_SIGNING_PASSPHRASE
profiles:
  release:
    targets:
      linux/amd64:
        formats: [archlinux]
        sign: true
```

Set the named environment variable through your local environment or CI secret
store, then run `wails3 build release`. Omit `credential` for an unencrypted
key. The Go signer runs on every supported desktop host and does not require
GPG. It publishes the original `.pkg.tar.zst` and a sibling `.pkg.tar.zst.sig`;
both are recorded separately in the artifact receipt. Key contents are not
staged or stored in the build cache. Signatures and their publication are never
reused from cache. Failed signing preserves the previous published pair, and
publication rolls both files back if installing either file fails. This does
not make two filesystem renames atomic to concurrent readers.

macOS app bundles include companion programs. Set
`darwin.microphone_usage_description` to add the usage text and, when no custom
entitlements file is supplied, the audio-input entitlement. Custom plist and
entitlements sources remain unchanged. An exported RSA PKCS#12 identity uses
Go signing on every desktop host; `identity: "-"` selects portable ad-hoc
signing. Keychain identities continue to use `codesign` on macOS. Ticket
stapling runs in Go. Notarization submission can use an existing macOS notary
keychain profile, or App Store Connect team API credentials:

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
darwin:
  signing:
    certificate: signing/developer-id.p12
    credential: APPLE_CERT_PASSWORD
    timestamp_server: https://timestamp.digicert.com
  notarization:
    api_key: signing/AuthKey.p8
    key_id: ABCDEFGHIJ
    issuer_id: 00000000-0000-0000-0000-000000000000
```

API mode submits, uploads and polls from Go without `notarytool`. Set all three
fields and omit `credential`; keep the private key out of version control.
The key is read only at signing time and is not copied to generated assets or
cached. This mode supports team API keys and archives up to 5 GiB; multipart
uploads and individual API keys are not implemented. Rejections retain the
submission ID and save Apple's diagnostic log to a private temporary file whose
path is reported in the error. A cancelled or rejected submission does not
replace previously published output. After acceptance, Wake retrieves and
verifies Apple's ticket and staples it in Go. With an exported identity and
team API credentials, signing, submission and stapling can run on Linux,
Windows or macOS. Keychain profiles still require macOS and `notarytool`.

Portable certificate signing supports RSA identities in AES/PBKDF2 or legacy
3DES PKCS#12 files. RC2-encrypted exports must be re-exported with a supported
cipher. Notarized portable signatures require an explicit SHA-2 RFC 3161
timestamp service; the Go verifier rejects legacy SHA-1 timestamp signatures.
The ticket service follows Apple's observed protocol and pins Apple Root G3;
Gatekeeper remains the authority for distribution trust and revocation.

Signed outputs retain their platform extension, for example
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

See the generated [`wails.yaml` field reference](/reference/wails-yaml-fields/)
for every field, type, default and applicability rule.

The [`wails3 run` configuration reference](/reference/wails-yaml-run/) describes
local launch defaults, platform overrides and the no-manifest `go run .` fallback.
