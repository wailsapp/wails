---
title: "Création de paquets macOS"
description: "Créez un paquet de votre application Wails pour la distribuer sous macOS"
slug: "guides/build/macos"
sourcePath: "guides/build/macos.md"
---

## API macOS privées

Wails v3 utilise par défaut les API macOS publiques. Pour activer des fonctionnalités nécessitant des API Apple non documentées, compilez votre application avec l’unique étiquette de compilation Go `private_mac_apis` :

```bash
wails3 build -tags private_mac_apis
EXTRA_TAGS=private_mac_apis wails3 dev
wails3 package GOOS=darwin EXTRA_TAGS=private_mac_apis
```

Pour une compilation Go directe, utilisez `go build -tags private_mac_apis .` (ou `-tags production,private_mac_apis` pour la production). Les applications existantes qui dépendent de comportements privés doivent ajouter cette étiquette afin de les conserver. Cette étiquette s’applique uniquement aux compilations d’applications de bureau pour macOS.

Consultez [API macOS privées](/guides/build/private-macos-apis/) pour obtenir la liste complète des fonctionnalités et des valeurs d’option concernées, les comportements de repli exacts des compilations publiques, les correspondances de styles Liquid Glass et les combinaisons de compilation de l’inspecteur. Sans cette étiquette, les opérations réservées aux API privées ne font rien ; l’API Go publique reste inchangée.

## Paquet d’application

Créez un paquet macOS `.app` standard pour votre application :

```bash
wails3 package GOOS=darwin
```

Cette opération crée `bin/<AppName>.app`, qui contient :

- Le binaire compilé dans `Contents/MacOS/`
- L’icône de l’application dans `Contents/Resources/` (provenant de `icons.icns` ou, s’il est présent, d’un catalogue de ressources `Assets.car`)
- `Info.plist` contenant les métadonnées de l’application

## Ressources du paquet

`Contents/Resources/` est l’emplacement standard des fichiers en lecture seule distribués avec une application macOS. Utilisez-le pour les modèles volumineux, les données initiales, les médias, les packs linguistiques ou les autres contenus qui devraient être ouverts à la demande plutôt que compilés dans l’exécutable Go avec `embed`.

Wails place déjà l’icône de l’application dans ce répertoire. Pour ajouter vos propres fichiers, placez-les dans un répertoire source tel que `build/resources/`, puis ajoutez une étape de copie à la tâche `create:app:bundle` dans `build/darwin/Taskfile.yml` :

```yaml
tasks:
  create:app:bundle:
    cmds:
      # Existing bundle creation commands...
      - |
          if [ -d build/resources ]; then
            cp -R build/resources/. "{{.BIN_DIR}}/{{.APP_NAME}}.app/Contents/Resources/"
          fi
```

Si vous utilisez la tâche `darwin:run` du Taskfile, ajoutez la commande équivalente à sa tâche `run`, avec `{{.BIN_DIR}}/{{.APP_NAME}}.dev.app/Contents/Resources/` comme destination.

### Lecture des ressources depuis Go

Importez le paquet propre à la plateforme macOS :

```go
import (
	"io/fs"

	"github.com/wailsapp/wails/v3/pkg/mac"
)
```

Pour les petits fichiers, utilisez `LoadResource` :

```go
func loadSplash() ([]byte, error) {
	return mac.LoadResource("images/splash.png")
}
```

Pour les fichiers plus volumineux, utilisez `ResourceFS`. Cette fonction renvoie un `io/fs.FS` dont la racine est `Contents/Resources`, ce qui permet au code appelant d’ouvrir et de lire une ressource en continu sans devoir d’abord la charger intégralement dans une tranche d’octets Go :

```go
func openCatalogue() (fs.File, error) {
	resources, err := mac.ResourceFS()
	if err != nil {
		return nil, err
	}

	return resources.Open("catalogue/defaults.json")
}
```

