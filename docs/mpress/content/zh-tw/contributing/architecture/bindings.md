---
title: "繫結系統"
description: "繫結系統如何收集、處理及產生 JavaScript/TypeScript 程式碼"
slug: "contributing/architecture/bindings"
sourcePath: "contributing/architecture/bindings.md"
---

本指南說明 Wails 繫結系統的內部運作方式，協助想要瞭解自動產生程式碼背後機制的開發人員深入掌握其原理。

## 架構概觀

Wails 繫結系統由三個主要元件組成：

1. **收集**：分析 Go 程式碼，以擷取服務、模型及其他宣告的相關資訊
2. **設定**：管理繫結產生程序的設定與選項
3. **轉譯**：根據收集到的資訊產生 JavaScript/TypeScript 程式碼

@filetree
- internal/generator/
  - collect/     # 套件分析與資訊擷取
  - config/      # 設定結構與介面
  - render/      # JS/TS 程式碼產生
@end

## 收集程序

收集程序負責分析 Go 套件，並擷取服務、模型及其他宣告的相關資訊。此程序由`collect`套件處理。

### 主要元件

- **收集器**：管理套件資訊，並快取收集到的資料
- **套件**：表示正在分析的 Go 套件，並儲存收集到的服務、模型及指示詞
- **服務**：收集服務型別及其方法的相關資訊
- **模型**：收集模型型別的詳細資訊，包括欄位、值及型別參數
- **指示詞**：剖析並解讀 Go 原始碼中的`//wails:`指示詞

### 收集流程

1. 收集器會掃描專案中指定的 Go 套件
2. 收集器會識別服務型別（具有方法且將公開給前端的結構）
3. 收集器會收集每個服務的方法資訊
4. 收集器會識別模型型別（在服務方法中用作參數或傳回值的結構）
5. 收集器會收集每個模型的欄位與型別參數資訊
6. 收集器會處理程式碼中找到的所有`//wails:`指示詞

## 轉譯程序

轉譯程序負責根據收集到的資訊產生 JavaScript/TypeScript 程式碼。此程序由`render`套件處理。

### 主要元件

- **轉譯器**：協調服務、模型及索引檔案的轉譯
- **模組**：表示單一產生的 JavaScript/TypeScript 模組
- **範本**：用於產生程式碼的文字範本

### 轉譯流程

1. 轉譯器會為每個服務產生一個 JavaScript/TypeScript 檔案，其中包含與服務方法對應的函式
2. 轉譯器會為每個模型產生一個與模型結構對應的 JavaScript/TypeScript 類別
3. 轉譯器會產生重新匯出所有服務與模型的索引檔案
4. 轉譯器會套用`//wails:inject`指示詞所指定的所有自訂程式碼注入

## 型別對應

繫結系統最重要的層面之一，是如何將 Go 型別對應至 JavaScript/TypeScript 型別。以下是對應關係摘要：

| Go 型別 | JavaScript 型別 | TypeScript 型別 |
| --- | --- | --- |
| `bool` | `boolean` | `boolean` |
| `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64` | `number` | `number` |
| `string` | `string` | `string` |
| `[]byte` | `Uint8Array` | `Uint8Array` |
| `[]T` | `Array<T>` | `T[]` |
| `map[string]V` | `Object` | `{ [_: string]: V }` |
| `map[K]V`（非字串`K`） | `Object` | `{ [_ in K]?: V }` |
| `struct` | `Object` | 自訂類別 |
| `interface{}` | `any` | `any` |
| `*T` | `T \| null` | `T \| null` |
| `func` | 不支援 | 不支援 |
| `chan` | 不支援 | 不支援 |

## 指示詞系統

繫結系統支援多種可用來自訂所產生程式碼的指示詞。這些指示詞會以註解的形式加入 Go 程式碼。

### 可用的指示詞

- `//wails:inject`：將自訂 JavaScript/TypeScript 程式碼注入所產生的繫結
- `//wails:include`：將其他檔案納入所產生的繫結
- `//wails:internal`：隱藏所產生 SDK 中的匯出項目；但不會阻止在執行階段叫用已匯出的 Go 方法
- `//wails:ignore`：在產生繫結時完全忽略某個方法
- `//wails:id`：為方法指定自訂 ID，覆寫預設以雜湊為基礎的 ID

這些指令會控制產生的程式碼。在一般建置中，執行階段註冊仍會反映已匯出的服務方法；私人操作請使用未匯出的方法或未繫結的型別，並在公開的方法中強制執行授權檢查。

### 指令處理

1. 在收集階段，收集器會識別並剖析 Go 程式碼中的指令
2. 這些指令會與對應的宣告（服務、方法、模型等）一起儲存
3. 在算繪階段，算繪器會套用這些指令，以自訂產生的程式碼

## 進階功能

### 條件式程式碼產生

繫結系統支援使用兩個字元的條件前置字元，為`include`與`inject`指令進行條件式程式碼產生：

```
<language><style>:<content>
```

其中：

- `<language>`可以是：
  - `*`－JavaScript 與 TypeScript 兩者
  - `j`－僅 JavaScript
  - `t`－僅 TypeScript


- `<style>`可以是：
  - `*`－類別與介面兩者
  - `c`－僅類別
  - `i`－僅介面


例如：

```go
//wails:inject j*:console.log("JavaScript only");
//wails:inject t*:console.log("TypeScript only");
```

### 自訂方法 ID

依預設，方法會以雜湊式 ID 識別。不過，您可以使用`//wails:id`指令指定自訂 ID：

```go
//wails:id 42
func (s *Service) CustomIDMethod() {}
```

在重構程式碼時，這有助於維持相容性。

## 效能考量

繫結產生器的設計著重效率，但仍有幾點需要注意：

1. 第一次執行時會較慢，因為需要建立待掃描套件的快取
2. 後續執行會使用快取資訊，因此速度較快
3. 產生器會處理專案中的所有套件，對大型專案而言可能相當耗時
4. 您可以使用`-clean`旗標，在產生前清理輸出目錄

## 偵錯

如果繫結產生發生問題，您可以使用`-v`旗標啟用偵錯輸出：

```bash
wails3 generate bindings -v
```

這會提供收集與算繪程序的詳細資訊，有助於找出問題來源。
