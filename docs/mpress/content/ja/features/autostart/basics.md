---
title: "自動起動"
description: "macOS、Windows、Linux でユーザーのログイン時にアプリケーションが起動するよう登録します"
slug: "features/autostart/basics"
sourcePath: "features/autostart/basics.md"
---

## 自動起動

`app.Autostart` は、ユーザーのログイン時にアプリケーションが自動的に起動するよう登録します。プラットフォームごとに適切なネイティブ機構を選択し、シンボリックリンクされたインストールパス（Homebrew、Scoop）を解決するため、バイナリをアップグレードしても登録が壊れません。

登録は直ちにではなく、**次回のログイン**時に有効になります。

## クイックスタート

```go
import "github.com/wailsapp/wails/v3/pkg/application"

// Register to launch at login
if err := app.Autostart.Enable(); err != nil {
    app.Logger.Error("autostart enable failed", "error", err)
}

// Stop launching at login
if err := app.Autostart.Disable(); err != nil {
    app.Logger.Error("autostart disable failed", "error", err)
}

// Check status
enabled, err := app.Autostart.IsEnabled()
```

## API

### `Enable`

デフォルトのオプションを使用して、ログイン時にアプリケーションが起動するよう登録します。

```go
func (m *AutostartManager) Enable() error
```

`Enable` を繰り返し呼び出しても安全です。呼び出すたびに登録が上書きされるため、ユーザーの設定を永続化している場合は、起動のたびに呼び出しても問題ありません。

### `EnableWithOptions`

カスタムオプションを使用して登録します。

```go
func (m *AutostartManager) EnableWithOptions(opts AutostartOptions) error
```

**`AutostartOptions`：**

| フィールド | 型 | 説明 |
| --- | --- | --- |
| `Identifier` | `string` | 自動的に導出された登録 ID を上書きします。後述の「識別子」を参照してください。 |
| `Arguments` | `[]string` | ログイン時の起動において、実行ファイルのパスに追加される引数（例：`--hidden`）。 |

### `Disable`

自動起動の登録を削除します。アプリケーションが登録されていなかった場合は `nil` を返します。無効化操作は冪等です。

```go
func (m *AutostartManager) Disable() error
```

### `IsEnabled`

登録が存在するかどうかを返します。高速ですが、登録されたパスは検証しません。

```go
func (m *AutostartManager) IsEnabled() (bool, error)
```

### `Status`

登録状態全体を返します。

```go
func (m *AutostartManager) Status() (AutostartStatus, error)
```

**`AutostartStatus`：**