Les noms de ressources sont des chemins séparés par des barres obliques et relatifs à `Contents/Resources`. `ResourceFS` et `LoadResource` renvoient `mac.ErrNotInAppBundle`, sauf si l’exécutable s’exécute depuis `.app/Contents/MacOS`.

Traitez les ressources du paquet comme immuables. La modification de fichiers à l’intérieur d’une application signée invalide sa signature de code ; stockez plutôt les données téléchargées, générées ou modifiables par l’utilisateur dans le répertoire Application Support de ce dernier.

### Binaire universel

Compilez pour les Mac équipés d’Apple Silicon comme pour ceux équipés de processeurs Intel :

```bash
wails3 task darwin:package:universal
```

Cette opération crée un seul `.app` qui s’exécute nativement sur les deux architectures. Les binaires universels peuvent être compilés sur n’importe quelle plateforme : sous Linux et Windows, `wails3 tool lipo` est utilisé automatiquement.

## Personnalisation du paquet

Modifiez `build/darwin/Info.plist` pour personnaliser les éléments suivants :

- Identifiant du paquet (`CFBundleIdentifier`)
- Nom et version de l’application
- Version minimale de macOS
- Associations de fichiers
- Schémas d’URL

L’icône de l’application est générée à partir des ressources du répertoire `build/`. Utilisez la tâche `generate:icons` :

```bash
wails3 task common:generate:icons
```

Cette tâche utilise `build/appicon.png` pour produire `darwin/icons.icns` et `windows/icon.ico`. Sous macOS, vous pouvez également fournir `build/appicon.icon` (format Icon Composer) : la tâche transmet `-iconcomposerinput appicon.icon -macassetdir darwin`, ce qui produit `Assets.car` et `darwin/icons.icns` à partir du fichier `.icon` (cette étape est ignorée sur les plateformes autres que macOS). Lorsque `Assets.car` est présent, exécutez la tâche `update:build-assets` afin que `Info.plist` et `CFBundleIconName` soient mis à jour en conséquence :

```bash
wails3 task common:update:build-assets
```

Pour exécuter manuellement la commande relative aux icônes depuis le répertoire `build/` :

```bash
cd build
wails3 generate icons -input appicon.png -macfilename darwin/icons.icns -windowsfilename windows/icon.ico -iconcomposerinput appicon.icon -macassetdir darwin
```

## Signature du code

Signez votre application pour la distribuer :

```bash
# Using the wrapper (auto-detects platform)
wails3 sign GOOS=darwin

# Or using the task directly
wails3 task darwin:sign
```

Configurez la signature dans `build/darwin/Taskfile.yml` :

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  ENTITLEMENTS: "build/darwin/entitlements.plist"
```

### Notarisation

Apple exige la notarisation des applications distribuées en dehors du Mac App Store :

```bash
wails3 task darwin:sign:notarize
```

Commencez par enregistrer vos identifiants. Exécutez l’assistant interactif (`wails3 setup signing`) ou appelez directement `notarytool` :

```bash
xcrun notarytool store-credentials "my-notarize-profile" \
  --apple-id "you@email.com" \
  --team-id "TEAMID" \
  --password "app-specific-password"
```

Effectuez la configuration dans `build/darwin/Taskfile.yml` :

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
```

Pour plus de détails, consultez [Signature des applications](/guides/build/signing/).

## Programme d’installation DMG

Le modèle fourni avec Wails 3 propose `wails3 task darwin:package:dmg`. Il crée d’abord le paquet `.app`, puis génère un DMG mis en forme à l’aide de la bibliothèque DMG. Par défaut, le DMG utilise un arrière-plan en dégradé aux couleurs de Wails, avec l’emblème du dragon rouge et le logotype WAILS.

```bash
wails3 task darwin:package:dmg
```

La tâche `darwin:create:dmg` de plus bas niveau crée un DMG à partir d’un paquet `.app` existant et peut être configurée directement depuis le Taskfile :

