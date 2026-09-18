---
title: "Go-Frontend-Bridge"
description: "Detaillierte Erläuterung, wie Wails die direkte Kommunikation zwischen Go und JavaScript ermöglicht"
slug: "concepts/bridge"
sourcePath: "concepts/bridge.md"
---

## Direkte Go-JavaScript-Kommunikation

Wails stellt eine **direkte In-Memory-Bridge** zwischen Go und JavaScript bereit und ermöglicht so eine nahtlose Kommunikation ohne HTTP-Overhead, Prozessgrenzen oder Serialisierungsengpässe.

## Das Gesamtbild

```d2
direction: right

Frontend: Frontend (JavaScript) {
  UI: React/Vue/Vanilla {
    shape: rectangle
    style.fill: "#8B5CF6"
  }

  Bindings: Automatisch generierte Bindings {
    shape: rectangle
    style.fill: "#A78BFA"
  }
}

Bridge: Wails-Bridge {
  Encoder: JSON-Encoder {
    shape: rectangle
    style.fill: "#10B981"
  }

  Router: Methoden-Router {
    shape: diamond
    style.fill: "#10B981"
  }

  Decoder: JSON-Decoder {
    shape: rectangle
    style.fill: "#10B981"
  }

  TypeGen: Typgenerator {
    shape: rectangle
    style.fill: "#10B981"
  }
}

Backend: Backend (Go) {
  Services: Ihre Dienste {
    shape: rectangle
    style.fill: "#00ADD8"
  }

  Registry: Dienstregister {
    shape: rectangle
    style.fill: "#00ADD8"
  }
}

Frontend.UI -> Frontend.Bindings: "import { Method }"
Frontend.Bindings -> Bridge.Encoder: "Method('arg') aufrufen"
Bridge.Encoder -> Bridge.Router: Als JSON codieren
Bridge.Router -> Backend.Registry: Dienst suchen
Backend.Registry -> Backend.Services: Methode aufrufen
Backend.Services -> Bridge.Decoder: Ergebnis zurückgeben
Bridge.Decoder -> Frontend.Bindings: In JS decodieren
Frontend.Bindings -> Frontend.UI: Promise wird erfüllt
Bridge.TypeGen -> Frontend.Bindings: Typen generieren
```

**Wichtige Erkenntnis:** Kein HTTP, keine IPC, keine Prozessgrenzen. Nur **direkte Funktionsaufrufe** mit **Typsicherheit**.

## Funktionsweise Schritt für Schritt

### 1. Service-Registrierung (beim Start)

Beim Start Ihrer Anwendung durchsucht Wails Ihre Services:

```go
type GreetService struct {
    prefix string
}

func (g *GreetService) Greet(name string) string {
    return g.prefix + name + "!"
}

func (g *GreetService) Add(a, b int) int {
    return a + b
}

// Register service
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&GreetService{prefix: "Hello, "}),
    },
})
```

**Was Wails ausführt:**

1. **Durchsucht die Struktur** nach exportierten Methoden
2. **Extrahiert Typinformationen** (Parameter, Rückgabetypen)
3. **Erstellt eine Registry**, die Methodennamen Funktionen zuordnet
4. **Generiert TypeScript-Bindings** mit vollständigen Typdefinitionen

### 2. Binding-Generierung (Build-Zeit)

Wails generiert automatisch TypeScript-Bindings:

```typescript
// Auto-generated: frontend/bindings/GreetService.ts
export function Greet(name: string): Promise<string>
export function Add(a: number, b: number): Promise<number>
```

**Typzuordnung:**

| Go-Typ | TypeScript-Typ |
| --- | --- |
| `string` | `string` |
| `int`, `int32`, `int64` | `number` |
| `float32`, `float64` | `number` |
| `bool` | `boolean` |
| `[]T` | `T[]` |
| `map[string]T` | `Record<string, T>` |
| `struct` | `interface` |
| `time.Time` | `Date` |
| `error` | Exception (ausgelöst) |

### 3. Frontend-Aufruf (Laufzeit)

Der Entwickler ruft die Go-Methode aus JavaScript auf:

```javascript
import { Greet, Add } from './bindings/GreetService'

// Call Go from JavaScript
const greeting = await Greet("World")
console.log(greeting)  // "Hello, World!"

const sum = await Add(5, 3)
console.log(sum)  // 8
```

**Was geschieht:**

