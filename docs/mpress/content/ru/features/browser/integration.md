---
title: "Интеграция с браузером"
description: "Открытие URL-адресов и файлов в браузере пользователя по умолчанию"
slug: "features/browser/integration"
sourcePath: "features/browser/integration.md"
---

Wails обеспечивает простую интеграцию с браузером через API BrowserManager, позволяя приложению открывать URL-адреса и файлы в браузере пользователя по умолчанию. Это удобно для открытия внешних ссылок, документации и файлов, которые должен обрабатывать браузер.

## Доступ к диспетчеру браузера

Доступ к диспетчеру браузера осуществляется через свойство `Browser` экземпляра приложения:

```go
app := application.New(application.Options{
    Name: "Browser Integration Demo",
})

// Access the browser manager
browser := app.Browser
```

## Открытие URL-адресов

### Открытие веб-адресов

Открывайте URL-адреса в браузере пользователя по умолчанию:

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

### Открытие локальных URL-адресов

Открывайте локальные серверы разработки и ресурсы локальной сети:

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

## Открытие файлов

### Открытие HTML-файлов

Открывайте локальные HTML-файлы в браузере:

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

### Открытие файлов других типов

Открывайте файлы различных типов, которые поддерживаются браузерами:

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

## Распространённые сценарии использования

### Справка и документация

Предоставьте удобный доступ к справочным ресурсам:

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

### Внешние ссылки в содержимом

Обрабатывайте внешние ссылки из содержимого приложения:

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

### Экспорт и просмотр отчётов

Создавайте и открывайте отчёты в браузере:

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

### Инструменты разработки

Открывайте ресурсы для разработчиков во время разработки:

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

**Закрытый API в macOS:** для `OpenDevTools()` требуется `private_mac_apis`; в противном случае вызов не выполняет никаких действий. Для производственных сборок также требуется `devtools`. См. [матрицу сборок Web Inspector](/guides/build/private-macos-apis/#web-inspector).

## Обработка ошибок

### Корректная обработка ошибок

Всегда обрабатывайте возможные ошибки при открытии URL-адресов или файлов:

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

### Обратная связь с пользователем

Сообщайте пользователю об успешном или неудачном выполнении операций:

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

## Особенности платформ

@tabs
[macOS]
В macOS:

- Для запуска браузера по умолчанию используется команда `open`
- Учитывается браузер по умолчанию, выбранный пользователем в «Системных настройках»
- Если приложение работает в песочнице, может появиться запрос разрешения
- URL-адреса `file://` для локальных файлов обрабатываются корректно

[Windows]
В Windows:

- Для открытия URL-адресов используется Windows Shell API
- Учитывается браузер по умолчанию, выбранный в «Параметрах Windows»
- Форматы путей Windows обрабатываются корректно
- При открытии ненадёжных URL-адресов могут отображаться предупреждения системы безопасности

[Linux]
В Linux:

- Сначала предпринимается попытка использовать `xdg-open`, а при неудаче применяются другие методы
- Поведение зависит от среды рабочего стола
- Если задана переменная среды `BROWSER`, её значение учитывается
- В минимальных установках могут потребоваться дополнительные пакеты

@end

## Рекомендации

1. **Всегда обрабатывайте ошибки**: операции с браузером могут завершаться с ошибкой по разным причинам:
  ```go
  if err := app.Browser.OpenURL(url); err != nil {
      app.Logger.Error("Failed to open browser", "error", err)
      // Provide fallback or user notification
  }
  ```


2. **Проверяйте URL-адреса**: перед открытием убедитесь, что URL-адреса имеют корректный формат:
  ```go
  func isValidHTTPURL(str string) bool {
      u, err := url.Parse(str)
      return err == nil && (u.Scheme == "http" || u.Scheme == "https")
  }
  ```


3. **Подтверждение пользователя**: подумайте, стоит ли запрашивать разрешение пользователя при открытии внешних ссылок:
  ```go
  // Show confirmation dialog before opening external links
  confirmAndOpen(app, "https://external-site.com")
  ```


4. **Безопасные пути к файлам**: при открытии файлов убедитесь, что пути безопасны:
  ```go
  func openSafeFile(app *application.App, filename string) error {
      // Ensure file exists and is readable
      if _, err := os.Stat(filename); err != nil {
          return err
      }
      return app.Browser.OpenFile(filename)
  }
  ```


## Полный пример

Ниже приведён полный пример различных вариантов интеграции с браузером:

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

@note{type="tip" title="Совет"}
Рассмотрите возможность предусмотреть резервные механизмы на случай сбоя операций с браузером, например копирование URL-адресов в буфер обмена или их отображение в диалоговом окне для открытия вручную.

@end

@note{type="danger" title="Предупреждение"}
Всегда проверяйте URL-адреса и пути к файлам перед открытием, чтобы предотвратить проблемы с безопасностью. Будьте осторожны при открытии URL-адресов, предоставленных пользователем, без проверки.

@end
