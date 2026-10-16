---
title: "Окна без рамки"
description: "Создание собственного оформления окон без рамки"
slug: "features/windows/frameless"
sourcePath: "features/windows/frameless.md"
---

## Окна без рамки

Wails предоставляет **поддержку окон без рамки** с задаваемыми через CSS областями перетаскивания и нативным для платформы поведением. Удалите нативную строку заголовка, чтобы получить полный контроль над оформлением окна, создавать собственный дизайн и уникальный пользовательский интерфейс, сохраняя такие важные функции, как перетаскивание, изменение размера и системные элементы управления.

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

По умолчанию окна без рамки сохраняют стандартные скруглённые углы macOS, предоставляемые AppKit. Задайте `Mac.CornerRadius`, чтобы использовать собственный радиус (в пунктах):

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Mac: application.MacWindow{
        CornerRadius: 16,
    },
})
```

Установите для `Mac.CornerType` значение `MacWindowCornerTypeSquare`, чтобы сделать углы прямыми. При этом `CornerRadius` игнорируется:

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

**Закрытый API в macOS:** задайте `Mac.Backdrop: application.MacBackdropTransparent` и выполните сборку с `-tags private_mac_apis`, чтобы сделать WebView прозрачным. Без этого тега нативный WebView остаётся непрозрачным, даже если фон HTML/CSS прозрачен. Сами `Frameless` и `TitleBar.AppearsTransparent` используют общедоступные API. См. раздел [«Закрытые API macOS»](/guides/build/private-macos-apis/#webview-transparency-and-background).

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

### Перетаскивание с помощью CSS

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
- `no-drag` — область нельзя перетаскивать (даже если родительскую область можно)

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

Windows может обрабатывать части пользовательской строки заголовка как нативные неклиентские области. Благодаря этому строку заголовка и кнопки управления окном можно оформить в HTML/CSS любым способом, сохранив нативное поведение Windows: область заголовка позволяет перетаскивать окно, кнопка разворачивания может отображать интерфейс Windows 11 Snap Assist / Snap Layouts, а кнопки сворачивания, разворачивания и закрытия получают нативное определение попадания и состояние мыши.

В видео ниже показана пользовательская строка заголовка на HTML/CSS с нативным определением попадания Windows, включая интерфейс Windows 11 Snap Assist / Snap Layouts для пользовательской кнопки разворачивания.

<video src="/assets/windows-native-non-client-regions/wails-app-region.mp4" controls muted playsInline></video>

Wails поддерживает два механизма, предназначенных для Windows:

- `app-region` через нативную поддержку неклиентских областей в WebView2
- `--wails-non-client-region` через отслеживание в среде выполнения Wails для пользовательских кнопок управления окном

### Выбор режима

@note{type="caution" title="Экспериментальная функция"}
`WebView2CompositionHosting` изменяет внутренний механизм размещения WebView2 в окне и взаимодействия с ним. Вместо стандартного контроллера WebView2, размещённого в HWND, Wails использует контроллер композиции и явно перенаправляет ввод. В этом режиме могут возникать проблемы с отрисовкой, вводом, фокусом или совместимостью со средой выполнения WebView2. Включайте его только тогда, когда требуется нативное поведение пользовательских кнопок управления окном, и тщательно тестируйте приложение с поддерживаемыми вами версиями Windows и среды выполнения WebView2.

@end

Выберите режим в зависимости от того, какие возможности Windows вам нужны:

- Используйте `NonClientRegionSupport` для простого нативного перетаскивания приложения с помощью `app-region: drag` и `app-region: no-drag` в WebView2.
- Используйте `WebView2CompositionHosting`, если ваши пользовательские кнопки сворачивания, разворачивания и закрытия должны вести себя как нативные кнопки управления окном Windows.
- Включите оба механизма, если одному и тому же окну необходимы нативная поддержка `app-region` в WebView2 и управляемые Wails области пользовательских кнопок управления окном.

`NonClientRegionSupport` — это облегчённая нативная альтернатива отслеживанию `--wails-draggable` в Wails. Вы отмечаете перетаскиваемые и неперетаскиваемые области с помощью CSS, WebView2 определяет, какие пиксели относятся к заголовку, а Wails при определении попадания запрашивает у WebView2 нативную область.

На сегодняшний день возможности этого режима этим исчерпываются. Он не обеспечивает для пользовательских кнопок сворачивания, разворачивания и закрытия поведение нативных кнопок управления окном Windows и не включает интерфейс Windows 11 Snap Assist / Snap Layouts для пользовательской кнопки разворачивания. Используйте его, когда требуется простое нативное перетаскивание приложения без дополнительных механизмов `--wails-draggable`.

`WebView2CompositionHosting` предназначен для пользовательских кнопок заголовка с нативным поведением. Wails отслеживает прямоугольные области DOM, помеченные `--wails-non-client-region`, сопоставляет их со значениями проверки попадания Windows, такими как `HTMINBUTTON`, `HTMAXBUTTON` и `HTCLOSE`, а затем передаёт ввод мыши обратно на поверхность WebView2, размещённую с помощью композиции. Благодаря этому пользовательская кнопка развёртывания может участвовать в работе функций Windows 11 Snap Assist и Snap Layouts, сохраняя любой выбранный вами визуальный дизайн.

Иными словами, `NonClientRegionSupport` — это нативная поддержка CSS-областей в WebView2. При использовании `WebView2CompositionHosting` Wails берёт на себя композицию, управляемую хостом, и пользовательскую проверку попадания в неклиентскую область.

### Области app-region в WebView2

Включите для окна нативную поддержку неклиентских областей WebView2:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        NonClientRegionSupport: true,
    },
})
```

