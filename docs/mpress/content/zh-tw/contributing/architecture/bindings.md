---
title: "繫結系統"
description: "繫結系統如何收集、處理及產生 JavaScript/TypeScript 程式碼"
slug: "contributing/architecture/bindings"
sourcePath: "contributing/architecture/bindings.md"
---

本指南說明 Wails 繫結系統的內部運作方式，協助想瞭解自動產生程式碼背後機制的開發人員掌握其原理。

## 架構概觀

Wails 繫結系統由三個主要元件組成：

1. **收集**：分析 Go 程式碼，以擷取服務、模型及其他宣告的相關資訊
2. **設定**：管理繫結產生程序的設定與選項
3. **轉譯**：根據收集到的資訊產生 JavaScript/TypeScript 程式碼

@filetree
- internal/generator/
  - collect/     # 套件分析與資訊擷取
  - config/      # 設定結構與介面
  - render/      # 產生 JS/TS 程式碼
@end

## 收集程序

收集程序負責分析 Go 套件，並擷取服務、模型及其他宣告的相關資訊。此程序由 `collect` 套件處理。

### 主要元件

- **Collector**：管理套件資訊，並快取收集到的資料
- **Package**：表示正在分析的 Go 套件，並儲存收集到的服務、模型及指示詞
- **Service**：收集服務型別及其方法的相關資訊
- **Model**：收集模型型別的詳細資訊，包括欄位、值及型別參數
- **Directive**：剖析並解讀 Go 原始碼中的 `//wails:` 指示詞

### 收集流程

1. 收集器掃描專案中指定的 Go 套件
2. 收集器識別服務型別（其方法將公開給前端的結構）
3. 收集器會收集每項服務的方法資訊
4. 收集器識別模型型別（在服務方法中用作參數或傳回值的結構）
5. 收集器會收集每個模型的欄位及型別參數資訊
6. 收集器處理在程式碼中找到的所有 `//wails:` 指示詞

## 轉譯程序

轉譯程序負責根據收集到的資訊產生 JavaScript/TypeScript 程式碼。此程序由 `render` 套件處理。

### 主要元件

- **Renderer**：統籌服務、模型及索引檔案的轉譯作業
- **Module**：表示單一產生的 JavaScript/TypeScript 模組
- **Templates**：用於產生程式碼的文字範本

### 轉譯流程

1. 轉譯器會為每項服務產生一個 JavaScript/TypeScript 檔案，其中包含與服務方法對應的函式
2. 轉譯器會為每個模型產生一個與模型結構對應的 JavaScript/TypeScript 類別
3. 轉譯器會產生重新匯出所有服務與模型的索引檔案
4. 轉譯器會套用 `//wails:inject` 指示詞所指定的所有自訂程式碼注入

## 型別對應

繫結系統最重要的層面之一，是如何將 Go 型別對應至 JavaScript/TypeScript 型別。以下是對應關係摘要：

| Go 型別 | JavaScript 型別 | TypeScript 型別 |
| --- | --- | --- |
| `bool` | `boolean` | `boolean` |
| `int`、`int8`、`int16`、`int32`、`int64`、`uint`、`uint8`、`uint16`、`uint32`、`uint64`、`float32`、`float64` | `number` | `number` |
| `string` | `string` | `string` |
| `[]byte` | `Uint8Array` | `Uint8Array` |
| `[]T` | `Array<T>` | `T[]` |
| `map[string]V` | `Object` | `{ [_: string]: V }` |
| `map[K]V`（非字串 `K`） | `Object` | `{ [_ in K]?: V }` |
| `struct` | `Object` | 自訂類別 |
| `interface{}` | `any` | `any` |
| `*T` | `T \| null` | `T \| null` |
| `func` | 不支援 | 不支援 |
| `chan` | 不支援 | 不支援 |

## 指令系統

繫結系統支援多種指令，可用來自訂產生的程式碼。這些指令會以註解形式加入 Go 程式碼中。

### 可用的指令

- `//wails:inject`：將自訂 JavaScript/TypeScript 程式碼注入產生的繫結
- `//wails:include`：在產生的繫結中包含其他檔案
- `//wails:internal`：將型別或方法標記為內部項目，防止其匯出至前端
- `//wails:ignore`：在產生繫結時完全忽略某個方法
- `//wails:id`：為方法指定自訂 ID，覆寫預設的雜湊式 ID

### 指令處理

1. 在收集階段，收集器會識別並剖析 Go 程式碼中的指令
2. 指令會與對應的宣告（服務、方法、模型等）一併儲存
3. 在轉譯階段，轉譯器會套用指令以自訂產生的程式碼

## 進階功能

### 條件式程式碼產生

繫結系統支援條件式程式碼產生，會為`include`和`inject`指令使用兩個字元的條件前置詞：

```
<language><style>:<content>
```

其中：

- `<language>`可以是：
  - `*`－JavaScript 和 TypeScript
  - `j`－僅限 JavaScript
  - `t`－僅限 TypeScript


- `<style>`可以是：
  - `*`－類別和介面
  - `c`－僅限類別
  - `i`－僅限介面


例如：

```go
//wails:inject j*:console.log("JavaScript only");
//wails:inject t*:console.log("TypeScript only");
```

### 自訂方法 ID

方法預設以雜湊式 ID 識別。不過，你可以使用`//wails:id`指令指定自訂 ID：

```go
//wails:id 42
func (s *Service) CustomIDMethod() {}
```

這有助於在重構程式碼時維持相容性。

## 效能考量

繫結產生器以高效率為設計目標，但仍有幾點需要注意：

1. 第一次執行時會建立待掃描套件的快取，因此速度較慢
2. 後續執行會使用快取資訊，因此速度較快
3. 產生器會處理專案中的所有套件；對大型專案而言，這可能相當耗時
4. 你可以使用`-clean`旗標，在產生前清除輸出目錄

## 偵錯

如果產生繫結時遇到問題，可以使用`-v`旗標啟用偵錯輸出：

```bash
wails3 generate bindings -v
```

這會提供收集和轉譯程序的詳細資訊，有助於找出問題來源。
