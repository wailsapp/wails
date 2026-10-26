---
title: "資料模型"
description: "在 Go 與 JavaScript 之間繫結複雜的資料結構"
slug: "features/bindings/models"
sourcePath: "features/bindings/models.md"
---

## 資料模型繫結

Wails 會根據 Go 結構**自動產生 JavaScript/TypeScript 類別**，在後端與前端之間傳遞複雜資料時提供完整的型別安全性。只需編寫 Go 結構並產生繫結，即可取得具有建構函式、型別註解及 JSDoc 註解的完整型別前端模型。

## 快速開始

**Go 結構：**

```go
type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"createdAt"`
}

func (s *UserService) GetUser(id int) (*User, error) {
    // Return user
}
```

**產生：**

```bash
wails3 generate bindings
```

**JavaScript：**

```javascript
import { GetUser } from './bindings/changeme/userservice'
import { User } from './bindings/changeme/models'

const user = await GetUser(1)
console.log(user.name)  // Type-safe!
```

<strong>就這麼簡單！</strong>在整個橋接層中享有完整的型別安全性。

## 定義模型

以下模型類別用來說明欄位的結構；請使用實際產生的繫結，而不要複製這些簡化類別作為轉換程式碼。

對於`time.Time`，產生器預設使用`-time-type=string`。若產生的類別需要將時間戳記轉換為 JavaScript `Date`值，請使用`-time-type=Date`。標準 JavaScript／類別範本建置工作會選用`Date`；TypeScript 介面工作則保留`string`。變更模式前，請先檢查產生的工作。Go 時間戳記在 JSON 中仍為 RFC3339Nano 字串。

在 Date 模式下，產生的轉換輔助函式會將這些字串轉換為日期。僅使用`Object.assign`並不會執行此轉換。使用字串繫結時，請先使用`new Date(timestamp)`明確轉換時間戳記，再使用 Date 方法。

### 基本結構

```go
type Person struct {
    Name string
    Age  int
}
```

**JavaScript 範例：**

```javascript
export class Person {
    /** @type {string} */
    Name = ""
    
    /** @type {number} */
    Age = 0
    
    constructor(source = {}) {
        Object.assign(this, source)
    }
    
    static createFrom(source = {}) {
        return new Person(source)
    }
}
```

### 使用 JSON 標籤

```go
type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"createdAt"`
}
```

**JavaScript 範例：**

```javascript
export class User {
    /** @type {number} */
    id = 0
    
    /** @type {string} */
    name = ""
    
    /** @type {string} */
    email = ""
    
    /** @type {Date} */
    createdAt = new Date()
    
    constructor(source = {}) {
        Object.assign(this, source)
    }
}
```

**JSON 標籤會控制 JavaScript 中的欄位名稱**。

### 自訂 JSON 封送處理

繫結方法的回傳結果會使用 Go 的標準`encoding/json`套件進行封送處理。請注意，只有在要封送的值可定址時，才會呼叫在指標接收器上宣告的`MarshalJSON`方法。這表示回傳切片中的元素可以呼叫自訂封送器，但以值回傳的結構則不會：

```go
func (n *Nullable[T]) MarshalJSON() ([]byte, error) {
    if n.Valid {
        return json.Marshal(n.V)
    }
    return []byte("null"), nil
}

func (s *UserService) GetUser() User {
    return User{} // pointer receiver may not be called
}

func (s *UserService) GetUsers() []User {
    return []User{{}} // pointer receiver is called for elements
}
```

為確保輸出一致，若該型別可安全複製，請在值接收器上定義`MarshalJSON`；或回傳指向該結構的指標（例如`*User`）。這是 Go 的`encoding/json`行為，而不是 Wails 繫結產生機制的差異。

### 使用註解

```go
// User represents an application user
type User struct {
    // Unique identifier
    ID int `json:"id"`
    
    // Full name of the user
    Name string `json:"name"`
    
    // Email address (must be unique)
    Email string `json:"email"`
}
```

**JavaScript 範例：**

