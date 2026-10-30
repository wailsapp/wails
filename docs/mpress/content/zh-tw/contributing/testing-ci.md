---
title: "測試與持續整合"
description: "Wails v3如何透過單元測試、整合測試套件、競爭偵測及 GitHub Actions CI 確保品質。"
slug: "contributing/testing-ci"
sourcePath: "contributing/testing-ci.md"
---

穩健的桌面框架需要堅實可靠的測試。Wails v3採用**分層策略**：

| 層級 | 目標 | 工具 |
| --- | --- | --- |
| 單元測試 | 針對獨立函式提供快速回饋 | `go test ./...` |
| 產生器／CLI 測試 | 驗證`wails3 generate bindings`及 CLI 內部串接 | `task test:generator`, `task test:cli` |
| 範本測試 | 確保每個隨附範本仍可建置 | `task test:templates` |
| 競爭偵測 | 找出執行階段與橋接層中的資料競爭 | `go test -race ./...` |
| CI 矩陣 | 確保每個 PR 在不同作業系統上都可靠 | GitHub Actions |

本文件說明<strong>測試位於何處</strong>、<strong>如何執行測試</strong>，以及<strong> Taskfile 會協調哪些工作</strong>。

`v3/TESTING.md`記錄跨平台範例的建置方式；它不是競爭抑制項目清單。

---

## 1. 目錄慣例

```
v3/
├── internal/.../_test.go     # Unit tests for internal packages
├── pkg/.../_test.go          # Public API tests
├── tasks/events/generate.go  # Code generator for event constants (NOT a test harness)
├── tests/                    # Top-level integration test harness
└── TESTING.md                # Cross-platform example-build guidance
```

準則：

- **將單元測試放在程式碼旁**（`foo.go` ↔ `foo_test.go`）。
- 若能改善 API 的整潔性，請對`pkg/`套件（`package application_test`）採用**黑箱風格**。
- 共用測試夾具應放在使用它們的位置（此程式碼樹中沒有中央`internal/testutil/`套件，請改為按套件串接輔助程式）。

---

## 2. 單元測試

### 撰寫測試

```go
func TestEventConstants(t *testing.T) {
    assert.NotEmpty(t, events.Common.WindowFocus)
}
```

建議：

