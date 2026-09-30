---
title: "自定义 URL 协议"
description: "注册自定义 URL 方案，以便通过链接启动应用程序"
slug: "guides/distribution/custom-protocols"
sourcePath: "guides/distribution/custom-protocols.md"
---

自定义 URL 协议（也称为 URL 方案）允许用户点击使用自定义协议的链接（例如`myapp://action`或`myapp://open/document`）时启动您的应用程序。

## 概述

自定义协议支持：

- **深层链接**：使用特定数据启动应用
- **浏览器集成**：处理来自网页的链接
- **电子邮件链接**：从电子邮件客户端打开应用
- **应用间通信**：从其他应用程序启动应用

**示例**：`myapp://open/document?id=123`会启动应用并打开文档123。

## 配置

在应用程序选项中定义自定义协议：

自定义协议在`build/config.yml`中声明（各平台的打包工具会在打包时使用这些声明，包括 Windows 上的 NSIS 宏、MSIX 清单、macOS 的`CFBundleURLTypes`以及 Linux 的`.desktop`/`xdg-mime`）。不存在`application.Protocol`类型，`application.Options`上也不存在`Protocols`字段。

```yaml
# build/config.yml
protocols:
  - scheme: myapp
    description: "My Application Protocol"
```

在 Go 代码中，通过`ApplicationLaunchedWithUrl`事件监听使用 URL 启动应用的操作：

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

## 协议处理程序

监听协议事件以处理传入的 URL：

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

## URL 结构

设计清晰且层次分明的 URL 结构：

```
myapp://action/resource?param=value

Examples:
myapp://open/document?id=123
myapp://settings/theme?mode=dark
myapp://user/profile?username=john
```

**最佳实践：**

- 方案名称使用小写字母
- 方案名称应简短易记
- 对资源使用分层路径
- 使用查询参数传递可选数据
- 对特殊字符进行 URL 编码

## 平台注册

各平台采用不同的方式注册自定义协议。

@tabs{sync-key="platform"}
[Windows]
### Windows NSIS 安装程序

使用 NSIS 安装程序时，**Wails v3 会自动注册自定义协议**。

#### 自动注册

使用`wails3 build`构建应用程序时，NSIS 安装程序会：

1. 自动注册`build/config.yml`的`protocols:`键下声明的所有协议
2. 将协议与应用程序可执行文件关联
3. 设置正确的注册表项
4. 卸载时移除协议关联

**无需额外配置！**

#### 工作原理

NSIS 模板包含内置宏：

- `wails.associateCustomProtocols` - 在安装期间注册协议
- `wails.unassociateCustomProtocols` - 在卸载期间移除协议

系统会根据您的`Protocols`配置自动调用这些宏。

#### 手动配置注册表（高级）

如果需要手动注册（不使用 NSIS）：

```batch
@echo off
REM Register custom protocol
REG ADD "HKEY_CURRENT_USER\SOFTWARE\Classes\myapp" /ve /d "URL:My Application Protocol" /f
REG ADD "HKEY_CURRENT_USER\SOFTWARE\Classes\myapp" /v "URL Protocol" /t REG_SZ /d "" /f
REG ADD "HKEY_CURRENT_USER\SOFTWARE\Classes\myapp\shell\open\command" /ve /d "\"%1\"" /f
```

#### 测试

测试协议注册：

```powershell
# Open protocol URL from PowerShell
Start-Process "myapp://test/action"

# Or from command prompt
start myapp://test/action
```

### Windows MSIX 包

使用 MSIX 打包时，也会自动注册自定义协议。

#### 自动注册

使用 MSIX 构建应用程序时，清单会自动包含`build/config.yml`协议配置中的协议注册信息。

生成的清单包含：

```xml
<uap:Extension Category="windows.protocol">
  <uap:Protocol Name="myapp">
    <uap:DisplayName>My Application Protocol</uap:DisplayName>
  </uap:Protocol>
</uap:Extension>
```

#### 通用链接（网页到应用链接）

Windows 支持<strong>网页到应用链接</strong>，其工作方式与 macOS 上的通用链接类似。将应用程序部署为 MSIX 包时，可以启用 HTTPS 链接来直接启动应用。

@note{type="note"}
网页到应用链接需要手动配置清单。自定义协议方案会根据`build/config.yml`自动配置，但必须手动将关联域添加到 MSIX 清单中。

@end

