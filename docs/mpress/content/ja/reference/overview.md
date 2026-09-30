---
title: "API リファレンス"
description: "Wails v3 の完全な API ドキュメント"
slug: "reference/overview"
sourcePath: "reference/overview.md"
---

## このリファレンスについて

これは Wails v3 の完全な API リファレンスです。フレームワークで利用できるすべての公開型、メソッド、オプションについて説明します。

**構成：**

- [アプリケーション](/reference/application/) - アプリケーションのコア API
- [ウィンドウ](/reference/window/) - ウィンドウの作成と管理
- [メニュー](/reference/menu/) - アプリケーションメニュー、コンテキストメニュー、システムトレイメニュー
- [イベント](/reference/events/) - イベントシステムと組み込みイベント
- [ダイアログ](/reference/dialogs/) - ファイルダイアログとメッセージダイアログ
- [フロントエンドランタイム](/reference/frontend-runtime/) - フロントエンドランタイム API
- [CLI](/reference/cli/) - コマンドラインインターフェース

## API の規約

@details{title="Go API の規則 - Go 初心者の開発者向け"}
### 命名規則

- <strong></strong>型<strong></strong>：PascalCase（例：`WebviewWindow`）
- <strong></strong>メソッド<strong></strong>：PascalCase（例：`SetTitle()`）
- <strong></strong>オプション<strong></strong>：PascalCase の構造体（例：`WindowOptions`）
- <strong></strong>定数<strong></strong>：PascalCase（例：`WindowStartStateMaximised`）

#### エラー処理

失敗する可能性があるほとんどのメソッドは、最後の戻り値として `error` を返します。`app.Run()` はアプリケーションが終了するまでブロックし、起動時に発生したエラーがあれば返します：

```go
if err := app.Run(); err != nil {
    log.Fatal(err)
}
```

ウィンドウの構築はエラーを返しません。`app.Window.New()` は `*WebviewWindow` を直接返します。

#### コンテキスト

サービスのライフサイクルメソッドは `context.Context` を受け取ります：

```go
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    // ctx is cancelled when the application is shutting down.
    return nil
}
```

アプリケーションのライフタイムコンテキストは `app.Context()` から取得できます。`RunWithContext` は存在しないため、`app.Run()` を呼び出してください。

#### オプションパターン

設定にはオプション構造体を使用します：

```go
app := application.New(application.Options{
    Name: "My App",
    Description: "A demo application",
    Services: []application.Service{
        application.NewService(&MyService{}),
    },
})
```

@end

### JavaScript API の規約

#### 命名規則

- **関数**：camelCase（例：`setTitle()`）
- **定数**：SCREAMING<em>SNAKE</em>CASE（例：`WINDOW_EVENT_FOCUS`）

#### デフォルトで非同期

すべての Go メソッド呼び出しは Promise を返します。

```javascript
// Async/await (recommended)
const result = await MyService.DoSomething()

// Promise chain
MyService.DoSomething()
    .then(result => console.log(result))
    .catch(error => console.error(error))
```

#### エラー処理

Go のエラーは JavaScript の例外になります。

```javascript
try {
    await MyService.MightFail()
} catch (error) {
    console.error('Go error:', error)
}
```

#### 型安全性

TypeScript の型定義は自動生成されます。

```typescript
// Fully typed
import { Greet } from './bindings/GreetService'

const message: string = await Greet("World")
```

## パッケージ構成

```
github.com/wailsapp/wails/v3/pkg/
├── application/          # Core application package
│   ├── application.go    # App type
│   ├── webview_window.go # Window management
│   ├── menu.go           # Menu types
│   ├── event_manager.go  # Event system
│   └── dialogs.go        # Dialog APIs
├── events/               # Event constants
└── services/             # Built-in services
    ├── dock/             # macOS dock (includes badge support)
    ├── fileserver/       # File-server service
    ├── kvstore/          # Key/value store
    ├── log/              # Structured logging service
    ├── notifications/    # Notifications service
    └── sqlite/           # SQLite service
```

## インポートパス

### Go

```go
import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)
```

### JavaScript

```javascript
// Auto-generated bindings
import { MyMethod } from './bindings/MyService'

// Runtime APIs
import { Events, Window } from '@wailsio/runtime'
```

## 型リファレンス

