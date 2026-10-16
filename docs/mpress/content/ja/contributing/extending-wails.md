---
title: "Wails の拡張"
description: "Wails v3 に新機能や新しいプラットフォームを追加するための実践ガイド"
slug: "contributing/extending-wails"
sourcePath: "contributing/extending-wails.md"
---

> Wails は<strong>自由に改造できる</strong>ように設計されています。
>
> 主要なサブシステムはすべて、読んで変更し、リリースできる Go コードで実装されています。
>
> このページでは、次の作業を行う際に、<em>どこ</em>から着手し、<em>どのように</em>クロスプラットフォーム対応を維持するかを説明します。

- <strong>サービス</strong>を追加する（通知、KV ストア、カスタム IPC など）
- <strong>新しい CLI コマンド</strong>を作成する（`wails3 <foo>`）
- <strong>ランタイム</strong>を拡張する（ウィンドウ API、ダイアログ、イベント）
- <strong>プラットフォーム機能</strong>を導入する（Wayland など）
- `//go:build`タグに埋もれることなく、<strong>クロスプラットフォーム互換性</strong>を維持する

---

## 1. サービスの追加

v3 における「サービス」とは、ユーザーが用意し、 `application.Options.Services`を通じて登録して、生成されたバインディング経由で JS に公開する Go の型です。 v3 のコードベースには次のものが含まれています。

- `internal/service/` — `wails3 generate service`用のスキャフォールディング：
  ```
  internal/service/
  ├── service.go              # Install(options *flags.ServiceInit)
  └── template/
      ├── README.tmpl.md
      ├── go.mod.tmpl
      ├── service.go.tmpl
      └── service.tmpl.yml
  ```

- `pkg/services/` — 現在すぐに登録して使用できる既製のサービス（notifications、kvstore、sqlite、log、fileserver、dock など）。

古い草稿で参照されていたジェネレーターおよび CLI ファイルの `internal/service/template/template.go`と `internal/generator/collect/services.go`は存在しません。スキャフォールダーは `internal/service/service.go`（エントリーポイントは `service.Install`）であり、サービスのバインディング メタデータは `internal/generator/collect/service.go` で収集されます。

### 1.1 サービスを定義する

```go
package chat

type Service struct {
    messages []string
}

func New() *Service { return &Service{} }

func (s *Service) Send(msg string) string {
    s.messages = append(s.messages, msg)
    return "ok"
}
```

### 1.2 ライフサイクルインターフェースを実装する（任意）

サービスでは、必要に応じて次のインターフェース（`pkg/application` で定義）を実装できます。

```go
func (s *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error { return nil }
func (s *Service) ServiceShutdown() error                                                       { return nil }
```

> **重要：**`ServiceShutdown`は<strong>引数を取りません</strong>。次の
>
> シグネチャ `ServiceShutdown(ctx context.Context) error` のメソッドは、このインターフェースを<strong>実装したことにはなりません</strong>。そのため、
>
> 何の通知もなく、呼び出されることはありません。

### 1.3 アプリケーションにサービスを登録する

グローバルな `services.Register(...)` 呼び出しはありません。サービスは実行時に `application.Options.Services`を介して登録します。

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(chat.New()),
    },
})
```

登録すると、`wails3 generate bindings`は、エクスポートされたメソッドをラップする ES モジュールを `frontend/bindings/<your import path>/...`配下に出力します。

### 1.4 JS から呼び出す

```js
import { Send } from "../bindings/github.com/you/yourapp/chat";

await Send("hi");
```

v3 にはグローバルな `window.backend.*` はありません。呼び出しは生成された ES モジュールを経由し、そのモジュールが `/wails/runtime.js` の `Call.ByID(...)` を呼び出します。

---

## 2. 新しい CLI コマンドの作成

v3 の CLI は cobra ではなく **`github.com/leaanthony/clir`** を使用します。配線処理は `v3/cmd/wails3/main.go`にあります。

```go
import "github.com/leaanthony/clir"

func main() {
    app := clir.NewCli("wails", "The Wails3 CLI", "v3")
    app.NewSubCommand("hello", "Prints Hello Wails").Action(func() error {
        fmt.Println("Hello Wails!")
        return nil
    })
    // ... other subcommands explicitly wired here
    _ = app.Run()
}
```

`init()`を利用した自動登録はありません。新しいサブコマンドを `cmd/wails3/main.go`に追加し、対応する関数を `internal/commands/` に追加してください（オプションを取る場合は、 `internal/flags/`配下にフラグ用の構造体も追加します）。CLI を再ビルドします。

```
cd v3
go install ./cmd/wails3
wails3 hello
```

コマンドに Taskfile との連携処理が必要な場合は、 `internal/commands/task_wrapper.go`（`wrapTask("yourtask", args)`）のヘルパーを再利用してください。

---

## 3. ランタイムの変更

一般的な変更理由：

- 新しいウィンドウ機能（`SetOpacity`、`Shake` など）
- 追加のダイアログ（`ColorPicker`）
- システムレベルの API（画面の明るさ）

### 3.1 公開 API

`pkg/application/webview_window.go`にメソッドを追加します（インターフェースは `window.go`にあります）。

```go
func (w *WebviewWindow) SetOpacity(o float32) Window {
    InvokeSync(func() { w.impl.setOpacity(o) })
    return w
}
```

呼び出しがメインスレッドで実行されるように、既存の `InvokeSync`／`InvokeAsync` ヘルパーを使用してください。

### 3.2 メッセージプロセッサー

JS から新しいメソッドを呼び出す必要がある場合は、該当する `pkg/application/messageprocessor_*.go`ファイルを拡張してください。メッセージプロセッサーでは、グローバルな `register(...)` 呼び出しではなく、 `MessageProcessor`の switch ベースのメソッドを使用します。

```go
// inside messageprocessor_window.go
case "setOpacity":
    var args struct {
        WindowID uint    `json:"windowID"`
        Opacity  float32 `json:"opacity"`
    }
    if err := json.Unmarshal(payload, &args); err != nil { ... }
    window, _ := m.app.Window.GetByID(args.WindowID)
    window.SetOpacity(args.Opacity)
