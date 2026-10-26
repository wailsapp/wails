---
title: "開始使用"
description: "如何開始為 Wails 做出貢獻v3"
slug: "contributing/getting-started"
sourcePath: "contributing/getting-started.md"
---

## 歡迎，貢獻者！

感謝你有興趣為 Wails 做出貢獻！本指南將協助你完成第一次貢獻。

## 必要條件

開始之前，請確認你已具備：

- 已安裝**Go 1.26+**（[下載](https://go.dev/dl/)）
- 已安裝**Node.js 20+**與**npm**（[下載](https://nodejs.org/)）
- 已為你的 GitHub 帳號設定**Git**
- 具備基本的 Go 與 JavaScript/TypeScript 知識

### 各平台的特定需求

**macOS：**

- Xcode 命令列工具：`xcode-select --install`

**Windows：**

- 建議使用 MSYS2 或類似的 Unix-like 環境
- WebView2 執行階段（通常已預先安裝於 Windows 11）

**Linux：**

- `gcc`、`pkg-config`、`libgtk-4-dev`、`libwebkitgtk-6.0-dev`（預設 GTK4 技術堆疊）
- 透過以下方式安裝：`sudo apt install build-essential pkg-config libgtk-4-dev libwebkitgtk-6.0-dev`（Debian/Ubuntu）
- 若使用舊版`-tags gtk3`建置路徑，還須安裝`libgtk-3-dev`與`libwebkit2gtk-4.1-dev`

## 貢獻流程概覽

一般的貢獻工作流程包含以下步驟：

1. **建立分支複本並複製**——建立你自己的 Wails 儲存庫副本
2. **設定**——建置 Wails CLI 並驗證你的環境
3. **建立分支**——為你的變更建立功能分支
4. **開發**——依照我們的程式碼撰寫標準進行變更
5. **測試**——執行測試以確認一切正常運作
6. **提交**——使用清楚且符合 Conventional Commits 規範的提交訊息
7. **送出**——開啟 Pull Request 以供審查
8. **反覆修正**——回應意見並進行調整
9. **合併**——核准後，你的變更就會成為 Wails 的一部分！

## 逐步指南

選擇你的貢獻類型：

@tabs
[錯誤修正]
@steps
### 尋找或回報錯誤
- 確認[GitHub Issues](https://github.com/wailsapp/wails/issues)中是否已回報此錯誤
- 若尚未回報，請建立新 Issue 並提供重現步驟
- 開始處理前，請等待確認

### 建立分支複本並複製
在[github.com/wailsapp/wails/fork](https://github.com/wailsapp/wails/fork)建立儲存庫的分支複本

複製你的分支複本：

```bash
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails
git remote add upstream https://github.com/wailsapp/wails.git
```

### 建置並驗證
建置 Wails，並確認你能重現此錯誤：

```bash
cd v3
go build -o ../wails3 ./cmd/wails3

# Reproduce the bug to understand it
```

### 建立錯誤修正分支
為你的修正建立分支：

```bash
git checkout -b fix/issue-123-window-crash
```

### 修正錯誤
- 僅進行修正此錯誤所需的最少變更
- 不要重構無關的程式碼
- 新增或更新測試，以防止問題再次發生

```bash
# Make your changes
# Add tests in *_test.go files
```

### 測試你的修正
執行測試，確認修正有效：

```bash
go test ./...

# Test the specific package
go test ./pkg/application -v

# Run with race detector
go test ./... -race
```

### 提交你的修正
使用清楚的訊息提交：

```bash
git commit -m "fix: prevent window crash when closing during initialization

Fixes #123"
```

### 送出 Pull Request
推送並建立 PR：

```bash
git push origin fix/issue-123-window-crash
```

在 PR 說明中：

- 說明錯誤及其根本原因
- 說明你的修正
- 引用相關 Issue：「Fixes #123」
- 說明修正前後的行為

### 回應意見
處理審查意見，並視需要更新你的 PR。

@end

[WEP（功能增強）]
@steps
### 撰寫 WEP
- 閱讀[WEP（Wails Enhancement Proposal）流程](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/README.md)
- 將 WEP 範本複製到`v3/wep/proposals/<name>/proposal.md`
- 開啟標題為`[WEP] <title>`且僅包含 WEP 的草稿 PR
- 實作前，請等待維護者做出決定

### 建立分支複本並複製
在[github.com/wailsapp/wails/fork](https://github.com/wailsapp/wails/fork)建立儲存庫的分支複本

複製你的分支複本：

```bash
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails
git remote add upstream https://github.com/wailsapp/wails.git
```

### 設定開發環境
建置 Wails 並驗證你的環境：

```bash
cd v3
go build -o ../wails3 ./cmd/wails3

# Run tests to ensure everything works
go test ./...
```

### 建立功能分支
建立名稱清楚表達用途的分支：

```bash
git checkout -b feat/window-transparency-support
```

### 實作功能
- 遵循我們的[程式碼撰寫標準](/contributing/standards/)
- 讓變更聚焦於此功能
- 撰寫簡潔且有完整文件的程式碼
- 新增全面的測試

```bash
# Example: Adding a new window method
# 1. Add to window.go interface
# 2. Implement in platform files (darwin, windows, linux)
# 3. Add tests
# 4. Update documentation
```

### 完整測試
測試你的功能：

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

### 為你的功能撰寫文件
- 為所有公開 API 新增文件註解
- 更新`/docs/mpress/content/`中的相關文件
- 如適用，請新增範例

### 依循規範提交
使用 Conventional Commits：

```bash
git commit -m "feat: add window transparency support

- Add SetTransparent() method to Window API
- Implement for macOS, Windows, and Linux
- Add tests and documentation

Closes #456"
```

### 送出 Pull Request
推送並建立 PR：

```bash
git push origin feat/window-transparency-support
```

在 PR 中：

- 說明此功能及其使用情境
- 提供範例或螢幕擷取畫面
- 列出所有破壞性變更
- 引用已接受的 WEP PR

### 根據審查意見反覆修正
維護者可能會要求變更。請耐心等候並保持合作。

@end

[文件]
歡迎直接提交修正 PR，無須先建立議題。僅修正文件時，不需要提供失敗的程式碼測試。請參閱[使用 M-Press 撰寫文件](/contributing/documentation/)，了解 M-Press 的安裝、原始檔路徑、預覽、驗證及 PR 步驟。

@end

## 尋找可處理的議題

- 尋找帶有[`good first issue`](https://github.com/wailsapp/wails/labels/good%20first%20issue)標籤的議題
- 查看[`help wanted`](https://github.com/wailsapp/wails/labels/help%20wanted)議題
- 瀏覽[未解決的議題](https://github.com/wailsapp/wails/issues)，並請求指派給您

## 取得協助

- <strong>Discord：</strong>加入[Wails Discord](https://discord.gg/JDdSxwjhGf)
- <strong>討論區：</strong>在[GitHub Discussions](https://github.com/wailsapp/wails/discussions)中發文
- <strong>議題：</strong>若是可重現的錯誤，請建立議題；如有問題，請使用 Discussions；若要提出功能改進，請提交 WEP PR

## 行為準則

請保持尊重、提供具建設性的意見，並友善接納他人。我們正共同建立一個友善的社群，專注於攜手打造優秀的軟體。

## 後續步驟

- 設定您的[開發環境](/contributing/setup/)
- 檢閱我們的[程式碼撰寫規範](/contributing/standards/)
- 探索[技術文件](/contributing/overview/)
