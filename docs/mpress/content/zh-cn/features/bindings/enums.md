---
title: "枚举"
description: "根据 Go 常量自动生成枚举"
slug: "features/bindings/enums"
sourcePath: "features/bindings/enums.md"
---

## 枚举绑定

Wails v3 绑定生成器<strong>会自动检测 Go 常量类型，并生成 TypeScript 枚举或 JavaScript const 对象</strong>。无需注册，也无需配置——只需在 Go 中定义类型和常量，其余工作由生成器完成。

@note{type="info"}
与 Wails v2 不同，**无需调用`EnumBind`**，也无需手动注册枚举。生成器会自动从源代码中发现它们。

@end

## 快速开始

**在 Go 中定义一个具名类型及其常量：**

```go
type Status string

const (
    StatusActive  Status = "active"
    StatusPending Status = "pending"
    StatusClosed  Status = "closed"
)
```

**在结构体或服务方法中使用该类型：**

```go
type Ticket struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Status Status `json:"status"`
}
```

**生成绑定：**

```bash
wails3 generate bindings
```

生成器的输出会在模型数量旁报告枚举数量：

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

<strong>就是这么简单！</strong>Go 和 JavaScript/TypeScript 中都会强制使用该枚举类型。

## 定义枚举

Wails 中的枚举是一个具有基础类型的<strong>具名类型</strong>，并配有该类型的<strong>const 声明</strong>。

### 字符串枚举

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

**生成的 TypeScript：**

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

**生成的 JavaScript：**

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

### 整数枚举

```go
type Priority int

const (
    PriorityLow    Priority = 0
    PriorityMedium Priority = 1
    PriorityHigh   Priority = 2
)
```

**生成的 TypeScript：**

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

### 类型别名枚举

Go 类型别名（`=`）同样有效，但生成的输出略有不同——它会生成类型定义和 const 对象，而不是原生 TypeScript `enum`：

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

**生成的 TypeScript：**

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

**生成的 JavaScript：**

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
**具名类型**（`type Title string`）会生成带有`$zero`成员的原生 TypeScript `enum`声明。 **类型别名**（`type Age = int`）会生成不含`$zero`的`type` + `const`命名空间对。

@end

## `$zero`值

每个具名类型枚举都包含一个特殊的`$zero`成员，表示其基础类型的<strong>Go 零值</strong>：

| 基础类型 | `$zero`值 |
| --- | --- |
| `string` | `""` |
| `int`、`int8`、`int16`、`int32`、`int64` | `0` |
| `uint`、`uint8`、`uint16`、`uint32`、`uint64` | `0` |
| `float32`、`float64` | `0` |
| `bool` | `false` |

当结构体字段使用枚举类型且未提供值时，构造函数会默认使用`$zero`：

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

这可确保生成类时进行类型安全的初始化——枚举字段绝不会是`undefined`。生成 TypeScript 接口时（使用`-i`），没有构造函数，字段仍可像往常一样缺省。

## 在结构体中使用枚举

当结构体字段采用枚举类型时，生成的代码会<strong>保留该类型</strong>，而不会回退为基本类型：

```go
type Person struct {
    Title Title
    Name  string
    Age   Age
}
```

**生成的 TypeScript：**

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

`Title`字段的类型是`Title`，而不是`string`。这样，IDE 就能为枚举值提供完整的自动补全和类型检查。

## 来自导入包的枚举

完全支持在其他包中定义的枚举。它们会生成到对应的包目录中：

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

`Title`枚举会生成到`services`模型文件中，并自动解析导入路径：

```typescript
// bindings/changeme/services/models.ts
export enum Title {
    $zero = "",
    Mister = "Mr",
    Miss = "Miss",
    Ms = "Ms",
}
```

## 枚举方法

可以在 Go 中为枚举类型添加方法。这些方法不会影响绑定生成，但可提供实用的服务端功能：

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

无论该类型是否存在 Go 方法，生成的枚举都完全相同。

## 注释和文档

生成器会在生成的输出中将 Go 注释保留为 JSDoc：

- <strong>类型注释</strong>会成为枚举的文档注释
- <strong>const 组注释</strong>会成为分节符
- <strong>各个 const 的注释</strong>会成为成员文档注释
- <strong>行内注释</strong>会尽可能予以保留

因此，将鼠标悬停在枚举值上时，IDE 会显示相应文档。

## 支持的基础类型

绑定生成器支持使用以下 Go 基础类型的枚举：

| Go 类型 | 可用作枚举 |
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

生成枚举时<strong>不</strong>支持以下内容：

- **泛型类型** — 类型参数会导致无法检测常量
- <strong>具有自定义`json.Marshaler`或`encoding.TextMarshaler`</strong>的类型 — 自定义序列化意味着生成的值可能与运行时行为不一致，因此生成器会跳过这些类型
- **值无法静态求值或表示的常量** — 常量在其底层类型中必须具有已知且可表示的值。标准`iota`模式可以正常使用，因为编译器会将它们解析为具体值
- **复数类型** — `complex64`和`complex128`不能作为枚举的底层类型

## 完整示例

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

## 后续步骤

@cards{cols="2"}
📖 数据模型
结构体、类型映射和模型生成。

[了解更多 →](/features/bindings/models/)

---
🚀 方法绑定
将 Go 方法绑定到前端。

[了解更多 →](/features/bindings/methods/)

---
⚙ 高级绑定
指令、代码注入和自定义 ID。

[了解更多 →](/features/bindings/advanced/)

---
✓ 最佳实践
绑定设计模式。

[了解更多 →](/features/bindings/best-practices/)

@end

---

<strong>有疑问？</strong>请在[Discord](https://discord.gg/JDdSxwjhGf)中提问，或查看[绑定示例](https://github.com/wailsapp/wails/tree/master/v3/examples/binding)。
