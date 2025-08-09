<p align="center" style="text-align: center">
  <img src="./assets/images/logo-universal.png" width="55%"><br/>
</p>

<p align="center">
  Собирайте Desktop приложения используя Go и Web технологии
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
[한국어](README.ko.md) · [Español](README.es.md) · [Русский](README.ru.md) · [Francais](README.fr.md) · [Uzbek](README.uz.md) · [Deutsch](README.de.md) ·
[Türkçe](README.tr.md)

</samp>
</strong>
</div>

## Содержание

- [Содержание](#содержание)
- [Вступление](#вступление)
- [Особенности](#особенности)
  - [Roadmap](#roadmap)
- [Быстрый старт](#быстрый-старт)
- [Спонсоры](#спонсоры)
- [FAQ](#faq)
- [График звёздочек](#график-звёздочек-репозитория-относительно-времени)
- [Контребьюторы](#контребьюторы)
- [Лицензия](#лицензия)
- [Вдохновение](#вдохновение)

## Вступление

Обычно, веб-интерфейсы для программ Go - это встроенный веб-сервер и веб-браузер.
У Walls другой подход: он оборачивает как код Go, так и веб-интерфейс в один бинарник (EXE файл).
Облегчает вам создание вашего приложения, управляя созданием, компиляцией и объединением проектов.
Все ограничивается лишь вашей фантазией!

## Особенности

- Использование Go для backend
- Поддержка любой frontend технологии, с которой вы уже знакомы для создания вашего UI
- Быстрое создание frontend для ваших программ, используя готовые шаблоны
- Очень лёгкий вызов функций Go из JavaScript
- Автогенерация TypeScript типов для Go структур и функций
- Нативные диалоги и меню
- Нативная поддержка тёмной и светлой темы
- Поддержка современных эффектов прозрачности и "матового окна"
- Единая система эвентов для Go и JavaScript
- Мощный CLI для быстрого создания ваших проектов
- Мультиплатформенность
- Использование нативного движка рендеринга - нет встроенному браузеру!

### Roadmap

Roadmap проекта вы можете найти [здесь](https://github.com/wailsapp/wails/discussions/1484).
Пожалуйста, проконсультируйтесь перед предложением улучшения.

## Быстрый старт

Инструкции по установке находятся на [официальном сайте](https://wails.io/docs/gettingstarted/installation).

## Спонсоры

Проект поддерживается этими добрыми людьми / компаниями:
<img src="website/static/img/sponsors.svg" style="width:100%;max-width:800px;"/>

<p align="center">
<img src="https://wails.io/img/sponsor/jetbrains-grayscale.webp" style="width: 100px"/>
</p>

## FAQ

- Это альтернатива Electron?

  Зависит от ваших требований. Wails разработан для легкого создания Desktop приложений или
  расширения интерфейсной части существующих приложений для программистов на Go. Wails действительно
  предлагает встроенные элементы, такие как меню и диалоги, так что его можно считать облегченной альтернативой Electron.

- Для кого предназначен этот проект?

  Для Golang программистов, которые хотят создавать приложения, используя HTML, JS и CSS,
  без создания веб-сервера и открытия браузера для их просмотра.

- Что это за название?

  Когда я увидел WebView, я подумал: "Что мне действительно нужно, так это инструменты для создания приложения WebView,
  немного похожие на Rails для Ruby". Изначально это была игра слов (Webview on Rails). Просто так получилось, что это
  также омофон английского названия для [Страны](https://en.wikipedia.org/wiki/Wales) от куда я родом. Так что это прижилось.

## График звёздочек репозитория по времени

[![График звёзд](https://api.star-history.com/svg?repos=wailsapp/wails&type=Date)](https://star-history.com/#wailsapp/wails&Date)

## Контрибьюторы

Список участников слишком велик для README! У всех замечательных людей, которые внесли свой вклад в этот
проект, есть своя [страничка](https://wails.io/credits#contributors).

## Лицензия

[![Статус FOSSA](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fwailsapp%2Fwails.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fwailsapp%2Fwails?ref=badge_large)

## Вдохновение

Этот проект был создан, в основном, под эти альбомы:

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
