---
title: "開始使用"
description: "如何開始為 Wails v3 做出貢獻"
slug: "contributing/getting-started"
sourcePath: "contributing/getting-started.md"
---

## 歡迎貢獻！

感謝您有興趣為 Wails 做出貢獻！本指南將協助您完成第一次貢獻。

## 必要條件

開始之前，請確認您已具備：

- 已安裝 **Go 1.25+**（[下載](https://go.dev/dl/)）
- **Node.js 20+** 和 **npm**（[下載](https://nodejs.org/)）
- 已使用您的 GitHub 帳戶設定 **Git**
- 具備 Go 和 JavaScript/TypeScript 的基本知識

### 各平台的特定要求

**macOS：**

- Xcode 命令列工具：`xcode-select --install`

**Windows：**

- 建議使用 MSYS2 或類似 Unix 的環境
- WebView2 執行階段（通常已預先安裝於 Windows 11）

**Linux：**

- `gcc`、`pkg-config`、`libgtk-4-dev`、`libwebkitgtk-6.0-dev`（預設 GTK4 技術堆疊）
- 安裝方式：`sudo apt install build-essential pkg-config libgtk-4-dev libwebkitgtk-6.0-dev`（Debian/Ubuntu）
- 若使用舊版 `-tags gtk3` 建置路徑，還須安裝 `libgtk-3-dev` 和 `libwebkit2gtk-4.1-dev`

## 貢獻流程概覽

一般的貢獻工作流程包含以下步驟：

1. **建立分支副本並複製**——建立您自己的 Wails 儲存庫副本
2. **設定**——建置 Wails CLI 並驗證您的環境
3. **建立分支**——為您的變更建立功能分支
4. **開發**——依照我們的程式碼規範進行變更
5. **測試**——執行測試，確保一切正常運作
6. **提交**——使用清楚且符合慣例的提交訊息提交變更
7. **送出**——建立拉取請求以供審查
8. **反覆修正**——回應意見並進行調整
9. **合併**——核准後，您的變更就會成為 Wails 的一部分！

## 逐步指南

選擇您的貢獻類型：

@tabs
[錯誤修正]
@steps
### 尋找或回報錯誤
- 檢查 [GitHub Issues](https://github.com/wailsapp/wails/issues) 中是否已有人回報此錯誤
- 若尚未回報，請建立新的議題並附上重現步驟
- 請等待確認後再開始進行

### 建立分支副本並複製
在 [github.com/wailsapp/wails/fork](https://github.com/wailsapp/wails/fork) 建立儲存庫的分支副本

複製您的分支副本：

```bash
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails
git remote add upstream https://github.com/wailsapp/wails.git
```

### 建置並驗證
建置 Wails，並確認您可以重現此錯誤：

```bash
cd v3
go build -o ../wails3 ./cmd/wails3

# Reproduce the bug to understand it
```

### 建立錯誤修正分支
為您的修正建立分支：

```bash
git checkout -b fix/issue-123-window-crash
```

### 修正錯誤
- 僅進行修正此錯誤所需的最少變更
- 請勿重構不相關的程式碼
- 新增或更新測試，以防止迴歸

```bash
# Make your changes
# Add tests in *_test.go files
```

### 測試您的修正
執行測試以確認修正有效：

```bash
go test ./...

# Test the specific package
go test ./pkg/application -v

# Run with race detector
go test ./... -race
```

### 提交您的修正
使用清楚的訊息提交：

```bash
git commit -m "fix: prevent window crash when closing during initialization

Fixes #123"
```

### 送出拉取請求
推送並建立 PR：

```bash
git push origin fix/issue-123-window-crash
```

請在 PR 說明中：

- 說明錯誤及其根本原因
- 說明您的修正
- 引用相關議題：「Fixes #123」
- 附上修正前後的行為

### 回應意見
處理審查意見，並視需要更新您的 PR。

@end

[WEP（增強功能）]
@steps
### 撰寫 WEP
- 閱讀 [WEP（Wails Enhancement Proposal）流程](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md)
- 將 WEP 範本複製至 `v3/wep/proposals/<name>/proposal.md`
- 建立標題為 `[WEP] <title>` 且僅包含 WEP 的草稿 PR
- 請等待維護者做出決定後再開始實作

### 建立分支副本並複製
在 [github.com/wailsapp/wails/fork](https://github.com/wailsapp/wails/fork) 建立儲存庫的分支副本

複製您的分支副本：

```bash
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails
git remote add upstream https://github.com/wailsapp/wails.git
```

### 設定開發環境
建置 Wails 並驗證您的環境：

```bash
cd v3
go build -o ../wails3 ./cmd/wails3

# Run tests to ensure everything works
go test ./...
```

### 建立功能分支
建立名稱清楚描述用途的分支：

```bash
git checkout -b feat/window-transparency-support
```

### 實作功能
- 遵循我們的[程式碼撰寫標準](/contributing/standards/)
- 讓變更聚焦於此功能
- 撰寫清晰且有文件說明的程式碼
- 新增完整的測試

```bash
# Example: Adding a new window method
# 1. Add to window.go interface
# 2. Implement in platform files (darwin, windows, linux)
# 3. Add tests
# 4. Update documentation
```

### 徹底測試
測試您的功能：

```bash
# Unit tests
go test ./pkg/application -v

# Integration test - create a test app
cd ..
./wails3 init -n feature-test
cd feature-test
# Add code using your new feature
../wails3 dev
```

### 為功能撰寫文件
- 為所有公開 API 新增文件字串
- 更新`/docs/mpress/content/`中的相關文件
- 如適用，請新增範例

### 依照慣例提交
使用約定式提交：

```bash
git commit -m "feat: add window transparency support

- Add SetTransparent() method to Window API
- Implement for macOS, Windows, and Linux
- Add tests and documentation

Closes #456"
```

### 提交拉取請求
推送並建立 PR：

```bash
git push origin feat/window-transparency-support
```

在您的 PR 中：

- 說明此功能及其使用情境
- 提供範例或螢幕擷取畫面
- 列出所有破壞性變更
- 引用已接受的 WEP PR

### 根據審查意見反覆改進
維護人員可能會要求變更。請保持耐心並積極協作。

@end

[文件]
歡迎直接提交修正 PR，無須先建立議題。僅修正文件時，不需要提供會失敗的程式碼測試。請參閱[修正文件](/contributing/documentation/)，了解來源檔案路徑，並依照指南進行 M-Press 安裝、預覽、驗證及 PR 提交。

@end

## 尋找可著手處理的議題

- 尋找[`good first issue`](https://github.com/wailsapp/wails/labels/good%20first%20issue)標籤
- 查看[`help wanted`](https://github.com/wailsapp/wails/labels/help%20wanted)議題
- 瀏覽[未結議題](https://github.com/wailsapp/wails/issues)，並要求將議題指派給您

## 取得協助

- <strong>Discord：</strong>加入[Wails Discord](https://discord.gg/JDdSxwjhGf)
- <strong>討論區：</strong>在[GitHub Discussions](https://github.com/wailsapp/wails/discussions)中發文
- <strong>議題：</strong>針對可重現的錯誤建立議題；如有疑問，請使用 Discussions；如要提出功能增強，請提交 WEP PR

## 行為準則

請以尊重、建設性且友善的態度待人。我們正共同打造一個友善的社群，專注於攜手建立優秀的軟體。

## 後續步驟

- 設定您的[開發環境](/contributing/setup/)
- 檢閱我們的[程式碼撰寫標準](/contributing/standards/)
- 探索[技術文件](/contributing/overview/)
