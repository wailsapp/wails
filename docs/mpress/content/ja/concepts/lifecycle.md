---
title: "アプリケーションのライフサイクル"
description: "起動からシャットダウンまでの Wails アプリケーションのライフサイクルを理解する"
slug: "concepts/lifecycle"
sourcePath: "concepts/lifecycle.md"
---

## アプリケーションのライフサイクルを理解する

デスクトップアプリケーションには、起動からシャットダウンまでのライフサイクルがあります。Wails v3 では、このライフサイクルを効果的に管理するために、**サービス**、**イベント**、<strong>フック</strong>が提供されています。

## ライフサイクルの各段階

```d2
direction: down

Start: アプリケーション開始 {
  shape: oval
  style.fill: "#10B981"
}

Init: 初期化 {
  Parse: オプションを解析 {
    shape: rectangle
  }
  Register: サービスを登録 {
    shape: rectangle
  }
  Setup: ランタイムをセットアップ {
    shape: rectangle
  }
}

AppRun: app.Run() {
  shape: rectangle
  style.fill: "#3B82F6"
}

ServiceStartup: サービスを起動 {
  shape: rectangle
  style.fill: "#8B5CF6"
}

EventLoop: イベントループ {
  Process: イベントを処理 {
    shape: rectangle
  }
  Handle: メッセージを処理 {
    shape: rectangle
  }
  Update: UIを更新 {
    shape: rectangle
  }
}

QuitSignal: 終了シグナル {
  shape: diamond
  style.fill: "#F59E0B"
}

ShouldQuit: ShouldQuitを確認 {
  shape: rectangle
  style.fill: "#3B82F6"
}

OnShutdown: OnShutdownコールバック {
  shape: rectangle
  style.fill: "#3B82F6"
}

ServiceShutdown: サービスを停止 {
  shape: rectangle
  style.fill: "#8B5CF6"
}

Cleanup: クリーンアップ {
  Close: ウィンドウを閉じる {
    shape: rectangle
  }
  Release: リソースを解放 {
    shape: rectangle
  }
}

End: アプリケーション終了 {
  shape: oval
  style.fill: "#EF4444"
}

Start -> Init.Parse
Init.Parse -> Init.Register
Init.Register -> Init.Setup
Init.Setup -> AppRun
AppRun -> ServiceStartup
ServiceStartup -> EventLoop.Process
EventLoop.Process -> EventLoop.Handle
EventLoop.Handle -> EventLoop.Update
EventLoop.Update -> EventLoop.Process: ループ
EventLoop.Process -> QuitSignal: ユーザーが終了
QuitSignal -> ShouldQuit: 許可されているか確認
ShouldQuit -> EventLoop.Process: 拒否
ShouldQuit -> OnShutdown: 許可
OnShutdown -> ServiceShutdown
ServiceShutdown -> Cleanup.Close
Cleanup.Close -> Cleanup.Release
Cleanup.Release -> End
```

### 1. アプリケーションの作成

`application.New()`を使用してアプリケーションを作成します。

```go
app := application.New(application.Options{
    Name:        "My App",
    Description: "An application built with Wails",
    Services: []application.Service{
        application.NewService(&MyService{}),
    },
    Assets: application.AssetOptions{
        Handler: application.BundledAssetFileServer(assets),
    },
})
```

**実行される処理：**

1. オプションが解析および検証される
2. サービスが登録される（この時点ではまだ起動されない）
3. アセットサーバーが構成される
4. ランタイムがセットアップされる

### 2. アプリケーションの実行

`app.Run()`を呼び出してアプリケーションを起動します。

```go
err := app.Run()  // Blocks until quit
if err != nil {
    log.Fatal(err)
}
```

**実行される処理：**

1. サービスが登録順に起動される
2. イベントリスナーが有効化される
3. ウィンドウを作成できるようになる
4. イベントループが開始される

### 3. イベントループ