```

コードベースには `messageprocessor_window_opacity.go` ファイルも、`init()`を利用した `register(MsgSetOpacity, ...)`パターンも<strong>存在しません</strong>。

### 3.3 プラットフォーム別の実装

`pkg/application/`配下にある OS 別の各ファイルに実装を追加してください。

```
pkg/application/
├── webview_window_darwin.go   //go:build darwin
├── webview_window_linux.go    //go:build linux
└── webview_window_windows.go  //go:build windows
```

あるプラットフォームでその機能をサポートできない場合は、何もしないスタブを作成してください。フレームワークには `ErrCapability`センチネルがありません。サポート状況はドキュメントで明示し、必要であれば、 `Options`またはプラットフォーム固有のオプション構造体の 該当するブール型フィールドを通じて公開してください。

### 3.4 機能フラグ（任意）

`internal/capabilities/`パッケージは、プラットフォーム別の 機能セットを宣言するために存在します。公開された `application.HasCapability`／ `application.CapOpacity` API はありません。実行時に確認可能な機能を追加する場合は、 `internal/capabilities/`配下に追加し、`pkg/application`から型付き getter を 公開してください。

---

## 4. 新しいプラットフォーム機能の追加

例：Linux でオプションの Wayland サポートを追加する場合。

1. 該当する `pkg/application/*_linux.go` ファイルを、`*_linux_x11.go`（`//go:build linux && !wayland`）と `*_linux_wayland.go`（`//go:build linux && wayland`）に分割してください。
2. ユーザーが `wails3 build --tags wayland` でオプトインできるようにしてください。追加のタグは、`internal/commands/task_wrapper.go` にある既存の `EXTRA_TAGS` の配線処理を通じて渡します。`dev`レベルの `--tags wayland` フラグはありません。`wails3 dev`が受け付けるのは、`--config`、`--port`、`-s`だけです。
3. ドキュメントと、`pkg/application/`配下にあるプラットフォーム固有の README を更新してください。

> デフォルトのビルドタグは最小限に留め、オプトイン方式のタグはニッチな機能にのみ使用してください。

---

## 5. クロスプラットフォーム互換性チェックリスト

| ✅ 手順 | 理由 |
| --- | --- |
| すべてのプラットフォーム用ファイルで、公開メソッドを<strong>すべて</strong>用意する（スタブでも可） | すべての OS でビルドが成功する状態を維持できる |
| OS ごとのグレースフルデグラデーションを文書化する | アプリが隠れたエラーなしで`runtime.GOOS`に基づいて処理を分岐できる |
| まず<strong>純粋な Go</strong>を使用し、Cgo は必要な場合にのみ使用する | クロスコンパイルが容易になる（Linux ではすでに Cgo のコストが発生している） |
| `task test:cli`、`task test:generator`、`task test:templates`を実行する | CI をローカルで再現できる |
| 新しいビルドタグをコントリビューター向けドキュメント／テンプレートの README に記載する | ユーザーにオプトインが必要であることを知らせる必要がある |

---

## 6. デバッグビルドと反復速度

- 詳細なランタイム動作を出力するには、`Options.LogLevel = slog.LevelDebug`（`Options.Logger = slog.Default()`）を使用します。`WAILS_LOG_LEVEL`という環境変数はありません。
- `wails3 dev`のフラグは、`--config`、`--port`、`-s`です。`-race`または`-verbose`というフラグはありません。レースディテクターを使用するには、`go test -race ./...`を実行するか、アプリを`go build -race`して直接実行します。
- レース／Cgo テストガイドは`v3/TESTING.md`にあります（古い草稿では、存在しない`pkg/application/RACE.md`を参照していました）。

---

## 7. アップストリームへのコントリビューション

1. 新機能または公開動作の変更については、アイデアと設計を議論するため、<strong>WEP（Wails Enhancement Proposal）</strong>のドラフト PR を作成します。Issue は、再現可能なバグまたはドキュメントの問題にのみ使用してください。
2. 上記の手法に従って実装します。
3. 以下を追加します。
  - 単体テスト（`*_test.go`）
  - ドキュメント（このファイルまたは関連する`docs/...`ページ）
  - バインディングジェネレーターを変更した場合は、`internal/generator/testcases/`配下に回帰テスト

4. プッシュする前に、`task precommit`と関連する`task test:*`ターゲットをローカルで実行します。

---

### クイックリンク

| 領域 | 場所 |
| --- | --- |
| 組み込みサービス | `pkg/services/` |
| サービススキャフォールダー | `internal/service/` |
| CLI の配線 | `v3/cmd/wails3/main.go` |
| CLI コマンド本体 | `internal/commands/` |
| OS ごとのランタイム | `pkg/application/*_{darwin,linux,windows}.go` |
| ケイパビリティ宣言 | `internal/capabilities/` |
| Taskfile DSL | `v3/Taskfile.yaml` |
| イベント定数ジェネレーター | `v3/tasks/events/generate.go` |

---

これで、Wails を思いどおりに拡張するための<strong>ロードマップ</strong>が整いました。サービスを追加する、 CLI に魔法をかける、ランタイムを改造する、あるいはまったく新しい OS 機能を導入することもできます。 拡張を楽しんでください！
