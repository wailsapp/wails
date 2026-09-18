---
title: "Informações da tela"
description: "Obtenha informações sobre telas e monitores"
slug: "features/screens/info"
sourcePath: "features/screens/info.md"
---

## Informações da tela

O Wails fornece uma **API unificada para telas** que funciona em todas as plataformas. Obtenha informações sobre telas, detecte vários monitores, consulte propriedades da tela (tamanho, posição e DPI), identifique a tela principal e gerencie o escalonamento de DPI com código consistente.

## Início rápido

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

**É só isso!** Informações de tela multiplataforma.

## Como obter informações da tela

### Todas as telas

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

### Tela principal

```go
primary := app.Screen.GetPrimary()

fmt.Printf("Primary screen: %s\n", primary.Name)
fmt.Printf("Resolution: %dx%d\n", primary.Size.Width, primary.Size.Height)
fmt.Printf("Scale factor: %.2f\n", primary.ScaleFactor)
```

### Tela atual

Obtenha a tela que contém uma janela:

```go
screen, err := window.GetScreen()
if err != nil {
    // handle error
}
fmt.Printf("Window is on: %s\n", screen.Name)
```

### Tela por ID

```go
screen := app.Screen.GetByID("screen-id")
if screen != nil {
    fmt.Printf("Found screen: %s\n", screen.Name)
}
```

## Propriedades da tela

### Estrutura da tela

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

Leia os tamanhos em pixels por meio de `screen.Size.Width` / `screen.Size.Height` — não há campos `Width`/`Height` no nível superior.

### Pixels físicos e lógicos

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

**Fatores de escala comuns:**

- `1.0` - DPI padrão (96 DPI)
- `1.25` - Escala de 125% (120 DPI)
- `1.5` - Escala de 150% (144 DPI)
- `2.0` - Escala de 200% (192 DPI) - Retina
- `3.0` - Escala de 300% (288 DPI) - 4K/5K

## Posicionamento de janelas

### Centralizar na tela

```go
func centreOnScreen(window *application.WebviewWindow, screen *Screen) {
    windowWidth, windowHeight := window.Size()
    
    x := screen.X + (screen.Size.Width-windowWidth)/2
    y := screen.Y + (screen.Size.Height-windowHeight)/2
    
    window.SetPosition(x, y)
}
```

### Posicionar em uma tela específica

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

### Posicionar em relação à tela

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

## Suporte a vários monitores

### Detectar vários monitores

```go
func hasMultipleMonitors() bool {
    return len(app.Screen.GetAll()) > 1
}

func getMonitorCount() int {
    return len(app.Screen.GetAll())
}
```

### Listar todos os monitores

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

### Escolher o monitor

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

## Exemplos completos

### Gerenciador de janelas para vários monitores

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

### Detecção de alterações nas telas

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

### Dimensionamento de janelas com reconhecimento de DPI

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

### Visualizador do layout das telas

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

## Práticas recomendadas

### ✅ Faça

- **Verifique a quantidade de telas** - Ofereça suporte a um ou vários monitores
- **Use pixels lógicos** - O Wails gerencia o DPI automaticamente
- **Centralize as janelas** - Isso proporciona uma experiência do usuário melhor do que posições fixas
- **Valide as posições** - Garanta que as janelas estejam visíveis
- **Trate alterações nas telas** - Monitores podem ser adicionados ou removidos
- **Teste com diferentes configurações de DPI** - 100%, 125%, 150%, 200%

### ❌ Não faça

- **Não fixe posições no código** - Use as dimensões da tela
- **Não pressuponha que haja apenas a tela principal** - O usuário pode ter várias telas
- **Não ignore o fator de escala** - Ele é importante para o reconhecimento de DPI
- **Não posicione janelas fora da tela** - Valide as coordenadas
- **Não se esqueça das alterações nas telas** - Notebooks podem ser conectados e desconectados de uma estação de acoplamento
- **Não use pixels físicos** - Use pixels lógicos

## Diferenças entre plataformas

### macOS

- Telas Retina (fator de escala 2x)
- O uso de várias telas é comum
- Sistema de coordenadas: (0,0) no canto inferior esquerdo
- Os Spaces (áreas de trabalho virtuais) afetam o posicionamento

### Windows

- Vários níveis de escala de DPI (100%, 125%, 150%, 200%)
- O uso de várias telas é comum
- Sistema de coordenadas: (0,0) no canto superior esquerdo
- Reconhecimento de DPI por monitor

### Linux

- Varia conforme o ambiente de desktop
- Diferenças entre X11 e Wayland
- O suporte ao escalonamento de DPI varia
- Há suporte a várias telas

## Próximas etapas

@cards{cols="2"}
▣ Janelas
Saiba mais sobre o gerenciamento de janelas.

[Saiba mais →](/features/windows/basics/)

---
⚙ Opções de janela
Configure a aparência da janela.

[Saiba mais →](/features/windows/options/)

---
◆ Várias janelas
Padrões para várias janelas.

[Saiba mais →](/features/windows/multiple/)

---
🚀 Bindings
Chame funções Go a partir do JavaScript.

[Saiba mais →](/features/bindings/methods/)

@end

---

**Tem alguma dúvida?** Pergunte no [Discord](https://discord.gg/JDdSxwjhGf) ou consulte os [exemplos de telas](https://github.com/wailsapp/wails/tree/master/v3/examples).
