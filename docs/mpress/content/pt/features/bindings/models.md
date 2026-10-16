---
title: "Modelos de dados"
description: "Vincule estruturas de dados complexas entre Go e JavaScript"
slug: "features/bindings/models"
sourcePath: "features/bindings/models.md"
---

## Vinculações de modelos de dados

O Wails **gera automaticamente classes JavaScript/TypeScript** a partir de structs do Go, garantindo total segurança de tipos ao passar dados complexos entre o backend e o frontend. Escreva structs em Go, gere as vinculações e obtenha modelos de frontend totalmente tipados, com construtores, anotações de tipo e comentários JSDoc.

## Início rápido

**Struct do Go:**

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

**Gere:**

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

**É só isso!** Segurança total de tipos em toda a ponte.

## Definição de modelos

### Struct básica

```go
type Person struct {
    Name string
    Age  int
}
```

**JavaScript gerado:**

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

### Com tags JSON

```go
type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"createdAt"`
}
```

**JavaScript gerado:**

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

**As tags JSON controlam os nomes dos campos** no JavaScript.

### Serialização JSON personalizada

Os resultados de métodos vinculados são serializados com o pacote padrão `encoding/json` do Go. Esteja ciente de que os métodos `MarshalJSON` declarados em um receptor de ponteiro só são chamados quando o valor que está sendo serializado é endereçável. Isso significa que um serializador personalizado pode ser chamado para elementos de um slice retornado, mas não para uma struct retornada por valor:

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

Para garantir uma saída consistente, defina `MarshalJSON` em um receptor de valor quando o tipo puder ser copiado com segurança ou retorne um ponteiro para a struct (por exemplo, `*User`). Esse é um comportamento de `encoding/json` do Go, e não uma diferença na geração de vinculações do Wails.

### Com comentários

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

**JavaScript gerado:**

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

**Os comentários se tornam JSDoc!** Sua IDE os exibe.

### Structs aninhadas

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

**JavaScript gerado:**

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

**Uso:**

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

### Arrays e slices

```go
type Team struct {
    Name    string   `json:"name"`
    Members []string `json:"members"`
}
```

**JavaScript gerado:**

```javascript
export class Team {
    /** @type {string} */
    name = ""
    
    /** @type {string[]} */
    members = []
}
```

**Uso:**

```javascript
const team = new Team({
    name: "Engineering",
    members: ["Alice", "Bob", "Charlie"]
})
```

### Mapas

```go
type Config struct {
    Settings map[string]string `json:"settings"`
}
```

**JavaScript gerado:**

```javascript
export class Config {
    /** @type {Record<string, string>} */
    settings = {}
}
```

**Uso:**

```javascript
const config = new Config({
    settings: {
        theme: "dark",
        language: "en"
    }
})
```

## Mapeamento de tipos

### Tipos primitivos

| Tipo Go | JavaScript/TypeScript |
| --- | --- |
| `string` | `string` |
| `bool` | `boolean` |
| `int`, `int8`, `int16`, `int32`, `int64` | `number` |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | `number` |
| `float32`, `float64` | `number` |
| `byte` | `number` |
| `rune` | `number` |

### Tipos complexos

| Tipo Go | JavaScript/TypeScript | Observações |
| --- | --- | --- |
| `[]T` | `T[]` | - |
| `[N]T` | `T[]` | - |
| `map[string]T` | `{ [_: string]: T }` | mapa com chaves do tipo string |
| `map[K]V` | `{ [_ in K]?: V }` | `K` que não seja string é renderizado como um tipo mapeado, e **não** como um `Map` do JS |
| `[]byte` | `string` | codificado em base64 |
| `struct` | `class` / `interface` | com campos |
| `time.Time` | `any` | serializado como string RFC3339Nano no runtime |
| `*T` | `T \| null` | ponteiro indica que o valor pode ser nulo |
| `any` / `interface{}` | `any` | - |
| `error` | `any` / `Exception` | Exception se for usado como valor de retorno; caso contrário, any |

### Tipos não compatíveis

- `chan T` (canais)
- `func()` (funções)
- Interfaces complexas (exceto `interface{}`)
- Campos não exportados (iniciados por letra minúscula)

## Como usar modelos

### Como criar instâncias

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

### Como passar dados para Go

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

### Como receber dados de Go

```javascript
import { GetUser } from './bindings/changeme/userservice'

const user = await GetUser(1)

// user is already a User instance
console.log(user.name)
console.log(user.email)
console.log(user.createdAt.toISOString())
```

### Como atualizar modelos

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

## Suporte a TypeScript

### TypeScript gerado

```bash
wails3 generate bindings -ts
```

**Gerado:**

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

### Uso no TypeScript

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

## Padrões avançados

### Campos opcionais

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

### Enums

O gerador de bindings detecta automaticamente tipos nomeados do Go com constantes e gera enums do TypeScript ou objetos const do JavaScript — incluindo um membro `$zero` para o valor zero do Go e preservando integralmente o JSDoc.

```go
type UserRole string

const (
    RoleAdmin UserRole = "admin"
    RoleUser  UserRole = "user"
    RoleGuest UserRole = "guest"
)
```

**TypeScript gerado:**

```typescript
export enum UserRole {
    $zero = "",
    RoleAdmin = "admin",
    RoleUser = "user",
    RoleGuest = "guest",
}
```

**Uso:**

```javascript
import { User, UserRole } from './bindings/changeme/models'

const admin = new User({
    name: "Admin",
    role: UserRole.RoleAdmin
})
```

Para uma abordagem completa de enums de string, enums de inteiros, aliases de tipo, enums de pacotes importados e limitações, consulte a página específica sobre **[Enums](/features/bindings/enums/)**.

### Validação

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

### Serialização

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

## Exemplo completo

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

## Práticas recomendadas

### ✅ Faça

- **Use tags JSON** — Controle os nomes dos campos
- **Adicione comentários** — Eles são convertidos em JSDoc
- **Use time.Time** — É convertido em Date
- **Valide no lado do Go** — Não confie no frontend
- **Mantenha os modelos simples** — Use-os apenas como contêineres de dados
- **Use ponteiros para campos opcionais** — `*string` para valores que podem ser nulos

### ❌ Não faça

- **Não adicione métodos às structs do Go** — Mantenha-as como dados
- **Não use campos não exportados** — Eles não serão incluídos no binding
- **Não use interfaces complexas** — Elas não são compatíveis
- **Não se esqueça das tags JSON** — Os nomes dos campos são importantes
- **Não crie níveis de aninhamento demais** — Mantenha a simplicidade

## Próximas etapas

@cards{cols="2"}
🚀 Bindings de métodos
Saiba como criar bindings para métodos do Go.

[Saiba mais →](/features/bindings/methods/)

---
◆ Serviços
Organize o código com serviços.

[Saiba mais →](/features/bindings/services/)

---
✓ Práticas recomendadas
Padrões de projeto para bindings.

[Saiba mais →](/features/bindings/best-practices/)

---
★ Ponte entre Go e frontend
Entenda o mecanismo da ponte.

[Saiba mais →](/concepts/bridge/)

@end

---

**Dúvidas?** Pergunte no [Discord](https://discord.gg/JDdSxwjhGf) ou consulte os [exemplos de bindings](https://github.com/wailsapp/wails/tree/master/v3/examples/binding).
