---
title: "Dock \u0026 Taskleiste"
description: "Sichtbarkeit des Dock-Symbols verwalten und Badges unter macOS und Windows anzeigen"
slug: "features/platform/dock"
sourcePath: "features/platform/dock.md"
---

## Einführung

Wails bietet einen plattformübergreifenden Dock-Dienst für Desktopanwendungen. Mit diesem Dienst können Sie:

- das Anwendungssymbol im macOS-Dock aus- und einblenden
- Badges auf der Anwendungskachel oder dem Dock-/Taskleistensymbol anzeigen (macOS und Windows)

## Grundlegende Verwendung

### Dienst erstellen

Initialisieren Sie zunächst den Dock-Dienst:

```go
import "github.com/wailsapp/wails/v3/pkg/application"
import "github.com/wailsapp/wails/v3/pkg/services/dock"

// Create a new Dock service
dockService := dock.New()

// Register the service with the application
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(dockService),
    },
})
```

### Dienst mit benutzerdefinierten Badge-Optionen erstellen (nur Windows)

Unter Windows können Sie das Erscheinungsbild des Badges mit verschiedenen Optionen anpassen:

```go
import "github.com/wailsapp/wails/v3/pkg/application"
import "github.com/wailsapp/wails/v3/pkg/services/dock"
import "image/color"

// Create a dock service with custom badge options
options := dock.BadgeOptions{
    TextColour:       color.RGBA{255, 255, 255, 255}, // White text
    BackgroundColour: color.RGBA{0, 0, 255, 255},     // Blue background
    FontName:         "consolab.ttf",                 // Bold Consolas font
    FontSize:         20,                             // Font size for single character
    SmallFontSize:    14,                             // Font size for multiple characters
}

dockService := dock.NewWithOptions(options)

// Register the service with the application
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(dockService),
    },
})
```

## Dock-Operationen

### Anwendungssymbol im Dock ausblenden

Blenden Sie das Anwendungssymbol im macOS-Dock aus:

```go
// Hide the app icon
dockService.HideAppIcon()
```

### Anwendungssymbol im Dock einblenden

Blenden Sie das Anwendungssymbol im macOS-Dock ein:

```go
// Show the app icon
dockService.ShowAppIcon()
```

## Badge-Operationen

### Badge festlegen

Legen Sie ein Badge auf der Anwendungskachel bzw. dem Dock-Symbol fest:

```go
// Set a default badge
dockService.SetBadge("")

// Set a numeric badge
dockService.SetBadge("3")

// Set a text badge
dockService.SetBadge("New")
```

### Benutzerdefiniertes Badge festlegen (nur Windows)

Legen Sie ein Badge mit einmalig angewendeten Optionen fest:

```go
options := dock.BadgeOptions{
    BackgroundColour: color.RGBA{0, 255, 255, 255},
    FontName:         "arialb.ttf", // System font
    FontSize:         16,
    SmallFontSize:    10,
    TextColour:       color.RGBA{0, 0, 0, 255},
}

// Set a default badge
dockService.SetCustomBadge("", options)

// Set a numeric badge
dockService.SetCustomBadge("3", options)

// Set a text badge
dockService.SetCustomBadge("New", options)
```

### Badge entfernen

Entfernen Sie das Badge vom Anwendungssymbol:

```go
dockService.RemoveBadge()
```

### Festgelegtes Badge abrufen

```go
dockService.GetBadge()
```

## Plattformspezifische Hinweise

@tabs
[macOS]
Unter macOS:

- Das Dock-Symbol kann **ausgeblendet** und **eingeblendet** werden
- Badges werden direkt auf dem Dock-Symbol angezeigt
- Badge-Optionen können **nicht angepasst werden** (an `NewWithOptions`/`SetCustomBadge` übergebene Optionen werden ignoriert)
- Es wird das standardmäßige Badge-Design des macOS-Docks verwendet, das sich automatisch an das Erscheinungsbild anpasst
- Das System behandelt zu lange Beschriftungen
- Bei einer leeren Beschriftung wird das Standard-Badge „●“ angezeigt

[Windows]
Unter Windows:

- Das Aus- und Einblenden des Taskleistensymbols wird von diesem Dienst derzeit nicht unterstützt
- Badges werden als überlagertes Symbol in der Taskleiste angezeigt
- Badges unterstützen Textwerte
- Das Erscheinungsbild des Badges kann über `BadgeOptions` angepasst werden
- Die Anwendung muss über ein Fenster verfügen, damit Badges angezeigt werden können
- Für mehrstellige Beschriftungen wird automatisch eine kleinere Schriftgröße verwendet
- Zu lange Beschriftungen werden nicht behandelt
- Anpassungsoptionen:
  - **TextColour**: Textfarbe (Standard: Weiß)
  - **BackgroundColour**: Hintergrundfarbe des Badges (Standard: Rot)
  - **FontName**: Name der Schriftdatei (Standard: „segoeuib.ttf“)
  - **FontSize**: Schriftgröße für ein einzelnes Zeichen (Standard: 18)
  - **SmallFontSize**: Schriftgröße für mehrere Zeichen (Standard: 14)


[Linux]
Unter Linux:

- Die Sichtbarkeitssteuerung des Dock-Symbols und die Badge-Funktion sind nicht verfügbar

@end

## Bewährte Verfahren

1. **Beim Ausblenden des Dock-Symbols (macOS):**
  - Stellen Sie sicher, dass Benutzer weiterhin auf Ihre Anwendung zugreifen können (z. B. über den [Infobereich](/features/menus/systray/))
  - Fügen Sie Ihrer alternativen Benutzeroberfläche eine Option zum Beenden hinzu
  - Die Anwendung erscheint nicht im Command+Tab-Umschalter
  - Geöffnete Fenster bleiben sichtbar und funktionsfähig
  - Das Schließen aller Fenster beendet die Anwendung möglicherweise nicht (das Verhalten variiert unter macOS)
  - Benutzern steht die übliche Möglichkeit zum Beenden über einen Rechtsklick im Dock nicht mehr zur Verfügung


2. **Verwenden Sie Badges sparsam:**
  - Zu häufige Badge-Aktualisierungen können Benutzer ablenken
  - Verwenden Sie Badges nur für wichtige Benachrichtigungen


3. **Halten Sie Badge-Texte kurz:**
  - Numerische Badges sind am wirkungsvollsten
  - Unter macOS sollten Text-Badges kurz sein


4. **Bei der Anpassung von Badges unter Windows:**
  - Achten Sie auf einen hohen Kontrast zwischen Text- und Hintergrundfarbe
  - Testen Sie unterschiedliche Textlängen, da die Schriftgröße mit zunehmender Länge abnimmt
  - Verwenden Sie gängige Systemschriftarten, um deren Verfügbarkeit sicherzustellen


## API-Referenz

### Dienstverwaltung

| Methode | Beschreibung |
| --- | --- |
| `New()` | Erstellt einen neuen Dock-Dienst |
| `NewWithOptions(options BadgeOptions)` | Erstellt einen neuen Dock-Dienst mit benutzerdefinierten Badge-Optionen (nur Windows; unter macOS und Linux werden die Optionen ignoriert) |

### Dock-Operationen

| Methode | Beschreibung |
| --- | --- |
| `HideAppIcon()` | Blendet das App-Symbol im macOS-Dock aus (nur macOS) |
| `ShowAppIcon()` | Zeigt das App-Symbol im macOS-Dock an (nur macOS) |

### Badge-Operationen

| Methode | Beschreibung |
| --- | --- |
| `SetBadge(label string) error` | Setzt ein Badge mit der angegebenen Beschriftung |
| `SetCustomBadge(label string, options BadgeOptions) error` | Setzt ein Badge mit der angegebenen Beschriftung und benutzerdefinierten Darstellungsoptionen (nur Windows) |
| `RemoveBadge() error` | Entfernt das Badge vom Anwendungssymbol |
| `GetBadge() *string` | Ruft das aktuelle Badge ab |

### Strukturen und Typen

```go
// Options for customizing badge appearance (Windows only)
type BadgeOptions struct {
    TextColour       color.RGBA  // Color of the badge text
    BackgroundColour color.RGBA  // Color of the badge background
    FontName         string      // Font file name (e.g., "segoeuib.ttf")
    FontSize         int         // Font size for single character
    SmallFontSize    int         // Font size for multiple characters
}
```
