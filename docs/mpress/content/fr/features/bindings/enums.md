---
title: "Énumérations"
description: "Génération automatique d’énumérations à partir de constantes Go"
slug: "features/bindings/enums"
sourcePath: "features/bindings/enums.md"
---

## Liaisons d’énumérations

Le générateur de liaisons v3 de Wails **détecte automatiquement les types de constantes Go et génère des énumérations TypeScript ou des objets constants JavaScript**. Aucun enregistrement ni aucune configuration ne sont nécessaires : définissez simplement vos types et vos constantes en Go, et le générateur s’occupe du reste.

@note{type="info"}
Contrairement à Wails v2, il n’est **pas nécessaire d’appeler `EnumBind`** ni d’enregistrer manuellement les énumérations. Le générateur les découvre automatiquement dans votre code source.

@end

## Démarrage rapide

**Définissez un type nommé avec des constantes en Go :**

```go
type Status string

const (
    StatusActive  Status = "active"
    StatusPending Status = "pending"
    StatusClosed  Status = "closed"
)
```

**Utilisez ce type dans une structure ou une méthode de service :**

```go
type Ticket struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Status Status `json:"status"`
}
```

**Générez les liaisons :**

```bash
wails3 generate bindings
```

La sortie du générateur indique le nombre d’énumérations en plus de celui des modèles :

```
3 Enums, 5 Models
```

**Utilisez-les dans votre interface :**

```javascript
import { Ticket, Status } from './bindings/changeme/models'

const ticket = new Ticket({
    id: 1,
    title: "Bug report",
    status: Status.StatusActive
})
```

**C’est tout !** Le type d’énumération est imposé à la fois en Go et en JavaScript/TypeScript.

## Définition des énumérations

Dans Wails, une énumération est un **type nommé** doté d’un type de base élémentaire, associé à des **déclarations de constantes** de ce type.

### Énumérations de chaînes

```go
// Title is a title
type Title string

const (
    // Mister is a title
    Mister Title = "Mr"
    Miss   Title = "Miss"
    Ms     Title = "Ms"
    Mrs    Title = "Mrs"
    Dr     Title = "Dr"
)
```

**TypeScript généré :**

```typescript
/**
 * Title is a title
 */
export enum Title {
    /**
     * The Go zero value for the underlying type of the enum.
     */
    $zero = "",

    /**
     * Mister is a title
     */
    Mister = "Mr",
    Miss = "Miss",
    Ms = "Ms",
    Mrs = "Mrs",
    Dr = "Dr",
}
```

**JavaScript généré :**

```javascript
/**
 * Title is a title
 * @readonly
 * @enum {string}
 */
export const Title = {
    /**
     * The Go zero value for the underlying type of the enum.
     */
    $zero: "",

    /**
     * Mister is a title
     */
    Mister: "Mr",
    Miss: "Miss",
    Ms: "Ms",
    Mrs: "Mrs",
    Dr: "Dr",
};
```

### Énumérations d’entiers

```go
type Priority int

const (
    PriorityLow    Priority = 0
    PriorityMedium Priority = 1
    PriorityHigh   Priority = 2
)
```

**TypeScript généré :**

```typescript
export enum Priority {
    /**
     * The Go zero value for the underlying type of the enum.
     */
    $zero = 0,

    PriorityLow = 0,
    PriorityMedium = 1,
    PriorityHigh = 2,
}
```

### Énumérations fondées sur des alias de type

Les alias de type Go (`=`) fonctionnent également, mais produisent une sortie légèrement différente : une définition de type accompagnée d’un objet constant, plutôt qu’une déclaration `enum` TypeScript native :

```go
// Age is an integer with some predefined values
type Age = int

const (
    NewBorn    Age = 0
    Teenager   Age = 12
    YoungAdult Age = 18

    // Oh no, some grey hair!
    MiddleAged Age = 50
    Mathusalem Age = 1000 // Unbelievable!
)
```

**TypeScript généré :**

```typescript
/**
 * Age is an integer with some predefined values
 */
export type Age = number;

/**
 * Predefined constants for type Age.
 * @namespace
 */
export const Age = {
    NewBorn: 0,
    Teenager: 12,
    YoungAdult: 18,

    /**
     * Oh no, some grey hair!
     */
    MiddleAged: 50,

    /**
     * Unbelievable!
     */
    Mathusalem: 1000,
};
```

