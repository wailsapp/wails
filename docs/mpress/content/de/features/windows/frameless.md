---
title: "Rahmenlose Fenster"
description: "Benutzerdefinierte Fensterdekorationen mit rahmenlosen Fenstern erstellen"
slug: "features/windows/frameless"
sourcePath: "features/windows/frameless.md"
---

## Rahmenlose Fenster

Wails bietet **Unterstützung für rahmenlose Fenster** mit CSS-basierten Ziehbereichen und plattformeigenem Verhalten. Entfernen Sie die plattformeigene Titelleiste, um Fensterdekoration, individuelle Designs und besondere Benutzererlebnisse vollständig selbst zu gestalten, ohne dabei auf wichtige Funktionen wie Verschieben, Größenänderung und Systemsteuerelemente zu verzichten.

![Die standardmäßige TypeScript-Starter-App von Wails v3 als rahmenloses Fenster mit nativen macOS-Ecken](/assets/screenshots/frameless-v3-native-corners-macos.png)

Das obige Beispiel zeigt die standardmäßige TypeScript-Starter-App von Wails v3 mit aktiviertem `Frameless: true`.

## Schnellstart

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:     "Frameless App",
    Width:     800,
    Height:    600,
    Frameless: true,
})
```

**CSS für eine ziehbare Titelleiste:**

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

**Das ist alles!** Sie haben jetzt eine benutzerdefinierte Titelleiste.

## Rahmenlose Fenster erstellen

### Eckenradius (macOS)

Rahmenlose Fenster behalten standardmäßig die üblichen abgerundeten macOS-Ecken von AppKit bei. Legen Sie mit `Mac.CornerRadius` einen benutzerdefinierten Radius (in Punkten) fest:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Mac: application.MacWindow{
        CornerRadius: 16,
    },
})
```

Setzen Sie `Mac.CornerType` für eckige Ecken auf `MacWindowCornerTypeSquare`. Dadurch wird `CornerRadius` ignoriert:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Mac: application.MacWindow{
        CornerType: application.MacWindowCornerTypeSquare,
    },
})
```

### Einfaches rahmenloses Fenster

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Width:     800,
    Height:    600,
})
```

**Das erhalten Sie:**

- Keine Titelleiste
- Keine Fensterrahmen
- Keine Systemschaltflächen
- Transparenter Hintergrund (optional)

**Das müssen Sie implementieren:**

- Ziehbarer Bereich
- Schaltflächen zum Schließen, Minimieren und Maximieren
- Größenänderungsbereiche (falls die Größe veränderbar ist)

### Mit transparentem Hintergrund

**Private API unter macOS:** Legen Sie `Mac.Backdrop: application.MacBackdropTransparent` fest und erstellen Sie den Build mit `-tags private_mac_apis`, um die Webview transparent darzustellen. Ohne dieses Tag bleibt die native Webview undurchsichtig, selbst wenn der HTML-/CSS-Hintergrund transparent ist. `Frameless` und `TitleBar.AppearsTransparent` selbst verwenden öffentliche APIs. Siehe [Private macOS-APIs](/guides/build/private-macos-apis/#webview-transparency-and-background).

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless:      true,
    BackgroundType: application.BackgroundTypeTransparent,
    Mac: application.MacWindow{
        Backdrop: application.MacBackdropTransparent, // Requires -tags private_mac_apis.
    },
})
```

**Anwendungsfälle:**

- Abgerundete Ecken
- Benutzerdefinierte Formen
- Overlay-Fenster
- Startbildschirme

## Ziehbereiche

### CSS-basiertes Ziehen

Verwenden Sie die CSS-Eigenschaft `--wails-draggable`:

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

**Werte:**

- `drag` – Der Bereich ist ziehbar
- `no-drag` – Der Bereich ist nicht ziehbar (auch wenn der übergeordnete Bereich ziehbar ist)

### Vollständiges Beispiel für eine Titelleiste

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

**JavaScript für die Schaltflächen:**

```javascript
import { Window } from '@wailsio/runtime'

