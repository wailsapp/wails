---
title: "성능 최적화"
description: "Wails 애플리케이션의 성능을 최대한으로 최적화합니다"
slug: "guides/performance"
sourcePath: "guides/performance.md"
---

## 개요

Wails 애플리케이션의 속도, 메모리 효율성 및 응답성을 최적화합니다.

## 프런트엔드 최적화

### 번들 크기

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

### 코드 분할

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

### 에셋 최적화

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

## 백엔드 최적화

### 효율적인 바인딩

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

### 캐싱

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

### 장시간 작업에 고루틴 사용

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

## 메모리 최적화

### 메모리 누수 방지

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

### 리소스 풀링

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

## 이벤트 최적화

### 이벤트 디바운싱

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

### 일괄 업데이트

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

## 빌드 최적화

### 바이너리 크기

`wails3 build` 자체에는 `-ldflags` 플래그가 없습니다. 함께 제공되는 Taskfile에서 이미 `go build`에 `-ldflags="-s -w"`를 전달합니다. 바이너리 크기를 더 줄이려면 빌드 작업을 수정하거나 `go build`을 직접 호출하세요:

```bash
# Strip debug symbols and trim file paths
go build -ldflags="-s -w" -trimpath -o bin/myapp
```

### 컴파일 속도

```bash
# Use build cache
go build -buildmode=default

# Parallel compilation
go build -p 8
```

## 프로파일링

### CPU 프로파일링

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

### 메모리 프로파일링

```go
import "runtime/pprof"

func profileMemory() {
    f, _ := os.Create("mem.prof")
    defer f.Close()
    
    runtime.GC()
    pprof.WriteHeapProfile(f)
}
```

### 프로파일 분석

```bash
# View CPU profile
go tool pprof cpu.prof

# View memory profile
go tool pprof mem.prof

# Web interface
go tool pprof -http=:8080 cpu.prof
```

## 권장 사항

### ✅ 권장

- 최적화 전에 프로파일링하세요
- 비용이 많이 드는 작업을 캐싱하세요
- 대규모 데이터 세트에는 페이지네이션을 사용하세요
- 빈번한 이벤트를 디바운싱하세요
- 리소스를 풀링하세요
- 고루틴을 정리하세요
- 번들 크기를 최적화하세요
- 지연 로딩을 사용하세요

### ❌ 금지

- 너무 일찍 최적화하지 마세요
- 메모리 누수를 무시하지 마세요
- 메인 스레드를 차단하지 마세요
- 방대한 데이터 세트를 반환하지 마세요
- 프로파일링을 건너뛰지 마세요
- 정리 작업을 잊지 마세요

## 성능 체크리스트

- [ ] 프런트엔드 번들 최적화 완료
- [ ] 이미지 압축 완료
- [ ] 코드 분할 구현 완료
- [ ] 백엔드 메서드에 페이지네이션 적용 완료
- [ ] 캐싱 구현 완료
- [ ] 고루틴 정리 완료
- [ ] 이벤트 디바운싱 완료
- [ ] 바이너리 크기 최적화 완료
- [ ] 프로파일링 완료
- [ ] 메모리 누수 수정 완료

## 다음 단계

- [아키텍처](/guides/architecture/) - 애플리케이션 아키텍처 패턴
- [테스트](/guides/testing/) - 애플리케이션 테스트
- [빌드](/guides/build/building/) - 최적화된 바이너리 빌드
