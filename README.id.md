<h1 align="center">Wails</h1>

<p align="center" style="text-align: center">
  <img src="./assets/images/logo-universal.png" width="55%"><br/>
</p>

<p align="center">
  Buat aplikasi desktop menggunakan Go & Teknologi Web.
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
    <img alt="Discord" src="https://img.shields.io/discord/1042734330029547630?logo=discord"/>
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
[Türkçe](README.tr.md) · [Bahasa Indonesia](README.id.md)

</samp>
</strong>
</div>

## Daftar Isi

- [Daftar Isi](#daftar-isi)
- [Pendahuluan](#pendahuluan)
- [Fitur](#fitur)
  - [Peta Jalan](#peta-jalan)
- [Memulai](#memulai)
- [Sponsor](#sponsor)
- [FAQ](#faq)
- [Stargazers dari Waktu ke Waktu](#stargazers-dari-waktu-ke-waktu)
- [Kontributor](#kontributor)
- [Lisensi](#lisensi)
- [Inspirasi](#inspirasi)

## Pendahuluan

Metode tradisional untuk menyediakan antarmuka web pada program Go adalah melalui server web bawaan. Wails menawarkan pendekatan yang berbeda: kemampuan untuk membungkus kode Go dan frontend web ke dalam satu binary. Alat-alat disediakan agar proses ini mudah bagi Anda dengan menangani pembuatan proyek, kompilasi, dan bundling. Yang perlu Anda lakukan hanya berkreasi!

## Fitur

- Gunakan Go standar untuk backend
- Gunakan teknologi frontend apa pun yang sudah Anda kenal untuk membangun UI
- Buat frontend yang kaya dengan cepat untuk program Go Anda menggunakan template bawaan
- Panggil metode Go dari Javascript dengan mudah
- Definisi Typescript yang dihasilkan otomatis untuk struct dan metode Go Anda
- Dialog & Menu Native
- Dukungan mode Gelap / Terang native
- Mendukung efek translucency modern dan "frosted window"
- Sistem eventing terpadu antara Go dan Javascript
- Alat CLI yang andal untuk membuat dan membangun proyek dengan cepat
- Multiplatform
- Menggunakan mesin rendering native - _tanpa browser embedded_!

### Peta Jalan

Peta jalan proyek dapat ditemukan [di sini](https://github.com/wailsapp/wails/discussions/1484). Harap
konsultasikan sebelum membuat permintaan peningkatan.

## Memulai

Wails memiliki dua versi aktif:

| Versi | Status | Instalasi | Dokumentasi |
|---|---|---|---|
| v2 | Stabil | `go install github.com/wailsapp/wails/v2/cmd/wails@latest` | [wails.io](https://wails.io/) |
| v3 | Alpha | `go install github.com/wailsapp/wails/v3/cmd/wails3@latest` | [v3.wails.io](https://v3.wails.io/) |

Petunjuk instalasi lengkap tersedia untuk [v2](https://wails.io/docs/gettingstarted/installation) dan [v3](https://v3.wails.io).

## Sponsor

Proyek ini didukung oleh orang-orang / perusahaan baik hati ini:
<img src="website/static/img/sponsors.svg" style="width:100%;max-width:800px;"/>

## Didukung Oleh

[![JetBrains logo.](https://resources.jetbrains.com/storage/products/company/brand/logos/jetbrains.svg)](https://jb.gg/OpenSource)

## FAQ

- Apakah ini alternatif Electron?

  Tergantung kebutuhan Anda. Wails dirancang agar programmer Go dapat dengan mudah membuat aplikasi desktop
  ringan atau menambahkan frontend ke aplikasi yang sudah ada. Wails menawarkan elemen native seperti menu
  dan dialog, sehingga dapat dianggap sebagai alternatif Electron yang ringan.

- Proyek ini ditujukan untuk siapa?

  Programmer Go yang ingin mengemas frontend HTML/JS/CSS dengan aplikasi mereka, tanpa perlu membuat
  server dan membuka browser untuk melihatnya.

- Apa arti namanya?

  Ketika saya melihat WebView, saya berpikir "Yang benar-benar saya inginkan adalah alat untuk membangun aplikasi WebView, seperti Rails untuk
  Ruby". Jadi awalnya ini permainan kata (Webview on Rails). Kebetulan juga homofon dengan
  nama bahasa Inggris untuk [Negara](https://en.wikipedia.org/wiki/Wales) asal saya. Jadi namanya tetap dipakai.

## Stargazers dari Waktu ke Waktu

<a href="https://star-history.com/#wailsapp/wails&Date">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=wailsapp/wails&type=Date&theme=dark" />
    <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=wailsapp/wails&type=Date" />
    <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=wailsapp/wails&type=Date" />
  </picture>
</a>

## Kontributor

Daftar kontributor sudah terlalu panjang untuk readme! Semua orang luar biasa yang telah berkontribusi pada
proyek ini memiliki halaman mereka sendiri [di sini](https://wails.io/credits#contributors).

## Lisensi

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fwailsapp%2Fwails.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fwailsapp%2Fwails?ref=badge_large)

## Inspirasi

Proyek ini sebagian besar dikodekan sambil mendengarkan album-album berikut:

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