document.querySelector('.minimize').addEventListener('click', () => Window.Minimise())
document.querySelector('.maximize').addEventListener('click', () => Window.Maximise())
document.querySelector('.close').addEventListener('click', () => Window.Close())
```

## Native Nicht-Client-Bereiche unter Windows

Windows kann Teile einer benutzerdefinierten Titelleiste als native Nicht-Client-Bereiche behandeln. Dadurch können Sie die Titelleiste und die Fensterschaltflächen mit einem beliebigen HTML-/CSS-Design darstellen und gleichzeitig das native Windows-Verhalten beibehalten: Der Titelleistenbereich verschiebt das Fenster, die Maximieren-Schaltfläche kann Windows 11 Snap Assist / Snap Layouts anzeigen und die Schaltflächen zum Minimieren, Maximieren und Schließen erhalten native Treffertests und Mauszustände.

Das folgende Video zeigt eine benutzerdefinierte HTML-/CSS-Titelleiste mit nativen Windows-Treffertests, einschließlich Windows 11 Snap Assist / Snap Layouts für eine benutzerdefinierte Maximieren-Schaltfläche.

<video src="/assets/windows-native-non-client-regions/wails-app-region.mp4" controls muted playsInline></video>

Wails unterstützt zwei Windows-spezifische Mechanismen:

- `app-region` über die native Unterstützung von WebView2 für Nicht-Client-Bereiche
- `--wails-non-client-region` über die Laufzeitverfolgung von Wails für benutzerdefinierte Fensterschaltflächen

### Modus auswählen

@note{type="caution" title="Experimentell"}
`WebView2CompositionHosting` verändert im Hintergrund, wie das Fenster WebView2 hostet und damit interagiert. Anstelle des standardmäßigen, in einem HWND gehosteten WebView2-Controllers verwendet Wails das Hosting eines Composition-Controllers und leitet Eingaben explizit weiter. In diesem Modus können Probleme mit Darstellung, Eingabe, Fokus oder der Kompatibilität mit der WebView2 Runtime auftreten. Aktivieren Sie ihn nur, wenn Sie natives Verhalten für benutzerdefinierte Fensterschaltflächen benötigen, und testen Sie Ihre App sorgfältig mit den von Ihnen unterstützten Windows- und WebView2-Runtime-Versionen.

@end

Wählen Sie den Modus danach aus, welche Windows-Funktionen Sie benötigen:

- Verwenden Sie `NonClientRegionSupport` für einfaches natives Verschieben der App mit `app-region: drag` und `app-region: no-drag` von WebView2.
- Verwenden Sie `WebView2CompositionHosting`, wenn sich Ihre benutzerdefinierten Schaltflächen zum Minimieren, Maximieren und Schließen wie native Windows-Fensterschaltflächen verhalten sollen.
- Aktivieren Sie beide, wenn dasselbe Fenster sowohl die native `app-region`-Unterstützung von WebView2 als auch von Wails verwaltete Bereiche für benutzerdefinierte Fensterschaltflächen benötigt.

`NonClientRegionSupport` ist die schlanke native Alternative zur `--wails-draggable`-Verfolgung von Wails. Sie kennzeichnen ziehbare und nicht ziehbare Bereiche mit CSS, WebView2 bestimmt, welche Pixel zur Titelleiste gehören, und Wails fragt WebView2 beim Treffertest nach dem nativen Bereich.

Das ist derzeit der gesamte Funktionsumfang dieses Modus. Er bewirkt weder, dass sich benutzerdefinierte Schaltflächen zum Minimieren, Maximieren oder Schließen wie native Windows-Fensterschaltflächen verhalten, noch aktiviert er Windows 11 Snap Assist / Snap Layouts für eine benutzerdefinierte Maximieren-Schaltfläche. Verwenden Sie ihn, wenn Sie einfaches natives Verschieben der App ohne den zusätzlichen Mechanismus von `--wails-draggable` benötigen.

`WebView2CompositionHosting` ist für benutzerdefinierte Titelleistenschaltflächen mit nativem Verhalten vorgesehen. Wails verfolgt mit `--wails-non-client-region` markierte DOM-Rechtecke, ordnet ihnen Windows-Hit-Test-Werte wie `HTMINBUTTON`, `HTMAXBUTTON` und `HTCLOSE` zu und leitet Mauseingaben zurück an die per Composition gehostete WebView2-Oberfläche. Dadurch kann eine benutzerdefinierte Maximieren-Schaltfläche Windows 11 Snap Assist/Snap Layouts unterstützen, während Sie das visuelle Design frei bestimmen können.

Anders ausgedrückt: `NonClientRegionSupport` bietet WebView2-native Unterstützung für CSS-Regionen. Bei `WebView2CompositionHosting` übernimmt Wails die Verantwortung für die vom Host verwaltete Composition und benutzerdefinierte Nicht-Client-Hit-Tests.

### WebView2-app-region

Aktivieren Sie für das Fenster die native Unterstützung von WebView2 für Nicht-Client-Regionen:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        NonClientRegionSupport: true,
    },
})
```

Markieren Sie anschließend ziehbare Bereiche mit der CSS-Eigenschaft `app-region`:

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

Verwenden Sie diesen Modus, wenn Sie nur natives Ziehen über die Titelleiste benötigen und die Steuerelemente Ihrer Titelleiste über normale Frontend-Klicks bedient werden.

