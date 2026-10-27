---
title: "Модели данных"
description: "Связывайте сложные структуры данных между Go и JavaScript"
slug: "features/bindings/models"
sourcePath: "features/bindings/models.md"
---

## Связывание моделей данных

Wails **автоматически создаёт классы JavaScript/TypeScript** из структур Go, обеспечивая полную типобезопасность при передаче сложных данных между бэкендом и фронтендом. Определите структуры Go, сгенерируйте привязки и получите полностью типизированные модели для фронтенда с конструкторами, аннотациями типов и комментариями JSDoc.

## Быстрый старт

**Структура Go:**

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

**Сгенерируйте:**

```bash
wails3 generate bindings
```

**JavaScript:**

```javascript
import { GetUser } from './bindings/changeme/userservice'
import { User } from './bindings/changeme/models'

const user = await GetUser(1)
console.log(user.name)  // Type-safe!
```

**Готово!** Полная типобезопасность при передаче данных через мост.

## Определение моделей

Приведённые ниже классы моделей иллюстрируют формы полей; используйте фактически сгенерированные привязки, а не копируйте эти упрощённые классы в качестве кода преобразования.

Для `time.Time` генератор по умолчанию использует `-time-type=string`. Используйте `-time-type=Date` для сгенерированных классов, преобразующих метки времени в значения JavaScript `Date`. Стандартные задачи сборки шаблонов JavaScript и классов выбирают `Date`, а задачи для интерфейсов TypeScript сохраняют `string`. Перед сменой режима проверьте сгенерированную задачу. В JSON метки времени Go остаются строками в формате RFC3339Nano.

В режиме Date сгенерированные вспомогательные функции преобразуют эти строки в даты. Сам по себе `Object.assign` такого преобразования не выполняет. При использовании строковых привязок явно преобразуйте метку времени с помощью `new Date(timestamp)`, прежде чем применять методы Date.

### Базовая структура

```go
type Person struct {
    Name string
    Age  int
}
```

**Пример на JavaScript:**

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

### С тегами JSON

```go
type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"createdAt"`
}
```

**Пример на JavaScript:**

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

**Теги JSON определяют имена полей** в JavaScript.

### Пользовательская сериализация JSON

Результаты связанных методов сериализуются стандартным пакетом Go `encoding/json`. Учтите, что методы `MarshalJSON`, объявленные для получателя-указателя, вызываются только тогда, когда сериализуемое значение является адресуемым. Поэтому пользовательский маршалер может вызываться для элементов возвращаемого среза, но не для структуры, возвращаемой по значению:

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

Чтобы обеспечить единообразный результат, либо определите `MarshalJSON` для получателя-значения, если тип можно безопасно копировать, либо возвращайте указатель на структуру (например, `*User`). Это особенность `encoding/json` в Go, а не отличие в генерации привязок Wails.

### С комментариями

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

**Пример на JavaScript:**

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

**Комментарии преобразуются в JSDoc!** Они отображаются в вашей IDE.

### Вложенные структуры

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

**Пример на JavaScript:**

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

**Использование:**

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

### Массивы и срезы

```go
type Team struct {
    Name    string   `json:"name"`
    Members []string `json:"members"`
}
```

**Пример на JavaScript:**

```javascript
export class Team {
    /** @type {string} */
    name = ""
    
    /** @type {string[]} */
    members = []
}
```

**Использование:**

```javascript
const team = new Team({
    name: "Engineering",
    members: ["Alice", "Bob", "Charlie"]
})
```

### Отображения

```go
type Config struct {
    Settings map[string]string `json:"settings"`
}
```

**Пример на JavaScript:**

```javascript
export class Config {
    /** @type {Record<string, string>} */
    settings = {}
}
```

**Использование:**

```javascript
const config = new Config({
    settings: {
        theme: "dark",
        language: "en"
    }
})
```

## Сопоставление типов

### Примитивные типы

| Тип Go | JavaScript/TypeScript |
| --- | --- |
| `string` | `string` |
| `bool` | `boolean` |
| `int`, `int8`, `int16`, `int32`, `int64` | `number` |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | `number` |
| `float32`, `float64` | `number` |
| `byte` | `number` |
| `rune` | `number` |

### Сложные типы

| Тип Go | JavaScript/TypeScript | Примечания |
| --- | --- | --- |
| `[]T` | `T[]` | - |
| `[N]T` | `T[]` | - |
| `map[string]T` | `{ [_: string]: T }` | отображение со строковыми ключами |
| `map[K]V` | `{ [_ in K]?: V }` | `K` с ключами нестрокового типа отображается как отображаемый тип, а **не** как `Map` в JS |
| `[]byte` | `string` | в кодировке base64 |
| `struct` | `class` / `interface` | с полями |
| `time.Time` | `string` (по умолчанию) или `Date` | `-time-type` задаёт генерируемое представление; в JSON используется строка в формате RFC3339Nano |
| `*T` | `T \| null` | указатель означает, что допустимо значение null |
| `any` / `interface{}` | `any` | - |
| `error` | `any` / `Exception` | Исключение при использовании в качестве возвращаемого значения, иначе — any |

### Неподдерживаемые типы

- `chan T` (каналы)
- `func()` (функции)
- Сложные интерфейсы (кроме `interface{}`)
- Неэкспортируемые поля (со строчной буквы)

## Использование моделей

### Создание экземпляров

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

### Передача в Go

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

### Получение из Go

```javascript
import { GetUser } from './bindings/changeme/userservice'

