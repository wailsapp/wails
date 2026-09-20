---
title: "Enums"
description: "Automatische Generierung von Enums aus Go-Konstanten"
slug: "features/bindings/enums"
sourcePath: "features/bindings/enums.md"
---

## Enum-Bindings

Der Binding-Generator von Wails v3 **erkennt Typen von Go-Konstanten automatisch und generiert TypeScript-Enums oder JavaScript-const-Objekte**. Weder Registrierung noch Konfiguration sind erforderlich: Definieren Sie einfach Ihre Typen und Konstanten in Go; der Generator erledigt den Rest.

@note{type="info"}
Anders als bei Wails v2 müssen Sie **weder `EnumBind`** aufrufen noch Enums manuell registrieren. Der Generator erkennt sie automatisch in Ihrem Quellcode.

@end

## Schnellstart

**Definieren Sie in Go einen benannten Typ mit Konstanten:**

```go
type Status string

const (
    StatusActive  Status = "active"
    StatusPending Status = "pending"
    StatusClosed  Status = "closed"
)
```

**Verwenden Sie den Typ in einer Struktur oder einer Servicemethode:**

```go
type Ticket struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Status Status `json:"status"`
}
```

**Generieren Sie die Bindings:**

```bash
wails3 generate bindings
```

Die Ausgabe des Generators enthält neben der Anzahl der Modelle auch die Anzahl der Enums:

```
3 Enums, 5 Models
```

**Verwenden Sie das Enum in Ihrem Frontend:**

```javascript
import { Ticket, Status } from './bindings/changeme/models'

const ticket = new Ticket({
    id: 1,
    title: "Bug report",
    status: Status.StatusActive
})
```

**Das ist alles!** Der Enum-Typ wird sowohl in Go als auch in JavaScript/TypeScript durchgesetzt.

## Enums definieren

Ein Enum in Wails besteht aus einem **benannten Typ** mit einem zugrunde liegenden Basistyp und aus **const-Deklarationen** dieses Typs.

### String-Enums

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

**Generiertes TypeScript:**

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

**Generiertes JavaScript:**

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

### Integer-Enums

```go
type Priority int

const (
    PriorityLow    Priority = 0
    PriorityMedium Priority = 1
    PriorityHigh   Priority = 2
)
```

**Generiertes TypeScript:**

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

### Enums mit Typaliasen

Go-Typaliase (`=`) funktionieren ebenfalls, erzeugen jedoch eine etwas andere Ausgabe: eine Typdefinition zusammen mit einem const-Objekt statt eines nativen TypeScript-`enum`:

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

**Generiertes TypeScript:**

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

**Generiertes JavaScript:**

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
**Benannte Typen** (`type Title string`) erzeugen native TypeScript-`enum`-Deklarationen mit einem `$zero`-Member. **Typaliase** (`type Age = int`) erzeugen ein Namespace-Paar aus `type` + `const` ohne `$zero`.

@end

## Der `$zero`-Wert

Jedes Enum mit benanntem Typ enthält einen speziellen `$zero`-Member, der den **Go-Nullwert** des zugrunde liegenden Typs darstellt:

| Zugrunde liegender Typ | `$zero`-Wert |
| --- | --- |
| `string` | `""` |
| `int`, `int8`, `int16`, `int32`, `int64` | `0` |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | `0` |
| `float32`, `float64` | `0` |
| `bool` | `false` |

Wenn ein Strukturfeld einen Enum-Typ verwendet und kein Wert angegeben ist, verwendet der Konstruktor standardmäßig `$zero`:

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

Dies gewährleistet bei der Generierung von Klassen eine typsichere Initialisierung: Enum-Felder sind niemals `undefined`. Bei der Generierung von TypeScript-Schnittstellen (mit `-i`) gibt es keinen Konstruktor, und Felder dürfen wie üblich fehlen.

## Enums in Strukturen verwenden

Wenn ein Strukturfeld einen Enum-Typ hat, **behält der generierte Code diesen Typ bei**, statt auf den primitiven Typ zurückzufallen:

```go
type Person struct {
    Title Title
    Name  string
    Age   Age
}
```

