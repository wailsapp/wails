---
title: "프레임 없는 창"
description: "프레임 없는 창으로 사용자 지정 창 장식 만들기"
slug: "features/windows/frameless"
sourcePath: "features/windows/frameless.md"
---

## 프레임 없는 창

Wails는 CSS 기반 드래그 영역과 플랫폼 네이티브 동작을 지원하는 <strong>프레임 없는 창 기능</strong>을 제공합니다. 플랫폼 네이티브 제목 표시줄을 제거하면 드래그, 크기 조절, 시스템 컨트롤과 같은 필수 기능을 유지하면서 창 장식, 사용자 지정 디자인, 고유한 사용자 경험을 완전히 제어할 수 있습니다.

![네이티브 macOS 모서리가 적용된 프레임 없는 창으로 실행 중인 기본 Wails v3 TypeScript 스타터 앱](/assets/screenshots/frameless-v3-native-corners-macos.png)

위 예시는 `Frameless: true`을 활성화한 기본 Wails v3 TypeScript 스타터 앱입니다.

## 빠른 시작

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:     "Frameless App",
    Width:     800,
    Height:    600,
    Frameless: true,
})
```

**드래그 가능한 제목 표시줄용 CSS:**

```css
.titlebar {
    --wails-draggable: drag;
    height: 40px;
    background: #333;
}

.titlebar button {
    --wails-draggable: no-drag;
}
```

**HTML:**

```html
<div class="titlebar">
    <span>My Application</span>
    <button onclick="window.close()">×</button>
</div>
```

**이것으로 끝입니다!** 이제 사용자 지정 제목 표시줄을 사용할 수 있습니다.

## 프레임 없는 창 만들기

### 모서리 반경(macOS)

프레임 없는 창에는 기본적으로 AppKit의 표준 둥근 macOS 모서리가 유지됩니다. 사용자 지정 반경(포인트 단위)을 사용하려면 `Mac.CornerRadius`을 설정하세요.

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Mac: application.MacWindow{
        CornerRadius: 16,
    },
})
```

각진 모서리를 사용하려면 `Mac.CornerType`을 `MacWindowCornerTypeSquare`로 설정하세요. 이 경우 `CornerRadius`은 무시됩니다.

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Mac: application.MacWindow{
        CornerType: application.MacWindowCornerTypeSquare,
    },
})
```

### 기본 프레임 없는 창

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Width:     800,
    Height:    600,
})
```

**제공되는 사항:**

- 제목 표시줄 없음
- 창 테두리 없음
- 시스템 버튼 없음
- 투명 배경(선택 사항)

**구현해야 하는 사항:**

- 드래그 가능한 영역
- 닫기/최소화/최대화 버튼
- 크기 조절 핸들(창 크기를 조절할 수 있는 경우)

### 투명 배경 사용

**macOS의 비공개 API:** WebView를 투명하게 만들려면 `Mac.Backdrop: application.MacBackdropTransparent`을 설정하고 `-tags private_mac_apis`로 빌드하세요. 이 태그가 없으면 HTML/CSS 배경이 투명하더라도 네이티브 WebView는 불투명하게 유지됩니다. `Frameless`과 `TitleBar.AppearsTransparent` 자체는 공개 API를 사용합니다. [macOS 비공개 API](/guides/build/private-macos-apis/#webview-transparency-and-background)를 참조하세요.

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless:      true,
    BackgroundType: application.BackgroundTypeTransparent,
    Mac: application.MacWindow{
        Backdrop: application.MacBackdropTransparent, // Requires -tags private_mac_apis.
    },
})
```

**사용 사례:**

- 둥근 모서리
- 사용자 지정 도형
- 오버레이 창
- 스플래시 화면

## 드래그 영역

### CSS 기반 드래그

`--wails-draggable` CSS 속성을 사용하세요.

```css
/* Draggable area */
.titlebar {
    --wails-draggable: drag;
}

