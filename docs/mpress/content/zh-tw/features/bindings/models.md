---
title: "資料模型"
description: "在 Go 與 JavaScript 之間繫結複雜的資料結構"
slug: "features/bindings/models"
sourcePath: "features/bindings/models.md"
---

## 資料模型繫結

Wails 會<strong>自動從 Go 結構產生 JavaScript/TypeScript 類別</strong>，讓您在後端與前端之間傳遞複雜資料時享有完整的型別安全。只需撰寫 Go 結構並產生繫結，即可取得具有建構函式、型別註解及 JSDoc 註解的完整型別化前端模型。

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
console.log(user.Name)  // Type-safe!
```

<strong>就是這麼簡單！</strong>跨橋接層的完整型別安全。

## 定義模型

### 基本結構

```go
type Person struct {
    Name string
    Age  int
}
```

**產生的 JavaScript：**

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

**產生的 JavaScript：**

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

**JSON 標籤會控制** JavaScript 中的欄位名稱。

### 自訂 JSON 編組

繫結方法的結果會使用 Go 的標準 `encoding/json` 套件進行編組。請注意，只有當正在編組的值可定址時，才會呼叫在指標接收器上宣告的 `MarshalJSON` 方法。這表示自訂編組器可針對傳回切片中的元素呼叫，但不會針對以值傳回的結構呼叫：

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

為確保輸出一致，若該型別可安全複製，請在值接收器上定義 `MarshalJSON`；或者傳回該結構的指標（例如 `*User`）。這是 Go 的 `encoding/json` 行為，而非 Wails 繫結產生機制的差異。

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

**產生的 JavaScript：**

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

<strong>註解會變成 JSDoc！</strong>您的 IDE 會顯示這些註解。

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

**產生的 JavaScript：**

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

**產生的 JavaScript：**

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

**產生的 JavaScript：**

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
| `int`、`int8`、`int16`、`int32`、`int64` | `number` |
| `uint`、`uint8`、`uint16`、`uint32`、`uint64` | `number` |
| `float32`、`float64` | `number` |
| `byte` | `number` |
| `rune` | `number` |

### 複合型別

| Go 型別 | JavaScript/TypeScript | 備註 |
| --- | --- | --- |
| `[]T` | `T[]` | - |
| `[N]T` | `T[]` | - |
| `map[string]T` | `{ [_: string]: T }` | 以字串為鍵的映射 |
| `map[K]V` | `{ [_ in K]?: V }` | 非字串的 `K` 會呈現為對應型別，**而不是** JS `Map` |
| `[]byte` | `string` | 以 base64 編碼 |
| `struct` | `class` / `interface` | 包含欄位 |
| `time.Time` | `any` | 在執行階段序列化為 RFC3339Nano 字串 |
| `*T` | `T \| null` | 指標表示可為 null |
| `any` / `interface{}` | `any` | - |
| `error` | `any` / `Exception` | 作為傳回值時為 Exception，否則為 any |

### 不支援的型別

- `chan T`（通道）
- `func()`（函式）
- 複雜介面（`interface{}`除外）
- 未匯出的欄位（小寫）

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

```bash
wails3 generate bindings -ts
```

**產生的內容：**

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
    console.log(user.createdAt.getTime()) // ✅ Date method
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

繫結產生器會自動偵測具有常數的 Go 具名型別，並產生 TypeScript 列舉或 JavaScript const 物件，其中包括代表 Go 零值的`$zero`成員，並完整保留 JSDoc。

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

如需全面瞭解字串列舉、整數列舉、型別別名、匯入套件的列舉及其限制，請參閱專門的<strong>[列舉](/features/bindings/enums/)</strong>頁面。

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
                <small>Created: ${user.createdAt.toLocaleDateString()}</small>
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
- **新增註解** — 註解會轉換為 JSDoc
- **使用 time.Time** — 會轉換為 Date
- **在 Go 端驗證** — 不要信任前端
- **保持模型簡單** — 僅作為資料容器
- **選用欄位使用指標** — 使用`*string`表示可為 null

### ❌ 請勿這樣做

- **不要為 Go 結構新增方法** — 僅將其用作資料
- **不要使用未匯出的欄位** — 這些欄位不會被繫結
- **不要使用複雜介面** — 不支援此類介面
- **不要忘記 JSON 標籤** — 欄位名稱很重要
- **不要巢狀太深** — 保持簡單

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
★ Go 與前端橋接
瞭解橋接機制。

[深入瞭解 →](/concepts/bridge/)

@end

---

<strong>有問題嗎？</strong>請在[Discord](https://discord.gg/JDdSxwjhGf)中提問，或查看[繫結範例](https://github.com/wailsapp/wails/tree/master/v3/examples/binding)。
