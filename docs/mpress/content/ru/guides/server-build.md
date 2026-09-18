---
title: "Серверная сборка"
description: "Запуск приложений Wails в качестве HTTP-серверов без нативного окна графического интерфейса"
slug: "guides/server-build"
sourcePath: "guides/server-build.md"
---

Wails v3 поддерживает серверный режим, в котором приложение можно запускать как обычный HTTP-сервер без создания нативных окон и зависимостей графического интерфейса. Это позволяет развёртывать одно и то же приложение Wails на серверах, в контейнерах и для доступа через веб-браузеры.

Серверный режим полезен для следующих сценариев:

- **Развёртывание в Docker и других контейнерах** — работа без зависимостей X11/Wayland
- **Серверные приложения** — развёртывание в виде веб-сервера, доступного через браузер
- **Доступ только через веб-интерфейс** — единая кодовая база для настольной и веб-версий
- **Тестирование в CI/CD** — запуск интеграционных тестов без сервера отображения
- **Микросервисы** — использование привязок Wails в серверных службах без графического интерфейса

## Быстрый старт

Серверный режим включается с помощью тега сборки `server`. Код приложения остаётся прежним — достаточно выполнить сборку с этим тегом:

```bash
# Using Taskfile (recommended)
wails3 task build:server
wails3 task run:server

# Or build directly with Go
go build -tags server -o myapp-server .
```

Минимальный пример:

```go
package main

import (
    "embed"
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name: "My App",
        // Server options are used when built with -tags server
        Server: application.ServerOptions{
            Host: "localhost",
            Port: 8080,
        },
        Services: []application.Service{
            application.NewService(&MyService{}),
        },
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
    })

    log.Println("Starting application...")
    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

Один и тот же код можно собрать для настольного режима (без тега) или серверного режима (с тегом `-tags server`).

## Конфигурация

### ServerOptions

Настройте HTTP-сервер с помощью `ServerOptions`:

```go
Server: application.ServerOptions{
    // Host to bind to. Default: "localhost"
    // Use "0.0.0.0" to listen on all interfaces
    Host: "localhost",

    // Port to listen on. Default: 8080
    Port: 8080,

    // Request read timeout. Default: 30s
    ReadTimeout: 30 * time.Second,

    // Response write timeout. Default: 30s
    WriteTimeout: 30 * time.Second,

    // Idle connection timeout. Default: 120s
    IdleTimeout: 120 * time.Second,

    // Graceful shutdown timeout. Default: 30s
    ShutdownTimeout: 30 * time.Second,

    // Additional origins allowed to open WebSocket connections.
    // Same-origin connections are always allowed.
    WebSocketOriginPatterns: []string{"app.example.com"},

    // Disable WebSocket origin checks. Unsafe; default: false.
    WebSocketAllowAllOrigins: false,

    // TLS configuration (optional)
    TLS: &application.TLSOptions{
        CertFile: "/path/to/cert.pem",
        KeyFile:  "/path/to/key.pem",
    },
},
```

## Возможности

### Конечная точка проверки работоспособности

Конечная точка проверки работоспособности автоматически доступна по адресу `/health`:

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

Она полезна для следующих задач:

- Проверки активности и готовности Kubernetes
- Проверки работоспособности балансировщиком нагрузки
- Системы мониторинга

### Привязки сервисов

Все привязки сервисов работают так же, как в настольном режиме:

```go
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello, " + name + "!"
}

// Register in options
Services: []application.Service{
    application.NewService(&GreetService{}),
},
```

Клиентская часть может вызывать эти привязки через стандартную среду выполнения Wails:

```javascript
const greeting = await wails.Call.ByName('main.GreetService.Greet', 'World');
```

### События

В серверном режиме события передаются в обоих направлениях:

- **Из клиентской части в серверную**: события, отправленные из браузера, передаются по HTTP и обрабатываются обработчиками событий Go
- **Из серверной части в клиентскую**: события, отправленные из Go, транслируются всем подключённым браузерам через WebSocket

Каждая вкладка браузера представлена как «окно» с уникальным именем (`browser-1`, `browser-2` и т. д.), доступное через `event.Sender`:

```go
// Listen for events from browsers
app.Event.On("user-action", func(event *application.CustomEvent) {
    log.Printf("Event from %s: %v", event.Sender, event.Data)
    // event.Sender will be "browser-1", "browser-2", etc.
})

// Emit events to all connected browsers
app.Event.Emit("server-update", data)
```

Из клиентской части:

```javascript
// Emit event to server (and all other browsers)
await wails.Events.Emit('user-action', { action: 'click' });