/* Non-draggable elements within draggable area */
.titlebar button {
    --wails-draggable: no-drag;
}
```

**값:**

- `drag` - 드래그 가능한 영역
- `no-drag` - 드래그할 수 없는 영역(부모 요소를 드래그할 수 있더라도 적용됨)

### 전체 제목 표시줄 예제

```html
<div class="titlebar">
    <div class="title">My Application</div>
    <div class="controls">
        <button class="minimize">−</button>
        <button class="maximize">□</button>
        <button class="close">×</button>
    </div>
</div>
```

```css
.titlebar {
    --wails-draggable: drag;
    display: flex;
    justify-content: space-between;
    align-items: center;
    height: 40px;
    background: #2c2c2c;
    color: white;
    padding: 0 16px;
}

.title {
    font-size: 14px;
    user-select: none;
}

.controls {
    display: flex;
    gap: 8px;
}

.controls button {
    --wails-draggable: no-drag;
    width: 32px;
    height: 32px;
    border: none;
    background: transparent;
    color: white;
    font-size: 16px;
    cursor: pointer;
    border-radius: 4px;
}

.controls button:hover {
    background: rgba(255, 255, 255, 0.1);
}

.controls .close:hover {
    background: #e81123;
}
```

**버튼용 JavaScript:**

```javascript
import { Window } from '@wailsio/runtime'

document.querySelector('.minimize').addEventListener('click', () => Window.Minimise())
document.querySelector('.maximize').addEventListener('click', () => Window.Maximise())
document.querySelector('.close').addEventListener('click', () => Window.Close())
```

## Windows의 네이티브 비클라이언트 영역

Windows에서는 사용자 지정 제목 표시줄의 일부를 네이티브 비클라이언트 영역으로 처리할 수 있습니다. 이를 사용하면 어떤 HTML/CSS 디자인으로든 제목 표시줄과 캡션 버튼을 그리면서도 Windows 네이티브 동작을 유지할 수 있습니다. 캡션 영역으로 창을 드래그할 수 있고, 최대화 버튼에 Windows 11 Snap Assist / Snap Layouts를 표시할 수 있으며, 최소화, 최대화 및 닫기 버튼에 네이티브 적중 테스트와 마우스 상태가 적용됩니다.

아래 동영상에서는 네이티브 Windows 적중 테스트를 사용하는 사용자 지정 HTML/CSS 제목 표시줄을 보여 줍니다. 여기에는 사용자 지정 최대화 버튼의 Windows 11 Snap Assist / Snap Layouts도 포함됩니다.

<video src="/assets/windows-native-non-client-regions/wails-app-region.mp4" controls muted playsInline></video>

Wails는 Windows 전용 메커니즘 두 가지를 지원합니다.

- WebView2의 네이티브 비클라이언트 영역 지원을 통한 `app-region`
- 사용자 지정 캡션 버튼을 위한 Wails 런타임 추적을 통한 `--wails-non-client-region`

### 모드 선택

@note{type="caution" title="실험적 기능"}
`WebView2CompositionHosting`은 내부적으로 창이 WebView2를 호스팅하고 상호 작용하는 방식을 변경합니다. Wails는 기본 HWND 호스팅 WebView2 컨트롤러 대신 컴포지션 컨트롤러 호스팅을 사용하고 입력을 명시적으로 전달합니다. 이 모드에서는 렌더링, 입력, 포커스 또는 WebView2 Runtime 호환성 문제가 발생할 수 있습니다. 네이티브 사용자 지정 캡션 버튼 동작이 필요한 경우에만 활성화하고, 지원하는 Windows 및 WebView2 Runtime 버전에서 앱을 철저히 테스트하세요.

@end

Windows에서 필요한 기능에 따라 선택하세요.

- WebView2의 `app-region: drag` 및 `app-region: no-drag`을 사용하여 간단한 네이티브 앱 드래그를 구현하려면 `NonClientRegionSupport`을 사용하세요.
- 사용자 지정 최소화, 최대화 및 닫기 버튼이 네이티브 Windows 캡션 버튼처럼 동작해야 한다면 `WebView2CompositionHosting`을 사용하세요.
- 동일한 창에 WebView2 네이티브 `app-region` 지원과 Wails에서 관리하는 사용자 지정 캡션 버튼 영역이 모두 필요하다면 두 기능을 모두 활성화하세요.

`NonClientRegionSupport`은 Wails의 `--wails-draggable` 추적을 대신하는 가벼운 네이티브 방식입니다. CSS로 드래그 가능한 영역과 드래그할 수 없는 영역을 표시하면 WebView2가 캡션에 속하는 픽셀을 판별하고, Wails는 적중 테스트를 수행할 때 WebView2에 네이티브 영역을 요청합니다.

현재 이 모드가 제공하는 기능은 이것이 전부입니다. 사용자 지정 최소화, 최대화 또는 닫기 버튼이 네이티브 Windows 캡션 버튼처럼 동작하게 하지는 않으며, 사용자 지정 최대화 버튼에서 Windows 11 Snap Assist / Snap Layouts를 활성화하지도 않습니다. `--wails-draggable`의 추가 메커니즘 없이 간단한 네이티브 앱 드래그가 필요할 때 사용하세요.

`WebView2CompositionHosting`은 네이티브 동작을 지원하는 사용자 지정 캡션 버튼에 사용합니다. Wails는 `--wails-non-client-region`로 표시된 DOM 사각형을 추적하고, 이를 `HTMINBUTTON`, `HTMAXBUTTON`, `HTCLOSE` 등의 Windows 적중 테스트 값에 매핑한 다음, 마우스 입력을 컴포지션으로 호스팅된 WebView2 표면으로 다시 전달합니다. 따라서 원하는 시각적 디자인을 유지하면서도 사용자 지정 최대화 버튼이 Windows 11 Snap Assist / Snap Layouts에 연동될 수 있습니다.

달리 표현하면, `NonClientRegionSupport`은 WebView2 네이티브 CSS 영역 지원입니다. `WebView2CompositionHosting`은 Wails가 호스트 소유 컴포지션과 사용자 지정 비클라이언트 적중 테스트를 담당하는 방식입니다.

### WebView2 app-region

창에 WebView2의 네이티브 비클라이언트 영역 지원을 활성화합니다:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        NonClientRegionSupport: true,
    },
})
```

