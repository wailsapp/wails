---
title: "메뉴 API"
description: "메뉴 API 전체 레퍼런스"
slug: "reference/menu"
sourcePath: "reference/menu.md"
---

## 개요

메뉴 API는 애플리케이션 메뉴, 컨텍스트 메뉴 및 시스템 트레이 메뉴를 만들고 관리하는 메서드를 제공합니다.

**메뉴 유형:**

- **애플리케이션 메뉴** - 상단 메뉴 모음(파일, 편집 등)
- **컨텍스트 메뉴** - 마우스 오른쪽 버튼 클릭 메뉴
- **시스템 트레이 메뉴** - 시스템 트레이/알림 영역의 메뉴

## 메뉴 만들기

### NewMenu()

새 메뉴를 만듭니다.

```go
func (a *App) NewMenu() *Menu
```

**예:**

```go
menu := app.NewMenu()
```

## 메뉴 메서드

### Add()

메뉴에 메뉴 항목을 추가합니다.

```go
func (m *Menu) Add(label string) *MenuItem
```

**매개변수:**

- `label` - 메뉴 항목에 표시할 텍스트

**반환값:** 생성된 메뉴 항목

**예:**

```go
item := menu.Add("Open File")
item.OnClick(func(ctx *application.Context) {
    // Handle click
})
```

### AddSubmenu()

메뉴에 하위 메뉴를 추가합니다.

```go
func (m *Menu) AddSubmenu(label string) *Menu
```

**매개변수:**

- `label` - 하위 메뉴 레이블

**반환값:** 생성된 하위 메뉴

**예:**

```go
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("New")
fileMenu.Add("Open")
fileMenu.Add("Save")
```

### AddSeparator()

메뉴 항목 사이에 시각적인 구분선을 추가합니다.

```go
func (m *Menu) AddSeparator()
```

**예:**

```go
menu.Add("Copy")
menu.Add("Paste")
menu.AddSeparator()
menu.Add("Select All")
```

**권장 사항:** 구분선을 사용하여 관련 메뉴 항목을 그룹화하세요.

### AddCheckbox()

선택 여부를 전환할 수 있는 메뉴 항목을 추가합니다.

```go
func (m *Menu) AddCheckbox(label string, checked bool) *MenuItem
```

**매개변수:**

- `label` - 체크박스 레이블
- `checked` - 초기 선택 상태

**예:**

```go
darkMode := menu.AddCheckbox("Dark Mode", false)
darkMode.OnClick(func(ctx *application.Context) {
    isChecked := darkMode.Checked()
    // Toggle dark mode
})
```

### AddRadio()

라디오 메뉴 항목(상호 배타적 그룹)을 추가합니다.

```go
func (m *Menu) AddRadio(label string, checked bool) *MenuItem
```

**매개변수:**

- `label` - 라디오 버튼 레이블
- `checked` - 초기 선택 상태

**예:**

```go
// Create radio group for view modes
viewMenu := menu.AddSubmenu("View")
listView := viewMenu.AddRadio("List View", true)
gridView := viewMenu.AddRadio("Grid View", false)
treeView := viewMenu.AddRadio("Tree View", false)

listView.OnClick(func(ctx *application.Context) {
    setViewMode("list")
})
gridView.OnClick(func(ctx *application.Context) {
    setViewMode("grid")
})
```

### Update()

메뉴 항목에 적용한 변경 사항이 반영되도록 메뉴를 업데이트합니다.

```go
func (m *Menu) Update()
```

**예:**

```go
item.SetEnabled(false)
menu.Update()  // Must call to apply changes
```

**중요:** 메뉴 항목의 속성을 수정한 후에는 항상 `Update()`을 호출하세요.

## 메뉴 항목 메서드

### OnClick()

메뉴 항목의 클릭 핸들러를 등록합니다.

```go
func (mi *MenuItem) OnClick(callback func(ctx *application.Context)) *MenuItem
```

**매개변수:**

- `callback` - 항목을 클릭할 때 호출되는 함수

**반환값:** 메뉴 항목(메서드 체이닝용)

**예:**