**Generiertes TypeScript:**

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

Das Feld `Title` hat den Typ `Title` und nicht `string`. Dadurch bietet Ihre IDE vollständige Autovervollständigung und Typprüfung für Enum-Werte.

## Enums aus importierten Paketen

In separaten Paketen definierte Enums werden vollständig unterstützt. Sie werden im Verzeichnis des jeweiligen Pakets generiert:

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

Das Enum `Title` wird in der Modelldatei des Pakets `services` generiert, und Importpfade werden automatisch aufgelöst:

```typescript
// bindings/changeme/services/models.ts
export enum Title {
    $zero = "",
    Mister = "Mr",
    Miss = "Miss",
    Ms = "Ms",
}
```

## Enum-Methoden

Sie können Ihren Enum-Typen in Go Methoden hinzufügen. Diese wirken sich nicht auf die Generierung der Bindings aus, bieten aber nützliche serverseitige Funktionen:

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

Das generierte Enum ist unabhängig davon identisch, ob für den Typ Go-Methoden vorhanden sind.

## Kommentare und Dokumentation

Der Generator übernimmt Go-Kommentare als JSDoc in die generierte Ausgabe:

- **Typkommentare** werden zum Dokumentationskommentar des Enums
- **Kommentare zu const-Gruppen** werden zu Abschnittstrennern
- **Kommentare zu einzelnen Konstanten** werden zu Dokumentationskommentaren der Member
- **Inline-Kommentare** werden nach Möglichkeit beibehalten

Dadurch zeigt Ihre IDE beim Bewegen des Mauszeigers über Enum-Werte deren Dokumentation an.

## Unterstützte zugrunde liegende Typen

Der Binding-Generator unterstützt Enums mit den folgenden zugrunde liegenden Go-Typen:

| Go-Typ | Als Enum unterstützt |
| --- | :---: |
| `string` | Ja |
| `int`, `int8`, `int16`, `int32`, `int64` | Ja |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | Ja |
| `float32`, `float64` | Ja |
| `byte` (`uint8`) | Ja |
| `rune` (`int32`) | Ja |
| `bool` | Ja |
| `complex64`, `complex128` | Nein |

## Einschränkungen

Folgendes wird bei der Enum-Generierung **nicht** unterstützt:

- **Generische Typen** — Typparameter verhindern die Erkennung von Konstanten
- **Typen mit benutzerdefiniertem `json.Marshaler` oder `encoding.TextMarshaler`** — Bei benutzerdefinierter Serialisierung stimmen die generierten Werte möglicherweise nicht mit dem Laufzeitverhalten überein. Daher überspringt der Generator diese Typen
- **Konstanten, deren Werte nicht statisch ausgewertet oder dargestellt werden können** — Konstanten müssen bekannte Werte besitzen, die sich in ihrem zugrunde liegenden Typ darstellen lassen. Übliche `iota`-Muster funktionieren problemlos, da der Compiler sie in konkrete Werte auflöst
- **Komplexe Zahlentypen** — `complex64` und `complex128` können nicht als zugrunde liegende Enum-Typen verwendet werden

## Vollständiges Beispiel

**Go:**

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

**Frontend (TypeScript):**

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

**Frontend (JavaScript):**

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

## Nächste Schritte

@cards{cols="2"}
📖 Datenmodelle
Structs, Typzuordnung und Modellgenerierung.

[Mehr erfahren →](/features/bindings/models/)

---
🚀 Methodenbindung
Go-Methoden an das Frontend binden.

[Mehr erfahren →](/features/bindings/methods/)

---
⚙ Erweiterte Bindung
Direktiven, Code-Injektion und benutzerdefinierte IDs.

[Mehr erfahren →](/features/bindings/advanced/)

---
✓ Bewährte Vorgehensweisen
Entwurfsmuster für Bindings.

[Mehr erfahren →](/features/bindings/best-practices/)

@end

---

**Fragen?** Stelle sie auf [Discord](https://discord.gg/JDdSxwjhGf) oder sieh dir die [Binding-Beispiele](https://github.com/wailsapp/wails/tree/master/v3/examples/binding) an.
