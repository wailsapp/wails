---
title: "Normes de codage"
description: "Style de code, conventions et bonnes pratiques pour Wails v3"
slug: "contributing/standards"
sourcePath: "contributing/standards.md"
---

## Style de code et conventions

Le respect de normes de codage cohérentes facilite la lecture, la maintenance et les contributions à la base de code.

## Normes relatives au code Go

### Mise en forme du code

Utilisez les outils de mise en forme Go standard :

```bash
# Format all code
gofmt -w .

# Use goimports for import organization
goimports -w .
```

**Obligatoire :** tout le code Go doit réussir les vérifications de `gofmt` et `goimports` avant tout commit.

### Conventions de nommage

**Packages :**

- Utilisez des minuscules et, si possible, un seul mot
- `package application`, `package events`
- Évitez les traits de soulignement et les mélanges de majuscules et de minuscules

**Noms exportés :**

- Utilisez PascalCase pour les types, les fonctions et les constantes
- `type WebviewWindow struct`, `func NewApplication()`

**Noms non exportés :**

- Utilisez camelCase pour les types, les fonctions et les variables internes
- `type windowImpl struct`, `func createWindow()`

**Interfaces :**

- Nommez-les selon leur comportement : `Reader`, `Writer`, `Handler`
- Pour les interfaces à une seule méthode, utilisez un nom doté du suffixe `-er`

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

### Gestion des erreurs

**Vérifiez toujours les erreurs :**

```go
// Good
result, err := doSomething()
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// Bad - ignoring errors
result, _ := doSomething()
```

**Encapsulez les erreurs :**

```go
// Wrap errors to provide context
if err := validate(); err != nil {
    return fmt.Errorf("validation failed: %w", err)
}
```

**Créez des types d’erreurs personnalisés si nécessaire :**

```go
type ValidationError struct {
    Field string
    Value string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("invalid value %q for field %q", e.Value, e.Field)
}
```

### Commentaires et documentation

**Commentaires de package :**

```go
// Package application provides the core Wails application runtime.
//
// It handles window management, event dispatching, and service lifecycle.
package application
```

**Déclarations exportées :**

```go
// NewApplication creates a new Wails application with the given options.
//
// The application must be started with Run() or RunWithContext().
func NewApplication(opts Options) *Application {
    // ...
}
```

**Commentaires sur l’implémentation :**

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

### Structure des fonctions et des méthodes

**Gardez chaque fonction centrée sur une seule tâche :**

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

**Effectuez des retours anticipés :**

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

### Concurrence

**Utilisez un contexte pour l’annulation :**

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

**Protégez l’état partagé à l’aide de mutex :**

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

**Évitez les fuites de goroutines :**

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

**Nommage des fichiers de test :**

```go
// Implementation: window.go
// Tests: window_test.go
```

**Tests pilotés par table :**

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

## Normes JavaScript/TypeScript

### Mise en forme du code

Utilisez Prettier pour assurer une mise en forme cohérente :

```json
{
  "semi": false,
  "singleQuote": true,
  "tabWidth": 2,
  "trailingComma": "es5"
}
```

### Conventions de nommage

**Variables et fonctions :**

- camelCase : `const userName = "John"`

**Classes et types :**

- PascalCase : `class WindowManager`

**Constantes :**

- UPPER<em>SNAKE</em>CASE : `const MAX_RETRIES = 3`

### TypeScript

**Utilisez des types explicites :**

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

**Définissez des interfaces :**

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

## Format des messages de commit

Utilisez les [Conventional Commits](https://www.conventionalcommits.org/) :

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types :**

- `feat` : nouvelle fonctionnalité
- `fix` : correction de bug
- `docs` : modifications de la documentation
- `refactor` : refactorisation du code
- `test` : ajout ou mise à jour de tests
- `chore` : tâches de maintenance

**Exemples :**

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

## Directives pour les demandes de tirage

### Avant de soumettre

- [ ] Le code passe `gofmt` et `goimports`
- [ ] Tous les tests réussissent (`go test ./...`)
- [ ] Le nouveau code est couvert par des tests
- [ ] La documentation est mise à jour si nécessaire
- [ ] Les messages de commit respectent les conventions
- [ ] Aucun conflit de fusion avec `master`

### Modèle de description de demande de tirage

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

## Processus de revue du code

### En tant que réviseur

- Soyez constructif et respectueux
- Concentrez-vous sur la qualité du code, pas sur vos préférences personnelles
- Expliquez pourquoi vous suggérez ces modifications
- Approuvez dès que vous êtes satisfait

### En tant qu’auteur

- Répondez à tous les commentaires
- Demandez des précisions si nécessaire
- Apportez les modifications demandées ou expliquez pourquoi vous ne le faites pas
- Soyez réceptif aux retours

## Bonnes pratiques

### Performances

- Évitez toute optimisation prématurée
- Effectuez un profilage avant d’optimiser
- Utilisez des benchmarks pour le code critique en matière de performances

```go
func BenchmarkProcess(b *testing.B) {
    for i := 0; i < b.N; i++ {
        process(testData)
    }
}
```

### Sécurité

- Validez toutes les entrées utilisateur
- Nettoyez les données avant de les afficher
- Utilisez `crypto/rand` pour les données aléatoires
- Ne consignez jamais d’informations sensibles dans les journaux

### Documentation

- Documentez les API exportées
- Incluez des exemples dans la documentation
- Mettez à jour la documentation lorsque vous modifiez les API
- Maintenez les fichiers README à jour

## Code propre à une plateforme

### Nommage des fichiers

```
window.go           // Common interface
window_darwin.go    // macOS implementation
window_windows.go   // Windows implementation
window_linux.go     // Linux implementation
```

### Balises de compilation

```go
//go:build darwin

package application

// macOS-specific code
```

## Analyse statique

Exécutez les outils d’analyse statique avant de créer un commit :

```bash
# golangci-lint (recommended)
golangci-lint run

# Individual linters
go vet ./...
staticcheck ./...
```

## Des questions ?

Si vous avez un doute sur l’une de ces normes :

- Consultez le code existant pour trouver des exemples
- Posez votre question sur [Discord](https://discord.gg/JDdSxwjhGf)
- Ouvrez une discussion sur GitHub