그런 다음 CSS `app-region` 속성으로 드래그 가능 영역을 표시합니다:

```css
.titlebar {
    app-region: drag;
}

.titlebar button,
.titlebar input,
.titlebar select,
.titlebar textarea {
    app-region: no-drag;
}
```

네이티브 캡션 드래그만 필요하고 제목 표시줄 컨트롤을 일반적인 프런트엔드 클릭으로 처리하는 경우에 사용합니다.

이 모드는 WebView2 자체의 비클라이언트 영역 지원 범위로 제한됩니다. 현재 WebView2 릴리스에서는 드래그 영역과 드래그 금지 영역만 지원한다는 의미입니다. 네이티브 최소화, 최대화, 닫기 역할이 각각 구분된 완전한 사용자 지정 프런트엔드 캡션 버튼을 모델링하기 위한 기능은 아닙니다.

### 네이티브 동작을 지원하는 사용자 지정 캡션 버튼

시스템 캡션 버튼처럼 동작해야 하는 사용자 지정 최소화, 최대화, 닫기 버튼을 만들려면 컴포지션 호스팅을 활성화합니다:

@note{type="caution" title="실험적 기능"}
`WebView2CompositionHosting`은 DirectComposition과 함께 WebView2 컴포지션 컨트롤러 호스팅을 사용합니다. 활성화하기 전에 [모드 선택](#--2)을 참조하세요.

@end

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        WebView2CompositionHosting: true,
    },
})
```

그런 다음 각 프런트엔드 영역을 `--wails-non-client-region`로 표시합니다:

```html
<div class="titlebar">
    <div class="title">My Application</div>
    <div class="window-controls">
        <button class="window-button minimize" aria-label="Minimize"></button>
        <button class="window-button maximize" aria-label="Maximize"></button>
        <button class="window-button close" aria-label="Close"></button>
    </div>
</div>
```

```css
.titlebar {
    --wails-non-client-region: caption;
    height: 40px;
}

.window-controls {
    display: flex;
    height: 100%;
}

.window-button {
    width: 46px;
    border: 0;
    background: transparent;
}

.window-button.minimize {
    --wails-non-client-region: minimize;
}

.window-button.maximize {
    --wails-non-client-region: maximize;
}

