---
title: "Liaison avancée"
description: "Techniques de liaison avancées, notamment les directives, l’injection de code et les identifiants personnalisés"
slug: "features/bindings/advanced"
sourcePath: "features/bindings/advanced.md"
---

Ce guide présente des techniques avancées pour personnaliser et optimiser le processus de génération des liaisons dans Wails v3.

## Personnaliser le code généré avec des directives

### Injecter du code personnalisé

La directive `//wails:inject` vous permet d’injecter du code JavaScript/TypeScript personnalisé dans les liaisons générées :

```go
//wails:inject console.log("Hello from Wails!");
type MyService struct {}

func (s *MyService) Greet(name string) string {
    return "Hello, " + name
}
```

Le code spécifié sera injecté dans le fichier JavaScript/TypeScript généré pour le service `MyService`.

Vous pouvez également utiliser l’injection conditionnelle pour cibler des formats de sortie précis :

```go
//wails:inject j*:console.log("Hello JS!");  // JavaScript only
//wails:inject t*:console.log("Hello TS!");  // TypeScript only
```

### Inclure des fichiers supplémentaires

La directive `//wails:include` vous permet d’inclure des fichiers supplémentaires avec les liaisons générées :

```go
//wails:include js/*.js
package mypackage
```

Cette directive s’utilise généralement dans les commentaires de documentation du paquet afin d’inclure des fichiers JavaScript/TypeScript supplémentaires avec les liaisons générées.

### Marquer des types et des méthodes comme internes

La directive `//wails:internal` empêche l’exportation d’un type ou d’une méthode dans les liaisons JavaScript/TypeScript générées :

```go
//wails:internal
type InternalModel struct {
    Field string
}

//wails:internal
func (s *MyService) InternalMethod() {}
```

Cela contrôle le SDK généré, et non l’accès à l’exécution. Dans les builds normaux, les méthodes exportées d’un service Go peuvent toujours être enregistrées et appelées par l’intermédiaire de l’environnement d’exécution, même si elles sont absentes du SDK généré. Gardez les opérations privées non exportées ou placez-les sur un type non lié, et appliquez les contrôles d’autorisation dans les méthodes exposées.

### Ignorer des méthodes

La directive `//wails:ignore` ignore complètement une méthode lors de la génération des liaisons :

```go
//wails:ignore
func (s *MyService) IgnoredMethod() {}
```

Cela omet le wrapper généré. Comme `//wails:internal`, cette directive ne constitue pas une barrière de contrôle d’accès pour les méthodes exportées d’un service Go.

### Identifiants de méthode personnalisés

La directive `//wails:id` spécifie un identifiant personnalisé pour une méthode et remplace l’identifiant par défaut fondé sur un hachage :

```go
//wails:id 42
func (s *MyService) CustomIDMethod() {}
```

Cela peut être utile pour préserver la compatibilité lors de la refactorisation du code.

## Utiliser des types complexes

### Structures imbriquées

Le générateur de liaisons gère automatiquement les structures imbriquées :

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

Le code JavaScript/TypeScript généré comprendra des classes pour `Person` et `Address`.

### Maps et slices

Les maps et les slices sont également gérées automatiquement :

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

En JavaScript, les maps sont représentées par des objets et les slices par des tableaux. En TypeScript, les maps sont représentées par `Record<K, V>` et les slices par `T[]`.

### Types génériques

Le générateur de liaisons prend en charge les types génériques :

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

Le code TypeScript généré comprendra une classe générique pour `Result` :

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

### Interfaces

Le générateur de liaisons peut générer des interfaces TypeScript plutôt que des classes à l’aide de l’option `-i` :

```bash
wails3 generate bindings -ts -i
```

Cela générera des interfaces TypeScript pour tous les modèles :

```typescript
export interface Person {
    Name: string;
    Attributes: Record<string, string>;
    Friends: string[];
}
```

## Optimiser la génération des liaisons

### Utiliser des noms plutôt que des identifiants

Par défaut, le générateur de liaisons utilise des identifiants fondés sur un hachage pour les appels de méthodes. Vous pouvez utiliser l’option `-names` pour employer des noms à la place :

```bash
wails3 generate bindings -names
```

Cela générera du code qui appelle la méthode liée par son nom **pleinement qualifié** (`<package>.<Service>.<Method>`) et utilise les noms de paramètres positionnels `$0`, `$1`, … :

```javascript
export function Greet($0) {
    let $resultPromise = $Call.ByName("main.GreetService.Greet", $0);
    return $resultPromise;
}
```

Cela peut rendre le code généré plus lisible et plus facile à déboguer, mais aussi légèrement moins efficace.

### Regrouper l’environnement d’exécution

Par défaut, le code généré importe l’environnement d’exécution Wails depuis le paquet npm `@wailsio/runtime`. Vous pouvez utiliser l’option `-b` pour regrouper l’environnement d’exécution avec le code généré :

```bash
wails3 generate bindings -b
```

Cela inclura directement le code de l’environnement d’exécution dans les fichiers générés, supprimant ainsi la nécessité du paquet npm.

### Désactiver les fichiers d’index

Si vous n’avez pas besoin des fichiers d’index, vous pouvez utiliser l’option `-noindex` pour désactiver leur génération :

```bash
wails3 generate bindings -noindex
```

Cela peut être utile si vous préférez importer les services et les modèles directement depuis leurs fichiers respectifs.

## Exemples concrets

### Service d’authentification

Voici un exemple de service d’authentification avec des directives personnalisées :

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

### Service de traitement des données

Voici un exemple de service de traitement des données avec des types génériques :

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

### Injection conditionnelle de code

Voici un exemple d’injection conditionnelle de code pour différents formats de sortie :

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

Cela injecte un code différent dans les sorties JavaScript et TypeScript, avec des annotations de type adaptées à chaque langage.
