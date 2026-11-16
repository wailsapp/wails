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

The wizard opens in your browser and guides you through dependency checks, global project defaults, optional cross-platform tools and signing credentials. These preferences can be reused across projects.

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

## Create a Wake project

With the Wake experiment enabled, `wails3 init -ui` collects your project-specific build choices and saves them in `wails.yaml` when creation succeeds. Follow the [Wake getting-started tutorial](/experimental/wake/) to choose a build backend, targets and package formats while creating a desktop project. `wails3 setup` keeps its global preferences role.

Existing Wake projects are configured in `wails.yaml`. Use the [YAML build guide](/experimental/yaml-builds/) to change their profiles and packaging settings. Existing Taskfile projects need the [reviewed one-way migration](/experimental/wake/#migrate-an-existing-project) before using Wake; setup does not translate or execute their Taskfiles.

## Reuse signing preferences

Signed Wake builds reuse missing signing settings from your user-level defaults, including certificate files outside the project. Explicit project settings take precedence. Private keys, passwords and global credential paths are not copied into `wails.yaml`.

For globally configured Windows and Linux keys, provide the password through the existing credential store or `WAILS_WINDOWS_CERT_PASSWORD` / `WAILS_PGP_PASSWORD`. For an encrypted Apple identity, set `darwin.signing.credential` to the name of an environment variable containing its password.

## Sub-Commands

```bash
wails3 signing      # Configure code signing
wails3 entitlements # Configure macOS entitlements
```

## Having Issues?

1. Run `wails3 doctor` to diagnose problems
2. Follow the [manual installation steps](/getting-started/installation/#platform-specific-dependencies)
3. [Report the issue](https://github.com/wailsapp/wails/issues/4904) with your `wails3 doctor` output