アプリケーションはイベントループに入り、実行時間の大半をそこで費やします。

- OS イベント（マウス、キーボード、ウィンドウのイベント）が処理される
- Go から JS へのメッセージが処理される
- JS から Go への呼び出しが実行される
- UI の更新がレンダリングされる

### 4. シャットダウン

アプリケーションが終了するときは、次の処理が行われます。

1. `ShouldQuit`コールバックが設定されている場合、その結果が確認される
2. `OnShutdown`コールバックが実行される
3. サービスが逆順にシャットダウンされる
4. ウィンドウが閉じられる
5. リソースが解放される

## サービスのライフサイクル

Wails v3 では、ライフサイクルの主な管理手段としてサービスを使用します。サービスは、インターフェースを通じて起動フックとシャットダウンフックを提供します。サービスの詳細については、[サービスガイド](/features/bindings/services/)を参照してください。

### サービスの作成

```go
type MyService struct {
    db *sql.DB
}

// ServiceStartup is called when the application starts
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    var err error
    s.db, err = sql.Open("sqlite3", "app.db")
    if err != nil {
        return err  // Startup aborts if error returned
    }

    // Run migrations
    if err := s.runMigrations(); err != nil {
        return err
    }

    return nil
}

// ServiceShutdown is called when the application shuts down
func (s *MyService) ServiceShutdown() error {
    if s.db != nil {
        return s.db.Close()
    }
    return nil
}
```