### 共通の型

@tabs{sync-key="lang"}
[Go]
```go
// Application
type App struct { /* ... */ }
type Options struct { /* ... */ }

// Window
type WebviewWindow struct { /* ... */ } // implements the Window interface
type WebviewWindowOptions struct { /* ... */ }

// Menu
type Menu struct { /* ... */ }
type MenuItem struct { /* ... */ }

// Events — there is no generic Event type; events are typed by source.
type ApplicationEvent struct { /* ... */ }
type WindowEvent struct { /* ... */ }
type CustomEvent struct { /* ... */ }
type EventListener struct { /* ... */ }

// Dialogs
type OpenFileDialogOptions struct { /* ... */ }
type SaveFileDialogOptions struct { /* ... */ }
```

[TypeScript]
```typescript
// Window runtime
interface WindowOptions {
    title?: string
    width?: number
    height?: number
    // ...
}

// Events
type EventCallback = (data: any) => void

// Bindings (auto-generated)
export function MyMethod(arg: string): Promise<string>
```

@end

## プラットフォームによる違い

一部の API は、プラットフォームによって動作が異なります。

| 機能 | Windows | macOS | Linux |
| --- | --- | --- | --- |
| **アプリケーションメニュー** | ウィンドウのメニューバー | グローバルメニューバー | ウィンドウのメニューバー |
| **システムトレイ** | 通知領域 | メニューバー | システムトレイ |
| **Dock** | 該当なし | ✅ 利用可能 | 該当なし |
| **ファイルダイアログ** | ネイティブ | ネイティブ | ネイティブ（GTK） |
| **透過** | ✅ 完全対応 | [`private_mac_apis`](/guides/build/private-macos-apis/#webview-transparency-and-background)が必要 | ⚠️ 制限あり |

プラットフォーム固有の動作については、各 API セクションで説明します。

## バージョニング

Wails v3 はセマンティックバージョニングに従います。

- **メジャー**（v3.x.x）：破壊的変更
- **マイナー**（v3.x.x）：後方互換性のある新機能
- **パッチ**（v3.x.x）：後方互換性のあるバグ修正

<strong>現在のステータス：</strong>ベータ（API は安定、改良は継続中）

## 非推奨化ポリシー

API が非推奨になった場合：

1. <strong>ドキュメントに明記</strong>し、非推奨であることを通知
2. <strong>代替手段を提示</strong>し、移行ガイドを提供
3. 削除前に<strong>1 メジャーバージョン分</strong>の期間は保守されます
4. **コンパイラ警告**（可能な場合）

## API の安定性

### 安定版 API ✅

以下の API は安定しており、本番環境で安全に使用できます。

- コアアプリケーション API
- ウィンドウ管理
- メニューシステム
- イベントシステム
- ファイルダイアログ
- サービスバインディング

### 不安定版 API ⚠️

以下の API は正式リリースまでに変更される可能性があります。

- 一部の高度なウィンドウオプション
- プラットフォーム固有の機能
- 実験的機能

不安定版 API はドキュメントに明記されています。

## ヘルプの利用

### API に関する質問

1. **このリファレンスを確認** — 完全な API ドキュメント
2. **サンプルを確認** — [GitHub のサンプル](https://github.com/wailsapp/wails/tree/master/v3/examples)
3. **Discord を検索** — [Discord サーバー](https://discord.gg/JDdSxwjhGf)
4. **コミュニティに質問** — Discord の #help チャンネル

### API の問題の報告

バグや不整合を見つけましたか？

1. **既存の Issue を確認** — [GitHub Issues](https://github.com/wailsapp/wails/issues)
2. **詳細なレポートを作成** — コード、エラー、プラットフォームを含める
3. **再現手順を提示** — 問題を再現する最小限のサンプル

## 関連ドキュメント

- [チュートリアル](/tutorials/overview/) — 実際のアプリケーションを構築しながら学習
- [ガイド](/guides/architecture/) — 一般的なシナリオ向けのタスク別ガイド
- [機能](/features/windows/basics/) — 機能ごとのドキュメント
- [サンプル](https://github.com/wailsapp/wails/tree/master/v3/examples) — GitHub 上の動作するコード例

---

<strong>API を参照：</strong>左側のナビゲーションを使用して、個々の API を確認できます。
