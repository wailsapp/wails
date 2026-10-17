---
title: "性能优化"
description: "优化 Wails 应用以获得最佳性能"
slug: "guides/performance"
sourcePath: "guides/performance.md"
---

## 概述

优化 Wails 应用的速度、内存效率和响应能力。

## 前端优化

### 包体积

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

### 代码拆分

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

### 资源优化

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

## 后端优化

### 高效绑定

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

### 缓存

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

### 使用 goroutine 执行耗时操作

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

## 内存优化

### 避免内存泄漏

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

### 资源池化

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

## 事件优化

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

### 批量更新

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

## 构建优化

### 二进制文件大小

`wails3 build`本身没有`-ldflags`标志——随附的 Taskfile 已向`go build`传递`-ldflags="-s -w"`。若要进一步缩小二进制文件，请编辑构建任务或直接调用`go build`：

```bash
# Strip debug symbols and trim file paths
go build -ldflags="-s -w" -trimpath -o bin/myapp
```

### 编译速度

```bash
# Use build cache
go build -buildmode=default

# Parallel compilation
go build -p 8
```

## 性能分析

### CPU 性能分析

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

### 内存性能分析

```go
import "runtime/pprof"

func profileMemory() {
    f, _ := os.Create("mem.prof")
    defer f.Close()
    
    runtime.GC()
    pprof.WriteHeapProfile(f)
}
```

### 分析性能数据

```bash
# View CPU profile
go tool pprof cpu.prof

# View memory profile
go tool pprof mem.prof

# Web interface
go tool pprof -http=:8080 cpu.prof
```

## 最佳实践

### ✅ 应该做

- 先进行性能分析，再优化
- 缓存开销较大的操作
- 对大型数据集使用分页
- 对频繁触发的事件进行防抖
- 使用资源池
- 清理 goroutine
- 优化包体积
- 使用延迟加载

### ❌ 不应该做

- 不要过早优化
- 不要忽视内存泄漏
- 不要阻塞主线程
- 不要返回庞大的数据集
- 不要跳过性能分析
- 不要忘记清理资源

## 性能检查清单

- [ ] 已优化前端包
- [ ] 已压缩图像
- [ ] 已实施代码拆分
- [ ] 后端方法已采用分页
- [ ] 已实施缓存
- [ ] 已清理 goroutine
- [ ] 已对事件进行防抖
- [ ] 已优化二进制文件大小
- [ ] 已完成性能分析
- [ ] 已修复内存泄漏

## 后续步骤

- [架构](/guides/architecture/) - 应用架构模式
- [测试](/guides/testing/) - 测试应用
- [构建](/guides/build/building/) - 构建经过优化的二进制文件
