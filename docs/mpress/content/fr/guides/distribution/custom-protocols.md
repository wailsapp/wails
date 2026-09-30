---
title: "Protocoles d’URL personnalisés"
description: "Enregistrez des schémas d’URL personnalisés pour lancer votre application depuis des liens"
slug: "guides/distribution/custom-protocols"
sourcePath: "guides/distribution/custom-protocols.md"
---

Les protocoles d’URL personnalisés (également appelés schémas d’URL) permettent de lancer votre application lorsque les utilisateurs cliquent sur des liens utilisant votre protocole personnalisé, tels que `myapp://action` ou `myapp://open/document`.

## Vue d’ensemble

Les protocoles personnalisés permettent :

- **Liens profonds** : lancez votre application avec des données précises
- **Intégration au navigateur** : gérez les liens provenant de pages web
- **Liens dans les e-mails** : ouvrez votre application depuis des clients de messagerie
- **Communication entre applications** : lancez votre application depuis d’autres applications

**Exemple** : `myapp://open/document?id=123` lance votre application et ouvre le document 123.

## Configuration

Définissez les protocoles personnalisés dans les options de votre application :

Les protocoles personnalisés sont déclarés dans `build/config.yml`, que les outils de création de paquets propres à chaque plateforme — macros NSIS sous Windows, manifeste MSIX, `CFBundleURLTypes` sous macOS, `.desktop`/`xdg-mime` sous Linux — utilisent lors de la création du paquet. Il n’existe ni type `application.Protocol` ni champ `Protocols` dans `application.Options`.

```yaml
# build/config.yml
protocols:
  - scheme: myapp
    description: "My Application Protocol"
```

Dans le code Go, écoutez l’événement `ApplicationLaunchedWithUrl` pour détecter un lancement avec une URL :

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

## Gestionnaire de protocole

Écoutez les événements de protocole pour traiter les URL entrantes :

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

## Structure des URL

Concevez des structures d’URL claires et hiérarchiques :

```
myapp://action/resource?param=value

Examples:
myapp://open/document?id=123
myapp://settings/theme?mode=dark
myapp://user/profile?username=john
```

**Bonnes pratiques :**

- Utilisez des noms de schéma en minuscules
- Choisissez des schémas courts et faciles à retenir
- Utilisez des chemins hiérarchiques pour les ressources
- Utilisez des paramètres de requête pour les données facultatives
- Encodez les caractères spéciaux dans les URL

## Enregistrement propre à chaque plateforme

Les protocoles personnalisés sont enregistrés différemment sur chaque plateforme.

@tabs{sync-key="platform"}
[Windows]
### Programme d’installation NSIS pour Windows

**Wails v3 enregistre automatiquement les protocoles personnalisés** lorsque vous utilisez un programme d’installation NSIS.

#### Enregistrement automatique

Lorsque vous compilez votre application avec `wails3 build`, le programme d’installation NSIS :

1. Enregistre automatiquement tous les protocoles déclarés dans `build/config.yml` sous la clé `protocols:`
2. Associe les protocoles au fichier exécutable de votre application
3. Crée les entrées de registre appropriées
4. Supprime les associations de protocoles lors de la désinstallation

**Aucune configuration supplémentaire n’est requise !**

#### Fonctionnement

Le modèle NSIS comprend des macros intégrées :

- `wails.associateCustomProtocols` - Enregistre les protocoles lors de l’installation
- `wails.unassociateCustomProtocols` - Supprime les protocoles lors de la désinstallation

Ces macros sont appelées automatiquement en fonction de votre configuration `Protocols`.

#### Registre manuel (avancé)

Si vous devez effectuer un enregistrement manuel (sans NSIS) :

```batch
@echo off
REM Register custom protocol
REG ADD "HKEY_CURRENT_USER\SOFTWARE\Classes\myapp" /ve /d "URL:My Application Protocol" /f
REG ADD "HKEY_CURRENT_USER\SOFTWARE\Classes\myapp" /v "URL Protocol" /t REG_SZ /d "" /f
REG ADD "HKEY_CURRENT_USER\SOFTWARE\Classes\myapp\shell\open\command" /ve /d "\"%1\"" /f
```

#### Test

Testez l’enregistrement de votre protocole :

```powershell
# Open protocol URL from PowerShell
Start-Process "myapp://test/action"

# Or from command prompt
start myapp://test/action
```

### Paquet MSIX pour Windows

Les protocoles personnalisés sont également enregistrés automatiquement lorsque vous utilisez un paquet MSIX.

#### Enregistrement automatique

Lorsque vous compilez votre application au format MSIX, le manifeste inclut automatiquement les enregistrements de protocoles provenant de votre configuration des protocoles `build/config.yml`.

