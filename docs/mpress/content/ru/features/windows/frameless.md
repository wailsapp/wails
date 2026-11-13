---
title: "Окна без рамки"
description: "Создание пользовательского оформления окна без рамки"
slug: "features/windows/frameless"
sourcePath: "features/windows/frameless.md"
---

## Окна без рамки

Wails предоставляет **поддержку окон без рамки** с областями перетаскивания на основе CSS и нативным для платформы поведением. Уберите нативную строку заголовка, чтобы полностью управлять оформлением окна, создавать собственный дизайн и уникальный пользовательский интерфейс, сохраняя при этом основные возможности: перетаскивание, изменение размера и системные элементы управления.

![Стандартное стартовое приложение Wails v3 на TypeScript, запущенное в окне без рамки с нативными скруглёнными углами macOS](/assets/screenshots/frameless-v3-native-corners-macos.png)

В примере выше показано стандартное стартовое приложение Wails v3 на TypeScript с включённым параметром `Frameless: true`.

## Быстрый старт

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:     "Frameless App",
    Width:     800,
    Height:    600,
    Frameless: true,
})
```

**CSS для перетаскиваемой строки заголовка:**

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

**Готово!** Теперь у вас есть собственная строка заголовка.

## Создание окон без рамки

### Радиус скругления углов (macOS)

По умолчанию окна без рамки сохраняют стандартные для AppKit скруглённые углы macOS. Чтобы задать собственный радиус в пунктах, установите `Mac.CornerRadius`:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Mac: application.MacWindow{
        CornerRadius: 16,
    },
})
```

Чтобы углы были прямыми, задайте для `Mac.CornerType` значение `MacWindowCornerTypeSquare`. При этом `CornerRadius` игнорируется:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Mac: application.MacWindow{
        CornerType: application.MacWindowCornerTypeSquare,
    },
})
```

### Базовое окно без рамки

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Width:     800,
    Height:    600,
})
```

**Что вы получаете:**

- Без строки заголовка
- Без границ окна
- Без системных кнопок
- Прозрачный фон (необязательно)

**Что необходимо реализовать:**

- Область перетаскивания
- Кнопки закрытия, сворачивания и разворачивания
- Маркеры изменения размера (если размер окна можно изменять)

### С прозрачным фоном

**Закрытый API в macOS:** для прозрачности WebView задайте `Mac.Backdrop: application.MacBackdropTransparent` и выполните сборку с `-tags private_mac_apis`. Без этого тега нативный WebView остаётся непрозрачным, даже если фон HTML/CSS прозрачен. Сами `Frameless` и `TitleBar.AppearsTransparent` используют общедоступные API. См. раздел [«Закрытые API macOS»](/guides/build/private-macos-apis/#webview-transparency-and-background).

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless:      true,
    BackgroundType: application.BackgroundTypeTransparent,
    Mac: application.MacWindow{
        Backdrop: application.MacBackdropTransparent, // Requires -tags private_mac_apis.
    },
})
```

**Варианты использования:**

- Скруглённые углы
- Произвольные формы
- Окна-наложения
- Заставки

## Области перетаскивания

### Перетаскивание на основе CSS

Используйте CSS-свойство `--wails-draggable`:

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

**Значения:**

- `drag` — область можно перетаскивать
- `no-drag` — область нельзя перетаскивать (даже если можно перетаскивать родительскую область)

### Полный пример строки заголовка

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

**JavaScript для кнопок:**

```javascript
import { Window } from '@wailsio/runtime'

