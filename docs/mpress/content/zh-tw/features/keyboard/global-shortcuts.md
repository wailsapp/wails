---
title: "全域快速鍵"
description: "註冊即使應用程式未取得焦點也能觸發的系統全域鍵盤快速鍵"
slug: "features/keyboard/global-shortcuts"
sourcePath: "features/keyboard/global-shortcuts.md"
---

只要 Wails 應用程式仍在執行，無論目前哪個應用程式取得焦點，全域快速鍵都能在整個系統中觸發。它們非常適合用於顯示／隱藏快速鍵、快速擷取工具、媒體控制，以及使用者預期可從任何地方存取的其他功能。

@note{type="info" title="全域快速鍵與按鍵綁定"}
[按鍵綁定](/features/keyboard/shortcuts/)（`app.KeyBinding`）只會在應用程式的其中一個視窗取得焦點時觸發。全域快速鍵（`app.GlobalShortcut`）則會在整個系統中觸發，即使應用程式在背景執行也一樣。請依需求選用。

@end

全域快速鍵直接以各平台的原生功能實作，不會加入任何第三方相依套件。

## 存取全域快速鍵管理員

管理員可透過應用程式執行個體的 `GlobalShortcut` 屬性取得：

```go
app := application.New(application.Options{
    Name: "Global Shortcuts Demo",
})

globalShortcuts := app.GlobalShortcut
```

## 註冊快速鍵

`Register` 接受一個加速鍵和一個回呼。每當按下快速鍵時，回呼會在自己的 goroutine 中執行。

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

你可以在呼叫 `app.Run` 前註冊快速鍵。應用程式啟動時，系統會自動完成與作業系統的綁定。

@note{type="tip" title="從回呼操作 UI"}
回呼不會在主執行緒上執行。如果回呼需要與視窗或其他 UI 互動，視窗方法會代為處理；若要執行自訂的主執行緒工作，請使用 `application.InvokeSync` 將其包裝起來。

@end

### 加速鍵格式

全域快速鍵使用與選單加速鍵及按鍵綁定相同的加速鍵格式：

```go
"CmdOrCtrl+Shift+G"  // Command on macOS, Control elsewhere
"Ctrl+Alt+K"         // Control + Alt + K
"Cmd+Option+Space"   // Command + Option + Space (macOS)
"Super+D"            // Super / Windows / Logo key + D
"Ctrl+Shift+F5"      // Function keys are supported
```

`CmdOrCtrl` 在 macOS 上會解析為 Command，在 Windows 和 Linux 上則會解析為 Control，因此很適合用於跨平台快速鍵。

## 管理快速鍵

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

應用程式結束時，所有已註冊的快速鍵都會自動釋放，因此不需要手動清理。

## 重複註冊同一快速鍵時會發生什麼事

這分為兩種不同情況，Wails 會採用不同方式處理。

### 同一應用程式重複註冊快速鍵

此情況由 Wails 自行處理，且在所有平台上的行為完全相同。第二次呼叫 `Register` 會傳回錯誤，並保留原有綁定（「報錯並保留」）。如此可使行為保持可預期，並明確呈現錯誤，而不會悄悄取代正常運作的快速鍵。

```go
app.GlobalShortcut.Register("Ctrl+Shift+G", showWindow)        // ok
err := app.GlobalShortcut.Register("Shift+Ctrl+G", doSomething) // err: already registered
// showWindow is still the active callback for this shortcut.
```

若要變更快速鍵的回呼，請先對它呼叫 `Unregister`，再使用 `Register` 重新註冊。

### 另一個應用程式已占用快速鍵

此情況由作業系統決定，因此結果因平台而異：

| 平台 | 另一個應用程式占用快速鍵時的行為 |
| --- | --- |
| **macOS** | 註冊成功。macOS 允許多個應用程式註冊同一組快速鍵，因此系統不會拒絕你的回呼，而是將它加入現有占用者的回呼 alongside。 |
| **Windows** | 註冊失敗，且 `Register` 會傳回錯誤。最先註冊快速鍵的應用程式會繼續保有它。 |
| **Linux（X11）** | 註冊失敗，且 `Register` 會傳回錯誤，因為 X 伺服器會拒絕再次擷取同一組按鍵組合。 |
| **Linux（Wayland）** | 由合成器仲裁。系統通常會要求使用者透過桌面環境的全域快速鍵對話方塊核准或選擇綁定。 |

由於存在這些差異，請一律檢查 `Register` 傳回的錯誤，並在無法占用快速鍵時提供備用快速鍵或使用者回饋。

## 平台注意事項

@tabs
[macOS]
全域快速鍵使用 Carbon Event Manager 的快速鍵 API。這是 macOS 上系統全域快速鍵的標準機制，不需要「輔助使用」權限。

字母快速鍵會感知鍵盤配置。Carbon API 綁定的是實體按鍵位置，因此 Wails 會將字母轉換為目前使用中的鍵盤配置上可輸出該字母的按鍵（透過 `UCKeyTranslate`）。因此，在 AZERTY、QWERTZ、Dvorak 及類似配置上，`Cmd+A` 會對應至標示為 **A** 的按鍵，而不是 QWERTY 配置中的位置。數字、標點符號及具名按鍵（方向鍵、功能鍵等）則使用其固定的實體位置。系統會在註冊快速鍵時解析綁定；若在應用程式執行期間切換鍵盤配置，現有的字母快速鍵仍會保留最初綁定的實體按鍵。

@note{type="caution" title="隱藏快速鍵與 `ApplicationShouldTerminateAfterLastWindowClosed`"}
`window.Hide()` 在 macOS 上使用 `orderOut:`，這會使視窗不可見。AppKit 會將最後一個不可見的視窗視為已關閉，因此如果你設定 `Mac.ApplicationShouldTerminateAfterLastWindowClosed: true`，並使用全域快速鍵隱藏唯一的視窗，應用程式將會結束，而不會繼續在背景執行。依賴顯示／隱藏快速鍵時，請不要設定該選項（預設即為未設定），如此視窗才能先隱藏，之後再重新叫出。

@end

[Windows]
全域快速鍵使用 Win32 `RegisterHotKey` API。系統會抑制自動重複，因此按住按鍵只會觸發一次回呼，不會反覆觸發。

如果另一個應用程式已占用該按鍵組合，註冊就會失敗，因此預設值應優先選用較不常見的組合。

[Linux]
在 **X11** 工作階段中，Wails 會直接向 X 伺服器擷取快速鍵，因此會完全按照指定內容綁定所要求的加速鍵。

在 **Wayland** 工作階段中，依其設計，應用程式無法直接擷取按鍵。Wails 會改用 XDG Desktop Portal 的 `org.freedesktop.portal.GlobalShortcuts` 介面。使用入口網站時，你傳入的加速鍵是*偏好*的觸發方式，最終按鍵組合則由合成器（以及最終由使用者）決定。快速鍵啟用時仍會觸發回呼，但實際按鍵不保證與要求相符，而且 `IsRegistered`／`GetAll` 回報的是你要求的內容，而不是合成器實際綁定的內容。

入口網站需要桌面環境實作全域快速鍵入口網站，例如較新版本的 GNOME 或 KDE Plasma。

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

@note{type="danger" title="避開重要的系統快速鍵"}
部分按鍵組合由作業系統或桌面環境保留，應用程式無法占用。請選擇不太可能發生衝突的預設值，並一律處理 `Register` 傳回的錯誤。

@end
