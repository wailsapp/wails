---
title: "瀏覽器整合"
description: "在使用者的預設網頁瀏覽器中開啟 URL 和檔案"
slug: "features/browser/integration"
sourcePath: "features/browser/integration.md"
---

Wails 透過 BrowserManager API 提供簡易的瀏覽器整合，讓應用程式能在使用者的預設網頁瀏覽器中開啟 URL 和檔案。這適合用於開啟外部連結、說明文件，或應由瀏覽器處理的檔案。

## 存取瀏覽器管理員

透過應用程式執行個體的`Browser`屬性存取瀏覽器管理員：

```go
app := application.New(application.Options{
    Name: "Browser Integration Demo",
})

// Access the browser manager
browser := app.Browser
```

## 開啟 URL

### 開啟網頁 URL

在使用者的預設網頁瀏覽器中開啟 URL：

```go
// Open a website
err := app.Browser.OpenURL("https://wails.io")
if err != nil {
    app.Logger.Error("Failed to open URL", "error", err)
}

// Open specific pages
err = app.Browser.OpenURL("https://github.com/wailsapp/wails")
if err != nil {
    app.Logger.Error("Failed to open GitHub", "error", err)
}
```

### 開啟本機 URL

開啟本機開發伺服器或區域網路資源：

```go
// Open local development server
err := app.Browser.OpenURL("http://localhost:3000")
if err != nil {
    app.Logger.Error("Failed to open local server", "error", err)
}

// Open network resource
err = app.Browser.OpenURL("http://192.168.1.100:8080")
if err != nil {
    app.Logger.Error("Failed to open network resource", "error", err)
}
```

## 開啟檔案

### 開啟 HTML 檔案

在瀏覽器中開啟本機 HTML 檔案：

```go
// Open an HTML file
err := app.Browser.OpenFile("/path/to/documentation.html")
if err != nil {
    app.Logger.Error("Failed to open HTML file", "error", err)
}

// Open generated reports
reportPath := "/tmp/report.html"
err = app.Browser.OpenFile(reportPath)
if err != nil {
    app.Logger.Error("Failed to open report", "error", err)
}
```

### 開啟其他檔案類型

開啟瀏覽器可處理的各種檔案類型：

```go
// Open PDF files
err := app.Browser.OpenFile("/path/to/document.pdf")
if err != nil {
    app.Logger.Error("Failed to open PDF", "error", err)
}

// Open image files
err = app.Browser.OpenFile("/path/to/image.png")
if err != nil {
    app.Logger.Error("Failed to open image", "error", err)
}

// Open text files
err = app.Browser.OpenFile("/path/to/readme.txt")
if err != nil {
    app.Logger.Error("Failed to open text file", "error", err)
}
```

## 常見使用情境

### 說明與文件

讓使用者可輕鬆存取說明資源：

```go
// Create help menu
func setupHelpMenu(app *application.App) {
    menu := app.Menu.New()
    helpMenu := menu.AddSubmenu("Help")
    
    helpMenu.Add("Online Documentation").OnClick(func(ctx *application.Context) {
        err := app.Browser.OpenURL("https://docs.yourapp.com")
        if err != nil {
            app.Dialog.Error().
                SetTitle("Error").
                SetMessage("Could not open documentation").
                Show()
        }
    })
    
    helpMenu.Add("GitHub Repository").OnClick(func(ctx *application.Context) {
        app.Browser.OpenURL("https://github.com/youruser/yourapp")
    })
    
    helpMenu.Add("Report Issue").OnClick(func(ctx *application.Context) {
        app.Browser.OpenURL("https://github.com/youruser/yourapp/issues/new")
    })
}
```

### 內容中的外部連結

處理應用程式內容中的外部連結：

```go
func handleExternalLink(app *application.App, url string) {
    // Validate the URL before opening
    if !isValidURL(url) {
        app.Logger.Warn("Invalid URL", "url", url)
        return
    }
    
    // Optionally confirm with user
    dialog := app.Dialog.Question()
    dialog.SetTitle("Open External Link")
    dialog.SetMessage(fmt.Sprintf("Open %s in your browser?", url))
    
    dialog.AddButton("Open").OnClick(func() {
        err := app.Browser.OpenURL(url)
        if err != nil {
            app.Logger.Error("Failed to open URL", "url", url, "error", err)
        }
    })
    
    dialog.AddButton("Cancel")
    dialog.Show()
}

func isValidURL(url string) bool {
    parsed, err := url.Parse(url)
    return err == nil && (parsed.Scheme == "http" || parsed.Scheme == "https")
}
```

