---
title: "Associations de fichiers"
description: "Configurez les associations de fichiers de votre application Wails"
slug: "guides/file-associations"
sourcePath: "guides/file-associations.md"
---

Plateformes concernées : <span class="mpress-badge mpress-badge-note">Windows</span> <span class="mpress-badge mpress-badge-success">macOS</span>

<br/>

Les associations de fichiers permettent à votre application de gérer certains types de fichiers lorsque les utilisateurs les ouvrent. Elles sont particulièrement utiles pour les éditeurs de texte, les visionneuses d’images ou toute application qui traite des formats de fichiers particuliers. Ce guide explique comment implémenter des associations de fichiers dans votre application Wails v3.

## Vue d’ensemble

Dans Wails v3, les associations de fichiers sont actuellement prises en charge pour :

- Windows (paquets d’installation NSIS)
- macOS (paquets d’application)

## Configuration

Les associations de fichiers se configurent dans le fichier `config.yml` situé dans le répertoire `build` de votre projet.

### Configuration de base

Pour configurer les associations de fichiers :

1. Ouvrez `build/config.yml`
2. Ajoutez vos associations de fichiers dans la section `fileAssociations`
3. Exécutez `wails3 update build-assets` pour mettre à jour les ressources de compilation
4. Définissez le champ `FileAssociations` dans les options de l’application
5. Créez le paquet de votre application avec `wails3 package`

Voici un exemple de configuration :

```yaml
fileAssociations:
  - ext: myapp
    name: MyApp Document
    description: MyApp Document File
    iconName: myappFileIcon
    role: Editor
  - ext: custom
    name: Custom Format
    description: Custom File Format
    iconName: customFileIcon
    role: Editor
```

### Propriétés de configuration

| Propriété | Description | Plateforme |
| --- | --- | --- |
| ext | Extension de fichier sans le point initial (par exemple, `txt`) | Toutes |
| name | Nom d’affichage du type de fichier | Toutes |
| description | Description affichée dans les propriétés du fichier | Windows |
| iconName | Nom du fichier d’icône (sans extension) dans le dossier build | Toutes |
| role | Rôle de l’application pour ce type de fichier (par exemple, `Editor` ou `Viewer`) | macOS |
| mimeType | Type MIME du fichier (par exemple, `image/jpeg`) | macOS |

## Écoute des événements d’ouverture de fichiers

Pour gérer les événements d’ouverture de fichiers dans votre application, vous pouvez écouter l’événement `events.Common.ApplicationOpenedWithFile` :

```go {title="main.go"}
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

func main() {
    app := application.New(application.Options{
        Name: "MyApp",
        FileAssociations: []string{".txt", ".md"}, // Specify supported extensions
    })

    // Listen for files being used to open the application
    app.Event.OnApplicationEvent(events.Common.ApplicationOpenedWithFile, func(event *application.ApplicationEvent) {
        associatedFile := event.Context().Filename()
        app.Dialog.Info().SetMessage("Application opened with file: " + associatedFile).Show()
    })

    // Create your window and run the app...
}

```

## Tutoriel pas à pas

Voyons comment configurer les associations de fichiers pour un éditeur de texte simple :

@steps
### Créer les icônes
- Créez des icônes pour votre type de fichier (dimensions recommandées : 16x16, 32x32, 48x48 et 256x256)
- Enregistrez les icônes dans le dossier `build` de votre projet
- Nommez-les conformément à votre configuration `iconName` (par exemple, `textFileIcon.png`)

@note{type="tip"}
Vous pouvez utiliser `wails3 generate icons` pour générer les icônes requises. Exécutez `wails3 generate icons --help` pour plus d’informations.

@end

- Pour macOS, ajoutez une instruction de copie telle que `cp build/darwin/documenticon.icns {{.BIN_DIR}}/{{.APP_NAME}}.app/Contents/Resources` dans la tâche `create:app:bundle:`.

### Configurer les associations de fichiers
Modifiez le fichier `build/config.yml` pour ajouter vos associations de fichiers :

```yaml
# build/config.yml
fileAssociations:
  - ext: txt
    name: Text Document
    description: Plain Text Document
    iconName: textFileIcon
    role: Editor
```

### Mettre à jour les ressources de compilation
Exécutez la commande suivante pour mettre à jour les ressources de compilation :

```bash
wails3 update build-assets
```

### Définir les associations de fichiers dans les options de l’application
Dans votre fichier `main.go`, définissez le champ `FileAssociations` dans les options de l’application :

```go
app := application.New(application.Options{
  Name: "MyApp",
  FileAssociations: []string{".txt", ".md"}, // Specify supported extensions
})
```

@note{type="tip" title="Pourquoi les extensions de fichiers sont-elles requises à la fois dans la configuration de l’application et dans config.yml ?"}
Sous Windows, lorsqu’un fichier est ouvert au moyen d’une association de fichiers, l’application est lancée avec le nom du fichier comme premier argument. L’application ne peut pas savoir si ce premier argument est un fichier ou un argument de ligne de commande. Elle utilise donc le champ `FileAssociations` des options de l’application pour déterminer si le premier argument est ou non un fichier associé.

@end

### Créer le paquet de votre application
Créez le paquet de votre application à l’aide de la commande suivante :

```bash
wails3 package
```

L’application empaquetée sera créée dans le répertoire `bin`. Vous pourrez ensuite l’installer et la tester.

## Remarques supplémentaires

- Les icônes devraient être fournies au format PNG dans le dossier build
- Pour tester les associations de fichiers, vous devez installer l’application empaquetée

@end
