<p align="center" style="text-align: center">
  <img src="./assets/images/logo-universal.png" width="55%"><br/>
</p>

<p align="center">
  Build desktop applications using Go & Web Technologies.<br/><br/>
  <a href="https://github.com/wailsapp/wails/blob/master/LICENSE">
    <img src="https://img.shields.io/badge/License-MIT-blue.svg" />
  </a>
  <a href="https://goreportcard.com/report/github.com/wailsapp/wails">
    <img src="https://goreportcard.com/badge/github.com/wailsapp/wails" />
  </a>
  <a href="http://godoc.org/github.com/wailsapp/wails">
    <img src="https://img.shields.io/badge/godoc-reference-blue.svg" />
  </a>
  <a href="https://www.codefactor.io/repository/github/wailsapp/wails">
    <img src="https://www.codefactor.io/repository/github/wailsapp/wails/badge" alt="CodeFactor" />
  </a>
  <a href="https://github.com/wailsapp/wails/issues">
    <img src="https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat" alt="CodeFactor" />
  </a>
  <a href="https://app.fossa.com/projects/git%2Bgithub.com%2Fwailsapp%2Fwails?ref=badge_shield" alt="FOSSA Status">
    <img src="https://app.fossa.com/api/projects/git%2Bgithub.com%2Fwailsapp%2Fwails.svg?type=shield" />
  </a>
  <a href="https://github.com/avelino/awesome-go" rel="nofollow">
    <img src="https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg" alt="Awesome" />
  </a>
  <a href="https://github.com/wailsapp/wails/workflows/release/badge.svg?branch=master" rel="nofollow">
    <img src="https://github.com/wailsapp/wails/workflows/release/badge.svg?branch=master" alt="Release Pipelines" />
  </a>
</p>

<hr/>
<h3 align="center">
<strong>
PLEASE NOTE: As we are approaching the v2 release, we are not accepting any new feature requests or bug reports for v1. If you have a critical issue, please open a ticket and state why it is critical.
</strong>
</h3>
<hr/>

## Internationalization

[English](README.md) | [简体中文](README.zh-Hans.md) | [日本語](README.ja.md)

## Table of Contents

<details>
  <summary>Click me to Open/Close the directory listing</summary>

