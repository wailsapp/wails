---
title: "Ciclo de vida da aplicação"
description: "Entenda o ciclo de vida da aplicação Wails, da inicialização ao encerramento"
slug: "concepts/lifecycle"
sourcePath: "concepts/lifecycle.md"
---

## Entendendo o ciclo de vida da aplicação

As aplicações desktop têm um ciclo de vida que vai da inicialização ao encerramento. O Wails v3 fornece **serviços**, **eventos** e **hooks** para gerenciar esse ciclo de vida com eficiência.

## Etapas do ciclo de vida

```d2
direction: down

Start: Início do aplicativo {
  shape: oval
  style.fill: "#10B981"
}

Init: Inicialização {
  Parse: Analisar opções {
    shape: rectangle
  }
  Register: Registrar serviços {
    shape: rectangle
  }
  Setup: Configurar o runtime {
    shape: rectangle
  }
}

AppRun: app.Run() {
  shape: rectangle
  style.fill: "#3B82F6"
}

ServiceStartup: Inicialização dos serviços {
  shape: rectangle
  style.fill: "#8B5CF6"
}

EventLoop: Loop de eventos {
  Process: Processar eventos {
    shape: rectangle
  }
  Handle: Tratar mensagens {
    shape: rectangle
  }
  Update: Atualizar a interface {
    shape: rectangle
  }
}

QuitSignal: Sinal de encerramento {
  shape: diamond
  style.fill: "#F59E0B"
}

ShouldQuit: Verificação de ShouldQuit {
  shape: rectangle
  style.fill: "#3B82F6"
}

OnShutdown: Callbacks de OnShutdown {
  shape: rectangle
  style.fill: "#3B82F6"
}

ServiceShutdown: Encerramento dos serviços {
  shape: rectangle
  style.fill: "#8B5CF6"
}

Cleanup: Limpeza {
  Close: Fechar janelas {
    shape: rectangle
  }
  Release: Liberar recursos {
    shape: rectangle
  }
}

End: Fim do aplicativo {
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
EventLoop.Update -> EventLoop.Process: Loop
EventLoop.Process -> QuitSignal: Usuário encerra
QuitSignal -> ShouldQuit: Verificar se é permitido?
ShouldQuit -> EventLoop.Process: Negado
ShouldQuit -> OnShutdown: Permitido
OnShutdown -> ServiceShutdown
ServiceShutdown -> Cleanup.Close
Cleanup.Close -> Cleanup.Release
Cleanup.Release -> End
```

### 1. Criação da aplicação

Crie sua aplicação com `application.New()`:

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

**O que acontece:**

1. As opções são analisadas e validadas
2. Os serviços são registrados (mas ainda não são iniciados)
3. O servidor de assets é configurado
4. O runtime é configurado

### 2. Execução da aplicação

Chame `app.Run()` para iniciar a aplicação:

```go
err := app.Run()  // Blocks until quit
if err != nil {
    log.Fatal(err)
}
```

**O que acontece:**

1. Os serviços são iniciados na ordem de registro
2. Os listeners de eventos são ativados
3. As janelas podem ser criadas
4. O loop de eventos é iniciado

### 3. Loop de eventos

A aplicação entra no loop de eventos, no qual passa a maior parte do tempo:

- Os eventos do sistema operacional são processados (eventos de mouse, teclado e janela)
- As mensagens de Go para JS são tratadas
- As chamadas de JS para Go são executadas
- As atualizações da interface são renderizadas

### 4. Encerramento

Quando a aplicação é encerrada:

1. O callback `ShouldQuit` é verificado (se estiver definido)
2. Os callbacks `OnShutdown` são executados
3. Os serviços são encerrados na ordem inversa
4. As janelas são fechadas
5. Os recursos são liberados

## Ciclo de vida dos serviços

Os serviços são a principal forma de gerenciar o ciclo de vida no Wails v3. Eles fornecem hooks de inicialização e encerramento por meio de interfaces. Para consultar a documentação completa sobre serviços, consulte o [guia de serviços](/features/bindings/services/).

### Criando um serviço

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

