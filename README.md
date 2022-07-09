<p align="center" style="text-align: center">
  <img src="logo-universal.png" width="55%"><br/>
</p>
<p align="center">
  Build desktop applications using Go & Web Technologies.<br/><br/>
  <a href="https://github.com/wailsapp/wails/blob/master/LICENSE">
    <img src="https://img.shields.io/badge/License-MIT-blue.svg">
  </a>
  <a href="https://goreportcard.com/report/github.com/wailsapp/wails">
    <img src="https://goreportcard.com/badge/github.com/wailsapp/wails"/>
  </a>
  <a href="http://godoc.org/github.com/wailsapp/wails">
    <img src="https://img.shields.io/badge/godoc-reference-blue.svg"/>
  </a>
  <a href="https://www.codefactor.io/repository/github/wailsapp/wails">
    <img src="https://www.codefactor.io/repository/github/wailsapp/wails/badge" alt="CodeFactor" />
  </a>
  <a href="https://github.com/wailsapp/wails/issues">
    <img src="https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat" alt="CodeFactor" />
  </a>
  <a href="https://app.fossa.io/projects/git%2Bgithub.com%2Fwailsapp%2Fwails?ref=badge_shield" alt="FOSSA Status">
    <img src="https://app.fossa.io/api/projects/git%2Bgithub.com%2Fwailsapp%2Fwails.svg?type=shield"/>
  </a>
  <a href="https://houndci.com">
    <img src="https://img.shields.io/badge/Reviewed_by-Hound-8E64B0.svg"/>
  </a>
  <a href="https://github.com/avelino/awesome-go" rel="nofollow">
    <img src="https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg" alt="Awesome"/>
  </a>
  <a href="https://github.com/wailsapp/wails/workflows/release/badge.svg?branch=master" rel="nofollow">
    <img src="https://github.com/wailsapp/wails/workflows/release/badge.svg?branch=master" alt="Release Pipelines"/>
  </a>
</p>

<span id="nav-1"></span>

## Internationalization

[English](README.md) | [简体中文](README.zh-Hans.md)

<span id="nav-2"></span>

## Table of Contents

<details>
  <summary>Click me to Open/Close the directory listing</summary>

- [1. Internationalization](#nav-1)
- [2. Table of Contents](#nav-2)
- [3. Introduction](#nav-3)
  - [3.1 Official Website](#nav-3-1)
- [4. Features](#nav-4)
- [5. Sponsors](#nav-5)
- [6. Installation](#nav-6)
- [7. FAQ](#nav-8)
- [8. Contributors](#nav-9)
- [9. Special Mentions](#nav-10)
- [10. Special Thanks](#nav-11)

</details>

<span id="nav-3"></span>

## Introduction

The traditional method of providing web interfaces to Go programs is via a built-in web server. Wails offers a different
approach: it provides the ability to wrap both Go code and a web frontend into a single binary. Tools are provided to
make this easy for you by handling project creation, compilation and bundling. All you have to do is get creative!

<span id="nav-3-1"></span>
<hr/>
<h3><strong>PLEASE NOTE: As we are approaching the v2 release, we are not accepting any new feature requests or bug reports for v1. If you have a critical issue, please open a ticket and state why it is critical.</strong></h3>
<hr/>

### Version 2

Wails v2 has been released in Beta for all 3 platforms. Check out the [new website](https://wails.io) if you are
interested in trying it out.

### Legacy v1 Website

The legacy v1 docs can be found at [https://wails.app](https://wails.app).

<span id="nav-4"></span>

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
- Uses native rendering engines - *no embedded browser*!

<span id="nav-5"></span>

## Sponsors

This project is supported by these kind people / companies:

<a href="https://github.com/sponsors/leaanthony" style="width:100px;">
  <img src="/img/silver%20sponsor.png" width="100"/>
</a>
<a href="https://github.com/selvindev" style="width:100px;">
  <img src="https://github.com/selvindev.png?size=100" width="100"/>
</a>
<br/>
<br/>
<a href="https://github.com/sponsors/leaanthony" style="width:100px;">
  <img src="img/bronze%20sponsor.png" width="100"/>
</a>
<a href="https://github.com/snider" style="width:100px;">
  <img src="https://github.com/snider.png?size=100" width="100"/>
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
<a href="https://www.jetbrains.com?from=Wails" style="width:100px">
  <img src="/assets/images/jetbrains-grayscale.png" width="100"/>
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
<a href="https://github.com/EdenNetworkItalia" style="width:65px">
  <img src="https://github.com/EdenNetworkItalia.png?size=65" width="65"/>
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

<span id="nav-6"></span>

## Roadmap

The project roadmap may be found [here](https://github.com/wailsapp/wails/discussions/1484). Please consult
this before open up an enhancement request.

## Installation

The installation instructions are on the [official website](https://wails.io/docs/gettingstarted/installation).

<span id="nav-8"></span>

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

<span id="nav-9"></span>

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
<a href="https://github.com/hi019"><img src="https://github.com/hi019.png?size=40" width="40"/></a>
<a href="https://github.com/Igogrek"><img src="https://github.com/Igogrek.png?size=40" width="40"/></a>
<a href="https://github.com/aschey"><img src="https://github.com/aschey.png?size=40" width="40"/></a>
<a href="https://github.com/akhudek"><img src="https://github.com/akhudek.png?size=40" width="40"/></a>
<a href="https://github.com/s12chung"><img src="https://github.com/s12chung.png?size=40" width="40"/></a>

<span id="nav-10"></span>

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

<span id="nav-11"></span>

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
