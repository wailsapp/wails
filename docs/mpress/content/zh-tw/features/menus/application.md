---
title: "應用程式選單"
description: "為桌面應用程式建立原生選單列"
slug: "features/menus/application"
sourcePath: "features/menus/application.md"
---

## 問題

專業桌面應用程式需要「檔案」、「編輯」、「顯示方式」和「輔助說明」等選單列。但選單在各平台上的運作方式不同：

- **macOS**：位於螢幕頂端的全域選單列
- **Windows**：位於視窗標題列中的選單列
- **Linux**：依桌面環境而異

手動建立符合各平台慣例的選單既繁瑣又容易出錯。

## Wails 解決方案

Wails 提供<strong>統一的 API</strong>，可自動建立各平台的原生選單。只需撰寫一次，即可在所有平台上獲得原生行為。

![macOS 上的 Wails 應用程式選單，包含標準項目、核取方塊項目、選項按鈕項目及子選單項目](/assets/screenshots/application-menu-macos.png)

在 macOS 上，應用程式選單會顯示於全域選單列中。此擷取畫面顯示由 Wails 選單 API 繪製的原生選單，其中包含停用的項目、核取方塊項目、選項按鈕項目及子選單項目。

## 快速開始

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

<strong>就這麼簡單！</strong>現在您已擁有包含標準項目的平台原生選單。`UseApplicationMenu`選項可確保 Windows 和 Linux 視窗無須額外程式碼即可顯示選單。

## 建立選單

### 建立基本選單

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

### 設定選單

**建議做法** — 使用`UseApplicationMenu`以確保跨平台一致性：

```go
// Set the application menu once
app.Menu.Set(menu)

// Create windows that inherit the menu on Windows/Linux
app.Window.NewWithOptions(application.WebviewWindowOptions{
    UseApplicationMenu: true,  // Window uses the app menu
})
```

此做法的效果如下：

- 在<strong>macOS</strong>上：選單會顯示於螢幕頂端（標準行為）
- 在<strong>Windows/Linux</strong>上：每個具有`UseApplicationMenu: true`的視窗都會顯示應用程式選單

**各平台的詳細資訊：**

@tabs{sync-key="platform"}
[macOS]
**全域選單列**（每個應用程式一個）：

```go
app.Menu.Set(menu)
```

選單會顯示於螢幕頂端，即使所有視窗都已關閉仍會保留。由於所有應用程式都使用全域選單，因此`UseApplicationMenu`選項在 macOS 上不會產生任何作用。

[Windows]
**各視窗的選單列**：

```go
// Option 1: Use application menu (recommended)
app.Menu.Set(menu)
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    UseApplicationMenu: true,
})

// Option 2: Set menu directly on window
window.SetMenu(menu)
```

每個視窗都能擁有自己的選單，或繼承應用程式選單。選單會顯示於視窗的標題列中。

[Linux]
**各視窗的選單列**（通常如此）：

```go
// Option 1: Use application menu (recommended)
app.Menu.Set(menu)
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    UseApplicationMenu: true,
})

// Option 2: Set menu directly on window
window.SetMenu(menu)
```

行為會依桌面環境而異。部分桌面環境（例如 Unity）支援全域選單。

@end

@note{type="tip" title="簡化跨平台選單"}
使用`UseApplicationMenu: true`後，就不需要如下所示的平台專用程式碼：

```go
// Old approach - no longer needed
if runtime.GOOS == "darwin" {
    app.Menu.Set(menu)
} else {
    window.SetMenu(menu)
}
```

@end

**各視窗的自訂選單：**

如果視窗需要不同於應用程式選單的選單，請直接為其設定：

```go
window.SetMenu(customMenu)  // Overrides UseApplicationMenu
```

## 選單角色

Wails 提供<strong>預先定義的選單角色</strong>，可自動建立符合各平台慣例的選單結構。

### 可用角色

| 角色 | 說明 | 平台注意事項 |
| --- | --- | --- |
| `AppMenu` | 包含「關於」、「偏好設定」和「結束」的應用程式選單 | **僅限 macOS** |
| `FileMenu` | 檔案操作（新增、開啟、儲存等） | 所有平台 |
| `EditMenu` | 文字編輯（還原、重做、剪下、複製、貼上） | 所有平台 |
| `WindowMenu` | 視窗管理（最小化、縮放等） | 所有平台 |
| `HelpMenu` | 說明與資訊 | 所有平台 |

### 使用角色

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

**您會獲得以下項目：**

@tabs{sync-key="platform"}
[macOS]
**AppMenu**（包含應用程式名稱）：

- 關於 [應用程式名稱]
- 偏好設定... (⌘,)
- ---
- 服務
- ---
- 隱藏 [應用程式名稱] (⌘H)
- 隱藏其他項目 (⌥⌘H)
- 全部顯示
- ---
- 結束 [應用程式名稱] (⌘Q)

**FileMenu**：

- 新增 (⌘N)
- 開啟... (⌘O)
- ---
- 關閉視窗 (⌘W)

**EditMenu**：

- 還原 (⌘Z)
- 重做 (⇧⌘Z)
- ---
- 剪下 (⌘X)
- 複製 (⌘C)
- 貼上 (⌘V)
- 全選 (⌘A)

**WindowMenu**：

- 最小化 (⌘M)
- 縮放
- ---
- 將所有視窗移至最前方

**HelpMenu**：

- [應用程式名稱] 說明

[Windows]
**FileMenu**：

- 新增 (Ctrl+N)
- 開啟... (Ctrl+O)
- ---
- 結束 (Alt+F4)

