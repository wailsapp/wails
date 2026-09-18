---
title: "Go–フロントエンドブリッジ"
description: "Wails が Go と JavaScript の直接通信を実現する仕組みを詳しく解説します"
slug: "concepts/bridge"
sourcePath: "concepts/bridge.md"
---

## Go と JavaScript の直接通信

Wails は Go と JavaScript の間に<strong>メモリ内で動作する直接ブリッジ</strong>を提供し、HTTP のオーバーヘッド、プロセス境界、シリアライズのボトルネックなしでシームレスな通信を実現します。

## 全体像

```d2
direction: right

Frontend: フロントエンド（JavaScript） {
  UI: React/Vue/Vanilla {
    shape: rectangle
    style.fill: "#8B5CF6"
  }

  Bindings: 自動生成バインディング {
    shape: rectangle
    style.fill: "#A78BFA"
  }
}

Bridge: Wails ブリッジ {
  Encoder: JSON エンコーダー {
    shape: rectangle
    style.fill: "#10B981"
  }

  Router: メソッドルーター {
    shape: diamond
    style.fill: "#10B981"
  }

  Decoder: JSON デコーダー {
    shape: rectangle
    style.fill: "#10B981"
  }

  TypeGen: 型ジェネレーター {
    shape: rectangle
    style.fill: "#10B981"
  }
}

Backend: バックエンド（Go） {
  Services: サービス {
    shape: rectangle
    style.fill: "#00ADD8"
  }

  Registry: サービスレジストリ {
    shape: rectangle
    style.fill: "#00ADD8"
  }
}

Frontend.UI -> Frontend.Bindings: "import { Method }"
Frontend.Bindings -> Bridge.Encoder: "Method('arg') を呼び出す"
Bridge.Encoder -> Bridge.Router: JSON にエンコード
Bridge.Router -> Backend.Registry: サービスを検索
Backend.Registry -> Backend.Services: メソッドを呼び出す
Backend.Services -> Bridge.Decoder: 結果を返す
Bridge.Decoder -> Frontend.Bindings: JS にデコード
Frontend.Bindings -> Frontend.UI: Promise が解決
Bridge.TypeGen -> Frontend.Bindings: 型を生成
```

<strong>要点：</strong>HTTP、IPC、プロセス境界はありません。<strong>型安全性</strong>を備えた<strong>直接的な関数呼び出し</strong>だけです。

## 仕組み：ステップごとの解説

### 1. サービスの登録（起動時）

アプリケーションの起動時に、Wails がサービスをスキャンします。

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

**Wails が行う処理：**

1. エクスポートされたメソッドを検出するため、**構造体をスキャンします**
2. **型情報を抽出します**（パラメーター、戻り値の型）
3. メソッド名を関数に対応付ける<strong>レジストリを構築します</strong>
4. 完全な型定義を含む<strong>TypeScript バインディングを生成します</strong>

### 2. バインディングの生成（ビルド時）

Wails は TypeScript バインディングを自動的に生成します。

```typescript
// Auto-generated: frontend/bindings/GreetService.ts
export function Greet(name: string): Promise<string>
export function Add(a: number, b: number): Promise<number>
```

**型のマッピング：**

| Go の型 | TypeScript の型 |
| --- | --- |
| `string` | `string` |
| `int`、`int32`、`int64` | `number` |
| `float32`、`float64` | `number` |
| `bool` | `boolean` |
| `[]T` | `T[]` |
| `map[string]T` | `Record<string, T>` |
| `struct` | `interface` |
| `time.Time` | `Date` |
| `error` | 例外（スローされる） |

### 3. フロントエンドからの呼び出し（実行時）

開発者が JavaScript から Go メソッドを呼び出します。

```javascript
import { Greet, Add } from './bindings/GreetService'

// Call Go from JavaScript
const greeting = await Greet("World")
console.log(greeting)  // "Hello, World!"

const sum = await Add(5, 3)
console.log(sum)  // 8
```

**実行される処理：**

