---
title: "Мост между Go и фронтендом"
description: "Подробный разбор того, как Wails обеспечивает прямое взаимодействие между Go и JavaScript"
slug: "concepts/bridge"
sourcePath: "concepts/bridge.md"
---

## Прямое взаимодействие между Go и JavaScript

Wails предоставляет **прямой мост в памяти** между Go и JavaScript, обеспечивая беспрепятственное взаимодействие без накладных расходов HTTP, границ между процессами и узких мест сериализации.

## Общая картина

```d2
direction: right

Frontend: Фронтенд (JavaScript) {
  UI: React/Vue/Vanilla {
    shape: rectangle
    style.fill: "#8B5CF6"
  }

  Bindings: Автоматически сгенерированные привязки {
    shape: rectangle
    style.fill: "#A78BFA"
  }
}

Bridge: Мост Wails {
  Encoder: Кодировщик JSON {
    shape: rectangle
    style.fill: "#10B981"
  }

  Router: Маршрутизатор методов {
    shape: diamond
    style.fill: "#10B981"
  }

  Decoder: Декодировщик JSON {
    shape: rectangle
    style.fill: "#10B981"
  }

  TypeGen: Генератор типов {
    shape: rectangle
    style.fill: "#10B981"
  }
}

Backend: Бэкенд (Go) {
  Services: Ваши сервисы {
    shape: rectangle
    style.fill: "#00ADD8"
  }

  Registry: Реестр сервисов {
    shape: rectangle
    style.fill: "#00ADD8"
  }
}

Frontend.UI -> Frontend.Bindings: "import { Method }"
Frontend.Bindings -> Bridge.Encoder: "Вызов Method('arg')"
Bridge.Encoder -> Bridge.Router: Кодирование в JSON
Bridge.Router -> Backend.Registry: Поиск сервиса
Backend.Registry -> Backend.Services: Вызов метода
Backend.Services -> Bridge.Decoder: Возврат результата
Bridge.Decoder -> Frontend.Bindings: Декодирование в JS
Frontend.Bindings -> Frontend.UI: Разрешение Promise
Bridge.TypeGen -> Frontend.Bindings: Генерация типов
```

**Главное:** никаких HTTP, IPC и границ между процессами. Только **прямые вызовы функций** с **типобезопасностью**.

## Как это работает: пошаговое описание

### 1. Регистрация сервисов (запуск)

При запуске приложения Wails сканирует ваши сервисы:

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

**Что делает Wails:**

1. **Сканирует структуру** на наличие экспортируемых методов
2. **Извлекает сведения о типах** (параметрах и возвращаемых типах)
3. **Создаёт реестр**, сопоставляющий имена методов с функциями
4. **Создаёт привязки TypeScript** с полными определениями типов

### 2. Создание привязок (этап сборки)

Wails автоматически создаёт привязки TypeScript:

```typescript
// Auto-generated: frontend/bindings/GreetService.ts
export function Greet(name: string): Promise<string>
export function Add(a: number, b: number): Promise<number>
```

**Сопоставление типов:**

| Тип Go | Тип TypeScript |
| --- | --- |
| `string` | `string` |
| `int`, `int32`, `int64` | `number` |
| `float32`, `float64` | `number` |
| `bool` | `boolean` |
| `[]T` | `T[]` |
| `map[string]T` | `Record<string, T>` |
| `struct` | `interface` |
| `time.Time` | `Date` |
| `error` | Исключение (выбрасывается) |

### 3. Вызов из фронтенда (выполнение)

Разработчик вызывает метод Go из JavaScript:

```javascript
import { Greet, Add } from './bindings/GreetService'

// Call Go from JavaScript
const greeting = await Greet("World")
console.log(greeting)  // "Hello, World!"

const sum = await Add(5, 3)
console.log(sum)  // 8
```

**Что происходит:**

1. **Вызывается функция привязки** — `Greet("World")`
2. **Создаётся сообщение** — `{ service: "GreetService", method: "Greet", args: ["World"] }`
3. **Отправка в мост** — через JavaScript-мост WebView
4. **Возврат Promise** — ожидает ответа

