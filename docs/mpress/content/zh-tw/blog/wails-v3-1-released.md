---
title: "Wails v3.1.0：更多 Go，更少建置障礙"
description: "實驗性的無 CGO 桌面建置，以及 Wake：以更簡單的方式建置、封裝及簽署 Wails 應用程式。"
authors: ["leaanthony"]
tags: ["wails", "v3", "release"]
date: "2026-11-01"
slug: "blog/wails-v3-1-released"
sourcePath: "blog/wails-v3-1-released.md"
---

Wails 本週進展繁多：無 CGO 建置、更容易的專案設定，以及更順暢的封裝與簽署流程。

Wails v3.1.0 開啟了每週功能版本的發布節奏，定期將新功能交到你手中。錯誤修正準備就緒後，仍會透過夜間修補更新持續發布。

## 本週的修正與最佳化

本週在框架、CLI、建置工具及 CI 中完成了 **100 項修正、21 項效能改進，以及四項安全性強化變更**。重點包括：混合 DPI 環境下更可靠的視窗定位、Linux 系統匣與深色模式修正、更精簡的資產提供機制，以及 MSIX 封裝期間更低的記憶體用量。[請參閱下方列出的部分變更](#heading-4)。

## 實驗性的無 CGO 建置

Purego 讓 Go 不必透過 CGo（Go 通常用來銜接 C 程式碼的橋樑），就能呼叫原生系統程式庫。我們在它剛問世時便開始探索，但維護兩套後端是一項挑戰。如今 purego 已更加成熟，而生成式 AI 也讓維護工作更容易管理，因此我們準備分享實驗性的 macOS 與 Linux 後端。

現在，你不需要 C 編譯器也能建置及交叉編譯桌面應用程式：

```sh
WAILS_PUREGO=1 wails3 build
```

建置前端資產後，你也可以直接使用 Go 建置：

```sh
CGO_ENABLED=0 go build -tags purego,production .
```

macOS 後端使用 AppKit 與 WebKit；Linux 支援 GTK4/WebKitGTK 6.0 及 GTK3。使用者仍需安裝原生執行階段程式庫，而封裝與簽署也各有個別需求。在我們收集意見回饋期間，CGo 仍是預設選項。

[實作與平台詳細資料](https://github.com/wailsapp/wails/pull/6521)

## 為什麼選擇 Wake？

Wails v2 隱藏了太多建置流程。在 v3 中，Taskfile 將一切開放出來，但多數人只想進行少量自訂，而不想維護整套建置系統。

Wake 將常用選項集中在 `wails.yaml` 中，包括應用程式、配套 CLI 工具、伺服器、目標及封裝。Task 仍是我們的核心建置系統；當你需要它的彈性時，依然可以使用。

Wake 現在支援即時重新載入、增量建置、快取輸出及平行封裝。已經有二進位檔了嗎？即使未安裝 Go 或前端工具，`--skip-build` 也能封裝它們。

使用 `WAILS_EXP_USE_WAKE=1` 啟用此功能。

### 引導式入門

啟用 Wake 後，`wails3 init -ui` 會引導你選擇後端、目標及套件，檢查所需工具，並儲存專案設定。

@image{src="/assets/screenshots/wake/application.png" alt="Wake 專案精靈顯示桌面應用程式與後端選項，其中包括實驗性的 Pure Go。" expand}

為每個目標選擇套件格式及簽署選項：

@image{src="/assets/screenshots/wake/packages.png" alt="Wake 的 Linux 套件選擇畫面，顯示 AppImage、DEB、RPM、Void 及 Arch Linux 格式，以及簽署發布輸出的選項。" expand}

下方的工具檢查顯示標準 Linux 建置：

@image{src="/assets/screenshots/wake/tools.png" alt="Wake 正在為標準 Linux 建置檢查 C 編譯器、Go、npm 及 pkg-config。" expand}

### 我的電腦可以建置哪些項目？

編譯、封裝及簽署各有不同需求。以下是你能在各種主機上執行的作業：

@table{header search sort column-separators align=["left", "centre", "centre", "centre", "left"]}
| 作業 | Windows 主機 | macOS 主機 | Linux 主機 | 所需項目 |
| 使用 purego 編譯桌面應用程式 | ✓ | ✓ | ✓ | Go，以及用於 macOS/Linux 目標的實驗性 purego 後端 |
| 建立 MSIX 套件 | ✓ | ✓ | ✓ | 已建置的應用程式；不需要 Windows SDK |
| 建立 NSIS 安裝程式 | ✓ | ✓ | ✓ | 已建置的應用程式及 `makensis` |
| 建立 macOS 應用程式及 DMG | ✓ | ✓ | ✓ | 已建置的應用程式；由 Go 處理套件建立作業 |
| 建立 DEB、RPM、Arch 及 Void 套件 | ✓ | ✓ | ✓ | 已建置的應用程式 |
| 建立 AppImage | ✓ | ✓ | ✓ | 目標 AppImage 執行階段；用於跨平台封裝的 Linux GUI 執行階段目錄樹 |
| 使用匯出的身分簽署 Windows 套件 | ✓ | ✓ | ✓ | 匯出的 RSA PFX 憑證及私密金鑰 |
| 使用匯出的身分簽署 Apple 應用程式及 DMG | ✓ | ✓ | ✓ | 匯出的 RSA PKCS#12 身分；請參閱下方的驗證說明 |
| 簽署 DEB、RPM 及 Arch 套件 | ✓ | ✓ | ✓ | 匯出的 OpenPGP 金鑰 |
| 提交 Apple 公證並附加票證 | ✓ | ✓ | ✓ | Apple 團隊 API 金鑰及經適當簽署的應用程式 |
| 使用 Apple 的 `codesign` 或鑰匙圈描述檔 | – | ✓ | – | macOS 及適當的簽署身分 |
@end

CGo 建置仍需要 C 工具鏈。會修改 DMG 的自訂掛鉤需要 macOS 及 `hdiutil`。

### 簽署與公證

Wake 支援 Windows、Apple 及 Linux 套件簽署，包括 Apple 臨時簽章。跨平台簽署使用匯出的身分；憑證存放區與鑰匙圈中的身分仍需要原生工具。使用者層級的認證資料可避免將私密金鑰放入專案設定。目前尚未支援 AppImage 與 Void 簽署。

Apple 公證可使用團隊 API 金鑰透過 Go 完成。我們已在 Linux 上測試提交及附加票證，Gatekeeper 也接受了結果。該應用程式使用 Apple 的 [`codesign`](https://developer.apple.com/library/archive/documentation/Security/Conceptual/CodeSigningGuide/Procedures/Procedures.html) 簽署；Apple 是否接受 Wake 由 Go 產生的 Developer ID 簽章，仍未經驗證。

### 移轉至 Wake

移轉功能支援全部九種內建範本。此程序不可逆：啟用 `wails.yaml` 會取代 Taskfile 路由。不支援的自訂命令及覆寫項目會遭拒絕，需進一步檢視。

[Wake 的功能、需求及移轉詳細資料](https://github.com/wailsapp/wails/pull/6531)

## 本週部分重點

- **視窗會正確出現在混合 DPI 顯示器上。** 在顯示縮放比例不同的 macOS 與 Linux 設定中，現在能正確處理螢幕邊界及位置。[#6346](https://github.com/wailsapp/wails/pull/6346)
- **Linux 系統匣與佈景主題偵測更加可靠。** 應用程式啟動後建立的系統匣不再被遺漏，工具提示也會正確填入內容，而且深色模式偵測會從啟動時起讀取桌面的外觀設定。
- **Go 事件回呼會依各接聽器維持順序。** 每個接聽器都會依事件發出順序循序執行；不同的接聽器可以並行執行，而並行發出的事件則沒有明確定義的順序。事件掛鉤會在鎖定期間建立快照，以避免資料競爭。[#6275](https://github.com/wailsapp/wails/pull/6275)、[#6276](https://github.com/wailsapp/wails/pull/6276)
- **建置快取能處理更多實際專案。** 加上引號或跳脫處理的空格不再造成 Go 輸出快取未命中，並行快取儲存也會循序執行，以防止檔案損毀。[#6495](https://github.com/wailsapp/wails/pull/6495)、[#6488](https://github.com/wailsapp/wails/pull/6488)
- **提供資產時會減少不必要的工作。** 已知的副檔名會略過 MIME 探測，而停用偵錯記錄後，也不再準備最終會被捨棄的訊息。這兩項變更都能減少請求路徑上的工作量。[#6414](https://github.com/wailsapp/wails/pull/6414)、[#6420](https://github.com/wailsapp/wails/pull/6420)
- **MSIX 封裝不再使用與二進位檔等大的緩衝區。** 可執行檔會以串流方式寫入套件，而不是完整讀入記憶體。先前複製一個具代表性的 16 MiB 檔案時，每次作業約會配置 16.8 MB。[#6472](https://github.com/wailsapp/wails/pull/6472)

## 親自試用

這兩項功能都仍處於實驗階段。請嘗試跨平台建置、封裝你的應用程式，並告訴我們你在哪些環節遇到阻礙。感謝所有參與測試、回報問題及貢獻的人。

特別感謝 Discord 社群的 beesophie；她打造了令人驚豔的建置系統 [Cook](https://github.com/SimTXNG/cook)，讓我們看見自身系統中的一些缺口。Cook 所實現的可能性幫助我們找出哪些地方能讓 Wails 建置更加容易，也確實啟發了這項工作。

謝謝你，Beesophie！
