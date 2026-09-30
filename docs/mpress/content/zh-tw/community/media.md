---
title: "媒體素材"
description: "Wails 官方按鈕預覽及 HTML 與 Markdown 嵌入程式碼"
slug: "community/media"
sourcePath: "community/media.md"
---

分享使用 Wails 建立的專案時，請使用 Wails 官方素材。下方按鈕連結至 [wails.io](https://wails.io)，圖片仍由 Wails 代管，因此不必將圖片複製到專案中。

## 使用 Wails 建立

請選擇與背景有足夠對比度的版本。每個預覽都使用訪客將看到的同一張附有連結的圖片。

### 淺色背景

<div class="wails-media-preview wails-media-preview-light">  
<p>標準尺寸 (<code>180</code> × <code>48</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-light.svg" width="180" height="48" alt="使用 Wails 建立" /></a>  
<p>小尺寸 (<code>90</code> × <code>24</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-light.svg" width="90" height="24" alt="使用 Wails 建立" /></a>  
</div>

#### HTML

```html
<a href="https://wails.io">
  <img src="https://wails.io/img/wails-button-light.svg" width="180" height="48" alt="Built with Wails">
</a>
```

小尺寸版本請使用 `width="90" height="24"`。

#### Markdown

```md
[![Built with Wails](https://wails.io/img/wails-button-light.svg)](https://wails.io)
```

### 深色背景

<div class="wails-media-preview wails-media-preview-dark">  
<p>標準尺寸 (<code>180</code> × <code>48</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-dark.svg" width="180" height="48" alt="使用 Wails 建立" /></a>  
<p>小尺寸 (<code>90</code> × <code>24</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-dark.svg" width="90" height="24" alt="使用 Wails 建立" /></a>  
</div>

#### HTML

```html
<a href="https://wails.io">
  <img src="https://wails.io/img/wails-button-dark.svg" width="180" height="48" alt="Built with Wails">
</a>
```

小尺寸版本請使用 `width="90" height="24"`。

#### Markdown

```md
[![Built with Wails](https://wails.io/img/wails-button-dark.svg)](https://wails.io)
```

Markdown 會以圖片的原始尺寸顯示，也就是寬 `180`、高 `48` 像素。如需設定其他顯示尺寸，請使用 HTML 版本。請維持按鈕原有的 `15:4` 寬高比。
