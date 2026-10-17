---
title: "ランタイムの内部構造"
description: "Wails v3 の起動、実行、OS との通信の仕組みを詳しく解説します"
slug: "contributing/runtime-internals"
sourcePath: "contributing/runtime-internals.md"
---

**runtime** は、通常の Go 関数をクロスプラットフォームの デスクトップアプリケーションへ変換するレイヤーです。 このドキュメントでは、ソースコードを追跡する際に目にする各構成要素について説明します。

---

## 1. アプリケーションのライフサイクル

| フェーズ | コードパス | 処理内容 |
| --- | --- | --- |
| **ブートストラップ** | `pkg/application/application.go:init()` | ビルド時のデータを登録し、グローバルな `application` シングルトンを作成します。 |
| **New()** | `application.New(...)` | `Options` を検証し、**AssetServer** を起動して、ロギングを初期化します。 |
| **Run()** | `application.(*App).Run()` | 1. プラットフォームの `mainthread.X()` を呼び出して、OS の UI スレッドに入ります。<br />2. **runtime**（`internal/runtime`）を起動します。<br />3. 最後のウィンドウが閉じるか、`Quit()` が呼び出されるまでブロックします。 |
| **シャットダウン** | `application.(*App).Quit()` | `application:shutdown` イベントをブロードキャストし、ログをフラッシュして、ウィンドウとサービスを終了します。 |

ライフサイクルは厳密に <strong>単一エントリ</strong>です。ウィンドウは複数作成できますが、 アプリケーションオブジェクト自体は一度だけ初期化されます。

---

## 2. ウィンドウ管理

### 公開 API

```go
win := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "Dashboard",
    Width:  1280,
    Height: 720,
})
win.Show()
```

> `app.Window.New()` は **引数を取りません**。`NewWithOptions(...)` を使用するのは、
>
> `application.WebviewWindowOptions` 構造体を値として渡す必要がある場合です。

`app.Window.New[WithOptions]()` は、プラットフォーム固有の実装が置かれている `pkg/application/webview_window_*.go` に処理を委譲します。

```
pkg/application/
├── webview_window_darwin.go    // WKWebView
├── webview_window_linux.go     // GTK + WebKitGTK (plus linux_cgo*.go)
└── webview_window_windows.go   // WebView2
```

各ファイルは、次の処理を行います。

1. ネイティブ WebView（WKWebView、WebKitGTK、WebView2）を作成します。
2. **Message Processor** コールバック（`pkg/application/messageprocessor*.go`）を登録します。
3. Wails イベント（`WindowDidResize`、`WindowFocus`、`WindowFilesDropped`、…）を `pkg/events` の定数にマッピングします。

`internal/runtime/` は、小規模なビルドタグ用グルーコード （`runtime{,_darwin,_linux,_windows,_android,_dev,_prod}.go`）と、`internal/runtime/desktop/` 配下に埋め込まれた JS ランタイム専用です。

アクティブなウィンドウは `pkg/application/window_manager.go` / `webview_window.go` によって追跡されます。`pkg/application/screenmanager.go` は <strong>ディスプレイ</strong>の メタデータ（解像度、スケール、作業領域）を扱うものであり、ウィンドウは管理しません。

---

## 3. メッセージ処理パイプライン

JavaScript と Go の間のブリッジは、`pkg/application/messageprocessor_*.go` にある **Message Processor** ファミリーによって実装されています。

処理フロー：

1. **JavaScript** は、`/wails/runtime.js` から `Call.ByID(<fnv-id>, ...args)`（`internal/runtime/desktop/@wailsio/runtime/src/calls.ts` で実装）を呼び出します。name-mode ビルドでは、代わりに `Call.ByName("pkg.Struct.Method", ...args)` を呼び出します。
2. ランタイムヘルパーは呼び出しをパッケージ化し、プラットフォームごとのネイティブブリッジを介して Go にディスパッチします。
3. **Go** は `pkg/application/messageprocessor_call.go` でメッセージを受信します。
4. プロセッサーは `pkg/application/bindings.go`（手書きの `reflect` ベースの実装）からバインド済みメソッドを検索し、呼び出します。
5. 結果またはエラーは JS にマーシャリングされて返され、そこで `Promise` が解決または拒否されます。

> JSON エンベロープの正確な形式は、JS 側ではランタイムヘルパーによって、
>
> Go 側では `messageprocessor_call.go` によって定義されています。このページの古い草稿では
>
> `{"t":"c","id":"123","m":"Greet","p":[…]}` という形式が記載されていましたが、これは
>
> 現在の実装と一致しません。ワイヤーフォーマットのバグを追跡する際は、両方のファイルを
>
> 併せて確認してください。

専用プロセッサー：

| ファイル | 用途 |
| --- | --- |
| `messageprocessor_window.go` | ウィンドウ操作（非表示、最大化、…） |
| `messageprocessor_dialog.go` | ネイティブダイアログ（`OpenFile`、`MessageBox`、…） |
| `messageprocessor_clipboard.go` | クリップボードの読み書き |
| `messageprocessor_events.go` | イベントの購読／発行 |
| `messageprocessor_browser.go` | ブラウザーのナビゲーション、開発者ツール |

