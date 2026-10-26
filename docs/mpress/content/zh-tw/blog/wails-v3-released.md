---
title: "Wails v3.0.0 正式發布"
description: "Wails v3.0.0 正式登場：支援多視窗、系統匣、原生選單、服務與產生的繫結、更新程式、行動裝置，以及全新的文件網站。"
authors: ["leaanthony"]
tags: ["wails","v3","release"]
date: "2026-10-26"
slug: "blog/wails-v3-released"
image: "/assets/videos/wails-v3-release-clean-poster.webp"
sourcePath: "blog/wails-v3-released.md"
---

@video{src="/assets/videos/wails-v3-release-clean-av1.mp4" poster="/assets/videos/wails-v3-release-clean-poster.webp" title="Wails v3.0.0 發布影片——建議以全螢幕觀看" width="1920" height="1080"}

## 正式登場！

今天，我要發布 Wails v3.0.0。

老實說，這一路走來很漫長。第一個 v3 alpha 標籤於 2023 年 1 月發布；自那之後，許多人一路參與建置、測試、找出問題並加以修正。我由衷感謝每一位參與者。如果你想了解這段歷程最坦率的版本，我在[v3 beta 公告](/blog/wails-v3-beta/)中寫過。

## Wails*究竟是*什麼？

如果你剛接觸 Wails：Wails 讓 Go 程式設計師能使用熟悉的 Web 前端工具建置桌面應用程式，並採用各平台的原生 WebView，而非內嵌瀏覽器。你可以用 Go 撰寫後端，用 HTML、CSS 與 JavaScript（或任何你偏好的框架）製作使用者介面，再由 Wails 將這些內容封裝成單一原生應用程式。

## 為什麼要推出 v3？

Wails v2 讓使用現代 Web 前端建置 Go 應用程式變得簡單直接，也一直為許多應用程式提供良好支援。但它採用單一視窗、由 context 驅動的執行階段，使某些常見的桌面應用開發工作變得比應有的更困難。

v3 採用不同的模型。應用程式、視窗、服務、事件與平台功能都是明確的物件。你會建立應用程式、註冊服務、建立視窗，並與擁有所需行為的物件互動。隨著應用程式成長，這種模型更容易理解；多視窗等功能也是模型中的常規部分，而非權宜之計。

## v3.0.0 包含哪些內容

