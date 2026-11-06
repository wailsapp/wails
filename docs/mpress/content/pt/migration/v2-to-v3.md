---
title: "Migração do v2 para o v3"
description: "Guia completo para migrar seu aplicativo Wails v2 para o v3"
slug: "migration/v2-to-v3"
sourcePath: "migration/v2-to-v3.md"
---

O Wails v3 é uma **reescrita completa**, com melhorias significativas na arquitetura, no desempenho e na experiência de desenvolvimento. Este guia ajuda você a migrar seu aplicativo v2 para o v3.

**Principais mudanças:**

- Nova estrutura do aplicativo
- Sistema de bindings aprimorado
- Gerenciamento de janelas aprimorado
- Sistema de eventos aprimorado
- Configuração simplificada

**Tempo de migração:** 1-4 horas para aplicativos típicos

## Alterações incompatíveis

### Inicialização do aplicativo

No v2, a configuração do aplicativo, a configuração da janela e a execução eram combinadas em uma única chamada a `wails.Run()`. Essa abordagem monolítica dificultava criar várias janelas, tratar erros em diferentes etapas ou testar componentes individuais do aplicativo.

O v3 separa essas responsabilidades em fases distintas: criação do aplicativo, criação da janela e execução. Essa separação oferece controle explícito sobre cada etapa do ciclo de vida do aplicativo e torna o código mais modular e testável.

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

**Por que isso é melhor:**

- **Suporte a várias janelas**: você pode criar janelas dinamicamente a qualquer momento, não apenas durante a inicialização
- **Melhor tratamento de erros**: cada fase pode ser validada separadamente, com o tratamento adequado de erros
- **Código mais claro**: a separação deixa evidente o que acontece em cada etapa
- **Mais testável**: você pode testar a configuração do aplicativo sem executar o loop de eventos
- **Mais flexível**: as janelas podem ser criadas, destruídas e recriadas ao longo do ciclo de vida do aplicativo

### Bindings

No v2, cada struct vinculada precisava de um campo de contexto e de um método `startup(ctx)` para receber o contexto de runtime. Isso criava um forte acoplamento entre a lógica de negócios e o runtime do Wails, dificultando os testes e a compreensão do código.

O v3 introduz o padrão de serviços, no qual suas structs são totalmente independentes e não precisam armazenar o contexto de runtime. Se um serviço precisar acessar a instância do aplicativo, ele a receberá explicitamente por injeção de dependência, em vez de depender da propagação implícita do contexto.

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

**Por que isso é melhor:**

- **Sem dependências implícitas**: os serviços são structs Go comuns, sem dependências ocultas do runtime
- **Testes mais fáceis**: você pode testar os métodos dos serviços sem simular um contexto do Wails
- **Código mais claro**: as dependências são explícitas — passadas como argumentos do construtor — em vez de ficarem ocultas em um campo de contexto
- **Melhor organização**: os serviços podem ser agrupados por domínio, em vez de todos residirem em uma única struct `App`
- **Inicialização adequada**: quando precisar de inicialização, use o método `ServiceStartup()` para torná-la explícita

### Runtime

No v2, todas as operações de runtime exigiam passar um contexto para funções globais do pacote `runtime`. Isso criava um forte acoplamento ao objeto de contexto em toda a base de código e fazia a API parecer procedural, em vez de orientada a objetos.

O v3 substitui o runtime baseado em contexto por chamadas diretas a métodos dos objetos de aplicativo e janela. As operações são chamadas diretamente nos objetos que elas afetam, tornando o código mais intuitivo e orientado a objetos.

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

**Por que isso é melhor:**

- **Design orientado a objetos**: os métodos são chamados nos objetos que eles afetam (janela, aplicativo, menu etc.)
- **Intenção mais clara**: `window.SetTitle()` é mais explícito que `runtime.WindowSetTitle(ctx, ...)`
- **Melhor suporte da IDE**: o preenchimento automático funciona corretamente quando os métodos pertencem aos objetos
- **Mais clareza com várias janelas**: quando há várias janelas, você escolhe explicitamente em qual delas operar
- **Sem propagação de contexto**: você não precisa passar o contexto por todas as funções

