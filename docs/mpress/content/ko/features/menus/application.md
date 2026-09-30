---
title: "애플리케이션 메뉴"
description: "데스크톱 애플리케이션용 네이티브 메뉴 막대 만들기"
slug: "features/menus/application"
sourcePath: "features/menus/application.md"
---

## 문제점

전문적인 데스크톱 애플리케이션에는 파일, 편집, 보기, 도움말 등의 메뉴 막대가 필요합니다. 하지만 메뉴는 플랫폼마다 다르게 동작합니다.

- **macOS**: 화면 상단의 전역 메뉴 막대
- **Windows**: 창 제목 표시줄의 메뉴 막대
- **Linux**: 데스크톱 환경에 따라 다름

플랫폼에 적합한 메뉴를 수동으로 만들려면 번거롭고 오류가 발생하기 쉽습니다.

## Wails 솔루션

Wails는 플랫폼 네이티브 메뉴를 자동으로 만드는 <strong>통합 API</strong>를 제공합니다. 한 번만 작성하면 모든 플랫폼에서 네이티브 동작을 구현할 수 있습니다.

![표준 항목, 체크박스 항목, 라디오 항목 및 하위 메뉴 항목이 포함된 macOS용 Wails 애플리케이션 메뉴](/assets/screenshots/application-menu-macos.png)

macOS에서는 애플리케이션 메뉴가 전역 메뉴 막대에 배치됩니다. 이 캡처에는 비활성화된 항목, 체크박스 항목, 라디오 항목 및 하위 메뉴 항목을 포함하여 Wails 메뉴 API가 렌더링한 네이티브 메뉴가 표시되어 있습니다.

## 빠른 시작

```go
package main

import (
    "runtime"
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "My App",
    })

    // Create menu
    menu := app.NewMenu()

    // Add standard menus (platform-appropriate)
    if runtime.GOOS == "darwin" {
        menu.AddRole(application.AppMenu)  // macOS only
    }
    menu.AddRole(application.FileMenu)
    menu.AddRole(application.EditMenu)
    menu.AddRole(application.WindowMenu)
    menu.AddRole(application.HelpMenu)

    // Set the application menu
    app.Menu.Set(menu)

    // Create window with UseApplicationMenu to inherit the menu on Windows/Linux
    app.Window.NewWithOptions(application.WebviewWindowOptions{
        UseApplicationMenu: true,
    })

    app.Run()
}
```

**이것으로 끝입니다!** 이제 표준 항목이 포함된 플랫폼 네이티브 메뉴를 사용할 수 있습니다. `UseApplicationMenu` 옵션을 사용하면 추가 코드 없이 Windows 및 Linux 창에 메뉴가 표시됩니다.

## 메뉴 만들기

### 기본 메뉴 만들기

```go
// Create a new menu
menu := app.NewMenu()

// Add a top-level menu
fileMenu := menu.AddSubmenu("File")

// Add menu items
fileMenu.Add("New").OnClick(func(ctx *application.Context) {
    // Handle New
})

fileMenu.Add("Open").OnClick(func(ctx *application.Context) {
    // Handle Open
})

fileMenu.AddSeparator()

fileMenu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})
```

### 메뉴 설정하기

**권장 방식** — 플랫폼 간 일관성을 위해 `UseApplicationMenu`를 사용하세요.

```go
// Set the application menu once
app.Menu.Set(menu)

// Create windows that inherit the menu on Windows/Linux
app.Window.NewWithOptions(application.WebviewWindowOptions{
    UseApplicationMenu: true,  // Window uses the app menu
})
```

이 방식은 다음과 같이 동작합니다.

- **macOS**: 메뉴가 화면 상단에 표시됩니다(표준 동작).
- **Windows/Linux**: `UseApplicationMenu: true`가 설정된 각 창에 애플리케이션 메뉴가 표시됩니다.

**플랫폼별 세부 정보:**

@tabs{sync-key="platform"}
[macOS]
**전역 메뉴 막대**(애플리케이션당 하나):

```go
app.Menu.Set(menu)
```

메뉴는 화면 상단에 표시되며 모든 창을 닫아도 유지됩니다. macOS에서는 모든 앱이 전역 메뉴를 사용하므로 `UseApplicationMenu` 옵션이 아무런 영향을 주지 않습니다.

[Windows]
**창별 메뉴 막대**:

```go
// Option 1: Use application menu (recommended)
app.Menu.Set(menu)
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    UseApplicationMenu: true,
})

// Option 2: Set menu directly on window
window.SetMenu(menu)
```

각 창에 자체 메뉴를 설정하거나 애플리케이션 메뉴를 상속하도록 할 수 있습니다. 메뉴는 창의 제목 표시줄에 나타납니다.

[Linux]
**창별 메뉴 막대**(일반적인 경우):

```go
// Option 1: Use application menu (recommended)
app.Menu.Set(menu)
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    UseApplicationMenu: true,
})

// Option 2: Set menu directly on window
window.SetMenu(menu)
```

동작은 데스크톱 환경에 따라 다릅니다. 일부 환경(예: Unity)은 전역 메뉴를 지원합니다.

@end

@note{type="tip" title="크로스 플랫폼 메뉴 간소화"}
`UseApplicationMenu: true`를 사용하면 다음과 같은 플랫폼별 코드가 필요하지 않습니다.

```go
// Old approach - no longer needed
if runtime.GOOS == "darwin" {
    app.Menu.Set(menu)
} else {
    window.SetMenu(menu)
}
```

@end

**창별 사용자 지정 메뉴:**

창에 애플리케이션 메뉴와 다른 메뉴가 필요한 경우 해당 창에 직접 설정하세요.

```go
window.SetMenu(customMenu)  // Overrides UseApplicationMenu
```

## 메뉴 역할

Wails는 플랫폼에 적합한 메뉴 구조를 자동으로 만드는 <strong>사전 정의된 메뉴 역할</strong>을 제공합니다.

### 사용 가능한 역할

| 역할 | 설명 | 플랫폼 참고 사항 |
| --- | --- | --- |
| `AppMenu` | 정보, 환경설정, 종료가 포함된 애플리케이션 메뉴 | **macOS 전용** |
| `FileMenu` | 파일 작업(새로 만들기, 열기, 저장 등) | 모든 플랫폼 |
| `EditMenu` | 텍스트 편집(실행 취소, 다시 실행, 잘라내기, 복사, 붙여넣기) | 모든 플랫폼 |
| `WindowMenu` | 창 관리(최소화, 확대/축소 등) | 모든 플랫폼 |
| `HelpMenu` | 도움말 및 정보 | 모든 플랫폼 |

### 역할 사용하기

```go
menu := app.NewMenu()

// macOS: Add application menu
if runtime.GOOS == "darwin" {
    menu.AddRole(application.AppMenu)
}

// All platforms: Add standard menus
menu.AddRole(application.FileMenu)
menu.AddRole(application.EditMenu)
menu.AddRole(application.WindowMenu)
menu.AddRole(application.HelpMenu)
```

**생성되는 항목:**

@tabs{sync-key="platform"}
[macOS]
**AppMenu**(앱 이름 포함):

- [App Name] 정보
- 환경설정... (⌘,)
- ---
- 서비스
- ---
- [App Name] 가리기 (⌘H)
- 다른 항목 가리기 (⌥⌘H)
- 모두 보기
- ---
- [App Name] 종료 (⌘Q)

**FileMenu**:

- 새로 만들기 (⌘N)
- 열기... (⌘O)
- ---
- 창 닫기 (⌘W)

**EditMenu**:

- 실행 취소 (⌘Z)
- 다시 실행 (⇧⌘Z)
- ---
- 잘라내기 (⌘X)
- 복사 (⌘C)
- 붙여넣기 (⌘V)
- 모두 선택 (⌘A)

**WindowMenu**:

- 최소화 (⌘M)
- 확대/축소
- ---
- 모두 앞으로 가져오기

**HelpMenu**:

- [App Name] 도움말

[Windows]
**FileMenu**:

- 새로 만들기 (Ctrl+N)
- 열기... (Ctrl+O)
- ---
- 종료 (Alt+F4)

**EditMenu**:

- 실행 취소 (Ctrl+Z)
- 다시 실행 (Ctrl+Y)
- ---
- 잘라내기 (Ctrl+X)
- 복사 (Ctrl+C)
- 붙여넣기 (Ctrl+V)
- 모두 선택 (Ctrl+A)

**WindowMenu**:

- 최소화
- 최대화

**HelpMenu**:

- [App Name] 정보

[Linux]
Windows와 유사하지만 키보드 단축키는 데스크톱 환경에 따라 다를 수 있습니다.

@end

