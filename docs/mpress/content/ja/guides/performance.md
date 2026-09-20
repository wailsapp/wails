---
title: "パフォーマンス最適化"
description: "Wails アプリケーションを最適化してパフォーマンスを最大限に高める"
slug: "guides/performance"
sourcePath: "guides/performance.md"
---

## 概要

Wails アプリケーションの速度、メモリ効率、応答性を最適化します。

## フロントエンドの最適化

### バンドルサイズ

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

### コード分割

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

### アセットの最適化

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

## バックエンドの最適化

### 効率的なバインディング

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

### キャッシュ

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

### 長時間実行する処理での goroutine の使用

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

## メモリの最適化

### メモリリークの防止

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

### リソースのプール

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

## イベントの最適化

### イベントのデバウンス

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

### 更新の一括処理

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

## ビルドの最適化

### バイナリサイズ

`wails3 build` 自体には `-ldflags` フラグがありません。配布されている Taskfile では、すでに `go build` に `-ldflags="-s -w"` を渡しています。バイナリをさらに縮小するには、ビルドタスクを編集するか、`go build` を直接呼び出します。

```bash
# Strip debug symbols and trim file paths
go build -ldflags="-s -w" -trimpath -o bin/myapp
```

### コンパイル速度

```bash
# Use build cache
go build -buildmode=default

# Parallel compilation
go build -p 8
```

## プロファイリング

### CPU プロファイリング

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

### メモリプロファイリング

```go
import "runtime/pprof"

func profileMemory() {
    f, _ := os.Create("mem.prof")
    defer f.Close()
    
    runtime.GC()
    pprof.WriteHeapProfile(f)
}
```

### プロファイルの解析

```bash
# View CPU profile
go tool pprof cpu.prof

# View memory profile
go tool pprof mem.prof

# Web interface
go tool pprof -http=:8080 cpu.prof
```

## ベストプラクティス

### ✅ 推奨事項

- 最適化する前にプロファイリングする
- コストの高い処理をキャッシュする
- 大規模なデータセットにはページネーションを使用する
- 頻繁に発生するイベントをデバウンスする
- リソースをプールする
- goroutine をクリーンアップする
- バンドルサイズを最適化する
- 遅延読み込みを使用する

### ❌ 禁止事項

- 時期尚早な最適化をしない
- メモリリークを放置しない
- メインスレッドをブロックしない
- 巨大なデータセットを返さない
- プロファイリングを省略しない
- クリーンアップを忘れない

## パフォーマンスチェックリスト

- [ ] フロントエンドバンドルを最適化済み
- [ ] 画像を圧縮済み
- [ ] コード分割を実装済み
- [ ] バックエンドメソッドにページネーションを実装済み
- [ ] キャッシュを実装済み
- [ ] goroutine をクリーンアップ済み
- [ ] イベントをデバウンス済み
- [ ] バイナリサイズを最適化済み
- [ ] プロファイリングを実施済み
- [ ] メモリリークを修正済み

## 次のステップ

- [アーキテクチャ](/guides/architecture/) - アプリケーションのアーキテクチャパターン
- [テスト](/guides/testing/) - アプリケーションをテストする
- [ビルド](/guides/build/building/) - 最適化されたバイナリをビルドする
