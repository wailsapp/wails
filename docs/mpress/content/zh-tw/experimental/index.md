---
title: "實驗性功能"
description: "Wails v3 中仍在開發的實驗功能專區——說明這些功能是什麼、為何存在，以及如何提供意見回饋。"
slug: "experimental"
sourcePath: "experimental/index.md"
---

@note{type="caution" title="此處皆為實驗功能"}
本節的所有內容，顧名思義，都是實驗性功能。這些功能必須主動啟用，預設為停用，且在不同版本之間可能改變形式、重新命名，甚至完全移除。除非已準備好因應頻繁變動，否則請勿將任何關鍵功能建立在這些實驗上。

@end

## 「實驗性」的意義

Wails 公開推出實驗性功能。實驗性功能是我們認為值得提早交到你手中的構想，但尚未完全承諾會正式採用。我們之所以將它公開，*是因為*希望先從實際使用情況中學習，再決定它是否會成為 Wails 永久且受支援的一部分。

這表示本節的所有功能都有以下幾項特性：

- <strong>必須主動啟用。</strong>實驗性功能絕不會改變`wails3`的預設行為。你必須刻意啟用它們（通常使用環境變數或建置旗標）；停用時，你現有的工作流程完全不受影響。
- <strong>它可能不會保留下來。</strong>有些實驗性功能會發展成穩定功能，有些則會大幅重製到面目全非，或遭到捨棄。我們寧可公開嘗試並快速學習，也不願只推出已經完全確定的功能。
- <strong>API 尚未定案。</strong>在實驗性功能逐漸成形的過程中，其名稱、旗標、預設值及行為都可能隨版本變更。版本資訊會特別說明這些變更，但請勿期待它具備穩定功能所提供的穩定性保證。

## 我們需要你的意見回饋

這是最重要的部分。實驗性功能能否延續，取決於實際使用者提供的意見回饋。如果你試用了其中一項功能，我們真心希望瞭解：

- 它是否適用於你的專案？在哪些地方無法正常運作？
- 它是否更快、更清楚或更好用，還是根本不值得切換？
- 必須滿足哪些條件，你才會預設使用它？

最有幫助的意見回饋應具體說明：你執行了什麼、預期會發生什麼，以及實際發生了什麼。每項實驗在 GitHub Discussions 的**實驗性功能**分類中都有專屬討論串：

@container{display="grid" columns="2" gap="1rem"}
@linkcard{title="實驗性功能討論" href="https://github.com/wailsapp/wails/discussions/categories/experiments" description="尋找你所用實驗功能的討論串，告訴我們使用結果、哪些地方發生問題，或還缺少什麼。"}
@end

## 目前的實驗性功能

@container{display="grid" columns="2" gap="1rem"}
@linkcard{title="Wake" href="/experimental/wake/" description="試用原生 YAML 建置系統：透過單一精靈建立並設定專案、建置及封裝專案，或檢視從 Taskfile 單向遷移的內容。"}
@linkcard{title="無 CGO 建置" href="/experimental/#cgo-free-builds-purego" description="使用 CGO_ENABLED=0 建置 macOS 與 Linux 應用程式：不需要 C 工具鏈，也能在不使用 Docker 的情況下交叉編譯。透過 WAILS_PUREGO=1 或 purego 建置標籤啟用。"}
@linkcard{title="LLM 控制（MCP）" href="/guides/mcp-service/" description="內建的 Model Context Protocol 伺服器，讓 LLM 代理程式能檢查、測試及操控執行中的 Wails 應用程式；不需要使用者程式碼，並可透過建置標籤啟用。"}
@end

## 無 CGO 建置（purego）

@note{type="caution" title="實驗性功能"}
無 CGO 後端必須透過`purego`建置標籤（或`WAILS_PUREGO=1`）主動啟用，**並非**預設後端。若未啟用，macOS 與 Linux 建置會和以往完全相同，使用 cgo 後端。其行為可能隨版本變更。

