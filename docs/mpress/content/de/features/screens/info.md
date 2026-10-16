---
title: "Bildschirminformationen"
description: "Informationen zu Displays und Monitoren abrufen"
slug: "features/screens/info"
sourcePath: "features/screens/info.md"
---

## Bildschirminformationen

Wails bietet eine **einheitliche Bildschirm-API**, die auf allen Plattformen funktioniert. Rufen Sie Bildschirminformationen ab, erkennen Sie mehrere Monitore, fragen Sie Bildschirmeigenschaften (Größe, Position, DPI) ab, ermitteln Sie das primäre Display und handhaben Sie die DPI-Skalierung mit einheitlichem Code.

## Schnellstart

```go
// Get all screens
screens := app.Screen.GetAll()

for _, screen := range screens {
    fmt.Printf("Screen: %s (%dx%d)\n", 
        screen.Name, screen.Size.Width, screen.Size.Height)
}

// Get primary screen
primary := app.Screen.GetPrimary()
fmt.Printf("Primary: %s\n", primary.Name)
```

**Das war's!** Plattformübergreifende Bildschirminformationen.

## Bildschirminformationen abrufen

### Alle Bildschirme

```go
screens := app.Screen.GetAll()

for _, screen := range screens {
    fmt.Printf("ID: %s\n", screen.ID)
    fmt.Printf("Name: %s\n", screen.Name)
    fmt.Printf("Size: %dx%d\n", screen.Size.Width, screen.Size.Height)
    fmt.Printf("Position: %d,%d\n", screen.X, screen.Y)
    fmt.Printf("Scale: %.2f\n", screen.ScaleFactor)
    fmt.Printf("Primary: %v\n", screen.IsPrimary)
    fmt.Println("---")
}
```

### Primärer Bildschirm

```go
primary := app.Screen.GetPrimary()

fmt.Printf("Primary screen: %s\n", primary.Name)
fmt.Printf("Resolution: %dx%d\n", primary.Size.Width, primary.Size.Height)
fmt.Printf("Scale factor: %.2f\n", primary.ScaleFactor)
```

### Aktueller Bildschirm

Rufen Sie den Bildschirm ab, auf dem sich ein Fenster befindet:

```go
screen, err := window.GetScreen()
if err != nil {
    // handle error
}
fmt.Printf("Window is on: %s\n", screen.Name)
```

### Bildschirm nach ID

```go
screen := app.Screen.GetByID("screen-id")
if screen != nil {
    fmt.Printf("Found screen: %s\n", screen.Name)
}
```

## Bildschirmeigenschaften

### Bildschirmstruktur

```go
type Screen struct {
    ID               string  // Unique identifier
    Name             string  // Display name
    ScaleFactor      float32 // DPI scale (1.0, 1.5, 2.0, etc.)
    X, Y             int     // Position (top-left, logical pixels)
    Size             Size    // Logical size (Width, Height)
    Bounds           Rect    // Logical bounds (X, Y, Width, Height)
    PhysicalBounds   Rect    // Physical bounds (scaled)
    WorkArea         Rect    // Logical work area excluding taskbars/menubars
    PhysicalWorkArea Rect    // Physical work area
    IsPrimary        bool    // Is this the primary screen?
    Rotation         float32 // Screen rotation in degrees (e.g. 0, 90, 180, 270)
}

// Size and Rect are simple value types:
type Size struct{ Width, Height int }
type Rect struct{ X, Y, Width, Height int }
```

Lesen Sie Pixelgrößen über `screen.Size.Width` / `screen.Size.Height` aus – es gibt keine übergeordneten Felder `Width`/`Height`.

### Physische und logische Pixel

```go
screen := app.Screen.GetPrimary()

// Logical pixels (what you use)
logicalWidth := screen.Size.Width
logicalHeight := screen.Size.Height

// Physical pixels (actual display)
physicalWidth := int(float32(screen.Size.Width) * screen.ScaleFactor)
physicalHeight := int(float32(screen.Size.Height) * screen.ScaleFactor)

fmt.Printf("Logical: %dx%d\n", logicalWidth, logicalHeight)
fmt.Printf("Physical: %dx%d\n", physicalWidth, physicalHeight)
fmt.Printf("Scale: %.2f\n", screen.ScaleFactor)
```

**Übliche Skalierungsfaktoren:**

- `1.0` – Standard-DPI (96 DPI)
- `1.25` – 125 % Skalierung (120 DPI)
- `1.5` – 150 % Skalierung (144 DPI)
- `2.0` – 200 % Skalierung (192 DPI) – Retina
- `3.0` – 300 % Skalierung (288 DPI) – 4K/5K

## Fensterpositionierung

### Auf dem Bildschirm zentrieren

```go
func centreOnScreen(window *application.WebviewWindow, screen *Screen) {
    windowWidth, windowHeight := window.Size()
    
    x := screen.X + (screen.Size.Width-windowWidth)/2
    y := screen.Y + (screen.Size.Height-windowHeight)/2
    
    window.SetPosition(x, y)
}
```

