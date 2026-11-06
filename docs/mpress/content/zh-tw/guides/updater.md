---
title: "更新程式"
description: "為 Wails v3 提供應用程式內自我更新功能——支援可插拔的提供者、密碼學驗證、原子交換，以及可自訂佈景主題或替換的預設 UI。"
slug: "guides/updater"
sourcePath: "guides/updater.md"
---

更新程式可在應用程式內發布軟體更新，無須強迫你自行建置下載／驗證／交換管線。它建構於`app.Updater`之上，可接受一或多個可插拔的`Provider`（GitHub Releases、keygen.sh、Sparkle AppCast、開放的 Wails Update Manifest 通訊協定，或你自行實作的提供者），使用已設定的公開金鑰驗證下載內容、安全地交換執行中的二進位檔，並透過標準 Wails 事件匯流排呈現每一次狀態轉換。

![處於「更新已就緒」狀態的預設更新程式視窗——包含依狀態變化的圖示、版本膠囊標籤（v1.0.0 → v2.0.1 · 8.8 MB）、以 Markdown 轉譯且包含 GFM 表格的版本資訊，以及單一主要操作。](/assets/updater/default-window-ready.png)

## 快速開始

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

這會開啟框架的更新視窗、檢查 GitHub、下載目前平台的成品、加以驗證、交換二進位檔，然後等待使用者重新啟動。

## 生命週期

`app.Updater`是一個包含下列狀態的狀態機（`updater.State`）：

| 狀態 | 時機 |
| --- | --- |
| `unconfigured` | 呼叫`Init`之前 |
| `idle` | 呼叫`Init`之後、進行任何檢查之前 |
| `checking` | 正在執行`Check` |
| `up-to-date` | 最近一次提供者回應指出呼叫端已是最新版本 |
| `available` | 已找到新版本，但尚未開始下載 |
| `downloading` | 正在從提供者串流位元組 |
| `verifying` | 下載完成，正在檢查簽章／摘要 |
| `installing` | 正在解封裝已驗證的位元組，並重新命名至暫存目錄 |
| `ready` | 更新已暫存；呼叫`Restart`以套用 |
| `error` | 先前任何步驟失敗 |

