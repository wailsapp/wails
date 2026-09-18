---
title: "Fortgeschrittene Bindings"
description: "Fortgeschrittene Binding-Techniken einschließlich Direktiven, Code-Injektion und benutzerdefinierten IDs"
slug: "features/bindings/advanced"
sourcePath: "features/bindings/advanced.md"
---

Dieser Leitfaden behandelt fortgeschrittene Techniken zum Anpassen und Optimieren der Binding-Generierung in Wails v3.

## Generierten Code mit Direktiven anpassen

### Benutzerdefinierten Code injizieren

Mit der Direktive `//wails:inject` können Sie benutzerdefinierten JavaScript-/TypeScript-Code in die generierten Bindings injizieren:

```go
//wails:inject console.log("Hello from Wails!");
type MyService struct {}

func (s *MyService) Greet(name string) string {
    return "Hello, " + name
}
```

Dadurch wird der angegebene Code in die generierte JavaScript-/TypeScript-Datei für den Dienst `MyService` injiziert.

Mit bedingter Injektion können Sie auch bestimmte Ausgabeformate ansprechen:

```go
//wails:inject j*:console.log("Hello JS!");  // JavaScript only
//wails:inject t*:console.log("Hello TS!");  // TypeScript only
```

### Zusätzliche Dateien einbinden

Mit der Direktive `//wails:include` können Sie zusätzliche Dateien in die generierten Bindings einbinden:

```go
//wails:include js/*.js
package mypackage
```

Diese Direktive wird üblicherweise in Dokumentationskommentaren eines Pakets verwendet, um zusätzliche JavaScript-/TypeScript-Dateien in die generierten Bindings einzubinden.

### Interne Typen und Methoden kennzeichnen

Die Direktive `//wails:internal` kennzeichnet einen Typ oder eine Methode als intern und verhindert dadurch den Export an das Frontend:

```go
//wails:internal
type InternalModel struct {
    Field string
}

//wails:internal
func (s *MyService) InternalMethod() {}
```

Dies ist für Typen und Methoden nützlich, die nur intern von Ihrem Go-Code verwendet werden und nicht für das Frontend verfügbar sein sollen.

### Methoden ignorieren

Die Direktive `//wails:ignore` ignoriert eine Methode bei der Binding-Generierung vollständig:

```go
//wails:ignore
func (s *MyService) IgnoredMethod() {}
```

Dies ähnelt `//wails:internal`, ignoriert die Methode jedoch vollständig, statt sie als intern zu kennzeichnen.

### Benutzerdefinierte Methoden-IDs

Die Direktive `//wails:id` legt eine benutzerdefinierte ID für eine Methode fest und überschreibt damit die standardmäßige hashbasierte ID:

```go
//wails:id 42
func (s *MyService) CustomIDMethod() {}
```

Dies kann beim Refactoring von Code helfen, die Kompatibilität zu erhalten.

## Mit komplexen Typen arbeiten

### Verschachtelte Structs

Der Binding-Generator verarbeitet verschachtelte Structs automatisch:

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

Der generierte JavaScript-/TypeScript-Code enthält Klassen für `Person` und `Address`.

### Maps und Slices

Maps und Slices werden ebenfalls automatisch verarbeitet:

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

In JavaScript werden Maps als Objekte und Slices als Arrays dargestellt. In TypeScript werden Maps als `Record<K, V>` und Slices als `T[]` dargestellt.

### Generische Typen

Der Binding-Generator unterstützt generische Typen:

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

Der generierte TypeScript-Code enthält eine generische Klasse für `Result`:

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

Mit dem Flag `-i` kann der Binding-Generator TypeScript-Interfaces anstelle von Klassen generieren:

```bash
wails3 generate bindings -ts -i
```

Dadurch werden TypeScript-Interfaces für alle Modelle generiert:

```typescript
export interface Person {
    Name: string;
    Attributes: Record<string, string>;
    Friends: string[];
}
```

## Binding-Generierung optimieren

### Namen anstelle von IDs verwenden

Standardmäßig verwendet der Binding-Generator hashbasierte IDs für Methodenaufrufe. Mit dem Flag `-names` können Sie stattdessen Namen verwenden:

```bash
wails3 generate bindings -names
```

Dadurch wird Code generiert, der die gebundene Methode über ihren **vollständig qualifizierten** Namen (`<package>.<Service>.<Method>`) aufruft und die positionsbezogenen Parameternamen `$0`, `$1`, … verwendet:

```javascript
export function Greet($0) {
    let $resultPromise = $Call.ByName("main.GreetService.Greet", $0);
    return $resultPromise;
}
```

Dies kann den generierten Code lesbarer und leichter zu debuggen machen, aber möglicherweise ist er geringfügig weniger effizient.

### Runtime bündeln

Standardmäßig importiert der generierte Code die Wails-Runtime aus dem npm-Paket `@wailsio/runtime`. Mit dem Flag `-b` können Sie die Runtime mit dem generierten Code bündeln:

```bash
wails3 generate bindings -b
```

Dadurch wird der Runtime-Code direkt in die generierten Dateien eingebunden, sodass das npm-Paket nicht mehr benötigt wird.

### Indexdateien deaktivieren

Wenn Sie die Indexdateien nicht benötigen, können Sie ihre Generierung mit dem Flag `-noindex` deaktivieren:

```bash
wails3 generate bindings -noindex
```

Dies kann nützlich sein, wenn Sie Dienste und Modelle lieber direkt aus den jeweiligen Dateien importieren.

## Praxisbeispiele

### Authentifizierungsdienst

Hier sehen Sie ein Beispiel für einen Authentifizierungsdienst mit benutzerdefinierten Direktiven:

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

### Datenverarbeitungsdienst

Hier sehen Sie ein Beispiel für einen Datenverarbeitungsdienst mit generischen Typen:

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

### Bedingte Code-Injektion

Hier sehen Sie ein Beispiel für die bedingte Code-Injektion für unterschiedliche Ausgabeformate:

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

Dadurch wird für die JavaScript- und TypeScript-Ausgabe jeweils unterschiedlicher Code eingefügt, der passende Typannotationen für die jeweilige Sprache bereitstellt.
