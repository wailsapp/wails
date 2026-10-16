---
title: "Modèles de données"
description: "Liez des structures de données complexes entre Go et JavaScript"
slug: "features/bindings/models"
sourcePath: "features/bindings/models.md"
---

## Liaisons de modèles de données

Wails **génère automatiquement des classes JavaScript/TypeScript** à partir des structures Go, assurant un typage entièrement sûr lors du transfert de données complexes entre le backend et le frontend. Écrivez des structures Go, générez les liaisons et obtenez des modèles frontend entièrement typés, avec leurs constructeurs, annotations de type et commentaires JSDoc.

## Démarrage rapide

**Structure Go :**

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

**Générez :**

```bash
wails3 generate bindings
```

**JavaScript :**

```javascript
import { GetUser } from './bindings/changeme/userservice'
import { User } from './bindings/changeme/models'

const user = await GetUser(1)
console.log(user.Name)  // Type-safe!
```

**C’est tout !** Vous bénéficiez d’un typage entièrement sûr de part et d’autre du pont.

## Définition des modèles

### Structure de base

```go
type Person struct {
    Name string
    Age  int
}
```

**JavaScript généré :**

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

### Avec des balises JSON

```go
type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"createdAt"`
}
```

**JavaScript généré :**

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

Dans JavaScript, **les balises JSON déterminent le nom des champs**.

### Sérialisation JSON personnalisée

Les résultats des méthodes liées sont sérialisés avec le paquet `encoding/json` standard de Go. Notez que les méthodes `MarshalJSON` déclarées sur un récepteur pointeur ne sont appelées que lorsque la valeur à sérialiser est adressable. Un sérialiseur personnalisé peut donc être appelé pour les éléments d’une tranche renvoyée, mais pas pour une structure renvoyée par valeur :

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

Pour garantir un résultat cohérent, définissez `MarshalJSON` sur un récepteur valeur lorsque le type peut être copié sans risque, ou renvoyez un pointeur vers la structure (par exemple, `*User`). Il s’agit d’un comportement du paquet `encoding/json` de Go, et non d’une différence dans la génération des liaisons par Wails.

### Avec des commentaires

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

**JavaScript généré :**

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

**Les commentaires deviennent des commentaires JSDoc !** Votre IDE les affiche.

### Structures imbriquées

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

**JavaScript généré :**

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

**Utilisation :**

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

### Tableaux et tranches

```go
type Team struct {
    Name    string   `json:"name"`
    Members []string `json:"members"`
}
```

**JavaScript généré :**

```javascript
export class Team {
    /** @type {string} */
    name = ""
    
    /** @type {string[]} */
    members = []
}
```

**Utilisation :**

```javascript
const team = new Team({
    name: "Engineering",
    members: ["Alice", "Bob", "Charlie"]
})
```

### Tables associatives

```go
type Config struct {
    Settings map[string]string `json:"settings"`
}
```

**JavaScript généré :**

```javascript
export class Config {
    /** @type {Record<string, string>} */
    settings = {}
}
```

**Utilisation :**

```javascript
const config = new Config({
    settings: {
        theme: "dark",
        language: "en"
    }
})
```

## Correspondance des types

### Types primitifs

| Type Go | JavaScript/TypeScript |
| --- | --- |
| `string` | `string` |
| `bool` | `boolean` |
| `int`, `int8`, `int16`, `int32`, `int64` | `number` |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | `number` |
| `float32`, `float64` | `number` |
| `byte` | `number` |
| `rune` | `number` |

### Types complexes

| Type Go | JavaScript/TypeScript | Remarques |
| --- | --- | --- |
| `[]T` | `T[]` | - |
| `[N]T` | `T[]` | - |
| `map[string]T` | `{ [_: string]: T }` | table associative à clés de type chaîne |
| `map[K]V` | `{ [_ in K]?: V }` | un type `K` autre qu’une chaîne est rendu sous forme de type mappé, et **non** sous forme de `Map` JS |
| `[]byte` | `string` | encodé en base64 |
| `struct` | `class` / `interface` | avec des champs |
| `time.Time` | `any` | sérialisé à l’exécution sous forme de chaîne RFC3339Nano |
| `*T` | `T \| null` | un pointeur indique une valeur nullable |
| `any` / `interface{}` | `any` | - |
| `error` | `any` / `Exception` | Exception si utilisé comme valeur de retour, sinon any |

### Types non pris en charge

- `chan T` (canaux)
- `func()` (fonctions)
- Interfaces complexes (sauf `interface{}`)
- Champs non exportés (en minuscules)

## Utiliser les modèles

### Créer des instances

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

### Transmettre à Go

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

### Recevoir depuis Go

```javascript
import { GetUser } from './bindings/changeme/userservice'