### 4. Обработка в мосте (во время выполнения)

Мост получает и обрабатывает сообщение:

```d2
direction: down

Receive: Получение сообщения {
  shape: rectangle
  style.fill: "#10B981"
}

Parse: Разбор JSON {
  shape: rectangle
}

Validate: Проверка {
  Check: Сервис существует? {
    shape: diamond
  }

  CheckMethod: Метод существует? {
    shape: diamond
  }

  CheckTypes: Типы корректны? {
    shape: diamond
  }
}

Invoke: Вызов метода Go {
  shape: rectangle
  style.fill: "#00ADD8"
}

Encode: Кодирование результата {
  shape: rectangle
}

Send: Отправка ответа {
  shape: rectangle
  style.fill: "#10B981"
}

Error: Отправка ошибки {
  shape: rectangle
  style.fill: "#EF4444"
}

Receive -> Parse
Parse -> Validate.Check
Validate.Check -> Validate.CheckMethod: Да
Validate.Check -> Error: Нет
Validate.CheckMethod -> Validate.CheckTypes: Да
Validate.CheckMethod -> Error: Нет
Validate.CheckTypes -> Invoke: Да
Validate.CheckTypes -> Error: Нет
Invoke -> Encode: Успешно
Invoke -> Error: Ошибка
Encode -> Send
```

**Безопасность:** можно вызывать только зарегистрированные сервисы и экспортируемые методы.

### 5. Выполнение Go-кода (во время выполнения)

Выполняется метод Go:

```go
func (g *GreetService) Greet(name string) string {
    // This runs in Go
    return g.prefix + name + "!"
}
```

**Контекст выполнения:**

- Выполняется в **горутине** (без блокировки)
- Имеет доступ ко **всем возможностям Go** (файловой системе, сети, базам данных)
- Может свободно вызывать **другой Go-код**
- Возвращает результат или ошибку

### 6. Ответ (во время выполнения)

Результат отправляется обратно в JavaScript:

```javascript
// Promise resolves with result
const greeting = await Greet("World")
// greeting = "Hello, World!"
```

**Обработка ошибок:**

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

## Характеристики производительности

### Скорость

**Типичные накладные расходы на вызов:** &lt;1 мс

```
Frontend Call → Bridge → Go Execution → Bridge → Frontend Response
     ↓            ↓           ↓            ↓            ↓
   &lt;0.1ms      &lt;0.1ms      [varies]     &lt;0.1ms      &lt;0.1ms
```

**Сравнение с альтернативами:**

- **HTTP/REST:** 5-50 мс (сетевой стек, сериализация)
- **IPC:** 1-10 мс (границы процессов, маршалинг)
- **Мост Wails:** &lt;1 мс (в памяти, прямой вызов)

### Память

**Накладные расходы на один вызов:** ~1 КБ (буфер сообщения)

**Оптимизация без копирования:** для больших объёмов данных (>1 МБ) по возможности используется общая память.

### Конкурентное выполнение

**Вызовы выполняются конкурентно:**

- Каждый вызов выполняется в отдельной горутине
- Несколько вызовов могут выполняться одновременно
- Вызовы не блокируют друг друга

```javascript
// These run concurrently
const [result1, result2, result3] = await Promise.all([
    SlowOperation1(),
    SlowOperation2(),
    SlowOperation3(),
])
```

## Система типов

### Поддерживаемые типы

#### Примитивные типы

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

#### Срезы и массивы

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

#### Отображения

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

#### Структуры

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

**Теги JSON:** используйте теги `json:`, чтобы задавать имена полей в TypeScript.

#### Время

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

#### Ошибки

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

### Неподдерживаемые типы

Эти типы **нельзя** передавать через мост:

- **Каналы** (`chan T`)
- **Функции** (`func()`)
- **Интерфейсы** (кроме `interface{}` / `any`)
- **Указатели** (кроме указателей на структуры)
- **Неэкспортируемые поля** (со строчной буквы)

