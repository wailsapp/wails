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

Sie können auch eine bedingte Injektion verwenden, um bestimmte Ausgabeformate anzusprechen:

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

Die Direktive `//wails:internal` unterdrückt den Export eines Typs oder einer Methode in den generierten JavaScript-/TypeScript-Bindings:

```go
//wails:internal
type InternalModel struct {
    Field string
}

//wails:internal
func (s *MyService) InternalMethod() {}
```

Dies steuert das generierte SDK, nicht den Laufzeitzugriff. In normalen Builds können exportierte Go-Dienstmethoden weiterhin registriert und über die Laufzeit aufgerufen werden, selbst wenn sie im generierten SDK fehlen. Lassen Sie private Operationen nicht exportiert oder definieren Sie sie auf einem nicht gebundenen Typ und erzwingen Sie die Autorisierung innerhalb offengelegter Methoden.

### Methoden ignorieren

Die Direktive `//wails:ignore` ignoriert eine Methode bei der Binding-Generierung vollständig:

```go
//wails:ignore
func (s *MyService) IgnoredMethod() {}
```

Dadurch wird der generierte Wrapper ausgelassen. Wie `//wails:internal` stellt dies keine Zugriffskontrollgrenze für exportierte Go-Dienstmethoden dar.

### Benutzerdefinierte Methoden-IDs

Die Direktive `//wails:id` legt eine benutzerdefinierte ID für eine Methode fest und überschreibt damit die standardmäßige Hash-basierte ID:

```go
//wails:id 42
func (s *MyService) CustomIDMethod() {}
```

Dies kann beim Refactoring von Code hilfreich sein, um die Kompatibilität beizubehalten.

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

Der generierte JavaScript-/TypeScript-Code enthält Klassen sowohl für `Person` als auch für `Address`.

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

Standardmäßig verwendet der Binding-Generator Hash-basierte IDs für Methodenaufrufe. Mit dem Flag `-names` können Sie stattdessen Namen verwenden:

```bash
wails3 generate bindings -names
```

Dadurch wird Code generiert, der die gebundene Methode über ihren **vollqualifizierten** Namen (`<package>.<Service>.<Method>`) aufruft und die positionsbasierten Parameternamen `$0`, `$1`, … verwendet:

```javascript
export function Greet($0) {
    let $resultPromise = $Call.ByName("main.GreetService.Greet", $0);
    return $resultPromise;
}
```

Dadurch kann der generierte Code lesbarer und leichter zu debuggen sein, allerdings möglicherweise auf Kosten einer geringfügig niedrigeren Effizienz.

### Laufzeit bündeln

Standardmäßig importiert der generierte Code die Wails-Laufzeit aus dem npm-Paket `@wailsio/runtime`. Mit dem Flag `-b` können Sie die Laufzeit zusammen mit dem generierten Code bündeln:

```bash
wails3 generate bindings -b
```

Dadurch wird der Laufzeitcode direkt in die generierten Dateien eingebunden, sodass das npm-Paket nicht mehr benötigt wird.

### Indexdateien deaktivieren

Wenn Sie die Indexdateien nicht benötigen, können Sie ihre Generierung mit dem Flag `-noindex` deaktivieren:

```bash
wails3 generate bindings -noindex
```

Dies kann hilfreich sein, wenn Sie Dienste und Modelle lieber direkt aus ihren jeweiligen Dateien importieren.

## Praxisbeispiele

### Authentifizierungsdienst

Das folgende Beispiel zeigt einen Authentifizierungsdienst mit benutzerdefinierten Direktiven:

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

Das folgende Beispiel zeigt einen Datenverarbeitungsdienst mit generischen Typen:

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

Das folgende Beispiel zeigt eine bedingte Code-Injektion für verschiedene Ausgabeformate:

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

Dadurch wird für JavaScript- und TypeScript-Ausgaben jeweils unterschiedlicher Code injiziert, der die passenden Typannotationen für die jeweilige Sprache bereitstellt.