**JavaScript généré :**

```javascript
/**
 * Age is an integer with some predefined values
 * @typedef {number} Age
 */

/**
 * Predefined constants for type Age.
 * @namespace
 */
export const Age = {
    NewBorn: 0,
    Teenager: 12,
    YoungAdult: 18,

    /**
     * Oh no, some grey hair!
     */
    MiddleAged: 50,

    /**
     * Unbelievable!
     */
    Mathusalem: 1000,
};
```

@note{type="tip"}
Les **types nommés** (`type Title string`) génèrent des déclarations `enum` TypeScript natives comportant un membre `$zero`. Les **alias de type** (`type Age = int`) génèrent une paire espace de noms `type` + `const` sans `$zero`.

@end

## La valeur `$zero`

Chaque énumération fondée sur un type nommé comprend un membre spécial `$zero` représentant la **valeur zéro de Go** pour le type de base :

| Type de base | Valeur `$zero` |
| --- | --- |
| `string` | `""` |
| `int`, `int8`, `int16`, `int32`, `int64` | `0` |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | `0` |
| `float32`, `float64` | `0` |
| `bool` | `false` |

Lorsqu’un champ de structure utilise un type d’énumération et qu’aucune valeur n’est fournie, le constructeur utilise `$zero` par défaut :

```typescript
export class Person {
    "Title": Title;

    constructor($$source: Partial<Person> = {}) {
        if (!("Title" in $$source)) {
            this["Title"] = Title.$zero;  // defaults to ""
        }
        Object.assign(this, $$source);
    }
}
```

Cela garantit une initialisation avec contrôle de type lors de la génération de classes : les champs d’énumération ne valent jamais `undefined`. Lors de la génération d’interfaces TypeScript (avec `-i`), aucun constructeur n’est créé et les champs peuvent être absents, comme d’habitude.

## Utilisation des énumérations dans les structures

Lorsqu’un champ de structure possède un type d’énumération, le code généré **conserve ce type** au lieu de revenir au type primitif :

```go
type Person struct {
    Title Title
    Name  string
    Age   Age
}
```

**TypeScript généré :**

```typescript
export class Person {
    "Title": Title;
    "Name": string;
    "Age": Age;

    constructor($$source: Partial<Person> = {}) {
        if (!("Title" in $$source)) {
            this["Title"] = Title.$zero;
        }
        if (!("Name" in $$source)) {
            this["Name"] = "";
        }
        if (!("Age" in $$source)) {
            this["Age"] = 0;
        }

        Object.assign(this, $$source);
    }
}
```

Le champ `Title` est typé `Title`, et non `string`. Votre IDE peut ainsi proposer l’autocomplétion et vérifier intégralement le type des valeurs d’énumération.

## Énumérations provenant de paquets importés

Les énumérations définies dans des paquets distincts sont entièrement prises en charge. Elles sont générées dans le répertoire du paquet correspondant :

```go
// services/types.go
package services

type Title string

const (
    Mister Title = "Mr"
    Miss   Title = "Miss"
    Ms     Title = "Ms"
)
```

```go
// main.go
package main

import "myapp/services"

func (*GreetService) Greet(name string, title services.Title) string {
    return "Hello " + string(title) + " " + name
}
```

L’énumération `Title` est générée dans le fichier de modèles `services`, et les chemins d’importation sont résolus automatiquement :

```typescript
// bindings/changeme/services/models.ts
export enum Title {
    $zero = "",
    Mister = "Mr",
    Miss = "Miss",
    Ms = "Ms",
}
```

## Méthodes des énumérations

Vous pouvez ajouter des méthodes à vos types d’énumération en Go. Elles n’ont aucune incidence sur la génération des liaisons, mais fournissent des fonctionnalités utiles côté serveur :

```go
type Title string

func (t Title) String() string {
    return string(t)
}

const (
    Mister Title = "Mr"
    Miss   Title = "Miss"
)
```

L’énumération générée est identique, que le type possède ou non des méthodes Go.

## Commentaires et documentation

Le générateur conserve les commentaires Go sous forme de JSDoc dans la sortie générée :

- Les **commentaires de type** deviennent le commentaire de documentation de l’énumération
- Les **commentaires de groupes de constantes** deviennent des séparateurs de sections
- Les **commentaires de constantes individuelles** deviennent les commentaires de documentation des membres
- Les **commentaires en fin de ligne** sont conservés lorsque cela est possible

