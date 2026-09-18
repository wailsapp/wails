---
title: "Media"
description: "Tombol resmi Wails dengan pratinjau serta kode sematan HTML dan Markdown"
slug: "community/media"
sourcePath: "community/media.md"
---

Gunakan aset resmi Wails saat membagikan proyek yang dibuat dengan Wails. Tombol di bawah mengarah ke [wails.io](https://wails.io) dan gambarnya tetap dihosting oleh Wails, sehingga Anda tidak perlu menyalin gambar ke proyek Anda.

## Dibuat dengan Wails

Pilih versi dengan kontras yang memadai terhadap latar belakangnya. Setiap pratinjau menampilkan gambar bertautan yang persis sama dengan yang akan dilihat pengunjung Anda.

### Latar belakang terang

<div class="wails-media-preview wails-media-preview-light">  
<p>Ukuran standar (<code>180</code> × <code>48</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-light.svg" width="180" height="48" alt="Dibuat dengan Wails" /></a>  
<p>Ukuran ringkas (<code>90</code> × <code>24</code> px):</p>  
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
<p>Ukuran standar (<code>180</code> × <code>48</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-dark.svg" width="180" height="48" alt="Dibuat dengan Wails" /></a>  
<p>Ukuran ringkas (<code>90</code> × <code>24</code> px):</p>  
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

Markdown menampilkan gambar pada ukuran aslinya, yaitu `180` kali `48` piksel. Gunakan versi HTML jika Anda perlu menentukan ukuran tampilan yang berbeda. Pertahankan rasio aspek asli tombol, yaitu `15:4`.