.window-button.close {
    --wails-non-client-region: close;
}
```

지원되는 `--wails-non-client-region` 값:

- `caption` - 드래그 가능한 캡션 영역
- `minimize` - 네이티브 최소화 버튼의 적중 대상
- `maximize` - Windows 11 Snap Assist / Snap Layouts의 마우스 오버 동작을 포함하는 네이티브 최대화 버튼의 적중 대상
- `close` - 네이티브 닫기 버튼의 적중 대상

Wails 런타임은 DOM, 스타일, 크기, 스크롤 및 뷰포트 변경을 관찰한 다음 영역 스냅샷을 네이티브 창으로 전송합니다. 영역의 기하 정보는 CSS 픽셀 단위로 측정되고 Windows 적중 테스트를 위해 물리적 픽셀로 변환됩니다.

시각적 디자인은 전적으로 개발자가 결정합니다. 영역은 각 사각형의 의미만 Windows에 알려 줍니다. 버튼의 모양, 아이콘, 색상, 간격, 마우스 오버 스타일 및 레이아웃은 계속 프런트엔드에서 정의합니다.

### 두 방식 함께 사용하기

같은 창에서 WebView2 `app-region` 지원과 Wails가 관리하는 캡션 버튼 영역을 모두 사용하려면 두 옵션을 함께 활성화할 수 있습니다:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        NonClientRegionSupport:    true,
        WebView2CompositionHosting: true,
    },
})
```

## 시스템 버튼

### 닫기/최소화/최대화 구현

**Go 측:**

```go
type WindowControls struct {
    window *application.WebviewWindow
}

func (wc *WindowControls) Minimise() {
    wc.window.Minimise()
}

func (wc *WindowControls) Maximise() {
    if wc.window.IsMaximised() {
        wc.window.UnMaximise()
    } else {
        wc.window.Maximise()
    }
}

func (wc *WindowControls) Close() {
    wc.window.Close()
}
```

**JavaScript 측:**

```javascript
import { Minimise, Maximise, Close } from './bindings/WindowControls'

document.querySelector('.minimize').addEventListener('click', Minimise)
document.querySelector('.maximize').addEventListener('click', Maximise)
document.querySelector('.close').addEventListener('click', Close)
```

**또는 런타임 메서드를 사용합니다:**

```javascript
import { Window } from '@wailsio/runtime'

document.querySelector('.minimize').addEventListener('click', () => Window.Minimise())
document.querySelector('.maximize').addEventListener('click', () => Window.Maximise())
document.querySelector('.close').addEventListener('click', () => Window.Close())
```

### 최대화 상태 전환

버튼 아이콘을 위해 최대화 상태를 추적합니다:

```javascript
import { Window } from '@wailsio/runtime'

async function toggleMaximise() {
    const isMaximised = await Window.IsMaximised()

    if (isMaximised) {
        await Window.Restore()
    } else {
        await Window.Maximise()
    }

    updateMaximiseButton()
}

async function updateMaximiseButton() {
    const isMaximised = await Window.IsMaximised()
    const button = document.querySelector('.maximize')
    button.textContent = isMaximised ? '❐' : '□'
}
```

## 크기 조절 핸들

### CSS 기반 크기 조절

Wails는 프레임 없는 창에 자동 크기 조절 핸들을 제공합니다:

```css
/* Enable resize on all edges */
body {
    --wails-resize: all;
}

/* Or specific edges */
.resize-top {
    --wails-resize: top;
}

.resize-bottom {
    --wails-resize: bottom;
}

.resize-left {
    --wails-resize: left;
}

.resize-right {
    --wails-resize: right;
}

/* Corners */
.resize-top-left {
    --wails-resize: top-left;
}

.resize-top-right {
    --wails-resize: top-right;
}

.resize-bottom-left {
    --wails-resize: bottom-left;
}

.resize-bottom-right {
    --wails-resize: bottom-right;
}
```

**값:**

- `all` - 모든 가장자리에서 크기 조절
- `top`, `bottom`, `left`, `right` - 특정 가장자리
- `top-left`, `top-right`, `bottom-left`, `bottom-right` - 모서리
- `none` - 크기 조절 안 함

