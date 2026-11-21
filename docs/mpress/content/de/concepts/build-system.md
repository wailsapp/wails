---
title: "Buildsystem"
description: "So erstellt und paketiert Wails Ihre Anwendung"
slug: "concepts/build-system"
sourcePath: "concepts/build-system.md"
---

## Einheitliches Buildsystem

Wails bietet ein **einheitliches Buildsystem**, das Go-Code kompiliert, Frontend-Assets bündelt, alles in eine einzelne ausführbare Datei einbettet und plattformspezifische Builds verarbeitet – alles mit einem einzigen Befehl.

```bash
wails3 build
```

**Ausgabe:** Native ausführbare Datei, in die alles eingebettet ist.

## Überblick über den Buildprozess

**[Platzhalter für Buildprozessdiagramm]**

## Buildphasen

### 1. Analysephase

Wails durchsucht Ihren Go-Code, um Ihre Services zu erfassen:

```go
type GreetService struct {
    prefix string
}

func (g *GreetService) Greet(name string) string {
    return g.prefix + name + "!"
}
```

**Von Wails extrahierte Informationen:**

- Servicename: `GreetService`
- Methodenname: `Greet`
- Parametertypen: `string`
- Rückgabetypen: `string`

**Verwendung:** Generieren von TypeScript-Bindings

### 2. Generierungsphase

#### TypeScript-Bindings

Wails generiert typsichere Bindings:

```javascript
// Auto-generated: frontend/bindings/<full-go-import-path>/greetservice.js
// (TypeScript is also generated when you pass `-ts`. The shape below is the real
// runtime call format — numeric IDs via $Call.ByID, imported from /wails/runtime.js.)
import { Call as $Call } from "/wails/runtime.js";

export function Greet($0) {
    return $Call.ByID(1234567890, $0);
}
```

**Vorteile:**

- Vollständige Typsicherheit
- Automatische Vervollständigung in der IDE
- Fehler zur Kompilierzeit
- JSDoc-Kommentare

#### Frontend-Build

Ihr Frontend-Bundler wird ausgeführt (Vite, webpack usw.):

```bash
# Vite example
vite build --outDir dist
```

**Ablauf:**

- JavaScript/TypeScript wird kompiliert
- CSS wird verarbeitet und minimiert
- Assets werden optimiert
- Source Maps werden generiert (nur bei der Entwicklung)
- Ausgabe nach `frontend/dist/`

### 3. Kompilierungsphase

#### Go-Kompilierung

Der Go-Code wird mit Optimierungen kompiliert:

```bash
go build -ldflags="-s -w" -o myapp.exe
```

**Flags:**

- `-s`: Symboltabelle entfernen
- `-w`: DWARF-Debuginformationen entfernen
- Ergebnis: kleinere Binärdatei (Reduzierung um ca. 30 %)

**Plattformspezifisch:**

- Windows: `.exe` mit eingebettetem Symbol
- macOS: `.app`-Bundle-Struktur
- Linux: ELF-Binärdatei

#### Einbetten von Assets

Frontend-Assets werden in die Go-Binärdatei eingebettet:

```go
//go:embed frontend/dist
var assets embed.FS
```

**Ergebnis:** eine einzelne ausführbare Datei, die alles enthält.

### 4. Ausgabe

**Einzelne native Binärdatei:**

- Windows: `myapp.exe` (ca. 15 MB)
- macOS: `myapp.app` (ca. 15 MB)
- Linux: `myapp` (ca. 15 MB)

**Keine Abhängigkeiten** (mit Ausnahme des systemeigenen WebViews).

## Entwicklung und Produktion

@tabs{sync-key="mode"}
[Entwicklung (wails3 dev)]
**Für Geschwindigkeit optimiert:**

```bash
wails3 dev
```

**Ablauf:**

1. Startet den Frontend-Entwicklungsserver (standardmäßig Vite auf Port 9245)
2. Kompiliert Go ohne Optimierungen
3. Startet die Anwendung mit Verweis auf den Entwicklungsserver
4. Aktiviert Hot Reload
5. Schließt Source Maps ein

**Merkmale:**

- **Schnelle erneute Builds** (&lt;1 s bei Frontend-Änderungen)
- **Keine Einbettung von Assets** (Bereitstellung über den Entwicklungsserver)
- **Debugsymbole** enthalten
- **Source Maps** aktiviert
- **Ausführliche Protokollierung**

**Dateigröße:** größer (ca. 50 MB mit Debugsymbolen)

[Produktion (wails3 build)]
**Für Größe und Leistung optimiert:**

```bash
wails3 build
```

**Ablauf:**

