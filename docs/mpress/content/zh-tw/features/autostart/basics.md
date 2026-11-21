---
title: "自動啟動"
description: "將應用程式註冊為在 macOS、Windows 和 Linux 使用者登入時啟動"
slug: "features/autostart/basics"
sourcePath: "features/autostart/basics.md"
---

## 自動啟動

`app.Autostart`會將您的應用程式註冊為在使用者登入時自動啟動。它會依平台選用正確的原生機制，並解析使用符號連結的安裝路徑（Homebrew、Scoop），避免升級二進位檔後註冊失效。

註冊會在**下次登入**時生效，而非立即生效。

## 快速開始

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

使用預設選項將應用程式註冊為在登入時啟動。

```go
func (m *AutostartManager) Enable() error
```

重複呼叫`Enable`是安全的——每次都會覆寫註冊。因此，如果您已持久儲存使用者的偏好設定，在每次啟動時呼叫它是正確的做法。

### `EnableWithOptions`

使用自訂選項進行註冊。

```go
func (m *AutostartManager) EnableWithOptions(opts AutostartOptions) error
```

**`AutostartOptions`：**

| 欄位 | 型別 | 說明 |
| --- | --- | --- |
| `Identifier` | `string` | 覆寫自動衍生的註冊 ID。請參閱下方的「識別碼」。 |
| `Arguments` | `[]string` | 登入時啟動應用程式時，附加至可執行檔路徑後方的額外引數（例如`--hidden`）。 |

### `Disable`

移除自動啟動註冊。如果應用程式尚未註冊，則傳回`nil`——停用操作具有冪等性。

```go
func (m *AutostartManager) Disable() error
```

### `IsEnabled`

回報註冊是否存在。此操作速度快，但不會驗證已註冊的路徑。

```go
func (m *AutostartManager) IsEnabled() (bool, error)
```

### `Status`

傳回完整的註冊狀態。

```go
func (m *AutostartManager) Status() (AutostartStatus, error)
```

**`AutostartStatus`：**

| 欄位 | 型別 | 說明 |
| --- | --- | --- |
| `Enabled` | `bool` | 註冊是否存在。 |
| `Path` | `string` | 註冊產物在磁碟上的位置（plist 路徑、`.desktop`路徑或登錄子機碼）。`Enabled`為 false 時為空。 |
| `Strategy` | `AutostartStrategy` | 用來註冊應用程式的機制（請參閱[平台行為](#heading-2)）。 |

## 平台行為

@tabs{sync-key="platform"}
[macOS]
依應用程式的封裝方式，會使用以下兩種機制之一：

- **macOS 13+，已封裝的`.app`**：`SMAppService.mainAppService`。適用於沙箱化應用程式和 Mac App Store 組建。不會出現 TCC 自動化提示（舊有的 AppleScript 方法會觸發此提示）。
- **早於13的 macOS，或未封裝的二進位檔**：將含有`RunAtLoad=true`的 LaunchAgent plist 寫入`~/Library/LaunchAgents/<identifier>.plist`。

`Status()`會傳回`AutostartStrategySMAppService`或`AutostartStrategyLaunchAgent`，讓呼叫端得知採用了哪一種方式。

應用程式從未封裝升級為已封裝時，`Status()`會檢查這兩種路徑，而`Disable()`會清理它們，避免孤立的 LaunchAgent 繼續啟動舊版組建。

[Windows]
在`HKCU\Software\Microsoft\Windows\CurrentVersion\Run`下新增登錄值，以自動啟動`Identifier`作為值名稱，並以加上引號的可執行檔路徑及任何`Arguments`作為資料。

引數引號遵循`CommandLineToArgvW`規則（引號前的反斜線加倍），因此含有空格或引號的路徑能正確往返轉換。

`Status().Strategy`會傳回`AutostartStrategyRegistryRun`。

[Linux]
將 XDG 自動啟動項目寫入`$XDG_CONFIG_HOME/autostart/<identifier>.desktop`（預設為`~/.config/autostart/`），其中包含：

```ini
Type=Application
Hidden=false
X-GNOME-Autostart-enabled=true
Exec=<executable> <arguments>
```

`Exec`欄位會依照[freedesktop.org Desktop Entry 規範](https://specifications.freedesktop.org/desktop-entry-spec/)進行逸出——保留字元（`"`、`` ` ``、`$`、`\\`）會以反斜線逸出；值含有空白字元時，則會以雙引號括住。

`Status().Strategy`會傳回`AutostartStrategyXDGAutostart`。

[iOS／Android／伺服器]
不支援。所有方法都會傳回`ErrAutostartNotSupported`。請使用`errors.Is(err, application.ErrAutostartNotSupported)`明確偵測此情況：

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

## 識別碼

如果`Options.Identifier`為空，系統會根據應用程式名稱衍生預設值：

| 平台 | 預設識別碼 |
| --- | --- |
| macOS（已封裝） | 應用程式的套件識別碼，例如`com.example.MyApp` |
| macOS（未封裝） | `wails.autostart.<slug>`，其中`<slug>`衍生自`application.Options.Name` |
| Windows | `application.Options.Name`的 slug（轉為小寫、移除非`A-Za-z0-9._-`字元，並將空格替換為連字號） |
| Linux | 與 Windows 使用相同的 slug |

識別碼必須符合`^[A-Za-z0-9._-]+$`，且不得超過200個字元。macOS 建議使用反向 DNS 格式（這符合 launchd Label 的慣用寫法）。

覆寫`AutostartOptions.Identifier`時，同一識別碼也會用作 Windows 上的登錄值名稱和 Linux 上的`.desktop`檔名，因此只需一個字串即可跨平台識別該註冊。

## 過時註冊偵測

`Disable()`和`Status()`不會查詢識別碼，而是透過**將已註冊的可執行檔路徑與`os.Executable()`（解析所有符號連結後）比對**來找出註冊。這表示：

- <strong>在不同版本之間變更識別碼是安全的。</strong>只要可執行檔路徑相同，`Status()`仍可找到舊註冊，而`Disable()`會將其清理。
- <strong>位於不同路徑的第二份應用程式不會覆寫第一份的註冊。</strong>每個二進位檔位置都會獨立追蹤。
- <strong>使用符號連結的安裝（Homebrew、Scoop）可保持穩定。</strong>比對前會對`os.Executable()`套用`filepath.EvalSymlinks`，因此 Homebrew 升級替換連結目標時，不會讓該項目失去對應的檔案。

此機制*不*涵蓋以下情況：如果使用者將二進位檔移動或重新命名至不相關的路徑，舊的註冊項目會變成孤立項目（它仍指向目前已不存在的檔案）。從固定安裝路徑發布的應用程式無須擔心此問題；以可攜式單一檔案二進位檔形式發布的應用程式，應在移動自身之前呼叫`Disable()`，或一律透過固定的符號連結啟動。

## 範例

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

包含狀態／啟用／停用按鈕且可完整執行的範例位於[`examples/autostart/`](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/autostart)。
