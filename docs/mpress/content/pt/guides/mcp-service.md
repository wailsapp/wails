---
title: "Controle por LLM (MCP)"
description: "Permita que agentes de LLM testem e controlem seu aplicativo por meio do Model Context Protocol"
slug: "guides/mcp-service"
sourcePath: "guides/mcp-service.md"
---

@note{type="caution" title="Recurso experimental"}
O servidor MCP integrado é experimental, e a API pode mudar em versões futuras.

@end

O Wails v3 tem um servidor integrado do [Model Context Protocol](https://modelcontextprotocol.io) (MCP) que permite que agentes de LLM — Claude Code, assistentes de IDE ou qualquer cliente MCP — inspecionem, testem e controlem um aplicativo Wails **em execução**.

Para automatizar o ciclo de vida do projeto, use o servidor de CLI `wails3 mcp` separado. Ele permite que agentes inspecionem e inicializem projetos, executem diagnósticos, iniciem builds e tarefas de desenvolvimento, gerem bindings, executem tarefas nomeadas do Taskfile e obtenham uma quantidade limitada da saída das tarefas. Por padrão, o servidor de CLI fica restrito ao diretório atual e não permite a execução arbitrária de comandos do shell. Consulte a [documentação do MCP da CLI](/guides/cli/#mcp) para obter detalhes sobre transporte, autenticação e ferramentas.

Quando esse recurso está habilitado, um agente conectado ao seu aplicativo pode:

- **Listar e controlar janelas** — tamanho, posição, foco, tela cheia, ferramentas de desenvolvimento, recarregamento, …
- **Inspecionar o DOM** — consultar elementos, obter o HTML e gerar um instantâneo estrutural
- **Avaliar JavaScript** — executar código arbitrário dentro de qualquer janela e obter o resultado
- **Simular entradas do usuário** — movimentos do mouse, cliques, operações de arrastar e rolagens, exibidos com um **cursor animado na tela** para que você possa acompanhar o trabalho do agente
- **Digitar e pressionar teclas** — eventos realistas para cada caractere, compatíveis com entradas controladas do React
- **Chamar métodos Go vinculados** e emitir ou aguardar eventos do aplicativo

## Como funciona

O servidor MCP só é compilado no aplicativo quando a tag de build **`mcp`** está presente. Sem a tag, o código do servidor fica completamente ausente do binário — sem sobrecarga em tempo de execução, sem portas abertas e sem superfície de ataque.

Quando a tag está presente, o servidor é iniciado automaticamente dentro de `App.Run()`, associa-se a `127.0.0.1:9099` por padrão e registra seu endpoint no log. Nenhum código do usuário é necessário.

## Tutorial

### Etapa 1 — escreva um aplicativo Wails normal

O MCP não exige importações nem registro. Crie seu aplicativo como faria normalmente:

```go {title="main.go"}
package main

import (
    "embed"
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name: "My App",
        Assets: application.AssetOptions{
            Handler: application.BundledAssetFileServer(assets),
        },
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: "My App",
        Width: 1024, Height: 768,
    })

    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

### Etapa 2 — faça o build ou execute com a tag `mcp`

@tabs
[CLI do Wails (recomendado)]
Defina `WAILS_MCP=1` para que a CLI do Wails adicione a tag `mcp` para você:

```shell
# Development
WAILS_MCP=1 wails3 dev

# Production build
WAILS_MCP=1 wails3 build
```

[Go diretamente]
Passe a tag diretamente para `go run` ou `go build`:

```shell
go run -tags mcp .
go build -tags mcp -o myapp .
```

[Windows (PowerShell)]
```powershell
$env:WAILS_MCP = "1"
wails3 dev
# or
wails3 build
```

@end

Na inicialização, o aplicativo registra o endpoint MCP no log:

```
INFO MCP server started. Connect MCP clients using the streamable HTTP transport.
     url=http://127.0.0.1:9099/mcp
```

### Etapa 3 — conecte um cliente

O servidor usa o **transporte HTTP com streaming do MCP**. Conecte-se com qualquer cliente compatível com MCP.

@tabs
[Claude Code]
```shell
claude mcp add --transport http my-app http://127.0.0.1:9099/mcp
```

Em seguida, peça ao Claude para interagir com seu aplicativo:

```
Click the "Submit" button, then verify a success toast appears.
```

[VS Code (GitHub Copilot)]
Adicione a `.vscode/settings.json`:

```json
{
  "github.copilot.chat.mcp.enabled": true,
  "mcp": {
    "servers": {
      "my-wails-app": {
        "type": "http",
        "url": "http://127.0.0.1:9099/mcp"
      }
    }
  }
}
```

[Outros clientes]
Direcione qualquer cliente MCP compatível com transporte HTTP com streaming para:

```
http://127.0.0.1:9099/mcp
```

@end

### Etapa 4 — execute uma sessão de teste

Peça ao agente para testar seu aplicativo. Veja alguns exemplos de prompts:

```
Take a DOM snapshot of the main window.
```

```
Click the "Add item" button, type "Hello world" in the input field,
then press Enter and verify the item appears in the list.
```

```
Call the bound method main.GreetService.Greet with argument ["World"]
and return the result.
```

```
Wait for the event "save:complete" while clicking the Save button.
```

## Configuração

Toda a configuração é feita por meio de variáveis de ambiente — nenhuma alteração no código é necessária.

| Variável de ambiente | Padrão | Descrição |
| --- | --- | --- |
| `WAILS_MCP` | (não definida) | Defina como `1`, `true`, `on` ou `yes` para adicionar automaticamente a tag de build `mcp` ao usar a CLI do Wails. |
| `WAILS_MCP_HOST` | `127.0.0.1` | Interface à qual o servidor será associado. Associações a interfaces que não sejam de loopback exigem `WAILS_MCP_TOKEN`. |
| `WAILS_MCP_TOKEN` | não definida | Token bearer opcional em loopback; obrigatório em outros endereços de associação. Os clientes enviam `Authorization: Bearer <token>`. |
| `WAILS_MCP_PORT` | `9099` | Porta de escuta. Defina como `0` para usar uma porta livre atribuída aleatoriamente, que será exibida no log. |
| `WAILS_MCP_TIMEOUT` | `30000` | Tempo limite padrão para avaliação de JS, em **milissegundos**. |
| `WAILS_MCP_HIDE_CURSOR` | (não definida) | Defina como `1` ou `true` para desabilitar a sobreposição do cursor animado. |

Exemplo — porta personalizada e tempo limite de 60 segundos:

```shell
WAILS_MCP=1 WAILS_MCP_PORT=9200 WAILS_MCP_TIMEOUT=60000 wails3 dev
```

## Ferramentas disponíveis

| Ferramenta | Finalidade |
| --- | --- |
| `app_info` | Informações do aplicativo: plataforma, arquitetura, todas as janelas e endpoint MCP |
| `windows_list` | Lista todas as janelas com geometria e estado |
| `window_control` | Focar, redimensionar, mover, colocar em tela cheia, abrir as ferramentas de desenvolvimento, recarregar, definir a URL, … (22 ações) |
| `js_eval` | Avalia JavaScript em uma janela (corpo assíncrono, `return` para o valor) |
| `dom_html` | Obter o HTML da página ou de um elemento específico |
| `dom_query` | Localizar elementos por seletor CSS — tag, texto, limites, visibilidade |
| `screenshot_dom` | Instantâneo estrutural da página visível (baseado no DOM, sem pixels) |
| `mouse_move` | Animar o cursor até um ponto ou seletor CSS |
| `mouse_click` | Clicar com o cursor animado (botão esquerdo/direito/do meio, clique duplo, modificadores) |
| `mouse_drag` | Arrastar com o cursor animado (compatível com elementos de arrastar e soltar do HTML5) |
| `mouse_scroll` | Rolar em um ponto ou elemento |
| `keyboard_type` | Digitar texto caractere por caractere com eventos realistas |
| `keyboard_press` | Pressionar uma única tecla (Enter, Tab, Escape, ArrowDown, …) com modificadores opcionais |
| `call_bound_method` | Chamar um método vinculado de um serviço Go, por exemplo, `main.GreetService.Greet` |
| `emit_event` | Emitir um evento da aplicação Wails |
| `wait_for_event` | Aguardar um evento da aplicação Wails e retornar seus dados |

### Suporte a várias janelas

Todas as ferramentas que atuam em uma janela aceitam um argumento opcional `window` contendo o **nome** da janela (definido por meio de `WebviewWindowOptions.Name`). Quando omitido, a ferramenta atua na janela que está em foco ou na primeira janela, se nenhuma estiver em foco.

```
List all windows, then click the "New" button in the window named "editor".
```

### Seleção de elementos

As ferramentas de mouse e teclado aceitam uma destas opções:

- **Seletor CSS** — `selector: "#submit-btn"` (o elemento é rolado automaticamente para a área visível)
- **Coordenadas** — `x: 400, y: 300` (pixels CSS relativos à viewport)

Para operações de arrastar, use os prefixos `from_` e `to_`:

```
Drag from selector: ".card" to selector: ".dropzone"
```

## Segurança

@note{type="caution"}
O servidor MCP oferece controle programático total sobre a sua aplicação. Qualquer pessoa que possa acessar as ferramentas dele pode ler o DOM, avaliar JavaScript, clicar em botões e chamar métodos Go.

@end

- Por padrão, o servidor escuta em `127.0.0.1`. As origens do navegador devem ser origens de loopback HTTP(S); origens opacas (`null`), malformadas e externas são rejeitadas.
- Clientes MCP nativos sem cabeçalhos continuam compatíveis. Sem `WAILS_MCP_TOKEN`, os processos locais e as origens locais permitidas do navegador são considerados confiáveis; as verificações de origem não constituem autenticação. Defina um token de alta entropia para exigir autenticação por bearer token em todas as chamadas `/mcp`. Configure o mesmo token no cabeçalho `Authorization: Bearer <token>` do cliente. As solicitações de preflight não exigem o token.
- O callback `/eval-result` usa IDs imprevisíveis específicos de cada avaliação em vez do bearer token do cliente, de modo que a entrega de resultados pela webview continue compatível.
- As builds de produção **não** devem incluir a tag `mcp`. A CLI do Wails só a adiciona quando `WAILS_MCP=1` é definido explicitamente, e a `wails3 build` padrão não contém nenhum código do servidor.
- Se precisar expor o servidor em uma interface que não seja de loopback (por exemplo, para testes em LAN), defina `WAILS_MCP_HOST=0.0.0.0` e um `WAILS_MCP_TOKEN` de alta entropia; a inicialização falhará sem um token. Use um túnel criptografado ou um proxy que encerre TLS em redes não confiáveis, pois o listener integrado usa HTTP.

## Aplicação de exemplo

Uma aplicação playground completa que demonstra todas as ferramentas está disponível em [`v3/examples/mcp`](https://github.com/wailsapp/wails/tree/releases/v3-beta/v3/examples/mcp). Ela inclui:

- Contador com botões para incrementar e redefinir
- Campo de nome com os métodos vinculados Greet, Add e Shout
- Origem e destino de arrastar e soltar do HTML5
- Lista rolável (50 itens)
- Registro de eventos

Execute-a com:

```shell
cd v3/examples/mcp
go run -tags mcp .
```

Em seguida, conecte o Claude Code ou qualquer cliente MCP a `http://127.0.0.1:9099/mcp` e peça que ele interaja com a interface.

## Feedback

O servidor MCP integrado é um experimento, e seu feedback decidirá o rumo dele. Se você testá-lo, queremos saber qual cliente e quais ferramentas usou, o que esperava e o que realmente aconteceu, e se permitir que um agente controlasse sua aplicação foi útil — os relatos mais úteis informam exatamente o que foi executado. Conte-nos na [discussão de feedback sobre o servidor MCP](https://github.com/wailsapp/wails/discussions/5692).