你可以隨時使用`app.Updater.State()`讀取目前狀態。每次狀態轉換也都會發出 Wails 事件（請參閱[事件](#heading-13)）。

`Restart`會等待輔助程式到達`application.New`，之後才要求執行中的應用程式結束。預設的啟動逾時時間為30秒。如果你的應用程式在`application.New`之前需要進行耗時的初始化，請將`Config.HelperReadyTimeout`設為較長的持續時間，例如`time.Minute`。零值會選用預設值；負持續時間會遭到拒絕。如果啟動逾時，`Restart`會傳回`updater.ErrHelperNotReady`，並讓執行中的應用程式保持開啟。

更新會暫存在已安裝二進位檔旁的`wails-update-*`目錄中（在 macOS 上則位於`.app`套件旁），因此該目錄必須可寫入。

在 AppImage 內，更新程式會替換並重新啟動由`$APPIMAGE`指定的`.AppImage`檔案，而不是唯讀掛載點。

預設視窗會自動反映目前狀態——例如，當`Check`傳回沒有可用升級時，使用者會看到此畫面，並以**關閉**將其關閉：

![處於「已是最新版本」狀態的預設更新程式視窗——綠色勾號、「你已是最新版本」標題，以及單一「關閉」按鈕。](/assets/updater/default-window-up-to-date.png)

## 提供者

任何符合此介面的項目都是`Provider`：

```go
type Provider interface {
    Name() string
    Check(ctx context.Context, req CheckRequest) (*Release, error)
    Download(ctx context.Context, r *Release, dst io.Writer, onProgress func(written, total int64)) error
}
```

原始碼樹內隨附四種實作。

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

預設的成品比對器會依檔名中的`GOOS` + `GOARCH`子字串選取，並辨識常見別名（`amd64` / `x86_64` / `x64`、`arm64` / `aarch64`、`386` / `i386` / `x86` / `ia32`）。若使用自訂命名方案：

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

`ChecksumAsset`是同一版本中的另一項成品名稱，其內容為`<sha256>  <filename>`格式的行（即`sha256sum`與`shasum -a 256`所產生的格式）。提供者會在`Check`期間擷取該成品，找出與所選成品相符的行，並填入`Release.Verification.Digest`，讓框架驗證下載內容。

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

提供者會自動將 keygen.sh 每項成品的 SHA-512總和檢查碼與 Ed25519ph 簽章對應至框架的`Release.Verification`區塊，無須額外接線。

**權杖格式**：keygen.sh 權杖帶有角色前綴（`admi-`／`prod-`／`envi-`／`user-`）。您在儀表板中看到的原始 UUID 是權杖的*識別碼*，而不是其祕密值；祕密值僅在建立權杖時顯示。如需詳細資訊，請參閱 keygen.sh 的[驗證文件](https://keygen.sh/docs/api/authentication/)。

### Sparkle AppCast — `updater/providers/appcast`

```go {title="appcast provider"}
ac, err := appcast.New(appcast.Config{
    URL:        "https://your.app/appcast.xml", // required
    Channel:    "stable",                       // optional sparkle:channel filter
    HTTPClient: nil,                            // optional client override
})
```

可直接接入現有的 Sparkle / WinSparkle 基礎架構，無須變更。它會從資訊來源讀取`sparkle:shortVersionString`、`<enclosure url type length sparkle:os sparkle:edSignature>`及`sparkle:channel`。

不支援 Sparkle 1的 DSA 簽章（`sparkle:dsaSignature`）；採用該簽署方案的專案應輪替至 EdDSA（Sparkle 2）。

### Wails Update Manifest — `updater/providers/endpoint`

```go {title="endpoint provider"}
ep, err := endpoint.New(endpoint.Config{
    URL:        "https://updates.example.com/check", // required; supports {{platform}} / {{arch}} / {{version}} / {{channel}} placeholders
    Channel:    "stable",                            // optional channel filter
    Headers:    nil,                                 // optional headers, e.g. {"Authorization": "License <key>"}
    HTTPClient: nil,                                 // optional client override; redirect-strip wrapper still applied
})
```

使用開放的[Wails Update Manifest 通訊協定](/reference/update-manifest/)：以單一 JSON 文件描述最新版本及其各平台成品，並內嵌總和檢查碼與簽章。同一份文件既可用於靜態檔案主機（S3、GitHub Pages、任何 CDN——每個通道發布一份列出所有平台的資訊清單），也可用於動態更新伺服器（提供者每次檢查都會傳送`platform`、`arch`、`version`及`channel`，因此伺服器可以只傳回一項成品，或依授權狀態決定是否提供更新）。

URL 預留位置讓靜態版面配置只需一行設定：

```go
ep, _ := endpoint.New(endpoint.Config{
    URL: "https://cdn.example.com/updates/{{platform}}/{{arch}}/stable.json",
})
```

每次資訊清單請求都會傳送已設定的標頭；只有在成品位於資訊清單本身的主機上，且未從`https`降級至`http`時，成品下載才會重複使用這些標頭。任何跨來源或降級重新導向都會移除`Authorization`標頭。

CLI 負責發布端：`wails3 updater manifest`可用單一命令計算摘要、簽署並描述你的版本檔案，而`wails3 updater verify`會在上傳前重新檢查結果。請參閱[使用 wails3 CLI 發布](/reference/update-manifest/#publishing-with-the-wails3-cli)。

### 後援鏈

`Config.Providers`具有順序。更新程式會依序走訪：第一個傳回版本的提供者勝出；第一個回報「已是最新版本」的提供者會使鏈短路（後援是為了處理「主要提供者無法連線」，而不是「提供者意見不一致」）。發生錯誤時則會前進至下一個提供者。

```go
app.Updater.Init(updater.Config{
    CurrentVersion: "1.0.0",
    Providers: []updater.Provider{
        kg, // primary: licensed customers
        gh, // fallback: public mirror
    },
})
```

### 自行撰寫提供者

僅需三個方法；一般實作約為150行。Updater 負責驗證、原子暫存、交換及視窗；提供者程式碼則解析下一個版本並串流位元組：

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

請以原始碼樹內的提供者作為參考——每個提供者都只有一個 Go 檔案。

## 密碼學驗證

框架的驗證器會使用`Config.PublicKey`作為信任根，以驗證版本的真實性：

```go
//go:embed publickey.pem
var publicKey []byte

app.Updater.Init(updater.Config{
    CurrentVersion: "1.0.0",
    Providers:      []updater.Provider{...},
    PublicKey:      publicKey,
})
```

支援的演算法（`Release.Verification.SignatureAlgo`）：

| 演算法 | 簽署內容 | 備註 |
| --- | --- | --- |
| `ed25519` | 成品的 SHA-256 摘要 | 供 Sparkle EdDSA 使用 |
| `ed25519ph` | 透過 Ed25519ph 的預雜湊簽署完整成品（內部使用 SHA-512） | 供 keygen.sh 使用 |
| `ecdsa-p256` | 成品的 SHA-256 摘要 | 接受原始`r∥s`與 DER 簽章 |

此外，當發行版本提供雜湊值但沒有簽章時，也支援僅摘要驗證（`DigestAlgo`：`sha256` / `sha512`）。

`Config.PublicKey`是簽章驗證的唯一信任錨點——發行來源無法替換成自己的金鑰。若發行版本含有`Signature`，但未設定`Config.PublicKey`，驗證將以安全方式失敗。驗證器會在下載時以串流方式計算摘要，因此即使是數 GB 的更新，驗證也不會增加額外的磁碟掃描作業。

@note{type="caution" title="僅摘要驗證 ≠ 密碼學驗證"}
僅含`Digest`的發行版本，是依據登錄服務的 TLS 以及該登錄服務本身提供的完整性保證進行驗證，而不是依據您控制的密碼學信任根。僅摘要驗證適合用來偵測位元腐化；若要抵抗遭入侵的發行管線所造成的竄改，請使用簽章。

@end

### 產生簽署金鑰

```bash
wails3 updater genkey
# updater.key       — keep secret, use to sign releases
# updater.key.pub   — bundle in your app via go:embed
```

私密金鑰採用 PKCS#8 PEM 格式，公開金鑰採用 PKIX PEM 格式；`Config.PublicKey`可直接接受`.pub`檔案（也接受原始32位元組金鑰或其 base64 編碼；`genkey`會輸出後者，以便內嵌）。請使用`wails3 updater manifest -key updater.key ...`或`wails3 updater sign`簽署發行版本；請參閱[使用 wails3 CLI 發佈](/reference/update-manifest/#publishing-with-the-wails3-cli)。

或者使用 Go：

```go
import "crypto/ed25519"
import "crypto/rand"

pub, priv, _ := ed25519.GenerateKey(rand.Reader)
// Persist `priv` securely (HSM, signing CI, etc.); embed `pub` in your binary.
```

## 成品格式

提供者會串流傳輸您發佈之檔案的位元組；框架接著會在置換前將其解壓縮：

- **單一二進位檔**（例如`myapp-linux-amd64`）——直接使用。常見於 Linux。
- **`.zip`**——就地解壓縮。封存檔必須恰好包含一個頂層項目（通常是 macOS `.app`套件或單一二進位檔）。這是 macOS 的建議封裝格式。
- **`.tar.gz`** / **`.tgz`**——依照相同的單一頂層項目規則就地解壓縮。適用於將執行階段目錄樹與二進位檔一併提供的 Linux 發行版。

含有多個頂層項目的封存檔會遭到拒絕：框架只會置換單一磁碟目標，因此當封存檔包含多個項目時，「將此封存檔置換到指定位置」的含義並不明確。`.dmg`與`.pkg`（macOS）以及`.msi`（Windows）不受v1支援——請改為散佈套件的`.zip`。解壓縮程序會防範 Zip Slip、拒絕指向封存檔根目錄以外位置的符號連結，並限制解壓縮後的總大小（2 GiB）及項目數量（50 000）。

## 預設視窗

`app.Updater.CheckAndInstall(ctx)`會開啟一個由框架擁有、大小為520×540的視窗，其中包含：

- 依狀態變化的主視覺圖示（有可用更新或正在下載時顯示藍色 ↓，準備就緒或已是最新版本時顯示綠色 ✓，發生錯誤時顯示紅色 !）
- 版本膠囊標籤：`v1.0.0 → v2.0.1 · 8.8 MB`
- 可捲動的發行說明面板，其中包含**轉譯後的 Markdown**（段落、粗體／斜體、清單、GFM 表格、行內程式碼、圍欄程式碼區塊、h1–h3、連結）
- 每種狀態僅提供一個主要動作（安裝／重新啟動並套用／再試一次）
- 幽靈樣式的次要動作（略過此版本／稍後提醒我）
- 透過`prefers-color-scheme`切換深色／淺色模式
- 總大小未知時顯示不確定進度的微光動畫

它會在 Wails 事件匯流排上監聽`updater:*`事件，並將`updater:user:*`動作發回 Go。

### 透過 CSS 變數設定佈景主題

```go
app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{
        CSS: `:root { --accent: #ff6f00; --bg: #1a1a1a; --fg: #fafafa; }`,
    },
})
```

預設樣式表公開下列變數——您可以覆寫其中任何一個：

| 變數 | 預設值（淺色） | 預設值（深色） |
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
| `--font` | 系統字型堆疊 | — |

### 取代範本

提供您自己的 HTML；它只需要監聽`wails:updater:*`事件並發出`wails:updater:user:*`動作：

```go
app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{
        HTML: myCustomTemplate,
    },
})
```

InitialHTML 視窗載入時沒有資產伺服器來源，因此無法動態擷取`/wails/runtime.js`。您有兩種方式可從此類視窗與主機通訊：

1. <strong>直接撰寫 HTML。</strong>框架會將最精簡的`window.wails.Events`墊片自動注入任何以`WebviewWindowOptions.AllowSimpleEventEmit = true`開啟且設定了`HTML`的視窗——更新程式的內建路徑和自備視窗路徑正是如此運作。不需要建置步驟。下方範例使用的就是此方式。
2. **使用您偏好的打包工具**（Vite、esbuild、Rollup）打包`@wailsio/runtime`，並在建置時將其匯入自訂 HTML。`Events.On`完全在用戶端運作，因此可直接使用；`Events.Emit`會經由執行階段的 fetch 傳輸，而 null origin 會使其失效——因此，請透過執行階段的[`setTransport`](/guides/custom-transport/)掛鉤安裝一個精簡的 postMessage 傳輸，並經由`window._wails.invoke("wails:event:emit:<name>")`進行路由。如果`window.wails.Events`已在作用域內，框架注入不會執行任何動作，因此兩種方法不會互相衝突。

無論採用哪種方式，您在自訂 HTML 中編寫的 JS 都會呼叫相同的`Events.On`／`Events.Emit` API：

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

此 shim 提供以裸名稱使用事件所需的現代執行階段功能子集：`Events.On(name, cb)`會傳回取消訂閱函式，而`Events.Emit(nameOrEventObject)`會透過受控的`wails:event:emit:` postMessage 路徑將事件路由至主機。它會在頁面載入時安裝一次，且早於您的任何行內指令碼執行。

如果您*想要*覆寫此 shim（或正以其他方式載入完整執行階段），請在頁面的第一個`<script>`標籤執行前設定`window.wails.Events`，注入程序便會自行略過。

### 視窗外觀

無須修改 HTML，即可覆寫視窗選項（大小、無框、永遠置頂）：

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

### 使用您自己的視窗

使用您自行建立的`*application.WebviewWindow`來驅動更新流程。更新器會在您的視窗上呼叫`Show()`／`Close()`／`EmitEvent()`，而您的 HTML 決定要呈現的內容：

![一個「使用您自己的」更新器視窗，具有粉紅至橘色的漸層背景、單張白色圓角卡片、自訂字型排印，並由相同的更新器事件驅動畫面狀態。示範預設 UI 可以多麼徹底地被取代。](/assets/updater/byo-custom-window.png)

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

您的 HTML 會像內建範本一樣使用`window.wails.Events.On`／`Events.Emit`；只要視窗設有`AllowSimpleEventEmit: true`，框架自動注入的 shim 就會隨視窗提供，無論該視窗由框架還是由您擁有。API 請參閱[取代範本](#heading-9)。

@note{type="caution" title="BYO 更新器視窗必須設定 `AllowSimpleEventEmit`"}
基於安全性，框架會依此欄位管控`wails:event:emit:` postMessage 捷徑：未設定此欄位的視窗無法合成主機端自訂事件。更新器的自訂 HTML shim 會透過該捷徑發出`updater:user:*`事件，因此 BYO 視窗若忘記設定此欄位，所有按鈕點擊都會被無聲忽略——使用者點選「安裝」後不會發生任何事。

任何載入您無法完全控制之 HTML（遠端 URL、使用者提供的內容）的視窗，都應讓`AllowSimpleEventEmit`保持**關閉**。啟用時，頁面中的任何 JavaScript（包括 XSS 注入點）都能觸發任何`app.Event.On(name, …)`處理常式。此捷徑只能傳遞裸名稱（沒有承載資料），且無法存取繫結／Call 路徑；但如果 Go 程式碼中的具權限自訂事件處理常式僅依事件名稱執行動作，它仍可觸發這些處理常式。

框架的*內建*更新器視窗已在內部設定此欄位，只有 BYO 呼叫端需要記得設定。

@end

### 無介面模式

```go
app.Updater.Init(updater.Config{
    // …
    Window: updater.WindowNone,
})
```

系統完全不會開啟視窗。請從您自己的 UI（或現有的主視窗）訂閱`updater:*`事件，並從按鈕處理常式呼叫`app.Updater.CheckAndInstall(ctx)`。這適合只有在找到更新時才應顯示介面的定期背景檢查，也適合將更新流程整合至自訂設定面板的應用程式。

## 事件

Go 和 JavaScript 都透過標準 Wails 事件匯流排進行訂閱。**請勿手動輸入線路字串**——請使用更新器套件（Go）或執行階段套件（JS）匯出的常數。兩層共用同一組名稱，並由迴歸測試確保同步。

### 從 Go

這些常數位於`github.com/wailsapp/wails/v3/pkg/updater`。請透過`app.Event.On(name, fn)`訂閱；回呼會收到`*application.CustomEvent`，其`Data`欄位為[事件參考](#heading-14)中列出的具型別承載資料——請執行型別斷言，不要進行 JSON 解碼：

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

所有可用的 Go 常數：

| 常數 | 線路字串 |
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

### 從 JavaScript

這些常數位於`@wailsio/runtime`中的`Updater.Events`之下。名稱與 Go 相同，並依子命名空間（`User.*`、`Window.*`）組織，方便透過自動完成探索：

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

自訂 HTML*回傳*給主機的使用者動作事件位於`Updater.Events.User`之下：

```js
import { Updater } from "@wailsio/runtime";

