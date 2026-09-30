---
title: "Entwicklungsumgebung einrichten"
description: "Entwicklungsumgebung für die Entwicklung von Wails v3 einrichten"
slug: "contributing/setup"
sourcePath: "contributing/setup.md"
---

## Entwicklungsumgebung einrichten

Diese Anleitung führt Sie durch die Einrichtung einer vollständigen Entwicklungsumgebung für die Arbeit an Wails v3.

## Erforderliche Werkzeuge

### Go-Entwicklung

1. **Installieren Sie Go 1.25 oder höher:**
  ```bash
  # Download from https://go.dev/dl/
  go version  # Verify installation
  ```


2. **Konfigurieren Sie die Go-Umgebung:**
  ```bash
  # Add to your shell profile (.bashrc, .zshrc, etc.)
  export GOPATH=$HOME/go
  export PATH=$PATH:$GOPATH/bin
  ```


3. **Installieren Sie nützliche Go-Werkzeuge:**
  ```bash
  go install golang.org/x/tools/cmd/goimports@latest
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  ```


### Node.js und npm

Nur für Beispiele zur Frontend-Integration erforderlich.

```bash
# Install Node.js 20+ and npm
node --version  # Should be 20+
npm --version
```

### Plattformspezifische Abhängigkeiten

**macOS:**

```bash
# Install Xcode Command Line Tools
xcode-select --install

# Verify installation
xcode-select -p  # Should output a path
```

**Windows:**

