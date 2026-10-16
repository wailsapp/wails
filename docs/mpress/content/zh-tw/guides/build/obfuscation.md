---
title: "混淆建置"
description: "使用 Garble 建置 Wails 應用程式，保護原始碼免遭逆向工程"
slug: "guides/build/obfuscation"
sourcePath: "guides/build/obfuscation.md"
---

[Garble](https://github.com/burrowers/garble) 是一項 Go 建置工具，可取代 `go build`，對符號重新命名、混淆常數，並從產生的二進位檔中移除偵錯資訊。Wails v3 透過兩個新命令提供 Garble 的原生支援。

## 先決條件

- **Go 1.26.2 或更新版本** — Garble v0.16.0 的必要條件
- **Garble v0.16.0**

```bash
go install mvdan.cc/garble@v0.16.0
```

@note{type="tip"}
Garble 所需的最低 Go 版本會隨版本變更。如果使用較舊的 Go 工具鏈，請在安裝前查看[Garble 發行版本頁面](https://github.com/burrowers/garble/releases)，找出與工具鏈相符的版本。

@end

## 必要：為服務型別新增 JSON 標記

繫結服務方法所傳回或接受的任何結構，都必須在每個匯出欄位上明確指定 JSON 標記：

```go
// Without tags — breaks under Garble
type OrderSummary struct {
    ID        int
    Total     float64
    LineItems []LineItem
}

// With tags — safe under Garble
type OrderSummary struct {
    ID        int       `json:"id"`
    Total     float64   `json:"total"`
    LineItems []LineItem `json:"lineItems"`
}
```

@note{type="caution"}
Garble 會重新命名匯出的結構欄位，而 Wails 會透過 Garble 無法靜態追蹤的 `interface{}` 參數，將這些結構傳給 `json.Marshal`。沒有 JSON 標記的混淆建置雖然能成功編譯，但前端會在執行階段收到遭混淆或空白的欄位名稱。請先新增標記，再執行混淆建置。

@end

Wails 本身的型別 — `Screen`、`Rect`、`Point`、`Size`、`EnvironmentInfo`、`OSInfo`、`Capabilities` — 已經有標記。您只需為自己的型別新增標記。

## 使用混淆進行建置

@steps
### 產生穩定 ID 檔案
每當新增、重新命名或移除繫結服務方法時，請執行此命令：

```bash
wails3 generate bindings -obfuscated
```

這會在主套件目錄中建立 `wails_obfuscated.gen.go` — 請提交此檔案。

### 使用 Garble 建置
```bash
wails3 build --obfuscated
```

使用經過混淆的繫結建置應用程式。

@end

## 將額外旗標傳遞給 Garble

使用 `--garbleargs` 將選項直接轉送至 `garble`：

```bash
# Obfuscate string literals and reduce binary size
wails3 build --obfuscated --garbleargs "-literals -tiny"

# Reproducible output — same seed produces the same binary
wails3 build --obfuscated --garbleargs "-seed=deadbeef"
```

如需支援旗標的完整清單，請參閱[Garble 文件](https://github.com/burrowers/garble#flags)。

## 進階：將 ID 檔案寫入其他套件

依預設，`wails_obfuscated.gen.go` 會寫入 `main` 套件旁。如果專案將服務放在由 `main` 匯入的子套件中，您可以改用 `-obfuscated-output` 將檔案寫入該處：

```bash
wails3 generate bindings -obfuscated -obfuscated-output ./internal/services
```

@note{type="caution"}
目的地套件必須由 `main` 套件直接或間接匯入，其 `init()` 才會在啟動時執行。如果無法連到該套件，穩定 ID 將永遠不會註冊，繫結呼叫也會失敗（例如在執行階段出現 `binding not found` 錯誤）。

@end

## 疑難排解

### `garble: command not found`

Garble 尚未安裝，或 `$(go env GOPATH)/bin` 不在您的 `PATH` 中。

```bash
go install mvdan.cc/garble@v0.16.0
export PATH="$PATH:$(go env GOPATH)/bin"
```

### 前端收到錯誤或空白的欄位值

服務傳回型別缺少 `json:"..."` 標記。請檢查繫結方法傳回的每個結構，並為每個匯出欄位新增明確的標記。

### 瀏覽器主控台出現 `binding not found` 錯誤

穩定 ID 檔案遺失或未編譯進建置中。請檢查：

- `wails_obfuscated.gen.go` 存在於主套件目錄中（或您傳給 `-obfuscated-output` 的目錄中）
- 您已執行 `wails3 build --obfuscated`，該命令會新增 `wails_obfuscated` 建置標記
- 如果使用了 `-obfuscated-output`，目的地套件已由 `main` 匯入

### Windows Defender 將建置標示為病毒

由於經 Garble 混淆的 Go 二進位檔缺少偵錯符號，且外觀類似經封裝的可執行檔，因此 Windows Defender 會在建置期間透過啟發式偵測將其標示為病毒。建置會失敗並顯示：

```
open C:\Users\...\AppData\Local\Temp\go-build...\a.out.exe: The file contains a virus or potentially unwanted software.
```

將暫存目錄（Go 寫入中繼建置成品的位置）和專案目錄新增至 Defender 的排除清單：

```powershell
Add-MpPreference -ExclusionPath "$env:TEMP"
Add-MpPreference -ExclusionPath "C:\path\to\your\project"
```

這些排除項目只套用於指定路徑，不會在全域停用 Defender。

### 建置因 `unsupported Go version` 而失敗

Garble v0.16.0 需要 Go 1.26.2 或更新版本。請升級 Go，或查閱[Garble 發行版本頁面](https://github.com/burrowers/garble/releases)，找出與工具鏈相容的版本。
