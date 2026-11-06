---
title: "管理器 API"
description: "透過功能明確的管理器介面提供條理分明的 API 結構"
slug: "concepts/manager-api"
sourcePath: "concepts/manager-api.md"
---

Wails v3 管理器 API 將功能明確的管理器結構分組置於`*application.App`的公開欄位下，讓您能以條理分明且易於探索的方式存取應用程式功能。Wails 3與 v2 完全分離——不提供逐次呼叫的包裝層來維持與舊式`app.NewWebviewWindow(...)`風格 API 的相容性，因此必須透過下列管理器來操作應用程式。

## 概覽

管理器 API 將應用程式功能劃分為十二個明確的領域（一個記錄器及十一個管理器）：

- **`app.Window`**－視窗建立、管理及回呼
- **`app.ContextMenu`**－內容功能表註冊與管理\
- **`app.KeyBinding`**－全域按鍵綁定管理
- **`app.Browser`**－瀏覽器整合（開啟 URL 與檔案）
- **`app.Env`**－環境資訊與系統狀態
- **`app.Dialog`**－檔案與訊息對話方塊操作
- **`app.Event`**－自訂事件處理與應用程式事件
- **`app.Menu`**－應用程式功能表管理
- **`app.Screen`**－螢幕管理與座標轉換
- **`app.Clipboard`**－剪貼簿文字操作
- **`app.SystemTray`**－系統匣圖示建立與管理
- **`app.Autostart`**－註冊應用程式，使其在使用者登入時啟動

## 優點

- **更容易探索**－IDE 自動完成功能會顯示條理分明的 API 介面
- **改善程式碼組織**－將相關方法歸為一組
- **提升可維護性**－各管理器各司其職
- **便於日後擴充**－更容易為特定領域新增功能

## 用法

管理器 API 以條理分明的方式提供對所有應用程式功能的存取：

```go
// Events and custom event handling
app.Event.Emit("custom", data)
app.Event.On("custom", func(e *CustomEvent) { ... })

// Window management
window, _ := app.Window.GetByName("main")
app.Window.OnCreate(func(window Window) { ... })

// Browser integration
app.Browser.OpenURL("https://wails.io")

// Menu management
menu := app.Menu.New()
app.Menu.Set(menu)

// System tray
systray := app.SystemTray.New()
```

## 管理器參考

### 視窗管理器

管理視窗的建立、擷取及生命週期回呼。

```go
// Create windows
window := app.Window.New()
window := app.Window.NewWithOptions(options)
current := app.Window.Current()

// Find windows
window, exists := app.Window.GetByName("main")
windows := app.Window.GetAll()

// Window callbacks
app.Window.OnCreate(func(window Window) {
    // Handle window creation
})
```

### 事件管理器

處理自訂事件及應用程式事件監聽。

```go
// Custom events
app.Event.Emit("userAction", data)
cancelFunc := app.Event.On("userAction", func(e *CustomEvent) {
    // Handle event
})
app.Event.Off("userAction")
app.Event.Reset() // Remove all listeners

// Application events
app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(e *ApplicationEvent) {
    // Handle system theme change
})
```

### 瀏覽器管理器

提供用於開啟 URL 與檔案的瀏覽器整合功能。

```go
// Open URLs and files in default browser
err := app.Browser.OpenURL("https://wails.io")
err := app.Browser.OpenFile("/path/to/document.pdf")
```

### 環境管理器

提供系統環境資訊的存取。

```go
// Get environment info
env := app.Env.Info()
fmt.Printf("OS: %s, Arch: %s\n", env.OS, env.Arch)

// Check system theme
if app.Env.IsDarkMode() {
    // Dark mode is active
}

// Open file manager
err := app.Env.OpenFileManager("/path/to/folder", false)
```

### 對話方塊管理器

以條理分明的方式存取檔案與訊息對話方塊。

```go
// File dialogs
result, err := app.Dialog.OpenFile().
    AddFilter("Text Files", "*.txt").
    PromptForSingleSelection()

result, err = app.Dialog.SaveFile().
    SetFilename("document.txt").
    PromptForSingleSelection()

// Message dialogs
app.Dialog.Info().
    SetTitle("Information").
    SetMessage("Operation completed successfully").
    Show()

app.Dialog.Error().
    SetTitle("Error").
    SetMessage("An error occurred").
    Show()
```

### 功能表管理器

建立及管理應用程式功能表。

```go
// Create and set application menu
menu := app.Menu.New()
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("New").OnClick(func(ctx *Context) {
    // Handle menu click
})

app.Menu.Set(menu)

// Show about dialog
app.Menu.ShowAbout()
```

### 按鍵綁定管理器

動態管理全域按鍵綁定。

```go
// Add key bindings
app.KeyBinding.Add("ctrl+n", func(window application.Window) {
    // Handle Ctrl+N
})

app.KeyBinding.Add("ctrl+q", func(window application.Window) {
    app.Quit()
})

// Remove key bindings
app.KeyBinding.Remove("ctrl+n")

// Get all bindings
bindings := app.KeyBinding.GetAll()
```

### 內容功能表管理器

進階內容功能表管理（供程式庫作者使用）。

```go
// Create and register context menu
menu := app.ContextMenu.New()
app.ContextMenu.Add("myMenu", menu)

// Retrieve context menu
menu, exists := app.ContextMenu.Get("myMenu")

// Remove context menu
app.ContextMenu.Remove("myMenu")
```

### 螢幕管理器

管理多螢幕設定中的螢幕與座標轉換。

```go
// Get screen information
screens := app.Screen.GetAll()
primary := app.Screen.GetPrimary()

// Coordinate transformations
physicalPoint := app.Screen.DipToPhysicalPoint(logicalPoint)
logicalPoint := app.Screen.PhysicalToDipPoint(physicalPoint)

// Screen detection
screen := app.Screen.ScreenNearestDipPoint(point)
screen = app.Screen.ScreenNearestDipRect(rect)
```

### 剪貼簿管理器

讀取及寫入文字的剪貼簿操作。

```go
// Set text to clipboard
success := app.Clipboard.SetText("Hello World")
if !success {
    // Handle error
}

// Get text from clipboard
text, ok := app.Clipboard.Text()
if !ok {
    // Handle error
} else {
    // Use the text
}
```

### SystemTray 管理器

建立及管理系統匣圖示。

```go
// Create system tray
systray := app.SystemTray.New()
systray.SetLabel("My App")
systray.SetIcon(iconBytes)

// Add menu to system tray
menu := app.Menu.New()
menu.Add("Open").OnClick(func(ctx *Context) {
    // Handle click
})
systray.SetMenu(menu)

// Destroy system tray when done
systray.Destroy()
```

### 自動啟動管理器

註冊應用程式，使其在使用者登入時啟動。系統會依平台選用適當的原生機制：macOS 使用 SMAppService 或 LaunchAgent plist，Windows 使用`HKCU\…\Run`登錄機碼，Linux 則使用 XDG `.desktop`項目。

```go
// Register to launch at login
err := app.Autostart.Enable()

// With extra launch-time arguments and a custom identifier
err = app.Autostart.EnableWithOptions(application.AutostartOptions{
    Identifier: "com.example.myapp",
    Arguments:  []string{"--hidden"},
})

// Check / remove
enabled, err := app.Autostart.IsEnabled()
status, err := app.Autostart.Status()  // includes Path + Strategy
err = app.Autostart.Disable()
```

如需瞭解各平台的行為、識別碼規則及過時項目偵測保證，請參閱[自動啟動功能頁面](/features/autostart/basics/)。
