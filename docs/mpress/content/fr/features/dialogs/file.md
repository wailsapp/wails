---
title: "Boîtes de dialogue de fichiers"
description: "Boîtes de dialogue d’ouverture, d’enregistrement et de sélection de dossiers"
slug: "features/dialogs/file"
sourcePath: "features/dialogs/file.md"
---

## Boîtes de dialogue de fichiers

Wails fournit des **boîtes de dialogue de fichiers natives**, dont l’apparence est adaptée à chaque plateforme, pour ouvrir et enregistrer des fichiers ainsi que sélectionner des dossiers. L’API, simple à utiliser, permet de filtrer les types de fichiers, d’effectuer des sélections multiples et de définir des emplacements par défaut.

![Un sélecteur de fichiers macOS natif ouvert par une application Wails](/assets/screenshots/file-dialog-macos.png)

L’API Wails délègue la sélection au système d’exploitation afin de conserver les fonctions familières de navigation, de filtrage et de sélection.

## Créer des boîtes de dialogue de fichiers

Accédez aux boîtes de dialogue de fichiers à l’aide du gestionnaire `app.Dialog` :

```go
app.Dialog.OpenFile()
app.Dialog.SaveFile()
```

## Boîte de dialogue Ouvrir un fichier

Sélectionnez les fichiers à ouvrir :

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Image").
    AddFilter("Images", "*.png;*.jpg;*.gif").
    AddFilter("All Files", "*.*").
    PromptForSingleSelection()

if err != nil || path == "" {
    return
}

openFile(path)
```

**Cas d’utilisation :**

- Ouvrir des documents
- Importer des fichiers
- Charger des images
- Sélectionner des fichiers de configuration

### Sélection d’un seul fichier

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Open Document").
    AddFilter("Text Files", "*.txt").
    AddFilter("All Files", "*.*").
    PromptForSingleSelection()

if err != nil || path == "" {
    // User cancelled or error occurred
    return
}

// Use selected file
data, _ := os.ReadFile(path)
```

### Sélection de plusieurs fichiers

```go
paths, err := app.Dialog.OpenFile().
    SetTitle("Select Images").
    AddFilter("Images", "*.png;*.jpg;*.jpeg;*.gif").
    PromptForMultipleSelection()

if err != nil {
    return
}

// Process all selected files
for _, path := range paths {
    processFile(path)
}
```

### Avec un répertoire par défaut

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Open File").
    SetDirectory("/Users/me/Documents").
    PromptForSingleSelection()
```

## Boîte de dialogue Enregistrer un fichier

Choisissez l’emplacement d’enregistrement :

```go
path, err := app.Dialog.SaveFile().
    SetFilename("document.txt").
    AddFilter("Text Files", "*.txt").
    AddFilter("All Files", "*.*").
    PromptForSingleSelection()

if err != nil || path == "" {
    return
}

saveFile(path, data)
```

**Cas d’utilisation :**

- Enregistrer des documents
- Exporter des données
- Créer des fichiers
- Enregistrer sous…

### Avec un nom de fichier par défaut

```go
path, err := app.Dialog.SaveFile().
    SetFilename("export.csv").
    AddFilter("CSV Files", "*.csv").
    PromptForSingleSelection()
```

### Avec un répertoire par défaut

```go
path, err := app.Dialog.SaveFile().
    SetDirectory("/Users/me/Documents").
    SetFilename("untitled.txt").
    PromptForSingleSelection()
```

### Confirmation du remplacement

```go
path, err := app.Dialog.SaveFile().
    SetFilename("document.txt").
    PromptForSingleSelection()

if err != nil || path == "" {
    return
}

// Check if file exists
if _, err := os.Stat(path); err == nil {
    dialog := app.Dialog.Question().
        SetTitle("Confirm Overwrite").
        SetMessage("File already exists. Overwrite?")

    overwriteBtn := dialog.AddButton("Overwrite")
    overwriteBtn.OnClick(func() {
        saveFile(path, data)
    })

    cancelBtn := dialog.AddButton("Cancel")
    dialog.SetDefaultButton(cancelBtn)
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
    return
}

saveFile(path, data)
```

## Boîte de dialogue Sélectionner un dossier

Choisissez un répertoire au moyen de la boîte de dialogue d’ouverture de fichier après avoir activé la sélection de répertoires :

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Output Folder").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()

if err != nil || path == "" {
    return
}

exportToFolder(path)
```

**Cas d’utilisation :**

- Choisir le répertoire de sortie
- Sélectionner l’espace de travail
- Choisir l’emplacement de sauvegarde
- Choisir le répertoire d’installation

### Avec un répertoire par défaut

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Folder").
    SetDirectory("/Users/me/Documents").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()