document.getElementById("install-btn").addEventListener("click", () => {
    // The framework window does this internally via the postMessage shim;
    // shown here for BYO templates that need to drive the flow themselves.
    window._wails.invoke("wails:event:emit:" + Updater.Events.User.Install);
});
```

### 事件參考

訂閱端（主機 → 頁面）：

| 常數（Go） | 常數（JS） | 承載資料 | 時機 |
| --- | --- | --- | --- |
| `updater.EventCheckStarted` | `Updater.Events.CheckStarted` | 無 | 每次`Check`往返前 |
| `updater.EventUpdateAvailable` | `Updater.Events.UpdateAvailable` | `*Release` | `Check`找到較新的版本 |
| `updater.EventNoUpdate` | `Updater.Events.NoUpdate` | 無 | `Check`確認已是最新版本 |
| `updater.EventDownloadStarted` | `Updater.Events.DownloadStarted` | `*Release` | 開始串流位元組 |
| `updater.EventDownloadProgress` | `Updater.Events.DownloadProgress` | `Progress` | 下載期間約為10 Hz |
| `updater.EventDownloadComplete` | `Updater.Events.DownloadComplete` | `*Release` | 所有位元組均已寫入，尚未驗證 |
| `updater.EventVerifying` | `Updater.Events.Verifying` | `*Release` | 開始檢查簽章／摘要 |
| `updater.EventInstalling` | `Updater.Events.Installing` | `*Release` | 開始解壓縮及暫存 |
| `updater.EventUpdateReady` | `Updater.Events.UpdateReady` | `*Release` | 等待重新啟動 |
| `updater.EventError` | `Updater.Events.Error` | `ErrorInfo` | 任一階段失敗 |
| `updater.EventMeta` | `Updater.Events.Meta` | `Meta` | 每個工作階段一次，於重播快照之前 |

頁面端（頁面 → 主機）— 若您撰寫自訂範本，程式碼會訂閱以下事件：

| 常數（Go） | 常數（JS） | 時機 |
| --- | --- | --- |
| `updater.EventWindowReady` | `Updater.Events.Window.Ready` | 視窗載入完成；主機重播目前狀態 |
| `updater.EventUserInstall` | `Updater.Events.User.Install` | 處於`available`狀態時的主要動作 |
| `updater.EventUserRestart` | `Updater.Events.User.Restart` | 處於`ready`狀態時的主要動作 |
| `updater.EventUserSkip` | `Updater.Events.User.Skip` | 「略過此版本」 |
| `updater.EventUserRemind` | `Updater.Events.User.Remind` | 「稍後提醒我」 |
| `updater.EventUserCancel` | `Updater.Events.User.Cancel` | 關閉按鈕 |

## API 參考

### `updater.Config`

| 欄位 | 型別 | 說明 |
| --- | --- | --- |
| `CurrentVersion` | `string` | <strong>必填。</strong>與您為發行版本加上的標籤使用相同字串（不含`v`前綴） |
| `Providers` | `[]updater.Provider` | <strong>必填。</strong>依序排列的備援鏈 |
| `PublicKey` | `[]byte` | PEM 或原始位元組。選填，但若未提供，已簽署的發行版本會以安全失敗方式拒絕更新 |
| `CheckInterval` | `time.Duration` | 非零值會啟動背景輪詢迴圈，呼叫`CheckAndInstall` |
| `Platform` | `string` | 覆寫`runtime.GOOS`以選取資產 |
| `Arch` | `string` | 覆寫`runtime.GOARCH`以選取資產 |
| `Channel` | `string` | 目前僅供參考；由提供者進行通道篩選 |
| `Window` | `updater.WindowOption` | `nil`（內建預設值）、`&BuiltinWindow{…}`、`BYOWindow(handle)`或`WindowNone` |

### `*updater.Updater`的方法

| 簽章 | 用途 |
| --- | --- |
| `Init(cfg Config) error` | 進行設定。第二次呼叫時傳回`ErrAlreadyConfigured` |
| `State() State` | 目前的生命週期階段 |
| `CurrentVersion() string` | 傳入`Init`的版本 |
| `Check(ctx) (*Release, error)` | 依序巡訪提供者鏈。`(rel, nil)` = 找到、`(nil, nil)` = 已是最新版本、`(nil, err)` = 全部失敗 |
| `DownloadAndInstall(ctx) error` | 串流、驗證、解壓縮（若為封存檔）並暫存。必須先呼叫`Check` |
| `CheckAndInstall(ctx) error` | 便利方法：開啟視窗、執行`Check`，若找到更新則執行`DownloadAndInstall` |
| `Restart(ctx) error` | 啟動輔助程式、呼叫`Host.Quit`後結束；輔助程式會置換並重新啟動應用程式 |
| `DownloadedPath() string` | 暫存更新在磁碟上的位置；若無則為`""` |
| `SkipVersion(v string)` | 將`v`記錄為已略過；後續`Check`會將其視為已是最新版本 |
| `SkippedVersion() string` | 讀取目前略過的版本 |
| `StopPeriodicCheck()` | 取消由`Config.CheckInterval`啟動的計時器，並等待迴圈返回 |

### 錯誤

| 哨兵值 | 傳回來源 |
| --- | --- |
| `ErrAlreadyConfigured` | 首次成功後的`Init` |
| `ErrNotConfigured` | 在`Init`之前進行的任何操作 |
| `ErrNoPendingRelease` | 未先呼叫`Check`就呼叫`DownloadAndInstall` |
| `ErrDownloadInProgress` | 另一個流程執行期間呼叫`DownloadAndInstall` |
| `ErrNotReady` | 沒有已暫存的更新時呼叫`Restart` |

## 置換的運作方式

`Restart`會在設定哨兵環境變數後重新執行目前的二進位檔。`application.New`會在啟動時偵測這些變數，並轉入輔助程式模式：

1. 輔助程式最多等待30秒，讓父程序 PID 結束（`platformIsAlive`在 Windows 上透過`syscall.OpenProcess` + `GetExitCodeProcess`輪詢，在 Unix 上則透過`os.FindProcess` + `proc.Signal(syscall.Signal(0))`輪詢）。
2. 輔助程式會備份目標（檔案直接複製；macOS `.app`套件目錄則遞迴複製）。
3. 輔助程式會以暫存成品取代目標，最多重試20次，每次嘗試之間退避500毫秒：
  - **Unix** — `os.RemoveAll(target)` + `os.Rename(newPath, target)`。指向舊 inode 的已開啟檔案描述元仍然有效。
  - **Windows** — `os.Rename(target, target.old.<nanos>)` + `os.Rename(newPath, target)`。Windows 允許重新命名映像仍處於對應狀態的檔案，但不允許刪除；下次更新時，輔助程式會清除所有殘留的同層`.old.*`項目，但僅限其所屬核心對應已解除者。

4. 輔助程式會在新的二進位檔上還原原始可執行檔的模式（下載的檔案是使用預設 umask 建立，因此在 Unix 上會移除`+x`；在 Windows 上則不會進行任何操作）。
5. 輔助程式會清除輔助程式模式的環境變數，並重新啟動（此時已遭取代的）二進位檔。
6. 輔助程式結束。

如果啟動失敗，輔助程式會還原備份。如果父程序未在30秒內結束，輔助程式會在接觸目標前中止（因此，即使關機對話方塊阻擋`Quit`，使用者仍可保有能正常運作的應用程式）。

對於以`.zip`散布的 macOS `.app`套件（建議的封裝方式），封存檔會在驗證完成與準備就緒之間解壓縮，讓輔助程式能以實際目錄進行置換。

## 定期檢查

```go
app.Updater.Init(updater.Config{
    // …
    CheckInterval: 6 * time.Hour,
})
```

當`CheckInterval > 0`時，背景 goroutine 會依照設定的間隔呼叫`CheckAndInstall`。如果另一個流程已在進行中（檢查／下載／驗證／安裝），此時到達的計時觸發會被捨棄，因為不支援並行狀態機。

若要在背景靜默輪詢，且僅於找到更新時才顯示資訊，請設定`Window: updater.WindowNone`，並從您自己的 UI 回應`EventUpdateAvailable`。

## 略過與稍後提醒

預設視窗的「略過此版本」按鈕會透過`SkipVersion(rel.Version)`記錄可用版本。後續的`Check`會找到相同版本，並將其視為已是最新版本（直到使用者更新`CurrentVersion`；成功執行`Restart`後會自動更新）。「稍後提醒我」只會關閉視窗，不會記錄任何內容。

```go
// Reading what the user skipped (e.g. to surface in app settings)
if v := app.Updater.SkippedVersion(); v != "" {
    log.Printf("user skipped %s", v)
}

