---
title: "Что нового в Wails v3"
description: "Узнайте об основных улучшениях и новых возможностях Wails v3"
slug: "whats-new"
sourcePath: "whats-new.md"
---

В Wails v3 внесены значительные изменения по сравнению с v2. Вместо декларативного API для работы с одним окном теперь используется более гибкий процедурный подход. Новый дизайн API повышает читаемость кода и упрощает разработку, особенно сложных многооконных приложений.

Wails v3 представляет собой существенный шаг вперёд в разработке настольных приложений с использованием Go и веб-технологий.

## Несколько окон

Wails v3 позволяет создавать несколько окон в одном приложении и управлять ими. Благодаря этой возможности разработчики могут создавать более сложные и универсальные пользовательские интерфейсы, не ограничиваясь однооконными приложениями.

Каждое окно можно настраивать независимо, включая его размер, положение, содержимое и поведение. Это позволяет создавать приложения с отдельными окнами для разных функций, например для основного интерфейса, панели настроек или вспомогательных представлений.

Разработчики могут программно создавать, изменять и контролировать эти окна, формируя динамические пользовательские интерфейсы, которые адаптируются к потребностям пользователя и состоянию приложения.

@note{type="tip" title="Несколько окон"}
@details{title="Пример"}
```go
package main

import (
   "embed"
   "log"
   
   "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/*
var assets embed.FS

func main() {

   app := application.New(application.Options{
        Name:   "Multi Window Demo",
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
   })
   
   window1 := app.Window.NewWithOptions(application.WebviewWindowOptions{
       Title:  "Window 1",
   })
   
   window2 := app.Window.NewWithOptions(application.WebviewWindowOptions{
       Title:  "Window 2",
   })
   
   // load the embedded html from the embed.FS
   window1.SetURL("/")
   window1.Center()
   
   // Load an external URL
   window2.SetURL("https://wails.io")
   
   err := app.Run()

   if err != nil {
	   log.Fatal(err.Error())
   }
}
```

@end

@end

## Интеграция с областью уведомлений

Wails v3 предоставляет развитую поддержку области уведомлений, позволяя приложению постоянно оставаться доступным на рабочем столе пользователя. Эта возможность особенно полезна для приложений, которым необходимо работать в фоновом режиме или предоставлять быстрый доступ к ключевым функциям.

Основные возможности интеграции Wails v3 с областью уведомлений:

1. Привязка окна: к значку в области уведомлений можно привязать окно. При активации оно будет располагаться по центру относительно значка, обеспечивая быстрый доступ к приложению.

2. Полноценная поддержка меню: создавайте функциональные интерактивные меню, доступные непосредственно через значок в области уведомлений. Пользователи смогут быстро выполнять действия, не открывая полное окно приложения.

3. Адаптивное отображение значка: поддержка значков для светлого и тёмного режимов обеспечивает хорошую видимость и привлекательный вид значка приложения в области уведомлений при разных системных темах. В macOS также поддерживаются шаблонные значки.

@note{type="tip" title="Область уведомлений"}
@details{title="Пример"}
```go
package main

import (
    "log"
    "runtime"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/icons"
)

func main() {
    app := application.New(application.Options{
        Name:        "Systray Demo",
        Mac: application.MacOptions{
            ActivationPolicy: application.ActivationPolicyAccessory,
        },
    })

    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Width:       500,
        Height:      800,
        Frameless:   true,
        AlwaysOnTop: true,
        Hidden:      true,
        Windows: application.WindowsWindow{
            HiddenOnTaskbar: true,
        },
    })

    systemTray := app.SystemTray.New()

    // Support for template icons on macOS
    if runtime.GOOS == "darwin" {
        systemTray.SetTemplateIcon(icons.SystrayMacTemplate)
    } else {
        // Support for light/dark mode icons
        systemTray.SetDarkModeIcon(icons.SystrayDark)
        systemTray.SetIcon(icons.SystrayLight)
    }

    // Support for menu
    myMenu := app.Menu.New()
    myMenu.Add("Hello World!").OnClick(func(_ *application.Context) {
        println("Hello World!")
    })
    systemTray.SetMenu(myMenu)

    // This will center the window to the systray icon with a 5px offset
    // It will automatically be shown when the systray icon is clicked
    // and hidden when the window loses focus
    systemTray.AttachWindow(window).WindowOffset(5)

    err := app.Run()
    if err != nil {
        log.Fatal(err)
    }
}
```

@end

@end

## Улучшенная генерация привязок

В Wails v3 существенно улучшена генерация привязок для проекта. Привязки соединяют серверную часть на Go с клиентской, обеспечивая беспрепятственный обмен данными между ними.

Теперь привязки генерируются с помощью развитого статического анализатора, который радикально улучшает этот процесс. Он повышает скорость и помогает сохранять качество кода, оставляя без изменений комментарии и имена параметров.

Процесс генерации привязок упрощён: теперь для него требуется всего одна команда — `wails3 generate bindings`.

@note{type="tip" title="Привязки"}
@details{title="Пример"}
```js
// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

// Generated layout (excerpt): frontend/bindings/<full-go-import-path>/greetservice.js
import { Call as $Call, Create as $Create } from "/wails/runtime.js";

/**
 * Greet greets a person
 * @param {string} $0
 * @returns {Promise<string>}
 */
export function Greet($0) {
    return $Call.ByID(1411160069, $0);
}

/**
 * GreetPerson greets a person
 * @param {main.Person} $0
 * @returns {Promise<string>}
 */
export function GreetPerson($0) {
    return $Call.ByID(4021313248, $0);
}
```

