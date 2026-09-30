---
title: "Leistungsoptimierung"
description: "Optimieren Sie Ihre Wails-Anwendung für maximale Leistung"
slug: "guides/performance"
sourcePath: "guides/performance.md"
---

## Übersicht

Optimieren Sie Ihre Wails-Anwendung hinsichtlich Geschwindigkeit, Speichereffizienz und Reaktionsfähigkeit.

## Frontend-Optimierung

### Bundle-Größe

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

### Codeaufteilung

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

### Asset-Optimierung

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

## Backend-Optimierung

### Effiziente Bindings

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

### Caching

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

### Goroutinen für lang laufende Vorgänge

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

## Speicheroptimierung

### Speicherlecks vermeiden

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

### Ressourcen-Pooling

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

## Ereignisoptimierung

### Ereignisse entprellen

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

### Aktualisierungen bündeln

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

## Build-Optimierung

### Binärdateigröße

`wails3 build` selbst verfügt über kein `-ldflags`-Flag – das mitgelieferte Taskfile übergibt bereits `-ldflags="-s -w"` an `go build`. Um die Binärdatei weiter zu verkleinern, bearbeiten Sie den Build-Task oder rufen Sie `go build` direkt auf:

```bash
# Strip debug symbols and trim file paths
go build -ldflags="-s -w" -trimpath -o bin/myapp
```

### Kompilierungsgeschwindigkeit

```bash
# Use build cache
go build -buildmode=default

# Parallel compilation
go build -p 8
```

## Profiling

### CPU-Profiling

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

### Speicher-Profiling

```go
import "runtime/pprof"

func profileMemory() {
    f, _ := os.Create("mem.prof")
    defer f.Close()
    
    runtime.GC()
    pprof.WriteHeapProfile(f)
}
```

### Profile analysieren

```bash
# View CPU profile
go tool pprof cpu.prof

# View memory profile
go tool pprof mem.prof

# Web interface
go tool pprof -http=:8080 cpu.prof
```

## Bewährte Methoden

### ✅ Empfohlen

- Vor dem Optimieren ein Profil erstellen
- Aufwendige Vorgänge zwischenspeichern
- Für große Datensätze Paginierung verwenden
- Häufige Ereignisse entprellen
- Ressourcen-Pooling verwenden
- Goroutinen bereinigen
- Bundle-Größe optimieren
- Lazy Loading verwenden

### ❌ Nicht empfohlen

- Nicht vorzeitig optimieren
- Speicherlecks nicht ignorieren
- Den Hauptthread nicht blockieren
- Keine riesigen Datensätze zurückgeben
- Profiling nicht auslassen
- Bereinigung nicht vergessen

## Leistungscheckliste

- [ ] Frontend-Bundle optimiert
- [ ] Bilder komprimiert
- [ ] Codeaufteilung implementiert
- [ ] Backend-Methoden paginiert
- [ ] Caching implementiert
- [ ] Goroutinen bereinigt
- [ ] Ereignisse entprellt
- [ ] Binärdateigröße optimiert
- [ ] Profiling durchgeführt
- [ ] Speicherlecks behoben

## Nächste Schritte

- [Architektur](/guides/architecture/) – Architekturmuster für Anwendungen
- [Tests](/guides/testing/) – Testen Sie Ihre Anwendung
- [Builds](/guides/build/building/) – Erstellen Sie optimierte Binärdateien
