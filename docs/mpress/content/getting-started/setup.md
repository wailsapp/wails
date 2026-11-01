---
title: "Setup"
slug: "getting-started/setup"
sourcePath: "getting-started/setup.md"
---

@note{type="caution" title="Experimental"}
The setup wizard is new and primarily tested on Linux. If you encounter issues, please [report them](https://github.com/wailsapp/wails/issues/4904) and follow the [manual installation steps](/getting-started/installation/#platform-specific-dependencies) instead.

@end

## Quick Start

```bash
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

# 2. Run the setup wizard
wails3 setup
```

The wizard opens in your browser and guides you through dependency checks, project defaults, and optional cross-platform build setup. The instructions above use the default setup flow; the experimental Wake project flow is described below.

Then you're ready for your first project:

```bash
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
```

## What It Does

- **Checks dependencies** - Verifies Go, npm, and platform tooling
- **Configures defaults** - Author info, bundle ID prefix, preferred templates
- **Cross-platform builds** - Optional Docker setup for building from any host
- **Code signing** - Optional setup for macOS, Windows, and Linux

Configuration is saved to `$XDG_CONFIG_HOME/wails/defaults.yaml`, falling back to `~/.config/wails/defaults.yaml`, and used by `wails3 init`. This is separate from a project’s `build/config.yml`.

## Configure a Wake project

Run setup from your project directory with the [Wake experiment enabled](/experimental/yaml-builds/#enable-yaml-builds):

```bash
export WAILS_EXP_USE_WAKE=1
wails3 setup
```

In PowerShell:

```powershell
$env:WAILS_EXP_USE_WAKE = '1'
wails3 setup
```

This opens a project-focused flow:

1. **Project** — enter project metadata when no native manifest exists.
2. **Application** — choose GUI, CLI or server, a build backend, or already-built executables. Pure Go is an explicit experimental choice.
3. **Targets** — select operating systems and architectures. iOS and Android remain experimental; configure their signing settings in YAML.
4. **Packages** — choose package formats and whether to sign them.
5. **Tools** — review requirements for the selected build and packaging path.
6. **Signing** — configure signing when selected, using the existing signing controls.
7. **Save** — write the project configuration and show the command to run.

Tool requirements depend on your choices. Pure-Go desktop compilation does not need Docker or a C compiler. CGo cross-builds need a compatible compiler and target libraries, or a supported container toolchain. NSIS still needs `makensis`; cross-host GUI AppImages need a supplied Linux runtime tree. Exported signing identities can use Wails' Go signers, while native keychain or certificate-store identities retain platform tooling. See the [packaging and signing requirements](/experimental/yaml-builds/#package-and-sign-for-distribution).

Already-built executables skip Go compilation, frontend installation/builds, bindings and UPX. They still need the selected package/signing tools and valid input binaries. See [prebuilt input paths](/experimental/yaml-builds/#use-prebuilt-desktop-executables). This choice applies to the command shown after saving; it is not a permanent YAML build setting.

Save writes `wails.yaml`, updates the application kind, frontend enablement and selected backend settings, and replaces the dedicated `profiles.setup` target/package/signing selection. Other profiles, frontend commands, companion executables, hooks and unrelated project settings are preserved. Backend settings are project-wide, so other profiles that inherit them also use the selected backend. New project metadata stays in the wizard until Save. Missing tools are reported so you can install them before building; saving valid configuration does not install tools or run a build.

Signed builds reuse missing signing settings from your user-level defaults, including certificate files outside the project. Explicit project settings take precedence. Setup does not copy private keys, passwords or global credential paths into `wails.yaml`. For globally configured Windows and Linux keys, provide the password through the existing credential store or `WAILS_WINDOWS_CERT_PASSWORD` / `WAILS_PGP_PASSWORD`. For an encrypted Apple identity, set `darwin.signing.credential` to the name of an environment variable containing its password.

The saved profile runs with `wails3 build setup`; the already-built option adds `--skip-build`.

Keep Wake enabled when running the displayed command. An active `wails.yaml` never falls back to Taskfiles when the flag is unset. Setup does not translate existing Taskfiles or legacy metadata: use the [reviewed migration workflow](/experimental/yaml-builds/#migrate-an-existing-project) before configuring those projects with Wake.

## Sub-Commands

```bash
wails3 signing      # Configure code signing
wails3 entitlements # Configure macOS entitlements
```

## Having Issues?

1. Run `wails3 doctor` to diagnose problems
2. Follow the [manual installation steps](/getting-started/installation/#platform-specific-dependencies)
3. [Report the issue](https://github.com/wailsapp/wails/issues/4904) with your `wails3 doctor` output
