---
title: "イベントガイド"
description: "アプリケーション内の通信とライフサイクル管理に Wails v3 のイベントを活用するための実践ガイド"
slug: "guides/events-reference"
sourcePath: "guides/events-reference.md"
---

**注: このガイドは現在作成中です**

## イベントガイド

イベントは、Wails アプリケーションにおける通信の要です。イベントを使うと、アプリケーションの各部分を密結合させることなく、相互に通信できます。このガイドでは、Wails アプリケーションでイベントを効果的に使用するために必要な知識をすべて説明します。

## Wails イベントの仕組み

イベントは、アプリケーション全体にブロードキャストされるメッセージと考えることができます。アプリケーションのどの部分でも、これらのメッセージをリッスンし、それに応じて処理できます。これは特に次の用途に役立ちます。

- **ウィンドウの変化への対応**: ウィンドウが最小化、最大化、または移動されたタイミングを把握する
- **システムイベントの処理**: テーマの変更や電源イベントに対応する
- **アプリケーション固有のロジック**: データの更新やユーザー操作などの機能に独自のイベントを作成する
- **コンポーネント間の通信**: アプリの各部分を直接依存させることなく通信させる

## イベントの命名規則

すべての Wails イベントは、発生元を明確に示す名前空間パターンに従います。

- `common:` - Windows、macOS、Linux で動作するクロスプラットフォームイベント
- `windows:` - Windows 固有のイベント
- `mac:` - macOS 固有のイベント\
- `linux:` - Linux 固有のイベント

例:

- `common:WindowFocus` - ウィンドウがフォーカスを得た（すべてのプラットフォームで動作）
- `windows:APMSuspend` - システムがサスペンドに移行している（Windows のみ）
- `mac:ApplicationDidBecomeActive` - アプリがアクティブになった（macOS のみ）

## イベントの使用を開始する

### イベントのリッスン（フロントエンド）

最も一般的な使用例は、フロントエンドコードでイベントをリッスンすることです。

```javascript
import { Events } from '@wailsio/runtime';

// Listen for when the window gains focus
Events.On('common:WindowFocus', () => {
    console.log('Window is now focused!');
    // Maybe refresh some data or resume animations
});

// Listen for theme changes
Events.On('common:ThemeChanged', (event) => {
    console.log('Theme changed:', event.data);
    // Update your app's theme accordingly
});

// Listen for custom events from your Go backend
Events.On('my-app:data-updated', (event) => {
    console.log('Data updated:', event.data);
    // Update your UI with the new data
});
```

### イベントの発行（バックエンド）

Go コードから、フロントエンドでリッスンできるイベントを発行できます。

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "time"
)

func (s *Service) UpdateData() {
    // Do some data processing...

    // Notify the frontend
    app := application.Get()
    app.Event.Emit("my-app:data-updated",
        map[string]interface{}{
            "timestamp": time.Now(),
            "count": 42,
        },
	)
}
```

### イベントの発行（フロントエンド）

一般的な使い方ではありませんが、フロントエンドから、Go コードでリッスンできるイベントを発行することもできます。

```javascript
import { Events } from '@wailsio/runtime';

// Event without data
Events.Emit('myapp:close-window')

// Event with data
Events.Emit('myapp:disconnect-requested', 'id-123')
```

フロントエンドで TypeScript を使用し、Go コードで[型付きイベントを登録](#heading-9)すると、イベント名の自動補完と検証、およびデータ型の検証を利用できます。

### イベントリスナーの削除

不要になったイベントリスナーは必ずクリーンアップしてください。

```javascript
import { Events } from '@wailsio/runtime';

// Store the handler reference
const focusHandler = () => {
    console.log('Window focused');
};

// Add the listener — capture the unsubscribe function it returns
const unsubscribe = Events.On('common:WindowFocus', focusHandler);

// Later, remove this specific listener via the returned unsubscribe
unsubscribe();