### 匯出及檢視報告

產生報告並在瀏覽器中開啟：

```go
import (
    "html/template"
    "os"
    "path/filepath"
)

func generateAndOpenReport(app *application.App, data interface{}) error {
    // Create temporary file for the report
    tmpDir := os.TempDir()
    reportPath := filepath.Join(tmpDir, "report.html")
    
    // Generate HTML report
    tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Application Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .header { border-bottom: 2px solid #333; padding-bottom: 10px; }
        .data { margin-top: 20px; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Application Report</h1>
        <p>Generated: {{.Timestamp}}</p>
    </div>
    <div class="data">
        <!-- Report content here -->
        {{range .Items}}
        <p>{{.}}</p>
        {{end}}
    </div>
</body>
</html>`
    
    // Write report to file
    file, err := os.Create(reportPath)
    if err != nil {
        return err
    }
    defer file.Close()
    
    t, err := template.New("report").Parse(tmpl)
    if err != nil {
        return err
    }
    
    err = t.Execute(file, data)
    if err != nil {
        return err
    }
    
    // Open in browser
    return app.Browser.OpenFile(reportPath)
}
```

### 開發工具

在開發期間開啟開發資源：

```go
func setupDevelopmentMenu(app *application.App) {
    if !app.Env.Info().Debug {
        return // Only show in debug mode
    }
    
    menu := app.Menu.New()
    devMenu := menu.AddSubmenu("Development")
    
    devMenu.Add("Open DevTools").OnClick(func(ctx *application.Context) {
        // This would open browser devtools if available
        window := app.Window.Current()
        if window != nil {
            window.OpenDevTools()
        }
    })
    
    devMenu.Add("View Source").OnClick(func(ctx *application.Context) {
        // Open source code repository
        app.Browser.OpenURL("https://github.com/youruser/yourapp")
    })
    
    devMenu.Add("API Documentation").OnClick(func(ctx *application.Context) {
        // Open local API docs
        app.Browser.OpenURL("http://localhost:8080/docs")
    })
}
```

**macOS 上的私有 API：**`OpenDevTools()`需要`private_mac_apis`；否則該呼叫不會執行任何操作。正式環境組建還需要`devtools`。請參閱[網頁檢查器組建矩陣](/guides/build/private-macos-apis/#web-inspector)。

## 錯誤處理

### 妥善處理錯誤

開啟 URL 或檔案時，務必處理可能發生的錯誤：

```go
func openURLWithFallback(app *application.App, url string, fallbackMessage string) {
    err := app.Browser.OpenURL(url)
    if err != nil {
        app.Logger.Error("Failed to open URL", "url", url, "error", err)
        
        // Show fallback dialog with URL
        dialog := app.Dialog.Info()
        dialog.SetTitle("Unable to Open Link")
        dialog.SetMessage(fmt.Sprintf("%s\n\nURL: %s", fallbackMessage, url))
        dialog.Show()
    }
}

// Usage
openURLWithFallback(app, 
    "https://docs.example.com", 
    "Please open the following URL manually in your browser:")
```

### 使用者回饋

操作成功或失敗時提供回饋：

```go
func openURLWithFeedback(app *application.App, url string) {
    err := app.Browser.OpenURL(url)
    if err != nil {
        // Show error dialog
        app.Dialog.Error().
            SetTitle("Browser Error").
            SetMessage(fmt.Sprintf("Could not open URL: %s", err.Error())).
            Show()
    } else {
        // Optionally show success notification
        app.Logger.Info("URL opened successfully", "url", url)
    }
}
```

## 平台注意事項

@tabs
[macOS]
在 macOS 上：

- 使用`open`命令啟動預設瀏覽器
- 遵循使用者在「系統偏好設定」中設定的預設瀏覽器
- 如果應用程式在沙箱中執行，可能會提示要求權限
- 能正確處理本機檔案的`file://` URL

[Windows]
在 Windows 上：

- 使用 Windows Shell API 開啟 URL
- 遵循 Windows 設定中的預設瀏覽器設定
- 能正確處理 Windows 路徑格式
- 開啟不受信任的 URL 時可能會顯示安全性警告