1. Erstellt das Frontend für die Produktion (minimiert)
2. Kompiliert Go mit Optimierungen
3. Entfernt Debugsymbole
4. Bettet Assets ein
5. Erstellt eine einzelne Binärdatei

**Merkmale:**

- **Optimierter Code** (minimiert, Tree Shaking angewendet)
- **Eingebettete Assets** (keine externen Dateien)
- **Debugsymbole entfernt**
- **Keine Source Maps**
- **Minimale Protokollierung**

**Dateigröße:** kleiner (ca. 15 MB)

@end

## Buildbefehle

### Einfacher Build

```bash
wails3 build
```

**Ausgabe:** `bin/<APP_NAME>` (oder `bin/<APP_NAME>.exe` unter Windows). Das Verzeichnis `bin/` befindet sich im Projektstammverzeichnis.

`wails3 build` ist ein schlanker Wrapper für `wails3 task build`. Das einzige zur Buildzeit weitergereichte Flag ist `--tags`, das zur Taskfile-Variable `EXTRA_TAGS` wird:

```bash
# Build with extra Go build tags
wails3 build --tags "myfeature,gtk4"
```

Für `wails3 build` gibt es keine Flags `-platform`, `-o`, `-skipbindings`, `-clean`, `-debug`, `-devbuild`, `-icon`, `-ldflags` oder `-package`. Cross-Kompilierung, Ausgabepfade, Symbole und Paketierung werden über das Taskfile des Projekts gesteuert (`Taskfile.yml` + `build/config.yml`).

### Plattformübergreifende und plattformspezifische Builds

Plattform-Builds stehen als Taskfile-Tasks in den Namensräumen `darwin:` / `windows:` / `linux:` zur Verfügung (definiert in `build/Taskfile.<platform>.yml`). Beispiel:

```bash
# macOS — universal binary
wails3 task darwin:build:universal

# macOS — current arch
wails3 task darwin:build

# Windows
wails3 task windows:build

# Linux
wails3 task linux:build
```

So zeigen Sie alle im aktuellen Projekt verfügbaren Tasks an:

```bash
wails3 task --list
```

### Symbole und Paketierung

Generieren Sie Plattformsymbole (`build/icons.icns`, `build/icon.ico` usw.) aus einer Quell-PNG-Datei:

```bash
wails3 generate icons -input appicon.png
```

Erstellen Sie plattformspezifische Installationsprogramme und Pakete:

```bash
wails3 package           # uses the current Go build env
wails3 task linux:create:deb
wails3 task windows:package
wails3 task darwin:package:universal
```

## Build-Konfiguration

### Taskfile.yml