```javascript
/**
 * User represents an application user
 */
export class User {
    /**
     * Unique identifier
     * @type {number}
     */
    id = 0
    
    /**
     * Full name of the user
     * @type {string}
     */
    name = ""
    
    /**
     * Email address (must be unique)
     * @type {string}
     */
    email = ""
}
```

<strong>註解會轉換成 JSDoc！</strong>您的 IDE 會顯示這些註解。

### 巢狀結構

```go
type Address struct {
    Street  string `json:"street"`
    City    string `json:"city"`
    Country string `json:"country"`
}

type User struct {
    ID      int     `json:"id"`
    Name    string  `json:"name"`
    Address Address `json:"address"`
}
```

**JavaScript 範例：**

```javascript
export class Address {
    /** @type {string} */
    street = ""
    
    /** @type {string} */
    city = ""
    
    /** @type {string} */
    country = ""
}

export class User {
    /** @type {number} */
    id = 0
    
    /** @type {string} */
    name = ""
    
    /** @type {Address} */
    address = new Address()
}
```

**用法：**

```javascript
const user = new User({
    id: 1,
    name: "Alice",
    address: new Address({
        street: "123 Main St",
        city: "Springfield",
        country: "USA"
    })
})
```

### 陣列與切片

```go
type Team struct {
    Name    string   `json:"name"`
    Members []string `json:"members"`
}
```

**JavaScript 範例：**

```javascript
export class Team {
    /** @type {string} */
    name = ""
    
    /** @type {string[]} */
    members = []
}
```

**用法：**

```javascript
const team = new Team({
    name: "Engineering",
    members: ["Alice", "Bob", "Charlie"]
})
```

### 映射

```go
type Config struct {
    Settings map[string]string `json:"settings"`
}
```

**JavaScript 範例：**

```javascript
export class Config {
    /** @type {Record<string, string>} */
    settings = {}
}
```

**用法：**

```javascript
const config = new Config({
    settings: {
        theme: "dark",
        language: "en"
    }
})
```

## 型別對應

### 基本型別

| Go 型別 | JavaScript/TypeScript |
| --- | --- |
| `string` | `string` |
| `bool` | `boolean` |
| `int`, `int8`, `int16`, `int32`, `int64` | `number` |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | `number` |
| `float32`, `float64` | `number` |
| `byte` | `number` |
| `rune` | `number` |

### 複雜型別

| Go 型別 | JavaScript/TypeScript | 備註 |
| --- | --- | --- |
| `[]T` | `T[]` | - |
| `[N]T` | `T[]` | - |
| `map[string]T` | `{ [_: string]: T }` | 以字串為鍵的對映 |
| `map[K]V` | `{ [_ in K]?: V }` | 非字串的`K`會呈現為對映型別，**而不是**JS 的`Map` |
| `[]byte` | `string` | 以 base64 編碼 |
| `struct` | `class` / `interface` | 包含欄位 |
| `time.Time` | `string`（預設）或`Date` | `-time-type`會選擇產生的表示形式；JSON 使用 RFC3339Nano 字串 |
| `*T` | `T \| null` | 指標表示可為 null |
| `any` / `interface{}` | `any` | - |
| `error` | `any` / `Exception` | 若作為回傳值則為 Exception，否則為 any |

### 不支援的型別

- `chan T`（通道）
- `func()`（函式）
- 複雜介面（`interface{}`除外）
- 未匯出的欄位（小寫開頭）

## 使用模型

### 建立執行個體

```javascript
import { User } from './bindings/changeme/models'

// Empty instance
const user1 = new User()

// With data
const user2 = new User({
    id: 1,
    name: "Alice",
    email: "alice@example.com"
})

// From JSON string
const user3 = User.createFrom('{"id":1,"name":"Alice"}')

// From object
const user4 = User.createFrom({ id: 1, name: "Alice" })
```

### 傳遞至 Go

```javascript
import { CreateUser } from './bindings/changeme/userservice'
import { User } from './bindings/changeme/models'

const user = new User({
    name: "Bob",
    email: "bob@example.com"
})

const created = await CreateUser(user)
console.log("Created user:", created.id)
```

### 從 Go 接收