プロセッサーは <strong>ステートレス</strong>です。必要な情報はすべて、各メッセージとともに 渡される `ApplicationContext` から取得します。

---

## 4. イベントシステム

イベントは名前空間付きの文字列であり、次の 3 つのレイヤーにわたってディスパッチされます。

1. **アプリケーションイベント**：グローバルライフサイクル（`application:ready`、`application:shutdown`）。
2. **ウィンドウイベント**：ウィンドウ単位（`window:focus`、`window:resize`）。
3. **カスタムイベント**：ユーザー定義（`chat:new-message`）。

実装の詳細：

- イベント定数は`pkg/events/`（`defaults.go`、`known_events.go`、`events.txt`）にあります。これらは`v3/tasks/events/generate.go`によって生成され、`events.Common.*`、`events.Mac.*`、`events.Windows.*`、`events.Linux.*`として公開されます。`wails3 generate constants`で再生成できます。
- Go 側（アプリケーションイベント）：
  ```go
  app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {})
  ```

- Go 側（ウィンドウイベント）：
  ```go
  window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {})
  ```

- Go 側（カスタムイベント）：
  ```go
  app.Event.On("chat:new-message", func(e *application.CustomEvent) {})
  ```

- JS 側：
  ```js
  import { Events } from "/wails/runtime.js";
  Events.On("chat:new-message", (e) => { /* … */ });
  ```


アプリケーション／ウィンドウ／カスタムの各イベントは、すべて `pkg/application/event_manager.go`を経由します。ウィンドウイベントのサブスクリプションは そのウィンドウに限定されるため、ウィンドウを閉じるとハンドラーの登録は自動的に解除されます。

---

## 5. プラットフォーム固有の実装

条件付きコンパイルにより、OS 固有の複雑さを隠しながら、公開 API を同一に保ちます。

| 項目 | Darwin | Linux | Windows |
| --- | --- | --- | --- |
| メインスレッド | `mainthread_darwin.go`（Cgo から Foundation を使用） | `mainthread_linux.go`（GTK） | `mainthread_windows.go`（Win32 `AttachThreadInput`） |
| ダイアログ | `dialogs_darwin.*`（NSAlert） | `dialogs_linux.go`（GtkFileChooser） | `dialogs_windows.go`（IFileOpenDialog） |
| クリップボード | `clipboard_darwin.go` | `clipboard_linux.go` | `clipboard_windows.go` |
| トレイアイコン | `systemtray_darwin.*` | `systemtray_linux.go`（DBus） | `systemtray_windows.go`（Shell_NotifyIcon） |

主な原則：

- <strong>macOS と Windows</strong>では Cgo の使用を最小限に抑えています（主に`pkg/mac/`と、`pkg/w32`内の`w32` Win32 ラッパーを介して使用します）。
- **Linux では必要上、Cgo を多用します**。`pkg/application/linux_cgo.go`（約69 KB）と`linux_cgo_gtk4.{c,go,h}`（約50 KB 以上）から GTK/WebKitGTK を直接操作します。
- OS ごとのファイルを読みやすく保つには、**ビルドタグ**（`//go:build darwin`、`//go:build linux`、…）を使用します。
- `internal/capabilities/`はプラットフォームごとの機能フラグ用に存在しますが、フレームワークは`ErrCapability`センチネルをエクスポート<strong>しません</strong>。機能の制御は、プラットフォーム固有のスタブが返す値によって行われます。

---

## 6. ファイルガイド

| ファイル | 変更する目的 |
| --- | --- |
| `internal/runtime/runtime_*.go` | 小規模なビルドタグ付きスタブ層（開発環境と本番環境の差異、OS 固有の連携部分）を変更します。 |
| `pkg/application/webview_window_*.go` | 新しいウィンドウヒントまたは動作を実装します。 |
| `pkg/application/messageprocessor*.go` | JS から呼び出せる新しいブリッジコマンドを追加します。 |
| `pkg/events/*.go` | 組み込みのイベント定義を拡張します（その後、`wails3 generate constants`を再実行します）。 |
| `internal/assetserver/*` | 開発環境／本番環境でのアセット処理を調整します。 |
| `internal/runtime/desktop/@wailsio/runtime/src/*` | 組み込みの JS ランタイム（呼び出し／イベントのディスパッチ、ダイアログ、ドラッグなど）を編集します。 |

---

## 7. デバッグのヒント

- `Options.LogLevel`（例：`slog.LevelDebug`）を設定し、`Options.Logger`の出力を確認します。`WAILS_LOG_LEVEL`環境変数はありません。
- `wails3 dev`のフラグは、`--config`、`--port`、`-s`（HTTPS を有効化）、およびグローバルな`--no-colour`です。`-verbose`フラグはありません。
- macOS では、Objective-C の例外を早期に検出するため、`lldb --`の下で実行します。
- Windows で Chromium の問題を調査する場合は、WebView2 のデバッグログを有効にします：`set WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS=--remote-debugging-port=9222`

---

