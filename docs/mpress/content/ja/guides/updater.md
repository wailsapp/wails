---
title: "アップデーター"
description: "Wails v3 向けのアプリ内セルフアップデート — 差し替え可能なプロバイダー、暗号学的検証、アトミックな入れ替え、テーマ設定や置き換えが可能なデフォルト UI を提供します。"
slug: "guides/updater"
sourcePath: "guides/updater.md"
---

アップデーターを使うと、ダウンロード、検証、入れ替えのパイプラインを独自に構築せずに、アプリ内でソフトウェア更新を配信できます。`app.Updater` 上で動作し、差し替え可能な 1 つ以上の `Provider`（GitHub Releases、keygen.sh、Sparkle AppCast、オープンな Wails Update Manifest プロトコル、または独自実装）を受け入れます。設定済みの公開鍵でダウンロードを認証し、実行中のバイナリを安全に入れ替え、すべての遷移を標準の Wails イベントバス経由で通知します。

![「更新の準備完了」状態のデフォルトアップデーターウィンドウ — 状態に応じたアイコン、バージョンピル（v1.0.0 → v2.0.1 · 8.8 MB）、GFM テーブルを含む Markdown でレンダリングされたリリースノート、および 1 つの主要アクション。](/assets/updater/default-window-ready.png)

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

これにより、フレームワークの更新ウィンドウが開き、GitHub を確認してプラットフォーム用アセットをダウンロードし、検証してバイナリを入れ替え、ユーザーが再起動するまで待機します。

## ライフサイクル

`app.Updater` は、次の状態を持つステートマシンです（`updater.State`）。

| 状態 | 該当するタイミング |
| --- | --- |
| `unconfigured` | `Init` が呼び出される前 |
| `idle` | `Init` の後、最初の確認を行う前 |
| `checking` | `Check` の実行中 |
| `up-to-date` | 直近のプロバイダー応答で、呼び出し元が最新であると報告されたとき |
| `available` | 新しいリリースが見つかり、まだダウンロードを開始していないとき |
| `downloading` | プロバイダーからバイトデータをストリーミング中 |
| `verifying` | ダウンロードが完了し、署名またはダイジェストを確認中 |
| `installing` | 検証済みのバイトデータを展開し、ステージングディレクトリ内へ名前を変更して配置中 |
| `ready` | 更新のステージングが完了し、適用するには `Restart` を呼び出す |
| `error` | それ以前のいずれかの手順が失敗したとき |

