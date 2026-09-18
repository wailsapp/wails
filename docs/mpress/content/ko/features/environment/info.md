---
title: "환경"
description: "Wails 애플리케이션에서 시스템 환경 정보에 접근합니다"
slug: "features/environment/info"
sourcePath: "features/environment/info.md"
---

Wails는 EnvironmentManager API를 통해 포괄적인 환경 정보를 제공합니다. 애플리케이션은 이를 사용해 시스템 속성과 테마 환경설정을 감지하고 운영 체제의 파일 관리자와 통합할 수 있습니다.

## 환경 관리자에 접근하기

애플리케이션 인스턴스의 `Env` 속성을 통해 환경 관리자에 접근합니다:

```go
app := application.New(application.Options{
    Name: "Environment Demo",
})

// Access the environment manager
env := app.Env
```

## 시스템 정보

### 환경 정보 가져오기

런타임 환경에 관한 포괄적인 정보를 가져옵니다:

```go
envInfo := app.Env.Info()

app.Logger.Info("Environment information",
    "os", envInfo.OS,           // "windows", "darwin", "linux"
    "arch", envInfo.Arch,       // "amd64", "arm64", etc.
    "debug", envInfo.Debug,     // Debug mode flag
)

// Operating system details
if envInfo.OSInfo != nil {
    app.Logger.Info("OS details",
        "name", envInfo.OSInfo.Name,
        "version", envInfo.OSInfo.Version,
    )
}

// Platform-specific information
for key, value := range envInfo.PlatformInfo {
    app.Logger.Info("Platform info", "key", key, "value", value)
}
```

### 환경 구조

환경 정보에는 몇 가지 중요한 필드가 포함됩니다:

```go
type EnvironmentInfo struct {
    OS           string                 // Operating system: "windows", "darwin", "linux"
    Arch         string                 // Architecture: "amd64", "arm64", "386", etc.
    Debug        bool                   // Whether running in debug mode
    OSInfo       *operatingsystem.OS    // Detailed OS information
    PlatformInfo map[string]any         // Platform-specific details
}
```

## 테마 감지

### 다크 모드 감지

시스템에서 다크 모드를 사용 중인지 감지합니다:

```go
if app.Env.IsDarkMode() {
    app.Logger.Info("System is in dark mode")
    // Apply dark theme to your application
    applyDarkTheme()
} else {
    app.Logger.Info("System is in light mode")
    // Apply light theme to your application
    applyLightTheme()
}
```

### 테마 변경 모니터링

애플리케이션을 동적으로 업데이트할 수 있도록 테마 변경을 수신합니다:

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Listen for theme changes
app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(event *application.ApplicationEvent) {
    if app.Env.IsDarkMode() {
        app.Logger.Info("Switched to dark mode")
        updateApplicationTheme("dark")
    } else {
        app.Logger.Info("Switched to light mode")
        updateApplicationTheme("light")
    }
})

func updateApplicationTheme(theme string) {
    // Update your application's theme
    // This could emit an event to the frontend
    app.Event.Emit("theme:changed", theme)
}
```

## 파일 관리자 통합

### 파일 관리자 열기

시스템 파일 관리자를 특정 위치에서 엽니다:

```go
// Open file manager at a directory
err := app.Env.OpenFileManager("/Users/username/Documents", false)
if err != nil {
    app.Logger.Error("Failed to open file manager", "error", err)
}

// Open file manager and select a specific file
err = app.Env.OpenFileManager("/Users/username/Documents/report.pdf", true)
if err != nil {
    app.Logger.Error("Failed to open file manager with selection", "error", err)
}
```

### 일반적인 사용 사례

애플리케이션에서 파일 관리자를 사용해 파일이나 폴더를 표시합니다:

```go
func showInFileManager(app *application.App, path string) {
    err := app.Env.OpenFileManager(path, true)
    if err != nil {
        // Fallback: try opening just the directory
        dir := filepath.Dir(path)
        err = app.Env.OpenFileManager(dir, false)
        if err != nil {
            app.Logger.Error("Could not open file manager", "path", path, "error", err)
            
            // Show error to user
            app.Dialog.Error().
                SetTitle("File Manager Error").
                SetMessage("Could not open file manager").
                Show()
        }
    }
}

