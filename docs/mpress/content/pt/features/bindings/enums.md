---
title: "Enums"
description: "Geração automática de enums a partir de constantes Go"
slug: "features/bindings/enums"
sourcePath: "features/bindings/enums.md"
---

## Bindings de enums

O gerador de bindings do Wails v3 **detecta automaticamente tipos de constantes Go e gera enums TypeScript ou objetos const JavaScript**. Não é necessário registrar nem configurar nada — basta definir seus tipos e constantes em Go, e o gerador cuida do restante.

@note{type="info"}
Ao contrário do Wails v2, **não é necessário chamar `EnumBind`** nem registrar enums manualmente. O gerador os descobre automaticamente no seu código-fonte.

@end

## Início rápido

**Defina um tipo nomeado com constantes em Go:**

```go
type Status string

const (
    StatusActive  Status = "active"
    StatusPending Status = "pending"
    StatusClosed  Status = "closed"
)
```

**Use o tipo em uma struct ou em um método de serviço:**

```go
type Ticket struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Status Status `json:"status"`
}
```

**Gere os bindings:**

```bash
wails3 generate bindings
```

A saída do gerador informará a quantidade de enums junto com a de modelos:

```
3 Enums, 5 Models
```

**Use no seu frontend:**

```javascript
import { Ticket, Status } from './bindings/changeme/models'

const ticket = new Ticket({
    id: 1,
    title: "Bug report",
    status: Status.StatusActive
})
```

**É só isso!** O tipo enum é verificado tanto em Go quanto em JavaScript/TypeScript.

## Definição de enums

No Wails, um enum é um **tipo nomeado** com um tipo básico subjacente, combinado com **declarações const** desse tipo.

### Enums de strings

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

**TypeScript gerado:**

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

**JavaScript gerado:**

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

### Enums de inteiros

```go
type Priority int

const (
    PriorityLow    Priority = 0
    PriorityMedium Priority = 1
    PriorityHigh   Priority = 2
)
```

**TypeScript gerado:**

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

### Enums de aliases de tipo

Aliases de tipo do Go (`=`) também funcionam, mas geram uma saída um pouco diferente — uma definição de tipo e um objeto const, em vez de um `enum` nativo do TypeScript:

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

**TypeScript gerado:**

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

**JavaScript gerado:**

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
**Tipos nomeados** (`type Title string`) geram declarações `enum` nativas do TypeScript com um membro `$zero`. **Aliases de tipo** (`type Age = int`) geram um par de namespaces `type` + `const` sem `$zero`.

@end

## O valor `$zero`

Todo enum de tipo nomeado inclui um membro especial `$zero` que representa o **valor zero do Go** para o tipo subjacente:

| Tipo subjacente | Valor `$zero` |
| --- | --- |
| `string` | `""` |
| `int`, `int8`, `int16`, `int32`, `int64` | `0` |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | `0` |
| `float32`, `float64` | `0` |
| `bool` | `false` |

Quando um campo de uma struct usa um tipo enum e nenhum valor é fornecido, o construtor usa `$zero` como padrão:

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

Isso garante uma inicialização com segurança de tipos ao gerar classes — campos enum nunca são `undefined`. Ao gerar interfaces TypeScript (usando `-i`), não há construtor, e os campos podem estar ausentes como de costume.

## Uso de enums em structs

Quando um campo de uma struct tem um tipo enum, o código gerado **preserva esse tipo** em vez de usar o tipo primitivo:

```go
type Person struct {
    Title Title
    Name  string
    Age   Age
}
```

**TypeScript gerado:**

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

O campo `Title` é tipado como `Title`, e não como `string`. Isso permite que sua IDE ofereça preenchimento automático completo e verificação de tipos para os valores do enum.

## Enums de pacotes importados

Enums definidos em pacotes separados são totalmente compatíveis. Eles são gerados no diretório do pacote correspondente:

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

O enum `Title` é gerado no arquivo de modelos `services`, e os caminhos de importação são resolvidos automaticamente:

```typescript
// bindings/changeme/services/models.ts
export enum Title {
    $zero = "",
    Mister = "Mr",
    Miss = "Miss",
    Ms = "Ms",
}
```

## Métodos de enums

Você pode adicionar métodos aos seus tipos enum em Go. Esses métodos não afetam a geração de bindings, mas oferecem funcionalidades úteis no lado do servidor:

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

O enum gerado é idêntico, independentemente de o tipo ter métodos Go.

## Comentários e documentação

O gerador preserva os comentários Go como JSDoc na saída gerada:

- **Comentários de tipo** tornam-se o comentário de documentação do enum
- **Comentários de grupos de constantes** tornam-se separadores de seção
- **Comentários de constantes individuais** tornam-se comentários de documentação dos membros
- **Comentários na mesma linha** são preservados quando possível

Isso significa que sua IDE exibirá a documentação dos valores do enum quando você passar o cursor sobre eles.

## Tipos subjacentes compatíveis

O gerador de bindings oferece suporte a enums com os seguintes tipos subjacentes do Go:

| Tipo Go | Funciona como enum |
| --- | :---: |
| `string` | Sim |
| `int`, `int8`, `int16`, `int32`, `int64` | Sim |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | Sim |
| `float32`, `float64` | Sim |
| `byte` (`uint8`) | Sim |
| `rune` (`int32`) | Sim |
| `bool` | Sim |
| `complex64`, `complex128` | Não |

## Limitações

Os itens a seguir **não** são compatíveis com a geração de enums:

- **Tipos genéricos** — Os parâmetros de tipo impedem a detecção de constantes
- **Tipos com `json.Marshaler` ou `encoding.TextMarshaler`** personalizados — A serialização personalizada significa que os valores gerados podem não corresponder ao comportamento em tempo de execução; por isso, o gerador ignora esses tipos
- **Constantes cujos valores não podem ser avaliados estaticamente ou representados** — As constantes devem ter valores conhecidos e representáveis em seu tipo subjacente. Os padrões `iota` padrão funcionam normalmente, pois o compilador os resolve em valores concretos
- **Tipos de números complexos** — `complex64` e `complex128` não podem ser tipos subjacentes de enums

## Exemplo completo

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

## Próximas etapas

@cards{cols="2"}
📖 Modelos de dados
Structs, mapeamento de tipos e geração de modelos.

[Saiba mais →](/features/bindings/models/)

---
🚀 Vinculação de métodos
Vincule métodos Go ao frontend.

[Saiba mais →](/features/bindings/methods/)

---
⚙ Vinculação avançada
Diretivas, injeção de código e IDs personalizados.

[Saiba mais →](/features/bindings/advanced/)

---
✓ Práticas recomendadas
Padrões de projeto para vinculação.

[Saiba mais →](/features/bindings/best-practices/)

@end

---

**Dúvidas?** Pergunte no [Discord](https://discord.gg/JDdSxwjhGf) ou consulte os [exemplos de vinculação](https://github.com/wailsapp/wails/tree/master/v3/examples/binding).