**EditMenu**：

- 復原 (Ctrl+Z)
- 重做 (Ctrl+Y)
- ---
- 剪下 (Ctrl+X)
- 複製 (Ctrl+C)
- 貼上 (Ctrl+V)
- 全選 (Ctrl+A)

**WindowMenu**：

- 最小化
- 最大化

**HelpMenu**：

- 關於 [應用程式名稱]

[Linux]
與 Windows 類似，但鍵盤快速鍵可能因桌面環境而異。

@end

### 自訂角色選單

`Menu.AddRole(role)`會傳回<strong>接收者</strong>選單（最上層選單），<strong>而不是</strong>該角色的子選單。若要將項目新增至該角色的子選單，請使用`FindByRole`查詢已插入的角色項目，並對其呼叫`GetSubmenu()`：

```go
menu.AddRole(application.FileMenu)

fileMenu := menu.FindByRole(application.FileMenu).GetSubmenu()
fileMenu.Add("Import...").OnClick(handleImport)
fileMenu.Add("Export...").OnClick(handleExport)
```

## 自訂選單

為應用程式特有的功能建立自己的選單：

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

**如需更多選單項目類型**，請參閱[選單參考](/features/menus/reference/)。

## 動態選單

根據應用程式狀態更新選單：

### 啟用／停用項目

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

@note{type="caution" title="一律呼叫 menu.Update()"}
變更選單狀態（啟用／停用、標籤、勾選狀態）後，**一律呼叫`menu.Update()`**。這在會重建選單的 Windows 上尤其重要。

如需詳細資訊，請參閱[選單參考](/features/menus/reference/#enabled-state)。

@end

### 變更標籤

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

### 重建選單

若要進行重大變更，請重建整個選單：

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

## 從選單控制視窗

選單項目可以控制視窗：

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

**取得作用中的視窗：**

```go
menuItem.OnClick(func(ctx *application.Context) {
    window := application.Get().Window.Current() // the window the menu was invoked from
    // Use window
})
```

## 各平台的注意事項

### macOS

**選單列行為：**

- 顯示於<strong>螢幕頂端</strong>（全域）
- 關閉所有視窗後仍會保留
- 第一個選單<strong>一律是應用程式選單</strong>
- 標準項目請使用`menu.AddRole(application.AppMenu)`

**標準位置：**

- **關於**：應用程式選單
- **偏好設定**：應用程式選單 (⌘,)
- **結束**：應用程式選單 (⌘Q)
- **說明**：說明選單

**範例：**

```go
if runtime.GOOS == "darwin" {
    menu.AddRole(application.AppMenu)  // Adds About, Preferences, Quit
    
    // Don't add Quit to File menu on macOS
    // Don't add About to Help menu on macOS
}
```

### Windows

**選單列行為：**

- 顯示於<strong>視窗標題列</strong>中
- 每個視窗都有自己的選單
- 沒有應用程式選單

**標準位置：**

- **結束**：檔案選單 (Alt+F4)
- **設定**：工具或編輯選單
- **關於**：說明選單

**範例：**

```go
if runtime.GOOS == "windows" {
    menu.AddRole(application.FileMenu) // Exit is added automatically
    menu.AddRole(application.HelpMenu) // About is added automatically
}
```

### Linux

**選單列行為：**

- 通常每個視窗各自擁有選單（與 Windows 類似）
- 某些桌面環境支援全域選單（Unity、安裝擴充功能的 GNOME）
- 外觀會因桌面環境而異

<strong>最佳實務：</strong>遵循 Windows 慣例，並在目標桌面環境上測試。

## 完整範例

以下是可用於正式環境的選單結構：

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

## 最佳實務

### ✅ 建議做法

- 標準選單（檔案、編輯等）應<strong>使用選單角色</strong>
- 選單結構應<strong>遵循平台慣例</strong>
- 為常用操作<strong>新增鍵盤快速鍵</strong>
- 變更選單狀態後，**呼叫 menu.Update()**
- **在所有平台上測試**——行為會因平台而異
- **讓選單層級保持精簡**——最多 2-3 層
- **使用清楚的標籤**——使用「儲存專案」，而不是「儲存」

### ❌ 請勿這樣做

- **請勿寫死各平台的快速鍵**——請使用 `CmdOrCtrl`
- **在 macOS 上，請勿將「結束」放在「檔案」選單中**——它應位於「應用程式」選單中
- **在 macOS 上，請勿將「關於」放在「說明」選單中**——它應位於「應用程式」選單中
- **請勿忘記呼叫 menu.Update()**——否則選單將無法正常運作
- **請勿巢狀設定過多層級**——使用者會迷失方向
- **請勿使用術語**——讓標籤保持淺顯易懂

## 後續步驟

@cards{cols="2"}
📖 選單參考
選單項目類型和屬性的完整參考資料。

[深入瞭解 →](/features/menus/reference/)

---
◆ 快顯選單
建立按一下滑鼠右鍵時顯示的快顯選單。

[深入瞭解 →](/features/menus/context/)

---
★ 系統匣選單
新增系統匣／選單列整合。

[深入瞭解 →](/features/menus/systray/)

---
📖 選單模式
常見的選單模式與最佳實務。

[深入瞭解 →](/guides/menus/)

@end

---

<strong>有疑問嗎？</strong>請在 [Discord](https://discord.gg/JDdSxwjhGf) 中提問，或查看[選單範例](https://github.com/wailsapp/wails/tree/master/v3/examples/menu)。
