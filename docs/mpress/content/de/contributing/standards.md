---
title: "Programmierstandards"
description: "Codestil, Konventionen und bewährte Methoden für Wails v3"
slug: "contributing/standards"
sourcePath: "contributing/standards.md"
---

## Codestil und Konventionen

Einheitliche Programmierstandards erleichtern es, die Codebasis zu lesen, zu pflegen und zu ihr beizutragen.

## Go-Codestandards

### Codeformatierung

Verwenden Sie die standardmäßigen Go-Formatierungswerkzeuge:

```bash
# Format all code
gofmt -w .

# Use goimports for import organization
goimports -w .
```

**Erforderlich:** Der gesamte Go-Code muss vor dem Committen `gofmt` und `goimports` erfolgreich durchlaufen.

### Namenskonventionen

**Pakete:**

- Kleinschreibung, möglichst ein einzelnes Wort
- `package application`, `package events`
- Vermeiden Sie Unterstriche und gemischte Groß- und Kleinschreibung

**Exportierte Namen:**

- PascalCase für Typen, Funktionen und Konstanten
- `type WebviewWindow struct`, `func NewApplication()`

**Nicht exportierte Namen:**

- camelCase für interne Typen, Funktionen und Variablen
- `type windowImpl struct`, `func createWindow()`

**Schnittstellen:**

- Benennen Sie sie nach ihrem Verhalten: `Reader`, `Writer`, `Handler`
- Schnittstellen mit nur einer Methode: Verwenden Sie Namen mit dem Suffix `-er`

```go
// Good
type Closer interface {
    Close() error
}

// Avoid
type CloseInterface interface {
    Close() error
}
```

### Fehlerbehandlung

**Prüfen Sie Fehler immer:**

```go
// Good
result, err := doSomething()
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// Bad - ignoring errors
result, _ := doSomething()
```

**Verwenden Sie Fehlerverkettung:**

```go
// Wrap errors to provide context
if err := validate(); err != nil {
    return fmt.Errorf("validation failed: %w", err)
}
```

**Erstellen Sie bei Bedarf benutzerdefinierte Fehlertypen:**

```go
type ValidationError struct {
    Field string
    Value string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("invalid value %q for field %q", e.Value, e.Field)
}
```

### Kommentare und Dokumentation

**Paketkommentare:**

```go
// Package application provides the core Wails application runtime.
//
// It handles window management, event dispatching, and service lifecycle.
package application
```

**Exportierte Deklarationen:**

```go
// NewApplication creates a new Wails application with the given options.
//
// The application must be started with Run() or RunWithContext().
func NewApplication(opts Options) *Application {
    // ...
}
```

**Implementierungskommentare:**

```go
// processEvent handles incoming events from the runtime.
// It dispatches to registered handlers and manages event lifecycle.
func (a *Application) processEvent(event *Event) {
    // Validate event before processing
    if event == nil {
        return
    }

    // Find and invoke handlers
    // ...
}
```

### Struktur von Funktionen und Methoden

**Halten Sie Funktionen fokussiert:**

```go
// Good - single responsibility
func (w *Window) setTitle(title string) {
    w.title = title
    w.updateNativeTitle()
}

// Bad - doing too much
func (w *Window) updateEverything() {
    w.setTitle(w.title)
    w.setSize(w.width, w.height)
    w.setPosition(w.x, w.y)
    // ... 20 more operations
}
```

**Verwenden Sie frühzeitige Rückgaben:**

```go
// Good
func validate(input string) error {
    if input == "" {
        return errors.New("empty input")
    }

    if len(input) > 100 {
        return errors.New("input too long")
    }

    return nil
}

// Avoid deep nesting
```

### Nebenläufigkeit

**Verwenden Sie einen Kontext für Abbrüche:**

```go
func (a *Application) RunWithContext(ctx context.Context) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-a.done:
        return nil
    }
}
```

**Schützen Sie gemeinsam genutzten Zustand mit Mutexen:**

```go
type SafeCounter struct {
    mu    sync.Mutex
    count int
}

func (c *SafeCounter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}
```

**Vermeiden Sie Goroutine-Leaks:**

```go
// Good - goroutine has exit condition
func (a *Application) startWorker(ctx context.Context) {
    go func() {
        for {
            select {
            case <-ctx.Done():
                return  // Clean exit
            case work := <-a.workChan:
                a.process(work)
            }
        }
    }()
}
```

### Tests

**Benennung von Testdateien:**

```go
// Implementation: window.go
// Tests: window_test.go
```

**Tabellengesteuerte Tests:**