### サービスの登録

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&MyService{}),
        application.NewService(&AnotherService{}),
    },
})
```

**要点：**

- サービスは登録順に起動される
- サービスは登録時とは<strong>逆</strong>の順序でシャットダウンされる
- サービスの`ServiceStartup`がエラーを返した場合、アプリケーションは中止される
- `ServiceStartup`に渡された`ctx`は、シャットダウンの開始時にキャンセルされる

### アプリケーションコンテキストの使用

`ServiceStartup` に渡されるコンテキストは、アプリケーションの存続期間中有効です。

```go
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    // Start a background task that respects shutdown
    go func() {
        ticker := time.NewTicker(5 * time.Minute)
        defer ticker.Stop()

        for {
            select {
            case <-ticker.C:
                s.performBackgroundSync()
            case <-ctx.Done():
                // Application is shutting down
                return
            }
        }
    }()

    return nil
}
```

アプリケーションインスタンスからコンテキストにアクセスすることもできます。

```go
app := application.Get()
ctx := app.Context()
```

## アプリケーションレベルのフック

これらは `application.Options` の便利なコールバックであり、完全なサービスを作成せずにアプリケーションのライフサイクルへフックできます。単純なクリーンアップ処理、終了確認、またはシャットダウンシーケンスの特定時点でコードを実行する必要がある場合に役立ちます。

起動ロジック、依存性注入、ステートフルなリソースを伴う、より複雑なライフサイクル管理には、代わりに [サービス](#heading-2) を使用してください。

### ShouldQuit

終了が要求されるたびに `ShouldQuit` コールバックが呼び出されます。これには、ユーザーが最後のウィンドウを閉じた場合、Cmd+Q（macOS）／Alt+F4（Windows）を押した場合、またはプログラムから `app.Quit()` を呼び出した場合が含まれます。

**戻り値：**

- 終了処理の続行を許可するには `true` を返します（アプリケーションはシャットダウンします）
- 終了をキャンセルするには `false` を返します（アプリケーションは実行を継続します）

ここで終了要求をインターセプトし、必要に応じて終了を阻止できます。たとえば、未保存の変更についてユーザーに確認できます。

```go
app := application.New(application.Options{
    ShouldQuit: func() bool {
        if !hasUnsavedChanges() {
            return true // No unsaved changes, allow quit
        }

        // Prompt the user — MessageDialog.Show() blocks and returns nothing.
        // The button's OnClick callback fires for whichever button the user picks.
        shouldQuit := false

        dlg := application.Get().Dialog.Question().
            SetTitle("Unsaved Changes").
            SetMessage("You have unsaved changes. Quit anyway?")

        quit := dlg.AddButton("Quit")
        cancel := dlg.AddButton("Cancel")
        dlg.SetDefaultButton(cancel)
        dlg.SetCancelButton(cancel)

        quit.OnClick(func() { shouldQuit = true })

        dlg.Show()
        return shouldQuit
    },
})
```

`ShouldQuit` が設定されていない場合、終了が要求されるとアプリケーションは直ちに終了します。

**ShouldQuit が呼び出される場合：**

- ユーザーが最後のウィンドウを閉じた場合（`DisableQuitOnLastWindowClosed` が設定されている場合を除く）
- ユーザーが macOS で Cmd+Q を押した場合
- ユーザーが Windows で Alt+F4 を押した場合（最後のウィンドウにフォーカスがあるとき）
- コードから `app.Quit()` を呼び出した場合

**ShouldQuit が呼び出されない場合：**

- プロセスが強制終了された場合（SIGKILL、タスクマネージャーによる強制終了）
- `os.Exit()` が直接呼び出された場合

### OnShutdown

アプリケーションの終了が確定すると、`OnShutdown` コールバックが呼び出されます（`ShouldQuit` が設定されている場合は、それが `true` を返した後）。状態の保存、データベース接続の切断、リソースの解放などのクリーンアップ処理に使用してください。

```go
app := application.New(application.Options{
    OnShutdown: func() {
        // Save application state
        saveState()

        // Close connections
        cleanup()
    },
})
```

アプリケーションのライフサイクル中はいつでも、プログラムから追加のシャットダウンコールバックを登録できます。

```go
app.OnShutdown(func() {
    log.Println("Application shutting down...")
})
```

複数のコールバックは登録順に実行されます。すべてのコールバックが完了するまで、シャットダウン処理はブロックされます。

**重要：** シャットダウンコールバックは短時間（1 秒未満）で完了させてください。終了に時間がかかりすぎるアプリケーションは、オペレーティングシステムによって強制終了される可能性があります。その場合、クリーンアップが中断され、データが失われるおそれがあります。

### PostShutdown

すべてのシャットダウン処理が完了した後、プロセスが終了する直前に `PostShutdown` コールバックが呼び出されます。この時点では、アプリケーションインスタンスは使用できません。ウィンドウは閉じられ、サービスはシャットダウンされ、リソースは解放されています。

これは主に次の用途に役立ちます。

- ほかのすべてのクリーンアップ後に行う必要がある最終ログ記録
- シャットダウン動作のテストとデバッグ
- `app.Run()` が戻らないプラットフォーム（このコールバックによりコードの実行が保証されます）

```go
app := application.New(application.Options{
    PostShutdown: func() {
        // Final logging
        log.Println("Application terminated cleanly")

        // Flush any buffered logs
        logger.Sync()
    },
})
```

**注記：** `PostShutdown` では、アプリケーション機能（ウィンドウやダイアログなど）を使用しないでください。これらはすでに利用できません。

## イベントベースのライフサイクル

Wails は、ウィンドウが開いたとき、アプリケーションが起動したとき、テーマが変更されたときなど、アプリケーション内で何かが発生した際に通知するイベントシステムを提供します。これらのイベントをリッスンすることで、処理をブロックしたりイベントをインターセプトしたりせずに、ライフサイクルの変化に対応できます。

ウィンドウイベントでは、`OnWindowEvent` の代わりに `RegisterHook` を使用してアクションをインターセプトし、キャンセルすることもできます。たとえば、ウィンドウが閉じるのを阻止できます。以下の「[ウィンドウフック](#heading-10)」を参照してください。

イベントシステムの完全なドキュメントについては、[イベントガイド](/features/events/system/)を参照してください。

### アプリケーションイベント

アプリケーションのライフサイクルイベントをリッスンします。

```go
app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(event *application.ApplicationEvent) {
    app.Logger.Info("Application has started!")
})
```

プラットフォーム固有のイベントも利用できます。

```go
// macOS
app.Event.OnApplicationEvent(events.Mac.ApplicationDidFinishLaunching, func(event *application.ApplicationEvent) {
    // Handle macOS launch
})

