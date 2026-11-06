---
title: "미디어"
description: "미리보기와 HTML 및 Markdown 임베드가 포함된 공식 Wails 버튼"
slug: "community/media"
sourcePath: "community/media.md"
---

Wails로 빌드한 프로젝트를 공유할 때는 공식 Wails 에셋을 사용하세요. 아래 버튼은 [wails.io](https://wails.io)로 연결되며 Wails에서 계속 호스팅하므로 이미지를 프로젝트에 복사할 필요가 없습니다.

## Wails로 빌드됨

배경과 충분한 대비를 이루는 버전을 선택하세요. 각 미리보기는 방문자에게 표시될 링크된 이미지와 정확히 동일합니다.

### 밝은 배경

<div class="wails-media-preview wails-media-preview-light">  
<p>표준 (<code>180</code> × <code>48</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-light.svg" width="180" height="48" alt="Wails로 빌드됨" /></a>  
<p>컴팩트 (<code>90</code> × <code>24</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-light.svg" width="90" height="24" alt="Wails로 빌드됨" /></a>  
</div>

#### HTML

```html
<a href="https://wails.io">
  <img src="https://wails.io/img/wails-button-light.svg" width="180" height="48" alt="Built with Wails">
</a>
```

컴팩트 버전에는 `width="90" height="24"`을 사용하세요.

#### Markdown

```md
[![Built with Wails](https://wails.io/img/wails-button-light.svg)](https://wails.io)
```

### 어두운 배경

<div class="wails-media-preview wails-media-preview-dark">  
<p>표준 (<code>180</code> × <code>48</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-dark.svg" width="180" height="48" alt="Wails로 빌드됨" /></a>  
<p>컴팩트 (<code>90</code> × <code>24</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-dark.svg" width="90" height="24" alt="Wails로 빌드됨" /></a>  
</div>

#### HTML

```html
<a href="https://wails.io">
  <img src="https://wails.io/img/wails-button-dark.svg" width="180" height="48" alt="Built with Wails">
</a>
```

컴팩트 버전에는 `width="90" height="24"`을 사용하세요.

#### Markdown

```md
[![Built with Wails](https://wails.io/img/wails-button-dark.svg)](https://wails.io)
```

Markdown은 이미지를 고유 크기인 `180` × `48`픽셀로 렌더링합니다. 다른 표시 크기를 설정해야 할 때는 HTML 버전을 사용하세요. 버튼의 원래 `15:4` 비율을 유지하세요.
