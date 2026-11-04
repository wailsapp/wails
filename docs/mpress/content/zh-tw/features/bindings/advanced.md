---
title: "進階繫結"
description: "進階繫結技術，包括指示詞、程式碼注入和自訂 ID"
slug: "features/bindings/advanced"
sourcePath: "features/bindings/advanced.md"
---

本指南涵蓋在 Wails v3 中自訂及最佳化繫結產生流程的進階技術。

## 使用指示詞自訂產生的程式碼

### 注入自訂程式碼

`//wails:inject` 指示詞可讓你將自訂 JavaScript/TypeScript 程式碼注入產生的繫結：

```go
//wails:inject console.log("Hello from Wails!");
type MyService struct {}

func (s *MyService) Greet(name string) string {
    return "Hello, " + name
}
```

這會將指定的程式碼注入為 `MyService` 服務產生的 JavaScript/TypeScript 檔案。

你也可以使用條件式注入，以特定輸出格式為目標：

```go
//wails:inject j*:console.log("Hello JS!");  // JavaScript only
//wails:inject t*:console.log("Hello TS!");  // TypeScript only
```

### 包含其他檔案

`//wails:include` 指示詞可讓你在產生繫結時包含其他檔案：

```go
//wails:include js/*.js
package mypackage
```

此指示詞通常用於套件文件註解，以便在產生繫結時包含其他 JavaScript/TypeScript 檔案。

### 將型別和方法標記為內部使用

`//wails:internal` 指示詞會禁止在產生的 JavaScript/TypeScript 繫結中匯出型別或方法：

```go
//wails:internal
type InternalModel struct {
    Field string
}

//wails:internal
func (s *MyService) InternalMethod() {}
```

這會控制產生的 SDK，而非執行階段存取。在一般建置中，即使已匯出的 Go 服務方法未出現在產生的 SDK 中，仍可透過執行階段註冊和呼叫。請將私有操作保持為未匯出狀態，或將其放在未繫結的型別上，並在公開的方法內強制執行授權檢查。

### 忽略方法

`//wails:ignore` 指示詞會在產生繫結時完全忽略某個方法：

```go
//wails:ignore
func (s *MyService) IgnoredMethod() {}
```

這會省略產生的包裝函式。和 `//wails:internal` 一樣，它並不是已匯出 Go 服務方法的存取控制界線。

### 自訂方法 ID

`//wails:id` 指示詞可為方法指定自訂 ID，覆寫預設的雜湊式 ID：

```go
//wails:id 42
func (s *MyService) CustomIDMethod() {}
```

這有助於在重構程式碼時維持相容性。

## 使用複雜型別

### 巢狀結構

繫結產生器會自動處理巢狀結構：

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

### 映射和切片

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

在 JavaScript 中，映射以物件表示，切片則以陣列表示。在 TypeScript 中，映射以 `Record<K, V>` 表示，切片則以 `T[]` 表示。

### 泛型型別

繫結產生器支援泛型型別：

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

繫結產生器可使用 `-i` 旗標產生 TypeScript 介面，而非類別：

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

## 最佳化繫結產生流程

### 使用名稱而非 ID

繫結產生器預設使用雜湊式 ID 進行方法呼叫。你可以使用 `-names` 旗標，改用名稱：

```bash
wails3 generate bindings -names
```

這會產生程式碼，透過繫結方法的**完整限定**名稱（`<package>.<Service>.<Method>`）呼叫該方法，並使用依位置命名的 `$0`、`$1`、… 參數名稱：

```javascript
export function Greet($0) {
    let $resultPromise = $Call.ByName("main.GreetService.Greet", $0);
    return $resultPromise;
}
```

這可讓產生的程式碼更易讀、也更容易偵錯，但效率可能略低。

### 封裝執行階段

產生的程式碼預設會從 `@wailsio/runtime` npm 套件匯入 Wails 執行階段。你可以使用 `-b` 旗標，將執行階段與產生的程式碼封裝在一起：

```bash
wails3 generate bindings -b
```

這會將執行階段程式碼直接包含在產生的檔案中，因此不再需要該 npm 套件。

### 停用索引檔案

如果不需要索引檔案，可以使用 `-noindex` 旗標停用其產生功能：

```bash
wails3 generate bindings -noindex
```

如果你偏好直接從服務和模型各自的檔案匯入，這會很有用。

## 實務範例

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

這會分別為 JavaScript 和 TypeScript 輸出注入不同的程式碼，並為每種語言提供適當的型別註解。