app.Event.OnApplicationEvent(events.Mac.ApplicationWillTerminate, func(event *application.ApplicationEvent) {
    // Handle macOS termination
})

// Windows
app.Event.OnApplicationEvent(events.Windows.ApplicationStarted, func(event *application.ApplicationEvent) {
    // Handle Windows start
})
```

### ウィンドウイベント

ウィンドウのライフサイクルイベントをリッスンします。

```go
window := app.Window.New()

window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window gained focus")
})

window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
    app.Logger.Info("Window is closing")
})
```

### ウィンドウフック（キャンセル可能なイベント）

イベントを<strong>キャンセル</strong>する必要がある場合は、`OnWindowEvent` の代わりに `RegisterHook` を使用します。

```go
window := app.Window.New()

var countdown = 3

window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    countdown--
    if countdown > 0 {
        app.Logger.Info("Not closing yet!", "remaining", countdown)
        e.Cancel()  // Prevent the window from closing
        return
    }
    app.Logger.Info("Window closing now")
})
```

**OnWindowEvent と RegisterHook の違い：**

- `OnWindowEvent`：イベントが発生したときに通知します（キャンセル不可）
- `RegisterHook`：イベントをインターセプトし、必要に応じてキャンセルできます

## ウィンドウのライフサイクル

ウィンドウには、作成から破棄まで独自のライフサイクルがあります。各ウィンドウはフロントエンドのコンテンツを個別に読み込み、いつでも表示、非表示、または閉じることができます。ユーザーがウィンドウを閉じようとしたときは、`RegisterHook`でその操作をインターセプトし、確認を求めたり、ウィンドウを破棄せずに非表示にしたりできます。

ウィンドウに関する完全なドキュメントについては、[ウィンドウガイド](/features/windows/basics/)を参照してください。

```d2
direction: down

Create: ウィンドウを作成 {
  shape: oval
  style.fill: "#10B981"
}

Load: フロントエンドを読み込む {
  shape: rectangle
}

Show: ウィンドウを表示 {
  shape: rectangle
}

Active: ウィンドウがアクティブ {
  Events: イベントを処理 {
    shape: rectangle
  }
}

CloseRequest: 閉じるリクエスト {
  shape: diamond
  style.fill: "#F59E0B"
}

Hook: WindowClosingフック {
  shape: rectangle
  style.fill: "#3B82F6"
}

Destroy: ウィンドウを破棄 {
  shape: rectangle
}

End: ウィンドウが閉じられた {
  shape: oval
  style.fill: "#EF4444"
}

Create -> Load
Load -> Show
Show -> Active.Events
Active.Events -> Active.Events: ループ
Active.Events -> CloseRequest: ユーザーが閉じる
CloseRequest -> Hook
Hook -> Active.Events: キャンセル
Hook -> Destroy: 許可
Destroy -> End
```

### ウィンドウの作成

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "My Window",
    Width:  800,
    Height: 600,
})
```

### ウィンドウが閉じられるのを防ぐ

```go
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    if !hasUnsavedChanges() {
        return
    }

    // MessageDialog.Show() returns nothing; per-button OnClick handlers fire.
    dlg := application.Get().Dialog.Question().
        SetTitle("Unsaved Changes").
        SetMessage("Save before closing?")

    save := dlg.AddButton("Save")
    discard := dlg.AddButton("Discard")
    cancel := dlg.AddButton("Cancel")
    dlg.SetDefaultButton(save)
    dlg.SetCancelButton(cancel)

    save.OnClick(func() { saveChanges() })
    cancel.OnClick(func() { e.Cancel() }) // Prevent close
    _ = discard                            // "Discard" falls through and allows close

    dlg.Show()
})
```

