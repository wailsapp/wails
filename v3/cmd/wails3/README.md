# The Wails CLI

The Wails CLI is a command line tool that allows you to create, build and run Wails applications.
There are a number of commands related to tooling, such as icon generation and asset bundling.

## Commands

### task

The `task` command is for running tasks defined in `Taskfile.yml`. It is a wrapper around [Task](https://taskfile.dev).

### generate

The `generate` command is used to generate resources and assets for your Wails project.
It can be used to generate many things including: 
  - application icons, 
  - resource files for Windows applications
  - Info.plist files for macOS deployments

#### icons

The `icons` command generates icons for your project.

| Flag                 | Type   | Description                                              | Default                       |
|----------------------|--------|----------------------------------------------------------|-------------------------------|
| `-example`           | bool   | Generate example icon file (`appicon.png`)                |                               |
| `-input`             | string | Input PNG file                                           | `build/appicon.png`           |
| `-sizes`             | string | Windows ICO sizes (comma-separated)                       | `256,128,64,48,32,16`         |
| `-windowsfilename`   | string | Windows icon output                                      | `build/windows/icon.ico`      |
| `-macfilename`       | string | macOS icon bundle output                                 | `build/darwin/icon.icns`      |
| `-iconcomposerinput` | string | Input Icon Composer file (`.icon`)                        |                               |
| `-macassetdir`       | string | Output directory for macOS `Assets.car` and `icons.icns` |                               |
| `-linuxoutputdir`    | string | Output directory for Linux hicolor PNGs                   |                               |
| `-linuxsizes`        | string | Linux PNG sizes (comma-separated)                         | `16,32,48,64,128,256,512`     |

```bash
wails3 generate icons \
  -input build/appicon.png \
  -windowsfilename build/windows/icon.ico \
  -macfilename build/darwin/icons.icns \
  -linuxoutputdir build/linux/icons
```

Linux images are resized without stretching: non-square input is centered on a transparent square canvas. The PNG output is deterministic and can be regenerated from `build/appicon.png` on any supported build host.

#### syso

The `syso` command generates a Windows resource file (aka `.syso`).

```bash
wails3 generate syso <options>
```

| Flag        | Type   | Description                                | Default          |
|-------------|--------|--------------------------------------------|------------------|
| `-example`  | bool   | Generates example manifest & info files    |                  |
| `-manifest` | string | The manifest file                          |                  |
| `-info`     | string | The info.json file                         |                  |
| `-icon`     | string | The icon file                              |                  |
| `-out`      | string | The output filename for the syso file      | `wails.exe.syso` |
| `-arch`     | string | The target architecture  (amd64,arm64,386) | `runtime.GOOS`   |

If `-example` is provided, the command will generate example manifest and info files
in the current directory and exit.

If `-manifest` is provided, the command will use the provided manifest file to generate
the syso file.

If `-info` is provided, the command will use the provided info.json file to set the version
information in the syso file.

NOTE: We use [winres](https://github.com/tc-hib/winres) to generate the syso file. Please
refer to the winres documentation for more information.

NOTE: Whilst the tool will work for 32-bit Windows, it is not supported. Please use 64-bit.

#### defaults

```bash
wails3 generate defaults      
```
This will generate all the default assets and resources in the current directory.

#### bindings

```bash
wails3 generate bindings
```

Generates bindings and models for your bound Go methods and structs.
