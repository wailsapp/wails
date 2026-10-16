---
title: "自訂 URL 通訊協定"
description: "註冊自訂 URL 通訊協定名稱，讓連結能啟動您的應用程式"
slug: "guides/distribution/custom-protocols"
sourcePath: "guides/distribution/custom-protocols.md"
---

自訂 URL 通訊協定（也稱為 URL 通訊協定名稱）可讓使用者點選採用自訂通訊協定的連結（例如`myapp://action`或`myapp://open/document`）時啟動您的應用程式。

## 概觀

自訂通訊協定可用於：

- **深層連結**：使用特定資料啟動您的應用程式
- **瀏覽器整合**：處理來自網頁的連結
- **電子郵件連結**：從電子郵件用戶端開啟您的應用程式
- **應用程式間通訊**：從其他應用程式啟動

**範例**：`myapp://open/document?id=123`會啟動您的應用程式，並開啟文件123。

## 設定

在應用程式選項中定義自訂通訊協定：

自訂通訊協定是在`build/config.yml`中宣告（各平台的封裝工具會在封裝時使用此設定，包括 Windows 上的 NSIS 巨集、MSIX 資訊清單、macOS 的`CFBundleURLTypes`，以及 Linux 的`.desktop`/`xdg-mime`）。`application.Protocol`型別並不存在，`application.Options`上也沒有`Protocols`欄位。

```yaml
# build/config.yml
protocols:
  - scheme: myapp
    description: "My Application Protocol"
```

在 Go 程式碼中，透過`ApplicationLaunchedWithUrl`事件監聽由 URL 觸發的啟動：

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

func main() {
    app := application.New(application.Options{
        Name:        "My Application",
        Description: "My awesome application",
    })

    app.Event.OnApplicationEvent(events.Common.ApplicationLaunchedWithUrl, func(e *application.ApplicationEvent) {
        handleCustomURL(e.Context().URL())
    })

    app.Run()
}

func handleCustomURL(url string) {
    // Parse and handle the custom URL
    // Example: myapp://open/document?id=123
    println("Received URL:", url)
}
```

## 通訊協定處理常式

監聽通訊協定事件以處理傳入的 URL：

```go
app.Event.OnApplicationEvent(events.Common.ApplicationLaunchedWithUrl, func(e *application.ApplicationEvent) {
    url := e.Context().URL()

    // Parse the URL
    parsedURL, err := parseCustomURL(url)
    if err != nil {
        app.Logger.Error("Failed to parse URL:", err)
        return
    }

    // Handle different actions
    switch parsedURL.Action {
    case "open":
        openDocument(parsedURL.DocumentID)
    case "settings":
        showSettings()
    case "user":
        showUser Profile(parsedURL.UserID)
    default:
        app.Logger.Warn("Unknown action:", parsedURL.Action)
    }
})
```

## URL 結構

設計清楚且具階層性的 URL 結構：

```
myapp://action/resource?param=value

Examples:
myapp://open/document?id=123
myapp://settings/theme?mode=dark
myapp://user/profile?username=john
```

**最佳做法：**

- 通訊協定名稱使用小寫字母
- 通訊協定名稱應簡短且容易記憶
- 使用階層式路徑表示資源
- 使用查詢參數傳遞選用資料
- 對特殊字元進行 URL 編碼

## 平台註冊

各平台註冊自訂通訊協定的方式不同。

@tabs{sync-key="platform"}
[Windows]
### Windows NSIS 安裝程式

使用 NSIS 安裝程式時，**Wails v3 會自動註冊自訂通訊協定**。

#### 自動註冊

使用`wails3 build`建置應用程式時，NSIS 安裝程式會：

1. 自動註冊`build/config.yml`中`protocols:`鍵下宣告的所有通訊協定
2. 將通訊協定與應用程式的可執行檔建立關聯
3. 建立正確的登錄項目
4. 解除安裝時移除通訊協定關聯

**不需要任何額外設定！**

#### 運作方式

NSIS 範本包含內建巨集：

- `wails.associateCustomProtocols`－在安裝期間註冊通訊協定
- `wails.unassociateCustomProtocols`－在解除安裝期間移除通訊協定

系統會根據您的`Protocols`設定自動呼叫這些巨集。

#### 手動設定登錄（進階）

如果需要手動註冊（不透過 NSIS）：

```batch
@echo off
REM Register custom protocol
REG ADD "HKEY_CURRENT_USER\SOFTWARE\Classes\myapp" /ve /d "URL:My Application Protocol" /f
REG ADD "HKEY_CURRENT_USER\SOFTWARE\Classes\myapp" /v "URL Protocol" /t REG_SZ /d "" /f
REG ADD "HKEY_CURRENT_USER\SOFTWARE\Classes\myapp\shell\open\command" /ve /d "\"%1\"" /f
```

#### 測試

測試通訊協定註冊：

```powershell
# Open protocol URL from PowerShell
Start-Process "myapp://test/action"