### 閉じる代わりに非表示にする

システムトレイアプリで一般的なパターンは次のとおりです。

```go
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    window.Hide()  // Hide instead of destroy
    e.Cancel()     // Prevent actual close
})
```

## 複数ウィンドウのライフサイクル

複数のウィンドウがある場合：

```go
mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Main Window",
})

settingsWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "Settings",
    Width:  400,
    Height: 600,
    Hidden: true,  // Start hidden
})
```

**デフォルトの動作はプラットフォームによって異なります：**

| プラットフォーム | 最後のウィンドウを閉じたときのデフォルト動作 |
| --- | --- |
| macOS | アプリは実行を継続する（メニューバーは残る） |
| Windows | アプリが終了する |
| Linux | アプリが終了する |

macOS はプラットフォーム固有の慣例に従い、ウィンドウがなくても通常はアプリケーションがメニューバーでアクティブな状態を維持します。Windows と Linux ではデフォルトで終了します。

**すべてのプラットフォームで、最後のウィンドウを閉じたときに終了させる：**

```go
app := application.New(application.Options{
    Mac: application.MacOptions{
        ApplicationShouldTerminateAfterLastWindowClosed: true,
    },
})
```

**すべてのプラットフォームで、最後のウィンドウを閉じた後も実行を継続させる：**

これは、システムトレイアプリケーションや、バックグラウンドで実行を継続する必要があるアプリに役立ちます。

```go
app := application.New(application.Options{
    Windows: application.WindowsOptions{
        DisableQuitOnLastWindowClosed: true,
    },
    Linux: application.LinuxOptions{
        DisableQuitOnLastWindowClosed: true,
    },
})
```

## 一般的なパターン

### パターン 1：データベースサービス

```go
type DatabaseService struct {
    db *sql.DB
}

func (s *DatabaseService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    var err error
    s.db, err = sql.Open("sqlite3", "app.db")
    if err != nil {
        return fmt.Errorf("failed to open database: %w", err)
    }

    if err := s.db.PingContext(ctx); err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }

    return nil
}

func (s *DatabaseService) ServiceShutdown() error {
    if s.db != nil {
        return s.db.Close()
    }
    return nil
}

// Exported methods are available to the frontend
func (s *DatabaseService) GetUsers() ([]User, error) {
    // Query implementation
}
```

### パターン 2：構成サービス

```go
type ConfigService struct {
    config *Config
    path   string
}

func (s *ConfigService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    s.path = "config.json"

    data, err := os.ReadFile(s.path)
    if err != nil {
        if os.IsNotExist(err) {
            s.config = &Config{} // Default config
            return nil
        }
        return err
    }

    return json.Unmarshal(data, &s.config)
}

func (s *ConfigService) ServiceShutdown() error {
    data, err := json.MarshalIndent(s.config, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(s.path, data, 0644)
}
```

### パターン 3：バックグラウンドワーカー

```go
type WorkerService struct {
    cancel context.CancelFunc
}

func (s *WorkerService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    workerCtx, cancel := context.WithCancel(ctx)
    s.cancel = cancel

    go s.runWorker(workerCtx)

    return nil
}

func (s *WorkerService) ServiceShutdown() error {
    if s.cancel != nil {
        s.cancel()
    }
    return nil
}

func (s *WorkerService) runWorker(ctx context.Context) {
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            s.doWork()
        case <-ctx.Done():
            return
        }
    }
}
```

## ライフサイクルリファレンス