Dieser Modus ist auf die Unterstützung beschränkt, die WebView2 selbst für Nicht-Client-Regionen bietet. In aktuellen WebView2-Versionen sind das ausschließlich ziehbare und nicht ziehbare Regionen. Der Modus ist nicht dazu vorgesehen, vollständig benutzerdefinierte Titelleistenschaltflächen im Frontend mit getrennten nativen Rollen für Minimieren, Maximieren und Schließen abzubilden.

### Benutzerdefinierte Titelleistenschaltflächen mit nativem Verhalten

Aktivieren Sie das Composition-Hosting für benutzerdefinierte Schaltflächen zum Minimieren, Maximieren und Schließen, die sich wie Titelleistenschaltflächen des Systems verhalten sollen:

@note{type="caution" title="Experimentell"}
`WebView2CompositionHosting` verwendet das Hosting des WebView2-Composition-Controllers mit DirectComposition. Lesen Sie vor der Aktivierung den Abschnitt [Modus auswählen](#modus-auswhlen).

@end

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        WebView2CompositionHosting: true,
    },
})
```

Markieren Sie anschließend jede Frontend-Region mit `--wails-non-client-region`:

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

Unterstützte Werte für `--wails-non-client-region`:

- `caption` – ziehbarer Titelleistenbereich
- `minimize` – Hit-Test-Ziel der nativen Minimieren-Schaltfläche
- `maximize` – Hit-Test-Ziel der nativen Maximieren-Schaltfläche, einschließlich des Hover-Verhaltens von Windows 11 Snap Assist/Snap Layouts
- `close` – Hit-Test-Ziel der nativen Schließen-Schaltfläche

Die Wails-Laufzeit überwacht Änderungen am DOM, an Stilen, Größe, Bildlauf und Viewport und sendet anschließend Momentaufnahmen der Regionen an das native Fenster. Die Geometrie der Regionen wird in CSS-Pixeln gemessen und für Windows-Hit-Tests in physische Pixel umgerechnet.

Das visuelle Design bestimmen weiterhin vollständig Sie. Die Regionen teilen Windows lediglich die Bedeutung der einzelnen Rechtecke mit; Form, Symbol, Farbe, Abstände, Hover-Stil und Layout der Schaltflächen stammen weiterhin aus Ihrem Frontend.

### Beide kombinieren

Sie können beide Optionen aktivieren, wenn Sie im selben Fenster WebView2-Unterstützung für `app-region` und von Wails verwaltete Regionen für Titelleistenschaltflächen verwenden möchten:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        NonClientRegionSupport:    true,
        WebView2CompositionHosting: true,
    },
})
```

## Systemschaltflächen

### Schließen, Minimieren und Maximieren implementieren

**Go-Seite:**

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

**JavaScript-Seite:**

```javascript
import { Minimise, Maximise, Close } from './bindings/WindowControls'

document.querySelector('.minimize').addEventListener('click', Minimise)
document.querySelector('.maximize').addEventListener('click', Maximise)
document.querySelector('.close').addEventListener('click', Close)
```

**Alternativ können Sie Laufzeitmethoden verwenden:**

```javascript
import { Window } from '@wailsio/runtime'

document.querySelector('.minimize').addEventListener('click', () => Window.Minimise())
document.querySelector('.maximize').addEventListener('click', () => Window.Maximise())
document.querySelector('.close').addEventListener('click', () => Window.Close())
```

### Maximierungszustand umschalten

Verfolgen Sie den Maximierungszustand für das Schaltflächensymbol:

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

## Größenänderungsbereiche

### CSS-basierte Größenänderung

Wails stellt automatische Größenänderungsbereiche für rahmenlose Fenster bereit:

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

**Werte:**

- `all` – Größenänderung an allen Rändern
- `top`, `bottom`, `left`, `right` – bestimmte Ränder
- `top-left`, `top-right`, `bottom-left`, `bottom-right` – Ecken
- `none` – keine Größenänderung

### Beispiel für Größenänderungsbereiche

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

## Plattformspezifisches Verhalten

