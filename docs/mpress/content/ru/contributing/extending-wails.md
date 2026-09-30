---
title: "Расширение Wails"
description: "Практическое руководство по добавлению новых возможностей и платформ в Wails v3"
slug: "contributing/extending-wails"
sourcePath: "contributing/extending-wails.md"
---

> Wails спроектирован так, чтобы его было **легко модифицировать**.
>
> Каждая основная подсистема реализована на Go: её код можно изучать, изменять и выпускать.
>
> На этой странице показано, *с чего* начать и *как* сохранить кроссплатформенность, если вы хотите:

- Добавить **сервис** (уведомления, хранилище ключей и значений, собственный IPC и т. д.)
- Создать **новую команду CLI** (`wails3 <foo>`)
- Расширить **среду выполнения** (API окон, диалоговые окна, события)
- Добавить **возможность платформы** (Wayland и т. д.)
- Сохранить **кроссплатформенную совместимость**, не утонув в тегах `//go:build`

---

## 1. Добавление сервиса

«Сервис» в v3 — это предоставленный пользователем тип Go, который регистрируется через `application.Options.Services` и становится доступен из JS посредством сгенерированных привязок. В кодовую базу v3 входят:

- `internal/service/` — каркас для `wails3 generate service`:
  ```
  internal/service/
  ├── service.go              # Install(options *flags.ServiceInit)
  └── template/
      ├── README.tmpl.md
      ├── go.mod.tmpl
      ├── service.go.tmpl
      └── service.tmpl.yml
  ```

- `pkg/services/` — готовые сервисы, которые можно зарегистрировать уже сейчас (notifications, kvstore, sqlite, log, fileserver, dock и т. д.).

Файлы генератора и CLI, которые в старых черновиках упоминались как `internal/service/template/template.go` и `internal/generator/collect/services.go`, не существуют: генератор каркаса — это `internal/service/service.go` (точка входа — `service.Install`), а метаданные привязок для сервисов собираются в `internal/generator/collect/service.go`.

### 1.1 Определение сервиса

```go
package chat

type Service struct {
    messages []string
}

func New() *Service { return &Service{} }

func (s *Service) Send(msg string) string {
    s.messages = append(s.messages, msg)
    return "ok"
}
```

### 1.2 Реализация интерфейсов жизненного цикла (необязательно)

Сервис может при необходимости реализовывать следующие интерфейсы (из `pkg/application`):

```go
func (s *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error { return nil }
func (s *Service) ServiceShutdown() error                                                       { return nil }
```

> **Важно:** `ServiceShutdown` **не принимает аргументов**. Метод с
>
> сигнатурой `ServiceShutdown(ctx context.Context) error` **не** реализует
>
> этот интерфейс и поэтому никогда не будет вызван, причём без каких-либо предупреждений.

### 1.3 Регистрация сервиса в приложении

Глобального вызова `services.Register(...)` нет. Сервисы регистрируются во время выполнения через `application.Options.Services`:

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(chat.New()),
    },
})
```

После регистрации `wails3 generate bindings` создаёт в каталоге `frontend/bindings/<your import path>/...` модули ES, предоставляющие обёртки для экспортированных методов.

### 1.4 Вызов из JS

```js
import { Send } from "../bindings/github.com/you/yourapp/chat";

await Send("hi");
```

В v3 нет глобального объекта `window.backend.*`: вызовы выполняются через сгенерированные модули ES, которые, в свою очередь, вызывают `Call.ByID(...)` из `/wails/runtime.js`.

---

## 2. Создание новой команды CLI

CLI в v3 использует **`github.com/leaanthony/clir`**, а не cobra. Подключение команд находится в `v3/cmd/wails3/main.go`:

```go
import "github.com/leaanthony/clir"

func main() {
    app := clir.NewCli("wails", "The Wails3 CLI", "v3")
    app.NewSubCommand("hello", "Prints Hello Wails").Action(func() error {
        fmt.Println("Hello Wails!")
        return nil
    })
    // ... other subcommands explicitly wired here
    _ = app.Run()
}
```

Автоматической регистрации на основе `init()` нет. Добавьте новую подкоманду в `cmd/wails3/main.go` и реализующую её функцию в `internal/commands/` (а если команда принимает параметры — также структуру флагов в `internal/flags/`). Пересоберите CLI:

```
cd v3
go install ./cmd/wails3
wails3 hello
```

Если вашей команде нужна интеграция с Taskfile, повторно используйте вспомогательные функции из `internal/commands/task_wrapper.go` (`wrapTask("yourtask", args)`).

---

## 3. Изменение среды выполнения

Типичные причины:

- Новая возможность окна (`SetOpacity`, `Shake` и т. д.)
- Дополнительное диалоговое окно (`ColorPicker`)
- API системного уровня (яркость экрана)

### 3.1 Публичный API

Добавьте метод в `pkg/application/webview_window.go` (интерфейс находится в `window.go`):

```go
func (w *WebviewWindow) SetOpacity(o float32) Window {
    InvokeSync(func() { w.impl.setOpacity(o) })
    return w
}
```

Используйте существующие вспомогательные функции `InvokeSync`/`InvokeAsync`, чтобы вызов выполнялся в главном потоке.

### 3.2 Обработчик сообщений

Если JS должен вызывать новый метод, расширьте соответствующий файл `pkg/application/messageprocessor_*.go`. Обработчик сообщений использует методы с конструкцией switch у `MessageProcessor`, а не глобальный вызов `register(...)`:

```go
// inside messageprocessor_window.go
case "setOpacity":
    var args struct {
        WindowID uint    `json:"windowID"`
        Opacity  float32 `json:"opacity"`
    }
    if err := json.Unmarshal(payload, &args); err != nil { ... }
    window, _ := m.app.Window.GetByID(args.WindowID)
    window.SetOpacity(args.Opacity)