```javascript
import { GetUser } from './bindings/changeme/userservice'

const user = await GetUser(1)

// user is already a User instance
console.log(user.name)
console.log(user.email)
console.log(user.createdAt.toISOString())
```

### 更新模型

```javascript
import { GetUser, UpdateUser } from './bindings/changeme/userservice'

// Get user
const user = await GetUser(1)

// Modify
user.name = "Alice Smith"
user.email = "alice.smith@example.com"

// Save
await UpdateUser(user)
```

## TypeScript 支援

### 產生的 TypeScript

本節說明具有 Date 轉換功能的類別繫結。未加其他選項的`-ts`命令會保留產生器預設的字串表示形式。此處顯示的類別經過簡化；產生器也會產生轉換輔助函式。

```bash
wails3 generate bindings -ts -time-type=Date
```

**產生結果：**

```typescript
/**
 * User represents an application user
 */
export class User {
    /**
     * Unique identifier
     */
    id: number = 0
    
    /**
     * Full name of the user
     */
    name: string = ""
    
    /**
     * Email address (must be unique)
     */
    email: string = ""
    
    /**
     * Account creation date
     */
    createdAt: Date = new Date()
    
    constructor(source: Partial<User> = {}) {
        Object.assign(this, source)
    }
    
    static createFrom(source: string | Partial<User> = {}): User {
        const parsedSource = typeof source === 'string' 
            ? JSON.parse(source) 
            : source
        return new User(parsedSource)
    }
}
```

### 在 TypeScript 中使用

```typescript
import { GetUser, CreateUser } from './bindings/changeme/userservice'
import { User } from './bindings/changeme/models'

async function example() {
    // Create user
    const newUser = new User({
        name: "Alice",
        email: "alice@example.com"
    })
    
    const created: User = await CreateUser(newUser)
    
    // Get user
    const user: User = await GetUser(created.id)
    
    // Type-safe access
    console.log(user.name.toUpperCase())  // ✅ string method
    console.log(user.id + 1)              // ✅ number operation
    console.log(user.createdAt.getTime()) // Date mode selected above
}
```

## 進階模式

### 選用欄位

```go
type User struct {
    ID       int     `json:"id"`
    Name     string  `json:"name"`
    Nickname *string `json:"nickname,omitempty"`
}
```

**JavaScript：**

```javascript
const user = new User({
    id: 1,
    name: "Alice",
    nickname: "Ally"  // Optional
})

// Check if set
if (user.nickname) {
    console.log("Nickname:", user.nickname)
}
```

### 列舉

繫結產生器會自動偵測具有常數的 Go 具名型別，並產生 TypeScript 列舉或 JavaScript const 物件，其中包括代表 Go 零值的`$zero`成員，且會完整保留 JSDoc。

```go
type UserRole string

const (
    RoleAdmin UserRole = "admin"
    RoleUser  UserRole = "user"
    RoleGuest UserRole = "guest"
)
```

**產生的 TypeScript：**

```typescript
export enum UserRole {
    $zero = "",
    RoleAdmin = "admin",
    RoleUser = "user",
    RoleGuest = "guest",
}
```

**用法：**

```javascript
import { User, UserRole } from './bindings/changeme/models'

const admin = new User({
    name: "Admin",
    role: UserRole.RoleAdmin
})
```

如需全面瞭解字串列舉、整數列舉、型別別名、匯入套件中的列舉及其限制，請參閱專門的<strong>[列舉](/features/bindings/enums/)</strong>頁面。

### 驗證

```javascript
class User {
    validate() {
        if (!this.name) {
            throw new Error("Name is required")
        }
        if (!this.email.includes('@')) {
            throw new Error("Invalid email")
        }
        return true
    }
}

// Use
const user = new User({ name: "Alice", email: "alice@example.com" })
user.validate()  // ✅

const invalid = new User({ name: "", email: "invalid" })
invalid.validate()  // ❌ Throws
```

### 序列化

```javascript
// To JSON
const json = JSON.stringify(user)

// From JSON
const user = User.createFrom(json)

// To plain object
const obj = { ...user }

// From plain object
const user2 = new User(obj)
```