// Or remove ALL listeners for one (or more) event names — Events.Off takes only event-name strings
Events.Off('common:WindowFocus');
// Events.Off('common:WindowFocus', 'common:WindowLostFocus'); // variadic
// Events.OffAll(); // remove every listener for every event (no args)
```

## 一般的な使用例

### 1. ウィンドウのフォーカスに応じた一時停止と再開

多くのアプリケーションでは、ウィンドウがフォーカスを失ったときに特定の処理を一時停止する必要があります。

```javascript
import { Events } from '@wailsio/runtime';

let animationRunning = true;

Events.On('common:WindowLostFocus', () => {
    animationRunning = false;
    pauseBackgroundTasks();
});

Events.On('common:WindowFocus', () => {
    animationRunning = true;
    resumeBackgroundTasks();
});
```

### 2. テーマ変更への対応

アプリのテーマをシステムテーマと同期させます。

```javascript
import { Events } from '@wailsio/runtime';

Events.On('common:ThemeChanged', (event) => {
    const isDarkMode = event.data.isDarkMode;

    if (isDarkMode) {
        document.body.classList.add('dark-theme');
        document.body.classList.remove('light-theme');
    } else {
        document.body.classList.add('light-theme');
        document.body.classList.remove('dark-theme');
    }
});
```

### 3. ファイルドロップの処理

ドラッグされたファイルをアプリで受け取れるようにします。

```javascript
import { Events } from '@wailsio/runtime';

Events.On('common:WindowFilesDropped', (event) => {
    const files = event.data.files;

    files.forEach(file => {
        console.log('File dropped:', file);
        // Process the dropped files
        handleFileUpload(file);
    });
});
```

### 4. ウィンドウのライフサイクル管理

ウィンドウの状態変化に対応します。

```javascript
import { Events } from '@wailsio/runtime';

Events.On('common:WindowClosing', () => {
    // Save user data before closing
    saveApplicationState();

    // You could also prevent closing by returning false
    // from a registered window close handler
});

Events.On('common:WindowMaximise', () => {
    // Adjust UI for maximized view
    adjustLayoutForMaximized();
});

Events.On('common:WindowRestore', () => {
    // Return UI to normal state
    adjustLayoutForNormal();
});
```

### 5. プラットフォーム固有の機能

必要に応じて、プラットフォーム固有のイベントを処理します。

```javascript
import { Events } from '@wailsio/runtime';

// Windows-specific power management
Events.On('windows:APMSuspend', () => {
    console.log('System is going to sleep');
    saveState();
});

Events.On('windows:APMResumeSuspend', () => {
    console.log('System woke up');
    refreshData();
});

// macOS-specific app lifecycle
Events.On('mac:ApplicationWillTerminate', () => {
    console.log('App is about to quit');
    performCleanup();
});
```

## カスタムイベントの作成

アプリケーション固有の要件に合わせて独自のイベントを作成できます。

### バックエンド（Go）

```go
// Emit a custom event when data changes

func (s *Service) ProcessUserData(userData UserData) error {
    // Process the data...

    app := application.Get()
    // Notify all listeners
    app.Event.Emit("user:data-processed",
        map[string]interface{}{
            "userId": userData.ID,
            "status": "completed",
            "timestamp": time.Now(),
        },
    )
    return nil
}

// Emit periodic updates
func (s *Service) StartMonitoring() {
    app := application.Get()
    ticker := time.NewTicker(5 * time.Second)
    go func() {
        for range ticker.C {
            stats := s.collectStats()
            app.Event.Emit("monitor:stats-updated", stats)
        }
    }()
}
```

### フロントエンド（JavaScript）

```javascript
import { Events } from '@wailsio/runtime';

// Listen for your custom events
Events.On('user:data-processed', (event) => {
    const { userId, status, timestamp } = event.data;

    showNotification(`User ${userId} processing ${status}`);
    updateUIWithNewData();
});

