---
title: "Ressources des versions GitHub pour le programme de mise à jour"
description: "Comment le programme de mise à jour de Wails sélectionne les artefacts d’application dans les versions GitHub et évite les programmes d’installation."
slug: "guides/updater-github-release-assets"
sourcePath: "guides/updater-github-release-assets.md"
---

Le fournisseur GitHub Releases sélectionne une ressource de version à l’aide du `AssetMatcher` configuré. Lorsque `AssetMatcher` vaut `nil`, il utilise `github.DefaultAssetMatcher`.

## Correspondance par défaut

Le mécanisme de correspondance par défaut recherche la plateforme et l’architecture actuelles dans le nom de fichier de chaque ressource. Il reconnaît les alias d’architecture courants, notamment :

- `amd64`, `x86_64` et `x64`
- `arm64` et `aarch64`
- `386`, `i386`, `x86` et `ia32`

Les fichiers annexes, tels que les signatures et les sommes de contrôle, sont ignorés.

## Ressources d’installation

Une version GitHub peut contenir à la fois le fichier binaire de l’application utilisé par le programme de mise à jour et un programme d’installation classique destiné à une première installation. Le mécanisme de correspondance par défaut ignore les ressources dont le nom de fichier en minuscules :

- contient `-installer.`
- contient `_installer.`
- est exactement `installer.exe`

Par exemple, avec les ressources Windows suivantes :

```text
myapp-windows-amd64.exe
myapp-windows-amd64-installer.exe
```

`DefaultAssetMatcher` sélectionne `myapp-windows-amd64.exe` et ignore le programme d’installation. Cela empêche le programme de mise à jour de remplacer l’application en cours d’exécution par un exécutable d’installation NSIS ou conditionné de façon similaire.

La vérification est volontairement ciblée. Les noms d’application qui contiennent simplement le mot `installer` restent valides, notamment :

```text
myinstaller.exe
installer-tool-windows-amd64.exe
myinstaller-windows-amd64.zip
```

## Schémas de nommage personnalisés

Configurez `AssetMatcher` lorsque vos ressources de version ne respectent pas la convention de nommage fondée sur la plateforme et l’architecture, ou lorsque vous avez besoin d’un filtrage différent des programmes d’installation :

```go
import (
    "strings"

    "github.com/wailsapp/wails/v3/pkg/updater"
    "github.com/wailsapp/wails/v3/pkg/updater/providers/github"
)

gh, err := github.New(github.Config{
    Repository: "myorg/myapp",
    AssetMatcher: func(req updater.CheckRequest, assets []github.ReleaseAsset) int {
        for i, asset := range assets {
            name := strings.ToLower(asset.Name)
            if strings.Contains(name, req.Platform) &&
                strings.Contains(name, req.Arch) &&
                !strings.Contains(name, "-setup.") {
                return i
            }
        }
        return -1
    },
})
```

Un mécanisme de correspondance personnalisé remplace entièrement `DefaultAssetMatcher`. Il lui incombe donc d’exclure les signatures, les sommes de contrôle, les programmes d’installation et toute autre ressource qui ne doit pas être installée en tant que mise à jour de l’application.
