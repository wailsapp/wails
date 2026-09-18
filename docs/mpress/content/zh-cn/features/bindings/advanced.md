---
title: "高级绑定"
description: "高级绑定技术，包括指令、代码注入和自定义 ID"
slug: "features/bindings/advanced"
sourcePath: "features/bindings/advanced.md"
---

本指南介绍在 Wails v3 中自定义和优化绑定生成过程的高级技术。

## 使用指令自定义生成的代码

### 注入自定义代码

`//wails:inject` 指令允许你将自定义 JavaScript/TypeScript 代码注入生成的绑定：

```go
//wails:inject console.log("Hello from Wails!");
type MyService struct {}

func (s *MyService) Greet(name string) string {
    return "Hello, " + name
}
```

这会将指定代码注入为 `MyService` 服务生成的 JavaScript/TypeScript 文件中。

你还可以使用条件注入来针对特定输出格式：

```go
//wails:inject j*:console.log("Hello JS!");  // JavaScript only
//wails:inject t*:console.log("Hello TS!");  // TypeScript only
```

### 包含其他文件

`//wails:include` 指令允许你在生成绑定时包含其他文件：

```go
//wails:include js/*.js
package mypackage
```

此指令通常用于包文档注释，以便在生成绑定时包含其他 JavaScript/TypeScript 文件。

### 将类型和方法标记为内部使用

`//wails:internal` 指令将类型或方法标记为内部使用，阻止其导出到前端：

```go
//wails:internal
type InternalModel struct {
    Field string
}

//wails:internal
func (s *MyService) InternalMethod() {}
```

这适用于仅供 Go 代码内部使用、不应向前端公开的类型和方法。

### 忽略方法

`//wails:ignore` 指令会在生成绑定时完全忽略某个方法：

```go
//wails:ignore
func (s *MyService) IgnoredMethod() {}
```

这与 `//wails:internal` 类似，但它会完全忽略该方法，而不是将其标记为内部使用。

### 自定义方法 ID

`//wails:id` 指令为方法指定自定义 ID，覆盖默认的基于哈希的 ID：

```go
//wails:id 42
func (s *MyService) CustomIDMethod() {}
```

这有助于在重构代码时保持兼容性。

## 处理复杂类型

### 嵌套结构体

绑定生成器会自动处理嵌套结构体：

```go
type Address struct {
    Street string
    City   string
    State  string
    Zip    string
}

type Person struct {
    Name    string
    Address Address
}

func (s *MyService) GetPerson() Person {
    return Person{
        Name: "John Doe",
        Address: Address{
            Street: "123 Main St",
            City:   "Anytown",
            State:  "CA",
            Zip:    "12345",
        },
    }
}
```

生成的 JavaScript/TypeScript 代码将同时包含 `Person` 和 `Address` 的类。

### 映射和切片

映射和切片也会自动处理：

```go
type Person struct {
    Name       string
    Attributes map[string]string
    Friends    []string
}

func (s *MyService) GetPerson() Person {
    return Person{
        Name: "John Doe",
        Attributes: map[string]string{
            "hair": "brown",
            "eyes": "blue",
        },
        Friends: []string{"Jane", "Bob", "Alice"},
    }
}
```

在 JavaScript 中，映射表示为对象，切片表示为数组。在 TypeScript 中，映射表示为 `Record<K, V>`，切片表示为 `T[]`。

### 泛型类型

绑定生成器支持泛型类型：

```go
type Result[T any] struct {
    Data  T
    Error string
}

func (s *MyService) GetResult() Result[string] {
    return Result[string]{
        Data:  "Hello, World!",
        Error: "",
    }
}
```

生成的 TypeScript 代码将包含 `Result` 的泛型类：

