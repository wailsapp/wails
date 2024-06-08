<p align="center" style="text-align: center">
  <img src="./assets/images/logo-universal.png" width="55%"><br/>
</p>

<p align="center">
  Créer des applications de bureau avec Go et les technologies Web.
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

## Sommaire

- [Sommaire](#sommaire)
- [Introduction](#introduction)
- [Fonctionnalités](#fonctionnalités)
  - [Feuille de route](#feuille-de-route)
- [Démarrage](#démarrage)
- [Les sponsors](#les-sponsors)
- [Foire aux questions](#foire-aux-questions)
- [Les étoiles au fil du temps](#les-étoiles-au-fil-du-temps)
- [Les contributeurs](#les-contributeurs)
- [License](#license)
- [Inspiration](#inspiration)

## Introduction

La méthode traditionnelle pour fournir des interfaces web aux programmes Go consiste à utiliser un serveur web intégré. Wails propose une approche différente : il offre la possibilité d'intégrer à la fois le code Go et une interface web dans un seul binaire. Des outils sont fournis pour vous faciliter la tâche en gérant la création, la compilation et le regroupement des projets. Il ne vous reste plus qu'à faire preuve de créativité!

## Fonctionnalités

- Utiliser Go pour le backend
- Utilisez n'importe quelle technologie frontend avec laquelle vous êtes déjà familier pour construire votre interface utilisateur.
- Créez rapidement des interfaces riches pour vos programmes Go à l'aide de modèles prédéfinis.
- Appeler facilement des méthodes Go à partir de Javascript
- Définitions Typescript auto-générées pour vos structures et méthodes Go
- Dialogues et menus natifs
- Prise en charge native des modes sombre et clair
- Prise en charge des effets modernes de translucidité et de "frosted window".
- Système d'événements unifié entre Go et Javascript
- Outil puissant pour générer et construire rapidement vos projets
- Multiplateforme
- Utilise des moteurs de rendu natifs - _pas de navigateur intégré_ !

### Feuille de route

La feuille de route du projet peut être consultée [ici](https://github.com/wailsapp/wails/discussions/1484). Veuillez consulter avant d'ouvrir une demande d'amélioration.

## Démarrage

Les instructions d'installation se trouvent sur le site [site officiel](https://wails.io/docs/gettingstarted/installation).

## Les sponsors

Ce projet est soutenu par ces personnes aimables et entreprises:
<img src="website/static/img/sponsors.svg" style="width:100%;max-width:800px;"/>

<p align="center">
<img src="https://wails.io/img/sponsor/jetbrains-grayscale.webp" style="width: 100px"/>
</p>

## Foire aux questions

- S'agit-il d'une alternative à Electron ?

  Cela dépend de vos besoins. Il est conçu pour permettre aux programmeurs Go de créer facilement des applications de bureau légères ou d'ajouter une interface à leurs applications existantes. Wails offre des éléments natifs tels que des menus et des boîtes de dialogue, il peut donc être considéré comme une alternative légère à electron.

- À qui s'adresse ce projet ?

  Les programmeurs Go qui souhaitent intégrer une interface HTML/JS/CSS à leurs applications, sans avoir à créer un serveur et à ouvrir un navigateur pour l'afficher.

- Pourquoi ce nom ??

  Lorsque j'ai vu WebView, je me suis dit : "Ce que je veux vraiment, c'est un outil pour construire une application WebView, un peu comme Rails l'est pour Ruby". Au départ, il s'agissait donc d'un jeu de mots (Webview on Rails). Il se trouve que c'est aussi un homophone du nom anglais du [Pays](https://en.wikipedia.org/wiki/Wales) d'où je viens. Il s'est donc imposé.

## Les étoiles au fil du temps

[![Graphique de l'histoire des étoiles](https://api.star-history.com/svg?repos=wailsapp/wails&type=Date)](https://star-history.com/#wailsapp/wails&Date)

## Les contributeurs

La liste des contributeurs devient trop importante pour le readme ! Toutes les personnes extraordinaires qui ont contribué à ce projet ont leur propre page [ici](https://wails.io/credits#contributors).

## License

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fwailsapp%2Fwails.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fwailsapp%2Fwails?ref=badge_large)

## Inspiration

Ce projet a été principalement codé sur les albums suivants :

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