[Linux]
在 Linux 上：

- 先嘗試使用`xdg-open`，失敗時改用其他方法
- 行為因桌面環境而異
- 如果已設定`BROWSER`環境變數，則會遵循其設定
- 在最小化安裝的環境中，可能需要額外的套件

@end

## 最佳實務

1. **務必處理錯誤**：瀏覽器操作可能因各種原因而失敗：
  ```go
  if err := app.Browser.OpenURL(url); err != nil {
      app.Logger.Error("Failed to open browser", "error", err)
      // Provide fallback or user notification
  }
  ```


2. **驗證 URL**：開啟 URL 前，請確認其格式正確：
  ```go
  func isValidHTTPURL(str string) bool {
      u, err := url.Parse(str)
      return err == nil && (u.Scheme == "http" || u.Scheme == "https")
  }
  ```


3. **使用者確認**：對於外部連結，請考慮先徵求使用者同意：
  ```go
  // Show confirmation dialog before opening external links
  confirmAndOpen(app, "https://external-site.com")
  ```


4. **確保檔案路徑安全**：開啟檔案時，請確認路徑安全：
  ```go
  func openSafeFile(app *application.App, filename string) error {
      // Ensure file exists and is readable
      if _, err := os.Stat(filename); err != nil {
          return err
      }
      return app.Browser.OpenFile(filename)
  }
  ```


## 完整範例

以下完整範例示範各種瀏覽器整合模式：

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "Browser Integration Demo",
    })

    // Setup menu with browser actions
    setupMenu(app)

    // Create main window
    window := app.Window.New()
    window.SetTitle("Browser Integration")

    err := app.Run()
    if err != nil {
        panic(err)
    }
}

func setupMenu(app *application.App) {
    menu := app.Menu.New()
    
    // File menu
    fileMenu := menu.AddSubmenu("File")
    fileMenu.Add("Generate Report").OnClick(func(ctx *application.Context) {
        generateHTMLReport(app)
    })
    
    // Help menu
    helpMenu := menu.AddSubmenu("Help")
    helpMenu.Add("Documentation").OnClick(func(ctx *application.Context) {
        openWithConfirmation(app, "https://docs.example.com")
    })
    helpMenu.Add("Support").OnClick(func(ctx *application.Context) {
        openWithConfirmation(app, "https://support.example.com")
    })
    
    app.Menu.Set(menu)
}

func openWithConfirmation(app *application.App, url string) {
    dialog := app.Dialog.Question()
    dialog.SetTitle("Open External Link")
    dialog.SetMessage(fmt.Sprintf("Open %s in your browser?", url))
    
    dialog.AddButton("Open").OnClick(func() {
        if err := app.Browser.OpenURL(url); err != nil {
            showError(app, "Failed to open URL", err)
        }
    })
    
    dialog.AddButton("Cancel")
    dialog.Show()
}

func generateHTMLReport(app *application.App) {
    // Create temporary HTML file
    tmpDir := os.TempDir()
    reportPath := filepath.Join(tmpDir, "demo_report.html")
    
    html := `
<!DOCTYPE html>
<html>
<head>
    <title>Demo Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .header { color: #333; border-bottom: 1px solid #ccc; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Application Report</h1>
        <p>This is a sample report generated by the application.</p>
    </div>
    <div class="content">
        <h2>Report Details</h2>
        <p>This report was generated to demonstrate browser integration.</p>
    </div>
</body>
</html>`
    
    err := os.WriteFile(reportPath, []byte(html), 0644)
    if err != nil {
        showError(app, "Failed to create report", err)
        return
    }
    
    // Open in browser
    err = app.Browser.OpenFile(reportPath)
    if err != nil {
        showError(app, "Failed to open report", err)
    }
}

func showError(app *application.App, message string, err error) {
    app.Dialog.Error().
        SetTitle("Error").
        SetMessage(fmt.Sprintf("%s: %v", message, err)).
        Show()
}
```

@note{type="tip" title="專業提示"}
請考慮為瀏覽器操作失敗的情況提供備援機制，例如將 URL 複製到剪貼簿，或在對話方塊中顯示 URL，讓使用者手動開啟。

@end

@note{type="danger" title="警告"}
開啟 URL 和檔案路徑前，務必先加以驗證，以免發生安全性問題。對於使用者提供的 URL，請勿未經驗證便直接開啟。

@end
