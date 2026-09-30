---
title: "Migração da v2 para a v3"
description: "Guia completo para migrar seu aplicativo Wails da v2 para a v3"
slug: "migration/v2-to-v3"
sourcePath: "migration/v2-to-v3.md"
---

O Wails v3 é uma **reescrita completa**, com melhorias significativas na arquitetura, no desempenho e na experiência de desenvolvimento. Este guia ajuda você a migrar seu aplicativo da v2 para a v3.

**Principais mudanças:**

- Nova estrutura de aplicativos
- Sistema de bindings aprimorado
- Gerenciamento de janelas aprimorado
- Sistema de eventos aprimorado
- Configuração simplificada

**Tempo de migração:** 1-4 horas para aplicativos típicos

## Mudanças incompatíveis

### Inicialização do aplicativo

Na v2, a configuração do aplicativo, a configuração da janela e a execução eram combinadas em uma única chamada a `wails.Run()`. Essa abordagem monolítica dificultava criar várias janelas, tratar erros em diferentes etapas ou testar componentes individuais do aplicativo.

A v3 separa essas responsabilidades em fases distintas: criação do aplicativo, criação das janelas e execução. Essa separação oferece controle explícito sobre cada etapa do ciclo de vida do aplicativo e torna o código mais modular e testável.

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

Na v2, cada struct vinculada exigia um campo de contexto e um método `startup(ctx)` para receber o contexto do runtime. Isso criava um forte acoplamento entre a lógica de negócios e o runtime do Wails, tornando o código mais difícil de testar e compreender.

A v3 introduz o padrão de serviços, no qual suas structs são totalmente independentes e não precisam armazenar o contexto do runtime. Se um serviço precisar acessar a instância do aplicativo, ele a receberá explicitamente por injeção de dependência, em vez de receber o contexto implicitamente ao longo da cadeia de chamadas.

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
- **Código mais claro**: as dependências são explícitas (passadas como argumentos do construtor), em vez de ficarem ocultas em um campo de contexto
- **Melhor organização**: os serviços podem ser agrupados por domínio, em vez de ficarem todos em uma única struct `App`
- **Inicialização adequada**: quando precisar de inicialização, use o método `ServiceStartup()` para torná-la explícita

### Runtime

Na v2, todas as operações do runtime exigiam passar um contexto para funções globais do pacote `runtime`. Isso criava um forte acoplamento com o objeto de contexto em toda a base de código e fazia a API parecer procedural, em vez de orientada a objetos.

A v3 substitui o runtime baseado em contexto por chamadas diretas a métodos dos objetos de aplicativo e janela. As operações são chamadas diretamente nos objetos que afetam, tornando o código mais intuitivo e orientado a objetos.

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

- **Design orientado a objetos**: os métodos são chamados nos objetos que afetam (janela, aplicativo, menu etc.)
- **Intenção mais clara**: `window.SetTitle()` é mais evidente do que `runtime.WindowSetTitle(ctx, ...)`
- **Melhor suporte da IDE**: o preenchimento automático funciona corretamente quando os métodos pertencem aos objetos
- **Mais clareza ao usar várias janelas**: quando há várias janelas, você escolhe explicitamente em qual delas a operação será realizada
- **Sem propagação de contexto**: você não precisa passar o contexto por todas as funções

### Bindings do frontend

Na v2, os bindings eram organizados pelo pacote Go e pelo nome da struct, geralmente resultando em caminhos como `wailsjs/go/main/App`. Essa estrutura não refletia um agrupamento lógico e dificultava encontrar funcionalidades relacionadas.

A v3 organiza os bindings pelo nome do serviço e pelo módulo do aplicativo, criando uma estrutura lógica mais clara. Os bindings são gerados em um diretório `bindings`, organizado pelo nome do aplicativo e pelos nomes dos serviços, facilitando a compreensão das funcionalidades disponíveis.

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

- **Organização lógica**: os bindings são agrupados pelo nome do serviço, em vez de pela estrutura de pacotes Go
- **Importações mais claras**: o caminho reflete a lógica do domínio (greetservice), não a estrutura de arquivos (main/App)
- **Mais facilidade para encontrar recursos**: Você pode navegar pelos bindings por funcionalidade, em vez de pela estrutura técnica
- **Nomenclatura consistente**: A organização baseada em serviços corresponde à arquitetura do seu backend
- **Caminhos mais simples**: Não é mais necessário usar o prefixo `../wailsjs/go` — basta usar `./bindings`

### Eventos

Na v2, os eventos usavam parâmetros variádicos `interface{}` e exigiam que o contexto fosse passado para todas as funções de evento. Os manipuladores de eventos recebiam dados sem tipo definido que exigiam asserções de tipo manuais, tornando o sistema de eventos propenso a erros e difícil de depurar.

A v3 introduz objetos de evento tipados e elimina a exigência de contexto. Os manipuladores de eventos recebem um objeto de evento adequado com dados tipados, tornando o sistema de eventos mais confiável e fácil de usar.

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

- **Segurança de tipos**: Os eventos usam objetos de evento adequados em vez de `...interface{}`
- **Depuração aprimorada**: Os objetos de evento contêm metadados, como o nome do evento, facilitando a depuração
- **API mais clara**: `app.Event.On()` e `app.Event.Emit()` são mais intuitivos do que as funções de runtime
- **Não requer contexto**: Os eventos operam diretamente no objeto da aplicação, sem precisar propagar o contexto
- **Manipuladores mais simples**: Os manipuladores de eventos têm uma assinatura clara em vez de parâmetros variádicos