### Auf einem bestimmten Bildschirm positionieren

```go
func moveToScreen(window *application.WebviewWindow, screenIndex int) {
    screens := app.Screen.GetAll()
    
    if screenIndex < 0 || screenIndex >= len(screens) {
        return
    }
    
    screen := screens[screenIndex]
    
    // Centre on target screen
    centreOnScreen(window, screen)
}
```

### Relativ zum Bildschirm positionieren

```go
// Top-left corner
func positionTopLeft(window *application.WebviewWindow, screen *Screen) {
    window.SetPosition(screen.X+10, screen.Y+10)
}

// Top-right corner
func positionTopRight(window *application.WebviewWindow, screen *Screen) {
    windowWidth, _ := window.Size()
    window.SetPosition(screen.X+screen.Size.Width-windowWidth-10, screen.Y+10)
}

// Bottom-right corner
func positionBottomRight(window *application.WebviewWindow, screen *Screen) {
    windowWidth, windowHeight := window.Size()
    window.SetPosition(
        screen.X+screen.Size.Width-windowWidth-10,
        screen.Y+screen.Size.Height-windowHeight-10,
    )
}
```

## Unterstützung mehrerer Monitore

### Mehrere Monitore erkennen

```go
func hasMultipleMonitors() bool {
    return len(app.Screen.GetAll()) > 1
}

func getMonitorCount() int {
    return len(app.Screen.GetAll())
}
```

### Alle Monitore auflisten

```go
func listMonitors() {
    screens := app.Screen.GetAll()
    
    fmt.Printf("Found %d monitor(s):\n", len(screens))
    
    for i, screen := range screens {
        primary := ""
        if screen.IsPrimary {
            primary = " (Primary)"
        }
        
        fmt.Printf("%d. %s%s\n", i+1, screen.Name, primary)
        fmt.Printf("   Resolution: %dx%d\n", screen.Size.Width, screen.Size.Height)
        fmt.Printf("   Position: %d,%d\n", screen.X, screen.Y)
        fmt.Printf("   Scale: %.2fx\n", screen.ScaleFactor)
    }
}
```

### Monitor auswählen

```go
func chooseMonitor() (*Screen, error) {
    screens := app.Screen.GetAll()
    
    if len(screens) == 1 {
        return screens[0], nil
    }
    
    // Show dialog to choose
    var options []string
    for i, screen := range screens {
        primary := ""
        if screen.IsPrimary {
            primary = " (Primary)"
        }
        options = append(options, 
            fmt.Sprintf("%d. %s%s - %dx%d", 
                i+1, screen.Name, primary, screen.Size.Width, screen.Size.Height))
    }
    
    // Use dialog to select
    // (Implementation depends on your dialog system)
    
    return screens[0], nil
}
```

## Vollständige Beispiele

### Fenstermanager für mehrere Monitore

```go
type MultiMonitorManager struct {
    app     *application.App
    windows map[int]*application.WebviewWindow
}

func NewMultiMonitorManager(app *application.App) *MultiMonitorManager {
    return &MultiMonitorManager{
        app:     app,
        windows: make(map[int]*application.WebviewWindow),
    }
}

func (m *MultiMonitorManager) CreateWindowOnScreen(screenIndex int) error {
    screens := m.app.Screen.GetAll()
    
    if screenIndex < 0 || screenIndex >= len(screens) {
        return errors.New("invalid screen index")
    }
    
    screen := screens[screenIndex]
    
    // Create window
    window := m.app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:  fmt.Sprintf("Window on %s", screen.Name),
        Width:  800,
        Height: 600,
    })
    
    // Centre on screen
    x := screen.X + (screen.Size.Width-800)/2
    y := screen.Y + (screen.Size.Height-600)/2
    window.SetPosition(x, y)
    
    window.Show()
    
    m.windows[screenIndex] = window
    return nil
}

func (m *MultiMonitorManager) CreateWindowOnEachScreen() {
    screens := m.app.Screen.GetAll()
    
    for i := range screens {
        m.CreateWindowOnScreen(i)
    }
}
```

### Bildschirmänderungen erkennen

```go
type ScreenMonitor struct {
    app           *application.App
    lastScreens   []*Screen
    changeHandler func([]*Screen)
}

func NewScreenMonitor(app *application.App) *ScreenMonitor {
    return &ScreenMonitor{
        app:         app,
        lastScreens: app.Screen.GetAll(),
    }
}

func (sm *ScreenMonitor) OnScreenChange(handler func([]*Screen)) {
    sm.changeHandler = handler
}

func (sm *ScreenMonitor) Start() {
    ticker := time.NewTicker(2 * time.Second)
    
    go func() {
        for range ticker.C {
            sm.checkScreens()
        }
    }()
}

func (sm *ScreenMonitor) checkScreens() {
    current := sm.app.Screen.GetAll()
    
    if len(current) != len(sm.lastScreens) {
        sm.lastScreens = current
        if sm.changeHandler != nil {
            sm.changeHandler(current)
        }
    }
}
```

