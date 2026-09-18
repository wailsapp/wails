---
title: "繫結系統"
description: "Wails v3 如何讓 Go 與 JavaScript 無須樣板程式碼即可互相呼叫"
slug: "contributing/binding-system"
sourcePath: "contributing/binding-system.md"
---

> 「繫結」是讓你能撰寫以下程式碼的<strong>型別安全合約</strong>：

```go
msg, err := chatService.Send("Hello")
```

在 Go 中<em>，以及</em>

```ts
import { Send } from "../bindings/github.com/you/yourapp/services/chatservice";

const msg = await Send("Hello");
```

在 TypeScript 中<strong>，完全不必手動撰寫任何 IPC 銜接程式碼</strong>。 本文將詳細說明這項機制<em>如何</em>運作：從建置階段的 **靜態分析**，經過<strong>程式碼產生</strong>，一直到透過 WebView 傳輸 位元組的<strong>執行階段橋接器</strong>。

> 如需以下內容，請參閱[`contributing/architecture/bindings`](/contributing/architecture/bindings/)：
>
> 產生器管線的權威深入解析——本頁則是
>
> 以貢獻者為對象的概覽。

---

## 1. 30 秒概覽

| 階段 | 元件 | 輸出 |
| --- | --- | --- |
| **收集／分析** | `internal/generator/collect/`、`internal/generator/analyse.go` | 已匯出 Go 服務、方法、參數、回傳型別及模型的記憶體內模型 |
| **產生** | `internal/generator/render/templates/*.tmpl`（`service.{js,ts}.tmpl`、`models.{js,ts}.tmpl`、`index.tmpl`、`eventcreate.js.tmpl`、`eventdata.d.ts.tmpl`、`newline.tmpl`） | 位於`frontend/bindings/<full Go import path>/...`下、各服務專屬的 ES 模組 |
| **執行階段** | `pkg/application/messageprocessor*.go`，以及位於`internal/runtime/desktop/@wailsio/runtime/src/`下的內嵌 JS 執行階段（`calls.ts`、`events.ts`……） | 透過 WebView 原生橋接器傳輸的呼叫／事件訊息 |

此流程由`wails3 generate bindings`命令協調；該命令會針對一組 Go 套件， 驅動`generator.Generate`（定義於`internal/generator/generate.go`）。

```
wails3 generate bindings
        │
        ▼
internal/generator/generate.go: Generator.Generate(patterns...)
        │
        ├── internal/generator/collect/   // load.go, collector.go, service.go, model.go, …
        ├── internal/generator/analyse.go // semantic checks
        └── internal/generator/render/    // template execution → frontend/bindings/**
```

---

## 2. 靜態分析

### 進入點

```
internal/generator/generate.go        // Generator + Generate(patterns…)
internal/generator/analyse.go         // semantic validation
internal/generator/collect/load.go    // go/packages loader
internal/generator/collect/collector.go
```

收集器階段會逐一巡覽所有已載入的套件，並記錄：

- `collect.ServiceInfo`——每個已匯出且已繫結的 Go 結構各有一個。
- `collect.ServiceMethodInfo`／`collect.MethodInfo`——各方法的簽章資訊（名稱、參數、結果、錯誤位置、接收者、文件）。
- `collect.ModelInfo`／`collect.StructInfo`——產生為 TS／JS 模型。
- 指令註解，例如`//wails:inject`、`//wails:include`、`//wails:internal`、`//wails:ignore`、`//wails:id <hex>`（請參閱`internal/generator/collect/directive.go`）。

不支援的型別會導致產生器錯誤，讓錯誤在建置階段浮現， 而非等到執行階段才發生。

### 模型識別碼

執行階段呼叫信封會使用方法完整限定名稱（`pkg.Struct.Method`）的 <strong>確定性 FNV-1a 雜湊</strong>來識別該方法。在產生的繫結中， 你會看到它表示為`$Call.ByID(<numeric-id>, …)`；若產生作業搭配`-names`執行， 則會表示為`$Call.ByName("pkg.Struct.Method", …)`。

---

## 3. 程式碼產生

### 範本

`internal/generator/render/templates/`：

| 範本 | 用途 |
| --- | --- |
| `service.js.tmpl` | 每個已繫結服務各有一個 JS 模組 |
| `service.ts.tmpl` | TypeScript 配套項目（使用`-ts`選項時產生） |
| `models.js.tmpl` | 模型類別輸出（各套件） |
| `models.ts.tmpl` | 模型`.d.ts`輸出（各套件） |
| `index.tmpl` | 各套件的`index.{js,ts}`彙總重新匯出 |
| `eventcreate.js.tmpl`／`eventdata.d.ts.tmpl` | 事件建構函式／承載資料型別定義 |
| `newline.tmpl` | 結尾換行字元正規化器 |

輸出會放在`frontend/bindings/<full Go import path>/...`下——例如， 定義於`github.com/you/yourapp/services/chat`的服務會放到 `frontend/bindings/github.com/you/yourapp/services/chat/`下。v3 中沒有 `frontend/src/wailsjs/`目錄。

### JavaScript 輸出

產生的繫結是 ES 模組，會從 `/wails/runtime.js`匯入執行階段輔助函式：