現在の状態は、いつでも `app.Updater.State()` で確認できます。また、遷移のたびに Wails イベントが発行されます（[イベント](#heading-13)を参照）。

`Restart` は、実行中のアプリに終了を要求する前に、ヘルパーが `application.New` に到達するまで待機します。デフォルトの起動タイムアウトは 30 秒です。アプリが `application.New` より前に時間のかかる初期化を行う場合は、`Config.HelperReadyTimeout` に `time.Minute` などの長い期間を設定してください。ゼロを指定するとデフォルトが選択され、負の期間は拒否されます。起動がタイムアウトすると、`Restart` は `updater.ErrHelperNotReady` を返し、実行中のアプリを開いたままにします。

更新は、インストール済みバイナリの隣（macOS では `.app` バンドルの隣）にある `wails-update-*` ディレクトリへステージングされるため、そのディレクトリには書き込み権限が必要です。

AppImage 内では、アップデーターは読み取り専用マウントではなく、`$APPIMAGE` で指定された `.AppImage` ファイルを置き換えて再起動します。

デフォルトウィンドウには現在の状態が自動的に反映されます。たとえば、`Check` がアップグレードなしを返すと、ユーザーには次の画面が表示され、**閉じる**で閉じることができます。

![「最新」状態のデフォルトアップデーターウィンドウ — 緑色のチェックマーク、「最新の状態です」という見出し、1 つの「閉じる」ボタン。](/assets/updater/default-window-up-to-date.png)

## プロバイダー

`Provider` は、次のインターフェースを満たす任意の実装です。

```go
type Provider interface {
    Name() string
    Check(ctx context.Context, req CheckRequest) (*Release, error)
    Download(ctx context.Context, r *Release, dst io.Writer, onProgress func(written, total int64)) error
}
```

リポジトリ内には 4 つの実装が含まれています。

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

デフォルトのアセットマッチャーは、ファイル名に含まれる `GOOS` と `GOARCH` の部分文字列に基づいて選択し、一般的な別名（`amd64` / `x86_64` / `x64`、`arm64` / `aarch64`、`386` / `i386` / `x86` / `ia32`）も認識します。独自の命名規則を使用する場合は、次のようにします。

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

`ChecksumAsset` は、内容が `<sha256>  <filename>` 形式の行で構成される、同じリリース内の別アセットの名前です（`sha256sum` と `shasum -a 256` が生成する形式です）。プロバイダーは `Check` の実行中にそれを取得し、選択したアーティファクトに一致する行を見つけて `Release.Verification.Digest` を設定するため、フレームワークがダウンロードを検証できます。

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

プロバイダーは、keygen.sh のアーティファクトごとの SHA-512 チェックサムと Ed25519ph 署名を、フレームワークの `Release.Verification` ブロックへ自動的にマッピングします。追加の接続設定は不要です。

**トークン形式：** keygen.sh のトークンには、ロールを示す接頭辞（`admi-` / `prod-` / `envi-` / `user-`）が付いています。ダッシュボードに表示される生の UUID はトークンの*識別子*であり、シークレット値ではありません。シークレットはトークン作成時にしか表示されません。詳細については、keygen.sh の[認証ドキュメント](https://keygen.sh/docs/api/authentication/)を参照してください。

### Sparkle AppCast — `updater/providers/appcast`

```go {title="appcast provider"}
ac, err := appcast.New(appcast.Config{
    URL:        "https://your.app/appcast.xml", // required
    Channel:    "stable",                       // optional sparkle:channel filter
    HTTPClient: nil,                            // optional client override
})
```

既存の Sparkle / WinSparkle インフラストラクチャへ変更なしで導入できます。フィードから `sparkle:shortVersionString`、`<enclosure url type length sparkle:os sparkle:edSignature>`、`sparkle:channel` を読み取ります。

Sparkle 1 の DSA 署名（`sparkle:dsaSignature`）はサポートされていません。この署名方式を使用しているプロジェクトは、EdDSA（Sparkle 2）へ移行してください。

### Wails Update Manifest — `updater/providers/endpoint`

```go {title="endpoint provider"}
ep, err := endpoint.New(endpoint.Config{
    URL:        "https://updates.example.com/check", // required; supports {{platform}} / {{arch}} / {{version}} / {{channel}} placeholders
    Channel:    "stable",                            // optional channel filter
    Headers:    nil,                                 // optional headers, e.g. {"Authorization": "License <key>"}
    HTTPClient: nil,                                 // optional client override; redirect-strip wrapper still applied
})
```

オープンな[Wails Update Manifest プロトコル](/reference/update-manifest/)を使用します。これは、最新リリースとプラットフォームごとのアーティファクトを、チェックサムおよび署名とともにインラインで記述する単一の JSON ドキュメントです。同じドキュメントを、静的ファイルホスト（S3、GitHub Pages、任意の CDN。全プラットフォームを列挙したマニフェストをチャネルごとに 1 つ公開）でも、動的な更新サーバーでも使用できます。動的な更新サーバーの場合、プロバイダーは確認のたびに `platform`、`arch`、`version`、`channel` を送信するため、サーバーはアーティファクトを 1 つだけ返したり、ライセンスに基づいて提供可否を制御したりできます。

URL プレースホルダーを使うと、静的なレイアウトを 1 行で設定できます。

```go
ep, _ := endpoint.New(endpoint.Config{
    URL: "https://cdn.example.com/updates/{{platform}}/{{arch}}/stable.json",
})
```

設定したヘッダーは、すべてのマニフェストリクエストで送信されます。アーティファクトのダウンロードでそれらを再利用するのは、マニフェストと同じホストで、かつ `https` から `http` へのダウングレードがない場合だけです。また、クロスオリジンまたはダウングレードを伴うリダイレクトでは、`Authorization` ヘッダーが削除されます。

公開側は CLI が担います。`wails3 updater manifest` は 1 つのコマンドでリリースファイルのダイジェストを生成し、署名して、その内容を記述します。`wails3 updater verify` は、アップロード前にその結果を再確認します。[wails3 CLI を使用した公開](/reference/update-manifest/#publishing-with-the-wails3-cli)を参照してください。

### フォールバックチェーン

`Config.Providers` には順序があります。アップデーターはこれを順番に処理します。最初にリリースを返したプロバイダーが採用され、最初に「最新」と報告したプロバイダーでチェーンの処理を打ち切ります（フォールバックは「プライマリに到達できない」場合のためのものであり、「プロバイダー間で結果が一致しない」場合のためではありません）。エラーが発生すると、次のプロバイダーへ進みます。

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

必要なメソッドは 3 つで、一般的な実装は約 150 行です。検証、アトミックなステージング、入れ替え、ウィンドウは Updater が担い、プロバイダーコードは次のリリースを特定してバイトデータをストリーミングします。

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

リポジトリ内のプロバイダーを参考にしてください。いずれも 1 つの Go ファイルで実装されています。

## 暗号学的検証

リリースは、`Config.PublicKey` を信頼の起点として、フレームワークの検証機能によって認証されます。

```go
//go:embed publickey.pem
var publicKey []byte

app.Updater.Init(updater.Config{
    CurrentVersion: "1.0.0",
    Providers:      []updater.Provider{...},
    PublicKey:      publicKey,
})
```

サポートされているアルゴリズム（`Release.Verification.SignatureAlgo`）：

| アルゴリズム | 署名対象 | 備考 |
| --- | --- | --- |
| `ed25519` | 成果物の SHA-256 ダイジェスト | Sparkle EdDSA で使用 |
| `ed25519ph` | Ed25519ph の事前ハッシュ（内部では SHA-512）を使用した成果物全体 | keygen.sh で使用 |
| `ecdsa-p256` | 成果物の SHA-256 ダイジェスト | 生の `r∥s` 署名と DER 署名の両方に対応 |

さらに、リリースにハッシュは含まれるものの署名がない場合は、ダイジェストのみ（`DigestAlgo`: `sha256` / `sha512`）にも対応します。

署名検証における唯一の信頼の基点は `Config.PublicKey` です。リリース元が独自の鍵に差し替えることはできません。`Signature` を含むリリースは、`Config.PublicKey` が設定されていない場合、安全側に倒して失敗します。検証機構はダウンロード中にストリーミング処理でダイジェストを計算するため、数 GB 規模の更新であっても、検証のためにディスクを追加で走査することはありません。

@note{type="caution" title="ダイジェストのみ ≠ 暗号学的検証"}
`Digest` のみを含むリリースは、レジストリの TLS と、レジストリ自体が提供する整合性保証によって認証されます。管理下にある暗号学的な信頼の基点に対して認証されるわけではありません。ビットロットの検出にはダイジェストのみを使用し、侵害されたリリースパイプラインによる改ざんへの耐性が必要な場合は署名を使用してください。

@end

### 署名鍵の生成

```bash
wails3 updater genkey
# updater.key       — keep secret, use to sign releases
# updater.key.pub   — bundle in your app via go:embed
```

秘密鍵は PKCS#8 PEM、公開鍵は PKIX PEM です。`Config.PublicKey` は `.pub` ファイルをそのまま受け付けます（生の 32 バイト鍵またはその base64 も受け付けます。インライン化用の値は `genkey` で出力できます）。リリースには `wails3 updater manifest -key updater.key ...` または `wails3 updater sign` で署名してください。[wails3 CLI を使用した公開](/reference/update-manifest/#publishing-with-the-wails3-cli)を参照してください。

または Go では、次のようにします。

```go
import "crypto/ed25519"
import "crypto/rand"

pub, priv, _ := ed25519.GenerateKey(rand.Reader)
// Persist `priv` securely (HSM, signing CI, etc.); embed `pub` in your binary.
```

## 成果物の形式

プロバイダーは公開されたファイルのバイト列をストリーミングし、フレームワークは入れ替え前にそのファイルを展開します。

- **単一バイナリ**（例: `myapp-linux-amd64`）— そのまま使用されます。Linux で一般的です。
- **`.zip`** — その場で展開されます。アーカイブには、トップレベルのエントリを厳密に 1 つだけ含める必要があります（通常は macOS の `.app` バンドルまたは単一バイナリ）。macOS で推奨されるパッケージ形式です。
- **`.tar.gz`** / **`.tgz`** — トップレベルのエントリを 1 つだけ含めるという同じ規則に従い、その場で展開されます。バイナリとともにランタイムツリーを配布する Linux ディストリビューションに便利です。

トップレベルのエントリを複数含むアーカイブは拒否されます。フレームワークがディスク上で入れ替える対象は 1 つだけなので、アーカイブに複数のものが含まれている場合、「このアーカイブを所定の場所に入れ替える」という指示が曖昧になるためです。`.dmg` と `.pkg`（macOS）、および `.msi`（Windows）は v1 ではサポートされていません。代わりに、バンドルの `.zip` を配布してください。展開時には zip-slip 対策が適用され、アーカイブルート外を指すシンボリックリンクは拒否されます。また、展開後の合計サイズ（2 GiB）とエントリ数（50 000）には上限があります。

## デフォルトウィンドウ

`app.Updater.CheckAndInstall(ctx)` は、フレームワークが管理する 520×540 のウィンドウを開きます。このウィンドウには次の要素があります。

- 状態に応じて変化するヒーローアイコン（利用可能／ダウンロード中は青色の ↓、準備完了／最新版は緑色の ✓、エラーは赤色の !）
- バージョンピル: `v1.0.0 → v2.0.1 · 8.8 MB`
- **レンダリング済み Markdown**（段落、太字／斜体、リスト、GFM テーブル、インラインコード、フェンス付きコードブロック、h1～h3、リンク）を表示する、スクロール可能なリリースノートパネル
- 各状態につき 1 つの主要アクション（インストール／再起動して適用／再試行）
- ゴーストスタイルの副次アクション（このバージョンをスキップ／後で通知）
- `prefers-color-scheme` によるダーク／ライトモード
- 合計サイズが不明な場合の不確定プログレスシマー

Wails イベントバス上の `updater:*` イベントをリッスンし、`updater:user:*` アクションを Go に返します。

### CSS 変数によるテーマ設定

```go
app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{
        CSS: `:root { --accent: #ff6f00; --bg: #1a1a1a; --fg: #fafafa; }`,
    },
})
```

デフォルトのスタイルシートでは次の変数が公開されており、どれでも上書きできます。

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
| `--font` | システムスタック | — |

### テンプレートの置き換え

独自の HTML を指定できます。必要なのは、`wails:updater:*` イベントをリッスンし、`wails:updater:user:*` アクションを送出することだけです。

```go
app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{
        HTML: myCustomTemplate,
    },
})
```

InitialHTML ウィンドウはアセットサーバーのオリジンなしで読み込まれるため、`/wails/runtime.js` を動的にフェッチできません。このようなウィンドウからホストと通信する方法は 2 つあります。

1. **HTML をそのまま記述します。** フレームワークは、`WebviewWindowOptions.AllowSimpleEventEmit = true` と `HTML` を設定して開いたすべてのウィンドウに、最小限の `window.wails.Events` シムを自動的に挿入します。アップデーターの組み込みパスと BYO パスはいずれも、まさにこの方法を使用します。ビルド手順は不要です。以下の例でもこの方法を使用しています。
2. **お好みのバンドラー**（Vite、esbuild、Rollup）で `@wailsio/runtime` をバンドルし、ビルド時にカスタム HTML へインポートします。`Events.On` は完全にクライアント側で動作するため、そのまま使用できます。一方、`Events.Emit` はランタイムの fetch トランスポートを経由しますが、null オリジンでは機能しません。そのため、`window._wails.invoke("wails:event:emit:<name>")` を経由してルーティングする小さな postMessage トランスポートを、ランタイムの [`setTransport`](/guides/custom-transport/) フックを使って組み込んでください。`window.wails.Events` がすでにスコープ内にある場合、フレームワークによる注入は何も行わないため、2 つの方法が競合することはありません。

どちらの方法でも、カスタム HTML に記述する JavaScript は同じ `Events.On` / `Events.Emit` API を呼び出します。

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

この shim は、修飾なしの名前を持つイベントに必要な最新ランタイムのサブセットを公開します。`Events.On(name, cb)` は購読解除関数を返し、`Events.Emit(nameOrEventObject)` は制限付きの `wails:event:emit:` postMessage 経路を介してホストへルーティングします。ページ読み込み時に、インラインスクリプトが実行される前に一度だけ組み込まれます。

shim を上書きすることを*意図している*場合（または別の方法で完全なランタイムを読み込む場合）は、ページ内の最初の `<script>` タグが実行される前に `window.wails.Events` を設定してください。これにより、注入処理はスキップされます。

### ウィンドウの装飾

HTML を変更せずに、ウィンドウオプション（サイズ、フレームレス、常に手前に表示）を上書きできます。

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

自分で作成した `*application.WebviewWindow` を対象に更新フローを動作させます。アップデーターはそのウィンドウで `Show()` / `Close()` / `EmitEvent()` を呼び出し、何をレンダリングするかは HTML 側で決定します。

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

HTML では、組み込みテンプレートと同様に `window.wails.Events.On` / `Events.Emit` を使用します。ウィンドウがフレームワーク所有か独自実装かにかかわらず、`AllowSimpleEventEmit: true` が設定されたすべてのウィンドウには、フレームワークによって shim が自動注入されます。API については、[テンプレートを置き換える](#heading-9)を参照してください。

@note{type="caution" title="BYO アップデーターウィンドウには `AllowSimpleEventEmit` が必要です"}
セキュリティ上の理由から、フレームワークは `wails:event:emit:` postMessage ショートカットの使用可否をこのフィールドで制限します。このフィールドが設定されていないウィンドウでは、ホスト側のカスタムイベントを生成できません。アップデーターのカスタム HTML 用 shim は、このショートカットを介して `updater:user:*` イベントを送出します。そのため、このフィールドを設定し忘れた BYO ウィンドウでは、すべてのボタンクリックが何の通知もなく無視されます。ユーザーが「インストール」をクリックしても、何も起こりません。

完全には制御できない HTML（リモート URL やユーザー提供コンテンツ）を読み込むウィンドウでは、`AllowSimpleEventEmit` を必ず**オフ**にしてください。オンにすると、XSS の攻撃経路を含め、ページ内のあらゆる JavaScript が任意の `app.Event.On(name, …)` ハンドラーを発火できます。このショートカットで伝達できるのは修飾なしの名前だけで、ペイロードはなく、binding/Call 経路にも到達できません。しかし、Go コードのハンドラーがイベント名だけを基に動作する場合は、権限を伴うカスタムイベントハンドラーを発火させる可能性があります。

フレームワークの*組み込み*アップデーターウィンドウでは内部的に設定済みです。設定を忘れないようにする必要があるのは BYO の呼び出し側だけです。

@end

### ヘッドレス

```go
app.Updater.Init(updater.Config{
    // …
    Window: updater.WindowNone,
})
```

ウィンドウは一切開かれません。独自の UI（または既存のメインウィンドウ）から `updater:*` イベントを購読し、ボタンハンドラーから `app.Updater.CheckAndInstall(ctx)` を呼び出します。何かが見つかった場合にのみ画面へ表示する定期的なバックグラウンドチェックや、更新フローを独自の設定パネルへ統合するアプリに便利です。

## イベント

Go と JavaScript は、どちらも標準の Wails イベントバスを介して購読します。**通信上の文字列を手入力しないでください**。updater パッケージ（Go）または runtime パッケージ（JS）がエクスポートする定数を使用してください。両レイヤーは同じ名前のセットを共有しており、リグレッションテストによって同期が保たれています。

### Go から

定数は `github.com/wailsapp/wails/v3/pkg/updater` にあります。`app.Event.On(name, fn)` を介して購読します。コールバックは `*application.CustomEvent` を受け取り、その `Data` フィールドには[イベントリファレンス](#heading-14)に記載された型付きペイロードが格納されます。JSON デコードは行わず、型アサーションを使用してください。

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

| 定数 | 通信上の文字列 |
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

定数は `@wailsio/runtime` の `Updater.Events` 以下にあります。Go と同じ名前が使用され、オートコンプリートで見つけやすいようにサブ名前空間（`User.*`、`Window.*`）ごとに整理されています。

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

カスタム HTML がホストへ*返送*するユーザー操作イベントは、`Updater.Events.User` 以下にあります。

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
| `updater.EventCheckStarted` | `Updater.Events.CheckStarted` | なし | `Check` の各往復処理の前 |
| `updater.EventUpdateAvailable` | `Updater.Events.UpdateAvailable` | `*Release` | `Check` が新しいリリースを検出したとき |
| `updater.EventNoUpdate` | `Updater.Events.NoUpdate` | なし | `Check` が最新版であることを確認したとき |
| `updater.EventDownloadStarted` | `Updater.Events.DownloadStarted` | `*Release` | バイトのストリーミング開始時 |
| `updater.EventDownloadProgress` | `Updater.Events.DownloadProgress` | `Progress` | ダウンロード中は約10 Hz |
| `updater.EventDownloadComplete` | `Updater.Events.DownloadComplete` | `*Release` | 全バイトの書き込み完了後、検証前 |
| `updater.EventVerifying` | `Updater.Events.Verifying` | `*Release` | 署名／ダイジェストの検査開始時 |
| `updater.EventInstalling` | `Updater.Events.Installing` | `*Release` | 展開とステージングの開始時 |
| `updater.EventUpdateReady` | `Updater.Events.UpdateReady` | `*Release` | 再起動待ち |
| `updater.EventError` | `Updater.Events.Error` | `ErrorInfo` | いずれかの段階が失敗したとき |
| `updater.EventMeta` | `Updater.Events.Meta` | `Meta` | スナップショットの再生前に、セッションごとに1回 |

ページ側（ページ → ホスト）— カスタムテンプレートを作成する場合は、コードで次のイベントを購読します。

| 定数（Go） | 定数（JS） | タイミング |
| --- | --- | --- |
| `updater.EventWindowReady` | `Updater.Events.Window.Ready` | ウィンドウの読み込み完了時。ホストが現在の状態を再生 |
| `updater.EventUserInstall` | `Updater.Events.User.Install` | `available`状態でのプライマリアクション |
| `updater.EventUserRestart` | `Updater.Events.User.Restart` | `ready`状態でのプライマリアクション |
| `updater.EventUserSkip` | `Updater.Events.User.Skip` | 「このバージョンをスキップ」 |
| `updater.EventUserRemind` | `Updater.Events.User.Remind` | 「後で通知」 |
| `updater.EventUserCancel` | `Updater.Events.User.Cancel` | 閉じるボタン |

## APIリファレンス

### `updater.Config`

| フィールド | 型 | 備考 |
| --- | --- | --- |
| `CurrentVersion` | `string` | <strong>必須。</strong>リリースのタグ付けに使用するものと同じ文字列（`v`プレフィックスなし） |
| `Providers` | `[]updater.Provider` | <strong>必須。</strong>順序付きフォールバックチェーン |
| `PublicKey` | `[]byte` | PEMまたは生のバイト列。省略可能ですが、これがない場合、署名付きリリースは安全側に倒して失敗します |
| `CheckInterval` | `time.Duration` | ゼロ以外の場合、`CheckAndInstall`を呼び出すバックグラウンドのポーリングループを開始 |
| `Platform` | `string` | アセット選択用の`runtime.GOOS`を上書き |
| `Arch` | `string` | アセット選択用の`runtime.GOARCH`を上書き |
| `Channel` | `string` | 現在は情報提供のみ。プロバイダー固有のチャンネルフィルタリング |
| `Window` | `updater.WindowOption` | `nil`（組み込みのデフォルト）、`&BuiltinWindow{…}`、`BYOWindow(handle)`、または`WindowNone` |

### `*updater.Updater`のメソッド

| シグネチャ | 目的 |
| --- | --- |
| `Init(cfg Config) error` | 構成します。2回目の呼び出しでは`ErrAlreadyConfigured`を返します |
| `State() State` | 現在のライフサイクルフェーズ |
| `CurrentVersion() string` | `Init`に渡されたバージョン |
| `Check(ctx) (*Release, error)` | プロバイダーチェーンを順に処理します。`(rel, nil)` = 検出、`(nil, nil)` = 最新、`(nil, err)` = すべて失敗 |
| `DownloadAndInstall(ctx) error` | ストリーミング、検証、（アーカイブの場合は）展開、ステージングを行います。事前に`Check`を呼び出す必要があります |
| `CheckAndInstall(ctx) error` | 便利なメソッド：ウィンドウを開き、`Check`を実行し、見つかった場合は`DownloadAndInstall`を実行します |
| `Restart(ctx) error` | ヘルパーを起動し、`Host.Quit`を呼び出して終了します。ヘルパーが入れ替えて再起動します |
| `DownloadedPath() string` | ステージング済みアップデートがディスク上に存在する場所。存在しない場合は`""` |
| `SkipVersion(v string)` | `v`をスキップ対象として記録します。以降の`Check`では最新として扱われます |
| `SkippedVersion() string` | 現在スキップされているバージョンを読み取ります |
| `StopPeriodicCheck()` | `Config.CheckInterval`によって開始されたタイマーをキャンセルし、ループが終了するまで待機します |

### エラー

| センチネル | 返される箇所 |
| --- | --- |
| `ErrAlreadyConfigured` | 最初に成功した後の`Init` |
| `ErrNotConfigured` | `Init`より前のすべての操作 |
| `ErrNoPendingRelease` | 事前に`Check`を実行せずに`DownloadAndInstall`を実行した場合 |
| `ErrDownloadInProgress` | 別の処理の実行中に`DownloadAndInstall`を呼び出した場合 |
| `ErrNotReady` | ステージ済みの更新がない状態で`Restart`を実行した場合 |

## 入れ替えの仕組み

`Restart`は、センチネル環境変数を設定して現在のバイナリを再実行します。`application.New`は起動時にそれらを検出し、ヘルパーモードに処理を切り替えます。

1. ヘルパーは、親PIDが終了するまで最大30秒待機します（`platformIsAlive`は、Windowsでは`syscall.OpenProcess` + `GetExitCodeProcess`、Unixでは`os.FindProcess` + `proc.Signal(syscall.Signal(0))`を使用してポーリングします）。
2. ヘルパーは対象をバックアップします（ファイルの場合はコピー、macOSの`.app`バンドルディレクトリの場合は再帰コピー）。
3. ヘルパーは、ステージ済みのアーティファクトで対象を置き換えます。試行の間に500ミリ秒のバックオフを挟み、最大20回再試行します。
  - **Unix** — `os.RemoveAll(target)` + `os.Rename(newPath, target)`。古いinodeを参照する開いたファイルディスクリプターは引き続き有効です。
  - **Windows** — `os.Rename(target, target.old.<nanos>)` + `os.Rename(newPath, target)`。Windowsでは、イメージがまだマップされているファイルの名前は変更できますが、削除はできません。次回の更新時に、ヘルパーは、所有元のカーネルマッピングが解放された残存`.old.*`兄弟ファイルをすべて削除します。

4. ヘルパーは、新しいバイナリに元の実行ファイルモードを復元します（ダウンロードしたファイルはデフォルトのumaskで作成されるため、Unixでは`+x`が失われます。Windowsでは何も行いません）。
5. ヘルパーはヘルパーモード用の環境変数を消去し、置き換え済みのバイナリを再起動します。
6. ヘルパーは終了します。

起動に失敗した場合、ヘルパーはバックアップを復元します。親プロセスが30秒以内に終了しない場合、ヘルパーは対象に手を加える前に中止します（そのため、シャットダウンダイアログが`Quit`をブロックしても、ユーザーは動作するアプリを維持できます）。

推奨パッケージ形式である`.zip`として配布されるmacOSの`.app`バンドルでは、検証から準備完了までの間にアーカイブが展開されるため、ヘルパーは実ディレクトリを所定の場所に入れ替えられます。

## 定期チェック

```go
app.Updater.Init(updater.Config{
    // …
    CheckInterval: 6 * time.Hour,
})
```

`CheckInterval > 0`の場合、バックグラウンドのgoroutineが設定された間隔で`CheckAndInstall`を呼び出します。別のフロー（チェック／ダウンロード／検証／インストール）がすでに進行中のときに到着したtickは破棄されます。状態機械の並行実行はサポートされていません。

何かが見つかった場合にのみ通知するサイレントなバックグラウンドポーリングを行うには、`Window: updater.WindowNone`を設定し、独自のUIで`EventUpdateAvailable`に反応するようにします。

## スキップと後で通知

デフォルトウィンドウの「このバージョンをスキップ」ボタンは、`SkipVersion(rel.Version)`を介して利用可能なバージョンを記録します。それ以降の`Check`は同じバージョンを検出すると最新版として扱います（ユーザーが`CurrentVersion`を更新するまで。これは`Restart`が成功すると自動的に行われます）。「後で通知」は何も記録せず、ウィンドウを閉じるだけです。

```go
// Reading what the user skipped (e.g. to surface in app settings)
if v := app.Updater.SkippedVersion(); v != "" {
    log.Printf("user skipped %s", v)
}

