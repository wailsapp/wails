---
title: "데이터 모델"
description: "Go와 JavaScript 간에 복잡한 데이터 구조 바인딩"
slug: "features/bindings/models"
sourcePath: "features/bindings/models.md"
---

## 데이터 모델 바인딩

Wails는 Go 구조체에서 <strong>JavaScript/TypeScript 클래스를 자동으로 생성</strong>하여 백엔드와 프런트엔드 간에 복잡한 데이터를 전달할 때 완전한 타입 안전성을 제공합니다. Go 구조체를 작성하고 바인딩을 생성하면 생성자, 타입 주석, JSDoc 주석을 모두 갖춘 완전한 타입 지정 프런트엔드 모델을 얻을 수 있습니다.

## 빠른 시작

**Go 구조체:**

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

**생성:**

```bash
wails3 generate bindings
```

**JavaScript:**

```javascript
import { GetUser } from './bindings/changeme/userservice'
import { User } from './bindings/changeme/models'

const user = await GetUser(1)
console.log(user.Name)  // Type-safe!
```

**이것으로 끝입니다!** 브리지 전체에서 완전한 타입 안전성이 제공됩니다.

## 모델 정의

### 기본 구조체

```go
type Person struct {
    Name string
    Age  int
}
```

**생성된 JavaScript:**

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

### JSON 태그 사용

```go
type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"createdAt"`
}
```

**생성된 JavaScript:**

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

JavaScript의 **필드 이름은 JSON 태그로 제어합니다**.

### 사용자 정의 JSON 마샬링

바인딩된 메서드의 결과는 Go 표준 `encoding/json` 패키지로 마샬링됩니다. 포인터 리시버에 선언된 `MarshalJSON` 메서드는 마샬링되는 값이 주소 지정 가능할 때만 호출된다는 점에 유의하세요. 따라서 반환된 슬라이스의 요소에는 사용자 정의 마샬러가 호출될 수 있지만, 값으로 반환된 구조체에는 호출되지 않습니다:

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

일관된 출력을 보장하려면 타입을 안전하게 복사할 수 있는 경우 값 리시버에 `MarshalJSON`을 정의하거나 구조체의 포인터를 반환하세요(예: `*User`). 이는 Wails의 바인딩 생성 방식 차이가 아니라 Go의 `encoding/json` 동작입니다.

### 주석 사용

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

**생성된 JavaScript:**

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

**주석이 JSDoc으로 변환됩니다!** IDE에서 해당 주석을 확인할 수 있습니다.

### 중첩 구조체

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

**생성된 JavaScript:**

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

**사용법:**

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

### 배열과 슬라이스

```go
type Team struct {
    Name    string   `json:"name"`
    Members []string `json:"members"`
}
```

**생성된 JavaScript:**

```javascript
export class Team {
    /** @type {string} */
    name = ""
    
    /** @type {string[]} */
    members = []
}
```

**사용법:**

```javascript
const team = new Team({
    name: "Engineering",
    members: ["Alice", "Bob", "Charlie"]
})
```

### 맵

```go
type Config struct {
    Settings map[string]string `json:"settings"`
}
```

**생성된 JavaScript:**

```javascript
export class Config {
    /** @type {Record<string, string>} */
    settings = {}
}
```

**사용법:**

```javascript
const config = new Config({
    settings: {
        theme: "dark",
        language: "en"
    }
})
```

## 타입 매핑

### 기본 타입

| Go 타입 | JavaScript/TypeScript |
| --- | --- |
| `string` | `string` |
| `bool` | `boolean` |
| `int`, `int8`, `int16`, `int32`, `int64` | `number` |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | `number` |
| `float32`, `float64` | `number` |
| `byte` | `number` |
| `rune` | `number` |

### 복합 타입

| Go 타입 | JavaScript/TypeScript | 참고 |
| --- | --- | --- |
| `[]T` | `T[]` | - |
| `[N]T` | `T[]` | - |
| `map[string]T` | `{ [_: string]: T }` | 문자열 키 맵 |
| `map[K]V` | `{ [_ in K]?: V }` | 문자열이 아닌 `K`은 매핑된 타입으로 렌더링되며, JS `Map`이 **아닙니다** |
| `[]byte` | `string` | base64로 인코딩됨 |
| `struct` | `class` / `interface` | 필드 포함 |
| `time.Time` | `any` | 런타임에서 RFC3339Nano 문자열로 직렬화됨 |
| `*T` | `T \| null` | 포인터는 null 허용을 의미함 |
| `any` / `interface{}` | `any` | - |
| `error` | `any` / `Exception` | 반환 값이면 Exception, 그 외에는 any |

### 지원되지 않는 타입

- `chan T`(채널)
- `func()`(함수)
- 복잡한 인터페이스(`interface{}` 제외)
- 내보내지 않은 필드(소문자로 시작)

## 모델 사용하기

### 인스턴스 생성하기

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

### Go로 전달하기

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

### Go에서 받기

```javascript
import { GetUser } from './bindings/changeme/userservice'

