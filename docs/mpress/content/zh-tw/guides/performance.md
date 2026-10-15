---
title: "效能最佳化"
description: "最佳化 Wails 應用程式以達到最高效能"
slug: "guides/performance"
sourcePath: "guides/performance.md"
---

## 概覽

最佳化 Wails 應用程式的速度、記憶體使用效率與回應速度。

## 前端最佳化

### 打包產物大小

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

### 程式碼分割

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

### 資源最佳化

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

## 後端最佳化

### 高效繫結

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

### 快取

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

### 使用 Goroutine 執行耗時操作

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

## 記憶體最佳化

### 避免記憶體洩漏

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

### 使用資源池

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

## 事件最佳化

### 事件防抖

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

### 批次更新

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

## 建置最佳化

### 二進位檔大小

`wails3 build`本身沒有`-ldflags`旗標——隨附的 Taskfile 已將`-ldflags="-s -w"`傳遞給`go build`。若要進一步縮小二進位檔，請編輯建置工作或直接呼叫`go build`：

```bash
# Strip debug symbols and trim file paths
go build -ldflags="-s -w" -trimpath -o bin/myapp
```

### 編譯速度

```bash
# Use build cache
go build -buildmode=default

# Parallel compilation
go build -p 8
```

## 效能分析

### CPU 效能分析

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

### 記憶體效能分析

```go
import "runtime/pprof"

func profileMemory() {
    f, _ := os.Create("mem.prof")
    defer f.Close()
    
    runtime.GC()
    pprof.WriteHeapProfile(f)
}
```

### 分析效能剖析資料

```bash
# View CPU profile
go tool pprof cpu.prof

# View memory profile
go tool pprof mem.prof

# Web interface
go tool pprof -http=:8080 cpu.prof
```

## 最佳實務

### ✅ 建議做法

- 先進行效能分析，再進行最佳化
- 快取成本高昂的操作結果
- 對大型資料集使用分頁
- 對頻繁事件進行防抖處理
- 集中管理並重複使用資源
- 清理 Goroutine
- 最佳化打包產物大小
- 使用延遲載入

### ❌ 避免做法

- 不要過早最佳化
- 不要忽略記憶體洩漏
- 不要阻塞主執行緒
- 不要傳回龐大的資料集
- 不要略過效能分析
- 不要忘記清理資源

## 效能檢查清單

- [ ] 已最佳化前端打包產物
- [ ] 已壓縮圖片
- [ ] 已實作程式碼分割
- [ ] 後端方法已使用分頁
- [ ] 已實作快取
- [ ] 已清理 Goroutine
- [ ] 已對事件進行防抖處理
- [ ] 已最佳化二進位檔大小
- [ ] 已完成效能分析
- [ ] 已修正記憶體洩漏

## 後續步驟

- [架構](/guides/architecture/) - 應用程式架構模式
- [測試](/guides/testing/) - 測試您的應用程式
- [建置](/guides/build/building/) - 建置最佳化的二進位檔