- [1. Internationalization](#internationalization)
- [2. Table of Contents](#table-of-contents)
- [3. Introduction](#introduction)
  - [3.1 Official Website](#official-website)
  - [3.2 Roadmap](#roadmap)
- [4. Features](#features)
- [5. Sponsors](#sponsors)
- [6. Getting Started](#getting-started)
- [7. FAQ](#faq)
- [8. Contributors](#contributors)
- [9. Special Mentions](#special-mentions)
- [10. Special Thanks](#special-thanks)
- [11. License](#license)

</details>

## Introduction

The traditional method of providing web interfaces to Go programs is via a built-in web server. Wails offers a different
approach: it provides the ability to wrap both Go code and a web frontend into a single binary. Tools are provided to
make this easy for you by handling project creation, compilation and bundling. All you have to do is get creative!

### Official Website

Version 2:

Wails v2 has been released in Beta for all 3 platforms. Check out the [new website](https://wails.io) if you are
interested in trying it out.

Legacy v1 Website:

The legacy v1 docs can be found at [https://wails.app](https://wails.app).

### Roadmap

The project roadmap may be found [here](https://github.com/wailsapp/wails/discussions/1484). Please consult
this before open up an enhancement request.

## Features

- Use standard Go for the backend
- Use any frontend technology you are already familiar with to build your UI
- Quickly create rich frontends for your Go programs using pre-built templates
- Easily call Go methods from Javascript
- Auto-generated Typescript definitions for your Go structs and methods
- Native Dialogs & Menus
- Supports modern translucency and "frosted window" effects
- Unified eventing system between Go and Javascript
- Powerful cli tool to quickly generate and build your projects
- Multiplatform
- Uses native rendering engines - _no embedded browser_!

## Sponsors

This project is supported by these kind people / companies:

<a href="https://github.com/sponsors/leaanthony" style="width:100px;">
  <img src="/assets/images/sponsors/silver-sponsor.png" width="100"/>
</a>
<a href="https://github.com/selvindev" style="width:100px;">
  <img src="https://github.com/selvindev.png?size=100" width="100"/>
</a>
<br/>
<br/>
<a href="https://github.com/sponsors/leaanthony" style="width:100px;">
  <img src="/assets/images/sponsors/bronze-sponsor.png" width="100"/>
</a>

<a href="https://github.com/codydbentley" style="width:100px">
  <img src="https://github.com/codydbentley.png?size=100" width="100"/>
</a>
<a href="https://www.easywebadv.it/" style="width:100px">
  <img src="website/static/img/easyweb.png" width="100"/>
</a>
<br/>
<br/>
<a href="https://github.com/matryer" style="width:100px">
  <img src="https://github.com/matryer.png" width="100"/>
</a>
<a href="https://github.com/tc-hib" style="width:55px">
 <img src="https://github.com/tc-hib.png?size=55" width="55"/>
</a>
<a href="https://github.com/picatz" style="width:50px">
  <img src="https://github.com/picatz.png?size=50" width="50"/>
</a>
<a href="https://github.com/tylertravisty" style="width:50px">
  <img src="https://github.com/tylertravisty.png?size=50" width="50"/>
</a>
<a href="https://github.com/akhudek" style="width:50px">
  <img src="https://github.com/akhudek.png?size=50" width="50"/>
</a>
<a href="https://github.com/trea" style="width:50px">
  <img src="https://github.com/trea.png?size=50" width="50"/>
</a>
<a href="https://github.com/LanguageAgnostic" style="width:55px">
  <img src="https://github.com/LanguageAgnostic.png?size=55" width="55"/>
</a>
<a href="https://github.com/fcjr" style="width:55px">
  <img src="https://github.com/fcjr.png?size=55" width="55"/>
</a>
<a href="https://github.com/nickarellano" style="width:60px">
  <img src="https://github.com/nickarellano.png?size=60" width="60"/>
</a>
<a href="https://github.com/bglw" style="width:65px">
  <img src="https://github.com/bglw.png?size=65" width="65"/>
</a>
<a href="https://github.com/marcus-crane" style="width:65px">
  <img src="https://github.com/marcus-crane.png?size=65" width="65"/>
</a>
<a href="https://github.com/bbergshaven" style="width:45px">
  <img src="https://github.com/bbergshaven.png?size=45" width="45"/>
</a>
<a href="https://github.com/Gilgames000" style="width:45px">
  <img src="https://github.com/Gilgames000.png?size=45" width="45"/>
</a>
<a href="https://github.com/ilgityildirim" style="width:50px">
  <img src="https://github.com/ilgityildirim.png?size=50" width="50"/>
</a>
<a href="https://github.com/questrail" style="width:50px">
  <img src="https://github.com/questrail.png?size=50" width="50"/>
</a>
<a href="https://github.com/DonTomato" style="width:45px">
  <img src="https://github.com/DonTomato.png?size=45" width="45"/>
</a>
<a href="https://github.com/taigrr" style="width:55px">
  <img src="https://github.com/taigrr.png?size=55" width="55"/>
</a>
<a href="https://github.com/charlie-dee" style="width:55px">
  <img src="https://github.com/charlie-dee.png?size=55" width="55"/>
</a>
<a href="https://github.com/michaelolson1996" style="width:55px">
  <img src="https://github.com/michaelolson1996.png?size=55" width="55"/>
</a>
<a href="https://github.com/GargantuaX" style="width:45px">
  <img src="https://github.com/GargantuaX.png?size=45" width="45"/>
</a>
<a href="https://github.com/CharlieGo88" style="width:55px">
  <img src="https://github.com/CharlieGo88.png?size=55" width="55"/>
</a>
<a href="https://github.com/Bironou" style="width:55px">
  <img src="https://github.com/Bironou.png?size=55" width="55"/>
</a>
<a href="https://github.com/Shackelford-Arden" style="width:55px">
  <img src="https://github.com/Shackelford-Arden.png?size=55" width="55"/>
</a>
<a href="https://github.com/boostchicken" style="width:65px">
  <img src="https://github.com/boostchicken.png?size=65" width="65"/>
</a>
<a href="https://github.com/iansinnott" style="width:55px">
  <img src="https://github.com/iansinnott.png?size=55" width="55"/>
</a>
<a href="https://github.com/Ilshidur" style="width:50px">
  <img src="https://github.com/Ilshidur.png?size=50" width="50"/>
</a>

## Getting Started

The installation instructions are on the [official website](https://wails.io/docs/gettingstarted/installation).

## FAQ

- Is this an alternative to Electron?

  Depends on your requirements. It's designed to make it easy for Go programmers to make lightweight desktop
  applications or add a frontend to their existing applications. Wails v2 does offer native elements such as menus
  and dialogs, so it is becoming a lightweight electron alternative.

- Who is this project aimed at?

  Go programmers who want to bundle an HTML/JS/CSS frontend with their applications, without resorting to creating a
  server and opening a browser to view it.

- What's with the name?

  When I saw WebView, I thought "What I really want is tooling around building a WebView app, a bit like Rails is to
  Ruby". So initially it was a play on words (Webview on Rails). It just so happened to also be a homophone of the
  English name for the [Country](https://en.wikipedia.org/wiki/Wales) I am from. So it stuck.

## Stargazers over time

[![Stargazers over time](https://starchart.cc/wailsapp/wails.svg)](https://starchart.cc/wailsapp/wails)

## Contributors

The contributors list is getting too big for the readme! All the amazing people who have contributed to this
project have their own page [here](https://wails.io/credits#contributors).

## Special Mentions

Without the following people, this project would never have existed:

- [Dustin Krysak](https://wiki.ubuntu.com/bashfulrobot) - His support and feedback has been immense. More patience than
  you can throw a stick at (Not long now Dustin!).
- [Serge Zaitsev](https://github.com/zserge) - Creator of [Webview](https://github.com/zserge/webview) which Wails uses
  for the windowing.
- [Byron](https://github.com/bh90210) - At times, Byron has single handedly kept this project alive. Without his
  incredible input, we never would have got to v1.

This project was mainly coded to the following albums:

- [Manic Street Preachers - Resistance Is Futile](https://open.spotify.com/album/1R2rsEUqXjIvAbzM0yHrxA)
- [Manic Street Preachers - This Is My Truth, Tell Me Yours](https://open.spotify.com/album/4VzCL9kjhgGQeKCiojK1YN)
- [The Midnight - Endless Summer](https://open.spotify.com/album/4Krg8zvprquh7TVn9OxZn8)
- [Gary Newman - Savage (Songs from a Broken World)](https://open.spotify.com/album/3kMfsD07Q32HRWKRrpcexr)
- [Steve Vai - Passion & Warfare](https://open.spotify.com/album/0oL0OhrE2rYVns4IGj8h2m)
- [Ben Howard - Every Kingdom](https://open.spotify.com/album/1nJsbWm3Yy2DW1KIc1OKle)
- [Ben Howard - Noonday Dream](https://open.spotify.com/album/6astw05cTiXEc2OvyByaPs)
- [Adwaith - Melyn](https://open.spotify.com/album/2vBE40Rp60tl7rNqIZjaXM)
- [Gwidaith Hen Fran - Cedors Hen Wrach](https://open.spotify.com/album/3v2hrfNGINPLuDP0YDTOjm)
- [Metallica - Metallica](https://open.spotify.com/album/2Kh43m04B1UkVcpcRa1Zug)
- [Bloc Party - Silent Alarm](https://open.spotify.com/album/6SsIdN05HQg2GwYLfXuzLB)
- [Maxthor - Another World](https://open.spotify.com/album/3tklE2Fgw1hCIUstIwPBJF)
- [Alun Tan Lan - Y Distawrwydd](https://open.spotify.com/album/0c32OywcLpdJCWWMC6vB8v)

## Special Thanks

<p align="center" style="text-align: center">
   <a href="https://pace.dev"><img src="/assets/images/pace.jpeg"/></a><br/>
   A <i>huge</i> thanks to <a href="https://pace.dev">Pace</a> for sponsoring the project and helping the efforts to get Wails ported to Apple Silicon!<br/><br/>
   If you are looking for a Project Management tool that's powerful but quick and easy to use, check them out!<br/><br/>
</p>

<p align="center" style="text-align: center">
   A special thank you to JetBrains for donating licenses to us!<br/><br/>
   Please click the logo to let them know your appreciation!<br/><br/>
   <a href="https://www.jetbrains.com?from=Wails"><img src="/assets/images/jetbrains-grayscale.png" width="30%"></a>
</p>

## License

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fwailsapp%2Fwails.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fwailsapp%2Fwails?ref=badge_large)
