---
title: "常見問題"
description: "解答使用 Wails 建置應用程式時最常見的問題 v3"
slug: "faq"
sourcePath: "faq.md"
---

## 一般問題

### Wails 是什麼？

Wails 是一套使用 Go 與 Web 技術建置桌面應用程式的框架。應用程式邏輯以 Go 撰寫，介面則以 HTML、CSS 和 JavaScript（或任何前端框架）建置，並由 Wails 在作業系統的原生 WebView 中呈現。最終得到的是體積小、速度快且具原生使用體驗的應用程式：不需內含瀏覽器、記憶體用量低，而且通常只會產生一個大小約為 10MB 的二進位檔。

### Wails 支援哪些平台？

| 平台 | 需求 |
| --- | --- |
| Windows | AMD64 與 ARM64。使用 [WebView2 執行階段](https://developer.microsoft.com/microsoft-edge/webview2/)。 |
| macOS | Intel 需為 10.15 以上版本（應用程式可將 10.13 以上版本設為目標），Apple Silicon 需為 11.0 以上版本。支援通用二進位檔。 |
| Linux | AMD64 與 ARM64。預設技術堆疊為 GTK4 搭配 WebKitGTK 6.0（Ubuntu 24.04 以上、Debian 13 以上、Fedora 40 以上及類似版本）。僅提供 WebKit2GTK 4.1 的發行版，例如 Ubuntu 22.04、Debian 12 和 RHEL 9，可透過舊版 `-tags gtk3` 建置獲得支援（提供至 v3.1）。不支援僅提供 WebKit2GTK 4.0 的發行版。請參閱 [Linux 建置指南](/guides/build/linux/)。 |
| iOS 與 Android | 實驗性支援。請參閱[行動平台指南](/guides/mobile/)。 |

您也可以使用[伺服器建置](/guides/server-build/)，將應用程式作為一般 Web 應用程式提供。

您可以隨時執行 `wails3 doctor`，檢查系統並取得適用於平台的安裝指示。

### 開始使用需要準備什麼？

- Go 1.26 或更新版本
- Node.js 與 npm（用於前端建置）
- 平台工具鏈：Windows 使用 WebView2（已預先安裝於 10/11）、macOS 使用 Xcode Command Line Tools，Linux 使用 `gcc` 以及 GTK/WebKit 開發套件

`wails3 doctor`會替您檢查所有項目，並明確指出缺少哪些項目。完整步驟請參閱[安裝](/quick-start/installation/)。

### Wails v3已可用於正式環境嗎？

Wails v3.0.0是適用於 Windows、macOS 和 Linux 桌面應用程式的穩定版本。Android 與 iOS 支援仍處於實驗階段，不在穩定版桌面相容性契約的範圍內。請參閱[專案狀態頁面](/status/)。

## 開發

### 我需要懂 Go 嗎？

具備基本的 Go 知識會有所幫助，但不需要成為專家。應用程式邏輯位於一般的 Go 方法中，而[教學](/tutorials/overview/)會引導您完成其餘所有內容。許多開發人員是在建置第一個 Wails 應用程式時學會 Go。

### 我可以使用偏好的前端框架嗎？

可以。只要能建置成 HTML、CSS 和 JavaScript，就能搭配 Wails 使用。Wails 提供 React、Vue、Svelte 和原生 JavaScript 的範本（每種範本都有 TypeScript 版本），其他框架也只需幾分鐘即可完成整合。請參閱[前端框架](/guides/dev/frontend-frameworks/)。

### 如何從 JavaScript 呼叫 Go 函式？

註冊服務後，Wails 會為其產生具型別的繫結：

```go
// Go
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello " + name
}
```

```javascript
// JavaScript
import { GreetService } from "./bindings/changeme";

const message = await GreetService.Greet("World");
```

在 `wails3 dev` 期間會自動重新產生繫結，您也可以使用 `wails3 generate bindings` 隨需產生。請參閱[服務](/features/bindings/services/)。

### 我可以使用 TypeScript 嗎？

可以。繫結產生器會為您的服務及其型別產生 TypeScript 定義，因此對 Go 的呼叫都具有完整型別。

### 如何在 Go 與 JavaScript 之間傳送事件？

```go
// Go
app.Event.Emit("time", time.Now().Format(time.RFC1123))
```

```javascript
// JavaScript
import { Events } from "@wailsio/runtime";

Events.On("time", (event) => {
    console.log(event.data);
});
```

事件名稱必須完全相符。請參閱[事件參考](/guides/events-reference/)。

### 如何偵錯應用程式？

執行 `wails3 dev`，然後在視窗中按一下滑鼠右鍵，即可像在 Web 上一樣開啟瀏覽器開發人員工具。開發伺服器也支援前端熱重新載入。請參閱[偵錯](/guides/dev/debugging/)。

## 建置與散發

### 如何建置正式版本？

```bash
wails3 build
```

二進位檔會輸出至 `bin/`。正式版本建置已套用合理的預設值（建置標籤、`-trimpath`、移除符號），因此不需額外旗標即可產生精簡的二進位檔。

### 可以交叉編譯嗎？

可以，但有一定限制。由於每個平台都使用原生 WebView 程式庫，因此無法使用純 Go 交叉編譯，但常見情境皆有良好支援：

```bash
# Different architecture, same OS
wails3 build GOOS=windows GOARCH=arm64

# macOS universal binary
wails3 task darwin:build:universal
```

從其他作業系統建置 Linux 版本時，會使用以 Docker 為基礎的工具鏈。完整支援矩陣請參閱[跨平台建置](/guides/build/cross-platform/)。

### 如何建立安裝程式或套件？

```bash
wails3 package
```

這會產生平台原生格式；[安裝程式指南](/guides/installers/)涵蓋 Windows 上的 NSIS、macOS 上的 `.app` 套件組合與 DMG，以及 Linux 套件。

### 如何為應用程式進行程式碼簽署？

[簽署指南](/guides/build/signing/)逐步說明 Windows 與 macOS 的簽署流程，包括公證。

## 功能

### 我可以建立多個視窗嗎？

可以，v3原生支援多視窗：

```go
window1 := app.Window.New()
window2 := app.Window.New()
```

請參閱[多個視窗](/features/windows/multiple/)。

### Wails 支援系統匣嗎？

支援，包括選單與點擊處理常式：

```go
systemTray := app.SystemTray.New()
systemTray.SetIcon(iconBytes)
systemTray.SetMenu(myMenu)
```

請參閱[系統匣](/features/menus/systray/)。

### 我可以使用原生對話方塊嗎？

可以。檔案對話方塊、訊息對話方塊和詢問對話方塊全都使用原生實作：

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select File").
    PromptForSingleSelection()
```

請參閱[對話方塊](/features/dialogs/overview/)。

### Wails 支援自動更新嗎？

支援。Wails v3內建自我更新程式（`app.Updater`），具備適用於 GitHub Releases、keygen.sh 和 Sparkle AppCast 的可插拔提供者、密碼編譯簽章驗證，以及可自訂主題或替換的預設 UI。請參閱[應用程式內更新程式](/guides/updater/)指南與[可自我更新的 Wails 應用程式](/tutorials/04-self-update-a-wails-app/)教學。

## 疑難排解

### 有些功能無法運作。我該從哪裡開始？

```bash
wails3 doctor
```

它會驗證您的工具鏈、列出缺少的相依套件及其安裝命令，並輸出您應在任何錯誤報告中附上的版本資訊。

### 建置失敗

常見的修正方式依序如下：

1. `go mod tidy`
2. `cd frontend && npm install`（最常見的原因是缺少 `node_modules`）
3. 更新 CLI：`go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0`
4. 在 Linux 上，檢查 `wails3 doctor` 是否列出缺少的 GTK/WebKit 套件

### 我的繫結遺失或未更新

```bash
wails3 generate bindings
```

在開發模式下，繫結會自動重新產生；如果您在`wails3 dev`之外新增了服務或變更了方法簽章，請手動重新產生繫結。

### 事件未觸發

Go 中`app.Event.Emit("name", ...)`的事件名稱必須與 JavaScript 中`Events.On("name", ...)`的事件名稱完全相符。請先檢查是否有拼字錯誤或大小寫差異。

### 我發現了錯誤

請[提出議題](https://github.com/wailsapp/wails/issues)，並附上`wails3 doctor`的輸出。[意見回饋指南](/feedback/)說明如何撰寫便於處理的問題報告。

## 從v2遷移

### 我應該從v2遷移至v3嗎？

v3帶來多視窗支援、更簡潔的服務型 API、內建更新程式、彈性高出許多的建置系統，以及更出色的效能。新專案應從v3開始。對於現有專案，[遷移指南](/migration/v2-to-v3/)會逐一說明各項差異。

### v2會繼續維護嗎？

會。Wails v2會繼續收到重大更新。新專案應使用 Wails v3.0.0。

### 我可以並行執行v2和v3嗎？

可以。這些 CLI 是各自獨立的二進位檔（`wails`和`wails3`），而且模組具有不同的匯入路徑，因此使用不同主要版本的專案可以順利共存於同一台電腦上。

## 社群

### 如何取得協助？

- 如需快速提問和討論，請使用[Discord](https://discord.gg/JDdSxwjhGf)
- 如需提出較深入的問題，請使用[GitHub Discussions](https://github.com/wailsapp/wails/discussions)
- 如需回報錯誤，請使用[GitHub Issues](https://github.com/wailsapp/wails/issues)

### 如何貢獻？

請參閱[貢獻指南](/contributing/)。隨時歡迎提交錯誤修正；新增功能及公開行為的變更須透過[WEP（Wails 增強提案）](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/README.md)草案 PR 提交。您可選擇在 Discord 或 GitHub Discussions 上進行非正式討論，但並非必要。

### 可以在哪裡找到範例？

此儲存庫隨附超過60個可執行的範例，涵蓋視窗、對話方塊、事件、系統匣、服務等內容：[v3/examples](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples)。

## 還有其他問題嗎？

請在[Discord](https://discord.gg/JDdSxwjhGf)中提問，或[發起討論](https://github.com/wailsapp/wails/discussions)。
