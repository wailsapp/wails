---
title: "Wails v2 MacOS Beta 版"
description: "Wails 的版本資訊與公告"
authors: ["leaanthony"]
tags: ["wails","v2"]
date: "2021-11-08"
slug: "blog/wails-v2-beta-for-mac"
image: "/assets/blog-images/wails-mac.webp"
sourcePath: "blog/wails-v2-beta-for-mac.md"
---

![wails-mac 螢幕截圖](/assets/blog-images/wails-mac.webp)

今天是 Wails v2 Mac 版首個 Beta 版本發布的日子！我們花了相當長的時間才走到這一步，希望今天發布的版本能為你提供相當實用的功能。一路走來經歷了不少波折，希望能在你的協助下解決各種小問題，完善 Mac 移植版，為 v2 正式版做好準備。

所以這還不能用於正式環境嗎？對你的使用案例而言，它很可能已經可以使用，但目前仍有一些已知問題，因此請持續關注[此專案看板](https://github.com/wailsapp/wails/projects/7)。如果你願意貢獻，我們非常歡迎！

那麼，Wails v2 Mac 版相較於 v1 有哪些新功能？提示：它和 Windows Beta 版相當類似 :wink:

## 新功能

![wails-menus-mac 螢幕截圖](/assets/blog-images/wails-menus-mac.webp)

許多人都希望支援原生選單。Wails 終於滿足了這項需求。現在可以使用應用程式選單，並支援大多數原生選單功能，包括標準選單項目、核取方塊、選項群組、子選單和分隔線。

v1 收到了大量希望能更深入控制視窗本身的需求。很高興宣布，我們為此新增了專用的執行階段 API。它功能豐富，並支援多顯示器設定。我們也改善了對話方塊 API：現在，你可以使用具備豐富設定選項的現代原生對話方塊，滿足所有對話方塊需求。

### Mac 專屬選項

除了常規的應用程式選項外，Wails v2 Mac 版還提供一些 Mac 專屬功能：

- 讓你的視窗變得時髦又半透明，就像那些漂亮的 Swift 應用程式一樣！
- 可高度自訂的標題列
- 支援應用程式的 NSAppearance 選項
- 透過簡單設定自動建立「關於」選單

### 不需要封裝資產

v1 的一大痛點，是必須將整個應用程式壓縮成單一 JS 與 CSS 檔案。很高興宣布，在 v2 中完全不需要以任何方式封裝資產。想載入本機圖片嗎？使用帶有本機 src 路徑的`<img>`標籤即可。想使用很酷的字型嗎？將它複製進來，然後在 CSS 中加入其路徑。

> 哇，聽起來就像 Web 伺服器……

沒錯，它的運作方式就像 Web 伺服器，只不過它並不是。

> 那我要如何納入自己的資產？

只要將包含所有資產的單一`embed.FS`傳入應用程式設定即可。資產甚至不需要位於頂層目錄，Wails 會自行處理。

### 全新的開發體驗

現在資產不再需要封裝，因此帶來了全新的開發體驗。新的`wails dev`命令會建置並執行你的應用程式，但它不會使用`embed.FS`中的資產，而是直接從磁碟載入。

它還提供下列額外功能：

- 熱重新載入——前端資產的任何變更都會觸發應用程式前端自動重新載入
- 自動重新建置——Go 程式碼的任何變更都會觸發應用程式重新建置並重新啟動

此外，系統會在連接埠34115啟動 Web 伺服器，向任何連線的瀏覽器提供你的應用程式。所有已連線的 Web 瀏覽器都會回應系統事件，例如資產變更時的熱重新載入。

使用 Go 時，我們習慣在應用程式中處理結構體。將結構體傳送至前端並用作應用程式的狀態通常很實用。在 v1 中，這個過程需要大量手動操作，也為開發人員帶來一些負擔。很高興宣布，在 v2 中，任何以開發模式執行的應用程式，都會針對作為繫結方法輸入或輸出參數的所有結構體，自動產生 TypeScript 模型。這讓兩端的資料模型得以無縫交換。

此外，系統還會動態產生另一個 JS 模組，包裝所有繫結方法。這會為你的方法提供 JSDoc，讓 IDE 能夠提供程式碼補全與提示。當你在包裝 Go 程式碼的自動產生模組中按下 Tab 鍵時，資料模型還會自動匯入，真的很酷！

### 遠端範本

![remote-mac 螢幕截圖](/assets/blog-images/remote-mac.webp)

快速啟動並執行應用程式，一直是 Wails 專案的主要目標。專案推出時，我們嘗試涵蓋當時許多現代化框架：react、vue 和 angular。前端開發領域充滿各種不同主張、變化迅速，而且很難持續跟進！因此，我們發現基礎範本很快就會過時，造成維護上的困擾。這也意味著，我們無法為最新、最出色的技術堆疊提供很酷的現代化範本。

在 v2 中，我希望賦予社群更多能力，讓你們可以自行建立及託管範本，而不必依賴 Wails 專案。因此，現在你可以使用社群支援的範本建立專案！我希望這能激勵開發人員打造充滿活力的專案範本生態系。我真的非常期待開發者社群能創造出什麼成果！

### 原生 M1 支援

感謝[Mat Ryer](https://github.com/matryer/)的大力支持，Wails 專案現在支援 M1 原生建置：

![build-darwin-arm 螢幕截圖](/assets/blog-images/build-darwin-arm.webp)

你也可以將`darwin/amd64`指定為目標：

![build-darwin-amd 螢幕截圖](/assets/blog-images/build-darwin-amd.webp)

噢，我差點忘了……你也可以使用`darwin/universal`…… :wink:

![build-darwin-universal 螢幕截圖](/assets/blog-images/build-darwin-universal.webp)

### 交叉編譯至 Windows

由於 Windows 版 Wails v2 完全以 Go 編寫，因此無須 docker 即可將 Windows 指定為建置目標。

![build-cross-windows 螢幕截圖](/assets/blog-images/build-cross-windows.webp) bu

### WKWebView 轉譯器

V1 採用一個現已棄用的 WebView 元件。V2 使用最新的 WKWebKit 元件，因此你可以享有 Apple 提供的最新、最出色功能。

### 結語

正如我在 Windows 版本資訊中所說，Wails v2 為此專案奠定了新的基礎。此版本旨在收集大家對新方法的意見，並在正式發布前修正所有錯誤。非常歡迎您提供意見！請將任何意見提交至[v2 Beta](https://github.com/wailsapp/wails/discussions/828)討論區。

最後，我要特別感謝所有[專案贊助者](/credits/#sponsors)，包括[JetBrains](https://www.jetbrains.com?from=Wails)。他們的支持在幕後以許多方式推動著專案發展。

我很期待看到大家在此專案令人振奮的下一階段中使用 Wails 打造出哪些成果！

Lea。

附註：Linux 使用者，接下來就輪到你們了！

再附註：如果您或貴公司覺得 Wails 很實用，請考慮[贊助此專案](https://github.com/sponsors/leaanthony)。謝謝！
