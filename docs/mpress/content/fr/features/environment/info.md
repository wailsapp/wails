---
title: "Environnement"
description: "Accédez aux informations sur l’environnement système dans votre application Wails"
slug: "features/environment/info"
sourcePath: "features/environment/info.md"
---

Wails fournit des informations complètes sur l’environnement par l’intermédiaire de l’API EnvironmentManager. Votre application peut ainsi détecter les propriétés du système et les préférences de thème, et s’intégrer au gestionnaire de fichiers du système d’exploitation.

## Accès au gestionnaire d’environnement

Accédez au gestionnaire d’environnement par l’intermédiaire de la propriété `Env` de votre instance d’application :

```go
app := application.New(application.Options{
    Name: "Environment Demo",
})

// Access the environment manager
env := app.Env
```

## Informations système

### Obtenir les informations sur l’environnement

Récupérez des informations complètes sur l’environnement d’exécution :

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

### Structure de l’environnement

Les informations sur l’environnement comprennent plusieurs champs importants :

```go
type EnvironmentInfo struct {
    OS           string                 // Operating system: "windows", "darwin", "linux"
    Arch         string                 // Architecture: "amd64", "arm64", "386", etc.
    Debug        bool                   // Whether running in debug mode
    OSInfo       *operatingsystem.OS    // Detailed OS information
    PlatformInfo map[string]any         // Platform-specific details
}
```

## Détection du thème

### Détection du mode sombre

Détectez si le système utilise le mode sombre :

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

### Surveillance des changements de thème

Détectez les changements de thème pour mettre à jour dynamiquement votre application :

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

## Intégration au gestionnaire de fichiers

### Ouvrir le gestionnaire de fichiers

Ouvrez le gestionnaire de fichiers du système à un emplacement précis :

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

### Cas d’usage courants

Affichez des fichiers ou des dossiers dans le gestionnaire de fichiers depuis votre application :

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

## Comportement propre à chaque plateforme

### Comportement adaptatif de l’application

Utilisez les informations sur l’environnement pour adapter le comportement de votre application :

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

## Gestion du mode débogage

### Développement et production

Utilisez les informations sur le mode débogage pour activer les fonctionnalités de développement :

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

**API privée sous macOS :** `OpenDevTools()` nécessite `private_mac_apis` ; sinon, l’appel n’a aucun effet. Les builds de production nécessitent également `devtools`. Consultez la [matrice des builds de l’inspecteur web](/guides/build/private-macos-apis/#web-inspector).

## Boîte de dialogue d’informations sur l’environnement

### Afficher les informations système

Créez une boîte de dialogue affichant les informations sur l’environnement :

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

## Considérations propres aux plateformes

@tabs
[macOS]
Sous macOS :

- La détection du mode sombre utilise les réglages d’apparence du système
- Les opérations du gestionnaire de fichiers utilisent le Finder
- Les informations sur la plateforme comprennent des détails sur la version de macOS
- L’architecture peut être « arm64 » sur les Mac équipés d’une puce Apple Silicon

```go
if envInfo.OS == "darwin" {
    // macOS-specific handling
    if envInfo.Arch == "arm64" {
        app.Logger.Info("Running on Apple Silicon")
    }
}
```

[Windows]
Sous Windows :

- La détection du mode sombre utilise les paramètres de thème de Windows
- Les opérations du gestionnaire de fichiers utilisent l’Explorateur Windows
- Les informations sur la plateforme comprennent des détails sur la version de Windows
- Elles peuvent contenir des informations supplémentaires propres à Windows

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
Sous Linux :

- La détection du mode sombre varie selon l’environnement de bureau
- Les opérations du gestionnaire de fichiers utilisent le gestionnaire de fichiers par défaut du système
- Les informations sur la plateforme comprennent des détails sur la distribution
- Le comportement peut varier d’une distribution Linux à l’autre

```go
if envInfo.OS == "linux" {
    // Linux-specific handling
    if distro, ok := envInfo.PlatformInfo["distribution"]; ok {
        app.Logger.Info("Linux distribution", "distro", distro)
    }
}
```

@end

## Bonnes pratiques

1. **Mettre en cache les informations sur l’environnement** : les informations sur l’environnement changent rarement pendant l’exécution :
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


2. **Gérer les changements de thème** : détectez les changements du thème système :
  ```go
  app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(event *application.ApplicationEvent) {
      updateTheme(app.Env.IsDarkMode())
  })
  ```


3. **Gérer sans interruption les échecs du gestionnaire de fichiers** : gérez toujours les erreurs du gestionnaire de fichiers :
  ```go
  func openFileManagerSafely(app *application.App, path string) {
      err := app.Env.OpenFileManager(path, false)
      if err != nil {
          // Provide fallback or user notification
          app.Logger.Warn("Could not open file manager", "path", path)
      }
  }
  ```


4. **Fonctionnalités propres à chaque plateforme** : utilisez les informations sur l’environnement pour activer les fonctionnalités de la plateforme :
  ```go
  envInfo := app.Env.Info()
  if envInfo.OS == "darwin" {
      // Enable macOS-specific features
  }
  ```


## Exemple complet

Voici un exemple complet illustrant la gestion de l’environnement :

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

@note{type="tip" title="Conseil de pro"}
Utilisez les informations sur l’environnement pour proposer une expérience utilisateur adaptée à la plateforme. Par exemple, utilisez les raccourcis avec la touche Command sous macOS et ceux avec la touche Control sous Windows et Linux.

@end

@note{type="danger" title="Avertissement"}
Les informations sur l’environnement restent généralement stables pendant l’exécution de l’application, mais les préférences de thème peuvent changer. Détectez toujours les événements de changement de thème afin de maintenir votre interface utilisateur synchronisée.

@end