Votre IDE affiche donc la documentation des valeurs d’énumération lorsque vous les survolez.

## Types de base pris en charge

Le générateur de liaisons prend en charge les énumérations utilisant les types de base Go suivants :

| Type Go | Utilisable comme énumération |
| --- | :---: |
| `string` | Oui |
| `int`, `int8`, `int16`, `int32`, `int64` | Oui |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | Oui |
| `float32`, `float64` | Oui |
| `byte` (`uint8`) | Oui |
| `rune` (`int32`) | Oui |
| `bool` | Oui |
| `complex64`, `complex128` | Non |

## Limitations

Les éléments suivants ne sont **pas** pris en charge pour la génération d’énumérations :

- **Types génériques** — Les paramètres de type empêchent la détection des constantes
- **Types dotés d’une implémentation personnalisée de `json.Marshaler` ou de `encoding.TextMarshaler`** — Avec une sérialisation personnalisée, les valeurs générées risquent de ne pas correspondre au comportement à l’exécution ; le générateur ignore donc ces types
- **Constantes dont les valeurs ne peuvent pas être évaluées statiquement ou représentées** — Les constantes doivent avoir des valeurs connues et représentables dans leur type sous-jacent. Les constructions `iota` standard fonctionnent correctement, car le compilateur les résout en valeurs concrètes
- **Types de nombres complexes** — `complex64` et `complex128` ne peuvent pas servir de types sous-jacents à une énumération

## Exemple complet

**Go :**

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
)

// BackgroundType defines the type of background
type BackgroundType string

const (
    BackgroundSolid    BackgroundType = "solid"
    BackgroundGradient BackgroundType = "gradient"
    BackgroundImage    BackgroundType = "image"
)

type BackgroundConfig struct {
    Type  BackgroundType `json:"type"`
    Value string         `json:"value"`
}

type ThemeService struct{}

func (*ThemeService) GetBackground() BackgroundConfig {
    return BackgroundConfig{
        Type:  BackgroundSolid,
        Value: "#ffffff",
    }
}

func (*ThemeService) SetBackground(config BackgroundConfig) error {
    // Apply background
    return nil
}

func main() {
    app := application.New(application.Options{
        Services: []application.Service{
            application.NewService(&ThemeService{}),
        },
    })
    app.Window.New()
    app.Run()
}
```

**Frontend (TypeScript) :**

```typescript
import { GetBackground, SetBackground } from './bindings/changeme/themeservice'
import { BackgroundConfig, BackgroundType } from './bindings/changeme/models'

// Get current background
const bg = await GetBackground()

// Check the type using enum values
if (bg.type === BackgroundType.BackgroundSolid) {
    console.log("Solid background:", bg.value)
}

// Set a new background
await SetBackground(new BackgroundConfig({
    type: BackgroundType.BackgroundGradient,
    value: "linear-gradient(to right, #000, #fff)"
}))
```

**Frontend (JavaScript) :**

```javascript
import { GetBackground, SetBackground } from './bindings/changeme/themeservice'
import { BackgroundConfig, BackgroundType } from './bindings/changeme/models'

// Use enum values for type-safe comparisons
const bg = await GetBackground()

switch (bg.type) {
    case BackgroundType.BackgroundSolid:
        applySolid(bg.value)
        break
    case BackgroundType.BackgroundGradient:
        applyGradient(bg.value)
        break
    case BackgroundType.BackgroundImage:
        applyImage(bg.value)
        break
}
```

## Étapes suivantes

@cards{cols="2"}
📖 Modèles de données
Structures, correspondance des types et génération de modèles.

[En savoir plus →](/features/bindings/models/)

---
🚀 Liaison de méthodes
Liez des méthodes Go au frontend.

[En savoir plus →](/features/bindings/methods/)

---
⚙ Liaison avancée
Directives, injection de code et identifiants personnalisés.

[En savoir plus →](/features/bindings/advanced/)

---
✓ Bonnes pratiques
Modèles de conception pour les liaisons.

[En savoir plus →](/features/bindings/best-practices/)

@end

---

**Des questions ?** Posez-les sur [Discord](https://discord.gg/JDdSxwjhGf) ou consultez les [exemples de liaisons](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/binding).
