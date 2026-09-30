---
title: "列舉"
description: "從 Go 常數自動產生列舉"
slug: "features/bindings/enums"
sourcePath: "features/bindings/enums.md"
---

## 列舉繫結

Wails v3 繫結產生器會<strong>自動偵測 Go 常數型別，並產生 TypeScript 列舉或 JavaScript const 物件</strong>。不需註冊，也不需設定，只要在 Go 中定義型別與常數，其餘工作便由產生器處理。

@note{type="info"}
與 Wails v2 不同，現在<strong>不需要呼叫`EnumBind`</strong>或手動註冊列舉。產生器會自動從原始碼中找出它們。

@end

## 快速開始

**在 Go 中定義具名型別及其常數：**

```go
type Status string

const (
    StatusActive  Status = "active"
    StatusPending Status = "pending"
    StatusClosed  Status = "closed"
)
```

**在結構或服務方法中使用該型別：**

```go
type Ticket struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Status Status `json:"status"`
}
```

**產生繫結：**

```bash
wails3 generate bindings
```

產生器的輸出會同時報告列舉與模型的數量：

```
3 Enums, 5 Models
```

**在前端中使用：**

```javascript
import { Ticket, Status } from './bindings/changeme/models'

const ticket = new Ticket({
    id: 1,
    title: "Bug report",
    status: Status.StatusActive
})
```

<strong>就這麼簡單！</strong>Go 和 JavaScript/TypeScript 都會強制執行列舉型別檢查。

## 定義列舉

Wails 中的列舉是具有基礎基本型別的<strong>具名型別</strong>，並搭配該型別的<strong>const 宣告</strong>。

### 字串列舉

```go
// Title is a title
type Title string

const (
    // Mister is a title
    Mister Title = "Mr"
    Miss   Title = "Miss"
    Ms     Title = "Ms"
    Mrs    Title = "Mrs"
    Dr     Title = "Dr"
)
```

**產生的 TypeScript：**

```typescript
/**
 * Title is a title
 */
export enum Title {
    /**
     * The Go zero value for the underlying type of the enum.
     */
    $zero = "",

    /**
     * Mister is a title
     */
    Mister = "Mr",
    Miss = "Miss",
    Ms = "Ms",
    Mrs = "Mrs",
    Dr = "Dr",
}
```

**產生的 JavaScript：**

```javascript
/**
 * Title is a title
 * @readonly
 * @enum {string}
 */
export const Title = {
    /**
     * The Go zero value for the underlying type of the enum.
     */
    $zero: "",

    /**
     * Mister is a title
     */
    Mister: "Mr",
    Miss: "Miss",
    Ms: "Ms",
    Mrs: "Mrs",
    Dr: "Dr",
};
```

### 整數列舉

```go
type Priority int

const (
    PriorityLow    Priority = 0
    PriorityMedium Priority = 1
    PriorityHigh   Priority = 2
)
```

**產生的 TypeScript：**

```typescript
export enum Priority {
    /**
     * The Go zero value for the underlying type of the enum.
     */
    $zero = 0,

    PriorityLow = 0,
    PriorityMedium = 1,
    PriorityHigh = 2,
}
```

### 型別別名列舉

Go 型別別名（`=`）也可使用，但產生的輸出略有不同：它會產生型別定義和 const 物件，而不是原生 TypeScript `enum`：

```go
// Age is an integer with some predefined values
type Age = int

const (
    NewBorn    Age = 0
    Teenager   Age = 12
    YoungAdult Age = 18

    // Oh no, some grey hair!
    MiddleAged Age = 50
    Mathusalem Age = 1000 // Unbelievable!
)
```

**產生的 TypeScript：**

```typescript
/**
 * Age is an integer with some predefined values
 */
export type Age = number;

/**
 * Predefined constants for type Age.
 * @namespace
 */
export const Age = {
    NewBorn: 0,
    Teenager: 12,
    YoungAdult: 18,

    /**
     * Oh no, some grey hair!
     */
    MiddleAged: 50,

    /**
     * Unbelievable!
     */
    Mathusalem: 1000,
};
```