// Usage examples
func setupFileMenu(app *application.App) {
    menu := app.Menu.New()
    fileMenu := menu.AddSubmenu("File")
    
    fileMenu.Add("Show Downloads Folder").OnClick(func(ctx *application.Context) {
        homeDir, _ := os.UserHomeDir()
        downloadsDir := filepath.Join(homeDir, "Downloads")
        showInFileManager(app, downloadsDir)
    })
    
    fileMenu.Add("Show Application Data").OnClick(func(ctx *application.Context) {
        configDir, _ := os.UserConfigDir()
        appDir := filepath.Join(configDir, "MyApp")
        showInFileManager(app, appDir)
    })
}
```

## 플랫폼별 동작

### 적응형 애플리케이션 동작

환경 정보를 사용해 애플리케이션의 동작을 조정합니다:

```go
func configureForPlatform(app *application.App) {
    envInfo := app.Env.Info()
    
    switch envInfo.OS {
    case "darwin":
        configureMacOS(app)
    case "windows":
        configureWindows(app)
    case "linux":
        configureLinux(app)
    }
    
    // Adapt to architecture
    if envInfo.Arch == "arm64" {
        app.Logger.Info("Running on ARM architecture")
        // Potentially optimize for ARM
    }
}

func configureMacOS(app *application.App) {
    app.Logger.Info("Configuring for macOS")
    
    // macOS-specific configuration
    menu := app.Menu.New()
    menu.AddRole(application.AppMenu) // Add standard macOS app menu
    
    // Handle dark mode
    if app.Env.IsDarkMode() {
        setMacOSDarkTheme()
    }
}

func configureWindows(app *application.App) {
    app.Logger.Info("Configuring for Windows")
    
    // Windows-specific configuration
    // Set up Windows-style menus, key bindings, etc.
}

func configureLinux(app *application.App) {
    app.Logger.Info("Configuring for Linux")
    
    // Linux-specific configuration
    // Adapt to different desktop environments
}
```

## 디버그 모드 처리

### 개발 환경과 프로덕션 환경

디버그 모드 정보를 사용해 개발 기능을 활성화합니다:

```go
func setupApplicationMode(app *application.App) {
    envInfo := app.Env.Info()
    
    if envInfo.Debug {
        app.Logger.Info("Running in debug mode")
        setupDevelopmentFeatures(app)
    } else {
        app.Logger.Info("Running in production mode")
        setupProductionFeatures(app)
    }
}

func setupDevelopmentFeatures(app *application.App) {
    // Enable development-only features
    menu := app.Menu.New()
    
    // Add development menu
    devMenu := menu.AddSubmenu("Development")
    devMenu.Add("Reload Application").OnClick(func(ctx *application.Context) {
        // Reload the application
        window := app.Window.Current()
        if window != nil {
            window.Reload()
        }
    })
    
    devMenu.Add("Open DevTools").OnClick(func(ctx *application.Context) {
        window := app.Window.Current()
        if window != nil {
            window.OpenDevTools()
        }
    })
    
    devMenu.Add("Show Environment").OnClick(func(ctx *application.Context) {
        showEnvironmentDialog(app)
    })
}

