---
title: "v2 から v3 への移行"
description: "Wails v2 アプリケーションを v3 に移行するための完全ガイド"
slug: "migration/v2-to-v3"
sourcePath: "migration/v2-to-v3.md"
---

Wails v3 は、アーキテクチャ、パフォーマンス、開発者体験が大幅に改善された**全面的な書き直し**です。このガイドでは、v2 アプリケーションを v3 に移行する方法を説明します。

**主な変更点：**

- 新しいアプリケーション構造
- 改善されたバインディングシステム
- 強化されたウィンドウ管理
- 改善されたイベントシステム
- 簡素化された設定

<strong>移行時間：</strong>一般的なアプリケーションでは 1-4 時間

## 破壊的変更

### アプリケーションの初期化

v2 では、アプリケーションのセットアップ、ウィンドウの設定、実行がすべて単一の `wails.Run()` 呼び出しにまとめられていました。このモノリシックな方式では、複数のウィンドウを作成したり、各段階でエラーを処理したり、アプリケーションの個々のコンポーネントをテストしたりすることが困難でした。

v3 では、これらの関心事がアプリケーションの作成、ウィンドウの作成、実行という個別の段階に分離されています。この分離により、アプリケーションのライフサイクルの各段階を明示的に制御でき、コードのモジュール性とテスト容易性が向上します。

**v2：**

```go
err := wails.Run(&options.App{
    Title:  "My App",
    Width:  1024,
    Height: 768,
    Bind: []interface{}{
        &GreetService{},
    },
})
```

**v3：**

```go
app := application.New(application.Options{
    Name: "My App",
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
})

window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "My App",
    Width:  1024,
    Height: 768,
})

app.Run()
```

**この方式が優れている理由：**

- **マルチウィンドウのサポート**：起動時だけでなく、いつでも動的にウィンドウを作成できます
- **エラー処理の改善**：適切なエラー処理を使用して、各段階を個別に検証できます
- **より明確なコード**：処理が分離されているため、各段階で何が起きているかが明確になります
- **テスト容易性の向上**：イベントループを実行せずにアプリケーションのセットアップをテストできます
- **柔軟性の向上**：アプリケーションのライフサイクル全体を通じて、ウィンドウを作成、破棄、再作成できます

### バインディング

v2 では、バインドされたすべての構造体にコンテキストフィールドと、ランタイムコンテキストを受け取るための `startup(ctx)` メソッドが必要でした。このため、ビジネスロジックと Wails ランタイムが密結合になり、コードのテストや理解が難しくなっていました。

v3 ではサービスパターンが導入され、構造体は完全に独立し、ランタイムコンテキストを保持する必要がなくなりました。サービスがアプリケーションインスタンスにアクセスする必要がある場合は、暗黙的にコンテキストを引き回すのではなく、依存性注入を通じて明示的に受け取ります。

**v2：**

```go
type App struct {
    ctx context.Context
}

func (a *App) startup(ctx context.Context) {
    a.ctx = ctx
}

func (a *App) Greet(name string) string {
    return "Hello " + name
}
```

**v3：**

```go
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello " + name
}

// Register as service
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
})
```

**この方式が優れている理由：**

- **暗黙的な依存関係の排除**：サービスは、隠れたランタイム依存関係を持たない通常の Go 構造体です
- **テストの簡素化**：Wails コンテキストをモックせずにサービスメソッドをテストできます
- **より明確なコード**：依存関係はコンテキストフィールドに隠されず、コンストラクター引数として明示的に渡されます
- **構成の改善**：すべてのサービスを単一の `App` 構造体に配置するのではなく、ドメイン別にグループ化できます
- **適切な初期化**：初期化が必要な場合は `ServiceStartup()` メソッドを使用し、その処理を明示します

### ランタイム

v2 では、すべてのランタイム操作で、`runtime` パッケージのグローバル関数にコンテキストを渡す必要がありました。このため、コードベース全体がコンテキストオブジェクトに密結合となり、API はオブジェクト指向ではなく手続き型のように感じられました。

v3 では、コンテキストベースのランタイムが、アプリケーションオブジェクトとウィンドウオブジェクトに対する直接のメソッド呼び出しに置き換えられています。操作の対象となるオブジェクト上でメソッドを直接呼び出すため、コードがより直感的でオブジェクト指向になります。

**v2：**

```go
import "github.com/wailsapp/wails/v2/pkg/runtime"

runtime.WindowSetTitle(a.ctx, "New Title")
runtime.EventsEmit(a.ctx, "event-name", data)
```

**v3：**

