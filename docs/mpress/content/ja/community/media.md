---
title: "メディア"
description: "Wails公式ボタンのプレビューとHTML・Markdownの埋め込みコード"
slug: "community/media"
sourcePath: "community/media.md"
---

Wailsで作成したプロジェクトを紹介するときは、Wailsの公式素材をご利用ください。以下のボタンは [wails.io](https://wails.io) にリンクしており、画像はWails側でホストされているため、プロジェクトにコピーする必要はありません。

## Wailsで作成

背景に対して十分なコントラストが得られる種類を選んでください。各プレビューには、訪問者に表示されるものと同じリンク付き画像を使用しています。

### 明るい背景

<div class="wails-media-preview wails-media-preview-light">  
<p>標準サイズ (<code>180</code> × <code>48</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-light.svg" width="180" height="48" alt="Wailsで作成" /></a>  
<p>コンパクトサイズ (<code>90</code> × <code>24</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-light.svg" width="90" height="24" alt="Wailsで作成" /></a>  
</div>

#### HTML

```html
<a href="https://wails.io">
  <img src="https://wails.io/img/wails-button-light.svg" width="180" height="48" alt="Built with Wails">
</a>
```

コンパクト版には `width="90" height="24"` を使用してください。

#### Markdown

```md
[![Built with Wails](https://wails.io/img/wails-button-light.svg)](https://wails.io)
```

### 暗い背景

<div class="wails-media-preview wails-media-preview-dark">  
<p>標準サイズ (<code>180</code> × <code>48</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-dark.svg" width="180" height="48" alt="Wailsで作成" /></a>  
<p>コンパクトサイズ (<code>90</code> × <code>24</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-dark.svg" width="90" height="24" alt="Wailsで作成" /></a>  
</div>

#### HTML

```html
<a href="https://wails.io">
  <img src="https://wails.io/img/wails-button-dark.svg" width="180" height="48" alt="Built with Wails">
</a>
```

コンパクト版には `width="90" height="24"` を使用してください。

#### Markdown

```md
[![Built with Wails](https://wails.io/img/wails-button-dark.svg)](https://wails.io)
```

Markdownでは、画像は元のサイズである横 `180` × 縦 `48` ピクセルで表示されます。表示サイズを変更したい場合は、HTML版を使用してください。ボタンの縦横比は元の `15:4` を維持してください。