1. Installieren Sie [MSYS2](https://www.msys2.org/) für eine Unix-ähnliche Umgebung
2. WebView2 Runtime (unter Windows 11 vorinstalliert; für Windows 10 [herunterladen](https://developer.microsoft.com/en-us/microsoft-edge/webview2/))
3. Optional: Installieren Sie [Git for Windows](https://git-scm.com/download/win)

**Linux (Debian/Ubuntu):**

```bash
sudo apt update
# Default GTK4 + WebKitGTK 6.0 stack (Ubuntu 24.04+ / Debian 13+)
sudo apt install build-essential pkg-config libgtk-4-dev libwebkitgtk-6.0-dev
# For the legacy -tags gtk3 path:
sudo apt install libgtk-3-dev libwebkit2gtk-4.1-dev
```

**Linux (Fedora/RHEL):**

```bash
# Default GTK4 stack
sudo dnf install gcc pkg-config gtk4-devel webkitgtk6.0-devel
# Legacy GTK3 path:
sudo dnf install gtk3-devel webkit2gtk4.1-devel
```

**Linux (Arch):**

```bash
# Default GTK4 stack
sudo pacman -S base-devel gtk4 webkitgtk-6.0
# Legacy GTK3 path:
sudo pacman -S gtk3 webkit2gtk-4.1
```

## Repository einrichten

### Klonen und konfigurieren

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails

# Add upstream remote
git remote add upstream https://github.com/wailsapp/wails.git

# Verify remotes
git remote -v
```

### Wails CLI erstellen

```bash
# Navigate to v3 directory
cd v3

# Build the CLI
go build -o ../wails3 ./cmd/wails3

# Test the build
cd ..
./wails3 version
```

### Zum PATH hinzufügen (optional)

**Linux/macOS:**

```bash
# Add to ~/.bashrc or ~/.zshrc
export PATH=$PATH:/path/to/wails
```

**Windows:**

Fügen Sie das Wails-Verzeichnis über die Systemeigenschaften zur Umgebungsvariable PATH hinzu.

## IDE einrichten

### VS Code (empfohlen)

1. **Installieren Sie VS Code:** [Herunterladen](https://code.visualstudio.com/)

2. **Installieren Sie Erweiterungen:**
  - Go (vom Go Team bei Google)
  - ESLint
  - Prettier
  - MDX (für die Dokumentation)


3. **Konfigurieren Sie die Arbeitsbereichseinstellungen** (`.vscode/settings.json`):
  ```json
  {
    "go.useLanguageServer": true,
    "go.lintTool": "golangci-lint",
    "go.lintOnSave": "workspace",
    "editor.formatOnSave": true,
    "go.formatTool": "goimports"
  }
  ```


### GoLand

1. **Installieren Sie GoLand:** [Herunterladen](https://www.jetbrains.com/go/)

2. **Konfigurieren Sie Folgendes:**
  - Aktivieren Sie die Unterstützung für Go-Module
  - Richten Sie Datei-Watcher für `goimports` ein
  - Konfigurieren Sie den Codestil gemäß den Projektkonventionen


## Einrichtung überprüfen

Führen Sie diese Befehle aus, um zu überprüfen, ob alles funktioniert:

```bash
# Go version check
go version

# Build Wails
cd v3
go build ./cmd/wails3

# Run tests
go test ./pkg/...

# Create a test app
cd ..
./wails3 init -n mytest -t vanilla
cd mytest
../wails3 dev
```

Wenn die Testanwendung erstellt und ausgeführt wird, ist Ihre Umgebung einsatzbereit!

## Tests ausführen

### Unit-Tests

```bash
cd v3
go test ./...
cd ..
```

### Tests für bestimmte Pakete

```bash
cd v3
go test ./pkg/application
go test ./pkg/events -v  # Verbose output
```

### Mit Codeabdeckung ausführen

```bash
cd v3
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Mit Race Detector ausführen

```bash
cd v3
go test ./... -race
```

## Mit der Dokumentation arbeiten

Die Wails-v3-Dokumentation ist in M-Press verfasst. Die englischen Quelldateien befinden sich in  
`docs/mpress/content/`; die Übersetzungen befinden sich in Sprachverzeichnissen wie  
`fr/` und `id/`.

Zeigen Sie Dokumentationsänderungen in der Vorschau an und validieren Sie sie im Stammverzeichnis des Repositorys:

```bash
go install github.com/leaanthony/mpress/cmd/mpress@v1.0.17
mpress dev
mpress build --strict --no-purge-css
mpress check
# Use `python` on Windows if `python3` is unavailable.
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

Die Produktionswebsite ist statisch. Für die lokale Arbeit an der Dokumentation sind Node.js, ein Übersetzungsanbieter und Cloudflare-Anmeldedaten nicht erforderlich.

## Debugging

### Go-Code debuggen

**VS Code:**

Erstellen Sie `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug Wails CLI",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/v3/cmd/wails3",
      "args": ["dev"]
    }
  ]
}
```

**Befehlszeile:**

```bash
# Use Delve debugger
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug ./cmd/wails3 -- dev
```

### Plattformcode debuggen

Für plattformspezifisches Debugging sind plattformspezifische Werkzeuge erforderlich:

- **macOS:** Xcode Instruments
- **Windows:** Visual Studio Debugger
- **Linux:** GDB

## Häufige Probleme

### „command not found: wails3“

Fügen Sie das Wails-Verzeichnis zu Ihrem PATH hinzu oder verwenden Sie `./wails3` aus dem Projektstammverzeichnis.

### „webkitgtk-6.0 not found“ oder „webkit2gtk not found“ (Linux)

Installieren Sie die Entwicklungspakete für den Stack, für den Sie die Anwendung erstellen:

```bash
# Default GTK4 stack (Debian/Ubuntu):
sudo apt install libwebkitgtk-6.0-dev

# Legacy GTK3 path:
sudo apt install libwebkit2gtk-4.1-dev
```

### Build schlägt aufgrund von Go-Modulfehlern fehl

```bash
cd v3
go mod tidy
go mod download
```

### „CGO_ENABLED“-Fehler unter Windows

Stellen Sie sicher, dass ein C-Compiler (MinGW-w64 über MSYS2) in Ihrem PATH vorhanden ist.

## Nächste Schritte

- Lesen Sie die [Programmierrichtlinien](/contributing/standards/).
- Sehen Sie sich die [technische Dokumentation](/contributing/) an.
- Suchen Sie nach einem Issue, an dem Sie arbeiten können: [Issues für den Einstieg](https://github.com/wailsapp/wails/labels/good%20first%20issue)
