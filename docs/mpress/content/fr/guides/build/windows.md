---
title: "Création de paquets Windows"
description: "Créez un paquet de votre application Wails pour la distribuer sous Windows"
slug: "guides/build/windows"
sourcePath: "guides/build/windows.md"
---

## Programme d’installation NSIS

Le format de paquet par défaut crée un programme d’installation NSIS :

```bash
wails3 package GOOS=windows
```

Cette opération exécute `wails3 task windows:package`, qui :

1. Compile l’application
2. Génère le programme d’amorçage de WebView2
3. Crée un programme d’installation NSIS

Sortie : `build/windows/nsis/<AppName>-installer.exe`

### Paquet MSIX

Pour une distribution via le Microsoft Store ou un déploiement Windows moderne :

```bash
wails3 package GOOS=windows FORMAT=msix
```

Sortie : `bin/<AppName>-<arch>.msix`

@note{type="note"}
MSIX nécessite soit `makeappx.exe` (Windows SDK), soit les outils MSIX autonomes. Le Taskfile Windows expose la tâche d’installation sous le nom `wails3 task install:msix:tools`.

@end

## Personnalisation du programme d’installation

La configuration de NSIS se trouve dans `build/windows/nsis/project.nsi`. Modifiez ce fichier pour personnaliser :

- L’interface utilisateur et l’identité visuelle du programme d’installation
- Le répertoire d’installation
- Les raccourcis du menu Démarrer et du bureau
- Les associations de fichiers
- Le contrat de licence

Les métadonnées de l’application proviennent de `build/windows/info.json` :

```json
{
  "fixed": {
    "file_version": "1.0.0"
  },
  "info": {
    "0000": {
      "ProductVersion": "1.0.0",
      "CompanyName": "My Company",
      "FileDescription": "My Application",
      "ProductName": "MyApp"
    }
  }
}
```

## Signature du code

Signez votre exécutable et votre programme d’installation pour éviter les avertissements SmartScreen :

```bash
# Using the wrapper (auto-detects platform)
wails3 sign GOOS=windows

# Or using tasks directly
wails3 task windows:sign
wails3 task windows:sign:installer
```

Configurez la signature dans `build/windows/Taskfile.yml` :

```yaml
vars:
  SIGN_CERTIFICATE: "path/to/certificate.pfx"
  # Or use thumbprint for certificates in Windows store
  SIGN_THUMBPRINT: "certificate-thumbprint"
  TIMESTAMP_SERVER: "http://timestamp.digicert.com"
```

Stockez le mot de passe de votre certificat de manière sécurisée :

```bash
wails3 setup signing
```

Pour plus de détails, consultez [Signature des applications](/guides/build/signing/).

## Compilation pour ARM

```bash
wails3 build GOOS=windows GOARCH=arm64
wails3 package GOOS=windows GOARCH=arm64
```

## Résolution des problèmes

### makensis introuvable

Installez NSIS :

```bash
# Windows
winget install NSIS.NSIS

# Or download from https://nsis.sourceforge.io/
```

### Avertissement SmartScreen

Votre exécutable n’est pas signé. Consultez la section [Signature du code](#signature-du-code) ci-dessus.

### WebView2 manquant

Le programme d’installation comprend un programme d’amorçage de WebView2 qui télécharge l’environnement d’exécution si nécessaire. Si vous devez effectuer une installation hors ligne, téléchargez Evergreen Standalone Installer auprès de Microsoft.
