---
title: "Wake"
description: "一個可感知 Wails 的實驗性建置執行器，能以更快的增量建置、結構化輸出及預設平行執行來執行您現有的 Taskfile。"
slug: "experimental/wake"
sourcePath: "experimental/wake.md"
---

@note{type="caution" title="實驗性功能"}
Wake 必須透過 `WAILS_USE_WAKE=true` 主動啟用，<strong>並非</strong>預設執行器。 未設定此變數時，`wails3 build / package / sign / task` 的行為與以往完全相同。 各版本之間的功能涵蓋範圍與行為可能有所變更。

@end

Wake 是 `wails3` 的<strong>實驗性替代建置執行器</strong>。它會讀取專案已有的相同 `Taskfile.yml`，支援相同的 task、dep、var、template、include 與平台命名空間語法，並透過可感知 Wails 的執行器來執行，而不是使用通用的 [Task](https://taskfile.dev) 執行階段。

其目標並非取代 Task，而是提供一個專為 Wails 專案實際建置方式打造的執行器，其語意、輸出與預設值皆與 `wails3` CLI 的其餘部分一致。**如果您始終只使用 Wake，便不必變更 Taskfile。**

## 存在理由

Wake 與 Task 執行階段都已編譯至 `wails3` 中，兩者都不需要另行安裝二進位檔。差別在於 Wake **瞭解其領域**。通用執行器會按照 Taskfile 指定的順序執行其中列出的所有步驟。Wake 知道 Wails 建置實際上<em>是</em>什麼：前端套件會嵌入二進位檔，二進位檔會封裝成平台專用成品，圖示與繫結也會一併產生；它會運用這些知識，以通用執行器無法採用的方式最佳化建置。

- <strong>它只會執行建置實際需要的工作。</strong>Wake 會自行追蹤每個步驟的實際輸入與輸出。對 Go 建置而言，這包括模組圖及其相依步驟的輸出，因此在沒有任何相關變更時，它會完全略過編譯器與連結器，而不會重新執行。只有當 Taskfile 預先明確列出要監控的確切檔案時，通用執行器才能略過某個步驟；Wake 則會根據它對建置已有的瞭解自行判斷。在沒有任何工作可做的重新建置中，耗時約為 **~20 ms（Wake），相較之下 Task 約為 ~316 ms**。冷建置的實際耗時相同，因為主要時間都花在 `npm install`、Vite 與 Go 編譯器上。

- <strong>它知道哪些工作可以同時執行。</strong>由於 Wake 瞭解哪些步驟彼此獨立，因此預設會平行執行這些步驟，結果行也會回報由此獲得的加速幅度。若同層步驟交錯的輸出會干擾問題調查，可使用 `WAKE_SERIAL=true` 停用此行為。

- <strong>由 wails3 控制的結構化輸出。</strong>Wake 透過 wails3 自己的報告器呈現輸出：每個規劃步驟各占一列、即時顯示狀態、最後提供以色彩區分的階段細目，並在失敗面板中提供可點選的 `file:line` 連結。`NO_COLOR` 與非 TTY 環境（CI 記錄）也能妥善降級。

- <strong>內建於其中，因此能隨 Wails 一同發展。</strong>由於 Wake 是`wails3`的一部分，而非第三方工具，因此可以直接疊加新的建置功能，不必等待其他專案實作。這也為未來以原生方式執行跨平台指令碼與工具創造了可能性；目前 Taskfile 在這些情況下會呼叫`wails3`二進位檔，且每次呼叫都會產生一個程序。將這些工作移至同一程序內可降低額外負擔，也表示未來還能進一步加速。

## 啟用 Wake

Wake 完全由 `WAILS_USE_WAKE=true` 環境變數控管。 若未設定此變數（或將其設為 `true` 以外的任何值），每個 `wails3` 命令都會與以往完全相同，使用內嵌的 Task 執行階段。

```bash
# Default: Task runtime, no Wake involvement
wails3 build

# Opt in: Wake drives build / package / sign / task <name>
WAILS_USE_WAKE=true wails3 build
WAILS_USE_WAKE=true wails3 package
WAILS_USE_WAKE=true wails3 task <some-task-name>
```

此旗標涵蓋 `wails3 build`、`wails3 package`、`wails3 sign` 與 `wails3 task <name>`。`wails3 dev` 目前<strong>尚未</strong>受到影響，開發監看程式仍使用自己的管線。

@note{type="tip" title="啟用 Wake 是安全的"}
如果 Wake 遇到尚未實作的 Taskfile 功能，便會在同一程序內將整次執行交給內嵌的 Task 執行階段，不需要安裝外部 `task` 二進位檔。最壞的情況下，您得到的行為也會與未使用此旗標時完全相同。

@end

## 分層本機覆寫

Wake 支援<strong>基礎 Taskfile 加上本機覆寫</strong>。將檔案放在 `Taskfile.yml` 旁，其定義便會優先套用：

| 檔案 | 用途 | 優先順序 |
| --- | --- | --- |
| `Taskfile.yml` | 基礎檔案，已提交 | 最低 |
| `Taskfile.override.yml` / `.yaml` | 已提交、全團隊共用的覆寫 | 中等 |
| `Taskfile.local.yml` / `.yaml` | 個人使用，通常由 Git 忽略 | 最高 |

**合併語意（以本機定義為準）：**

- <strong>同名</strong>的 task 會覆寫基礎 task。當覆寫檔提供清單欄位（`cmds`、`deps`、`sources`、`generates`、`platforms`、`status`、`preconditions`、`aliases`）時，這些欄位會<strong>取代</strong>基礎定義；覆寫檔省略的欄位則會保留基礎定義中的值。
- `env` 與 `vars` 會<strong>逐一依鍵合併</strong>，發生衝突時以覆寫值為準。
- 如果某個 task <strong>只</strong>存在於覆寫檔中，便會予以<strong>新增</strong>。

例如，如果已提交的 `Taskfile.yml` 使用開發旗標進行建置，但您的機器應一律進行正式環境建置：

```yaml
# Taskfile.local.yml (git-ignored, yours)
tasks:
  build:
    cmds:
      - go build -tags production -o bin/app .
  smoke:
    cmds:
      - ./bin/app --selftest
```

現在，`build` 會執行您的正式環境命令，且 `smoke` 可供使用，而不必變更已提交的 Taskfile。

@note{type="note" title="信任模型"}
系統會自動探索並套用覆寫檔，不會顯示提示。這不會授予任何新能力：Taskfile 本來就能執行任意 shell 命令，因此覆寫檔能做的事不會超出直接編輯 `Taskfile.yml` 所能做到的範圍。已提交的 `Taskfile.override.*` 會出現在 PR 差異中；`Taskfile.local.*` 則是在您自己的機器上建立。格式錯誤的覆寫檔會中止執行，而不會遭到無聲略過。若要進行確定性 CI 建置，請設定 `WAILS_NO_OVERRIDES=true`，以完全略過覆寫檔探索。

@end

## 自動回退

如果 Wake 遇到尚未實作的 Taskfile 功能，便會將整次執行交給內嵌的 Task 執行階段。目前下列項目會觸發此回退機制：

- Taskfile 層級的 `dotenv`
- `interleaved` 以外的 `output` 模式
- `requires` 區塊
- `interval`（Taskfile 或 task 層級）
- `always` 以外的 `run` 模式
- task 中的 `short`
- 工作中的`defer`

## 環境變數

| 變數 | 效果 |
| --- | --- |
| `WAILS_USE_WAKE` | 設為`true`時，會為可路由的`wails3`動詞啟用 Wake；設為其他值時則使用 Task 執行階段 |
| `WAILS_NO_OVERRIDES` | `true`會略過`Taskfile.local.*`／`.override.*`探索（確保建置結果可重現） |
| `WAKE_VERBOSE` | 即時串流子行程的 stdout/stderr，而非擷取後僅在失敗時顯示 |
| `WAKE_SILENT` | 完全隱藏工作輸出 |
| `WAKE_SERIAL` | `true`會停用`deps:`的平行分派（預設為平行執行） |
| `WAKE_FORCE` | `true`會略過所有快取，以執行真正徹底的重新建置 |
| `WAKE_DEBUG` | 記錄解析器內部資訊（DAG、相依項目、變數參照、執行路由） |
| `WAKE_NOTICE` | 設為`off`可隱藏每次執行時的「wake (experimental)」通知 |

建置快取位於`.wake/cache.json`（Task 使用`.task/`）。

## 意見回饋

Wake 是一項實驗，而您的意見回饋將決定其未來方向。如果您試用了 Wake，我們很希望知道它是否更快、更清楚，以及是否有任何功能發生問題——最有幫助的報告會說明您執行了什麼、預期結果為何，以及實際發生的情況。請到[Wake 意見回饋討論](https://github.com/wailsapp/wails/discussions/5679)告訴我們。
