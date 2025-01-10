---
title: What's New in Wails v3
---

Wails v3 introduces significant changes from v2. It replaces the
single-window, declarative API with a more flexible procedural approach. This
new API design improves code readability and simplifies development, especially
for complex multi-window applications.

Wails v3 represents a substantial evolution in how desktop applications
can be built using Go and web technologies.

## Multiple Windows

Wails v3 introduces the ability to create and manage multiple windows within a
single application. This feature allows developers to design more complex and
versatile user interfaces, moving beyond the limitations of single-window
applications.

Each window can be independently configured, providing flexibility in terms of
size, position, content, and behavior. This enables the creation of applications
with separate windows for different functionalities, such as main interfaces,
settings panels, or auxiliary views.

Developers can create, manipulate, and manage these windows programmatically,
allowing for dynamic user interfaces that adapt to user needs and application
states.

:::tip[Multiple Windows]
<details><summary>Example</summary>

```go
package main

import (
   _ "embed"
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
   
   window1 := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
       Title:  "Window 1",
   })
   
   window2 := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
       Title:  "Window 2",
   })
   
   // load the embedded html from the embed.FS
   window1.SetURL("/")
   window1.Center()
   
   // Load an external URL
   window2.SetURL("https://wails.app")
   
   err := app.Run()

   if err != nil {
	   log.Fatal(err.Error())
   }
}
```
</details>

:::

## System Tray Integration

Wails v3 introduces robust support for system tray functionality, allowing your
application to maintain a persistent presence on the user's desktop. This
feature is particularly useful for applications that need to run in the
background or provide quick access to key functions.

Key features of the Wails v3 system tray integration include:

1. Window Attachment: You can associate a window with the system tray icon. When
   activated, this window will be centered relative to the icon's position,
   providing a great way to quickly access your application.

2. Comprehensive Menu Support: Create rich, interactive menus that users can
   access directly from the system tray icon. This allows for quick actions
   without needing to open the full application window.

3. Adaptive Icon Display: Support for both light and dark mode icons ensures
   your application's system tray icon remains visible and aesthetically
   pleasing across different system themes. Template icons are also supported on macOS.

:::tip[Systray]

<details><summary>Example</summary>


```go
package main

import (
    _ "embed"
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

    window := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
        Width:       500,
        Height:      800,
        Frameless:   true,
        AlwaysOnTop: true,
        Hidden:      true,
        Windows: application.WindowsWindow{
            HiddenOnTaskbar: true,
        },
    })

    systemTray := app.NewSystemTray()

    // Support for template icons on macOS
    if runtime.GOOS == "darwin" {
        systemTray.SetTemplateIcon(icons.SystrayMacTemplate)
    } else {
        // Support for light/dark mode icons
        systemTray.SetDarkModeIcon(icons.SystrayDark)
        systemTray.SetIcon(icons.SystrayLight)
    }

    // Support for menu
    myMenu := app.NewMenu()
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
</details>
:::

## Improved bindings generation

Wails v3 introduces a significant improvement in how bindings are generated for
your project. Bindings are the glue that connects your Go backend to your
frontend, allowing seamless communication between the two.

Binding generation is now done using a sophisticated static analyzer that
radically improves the binding generation process. It offers enhanced speed and
preserves code quality by maintaining comments and parameter names.

The binding generation process has been simplified, requiring only a single
command: `wails3 generate bindings`.

:::tip[Bindings]

<details><summary>Example</summary>

```js
// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

import { main } from "./models";

window.go = window.go || {};
window.go.main = {
  GreetService: {
    /**
     * GreetService.Greet
     * Greet greets a person
     * @param name {string}
     * @returns {Promise<string>}
     **/
    Greet: function (name) {
      wails.CallByID(1411160069, ...Array.prototype.slice.call(arguments, 0));
    },

    /**
     * GreetService.GreetPerson
     * GreetPerson greets a person
     * @param person {main.Person}
     * @returns {Promise<string>}
     **/
    GreetPerson: function (person) {
      wails.CallByID(4021313248, ...Array.prototype.slice.call(arguments, 0));
    },
  },
};
```

</details>

:::

## Improved build system

Wails v3 introduces a more flexible and transparent build system, addressing the
limitations of its predecessor. In v2, the build process was largely opaque and
difficult to customise, which could be frustrating for developers seeking more
control over their project's build process.

All the heavy lifting that the v2 build system did, such as icon generation and
manifest creation, have been added as tool commands in the CLI. We have
incorporated [Taskfile](https://taskfile.dev) into the CLI to orchestrate these
calls to bring the same developer experience as v2. However, this approach
brings the ultimate balance of flexibility and ease of use as you can now
customise the build process to your needs.

You can even use make if that's your thing!

:::tip[Taskfile.yml]

<details><summary>Example</summary>

```yaml "Snippet from Taskfile.yml"
build:darwin:
  summary: Builds the application
  platforms:
    - darwin
  cmds:
    - task: pre-build
    - task: build-frontend
    - go build -gcflags=all="-N -l" -o bin/{{.APP_NAME}}
    - task: post-build
  env:
  CGO_CFLAGS: "-mmacosx-version-min=10.13"
  CGO_LDFLAGS: "-mmacosx-version-min=10.13"
  MACOSX_DEPLOYMENT_TARGET: "10.13"
```
</details>

:::

## Improved events

Wails now emits events for various runtime operations and system activities.
This allows your application to respond to these events in real-time.
Additionally, cross-platform (common) events are available, enabling you to
write consistent event handling methods that work across different operating
systems.

Event hooks can be registered to handle specific events synchronously. Unlike
the `On` method, these hooks allow you to cancel the event if needed. A common
use case is displaying a confirmation dialog before closing a window. This gives
you more control over the event flow and user experience.

:::tip[Example of event handling]

<details><summary>Example</summary>

```go
package main

import (
    _ "embed"
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

    // Custom event handling
    app.Events.On("myevent", func(e *application.WailsEvent) {
        log.Printf("[Go] WailsEvent received: %+v\n", e)
    })

    // OS specific application events
    app.On(events.Mac.ApplicationDidFinishLaunching, func(event *application.Event) {
        println("events.Mac.ApplicationDidFinishLaunching fired!")
    })

    // Platform agnostic events
    app.On(events.Common.ApplicationStarted, func(event *application.Event) {
        println("events.Common.ApplicationStarted fired!")
    })

    win1 := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
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

    win1.On(events.Common.WindowFocus, func(e *application.WindowEvent) {
        println("[Event] Window focus!")
    })

    err := app.Run()

    if err != nil {
        log.Fatal(err.Error())
    }
}
```

</details>

:::

## Wails Markup Language (wml)

An experimental feature to call runtime methods using plain html, similar to
[htmx](https://htmx.org).

:::tip[Example of wml]

<details><summary>Example</summary>

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

</details>

:::

## Examples

There are more examples available in the
[examples](https://github.com/wailsapp/wails/tree/v3-alpha/v3/examples)
directory. Check them out!
