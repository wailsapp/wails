---
title: "Создание собственного транспортного уровня"
description: "Узнайте, как создать и настроить собственный транспортный уровень IPC для Wails v3"
slug: "guides/custom-transport"
sourcePath: "guides/custom-transport.md"
---

Wails v3 позволяет предоставить собственный транспортный уровень IPC, сохранив все сгенерированные привязки и обмен событиями. Благодаря этому стандартный транспорт на основе HTTP Fetch можно заменить протоколом WebSocket, собственными протоколами или любым другим транспортным механизмом.

## Обзор

По умолчанию для обмена данными между фронтендом и бэкендом через `/wails/runtime` Wails использует HTTP-запросы Fetch, отправляемые из фронтенда. API собственного транспорта позволяет:

- Заменить HTTP-транспорт на WebSocket, gRPC или любой собственный протокол
- Сохранить полную совместимость с генерацией кода Wails
- Сохранить все существующие привязки, события, диалоговые окна и другие возможности Wails
- Реализовать собственное управление подключениями, аутентификацию и обработку ошибок

## Архитектура

```text
┌─────────────────────────────────────────────────┐
│  Frontend (TypeScript)                          │
│  - Generated bindings still work                │
│  - Your custom client transport                 │
└──────────────────┬──────────────────────────────┘
                   │
                   │ Your Protocol (WebSocket/etc)
                   │
┌──────────────────▼──────────────────────────────┐
│  Backend (Go)                                   │
│  - Your Transport implementation                │
│  - Wails MessageProcessor                       │
│  - All existing Wails infrastructure            │
└─────────────────────────────────────────────────┘
```

## Запуск полного примера с WebSocket

[Пример транспорта на основе WebSocket](https://github.com/wailsapp/wails/tree/master/v3/examples/websocket-transport) включает обе стороны протокола: идентификаторы запросов и ответы, управление жизненным циклом, клиент JavaScript и доставку событий. Используйте эти файлы как отправную точку вместо транспорта, реализующего только `Start` и `Stop`.

В рабочей копии репозитория Wails с установленными нативными зависимостями платформы выполните:

```bash
cd v3/examples/websocket-transport
go run .
```

Пример открывает нативное окно и запускает конечную точку WebSocket на порту 9099. Используйте его в доверенной локальной среде разработки. Перед развёртыванием варианта, доступного по сети, ограничьте прослушиваемые адреса, проверяйте источники запросов и аутентифицируйте клиентов. Политика проверки источников в примере допускает любые источники и не является политикой безопасности для рабочей среды.

### Контракт бэкенда

Транспорт должен реализовывать все три обязательных метода. В примере `WebSocketTransport` реализует `Start`, `JSClient` и `Stop`, а `JSClient` предоставляет JavaScript-адаптер, который Wails внедряет в страницу. Добавьте в свою реализацию проверку соответствия интерфейсу на этапе компиляции:

```go
var _ application.Transport = (*WebSocketTransport)(nil)
```

Настройте приложение с помощью конструктора из примера:

```go
Transport: NewWebSocketTransport("127.0.0.1:9099"),
```

Это поле должно находиться в `application.Options` рядом с обычными ресурсами и сервисами. Скопируйте как файлы реализации, так и вызов конструктора. `Start` получает `MessageProcessor`; `Stop` должен закрывать прослушиватели, активные соединения и незавершённые операции.

### Контракт фронтенда

Файл `websocket-transport.js` из примера реализует транспорт среды выполнения, включая сопоставление запросов с ответами и отклонение завершившихся ошибкой вызовов. После того как этот файл станет доступен через ваши ресурсы, инициализируйте его перед вызовом сгенерированных привязок сервисов:

```javascript
import { setTransport } from '/wails/runtime.js'
import { createWebSocketTransport } from '/websocket-transport.js'

setTransport(createWebSocketTransport('ws://127.0.0.1:9099/wails/ws'))
```

Сгенерированные привязки сервисов остаются без изменений. Если ваш транспорт заменяет клиент, он должен поддерживать методы, используемые приложением, и сохранять их семантику ошибок. Для доставки событий требуется необязательный контракт транспорта событий; произвольная реализация `Transport` не предоставляет его автоматически. Это показано в полном примере с WebSocket.

При удалённом развёртывании используйте TLS, а также явно заданные политики аутентификации и проверки источников. Если приложение размещается только для браузеров, необходимо дополнительно предоставлять его ресурсы, как описано ниже.

## Справочник по API

### Интерфейс Transport

```go
type Transport interface {
    Start(ctx context.Context, messageProcessor *application.MessageProcessor) error
    // JSClient returns the JavaScript shim that the runtime injects into the
    // window so that frontend code can call into the transport.
    JSClient() []byte
    Stop() error
}
```

### Интерфейс AssetServerTransport (необязательно)

Для развёртывания в браузере или предоставления ресурсов и IPC через собственный транспорт реализуйте интерфейс `AssetServerTransport`:

```go
type AssetServerTransport interface {
    Transport

    // ServeAssets configures the transport to serve assets alongside IPC.
    // The assetHandler is Wails' internal asset server that handles all assets,
    // runtime.js, capabilities, flags, etc.
    ServeAssets(assetHandler http.Handler) error
}
```

**Когда следует реализовывать этот интерфейс:**

- Запуск приложения в браузере вместо WebView
- Предоставление ресурсов по HTTP вместе с собственным транспортом IPC
- Создание приложений, доступных по сети

**Пример реализации:**

```go
func (t *MyTransport) ServeAssets(assetHandler http.Handler) error {
    mux := http.NewServeMux()

    // Mount your IPC endpoint
    mux.HandleFunc("/my/ipc/endpoint", t.handleIPC)

    // Mount Wails asset server for everything else
    mux.Handle("/", assetHandler)

    // Start HTTP server
    t.httpServer.Handler = mux
    go t.httpServer.ListenAndServe()

    return nil
}
```

При вызове `ServeAssets()` обработчик assetHandler предоставляет:

- Все статические ресурсы (HTML, CSS, JS, изображения и т. д.)
- `/wails/runtime.js` — библиотеку среды выполнения Wails

## См. также

- `transport.go` — основные транспортные интерфейсы и типы
- `messageprocessor.go` — базовый обработчик сообщений, обрабатывающий весь IPC-трафик Wails
