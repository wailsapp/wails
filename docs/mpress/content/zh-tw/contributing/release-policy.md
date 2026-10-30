---
title: "發布政策"
description: "Wails v3 如何為夜間版與每週版編訂版本、維護者如何發布這些版本，以及使用者為何應固定版本。"
slug: "contributing/release-policy"
sourcePath: "contributing/release-policy.md"
---

## 以一張表說明整套模型

Wails v3依**固定時程，而非依功能**發布版本。在`v3.0.0`基準版本之後：

| 版本類型 | 標籤 | 發布頻率 | GitHub 發布版本 | npm `@wailsio/runtime` |
| --- | --- | --- | --- | --- |
| 夜間版 | `v3.<next>.0-nightly.<YYYYMMDD>`（UTC 日期；若同一天重新執行，則附加`.N`） | 每日自動執行，僅在`UNRELEASED_CHANGELOG.md`有內容時發布 | 預發布版本，**不會**標為 Latest | 相同版本但不含`v`，dist-tag 為`nightly` |
| 每週版 | `v3.<next>.0` | 於 Australia/Sydney 時區星期日20:00自動執行<strong>（日期尚未確認）</strong> | 一般發布版本，標為 Latest | 相同版本，dist-tag 為`latest` |
| 緊急修補版 | `v3.<current>.<N>`, `N > 0` | 很少發布，手動執行 | 一般版本、Latest | 相同版本，dist-tag 為`latest` |
| 主要版本 | `v4.0.0` | 僅用於破壞性變更 | 一般發布版本 | `latest` |

`<next>`是最高的**已發布**v3次要版本加一。nightly 是**下一個每週版本**的預覽，因此其基準是最高的發布標籤（忽略預發布版本）：在修補版本`v3.1.1`之後，下一個 nightly 仍是`v3.2.0-nightly.*`。它必須使用較高的版本號，因為預發布版本的排序會低於其對應的正式版本：`v3.0.0-nightly.<date>`會排在`v3.0.0`之前。

**提高次要版本號並不代表功能里程碑。** `v3.4.0`代表「`v3.3.0`之後的每週發布版本」。請閱讀版本資訊而非僅查看版本號，以瞭解變更內容。

政策依版本的**格式**決定：帶有`-nightly.`尾碼的是夜間版；帶有任何其他尾碼（`-alpha`、`-beta`、`-rc`）的是正式發布前的舊版建置；沒有尾碼的則是發布版本。

### 同步版本

CLI、Go 模組與 npm 套件`@wailsio/runtime`採用**相同版本**，取自`v3/internal/version/version.txt`（標籤`v3.2.0-nightly.20261102`在 npm 上為`3.2.0-nightly.20261102`）。以`wails3 init`建立的專案會將`@wailsio/runtime`固定為 CLI 本身的版本；CLI 的開發建置則會退回使用`latest`。

## 完整範例

| 日期 | 事件 | 標籤 | `go install ...@latest`解析為 |
| --- | --- | --- | --- |
| 10月25日（星期日） | 正式版，手動發布 | `v3.0.0` | `v3.0.0` |
| 10月26日（星期一）至10月31日（星期六） | 夜間版 | `v3.1.0-nightly.20261026` ... `v3.1.0-nightly.20261031` | `v3.0.0` |
| 11月1日（星期日）20:00（雪梨時間） | 每週版 | `v3.1.0` | `v3.1.0` |
| 11月2日（星期一） | 夜間版（預覽下一個每週版） | `v3.2.0-nightly.20261102` | `v3.1.0` |
| 11月4日（星期三） | 從發布分支手動發布緊急修補版 | `v3.1.1` | `v3.1.1` |
| 11月8日（星期日） | 每週版 | `v3.2.0` | `v3.2.0` |

每週版的排序一律高於同版本號的所有夜間版，而夜間版絕不會變成`@latest`。

## 使用者執行方式

| 目標 | 命令 |
| --- | --- |
| 最新發布版本（每週版或緊急修補版） | `go install github.com/wailsapp/wails/v3/cmd/wails3@latest`或`wails3 update cli` |
| **可重現（建議）** | 固定確切標籤：在`go.mod`、CI 與文件中使用`v3.1.0` |
| 夜間版 | 使用其確切標籤（`@v3.2.0-nightly.20261102`）、`wails3 update cli -nightly` |
| 任何類型的最新建置 | `wails3 update cli -pre` |
| Runtime npm 套件 | `@wailsio/runtime`（`latest` = 每週版／緊急修補版，`nightly` = 夜間版） |

其原理如下：Go 會將`@latest`解析為**不含預發布尾碼**的最高標籤，而且只有在不存在發布版本時才會考慮預發布版本。Go Proxy 無法辨識 GitHub 的「pre-release」核取方塊，因此夜間版使用 semver 預發布尾碼，而不依賴該核取方塊。在正式環境中一律固定版本。

## 分支

| 分支 | 用途 |
| --- | --- |
| `master` | 開發與夜間版。 |
| `releases/**` | 用於一次性發布的短期分支：`v3.0.0`正式版與緊急修補版。必須執行與`master`相同的 CI。 |