## 8. ランタイムの拡張

1. 任意：新しい機能フラグがあれば、`internal/capabilities/`で宣言します。
2. ビルドタグを使用して、`pkg/application/*_{darwin,linux,windows}.go`の各バリアントに機能を実装します。サポートしないプラットフォームにはスタブを用意します。
3. `pkg/application`に公開 API（インターフェースと具象`WebviewWindow`メソッド、オプション構造体など）を追加します。
4. JS から呼び出す必要がある場合は、新しいメッセージプロセッサーメソッド（`pkg/application/messageprocessor*.go`）を登録し、対応するヘルパーを JS ランタイムに登録します。
5. イベントを追加する場合は、その定数を `pkg/events/` で宣言し、`wails3 generate constants` を実行して生成済みファイルを更新します。

このチェックリストに従えば、クロスプラットフォームの契約を維持できます。

---

## 9. ドラッグ＆ドロップ

ファイルのドラッグ＆ドロップには、すべてのプラットフォームで <strong>JavaScript ファーストの方式</strong>を使用します。ネイティブ層が OS のドラッグイベントをインターセプトしますが、実際のドロップ処理と DOM の操作は JavaScript で行われます。

### 処理フロー

1. ユーザーが OS から Wails ウィンドウ上へファイルをドラッグする
2. ネイティブ層がドラッグを検出し、ホバー効果のために JavaScript へ通知する
3. ユーザーがファイルをドロップする
4. ネイティブ層がファイルパスと座標を JavaScript へ送信する
5. JavaScript がドロップ先の要素（`data-file-drop-target`）を特定する
6. JavaScript がファイルパスと要素の詳細を Go バックエンドへ送信する
7. Go が完全なコンテキストを含む `WindowFilesDropped` イベントを発行する

### プラットフォーム別の実装

| プラットフォーム | ネイティブ層 | 主な課題 |
| --- | --- | --- |
| **Windows** | WebView2 の組み込みドラッグサポート | 座標は CSS ピクセル単位のため、変換は不要 |
| **macOS** | NSWindow のドラッグデリゲート | ウィンドウ相対座標を WebView 相対座標へ変換 |
| **Linux** | GTK3 のドラッグシグナル | ファイルのドラッグと内部の HTML5 ドラッグを区別する必要がある |

### Linux：ドラッグ種別の判別

GTK と WebKit は、どちらもドラッグイベントを処理しようとします。重要なのは、ドラッグ対象の種類を確認することです。

```c
static gboolean is_file_drag(GdkDragContext *context) {
    GList *targets = gdk_drag_context_list_targets(context);
    for (GList *l = targets; l != NULL; l = l->next) {
        GdkAtom atom = GDK_POINTER_TO_ATOM(l->data);
        gchar *name = gdk_atom_name(atom);
        if (name && g_strcmp0(name, "text/uri-list") == 0) {
            g_free(name);
            return TRUE;  // External file drag
        }
        g_free(name);
    }
    return FALSE;  // Internal HTML5 drag
}
```

シグナルハンドラーは、内部ドラッグの場合は `FALSE` を返して WebKit に処理させ、ファイルのドラッグの場合は `TRUE` を返して独自に処理します。

### ファイルドロップのブロック

`EnableFileDrop` が `false` の場合でも、ドロップされたファイルへブラウザーが移動しないようにする必要があります。各プラットフォームでは、これを次のように処理します。

- **Windows**：ドラッグイベントで JavaScript が `preventDefault()` を呼び出す
- **macOS**：ドラッグイベントで JavaScript が `preventDefault()` を呼び出す\
- **Linux**：GTK シグナルハンドラーがネイティブ層でファイルのドラッグをインターセプトして拒否する

### 主要ファイル

| ファイル | 用途 |
| --- | --- |
| `pkg/application/linux_cgo.go` | GTK ドラッグシグナルハンドラー（cgo プリアンブル内の C コード） |
| `pkg/application/webview_window_darwin.go` | macOS のドラッグデリゲート |
| `pkg/application/webview_window_windows.go` | WebView2 のメッセージ処理 |
| `internal/runtime/desktop/@wailsio/runtime/src/window.ts` | JavaScript のドロップ処理 |

### デバッグ

- **Linux**：C コードに `printf` を追加する（`fflush(stdout)` を忘れないこと）
- **Windows**：`globalApplication.debug()` を使用する
- **JavaScript**：ブラウザーコンソールを確認し、デバッグモードを有効にする

よくある問題：

1. **内部の HTML5 ドラッグが機能しない**：ネイティブハンドラーがインターセプトしている（ファイル以外のドラッグでは `FALSE` を返す）
2. **ホバー効果が表示されない**：JavaScript ハンドラーが呼び出されていない
3. **座標が正しくない**：座標空間の変換を確認する

---

これで、ランタイム内部のガイドツアーは完了です。この知識を <strong>コードベースの構成</strong>マップおよび <strong>アセットサーバー</strong>のドキュメントと組み合わせれば、迷わずコードベースを把握し、 大きな効果をもたらす貢献ができます。コーディングを楽しんでください！