### 크기 조절 핸들 예제

```html
<div class="window">
    <div class="titlebar">...</div>
    <div class="content">...</div>
    <div class="resize-handle resize-bottom-right"></div>
</div>
```

```css
.resize-handle {
    position: absolute;
    width: 16px;
    height: 16px;
}

.resize-bottom-right {
    --wails-resize: bottom-right;
    bottom: 0;
    right: 0;
    cursor: nwse-resize;
}
```

## 플랫폼별 동작

@tabs{sync-key="platform"}
[Windows]
**Windows의 프레임 없는 창:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        DisableFramelessWindowDecorations: false,
    },
})
```

**기능:**

- 자동 그림자
- Snap Layouts 지원(Windows 11)
- Aero Snap 지원
- DPI 배율 조정

**창 장식 비활성화:**

```go
Windows: application.WindowsWindow{
    DisableFramelessWindowDecorations: true,
},
```

**Snap Assist:**

```go
// Trigger Windows 11 Snap Assist
window.SnapAssist()
```

이렇게 하면 Windows 단축키 경로를 통해 Snap Layouts가 실행됩니다. 마우스 오버 시 네이티브 Snap Layouts를 표시하는 사용자 지정 HTML 최대화 버튼이 필요하면 대신 [Windows의 네이티브 비클라이언트 영역](#windows---)을 사용하세요.

**사용자 지정 제목 표시줄 높이:** Windows는 CSS에서 드래그 영역을 자동으로 감지합니다.

[macOS]
**macOS의 프레임 없는 창:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Mac: application.MacWindow{
        TitleBar: application.MacTitleBar{
            AppearsTransparent: true,
        },
        InvisibleTitleBarHeight: 40,
    },
})
```

**기능:**

- 네이티브 전체 화면 지원
- 신호등 버튼(선택 사항)
- 바이브런시 효과
- 투명한 제목 표시줄

**제목 표시줄을 완전히 숨기기**(`application` 패키지에서 내보낸 프리셋 변형을 사용하세요. `TitleBarStyle` 필드나 `MacTitleBarStyleHidden` 상수는 없습니다):

```go
Mac: application.MacWindow{
    TitleBar: application.MacTitleBarHidden,
},
```

그 밖의 프리셋으로는 `MacTitleBarDefault`, `MacTitleBarHiddenInset`, `MacTitleBarHiddenInsetUnified` 등이 있습니다.

**보이지 않는 제목 표시줄:** 제목 표시줄을 숨긴 상태에서도 창을 드래그할 수 있습니다. 창이 프레임리스이거나 `AppearsTransparent`을 사용할 때만 적용됩니다:

```go
Mac: application.MacWindow{
    InvisibleTitleBarHeight: 40,
},
```

[Linux]
**Linux 프레임리스 창:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
})
```

**기능:**

- 기본적인 프레임리스 지원
- CSS 드래그 영역
- 데스크톱 환경에 따라 다름

**데스크톱 환경 참고 사항:**

- **GNOME:** 지원 수준 좋음
- **KDE Plasma:** 지원 수준 좋음
- **XFCE:** 기본 지원
- **타일링 WM:** 제한적으로 지원

**컴포지터 필요:** 투명도를 사용하려면 컴포지터가 필요합니다(대부분의 최신 데스크톱 환경에는 컴포지터가 있습니다).

@end

## 일반적인 패턴

### 패턴 1: 현대적인 제목 표시줄

```html
<div class="modern-titlebar">
    <div class="app-icon">
        <img src="/icon.png" alt="App Icon">
    </div>
    <div class="title">My Application</div>
    <div class="controls">
        <button class="minimize">−</button>
        <button class="maximize">□</button>
        <button class="close">×</button>
    </div>