### Bindings do frontend

No v2, os bindings eram organizados pelo pacote Go e pelo nome da struct, normalmente resultando em caminhos como `wailsjs/go/main/App`. Essa estrutura não refletia o agrupamento lógico e dificultava encontrar funcionalidades relacionadas.

O v3 organiza os bindings pelo nome do serviço e pelo módulo do aplicativo, criando uma estrutura lógica mais clara. Os bindings são gerados em um diretório `bindings`, organizado pelo nome do aplicativo e pelos nomes dos serviços, facilitando a compreensão das funcionalidades disponíveis.

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

**Por que isso é melhor:**

- **Organização lógica**: os bindings são agrupados pelo nome do serviço, e não pela estrutura de pacotes Go
- **Importações mais claras**: o caminho reflete a lógica do domínio (greetservice), e não a estrutura de arquivos (main/App)
- **Funcionalidades mais fáceis de encontrar**: você pode navegar pelos bindings por funcionalidade, em vez de pela estrutura técnica
- **Nomenclatura consistente**: a organização baseada em serviços corresponde à arquitetura do backend
- **Caminhos mais simples**: não há mais o prefixo `../wailsjs/go`, apenas `./bindings`

### Eventos

No v2, os eventos usavam parâmetros variádicos `interface{}` e exigiam que o contexto fosse passado para cada função de evento. Os manipuladores de eventos recebiam dados sem tipo definido, que exigiam asserções de tipo manuais, tornando o sistema de eventos propenso a erros e difícil de depurar.

O v3 introduz objetos de evento tipados e elimina a exigência de contexto. Os manipuladores recebem um objeto de evento adequado, com dados tipados, tornando o sistema de eventos mais confiável e fácil de usar.

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

**Por que isso é melhor:**

- **Segurança de tipos**: os eventos usam objetos de evento adequados em vez de `...interface{}`
- **Depuração aprimorada**: os objetos de evento contêm metadados, como o nome do evento, facilitando a depuração
- **API mais clara**: `app.Event.On()` e `app.Event.Emit()` são mais intuitivos do que as funções de runtime
- **Não requer contexto**: os eventos funcionam diretamente no objeto do aplicativo, sem propagar o contexto entre threads
- **Manipuladores mais simples**: os manipuladores de eventos têm uma assinatura clara em vez de parâmetros variádicos

### Janelas

v2 permitia apenas uma janela por aplicativo. A janela era criada na inicialização, e todas as operações de janela eram realizadas por funções de runtime que tinham essa única janela como destino implícito.

v3 introduz suporte nativo a várias janelas como um recurso central. Cada janela é um objeto de primeira classe, com métodos e ciclo de vida próprios. Você pode criar, gerenciar e destruir dinamicamente várias janelas durante todo o ciclo de vida do aplicativo.

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

**Por que isso é melhor:**

- **Aplicativos com várias janelas**: crie aplicativos com várias janelas independentes (painéis, preferências, ferramentas etc.)
- **Referências explícitas a janelas**: cada janela é um objeto que você pode armazenar e manipular diretamente
- **Criação dinâmica de janelas**: crie e destrua janelas a qualquer momento durante a execução
- **Estado independente por janela**: cada janela tem seus próprios eventos, propriedades e ciclo de vida
- **Arquitetura aprimorada**: o gerenciamento de janelas é orientado a objetos, em vez de baseado em contexto

## Etapas da migração

### Etapa 1: atualize as dependências

**go.mod:**

```go
module myapp

go 1.25.0

require (
    github.com/wailsapp/wails/v3 v3.0.0
)
```

**Atualize:**

```bash
go get github.com/wailsapp/wails/v3@v3.0.0
go mod tidy
```

### Etapa 2: atualize main.go

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

### Etapa 3: converta a struct do aplicativo em um serviço

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

### Etapa 4: atualize as chamadas de runtime

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

### Etapa 5: atualize o frontend

**Gere novos bindings:**

```bash
wails3 generate bindings
```

**Atualize as importações:**

```javascript
// v2
import { Greet } from '../wailsjs/go/main/App'

// v3
import { Greet } from './bindings/changeme/myservice'
```

