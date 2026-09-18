---
title: "自動起動"
description: "macOS、Windows、Linuxでユーザーのログイン時にアプリケーションを起動するよう登録します"
slug: "features/autostart/basics"
sourcePath: "features/autostart/basics.md"
---

## 自動起動

`app.Autostart`は、ユーザーのログイン時にアプリケーションが自動的に起動するよう登録します。プラットフォームごとに適切なネイティブの仕組みを選択し、シンボリックリンクされたインストールパス（Homebrew、Scoop）を解決するため、バイナリをアップグレードしても登録が無効になりません。

登録が有効になるのは即時ではなく、<strong>次回のログイン時</strong>です。

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

`Enable`は繰り返し呼び出しても安全です。呼び出すたびに登録が上書きされるため、ユーザーの設定を永続化している場合は、起動のたびに呼び出しても問題ありません。

### `EnableWithOptions`

カスタムオプションを使用して登録します。

```go
func (m *AutostartManager) EnableWithOptions(opts AutostartOptions) error
```

**`AutostartOptions`：**

| フィールド | 型 | 説明 |
| --- | --- | --- |
| `Identifier` | `string` | 自動生成された登録IDを上書きします。後述の「識別子」を参照してください。 |
| `Arguments` | `[]string` | ログイン時の起動において、実行可能ファイルのパスの後ろに追加される引数（例：`--hidden`）。 |

### `Disable`

自動起動の登録を削除します。アプリケーションが登録されていなかった場合は`nil`を返します。無効化操作は冪等です。

```go
func (m *AutostartManager) Disable() error
```

### `IsEnabled`

登録が存在するかどうかを報告します。処理は高速ですが、登録済みのパスは検証しません。

```go
func (m *AutostartManager) IsEnabled() (bool, error)
```

### `Status`

登録状態の全情報を返します。

```go
func (m *AutostartManager) Status() (AutostartStatus, error)
```

**`AutostartStatus`：**

| フィールド | 型 | 説明 |
| --- | --- | --- |
| `Enabled` | `bool` | 登録が存在するかどうか。 |
| `Path` | `string` | 登録アーティファクトのディスク上の場所（plist のパス、`.desktop`のパス、またはレジストリのサブキー）。`Enabled`が false の場合は空です。 |
| `Strategy` | `AutostartStrategy` | アプリの登録に使用された仕組み（[プラットフォームごとの動作](#heading-2)を参照）。 |

## プラットフォームごとの動作

@tabs{sync-key="platform"}
[macOS]
アプリのパッケージ化方法に応じて、次の 2 つの仕組みを使用します。

- **macOS 13 以降、バンドルされた`.app`**：`SMAppService.mainAppService`。サンドボックス化されたアプリおよび Mac App Store ビルドで動作します。TCC のオートメーション許可プロンプトは表示されません（以前の AppleScript を使用する方式では表示されていました）。
- **13より前の macOS、またはバンドルされていないバイナリ**：`RunAtLoad=true`を指定した LaunchAgent の plist を`~/Library/LaunchAgents/<identifier>.plist`に書き込みます。

呼び出し元がどちらの処理経路を使用したか判別できるように、`Status()`は`AutostartStrategySMAppService`または`AutostartStrategyLaunchAgent`を返します。

アプリをバンドルなしからバンドルありへアップグレードした場合は、孤立した LaunchAgent が古いビルドを起動し続けないように、`Status()`で両方の経路を確認し、`Disable()`でクリーンアップします。

[Windows]
`HKCU\Software\Microsoft\Windows\CurrentVersion\Run`の下にレジストリ値を追加します。値の名前には自動起動の`Identifier`を設定し、データには引用符で囲んだ実行ファイルのパスと、指定されている場合は`Arguments`を設定します。

引数の引用符処理は`CommandLineToArgvW`の規則（引用符の前にあるバックスラッシュを二重化）に従うため、空白や引用符を含むパスも正しくラウンドトリップできます。

`Status().Strategy`は`AutostartStrategyRegistryRun`を返します。

[Linux]
次の内容を持つ XDG 自動起動エントリを`$XDG_CONFIG_HOME/autostart/<identifier>.desktop`（デフォルトは`~/.config/autostart/`）に書き込みます。

```ini
Type=Application
Hidden=false
X-GNOME-Autostart-enabled=true
Exec=<executable> <arguments>
```

`Exec`フィールドは、[freedesktop.org Desktop Entry 仕様](https://specifications.freedesktop.org/desktop-entry-spec/)に従ってエスケープします。予約文字（`"`、`` ` ``、`$`、`\\`）はバックスラッシュでエスケープし、値に空白文字が含まれる場合は二重引用符で囲みます。

`Status().Strategy`は`AutostartStrategyXDGAutostart`を返します。

[iOS / Android / サーバー]
サポートされていません。すべてのメソッドが`ErrAutostartNotSupported`を返します。これを明確に検出するには`errors.Is(err, application.ErrAutostartNotSupported)`を使用します。

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

`Options.Identifier`が空の場合は、アプリ名からデフォルト値を生成します。

| プラットフォーム | デフォルトの識別子 |
| --- | --- |
| macOS（バンドルあり） | アプリのバンドル識別子（例：`com.example.MyApp`） |
| macOS（バンドルなし） | `wails.autostart.<slug>`。ここで、`<slug>`は`application.Options.Name`から生成されます |
| Windows | `application.Options.Name`のスラッグ（小文字に変換し、`A-Za-z0-9._-`以外を削除し、空白をハイフンに置換） |
| Linux | Windows と同じスラッグ |

識別子は`^[A-Za-z0-9._-]+$`に一致し、200文字以下でなければなりません。macOS ではドメイン名の要素を逆順に並べた形式を推奨します（launchd の Label で慣例的に使用される表記と一致します）。

`AutostartOptions.Identifier`を上書きすると、同じ識別子が Windows ではレジストリ値名として、Linux では`.desktop`ファイル名として再利用されるため、単一の文字列でプラットフォームをまたいで登録を識別できます。

## 古い登録の検出

`Disable()`と`Status()`は、識別子を検索するのではなく、<strong>登録済みの実行可能ファイルのパスを`os.Executable()`（すべてのシンボリックリンクを解決したもの）</strong>と照合することで登録を特定します。これは次のことを意味します。

- **リリース間で識別子を変更しても安全です。** 実行可能ファイルのパスが同じである限り、古い登録は引き続き`Status()`で検出でき、`Disable()`で削除されます。
- **別のパスにあるアプリの2つ目のコピーが、最初のコピーの登録を上書きすることはありません。** 各バイナリの場所は個別に追跡されます。
- **シンボリックリンクを使用したインストール（Homebrew、Scoop）でも安定して動作します。** 照合前に`os.Executable()`へ`filepath.EvalSymlinks`が適用されるため、Homebrew のアップグレードでリンク先が切り替わっても、登録エントリが取り残されることはありません。

この仕組みで対応<em>できない</em>ケースがあります。ユーザーがバイナリを無関係なパスへ移動したり名前を変更したりすると、古い登録は孤立し、存在しなくなったファイルを指すことになります。安定したインストールパスから提供されるアプリでは、この点を気にする必要はありません。ポータブルな単一ファイルのバイナリとして提供されるアプリは、自身を移動する前に`Disable()`を呼び出すか、常に安定したシンボリックリンクを介して起動してください。

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

ステータス／有効化／無効化ボタンを備えた、完全に実行可能な例は[`examples/autostart/`](https://github.com/wailsapp/wails/tree/master/v3/examples/autostart)にあります。