```

## Filtres de fichiers

Utilisez la méthode `AddFilter()` pour ajouter des filtres de types de fichiers aux boîtes de dialogue. Chaque appel ajoute une nouvelle option de filtrage.

### Filtres de base

```go
path, _ := app.Dialog.OpenFile().
    AddFilter("Text Files", "*.txt").
    AddFilter("All Files", "*.*").
    PromptForSingleSelection()
```

### Extensions multiples

Utilisez des points-virgules pour spécifier plusieurs extensions dans un même filtre :

```go
dialog := app.Dialog.OpenFile().
    AddFilter("Images", "*.png;*.jpg;*.jpeg;*.gif;*.bmp").
    AddFilter("Documents", "*.txt;*.doc;*.docx;*.pdf").
    AddFilter("All Files", "*.*")
```

### Format des motifs

Utilisez des **points-virgules** pour séparer plusieurs extensions dans un même filtre :

```go
// Multiple extensions separated by semicolons
AddFilter("Images", "*.png;*.jpg;*.gif")
```

## Exemples complets

### Ouvrir un fichier image

```go
func openImage(app *application.App) (image.Image, error) {
    path, err := app.Dialog.OpenFile().
        SetTitle("Select Image").
        AddFilter("Images", "*.png;*.jpg;*.jpeg;*.gif;*.bmp").
        PromptForSingleSelection()

    if err != nil {
        return nil, err
    }

    if path == "" {
        return nil, errors.New("no file selected")
    }

    // Open and decode image
    file, err := os.Open(path)
    if err != nil {
        app.Dialog.Error().
            SetTitle("Open Failed").
            SetMessage(err.Error()).
            Show()
        return nil, err
    }
    defer file.Close()

    img, _, err := image.Decode(file)
    if err != nil {
        app.Dialog.Error().
            SetTitle("Invalid Image").
            SetMessage("Could not decode image file.").
            Show()
        return nil, err
    }

    return img, nil
}
```

### Enregistrer un document avec validation

```go
func saveDocument(app *application.App, content string) {
    path, err := app.Dialog.SaveFile().
        SetFilename("document.txt").
        AddFilter("Text Files", "*.txt").
        AddFilter("Markdown Files", "*.md").
        AddFilter("All Files", "*.*").
        PromptForSingleSelection()

    if err != nil || path == "" {
        return
    }

    // Validate extension
    ext := filepath.Ext(path)
    if ext != ".txt" && ext != ".md" {
        dialog := app.Dialog.Question().
            SetTitle("Confirm Extension").
            SetMessage(fmt.Sprintf("Save as %s file?", ext))

        saveBtn := dialog.AddButton("Save")
        saveBtn.OnClick(func() {
            doSave(app, path, content)
        })

        cancelBtn := dialog.AddButton("Cancel")
        dialog.SetDefaultButton(cancelBtn)
        dialog.SetCancelButton(cancelBtn)
        dialog.Show()
        return
    }

    doSave(app, path, content)
}