1. **Binding-Funktion aufgerufen** – `Greet("World")`
2. **Nachricht erstellt** – `{ service: "GreetService", method: "Greet", args: ["World"] }`
3. **An die Bridge gesendet** – über die JavaScript-Bridge des WebViews
4. **Promise zurückgegeben** – wartet auf die Antwort

### 4. Verarbeitung durch die Bridge (Laufzeit)

Die Bridge empfängt und verarbeitet die Nachricht:

```d2
direction: down

Receive: Nachricht empfangen {
  shape: rectangle
  style.fill: "#10B981"
}

Parse: JSON parsen {
  shape: rectangle
}

Validate: Validieren {
  Check: Dienst vorhanden? {
    shape: diamond
  }

  CheckMethod: Methode vorhanden? {
    shape: diamond
  }

  CheckTypes: Typen korrekt? {
    shape: diamond
  }
}

Invoke: Go-Methode aufrufen {
  shape: rectangle
  style.fill: "#00ADD8"
}

Encode: Ergebnis codieren {
  shape: rectangle
}

Send: Antwort senden {
  shape: rectangle
  style.fill: "#10B981"
}

Error: Fehler senden {
  shape: rectangle
  style.fill: "#EF4444"
}

Receive -> Parse
Parse -> Validate.Check
Validate.Check -> Validate.CheckMethod: Ja
Validate.Check -> Error: Nein
Validate.CheckMethod -> Validate.CheckTypes: Ja
Validate.CheckMethod -> Error: Nein
Validate.CheckTypes -> Invoke: Ja
Validate.CheckTypes -> Error: Nein
Invoke -> Encode: Erfolg
Invoke -> Error: Fehler
Encode -> Send
```

**Sicherheit:** Nur registrierte Dienste und exportierte Methoden können aufgerufen werden.

### 5. Ausführung in Go (Laufzeit)

Die Go-Methode wird ausgeführt:

```go
func (g *GreetService) Greet(name string) string {
    // This runs in Go
    return g.prefix + name + "!"
}
```

**Ausführungskontext:**

- Wird in einer **Goroutine** ausgeführt (nicht blockierend)
- Hat Zugriff auf **alle Go-Funktionen** (Dateisystem, Netzwerk, Datenbanken)
- Kann uneingeschränkt **anderen Go-Code** aufrufen
- Gibt ein Ergebnis oder einen Fehler zurück

### 6. Antwort (Laufzeit)

Das Ergebnis wird an JavaScript zurückgesendet:

```javascript
// Promise resolves with result
const greeting = await Greet("World")
// greeting = "Hello, World!"
```

**Fehlerbehandlung:**

```go
func (g *GreetService) Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}
```

```javascript
try {
    const result = await Divide(10, 0)
} catch (error) {
    console.error("Go error:", error)  // "division by zero"
}
```

## Leistungsmerkmale

### Geschwindigkeit

**Typischer Aufruf-Overhead:** &lt;1 ms

```
Frontend Call → Bridge → Go Execution → Bridge → Frontend Response
     ↓            ↓           ↓            ↓            ↓
   &lt;0.1ms      &lt;0.1ms      [varies]     &lt;0.1ms      &lt;0.1ms
```

**Im Vergleich zu Alternativen:**

- **HTTP/REST:** 5-50 ms (Netzwerk-Stack, Serialisierung)
- **IPC:** 1-10 ms (Prozessgrenzen, Marshalling)
- **Wails-Bridge:** &lt;1 ms (im Arbeitsspeicher, direkter Aufruf)

### Arbeitsspeicher

**Overhead pro Aufruf:** ~1 KB (Nachrichtenpuffer)

**Zero-Copy-Optimierung:** Große Datenmengen (>1 MB) verwenden nach Möglichkeit gemeinsam genutzten Speicher.

### Nebenläufigkeit

**Aufrufe erfolgen nebenläufig:**

- Jeder Aufruf wird in einer eigenen Goroutine ausgeführt
- Mehrere Aufrufe können gleichzeitig ausgeführt werden
- Aufrufe blockieren einander nicht

```javascript
// These run concurrently
const [result1, result2, result3] = await Promise.all([
    SlowOperation1(),
    SlowOperation2(),
    SlowOperation3(),
])
```

## Typsystem

### Unterstützte Typen

#### Primitive Typen

```go
// Go
func Example(
    s string,
    i int,
    f float64,
    b bool,
) (string, int, float64, bool) {
    return s, i, f, b
}
```