# Or from command prompt
start myapp://test/action
```

### Windows MSIX 套件

使用 MSIX 封裝時，也會自動註冊自訂通訊協定。

#### 自動註冊

使用 MSIX 建置應用程式時，資訊清單會自動納入`build/config.yml`通訊協定設定中的通訊協定註冊資訊。

產生的資訊清單包含：

```xml
<uap:Extension Category="windows.protocol">
  <uap:Protocol Name="myapp">
    <uap:DisplayName>My Application Protocol</uap:DisplayName>
  </uap:Protocol>
</uap:Extension>
```

#### 通用連結（網頁至應用程式連結）

Windows 支援<strong>網頁至應用程式連結</strong>，運作方式類似 macOS 上的通用連結。將應用程式部署為 MSIX 套件時，您可以啟用 HTTPS 連結，直接啟動應用程式。

@note{type="note"}
網頁至應用程式連結需要手動設定資訊清單。自訂通訊協定名稱會依據`build/config.yml`自動設定，但關聯網域必須手動加入您的 MSIX 資訊清單。

@end

若要啟用網頁至應用程式連結，請依照[Microsoft 的網頁至應用程式連結指南](https://learn.microsoft.com/en-us/windows/apps/develop/launch/web-to-app-linking)操作。您需要：

1. **手動將 App URI Handler 加入 MSIX 資訊清單**（`build/windows/msix/app_manifest.xml`）：
  ```xml
  <uap3:Extension Category="windows.appUriHandler">
    <uap3:AppUriHandler>
      <uap3:Host Name="myawesomeapp.com"/>
    </uap3:AppUriHandler>
  </uap3:Extension>
  ```


2. <strong>在您的網站上設定`windows-app-web-link`：</strong>在`https://myawesomeapp.com/.well-known/windows-app-web-link`代管`windows-app-web-link`檔案。此檔案應包含應用程式的套件資訊及其處理的路徑。

當網頁至應用程式連結啟動您的應用程式時，您會收到與透過自訂通訊協定啟動時相同的`ApplicationLaunchedWithUrl`事件。

[macOS]
### Info.plist 設定

在 macOS 上，通訊協定是透過您的`Info.plist`檔案註冊。

#### 自動設定

使用`wails3 build`建置時，Wails 會自動產生包含通訊協定的`Info.plist`。

在`build/config.yml`中宣告的通訊協定會新增至：

```xml
<key>CFBundleURLTypes</key>
<array>
    <dict>
        <key>CFBundleURLName</key>
        <string>My Application Protocol</string>
        <key>CFBundleURLSchemes</key>
        <array>
            <string>myapp</string>
        </array>
        <key>CFBundleTypeRole</key>
        <string>Editor</string>
    </dict>
</array>
```

#### 測試

```bash
# Open protocol URL from terminal
open "myapp://test/action"

# Check registered handlers
/System/Library/Frameworks/CoreServices.framework/Versions/A/Frameworks/LaunchServices.framework/Versions/A/Support/lsregister -dump | grep myapp
```

#### 通用連結

除了自訂通訊協定外，macOS 也支援<strong>通用連結</strong>，可讓一般 HTTPS 連結（例如`https://myawesomeapp.com/path`）啟動您的應用程式。通用連結可在您的網頁與桌面應用程式之間提供流暢的使用者體驗。

@note{type="caution"}
通用連結要求您的 macOS 應用程式使用有效的 Apple Developer 憑證和佈建描述檔進行<strong>程式碼簽署</strong>。未簽署或臨時簽署的建置版本無法開啟通用連結。測試前，請確保您的應用程式已正確簽署。

@end