要启用网页到应用链接，请参阅[Microsoft 网页到应用链接指南](https://learn.microsoft.com/en-us/windows/apps/develop/launch/web-to-app-linking)。您需要：

1. **手动将 App URI Handler 添加到 MSIX 清单**（`build/windows/msix/app_manifest.xml`）：
  ```xml
  <uap3:Extension Category="windows.appUriHandler">
    <uap3:AppUriHandler>
      <uap3:Host Name="myawesomeapp.com"/>
    </uap3:AppUriHandler>
  </uap3:Extension>
  ```


2. <strong>在网站上配置`windows-app-web-link`：</strong>在`https://myawesomeapp.com/.well-known/windows-app-web-link`托管`windows-app-web-link`文件。此文件应包含应用的软件包信息及其处理的路径。

当网页到应用链接启动您的应用程序时，您会收到与自定义协议方案相同的`ApplicationLaunchedWithUrl`事件。

[macOS]
### Info.plist 配置

在 macOS 上，通过`Info.plist`文件注册协议。

#### 自动配置

使用`wails3 build`构建时，Wails 会自动生成包含协议配置的`Info.plist`。

在`build/config.yml`中声明的协议会添加到：

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

#### 测试

```bash
# Open protocol URL from terminal
open "myapp://test/action"

# Check registered handlers
/System/Library/Frameworks/CoreServices.framework/Versions/A/Frameworks/LaunchServices.framework/Versions/A/Support/lsregister -dump | grep myapp
```

#### 通用链接

除自定义协议方案外，macOS 还支持<strong>通用链接</strong>，让常规 HTTPS 链接（例如`https://myawesomeapp.com/path`）也能启动应用。通用链接可在 Web 应用与桌面应用之间提供无缝的用户体验。

@note{type="caution"}
通用链接要求使用有效的 Apple Developer 证书和配置描述文件对 macOS 应用进行<strong>代码签名</strong>。未签名或临时签名的构建无法打开通用链接。测试前，请确保应用已正确签名。

@end

要启用通用链接，请按照[在应用中支持通用链接的 Apple 指南](https://developer.apple.com/documentation/xcode/supporting-universal-links-in-your-app)操作。你需要：

1. 在`entitlements.plist`中<strong>添加权利配置</strong>：
  ```xml
  <key>com.apple.developer.associated-domains</key>
  <array>
    <string>applinks:myawesomeapp.com</string>
  </array>
  ```


2. **将 NSUserActivityTypes 添加到 Info.plist**：
  ```xml
  <key>NSUserActivityTypes</key>
  <array>
    <string>NSUserActivityTypeBrowsingWeb</string>
  </array>
  ```


3. <strong>在网站上配置`apple-app-site-association`：</strong>在`https://myawesomeapp.com/.well-known/apple-app-site-association`托管一个`apple-app-site-association`文件。

通用链接触发应用时，你会收到相同的`ApplicationLaunchedWithUrl`事件，因此其处理代码与自定义协议方案完全相同。

[Linux]
### 桌面条目

在 Linux 上，协议通过`.desktop`文件注册。

#### 自动配置

使用`wails3 build`构建时，Wails 会生成包含协议处理程序的桌面条目文件。

**已在 v3 中修复**：Linux 桌面模板现已正确包含协议处理配置。

生成的桌面文件包含：

```ini
[Desktop Entry]
Type=Application
Name=My Application
Exec=/usr/bin/myapp %u
MimeType=x-scheme-handler/myapp;
```

#### 手动注册

如有需要，请手动安装桌面文件：

```bash
# Copy desktop file
cp myapp.desktop ~/.local/share/applications/

# Update desktop database
update-desktop-database ~/.local/share/applications/

# Register protocol handler
xdg-mime default myapp.desktop x-scheme-handler/myapp
```

#### 测试

```bash
# Open protocol URL
xdg-open "myapp://test/action"

# Check registered handler
xdg-mime query default x-scheme-handler/myapp
```

@end

## 完整示例

下面是一个处理多种协议操作的完整示例：

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

## 前端集成

在前端处理导航事件：

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

## 安全注意事项

### 验证所有输入

始终验证并清理来自外部来源的 URL：

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

### 防止注入攻击

切勿直接将 URL 作为代码或 SQL 执行：

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

## 测试

### 手动测试

在开发期间测试协议处理程序：

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

### HTML 测试

创建一个测试用 HTML 页面：

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

## 故障排除

### 协议未注册

**Windows：**

- 检查注册表：`HKEY_CURRENT_USER\SOFTWARE\Classes\<scheme>`
- 使用 NSIS 安装程序重新安装
- 确认安装程序以适当的权限运行

**macOS：**

- 使用`wails3 build`重新构建应用
- 检查应用包中的`Info.plist`：`MyApp.app/Contents/Info.plist`
- 重置 Launch Services：`/System/Library/Frameworks/CoreServices.framework/Frameworks/LaunchServices.framework/Support/lsregister -kill`

**Linux：**

- 检查桌面文件：`~/.local/share/applications/myapp.desktop`
- 更新数据库：`update-desktop-database ~/.local/share/applications/`
- 验证处理程序：`xdg-mime query default x-scheme-handler/myapp`

### 应用无法启动

**检查日志：**

```go
app := application.New(application.Options{
    LogLevel: slog.LevelDebug, // requires `import "log/slog"`
    // ...
})
```

**常见问题：**

- 应用未安装在预期位置
- 注册信息中的可执行文件路径与实际位置不一致
- 权限问题

## 最佳实践

### ✅ 推荐做法

- **使用描述性方案名称**——使用`mycompany-myapp`，而不是`mca`
- **验证所有输入**——切勿信任来自外部来源的 URL
- **妥善处理错误**——记录无效 URL，不要让应用崩溃
- **向用户提供反馈**——显示触发了什么操作
- **在所有平台上测试**——协议处理方式因平台而异
- **记录 URL 结构**——为用户和集成人员提供帮助

### ❌ 不要这样做

- **不要使用常见的协议方案名称**——避免使用`http`、`file`、`app`等名称
- **不要将 URL 作为代码执行**——这会带来巨大的安全风险
- **不要暴露敏感操作**——执行破坏性操作前必须要求用户确认
- **不要假定协议在任何环境中都能正常工作**——应提供备用机制
- **不要忘记 URL 编码**——正确处理特殊字符

## 后续步骤

- [Windows 打包](/guides/build/windows/)——了解 NSIS 安装程序选项
- [文件关联](/guides/file-associations/)——使用你的应用打开文件
- [单实例](/guides/single-instance/)——防止应用运行多个实例

---

<strong>有疑问？</strong>请在[Discord](https://discord.gg/JDdSxwjhGf)中提问，或查看[示例](https://github.com/wailsapp/wails/tree/master/v3/examples)。
