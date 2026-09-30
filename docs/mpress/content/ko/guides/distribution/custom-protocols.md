---
title: "사용자 지정 URL 프로토콜"
description: "링크에서 애플리케이션을 실행하도록 사용자 지정 URL 스킴 등록"
slug: "guides/distribution/custom-protocols"
sourcePath: "guides/distribution/custom-protocols.md"
---

사용자 지정 URL 프로토콜(URL 스킴이라고도 함)을 사용하면 사용자가 `myapp://action` 또는 `myapp://open/document` 같은 사용자 지정 프로토콜이 포함된 링크를 클릭할 때 애플리케이션을 실행할 수 있습니다.

## 개요

사용자 지정 프로토콜로 다음 기능을 사용할 수 있습니다:

- **딥 링크**: 특정 데이터와 함께 앱 실행
- **브라우저 통합**: 웹 페이지의 링크 처리
- **이메일 링크**: 이메일 클라이언트에서 앱 열기
- **앱 간 통신**: 다른 애플리케이션에서 실행

**예시**: `myapp://open/document?id=123`은 앱을 실행하고 문서 123을 엽니다.

## 구성

애플리케이션 옵션에서 사용자 지정 프로토콜을 정의합니다:

사용자 지정 프로토콜은 `build/config.yml`에 선언합니다. 플랫폼 패키저(Windows의 NSIS 매크로, MSIX 매니페스트, macOS의 `CFBundleURLTypes`, Linux의 `.desktop`/`xdg-mime`)가 패키징할 때 이 선언을 사용합니다. `application.Protocol` 타입은 없으며 `application.Options`에는 `Protocols` 필드가 없습니다.

```yaml
# build/config.yml
protocols:
  - scheme: myapp
    description: "My Application Protocol"
```

Go 코드에서는 `ApplicationLaunchedWithUrl` 이벤트를 통해 URL을 사용한 실행을 수신합니다:

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

## 프로토콜 핸들러

수신 URL을 처리하려면 프로토콜 이벤트를 수신합니다:

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

## URL 구조

명확한 계층형 URL 구조를 설계합니다:

```
myapp://action/resource?param=value

Examples:
myapp://open/document?id=123
myapp://settings/theme?mode=dark
myapp://user/profile?username=john
```

**권장 사항:**

- 스킴 이름은 소문자로 작성합니다
- 스킴은 짧고 기억하기 쉽게 만듭니다
- 리소스에는 계층형 경로를 사용합니다
- 선택적 데이터에는 쿼리 매개변수를 사용합니다
- 특수 문자를 URL 인코딩합니다

## 플랫폼 등록

사용자 지정 프로토콜의 등록 방식은 플랫폼마다 다릅니다.

@tabs{sync-key="platform"}
[Windows]
### Windows NSIS 설치 프로그램

NSIS 설치 프로그램을 사용하면 **Wails v3가 사용자 지정 프로토콜을 자동으로 등록합니다**.

#### 자동 등록

`wails3 build`으로 애플리케이션을 빌드하면 NSIS 설치 프로그램이 다음 작업을 수행합니다:

1. `build/config.yml`의 `protocols:` 키 아래에 선언된 모든 프로토콜을 자동으로 등록합니다
2. 프로토콜을 애플리케이션 실행 파일에 연결합니다
3. 올바른 레지스트리 항목을 설정합니다
4. 제거할 때 프로토콜 연결을 삭제합니다

**추가 구성이 필요하지 않습니다!**

#### 작동 방식

NSIS 템플릿에는 다음과 같은 기본 제공 매크로가 포함되어 있습니다:

- `wails.associateCustomProtocols` - 설치 중 프로토콜을 등록합니다
- `wails.unassociateCustomProtocols` - 제거 중 프로토콜을 삭제합니다

이러한 매크로는 `Protocols` 구성에 따라 자동으로 호출됩니다.

#### 수동 레지스트리 등록(고급)

NSIS 외부에서 수동으로 등록해야 하는 경우:

```batch
@echo off
REM Register custom protocol
REG ADD "HKEY_CURRENT_USER\SOFTWARE\Classes\myapp" /ve /d "URL:My Application Protocol" /f
REG ADD "HKEY_CURRENT_USER\SOFTWARE\Classes\myapp" /v "URL Protocol" /t REG_SZ /d "" /f
REG ADD "HKEY_CURRENT_USER\SOFTWARE\Classes\myapp\shell\open\command" /ve /d "\"%1\"" /f
```

#### 테스트

프로토콜 등록을 테스트합니다:

```powershell
# Open protocol URL from PowerShell
Start-Process "myapp://test/action"

# Or from command prompt
start myapp://test/action
```

### Windows MSIX 패키지

MSIX로 패키징할 때도 사용자 지정 프로토콜이 자동으로 등록됩니다.

#### 자동 등록

MSIX로 애플리케이션을 빌드하면 매니페스트에 `build/config.yml` 프로토콜 구성의 프로토콜 등록 정보가 자동으로 포함됩니다.

생성된 매니페스트에는 다음이 포함됩니다:

```xml
<uap:Extension Category="windows.protocol">
  <uap:Protocol Name="myapp">
    <uap:DisplayName>My Application Protocol</uap:DisplayName>
  </uap:Protocol>
</uap:Extension>
```

#### 유니버설 링크(웹-앱 연결)

Windows는 macOS의 유니버설 링크와 유사하게 작동하는 <strong>웹-앱 연결</strong>을 지원합니다. 애플리케이션을 MSIX 패키지로 배포할 때 HTTPS 링크를 통해 앱이 직접 실행되도록 설정할 수 있습니다.

@note{type="note"}
웹-앱 연결을 사용하려면 매니페스트를 수동으로 구성해야 합니다. 사용자 지정 프로토콜 스킴은 `build/config.yml`에서 자동으로 구성되지만 연결된 도메인은 MSIX 매니페스트에 직접 추가해야 합니다.

@end

웹-앱 연결을 사용하려면 [Microsoft의 웹-앱 연결 가이드](https://learn.microsoft.com/en-us/windows/apps/develop/launch/web-to-app-linking)를 따르세요. 다음 작업이 필요합니다:

1. **MSIX 매니페스트에 App URI Handler를 수동으로 추가합니다**(`build/windows/msix/app_manifest.xml`):
  ```xml
  <uap3:Extension Category="windows.appUriHandler">
    <uap3:AppUriHandler>
      <uap3:Host Name="myawesomeapp.com"/>
    </uap3:AppUriHandler>
  </uap3:Extension>
  ```


2. **웹사이트에서 `windows-app-web-link`을 구성합니다:** `https://myawesomeapp.com/.well-known/windows-app-web-link`에 `windows-app-web-link` 파일을 호스팅합니다. 이 파일에는 앱의 패키지 정보와 앱이 처리하는 경로가 포함되어야 합니다.

웹-앱 링크가 애플리케이션을 실행하면 사용자 지정 프로토콜 스킴을 사용할 때와 동일한 `ApplicationLaunchedWithUrl` 이벤트를 받습니다.

[macOS]
### Info.plist 구성

macOS에서는 `Info.plist` 파일을 통해 프로토콜을 등록합니다.

#### 자동 구성

Wails는 `wails3 build`로 빌드할 때 프로토콜이 포함된 `Info.plist`을 자동으로 생성합니다.

`build/config.yml`에 선언된 프로토콜은 다음 항목에 추가됩니다.

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

#### 테스트

```bash
# Open protocol URL from terminal
open "myapp://test/action"

# Check registered handlers
/System/Library/Frameworks/CoreServices.framework/Versions/A/Frameworks/LaunchServices.framework/Versions/A/Support/lsregister -dump | grep myapp
```

#### Universal Links

macOS는 사용자 지정 프로토콜 스킴뿐만 아니라 일반 HTTPS 링크(예: `https://myawesomeapp.com/path`)로 앱을 실행할 수 있는 <strong>Universal Links</strong>도 지원합니다. Universal Links를 사용하면 웹 앱과 데스크톱 앱을 매끄럽게 오가는 사용자 경험을 제공할 수 있습니다.

@note{type="caution"}
Universal Links를 사용하려면 유효한 Apple Developer 인증서와 프로비저닝 프로파일로 macOS 앱에 <strong>코드 서명</strong>을 해야 합니다. 서명되지 않았거나 ad-hoc 방식으로 서명된 빌드에서는 Universal Links를 열 수 없습니다. 테스트하기 전에 앱이 올바르게 서명되었는지 확인하세요.

@end

Universal Links를 활성화하려면 [앱에서 Universal Links를 지원하는 방법에 관한 Apple 가이드](https://developer.apple.com/documentation/xcode/supporting-universal-links-in-your-app)를 따르세요. 다음 작업이 필요합니다.

1. `entitlements.plist`에 **엔타이틀먼트를 추가합니다**.
  ```xml
  <key>com.apple.developer.associated-domains</key>
  <array>
    <string>applinks:myawesomeapp.com</string>
  </array>
  ```


2. **Info.plist에 NSUserActivityTypes를 추가합니다**.
  ```xml
  <key>NSUserActivityTypes</key>
  <array>
    <string>NSUserActivityTypeBrowsingWeb</string>
  </array>
  ```


3. **웹사이트에서 `apple-app-site-association`을 구성합니다.** `https://myawesomeapp.com/.well-known/apple-app-site-association`에 `apple-app-site-association` 파일을 호스팅하세요.

Universal Link가 앱을 실행하면 사용자 지정 프로토콜 스킴과 동일한 `ApplicationLaunchedWithUrl` 이벤트가 전달되므로 처리 코드를 똑같이 사용할 수 있습니다.

[Linux]
### 데스크톱 항목

Linux에서는 `.desktop` 파일을 통해 프로토콜을 등록합니다.

#### 자동 구성

`wails3 build`로 빌드하면 Wails가 프로토콜 핸들러가 포함된 데스크톱 항목 파일을 생성합니다.

**v3에서 수정됨**: 이제 Linux 데스크톱 템플릿에 프로토콜 처리 기능이 올바르게 포함됩니다.

생성된 데스크톱 파일에는 다음 항목이 포함됩니다.

```ini
[Desktop Entry]
Type=Application
Name=My Application
Exec=/usr/bin/myapp %u
MimeType=x-scheme-handler/myapp;
```

#### 수동 등록

필요한 경우 데스크톱 파일을 수동으로 설치하세요.

```bash
# Copy desktop file
cp myapp.desktop ~/.local/share/applications/

# Update desktop database
update-desktop-database ~/.local/share/applications/

# Register protocol handler
xdg-mime default myapp.desktop x-scheme-handler/myapp
```

#### 테스트

```bash
# Open protocol URL
xdg-open "myapp://test/action"

# Check registered handler
xdg-mime query default x-scheme-handler/myapp
```

@end

## 전체 예제

다음은 여러 프로토콜 작업을 처리하는 전체 예제입니다.

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

## 프런트엔드 통합

프런트엔드에서 탐색 이벤트를 처리하세요.

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

## 보안 고려 사항

### 모든 입력 검증

외부 소스에서 받은 URL은 항상 검증하고 안전하게 정제하세요.

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

### 인젝션 공격 방지

URL을 코드나 SQL로 직접 실행해서는 안 됩니다.

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

## 테스트

### 수동 테스트

개발 중에 프로토콜 핸들러를 테스트하세요.

**Windows:**

```powershell
Start-Process "myapp://test/action?id=123"
```

**macOS:**

```bash
open "myapp://test/action?id=123"
```

**Linux:**

```bash
xdg-open "myapp://test/action?id=123"
```

### HTML 테스트

테스트용 HTML 페이지를 만드세요.

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

## 문제 해결

### 프로토콜이 등록되지 않음

**Windows:**

- 레지스트리를 확인하세요: `HKEY_CURRENT_USER\SOFTWARE\Classes\<scheme>`
- NSIS 설치 프로그램으로 다시 설치하세요.
- 설치 프로그램이 적절한 권한으로 실행되었는지 확인하세요.

**macOS:**

- `wails3 build`로 애플리케이션을 다시 빌드하세요.
- 앱 번들에서 `Info.plist`을 확인하세요: `MyApp.app/Contents/Info.plist`
- Launch Services를 재설정하세요: `/System/Library/Frameworks/CoreServices.framework/Frameworks/LaunchServices.framework/Support/lsregister -kill`

**Linux:**

- 데스크톱 파일을 확인하세요: `~/.local/share/applications/myapp.desktop`
- 데이터베이스를 업데이트하세요: `update-desktop-database ~/.local/share/applications/`
- 핸들러를 확인하세요: `xdg-mime query default x-scheme-handler/myapp`

### 애플리케이션이 실행되지 않음

**로그 확인:**

```go
app := application.New(application.Options{
    LogLevel: slog.LevelDebug, // requires `import "log/slog"`
    // ...
})
```

**일반적인 문제:**

- 애플리케이션이 예상 위치에 설치되지 않음
- 등록된 실행 파일 경로가 실제 위치와 일치하지 않음
- 권한 문제

## 권장 사례

### ✅ 권장 사항

- **설명적인 스킴 이름을 사용하세요**. `mca` 대신 `mycompany-myapp`을 사용하세요.
- **모든 입력을 검증하세요**. 외부 소스에서 받은 URL을 절대 신뢰하지 마세요.
- **오류를 적절히 처리하세요**. 잘못된 URL은 로그에 기록하고 비정상 종료하지 마세요.
- **사용자에게 피드백 제공** - 어떤 동작이 실행되었는지 표시하세요
- **모든 플랫폼에서 테스트** - 프로토콜 처리 방식은 플랫폼마다 다릅니다
- **URL 구조 문서화** - 사용자와 통합 개발자가 이해할 수 있도록 도우세요

### ❌ 하지 말아야 할 사항

- **흔한 스킴 이름을 사용하지 마세요** - `http`, `file`, `app` 등의 이름은 피하세요.
- **URL을 코드로 실행하지 마세요** - 매우 심각한 보안 위험이 있습니다
- **민감한 작업을 노출하지 마세요** - 파괴적인 작업에는 확인을 요구하세요
- **프로토콜이 어디서나 작동한다고 가정하지 마세요** - 대체 메커니즘을 마련하세요
- **URL 인코딩을 잊지 마세요** - 특수 문자를 올바르게 처리하세요

## 다음 단계

- [Windows 패키징](/guides/build/windows/) - NSIS 설치 프로그램 옵션을 알아보세요
- [파일 연결](/guides/file-associations/) - 앱으로 파일을 여는 방법을 알아보세요
- [단일 인스턴스](/guides/single-instance/) - 앱 인스턴스가 여러 개 실행되지 않도록 방지하세요

---

**궁금한 점이 있으신가요?** [Discord](https://discord.gg/JDdSxwjhGf)에서 질문하거나 [예제](https://github.com/wailsapp/wails/tree/master/v3/examples)를 확인하세요.
