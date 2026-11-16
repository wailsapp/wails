---
title: "고급 바인딩"
description: "지시문, 코드 삽입, 사용자 지정 ID를 포함한 고급 바인딩 기법"
slug: "features/bindings/advanced"
sourcePath: "features/bindings/advanced.md"
---

이 가이드에서는 Wails v3의 바인딩 생성 프로세스를 사용자 지정하고 최적화하는 고급 기법을 설명합니다.

## 지시문으로 생성 코드 사용자 지정하기

### 사용자 지정 코드 삽입하기

`//wails:inject` 지시문을 사용하면 생성된 바인딩에 사용자 지정 JavaScript/TypeScript 코드를 삽입할 수 있습니다.

```go
//wails:inject console.log("Hello from Wails!");
type MyService struct {}

func (s *MyService) Greet(name string) string {
    return "Hello, " + name
}
```

이렇게 하면 지정한 코드가 `MyService` 서비스용으로 생성된 JavaScript/TypeScript 파일에 삽입됩니다.

조건부 삽입을 사용하여 특정 출력 형식만 대상으로 지정할 수도 있습니다.

```go
//wails:inject j*:console.log("Hello JS!");  // JavaScript only
//wails:inject t*:console.log("Hello TS!");  // TypeScript only
```

### 추가 파일 포함하기

`//wails:include` 지시문을 사용하면 생성된 바인딩에 추가 파일을 포함할 수 있습니다.

```go
//wails:include js/*.js
package mypackage
```

일반적으로 패키지 문서 주석에서 이 지시문을 사용하여 생성된 바인딩에 추가 JavaScript/TypeScript 파일을 포함합니다.

### 내부 형식 및 메서드 표시하기

`//wails:internal` 지시문은 생성된 JavaScript/TypeScript 바인딩에서 형식이나 메서드가 내보내지지 않도록 합니다.

```go
//wails:internal
type InternalModel struct {
    Field string
}

//wails:internal
func (s *MyService) InternalMethod() {}
```

이 지시문은 런타임 접근이 아니라 생성된 SDK를 제어합니다. 일반 빌드에서는 내보낸 Go 서비스 메서드가 생성된 SDK에 없더라도 여전히 런타임을 통해 등록하고 호출할 수 있습니다. 비공개 작업은 내보내지 않거나 바인딩되지 않은 형식에 두고, 노출된 메서드 내부에서 권한 부여를 적용하세요.

### 메서드 무시하기

`//wails:ignore` 지시문은 바인딩을 생성할 때 메서드를 완전히 무시합니다.

```go
//wails:ignore
func (s *MyService) IgnoredMethod() {}
```

이렇게 하면 생성되는 래퍼가 생략됩니다. `//wails:internal`과 마찬가지로, 이는 내보낸 Go 서비스 메서드에 대한 접근 제어 경계가 아닙니다.

### 사용자 지정 메서드 ID

`//wails:id` 지시문은 기본 해시 기반 ID를 재정의하여 메서드의 사용자 지정 ID를 지정합니다.

```go
//wails:id 42
func (s *MyService) CustomIDMethod() {}
```

이는 코드를 리팩터링할 때 호환성을 유지하는 데 유용할 수 있습니다.

## 복합 형식 다루기

### 중첩 구조체

바인딩 생성기는 중첩 구조체를 자동으로 처리합니다.

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

생성된 JavaScript/TypeScript 코드에는 `Person` 및 `Address`용 클래스가 모두 포함됩니다.

### 맵과 슬라이스

맵과 슬라이스도 자동으로 처리됩니다.

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

JavaScript에서 맵은 객체로, 슬라이스는 배열로 표현됩니다. TypeScript에서 맵은 `Record<K, V>`로, 슬라이스는 `T[]`로 표현됩니다.

### 제네릭 형식

바인딩 생성기는 제네릭 형식을 지원합니다.

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

생성된 TypeScript 코드에는 `Result`용 제네릭 클래스가 포함됩니다.

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

### 인터페이스

`-i` 플래그를 사용하면 바인딩 생성기가 클래스 대신 TypeScript 인터페이스를 생성하도록 할 수 있습니다.

```bash
wails3 generate bindings -ts -i
```

이렇게 하면 모든 모델에 대해 TypeScript 인터페이스가 생성됩니다.

```typescript
export interface Person {
    Name: string;
    Attributes: Record<string, string>;
    Friends: string[];
}
```

## 바인딩 생성 최적화하기

### ID 대신 이름 사용하기

기본적으로 바인딩 생성기는 메서드 호출에 해시 기반 ID를 사용합니다. 대신 이름을 사용하려면 `-names` 플래그를 사용할 수 있습니다.

```bash
wails3 generate bindings -names
```

이렇게 하면 바인딩된 메서드를 **정규화된 전체** 이름(`<package>.<Service>.<Method>`)으로 호출하고 위치 기반 매개변수 이름 `$0`, `$1`, …을 사용하는 코드가 생성됩니다.

```javascript
export function Greet($0) {
    let $resultPromise = $Call.ByName("main.GreetService.Greet", $0);
    return $resultPromise;
}
```

이렇게 하면 생성된 코드를 더 쉽게 읽고 디버깅할 수 있지만 효율성이 약간 떨어질 수 있습니다.

### 런타임 번들링하기

기본적으로 생성된 코드는 `@wailsio/runtime` npm 패키지에서 Wails 런타임을 가져옵니다. `-b` 플래그를 사용하면 런타임을 생성된 코드와 함께 번들링할 수 있습니다.

```bash
wails3 generate bindings -b
```

이렇게 하면 런타임 코드가 생성된 파일에 직접 포함되므로 npm 패키지가 필요하지 않습니다.

### 인덱스 파일 비활성화하기

인덱스 파일이 필요하지 않다면 `-noindex` 플래그를 사용하여 생성을 비활성화할 수 있습니다.

```bash
wails3 generate bindings -noindex
```

서비스와 모델을 각각의 파일에서 직접 가져오려는 경우 유용할 수 있습니다.

## 실제 사용 예제

### 인증 서비스

다음은 사용자 지정 지시문을 사용하는 인증 서비스의 예입니다.

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

### 데이터 처리 서비스

다음은 제네릭 형식을 사용하는 데이터 처리 서비스의 예입니다.

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

### 조건부 코드 삽입

다음은 출력 형식에 따라 코드를 조건부로 삽입하는 예입니다.

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

이렇게 하면 JavaScript 및 TypeScript 출력에 서로 다른 코드가 삽입되어 각 언어에 적합한 형식 주석이 제공됩니다.
