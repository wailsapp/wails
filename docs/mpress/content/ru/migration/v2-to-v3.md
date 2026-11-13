---
title: "Миграция с v2 на v3"
description: "Полное руководство по миграции приложения Wails с v2 на v3"
slug: "migration/v2-to-v3"
sourcePath: "migration/v2-to-v3.md"
---

Wails v3 — это **полностью переписанная версия** со значительными улучшениями архитектуры, производительности и удобства разработки. Это руководство поможет перенести приложение с v2 на v3.

**Основные изменения:**

- Новая структура приложения
- Улучшенная система привязок
- Улучшенное управление окнами
- Улучшенная система событий
- Упрощённая конфигурация

**Время миграции:** 1-4 часов для типичных приложений

## Критические изменения

### Инициализация приложения

В v2 настройка приложения, конфигурация окна и запуск были объединены в один вызов `wails.Run()`. Такой монолитный подход затруднял создание нескольких окон, обработку ошибок на разных этапах и тестирование отдельных компонентов приложения.

В v3 эти задачи разделены на отдельные этапы: создание приложения, создание окна и запуск. Такое разделение даёт явный контроль над каждым этапом жизненного цикла приложения, а код становится более модульным и удобным для тестирования.

**v2:**

```go
err := wails.Run(&options.App{
    Title:  "My App",
    Width:  1024,
    Height: 768,
    Bind: []interface{}{
        &GreetService{},
    },
})
```

**v3:**

```go
app := application.New(application.Options{
    Name: "My App",
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
})

window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "My App",
    Width:  1024,
    Height: 768,
})

app.Run()
```

**Преимущества этого подхода:**

- **Поддержка нескольких окон**: окна можно создавать динамически в любой момент, а не только при запуске
- **Улучшенная обработка ошибок**: каждый этап можно проверять отдельно и корректно обрабатывать возникающие ошибки
- **Более понятный код**: благодаря разделению ясно, что происходит на каждом этапе
- **Упрощённое тестирование**: настройку приложения можно тестировать без запуска цикла событий
- **Больше гибкости**: окна можно создавать, уничтожать и создавать заново на протяжении всего жизненного цикла приложения

### Привязки

В v2 каждой привязанной структуре требовались поле контекста и метод `startup(ctx)` для получения контекста среды выполнения. Это приводило к тесной связи бизнес-логики со средой выполнения Wails и затрудняло тестирование и понимание кода.

В v3 появилась модель сервисов, в которой структуры полностью автономны и не должны хранить контекст среды выполнения. Если сервису нужен доступ к экземпляру приложения, он получает его явно посредством внедрения зависимостей, а не через неявную передачу контекста.

**v2:**

```go
type App struct {
    ctx context.Context
}

func (a *App) startup(ctx context.Context) {
    a.ctx = ctx
}

func (a *App) Greet(name string) string {
    return "Hello " + name
}
```

**v3:**

```go
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello " + name
}

// Register as service
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
})
```

**Преимущества этого подхода:**

- **Нет неявных зависимостей**: сервисы представляют собой обычные структуры Go без скрытых зависимостей от среды выполнения
- **Упрощённое тестирование**: методы сервисов можно тестировать без создания имитации контекста Wails
- **Более понятный код**: зависимости задаются явно — передаются как аргументы конструктора, а не скрываются в поле контекста
- **Улучшенная организация**: сервисы можно группировать по предметным областям, а не размещать все в одной структуре `App`
- **Корректная инициализация**: когда требуется инициализация, используйте метод `ServiceStartup()`, чтобы выполнять её явно

### Среда выполнения

В v2 для всех операций среды выполнения требовалось передавать контекст глобальным функциям из пакета `runtime`. Из-за этого весь код был тесно связан с объектом контекста, а API воспринимался скорее как процедурный, чем как объектно-ориентированный.

В v3 среда выполнения на основе контекста заменена прямыми вызовами методов объектов приложения и окон. Операции вызываются непосредственно у объектов, на которые они воздействуют, благодаря чему код становится понятнее и соответствует объектно-ориентированному подходу.

**v2:**

```go
import "github.com/wailsapp/wails/v2/pkg/runtime"

runtime.WindowSetTitle(a.ctx, "New Title")
runtime.EventsEmit(a.ctx, "event-name", data)
```

**v3:**

```go
// Store app reference
type MyService struct {
    app *application.App
}

func (s *MyService) UpdateTitle() {
    window := s.app.Window.Current()
    window.SetTitle("New Title")
}

func (s *MyService) EmitEvent() {
    s.app.Event.Emit("event-name", data)
}
```

