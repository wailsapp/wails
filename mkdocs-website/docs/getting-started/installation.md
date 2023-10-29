# Installation

To install the Wails CLI, ensure you have [Go 1.21+](https://go.dev/dl/)
installed and run:

```shell
git clone https://github.com/wailsapp/wails.git
git checkout v3-alpha
cd wails/cmd/wails3
go install
```

## Supported Platforms

- Windows 10/11 AMD64/ARM64
- MacOS 10.13+ AMD64
- MacOS 11.0+ ARM64
- Linux AMD64/ARM64

## Dependencies

Wails has a number of common dependencies that are required before installation:

=== "Go 1.21+"

    Download Go from the [Go Downloads Page](https://go.dev/dl/).

    Ensure that you follow the official [Go installation instructions](https://go.dev/doc/install). You will also need to ensure that your `PATH` environment variable also includes the path to your `~/go/bin` directory. Restart your terminal and do the following checks:

    - Check Go is installed correctly: `go version`
    - Check `~/go/bin` is in your PATH variable: `echo $PATH | grep go/bin`

=== "npm (Optional)"

    Although Wails doesn't require npm to be installed, it is needed if you want to use the bundled templates.

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

    Linux requires the standard `gcc` build tools plus `libgtk3` and `libwebkit`. Rather than list a ton of commands for different distros, Wails can try to determine what the installation commands are for your specific distribution. Run <code>wails doctor</code> after installation to be shown how to install the dependencies. If your distro/package manager is not supported, please let us know on discord.

## System Check

Running `wails3 doctor` will check if you have the correct dependencies
installed. If not, it will advise on what is missing and help on how to rectify
any problems.

## The `wails3` command appears to be missing?

If your system is reporting that the `wails3` command is missing, check the
following:

- Make sure you have followed the Go installation guide correctly.
- Check that the `go/bin` directory is in the `PATH` environment variable.
- Close/Reopen current terminals to pick up the new `PATH` variable.
