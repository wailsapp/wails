---
title: "Padrões de codificação"
description: "Estilo de código, convenções e práticas recomendadas para o Wails v3"
slug: "contributing/standards"
sourcePath: "contributing/standards.md"
---

## Estilo e convenções de código

Seguir padrões de codificação consistentes facilita a leitura, a manutenção e a contribuição para a base de código.

## Padrões de código Go

### Formatação de código

Use as ferramentas padrão de formatação do Go:

```bash
# Format all code
gofmt -w .

# Use goimports for import organization
goimports -w .
```

**Obrigatório:** todo código Go deve passar por `gofmt` e `goimports` antes do commit.

### Convenções de nomenclatura

**Pacotes:**

- Use letras minúsculas e, quando possível, uma única palavra
- `package application`, `package events`
- Evite sublinhados ou a combinação de letras maiúsculas e minúsculas

**Nomes exportados:**

- Use PascalCase para tipos, funções e constantes
- `type WebviewWindow struct`, `func NewApplication()`

**Nomes não exportados:**

- Use camelCase para tipos, funções e variáveis internos
- `type windowImpl struct`, `func createWindow()`

**Interfaces:**

- Nomeie de acordo com o comportamento: `Reader`, `Writer`, `Handler`
- Para interfaces com um único método, use nomes com o sufixo `-er`

```go
// Good
type Closer interface {
    Close() error
}

// Avoid
type CloseInterface interface {
    Close() error
}
```

### Tratamento de erros

**Sempre verifique os erros:**

```go
// Good
result, err := doSomething()
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// Bad - ignoring errors
result, _ := doSomething()
```

**Use o encapsulamento de erros:**

```go
// Wrap errors to provide context
if err := validate(); err != nil {
    return fmt.Errorf("validation failed: %w", err)
}
```

**Crie tipos de erro personalizados quando necessário:**

```go
type ValidationError struct {
    Field string
    Value string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("invalid value %q for field %q", e.Value, e.Field)
}
```

### Comentários e documentação

**Comentários de pacote:**

```go
// Package application provides the core Wails application runtime.
//
// It handles window management, event dispatching, and service lifecycle.
package application
```

**Declarações exportadas:**

```go
// NewApplication creates a new Wails application with the given options.
//
// The application must be started with Run() or RunWithContext().
func NewApplication(opts Options) *Application {
    // ...
}
```

**Comentários de implementação:**

```go
// processEvent handles incoming events from the runtime.
// It dispatches to registered handlers and manages event lifecycle.
func (a *Application) processEvent(event *Event) {
    // Validate event before processing
    if event == nil {
        return
    }

    // Find and invoke handlers
    // ...
}
```

### Estrutura de funções e métodos

**Mantenha as funções focadas:**

```go
// Good - single responsibility
func (w *Window) setTitle(title string) {
    w.title = title
    w.updateNativeTitle()
}

// Bad - doing too much
func (w *Window) updateEverything() {
    w.setTitle(w.title)
    w.setSize(w.width, w.height)
    w.setPosition(w.x, w.y)
    // ... 20 more operations
}
```

**Use retornos antecipados:**

```go
// Good
func validate(input string) error {
    if input == "" {
        return errors.New("empty input")
    }

    if len(input) > 100 {
        return errors.New("input too long")
    }

    return nil
}

// Avoid deep nesting
```

### Concorrência

**Use context para cancelamento:**

```go
func (a *Application) RunWithContext(ctx context.Context) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-a.done:
        return nil
    }
}
```

**Proteja o estado compartilhado com mutexes:**

```go
type SafeCounter struct {
    mu    sync.Mutex
    count int
}

func (c *SafeCounter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}
```

**Evite vazamentos de goroutines:**

```go
// Good - goroutine has exit condition
func (a *Application) startWorker(ctx context.Context) {
    go func() {
        for {
            select {
            case <-ctx.Done():
                return  // Clean exit
            case work := <-a.workChan:
                a.process(work)
            }
        }
    }()
}
```

### Testes

**Nomenclatura de arquivos de teste:**

```go
// Implementation: window.go
// Tests: window_test.go
```