// Programmatically clearing the skip:
app.Updater.SkipVersion("")
```

## 配布チェックリスト

アップデーターがインストールするリリースを公開する前に、次を確認してください。

1. **適切なアーカイブ形式を選択します。** macOS：`.app`バンドルの`.zip`。Linux：単一バイナリまたは`.tar.gz`。Windows：単一の`.exe`または`.zip`。`.dmg`／`.msi`／`.pkg`はサポートされていません。
2. **`Config.PublicKey`に対応する秘密鍵でアーティファクトに署名します**。プロバイダーが公開するフィード（keygen.sh、AppCast）では、各プロバイダーの署名手順に従ってください。`ChecksumAsset`を使用するGitHub Releasesでは、`sha256sum`／`shasum -a 256`を使って`SHA256SUMS`ファイルを生成します。
3. **バージョン文字列を一致させます。** `Config.CurrentVersion`とリリースのバージョンタグは完全に一致する必要があります（例：`1.0.0` ↔ タグ`v1.0.0`。先頭の`v`はプロバイダー側で除去されます）。
4. **出荷前に、対象プラットフォームで入れ替えを少なくとも一度テストします**。コード署名、公証、Gatekeeperへの対応はプラットフォーム固有であり、アップデーター自体の対象外です。

## トラブルシューティング

**「signature requires a public key but none configured」** — リリースに`Signature`フィールドがありますが、`Config.PublicKey`が空です。公開鍵を設定するか、署名を含めないようにリリースパイプラインを変更してください。

**「digest mismatch」** — ダウンロードしたバイト列がプロバイダーの提示内容と一致しません。通常は、ネットワークの一時的な不具合による不完全なダウンロードか、破損したアーティファクトが原因です。再実行すると解決することがよくあります。

**ウィンドウが開いてもすぐに消え、Markdownも進行状況も表示されない** — カスタムHTMLが`wails:runtime:ready`を呼び出していません。[テンプレートを置き換える](#heading-9)のshimを参照してください。

**Windowsの更新が完了せず、ヘルパーログに「remove old (attempt N): Access is denied」と表示される** — このPRの`de764fb`より前のバージョンでのみ発生します。現在の実装ではこの問題が発生しないrename-aside方式を使用しています。アップグレードしてください。

**macOSのGatekeeperが入れ替え後のバイナリをブロックする** — コード署名をエンドツーエンドで維持する必要があります。元の`.app`に署名し、更新時にビルドパイプラインがentitlementを変更する場合は、再起動するバイナリにも*署名し直して*ください。

## 関連項目

- 実行可能なサンプル：[`v3/examples/updater`](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/updater)
- テスト用デモリポジトリ：[`wailsapp/updater-demo`](https://github.com/wailsapp/updater-demo)
- チュートリアル：[Wailsアプリに自己更新を追加する](/tutorials/04-self-update-a-wails-app/)