Wails-3-Projekte verwenden das [Taskfile](https://taskfile.dev/) zur Build-Orchestrierung. Das `Taskfile.yml` im Stammverzeichnis bindet plattformspezifische Task-Dateien aus `build/` ein:

```yaml
# Taskfile.yml (excerpt — the real templates are richer)
version: '3'

includes:
  common: ./build/Taskfile.yml
  darwin: ./build/Taskfile.darwin.yml
  windows: ./build/Taskfile.windows.yml
  linux: ./build/Taskfile.linux.yml

tasks:
  build:
    desc: Build the application
    cmds:
      - task: "{{OS}}:build"
```

Führen Sie Tasks mit `wails3 task <name>` oder `task <name>` aus:

```bash
wails3 task windows:build
wails3 task darwin:package:universal
wails3 task linux:create:appimage
```

### Projektkonfiguration: `build/config.yml`

Projektmetadaten (Name, Kennung, Version, Info.plist-Werte, NSIS-Einstellungen, `.desktop`-Felder, benutzerdefinierte Protokolle usw.) befinden sich in `build/config.yml`. Das Taskfile liest diese Datei beim Generieren von Symbolen, Manifesten, Installationsprogrammen und Ähnlichem. In Wails 3 gibt es **keine** `build/build.json`-Datei.

```yaml
# build/config.yml (illustrative)
info:
  productName: "My App"
  productIdentifier: "com.example.myapp"
  productVersion: "1.0.0"
  companyName: "Example Ltd."
  productDescription: "An application built with Wails"
```

Führen Sie `wails3 generate build-assets` (oder `wails3 update build-assets`) aus, um die plattformspezifischen Build-Ressourcen anhand dieser Konfiguration zu aktualisieren.

## Einbetten von Ressourcen

### Funktionsweise

Wails verwendet das Go-Paket `embed`:

```go
package main

import (
    "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name:   "My App",
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
    })
    
    app.Window.New()
    app.Run()
}
```

**Zur Build-Zeit:**

1. Frontend wird nach `frontend/dist/` gebaut
2. Die `//go:embed`-Direktive bindet Dateien ein
3. Dateien werden in die Binärdatei kompiliert
4. Die Binärdatei enthält alles

**Zur Laufzeit:**

1. Anwendung startet
2. Ressourcen werden aus dem Arbeitsspeicher bereitgestellt
3. Keine Datenträger-E/A für Ressourcen
4. Schnelles Laden

### Benutzerdefinierte Ressourcen

Betten Sie zusätzliche Dateien ein:

```go
//go:embed frontend/dist
var frontendAssets embed.FS

//go:embed data/*.json
var dataAssets embed.FS

//go:embed templates/*.html
var templateAssets embed.FS
```

## Build-Optimierungen

### Frontend-Optimierungen

**Vite (Standard):**

```javascript
// vite.config.js
export default {
  build: {
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: true,  // Remove console.log
        drop_debugger: true,
      },
    },
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['react', 'react-dom'],  // Separate vendor bundle
        },
      },
    },
  },
}
```

**Ergebnisse:**

- JavaScript minimiert (Reduzierung um ca. 70 %)
- CSS minimiert (Reduzierung um ca. 60 %)
- Bilder optimiert
- Tree-Shaking angewendet

### Go-Optimierungen

**Compiler-Flags:**

```bash
-ldflags="-s -w"
```

- `-s`: Symboltabelle entfernen (Reduzierung um ca. 10 %)
- `-w`: DWARF-Debuginformationen entfernen (Reduzierung um ca. 20 %)

**Zusätzliche Optimierungen:**

```bash
-ldflags="-s -w -X main.version=1.0.0"
```

- `-X`: Variablenwerte zur Build-Zeit festlegen
- Nützlich für Versionsnummern und Build-Daten

### Binärkomprimierung

**UPX (optional):**

```bash
# After building
upx --best bin/myapp.exe
```

**Ergebnisse:**

- Größenreduzierung um ca. 50 %
- Etwas langsamerer Start (ca. 100 ms)
- Für macOS nicht empfohlen (Probleme mit der Codesignierung)

## Plattformspezifische Builds

### Windows

**Ausgabe:** `myapp.exe`

**Enthält:**

- Anwendungssymbol
- Versionsinformationen
- Manifest (UAC-Einstellungen)

**Symbol:**

```bash
# Generate platform icons from a source PNG
wails3 generate icons -input appicon.png -windowsfilename build/icon.ico
```

Der Windows-Schritt `tool package` bettet anschließend die generierte Datei `.ico` in die ausführbare Datei ein.

**Manifest:**

```xml
<!-- build/windows/manifest.xml -->
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">
  <assemblyIdentity version="1.0.0.0" name="MyApp"/>
  <trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
      <requestedPrivileges>
        <requestedExecutionLevel level="asInvoker" uiAccess="false"/>
      </requestedPrivileges>
    </security>
  </trustInfo>
</assembly>
```

### macOS

**Ausgabe:** `myapp.app` (Anwendungspaket)

**Struktur:**

```
myapp.app/
├── Contents/
│   ├── Info.plist          # App metadata
│   ├── MacOS/
│   │   └── myapp           # Binary
│   ├── Resources/
│   │   └── icon.icns       # Icon
│   └── _CodeSignature/     # Code signature (if signed)
```

**Info.plist:**

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleName</key>
    <string>My App</string>
    <key>CFBundleIdentifier</key>
    <string>com.example.myapp</string>
    <key>CFBundleVersion</key>
    <string>1.0.0</string>
</dict>
</plist>
```

**Universal Binary:**

Das macOS-Taskfile stellt einen Task `darwin:build:universal` (sowie `darwin:package:universal`) bereit, der beide Architekturen baut und sie mit `wails3 tool lipo` kombiniert:

```bash
wails3 task darwin:build:universal
```

### Linux

**Ausgabe:** `myapp` (ELF-Binärdatei)

**Abhängigkeiten:**

- GTK3
- WebKitGTK

**Desktop-Datei:**

```ini
# myapp.desktop
[Desktop Entry]
Name=My App
Exec=/usr/bin/myapp
Icon=myapp
Type=Application
Categories=Utility;
```

**Installation:**

```bash
# Copy binary
sudo cp myapp /usr/bin/

# Copy desktop file
sudo cp myapp.desktop /usr/share/applications/

# Copy icon
sudo cp icon.png /usr/share/icons/hicolor/256x256/apps/myapp.png
```

## Build-Performance

### Typische Build-Zeiten

| Phase | Dauer | Hinweise |
| --- | --- | --- |
| Analyse | &lt;1 s | Scannen des Go-Codes |
| Generierung der Bindings | &lt;1 s | TypeScript-Generierung |
| Frontend-Build | 5-30 s | Abhängig von der Projektgröße |
| Go-Kompilierung | 2-10 s | Abhängig von der Codegröße |
| Einbetten von Ressourcen | &lt;1 s | Einbetten des Frontends |
| **Gesamt** | **10-45 s** | Erster Build |
| **Inkrementell** | **5-15 s** | Nachfolgende Builds |

### Builds beschleunigen

**1. Build-Cache verwenden:**

```bash
# Go build cache is automatic
# Frontend cache (Vite)
npm run build  # Uses cache by default
```

**2. Nur das Nötige ausführen:**

```bash
# Pick the specific Taskfile target you actually need
wails3 task common:build:frontend   # rebuild only the frontend
wails3 task windows:build           # rebuild only the Windows binary
```

**3. Parallele Builds (mehrere Rechner/CI):**

Die Cross-Kompilierung zwischen Linux, Windows und macOS in v3 erfolgt im Allgemeinen in einem Docker-`wails-cross`-Container oder auf dedizierten Runnern für die jeweilige Plattform – `wails3 build` selbst ist auf das Hostbetriebssystem ausgerichtet. Informationen zu den unterstützten Workflows finden Sie unter [Plattformübergreifende Builds](/guides/build/cross-platform/).

**4. Schnellere Tools verwenden:**

```bash
# Use esbuild instead of webpack
# (Vite uses esbuild by default)
```

## Fehlerbehebung

### Build schlägt fehl

**Symptom:** `wails3 build` wird mit einem Fehler beendet

**Häufige Ursachen:**

1. **Go-Kompilierungsfehler**
  ```bash
  # Check Go code compiles
  go build
  ```


2. **Frontend-Build-Fehler**
  ```bash
  # Check frontend builds
  cd frontend
  npm run build
  ```


3. **Fehlende Abhängigkeiten**
  ```bash
  # Install dependencies
  npm install
  go mod download
  ```


### Binärdatei zu groß

**Symptom:** Die Binärdatei ist >50 MB groß

**Lösungen:**

1. **Debug-Symbole entfernen** (das mitgelieferte Taskfile übergibt bereits `-ldflags="-s -w"` an `go build`).

2. **Eingebettete Assets prüfen**
  ```bash
  # Remove unnecessary files from frontend/dist/
  # Check for large images, videos, etc.
  ```


3. **UPX-Komprimierung verwenden**
  ```bash
  upx --best bin/myapp.exe
  ```


### Langsame Builds

**Symptom:** Builds dauern länger als 1 Minute

**Lösungen:**

1. **Build-Cache verwenden**
  - Der Go-Cache wird automatisch verwendet
  - Der Frontend-Cache (Vite) wird automatisch verwendet


2. **Nur den benötigten Task ausführen**
  ```bash
  wails3 task common:build:frontend
  wails3 task windows:build
  ```


3. **Frontend-Build optimieren**
  ```javascript
  // vite.config.js
  export default {
    build: {
      minify: 'esbuild',  // Faster than terser
    },
  }
  ```


## Bewährte Vorgehensweisen

### ✅ Empfohlen

- **Während der Entwicklung `wails3 dev` verwenden** – Schnelle Iterationen
- **Für Releases `wails3 build` verwenden** – Optimierte Ausgabe
- **Builds versionieren** – Die Version mit `-ldflags` einbetten
- **Builds auf den Zielplattformen testen** – Cross-Kompilierung ist nicht perfekt
- **Frontend-Builds schnell halten** – Bundler-Konfiguration optimieren
- **Build-Cache verwenden** – Beschleunigt nachfolgende Builds

### ❌ Nicht empfohlen

- **Das Verzeichnis `build/` nicht committen** – Zu `.gitignore` hinzufügen
- **Build-Tests nicht überspringen** – Vor dem Release immer testen
- **Keine unnötigen Assets einbetten** – Binärdateien klein halten
- **Keine Debug-Builds für die Produktion verwenden** – Optimierte Builds verwenden
- **Code-Signierung nicht vergessen** – Für die Distribution erforderlich

## Nächste Schritte

**Anwendungen erstellen** – Ausführliche Anleitung zum Erstellen und Paketieren [Mehr erfahren →](/guides/build/building/)

**Plattformübergreifende Builds** – Von einem Rechner aus für alle Plattformen erstellen [Mehr erfahren →](/guides/build/cross-platform/)

**Installationsprogramme erstellen** – Installationsprogramme für Endbenutzer erstellen [Mehr erfahren →](/guides/installers/)

---

**Fragen zum Erstellen von Builds?** Fragen Sie auf [Discord](https://discord.gg/JDdSxwjhGf) nach oder sehen Sie sich die [Build-Beispiele](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/build) an.