**Testes orientados por tabelas:**

```go
func TestValidate(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"empty input", "", true},
        {"valid input", "hello", false},
        {"too long", strings.Repeat("a", 101), true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validate(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Padrões de JavaScript/TypeScript

### Formatação de código

Use o Prettier para manter a formatação consistente:

```json
{
  "semi": false,
  "singleQuote": true,
  "tabWidth": 2,
  "trailingComma": "es5"
}
```

### Convenções de nomenclatura

**Variáveis e funções:**

- camelCase: `const userName = "John"`

**Classes e tipos:**

- PascalCase: `class WindowManager`

**Constantes:**

- UPPER<em>SNAKE</em>CASE: `const MAX_RETRIES = 3`

### TypeScript

**Use tipos explícitos:**

```typescript
// Good
function greet(name: string): string {
    return `Hello, ${name}`
}

// Avoid implicit any
function process(data) {  // Bad
    return data
}
```

**Defina interfaces:**

```typescript
interface WindowOptions {
    title: string
    width: number
    height: number
}

function createWindow(options: WindowOptions): void {
    // ...
}
```

## Formato das mensagens de commit

Use [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Tipos:**

- `feat`: novo recurso
- `fix`: correção de bug
- `docs`: alterações na documentação
- `refactor`: refatoração de código
- `test`: adição ou atualização de testes
- `chore`: tarefas de manutenção

**Exemplos:**

```
feat(window): add SetAlwaysOnTop method

Implement SetAlwaysOnTop for keeping windows above others.
Adds platform implementations for macOS, Windows, and Linux.

Closes #123
```

```
fix(events): prevent event handler memory leak

Event listeners were not being properly cleaned up when
windows were closed. This adds explicit cleanup in the
window destructor.
```

## Diretrizes para pull requests

### Antes de enviar

- [ ] O código passa por `gofmt` e `goimports`
- [ ] Todos os testes passam (`go test ./...`)
- [ ] O novo código tem testes
- [ ] A documentação foi atualizada, se necessário
- [ ] As mensagens de commit seguem as convenções
- [ ] Não há conflitos de merge com `master`

### Modelo de descrição do PR

```markdown
## Description
Brief description of what this PR does.

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
How was this tested?

## Checklist
- [ ] Tests pass
- [ ] Documentation updated
- [ ] No breaking changes (or documented)
```

## Processo de revisão de código

### Como revisor

- Seja construtivo e respeitoso
- Concentre-se na qualidade do código, não em preferências pessoais
- Explique por que as alterações são sugeridas
- Aprove quando estiver satisfeito

### Como autor

- Responda a todos os comentários
- Peça esclarecimentos, se necessário
- Faça as alterações solicitadas ou explique por que não as fará
- Esteja aberto a feedback

## Práticas recomendadas

### Desempenho

- Evite otimizações prematuras
- Faça a análise de desempenho antes de otimizar
- Use benchmarks em código crítico para o desempenho

```go
func BenchmarkProcess(b *testing.B) {
    for i := 0; i < b.N; i++ {
        process(testData)
    }
}
```

### Segurança

- Valide todas as entradas do usuário
- Sanitize os dados antes de exibi-los
- Use `crypto/rand` para dados aleatórios
- Nunca registre informações confidenciais em logs

### Documentação

- Documente as APIs exportadas
- Inclua exemplos na documentação
- Atualize a documentação ao alterar APIs
- Mantenha os arquivos README atualizados

## Código específico de plataforma

### Nomenclatura de arquivos

```
window.go           // Common interface
window_darwin.go    // macOS implementation
window_windows.go   // Windows implementation
window_linux.go     // Linux implementation
```

### Tags de build

```go
//go:build darwin

package application

// macOS-specific code
```

## Linting

Execute os linters antes de fazer o commit:

```bash
# golangci-lint (recommended)
golangci-lint run

# Individual linters
go vet ./...
staticcheck ./...
```

## Dúvidas?

Se tiver dúvidas sobre alguma norma:

- Consulte o código existente para ver exemplos
- Pergunte no [Discord](https://discord.gg/JDdSxwjhGf)
- Abra uma discussão no GitHub
