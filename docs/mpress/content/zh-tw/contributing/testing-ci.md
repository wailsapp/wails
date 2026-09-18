---
title: "測試與持續整合"
description: "Wails v3 如何透過單元測試、整合測試套件、資料競爭偵測與 GitHub Actions CI 確保品質。"
slug: "contributing/testing-ci"
sourcePath: "contributing/testing-ci.md"
---

穩健的桌面框架需要堅若磐石的測試。 Wails v3 採用<strong>分層策略</strong>：

| 層級 | 目標 | 工具 |
| --- | --- | --- |
| 單元測試 | 針對獨立函式提供快速回饋 | `go test ./...` |
| 產生器／CLI 測試 | 驗證`wails3 generate bindings`與 CLI 的串接機制 | `task test:generator`、`task test:cli` |
| 範本測試 | 確保每個隨附範本仍可建置 | `task test:templates` |
| 資料競爭偵測 | 找出執行階段與橋接層中的資料競爭 | `go test -race ./...` |
| CI 矩陣 | 確保每個 PR 在不同作業系統上皆可運作 | GitHub Actions |

本文件說明<strong>測試位於何處</strong>、**如何執行測試**，以及<strong>Taskfile 會協調哪些工作</strong>。

> 舊版草稿以`pkg/application/RACE.md`引用的資料競爭指南
>
> 目前位於`v3/TESTING.md`。

---

## 1. 目錄慣例

```
v3/
├── internal/.../_test.go     # Unit tests for internal packages
├── pkg/.../_test.go          # Public API tests
├── tasks/events/generate.go  # Code generator for event constants (NOT a test harness)
├── tests/                    # Top-level integration test harness
└── TESTING.md                # Race / Cgo testing guidance
```

準則：

- **將單元測試放在程式碼旁**（`foo.go` ↔ `foo_test.go`）。
- 若能改善 API 的整潔性，請對`pkg/`套件（`package application_test`）使用<strong>黑箱樣式</strong>。
- 共用測試資料應放在使用它們的位置（此工作樹中沒有集中管理的`internal/testutil/`套件，請改為在各套件中個別串接輔助程式）。

---

## 2. 單元測試

### 撰寫測試

```go
func TestEventConstants(t *testing.T) {
    assert.NotEmpty(t, events.Common.WindowFocus)
}
```

建議：

