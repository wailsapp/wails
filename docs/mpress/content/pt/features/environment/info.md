---
title: "Ambiente"
description: "Acesse informações sobre o ambiente do sistema no seu aplicativo Wails"
slug: "features/environment/info"
sourcePath: "features/environment/info.md"
---

O Wails fornece informações abrangentes sobre o ambiente por meio da API EnvironmentManager. Isso permite que seu aplicativo detecte propriedades do sistema e preferências de tema, além de se integrar ao gerenciador de arquivos do sistema operacional.

## Acessar o gerenciador de ambiente

O gerenciador de ambiente é acessado por meio da propriedade `Env` da instância do seu aplicativo:

```go
app := application.New(application.Options{
    Name: "Environment Demo",
})

// Access the environment manager
env := app.Env
```

## Informações do sistema

### Obter informações do ambiente

Obtenha informações abrangentes sobre o ambiente de execução:

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

### Estrutura do ambiente

As informações do ambiente incluem vários campos importantes:

```go
type EnvironmentInfo struct {
    OS           string                 // Operating system: "windows", "darwin", "linux"
    Arch         string                 // Architecture: "amd64", "arm64", "386", etc.
    Debug        bool                   // Whether running in debug mode
    OSInfo       *operatingsystem.OS    // Detailed OS information
    PlatformInfo map[string]any         // Platform-specific details
}
```

## Detecção de tema

### Detecção do modo escuro

Detecte se o sistema está usando o modo escuro:

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

### Monitoramento de alterações de tema

Monitore alterações de tema para atualizar seu aplicativo dinamicamente:

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

## Integração com o gerenciador de arquivos

### Abrir o gerenciador de arquivos

Abra o gerenciador de arquivos do sistema em um local específico:

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

### Casos de uso comuns

Exiba arquivos ou pastas no gerenciador de arquivos a partir do seu aplicativo:

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

## Comportamento específico da plataforma

### Comportamento adaptativo do aplicativo

Use as informações do ambiente para adaptar o comportamento do seu aplicativo:

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

## Tratamento do modo de depuração

### Desenvolvimento e produção

Use as informações do modo de depuração para habilitar recursos de desenvolvimento:

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

**API privada no macOS:** `OpenDevTools()` requer `private_mac_apis`; caso contrário, a chamada não realiza nenhuma operação. Builds de produção também exigem `devtools`. Consulte a [matriz de builds do Web Inspector](/guides/build/private-macos-apis/#web-inspector).

## Caixa de diálogo de informações do ambiente

### Exibir informações do sistema

Crie uma caixa de diálogo que exiba informações do ambiente:

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

## Considerações sobre as plataformas

@tabs
[macOS]
No macOS:

- A detecção do modo escuro usa as configurações de aparência do sistema
- As operações do gerenciador de arquivos usam o Finder
- As informações da plataforma incluem detalhes da versão do macOS
- A arquitetura pode ser "arm64" em Macs com Apple Silicon

```go
if envInfo.OS == "darwin" {
    // macOS-specific handling
    if envInfo.Arch == "arm64" {
        app.Logger.Info("Running on Apple Silicon")
    }
}
```

[Windows]
No Windows:

- A detecção do modo escuro usa as configurações de tema do Windows
- As operações do gerenciador de arquivos usam o Explorador de Arquivos do Windows
- As informações da plataforma incluem detalhes da versão do Windows
- Podem incluir informações adicionais específicas do Windows

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
No Linux:

- A detecção do modo escuro varia de acordo com o ambiente de desktop
- As operações do gerenciador de arquivos usam o gerenciador de arquivos padrão do sistema
- As informações da plataforma incluem detalhes da distribuição
- O comportamento pode variar entre diferentes distribuições Linux

```go
if envInfo.OS == "linux" {
    // Linux-specific handling
    if distro, ok := envInfo.PlatformInfo["distribution"]; ok {
        app.Logger.Info("Linux distribution", "distro", distro)
    }
}
```

@end

## Práticas recomendadas

1. **Armazene em cache as informações do ambiente**: elas raramente mudam durante a execução:
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


2. **Trate as alterações de tema**: monitore as alterações do tema do sistema:
  ```go
  app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(event *application.ApplicationEvent) {
      updateTheme(app.Env.IsDarkMode())
  })
  ```


3. **Trate corretamente as falhas do gerenciador de arquivos**: sempre trate os erros do gerenciador de arquivos:
  ```go
  func openFileManagerSafely(app *application.App, path string) {
      err := app.Env.OpenFileManager(path, false)
      if err != nil {
          // Provide fallback or user notification
          app.Logger.Warn("Could not open file manager", "path", path)
      }
  }
  ```


4. **Recursos específicos da plataforma**: use as informações do ambiente para habilitar recursos da plataforma:
  ```go
  envInfo := app.Env.Info()
  if envInfo.OS == "darwin" {
      // Enable macOS-specific features
  }
  ```


## Exemplo completo

Veja um exemplo completo que demonstra o gerenciamento do ambiente:

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

@note{type="tip" title="Dica profissional"}
Use as informações do ambiente para oferecer experiências de usuário adequadas a cada plataforma. Por exemplo, use atalhos com a tecla Command no macOS e com a tecla Control no Windows/Linux.

@end

@note{type="danger" title="Aviso"}
Em geral, as informações do ambiente permanecem estáveis durante a execução do aplicativo, mas as preferências de tema podem mudar. Sempre monitore os eventos de alteração de tema para manter sua interface sincronizada.

@end
