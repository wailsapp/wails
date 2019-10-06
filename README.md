<p align="center" style="text-align: center">
   <img src="https://github.com/wailsapp/docs/raw/master/.vuepress/public/media/logo_cropped.png" width="40%"><br/>
</p>
<p align="center">
   A framework for building desktop applications using Go & Web Technologies.<br/><br/>
   <a href="https://github.com/wailsapp/wails/blob/master/LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg"></a>
   <a href="https://goreportcard.com/report/github.com/wailsapp/wails"><img src="https://goreportcard.com/badge/github.com/wailsapp/wails"/></a>
   <a href="http://godoc.org/github.com/wailsapp/wails"><img src="https://img.shields.io/badge/godoc-reference-blue.svg"/></a>
   <a href="https://www.codefactor.io/repository/github/wailsapp/wails"><img src="https://www.codefactor.io/repository/github/wailsapp/wails/badge" alt="CodeFactor" /></a>
   <a href="https://github.com/wailsapp/wails/issues"><img src="https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat" alt="CodeFactor" /></a>
   <a href="https://houndci.com"><img src="https://img.shields.io/badge/Reviewed_by-Hound-8E64B0.svg"/></a>
   <a href="https://github.com/avelino/awesome-go" rel="nofollow"><img src="https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg" alt="Awesome"></a>
   <a href="https://dashboard.guardrails.io/default/gh/wailsapp/wails"><img src="https://badges.guardrails.io/wailsapp/wails.svg?token=53657bc22ec360d7673c894fdd70568e918ec581d10d84427ed4de5fe1eeff1a"></a>
</p>

The traditional method of providing web interfaces to Go programs is via a built-in web server. Wails offers a different approach: it provides the ability to wrap both Go code and a web frontend into a single binary. Tools are provided to make this easy for you by handling project creation, compilation and bundling. All you have to do is get creative!

## Features

- Use standard Go libraries/frameworks for the backend
- Use any frontend technology to build your UI
- Quickly create Vue, Vuetify or React frontends for your Go programs
- Expose Go methods/functions to the frontend via a single bind command
- Uses native rendering engines - no embedded browser
- Shared events system
- Native file dialogs
- Powerful cli tool
- Multiplatform

## Project Status

Wails is currently in Beta. Please make sure you read the [Project Status](https://wails.app/project_status.html) if you are interested in using this project.

## Installation

Wails uses cgo to bind to the native rendering engines so a number of platform dependent libraries are needed as well as an installation of Go. The basic requirements are:

- Go 1.12
- npm

### MacOS

Make sure you have the xcode command line tools installed. This can be done by running:

`xcode-select --install`

### Linux

#### Debian/Ubuntu

`sudo apt install libgtk-3-dev libwebkit2gtk-4.0-dev`

_Debian: 8, 9, 10_

_Ubuntu: 16.04, 18.04, 19.04_

_Also succesfully tested on: Zorin 15, Parrot 4.7, Linuxmint 19, Elementary 5, Kali, Neon_

#### Arch Linux

`sudo pacman -S webkit2gtk gtk3`

_Also succesfully test on: Manjaro & ArcoLinux_

#### Centos

`sudo yum install webkitgtk3-devel gtk3-devel`

_CentOS 6, 7_

#### Fedora

`sudo yum install webkit2gtk3-devel gtk3-devel`

_Fedora 29, 30_

#### VoidLinux & VoidLinux-musl

`xbps-install gtk+3-devel webkit2gtk-devel`

#### Gentoo

`sudo emerge gtk+:3 webkit-gtk`

### Windows

Windows requires gcc and related tooling. The recommended download is from [http://tdm-gcc.tdragon.net/download](http://tdm-gcc.tdragon.net/download). Once this is installed, you are good to go.

## Installation

**Ensure Go modules are enabled: GO111MODULE=on and go/bin is in your PATH variable.**

Installation is as simple as running the following command:

<pre style='color:white'>
go get github.com/wailsapp/wails/cmd/wails
</pre>

## Next Steps

It is recommended at this stage to read the comprehensive documentation at [https://wails.app](https://wails.app).

## FAQ

 * Is this an alternative to Electron?

   Depends on your requirements. It's designed to make it easy for Go programmers to make lightweight desktop applications or add a frontend to their existing applications. Whilst Wails does not currently offer hooks into native elements such as menus, this may change in the future.

 * Who is this project aimed at?

   Go programmers who want to bundle an HTML/JS/CSS frontend with their applications, without resorting to creating a server and opening a browser to view it.

 * What's with the name?

   When I saw WebView, I thought "What I really want is tooling around building a WebView app, a bit like Rails is to Ruby". So initially it was a play on words (Webview on Rails). It just so happened to also be a homophone of the English name for the [Country](https://en.wikipedia.org/wiki/Wales) I am from. So it stuck.

## Shoulders of Giants

Without the following people, this project would never have existed:

  * [Dustin Krysak](https://wiki.ubuntu.com/bashfulrobot) - His support and feedback has been immense. More patience than you can throw a stick at (Not long now Dustin!).
  * [Serge Zaitsev](https://github.com/zserge) - Creator of [Webview](https://github.com/zserge/webview) which Wails uses for the windowing.

And without [these people](CONTRIBUTORS.md), it wouldn't be what it is today. A huge thank you to each and every one of you!

Special Mentions:

  * [Bill Kennedy](https://twitter.com/goinggodotnet) - Go guru, encourager and all-round nice guy, whose infectious energy and inspiration powered me on when I had none left.
  * [Mark Bates](https://github.com/markbates) - Creator of [Packr](https://github.com/gobuffalo/packr), inspiration for packing strategies which fed into some of the tooling.

This project was mainly coded to the following albums:

  * [Manic Street Preachers - Resistance Is Futile](https://open.spotify.com/album/1R2rsEUqXjIvAbzM0yHrxA)
  * [Manic Street Preachers - This Is My Truth, Tell Me Yours](https://open.spotify.com/album/4VzCL9kjhgGQeKCiojK1YN)
  * [The Midnight - Endless Summer](https://open.spotify.com/album/4Krg8zvprquh7TVn9OxZn8)
  * [Gary Newman - Savage (Songs from a Broken World)](https://open.spotify.com/album/3kMfsD07Q32HRWKRrpcexr)
  * [Steve Vai - Passion & Warfare](https://open.spotify.com/album/0oL0OhrE2rYVns4IGj8h2m)
  * [Ben Howard - Every Kingdom](https://open.spotify.com/album/1nJsbWm3Yy2DW1KIc1OKle)
  * [Ben Howard - Noonday Dream](https://open.spotify.com/album/6astw05cTiXEc2OvyByaPs)
  * [Adwaith - Melyn](https://open.spotify.com/album/2vBE40Rp60tl7rNqIZjaXM)
  * [Gwidaith Hen Fran - Cedors Hen Wrach](https://open.spotify.com/album/3v2hrfNGINPLuDP0YDTOjm)
  * [Metallica - Metallica](https://open.spotify.com/album/2Kh43m04B1UkVcpcRa1Zug)
  * [Bloc Party - Silent Alarm](https://open.spotify.com/album/6SsIdN05HQg2GwYLfXuzLB)
  * [Maxthor - Another World](https://open.spotify.com/album/3tklE2Fgw1hCIUstIwPBJF)
  * [Alun Tan Lan - Y Distawrwydd](https://open.spotify.com/album/0c32OywcLpdJCWWMC6vB8v)