**Atualize o tratamento de eventos:**

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

### Etapa 6: atualize a configuração

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

## Mapeamento de recursos

### Caixas de diálogo

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

### Menus

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

### Bandeja do sistema

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

## Problemas comuns

### Problema: bindings não encontrados

**Problema:** erros de importação após a migração

**Solução:**

```bash
# Regenerate bindings
wails3 generate bindings

# Check output directory
ls frontend/bindings
```

### Problema: erros de contexto

**Problema:** `ctx` não está disponível

**Solução:**

Em vez disso, armazene a referência ao aplicativo:

```go
type MyService struct {
    app *application.App
}

func NewMyService(app *application.App) *MyService {
    return &MyService{app: app}
}
```

### Problema: os métodos de janela não funcionam

**Problema:** `runtime.WindowSetTitle()` não existe

**Solução:**

Use os métodos da janela diretamente:

```go
window := s.app.Window.Current()
window.SetTitle("New Title")
```

### Problema: os eventos não são disparados

**Problema:** os eventos são registrados, mas não são recebidos

**Solução:**

Verifique se os nomes dos eventos correspondem exatamente:

```go
// Go
app.Event.Emit("my-event", data)

// JavaScript
OnEvent("my-event", handler)  // Must match exactly
```

## Teste da migração

### Lista de verificação

- [ ] O aplicativo inicia sem erros
- [ ] Todos os bindings funcionam
- [ ] Os eventos são enviados e recebidos
- [ ] As janelas abrem e fecham corretamente
- [ ] Os menus funcionam (se aplicável)
- [ ] As caixas de diálogo funcionam (se aplicável)
- [ ] A bandeja do sistema funciona (se aplicável)
- [ ] O processo de build funciona
- [ ] O build de produção funciona

### Comandos de teste

```bash
# Development
wails3 dev

# Build
wails3 build

# Generate bindings
wails3 generate bindings
```

## Benefícios do v3

### Desempenho

- **Inicialização mais rápida** — inicialização otimizada
- **Menor uso de memória** — uso eficiente de recursos
- **Bridge aprimorada** — overhead de &lt;1 ms por chamada

### Recursos

- **Várias janelas** — suporte nativo
- **Bandeja do sistema** — integrada
- **Eventos melhores** — API tipada e mais simples
- **Serviços** — Melhor organização do código

### Experiência de desenvolvimento

- **Segurança de tipos** — Suporte completo ao TypeScript
- **Erros melhores** — Mensagens de erro claras
- **Recarga automática** — Desenvolvimento mais rápido
- **Documentação melhor** — Guias abrangentes

## Como obter ajuda

### Recursos

- [Documentação](/quick-start/why-wails/)
- [Comunidade no Discord](https://discord.gg/JDdSxwjhGf)
- [Issues do GitHub](https://github.com/wailsapp/wails/issues)
- [Exemplos](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples)

### Perguntas frequentes

**P: Posso executar v2 e v3 lado a lado?** R: Sim, eles usam caminhos de importação diferentes.

**P: O v3 está pronto para produção?** O Wails v3.0.0 é a versão estável para aplicativos desktop no Windows, macOS e Linux. O suporte a Android e iOS continua experimental e não faz parte do contrato de compatibilidade estável para desktop.

**P: O v2 continuará recebendo manutenção?** R: Sim, o v2 receberá atualizações críticas.

**P: Quanto tempo leva a migração?** R: 1-4 horas para aplicativos típicos.

## Próximas etapas

@cards{cols="2"}
🚀 Início rápido
Comece a usar o Wails v3.

[Saiba mais →](/quick-start/installation/)

---
★ Conceitos fundamentais
Entenda a arquitetura do v3.

[Saiba mais →](/concepts/architecture/)

---
◆ Bindings
Conheça o novo sistema de bindings.

[Saiba mais →](/features/bindings/methods/)

---
📖 Exemplos
Veja exemplos completos do v3.

[Ver exemplos →](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples)

@end

---

**Dúvidas?** Pergunte no [Discord](https://discord.gg/JDdSxwjhGf) ou [abra uma issue](https://github.com/wailsapp/wails/issues).