**Обходное решение:** используйте идентификаторы или дескрипторы:

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

## Расширенные приёмы

### Передача контекста

Сервисы могут обращаться к контексту вызова:

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

**Контекст предоставляет:**

- Окно, из которого выполнен вызов
- Экземпляр приложения
- Метаданные запроса

### Потоковая передача данных

Для больших объёмов данных используйте события вместо возвращаемых значений:

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

### Отмена операций

Для операций с возможностью отмены используйте контекст:

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

**Примечание:** контекст автоматически отменяется при отключении фронтенда.

### Пакетные операции

Сократите накладные расходы моста, объединяя операции в пакеты:

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

## Отладка моста

### Включение отладочного журналирования

```go
app := application.New(application.Options{
    Name:     "My App",
    LogLevel: slog.LevelDebug, // requires `import "log/slog"`
    // `Logger` is an optional *slog.Logger; the default-logger helper is
    // application.DefaultLogger(slog.Leveler) if you want to construct one explicitly.
})
```

**В выводе отображаются:**

- Вызовы методов
- Параметры
- Возвращаемые значения
- Ошибки
- Информация о времени выполнения

### Проверка сгенерированных привязок

Чтобы просмотреть сгенерированный TypeScript, проверьте `frontend/bindings/`:

```javascript
// frontend/bindings/<full-go-import-path>/myservice.js (real generated shape)
import { Call as $Call } from "/wails/runtime.js";

export function MyMethod($0) {
    return $Call.ByID(1234567890, $0); // numeric method ID assigned by the generator
}
```

### Непосредственное тестирование сервисов

Тестируйте сервисы Go без фронтенда:

```go
func TestGreetService(t *testing.T) {
    service := &GreetService{prefix: "Hello, "}
    result := service.Greet("Test")
    if result != "Hello, Test!" {
        t.Errorf("Expected 'Hello, Test!', got '%s'", result)
    }
}
```

## Советы по производительности

### ✅ Рекомендуется

- **Объединяйте операции в пакеты** — сокращайте количество вызовов моста
- **Используйте события для потоковой передачи** — не возвращайте большие массивы
- **Обеспечивайте быстрое выполнение методов** — в идеале &lt;100 мс
- **Используйте горутины** — для длительных операций
- **Кешируйте данные на стороне Go** — избегайте повторных вычислений

### ❌ Не рекомендуется

- **Не выполняйте избыточные вызовы** — по возможности объединяйте их в пакеты
- **Не возвращайте огромные объёмы данных** — используйте пагинацию или потоковую передачу
- **Не блокируйте выполнение** — используйте горутины для длительных операций
- **Не передавайте сложные типы** — используйте простые типы
- **Не игнорируйте ошибки** — всегда обрабатывайте их

## Безопасность

Мост по умолчанию безопасен:

1. **Только разрешённые сервисы** — вызывать можно лишь зарегистрированные сервисы
2. **Проверка типов** — аргументы проверяются на соответствие типам Go
3. **Без eval()** — фронтенд не может выполнять произвольный код Go
4. **Без злоупотребления рефлексией** — доступны только экспортируемые методы

**Рекомендации:**

- **Проверяйте входные данные** в Go (не доверяйте фронтенду)
- **Используйте контекст** для аутентификации и авторизации
- **Ограничивайте частоту выполнения** ресурсоёмких операций
- **Очищайте** пути к файлам и пользовательский ввод

## Дальнейшие действия

**Система сборки** — узнайте, как Wails собирает ваше приложение и объединяет его ресурсы в пакет [Подробнее →](/concepts/build-system/)

**Сервисы** — подробно изучите систему сервисов [Подробнее →](/features/bindings/services/)

**События** — используйте события для обмена данными по модели публикации и подписки [Подробнее →](/features/events/system/)

---

**Есть вопросы о мосте?** Задайте их в [Discord](https://discord.gg/JDdSxwjhGf) или ознакомьтесь с [примерами привязок](https://github.com/wailsapp/wails/tree/master/v3/examples/binding).
