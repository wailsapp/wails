---
title: "Novidades do Wails v3"
description: "Conheça as principais melhorias e os novos recursos do Wails v3"
slug: "whats-new"
sourcePath: "whats-new.md"
---

O Wails v3 introduz mudanças significativas em relação ao v2. Ele substitui a API declarativa de janela única por uma abordagem procedural mais flexível. Esse novo design de API melhora a legibilidade do código e simplifica o desenvolvimento, especialmente para aplicações complexas com várias janelas.

O Wails v3 representa uma evolução substancial na maneira como aplicações para desktop podem ser criadas usando Go e tecnologias web.

## Várias janelas

O Wails v3 introduz a capacidade de criar e gerenciar várias janelas em uma única aplicação. Esse recurso permite que os desenvolvedores criem interfaces de usuário mais complexas e versáteis, superando as limitações das aplicações de janela única.

Cada janela pode ser configurada de forma independente, oferecendo flexibilidade de tamanho, posição, conteúdo e comportamento. Isso permite criar aplicações com janelas separadas para diferentes funcionalidades, como interfaces principais, painéis de configurações ou visualizações auxiliares.

Os desenvolvedores podem criar, manipular e gerenciar essas janelas por meio de código, possibilitando interfaces de usuário dinâmicas que se adaptam às necessidades do usuário e aos estados da aplicação.

@note{type="tip" title="Várias janelas"}
@details{title="Exemplo"}
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

## Integração com a bandeja do sistema

O Wails v3 introduz suporte robusto aos recursos da bandeja do sistema, permitindo que sua aplicação mantenha uma presença constante na área de trabalho do usuário. Esse recurso é particularmente útil para aplicações que precisam ser executadas em segundo plano ou oferecer acesso rápido a funções essenciais.

Os principais recursos da integração do Wails v3 com a bandeja do sistema incluem:

1. Vinculação de janela: você pode associar uma janela ao ícone da bandeja do sistema. Quando ativada, essa janela será centralizada em relação à posição do ícone, oferecendo uma ótima maneira de acessar rapidamente sua aplicação.

2. Suporte completo a menus: crie menus avançados e interativos que os usuários possam acessar diretamente pelo ícone da bandeja do sistema. Isso permite executar ações rápidas sem precisar abrir a janela completa da aplicação.

3. Exibição adaptável de ícones: o suporte a ícones para os modos claro e escuro garante que o ícone da sua aplicação na bandeja do sistema permaneça visível e esteticamente agradável em diferentes temas do sistema. A integração também oferece suporte a ícones de modelo no macOS.

@note{type="tip" title="Bandeja do sistema"}
@details{title="Exemplo"}
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

## Geração de bindings aprimorada

O Wails v3 introduz uma melhoria significativa na forma como os bindings são gerados para o seu projeto. Os bindings são a camada que conecta o backend em Go ao frontend, permitindo uma comunicação transparente entre ambos.

A geração de bindings agora é feita por um analisador estático sofisticado, que melhora radicalmente o processo. Ele oferece mais velocidade e preserva a qualidade do código ao manter comentários e nomes de parâmetros.

O processo de geração de bindings foi simplificado e agora requer apenas um único comando: `wails3 generate bindings`.

@note{type="tip" title="Bindings"}
@details{title="Exemplo"}
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

## Sistema de build aprimorado

O Wails v3 introduz um sistema de build mais flexível e transparente, solucionando as limitações de seu antecessor. No v2, o processo de build era em grande parte opaco e difícil de personalizar, o que podia ser frustrante para desenvolvedores que buscavam mais controle sobre o processo de build de seus projetos.

Todo o trabalho pesado realizado pelo sistema de build do v2, como a geração de ícones e a criação de manifestos, foi adicionado à CLI na forma de comandos de ferramentas. Incorporamos o [Taskfile](https://taskfile.dev) à CLI para orquestrar essas chamadas e oferecer a mesma experiência de desenvolvimento do v2. No entanto, essa abordagem proporciona o equilíbrio ideal entre flexibilidade e facilidade de uso, pois agora você pode personalizar o processo de build de acordo com suas necessidades.

Você pode até usar o make, se preferir!

@note{type="tip" title="Taskfile.yml"}
@details{title="Exemplo"}
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

## Eventos aprimorados

O Wails agora emite eventos para diversas operações de tempo de execução e atividades do sistema. Isso permite que seu aplicativo responda a esses eventos em tempo real. Além disso, há eventos multiplataforma (comuns), permitindo que você escreva métodos consistentes de tratamento de eventos que funcionem em diferentes sistemas operacionais.

É possível registrar hooks de eventos para tratar eventos específicos de forma síncrona. Diferentemente do método `On`, esses hooks permitem cancelar o evento, se necessário. Um caso de uso comum é exibir uma caixa de diálogo de confirmação antes de fechar uma janela. Isso proporciona maior controle sobre o fluxo de eventos e a experiência do usuário.

@note{type="tip" title="Exemplo de tratamento de eventos"}
@details{title="Exemplo"}
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

## Linguagem de Marcação do Wails (wml)

Um recurso experimental para chamar métodos de tempo de execução usando HTML simples, semelhante ao  [htmx](https://htmx.org).

@note{type="tip" title="Exemplo de wml"}
@details{title="Exemplo"}
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

## Exemplos

Há mais exemplos disponíveis no diretório  [examples](https://github.com/wailsapp/wails/tree/master/v3/examples). Confira!
