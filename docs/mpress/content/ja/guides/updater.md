---
title: "アップデーター"
description: "Wails v3向けのアプリ内自己更新機能 — 差し替え可能なプロバイダー、暗号学的検証、アトミックな置換、テーマ設定や置き換えが可能な既定のUIを提供します。"
slug: "guides/updater"
sourcePath: "guides/updater.md"
---

アップデーターを使うと、ダウンロード、検証、置換のパイプラインを独自に構築しなくても、アプリ内でソフトウェア更新を配信できます。`app.Updater`上で動作し、1つ以上の差し替え可能な`Provider`（GitHub Releases、keygen.sh、Sparkle AppCast、オープンなWails Update Manifestプロトコル、または独自実装）を受け入れ、設定済みの公開鍵でダウンロードを認証し、実行中のバイナリを安全に置き換え、すべての状態遷移を標準のWailsイベントバス経由で通知します。

![「更新の準備完了」状態の既定のアップデーターウィンドウ — 状態に応じたアイコン、バージョン表示（v1.0.0 → v2.0.1 · 8.8 MB）、GFMテーブルを含むMarkdown形式で表示されたリリースノート、単一の主要アクション。](/assets/updater/default-window-ready.png)

## クイックスタート

```go {title="main.go"}
package main

import (
    "context"
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/updater"
    "github.com/wailsapp/wails/v3/pkg/updater/providers/github"
)

func main() {
    app := application.New(application.Options{Name: "Demo"})

    gh, _ := github.New(github.Config{Repository: "myorg/myapp"})
    if err := app.Updater.Init(updater.Config{
        CurrentVersion: "1.0.0",
        Providers:      []updater.Provider{gh},
    }); err != nil {
        log.Fatal(err)
    }

    if err := app.Updater.CheckAndInstall(context.Background()); err != nil {
        log.Printf("update: %v", err)
    }

    _ = app.Run()
}
```

これにより、フレームワークの更新ウィンドウが開き、GitHubを確認してプラットフォーム用アセットをダウンロードし、検証してバイナリを置き換えた後、ユーザーが再起動するのを待ちます。

## ライフサイクル

`app.Updater`は、次の状態を持つステートマシンです（`updater.State`）。

| 状態 | 該当する状況 |
| --- | --- |
| `unconfigured` | `Init`が呼び出される前 |
| `idle` | `Init`の後、最初の確認を行う前 |
| `checking` | `Check`の処理中 |
| `up-to-date` | 最新のプロバイダー応答で、呼び出し元が最新であると報告されたとき |
| `available` | 新しいリリースが見つかったが、まだダウンロードを開始していないとき |
| `downloading` | プロバイダーからバイト列をストリーミングしているとき |
| `verifying` | ダウンロードが完了し、署名またはダイジェストを検証しているとき |
| `installing` | 検証済みのバイト列を展開し、名前を変更してステージングディレクトリに配置しているとき |
| `ready` | 更新のステージングが完了した状態。適用するには`Restart`を呼び出します |
| `error` | それまでのいずれかの手順が失敗したとき |

