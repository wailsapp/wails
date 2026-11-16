---
title: "意見回饋"
description: "如何提供意見回饋及回報 Wails v3的問題"
slug: "feedback"
sourcePath: "feedback.md"
---

我們歡迎（也鼓勵）您提供意見回饋！建立新的問題或討論之前，請先搜尋現有的問題或討論。您可以透過下列方式參與貢獻：

@tabs
[錯誤]
如果發現錯誤，請使用錯誤回報範本在 GitHub 上[建立問題](https://github.com/wailsapp/wails/issues/new/choose)。

- 請提供簡單且可重現的範例，清楚描述錯誤。如果文件未清楚說明*應該*發生什麼情況，也請在回報中註明。
- 請在回報中附上`wails3 doctor`的輸出。
- 如果該錯誤是與目前文件不符的行為，另請執行下列步驟：
  - 更新`v3/examples`目錄中的現有範例，或建立一個能清楚呈現問題的新範例。
  - 建立一個引用該問題的[PR](https://github.com/wailsapp/wails/pulls)。


@note{type="caution"}
*請記住*，非預期行為不一定是錯誤，也可能只是行為不符合您的預期。這種情況請使用`Suggestions`。

@end

也歡迎在 Discord 的[#v3](https://discord.gg/bdj28QNHmT)頻道中討論錯誤。

[修正]
如果您有錯誤修正或文件改進，請：

- 依照[貢獻指南](https://github.com/wailsapp/wails/blob/master/CONTRIBUTING.md)，在[Wails 儲存庫](https://github.com/wailsapp/wails)中建立提取要求。
- 請在 PR 說明中引用所有相關問題。

[功能增強]
新功能及公開行為的變更應透過<strong>WEP（Wails Enhancement Proposal）</strong>草稿提取要求提出，而不是建立功能要求問題。

- 閱讀[WEP 流程](https://github.com/wailsapp/wails/blob/v3.0.0/v3/wep/README.md)。
- 複製範本，然後建立標題為`[WEP] <title>`的草稿 PR，其中僅包含 WEP 及支援資料。
- 您可以先在[GitHub Discussions](https://github.com/wailsapp/wails/discussions)或 Discord 的[#v3](https://discord.gg/bdj28QNHmT)頻道中非正式討論想法，但若要由維護者作出決定，仍必須建立 WEP PR。

[按讚支持]
- 請在 GitHub 上使用 :thumbsup: 回應，表達對錯誤、WEP 和討論的支持。
- 請*不要*只留下「+1」或「我也是」之類的留言。
- 如果您有實質內容可以補充，例如「這個錯誤也會影響 ARM 組建」或「另一種做法是……」，請務必留言。

@end

您可以在[這裡](https://github.com/orgs/wailsapp/projects/6)查看已知問題和進行中的工作。
