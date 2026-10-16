---
title: "自動啟動"
description: "註冊您的應用程式，使其在使用者登入 macOS、Windows 和 Linux 時啟動"
slug: "features/autostart/basics"
sourcePath: "features/autostart/basics.md"
---

## 自動啟動

`app.Autostart`會註冊您的應用程式，使其在使用者登入時自動啟動。它會為各平台選用適當的原生機制，並解析使用符號連結的安裝路徑（Homebrew、Scoop），避免二進位檔升級後註冊失效。

註冊會在<strong>下次登入</strong>時生效，而非立即生效。

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

使用預設選項註冊應用程式，使其在登入時啟動。

```go
func (m *AutostartManager) Enable() error
```

重複呼叫`Enable`是安全的——每次都會覆寫註冊。因此，如果您已保存使用者的偏好設定，可以在每次啟動時呼叫它。

### `EnableWithOptions`

使用自訂選項進行註冊。

```go
func (m *AutostartManager) EnableWithOptions(opts AutostartOptions) error
```

**`AutostartOptions`：**

| 欄位 | 型別 | 說明 |
| --- | --- | --- |
| `Identifier` | `string` | 覆寫自動衍生的註冊 ID。請參閱下方的「識別碼」。 |
| `Arguments` | `[]string` | 登入時啟動應用程式，附加於可執行檔路徑之後的額外引數（例如`--hidden`）。 |

### `Disable`

移除自動啟動註冊。若應用程式尚未註冊，則傳回`nil`——停用操作具有冪等性。

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
| `Path` | `string` | 註冊成品在磁碟上的位置（plist 路徑、`.desktop`路徑或登錄子機碼）。`Enabled`為 false 時為空。 |
| `Strategy` | `AutostartStrategy` | 用來註冊應用程式的機制（請參閱[平台行為](#heading-2)）。 |

## 平台行為

@tabs{sync-key="platform"}
[macOS]
根據應用程式的封裝方式，會使用以下兩種機制之一：

- **macOS 13+，已封裝的`.app`**：`SMAppService.mainAppService`。適用於沙箱化應用程式和 Mac App Store 組建。不會顯示 TCC 自動化提示（過去的 AppleScript 方法會觸發此提示）。
- **早於 macOS 13的版本，或未封裝的二進位檔**：將 LaunchAgent plist 寫入`~/Library/LaunchAgents/<identifier>.plist`，並使用`RunAtLoad=true`。

`Status()`會傳回`AutostartStrategySMAppService`或`AutostartStrategyLaunchAgent`，讓呼叫端能判斷採用了哪一種途徑。

應用程式從未封裝版本升級為已封裝版本時，`Status()`會檢查兩種途徑，而`Disable()`會清理它們，避免遺留的 LaunchAgent 繼續啟動舊組建。

[Windows]
在`HKCU\Software\Microsoft\Windows\CurrentVersion\Run`下新增登錄值，以自動啟動`Identifier`作為值名稱，並以加上引號的可執行檔路徑及任何`Arguments`作為資料。

引數引號處理遵循`CommandLineToArgvW`規則（引號前的反斜線會加倍），讓包含空格或引號的路徑能正確往返轉換。

`Status().Strategy`會傳回`AutostartStrategyRegistryRun`。

[Linux]
將 XDG 自動啟動項目寫入`$XDG_CONFIG_HOME/autostart/<identifier>.desktop`（預設為`~/.config/autostart/`），內容如下：

```ini
Type=Application
Hidden=false
X-GNOME-Autostart-enabled=true
Exec=<executable> <arguments>
```

`Exec`欄位會依照[freedesktop.org Desktop Entry 規格](https://specifications.freedesktop.org/desktop-entry-spec/)進行逸出——保留字元（`"`、`` ` ``、`$`、`\\`）會以反斜線逸出；若值包含空白字元，則會以雙引號括住。

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

若`Options.Identifier`為空，系統會根據您的應用程式名稱衍生預設值：

| 平台 | 預設識別碼 |
| --- | --- |
| macOS（已封裝） | 應用程式的套件識別碼，例如`com.example.MyApp` |
| macOS（未封裝） | `wails.autostart.<slug>`，其中`<slug>`衍生自`application.Options.Name` |
| Windows | `application.Options.Name`的 slug（轉為小寫、移除非`A-Za-z0-9._-`字元，並將空格轉為連字號） |
| Linux | 與 Windows 相同的 slug |

識別碼必須符合`^[A-Za-z0-9._-]+$`，且不得超過200個字元。macOS 建議使用反向 DNS 格式（這符合 launchd Label 的慣用寫法）。

覆寫`AutostartOptions.Identifier`時，同一個識別碼會在 Windows 上重複用作登錄值名稱，並在 Linux 上用作`.desktop`檔名，因此只需一個字串即可跨平台識別該項註冊。

## 過時項目偵測

`Disable()`和`Status()`會透過<strong>將已註冊的可執行檔路徑與`os.Executable()`比對（解析所有符號連結後）</strong>來找出註冊項目，而不是查詢識別碼。這表示：

- <strong>在不同版本之間變更識別碼是安全的。</strong>只要可執行檔路徑相同，`Status()`仍可找出舊的註冊項目，並由`Disable()`將其清除。
- <strong>位於不同路徑的第二份應用程式副本不會覆寫第一份副本的註冊項目。</strong>每個二進位檔位置都會分別追蹤。
- <strong>透過符號連結安裝的版本（Homebrew、Scoop）具有穩定性。</strong>比對前會先對`os.Executable()`套用`filepath.EvalSymlinks`，因此 Homebrew 升級時即使替換連結目標，也不會留下失效的項目。

此機制<em>不會</em>處理以下情況：如果使用者將二進位檔移至不相關的路徑或重新命名，舊的註冊項目會成為孤立項目（指向目前已不存在的檔案）。從穩定安裝路徑發佈的應用程式不必擔心此問題；以可攜式單一檔案二進位檔形式發佈的應用程式，應在移動自身之前呼叫`Disable()`，或一律透過穩定的符號連結啟動。

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

[`examples/autostart/`](https://github.com/wailsapp/wails/tree/master/v3/examples/autostart)中提供了包含狀態、啟用及停用按鈕，且可完整執行的範例。
