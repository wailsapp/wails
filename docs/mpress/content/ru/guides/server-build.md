---
title: "Серверная сборка"
description: "Запуск приложений Wails в качестве HTTP-серверов без нативного окна графического интерфейса"
slug: "guides/server-build"
sourcePath: "guides/server-build.md"
---

Wails v3 поддерживает серверный режим, позволяющий запускать приложение как обычный HTTP-сервер без создания нативных окон и зависимостей графического интерфейса. Благодаря этому одно и то же приложение Wails можно развёртывать на серверах, в контейнерах и веб-браузерах.

Серверный режим подходит для следующих сценариев:

- **Развёртывание в Docker и других контейнерах** — работа без зависимостей X11/Wayland
- **Серверные приложения** — развёртывание в виде веб-сервера, доступного через браузер
- **Доступ только через веб-интерфейс** — единая кодовая база для настольного приложения и веб-версии
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

Ниже приведён минимальный пример:

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

    // Additional origins allowed to open WebSocket connections and call
    // /wails/runtime. Same-origin requests are always allowed.
    WebSocketOriginPatterns: []string{"app.example.com"},

    // Disable WebSocket and runtime origin checks. Unsafe; default: false.
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

Она подходит для следующих задач:

- Проверки liveness и readiness в Kubernetes
- Проверки работоспособности балансировщика нагрузки
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

- **Из клиентской части в серверную**: события, отправленные из браузера, передаются по HTTP и поступают обработчикам событий Go
- **Из серверной части в клиентскую**: события, отправленные из Go, транслируются по WebSocket всем подключённым браузерам

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
4. Завершает работу сервисов в обратном порядке

## Отличия от настольного режима

| Возможность | Настольный режим | Серверный режим |
| --- | --- | --- |
| Нативные окна | Создаются | Окна браузера (`browser-N`) |
| Системный трей | Доступен | Недоступен |
| Нативные диалоговые окна | Доступны | Недоступны |
| Меню приложения | Доступно | Недоступно |
| Сведения об экране | Доступны | Нативного экрана нет; `GetPrimary()` возвращает `nil` |
| Привязки сервисов | Работают | Работают |
| События | Работают | Работают (через WebSocket) |
| Ресурсы | Через WebView | Через HTTP |
| Требуется CGO | Да | Нет |

### Поведение API окон

В серверном режиме API, связанные с окнами, обрабатываются безопасно:

- `app.Window.NewWithOptions(options)` — возвращает ненулевой объект окна Go; нативное окно рабочего стола при этом не создаётся
- `app.Hide()` / `app.Show()` — не выполняют никаких действий
- `app.Screen.GetPrimary()` — возвращает `*application.Screen`, равное `nil`, если основной экран недоступен; проверяйте значение перед разыменованием

Различайте объект окна Go и наличие нативного окна. При сборке сервера защищайте код, предназначенный только для настольного режима: ненулевой объект окна Go не означает, что существует нативное окно или экран.

### Поддерживаемые платформы

Серверный режим полностью реализован на Go и без CGO компилируется для всех поддерживаемых Wails значений `GOOS`, включая `android` и `ios`. Для мобильных целевых платформ он создаёт только сервер без графического интерфейса:

- Нативные мобильные API `application.Android` и `application.IOS` недоступны. Код, который их использует, необходимо исключить из серверных сборок с помощью ограничения сборки `!server`.
- По умолчанию используется средство ведения журнала для настольного режима, которое в сборках для разработки записывает данные в stderr, а не в logcat или консоль WebView.

## Сборка для рабочей среды

### Использование Task (рекомендуется)

Проекты, созданные с помощью `wails3 init`, содержат задачу `build:server`:

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

Проекты Wails содержат готовую конфигурацию Docker. Чтобы собрать и запустить приложение в контейнере:

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

Созданный файл `Dockerfile.server` формирует минимальный образ на базе distroless. Привязка к сетевому интерфейсу настраивается автоматически, поэтому приложение будет доступно извне контейнера.

### Docker Compose

Для более сложных развёртываний ниже приведена конфигурация Docker Compose с проверками работоспособности:

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
В примере проверки работоспособности используется `wget`. Если базовый образ не содержит дистрибутив, необходимо либо включить в образ исполняемый файл для проверки работоспособности, либо использовать внешний механизм такой проверки (например, параметр Docker `curl` или вспомогательный контейнер).

@end

### Собственный Dockerfile

Если требуется больше контроля, можно создать собственный Dockerfile. Главное — задать `WAILS_SERVER_HOST=0.0.0.0`, чтобы сервер принимал подключения извне контейнера:

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

