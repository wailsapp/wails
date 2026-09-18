---
title: "Ponte entre Go e frontend"
description: "Análise detalhada de como o Wails permite a comunicação direta entre Go e JavaScript"
slug: "concepts/bridge"
sourcePath: "concepts/bridge.md"
---

## Comunicação direta entre Go e JavaScript

O Wails fornece uma ponte **direta e em memória** entre Go e JavaScript, permitindo uma comunicação integrada sem a sobrecarga do HTTP, limites entre processos ou gargalos de serialização.

## Visão geral

```d2
direction: right

Frontend: Frontend (JavaScript) {
  UI: React/Vue/JavaScript puro {
    shape: rectangle
    style.fill: "#8B5CF6"
  }

  Bindings: Bindings gerados automaticamente {
    shape: rectangle
    style.fill: "#A78BFA"
  }
}

Bridge: Ponte do Wails {
  Encoder: Codificador JSON {
    shape: rectangle
    style.fill: "#10B981"
  }

  Router: Roteador de métodos {
    shape: diamond
    style.fill: "#10B981"
  }

  Decoder: Decodificador JSON {
    shape: rectangle
    style.fill: "#10B981"
  }

  TypeGen: Gerador de tipos {
    shape: rectangle
    style.fill: "#10B981"
  }
}

Backend: Backend (Go) {
  Services: Seus serviços {
    shape: rectangle
    style.fill: "#00ADD8"
  }

  Registry: Registro de serviços {
    shape: rectangle
    style.fill: "#00ADD8"
  }
}

Frontend.UI -> Frontend.Bindings: "import { Method }"
Frontend.Bindings -> Bridge.Encoder: "Chamar Method('arg')"
Bridge.Encoder -> Bridge.Router: Codificar como JSON
Bridge.Router -> Backend.Registry: Localizar serviço
Backend.Registry -> Backend.Services: Invocar método
Backend.Services -> Bridge.Decoder: Retornar resultado
Bridge.Decoder -> Frontend.Bindings: Decodificar para JS
Frontend.Bindings -> Frontend.UI: Promise é resolvida
Bridge.TypeGen -> Frontend.Bindings: Gerar tipos
```

**Ponto principal:** sem HTTP, IPC ou limites entre processos. Apenas **chamadas diretas de funções** com **segurança de tipos**.

## Como funciona: passo a passo

### 1. Registro de serviços (inicialização)

Quando seu aplicativo é iniciado, o Wails examina seus serviços:

```go
type GreetService struct {
    prefix string
}

func (g *GreetService) Greet(name string) string {
    return g.prefix + name + "!"
}

func (g *GreetService) Add(a, b int) int {
    return a + b
}

// Register service
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&GreetService{prefix: "Hello, "}),
    },
})
```

**O que o Wails faz:**

1. **Examina a struct** em busca de métodos exportados
2. **Extrai informações de tipos** (parâmetros e tipos de retorno)
3. **Cria um registro** que associa nomes de métodos a funções
4. **Gera bindings TypeScript** com definições completas de tipos

### 2. Geração de bindings (durante a compilação)

O Wails gera bindings TypeScript automaticamente:

```typescript
// Auto-generated: frontend/bindings/GreetService.ts
export function Greet(name: string): Promise<string>
export function Add(a: number, b: number): Promise<number>
```

**Mapeamento de tipos:**

| Tipo Go | Tipo TypeScript |
| --- | --- |
| `string` | `string` |
| `int`, `int32`, `int64` | `number` |
| `float32`, `float64` | `number` |
| `bool` | `boolean` |
| `[]T` | `T[]` |
| `map[string]T` | `Record<string, T>` |
| `struct` | `interface` |
| `time.Time` | `Date` |
| `error` | Exceção (lançada) |

### 3. Chamada do frontend (em tempo de execução)

O desenvolvedor chama o método Go no JavaScript:

```javascript
import { Greet, Add } from './bindings/GreetService'

// Call Go from JavaScript
const greeting = await Greet("World")
console.log(greeting)  // "Hello, World!"

const sum = await Add(5, 3)
console.log(sum)  // 8
```