```go
// Store app reference
type MyService struct {
    app *application.App
}

func (s *MyService) UpdateTitle() {
    window := s.app.Window.Current()
    window.SetTitle("New Title")
}

func (s *MyService) EmitEvent() {
    s.app.Event.Emit("event-name", data)
}
```

**この方式が優れている理由：**

- **オブジェクト指向設計**：操作対象のオブジェクト（ウィンドウ、アプリ、メニューなど）上でメソッドを呼び出します
- **より明確な意図**：`window.SetTitle()` は `runtime.WindowSetTitle(ctx, ...)` よりも意図が明確です
- **IDE サポートの改善**：メソッドがオブジェクト上にあるため、オートコンプリートが適切に機能します
- **マルチウィンドウ操作の明確化**：複数のウィンドウがある場合、操作対象のウィンドウを明示的に選択します
- **コンテキストの引き回しが不要**：すべての関数にコンテキストを渡す必要がありません

### フロントエンドバインディング

v2 では、バインディングは Go パッケージ名と構造体名に基づいて構成され、通常は `wailsjs/go/main/App` のようなパスになっていました。この構造は論理的なグループ分けを反映しておらず、関連する機能を見つけることが困難でした。

v3 では、バインディングがサービス名とアプリケーションモジュール別に整理され、より明確で論理的な構造になります。バインディングは、アプリケーション名とサービス名で構成された `bindings` ディレクトリに生成されるため、利用可能な機能を把握しやすくなります。

**v2：**

```javascript
import { Greet } from '../wailsjs/go/main/App'

const result = await Greet("World")
```

**v3：**

```javascript
import { Greet } from './bindings/changeme/greetservice'

const result = await Greet("World")
```

**この方式が優れている理由：**

- **論理的な構成**：バインディングは Go パッケージ構造ではなく、サービス名別にグループ化されます
- **より明確なインポート**：パスにはファイル構造（main/App）ではなく、ドメインロジック（greetservice）が反映されます
- **見つけやすさの向上**：技術的な構造ではなく、機能別にバインディングをたどれます
- **一貫した命名**：サービスベースの構成がバックエンドアーキテクチャと一致します
- **より単純なパス**：`../wailsjs/go` プレフィックスは不要になり、`./bindings` だけになります

### イベント

v2 では、イベントに可変長の `interface{}` パラメーターが使用され、すべてのイベント関数にコンテキストを渡す必要がありました。イベントハンドラーが受け取るデータには型がなく、手動で型アサーションを行う必要があったため、イベントシステムではエラーが発生しやすく、デバッグも困難でした。

v3 では、型付きイベントオブジェクトが導入され、コンテキストが不要になります。イベントハンドラーは型付きデータを含む適切なイベントオブジェクトを受け取るため、イベントシステムの信頼性と使いやすさが向上します。

**v2：**

```go
runtime.EventsOn(ctx, "event-name", func(data ...interface{}) {
    // Handle event
})

runtime.EventsEmit(ctx, "event-name", data)
```

**v3：**

```go
app.Event.On("event-name", func(e *application.CustomEvent) {
    data := e.Data
    // Handle event
})

app.Event.Emit("event-name", data)
```

**この方式が優れている理由：**

- **型安全性**：イベントでは、`...interface{}`の代わりに適切なイベントオブジェクトを使用します
- **デバッグの向上**：イベントオブジェクトにはイベント名などのメタデータが含まれるため、デバッグが容易になります
- **より明確な API**：`app.Event.On()`と`app.Event.Emit()`は、runtime 関数より直感的です
- **コンテキストが不要**：イベントは、コンテキストをスレッド間で受け渡さなくても app オブジェクト上で直接機能します
- **よりシンプルなハンドラー**：イベントハンドラーでは、可変長引数の代わりに明確なシグネチャを使用します

### ウィンドウ

v2では、アプリケーションごとにウィンドウを 1 つだけ使用できました。ウィンドウは起動時に作成され、すべてのウィンドウ操作は、その単一のウィンドウを暗黙的に対象とする runtime 関数を通じて実行されていました。

v3では、ネイティブのマルチウィンドウ対応が中核機能として導入されています。各ウィンドウは、固有のメソッドとライフサイクルを持つ第一級オブジェクトです。アプリケーションの実行中、複数のウィンドウを動的に作成、管理、破棄できます。

**v2：**

```go
// Single window only
runtime.WindowSetSize(ctx, 800, 600)
```

**v3：**

```go
// Multiple windows supported
window1 := app.Window.New()
window1.SetSize(800, 600)

window2 := app.Window.New()
window2.SetSize(1024, 768)
```

**改善点：**