Затем пометьте перетаскиваемые области с помощью CSS-свойства `app-region`:

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

Используйте этот режим, если вам нужно только нативное перетаскивание за заголовок, а элементы управления заголовка обрабатываются обычными щелчками во фронтенде.

Ограничение этого режима состоит в том, что его возможности определяются собственной поддержкой неклиентских областей в WebView2. В текущих выпусках WebView2 доступны только области перетаскивания и области, не допускающие перетаскивание. Этот режим не предназначен для моделирования полностью пользовательских кнопок заголовка во фронтенде с отдельными нативными ролями свёртывания, развёртывания и закрытия.

### Пользовательские кнопки заголовка с нативным поведением

Для пользовательских кнопок свёртывания, развёртывания и закрытия, которые должны вести себя как системные кнопки заголовка, включите размещение с помощью композиции:

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
- `minimize` — нативная целевая область кнопки свёртывания
- `maximize` — нативная целевая область кнопки развёртывания, включая поведение функций Windows 11 Snap Assist и Snap Layouts при наведении указателя
- `close` — нативная целевая область кнопки закрытия

Среда выполнения Wails отслеживает изменения DOM, стилей, размеров, прокрутки и области просмотра, а затем отправляет снимки областей нативному окну. Геометрия областей измеряется в CSS-пикселях и преобразуется в физические пиксели для проверки попадания в Windows.

Визуальный дизайн полностью остаётся на ваше усмотрение. Области лишь сообщают Windows назначение каждого прямоугольника; форма, значок, цвет, интервалы, стиль при наведении и компоновка кнопок по-прежнему задаются во фронтенде.

### Совместное использование обоих режимов

Можно включить оба параметра, если в одном окне требуется поддержка `app-region` в WebView2 и области кнопок заголовка, управляемые Wails:

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

### Реализация закрытия, свёртывания и развёртывания

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

### Переключение состояния развёртывания

Отслеживайте состояние развёртывания, чтобы менять значок кнопки:

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

- `all` — изменение размера за любую границу
- `top`, `bottom`, `left`, `right` — отдельные границы
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

## Поведение на разных платформах

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

Это вызывает Snap Layouts посредством сочетания клавиш Windows. Если пользовательской HTML-кнопке развёртывания требуется нативное отображение Snap Layouts при наведении, вместо этого используйте [нативные неклиентские области в Windows](#----windows).

**Пользовательская высота заголовка:** Windows автоматически обнаруживает области перетаскивания по CSS.

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
- Кнопки управления окном («светофор»), необязательно
- Эффекты полупрозрачности и смешивания с фоном
- Прозрачная строка заголовка

**Полностью скройте строку заголовка** (используйте предустановленные варианты, экспортируемые из пакета `application` — поля `TitleBarStyle` и константы `MacTitleBarStyleHidden` не существует):

```go
Mac: application.MacWindow{
    TitleBar: application.MacTitleBarHidden,
},
```

Другие предустановленные варианты: `MacTitleBarDefault`, `MacTitleBarHiddenInset` и `MacTitleBarHiddenInsetUnified`.

**Невидимая строка заголовка:** Позволяет перетаскивать окно, скрывая строку заголовка. Действует только для окна без рамки или окна, использующего `AppearsTransparent`:

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
- CSS-области перетаскивания
- Зависит от среды рабочего стола

**Примечания о средах рабочего стола:**

- **GNOME:** хорошая поддержка
- **KDE Plasma:** хорошая поддержка
- **XFCE:** базовая поддержка
- **Тайловые оконные менеджеры:** ограниченная поддержка

**Требуется композитный менеджер:** Для прозрачности необходим композитный менеджер (он есть в большинстве современных сред рабочего стола).

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

### Шаблон 2: экран-заставка

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

Ниже приведён готовый к использованию в рабочей среде пример окна без рамки:

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

### ✅ Следует

- **Предусмотрите область перетаскивания** — пользователям нужно перемещать окно
- **Реализуйте системные кнопки** — закрытие, сворачивание и разворачивание
- **Задайте минимальный размер** — это предотвратит появление непригодных для использования макетов
- **Тестируйте на всех платформах** — поведение различается
- **Используйте CSS для областей перетаскивания** — это гибкое и удобное в сопровождении решение
- **Обеспечьте визуальную обратную связь** — состояния кнопок при наведении указателя

### ❌ Не следует

- **Не забудьте о маркерах изменения размера** — если размер окна можно изменять
- **Не делайте всё окно областью перетаскивания** — это препятствует взаимодействию
- **Не забудьте отключить перетаскивание для кнопок** — иначе они не будут работать
- **Не делайте области перетаскивания слишком маленькими** — за них трудно ухватиться
- **Не забывайте о различиях между платформами** — тщательно тестируйте приложение

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

## Дальнейшие шаги

@cards{cols="2"}
▣ Основы работы с окнами
Изучите основы управления окнами.

[Подробнее →](/features/windows/basics/)

---
⚙ Параметры окна
Полный справочник параметров окна.

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

**Остались вопросы?** Задайте их в [Discord](https://discord.gg/JDdSxwjhGf) или ознакомьтесь с [примером окна без рамки](https://github.com/wailsapp/wails/tree/master/v3/examples/frameless).