**O que acontece:**

1. **Função de binding chamada** — `Greet("World")`
2. **Mensagem criada** — `{ service: "GreetService", method: "Greet", args: ["World"] }`
3. **Enviada para a ponte** — pela ponte JavaScript da WebView
4. **Promise retornada** — aguarda a resposta

### 4. Processamento da ponte (tempo de execução)

A ponte recebe a mensagem e a processa:

```d2
direction: down

Receive: Receber mensagem {
  shape: rectangle
  style.fill: "#10B981"
}

Parse: Analisar JSON {
  shape: rectangle
}

Validate: Validar {
  Check: O serviço existe? {
    shape: diamond
  }

  CheckMethod: O método existe? {
    shape: diamond
  }

  CheckTypes: Os tipos estão corretos? {
    shape: diamond
  }
}

Invoke: Invocar método Go {
  shape: rectangle
  style.fill: "#00ADD8"
}

Encode: Codificar resultado {
  shape: rectangle
}

Send: Enviar resposta {
  shape: rectangle
  style.fill: "#10B981"
}

Error: Enviar erro {
  shape: rectangle
  style.fill: "#EF4444"
}

Receive -> Parse
Parse -> Validate.Check
Validate.Check -> Validate.CheckMethod: Sim
Validate.Check -> Error: Não
Validate.CheckMethod -> Validate.CheckTypes: Sim
Validate.CheckMethod -> Error: Não
Validate.CheckTypes -> Invoke: Sim
Validate.CheckTypes -> Error: Não
Invoke -> Encode: Sucesso
Invoke -> Error: Erro
Encode -> Send
```

**Segurança:** somente serviços registrados e métodos exportados podem ser chamados.

### 5. Execução em Go (tempo de execução)

O método Go é executado:

```go
func (g *GreetService) Greet(name string) string {
    // This runs in Go
    return g.prefix + name + "!"
}
```

**Contexto de execução:**

- É executado em uma **goroutine** (sem bloqueio)
- Tem acesso a **todos os recursos do Go** (sistema de arquivos, rede e bancos de dados)
- Pode chamar livremente **outro código Go**
- Retorna um resultado ou erro

### 6. Resposta (tempo de execução)

O resultado é enviado de volta ao JavaScript:

```javascript
// Promise resolves with result
const greeting = await Greet("World")
// greeting = "Hello, World!"
```

**Tratamento de erros:**

```go
func (g *GreetService) Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}
```

```javascript
try {
    const result = await Divide(10, 0)
} catch (error) {
    console.error("Go error:", error)  // "division by zero"
}
```

## Características de desempenho

### Velocidade

**Sobrecarga típica por chamada:** &lt;1 ms

```
Frontend Call → Bridge → Go Execution → Bridge → Frontend Response
     ↓            ↓           ↓            ↓            ↓
   &lt;0.1ms      &lt;0.1ms      [varies]     &lt;0.1ms      &lt;0.1ms
```

**Comparação com as alternativas:**

- **HTTP/REST:** 5-50 ms (pilha de rede, serialização)
- **IPC:** 1-10 ms (limites entre processos, marshaling)
- **Ponte do Wails:** &lt;1 ms (em memória, chamada direta)

### Memória

**Sobrecarga por chamada:** ~1 KB (buffer de mensagens)

**Otimização sem cópia:** dados grandes (>1 MB) usam memória compartilhada sempre que possível.

### Concorrência

**As chamadas são concorrentes:**

- Cada chamada é executada em sua própria goroutine
- Várias chamadas podem ser executadas simultaneamente
- Não há bloqueio entre chamadas

```javascript
// These run concurrently
const [result1, result2, result3] = await Promise.all([
    SlowOperation1(),
    SlowOperation2(),
    SlowOperation3(),
])
```

## Sistema de tipos

### Tipos compatíveis

#### Tipos primitivos

```go
// Go
func Example(
    s string,
    i int,
    f float64,
    b bool,
) (string, int, float64, bool) {
    return s, i, f, b
}
```