## 完整範例

**Go：**

```go
package main

import (
    "time"
    "github.com/wailsapp/wails/v3/pkg/application"
)

type Address struct {
    Street  string `json:"street"`
    City    string `json:"city"`
    Country string `json:"country"`
}

type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Address   Address   `json:"address"`
    CreatedAt time.Time `json:"createdAt"`
}

type UserService struct {
    users []User
}

func (s *UserService) GetAll() []User {
    return s.users
}

func (s *UserService) GetByID(id int) (*User, error) {
    for _, user := range s.users {
        if user.ID == id {
            return &user, nil
        }
    }
    return nil, fmt.Errorf("user %d not found", id)
}

func (s *UserService) Create(user User) User {
    user.ID = len(s.users) + 1
    user.CreatedAt = time.Now()
    s.users = append(s.users, user)
    return user
}

func (s *UserService) Update(user User) error {
    for i, u := range s.users {
        if u.ID == user.ID {
            s.users[i] = user
            return nil
        }
    }
    return fmt.Errorf("user %d not found", user.ID)
}

func main() {
    app := application.New(application.Options{
        Services: []application.Service{
            application.NewService(&UserService{}),
        },
    })
    
    app.Window.New()
    app.Run()
}
```

**JavaScript：**

```javascript
import { GetAll, GetByID, Create, Update } from './bindings/changeme/userservice'
import { User, Address } from './bindings/changeme/models'

class UserManager {
    async loadUsers() {
        const users = await GetAll()
        this.renderUsers(users)
    }
    
    async createUser(name, email, address) {
        const user = new User({
            name,
            email,
            address: new Address(address)
        })
        
        try {
            const created = await Create(user)
            console.log("Created user:", created.id)
            this.loadUsers()
        } catch (error) {
            console.error("Failed to create user:", error)
        }
    }
    
    async updateUser(id, updates) {
        try {
            const user = await GetByID(id)
            Object.assign(user, updates)
            await Update(user)
            this.loadUsers()
        } catch (error) {
            console.error("Failed to update user:", error)
        }
    }
    
    renderUsers(users) {
        const list = document.getElementById('users')
        list.innerHTML = users.map(user => `
            <div class="user">
                <h3>${user.name}</h3>
                <p>${user.email}</p>
                <p>${user.address.city}, ${user.address.country}</p>
                <small>Created: ${new Date(user.createdAt).toLocaleDateString()}</small>
            </div>
        `).join('')
    }
}

const manager = new UserManager()
manager.loadUsers()
```

## 最佳實務

### ✅ 建議做法

- **使用 JSON 標籤** — 控制欄位名稱
- **新增註解** — 註解會成為 JSDoc
- **選擇時間表示形式** — `time.Time`預設會轉為`string`；`-time-type=Date`會啟用產生的 Date 轉換功能
- **在 Go 端進行驗證** — 不要信任前端
- **保持模型簡潔** — 僅作為資料容器
- **選用值使用指標** — 使用`*string`表示可為 null

### ❌ 避免做法

- **不要為 Go 結構新增方法** — 讓它們只承載資料
- **不要使用未匯出的欄位** — 這些欄位不會被繫結
- **不要使用複雜介面** — 不受支援
- **不要忘記 JSON 標籤** — 欄位名稱很重要
- **不要巢狀太多層** — 保持簡潔

## 後續步驟

@cards{cols="2"}
🚀 方法繫結
瞭解如何繫結 Go 方法。

[深入瞭解 →](/features/bindings/methods/)

---
◆ 服務
使用服務組織程式碼。

[深入瞭解 →](/features/bindings/services/)

---
✓ 最佳實務
繫結設計模式。

[深入瞭解 →](/features/bindings/best-practices/)

---
★ Go 與前端的橋接
瞭解橋接機制。

[深入瞭解 →](/concepts/bridge/)

@end

---

<strong>有問題嗎？</strong>請到[Discord](https://discord.gg/JDdSxwjhGf)提問，或查看[繫結範例](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples/binding)。
