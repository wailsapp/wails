---
title: "全域快捷鍵"
description: "註冊系統範圍的鍵盤快捷鍵，即使應用程式未取得焦點也能觸發"
slug: "features/keyboard/global-shortcuts"
sourcePath: "features/keyboard/global-shortcuts.md"
---

只要 Wails 應用程式仍在執行，無論目前哪個應用程式取得焦點，全域快捷鍵都能在整個系統中觸發。它們非常適合用於顯示／隱藏快速鍵、快速擷取工具、媒體控制，以及使用者期望能從任何位置存取的其他功能。

@note{type="info" title="全域快捷鍵與按鍵繫結的比較"}
[按鍵繫結](/features/keyboard/shortcuts/)（`app.KeyBinding`）只會在應用程式的其中一個視窗取得焦點時觸發。全域快捷鍵（`app.GlobalShortcut`）會在整個系統中觸發，即使應用程式位於背景也一樣。請依需求選用。

@end

全域快捷鍵直接建置於各平台的原生功能之上，不會新增任何第三方相依性。

## 存取全域快捷鍵管理員

您可以透過應用程式執行個體的`GlobalShortcut`屬性存取管理員：

```go
app := application.New(application.Options{
    Name: "Global Shortcuts Demo",
})

globalShortcuts := app.GlobalShortcut
```

## 註冊快捷鍵

`Register`接受一個快捷鍵組合和一個回呼函式。每次按下快捷鍵時，回呼函式都會在自己的 goroutine 上執行。

```go
err := app.GlobalShortcut.Register("CmdOrCtrl+Shift+G", func() {
    // Runs even when another application is focused.
    window.Show()
    window.Focus()
})
if err != nil {
    app.Logger.Error("could not register shortcut", "error", err)
}
```

您可以在呼叫`app.Run`之前註冊快捷鍵。應用程式啟動時，系統會自動完成與作業系統的繫結。

@note{type="tip" title="從回呼函式操作 UI"}
回呼函式不在主執行緒上執行。如果回呼函式需要與視窗或其他 UI 互動，視窗方法會替您處理；但若要執行自訂的主執行緒工作，請用`application.InvokeSync`包裝。

@end

### 快捷鍵組合格式

全域快捷鍵使用與選單快捷鍵和按鍵繫結相同的快捷鍵組合格式：

```go
"CmdOrCtrl+Shift+G"  // Command on macOS, Control elsewhere
"Ctrl+Alt+K"         // Control + Alt + K
"Cmd+Option+Space"   // Command + Option + Space (macOS)
"Super+D"            // Super / Windows / Logo key + D
"Ctrl+Shift+F5"      // Function keys are supported
```

`CmdOrCtrl`在 macOS 上會解析為 Command，在 Windows 和 Linux 上則會解析為 Control，方便建立跨平台快捷鍵。

## 管理快捷鍵

```go
// Check whether a shortcut is registered (modifier order does not matter).
registered := app.GlobalShortcut.IsRegistered("Ctrl+Shift+G")

// List every shortcut this application has registered.
for _, accelerator := range app.GlobalShortcut.GetAll() {
    app.Logger.Info("global shortcut", "accelerator", accelerator)
}

// Release a single shortcut.
app.GlobalShortcut.Unregister("Ctrl+Shift+G")

// Release everything (also done automatically on shutdown).
app.GlobalShortcut.UnregisterAll()
```

應用程式結束時會自動釋放所有已註冊的快捷鍵，因此不需要手動清理。

## 重複註冊同一個快捷鍵時會發生什麼事

這分為兩種不同的情況，Wails 會以不同方式處理。

### 同一個應用程式重複註冊快捷鍵

此情況由 Wails 自行處理，且在所有平台上的行為都相同。第二次呼叫`Register`會傳回錯誤，並保留原有繫結（「發生錯誤並保留」）。這可確保行為可預期，並明確呈現錯誤，而不會無聲地取代運作中的快捷鍵。

```go
app.GlobalShortcut.Register("Ctrl+Shift+G", showWindow)        // ok
err := app.GlobalShortcut.Register("Shift+Ctrl+G", doSomething) // err: already registered
// showWindow is still the active callback for this shortcut.
```