const user = await GetUser(1)

// user is already a User instance
console.log(user.name)
console.log(user.email)
console.log(user.createdAt.toISOString())
```

### 모델 업데이트하기

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

## TypeScript 지원

### 생성된 TypeScript

```bash
wails3 generate bindings -ts
```

**생성 결과:**

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

### TypeScript에서 사용하기

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

## 고급 패턴

### 선택적 필드

```go
type User struct {
    ID       int     `json:"id"`
    Name     string  `json:"name"`
    Nickname *string `json:"nickname,omitempty"`
}
```

**JavaScript:**

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

### 열거형

바인딩 생성기는 상수가 있는 Go 명명 타입을 자동으로 감지하고 TypeScript 열거형 또는 JavaScript const 객체를 생성합니다. 여기에는 Go의 제로 값에 해당하는 `$zero` 멤버와 완전히 보존된 JSDoc도 포함됩니다.

```go
type UserRole string

const (
    RoleAdmin UserRole = "admin"
    RoleUser  UserRole = "user"
    RoleGuest UserRole = "guest"
)
```

**생성된 TypeScript:**

```typescript
export enum UserRole {
    $zero = "",
    RoleAdmin = "admin",
    RoleUser = "user",
    RoleGuest = "guest",
}
```

**사용법:**

```javascript
import { User, UserRole } from './bindings/changeme/models'

const admin = new User({
    name: "Admin",
    role: UserRole.RoleAdmin
})
```

문자열 열거형, 정수 열거형, 타입 별칭, 가져온 패키지의 열거형 및 제한 사항을 포괄적으로 알아보려면 전용 **[열거형](/features/bindings/enums/)** 페이지를 참조하세요.

### 유효성 검사

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

### 직렬화

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

## 전체 예제

**Go:**

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

**JavaScript:**

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

## 모범 사례

### ✅ 권장 사항

- **JSON 태그를 사용하세요** - 필드 이름을 제어할 수 있습니다
- **주석을 추가하세요** - JSDoc으로 변환됩니다
- **time.Time을 사용하세요** - Date로 변환됩니다
- **Go 측에서 유효성을 검사하세요** - 프런트엔드를 신뢰하지 마세요
- **모델을 단순하게 유지하세요** - 데이터 컨테이너로만 사용하세요
- **선택적 필드에는 포인터를 사용하세요** - null을 허용하려면 `*string`을 사용하세요

### ❌ 금지 사항

- **Go 구조체에 메서드를 추가하지 마세요** - 데이터로만 유지하세요
- **내보내지 않은 필드를 사용하지 마세요** - 바인딩되지 않습니다
- **복잡한 인터페이스를 사용하지 마세요** - 지원되지 않습니다
- **JSON 태그를 빠뜨리지 마세요** - 필드 이름은 중요합니다
- **너무 깊게 중첩하지 마세요** - 단순하게 유지하세요

## 다음 단계

@cards{cols="2"}
🚀 메서드 바인딩
Go 메서드를 바인딩하는 방법을 알아보세요.

[자세히 알아보기 →](/features/bindings/methods/)

---
◆ 서비스
서비스로 코드를 구성하세요.

[자세히 알아보기 →](/features/bindings/services/)

---
✓ 모범 사례
바인딩 디자인 패턴을 알아보세요.

[자세히 알아보기 →](/features/bindings/best-practices/)

---
★ Go-프런트엔드 브리지
브리지 메커니즘을 이해하세요.

[자세히 알아보기 →](/concepts/bridge/)

@end

---

**궁금한 점이 있으신가요?** [Discord](https://discord.gg/JDdSxwjhGf)에서 질문하거나 [바인딩 예제](https://github.com/wailsapp/wails/tree/master/v3/examples/binding)를 확인하세요.
