---
title: "Media"
description: "Tombol resmi Wails dengan pratinjau serta sematan HTML dan Markdown"
slug: "community/media"
sourcePath: "community/media.md"
---

Gunakan aset resmi Wails saat membagikan proyek yang dibuat dengan Wails. Tombol di bawah ini tertaut ke [wails.io](https://wails.io) dan tetap dihosting oleh Wails, sehingga Anda tidak perlu menyalin gambar tersebut ke dalam proyek Anda.

## Dibuat dengan Wails

Pilih versi yang memiliki kontras memadai dengan latar belakangnya. Setiap pratinjau adalah gambar tertaut persis yang akan dilihat pengunjung Anda.

### Latar belakang terang

<div class="wails-media-preview wails-media-preview-light">  
<p>Standar (<code>180</code> × <code>48</code> piksel):</p>  
<a href="https://wails.io"><img src="/img/wails-button-light.svg" width="180" height="48" alt="Dibuat dengan Wails" /></a>  
<p>Ringkas (<code>90</code> × <code>24</code> piksel):</p>  
<a href="https://wails.io"><img src="/img/wails-button-light.svg" width="90" height="24" alt="Dibuat dengan Wails" /></a>  
</div>

#### HTML

```html
<a href="https://wails.io">
  <img src="https://wails.io/img/wails-button-light.svg" width="180" height="48" alt="Built with Wails">
</a>
```

Untuk versi ringkas, gunakan `width="90" height="24"`.

#### Markdown

```md
[![Built with Wails](https://wails.io/img/wails-button-light.svg)](https://wails.io)
```

### Latar belakang gelap

<div class="wails-media-preview wails-media-preview-dark">  
<p>Standar (<code>180</code> × <code>48</code> piksel):</p>  
<a href="https://wails.io"><img src="/img/wails-button-dark.svg" width="180" height="48" alt="Dibuat dengan Wails" /></a>  
<p>Ringkas (<code>90</code> × <code>24</code> piksel):</p>  
<a href="https://wails.io"><img src="/img/wails-button-dark.svg" width="90" height="24" alt="Dibuat dengan Wails" /></a>  
</div>

#### HTML

```html
<a href="https://wails.io">
  <img src="https://wails.io/img/wails-button-dark.svg" width="180" height="48" alt="Built with Wails">
</a>
```

Untuk versi ringkas, gunakan `width="90" height="24"`.

#### Markdown

```md
[![Built with Wails](https://wails.io/img/wails-button-dark.svg)](https://wails.io)
```

Markdown merender gambar sesuai ukuran intrinsiknya, yaitu `180` × `48` piksel. Gunakan versi HTML jika Anda perlu menetapkan ukuran tampilan yang berbeda. Pertahankan proporsi asli tombol, yaitu `15:4`.
