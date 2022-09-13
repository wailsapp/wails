<h1 align="center">Wails</h1>

<p align="center" style="text-align: center">
  <img src="./assets/images/logo-universal.png" width="55%"><br/>
</p>

<p align="center">
  使用 Go 和 Web 技术构建桌面应用程序。
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
  <br/>
  <a href="https://github.com/wailsapp/wails/actions/workflows/build.yml" rel="nofollow">
    <img src="https://img.shields.io/github/workflow/status/wailsapp/wails/Build?logo=github" alt="Build" />
  </a>
  <a href="https://github.com/misitebao/wails/tags" rel="nofollow">
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

<hr/>
<h3 align="center">
<strong>
请注意：随着我们接近 v2 版本，我们不接受 v1 的任何新功能请求或错误报告。如果您有一个关键问题，请开一个Issue并说明为什么它很关键。
</strong>
</h3>
<hr/>

## 内容目录

<details>
  <summary>点我 打开/关闭 目录列表</summary>

- [内容目录](#内容目录)
- [项目介绍](#项目介绍)
  - [官方网站](#官方网站)
  - [路线图](#路线图)
- [功能](#功能)
- [赞助商](#赞助商)
- [快速入门](#快速入门)
- [常见问题](#常见问题)
- [贡献者](#贡献者)
- [特别提及](#特别提及)
- [特别感谢](#特别感谢)
- [许可证](#许可证)

</details>

## 项目介绍

为 Go 程序提供 Web 界面的传统方法是通过内置 Web 服务器。Wails 提供了一种不同的方法：它提供了将 Go 代码和 Web
前端一起打包成单个二进制文件的能力。通过提供的工具，可以很轻松的完成项目的创建、编译和打包。你所要做的就是发挥创造力！

### 官方网站

V2：

Wails v2 已针对所有 3 个平台发布了 Beta 版。如果您有兴趣尝试一下，请查看[新网站](https://wails.io)。

旧版 V1：

旧版 v1 文档可以在[https://wails.app](https://wails.app)找到。

### 路线图

项目路线图可在[此处](https://github.com/wailsapp/wails/discussions/1484)找到。在提出增强请求之前请查阅此内容。

## 功能

- 后端使用标准 Go
- 使用您已经熟悉的任何前端技术来构建您的 UI
- 使用内置模板为您的 Go 程序快速创建丰富的前端
- 从 Javascript 轻松调用 Go 方法
- 为您的 Go 结构体和方法自动生成 Typescript 声明
- 原生对话框和菜单
- 支持现代半透明和“磨砂窗”效果
- Go 和 Javascript 之间的统一事件系统
- 强大的 CLI 工具，可快速生成和构建您的项目
- 跨平台
- 使用原生渲染引擎 - _没有嵌入式浏览器_！

## 赞助商

这个项目由以下这些人或者公司支持：

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
<a href="https://github.com/KiddoV" style="width:45px">
  <img src="https://github.com/KiddoV.png?size=45" width="45"/>
</a>

## 快速入门

使用说明在[官网](https://wails.io/docs/gettingstarted/installation)。

## 常见问题

- 它是 Electron 的替代品吗?

  取决于您的要求。它旨在使 Go 程序员可以轻松制作轻量级桌面应用程序或在其现有应用程序中添加前端。尽管 Wails 当前不提供对诸如菜单之类的原生元素的钩子，但将来可能会改变。

- 这个项目针对的是哪些人?

  希望将 HTML / JS / CSS 前端与其应用程序捆绑在一起的程序员，而不是借助创建服务并打开浏览器进行查看的方式。

- 名字怎么来的?

  当我看到 WebView 时，我想"我真正想要的是围绕构建 WebView 应用程序工作，有点像 Rails 对于 Ruby"。因此，最初它是一个文字游戏（Webview on
  Rails）。碰巧也是我来自的 [国家](https://en.wikipedia.org/wiki/Wales) 的英文名字的同音。所以就是它了。

## 星星增长趋势

[![星星增长趋势](https://starchart.cc/wailsapp/wails.svg)](https://starchart.cc/wailsapp/wails)

## 贡献者

贡献者列表对于 README 文件来说太大了！所有为这个项目做出贡献的了不起的人在[这里](https://wails.io/credits#contributors)都有自己的页面。

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

## 许可证

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fwailsapp%2Fwails.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fwailsapp%2Fwails?ref=badge_large)
