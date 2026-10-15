---
title: "Browserintegration"
description: "URLs und Dateien im Standard-Webbrowser des Benutzers öffnen"
slug: "features/browser/integration"
sourcePath: "features/browser/integration.md"
---

Wails bietet über die BrowserManager-API eine einfache Browserintegration, mit der Ihre Anwendung URLs und Dateien im Standard-Webbrowser des Benutzers öffnen kann. Dies ist nützlich, um externe Links, Dokumentation oder Dateien zu öffnen, die vom Browser verarbeitet werden sollen.

## Auf den Browser-Manager zugreifen

Greifen Sie über die Eigenschaft `Browser` Ihrer Anwendungsinstanz auf den Browser-Manager zu:

```go
app := application.New(application.Options{
    Name: "Browser Integration Demo",
})

// Access the browser manager
browser := app.Browser
```

## URLs öffnen

### Web-URLs öffnen

Öffnen Sie URLs im Standard-Webbrowser des Benutzers:

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

### Lokale URLs öffnen

Öffnen Sie lokale Entwicklungsserver oder Ressourcen im lokalen Netzwerk:

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

## Dateien öffnen

### HTML-Dateien öffnen

Öffnen Sie lokale HTML-Dateien im Browser:

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

### Andere Dateitypen öffnen

Öffnen Sie verschiedene Dateitypen, die Browser verarbeiten können:

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

## Häufige Anwendungsfälle

### Hilfe und Dokumentation

Ermöglichen Sie einfachen Zugriff auf Hilferessourcen:

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

### Externe Links in Inhalten

Verarbeiten Sie externe Links aus den Inhalten Ihrer Anwendung:

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

### Berichte exportieren und anzeigen

Erstellen und öffnen Sie Berichte im Browser:

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

### Entwicklungswerkzeuge

Öffnen Sie während der Entwicklung Entwicklungsressourcen:

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

**Private API unter macOS:** `OpenDevTools()` erfordert `private_mac_apis`; andernfalls bleibt der Aufruf ohne Wirkung. Produktions-Builds erfordern außerdem `devtools`. Weitere Informationen finden Sie in der [Build-Matrix für den Web Inspector](/guides/build/private-macos-apis/#web-inspector).

## Fehlerbehandlung

### Robuste Fehlerbehandlung

Behandeln Sie beim Öffnen von URLs oder Dateien stets mögliche Fehler:

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

### Benutzerrückmeldung

Geben Sie eine Rückmeldung aus, wenn Vorgänge erfolgreich sind oder fehlschlagen:

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

## Plattformspezifische Aspekte

@tabs
[macOS]
Unter macOS:

- Verwendet den Befehl `open`, um den Standardbrowser zu starten
- Berücksichtigt die Standardbrowser-Einstellung des Benutzers in den Systemeinstellungen
- Kann eine Berechtigungsabfrage anzeigen, wenn die Anwendung in einer Sandbox ausgeführt wird
- Verarbeitet `file://`-URLs für lokale Dateien korrekt

[Windows]
Unter Windows:

- Verwendet die Windows-Shell-API zum Öffnen von URLs
- Berücksichtigt die Standardbrowser-Einstellung in den Windows-Einstellungen
- Verarbeitet Windows-Pfadformate korrekt
- Kann bei nicht vertrauenswürdigen URLs Sicherheitswarnungen anzeigen

[Linux]
Unter Linux:

- Versucht zunächst, `xdg-open` zu verwenden, und greift andernfalls auf andere Methoden zurück
- Das Verhalten variiert je nach Desktop-Umgebung
- Berücksichtigt die Umgebungsvariable `BROWSER`, falls sie gesetzt ist
- Kann bei Minimalinstallationen zusätzliche Pakete erfordern

@end

## Bewährte Vorgehensweisen

1. **Fehler stets behandeln**: Browservorgänge können aus verschiedenen Gründen fehlschlagen:
  ```go
  if err := app.Browser.OpenURL(url); err != nil {
      app.Logger.Error("Failed to open browser", "error", err)
      // Provide fallback or user notification
  }
  ```


2. **URLs validieren**: Stellen Sie vor dem Öffnen sicher, dass URLs korrekt formatiert sind:
  ```go
  func isValidHTTPURL(str string) bool {
      u, err := url.Parse(str)
      return err == nil && (u.Scheme == "http" || u.Scheme == "https")
  }
  ```


3. **Benutzerbestätigung**: Erwägen Sie bei externen Links, den Benutzer um Erlaubnis zu bitten:
  ```go
  // Show confirmation dialog before opening external links
  confirmAndOpen(app, "https://external-site.com")
  ```


4. **Sichere Dateipfade**: Stellen Sie beim Öffnen von Dateien sicher, dass die Pfade sicher sind:
  ```go
  func openSafeFile(app *application.App, filename string) error {
      // Ensure file exists and is readable
      if _, err := os.Stat(filename); err != nil {
          return err
      }
      return app.Browser.OpenFile(filename)
  }
  ```


## Vollständiges Beispiel

Das folgende vollständige Beispiel zeigt verschiedene Muster für die Browserintegration:

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

@note{type="tip" title="Profi-Tipp"}
Erwägen Sie Ausweichmechanismen für den Fall, dass Browservorgänge fehlschlagen, etwa das Kopieren von URLs in die Zwischenablage oder ihre Anzeige in einem Dialog zum manuellen Öffnen.

@end

@note{type="danger" title="Warnung"}
Validieren Sie URLs und Dateipfade stets vor dem Öffnen, um Sicherheitsprobleme zu vermeiden. Öffnen Sie von Benutzern bereitgestellte URLs nicht unbedacht ohne vorherige Validierung.

@end
