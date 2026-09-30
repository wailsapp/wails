---
title: "Жизненный цикл приложения"
description: "Жизненный цикл приложения Wails от запуска до завершения работы"
slug: "concepts/lifecycle"
sourcePath: "concepts/lifecycle.md"
---

## Жизненный цикл приложения

У настольных приложений есть жизненный цикл от запуска до завершения работы. Для эффективного управления этим жизненным циклом Wails v3 предоставляет **сервисы**, **события** и **хуки**.

## Этапы жизненного цикла

```d2
direction: down

Start: Запуск приложения {
  shape: oval
  style.fill: "#10B981"
}

Init: Инициализация {
  Parse: Разбор параметров {
    shape: rectangle
  }
  Register: Регистрация сервисов {
    shape: rectangle
  }
  Setup: Настройка среды выполнения {
    shape: rectangle
  }
}

AppRun: app.Run() {
  shape: rectangle
  style.fill: "#3B82F6"
}

ServiceStartup: Запуск сервисов {
  shape: rectangle
  style.fill: "#8B5CF6"
}

EventLoop: Цикл обработки событий {
  Process: Обработка событий {
    shape: rectangle
  }
  Handle: Обработка сообщений {
    shape: rectangle
  }
  Update: Обновление интерфейса {
    shape: rectangle
  }
}

QuitSignal: Сигнал завершения {
  shape: diamond
  style.fill: "#F59E0B"
}

ShouldQuit: Проверка ShouldQuit {
  shape: rectangle
  style.fill: "#3B82F6"
}

OnShutdown: Обратные вызовы OnShutdown {
  shape: rectangle
  style.fill: "#3B82F6"
}

ServiceShutdown: Остановка сервисов {
  shape: rectangle
  style.fill: "#8B5CF6"
}

Cleanup: Очистка {
  Close: Закрытие окон {
    shape: rectangle
  }
  Release: Освобождение ресурсов {
    shape: rectangle
  }
}

End: Завершение приложения {
  shape: oval
  style.fill: "#EF4444"
}

Start -> Init.Parse
Init.Parse -> Init.Register
Init.Register -> Init.Setup
Init.Setup -> AppRun
AppRun -> ServiceStartup
ServiceStartup -> EventLoop.Process
EventLoop.Process -> EventLoop.Handle
EventLoop.Handle -> EventLoop.Update
EventLoop.Update -> EventLoop.Process: Цикл
EventLoop.Process -> QuitSignal: Пользователь завершает работу
QuitSignal -> ShouldQuit: Завершение разрешено?
ShouldQuit -> EventLoop.Process: Запрещено
ShouldQuit -> OnShutdown: Разрешено
OnShutdown -> ServiceShutdown
ServiceShutdown -> Cleanup.Close
Cleanup.Close -> Cleanup.Release
Cleanup.Release -> End
```

### 1. Создание приложения

Создайте приложение с помощью `application.New()`:

```go
app := application.New(application.Options{
    Name:        "My App",
    Description: "An application built with Wails",
    Services: []application.Service{
        application.NewService(&MyService{}),
    },
    Assets: application.AssetOptions{
        Handler: application.BundledAssetFileServer(assets),
    },
})
```

**Что происходит:**

1. Параметры разбираются и проверяются
2. Сервисы регистрируются (но пока не запускаются)
3. Настраивается сервер ресурсов
4. Настраивается среда выполнения

### 2. Запуск приложения

Чтобы запустить приложение, вызовите `app.Run()`:

```go
err := app.Run()  // Blocks until quit
if err != nil {
    log.Fatal(err)
}
```

**Что происходит:**

1. Сервисы запускаются в порядке регистрации
2. Активируются обработчики событий
3. Можно создавать окна
4. Запускается цикл событий

### 3. Цикл событий

Приложение переходит в цикл событий, в котором проводит большую часть времени:

- Обрабатываются события ОС (мышь, клавиатура, события окон)
- Обрабатываются сообщения из Go в JS
- Выполняются вызовы из JS в Go
- Отрисовываются обновления интерфейса

### 4. Завершение работы

При завершении работы приложения:

1. Проверяется функция обратного вызова `ShouldQuit` (если она задана)
2. Выполняются функции обратного вызова `OnShutdown`
3. Работа сервисов завершается в обратном порядке
4. Окна закрываются
5. Ресурсы освобождаются

## Жизненный цикл сервисов

Сервисы — основной механизм управления жизненным циклом в Wails v3. Через интерфейсы они предоставляют хуки запуска и завершения работы. Полную документацию по сервисам см. в [руководстве по сервисам](/features/bindings/services/).

