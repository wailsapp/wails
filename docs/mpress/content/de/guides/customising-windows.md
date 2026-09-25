---
title: "Fenster in Wails anpassen"
description: "Passen Sie das Erscheinungsbild und Verhalten von Fenstern in Ihren Wails-Anwendungen an"
slug: "guides/customising-windows"
sourcePath: "guides/customising-windows.md"
---

Relevante Plattformen: <span class="mpress-badge mpress-badge-note">Windows</span> <span class="mpress-badge mpress-badge-success">macOS</span>

<br/>

Wails stellt eine API bereit, mit der sich das Erscheinungsbild und die Funktionalität der Bedienelemente eines Fensters steuern lassen. Diese Funktionalität ist unter Windows und macOS verfügbar, jedoch nicht unter Linux.

## Zustände der Fensterschaltflächen festlegen

Die Schaltflächenzustände werden durch die Enumeration `ButtonState` definiert:

```go
type ButtonState int

const (
    ButtonEnabled   ButtonState = 0
    ButtonDisabled  ButtonState = 1
    ButtonHidden    ButtonState = 2
)
```

- `ButtonEnabled`: Die Schaltfläche ist aktiviert und sichtbar.
- `ButtonDisabled`: Die Schaltfläche ist sichtbar, aber deaktiviert (ausgegraut).
- `ButtonHidden`: Die Schaltfläche wird in der Titelleiste ausgeblendet.

Die Schaltflächenzustände können beim Erstellen des Fensters oder zur Laufzeit festgelegt werden.

### Schaltflächenzustände beim Erstellen des Fensters festlegen

Beim Erstellen eines neuen Fensters können Sie den Anfangszustand der Schaltflächen mit der Struktur `WebviewWindowOptions` festlegen:

```go {title="main.go"}
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name:        "My Application",
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        MinimiseButtonState:   application.ButtonHidden,
        MaximiseButtonState:   application.ButtonDisabled,
        CloseButtonState:      application.ButtonEnabled,
        FullscreenButtonState: application.ButtonEnabled,
    })

    app.Run()
}
```

Im obigen Beispiel ist die Minimieren-Schaltfläche ausgeblendet, die Maximieren-Schaltfläche deaktiviert (ausgegraut) und die Schließen-Schaltfläche aktiviert.

### Schaltflächenzustände zur Laufzeit festlegen

Sie können die Schaltflächenzustände auch zur Laufzeit mit den folgenden Methoden der Schnittstelle `Window` ändern:

```go
window.SetMinimiseButtonState(wails.ButtonHidden)
window.SetMaximiseButtonState(wails.ButtonEnabled)
window.SetCloseButtonState(wails.ButtonDisabled)
window.SetFullscreenButtonState(wails.ButtonEnabled)
```

### macOS: MaximiseButtonState und FullscreenButtonState verwenden dieselbe Schaltfläche

Unter macOS ist die grüne Ampelschaltfläche (`NSWindowZoomButton`) für das Maximieren und den Vollbildmodus dasselbe physische Bedienelement. Würden `MaximiseButtonState` und `FullscreenButtonState` beim Erstellen des Fensters auf unterschiedliche Werte gesetzt, käme es andernfalls unbemerkt zu einer Überschreibung, bei der der zuletzt gesetzte Wert Vorrang hat.

Um dies zu vermeiden, verwendet Wails bei der Initialisierung den **restriktiveren** der beiden Zustände, und zwar in der Reihenfolge `ButtonEnabled` < `ButtonDisabled` < `ButtonHidden`.

| `MaximiseButtonState` | `FullscreenButtonState` | Wirksamer Zustand unter macOS |
| --- | --- | --- |
| `ButtonEnabled` | `ButtonEnabled` | `ButtonEnabled` |
| `ButtonDisabled` | `ButtonEnabled` | `ButtonDisabled` |
| `ButtonEnabled` | `ButtonHidden` | `ButtonHidden` |
| `ButtonDisabled` | `ButtonHidden` | `ButtonHidden` |

Zur Laufzeit steuern `SetMaximiseButtonState` und `SetFullscreenButtonState` unter macOS beide `NSWindowZoomButton` an, sodass der letzte Aufruf Vorrang hat.

### Plattformunterschiede

Die Funktionalität für Schaltflächenzustände verhält sich unter Windows und macOS geringfügig unterschiedlich:

|  | Windows | Mac |
| --- | --- | --- |
| Minimieren/Maximieren/Schließen deaktivieren | Deaktiviert Minimieren/Maximieren/Schließen | Deaktiviert Minimieren/Maximieren/Schließen |
| Minimieren ausblenden | Deaktiviert Minimieren | Blendet die Minimieren-Schaltfläche aus |
| Maximieren ausblenden | Deaktiviert Maximieren | Blendet die Maximieren-Schaltfläche aus |
| Schließen ausblenden | Blendet alle Bedienelemente aus | Blendet Schließen aus |
| `FullscreenButtonState` | Keine Auswirkung | Steuert die Zoom-Schaltfläche (grün) |

Hinweis: Unter Windows können die Schaltflächen zum Minimieren und Maximieren nicht einzeln ausgeblendet werden. Wenn Sie jedoch beide deaktivieren, werden beide Bedienelemente ausgeblendet und nur die Schaltfläche zum Schließen angezeigt. Windows verfügt in der Standardtitelleiste über keine eigene Vollbildschaltfläche, daher hat `FullscreenButtonState` dort keine Wirkung.

### Fensterstil steuern (Windows)

Um den Stil der Titelleiste unter Windows zu steuern, verwenden Sie das Feld `ExStyle` in der Struktur `WebviewWindowOptions`:

Beispiel:

```go {title="main.go"}
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/w32"
)

func main() {
    app := application.New(application.Options{
        Name: "My Application",
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Windows: application.WindowsWindow{
            ExStyle: w32.WS_EX_TOOLWINDOW | w32.WS_EX_NOREDIRECTIONBITMAP | w32.WS_EX_TOPMOST,
        },
    })

    app.Run()
}
```

Andere Optionen, die den erweiterten Stil eines Fensters beeinflussen, werden durch diese Einstellung überschrieben:

- HiddenOnTaskbar
- AlwaysOnTop
- IgnoreMouseEvents
- BackgroundType
