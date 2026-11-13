---
title: "Wails の仕組み"
description: "Wails のアーキテクチャと、ネイティブ性能を実現する仕組みを理解する"
slug: "concepts/architecture"
sourcePath: "concepts/architecture.md"
---

Wails は、**バックエンドに Go**、**フロントエンドに Web 技術**を使用してデスクトップアプリケーションを構築するためのフレームワークです。ただし Electron とは異なり、ブラウザーを同梱せず、**オペレーティングシステムのネイティブ WebView**を使用します。

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
    frontend."Render Site via OS-native WebView": OS ネイティブの WebView でサイトを描画
  }
  Regular Communication: 通常の通信 {
    shape: sequence_diagram
    frontend."Make API-style call": API 形式の呼び出しを実行
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
| **ブラウザー** | OS が提供する WebView | 同梱の Chromium（約100MB） |
| **バックエンド** | Go（コンパイル済み） | Node.js（インタープリター実行） |
| **通信** | インメモリブリッジ | IPC（プロセス間通信） |
| **バンドルサイズ** | 約15MB | 約150MB |
| **メモリ** | 約10MB | 約100MB 以上 |
| **起動時間** | &lt;0.5秒 | 2-3秒 |

## コアコンポーネント

### 1. ネイティブ WebView

Wails は、オペレーティングシステムに組み込まれた次の Web レンダリングエンジンを使用します。

@tabs{sync-key="platform"}
[Windows]
**WebView2**（Microsoft Edge WebView2）

- Chromium ベース（Edge ブラウザーと同じ）
- Windows 10/11にプリインストール済み
- Windows Update による自動更新
- 最新の Web 標準を完全にサポート

[macOS]
**WebKit**（Safari のレンダリングエンジン）

- macOS に組み込み済み
- Safari ブラウザーと同じエンジン
- 優れたパフォーマンスとバッテリー駆動時間
- 最新の Web 標準を完全にサポート

[Linux]
**WebKitGTK**（WebKit の GTK ポート）

- パッケージマネージャー経由でインストール
- GNOME Web（Epiphany）と同じエンジン
- Web 標準を幅広くサポート
- 軽量で高性能

@end

**これが重要な理由：**

- **ブラウザーを同梱しない** → アプリのサイズを削減
- **OS ネイティブ** → 統合性とパフォーマンスが向上
- **自動更新** → OS の更新からセキュリティパッチを取得
- **使い慣れたレンダリング** → システムブラウザーと同じ

### 2. Wails ブリッジ

ブリッジは Wails の中核であり、Go と JavaScript の間の**直接通信**を可能にします。

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
2. **ブリッジが呼び出しを JSON にエンコードする**（メソッド名と引数）
3. **ルーターが登録済みサービスから Go メソッドを見つける**
4. **Go メソッドが実行され**、値を返す
5. **ブリッジが結果をデコードし**、フロントエンドへ返す
6. JavaScript で**Promise が結果によって解決される**

**パフォーマンス特性：**

- **インメモリ**：ネットワークのオーバーヘッドも HTTP もなし
- 可能な場合は**ゼロコピー**（大規模データ向け）
- **デフォルトで非同期**：両側ともノンブロッキング
- **型安全**：TypeScript 定義を自動生成

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

- Wails は起動時に**構造体をスキャンします**
- **エクスポートされたメソッド**はフロントエンドから呼び出せるようになります
- TypeScript バインディング用に**型情報**が抽出されます
- **エラー処理**は自動です（Go のエラー → JS の例外）

**生成された TypeScript バインディング：**

```typescript
// Auto-generated in frontend/bindings/GreetService.ts
export function Greet(name: string): Promise<string>
export function GetTime(): Promise<Date>
```

**サービスを使用する理由**

- **型安全**：TypeScript を完全にサポート
- **自動検出**：メソッドを手動で登録する必要がありません
- **体系的**：関連する機能をグループ化できます
- **テスト可能**：サービスは通常の Go 構造体です

[サービスの詳細を見る →](/features/bindings/services/)

### 4. イベントシステム

イベントを使用すると、コンポーネント間で**パブリッシュ／サブスクライブ通信**を行えます。

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
    backend.a."App Emit('data-updated', data)": "App Emit('data-updated', data)"
    backend.a -> window1.a: JSON イベントバス
    backend.a -> window2: JSON イベントバス
    window1.a."Subscriber processes On('data-updated', handler)": "サブスクライバーが On('data-updated', handler) を処理"
    window2."Subscriber processes On('data-updated', handler)": "サブスクライバーが On('data-updated', handler) を処理"
  }
  backend.a.label: a
  window1.a.label: a
}
```

**ユースケース：**

- **ウィンドウ間通信**：あるウィンドウからほかのウィンドウに通知します
- **バックグラウンドタスク**：Go サービスから UI に進捗を通知します
- **状態の同期**：複数のウィンドウの同期を維持します
- **疎結合**：コンポーネント間で直接参照する必要がありません

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

[イベントの詳細を見る →](/features/events/system/)

## アプリケーションのライフサイクル

ライフサイクルを理解すると、リソースを初期化するタイミングとクリーンアップするタイミングを把握できます。

```d2
direction: down

