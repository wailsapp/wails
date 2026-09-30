---
title: "GitHub-Release-Assets für den Updater"
description: "Wie der Wails-Updater Anwendungsartefakte aus GitHub Releases auswählt und Installationspakete vermeidet."
slug: "guides/updater-github-release-assets"
sourcePath: "guides/updater-github-release-assets.md"
---

Der GitHub-Releases-Provider wählt ein Release-Asset anhand des konfigurierten `AssetMatcher` aus. Wenn `AssetMatcher` den Wert `nil` hat, verwendet er `github.DefaultAssetMatcher`.

## Standardmäßiger Abgleich

Der Standard-Matcher sucht in jedem Asset-Dateinamen nach der aktuellen Plattform und Architektur. Er erkennt gängige Architektur-Aliasse, darunter:

- `amd64`, `x86_64` und `x64`
- `arm64` und `aarch64`
- `386`, `i386`, `x86` und `ia32`

Begleitdateien wie Signaturen und Prüfsummen werden ignoriert.

## Installer-Assets

Ein GitHub Release kann sowohl die vom Updater verwendete Anwendungsbinärdatei als auch einen herkömmlichen Installer für die Erstinstallation enthalten. Der Standard-Matcher ignoriert Assets, deren in Kleinbuchstaben umgewandelter Dateiname:

- `-installer.` enthält
- `_installer.` enthält
- genau `installer.exe` lautet

Beispiel mit folgenden Windows-Assets:

```text
myapp-windows-amd64.exe
myapp-windows-amd64-installer.exe
```

`DefaultAssetMatcher` wählt `myapp-windows-amd64.exe` aus und ignoriert den Installer. Dadurch wird verhindert, dass der Updater die laufende Anwendung durch eine mit NSIS oder ähnlich paketierte ausführbare Installer-Datei ersetzt.

Die Prüfung ist bewusst eng gefasst. Anwendungsnamen, die lediglich das Wort `installer` enthalten, bleiben gültig, darunter:

```text
myinstaller.exe
installer-tool-windows-amd64.exe
myinstaller-windows-amd64.zip
```

## Benutzerdefinierte Benennungsschemata

Konfigurieren Sie `AssetMatcher`, wenn Ihre Release-Assets nicht der Benennungskonvention aus Plattform und Architektur entsprechen oder wenn Sie Installer anders filtern müssen:

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

Ein benutzerdefinierter Matcher ersetzt `DefaultAssetMatcher` vollständig und ist daher dafür verantwortlich, Signaturen, Prüfsummen, Installer und alle anderen Assets auszuschließen, die nicht als Anwendungsupdate installiert werden sollen.
