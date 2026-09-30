---
title: "コーディング規約"
description: "Wails v3 のコードスタイル、規約、ベストプラクティス"
slug: "contributing/standards"
sourcePath: "contributing/standards.md"
---

## コードスタイルと規約

一貫したコーディング規約に従うことで、コードベースが読みやすく、保守や貢献もしやすくなります。

## Go コード規約

### コードのフォーマット

標準の Go フォーマットツールを使用します。

```bash
# Format all code
gofmt -w .

# Use goimports for import organization
goimports -w .
```

<strong>必須：</strong>コミットする前に、すべての Go コードが `gofmt` と `goimports` のチェックに合格している必要があります。

### 命名規則

**パッケージ：**

- 小文字を使用し、可能であれば単語を 1 つにします
- `package application`、`package events`
- アンダースコアや大文字と小文字の混在は避けます

**エクスポートされる名前：**

- 型、関数、定数には PascalCase を使用します
- `type WebviewWindow struct`、`func NewApplication()`

**エクスポートされない名前：**

- 内部の型、関数、変数には camelCase を使用します
- `type windowImpl struct`、`func createWindow()`

**インターフェース：**

- 振る舞いに基づいて命名します：`Reader`、`Writer`、`Handler`
- メソッドが 1 つだけのインターフェース：名前に接尾辞 `-er` を付けます

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

### エラー処理

**エラーを必ず確認する：**

```go
// Good
result, err := doSomething()
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// Bad - ignoring errors
result, _ := doSomething()
```

**エラーラッピングを使用する：**

```go
// Wrap errors to provide context
if err := validate(); err != nil {
    return fmt.Errorf("validation failed: %w", err)
}
```

**必要に応じてカスタムエラー型を作成する：**

```go
type ValidationError struct {
    Field string
    Value string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("invalid value %q for field %q", e.Value, e.Field)
}
```

### コメントとドキュメント

**パッケージコメント：**

```go
// Package application provides the core Wails application runtime.
//
// It handles window management, event dispatching, and service lifecycle.
package application
```

**エクスポートされる宣言：**

```go
// NewApplication creates a new Wails application with the given options.
//
// The application must be started with Run() or RunWithContext().
func NewApplication(opts Options) *Application {
    // ...
}
```

**実装コメント：**

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

### 関数とメソッドの構造

**関数の責務を絞る：**

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

**早期リターンを使用する：**

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

### 並行処理

**キャンセルには context を使用する：**

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

**共有状態をミューテックスで保護する：**

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

**ゴルーチンのリークを防ぐ：**

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

### テスト

**テストファイルの命名：**

```go
// Implementation: window.go
// Tests: window_test.go
```

**テーブル駆動テスト：**

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

## JavaScript／TypeScript 規約

### コードのフォーマット

一貫したフォーマットには Prettier を使用します。

```json
{
  "semi": false,
  "singleQuote": true,
  "tabWidth": 2,
  "trailingComma": "es5"
}
```

### 命名規則

**変数と関数：**

- camelCase：`const userName = "John"`

**クラスと型：**

- PascalCase：`class WindowManager`

**定数：**

- UPPER<em>SNAKE</em>CASE：`const MAX_RETRIES = 3`

### TypeScript

**明示的な型を使用する：**

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

**インターフェースを定義する：**

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

## コミットメッセージの形式

[Conventional Commits](https://www.conventionalcommits.org/) を使用します。

```
<type>(<scope>): <subject>

<body>

<footer>
```

**種類：**

- `feat`：新機能
- `fix`：バグ修正
- `docs`：ドキュメントの変更
- `refactor`：コードのリファクタリング
- `test`：テストの追加または更新
- `chore`：メンテナンス作業

**例：**

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

## プルリクエストのガイドライン

### 提出前の確認

- [ ] コードが`gofmt`と`goimports`に合格している
- [ ] すべてのテスト（`go test ./...`）に合格している
- [ ] 新しいコードにテストがある
- [ ] 必要に応じてドキュメントを更新している
- [ ] コミットメッセージが規約に従っている
- [ ] `master`とのマージ競合がない

### PR の説明テンプレート

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

## コードレビュープロセス

### レビュー担当者として

- 建設的かつ敬意を持って対応する
- 個人的な好みではなく、コードの品質を重視する
- 変更を提案する理由を説明する
- 内容に納得したら承認する

### 作成者として

- すべてのコメントに返信する
- 必要に応じて説明を求める
- 要求された変更を行うか、行わない理由を説明する
- フィードバックを柔軟に受け入れる

## ベストプラクティス

### パフォーマンス

- 早すぎる最適化を避ける
- 最適化する前にプロファイリングを行う
- パフォーマンスが重要なコードにはベンチマークを使用する

```go
func BenchmarkProcess(b *testing.B) {
    for i := 0; i < b.N; i++ {
        process(testData)
    }
}
```

### セキュリティ

- すべてのユーザー入力を検証する
- データを表示する前にサニタイズする
- ランダムデータには`crypto/rand`を使用する
- 機密情報を決してログに記録しない

### ドキュメント

- エクスポートされた API を文書化する
- ドキュメントに例を含める
- API を変更したらドキュメントを更新する
- README ファイルを最新の状態に保つ

## プラットフォーム固有のコード

### ファイル名

```
window.go           // Common interface
window_darwin.go    // macOS implementation
window_windows.go   // Windows implementation
window_linux.go     // Linux implementation
```

### ビルドタグ

```go
//go:build darwin

package application

// macOS-specific code
```

## Lint

コミットする前に Linter を実行します：

```bash
# golangci-lint (recommended)
golangci-lint run

# Individual linters
go vet ./...
staticcheck ./...
```

## 質問がある場合

いずれかの規約について不明な点がある場合：

- 既存のコードで例を確認する
- [Discord](https://discord.gg/JDdSxwjhGf)で質問する
- GitHub でディスカッションを開始する
