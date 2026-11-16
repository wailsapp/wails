---
title: "メディア"
description: "公式 Wails ボタンのプレビューと、HTML および Markdown の埋め込みコード"
slug: "community/media"
sourcePath: "community/media.md"
---

Wails で構築したプロジェクトを共有する際は、Wails の公式アセットを使用してください。以下のボタンは [wails.io](https://wails.io) にリンクしており、画像は引き続き Wails によってホストされるため、プロジェクトにコピーする必要はありません。

## Wails で構築

背景に対して十分なコントラストが得られるバージョンを選択してください。各プレビューは、訪問者に表示されるリンク先の画像そのものです。

### 明るい背景

<div class="wails-media-preview wails-media-preview-light">  
<p>標準（<code>180</code> × <code>48</code> px）：</p>  
<a href="https://wails.io"><img src="/img/wails-button-light.svg" width="180" height="48" alt="Wails で構築" /></a>  
<p>コンパクト（<code>90</code> × <code>24</code> px）：</p>  
<a href="https://wails.io"><img src="/img/wails-button-light.svg" width="90" height="24" alt="Wails で構築" /></a>  
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
<p>標準（<code>180</code> × <code>48</code> px）：</p>  
<a href="https://wails.io"><img src="/img/wails-button-dark.svg" width="180" height="48" alt="Wails で構築" /></a>  
<p>コンパクト（<code>90</code> × <code>24</code> px）：</p>  
<a href="https://wails.io"><img src="/img/wails-button-dark.svg" width="90" height="24" alt="Wails で構築" /></a>  
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

Markdown では、画像が本来の `180` × `48` ピクセルのサイズでレンダリングされます。別の表示サイズを設定する必要がある場合は、HTML 版を使用してください。ボタンは元の `15:4` の比率を維持してください。