Le manifeste généré comprend :

```xml
<uap:Extension Category="windows.protocol">
  <uap:Protocol Name="myapp">
    <uap:DisplayName>My Application Protocol</uap:DisplayName>
  </uap:Protocol>
</uap:Extension>
```

#### Liens universels (liens entre le Web et l’application)

Windows prend en charge les **liens entre le Web et l’application**, qui fonctionnent de manière similaire aux liens universels sous macOS. Lorsque vous déployez votre application sous forme de paquet MSIX, vous pouvez autoriser des liens HTTPS à lancer directement votre application.

@note{type="note"}
Les liens entre le Web et l’application nécessitent une configuration manuelle du manifeste. Les schémas de protocoles personnalisés sont configurés automatiquement depuis `build/config.yml`, mais les domaines associés doivent être ajoutés manuellement à votre manifeste MSIX.

@end

Pour activer les liens entre le Web et l’application, suivez le [guide de Microsoft sur les liens entre le Web et l’application](https://learn.microsoft.com/en-us/windows/apps/develop/launch/web-to-app-linking). Vous devrez :

1. **Ajouter manuellement App URI Handler à votre manifeste MSIX** (`build/windows/msix/app_manifest.xml`) :
  ```xml
  <uap3:Extension Category="windows.appUriHandler">
    <uap3:AppUriHandler>
      <uap3:Host Name="myawesomeapp.com"/>
    </uap3:AppUriHandler>
  </uap3:Extension>
  ```


2. **Configurez `windows-app-web-link` sur votre site web :** hébergez un fichier `windows-app-web-link` à l’adresse `https://myawesomeapp.com/.well-known/windows-app-web-link`. Ce fichier devrait contenir les informations sur le paquet de votre application et les chemins qu’elle gère.

Lorsqu’un lien entre le Web et l’application lance votre application, vous recevez le même événement `ApplicationLaunchedWithUrl` qu’avec les schémas de protocoles personnalisés.

[macOS]
### Configuration du fichier Info.plist

Sous macOS, les protocoles sont enregistrés au moyen de votre fichier `Info.plist`.

#### Configuration automatique

Wails génère automatiquement le fichier `Info.plist` avec vos protocoles lorsque vous effectuez une compilation avec `wails3 build`.

Les protocoles déclarés dans `build/config.yml` sont ajoutés à :

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

#### Tests

```bash
# Open protocol URL from terminal
open "myapp://test/action"

# Check registered handlers
/System/Library/Frameworks/CoreServices.framework/Versions/A/Frameworks/LaunchServices.framework/Versions/A/Support/lsregister -dump | grep myapp
```

#### Liens universels

Outre les schémas de protocole personnalisés, macOS prend également en charge les **liens universels**, qui permettent de lancer votre application à l’aide de liens HTTPS ordinaires (par exemple, `https://myawesomeapp.com/path`). Les liens universels offrent une expérience utilisateur fluide entre votre application web et votre application de bureau.

@note{type="caution"}
Pour prendre en charge les liens universels, votre application macOS doit être **signée numériquement** à l’aide d’un certificat Apple Developer et d’un profil de provisionnement valides. Les versions non signées ou dotées d’une signature ad hoc ne pourront pas ouvrir de liens universels. Assurez-vous que votre application est correctement signée avant de la tester.

@end

Pour activer les liens universels, suivez le [guide d’Apple sur leur prise en charge dans votre application](https://developer.apple.com/documentation/xcode/supporting-universal-links-in-your-app). Vous devrez :

1. **Ajouter des droits** dans votre fichier `entitlements.plist` :
  ```xml
  <key>com.apple.developer.associated-domains</key>
  <array>
    <string>applinks:myawesomeapp.com</string>
  </array>
  ```


2. **Ajouter NSUserActivityTypes à Info.plist** :
  ```xml
  <key>NSUserActivityTypes</key>
  <array>
    <string>NSUserActivityTypeBrowsingWeb</string>
  </array>
  ```


3. **Configurer `apple-app-site-association` sur votre site web :** hébergez un fichier `apple-app-site-association` à l’emplacement `https://myawesomeapp.com/.well-known/apple-app-site-association`.

Lorsqu’un lien universel déclenche votre application, vous recevez le même événement `ApplicationLaunchedWithUrl` ; le code de traitement est donc identique à celui des schémas de protocole personnalisés.

[Linux]
### Entrée de bureau

Sous Linux, les protocoles sont enregistrés au moyen de fichiers `.desktop`.

#### Configuration automatique

Wails génère un fichier d’entrée de bureau contenant les gestionnaires de protocoles lorsque vous effectuez une compilation avec `wails3 build`.

**Corrigé dans la v3** : le modèle de bureau Linux inclut désormais correctement la gestion des protocoles.

Le fichier de bureau généré contient :

```ini
[Desktop Entry]
Type=Application
Name=My Application
Exec=/usr/bin/myapp %u
MimeType=x-scheme-handler/myapp;
```

#### Enregistrement manuel

Si nécessaire, installez manuellement le fichier de bureau :

```bash
# Copy desktop file
cp myapp.desktop ~/.local/share/applications/

# Update desktop database
update-desktop-database ~/.local/share/applications/

# Register protocol handler
xdg-mime default myapp.desktop x-scheme-handler/myapp
```

#### Tests

```bash
# Open protocol URL
xdg-open "myapp://test/action"

# Check registered handler
xdg-mime query default x-scheme-handler/myapp
```

@end

## Exemple complet

Voici un exemple complet qui traite plusieurs actions de protocole :

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

## Intégration au frontend

Gérez les événements de navigation dans votre frontend :

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

## Considérations de sécurité

### Valider toutes les entrées

Validez et nettoyez toujours les URL provenant de sources externes :

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

### Prévenir les attaques par injection

N’exécutez jamais directement les URL en tant que code ou requêtes SQL :

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

## Tests

### Tests manuels

Testez les gestionnaires de protocoles pendant le développement :

**Windows :**

```powershell
Start-Process "myapp://test/action?id=123"
```

**macOS :**

```bash
open "myapp://test/action?id=123"
```

**Linux :**

```bash
xdg-open "myapp://test/action?id=123"
```

### Tests HTML

Créez une page HTML de test :

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

## Dépannage

### Protocole non enregistré

**Windows :**

- Vérifiez le registre : `HKEY_CURRENT_USER\SOFTWARE\Classes\<scheme>`
- Réinstallez l’application avec le programme d’installation NSIS
- Vérifiez que le programme d’installation a été exécuté avec les autorisations appropriées

**macOS :**

- Recompilez l’application avec `wails3 build`
- Vérifiez `Info.plist` dans le paquet de l’application : `MyApp.app/Contents/Info.plist`
- Réinitialisez Launch Services : `/System/Library/Frameworks/CoreServices.framework/Frameworks/LaunchServices.framework/Support/lsregister -kill`

**Linux :**

- Vérifiez le fichier de bureau : `~/.local/share/applications/myapp.desktop`
- Mettez à jour la base de données : `update-desktop-database ~/.local/share/applications/`
- Vérifiez le gestionnaire : `xdg-mime query default x-scheme-handler/myapp`

### L’application ne se lance pas

**Consultez les journaux :**

```go
app := application.New(application.Options{
    LogLevel: slog.LevelDebug, // requires `import "log/slog"`
    // ...
})
```

**Problèmes courants :**

- L’application n’est pas installée à l’emplacement attendu
- Le chemin de l’exécutable indiqué lors de l’enregistrement ne correspond pas à son emplacement réel
- Problèmes d’autorisations

## Bonnes pratiques

### ✅ À faire

- **Utilisez des noms de schéma explicites** — `mycompany-myapp` plutôt que `mca`
- **Validez toutes les entrées** — Ne faites jamais confiance aux URL provenant de sources externes
- **Gérez les erreurs sans interrompre l’application** — Consignez les URL non valides dans les journaux ; ne provoquez pas de plantage
- **Fournissez un retour à l’utilisateur** — Indiquez quelle action a été déclenchée
- **Testez sur toutes les plateformes** — La gestion des protocoles varie selon la plateforme
- **Documentez la structure de vos URL** — Aidez les utilisateurs et les intégrateurs

### ❌ À éviter

- **N’utilisez pas de noms de schéma courants** — Évitez `http`, `file`, `app`, etc.
- **N’exécutez pas les URL comme du code** — Cela présente un risque de sécurité majeur
- **N’exposez pas d’opérations sensibles** — Exigez une confirmation pour les actions destructrices
- **Ne supposez pas que les protocoles fonctionnent partout** — Prévoyez des mécanismes de secours
- **N’oubliez pas l’encodage des URL** — Gérez correctement les caractères spéciaux

## Étapes suivantes

- [Création de packages Windows](/guides/build/windows/) — Découvrez les options du programme d’installation NSIS
- [Associations de fichiers](/guides/file-associations/) — Ouvrez des fichiers avec votre application
- [Instance unique](/guides/single-instance/) — Empêchez l’exécution de plusieurs instances de l’application

---

**Des questions ?** Posez-les sur [Discord](https://discord.gg/JDdSxwjhGf) ou consultez les [exemples](https://github.com/wailsapp/wails/tree/master/v3/examples).