- **マルチウィンドウアプリケーション**：独立した複数のウィンドウを持つアプリ（ダッシュボード、環境設定、ツールなど）を構築できます
- **明示的なウィンドウ参照**：各ウィンドウは、保存して直接操作できるオブジェクトです
- **動的なウィンドウ作成**：実行中にいつでもウィンドウを作成、破棄できます
- **独立したウィンドウ状態**：各ウィンドウは、固有のイベント、プロパティ、ライフサイクルを持ちます
- **より優れたアーキテクチャ**：ウィンドウ管理はコンテキストベースではなく、オブジェクト指向になっています

## 移行手順

### ステップ 1：依存関係を更新する

**go.mod：**

```go
module myapp

go 1.25.0

require (
    github.com/wailsapp/wails/v3 v3.0.0
)
```

**更新後：**

```bash
go get github.com/wailsapp/wails/v3@v3.0.0
go mod tidy
```

### ステップ 2：main.go を更新する

**v2：**

```go
package main

import (
    "embed"
    "github.com/wailsapp/wails/v2/pkg/options"
    "github.com/wailsapp/wails/v2/pkg/options/assetserver"
    "github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
    app := NewApp()

    err := wails.Run(&options.App{
        Title:  "My App",
        Width:  1024,
        Height: 768,
        AssetServer: &assetserver.Options{
            Assets: assets,
        },
        Bind: []interface{}{
            app,
        },
        Windows: &windows.Options{
            WebviewIsTransparent: false,
        },
    })

    if err != nil {
        println("Error:", err.Error())
    }
}
```

**v3：**

```go
package main

import (
    "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name: "My App",
        Services: []application.Service{
            application.NewService(&MyService{}),
        },
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:  "My App",
        Width:  1024,
        Height: 768,
    })

    err := app.Run()
    if err != nil {
        panic(err)
    }
}
```

### ステップ 3：App 構造体をサービスに変換する

**v2：**

```go
type App struct {
    ctx context.Context
}

func NewApp() *App {
    return &App{}
}

func (a *App) startup(ctx context.Context) {
    a.ctx = ctx
    // Initialisation
}

func (a *App) Greet(name string) string {
    return "Hello " + name
}
```

**v3：**

```go
type MyService struct {
    app *application.App
}

func NewMyService(app *application.App) *MyService {
    return &MyService{app: app}
}

func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    // Initialisation
    return nil
}

func (s *MyService) Greet(name string) string {
    return "Hello " + name
}

// Register after app creation
app := application.New(application.Options{})
app.RegisterService(application.NewService(NewMyService(app)))
```

### ステップ 4：Runtime 呼び出しを更新する

**v2：**

```go
func (a *App) DoSomething() {
    runtime.WindowSetTitle(a.ctx, "New Title")
    runtime.EventsEmit(a.ctx, "update", data)
    runtime.LogInfo(a.ctx, "Message")
}
```

**v3：**

```go
func (s *MyService) DoSomething() {
    window := s.app.Window.Current()
    window.SetTitle("New Title")
    
    s.app.Event.Emit("update", data)
    
    s.app.Logger.Info("Message")
}
```

### ステップ 5：フロントエンドを更新する

**新しいバインディングを生成します：**

```bash
wails3 generate bindings
```

**import を更新します：**

```javascript
// v2
import { Greet } from '../wailsjs/go/main/App'

// v3
import { Greet } from './bindings/changeme/myservice'
```

**イベント処理を更新します：**

```javascript
// v2
import { EventsOn, EventsEmit } from '../wailsjs/runtime/runtime'

EventsOn("update", (data) => {
    console.log(data)
})

EventsEmit("action", data)

// v3
import { Events } from '@wailsio/runtime'

Events.On("update", (data) => {
    console.log(data)
})

Events.Emit("action", data)
```

### ステップ 6：設定を更新する

**v2（wails.json）：**

```json
{
  "name": "myapp",
  "outputfilename": "myapp",
  "frontend:install": "npm install",
  "frontend:build": "npm run build",
  "frontend:dev:watcher": "npm run dev",
  "frontend:dev:serverUrl": "auto"
}
```

**v3（wails.json）：**

```json
{
  "name": "myapp",
  "frontend": {
    "dir": "./frontend",
    "install": "npm install",
    "build": "npm run build",
    "dev": "npm run dev",
    "devServerUrl": "http://localhost:5173"
  }
}
```

## 機能の対応関係

### ダイアログ

**v2：**

```go
selection, err := runtime.OpenFileDialog(ctx, runtime.OpenDialogOptions{
    Title: "Select File",
})
```

**v3：**

```go
selection, err := app.Dialog.OpenFileWithOptions(&application.OpenFileDialogOptions{
    Title: "Select File",
}).PromptForSingleSelection()
```

### メニュー

**v2：**

