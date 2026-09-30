---
title: "미디어"
description: "Wails 공식 버튼 미리보기와 HTML 및 Markdown 삽입 코드"
slug: "community/media"
sourcePath: "community/media.md"
---

Wails로 만든 프로젝트를 공유할 때 Wails 공식 자료를 사용하세요. 아래 버튼은 [wails.io](https://wails.io)로 연결되며 이미지는 Wails에서 계속 호스팅하므로 프로젝트에 이미지를 복사할 필요가 없습니다.

## Wails로 제작

배경과 충분한 대비를 이루는 버전을 선택하세요. 각 미리보기는 방문자에게 표시될 것과 동일한 링크가 있는 이미지를 보여 줍니다.

### 밝은 배경

<div class="wails-media-preview wails-media-preview-light">  
<p>표준 크기 (<code>180</code> × <code>48</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-light.svg" width="180" height="48" alt="Wails로 제작" /></a>  
<p>작은 크기 (<code>90</code> × <code>24</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-light.svg" width="90" height="24" alt="Wails로 제작" /></a>  
</div>

#### HTML

```html
<a href="https://wails.io">
  <img src="https://wails.io/img/wails-button-light.svg" width="180" height="48" alt="Built with Wails">
</a>
```

작은 버전에는 `width="90" height="24"`를 사용하세요.

#### Markdown

```md
[![Built with Wails](https://wails.io/img/wails-button-light.svg)](https://wails.io)
```

### 어두운 배경

<div class="wails-media-preview wails-media-preview-dark">  
<p>표준 크기 (<code>180</code> × <code>48</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-dark.svg" width="180" height="48" alt="Wails로 제작" /></a>  
<p>작은 크기 (<code>90</code> × <code>24</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-dark.svg" width="90" height="24" alt="Wails로 제작" /></a>  
</div>

#### HTML

```html
<a href="https://wails.io">
  <img src="https://wails.io/img/wails-button-dark.svg" width="180" height="48" alt="Built with Wails">
</a>
```

작은 버전에는 `width="90" height="24"`를 사용하세요.

#### Markdown

```md
[![Built with Wails](https://wails.io/img/wails-button-dark.svg)](https://wails.io)
```

Markdown은 이미지를 원래 크기인 가로 `180`, 세로 `48`픽셀로 표시합니다. 다른 표시 크기를 지정하려면 HTML 버전을 사용하세요. 버튼의 원래 가로세로 비율인 `15:4`를 유지하세요.