func setupProductionFeatures(app *application.App) {
    // Production-only features
    // Disable debug logging, enable analytics, etc.
}
```

**macOS의 비공개 API:** `OpenDevTools()`에는 `private_mac_apis`이 필요하며, 그렇지 않으면 호출해도 아무 작업도 수행되지 않습니다. 프로덕션 빌드에는 `devtools`도 필요합니다. [Web Inspector 빌드 매트릭스](/guides/build/private-macos-apis/#web-inspector)를 참조하세요.

## 환경 정보 대화 상자

### 시스템 정보 표시

환경 정보를 표시하는 대화 상자를 만듭니다:

```go
func showEnvironmentDialog(app *application.App) {
    envInfo := app.Env.Info()
    
    details := fmt.Sprintf(`Environment Information:

Operating System: %s
Architecture: %s
Debug Mode: %t

Dark Mode: %t

Platform Information:`, 
        envInfo.OS, 
        envInfo.Arch, 
        envInfo.Debug,
        app.Env.IsDarkMode())
    
    // Add platform-specific details
    for key, value := range envInfo.PlatformInfo {
        details += fmt.Sprintf("\n%s: %v", key, value)
    }
    
    if envInfo.OSInfo != nil {
        details += fmt.Sprintf("\n\nOS Details:\nName: %s\nVersion: %s", 
            envInfo.OSInfo.Name, 
            envInfo.OSInfo.Version)
    }
    
    dialog := app.Dialog.Info()
    dialog.SetTitle("Environment Information")
    dialog.SetMessage(details)
    dialog.Show()
}
```

## 플랫폼 고려 사항

@tabs
[macOS]
macOS에서는 다음과 같이 동작합니다:

- 다크 모드 감지에 시스템 화면 모드 설정을 사용합니다
- 파일 관리자 작업에 Finder를 사용합니다
- 플랫폼 정보에 macOS 버전 세부 정보가 포함됩니다
- Apple Silicon Mac에서는 아키텍처가 "arm64"일 수 있습니다

```go
if envInfo.OS == "darwin" {
    // macOS-specific handling
    if envInfo.Arch == "arm64" {
        app.Logger.Info("Running on Apple Silicon")
    }
}
```

[Windows]
Windows에서는 다음과 같이 동작합니다:

- 다크 모드 감지에 Windows 테마 설정을 사용합니다
- 파일 관리자 작업에 Windows Explorer를 사용합니다
- 플랫폼 정보에 Windows 버전 세부 정보가 포함됩니다
- Windows 관련 추가 정보가 포함될 수 있습니다

```go
if envInfo.OS == "windows" {
    // Windows-specific handling
    for key, value := range envInfo.PlatformInfo {
        if key == "windows_version" {
            app.Logger.Info("Windows version", "version", value)
        }
    }
}
```

[Linux]
Linux에서는 다음과 같이 동작합니다:

- 다크 모드 감지 방식은 데스크톱 환경에 따라 다릅니다
- 파일 관리자 작업에 시스템 기본 파일 관리자를 사용합니다
- 플랫폼 정보에 배포판 세부 정보가 포함됩니다
- Linux 배포판에 따라 동작이 다를 수 있습니다

```go
if envInfo.OS == "linux" {
    // Linux-specific handling
    if distro, ok := envInfo.PlatformInfo["distribution"]; ok {
        app.Logger.Info("Linux distribution", "distro", distro)
    }
}
```

@end

## 모범 사례

1. **환경 정보 캐시하기**: 런타임 중에는 환경 정보가 거의 변경되지 않습니다:
  ```go
  type App struct {
      envInfo *application.EnvironmentInfo
  }

  func (a *App) getEnvInfo() application.EnvironmentInfo {
      if a.envInfo == nil {
          info := a.app.Env.Info()
          a.envInfo = &info
      }
      return *a.envInfo
  }
  ```


2. **테마 변경 처리하기**: 시스템 테마 변경을 수신합니다:
  ```go
  app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(event *application.ApplicationEvent) {
      updateTheme(app.Env.IsDarkMode())
  })
  ```


3. **파일 관리자 오류를 안정적으로 처리하기**: 파일 관리자 오류를 항상 처리합니다:
  ```go
  func openFileManagerSafely(app *application.App, path string) {
      err := app.Env.OpenFileManager(path, false)
      if err != nil {
          // Provide fallback or user notification
          app.Logger.Warn("Could not open file manager", "path", path)
      }
  }
  ```


4. **플랫폼별 기능**: 환경 정보를 사용해 플랫폼 기능을 활성화합니다:
  ```go
  envInfo := app.Env.Info()
  if envInfo.OS == "darwin" {
      // Enable macOS-specific features
  }
  ```


## 전체 예제

다음은 환경 관리를 보여 주는 전체 예제입니다:

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "runtime"
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

func main() {
    app := application.New(application.Options{
        Name: "Environment Demo",
    })

    // Setup application based on environment
    setupForEnvironment(app)

    // Monitor theme changes
    monitorThemeChanges(app)

    // Create menu with environment features
    setupEnvironmentMenu(app)

    // Create main window
    window := app.Window.New()
    window.SetTitle("Environment Demo")

    err := app.Run()
    if err != nil {
        panic(err)
    }
}

func setupForEnvironment(app *application.App) {
    envInfo := app.Env.Info()
    
    app.Logger.Info("Application environment",
        "os", envInfo.OS,
        "arch", envInfo.Arch,
        "debug", envInfo.Debug,
        "darkMode", app.Env.IsDarkMode(),
    )

    // Configure for platform
    switch envInfo.OS {
    case "darwin":
        app.Logger.Info("Configuring for macOS")
        // macOS-specific setup
    case "windows":
        app.Logger.Info("Configuring for Windows")
        // Windows-specific setup
    case "linux":
        app.Logger.Info("Configuring for Linux")
        // Linux-specific setup
    }

    // Apply initial theme
    if app.Env.IsDarkMode() {
        applyDarkTheme(app)
    } else {
        applyLightTheme(app)
    }
}

func monitorThemeChanges(app *application.App) {
    app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(event *application.ApplicationEvent) {
        if app.Env.IsDarkMode() {
            app.Logger.Info("System switched to dark mode")
            applyDarkTheme(app)
        } else {
            app.Logger.Info("System switched to light mode")
            applyLightTheme(app)
        }
    })
}

func setupEnvironmentMenu(app *application.App) {
    menu := app.Menu.New()
    
    // Add platform-specific app menu
    if runtime.GOOS == "darwin" {
        menu.AddRole(application.AppMenu)
    }
    
    // Tools menu
    toolsMenu := menu.AddSubmenu("Tools")
    
    toolsMenu.Add("Show Environment Info").OnClick(func(ctx *application.Context) {
        showEnvironmentInfo(app)
    })
    
    toolsMenu.Add("Open Downloads Folder").OnClick(func(ctx *application.Context) {
        openDownloadsFolder(app)
    })
    
    toolsMenu.Add("Toggle Theme").OnClick(func(ctx *application.Context) {
        // This would typically be handled by the system
        // but shown here for demonstration
        toggleTheme(app)
    })
    
    app.Menu.Set(menu)
}

func showEnvironmentInfo(app *application.App) {
    envInfo := app.Env.Info()
    
    message := fmt.Sprintf(`Environment Information:

Operating System: %s
Architecture: %s
Debug Mode: %t
Dark Mode: %t

Platform Details:`,
        envInfo.OS,
        envInfo.Arch,
        envInfo.Debug,
        app.Env.IsDarkMode())
    
    for key, value := range envInfo.PlatformInfo {
        message += fmt.Sprintf("\n%s: %v", key, value)
    }
    
    if envInfo.OSInfo != nil {
        message += fmt.Sprintf("\n\nOS Information:\nName: %s\nVersion: %s",
            envInfo.OSInfo.Name,
            envInfo.OSInfo.Version)
    }
    
    dialog := app.Dialog.Info()
    dialog.SetTitle("Environment Information")
    dialog.SetMessage(message)
    dialog.Show()
}

func openDownloadsFolder(app *application.App) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        app.Logger.Error("Could not get home directory", "error", err)
        return
    }
    
    downloadsDir := filepath.Join(homeDir, "Downloads")
    err = app.Env.OpenFileManager(downloadsDir, false)
    if err != nil {
        app.Logger.Error("Could not open Downloads folder", "error", err)
        
        app.Dialog.Error().
            SetTitle("File Manager Error").
            SetMessage("Could not open Downloads folder").
            Show()
    }
}

func applyDarkTheme(app *application.App) {
    app.Logger.Info("Applying dark theme")
    // Emit theme change to frontend
    app.Event.Emit("theme:apply", "dark")
}

func applyLightTheme(app *application.App) {
    app.Logger.Info("Applying light theme")
    // Emit theme change to frontend
    app.Event.Emit("theme:apply", "light")
}

func toggleTheme(app *application.App) {
    // This is just for demonstration
    // Real theme changes should come from the system
    currentlyDark := app.Env.IsDarkMode()
    if currentlyDark {
        applyLightTheme(app)
    } else {
        applyDarkTheme(app)
    }
}
```

@note{type="tip" title="유용한 팁"}
환경 정보를 사용해 플랫폼에 적합한 사용자 경험을 제공하세요. 예를 들어 macOS에서는 Command 키 단축키를 사용하고 Windows/Linux에서는 Control 키 단축키를 사용하세요.

@end

@note{type="danger" title="경고"}
환경 정보는 일반적으로 애플리케이션 런타임 중에 변경되지 않지만 테마 환경설정은 변경될 수 있습니다. UI를 동기화된 상태로 유지하려면 항상 테마 변경 이벤트를 수신하세요.

@end