每週版從`master`發布，但**不是從其分支頂端**發布：工作流程會選取自上一個每週版以來，所有必要檢查（`v3 Build Results`、`Run Go Tests`、`Bundled runtime matches source`）皆通過的最新提交，並為該提交加上標籤。若沒有任何提交通過檢查，就不會發布。夜間版從`master`發布；HEAD 未通過檢查時只會顯示警告。

## 如何發布版本

發布權限僅授予專案負責人（請參閱治理文件）。貢獻者絕不建立標籤或發布版本；他們只將項目新增至`v3/UNRELEASED_CHANGELOG.md`。

所有作業都從**Nightly + Weekly Release v3**（`.github/workflows/nightly-release-v3.yml`）執行：

- <strong>夜間版（自動）</strong>於 UTC 15:00執行。未發布的變更記錄為空時會略過。
- <strong>每週（自動）</strong>於雪梨時間20:00執行。GitHub cron 僅支援 UTC，因此會排定兩個 cron（UTC 09:00和10:00），而防護工作只會在雪梨時間為20:00時繼續執行（會考量日光節約時間）。每週版本會通過發布閘門、選取最新且檢查通過的提交、遞增次要版本，並為該提交加上標籤而不變更變更記錄；接著只更新`master`上的版本欄位（`version.txt`和執行階段套件版本），因此會保留在所選提交之後合併的相依套件與鎖定檔更新。npm 從標籤發布，而不是從`master`發布。
- <strong>手動。</strong>執行工作流程時使用`channel`、`dry_run`（預設開啟），並可選擇使用`release_version`。試執行會輸出版本、類型、GitHub 旗標、npm dist-tag 與版本資訊，而不會變更 GitHub 或 npm。
- <strong>GA 或修補版本。</strong>設定`release_version`（`v3.0.0`、`v3.1.1`）。較舊版本線的修補版本也可以設定`ref: releases/...`；其發布提交會以 cherry-pick 方式帶回`master`。

本機模擬執行：`cd v3/tasks/release && go run release.go --dry-run --channel weekly`；`go run release.go --print-policy v3.2.0-nightly.20261102`會顯示推導出的旗標。

### 安全防護機制

- 自動產生的版本必須嚴格大於所有現有標籤（依語意化版本順序，預發行版小於正式版），否則工具會拒絕執行。
- 同一時間只能執行一個發行作業（共用並行群組）。
- npm 發布由發行**標籤**觸發，並建置該標籤所指向的提交。

## 發行閘門

只要存在以下任何項目，**每週發行**就會拒絕建立標籤（每夜發行絕不會如此）：

1. 具有`P0`或`release-blocker`標籤的未結問題或提取要求；
2. 自上次每週發行以來建立、具有`regression`標籤的未結問題；
3. 嚴重性為重大或高的未結 Dependabot 警示、嚴重性為重大或高的未結程式碼掃描警示，以及任何未結的祕密掃描警示，且**僅限已發布的程式碼**。若警示只位於`v3/examples/`、`v2/examples/`、`website/`、`docs/`、`.github/`或任何`examples/`目錄下，則會忽略該警示（`v3/tasks/releasegate/main.go`中的`SHIPPED_ALERT_IGNORE`，**由擁有者審查**）；
4. 仍處於`triage`狀態的儲存庫安全公告。

**誰可以設定阻擋項目。**只有在某標籤最近一次的`labeled`事件（之後沒有`unlabeled`事件）是由`wailsapp/maintainers`團隊成員，或具有`maintain`／`admin`權限的協作者執行時，該標籤才會生效。其他任何人套用的相同標籤都會被忽略並記錄，因此外部人員無法阻擋或解除阻擋版本發布。自動套用的`high-priority`標籤刻意**不**採用。`taliesin-ai`等機器人帳號若具備`maintain`，即可通過允許清單；這是預期行為，分流代理程式可在提供佐證的情況下套用阻擋標籤。

**失敗時的行為。**閘門採取**預設拒絕**策略：若無法讀取警示 API 或允許清單，每週發行就會停止。失敗時不會建立任何標籤，會將清單（含連結）寫入作業摘要；若已設定儲存庫變數`WEEKLY_STATUS_ISSUE`，還會更新唯一一個「每週發行狀態」問題中的留言。

**擁有者必須在 GitHub 上完成的設定**（系統不會自動建立任何項目）：建立`P0`、`release-blocker`和`regression`標籤；建立「每週發行狀態」問題，並將其編號設為`WEEKLY_STATUS_ISSUE`；授予`WAILS_REPO_TOKEN`對 Issues、Pull requests、Administration（協作者權限）、Dependabot alerts、Code scanning alerts、Secret scanning alerts、Repository security advisories 及組織 Members 的讀取權限，並授予 Issues 寫入權限以張貼狀態留言。預設的`GITHUB_TOKEN`無法讀取這些項目。

**略過。**`override_reason`（工作流程手動派送輸入）會略過閘門。略過原因、執行者及每個被略過的阻擋項目都會寫入摘要和狀態問題。只有發行負責人應使用此功能。GitHub 無法將某項輸入限制為僅供一人使用，因此請在設有必要審核者的環境中執行發行作業（建議），讓略過閘門必須獲得明確核准。

## 相關內容

- [測試與持續整合](/contributing/testing-ci/)
- [維護者分流處理](/contributing/maintainer-triage/)
