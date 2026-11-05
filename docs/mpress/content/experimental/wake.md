---
title: "Wake"
description: "Create your first Wake project, follow the project creation wizard, build and run it, and migrate an existing Taskfile project."
slug: "experimental/wake"
sourcePath: "experimental/wake.md"
---

Wake is Wails' experimental build system. One `wails.yaml` describes your
application, frontend, target platforms, packages and signing settings. Wake
uses that configuration to plan and run builds, cache unchanged work, and
report the resulting files.

In this tutorial, you will create a small desktop application, configure it in
the browser, build it for your own computer, and run it in development mode.
Start with one local target and an unsigned build. Add cross-compilation and
distribution after that first build works.

Already have an application? Go to [Migrate an existing project](#migrate-an-existing-project).
For individual fields and advanced workflows, use the
[YAML build guide](/experimental/yaml-builds/) and
[manifest reference](/reference/wails-yaml-fields/).

## What to expect from the experiment

@note{type="caution" title="Use a test project or a separate branch"}
Wake is opt-in and its configuration, commands and behaviour may change between
preview builds. It does not yet carry the compatibility guarantees of a stable
build system. Keep a known-working revision of your project and record the CLI
and framework versions you test. Test packages on their destination systems
before distributing them.
@end

Wake and the **Pure Go** backend are separate choices. Enabling Wake does not
automatically enable Pure Go. This lesson uses the standard backend; macOS and
Linux GUI builds therefore need their usual native development tools.

The wizard checks requirements and saves configuration. Saving is not a build,
and it is not proof that every selected target is ready. Missing tools or
credentials can be addressed later, for example on a CI host. Invalid target or
package combinations must be corrected before saving.

Migration is **one-way**. Wake does not execute Taskfiles, expand Task templates,
or translate YAML back into Taskfiles. Unmigrated projects still use the separate
legacy Task runtime. An active `wails.yaml` never silently falls back to it.

## Before you start

You need Go, Node.js with npm for the starter frontend, and the
[platform dependencies](/getting-started/installation/#platform-specific-dependencies)
for your computer. Start from a directory where you want the new project folder
to be created. You also need a browser for the wizard.

Use a CLI and framework revision that include this preview. Installing the
stable `v3.0.0` CLI does not add the Wake project creation screens described here. The
following commands pin the preview from [PR #6531](https://github.com/wailsapp/wails/pull/6531):

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@b1daf95f3
wails3 version
```

Make sure the Go executable directory is on `PATH`; the
[installation guide](/getting-started/installation/) explains this. Later in
the tutorial, you will pin the application's Go dependency to the same revision.
When testing a newer preview, update both revision arguments together.

### Enable Wake in your terminal

In Bash or Zsh:

```bash
export WAILS_EXP_USE_WAKE=1
```

In PowerShell:

```powershell
$env:WAILS_EXP_USE_WAKE = '1'
```

Keep this terminal open for the rest of the tutorial. Set the variable again in
any new terminal or CI job that runs Wake commands. Its **presence** enables
Wake: setting it to `false` or an empty string does not disable it. The old
`WAILS_USE_WAKE` name is ignored.

## Create the application

Run the project creation wizard:

```sh
wails3 init -ui
```

The browser opens a project wizard. With Wake enabled, it collects the source
and build choices in one flow. Nothing is created until the final **Create
Project** button. `wails3 setup` is a separate wizard for global preferences,
tools and signing credentials; it does not configure this project's Wake profile.

| Screen | Choose for this tutorial | What happens next |
| --- | --- | --- |
| **Project — Create a Wails project** | Enter `wake-demo`, then **Continue**. | The project will be created in a new `wake-demo` folder under the displayed directory. |
| **Framework — Choose a framework** | Select **Vanilla**, then **Continue**. | You get a small frontend without a framework-specific learning step. |
| **Language — Language Preference** | Select **JavaScript**, then **Continue**. | The TypeScript-only **Bindings** screen is skipped. |
| **Details — Project details** | Use product name `Wake Demo`, bundle identifier `com.example.wakedemo`, and version `0.1.0`. Click **Continue**. | Your metadata stays in the wizard while you choose the build settings. |

Do not use **Use global defaults** for this lesson: it skips the framework and
language choices. If you choose TypeScript in another project, the **Bindings**
screen lets you choose interfaces or classes for generated Go data bindings.

## Choose the build settings before creation

Continue in the same browser window. The screenshots below show this walkthrough
on Linux. Your host target and required tools will differ on Windows or macOS.

### Application — What are you building?

Keep **Desktop application** and choose **Standard backend**. Click **Continue**.

![Application screen in the project creation wizard with Desktop application and Standard backend selected.](/assets/screenshots/wake/application.png)

This lesson creates a GUI with a web frontend. **Pure Go (experimental)** is a
separate opt-in. It removes the C compiler requirement for supported desktop
compilation, provided your dependencies also build without CGo. It does not
remove target runtime libraries or all packaging and signing tools.

The backend is a **project-wide setting**. Other profiles that inherit those
settings use the selected backend too.

### Targets — Where should it run?

Keep only the target for the computer you are using, then click **Continue**.
The wizard initially selects the host operating system and architecture.

![Targets screen with only linux/amd64 selected for the tutorial host.](/assets/screenshots/wake/targets.png)

Target names combine an operating system and CPU architecture: for example,
`windows/amd64`, `darwin/arm64` for an Apple Silicon Mac, or `linux/amd64`.
`darwin` means macOS. Selecting another target requests an executable for that
machine; it does not make that executable runnable on your current computer.

Leave additional architectures, macOS universal builds and experimental mobile
targets for later. You can add them to `wails.yaml` after the first local build.

### Packages — Choose release outputs

Leave the package format boxes unchecked and **Sign release outputs** off.
Click **Continue**.

![Packages screen with no additional Linux package formats selected and signing off.](/assets/screenshots/wake/packages.png)

With no additional formats selected, Wake produces the target's default output:
a desktop executable, or a macOS application bundle for a macOS GUI. This is
enough to check the application locally. Selecting formats adds distribution
outputs such as DMG, MSIX, NSIS, DEB or AppImage where supported.

Signing is optional for this first local build. A successful unsigned build
does not mean the application is ready for public distribution. For a signed
project, configure reusable credentials with `wails3 setup` before creation, or
configure signing in `wails.yaml` afterwards. The creation wizard uses your
configured credentials; it does not change your global signing preferences.

Wake can use configured certificate files outside the project without copying
private keys or passwords into `wails.yaml`. Explicit project signing settings
take precedence. Follow the [setup signing guidance](/getting-started/setup/)
and [packaging and signing guide](/experimental/yaml-builds/#package-and-sign-for-distribution)
before preparing a public release. Mobile signing is configured in YAML.

### Tools — Tools for your build

Read the requirements for your selections. The list comes from Wake's build
plan; it is not a request to install every tool on every platform.

![Tools screen showing the requirements resolved for the local Linux source build.](/assets/screenshots/wake/tools.png)

For this lesson, expect Go and the frontend tools, plus the native development
tools required by your host's standard GUI backend. Docker is not a blanket
requirement. A cross-build with CGo can require a compiler and target libraries
or a supported container toolchain.

Resolve rows marked **Not ready**, **Needs configuration** or **Update needed**
before trying the local build. Use the displayed instructions where available,
or the [platform installation guide](/getting-started/installation/), then click
**Recheck requirements**. If an installation changed the terminal's `PATH`,
restart `wails3 init -ui` from an updated terminal and repeat your selections.

The wizard can create a valid project even when this computer lacks tools or
signing credentials. That is useful for CI configuration, but for this tutorial
you want a build that can run here. Invalid configuration must be corrected
before creation.

### Create the project

Check your choices: desktop application, standard backend, your host target,
and signing disabled. On **Tools**, click **Create Project** and wait for
**Project choices accepted**. Return to the terminal and wait for the CLI to
finish creating the project; the browser confirmation only accepts your choices.

![Project choices accepted screen showing the planned application path and commands to use after creation succeeds.](/assets/screenshots/wake/save.png)

The wizard creates the source and saves the target, package and signing choices
under `profiles.setup` in `wails.yaml`. It also sets the application kind,
frontend enablement and project-wide backend. You can close the browser and
return to the terminal. Wait for project creation to finish, then run:

```sh
cd wake-demo
go get github.com/wailsapp/wails/v3@b1daf95f3
```

You should now have `main.go`, a Go module, a `frontend` directory and a native
`wails.yaml`. The built-in template uses Wake directly; you do not need to
migrate this new project. A community template can have custom legacy tasks
and may instead require the migration review described below.

The completion screen shows `wails3 build setup` and `wails3 dev`. The profile
name matters: `wails3 build setup` uses the selections you just saved. Plain
`wails3 build` does not automatically select the `setup` profile. The profile
name does not mean you need to run `wails3 setup`.

## Build and run your first application

Run the command shown by the wizard:

```sh
wails3 build setup
```

The first build installs frontend dependencies, generates bindings and compiles
the application. Later builds can reuse unchanged work. Wait for **build
succeeded** and read the output paths printed by the CLI. The default output
directory is `bin`; a macOS GUI is assembled into an `.app` bundle.

Start development mode on your own computer:

```sh
wails3 dev
```

You should see the starter application window. Enter your name and try its
greeting. Change the name input's placeholder in `frontend/index.html` and save
the file; the frontend should refresh. Edit the greeting in `greetservice.go`
to exercise a Go rebuild. Press **Ctrl+C** in the terminal when finished.

Development is a local host workflow. A release profile can contain several
cross-build targets without changing which machine can run the development
window. To build and launch a production application locally, use
`wails3 run`; see the [run reference](/reference/wails-yaml-run/).

You now have a source project, a reviewed build profile and a working local
development loop. Commit the source and `wails.yaml` to version control.
Generated `.wails/` working files are not your configuration source.

## What you can do next

| Goal | Wake support | What to check first |
| --- | --- | --- |
| Build for another desktop | Windows, macOS and Linux targets, including both macOS architectures and universal output. | Pure Go and CGo have different compiler requirements. Test on the destination OS and architecture. |
| Make installers or packages | MSIX, NSIS, macOS app/DMG, DEB, RPM, Arch, Void and AppImage. | NSIS needs `makensis`. AppImage needs a target runtime; GUI cross-packaging also needs a supplied Linux runtime tree. |
| Sign a release | Portable exported-key signing for Windows, macOS and DEB/RPM/Arch; native identity stores remain available. | Each identity, format and host has limits. AppImage and Void have no built-in signing flow. |
| Notarise a macOS app | Team API-key submission and ticket stapling, or a macOS keychain profile. | A suitable Developer ID signature and timestamp are still required. Go-generated test signatures passing native checks is not proof of Apple accepting your release. |
| Build several programs | A GUI, CLI or server primary program plus named companion programs. | Declare companion packages and kinds in YAML. The wizard does not generate their source. |
| Package files built elsewhere | `wails3 build setup --skip-build`. | Supply matching executable names and architectures in `build.output`. This flag is a one-time command choice, not a permanent profile setting. |
| Add release automation | Named profiles, bounded lifecycle hooks, caching and parallel package work. | Use literal YAML and declared project-owned scripts. Arbitrary Taskfile behaviour is not a Wake extension mechanism. |

Edit `profiles.setup`, other named profiles and advanced settings in
`wails.yaml`. Use `wails3 setup` to change shared preferences and credentials.
Use the
[YAML build guide](/experimental/yaml-builds/) for complete examples and package
requirements. GUI applications still need the destination's system libraries;
Pure Go does not turn them into completely static applications.

iOS and Android remain experimental, need their platform SDKs and CGo, and are
outside this first desktop lesson. The wizard does not make mobile builds
portable across arbitrary hosts.

## Migrate an existing project

Work on a separate branch or copy of a working **Wails v3** project. This is a
build-configuration migration, not an upgrade from Wails v2. Keep your existing
build result and tests so you can compare behaviour after activation.

@note{type="note" title="Start with a reviewed migration"}
Fresh projects generated from the current built-in legacy templates can be
migrated before their frontend has been built. Wake recognises their stock task
selectors and development lifecycle without executing Taskfiles or evaluating
custom expressions.

Custom tasks, variable overrides and development policies still require review.
If analysis cannot represent your project, leave it on the legacy build and
report the diagnostics. Do not remove behaviour you need just to make activation
pass. The new-project tutorial above does not need migration.
@end

### 1. Analyse without changing the active build

Enable Wake and use the preview CLI from the start of this page. From the
existing project root, run:

```sh
wails3 migrate --dry-run
```

Read the reported supported settings and blockers. The analyser reads legacy
declarations; it does not execute their commands. Dry-run does not create an
active manifest. For machine-readable diagnostics, use `wails3 migrate --json`;
use `--dry-run --json` if you also want to avoid writing a draft.

### 2. Create and review the inactive proposal

```sh
wails3 migrate
```

This creates **`wails.migrated.yaml`**, not `wails.yaml`. Inspect it alongside
your existing Taskfiles and `build/config.yml`. Check project metadata,
frontend commands, target architectures, package contents, assets and signing
settings. Keep secrets in environment variables or your credential store.

Known built-in tasks can be mapped to native settings. Conservative,
argument-free lifecycle tasks that call a project-owned script can become
hooks. Other reachable custom task behaviour is a blocker, not something Wake
silently drops or runs through Task as a fallback.

If a draft already exists, migration does not overwrite your reviewed edits.
Keep a copy before deliberately replacing or regenerating it. Do not rename a
blocked draft to `wails.yaml` to bypass the activation checks. Reduce or replace
unsupported legacy behaviour, and rerun analysis before activation. If the
project cannot be represented yet, keep its legacy build and report the gap.

### 3. Activate explicitly

After reviewing the draft and resolving blockers, run:

```sh
wails3 migrate --activate --dry-run
wails3 migrate --activate
```

The first command checks activation without publishing it. The second rechecks
the legacy declarations and validates the draft before atomically creating
`wails.yaml`; successful activation removes the inactive `wails.migrated.yaml`.
It leaves Taskfiles, `build/config.yml`, scripts and assets untouched. From this
point, the active manifest controls Wake builds. The old files are not a second
configuration source.

Some older projects already have a metadata-only `wails.yaml` with a top-level
`info:` section. Although it is not a native Wake manifest, it still occupies
the activation destination. Preserve and deliberately move that legacy file
before activation; the migrator will not overwrite or rename it for you.

### 4. Verify the migrated application

```sh
wails3 config check
wails3 build
wails3 dev
```

Compare the application and package contents with the previous build, including
your custom release behaviour. Run your project tests. Check each migrated
profile explicitly, for example `wails3 config check release` and
`wails3 build release` when the draft contains a `release` profile.

Edit the migrated profiles directly in `wails.yaml`. `wails3 init -ui` creates
new projects; `wails3 setup` manages global preferences. Neither is an editor
for an existing project's release profile.

### If you need to stop the experiment

Unsetting `WAILS_EXP_USE_WAKE` disables the opt-in, but an active `wails.yaml`
still blocks automatic Taskfile routing. Invalid YAML also stops with an error.
There is no reverse migration command.

To return to the old workflow, restore the complete pre-migration project
revision in version control, including the absence of the active manifest, and
use its compatible CLI and dependencies. Preserve any new source changes you
need before doing that. Keeping a separate trial branch makes this deliberate
recovery much easier.

## Diagnose a problem and give feedback

| Symptom | First check |
| --- | --- |
| Project creation skips Application, Targets, Packages and Tools | Set `WAILS_EXP_USE_WAKE` in the launching terminal and check which `wails3` executable is on `PATH`. Global `wails3 setup` always keeps its preferences flow. |
| A community template needs migration | Review its legacy configuration with the migration commands before activation. |
| A saved build reports missing tools | Creating a project records intent. Install the reported host requirements before building; creation does not install every selected tool. |
| A cross-built executable will not start here | Run it on the selected operating system and architecture; use the host target for local development. |
| Prebuilt mode cannot find a binary | Check the [prebuilt filenames and paths](/experimental/yaml-builds/#use-prebuilt-desktop-executables), including multi-target subdirectories. |
| A manifest change fails | Run `wails3 config check` for the relevant profile and follow the field diagnostic. Wake does not fall back to Taskfiles. |

Inspect a build before running it with `wails3 build setup --plan --json`.
`wails3 config check` validates the configuration, and `wails3 eject` writes an
inactive `wails.ejected.yaml` snapshot without changing the active manifest.
Review command output before sharing it: project paths and credential
locations can identify your local environment.

Share feedback in [Experiments Discussions](https://github.com/wailsapp/wails/discussions/categories/experiments)
or [file a reproducible bug](https://github.com/wailsapp/wails/issues).
Include your CLI and framework revisions, host OS/architecture, selected target,
the exact command, whether you selected Pure Go, and the relevant diagnostic.
Say whether the failure was in compilation, packaging, signing or running the
application. Include a minimal manifest with secrets removed and `wails3 doctor`
output where useful.
