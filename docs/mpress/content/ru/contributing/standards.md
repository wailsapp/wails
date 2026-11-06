---
title: "Стандарты написания кода"
description: "Стиль кода, соглашения и рекомендации для Wails v3"
slug: "contributing/standards"
sourcePath: "contributing/standards.md"
---

## Стиль кода и соглашения

Соблюдение единых стандартов написания кода упрощает чтение и сопровождение кодовой базы, а также внесение в неё изменений.

## Стандарты кода Go

### Форматирование кода

Используйте стандартные инструменты форматирования Go:

```bash
# Format all code
gofmt -w .

# Use goimports for import organization
goimports -w .
```

**Обязательное требование:** перед фиксацией изменений весь код Go должен успешно проходить проверки `gofmt` и `goimports`.

### Соглашения об именовании

**Пакеты:**

- Используйте нижний регистр и по возможности одно слово
- `package application`, `package events`
- Не используйте символы подчёркивания и сочетания разных регистров

**Экспортируемые имена:**

- Используйте PascalCase для типов, функций и констант
- `type WebviewWindow struct`, `func NewApplication()`

**Неэкспортируемые имена:**

- Используйте camelCase для внутренних типов, функций и переменных
- `type windowImpl struct`, `func createWindow()`

**Интерфейсы:**

- Называйте по поведению: `Reader`, `Writer`, `Handler`
- Для интерфейсов с одним методом используйте имена с суффиксом `-er`

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

### Обработка ошибок

**Всегда проверяйте ошибки:**

```go
// Good
result, err := doSomething()
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// Bad - ignoring errors
result, _ := doSomething()
```

**Используйте обёртывание ошибок:**

```go
// Wrap errors to provide context
if err := validate(); err != nil {
    return fmt.Errorf("validation failed: %w", err)
}
```

**При необходимости создавайте собственные типы ошибок:**

```go
type ValidationError struct {
    Field string
    Value string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("invalid value %q for field %q", e.Value, e.Field)
}
```

### Комментарии и документация

**Комментарии к пакетам:**

```go
// Package application provides the core Wails application runtime.
//
// It handles window management, event dispatching, and service lifecycle.
package application
```

**Экспортируемые объявления:**

```go
// NewApplication creates a new Wails application with the given options.
//
// The application must be started with Run() or RunWithContext().
func NewApplication(opts Options) *Application {
    // ...
}
```

**Комментарии к реализации:**

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

### Структура функций и методов

**Каждая функция должна решать одну конкретную задачу:**

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

**Используйте ранний возврат:**

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

### Конкурентное выполнение

**Используйте контекст для отмены операций:**

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

**Защищайте общее состояние с помощью мьютексов:**

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

**Не допускайте утечек горутин:**

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

### Тестирование

**Именование файлов тестов:**

```go
// Implementation: window.go
// Tests: window_test.go
```

**Тесты на основе таблиц:**

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

## Стандарты JavaScript/TypeScript

### Форматирование кода

Для единообразного форматирования используйте Prettier:

```json
{
  "semi": false,
  "singleQuote": true,
  "tabWidth": 2,
  "trailingComma": "es5"
}
```

### Соглашения об именовании

**Переменные и функции:**

- camelCase: `const userName = "John"`

**Классы и типы:**

- PascalCase: `class WindowManager`

**Константы:**

- UPPER<em>SNAKE</em>CASE: `const MAX_RETRIES = 3`

### TypeScript

**Указывайте типы явно:**

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

**Определяйте интерфейсы:**

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

## Формат сообщений коммитов

Используйте формат [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Типы:**

- `feat`: новая функциональность
- `fix`: исправление ошибки
- `docs`: изменения документации
- `refactor`: рефакторинг кода
- `test`: добавление или обновление тестов
- `chore`: задачи по сопровождению

**Примеры:**

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

## Рекомендации по запросам на включение изменений

### Перед отправкой

- [ ] Код успешно проходит `gofmt` и `goimports`
- [ ] Все тесты проходят успешно (`go test ./...`)
- [ ] Для нового кода добавлены тесты
- [ ] При необходимости документация обновлена
- [ ] Сообщения коммитов соответствуют принятым соглашениям
- [ ] Нет конфликтов слияния с `master`

### Шаблон описания запроса на включение изменений

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

## Процесс проверки кода

### Для рецензента

- Давайте конструктивные и уважительные отзывы
- Сосредоточьтесь на качестве кода, а не на личных предпочтениях
- Объясняйте причины предлагаемых изменений
- Одобрите изменения, когда будете удовлетворены результатом

### Для автора

- Отвечайте на все комментарии
- При необходимости просите разъяснений
- Внесите запрошенные изменения или объясните, почему этого делать не следует
- Будьте открыты к обратной связи

## Рекомендации

### Производительность

- Избегайте преждевременной оптимизации
- Перед оптимизацией выполните профилирование
- Используйте бенчмарки для критичного к производительности кода

```go
func BenchmarkProcess(b *testing.B) {
    for i := 0; i < b.N; i++ {
        process(testData)
    }
}
```

### Безопасность

- Проверяйте все данные, вводимые пользователем
- Очищайте данные перед отображением
- Используйте `crypto/rand` для генерации случайных данных
- Никогда не записывайте конфиденциальную информацию в журнал

### Документация

- Документируйте экспортируемые API
- Добавляйте примеры в документацию
- Обновляйте документацию при изменении API
- Поддерживайте файлы README в актуальном состоянии

## Код для конкретных платформ

### Именование файлов

```
window.go           // Common interface
window_darwin.go    // macOS implementation
window_windows.go   // Windows implementation
window_linux.go     // Linux implementation
```

### Теги сборки

```go
//go:build darwin

package application

// macOS-specific code
```

## Статический анализ

Перед коммитом запустите линтеры:

```bash
# golangci-lint (recommended)
golangci-lint run

# Individual linters
go vet ./...
staticcheck ./...
```

## Остались вопросы?

Если вы не уверены в каких-либо стандартах:

- Найдите примеры в существующем коде
- Задайте вопрос в [Discord](https://discord.gg/JDdSxwjhGf)
- Создайте обсуждение на GitHub
