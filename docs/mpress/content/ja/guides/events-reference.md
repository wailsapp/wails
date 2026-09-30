---
title: "イベントガイド"
description: "Wails v3でイベントを使用してアプリケーション内の通信とライフサイクルを管理するための実践ガイド"
slug: "guides/events-reference"
sourcePath: "guides/events-reference.md"
---

**注: このガイドは作成中です**

## イベントガイド

イベントは、Wailsアプリケーションにおける通信の要です。イベントを使用すると、アプリケーションの各部分を密結合させることなく、相互に通信できます。このガイドでは、Wailsアプリケーションでイベントを効果的に使用するために必要な知識をすべて説明します。

## Wailsのイベントについて

イベントは、アプリケーション全体にブロードキャストされるメッセージと考えることができます。アプリケーションのどの部分でも、これらのメッセージをリッスンし、それに応じて処理できます。イベントは、特に次の用途に役立ちます。

- **ウィンドウの変化への対応**: ウィンドウが最小化、最大化、または移動されたことを検知する
- **システムイベントの処理**: テーマの変更や電源イベントに対応する
- **アプリケーション固有のロジック**: データの更新やユーザー操作などの機能向けに独自のイベントを作成する
- **コンポーネント間の通信**: アプリの各部分が直接依存することなく通信できるようにする

## イベントの命名規則

すべてのWailsイベントは、その発生元を明確に示すために名前空間パターンに従います。

- `common:` - Windows、macOS、Linuxで動作するクロスプラットフォームイベント
- `windows:` - Windows固有のイベント
- `mac:` - macOS固有のイベント\
- `linux:` - Linux固有のイベント

例:

- `common:WindowFocus` - ウィンドウがフォーカスを得た（すべてのプラットフォームで動作）
- `windows:APMSuspend` - システムがサスペンドに移行している（Windowsのみ）
- `mac:ApplicationDidBecomeActive` - アプリがアクティブになった（macOSのみ）

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

Goコードから、フロントエンドがリッスンできるイベントを発行できます。

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

一般的な用途ではありませんが、フロントエンドから、Goコードがリッスンできるイベントを発行することもできます。

```javascript
import { Events } from '@wailsio/runtime';

// Event without data
Events.Emit('myapp:close-window')

// Event with data
Events.Emit('myapp:disconnect-requested', 'id-123')
```

