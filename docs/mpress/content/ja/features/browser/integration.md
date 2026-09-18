---
title: "ブラウザー統合"
description: "URL とファイルをユーザーの既定の Web ブラウザーで開く"
slug: "features/browser/integration"
sourcePath: "features/browser/integration.md"
---

Wails は BrowserManager API を通じてシンプルなブラウザー統合を提供し、アプリケーションから URL やファイルをユーザーの既定の Web ブラウザーで開けるようにします。これは、外部リンクやドキュメント、ブラウザーで処理すべきファイルを開く場合に便利です。

## ブラウザーマネージャーへのアクセス

アプリケーションインスタンスの `Browser` プロパティを介してブラウザーマネージャーにアクセスします。

```go
app := application.New(application.Options{
    Name: "Browser Integration Demo",
})

// Access the browser manager
browser := app.Browser
```

## URL を開く

### Web URL を開く

URL をユーザーの既定の Web ブラウザーで開きます。

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

### ローカル URL を開く

ローカルの開発サーバーまたはローカルネットワーク上のリソースを開きます。

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

## ファイルを開く

### HTML ファイルを開く

ローカルの HTML ファイルをブラウザーで開きます。

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

### その他のファイル形式を開く

ブラウザーで処理できるさまざまな形式のファイルを開きます。

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

## 一般的なユースケース

### ヘルプとドキュメント

ヘルプリソースへ簡単にアクセスできるようにします。

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

### コンテンツ内の外部リンク

アプリケーションのコンテンツ内にある外部リンクを処理します。

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

### レポートのエクスポートと表示

レポートを生成し、ブラウザーで開きます。

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

### 開発ツール

開発中に開発リソースを開きます。

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

**macOS での非公開 API：** `OpenDevTools()` には `private_mac_apis` が必要です。指定されていない場合、この呼び出しでは何も行われません。本番ビルドでは `devtools` も必要です。[Web Inspector のビルドマトリックス](/guides/build/private-macos-apis/#web-inspector)を参照してください。

## エラー処理

### 適切なエラー処理

URL またはファイルを開くときは、発生する可能性があるエラーを必ず処理してください。

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

### ユーザーへのフィードバック

操作の成功または失敗をユーザーに通知します。

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

## プラットフォーム固有の考慮事項

@tabs
[macOS]
macOS では、次のように動作します。

- `open` コマンドを使用して既定のブラウザーを起動します
- 「システム環境設定」で指定されたユーザーの既定のブラウザー設定に従います
- アプリケーションがサンドボックス化されている場合、権限を求めることがあります
- ローカルファイルの `file://` URL を正しく処理します

[Windows]
Windows では、次のように動作します。

- Windows Shell API を使用して URL を開きます
- 「Windows の設定」で指定された既定のブラウザー設定に従います
- Windows のパス形式を正しく処理します
- 信頼されていない URL に対してセキュリティ警告が表示されることがあります

[Linux]
Linux では、次のように動作します。

- 最初に `xdg-open` の使用を試み、失敗した場合はほかの方法を使用します
- 動作はデスクトップ環境によって異なります
- `BROWSER` 環境変数が設定されている場合は、その設定に従います
- 最小構成の環境では追加のパッケージが必要になる場合があります

@end

## ベストプラクティス

1. **エラーを必ず処理する**：ブラウザー操作はさまざまな理由で失敗する可能性があります。
  ```go
  if err := app.Browser.OpenURL(url); err != nil {
      app.Logger.Error("Failed to open browser", "error", err)
      // Provide fallback or user notification
  }
  ```


2. **URL を検証する**：URL を開く前に、形式が正しいことを確認してください。
  ```go
  func isValidHTTPURL(str string) bool {
      u, err := url.Parse(str)
      return err == nil && (u.Scheme == "http" || u.Scheme == "https")
  }
  ```


3. **ユーザーの確認を得る**：外部リンクを開く場合は、ユーザーに許可を求めることを検討してください。
  ```go
  // Show confirmation dialog before opening external links
  confirmAndOpen(app, "https://external-site.com")
  ```


4. **ファイルパスを保護する**：ファイルを開くときは、パスが安全であることを確認してください。
  ```go
  func openSafeFile(app *application.App, filename string) error {
      // Ensure file exists and is readable
      if _, err := os.Stat(filename); err != nil {
          return err
      }
      return app.Browser.OpenFile(filename)
  }
  ```


## 完全な例

さまざまなブラウザー統合パターンを示す完全な例を次に示します。

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

@note{type="tip" title="ヒント"}
ブラウザー操作が失敗した場合に備えて、URL をクリップボードにコピーしたり、手動で開けるようにダイアログに表示したりするフォールバック手段を用意することを検討してください。

@end

@note{type="danger" title="警告"}
セキュリティ上の問題を防ぐため、URL とファイルパスは開く前に必ず検証してください。ユーザーが指定した URL を検証せずに開く場合は注意してください。

@end
