---
title: "Wails v2 正式發布"
description: "Wails 的版本資訊與公告"
authors: ["leaanthony"]
tags: ["wails","v2"]
date: "2022-09-22"
slug: "blog/wails-v2-released"
image: "/assets/blog-images/montage.png"
sourcePath: "blog/wails-v2-released.md"
---

![蒙太奇截圖](/assets/blog-images/montage.png)

## 正式登場！

今天是[Wails](https://wails.io) v2 正式發布的日子。距離第一個 v2 Alpha 版發布至今約有18個月，距離第一個 Beta 版發布則約有一年。衷心感謝每一位參與推動本專案發展的人。

之所以耗時如此之久，部分原因是我們希望先達到某種程度的完整性，再正式將它稱為 v2。事實上，標記版本永遠沒有完美的時機——總會有尚未解決的問題，或是想再塞進去的「最後一項」功能。然而，即使標記的是不完美的主要版本，仍能為專案使用者提供一定的穩定性，也讓開發者得以重新整頓再出發。

這個版本遠遠超出我原先的期待。希望它能帶給你如同我們開發它時所感受到的那般喜悅。

## Wails<em>是</em>什麼？

如果你不熟悉 Wails，它是一個讓 Go 程式設計師能使用熟悉的 Web 技術，為 Go 程式打造功能豐富前端的專案。它是以 Go 打造的輕量級 Electron 替代方案。如需更多資訊，請造訪[官方網站](https://wails.io/docs/introduction)。

## 有哪些新功能？

v2 版本是本專案的一次重大躍進，解決了 v1 的許多痛點。如果你尚未閱讀關於[macOS](/blog/wails-v2-beta-for-mac/)、[Windows](/blog/wails-v2-beta-for-windows/)或[Linux](/blog/wails-v2-beta-for-linux/) Beta 版的任何部落格文章，建議你閱讀；其中更詳細地介紹了所有主要變更。摘要如下：

- Windows 採用 Webview2 元件，支援現代 Web 標準與偵錯功能。
- Windows 支援[深色／淺色佈景主題](https://wails.io/docs/reference/options#theme)及[自訂佈景主題](https://wails.io/docs/reference/options#customtheme)。
- Windows 現在不再需要 CGO。
- 內建支援 Svelte、Vue、React、Preact、Lit 與 Vanilla 專案範本。
- 整合[Vite](https://vitejs.dev/)，為你的應用程式提供支援熱重新載入的開發環境。
- 原生應用程式[選單](https://wails.io/docs/guides/application-development#application-menu)與[對話方塊](https://wails.io/docs/reference/runtime/dialog)。
- [Windows](https://wails.io/docs/reference/options#windowistranslucent)與[macOS](https://wails.io/docs/reference/options#windowistranslucent-1)的原生視窗半透明效果。支援 Mica 與 Acrylic 背景材質。
- 可輕鬆產生用於 Windows 部署的[NSIS 安裝程式](https://wails.io/docs/guides/windows-installer)。
- 功能豐富的[執行階段程式庫](https://wails.io/docs/reference/runtime/intro)，提供視窗操作、事件處理、對話方塊、選單與記錄等公用方法。
- 支援使用[garble](https://github.com/burrowers/garble)對應用程式進行[混淆](https://wails.io/docs/guides/obfuscated)。
- 支援使用[UPX](https://upx.github.io/)壓縮應用程式。
- 可從 Go 結構自動產生 TypeScript。更多資訊請見[此處](https://wails.io/docs/howdoesitwork#calling-bound-go-methods)。
- 無論在哪個平台，都不需要隨應用程式附帶任何額外的程式庫或 DLL。
- 不需要打包前端資產。只需像開發任何其他 Web 應用程式一樣開發即可。

## 致謝

一路走到 v2，我們付出了巨大的努力。從最初的 alpha 版本到今天正式發布，共有 89 位貢獻者提交了約 2200 次 commit；此外，還有非常多人提供了翻譯、測試、意見回饋，並在 討論論壇和問題追蹤器上協助大家。我由衷感謝你們每一位。我也要特別感謝所有 專案贊助者提供指導、建議和意見回饋。你們所做的一切都讓我深深感激。

我想特別提到幾位：

首先，要向[@stffabi](https://github.com/stffabi)致上<strong>萬分</strong>感謝。他做出了許多讓所有人受益的貢獻，也為眾多問題提供了大量支援。他提供了多項關鍵功能，例如外部開發伺服器支援；這讓我們得以運用[Vite](https://vitejs.dev/)的強大能力，徹底改變了我們提供的開發模式。可以毫不誇張地說，如果沒有他的[卓越貢獻](https://github.com/wailsapp/wails/commits?author=stffabi&since=2020-01-04)，Wails v2 會是個遜色許多的版本。真的非常感謝你，@stffabi！

我也要大力感謝[@misitebao](https://github.com/misitebao)。他不辭辛勞地維護網站、提供中文翻譯、管理 Crowdin，並協助新譯者快速上手。這是一項極其重要的工作，我非常感謝他為此投入的所有時間與心力！你太厲害了！

最後同樣重要的是，非常感謝 Mat Ryer 在 v2 開發期間提供建議與支援。我們曾使用 v2 的早期 Alpha 版共同開發 xBar，這不僅有助於確立 v2 的發展方向，也讓我了解早期版本中的一些設計缺陷。我很高興宣布，自今天起，我們將開始把 xBar 移植到 Wails v2，它也將成為本專案的旗艦應用程式。謝謝你，Mat！

## 經驗與教訓

在邁向 v2 的過程中，我們汲取了多項經驗與教訓，這些都將影響未來的開發方向。

## 規模更小、速度更快、目標更明確的版本

在開發 v2 的過程中，許多功能與錯誤修正都是臨時開發的。這導致發布週期更長，也更難偵錯。未來，我們將更頻繁地發布版本，每個版本納入的功能數量也會減少。每次發布都會包含文件更新與完整測試。希望這些規模更小、速度更快且目標明確的版本，能減少迴歸問題並提升文件品質。

## 鼓勵參與

剛開始這個專案時，我希望能立即幫助每一位遇到問題的人。我把每個問題都視為「自己的事」，希望盡快解決。然而，這種做法無法長久維持，最終反而不利於專案的長期發展。未來，我會留出更多空間，讓其他人參與回答問題及分流處理議題。若能有一些工具協助處理這些工作會很有幫助，因此如果你有任何建議，歡迎在[此處](https://github.com/wailsapp/wails/discussions/1855)參與討論。

## 學會說「不」

參與開放原始碼專案的人越多，新增功能的要求也會越多，而這些功能可能對大多數人有用，也可能沒有用。這些功能不僅需要投入前期時間進行開發與偵錯，往後也會持續產生維護成本。在這方面，我自己最容易犯這個毛病，常常想要一次做到面面俱到，而不是只提供最小可行功能。未來，對於新增核心功能，我們需要更常說「不」，並把精力集中在如何讓開發者自行提供這些功能。我們正認真考慮以外掛程式來處理這種情境。這將讓任何人都能依自身需求擴充本專案，同時也提供一種簡便的專案貢獻方式。

## 展望未來

我們已經在研究許多核心功能，準備在下一個主要開發週期加入 Wails。[開發藍圖](https://github.com/wailsapp/wails/discussions/1484)中充滿了有趣的構想，我迫不及待想開始著手實現。其中呼聲很高的一項需求是支援多視窗。這項功能相當棘手；為了妥善實現它，我們可能需要考慮提供替代 API，因為目前的 API 在設計時並未將此需求納入考量。根據一些初步構想與意見回饋，我想你會喜歡我們正在探索的方向。

我個人非常期待讓 Wails 應用程式在行動裝置上執行。我們已經有一個示範專案，證明 Wails 應用程式可以在 Android 上執行，因此我非常想探索這方面還能有什麼發展！

最後，我想談談功能一致性。長久以來，我們的一項核心原則是：除非某項功能能獲得完整的跨平台支援，否則不會將它加入專案。儘管目前為止已證明這項原則（大致上）可行，但它確實阻礙了專案發布新功能。今後，我們將採用稍微不同的做法：任何無法立即在所有平台上發布的新功能，都會透過實驗性設定或 API 發布。如此一來，特定平台上的早期採用者便能試用該功能並提供意見，這些回饋將納入功能的最終設計。當然，這表示在該功能於所有可支援的平台上獲得完整支援之前，API 的穩定性不受任何保證；但至少能讓開發工作不再受阻。

## 最後的話

我真的很以我們在 V2 版本中取得的成果為榮。看到大家至今已經能使用 Beta 版本打造出令人驚豔的成果，實在令人振奮，其中包括 [Varly](https://varly.app/)、[Surge](https://getsurge.io/) 和 [October](https://october.utf9k.net/) 這些優質應用程式。歡迎你去看看。

這個版本是許多貢獻者辛勤努力的成果。雖然任何人都能免費下載及使用，但它並非毫無成本地誕生。請不要誤會，這個專案付出了相當可觀的代價。付出的不只是我的時間與每一位貢獻者的時間，還包括這些人因投入專案而無法陪伴親友的代價。因此，我由衷感謝大家為實現這個專案投入的每一分每一秒。貢獻者越多，我們就越能分攤這份心力，也越能攜手完成更多成果。我想鼓勵大家各自挑選一件力所能及的事來貢獻，無論是協助確認他人回報的錯誤、提出修正建議、修改文件，或幫助有需要的人。這些微小的付出都能產生極大的影響！如果在邁向 v3 的故事中也有你的一份力量，那就太棒了。

祝你使用愉快！

&dash; Lea

附註：如果你或你的公司認為 Wails 很實用，請考慮[贊助本專案](https://github.com/sponsors/leaanthony)。謝謝！
