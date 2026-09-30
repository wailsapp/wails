---
title: "Protokol URL Kustom"
description: "Daftarkan skema URL kustom untuk meluncurkan aplikasi Anda dari tautan"
slug: "guides/distribution/custom-protocols"
sourcePath: "guides/distribution/custom-protocols.md"
---

Protokol URL kustom (juga disebut skema URL) memungkinkan aplikasi Anda diluncurkan saat pengguna mengeklik tautan dengan protokol kustom Anda, seperti `myapp://action` atau `myapp://open/document`.

## Ikhtisar

Protokol kustom memungkinkan:

- **Tautan dalam**: Meluncurkan aplikasi Anda dengan data tertentu
- **Integrasi browser**: Menangani tautan dari halaman web
- **Tautan email**: Membuka aplikasi Anda dari klien email
- **Komunikasi antar-aplikasi**: Meluncurkan aplikasi dari aplikasi lain

**Contoh**: `myapp://open/document?id=123` meluncurkan aplikasi Anda dan membuka dokumen 123.

## Konfigurasi

Tentukan protokol kustom dalam opsi aplikasi Anda:

Protokol kustom dideklarasikan dalam `build/config.yml` (yang digunakan oleh pembuat paket platform—makro NSIS di Windows, manifes MSIX, `CFBundleURLTypes` di macOS, serta `.desktop`/`xdg-mime` di Linux—saat pembuatan paket). Tidak ada tipe `application.Protocol` maupun bidang `Protocols` pada `application.Options`.

```yaml
# build/config.yml
protocols:
  - scheme: myapp
    description: "My Application Protocol"
```

Dalam kode Go, pantau peluncuran dengan URL melalui peristiwa `ApplicationLaunchedWithUrl`:

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

## Penangan Protokol

Pantau peristiwa protokol untuk menangani URL yang masuk:

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

## Struktur URL

Rancang struktur URL hierarkis yang jelas:

```
myapp://action/resource?param=value

Examples:
myapp://open/document?id=123
myapp://settings/theme?mode=dark
myapp://user/profile?username=john
```

**Praktik terbaik:**

- Gunakan nama skema dengan huruf kecil
- Buat skema tetap singkat dan mudah diingat
- Gunakan jalur hierarkis untuk sumber daya
- Sertakan parameter kueri untuk data opsional
- Enkodekan karakter khusus untuk URL

## Pendaftaran Platform

Protokol kustom didaftarkan secara berbeda pada setiap platform.

@tabs{sync-key="platform"}
[Windows]
### Penginstal NSIS Windows

**Wails v3 secara otomatis mendaftarkan protokol kustom** saat menggunakan penginstal NSIS.

#### Pendaftaran Otomatis

Saat Anda membangun aplikasi dengan `wails3 build`, penginstal NSIS akan:

1. Secara otomatis mendaftarkan semua protokol yang dideklarasikan dalam `build/config.yml` di bawah kunci `protocols:`
2. Mengaitkan protokol dengan berkas aplikasi Anda yang dapat dieksekusi
3. Menyiapkan entri registri yang tepat
4. Menghapus pengaitan protokol selama penghapusan instalasi

**Tidak diperlukan konfigurasi tambahan!**

#### Cara Kerjanya

Templat NSIS menyertakan makro bawaan:

- `wails.associateCustomProtocols` - Mendaftarkan protokol selama instalasi
- `wails.unassociateCustomProtocols` - Menghapus protokol selama penghapusan instalasi

Makro ini dipanggil secara otomatis berdasarkan konfigurasi `Protocols` Anda.

#### Registri Manual (Lanjutan)

Jika Anda memerlukan pendaftaran manual (di luar NSIS):

```batch
@echo off
REM Register custom protocol
REG ADD "HKEY_CURRENT_USER\SOFTWARE\Classes\myapp" /ve /d "URL:My Application Protocol" /f
REG ADD "HKEY_CURRENT_USER\SOFTWARE\Classes\myapp" /v "URL Protocol" /t REG_SZ /d "" /f
REG ADD "HKEY_CURRENT_USER\SOFTWARE\Classes\myapp\shell\open\command" /ve /d "\"%1\"" /f
```

#### Pengujian

Uji pendaftaran protokol Anda:

```powershell
# Open protocol URL from PowerShell
Start-Process "myapp://test/action"

# Or from command prompt
start myapp://test/action
```

### Paket MSIX Windows

Protokol kustom juga didaftarkan secara otomatis saat menggunakan paket MSIX.

#### Pendaftaran Otomatis

Saat Anda membangun aplikasi dengan MSIX, manifes secara otomatis menyertakan pendaftaran protokol dari konfigurasi protokol `build/config.yml` Anda.

