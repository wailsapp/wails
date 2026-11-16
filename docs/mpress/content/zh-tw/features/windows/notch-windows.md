---
title: "瀏海視窗"
description: "建立附著於相機模組的原生 macOS 視窗。"
slug: "features/windows/notch-windows"
sourcePath: "features/windows/notch-windows.md"
---

`NewNotchWindow` 會建立附著於相機模組、具有特殊形狀且不會啟用應用程式的 macOS 面板。原生位置、黑色連接翼、透明外層畫布、視窗層級、Spaces 行為，以及選用的顯示與隱藏動畫都由 Wails 管理。網頁內容只會占用要求的內部矩形區域。將指標移入視窗後，其 webview 會立即成為鍵盤輸入目標，而不會啟用應用程式。

@note{type="caution" title="Webview 透明效果使用私有 API"}
`NewNotchWindow` 需要 `-tags private_mac_apis`，才能讓透明網頁內容顯露原生瀏海形狀。若未加上此標籤，面板、位置與動畫仍可正常運作，但 webview 會維持不透明。請參閱[私有 macOS API](/guides/build/private-macos-apis/#webview-transparency-and-background)。

@end

<figure>
  <img
    src="/images/notch-notification.gif"
    alt="從 MacBook 相機模組向下滑出的 Wails 瀏海通知，顯示即時系統指標後再隱藏回瀏海下方"
    loading="lazy"
    decoding="async"
    style="width: 100%; border-radius: 0.75rem"
  />
  <figcaption>
    使用持續存在的 webview，並搭配原生顯示與隱藏轉場的動畫瀏海通知。
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
| `Width` | `int` | `660` | 原生特殊形狀邊緣內可用的 webview 寬度。 |
| `Height` | `int` | `92` | 原生特殊形狀邊緣內可用的 webview 高度。 |
| `Animated` | `bool` | `false` | 在 `Show` 時將視窗向下滑入，在 `Hide` 時向上滑出。 |
| `AnimationSpeed` | `time.Duration` | `420ms` | 顯示動畫的持續時間。隱藏動畫使用此值的三分之二，預設為 `280ms`。 |
| `Screen` | `*Screen` | 主要顯示器 | 指定目標顯示器。若未指定，Wails 會使用主要顯示器。 |
| `WindowOptions` | `WebviewWindowOptions` | 預設值 | 提供名稱、URL 或 HTML、CSS、JavaScript、按鍵繫結及其他 webview 行為。 |

`NewNotchWindow` 會管理外部尺寸、位置、框架、透明度、調整大小原則、原生面板類別、視窗層級及集合行為。`WindowOptions` 中這些欄位的值會刻意被取代，其他欄位則會保留。原生背景拖曳與 CSS 拖曳區域會停用，使視窗保持附著於相機模組。

傳回的 `NotchWindow` 刻意只公開 `Show`、`Hide`、`Visibility` 和 `Close`；無法透過高階控制代碼變更原生幾何配置。

@note{type="note"}
瀏海視窗需要 macOS。在沒有相機模組的 Mac 上，Wails 會將視窗置於選單列下方的頂端中央。在不支援的平台上，`NewNotchWindow` 會傳回非作用中的控制代碼；其生命週期方法是可安全呼叫且不執行任何操作的函式，而其 `Visibility` 永遠為 false。

@end

## 生命週期

- `Show` 會顯示既有的原生視窗。啟用動畫時，視窗會從顯示器上方往下滑入。
- `Hide` 會讓視窗保持存續，以便重複使用。啟用動畫時，視窗會先滑回顯示器上方，再從視窗顯示順序中移除。webview、JavaScript 狀態、繫結與事件監聽器仍會保持載入，但隱藏的視窗沒有可供游標懸停的目標；應用程式必須呼叫 `Show` 才能再次顯示。
- `Visibility` 會回報目前的原生可見性。
- `Close` 會永久銷毀原生視窗。再次顯示該通知前，請先建立新視窗。
- 指標進入時，該瀏海視窗會移至最前方，且其 webview 會取得焦點以立即接受鍵盤操作，同時維持不啟用應用程式的面板行為。

每次呼叫 `NewNotchWindow` 都會建立一個獨立視窗，並擁有自己的內容、可見性與動畫狀態。所有視窗使用相同的原生層級，因此最近顯示或有指標進入的執行個體會出現在最前方，並可能覆蓋較早的執行個體。當視窗以不同螢幕為目標時，每個視窗都會置中於該螢幕的相機模組。macOS 不會協調分屬不同應用程式的瀏海視窗；如果不同應用程式在相同位置與層級顯示視窗，最近排入顯示順序的視窗會出現在最前方。

對於通知工作負載，應用程式通常會重複使用隱藏的視窗，或在單一 webview 中自行維護佇列。Wails 不會強制採用佇列、取代、自動關閉或單一視窗原則。

@note{type="note"}
可在游標懸停時重新開啟的持續摺疊介面，刻意與通知隱藏機制分開處理，並在 [#6009](https://github.com/wailsapp/wails/issues/6009) 中追蹤。

@end

請參閱[瀏海通知範例](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/notch-notification)；這是一個精簡的系統監視器，其即時 JavaScript 狀態可在反覆顯示及隱藏通知的循環中持續保留。
