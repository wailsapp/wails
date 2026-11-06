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

Para automatizar o ciclo de vida do projeto, use o servidor separado da CLI `wails3 mcp`. Ele permite que os agentes inspecionem e inicializem projetos, executem diagnósticos, iniciem compilações e tarefas de desenvolvimento, gerem bindings, executem tarefas nomeadas do Taskfile e obtenham uma quantidade limitada de saída das tarefas. Por padrão, o servidor da CLI fica restrito ao diretório atual e não permite a execução arbitrária de comandos do shell. Consulte a [documentação do MCP da CLI](/guides/cli/#mcp) para obter detalhes sobre transporte, autenticação e ferramentas.

Quando esse recurso está habilitado, um agente conectado ao seu aplicativo pode:

- **Listar e controlar janelas** — tamanho, posição, foco, tela cheia, ferramentas de desenvolvimento, recarregamento, …
- **Inspecionar o DOM** — consultar elementos, obter HTML e capturar um snapshot estrutural
- **Avaliar JavaScript** — executar código arbitrário dentro de qualquer janela e obter o resultado
- **Simular entradas do usuário** — movimentos, cliques, ações de arrastar e rolagens do mouse, renderizados com um **cursor animado na tela** para que você possa acompanhar o trabalho do agente
- **Digitar e pressionar teclas** — eventos realistas para cada caractere que funcionam com entradas controladas do React
- **Chamar métodos Go vinculados** e emitir ou aguardar eventos do aplicativo

## Como funciona

O servidor MCP só é compilado no seu aplicativo quando a tag de compilação **`mcp`** está presente. Sem essa tag, o código do servidor fica completamente ausente do binário — sem sobrecarga em tempo de execução, portas abertas ou superfície de ataque.

Quando a tag está presente, o servidor é iniciado automaticamente dentro de `App.Run()`, vincula-se por padrão a `127.0.0.1:9099` e registra seu endpoint no log. Nenhum código do usuário é necessário.

## Tutorial

### Etapa 1 — crie um aplicativo Wails normal

O MCP não exige importações nem registros. Crie seu aplicativo como de costume:

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

### Etapa 2 — compile ou execute com a tag `mcp`

@tabs
[CLI do Wails (recomendado)]
Defina `WAILS_MCP=1`, e a CLI do Wails adicionará a tag `mcp` para você:

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

O servidor usa o **transporte HTTP com streaming do MCP**. Conecte-o a qualquer cliente compatível com MCP.

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
Direcione qualquer cliente MCP compatível com o transporte HTTP com streaming para:

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
| `WAILS_MCP` | (não definida) | Defina como `1`, `true`, `on` ou `yes` para adicionar automaticamente a tag de compilação `mcp` ao usar a CLI do Wails. |
| `WAILS_MCP_HOST` | `127.0.0.1` | Interface à qual o servidor será vinculado. Vínculos a interfaces que não sejam de loopback exigem `WAILS_MCP_TOKEN`. |
| `WAILS_MCP_TOKEN` | não definida | Token bearer opcional em loopback; obrigatório em outros endereços de vínculo. Os clientes enviam `Authorization: Bearer <token>`. |
| `WAILS_MCP_PORT` | `9099` | Porta de escuta. Defina como `0` para usar uma porta livre atribuída aleatoriamente, que será exibida no log. |
| `WAILS_MCP_TIMEOUT` | `30000` | Tempo limite padrão da avaliação de JS em **milissegundos**. |
| `WAILS_MCP_HIDE_CURSOR` | (não definida) | Defina como `1` ou `true` para desativar a sobreposição do cursor animado. |

Exemplo — porta personalizada e tempo limite de 60 segundos:

```shell
WAILS_MCP=1 WAILS_MCP_PORT=9200 WAILS_MCP_TIMEOUT=60000 wails3 dev
```

## Ferramentas disponíveis

| Ferramenta | Finalidade |
| --- | --- |
| `app_info` | Informações do aplicativo: plataforma, arquitetura, todas as janelas e endpoint MCP |
| `windows_list` | Lista todas as janelas com geometria e estado |
| `window_control` | Foco, redimensionamento, movimentação, tela cheia, ferramentas de desenvolvimento, recarregamento, definição de URL, … (22 ações) |
| `js_eval` | Avalia JavaScript em uma janela (corpo assíncrono, `return` para o valor) |
| `dom_html` | Obtém o HTML da página ou de um elemento específico |
| `dom_query` | Localiza elementos por seletor CSS — tag, texto, limites e visibilidade |
| `screenshot_dom` | Snapshot estrutural da página visível (baseado no DOM, sem pixels) |
| `mouse_move` | Anima o cursor até um ponto ou seletor CSS |
| `mouse_click` | Clica com o cursor animado (botão esquerdo, direito ou do meio, clique duplo e modificadores) |
| `mouse_drag` | Arrasta com o cursor animado (compatível com elementos de arrastar e soltar do HTML5) |
| `mouse_scroll` | Rola em um ponto ou elemento |
| `keyboard_type` | Digita texto caractere por caractere com eventos realistas |
| `keyboard_press` | Pressionar uma única tecla (Enter, Tab, Escape, ArrowDown, …), com modificadores opcionais |
| `call_bound_method` | Chamar um método de serviço Go vinculado, por exemplo, `main.GreetService.Greet` |
| `emit_event` | Emitir um evento da aplicação Wails |
| `wait_for_event` | Aguardar um evento da aplicação Wails e retornar seus dados |

### Suporte a várias janelas

Todas as ferramentas que atuam em uma janela aceitam um argumento opcional `window` contendo o **nome** da janela (definido por meio de `WebviewWindowOptions.Name`). Quando ele é omitido, a ferramenta atua na janela que está com o foco ou na primeira janela, se nenhuma estiver com o foco.

```
List all windows, then click the "New" button in the window named "editor".
```

### Seleção de elementos

As ferramentas de mouse e teclado aceitam uma destas opções:

- **Seletor CSS** — `selector: "#submit-btn"` (o elemento é rolado automaticamente para a área visível)
- **Coordenadas** — `x: 400, y: 300` (pixels CSS relativos à área de visualização)

Para operações de arrastar, use os prefixos `from_` e `to_`:

```
Drag from selector: ".card" to selector: ".dropzone"
```

## Segurança

@note{type="caution"}
O servidor MCP concede controle programático total sobre sua aplicação. Qualquer pessoa que tenha acesso às ferramentas pode ler o DOM, avaliar JavaScript, clicar em botões e chamar métodos Go.

@end

- Por padrão, o servidor se vincula a `127.0.0.1`. As origens do navegador devem ser origens de loopback HTTP(S); origens opacas (`null`), malformadas e externas são rejeitadas.
- Clientes MCP nativos sem cabeçalhos continuam sendo compatíveis. Sem `WAILS_MCP_TOKEN`, processos locais e origens locais permitidas do navegador são considerados confiáveis; as verificações de origem não constituem autenticação. Defina um token de alta entropia para exigir autenticação por bearer token em todas as chamadas `/mcp`. Configure o mesmo token no cabeçalho `Authorization: Bearer <token>` do cliente. As solicitações de preflight não exigem o token.
- O callback `/eval-result` usa IDs imprevisíveis para cada avaliação, em vez do bearer token do cliente, mantendo compatível a entrega de resultados da webview.
- As compilações de produção **não** devem incluir a tag `mcp`. A CLI do Wails só a adiciona quando `WAILS_MCP=1` é definido explicitamente, e o `wails3 build` padrão não contém nenhum código do servidor.
- Se precisar expor o servidor em uma interface que não seja de loopback (por exemplo, para testes em uma LAN), defina `WAILS_MCP_HOST=0.0.0.0` e um `WAILS_MCP_TOKEN` de alta entropia; a inicialização falhará sem um token. Em redes não confiáveis, use um túnel criptografado ou um proxy com terminação TLS, pois o listener integrado usa HTTP.

## Aplicação de exemplo

Uma aplicação de experimentação completa que demonstra todas as ferramentas está disponível em [`v3/examples/mcp`](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/mcp). Ela inclui:

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

Em seguida, conecte o Claude Code ou qualquer cliente MCP a `http://127.0.0.1:9099/mcp` e peça que ele teste a interface.

## Feedback

O servidor MCP integrado é um experimento, e seu feedback definirá os próximos passos. Se você testá-lo, gostaríamos de saber qual cliente e quais ferramentas usou, o que esperava e o que realmente aconteceu, além de saber se permitir que um agente controlasse sua aplicação foi útil — os relatos mais úteis informam exatamente o que foi executado. Conte-nos na [discussão de feedback sobre o servidor MCP](https://github.com/wailsapp/wails/discussions/5692).
