---
title: "可自行更新的 Wails 應用程式"
description: "建置可透過 GitHub Releases 自行更新的 Wails v3應用程式——涵蓋從`wails3 init`、已簽署版本驗證到輔助程式模式置換的完整流程。"
slug: "tutorials/04-self-update-a-wails-app"
sourcePath: "tutorials/04-self-update-a-wails-app.md"
---

在本教學中，你將為全新的 Wails v3應用程式加入應用程式內更新程式。完成後，應用程式將能：

- 依需求檢查 GitHub Releases（也可選擇透過計時器定期檢查）。
- 下載適用於目前執行中作業系統與架構的資產。
- 根據下載的位元組驗證 SHA-256摘要（也可選擇驗證 Ed25519 簽章）。
- 在框架的預設更新視窗中顯示版本資訊。
- 置換執行中的二進位檔並重新啟動，完全不必隨附個別的輔助執行檔。

我們將使用<strong>GitHub Releases</strong>作為更新來源，因為它免費且不需要任何基礎架構。同樣的模式也適用於[keygen.sh](/guides/updater/#keygensh--updaterproviderskeygen)和[Sparkle AppCast](/guides/updater/#sparkle-appcast--updaterprovidersappcast)。完成後，請參閱[更新程式指南](/guides/updater/)。

@note{type="tip" title="事前準備"}
- Go 1.25或更新版本
- 已安裝`wails3` CLI（`go install github.com/wailsapp/wails/v3/cmd/wails3@latest`）
- 一個可供你推送發行版本的 GitHub 儲存庫
- 熟悉[QR Code 服務教學](/tutorials/01-creating-a-service/)會有所幫助，但並非必要條件

@end

<br/>

@steps
### 從全新的 Wails 應用程式開始
使用 vanilla 範本建立新專案的骨架：

```bash
wails3 init -n updater-tutorial -t vanilla
cd updater-tutorial
```

現在應該會有一個包含`main.go`、`frontend/`和`Taskfile.yml`的目錄。確認專案可以建置並啟動：

```bash
wails3 task dev
```

此時應開啟一個空白的 Wails 視窗。結束應用程式後繼續。

### 加入更新程式的匯入項目
開啟`main.go`，並將兩個更新程式套件加入匯入項目：

```go {title="main.go" ins="6-7"}
package main

import (
    _ "embed"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/updater"
    "github.com/wailsapp/wails/v3/pkg/updater/providers/github"
)
```

這些套件會引入 Updater 本身及 GitHub Releases 提供者。

### 設定 Updater
每個`*application.App`都已接好`app.Updater`，你只需呼叫`Init`：

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

將此程式碼放在`application.New`之後、`app.Run()`之前。

@note{type="note" title="版本字串格式"}
傳入與發行版本標籤相同的版本，但<strong>不要</strong>包含開頭的`v`。提供者會自行移除標籤名稱中的`v`。此處的`1.0.0` ↔ GitHub 上的`v1.0.0`。

@end

### 加入觸發更新的選單項目
在同一個`main.go`中加入「檢查更新…」選單項目：

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

`CheckAndInstall`會開啟框架的更新視窗、執行`Check`，若找到發行版本，則自動執行`DownloadAndInstall`。沒有新版本時，視窗會停留在「已是最新版本」狀態，由使用者按下<strong>關閉</strong>按鈕將其關閉。

@note{type="caution" title="在 goroutine 中執行"}
`CheckAndInstall`會阻塞，直到驗證與安裝完成為止。直接從選單點擊事件呼叫它會阻塞 UI 執行緒。請將它包在`go func()`中。

@end

### 在沒有發行版本的情況下執行一次
```bash
wails3 task dev
```

按一下<strong>App → 檢查更新…</strong>。更新視窗應會短暫開啟、呼叫 GitHub API、找不到比`1.0.0`更新的發行版本，最後停留在帶有綠色 ✓ 的<strong>已是最新版本</strong>狀態。

如果此處發生錯誤，通常是以下原因之一：

| 症狀 | 修正方式 |
| --- | --- |
| `404 Not Found` | `Repository`欄位不正確——必須是`owner/repo` |
| `403 rate-limited` | 將`Token: "ghp_…"`加入 github.Config（使用具備`public_repo`範圍的 PAT） |
| 網路錯誤 | 確認執行中的應用程式可以連線至`api.github.com` |

### 發布測試用發行版本
將`main.go`中的`currentVersion`遞增至`1.0.0`（或維持不變）。針對一個平台進行建置，以取得可附加至發行版本的二進位檔：

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

在二進位檔旁產生一個`SHA256SUMS`檔案：

```bash
cd bin
shasum -a 256 updater-tutorial-* > SHA256SUMS
cat SHA256SUMS
```

你應該會看到一或多行類似以下的內容：

```
abc123…  updater-tutorial-darwin-arm64.zip
```

現在將其作為<strong>v2.0.0</strong>發布至你的 GitHub 儲存庫：

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

@note{type="note" title="資產命名"}
預設資產比對器會依檔名中的`GOOS`與`GOARCH`子字串挑選資產。只要資產名稱包含`darwin`（或`linux`／`windows`）以及`arm64`（或`amd64`／`386`），比對器就能找到它。若要使用自訂比對器，請參閱[更新程式指南](/guides/updater/#github-releases--updaterprovidersgithub)。

@end

### 執行應用程式並驗證更新
在`currentVersion`仍為`1.0.0`時，再次執行應用程式：

```bash
wails3 task dev
```

按一下<strong>App → 檢查更新…</strong>。這次應該會看到類似以下的內容：

![處於「更新已就緒」狀態的預設更新程式視窗，其中顯示版本標籤、以 Markdown 轉譯的版本資訊，以及「重新啟動並套用」主要按鈕。](/assets/updater/default-window-ready.png)

- 主圖示會從藍色 ↓（「有可用更新」）變成綠色 ✓（「更新已就緒」）。
- 副標題顯示`v1.0.0 → v2.0.0 · <size>`。
- 版本資訊面板會轉譯你的 Markdown，包括粗體、行內程式碼和表格。
- 下載期間進度列會逐漸填滿（二進位檔很小，因此速度會很快）。

Updater 會將新的二進位檔暫存於臨時目錄。若要完成更新：

- 按一下<strong>重新啟動並套用</strong>。
- 應用程式會結束，輔助程式會置換二進位檔，接著重新啟動新的二進位檔。
- 重新啟動的應用程式會回報`currentVersion = "1.0.0"`（因為我們將它寫死了），但磁碟上的位元組與v2.0.0組建相符。

在實際的應用程式中，`currentVersion`會在建置時透過`-ldflags`設定，使新的二進位檔知道自己目前是v2.0.0，而後續檢查也不會找到更新。

### 將`currentVersion`連結至建置流程
將常數替換為建置時變數：

```go {title="main.go" ins="2,4"}
var (
    currentVersion = "dev" // overridden by -ldflags at release time
)
```

接著在建置命令中：

```bash
wails3 task build:darwin -- -ldflags "-X main.currentVersion=2.0.0"
```

或者，將`-ldflags`加入`Taskfile.yml`，使其從`git describe --tags`取得值。

### 加入密碼學簽章（建議用於正式環境）
SHA256SUMS 路徑會驗證<em>完整性</em>（位元組與 GitHub 所儲存的內容相符），但不會驗證<em>真實性</em>（亦即這些位元組是由你的發行流程產生，而不是來自遭入侵的維護者帳號）。為了防止竄改，請使用 Ed25519 金鑰簽署每個發行版本：

```bash
# One-time: generate the keypair
ssh-keygen -t ed25519 -f updater-key -N "" -C "wails-updater"
#   updater-key      — keep secret (build server, HSM, password manager)
#   updater-key.pub  — bundle in your app
```

每次發行時，請使用私鑰簽署每個資產的 SHA-256 摘要。以下是一個小型 Go 輔助程式：

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

預設的 GitHub 提供者目前不會擷取個別的簽章檔案；你可以[撰寫自訂提供者](/guides/updater/#writing-your-own-provider)來擷取，或改用 **keygen.sh**。後者會在伺服器端簽署每個成品，並透過其 API 提供摘要與簽章。

將公開金鑰嵌入你的應用程式：

```go {title="main.go" ins="1,6"}
//go:embed updater-key.pub
var updaterPublicKey []byte

app.Updater.Init(updater.Config{
    CurrentVersion: currentVersion,
    PublicKey:      updaterPublicKey,
    Providers:      []updater.Provider{gh},
})
```

設定`PublicKey`後，任何附帶`Signature`的發行版本都必須通過此金鑰的驗證。發行來源無法替換成自己的金鑰；這正是在建置時透過頻外方式固定金鑰的目的。

### 自訂視窗
預設視窗適用於常見情境。如果需要更多控制，可依照你想自訂的程度，從以下三種方式中選擇一種：

@tabs
[僅使用 CSS]
```go
app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{
        CSS: `:root { --accent: #ff6f00; --radius: 16px; }`,
    },
})
```

如需完整的變數清單，請參閱[透過 CSS 變數設定主題](/guides/updater/#theme-via-css-variables)一節。

[自訂 HTML]
```go
//go:embed updater-window.html
var updaterHTML string

app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{HTML: updaterHTML},
})
```

你的 HTML 必須訂閱`updater:*`事件，並透過 Wails 事件通道發出`updater:user:*`動作。如需 JS 轉接程式，請參閱[替換範本](/guides/updater/#replace-the-template)。

[使用你自己的視窗]
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

如果你已有自己的視窗基礎架構，並希望更新程式驅動該視窗，而不是另開一個視窗，這種方式會很實用。以下是一個完全自訂的 HTML 範本，由與預設範本相同的更新程式事件驅動：

![一個自備的更新程式視窗，採用粉紅至橘色的漸層背景與自訂圓角卡片版面配置，展示預設 UI 可以完全替換。](/assets/updater/byo-custom-window.png)

@note{type="caution" title="必須啟用 `AllowSimpleEventEmit`"}
更新程式的自訂 HTML 轉接程式會透過`wails:event:emit:` postMessage 捷徑驅動「安裝／略過／稍後提醒／重新啟動」動作。基於安全考量，必須啟用此欄位才能使用該捷徑。若忘記啟用，這些按鈕將毫無反應。請勿在載入非由你完全控制之 HTML 的視窗上啟用此功能；如需瞭解威脅模型，請參閱指南中的[使用你自己的視窗](/guides/updater/#bring-your-own-window)一節。

@end

@end

### 在背景執行自動檢查
若要依計時器檢查，而不是透過選單點擊來檢查（或兩者並用）：

```go {ins="5"}
app.Updater.Init(updater.Config{
    CurrentVersion: currentVersion,
    Providers:      []updater.Provider{gh},
    PublicKey:      updaterPublicKey,
    CheckInterval:  6 * time.Hour,
})
```

每次計時器觸發時，都會執行與手動點擊相同的`CheckAndInstall`流程。如果你希望定期檢查在實際找到更新前保持靜默，請設定`Window: updater.WindowNone`，然後自行訂閱`EventUpdateAvailable`，以決定要顯示何種使用者體驗。

@end

## 大功告成

你現在已有一個具備以下功能的 Wails 應用程式：

- 可依需求或按照計時器檢查 GitHub Releases 中的更新。
- 在精美的預設視窗中將版本資訊呈現為 Markdown。
- 使用你發布的 SHA-256 摘要驗證下載內容。
- 可選擇使用建置時嵌入的公開金鑰來驗證 Ed25519 簽章。
- 原地替換執行中的二進位檔，並自動重新啟動。

## 後續步驟

- [更新程式指南](/guides/updater/)包含完整的 API 參考、所有事件、所有設定選項，以及輔助模式的替換機制。
- 請查看[`v3/examples/updater`](https://github.com/wailsapp/wails/tree/master/v3/examples/updater)，其中提供可供複製的完整運作範例。
- 測試目標儲存庫[`wailsapp/updater-demo`](https://github.com/wailsapp/updater-demo)展示了建議的發行資產配置。

## 正式環境中的注意事項

- **macOS 上的程式碼簽署** — Gatekeeper 要求替換後的二進位檔必須經過簽署與公證。在將`.app`套件壓縮成發行用 zip 檔<em>之前</em>，請先簽署套件。更新程式會逐位元組原樣保留內容，不會重新簽署任何項目。
- **Windows 上的防毒軟體** — 從網際網路下載且未簽署的`.exe`檔案可能會觸發 SmartScreen 警告。請使用 Authenticode 憑證簽署二進位檔，否則必須接受使用受嚴格管控之電腦的使用者可能需要將你的應用程式加入允許清單。
- **不可分割的發行作業** — 請一起發布`SHA256SUMS`與二進位檔，不要分成不同的提交。更新程式會分別擷取附屬檔案與二進位檔；若兩者不同步，摘要檢查將以失敗方式安全終止。
- **略過的版本** — 預設視窗中的「略過此版本」按鈕會將略過設定記錄在本機。如果發布重大安全性更新，請使用新的版本號碼，以免先前略過某個發行版本的使用者也自動略過此更新。