// Listen for events from server
wails.Events.On('server-update', (event) => {
    console.log('Update from server:', event.data);
});
```

### Корректное завершение работы

Сервер корректно обрабатывает сигналы `SIGINT` и `SIGTERM`:

1. Прекращает принимать новые подключения
2. Ожидает завершения активных запросов (не более `ShutdownTimeout`)
3. Выполняет обработчики `OnShutdown`
4. Останавливает сервисы в обратном порядке

## Отличия от настольного режима

| Возможность | Настольный режим | Серверный режим |
| --- | --- | --- |
| Нативные окна | Создаются | Окна браузера (`browser-N`) |
| Область уведомлений | Доступна | Недоступна |
| Нативные диалоговые окна | Доступны | Недоступны |
| Меню приложения | Доступно | Недоступно |
| Сведения об экране | Доступны | Возвращается ошибка |
| Привязки сервисов | Работают | Работают |
| События | Работают | Работают (через WebSocket) |
| Ресурсы | Через webview | Через HTTP |
| Требуется CGO | Да | Нет |

### Поведение API окон

В серверном режиме API, связанные с окнами, обрабатываются безопасно:

- `app.Window.NewWithOptions()` — записывает предупреждение в журнал и возвращает nil
- `app.Hide()` / `app.Show()` — не выполняют никаких действий
- `app.Screen.GetPrimary()` — возвращает ошибку

Благодаря этому код, обращающийся к окнам, выполняется без аварийного завершения, хотя операции с окнами не дают никакого эффекта.

## Сборка для промышленной эксплуатации

### С помощью Task (рекомендуется)

Проекты, созданные с помощью `wails3 init`, включают задачу `build:server`:

```bash
# Build for server mode
wails3 task build:server

# Build and run
wails3 task run:server
```

### Сборка вручную

```bash
# Build with server mode
go build -tags server -o myapp-server .
```

### Docker

Проекты Wails включают готовую к использованию конфигурацию Docker. Чтобы собрать и запустить приложение в контейнере:

```bash
# Build the Docker image
wails3 task build:docker

# Run it
wails3 task run:docker
```

Готово! Приложение будет доступно по адресу `http://localhost:8080`.

Сборку можно настроить с помощью нескольких параметров:

```bash
# Use a custom image tag
wails3 task build:docker TAG=myapp:v1.0.0

# Run on a different port
wails3 task run:docker PORT=3000
```

Созданный `Dockerfile.server` формирует минимальный образ на основе distroless. Привязка к сетевому интерфейсу выполняется автоматически, поэтому приложение будет доступно извне контейнера.

### Docker Compose

Для более сложных развёртываний используйте следующую конфигурацию Docker Compose с проверками работоспособности:

```yaml
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - WAILS_SERVER_HOST=0.0.0.0
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

@note{type="info"}
В примере проверки работоспособности используется `wget`. Если вы используете базовый образ distroless, добавьте в образ исполняемый файл для проверки работоспособности либо используйте внешний механизм проверки (например, параметр Docker `curl` или контейнер-сайдкар).

@end

### Собственный Dockerfile

Если вам требуется больше контроля, создайте собственный Dockerfile. Главное — задать `WAILS_SERVER_HOST=0.0.0.0`, чтобы сервер принимал подключения извне контейнера:

```dockerfile
# Build stage
FROM golang:alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY . .
RUN go mod tidy
RUN go build -tags server -ldflags="-s -w" -o server .

# Runtime stage
FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/server /server
COPY --from=builder /app/frontend/dist /frontend/dist
EXPOSE 8080
ENV WAILS_SERVER_HOST=0.0.0.0
ENTRYPOINT ["/server"]
```

## Рекомендации по безопасности

При развёртывании приложений в серверном режиме:

1. **По умолчанию выполняйте привязку к localhost** — используйте `0.0.0.0` только при необходимости
2. **Используйте TLS в рабочей среде** — настройте `ServerOptions.TLS`
3. **Размещайте приложение за обратным прокси-сервером** — используйте nginx/traefik для дополнительной защиты
4. **Оставляйте WebSocket-подключения в пределах одного источника** — добавляйте с помощью `WebSocketOriginPatterns` только доверенные источники; избегайте `WebSocketAllowAllOrigins`
5. **Проверяйте все входные данные** — применяйте те же методы обеспечения безопасности, что и для любого веб-приложения

## Пример

Полный пример доступен в `v3/examples/server/`:

```bash
cd v3/examples/server

# Using Taskfile
task dev

# Or run directly
go run -tags server .

# Open http://localhost:8080 in browser
```

## Переменные окружения

Для сценариев развёртывания, в которых требуется переопределить конфигурацию сервера без изменения кода, Wails распознаёт следующие переменные окружения:

| Переменная | Описание | Значение по умолчанию |
| --- | --- | --- |
| `WAILS_SERVER_HOST` | Сетевой интерфейс для привязки | `localhost` |
| `WAILS_SERVER_PORT` | Порт для приёма подключений | `8080` |

Эти переменные имеют приоритет над `ServerOptions` в вашем коде. Поэтому в примерах Docker задаётся `WAILS_SERVER_HOST=0.0.0.0`: это позволяет контейнеру принимать внешние подключения без каких-либо изменений в приложении.

## См. также

- [Пользовательский транспорт](/guides/custom-transport/) — расширенная настройка IPC
- [Сервисы](/features/bindings/services/) — документация по привязке сервисов
- [События](/guides/events-reference/) — документация по системе событий

### Размер запросов среды выполнения

Размер запросов к `/wails/runtime` до обработки JSON ограничен 64 МиБ. Обычный запрос, превышающий этот предел, получает ответ HTTP 413; это относится и к запросам без `Content-Length`. Для фрагментированных загрузок среды выполнения сохраняются отдельные ограничения: 1 МиБ на фрагмент и 64 МиБ на собранную полезную нагрузку. При необходимости установите более низкий предел с помощью промежуточного ПО приложения или обратного прокси-сервера.
