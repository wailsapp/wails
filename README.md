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
<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->
[![All Contributors](https://img.shields.io/badge/all_contributors-0-orange.svg?style=flat-square)](#contributors-)
<!-- ALL-CONTRIBUTORS-BADGE:END -->
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
<a href="https://github.com/boostchicken" style="width:65px">
  <img src="https://github.com/boostchicken.png?size=65" width="65"/>
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

## Stargazers over time

[![Stargazers over time](https://starchart.cc/wailsapp/wails.svg)](https://starchart.cc/wailsapp/wails)

<span id="nav-9"></span>

## Contributors

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tr>
    <td align="center"><a href="https://github.com/leaanthony"><img src="https://avatars.githubusercontent.com/u/1943904?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Lea Anthony</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=leaanthony" title="Code">💻</a> <a href="#ideas-leaanthony" title="Ideas, Planning, & Feedback">🤔</a> <a href="#design-leaanthony" title="Design">🎨</a> <a href="#content-leaanthony" title="Content">🖋</a> <a href="#example-leaanthony" title="Examples">💡</a> <a href="#mentoring-leaanthony" title="Mentoring">🧑‍🏫</a> <a href="#projectManagement-leaanthony" title="Project Management">📆</a> <a href="#tool-leaanthony" title="Tools">🔧</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3Aleaanthony" title="Bug reports">🐛</a> <a href="#blog-leaanthony" title="Blogposts">📝</a> <a href="#maintenance-leaanthony" title="Maintenance">🚧</a> <a href="#platform-leaanthony" title="Packaging/porting to new platform">📦</a> <a href="https://github.com/wailsapp/wails/pulls?q=is%3Apr+reviewed-by%3Aleaanthony" title="Reviewed Pull Requests">👀</a> <a href="#question-leaanthony" title="Answering Questions">💬</a> <a href="#research-leaanthony" title="Research">🔬</a> <a href="https://github.com/wailsapp/wails/commits?author=leaanthony" title="Tests">⚠️</a> <a href="#tutorial-leaanthony" title="Tutorials">✅</a> <a href="#talk-leaanthony" title="Talks">📢</a> <a href="https://github.com/wailsapp/wails/pulls?q=is%3Apr+reviewed-by%3Aleaanthony" title="Reviewed Pull Requests">👀</a></td>
    <td align="center"><a href="https://github.com/stffabi"><img src="https://avatars.githubusercontent.com/u/9464631?v=4?s=75" width="75px;" alt=""/><br /><sub><b>stffabi</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=stffabi" title="Code">💻</a> <a href="#ideas-stffabi" title="Ideas, Planning, & Feedback">🤔</a> <a href="#design-stffabi" title="Design">🎨</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3Astffabi" title="Bug reports">🐛</a> <a href="#maintenance-stffabi" title="Maintenance">🚧</a> <a href="#platform-stffabi" title="Packaging/porting to new platform">📦</a> <a href="https://github.com/wailsapp/wails/pulls?q=is%3Apr+reviewed-by%3Astffabi" title="Reviewed Pull Requests">👀</a> <a href="#question-stffabi" title="Answering Questions">💬</a> <a href="#research-stffabi" title="Research">🔬</a> <a href="https://github.com/wailsapp/wails/pulls?q=is%3Apr+reviewed-by%3Astffabi" title="Reviewed Pull Requests">👀</a> <a href="https://github.com/wailsapp/wails/commits?author=stffabi" title="Tests">⚠️</a></td>
    <td align="center"><a href="https://github.com/tmclane"><img src="https://avatars.githubusercontent.com/u/511975?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Travis McLane</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=tmclane" title="Code">💻</a></td>
    <td align="center"><a href="https://misitebao.com/"><img src="https://avatars.githubusercontent.com/u/28185258?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Misite Bao</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=misitebao" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/apps/dependabot"><img src="https://avatars.githubusercontent.com/in/29110?v=4?s=75" width="75px;" alt=""/><br /><sub><b>dependabot[bot]</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=dependabot[bot]" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/bh90210"><img src="https://avatars.githubusercontent.com/u/22690219?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Byron Chris</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=bh90210" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/konez2k"><img src="https://avatars.githubusercontent.com/u/32417933?v=4?s=75" width="75px;" alt=""/><br /><sub><b>konez2k</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=konez2k" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/dedo1911"><img src="https://avatars.githubusercontent.com/u/1364496?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Dario Emerson</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=dedo1911" title="Code">💻</a></td>
    <td align="center"><a href="https://ianmjones.com/"><img src="https://avatars.githubusercontent.com/u/4710?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Ian M. Jones</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=ianmjones" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/marktohark"><img src="https://avatars.githubusercontent.com/u/19359934?v=4?s=75" width="75px;" alt=""/><br /><sub><b>marktohark</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=marktohark" title="Code">💻</a></td>
  </tr>
  <tr>
    <td align="center"><a href="https://github.com/rh12503"><img src="https://avatars.githubusercontent.com/u/48951973?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Ryan H</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=rh12503" title="Code">💻</a></td>
    <td align="center"><a href="https://codybentley.dev/"><img src="https://avatars.githubusercontent.com/u/6968902?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Cody Bentley</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=codydbentley" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/napalu"><img src="https://avatars.githubusercontent.com/u/6690378?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Florent</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=napalu" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/akhudek"><img src="https://avatars.githubusercontent.com/u/147633?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Alexander Hudek</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=akhudek" title="Code">💻</a></td>
    <td align="center"><a href="https://twitter.com/timkippdev"><img src="https://avatars.githubusercontent.com/u/37030721?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Tim Kipp</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=timkippdev" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/gelleson"><img src="https://avatars.githubusercontent.com/u/44272887?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Altynbek Kaliakbarov</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=gelleson" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/Chronophylos"><img src="https://avatars.githubusercontent.com/u/14890588?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Nikolai Zimmermann</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=Chronophylos" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/k-muchmore"><img src="https://avatars.githubusercontent.com/u/16393095?v=4?s=75" width="75px;" alt=""/><br /><sub><b>k-muchmore</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=k-muchmore" title="Code">💻</a></td>
    <td align="center"><a href="https://peakd.com/@snider"><img src="https://avatars.githubusercontent.com/u/631881?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Snider</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=Snider" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/albert-sun"><img src="https://avatars.githubusercontent.com/u/54585592?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Albert Sun</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=albert-sun" title="Code">💻</a></td>
  </tr>
  <tr>
    <td align="center"><a href="https://github.com/adalessa"><img src="https://avatars.githubusercontent.com/u/7914601?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Ariel</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=adalessa" title="Code">💻</a></td>
    <td align="center"><a href="https://triplebits.com/"><img src="https://avatars.githubusercontent.com/u/4365245?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Ilgıt Yıldırım</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=ilgityildirim" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/Vaelatern"><img src="https://avatars.githubusercontent.com/u/7906072?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Toyam Cox</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=Vaelatern" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/hi019"><img src="https://avatars.githubusercontent.com/u/65871571?v=4?s=75" width="75px;" alt=""/><br /><sub><b>hi019</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=hi019" title="Code">💻</a></td>
    <td align="center"><a href="https://artooro.com/"><img src="https://avatars.githubusercontent.com/u/393395?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Arthur Wiebe</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=artooro" title="Code">💻</a></td>
    <td align="center"><a href="https://sectcs.com/"><img src="https://avatars.githubusercontent.com/u/16898783?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Balakrishna Prasad Ganne</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=aayush420" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/BillBuilt"><img src="https://avatars.githubusercontent.com/u/28831382?v=4?s=75" width="75px;" alt=""/><br /><sub><b>BillBuilt</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=BillBuilt" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/Juneezee"><img src="https://avatars.githubusercontent.com/u/20135478?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Eng Zer Jun</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=Juneezee" title="Code">💻</a></td>
    <td align="center"><a href="https://lgiki.net/"><img src="https://avatars.githubusercontent.com/u/20807713?v=4?s=75" width="75px;" alt=""/><br /><sub><b>LGiki</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=LGiki" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/lontten"><img src="https://avatars.githubusercontent.com/u/30745595?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Lontten</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=lontten" title="Code">💻</a></td>
  </tr>
  <tr>
    <td align="center"><a href="https://github.com/phoenix147"><img src="https://avatars.githubusercontent.com/u/809358?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Lukas Crepaz</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=phoenix147" title="Code">💻</a></td>
    <td align="center"><a href="https://utf9k.net/"><img src="https://avatars.githubusercontent.com/u/14816406?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Marcus Crane</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=marcus-crane" title="Code">💻</a></td>
    <td align="center"><a href="https://qaisjp.com/"><img src="https://avatars.githubusercontent.com/u/923242?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Qais Patankar</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=qaisjp" title="Code">💻</a></td>
    <td align="center"><a href="https://wakefulcloud.dev/"><img src="https://avatars.githubusercontent.com/u/38930607?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Wakeful-Cloud</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=Wakeful-Cloud" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/Lyimmi"><img src="https://avatars.githubusercontent.com/u/8627125?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Zámbó, Levente</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=Lyimmi" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/Ironpark"><img src="https://avatars.githubusercontent.com/u/4973597?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Ironpark</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=Ironpark" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/mondy"><img src="https://avatars.githubusercontent.com/u/3961824?v=4?s=75" width="75px;" alt=""/><br /><sub><b>mondy</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=mondy" title="Code">💻</a></td>
    <td align="center"><a href="https://ryben.dev/"><img src="https://avatars.githubusercontent.com/u/6241454?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Benjamin Ryan</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=redraskal" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/fallendusk"><img src="https://avatars.githubusercontent.com/u/565631?v=4?s=75" width="75px;" alt=""/><br /><sub><b>fallendusk</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=fallendusk" title="Code">💻</a></td>
    <td align="center"><a href="https://twitter.com/matryer"><img src="https://avatars.githubusercontent.com/u/101659?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Mat Ryer</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=matryer" title="Code">💻</a></td>
  </tr>
  <tr>
    <td align="center"><a href="https://github.com/abtin"><img src="https://avatars.githubusercontent.com/u/441372?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Abtin</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=abtin" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/lanzafame"><img src="https://avatars.githubusercontent.com/u/5924712?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Adrian Lanzafame</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=lanzafame" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/polikow"><img src="https://avatars.githubusercontent.com/u/58259700?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Aleksey Polyakov</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=polikow" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/alexmat"><img src="https://avatars.githubusercontent.com/u/745421?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Alexander Matviychuk</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=alexmat" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/AlienRecall"><img src="https://avatars.githubusercontent.com/u/68950287?v=4?s=75" width="75px;" alt=""/><br /><sub><b>AlienRecall</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=AlienRecall" title="Code">💻</a></td>
    <td align="center"><a href="https://blog.checkyo.tech/"><img src="https://avatars.githubusercontent.com/u/17457975?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Aman</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=achhabra2" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/amaury-tobias"><img src="https://avatars.githubusercontent.com/u/37311888?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Amaury Tobias Quiroz</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=amaury-tobias" title="Code">💻</a></td>
    <td align="center"><a href="http://blog.nms.de/"><img src="https://avatars.githubusercontent.com/u/51517?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Andreas Wenk</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=andywenk" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/stankovic98"><img src="https://avatars.githubusercontent.com/u/29852655?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Antonio Stanković</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=stankovic98" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/antimatter96"><img src="https://avatars.githubusercontent.com/u/12068176?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Arpit Jain</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=antimatter96" title="Code">💻</a></td>
  </tr>
  <tr>
    <td align="center"><a href="https://github.com/aschey"><img src="https://avatars.githubusercontent.com/u/5882266?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Austin Schey</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=aschey" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/benjamin-thomas"><img src="https://avatars.githubusercontent.com/u/1557738?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Benjamin Thomas</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=benjamin-thomas" title="Code">💻</a></td>
    <td align="center"><a href="https://www.bertramtruong.com/"><img src="https://avatars.githubusercontent.com/u/1100843?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Bertram Truong</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=bt" title="Code">💻</a></td>
    <td align="center"><a href="http://techwizworld.net/"><img src="https://avatars.githubusercontent.com/u/175873?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Blake Bourque</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=TechplexEngineer" title="Code">💻</a></td>
    <td align="center"><a href="http://vk.com/raitonoberu"><img src="https://avatars.githubusercontent.com/u/64320078?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Denis</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=raitonoberu" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/diogox"><img src="https://avatars.githubusercontent.com/u/13244408?v=4?s=75" width="75px;" alt=""/><br /><sub><b>diogox</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=diogox" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/kyoto44"><img src="https://avatars.githubusercontent.com/u/17720761?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Dmitry Gomzyakov</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=kyoto44" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/edwardbrowncross"><img src="https://avatars.githubusercontent.com/u/35063432?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Edward Browncross</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=edwardbrowncross" title="Code">💻</a></td>
    <td align="center"><a href="http://pr0gramming.ca/"><img src="https://avatars.githubusercontent.com/u/14944216?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Elie Grenon</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=elie-g" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/fdidron"><img src="https://avatars.githubusercontent.com/u/1848786?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Florian Didron</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=fdidron" title="Code">💻</a></td>
  </tr>
  <tr>
    <td align="center"><a href="https://github.com/GargantuaX"><img src="https://avatars.githubusercontent.com/u/14013111?v=4?s=75" width="75px;" alt=""/><br /><sub><b>GargantuaX</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=GargantuaX" title="Code">💻</a></td>
    <td align="center"><a href="https://bednya.ga/"><img src="https://avatars.githubusercontent.com/u/12101721?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Igor Minin</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=Igogrek" title="Code">💻</a></td>
    <td align="center"><a href="https://www.jae-sung.com/"><img src="https://avatars.githubusercontent.com/u/39658806?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Jae-Sung Lee</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=jaesung9507" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/Jarek-SRT"><img src="https://avatars.githubusercontent.com/u/3391365?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Jarek</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=Jarek-SRT" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/Junkher"><img src="https://avatars.githubusercontent.com/u/85776620?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Junker</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=Junkher" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/kraney"><img src="https://avatars.githubusercontent.com/u/5760081?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Kris Raney</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=kraney" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/LukenSkyne"><img src="https://avatars.githubusercontent.com/u/29918069?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Luken</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=LukenSkyne" title="Code">💻</a></td>
    <td align="center"><a href="https://markstenglein.com/"><img src="https://avatars.githubusercontent.com/u/9255772?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Mark Stenglein</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=ocelotsloth" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/buddyabaddon"><img src="https://avatars.githubusercontent.com/u/33861511?v=4?s=75" width="75px;" alt=""/><br /><sub><b>buddyabaddon</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=buddyabaddon" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/MikeSchaap"><img src="https://avatars.githubusercontent.com/u/35368821?v=4?s=75" width="75px;" alt=""/><br /><sub><b>MikeSchaap</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=MikeSchaap" title="Code">💻</a></td>
  </tr>
  <tr>
    <td align="center"><a href="https://github.com/Orijhins"><img src="https://avatars.githubusercontent.com/u/47521598?v=4?s=75" width="75px;" alt=""/><br /><sub><b>NYSSEN Michaël</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=Orijhins" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/NanoNik"><img src="https://avatars.githubusercontent.com/u/11991329?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Nan0</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=NanoNik" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/marcio199226"><img src="https://avatars.githubusercontent.com/u/10244404?v=4?s=75" width="75px;" alt=""/><br /><sub><b>oskar</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=marcio199226" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/pierrejoye"><img src="https://avatars.githubusercontent.com/u/282408?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Pierre Joye</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=pierrejoye" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/Rested"><img src="https://avatars.githubusercontent.com/u/2003608?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Reuben Thomas-Davis</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=Rested" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/mewmew"><img src="https://avatars.githubusercontent.com/u/1414531?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Robin</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=mewmew" title="Code">💻</a></td>
    <td align="center"><a href="https://threema.id/YSB3TVF7"><img src="https://avatars.githubusercontent.com/u/70367451?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Sebastian Bauer</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=sebastian0x62" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/sidwebworks"><img src="https://avatars.githubusercontent.com/u/58144379?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Sidharth Rathi</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=sidwebworks" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/sithembiso"><img src="https://avatars.githubusercontent.com/u/6559905?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Sithembiso Khumalo</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=sithembiso" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/LanguageAgnostic"><img src="https://avatars.githubusercontent.com/u/19310562?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Soheib El-Harrache</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=LanguageAgnostic" title="Code">💻</a></td>
  </tr>
  <tr>
    <td align="center"><a href="https://www.sophieau.com/"><img src="https://avatars.githubusercontent.com/u/11145039?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Sophie Au</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=SophieAu" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/stefpap"><img src="https://avatars.githubusercontent.com/u/22637722?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Stefanos Papadakis</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=stefpap" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/s12chung"><img src="https://avatars.githubusercontent.com/u/263394?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Steve Chung</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=s12chung" title="Code">💻</a></td>
    <td align="center"><a href="https://tortloff.de/"><img src="https://avatars.githubusercontent.com/u/41272726?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Timm Ortloff</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=TAINCER" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/tomanagle"><img src="https://avatars.githubusercontent.com/u/8683577?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Tom</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=tomanagle" title="Code">💻</a></td>
    <td align="center"><a href="https://www.linkedin.com/in/valentintrinque"><img src="https://avatars.githubusercontent.com/u/4662842?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Valentin Trinqué</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=ValentinTrinque" title="Code">💻</a></td>
    <td align="center"><a href="https://mattn.kaoriya.net/"><img src="https://avatars.githubusercontent.com/u/10111?v=4?s=75" width="75px;" alt=""/><br /><sub><b>mattn</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=mattn" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/bearsh"><img src="https://avatars.githubusercontent.com/u/1089356?v=4?s=75" width="75px;" alt=""/><br /><sub><b>bearsh</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=bearsh" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/chenxiao1990"><img src="https://avatars.githubusercontent.com/u/16933565?v=4?s=75" width="75px;" alt=""/><br /><sub><b>chenxiao</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=chenxiao1990" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/fengweiqiang"><img src="https://avatars.githubusercontent.com/u/22905300?v=4?s=75" width="75px;" alt=""/><br /><sub><b>fengweiqiang</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=fengweiqiang" title="Code">💻</a></td>
  </tr>
  <tr>
    <td align="center"><a href="https://github.com/flin7"><img src="https://avatars.githubusercontent.com/u/58138185?v=4?s=75" width="75px;" alt=""/><br /><sub><b>flin7</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=flin7" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/fred21O4"><img src="https://avatars.githubusercontent.com/u/67189813?v=4?s=75" width="75px;" alt=""/><br /><sub><b>fred21O4</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=fred21O4" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/gardc"><img src="https://avatars.githubusercontent.com/u/41453409?v=4?s=75" width="75px;" alt=""/><br /><sub><b>gardc</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=gardc" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/rayshoo"><img src="https://avatars.githubusercontent.com/u/52561899?v=4?s=75" width="75px;" alt=""/><br /><sub><b>rayshoo</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=rayshoo" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/Yz4230"><img src="https://avatars.githubusercontent.com/u/38999742?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Ishiyama Yuzuki</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=Yz4230" title="Code">💻</a></td>
    <td align="center"><a href="https://baiyue.one/"><img src="https://avatars.githubusercontent.com/u/43716063?v=4?s=75" width="75px;" alt=""/><br /><sub><b>佰阅</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=Baiyuetribe" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/daodao97"><img src="https://avatars.githubusercontent.com/u/15009280?v=4?s=75" width="75px;" alt=""/><br /><sub><b>刀刀</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=daodao97" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/jicg"><img src="https://avatars.githubusercontent.com/u/6479672?v=4?s=75" width="75px;" alt=""/><br /><sub><b>归位</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=jicg" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/skamensky"><img src="https://avatars.githubusercontent.com/u/19151369?v=4?s=75" width="75px;" alt=""/><br /><sub><b>skamensky</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=skamensky" title="Code">💻</a> <a href="#ideas-skamensky" title="Ideas, Planning, & Feedback">🤔</a> <a href="https://github.com/wailsapp/wails/commits?author=skamensky" title="Documentation">📖</a></td>
  </tr>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

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