document.querySelector('.minimize').addEventListener('click', () => Window.Minimise())
document.querySelector('.maximize').addEventListener('click', () => Window.Maximise())
document.querySelector('.close').addEventListener('click', () => Window.Close())
```

## Нативные неклиентские области в Windows

Windows может обрабатывать части пользовательской строки заголовка как нативные неклиентские области. Это позволяет оформлять строку заголовка и кнопки управления окном с помощью любого дизайна HTML/CSS, сохраняя нативное поведение Windows: область заголовка перетаскивает окно, кнопка разворачивания может отображать Windows 11 Snap Assist / Snap Layouts, а кнопки сворачивания, разворачивания и закрытия получают нативную проверку попадания и обработку состояния мыши.

В видео ниже показана пользовательская строка заголовка на HTML/CSS, использующая нативную проверку попадания Windows, включая Windows 11 Snap Assist / Snap Layouts для пользовательской кнопки разворачивания.

<video src="/assets/windows-native-non-client-regions/wails-app-region.mp4" controls muted playsInline></video>

Wails поддерживает два механизма, предназначенных для Windows:

- `app-region` через нативную поддержку неклиентских областей в WebView2
- `--wails-non-client-region` через отслеживание пользовательских кнопок управления окном средой выполнения Wails

### Выбор режима

@note{type="caution" title="Экспериментальная возможность"}
`WebView2CompositionHosting` изменяет внутренний механизм размещения WebView2 в окне и взаимодействия с ним. Вместо стандартного контроллера WebView2, размещённого в HWND, Wails использует размещение с помощью контроллера композиции и явно перенаправляет ввод. В этом режиме возможны проблемы с отрисовкой, вводом, фокусом или совместимостью со средой выполнения WebView2. Включайте его только тогда, когда требуется нативное поведение пользовательских кнопок управления окном, и тщательно тестируйте приложение с поддерживаемыми версиями Windows и среды выполнения WebView2.

@end

Выберите режим в зависимости от требуемых возможностей Windows:

- Для простого нативного перетаскивания приложения используйте `NonClientRegionSupport` вместе с `app-region: drag` и `app-region: no-drag` из WebView2.
- Используйте `WebView2CompositionHosting`, если пользовательские кнопки сворачивания, разворачивания и закрытия должны вести себя как нативные кнопки управления окном Windows.
- Включите оба режима, если одному окну одновременно требуются нативная для WebView2 поддержка `app-region` и управляемые Wails пользовательские области кнопок управления окном.

`NonClientRegionSupport` — это облегчённая нативная альтернатива отслеживанию `--wails-draggable` в Wails. Вы помечаете перетаскиваемые и неперетаскиваемые области с помощью CSS, WebView2 определяет, какие пиксели относятся к области заголовка, а Wails при проверке попадания запрашивает у WebView2 нативную область.

На сегодняшний день возможности этого режима этим исчерпываются. Он не заставляет пользовательские кнопки сворачивания, разворачивания и закрытия вести себя как нативные кнопки управления окном Windows и не включает Windows 11 Snap Assist / Snap Layouts для пользовательской кнопки разворачивания. Используйте его, если требуется простое нативное перетаскивание приложения без дополнительных механизмов `--wails-draggable`.

`WebView2CompositionHosting` предназначен для пользовательских кнопок управления окном с нативным поведением. Wails отслеживает прямоугольные области DOM, помеченные с помощью `--wails-non-client-region`, сопоставляет их со значениями проверки попадания Windows, такими как `HTMINBUTTON`, `HTMAXBUTTON` и `HTCLOSE`, и перенаправляет ввод мыши обратно на поверхность WebView2, размещённую средствами композиции. Благодаря этому пользовательская кнопка разворачивания может участвовать в работе Windows 11 Snap Assist / Snap Layouts независимо от выбранного визуального дизайна.

Иными словами, `NonClientRegionSupport` — это нативная для WebView2 поддержка CSS-областей. При использовании `WebView2CompositionHosting` Wails отвечает за композицию, которой управляет хост, и пользовательскую проверку попадания в неклиентские области.

### app-region в WebView2

Включите для окна нативную поддержку неклиентских областей WebView2:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        NonClientRegionSupport: true,
    },
})
```

Затем пометьте перетаскиваемые области CSS-свойством `app-region`:

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

Используйте этот режим, если требуется только нативное перетаскивание за область заголовка, а элементы управления строки заголовка обрабатываются обычными нажатиями во фронтенде.

