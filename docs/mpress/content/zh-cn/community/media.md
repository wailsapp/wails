---
title: "媒体"
description: "Wails 官方按钮，提供预览以及 HTML 和 Markdown 嵌入代码"
slug: "community/media"
sourcePath: "community/media.md"
---

分享使用 Wails 构建的项目时，请使用 Wails 官方资源。下方按钮链接到[wails.io](https://wails.io)，并继续由 Wails 托管，因此无需将图片复制到项目中。

## 使用 Wails 构建

请选择与背景具有足够对比度的版本。每个预览都与访客将看到的链接图片完全一致。

### 浅色背景

<div class="wails-media-preview wails-media-preview-light">  
<p>标准版（<code>180</code> × <code>48</code> 像素）：</p>  
<a href="https://wails.io"><img src="/img/wails-button-light.svg" width="180" height="48" alt="使用 Wails 构建" /></a>  
<p>紧凑版（<code>90</code> × <code>24</code> 像素）：</p>  
<a href="https://wails.io"><img src="/img/wails-button-light.svg" width="90" height="24" alt="使用 Wails 构建" /></a>  
</div>

#### HTML

```html
<a href="https://wails.io">
  <img src="https://wails.io/img/wails-button-light.svg" width="180" height="48" alt="Built with Wails">
</a>
```

紧凑版请使用`width="90" height="24"`。

#### Markdown

```md
[![Built with Wails](https://wails.io/img/wails-button-light.svg)](https://wails.io)
```

### 深色背景

<div class="wails-media-preview wails-media-preview-dark">  
<p>标准版（<code>180</code> × <code>48</code> 像素）：</p>  
<a href="https://wails.io"><img src="/img/wails-button-dark.svg" width="180" height="48" alt="使用 Wails 构建" /></a>  
<p>紧凑版（<code>90</code> × <code>24</code> 像素）：</p>  
<a href="https://wails.io"><img src="/img/wails-button-dark.svg" width="90" height="24" alt="使用 Wails 构建" /></a>  
</div>

#### HTML

```html
<a href="https://wails.io">
  <img src="https://wails.io/img/wails-button-dark.svg" width="180" height="48" alt="Built with Wails">
</a>
```

紧凑版请使用`width="90" height="24"`。

#### Markdown

```md
[![Built with Wails](https://wails.io/img/wails-button-dark.svg)](https://wails.io)
```

Markdown 会按图片固有的`180`×`48`像素尺寸渲染图片。如需设置其他显示尺寸，请使用 HTML 版本。请保持按钮原始的`15:4`比例。
