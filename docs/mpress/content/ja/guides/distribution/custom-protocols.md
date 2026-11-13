---
title: "カスタム URL プロトコル"
description: "カスタム URL スキームを登録し、リンクからアプリケーションを起動する"
slug: "guides/distribution/custom-protocols"
sourcePath: "guides/distribution/custom-protocols.md"
---

カスタム URL プロトコル（URL スキームとも呼ばれます）を使用すると、ユーザーが `myapp://action` や `myapp://open/document` などのカスタムプロトコルを持つリンクをクリックしたときに、アプリケーションを起動できます。

## 概要

カスタムプロトコルでは、次のことが可能です。

- **ディープリンク**：特定のデータを渡してアプリを起動する
- **ブラウザー連携**：Web ページからのリンクを処理する
- **メール内リンク**：メールクライアントからアプリを開く
- **アプリ間通信**：ほかのアプリケーションから起動する

**例**：`myapp://open/document?id=123` はアプリを起動し、ドキュメント 123 を開きます。

## 設定

アプリケーションオプションでカスタムプロトコルを定義します。

カスタムプロトコルは `build/config.yml` で宣言します。各プラットフォームのパッケージャー（Windows の NSIS マクロ、MSIX マニフェスト、macOS の `CFBundleURLTypes`、Linux の `.desktop`/`xdg-mime`）は、パッケージ作成時にこの宣言を使用します。`application.Protocol` 型は存在せず、`application.Options` に `Protocols` フィールドもありません。

```yaml
# build/config.yml
protocols:
  - scheme: myapp
    description: "My Application Protocol"
```

Go コードでは、`ApplicationLaunchedWithUrl` イベントを介して URL による起動を待ち受けます。

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

## プロトコルハンドラー

受信した URL を処理するには、プロトコルイベントを待ち受けます。

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

## URL 構造

明確な階層構造を持つ URL を設計します。

```
myapp://action/resource?param=value

Examples:
myapp://open/document?id=123
myapp://settings/theme?mode=dark
myapp://user/profile?username=john
```

**ベストプラクティス：**

- スキーム名には小文字を使用する
- スキームは短く覚えやすいものにする
- リソースには階層的なパスを使用する
- 任意指定のデータにはクエリパラメーターを使用する
- 特殊文字を URL エンコードする

## プラットフォームへの登録

カスタムプロトコルの登録方法は、プラットフォームごとに異なります。

@tabs{sync-key="platform"}
[Windows]
### Windows NSIS インストーラー

NSIS インストーラーを使用すると、**Wails v3 はカスタムプロトコルを自動的に登録します**。

#### 自動登録

`wails3 build` を使用してアプリケーションをビルドすると、NSIS インストーラーは次の処理を行います。

1. `build/config.yml` で宣言されたすべてのプロトコルを、`protocols:` キーの配下に自動登録する
2. プロトコルをアプリケーションの実行ファイルに関連付ける
3. 適切なレジストリエントリを設定する
4. アンインストール時にプロトコルの関連付けを削除する

**追加の設定は不要です。**

#### 仕組み

NSIS テンプレートには、次の組み込みマクロが含まれています。

- `wails.associateCustomProtocols` - インストール時にプロトコルを登録する
- `wails.unassociateCustomProtocols` - アンインストール時にプロトコルを削除する

これらのマクロは、`Protocols` の設定に基づいて自動的に呼び出されます。

#### レジストリへの手動登録（上級者向け）

NSIS を使用せず、手動で登録する必要がある場合：

```batch
@echo off
REM Register custom protocol
REG ADD "HKEY_CURRENT_USER\SOFTWARE\Classes\myapp" /ve /d "URL:My Application Protocol" /f
REG ADD "HKEY_CURRENT_USER\SOFTWARE\Classes\myapp" /v "URL Protocol" /t REG_SZ /d "" /f
REG ADD "HKEY_CURRENT_USER\SOFTWARE\Classes\myapp\shell\open\command" /ve /d "\"%1\"" /f
```

#### テスト

プロトコルが登録されていることをテストします。

```powershell
# Open protocol URL from PowerShell
Start-Process "myapp://test/action"

# Or from command prompt
start myapp://test/action
```

### Windows MSIX パッケージ

MSIX パッケージを使用する場合も、カスタムプロトコルは自動的に登録されます。

#### 自動登録

