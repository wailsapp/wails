---
title: "Componentes internos do runtime"
description: "Uma análise detalhada de como o Wails v3 inicializa, executa e se comunica com o sistema operacional"
slug: "contributing/runtime-internals"
sourcePath: "contributing/runtime-internals.md"
---

O **runtime** é a camada que transforma funções Go comuns em um aplicativo desktop multiplataforma. Este documento explica os componentes que você encontrará ao acompanhar o código-fonte.

---

## 1. Ciclo de vida do aplicativo

| Fase | Caminho do código | O que acontece |
| --- | --- | --- |
| **Inicialização** | `pkg/application/application.go:init()` | Registra dados do momento da compilação e cria uma instância singleton global de `application`. |
| **New()** | `application.New(...)` | Valida `Options`, inicia o **AssetServer** e inicializa o sistema de logs. |
| **Run()** | `application.(*App).Run()` | 1. Chama o `mainthread.X()` da plataforma para entrar na thread da interface do sistema operacional.<br />2. Inicializa o **runtime** (`internal/runtime`).<br />3. Bloqueia até que a última janela seja fechada ou `Quit()` seja chamado. |
| **Encerramento** | `application.(*App).Quit()` | Transmite o evento `application:shutdown`, descarrega o log e encerra janelas e serviços. |

O ciclo de vida permite estritamente uma **única entrada**: você pode criar várias janelas, mas o objeto do aplicativo é inicializado apenas uma vez.

---

## 2. Gerenciamento de janelas

### API pública

```go
win := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "Dashboard",
    Width:  1280,
    Height: 720,
})
win.Show()
```

> `app.Window.New()` não recebe **nenhum argumento**; use `NewWithOptions(...)` quando você
>
> precisar passar uma struct `application.WebviewWindowOptions` (por valor).

`app.Window.New[WithOptions]()` delega para `pkg/application/webview_window_*.go`, onde ficam as implementações específicas de cada plataforma:

```
pkg/application/
├── webview_window_darwin.go    // WKWebView
├── webview_window_linux.go     // GTK + WebKitGTK (plus linux_cgo*.go)
└── webview_window_windows.go   // WebView2
```

Cada arquivo:

1. Cria uma webview nativa (WKWebView, WebKitGTK, WebView2).
2. Registra um callback do **Processador de Mensagens** (`pkg/application/messageprocessor*.go`).
3. Mapeia eventos do Wails (`WindowDidResize`, `WindowFocus`, `WindowFilesDropped`, …) para as constantes em `pkg/events`.

`internal/runtime/` é reservado para o pequeno código de integração das tags de compilação (`runtime{,_darwin,_linux,_windows,_android,_dev,_prod}.go`) e para o runtime JS incorporado em `internal/runtime/desktop/`.

As janelas ativas são rastreadas por `pkg/application/window_manager.go` / `webview_window.go`. `pkg/application/screenmanager.go` destina-se aos metadados de **tela** (resolução, escala e área de trabalho) — ele não gerencia janelas.

---

## 3. Pipeline de processamento de mensagens

A ponte entre JavaScript e Go é implementada pela família de **Processadores de Mensagens** em `pkg/application/messageprocessor_*.go`.

Fluxo:

1. O **JavaScript** chama `Call.ByID(<fnv-id>, ...args)` a partir de `/wails/runtime.js` (implementado em `internal/runtime/desktop/@wailsio/runtime/src/calls.ts`) — ou `Call.ByName("pkg.Struct.Method", ...args)` para compilações no modo de nomes.
2. O auxiliar do runtime empacota a chamada e a encaminha para o Go por meio da ponte nativa específica de cada plataforma.
3. O **Go** recebe a mensagem em `pkg/application/messageprocessor_call.go`.
4. O processador procura o método vinculado em `pkg/application/bindings.go` (escrito manualmente e baseado em `reflect`) e o invoca.
5. O resultado ou erro é serializado e enviado de volta ao JS, onde uma `Promise` é resolvida ou rejeitada.

> O envelope JSON exato é definido pelo auxiliar do runtime no lado do JS e
>
> por `messageprocessor_call.go` no lado do Go — versões anteriores desta página
>
> mencionavam um formato `{"t":"c","id":"123","m":"Greet","p":[…]}`, mas ele não
>
> corresponde à implementação atual. Leia os dois arquivos em conjunto ao investigar um
>
> bug no formato de transmissão.

Processadores especializados:

| Arquivo | Finalidade |
| --- | --- |
| `messageprocessor_window.go` | Ações de janela (ocultar, maximizar, …) |
| `messageprocessor_dialog.go` | Caixas de diálogo nativas (`OpenFile`, `MessageBox`, …) |
| `messageprocessor_clipboard.go` | Leitura/gravação da área de transferência |
| `messageprocessor_events.go` | Inscrição/emissão de eventos |
| `messageprocessor_browser.go` | Navegação no navegador e ferramentas de desenvolvimento |

