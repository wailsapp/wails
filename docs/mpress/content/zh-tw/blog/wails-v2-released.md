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

## 它來了！

今天，[Wails](https://wails.io) v2 正式發布。距離第一個 18 Alpha 版本大約已有 v2 個月，距離第一個 Beta 版本則約有一年。我由衷感謝所有參與推動專案發展的人。

之所以花了這麼長時間，部分原因是希望先讓專案達到某種完整程度，再正式稱之為 v2。事實上，永遠不會有標記版本的完美時機——總會有尚未解決的問題，或是想再塞進去的「最後一項」功能。然而，為一個尚不完美的主要版本加上標籤，確實能為專案使用者提供一些穩定性，也讓開發者得以重新整頓、再出發。

這個版本遠遠超出我原先的期待。希望它能帶給你如同我們開發它時所感受到的喜悅。

## Wails *是*什麼？

如果你不熟悉 Wails，它是一個讓 Go 程式設計師能使用熟悉的 Web 技術，為 Go 程式打造豐富前端的專案。它是 Electron 的輕量級 Go 替代方案。你可以在[官方網站](https://v2.wails.io/docs/introduction)找到更多資訊。

## 有哪些新功能？

v2 版本是本專案的一大躍進，解決了 v1 的許多痛點。如果你還沒讀過關於 [macOS](/blog/wails-v2-beta-for-mac/)、[Windows](/blog/wails-v2-beta-for-windows/) 或 [Linux](/blog/wails-v2-beta-for-linux/) Beta 版本的部落格文章，我建議你閱讀；其中更詳細地介紹了所有重大變更。簡而言之：

- Windows 採用 Webview2 元件，支援現代 Web 標準與偵錯功能。
- Windows 支援[深色／淺色主題](https://v2.wails.io/docs/reference/options#theme)及[自訂主題](https://v2.wails.io/docs/reference/options#customtheme)。
- Windows 現在不再需要 CGO。
- 內建支援 Svelte、Vue、React、Preact、Lit 與 Vanilla 專案範本。
- 整合[Vite](https://vitejs.dev/)，為應用程式提供可熱重載的開發環境。
- 原生應用程式[選單](https://v2.wails.io/docs/guides/application-development#application-menu)與[對話方塊](https://v2.wails.io/docs/reference/runtime/dialog)。
- 為[Windows](https://v2.wails.io/docs/reference/options#windowistranslucent)與[macOS](https://v2.wails.io/docs/reference/options#windowistranslucent-1)提供原生視窗半透明效果，並支援 Mica 與 Acrylic 背景材質。
- 輕鬆產生用於 Windows 部署的[NSIS 安裝程式](https://v2.wails.io/docs/guides/windows-installer)。
- 功能豐富的[執行階段程式庫](https://v2.wails.io/docs/reference/runtime/intro)，提供視窗操作、事件處理、對話方塊、選單及記錄功能的公用方法。
- 支援使用[garble](https://github.com/burrowers/garble)對應用程式進行[混淆](https://v2.wails.io/docs/guides/obfuscated)。
- 支援使用[UPX](https://upx.github.io/)壓縮應用程式。
- 自動為 Go 結構產生 TypeScript。更多資訊請見[此處](https://v2.wails.io/docs/howdoesitwork#calling-bound-go-methods)。
- 無論在哪個平台，都不需要隨應用程式一併提供額外的程式庫或 DLL。
- 不必封裝前端資產，只需像開發其他 Web 應用程式一樣開發你的應用程式。

## 致謝

一路走到 v2，大家付出了極大的努力。從最初的 Alpha 版本到今天正式發布，共有 89 位貢獻者完成了約 2200 次提交；此外，還有更多人提供翻譯、測試、意見回饋，並在討論區與問題追蹤器上給予協助。我對每一位都感激不盡。我也要特別感謝所有提供指導、建議與回饋的專案贊助者。你們所做的一切都非常值得感謝。

我想特別提及幾位貢獻者：

首先，**非常**感謝 [@stffabi](https://github.com/stffabi)。他做出了許多讓我們所有人受益的貢獻，也協助處理了許多問題。他提供了多項關鍵功能，例如支援外部開發伺服器，讓我們得以運用[Vite](https://vitejs.dev/)的強大能力，徹底改變了我們的開發模式。可以毫不誇張地說，如果沒有他的[傑出貢獻](https://github.com/wailsapp/wails/commits?author=stffabi&since=2020-01-04)，Wails v2 版本絕不會如此令人振奮。萬分感謝你，@stffabi！

我也要大力感謝 [@misitebao](https://github.com/misitebao)。他不辭辛勞地維護網站、提供中文翻譯、管理 Crowdin，並協助新加入的翻譯人員快速上手。這是一項極為重要的工作，我由衷感謝他為此投入的所有時間與心力！你太厲害了！

最後同樣重要的是，非常感謝 Mat Ryer 在 v2 開發期間提供建議與支援。使用 v2 的早期 Alpha 版本一起開發 xBar，不僅有助於確立 v2 的發展方向，也讓我了解早期版本中的一些設計缺陷。我很高興地宣布，從今天起，我們將開始把 xBar 移植到 Wails v2，它也會成為本專案的旗艦應用程式。謝謝你，Mat！

## 經驗與教訓

一路走到 v2 的過程中，我們學到了許多經驗與教訓，這些都將影響未來的開發方向。

## 更小、更快、更聚焦的版本

在開發 v2 的過程中，許多功能與錯誤修正都是臨時開發的。這不僅拉長了發布週期，也提高了偵錯難度。今後，我們會更頻繁地發布版本，且每個版本包含的功能數量會減少。每次發布都會包括文件更新與全面測試。希望這些更小、更快、更聚焦的版本能減少迴歸問題，並提升文件品質。

## 鼓勵參與

剛開始這個專案時，我想立即幫助每一位遇到問題的人。我把問題視為「自己的事」，希望盡快解決它們。這種做法無法持續，最終也不利於專案的長久發展。今後，我會留出更多空間，讓其他人參與回答問題與分流處理議題。如果能有一些工具協助這項工作會很好；若你有任何建議，歡迎到[這裡](https://github.com/wailsapp/wails/discussions/1855)參與討論。

## 學會說不

參與開放原始碼專案的人越多，新增功能的要求就越多，而這些功能對大多數人而言未必實用。開發與偵錯這些功能一開始就需要投入時間，之後還會持續產生維護成本。在這方面，我自己最需要反省：我常想著「煮沸整片海洋」，而不是先提供最小可行的功能。今後，面對新增核心功能的要求，我們需要更常說「不」，並把心力集中在讓開發者能自行提供這些功能的方法上。我們正認真研究以外掛程式處理這類情境。如此一來，任何人都能依需求擴充專案，也能以簡便的方式為專案作出貢獻。

## 展望未來

我們已經在考慮於下一個主要開發週期中，為 Wails 加入許多核心功能。[開發藍圖](https://github.com/wailsapp/wails/discussions/1484)中充滿了有趣的構想，我也迫不及待想著手實作。其中一項呼聲很高的需求是支援多視窗。這很棘手；為了妥善實作，我們可能需要考慮提供另一套 API，因為目前的 API 並非以此為設計前提。根據一些初步構想與回饋，我想你會喜歡我們規劃的方向。

就個人而言，我對讓 Wails 應用程式在行動裝置上執行的前景感到非常興奮。我們已經有一個示範專案，證明 Wails 應用程式可以在 Android 上執行，因此我非常期待探索未來能把它發展到什麼程度！

最後，我想談談功能一致性。長期以來，我們一直秉持一項核心原則：除非某項功能能獲得完整的跨平台支援，否則不會將其加入專案。儘管到目前為止，這項原則（大致上）都能落實，卻也嚴重阻礙了專案發布新功能。今後，我們將採取略有不同的方式：任何無法立即在所有平台上發布的新功能，都會透過實驗性設定或 API 發布。這可讓特定平台上的早期採用者試用該功能並提供意見，這些回饋將納入功能的最終設計。當然，這表示在該功能獲得所有可支援平台的完整支援之前，無法保證 API 的穩定性，但至少不會再阻礙開發進度。

## 結語

我對我們在V2版本中所取得的成果深感自豪。看到大家至今已能使用測試版打造出各種成果，實在令人驚喜，其中包括[Varly](https://varly.app/)、[Surge](https://getsurge.io/)和[October](https://october.utf9k.net/)等優質應用程式。我鼓勵大家去看看。

這個版本是許多貢獻者辛勤努力的成果。雖然任何人都能免費下載及使用，但它並非毫無成本。請不要誤會，這個專案付出了相當可觀的代價。這不僅包括我和每一位貢獻者投入的時間，也包括我們每個人因而無法陪伴親友與家人的時間。因此，我由衷感謝大家為實現這個專案所投入的每一秒。貢獻者越多，就越能分擔這些工作，也越能攜手完成更多事情。我想鼓勵大家挑一件自己能貢獻的事，無論是確認他人回報的錯誤、提出修正建議、修改文件，或是協助需要幫助的人。這些看似微小的事情，都能帶來極大的影響！如果在邁向v3的故事中也有你的一份，那就太棒了。

祝你使用愉快！

&dash; Lea

附註：如果你或你的公司認為 Wails 很實用，請考慮[贊助本專案](https://github.com/sponsors/leaanthony)。謝謝！