Events.On('monitor:stats-updated', (event) => {
    updateDashboard(event.data);
});
```

## 型安全な型付きイベント

Wails v3 は、イベントの登録とバインディングの自動生成により、TypeScript で完全に型安全な型付きイベントをサポートします。

### カスタムイベントの登録

初期化時に `application.RegisterEvent` を呼び出し、カスタムイベント名とそのデータ型を登録します。

```go
package main

import "github.com/wailsapp/wails/v3/pkg/application"

type UserLoginData struct {
    UserID   string
    Username string
    LoginTime string
}

type MonitorStats struct {
    CPUUsage    float64
    MemoryUsage float64
}

func init() {
    // Register events with their data types
    application.RegisterEvent[UserLoginData]("user:login")
    application.RegisterEvent[MonitorStats]("monitor:stats")

    // Register events without data (void events)
    application.RegisterEvent[application.Void]("app:ready")
}
```

@note{type="caution"}
`RegisterEvent` は初期化時に呼び出すことを想定しており、次の場合はパニックします。

- 引数が無効である
- 同じイベント名が異なるデータ型で重複して登録されている

@end

@note{type="info"}
データ型が常に同じであれば、同じイベントを複数回登録しても安全です。これは、複数のパッケージのいずれかが読み込まれたときにイベントが確実に登録されるようにする場合に役立ちます。

@end

### イベント登録の利点

登録後、`Event.Emit` に渡されたデータ引数は、指定された型に対して型チェックされます。型が一致しない場合は、次のようになります。

- エラーが発行され、ログに記録される（または登録済みのエラーハンドラーに渡される）
- 問題のあるイベントは伝播されない
- これにより、登録済みイベントのデータフィールドには、宣言された型に代入可能な値だけが常に設定されます

### 厳格モード

開発時に未登録イベントの警告を有効にするには、`strictevents` ビルドタグを使用します。

```bash
go build -tags strictevents
```

厳格モードを有効にすると、ログが警告で埋まるのを防ぐため、ランタイムが発行する警告は未登録のイベント名ごとに最大 1 回です。

### TypeScript バインディングの生成

バインディングジェネレーターは、フロントエンドで型付きイベントを透過的にサポートするための TypeScript 定義とグルーコードを出力します。

#### 1. Vite プラグインの設定

`vite.config.ts` で次のように設定します。

```typescript
import { defineConfig } from 'vite'
import wails from '@wailsio/runtime/plugins/vite'

export default defineConfig({
  plugins: [wails()],
})
```

#### 2. バインディングの生成

バインディングジェネレーターを実行します。

```bash
wails3 generate bindings
```

これにより、型付きイベントデータの定義を含む TypeScript ファイルがフロントエンドディレクトリに作成されます。さらに、登録済みのイベント名ごとに、ジェネレーターが検出可能なコンパイル時定数を 1 つエクスポートする `bindings/events` モジュールも作成されます（たとえば、`user:login` は `UserLogin` になります）。

#### 3. フロントエンドでの型付きイベントの使用

```typescript
import { Events } from '@wailsio/runtime'
import { UserLogin, MonitorStats } from './bindings/events'

// Type-safe event emission with autocomplete
Events.Emit(UserLogin, {
    UserID: "123",
    Username: "john_doe",
    LoginTime: new Date().toISOString()
})

// Type-safe event listening
Events.On(UserLogin, (event) => {
    // event.data is typed as UserLoginData
    console.log(`User ${event.data.Username} logged in`)
})