@tabs{sync-key="platform"}
[Windows]
**Rahmenlose Fenster unter Windows:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        DisableFramelessWindowDecorations: false,
    },
})
```

**Funktionen:**

- Automatischer Schlagschatten
- Unterstützung für Snap Layouts (Windows 11)
- Unterstützung für Aero Snap
- DPI-Skalierung

**Fensterdekorationen deaktivieren:**

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

Dadurch werden Snap Layouts über den Windows-Tastenkürzelpfad ausgelöst. Verwenden Sie für eine benutzerdefinierte HTML-Maximieren-Schaltfläche mit nativen Snap Layouts beim Darüberfahren stattdessen [Native Nicht-Client-Regionen unter Windows](#native-nicht-client-bereiche-unter-windows).

**Benutzerdefinierte Titelleistenhöhe:** Windows erkennt ziehbare Regionen automatisch anhand des CSS.

[macOS]
**Rahmenlose Fenster unter macOS:**

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

**Funktionen:**

- Native Vollbildunterstützung
- Ampelschaltflächen (optional)
- Vibrancy-Effekte
- Transparente Titelleiste

**Blenden Sie die Titelleiste vollständig aus** (verwenden Sie die aus dem Paket `application` exportierten Vorgabevarianten – es gibt weder ein Feld `TitleBarStyle` noch eine Konstante `MacTitleBarStyleHidden`):

```go
Mac: application.MacWindow{
    TitleBar: application.MacTitleBarHidden,
},
```

Weitere Vorgaben sind `MacTitleBarDefault`, `MacTitleBarHiddenInset` und `MacTitleBarHiddenInsetUnified`.

**Unsichtbare Titelleiste:** Ermöglicht das Ziehen bei ausgeblendeter Titelleiste. Dies wirkt sich nur aus, wenn das Fenster rahmenlos ist oder `AppearsTransparent` verwendet:

```go
Mac: application.MacWindow{
    InvisibleTitleBarHeight: 40,
},
```

[Linux]
**Rahmenlose Fenster unter Linux:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
})
```

**Funktionen:**

- Grundlegende Unterstützung für rahmenlose Fenster
- CSS-Ziehbereiche
- Je nach Desktop-Umgebung unterschiedlich

**Hinweise zu Desktop-Umgebungen:**

- **GNOME:** Gute Unterstützung
- **KDE Plasma:** Gute Unterstützung
- **XFCE:** Grundlegende Unterstützung
- **Kachelnde Fenstermanager:** Eingeschränkte Unterstützung

**Compositor erforderlich:** Transparenz erfordert einen Compositor (die meisten modernen Desktop-Umgebungen verfügen über einen).

@end

## Gängige Muster

### Muster 1: Moderne Titelleiste

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

### Muster 2: Startbildschirm

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

### Muster 3: Abgerundetes Fenster

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

### Muster 4: Overlay-Fenster

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

## Vollständiges Beispiel

Hier sehen Sie ein produktionsreifes rahmenloses Fenster:

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

## Bewährte Vorgehensweisen

### ✅ Empfohlen

- **Stellen Sie einen Ziehbereich bereit** – Benutzer müssen das Fenster verschieben können
- **Implementieren Sie Systemschaltflächen** – Schließen, Minimieren, Maximieren
- **Legen Sie eine Mindestgröße fest** – So verhindern Sie unbrauchbare Layouts
- **Testen Sie auf allen Plattformen** – Das Verhalten variiert
- **Verwenden Sie CSS für Ziehbereiche** – Flexibel und wartungsfreundlich
- **Geben Sie visuelles Feedback** – Hover-Zustände für Schaltflächen

### ❌ Nicht empfohlen

- **Vergessen Sie die Griffe zur Größenänderung nicht** – Wenn die Fenstergröße veränderbar ist
- **Machen Sie nicht das gesamte Fenster ziehbar** – Dies verhindert Interaktionen
- **Vergessen Sie bei Schaltflächen nicht, das Ziehen zu deaktivieren** – Andernfalls funktionieren sie nicht
- **Verwenden Sie keine winzigen Ziehbereiche** – Sie sind schwer zu treffen
- **Vergessen Sie die Plattformunterschiede nicht** – Testen Sie gründlich

## Fehlerbehebung

### Fenster lässt sich nicht ziehen

**Ursache:** `--wails-draggable: drag` fehlt

**Lösung:**

```css
.titlebar {
    --wails-draggable: drag;
}
```

### Schaltflächen funktionieren nicht

**Ursache:** Die Schaltflächen befinden sich in einem Ziehbereich

**Lösung:**

```css
.titlebar button {
    --wails-draggable: no-drag;
}
```

### Fenstergröße lässt sich nicht ändern

**Ursache:** Griffe zur Größenänderung fehlen

**Lösung:**

```css
body {
    --wails-resize: all;
}
```

## Nächste Schritte

@cards{cols="2"}
▣ Fenstergrundlagen
Lernen Sie die Grundlagen der Fensterverwaltung kennen.

[Mehr erfahren →](/features/windows/basics/)

---
⚙ Fensteroptionen
Vollständige Referenz der Fensteroptionen.

[Mehr erfahren →](/features/windows/options/)

---
🚀 Fensterereignisse
Verarbeiten Sie Ereignisse im Lebenszyklus eines Fensters.

[Mehr erfahren →](/features/windows/events/)

---
◆ Mehrere Fenster
Muster für Anwendungen mit mehreren Fenstern.

[Mehr erfahren →](/features/windows/multiple/)

@end

---

**Fragen?** Stellen Sie sie auf [Discord](https://discord.gg/JDdSxwjhGf) oder sehen Sie sich das [Beispiel für rahmenlose Fenster](https://github.com/wailsapp/wails/tree/master/v3/examples/frameless) an.
