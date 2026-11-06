---
title: "刘海屏窗口"
description: "创建附着在摄像头模组上的原生 macOS 窗口。"
slug: "features/windows/notch-windows"
sourcePath: "features/windows/notch-windows.md"
---

`NewNotchWindow`会创建一个附着在摄像头模组上的异形、非激活式 macOS 面板。Wails 负责原生定位、黑色连接翼、透明外层画布、窗口层级、Spaces 行为，以及可选的显示和隐藏动画。你的 Web 内容仅占用所请求的内部矩形区域。将指针移入窗口后，其 WebView 会立即成为键盘输入目标，而不会激活应用程序。

@note{type="caution" title="WebView 透明效果使用私有 API"}
`NewNotchWindow`需要使用`-tags private_mac_apis`，透明 Web 内容才能显露原生刘海形状。如果没有该标签，面板、定位和动画仍可正常工作，但 WebView 会保持不透明。请参阅[macOS 私有 API](/guides/build/private-macos-apis/#webview-transparency-and-background)。

@end

<figure>
  <img
    src="/images/notch-notification.gif"
    alt="一条 Wails 刘海通知从 MacBook 摄像头模组向下滑出，显示实时系统指标，然后隐藏回刘海下方"
    loading="lazy"
    decoding="async"
    style="width: 100%; border-radius: 0.75rem"
  />
  <figcaption>
    使用持久化 WebView，并通过原生显示和隐藏过渡效果呈现的刘海通知动画。
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
| `Width` | `int` | `660` | 原生异形边缘内可用的 WebView 宽度。 |
| `Height` | `int` | `92` | 原生异形边缘内可用的 WebView 高度。 |
| `Animated` | `bool` | `false` | 执行`Show`时将窗口向下滑出，执行`Hide`时将其向上滑回。 |
| `AnimationSpeed` | `time.Duration` | `420ms` | 显示动画的持续时间。隐藏动画使用该值的三分之二，默认值为`280ms`。 |
| `Screen` | `*Screen` | 主显示器 | 指定目标显示器。否则，Wails 使用主显示器。 |
| `WindowOptions` | `WebviewWindowOptions` | 默认值 | 提供名称、URL 或 HTML、CSS、JavaScript、按键绑定及其他 WebView 行为。 |

`NewNotchWindow`负责外部尺寸、位置、边框、透明度、调整大小策略、原生面板类、窗口层级和集合行为。`WindowOptions`中这些字段的值会被有意替换，其他字段则会保留。原生背景拖动和 CSS 拖动区域均被禁用，以确保窗口始终附着在摄像头模组上。

返回的`NotchWindow`有意仅公开`Show`、`Hide`、`Visibility`和`Close`；无法通过高级句柄更改原生几何属性。

@note{type="note"}
刘海屏窗口需要 macOS。在没有摄像头模组的 Mac 上，Wails 会将窗口放置在菜单栏下方的顶部中央位置。在不受支持的平台上，`NewNotchWindow`会返回一个惰性句柄；其生命周期方法均为可安全调用的空操作，并且其`Visibility`始终为 false。

@end

## 生命周期

- `Show`会显示现有的原生窗口。启用动画后，窗口会从显示器上方向下滑出。
- `Hide`会保留窗口，以便重复使用。启用动画后，窗口会先滑回显示器上方，再从窗口顺序中移除。WebView、JavaScript 状态、绑定和事件侦听器仍会保持加载，但隐藏的窗口没有悬停目标；应用程序必须调用`Show`才能再次显示该窗口。
- `Visibility`报告当前的原生可见性状态。
- `Close`会永久销毁原生窗口。再次显示该通知前，请创建一个新窗口。
- 指针进入时，会将该刘海屏窗口移到最前，并聚焦其 WebView，以便立即进行键盘交互，同时保留面板的非激活行为。

每次调用`NewNotchWindow`都会创建一个独立窗口，该窗口拥有自己的内容、可见性和动画状态。各窗口使用相同的原生层级，因此最近显示或指针最近进入的实例会出现在最前面，并可能覆盖较早的实例。当窗口以不同屏幕为目标时，每个窗口都会在对应屏幕的摄像头模组上居中。macOS 不会协调属于不同应用程序的刘海屏窗口；如果多个独立应用在相同位置和层级显示窗口，最近加入窗口顺序的窗口会出现在最前面。

对于通知工作负载，应用程序通常会复用隐藏窗口，或在单个 WebView 中维护自己的队列。Wails 不强制实施排队、替换、自动关闭或单窗口策略。

@note{type="note"}
可通过悬停重新打开的持久化折叠界面有意与通知隐藏功能分开，并在[#6009](https://github.com/wailsapp/wails/issues/6009)中跟踪。

@end

请参阅[notch-notification 示例](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/notch-notification)，其中展示了一个紧凑的系统监视器；即使通知经过多次显示和隐藏，其实时 JavaScript 状态仍会保留。
