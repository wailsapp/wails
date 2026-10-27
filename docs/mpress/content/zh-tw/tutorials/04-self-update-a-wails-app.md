---
title: "可自我更新的 Wails 應用程式"
description: "建立一個可從 GitHub Releases 自我更新的 Wails v3應用程式——從`wails3 init`一路到已簽署版本的驗證與輔助模式置換。"
slug: "tutorials/04-self-update-a-wails-app"
sourcePath: "tutorials/04-self-update-a-wails-app.md"
---

在本教學中，你將為全新的 Wails v3應用程式加入內建更新程式。完成後，應用程式將能：

- 依需求檢查 GitHub Releases（也可選擇透過計時器定期檢查）。
- 下載符合目前執行中作業系統與架構的資產。
- 根據下載的位元組驗證 SHA-256摘要（也可選擇驗證 Ed25519 簽章）。
- 在框架的預設更新視窗中顯示版本資訊。
- 置換執行中的二進位檔並重新啟動，而且完全不必另外隨附輔助執行檔。

我們將使用**GitHub Releases**作為更新來源，因為它免費且不需要任何基礎設施。同樣的模式也適用於[keygen.sh](/guides/updater/#keygensh--updaterproviderskeygen)和[Sparkle AppCast](/guides/updater/#sparkle-appcast--updaterprovidersappcast)。完成後，請參閱[更新程式指南](/guides/updater/)。

@note{type="tip" title="必要條件"}
- Go 1.26或更新版本
- 已安裝`wails3` CLI（`go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0`）
- 一個你可以推送版本的 GitHub 儲存庫
- 熟悉[QR Code Service 教學](/tutorials/01-creating-a-service/)會有幫助，但不是必要條件

@end

<br/>

@steps
### 從全新的 Wails 應用程式開始
使用 vanilla 範本建立新專案的基本架構：

```bash
wails3 init -n updater-tutorial -t vanilla
cd updater-tutorial
```

現在應該會有一個包含`main.go`、`frontend/`及`Taskfile.yml`的目錄。確認它能成功建置並啟動：

```bash
wails3 task dev
```

此時應該會開啟一個空白的 Wails 視窗。結束應用程式後繼續。

### 加入更新程式的匯入項目
開啟`main.go`，將兩個更新程式套件加入匯入項目：

```go {title="main.go" ins="6-7"}
package main

import (
    _ "embed"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/updater"
    "github.com/wailsapp/wails/v3/pkg/updater/providers/github"
)
```

這些套件會引入更新程式本身及 GitHub Releases 提供者。

### 設定更新程式
每個`*application.App`都已連接`app.Updater`，你只需要呼叫`Init`：

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
傳入與版本標籤相同的版本號，但**不要**包含開頭的`v`。提供者會在其端點移除標籤名稱中的`v`。此處的`1.0.0` ↔ GitHub 上的`v1.0.0`。

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

`CheckAndInstall`會開啟框架的更新視窗、執行`Check`，並在找到新版本時自動執行`DownloadAndInstall`。若沒有新版本，視窗會保持開啟並顯示「已是最新版本」狀態；使用者可按**關閉**按鈕將其關閉。

@note{type="caution" title="在 goroutine 中執行"}
`CheckAndInstall`會阻塞，直到驗證和安裝完成。若直接從選單點擊處呼叫它，將會阻塞 UI 執行緒。請用`go func()`包裝。

@end

### 在沒有任何版本時執行一次
```bash
wails3 task dev
```

按一下**應用程式 → 檢查更新…**。你應該會看到更新視窗短暫開啟、呼叫 GitHub API、找不到比`1.0.0`更新的版本，最後停留在帶有綠色 ✓ 的**已是最新版本**狀態。

如果此處發生錯誤，通常是下列其中一種情況：

| 症狀 | 修正方式 |
| --- | --- |
| `404 Not Found` | `Repository`欄位不正確——必須是`owner/repo` |
| `403 rate-limited` | 將`Token: "ghp_…"`加入 github.Config（使用具備`public_repo`範圍的 PAT） |
| 網路錯誤 | 確認執行中的應用程式可以連線至`api.github.com` |

### 發布測試版本
將`main.go`中的`currentVersion`提高到`1.0.0`（也可以維持不變）。針對一個平台進行建置，以取得可附加至版本的二進位檔：

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

在二進位檔旁產生`SHA256SUMS`檔案：

```bash
cd bin
shasum -a 256 updater-tutorial-* > SHA256SUMS
cat SHA256SUMS
```

你應該會看到一行或多行類似以下的內容：

```
abc123…  updater-tutorial-darwin-arm64.zip
```

現在將它發布至你的 GitHub 儲存庫，版本為**v2.0.0**：

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
預設資產比對器會依據檔名中的`GOOS`與`GOARCH`子字串選取資產。只要資產名稱包含`darwin`（或`linux`／`windows`）和`arm64`（或`amd64`／`386`），比對器就能找到它。如需自訂比對器，請參閱[更新程式指南](/guides/updater/#github-releases--updaterprovidersgithub)。

@end

### 執行應用程式並驗證更新
在`currentVersion`仍為`1.0.0`時，再次執行應用程式：

```bash
wails3 task dev
```

按一下**應用程式 → 檢查更新…**。這次應該會看到類似以下的內容：

![處於「更新已就緒」狀態的預設更新程式視窗，其中顯示版本標籤、以 Markdown 呈現的版本資訊，以及主要的「重新啟動並套用」按鈕。](/assets/updater/default-window-ready.png)

- 主圖示會從藍色 ↓（「有可用更新」）變成綠色 ✓（「更新已就緒」）。
- 副標題會顯示`v1.0.0 → v2.0.0 · <size>`。
- 版本資訊面板會呈現你的 Markdown，包括粗體、行內程式碼和表格。
- 下載期間進度列會逐漸填滿（二進位檔很小，因此下載速度會很快）。

更新程式會將新的二進位檔暫存到臨時目錄。若要完成更新：

- 按一下**重新啟動並套用**。
- 應用程式會結束，輔助程式會置換二進位檔，接著重新啟動新的二進位檔。
- 重新啟動的應用程式會回報`currentVersion = "1.0.0"`（因為我們將它寫死了），但磁碟上的位元組會與v2.0.0建置版本相符。

在實際的應用程式中，`currentVersion`會在建置時透過`-ldflags`設定，讓新的二進位檔知道自己現在是v2.0.0，而後續檢查將找不到可用更新。

### 將`currentVersion`連接至建置流程
將常數改為建置時變數：

```go {title="main.go" ins="2,4"}
var (
    currentVersion = "dev" // overridden by -ldflags at release time
)
```

接著在建置命令中：

```bash
wails3 task build:darwin -- -ldflags "-X main.currentVersion=2.0.0"
```

或者，將`-ldflags`加入`Taskfile.yml`，讓它從`git describe --tags`取得值。

### 加入密碼編譯簽章（建議用於正式環境）
SHA256SUMS 路徑會驗證*完整性*（位元組與 GitHub 儲存的內容相符），但不會驗證*真實性*（這些位元組是由你的發布管線產生，而不是來自遭入侵的維護者帳號）。為了防止竄改，請使用 Ed25519 金鑰簽署每個版本：

```bash
# One-time: generate the keypair
ssh-keygen -t ed25519 -f updater-key -N "" -C "wails-updater"
#   updater-key      — keep secret (build server, HSM, password manager)
#   updater-key.pub  — bundle in your app
```

每次發佈時，請使用私鑰簽署每個資產的 SHA-256 摘要。以下是一個小型 Go 輔助工具：

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

預設的 GitHub 提供者目前不會擷取獨立的簽章檔案；你可以[撰寫自訂提供者](/guides/updater/#writing-your-own-provider)來擷取，或改用**keygen.sh**。它會在伺服器端簽署每個成品，並透過其 API 提供摘要和簽章。

將公開金鑰嵌入應用程式：

```go {title="main.go" ins="1,6"}
//go:embed updater-key.pub
var updaterPublicKey []byte

app.Updater.Init(updater.Config{
    CurrentVersion: currentVersion,
    PublicKey:      updaterPublicKey,
    Providers:      []updater.Provider{gh},
})
```

設定`PublicKey`後，任何隨附`Signature`的發佈版本都必須使用此金鑰通過驗證。發佈來源無法替換成自己的金鑰；這正是在建置時透過頻道外方式固定金鑰的用意。

### 自訂視窗
預設視窗適用於常見情境。如果需要更多控制，可使用以下三種替代方案；請依所需的自訂程度選擇一種：

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

如需完整的變數清單，請參閱[透過 CSS 變數設定佈景主題](/guides/updater/#theme-via-css-variables)一節。

[自訂 HTML]
```go
//go:embed updater-window.html
var updaterHTML string

app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{HTML: updaterHTML},
})
```

你的 HTML 必須訂閱`updater:*`事件，並透過 Wails 事件通道發出`updater:user:*`動作。JS 墊片請參閱[替換範本](/guides/updater/#replace-the-template)。

[使用自己的視窗]
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

如果你已有自己的視窗基礎架構，並希望更新程式驅動該視窗，而不是另開視窗，這個方式會很實用。以下是一個完全自訂的 HTML 範本，它與預設視窗一樣由相同的更新程式事件驅動：

![一個採用粉紅至橘色漸層背景與自訂圓角卡片版面的自備更新程式視窗，展示預設 UI 可完全替換。](/assets/updater/byo-custom-window.png)

@note{type="caution" title="必須啟用 `AllowSimpleEventEmit`"}
更新程式的自訂 HTML 墊片會透過`wails:event:emit:` postMessage 捷徑驅動「安裝／略過／稍後提醒／重新啟動」動作；為了安全起見，只有啟用此欄位才能使用該捷徑。若忘記啟用，按鈕將毫無反應，也不會顯示錯誤。請勿在載入非你完全掌控之 HTML 的視窗上啟用此功能；威脅模型請參閱指南的[使用自己的視窗](/guides/updater/#bring-your-own-window)一節。

@end

@end

### 在背景自動執行檢查
若要依計時器檢查，而不是（或除了）點選選單，請使用：

```go {ins="5"}
app.Updater.Init(updater.Config{
    CurrentVersion: currentVersion,
    Providers:      []updater.Provider{gh},
    PublicKey:      updaterPublicKey,
    CheckInterval:  6 * time.Hour,
})
```

每次計時器觸發時，都會執行與手動點選相同的`CheckAndInstall`流程。若希望定期檢查在實際找到更新前保持靜默，請設定`Window: updater.WindowNone`，然後自行訂閱`EventUpdateAvailable`，以決定要呈現的使用者體驗。

@end

## 完成

你現在已有一個能夠執行以下操作的 Wails 應用程式：

- 依需求及定時檢查 GitHub Releases 是否有更新。
- 在精緻的預設視窗中，將發佈說明呈現為 Markdown。
- 使用你發佈的 SHA-256 摘要驗證下載內容。
- 可選擇使用建置時嵌入的公開金鑰驗證 Ed25519 簽章。
- 就地替換正在執行的二進位檔，並自動重新啟動。

## 後續步驟

- [更新程式指南](/guides/updater/)包含完整的 API 參考、所有事件與設定選項，以及輔助程式模式的替換機制。
- 請參閱[`v3/examples/updater`](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples/updater)，其中提供可供複製的完整運作範例。
- 測試目標儲存庫[`wailsapp/updater-demo`](https://github.com/wailsapp/updater-demo)展示了建議的發佈資產配置。

## 正式環境中的注意事項

- **macOS 上的程式碼簽署** — Gatekeeper 要求替換後的二進位檔必須經過簽署與公證。請在將`.app`套件壓縮成發佈檔案*之前*完成簽署。更新程式會原封不動保留其位元組，不會重新簽署任何內容。
- **Windows 上的防毒軟體** — 從網際網路下載且未簽署的`.exe`檔案可能會觸發 SmartScreen 警告。請使用 Authenticode 憑證簽署二進位檔，否則就必須接受受嚴格管控之電腦的使用者可能需要將你的應用程式加入允許清單。
- **不可分割的發佈** — 請一起發佈`SHA256SUMS`（以及你的二進位檔），不要分成不同的提交。更新程式會分別擷取附屬檔案與二進位檔；若兩者不同步，摘要檢查將以安全方式失敗並拒絕更新。
- **略過的版本** — 預設視窗的「略過此版本」按鈕會在本機記錄略過狀態。如果你發佈重大安全性更新，請為其指定新的版本號碼，以免先前略過某個發佈版本的使用者也自動略過此更新。
