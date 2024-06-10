<p align="center" style="text-align: center">
  <img src="./assets/images/logo-universal.png" width="55%"><br/>
</p>

<p align="center">
  Go ve Web Teknolojilerini kullanarak masaüstü uygulamaları oluşturun.
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
[Русский](README.ru.md) · [Francais](README.fr.md) · [Uzbek](README.uz.md) ·
[Türkçe](README.tr.md)

</samp>
</strong>
</div>

## İçerik

- [İçerik](#içerik)
- [Giriş](#giriş)
- [Özellikler](#özellikler)
  - [Yol Haritası](#yol-haritası)
- [Başlarken](#başlarken)
- [Sponsorlar](#sponsorlar)
- [Sıkça sorulan sorular](#sıkça-sorulan-sorular)
- [Zaman içinda yıldızlayanlar](#zaman-içinde-yıldızlayanlar)
- [Katkıda bulunanlar](#katkıda-bulunanlar)
- [Lisans](#lisans)
- [İlham](#ilham)

## Giriş

Go programlarına web arayüzleri sağlamak için geleneksel yöntem, yerleşik bir web sunucusu kullanmaktır. Wails, farklı bir yaklaşım sunar: Hem Go kodunu hem de bir web ön yüzünü tek bir ikili dosyada paketleme yeteneği sağlar. Proje oluşturma, derleme ve paketleme işlemlerini kolaylaştıran araçlar sunar. Tek yapmanız gereken yaratıcı olmaktır!

## Özellikler

- Backend için standart Go kullanın
- Kullanıcı arayüzünüzü oluşturmak için zaten aşina olduğunuz herhangi bir frontend teknolojisini kullanın
- Hazır şablonlar kullanarak Go programlarınız için hızlıca zengin ön yüzler oluşturun
- Javascript'ten Go metodlarını kolayca çağırın
- Go yapı ve metodlarınız için otomatik oluşturulan Typescript tanımları
- Yerel Diyaloglar ve Menüler
- Yerel Karanlık / Aydınlık mod desteği
- Modern saydamlık ve "buzlu cam" efektlerini destekler
- Go ve Javascript arasında birleşik olay sistemi
- Projelerinizi hızlıca oluşturmak ve derlemek için güçlü bir komut satırı aracı
- Çoklu platform desteği
- Yerel render motorlarını kullanır - _gömülü tarayıcı yok_!


### Yol Haritesı

Proje yol haritasına [buradan](https://github.com/wailsapp/wails/discussions/1484) ulaşabilirsiniz. Lütfen bir iyileştirme talebi oluşturmadan önce danışın.


## Başlarken

Kurulum talimatları [resmi web sitesinde](https://wails.io/docs/gettingstarted/installation) bulunmaktadır.


## Sponsorlar

Bu proje, aşağıdaki nazik insanlar / şirketler tarafından desteklenmektedir:
<img src="website/static/img/sponsors.svg" style="width:100%;max-width:800px;"/>

<p align="center">
<img src="https://wails.io/img/sponsor/jetbrains-grayscale.webp" style="width: 100px"/>
</p>

## Sıkça Sorulan Sorular

- Bu Electron'a alternatif mi?

  Gereksinimlerinize bağlıdır. Go programcılarının hafif masaüstü uygulamaları yapmasını veya mevcut uygulamalarına bir ön yüz eklemelerini kolaylaştırmak için tasarlanmıştır. Wails, menüler ve diyaloglar gibi yerel öğeler sunduğundan, hafif bir Electron alternatifi olarak kabul edilebilir.

- Bu proje kimlere yöneliktir?

  HTML/JS/CSS ön yüzünü uygulamalarıyla birlikte paketlemek isteyen, ancak bir sunucu oluşturup bir tarayıcı açmaya başvurmadan bunu yapmak isteyen Go programcıları için.

- İsmin anlamı nedir?

  WebView'i gördüğümde, "Aslında istediğim şey, WebView uygulaması oluşturmak için araçlar, biraz Rails'in Ruby için olduğu gibi" diye düşündüm. Bu nedenle başlangıçta kelime oyunu (Rails üzerinde Webview) olarak ortaya çıktı. Ayrıca, benim geldiğim [ülkenin](https://en.wikipedia.org/wiki/Wales) İngilizce adıyla homofon olması tesadüf oldu. Bu yüzden bu isim kaldı.


## Zaman içinda yıldızlayanlar

<a href="https://star-history.com/#wailsapp/wails&Date">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=wailsapp/wails&type=Date&theme=dark" />
    <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=wailsapp/wails&type=Date" />
    <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=wailsapp/wails&type=Date" />
  </picture>
</a>

## Katkıda Bulunanlar

Katkıda bulunanların listesi, README için çok büyük hale geldi! Bu projeye katkıda bulunan tüm harika insanların kendi sayfaları [burada](https://wails.io/credits#contributors) bulunmaktadır.


## Lisans

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fwailsapp%2Fwails.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fwailsapp%2Fwails?ref=badge_large)

## İlham

Bu proje esas olarak aşağıdaki albümler dinlenilerek kodlandı:

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