### Создание сервиса

```go
type MyService struct {
    db *sql.DB
}

// ServiceStartup is called when the application starts
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    var err error
    s.db, err = sql.Open("sqlite3", "app.db")
    if err != nil {
        return err  // Startup aborts if error returned
    }

    // Run migrations
    if err := s.runMigrations(); err != nil {
        return err
    }

    return nil
}

// ServiceShutdown is called when the application shuts down
func (s *MyService) ServiceShutdown() error {
    if s.db != nil {
        return s.db.Close()
    }
    return nil
}
```

### Регистрация сервисов

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&MyService{}),
        application.NewService(&AnotherService{}),
    },
})
```

**Основные моменты:**

- Сервисы запускаются в порядке регистрации
- Работа сервисов завершается в **обратном** порядке регистрации
- Если метод `ServiceStartup` сервиса возвращает ошибку, приложение аварийно завершает запуск
- Контекст `ctx`, переданный в `ServiceStartup`, отменяется при начале завершения работы

### Использование контекста приложения

Контекст, переданный в `ServiceStartup`, действителен в течение всего жизненного цикла приложения:

```go
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    // Start a background task that respects shutdown
    go func() {
        ticker := time.NewTicker(5 * time.Minute)
        defer ticker.Stop()

        for {
            select {
            case <-ticker.C:
                s.performBackgroundSync()
            case <-ctx.Done():
                // Application is shutting down
                return
            }
        }
    }()

    return nil
}
```

Контекст также можно получить из экземпляра приложения:

```go
app := application.Get()
ctx := app.Context()
```

## Хуки уровня приложения

Это вспомогательные функции обратного вызова в `application.Options`, которые позволяют подключаться к жизненному циклу приложения без создания полноценного сервиса. Они полезны для простых задач очистки, подтверждения выхода или выполнения кода в определённые моменты последовательности завершения работы.

Для более сложного управления жизненным циклом, включающего логику запуска, внедрение зависимостей или ресурсы с состоянием, вместо этого используйте [сервисы](#---2).

### ShouldQuit

Функция обратного вызова `ShouldQuit` вызывается при каждом запросе на выход — независимо от того, закрыл ли пользователь последнее окно, нажал Cmd+Q (macOS) / Alt+F4 (Windows) или был программно вызван `app.Quit()`.

**Возвращаемое значение:**

- Верните `true`, чтобы разрешить выход (работа приложения будет завершена)
- Верните `false`, чтобы отменить выход (приложение продолжит работу)

Здесь можно перехватить запрос на выход и при необходимости отменить его, например чтобы предупредить пользователя о несохранённых изменениях:

```go
app := application.New(application.Options{
    ShouldQuit: func() bool {
        if !hasUnsavedChanges() {
            return true // No unsaved changes, allow quit
        }

        // Prompt the user — MessageDialog.Show() blocks and returns nothing.
        // The button's OnClick callback fires for whichever button the user picks.
        shouldQuit := false

        dlg := application.Get().Dialog.Question().
            SetTitle("Unsaved Changes").
            SetMessage("You have unsaved changes. Quit anyway?")

        quit := dlg.AddButton("Quit")
        cancel := dlg.AddButton("Cancel")
        dlg.SetDefaultButton(cancel)
        dlg.SetCancelButton(cancel)

        quit.OnClick(func() { shouldQuit = true })

        dlg.Show()
        return shouldQuit
    },
})
```

Если `ShouldQuit` не задан, приложение немедленно завершит работу при поступлении запроса на выход.

**Когда вызывается ShouldQuit:**

- Пользователь закрывает последнее окно (если не задан `DisableQuitOnLastWindowClosed`)
- Пользователь нажимает Cmd+Q в macOS
- Пользователь нажимает Alt+F4 в Windows (когда последнее окно находится в фокусе)
- Код вызывает `app.Quit()`

**Когда ShouldQuit НЕ вызывается:**

- Процесс принудительно завершается (SIGKILL или принудительное завершение через Диспетчер задач)
- `os.Exit()` вызывается напрямую

### OnShutdown

Функция обратного вызова `OnShutdown` вызывается после подтверждения выхода из приложения (после того как `ShouldQuit` возвращает `true`, если эта функция задана). Используйте её для таких задач очистки, как сохранение состояния, закрытие подключений к базе данных или освобождение ресурсов.

```go
app := application.New(application.Options{
    OnShutdown: func() {
        // Save application state
        saveState()

        // Close connections
        cleanup()
    },
})
```

Дополнительные функции обратного вызова для завершения работы также можно программно регистрировать в любой момент жизненного цикла приложения:

```go
app.OnShutdown(func() {
    log.Println("Application shutting down...")
})
```

Несколько функций обратного вызова выполняются в порядке их регистрации. Процесс завершения работы блокируется до тех пор, пока не завершатся все функции обратного вызова.

**Важно:** функции обратного вызова при завершении работы должны выполняться быстро (менее 1 секунды). Операционная система может принудительно завершить приложения, которым требуется слишком много времени для выхода. Это может прервать очистку и привести к потере данных.

### PostShutdown

Функция обратного вызова `PostShutdown` вызывается после выполнения всех задач завершения работы, непосредственно перед завершением процесса. На этом этапе экземпляр приложения уже нельзя использовать: окна закрыты, службы остановлены, а ресурсы освобождены.

Это в первую очередь полезно для:

- Окончательного журналирования, которое должно выполняться после всех остальных операций очистки
- Тестирования и отладки поведения при завершении работы
- Платформ, на которых `app.Run()` не возвращает управление (функция обратного вызова гарантирует выполнение вашего кода)

```go
app := application.New(application.Options{
    PostShutdown: func() {
        // Final logging
        log.Println("Application terminated cleanly")

        // Flush any buffered logs
        logger.Sync()
    },
})
```

**Примечание:** не пытайтесь использовать функции приложения (окна, диалоговые окна и т. д.) в `PostShutdown` — они уже недоступны.

## Жизненный цикл на основе событий

Wails предоставляет систему событий, которая уведомляет вас о происходящем в приложении: открытии окон, запуске приложения, изменении темы и многом другом. Вы можете прослушивать эти события и реагировать на изменения жизненного цикла, не блокируя и не перехватывая их.

Для событий окна вместо `OnWindowEvent` также можно использовать `RegisterHook`, чтобы перехватывать и отменять действия — например, предотвращать закрытие окна. См. раздел [«Перехватчики окон»](#---) ниже.

Полную документацию по системе событий см. в [руководстве по событиям](/features/events/system/).

### События приложения

Прослушивайте события жизненного цикла приложения:

```go
app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(event *application.ApplicationEvent) {
    app.Logger.Info("Application has started!")
})
```

Также доступны события для отдельных платформ:

```go
// macOS
app.Event.OnApplicationEvent(events.Mac.ApplicationDidFinishLaunching, func(event *application.ApplicationEvent) {
    // Handle macOS launch
})