func doSave(app *application.App, path, content string) {
    if err := os.WriteFile(path, []byte(content), 0644); err != nil {
        app.Dialog.Error().
            SetTitle("Save Failed").
            SetMessage(err.Error()).
            Show()
        return
    }

    app.Dialog.Info().
        SetTitle("Saved").
        SetMessage("Document saved successfully!").
        Show()
}
```

### Traitement de fichiers par lots

```go
func processMultipleFiles(app *application.App) {
    paths, err := app.Dialog.OpenFile().
        SetTitle("Select Files to Process").
        AddFilter("Images", "*.png;*.jpg").
        PromptForMultipleSelection()

    if err != nil || len(paths) == 0 {
        return
    }

    // Confirm processing
    dialog := app.Dialog.Question().
        SetTitle("Confirm Processing").
        SetMessage(fmt.Sprintf("Process %d file(s)?", len(paths)))

    processBtn := dialog.AddButton("Process")
    processBtn.OnClick(func() {
        // Process files
        var errs []error
        for i, path := range paths {
            if err := processFile(path); err != nil {
                errs = append(errs, err)
            }

            // Update progress
            // app.Event.Emit("progress", map[string]interface{}{
            //     "current": i + 1,
            //     "total":   len(paths),
            // })
            _ = i // suppress unused variable warning in example
        }

        // Show results
        if len(errs) > 0 {
            app.Dialog.Warning().
                SetTitle("Processing Complete").
                SetMessage(fmt.Sprintf("Processed %d files with %d errors.",
                    len(paths), len(errs))).
                Show()
        } else {
            app.Dialog.Info().
                SetTitle("Success").
                SetMessage(fmt.Sprintf("Processed %d files successfully!", len(paths))).
                Show()
        }
    })

    cancelBtn := dialog.AddButton("Cancel")
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
}
```

### Exporter après sélection d’un dossier

```go
func exportData(app *application.App, data []byte) {
    // Select output folder
    folder, err := app.Dialog.OpenFile().
        SetTitle("Select Export Folder").
        SetDirectory(getDefaultExportFolder()).
        CanChooseDirectories(true).
        CanChooseFiles(false).
        PromptForSingleSelection()

    if err != nil || folder == "" {
        return
    }

    // Generate filename
    filename := fmt.Sprintf("export_%s.csv",
        time.Now().Format("2006-01-02_15-04-05"))
    path := filepath.Join(folder, filename)

    // Save file
    if err := os.WriteFile(path, data, 0644); err != nil {
        app.Dialog.Error().
            SetTitle("Export Failed").
            SetMessage(err.Error()).
            Show()
        return
    }

    // Show success with option to open folder
    dialog := app.Dialog.Question().
        SetTitle("Export Complete").
        SetMessage(fmt.Sprintf("Exported to %s", filename))

    openBtn := dialog.AddButton("Open Folder")
    openBtn.OnClick(func() {
        openFolder(folder)
    })

    dialog.AddButton("OK")
    dialog.Show()
}
```

### Importer avec validation

```go
func importConfiguration(app *application.App) {
    path, err := app.Dialog.OpenFile().
        SetTitle("Import Configuration").
        AddFilter("JSON Files", "*.json").
        AddFilter("YAML Files", "*.yaml;*.yml").
        PromptForSingleSelection()

    if err != nil || path == "" {
        return
    }

    // Read file
    data, err := os.ReadFile(path)
    if err != nil {
        app.Dialog.Error().
            SetTitle("Read Failed").
            SetMessage(err.Error()).
            Show()
        return
    }

    // Validate configuration
    config, err := parseConfig(data)
    if err != nil {
        app.Dialog.Error().
            SetTitle("Invalid Configuration").
            SetMessage("File is not a valid configuration.").
            Show()
        return
    }

    // Confirm import
    dialog := app.Dialog.Question().
        SetTitle("Confirm Import").
        SetMessage("Import this configuration?")

    importBtn := dialog.AddButton("Import")
    importBtn.OnClick(func() {
        // Apply configuration
        if err := applyConfig(config); err != nil {
            app.Dialog.Error().
                SetTitle("Import Failed").
                SetMessage(err.Error()).
                Show()
            return
        }

        app.Dialog.Info().
            SetTitle("Success").
            SetMessage("Configuration imported successfully!").
            Show()
    })

    cancelBtn := dialog.AddButton("Cancel")
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
}
```

## Bonnes pratiques

### ✅ À faire

- **Fournissez des filtres de fichiers** – Aidez les utilisateurs à trouver leurs fichiers
- **Définissez des titres appropriés** – Donnez un contexte clair
- **Utilisez des répertoires par défaut** – Commencez dans un emplacement logique
- **Validez les sélections** – Vérifiez les types de fichiers
- **Gérez l’annulation** – L’utilisateur peut annuler l’opération
- **Demandez confirmation** – Pour les actions destructrices
- **Informez l’utilisateur du résultat** – Affichez des messages de réussite ou d’erreur

### ❌ À ne pas faire

- **Ne négligez pas la validation** – Vérifiez les types de fichiers
- **N’ignorez pas les erreurs** – Gérez l’annulation
- **N’utilisez pas de filtres génériques** – Soyez précis
- **N’oubliez pas « Tous les fichiers »** – Proposez toujours cette option
- **Ne codez pas les chemins en dur** – Utilisez le répertoire personnel de l’utilisateur
- **Ne supposez pas que le fichier existe** – Vérifiez son existence avant de l’ouvrir

## Différences entre les plateformes

### macOS

- NSOpenPanel/NSSavePanel natifs
- Présentation sous forme de feuille lorsque la boîte de dialogue est rattachée à une fenêtre
- Respecte le thème du système
- Prend en charge l’aperçu Quick Look
- Intégration des étiquettes et des favoris

### Windows

- Boîtes de dialogue natives d’ouverture et d’enregistrement de fichiers
- Respecte le thème du système
- Intégration des fichiers récents
- Prise en charge des emplacements réseau

### Linux

- Sélecteur de fichiers GTK
- Varie selon l’environnement de bureau
- Suit le thème du bureau
- Prise en charge des fichiers récents

## Étapes suivantes

@cards{cols="2"}
ℹ Boîtes de dialogue de message
Boîtes de dialogue d’information, d’avertissement et d’erreur.

[En savoir plus →](/features/dialogs/message/)

---
◆ Boîtes de dialogue personnalisées
Créez des fenêtres de dialogue personnalisées.

[En savoir plus →](/features/dialogs/custom/)

---
🚀 Liaisons
Appelez des fonctions Go depuis JavaScript.

[En savoir plus →](/features/bindings/methods/)

---
★ Événements
Utilisez des événements pour signaler la progression.

[En savoir plus →](/features/events/system/)

@end

---

**Des questions ?** Posez-les sur [Discord](https://discord.gg/JDdSxwjhGf) ou consultez les [exemples de boîtes de dialogue de fichiers](https://github.com/wailsapp/wails/tree/master/v3/examples/dialogs).
