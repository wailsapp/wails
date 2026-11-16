---
title: "Benutzerdefinierte URL-Protokolle"
description: "Benutzerdefinierte URL-Schemata registrieren, um Ihre Anwendung über Links zu starten"
slug: "guides/distribution/custom-protocols"
sourcePath: "guides/distribution/custom-protocols.md"
---

Mit benutzerdefinierten URL-Protokollen (auch URL-Schemata genannt) kann Ihre Anwendung gestartet werden, wenn Benutzer auf Links mit Ihrem benutzerdefinierten Protokoll klicken, beispielsweise `myapp://action` oder `myapp://open/document`.

## Übersicht

Benutzerdefinierte Protokolle ermöglichen:

- **Deep Linking**: Ihre App mit bestimmten Daten starten
- **Browserintegration**: Links von Webseiten verarbeiten
- **E-Mail-Links**: Ihre App aus E-Mail-Clients öffnen
- **Kommunikation zwischen Apps**: Ihre App aus anderen Anwendungen starten

**Beispiel**: `myapp://open/document?id=123` startet Ihre App und öffnet das Dokument 123.

## Konfiguration

Definieren Sie benutzerdefinierte Protokolle in den Optionen Ihrer Anwendung:

Benutzerdefinierte Protokolle werden in `build/config.yml` deklariert. Die Paketierungswerkzeuge der Plattformen – NSIS-Makros unter Windows, das MSIX-Manifest, `CFBundleURLTypes` unter macOS sowie `.desktop`/`xdg-mime` unter Linux – verarbeiten sie beim Paketieren. Es gibt weder einen Typ `application.Protocol` noch ein Feld `Protocols` in `application.Options`.

```yaml
# build/config.yml
protocols:
  - scheme: myapp
    description: "My Application Protocol"
```

Überwachen Sie im Go-Code das Ereignis `ApplicationLaunchedWithUrl`, um Starts über eine URL zu erkennen:

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

## Protokoll-Handler

Überwachen Sie Protokollereignisse, um eingehende URLs zu verarbeiten:

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

## URL-Struktur

Entwerfen Sie klare, hierarchische URL-Strukturen:

```
myapp://action/resource?param=value

Examples:
myapp://open/document?id=123
myapp://settings/theme?mode=dark
myapp://user/profile?username=john
```

**Bewährte Vorgehensweisen:**

- Verwenden Sie kleingeschriebene Schemanamen
- Halten Sie Schemata kurz und einprägsam
- Verwenden Sie hierarchische Pfade für Ressourcen
- Verwenden Sie Abfrageparameter für optionale Daten
- Kodieren Sie Sonderzeichen URL-konform

## Plattformregistrierung

Benutzerdefinierte Protokolle werden auf jeder Plattform unterschiedlich registriert.

@tabs{sync-key="platform"}
[Windows]
### Windows-NSIS-Installationsprogramm

**Wails v3 registriert benutzerdefinierte Protokolle bei Verwendung von NSIS-Installationsprogrammen automatisch**.

#### Automatische Registrierung

Wenn Sie Ihre Anwendung mit `wails3 build` erstellen, führt das NSIS-Installationsprogramm Folgendes aus:

1. Es registriert automatisch alle in `build/config.yml` deklarierten Protokolle unter dem Schlüssel `protocols:`
2. Es verknüpft die Protokolle mit der ausführbaren Datei Ihrer Anwendung
3. Es richtet die erforderlichen Registrierungseinträge ein
4. Es entfernt die Protokollzuordnungen bei der Deinstallation

**Keine zusätzliche Konfiguration erforderlich!**

#### Funktionsweise

Die NSIS-Vorlage enthält integrierte Makros:

- `wails.associateCustomProtocols` – Registriert Protokolle während der Installation
- `wails.unassociateCustomProtocols` – Entfernt Protokolle während der Deinstallation

Diese Makros werden anhand Ihrer `Protocols`-Konfiguration automatisch aufgerufen.

#### Manuelle Registrierung (für Fortgeschrittene)

Wenn Sie eine manuelle Registrierung benötigen (außerhalb von NSIS):

```batch
@echo off
REM Register custom protocol
REG ADD "HKEY_CURRENT_USER\SOFTWARE\Classes\myapp" /ve /d "URL:My Application Protocol" /f
REG ADD "HKEY_CURRENT_USER\SOFTWARE\Classes\myapp" /v "URL Protocol" /t REG_SZ /d "" /f
REG ADD "HKEY_CURRENT_USER\SOFTWARE\Classes\myapp\shell\open\command" /ve /d "\"%1\"" /f
```

#### Testen

Testen Sie Ihre Protokollregistrierung:

```powershell
# Open protocol URL from PowerShell
Start-Process "myapp://test/action"

# Or from command prompt
start myapp://test/action
```

### Windows-MSIX-Paket

Auch bei der MSIX-Paketierung werden benutzerdefinierte Protokolle automatisch registriert.

#### Automatische Registrierung

Wenn Sie Ihre Anwendung mit MSIX erstellen, übernimmt das Manifest automatisch die Protokollregistrierungen aus Ihrer `build/config.yml`-Protokollkonfiguration.