**Преимущества этого подхода:**

- **Объектно-ориентированная архитектура**: методы вызываются у объектов, на которые они воздействуют (окно, приложение, меню и т. д.)
- **Более очевидное назначение**: `window.SetTitle()` понятнее, чем `runtime.WindowSetTitle(ctx, ...)`
- **Улучшенная поддержка в IDE**: автодополнение корректно работает, когда методы принадлежат объектам
- **Понятная работа с несколькими окнами**: при наличии нескольких окон можно явно выбрать окно, над которым выполняется операция
- **Не нужно передавать контекст по цепочке**: контекст не требуется передавать через каждую функцию

### Привязки фронтенда

В v2 привязки группировались по имени пакета Go и структуры, из-за чего обычно получались пути наподобие `wailsjs/go/main/App`. Такая структура не отражала логическую группировку и затрудняла поиск связанной функциональности.

В v3 привязки группируются по имени сервиса и модулю приложения, образуя более понятную логическую структуру. Они создаются в каталоге `bindings` и организуются по имени приложения и именам сервисов, поэтому доступную функциональность проще понять.

**v2:**

```javascript
import { Greet } from '../wailsjs/go/main/App'

const result = await Greet("World")
```

**v3:**

```javascript
import { Greet } from './bindings/changeme/greetservice'

const result = await Greet("World")
```

**Преимущества этого подхода:**

- **Логичная организация**: привязки группируются по имени сервиса, а не по структуре пакетов Go
- **Более понятные пути импорта**: путь отражает логику предметной области (greetservice), а не структуру файлов (main/App)
- **Упрощённый поиск**: по привязкам можно перемещаться по функциональности, а не по технической структуре
- **Единообразное именование**: организация по сервисам соответствует архитектуре бэкенда
- **Более простые пути**: больше не нужен префикс `../wailsjs/go` — только `./bindings`

### События

В v2 события использовали вариативные параметры `interface{}`, а каждой функции событий требовалось передавать контекст. Обработчики событий получали нетипизированные данные, для которых приходилось вручную выполнять приведение типов. Из-за этого система событий была подвержена ошибкам, а её отладка была затруднена.

В v3 появились типизированные объекты событий, а требование передавать контекст было устранено. Обработчики получают полноценный объект события с типизированными данными, благодаря чему система событий становится надёжнее и проще в использовании.

**v2:**

```go
runtime.EventsOn(ctx, "event-name", func(data ...interface{}) {
    // Handle event
})

runtime.EventsEmit(ctx, "event-name", data)
```

**v3:**

```go
app.Event.On("event-name", func(e *application.CustomEvent) {
    data := e.Data
    // Handle event
})

app.Event.Emit("event-name", data)
```

**Преимущества этого подхода:**

- **Безопасность типов**: вместо `...interface{}` события используют полноценные объекты событий
- **Удобнее отлаживать**: объекты событий содержат метаданные, например имя события, что упрощает отладку
- **Более понятный API**: `app.Event.On()` и `app.Event.Emit()` интуитивно понятнее функций среды выполнения
- **Контекст не требуется**: события работают непосредственно с объектом приложения без передачи контекста между потоками
- **Более простые обработчики**: у обработчиков событий чёткая сигнатура без вариативных параметров

### Окна

v2 поддерживал только одно окно на приложение. Окно создавалось при запуске, а все операции с ним выполнялись через функции среды выполнения, неявно обращавшиеся к этому единственному окну.

v3 добавляет встроенную поддержку нескольких окон в качестве основной возможности. Каждое окно представляет собой полноценный объект с собственными методами и жизненным циклом. На протяжении жизненного цикла приложения можно динамически создавать несколько окон, управлять ими и уничтожать их.

**v2:**

```go
// Single window only
runtime.WindowSetSize(ctx, 800, 600)
```

**v3:**

```go
// Multiple windows supported
window1 := app.Window.New()
window1.SetSize(800, 600)

window2 := app.Window.New()
window2.SetSize(1024, 768)
```

**Почему это лучше:**

- **Многооконные приложения**: создавайте приложения с несколькими независимыми окнами (панелями мониторинга, настройками, инструментами и т. д.)
- **Явные ссылки на окна**: каждое окно — это объект, который можно хранить и изменять напрямую
- **Динамическое создание окон**: создавайте и уничтожайте окна в любой момент во время работы приложения
- **Независимое состояние окон**: у каждого окна есть собственные события, свойства и жизненный цикл
- **Улучшенная архитектура**: управление окнами основано на объектах, а не на контексте

