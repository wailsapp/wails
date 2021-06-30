<p align="center" style="text-align: center">
   <img src="logo_cropped.png" width="40%"><br/>
</p>
<p align="center">
   使用 Go 和 Web 技术构建桌面应用程序。<br/><br/>
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

<span id="nav-1"></span>

## 国际化

[English](README.md) | 简体中文

向 Go 程序提供 Web 接口的传统方法是通过内置 Web 服务器。Wails 提供了一种不同的方法：它提供了将 Go 代码和 Web 前端都包装成单个二进制文件的能力。通过处理项目创建、编译和打包，可为您提供工具，使您轻松做到这一点。你所要做的就是发挥创造力！

官方文档可以在 [https://wails.app](https://wails.app) 中找到。
国内镜像站点 [https://wails.top](https://wails.top)。

<span id="nav-2"></span>

## 内容目录

- [1. 国际化](#nav-1)
- [2. 内容目录](#nav-2)
- [3. 特征](#nav-3)
- [4. 赞助商](#nav-4)
- [5. 安装](#nav-5)
  - [5.1 MacOS](#nav-5-1)
  - [5.2 Linux](#nav-5-2)
    - [5.2.1 Debian/Ubuntu](#nav-5-2-1)
    - [5.2.2 Arch Linux / ArchLabs / Ctlos Linux](#nav-5-2-2)
    - [5.2.3 Centos](#nav-5-2-3)
    - [5.2.4 Fedora](#nav-5-2-4)
    - [5.2.5 VoidLinux & VoidLinux-musl](#nav-5-2-5)
    - [5.2.6 Gentoo](#nav-5-2-6)
  - [5.3 Windows](#nav-5-3)
- [6. 安装](#nav-6)
- [7. 下一步](#nav-7)
- [8. 常见问题](#nav-8)
- [9. 贡献者](#nav-9)
- [10. 特别提及](#nav-10)
- [11. 许可协议](#nav-11)
- [12. 特别感谢](#nav-12)

<span id="nav-3"></span>

## 特征

- 后端使用标准 Go
- 使用任意前端技术构建 UI 界面
- 快速为您的 Go 应用生成 Vue、Vuetify、React 前端代码
- 通过简单的绑定命令将 Go 方法暴露到前端
- 使用原生渲染引擎 - 无嵌入式浏览器
- 共享事件系统
- 原生文件系统对话框
- 强大的命令行工具
- 跨多个平台

<span id="nav-4"></span>

## 赞助商

这个项目由以下这些人或者公司支持：

<a href="https://github.com/matryer" style="width:100px"><img src="https://github.com/matryer.png" width="100"/></a>
<a href="https://www.jetbrains.com?from=Wails" style="width:100px"><img src="jetbrains-grayscale.png" width="100"/></a>
<a href="https://github.com/tc-hib" style="width:55px;border-radius: 50%">
<img src="https://github.com/tc-hib.png?size=55" width="55" style="border-radius: 50%"/>
</a>
<a href="https://github.com/picatz" style="width:50px;border-radius: 50%">
<img src="https://github.com/picatz.png?size=50" width="50" style="border-radius: 50%"/>
</a>
<a href="https://github.com/tylertravisty" style="width:50px;border-radius: 50%">
<img src="https://github.com/tylertravisty.png?size=50" width="50" style="border-radius: 50%"/>
</a>
<a href="https://github.com/akhudek" style="width:50px;border-radius: 50%">
<img src="https://github.com/akhudek.png?size=50" width="50" style="border-radius: 50%"/>
</a>
<a href="https://github.com/akhudek" style="width:50px;border-radius: 50%">
<img src="https://github.com/akhudek.png?size=50" width="50" style="border-radius: 50%"/>
</a>
<a href="https://github.com/trea" style="width:50px;border-radius: 50%">
<img src="https://github.com/trea.png?size=50" width="50" style="border-radius: 50%"/>
</a>
<a href="https://github.com/LanguageAgnostic" style="width:55px;border-radius: 50%">
<img src="https://github.com/LanguageAgnostic.png?size=55" width="55" style="border-radius: 50%"/>
</a>

<span id="nav-5"></span>

## 安装

Wails 使用 cgo 与原生渲染引擎结合，因此需要一些依赖平台的库以及 Go 的安装。基本要求是：

- Go 1.16
- npm

<span id="nav-5-1"></span>

### MacOS

请确保已安装 xcode 命令行工具。这可以通过运行下面的命令来完成：

`xcode-select --install`

<span id="nav-5-2"></span>

### Linux

<span id="nav-5-2-1"></span>

#### Debian/Ubuntu

`sudo apt install libgtk-3-dev libwebkit2gtk-4.0-dev`

_Debian: 8, 9, 10_

_Ubuntu: 16.04, 18.04, 19.04_

_Also succesfully tested on: Zorin 15, Parrot 4.7, Linuxmint 19, Elementary 5, Kali, Neon_, Pop!\_OS

<span id="nav-5-2-2"></span>

#### Arch Linux / ArchLabs / Ctlos Linux

`sudo pacman -S webkit2gtk gtk3`

_Also succesfully test on: Manjaro & ArcoLinux_

<span id="nav-5-2-3"></span>

#### Centos

`sudo yum install webkitgtk3-devel gtk3-devel`

_CentOS 6, 7_

<span id="nav-5-2-4"></span>

#### Fedora

`sudo yum install webkit2gtk3-devel gtk3-devel`

_Fedora 29, 30_

<span id="nav-5-2-5"></span>

#### VoidLinux & VoidLinux-musl

`xbps-install gtk+3-devel webkit2gtk-devel`

<span id="nav-5-2-6"></span>

#### Gentoo

`sudo emerge gtk+:3 webkit-gtk`

<span id="nav-5-3"></span>

### Windows

Windows 需要 GCC 和相关工具。 建议从 [http://tdm-gcc.tdragon.net/download](http://tdm-gcc.tdragon.net/download) 下载， 安装完成，您就可以开始了。

<span id="nav-6"></span>

## 安装

**确保 Go modules 是开启的: GO111MODULE=on 并且 go/bin 在您的 PATH 变量中.**

安装很简单，运行以下命令：

<pre style='color:white'>
go get -u github.com/wailsapp/wails/cmd/wails
</pre>

<span id="nav-7"></span>

## 下一步

建议在此时阅读 [https://wails.app](https://wails.app) 上面的文档.

<span id="nav-8"></span>

## 常见问题

- 它是 Electron 的替代品吗?

  取决于您的要求。它旨在使 Go 程序员可以轻松制作轻量级桌面应用程序或在其现有应用程序中添加前端。尽管 Wails 当前不提供对诸如菜单之类的原生元素的钩子，但将来可能会改变。

- 这个项目针对的是谁?

  希望将 HTML / JS / CSS 前端与其应用程序捆绑在一起的程序员，而无需借助创建服务并打开浏览器进行查看的方式。

- 名字怎么来的?

  当我看到 WebView 时，我想"我真正想要的是围绕构建 WebView 应用程序工作，有点像 Rails 对于 Ruby"。因此，最初它是一个文字游戏（Webview on Rails）。碰巧也是我来自的 [国家](https://en.wikipedia.org/wiki/Wales) 的英文名字的同音。所以就是他了。

<span id="nav-9"></span>

## 贡献者

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
<a href="https://github.com/akhudek"><img src="https://github.com/akhudek.png?size=40" width="40"/></a></a>

<span id="nav-10"></span>

## 特别提及

如果没有以下人员，此项目将永远不会存在：

- [Dustin Krysak](https://wiki.ubuntu.com/bashfulrobot) - His support and feedback has been immense. More patience than you can throw a stick at (Not long now Dustin!).
- [Serge Zaitsev](https://github.com/zserge) - Creator of [Webview](https://github.com/zserge/webview) which Wails uses for the windowing.
- [Byron](https://github.com/bh90210) - At times, Byron has single handedly kept this project alive. Without his incredible input, we never would have got to v1.

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

## 许可协议

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fwailsapp%2Fwails.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fwailsapp%2Fwails?ref=badge_large)

<span id="nav-12"></span>

## 特别感谢

<p align="center" style="text-align: center">
   <a href="https://pace.dev"><img src="pace.jpeg"/></a><br/>
   非常感谢<a href="https://pace.dev">Pace</a>对项目的赞助，并帮助将Wails移植到Apple Silicon<br/><br/>
   如果您正在寻找一个强大的项目管理工具，并且快速和易于使用，可以看看他们！<br/><br/>
</p>

<p align="center" style="text-align: center">
   特别感谢JetBrains向我们捐赠许可！<br/><br/>
   请点击logo让他们知道你的感激之情！<br/><br/>
   <a href="https://www.jetbrains.com?from=Wails"><img src="jetbrains-grayscale.png" width="30%"></a>
</p>
