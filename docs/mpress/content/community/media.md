---
title: "Media"
description: "Official Wails buttons with previews and HTML and Markdown embeds"
slug: "community/media"
sourcePath: "community/media.md"
---

Use the official Wails assets when sharing a project built with Wails. The buttons below link to [wails.io](https://wails.io) and remain hosted by Wails, so you do not need to copy the image into your project.

## Built with Wails

Choose the version with sufficient contrast for its background. Each preview is the exact linked image that your visitors will see.

### Light background

<div class="wails-media-preview wails-media-preview-light">  
<p>Standard (<code>180</code> × <code>48</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-light.svg" width="180" height="48" alt="Built with Wails" /></a>  
<p>Compact (<code>90</code> × <code>24</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-light.svg" width="90" height="24" alt="Built with Wails" /></a>  
</div>

#### HTML

```html
<a href="https://wails.io">
  <img src="https://wails.io/img/wails-button-light.svg" width="180" height="48" alt="Built with Wails">
</a>
```

For the compact version, use `width="90" height="24"`.

#### Markdown

```md
[![Built with Wails](https://wails.io/img/wails-button-light.svg)](https://wails.io)
```

### Dark background

<div class="wails-media-preview wails-media-preview-dark">  
<p>Standard (<code>180</code> × <code>48</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-dark.svg" width="180" height="48" alt="Built with Wails" /></a>  
<p>Compact (<code>90</code> × <code>24</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-dark.svg" width="90" height="24" alt="Built with Wails" /></a>  
</div>

#### HTML

```html
<a href="https://wails.io">
  <img src="https://wails.io/img/wails-button-dark.svg" width="180" height="48" alt="Built with Wails">
</a>
```

For the compact version, use `width="90" height="24"`.

#### Markdown

```md
[![Built with Wails](https://wails.io/img/wails-button-dark.svg)](https://wails.io)
```

Markdown renders the image at its intrinsic `180` by `48` pixel size. Use the HTML version when you need to set a different display size. Keep the buttons at their original `15:4` proportions.
