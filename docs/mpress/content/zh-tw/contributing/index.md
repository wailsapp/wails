---
title: "參與貢獻"
description: "為 Wails 貢獻"
slug: "contributing"
sourcePath: "contributing/index.md"
---

## 歡迎貢獻者！

歡迎為 Wails 做出貢獻！無論是修正錯誤、新增功能，還是改善文件，我們都非常感謝您的協助。

## 貢獻方式

### 1. 回報問題

發現錯誤了嗎？請[建立議題](https://github.com/wailsapp/wails/issues/new)，並提供：

- 清楚的說明
- 重現步驟
- 預期行為與實際行為
- 系統資訊
- 程式碼範例

### 2. 改善文件

歡迎直接提交修正文件的 PR，不必先建立議題，也不需要有失敗的程式碼測試。  
請依照[修正文件](/contributing/documentation/)的指引，使用 M-Press 預覽並驗證變更。

我們隨時歡迎改善文件：

- 修正錯字與錯誤
- 新增範例
- 釐清說明
- 翻譯內容

### 3. 提交程式碼

透過提取要求貢獻程式碼：

- 錯誤修正
- 新功能
- 效能改善
- 測試

### 4. 提出增強提案（WEP）

新功能及公開行為的變更須採用 Wails 增強提案（WEP）流程。此流程可讓功能開發保持透明，並確保每個獲接受的提案都有實作者。請勿建立功能請求議題。

1. 您可以選擇先在 GitHub Discussions 的[構想](https://github.com/wailsapp/wails/discussions/categories/ideas)分類或[Discord](https://discord.gg/JDdSxwjhGf)上提出構想，以評估社群興趣。
2. 將[`v3/wep/WEP_TEMPLATE.md`](https://github.com/wailsapp/wails/blob/master/v3/wep/WEP_TEMPLATE.md)複製到`v3/wep/proposals/<proposal name>/proposal.md`，並填寫每個章節。
3. 建立標題為`[WEP] <title>`的草稿提取要求，其中只能包含提案。該 PR 是討論提案的正式場所。
4. 收集意見與支持（PR 上的留言與豎起大拇指反應）。請至少保留兩週供大家討論，並就由誰實作提案達成共識。
5. 將 PR 標示為可供審查。維護者會做出最終決定：獲接受的提案將獲指派 WEP 編號並合併。

完整流程記載於[`v3/wep/README.md`](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md)。

## 開始使用

### 建立分支版本並複製儲存庫

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

測試是變更的一部分，而不是最後才進行的冒煙測試。請將單元測試放在其所測試的程式碼旁；若要用多種輸入或邊界案例檢查同一項行為，應優先採用表格驅動測試。請為每個案例命名，讓失敗訊息能說明相應情境。

新增及變更的邏輯應獲得完整涵蓋。PR 新增或變更的程式碼，其 Go 陳述式涵蓋率應以100% 為目標；請勿以整個儲存庫的涵蓋率百分比取代對該變更的測試。涵蓋缺口有時可能合理，例如僅在特定作業系統出現的錯誤路徑，或沒有實體硬體便難以重現的情況；但請在 PR 說明中解釋此缺口，以及無法合理測試的原因。

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

如需整合測試套件、競態偵測，以及與完整 CI 等效的命令，請參閱[測試與持續整合](/contributing/testing-ci/)。

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

1. 依照 Go 慣例<strong>撰寫程式碼</strong>
2. 為新功能<strong>新增測試</strong>
3. 視需要<strong>更新文件</strong>
4. **執行測試**，確保沒有任何功能損壞
5. 以清楚的訊息<strong>提交變更</strong>

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

### 提交提取要求

```bash
# Push to your fork
git push origin feature/my-feature

# Open pull request on GitHub
# Provide clear description
# Reference related issues
```

## 提取要求準則

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
- [ ] 已新增／更新測試
- [ ] 已更新文件
- [ ] 所有測試均通過
- [ ] 沒有破壞性變更（或已加以記載）
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

文件使用 M-Press。請編輯 `docs/mpress/content/` 下的 `.md` 檔案，然後從儲存庫根目錄預覽並驗證：

```bash
go install github.com/leaanthony/mpress/cmd/mpress@v1.0.17
mpress dev
mpress build --strict --no-purge-css
mpress check
# Use `python` on Windows if `python3` is unavailable.
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

### 文件風格

- 使用國際英語拼寫
- 從問題著手
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

我們會在以下項目中列出貢獻者：

- 版本資訊
- 貢獻者名單
- GitHub 洞察報告

感謝您為 Wails 做出貢獻！🎉

## 後續步驟

@cards{cols="2"}
◆ GitHub 儲存庫
造訪 Wails 儲存庫。

[在 GitHub 上檢視 →](https://github.com/wailsapp/wails)

---
◆ Discord 社群
加入社群。

[加入 Discord →](https://discord.gg/JDdSxwjhGf)

---
📖 文件
閱讀文件。

[瀏覽文件 →](/quick-start/why-wails/)

@end
