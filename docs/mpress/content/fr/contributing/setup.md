---
title: "Configuration de l’environnement de développement"
description: "Configurez votre environnement de développement pour contribuer au développement de Wails v3"
slug: "contributing/setup"
sourcePath: "contributing/setup.md"
---

## Configuration de l’environnement de développement

Ce guide vous accompagne dans la configuration d’un environnement de développement complet pour travailler sur Wails v3.

## Outils requis

### Développement en Go

1. **Installez Go 1.26 ou une version ultérieure :**
  ```bash
  # Download from https://go.dev/dl/
  go version  # Verify installation
  ```


2. **Configurez l’environnement Go :**
  ```bash
  # Add to your shell profile (.bashrc, .zshrc, etc.)
  export GOPATH=$HOME/go
  export PATH=$PATH:$GOPATH/bin
  ```


3. **Installez les outils Go utiles :**
  ```bash
  go install golang.org/x/tools/cmd/goimports@latest
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  ```


### Node.js et npm

Requis uniquement pour les exemples d’intégration d’interfaces frontend.

```bash
# Install Node.js 20+ and npm
node --version  # Should be 20+
npm --version
```

### Dépendances propres à chaque plateforme

**macOS :**

```bash
# Install Xcode Command Line Tools
xcode-select --install

# Verify installation
xcode-select -p  # Should output a path
```

**Windows :**

1. Installez [MSYS2](https://www.msys2.org/) pour disposer d’un environnement de type Unix
2. WebView2 Runtime (préinstallé sur Windows 11 ; pour Windows 10, procédez au [téléchargement](https://developer.microsoft.com/en-us/microsoft-edge/webview2/))
3. Facultatif : installez [Git for Windows](https://git-scm.com/download/win)

**Linux (Debian/Ubuntu) :**

```bash
sudo apt update
# Default GTK4 + WebKitGTK 6.0 stack (Ubuntu 24.04+ / Debian 13+)
sudo apt install build-essential pkg-config libgtk-4-dev libwebkitgtk-6.0-dev
# For the legacy -tags gtk3 path:
sudo apt install libgtk-3-dev libwebkit2gtk-4.1-dev
```

**Linux (Fedora/RHEL) :**

```bash
# Default GTK4 stack
sudo dnf install gcc pkg-config gtk4-devel webkitgtk6.0-devel
# Legacy GTK3 path:
sudo dnf install gtk3-devel webkit2gtk4.1-devel
```

**Linux (Arch) :**

```bash
# Default GTK4 stack
sudo pacman -S base-devel gtk4 webkitgtk-6.0
# Legacy GTK3 path:
sudo pacman -S gtk3 webkit2gtk-4.1
```

## Configuration du dépôt

### Cloner et configurer

```bash
# Clone your fork
git clone --branch releases/v3.0.0 https://github.com/YOUR_USERNAME/wails.git
cd wails

# Add upstream remote
git remote add upstream https://github.com/wailsapp/wails.git

# Verify remotes
git remote -v
```

### Compiler la CLI Wails

```bash
# Navigate to v3 directory
cd v3

# Build the CLI
go build -o ../wails3 ./cmd/wails3

# Test the build
cd ..
./wails3 version
```

### Ajouter au PATH (facultatif)

**Linux/macOS :**

```bash
# Add to ~/.bashrc or ~/.zshrc
export PATH=$PATH:/path/to/wails
```

**Windows :**

Ajoutez le répertoire Wails à votre variable d’environnement PATH dans les propriétés système.

## Configuration de l’IDE

### VS Code (recommandé)

1. **Installez VS Code :** [Télécharger](https://code.visualstudio.com/)

2. **Installez les extensions :**
  - Go (par l’équipe Go de Google)
  - ESLint
  - Prettier
  - MDX (pour la documentation)


3. **Configurez les paramètres de l’espace de travail** (`.vscode/settings.json`) :
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

1. **Installez GoLand :** [Télécharger](https://www.jetbrains.com/go/)

2. **Configurez GoLand :**
  - Activez la prise en charge des modules Go
  - Configurez des observateurs de fichiers pour `goimports`
  - Configurez le style du code conformément aux conventions du projet


## Vérifier votre configuration

Exécutez les commandes suivantes pour vérifier que tout fonctionne :

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

Si l’application de test se compile et s’exécute, votre environnement est prêt !

## Exécuter les tests

### Tests unitaires

```bash
cd v3
go test ./...
cd ..
```

### Tests d’un paquet précis

```bash
cd v3
go test ./pkg/application
go test ./pkg/events -v  # Verbose output
```

### Exécuter avec la mesure de couverture

```bash
cd v3
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Exécuter avec le détecteur de situations de concurrence

```bash
cd v3
go test ./... -race
```

## Travailler sur la documentation

Les personnes qui contribuent à la documentation peuvent utiliser M-Press sans installer la CLI Wails ni les dépendances de compilation des applications de bureau. Consultez [Rédiger la documentation avec M-Press](/contributing/documentation/) pour connaître la procédure de téléchargement des exécutables aux versions épinglées, d’installation, d’aperçu en direct, de validation et de création des PR destinées aux branches de publication.

Les fichiers Markdown en anglais se trouvent dans `docs/mpress/content/` ; les traductions se trouvent dans les répertoires propres à chaque langue. Depuis la racine du dépôt, `mpress dev` met l’aperçu à disposition à l’adresse `http://127.0.0.1:3000`. Arrêtez-le avec Ctrl+C avant d’exécuter les commandes de validation du guide.

## Débogage

### Déboguer le code Go

**VS Code :**

Créez `.vscode/launch.json` :

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

**Ligne de commande :**

```bash
# Use Delve debugger
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug ./cmd/wails3 -- dev
```

### Déboguer le code propre à chaque plateforme

Le débogage du code propre à une plateforme nécessite les outils correspondants :

- **macOS :** Xcode Instruments
- **Windows :** Visual Studio Debugger
- **Linux :** GDB

## Problèmes courants

### « command not found: wails3 »

Ajoutez le répertoire Wails à votre PATH ou utilisez `./wails3` depuis la racine du projet.

### « webkitgtk-6.0 not found » ou « webkit2gtk not found » (Linux)

Installez les paquets de développement correspondant à la pile pour laquelle vous compilez :

```bash
# Default GTK4 stack (Debian/Ubuntu):
sudo apt install libwebkitgtk-6.0-dev

# Legacy GTK3 path:
sudo apt install libwebkit2gtk-4.1-dev
```

### Échec de la compilation en raison d’erreurs de modules Go

```bash
cd v3
go mod tidy
go mod download
```

### Erreurs « CGO_ENABLED » sous Windows

Vérifiez qu’un compilateur C (MinGW-w64 via MSYS2) est présent dans votre PATH.

## Étapes suivantes

- Consultez les [Normes de codage](/contributing/standards/)
- Explorez la [Documentation technique](/contributing/)
- Trouvez un ticket sur lequel travailler : [Tickets pour débuter](https://github.com/wailsapp/wails/labels/good%20first%20issue)