### 역할 메뉴 사용자 지정

`Menu.AddRole(role)`은 역할의 하위 메뉴가 **아니라** **리시버** 메뉴(최상위 메뉴)를 반환합니다. 역할의 하위 메뉴에 항목을 추가하려면 `FindByRole`으로 삽입된 역할 항목을 찾은 다음 해당 항목에서 `GetSubmenu()`을 호출하세요:

```go
menu.AddRole(application.FileMenu)

fileMenu := menu.FindByRole(application.FileMenu).GetSubmenu()
fileMenu.Add("Import...").OnClick(handleImport)
fileMenu.Add("Export...").OnClick(handleExport)
```

## 사용자 지정 메뉴

애플리케이션별 기능을 위한 메뉴를 직접 만드세요:

```go
// Add a custom top-level menu
toolsMenu := menu.AddSubmenu("Tools")

// Add items
toolsMenu.Add("Settings").OnClick(func(ctx *application.Context) {
    showSettingsWindow()
})

toolsMenu.AddSeparator()

// Add checkbox
toolsMenu.AddCheckbox("Dark Mode", false).OnClick(func(ctx *application.Context) {
    isDark := ctx.ClickedMenuItem().Checked()
    setTheme(isDark)
})

// Add radio group
toolsMenu.AddRadio("Small", true).OnClick(handleFontSize)
toolsMenu.AddRadio("Medium", false).OnClick(handleFontSize)
toolsMenu.AddRadio("Large", false).OnClick(handleFontSize)

// Add submenu
advancedMenu := toolsMenu.AddSubmenu("Advanced")
advancedMenu.Add("Configure...").OnClick(showAdvancedSettings)
```

<strong>더 많은 메뉴 항목 유형</strong>은 [메뉴 레퍼런스](/features/menus/reference/)를 참조하세요.

## 동적 메뉴

애플리케이션 상태에 따라 메뉴를 업데이트하세요:

### 항목 활성화/비활성화

```go
var saveMenuItem *application.MenuItem

func createMenu() {
    menu := app.NewMenu()
    fileMenu := menu.AddSubmenu("File")
    
    saveMenuItem = fileMenu.Add("Save")
    saveMenuItem.SetEnabled(false)  // Initially disabled
    saveMenuItem.OnClick(handleSave)
    
    app.Menu.Set(menu)
}

func onDocumentChanged() {
    saveMenuItem.SetEnabled(hasUnsavedChanges())
    menu.Update()  // Important!
}
```

@note{type="caution" title="항상 menu.Update() 호출"}
메뉴 상태(활성화/비활성화, 레이블, 선택 상태)를 변경한 후에는 <strong>항상 `menu.Update()`</strong>을 호출하세요. 메뉴가 재구성되는 Windows에서는 특히 중요합니다.

