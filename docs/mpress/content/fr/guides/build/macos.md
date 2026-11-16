---
title: "Création de paquets macOS"
description: "Créez un paquet de votre application Wails pour la distribuer sous macOS"
slug: "guides/build/macos"
sourcePath: "guides/build/macos.md"
---

## API macOS privées

Par défaut, Wails v3 utilise les API macOS publiques. Pour activer les fonctionnalités qui nécessitent des API Apple non documentées, compilez votre application avec l’unique balise de compilation Go `private_mac_apis` :

```bash
wails3 build -tags private_mac_apis
EXTRA_TAGS=private_mac_apis wails3 dev
wails3 package GOOS=darwin EXTRA_TAGS=private_mac_apis
```

Pour une compilation Go directe, utilisez `go build -tags private_mac_apis .` (ou `-tags production,private_mac_apis` en production). Les applications existantes qui dépendent d’un comportement privé doivent ajouter cette balise pour le conserver. Cette balise s’applique uniquement aux compilations d’applications de bureau pour macOS.

Consultez [API macOS privées](/guides/build/private-macos-apis/) pour obtenir la liste complète des fonctionnalités et valeurs d’option concernées, les comportements de repli exacts des compilations publiques, les correspondances de styles Liquid Glass et les combinaisons de compilation de l’inspecteur. Sans cette balise, les opérations exclusivement privées ne produisent aucun effet ; l’API Go publique reste inchangée.

## Paquet d’application

Créez un paquet `.app` macOS standard pour votre application :

```bash
wails3 package GOOS=darwin
```

Cette commande crée `bin/<AppName>.app`, qui contient :

- Le binaire compilé dans `Contents/MacOS/`
- L’icône de l’application dans `Contents/Resources/` (provenant de `icons.icns` ou, s’il est présent, d’un catalogue de ressources `Assets.car`)
- `Info.plist` avec les métadonnées de l’application

### Projets dans des dossiers synchronisés

Certains fournisseurs de fichiers macOS ajoutent des métadonnées Finder aux répertoires `.app`. `codesign` rejette alors le paquet, ou sa signature devient non valide après la signature. Supprimer uniquement les métadonnées ne suffit pas si le fournisseur les ajoute de nouveau.

Les tâches de compilation macOS générées recherchent ce comportement. Lorsqu’il est détecté, Wails déplace le répertoire de sortie généré du projet vers `~/Library/Application Support/Wails/build-output/<project>-<hash>` et remplace le répertoire de sortie d’origine par un lien symbolique. Les fichiers de sortie existants sont conservés, et `bin/<AppName>.app` reste le chemin utilisé pour la création du paquet, la signature et le mode de développement. L’interface en ligne de commande affiche l’emplacement de stockage local lorsqu’elle effectue le déplacement.

Définissez `WAILS_MACOS_OUTPUT_ROOT` pour choisir un autre emplacement de stockage local non synchronisé. La racine de relocalisation doit se trouver sur le même système de fichiers que le projet. Une valeur `BIN_DIR` explicite située hors du projet est utilisée telle qu’elle est configurée. Pour bénéficier de cette gestion, les projets existants doivent utiliser le fichier `build/darwin/Taskfile.yml` mis à jour.

Pour la distribution, copiez `bin/<AppName>.app` ou créez un DMG. Copier uniquement le lien symbolique `bin` n’inclut pas les artefacts de compilation. La commande directe `wails3 tool sign` signale une erreur si son dossier d’entrée reçoit de nouveau des métadonnées ; déplacez ce paquet vers un dossier non synchronisé avant de le signer.

## Ressources du paquet

`Contents/Resources/` est l’emplacement standard des fichiers en lecture seule distribués avec une application macOS. Utilisez-le pour les modèles volumineux, les données initiales, les contenus multimédias, les modules linguistiques ou les autres ressources qui doivent être ouvertes à la demande plutôt que compilées dans l’exécutable Go avec `embed`.

Wails place déjà l’icône de l’application dans ce répertoire. Pour ajouter vos propres fichiers, placez-les dans un répertoire source tel que `build/resources/`, puis ajoutez une étape de copie à la tâche `create:app:bundle` dans `build/darwin/Taskfile.yml` :

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

Importez le paquet de la plateforme macOS :

```go
import (
	"io/fs"

	"github.com/wailsapp/wails/v3/pkg/mac"
)
```

Pour les petits fichiers, utilisez `LoadResource` :

```go
func loadSplash() ([]byte, error) {
	return mac.LoadResource("images/splash.png")
}
```

Pour les fichiers plus volumineux, utilisez `ResourceFS`. Cette fonction renvoie un `io/fs.FS` dont la racine est `Contents/Resources`, afin que le code appelant puisse ouvrir et lire une ressource en flux sans devoir d’abord la charger entièrement dans une tranche d’octets Go :

```go
func openCatalogue() (fs.File, error) {
	resources, err := mac.ResourceFS()
	if err != nil {
		return nil, err
	}

	return resources.Open("catalogue/defaults.json")
}
```

Les noms de ressources sont des chemins séparés par des barres obliques et relatifs à `Contents/Resources`. `ResourceFS` et `LoadResource` renvoient `mac.ErrNotInAppBundle` sauf si l’exécutable est lancé depuis `.app/Contents/MacOS`.

Traitez les ressources du paquet comme immuables. La modification de fichiers à l’intérieur d’une application signée invalide sa signature de code ; stockez plutôt les données téléchargées, générées ou modifiables par l’utilisateur dans le répertoire Application Support de ce dernier.

