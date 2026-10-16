---
title: "Необработанные сообщения"
description: "Реализация пользовательского обмена данными между фронтендом и бэкендом для приложений, критичных к производительности"
slug: "guides/raw-messages"
sourcePath: "guides/raw-messages.md"
---

Необработанные сообщения обеспечивают низкоуровневый канал связи между фронтендом и бэкендом в обход стандартной системы привязок. Это повышает скорость за счёт удобства.

## Когда использовать необработанные сообщения

Необработанные сообщения лучше всего подходят для крайне редких особых случаев:

- **Сверхчастые обновления** — тысячи сообщений в секунду, когда важна каждая микросекунда
- **Пользовательские протоколы обмена сообщениями** — когда необходим полный контроль над форматом передаваемых данных

@note{type="tip"}
Почти во всех случаях рекомендуется использовать стандартные [привязки сервисов](/features/bindings/services/): они обеспечивают типобезопасность, автоматическую сериализацию и более удобную разработку при пренебрежимо малых накладных расходах.

@end

## Настройка бэкенда

Настройте `RawMessageHandler` в параметрах приложения:

```go
package main

import (
    "encoding/json"
    "fmt"

    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "Raw Message Demo",
        Assets: application.AssetOptions{
            Handler: application.BundledAssetFileServer(assets),
        },
        RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
            fmt.Printf("Raw message from window '%s': %s (origin: %+v)\n", window.Name(), message, originInfo.Origin)

            // Process the message and respond via events
            response := processMessage(message)
            window.EmitEvent("raw-response", response)
        },
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: "My App",
        Name:  "main",
    })

    app.Run()
}

func processMessage(message string) map[string]any {
    // Your custom message processing logic
    return map[string]any{
        "received": message,
        "status":   "processed",
    }
}
```

### Сигнатура обработчика

```go
RawMessageHandler func(window Window, message string, originInfo *application.OriginInfo)
```

| Параметр | Тип | Описание |
| --- | --- | --- |
| `window` | `Window` | Окно, отправившее сообщение |
| `message` | `string` | Необработанное содержимое сообщения |
| `originInfo` | `*application.OriginInfo` | Сведения об источнике сообщения |

#### Структура OriginInfo

```go
type OriginInfo struct {
	Origin      string
	TopOrigin   string
	IsMainFrame bool
}
```

| Поле | Тип | Описание |
| --- | --- | --- |
| `Origin` | `string` | URL источника документа, отправившего сообщение |
| `TopOrigin` | `string` | URL источника верхнего уровня (во фреймах iframe может отличаться от Origin) |
| `IsMainFrame` | `bool` | Указывает, было ли сообщение отправлено из главного фрейма |

#### Доступность в зависимости от платформы

- **macOS**: предоставляются `Origin` и `IsMainFrame`
- **Windows**: предоставляются `Origin` и `TopOrigin`
- **Linux**: предоставляется только `Origin`

### Проверка источника

@note{type="caution"}
Никогда не считайте сообщение безопасным лишь потому, что оно поступило в обработчик. Перед выполнением операций, важных с точки зрения безопасности, или операций, изменяющих состояние, необходимо проверить сведения об источнике.

@end

**Всегда проверяйте источник входящих сообщений перед их обработкой.** Параметр `originInfo` содержит критически важные сведения для обеспечения безопасности, которые необходимо проверять для предотвращения несанкционированного доступа. Необработанные сообщения могут отправлять вредоносное или скомпрометированное содержимое, а также непредусмотренные скрипты. Без проверки источника вы можете обработать команды из ненадёжных источников. Используйте `originInfo`, чтобы убедиться, что сообщения поступают из ожидаемых источников.

### Основные аспекты проверки

- **Всегда проверяйте `Origin`** — убедитесь, что источник соответствует ожидаемым доверенным источникам (обычно `wails://wails` или `http://wails.localhost` для локальных ресурсов либо конкретный источник вашего приложения)
- **Проверяйте `IsMainFrame`** (macOS) — учитывайте, поступило ли сообщение из iframe, поскольку это может указывать на встроенное содержимое с другим контекстом безопасности
- **Используйте `TopOrigin`** (Windows) — при работе с содержимым во фреймах проверяйте источник верхнего уровня
- **Отклоняйте сообщения из неожиданных источников** — для безопасного отказа отклоняйте сообщения из источников, которые не разрешены явно

@note{type="info"}
Сообщения с префиксом `wails:` зарезервированы для внутреннего обмена данными Wails и не передаются вашему обработчику.

@end

## Настройка фронтенда

Отправляйте необработанные сообщения с помощью `System.invoke()`:

```html
<!DOCTYPE html>
<html>
<head>
    <script type="module">
        import { System, Events } from '@wailsio/runtime'

        // Send raw message
        document.getElementById('send').addEventListener('click', () => {
            const message = document.getElementById('input').value
            System.invoke(message)
        })

        // Listen for response
        Events.On('raw-response', (event) => {
            console.log('Response:', event.data)
        })
    </script>
</head>
<body>
    <input type="text" id="input" placeholder="Enter message" />
    <button id="send">Send</button>
</body>
</html>
```

### Использование готового пакета

Если вы не используете npm, обращайтесь к `invoke` через глобальный объект `wails`:

```html
<script type="module" src="/wails/runtime.js"></script>
<script>
    window.onload = function() {
        document.getElementById('send').onclick = function() {
            wails.System.invoke('my-message')
        }
    }
</script>
```