- <strong>多視窗。</strong>可在執行階段建立、管理及關閉視窗，每個視窗都有各自的設定與生命週期。
- <strong>系統匣。</strong>讓應用程式常駐系統匣，支援選單、亮色與暗色圖示、macOS 樣板圖示，並可選擇將視窗附加至系統匣圖示。
- <strong>原生選單與對話方塊。</strong>支援各平台的應用程式選單、內容選單及原生對話方塊。
- <strong>服務與產生的繫結。</strong>後端仍是一般的 Go 程式碼。繫結會透過靜態原始碼分析產生，因此只需一個 `wails3 generate bindings` 命令，註解與具意義的參數名稱就能一併帶到前端。
- <strong>看得見的建置系統。</strong>專案採用以 [Taskfile](https://taskfile.dev) 為基礎的建置方式，讓你能檢查、擴充及偵錯，而不是依賴一個不透明的命令。
- <strong>跨平台。</strong>支援 Windows（amd64 與 arm64）、macOS（Intel 與 Apple Silicon）及 Linux（amd64 與 arm64）；Linux 預設使用 GTK4 與 WebKitGTK 6.0。在 v3.1 之前，仍可透過 `-tags gtk3` 使用舊版 GTK3 技術堆疊。
- <strong>iOS 與 Android。</strong>相同的 `main.go` 和前端也能在行動裝置上執行。文件仍將行動裝置支援標示為實驗性功能，因此應將其視為可供探索及提供意見回饋的功能，而非穩定性承諾。
- <strong>使用 Go 1.26 建置。</strong>v3 模組使用 Go 1.26 建置。
- <strong>應用程式內更新程式。</strong>支援可插拔的提供者、密碼學驗證、不可分割的置換作業，以及可自訂佈景主題或替換的預設使用者介面。請參閱[更新程式指南](/guides/updater/)。
- <strong>封裝與簽署。</strong>提供各平台的[程式碼簽署](/guides/build/signing/)、[安裝程式](/guides/installers/)及[MSIX 封裝](/guides/build/msix/)指南。
- <strong>伺服器建置。</strong>無須原生視窗，即可將相同的應用程式與服務作為 HTTP 伺服器執行。

## Beta 版發布以來的重點

Beta 版文章於 8 月 2 日發布，此後變更記錄新增了許多內容。以下是我特別想指出的幾項：

- <strong>全新的文件網站。</strong>文件網站現在使用我們自行開發的文件工具 [M-Press](https://m-press.me)，翻譯也已完成並從社群移轉至新網站。英文文件仍是唯一的權威依據。
- <strong>更強韌的 Windows 使用體驗。</strong>現在 WebView2 處理程序失敗時，Wails 能夠復原，不會只留下空白視窗；同時也修正了多項 WebView2 COM 生命週期問題。
- <strong>配合鍵盤配置的全域快速鍵。</strong>在 macOS 上，字母快速鍵現在會依循 AZERTY、QWERTZ 及 Dvorak 等配置中按鍵標示的字母。
- <strong>強化更新程式。</strong>修正了 Linux 上跨檔案系統置換二進位檔，以及 AppImage 更新等問題。
- **設定與工具。**`wails3 setup`精靈現已完成，而 `wails3 dev` 現在會在應用程式結束時停止其背景處理程序。
- <strong>macOS 原生 Liquid Glass。</strong>支援 macOS 26 及更新版本的 Liquid Glass 效果，舊版則會使用備援效果。

## 從 v2 移轉

v3 是新的主要版本，因此移轉需要實際移植，而不只是更新版本號。主要變更包括應用程式與視窗的生命週期、以服務取代與 context 綁定的繫結、以直接的應用程式與視窗 API 取代 v2 runtime 套件，以及重新產生前端繫結。

[v2 到 v3 的移轉指南](/migration/v2-to-v3/)以移轉前後的範例逐一說明所有變更。請有計畫地進行移植、在過程中持續測試，並在 v3 應用程式準備就緒前繼續保留 v2。

<strong>Wails v2 不會消失。</strong>它仍會持續獲得重大更新。新專案請使用 v3.0.0。兩個 CLI 是不同的二進位檔（`wails` 與 `wails3`），模組匯入路徑也不同，因此可以同時安裝在同一部機器上。

## 安裝

使用下列命令安裝 v3.0.0 CLI：

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0
```

接著執行引導式設定精靈；它會檢查你的環境，並協助處理 Wails 所需的相依套件：

```sh
wails3 setup
```

![深色模式下的 Wails v3 設定精靈](/assets/screenshots/wails3-setup-wizard-dark.png)

然後建立第一個專案：

```sh
wails3 init
```

- 文件：<https://v3.wails.io/>
- 移轉指南：<https://v3.wails.io/migration/v2-to-v3/>
- 版本資訊：GitHub 上的[Wails v3.0.0](https://github.com/wailsapp/wails/releases/tag/v3.0.0)

## 支援與版本發布

Wails 的成長遠遠超出我的預期，而議題與提取要求的數量也一直令人難以負荷。過去發布版本有時缺乏固定規劃，修正與貢獻也等了太久才送到使用者手中。我們已建立系統，將分流處理與發布準備的更多基礎工作自動化，以便更快回應問題，並投入更多時間推出改進與新功能。

在 v3.0.0 之後，我們將改為發布帶有日期標籤（例如 `v3.1.0-nightly.20261026`）的**每夜預發布版本**，同時每週發布**次要版本**，例如 `v3.1.0` 與 `v3.2.0`。每夜版本能讓近期修正更快提供使用，而每週版本則會以固定節奏提供經過測試的更新。

如果你發現可重現的錯誤，請附上 `wails3 doctor` 的輸出結果與最小範例進行回報。如果想提議新功能或變更公開行為，請先建立 WEP（Wails Enhancement Proposal，Wails 增強提案）草稿 PR，而非功能要求議題。

## 翻譯

衷心感謝每一位協助翻譯 Wails v2 文件的人。你們付出時間，讓開發人員能以自己的語言閱讀文件，這對許多人而言意義重大。我們由衷感激。

我們發現 Crowdin 特別難以管理，因此已停止使用。從 v3 起，文件會自動翻譯，譯文與英文頁面一同存放在儲存庫中。

如果你希望新增我們尚未提供的語言譯本，或發現現有譯本中有錯誤，請建立 issue 告訴我們。我們非常期待收到你的意見，也感謝每一份回報。

## 致謝

Wails v3 得以問世，要歸功於所有測試未完成版本、回報棘手錯誤、翻譯文件、回答問題及貢獻程式碼的人。謝謝你們。完整的貢獻者名單列在[致謝頁面](/credits/)上，我誠摯邀請你前往瀏覽。

首先，要特別感謝 Wails 維護者[Atterpac](https://github.com/atterpac)。他在開發工具、Linux 桌面整合與社群支援方面的付出，協助 v3 走到今天。

特別感謝[Fabio Massaioli（fbbdev）](https://github.com/fbbdev)打造出如此出色的靜態原始碼分析器與繫結產生器。它是 v3 得以運作的重要核心：一般的 Go 服務會轉換成 JavaScript 與 TypeScript API，型別和文件也會一併帶入前端。

也要感謝[Travis McLane](https://github.com/tmclane)早期在平台方面投入大量工作；[Stffabi](https://github.com/stffabi)奠定資產伺服器與原生平台的基礎；[Zach Botterman（popaprozac）](https://github.com/popaprozac)實作通知、Dock 整合及其他桌面功能；[Ian VanSchooten（IanVS）](https://github.com/IanVS)改善強型別事件與執行階段；以及[Lyimmi](https://github.com/Lyimmi)在 Linux 方面的工作。這些貢獻，以及每一位協助建構與完善專案的人所付出的成果，至今仍是 v3 的一部分。

衷心感謝在這段漫長轉型期間一路支持 Wails 的贊助者。你們的支持讓我們有時間投入架構、工具、測試及文件工作，成就了今天的成果。也特別感謝 JetBrains 和 ZSA，以及透過程式碼審查支持本專案的[CodeRabbit](https://www.coderabbit.ai/oss)。

最後，特別感謝 Valerio Tesei 多年來對本專案的全力支持。

## 結語

我們已有許多工作蓄勢待發。科技領域正出現一些令人振奮的發展，我們也已準備好探索新的典範，為 Wails 帶來全新的功能。不過，我們首先必須讓 v3 成為你可以信賴的穩定基礎。把這件事做好，比推出下一階段的內容更重要。

我真的為 v3 感到自豪。我希望使用它進行開發，也能讓你感受到如同我們打造它時一樣的樂趣。試用看看、動手做些東西，並告訴我們你的想法。歡迎到 [Discord](https://discord.gg/JDdSxwjhGf)、[GitHub Discussions](https://github.com/wailsapp/wails/discussions) 或 [subreddit](https://www.reddit.com/r/wails/) 向我們打聲招呼。

祝你使用愉快！

&dash;&nbsp;Lea

附註：如果你或你的公司認為 Wails 很實用，請考慮[贊助本專案](https://github.com/sponsors/leaanthony)。謝謝！