자세한 내용은 [메뉴 레퍼런스](/features/menus/reference/#enabled-state)를 참조하세요.

@end

### 레이블 변경

```go
updateMenuItem := menu.Add("Check for Updates")

updateMenuItem.OnClick(func(ctx *application.Context) {
    updateMenuItem.SetLabel("Checking...")
    menu.Update()
    
    checkForUpdates()
    
    updateMenuItem.SetLabel("Check for Updates")
    menu.Update()
})
```

### 메뉴 다시 빌드

큰 변경 사항이 있는 경우 전체 메뉴를 다시 빌드하세요:

```go
func rebuildFileMenu() {
    menu := app.NewMenu()
    fileMenu := menu.AddSubmenu("File")
    
    fileMenu.Add("New").OnClick(handleNew)
    fileMenu.Add("Open").OnClick(handleOpen)
    
    // Add recent files dynamically
    if hasRecentFiles() {
        recentMenu := fileMenu.AddSubmenu("Open Recent")
        for _, file := range getRecentFiles() {
            filePath := file  // Capture for closure
            recentMenu.Add(filepath.Base(file)).OnClick(func(ctx *application.Context) {
                openFile(filePath)
            })
        }
        recentMenu.AddSeparator()
        recentMenu.Add("Clear Recent").OnClick(clearRecentFiles)
    }
    
    fileMenu.AddSeparator()
    fileMenu.Add("Quit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })
    
    app.Menu.Set(menu)
}
```

## 메뉴에서 창 제어

메뉴 항목으로 창을 제어할 수 있습니다:

```go
viewMenu := menu.AddSubmenu("View")

// Toggle fullscreen
viewMenu.Add("Toggle Fullscreen").OnClick(func(ctx *application.Context) {
    if window, ok := app.Window.GetByName("main"); ok {
        window.ToggleFullscreen()
    }
})

// Zoom controls
viewMenu.Add("Zoom In").SetAccelerator("CmdOrCtrl++").OnClick(func(ctx *application.Context) {
    // Increase zoom
})

viewMenu.Add("Zoom Out").SetAccelerator("CmdOrCtrl+-").OnClick(func(ctx *application.Context) {
    // Decrease zoom
})

viewMenu.Add("Reset Zoom").SetAccelerator("CmdOrCtrl+0").OnClick(func(ctx *application.Context) {
    // Reset zoom
})
```

**활성 창 가져오기:**

```go
menuItem.OnClick(func(ctx *application.Context) {
    window := application.Get().Window.Current() // the window the menu was invoked from
    // Use window
})
```

## 플랫폼별 고려 사항

### macOS

**메뉴 막대 동작:**

- <strong>화면 상단</strong>에 전역으로 표시됩니다
- 모든 창을 닫아도 유지됩니다
- 첫 번째 메뉴는 **항상 애플리케이션 메뉴입니다**
- 표준 항목에는 `menu.AddRole(application.AppMenu)`을 사용하세요

**표준 위치:**

- **정보**: 애플리케이션 메뉴
- **환경설정**: 애플리케이션 메뉴 (⌘,)
- **종료**: 애플리케이션 메뉴 (⌘Q)
- **도움말**: 도움말 메뉴

**예:**

```go
if runtime.GOOS == "darwin" {
    menu.AddRole(application.AppMenu)  // Adds About, Preferences, Quit
    
    // Don't add Quit to File menu on macOS
    // Don't add About to Help menu on macOS
}
```

### Windows

**메뉴 막대 동작:**

- <strong>창 제목 표시줄</strong>에 표시됩니다
- 각 창에는 자체 메뉴가 있습니다
- 애플리케이션 메뉴가 없습니다

**표준 위치:**

- **종료**: 파일 메뉴 (Alt+F4)
- **설정**: 도구 또는 편집 메뉴
- **정보**: 도움말 메뉴

**예:**

```go
if runtime.GOOS == "windows" {
    menu.AddRole(application.FileMenu) // Exit is added automatically
    menu.AddRole(application.HelpMenu) // About is added automatically
}
```

### Linux

**메뉴 막대 동작:**

- 일반적으로 창별로 표시됩니다(Windows와 유사)
- 일부 데스크톱 환경은 전역 메뉴를 지원합니다(Unity, 확장 기능을 사용하는 GNOME).
- 모양은 데스크톱 환경에 따라 달라집니다.

**권장 사항:** Windows 규칙을 따르고 대상 데스크톱 환경에서 테스트하세요.

## 전체 예제

다음은 프로덕션 환경에서 바로 사용할 수 있는 메뉴 구조입니다.

```go
package main

import (
    "runtime"
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "My Application",
    })

    // Create and set menu
    createMenu(app)

    // Create main window with UseApplicationMenu for cross-platform menu support
    app.Window.NewWithOptions(application.WebviewWindowOptions{
        UseApplicationMenu: true,
    })

    app.Run()
}

func createMenu(app *application.App) {
    menu := app.NewMenu()

    // Platform-specific application menu (macOS only)
    if runtime.GOOS == "darwin" {
        menu.AddRole(application.AppMenu)
    }

    // File menu — AddRole returns the receiver menu, not the role submenu.
    // To add items into the File submenu, look it up via FindByRole + GetSubmenu.
    menu.AddRole(application.FileMenu)
    fileMenu := menu.FindByRole(application.FileMenu).GetSubmenu()
    fileMenu.Add("Import...").SetAccelerator("CmdOrCtrl+I").OnClick(handleImport)
    fileMenu.Add("Export...").SetAccelerator("CmdOrCtrl+E").OnClick(handleExport)

    // Edit menu
    menu.AddRole(application.EditMenu)

    // View menu
    viewMenu := menu.AddSubmenu("View")
    viewMenu.Add("Toggle Fullscreen").SetAccelerator("F11").OnClick(toggleFullscreen)
    viewMenu.AddSeparator()
    viewMenu.AddCheckbox("Show Sidebar", true).OnClick(toggleSidebar)
    viewMenu.AddCheckbox("Show Toolbar", true).OnClick(toggleToolbar)

    // Tools menu
    toolsMenu := menu.AddSubmenu("Tools")
    
    // Settings location varies by platform
    if runtime.GOOS == "darwin" {
        // On macOS, Preferences is in Application menu (added by AppMenu role)
    } else {
        toolsMenu.Add("Settings").SetAccelerator("CmdOrCtrl+,").OnClick(showSettings)
    }
    
    toolsMenu.AddSeparator()
    toolsMenu.AddCheckbox("Dark Mode", false).OnClick(toggleDarkMode)

    // Window menu
    menu.AddRole(application.WindowMenu)

    // Help menu
    helpMenu := menu.AddRole(application.HelpMenu)
    helpMenu.Add("Documentation").OnClick(openDocumentation)
    
    // About location varies by platform
    if runtime.GOOS == "darwin" {
        // On macOS, About is in Application menu (added by AppMenu role)
    } else {
        helpMenu.AddSeparator()
        helpMenu.Add("About").OnClick(showAbout)
    }

    // Set the application menu
    app.Menu.Set(menu)
}

func handleImport(ctx *application.Context) {
    // Implementation
}

func handleExport(ctx *application.Context) {
    // Implementation
}

func toggleFullscreen(ctx *application.Context) {
    window := application.Get().Window.Current()
    window.ToggleFullscreen()
}

func toggleSidebar(ctx *application.Context) {
    // Implementation
}

func toggleToolbar(ctx *application.Context) {
    // Implementation
}

func showSettings(ctx *application.Context) {
    // Implementation
}

func toggleDarkMode(ctx *application.Context) {
    isDark := ctx.ClickedMenuItem().Checked()
    // Apply theme
}

func openDocumentation(ctx *application.Context) {
    // Open browser
}

func showAbout(ctx *application.Context) {
    // Show about dialog
}
```

## 권장 사항

### ✅ 해야 할 일

- 표준 메뉴(파일, 편집 등)에는 **메뉴 역할을 사용하세요**.
- 메뉴 구조는 **플랫폼 규칙을 따르세요**.
- 자주 사용하는 작업에 **키보드 단축키를 추가하세요**.
- 메뉴 상태를 변경한 후에는 **menu.Update()를 호출하세요**.
- 동작 방식이 다를 수 있으므로 **모든 플랫폼에서 테스트하세요**.
- **메뉴 계층을 얕게 유지하세요**. 최대 2-3단계까지만 사용하세요.
- **명확한 레이블을 사용하세요**. "저장"이 아니라 "프로젝트 저장"을 사용하세요.

### ❌ 하지 말아야 할 일

- **플랫폼별 단축키를 하드코딩하지 마세요**. `CmdOrCtrl`을 사용하세요.
- **macOS에서는 파일 메뉴에 종료를 넣지 마세요**. 종료는 애플리케이션 메뉴에 있습니다.
- **macOS에서는 도움말 메뉴에 정보를 넣지 마세요**. 정보는 애플리케이션 메뉴에 있습니다.
- **menu.Update() 호출을 잊지 마세요**. 그렇지 않으면 메뉴가 제대로 작동하지 않습니다.
- **메뉴를 너무 깊게 중첩하지 마세요**. 사용자가 길을 잃을 수 있습니다.
- **전문 용어를 사용하지 마세요**. 사용자가 이해하기 쉬운 레이블을 사용하세요.

## 다음 단계

@cards{cols="2"}
📖 메뉴 참조
메뉴 항목 유형과 속성에 관한 전체 참조입니다.

[자세히 알아보기 →](/features/menus/reference/)

---
◆ 컨텍스트 메뉴
마우스 오른쪽 버튼 클릭 컨텍스트 메뉴를 만드세요.

[자세히 알아보기 →](/features/menus/context/)

---
★ 시스템 트레이 메뉴
시스템 트레이/메뉴 막대 통합을 추가하세요.

[자세히 알아보기 →](/features/menus/systray/)

---
📖 메뉴 패턴
일반적인 메뉴 패턴과 권장 사항입니다.

[자세히 알아보기 →](/guides/menus/)

@end

---

**질문이 있으신가요?** [Discord](https://discord.gg/JDdSxwjhGf)에서 질문하거나 [메뉴 예제](https://github.com/wailsapp/wails/tree/master/v3/examples/menu)를 확인하세요.