1. **По умолчанию привязывайтесь к localhost** — используйте `0.0.0.0` только при необходимости
2. **Используйте TLS в рабочей среде** — настройте `ServerOptions.TLS`
3. **Размещайте приложение за обратным прокси-сервером** — используйте nginx/traefik для дополнительной защиты
4. **Используйте один источник для WebSocket и вызовов среды выполнения** — добавляйте с помощью `WebSocketOriginPatterns` только доверенные источники; избегайте `WebSocketAllowAllOrigins`
5. **Проверяйте все входные данные** — применяйте те же меры безопасности, что и для любого веб-приложения

## Пример

Полный пример доступен по адресу `v3/examples/server/`:

```bash
cd v3/examples/server

# Using Taskfile
task dev

# Or run directly
go run -tags server .

# Open http://localhost:8080 in browser
```

## Переменные среды

Для сценариев развёртывания, в которых конфигурацию сервера требуется переопределять без изменения кода, Wails распознаёт следующие переменные среды:

| Переменная | Описание | Значение по умолчанию |
| --- | --- | --- |
| `WAILS_SERVER_HOST` | Сетевой интерфейс для привязки | `localhost` |
| `WAILS_SERVER_PORT` | Порт для прослушивания | `8080` |

Эти переменные имеют приоритет над `ServerOptions` в коде. Поэтому в примерах Docker задаётся `WAILS_SERVER_HOST=0.0.0.0`: это позволяет контейнеру принимать внешние подключения без каких-либо изменений приложения.

## См. также

- [Собственный транспорт](/guides/custom-transport/) — расширенная настройка IPC
- [Сервисы](/features/bindings/services/) — документация по привязке сервисов
- [События](/guides/events-reference/) — документация по системе событий

### Размер запроса к среде выполнения

Размер запросов к `/wails/runtime` до обработки JSON ограничен 64 МиБ. На обычный запрос, превышающий это ограничение, возвращается HTTP 413, в том числе если в запросе отсутствует `Content-Length`. Для фрагментированной загрузки в среду выполнения сохраняются отдельные ограничения: 1 МиБ на фрагмент и 64 МиБ на собранные данные. При необходимости установите более низкое ограничение с помощью промежуточного ПО приложения или обратного прокси-сервера.

### Источник запроса к среде выполнения

`/wails/runtime` принимает только запросы `POST`; для других методов возвращается HTTP 405. Для запроса браузера, в котором `Origin` не совпадает ни с хостом запроса, ни с `WebSocketOriginPatterns`, возвращается HTTP 403. Такой же ответ возвращается для запроса без `Origin`, если его `Sec-Fetch-Site` имеет значение `cross-site` или `same-site`. Запросы без обоих заголовков, например от `curl` или других клиентов, не являющихся браузерами, принимаются. Если обратный прокси-сервер перезаписывает заголовок `Host`, добавьте общедоступный источник в `WebSocketOriginPatterns`.

`WebSocketOriginPatterns` предназначен для развёртываний с обратным прокси-сервером, в которых общедоступный источник отличается от `Host`, видимого приложению. Интерфейсы, обслуживаемые из другого источника и напрямую вызывающие `/wails/runtime`, не поддерживаются: конечная точка не отправляет заголовки CORS и отвечает на предварительные запросы `OPTIONS` кодом HTTP 405. Обслуживайте интерфейс с сервера Wails либо разместите их оба за одним источником прокси-сервера.

### Проверка Host для кольцевого интерфейса

Когда сервер привязан к кольцевому адресу (к используемому по умолчанию `localhost`, к любому адресу `127.0.0.0/8` или к `::1`, независимо от того, задан ли он в `ServerOptions.Host` или `WAILS_SERVER_HOST`), каждый запрос должен содержать заголовок `Host`, указывающий `localhost`, IP-адрес в буквенной записи, хост привязки или хост, совпадающий с частью хоста записи `WebSocketOriginPatterns`; при любом другом значении `Host` возвращается HTTP 403. Это предотвращает перепривязку DNS, при которой вредоносная веб-страница направляет собственное имя хоста на `127.0.0.1`, чтобы её запросы проходили описанные выше проверки одинакового источника. Порт не сравнивается, а варианты записи наподобие `localhost.` (с точкой в конце) или `127.1` отклоняются. Проверка продолжает действовать и при заданном `WebSocketAllowAllOrigins`. Если локальный обратный прокси-сервер передаёт свой общедоступный `Host`, добавьте этот хост в `WebSocketOriginPatterns`. Привязки к другим адресам, например к `0.0.0.0` в контейнерах, не затрагиваются.
