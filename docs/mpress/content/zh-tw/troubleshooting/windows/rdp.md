---
title: "透過遠端桌面（RDP）使用時 WebView2 卡住"
description: "修正在 Wails 應用程式透過 RDP 工作階段執行，且工作階段中途變更顯示器 DPI 時，WebView2 使用者介面停滯數秒的問題。"
slug: "troubleshooting/windows/rdp"
sourcePath: "troubleshooting/windows/rdp.md"
---

## 問題

透過遠端桌面（RDP）工作階段使用 Wails 應用程式時，使用者介面可能會在下列常見互動中停滯數秒：

- 從按一下到快顯視窗顯示內容，大約需要 4 到 8 秒。
- 關閉視窗會使父視窗封鎖約 2 秒。
- 重新連線後，緩慢狀態仍會持續，且只有重新啟動主機後才會消除。

此問題最常出現在 iOS 上的 Microsoft Remote Desktop 用戶端，該用戶端會在工作階段進行途中佈建針對 Retina 最佳化的虛擬顯示器。任何會在工作階段中途引入具有不同 DPI 上下文之顯示器的 RDP 用戶端，都可能觸發相同行為。

## 發生原因

WebView2 預設使用視窗式裝載，其合成器的繪圖表面位於子視窗中。當 RDP 用戶端引入 DPI 上下文與工作階段不同的顯示器時，每次呼叫 WebView2 控制器（`PutIsVisible`、`MoveFocus`、首次繪製和釋放繪圖表面）都會強制執行同步 DirectComposition 重新封送處理。每次重新封送處理都會封鎖 UI 執行緒約2秒，因此大量使用快顯視窗的應用程式會累積多次停滯。

同一台機器上的原生 Win32 加 WebView2 應用程式不受影響，因為它採用視覺裝載。這表示原因在於裝載模式，而不是一般性的 WebView2 或 Windows 合成器問題。

## 解決方案

在 Windows 選項中設定`UseVisualHosting`，以啟用視覺裝載。使用視覺裝載時，主機會透過其擁有的 DirectComposition 視覺物件，持有 WebView2 合成器繪圖表面的擁有權，因此 DPI 上下文變更不再觸發同步重新封送處理。

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

啟用後，快顯視窗會在正常導覽時間內開啟（約 150 到 500 毫秒），且關閉視窗不再封鎖父視窗。

@note{type="caution"}
必須在 `app.Run()` 之前設定 `UseVisualHosting`。Wails 會在應用程式啟動期間讀取此設定，並在 WebView2 環境初始化之前，將 `COREWEBVIEW2_FORCED_HOSTING_MODE` 環境變數設為 `COREWEBVIEW2_HOSTING_MODE_WINDOW_TO_VISUAL`。之後才設定不會生效。

@end

此選項預設為 `false`，因此視窗式裝載仍是預設模式。除非現有應用程式選擇啟用，否則其行為不會改變。

## 何時啟用

如果您的應用程式經常透過 RDP 使用（尤其是使用 iOS 上的 Microsoft Remote Desktop 用戶端），而且在開啟或關閉視窗時會停滯數秒，請設定 `UseVisualHosting: true`。如果您的應用程式不會透過 RDP 執行，則不需要此選項，可以保留預設設定。

## 參考資料

- [WebView2：視窗式裝載與視覺裝載](https://learn.microsoft.com/en-us/microsoft-edge/webview2/concepts/windowed-vs-visual-hosting)
- [WebView2Feedback 議題 #5248](https://github.com/MicrosoftEdge/WebView2Feedback/issues/5248)
- [WebView2Feedback 議題 #4485](https://github.com/MicrosoftEdge/WebView2Feedback/issues/4485)