```yaml
vars:
  # These are the template defaults; override them when needed.
  DMG_BACKGROUND: build/darwin/dmg-background.png
  DMG_VOLUME_ICON: build/darwin/icons.icns
  DMG_FILE_ICON: build/darwin/dmg-file-icon.icns
  DMG_WINDOW_WIDTH: 540
  DMG_WINDOW_HEIGHT: 380
  DMG_FILES: "Install.command=build/darwin/Install.command,README.txt=README.md"
```

### Disposition par défaut

Le DMG généré contient :

- Le paquet d’application à gauche
- Un lien `Applications` à droite
- Une fenêtre du Finder de 540 × 380 pixels
- Des icônes de 96 points avec leur libellé en dessous
- L’arrière-plan aux couleurs de Wails provenant de `build/darwin/dmg-background.png`

Les icônes de l’application et de `Applications` sont positionnées par rapport aux dimensions configurées de la fenêtre. Ainsi, lorsque vous modifiez `DMG_WINDOW_WIDTH` ou `DMG_WINDOW_HEIGHT`, l’espacement proportionnel de la disposition à deux icônes par défaut est conservé. Pour un résultat optimal, utilisez une image d’arrière-plan ayant les mêmes dimensions en pixels que la fenêtre du Finder.

### Remplacement des ressources du DMG

Les fichiers générés sous `build/darwin/` sont des ressources de projet ordinaires et peuvent être remplacés :

- `DMG_BACKGROUND` définit l’image affichée derrière le contenu de la fenêtre du Finder.
- `DMG_VOLUME_ICON` définit l’icône affichée pour le volume monté.
- `DMG_FILE_ICON` définit l’icône affichée dans le Finder pour le fichier `.dmg` obtenu.

L’icône du volume et celle du fichier DMG sont des ressources distinctes. Le remplacement de l’icône de l’application ne remplace automatiquement aucune des deux.

### Ajout de fichiers supplémentaires

Utilisez `DMG_FILES` pour inclure des scripts d’installation, des notes de version, des licences ou d’autres ressources aux côtés de l’application. La valeur est une liste de paires `name=path` séparées par des virgules :

```yaml
vars:
  DMG_FILES: "Install.command=build/darwin/Install.command,README.txt=README.md"
```

Le nom précédant `=` est le nom de fichier affiché dans le DMG. Le chemin suivant `=` désigne le fichier source dans le projet. Les espaces en début et en fin sont ignorés.

Chaque nom affiché doit être unique. Les fichiers supplémentaires ne peuvent pas remplacer les entrées déjà créées par l’outil de création du paquet, notamment le bundle de l’application ou l’entrée `Applications`. En cas de conflit de noms, la création du paquet échoue avec une erreur au lieu de produire un DMG défectueux.

@note{type="note"}
La création de DMG n’est prise en charge que sous macOS, car elle utilise les outils d’image disque et du Finder de macOS. Les bundles `.app` compilés de manière croisée peuvent être créés ailleurs, mais le DMG final doit être produit sur un Mac.

@end

## Dépannage

### « L’app est endommagée et ne peut pas être ouverte »

L’application n’est pas signée. Signez-la avec un certificat Developer ID, ou les utilisateurs peuvent contourner Gatekeeper :

```bash
xattr -cr /path/to/YourApp.app
```

### Échec de la notarisation

Problèmes courants :

- **Identifiants non valides** : exécutez de nouveau `xcrun notarytool store-credentials` (ou `wails3 setup signing`)
- **Environnement d’exécution renforcé requis** : vérifiez que les autorisations incluent `com.apple.security.cs.allow-unsigned-executable-memory` si nécessaire
- **Horodatage manquant** : le processus de signature devrait inclure automatiquement un horodatage

### L’application compilée de manière croisée ne s’exécute pas

Les fichiers binaires macOS compilés de manière croisée ne sont pas signés. Transférez-les sur un Mac et signez-les avant de les tester :

```bash
codesign --force --deep --sign - YourApp.app
```
