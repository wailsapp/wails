---
title: "意見回饋"
description: "如何針對 Wails v3 提供意見回饋及回報問題"
slug: "feedback"
sourcePath: "feedback.md"
---

我們歡迎（也鼓勵）您提供意見回饋！建立新議題或討論前，請先搜尋現有的議題或討論。您可以透過以下不同方式參與貢獻：

@tabs
[錯誤]
如果發現錯誤，請使用錯誤回報範本，在 GitHub 上[建立議題](https://github.com/wailsapp/wails/issues/new/choose)。

- 請提供簡單且可重現的範例，清楚描述錯誤。如果文件未清楚說明<em>應有的</em>行為，也請在回報中註明。
- 請在回報中附上`wails3 doctor`的輸出。
- 如果錯誤行為與目前文件的說明不符，另請執行以下操作：
  - 更新`v3/examples`目錄中的現有範例，或建立清楚呈現該問題的新範例。
  - 建立一個參照該議題的[PR](https://github.com/wailsapp/wails/pulls)。


@note{type="caution"}
*請記住*，非預期的行為不一定是錯誤——它可能只是未按照您的預期運作。這種情況請使用`Suggestions`。

@end

也歡迎您在 Discord 的[#v3](https://discord.gg/bdj28QNHmT)頻道討論錯誤。

[修正]
如果您有錯誤修正或文件改善內容，請：

- 依照[貢獻指南](https://github.com/wailsapp/wails/blob/master/CONTRIBUTING.md)，在[Wails 儲存庫](https://github.com/wailsapp/wails)建立提取要求。
- 請在 PR 說明中參照所有相關議題。

[功能增強]
新功能及公開行為的變更應透過<strong>WEP（Wails Enhancement Proposal）</strong>草稿提取要求提出，而非建立功能要求議題。

- 閱讀[WEP 流程](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md)。
- 複製範本，接著建立標題為`[WEP] <title>`的草稿 PR，其中僅包含 WEP 及其支援資料。
- 您可以先在[GitHub Discussions](https://github.com/wailsapp/wails/discussions)或 Discord 的[#v3](https://discord.gg/bdj28QNHmT)頻道非正式討論構想，但維護者必須收到 WEP PR 才能作出決定。

[投贊成票]
- 請在 GitHub 上使用 :thumbsup: 回應，表示支持錯誤議題、WEP 和討論。
- 請<em>不要</em>只新增「+1」或「我也是」之類的留言。
- 如果您有實質內容可以補充，請留言，例如「這個錯誤也會影響 ARM 組建」或「另一種方法是……」。

@end

您可以在[這裡](https://github.com/orgs/wailsapp/projects/6)找到已知問題和進行中的工作。
