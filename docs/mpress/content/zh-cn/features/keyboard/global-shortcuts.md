---
title: "全局快捷键"
description: "注册系统级键盘快捷键，即使应用未获得焦点也能触发"
slug: "features/keyboard/global-shortcuts"
sourcePath: "features/keyboard/global-shortcuts.md"
---

全局快捷键是系统级键盘快捷键。只要 Wails 应用仍在运行，无论当前哪个应用获得焦点，它们都能触发。它们非常适合用于显示/隐藏热键、快速捕获工具、媒体控制，以及用户希望能从任何位置访问的其他功能。

@note{type="info" title="全局快捷键与按键绑定"}
[按键绑定](/features/keyboard/shortcuts/)（`app.KeyBinding`）仅在应用的某个窗口获得焦点时触发。全局快捷键（`app.GlobalShortcut`）可在整个系统范围内触发，即使应用在后台运行也是如此。请根据需要选择。

@end

全局快捷键直接基于各平台的原生功能实现，不会引入任何第三方依赖。

## 访问全局快捷键管理器

可通过应用实例的`GlobalShortcut`属性访问该管理器：

```go
app := application.New(application.Options{
    Name: "Global Shortcuts Demo",
})

globalShortcuts := app.GlobalShortcut
```

## 注册快捷键

`Register`接收一个快捷键组合和一个回调。每当按下该快捷键时，回调都会在其自身的 goroutine 中运行。

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

可以在调用`app.Run`之前注册快捷键。应用启动时，系统会自动完成与操作系统的绑定。

@note{type="tip" title="在回调中操作 UI"}
回调不在主线程上运行。如果回调需要与窗口或其他 UI 交互，窗口方法会替你处理线程调度；但对于自定义的主线程操作，请使用`application.InvokeSync`将其包装起来。

@end

### 快捷键组合格式

全局快捷键使用与菜单快捷键和按键绑定相同的快捷键组合格式：

```go
"CmdOrCtrl+Shift+G"  // Command on macOS, Control elsewhere
"Ctrl+Alt+K"         // Control + Alt + K
"Cmd+Option+Space"   // Command + Option + Space (macOS)
"Super+D"            // Super / Windows / Logo key + D
"Ctrl+Shift+F5"      // Function keys are supported
```

`CmdOrCtrl`在 macOS 上解析为 Command，在 Windows 和 Linux 上解析为 Control，便于定义跨平台快捷键。

## 管理快捷键

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

应用退出时会自动释放所有已注册的快捷键，因此无需手动清理。

## 重复注册同一快捷键时会发生什么

这里有两种不同的情况，Wails 会以不同方式处理。

### 同一应用重复注册快捷键

这种情况由 Wails 自行处理，在所有平台上的行为都相同。第二次调用`Register`会返回错误，并保留原有绑定（“报错并保留”）。这样可以确保行为可预测，并显式指出错误，而不会悄然替换仍可用的快捷键。

```go
app.GlobalShortcut.Register("Ctrl+Shift+G", showWindow)        // ok
err := app.GlobalShortcut.Register("Shift+Ctrl+G", doSomething) // err: already registered
// showWindow is still the active callback for this shortcut.
```

如果要更改快捷键的回调，请先对其调用`Unregister`，然后再次调用`Register`进行注册。

### 快捷键已被另一个应用占用

这种情况由操作系统决定，因此结果因平台而异：

| 平台 | 快捷键被另一个应用占用时的行为 |
| --- | --- |
| **macOS** | 注册成功。macOS 允许多个应用注册同一热键，因此你的回调会与现有占用方的回调并存，而不会被拒绝。 |
| **Windows** | 注册失败，`Register`返回错误。最先注册该快捷键的应用会继续占用它。 |
| **Linux (X11)** | 注册失败，`Register`返回错误，因为 X 服务器会拒绝再次抓取同一组合键。 |
| **Linux (Wayland)** | 由合成器协调。通常会通过桌面的全局快捷键对话框要求用户批准或选择绑定。 |

鉴于这些差异，请始终检查`Register`返回的错误，并在无法占用快捷键时提供备用快捷键或向用户反馈。

## 平台注意事项

@tabs
[macOS]
全局快捷键使用 Carbon Event Manager 的热键 API。这是 macOS 上实现系统级热键的标准机制，无需辅助功能权限。

热键绑定到物理按键位置，因此在非 QWERTY 键盘布局上，快捷键会映射到标准 ANSI/QWERTY 布局中相应位置的按键。

@note{type="caution" title="隐藏快捷键与`ApplicationShouldTerminateAfterLastWindowClosed`"}
在 macOS 上，`window.Hide()`使用`orderOut:`，这会使窗口不可见。AppKit 会将最后一个不可见窗口视为已关闭，因此，如果设置了`Mac.ApplicationShouldTerminateAfterLastWindowClosed: true`，并使用全局快捷键隐藏唯一的窗口，应用将退出，而不是继续在后台运行。依赖隐藏/显示热键时，请不要设置该选项（默认即为未设置），这样窗口便可隐藏，并在之后重新调出。

@end

[Windows]
全局快捷键使用 Win32 `RegisterHotKey` API。系统会抑制自动重复，因此按住组合键只会触发一次回调，而不会反复触发。

如果组合键已被另一个应用占用，注册就会失败，因此请为默认快捷键优先选择不太常见的组合键。

[Linux]
在<strong>X11</strong>会话中，Wails 直接从 X 服务器获取快捷键，因此请求的加速键会完全按照指定方式绑定。

在<strong>Wayland</strong>会话中，按照其设计，应用程序无法直接获取按键。Wails 改用 XDG Desktop Portal 的`org.freedesktop.portal.GlobalShortcuts`接口。使用该门户时，传入的加速键只是一个<em>首选</em>触发键，最终的按键组合由合成器（最终由用户）决定。快捷键激活时，回调仍会触发，但无法保证实际按键与请求的按键完全一致，并且`IsRegistered`/`GetAll`报告的是请求的按键，而不是合成器绑定的按键。

该门户要求桌面环境实现全局快捷键门户（例如较新版本的 GNOME 或 KDE Plasma）。

@end

## 完整示例

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

@note{type="danger" title="避免使用关键系统快捷键"}
某些按键组合由操作系统或桌面环境保留，应用程序无法占用。请选择不易发生冲突的默认组合，并始终处理`Register`返回的错误。

@end
