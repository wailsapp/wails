---
title: "Système de build"
description: "Comprendre comment Wails compile et empaquette votre application"
slug: "concepts/build-system"
sourcePath: "concepts/build-system.md"
---

## Système de build unifié

Wails fournit un **système de build unifié** qui compile le code Go, regroupe les ressources frontend, incorpore le tout dans un exécutable unique et gère les builds propres à chaque plateforme, le tout avec une seule commande.

```bash
wails3 build
```

**Sortie :** exécutable natif contenant tous les éléments incorporés.

## Vue d’ensemble du processus de build

**[Espace réservé au diagramme du processus de build]**

## Phases du build

### 1. Phase d’analyse

Wails analyse votre code Go pour comprendre vos services :

```go
type GreetService struct {
    prefix string
}

func (g *GreetService) Greet(name string) string {
    return g.prefix + name + "!"
}
```

**Informations extraites par Wails :**

- Nom du service : `GreetService`
- Nom de la méthode : `Greet`
- Types des paramètres : `string`
- Types de retour : `string`

**Utilisation :** génération des liaisons TypeScript

### 2. Phase de génération

#### Liaisons TypeScript

Wails génère des liaisons avec vérification des types :

```javascript
// Auto-generated: frontend/bindings/<full-go-import-path>/greetservice.js
// (TypeScript is also generated when you pass `-ts`. The shape below is the real
// runtime call format — numeric IDs via $Call.ByID, imported from /wails/runtime.js.)
import { Call as $Call } from "/wails/runtime.js";

export function Greet($0) {
    return $Call.ByID(1234567890, $0);
}
```

**Avantages :**

- Sécurité complète des types
- Autocomplétion dans l’IDE
- Erreurs détectées à la compilation
- Commentaires JSDoc

#### Build du frontend

Votre outil de regroupement frontend s’exécute (Vite, webpack, etc.) :

```bash
# Vite example
vite build --outDir dist
```

**Opérations effectuées :**

- Compilation du JavaScript/TypeScript
- Traitement et minification du CSS
- Optimisation des ressources
- Génération des cartes de sources (uniquement en développement)
- Sortie dans `frontend/dist/`

### 3. Phase de compilation

#### Compilation Go

Le code Go est compilé avec des optimisations :

```bash
go build -ldflags="-s -w" -o myapp.exe
```

**Options :**

- `-s` : supprime la table des symboles
- `-w` : supprime les informations de débogage DWARF
- Résultat : binaire plus petit (réduction d’environ 30 %)

**Selon la plateforme :**

- Windows : `.exe` avec icône incorporée
- macOS : structure de bundle `.app`
- Linux : binaire ELF

#### Incorporation des ressources

Les ressources frontend sont incorporées dans le binaire Go :

```go
//go:embed frontend/dist
var assets embed.FS
```

**Résultat :** un exécutable unique qui contient tout.

### 4. Sortie

**Un seul binaire natif :**

- Windows : `myapp.exe` (environ 15 Mo)
- macOS : `myapp.app` (environ 15 Mo)
- Linux : `myapp` (environ 15 Mo)

**Aucune dépendance** (à l’exception de la WebView du système).

## Développement et production

@tabs{sync-key="mode"}
[Développement (wails3 dev)]
**Optimisé pour la rapidité :**

```bash
wails3 dev
```

**Opérations effectuées :**

1. Démarre le serveur de développement frontend (Vite sur le port 9245 par défaut)
2. Compile Go sans optimisations
3. Lance l’application en la faisant pointer vers le serveur de développement
4. Active le rechargement à chaud
5. Inclut les cartes de sources

**Caractéristiques :**

- **Rebuilds rapides** (&lt;1 s pour les modifications du frontend)
- **Aucune incorporation des ressources** (elles sont servies par le serveur de développement)
- **Symboles de débogage** inclus
- **Cartes de sources** activées
- **Journalisation détaillée**