Os processadores são **sem estado** — eles obtêm tudo de que precisam do `ApplicationContext` passado com cada mensagem.

---

## 4. Sistema de eventos

Os eventos são strings com namespace, encaminhadas por três camadas:

1. **Eventos do aplicativo**: ciclo de vida global (`application:ready`, `application:shutdown`).
2. **Eventos de janela**: por janela (`window:focus`, `window:resize`).
3. **Eventos personalizados**: definidos pelo usuário (`chat:new-message`).

Detalhes da implementação:

- As constantes de eventos ficam em `pkg/events/` (`defaults.go`, `known_events.go`, `events.txt`). Elas são geradas por `v3/tasks/events/generate.go` e expostas como `events.Common.*`, `events.Mac.*`, `events.Windows.*` e `events.Linux.*`. Você pode gerá-las novamente com `wails3 generate constants`.
- No lado do Go (eventos da aplicação):
  ```go
  app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {})
  ```

- No lado do Go (eventos de janela):
  ```go
  window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {})
  ```

- No lado do Go (eventos personalizados):
  ```go
  app.Event.On("chat:new-message", func(e *application.CustomEvent) {})
  ```

- No lado do JS:
  ```js
  import { Events } from "/wails/runtime.js";
  Events.On("chat:new-message", (e) => { /* … */ });
  ```


Os eventos da aplicação, de janela e personalizados passam por `pkg/application/event_manager.go`. As assinaturas de eventos de janela têm escopo restrito à respectiva janela; portanto, fechar a janela cancela automaticamente o registro de seus manipuladores.

---

## 5. Implementações específicas de plataforma

A compilação condicional mantém a API pública idêntica enquanto oculta as peculiaridades dos sistemas operacionais.

| Aspecto | Darwin | Linux | Windows |
| --- | --- | --- | --- |
| Thread principal | `mainthread_darwin.go` (Cgo para Foundation) | `mainthread_linux.go` (GTK) | `mainthread_windows.go` (Win32 `AttachThreadInput`) |
| Caixas de diálogo | `dialogs_darwin.*` (NSAlert) | `dialogs_linux.go` (GtkFileChooser) | `dialogs_windows.go` (IFileOpenDialog) |
| Área de transferência | `clipboard_darwin.go` | `clipboard_linux.go` | `clipboard_windows.go` |
| Ícones da bandeja do sistema | `systemtray_darwin.*` | `systemtray_linux.go` (DBus) | `systemtray_windows.go` (Shell_NotifyIcon) |

Princípios fundamentais:

- **macOS e Windows** usam Cgo com moderação (principalmente por meio de `pkg/mac/` e dos wrappers Win32 `w32` em `pkg/w32`).
- **O Linux usa Cgo intensivamente** por necessidade — `pkg/application/linux_cgo.go` (~69 KB) e `linux_cgo_gtk4.{c,go,h}` (~50 KB+) controlam GTK/WebKitGTK diretamente.
- Use **tags de compilação** (`//go:build darwin`, `//go:build linux`, …) para manter legíveis os arquivos específicos de cada sistema operacional.
- `internal/capabilities/` existe para sinalizadores de recursos específicos de cada plataforma, mas o framework **não** exporta uma sentinela `ErrCapability` — a habilitação condicional de recursos é feita por retornos de stubs específicos da plataforma.

---

## 6. Guia de arquivos

| Arquivo | Por que você o modificaria |
| --- | --- |
| `internal/runtime/runtime_*.go` | Alterar a pequena camada de stubs com tags de compilação (desenvolvimento versus produção e integração específica do sistema operacional). |
| `pkg/application/webview_window_*.go` | Implemente uma nova indicação ao gerenciador de janelas ou um novo comportamento de janela. |
| `pkg/application/messageprocessor*.go` | Adicionar um novo comando da ponte que possa ser chamado pelo JS. |
| `pkg/events/*.go` | Estender as definições de eventos integradas (depois, execute `wails3 generate constants` novamente). |
| `internal/assetserver/*` | Ajustar o processamento de ativos em desenvolvimento/produção. |
| `internal/runtime/desktop/@wailsio/runtime/src/*` | Editar o runtime JS incorporado (despacho de chamadas/eventos, caixas de diálogo, arraste etc.). |

---

## 7. Dicas de depuração

- Configure `Options.LogLevel` (por exemplo, `slog.LevelDebug`) e inspecione a saída de `Options.Logger` — não existe uma variável de ambiente `WAILS_LOG_LEVEL`.
- Os sinalizadores de `wails3 dev` são `--config`, `--port`, `-s` (habilita HTTPS) e o `--no-colour` global. Não existe um sinalizador `-verbose`.
- No macOS, execute sob `lldb --` para detectar exceções do Objective-C antecipadamente.
- Para problemas com o Chromium no Windows, habilite os logs de depuração do WebView2: `set WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS=--remote-debugging-port=9222`

---

## 8. Extensão do runtime