// Programmatically clearing the skip:
app.Updater.SkipVersion("")
```

## 散布檢查清單

發布將由更新程式安裝的版本之前：

1. <strong>選擇正確的封存格式。</strong>macOS：`.app`套件的`.zip`。Linux：單一二進位檔或`.tar.gz`。Windows：單一`.exe`或`.zip`。不支援`.dmg`／`.msi`／`.pkg`。
2. **簽署成品**時，請使用與`Config.PublicKey`相符的私密金鑰。對於由供應商發布的資訊提供（keygen.sh、AppCast），請遵循各供應商的簽署流程。對於使用`ChecksumAsset`的 GitHub Releases，請以`sha256sum`／`shasum -a 256`產生`SHA256SUMS`檔案。
3. **比對版本字串。**`Config.CurrentVersion`必須與該版本的版本標籤完全相符（例如`1.0.0` ↔ 標籤`v1.0.0`；開頭的`v`會在供應商端移除）。
4. **在目標平台上測試置換**，並在發布前至少測試一次。程式碼簽署、公證與 Gatekeeper 處理方式均依平台而異，更新程式本身不涵蓋這些項目。

## 疑難排解

**「signature requires a public key but none configured」** — 此版本有`Signature`欄位，但`Config.PublicKey`為空。請設定公開金鑰，或變更版本發布管線，使其不包含簽章。

**「digest mismatch」** — 下載的位元組與供應商承諾的內容不符。通常是下載不完整（網路短暫中斷）或成品損毀。重新執行通常可以解決問題。

**視窗開啟後立即消失，沒有 Markdown，也沒有進度資訊** — 您的自訂 HTML 未叫用`wails:runtime:ready`。請參閱[取代範本](#heading-9)的相容層。

**Windows 更新始終無法完成；輔助程式記錄顯示「remove old (attempt N): Access is denied」** — 這只會發生於此 PR 的`de764fb`之前版本；目前的實作會先重新命名並移至一旁，因此沒有這個問題。請升級。

**macOS Gatekeeper 封鎖已置換的二進位檔** — 必須在整個流程中保留程式碼簽章。簽署原始`.app`，*並且*在您的建置管線於更新時修改權利設定的情況下，重新簽署重新啟動的二進位檔。

## 另請參閱

- 可執行的範例：[`v3/examples/updater`](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/updater)
- 測試示範存放庫：[`wailsapp/updater-demo`](https://github.com/wailsapp/updater-demo)
- 教學：[為 Wails 應用程式新增自我更新功能](/tutorials/04-self-update-a-wails-app/)