**Taille du fichier :** plus grande (environ 50 Mo avec les symboles de débogage)

[Production (wails3 build)]
**Optimisé pour réduire la taille et améliorer les performances :**

```bash
wails3 build
```

**Opérations effectuées :**

1. Compile le frontend pour la production (avec minification)
2. Compile Go avec des optimisations
3. Supprime les symboles de débogage
4. Incorpore les ressources
5. Crée un binaire unique

**Caractéristiques :**

- **Code optimisé** (minifié, avec élimination du code inutilisé)
- **Ressources incorporées** (aucun fichier externe)
- **Symboles de débogage supprimés**
- **Aucune carte de sources**
- **Journalisation minimale**

**Taille du fichier :** plus petite (environ 15 Mo)

@end

## Commandes de build

### Build de base

```bash
wails3 build
```

**Sortie :** `bin/<APP_NAME>` (ou `bin/<APP_NAME>.exe` sous Windows). Le répertoire `bin/` se trouve à la racine du projet.

`wails3 build` est une fine surcouche autour de `wails3 task build`. La seule option de build qu’elle transmet est `--tags`, qui devient la variable Taskfile `EXTRA_TAGS` :

```bash
# Build with extra Go build tags
wails3 build --tags "myfeature,gtk4"
```

`wails3 build` ne propose aucune option `-platform`, `-o`, `-skipbindings`, `-clean`, `-debug`, `-devbuild`, `-icon`, `-ldflags` ou `-package`. La compilation croisée, les chemins de sortie, les icônes et l’empaquetage sont contrôlés par le Taskfile du projet (`Taskfile.yml` + `build/config.yml`).

### Builds multiplateformes et propres à chaque plateforme

Les builds par plateforme sont proposés sous forme de tâches Taskfile dans les espaces de noms `darwin:` / `windows:` / `linux:` (définis dans `build/Taskfile.<platform>.yml`). Par exemple :

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

Pour afficher toutes les tâches disponibles dans le projet actuel :

```bash
wails3 task --list
```

### Icônes et création de paquets

Générez les icônes propres à chaque plateforme (`build/icons.icns`, `build/icon.ico`, etc.) à partir d’un fichier PNG source :

```bash
wails3 generate icons -input appicon.png
```

Créez les programmes d’installation ou paquets propres à chaque plateforme :

```bash
wails3 package           # uses the current Go build env
wails3 task linux:create:deb
wails3 task windows:package
wails3 task darwin:package:universal
```

## Configuration du build

### Taskfile.yml