Manifes yang dihasilkan mencakup:

```xml
<uap:Extension Category="windows.protocol">
  <uap:Protocol Name="myapp">
    <uap:DisplayName>My Application Protocol</uap:DisplayName>
  </uap:Protocol>
</uap:Extension>
```

#### Universal Links (Penautan Web ke Aplikasi)

Windows mendukung **penautan Web ke Aplikasi**, yang bekerja serupa dengan Universal Links di macOS. Saat menerapkan aplikasi sebagai paket MSIX, Anda dapat mengaktifkan tautan HTTPS untuk meluncurkan aplikasi secara langsung.

@note{type="note"}
Penautan Web ke Aplikasi memerlukan konfigurasi manifes secara manual. Skema protokol kustom dikonfigurasi secara otomatis dari `build/config.yml`, tetapi domain terkait harus ditambahkan secara manual ke manifes MSIX Anda.

@end

Untuk mengaktifkan penautan Web ke Aplikasi, ikuti [panduan Microsoft tentang penautan web ke aplikasi](https://learn.microsoft.com/en-us/windows/apps/develop/launch/web-to-app-linking). Anda perlu:

1. **Menambahkan App URI Handler secara manual ke manifes MSIX Anda** (`build/windows/msix/app_manifest.xml`):
  ```xml
  <uap3:Extension Category="windows.appUriHandler">
    <uap3:AppUriHandler>
      <uap3:Host Name="myawesomeapp.com"/>
    </uap3:AppUriHandler>
  </uap3:Extension>
  ```


2. **Konfigurasikan `windows-app-web-link` di situs web Anda:** Sediakan berkas `windows-app-web-link` di `https://myawesomeapp.com/.well-known/windows-app-web-link`. Berkas ini sebaiknya berisi informasi paket aplikasi Anda dan jalur yang ditanganinya.

Saat tautan Web ke Aplikasi meluncurkan aplikasi Anda, Anda akan menerima peristiwa `ApplicationLaunchedWithUrl` yang sama seperti pada skema protokol kustom.

[macOS]
### Konfigurasi Info.plist

Di macOS, protokol didaftarkan melalui berkas `Info.plist` Anda.

#### Konfigurasi Otomatis

Wails secara otomatis menghasilkan `Info.plist` yang berisi protokol Anda saat Anda membangun dengan `wails3 build`.

Protokol yang dideklarasikan dalam `build/config.yml` ditambahkan ke:

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

#### Pengujian

```bash
# Open protocol URL from terminal
open "myapp://test/action"

# Check registered handlers
/System/Library/Frameworks/CoreServices.framework/Versions/A/Frameworks/LaunchServices.framework/Versions/A/Support/lsregister -dump | grep myapp
```

#### Universal Links

Selain skema protokol khusus, macOS juga mendukung **Universal Links**, yang memungkinkan aplikasi Anda diluncurkan melalui tautan HTTPS biasa (misalnya, `https://myawesomeapp.com/path`). Universal Links menghadirkan pengalaman pengguna yang mulus antara aplikasi web dan aplikasi desktop Anda.

@note{type="caution"}
Universal Links mengharuskan aplikasi macOS Anda **ditandatangani secara digital** dengan sertifikat Apple Developer dan profil penyediaan yang valid. Build yang tidak ditandatangani atau ditandatangani secara ad hoc tidak akan dapat membuka Universal Links. Pastikan aplikasi Anda ditandatangani dengan benar sebelum melakukan pengujian.

@end

Untuk mengaktifkan Universal Links, ikuti [panduan Apple tentang dukungan Universal Links di aplikasi Anda](https://developer.apple.com/documentation/xcode/supporting-universal-links-in-your-app). Anda perlu:

1. **Tambahkan entitlement** dalam `entitlements.plist` Anda:
  ```xml
  <key>com.apple.developer.associated-domains</key>
  <array>
    <string>applinks:myawesomeapp.com</string>
  </array>
  ```


2. **Tambahkan NSUserActivityTypes ke Info.plist**:
  ```xml
  <key>NSUserActivityTypes</key>
  <array>
    <string>NSUserActivityTypeBrowsingWeb</string>
  </array>
  ```


3. **Konfigurasikan `apple-app-site-association` di situs web Anda:** Sediakan file `apple-app-site-association` di `https://myawesomeapp.com/.well-known/apple-app-site-association`.

Saat Universal Link memicu aplikasi Anda, Anda akan menerima peristiwa `ApplicationLaunchedWithUrl` yang sama, sehingga kode penanganannya identik dengan skema protokol khusus.

[Linux]
### Entri Desktop

Di Linux, protokol didaftarkan melalui file `.desktop`.

#### Konfigurasi Otomatis

Wails menghasilkan file entri desktop yang berisi penangan protokol saat Anda membangun dengan `wails3 build`.

**Diperbaiki di v3**: Templat desktop Linux kini menyertakan penanganan protokol dengan benar.

File desktop yang dihasilkan mencakup:

```ini
[Desktop Entry]
Type=Application
Name=My Application
Exec=/usr/bin/myapp %u
MimeType=x-scheme-handler/myapp;
```

#### Pendaftaran Manual

Jika diperlukan, instal file desktop secara manual:

```bash
# Copy desktop file
cp myapp.desktop ~/.local/share/applications/

# Update desktop database
update-desktop-database ~/.local/share/applications/

# Register protocol handler
xdg-mime default myapp.desktop x-scheme-handler/myapp
```

#### Pengujian

```bash
# Open protocol URL
xdg-open "myapp://test/action"

# Check registered handler
xdg-mime query default x-scheme-handler/myapp
```

@end

## Contoh Lengkap

Berikut contoh lengkap untuk menangani beberapa tindakan protokol:

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

## Integrasi Frontend

Tangani peristiwa navigasi di frontend Anda:

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

## Pertimbangan Keamanan

### Validasi Semua Input

Selalu validasi dan sanitasi URL dari sumber eksternal:

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

### Cegah Serangan Injeksi

Jangan pernah mengeksekusi URL secara langsung sebagai kode atau SQL:

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

## Pengujian

### Pengujian Manual

Uji penangan protokol selama pengembangan:

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

### Pengujian HTML

Buat halaman HTML pengujian:

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

## Pemecahan Masalah

### Protokol Tidak Terdaftar

**Windows:**

- Periksa registri: `HKEY_CURRENT_USER\SOFTWARE\Classes\<scheme>`
- Instal ulang menggunakan penginstal NSIS
- Pastikan penginstal dijalankan dengan izin yang sesuai

**macOS:**

- Bangun ulang aplikasi dengan `wails3 build`
- Periksa `Info.plist` dalam bundel aplikasi: `MyApp.app/Contents/Info.plist`
- Atur ulang Launch Services: `/System/Library/Frameworks/CoreServices.framework/Frameworks/LaunchServices.framework/Support/lsregister -kill`

**Linux:**

- Periksa file desktop: `~/.local/share/applications/myapp.desktop`
- Perbarui basis data: `update-desktop-database ~/.local/share/applications/`
- Verifikasi penangan: `xdg-mime query default x-scheme-handler/myapp`

### Aplikasi Tidak Dapat Diluncurkan

**Periksa log:**

```go
app := application.New(application.Options{
    LogLevel: slog.LevelDebug, // requires `import "log/slog"`
    // ...
})
```

**Masalah umum:**

- Aplikasi tidak terinstal di lokasi yang semestinya
- Jalur berkas yang dapat dieksekusi dalam pendaftaran tidak sesuai dengan lokasi sebenarnya
- Masalah izin

## Praktik Terbaik

### ✅ Lakukan

- **Gunakan nama skema yang deskriptif** — gunakan `mycompany-myapp`, bukan `mca`
- **Validasi semua input** — Jangan pernah memercayai URL dari sumber eksternal
- **Tangani kesalahan dengan baik** — Catat URL yang tidak valid ke log; jangan biarkan aplikasi mengalami crash
- **Berikan umpan balik kepada pengguna** - Tampilkan tindakan yang dipicu
- **Uji di semua platform** - Penanganan protokol berbeda-beda
- **Dokumentasikan struktur URL Anda** - Bantu pengguna dan integrator

### ❌ Jangan Lakukan

- **Jangan gunakan nama skema yang umum** - Hindari `http`, `file`, `app`, dan sebagainya.
- **Jangan jalankan URL sebagai kode** - Risikonya terhadap keamanan sangat besar
- **Jangan ekspos operasi sensitif** - Wajibkan konfirmasi untuk tindakan destruktif
- **Jangan berasumsi bahwa protokol berfungsi di semua lingkungan** - Sediakan mekanisme cadangan
- **Jangan lupa melakukan pengodean URL** - Tangani karakter khusus dengan benar

## Langkah Berikutnya

- [Pemaketan Windows](/guides/build/windows/) - Pelajari opsi penginstal NSIS
- [Asosiasi File](/guides/file-associations/) - Buka file dengan aplikasi Anda
- [Instans Tunggal](/guides/single-instance/) - Cegah beberapa instans aplikasi berjalan sekaligus

---

**Ada pertanyaan?** Tanyakan di [Discord](https://discord.gg/JDdSxwjhGf) atau lihat [contoh](https://github.com/wailsapp/wails/tree/master/v3/examples).