若要啟用通用連結，請依照[在應用程式中支援通用連結的 Apple 指南](https://developer.apple.com/documentation/xcode/supporting-universal-links-in-your-app)操作。您需要：

1. 在`entitlements.plist`中<strong>新增權利</strong>：
  ```xml
  <key>com.apple.developer.associated-domains</key>
  <array>
    <string>applinks:myawesomeapp.com</string>
  </array>
  ```


2. **將 NSUserActivityTypes 新增至 Info.plist**：
  ```xml
  <key>NSUserActivityTypes</key>
  <array>
    <string>NSUserActivityTypeBrowsingWeb</string>
  </array>
  ```


3. <strong>在您的網站上設定`apple-app-site-association`：</strong>將`apple-app-site-association`檔案託管於`https://myawesomeapp.com/.well-known/apple-app-site-association`。

當通用連結觸發您的應用程式時，您會收到相同的`ApplicationLaunchedWithUrl`事件，因此可使用與自訂通訊協定完全相同的處理程式碼。

[Linux]
### 桌面項目

在 Linux 上，通訊協定是透過`.desktop`檔案註冊的。

#### 自動設定

使用`wails3 build`建置時，Wails 會產生包含通訊協定處理常式的桌面項目檔案。

**已在 v3 修正**：Linux 桌面範本現在會正確包含通訊協定處理功能。

產生的桌面檔案包含：

```ini
[Desktop Entry]
Type=Application
Name=My Application
Exec=/usr/bin/myapp %u
MimeType=x-scheme-handler/myapp;
```

#### 手動註冊

如有需要，請手動安裝桌面檔案：

```bash
# Copy desktop file
cp myapp.desktop ~/.local/share/applications/

# Update desktop database
update-desktop-database ~/.local/share/applications/

# Register protocol handler
xdg-mime default myapp.desktop x-scheme-handler/myapp
```

#### 測試

```bash
# Open protocol URL
xdg-open "myapp://test/action"

# Check registered handler
xdg-mime query default x-scheme-handler/myapp
```

@end

## 完整範例

以下是處理多個通訊協定動作的完整範例：

```go
package main

import (
    "fmt"
    "net/url"
    "strings"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

type App struct {
    app    *application.App
    window *application.WebviewWindow
}

func main() {
    // Protocol registration lives in build/config.yml (the platform packagers
    // consume it); the application code just listens for the launch event.
    app := application.New(application.Options{
        Name:        "DeepLink Demo",
        Description: "Custom protocol demonstration",
    })

    myApp := &App{app: app}
    myApp.setup()

    app.Run()
}

func (a *App) setup() {
    // Create window
    a.window = a.app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:  "DeepLink Demo",
        Width:  800,
        Height: 600,
        URL:    "http://wails.localhost/",
    })

    // Handle custom protocol URLs
    a.app.Event.OnApplicationEvent(events.Common.ApplicationLaunchedWithUrl, func(e *application.ApplicationEvent) {
        a.handleDeepLink(e.Context().URL())
    })
}

func (a *App) handleDeepLink(rawURL string) {
    // Parse URL
    parsedURL, err := url.Parse(rawURL)
    if err != nil {
        a.app.Logger.Error("Failed to parse URL:", err)
        return
    }

    // Bring window to front
    a.window.Show()
    a.window.Focus()

    // Extract path and query
    path := strings.Trim(parsedURL.Path, "/")
    query := parsedURL.Query()

    // Handle different actions
    parts := strings.Split(path, "/")
    if len(parts) == 0 {
        return
    }

    action := parts[0]

    switch action {
    case "open":
        if len(parts) >= 2 {
            resource := parts[1]
            id := query.Get("id")
            a.openResource(resource, id)
        }

    case "settings":
        section := ""
        if len(parts) >= 2 {
            section = parts[1]
        }
        a.openSettings(section)

    case "user":
        if len(parts) >= 2 {
            username := parts[1]
            a.openUserProfile(username)
        }

    default:
        a.app.Logger.Warn("Unknown action:", action)
    }
}

func (a *App) openResource(resourceType, id string) {
    fmt.Printf("Opening %s with ID: %s\n", resourceType, id)
    // Emit event to frontend
    a.app.Event.Emit("navigate", map[string]string{
        "type": resourceType,
        "id":   id,
    })
}

func (a *App) openSettings(section string) {
    fmt.Printf("Opening settings section: %s\n", section)
    a.app.Event.Emit("navigate", map[string]string{
        "page":    "settings",
        "section": section,
    })
}

func (a *App) openUserProfile(username string) {
    fmt.Printf("Opening user profile: %s\n", username)
    a.app.Event.Emit("navigate", map[string]string{
        "page": "user",
        "user": username,
    })
}
```

## 前端整合

在前端處理導覽事件：

```javascript
import { Events } from '@wailsio/runtime'

// Listen for navigation events from protocol handler
Events.On('navigate', (event) => {
    const { type, id, page, section, user } = event.data

    if (type === 'document') {
        // Open document with ID
        router.push(`/document/${id}`)
    } else if (page === 'settings') {
        // Open settings
        router.push(`/settings/${section}`)
    } else if (page === 'user') {
        // Open user profile
        router.push(`/user/${user}`)
    }
})
```

## 安全性考量

### 驗證所有輸入

務必驗證並清理來自外部來源的 URL：

```go
func (a *App) handleDeepLink(rawURL string) {
    // Parse URL
    parsedURL, err := url.Parse(rawURL)
    if err != nil {
        a.app.Logger.Error("Invalid URL:", err)
        return
    }

    // Validate scheme
    if parsedURL.Scheme != "myapp" {
        a.app.Logger.Warn("Invalid scheme:", parsedURL.Scheme)
        return
    }

    // Validate path
    path := strings.Trim(parsedURL.Path, "/")
    if !isValidPath(path) {
        a.app.Logger.Warn("Invalid path:", path)
        return
    }

    // Sanitize parameters
    params := sanitizeQueryParams(parsedURL.Query())

    // Process validated URL
    a.processDeepLink(path, params)
}

func isValidPath(path string) bool {
    // Only allow alphanumeric and forward slashes
    validPath := regexp.MustCompile(`^[a-zA-Z0-9/]+$`)
    return validPath.MatchString(path)
}

func sanitizeQueryParams(query url.Values) map[string]string {
    sanitized := make(map[string]string)
    for key, values := range query {
        if len(values) > 0 {
            // Take first value and sanitize
            sanitized[key] = sanitizeString(values[0])
        }
    }
    return sanitized
}
```

### 防止注入攻擊

絕不可將 URL 直接當作程式碼或 SQL 執行：

```go
// ❌ DON'T: Execute URL content
func badHandler(url string) {
    exec.Command("sh", "-c", url).Run() // DANGEROUS!
}

// ✅ DO: Parse and validate
func goodHandler(url string) {
    parsed, _ := url.Parse(url)
    action := parsed.Query().Get("action")

    // Whitelist allowed actions
    allowed := map[string]bool{
        "open":     true,
        "settings": true,
        "help":     true,
    }

    if allowed[action] {
        handleAction(action)
    }
}
```

## 測試

### 手動測試

在開發期間測試通訊協定處理常式：

**Windows：**

```powershell
Start-Process "myapp://test/action?id=123"
```

**macOS：**

```bash
open "myapp://test/action?id=123"
```

**Linux：**

```bash
xdg-open "myapp://test/action?id=123"
```

### HTML 測試

建立測試用 HTML 頁面：

```html
<!DOCTYPE html>
<html>
<head>
    <title>Protocol Test</title>
</head>
<body>
    <h1>Custom Protocol Test Links</h1>

    <ul>
        <li><a href="myapp://open/document?id=123">Open Document 123</a></li>
        <li><a href="myapp://settings/theme?mode=dark">Dark Mode Settings</a></li>
        <li><a href="myapp://user/profile?username=john">User Profile</a></li>
    </ul>
</body>
</html>
```

## 疑難排解

### 通訊協定未註冊

**Windows：**

- 檢查登錄：`HKEY_CURRENT_USER\SOFTWARE\Classes\<scheme>`
- 使用 NSIS 安裝程式重新安裝
- 確認安裝程式以適當權限執行

**macOS：**

- 使用`wails3 build`重新建置應用程式
- 檢查應用程式套件中的`Info.plist`：`MyApp.app/Contents/Info.plist`
- 重設 Launch Services：`/System/Library/Frameworks/CoreServices.framework/Frameworks/LaunchServices.framework/Support/lsregister -kill`

**Linux：**

- 檢查桌面檔案：`~/.local/share/applications/myapp.desktop`
- 更新資料庫：`update-desktop-database ~/.local/share/applications/`
- 確認處理常式：`xdg-mime query default x-scheme-handler/myapp`

### 應用程式未啟動

**檢查記錄：**

```go
app := application.New(application.Options{
    LogLevel: slog.LevelDebug, // requires `import "log/slog"`
    // ...
})
```

**常見問題：**

- 應用程式未安裝在預期位置
- 註冊資訊中的可執行檔路徑與實際位置不符
- 權限問題

## 最佳實務

### ✅ 建議做法

- **使用描述清楚的通訊協定名稱**——使用`mycompany-myapp`，而非`mca`
- **驗證所有輸入**——絕不要信任來自外部來源的 URL
- **妥善處理錯誤**——記錄無效的 URL，不要讓程式當機
- **提供使用者回饋**——顯示觸發了哪項動作
- **在所有平台上測試**——通訊協定的處理方式各有不同
- **記錄 URL 結構**——協助使用者和整合開發人員

### ❌ 請勿這樣做

- **請勿使用常見的通訊協定名稱**——避免使用`http`、`file`、`app`等名稱
- **請勿將 URL 當作程式碼執行**——這會造成極大的安全風險
- **請勿公開敏感操作**——破壞性動作必須要求確認
- **請勿假設通訊協定在所有環境中都能運作**——應提供備援機制
- **別忘了進行 URL 編碼**——正確處理特殊字元

## 後續步驟

- [Windows 封裝](/guides/build/windows/)——瞭解 NSIS 安裝程式選項
- [檔案關聯](/guides/file-associations/)——使用您的應用程式開啟檔案
- [單一執行個體](/guides/single-instance/)——防止應用程式同時執行多個執行個體

---

<strong>有問題嗎？</strong>請到 [Discord](https://discord.gg/JDdSxwjhGf) 提問，或查看[範例](https://github.com/wailsapp/wails/tree/master/v3/examples)。
