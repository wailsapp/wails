---
title: "Перечисления"
description: "Автоматическое создание перечислений из констант Go"
slug: "features/bindings/enums"
sourcePath: "features/bindings/enums.md"
---

## Привязки перечислений

Генератор привязок Wails v3 **автоматически обнаруживает типы констант Go и создаёт перечисления TypeScript или объекты const JavaScript**. Не требуется ни регистрация, ни настройка: просто определите типы и константы в Go, а генератор сделает всё остальное.

@note{type="info"}
В отличие от Wails v2, **не нужно вызывать `EnumBind`** или регистрировать перечисления вручную. Генератор автоматически обнаруживает их в исходном коде.

@end

## Быстрый старт

**Определите в Go именованный тип с константами:**

```go
type Status string

const (
    StatusActive  Status = "active"
    StatusPending Status = "pending"
    StatusClosed  Status = "closed"
)
```

**Используйте этот тип в структуре или методе сервиса:**

```go
type Ticket struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Status Status `json:"status"`
}
```

**Создайте привязки:**

```bash
wails3 generate bindings
```

В выводе генератора количество перечислений будет указано рядом с количеством моделей:

```
3 Enums, 5 Models
```

**Используйте во фронтенде:**

```javascript
import { Ticket, Status } from './bindings/changeme/models'

const ticket = new Ticket({
    id: 1,
    title: "Bug report",
    status: Status.StatusActive
})
```

**Вот и всё!** Тип перечисления проверяется как в Go, так и в JavaScript/TypeScript.

## Определение перечислений

Перечисление в Wails — это **именованный тип** с базовым простым типом в сочетании с **объявлениями const** этого типа.

### Строковые перечисления

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

**Созданный код TypeScript:**

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

**Созданный код JavaScript:**

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

### Целочисленные перечисления

```go
type Priority int

const (
    PriorityLow    Priority = 0
    PriorityMedium Priority = 1
    PriorityHigh   Priority = 2
)
```

**Созданный код TypeScript:**

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

### Перечисления на основе псевдонимов типов

Псевдонимы типов Go (`=`) также поддерживаются, но для них создаётся немного другой код: определение типа и объект const вместо собственного `enum` TypeScript:

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

**Созданный код TypeScript:**

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

**Созданный код JavaScript:**

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
Для **именованных типов** (`type Title string`) создаются собственные объявления `enum` TypeScript с элементом `$zero`. Для **псевдонимов типов** (`type Age = int`) создаётся пара объявлений `type` + `const` в пространстве имён без `$zero`.

@end

## Значение `$zero`

Каждое перечисление на основе именованного типа содержит специальный элемент `$zero`, представляющий **нулевое значение Go** для базового типа:

| Базовый тип | Значение `$zero` |
| --- | --- |
| `string` | `""` |
| `int`, `int8`, `int16`, `int32`, `int64` | `0` |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | `0` |
| `float32`, `float64` | `0` |
| `bool` | `false` |

Если поле структуры имеет тип перечисления и значение не задано, конструктор по умолчанию использует `$zero`:

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

Это обеспечивает типобезопасную инициализацию при создании классов: поля перечислений никогда не имеют значения `undefined`. При создании интерфейсов TypeScript (с использованием `-i`) конструктора нет, поэтому поля, как обычно, могут отсутствовать.

## Использование перечислений в структурах

Если поле структуры имеет тип перечисления, созданный код **сохраняет этот тип**, а не заменяет его простым типом:

```go
type Person struct {
    Title Title
    Name  string
    Age   Age
}
```

**Созданный код TypeScript:**

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

Поле `Title` имеет тип `Title`, а не `string`. Благодаря этому IDE обеспечивает полноценное автодополнение и проверку типов для значений перечисления.

## Перечисления из импортированных пакетов

Перечисления, определённые в отдельных пакетах, полностью поддерживаются. Код для них создаётся в каталоге соответствующего пакета:

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

Перечисление `Title` создаётся в файле моделей `services`, а пути импорта определяются автоматически:

```typescript
// bindings/changeme/services/models.ts
export enum Title {
    $zero = "",
    Mister = "Mr",
    Miss = "Miss",
    Ms = "Ms",
}
```

## Методы перечислений

К типам перечислений в Go можно добавлять методы. Они не влияют на создание привязок, но предоставляют полезные возможности на стороне сервера:

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

Созданное перечисление будет одинаковым независимо от наличия у типа методов Go.

## Комментарии и документация

Генератор сохраняет комментарии Go в созданном коде в формате JSDoc:

- **Комментарии к типам** становятся документирующими комментариями к перечислениям
- **Комментарии к группам констант** становятся разделителями секций
- **Комментарии к отдельным константам** становятся документирующими комментариями к элементам
- **Встроенные комментарии** сохраняются, где это возможно

Благодаря этому IDE показывает документацию для значений перечисления при наведении указателя.

## Поддерживаемые базовые типы

Генератор привязок поддерживает перечисления со следующими базовыми типами Go:

| Тип Go | Поддерживается как перечисление |
| --- | :---: |
| `string` | Да |
| `int`, `int8`, `int16`, `int32`, `int64` | Да |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | Да |
| `float32`, `float64` | Да |
| `byte` (`uint8`) | Да |
| `rune` (`int32`) | Да |
| `bool` | Да |
| `complex64`, `complex128` | Нет |

## Ограничения

При создании перечислений **не** поддерживаются:

- **Обобщённые типы** — параметры типов препятствуют обнаружению констант
- **Типы с пользовательскими `json.Marshaler` или `encoding.TextMarshaler`** — при пользовательской сериализации сгенерированные значения могут не соответствовать поведению во время выполнения, поэтому генератор пропускает такие типы
- **Константы, значения которых невозможно вычислить статически или представить** — значения констант должны быть известны и представимы в их базовом типе. Стандартные шаблоны `iota` работают корректно, поскольку компилятор преобразует их в конкретные значения
- **Типы комплексных чисел** — `complex64` и `complex128` не могут быть базовыми типами перечислений

## Полный пример

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

**Фронтенд (TypeScript):**

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

**Фронтенд (JavaScript):**

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

## Дальнейшие шаги

@cards{cols="2"}
📖 Модели данных
Структуры, сопоставление типов и создание моделей.

[Подробнее →](/features/bindings/models/)

---
🚀 Привязка методов
Привязка методов Go к фронтенду.

[Подробнее →](/features/bindings/methods/)

---
⚙ Расширенная привязка
Директивы, внедрение кода и пользовательские идентификаторы.

[Подробнее →](/features/bindings/advanced/)

---
✓ Рекомендации
Шаблоны проектирования привязок.

[Подробнее →](/features/bindings/best-practices/)

@end

---

**Есть вопросы?** Задайте их в [Discord](https://discord.gg/JDdSxwjhGf) или ознакомьтесь с [примерами привязки](https://github.com/wailsapp/wails/tree/master/v3/examples/binding).
