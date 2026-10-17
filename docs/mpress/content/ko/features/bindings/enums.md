---
title: "열거형"
description: "Go 상수에서 열거형 자동 생성"
slug: "features/bindings/enums"
sourcePath: "features/bindings/enums.md"
---

## 열거형 바인딩

Wails v3 바인딩 생성기는 **Go 상수 형식을 자동으로 감지하여 TypeScript 열거형 또는 JavaScript const 객체를 생성합니다**. 등록이나 설정은 필요하지 않습니다. Go에서 형식과 상수를 정의하기만 하면 나머지는 생성기가 처리합니다.

@note{type="info"}
Wails v2와 달리 <strong>`EnumBind`</strong>을 호출하거나 열거형을 수동으로 등록할 필요가 없습니다. 생성기가 소스 코드에서 열거형을 자동으로 검색합니다.

@end

## 빠른 시작

**Go에서 상수와 함께 명명된 형식을 정의합니다.**

```go
type Status string

const (
    StatusActive  Status = "active"
    StatusPending Status = "pending"
    StatusClosed  Status = "closed"
)
```

**구조체 또는 서비스 메서드에서 해당 형식을 사용합니다.**

```go
type Ticket struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Status Status `json:"status"`
}
```

**바인딩을 생성합니다.**

```bash
wails3 generate bindings
```

생성기 출력에는 모델 수와 함께 열거형 수도 표시됩니다.

```
3 Enums, 5 Models
```

**프런트엔드에서 사용합니다.**

```javascript
import { Ticket, Status } from './bindings/changeme/models'

const ticket = new Ticket({
    id: 1,
    title: "Bug report",
    status: Status.StatusActive
})
```

**이것으로 끝입니다!** 열거형 형식은 Go와 JavaScript/TypeScript 모두에서 강제됩니다.

## 열거형 정의

Wails의 열거형은 기본 형식을 기반 형식으로 하는 <strong>명명된 형식</strong>과 해당 형식의 <strong>const 선언</strong>을 결합한 것입니다.

### 문자열 열거형

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

**생성된 TypeScript:**

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

**생성된 JavaScript:**

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

### 정수 열거형

```go
type Priority int

const (
    PriorityLow    Priority = 0
    PriorityMedium Priority = 1
    PriorityHigh   Priority = 2
)
```

**생성된 TypeScript:**

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

### 형식 별칭 열거형

Go 형식 별칭(`=`)도 사용할 수 있지만, 출력은 약간 다릅니다. 네이티브 TypeScript `enum` 대신 형식 정의와 const 객체가 생성됩니다.

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

**생성된 TypeScript:**

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

**생성된 JavaScript:**

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
**명명된 형식**(`type Title string`)은 `$zero` 멤버가 포함된 네이티브 TypeScript `enum` 선언을 생성합니다. **형식 별칭**(`type Age = int`)은 `$zero` 없이 `type` + `const` 네임스페이스 쌍을 생성합니다.

@end

## `$zero` 값

모든 명명된 형식의 열거형에는 기반 형식의 <strong>Go 제로 값</strong>을 나타내는 특수 `$zero` 멤버가 포함됩니다.

| 기반 형식 | `$zero` 값 |
| --- | --- |
| `string` | `""` |
| `int`, `int8`, `int16`, `int32`, `int64` | `0` |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | `0` |
| `float32`, `float64` | `0` |
| `bool` | `false` |

구조체 필드가 열거형 형식을 사용하며 값이 제공되지 않으면 생성자는 기본값으로 `$zero`을 사용합니다.

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

따라서 클래스를 생성할 때 형식 안전성이 보장되도록 초기화되며, 열거형 필드는 절대 `undefined`이 되지 않습니다. TypeScript 인터페이스를 생성할 때(`-i` 사용)는 생성자가 없으며, 평소와 마찬가지로 필드가 없을 수 있습니다.

## 구조체에서 열거형 사용

구조체 필드가 열거형 형식이면 생성된 코드는 기본 형식으로 대체하지 않고 **해당 형식을 유지합니다**.

```go
type Person struct {
    Title Title
    Name  string
    Age   Age
}
```

