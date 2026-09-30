---
title: "Pont Go-interface utilisateur"
description: "Présentation détaillée de la communication directe entre Go et JavaScript rendue possible par Wails"
slug: "concepts/bridge"
sourcePath: "concepts/bridge.md"
---

## Communication directe entre Go et JavaScript

Wails fournit un pont **direct en mémoire** entre Go et JavaScript, permettant une communication fluide sans surcharge HTTP, frontière entre processus ni goulot d’étranglement lié à la sérialisation.

## Vue d’ensemble

```d2
direction: right

Frontend: Frontend (JavaScript) {
  UI: React/Vue/JavaScript natif {
    shape: rectangle
    style.fill: "#8B5CF6"
  }

  Bindings: Liaisons générées automatiquement {
    shape: rectangle
    style.fill: "#A78BFA"
  }
}

Bridge: Pont Wails {
  Encoder: Encodeur JSON {
    shape: rectangle
    style.fill: "#10B981"
  }

  Router: Routeur de méthodes {
    shape: diamond
    style.fill: "#10B981"
  }

  Decoder: Décodeur JSON {
    shape: rectangle
    style.fill: "#10B981"
  }

  TypeGen: Générateur de types {
    shape: rectangle
    style.fill: "#10B981"
  }
}

Backend: Backend (Go) {
  Services: Vos services {
    shape: rectangle
    style.fill: "#00ADD8"
  }

  Registry: Registre des services {
    shape: rectangle
    style.fill: "#00ADD8"
  }
}

Frontend.UI -> Frontend.Bindings: "import { Method }"
Frontend.Bindings -> Bridge.Encoder: "Appeler Method('arg')"
Bridge.Encoder -> Bridge.Router: Encoder en JSON
Bridge.Router -> Backend.Registry: Rechercher le service
Backend.Registry -> Backend.Services: Invoquer la méthode
Backend.Services -> Bridge.Decoder: Renvoyer le résultat
Bridge.Decoder -> Frontend.Bindings: Décoder en JS
Frontend.Bindings -> Frontend.UI: Résoudre la promesse
Bridge.TypeGen -> Frontend.Bindings: Générer les types
```

**Point essentiel :** ni HTTP, ni IPC, ni frontière entre processus. Uniquement des **appels directs de fonctions** avec **sûreté des types**.

## Fonctionnement étape par étape

### 1. Enregistrement des services (démarrage)

Au démarrage de votre application, Wails analyse vos services :

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

**Ce que fait Wails :**

1. **Analyse la structure** pour rechercher les méthodes exportées
2. **Extrait les informations de type** (paramètres et types de retour)
3. **Crée un registre** qui associe les noms de méthodes aux fonctions
4. **Génère les liaisons TypeScript** avec les définitions de types complètes

### 2. Génération des liaisons (compilation)

Wails génère automatiquement les liaisons TypeScript :

```typescript
// Auto-generated: frontend/bindings/GreetService.ts
export function Greet(name: string): Promise<string>
export function Add(a: number, b: number): Promise<number>
```

**Correspondance des types :**

| Type Go | Type TypeScript |
| --- | --- |
| `string` | `string` |
| `int`, `int32`, `int64` | `number` |
| `float32`, `float64` | `number` |
| `bool` | `boolean` |
| `[]T` | `T[]` |
| `map[string]T` | `Record<string, T>` |
| `struct` | `interface` |
| `time.Time` | `Date` |
| `error` | Exception (levée) |

### 3. Appel depuis l’interface utilisateur (exécution)

Le développeur appelle la méthode Go depuis JavaScript :

```javascript
import { Greet, Add } from './bindings/GreetService'

// Call Go from JavaScript
const greeting = await Greet("World")
console.log(greeting)  // "Hello, World!"

const sum = await Add(5, 3)
console.log(sum)  // 8
```

**Ce qui se passe :**

1. **Fonction de liaison appelée** — `Greet("World")`
2. **Message créé** — `{ service: "GreetService", method: "Greet", args: ["World"] }`
3. **Envoi au pont** — via le pont JavaScript de la WebView
4. **Promesse renvoyée** — attend la réponse