```go
item.OnClick(func(ctx *application.Context) {
    fmt.Println("Menu item clicked")
    app.Logger.Info("User clicked menu item")
})
```

### SetLabel()

메뉴 항목의 레이블을 변경합니다.

```go
func (mi *MenuItem) SetLabel(label string) *MenuItem
```

**예:**

```go
item.SetLabel("Save As...")
menu.Update()
```

### SetEnabled()

메뉴 항목을 활성화하거나 비활성화합니다.

```go
func (mi *MenuItem) SetEnabled(enabled bool) *MenuItem
```

**예:**

```go
// Disable save when no document is open
saveItem.SetEnabled(hasOpenDocument)
menu.Update()
```

**일반적인 패턴:**

```go
// Update menu state based on application state
func updateMenuState() {
    saveItem.SetEnabled(hasUnsavedChanges)
    undoItem.SetEnabled(canUndo)
    redoItem.SetEnabled(canRedo)
    menu.Update()
}
```

### SetChecked()

체크박스/라디오 메뉴 항목의 선택 상태를 설정합니다.

```go
func (mi *MenuItem) SetChecked(checked bool) *MenuItem
```

**예:**

```go
darkModeItem.SetChecked(isDarkModeEnabled)
menu.Update()
```

### Checked()

현재 선택 상태를 반환합니다.

```go
func (mi *MenuItem) Checked() bool
```

**예:**

```go
if darkModeItem.Checked() {
    // Dark mode is enabled
}
```

### SetAccelerator()

메뉴 항목의 키보드 단축키를 설정합니다.

```go
func (mi *MenuItem) SetAccelerator(accelerator string) *MenuItem
```

**매개변수:**

- `accelerator` - 키보드 단축키(예: "Ctrl+S", "Cmd+Q")

**단축키 형식:**

- **보조 키:** `Ctrl`, `Cmd`, `Alt`, `Shift`
- **키:** `A-Z`, `0-9`, `F1-F12`, `Enter`, `Backspace` 등
- **플랫폼:** macOS에서는 `Cmd`을, Windows/Linux에서는 `Ctrl`을 사용하세요.

**예:**

```go
saveItem.SetAccelerator("Ctrl+S")
quitItem.SetAccelerator("Ctrl+Q")
newItem.SetAccelerator("Ctrl+N")
```

**플랫폼을 고려한 예:**

```go
import "runtime"

var quitShortcut string
if runtime.GOOS == "darwin" {
    quitShortcut = "Cmd+Q"
} else {
    quitShortcut = "Ctrl+Q"
}
quitItem.SetAccelerator(quitShortcut)
```

### SetTooltip()

메뉴 항목 위에 포인터를 올리면 표시되는 도구 설명을 설정합니다.

```go
func (mi *MenuItem) SetTooltip(tooltip string) *MenuItem
```

**예:**

```go
item.SetTooltip("Opens a file from disk")
```

### SetHidden()

메뉴 항목을 표시하거나 숨깁니다.

```go
func (mi *MenuItem) SetHidden(hidden bool) *MenuItem
```

**예:**

```go
// Hide debug menu in production
debugItem.SetHidden(!isDevelopment)
menu.Update()
```

## 애플리케이션 메뉴

### app.Menu.Set()

애플리케이션의 기본 메뉴 모음을 설정합니다.

```go
func (mm *MenuManager) Set(menu *Menu)
```

**예:**

```go
menu := app.NewMenu()

// File menu
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("New").SetAccelerator("Ctrl+N").OnClick(newFile)
fileMenu.Add("Open").SetAccelerator("Ctrl+O").OnClick(openFile)
fileMenu.Add("Save").SetAccelerator("Ctrl+S").OnClick(saveFile)
fileMenu.AddSeparator()
fileMenu.Add("Exit").SetAccelerator("Ctrl+Q").OnClick(func(ctx *application.Context) {
    app.Quit()
})

// Edit menu
editMenu := menu.AddSubmenu("Edit")
editMenu.Add("Undo").SetAccelerator("Ctrl+Z").OnClick(undo)
editMenu.Add("Redo").SetAccelerator("Ctrl+Y").OnClick(redo)
editMenu.AddSeparator()
editMenu.Add("Cut").SetAccelerator("Ctrl+X").OnClick(cut)
editMenu.Add("Copy").SetAccelerator("Ctrl+C").OnClick(copy)
editMenu.Add("Paste").SetAccelerator("Ctrl+V").OnClick(paste)

app.Menu.Set(menu)
```