| フィールド | 型 | 説明 |
| --- | --- | --- |
| `Enabled` | `bool` | 登録が存在するかどうか。 |
| `Path` | `string` | 登録アーティファクトのディスク上の場所（plist のパス、`.desktop` のパス、またはレジストリのサブキー）。`Enabled` が false の場合は空です。 |
| `Strategy` | `AutostartStrategy` | アプリケーションの登録に使用された機構（[プラットフォームごとの動作](#heading-2)を参照）。 |

## プラットフォームごとの動作

@tabs{sync-key="platform"}
[macOS]
アプリケーションのパッケージ方法に応じて、次の 2 つの機構を使い分けます。

- **macOS 13 以降、バンドルされた `.app`**：`SMAppService.mainAppService`。サンドボックス化されたアプリケーションおよび Mac App Store ビルドで動作します。TCC のオートメーション許可プロンプトは表示されません（従来の AppleScript を使用する方法では表示されていました）。
- **13 より前の macOS、またはバンドルされていないバイナリ**：`RunAtLoad=true` を指定した LaunchAgent plist を `~/Library/LaunchAgents/<identifier>.plist` に書き込みます。

呼び出し元がどちらの経路を使用したか判別できるように、`Status()` は `AutostartStrategySMAppService` または `AutostartStrategyLaunchAgent` を返します。

アプリケーションをバンドルされていない形式からバンドル形式へアップグレードした場合、`Status()` は両方の経路を確認し、`Disable()` は両方をクリーンアップします。これにより、孤立した LaunchAgent が古いビルドを起動し続けることを防ぎます。

[Windows]
`HKCU\Software\Microsoft\Windows\CurrentVersion\Run` の下にレジストリ値を追加します。値の名前には自動起動の `Identifier` を設定し、データには引用符で囲んだ実行ファイルのパスと `Arguments`（指定されている場合）を設定します。

引数の引用符処理は `CommandLineToArgvW` の規則（引用符の前にあるバックスラッシュを二重化）に従うため、空白や引用符を含むパスも正しく往復変換できます。

`Status().Strategy` は `AutostartStrategyRegistryRun` を返します。

[Linux]
次の内容を含む XDG 自動起動エントリを `$XDG_CONFIG_HOME/autostart/<identifier>.desktop`（デフォルトは `~/.config/autostart/`）に書き込みます。

```ini
Type=Application
Hidden=false
X-GNOME-Autostart-enabled=true
Exec=<executable> <arguments>
```

`Exec` フィールドは、[freedesktop.org Desktop Entry 仕様](https://specifications.freedesktop.org/desktop-entry-spec/)に従ってエスケープされます。予約文字（`"`、`` ` ``、`$`、`\\`）はバックスラッシュでエスケープされ、値に空白が含まれる場合は二重引用符で囲まれます。

`Status().Strategy` は `AutostartStrategyXDGAutostart` を返します。

[iOS / Android / サーバー]
サポートされていません。すべてのメソッドが `ErrAutostartNotSupported` を返します。これを明確に検出するには `errors.Is(err, application.ErrAutostartNotSupported)` を使用します。

```go
if err := app.Autostart.Enable(); err != nil {
    if errors.Is(err, application.ErrAutostartNotSupported) {
        // hide the toggle in the UI
        return
    }
    // real failure — surface it
}
```

@end

## 識別子

`Options.Identifier` が空の場合、アプリケーション名からデフォルト値が導出されます。

| プラットフォーム | デフォルトの識別子 |
| --- | --- |
| macOS（バンドル形式） | アプリケーションのバンドル識別子（例：`com.example.MyApp`） |
| macOS（非バンドル形式） | `wails.autostart.<slug>`。`<slug>` は `application.Options.Name` から導出されます |
| Windows | `application.Options.Name` のスラッグ（小文字化し、`A-Za-z0-9._-` 以外を削除し、空白をハイフンに変換） |
| Linux | Windows と同じスラッグ |

識別子は `^[A-Za-z0-9._-]+$` に一致し、200 文字以下である必要があります。macOS では逆引き DNS 形式を推奨します（launchd の Label で慣例的に使用される表記と一致します）。

`AutostartOptions.Identifier` を上書きすると、Windows では同じ識別子がレジストリ値の名前として、Linux では `.desktop` のファイル名として再利用されます。そのため、単一の文字列でプラットフォームをまたいで登録を識別できます。

## 古い登録の検出

`Disable()` と `Status()` は、識別子を検索するのではなく、**登録された実行ファイルのパスを `os.Executable()`（シンボリックリンクをすべて解決したもの）と照合して**登録を特定します。これは次のことを意味します。

- **リリース間で識別子を変更しても安全です。** 実行ファイルのパスが同じである限り、古い登録も引き続き `Status()` で検出でき、`Disable()` でクリーンアップされます。
- **別のパスにあるアプリケーションの 2 つ目のコピーが、最初のコピーの登録を上書きすることはありません。** 各バイナリの場所は個別に追跡されます。
- **シンボリックリンクを使用したインストール（Homebrew、Scoop）は安定して扱われます。** 照合前に `os.Executable()` へ `filepath.EvalSymlinks` が適用されるため、Homebrew のアップグレードでリンク先が差し替えられても、エントリが取り残されることはありません。

これが対象としないケース: ユーザーがバイナリを無関係なパスへ移動または名前変更すると、古い登録は孤立します（存在しなくなったファイルを指したままになります）。安定したインストールパスから配布されるアプリでは、この点を心配する必要はありません。ポータブルな単一ファイルのバイナリとして配布されるアプリは、自身を移動する前に `Disable()` を呼び出すか、常に安定したシンボリックリンク経由で起動してください。*対象外*

## 例

```go
package main

import (
    "errors"

    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "My App",
    })

    // Restore the user's preference on startup
    if userPrefersAutostart() {
        if err := app.Autostart.Enable(); err != nil {
            if !errors.Is(err, application.ErrAutostartNotSupported) {
                app.Logger.Error("autostart", "error", err)
            }
        }
    }

    app.Run()
}
```

ステータス、有効化、無効化の各ボタンを備えた、完全に実行可能なサンプルは [`examples/autostart/`](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/autostart) にあります。