Start: アプリケーションの開始 {
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
  Cleanup: リソースをクリーンアップ {
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

`application.Options`には`OnStartup`フィールドがありません。起動時の処理は、サービスの`ServiceStartup(ctx, options)`、`app.Event.OnApplicationEvent(events.Common.ApplicationStarted, ...)`を介して登録したコールバック、または単に`app.Run()`より前に配置します。

[ライフサイクルの詳細を見る →](/concepts/lifecycle/)

## ビルドプロセス

Wails がアプリケーションをビルドする仕組みは次のとおりです。

```d2
direction: down

Source: ソースコード {
  Go: "Go コード\n(main.go、サービス)" {
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

  Embed: アセットを埋め込み {
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

1. **Go コードを解析**
  - サービスをスキャンしてエクスポートされたメソッドを検出します
  - パラメーターと戻り値の型を抽出します
  - メソッドシグネチャを生成します


2. **TypeScript バインディングを生成**
  - サービスごとに`.ts`ファイルを作成します
  - 完全な型定義を含めます
  - JSDoc コメントを追加します


3. **フロントエンドをビルド**
  - バンドラー（Vite、webpack など）を実行します
  - ミニファイして最適化します
  - `frontend/dist/`に出力します


4. **Go をコンパイル**
  - 最適化を有効にしてコンパイルします（`-ldflags="-s -w"`）
  - ビルドメタデータを含めます
  - プラットフォーム固有のコンパイルを行います


5. **アセットを埋め込み**
  - フロントエンドファイルを Go バイナリに埋め込みます
  - アセットを圧縮します
  - 単一の実行可能ファイルを作成します


<strong>結果：</strong>すべてが埋め込まれた単一のネイティブ実行可能ファイル。

[ビルドの詳細を見る →](/guides/build/building/)

## 開発環境と本番環境

Wails の動作は開発環境と本番環境で異なります。

@tabs{sync-key="mode"}
[開発環境（wails3 dev）]
**特徴：**

- **ホットリロード**：フロントエンドの変更を即座に再読み込み
- **ソースマップ**：元のソースコードでデバッグ可能
- **DevTools**：ブラウザーの DevTools を利用可能
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

- 変更結果を即座に確認可能
- 完全なデバッグ機能
- より迅速な反復開発

[本番環境（wails3 build）]
**特徴：**

- **埋め込みアセット**：フロントエンドをバイナリに組み込み
- **最適化済み**：縮小および圧縮済み
- **DevTools なし**：デフォルトで無効
- **最小限のログ出力**：エラーのみ
- **単一ファイル**：すべてを1つの実行ファイルに格納

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

- 単一ファイルで配布可能
- より小さいサイズ（縮小済み）
- パフォーマンスの向上
- 外部依存関係なし

@end

## メモリモデル

メモリ使用量を理解すると、効率的なアプリケーションを構築できます。

**メモリ領域：**

1. **Go ヒープ**
  - サービスとアプリケーションの状態
  - Go のガベージコレクターが管理
  - 単純なアプリでは通常 5-10MB


2. **WebView メモリ**
  - DOM、JavaScript ヒープ、CSS
  - WebView のエンジンが管理
  - 単純なアプリでは通常 10-20MB


3. **ブリッジメモリ**
  - 通信用のメッセージバッファー
  - 最小限のオーバーヘッド（<1MB）
  - 可能な場合は大容量データをゼロコピーで処理


**最適化のヒント：**

- **大容量データの転送を避ける**：ID を渡し、必要に応じて詳細を取得する
- **更新にはイベントを使用する**：フロントエンドからポーリングしない
- **大きなファイルはストリーミングする**：ファイル全体をメモリに読み込まない
- **リスナーをクリーンアップする**：使用後はイベントリスナーを削除する

[パフォーマンスの詳細 →](/guides/performance/)

## セキュリティモデル

Wails は、デフォルトで安全なアーキテクチャを提供します。

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

Backend: バックエンド（信頼済み） {
  shape: rectangle
  style.fill: "#10B981"
}

Frontend -> Bridge: メソッドを呼び出す
Bridge -> Bridge: "検証：\n- メソッドが存在するか？\n- 型が正しいか？\n- アクセスが許可されているか？"
Bridge -> Backend: 有効な場合に実行
Backend -> Bridge: 結果を返す
Bridge -> Frontend: レスポンスを送信
```

**セキュリティ機能：**

1. **メソッドの許可リスト**
  - エクスポートされたメソッドのみ呼び出し可能
  - 非公開メソッドにはアクセス不可
  - サービスの明示的な登録が必要


2. **型の検証**
  - 引数を Go の型と照合
  - 無効な型を拒否
  - インジェクション攻撃を防止


3. **eval() なし**
  - フロントエンドから任意の Go コードを実行することは不可能
  - 事前定義されたメソッドのみ呼び出し可能
  - 動的なコード実行なし


4. **コンテキストの分離**
  - 各ウィンドウが独自のコンテキストを保持
  - サービスは呼び出し元のコンテキストを確認可能
  - ウィンドウごとに権限を設定可能


**ベストプラクティス：**

- Go で**ユーザー入力を検証する**（フロントエンドを信頼しない）
- 認証と認可には**コンテキストを使用する**
- ファイル操作の前に**ファイルパスをサニタイズする**
- 高コストな操作を**レート制限する**

[セキュリティの詳細 →](/guides/security/)

## 次のステップ

**アプリケーションのライフサイクル** - 起動、終了、ライフサイクルフックについて理解する [詳細を見る →](/concepts/lifecycle/)

**Go とフロントエンドのブリッジ** - ブリッジの仕組みを詳しく理解する [詳細を見る →](/concepts/bridge/)

**ビルドシステム** - Wails がアプリケーションをビルドする仕組みを理解する [詳細を見る →](/concepts/build-system/)

**開発を始める** - チュートリアルで学んだことを実践しましょう [チュートリアル →](/tutorials/03-notes-vanilla/)

---

**アーキテクチャについて質問がありますか？** [Discord](https://discord.gg/JDdSxwjhGf) で質問するか、[API リファレンス](/reference/overview/)を確認してください。
