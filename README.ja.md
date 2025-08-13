<h1 align="center">Wails</h1>

<p align="center" style="text-align: center">
  <img src="./assets/images/logo-universal.png" width="55%"><br/>
</p>

<p align="center">
  GoとWebの技術を用いてデスクトップアプリケーションを構築します。
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
[한국어](README.ko.md) · [Español](README.es.md) · [Português](README.pt-br.md) ·
[Русский](README.ru.md) · [Francais](README.fr.md) · [Uzbek](README.uz.md) · [Deutsch](README.de.md) ·
[Türkçe](README.tr.md)

</samp>
</strong>
</div>

## 目次

- [目次](#目次)
- [はじめに](#はじめに)
- [特徴](#特徴)
  - [ロードマップ](#ロードマップ)
- [始め方](#始め方)
- [スポンサー](#スポンサー)
- [FAQ](#faq)
- [スター数の推移](#スター数の推移)
- [コントリビューター](#コントリビューター)
- [ライセンス](#ライセンス)
- [インスピレーション](#インスピレーション)


## はじめに

Go プログラムにウェブインタフェースを提供する従来の方法は内蔵のウェブサーバを経由するものですが、 Wails では異なるアプローチを提供します。
Wails では Go のコードとウェブフロントエンドを単一のバイナリにまとめる機能を提供します。
また、プロジェクトの作成、コンパイル、ビルドを行うためのツールが提供されています。あなたがすべきことは創造性を発揮することです！

## 特徴

- バックエンドには Go を利用しています
- 使い慣れたフロントエンド技術を利用して UI を構築できます
- あらかじめ用意されたテンプレートを利用することで、リッチなフロントエンドを備えた Go プログラムを素早く作成できます
- JavaScript から Go のメソッドを簡単に呼び出すことができます
- あなたの書いた Go の構造体やメソットに応じた TypeScript の定義が自動生成されます
- ネイティブのダイアログとメニューが利用できます
- ネイティブなダーク/ライトモードをサポートします
- モダンな半透明や「frosted window」エフェクトをサポートしています
- Go と JavaScript 間で統一されたイベント・システムを備えています
- プロジェクトを素早く生成して構築する強力な cli ツールを用意しています
- マルチプラットフォームに対応しています
- ネイティブなレンダリングエンジンを使用しています - _つまりブラウザを埋め込んでいるわけではありません！_

### ロードマップ

プロジェクトのロードマップは[こちら](https://github.com/wailsapp/wails/discussions/1484)になります。  
機能拡張のリクエストを出す前にご覧ください。

## 始め方

インストール方法は[公式サイト](https://wails.io/docs/gettingstarted/installation)に掲載されています。

## スポンサー

このプロジェクトは、以下の方々・企業によって支えられています。
<img src="website/static/img/sponsors.svg" style="width:100%;max-width:800px;"/>

## FAQ

- Electron の代替品になりますか？

  それはあなたの求める要件によります。Wails は Go プログラマーが簡単に軽量のデスクトップアプリケーションを作成したり、既存のアプリケーションにフロントエンドを追加できるように設計されています。
  Wails v2 ではメニューやダイアログといったネイティブな要素を提供するようになったため、軽量な Electron の代替となりつつあります。

- このプロジェクトは誰に向けたものですか？

  HTML/JS/CSS のフロントエンド技術をアプリケーションにバンドルさせることで、サーバーを作成してブラウザ経由で表示させることなくアプリケーションを利用したい Go プログラマにおすすめです。

- 名前の由来を教えて下さい

  WebView を見たとき、私はこう思いました。  
  「私が本当に欲しいのは、WebView アプリを構築するためのツールであり、Ruby に対する Rails のようなものである」と。  
  そのため、最初は言葉遊びのつもりでした（Webview on Rails）。  
  また、私の[出身国](https://en.wikipedia.org/wiki/Wales)の英語名と同音異義語でもあります。そしてこの名前が定着しました。

## スター数の推移

[![Star History Chart](https://api.star-history.com/svg?repos=wailsapp/wails&type=Date)](https://star-history.com/#wailsapp/wails&Date)

## コントリビューター

貢献してくれた方のリストが大きくなりすぎて、readme に入りきらなくなりました！  
このプロジェクトに貢献してくれた素晴らしい方々のページは[こちら](https://wails.io/credits#contributors)です。

## ライセンス

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fwailsapp%2Fwails.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fwailsapp%2Fwails?ref=badge_large)

## インスピレーション

プロジェクトを進める際に、以下のアルバムたちも支えてくれています。

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