```typescript
// TypeScript (auto-generated)
function Example(
    s: string,
    i: number,
    f: number,
    b: boolean,
): Promise<[string, number, number, boolean]>
```

#### Slices und Arrays

```go
// Go
func Sum(numbers []int) int {
    total := 0
    for _, n := range numbers {
        total += n
    }
    return total
}
```

```typescript
// TypeScript
function Sum(numbers: number[]): Promise<number>

// Usage
const total = await Sum([1, 2, 3, 4, 5])  // 15
```

#### Maps

```go
// Go
func GetConfig() map[string]interface{} {
    return map[string]interface{}{
        "theme": "dark",
        "fontSize": 14,
        "enabled": true,
    }
}
```

```typescript
// TypeScript
function GetConfig(): Promise<Record<string, any>>

// Usage
const config = await GetConfig()
console.log(config.theme)  // "dark"
```

#### Structs

```go
// Go
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func GetUser(id int) (*User, error) {
    return &User{
        ID:    id,
        Name:  "Alice",
        Email: "alice@example.com",
    }, nil
}
```

```typescript
// TypeScript (auto-generated)
interface User {
    id: number
    name: string
    email: string
}

function GetUser(id: number): Promise<User>

// Usage
const user = await GetUser(1)
console.log(user.name)  // "Alice"
```

**JSON-Tags:** Verwenden Sie `json:`-Tags, um die Feldnamen in TypeScript festzulegen.

#### Zeit

```go
// Go
func GetTimestamp() time.Time {
    return time.Now()
}
```

```typescript
// TypeScript
function GetTimestamp(): Promise<Date>

// Usage
const timestamp = await GetTimestamp()
console.log(timestamp.toISOString())
```

#### Fehler

```go
// Go
func Validate(input string) error {
    if input == "" {
        return errors.New("input cannot be empty")
    }
    return nil
}
```

```typescript
// TypeScript
function Validate(input: string): Promise<void>

// Usage
try {
    await Validate("")
} catch (error) {
    console.error(error)  // "input cannot be empty"
}
```

### Nicht unterstützte Typen

Diese Typen können **nicht** über die Bridge übergeben werden:

- **Kanäle** (`chan T`)
- **Funktionen** (`func()`)
- **Interfaces** (außer `interface{}` / `any`)
- **Zeiger** (außer auf Structs)
- **Nicht exportierte Felder** (kleingeschrieben)

**Problemumgehung:** IDs oder Handles verwenden:

```go
// ❌ Can't pass file handle
func OpenFile(path string) (*os.File, error) {
    return os.Open(path)
}

// ✅ Return file ID instead
var files = make(map[string]*os.File)

func OpenFile(path string) (string, error) {
    file, err := os.Open(path)
    if err != nil {
        return "", err
    }
    id := generateID()
    files[id] = file
    return id, nil
}

func ReadFile(id string) ([]byte, error) {
    file := files[id]
    return io.ReadAll(file)
}

func CloseFile(id string) error {
    file := files[id]
    delete(files, id)
    return file.Close()
}
```

## Fortgeschrittene Muster

### Kontextübergabe

Services können auf den Aufrufkontext zugreifen:

```go
type UserService struct{}

func (s *UserService) GetCurrentUser(ctx context.Context) (*User, error) {
    // Access the calling window via the context value
    window, _ := ctx.Value(application.WindowKey).(application.Window)
    _ = window

    // Access the application
    app := application.Get()
    _ = app

    // Your logic
    return getCurrentUser(), nil
}
```

**Der Kontext stellt Folgendes bereit:**

- Fenster, das den Aufruf ausgelöst hat
- Anwendungsinstanz
- Anfragemetadaten

### Datenstreaming

Für große Datenmengen Ereignisse statt Rückgabewerten verwenden:

```go
func ProcessLargeFile(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    lineNum := 0

    for scanner.Scan() {
        lineNum++
        // Emit progress events
        app.Event.Emit("file-progress", map[string]interface{}{
            "line": lineNum,
            "text": scanner.Text(),
        })
    }

    return scanner.Err()
}
```

```javascript
import { Events } from '@wailsio/runtime'
import { ProcessLargeFile } from './bindings/FileService'

// Listen for progress
Events.On('file-progress', (data) => {
    console.log(`Line ${data.line}: ${data.text}`)
})

// Start processing
await ProcessLargeFile('/path/to/large/file.txt')
```

### Abbruch

Für abbrechbare Vorgänge den Kontext verwenden:

```go
func LongRunningTask(ctx context.Context) error {
    for i := 0; i < 1000; i++ {
        // Check if cancelled
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            // Continue work
            time.Sleep(100 * time.Millisecond)
        }
    }
    return nil
}
```

**Hinweis:** Wenn das Frontend die Verbindung trennt, wird der Kontext automatisch abgebrochen.

### Stapeloperationen

Durch Stapelverarbeitung den Bridge-Overhead reduzieren:

```go
// ❌ Inefficient: N bridge calls
for _, item := range items {
    await ProcessItem(item)
}

// ✅ Efficient: 1 bridge call
await ProcessItems(items)
```

```go
func ProcessItems(items []Item) ([]Result, error) {
    results := make([]Result, len(items))
    for i, item := range items {
        results[i] = processItem(item)
    }
    return results, nil
}
```

## Bridge debuggen

### Debug-Protokollierung aktivieren

```go
app := application.New(application.Options{
    Name:     "My App",
    LogLevel: slog.LevelDebug, // requires `import "log/slog"`
    // `Logger` is an optional *slog.Logger; the default-logger helper is
    // application.DefaultLogger(slog.Leveler) if you want to construct one explicitly.
})
```

**Die Ausgabe zeigt:**

- Methodenaufrufe
- Parameter
- Rückgabewerte
- Fehler
- Zeitinformationen

### Generierte Bindings prüfen

In `frontend/bindings/` das generierte TypeScript prüfen:

```javascript
// frontend/bindings/<full-go-import-path>/myservice.js (real generated shape)
import { Call as $Call } from "/wails/runtime.js";

export function MyMethod($0) {
    return $Call.ByID(1234567890, $0); // numeric method ID assigned by the generator
}
```

### Services direkt testen

Go-Services ohne Frontend testen:

```go
func TestGreetService(t *testing.T) {
    service := &GreetService{prefix: "Hello, "}
    result := service.Greet("Test")
    if result != "Hello, Test!" {
        t.Errorf("Expected 'Hello, Test!', got '%s'", result)
    }
}
```

## Tipps zur Performance

### ✅ Empfohlen

- **Operationen bündeln** – Anzahl der Bridge-Aufrufe reduzieren
- **Ereignisse zum Streamen verwenden** – Keine großen Arrays zurückgeben
- **Methoden schnell halten** – Ideal sind &lt;100 ms
- **Goroutinen verwenden** – Für lang laufende Vorgänge
- **Auf Go-Seite zwischenspeichern** – Wiederholte Berechnungen vermeiden

### ❌ Nicht empfohlen

- **Nicht übermäßig viele Aufrufe durchführen** – Wenn möglich bündeln
- **Keine riesigen Datenmengen zurückgeben** – Paginierung oder Streaming verwenden
- **Nicht blockieren** – Für lang laufende Vorgänge Goroutinen verwenden
- **Keine komplexen Typen übergeben** – Einfach halten
- **Fehler nicht ignorieren** – Immer behandeln

## Sicherheit

Die Bridge ist standardmäßig sicher:

1. **Nur Positivliste** – Nur registrierte Services können aufgerufen werden
2. **Typvalidierung** – Argumente werden anhand der Go-Typen geprüft
3. **Kein eval()** – Das Frontend kann keinen beliebigen Go-Code ausführen
4. **Kein Missbrauch von Reflection** – Nur exportierte Methoden sind zugänglich

**Bewährte Vorgehensweisen:**

- **Eingaben** in Go validieren (dem Frontend nicht vertrauen)
- **Kontext** für Authentifizierung und Autorisierung verwenden
- **Aufrufrate** aufwendiger Vorgänge begrenzen
- **Dateipfade und Benutzereingaben** bereinigen

## Nächste Schritte

**Build-System** – Erfahre, wie Wails deine Anwendung erstellt und bündelt [Mehr erfahren →](/concepts/build-system/)

**Services** – Vertiefe dein Wissen über das Service-System [Mehr erfahren →](/features/bindings/services/)

**Ereignisse** – Verwenden Sie Ereignisse für die Publish/Subscribe-Kommunikation [Mehr erfahren →](/features/events/system/)

---

**Fragen zur Bridge?** Fragen Sie auf [Discord](https://discord.gg/JDdSxwjhGf) nach oder sehen Sie sich die [Binding-Beispiele](https://github.com/wailsapp/wails/tree/master/v3/examples/binding) an.