Events.On(MonitorStats, (event) => {
    // event.data is typed as MonitorStats
    updateDashboard({
        cpu: event.data.CPUUsage,
        memory: event.data.MemoryUsage
    })
})
```

型付きイベントでは、次の機能を利用できます。

- イベント名の**自動補完**
- イベントデータの**型チェック**
- データ型が一致しない場合の**コンパイル時エラー**
- **IntelliSense** ドキュメント

## イベントリファレンス

### 共通イベント（クロスプラットフォーム）

これらのイベントは、すべてのプラットフォームで動作します。

| イベント | 説明 | 使用するタイミング |
| --- | --- | --- |
| `common:ApplicationStarted` | アプリケーションが完全に起動した | アプリを初期化し、保存済みの状態を読み込む |
| `common:WindowRuntimeReady` | Wails ランタイムの準備が完了した | Wails API の呼び出しを開始する |
| `common:ThemeChanged` | システムテーマが変更された | アプリの外観を更新する |
| `common:SystemWillSleep` | システムがまもなくサスペンドする | 状態をフラッシュし、ソケットを閉じる |
| `common:SystemDidWake` | システムがサスペンドから復帰した | 再接続し、古くなったデータを更新する |
| `common:WindowFocus` | ウィンドウがフォーカスを得た | 処理を再開し、データを更新する |
| `common:WindowLostFocus` | ウィンドウがフォーカスを失った | 処理を一時停止し、状態を保存する |
| `common:WindowMinimise` | ウィンドウが最小化された | レンダリングを一時停止し、リソース使用量を削減する |
| `common:WindowMaximise` | ウィンドウが最大化された | 全画面表示に合わせてレイアウトを調整する |
| `common:WindowRestore` | ウィンドウが最小化または最大化された状態から復元された | 通常のレイアウトに戻す |
| `common:WindowClosing` | ウィンドウがまもなく閉じる | データを保存し、リソースをクリーンアップする |
| `common:WindowFilesDropped` | ウィンドウにファイルがドロップされた | ファイルのインポートを処理する |
| `common:WindowDidResize` | ウィンドウのサイズが変更された | レイアウトを調整し、チャートを再レンダリングする |
| `common:WindowDidMove` | ウィンドウが移動された | 位置に依存する機能を更新する |

### プラットフォーム固有のイベント

#### Windows のイベント

Windows アプリケーションの主なイベント：

| イベント | 説明 | ユースケース |
| --- | --- | --- |
| `windows:SystemThemeChanged` | Windows のテーマが変更された | アプリの配色を更新する |
| `windows:APMSuspend` | システムがサスペンド中 | 状態を保存し、処理を一時停止する |
| `windows:APMResumeAutomatic` | システムが復帰した（復帰時には必ず発生） | 状態を復元し、データを更新する |
| `windows:APMResumeSuspend` | ユーザー入力によってシステムが復帰した（`APMResumeAutomatic` の後） | ユーザー操作による復帰を識別する |
| `windows:APMPowerStatusChange` | 電源状態が変更された | パフォーマンス設定を調整する |

#### macOS のイベント

macOS アプリケーションの重要なイベント：

| イベント | 説明 | ユースケース |
| --- | --- | --- |
| `mac:ApplicationDidBecomeActive` | アプリがアクティブになった | 処理を再開する |
| `mac:ApplicationDidResignActive` | アプリが非アクティブになった | 処理を一時停止する |
| `mac:ApplicationWillTerminate` | アプリが終了しようとしている | 最終的なクリーンアップを行う |
| `mac:ApplicationWillSleep` | システムがまもなくサスペンドする | 状態を保存し、ソケットを閉じる |
| `mac:ApplicationDidWake` | システムが復帰した | 再接続し、更新する |
| `mac:ApplicationScreensDidSleep` | ディスプレイがスリープ状態になった | レンダリングを一時停止する（システムのスリープとは別） |
| `mac:ApplicationScreensDidWake` | ディスプレイがスリープから復帰した | レンダリングを再開する |
| `mac:WindowDidEnterFullScreen` | 全画面表示になった | 全画面表示に合わせて UI を調整する |
| `mac:WindowDidExitFullScreen` | 全画面表示を終了した | 通常の UI に戻す |

#### Linux のイベント

Linux の主要なウィンドウイベント：

| イベント | 説明 | ユースケース |
| --- | --- | --- |
| `linux:SystemThemeChanged` | デスクトップテーマが変更された | アプリのテーマを更新する |
| `linux:SystemWillSleep` | システムがまもなくサスペンドする（logind） | 状態を保存する |
| `linux:SystemDidWake` | システムが復帰した（logind） | 再接続し、更新する |
| `linux:WindowFocusIn` | ウィンドウがフォーカスを得た | 処理を再開する |
| `linux:WindowFocusOut` | ウィンドウがフォーカスを失った | 処理を一時停止する |
| `linux:WindowLoadStarted` | WebView が読み込みを開始した | 読み込みインジケーターを表示する |
| `linux:WindowLoadRedirected` | WebView がリダイレクトされた | ナビゲーションのリダイレクトを追跡 |
| `linux:WindowLoadCommitted` | WebView が読み込みを確定 | コンテンツを受信中 |
| `linux:WindowLoadFinished` | WebView の読み込みが完了 | 読み込みインジケーターを非表示にし、JS/CSS を挿入 |

## ベストプラクティス

### 1. イベント名前空間を使用する

カスタムイベントを作成する際は、競合を避けるために名前空間を使用します。

```javascript
import { Events } from '@wailsio/runtime';

