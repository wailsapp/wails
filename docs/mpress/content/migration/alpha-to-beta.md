---
title: "Upgrade from a v3 alpha"
description: "Move an existing Wails v3 alpha project to a pinned beta release"
slug: "migration/alpha-to-beta"
sourcePath: "migration/alpha-to-beta.md"
---

This guide is for existing v3 alpha projects. For Wails v2, use the [v2-to-v3 guide](/migration/v2-to-v3/).

## Before upgrading

Commit or back up your project. Read the [changelog](/changelog/) between your alpha version and the chosen beta: source code, APIs or build configuration may need changes. Check the [desktop compatibility policy](/status/) and your platform requirements.

The commands below use the published `v3.0.0-beta.23` release as an exact-version example, not a recommendation to track the latest release. If choosing another release, verify its CLI, Go module and npm runtime versions and update the commands together. In this example, the npm version matches the Go version without its `v` prefix.

## 1. Update the CLI

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0-beta.23
wails3 version
```

Check that `wails3 version` reports the version you installed. An older binary earlier on `PATH` can hide the new CLI.

## 2. Update the Go module

Run from your project root. Review the dependency changes; do not broadly upgrade unrelated modules.

```sh
go get github.com/wailsapp/wails/v3@v3.0.0-beta.23
go mod tidy
```

## 3. Update the frontend runtime

For projects using npm and a `frontend` directory:

```sh
cd frontend
npm install --save-exact @wailsio/runtime@3.0.0-beta.23
cd ..
```

Keep the lockfile and review its changes. If your frontend uses another package manager or directory, adapt this step while retaining an exact runtime version.

## 4. Regenerate, build and test

From the project root, regenerate bindings from your Go services and build:

```sh
wails3 generate bindings
wails3 build
```

Run the built application and test your workflows on each supported platform you ship. Review and commit the source, generated bindings, module files and frontend lockfile changes together.

## If the upgrade fails

Check the CLI on `PATH`, the module version with `go list -m github.com/wailsapp/wails/v3`, and the installed runtime with `npm --prefix frontend ls @wailsio/runtime`. Regenerate bindings after resolving version mismatches. Do not assume every alpha can upgrade without code changes.

If the problem remains, [report a reproducible issue](https://github.com/wailsapp/wails/issues/new/choose) with your old and new versions, exact error and `wails3 doctor` output. Follow the [security policy](https://github.com/wailsapp/wails/blob/master/SECURITY.md) for vulnerabilities.
