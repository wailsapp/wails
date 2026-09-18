---
title: "Справочник по API"
description: "Полная документация по API Wails v3"
slug: "reference/overview"
sourcePath: "reference/overview.md"
---

## Об этом справочнике

Это полный справочник по API Wails v3. В нём документированы все общедоступные типы, методы и параметры, имеющиеся во фреймворке.

**Структура:**

- [Приложение](/reference/application/) — основные API приложения
- [Окно](/reference/window/) — создание окон и управление ими
- [Меню](/reference/menu/) — меню приложения, контекстные меню и меню в области уведомлений
- [События](/reference/events/) — система событий и встроенные события
- [Диалоговые окна](/reference/dialogs/) — диалоговые окна для работы с файлами и сообщениями
- [Среда выполнения фронтенда](/reference/frontend-runtime/) — API среды выполнения фронтенда
- [CLI](/reference/cli/) — интерфейс командной строки

## Соглашения API

@details{title="Соглашения Go API — для разработчиков, впервые работающих с Go"}
### Именование

- <strong></strong>Типы<strong></strong>: PascalCase (например, `WebviewWindow`)
- <strong></strong>Методы<strong></strong>: PascalCase (например, `SetTitle()`)
- <strong></strong>Параметры<strong></strong>: структуры с именами в стиле PascalCase (например, `WindowOptions`)
- <strong></strong>Константы<strong></strong>: PascalCase (например, `WindowStartStateMaximised`)

#### Обработка ошибок

Большинство методов, при выполнении которых может произойти ошибка, возвращают `error` последним возвращаемым значением. `app.Run()` блокирует выполнение до завершения работы приложения и возвращает ошибку запуска, если она возникла:

```go
if err := app.Run(); err != nil {
    log.Fatal(err)
}
```

Создание окна не возвращает ошибку: `app.Window.New()` непосредственно возвращает `*WebviewWindow`.

#### Контекст

Методы жизненного цикла сервиса получают `context.Context`:

```go
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    // ctx is cancelled when the application is shutting down.
    return nil
}
```

Контекст времени жизни приложения доступен через `app.Context()`. Метода `RunWithContext` нет — вызывайте `app.Run()`.

#### Шаблон параметров

Для конфигурации используются структуры параметров:

```go
app := application.New(application.Options{
    Name: "My App",
    Description: "A demo application",
    Services: []application.Service{
        application.NewService(&MyService{}),
    },
})
```

@end

### Соглашения API JavaScript

#### Именование

- **Функции**: camelCase (например, `setTitle()`)
- **Константы**: SCREAMING<em>SNAKE</em>CASE (например, `WINDOW_EVENT_FOCUS`)

#### Асинхронность по умолчанию

Все вызовы методов Go возвращают Promise:

```javascript
// Async/await (recommended)
const result = await MyService.DoSomething()

// Promise chain
MyService.DoSomething()
    .then(result => console.log(result))
    .catch(error => console.error(error))
```

#### Обработка ошибок

Ошибки Go преобразуются в исключения JavaScript:

```javascript
try {
    await MyService.MightFail()
} catch (error) {
    console.error('Go error:', error)
}
```

#### Безопасность типов

Определения TypeScript создаются автоматически:

```typescript
// Fully typed
import { Greet } from './bindings/GreetService'

const message: string = await Greet("World")
```

## Структура пакетов

```
github.com/wailsapp/wails/v3/pkg/
├── application/          # Core application package
│   ├── application.go    # App type
│   ├── webview_window.go # Window management
│   ├── menu.go           # Menu types
│   ├── event_manager.go  # Event system
│   └── dialogs.go        # Dialog APIs
├── events/               # Event constants
└── services/             # Built-in services
    ├── dock/             # macOS dock (includes badge support)
    ├── fileserver/       # File-server service
    ├── kvstore/          # Key/value store
    ├── log/              # Structured logging service
    ├── notifications/    # Notifications service
    └── sqlite/           # SQLite service
```

## Пути импорта

### Go

```go
import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)
```

### JavaScript

```javascript
// Auto-generated bindings
import { MyMethod } from './bindings/MyService'

// Runtime APIs
import { Events, Window } from '@wailsio/runtime'
```

## Справочник типов

### Общие типы