1. Opcional: declare qualquer novo sinalizador de recurso em `internal/capabilities/`.
2. Implemente o recurso em cada variante de `pkg/application/*_{darwin,linux,windows}.go` usando tags de compilação. Forneça um stub nas plataformas que não oferecem suporte a ele.
3. Adicione a API pública em `pkg/application` (interface + métodos concretos de `WebviewWindow`, struct de opções etc.).
4. Registre um novo método do processador de mensagens (`pkg/application/messageprocessor*.go`) e um helper correspondente no runtime JS, caso o JS precise chamá-lo.
5. Se você adicionar um evento, declare a constante dele em `pkg/events/` e execute `wails3 generate constants` para atualizar os arquivos gerados.

Siga esta lista de verificação para manter intacto o contrato multiplataforma.

---

## 9. Arrastar e soltar

O recurso de arrastar e soltar arquivos usa uma **abordagem centrada em JavaScript** em todas as plataformas. A camada nativa intercepta os eventos de arrastar do sistema operacional, mas o processamento efetivo da soltura e a interação com o DOM ocorrem no JavaScript.

### Fluxo

1. O usuário arrasta arquivos do sistema operacional sobre a janela do Wails
2. A camada nativa detecta o arraste e notifica o JavaScript para aplicar os efeitos visuais ao passar sobre uma área
3. O usuário solta os arquivos
4. A camada nativa envia os caminhos dos arquivos e as coordenadas ao JavaScript
5. O JavaScript encontra o elemento de destino da soltura (`data-file-drop-target`)
6. O JavaScript envia os caminhos dos arquivos e os detalhes do elemento ao backend em Go
7. O Go emite o evento `WindowFilesDropped` com o contexto completo

### Implementações específicas de cada plataforma

| Plataforma | Camada nativa | Principal desafio |
| --- | --- | --- |
| **Windows** | Suporte integrado do WebView2 a arrastar | Coordenadas em pixels CSS, sem necessidade de conversão |
| **macOS** | Delegados de arraste do NSWindow | Converter coordenadas relativas à janela em coordenadas relativas à webview |
| **Linux** | Sinais de arraste do GTK3 | É necessário distinguir o arraste de arquivos do arraste interno do HTML5 |

### Linux: distinção entre tipos de arraste

Tanto o GTK quanto o WebKit querem processar os eventos de arraste. O segredo é verificar o tipo do destino do arraste:

```c
static gboolean is_file_drag(GdkDragContext *context) {
    GList *targets = gdk_drag_context_list_targets(context);
    for (GList *l = targets; l != NULL; l = l->next) {
        GdkAtom atom = GDK_POINTER_TO_ATOM(l->data);
        gchar *name = gdk_atom_name(atom);
        if (name && g_strcmp0(name, "text/uri-list") == 0) {
            g_free(name);
            return TRUE;  // External file drag
        }
        g_free(name);
    }
    return FALSE;  // Internal HTML5 drag
}
```

Os manipuladores de sinais retornam `FALSE` para arrastes internos (permitindo que o WebKit os processe) e `TRUE` para arrastes de arquivos (que nós mesmos processamos).

### Bloqueio da soltura de arquivos

Quando `EnableFileDrop` é `false`, ainda precisamos impedir que o navegador abra os arquivos soltos. Cada plataforma trata isso de maneira diferente:

- **Windows**: o JavaScript chama `preventDefault()` nos eventos de arraste
- **macOS**: o JavaScript chama `preventDefault()` nos eventos de arraste\
- **Linux**: os manipuladores de sinais do GTK interceptam e rejeitam o arraste de arquivos na camada nativa

### Arquivos principais

| Arquivo | Finalidade |
| --- | --- |
| `pkg/application/linux_cgo.go` | Manipuladores de sinais de arraste do GTK (código C no preâmbulo do cgo) |
| `pkg/application/webview_window_darwin.go` | Delegados de arraste do macOS |
| `pkg/application/webview_window_windows.go` | Processamento de mensagens do WebView2 |
| `internal/runtime/desktop/@wailsio/runtime/src/window.ts` | Processamento da soltura no JavaScript |

### Depuração

- **Linux**: adicione `printf` ao código C (lembre-se de `fflush(stdout)`)
- **Windows**: use `globalApplication.debug()`
- **JavaScript**: verifique o console do navegador e ative o modo de depuração

Problemas comuns:

1. **O arraste interno do HTML5 não funciona**: o manipulador nativo o está interceptando (retorne `FALSE` para arrastes que não sejam de arquivos)
2. **Os efeitos visuais ao passar sobre um elemento não aparecem**: os manipuladores JavaScript não estão sendo chamados
3. **Coordenadas incorretas**: verifique as conversões entre sistemas de coordenadas

---

Agora você concluiu uma visita guiada aos componentes internos do runtime. Combine esse conhecimento com o mapa **Estrutura da base de código** e a documentação do **Servidor de ativos** para navegar com confiança e fazer contribuições relevantes. Bom trabalho!