app.Event.OnApplicationEvent(events.Mac.ApplicationWillTerminate, func(event *application.ApplicationEvent) {
    // Handle macOS termination
})

// Windows
app.Event.OnApplicationEvent(events.Windows.ApplicationStarted, func(event *application.ApplicationEvent) {
    // Handle Windows start
})
```

### События окон

Прослушивайте события жизненного цикла окон:

```go
window := app.Window.New()

window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window gained focus")
})

window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
    app.Logger.Info("Window is closing")
})
```

### Перехватчики окон (отменяемые события)

Когда нужно **отменить** событие, используйте `RegisterHook` вместо `OnWindowEvent`:

```go
window := app.Window.New()

var countdown = 3

window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    countdown--
    if countdown > 0 {
        app.Logger.Info("Not closing yet!", "remaining", countdown)
        e.Cancel()  // Prevent the window from closing
        return
    }
    app.Logger.Info("Window closing now")
})
```

**Различие между OnWindowEvent и RegisterHook:**

- `OnWindowEvent`: уведомляет о наступлении события (отменить его нельзя)
- `RegisterHook`: позволяет перехватить и при необходимости отменить событие

## Жизненный цикл окна

У окон есть собственный жизненный цикл — от создания до уничтожения. Каждое окно загружает своё содержимое фронтенда независимо от других; его можно в любой момент показать, скрыть или закрыть. Когда пользователь пытается закрыть окно, это действие можно перехватить с помощью `RegisterHook`, чтобы запросить подтверждение или скрыть окно вместо его уничтожения.

Полную документацию по окнам см. в [руководстве по окнам](/features/windows/basics/).

```d2
direction: down

Create: Создание окна {
  shape: oval
  style.fill: "#10B981"
}

Load: Загрузка интерфейса {
  shape: rectangle
}

Show: Отображение окна {
  shape: rectangle
}

Active: Окно активно {
  Events: Обработка событий {
    shape: rectangle
  }
}