| フック／インターフェース | 呼び出されるタイミング | キャンセル可能か | 用途 |
| --- | --- | --- | --- |
| `ServiceStartup` | `app.Run()`の実行中（イベントループの開始前） | 不可（エラーを返すと中止） | 初期化 |
| `ServiceShutdown` | シャットダウン中（`OnShutdown`の後） | 不可 | クリーンアップ |
| `OnShutdown` | 終了が確定したとき | 不可 | アプリケーションのクリーンアップ |
| `ShouldQuit` | 終了が要求されたとき | 可（false を返す） | 終了の確認 |
| `RegisterHook(WindowClosing)` | ウィンドウを閉じるよう要求されたとき | 可（`e.Cancel()`） | ウィンドウが閉じられるのを防ぐ |
| `OnWindowEvent` | イベントが発生したとき | 不可 | イベントへの応答 |
| `OnApplicationEvent` | イベントが発生したとき | 不可 | イベントへの応答 |

## プラットフォームによる違い

### macOS

- ウィンドウがなくても<strong>アプリケーションメニュー</strong>は表示され続ける
- <strong>Cmd+Q</strong>で終了がトリガーされる（`ShouldQuit`を経由する）
- 非表示にしない限り、<strong>Dockアイコン</strong>は表示され続ける
- 終了動作を制御するには`ApplicationShouldTerminateAfterLastWindowClosed`を使用する

### Windows

- ウィンドウがない場合、**アプリケーションメニューは表示されない**
- <strong>Alt+F4</strong>でウィンドウが閉じる（`RegisterHook`で防止可能）
- <strong>システムトレイ</strong>を使用すると、アプリケーションを実行し続けることができる

### Linux

- デスクトップ環境によって<strong>動作が異なる</strong>
- **通常はWindowsと同様**

## ライフサイクルの問題をデバッグする

### 問題：アプリケーションが終了しない

**原因：**

1. `ShouldQuit`が`false`を返している
2. `OnShutdown`に時間がかかりすぎている
3. バックグラウンドのgoroutineが停止していない

**解決策：**

```go
// 1. Check ShouldQuit logic
ShouldQuit: func() bool {
    log.Println("ShouldQuit called")
    return true
}

// 2. Keep OnShutdown fast
OnShutdown: func() {
    log.Println("OnShutdown started")
    // Fast cleanup only
    log.Println("OnShutdown finished")
}

// 3. Use context for background tasks
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    go func() {
        <-ctx.Done()
        log.Println("Context cancelled, stopping background work")
    }()
    return nil
}
```

### 問題：サービスの起動に失敗する

**解決策：** 内容が分かるエラーを返す：

```go
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    if err := s.init(); err != nil {
        return fmt.Errorf("failed to initialise: %w", err)
    }
    return nil
}
```

エラーがログに記録され、アプリケーションは起動しません。

## ベストプラクティス

### 推奨事項

- **ライフサイクル管理にはサービスを使用する** — 適切な起動・終了フックが提供される
- **終了処理を高速に保つ** — すべてのクリーンアップを1秒未満で完了することを目標にする
- **キャンセルにはcontextを使用する** — バックグラウンドタスクを適切に停止する
- **起動時のエラーを処理する** — エラーを返して正常に中止する
- **ライフサイクルイベントをログに記録する** — デバッグに役立つ

### 禁止事項

- **サービスの起動処理をブロックしない** — 初期化を高速に保つ（2秒未満）
- **終了処理中にダイアログを表示しない** — アプリケーションは終了中であり、UIが機能しない可能性がある
- **contextを無視しない** — goroutine内では常に`ctx.Done()`を確認する
- **リソースをリークさせない** — 必ず`ServiceShutdown`を実装する

## 次のステップ

**サービス** — サービスシステムについて詳しく学ぶ [詳細を見る →](/features/bindings/services/)

**イベントシステム** — 通信にイベントを使用する [詳細を見る →](/features/events/system/)

**ウィンドウ管理** — ウィンドウを作成および管理する [詳細を見る →](/features/windows/basics/)

---

**ライフサイクルについて質問がありますか？** [Discord](https://discord.gg/JDdSxwjhGf)で質問するか、[サンプル](https://github.com/wailsapp/wails/tree/master/v3/examples)を確認してください。