若要變更快捷鍵的回呼函式，請先對它呼叫`Unregister`，再使用`Register`重新註冊。

### 另一個應用程式已占用快捷鍵

此情況由作業系統決定，因此結果因平台而異：

| 平台 | 另一個應用程式已占用快捷鍵時的行為 |
| --- | --- |
| **macOS** | 註冊成功。macOS 允許多個應用程式註冊同一個快速鍵，因此您的回呼函式會與現有擁有者並存，而不會遭到拒絕。 |
| **Windows** | 註冊失敗，且`Register`會傳回錯誤。最先註冊此快捷鍵的應用程式會繼續持有它。 |
| **Linux (X11)** | 註冊失敗，且`Register`會傳回錯誤，因為 X 伺服器會拒絕再次擷取相同的按鍵組合。 |
| **Linux (Wayland)** | 由合成器仲裁。通常會要求使用者透過桌面環境的全域快捷鍵對話方塊核准或選擇繫結。 |

由於存在這些差異，請務必檢查`Register`傳回的錯誤，並在無法取得快捷鍵時提供備用快捷鍵或使用者回饋。

## 平台注意事項

@tabs
[macOS]
全域快捷鍵使用 Carbon Event Manager 的快速鍵 API。這是 macOS 上系統範圍快速鍵的標準機制，不需要「輔助使用」權限。

快速鍵會繫結至實體按鍵位置，因此在非 QWERTY 鍵盤配置上，快捷鍵會對應至標準 ANSI/QWERTY 位置上的按鍵。

@note{type="caution" title="隱藏快捷鍵與 `ApplicationShouldTerminateAfterLastWindowClosed`"}
在 macOS 上，`window.Hide()`會使用`orderOut:`，使視窗變為不可見。AppKit 會將最後一個不可見視窗視為已關閉，因此，如果您設定了`Mac.ApplicationShouldTerminateAfterLastWindowClosed: true`，並使用全域快捷鍵隱藏唯一的視窗，應用程式將會結束，而不是繼續在背景執行。若您仰賴顯示／隱藏快速鍵，請勿設定該選項（預設即為未設定），如此便能隱藏視窗，並在之後再次將其叫出。

@end

[Windows]
全域快捷鍵使用 Win32 `RegisterHotKey` API。系統會抑制自動重複，因此持續按住按鍵只會觸發一次回呼函式，而不會重複觸發。

如果另一個應用程式已占用該按鍵組合，註冊就會失敗，因此預設值應優先選用較不常見的組合。

[Linux]
在<strong>X11</strong>工作階段中，Wails 會直接向 X 伺服器擷取快捷鍵，因此要求的快速鍵會完全依照指定內容繫結。

在<strong>Wayland</strong>工作階段中，系統在設計上不允許應用程式直接擷取按鍵。Wails 會改用 XDG Desktop Portal 的`org.freedesktop.portal.GlobalShortcuts`介面。使用該介面時，您傳入的快速鍵只是<em>偏好的</em>觸發方式，最終按鍵組合由合成器（以及最終由使用者）決定。啟用快捷鍵時仍會觸發您的回呼，但實際按鍵不保證與要求的內容相符，而且`IsRegistered`/`GetAll`回報的是您要求的內容，而不是合成器所繫結的內容。

此 Portal 介面需要桌面環境實作全域快捷鍵 Portal 介面（例如較新版本的 GNOME 或 KDE Plasma）。

@end

## 完整範例

```go
package main

import (
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "Global Shortcuts Demo",
    })

    window := app.Window.New()

    // Bring the window to the front from anywhere.
    if err := app.GlobalShortcut.Register("CmdOrCtrl+Shift+G", func() {
        window.Show()
        window.Focus()
    }); err != nil {
        log.Printf("could not register show shortcut: %v", err)
    }

    // Hide the window from anywhere.
    if err := app.GlobalShortcut.Register("CmdOrCtrl+Shift+H", func() {
        window.Hide()
    }); err != nil {
        log.Printf("could not register hide shortcut: %v", err)
    }

    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

@note{type="danger" title="避免使用重要的系統快捷鍵"}
部分按鍵組合由作業系統或桌面環境保留，應用程式無法註冊。請選擇不易發生衝突的預設組合，並一律處理`Register`傳回的錯誤。

@end
