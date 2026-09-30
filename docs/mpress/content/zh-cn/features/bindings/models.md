---
title: "数据模型"
description: "在 Go 与 JavaScript 之间绑定复杂数据结构"
slug: "features/bindings/models"
sourcePath: "features/bindings/models.md"
---

## 数据模型绑定

Wails **会根据 Go 结构体自动生成 JavaScript/TypeScript 类**，在后端与前端之间传递复杂数据时提供完整的类型安全保障。编写 Go 结构体并生成绑定，即可获得具有构造函数、类型注解和 JSDoc 注释的完全类型化前端模型。

## 快速入门

**Go 结构体：**

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

**生成：**

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

<strong>就是这样！</strong>桥接两端均具有完整的类型安全保障。

## 定义模型

### 基本结构体

```go
type Person struct {
    Name string
    Age  int
}
```

**生成的 JavaScript：**

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

### 使用 JSON 标签

```go
type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"createdAt"`
}
```

**生成的 JavaScript：**

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

**JSON 标签控制 JavaScript 中的字段名称**。

### 自定义 JSON 编组

绑定方法的结果使用 Go 的标准`encoding/json`包进行编组。请注意，只有当被编组的值可寻址时，才会调用在指针接收者上声明的`MarshalJSON`方法。这意味着，对于返回的切片中的元素，可以调用自定义编组器，但对于按值返回的结构体则不会调用：

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

为确保输出一致，如果该类型能够安全复制，请在值接收者上定义`MarshalJSON`；或者返回指向该结构体的指针（例如`*User`）。这是 Go 的`encoding/json`所具有的行为，而不是 Wails 绑定生成方式的差异。

### 使用注释

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

**生成的 JavaScript：**

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

<strong>注释会转换为 JSDoc！</strong>你的 IDE 会显示这些注释。

### 嵌套结构体

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

**生成的 JavaScript：**

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

### 数组和切片

```go
type Team struct {
    Name    string   `json:"name"`
    Members []string `json:"members"`
}
```

**生成的 JavaScript：**

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

**生成的 JavaScript：**

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

## 类型映射

### 基本类型

| Go 类型 | JavaScript/TypeScript |
| --- | --- |
| `string` | `string` |
| `bool` | `boolean` |
| `int`、`int8`、`int16`、`int32`、`int64` | `number` |
| `uint`、`uint8`、`uint16`、`uint32`、`uint64` | `number` |
| `float32`、`float64` | `number` |
| `byte` | `number` |
| `rune` | `number` |

### 复杂类型

| Go 类型 | JavaScript/TypeScript | 说明 |
| --- | --- | --- |
| `[]T` | `T[]` | - |
| `[N]T` | `T[]` | - |
| `map[string]T` | `{ [_: string]: T }` | 以字符串为键的映射 |
| `map[K]V` | `{ [_ in K]?: V }` | 非字符串`K`会呈现为映射类型，<strong>而不是</strong>JS 的`Map` |
| `[]byte` | `string` | 使用 base64 编码 |
| `struct` | `class` / `interface` | 包含字段 |
| `time.Time` | `any` | 在运行时序列化为 RFC3339Nano 字符串 |
| `*T` | `T \| null` | 指针表示可为 null |
| `any` / `interface{}` | `any` | - |
| `error` | `any` / `Exception` | 作为返回值时为 Exception，否则为 any |

### 不支持的类型

- `chan T`（通道）
- `func()`（函数）
- 复杂接口（`interface{}`除外）
- 未导出的字段（小写）

## 使用模型

### 创建实例

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

### 传递给 Go

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

### 从 Go 接收

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

## TypeScript 支持

### 生成的 TypeScript

```bash
wails3 generate bindings -ts
```

**生成结果：**

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

## 高级模式

### 可选字段

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

### 枚举

绑定生成器会自动检测带常量的 Go 命名类型，并生成 TypeScript 枚举或 JavaScript const 对象，其中包括表示 Go 零值的`$zero`成员，并完整保留 JSDoc。

```go
type UserRole string

const (
    RoleAdmin UserRole = "admin"
    RoleUser  UserRole = "user"
    RoleGuest UserRole = "guest"
)
```

**生成的 TypeScript：**

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

有关字符串枚举、整数枚举、类型别名、导入包中的枚举以及相关限制的全面说明，请参阅专门的<strong>[枚举](/features/bindings/enums/)</strong>页面。

### 验证

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

## 完整示例

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

## 最佳实践

### ✅ 应该做

- **使用 JSON 标签**——控制字段名称
- **添加注释**——注释会转换为 JSDoc
- **使用 time.Time**——它会转换为 Date
- **在 Go 端进行验证**——不要信任前端
- **保持模型简单**——仅用作数据容器
- **对可选字段使用指针**——使用`*string`表示可为 null

### ❌ 不要做

- **不要向 Go 结构体添加方法**——让它们仅用于存放数据
- **不要使用未导出的字段**——它们不会被绑定
- **不要使用复杂接口**——不受支持
- **不要忘记 JSON 标签**——字段名称很重要
- **不要嵌套得太深**——保持简单

## 后续步骤

@cards{cols="2"}
🚀 方法绑定
了解如何绑定 Go 方法。

[了解更多 →](/features/bindings/methods/)

---
◆ 服务
使用服务组织代码。

[了解更多 →](/features/bindings/services/)

---
✓ 最佳实践
绑定设计模式。

[了解更多 →](/features/bindings/best-practices/)

---
★ Go-前端桥接
了解桥接机制。

[了解更多 →](/concepts/bridge/)

@end

---

<strong>有疑问？</strong>请在[Discord](https://discord.gg/JDdSxwjhGf)中提问，或查看[绑定示例](https://github.com/wailsapp/wails/tree/master/v3/examples/binding)。