- 使用[`stretchr/testify`](https://github.com/stretchr/testify)，它已經位於`go.mod`中。
- 當多組輸入、邊界情況或預期結果會測試相同行為時，請優先採用**表格驅動**測試。為每個測試案例指定具描述性的名稱。
- 必要時，請在建置標籤（`foo_windows_test.go`、`foo_darwin_test.go`……）後方以虛設實作隔離平台特有行為。

### 涵蓋率要求

新增及變更的邏輯應達到100% 的 Go 陳述式涵蓋率。請測量您所變更的套件，而不要仰賴整個儲存庫的百分比：

```bash
go test ./pkg/application -coverprofile=coverage.out
go tool cover -func=coverage.out
```

有些路徑確實無法在一般測試環境中合理測試，例如僅特定平台會發生的失敗、依賴硬體的行為，或無法安全觸發的防禦性後援路徑。請將這類例外限縮在最小範圍，並在 PR 說明中解釋每一條未涵蓋的路徑。

### 在本機執行

```bash
cd v3
go test ./... -cover
```

也可以透過 Taskfile 執行（請使用實際目標，沒有`task test`捷徑）：

```
task test:cli            # CLI plumbing tests
task test:generator      # bindings generator round-trip tests
task test:templates      # build every shipped template
task test:infrastructure # supporting helpers
task test:examples       # exercise the example matrix (downloads as needed)
task test:all            # everything above
task sanity              # default GTK4 smoke check (legacy: sanity:gtk3)
task precommit           # what you should run before pushing
```

---

## 3. 整合測試

`v3/tests/`包含跨套件整合測試框架。Taskfile 的`test:example:*`與`test:examples:*`目標會在 darwin／windows／linux 上執行編譯檢查（包括 Linux 上以 Docker 為基礎的 GTK3／GTK4 矩陣）。

> 可執行的範例位於`v3/examples/`。測試目標會選取並建置
>
> 適合主機平台或 CI 矩陣的範例。

使用下列命令執行主機平台的冒煙測試套件：

```
task test:examples       # host
task test:examples:all   # full matrix (slow)
```

---

## 4. 競爭偵測

資料競爭對 GUI 執行階段而言是致命問題。

### 調查競爭問題

檢查兩個互相衝突的存取，並在受影響的平台上重現問題。請勿在沒有證據的情況下，將報告視為無害或加以抑制。`v3/TESTING.md`說明的是範例編譯；這些建置成功並不能證明不存在競爭，也不能證明原生執行階段的行為。

### 本機競爭測試套件

```
go test -race ./...
```

> `wails3 dev`沒有`-race`旗標；其 CLI 旗標為`--config`、`--port`，
>
> 以及`-s`（啟用 HTTPS）。若要在競爭偵測器下執行執行階段，
>
> 請使用`go build -race`建置測試應用程式，然後直接執行。

---

## 5. GitHub Actions 工作流程

`.github/workflows/`下實際存在的工作流程檔案（已對照程式碼樹驗證）：

| 檔案 | 用途 |
| --- | --- |
| `build-and-test-v3.yml` | Go 1.26/1.27矩陣、執行階段與產生器檢查、範本／範例編譯，以及應用程式建置檢查。Linux 預設使用 GTK4/WebKitGTK 6.0；舊版變體則會設定`BUILD_TAGS=gtk3`。原生啟動及互動必須另行檢查。 |
| `cross-compile-test-v3.yml` | 交叉編譯健全性檢查 |
| `auto-changelog-v3.yml`, `changelog-v3.yml` | 變更記錄自動化 |
| `nightly-release-v3.yml` | 每夜v3發行成品 |
| `bump-webview2-v3.yml`, `release-webview2.yml` | WebView2 相依套件／版本管理 |
| `build-cross-image.yml` | 建置交叉編譯器容器映像 |
| `publish-npm.yml` | 將內嵌的`@wailsio/runtime` JS 執行階段發佈至 npm |
| `pr-master.yml` | 在 PR 端針對`master`分支執行檢查 |
| `semgrep.yml` | Semgrep 靜態分析 |
| `stale-issues.yml`, `generate-sponsor-image.yml`, `build-and-test.yml`, `weekly-release-v2.yml` | 儲存庫維護及v2工作流程；舊有的 Crowdin 同步／上傳工作流程已停用。 |

CI 步驟對應上述 Taskfile 目標（`task test:cli`、`task test:generator`、`task test:templates`、`task test:examples`……），因此您可以在本機逐一重現 CI。`build-and-test-v3.yml`中的冒煙測試`wails3 build`步驟會以**不附加任何旗標**的方式叫用；`wails3 build`沒有`-skip-package`旗標。

---

## 6. 本機 CI 一致性

沒有單一的`task ci`整合目標。請串接實際目標來重現 CI：

```
task precommit
task test:cli
task test:generator
task test:templates
task test:examples
```

---

## 7. 疑難排解測試失敗

| 症狀 | 可能原因 | 修正方式 |
| --- | --- | --- |
| **`webview_window_darwin.go`中的競爭條件** | 在主執行緒以外變更視窗狀態 | 透過`application.InvokeAsync` / `Invoke`轉送，讓呼叫在主執行緒上執行 |
| **Linux 測試在無頭 CI 中停滯** | GTK 需要顯示器 | 在`xvfb-run`下執行，例如`xvfb-run task test:examples:linux` |
| **範本建置失敗** | 前端鎖定檔已過期 | 對乾淨的目錄重新執行`wails3 init`以更新範本 |
| **Coverpkg 錯誤** | 整合測試匯入了`main` | 改用建置標籤`//go:build integration`，並以條件限制該匯入 |

---

## 8. 新增測試

1. **單元測試** — 建立`*_test.go`，執行`go test ./...`
2. **產生器 / CLI** — 擴充`internal/generator/testcases/`或`internal/commands/*_test.go`下的測試案例，並重新執行`task test:generator` / `task test:cli`
3. **範本 / 範例** — 確認隨附的範本仍可使用`task test:templates`建置

---

## 9. 重要檔案對照表

| 項目 | 路徑 |
| --- | --- |
| 產生器往返測試 | `internal/generator/generate_test.go` |
| 資產建置測試 | `internal/commands/build-assets_test.go` |
| 跨平台範例建置 | `v3/TESTING.md` |
| Taskfile 測試目標 | `v3/Taskfile.yaml` |
| 事件常數產生器 | `v3/tasks/events/generate.go` |
| CI 工作流程 | `.github/workflows/build-and-test-v3.yml`（透過`actions/setup-go@v5`執行 Go 1.26/1.27） |
| 靜態分析 | `.github/workflows/semgrep.yml` |
| 發布執行階段 npm 套件 | `.github/workflows/publish-npm.yml` |

---

在 Wails v3中，品質絕非事後才考慮的問題。透過單元測試、產生器與範本測試套件、競爭條件偵測，以及跨平台 CI 矩陣，你可以放心貢獻，並確信你的變更能在我們支援的每個作業系統上順利通過測試。祝測試順利！
