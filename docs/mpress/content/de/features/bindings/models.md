---
title: "Datenmodelle"
description: "Komplexe Datenstrukturen zwischen Go und JavaScript binden"
slug: "features/bindings/models"
sourcePath: "features/bindings/models.md"
---

## Datenmodellbindungen

Wails **generiert automatisch JavaScript-/TypeScript-Klassen** aus Go-Structs und gewährleistet vollständige Typsicherheit bei der Übergabe komplexer Daten zwischen Backend und Frontend. Schreiben Sie Go-Structs, generieren Sie Bindungen und erhalten Sie vollständig typisierte Frontend-Modelle einschließlich Konstruktoren, Typannotationen und JSDoc-Kommentaren.

## Schnellstart

**Go-Struct:**

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

**Generieren:**

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

**Das ist alles!** Vollständige Typsicherheit über die Bridge hinweg.

## Modelle definieren

### Einfacher Struct

```go
type Person struct {
    Name string
    Age  int
}
```

**Generiertes JavaScript:**

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

### Mit JSON-Tags

```go
type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"createdAt"`
}
```

**Generiertes JavaScript:**

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

**JSON-Tags bestimmen die Feldnamen** in JavaScript.

### Benutzerdefiniertes JSON-Marshalling

Ergebnisse gebundener Methoden werden mit dem Go-Standardpaket `encoding/json` gemarshallt. Beachten Sie, dass für einen Pointer-Receiver deklarierte `MarshalJSON`-Methoden nur aufgerufen werden, wenn der zu marshallende Wert adressierbar ist. Das bedeutet, dass ein benutzerdefinierter Marshaler für Elemente eines zurückgegebenen Slices aufgerufen werden kann, nicht jedoch für einen als Wert zurückgegebenen Struct:

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

Um eine konsistente Ausgabe sicherzustellen, definieren Sie `MarshalJSON` entweder für einen Value-Receiver, wenn der Typ gefahrlos kopiert werden kann, oder geben Sie einen Pointer auf den Struct zurück (zum Beispiel `*User`). Dies ist ein Verhalten von Gos `encoding/json` und kein Unterschied bei der Bindungsgenerierung von Wails.

### Mit Kommentaren

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

**Generiertes JavaScript:**

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

**Kommentare werden zu JSDoc!** Ihre IDE zeigt sie an.

### Verschachtelte Structs

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

**Generiertes JavaScript:**

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

**Verwendung:**

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

### Arrays und Slices

```go
type Team struct {
    Name    string   `json:"name"`
    Members []string `json:"members"`
}
```

**Generiertes JavaScript:**

```javascript
export class Team {
    /** @type {string} */
    name = ""
    
    /** @type {string[]} */
    members = []
}
```

**Verwendung:**

```javascript
const team = new Team({
    name: "Engineering",
    members: ["Alice", "Bob", "Charlie"]
})
```

### Maps

```go
type Config struct {
    Settings map[string]string `json:"settings"`
}
```

**Generiertes JavaScript:**

```javascript
export class Config {
    /** @type {Record<string, string>} */
    settings = {}
}
```

**Verwendung:**

```javascript
const config = new Config({
    settings: {
        theme: "dark",
        language: "en"
    }
})
```

## Typzuordnung

### Primitive Typen

| Go-Typ | JavaScript/TypeScript |
| --- | --- |
| `string` | `string` |
| `bool` | `boolean` |
| `int`, `int8`, `int16`, `int32`, `int64` | `number` |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | `number` |
| `float32`, `float64` | `number` |
| `byte` | `number` |
| `rune` | `number` |

### Komplexe Typen

| Go-Typ | JavaScript/TypeScript | Hinweise |
| --- | --- | --- |
| `[]T` | `T[]` | - |
| `[N]T` | `T[]` | - |
| `map[string]T` | `{ [_: string]: T }` | Map mit String-Schlüsseln |
| `map[K]V` | `{ [_ in K]?: V }` | `K`, die keine Strings sind, werden als zugeordneter Typ dargestellt, **nicht** als JavaScript-`Map` |
| `[]byte` | `string` | Base64-codiert |
| `struct` | `class` / `interface` | mit Feldern |
| `time.Time` | `any` | wird zur Laufzeit als RFC3339Nano-Zeichenfolge serialisiert |
| `*T` | `T \| null` | Zeiger bedeutet, dass null zulässig ist |
| `any` / `interface{}` | `any` | - |
| `error` | `any` / `Exception` | Als Rückgabewert Exception, andernfalls any |

### Nicht unterstützte Typen

- `chan T` (Kanäle)
- `func()` (Funktionen)
- Komplexe Interfaces (außer `interface{}`)
- Nicht exportierte Felder (kleingeschrieben)

## Modelle verwenden

### Instanzen erstellen

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

### An Go übergeben

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

### Von Go empfangen

```javascript
import { GetUser } from './bindings/changeme/userservice'