### Binaire universel

Compilez l’application pour les Mac Apple Silicon et Intel :

```bash
wails3 task darwin:package:universal
```

Cette commande crée un seul `.app` qui s’exécute nativement sur les deux architectures. Les binaires universels peuvent être compilés sur n’importe quelle plateforme — sous Linux et Windows, `wails3 tool lipo` est utilisé automatiquement.

## Personnalisation du paquet

Modifiez `build/darwin/Info.plist` pour personnaliser :

- L’identifiant du paquet (`CFBundleIdentifier`)
- Le nom et la version de l’application
- La version minimale de macOS
- Les associations de fichiers
- Les schémas d’URL

L’icône de l’application est générée à partir des ressources du répertoire `build/`. Utilisez la tâche `generate:icons` :

```bash
wails3 task common:generate:icons
```

Cette tâche utilise `build/appicon.png` pour produire `darwin/icons.icns` et `windows/icon.ico`. Sous macOS, vous pouvez également fournir `build/appicon.icon` (format Icon Composer) : la tâche transmet `-iconcomposerinput appicon.icon -macassetdir darwin`, qui produit `Assets.car` et `darwin/icons.icns` à partir du fichier `.icon` (cette étape est ignorée sur les plateformes autres que macOS). Lorsque `Assets.car` est présent, exécutez la tâche `update:build-assets` afin de mettre à jour `Info.plist` et `CFBundleIconName` en conséquence :

```bash
wails3 task common:update:build-assets
```

Pour exécuter manuellement la commande d’icône depuis le répertoire `build/` :

```bash
cd build
wails3 generate icons -input appicon.png -macfilename darwin/icons.icns -windowsfilename windows/icon.ico -iconcomposerinput appicon.icon -macassetdir darwin
```

## Signature de code

Signez votre application pour la distribuer :

```bash
# Using the wrapper (auto-detects platform)
wails3 sign GOOS=darwin

# Or using the task directly
wails3 task darwin:sign
```

Configurez la signature dans `build/darwin/Taskfile.yml` :

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  ENTITLEMENTS: "build/darwin/entitlements.plist"
```

### Notarisation

Apple exige la notarisation des applications distribuées hors du Mac App Store :

```bash
wails3 task darwin:sign:notarize
```

Commencez par enregistrer vos identifiants. Exécutez l’assistant interactif (`wails3 signing`) ou appelez directement `notarytool` :

```bash
xcrun notarytool store-credentials "my-notarize-profile" \
  --apple-id "you@email.com" \
  --team-id "TEAMID" \
  --password "app-specific-password"
```

Effectuez la configuration dans `build/darwin/Taskfile.yml` :

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
```

Consultez [Signature des applications](/guides/build/signing/) pour plus de détails.

## Programme d’installation DMG

Le modèle 3 fourni avec Wails propose `wails3 task darwin:package:dmg`. Il crée d’abord `.app`, puis génère un DMG stylisé avec la bibliothèque DMG. Par défaut, le DMG utilise un arrière-plan en dégradé aux couleurs de Wails, avec l’emblème du dragon rouge et le logotype WAILS.

```bash
wails3 task darwin:package:dmg
```

La tâche de plus bas niveau `darwin:create:dmg` crée un DMG à partir d’un paquet `.app` existant et peut être configurée directement dans le Taskfile :

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

Le DMG généré contient :

- Le paquet de l’application à gauche
- Un lien `Applications` à droite
- Une fenêtre Finder de 540 × 380 pixels
- Des icônes d’une taille de 96 points, avec leur libellé en dessous
- L’arrière-plan aux couleurs de Wails provenant de `build/darwin/dmg-background.png`

Les icônes de l’application et de `Applications` sont positionnées par rapport aux dimensions configurées de la fenêtre. Ainsi, la modification de `DMG_WINDOW_WIDTH` ou de `DMG_WINDOW_HEIGHT` conserve un espacement proportionnel dans la disposition par défaut à deux icônes. Pour un résultat optimal, utilisez une image d’arrière-plan ayant les mêmes dimensions en pixels que la fenêtre Finder.

### Remplacement des ressources du DMG

Les fichiers générés sous `build/darwin/` sont des ressources de projet ordinaires et peuvent être remplacés :

- `DMG_BACKGROUND` définit l’image affichée derrière le contenu de la fenêtre Finder.
- `DMG_VOLUME_ICON` définit l’icône affichée pour le volume monté.
- `DMG_FILE_ICON` définit l’icône affichée dans Finder pour le fichier `.dmg` obtenu.

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

### « L’application est endommagée et ne peut pas être ouverte »

L’application n’est pas signée. Signez-la avec un certificat Developer ID, ou indiquez aux utilisateurs comment contourner Gatekeeper :

```bash
xattr -cr /path/to/YourApp.app
```

### Échec de la notarisation

Problèmes courants :

- **Identifiants non valides** : réexécutez `xcrun notarytool store-credentials` (ou `wails3 signing`)
- **Runtime renforcé requis** : vérifiez que les droits incluent `com.apple.security.cs.allow-unsigned-executable-memory` si nécessaire
- **Horodatage manquant** : le processus de signature devrait inclure automatiquement un horodatage

### L’application compilée de manière croisée ne s’exécute pas

Les binaires macOS compilés de manière croisée ne sont pas signés. Transférez-les sur un Mac et signez-les avant de les tester :

```bash
codesign --force --deep --sign - YourApp.app
```
