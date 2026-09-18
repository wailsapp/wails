---
title: "Como o Wails funciona"
description: "Entenda a arquitetura do Wails e como ela alcança desempenho nativo"
slug: "concepts/architecture"
sourcePath: "concepts/architecture.md"
---

Wails é um framework para criar aplicativos desktop usando **Go no backend** e **tecnologias web no frontend**. Porém, ao contrário do Electron, o Wails não inclui um navegador no pacote — ele usa a **WebView nativa do sistema operacional**.

```d2
direction: left

Wails App: {
  shape: sequence_diagram
  label: Aplicativo Wails

  frontend: Frontend
  backend: Backend em Go
  os: Sistema operacional

  Initialisation: Inicialização {
    shape: sequence_diagram
    backend."Serves Static Web App": Serve o aplicativo web estático
    backend -> frontend: HTML / JS / CSS
    frontend."Render Site via OS-native WebView": Renderiza o site pela WebView nativa do sistema operacional
  }
  Regular Communication: Comunicação regular {
    shape: sequence_diagram
    frontend."Make API-style call": Faz uma chamada no estilo de API
    frontend -> backend.a: JSON
    backend.a."Service processes request": O serviço processa a solicitação
    backend.a -> os: Chama APIs do sistema
    backend.a."Generate Response": Gera a resposta
    backend.a -> frontend: JSON
    frontend."Process response": Processa a resposta
  }
  backend.a.label: a
}
```

**Principais diferenças em relação ao Electron:**

| Aspecto | Wails | Electron |
| --- | --- | --- |
| **Navegador** | WebView fornecida pelo sistema operacional | Chromium incluído no pacote (~100 MB) |
| **Backend** | Go (compilado) | Node.js (interpretado) |
| **Comunicação** | Ponte em memória | IPC (entre processos) |
| **Tamanho do pacote** | ~15 MB | ~150 MB |
| **Memória** | ~10 MB | ~100 MB+ |
| **Inicialização** | &lt;0.5s | 2-3 s |

## Componentes principais

### 1. WebView nativa

O Wails usa o mecanismo de renderização web integrado ao sistema operacional:

@tabs{sync-key="platform"}
[Windows]
**WebView2** (Microsoft Edge WebView2)

- Baseada no Chromium (o mesmo mecanismo do navegador Edge)
- Pré-instalada no Windows 10/11
- Atualizações automáticas pelo Windows Update
- Compatibilidade total com os padrões web modernos

[macOS]
**WebKit** (mecanismo de renderização do Safari)

- Integrado ao macOS
- O mesmo mecanismo do navegador Safari
- Excelente desempenho e duração da bateria
- Compatibilidade total com os padrões web modernos

[Linux]
**WebKitGTK** (port do WebKit para GTK)

- Instalado pelo gerenciador de pacotes
- O mesmo mecanismo do GNOME Web (Epiphany)
- Boa compatibilidade com os padrões
- Leve e eficiente

@end

**Por que isso é importante:**

- **Nenhum navegador incluído no pacote** → Aplicativo menor
- **Nativo do sistema operacional** → Melhor integração e desempenho
- **Atualizações automáticas** → Correções de segurança fornecidas pelas atualizações do sistema operacional
- **Renderização familiar** → A mesma do navegador do sistema

### 2. A ponte do Wails

A ponte é o núcleo do Wails — ela permite a **comunicação direta** entre Go e JavaScript.

```d2
direction: down

Frontend: Frontend (JavaScript) {
  shape: rectangle
  style.fill: "#8B5CF6"
}

Bridge: Ponte do Wails {
  Encoder: Codificador JSON {
    shape: rectangle
  }

  Router: Roteador de métodos {
    shape: diamond
    style.fill: "#10B981"
  }

  Decoder: Decodificador JSON {
    shape: rectangle
  }
}

Backend: Backend (Go) {
  Services: Serviços registrados {
    shape: rectangle
    style.fill: "#00ADD8"
  }
}

Frontend -> Bridge.Encoder: "1. Chama o método Go\nGreet('Alice')"
Bridge.Encoder -> Bridge.Router: "2. Codifica como JSON\n{method: 'Greet', args: ['Alice']}"
Bridge.Router -> Backend.Services: "3. Encaminha ao serviço\nGreetService.Greet('Alice')"
Backend.Services -> Bridge.Decoder: "4. Retorna o resultado\n'Hello, Alice!'"
Bridge.Decoder -> Frontend: "5. Decodifica para JS\nA Promise é resolvida"
```

**Como funciona:**

