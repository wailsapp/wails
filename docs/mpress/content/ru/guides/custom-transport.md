---
title: "Создание пользовательского транспортного уровня"
description: "Узнайте, как создать и настроить собственный пользовательский транспортный уровень IPC в Wails v3"
slug: "guides/custom-transport"
sourcePath: "guides/custom-transport.md"
---

Wails v3 позволяет использовать пользовательский транспортный уровень IPC, сохраняя все сгенерированные привязки и обмен событиями. Это позволяет заменить стандартный транспорт на основе HTTP fetch транспортом WebSocket, пользовательскими протоколами или любым другим транспортным механизмом.

## Обзор

По умолчанию для связи фронтенда с бэкендом через `/wails/runtime` Wails использует HTTP-запросы fetch. API пользовательского транспорта позволяет:

- Заменить HTTP-транспорт на WebSocket, gRPC или любой пользовательский протокол
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

## Использование

### 1. Реализуйте интерфейс транспорта

Создайте пользовательский транспорт, реализовав интерфейс `Transport`:

```go
package main

import (
    "context"
    "github.com/wailsapp/wails/v3/pkg/application"
)

type MyCustomTransport struct {
    // Your fields
}

func (t *MyCustomTransport) Start(ctx context.Context, processor *application.MessageProcessor) error {
    // Initialize your transport (WebSocket server, gRPC server, etc.)
    // When you receive requests, call processor.HandleRuntimeCallWithIDs()
    return nil
}

func (t *MyCustomTransport) Stop() error {
    // Clean up your transport
    return nil
}
```

### 2. Настройте приложение

Передайте пользовательский транспорт в параметры приложения:

```go
func main() {
    app := application.New(application.Options{
        Name: "My App",
        Transport: &MyCustomTransport{},
        // ... other options
    })

    err := app.Run()
    if err != nil {
        log.Fatal(err)
    }
}
```

### 3. Измените среду выполнения фронтенда

При использовании пользовательского транспорта необходимо изменить среду выполнения фронтенда, чтобы она использовала ваш транспорт вместо HTTP fetch. Реализуйте интерфейс `RuntimeTransport`, который будет использоваться для обработки запросов:

```typescript
const { setTransport } = await import('/wails/runtime.js');

class MyRuntimeTransport {
  call(objectID: number, method: number, windowName: string, args: any): Promise<any> {
    // TODO: implement IPC call with your transport protocol

    return resp;
  }
}

const myTransport = new MyRuntimeTransport();
setTransport(myTransport);
```

## Примечания

- Если пользовательский транспорт не указан, стандартный HTTP-транспорт продолжает работать
- Сгенерированные привязки остаются неизменными — меняется только транспортный уровень
- События, диалоговые окна, буфер обмена и все остальные возможности Wails работают прозрачно
- Вы отвечаете за обработку ошибок, логику повторного подключения и безопасность пользовательского транспорта
- Предоставленный пример с WebSocket предназначен для демонстрации и может потребовать усиления защиты перед использованием в рабочей среде

## Справочник API

### Интерфейс транспорта

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

Для развёртывания в браузере или передачи ресурсов и IPC через пользовательский транспорт реализуйте интерфейс `AssetServerTransport`:

```go
type AssetServerTransport interface {
    Transport

    // ServeAssets configures the transport to serve assets alongside IPC.
    // The assetHandler is Wails' internal asset server that handles all assets,
    // runtime.js, capabilities, flags, etc.
    ServeAssets(assetHandler http.Handler) error
}
```

**Когда следует реализовать этот интерфейс:**

- Запуск приложения в браузере вместо веб-представления
- Передача ресурсов по HTTP наряду с пользовательским транспортом IPC
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
- `/wails/runtime.js` — библиотека среды выполнения Wails

## См. также

- `transport.go` — основные интерфейсы и типы транспорта
- `messageprocessor.go` — базовый обработчик сообщений, обрабатывающий весь IPC-трафик Wails