### Registrando serviços

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&MyService{}),
        application.NewService(&AnotherService{}),
    },
})
```

**Pontos principais:**

- Os serviços são iniciados na ordem de registro
- Os serviços são encerrados na ordem **inversa** à de registro
- Se `ServiceStartup` de um serviço retornar um erro, a aplicação será abortada
- O `ctx` passado para `ServiceStartup` é cancelado quando o encerramento começa

### Usando o contexto da aplicação

O contexto passado para `ServiceStartup` é válido durante todo o ciclo de vida do aplicativo:

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

Você também pode acessar o contexto pela instância do aplicativo:

```go
app := application.Get()
ctx := app.Context()
```

## Hooks no nível do aplicativo

Esses são callbacks de conveniência em `application.Options` que permitem interagir com o ciclo de vida do aplicativo sem criar um serviço completo. Eles são úteis para tarefas simples de limpeza, confirmação de encerramento ou quando você precisa executar código em pontos específicos da sequência de desligamento.

Para um gerenciamento mais complexo do ciclo de vida, com lógica de inicialização, injeção de dependências ou recursos com estado, use [Serviços](#ciclo-de-vida-dos-servios).

### ShouldQuit

O callback `ShouldQuit` é chamado sempre que o encerramento é solicitado, seja quando o usuário fecha a última janela, pressiona Cmd+Q (macOS) ou Alt+F4 (Windows), seja quando `app.Quit()` é chamado programaticamente.

**Valor de retorno:**

- Retorne `true` para permitir que o encerramento prossiga (o aplicativo será desligado)
- Retorne `false` para cancelar o encerramento (o aplicativo continuará em execução)

Essa é a oportunidade de interceptar solicitações de encerramento e, opcionalmente, impedi-las, por exemplo, para avisar o usuário sobre alterações não salvas:

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

Se `ShouldQuit` não estiver definido, o aplicativo será encerrado imediatamente quando solicitado.

**Quando ShouldQuit é chamado:**

- O usuário fecha a última janela (a menos que `DisableQuitOnLastWindowClosed` esteja definido)
- O usuário pressiona Cmd+Q no macOS
- O usuário pressiona Alt+F4 no Windows (quando o foco está na última janela)
- O código chama `app.Quit()`

**Quando ShouldQuit NÃO é chamado:**

- O processo é encerrado à força (SIGKILL ou encerramento forçado pelo Gerenciador de Tarefas)
- `os.Exit()` é chamado diretamente

### OnShutdown

O callback `OnShutdown` é chamado quando se confirma que o aplicativo será encerrado (depois que `ShouldQuit` retorna `true`, caso esteja definido). Use-o para tarefas de limpeza, como salvar o estado, fechar conexões com o banco de dados ou liberar recursos.

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

Você também pode registrar callbacks adicionais de desligamento programaticamente a qualquer momento durante o ciclo de vida do aplicativo:

```go
app.OnShutdown(func() {
    log.Println("Application shutting down...")
})
```

Vários callbacks são executados na ordem em que foram registrados. O processo de desligamento fica bloqueado até que todos os callbacks sejam concluídos.

**Importante:** Mantenha os callbacks de desligamento rápidos (menos de 1 segundo). O sistema operacional pode encerrar à força aplicativos que demoram demais para fechar, o que pode interromper a limpeza e causar perda de dados.

### PostShutdown

O callback `PostShutdown` é chamado depois que todas as tarefas de desligamento são concluídas, pouco antes do término do processo. Nesse ponto, a instância do aplicativo não pode mais ser usada: as janelas estão fechadas, os serviços foram desligados e os recursos foram liberados.

Isso é útil principalmente para:

- Registro final de logs que precisa ocorrer depois de todas as outras tarefas de limpeza
- Testes e depuração do comportamento de desligamento
- Plataformas nas quais `app.Run()` não retorna (o callback garante que seu código seja executado)

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

**Observação:** Não tente usar recursos do aplicativo (janelas, caixas de diálogo etc.) em `PostShutdown`, pois eles não estão mais disponíveis.

## Ciclo de vida baseado em eventos

O Wails fornece um sistema de eventos que notifica você quando algo acontece no aplicativo, como a abertura de janelas, a inicialização do aplicativo, alterações de tema e muito mais. Você pode escutar esses eventos para reagir a mudanças no ciclo de vida sem bloqueá-las nem interceptá-las.

Para eventos de janela, você também pode usar `RegisterHook` em vez de `OnWindowEvent` para interceptar e cancelar ações, por exemplo, para impedir que uma janela seja fechada. Consulte [Hooks de janela](#hooks-de-janela-eventos-cancelveis) abaixo.

Para consultar a documentação completa do sistema de eventos, consulte o [guia de eventos](/features/events/system/).

### Eventos do aplicativo

Escute os eventos do ciclo de vida do aplicativo:

```go
app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(event *application.ApplicationEvent) {
    app.Logger.Info("Application has started!")
})
```

Também há eventos específicos de cada plataforma:

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

### Eventos de janela

Escute os eventos do ciclo de vida das janelas:

```go
window := app.Window.New()