const user = await GetUser(1)

// user is already a User instance
console.log(user.name)
console.log(user.email)
console.log(user.createdAt.toISOString())
```

### Mettre à jour les modèles

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

## Prise en charge de TypeScript

### TypeScript généré

```bash
wails3 generate bindings -ts
```

**Généré :**

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

### Utilisation dans TypeScript

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

## Modèles avancés

### Champs facultatifs

```go
type User struct {
    ID       int     `json:"id"`
    Name     string  `json:"name"`
    Nickname *string `json:"nickname,omitempty"`
}
```

**JavaScript :**

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

### Énumérations

Le générateur de liaisons détecte automatiquement les types nommés Go associés à des constantes et génère des énumérations TypeScript ou des objets JavaScript const, avec notamment un membre `$zero` pour la valeur zéro de Go et la conservation intégrale de JSDoc.

```go
type UserRole string

const (
    RoleAdmin UserRole = "admin"
    RoleUser  UserRole = "user"
    RoleGuest UserRole = "guest"
)
```

**TypeScript généré :**

```typescript
export enum UserRole {
    $zero = "",
    RoleAdmin = "admin",
    RoleUser = "user",
    RoleGuest = "guest",
}
```

**Utilisation :**

```javascript
import { User, UserRole } from './bindings/changeme/models'

const admin = new User({
    name: "Admin",
    role: UserRole.RoleAdmin
})
```

Pour une présentation complète des énumérations de chaînes, des énumérations d’entiers, des alias de types, des énumérations de paquets importés et des limitations, consultez la page dédiée aux **[énumérations](/features/bindings/enums/)**.

### Validation

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

### Sérialisation

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

## Exemple complet

**Go :**

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

**JavaScript :**

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

## Bonnes pratiques

### ✅ À faire

- **Utilisez des balises JSON** — Contrôlez les noms des champs
- **Ajoutez des commentaires** — Ils deviennent des commentaires JSDoc
- **Utilisez time.Time** — Il est converti en Date
- **Validez côté Go** — Ne faites pas confiance au frontend
- **Gardez les modèles simples** — Utilisez-les uniquement comme conteneurs de données
- **Utilisez des pointeurs pour les champs facultatifs** — `*string` pour les valeurs nullables

### ❌ À ne pas faire

- **N’ajoutez pas de méthodes aux structures Go** — Conservez-les comme données
- **N’utilisez pas de champs non exportés** — Ils ne seront pas liés
- **N’utilisez pas d’interfaces complexes** — Elles ne sont pas prises en charge
- **N’oubliez pas les balises JSON** — Les noms des champs sont importants
- **N’imbriquez pas trop profondément** — Restez simple

## Étapes suivantes

@cards{cols="2"}
🚀 Liaisons de méthodes
Découvrez comment lier des méthodes Go.

[En savoir plus →](/features/bindings/methods/)

---
◆ Services
Organisez le code avec des services.

[En savoir plus →](/features/bindings/services/)

---
✓ Bonnes pratiques
Modèles de conception des liaisons.

[En savoir plus →](/features/bindings/best-practices/)

---
★ Pont Go-frontend
Comprenez le mécanisme du pont.

[En savoir plus →](/concepts/bridge/)

@end

---

**Des questions ?** Posez-les sur [Discord](https://discord.gg/JDdSxwjhGf) ou consultez les [exemples de liaisons](https://github.com/wailsapp/wails/tree/master/v3/examples/binding).
