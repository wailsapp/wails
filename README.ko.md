<h1 align="center">Wails</h1>

<p align="center" style="text-align: center">
  <img src="./assets/images/logo-universal.png" width="55%"><br/>
</p>

<p align="center">
  Go & Web 기술을 사용하여 데스크탑 애플리케이션을 빌드하세요.
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

## 목차

- [목차](#목차)
- [소개](#소개)
- [기능](#기능)
  - [로드맵](#로드맵)
- [시작하기](#시작하기)
- [스폰서](#스폰서)
- [FAQ](#faq)
- [Stargazers 성장 추세](#stargazers-성장-추세)
- [기여자](#기여자)
- [라이센스](#라이센스)
- [영감](#영감)

## 소개

Go 프로그램에 웹 인터페이스를 제공하는 전통적인 방법은 내장 웹 서버를 이용하는
것입니다. Wails는 다르게 접근합니다: Go 코드와 웹 프론트엔드를 단일 바이너리로
래핑하는 기능을 제공합니다. 프로젝트 생성, 컴파일 및 번들링을 처리하여 이를 쉽게
수행할 수 있도록 도구가 제공됩니다. 창의력을 발휘하기만 하면 됩니다!

## 기능

- 백엔드에 표준 Go 사용
- 이미 익숙한 프론트엔드 기술을 사용하여 UI 구축
- 사전 구축된 템플릿을 사용하여 Go 프로그램을 위한 풍부한 프론트엔드를 빠르게 생
  성
- Javascript에서 Go 메서드를 쉽게 호출
- Go 구조체 및 메서드에 대한 자동 생성된 Typescript 정의
- 기본 대화 및 메뉴
- 네이티브 다크/라이트 모드 지원
- 최신 반투명도 및 "반투명 창" 효과 지원
- Go와 Javascript 간의 통합 이벤트 시스템
- 프로젝트를 빠르게 생성하고 구축하는 강력한 CLI 도구
- 멀티플랫폼
- 기본 렌더링 엔진 사용 - _내장 브라우저 없음_!

### 로드맵

프로젝트 로드맵은 [여기](https://github.com/wailsapp/wails/discussions/1484)에서
확인할 수 있습니다. 개선 요청을 하기 전에 이것을 참조하십시오.

## 시작하기

설치 지침은 [공식 웹사이트](https://wails.io/docs/gettingstarted/installation)에
있습니다.

## 스폰서

이 프로젝트는 친절한 사람들 / 회사들이 지원합니다.
<img src="website/static/img/sponsors.svg" style="width:100%;max-width:800px;"/>

## FAQ

- 이것은 Electron의 대안인가요?

  요구 사항에 따라 다릅니다. Go 프로그래머가 쉽게 가벼운 데스크톱 애플리케이션을
  만들거나 기존 애플리케이션에 프론트엔드를 추가할 수 있도록 설계되었습니다.
  Wails는 메뉴 및 대화 상자와 같은 기본 요소를 제공하므로 가벼운 Electron 대안으
  로 간주될 수 있습니다.

- 이 프로젝트는 누구를 대상으로 하나요?

  서버를 생성하고 이를 보기 위해 브라우저를 열 필요 없이 HTML/JS/CSS 프런트엔드
  를 애플리케이션과 함께 묶고자 하는 프로그래머를 대상으로 합니다.

- Wails 이름의 의미는 무엇인가요?

  WebView를 보았을 때 저는 "내가 정말로 원하는 것은 WebView 앱을 구축하기 위한
  도구를 사용하는거야. 마치 Ruby on Rails 처럼 말이야."라고 생각했습니다. 그래서
  처음에는 말장난(Webview on Rails)이었습니다.
  [국가](https://en.wikipedia.org/wiki/Wales)에 대한 영어 이름의 동음이의어이기
  도 하여 정했습니다.

## Stargazers 성장 추세

[![Star History Chart](https://api.star-history.com/svg?repos=wailsapp/wails&type=Date)](https://star-history.com/#wailsapp/wails&Date)

## 기여자

기여자 목록이 추가 정보에 비해 너무 커지고 있습니다! 이 프로젝트에 기여한 모든
놀라운 사람들은 [여기](https://wails.io/credits#contributors)에 자신의 페이지를
가지고 있습니다.

## 라이센스

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fwailsapp%2Fwails.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fwailsapp%2Fwails?ref=badge_large)

## 영감

이 프로젝트는 주로 다음 앨범을 들으며 코딩되었습니다.

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
