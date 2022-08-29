<p align="center" style="text-align: center">
  <img src="./assets/images/logo-universal.png" width="55%"><br/>
</p>

<p align="center">
  GoとWebの技術を用いてデスクトップアプリケーションを構築します。<br/><br/>
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
注意：v2のリリースが近づいているため、v1の新しい機能リクエストやバグレポートは受け付けておりません。重要な問題がある場合はチケットを開き、なぜそれが重要なのかを明記してください。
</strong>
</h3>
<hr/>

## 国際化

[English](README.md) | [简体中文](README.zh-Hans.md) | [日本語](README.ja.md)

## 目次

<details>
  <summary>クリックすることで、ディレクトリ一覧の開閉が可能です。</summary>

- [1. 国際化](#国際化)
- [2. 目次](#目次)
- [3. はじめに](#はじめに)
  - [3.1 公式サイト](#公式サイト)
  - [3.2 ロードマップ](#ロードマップ)
- [4. 特徴](#特徴)
- [5. スポンサー](#スポンサー)
- [6. 始め方](#始め方)
- [7. FAQ](#faq)
- [8. コントリビューター](#コントリビューター)
- [9. 特記事項](#特記事項)
- [10. スペシャルサンクス](#スペシャルサンクス)
- [11. ライセンス](#ライセンス)

</details>

## はじめに

Goプログラムにウェブインタフェースを提供する従来の方法は内蔵のウェブサーバを経由するものですが、 Wailsでは異なるアプローチを提供します。
Wailsでは Go のコードとウェブフロントエンドを単一のバイナリにまとめる機能を提供します。
また、プロジェクトの作成、コンパイル、ビルドを行うためのツールが提供されています。あなたがすべきことは創造性を発揮することです！

### 公式サイト

Version 2:

Wails v2が3つのプラットフォームでベータ版としてリリースされました。興味のある方は[新しいウェブサイト](https://wails.io)をご覧ください。

レガシー版 v1:

レガシー版 v1のドキュメントは[https://wails.app](https://wails.app)で見ることができます。

### ロードマップ

プロジェクトのロードマップは[こちら](https://github.com/wailsapp/wails/discussions/1484)になります。  
機能拡張のリクエストを出す前にご覧ください。

## 特徴

- バックエンドにはGoを利用しています
- 使い慣れたフロントエンド技術を利用してUIを構築できます
- あらかじめ用意されたテンプレートを利用することで、リッチなフロントエンドを備えたGoプログラムを作成できます
- JavaScriptからGoのメソッドを簡単に呼び出すことができます
- あなたの書いたGoの構造体やメソットに応じたTypeScriptの定義が自動生成されます
- ネイティブのダイアログとメニューが利用できます
- モダンな半透明や「frosted window」エフェクトをサポートしています
- GoとJavaScript間で統一されたイベント・システムを備えています
- プロジェクトを素早く生成して構築する強力なcliツールを用意しています
- マルチプラットフォームに対応しています
- ネイティブなレンダリングエンジンを使用しています - _つまりブラウザを埋め込んでいるわけではありません！_


## スポンサー

このプロジェクトは、以下の方々・企業によって支えられています。

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

## 始め方

インストール方法は[公式サイト](https://wails.io/docs/gettingstarted/installation)に掲載されています。

## FAQ

- Electronの代替品になりますか？

  それはあなたの求める要件によります。WailsはGoプログラマーが簡単に軽量のデスクトップアプリケーションを作成したり、既存のアプリケーションにフロントエンドを追加できるように設計されています。
  Wails v2ではメニューやダイアログといったネイティブな要素を提供するようになったため、軽量なElectronの代替となりつつあります。

- このプロジェクトは誰に向けたものですか？

  HTML/JS/CSSのフロントエンド技術をアプリケーションにバンドルさせることで、サーバーを作成してブラウザ経由で表示させることなくアプリケーションを利用したいGoプログラマにおすすめです。

- 名前の由来を教えて下さい

  WebViewを見たとき、私はこう思いました。  
  「私が本当に欲しいのは、WebViewアプリを構築するためのツールであり、Rubyに対するRailsのようなものである」と。  
  そのため、最初は言葉遊びのつもりでした（Webview on Rails）。  
  また、私の[出身国](https://en.wikipedia.org/wiki/Wales)の英語名と同音異義語でもあります。そしてこの名前が定着しました。

## スター数の推移

[![スター数の推移](https://starchart.cc/wailsapp/wails.svg)](https://starchart.cc/wailsapp/wails)

## コントリビューター

貢献してくれた方のリストが大きくなりすぎて、readmeに入りきらなくなりました！  
このプロジェクトに貢献してくれた素晴らしい方々のページは[こちら](https://wails.io/credits#contributors)です。

## 特記事項

このプロジェクトは以下の方々の協力がなければ、実現しなかったと思います。

- [Dustin Krysak](https://wiki.ubuntu.com/bashfulrobot) - 彼のサポートとフィードバックはとても大きいものでした。
- [Serge Zaitsev](https://github.com/zserge) - Wailsのウィンドウで使用している[Webview](https://github.com/zserge/webview)の作者です。
- [Byron](https://github.com/bh90210) - 時にはByronが一人でこのプロジェクトを存続させてくれたこともありました。彼の素晴らしいインプットがなければv1に到達することはなかったでしょう。

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

## スペシャルサンクス

<p align="center" style="text-align: center">
   <a href="https://pace.dev"><img src="/assets/images/pace.jpeg"/></a><br/>
   このプロジェクトを後援し、WailsをApple Siliconに移植する取り組みを支援してくれた <a href="https://pace.dev">Pace</a> に <i>とても</i>感謝しています！<br/><br/>
   パワフルで素早く簡単に使えるプロジェクト管理ツールをお探しなら、ぜひチェックしてみてください！<br/><br/>
</p>

<p align="center" style="text-align: center">
   ライセンスを提供していただいたJetBrains社に感謝します！<br/><br/>
   ロゴをクリックして、感謝の気持ちを伝えてください！<br/><br/>
   <a href="https://www.jetbrains.com?from=Wails"><img src="/assets/images/jetbrains-grayscale.png" width="30%"></a>
</p>

## ライセンス

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fwailsapp%2Fwails.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fwailsapp%2Fwails?ref=badge_large)
