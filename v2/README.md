<p align="center" style="text-align: center">
  <img src="../assets/images/logo-universal.png" width="55%"><br/>
</p>

<p align="center">
  Build desktop applications using Go & Web Technologies.
  <br/>
  <br/>
  <a href="https://github.com/wailsapp/wails/blob/master/LICENSE">
    <img alt="GitHub" src="https://img.shields.io/github/license/wailsapp/wails"/>
  </a>
  <a href="https://goreportcard.com/report/github.com/wailsapp/wails">
    <img src="https://goreportcard.com/badge/github.com/wailsapp/wails" />
  </a>
  <a href="https://pkg.go.dev/github.com/wailsapp/wails">
    <img src="https://pkg.go.dev/badge/github.com/wailsapp/wails.svg" alt="Go Reference"/>
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
  <br/>
  <a href="https://github.com/wailsapp/wails/actions/workflows/build.yml" rel="nofollow">
    <img src="https://img.shields.io/github/workflow/status/wailsapp/wails/Build?logo=github" alt="Build" />
  </a>
  <a href="https://github.com/wailsapp/wails/tags" rel="nofollow">
    <img alt="GitHub tag (latest SemVer pre-release)" src="https://img.shields.io/github/v/tag/wailsapp/wails?include_prereleases&label=version"/>
  </a>
</p>

<div align="center">
<strong>
<samp>

[English](README.md) · [简体中文](README.zh-Hans.md) · [日本語](README.ja.md)

</samp>
</strong>
</div>

## Table of Contents

<details>
  <summary>Click me to Open/Close the directory listing</summary>

- [Table of Contents](#table-of-contents)
- [Introduction](#introduction)
  - [Roadmap](#roadmap)
- [Features](#features)
- [Sponsors](#sponsors)
- [Getting Started](#getting-started)
- [FAQ](#faq)
- [Contributors](#contributors)
- [License](#license)

</details>

## Introduction

The traditional method of providing web interfaces to Go programs is via a built-in web server. Wails offers a different
approach: it provides the ability to wrap both Go code and a web frontend into a single binary. Tools are provided to
make this easy for you by handling project creation, compilation and bundling. All you have to do is get creative!

## Features

- Use standard Go for the backend
- Use any frontend technology you are already familiar with to build your UI
- Quickly create rich frontends for your Go programs using pre-built templates
- Easily call Go methods from Javascript
- Auto-generated Typescript definitions for your Go structs and methods
- Native Dialogs & Menus
- Native Dark / Light mode support
- Supports modern translucency and "frosted window" effects
- Unified eventing system between Go and Javascript
- Powerful cli tool to quickly generate and build your projects
- Multiplatform
- Uses native rendering engines - _no embedded browser_!

### Roadmap

The project roadmap may be found [here](https://github.com/wailsapp/wails/discussions/1484). Please consult
this before open up an enhancement request.

## Sponsors

This project is supported by these kind people / companies:

<a href="https://github.com/sponsors/leaanthony" style="width:100px;">
  <img src="../website/static/img/silver%20sponsor.webp" width="100"/>
</a>
<a href="https://github.com/selvindev" style="width:100px;">
  <img src="https://github.com/selvindev.png?size=100" width="100"/>
</a>
<br/>
<br/>
<a href="https://github.com/sponsors/leaanthony" style="width:100px;">
  <img src="../website/static/img/bronze%20sponsor.webp" width="100"/>
</a>

<a href="https://github.com/codydbentley" style="width:100px">
  <img src="https://github.com/codydbentley.png?size=100" width="100"/>
</a>
<a href="https://www.easywebadv.it/" style="width:100px">
  <img src="../website/static/img/easyweb.png" width="100"/>
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
<a href="https://github.com/KiddoV" style="width:45px">
  <img src="https://github.com/KiddoV.png?size=45" width="45"/>
</a>

## Getting Started

The installation instructions are on the [official website](https://wails.io/docs/gettingstarted/installation).

## FAQ

- Is this an alternative to Electron?

  Depends on your requirements. It's designed to make it easy for Go programmers to make lightweight desktop
  applications or add a frontend to their existing applications. Wails does offer native elements such as menus
  and dialogs, so it could be considered a lightweight electron alternative.

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

## License

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fwailsapp%2Fwails.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fwailsapp%2Fwails?ref=badge_large)

## Inspiration

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
