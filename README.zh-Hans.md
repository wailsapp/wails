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
  <a href="https://github.com/wailsapp/wails/issues">
    <img src="https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat" alt="CodeFactor" />
  </a>
  <a href="https://app.fossa.com/projects/git%2Bgithub.com%2Fwailsapp%2Fwails?ref=badge_shield" alt="FOSSA Status">
    <img src="https://app.fossa.com/api/projects/git%2Bgithub.com%2Fwailsapp%2Fwails.svg?type=shield" />
  </a>
  <a href="https://github.com/avelino/awesome-go" rel="nofollow">
    <img src="https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg" alt="Awesome" />
  </a>
  <a href="https://discord.gg/BrRSWTaxVK">
    <img alt="Discord" src="https://dcbadge.vercel.app/api/server/BrRSWTaxVK?style=flat"/>
  </a>
  <br/>
  <a href="https://github.com/wailsapp/wails/actions/workflows/build-and-test.yml" rel="nofollow">
    <img src="https://img.shields.io/github/actions/workflow/status/wailsapp/wails/build-and-test.yml?branch=master&logo=Github" alt="Build" />
  </a>
  <a href="https://github.com/wailsapp/wails/tags" rel="nofollow">
    <img alt="GitHub tag (latest SemVer pre-release)" src="https://img.shields.io/github/v/tag/wailsapp/wails?include_prereleases&label=version"/>
  </a>
</p>

<div align="center">
<strong>
<samp>

[English](README.md) · [简体中文](README.zh-Hans.md) · [日本語](README.ja.md) ·
[한국어](README.ko.md)

</samp>
</strong>
</div>

## 内容目录

- [内容目录](#内容目录)
- [项目介绍](#项目介绍)
- [功能](#功能)
  - [路线图](#路线图)
- [快速入门](#快速入门)
- [赞助商](#赞助商)
- [常见问题](#常见问题)
- [星星增长趋势](#星星增长趋势)
- [贡献者](#贡献者)
- [许可证](#许可证)
- [灵感](#灵感)

## 项目介绍

为 Go 程序提供 Web 界面的传统方法是通过内置 Web 服务器。Wails 提供了一种不同的方
法：它提供了将 Go 代码和 Web 前端一起打包成单个二进制文件的能力。通过提供的工具
，可以很轻松的完成项目的创建、编译和打包。你所要做的就是发挥创造力！

## 功能

- 后端使用标准 Go
- 使用您已经熟悉的任何前端技术来构建您的 UI
- 使用内置模板为您的 Go 程序快速创建丰富的前端
- 从 Javascript 轻松调用 Go 方法
- 为您的 Go 结构体和方法自动生成 Typescript 声明
- 原生对话框和菜单
- 支持现代半透明和“磨砂窗”效果
- Go 和 Javascript 之间统一的事件系统
- 强大的命令行工具，可快速生成和构建您的项目
- 跨平台
- 使用原生渲染引擎 - _没有嵌入浏览器_！

### 路线图

项目路线图可在 [此处](https://github.com/wailsapp/wails/discussions/1484) 找到。
在提出增强请求之前请查阅此内容。

## 快速入门

使用说明在 [官网](https://wails.io/docs/gettingstarted/installation)。

## 赞助商

这个项目由以下这些人或者公司支持：

<img src="website/static/img/sponsors.svg" style="width:100%;max-width:800px;"/>

## 常见问题

- 它是 Electron 的替代品吗?

  取决于您的要求。它旨在使 Go 程序员可以轻松制作轻量级桌面应用程序或在其现有应用
  程序中添加前端。尽管 Wails 当前不提供对诸如菜单之类的原生元素的钩子，但将来可
  能会改变。

- 这个项目针对的是哪些人?

  希望将 HTML / JS / CSS 前端与其应用程序捆绑在一起的程序员，而不是借助创建服务
  并打开浏览器进行查看的方式。

- 名字怎么来的?

  当我看到 WebView 时，我想"我真正想要的是围绕构建 WebView 应用程序工作，有点像
  Rails 对于 Ruby"。因此，最初它是一个文字游戏（Webview on Rails）。碰巧也是我来
  自的 [国家](https://en.wikipedia.org/wiki/Wales) 的英文名字的同音。所以就是它
  了。

## 星星增长趋势

[![星星增长趋势](https://api.star-history.com/svg?repos=wailsapp/wails&type=Date)](https://star-history.com/#wailsapp/wails&Date)

## 贡献者

贡献者列表对于 README 文件来说太大了！所有为这个项目做出贡献的了不起的人
在[这里](https://wails.io/credits#contributors)都有自己的页面。

## 许可证

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fwailsapp%2Fwails.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fwailsapp%2Fwails?ref=badge_large)

## 灵感

项目灵感主要来自以下专辑：

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
