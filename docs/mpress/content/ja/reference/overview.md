---
title: "API リファレンス"
description: "Wails の型、メソッド、オプション、CLI コマンドを調べる"
slug: "reference/overview"
sourcePath: "reference/overview.md"
---

## このリファレンスについて

アプリケーションの開発中に Wails の型、メソッド、オプション、コマンドを調べるには、このセクションを使用します。手順に沿って学ぶ場合は、[チュートリアル](/tutorials/overview/)から始めてください。タスクを完了するための手順については、[ハウツーガイド](/quick-start/next-steps/#features)を参照してください。設計の背景については、[コアコンセプト](/concepts/architecture/)を参照してください。

**構成：**

- [アプリケーション](/reference/application/) - アプリケーションのコア API
- [ウィンドウ](/reference/window/) - ウィンドウの作成と管理。構成フィールドについては、[ウィンドウオプション](/features/windows/options/)を参照してください
- [メニュー](/reference/menu/) - アプリケーションメニュー、コンテキストメニュー、システムトレイの API。項目の種類とプロパティについては、[メニュー項目リファレンス](/features/menus/reference/)を参照してください
- [イベント](/reference/events/) - イベントシステムと組み込みイベント
- [ダイアログ](/reference/dialogs/) - ファイルダイアログとメッセージダイアログ
- [フロントエンドランタイム](/reference/frontend-runtime/) - フロントエンドランタイム API
- [CLI コマンドの概要](/reference/cli/) - 目的別に分類されたコマンド。各コマンドの詳細については、[フラグと使用例](/guides/cli/)を参照してください
- [更新マニフェストプロトコル](/reference/update-manifest/) - 更新メタデータと検証フィールド
- [モバイル API](/guides/mobile/mobile-api/) - 試験的なモバイル API

## API の規約

@details{title="Go API の規約 - Go を初めて使用する開発者向け"}
### 命名規則

- <strong>型</strong>：PascalCase（例：`WebviewWindow`）
- <strong>メソッド</strong>：PascalCase（例：`SetTitle()`）
- <strong>オプション</strong>：PascalCase の構造体（例：`WindowOptions`）
- <strong>定数</strong>：PascalCase（例：`WindowStartStateMaximised`）

#### エラー処理

失敗する可能性があるほとんどのメソッドは、最後の戻り値として `error` を返します。`app.Run()` はアプリケーションが終了するまでブロックし、起動時に発生したエラーがあればそれを返します：

```go
if err := app.Run(); err != nil {
    log.Fatal(err)
}
```

ウィンドウの構築ではエラーは返されません。`app.Window.New()` は `*WebviewWindow` を直接返します。

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

構成にはオプション構造体を使用します：

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

すべての Go メソッド呼び出しは Promise を返します：

```javascript
// Async/await (recommended)
const result = await MyService.DoSomething()

// Promise chain
MyService.DoSomething()
    .then(result => console.log(result))
    .catch(error => console.error(error))
```

#### エラー処理

Go のエラーは JavaScript の例外になります：

```javascript
try {
    await MyService.MightFail()
} catch (error) {
    console.error('Go error:', error)
}
```

#### 型安全性

TypeScript の型定義は自動生成されます：

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

## プラットフォーム間の相違点

一部の API はプラットフォームによって動作が異なります：

| 機能 | Windows | macOS | Linux |
| --- | --- | --- | --- |
| **アプリケーションメニュー** | ウィンドウのメニューバー | グローバルメニューバー | ウィンドウのメニューバー |
| **システムトレイ** | 通知領域 | メニューバー | システムトレイ |
| **Dock** | 該当なし | ✅ 利用可能 | 該当なし |
| **ファイルダイアログ** | ネイティブ | ネイティブ | ネイティブ（GTK） |
| **透過表示** | ✅ 完全対応 | [`private_mac_apis`](/guides/build/private-macos-apis/#webview-transparency-and-background) が必要 | ⚠️ 制限あり |

プラットフォーム固有の動作については、各 API セクションに記載されています。

## バージョニング

Wails v3 はセマンティックバージョニングに従います：

- **メジャー**（v4.0.0）：破壊的変更
- **マイナー**（v3.1.0）：後方互換性のある新機能
- **パッチ**（v3.0.1）：後方互換性のあるバグ修正

**Wails v3.0.0：** 安定版

## 非推奨ポリシー

API が非推奨になる場合：

1. **ドキュメントに記載**し、非推奨である旨を通知
2. **代替手段を提供**し、移行ガイドを用意
3. 削除まで **1 メジャーバージョン**にわたって保守
4. **コンパイラ警告**（可能な場合）

## API の安定性

### 安定版 API ✅

以下の API は安定しており、本番環境で安全に使用できます：

- アプリケーションのコア API
- ウィンドウ管理
- メニューシステム
- イベントシステム
- ファイルダイアログ
- サービスバインディング

### 実験的 API ⚠️

安定版デスクトップの互換性契約の対象外となるのは、実験的であると明記された以下の API のみです。

- Android サポート
- iOS サポート
- その他、実験的であると明記された API

実験的 API はドキュメント内で明示されています。

## ヘルプ

### API に関する質問

1. **このリファレンスを確認** - 完全な API ドキュメント
2. **サンプルを確認** - [GitHub のサンプル](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples)
3. **Discord を検索** - [Discord サーバー](https://discord.gg/JDdSxwjhGf)
4. **コミュニティに質問** - Discord の #help チャンネル

### API の問題を報告する

バグや不整合を見つけましたか？

1. **既存の Issue を確認** - [GitHub Issues](https://github.com/wailsapp/wails/issues)
2. **詳細なレポートを作成** - コード、エラー、プラットフォームを含める
3. **再現手順を提供** - 問題を再現できる最小限のサンプル

## 関連ドキュメント

- [チュートリアル](/tutorials/overview/) - 実際のアプリケーションを構築しながら学ぶ
- [ガイド](/guides/architecture/) - 一般的なシナリオに対応するタスク指向のガイド
- [機能](/features/windows/basics/) - 機能別のドキュメント
- [サンプル](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples) - GitHub 上の動作するコードサンプル

---

**API を参照：** 左側のナビゲーションを使用して、個別の API を確認できます。