現在の状態は、いつでも`app.Updater.State()`で確認できます。また、状態が遷移するたびにWailsイベントが発行されます（[イベント](#heading-13)を参照）。

`Restart` は、実行中のアプリケーションに終了を要求する前に、ヘルパーが `application.New` に到達するまで待機します。起動タイムアウトの既定値は 30 秒です。アプリケーションが `application.New` より前に時間のかかる初期化を行う場合は、`Config.HelperReadyTimeout` に、例えば `time.Minute` のような長い時間を設定してください。ゼロは既定値を選択し、負の時間は拒否されます。起動がタイムアウトすると、`Restart` は `updater.ErrHelperNotReady` を返し、実行中のアプリケーションを終了させません。

既定のウィンドウには現在の状態が自動的に反映されます。たとえば、`Check`が更新なしを返した場合は次のように表示され、ユーザーは<strong>閉じる</strong>でウィンドウを閉じます。

![「最新」状態の既定のアップデーターウィンドウ — 緑色のチェックマーク、「最新の状態です」という見出し、単一の「閉じる」ボタン。](/assets/updater/default-window-up-to-date.png)

## プロバイダー

次のインターフェースを満たすものは、すべて`Provider`として使用できます。

```go
type Provider interface {
    Name() string
    Check(ctx context.Context, req CheckRequest) (*Release, error)
    Download(ctx context.Context, r *Release, dst io.Writer, onProgress func(written, total int64)) error
}
```

リポジトリには4つの実装が含まれています。

### GitHub Releases — `updater/providers/github`

```go {title="github provider"}
gh, err := github.New(github.Config{
    Repository:    "myorg/myapp",     // your owner/repo (required)
    Token:         "ghp_…",           // optional; raises rate limit + private repos
    Prerelease:    false,             // include pre-releases in latest lookup
    ChecksumAsset: "SHA256SUMS",      // optional sibling asset for digest verification
    BaseURL:       "",                // optional override (e.g. GitHub Enterprise)
    AssetMatcher:  nil,               // optional custom asset-picker; nil uses DefaultAssetMatcher
    HTTPClient:    nil,               // optional client override
})
```

既定のアセット照合では、ファイル名に含まれる`GOOS`と`GOARCH`の部分文字列を基に選択し、一般的な別名（`amd64` / `x86_64` / `x64`、`arm64` / `aarch64`、`386` / `i386` / `x86` / `ia32`）も認識します。独自の命名規則を使用する場合は、次のように設定します。

```go
gh, _ := github.New(github.Config{
    Repository: "myorg/myapp",
    AssetMatcher: func(req updater.CheckRequest, assets []github.ReleaseAsset) int {
        for i, a := range assets {
            if strings.Contains(a.Name, "my-naming-convention") &&
               strings.Contains(a.Name, req.Platform) {
                return i
            }
        }
        return -1 // no match
    },
})
```

`ChecksumAsset`は、内容が`<sha256>  <filename>`形式の行で構成された、同じリリース内の別アセットの名前です（`sha256sum`と`shasum -a 256`が生成する形式です）。プロバイダーは`Check`の実行中にこれを取得し、選択したアーティファクトに一致する行を見つけて`Release.Verification.Digest`を設定するため、フレームワークがダウンロードを検証できます。

### keygen.sh — `updater/providers/keygen`

```go {title="keygen provider"}
kg, err := keygen.New(keygen.Config{
    Account:    "your-account-slug", // required
    Product:    "product-uuid",      // optional but recommended when account has multiple products
    Package:    "",                  // optional further narrowing
    Channel:    "stable",            // "stable" / "rc" / "beta" / "alpha" / "dev"
    Filetype:   "",                  // optional artifact filetype filter ("dmg", "exe", …)
    Token:      "prod-…",            // product / environment / user / admin token; wins over LicenseKey
    LicenseKey: "",                  // license key auth (used only when Token is empty)
    BaseURL:    "",                  // optional API base override
    HTTPClient: nil,                 // optional client override; redirect-strip wrapper still applied
})
```

プロバイダーは、keygen.shのアーティファクトごとのSHA-512チェックサムとEd25519ph署名を、フレームワークの`Release.Verification`ブロックに自動的に対応付けます。追加の接続設定は不要です。

<strong>トークン形式：</strong>keygen.shのトークンには、ロールを示す接頭辞（`admi-` / `prod-` / `envi-` / `user-`）が付いています。ダッシュボードに表示される未加工のUUIDは、トークンの<em>識別子</em>であり、シークレット値ではありません。シークレットが表示されるのは、トークンの作成時だけです。詳細については、keygen.shの[認証ドキュメント](https://keygen.sh/docs/api/authentication/)を参照してください。

### Sparkle AppCast — `updater/providers/appcast`

```go {title="appcast provider"}
ac, err := appcast.New(appcast.Config{
    URL:        "https://your.app/appcast.xml", // required
    Channel:    "stable",                       // optional sparkle:channel filter
    HTTPClient: nil,                            // optional client override
})
```

既存のSparkle / WinSparkleインフラストラクチャに変更を加えず、そのまま導入できます。フィードから`sparkle:shortVersionString`、`<enclosure url type length sparkle:os sparkle:edSignature>`、`sparkle:channel`を読み取ります。

Sparkle 1のDSA署名（`sparkle:dsaSignature`）はサポートされていません。この署名方式を使用しているプロジェクトは、EdDSA（Sparkle 2）に移行することを推奨します。

### Wails Update Manifest — `updater/providers/endpoint`

```go {title="endpoint provider"}
ep, err := endpoint.New(endpoint.Config{
    URL:        "https://updates.example.com/check", // required; supports {{platform}} / {{arch}} / {{version}} / {{channel}} placeholders
    Channel:    "stable",                            // optional channel filter
    Headers:    nil,                                 // optional headers, e.g. {"Authorization": "License <key>"}
    HTTPClient: nil,                                 // optional client override; redirect-strip wrapper still applied
})
```

オープンな[Wails Update Manifestプロトコル](/reference/update-manifest/)を使用します。これは、最新リリースとプラットフォーム別のアーティファクトを、チェックサムおよび署名とともにインラインで記述する単一のJSONドキュメントです。同じドキュメントを静的ファイルホスト（S3、GitHub Pages、任意のCDN。すべてのプラットフォームを列挙したマニフェストをチャネルごとに1つ公開します）でも、動的更新サーバーでも使用できます（プロバイダーは確認のたびに`platform`、`arch`、`version`、`channel`を送信するため、サーバーはアーティファクトを1つだけ返したり、ライセンスに基づいてアクセスを制御したりできます）。

URLプレースホルダーを使用すると、静的なレイアウトを1行で設定できます。

```go
ep, _ := endpoint.New(endpoint.Config{
    URL: "https://cdn.example.com/updates/{{platform}}/{{arch}}/stable.json",
})
```

設定したヘッダーは、マニフェストへのすべてのリクエストで送信されます。アーティファクトのダウンロードでそれらを再利用するのは、マニフェストと同じホスト上で、`https`から`http`へのダウングレードがない場合に限られます。また、オリジンをまたぐリダイレクトやダウングレードを伴うリダイレクトでは、`Authorization`ヘッダーが削除されます。

公開処理はCLIが担います。`wails3 updater manifest`は、1つのコマンドでリリースファイルのダイジェストを計算し、署名して、その内容を記述します。`wails3 updater verify`は、アップロード前にその結果を再検証します。[wails3 CLIを使用した公開](/reference/update-manifest/#publishing-with-the-wails3-cli)を参照してください。

### フォールバックチェーン

`Config.Providers`には順序があります。アップデーターは先頭から順に処理し、最初にリリースを返したプロバイダーを採用します。最初に「最新」と報告したプロバイダーがあれば、その時点でチェーンの処理を打ち切ります（フォールバックは「プライマリに到達できない」場合のためのものであり、「プロバイダー間で結果が一致しない」場合のためのものではありません）。エラーが発生した場合は、次のプロバイダーに進みます。

```go
app.Updater.Init(updater.Config{
    CurrentVersion: "1.0.0",
    Providers: []updater.Provider{
        kg, // primary: licensed customers
        gh, // fallback: public mirror
    },
})
```

### 独自プロバイダーの作成

メソッドは3つだけで、通常の実装なら約150行です。検証、アトミックなステージング、入れ替え、ウィンドウは Updater が担い、プロバイダーのコードは次のリリースを解決して、そのバイト列をストリーミングします。

```go
type CustomProvider struct { /* config */ }

func (p *CustomProvider) Name() string { return "custom" }

func (p *CustomProvider) Check(ctx context.Context, req updater.CheckRequest) (*updater.Release, error) {
    // Hit your update endpoint, decide whether req.CurrentVersion is current,
    // and return either nil (no upgrade) or a *Release with Artifact + optional
    // Verification populated.
    // Errors here drop through to the next provider in Config.Providers.
}

func (p *CustomProvider) Download(ctx context.Context, r *updater.Release, dst io.Writer, onProgress func(int64, int64)) error {
    // Stream the artifact's bytes to dst. Call onProgress(written, total) as
    // bytes flow past; the Updater debounces emits to ~10 Hz on the event bus.
}
```

ツリー内のプロバイダーを参考にしてください。いずれも単一の Go ファイルです。

## 暗号学的検証

リリースは、`Config.PublicKey`を信頼の起点として、フレームワークの検証機能によって認証されます。

```go
//go:embed publickey.pem
var publicKey []byte

app.Updater.Init(updater.Config{
    CurrentVersion: "1.0.0",
    Providers:      []updater.Provider{...},
    PublicKey:      publicKey,
})
```

対応アルゴリズム（`Release.Verification.SignatureAlgo`）：

| アルゴリズム | 署名対象 | 備考 |
| --- | --- | --- |
| `ed25519` | アーティファクトの SHA-256 ダイジェスト | Sparkle EdDSA で使用 |
| `ed25519ph` | Ed25519ph の事前ハッシュ（内部では SHA-512）を介したアーティファクト全体 | keygen.sh で使用 |
| `ecdsa-p256` | アーティファクトの SHA-256 ダイジェスト | 生の`r∥s`署名と DER 署名の両方に対応 |

さらに、リリースにハッシュは含まれるものの署名がない場合は、ダイジェストのみ（`DigestAlgo`：`sha256` / `sha512`）にも対応します。

署名検証における唯一の信頼アンカーは`Config.PublicKey`です。リリース元が独自の鍵に差し替えることはできません。`Signature`を含むリリースは、`Config.PublicKey`が設定されていなければ安全側に倒して失敗します。検証機能はダウンロード中の1回のストリーミング処理でダイジェストを計算するため、数 GB 規模の更新でも、検証のためにディスクを追加で走査することはありません。

@note{type="caution" title="ダイジェストのみ ≠ 暗号学的検証"}
`Digest`だけを含むリリースは、自分で管理する暗号学的な信頼の起点ではなく、レジストリの TLS と、レジストリ自体が提供する整合性保証によって認証されます。ビット腐敗の検出にはダイジェストのみの方式を使用し、侵害されたリリースパイプラインによる改ざんへの耐性が必要な場合は署名を使用してください。

@end

### 署名鍵の生成

```bash
wails3 updater genkey
# updater.key       — keep secret, use to sign releases
# updater.key.pub   — bundle in your app via go:embed
```

秘密鍵は PKCS#8 PEM、公開鍵は PKIX PEM です。`Config.PublicKey`には`.pub`ファイルをそのまま指定できます（生の32バイト鍵またはその base64 表現にも対応しており、インライン化用の値は`genkey`で出力できます）。リリースには`wails3 updater manifest -key updater.key ...`または`wails3 updater sign`で署名してください。詳しくは「[wails3 CLI を使用した公開](/reference/update-manifest/#publishing-with-the-wails3-cli)」を参照してください。

または Go では、次のようにします。

```go
import "crypto/ed25519"
import "crypto/rand"

pub, priv, _ := ed25519.GenerateKey(rand.Reader)
// Persist `priv` securely (HSM, signing CI, etc.); embed `pub` in your binary.
```

## アーティファクト形式

プロバイダーは、公開したファイルのバイト列をストリーミングします。その後、フレームワークが入れ替え前にファイルを展開します。

- **単一のバイナリ**（例：`myapp-linux-amd64`）— そのまま使用されます。Linux で一般的です。
- **`.zip`** — その場で展開されます。アーカイブには、トップレベルのエントリが正確に1つだけ含まれている必要があります（通常は macOS の`.app`バンドルまたは単一のバイナリ）。macOS で推奨されるパッケージ形式です。
- **`.tar.gz`** / **`.tgz`** — トップレベルのエントリを1つだけとする同じ規則に従い、その場で展開されます。バイナリとともにランタイムツリーを配布する Linux ディストリビューションに適しています。

トップレベルのエントリが複数あるアーカイブは拒否されます。フレームワークがディスク上で入れ替える対象は1つだけなので、アーカイブに複数のものが含まれている場合、「このアーカイブを所定の場所に入れ替える」という指示が曖昧になるためです。`.dmg`と`.pkg`（macOS）、および`.msi`（Windows）は v1 ではサポートされません。代わりに、バンドルの`.zip`を配布してください。展開時には zip-slip 対策が適用され、アーカイブルートの外部を指すシンボリックリンクは拒否されます。また、展開後の合計サイズ（2 GiB）とエントリ数（50 000）に上限が設けられます。

## デフォルトウィンドウ

`app.Updater.CheckAndInstall(ctx)`は、フレームワークが所有する520×540のウィンドウを開きます。このウィンドウには次の要素があります。

- 状態に応じて変化するヒーローアイコン（利用可能／ダウンロード中は青い「↓」、準備完了／最新版は緑の「✓」、エラーは赤い「!」）
- バージョンを示すピル：`v1.0.0 → v2.0.1 · 8.8 MB`
- <strong>レンダリング済み Markdown</strong>を表示する、スクロール可能なリリースノートパネル（段落、太字／斜体、リスト、GFM テーブル、インラインコード、フェンス付きコードブロック、h1～h3、リンク）
- 状態ごとに1つの主要アクション（インストール／再起動して適用／再試行）
- ゴーストスタイルの副次アクション（このバージョンをスキップ／後で通知）
- `prefers-color-scheme`によるダーク／ライトモード
- 合計サイズが不明な場合の不確定型プログレスシマー

Wails イベントバス上の`updater:*`イベントをリッスンし、`updater:user:*`アクションを Go に送り返します。

### CSS 変数によるテーマ設定

```go
app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{
        CSS: `:root { --accent: #ff6f00; --bg: #1a1a1a; --fg: #fafafa; }`,
    },
})
```

デフォルトのスタイルシートでは、次の変数が公開されています。どの変数も上書きできます。

| 変数 | デフォルト（ライト） | デフォルト（ダーク） |
| --- | --- | --- |
| `--bg` | `#f8f8fa` | `#1a1a1c` |
| `--surface` | `#ffffff` | `#232326` |
| `--surface-2` | `#f0f0f3` | `#2c2c30` |
| `--fg` | `#1d1d1f` | `#f5f5f7` |
| `--fg-dim` | `#6b6b73` | `#b0b0b8` |
| `--fg-faint` | `#99999f` | `#7a7a82` |
| `--border` | `#d6d6dc` | `#3a3a3e` |
| `--accent` | `#0a84ff` | `#0a84ff` |
| `--accent-fg` | `#ffffff` | — |
| `--success` | `#34c759` | — |
| `--error` | `#ff3b30` | — |
| `--radius` | `10px` | — |
| `--font` | システムフォントスタック | — |

### テンプレートを置き換える

独自の HTML を用意します。必要なのは、`wails:updater:*` イベントをリッスンし、`wails:updater:user:*` アクションを発行することだけです。

```go
app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{
        HTML: myCustomTemplate,
    },
})
```

InitialHTML ウィンドウはアセットサーバーのオリジンなしで読み込まれるため、`/wails/runtime.js` を動的に取得できません。このウィンドウからホストと通信する方法は 2 つあります。

1. **HTML をそのまま記述します。** フレームワークは、`WebviewWindowOptions.AllowSimpleEventEmit = true` と `HTML` が設定されたウィンドウに最小限の `window.wails.Events` shim を自動挿入します。これはまさに、アップデーターの組み込みパスと BYO パスが行う設定です。ビルド手順は不要です。以下の例ではこの方法を使用します。
2. **任意のバンドラー**（Vite、esbuild、Rollup）で `@wailsio/runtime` をバンドルし、ビルド時にカスタム HTML へインポートします。`Events.On` は完全にクライアント側で動作するため、そのまま使用できます。一方、`Events.Emit` はランタイムの fetch トランスポートを経由しますが、null オリジンでは機能しません。そのため、ランタイムの [`setTransport`](https://wails.io/wails/runtime.js) フックを使用し、`window._wails.invoke("wails:event:emit:<name>")` 経由でルーティングする小さな postMessage トランスポートを組み込みます。`window.wails.Events` がすでにスコープ内にある場合、フレームワークによる挿入は何も行わないため、2 つの方法が競合することはありません。

どちらの方法でも、カスタム HTML に記述する JS は同じ `Events.On` / `Events.Emit` API を呼び出します。

```html
<script>
const { On, Emit } = window.wails.Events;

On("wails:updater:update-available", (e) => {
    const rel = e.data ?? e;
    document.getElementById("ver").textContent = rel.version;
});

document.getElementById("install").addEventListener("click",
    () => Emit("wails:updater:user:install"));

// Ask the host to replay the current state so we paint correctly on (re)open.
Emit("wails:updater:window:ready");
</script>
```

この shim は、ベア名イベントに必要な最新ランタイムのサブセットを公開します。`Events.On(name, cb)` は購読解除関数を返し、`Events.Emit(nameOrEventObject)` は制限付きの `wails:event:emit:` postMessage パスを介してホストへルーティングします。ページ読み込み時に、インラインスクリプトが実行される前に一度だけインストールされます。

shim を<em>上書きしたい</em>場合（または別の方法で完全なランタイムを読み込む場合）は、ページ内の最初の `<script>` タグが実行される前に `window.wails.Events` を設定してください。これにより、自動挿入はスキップされます。

### ウィンドウ装飾

HTML に手を加えずに、ウィンドウオプション（サイズ、フレームなし、常に最前面）を上書きします。

```go
Window: &updater.BuiltinWindow{
    Options: updater.WindowOptions{
        Title:         "My App Updater",
        Width:         640,
        Height:        480,
        Frameless:     true,
        AlwaysOnTop:   true,
        DisableResize: false,
    },
},
```

### 独自のウィンドウを使用する

自分で作成した `*application.WebviewWindow` を使用して更新フローを動かします。アップデーターはそのウィンドウで `Show()` / `Close()` / `EmitEvent()` を呼び出し、何を描画するかは HTML 側で決定します。

![ピンクからオレンジへのグラデーション背景、白い角丸カード 1 枚、カスタムタイポグラフィを備え、同じアップデーターイベントによって表示状態が変化する「独自実装」のアップデーターウィンドウ。デフォルト UI を完全に置き換えられることを示しています。](/assets/updater/byo-custom-window.png)

```go
myWin := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:                "My Updater",
    Width:                520, Height: 460,
    HTML:                 myCustomHTML,
    AllowSimpleEventEmit: true,           // see security note below
})
app.Updater.Init(updater.Config{
    // …
    Window: updater.BYOWindow(myWin.AsUpdaterWindow()),
})
```

HTML では、組み込みテンプレートと同様に `window.wails.Events.On` / `Events.Emit` を使用します。ウィンドウがフレームワーク所有か独自作成かにかかわらず、`AllowSimpleEventEmit: true` が設定されたすべてのウィンドウには、フレームワークによって shim が自動挿入されます。API については、[テンプレートを置き換える](#heading-9)を参照してください。

@note{type="caution" title="BYO アップデーターウィンドウには `AllowSimpleEventEmit` が必要です"}
セキュリティのため、フレームワークはこのフィールドによって `wails:event:emit:` postMessage ショートカットの使用を制限します。このフィールドが設定されていないウィンドウは、ホスト側のカスタムイベントを生成できません。アップデーターのカスタム HTML shim は、このショートカットを介して `updater:user:*` イベントを発行します。そのため、このフィールドを設定し忘れた BYO ウィンドウでは、すべてのボタンクリックが通知なく破棄されます。ユーザーが「インストール」をクリックしても何も起こりません。

完全には制御できない HTML（リモート URL やユーザー提供コンテンツ）を読み込むウィンドウでは、`AllowSimpleEventEmit` を必ず<strong>無効</strong>にしてください。有効にすると、XSS シンクを含め、ページ内のあらゆる JavaScript が任意の `app.Event.On(name, …)` ハンドラーを起動できます。このショートカットが伝達できるのはベア名だけであり（ペイロードは不可）、バインディングの Call パスには到達できません。しかし、Go コードのハンドラーがイベント名だけに基づいて処理を実行する場合は、権限を必要とするカスタムイベントハンドラーを起動できてしまいます。

フレームワークの<em>組み込み</em>アップデーターウィンドウでは、このフィールドが内部で設定されています。覚えておく必要があるのは BYO の呼び出し元だけです。

@end

### ヘッドレス

```go
app.Updater.Init(updater.Config{
    // …
    Window: updater.WindowNone,
})
```

ウィンドウは一切開きません。独自の UI（または既存のメインウィンドウ）から `updater:*` イベントを購読し、ボタンハンドラーから `app.Updater.CheckAndInstall(ctx)` を呼び出します。更新が見つかった場合にだけ表示する定期的なバックグラウンドチェックや、更新フローをカスタム設定パネルに統合するアプリに適しています。

## イベント

Go と JavaScript はどちらも、標準の Wails イベントバスを介して購読します。**通信文字列を手入力しないでください**。updater パッケージ（Go）または runtime パッケージ（JS）がエクスポートする定数を使用してください。両レイヤーは同じ名前のセットを共有し、回帰テストによって同期が維持されています。

### Go から

定数は `github.com/wailsapp/wails/v3/pkg/updater` にあります。`app.Event.On(name, fn)` を介して購読します。コールバックは `*application.CustomEvent` を受け取り、その `Data` フィールドには[イベントリファレンス](#heading-14)に記載された型付きペイロードが格納されています。JSON デコードではなく、型アサーションを使用してください。

```go {title="main.go"}
import (
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/updater"
)

// …

app.Event.On(updater.EventUpdateAvailable, func(e *application.CustomEvent) {
    rel, ok := e.Data.(*updater.Release)
    if !ok { return }
    log.Printf("update found: %s", rel.Version)
})

app.Event.On(updater.EventDownloadProgress, func(e *application.CustomEvent) {
    p, ok := e.Data.(updater.Progress)
    if !ok { return }
    log.Printf("%d / %d bytes (%.0f KB/s)", p.Written, p.Total, p.Rate/1024)
})

app.Event.On(updater.EventError, func(e *application.CustomEvent) {
    info, ok := e.Data.(updater.ErrorInfo)
    if !ok { return }
    log.Printf("update failed during %s: %s", info.Stage, info.Message)
})
```

使用可能なすべての Go 定数：

| 定数 | 通信文字列 |
| --- | --- |
| `updater.EventCheckStarted` | `wails:updater:check-started` |
| `updater.EventUpdateAvailable` | `wails:updater:update-available` |
| `updater.EventNoUpdate` | `wails:updater:no-update` |
| `updater.EventDownloadStarted` | `wails:updater:download-started` |
| `updater.EventDownloadProgress` | `wails:updater:download-progress` |
| `updater.EventDownloadComplete` | `wails:updater:download-complete` |
| `updater.EventVerifying` | `wails:updater:verifying` |
| `updater.EventInstalling` | `wails:updater:installing` |
| `updater.EventUpdateReady` | `wails:updater:update-ready` |
| `updater.EventError` | `wails:updater:error` |
| `updater.EventMeta` | `wails:updater:meta` |
| `updater.EventWindowReady` | `wails:updater:window:ready` |
| `updater.EventUserInstall` | `wails:updater:user:install` |
| `updater.EventUserSkip` | `wails:updater:user:skip` |
| `updater.EventUserRemind` | `wails:updater:user:remind` |
| `updater.EventUserCancel` | `wails:updater:user:cancel` |
| `updater.EventUserRestart` | `wails:updater:user:restart` |

### JavaScript から

定数は `@wailsio/runtime` の `Updater.Events` 配下にあります。名前は Go と同じで、オートコンプリートで見つけやすいようサブ名前空間（`User.*`、`Window.*`）別に整理されています：

```js
import { Events, Updater } from "@wailsio/runtime";

Events.On(Updater.Events.UpdateAvailable, (e) => {
    console.log("update found:", e.data.version);
});

Events.On(Updater.Events.DownloadProgress, (e) => {
    const p = e.data;
    console.log(`${p.written} / ${p.total} bytes (${(p.rate/1024).toFixed(1)} KB/s)`);
});

Events.On(Updater.Events.Error, (e) => {
    const info = e.data;
    console.error(`update failed during ${info.stage}: ${info.message}`);
});
```

カスタム HTML がホストへ<em>返す</em>ユーザー操作イベントは、`Updater.Events.User` 配下にあります：

```js
import { Updater } from "@wailsio/runtime";

document.getElementById("install-btn").addEventListener("click", () => {
    // The framework window does this internally via the postMessage shim;
    // shown here for BYO templates that need to drive the flow themselves.
    window._wails.invoke("wails:event:emit:" + Updater.Events.User.Install);
});
```

### イベントリファレンス

購読側（ホスト → ページ）：

| 定数（Go） | 定数（JS） | ペイロード | 発生タイミング |
| --- | --- | --- | --- |
| `updater.EventCheckStarted` | `Updater.Events.CheckStarted` | なし | 各 `Check` ラウンドトリップの前 |
| `updater.EventUpdateAvailable` | `Updater.Events.UpdateAvailable` | `*Release` | `Check` が新しいリリースを検出したとき |
| `updater.EventNoUpdate` | `Updater.Events.NoUpdate` | なし | `Check` が最新であることを確認したとき |
| `updater.EventDownloadStarted` | `Updater.Events.DownloadStarted` | `*Release` | バイトのストリーミング開始時 |
| `updater.EventDownloadProgress` | `Updater.Events.DownloadProgress` | `Progress` | ダウンロード中は約 10 Hz |
| `updater.EventDownloadComplete` | `Updater.Events.DownloadComplete` | `*Release` | すべてのバイトの書き込み後、検証前 |
| `updater.EventVerifying` | `Updater.Events.Verifying` | `*Release` | 署名／ダイジェストのチェック開始時 |
| `updater.EventInstalling` | `Updater.Events.Installing` | `*Release` | 展開とステージングの開始時 |
| `updater.EventUpdateReady` | `Updater.Events.UpdateReady` | `*Release` | 再起動待機中 |
| `updater.EventError` | `Updater.Events.Error` | `ErrorInfo` | いずれかの段階で失敗したとき |
| `updater.EventMeta` | `Updater.Events.Meta` | `Meta` | スナップショットの再生前に、セッションごとに1回 |

ページ側（ページ → ホスト）— カスタムテンプレートを作成する場合は、コードで購読します：

| 定数（Go） | 定数（JS） | タイミング |
| --- | --- | --- |
| `updater.EventWindowReady` | `Updater.Events.Window.Ready` | ウィンドウの読み込み完了時。ホストが現在の状態を再送します |
| `updater.EventUserInstall` | `Updater.Events.User.Install` | `available`状態でのプライマリアクション |
| `updater.EventUserRestart` | `Updater.Events.User.Restart` | `ready`状態でのプライマリアクション |
| `updater.EventUserSkip` | `Updater.Events.User.Skip` | 「このバージョンをスキップ」 |
| `updater.EventUserRemind` | `Updater.Events.User.Remind` | 「後で通知」 |
| `updater.EventUserCancel` | `Updater.Events.User.Cancel` | 閉じるボタン |

## APIリファレンス

### `updater.Config`

| フィールド | 型 | 備考 |
| --- | --- | --- |
| `CurrentVersion` | `string` | **必須。** リリースのタグに使用するものと同じ文字列（`v`プレフィックスなし） |
| `Providers` | `[]updater.Provider` | **必須。** 順序付きフォールバックチェーン |
| `PublicKey` | `[]byte` | PEMまたは生のバイト列。省略可能ですが、署名付きリリースではこれがないと安全側に倒して失敗します |
| `CheckInterval` | `time.Duration` | ゼロ以外の場合、`CheckAndInstall`を呼び出すバックグラウンドのポーリングループを開始します |
| `Platform` | `string` | アセットの選択に使用する`runtime.GOOS`をオーバーライドします |
| `Arch` | `string` | アセットの選択に使用する`runtime.GOARCH`をオーバーライドします |
| `Channel` | `string` | 現在は情報提供のみ。プロバイダー固有のチャンネルフィルタリングに使用します |
| `Window` | `updater.WindowOption` | `nil`（組み込みのデフォルト）、`&BuiltinWindow{…}`、`BYOWindow(handle)`、または`WindowNone` |

### `*updater.Updater`のメソッド

| シグネチャ | 用途 |
| --- | --- |
| `Init(cfg Config) error` | 構成します。2回目の呼び出しでは`ErrAlreadyConfigured`を返します |
| `State() State` | 現在のライフサイクルフェーズ |
| `CurrentVersion() string` | `Init`に渡されたバージョン |
| `Check(ctx) (*Release, error)` | プロバイダーチェーンを順に調べます。`(rel, nil)` = 検出、`(nil, nil)` = 最新、`(nil, err)` = すべて失敗 |
| `DownloadAndInstall(ctx) error` | ストリーミング、検証、展開（アーカイブの場合）、ステージングを行います。事前に`Check`を実行しておく必要があります |
| `CheckAndInstall(ctx) error` | 簡便メソッド：ウィンドウを開き、`Check`を実行し、見つかった場合は`DownloadAndInstall`を実行します |
| `Restart(ctx) error` | ヘルパーを起動し、`Host.Quit`を呼び出して終了します。ヘルパーが入れ替えと再起動を行います |
| `DownloadedPath() string` | ステージング済みの更新がディスク上に保存されている場所。存在しない場合は`""` |
| `SkipVersion(v string)` | `v`をスキップ対象として記録します。以降の`Check`では、そのバージョンを最新として扱います |
| `SkippedVersion() string` | 現在スキップ対象になっているバージョンを読み取ります |
| `StopPeriodicCheck()` | `Config.CheckInterval`によって開始されたタイマーをキャンセルし、ループが終了するまで待機します |

### エラー

| センチネル | 返される条件 |
| --- | --- |
| `ErrAlreadyConfigured` | 最初に成功した後の`Init` |
| `ErrNotConfigured` | `Init`より前のすべての操作 |
| `ErrNoPendingRelease` | 事前の`Check`なしでの`DownloadAndInstall` |
| `ErrDownloadInProgress` | `DownloadAndInstall` が、同メソッドの別の呼び出しの実行中に呼び出された |
| `ErrNotReady` | ステージ済みの更新がない状態での`Restart` |

## 入れ替えの仕組み

`Restart`は、センチネル環境変数を設定して現在のバイナリを再実行します。`application.New`は起動時にそれらを検出し、ヘルパーモードへ処理を切り替えます。

1. ヘルパーは、親PIDが終了するまで最大30秒待機します（`platformIsAlive`は、Windowsでは`syscall.OpenProcess` + `GetExitCodeProcess`、Unixでは`os.FindProcess` + `proc.Signal(syscall.Signal(0))`を使用してポーリングします）。
2. ヘルパーはターゲットをバックアップします（ファイルの場合はコピー、macOSの`.app`バンドルディレクトリの場合は再帰コピー）。
3. ヘルパーは、試行間に500ミリ秒のバックオフを挟みながら最大20回再試行し、ターゲットをステージ済みアーティファクトに置き換えます。
  - **Unix** — `os.RemoveAll(target)` + `os.Rename(newPath, target)`。古いinodeを参照する開いたファイルディスクリプターは引き続き有効です。
  - **Windows** — `os.Rename(target, target.old.<nanos>)` + `os.Rename(newPath, target)`。Windowsでは、イメージがまだマッピングされているファイルの名前は変更できますが、削除はできません。次回の更新時に、ヘルパーは、所有元のカーネルマッピングが解放済みの残存する`.old.*`兄弟ファイルを削除します。

4. ヘルパーは、新しいバイナリに元の実行ファイルモードを復元します（ダウンロードしたファイルはデフォルトのumaskで作成され、Unixでは`+x`が失われます。Windowsでは何も行いません）。
5. ヘルパーは、ヘルパーモード用の環境変数を消去し、置き換え後のバイナリを再起動します。
6. ヘルパーが終了します。

起動に失敗した場合、ヘルパーはバックアップを復元します。親プロセスが30秒以内に終了しなければ、ヘルパーはターゲットに手を加える前に中止します（そのため、シャットダウンダイアログが`Quit`をブロックしても、ユーザーは動作するアプリを維持できます）。

推奨されるパッケージ形式である`.zip`として配布されたmacOSの`.app`バンドルは、検証から準備完了までの間にアーカイブが展開されるため、ヘルパーは実体のあるディレクトリを所定の場所へ入れ替えられます。

## 定期的な確認

```go
app.Updater.Init(updater.Config{
    // …
    CheckInterval: 6 * time.Hour,
})
```

`CheckInterval > 0`の場合、バックグラウンドのgoroutineが設定された間隔で`CheckAndInstall`を呼び出します。別のフロー（確認／ダウンロード／検証／インストール）がすでに進行中の場合、その間に到来したティックは破棄されます。同時実行するステートマシンはサポートされていません。

何かが見つかった場合にのみ表示するサイレントなバックグラウンドポーリングを行うには、`Window: updater.WindowNone`を設定し、独自のUIから`EventUpdateAvailable`に応答します。

## スキップと後で通知

デフォルトウィンドウの「このバージョンをスキップ」ボタンは、`SkipVersion(rel.Version)`を使用して利用可能なバージョンを記録します。それ以降の`Check`は同じバージョンを検出すると最新版として扱います（ユーザーが`CurrentVersion`を更新するまで。これは`Restart`が成功すると自動的に行われます）。「後で通知」は何も記録せず、ウィンドウを閉じるだけです。

```go
// Reading what the user skipped (e.g. to surface in app settings)
if v := app.Updater.SkippedVersion(); v != "" {
    log.Printf("user skipped %s", v)
}

// Programmatically clearing the skip:
app.Updater.SkipVersion("")
```

## 配布チェックリスト

アップデーターがインストールするリリースを公開する前に、次の項目を確認してください。

1. **適切なアーカイブ形式を選択します。** macOS：`.app`バンドルの`.zip`。Linux：単一バイナリまたは`.tar.gz`。Windows：単一の`.exe`または`.zip`。`.dmg`／`.msi`／`.pkg`はサポートされていません。
2. **アーティファクトに署名します**。`Config.PublicKey`に対応する秘密鍵を使用してください。プロバイダーが公開するフィード（keygen.sh、AppCast）では、各プロバイダーの署名手順に従います。`ChecksumAsset`を使用するGitHub Releasesでは、`sha256sum`／`shasum -a 256`を使用して`SHA256SUMS`ファイルを生成します。
3. **バージョン文字列を一致させます。** `Config.CurrentVersion`とリリースのバージョンタグは完全に一致する必要があります（例：`1.0.0` ↔ タグ`v1.0.0`。先頭の`v`はプロバイダー側で削除されます）。
4. **ターゲットプラットフォームで入れ替えをテストします**。リリース前に少なくとも一度実施してください。コード署名、公証、Gatekeeperへの対応はプラットフォーム固有であり、アップデーター自体の対象外です。

## トラブルシューティング

**「signature requires a public key but none configured」** — リリースに`Signature`フィールドがありますが、`Config.PublicKey`が空です。公開鍵を設定するか、署名を含めないようにリリースパイプラインを変更してください。

**「digest mismatch」** — ダウンロードしたバイト列がプロバイダーの提示した内容と一致しません。通常は、ネットワーク障害による不完全なダウンロードか、破損したアーティファクトが原因です。再実行すると解決することがよくあります。

**ウィンドウが開いてもすぐに消え、Markdownも進捗も表示されない** — カスタムHTMLが`wails:runtime:ready`を呼び出していません。[テンプレートを置き換える](#heading-9)のshimを参照してください。

**Windowsの更新が完了せず、ヘルパーログに「remove old (attempt N): Access is denied」と表示される** — この問題は、このPRの`de764fb`より前のバージョンでのみ発生します。現在の実装では、この問題が発生しないように対象を別名へ変更する方式を使用しています。アップグレードしてください。

**macOS Gatekeeperが入れ替え後のバイナリをブロックする** — コード署名をエンドツーエンドで維持する必要があります。元の`.app`に署名し、ビルドパイプラインが更新時にエンタイトルメントを変更する場合は、再起動されるバイナリにも<em>さらに</em>再署名してください。

## 関連項目

- 実行可能な例：[`v3/examples/updater`](https://github.com/wailsapp/wails/tree/master/v3/examples/updater)
- テスト用デモリポジトリ：[`wailsapp/updater-demo`](https://github.com/wailsapp/updater-demo)
- チュートリアル：[Wailsアプリに自己更新機能を追加する](/tutorials/04-self-update-a-wails-app/)
