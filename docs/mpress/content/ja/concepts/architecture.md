---
title: "Wails の仕組み"
description: "Wails のアーキテクチャと、ネイティブパフォーマンスを実現する仕組みを理解する"
slug: "concepts/architecture"
sourcePath: "concepts/architecture.md"
---

Wails は、**バックエンドに Go**、<strong>フロントエンドに Web 技術</strong>を使用してデスクトップアプリケーションを構築するためのフレームワークです。ただし Electron とは異なり、ブラウザーを同梱せず、<strong>オペレーティングシステムのネイティブ WebView</strong>を使用します。

```d2
direction: left

Wails App: {
  shape: sequence_diagram
  label: Wails アプリ

  frontend: フロントエンド
  backend: Go バックエンド
  os: オペレーティングシステム

  Initialisation: 初期化 {
    shape: sequence_diagram
    backend."Serves Static Web App": 静的 Web アプリを配信
    backend -> frontend: HTML / JS / CSS
    frontend."Render Site via OS-native WebView": OS ネイティブの WebView でサイトをレンダリング
  }
  Regular Communication: 通常の通信 {
    shape: sequence_diagram
    frontend."Make API-style call": API 形式で呼び出す
    frontend -> backend.a: JSON
    backend.a."Service processes request": サービスがリクエストを処理
    backend.a -> os: システム API を呼び出す
    backend.a."Generate Response": レスポンスを生成
    backend.a -> frontend: JSON
    frontend."Process response": レスポンスを処理
  }
  backend.a.label: a
}
```

**Electron との主な違い：**

| 項目 | Wails | Electron |
| --- | --- | --- |
| **ブラウザー** | OS が提供する WebView | 同梱の Chromium（約 100MB） |
| **バックエンド** | Go（コンパイル済み） | Node.js（インタープリター実行） |
| **通信** | インメモリブリッジ | IPC（プロセス間通信） |
| **バンドルサイズ** | 約 15MB | 約 150MB |
| **メモリ** | 約 10MB | 約 100MB 以上 |
| **起動時間** | &lt;0.5s | 2-3s |

## コアコンポーネント

### 1. ネイティブ WebView

Wails は、オペレーティングシステムに組み込まれた Web レンダリングエンジンを使用します。

@tabs{sync-key="platform"}
[Windows]
**WebView2**（Microsoft Edge WebView2）

- Chromium ベース（Edge ブラウザーと同じ）
- Windows 10/11 にプリインストール済み
- Windows Update による自動更新
- 最新の Web 標準を完全にサポート

[macOS]
**WebKit**（Safari のレンダリングエンジン）

- macOS に組み込み済み
- Safari ブラウザーと同じエンジン
- 優れたパフォーマンスとバッテリー駆動時間
- 最新の Web 標準を完全にサポート

[Linux]
**WebKitGTK**（WebKit の GTK 移植版）

- パッケージマネージャーでインストール
- GNOME Web（Epiphany）と同じエンジン
- Web 標準を十分にサポート
- 軽量で高性能

@end

**これが重要な理由：**

- **ブラウザーを同梱しない** → アプリのサイズを縮小
- **OS ネイティブ** → 統合性とパフォーマンスが向上
- **自動更新** → OS の更新によってセキュリティパッチを適用
- **使い慣れたレンダリング** → システムブラウザーと同じ

### 2. Wails ブリッジ

ブリッジは Wails の中核であり、Go と JavaScript の間の<strong>直接通信</strong>を可能にします。

```d2
direction: down

Frontend: フロントエンド（JavaScript） {
  shape: rectangle
  style.fill: "#8B5CF6"
}

Bridge: Wails ブリッジ {
  Encoder: JSON エンコーダー {
    shape: rectangle
  }

  Router: メソッドルーター {
    shape: diamond
    style.fill: "#10B981"
  }

  Decoder: JSON デコーダー {
    shape: rectangle
  }
}

Backend: バックエンド（Go） {
  Services: 登録済みサービス {
    shape: rectangle
    style.fill: "#00ADD8"
  }
}

Frontend -> Bridge.Encoder: "1. Go メソッドを呼び出す\nGreet('Alice')"
Bridge.Encoder -> Bridge.Router: "2. JSON にエンコード\n{method: 'Greet', args: ['Alice']}"
Bridge.Router -> Backend.Services: "3. サービスにルーティング\nGreetService.Greet('Alice')"
Backend.Services -> Bridge.Decoder: "4. 結果を返す\n'Hello, Alice!'"
Bridge.Decoder -> Frontend: "5. JS にデコード\nPromise が解決される"
```

**仕組み：**