1. **バインディング関数が呼び出されます** - `Greet("World")`
2. **メッセージが作成されます** - `{ service: "GreetService", method: "Greet", args: ["World"] }`
3. **ブリッジに送信** — WebView の JavaScript ブリッジ経由
4. **Promise を返す** — レスポンスを待機

### 4. ブリッジでの処理（実行時）

ブリッジはメッセージを受信し、処理します。

```d2
direction: down

Receive: メッセージを受信 {
  shape: rectangle
  style.fill: "#10B981"
}

Parse: JSON を解析 {
  shape: rectangle
}

Validate: 検証 {
  Check: サービスが存在するか？ {
    shape: diamond
  }

  CheckMethod: メソッドが存在するか？ {
    shape: diamond
  }

  CheckTypes: 型が正しいか？ {
    shape: diamond
  }
}

Invoke: Go メソッドを呼び出す {
  shape: rectangle
  style.fill: "#00ADD8"
}

Encode: 結果をエンコード {
  shape: rectangle
}

Send: レスポンスを送信 {
  shape: rectangle
  style.fill: "#10B981"
}

Error: エラーを送信 {
  shape: rectangle
  style.fill: "#EF4444"
}

Receive -> Parse
Parse -> Validate.Check
Validate.Check -> Validate.CheckMethod: はい
Validate.Check -> Error: いいえ
Validate.CheckMethod -> Validate.CheckTypes: はい
Validate.CheckMethod -> Error: いいえ
Validate.CheckTypes -> Invoke: はい
Validate.CheckTypes -> Error: いいえ
Invoke -> Encode: 成功
Invoke -> Error: エラー
Encode -> Send
```

<strong>セキュリティ：</strong>呼び出せるのは、登録済みのサービスとエクスポートされたメソッドだけです。

### 5. Go での実行（実行時）

Go メソッドが実行されます。

```go
func (g *GreetService) Greet(name string) string {
    // This runs in Go
    return g.prefix + name + "!"
}
```

**実行コンテキスト：**

- **goroutine** で実行（ノンブロッキング）
- **Go のすべての機能**（ファイルシステム、ネットワーク、データベース）を利用可能
- <strong>他の Go コード</strong>を自由に呼び出し可能
- 結果またはエラーを返す

### 6. レスポンス（実行時）

結果が JavaScript に返されます。

```javascript
// Promise resolves with result
const greeting = await Greet("World")
// greeting = "Hello, World!"
```

**エラー処理：**

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

## パフォーマンス特性

### 速度

**一般的な呼び出しオーバーヘッド：** &lt;1ms

```
Frontend Call → Bridge → Go Execution → Bridge → Frontend Response
     ↓            ↓           ↓            ↓            ↓
   &lt;0.1ms      &lt;0.1ms      [varies]     &lt;0.1ms      &lt;0.1ms
```

**他の方式との比較：**

- **HTTP/REST：** 5-50ms（ネットワークスタック、シリアル化）
- **IPC：** 1-10ms（プロセス境界、マーシャリング）
- **Wails Bridge：** &lt;1ms（インメモリ、直接呼び出し）

### メモリ

**呼び出しごとのオーバーヘッド：** 約1KB（メッセージバッファ）

<strong>ゼロコピー最適化：</strong>大容量データ（>1MB）では、可能な場合に共有メモリを使用します。

### 並行処理

**呼び出しは並行して実行されます：**

- 各呼び出しは個別の goroutine で実行される
- 複数の呼び出しを同時に実行可能
- 呼び出し間でブロッキングが発生しない

```javascript
// These run concurrently
const [result1, result2, result3] = await Promise.all([
    SlowOperation1(),
    SlowOperation2(),
    SlowOperation3(),
])
```

## 型システム

### サポートされる型

#### プリミティブ型

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

#### スライスと配列

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

#### マップ

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

#### 構造体

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

<strong>JSON タグ：</strong>TypeScript でのフィールド名を制御するには、`json:` タグを使用します。

#### 時刻

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

#### エラー

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

### サポートされない型

次の型はブリッジ経由で渡すことが<strong>できません</strong>。

- **チャネル**（`chan T`）
- **関数**（`func()`）
- **インターフェース**（`interface{}` / `any` を除く）
- **ポインター**（構造体へのポインターを除く）
- **非公開フィールド**（小文字で始まるフィールド）