1. **O frontend chama um método Go** (por meio de um binding gerado automaticamente)
2. **A ponte codifica a chamada** em JSON (nome do método + argumentos)
3. **O roteador localiza o método Go** nos serviços registrados
4. **O método Go é executado** e retorna um valor
5. **A ponte decodifica o resultado** e o envia de volta ao frontend
6. **A Promise é resolvida** no JavaScript com o resultado

**Características de desempenho:**

- **Em memória**: sem sobrecarga de rede e sem HTTP
- **Sem cópias** quando possível (para grandes volumes de dados)
- **Assíncrona por padrão**: não bloqueia nenhum dos lados
- **Com segurança de tipos**: definições TypeScript geradas automaticamente

### 3. Sistema de serviços

Os serviços são a maneira recomendada de expor funcionalidades Go ao frontend.

```go
// Define a service (just a regular Go struct)
type GreetService struct {
    prefix string
}

// Methods with exported names are automatically available
func (g *GreetService) Greet(name string) string {
    return g.prefix + name + "!"
}

func (g *GreetService) GetTime() time.Time {
    return time.Now()
}

// Register the service
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&GreetService{prefix: "Hello, "}),
    },
})
```

**Descoberta de serviços:**

- O Wails **inspeciona sua struct** durante a inicialização
- **Os métodos exportados** podem ser chamados pelo frontend
- **As informações de tipos** são extraídas para os bindings TypeScript
- **O tratamento de erros** é automático (erros Go → exceções JS)

**Binding TypeScript gerado:**

```typescript
// Auto-generated in frontend/bindings/GreetService.ts
export function Greet(name: string): Promise<string>
export function GetTime(): Promise<Date>
```

**Por que usar serviços?**

- **Tipagem segura**: suporte completo ao TypeScript
- **Descoberta automática**: não é necessário registrar métodos manualmente
- **Organização**: agrupe funcionalidades relacionadas
- **Testabilidade**: os serviços são apenas structs Go

[Saiba mais sobre serviços →](/features/bindings/services/)

### 4. Sistema de eventos

Os eventos permitem a **comunicação por publicação/assinatura** entre componentes.

```d2
direction: left

Wails Event System: Sistema de eventos do Wails {
  shape: sequence_diagram

  window1: Janela 1
  window2: Janela 2
  backend: Backend em Go

  Event Driver: Gerenciador de eventos {
    shape: sequence_diagram
    window1."Subscribe to 'data-updated' events": "Assina os eventos 'data-updated'"
    window2."Subscribe to 'data-updated' events": "Assina os eventos 'data-updated'"
    backend.a."App Emit('data-updated', data)": "O aplicativo emite Emit('data-updated', data)"
    backend.a -> window1.a: Barramento de eventos JSON
    backend.a -> window2: Barramento de eventos JSON
    window1.a."Subscriber processes On('data-updated', handler)": "O assinante processa On('data-updated', handler)"
    window2."Subscriber processes On('data-updated', handler)": "O assinante processa On('data-updated', handler)"
  }
  backend.a.label: a
  window1.a.label: a
}
```

**Casos de uso:**

- **Comunicação entre janelas**: uma janela notifica as demais
- **Tarefas em segundo plano**: o serviço Go notifica a interface sobre o progresso
- **Sincronização de estado**: mantenha várias janelas sincronizadas
- **Baixo acoplamento**: os componentes não precisam de referências diretas

**Exemplo:**

```go
// Go: Emit an event
app.Event.Emit("user-logged-in", user)
```

```javascript
// JavaScript: Listen for event
import { Events } from '@wailsio/runtime'

Events.On('user-logged-in', (user) => {
    console.log('User logged in:', user)
})
```

[Saiba mais sobre eventos →](/features/events/system/)

## Ciclo de vida da aplicação

Compreender o ciclo de vida ajuda você a saber quando inicializar e liberar recursos.

```d2
direction: down

Start: Início do aplicativo {
  shape: oval
  style.fill: "#10B981"
}

Init: Inicialização {
  Create: Cria o aplicativo {
    shape: rectangle
  }

  Register: Registra os serviços {
    shape: rectangle
  }

  Setup: Configura janelas e menus {
    shape: rectangle
  }
}

Run: Loop de eventos {
  Events: Processa eventos {
    shape: rectangle
  }

  Messages: Trata mensagens {
    shape: rectangle
  }

  Render: Atualiza a interface {
    shape: rectangle
  }
}

Shutdown: Encerramento {
  Cleanup: Libera recursos {
    shape: rectangle
  }

  Save: Salva o estado {
    shape: rectangle
  }
}

End: Fim do aplicativo {
  shape: oval
  style.fill: "#EF4444"
}

Start -> Init.Create
Init.Create -> Init.Register
Init.Register -> Init.Setup
Init.Setup -> Run.Events
Run.Events -> Run.Messages
Run.Messages -> Run.Render
Run.Render -> Run.Events: Loop
Run.Events -> Shutdown.Cleanup: Sinal de encerramento
Shutdown.Cleanup -> Shutdown.Save
Shutdown.Save -> End
```