### 4. Traitement par le pont (exécution)

Le pont reçoit le message et le traite :

```d2
direction: down

Receive: Recevoir le message {
  shape: rectangle
  style.fill: "#10B981"
}

Parse: Analyser le JSON {
  shape: rectangle
}

Validate: Valider {
  Check: Le service existe-t-il ? {
    shape: diamond
  }

  CheckMethod: La méthode existe-t-elle ? {
    shape: diamond
  }

  CheckTypes: Les types sont-ils corrects ? {
    shape: diamond
  }
}

Invoke: Invoquer la méthode Go {
  shape: rectangle
  style.fill: "#00ADD8"
}

Encode: Encoder le résultat {
  shape: rectangle
}

Send: Envoyer la réponse {
  shape: rectangle
  style.fill: "#10B981"
}

Error: Envoyer l’erreur {
  shape: rectangle
  style.fill: "#EF4444"
}

Receive -> Parse
Parse -> Validate.Check
Validate.Check -> Validate.CheckMethod: Oui
Validate.Check -> Error: Non
Validate.CheckMethod -> Validate.CheckTypes: Oui
Validate.CheckMethod -> Error: Non
Validate.CheckTypes -> Invoke: Oui
Validate.CheckTypes -> Error: Non
Invoke -> Encode: Succès
Invoke -> Error: Erreur
Encode -> Send
```

**Sécurité :** seuls les services enregistrés et les méthodes exportées peuvent être appelés.

### 5. Exécution Go (exécution)

La méthode Go s’exécute :

```go
func (g *GreetService) Greet(name string) string {
    // This runs in Go
    return g.prefix + name + "!"
}
```

**Contexte d’exécution :**

- S’exécute dans une **goroutine** (sans blocage)
- A accès à **toutes les fonctionnalités de Go** (système de fichiers, réseau, bases de données)
- Peut appeler librement **d’autre code Go**
- Renvoie un résultat ou une erreur

### 6. Réponse (exécution)

Le résultat est renvoyé à JavaScript :

```javascript
// Promise resolves with result
const greeting = await Greet("World")
// greeting = "Hello, World!"
```

**Gestion des erreurs :**

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

## Caractéristiques de performance

### Vitesse

**Surcoût habituel d’un appel :** &lt;1 ms

```
Frontend Call → Bridge → Go Execution → Bridge → Frontend Response
     ↓            ↓           ↓            ↓            ↓
   &lt;0.1ms      &lt;0.1ms      [varies]     &lt;0.1ms      &lt;0.1ms
```

**Comparaison avec les autres solutions :**

- **HTTP/REST :** 5-50 ms (pile réseau, sérialisation)
- **IPC :** 1-10 ms (frontières entre processus, marshalling)
- **Pont Wails :** &lt;1 ms (en mémoire, appel direct)

### Mémoire

**Surcoût par appel :** ~1 Ko (tampon de messages)

**Optimisation sans copie :** les données volumineuses (>1 Mo) utilisent la mémoire partagée lorsque cela est possible.

### Concurrence

**Les appels sont concurrents :**

- Chaque appel s’exécute dans sa propre goroutine
- Plusieurs appels peuvent s’exécuter simultanément
- Aucun blocage entre les appels

```javascript
// These run concurrently
const [result1, result2, result3] = await Promise.all([
    SlowOperation1(),
    SlowOperation2(),
    SlowOperation3(),
])
```

## Système de types

### Types pris en charge

#### Types primitifs

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

#### Tranches et tableaux

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

#### Dictionnaires

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

#### Structures

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

**Balises JSON :** utilisez les balises `json:` pour contrôler les noms des champs dans TypeScript.

#### Temps

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

#### Erreurs

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

### Types non pris en charge

Ces types ne peuvent **pas** être transmis par le pont :

- **Canaux** (`chan T`)
- **Fonctions** (`func()`)
- **Interfaces** (sauf `interface{}` / `any`)
- **Pointeurs** (sauf vers des structures)
- **Champs non exportés** (en minuscules)

