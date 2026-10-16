---
title: "瀏海視窗"
description: "建立附著於相機模組外殼的原生 macOS 視窗。"
slug: "features/windows/notch-windows"
sourcePath: "features/windows/notch-windows.md"
---

`NewNotchWindow` 會建立一個附著於相機模組外殼、具有特定形狀且不會啟用應用程式的 macOS 面板。Wails 負責原生定位、黑色連接翼、透明外層畫布、視窗層級、Spaces 行為，以及可選的顯示與隱藏動畫。網頁內容只會占用所要求的內部矩形區域。將指標移入視窗後，其 WebView 會立即成為鍵盤輸入目標，但不會啟用應用程式。

@note{type="caution" title="WebView 透明效果使用私有 API"}
`NewNotchWindow` 需要 `-tags private_mac_apis`，才能讓透明網頁內容顯露原生瀏海形狀。若沒有此標記，面板、定位和動畫仍可運作，但 WebView 會維持不透明。請參閱[私有 macOS API](/guides/build/private-macos-apis/#webview-transparency-and-background)。

@end

<figure>
  <img
    src="/images/notch-notification.gif"
    alt="Wails 瀏海通知從 MacBook 相機模組外殼向下滑出，顯示即時系統指標，然後隱藏回瀏海下方"
    loading="lazy"
    decoding="async"
    style="width: 100%; border-radius: 0.75rem"
  />
  <figcaption>
    使用持續存在的 WebView，並透過原生顯示與隱藏轉場呈現動畫效果的瀏海通知。
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

## 選項

| 欄位 | 型別 | 預設值 | 說明 |
| --- | --- | --- | --- |
| `Width` | `int` | `660` | 原生造型邊緣內可用的 WebView 寬度。 |
| `Height` | `int` | `92` | 原生造型邊緣內可用的 WebView 高度。 |
| `Animated` | `bool` | `false` | 在 `Show` 時將視窗向下滑出，並在 `Hide` 時向上滑回。 |
| `AnimationSpeed` | `time.Duration` | `420ms` | 顯示持續時間。隱藏會使用此值的三分之二，預設為 `280ms`。 |
| `Screen` | `*Screen` | 主要顯示器 | 指定特定顯示器。否則，Wails 會使用主要顯示器。 |
| `WindowOptions` | `WebviewWindowOptions` | 預設值 | 提供名稱、URL 或 HTML、CSS、JavaScript、按鍵繫結，以及其他 WebView 行為。 |

`NewNotchWindow` 負責外部尺寸、位置、框架、透明度、調整大小原則、原生面板類別、視窗層級及集合行為。`WindowOptions` 中這些欄位的值會刻意遭到取代，其他欄位則會保留。原生背景拖曳與 CSS 拖曳區域皆會停用，使視窗保持附著於相機模組外殼。

傳回的 `NotchWindow` 刻意只公開 `Show`、`Hide`、`Visibility` 和 `Close`；無法透過高階控制代碼變更原生幾何配置。

@note{type="note"}
瀏海視窗需要 macOS。在沒有相機模組外殼的 Mac 上，Wails 會將視窗放置於選單列下方的頂端中央。在不支援的平台上，`NewNotchWindow` 會傳回不具作用的控制代碼；呼叫其生命週期方法可安全地不執行任何操作，且其 `Visibility` 一律為 false。

@end

## 生命週期

- `Show` 會顯示現有的原生視窗。啟用動畫時，視窗會從顯示器上方往下滑出。
- `Hide` 會讓視窗保持存續以供重複使用。啟用動畫時，視窗會先滑回顯示器上方，再從視窗顯示順序中移除。WebView、JavaScript 狀態、繫結和事件接聽器都會維持載入，但隱藏的視窗沒有可供懸停的目標；應用程式必須呼叫 `Show` 才能再次顯示它。
- `Visibility` 會回報目前的原生可見狀態。
- `Close` 會永久銷毀原生視窗。再次顯示該通知前，請先建立新視窗。
- 指標進入後，該瀏海視窗會移至最前方，並讓其 WebView 取得焦點以立即接收鍵盤操作，同時維持不啟用應用程式的面板行為。

每次呼叫 `NewNotchWindow` 都會建立獨立視窗，各自擁有自己的內容、可見狀態和動畫狀態。這些視窗使用相同的原生層級，因此最近顯示或指標進入的執行個體會出現在最前方，並可能覆蓋較早的執行個體。當視窗以不同螢幕為目標時，每個視窗都會在其所屬螢幕的相機模組外殼上置中。macOS 不會協調分屬不同應用程式的瀏海視窗；若不同應用程式在相同位置與層級顯示視窗，最近排入顯示順序的視窗會出現在最前方。

對於通知工作負載，應用程式通常會重複使用隱藏的視窗，或在單一 WebView 中自行維護佇列。Wails 不會強制採用佇列、取代、自動關閉或單一視窗原則。

@note{type="note"}
懸停時重新開啟的持續性收合介面，刻意與通知隱藏功能分開，並於 [#6009](https://github.com/wailsapp/wails/issues/6009) 中追蹤。

@end

請參閱[notch-notification 範例](https://github.com/wailsapp/wails/tree/master/v3/examples/notch-notification)，其中示範一個精簡的系統監視器；其即時 JavaScript 狀態會在反覆顯示與隱藏通知的週期中持續保留。