@end

### 功能說明

在 macOS 和 Linux 上，Wails 通常透過 cgo 與系統的 UI 程式庫通訊（macOS 使用 AppKit 和 WKWebView，Linux 使用 GTK 和 WebKitGTK）。無 CGO 後端使用[purego](https://github.com/ebitengine/purego)以純 Go 發出相同呼叫，並在執行階段載入系統程式庫。

公開的應用程式 API 維持不變。建置時選擇後端即可，無須修改應用程式程式碼。散佈應用程式前，請先在目標系統上進行測試。

你可以獲得以下好處：

- <strong>Go 編譯不需要 C 工具鏈。</strong>Go 建置不需要 Xcode 命令列工具、`gcc`或 Linux 開發標頭檔。前端建置仍需要其慣用的 JavaScript 工具。
- <strong>不使用 Docker 即可交叉編譯。</strong>只需使用`go build`，即可在任何作業系統上建置 macOS 或 Linux 二進位檔。

Windows 原本就不使用 cgo 建置，因此不受影響。

### 啟用前須知

請使用皆包含此實驗性功能的 CLI 與框架版本。只更新 CLI，並不會將此後端加入應用程式`go.mod`中的舊版框架相依項目。若要試用開發版本，請在應用程式目錄中執行下列命令：

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@master
go get github.com/wailsapp/wails/v3@master
```

這兩個命令都會追蹤持續變動的`master`分支。若要建立可重現的建置，CLI 與框架請使用同一個經測試的提交。請依照[安裝指南](/getting-started/installation/)設定 Go 工具鏈與 CLI。

實驗性的 macOS 與 Linux 後端支援`amd64`和`arm64`。交叉編譯會為指定目標建立可執行檔，但不會在建置機器上執行該檔案。請在目標作業系統及架構上進行測試。

### 使用方式

請選擇下列其中一種方式指定後端。

#### 搭配 Wails CLI 使用

在使用此版本 CLI 產生的專案中，啟用`WAILS_PUREGO`後，`build`、`package`及`dev`會使用實驗性後端。在 POSIX shell 中：

```sh
WAILS_PUREGO=1 wails3 dev
WAILS_PUREGO=1 wails3 build
WAILS_PUREGO=1 wails3 package
```

在 PowerShell 中：

```powershell
$env:WAILS_PUREGO = "1"
wails3 build
```

CLI 可識別`1`、`true`、`on`及`yes`，且不區分大小寫，並會忽略前後空白。空值不會啟用此實驗性功能。

或者，將該標籤傳給單次建置或封裝命令：

```sh
wails3 build -tags purego
wails3 package -tags purego
```

設定此標籤後，更新版的預設 Taskfile 會使用主機的 Go 工具鏈，透過`CGO_ENABLED=0`進行建置。例如，從 macOS 或 Windows 建置 Linux ARM64 可執行檔：

```sh
wails3 build -tags purego GOOS=linux ARCH=arm64
```

舊版專案會保留現有的 Taskfile。請將自訂建置工作與更新後的範本比較，或使用純 Go 測試後端，而不變更這些工作。此環境變數會在 Wails CLI 命令中選取該標籤；直接使用`go build`時，必須明確指定標籤。

#### 搭配純 Go 使用

先建置前端，確保`go:embed`所參照的所有資產都存在，再編譯應用程式：

```sh
CGO_ENABLED=0 go build -tags purego,production .
```

`CGO_ENABLED=1 go build -tags purego`也可使用：後端相同，只是程式的其餘部分仍可使用 cgo。

進行交叉建置時，還需設定`GOOS`與`GOARCH`：

```sh
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -tags purego,production -o myapp .
```

在 Linux 上，使用`-tags purego,gtk3,production`選取 GTK3 技術堆疊。

若要改回預設後端，請移除`purego`標籤並取消設定`WAILS_PUREGO`（在 POSIX shell 中使用`unset WAILS_PUREGO`，在 PowerShell 中則使用`Remove-Item Env:WAILS_PUREGO`）。改回 macOS 或 Linux 的 cgo 後端時，也請移除明確設定的`CGO_ENABLED=0`。

如果使用`CGO_ENABLED=0`建置卻忘記加入標籤，建置會停止並顯示一則指出修正方式的錯誤：

```
undefined: requires_CGO_ENABLED_1_or_build_tag_purego
```

### 執行階段需求

應用程式啟動時會載入系統程式庫，因此目標機器必須已安裝這些程式庫，與 cgo 建置相同：

| 平台 | 程式庫 | 最舊支援版本 |
|---|---|---|
| macOS | AppKit、WebKit（作業系統的一部分） | 與 cgo 後端相同 |
| Linux（預設） | GTK 4、WebKitGTK 6.0 | GTK 4.8、WebKitGTK 2.40（Debian 12） |
| Linux（`gtk3`） | GTK 3、WebKit2GTK 4.1 | GTK 3.24、WebKit2GTK 2.36（Ubuntu 22.04） |

應用程式啟動時會查找高於這些最低版本的函式，且僅在已安裝的程式庫提供這些函式時才使用，因此同一個二進位檔可同時在新舊系統上執行。

應用程式仍會使用原生 UI 程式庫，並且可能動態連結 Linux 的 libc 等系統程式庫。`CGO_ENABLED=0`並不表示執行檔完全採用靜態連結。若缺少必要的執行階段程式庫，請在目標系統上安裝相符的 GTK/WebKit 軟體堆疊；在建置機器上安裝 C 標頭檔，無法解決缺少執行階段程式庫的問題。

### 封裝與簽署

此後端讓編譯不再需要 C 工具鏈。各種套件格式仍有各自的要求。尤其是，選取`purego`並不會讓現有的 AppImage 部署工具支援跨平台。請參閱[Linux 封裝](/guides/build/linux/)。

交叉建置的 macOS 應用程式在發佈前仍需完成適當的簽署與公證。Windows 安裝程式仍需使用其安裝程式工具。請參閱[macOS 封裝](/guides/build/macos/)與[Windows 封裝](/guides/build/windows/)。對於使用 cgo 的相依套件，`CGO_ENABLED=1`也需要 C 編譯器，以及目標平台所需的任何標頭檔或程式庫；選取 Wails 的純 Go 後端並不會免除這些相依套件的要求。

### 驗證

產生的繫結檢查與`puregocheck`會檢查原生符號、選用版本防護機制及後端涵蓋範圍。Linux 同等性工作流程會使用兩種後端測試 GTK4 與 GTK3，包括資產請求、二進位請求本文及原生分派；也會交叉編譯 Linux ARM64 目標。交叉編譯檢查用來確認程式碼能夠建置，而執行階段檢查則會實際測試安裝於測試機器上的原生程式庫。

[後端測試指南](https://github.com/wailsapp/wails/blob/master/v3/TESTING.md)說明同等性探查與端對端輸入測試。回報問題時，請區分編譯失敗與執行階段或封裝失敗。

### 已知差異

- <strong>Linux 上的選單快速鍵採用應用程式層級的註冊方式。</strong>請避免為不同選單指派互相衝突的快速鍵。使用`SetMenu(nil)`明確移除視窗的選單，也會防止已移除選單的快速鍵在該視窗中觸發；其他視窗仍會保留其選單與快速鍵。
- **macOS 的開啟與儲存面板**會在獨立的系統程序中執行，與使用 cgo 時相同；測試會檢查面板能否正常開啟及取消。

還發現其他問題嗎？請告訴我們。

### 意見回饋

請在[實驗](https://github.com/wailsapp/wails/discussions/categories/experiments)下的**無 CGO 建置**討論串中告訴我們使用情況：您的作業系統及版本、建置方式，以及哪些部分正常或未正常運作。提供失敗案例及`wails3 doctor`的輸出，是對我們最有幫助的回報。
