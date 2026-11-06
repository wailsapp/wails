---
title: "Оптимизация производительности"
description: "Оптимизируйте приложение Wails для максимальной производительности"
slug: "guides/performance"
sourcePath: "guides/performance.md"
---

## Обзор

Оптимизируйте приложение Wails, чтобы повысить его быстродействие, эффективность использования памяти и скорость отклика.

## Оптимизация фронтенда

### Размер пакета

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

### Разделение кода

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

### Оптимизация ресурсов

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

## Оптимизация бэкенда

### Эффективные привязки

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

### Кэширование

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

### Горутины для длительных операций

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

## Оптимизация памяти

### Предотвращение утечек памяти

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

### Пулы ресурсов

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

## Оптимизация событий

### Устранение дребезга событий

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

### Пакетные обновления

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

## Оптимизация сборки

### Размер исполняемого файла

У самой команды `wails3 build` нет флага `-ldflags` — поставляемый Taskfile уже передаёт `-ldflags="-s -w"` в `go build`. Чтобы дополнительно уменьшить исполняемый файл, измените задачу сборки или вызовите `go build` напрямую:

```bash
# Strip debug symbols and trim file paths
go build -ldflags="-s -w" -trimpath -o bin/myapp
```

### Скорость компиляции

```bash
# Use build cache
go build -buildmode=default

# Parallel compilation
go build -p 8
```

## Профилирование

### Профилирование ЦП

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

### Профилирование памяти

```go
import "runtime/pprof"

func profileMemory() {
    f, _ := os.Create("mem.prof")
    defer f.Close()
    
    runtime.GC()
    pprof.WriteHeapProfile(f)
}
```

### Анализ профилей

```bash
# View CPU profile
go tool pprof cpu.prof

# View memory profile
go tool pprof mem.prof

# Web interface
go tool pprof -http=:8080 cpu.prof
```

## Рекомендации

### ✅ Рекомендуется

- Выполняйте профилирование перед оптимизацией
- Кэшируйте результаты ресурсоёмких операций
- Используйте разбиение на страницы для больших наборов данных
- Устраняйте дребезг часто возникающих событий
- Используйте пулы ресурсов
- Завершайте горутины корректно
- Оптимизируйте размер пакета
- Используйте отложенную загрузку

### ❌ Не рекомендуется

- Не выполняйте преждевременную оптимизацию
- Не игнорируйте утечки памяти
- Не блокируйте основной поток
- Не возвращайте огромные наборы данных
- Не пропускайте профилирование
- Не забывайте освобождать ресурсы

## Контрольный список производительности

- [ ] Пакет фронтенда оптимизирован
- [ ] Изображения сжаты
- [ ] Разделение кода реализовано
- [ ] В методах бэкенда реализовано разбиение на страницы
- [ ] Кэширование реализовано
- [ ] Горутины завершаются корректно
- [ ] Дребезг событий устранён
- [ ] Размер исполняемого файла оптимизирован
- [ ] Профилирование выполнено
- [ ] Утечки памяти устранены

## Дальнейшие действия

- [Архитектура](/guides/architecture/) — шаблоны архитектуры приложений
- [Тестирование](/guides/testing/) — протестируйте приложение
- [Сборка](/guides/build/building/) — соберите оптимизированные исполняемые файлы