```

В кодовой базе **нет** файла `messageprocessor_window_opacity.go` или основанного на `init()` шаблона `register(MsgSetOpacity, ...)`.

### 3.3 Реализация для платформы

Добавьте реализацию в каждый файл для отдельной ОС в `pkg/application/`:

```
pkg/application/
├── webview_window_darwin.go   //go:build darwin
├── webview_window_linux.go    //go:build linux
└── webview_window_windows.go  //go:build windows
```

Если платформа не поддерживает эту возможность, реализуйте пустую заглушку. Во фреймворке нет специального значения `ErrCapability`: укажите поддержку в документации и при необходимости — в соответствующем логическом поле `Options` или структуры параметров конкретной платформы.

### 3.4 Флаг возможности (необязательно)

Пакет `internal/capabilities/` предназначен для объявления наборов возможностей для отдельных платформ. Публичного API `application.HasCapability` / `application.CapOpacity` нет. Если возможность должна проверяться во время выполнения, добавьте её в `internal/capabilities/` и предоставьте типизированный метод получения из `pkg/application`.

---

## 4. Добавление новых возможностей платформы

Пример: необязательная поддержка Wayland в Linux.

1. Разделите соответствующий файл `pkg/application/*_linux.go` на `*_linux_x11.go` (`//go:build linux && !wayland`) и `*_linux_wayland.go` (`//go:build linux && wayland`).
2. Пользователи должны включать эту возможность с помощью `wails3 build --tags wayland`. Передавайте дополнительные теги через существующий механизм `EXTRA_TAGS` в `internal/commands/task_wrapper.go`. На уровне `dev` нет флага `--tags wayland`: `wails3 dev` принимает только `--config`, `--port` и `-s`.
3. Обновите документацию и все относящиеся к конкретным платформам файлы README в `pkg/application/`.

> Сведите набор тегов сборки по умолчанию к минимуму; используйте подключаемые по запросу теги только для узкоспециализированных возможностей.

---

## 5. Контрольный список кроссплатформенной совместимости

| ✅ Шаг | Зачем |
| --- | --- |
| Реализуйте **каждый** публичный метод во всех файлах для отдельных платформ (даже если это заглушка) | Обеспечивает успешную сборку в каждой ОС |
| Документируйте корректную деградацию функциональности для каждой ОС | Приложения смогут выполнять ветвление по `runtime.GOOS` без скрытых ошибок |
| Сначала используйте **чистый Go**, а Cgo — только при необходимости | Упрощает кросс-компиляцию (в Linux издержки Cgo уже неизбежны) |
| Запустите `task test:cli`, `task test:generator` и `task test:templates` | Воспроизводит среду CI локально |
| Документируйте новые теги сборки в документации для участников проекта или в README шаблона проекта | Пользователи должны знать о функциях, требующих явного включения |

---

## 6. Отладочные сборки и скорость итераций

- Используйте `Options.LogLevel = slog.LevelDebug` (`Options.Logger = slog.Default()`), чтобы выводить подробные сведения о работе среды выполнения. Переменной окружения `WAILS_LOG_LEVEL` не существует.
- Доступные флаги `wails3 dev`: `--config`, `--port` и `-s`. Флагов `-race` и `-verbose` не существует — запускайте детектор гонок с помощью `go test -race ./...` либо выполните для приложения `go build -race` и запустите его напрямую.
- Руководство по тестированию с детектором гонок и Cgo находится в `v3/TESTING.md` (в старых черновиках указывался несуществующий файл `pkg/application/RACE.md`).

---

## 7. Вклад в основной проект

1. Для новой функциональности или изменения публичного поведения откройте черновой PR с **WEP (предложением по улучшению Wails)**, чтобы обсудить идею и дизайн. Создавайте issue только для воспроизводимой ошибки или проблемы в документации.
2. При реализации следуйте описанным выше подходам.
3. Добавьте:
  - Модульные тесты (`*_test.go`)
  - Документацию (этот файл или соответствующую страницу `docs/...`)
  - Регрессионный тест в `internal/generator/testcases/`, если вы изменили генератор привязок

4. Перед отправкой изменений запустите локально `task precommit` и соответствующие цели `task test:*`.

---

### Быстрые ссылки

| Область | Расположение |
| --- | --- |
| Встроенные сервисы | `pkg/services/` |
| Генератор каркаса сервиса | `internal/service/` |
| Подключение CLI | `v3/cmd/wails3/main.go` |
| Реализации команд CLI | `internal/commands/` |
| Среда выполнения для отдельных ОС | `pkg/application/*_{darwin,linux,windows}.go` |
| Объявления возможностей | `internal/capabilities/` |
| DSL файла Taskfile | `v3/Taskfile.yaml` |
| Генератор констант событий | `v3/tasks/events/generate.go` |

---

Теперь у вас есть **план действий**, который поможет настроить Wails под свои задачи: добавляйте сервисы, творите магию с CLI, изменяйте среду выполнения или реализуйте совершенно новые возможности ОС. Успешного расширения!