<strong>回避策：</strong>ID またはハンドルを使用します：

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

## 高度なパターン

### コンテキストの受け渡し

サービスは呼び出しコンテキストにアクセスできます：

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

**コンテキストから取得できる情報：**

- 呼び出し元のウィンドウ
- アプリケーションインスタンス
- リクエストのメタデータ

### データのストリーミング

大量のデータには、戻り値ではなくイベントを使用します：

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

### キャンセル

キャンセル可能な処理にはコンテキストを使用します：

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

<strong>注：</strong>フロントエンドの切断時には、コンテキストが自動的にキャンセルされます。

### バッチ処理

処理をバッチ化してブリッジのオーバーヘッドを削減します：

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

## ブリッジのデバッグ

### デバッグログの有効化

```go
app := application.New(application.Options{
    Name:     "My App",
    LogLevel: slog.LevelDebug, // requires `import "log/slog"`
    // `Logger` is an optional *slog.Logger; the default-logger helper is
    // application.DefaultLogger(slog.Leveler) if you want to construct one explicitly.
})
```

**出力に表示される情報：**

- メソッド呼び出し
- 引数
- 戻り値
- エラー
- タイミング情報

### 生成されたバインディングの確認

生成された TypeScript を確認するには、`frontend/bindings/`を調べます：

```javascript
// frontend/bindings/<full-go-import-path>/myservice.js (real generated shape)
import { Call as $Call } from "/wails/runtime.js";

export function MyMethod($0) {
    return $Call.ByID(1234567890, $0); // numeric method ID assigned by the generator
}
```

### サービスの直接テスト

フロントエンドを使用せずに Go サービスをテストします：

```go
func TestGreetService(t *testing.T) {
    service := &GreetService{prefix: "Hello, "}
    result := service.Greet("Test")
    if result != "Hello, Test!" {
        t.Errorf("Expected 'Hello, Test!', got '%s'", result)
    }
}
```

## パフォーマンスのヒント

### ✅ 推奨事項

- **処理をバッチ化する** — ブリッジの呼び出し回数を減らす
- **ストリーミングにはイベントを使用する** — 大きな配列を返さない
- **メソッドを高速に保つ** — 理想は&lt;100ms
- **goroutine を使用する** — 長時間の処理に使用する
- **Go 側でキャッシュする** — 計算の繰り返しを避ける

### ❌ 禁止事項

- **過剰に呼び出さない** — 可能な場合はバッチ化する
- **巨大なデータを返さない** — ページネーションまたはストリーミングを使用する
- **ブロックしない** — 長時間の処理には goroutine を使用する
- **複雑な型を渡さない** — 単純な型を使用する
- **エラーを無視しない** — 必ず処理する

## セキュリティ

ブリッジはデフォルトで安全です：

1. **ホワイトリストのみ** — 呼び出せるのは登録済みのサービスだけです
2. **型の検証** — 引数が Go の型に照らして検証されます
3. **eval() 不使用** — フロントエンドから任意の Go コードを実行することはできません
4. **リフレクションの悪用を防止** — アクセスできるのは公開メソッドだけです

**ベストプラクティス：**

- Go で<strong>入力を検証する</strong>（フロントエンドを信頼しない）
- 認証／認可には<strong>コンテキストを使用する</strong>
- 高コストな処理を<strong>レート制限する</strong>
- ファイルパスとユーザー入力を<strong>サニタイズする</strong>

## 次のステップ

**ビルドシステム** — Wails がアプリケーションをビルドしてバンドルする仕組みを学びます [詳細を見る →](/concepts/build-system/)

**サービス** — サービスシステムを詳しく掘り下げます [詳細を見る →](/features/bindings/services/)

**イベント** - パブリッシュ／サブスクライブ通信にイベントを使用します [詳細を見る →](/features/events/system/)

---

**ブリッジについて質問がありますか？** [Discord](https://discord.gg/JDdSxwjhGf)で質問するか、[バインディングの例](https://github.com/wailsapp/wails/tree/master/v3/examples/binding)を確認してください。
