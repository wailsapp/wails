---
title: "Otimização de desempenho"
description: "Otimize seu aplicativo Wails para obter o máximo desempenho"
slug: "guides/performance"
sourcePath: "guides/performance.md"
---

## Visão geral

Otimize seu aplicativo Wails para melhorar a velocidade, a eficiência de memória e a capacidade de resposta.

## Otimização do frontend

### Tamanho do bundle

```javascript
// vite.config.js
export default {
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['react', 'react-dom'],
        },
      },
    },
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: true,
      },
    },
  },
}
```

### Divisão de código

```javascript
// Lazy load components
const Settings = lazy(() => import('./Settings'))

function App() {
  return (
    <Suspense fallback={<Loading />}>
      <Settings />
    </Suspense>
  )
}
```

### Otimização de recursos

```javascript
// Optimise images
import { defineConfig } from 'vite'
import imagemin from 'vite-plugin-imagemin'

export default defineConfig({
  plugins: [
    imagemin({
      gifsicle: { optimizationLevel: 3 },
      optipng: { optimizationLevel: 7 },
      svgo: { plugins: [{ removeViewBox: false }] },
    }),
  ],
})
```

## Otimização do backend

### Bindings eficientes

```go
// ❌ Bad: Return everything
func (s *Service) GetAllData() []Data {
    return s.db.FindAll() // Could be huge
}

// ✅ Good: Paginate
func (s *Service) GetData(page, size int) (*PagedData, error) {
    return s.db.FindPaged(page, size)
}
```

### Cache

```go
type CachedService struct {
    cache *lru.Cache
    ttl   time.Duration
}

func (s *CachedService) GetData(key string) (interface{}, error) {
    // Check cache
    if val, ok := s.cache.Get(key); ok {
        return val, nil
    }
    
    // Fetch and cache
    data, err := s.fetchData(key)
    if err != nil {
        return nil, err
    }
    
    s.cache.Add(key, data)
    return data, nil
}
```

### Goroutines para operações demoradas

```go
func (s *Service) ProcessLargeFile(path string) error {
    // Process in background
    go func() {
        result, err := s.process(path)
        if err != nil {
            s.app.Event.Emit("process-error", err.Error())
            return
        }
        s.app.Event.Emit("process-complete", result)
    }()
    
    return nil
}
```

## Otimização de memória

### Evite vazamentos de memória

```go
// ❌ Bad: Goroutine leak
func (s *Service) StartPolling() {
    ticker := time.NewTicker(1 * time.Second)
    go func() {
        for range ticker.C {
            s.poll()
        }
    }()
    // ticker never stopped!
}

// ✅ Good: Proper cleanup
func (s *Service) StartPolling() {
    ticker := time.NewTicker(1 * time.Second)
    s.stopChan = make(chan bool)
    
    go func() {
        for {
            select {
            case <-ticker.C:
                s.poll()
            case <-s.stopChan:
                ticker.Stop()
                return
            }
        }
    }()
}

func (s *Service) StopPolling() {
    close(s.stopChan)
}
```

### Agrupamento de recursos

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func processData(data []byte) []byte {
    buf := bufferPool.Get().(*bytes.Buffer)
    defer bufferPool.Put(buf)
    
    buf.Reset()
    buf.Write(data)
    // Process...
    return buf.Bytes()
}
```

## Otimização de eventos

### Aplique debounce aos eventos

```javascript
// Debounce frequent events
let debounceTimer
function handleInput(value) {
    clearTimeout(debounceTimer)
    debounceTimer = setTimeout(() => {
        UpdateData(value)
    }, 300)
}
```

### Atualizações em lote

```go
type BatchProcessor struct {
    items []Item
    mu    sync.Mutex
    timer *time.Timer
}

func (b *BatchProcessor) Add(item Item) {
    b.mu.Lock()
    defer b.mu.Unlock()
    
    b.items = append(b.items, item)
    
    if b.timer == nil {
        b.timer = time.AfterFunc(100*time.Millisecond, b.flush)
    }
}

func (b *BatchProcessor) flush() {
    b.mu.Lock()
    items := b.items
    b.items = nil
    b.timer = nil
    b.mu.Unlock()
    
    // Process batch
    processBatch(items)
}
```

## Otimização da compilação

### Tamanho do binário

O próprio `wails3 build` não tem uma opção `-ldflags` — o Taskfile fornecido já passa `-ldflags="-s -w"` para `go build`. Para reduzir ainda mais o binário, edite a tarefa de compilação ou chame `go build` diretamente:

```bash
# Strip debug symbols and trim file paths
go build -ldflags="-s -w" -trimpath -o bin/myapp
```

### Velocidade de compilação

```bash
# Use build cache
go build -buildmode=default

# Parallel compilation
go build -p 8
```

## Análise de desempenho

### Análise de CPU

```go
import "runtime/pprof"

func profileCPU() {
    f, _ := os.Create("cpu.prof")
    defer f.Close()
    
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()
    
    // Code to profile
}
```

### Análise de memória

```go
import "runtime/pprof"

func profileMemory() {
    f, _ := os.Create("mem.prof")
    defer f.Close()
    
    runtime.GC()
    pprof.WriteHeapProfile(f)
}
```

### Analise os perfis

```bash
# View CPU profile
go tool pprof cpu.prof

# View memory profile
go tool pprof mem.prof

# Web interface
go tool pprof -http=:8080 cpu.prof
```

## Práticas recomendadas

### ✅ Faça

- Analise o desempenho antes de otimizar
- Armazene em cache as operações de alto custo
- Use paginação para grandes conjuntos de dados
- Aplique debounce aos eventos frequentes
- Agrupe recursos
- Encerre corretamente as goroutines
- Otimize o tamanho do bundle
- Use carregamento sob demanda

### ❌ Não faça

- Não otimize prematuramente
- Não ignore vazamentos de memória
- Não bloqueie a thread principal
- Não retorne conjuntos de dados enormes
- Não deixe de analisar o desempenho
- Não se esqueça da limpeza

## Lista de verificação de desempenho

- [ ] Bundle do frontend otimizado
- [ ] Imagens compactadas
- [ ] Divisão de código implementada
- [ ] Paginação implementada nos métodos do backend
- [ ] Cache implementado
- [ ] Goroutines encerradas corretamente
- [ ] Debounce aplicado aos eventos
- [ ] Tamanho do binário otimizado
- [ ] Análise de desempenho concluída
- [ ] Vazamentos de memória corrigidos

## Próximas etapas

- [Arquitetura](/guides/architecture/) — Padrões de arquitetura de aplicativos
- [Testes](/guides/testing/) — Teste seu aplicativo
- [Compilação](/guides/build/building/) — Compile binários otimizados
