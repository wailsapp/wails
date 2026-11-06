---
title: "通过远程桌面 (RDP) 使用时 WebView2 卡顿"
description: "修复 Wails 应用通过 RDP 会话运行且会话中途显示器 DPI 发生变化时，WebView2 UI 持续数秒的卡顿问题。"
slug: "troubleshooting/windows/rdp"
sourcePath: "troubleshooting/windows/rdp.md"
---

## 问题

通过远程桌面 (RDP) 会话使用 Wails 应用时，以下常见交互可能导致 UI 卡顿数秒：

- 从点击到弹出窗口内容可见大约需要 4 到 8 秒。
- 关闭窗口会使父窗口阻塞大约 2 秒。
- 重新连接后，卡顿状态仍会持续，只有重启主机才能消除。

此问题最常见于 iOS 上的 Microsoft Remote Desktop 客户端，该客户端会在会话进行过程中配置针对 Retina 优化的虚拟显示器。任何在会话中途引入具有不同 DPI 上下文的显示器的 RDP 客户端，都可能触发相同的行为。

## 原因

默认情况下，WebView2 使用窗口托管模式，其合成器表面位于子窗口中。当 RDP 客户端引入 DPI 上下文与会话不同的显示器时，每次 WebView2 控制器调用（`PutIsVisible`、`MoveFocus`、首次绘制和表面释放）都会强制执行同步的 DirectComposition 重新封送。每次重新封送都会阻塞 UI 线程大约 2 秒，因此在大量使用弹出窗口的应用中，卡顿会不断叠加。

同一台计算机上的原生 Win32 + WebView2 应用不受影响，因为它使用可视化托管。这表明问题的原因是托管模式，而不是 WebView2 或 Windows 合成器的普遍问题。

## 解决方案

在 Windows 选项中设置 `UseVisualHosting`，以启用可视化托管。使用可视化托管时，WebView2 的合成器表面通过主机拥有的 DirectComposition 可视化对象进行管理，因此 DPI 上下文变化不再触发同步重新封送。

```go
package main

import "github.com/wailsapp/wails/v3/pkg/application"

func main() {
    app := application.New(application.Options{
        Windows: application.WindowsOptions{
            UseVisualHosting: true,
        },
    })

    // ... create your windows, then:
    app.Run()
}
```

启用后，弹出窗口会在正常的导航时间内打开（约为 150 到 500 毫秒），关闭窗口也不再阻塞父窗口。

@note{type="caution"}
必须在 `app.Run()` 之前设置 `UseVisualHosting`。Wails 会在应用启动期间读取此选项，并在初始化 WebView2 环境之前，将 `COREWEBVIEW2_FORCED_HOSTING_MODE` 环境变量设置为 `COREWEBVIEW2_HOSTING_MODE_WINDOW_TO_VISUAL`。之后再设置不会生效。

@end

该选项默认为 `false`，因此窗口托管模式仍是默认模式。除非现有应用主动启用此选项，否则其行为不会发生变化。

## 何时启用

如果你的应用经常通过 RDP 使用，尤其是通过 iOS 上的 Microsoft Remote Desktop 客户端使用，并且在打开或关闭窗口时出现持续数秒的卡顿，请设置 `UseVisualHosting: true`。如果你的应用不通过 RDP 运行，则无需使用此选项，可以保留默认设置。

## 参考资料

- [WebView2：窗口托管与可视化托管](https://learn.microsoft.com/en-us/microsoft-edge/webview2/concepts/windowed-vs-visual-hosting)
- [WebView2Feedback 问题 #5248](https://github.com/MicrosoftEdge/WebView2Feedback/issues/5248)
- [WebView2Feedback 问题 #4485](https://github.com/MicrosoftEdge/WebView2Feedback/issues/4485)