**플랫폼 참고 사항:**

- **macOS:** 메뉴가 화면 상단의 메뉴 모음에 표시됩니다.
- **Windows/Linux:** 메뉴가 창 제목 표시줄에 표시됩니다.
- **macOS:** 애플리케이션 이름이 지정된 애플리케이션 메뉴를 자동으로 추가합니다.

## 컨텍스트 메뉴

### app.ContextMenu.New() / app.ContextMenu.Add()

관리자를 통해 `*ContextMenu`을 생성하고 이름을 지정해 등록합니다. `ContextMenuManager.Add`에는 `*Menu`가 **아닌** `*ContextMenu`을 전달하며, `app.RegisterContextMenu` 메서드는 **없습니다**.

```go
func (cm *ContextMenuManager) New() *ContextMenu
func (cm *ContextMenuManager) Add(name string, menu *ContextMenu)
func (cm *ContextMenuManager) Get(name string) (*ContextMenu, bool)
func (cm *ContextMenuManager) Remove(name string)
```

또는 패키지 수준의 `application.NewContextMenu(name string) *ContextMenu`을 사용하면 컨텍스트 메뉴를 한 번에 생성하고 등록할 수 있습니다.

**Go:**

```go
// Build the context menu via the manager
contextMenu := app.ContextMenu.New()
contextMenu.Add("Cut").OnClick(cut)
contextMenu.Add("Copy").OnClick(copy)
contextMenu.Add("Paste").OnClick(paste)
contextMenu.AddSeparator()
contextMenu.Add("Select All").OnClick(selectAll)

// Register under a name; HTML opts into it via the CSS custom property below.
app.ContextMenu.Add("editor", contextMenu)
```

**HTML / CSS:**

오른쪽 클릭 대상 또는 그 상위 요소에 `--custom-contextmenu` CSS 사용자 지정 속성이 해당 이름으로 설정되어 있으면 런타임이 등록된 컨텍스트 메뉴를 엽니다. 선택 사항인 `--custom-contextmenu-data` 속성은 `ctx.ContextMenuData()`을 통해 Go 콜백으로 전달됩니다. 브라우저의 기본 컨텍스트 메뉴가 표시되지 않게 하려면 `--default-contextmenu: hide` 또는 `auto`/`show`을 설정합니다.

```html
<!-- Trigger context menu on right-click -->
<div style="--custom-contextmenu: editor; --default-contextmenu: hide">
    Right-click here for context menu
</div>
```

`data-wails-context-menu="..."` 속성은 없습니다. 이 속성은 런타임에 연결된 적이 없습니다.

**동적 컨텍스트 메뉴:**

```go
// Update context menu based on selection
func updateContextMenu() {
    contextMenu := app.ContextMenu.New()

    if hasSelection {
        contextMenu.Add("Cut").OnClick(cut)
        contextMenu.Add("Copy").OnClick(copy)
    }

    contextMenu.Add("Paste").SetEnabled(hasClipboardContent).OnClick(paste)

    app.ContextMenu.Add("editor", contextMenu)
}
```

## 시스템 트레이 메뉴

### app.SystemTray.New()

새 시스템 트레이 아이콘을 생성합니다.

```go
func (sm *SystemTrayManager) New() *SystemTray
```

**예:**

```go
tray := app.SystemTray.New()
```

### SetIcon()

시스템 트레이 아이콘을 설정합니다.

```go
func (st *SystemTray) SetIcon(icon []byte) *SystemTray
```

**예:**

```go
iconData, _ := os.ReadFile("icon.png")
tray.SetIcon(iconData)
```

### SetMenu()

시스템 트레이의 메뉴를 설정합니다.