const user = await GetUser(1)

// user is already a User instance
console.log(user.name)
console.log(user.email)
console.log(user.createdAt.toISOString())
```

### Обновление моделей

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

## Поддержка TypeScript

### Сгенерированный код TypeScript

В этом разделе показаны привязки классов с преобразованием Date. Обычная команда `-ts` сохраняет используемое генератором по умолчанию строковое представление. Показанный класс упрощён; генератор также создаёт вспомогательные функции преобразования.

```bash
wails3 generate bindings -ts -time-type=Date
```

**Сгенерированный код:**

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

### Использование в TypeScript

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

## Расширенные шаблоны

### Необязательные поля

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

### Перечисления

Генератор привязок автоматически обнаруживает именованные типы Go с константами и создаёт перечисления TypeScript или константные объекты JavaScript, включая элемент `$zero` для нулевого значения Go и полностью сохраняя JSDoc.

```go
type UserRole string

const (
    RoleAdmin UserRole = "admin"
    RoleUser  UserRole = "user"
    RoleGuest UserRole = "guest"
)
```

**Сгенерированный код TypeScript:**

```typescript
export enum UserRole {
    $zero = "",
    RoleAdmin = "admin",
    RoleUser = "user",
    RoleGuest = "guest",
}
```

**Использование:**

```javascript
import { User, UserRole } from './bindings/changeme/models'

const admin = new User({
    name: "Admin",
    role: UserRole.RoleAdmin
})
```

Подробные сведения о строковых и целочисленных перечислениях, псевдонимах типов, перечислениях из импортированных пакетов и ограничениях см. на отдельной странице **[«Перечисления»](/features/bindings/enums/)**.

### Валидация

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

### Сериализация

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

## Полный пример

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
                <small>Created: ${new Date(user.createdAt).toLocaleDateString()}</small>
            </div>
        `).join('')
    }
}

const manager = new UserManager()
manager.loadUsers()
```

## Рекомендации

### ✅ Рекомендуется

- **Используйте теги JSON** — они позволяют управлять именами полей
- **Добавляйте комментарии** — они преобразуются в JSDoc
- **Выберите представление времени** — по умолчанию `time.Time` преобразуется в `string`; `-time-type=Date` включает генерируемое преобразование Date
- **Проверяйте данные на стороне Go** — не доверяйте фронтенду
- **Сохраняйте модели простыми** — используйте их только как контейнеры данных
- **Используйте указатели для необязательных значений** — `*string` для значений, допускающих null

### ❌ Не рекомендуется

- **Не добавляйте методы в структуры Go** — оставляйте их структурами данных
- **Не используйте неэкспортируемые поля** — для них не будут созданы привязки
- **Не используйте сложные интерфейсы** — они не поддерживаются
- **Не забывайте о тегах JSON** — имена полей имеют значение
- **Избегайте чрезмерной вложенности** — сохраняйте структуру простой

## Дальнейшие шаги

@cards{cols="2"}
🚀 Привязки методов
Узнайте, как привязывать методы Go.

[Подробнее →](/features/bindings/methods/)

---
◆ Сервисы
Организуйте код с помощью сервисов.

[Подробнее →](/features/bindings/services/)

---
✓ Рекомендации
Шаблоны проектирования привязок.

[Подробнее →](/features/bindings/best-practices/)

---
★ Мост между Go и фронтендом
Разберитесь в механизме работы моста.

[Подробнее →](/concepts/bridge/)

@end

---

**Есть вопросы?** Задайте их в [Discord](https://discord.gg/JDdSxwjhGf) или ознакомьтесь с [примерами привязок](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples/binding).