Les projets Wails 3 utilisent [Taskfile](https://taskfile.dev/) pour orchestrer leur build. Le fichier `Taskfile.yml` racine inclut les fichiers de tâches propres à chaque plateforme depuis `build/` :

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

Exécutez les tâches avec `wails3 task <name>` ou `task <name>` :

```bash
wails3 task windows:build
wails3 task darwin:package:universal
wails3 task linux:create:appimage
```

### Configuration du projet : `build/config.yml`

Les métadonnées du projet (nom, identifiant, version, valeurs d’Info.plist, paramètres NSIS, champs `.desktop`, protocoles personnalisés, etc.) se trouvent dans `build/config.yml`. Le Taskfile lit ce fichier lors de la génération des icônes, manifestes, programmes d’installation et autres éléments similaires. Il n’existe **aucun** fichier `build/build.json` dans Wails 3.

```yaml
# build/config.yml (illustrative)
info:
  productName: "My App"
  productIdentifier: "com.example.myapp"
  productVersion: "1.0.0"
  companyName: "Example Ltd."
  productDescription: "An application built with Wails"
```

Exécutez `wails3 generate build-assets` (ou `wails3 update build-assets`) pour actualiser, à partir de cette configuration, les ressources de build propres à chaque plateforme.

## Incorporation des ressources

### Fonctionnement

Wails utilise le paquet `embed` de Go :

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

**Lors du build :**

1. Le frontend est généré dans `frontend/dist/`
2. La directive `//go:embed` inclut les fichiers
3. Les fichiers sont compilés dans le binaire
4. Le binaire contient tout

**Lors de l’exécution :**

1. L’application démarre
2. Les ressources sont servies depuis la mémoire
3. Aucune E/S disque pour les ressources
4. Chargement rapide

### Ressources personnalisées

Incorporez des fichiers supplémentaires :

```go
//go:embed frontend/dist
var frontendAssets embed.FS

//go:embed data/*.json
var dataAssets embed.FS

//go:embed templates/*.html
var templateAssets embed.FS
```

## Optimisations du build

### Optimisations du frontend

**Vite (par défaut) :**

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

**Résultats :**

- JavaScript minifié (réduction d’environ 70 %)
- CSS minifié (réduction d’environ 60 %)
- Images optimisées
- Élimination du code mort appliquée

### Optimisations de Go

**Options du compilateur :**

```bash
-ldflags="-s -w"
```

- `-s` : supprime la table des symboles (réduction d’environ 10 %)
- `-w` : supprime les informations de débogage DWARF (réduction d’environ 20 %)

**Optimisations supplémentaires :**

```bash
-ldflags="-s -w -X main.version=1.0.0"
```

- `-X` : définit les valeurs des variables lors du build
- Utile pour les numéros de version et les dates de build

### Compression du binaire

**UPX (facultatif) :**

```bash
# After building
upx --best bin/myapp.exe
```

**Résultats :**

- Réduction de la taille d’environ 50 %
- Démarrage légèrement plus lent (environ 100 ms)
- Déconseillé sous macOS (problèmes de signature du code)

## Builds propres à chaque plateforme

### Windows

**Sortie :** `myapp.exe`

**Comprend :**

- Icône de l’application
- Informations de version
- Manifeste (paramètres UAC)

**Icône :**

```bash
# Generate platform icons from a source PNG
wails3 generate icons -input appicon.png -windowsfilename build/icon.ico
```

L’étape Windows `tool package` incorpore ensuite le fichier `.ico` généré dans l’exécutable.

**Manifeste :**

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

**Sortie :** `myapp.app` (paquet d’application)

**Structure :**

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

**Info.plist :**

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

**Binaire universel :**

Le Taskfile de macOS fournit une tâche `darwin:build:universal` (ainsi qu’une tâche `darwin:package:universal`) qui génère les deux architectures et les combine avec `wails3 tool lipo` :

```bash
wails3 task darwin:build:universal
```

### Linux

**Sortie :** `myapp` (binaire ELF)

**Dépendances :**

- GTK3
- WebKitGTK

**Fichier desktop :**

```ini
# myapp.desktop
[Desktop Entry]
Name=My App
Exec=/usr/bin/myapp
Icon=myapp
Type=Application
Categories=Utility;
```

**Installation :**

```bash
# Copy binary
sudo cp myapp /usr/bin/

# Copy desktop file
sudo cp myapp.desktop /usr/share/applications/

# Copy icon
sudo cp icon.png /usr/share/icons/hicolor/256x256/apps/myapp.png
```

## Performances du build

### Durées de build habituelles

| Phase | Durée | Remarques |
| --- | --- | --- |
| Analyse | &lt;1 s | Analyse du code Go |
| Génération des liaisons | &lt;1 s | Génération du TypeScript |
| Build du frontend | 5-30 s | Dépend de la taille du projet |
| Compilation Go | 2-10 s | Dépend de la taille du code |
| Incorporation des ressources | &lt;1 s | Incorporation du frontend |
| **Total** | **10-45 s** | Premier build |
| **Incrémentiel** | **5-15 s** | Builds suivants |

### Accélérer les builds

**1. Utilisez le cache de build :**

```bash
# Go build cache is automatic
# Frontend cache (Vite)
npm run build  # Uses cache by default
```

**2. Exécutez uniquement ce dont vous avez besoin :**

```bash
# Pick the specific Taskfile target you actually need
wails3 task common:build:frontend   # rebuild only the frontend
wails3 task windows:build           # rebuild only the Windows binary
```

**3. Parallélisez les builds (plusieurs machines/CI) :**

La compilation croisée entre Linux, Windows et macOS dans v3 s’effectue généralement dans un conteneur Docker `wails-cross` ou sur des runners dédiés à chaque plateforme ; `wails3 build` lui-même cible le système d’exploitation hôte. Consultez [Builds multiplateformes](/guides/build/cross-platform/) pour connaître les workflows pris en charge.

**4. Utilisez des outils plus rapides :**

```bash
# Use esbuild instead of webpack
# (Vite uses esbuild by default)
```

## Dépannage

### Échec du build

**Symptôme :** `wails3 build` se termine avec une erreur

**Causes courantes :**

1. **Erreur de compilation Go**
  ```bash
  # Check Go code compiles
  go build
  ```


2. **Erreur de build du frontend**
  ```bash
  # Check frontend builds
  cd frontend
  npm run build
  ```


3. **Dépendances manquantes**
  ```bash
  # Install dependencies
  npm install
  go mod download
  ```


### Binaire trop volumineux

**Symptôme :** le binaire dépasse 50 Mo

**Solutions :**

1. **Supprimez les symboles de débogage** (le Taskfile fourni transmet déjà `-ldflags="-s -w"` à `go build`).

2. **Vérifiez les ressources incorporées**
  ```bash
  # Remove unnecessary files from frontend/dist/
  # Check for large images, videos, etc.
  ```


3. **Utilisez la compression UPX**
  ```bash
  upx --best bin/myapp.exe
  ```


### Builds lents

**Symptôme :** les builds prennent plus de 1 minute

**Solutions :**

1. **Utilisez le cache de build**
  - Le cache Go est automatique
  - Le cache du frontend (Vite) est automatique


2. **Exécutez uniquement la tâche dont vous avez besoin**
  ```bash
  wails3 task common:build:frontend
  wails3 task windows:build
  ```


3. **Optimisez le build du frontend**
  ```javascript
  // vite.config.js
  export default {
    build: {
      minify: 'esbuild',  // Faster than terser
    },
  }
  ```


## Bonnes pratiques

### ✅ À faire

- **Utilisez `wails3 dev` pendant le développement** — Itérations rapides
- **Utilisez `wails3 build` pour les versions publiées** — Sortie optimisée
- **Versionnez vos builds** — Utilisez `-ldflags` pour incorporer la version
- **Testez les builds sur les plateformes cibles** — La compilation croisée n’est pas parfaite
- **Veillez à la rapidité des builds du frontend** — Optimisez la configuration du bundler
- **Utilisez le cache de build** — Accélère les builds suivants

### ❌ À ne pas faire

- **Ne validez pas le répertoire `build/` dans le dépôt** — Ajoutez-le à `.gitignore`
- **Ne négligez pas les tests des builds** — Testez toujours avant la publication
- **N’incorporez pas de ressources inutiles** — Limitez la taille des binaires
- **N’utilisez pas de builds de débogage en production** — Utilisez des builds optimisés
- **N’oubliez pas la signature du code** — Elle est requise pour la distribution

## Étapes suivantes

**Créer des applications** — Guide détaillé sur la création et l’empaquetage [En savoir plus →](/guides/build/building/)

**Builds multiplateformes** — Créez des builds pour toutes les plateformes depuis une seule machine [En savoir plus →](/guides/build/cross-platform/)

**Créer des programmes d’installation** — Créez des programmes d’installation pour les utilisateurs finaux [En savoir plus →](/guides/installers/)

---

**Des questions sur le processus de build ?** Posez-les sur [Discord](https://discord.gg/JDdSxwjhGf) ou consultez les [exemples de build](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/build).