CloseRequest: Запрос на закрытие {
  shape: diamond
  style.fill: "#F59E0B"
}

Hook: Хук WindowClosing {
  shape: rectangle
  style.fill: "#3B82F6"
}

Destroy: Уничтожение окна {
  shape: rectangle
}

End: Окно закрыто {
  shape: oval
  style.fill: "#EF4444"
}

Create -> Load
Load -> Show
Show -> Active.Events
Active.Events -> Active.Events: Цикл
Active.Events -> CloseRequest: Пользователь закрывает окно
CloseRequest -> Hook
Hook -> Active.Events: Отменено
Hook -> Destroy: Разрешено
Destroy -> End
```

### Создание окон

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "My Window",
    Width:  800,
    Height: 600,
})
```

### Предотвращение закрытия окна

```go
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    if !hasUnsavedChanges() {
        return
    }

    // MessageDialog.Show() returns nothing; per-button OnClick handlers fire.
    dlg := application.Get().Dialog.Question().
        SetTitle("Unsaved Changes").
        SetMessage("Save before closing?")

    save := dlg.AddButton("Save")
    discard := dlg.AddButton("Discard")
    cancel := dlg.AddButton("Cancel")
    dlg.SetDefaultButton(save)
    dlg.SetCancelButton(cancel)

    save.OnClick(func() { saveChanges() })
    cancel.OnClick(func() { e.Cancel() }) // Prevent close
    _ = discard                            // "Discard" falls through and allows close

    dlg.Show()
})
```

### Скрытие вместо закрытия

Распространённый шаблон для приложений в области уведомлений:

```go
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    window.Hide()  // Hide instead of destroy
    e.Cancel()     // Prevent actual close
})
```

## Жизненный цикл нескольких окон

При наличии нескольких окон:

```go
mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Main Window",
})

settingsWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "Settings",
    Width:  400,
    Height: 600,
    Hidden: true,  // Start hidden
})
```

**Поведение по умолчанию зависит от платформы:**

| Платформа | Поведение по умолчанию после закрытия последнего окна |
| --- | --- |
| macOS | Приложение продолжает работать (строка меню остаётся) |
| Windows | Приложение завершает работу |
| Linux | Приложение завершает работу |

macOS следует принятым на платформе соглашениям: приложения обычно остаются активными в строке меню, даже если не осталось ни одного окна. В Windows и Linux приложение по умолчанию завершает работу.

**Завершение работы на всех платформах после закрытия последнего окна:**

```go
app := application.New(application.Options{
    Mac: application.MacOptions{
        ApplicationShouldTerminateAfterLastWindowClosed: true,
    },
})
```

**Продолжение работы на всех платформах после закрытия последнего окна:**

Это полезно для приложений в области уведомлений и приложений, которые должны продолжать работать в фоновом режиме.

```go
app := application.New(application.Options{
    Windows: application.WindowsOptions{
        DisableQuitOnLastWindowClosed: true,
    },
    Linux: application.LinuxOptions{
        DisableQuitOnLastWindowClosed: true,
    },
})
```

## Распространённые шаблоны

### Шаблон 1: служба базы данных

```go
type DatabaseService struct {
    db *sql.DB
}

func (s *DatabaseService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    var err error
    s.db, err = sql.Open("sqlite3", "app.db")
    if err != nil {
        return fmt.Errorf("failed to open database: %w", err)
    }

    if err := s.db.PingContext(ctx); err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }

    return nil
}

func (s *DatabaseService) ServiceShutdown() error {
    if s.db != nil {
        return s.db.Close()
    }
    return nil
}

// Exported methods are available to the frontend
func (s *DatabaseService) GetUsers() ([]User, error) {
    // Query implementation
}
```

### Шаблон 2: служба конфигурации

```go
type ConfigService struct {
    config *Config
    path   string
}

func (s *ConfigService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    s.path = "config.json"

    data, err := os.ReadFile(s.path)
    if err != nil {
        if os.IsNotExist(err) {
            s.config = &Config{} // Default config
            return nil
        }
        return err
    }

    return json.Unmarshal(data, &s.config)
}

func (s *ConfigService) ServiceShutdown() error {
    data, err := json.MarshalIndent(s.config, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(s.path, data, 0644)
}
```

### Шаблон 3: фоновый обработчик

