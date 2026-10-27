---
title: "貢獻指南"
description: "為 Wails 做出貢獻"
slug: "contributing"
sourcePath: "contributing/index.md"
---

## 歡迎貢獻者！

歡迎為 Wails 做出貢獻！無論是修正錯誤、新增功能，還是改善文件，我們都非常感謝您的協助。

## 貢獻方式

### 1. 回報問題

發現錯誤了嗎？請[建立議題](https://github.com/wailsapp/wails/issues/new)，並附上：

- 清楚的說明
- 重現步驟
- 預期行為與實際行為
- 系統資訊
- 程式碼範例

### 2. 改善文件

歡迎直接提交修正 PR，無須先建立議題或提供失敗的程式碼測試。  
如需下載、在 macOS/Linux/Windows 上安裝、即時預覽及提交變更的相關說明，請參閱[使用 M-Press 撰寫文件](/contributing/documentation/)。

我們隨時歡迎以下文件改善：

- 修正錯字與錯誤
- 新增範例
- 釐清說明
- 翻譯內容

### 3. 提交程式碼

透過 Pull Request 貢獻程式碼：

- 錯誤修正
- 新功能
- 效能改善
- 測試

### 4. 提出改進提案（WEP）

新功能及公開行為的變更須遵循 Wails Enhancement Proposal（WEP）流程。此流程可讓功能開發保持透明，並確保每項獲接受的提案都有實作者。請勿建立功能請求議題。

1. 您可以選擇先在 GitHub Discussions 的[構想](https://github.com/wailsapp/wails/discussions/categories/ideas)分類或[Discord](https://discord.gg/JDdSxwjhGf)上提出構想，以了解大家的興趣。
2. 將[`v3/wep/WEP_TEMPLATE.md`](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/WEP_TEMPLATE.md)複製到`v3/wep/proposals/<proposal name>/proposal.md`，並填寫每個章節。
3. 建立標題為`[WEP] <title>`的草稿 Pull Request，且其中只能包含提案。此 PR 是討論該提案的正式場所。
4. 收集意見與支持（PR 上的留言及按讚回應）。請預留至少兩週進行討論，並就提案的實作者達成共識。
5. 將 PR 標示為可供審查。維護者會做出最終決定：獲接受的提案會取得 WEP 編號並合併。

完整流程記載於[`v3/wep/README.md`](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/README.md)。

## 開始使用

### 建立分支並複製儲存庫

```bash
# Fork the repository on GitHub
# Then clone your fork
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails

# Add upstream remote
git remote add upstream https://github.com/wailsapp/wails.git
```

### 從原始碼建置

```bash
# Build the v3 CLI. Go downloads any required modules automatically.
cd v3
go build -o ../wails3 ./cmd/wails3

# Confirm the built CLI runs and reports its version.
../wails3 version
```

### 執行測試

測試是變更的一部分，而不是最後才執行的冒煙測試。請將單元測試放在其所測試的程式碼旁；當需要以多個輸入或邊界情況檢查同一項行為時，應優先採用表格驅動測試。請為每個測試案例命名，讓失敗訊息能說明其情境。

新增及變更的邏輯應獲得完整涵蓋。您的 PR 所新增或變更的程式碼，Go 陳述式覆蓋率應以100% 為目標；請勿以整個儲存庫的覆蓋率百分比取代對該變更的測試。覆蓋缺口可能有合理原因，例如僅限特定作業系統的錯誤路徑，或不使用實體硬體就難以重現的情況；但請在 PR 說明中解釋該缺口，以及無法以合理方式測試的原因。

```bash
# Run all v3 tests
cd v3
go test ./...

# Run specific package tests
go test ./pkg/application

# Inspect coverage for the packages you changed
go test ./pkg/application -coverprofile=coverage.out
go tool cover -func=coverage.out
cd ..
```

如需整合測試套件、競爭偵測及完整的 CI 等效命令，請參閱[測試與持續整合](/contributing/testing-ci/)。

## 進行變更

### 建立分支

```bash
# Update master
git checkout master
git pull upstream master

# Create feature branch
git checkout -b feature/my-feature
```

### 進行變更

1. 遵循 Go 慣例**撰寫程式碼**
2. 為新功能**新增測試**
3. 視需要**更新文件**
4. **執行測試**，確保沒有破壞任何功能
5. 以清楚的訊息**提交變更**

### 提交準則

```bash
# Good commit messages
git commit -m "fix: resolve window focus issue on macOS"
git commit -m "feat: add support for custom window chrome"
git commit -m "docs: improve bindings documentation"

# Use conventional commits:
# - feat: New feature
# - fix: Bug fix
# - docs: Documentation
# - test: Tests
# - refactor: Code refactoring
# - chore: Maintenance
```

### 提交 Pull Request

```bash
# Push to your fork
git push origin feature/my-feature

# Open pull request on GitHub
# Provide clear description
# Reference related issues
```

## Pull Request 準則

### 良好的 PR 說明

```markdown
## Description
Brief description of changes

## Changes
- Added feature X
- Fixed bug Y
- Updated documentation

## Testing
- Tested on macOS 14
- Tested on Windows 11
- All tests passing

## Related Issues
Fixes #123
```

### PR 檢查清單

- [ ] 程式碼遵循 Go 慣例
- [ ] 已新增或更新測試
- [ ] 已更新文件
- [ ] 所有測試均已通過
- [ ] 沒有破壞性變更（或已記載）
- [ ] 提交訊息清楚明確

## 程式碼準則

### Go 程式碼風格

```go
// ✅ Good: Clear, documented, tested
// ProcessData processes the input data and returns the result.
// It returns an error if the data is invalid.
func ProcessData(data string) (string, error) {
    if data == "" {
        return "", errors.New("data cannot be empty")
    }
    
    result := process(data)
    return result, nil
}

// ❌ Bad: No docs, no error handling
func ProcessData(data string) string {
    return process(data)
}
```

### 測試

```go
func TestProcessData(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid input", "test", "processed", false},
        {"empty input", "", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ProcessData(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ProcessData() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("ProcessData() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## 文件

### 撰寫文件

完整的文件工作流程請參閱[使用 M-Press 撰寫文件](/contributing/documentation/)。其中涵蓋固定版本、可執行檔與 Go 的安裝選項、可編輯的原始檔、導覽、翻譯檢查，以及針對`master`提交 PR。

僅涉及文件的工作不需要建置 Wails 應用程式，也不需要平台相依項目。M-Press 可在儲存庫根目錄使用`mpress dev`預覽現有網站。

### 文件風格

- 使用國際英語拼字
- 先說明問題
- 提供可運作的範例
- 納入疑難排解資訊
- 交叉參照相關內容

## 社群

### 取得協助

- **Discord：**[加入我們的社群](https://discord.gg/JDdSxwjhGf)
- <strong>GitHub Discussions：</strong>提出問題
- <strong>GitHub Issues：</strong>回報錯誤

### 行為準則

請保持尊重、包容與專業。我們齊聚於此，是為了共同打造優秀的軟體。詳情請參閱[行為準則](https://github.com/wailsapp/wails/blob/master/CODE_OF_CONDUCT.md)。

## 貢獻者表揚

貢獻者會在以下位置獲得表揚：

- 版本資訊
- 貢獻者名單
- GitHub Insights

感謝您為 Wails 做出貢獻！🎉

## 後續步驟

@cards{cols="2"}
◆ GitHub 儲存庫
造訪 Wails 儲存庫。

[在 GitHub 上查看 →](https://github.com/wailsapp/wails)

---
◆ Discord 社群
加入社群。

[加入 Discord →](https://discord.gg/JDdSxwjhGf)

---
📖 說明文件
閱讀說明文件。

[瀏覽說明文件 →](/quick-start/why-wails/)

@end
