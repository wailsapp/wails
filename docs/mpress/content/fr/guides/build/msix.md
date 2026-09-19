---
title: "Création de packages MSIX"
description: "Création d’un package MSIX pour votre application Wails v3"
slug: "guides/build/msix"
sourcePath: "guides/build/msix.md"
---

MSIX est le format moderne de création de packages d’applications Windows. Wails peut générer un package MSIX dans le cadre de votre build Windows.

Les instructions relatives à la création de packages MSIX sont disponibles dans le guide [Création de packages Windows](/guides/build/windows/#msix-package).

## Création de packages avec la CLI

Exécutez l’outil MSIX sous Windows, y compris en CI. Le moteur par défaut utilise `MakeAppx.exe` ; la signature nécessite également `signtool.exe`. Tous deux font partie du SDK Windows. L’assistant d’installation ouvre le Microsoft Store et, si nécessaire, la page de téléchargement du SDK ; terminez l’installation avant de créer le package.

Définissez l’identité dans `build/config.yml`, puis compilez votre exécutable et créez son package :

```yaml
info:
  companyName: "Example Corp"
  productName: "MyApp"
  productIdentifier: "com.example.myapp"
  description: "MyApp"
  version: "1.0.0"
```

```powershell
wails3 tool msix-install-tools
wails3 build GOOS=windows
wails3 tool msix --executable bin/myapp.exe --name myapp.exe
```

La commande directe crée `MyApp.msix` dans le répertoire courant. La [tâche de création de packages Windows](/guides/build/windows/#msix-package) fournit son propre chemin de sortie.

## Options de la CLI

| Option | Signification |
| --- | --- |
| `--config` | Fichier de configuration ; valeur par défaut : `build/config.yml`. |
| `--executable`, `--name` | Exécutable existant et son nom de fichier dans le package ; les deux sont obligatoires. |
| `--out` | Fichier de sortie ; valeur par défaut : `<ProductName>.msix`. |
| `--arch` | Architecture du package : `x64` (par défaut), `x86`, `arm`, `arm64`, `x86a64` ou `neutral`. Les alias Go `amd64` et `386` sont acceptés. Utilisez l’architecture de l’exécutable. |
| `--publisher` | Identité de l’éditeur ; valeur par défaut : `CN=<companyName>`. |
| `--cert`, `--cert-password` | Chemin du certificat PFX et mot de passe pour la signature. |
| `--use-makeappx` | Utiliser l’outil de création de packages par défaut du SDK Windows. |
| `--use-msix-tool` | Choisir explicitement `MsixPackagingTool.exe`, qui doit figurer dans `PATH`. |

## Signature et CI

Pour une distribution hors du Store, signez avec un certificat approuvé sur la machine cible. Le champ Subject du certificat doit correspondre exactement à `--publisher`. Le moteur MakeAppx appelle SignTool avec SHA256 lorsque `--cert` est fourni. Consultez le [guide de signature de Microsoft](https://learn.microsoft.com/en-us/windows/msix/package/sign-msix-package-guide).

Cette étape de workflow Windows suppose que Wails et le SDK sont installés et qu’une étape précédente a déposé de façon sécurisée le fichier PFX à l’emplacement `CERT_PATH`. Un secret contenant un chemin ne transfère pas le certificat à lui seul :

```yaml
- name: MSIX
  if: runner.os == 'Windows'
  shell: pwsh
  run: |
    wails3 build GOOS=windows
    wails3 tool msix --executable bin/myapp.exe --name myapp.exe --publisher "$env:MSIX_PUBLISHER" --cert "$env:CERT_PATH" --cert-password "$env:CERT_PASSWORD"
  env:
    MSIX_PUBLISHER: ${{ vars.MSIX_PUBLISHER }}
    CERT_PATH: ${{ secrets.WINDOWS_CERT_PATH }}
    CERT_PASSWORD: ${{ secrets.WINDOWS_CERT_PASSWORD }}
```

## Associations de fichiers et ressources

Ajoutez les extensions sans point initial dans `build/config.yml` ; le manifeste généré ajoute le point :

```yaml
fileAssociations:
  - ext: myext
    name: MyApp Document
    description: MyApp Document
    iconName: fileicon
```

Gérez l’ouverture des fichiers à l’exécution comme décrit dans [Associations de fichiers](/guides/file-associations/). Le moteur MakeAppx copie actuellement uniquement l’exécutable et génère des images transparentes de remplacement. Il n’importe pas les fichiers `Assets/` du projet et ne convertit pas les icônes `iconName`. Utilisez un processus de création de packages personnalisé pour vos ressources graphiques ou des DLL supplémentaires.

| Ressource générée | Dimensions (pixels) |
| --- | --- |
| `Square150x150Logo.png` | 150×150 |
| `Square44x44Logo.png` | 44×44 |
| `Wide310x150Logo.png` | 310×150 |
| `StoreLogo.png` | 50×50 |
| `SplashScreen.png` | 620×300 |
| `FileIcon.png` | 44×44 |

`FileIcon.png` est généré uniquement lorsque des associations de fichiers sont configurées. Ces fichiers se trouvent dans le répertoire `Assets/` du package.

## Soumission au Store et dépannage

Réservez votre application dans le [portail Partner Center](https://partner.microsoft.com/dashboard) et utilisez l’identité de package et l’éditeur qui y figurent pour préparer la soumission. Le Store signe les packages MSIX lors de la soumission ; vous n’avez pas besoin d’acheter un certificat de signature pour cette voie de distribution. Consultez les [exigences de Microsoft pour les packages](https://learn.microsoft.com/en-us/windows/apps/publish/publish-your-app/msix/app-package-requirements).

Si `MakeAppx.exe` ou `signtool.exe` ne peut pas être trouvé, installez ou réparez le SDK Windows. Wails recherche dans `PATH` et aux emplacements standard du SDK. En cas d’échec de signature, vérifiez le champ Subject du certificat, sa validité et son approbation sur la machine cible ; consultez le [dépannage MSIX](https://learn.microsoft.com/en-us/windows/msix/msix-troubleshooting-guide).