Ограничение этого режима заключается в том, что его возможности определяются собственной поддержкой неклиентских областей в WebView2. В текущих версиях WebView2 это означает поддержку только перетаскиваемых и неперетаскиваемых областей. Этот режим не предназначен для моделирования полностью пользовательских кнопок управления окном во фронтенде с отдельными нативными ролями сворачивания, разворачивания и закрытия.

### Пользовательские кнопки управления окном с нативным поведением

Для пользовательских кнопок сворачивания, разворачивания и закрытия, которые должны вести себя как системные кнопки управления окном, включите размещение с помощью композиции:

@note{type="caution" title="Экспериментальная возможность"}
`WebView2CompositionHosting` использует размещение контроллера композиции WebView2 с DirectComposition. Перед включением ознакомьтесь с разделом [«Выбор режима»](#--2).

@end

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        WebView2CompositionHosting: true,
    },
})
```

Затем пометьте каждую область фронтенда с помощью `--wails-non-client-region`:

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

Поддерживаемые значения `--wails-non-client-region`:

- `caption` — перетаскиваемая область заголовка
- `minimize` — область нажатия нативной кнопки сворачивания
- `maximize` — область нажатия нативной кнопки разворачивания, включая поведение Windows 11 при наведении для Snap Assist и Snap Layouts
- `close` — область нажатия нативной кнопки закрытия

Среда выполнения Wails отслеживает изменения DOM, стилей, размеров, прокрутки и области просмотра, а затем отправляет снимки областей нативному окну. Геометрия областей измеряется в CSS-пикселях и преобразуется в физические пиксели для проверки попадания в Windows.

Визуальный дизайн полностью остаётся на ваше усмотрение. Области лишь сообщают Windows назначение каждого прямоугольника; форма кнопок, значки, цвета, интервалы, стили при наведении и компоновка по-прежнему задаются во фронтенде.

### Совместное использование

Оба параметра можно включить, если в одном окне нужны поддержка WebView2 `app-region` и управляемые Wails области кнопок заголовка:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        NonClientRegionSupport:    true,
        WebView2CompositionHosting: true,
    },
})
```

## Системные кнопки

### Реализация закрытия, сворачивания и разворачивания

**На стороне Go:**

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

**На стороне JavaScript:**

```javascript
import { Minimise, Maximise, Close } from './bindings/WindowControls'

document.querySelector('.minimize').addEventListener('click', Minimise)
document.querySelector('.maximize').addEventListener('click', Maximise)
document.querySelector('.close').addEventListener('click', Close)
```

**Либо используйте методы среды выполнения:**

```javascript
import { Window } from '@wailsio/runtime'

document.querySelector('.minimize').addEventListener('click', () => Window.Minimise())
document.querySelector('.maximize').addEventListener('click', () => Window.Maximise())
document.querySelector('.close').addEventListener('click', () => Window.Close())
```

### Переключение развёрнутого состояния

Отслеживайте развёрнутое состояние, чтобы изменять значок кнопки:

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

## Области изменения размера

### Изменение размера с помощью CSS

Wails предоставляет автоматические области изменения размера для окон без рамки:

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

**Значения:**

- `all` — изменение размера со всех краёв
- `top`, `bottom`, `left`, `right` — отдельные края
- `top-left`, `top-right`, `bottom-left`, `bottom-right` — углы
- `none` — изменение размера отключено

### Пример области изменения размера

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

## Особенности поведения на разных платформах

