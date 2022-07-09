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
