---
title: "自己更新する Wails アプリ"
description: "GitHub Releases から自己更新する Wails v3 アプリケーションを構築します。`wails3 init` から、署名付きリリースの検証とヘルパーモードによる入れ替えまでを扱います。"
slug: "tutorials/04-self-update-a-wails-app"
sourcePath: "tutorials/04-self-update-a-wails-app.md"
---

このチュートリアルでは、新しい Wails v3 アプリケーションにアプリ内アップデーターを追加します。完成したアプリでは、次のことが可能になります。

- 必要に応じて（オプションでタイマーを使用して）GitHub Releases を確認する。
- 実行中の OS とアーキテクチャに適したアセットをダウンロードする。
- ダウンロードしたバイト列に対して SHA-256 ダイジェスト（オプションで Ed25519 署名も）を検証する。
- フレームワークのデフォルトの更新ウィンドウにリリースノートを表示する。
- 別個のヘルパー実行ファイルを配布せずに、実行中のバイナリを入れ替えて再起動する。

更新元には、無料でインフラも不要な **GitHub Releases** を使用します。同じパターンは [keygen.sh](/guides/updater/#keygensh--updaterproviderskeygen) と [Sparkle AppCast](/guides/updater/#sparkle-appcast--updaterprovidersappcast) にも適用できます。完了後は [Updater ガイド](/guides/updater/)を参照してください。

@note{type="tip" title="前提条件"}
- Go 1.26 以降
- `wails3` CLI がインストール済みであること（`go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0`）
- リリースをプッシュできる GitHub リポジトリ
- [QR Code Service チュートリアル](/tutorials/01-creating-a-service/)の知識があると役立ちますが、必須ではありません

@end

<br/>

@steps
### 新しい Wails アプリから始める
vanilla テンプレートを使用して新しいプロジェクトのひな形を生成します。

```bash
wails3 init -n updater-tutorial -t vanilla
cd updater-tutorial
```

これで、`main.go`、`frontend/`、および `Taskfile.yml` を含むディレクトリが作成されます。ビルドして起動できることを確認します。

```bash
wails3 task dev
```

空の Wails ウィンドウが開きます。終了して、次に進みます。

### updater の import を追加する
`main.go` を開き、import に 2 つの updater パッケージを追加します。

```go {title="main.go" ins="6-7"}
package main

import (
    _ "embed"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/updater"
    "github.com/wailsapp/wails/v3/pkg/updater/providers/github"
)
```

これにより、Updater 本体と GitHub Releases プロバイダーが取り込まれます。

### Updater を設定する
`app.Updater` は、すでにすべての `*application.App` に組み込まれています。必要なのは `Init` の呼び出しだけです。

```go {title="main.go"}
const currentVersion = "1.0.0"

gh, err := github.New(github.Config{
    Repository:    "yourorg/your-repo",   // ← change this
    ChecksumAsset: "SHA256SUMS",          // sibling file with sha256 digests
})
if err != nil {
    log.Fatalf("github.New: %v", err)
}

if err := app.Updater.Init(updater.Config{
    CurrentVersion: currentVersion,
    Providers:      []updater.Provider{gh},
}); err != nil {
    log.Fatalf("Updater.Init: %v", err)
}
```

これを `application.New` の後、`app.Run()` の前に配置します。

@note{type="note" title="バージョン文字列の形式"}
リリースのタグ付けに使用するものと同じバージョンを、先頭の `v` は**付けずに**渡します。プロバイダー側では、タグ名から `v` が取り除かれます。ここでの `1.0.0` ↔ GitHub 上の `v1.0.0` です。

@end

### 更新を実行するメニュー項目を追加する
同じ `main.go` に「更新を確認…」メニュー項目を追加します。

```go {title="main.go"}
menu := app.Menu.New()
app.Menu.SetApplicationMenu(menu)
appMenu := menu.AddSubmenu("App")
appMenu.Add("Check for Updates…").OnClick(func(*application.Context) {
    go func() {
        if err := app.Updater.CheckAndInstall(context.Background()); err != nil {
            app.Logger.Error("update", "error", err)
        }
    }()
})
```

`CheckAndInstall` はフレームワークの更新ウィンドウを開いて `Check` を実行し、リリースが見つかった場合は `DownloadAndInstall` を自動的に実行します。新しいリリースがない場合、ウィンドウは「最新です」状態のまま開き続けます。ユーザーは **閉じる** ボタンでウィンドウを閉じます。

@note{type="caution" title="goroutine で実行する"}
`CheckAndInstall` は検証とインストールが完了するまでブロックします。メニューのクリック処理から直接呼び出すと UI スレッドがブロックされるため、`go func()` で囲みます。

@end

### リリースがない状態で一度実行する
```bash
wails3 task dev
```

<strong>アプリ → 更新を確認…</strong>をクリックします。更新ウィンドウが短時間開き、GitHub API にアクセスして、`1.0.0` より新しいリリースがないことを確認した後、緑色の ✓ が付いた <strong>最新です</strong> 状態になります。

ここでエラーが発生した場合、通常は次のいずれかが原因です。

| 症状 | 対処方法 |
| --- | --- |
| `404 Not Found` | `Repository` フィールドが誤っています。`owner/repo` でなければなりません |
| `403 rate-limited` | github.Config に `Token: "ghp_…"` を追加します（`public_repo` スコープを持つ PAT を使用してください） |
| ネットワークエラー | 実行中のアプリから `api.github.com` にアクセスできることを確認します |

### テストリリースを公開する
`main.go` 内の `currentVersion` を `1.0.0` に上げます（そのままでも構いません）。リリースに添付するバイナリを生成するため、1 つのプラットフォーム向けにビルドします。

@tabs
[macOS]
```bash
wails3 task build:darwin
# produces bin/updater-tutorial.app
# zip it for the release asset:
cd bin && zip -r updater-tutorial-darwin-arm64.zip updater-tutorial.app && cd ..
```

[Linux]
```bash
wails3 task build:linux
# produces bin/updater-tutorial
mv bin/updater-tutorial bin/updater-tutorial-linux-amd64
```

[Windows]
```bash
wails3 task build:windows
# produces bin/updater-tutorial.exe
mv bin/updater-tutorial.exe bin/updater-tutorial-windows-amd64.exe
```

@end

バイナリと同じ場所に `SHA256SUMS` ファイルを生成します。

```bash
cd bin
shasum -a 256 updater-tutorial-* > SHA256SUMS
cat SHA256SUMS
```

次のような行が 1 行以上表示されます。

```
abc123…  updater-tutorial-darwin-arm64.zip
```

これを GitHub リポジトリで **v2.0.0** として公開します。

```bash
gh release create v2.0.0 \
    --title "v2.0.0" \
    --notes "First update for the self-update tutorial.

- **Bold** Markdown renders in the update window
- \`Code spans\` too
- Lists work
- GFM tables work" \
    bin/SHA256SUMS bin/updater-tutorial-*
```

@note{type="note" title="アセットの命名"}
デフォルトのアセットマッチャーは、ファイル名に含まれる `GOOS` と `GOARCH` の部分文字列を基に選択します。アセット名に `darwin`（または `linux` / `windows`）と `arm64`（または `amd64` / `386`）が含まれていれば、マッチャーで検出できます。カスタムマッチャーについては、[Updater ガイド](/guides/updater/#github-releases--updaterprovidersgithub)を参照してください。

@end

### アプリを実行して更新を確認する
`currentVersion` がまだ `1.0.0` の状態で、アプリをもう一度実行します。

```bash
wails3 task dev
```

<strong>アプリ → 更新を確認…</strong>をクリックします。今回は、次のような表示になります。

![「更新の準備完了」状態のデフォルトの updater ウィンドウ。バージョンのピル、Markdown でレンダリングされたリリースノート、および主要ボタンの「再起動して適用」が表示されています。](/assets/updater/default-window-ready.png)

- ヒーローアイコンが青色の ↓（「更新があります」）から緑色の ✓（「更新の準備完了」）に切り替わります。
- サブタイトルに `v1.0.0 → v2.0.0 · <size>` が表示されます。
- リリースノートパネルでは、太字、コードスパン、表を含む Markdown がレンダリングされます。
- ダウンロード中にプログレスバーが進みます（バイナリが小さいため、すぐに完了します）。

Updater は新しいバイナリを一時ディレクトリにステージングします。更新を完了するには、次の操作を行います。

- **再起動して適用**をクリックします。
- アプリが終了し、ヘルパーがバイナリを入れ替え、新しいバイナリが再起動します。
- 再起動したアプリは `currentVersion = "1.0.0"` と報告します（ハードコードしたため）が、ディスク上のバイト列は v2.0.0 ビルドと一致します。

実際のアプリでは、`currentVersion` をビルド時に `-ldflags` で設定します。これにより、新しいバイナリは自身が v2.0.0 になったことを認識し、次回の確認では更新が見つかりません。

### `currentVersion` をビルドに連携させる
定数をビルド時変数に置き換えます。

```go {title="main.go" ins="2,4"}
var (
    currentVersion = "dev" // overridden by -ldflags at release time
)
```

次に、ビルドコマンドで次のように指定します。

```bash
wails3 task build:darwin -- -ldflags "-X main.currentVersion=2.0.0"
```

または、`git describe --tags` から値を取得するように、`Taskfile.yml` に `-ldflags` を追加します。

### 暗号署名を追加する（本番環境では推奨）
SHA256SUMS を使用する方法では*完全性*（バイト列が GitHub に保存されたものと一致すること）は検証できますが、*真正性*（そのバイト列が、侵害されたメンテナーアカウントではなく、リリースパイプラインによって生成されたこと）は検証できません。改ざん耐性を確保するには、各リリースに Ed25519 鍵で署名します。

```bash
# One-time: generate the keypair
ssh-keygen -t ed25519 -f updater-key -N "" -C "wails-updater"
#   updater-key      — keep secret (build server, HSM, password manager)
#   updater-key.pub  — bundle in your app
```

各リリースで、各アセットの SHA-256 ダイジェストを秘密鍵で署名します。次のような小さな Go ヘルパーを使用できます。

```go {title="cmd/sign-release/main.go"}
package main

import (
    "crypto/ed25519"
    "crypto/sha256"
    "encoding/base64"
    "fmt"
    "io"
    "os"
)

func main() {
    priv, _ := os.ReadFile("updater-key")
    key := ed25519.PrivateKey(priv) // raw 64-byte private key

    f, _ := os.Open(os.Args[1])
    defer f.Close()
    h := sha256.New()
    _, _ = io.Copy(h, f)
    sig := ed25519.Sign(key, h.Sum(nil))
    fmt.Println(base64.StdEncoding.EncodeToString(sig))
}
```

デフォルトの GitHub プロバイダーは、現時点では個別の署名ファイルを取得しません。署名ファイルを取得する[カスタムプロバイダーを作成](/guides/updater/#writing-your-own-provider)するか、すべてのアーティファクトをサーバー側で署名し、そのダイジェストと署名の両方を API 経由で公開する**keygen.sh**に切り替えることができます。

公開鍵をアプリに埋め込みます。

```go {title="main.go" ins="1,6"}
//go:embed updater-key.pub
var updaterPublicKey []byte

app.Updater.Init(updater.Config{
    CurrentVersion: currentVersion,
    PublicKey:      updaterPublicKey,
    Providers:      []updater.Provider{gh},
})
```

`PublicKey`を設定すると、`Signature`を同梱するすべてのリリースは、この鍵による検証に成功する必要があります。リリース元が独自の鍵に差し替えることはできません。ビルド時に帯域外で鍵を固定する目的は、まさにそこにあります。

### ウィンドウをカスタマイズする
デフォルトのウィンドウで一般的なユースケースに対応できます。さらに細かく制御する必要がある場合は、カスタマイズの度合いに応じて、次の 3 つの方法から 1 つを選択してください。

@tabs
[CSS のみ]
```go
app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{
        CSS: `:root { --accent: #ff6f00; --radius: 16px; }`,
    },
})
```

変数の完全な一覧については、[CSS 変数によるテーマ設定](/guides/updater/#theme-via-css-variables)のセクションを参照してください。

[カスタム HTML]
```go
//go:embed updater-window.html
var updaterHTML string

app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{HTML: updaterHTML},
})
```

HTML では、Wails のイベントチャネルを介して`updater:*`イベントを購読し、`updater:user:*`アクションを送信する必要があります。JS シムについては、[テンプレートを置き換える](/guides/updater/#replace-the-template)を参照してください。

[独自のウィンドウを使用]
```go
myWin := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:                "My App Updater",
    Width:                520, Height: 460,
    HTML:                 updaterHTML,
    AllowSimpleEventEmit: true,  // required — see security note
})
app.Updater.Init(updater.Config{
    // …
    Window: updater.BYOWindow(myWin.AsUpdaterWindow()),
})
```

すでに独自のウィンドウ基盤があり、別のウィンドウを開かずにアップデーターからそのウィンドウを制御したい場合に便利です。デフォルトと同じアップデーターイベントで駆動する、完全にカスタムの HTML テンプレートは次のようになります。

![ピンクからオレンジへのグラデーション背景と独自の角丸カードレイアウトを採用した、独自実装のアップデーターウィンドウ。デフォルト UI を完全に置き換えられることを示しています。](/assets/updater/byo-custom-window.png)

@note{type="caution" title="`AllowSimpleEventEmit` は必須です"}
アップデーターのカスタム HTML シムは、`wails:event:emit:`の postMessage ショートカットを介して Install / Skip / Remind / Restart を実行します。このショートカットはセキュリティ上、このフィールドによって制限されています。設定し忘れると、ボタンを押しても何も起こらず、エラーも表示されません。完全には管理できない HTML を読み込むウィンドウでは有効にしないでください。脅威モデルについては、ガイドの[独自のウィンドウを使用](/guides/updater/#bring-your-own-window)セクションを参照してください。

@end

@end

### バックグラウンドで自動チェックを実行する
メニューのクリックの代わりに、またはクリックに加えて、タイマーでチェックするには、次のようにします。

```go {ins="5"}
app.Updater.Init(updater.Config{
    CurrentVersion: currentVersion,
    Providers:      []updater.Provider{gh},
    PublicKey:      updaterPublicKey,
    CheckInterval:  6 * time.Hour,
})
```

タイマーが作動するたびに、手動クリックと同じ`CheckAndInstall`フローが実行されます。実際に更新が見つかるまで定期チェックで何も表示したくない場合は、`Window: updater.WindowNone`を設定します。そのうえで、表示する UX を決めるために`EventUpdateAvailable`を自分で購読してください。

@end

## これで完了です

これで、Wails アプリは次の機能を備えました。

- オンデマンドおよびタイマーで GitHub Releases の更新をチェックします。
- リリースノートを Markdown として、洗練されたデフォルトウィンドウに表示します。
- 公開した SHA-256 ダイジェストを使用してダウンロードを検証します。
- 必要に応じて、ビルド時に埋め込んだ公開鍵を使用して Ed25519 署名を検証します。
- 実行中のバイナリをその場で置き換え、自動的に再起動します。

## 次のステップ

- [アップデーターガイド](/guides/updater/)には、完全な API リファレンス、すべてのイベントと設定オプション、およびヘルパーモードでの置換処理の仕組みが記載されています。
- クローンできる完全な動作例については、[`v3/examples/updater`](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples/updater)を参照してください。
- テスト対象リポジトリ[`wailsapp/updater-demo`](https://github.com/wailsapp/updater-demo)では、推奨されるリリースアセットの構成を確認できます。

## 本番環境での注意事項

- **macOS でのコード署名** — Gatekeeper では、置き換えられるバイナリが署名および公証されている必要があります。リリース用に zip 圧縮する*前に*、`.app`バンドルへ署名してください。アップデーターはバイト列をそのまま保持し、再署名は行いません。
- **Windows のウイルス対策ソフト** — インターネットからダウンロードした未署名の`.exe`ファイルは、SmartScreen の警告を引き起こす可能性があります。バイナリを Authenticode 証明書で署名してください。署名しない場合は、制限の厳しいマシンを使用するユーザーがアプリを許可リストに追加しなければならない可能性があります。
- **アトミックなリリース** — `SHA256SUMS`とバイナリは、別々のコミットではなく一緒に公開してください。アップデーターはサイドカーファイルとバイナリを別々に取得するため、両者に不整合があるとダイジェストチェックが安全側に失敗し、更新は拒否されます。
- **スキップしたバージョン** — デフォルトウィンドウの「このバージョンをスキップ」ボタンを押すと、スキップしたことがローカルに記録されます。重大なセキュリティ更新を公開する場合は、以前のリリースを閉じたユーザーに自動的にスキップされないよう、新しいバージョン番号を付けてください。
