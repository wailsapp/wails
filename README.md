# Wails

<p align="center"><img src="https://github.com/wailsapp/docs/raw/master/.vuepress/public/media/logo.png" width="50%"></p>

A framework for building desktop applications using Go & Web Technologies.

## About

The traditional method of providing web interfaces to Go programs is via a built-in web server. Wails offers a different approach: it provides the ability to wrap both Go code and a web frontend into a single binary. Tools are provided to make this easy for you by handling project creation, compilation and bundling. All you have to do is get creative!

## Features

- Use standard Go libraries/frameworks for the backend
- Use any frontend technology to build your UI
- Expose Go methods/functions to the frontend via a single bind command
- Uses native rendering engines - no embedded browser
- Shared events system
- Native file dialogs
- Powerful cli tool
- Multiplatform

## Installation

Wails uses cgo to bind to the native rendering engines so a number of platform dependent libraries are needed as well as an installation of Go. The basic requirements are:

- Go 1.11 or above
- npm

### MacOS

Make sure you have the xcode command line tools installed. This can be done by running:

`xcode-select --install`

### Linux

#### Ubuntu 18.04

`sudo apt install pkg-config build-essential libgtk-3-dev libwebkit2gtk-4.0-dev`

::: tip
If you have successfully installed these dependencies on a different flavour of Linux, please consider submitting a PR.
:::

### Windows

Windows requires gcc and related tooling. The recommended download is from [http://tdm-gcc.tdragon.net/download](http://tdm-gcc.tdragon.net/download). Once this is installed, you are good to go.

## Installation

Installation is as simple as running the following command:

<pre style='color:white'>
go get -u github.com/wailsapp/wails/cmd/wails
</pre>

## Next Steps

It is recommended at this stage to read the comprehensive documentation at [https://wails.app](https://wails.app).

## FAQ

 * Is this an alternative to Electron?

   Depends on your requirements. It's designed to make it easy for Go programmers to make lightweight desktop applications or add a frontend to their existing applications. Whilst Wails does not currently offer hooks into native elements such as menus, this may change in the future.

 * Who is this project aimed at?

   Go programmers who want to bundle an HMTL frontend with their applications, without resorting to creating a server and opening a browser to view it.

## Shoulders of Giants

Without the following people, this project would never have existed:

  * [Dustin Krysak](https://wiki.ubuntu.com/bashfulrobot) - His support and feedback has been immense. More patience than you can throw a stick at.
  * [Serge Zaitsev](https://github.com/zserge) - Creator of [Webview](https://github.com/zserge/webview) which Wails uses for the windowing.

Special Mentions:

  * [Bill Kennedy](https://twitter.com/goinggodotnet) - Go guru, encourager and all-round nice guy, whose energy and inspiration powered me on when I had none left.
  * [Mark Bates](https://github.com/markbates) - Creator of [Packr](https://github.com/gobuffalo/packr), inspiration for packing strategies.