## Этапы миграции

### Шаг 1: обновите зависимости

**go.mod:**

```go
module myapp

go 1.25.0

require (
    github.com/wailsapp/wails/v3 v3.0.0
)
```

**Обновите:**

```bash
go get github.com/wailsapp/wails/v3@v3.0.0
go mod tidy
```

### Шаг 2: обновите main.go

**v2:**

```go
package main

import (
    "embed"
    "github.com/wailsapp/wails/v2/pkg/options"
    "github.com/wailsapp/wails/v2/pkg/options/assetserver"
    "github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
    app := NewApp()

    err := wails.Run(&options.App{
        Title:  "My App",
        Width:  1024,
        Height: 768,
        AssetServer: &assetserver.Options{
            Assets: assets,
        },
        Bind: []interface{}{
            app,
        },
        Windows: &windows.Options{
            WebviewIsTransparent: false,
        },
    })

    if err != nil {
        println("Error:", err.Error())
    }
}
```

**v3:**

```go
package main

import (
    "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name: "My App",
        Services: []application.Service{
            application.NewService(&MyService{}),
        },
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:  "My App",
        Width:  1024,
        Height: 768,
    })

    err := app.Run()
    if err != nil {
        panic(err)
    }
}
```

### Шаг 3: преобразуйте структуру App в сервис

**v2:**

```go
type App struct {
    ctx context.Context
}

func NewApp() *App {
    return &App{}
}

func (a *App) startup(ctx context.Context) {
    a.ctx = ctx
    // Initialisation
}

func (a *App) Greet(name string) string {
    return "Hello " + name
}
```

**v3:**

```go
type MyService struct {
    app *application.App
}

func NewMyService(app *application.App) *MyService {
    return &MyService{app: app}
}

func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    // Initialisation
    return nil
}

func (s *MyService) Greet(name string) string {
    return "Hello " + name
}

// Register after app creation
app := application.New(application.Options{})
app.RegisterService(application.NewService(NewMyService(app)))
```

### Шаг 4: обновите вызовы среды выполнения

**v2:**

```go
func (a *App) DoSomething() {
    runtime.WindowSetTitle(a.ctx, "New Title")
    runtime.EventsEmit(a.ctx, "update", data)
    runtime.LogInfo(a.ctx, "Message")
}
```

**v3:**

```go
func (s *MyService) DoSomething() {
    window := s.app.Window.Current()
    window.SetTitle("New Title")
    
    s.app.Event.Emit("update", data)
    
    s.app.Logger.Info("Message")
}
```

### Шаг 5: обновите фронтенд

**Сгенерируйте новые привязки:**

```bash
wails3 generate bindings
```

**Обновите импорты:**

```javascript
// v2
import { Greet } from '../wailsjs/go/main/App'

// v3
import { Greet } from './bindings/changeme/myservice'
```

**Обновите обработку событий:**

```javascript
// v2
import { EventsOn, EventsEmit } from '../wailsjs/runtime/runtime'

EventsOn("update", (data) => {
    console.log(data)
})

EventsEmit("action", data)

// v3
import { Events } from '@wailsio/runtime'

Events.On("update", (data) => {
    console.log(data)
})

Events.Emit("action", data)
```

### Шаг 6: обновите конфигурацию

**v2 (wails.json):**

```json
{
  "name": "myapp",
  "outputfilename": "myapp",
  "frontend:install": "npm install",
  "frontend:build": "npm run build",
  "frontend:dev:watcher": "npm run dev",
  "frontend:dev:serverUrl": "auto"
}
```

**v3 (wails.json):**

```json
{
  "name": "myapp",
  "frontend": {
    "dir": "./frontend",
    "install": "npm install",
    "build": "npm run build",
    "dev": "npm run dev",
    "devServerUrl": "http://localhost:5173"
  }
}
```

## Сопоставление возможностей

### Диалоговые окна

**v2:**

```go
selection, err := runtime.OpenFileDialog(ctx, runtime.OpenDialogOptions{
    Title: "Select File",
})
```

**v3:**

```go
selection, err := app.Dialog.OpenFileWithOptions(&application.OpenFileDialogOptions{
    Title: "Select File",
}).PromptForSingleSelection()
```

### Меню

**v2:**