```typescript
// TypeScript (auto-generated)
function Example(
    s: string,
    i: number,
    f: number,
    b: boolean,
): Promise<[string, number, number, boolean]>
```

#### Slices e arrays

```go
// Go
func Sum(numbers []int) int {
    total := 0
    for _, n := range numbers {
        total += n
    }
    return total
}
```

```typescript
// TypeScript
function Sum(numbers: number[]): Promise<number>

// Usage
const total = await Sum([1, 2, 3, 4, 5])  // 15
```

#### Maps

```go
// Go
func GetConfig() map[string]interface{} {
    return map[string]interface{}{
        "theme": "dark",
        "fontSize": 14,
        "enabled": true,
    }
}
```

```typescript
// TypeScript
function GetConfig(): Promise<Record<string, any>>

// Usage
const config = await GetConfig()
console.log(config.theme)  // "dark"
```

#### Structs

```go
// Go
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func GetUser(id int) (*User, error) {
    return &User{
        ID:    id,
        Name:  "Alice",
        Email: "alice@example.com",
    }, nil
}
```

```typescript
// TypeScript (auto-generated)
interface User {
    id: number
    name: string
    email: string
}

function GetUser(id: number): Promise<User>

// Usage
const user = await GetUser(1)
console.log(user.name)  // "Alice"
```

**Tags JSON:** use tags `json:` para controlar os nomes dos campos no TypeScript.

#### Tempo

```go
// Go
func GetTimestamp() time.Time {
    return time.Now()
}
```

```typescript
// TypeScript
function GetTimestamp(): Promise<Date>

// Usage
const timestamp = await GetTimestamp()
console.log(timestamp.toISOString())
```

#### Erros

```go
// Go
func Validate(input string) error {
    if input == "" {
        return errors.New("input cannot be empty")
    }
    return nil
}
```

```typescript
// TypeScript
function Validate(input: string): Promise<void>

// Usage
try {
    await Validate("")
} catch (error) {
    console.error(error)  // "input cannot be empty"
}
```

### Tipos não compatíveis

Estes tipos **não podem** ser transmitidos pela ponte:

- **Canais** (`chan T`)
- **Funções** (`func()`)
- **Interfaces** (exceto `interface{}` / `any`)
- **Ponteiros** (exceto para structs)
- **Campos não exportados** (letras minúsculas)

**Solução alternativa:** use IDs ou identificadores:

```go
// ❌ Can't pass file handle
func OpenFile(path string) (*os.File, error) {
    return os.Open(path)
}

// ✅ Return file ID instead
var files = make(map[string]*os.File)

func OpenFile(path string) (string, error) {
    file, err := os.Open(path)
    if err != nil {
        return "", err
    }
    id := generateID()
    files[id] = file
    return id, nil
}

func ReadFile(id string) ([]byte, error) {
    file := files[id]
    return io.ReadAll(file)
}

func CloseFile(id string) error {
    file := files[id]
    delete(files, id)
    return file.Close()
}
```

## Padrões avançados

### Passagem de contexto

Os serviços podem acessar o contexto da chamada:

```go
type UserService struct{}

func (s *UserService) GetCurrentUser(ctx context.Context) (*User, error) {
    // Access the calling window via the context value
    window, _ := ctx.Value(application.WindowKey).(application.Window)
    _ = window

    // Access the application
    app := application.Get()
    _ = app

    // Your logic
    return getCurrentUser(), nil
}
```

**O contexto fornece:**

- Janela que fez a chamada
- Instância do aplicativo
- Metadados da solicitação

### Transmissão de dados

Para grandes volumes de dados, use eventos em vez de valores de retorno:

```go
func ProcessLargeFile(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    lineNum := 0

    for scanner.Scan() {
        lineNum++
        // Emit progress events
        app.Event.Emit("file-progress", map[string]interface{}{
            "line": lineNum,
            "text": scanner.Text(),
        })
    }

    return scanner.Err()
}
```

```javascript
import { Events } from '@wailsio/runtime'
import { ProcessLargeFile } from './bindings/FileService'

// Listen for progress
Events.On('file-progress', (data) => {
    console.log(`Line ${data.line}: ${data.text}`)
})

// Start processing
await ProcessLargeFile('/path/to/large/file.txt')
```