## Структурированные сообщения

Для сложных данных используйте сериализацию в JSON:

### Фронтенд

```javascript
import { System } from '@wailsio/runtime'

const command = {
    action: 'update',
    payload: {
        id: 123,
        value: 'new value'
    }
}

System.invoke(JSON.stringify(command))
```

### Бэкенд

```go
RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
    var cmd struct {
        Action  string `json:"action"`
        Payload struct {
            ID    int    `json:"id"`
            Value string `json:"value"`
        } `json:"payload"`
    }

    if err := json.Unmarshal([]byte(message), &cmd); err != nil {
        window.EmitEvent("error", err.Error())
        return
    }

    switch cmd.Action {
    case "update":
        // Handle update
        result := handleUpdate(cmd.Payload.ID, cmd.Payload.Value)
        window.EmitEvent("update-complete", result)
    default:
        window.EmitEvent("error", "unknown action")
    }
}
```

## Сравнение производительности

| Подход | Накладные расходы | Типобезопасность | Сценарий использования |
| --- | --- | --- | --- |
| Привязки сервисов | Выше | Полная | Общее назначение |
| Необработанные сообщения | Минимальные | Вручную | Высокая частота, критичные требования к производительности |

### Пример теста производительности

При простой полезной нагрузке необработанные сообщения позволяют обрабатывать значительно больше сообщений в секунду, чем привязки сервисов:

```go
// Raw message handler - minimal overhead
RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
    // Direct string processing, no reflection or marshaling
    counter++
}
```

## Полный пример

Ниже приведён полный пример реализации простого протокола команд:

### main.go

```go
package main

import (
    "embed"
    "encoding/json"
    "fmt"
    "time"

    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

type Command struct {
    Type string          `json:"type"`
    Data json.RawMessage `json:"data"`
}

func main() {
    app := application.New(application.Options{
        Name: "Raw Message Demo",
        Assets: application.AssetOptions{
            Handler: application.BundledAssetFileServer(assets),
        },
        Mac: application.MacOptions{
            ApplicationShouldTerminateAfterLastWindowClosed: true,
        },
        RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
            var cmd Command
            if err := json.Unmarshal([]byte(message), &cmd); err != nil {
                window.EmitEvent("error", map[string]string{"error": err.Error()})
                return
            }

            switch cmd.Type {
            case "ping":
                window.EmitEvent("pong", map[string]any{
                    "time":   time.Now().UnixMilli(),
                    "window": window.Name(),
                })
            case "echo":
                var text string
                json.Unmarshal(cmd.Data, &text)
                window.EmitEvent("echo", text)
            default:
                window.EmitEvent("error", map[string]string{
                    "error": fmt.Sprintf("unknown command: %s", cmd.Type),
                })
            }
        },
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: "Raw Message Demo",
        Name:  "main",
        Width: 400,
        Height: 300,
    })

    app.Run()
}
```

### assets/index.html

```html
<!DOCTYPE html>
<html>
<head>
    <title>Raw Message Demo</title>
    <style>
        body { font-family: sans-serif; padding: 20px; }
        button { margin: 5px; padding: 10px 20px; }
        #output { margin-top: 20px; padding: 10px; background: #f0f0f0; }
    </style>
</head>
<body>
    <h1>Raw Message Demo</h1>

    <button id="ping">Ping</button>
    <button id="echo">Echo "Hello"</button>

    <div id="output">Waiting for response...</div>

    <script type="module">
        import { System, Events } from '@wailsio/runtime'

        const output = document.getElementById('output')

        function send(type, data) {
            System.invoke(JSON.stringify({ type, data }))
        }

        document.getElementById('ping').onclick = () => send('ping')
        document.getElementById('echo').onclick = () => send('echo', 'Hello')

        Events.On('pong', (e) => {
            output.textContent = `Pong from ${e.data.window} at ${e.data.time}`
        })

        Events.On('echo', (e) => {
            output.textContent = `Echo: ${e.data}`
        })

        Events.On('error', (e) => {
            output.textContent = `Error: ${e.data.error}`
        })
    </script>
</body>
</html>
```

## Рекомендации

### Что следует делать

- Используйте необработанные сообщения для участков, действительно критичных к производительности
- Реализуйте в обработчике надлежащую обработку ошибок
- Используйте события для отправки ответов во фронтенд
- Рассмотрите возможность использования JSON для структурированных данных
- Обрабатывайте сообщения быстро, чтобы избежать блокировки

### Чего не следует делать

- Не используйте необработанные сообщения, если достаточно привязок сервисов
- Не забывайте проверять входящие сообщения
- Не блокируйте обработчик длительными операциями (используйте горутины)
- Не игнорируйте параметр окна, если ответы нужно направлять определённым окнам

## Особенности работы с несколькими окнами

Параметр `window` определяет, какое окно отправило сообщение, и позволяет:

- Отправлять ответы в нужное окно
- Реализовать поведение для отдельных окон
- Отслеживать источники сообщений при отладке

```go
RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
    // Respond only to the sending window
    window.EmitEvent("response", result)

    // Or broadcast to all windows
    app.Event.Emit("broadcast", result)
}
```

## Дальнейшие шаги

- [Привязки сервисов](/features/bindings/services/) — стандартный подход для большинства приложений
- [События](/guides/events-reference/) — система событий для передачи данных от бэкенда к фронтенду
- [Производительность](/guides/performance/) — общая оптимизация производительности