**Hooks do ciclo de vida:**

```go
app := application.New(application.Options{
    Name: "My App",

    // Cleanly intercept quit requests (e.g. unsaved changes).
    ShouldQuit: func() bool { return true },

    // Called when the app is confirmed to be quitting — save state, close connections, etc.
    OnShutdown: func() {},
})
```

Não existe um campo `OnStartup` em `application.Options`. O trabalho de inicialização deve ficar no `ServiceStartup(ctx, options)` de um serviço, em um callback registrado por meio de `app.Event.OnApplicationEvent(events.Common.ApplicationStarted, ...)` ou simplesmente antes de `app.Run()`.

[Saiba mais sobre o ciclo de vida →](/concepts/lifecycle/)

## Processo de build

Entenda como o Wails faz o build da sua aplicação:

```d2
direction: down

Source: Código-fonte {
  Go: "Código Go\n(main.go, serviços)" {
    shape: rectangle
    style.fill: "#00ADD8"
  }

  Frontend: "Código do frontend\n(HTML/CSS/JS)" {
    shape: rectangle
    style.fill: "#8B5CF6"
  }
}

Build: Processo de build {
  AnalyseGo: Analisa o código Go {
    shape: rectangle
  }

  GenerateBindings: Gera os bindings {
    shape: rectangle
  }

  BuildFrontend: Faz o build do frontend {
    shape: rectangle
  }

  CompileGo: Compila o código Go {
    shape: rectangle
  }

  Embed: Incorpora os recursos {
    shape: rectangle
  }
}

Output: Saída {
  Binary: "Binário nativo\n(myapp.exe/.app)" {
    shape: rectangle
    style.fill: "#10B981"
  }
}

Source.Go -> Build.AnalyseGo
Build.AnalyseGo -> Build.GenerateBindings: Extrai os tipos
Build.GenerateBindings -> Source.Frontend: Bindings TypeScript
Source.Frontend -> Build.BuildFrontend: Compila (Vite/webpack)
Build.BuildFrontend -> Build.Embed: Recursos empacotados
Source.Go -> Build.CompileGo
Build.CompileGo -> Build.Embed
Build.Embed -> Output.Binary
```

**Etapas do build:**

1. **Analisar o código Go**
  - Examinar os serviços em busca de métodos exportados
  - Extrair os tipos dos parâmetros e dos valores de retorno
  - Gerar assinaturas de métodos


2. **Gerar bindings TypeScript**
  - Criar arquivos `.ts` para cada serviço
  - Incluir definições completas de tipos
  - Adicionar comentários JSDoc


3. **Fazer o build do frontend**
  - Executar seu empacotador (Vite, webpack etc.)
  - Minificar e otimizar
  - Gerar a saída em `frontend/dist/`


4. **Compilar o código Go**
  - Compilar com otimizações (`-ldflags="-s -w"`)
  - Incluir metadados do build
  - Compilar especificamente para cada plataforma


5. **Incorporar os recursos**
  - Incorporar os arquivos do frontend ao binário Go
  - Compactar os recursos
  - Criar um único executável


**Resultado:** um único executável nativo com tudo incorporado.

[Saiba mais sobre o processo de build →](/guides/build/building/)

## Desenvolvimento vs. produção

O Wails se comporta de maneira diferente nos ambientes de desenvolvimento e produção:

@tabs{sync-key="mode"}
[Desenvolvimento (wails3 dev)]
**Características:**

- **Recarregamento automático**: as alterações no frontend são recarregadas imediatamente
- **Mapas de código-fonte**: depuração com o código-fonte original
- **DevTools**: as ferramentas de desenvolvimento do navegador ficam disponíveis
- **Logs**: logs detalhados ficam habilitados
- **Frontend externo**: servido pelo servidor de desenvolvimento (Vite)

**Como funciona:**

```d2
direction: right

WailsApp: Aplicativo Wails {
  shape: rectangle
  style.fill: "#00ADD8"
}

DevServer: "Servidor de desenvolvimento do Vite\n(localhost:5173)" {
  shape: rectangle
  style.fill: "#8B5CF6"
}

WebView: WebView {
  shape: rectangle
  style.fill: "#6B7280"
}

WailsApp -> DevServer: Encaminha solicitações por proxy
DevServer -> WebView: Serve com HMR
WebView -> WailsApp: Chama métodos Go
```

**Benefícios:**