MSIX を使用してアプリケーションをビルドすると、`build/config.yml` のプロトコル設定に基づくプロトコル登録がマニフェストに自動的に追加されます。

生成されるマニフェストには、次の内容が含まれます。

```xml
<uap:Extension Category="windows.protocol">
  <uap:Protocol Name="myapp">
    <uap:DisplayName>My Application Protocol</uap:DisplayName>
  </uap:Protocol>
</uap:Extension>
```

#### Universal Links（Web-to-App リンク）

Windows は、macOS の Universal Links と同様に機能する **Web-to-App リンク**をサポートしています。アプリケーションを MSIX パッケージとしてデプロイする場合、HTTPS リンクからアプリを直接起動できるように設定できます。

@note{type="note"}
Web-to-App リンクには、マニフェストの手動設定が必要です。カスタムプロトコルスキームは `build/config.yml` から自動的に設定されますが、関連付けるドメインは MSIX マニフェストに手動で追加する必要があります。

@end

Web-to-App リンクを有効にするには、[Web-to-App リンクに関する Microsoft のガイド](https://learn.microsoft.com/en-us/windows/apps/develop/launch/web-to-app-linking)に従ってください。次の作業が必要です。

1. **MSIX マニフェストに App URI Handler を手動で追加する**（`build/windows/msix/app_manifest.xml`）：
  ```xml
  <uap3:Extension Category="windows.appUriHandler">
    <uap3:AppUriHandler>
      <uap3:Host Name="myawesomeapp.com"/>
    </uap3:AppUriHandler>
  </uap3:Extension>
  ```


2. **Web サイトで `windows-app-web-link` を設定する：** `https://myawesomeapp.com/.well-known/windows-app-web-link` に `windows-app-web-link` ファイルをホストします。このファイルには、アプリのパッケージ情報と、そのアプリが処理するパスを記載する必要があります。

Web-to-App リンクによってアプリケーションが起動されると、カスタムプロトコルスキームの場合と同じ `ApplicationLaunchedWithUrl` イベントを受信します。

[macOS]
### Info.plist の設定

macOS では、`Info.plist` ファイルを介してプロトコルを登録します。

#### 自動設定

`wails3 build` を使用してビルドすると、Wails はプロトコルを含む `Info.plist` を自動的に生成します。

`build/config.yml` で宣言したプロトコルは、次の場所に追加されます。

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

#### テスト

```bash
# Open protocol URL from terminal
open "myapp://test/action"

# Check registered handlers
/System/Library/Frameworks/CoreServices.framework/Versions/A/Frameworks/LaunchServices.framework/Versions/A/Support/lsregister -dump | grep myapp
```

#### Universal Links

macOS は、カスタムプロトコルスキームに加えて、通常の HTTPS リンク（例：`https://myawesomeapp.com/path`）からアプリを起動できる **Universal Links** もサポートしています。Universal Links を使用すると、Web アプリとデスクトップアプリの間でシームレスなユーザー体験を提供できます。

@note{type="caution"}
Universal Links を使用するには、有効な Apple Developer 証明書とプロビジョニングプロファイルを使用して、macOS アプリを **コード署名**する必要があります。未署名またはアドホック署名されたビルドでは、Universal Links を開けません。テストする前に、アプリが適切に署名されていることを確認してください。

@end

Universal Links を有効にするには、[アプリで Universal Links をサポートするための Apple のガイド](https://developer.apple.com/documentation/xcode/supporting-universal-links-in-your-app)に従ってください。次の作業が必要です。

1. `entitlements.plist` に**エンタイトルメントを追加する**：
  ```xml
  <key>com.apple.developer.associated-domains</key>
  <array>
    <string>applinks:myawesomeapp.com</string>
  </array>
  ```


2. **Info.plist に NSUserActivityTypes を追加する**：
  ```xml
  <key>NSUserActivityTypes</key>
  <array>
    <string>NSUserActivityTypeBrowsingWeb</string>
  </array>
  ```


3. **Web サイトで `apple-app-site-association` を設定する：** `https://myawesomeapp.com/.well-known/apple-app-site-association` に `apple-app-site-association` ファイルをホストします。

Universal Link によってアプリが起動されると、同じ `ApplicationLaunchedWithUrl` イベントを受信するため、処理コードはカスタムプロトコルスキームの場合と同一になります。

[Linux]
### デスクトップエントリ

Linux では、`.desktop` ファイルを介してプロトコルを登録します。

#### 自動設定

`wails3 build` を使用してビルドすると、Wails はプロトコルハンドラーを含むデスクトップエントリファイルを生成します。

**v3**で修正済み：Linux デスクトップテンプレートにプロトコル処理が正しく含まれるようになりました。

生成されるデスクトップファイルには、次の内容が含まれます：

```ini
[Desktop Entry]
Type=Application
Name=My Application
Exec=/usr/bin/myapp %u
MimeType=x-scheme-handler/myapp;
```

#### 手動登録

必要に応じて、デスクトップファイルを手動でインストールします：

```bash
# Copy desktop file
cp myapp.desktop ~/.local/share/applications/

# Update desktop database
update-desktop-database ~/.local/share/applications/

# Register protocol handler
xdg-mime default myapp.desktop x-scheme-handler/myapp
```

#### テスト

```bash
# Open protocol URL
xdg-open "myapp://test/action"

# Check registered handler
xdg-mime query default x-scheme-handler/myapp
```

@end

## 完全な例

複数のプロトコルアクションを処理する完全な例を次に示します：

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

## フロントエンドとの統合

フロントエンドでナビゲーションイベントを処理します：

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

## セキュリティ上の考慮事項

### すべての入力を検証する

外部ソースから受け取った URL は、必ず検証してサニタイズしてください：

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

### インジェクション攻撃を防止する

URL をコードや SQL として直接実行しないでください：

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

## テスト

### 手動テスト

開発中にプロトコルハンドラーをテストします：

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

### HTML でのテスト

テスト用の HTML ページを作成します：

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

## トラブルシューティング

### プロトコルが登録されていない

**Windows：**

- レジストリを確認します：`HKEY_CURRENT_USER\SOFTWARE\Classes\<scheme>`
- NSIS インストーラーで再インストールします
- インストーラーが適切な権限で実行されたことを確認します

**macOS：**

- `wails3 build`を使用してアプリケーションを再ビルドします
- アプリバンドル内の`Info.plist`を確認します：`MyApp.app/Contents/Info.plist`
- Launch Services をリセットします：`/System/Library/Frameworks/CoreServices.framework/Frameworks/LaunchServices.framework/Support/lsregister -kill`

**Linux：**

- デスクトップファイルを確認します：`~/.local/share/applications/myapp.desktop`
- データベースを更新します：`update-desktop-database ~/.local/share/applications/`
- ハンドラーを確認します：`xdg-mime query default x-scheme-handler/myapp`

### アプリケーションが起動しない

**ログを確認します：**

```go
app := application.New(application.Options{
    LogLevel: slog.LevelDebug, // requires `import "log/slog"`
    // ...
})
```

**よくある問題：**

- アプリケーションが想定した場所にインストールされていない
- 登録されている実行ファイルのパスが実際の場所と一致していない
- 権限に問題がある

## ベストプラクティス

### ✅ 推奨事項

- **わかりやすいスキーム名を使用する** — `mca`ではなく`mycompany-myapp`を使用します
- **すべての入力を検証する** — 外部ソースからの URL を決して信頼しないでください
- **エラーを適切に処理する** — 無効な URL はログに記録し、クラッシュさせないでください
- **ユーザーにフィードバックを提供する** — どのアクションがトリガーされたかを表示します
- **すべてのプラットフォームでテストする** — プロトコルの処理方法はプラットフォームによって異なります
- **URL 構造を文書化する** — ユーザーやインテグレーターが利用しやすいようにします

### ❌ 禁止事項

- **一般的なスキーム名を使用しない** — `http`、`file`、`app`などは避けてください
- **URL をコードとして実行しない** — 非常に大きなセキュリティリスクがあります
- **機密性の高い操作を公開しない** — 破壊的な操作には確認を必須としてください
- **プロトコルがどこでも動作すると思い込まない** — フォールバック手段を用意してください
- **URL エンコーディングを忘れない** — 特殊文字を適切に処理してください

## 次のステップ

- [Windows 向けパッケージング](/guides/build/windows/) — NSIS インストーラーのオプションについて学びます
- [ファイルの関連付け](/guides/file-associations/) — アプリでファイルを開きます
- [シングルインスタンス](/guides/single-instance/) — アプリの複数インスタンスが起動しないようにします

---

**ご質問がありますか？** [Discord](https://discord.gg/JDdSxwjhGf)で質問するか、[サンプル](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples)を確認してください。