@end

@end

## Улучшенная система сборки

В Wails v3 представлена более гибкая и прозрачная система сборки, устраняющая ограничения предыдущей версии. В v2 процесс сборки был преимущественно непрозрачным и с трудом поддавался настройке, что могло создавать неудобства разработчикам, которым требовалось больше контроля над сборкой проекта.

Все сложные операции, которые выполняла система сборки v2, например генерация значков и создание манифеста, добавлены в CLI в виде команд инструментов. Для управления их вызовами мы встроили [Taskfile](https://taskfile.dev) в CLI, сохранив тот же уровень удобства для разработчиков, что и в v2. При этом такой подход обеспечивает оптимальный баланс гибкости и простоты использования, поскольку теперь процесс сборки можно адаптировать к своим потребностям.

Если вам нравится make, можно использовать даже его!

@note{type="tip" title="Taskfile.yml"}
@details{title="Пример"}
```yaml {title="build/Taskfile.darwin.yml"}
darwin:build:
  summary: Builds the application for macOS
  platforms:
    - darwin
  cmds:
    - task: common:go:mod:tidy
    - task: common:build:frontend
    - task: common:generate:icons
    - task: darwin:build:app
  env:
    CGO_CFLAGS: "-mmacosx-version-min=10.15"
    CGO_LDFLAGS: "-mmacosx-version-min=10.15"
    MACOSX_DEPLOYMENT_TARGET: "10.15"
```

@end

@end

## Улучшенная система событий

Теперь Wails генерирует события для различных операций среды выполнения и системных действий. Это позволяет приложению реагировать на них в реальном времени. Кроме того, доступны кроссплатформенные (общие) события, благодаря которым можно создавать единообразные методы обработки событий, работающие в разных операционных системах.

Для синхронной обработки определённых событий можно регистрировать перехватчики. В отличие от метода `On` эти перехватчики позволяют при необходимости отменить событие. Типичный сценарий использования — отображение диалогового окна подтверждения перед закрытием окна. Это даёт больше контроля над ходом обработки событий и взаимодействием с пользователем.

@note{type="tip" title="Пример обработки событий"}
@details{title="Пример"}
```go
package main

import (
    "embed"
    "log"
    "time"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed assets
var assets embed.FS

func main() {

    app := application.New(application.Options{
        Name:        "Events Demo",
        Description: "A demo of the Events API",
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
        Mac: application.MacOptions{
            ApplicationShouldTerminateAfterLastWindowClosed: true,
        },
    })

    // Custom event handling — App.Event.On(name, func(e *CustomEvent))
    app.Event.On("myevent", func(e *application.CustomEvent) {
        log.Printf("[Go] CustomEvent received: %+v\n", e)
    })

    // OS-specific application events — App.Event.OnApplicationEvent(eventType, func(e *ApplicationEvent))
    app.Event.OnApplicationEvent(events.Mac.ApplicationDidFinishLaunching, func(event *application.ApplicationEvent) {
        println("events.Mac.ApplicationDidFinishLaunching fired!")
    })

    // Platform-agnostic events
    app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(event *application.ApplicationEvent) {
        println("events.Common.ApplicationStarted fired!")
    })

    win1 := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: "Takes 3 attempts to close me!",
    })

    var countdown = 3

    // Register a hook to cancel the window closing
    win1.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
        countdown--
        if countdown == 0 {
            println("Closing!")
            return
        }
        println("Nope! Not closing!")
        e.Cancel()
    })

    win1.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
        println("[Event] Window focus!")
    })

    err := app.Run()

    if err != nil {
        log.Fatal(err.Error())
    }
}
```

@end

@end

## Язык разметки Wails (wml)

Экспериментальная функция для вызова методов среды выполнения с помощью обычного HTML, аналогично [htmx](https://htmx.org).

@note{type="tip" title="Пример использования wml"}
@details{title="Пример"}
```html
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <title>Wails ML Demo</title>
  </head>
  <body style="margin-top:50px; color: white; background-color: #191919">
    <h2>Wails ML Demo</h2>
    <p>This application contains no Javascript!</p>
    <button wml-event="button-pressed">Press me!</button>
    <button wml-event="delete-things" wml-confirm="Are you sure?">
      Delete all the things!
    </button>
    <button wml-window="Close" wml-confirm="Are you sure?">
      Close the Window?
    </button>
    <button wml-window="Center">Center</button>
    <button wml-window="Minimise">Minimise</button>
    <button wml-window="Maximise">Maximise</button>
    <button wml-window="UnMaximise">UnMaximise</button>
    <button wml-window="Fullscreen">Fullscreen</button>
    <button wml-window="UnFullscreen">UnFullscreen</button>
    <button wml-window="Restore">Restore</button>
    <div
      style="width: 200px; height: 200px; border: 2px solid white;"
      wml-event="hover"
      wml-trigger="mouseover"
    >
      Hover over me
    </div>
  </body>
</html>
```

@end

@end

## Примеры

Дополнительные примеры находятся в каталоге [examples](https://github.com/wailsapp/wails/tree/master/v3/examples). Ознакомьтесь с ними!