```go
menu := menu.NewMenu()
menu.Append(menu.Text("File", nil, []*menu.MenuItem{
    menu.Text("Quit", nil, func(_ *menu.CallbackData) {
        runtime.Quit(ctx)
    }),
}))
```

**v3：**

```go
menu := app.NewMenu()
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})
```

### システムトレイ

**v2：**

```go
// Not available in v2
```

**v3：**

```go
systray := app.SystemTray.New()
systray.SetIcon(iconBytes)
systray.SetLabel("My App")

menu := app.NewMenu()
menu.Add("Show").OnClick(showWindow)
menu.Add("Quit").OnClick(app.Quit)
systray.SetMenu(menu)
```

## よくある問題

### 問題：バインディングが見つからない

<strong>問題：</strong>移行後に import エラーが発生する

**解決策：**

```bash
# Regenerate bindings
wails3 generate bindings

# Check output directory
ls frontend/bindings
```

### 問題：コンテキストエラー

**問題：**`ctx`を使用できない

**解決策：**

代わりに app への参照を保存します：

```go
type MyService struct {
    app *application.App
}

func NewMyService(app *application.App) *MyService {
    return &MyService{app: app}
}
```

### 問題：ウィンドウメソッドが機能しない

**問題：**`runtime.WindowSetTitle()`が存在しない

**解決策：**

ウィンドウメソッドを直接使用します：

```go
window := s.app.Window.Current()
window.SetTitle("New Title")
```

### 問題：イベントが発生しない

<strong>問題：</strong>イベントは登録されているが受信されない

**解決策：**

イベント名が完全に一致していることを確認します：

```go
// Go
app.Event.Emit("my-event", data)

// JavaScript
OnEvent("my-event", handler)  // Must match exactly
```

## 移行後のテスト

### チェックリスト

- [ ] アプリケーションがエラーなく起動する
- [ ] すべてのバインディングが機能する
- [ ] イベントが送受信される
- [ ] ウィンドウが正しく開閉する
- [ ] メニューが機能する（該当する場合）
- [ ] ダイアログが機能する（該当する場合）
- [ ] システムトレイが機能する（該当する場合）
- [ ] ビルドプロセスが機能する
- [ ] 本番用ビルドが機能する

### テストコマンド

```bash
# Development
wails3 dev

# Build
wails3 build

# Generate bindings
wails3 generate bindings
```

## v3の利点

### パフォーマンス

- **起動の高速化** - 初期化を最適化
- **メモリ使用量の削減** - リソースを効率的に使用
- **ブリッジの改善** - 呼び出しのオーバーヘッドは&lt;1ms

### 機能

- **マルチウィンドウ** - ネイティブ対応
- **システムトレイ** - 組み込み
- **改善されたイベント** - 型付きの、よりシンプルな API
- **サービス** - コード構成を改善

### 開発者エクスペリエンス

- **型安全性** - TypeScript を完全にサポート
- **改善されたエラー** - 明確なエラーメッセージ
- **ホットリロード** - 開発を高速化
- **改善されたドキュメント** - 網羅的なガイド

## ヘルプを得る

### リソース

- [ドキュメント](/quick-start/why-wails/)
- [Discord コミュニティ](https://discord.gg/JDdSxwjhGf)
- [GitHub Issues](https://github.com/wailsapp/wails/issues)
- [サンプル](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples)

### よくある質問

**Q: v2 と v3を併用できますか？** A: はい。異なるインポートパスを使用するため、併用できます。

**Q: v3は本番環境で使用できますか？** Wails v3.0.0は、Windows、macOS、Linux 向けデスクトップアプリケーションの安定版です。Android と iOS のサポートは引き続き実験的であり、安定版デスクトップの互換性保証の対象外です。

**Q: v2のメンテナンスは継続されますか？** A: はい。v2には重要なアップデートが提供されます。

**Q: 移行にはどのくらい時間がかかりますか？** A: 一般的なアプリケーションでは 1-4時間です。

## 次のステップ

@cards{cols="2"}
🚀 クイックスタート
Wails v3を使い始めましょう。

[詳細を見る →](/quick-start/installation/)

---
★ コアコンセプト
v3のアーキテクチャを理解しましょう。

[詳細を見る →](/concepts/architecture/)

---
◆ バインディング
新しいバインディングシステムを学びましょう。

[詳細を見る →](/features/bindings/methods/)

---
📖 サンプル
v3の完全なサンプルを確認しましょう。

[サンプルを見る →](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples)

@end

---

**質問がありますか？** [Discord](https://discord.gg/JDdSxwjhGf) で質問するか、[Issue を作成](https://github.com/wailsapp/wails/issues)してください。
