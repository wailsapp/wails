# The Wails CLI

The Wails CLI is a command line tool that allows you to create, build and run
Wails applications. There are a number of commands related to tooling, such as
icon generation and asset bundling.

## Commands

### run

The run command is for running tasks defined in `Taskfile.yml`.

| Flag | Type   | Description                 | Default |
| ---- | ------ | --------------------------- | ------- |
| `-t` | string | The name of the task to run |         |

### generate

The `generate` command is used to generate resources and assets for your Wails
project. It can be used to generate many things including:

- application icons,
- resource files for Windows applications
- Info.plist files for macOS deployments

#### icon

The `icon` command generates icons for your project.

| Flag               | Type   | Description                                          | Default               |
| ------------------ | ------ | ---------------------------------------------------- | --------------------- |
| `-example`         | bool   | Generates example icon file (appicon.png)            |                       |
| `-input`           | string | The input image file                                 |                       |
| `-sizes`           | string | The sizes to generate in .ico file (comma separated) | "256,128,64,48,32,16" |
| `-windowsFilename` | string | The output filename for the Windows icon             | icons.ico             |
| `-macFilename`     | string | The output filename for the Mac icon bundle          | icons.icns            |

```bash
wails generate icon -input myicon.png -sizes "32,64,128" -windowsFilename myicon.ico -macFilename myicon.icns
```

This will generate icons for mac and windows and save them in the current
directory as `myicon.ico` and `myicons.icns`.

#### syso

The `syso` command generates a Windows resource file (aka `.syso`).

```bash
wails generate syso <options>
```

| Flag        | Type   | Description                               | Default          |
| ----------- | ------ | ----------------------------------------- | ---------------- |
| `-example`  | bool   | Generates example manifest & info files   |                  |
| `-manifest` | string | The manifest file                         |                  |
| `-info`     | string | The info.json file                        |                  |
| `-icon`     | string | The icon file                             |                  |
| `-out`      | string | The output filename for the syso file     | `wails.exe.syso` |
| `-arch`     | string | The target architecture (amd64,arm64,386) | `runtime.GOOS`   |

If `-example` is provided, the command will generate example manifest and info
files in the current directory and exit.

If `-manifest` is provided, the command will use the provided manifest file to
generate the syso file.

If `-info` is provided, the command will use the provided info.json file to set
the version information in the syso file.

NOTE: We use [winres](https://github.com/tc-hib/winres) to generate the syso
file. Please refer to the winres documentation for more information.

NOTE: Whilst the tool will work for 32-bit Windows, it is not supported. Please
use 64-bit.

#### defaults

```bash
wails generate defaults
```

This will generate all the default assets and resources in the current
directory.

#### bindings

```bash
wails generate bindings
```

Generates bindings and models for your bound Go methods and structs.