```typescript
export class Result<T> {
    "Data": T;
    "Error": string;

    constructor(source: Partial<Result<T>> = {}) {
        if (!("Data" in source)) {
            this["Data"] = null as any;
        }
        if (!("Error" in source)) {
            this["Error"] = "";
        }

        Object.assign(this, source);
    }

    static createFrom<T>(source: string | object = {}): Result<T> {
        let parsedSource = typeof source === "string" ? JSON.parse(source) : source;
        return new Result<T>(parsedSource as Partial<Result<T>>);
    }
}
```

### 接口

绑定生成器可使用 `-i` 标志生成 TypeScript 接口而非类：

```bash
wails3 generate bindings -ts -i
```

这会为所有模型生成 TypeScript 接口：

```typescript
export interface Person {
    Name: string;
    Attributes: Record<string, string>;
    Friends: string[];
}
```

## 优化绑定生成

### 使用名称而非 ID

默认情况下，绑定生成器使用基于哈希的 ID 调用方法。你可以使用 `-names` 标志改为使用名称：

```bash
wails3 generate bindings -names
```

这会生成通过绑定方法的<strong>完全限定</strong>名称（`<package>.<Service>.<Method>`）调用该方法的代码，并使用按位置命名的参数名 `$0`、`$1`、……：

```javascript
export function Greet($0) {
    let $resultPromise = $Call.ByName("main.GreetService.Greet", $0);
    return $resultPromise;
}
```

这可以提高生成代码的可读性并使其更易于调试，但效率可能略低。

### 捆绑运行时

默认情况下，生成的代码从 `@wailsio/runtime` npm 包导入 Wails 运行时。你可以使用 `-b` 标志，将运行时与生成的代码捆绑在一起：

```bash
wails3 generate bindings -b
```

这会将运行时代码直接包含在生成的文件中，从而不再需要该 npm 包。

### 禁用索引文件

如果不需要索引文件，可以使用 `-noindex` 标志禁止生成这些文件：

```bash
wails3 generate bindings -noindex
```

如果你更喜欢直接从服务和模型各自的文件中导入它们，这会很有用。

## 实际应用示例

### 身份验证服务

下面是一个使用自定义指令的身份验证服务示例：

```go
package auth

//wails:inject console.log("Auth service initialized");
type AuthService struct {
    // Private fields
    users map[string]User
}

type User struct {
    Username string
    Email    string
    Role     string
}

type LoginRequest struct {
    Username string
    Password string
}

type LoginResponse struct {
    Success bool
    User    User
    Token   string
    Error   string
}

// Login authenticates a user
func (s *AuthService) Login(req LoginRequest) LoginResponse {
    // Implementation...
}

// GetCurrentUser returns the current user
func (s *AuthService) GetCurrentUser() User {
    // Implementation...
}

// Internal helper method
//wails:internal
func (s *AuthService) validateCredentials(username, password string) bool {
    // Implementation...
}
```

### 数据处理服务

下面是一个使用泛型类型的数据处理服务示例：

```go
package data

type ProcessingResult[T any] struct {
    Data  T
    Error string
}

type DataService struct {}

// Process processes data and returns a result
func (s *DataService) Process(data string) ProcessingResult[map[string]int] {
    // Implementation...
}

// ProcessBatch processes multiple data items
func (s *DataService) ProcessBatch(data []string) ProcessingResult[[]map[string]int] {
    // Implementation...
}

// Internal helper method
//wails:internal
func (s *DataService) parseData(data string) (map[string]int, error) {
    // Implementation...
}
```

### 条件代码注入

下面是一个针对不同输出格式进行条件代码注入的示例：

```go
//wails:inject j*:/**
//wails:inject j*: * @param {string} arg
//wails:inject j*: * @returns {Promise<void>}
//wails:inject j*: */
//wails:inject j*:export async function CustomMethod(arg) {
//wails:inject t*:export async function CustomMethod(arg: string): Promise<void> {
//wails:inject     await InternalMethod("Hello " + arg + "!");
//wails:inject }
type Service struct{}
```

这会为 JavaScript 和 TypeScript 输出注入不同的代码，并为每种语言提供适当的类型注解。
