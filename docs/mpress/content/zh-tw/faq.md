---
title: "常見問題"
description: "Wails v3 應用程式建置常見問題解答"
slug: "faq"
sourcePath: "faq.md"
---

## 一般問題

### Wails 是什麼？

Wails 是一套使用 Go 與 Web 技術建置桌面應用程式的框架。應用程式邏輯以 Go 撰寫，介面則使用 HTML、CSS 和 JavaScript（或任何前端框架）建置，再由 Wails 在作業系統的原生 WebView 中呈現。最終得到的是小巧、快速且具原生體驗的應用程式：不綑綁瀏覽器、記憶體用量低，而且通常只有約 10 MB 的單一二進位檔。

### Wails 支援哪些平台？

| 平台 | 需求 |
| --- | --- |
| Windows | AMD64 與 ARM64。使用 [WebView2 執行階段](https://developer.microsoft.com/microsoft-edge/webview2/)。 |
| macOS | Intel 上需為 10.15 以上版本（應用程式可將 10.13 以上版本設為目標），Apple Silicon 上需為 11.0 以上版本。支援通用二進位檔。 |
| Linux | AMD64 與 ARM64。預設技術堆疊為 GTK4 搭配 WebKitGTK 6.0（Ubuntu 24.04 以上版本、Debian 13 以上版本、Fedora 40 以上版本及類似發行版）。僅提供 WebKit2GTK 4.1 的發行版（例如 Ubuntu 22.04、Debian 12 和 RHEL 9）可透過舊版 `-tags gtk3` 建置獲得支援（提供至 v3.1）。不支援僅提供 WebKit2GTK 4.0 的發行版。請參閱[Linux 建置指南](/guides/build/linux/)。 |
| iOS 與 Android | 實驗性支援。請參閱[行動裝置指南](/guides/mobile/)。 |

您也可以使用[伺服器建置](/guides/server-build/)，將應用程式作為一般 Web 應用程式提供服務。

您可以隨時執行 `wails3 doctor` 來檢查系統，並取得平台專用的安裝指示。

### 開始使用需要準備什麼？

- Go 1.25 或更新版本
- Node.js 與 npm（用於前端建置）
- 平台工具鏈：Windows 上的 WebView2（已預先安裝於 10/11）、macOS 上的 Xcode Command Line Tools，以及 Linux 上的 `gcc` 與 GTK/WebKit 開發套件

`wails3 doctor` 會替您檢查上述所有項目，並明確指出缺少的項目。完整步驟請參閱[安裝](/quick-start/installation/)。

### Wails v3 已可用於正式環境嗎？

Wails v3 是測試版軟體，具備穩定的桌面 API。目前已有應用程式使用它在正式環境中執行，但在我們完成 3.0 的最後修整期間，您仍應在部署前進行完整測試。最新狀況請參閱[專案狀態頁面](/status/)。Wails v2 是目前的穩定版本，並會繼續獲得修正。

## 開發

### 我需要懂 Go 嗎？

具備基本的 Go 知識會有幫助，但您不需要成為專家。應用程式邏輯位於一般的 Go 方法中，其他所有內容則可依照[教學](/tutorials/overview/)逐步完成。許多開發人員會在建置第一個 Wails 應用程式的過程中學會 Go。

### 可以使用我偏好的前端框架嗎？

可以。只要能建置為 HTML、CSS 和 JavaScript，就能與 Wails 搭配使用。Wails 隨附 React、Vue、Svelte 和原生 JavaScript 的範本（各自皆有 TypeScript 版本），其他框架也都能在幾分鐘內完成整合。請參閱[前端框架](/guides/dev/frontend-frameworks/)。

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

繫結會在執行 `wails3 dev` 期間自動重新產生，也可使用 `wails3 generate bindings` 隨需產生。請參閱[服務](/features/bindings/services/)。

### 可以使用 TypeScript 嗎？

可以。繫結產生器會為您的服務及其型別產生 TypeScript 定義，因此對 Go 的呼叫都有完整的型別資訊。

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

執行 `wails3 dev`，然後在視窗中按一下滑鼠右鍵以開啟瀏覽器開發人員工具，操作方式與 Web 開發完全相同。開發伺服器也支援前端熱重新載入。請參閱[偵錯](/guides/dev/debugging/)。

## 建置與發佈

### 如何建置正式環境版本？

```bash
wails3 build
```

您的二進位檔會產生在 `bin/`。正式環境建置已套用合理的預設值（建置標籤、`-trimpath`、移除符號），因此不需要額外的旗標即可產生精簡的二進位檔。

### 可以交叉編譯嗎？

可以，但有一定限制。由於各平台都使用原生 WebView 程式庫，因此無法直接套用純 Go 交叉編譯；不過，常見情境皆有良好支援：

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

這會產生平台原生格式；[安裝程式指南](/guides/installers/)則涵蓋 Windows 上的 NSIS、macOS 上的 `.app`套件組合與 DMG，以及 Linux 套件。

### 如何為應用程式進行程式碼簽署？

[簽署指南](/guides/build/signing/)逐步說明 Windows 與 macOS 的簽署流程，包括公證。

## 功能

### 可以建立多個視窗嗎？

是的，v3 原生支援多視窗：

```go
window1 := app.Window.New()
window2 := app.Window.New()
```

請參閱[多視窗](/features/windows/multiple/)。

### Wails 支援系統匣嗎？

是的，包括選單和點擊處理常式：

```go
systemTray := app.SystemTray.New()
systemTray.SetIcon(iconBytes)
systemTray.SetMenu(myMenu)
```

請參閱[系統匣](/features/menus/systray/)。

### 我可以使用原生對話方塊嗎？

可以。檔案對話方塊、訊息對話方塊和詢問對話方塊均採用原生實作：

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select File").
    PromptForSingleSelection()
```

請參閱[對話方塊](/features/dialogs/overview/)。

### Wails 支援自動更新嗎？

是的。Wails v3 內建自我更新程式（`app.Updater`），可插接 GitHub Releases、keygen.sh 和 Sparkle AppCast 的提供者，支援密碼編譯簽章驗證，並提供可自訂主題或替換的預設 UI。請參閱[應用程式內更新程式](/guides/updater/)指南和[可自我更新的 Wails 應用程式](/tutorials/04-self-update-a-wails-app/)教學。

## 疑難排解

### 有些功能無法正常運作。我該從哪裡著手？

```bash
wails3 doctor
```

它會驗證您的工具鏈、列出缺少的相依套件及其安裝命令，並輸出您應在任何錯誤報告中附上的版本資訊。

### 建置失敗

常見的修正方式依序如下：

1. `go mod tidy`
2. `cd frontend && npm install`（最常見的原因是缺少`node_modules`）
3. 更新 CLI：`go install github.com/wailsapp/wails/v3/cmd/wails3@latest`
4. 在 Linux 上，檢查`wails3 doctor`是否顯示缺少 GTK/WebKit 套件

### 繫結遺失或已過期

```bash
wails3 generate bindings
```

在開發模式下，繫結會自動重新產生；如果您在`wails3 dev`之外新增了服務或變更了方法簽章，請手動重新產生繫結。

### 事件未觸發

Go 中`app.Event.Emit("name", ...)`的事件名稱必須與 JavaScript 中`Events.On("name", ...)`的事件名稱完全相符。請先檢查拼寫錯誤和大小寫差異。

### 我發現了錯誤

請[建立議題](https://github.com/wailsapp/wails/issues)並附上`wails3 doctor`的輸出。[意見回饋指南](/feedback/)說明了如何撰寫便於採取行動的報告。

## 從 v2 遷移

### 我應該從 v2 遷移至 v3 嗎？

v3 帶來多視窗支援、更簡潔的服務型 API、內建更新程式、靈活許多的建置系統，以及更佳的效能。新專案應從 v3 開始。對於現有專案，[遷移指南](/migration/v2-to-v3/)會逐一說明其中的差異。

### v2 會繼續維護嗎？

會。在 v3 邁向穩定版本期間，v2 會繼續獲得修正。

### 我可以並行使用 v2 和 v3 嗎？

可以。兩者的 CLI 是不同的二進位檔（`wails`和`wails3`），模組也有不同的匯入路徑，因此採用不同主要版本的專案可以在同一台機器上順利共存。

## 社群

### 如何取得協助？

- 如有簡短問題或想參與討論，請使用[Discord](https://discord.gg/JDdSxwjhGf)
- 如有需要詳述的問題，請使用[GitHub Discussions](https://github.com/wailsapp/wails/discussions)
- 如要回報錯誤，請使用[GitHub Issues](https://github.com/wailsapp/wails/issues)

### 如何貢獻？

請參閱[貢獻指南](/contributing/)。隨時歡迎提交錯誤修正；新增功能和變更公開行為時，則需透過[WEP（Wails Enhancement Proposal）](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md)草案 PR。您也可以選擇先在 Discord 或 GitHub Discussions 上進行非正式討論。

### 在哪裡可以找到範例？

儲存庫隨附超過60個可執行範例，涵蓋視窗、對話方塊、事件、系統匣、服務等內容：[v3/examples](https://github.com/wailsapp/wails/tree/master/v3/examples)。

## 仍有疑問嗎？

請在[Discord](https://discord.gg/JDdSxwjhGf)上提問，或[發起討論](https://github.com/wailsapp/wails/discussions)。