フロントエンドでTypeScriptを使用し、Goコードで[型付きイベントを登録](#heading-9)すると、イベント名の自動補完と検証、およびデータ型の検証を利用できます。

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

アプリをシステムテーマと同期させます。

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

ドラッグされたファイルをアプリで受け入れられるようにします。

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

アプリケーション固有の要件に合わせて、独自のイベントを作成できます。

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

Wails v3では、イベントの登録とバインディングの自動生成により、TypeScriptで完全な型安全性を備えた型付きイベントを使用できます。

### カスタムイベントの登録

初期化時に`application.RegisterEvent`を呼び出し、カスタムイベント名とそのデータ型を登録します。

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
`RegisterEvent`は初期化時に呼び出すことを想定しており、次の場合はpanicします。

- 引数が有効でない
- 同じイベント名が異なるデータ型で2回登録されている

@end

@note{type="info"}
データ型が常に同じであれば、同じイベントを複数回登録しても安全です。これは、複数のパッケージのいずれかが読み込まれたときにイベントが確実に登録されるようにする場合に役立ちます。

@end

### イベント登録の利点

登録後、`Event.Emit`に渡されたデータ引数は、指定された型に対して型チェックされます。型が一致しない場合は、次のように処理されます。

- エラーが発行されてログに記録される（または登録済みのエラーハンドラーに渡される）
- 問題のあるイベントは伝播されない
- これにより、登録済みイベントのデータフィールドが、宣言された型に常に代入可能であることが保証される

### 厳格モード

開発時に未登録イベントの警告を有効にするには、`strictevents` ビルドタグを使用します。

```bash
go build -tags strictevents
```

strict モードを有効にすると、ログが警告で埋まるのを避けるため、ランタイムが出力する警告は未登録のイベント名ごとに最大1回となります。

### TypeScript バインディングの生成

バインディングジェネレーターは、フロントエンドで透過的に型付きイベントをサポートするための TypeScript 定義とグルーコードを出力します。

#### 1. Vite プラグインをセットアップする

`vite.config.ts` で次のように設定します。

```typescript
import { defineConfig } from 'vite'
import wails from '@wailsio/runtime/plugins/vite'

export default defineConfig({
  plugins: [wails()],
})
```

#### 2. バインディングを生成する

バインディングジェネレーターを実行します。

```bash
wails3 generate bindings
```

これにより、型付きイベントクリエーターとデータインターフェースを含む TypeScript ファイルがフロントエンドディレクトリに作成されます。

#### 3. フロントエンドで型付きイベントを使用する

```typescript
import { Events } from '@wailsio/runtime'
import { UserLogin, MonitorStats } from './bindings/events'

// Type-safe event emission with autocomplete
Events.Emit(UserLogin({
    UserID: "123",
    Username: "john_doe",
    LoginTime: new Date().toISOString()
}))

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

型付きイベントには、次の機能があります。

- イベント名の<strong>自動補完</strong>
- イベントデータの<strong>型チェック</strong>
- データ型が一致しない場合の<strong>コンパイル時エラー</strong>
- **IntelliSense** ドキュメント

## イベントリファレンス

### 共通イベント（クロスプラットフォーム）

次のイベントはすべてのプラットフォームで動作します。

| イベント | 説明 | 使用場面 |
| --- | --- | --- |
| `common:ApplicationStarted` | アプリケーションが完全に起動した | アプリを初期化し、保存済みの状態を読み込む |
| `common:WindowRuntimeReady` | Wails ランタイムの準備が完了した | Wails API の呼び出しを開始する |
| `common:ThemeChanged` | システムテーマが変更された | アプリの外観を更新する |
| `common:SystemWillSleep` | システムがまもなくサスペンドする | 状態を永続化し、ソケットを閉じる |
| `common:SystemDidWake` | システムがサスペンドから復帰した | 再接続し、古くなったデータを更新する |
| `common:WindowFocus` | ウィンドウがフォーカスを得た | 処理を再開し、データを更新する |
| `common:WindowLostFocus` | ウィンドウがフォーカスを失った | 処理を一時停止し、状態を保存する |
| `common:WindowMinimise` | ウィンドウが最小化された | レンダリングを一時停止し、リソース使用量を削減する |
| `common:WindowMaximise` | ウィンドウが最大化された | 全画面表示に合わせてレイアウトを調整する |
| `common:WindowRestore` | ウィンドウが最小化または最大化された状態から復元された | 通常のレイアウトに戻す |
| `common:WindowClosing` | ウィンドウがまもなく閉じる | データを保存し、リソースを解放する |
| `common:WindowFilesDropped` | ウィンドウにファイルがドロップされた | ファイルのインポートを処理する |
| `common:WindowDidResize` | ウィンドウのサイズが変更された | レイアウトを調整し、グラフを再レンダリングする |
| `common:WindowDidMove` | ウィンドウが移動された | 位置に依存する機能を更新する |

### プラットフォーム固有のイベント

#### Windows イベント

Windowsアプリケーションの主要なイベント：

| イベント | 説明 | ユースケース |
| --- | --- | --- |
| `windows:SystemThemeChanged` | Windowsのテーマが変更された | アプリの配色を更新する |
| `windows:APMSuspend` | システムがサスペンドに移行する | 状態を保存し、処理を一時停止する |
| `windows:APMResumeAutomatic` | システムが再開した（再開時に必ず発生） | 状態を復元し、データを更新する |
| `windows:APMResumeSuspend` | ユーザー入力によってシステムが再開した（`APMResumeAutomatic`の後） | ユーザー操作による復帰を判別する |
| `windows:APMPowerStatusChange` | 電源状態が変更された | パフォーマンス設定を調整する |

#### macOSイベント

macOSアプリケーションの重要なイベント：

| イベント | 説明 | ユースケース |
| --- | --- | --- |
| `mac:ApplicationDidBecomeActive` | アプリがアクティブになった | 処理を再開する |
| `mac:ApplicationDidResignActive` | アプリが非アクティブになった | 処理を一時停止する |
| `mac:ApplicationWillTerminate` | アプリが終了しようとしている | 最終的なクリーンアップを行う |
| `mac:ApplicationWillSleep` | システムがサスペンドに移行しようとしている | 状態を保存し、ソケットを閉じる |
| `mac:ApplicationDidWake` | システムが再開した | 再接続し、更新する |
| `mac:ApplicationScreensDidSleep` | ディスプレイがスリープ状態になった | レンダリングを一時停止する（システムのスリープとは異なる） |
| `mac:ApplicationScreensDidWake` | ディスプレイがスリープから復帰した | レンダリングを再開する |
| `mac:WindowDidEnterFullScreen` | フルスクリーン表示になった | フルスクリーン表示に合わせてUIを調整する |
| `mac:WindowDidExitFullScreen` | フルスクリーン表示を終了した | 通常のUIに戻す |

#### Linuxイベント

Linuxの主要なウィンドウイベント：

| イベント | 説明 | ユースケース |
| --- | --- | --- |
| `linux:SystemThemeChanged` | デスクトップのテーマが変更された | アプリのテーマを更新する |
| `linux:SystemWillSleep` | システムがサスペンドに移行しようとしている（logind） | 状態を保存する |
| `linux:SystemDidWake` | システムが再開した（logind） | 再接続し、更新する |
| `linux:WindowFocusIn` | ウィンドウがフォーカスを得た | アクティビティを再開する |
| `linux:WindowFocusOut` | ウィンドウがフォーカスを失った | 処理を一時停止する |
| `linux:WindowLoadStarted` | WebView が読み込みを開始した | 読み込みインジケーターを表示する |
| `linux:WindowLoadRedirected` | WebView がリダイレクトされた | ナビゲーションのリダイレクトを追跡する |
| `linux:WindowLoadCommitted` | WebView が読み込みを確定した | コンテンツを受信中 |
| `linux:WindowLoadFinished` | WebView の読み込みが完了した | 読み込みインジケーターを非表示にし、JS/CSS を挿入する |

## ベストプラクティス

### 1. イベント名前空間を使用する

カスタムイベントを作成するときは、競合を避けるために名前空間を使用します。

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

コンポーネントがアンマウントされるときは、必ずイベントリスナーを削除します。

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

プラットフォーム固有のイベントを使用するときは、そのプラットフォームで利用可能か確認します。

```javascript
import { Events } from '@wailsio/runtime';

// Platform-specific events can be registered unconditionally;
// they will simply never fire on unsupported platforms.
Events.On('windows:APMSuspend', handleSuspend);
Events.On('mac:ApplicationWillTerminate', handleTerminate);
```

### 4. イベントを多用しない

イベントは強力ですが、あらゆる処理に使用しないでください。

- ✅ イベントの使用に適しているもの：システム通知、ライフサイクルの変更、更新のブロードキャスト
- ❌ イベントの使用を避けるもの：関数からの直接の戻り値、単一コンポーネントの更新、同期処理

## イベントのデバッグ

イベントに関する問題をデバッグするには、次の手順を実行します。

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

- フロントエンドのイベント：[`v3/internal/runtime/desktop/@wailsio/runtime/src/event_types.ts`](https://github.com/wailsapp/wails/blob/main/v3/internal/runtime/desktop/@wailsio/runtime/src/event_types.ts)
- バックエンドのイベント：[`v3/pkg/events/events.go`](https://github.com/wailsapp/wails/blob/main/v3/pkg/events/events.go)

最新のイベント名と利用可否については、必ずこれらのファイルを参照してください。

## まとめ

Wails のイベントを使用すると、アプリケーション内の通信を強力かつ疎結合な方法で処理できます。このガイドのパターンとプラクティスに従うことで、システムの変更やユーザー操作に円滑に反応する、応答性が高くプラットフォームを考慮したアプリケーションを構築できます。

重要：クロスプラットフォーム互換性を確保するために、まず共通イベントを使用し、必要に応じてプラットフォーム固有のイベントを追加してください。また、メモリリークを防ぐため、イベントリスナーは必ずクリーンアップしてください。