```js
import { Call as $Call, Create as $Create } from "/wails/runtime.js";

/**
 * @param {string} msg
 * @returns {Promise<string> & { cancel(): void }}
 */
export function Send(msg) {
    return $Call.ByID(2042131923, msg);
}
```

當產生作業搭配`-names`執行時，會改為產生`$Call.ByName("pkg.Struct.Method", ...)`—— 一律為<strong>完整限定</strong>，絕不會只有`"Method"`。

產生的模型類別會使用`$$source`建構函式模式，包含各欄位的 `if (!("X" in $$source))`預設值、加上引號的欄位名稱，以及一個會對字串輸入執行 `JSON.parse`的`static createFrom(...)`。

### 型別對應重點

已根據`internal/generator/render/`驗證：

| Go | TypeScript |
| --- | --- |
| `map[string]V` | `{ [_: string]: V }` |
| `map[K]V`（非字串`K`） | `{ [_ in K]?: V }`（不是`Map<K, V>`，也不是`Record<K, V>`） |
| `[]byte` | `Uint8Array` |
| `time.Time` | `string`（JSON ISO 8601） |
| `error`（回傳位置） | 遭拒絕的 Promise |

### 反射注意事項

`pkg/application/bindings.go`是<strong>手寫的</strong>，並使用`reflect`，根據`BoundMethod`註冊表進行方法分派。請勿過度從字面解讀舊有的「執行階段零反射」說法——產生器會避免使用反射，但執行階段分派器仍會使用反射。

---

## 4. 執行階段呼叫通訊協定

### JavaScript 端

```ts
import { Call } from "/wails/runtime.js";

await Call.ByID(0x7a1201d3 /* ChatService.Send */, "Hello");
// or, with -names:
await Call.ByName("chatservice.ChatService.Send", "Hello");
```

執行階段輔助程式位於`internal/runtime/desktop/@wailsio/runtime/src/calls.ts`（呼叫分派）、`events.ts`（事件）及相關檔案中——此工作樹中沒有`invoke.ts`或`errors.ts`。確切的線路封裝格式由 JS 端的`calls.ts`編碼，並由 Go 端的`pkg/application/messageprocessor_call.go`解碼；偵錯橋接器時，請一併查閱這兩個檔案。

### Go 端

1. `pkg/application/messageprocessor_call.go`接收呼叫訊息。
2. 在`pkg/application/bindings.go`中依 ID 或名稱查找繫結的方法（由`reflect`驅動）。
3. 呼叫繫結的方法，並將`{result, error}`序列化後傳回 JS。

### 錯誤對應

| Go | JavaScript |
| --- | --- |
| `error == nil` | `Promise`以結果完成 |
| `error != nil` | `Promise`以`Error`拒絕，其`message`為 Go 錯誤字串 |

---

## 5. 從 Go 呼叫 JavaScript

繫結產生器是單向的（將 Go 方法公開給 JS）。若要從 Go → JS 通訊，請使用事件匯流排，或在視窗中執行 JS：

```go
app.Event.Emit("chat:new-message", msg)
window.ExecJS(`window.dispatchEvent(new CustomEvent("ping"))`)
```

在 JS 端，使用`/wails/runtime.js`中的`Events.On(name, cb)`訂閱。

---

## 6. 擴充與疑難排解

### 不支援的型別錯誤

```
error: field "Client" uses unsupported type: chan struct{}
```

→ 將通道封裝在方法 API 後方，或以`//wails:internal`標記該欄位，讓產生器略過它。

### 過期的繫結

每次執行`wails3 generate bindings`／`wails3 dev`／`wails3 build`時，都會覆寫產生的輸出。若 IDE 的 IntelliSense 顯示過期的虛設常式，請刪除`frontend/bindings/`並重新執行產生器。`-clean`旗標（目前組建中的預設值為`true`）會在每次執行前清空繫結目錄。

### 效能提示

- 避免透過橋接器串流大型位元組切片——請改由資產伺服器提供。
- 延遲至關重要時，請將多個快速呼叫批次合併為一個方法。
- 小型參數結構請優先使用值接收器，以減少配置。

---

## 7. 關鍵檔案一覽

| 關注項目 | 檔案 |
| --- | --- |
| 產生器協調 | `internal/generator/generate.go` |
| 語意檢查 | `internal/generator/analyse.go` |
| 收集（服務、方法、模型） | `internal/generator/collect/{service,method,model,struct,package}.go` |
| 轉譯範本 | `internal/generator/render/templates/*.tmpl` |
| 產生的繫結位置 | `frontend/bindings/<full Go import path>/...` |
| Go 端分派器 | `pkg/application/bindings.go`、`messageprocessor_call.go` |
| JS 執行階段 | `internal/runtime/desktop/@wailsio/runtime/src/{calls,events,index}.ts` |

追查橋接器錯誤時，請將這份速查表放在手邊。

---

## 8. 重點回顧

1. <strong>收集器</strong>掃描 Go 程式碼 → 記憶體內語意模型。
2. <strong>範本</strong>會為每項服務輸出 ES 模組，並為每個套件輸出模型／索引檔案。
3. <strong>訊息處理器</strong>透過繫結註冊表在 Go 端分派呼叫。
4. <strong>JS 執行階段</strong>將這一切封裝成符合慣用寫法且支援取消的 Promise。

完全不必自行撰寫任何 IPC 樣板程式碼。這就是 Wails v3 繫結系統。現在就開始繫結吧！
