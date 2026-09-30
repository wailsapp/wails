---
title: "Lingkungan"
description: "Akses informasi lingkungan sistem di aplikasi Wails Anda"
slug: "features/environment/info"
sourcePath: "features/environment/info.md"
---

Wails menyediakan informasi lingkungan yang lengkap melalui API EnvironmentManager. Dengan API ini, aplikasi Anda dapat mendeteksi properti sistem dan preferensi tema, serta terintegrasi dengan pengelola file sistem operasi.

## Mengakses Pengelola Lingkungan

Pengelola lingkungan diakses melalui properti `Env` pada instans aplikasi Anda:

```go
app := application.New(application.Options{
    Name: "Environment Demo",
})

// Access the environment manager
env := app.Env
```

## Informasi Sistem

### Mendapatkan Informasi Lingkungan

Dapatkan informasi lengkap tentang lingkungan runtime:

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

### Struktur Lingkungan

Informasi lingkungan mencakup beberapa bidang penting:

```go
type EnvironmentInfo struct {
    OS           string                 // Operating system: "windows", "darwin", "linux"
    Arch         string                 // Architecture: "amd64", "arm64", "386", etc.
    Debug        bool                   // Whether running in debug mode
    OSInfo       *operatingsystem.OS    // Detailed OS information
    PlatformInfo map[string]any         // Platform-specific details
}
```

## Deteksi Tema

### Deteksi Mode Gelap

Deteksi apakah sistem menggunakan mode gelap:

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

### Pemantauan Perubahan Tema

Pantau perubahan tema untuk memperbarui aplikasi Anda secara dinamis:

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

## Integrasi Pengelola File

### Membuka Pengelola File

Buka pengelola file sistem di lokasi tertentu:

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

### Kasus Penggunaan Umum

Tampilkan file atau folder di pengelola file dari aplikasi Anda:

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

## Perilaku Khusus Platform

### Perilaku Aplikasi Adaptif

Gunakan informasi lingkungan untuk menyesuaikan perilaku aplikasi Anda:

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

## Penanganan Mode Debug

### Pengembangan vs Produksi

Gunakan informasi mode debug untuk mengaktifkan fitur pengembangan:

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

**API privat di macOS:** `OpenDevTools()` memerlukan `private_mac_apis`; jika tidak, pemanggilan tersebut tidak melakukan apa pun. Build produksi juga memerlukan `devtools`. Lihat [matriks build Web Inspector](/guides/build/private-macos-apis/#web-inspector).

## Dialog Informasi Lingkungan

### Menampilkan Informasi Sistem

Buat dialog yang menampilkan informasi lingkungan:

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

## Pertimbangan Platform

@tabs
[macOS]
Di macOS:

- Deteksi mode gelap menggunakan pengaturan tampilan sistem
- Operasi pengelola file menggunakan Finder
- Informasi platform mencakup detail versi macOS
- Arsitektur dapat berupa "arm64" pada Mac dengan Apple Silicon

```go
if envInfo.OS == "darwin" {
    // macOS-specific handling
    if envInfo.Arch == "arm64" {
        app.Logger.Info("Running on Apple Silicon")
    }
}
```

[Windows]
Di Windows:

- Deteksi mode gelap menggunakan pengaturan tema Windows
- Operasi pengelola file menggunakan Windows Explorer
- Informasi platform mencakup detail versi Windows
- Dapat mencakup informasi tambahan khusus Windows

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
Di Linux:

- Deteksi mode gelap berbeda-beda menurut lingkungan desktop
- Operasi pengelola file menggunakan pengelola file default sistem
- Informasi platform mencakup detail distribusi
- Perilaku dapat berbeda di antara berbagai distribusi Linux

```go
if envInfo.OS == "linux" {
    // Linux-specific handling
    if distro, ok := envInfo.PlatformInfo["distribution"]; ok {
        app.Logger.Info("Linux distribution", "distro", distro)
    }
}
```

@end

## Praktik Terbaik

1. **Simpan Informasi Lingkungan dalam Cache**: Informasi lingkungan jarang berubah selama runtime:
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


2. **Tangani Perubahan Tema**: Pantau perubahan tema sistem:
  ```go
  app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(event *application.ApplicationEvent) {
      updateTheme(app.Env.IsDarkMode())
  })
  ```


3. **Tangani Kegagalan Pengelola File dengan Baik**: Selalu tangani kesalahan pengelola file:
  ```go
  func openFileManagerSafely(app *application.App, path string) {
      err := app.Env.OpenFileManager(path, false)
      if err != nil {
          // Provide fallback or user notification
          app.Logger.Warn("Could not open file manager", "path", path)
      }
  }
  ```


4. **Fitur Khusus Platform**: Gunakan informasi lingkungan untuk mengaktifkan fitur platform:
  ```go
  envInfo := app.Env.Info()
  if envInfo.OS == "darwin" {
      // Enable macOS-specific features
  }
  ```


## Contoh Lengkap

Berikut contoh lengkap yang mendemonstrasikan pengelolaan lingkungan:

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

@note{type="tip" title="Kiat Pro"}
Gunakan informasi lingkungan untuk memberikan pengalaman pengguna yang sesuai dengan platform. Misalnya, gunakan pintasan dengan tombol Command di macOS dan pintasan dengan tombol Control di Windows/Linux.

@end

@note{type="danger" title="Peringatan"}
Informasi lingkungan umumnya tetap stabil selama runtime aplikasi, tetapi preferensi tema dapat berubah. Selalu pantau peristiwa perubahan tema agar UI Anda tetap tersinkronisasi.

@end
