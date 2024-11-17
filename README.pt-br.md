<p align="center" style="text-align: center">
  <img src="./assets/images/logo-universal.png" width="55%"><br/>
</p>

<p align="center">
  Crie aplicativos de desktop usando Go e tecnologias Web.
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
[한국어](README.ko.md) · [Español](README.es.md) · [Português](README.pt-br.md) · [Francais](README.fr.md) · [Uzbek](README.uz.md) · [Deutsch](README.de.md) ·
[Türkçe](README.tr.md)

</samp>
</strong>
</div>

## Índice

- [Índice](#índice)
- [Introdução](#introdução)
- [Recursos e funcionalidades](#recursos-e-funcionalidades)
  - [Plano de trabalho](#plano-de-trabalho)
- [Iniciando](#iniciando)
- [Patrocinadores](#patrocinadores)
- [Perguntas frequentes](#perguntas-frequentes)
- [Estrelas ao longo do tempo](#estrelas-ao-longo-do-tempo)
- [Colaboradores](#colaboradores)
- [Licença](#licença)
- [Inspiração](#inspiração)

## Introdução

O método tradicional de fornecer interfaces da Web para programas Go é por meio de um servidor da Web integrado. Wails oferece uma
abordagem: fornece a capacidade de agrupar o código Go e um front-end da Web em um único binário. As ferramentas são fornecidas para
que torne isso mais fácil para você lidando com a criação, compilação e agrupamento de projetos. Tudo o que você precisa fazer é ser criativo!

## Recursos e funcionalidades

- Use Go padrão para o back-end
- Use qualquer tecnologia de front-end com a qual você já esteja familiarizado para criar sua interface do usuário
- Crie rapidamente um front-end avançado para seus programas Go usando modelos pré-construídos
- Chame facilmente métodos Go com JavaScript
- Definições TypeScript geradas automaticamente para suas estruturas e métodos Go
- Diálogos e menus nativos
- Suporte nativo ao modo escuro/claro
- Suporta translucidez moderna e efeitos de "janela fosca"
- Sistema de eventos unificado entre Go e JavaScript
- Poderosa ferramenta cli para gerar e construir rapidamente seus projetos
- Multiplataforma
- Usa mecanismos de renderização nativos - _sem navegador incorporado_!

### Plano de trabalho

O plano de trabalho do projeto pode ser encontrado [aqui](https://github.com/wailsapp/wails/discussions/1484). Por favor consulte
isso antes de abrir um pedido de melhoria.

## Iniciando

As instruções de instalação estão no [site oficial](https://wails.io/docs/gettingstarted/installation).

## Patrocinadores

Este projeto é apoiado por estas simpáticas pessoas/empresas:
<img src="website/static/img/sponsors.svg" style="width:100%;max-width:800px;"/>

<p align="center">
<img src="https://wails.io/img/sponsor/jetbrains-grayscale.webp" style="width: 100px"/>
</p>

## Perguntas frequentes

- Esta é uma alternativa ao Electron?

  Depende de seus requisitos. Ele foi projetado para tornar mais fácil para os programadores Go criar aplicações desktop
  e adicionar um front-end aos seus aplicativos existentes. O Wails oferece elementos nativos, como menus
  e diálogos, por isso pode ser considerada uma alternativa leve, se comparado ao Electron.

- A quem se destina este projeto?

  Programadores Go que desejam agrupar um front-end HTML/JS/CSS com seus aplicativos, sem recorrer à criação de um
  servidor e abrir um navegador para visualizá-lo.

- Qual é o significado do nome?

  Quando vi o WebView, pensei "O que eu realmente quero é ferramentas para construir um aplicativo WebView, algo semelhante ao que Rails é para Ruby". Portanto, inicialmente era um jogo de palavras (WebView on Rails). Por acaso, também era um homófono do
  Nome em inglês para o [país](https://en.wikipedia.org/wiki/Wales) de onde eu sou. Então ficou com esse nome.

## Estrelas ao longo do tempo

[![Star History Chart](https://api.star-history.com/svg?repos=wailsapp/wails&type=Date)](https://star-history.com/#wailsapp/wails&Date)

## Colaboradores

A lista de colaboradores está ficando grande demais para o arquivo readme! Todas as pessoas incríveis que contribuíram para o
projeto tem sua própria página [aqui](https://wails.io/credits#contributors).

## Licença

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fwailsapp%2Fwails.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fwailsapp%2Fwails?ref=badge_large)

## Inspiração

Este projeto foi construído ouvindo esses álbuns:

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