// Good - namespaced events
Events.Emit('myapp:user:login');
Events.Emit('myapp:data:updated');
Events.Emit('myapp:network:connected');

// Avoid - generic names that might conflict
Events.Emit('login');
Events.Emit('update');
```

### 2. リスナーをクリーンアップする

コンポーネントがアンマウントされる際は、必ずイベントリスナーを削除します。

```javascript
import { Events } from '@wailsio/runtime';

// React example
useEffect(() => {
    const handler = (event) => {
        // Handle event
    };

    const off = Events.On('common:WindowDidResize', handler);

    // Cleanup — call the unsubscribe returned by Events.On
    return () => {
        off();
    };
}, []);
```

### 3. プラットフォームの違いに対応する

プラットフォーム固有のイベントを使用する際は、そのプラットフォームで利用できるか確認します。

```javascript
import { Events } from '@wailsio/runtime';

// Platform-specific events can be registered unconditionally;
// they will simply never fire on unsupported platforms.
Events.On('windows:APMSuspend', handleSuspend);
Events.On('mac:ApplicationWillTerminate', handleTerminate);
```

### 4. イベントを多用しない

イベントは強力ですが、あらゆる処理に使用しないでください。

- ✅ イベントの使用に適したもの：システム通知、ライフサイクルの変更、更新のブロードキャスト
- ❌ イベントの使用を避けるもの：関数からの直接の戻り値、単一コンポーネントの更新、同期処理

## イベントのデバッグ

イベントの問題をデバッグするには、次の手順を実行します。

```javascript
import { Events } from '@wailsio/runtime';

// Log all events (development only)
if (isDevelopment) {
    const originalOn = Events.On;
    Events.On = function(eventName, handler) {
        console.log(`[Event Registered] ${eventName}`);
        return originalOn.call(this, eventName, function(event) {
            console.log(`[Event Fired] ${eventName}`, event);
            return handler(event);
        });
    };
}
```

## 信頼できる情報源

利用可能なイベントの完全な一覧は、Wails のソースコードで確認できます。

- フロントエンドイベント：[`v3/internal/runtime/desktop/@wailsio/runtime/src/event_types.ts`](https://github.com/wailsapp/wails/blob/main/v3/internal/runtime/desktop/@wailsio/runtime/src/event_types.ts)
- バックエンドイベント：[`v3/pkg/events/events.go`](https://github.com/wailsapp/wails/blob/main/v3/pkg/events/events.go)

最新のイベント名と利用可否については、必ずこれらのファイルを参照してください。

## まとめ

Wails のイベントは、アプリケーション内の通信を処理するための強力で疎結合な手段です。このガイドのパターンとプラクティスに従うことで、システムの変更やユーザー操作に円滑に反応する、応答性が高くプラットフォームを考慮したアプリケーションを構築できます。

覚えておいてください。クロスプラットフォーム互換性のために共通イベントから始め、必要に応じてプラットフォーム固有のイベントを追加し、メモリリークを防ぐためにイベントリスナーを必ずクリーンアップします。
