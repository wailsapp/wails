---
title: "Расширенные возможности привязки"
description: "Расширенные методы привязки, включая директивы, внедрение кода и пользовательские идентификаторы"
slug: "features/bindings/advanced"
sourcePath: "features/bindings/advanced.md"
---

В этом руководстве рассматриваются расширенные методы настройки и оптимизации процесса генерации привязок в Wails v3.

## Настройка генерируемого кода с помощью директив

### Внедрение пользовательского кода

Директива `//wails:inject` позволяет внедрять пользовательский код JavaScript/TypeScript в генерируемые привязки:

```go
//wails:inject console.log("Hello from Wails!");
type MyService struct {}

func (s *MyService) Greet(name string) string {
    return "Hello, " + name
}
```

Указанный код будет внедрён в сгенерированный файл JavaScript/TypeScript для сервиса `MyService`.

Для отдельных форматов выходных данных также можно использовать условное внедрение:

```go
//wails:inject j*:console.log("Hello JS!");  // JavaScript only
//wails:inject t*:console.log("Hello TS!");  // TypeScript only
```

### Включение дополнительных файлов

Директива `//wails:include` позволяет включать дополнительные файлы в генерируемые привязки:

```go
//wails:include js/*.js
package mypackage
```

Эта директива обычно используется в комментариях к документации пакета, чтобы включать дополнительные файлы JavaScript/TypeScript в генерируемые привязки.

### Пометка внутренних типов и методов

Директива `//wails:internal` запрещает экспорт типа или метода в генерируемых привязках JavaScript/TypeScript:

```go
//wails:internal
type InternalModel struct {
    Field string
}

//wails:internal
func (s *MyService) InternalMethod() {}
```

Эта директива управляет генерируемым SDK, а не доступом во время выполнения. В обычных сборках экспортируемые методы сервисов Go по-прежнему можно зарегистрировать и вызвать через среду выполнения, даже если они отсутствуют в сгенерированном SDK. Не экспортируйте закрытые операции либо размещайте их в непривязанном типе, а в открытых методах выполняйте проверку авторизации.

### Игнорирование методов

Директива `//wails:ignore` полностью исключает метод из процесса генерации привязок:

```go
//wails:ignore
func (s *MyService) IgnoredMethod() {}
```

При этом генерируемая обёртка не создаётся. Как и `//wails:internal`, эта директива не является границей контроля доступа для экспортируемых методов сервисов Go.

### Пользовательские идентификаторы методов

Директива `//wails:id` задаёт пользовательский идентификатор метода вместо используемого по умолчанию идентификатора на основе хеша:

```go
//wails:id 42
func (s *MyService) CustomIDMethod() {}
```

Это может быть полезно для сохранения совместимости при рефакторинге кода.

## Работа со сложными типами

### Вложенные структуры

Генератор привязок автоматически обрабатывает вложенные структуры:

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

Сгенерированный код JavaScript/TypeScript будет содержать классы как для `Person`, так и для `Address`.

### Отображения и срезы

Отображения и срезы также обрабатываются автоматически:

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

В JavaScript отображения представлены объектами, а срезы — массивами. В TypeScript отображения представлены как `Record<K, V>`, а срезы — как `T[]`.

### Обобщённые типы

Генератор привязок поддерживает обобщённые типы:

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

Сгенерированный код TypeScript будет содержать обобщённый класс для `Result`:

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

### Интерфейсы

С помощью флага `-i` генератор привязок может создавать интерфейсы TypeScript вместо классов:

```bash
wails3 generate bindings -ts -i
```

В результате для всех моделей будут созданы интерфейсы TypeScript:

```typescript
export interface Person {
    Name: string;
    Attributes: Record<string, string>;
    Friends: string[];
}
```

## Оптимизация генерации привязок

### Использование имён вместо идентификаторов

По умолчанию генератор привязок использует для вызовов методов идентификаторы на основе хеша. Чтобы вместо них использовать имена, укажите флаг `-names`:

```bash
wails3 generate bindings -names
```

В результате будет создан код, который вызывает привязанный метод по его **полному** имени (`<package>.<Service>.<Method>`) и использует позиционные имена параметров `$0`, `$1`, …:

```javascript
export function Greet($0) {
    let $resultPromise = $Call.ByName("main.GreetService.Greet", $0);
    return $resultPromise;
}
```

Это может сделать сгенерированный код более понятным и упростить его отладку, однако он может работать немного менее эффективно.

### Включение среды выполнения в пакет

По умолчанию сгенерированный код импортирует среду выполнения Wails из npm-пакета `@wailsio/runtime`. С помощью флага `-b` можно включить среду выполнения непосредственно в генерируемый код:

```bash
wails3 generate bindings -b
```

Код среды выполнения будет включён непосредственно в генерируемые файлы, поэтому npm-пакет больше не потребуется.

### Отключение индексных файлов

Если индексные файлы не нужны, отключите их генерацию с помощью флага `-noindex`:

```bash
wails3 generate bindings -noindex
```

Это может быть полезно, если вы предпочитаете импортировать сервисы и модели непосредственно из соответствующих файлов.

## Практические примеры

### Сервис аутентификации

Ниже приведён пример сервиса аутентификации с пользовательскими директивами:

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

### Сервис обработки данных

Ниже приведён пример сервиса обработки данных с обобщёнными типами:

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

### Условное внедрение кода

Ниже приведён пример условного внедрения кода для разных форматов выходных данных:

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

Для выходных данных JavaScript и TypeScript внедряется разный код с подходящими аннотациями типов для каждого языка.