@tabs{sync-key="platform"}
[Windows]
**Окна без рамки в Windows:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        DisableFramelessWindowDecorations: false,
    },
})
```

**Возможности:**

- Автоматическая тень
- Поддержка Snap Layouts (Windows 11)
- Поддержка Aero Snap
- Масштабирование с учётом DPI

**Отключение оформления окна:**

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

Это вызывает Snap Layouts через сочетание клавиш Windows. Если нужна пользовательская HTML-кнопка разворачивания с нативным отображением Snap Layouts при наведении, вместо этого используйте раздел [Нативные неклиентские области в Windows](#----windows).

**Пользовательская высота заголовка:** Windows автоматически определяет области перетаскивания из CSS.

[macOS]
**Окна без рамки в macOS:**

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

**Возможности:**

- Нативная поддержка полноэкранного режима
- Кнопки «светофора» (необязательно)
- Эффекты полупрозрачности
- Прозрачная строка заголовка

**Полностью скройте строку заголовка** (используйте предустановленные варианты, экспортируемые из пакета `application`, — поля `TitleBarStyle` и константы `MacTitleBarStyleHidden` не существует):

```go
Mac: application.MacWindow{
    TitleBar: application.MacTitleBarHidden,
},
```

Среди других предустановок доступны `MacTitleBarDefault`, `MacTitleBarHiddenInset` и `MacTitleBarHiddenInsetUnified`.

**Невидимая строка заголовка:** позволяет перетаскивать окно при скрытой строке заголовка. Параметр действует, только если окно не имеет рамки или использует `AppearsTransparent`:

```go
Mac: application.MacWindow{
    InvisibleTitleBarHeight: 40,
},
```

[Linux]
**Окна без рамки в Linux:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
})
```

**Возможности:**

- Базовая поддержка окон без рамки
- Области перетаскивания, заданные в CSS
- Поведение зависит от среды рабочего стола

**Примечания по средам рабочего стола:**

- **GNOME:** хорошая поддержка
- **KDE Plasma:** хорошая поддержка
- **XFCE:** базовая поддержка
- **Тайловые оконные менеджеры:** ограниченная поддержка

**Требуется композитор:** для прозрачности необходим композитор (он есть в большинстве современных сред рабочего стола).

@end

## Распространённые шаблоны

### Шаблон 1: современная строка заголовка

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

### Шаблон 2: заставка

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

### Шаблон 3: окно со скруглёнными углами

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

### Шаблон 4: окно-оверлей

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

## Полный пример

Ниже приведён готовый к эксплуатации пример окна без рамки:

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

## Рекомендации

### ✅ Делайте так

- **Предусмотрите область перетаскивания** — пользователям нужно перемещать окно
- **Реализуйте системные кнопки** — закрытие, сворачивание и разворачивание
- **Задайте минимальный размер** — это предотвратит появление непригодных для использования компоновок
- **Проверяйте на всех платформах** — поведение различается
- **Задавайте области перетаскивания с помощью CSS** — это гибкое и удобное в сопровождении решение
- **Добавьте визуальную обратную связь** — состояния кнопок при наведении

### ❌ Не делайте так

- **Не забывайте об областях изменения размера** — если размер окна можно изменять
- **Не делайте всё окно перетаскиваемым** — это препятствует взаимодействию
- **Не забывайте отключать перетаскивание для кнопок** — иначе они не будут работать
- **Не делайте области перетаскивания слишком маленькими** — за них будет трудно ухватиться
- **Не забывайте о различиях между платформами** — тщательно всё проверяйте

## Устранение неполадок

### Окно не перетаскивается

**Причина:** отсутствует `--wails-draggable: drag`

**Решение:**

```css
.titlebar {
    --wails-draggable: drag;
}
```

### Кнопки не работают

**Причина:** кнопки находятся в области перетаскивания

**Решение:**

```css
.titlebar button {
    --wails-draggable: no-drag;
}
```

### Не удаётся изменить размер окна

**Причина:** отсутствуют маркеры изменения размера

**Решение:**

```css
body {
    --wails-resize: all;
}
```

## Следующие шаги

@cards{cols="2"}
▣ Основы работы с окнами
Изучите основы управления окнами.

[Подробнее →](/features/windows/basics/)

---
⚙ Параметры окна
Полный справочник по параметрам окна.

[Подробнее →](/features/windows/options/)

---
🚀 События окна
Обрабатывайте события жизненного цикла окна.

[Подробнее →](/features/windows/events/)

---
◆ Несколько окон
Шаблоны для многооконных приложений.

[Подробнее →](/features/windows/multiple/)

@end

---

**Есть вопросы?** Задайте их в [Discord](https://discord.gg/JDdSxwjhGf) или ознакомьтесь с [примером безрамочного окна](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/frameless).