```go
type WorkerService struct {
    cancel context.CancelFunc
}

func (s *WorkerService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    workerCtx, cancel := context.WithCancel(ctx)
    s.cancel = cancel

    go s.runWorker(workerCtx)

    return nil
}

func (s *WorkerService) ServiceShutdown() error {
    if s.cancel != nil {
        s.cancel()
    }
    return nil
}

func (s *WorkerService) runWorker(ctx context.Context) {
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            s.doWork()
        case <-ctx.Done():
            return
        }
    }
}
```

## Справочник по жизненному циклу

| Хук/интерфейс | Когда вызывается | Можно отменить? | Назначение |
| --- | --- | --- | --- |
| `ServiceStartup` | Во время `app.Run()`, до запуска цикла событий | Нет (верните ошибку, чтобы прервать запуск) | Инициализация |
| `ServiceShutdown` | Во время завершения работы, после `OnShutdown` | Нет | Освобождение ресурсов |
| `OnShutdown` | После подтверждения выхода | Нет | Освобождение ресурсов приложения |
| `ShouldQuit` | При запросе на выход | Да (верните false) | Подтверждение выхода |
| `RegisterHook(WindowClosing)` | При запросе на закрытие окна | Да (`e.Cancel()`) | Предотвращение закрытия окна |
| `OnWindowEvent` | При возникновении события | Нет | Реагирование на события |
| `OnApplicationEvent` | При возникновении события | Нет | Реагирование на события |

## Различия между платформами

### macOS

- **Меню приложения** остаётся доступным, даже если нет открытых окон
- **Cmd+Q** запускает процесс завершения работы приложения (через `ShouldQuit`)
- **Значок в Dock** остаётся видимым, если его не скрыть
- Используйте `ApplicationShouldTerminateAfterLastWindowClosed` для управления поведением при завершении работы

### Windows

- **Без окна меню приложения отсутствует**
- **Alt+F4** закрывает окно (это можно предотвратить с помощью `RegisterHook`)
- **Область уведомлений** позволяет приложению продолжать работу

### Linux

- **Поведение различается** в зависимости от среды рабочего стола
- **В целом поведение аналогично Windows**

## Отладка проблем жизненного цикла

### Проблема: приложение не завершает работу

**Причины:**

1. `ShouldQuit` возвращает `false`
2. `OnShutdown` выполняется слишком долго
3. Фоновые горутины не останавливаются

**Решение:**

```go
// 1. Check ShouldQuit logic
ShouldQuit: func() bool {
    log.Println("ShouldQuit called")
    return true
}

// 2. Keep OnShutdown fast
OnShutdown: func() {
    log.Println("OnShutdown started")
    // Fast cleanup only
    log.Println("OnShutdown finished")
}

// 3. Use context for background tasks
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    go func() {
        <-ctx.Done()
        log.Println("Context cancelled, stopping background work")
    }()
    return nil
}
```

### Проблема: не удаётся запустить сервис

**Решение:** возвращайте информативные ошибки:

```go
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    if err := s.init(); err != nil {
        return fmt.Errorf("failed to initialise: %w", err)
    }
    return nil
}
```

Ошибка будет записана в журнал, и приложение не запустится.

## Рекомендации

### Что следует делать

- **Используйте сервисы для управления жизненным циклом** — они предоставляют необходимые хуки запуска и завершения работы
- **Быстро завершайте работу** — старайтесь, чтобы вся очистка занимала менее 1 секунды
- **Используйте контекст для отмены** — корректно останавливайте фоновые задачи
- **Обрабатывайте ошибки при запуске** — возвращайте ошибки, чтобы корректно прервать запуск
- **Записывайте события жизненного цикла в журнал** — это помогает при отладке

### Чего не следует делать

- **Не блокируйте выполнение при запуске сервиса** — инициализация должна быть быстрой (менее 2 секунд)
- **Не показывайте диалоговые окна при завершении работы** — приложение закрывается, поэтому интерфейс может не работать
- **Не игнорируйте контекст** — всегда проверяйте `ctx.Done()` в горутинах
- **Не допускайте утечек ресурсов** — всегда реализуйте `ServiceShutdown`

## Дальнейшие действия

**Сервисы** — узнайте больше о системе сервисов [Подробнее →](/features/bindings/services/)

**Система событий** — используйте события для обмена данными [Подробнее →](/features/events/system/)

**Управление окнами** — создавайте окна и управляйте ими [Подробнее →](/features/windows/basics/)

---

**Есть вопросы о жизненном цикле?** Задайте их в [Discord](https://discord.gg/JDdSxwjhGf) или ознакомьтесь с [примерами](https://github.com/wailsapp/wails/tree/master/v3/examples).
