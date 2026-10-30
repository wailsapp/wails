---
title: "データモデル"
description: "Go と JavaScript の間で複雑なデータ構造をバインドする"
slug: "features/bindings/models"
sourcePath: "features/bindings/models.md"
---

## データモデルのバインディング

Wails は、Go の構造体から**JavaScript/TypeScript クラスを自動生成し**、バックエンドとフロントエンドの間で複雑なデータを受け渡す際の完全な型安全性を提供します。Go の構造体を記述してバインディングを生成するだけで、コンストラクター、型注釈、JSDoc コメントを備えた、完全に型付けされたフロントエンドモデルを利用できます。

## クイックスタート

**Go の構造体：**

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
console.log(user.name)  // Type-safe!
```

<strong>これだけです。</strong>ブリッジ全体で完全な型安全性が得られます。

## モデルの定義

以下のモデルクラスはフィールドの構造を示すためのものです。これらの簡略化したクラスを変換コードとしてコピーするのではなく、実際に生成されたバインディングを使用してください。

`time.Time`の場合、ジェネレーターのデフォルトは`-time-type=string`です。タイムスタンプをJavaScriptの`Date`値に変換する生成クラスには、`-time-type=Date`を使用してください。標準のJavaScript／クラステンプレートのビルドタスクでは`Date`が選択され、TypeScriptインターフェースのタスクでは`string`が維持されます。モードを変更する前に、生成されたタスクを確認してください。Goのタイムスタンプは、JSON内では引き続きRFC3339Nano文字列です。

Dateモードでは、生成された変換ヘルパーがこれらの文字列を日付に変換します。`Object.assign`だけでは、この変換は行われません。文字列バインディングでは、Dateのメソッドを使用する前に、`new Date(timestamp)`でタイムスタンプを明示的に変換してください。

### 基本的な構造体

```go
type Person struct {
    Name string
    Age  int
}
```

**JavaScriptの例：**

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

### JSON タグを使用する

```go
type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"createdAt"`
}
```

**JavaScriptの例：**

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

JavaScript での**フィールド名は JSON タグによって制御されます**。

### カスタム JSON マーシャリング

バインドされたメソッドの結果は、Go の標準`encoding/json`パッケージでマーシャリングされます。ポインターレシーバーに宣言された`MarshalJSON`メソッドは、マーシャリング対象の値がアドレス指定可能な場合にのみ呼び出されることに注意してください。つまり、返されたスライスの要素ではカスタムマーシャラーが呼び出されますが、値として返された構造体では呼び出されない場合があります。

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

一貫した出力を確保するには、その型を安全にコピーできる場合は値レシーバーに`MarshalJSON`を定義するか、構造体へのポインター（たとえば`*User`）を返してください。これは Wails のバインディング生成の違いではなく、Go の`encoding/json`の動作です。

### コメントを使用する

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

**JavaScriptの例：**

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

<strong>コメントは JSDoc になります。</strong>IDE にコメントが表示されます。

### ネストされた構造体

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

**JavaScriptの例：**

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

**使用例：**

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

### 配列とスライス

```go
type Team struct {
    Name    string   `json:"name"`
    Members []string `json:"members"`
}
```

**JavaScriptの例：**

```javascript
export class Team {
    /** @type {string} */
    name = ""
    
    /** @type {string[]} */
    members = []
}
```

**使用例：**

```javascript
const team = new Team({
    name: "Engineering",
    members: ["Alice", "Bob", "Charlie"]
})
```

### マップ

```go
type Config struct {
    Settings map[string]string `json:"settings"`
}
```

**JavaScriptの例：**

```javascript
export class Config {
    /** @type {Record<string, string>} */
    settings = {}
}
```

**使用例：**

```javascript
const config = new Config({
    settings: {
        theme: "dark",
        language: "en"
    }
})
```

## 型のマッピング

### プリミティブ型

| Go の型 | JavaScript/TypeScript |
| --- | --- |
| `string` | `string` |
| `bool` | `boolean` |
| `int`, `int8`, `int16`, `int32`, `int64` | `number` |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | `number` |
| `float32`, `float64` | `number` |
| `byte` | `number` |
| `rune` | `number` |

### 複合型