```go
menu := menu.NewMenu()
menu.Append(menu.Text("File", nil, []*menu.MenuItem{
    menu.Text("Quit", nil, func(_ *menu.CallbackData) {
        runtime.Quit(ctx)
    }),
}))
```

**v3:**

```go
menu := app.NewMenu()
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})
```

### Область уведомлений

**v2:**

```go
// Not available in v2
```

**v3:**

```go
systray := app.SystemTray.New()
systray.SetIcon(iconBytes)
systray.SetLabel("My App")

menu := app.NewMenu()
menu.Add("Show").OnClick(showWindow)
menu.Add("Quit").OnClick(app.Quit)
systray.SetMenu(menu)
```

## Распространённые проблемы

### Проблема: привязки не найдены

**Проблема:** ошибки импорта после миграции

**Решение:**

```bash
# Regenerate bindings
wails3 generate bindings

# Check output directory
ls frontend/bindings
```

### Проблема: ошибки контекста

**Проблема:** `ctx` недоступен

**Решение:**

Вместо этого сохраните ссылку на приложение:

```go
type MyService struct {
    app *application.App
}

func NewMyService(app *application.App) *MyService {
    return &MyService{app: app}
}
```

### Проблема: методы окна не работают

**Проблема:** `runtime.WindowSetTitle()` не существует

**Решение:**

Вызывайте методы окна напрямую:

```go
window := s.app.Window.Current()
window.SetTitle("New Title")
```

### Проблема: события не срабатывают

**Проблема:** события зарегистрированы, но не поступают

**Решение:**

Убедитесь, что имена событий совпадают в точности:

```go
// Go
app.Event.Emit("my-event", data)

// JavaScript
OnEvent("my-event", handler)  // Must match exactly
```

## Тестирование миграции

### Контрольный список

- [ ] Приложение запускается без ошибок
- [ ] Все привязки работают
- [ ] События отправляются и принимаются
- [ ] Окна корректно открываются и закрываются
- [ ] Меню работают (если применимо)
- [ ] Диалоговые окна работают (если применимо)
- [ ] Область уведомлений работает (если применимо)
- [ ] Процесс сборки работает
- [ ] Рабочая сборка успешно запускается

### Команды тестирования

```bash
# Development
wails3 dev

# Build
wails3 build

# Generate bindings
wails3 generate bindings
```

## Преимущества v3

### Производительность

- **Более быстрый запуск** — оптимизированная инициализация
- **Меньше памяти** — эффективное использование ресурсов
- **Улучшенный мост** — накладные расходы на вызов составляют &lt;1мс

### Возможности

- **Несколько окон** — встроенная поддержка
- **Область уведомлений** — встроенная поддержка
- **Улучшенные события** — типизированный и более простой API
- **Сервисы** — улучшенная организация кода

### Удобство разработки

- **Безопасность типов** — полная поддержка TypeScript
- **Улучшенные ошибки** — понятные сообщения об ошибках
- **Горячая перезагрузка** — более быстрая разработка
- **Улучшенная документация** — исчерпывающие руководства

## Получение помощи

### Ресурсы

- [Документация](/quick-start/why-wails/)
- [Сообщество в Discord](https://discord.gg/JDdSxwjhGf)
- [Задачи на GitHub](https://github.com/wailsapp/wails/issues)
- [Примеры](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples)

### Частые вопросы

**В: Можно ли одновременно использовать v2 и v3?** О: Да, у них разные пути импорта.

**В: Готов ли v3 к использованию в рабочей среде?** Wails v3.0.0 — стабильный выпуск для настольных приложений на Windows, macOS и Linux. Поддержка Android и iOS остаётся экспериментальной и не входит в гарантии стабильной совместимости для настольных платформ.

**В: Будет ли поддерживаться v2?** О: Да, v2 будет получать критически важные обновления.

**В: Сколько времени занимает миграция?** О: Для типичных приложений — 1-4 ч.

## Дальнейшие шаги

@cards{cols="2"}
🚀 Быстрый старт
Начните работу с Wails v3.

[Подробнее →](/quick-start/installation/)

---
★ Основные понятия
Разберитесь в архитектуре v3.

[Подробнее →](/concepts/architecture/)

---
◆ Привязки
Изучите новую систему привязок.

[Подробнее →](/features/bindings/methods/)

---
📖 Примеры
Ознакомьтесь с полными примерами v3.

[Посмотреть примеры →](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples)

@end

---

**Остались вопросы?** Задайте их в [Discord](https://discord.gg/JDdSxwjhGf) или [создайте задачу](https://github.com/wailsapp/wails/issues).