**Solution de contournement :** utilisez des identifiants ou des handles :

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

## Modèles avancés

### Transmission du contexte

Les services peuvent accéder au contexte de l’appel :

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

**Le contexte fournit :**

- Fenêtre à l’origine de l’appel
- Instance de l’application
- Métadonnées de la requête

### Transmission de données en continu

Pour les volumes de données importants, utilisez des événements plutôt que des valeurs de retour :

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

### Annulation

Utilisez le contexte pour les opérations annulables :

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

**Remarque :** le contexte est automatiquement annulé lorsque le frontend se déconnecte.

### Opérations par lots

Réduisez la surcharge liée au pont en regroupant les opérations :

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

## Débogage du pont

### Activer la journalisation de débogage

```go
app := application.New(application.Options{
    Name:     "My App",
    LogLevel: slog.LevelDebug, // requires `import "log/slog"`
    // `Logger` is an optional *slog.Logger; the default-logger helper is
    // application.DefaultLogger(slog.Leveler) if you want to construct one explicitly.
})
```

**La sortie affiche :**

- Appels de méthodes
- Paramètres
- Valeurs de retour
- Erreurs
- Informations sur les temps d’exécution

### Inspecter les liaisons générées

Consultez `frontend/bindings/` pour voir le code TypeScript généré :

```javascript
// frontend/bindings/<full-go-import-path>/myservice.js (real generated shape)
import { Call as $Call } from "/wails/runtime.js";

export function MyMethod($0) {
    return $Call.ByID(1234567890, $0); // numeric method ID assigned by the generator
}
```

### Tester directement les services

Testez les services Go sans le frontend :

```go
func TestGreetService(t *testing.T) {
    service := &GreetService{prefix: "Hello, "}
    result := service.Greet("Test")
    if result != "Hello, Test!" {
        t.Errorf("Expected 'Hello, Test!', got '%s'", result)
    }
}
```

## Conseils de performance

### ✅ À faire

- **Regrouper les opérations** — Réduit le nombre d’appels au pont
- **Utiliser des événements pour la transmission en continu** — Ne renvoyez pas de grands tableaux
- **Veiller à la rapidité des méthodes** — Idéalement &lt;100 ms
- **Utiliser des goroutines** — Pour les opérations longues
- **Mettre en cache côté Go** — Évite de répéter les calculs

### ❌ À ne pas faire

- **Ne pas multiplier les appels** — Regroupez-les lorsque c’est possible
- **Ne pas renvoyer d’énormes volumes de données** — Utilisez la pagination ou la transmission en continu
- **Ne pas bloquer** — Utilisez des goroutines pour les opérations longues
- **Ne pas transmettre de types complexes** — Faites simple
- **Ne pas ignorer les erreurs** — Gérez-les toujours

## Sécurité

Le pont est sécurisé par défaut :

1. **Liste d’autorisation uniquement** — Seuls les services enregistrés peuvent être appelés
2. **Validation des types** — Les arguments sont vérifiés par rapport aux types Go
3. **Pas d’eval()** — Le frontend ne peut pas exécuter de code Go arbitraire
4. **Pas d’utilisation abusive de la réflexion** — Seules les méthodes exportées sont accessibles

**Bonnes pratiques :**

- **Validez les entrées** dans Go (ne faites pas confiance au frontend)
- **Utilisez le contexte** pour l’authentification et l’autorisation
- **Limitez la fréquence** des opérations coûteuses
- **Nettoyez** les chemins de fichiers et les entrées utilisateur

## Étapes suivantes

**Système de build** — Découvrez comment Wails compile et regroupe votre application [En savoir plus →](/concepts/build-system/)

**Services** — Explorez en détail le système de services [En savoir plus →](/features/bindings/services/)

**Événements** - Utilisez les événements pour la communication par publication/abonnement [En savoir plus →](/features/events/system/)

---

**Des questions sur le pont ?** Posez-les sur [Discord](https://discord.gg/JDdSxwjhGf) ou consultez les [exemples de liaisons](https://github.com/wailsapp/wails/tree/master/v3/examples/binding).