</div>
```

```css
.modern-titlebar {
    --wails-draggable: drag;
    display: flex;
    align-items: center;
    height: 40px;
    background: linear-gradient(to bottom, #3a3a3a, #2c2c2c);
    border-bottom: 1px solid #1a1a1a;
    padding: 0 16px;
}

.app-icon {
    --wails-draggable: no-drag;
    width: 24px;
    height: 24px;
    margin-right: 12px;
}

.title {
    flex: 1;
    font-size: 13px;
    color: #e0e0e0;
    user-select: none;
}

.controls {
    display: flex;
    gap: 1px;
}

.controls button {
    --wails-draggable: no-drag;
    width: 46px;
    height: 32px;
    border: none;
    background: transparent;
    color: #e0e0e0;
    font-size: 14px;
    cursor: pointer;
    transition: background 0.2s;
}

.controls button:hover {
    background: rgba(255, 255, 255, 0.1);
}

.controls .close:hover {
    background: #e81123;
    color: white;
}
```

### 패턴 2: 시작 화면

```go
splash := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:          "Loading...",
    Width:          400,
    Height:         300,
    Frameless:      true,
    AlwaysOnTop:    true,
    BackgroundType: application.BackgroundTypeTransparent,
    Mac: application.MacWindow{
        Backdrop: application.MacBackdropTransparent, // Requires -tags private_mac_apis.
    },
})
```

```css
body {
    background: transparent;
    display: flex;
    justify-content: center;
    align-items: center;
}

.splash {
    background: white;
    border-radius: 12px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
    padding: 40px;
    text-align: center;
}
```

### 패턴 3: 모서리가 둥근 창

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless:      true,
    BackgroundType: application.BackgroundTypeTransparent,
    Mac: application.MacWindow{
        Backdrop: application.MacBackdropTransparent, // Requires -tags private_mac_apis.
    },
})
```

```css
body {
    background: transparent;
    margin: 8px;
}

.window {
    background: white;
    border-radius: 16px;
    box-shadow: 0 4px 24px rgba(0, 0, 0, 0.15);
    overflow: hidden;
    height: calc(100vh - 16px);
}

.titlebar {
    --wails-draggable: drag;
    background: #f5f5f5;
    border-bottom: 1px solid #e0e0e0;
}
```

### 패턴 4: 오버레이 창

```go
overlay := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless:      true,
    AlwaysOnTop:    true,
    BackgroundType: application.BackgroundTypeTransparent,
    Mac: application.MacWindow{
        Backdrop: application.MacBackdropTransparent, // Requires -tags private_mac_apis.
    },
})
```

```css
body {
    background: transparent;
}

.overlay {
    background: rgba(0, 0, 0, 0.8);
    backdrop-filter: blur(10px);
    border-radius: 8px;
    padding: 20px;
}
```

## 전체 예제

다음은 프로덕션에 사용할 수 있는 프레임리스 창입니다:

**Go:**

```go
package main

import (
    _ "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name: "Frameless App",
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
    })

    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:     "Frameless Application",
        Width:     1000,
        Height:    700,
        MinWidth:  800,
        MinHeight: 600,
        Frameless: true,

        Mac: application.MacWindow{
            TitleBar: application.MacTitleBar{
                AppearsTransparent: true,
            },
            InvisibleTitleBarHeight: 40,
        },

        Windows: application.WindowsWindow{
            DisableFramelessWindowDecorations: false,
        },
    })

    window.Center()
    window.Show()

    app.Run()
}
```

**HTML:**

```html
<!DOCTYPE html>
<html>
<head>
    <link rel="stylesheet" href="/style.css">
</head>
<body>
    <div class="window">
        <div class="titlebar">
            <div class="title">Frameless Application</div>
            <div class="controls">
                <button class="minimize" title="Minimise">−</button>
                <button class="maximize" title="Maximise">□</button>
                <button class="close" title="Close">×</button>
            </div>
        </div>
        <div class="content">
            <h1>Hello from Frameless Window!</h1>
        </div>
    </div>
    <script src="/main.js" type="module"></script>
</body>
</html>
```

**CSS:**