### Cancelamento

Use o contexto para operações canceláveis:

```go
func LongRunningTask(ctx context.Context) error {
    for i := 0; i < 1000; i++ {
        // Check if cancelled
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            // Continue work
            time.Sleep(100 * time.Millisecond)
        }
    }
    return nil
}
```

**Observação:** o contexto é cancelado automaticamente quando o frontend se desconecta.

### Operações em lote

Reduza a sobrecarga da ponte agrupando operações em lotes:

```go
// ❌ Inefficient: N bridge calls
for _, item := range items {
    await ProcessItem(item)
}

// ✅ Efficient: 1 bridge call
await ProcessItems(items)
```

```go
func ProcessItems(items []Item) ([]Result, error) {
    results := make([]Result, len(items))
    for i, item := range items {
        results[i] = processItem(item)
    }
    return results, nil
}
```

## Depuração da ponte

### Ativar logs de depuração

```go
app := application.New(application.Options{
    Name:     "My App",
    LogLevel: slog.LevelDebug, // requires `import "log/slog"`
    // `Logger` is an optional *slog.Logger; the default-logger helper is
    // application.DefaultLogger(slog.Leveler) if you want to construct one explicitly.
})
```

**A saída mostra:**

- Chamadas de métodos
- Parâmetros
- Valores de retorno
- Erros
- Informações de tempo de execução

### Inspecionar os bindings gerados

Verifique `frontend/bindings/` para consultar o TypeScript gerado:

```javascript
// frontend/bindings/<full-go-import-path>/myservice.js (real generated shape)
import { Call as $Call } from "/wails/runtime.js";

export function MyMethod($0) {
    return $Call.ByID(1234567890, $0); // numeric method ID assigned by the generator
}
```

### Testar os serviços diretamente

Teste os serviços Go sem o frontend:

```go
func TestGreetService(t *testing.T) {
    service := &GreetService{prefix: "Hello, "}
    result := service.Greet("Test")
    if result != "Hello, Test!" {
        t.Errorf("Expected 'Hello, Test!', got '%s'", result)
    }
}
```

## Dicas de desempenho

### ✅ Faça

- **Agrupe operações em lotes** — reduza as chamadas à ponte
- **Use eventos para transmissão** — não retorne arrays grandes
- **Mantenha os métodos rápidos** — o ideal é &lt;100ms
- **Use goroutines** — para operações demoradas
- **Use cache no lado do Go** — evite cálculos repetidos

### ❌ Não faça

- **Não faça chamadas em excesso** — agrupe-as em lotes quando possível
- **Não retorne volumes enormes de dados** — use paginação ou transmissão
- **Não bloqueie** — use goroutines para operações demoradas
- **Não passe tipos complexos** — mantenha a simplicidade
- **Não ignore erros** — sempre trate-os

## Segurança

A ponte é segura por padrão:

1. **Somente lista de permissões** — apenas serviços registrados podem ser chamados
2. **Validação de tipos** — os argumentos são verificados em relação aos tipos Go
3. **Sem eval()** — o frontend não pode executar código Go arbitrário
4. **Sem uso indevido de reflexão** — apenas métodos exportados são acessíveis

**Práticas recomendadas:**

- **Valide as entradas** no Go (não confie no frontend)
- **Use o contexto** para autenticação/autorização
- **Limite a taxa** de operações dispendiosas
- **Sanitize** caminhos de arquivos e entradas do usuário

## Próximos passos

**Sistema de build** — saiba como o Wails compila e empacota seu aplicativo [Saiba mais →](/concepts/build-system/)

**Serviços** — aprofunde-se no sistema de serviços [Saiba mais →](/features/bindings/services/)

**Eventos** — Use eventos para comunicação por publicação/assinatura [Saiba mais →](/features/events/system/)

---

**Dúvidas sobre a ponte?** Pergunte no [Discord](https://discord.gg/JDdSxwjhGf) ou consulte os [exemplos de bindings](https://github.com/wailsapp/wails/tree/master/v3/examples/binding).