**產生的 JavaScript：**

```javascript
/**
 * Age is an integer with some predefined values
 * @typedef {number} Age
 */

/**
 * Predefined constants for type Age.
 * @namespace
 */
export const Age = {
    NewBorn: 0,
    Teenager: 12,
    YoungAdult: 18,

    /**
     * Oh no, some grey hair!
     */
    MiddleAged: 50,

    /**
     * Unbelievable!
     */
    Mathusalem: 1000,
};
```

@note{type="tip"}
**具名型別**（`type Title string`）會產生原生 TypeScript `enum` 宣告，其中包含 `$zero` 成員。 **型別別名**（`type Age = int`）則會產生不含 `$zero` 的 `type` + `const` 命名空間組合。

@end

## `$zero` 值

每個具名型別列舉都包含特殊的 `$zero` 成員，代表基礎型別的<strong>Go 零值</strong>：

| 基礎型別 | `$zero` 值 |
| --- | --- |
| `string` | `""` |
| `int`、`int8`、`int16`、`int32`、`int64` | `0` |
| `uint`、`uint8`、`uint16`、`uint32`、`uint64` | `0` |
| `float32`、`float64` | `0` |
| `bool` | `false` |

當結構欄位使用列舉型別且未提供值時，建構函式會預設使用 `$zero`：

```typescript
export class Person {
    "Title": Title;

    constructor($$source: Partial<Person> = {}) {
        if (!("Title" in $$source)) {
            this["Title"] = Title.$zero;  // defaults to ""
        }
        Object.assign(this, $$source);
    }
}
```

這可確保產生類別時以型別安全的方式初始化，列舉欄位絕不會是 `undefined`。產生 TypeScript 介面時（使用 `-i`），則不會有建構函式，欄位也可像平常一樣不存在。

## 在結構中使用列舉

當結構欄位具有列舉型別時，產生的程式碼會<strong>保留該型別</strong>，而不會退回基本型別：

```go
type Person struct {
    Title Title
    Name  string
    Age   Age
}
```

**產生的 TypeScript：**

```typescript
export class Person {
    "Title": Title;
    "Name": string;
    "Age": Age;

    constructor($$source: Partial<Person> = {}) {
        if (!("Title" in $$source)) {
            this["Title"] = Title.$zero;
        }
        if (!("Name" in $$source)) {
            this["Name"] = "";
        }
        if (!("Age" in $$source)) {
            this["Age"] = 0;
        }

        Object.assign(this, $$source);
    }
}
```

`Title` 欄位的型別是 `Title`，而不是 `string`。因此，IDE 可針對列舉值提供完整的自動完成與型別檢查。

## 來自匯入套件的列舉

完整支援在其他套件中定義的列舉。它們會產生至對應的套件目錄：

```go
// services/types.go
package services

type Title string

const (
    Mister Title = "Mr"
    Miss   Title = "Miss"
    Ms     Title = "Ms"
)
```

```go
// main.go
package main

import "myapp/services"

func (*GreetService) Greet(name string, title services.Title) string {
    return "Hello " + string(title) + " " + name
}
```

`Title` 列舉會產生在 `services` 模型檔案中，且系統會自動解析匯入路徑：

```typescript
// bindings/changeme/services/models.ts
export enum Title {
    $zero = "",
    Mister = "Mr",
    Miss = "Miss",
    Ms = "Ms",
}
```

## 列舉方法

你可以在 Go 中為列舉型別新增方法。這些方法不會影響繫結產生，但可提供實用的伺服器端功能：

```go
type Title string

func (t Title) String() string {
    return string(t)
}

const (
    Mister Title = "Mr"
    Miss   Title = "Miss"
)
```

無論該型別是否具有 Go 方法，產生的列舉都完全相同。

## 註解與文件

產生器會將 Go 註解保留為產生輸出中的 JSDoc：

- <strong>型別註解</strong>會成為列舉的文件註解
- <strong>const 群組註解</strong>會成為區段分隔標記
- <strong>個別 const 註解</strong>會成為成員的文件註解
- <strong>行內註解</strong>會盡可能保留