| Go の型 | JavaScript/TypeScript | 注記 |
| --- | --- | --- |
| `[]T` | `T[]` | - |
| `[N]T` | `T[]` | - |
| `map[string]T` | `{ [_: string]: T }` | 文字列をキーとするマップ |
| `map[K]V` | `{ [_ in K]?: V }` | 文字列以外の`K`は、JS の`Map`では**なく**、マップ型としてレンダリングされる |
| `[]byte` | `string` | base64 エンコード済み |
| `struct` | `class` / `interface` | フィールドを持つ |
| `time.Time` | `string`（デフォルト）または`Date` | `-time-type`で生成される表現を選択する。JSON では RFC3339Nano 文字列を使用する |
| `*T` | `T \| null` | ポインターは null 許容を意味する |
| `any` / `interface{}` | `any` | - |
| `error` | `any` / `Exception` | 戻り値の場合は例外、それ以外は any |

### サポートされていない型

- `chan T`（チャネル）
- `func()`（関数）
- 複雑なインターフェース（`interface{}`を除く）
- 非公開フィールド（小文字で始まるもの）

## モデルの使用

### インスタンスの作成

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

### Go への受け渡し

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

### Go からの受け取り

```javascript
import { GetUser } from './bindings/changeme/userservice'

const user = await GetUser(1)

// user is already a User instance
console.log(user.name)
console.log(user.email)
console.log(user.createdAt.toISOString())
```

### モデルの更新

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

## TypeScript のサポート

### 生成される TypeScript

このセクションでは、Date 変換を伴うクラスバインディングについて説明します。通常の`-ts`コマンドでは、ジェネレーターのデフォルトである文字列表現が維持されます。表示しているクラスは簡略化されています。ジェネレーターは変換ヘルパーも生成します。

```bash
wails3 generate bindings -ts -time-type=Date
```

**生成結果：**

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

### TypeScript での使用

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

## 高度なパターン

### オプションフィールド

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

### 列挙型

バインディングジェネレーターは、定数を持つ Go の名前付き型を自動検出し、TypeScript の列挙型または JavaScript の const オブジェクトを生成します。これには、Go のゼロ値に対応する`$zero`メンバーと、完全に保持された JSDoc も含まれます。

```go
type UserRole string

const (
    RoleAdmin UserRole = "admin"
    RoleUser  UserRole = "user"
    RoleGuest UserRole = "guest"
)
```

**生成された TypeScript：**

```typescript
export enum UserRole {
    $zero = "",
    RoleAdmin = "admin",
    RoleUser = "user",
    RoleGuest = "guest",
}
```

**使用例：**

```javascript
import { User, UserRole } from './bindings/changeme/models'

const admin = new User({
    name: "Admin",
    role: UserRole.RoleAdmin
})
```

文字列列挙型、整数列挙型、型エイリアス、インポートしたパッケージの列挙型、および制限事項についての包括的な説明は、専用の<strong>[列挙型](/features/bindings/enums/)</strong>ページを参照してください。

### 検証

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

### シリアライズ

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

## 完全な例

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

## ベストプラクティス

### ✅ 推奨事項

- **JSON タグを使用する** - フィールド名を制御できます
- **コメントを追加する** - JSDoc に変換されます
- **時刻の表現を選択する** - `time.Time`はデフォルトで`string`になります。`-time-type=Date`を使用すると、生成される Date 変換が有効になります
- **Go 側で検証する** - フロントエンドを信用しないでください
- **モデルをシンプルに保つ** - データコンテナーとしてのみ使用してください
- **オプション項目にはポインターを使用する** - null 許容には`*string`を使用します

### ❌ 非推奨事項

- **Go 構造体にメソッドを追加しない** - データとして維持してください
- **非公開フィールドを使用しない** - バインドされません
- **複雑なインターフェースを使用しない** - サポートされていません
- **JSON タグを付け忘れない** - フィールド名は重要です
- **ネストを深くしすぎない** - シンプルに保ってください

## 次のステップ

@cards{cols="2"}
🚀 メソッドバインディング
Go メソッドをバインドする方法を学びます。

[詳細を見る →](/features/bindings/methods/)

---
◆ サービス
サービスを使用してコードを整理します。

[詳細を見る →](/features/bindings/services/)

---
✓ ベストプラクティス
バインディングの設計パターン。

[詳細を見る →](/features/bindings/best-practices/)

---
★ Go とフロントエンドのブリッジ
ブリッジの仕組みを理解します。

[詳細を見る →](/concepts/bridge/)

@end

---

**質問がありますか？** [Discord](https://discord.gg/JDdSxwjhGf)で質問するか、[バインディングの例](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples/binding)を確認してください。