@tabs{sync-key="lang"}
[Go]
```go
// Application
type App struct { /* ... */ }
type Options struct { /* ... */ }

// Window
type WebviewWindow struct { /* ... */ } // implements the Window interface
type WebviewWindowOptions struct { /* ... */ }

// Menu
type Menu struct { /* ... */ }
type MenuItem struct { /* ... */ }

// Events — there is no generic Event type; events are typed by source.
type ApplicationEvent struct { /* ... */ }
type WindowEvent struct { /* ... */ }
type CustomEvent struct { /* ... */ }
type EventListener struct { /* ... */ }

// Dialogs
type OpenFileDialogOptions struct { /* ... */ }
type SaveFileDialogOptions struct { /* ... */ }
```

[TypeScript]
```typescript
// Window runtime
interface WindowOptions {
    title?: string
    width?: number
    height?: number
    // ...
}

// Events
type EventCallback = (data: any) => void

// Bindings (auto-generated)
export function MyMethod(arg: string): Promise<string>
```

@end

## Различия между платформами

Поведение некоторых API различается в зависимости от платформы:

| Возможность | Windows | macOS | Linux |
| --- | --- | --- | --- |
| **Меню приложения** | Строка меню окна | Глобальная строка меню | Строка меню окна |
| **Область уведомлений** | Область уведомлений | Строка меню | Область уведомлений |
| **Dock** | Неприменимо | ✅ Доступно | Неприменимо |
| **Диалоговые окна для работы с файлами** | Нативные | Нативные | Нативные (GTK) |
| **Прозрачность** | ✅ Полная | Требуется [`private_mac_apis`](/guides/build/private-macos-apis/#webview-transparency-and-background) | ⚠️ Ограниченная |

Особенности поведения на разных платформах документированы в каждом разделе API.

## Управление версиями

Wails v3 следует правилам семантического версионирования:

- **Мажорная версия** (v3.x.x): изменения, нарушающие обратную совместимость
- **Минорная версия** (v3.x.x): новые возможности с сохранением обратной совместимости
- **Патч-версия** (v3.x.x): исправления ошибок с сохранением обратной совместимости

**Текущий статус:** бета-версия (API стабилен, доработка продолжается)

## Политика устаревания

Когда API объявляются устаревшими:

1. **Помечаются в документации** уведомлением об устаревании
2. **Предоставляется альтернатива** с руководством по переходу
3. **Поддерживаются в течение 1 основной версии** перед удалением
4. **Предупреждения компилятора** (где возможно)

## Стабильность API

### Стабильные API ✅

Эти API стабильны и подходят для использования в рабочей среде:

- Основные API приложения
- Управление окнами
- Система меню
- Система событий
- Диалоговые окна выбора файлов
- Привязки сервисов

### Нестабильные API ⚠️

Эти API могут измениться до финального выпуска:

- Некоторые расширенные параметры окон
- Возможности, зависящие от платформы
- Экспериментальные возможности

Нестабильные API помечены в документации.

## Получение помощи

### Вопросы об API

1. **Ознакомьтесь с этим справочником** — полной документацией по API
2. **Изучите примеры** — [примеры на GitHub](https://github.com/wailsapp/wails/tree/master/v3/examples)
3. **Выполните поиск в Discord** — [сервер Discord](https://discord.gg/JDdSxwjhGf)
4. **Задайте вопрос сообществу** — в канале Discord #help

### Сообщение о проблемах с API

Обнаружили ошибку или несоответствие?

1. **Проверьте существующие обращения** — [обращения на GitHub](https://github.com/wailsapp/wails/issues)
2. **Составьте подробный отчёт** — укажите код, ошибку и платформу
3. **Предоставьте способ воспроизведения** — минимальный пример, демонстрирующий проблему

## Связанная документация

- [Учебные материалы](/tutorials/overview/) — обучение на примере создания реальных приложений
- [Руководства](/guides/architecture/) — практические инструкции для типовых сценариев
- [Возможности](/features/windows/basics/) — документация по каждой возможности
- [Примеры](https://github.com/wailsapp/wails/tree/master/v3/examples) — рабочие примеры кода на GitHub

---

**Просмотр API:** используйте навигацию слева, чтобы изучить конкретные API.