- 使用[`stretchr/testify`](https://github.com/stretchr/testify)——已包含在`go.mod`中。
- 當多個輸入、邊界案例或預期結果會驗證相同行為時，優先使用<strong>表格驅動</strong>測試。為每個案例指定具描述性的名稱。
- 必要時，請在建置標籤（`foo_windows_test.go`、`foo_darwin_test.go`等）後方以虛設實作隔離平台特有行為。

### 涵蓋率要求

新增及變更的邏輯預期應達到100% 的 Go 陳述式涵蓋率。請測量您所變更的套件，而不要依賴整個儲存庫的百分比：

```bash
go test ./pkg/application -coverprofile=coverage.out
go tool cover -func=coverage.out
```

有些路徑無法在一般測試環境中合理測試，例如僅限特定平台的失敗、取決於硬體的行為，或無法安全觸發的防禦性備援。請將這類例外限縮至最小範圍，並在 PR 說明中解釋每一條未涵蓋的路徑。

### 在本機執行

```bash
cd v3
go test ./... -cover
```

也可透過 Taskfile 執行（請使用實際目標；不存在`task test`捷徑）：

```
task test:cli            # CLI plumbing tests
task test:generator      # bindings generator round-trip tests
task test:templates      # build every shipped template
task test:infrastructure # supporting helpers
task test:examples       # exercise the example matrix (downloads as needed)
task test:all            # everything above
task sanity              # quick smoke check (also: sanity:gtk4)
task precommit           # what you should run before pushing
```

---

## 3. 整合測試

`v3/tests/`包含跨套件整合測試框架。Taskfile 的`test:example:*`與`test:examples:*`目標會在 darwin / windows / linux 上執行建置及啟動檢查（包括 Linux 上基於 Docker 的 GTK3 / GTK4 矩陣）。

> 可執行的範例位於`v3/examples/`。測試目標會選取並建置
>
> 適用於主機平台或 CI 矩陣的範例。

使用以下命令執行主機平台的冒煙測試套件：

```
task test:examples       # host
task test:examples:all   # full matrix (slow)
```

---

## 4. 資料競爭偵測

資料競爭對 GUI 執行階段而言是致命問題。

### 資料競爭指南

請參閱`v3/TESTING.md`以瞭解：

- 已知的良性資料競爭及抑制這些競爭的理由
- 如何解讀跨越 Cgo 邊界的堆疊追蹤（Linux GTK + WebKit2GTK）

### 本機資料競爭測試套件

```
go test -race ./...
```

> `wails3 dev`沒有`-race`旗標；其 CLI 旗標為`--config`、`--port`，
>
> 以及`-s`（啟用 HTTPS）。若要在資料競爭偵測器下測試執行階段，
>
> 請使用`go build -race`建置測試應用程式，然後直接執行。

---

## 5. GitHub Actions 工作流程

`.github/workflows/`下的實際工作流程檔案（已對照工作樹驗證）：

| 檔案 | 用途 |
| --- | --- |
| `build-and-test-v3.yml` | 主要的 v3 建置與測試矩陣。搭配`go-version: 1.25`使用`actions/setup-go@v5`。依序執行`task runtime:check`、`task runtime:test`、`task runtime:build`、`task test:examples`（GTK4 路徑還會執行`BUILD_TAGS=gtk4 task test:examples`）、`task generator:test:check`、`task install`，最後以`wails3 build`進行冒煙檢查。Linux 工作會安裝`libgtk-3-dev libwebkit2gtk-4.1-dev libwayland-dev build-essential pkg-config xvfb x11-xserver-utils at-spi2-core xdg-desktop-portal-gtk`，並在`dbus-run-session -- xvfb-run`下執行測試套件。 |
| `cross-compile-test-v3.yml` | 交叉編譯健全性檢查 |
| `auto-changelog-v3.yml`、`changelog-v3.yml` | 變更記錄自動化 |
| `nightly-release-v3.yml` | 每夜建置的 v3 發行成品 |
| `bump-webview2-v3.yml`、`release-webview2.yml` | WebView2 相依性／版本發布管理 |
| `build-cross-image.yml` | 建置交叉編譯器容器映像 |
| `publish-npm.yml` | 將內嵌的`@wailsio/runtime` JS 執行階段發布至 npm |
| `pr-master.yml` | 針對`master`分支執行 PR 端檢查 |
| `semgrep.yml` | Semgrep 靜態分析 |
| `stale-issues.yml`、`issue-labeler.yml`、`file-labeler.yml`、`claude.yml`、`generate-sponsor-image.yml`、`sync-translated-documents.yml`、`upload-source-documents.yml`、`build-and-test.yml`、`weekly-release-v2.yml` | 儲存庫維護／v2 端流程 |

此工作樹中<strong>沒有</strong>`qodana.yaml`，也<strong>沒有</strong>`runtime.yml`——此頁面的舊版草稿曾提及兩者，但只有`semgrep.yml`涵蓋靜態分析，而執行階段 JS 套件則透過`publish-npm.yml`發布。

CI 步驟與上述 Taskfile 目標（`task test:cli`、`task test:generator`、`task test:templates`、`task test:examples`……）一致，因此你可以在本機逐一重現 CI。`build-and-test-v3.yml`中的冒煙`wails3 build`步驟在叫用時<strong>不會加上任何額外旗標</strong>——`wails3 build`沒有`-skip-package`旗標。

---

## 6. 本機 CI 對等性

沒有單一的`task ci`統括目標。請串接實際目標來重現 CI：

```
task precommit
task test:cli
task test:generator
task test:templates
task test:examples
```

---

## 7. 疑難排解失敗的測試

| 症狀 | 可能原因 | 修正方式 |
| --- | --- | --- |
| <strong>`webview_window_darwin.go`</strong>中的競爭情況 | 在主執行緒之外變更視窗狀態 | 透過`application.InvokeAsync`／`Invoke`進行封送，讓呼叫在主執行緒上執行 |
| **Linux 測試在無頭 CI 中停滯** | GTK 需要圖形顯示環境 | 在`xvfb-run`下執行，例如`xvfb-run task test:examples:linux` |
| **範本建置失敗** | 前端鎖定檔已過期 | 對乾淨目錄重新執行`wails3 init`，以重新整理範本 |
| **Coverpkg 錯誤** | 整合測試匯入了`main` | 改用建置標籤`//go:build integration`，並以條件限制該匯入 |

---

## 8. 新增測試

1. **單元測試**——建立`*_test.go`，然後執行`go test ./...`
2. **產生器／CLI**——擴充`internal/generator/testcases/`或`internal/commands/*_test.go`下的測試案例，然後重新執行`task test:generator`／`task test:cli`
3. **範本／範例**——確認隨附範本仍可使用`task test:templates`建置

---

## 9. 重要檔案對照表

| 項目 | 路徑 |
| --- | --- |
| 產生器往返測試 | `internal/generator/generate_test.go` |
| 建置資產測試 | `internal/commands/build-assets_test.go` |
| 競爭情況／Cgo 指南 | `v3/TESTING.md` |
| Taskfile 測試目標 | `v3/Taskfile.yaml` |
| 事件常數產生器 | `v3/tasks/events/generate.go` |
| CI 工作流程 | `.github/workflows/build-and-test-v3.yml`（透過`actions/setup-go@v5`執行 Go 1.25） |
| 靜態分析 | `.github/workflows/semgrep.yml` |
| 發布執行階段 npm 套件 | `.github/workflows/publish-npm.yml` |

---

在 Wails v3 中，品質絕非事後才考量的事項。透過單元測試、產生器／範本測試套件、競爭情況偵測，以及跨平台 CI 矩陣，你可以放心貢獻，因為你的變更能在我們支援的每個作業系統上順利通過測試。祝測試順利！