### DPI-gerechte Fenstergröße

```go
func createDPIAwareWindow(screen *Screen) *application.WebviewWindow {
    // Base size at 1.0 scale
    baseWidth := 800
    baseHeight := 600
    
    // Adjust for DPI
    width := int(float32(baseWidth) * screen.ScaleFactor)
    height := int(float32(baseHeight) * screen.ScaleFactor)
    
    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:  "DPI-Aware Window",
        Width:  width,
        Height: height,
    })
    
    // Centre on screen
    x := screen.X + (screen.Size.Width-width)/2
    y := screen.Y + (screen.Size.Height-height)/2
    window.SetPosition(x, y)
    
    return window
}
```

### Visualisierung der Bildschirmanordnung

```go
func visualiseScreenLayout() string {
    screens := app.Screen.GetAll()
    
    var layout strings.Builder
    layout.WriteString("Screen Layout:\n\n")
    
    for i, screen := range screens {
        primary := ""
        if screen.IsPrimary {
            primary = " [PRIMARY]"
        }
        
        layout.WriteString(fmt.Sprintf("Screen %d: %s%s\n", i+1, screen.Name, primary))
        layout.WriteString(fmt.Sprintf("  Position: (%d, %d)\n", screen.X, screen.Y))
        layout.WriteString(fmt.Sprintf("  Size: %dx%d\n", screen.Size.Width, screen.Size.Height))
        layout.WriteString(fmt.Sprintf("  Scale: %.2fx\n", screen.ScaleFactor))
        layout.WriteString(fmt.Sprintf("  Physical: %dx%d\n", 
            int(float32(screen.Size.Width)*screen.ScaleFactor),
            int(float32(screen.Size.Height)*screen.ScaleFactor)))
        layout.WriteString("\n")
    }
    
    return layout.String()
}
```

## Bewährte Vorgehensweisen

### ✅ Empfohlen

- **Anzahl der Bildschirme prüfen** – Einzelne und mehrere Monitore berücksichtigen
- **Logische Pixel verwenden** – Wails handhabt DPI automatisch
- **Fenster zentrieren** – Bessere Benutzerfreundlichkeit als feste Positionen
- **Positionen validieren** – Sicherstellen, dass Fenster sichtbar sind
- **Bildschirmänderungen behandeln** – Monitore können hinzugefügt oder entfernt werden
- **Mit unterschiedlichen DPI-Einstellungen testen** – 100 %, 125 %, 150 %, 200 %

### ❌ Nicht empfohlen

- **Positionen nicht fest codieren** – Bildschirmabmessungen verwenden
- **Nicht vom primären Bildschirm ausgehen** – Benutzer haben möglicherweise mehrere Bildschirme
- **Skalierungsfaktor nicht ignorieren** – Wichtig für DPI-Bewusstsein
- **Fenster nicht außerhalb des Bildschirms positionieren** – Koordinaten validieren
- **Bildschirmänderungen nicht vergessen** – Laptops werden an Dockingstationen angeschlossen und davon getrennt
- **Keine physischen Pixel verwenden** – Logische Pixel verwenden

## Plattformunterschiede

### macOS

- Retina-Displays (2-facher Skalierungsfaktor)
- Mehrere Displays sind üblich
- Koordinatensystem: (0,0) unten links
- Spaces (virtuelle Desktops) beeinflussen die Positionierung

### Windows

- Verschiedene DPI-Skalierungen (100 %, 125 %, 150 %, 200 %)
- Mehrere Displays sind üblich
- Koordinatensystem: (0,0) oben links
- DPI-Bewusstsein pro Monitor

### Linux

- Je nach Desktop-Umgebung unterschiedlich
- Unterschiede zwischen X11 und Wayland
- Unterstützung für DPI-Skalierung variiert
- Mehrere Displays werden unterstützt

## Nächste Schritte

@cards{cols="2"}
▣ Fenster
Erfahren Sie mehr über die Fensterverwaltung.

[Mehr erfahren →](/features/windows/basics/)

---
⚙ Fensteroptionen
Konfigurieren Sie das Erscheinungsbild des Fensters.

[Mehr erfahren →](/features/windows/options/)

---
◆ Mehrere Fenster
Muster für Anwendungen mit mehreren Fenstern.

[Mehr erfahren →](/features/windows/multiple/)

---
🚀 Bindings
Rufen Sie Go-Funktionen aus JavaScript auf.

[Mehr erfahren →](/features/bindings/methods/)

@end

---

**Fragen?** Stellen Sie sie auf [Discord](https://discord.gg/JDdSxwjhGf) oder sehen Sie sich die [Bildschirmbeispiele](https://github.com/wailsapp/wails/tree/master/v3/examples) an.