Das generierte Manifest enthält:

```xml
<uap:Extension Category="windows.protocol">
  <uap:Protocol Name="myapp">
    <uap:DisplayName>My Application Protocol</uap:DisplayName>
  </uap:Protocol>
</uap:Extension>
```

#### Universal Links (Verknüpfung vom Web zur App)

Windows unterstützt die **Verknüpfung vom Web zur App**, die ähnlich wie Universal Links unter macOS funktioniert. Wenn Sie Ihre Anwendung als MSIX-Paket bereitstellen, können Sie festlegen, dass HTTPS-Links Ihre App direkt starten.

@note{type="note"}
Die Verknüpfung vom Web zur App erfordert eine manuelle Manifestkonfiguration. Benutzerdefinierte Protokollschemata werden automatisch anhand von `build/config.yml` konfiguriert, zugeordnete Domains müssen Sie Ihrem MSIX-Manifest jedoch manuell hinzufügen.

@end

Um die Verknüpfung vom Web zur App zu aktivieren, folgen Sie dem [Microsoft-Leitfaden zur Verknüpfung vom Web zur App](https://learn.microsoft.com/en-us/windows/apps/develop/launch/web-to-app-linking). Sie müssen:

1. **App URI Handler manuell Ihrem MSIX-Manifest hinzufügen** (`build/windows/msix/app_manifest.xml`):
  ```xml
  <uap3:Extension Category="windows.appUriHandler">
    <uap3:AppUriHandler>
      <uap3:Host Name="myawesomeapp.com"/>
    </uap3:AppUriHandler>
  </uap3:Extension>
  ```


2. **Konfigurieren Sie `windows-app-web-link` auf Ihrer Website:** Stellen Sie unter `https://myawesomeapp.com/.well-known/windows-app-web-link` eine Datei `windows-app-web-link` bereit. Diese Datei sollte die Paketinformationen Ihrer App und die von ihr verarbeiteten Pfade enthalten.

Wenn Ihre Anwendung über einen Web-to-App-Link gestartet wird, empfangen Sie dasselbe Ereignis `ApplicationLaunchedWithUrl` wie bei benutzerdefinierten Protokollschemata.

[macOS]
### Info.plist-Konfiguration

Unter macOS werden Protokolle über Ihre Datei `Info.plist` registriert.

#### Automatische Konfiguration

Wenn Sie mit `wails3 build` erstellen, generiert Wails automatisch die Datei `Info.plist` mit Ihren Protokollen.

Die in `build/config.yml` deklarierten Protokolle werden hinzugefügt zu:

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

#### Testen

```bash
# Open protocol URL from terminal
open "myapp://test/action"

# Check registered handlers
/System/Library/Frameworks/CoreServices.framework/Versions/A/Frameworks/LaunchServices.framework/Versions/A/Support/lsregister -dump | grep myapp
```

#### Universal Links

Neben benutzerdefinierten Protokollschemata unterstützt macOS auch **Universal Links**, über die Ihre App durch gewöhnliche HTTPS-Links gestartet werden kann (z. B. `https://myawesomeapp.com/path`). Universal Links ermöglichen eine nahtlose Benutzererfahrung zwischen Ihrer Web- und Desktop-App.

@note{type="caution"}
Für Universal Links muss Ihre macOS-App mit einem gültigen Apple-Developer-Zertifikat und Bereitstellungsprofil **codesigniert** sein. Nicht signierte oder ad hoc signierte Builds können Universal Links nicht öffnen. Stellen Sie vor dem Testen sicher, dass Ihre App ordnungsgemäß signiert ist.

@end

Um Universal Links zu aktivieren, folgen Sie dem [Apple-Leitfaden zur Unterstützung von Universal Links in Ihrer App](https://developer.apple.com/documentation/xcode/supporting-universal-links-in-your-app). Sie müssen:

1. **Entitlements** in Ihrer Datei `entitlements.plist` hinzufügen:
  ```xml
  <key>com.apple.developer.associated-domains</key>
  <array>
    <string>applinks:myawesomeapp.com</string>
  </array>
  ```


2. **NSUserActivityTypes zu Info.plist hinzufügen**:
  ```xml
  <key>NSUserActivityTypes</key>
  <array>
    <string>NSUserActivityTypeBrowsingWeb</string>
  </array>
  ```


3. **Konfigurieren Sie `apple-app-site-association` auf Ihrer Website:** Stellen Sie unter `https://myawesomeapp.com/.well-known/apple-app-site-association` eine Datei `apple-app-site-association` bereit.

Wenn Ihre App durch einen Universal Link ausgelöst wird, empfangen Sie dasselbe Ereignis `ApplicationLaunchedWithUrl`. Daher ist der Verarbeitungscode mit dem für benutzerdefinierte Protokollschemata identisch.

[Linux]
### Desktop-Eintrag

Unter Linux werden Protokolle über `.desktop`-Dateien registriert.

#### Automatische Konfiguration

Wenn Sie mit `wails3 build` erstellen, generiert Wails eine Desktop-Eintragsdatei mit Protokoll-Handlern.

**Behoben in v3**: Die Linux-Desktopvorlage unterstützt jetzt Protokoll-Handler korrekt.

Die generierte Desktopdatei enthält:

```ini
[Desktop Entry]
Type=Application
Name=My Application
Exec=/usr/bin/myapp %u
MimeType=x-scheme-handler/myapp;
```

#### Manuelle Registrierung

Installieren Sie die Desktopdatei bei Bedarf manuell:

```bash
# Copy desktop file
cp myapp.desktop ~/.local/share/applications/

# Update desktop database
update-desktop-database ~/.local/share/applications/

# Register protocol handler
xdg-mime default myapp.desktop x-scheme-handler/myapp
```

#### Testen

```bash
# Open protocol URL
xdg-open "myapp://test/action"

# Check registered handler
xdg-mime query default x-scheme-handler/myapp
```

@end

## Vollständiges Beispiel

Das folgende vollständige Beispiel verarbeitet mehrere Protokollaktionen:

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

## Frontend-Integration

Verarbeiten Sie Navigationsereignisse in Ihrem Frontend:

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

## Sicherheitsaspekte

### Alle Eingaben validieren

Validieren und bereinigen Sie URLs aus externen Quellen immer:

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

### Injection-Angriffe verhindern

Führen Sie URLs niemals direkt als Code oder SQL aus:

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

## Testen

### Manuelles Testen

Testen Sie Protokoll-Handler während der Entwicklung:

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

### HTML-Test

Erstellen Sie eine HTML-Testseite:

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

## Fehlerbehebung

### Protokoll nicht registriert

**Windows:**

- Überprüfen Sie die Registrierung: `HKEY_CURRENT_USER\SOFTWARE\Classes\<scheme>`
- Installieren Sie die Anwendung mit dem NSIS-Installer neu
- Prüfen Sie, ob der Installer mit den erforderlichen Berechtigungen ausgeführt wurde

**macOS:**

- Erstellen Sie die Anwendung mit `wails3 build` neu
- Überprüfen Sie `Info.plist` im App-Bundle: `MyApp.app/Contents/Info.plist`
- Setzen Sie Launch Services zurück: `/System/Library/Frameworks/CoreServices.framework/Frameworks/LaunchServices.framework/Support/lsregister -kill`

**Linux:**

- Überprüfen Sie die Desktopdatei: `~/.local/share/applications/myapp.desktop`
- Aktualisieren Sie die Datenbank: `update-desktop-database ~/.local/share/applications/`
- Überprüfen Sie den Handler: `xdg-mime query default x-scheme-handler/myapp`

### Anwendung startet nicht

**Überprüfen Sie die Protokolle:**

```go
app := application.New(application.Options{
    LogLevel: slog.LevelDebug, // requires `import "log/slog"`
    // ...
})
```

**Häufige Probleme:**

- Die Anwendung ist nicht am erwarteten Speicherort installiert
- Der registrierte Pfad zur ausführbaren Datei stimmt nicht mit dem tatsächlichen Speicherort überein
- Berechtigungsprobleme

## Bewährte Vorgehensweisen

### ✅ Empfohlen

- **Verwenden Sie aussagekräftige Schemanamen** – `mycompany-myapp` statt `mca`
- **Validieren Sie alle Eingaben** – Vertrauen Sie niemals URLs aus externen Quellen
- **Behandeln Sie Fehler ordnungsgemäß** – Protokollieren Sie ungültige URLs, statt die Anwendung abstürzen zu lassen
- **Geben Sie den Benutzern Rückmeldung** – Zeigen Sie, welche Aktion ausgelöst wurde
- **Testen Sie auf allen Plattformen** – Die Protokollverarbeitung unterscheidet sich je nach Plattform
- **Dokumentieren Sie Ihre URL-Struktur** – Unterstützen Sie Benutzer und Integratoren

### ❌ Nicht empfohlen

- **Verwenden Sie keine gängigen Schemanamen** – Vermeiden Sie `http`, `file`, `app` usw.
- **Führen Sie URLs nicht als Code aus** – Dies stellt ein enormes Sicherheitsrisiko dar
- **Machen Sie keine sensiblen Vorgänge zugänglich** – Verlangen Sie bei destruktiven Aktionen eine Bestätigung
- **Gehen Sie nicht davon aus, dass Protokolle überall funktionieren** – Stellen Sie Ausweichmechanismen bereit
- **Vergessen Sie die URL-Codierung nicht** – Verarbeiten Sie Sonderzeichen korrekt

## Nächste Schritte

- [Windows-Paketierung](/guides/build/windows/) – Erfahren Sie mehr über die Optionen des NSIS-Installers
- [Dateizuordnungen](/guides/file-associations/) – Öffnen Sie Dateien mit Ihrer Anwendung
- [Einzelinstanz](/guides/single-instance/) – Verhindern Sie mehrere Anwendungsinstanzen

---

**Fragen?** Fragen Sie auf [Discord](https://discord.gg/JDdSxwjhGf) oder sehen Sie sich die [Beispiele](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples) an.
