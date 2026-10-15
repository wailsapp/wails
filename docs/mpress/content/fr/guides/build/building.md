---
title: "Création d’applications"
description: "Compilez et empaquetez votre application Wails"
slug: "guides/build/building"
sourcePath: "guides/build/building.md"
---

Wails v3 utilise [Task](https://taskfile.dev) comme système de build. Les commandes `wails3 build` et `wails3 package` sont des interfaces pratiques autour de Task.

## Compilation

Compilez pour la plateforme actuelle :

```bash
wails3 build
```

Compilez pour une plateforme spécifique :

```bash
wails3 build GOOS=windows
wails3 build GOOS=darwin
wails3 build GOOS=linux

# With architecture
wails3 build GOOS=darwin GOARCH=arm64

# Environment variable style works too
GOOS=windows wails3 build
```

La sortie est placée dans le répertoire `bin/`.

@note{type="tip"}
La compilation croisée vers macOS ou Linux depuis une autre plateforme nécessite Docker. Pour la configuration, consultez [Compilations multiplateformes](/guides/build/cross-platform/).

@end

## Développement

Exécutez votre application avec le rechargement à chaud :

```bash
wails3 dev
```

Cette commande lance un observateur de fichiers qui recompile et redémarre votre application à chaque modification. Par défaut, le serveur de développement du frontend s’exécute sur le port 9245.

```bash
# Custom port
wails3 dev -port 3000

# Enable HTTPS
wails3 dev -s
```

## Empaquetage

Empaquetez votre application pour la distribuer :

```bash
wails3 package
wails3 package GOOS=windows
wails3 package GOOS=darwin
wails3 package GOOS=linux
```

Cette commande crée des paquets propres à chaque plateforme :

- **Windows** : programme d’installation NSIS — consultez [Empaquetage pour Windows](/guides/build/windows/)
- **macOS** : paquet d’application (`.app`) — consultez [Empaquetage pour macOS](/guides/build/macos/)
- **Linux** : AppImage, deb et rpm — consultez [Empaquetage pour Linux](/guides/build/linux/)

## Balises de build personnalisées

Transmettez des balises de build Go personnalisées avec l’option `-tags` :

```bash
# Build with legacy GTK3 + WebKit2GTK 4.1 on Linux (default is GTK4 + WebKitGTK 6.0)
wails3 build -tags gtk3

# Build in server mode (no GUI, CGO-free)
wails3 build -tags server

# Combine multiple tags
wails3 build -tags gtk3,customtag
```

Les balises sont transmises au Taskfile sous la forme `EXTRA_TAGS`. Pour plus de détails, consultez [Build du serveur](/guides/server-build/) et [Empaquetage pour Linux — prise en charge de l’ancien GTK3](/guides/build/linux/#legacy-gtk3-support).

## Utilisation directe de Task

Pour davantage de contrôle, utilisez directement Task :

```bash
# List available tasks
wails3 task --list

# Verbose output
wails3 task build -v

# Dry run
wails3 task --dry

# Force rebuild
wails3 task build -f

# Pass variables
wails3 task darwin:build ARCH=amd64
```

Les tâches propres à une plateforme, telles que `linux:create:deb` ou `darwin:build:universal`, sont uniquement disponibles via Task.

## Génération des ressources

Régénérez les icônes ou mettez à jour la configuration du build :

```bash
wails3 generate icons -input build/appicon.png
wails3 update build-assets -name "MyApp" -config build/config.yml -dir build
```