```go
func (st *SystemTray) SetMenu(menu *Menu) *SystemTray
```

**예:**

```go
trayMenu := app.NewMenu()
trayMenu.Add("Show Window").OnClick(func(ctx *application.Context) {
    window.Show()
    window.Focus()
})
trayMenu.Add("Settings").OnClick(openSettings)
trayMenu.AddSeparator()
trayMenu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})

tray.SetMenu(trayMenu)
```

### SetTooltip()

트레이 아이콘 위에 포인터를 올리면 표시되는 도구 설명을 설정합니다. 반환값은 없습니다.

```go
func (st *SystemTray) SetTooltip(tooltip string)
```

**예:**

```go
tray.SetTooltip("My Application - Running")
```

### OnClick()

트레이 아이콘의 왼쪽 클릭을 처리합니다.

```go
func (st *SystemTray) OnClick(callback func()) *SystemTray
```

**예:**

```go
tray.OnClick(func() {
    if window.IsVisible() {
        window.Hide()
    } else {
        window.Show()
        window.Focus()
    }
})
```

## 전체 예제

### 표준 애플리케이션 메뉴

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
)

func createMenu(app *application.App) *application.Menu {
    menu := app.NewMenu()

    // File menu
    fileMenu := menu.AddSubmenu("File")
    fileMenu.Add("New").
        SetAccelerator("Ctrl+N").
        OnClick(func(ctx *application.Context) {
            // Create new document
        })
    fileMenu.Add("Open").
        SetAccelerator("Ctrl+O").
        OnClick(func(ctx *application.Context) {
            // Open file dialog
        })
    fileMenu.Add("Save").
        SetAccelerator("Ctrl+S").
        OnClick(func(ctx *application.Context) {
            // Save document
        })
    fileMenu.AddSeparator()
    fileMenu.Add("Exit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })

    // Edit menu
    editMenu := menu.AddSubmenu("Edit")
    editMenu.Add("Undo").SetAccelerator("Ctrl+Z")
    editMenu.Add("Redo").SetAccelerator("Ctrl+Y")
    editMenu.AddSeparator()
    editMenu.Add("Cut").SetAccelerator("Ctrl+X")
    editMenu.Add("Copy").SetAccelerator("Ctrl+C")
    editMenu.Add("Paste").SetAccelerator("Ctrl+V")

    // View menu
    viewMenu := menu.AddSubmenu("View")
    darkMode := viewMenu.AddCheckbox("Dark Mode", false)
    darkMode.OnClick(func(ctx *application.Context) {
        // Toggle dark mode
        isChecked := darkMode.Checked()
        app.Logger.Info("Dark mode", "enabled", isChecked)
    })
    viewMenu.AddSeparator()
    viewMenu.AddRadio("List View", true)
    viewMenu.AddRadio("Grid View", false)
    viewMenu.AddRadio("Detail View", false)

    // Help menu
    helpMenu := menu.AddSubmenu("Help")
    helpMenu.Add("Documentation").OnClick(func(ctx *application.Context) {
        // Open docs
    })
    helpMenu.Add("About").OnClick(func(ctx *application.Context) {
        // Show about dialog
    })

    return menu
}

func main() {
    app := application.New(application.Options{
        Name: "Menu Demo",
    })

    menu := createMenu(app)
    app.Menu.Set(menu)

    window := app.Window.New()
    window.Show()

    app.Run()
}
```

### 시스템 트레이 애플리케이션

```go
func setupSystemTray(app *application.App, window application.Window) {
    // Create system tray
    tray := app.SystemTray.New()

    // Set icon
    iconData, _ := os.ReadFile("icon.png")
    tray.SetIcon(iconData)
    tray.SetTooltip("My App - Running")

    // Handle left-click on tray icon
    tray.OnClick(func() {
        if window.IsVisible() {
            window.Hide()
        } else {
            window.Show()
            window.Focus()
        }
    })

    // Create tray menu
    trayMenu := app.NewMenu()

    showItem := trayMenu.Add("Show Window")
    showItem.OnClick(func(ctx *application.Context) {
        window.Show()
        window.Focus()
    })

    trayMenu.AddSeparator()

    trayMenu.Add("Settings").OnClick(func(ctx *application.Context) {
        // Open settings window
    })

    trayMenu.Add("About").OnClick(func(ctx *application.Context) {
        // Show about dialog
    })

    trayMenu.AddSeparator()

    trayMenu.Add("Quit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })

    tray.SetMenu(trayMenu)
}
```

### 동적 메뉴 업데이트

```go
type Editor struct {
    app         *application.App
    menu        *application.Menu
    undoItem    *application.MenuItem
    redoItem    *application.MenuItem
    saveItem    *application.MenuItem
    undoStack   []string
    redoStack   []string
    hasChanges  bool
}

func (e *Editor) createMenu() {
    e.menu = e.app.NewMenu()

    fileMenu := e.menu.AddSubmenu("File")
    e.saveItem = fileMenu.Add("Save").SetAccelerator("Ctrl+S")
    e.saveItem.OnClick(func(ctx *application.Context) {
        e.save()
    })

    editMenu := e.menu.AddSubmenu("Edit")
    e.undoItem = editMenu.Add("Undo").SetAccelerator("Ctrl+Z")
    e.undoItem.OnClick(func(ctx *application.Context) {
        e.undo()
    })

    e.redoItem = editMenu.Add("Redo").SetAccelerator("Ctrl+Y")
    e.redoItem.OnClick(func(ctx *application.Context) {
        e.redo()
    })

    e.updateMenuState()
    e.app.Menu.Set(e.menu)
}

func (e *Editor) updateMenuState() {
    // Update menu items based on current state
    e.saveItem.SetEnabled(e.hasChanges)
    e.undoItem.SetEnabled(len(e.undoStack) > 0)
    e.redoItem.SetEnabled(len(e.redoStack) > 0)
    e.menu.Update()
}

func (e *Editor) onChange() {
    e.hasChanges = true
    e.updateMenuState()
}

func (e *Editor) save() {
    // Save logic
    e.hasChanges = false
    e.updateMenuState()
}
```

## 권장 사항

### ✅ 권장 사항

- **표준 단축키를 사용하세요** - 플랫폼 관례를 따르세요(예: 복사에는 Ctrl+C).
- **변경 후 Update()를 호출하세요** - 호출하지 않으면 변경 사항이 메뉴에 반영되지 않습니다.
- **관련 항목을 그룹화하세요** - 구분선을 사용해 메뉴 항목을 정리하세요.
- **사용할 수 없는 동작을 비활성화하세요** - 숨기지 말고 SetEnabled(false)로 비활성화하세요.
- **명확한 레이블을 사용하세요** - 간결하면서도 의미가 분명해야 합니다.
- **플랫폼 관례를 따르세요** - macOS와 Windows/Linux의 메뉴 패턴은 서로 다릅니다.

### ❌ 피해야 할 사항

- **Update() 호출을 잊지 마세요** - 가장 흔한 실수입니다.
- **메뉴를 너무 깊게 중첩하지 마세요** - 메뉴 깊이는 최대 2-3단계로 유지하세요.
- **모호한 레이블을 사용하지 마세요** - "문서 처리" 대신 "처리"처럼 쓰지 마세요.
- **지나치게 복잡하게 만들지 마세요** - 메뉴를 단순하고 목적에 집중하도록 구성하세요.
- **서로 다른 표현 체계를 혼용하지 마세요** - 이름과 구성을 일관되게 유지하세요.

## 플랫폼별 참고 사항

### macOS

- 애플리케이션 이름이 지정된 애플리케이션 메뉴가 자동으로 추가됩니다.
- 단축키에는 `Ctrl` 대신 `Cmd`을 사용하세요.
- 기본적으로 애플리케이션 메뉴에 "정보", "환경설정", "종료"가 포함됩니다.

### Windows/Linux

- 애플리케이션 메뉴가 자동으로 생성되지 않음
- 단축키에는 `Ctrl` 사용
- 일반적으로 파일 메뉴에 "종료" 배치
