---
title: "Informations sur les écrans"
description: "Obtenir des informations sur les écrans et les moniteurs"
slug: "features/screens/info"
sourcePath: "features/screens/info.md"
---

## Informations sur les écrans

Wails fournit une **API unifiée pour les écrans** qui fonctionne sur toutes les plateformes. Obtenez des informations sur les écrans, détectez plusieurs moniteurs, interrogez les propriétés des écrans (taille, position, DPI), identifiez l’écran principal et gérez la mise à l’échelle selon les DPI avec un code cohérent.

## Démarrage rapide

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

**C’est tout !** Vous disposez désormais d’informations multiplateformes sur les écrans.

## Obtenir des informations sur les écrans

### Tous les écrans

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

### Écran principal

```go
primary := app.Screen.GetPrimary()

fmt.Printf("Primary screen: %s\n", primary.Name)
fmt.Printf("Resolution: %dx%d\n", primary.Size.Width, primary.Size.Height)
fmt.Printf("Scale factor: %.2f\n", primary.ScaleFactor)
```

### Écran actuel

Obtenez l’écran contenant une fenêtre :

```go
screen, err := window.GetScreen()
if err != nil {
    // handle error
}
fmt.Printf("Window is on: %s\n", screen.Name)
```

### Écran par ID

```go
screen := app.Screen.GetByID("screen-id")
if screen != nil {
    fmt.Printf("Found screen: %s\n", screen.Name)
}
```

## Propriétés de l’écran

### Structure de l’écran

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

Lisez les tailles en pixels via `screen.Size.Width` / `screen.Size.Height` — il n’existe aucun champ `Width`/`Height` de premier niveau.

### Pixels physiques et pixels logiques

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

**Facteurs d’échelle courants :**

- `1.0` - DPI standard (96 DPI)
- `1.25` - Mise à l’échelle de 125 % (120 DPI)
- `1.5` - Mise à l’échelle de 150 % (144 DPI)
- `2.0` - Mise à l’échelle de 200 % (192 DPI) - Retina
- `3.0` - Mise à l’échelle de 300 % (288 DPI) - 4K/5K

## Positionnement des fenêtres

### Centrer sur l’écran

```go
func centreOnScreen(window *application.WebviewWindow, screen *Screen) {
    windowWidth, windowHeight := window.Size()
    
    x := screen.X + (screen.Size.Width-windowWidth)/2
    y := screen.Y + (screen.Size.Height-windowHeight)/2
    
    window.SetPosition(x, y)
}
```

### Positionner sur un écran précis

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

### Positionner par rapport à l’écran

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

## Prise en charge de plusieurs moniteurs

### Détecter plusieurs moniteurs

```go
func hasMultipleMonitors() bool {
    return len(app.Screen.GetAll()) > 1
}

func getMonitorCount() int {
    return len(app.Screen.GetAll())
}
```

### Répertorier tous les moniteurs

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

### Choisir un moniteur

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

## Exemples complets

### Gestionnaire de fenêtres multimoniteur

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

### Détection des changements d’écran

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

### Dimensionnement des fenêtres tenant compte des DPI

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

### Visualiseur de la disposition des écrans

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

## Bonnes pratiques

### ✅ À faire

- **Vérifiez le nombre d’écrans** - Gérez les configurations à un ou plusieurs moniteurs
- **Utilisez des pixels logiques** - Wails gère automatiquement les DPI
- **Centrez les fenêtres** - L’expérience utilisateur est meilleure qu’avec des positions fixes
- **Validez les positions** - Assurez-vous que les fenêtres sont visibles
- **Gérez les changements d’écran** - Des moniteurs peuvent être ajoutés ou retirés
- **Testez avec différentes valeurs de DPI** - 100 %, 125 %, 150 %, 200 %

### ❌ À ne pas faire

- **Ne codez pas les positions en dur** - Utilisez les dimensions de l’écran
- **Ne supposez pas que l’écran principal est utilisé** - L’utilisateur peut en avoir plusieurs
- **N’ignorez pas le facteur d’échelle** - Il est important pour la prise en compte des DPI
- **Ne positionnez pas les fenêtres hors de l’écran** - Validez les coordonnées
- **N’oubliez pas les changements d’écran** - Les ordinateurs portables peuvent être connectés à une station d’accueil ou déconnectés de celle-ci
- **N’utilisez pas de pixels physiques** - Utilisez des pixels logiques

## Différences entre les plateformes

### macOS

- Écrans Retina (facteur d’échelle de 2x)
- Configurations à plusieurs écrans courantes
- Système de coordonnées : (0,0) en bas à gauche
- Les Spaces (bureaux virtuels) influent sur le positionnement

### Windows

- Différentes mises à l’échelle selon les DPI (100 %, 125 %, 150 %, 200 %)
- Configurations à plusieurs écrans courantes
- Système de coordonnées : (0,0) en haut à gauche
- Prise en compte des DPI propre à chaque moniteur

### Linux

- Varie selon l’environnement de bureau
- Différences entre X11 et Wayland
- La prise en charge de la mise à l’échelle selon les DPI varie
- Plusieurs écrans pris en charge

## Étapes suivantes

@cards{cols="2"}
▣ Fenêtres
Découvrez la gestion des fenêtres.

[En savoir plus →](/features/windows/basics/)

---
⚙ Options des fenêtres
Configurez l’apparence des fenêtres.

[En savoir plus →](/features/windows/options/)

---
◆ Fenêtres multiples
Découvrez les modèles à plusieurs fenêtres.

[En savoir plus →](/features/windows/multiple/)

---
🚀 Liaisons
Appelez des fonctions Go depuis JavaScript.

[En savoir plus →](/features/bindings/methods/)

@end

---

**Des questions ?** Posez-les sur [Discord](https://discord.gg/JDdSxwjhGf) ou consultez les [exemples relatifs aux écrans](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples).
