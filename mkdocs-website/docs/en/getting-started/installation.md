# Installation

To install the Wails CLI, first ensure you have the correct dependencies installed:

## Supported Platforms

- Windows 10/11 AMD64/ARM64
- MacOS 10.13+ AMD64
- MacOS 11.0+ ARM64
- Ubuntu 22.04 AMD64/ARM64 (other Linux may work too!)

## Dependencies

Wails has a number of common dependencies that are required before installation:

=== "Go (At least 1.22.4)"

    Download Go from the [Go Downloads Page](https://go.dev/dl/).

    Ensure that you follow the official [Go installation instructions](https://go.dev/doc/install). You will also need to ensure that your `PATH` environment variable also includes the path to your `~/go/bin` directory. Restart your terminal and do the following checks:

    - Check Go is installed correctly: `go version`
    - Check `~/go/bin` is in your PATH variable
        - Mac / Linux: `echo $PATH | grep go/bin`
        - Windows: `$env:PATH -split ';' | Where-Object { $_ -like '*\go\bin' }`

=== "npm (Optional)"

    Although Wails doesn't require npm to be installed, it is needed by most of the bundled templates.

    Download the latest node installer from the [Node Downloads Page](https://nodejs.org/en/download/). It is best to use the latest release as that is what we generally test against.

    Run `npm --version` to verify.

=== "Task (Optional)"

    The Wails CLI embeds a task runner called [Task](https://taskfile.dev/#/installation). It is optional, but recommended. If you do not wish to install Task, you can use the `wails3 task` command instead of `task`.
    Installing Task will give you the greatest flexibility.

## Platform Specific Dependencies

You will also need to install platform specific dependencies:

=== "Mac"

    Wails requires that the xcode command line tools are installed. This can be
    done by running:

    ```
    xcode-select --install
    ```

=== "Windows"

    Wails requires that the [WebView2 Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/) is installed. Some Windows installations will already have this installed. You can check using the `wails doctor` command.

=== "Linux"

    Linux requires the standard `gcc` build tools plus `gtk3` and `webkit2gtk`. Run <code>wails doctor</code> after installation to be shown how to install the dependencies. If your distro/package manager is not supported, please let us know on discord.

## Installation

To install the Wails CLI using Go Modules, run the following commands:

```shell
go install -v github.com/wailsapp/wails/v3/cmd/wails3@latest
```

If you would like to install the latest development version, run the following commands:

```shell
git clone https://github.com/wailsapp/wails.git
cd wails
git checkout v3-alpha
cd v3/cmd/wails3
go install
```

When using the development version, all generated projects will use Go's [replace(https://go.dev/ref/mod#go-mod-file-replace) directive
to ensure projects use the development version of Wails.

## System Check

Running `wails3 doctor` will check if you have the correct dependencies
installed. If not, it will advise on what is missing and help on how to rectify
any problems.

## The `wails3` command appears to be missing?

If your system is reporting that the `wails3` command is missing, check the
following:

- Make sure you have followed the [Go installation guide](#__tabbed_1_1) correctly and that the `go/bin` directory is in the `PATH` environment variable.
- Close/Reopen current terminals to pick up the new `PATH` variable.