### Janelas

A v2 permitia apenas uma janela por aplicação. A janela era criada na inicialização, e todas as operações de janela eram realizadas por funções de runtime que tinham essa única janela como alvo implícito.

A v3 introduz suporte nativo a várias janelas como funcionalidade central. Cada janela é um objeto de primeira classe, com métodos e ciclo de vida próprios. Você pode criar, gerenciar e destruir dinamicamente várias janelas durante todo o ciclo de vida da aplicação.

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

- **Aplicações com várias janelas**: Crie aplicações com várias janelas independentes (painéis, preferências, ferramentas etc.)
- **Referências explícitas a janelas**: Cada janela é um objeto que você pode armazenar e manipular diretamente
- **Criação dinâmica de janelas**: Crie e destrua janelas a qualquer momento durante a execução
- **Estado independente por janela**: Cada janela tem eventos, propriedades e ciclo de vida próprios
- **Arquitetura aprimorada**: O gerenciamento de janelas é orientado a objetos, em vez de baseado em contexto

## Etapas da migração

### Etapa 1: Atualizar as dependências

**go.mod:**

```go
module myapp

go 1.25.0

require (
    github.com/wailsapp/wails/v3 v3.0.0-beta.0
)
```

**Atualize:**

```bash
go get github.com/wailsapp/wails/v3@latest
go mod tidy
```

### Etapa 2: Atualizar main.go

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

### Etapa 3: Converter a struct App em um serviço

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

### Etapa 4: Atualizar as chamadas de runtime

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

### Etapa 5: Atualizar o frontend

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

### Etapa 6: Atualizar a configuração

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

## Mapeamento de funcionalidades

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

### Área de notificação

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

**Problema:** Erros de importação após a migração

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

Em vez disso, armazene uma referência ao aplicativo:

```go
type MyService struct {
    app *application.App
}

func NewMyService(app *application.App) *MyService {
    return &MyService{app: app}
}
```

### Problema: os métodos da janela não funcionam

**Problema:** `runtime.WindowSetTitle()` não existe

**Solução:**

Use os métodos da janela diretamente:

```go
window := s.app.Window.Current()
window.SetTitle("New Title")
```

### Problema: os eventos não são disparados

**Problema:** Os eventos são registrados, mas não são recebidos

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
- [ ] Todas as vinculações funcionam
- [ ] Os eventos são enviados e recebidos
- [ ] As janelas abrem e fecham corretamente
- [ ] Os menus funcionam (se aplicável)
- [ ] As caixas de diálogo funcionam (se aplicável)
- [ ] A bandeja do sistema funciona (se aplicável)
- [ ] O processo de compilação funciona
- [ ] A compilação para produção funciona

### Comandos de teste

```bash
# Development
wails3 dev

# Build
wails3 build

# Generate bindings
wails3 generate bindings
```

## Benefícios da v3

### Desempenho

- **Inicialização mais rápida** — Inicialização otimizada
- **Menor uso de memória** — Uso eficiente de recursos
- **Ponte aprimorada** — Sobrecarga de &lt;1 ms por chamada

### Recursos

- **Múltiplas janelas** — Suporte nativo
- **Bandeja do sistema** — Integrada
- **Eventos aprimorados** — API tipada e mais simples
- **Serviços** — Melhor organização do código

### Experiência de desenvolvimento

- **Segurança de tipos** — Suporte completo a TypeScript
- **Erros aprimorados** — Mensagens de erro claras
- **Recarregamento automático** — Desenvolvimento mais rápido
- **Documentação aprimorada** — Guias abrangentes

## Como obter ajuda

### Recursos

- [Documentação](/quick-start/why-wails/)
- [Comunidade no Discord](https://discord.gg/JDdSxwjhGf)
- [Issues do GitHub](https://github.com/wailsapp/wails/issues)
- [Exemplos](https://github.com/wailsapp/wails/tree/master/v3/examples)

### Perguntas frequentes

**P: Posso executar a v2 e a v3 lado a lado?** R: Sim, elas usam caminhos de importação diferentes.

**P: A v3 está pronta para produção?** R: A v3 é um software beta com uma API estável para desktop. Há aplicativos em execução em produção com ela, mas faça testes completos antes da implantação. A v2 continua sendo a versão estável atual.

**P: A v2 continuará recebendo manutenção?** R: Sim, a v2 receberá atualizações críticas.

**P: Quanto tempo leva a migração?** R: 1-4 horas para aplicativos típicos.

## Próximas etapas

@cards{cols="2"}
🚀 Início rápido
Comece a usar o Wails v3.

[Saiba mais →](/quick-start/installation/)

---
★ Conceitos fundamentais
Entenda a arquitetura da v3.

[Saiba mais →](/concepts/architecture/)

---
◆ Vinculações
Conheça o novo sistema de vinculações.

[Saiba mais →](/features/bindings/methods/)

---
📖 Exemplos
Veja exemplos completos da v3.

[Ver exemplos →](https://github.com/wailsapp/wails/tree/master/v3/examples)

@end

---

**Dúvidas?** Pergunte no [Discord](https://discord.gg/JDdSxwjhGf) ou [abra uma issue](https://github.com/wailsapp/wails/issues).