1. **フロントエンドが Go メソッドを呼び出す**（自動生成されたバインディング経由）
2. **ブリッジが呼び出しを** JSON にエンコードする（メソッド名と引数）
3. **ルーターが登録済みサービス内から Go メソッドを見つける**
4. **Go メソッドが実行され**、値を返す
5. **ブリッジが結果をデコードし**、フロントエンドに返す
6. JavaScript で<strong>結果を受け取って Promise が解決される</strong>

**パフォーマンス特性：**

- **インメモリ**：ネットワークのオーバーヘッドも HTTP もなし
- 可能な場合は<strong>ゼロコピー</strong>（大容量データ向け）
- **デフォルトで非同期**：どちら側もブロックしない
- **型安全**：TypeScript の型定義を自動生成

### 3. サービスシステム

フロントエンドに Go の機能を公開するには、サービスの使用を推奨します。

```go
// Define a service (just a regular Go struct)
type GreetService struct {
    prefix string
}

// Methods with exported names are automatically available
func (g *GreetService) Greet(name string) string {
    return g.prefix + name + "!"
}

func (g *GreetService) GetTime() time.Time {
    return time.Now()
}

// Register the service
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&GreetService{prefix: "Hello, "}),
    },
})
```

**サービスの検出：**

- Wails は起動時に<strong>構造体をスキャンします</strong>
- <strong>エクスポートされたメソッド</strong>をフロントエンドから呼び出せるようになります
- TypeScript バインディング用に<strong>型情報</strong>が抽出されます
- <strong>エラー処理</strong>は自動です（Go のエラー → JS の例外）

**生成された TypeScript バインディング：**

```typescript
// Auto-generated in frontend/bindings/GreetService.ts
export function Greet(name: string): Promise<string>
export function GetTime(): Promise<Date>
```

**サービスを使用する理由：**

- **型安全**：TypeScriptを完全にサポート
- **自動検出**：メソッドを手動で登録する必要はありません
- **体系的な整理**：関連する機能をグループ化
- **テスト可能**：サービスは単なるGo構造体です

[サービスの詳細 →](/features/bindings/services/)

### 4. イベントシステム

イベントにより、コンポーネント間で<strong>パブリッシュ／サブスクライブ通信</strong>を行えます。

```d2
direction: left

Wails Event System: Wails イベントシステム {
  shape: sequence_diagram

  window1: ウィンドウ 1
  window2: ウィンドウ 2
  backend: Go バックエンド

  Event Driver: イベントドライバー {
    shape: sequence_diagram
    window1."Subscribe to 'data-updated' events": "'data-updated' イベントを購読"
    window2."Subscribe to 'data-updated' events": "'data-updated' イベントを購読"
    backend.a."App Emit('data-updated', data)": "アプリが Emit('data-updated', data) を実行"
    backend.a -> window1.a: JSON イベントバス
    backend.a -> window2: JSON イベントバス
    window1.a."Subscriber processes On('data-updated', handler)": "購読側が On('data-updated', handler) で処理"
    window2."Subscriber processes On('data-updated', handler)": "購読側が On('data-updated', handler) で処理"
  }
  backend.a.label: a
  window1.a.label: a
}
```

**ユースケース：**

- **ウィンドウ間通信**：あるウィンドウから別のウィンドウへ通知
- **バックグラウンドタスク**：GoサービスからUIへ進捗を通知
- **状態の同期**：複数のウィンドウを同期した状態に維持
- **疎結合**：コンポーネント間の直接参照が不要

**例：**

```go
// Go: Emit an event
app.Event.Emit("user-logged-in", user)
```

```javascript
// JavaScript: Listen for event
import { Events } from '@wailsio/runtime'

Events.On('user-logged-in', (user) => {
    console.log('User logged in:', user)
})
```

[イベントの詳細 →](/features/events/system/)

## アプリケーションのライフサイクル

ライフサイクルを理解すると、リソースを初期化するタイミングとクリーンアップするタイミングが分かります。

```d2
direction: down

Start: アプリケーションの起動 {
  shape: oval
  style.fill: "#10B981"
}

Init: 初期化 {
  Create: アプリケーションを作成 {
    shape: rectangle
  }

  Register: サービスを登録 {
    shape: rectangle
  }

  Setup: ウィンドウ／メニューを設定 {
    shape: rectangle
  }
}

Run: イベントループ {
  Events: イベントを処理 {
    shape: rectangle
  }

  Messages: メッセージを処理 {
    shape: rectangle
  }

  Render: UI を更新 {
    shape: rectangle
  }
}

Shutdown: シャットダウン {
  Cleanup: リソースを解放 {
    shape: rectangle
  }

  Save: 状態を保存 {
    shape: rectangle
  }
}

End: アプリケーションの終了 {
  shape: oval
  style.fill: "#EF4444"
}

Start -> Init.Create
Init.Create -> Init.Register
Init.Register -> Init.Setup
Init.Setup -> Run.Events
Run.Events -> Run.Messages
Run.Messages -> Run.Render
Run.Render -> Run.Events: ループ
Run.Events -> Shutdown.Cleanup: 終了シグナル
Shutdown.Cleanup -> Shutdown.Save
Shutdown.Save -> End
```

