---
title: "Справочник API"
description: "Справочник по типам, методам, параметрам и командам CLI Wails"
slug: "reference/overview"
sourcePath: "reference/overview.md"
---

## Об этом справочнике

Используйте этот раздел как справочник по типам, методам, параметрам и командам Wails при разработке приложения. Если вам нужен пошаговый учебный материал, начните с раздела [«Учебные руководства»](/tutorials/overview/). Инструкции по выполнению конкретных задач см. в разделе [«Практические руководства»](/quick-start/next-steps/#features). Основные сведения о проектировании см. в разделе [«Ключевые концепции»](/concepts/architecture/).

**Структура:**

- [Приложение](/reference/application/) — основные API приложения
- [Окно](/reference/window/) — создание окон и управление ими; поля конфигурации перечислены в разделе [«Параметры окна»](/features/windows/options/)
- [Меню](/reference/menu/) — API меню приложения, контекстного меню и системного лотка; типы и свойства элементов описаны в [справочнике по элементам меню](/features/menus/reference/)
- [События](/reference/events/) — система событий и встроенные события
- [Диалоговые окна](/reference/dialogs/) — диалоговые окна для работы с файлами и сообщениями
- [Среда выполнения фронтенда](/reference/frontend-runtime/) — API среды выполнения фронтенда
- [Обзор команд CLI](/reference/cli/) — команды, сгруппированные по назначению; сведения об отдельных командах приведены в разделе [«Флаги и примеры»](/guides/cli/)
- [Протокол манифеста обновлений](/reference/update-manifest/) — метаданные обновлений и поля проверки
- [Мобильный API](/guides/mobile/mobile-api/) — экспериментальные мобильные API

## Соглашения API

@details{title="Соглашения Go API — для разработчиков, начинающих работать с Go"}
### Именование

- <strong>Типы</strong>: PascalCase (например, `WebviewWindow`)
- <strong>Методы</strong>: PascalCase (например, `SetTitle()`)
- <strong>Параметры</strong>: структуры с именами в PascalCase (например, `WindowOptions`)
- <strong>Константы</strong>: PascalCase (например, `WindowStartStateMaximised`)

#### Обработка ошибок

Большинство методов, при выполнении которых может возникнуть ошибка, возвращают `error` как последнее возвращаемое значение. `app.Run()` блокирует выполнение до завершения работы приложения и возвращает любую ошибку запуска:

```go
if err := app.Run(); err != nil {
    log.Fatal(err)
}
```

При создании окна ошибка не возвращается: `app.Window.New()` непосредственно возвращает `*WebviewWindow`.

#### Контекст

Методы жизненного цикла сервиса получают `context.Context`:

```go
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    // ctx is cancelled when the application is shutting down.
    return nil
}
```

Контекст жизненного цикла приложения доступен через `app.Context()`. `RunWithContext` отсутствует — вызывайте `app.Run()`.

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

### Соглашения JavaScript API

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

## Справочник по типам

### Распространённые типы

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

Поведение некоторых API зависит от платформы:

| Возможность | Windows | macOS | Linux |
| --- | --- | --- | --- |
| **Меню приложения** | Строка меню окна | Глобальная строка меню | Строка меню окна |
| **Системный лоток** | Область уведомлений | Строка меню | Системный лоток |
| **Dock** | Неприменимо | ✅ Доступно | Неприменимо |
| **Диалоговые окна для работы с файлами** | Нативные | Нативные | Нативные (GTK) |
| **Прозрачность** | ✅ Полная | Требуется [`private_mac_apis`](/guides/build/private-macos-apis/#webview-transparency-and-background) | ⚠️ Ограниченная |

Особенности поведения на разных платформах описаны в каждом разделе API.

## Версионирование

Wails v3 следует правилам семантического версионирования:

- **Основная версия** (v4.0.0): несовместимые изменения
- **Дополнительная версия** (v3.1.0): новые возможности с обратной совместимостью
- **Версия исправления** (v3.0.1): исправления ошибок с обратной совместимостью

**Wails v3.0.0:** стабильная версия

## Политика устаревания

Когда API объявляются устаревшими:

1. **В документации добавляется отметка** с уведомлением об устаревании
2. **Предоставляется альтернатива** с руководством по миграции
3. **Поддержка сохраняется в течение 1 основной версии** до удаления
4. **Предупреждения компилятора** (где это возможно)

## Стабильность API

### Стабильные API ✅

Эти API стабильны и подходят для использования в рабочей среде:

- Основные API приложения
- Управление окнами
- Система меню
- Система событий
- Диалоговые окна для работы с файлами
- Привязки сервисов

### Экспериментальные API ⚠️

За рамки контракта стабильной совместимости для настольных платформ выходят только API, явно указанные в документации как экспериментальные:

- Поддержка Android
- Поддержка iOS
- Другие API, явно помеченные как экспериментальные

Экспериментальные API имеют соответствующую пометку в документации.

## Получение помощи

### Вопросы об API

1. **Ознакомьтесь с этим справочником** — полная документация по API
2. **Посмотрите примеры** — [примеры на GitHub](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples)
3. **Выполните поиск в Discord** — [сервер Discord](https://discord.gg/JDdSxwjhGf)
4. **Задайте вопрос сообществу** — канал #help в Discord

### Сообщение о проблемах с API

Обнаружили ошибку или несоответствие?

1. **Проверьте существующие задачи** — [задачи на GitHub](https://github.com/wailsapp/wails/issues)
2. **Создайте подробный отчёт** — приложите код, текст ошибки и укажите платформу
3. **Предоставьте пример для воспроизведения** — минимальный пример, демонстрирующий проблему

## Связанная документация

- [Учебные материалы](/tutorials/overview/) — учитесь, создавая реальные приложения
- [Практические руководства](/guides/architecture/) — ориентированные на задачи руководства для распространённых сценариев
- [Возможности](/features/windows/basics/) — документация по каждой возможности
- [Примеры](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples) — рабочие примеры кода на GitHub

---

**Обзор API:** используйте навигацию слева, чтобы ознакомиться с конкретными API.