```css
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    background: #f5f5f5;
}

.window {
    height: 100vh;
    display: flex;
    flex-direction: column;
}

.titlebar {
    --wails-draggable: drag;
    display: flex;
    justify-content: space-between;
    align-items: center;
    height: 40px;
    background: #ffffff;
    border-bottom: 1px solid #e0e0e0;
    padding: 0 16px;
}

.title {
    font-size: 13px;
    font-weight: 500;
    color: #333;
    user-select: none;
}

.controls {
    display: flex;
    gap: 8px;
}

.controls button {
    --wails-draggable: no-drag;
    width: 32px;
    height: 32px;
    border: none;
    background: transparent;
    color: #666;
    font-size: 16px;
    cursor: pointer;
    border-radius: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s;
}

.controls button:hover {
    background: #f0f0f0;
    color: #333;
}

.controls .close:hover {
    background: #e81123;
    color: white;
}

.content {
    flex: 1;
    padding: 40px;
    overflow: auto;
}
```

**JavaScript:**

```javascript
import { Window } from '@wailsio/runtime'

// Minimise button
document.querySelector('.minimize').addEventListener('click', () => {
    Window.Minimise()
})

// Maximise/restore button
const maximiseBtn = document.querySelector('.maximize')
maximiseBtn.addEventListener('click', async () => {
    const isMaximised = await Window.IsMaximised()

    if (isMaximised) {
        await Window.Restore()
    } else {
        await Window.Maximise()
    }

    updateMaximiseButton()
})

// Close button
document.querySelector('.close').addEventListener('click', () => {
    Window.Close()
})

// Update maximise button icon
async function updateMaximiseButton() {
    const isMaximised = await Window.IsMaximised()
    maximiseBtn.textContent = isMaximised ? '❐' : '□'
    maximiseBtn.title = isMaximised ? 'Restore' : 'Maximise'
}

// Initial state
updateMaximiseButton()
```

## 모범 사례

### ✅ 권장 사항

- **드래그 가능한 영역을 제공하세요** - 사용자가 창을 이동할 수 있어야 합니다
- **시스템 버튼을 구현하세요** - 닫기, 최소화, 최대화
- **최소 크기를 설정하세요** - 사용할 수 없는 레이아웃이 되지 않도록 합니다
- **모든 플랫폼에서 테스트하세요** - 동작이 플랫폼마다 다릅니다
- **드래그 영역에는 CSS를 사용하세요** - 유연하고 유지 관리하기 쉽습니다
- **시각적 피드백을 제공하세요** - 버튼에 호버 상태를 적용합니다

### ❌ 피해야 할 사항

- **크기 조절 핸들을 빠뜨리지 마세요** - 창 크기를 조절할 수 있는 경우에 필요합니다
- **창 전체를 드래그 가능하게 만들지 마세요** - 사용자 조작을 방해합니다
- **버튼에 드래그 금지 영역을 지정하는 것을 잊지 마세요** - 지정하지 않으면 버튼이 작동하지 않습니다
- **드래그 영역을 너무 작게 만들지 마세요** - 잡기 어렵습니다
- **플랫폼별 차이를 간과하지 마세요** - 철저히 테스트하세요

## 문제 해결

### 창을 드래그할 수 없음

**원인:** `--wails-draggable: drag` 누락

**해결 방법:**

```css
.titlebar {
    --wails-draggable: drag;
}
```

### 버튼이 작동하지 않음

**원인:** 버튼이 드래그 가능 영역 안에 있음

**해결 방법:**

```css
.titlebar button {
    --wails-draggable: no-drag;
}
```

### 창 크기를 조절할 수 없음

**원인:** 크기 조절 핸들 누락

**해결 방법:**

```css
body {
    --wails-resize: all;
}
```

## 다음 단계

@cards{cols="2"}
▣ 창 기본 사항
창 관리의 기본 사항을 알아보세요.

[자세히 알아보기 →](/features/windows/basics/)

---
⚙ 창 옵션
창 옵션에 관한 전체 레퍼런스입니다.

[자세히 알아보기 →](/features/windows/options/)

---
🚀 창 이벤트
창 수명 주기 이벤트를 처리하세요.

[자세히 알아보기 →](/features/windows/events/)

---
◆ 여러 창
다중 창 애플리케이션을 위한 패턴입니다.

[자세히 알아보기 →](/features/windows/multiple/)

@end

---

**궁금한 점이 있으신가요?** [Discord](https://discord.gg/JDdSxwjhGf)에서 질문하거나 [프레임리스 예제](https://github.com/wailsapp/wails/tree/master/v3/examples/frameless)를 확인하세요.