**ライフサイクルフック：**

```go
app := application.New(application.Options{
    Name: "My App",

    // Cleanly intercept quit requests (e.g. unsaved changes).
    ShouldQuit: func() bool { return true },

    // Called when the app is confirmed to be quitting — save state, close connections, etc.
    OnShutdown: func() {},
})
```

`application.Options`には`OnStartup`フィールドはありません。起動時の処理は、サービスの`ServiceStartup(ctx, options)`、`app.Event.OnApplicationEvent(events.Common.ApplicationStarted, ...)`で登録したコールバック、または単に`app.Run()`より前に配置します。

[ライフサイクルの詳細 →](/concepts/lifecycle/)

## ビルドプロセス

Wailsがアプリケーションをビルドする仕組みを見てみましょう。

```d2
direction: down

Source: ソースコード {
  Go: "Go コード\n(main.go, サービス)" {
    shape: rectangle
    style.fill: "#00ADD8"
  }

  Frontend: "フロントエンドコード\n(HTML/CSS/JS)" {
    shape: rectangle
    style.fill: "#8B5CF6"
  }
}

Build: ビルドプロセス {
  AnalyseGo: Go コードを解析 {
    shape: rectangle
  }

  GenerateBindings: バインディングを生成 {
    shape: rectangle
  }

  BuildFrontend: フロントエンドをビルド {
    shape: rectangle
  }

  CompileGo: Go をコンパイル {
    shape: rectangle
  }

  Embed: アセットを埋め込む {
    shape: rectangle
  }
}

Output: 出力 {
  Binary: "ネイティブバイナリ\n(myapp.exe/.app)" {
    shape: rectangle
    style.fill: "#10B981"
  }
}

Source.Go -> Build.AnalyseGo
Build.AnalyseGo -> Build.GenerateBindings: 型を抽出
Build.GenerateBindings -> Source.Frontend: TypeScript バインディング
Source.Frontend -> Build.BuildFrontend: コンパイル（Vite/webpack）
Build.BuildFrontend -> Build.Embed: バンドル済みアセット
Source.Go -> Build.CompileGo
Build.CompileGo -> Build.Embed
Build.Embed -> Output.Binary
```

**ビルド手順：**

1. **Goコードを解析**
  - サービスを走査してエクスポート済みメソッドを検出
  - パラメーターと戻り値の型を抽出
  - メソッドシグネチャを生成


2. **TypeScriptバインディングを生成**
  - サービスごとに`.ts`ファイルを作成
  - 完全な型定義を含める
  - JSDocコメントを追加


3. **フロントエンドをビルド**
  - バンドラー（Vite、webpackなど）を実行
  - ミニファイして最適化
  - `frontend/dist/`に出力


4. **Goをコンパイル**
  - 最適化を有効にしてコンパイル（`-ldflags="-s -w"`）
  - ビルドメタデータを含める
  - プラットフォーム固有のコンパイルを実行


5. **アセットを埋め込む**
  - フロントエンドファイルをGoバイナリに埋め込む
  - アセットを圧縮
  - 単一の実行可能ファイルを作成


<strong>結果：</strong>すべてが埋め込まれた単一のネイティブ実行可能ファイル。

[ビルドの詳細 →](/guides/build/building/)

## 開発環境と本番環境

Wailsの動作は、開発環境と本番環境で異なります。

@tabs{sync-key="mode"}
[開発環境（wails3 dev）]
**特性：**

- **ホットリロード**：フロントエンドの変更を即座に再読み込み
- **ソースマップ**：元のソースコードを使用してデバッグ
- **DevTools**：ブラウザーのDevToolsを使用可能
- **ログ出力**：詳細なログ出力を有効化
- **外部フロントエンド**：開発サーバー（Vite）から配信

**仕組み：**

```d2
direction: right

WailsApp: Wails アプリ {
  shape: rectangle
  style.fill: "#00ADD8"
}

DevServer: "Vite 開発サーバー\n(localhost:5173)" {
  shape: rectangle
  style.fill: "#8B5CF6"
}

WebView: WebView {
  shape: rectangle
  style.fill: "#6B7280"
}

WailsApp -> DevServer: リクエストをプロキシ
DevServer -> WebView: HMR で配信
WebView -> WailsApp: Go メソッドを呼び出す
```

**利点：**

- 変更結果を即座に確認
- すべてのデバッグ機能を利用可能
- 反復作業を高速化

[本番環境（wails3 build）]
**特性：**

