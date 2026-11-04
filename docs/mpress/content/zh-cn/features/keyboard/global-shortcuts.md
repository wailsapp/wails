---
title: "全局快捷键"
description: "注册系统级键盘快捷键，即使应用程序未获得焦点也能触发"
slug: "features/keyboard/global-shortcuts"
sourcePath: "features/keyboard/global-shortcuts.md"
---

只要 Wails 应用程序仍在运行，无论当前哪个应用程序获得焦点，全局快捷键都能触发。它们非常适合用于显示/隐藏热键、快速捕获工具、媒体控制，以及用户希望能从任何位置访问的其他功能。

@note{type="info" title="全局快捷键与按键绑定"}
[按键绑定](/features/keyboard/shortcuts/)（`app.KeyBinding`）仅在应用程序的某个窗口获得焦点时触发。全局快捷键（`app.GlobalShortcut`）可在系统范围内触发，即使应用程序处于后台也是如此。请根据需要选择。

@end

全局快捷键直接基于各平台的原生功能实现，不会添加任何第三方依赖项。

## 访问全局快捷键管理器

可通过应用程序实例的`GlobalShortcut`属性访问该管理器：

```go
app := application.New(application.Options{
    Name: "Global Shortcuts Demo",
})

globalShortcuts := app.GlobalShortcut
```

## 注册快捷键

`Register`接受一个加速键和一个回调。每当按下该快捷键时，回调都会在自己的 goroutine 中运行。

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

可以在调用`app.Run`之前注册快捷键。应用程序启动时，系统会自动与操作系统建立绑定。

@note{type="tip" title="从回调中操作 UI"}
回调在主线程之外运行。如果回调需要与窗口或其他 UI 交互，窗口方法会代为处理；但对于自定义的主线程任务，请使用`application.InvokeSync`将其包装起来。

@end

### 加速键格式

全局快捷键使用与菜单加速键和按键绑定相同的加速键格式：

```go
"CmdOrCtrl+Shift+G"  // Command on macOS, Control elsewhere
"Ctrl+Alt+K"         // Control + Alt + K
"Cmd+Option+Space"   // Command + Option + Space (macOS)
"Super+D"            // Super / Windows / Logo key + D
"Ctrl+Shift+F5"      // Function keys are supported
```

`CmdOrCtrl`在 macOS 上解析为 Command，在 Windows 和 Linux 上解析为 Control，因此很适合定义跨平台快捷键。

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

应用程序退出时会自动释放所有已注册的快捷键，因此无需手动清理。

## 重复注册同一快捷键时会发生什么

这里有两种不同的情况，Wails 会以不同方式处理。

### 同一应用程序重复注册快捷键

这种情况由 Wails 自行处理，并且在所有平台上的行为都相同。第二次调用`Register`会返回错误，同时保留原有绑定（“报错并保留”）。这样既可确保行为可预测，也会明确暴露错误，而不是悄然替换一个正常工作的快捷键。

```go
app.GlobalShortcut.Register("Ctrl+Shift+G", showWindow)        // ok
err := app.GlobalShortcut.Register("Shift+Ctrl+G", doSomething) // err: already registered
// showWindow is still the active callback for this shortcut.
```

如果要更改快捷键的回调，请先对其调用`Unregister`，然后再次调用`Register`。

### 快捷键已被另一个应用程序占用

这种情况由操作系统决定，因此结果因平台而异：

| 平台 | 快捷键被另一个应用程序占用时的行为 |
| --- | --- |
| **macOS** | 注册成功。macOS 允许多个应用程序注册同一热键，因此系统会将你的回调与现有占用者并列添加，而不会拒绝注册。 |
| **Windows** | 注册失败，`Register`返回错误。最先注册该快捷键的应用程序继续保有它。 |
| **Linux (X11)** | 注册失败，`Register`返回错误，因为 X 服务器会拒绝再次抓取同一组合键。 |
| **Linux (Wayland)** | 由合成器进行仲裁。通常会要求用户通过桌面的全局快捷键对话框批准或选择绑定。 |

鉴于这些差异，请始终检查`Register`返回的错误，并在无法占用快捷键时提供备用快捷键或向用户反馈。

## 平台注意事项

@tabs
[macOS]
全局快捷键使用 Carbon Event Manager 的热键 API。这是 macOS 上实现系统级热键的标准机制，不需要辅助功能权限。

字母快捷键会感知键盘布局。Carbon API 绑定的是物理按键位置，因此 Wails 会将字母转换为当前键盘布局中能够输入该字母的按键（通过`UCKeyTranslate`）。因此，在 AZERTY、QWERTZ、Dvorak 及类似布局上，`Cmd+A`会落到标有**A**的按键，而不是 QWERTY 布局中的对应位置。数字、标点符号和命名按键（方向键、功能键等）使用其固定的物理位置。注册快捷键时会确定绑定；如果在应用程序运行期间切换键盘布局，已有的字母快捷键仍会保留最初绑定的物理按键。

@note{type="caution" title="隐藏快捷键与 `ApplicationShouldTerminateAfterLastWindowClosed`"}
在 macOS 上，`window.Hide()`使用`orderOut:`，这会使窗口不可见。AppKit 会将最后一个不可见窗口视为已关闭，因此，如果设置了`Mac.ApplicationShouldTerminateAfterLastWindowClosed: true`，并使用全局快捷键隐藏唯一的窗口，应用程序将会退出，而不是继续在后台运行。依赖显示/隐藏热键时，请不要设置该选项（默认即为未设置），这样便可隐藏窗口，并在之后再次将其调出。

@end

[Windows]
全局快捷键使用 Win32 `RegisterHotKey` API。系统会抑制自动重复，因此按住按键只会触发一次回调，而不会反复触发。

如果组合键已被另一个应用程序占用，注册将会失败，因此默认快捷键应优先选用不太常见的组合键。

[Linux]
在**X11**会话中，Wails 直接从 X 服务器抓取快捷键，因此请求的加速键会完全按指定内容绑定。

按照设计，在**Wayland**会话中，应用程序无法直接抓取按键。Wails 转而使用 XDG Desktop Portal 的`org.freedesktop.portal.GlobalShortcuts`接口。使用该门户时，传入的加速键只是一个*首选*触发键，最终组合键由合成器（归根结底由用户）决定。激活快捷键时仍会触发回调，但不能保证实际按键与你请求的按键完全一致，并且`IsRegistered`/`GetAll`报告的是你请求的内容，而不是合成器实际绑定的内容。

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

@note{type="danger" title="避开关键系统快捷键"}
某些组合键由操作系统或桌面环境保留，应用程序无法占用。请选择不太可能发生冲突的默认组合键，并始终处理`Register`返回的错误。

@end
