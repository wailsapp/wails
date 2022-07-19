<p align="center" style="text-align: center">
   <img src="logo-universal.png" width="55%"><br/>
</p>
<p align="center">
  使用 Go 和 Web 技术构建桌面应用程序。<br/><br/>
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
<a href='#contributors'><img src='https://img.shields.io/badge/all_contributors-1-orange.svg?style=flat' alt='Contributors'/></a>
<!-- ALL-CONTRIBUTORS-BADGE:END -->
</p>

<span id="nav-1"></span>

## 国际化

[English](README.md) | [简体中文](README.zh-Hans.md)

<span id="nav-2"></span>

## 内容目录

<details>
  <summary>点我 打开/关闭 目录列表</summary>

- [1. 国际化](#nav-1)
- [2. 内容目录](#nav-2)
- [3. 项目介绍](#nav-3)
  - [3.1 官方网站](#nav-3-1)
- [4. 功能](#nav-4)
- [5. 赞助商](#nav-5)
- [6. 安装](#nav-6)
  - [6.1 MacOS](#nav-6-1)
  - [6.2 Linux](#nav-6-2)
    - [6.2.1 Debian/Ubuntu](#nav-6-2-1)
    - [6.2.2 Arch Linux / ArchLabs / Ctlos Linux](#nav-6-2-2)
    - [6.2.3 Centos](#nav-6-2-3)
    - [6.2.4 Fedora](#nav-6-2-4)
    - [6.2.5 VoidLinux & VoidLinux-musl](#nav-6-2-5)
    - [6.2.6 Gentoo](#nav-6-2-6)
  - [6.3 Windows](#nav-6-3)
- [7. 使用方法](#nav-7)
  - [7.1 下一步](#nav-7-1)
- [8. 常见问题](#nav-8)
- [9. 贡献者](#nav-9)
- [10. 特别提及](#nav-10)
- [12. 特别感谢](#nav-11)

</details>

<span id="nav-3"></span>

## 项目介绍

为 Go 程序提供 Web 界面的传统方法是通过内置 Web 服务器。Wails 提供了一种不同的方法：它提供了将 Go 代码和 Web
前端一起打包成单个二进制文件的能力。通过提供的工具，可以很轻松的完成项目的创建、编译和打包。你所要做的就是发挥想象力！

<span id="nav-3-1"></span>

### 官方网站

#### v1

官方文档可以在 [https://wails.app](https://wails.app) 中找到。

#### v2

Wails v2 已针对所有 3 个平台发布了 Beta 版。如果您有兴趣尝试一下，请查看[新网站](https://wails.io)。

镜像网站：

- [中国大陆镜像站点 - https://wails.top](https://wails.top)

<span id="nav-4"></span>

## 功能

- 后端使用标准 Go
- 使用任意前端技术构建 UI 界面
- 快速为您的 Go 应用生成 Vue、Vuetify、React 前端代码
- 通过简单的绑定命令将 Go 方法暴露到前端
- 使用原生渲染引擎 - 无嵌入式浏览器
- 共享事件系统
- 原生文件系统对话框
- 强大的命令行工具
- 跨多个平台

<span id="nav-5"></span>

## 赞助商

这个项目由以下这些人或者公司支持：

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
directory
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
<span id="nav-6"></span>

## 安装

Wails 使用 cgo 与原生渲染引擎结合，因此需要依赖一些平台的库以及 Go 的安装。基本要求是：

- Go 1.16
- npm

<span id="nav-6-1"></span>

### MacOS

请确保已安装 xcode 命令行工具。这可以通过运行下面的命令来完成：

`xcode-select --install`

<span id="nav-6-2"></span>

### Linux

<span id="nav-6-2-1"></span>

#### Debian/Ubuntu

`sudo apt install libgtk-3-dev libwebkit2gtk-4.0-dev`

_Debian: 8, 9, 10_

_Ubuntu: 16.04, 18.04, 19.04_

_也成功测试了: Zorin 15, Parrot 4.7, Linuxmint 19, Elementary 5, Kali, Neon_, Pop!\_OS

<span id="nav-6-2-2"></span>

#### Arch Linux / ArchLabs / Ctlos Linux

`sudo pacman -S webkit2gtk gtk3`

_也成功测试了: Manjaro & ArcoLinux_

<span id="nav-6-2-3"></span>

#### Centos

`sudo yum install webkitgtk3-devel gtk3-devel`

_CentOS 6, 7_

<span id="nav-6-2-4"></span>

#### Fedora

`sudo yum install webkit2gtk3-devel gtk3-devel`

_Fedora 29, 30_

<span id="nav-6-2-5"></span>

#### VoidLinux & VoidLinux-musl

`xbps-install gtk+3-devel webkit2gtk-devel`

<span id="nav-6-2-6"></span>

#### Gentoo

`sudo emerge gtk+:3 webkit-gtk`

<span id="nav-6-3"></span>

### Windows

Windows 需要 GCC 和相关工具。 建议从 [http://tdm-gcc.tdragon.net/download](http://tdm-gcc.tdragon.net/download) 下载， 安装完成，您就可以开始了。

<span id="nav-7"></span>

## 使用方法

**确保 Go modules 是开启的：GO111MODULE=on 并且 go/bin 在您的 PATH 变量中。**

安装很简单，运行以下命令：

```
go get -u github.com/wailsapp/wails/cmd/wails
```

<span id="nav-7-1"></span>

### 下一步

建议在此时阅读 [https://wails.app](https://wails.app) 上面的文档.

<span id="nav-8"></span>

## 常见问题

- 它是 Electron 的替代品吗?

  取决于您的要求。它旨在使 Go 程序员可以轻松制作轻量级桌面应用程序或在其现有应用程序中添加前端。尽管 Wails 当前不提供对诸如菜单之类的原生元素的钩子，但将来可能会改变。

- 这个项目针对的是哪些人?

  希望将 HTML / JS / CSS 前端与其应用程序捆绑在一起的程序员，而不是借助创建服务并打开浏览器进行查看的方式。

- 名字怎么来的?

  当我看到 WebView 时，我想"我真正想要的是围绕构建 WebView 应用程序工作，有点像 Rails 对于 Ruby"。因此，最初它是一个文字游戏（Webview on
  Rails）。碰巧也是我来自的 [国家](https://en.wikipedia.org/wiki/Wales) 的英文名字的同音。所以就是它了。

<span id="nav-9"></span>

## 贡献者

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tr>
    <td align="center"><a href="https://github.com/leaanthony"><img src="https://avatars.githubusercontent.com/u/1943904?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Lea Anthony</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=leaanthony" title="Code">💻</a> <a href="#ideas-leaanthony" title="Ideas, Planning, & Feedback">🤔</a> <a href="#design-leaanthony" title="Design">🎨</a> <a href="#content-leaanthony" title="Content">🖋</a> <a href="#example-leaanthony" title="Examples">💡</a> <a href="#mentoring-leaanthony" title="Mentoring">🧑‍🏫</a> <a href="#projectManagement-leaanthony" title="Project Management">📆</a> <a href="#tool-leaanthony" title="Tools">🔧</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3Aleaanthony" title="Bug reports">🐛</a> <a href="#blog-leaanthony" title="Blogposts">📝</a> <a href="#maintenance-leaanthony" title="Maintenance">🚧</a> <a href="#platform-leaanthony" title="Packaging/porting to new platform">📦</a> <a href="https://github.com/wailsapp/wails/pulls?q=is%3Apr+reviewed-by%3Aleaanthony" title="Reviewed Pull Requests">👀</a> <a href="#question-leaanthony" title="Answering Questions">💬</a> <a href="#research-leaanthony" title="Research">🔬</a> <a href="https://github.com/wailsapp/wails/commits?author=leaanthony" title="Tests">⚠️</a> <a href="#tutorial-leaanthony" title="Tutorials">✅</a> <a href="#talk-leaanthony" title="Talks">📢</a> <a href="https://github.com/wailsapp/wails/pulls?q=is%3Apr+reviewed-by%3Aleaanthony" title="Reviewed Pull Requests">👀</a> <a href="https://github.com/wailsapp/wails/commits?author=leaanthony" title="Documentation">📖</a></td>
    <td align="center"><a href="https://github.com/stffabi"><img src="https://avatars.githubusercontent.com/u/9464631?v=4?s=75" width="75px;" alt=""/><br /><sub><b>stffabi</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=stffabi" title="Code">💻</a> <a href="#ideas-stffabi" title="Ideas, Planning, & Feedback">🤔</a> <a href="#design-stffabi" title="Design">🎨</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3Astffabi" title="Bug reports">🐛</a> <a href="#maintenance-stffabi" title="Maintenance">🚧</a> <a href="#platform-stffabi" title="Packaging/porting to new platform">📦</a> <a href="https://github.com/wailsapp/wails/pulls?q=is%3Apr+reviewed-by%3Astffabi" title="Reviewed Pull Requests">👀</a> <a href="#question-stffabi" title="Answering Questions">💬</a> <a href="#research-stffabi" title="Research">🔬</a> <a href="https://github.com/wailsapp/wails/pulls?q=is%3Apr+reviewed-by%3Astffabi" title="Reviewed Pull Requests">👀</a> <a href="https://github.com/wailsapp/wails/commits?author=stffabi" title="Documentation">📖</a> <a href="https://github.com/wailsapp/wails/commits?author=stffabi" title="Tests">⚠️</a></td>
    <td align="center"><a href="https://github.com/tmclane"><img src="https://avatars.githubusercontent.com/u/511975?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Travis McLane</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=tmclane" title="Code">💻</a> <a href="#research-tmclane" title="Research">🔬</a> <a href="#platform-tmclane" title="Packaging/porting to new platform">📦</a> <a href="#ideas-tmclane" title="Ideas, Planning, & Feedback">🤔</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3Atmclane" title="Bug reports">🐛</a> <a href="https://github.com/wailsapp/wails/pulls?q=is%3Apr+reviewed-by%3Atmclane" title="Reviewed Pull Requests">👀</a> <a href="https://github.com/wailsapp/wails/commits?author=tmclane" title="Tests">⚠️</a> <a href="#question-tmclane" title="Answering Questions">💬</a> <a href="https://github.com/wailsapp/wails/commits?author=tmclane" title="Documentation">📖</a></td>
    <td align="center"><a href="https://misitebao.com/"><img src="https://avatars.githubusercontent.com/u/28185258?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Misite Bao</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=misitebao" title="Documentation">📖</a> <a href="#translation-misitebao" title="Translation">🌍</a> <a href="#research-misitebao" title="Research">🔬</a> <a href="#maintenance-misitebao" title="Maintenance">🚧</a></td>
    <td align="center"><a href="https://github.com/bh90210"><img src="https://avatars.githubusercontent.com/u/22690219?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Byron Chris</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=bh90210" title="Code">💻</a> <a href="#research-bh90210" title="Research">🔬</a> <a href="#maintenance-bh90210" title="Maintenance">🚧</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3Abh90210" title="Bug reports">🐛</a> <a href="https://github.com/wailsapp/wails/pulls?q=is%3Apr+reviewed-by%3Abh90210" title="Reviewed Pull Requests">👀</a> <a href="https://github.com/wailsapp/wails/commits?author=bh90210" title="Tests">⚠️</a> <a href="#question-bh90210" title="Answering Questions">💬</a> <a href="#ideas-bh90210" title="Ideas, Planning, & Feedback">🤔</a> <a href="#design-bh90210" title="Design">🎨</a> <a href="#platform-bh90210" title="Packaging/porting to new platform">📦</a> <a href="#infra-bh90210" title="Infrastructure (Hosting, Build-Tools, etc)">🚇</a></td>
    <td align="center"><a href="https://github.com/konez2k"><img src="https://avatars.githubusercontent.com/u/32417933?v=4?s=75" width="75px;" alt=""/><br /><sub><b>konez2k</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=konez2k" title="Code">💻</a> <a href="#platform-konez2k" title="Packaging/porting to new platform">📦</a> <a href="#ideas-konez2k" title="Ideas, Planning, & Feedback">🤔</a></td>
    <td align="center"><a href="https://github.com/dedo1911"><img src="https://avatars.githubusercontent.com/u/1364496?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Dario Emerson</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=dedo1911" title="Code">💻</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3Adedo1911" title="Bug reports">🐛</a> <a href="#ideas-dedo1911" title="Ideas, Planning, & Feedback">🤔</a> <a href="https://github.com/wailsapp/wails/commits?author=dedo1911" title="Tests">⚠️</a></td>
    <td align="center"><a href="https://ianmjones.com/"><img src="https://avatars.githubusercontent.com/u/4710?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Ian M. Jones</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=ianmjones" title="Code">💻</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3Aianmjones" title="Bug reports">🐛</a> <a href="#ideas-ianmjones" title="Ideas, Planning, & Feedback">🤔</a> <a href="https://github.com/wailsapp/wails/commits?author=ianmjones" title="Tests">⚠️</a> <a href="https://github.com/wailsapp/wails/pulls?q=is%3Apr+reviewed-by%3Aianmjones" title="Reviewed Pull Requests">👀</a> <a href="#platform-ianmjones" title="Packaging/porting to new platform">📦</a></td>
  </tr>
  <tr>
    <td align="center"><a href="https://github.com/marktohark"><img src="https://avatars.githubusercontent.com/u/19359934?v=4?s=75" width="75px;" alt=""/><br /><sub><b>marktohark</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=marktohark" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/rh12503"><img src="https://avatars.githubusercontent.com/u/48951973?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Ryan H</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=rh12503" title="Code">💻</a></td>
    <td align="center"><a href="https://codybentley.dev/"><img src="https://avatars.githubusercontent.com/u/6968902?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Cody Bentley</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=codydbentley" title="Code">💻</a> <a href="#platform-codydbentley" title="Packaging/porting to new platform">📦</a> <a href="#ideas-codydbentley" title="Ideas, Planning, & Feedback">🤔</a> <a href="#financial-codydbentley" title="Financial">💵</a></td>
    <td align="center"><a href="https://github.com/napalu"><img src="https://avatars.githubusercontent.com/u/6690378?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Florent</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=napalu" title="Code">💻</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3Anapalu" title="Bug reports">🐛</a></td>
    <td align="center"><a href="https://github.com/akhudek"><img src="https://avatars.githubusercontent.com/u/147633?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Alexander Hudek</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=akhudek" title="Code">💻</a> <a href="#financial-akhudek" title="Financial">💵</a></td>
    <td align="center"><a href="https://twitter.com/timkippdev"><img src="https://avatars.githubusercontent.com/u/37030721?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Tim Kipp</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=timkippdev" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/gelleson"><img src="https://avatars.githubusercontent.com/u/44272887?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Altynbek Kaliakbarov</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=gelleson" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/Chronophylos"><img src="https://avatars.githubusercontent.com/u/14890588?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Nikolai Zimmermann</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=Chronophylos" title="Code">💻</a></td>
  </tr>
  <tr>
    <td align="center"><a href="https://github.com/k-muchmore"><img src="https://avatars.githubusercontent.com/u/16393095?v=4?s=75" width="75px;" alt=""/><br /><sub><b>k-muchmore</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=k-muchmore" title="Code">💻</a></td>
    <td align="center"><a href="https://peakd.com/@snider"><img src="https://avatars.githubusercontent.com/u/631881?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Snider</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=Snider" title="Code">💻</a> <a href="#ideas-Snider" title="Ideas, Planning, & Feedback">🤔</a> <a href="https://github.com/wailsapp/wails/commits?author=Snider" title="Documentation">📖</a> <a href="#financial-Snider" title="Financial">💵</a></td>
    <td align="center"><a href="https://github.com/albert-sun"><img src="https://avatars.githubusercontent.com/u/54585592?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Albert Sun</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=albert-sun" title="Code">💻</a> <a href="https://github.com/wailsapp/wails/commits?author=albert-sun" title="Tests">⚠️</a></td>
    <td align="center"><a href="https://github.com/adalessa"><img src="https://avatars.githubusercontent.com/u/7914601?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Ariel</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=adalessa" title="Code">💻</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3Aadalessa" title="Bug reports">🐛</a></td>
    <td align="center"><a href="https://triplebits.com/"><img src="https://avatars.githubusercontent.com/u/4365245?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Ilgıt Yıldırım</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=ilgityildirim" title="Code">💻</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3Ailgityildirim" title="Bug reports">🐛</a> <a href="#financial-ilgityildirim" title="Financial">💵</a></td>
    <td align="center"><a href="https://github.com/Vaelatern"><img src="https://avatars.githubusercontent.com/u/7906072?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Toyam Cox</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=Vaelatern" title="Code">💻</a> <a href="#platform-Vaelatern" title="Packaging/porting to new platform">📦</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3AVaelatern" title="Bug reports">🐛</a></td>
    <td align="center"><a href="https://github.com/hi019"><img src="https://avatars.githubusercontent.com/u/65871571?v=4?s=75" width="75px;" alt=""/><br /><sub><b>hi019</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=hi019" title="Code">💻</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3Ahi019" title="Bug reports">🐛</a></td>
    <td align="center"><a href="https://artooro.com/"><img src="https://avatars.githubusercontent.com/u/393395?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Arthur Wiebe</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=artooro" title="Code">💻</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3Aartooro" title="Bug reports">🐛</a></td>
  </tr>
  <tr>
    <td align="center"><a href="https://sectcs.com/"><img src="https://avatars.githubusercontent.com/u/16898783?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Balakrishna Prasad Ganne</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=aayush420" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/BillBuilt"><img src="https://avatars.githubusercontent.com/u/28831382?v=4?s=75" width="75px;" alt=""/><br /><sub><b>BillBuilt</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=BillBuilt" title="Code">💻</a> <a href="#platform-BillBuilt" title="Packaging/porting to new platform">📦</a> <a href="#ideas-BillBuilt" title="Ideas, Planning, & Feedback">🤔</a> <a href="#question-BillBuilt" title="Answering Questions">💬</a> <a href="#financial-BillBuilt" title="Financial">💵</a></td>
    <td align="center"><a href="https://github.com/Juneezee"><img src="https://avatars.githubusercontent.com/u/20135478?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Eng Zer Jun</b></sub></a><br /><a href="#maintenance-Juneezee" title="Maintenance">🚧</a> <a href="https://github.com/wailsapp/wails/commits?author=Juneezee" title="Code">💻</a></td>
    <td align="center"><a href="https://lgiki.net/"><img src="https://avatars.githubusercontent.com/u/20807713?v=4?s=75" width="75px;" alt=""/><br /><sub><b>LGiki</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=LGiki" title="Documentation">📖</a></td>
    <td align="center"><a href="https://github.com/lontten"><img src="https://avatars.githubusercontent.com/u/30745595?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Lontten</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=lontten" title="Documentation">📖</a></td>
    <td align="center"><a href="https://github.com/phoenix147"><img src="https://avatars.githubusercontent.com/u/809358?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Lukas Crepaz</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=phoenix147" title="Code">💻</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3Aphoenix147" title="Bug reports">🐛</a></td>
    <td align="center"><a href="https://utf9k.net/"><img src="https://avatars.githubusercontent.com/u/14816406?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Marcus Crane</b></sub></a><br /><a href="https://github.com/wailsapp/wails/issues?q=author%3Amarcus-crane" title="Bug reports">🐛</a> <a href="https://github.com/wailsapp/wails/commits?author=marcus-crane" title="Documentation">📖</a> <a href="#financial-marcus-crane" title="Financial">💵</a></td>
    <td align="center"><a href="https://qaisjp.com/"><img src="https://avatars.githubusercontent.com/u/923242?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Qais Patankar</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=qaisjp" title="Documentation">📖</a></td>
  </tr>
  <tr>
    <td align="center"><a href="https://wakefulcloud.dev/"><img src="https://avatars.githubusercontent.com/u/38930607?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Wakeful-Cloud</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=Wakeful-Cloud" title="Code">💻</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3AWakeful-Cloud" title="Bug reports">🐛</a></td>
    <td align="center"><a href="https://github.com/Lyimmi"><img src="https://avatars.githubusercontent.com/u/8627125?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Zámbó, Levente</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=Lyimmi" title="Code">💻</a> <a href="#platform-Lyimmi" title="Packaging/porting to new platform">📦</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3ALyimmi" title="Bug reports">🐛</a> <a href="https://github.com/wailsapp/wails/commits?author=Lyimmi" title="Tests">⚠️</a></td>
    <td align="center"><a href="https://github.com/Ironpark"><img src="https://avatars.githubusercontent.com/u/4973597?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Ironpark</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=Ironpark" title="Code">💻</a> <a href="#ideas-Ironpark" title="Ideas, Planning, & Feedback">🤔</a></td>
    <td align="center"><a href="https://github.com/mondy"><img src="https://avatars.githubusercontent.com/u/3961824?v=4?s=75" width="75px;" alt=""/><br /><sub><b>mondy</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=mondy" title="Code">💻</a> <a href="https://github.com/wailsapp/wails/commits?author=mondy" title="Documentation">📖</a></td>
    <td align="center"><a href="https://ryben.dev/"><img src="https://avatars.githubusercontent.com/u/6241454?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Benjamin Ryan</b></sub></a><br /><a href="https://github.com/wailsapp/wails/issues?q=author%3Aredraskal" title="Bug reports">🐛</a></td>
    <td align="center"><a href="https://github.com/fallendusk"><img src="https://avatars.githubusercontent.com/u/565631?v=4?s=75" width="75px;" alt=""/><br /><sub><b>fallendusk</b></sub></a><br /><a href="#platform-fallendusk" title="Packaging/porting to new platform">📦</a> <a href="https://github.com/wailsapp/wails/commits?author=fallendusk" title="Code">💻</a></td>
    <td align="center"><a href="https://twitter.com/matryer"><img src="https://avatars.githubusercontent.com/u/101659?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Mat Ryer</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=matryer" title="Code">💻</a> <a href="#ideas-matryer" title="Ideas, Planning, & Feedback">🤔</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3Amatryer" title="Bug reports">🐛</a></td>
    <td align="center"><a href="https://github.com/abtin"><img src="https://avatars.githubusercontent.com/u/441372?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Abtin</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=abtin" title="Code">💻</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3Aabtin" title="Bug reports">🐛</a></td>
  </tr>
  <tr>
    <td align="center"><a href="https://github.com/lanzafame"><img src="https://avatars.githubusercontent.com/u/5924712?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Adrian Lanzafame</b></sub></a><br /><a href="#platform-lanzafame" title="Packaging/porting to new platform">📦</a> <a href="https://github.com/wailsapp/wails/commits?author=lanzafame" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/polikow"><img src="https://avatars.githubusercontent.com/u/58259700?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Aleksey Polyakov</b></sub></a><br /><a href="https://github.com/wailsapp/wails/issues?q=author%3Apolikow" title="Bug reports">🐛</a> <a href="https://github.com/wailsapp/wails/commits?author=polikow" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/alexmat"><img src="https://avatars.githubusercontent.com/u/745421?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Alexander Matviychuk</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=alexmat" title="Code">💻</a> <a href="#platform-alexmat" title="Packaging/porting to new platform">📦</a></td>
    <td align="center"><a href="https://github.com/AlienRecall"><img src="https://avatars.githubusercontent.com/u/68950287?v=4?s=75" width="75px;" alt=""/><br /><sub><b>AlienRecall</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=AlienRecall" title="Code">💻</a> <a href="#platform-AlienRecall" title="Packaging/porting to new platform">📦</a></td>
    <td align="center"><a href="https://blog.checkyo.tech/"><img src="https://avatars.githubusercontent.com/u/17457975?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Aman</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=achhabra2" title="Documentation">📖</a></td>
    <td align="center"><a href="https://github.com/amaury-tobias"><img src="https://avatars.githubusercontent.com/u/37311888?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Amaury Tobias Quiroz</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=amaury-tobias" title="Code">💻</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3Aamaury-tobias" title="Bug reports">🐛</a></td>
    <td align="center"><a href="http://blog.nms.de/"><img src="https://avatars.githubusercontent.com/u/51517?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Andreas Wenk</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=andywenk" title="Documentation">📖</a></td>
    <td align="center"><a href="https://github.com/stankovic98"><img src="https://avatars.githubusercontent.com/u/29852655?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Antonio Stanković</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=stankovic98" title="Code">💻</a> <a href="#platform-stankovic98" title="Packaging/porting to new platform">📦</a></td>
  </tr>
  <tr>
    <td align="center"><a href="https://github.com/antimatter96"><img src="https://avatars.githubusercontent.com/u/12068176?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Arpit Jain</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=antimatter96" title="Documentation">📖</a></td>
    <td align="center"><a href="https://github.com/aschey"><img src="https://avatars.githubusercontent.com/u/5882266?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Austin Schey</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=aschey" title="Code">💻</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3Aaschey" title="Bug reports">🐛</a></td>
    <td align="center"><a href="https://github.com/benjamin-thomas"><img src="https://avatars.githubusercontent.com/u/1557738?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Benjamin Thomas</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=benjamin-thomas" title="Code">💻</a> <a href="#platform-benjamin-thomas" title="Packaging/porting to new platform">📦</a> <a href="#ideas-benjamin-thomas" title="Ideas, Planning, & Feedback">🤔</a></td>
    <td align="center"><a href="https://www.bertramtruong.com/"><img src="https://avatars.githubusercontent.com/u/1100843?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Bertram Truong</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=bt" title="Code">💻</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3Abt" title="Bug reports">🐛</a></td>
    <td align="center"><a href="http://techwizworld.net/"><img src="https://avatars.githubusercontent.com/u/175873?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Blake Bourque</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=TechplexEngineer" title="Documentation">📖</a></td>
    <td align="center"><a href="http://vk.com/raitonoberu"><img src="https://avatars.githubusercontent.com/u/64320078?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Denis</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=raitonoberu" title="Documentation">📖</a></td>
    <td align="center"><a href="https://github.com/diogox"><img src="https://avatars.githubusercontent.com/u/13244408?v=4?s=75" width="75px;" alt=""/><br /><sub><b>diogox</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=diogox" title="Code">💻</a> <a href="#platform-diogox" title="Packaging/porting to new platform">📦</a></td>
    <td align="center"><a href="https://github.com/kyoto44"><img src="https://avatars.githubusercontent.com/u/17720761?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Dmitry Gomzyakov</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=kyoto44" title="Code">💻</a> <a href="#platform-kyoto44" title="Packaging/porting to new platform">📦</a></td>
  </tr>
  <tr>
    <td align="center"><a href="https://github.com/edwardbrowncross"><img src="https://avatars.githubusercontent.com/u/35063432?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Edward Browncross</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=edwardbrowncross" title="Code">💻</a></td>
    <td align="center"><a href="http://pr0gramming.ca/"><img src="https://avatars.githubusercontent.com/u/14944216?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Elie Grenon</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=elie-g" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/fdidron"><img src="https://avatars.githubusercontent.com/u/1848786?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Florian Didron</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=fdidron" title="Code">💻</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3Afdidron" title="Bug reports">🐛</a> <a href="#ideas-fdidron" title="Ideas, Planning, & Feedback">🤔</a> <a href="https://github.com/wailsapp/wails/commits?author=fdidron" title="Tests">⚠️</a> <a href="https://github.com/wailsapp/wails/pulls?q=is%3Apr+reviewed-by%3Afdidron" title="Reviewed Pull Requests">👀</a> <a href="#platform-fdidron" title="Packaging/porting to new platform">📦</a></td>
    <td align="center"><a href="https://github.com/GargantuaX"><img src="https://avatars.githubusercontent.com/u/14013111?v=4?s=75" width="75px;" alt=""/><br /><sub><b>GargantuaX</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=GargantuaX" title="Documentation">📖</a> <a href="#financial-GargantuaX" title="Financial">💵</a></td>
    <td align="center"><a href="https://bednya.ga/"><img src="https://avatars.githubusercontent.com/u/12101721?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Igor Minin</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=Igogrek" title="Code">💻</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3AIgogrek" title="Bug reports">🐛</a></td>
    <td align="center"><a href="https://www.jae-sung.com/"><img src="https://avatars.githubusercontent.com/u/39658806?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Jae-Sung Lee</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=jaesung9507" title="Code">💻</a> <a href="#ideas-jaesung9507" title="Ideas, Planning, & Feedback">🤔</a></td>
    <td align="center"><a href="https://github.com/Jarek-SRT"><img src="https://avatars.githubusercontent.com/u/3391365?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Jarek</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=Jarek-SRT" title="Code">💻</a> <a href="#platform-Jarek-SRT" title="Packaging/porting to new platform">📦</a></td>
    <td align="center"><a href="https://github.com/Junkher"><img src="https://avatars.githubusercontent.com/u/85776620?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Junker</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=Junkher" title="Documentation">📖</a></td>
  </tr>
  <tr>
    <td align="center"><a href="https://github.com/kraney"><img src="https://avatars.githubusercontent.com/u/5760081?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Kris Raney</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=kraney" title="Code">💻</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3Akraney" title="Bug reports">🐛</a></td>
    <td align="center"><a href="https://github.com/LukenSkyne"><img src="https://avatars.githubusercontent.com/u/29918069?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Luken</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=LukenSkyne" title="Documentation">📖</a></td>
    <td align="center"><a href="https://markstenglein.com/"><img src="https://avatars.githubusercontent.com/u/9255772?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Mark Stenglein</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=ocelotsloth" title="Code">💻</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3Aocelotsloth" title="Bug reports">🐛</a></td>
    <td align="center"><a href="https://github.com/buddyabaddon"><img src="https://avatars.githubusercontent.com/u/33861511?v=4?s=75" width="75px;" alt=""/><br /><sub><b>buddyabaddon</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=buddyabaddon" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/MikeSchaap"><img src="https://avatars.githubusercontent.com/u/35368821?v=4?s=75" width="75px;" alt=""/><br /><sub><b>MikeSchaap</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=MikeSchaap" title="Code">💻</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3AMikeSchaap" title="Bug reports">🐛</a></td>
    <td align="center"><a href="https://github.com/Orijhins"><img src="https://avatars.githubusercontent.com/u/47521598?v=4?s=75" width="75px;" alt=""/><br /><sub><b>NYSSEN Michaël</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=Orijhins" title="Code">💻</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3AOrijhins" title="Bug reports">🐛</a></td>
    <td align="center"><a href="https://github.com/NanoNik"><img src="https://avatars.githubusercontent.com/u/11991329?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Nan0</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=NanoNik" title="Code">💻</a> <a href="#ideas-NanoNik" title="Ideas, Planning, & Feedback">🤔</a> <a href="https://github.com/wailsapp/wails/commits?author=NanoNik" title="Tests">⚠️</a> <a href="https://github.com/wailsapp/wails/pulls?q=is%3Apr+reviewed-by%3ANanoNik" title="Reviewed Pull Requests">👀</a></td>
    <td align="center"><a href="https://github.com/marcio199226"><img src="https://avatars.githubusercontent.com/u/10244404?v=4?s=75" width="75px;" alt=""/><br /><sub><b>oskar</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=marcio199226" title="Documentation">📖</a></td>
  </tr>
  <tr>
    <td align="center"><a href="https://github.com/pierrejoye"><img src="https://avatars.githubusercontent.com/u/282408?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Pierre Joye</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=pierrejoye" title="Code">💻</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3Apierrejoye" title="Bug reports">🐛</a> <a href="#ideas-pierrejoye" title="Ideas, Planning, & Feedback">🤔</a> <a href="https://github.com/wailsapp/wails/commits?author=pierrejoye" title="Tests">⚠️</a></td>
    <td align="center"><a href="https://github.com/Rested"><img src="https://avatars.githubusercontent.com/u/2003608?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Reuben Thomas-Davis</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=Rested" title="Code">💻</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3ARested" title="Bug reports">🐛</a></td>
    <td align="center"><a href="https://github.com/mewmew"><img src="https://avatars.githubusercontent.com/u/1414531?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Robin</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=mewmew" title="Code">💻</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3Amewmew" title="Bug reports">🐛</a></td>
    <td align="center"><a href="https://threema.id/YSB3TVF7"><img src="https://avatars.githubusercontent.com/u/70367451?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Sebastian Bauer</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=sebastian0x62" title="Code">💻</a> <a href="#ideas-sebastian0x62" title="Ideas, Planning, & Feedback">🤔</a> <a href="https://github.com/wailsapp/wails/commits?author=sebastian0x62" title="Tests">⚠️</a> <a href="https://github.com/wailsapp/wails/pulls?q=is%3Apr+reviewed-by%3Asebastian0x62" title="Reviewed Pull Requests">👀</a> <a href="#question-sebastian0x62" title="Answering Questions">💬</a></td>
    <td align="center"><a href="https://github.com/sidwebworks"><img src="https://avatars.githubusercontent.com/u/58144379?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Sidharth Rathi</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=sidwebworks" title="Documentation">📖</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3Asidwebworks" title="Bug reports">🐛</a></td>
    <td align="center"><a href="https://github.com/sithembiso"><img src="https://avatars.githubusercontent.com/u/6559905?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Sithembiso Khumalo</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=sithembiso" title="Code">💻</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3Asithembiso" title="Bug reports">🐛</a></td>
    <td align="center"><a href="https://github.com/LanguageAgnostic"><img src="https://avatars.githubusercontent.com/u/19310562?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Soheib El-Harrache</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=LanguageAgnostic" title="Code">💻</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3ALanguageAgnostic" title="Bug reports">🐛</a> <a href="#financial-LanguageAgnostic" title="Financial">💵</a></td>
    <td align="center"><a href="https://www.sophieau.com/"><img src="https://avatars.githubusercontent.com/u/11145039?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Sophie Au</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=SophieAu" title="Code">💻</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3ASophieAu" title="Bug reports">🐛</a></td>
  </tr>
  <tr>
    <td align="center"><a href="https://github.com/stefpap"><img src="https://avatars.githubusercontent.com/u/22637722?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Stefanos Papadakis</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=stefpap" title="Code">💻</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3Astefpap" title="Bug reports">🐛</a></td>
    <td align="center"><a href="https://github.com/s12chung"><img src="https://avatars.githubusercontent.com/u/263394?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Steve Chung</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=s12chung" title="Code">💻</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3As12chung" title="Bug reports">🐛</a></td>
    <td align="center"><a href="https://tortloff.de/"><img src="https://avatars.githubusercontent.com/u/41272726?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Timm Ortloff</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=TAINCER" title="Documentation">📖</a></td>
    <td align="center"><a href="https://github.com/tomanagle"><img src="https://avatars.githubusercontent.com/u/8683577?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Tom</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=tomanagle" title="Code">💻</a></td>
    <td align="center"><a href="https://www.linkedin.com/in/valentintrinque"><img src="https://avatars.githubusercontent.com/u/4662842?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Valentin Trinqué</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=ValentinTrinque" title="Code">💻</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3AValentinTrinque" title="Bug reports">🐛</a></td>
    <td align="center"><a href="https://mattn.kaoriya.net/"><img src="https://avatars.githubusercontent.com/u/10111?v=4?s=75" width="75px;" alt=""/><br /><sub><b>mattn</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=mattn" title="Code">💻</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3Amattn" title="Bug reports">🐛</a></td>
    <td align="center"><a href="https://github.com/bearsh"><img src="https://avatars.githubusercontent.com/u/1089356?v=4?s=75" width="75px;" alt=""/><br /><sub><b>bearsh</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=bearsh" title="Code">💻</a> <a href="#ideas-bearsh" title="Ideas, Planning, & Feedback">🤔</a> <a href="https://github.com/wailsapp/wails/commits?author=bearsh" title="Documentation">📖</a></td>
    <td align="center"><a href="https://github.com/chenxiao1990"><img src="https://avatars.githubusercontent.com/u/16933565?v=4?s=75" width="75px;" alt=""/><br /><sub><b>chenxiao</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=chenxiao1990" title="Code">💻</a> <a href="#ideas-chenxiao1990" title="Ideas, Planning, & Feedback">🤔</a> <a href="https://github.com/wailsapp/wails/commits?author=chenxiao1990" title="Documentation">📖</a></td>
  </tr>
  <tr>
    <td align="center"><a href="https://github.com/fengweiqiang"><img src="https://avatars.githubusercontent.com/u/22905300?v=4?s=75" width="75px;" alt=""/><br /><sub><b>fengweiqiang</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=fengweiqiang" title="Code">💻</a> <a href="#platform-fengweiqiang" title="Packaging/porting to new platform">📦</a></td>
    <td align="center"><a href="https://github.com/flin7"><img src="https://avatars.githubusercontent.com/u/58138185?v=4?s=75" width="75px;" alt=""/><br /><sub><b>flin7</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=flin7" title="Documentation">📖</a></td>
    <td align="center"><a href="https://github.com/fred21O4"><img src="https://avatars.githubusercontent.com/u/67189813?v=4?s=75" width="75px;" alt=""/><br /><sub><b>fred21O4</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=fred21O4" title="Documentation">📖</a></td>
    <td align="center"><a href="https://github.com/gardc"><img src="https://avatars.githubusercontent.com/u/41453409?v=4?s=75" width="75px;" alt=""/><br /><sub><b>gardc</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=gardc" title="Documentation">📖</a> <a href="#tutorial-gardc" title="Tutorials">✅</a></td>
    <td align="center"><a href="https://github.com/rayshoo"><img src="https://avatars.githubusercontent.com/u/52561899?v=4?s=75" width="75px;" alt=""/><br /><sub><b>rayshoo</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=rayshoo" title="Documentation">📖</a></td>
    <td align="center"><a href="https://github.com/Yz4230"><img src="https://avatars.githubusercontent.com/u/38999742?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Ishiyama Yuzuki</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=Yz4230" title="Code">💻</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3AYz4230" title="Bug reports">🐛</a></td>
    <td align="center"><a href="https://baiyue.one/"><img src="https://avatars.githubusercontent.com/u/43716063?v=4?s=75" width="75px;" alt=""/><br /><sub><b>佰阅</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=Baiyuetribe" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/daodao97"><img src="https://avatars.githubusercontent.com/u/15009280?v=4?s=75" width="75px;" alt=""/><br /><sub><b>刀刀</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=daodao97" title="Documentation">📖</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3Adaodao97" title="Bug reports">🐛</a></td>
  </tr>
  <tr>
    <td align="center"><a href="https://github.com/jicg"><img src="https://avatars.githubusercontent.com/u/6479672?v=4?s=75" width="75px;" alt=""/><br /><sub><b>归位</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=jicg" title="Code">💻</a> <a href="https://github.com/wailsapp/wails/issues?q=author%3Ajicg" title="Bug reports">🐛</a></td>
    <td align="center"><a href="https://github.com/skamensky"><img src="https://avatars.githubusercontent.com/u/19151369?v=4?s=75" width="75px;" alt=""/><br /><sub><b>skamensky</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=skamensky" title="Code">💻</a> <a href="#ideas-skamensky" title="Ideas, Planning, & Feedback">🤔</a> <a href="https://github.com/wailsapp/wails/commits?author=skamensky" title="Documentation">📖</a></td>
    <td align="center"><a href="https://github.com/apps/dependabot"><img src="https://avatars.githubusercontent.com/in/29110?v=4?s=75" width="75px;" alt=""/><br /><sub><b>dependabot[bot]</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=dependabot[bot]" title="Code">💻</a> <a href="#maintenance-dependabot[bot]" title="Maintenance">🚧</a></td>
    <td align="center"><a href="https://www.linkedin.com/in/dsieradzki/"><img src="https://avatars.githubusercontent.com/u/10297559?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Damian Sieradzki</b></sub></a><br /><a href="#financial-dsieradzki" title="Financial">💵</a></td>
    <td align="center"><a href="https://github.com/boostchicken"><img src="https://avatars.githubusercontent.com/u/427295?v=4?s=75" width="75px;" alt=""/><br /><sub><b>John Dorman</b></sub></a><br /><a href="#financial-boostchicken" title="Financial">💵</a></td>
    <td align="center"><a href="https://blog.iansinnott.com/"><img src="https://avatars.githubusercontent.com/u/3154865?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Ian Sinnott</b></sub></a><br /><a href="#financial-iansinnott" title="Financial">💵</a></td>
    <td align="center"><a href="https://github.com/Shackelford-Arden"><img src="https://avatars.githubusercontent.com/u/7362263?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Arden Shackelford</b></sub></a><br /><a href="#financial-Shackelford-Arden" title="Financial">💵</a></td>
    <td align="center"><a href="https://github.com/Bironou"><img src="https://avatars.githubusercontent.com/u/107761511?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Bironou</b></sub></a><br /><a href="#financial-Bironou" title="Financial">💵</a></td>
  </tr>
  <tr>
    <td align="center"><a href="https://github.com/CharlieGo19"><img src="https://avatars.githubusercontent.com/u/62405980?v=4?s=75" width="75px;" alt=""/><br /><sub><b>CharlieGo_</b></sub></a><br /><a href="#financial-CharlieGo19" title="Financial">💵</a></td>
    <td align="center"><a href="https://github.com/overnet"><img src="https://avatars.githubusercontent.com/u/6376126?v=4?s=75" width="75px;" alt=""/><br /><sub><b>overnet</b></sub></a><br /><a href="#financial-overnet" title="Financial">💵</a></td>
    <td align="center"><a href="https://jugglingjsons.dev/"><img src="https://avatars.githubusercontent.com/u/20739064?v=4?s=75" width="75px;" alt=""/><br /><sub><b>jugglingjsons</b></sub></a><br /><a href="#financial-jugglingjsons" title="Financial">💵</a></td>
    <td align="center"><a href="https://selvin.dev/"><img src="https://avatars.githubusercontent.com/u/1922523?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Selvin Ortiz</b></sub></a><br /><a href="#financial-selvindev" title="Financial">💵</a></td>
    <td align="center"><a href="https://github.com/zandercodes"><img src="https://avatars.githubusercontent.com/u/46308805?v=4?s=75" width="75px;" alt=""/><br /><sub><b>ZanderCodes</b></sub></a><br /><a href="#financial-zandercodes" title="Financial">💵</a></td>
    <td align="center"><a href="https://github.com/DonTomato"><img src="https://avatars.githubusercontent.com/u/1098084?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Michael Voronov</b></sub></a><br /><a href="#financial-DonTomato" title="Financial">💵</a></td>
    <td align="center"><a href="https://lt.hn/"><img src="https://avatars.githubusercontent.com/u/83868036?v=4?s=75" width="75px;" alt=""/><br /><sub><b>letheanVPN</b></sub></a><br /><a href="#financial-letheanVPN" title="Financial">💵</a></td>
    <td align="center"><a href="https://taigrr.com/"><img src="https://avatars.githubusercontent.com/u/8261498?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Tai Groot</b></sub></a><br /><a href="#financial-taigrr" title="Financial">💵</a></td>
  </tr>
  <tr>
    <td align="center"><a href="https://github.com/easy-web-it"><img src="https://avatars.githubusercontent.com/u/95484991?v=4?s=75" width="75px;" alt=""/><br /><sub><b>easy-web-it</b></sub></a><br /><a href="#financial-easy-web-it" title="Financial">💵</a></td>
    <td align="center"><a href="https://michaelolson1996.github.io/portfolio"><img src="https://avatars.githubusercontent.com/u/45323107?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Michael Olson</b></sub></a><br /><a href="#financial-michaelolson1996" title="Financial">💵</a></td>
    <td align="center"><a href="https://eden.network/"><img src="https://avatars.githubusercontent.com/u/4912777?v=4?s=75" width="75px;" alt=""/><br /><sub><b>EdenNetwork Italia</b></sub></a><br /><a href="#financial-EdenNetworkItalia" title="Financial">💵</a></td>
    <td align="center"><a href="https://github.com/ondoki"><img src="https://avatars.githubusercontent.com/u/88536792?v=4?s=75" width="75px;" alt=""/><br /><sub><b>ondoki</b></sub></a><br /><a href="#financial-ondoki" title="Financial">💵</a></td>
    <td align="center"><a href="https://github.com/questrail"><img src="https://avatars.githubusercontent.com/u/3536569?v=4?s=75" width="75px;" alt=""/><br /><sub><b>QuEST Rail LLC</b></sub></a><br /><a href="#financial-questrail" title="Financial">💵</a></td>
    <td align="center"><a href="https://github.com/Gilgames000"><img src="https://avatars.githubusercontent.com/u/22778436?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Gilgameš</b></sub></a><br /><a href="#financial-Gilgames000" title="Financial">💵</a></td>
    <td align="center"><a href="https://github.com/bbergshaven"><img src="https://avatars.githubusercontent.com/u/4091634?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Bernt-Johan Bergshaven</b></sub></a><br /><a href="#financial-bbergshaven" title="Financial">💵</a></td>
    <td align="center"><a href="https://github.com/bglw"><img src="https://avatars.githubusercontent.com/u/40188355?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Liam Bigelow</b></sub></a><br /><a href="#financial-bglw" title="Financial">💵</a></td>
  </tr>
  <tr>
    <td align="center"><a href="https://github.com/nickarellano"><img src="https://avatars.githubusercontent.com/u/13930605?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Nick Arellano</b></sub></a><br /><a href="#financial-nickarellano" title="Financial">💵</a></td>
    <td align="center"><a href="https://github.com/fcjr"><img src="https://avatars.githubusercontent.com/u/2053002?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Frank Chiarulli Jr.</b></sub></a><br /><a href="#financial-fcjr" title="Financial">💵</a></td>
    <td align="center"><a href="https://github.com/tylertravisty"><img src="https://avatars.githubusercontent.com/u/8620352?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Tyler</b></sub></a><br /><a href="#financial-tylertravisty" title="Financial">💵</a></td>
    <td align="center"><a href="https://github.com/trea"><img src="https://avatars.githubusercontent.com/u/1181448?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Trea Hauet</b></sub></a><br /><a href="#financial-trea" title="Financial">💵</a></td>
    <td align="center"><a href="https://picatz.github.io/"><img src="https://avatars.githubusercontent.com/u/14850816?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Kent 'picat' Gruber</b></sub></a><br /><a href="#financial-picatz" title="Financial">💵</a></td>
    <td align="center"><a href="https://github.com/tc-hib"><img src="https://avatars.githubusercontent.com/u/55949036?v=4?s=75" width="75px;" alt=""/><br /><sub><b>tc-hib</b></sub></a><br /><a href="#financial-tc-hib" title="Financial">💵</a></td>
    <td align="center"><a href="https://github.com/acheong08"><img src="https://avatars.githubusercontent.com/u/36258159?v=4?s=75" width="75px;" alt=""/><br /><sub><b>Antonio</b></sub></a><br /><a href="https://github.com/wailsapp/wails/commits?author=acheong08" title="Documentation">📖</a></td>
  </tr>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

<span id="nav-10"></span>

## 特别提及

如果没有以下人员，此项目或许永远不会存在：

- [Dustin Krysak](https://wiki.ubuntu.com/bashfulrobot) - 他的支持和反馈是巨大的。
- [Serge Zaitsev](https://github.com/zserge) - Wails 窗口所使用的 [Webview](https://github.com/zserge/webview) 的作者。
- [Byron](https://github.com/bh90210) - 有时，Byron 一个人保持这个项目活跃着。没有他令人难以置信的投入，我们永远不会得到 v1 。

编写项目代码时伴随着以下专辑：

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

## 特别感谢

<p align="center" style="text-align: center">
   <a href="https://pace.dev"><img src="/assets/images/pace.jpeg"/></a><br/>
   <i>非常</i> 感谢<a href="https://pace.dev">Pace</a>对项目的赞助，并帮助将 Wails 移植到 Apple Silicon !<br/><br/>
   如果您正在寻找一个强大并且快速和易于使用的项目管理工具，可以看看他们！<br/><br/>
</p>

<p align="center" style="text-align: center">
   特别感谢 JetBrains 向我们捐赠许可！<br/><br/>
   请点击 logo 让他们知道你的感激之情！<br/><br/>
   <a href="https://www.jetbrains.com?from=Wails"><img src="/assets/images/jetbrains-grayscale.png" width="30%"></a>
</p>