- **埋め込みアセット**：フロントエンドをバイナリに組み込み
- **最適化済み**：ミニファイおよび圧縮済み
- **DevToolsなし**：デフォルトで無効
- **最小限のログ出力**：エラーのみ
- **単一ファイル**：すべてを1つの実行可能ファイルに格納

**仕組み：**

```d2
direction: right

Binary: "単一バイナリ\n(myapp.exe)" {
  GoCode: コンパイル済み Go コード {
    shape: rectangle
    style.fill: "#00ADD8"
  }

  Assets: "埋め込みアセット\n(HTML/CSS/JS)" {
    shape: rectangle
    style.fill: "#8B5CF6"
  }
}

WebView: WebView {
  shape: rectangle
  style.fill: "#6B7280"
}

Binary.Assets -> WebView: メモリから配信
WebView -> Binary.GoCode: Go メソッドを呼び出す
```

**利点：**

- 単一ファイルで配布
- サイズを縮小（ミニファイ済み）
- パフォーマンスの向上
- 外部依存関係なし

@end

## メモリモデル

メモリ使用量を理解すると、効率的なアプリケーションを構築できます。

**メモリ領域：**

1. **Goヒープ**
  - サービスとアプリケーションの状態
  - Goのガベージコレクターによって管理
  - 単純なアプリでは通常5-10MB


2. **WebViewのメモリ**
  - DOM、JavaScriptヒープ、CSS
  - WebViewのエンジンによって管理
  - 単純なアプリでは通常10-20MB


3. **ブリッジのメモリ**
  - 通信用のメッセージバッファ
  - 最小限のオーバーヘッド（<1MB）
  - 可能な場合は大容量データをゼロコピーで処理


**最適化のヒント：**

- **大容量データの転送を避ける**：IDを渡し、詳細は必要に応じて取得します
- **更新にはイベントを使用する**：フロントエンドからポーリングしないでください
- **大きなファイルをストリーミングする**：ファイル全体をメモリに読み込まないでください
- **リスナーをクリーンアップする**：不要になったイベントリスナーは削除してください

[パフォーマンスについて詳しく見る →](/guides/performance/)

## セキュリティモデル

Wailsは、デフォルトで安全なアーキテクチャを提供します：

```d2
direction: down

Frontend: フロントエンド（信頼されていない） {
  shape: rectangle
  style.fill: "#EF4444"
}

Bridge: Wails ブリッジ（検証） {
  shape: diamond
  style.fill: "#F59E0B"
}

Backend: バックエンド（信頼されている） {
  shape: rectangle
  style.fill: "#10B981"
}

Frontend -> Bridge: メソッドを呼び出す
Bridge -> Bridge: "検証：\n- メソッドは存在するか？\n- 型は正しいか？\n- アクセスは許可されているか？"
Bridge -> Backend: 有効な場合に実行
Backend -> Bridge: 結果を返す
Bridge -> Frontend: レスポンスを送信
```

**セキュリティ機能：**

1. **メソッドの許可リスト**
  - 呼び出せるのはエクスポートされたメソッドのみ
  - 非公開メソッドにはアクセス不可
  - サービスの明示的な登録が必要


2. **型の検証**
  - 引数をGoの型と照合
  - 無効な型を拒否
  - インジェクション攻撃を防止


3. **eval()を使用しない**
  - フロントエンドから任意のGoコードを実行することは不可
  - 呼び出せるのは事前定義されたメソッドのみ
  - 動的なコード実行なし


4. **コンテキストの分離**
  - 各ウィンドウが独自のコンテキストを保持
  - サービスで呼び出し元のコンテキストを確認可能
  - ウィンドウごとの権限設定が可能


**ベストプラクティス：**

- Goで<strong>ユーザー入力を検証する</strong>（フロントエンドを信頼しないでください）
- 認証と認可には<strong>コンテキストを使用する</strong>
- ファイル操作の前に<strong>ファイルパスをサニタイズする</strong>
- 高コストな処理を<strong>レート制限する</strong>

[セキュリティについて詳しく見る →](/guides/security/)

## 次のステップ

**アプリケーションのライフサイクル** - 起動、終了、ライフサイクルフックについて理解します [詳しく見る →](/concepts/lifecycle/)

**Goとフロントエンド間のブリッジ** - ブリッジの仕組みを詳しく掘り下げます [詳しく見る →](/concepts/bridge/)

**ビルドシステム** - Wailsがアプリケーションをビルドする仕組みを理解します [詳しく見る →](/concepts/build-system/)

**構築を始める** - ここまでに学んだことをチュートリアルで実践します [チュートリアル →](/tutorials/03-notes-vanilla/)

---

**アーキテクチャについて質問がありますか？** [Discord](https://discord.gg/JDdSxwjhGf)で質問するか、[APIリファレンス](/reference/overview/)を確認してください。