window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window gained focus")
})

window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
    app.Logger.Info("Window is closing")
})
```

### Hooks de janela (eventos canceláveis)

Use `RegisterHook` em vez de `OnWindowEvent` quando precisar **cancelar** um evento:

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

**Diferença entre OnWindowEvent e RegisterHook:**

- `OnWindowEvent`: notifica você quando ocorre um evento (não permite cancelá-lo)
- `RegisterHook`: permite interceptar e, possivelmente, cancelar o evento

## Ciclo de vida das janelas

As janelas têm seu próprio ciclo de vida, desde a criação até a destruição. Cada janela carrega o conteúdo do frontend de forma independente e pode ser exibida, ocultada ou fechada a qualquer momento. Quando um usuário tenta fechar uma janela, você pode interceptar essa ação com um `RegisterHook` para solicitar confirmação ou ocultar a janela em vez de destruí-la.

Para consultar a documentação completa sobre janelas, veja o [guia de janelas](/features/windows/basics/).

```d2
direction: down

Create: Criar janela {
  shape: oval
  style.fill: "#10B981"
}

Load: Carregar o frontend {
  shape: rectangle
}

Show: Exibir janela {
  shape: rectangle
}

Active: Janela ativa {
  Events: Tratar eventos {
    shape: rectangle
  }
}

CloseRequest: Solicitação de fechamento {
  shape: diamond
  style.fill: "#F59E0B"
}

Hook: Hook WindowClosing {
  shape: rectangle
  style.fill: "#3B82F6"
}

Destroy: Destruir janela {
  shape: rectangle
}

End: Janela fechada {
  shape: oval
  style.fill: "#EF4444"
}

Create -> Load
Load -> Show
Show -> Active.Events
Active.Events -> Active.Events: Loop
Active.Events -> CloseRequest: Usuário fecha
CloseRequest -> Hook
Hook -> Active.Events: Cancelado
Hook -> Destroy: Permitido
Destroy -> End
```

### Criação de janelas

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "My Window",
    Width:  800,
    Height: 600,
})
```

### Como impedir o fechamento da janela

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

### Ocultar em vez de fechar

Um padrão comum para aplicativos da bandeja do sistema:

```go
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    window.Hide()  // Hide instead of destroy
    e.Cancel()     // Prevent actual close
})
```

## Ciclo de vida com várias janelas

Com várias janelas:

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

**O comportamento padrão varia conforme a plataforma:**

| Plataforma | Comportamento padrão ao fechar a última janela |
| --- | --- |
| macOS | O aplicativo continua em execução (a barra de menus permanece) |
| Windows | O aplicativo é encerrado |
| Linux | O aplicativo é encerrado |

O macOS segue as convenções nativas da plataforma, nas quais os aplicativos normalmente permanecem ativos na barra de menus mesmo sem nenhuma janela. No Windows e no Linux, eles são encerrados por padrão.

**Fazer o aplicativo ser encerrado em todas as plataformas quando a última janela for fechada:**

```go
app := application.New(application.Options{
    Mac: application.MacOptions{
        ApplicationShouldTerminateAfterLastWindowClosed: true,
    },
})
```