const user = await GetUser(1)

// user is already a User instance
console.log(user.name)
console.log(user.email)
console.log(user.createdAt.toISOString())
```

### Modelle aktualisieren

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

## TypeScript-Unterstützung

### Generiertes TypeScript

```bash
wails3 generate bindings -ts
```

**Generiert:**

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

### Verwendung in TypeScript

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

## Fortgeschrittene Muster

### Optionale Felder

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

### Enumerationen

Der Binding-Generator erkennt automatisch benannte Go-Typen mit Konstanten und generiert TypeScript-Enumerationen oder JavaScript-const-Objekte – einschließlich eines `$zero`-Members für den Nullwert von Go und unter vollständiger Beibehaltung von JSDoc.

```go
type UserRole string

const (
    RoleAdmin UserRole = "admin"
    RoleUser  UserRole = "user"
    RoleGuest UserRole = "guest"
)
```

**Generiertes TypeScript:**

```typescript
export enum UserRole {
    $zero = "",
    RoleAdmin = "admin",
    RoleUser = "user",
    RoleGuest = "guest",
}
```

**Verwendung:**

```javascript
import { User, UserRole } from './bindings/changeme/models'

const admin = new User({
    name: "Admin",
    role: UserRole.RoleAdmin
})
```

Eine umfassende Beschreibung von Zeichenfolgen- und Ganzzahl-Enumerationen, Typaliasen, Enumerationen aus importierten Paketen sowie Einschränkungen finden Sie auf der gesonderten Seite **[Enumerationen](/features/bindings/enums/)**.

### Validierung

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

### Serialisierung

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

## Vollständiges Beispiel

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

## Bewährte Vorgehensweisen

### ✅ Empfohlen

- **JSON-Tags verwenden** – Feldnamen steuern
- **Kommentare hinzufügen** – Sie werden zu JSDoc
- **time.Time verwenden** – Wird in Date konvertiert
- **Auf der Go-Seite validieren** – Dem Frontend nicht vertrauen
- **Modelle einfach halten** – Nur als Datencontainer verwenden
- **Zeiger für optionale Felder verwenden** – `*string` verwenden, wenn null zulässig sein soll

### ❌ Nicht tun

- **Go-Structs keine Methoden hinzufügen** – Als reine Datenstrukturen belassen
- **Keine nicht exportierten Felder verwenden** – Für sie werden keine Bindings erstellt
- **Keine komplexen Interfaces verwenden** – Sie werden nicht unterstützt
- **JSON-Tags nicht vergessen** – Feldnamen sind wichtig
- **Nicht zu tief verschachteln** – Einfach halten

## Nächste Schritte

@cards{cols="2"}
🚀 Methoden-Bindings
Erfahren Sie, wie Sie Go-Methoden binden.

[Mehr erfahren →](/features/bindings/methods/)

---
◆ Dienste
Organisieren Sie Code mithilfe von Diensten.

[Mehr erfahren →](/features/bindings/services/)

---
✓ Bewährte Vorgehensweisen
Entwurfsmuster für Bindings.

[Mehr erfahren →](/features/bindings/best-practices/)

---
★ Go-Frontend-Brücke
Lernen Sie den Brückenmechanismus kennen.

[Mehr erfahren →](/concepts/bridge/)

@end

---

**Fragen?** Stellen Sie sie auf [Discord](https://discord.gg/JDdSxwjhGf) oder sehen Sie sich die [Binding-Beispiele](https://github.com/wailsapp/wails/tree/master/v3/examples/binding) an.
