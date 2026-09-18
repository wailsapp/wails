---
title: "進階綁定"
description: "進階綁定技術，包括指示詞、程式碼注入和自訂 ID"
slug: "features/bindings/advanced"
sourcePath: "features/bindings/advanced.md"
---

本指南介紹在 Wails v3 中自訂及最佳化綁定產生流程的進階技術。

## 使用指示詞自訂產生的程式碼

### 注入自訂程式碼

`//wails:inject` 指示詞可讓你將自訂 JavaScript/TypeScript 程式碼注入產生的綁定：

```go
//wails:inject console.log("Hello from Wails!");
type MyService struct {}

func (s *MyService) Greet(name string) string {
    return "Hello, " + name
}
```

這會將指定的程式碼注入為 `MyService` 服務產生的 JavaScript/TypeScript 檔案中。

你也可以使用條件式注入，以指定特定的輸出格式：

```go
//wails:inject j*:console.log("Hello JS!");  // JavaScript only
//wails:inject t*:console.log("Hello TS!");  // TypeScript only
```

### 包含其他檔案

`//wails:include` 指示詞可讓你在產生綁定時包含其他檔案：

```go
//wails:include js/*.js
package mypackage
```

此指示詞通常用於套件文件註解，以便在產生綁定時包含其他 JavaScript/TypeScript 檔案。

### 將型別和方法標記為內部使用

`//wails:internal` 指示詞會將型別或方法標記為內部使用，防止其匯出至前端：

```go
//wails:internal
type InternalModel struct {
    Field string
}

//wails:internal
func (s *MyService) InternalMethod() {}
```

這適用於僅供 Go 程式碼內部使用、不應公開給前端的型別和方法。

### 忽略方法

`//wails:ignore` 指示詞會在產生綁定時完全忽略某個方法：

```go
//wails:ignore
func (s *MyService) IgnoredMethod() {}
```

這與 `//wails:internal` 類似，但它會完全忽略該方法，而不是將其標記為內部使用。

### 自訂方法 ID

`//wails:id` 指示詞可為方法指定自訂 ID，覆寫預設的雜湊式 ID：

```go
//wails:id 42
func (s *MyService) CustomIDMethod() {}
```

這有助於在重構程式碼時維持相容性。

## 處理複雜型別

### 巢狀結構體

綁定產生器會自動處理巢狀結構體：

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

產生的 JavaScript/TypeScript 程式碼會同時包含 `Person` 和 `Address` 的類別。

### 映射與切片

映射和切片也會自動處理：

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

在 JavaScript 中，映射以物件表示，切片以陣列表示。在 TypeScript 中，映射以 `Record<K, V>` 表示，切片以 `T[]` 表示。

### 泛型型別

綁定產生器支援泛型型別：

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

產生的 TypeScript 程式碼會包含 `Result` 的泛型類別：

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

### 介面

綁定產生器可使用 `-i` 旗標產生 TypeScript 介面，而不是類別：

```bash
wails3 generate bindings -ts -i
```

這會為所有模型產生 TypeScript 介面：

```typescript
export interface Person {
    Name: string;
    Attributes: Record<string, string>;
    Friends: string[];
}
```

## 最佳化綁定產生流程

### 使用名稱而非 ID

綁定產生器預設使用雜湊式 ID 呼叫方法。你可以使用 `-names` 旗標，改用名稱：

```bash
wails3 generate bindings -names
```

這會產生依繫結方法的 <strong>完整限定</strong>名稱（`<package>.<Service>.<Method>`）呼叫該方法的程式碼，並使用依位置命名的 `$0`、`$1`……參數名稱：

```javascript
export function Greet($0) {
    let $resultPromise = $Call.ByName("main.GreetService.Greet", $0);
    return $resultPromise;
}
```

這可讓產生的程式碼更容易閱讀和偵錯，但效率可能會稍低。

### 封裝執行階段

產生的程式碼預設會從 `@wailsio/runtime` npm 套件匯入 Wails 執行階段。你可以使用 `-b` 旗標，將執行階段封裝到產生的程式碼中：

```bash
wails3 generate bindings -b
```

這會將執行階段程式碼直接納入產生的檔案中，因而不再需要 npm 套件。

### 停用索引檔案

如果不需要索引檔案，可以使用 `-noindex` 旗標停用其產生：

```bash
wails3 generate bindings -noindex
```

如果你偏好直接從各自的檔案匯入服務和模型，這會很實用。

## 實際案例

### 驗證服務

以下是使用自訂指示詞的驗證服務範例：

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

### 資料處理服務

以下是使用泛型型別的資料處理服務範例：

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

### 條件式程式碼注入

以下是針對不同輸出格式進行條件式程式碼注入的範例：

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

這會針對 JavaScript 與 TypeScript 輸出注入不同的程式碼，並為各語言提供適當的型別註記。