```go
func TestValidate(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"empty input", "", true},
        {"valid input", "hello", false},
        {"too long", strings.Repeat("a", 101), true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validate(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## JavaScript-/TypeScript-Standards

### Codeformatierung

Verwenden Sie Prettier für eine einheitliche Formatierung:

```json
{
  "semi": false,
  "singleQuote": true,
  "tabWidth": 2,
  "trailingComma": "es5"
}
```

### Namenskonventionen

**Variablen und Funktionen:**

- camelCase: `const userName = "John"`

**Klassen und Typen:**

- PascalCase: `class WindowManager`

**Konstanten:**

- UPPER<em>SNAKE</em>CASE: `const MAX_RETRIES = 3`

### TypeScript

**Verwenden Sie explizite Typen:**

```typescript
// Good
function greet(name: string): string {
    return `Hello, ${name}`
}

// Avoid implicit any
function process(data) {  // Bad
    return data
}
```

**Definieren Sie Schnittstellen:**

```typescript
interface WindowOptions {
    title: string
    width: number
    height: number
}

function createWindow(options: WindowOptions): void {
    // ...
}
```

## Format von Commit-Nachrichten

Verwenden Sie [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Typen:**

- `feat`: Neue Funktion
- `fix`: Fehlerbehebung
- `docs`: Dokumentationsänderungen
- `refactor`: Code-Refactoring
- `test`: Tests hinzufügen oder aktualisieren
- `chore`: Wartungsaufgaben

**Beispiele:**

```
feat(window): add SetAlwaysOnTop method

Implement SetAlwaysOnTop for keeping windows above others.
Adds platform implementations for macOS, Windows, and Linux.

Closes #123
```

```
fix(events): prevent event handler memory leak

Event listeners were not being properly cleaned up when
windows were closed. This adds explicit cleanup in the
window destructor.
```

## Richtlinien für Pull Requests

### Vor dem Einreichen

- [ ] Der Code besteht `gofmt` und `goimports`
- [ ] Alle Tests sind erfolgreich (`go test ./...`)
- [ ] Für neuen Code sind Tests vorhanden
- [ ] Die Dokumentation wurde bei Bedarf aktualisiert
- [ ] Commit-Nachrichten entsprechen den Konventionen
- [ ] Keine Merge-Konflikte mit `master`

### Vorlage für die PR-Beschreibung

```markdown
## Description
Brief description of what this PR does.

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
How was this tested?

## Checklist
- [ ] Tests pass
- [ ] Documentation updated
- [ ] No breaking changes (or documented)
```

## Code-Review-Prozess

### Als Reviewer

- Seien Sie konstruktiv und respektvoll
- Konzentrieren Sie sich auf die Codequalität, nicht auf persönliche Vorlieben
- Begründen Sie vorgeschlagene Änderungen
- Erteilen Sie Ihre Freigabe, sobald Sie zufrieden sind

### Als Autor

- Reagieren Sie auf alle Kommentare
- Bitten Sie bei Bedarf um Klärung
- Nehmen Sie die angeforderten Änderungen vor oder erklären Sie, warum nicht
- Seien Sie offen für Feedback

## Bewährte Vorgehensweisen

### Performance

- Vermeiden Sie verfrühte Optimierungen
- Erstellen Sie vor der Optimierung ein Profil
- Verwenden Sie Benchmarks für performancekritischen Code

```go
func BenchmarkProcess(b *testing.B) {
    for i := 0; i < b.N; i++ {
        process(testData)
    }
}
```

### Sicherheit

- Validieren Sie alle Benutzereingaben
- Bereinigen Sie Daten vor der Anzeige
- Verwenden Sie `crypto/rand` für Zufallsdaten
- Protokollieren Sie niemals vertrauliche Informationen

### Dokumentation

- Dokumentieren Sie exportierte APIs
- Fügen Sie der Dokumentation Beispiele hinzu
- Aktualisieren Sie bei API-Änderungen die Dokumentation
- Halten Sie README-Dateien aktuell

## Plattformspezifischer Code

### Dateibenennung

```
window.go           // Common interface
window_darwin.go    // macOS implementation
window_windows.go   // Windows implementation
window_linux.go     // Linux implementation
```

### Build-Tags

```go
//go:build darwin

package application

// macOS-specific code
```

## Linting

Führen Sie vor dem Committen die Linter aus:

```bash
# golangci-lint (recommended)
golangci-lint run

# Individual linters
go vet ./...
staticcheck ./...
```

## Fragen?

Wenn Sie sich bei einer Vorgabe unsicher sind:

- Suchen Sie im vorhandenen Code nach Beispielen
- Fragen Sie auf [Discord](https://discord.gg/JDdSxwjhGf) nach
- Eröffnen Sie eine Diskussion auf GitHub
