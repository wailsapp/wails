---
title: "刘海窗口"
description: "创建附着于摄像头外壳的原生 macOS 窗口。"
slug: "features/windows/notch-windows"
sourcePath: "features/windows/notch-windows.md"
---

`NewNotchWindow` 创建一个附着于摄像头外壳、具有特定形状且不会激活应用程序的 macOS 面板。Wails 负责原生定位、黑色连接翼、透明外层画布、窗口层级、空间行为以及可选的显示和隐藏动画。你的 Web 内容仅占用所请求的内部矩形区域。将指针移入窗口后，其 WebView 会立即成为按键事件的接收目标，而不会激活应用程序。

@note{type="caution" title="WebView 透明效果使用私有 API"}
`NewNotchWindow` 需要 `-tags private_mac_apis`，才能让透明的 Web 内容显露原生刘海形状。没有该标签时，面板、定位和动画仍然有效，但 WebView 会保持不透明。请参阅[私有 macOS API](/guides/build/private-macos-apis/#webview-transparency-and-background)。

@end

<figure>
  <img
    src="/images/notch-notification.gif"
    alt="一个从 MacBook 摄像头外壳向下滑出的 Wails 刘海通知，显示实时系统指标，然后重新隐藏到刘海下方"
    loading="lazy"
    decoding="async"
    style="width: 100%; border-radius: 0.75rem"
  />
  <figcaption>
    一个使用持久 WebView 并带有原生显示和隐藏过渡效果的动态刘海通知。
  </figcaption>
</figure>

```go
alert := app.Window.NewNotchWindow(application.NotchWindowOptions{
    Width:    660,
    Height:   92,
    Animated: true,
    WindowOptions: application.WebviewWindowOptions{
        Name: "alert",
        URL:  "/alert",
    },
})

alert.Show()
alert.Hide()
visible := alert.Visibility()
alert.Close()
```

## 选项

| 字段 | 类型 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `Width` | `int` | `660` | 原生形状边缘内的可用 WebView 宽度。 |
| `Height` | `int` | `92` | 原生形状边缘内的可用 WebView 高度。 |
| `Animated` | `bool` | `false` | 执行 `Show` 时窗口向下滑动，执行 `Hide` 时向上滑动。 |
| `AnimationSpeed` | `time.Duration` | `420ms` | 显示持续时间。隐藏持续时间为此值的三分之二，默认是 `280ms`。 |
| `Screen` | `*Screen` | 主显示器 | 指定目标显示器。否则，Wails 使用主显示器。 |
| `WindowOptions` | `WebviewWindowOptions` | 默认值 | 提供名称、URL 或 HTML、CSS、JavaScript、按键绑定以及其他 WebView 行为。 |

`NewNotchWindow` 负责外部尺寸、位置、边框、透明度、调整大小策略、原生面板类、窗口层级和集合行为。`WindowOptions` 中这些字段的值会被有意替换，其他字段则予以保留。原生背景拖动和 CSS 拖动区域均被禁用，以使窗口保持附着于摄像头外壳。

返回的 `NotchWindow` 有意仅公开 `Show`、`Hide`、`Visibility` 和 `Close`；无法通过高级句柄更改原生几何属性。

@note{type="note"}
刘海窗口需要 macOS。在没有摄像头外壳的 Mac 上，Wails 会将窗口置于菜单栏下方的顶部中央。在不受支持的平台上，`NewNotchWindow` 返回一个惰性句柄；调用其生命周期方法不会执行任何操作且是安全的，其 `Visibility` 始终为 false。

@end

## 生命周期

- `Show` 显示现有的原生窗口。启用动画后，窗口会从显示器上方滑下。
- `Hide` 会让窗口保持存活，以供复用。启用动画后，窗口会先滑回显示器上方，再从屏幕显示列表中移除。WebView、JavaScript 状态、绑定和事件监听器仍会保持加载，但隐藏的窗口没有悬停目标；应用程序必须调用 `Show` 才能再次显示它。
- `Visibility` 报告当前的原生可见性。
- `Close` 会永久销毁原生窗口。再次显示该通知前，请创建一个新窗口。
- 指针进入后，该刘海窗口会被置于最前方，其 WebView 会获得焦点，以便立即进行键盘交互，同时保持面板不会激活应用程序的行为。

每次调用 `NewNotchWindow` 都会创建一个独立窗口，各自拥有自己的内容、可见性和动画状态。各窗口使用相同的原生层级，因此最近显示或最近被指针进入的实例会出现在最前方，并可能覆盖更早的实例。当窗口以不同的屏幕为目标时，每个窗口都会在相应屏幕的摄像头外壳处居中。macOS 不会协调属于不同应用程序的刘海窗口；如果不同应用在相同位置和层级显示窗口，最近调整显示顺序的窗口会出现在最前方。

对于通知类工作负载，应用程序通常会复用隐藏的窗口，或在一个 WebView 中维护自己的队列。Wails 不会强制实施排队、替换、自动消失或单窗口策略。

@note{type="note"}
可在悬停时重新打开的持久折叠界面与通知隐藏功能有意分开，相关工作在 [#6009](https://github.com/wailsapp/wails/issues/6009) 中跟踪。

@end

请参阅[notch-notification 示例](https://github.com/wailsapp/wails/tree/master/v3/examples/notch-notification)。它是一个紧凑的系统监视器，其 JavaScript 运行时状态在反复显示和隐藏通知后仍会保留。