因此，將游標停留在列舉值上時，IDE 會顯示相關文件。

## 支援的基礎型別

繫結產生器支援以下列 Go 基礎型別建立的列舉：

| Go 型別 | 可作為列舉 |
| --- | :---: |
| `string` | 是 |
| `int`、`int8`、`int16`、`int32`、`int64` | 是 |
| `uint`、`uint8`、`uint16`、`uint32`、`uint64` | 是 |
| `float32`、`float64` | 是 |
| `byte`（`uint8`） | 是 |
| `rune`（`int32`） | 是 |
| `bool` | 是 |
| `complex64`、`complex128` | 否 |

## 限制

產生列舉時，<strong>不</strong>支援以下項目：

- **泛型型別** — 型別參數會使常數無法被偵測
- <strong>具有自訂`json.Marshaler`或`encoding.TextMarshaler`</strong>的型別 — 自訂序列化表示產生的值可能與執行階段行為不符，因此產生器會略過這些型別
- **值無法靜態求值或表示的常數** — 常數在其基礎型別中必須具有已知且可表示的值。標準`iota`模式可以正常運作，因為編譯器會將它們解析為具體值
- **複數型別** — `complex64`和`complex128`不能作為列舉的基礎型別

## 完整範例

**Go：**

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
)

// BackgroundType defines the type of background
type BackgroundType string

const (
    BackgroundSolid    BackgroundType = "solid"
    BackgroundGradient BackgroundType = "gradient"
    BackgroundImage    BackgroundType = "image"
)

type BackgroundConfig struct {
    Type  BackgroundType `json:"type"`
    Value string         `json:"value"`
}

type ThemeService struct{}

func (*ThemeService) GetBackground() BackgroundConfig {
    return BackgroundConfig{
        Type:  BackgroundSolid,
        Value: "#ffffff",
    }
}

func (*ThemeService) SetBackground(config BackgroundConfig) error {
    // Apply background
    return nil
}

func main() {
    app := application.New(application.Options{
        Services: []application.Service{
            application.NewService(&ThemeService{}),
        },
    })
    app.Window.New()
    app.Run()
}
```

**前端（TypeScript）：**

```typescript
import { GetBackground, SetBackground } from './bindings/changeme/themeservice'
import { BackgroundConfig, BackgroundType } from './bindings/changeme/models'

// Get current background
const bg = await GetBackground()

// Check the type using enum values
if (bg.type === BackgroundType.BackgroundSolid) {
    console.log("Solid background:", bg.value)
}

// Set a new background
await SetBackground(new BackgroundConfig({
    type: BackgroundType.BackgroundGradient,
    value: "linear-gradient(to right, #000, #fff)"
}))
```

**前端（JavaScript）：**

```javascript
import { GetBackground, SetBackground } from './bindings/changeme/themeservice'
import { BackgroundConfig, BackgroundType } from './bindings/changeme/models'

// Use enum values for type-safe comparisons
const bg = await GetBackground()

switch (bg.type) {
    case BackgroundType.BackgroundSolid:
        applySolid(bg.value)
        break
    case BackgroundType.BackgroundGradient:
        applyGradient(bg.value)
        break
    case BackgroundType.BackgroundImage:
        applyImage(bg.value)
        break
}
```

## 後續步驟

@cards{cols="2"}
📖 資料模型
結構體、型別對應與模型產生。

[深入瞭解 →](/features/bindings/models/)

---
🚀 方法繫結
將 Go 方法繫結至前端。

[深入瞭解 →](/features/bindings/methods/)

---
⚙ 進階繫結
指令、程式碼注入與自訂 ID。

[深入瞭解 →](/features/bindings/advanced/)

---
✓ 最佳實務
繫結設計模式。

[深入瞭解 →](/features/bindings/best-practices/)

@end

---

<strong>有問題嗎？</strong>請在[Discord](https://discord.gg/JDdSxwjhGf)中提問，或查看[繫結範例](https://github.com/wailsapp/wails/tree/master/v3/examples/binding)。
