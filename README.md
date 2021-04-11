<p align="center" style="text-align: center">
   <img src="logo_cropped.png" width="40%"><br/>
</p>
<p align="center">
   Build desktop applications using Go & Web Technologies.<br/><br/>
   <a href="https://github.com/wailsapp/wails/blob/master/LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg"></a>
   <a href="https://goreportcard.com/report/github.com/wailsapp/wails"><img src="https://goreportcard.com/badge/github.com/wailsapp/wails"/></a>
   <a href="http://godoc.org/github.com/wailsapp/wails"><img src="https://img.shields.io/badge/godoc-reference-blue.svg"/></a>
   <a href="https://www.codefactor.io/repository/github/wailsapp/wails"><img src="https://www.codefactor.io/repository/github/wailsapp/wails/badge" alt="CodeFactor" /></a>
   <a href="https://github.com/wailsapp/wails/issues"><img src="https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat" alt="CodeFactor" /></a>
   <a href="https://app.fossa.io/projects/git%2Bgithub.com%2Fwailsapp%2Fwails?ref=badge_shield" alt="FOSSA Status"><img src="https://app.fossa.io/api/projects/git%2Bgithub.com%2Fwailsapp%2Fwails.svg?type=shield"/></a>
   <a href="https://houndci.com"><img src="https://img.shields.io/badge/Reviewed_by-Hound-8E64B0.svg"/></a>
   <a href="https://github.com/avelino/awesome-go" rel="nofollow"><img src="https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg" alt="Awesome"></a>
   <a href="https://github.com/wailsapp/wails/workflows/release/badge.svg?branch=master" rel="nofollow"><img src="https://github.com/wailsapp/wails/workflows/release/badge.svg?branch=master" alt="Release Pipelines"></a>
</p>

The traditional method of providing web interfaces to Go programs is via a built-in web server. Wails offers a different approach: it provides the ability to wrap both Go code and a web frontend into a single binary. Tools are provided to make this easy for you by handling project creation, compilation and bundling. All you have to do is get creative!