**Fazer o aplicativo continuar em execução em todas as plataformas quando a última janela for fechada:**

Isso é útil para aplicativos da bandeja do sistema ou aplicativos que devem continuar em execução em segundo plano.

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

## Padrões comuns

### Padrão 1: serviço de banco de dados

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

### Padrão 2: serviço de configuração

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

### Padrão 3: worker em segundo plano

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

## Referência do ciclo de vida

| Hook/interface | Quando é chamado | Pode cancelar? | Finalidade |
| --- | --- | --- | --- |
| `ServiceStartup` | Durante `app.Run()`, antes do loop de eventos | Não (retorne um erro para interromper) | Inicialização |
| `ServiceShutdown` | Durante o encerramento, após `OnShutdown` | Não | Limpeza |
| `OnShutdown` | Quando o encerramento é confirmado | Não | Limpeza do aplicativo |
| `ShouldQuit` | Quando o encerramento é solicitado | Sim (retorne false) | Confirmar encerramento |
| `RegisterHook(WindowClosing)` | Quando o fechamento da janela é solicitado | Sim (`e.Cancel()`) | Impedir o fechamento da janela |
| `OnWindowEvent` | Quando o evento ocorre | Não | Reagir a eventos |
| `OnApplicationEvent` | Quando o evento ocorre | Não | Reagir a eventos |

## Diferenças entre plataformas

### macOS

- O **menu do aplicativo** permanece disponível mesmo sem nenhuma janela
- **Cmd+Q** inicia o encerramento do aplicativo (passando por `ShouldQuit`)
- O **ícone no Dock** permanece visível, a menos que seja ocultado
- Use `ApplicationShouldTerminateAfterLastWindowClosed` para controlar o comportamento de encerramento

### Windows

- **Não há menu do aplicativo** sem uma janela
- **Alt+F4** fecha a janela (isso pode ser impedido com `RegisterHook`)
- A **bandeja do sistema** pode manter o aplicativo em execução

### Linux

- O **comportamento varia** de acordo com o ambiente de desktop
- **Geralmente semelhante ao Windows**

## Depuração de problemas do ciclo de vida

### Problema: o aplicativo não encerra

**Causas:**

1. `ShouldQuit` retornando `false`
2. `OnShutdown` demorando demais
3. Goroutines em segundo plano não são interrompidas

**Solução:**

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

### Problema: falha na inicialização do serviço

**Solução:** Retorne erros descritivos:

```go
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    if err := s.init(); err != nil {
        return fmt.Errorf("failed to initialise: %w", err)
    }
    return nil
}
```

O erro será registrado, e o aplicativo não será iniciado.

## Práticas recomendadas

### Faça

- **Use serviços para gerenciar o ciclo de vida** — Eles fornecem hooks adequados de inicialização e encerramento
- **Mantenha o encerramento rápido** — Procure concluir toda a limpeza em menos de 1 segundo
- **Use o contexto para cancelamento** — Interrompa corretamente as tarefas em segundo plano
- **Trate erros durante a inicialização** — Retorne os erros para interromper o processo de forma segura
- **Registre os eventos do ciclo de vida** — Isso ajuda na depuração

### Não faça

- **Não bloqueie durante a inicialização do serviço** — Mantenha a inicialização rápida (menos de 2 segundos)
- **Não exiba caixas de diálogo durante o encerramento** — O aplicativo está sendo encerrado, e a interface pode não funcionar
- **Não ignore o contexto** — Sempre verifique `ctx.Done()` nas goroutines
- **Não cause vazamento de recursos** — Sempre implemente `ServiceShutdown`

## Próximas etapas

**Serviços** — Saiba mais sobre o sistema de serviços [Saiba mais →](/features/bindings/services/)

**Sistema de eventos** — Use eventos para comunicação [Saiba mais →](/features/events/system/)

**Gerenciamento de janelas** — Crie e gerencie janelas [Saiba mais →](/features/windows/basics/)

---

**Tem dúvidas sobre o ciclo de vida?** Pergunte no [Discord](https://discord.gg/JDdSxwjhGf) ou consulte os [exemplos](https://github.com/wailsapp/wails/tree/master/v3/examples).
