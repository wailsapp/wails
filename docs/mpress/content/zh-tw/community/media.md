---
title: "媒體"
description: "Wails 官方按鈕，附預覽以及 HTML 和 Markdown 嵌入碼"
slug: "community/media"
sourcePath: "community/media.md"
---

分享使用 Wails 建置的專案時，請使用 Wails 官方素材。下方按鈕會連結至[wails.io](https://wails.io)，且圖片仍由 Wails 託管，因此不必將圖片複製到專案中。

## 使用 Wails 建置

請選擇與背景有足夠對比的版本。每個預覽都是訪客實際會看到的連結圖片。

### 淺色背景

<div class="wails-media-preview wails-media-preview-light">  
<p>標準（<code>180</code> × <code>48</code> 像素）：</p>  
<a href="https://wails.io"><img src="/img/wails-button-light.svg" width="180" height="48" alt="使用 Wails 建置" /></a>  
<p>精簡（<code>90</code> × <code>24</code> 像素）：</p>  
<a href="https://wails.io"><img src="/img/wails-button-light.svg" width="90" height="24" alt="使用 Wails 建置" /></a>  
</div>

#### HTML

```html
<a href="https://wails.io">
  <img src="https://wails.io/img/wails-button-light.svg" width="180" height="48" alt="Built with Wails">
</a>
```

精簡版請使用`width="90" height="24"`。

#### Markdown

```md
[![Built with Wails](https://wails.io/img/wails-button-light.svg)](https://wails.io)
```

### 深色背景

<div class="wails-media-preview wails-media-preview-dark">  
<p>標準（<code>180</code> × <code>48</code> 像素）：</p>  
<a href="https://wails.io"><img src="/img/wails-button-dark.svg" width="180" height="48" alt="使用 Wails 建置" /></a>  
<p>精簡（<code>90</code> × <code>24</code> 像素）：</p>  
<a href="https://wails.io"><img src="/img/wails-button-dark.svg" width="90" height="24" alt="使用 Wails 建置" /></a>  
</div>

#### HTML

```html
<a href="https://wails.io">
  <img src="https://wails.io/img/wails-button-dark.svg" width="180" height="48" alt="Built with Wails">
</a>
```

精簡版請使用`width="90" height="24"`。

#### Markdown

```md
[![Built with Wails](https://wails.io/img/wails-button-dark.svg)](https://wails.io)
```

Markdown 會以圖片原有的`180`乘以`48`像素尺寸呈現圖片。需要設定不同的顯示尺寸時，請使用 HTML 版本。請維持按鈕原有的`15:4`比例。