The official docs can be found at [https://wails.app](https://wails.app).

## Features

- Use standard Go for the backend
- Use any frontend technology to build your UI
- Quickly create Vue, Vuetify or React frontends for your Go programs
- Expose Go methods/functions to the frontend via a single bind command
- Uses native rendering engines - no embedded browser
- Shared events system
- Native file dialogs
- Powerful cli tool
- Multiplatform

## Sponsors

This project is supported by these kind people / companies:

<a href="https://www.jetbrains.com?from=Wails" style="width:100px"><img src="jetbrains-grayscale.png" width="100"/></a>
<a href="https://pace.dev" style="width:100px"><img src="pace.jpeg" width="100"/></a>
<a href="https://github.com/tc-hib" style="width:50px;border-radius: 50%">
  <img src="https://github.com/tc-hib.png?size=50" width="50" style="border-radius: 50%"/>
</a>
<a href="https://github.com/picatz" style="width:50px;border-radius: 50%">
  <img src="https://github.com/picatz.png?size=50" width="50" style="border-radius: 50%"/>
</a>
<a href="https://github.com/tylertravisty" style="width:50px;border-radius: 50%">
  <img src="https://github.com/tylertravisty.png?size=50" width="50" style="border-radius: 50%"/>
</a>

## Installation

Wails uses cgo to bind to the native rendering engines so a number of platform dependent libraries are needed as well as an installation of Go. The basic requirements are:

- Go 1.13
- npm

### MacOS

Make sure you have the xcode command line tools installed. This can be done by running:

`xcode-select --install`

### Linux

#### Debian/Ubuntu

`sudo apt install libgtk-3-dev libwebkit2gtk-4.0-dev`

_Debian: 8, 9, 10_

_Ubuntu: 16.04, 18.04, 19.04_

_Also succesfully tested on: Zorin 15, Parrot 4.7, Linuxmint 19, Elementary 5, Kali, Neon_, Pop!_OS

#### Arch Linux / ArchLabs / Ctlos Linux

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
go get -u github.com/wailsapp/wails/cmd/wails
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

## Contributors

<a href="https://github.com/qaisjp"><img src="https://github.com/qaisjp.png?size=40" width="40"/></a>
<a href="https://github.com/alee792"><img src="https://github.com/alee792.png?size=40" width="40"/></a>
<a href="https://github.com/lanzafame"><img src="https://github.com/lanzafame.png?size=40" width="40"/></a>
<a href="https://github.com/mattn"><img src="https://github.com/mattn.png?size=40" width="40"/></a>
<a href="https://github.com/0xflotus"><img src="https://github.com/0xflotus.png?size=40" width="40"/></a>
<a href="https://github.com/mdhender"><img src="https://github.com/mdhender.png?size=40" width="40"/></a>
<a href="https://github.com/fishfishfish2104"><img src="https://github.com/fishfishfish2104.png?size=40" width="40"/></a>
<a href="https://github.com/intelwalk"><img src="https://github.com/intelwalk.png?size=40" width="40"/></a>
<a href="https://github.com/ocelotsloth"><img src="https://github.com/ocelotsloth.png?size=40" width="40"/></a>
<a href="https://github.com/bh90210"><img src="https://github.com/bh90210.png?size=40" width="40"/></a>
<a href="https://github.com/iceleo-com"><img src="https://github.com/iceleo-com.png?size=40" width="40"/></a>
<a href="https://github.com/fallendusk"><img src="https://github.com/fallendusk.png?size=40" width="40"/></a>
<a href="https://github.com/Chronophylos"><img src="https://github.com/Chronophylos.png?size=40" width="40"/></a>
<a href="https://github.com/Vaelatern"><img src="https://github.com/Vaelatern.png?size=40" width="40"/></a>
<a href="https://github.com/mewmew"><img src="https://github.com/mewmew.png?size=40" width="40"/></a>
<a href="https://github.com/kraney"><img src="https://github.com/kraney.png?size=40" width="40"/></a>
<a href="https://github.com/JackMordaunt"><img src="https://github.com/JackMordaunt.png?size=40" width="40"/></a>
<a href="https://github.com/MichaelHipp"><img src="https://github.com/MichaelHipp.png?size=40" width="40"/></a>
<a href="https://github.com/tmclane"><img src="https://github.com/tmclane.png?size=40" width="40"/></a>
<a href="https://github.com/Rested"><img src="https://github.com/Rested.png?size=40" width="40"/></a>
<a href="https://github.com/Jarek-SRT"><img src="https://github.com/Jarek-SRT.png?size=40" width="40"/></a>
<a href="https://github.com/konez2k"><img src="https://github.com/konez2k.png?size=40" width="40"/></a>
<a href="https://github.com/sayuthisobri"><img src="https://github.com/sayuthisobri.png?size=40" width="40"/></a>
<a href="https://github.com/dedo1911"><img src="https://github.com/dedo1911.png?size=40" width="40"/></a>
<a href="https://github.com/fdidron"><img src="https://github.com/fdidron.png?size=40" width="40"/></a>
<a href="https://github.com/Splode"><img src="https://github.com/Splode.png?size=40" width="40"/></a>
<a href="https://github.com/Lyimmi"><img src="https://github.com/Lyimmi.png?size=40" width="40"/></a>
<a href="https://github.com/Unix4ever"><img src="https://github.com/Unix4ever.png?size=40" width="40"/></a>
<a href="https://github.com/timkippdev"><img src="https://github.com/timkippdev.png?size=40" width="40"/></a>
<a href="https://github.com/kyoto44"><img src="https://github.com/kyoto44.png?size=40" width="40"/></a>
<a href="https://github.com/artooro"><img src="https://github.com/artooro.png?size=40" width="40"/></a>
<a href="https://github.com/ilgityildirim"><img src="https://github.com/ilgityildirim.png?size=40" width="40"/></a>
<a href="https://github.com/gelleson"><img src="https://github.com/gelleson.png?size=40" width="40"/></a>
<a href="https://github.com/kmuchmore"><img src="https://github.com/kmuchmore.png?size=40" width="40"/></a>
<a href="https://github.com/aayush420"><img src="https://github.com/aayush420.png?size=40" width="40"/></a>
<a href="https://github.com/Rezrazi"><img src="https://github.com/Rezrazi.png?size=40" width="40"/></a>
<a href="https://github.com/misitebao"><img src="https://github.com/misitebao.png?size=40" width="40"/></a>
<a href="https://github.com/DrunkenPoney"><img src="https://github.com/DrunkenPoney.png?size=40" width="40"/></a>
<a href="https://github.com/SophieAu"><img src="https://github.com/SophieAu.png?size=40" width="40"/></a>
<a href="https://github.com/alexmat"><img src="https://github.com/alexmat.png?size=40" width="40"/></a>
<a href="https://github.com/RH12503"><img src="https://github.com/RH12503.png?size=40" width="40"/></a>
<a href="https://github.com/hi019"><img src="https://github.com/hi019.png?size=40" width="40"/></a></a>
<a href="https://github.com/Igogrek"><img src="https://github.com/Igogrek.png?size=40" width="40"/></a></a>
<a href="https://github.com/aschey"><img src="https://github.com/aschey.png?size=40" width="40"/></a></a>

## Special Mentions

Without the following people, this project would never have existed:

  * [Dustin Krysak](https://wiki.ubuntu.com/bashfulrobot) - His support and feedback has been immense. More patience than you can throw a stick at (Not long now Dustin!).
  * [Serge Zaitsev](https://github.com/zserge) - Creator of [Webview](https://github.com/zserge/webview) which Wails uses for the windowing.
  * [Byron](https://github.com/bh90210) - At times, Byron has single handedly kept this project alive. Without his incredible input, we never would have got to v1.

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

## Licensing

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fwailsapp%2Fwails.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fwailsapp%2Fwails?ref=badge_large)

## Special Thanks

<p align="center" style="text-align: center">
   <a href="https://pace.dev"><img src="pace.jpeg"/></a><br/>
   A *huge* thanks to <a href="https://pace.dev">Pace</a> for sponsoring the project and helping the efforts to get Wails ported to Apple Silicon!<br/><br/>
   If you are looking for a Project Management tool that's powerful but quick and easy to use, check them out!<br/><br/>
</p>

<p align="center" style="text-align: center">
   A special thank you to JetBrains for donating licenses to us!<br/><br/>
   Please click the logo to let them know your appreciation!<br/><br/>
   <a href="https://www.jetbrains.com?from=Wails"><img src="jetbrains-grayscale.png" width="30%"></a>
</p>