- Retorno imediato sobre as alterações
- Recursos completos de depuração
- Iterações mais rápidas

[Produção (wails3 build)]
**Características:**

- **Recursos incorporados**: o frontend é incluído no binário durante o build
- **Otimizado**: minificado e compactado
- **Sem DevTools**: desabilitadas por padrão
- **Logs mínimos**: somente erros
- **Arquivo único**: tudo em um único executável

**Como funciona:**

```d2
direction: right

Binary: "Binário único\n(myapp.exe)" {
  GoCode: Código Go compilado {
    shape: rectangle
    style.fill: "#00ADD8"
  }

  Assets: "Recursos incorporados\n(HTML/CSS/JS)" {
    shape: rectangle
    style.fill: "#8B5CF6"
  }
}

WebView: WebView {
  shape: rectangle
  style.fill: "#6B7280"
}

Binary.Assets -> WebView: Serve a partir da memória
WebView -> Binary.GoCode: Chama métodos Go
```

**Benefícios:**

- Distribuição em um único arquivo
- Tamanho menor (minificado)
- Melhor desempenho
- Sem dependências externas

@end

## Modelo de memória

Entender o uso da memória ajuda você a criar aplicações eficientes.

**Regiões da memória:**

1. **Heap do Go**
  - Seus serviços e o estado da aplicação
  - Gerenciado pelo coletor de lixo do Go
  - Normalmente 5-10MB para aplicações simples


2. **Memória do WebView**
  - DOM, heap do JavaScript e CSS
  - Gerenciada pelo mecanismo do WebView
  - Normalmente 10-20MB para aplicações simples


3. **Memória da ponte**
  - Buffers de mensagens para comunicação
  - Sobrecarga mínima (<1MB)
  - Cópia zero para grandes volumes de dados, quando possível


**Dicas de otimização:**

- **Evite transferir grandes volumes de dados**: passe IDs e busque os detalhes sob demanda
- **Use eventos para atualizações**: não faça polling no frontend
- **Transmita arquivos grandes por streaming**: não os carregue por completo na memória
- **Remova os listeners**: remova os listeners de eventos quando não forem mais necessários

[Saiba mais sobre desempenho →](/guides/performance/)

## Modelo de segurança

O Wails oferece uma arquitetura segura por padrão:

```d2
direction: down

Frontend: Frontend (não confiável) {
  shape: rectangle
  style.fill: "#EF4444"
}

Bridge: Ponte do Wails (validação) {
  shape: diamond
  style.fill: "#F59E0B"
}

Backend: Backend (confiável) {
  shape: rectangle
  style.fill: "#10B981"
}

Frontend -> Bridge: Chama o método
Bridge -> Bridge: "Validar:\n- O método existe?\n- Os tipos estão corretos?\n- O acesso é permitido?"
Bridge -> Backend: Executar se for válido
Backend -> Bridge: Retornar o resultado
Bridge -> Frontend: Enviar a resposta
```

**Recursos de segurança:**

1. **Lista de métodos permitidos**
  - Somente métodos exportados podem ser chamados
  - Métodos privados são inacessíveis
  - É necessário registrar explicitamente os serviços


2. **Validação de tipos**
  - Os argumentos são verificados em relação aos tipos do Go
  - Tipos inválidos são rejeitados
  - Impede ataques de injeção


3. **Sem eval()**
  - O frontend não pode executar código Go arbitrário
  - Somente métodos predefinidos podem ser chamados
  - Sem execução dinâmica de código


4. **Isolamento de contexto**
  - Cada janela tem seu próprio contexto
  - Os serviços podem verificar o contexto do chamador
  - É possível definir permissões por janela


**Práticas recomendadas:**

- **Valide a entrada do usuário** no Go (não confie no frontend)
- **Use o contexto** para autenticação e autorização
- **Sanitize os caminhos de arquivos** antes das operações com arquivos
- **Limite a frequência** de operações custosas

[Saiba mais sobre segurança →](/guides/security/)

## Próximos passos

**Ciclo de vida da aplicação** - Entenda a inicialização, o encerramento e os hooks do ciclo de vida [Saiba mais →](/concepts/lifecycle/)

**Ponte entre Go e frontend** - Explore em detalhes como a ponte funciona [Saiba mais →](/concepts/bridge/)

**Sistema de build** - Entenda como o Wails compila sua aplicação [Saiba mais →](/concepts/build-system/)

**Comece a desenvolver** - Aplique em um tutorial o que você aprendeu [Tutoriais →](/tutorials/03-notes-vanilla/)

---

**Dúvidas sobre a arquitetura?** Pergunte no [Discord](https://discord.gg/JDdSxwjhGf) ou consulte a [referência da API](/reference/overview/).