**생성된 TypeScript:**

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

`Title` 필드의 형식은 `string`이 아니라 `Title`입니다. 따라서 IDE에서 열거형 값에 대한 완전한 자동 완성과 형식 검사를 사용할 수 있습니다.

## 가져온 패키지의 열거형

별도 패키지에 정의된 열거형도 완전히 지원됩니다. 열거형은 해당 패키지 디렉터리에 생성됩니다.

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

`Title` 열거형은 `services` 모델 파일에 생성되며, 가져오기 경로는 자동으로 확인됩니다.

```typescript
// bindings/changeme/services/models.ts
export enum Title {
    $zero = "",
    Mister = "Mr",
    Miss = "Miss",
    Ms = "Ms",
}
```

## 열거형 메서드

Go에서 열거형 형식에 메서드를 추가할 수 있습니다. 이러한 메서드는 바인딩 생성에 영향을 주지 않지만 유용한 서버 측 기능을 제공합니다.

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

Go 메서드가 해당 형식에 존재하는지와 관계없이 생성되는 열거형은 동일합니다.

## 주석과 문서

생성기는 Go 주석을 생성된 출력의 JSDoc으로 보존합니다.

- <strong>형식 주석</strong>은 열거형의 문서 주석이 됩니다.
- <strong>const 그룹 주석</strong>은 섹션 구분자가 됩니다.
- <strong>개별 const 주석</strong>은 멤버 문서 주석이 됩니다.
- <strong>인라인 주석</strong>은 가능한 경우 보존됩니다.

따라서 IDE에서 열거형 값 위에 마우스 포인터를 올리면 관련 문서가 표시됩니다.

## 지원되는 기반 형식

바인딩 생성기는 다음 Go 기반 형식의 열거형을 지원합니다.

| Go 형식 | 열거형으로 사용 가능 |
| --- | :---: |
| `string` | 예 |
| `int`, `int8`, `int16`, `int32`, `int64` | 예 |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | 예 |
| `float32`, `float64` | 예 |
| `byte`(`uint8`) | 예 |
| `rune`(`int32`) | 예 |
| `bool` | 예 |
| `complex64`, `complex128` | 아니요 |

## 제한 사항

다음 항목은 열거형 생성에서 **지원되지 않습니다**:

- **제네릭 타입** — 타입 매개변수로 인해 상수를 감지할 수 없습니다
- <strong>사용자 정의 `json.Marshaler` 또는 `encoding.TextMarshaler`</strong>가 있는 타입 — 사용자 정의 직렬화를 사용하면 생성된 값이 런타임 동작과 일치하지 않을 수 있으므로 생성기가 이러한 타입을 건너뜁니다
- **값을 정적으로 평가하거나 표현할 수 없는 상수** — 상수는 기반 타입에서 값이 명확하고 표현 가능해야 합니다. 컴파일러가 구체적인 값으로 해석하므로 표준 `iota` 패턴은 정상적으로 작동합니다
- **복소수 타입** — `complex64` 및 `complex128`는 열거형의 기반 타입으로 사용할 수 없습니다

## 전체 예제

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

**프런트엔드(TypeScript):**

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

**프런트엔드(JavaScript):**

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

## 다음 단계

@cards{cols="2"}
📖 데이터 모델
구조체, 타입 매핑 및 모델 생성.

[자세히 알아보기 →](/features/bindings/models/)

---
🚀 메서드 바인딩
Go 메서드를 프런트엔드에 바인딩합니다.

[자세히 알아보기 →](/features/bindings/methods/)

---
⚙ 고급 바인딩
지시문, 코드 삽입 및 사용자 정의 ID.

[자세히 알아보기 →](/features/bindings/advanced/)

---
✓ 모범 사례
바인딩 설계 패턴.

[자세히 알아보기 →](/features/bindings/best-practices/)

@end

---

**궁금한 점이 있으신가요?** [Discord](https://discord.gg/JDdSxwjhGf)에서 질문하거나 [바인딩 예제](https://github.com/wailsapp/wails/tree/master/v3/examples/binding)를 확인하세요.
