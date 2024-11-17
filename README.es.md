<p align="center" style="text-align: center">
  <img src="./assets/images/logo-universal.png" width="55%"><br/>
</p>

<p align="center">
  Construye aplicaciones de escritorio usando Go y tecnologías web.
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

## Tabla de Contenidos

- [Tabla de Contenidos](#tabla-de-contenidos)
- [Introducción](#introducción)
- [Funcionalidades](#funcionalidades)
  - [Plan de Trabajo](#plan-de-trabajo)
- [Empezando](#empezando)
- [Patrocinadores](#patrocinadores)
- [Preguntas Frecuentes](#preguntas-frecuentes)
- [Estrellas a lo Largo del Tiempo](#estrellas-a-lo-largo-del-tiempo)
- [Colaboradores](#colaboradores)
- [Licencia](#licencia)
- [Inspiración](#inspiración)

## Introducción

El método tradicional para proveer una interfaz web en programas hechos con Go
es a través del servidor web incorporado. Wails ofrece un enfoque diferente al
permitir combinar el código hecho en Go con un frontend web en un solo archivo
binario. Las herramientas que proporcionamos facilitan este trabajo para ti, al
crear, compilar y empaquetar tu proyecto. ¡Lo único que debes hacer es ponerte
creativo!

## Funcionalidades

- Utiliza Go estándar para el backend
- Utiliza cualquier tecnología frontend con la que ya estés familiarizado para
  construir tu interfaz de usuario
- Crea rápidamente interfaces de usuario enriquecidas para tus programas en Go
  utilizando plantillas predefinidas
- Invoca fácilmente métodos de Go desde Javascript
- Definiciones de Typescript generadas automáticamente para tus structs y
  métodos de Go
- Diálogos y menús nativos
- Soporte nativo de modo oscuro / claro
- Soporte de translucidez y efectos de ventana esmerilada
- Sistema de eventos unificado entre Go y Javascript
- Herramienta CLI potente para generar y construir tus proyectos rápidamente
- Multiplataforma
- Usa motores de renderizado nativos - ¡_sin navegador integrado_!

### Plan de Trabajo

El plan de trabajo se puede encontrar
[aqui](https://github.com/wailsapp/wails/discussions/1484). Por favor,
consúltalo antes de abrir una solicitud de mejora.

## Empezando

Las instrucciones de instalacion se encuentran en nuestra
[pagina web oficial](https://wails.io/docs/gettingstarted/installation).

## Patrocinadores

Este Proyecto cuenta con el apoyo de estas amables personas/ compañías:
<img src="website/static/img/sponsors.svg" style="width:100%;max-width:800px;"/>

<p align="center">
<img src="https://wails.io/img/sponsor/jetbrains-grayscale.webp" style="width: 100px"/>
</p>

## Preguntas Frecuentes

- ¿Es esta una alternativa a Electron?

  Depende de tus requisitos. Está diseñado para facilitar a los programadores de
  Go la creación de aplicaciones de escritorio livianas o agregar una interfaz
  gráfica a sus aplicaciones existentes. Wails ofrece elementos nativos como
  menús y diálogos, por lo que podría considerarse una alternativa liviana a
  Electron.

- ¿A quien esta dirigido este proyecto?

  El proyecto esta dirigido a programadores de Go que desean integrar una
  interfaz HMTL/JS/CSS en sus aplicaciones, sin tener que recurrir a la creación
  de un servidor y abrir el navegador para visualizarla.

- ¿Cual es el significado del nombre?

  Cuando vi WebView, pensé: "Lo que realmente quiero es una herramienta para
  construir una aplicación WebView, algo similar a lo que Rails es para Ruby".
  Así que inicialmente fue un juego de palabras (WebView en Rails). Además, por
  casualidad, también es homófono del nombre en inglés del
  [país](https://en.wikipedia.org/wiki/Wales) del que provengo. Así que se quedó
  con ese nombre.

## Estrellas a lo Largo del Tiempo

[![Star History Chart](https://api.star-history.com/svg?repos=wailsapp/wails&type=Date)](https://star-history.com/#wailsapp/wails&Date)

## Colaboradores

¡La lista de colaboradores se está volviendo demasiado grande para el archivo
readme! Todas las personas increíbles que han contribuido a este proyecto tienen
su propia página [aqui](https://wails.io/credits#contributors).

## Licencia

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fwailsapp%2Fwails.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fwailsapp%2Fwails?ref=badge_large)

## Inspiración

Este proyecto fue construido mientras se escuchaban estos álbumes:

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
  [Alun Tan Lan - Y Distawrwydd](https://open.spotify.com/album/0c32OywcLpdJCWWMC6vB8v)
