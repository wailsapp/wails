---
title: "Présentation des boîtes de dialogue"
description: "Affichez des boîtes de dialogue système natives dans votre application"
slug: "features/dialogs/overview"
sourcePath: "features/dialogs/overview.md"
---

## Boîtes de dialogue natives

Wails fournit des **boîtes de dialogue système natives** compatibles avec toutes les plateformes : boîtes de message (information, avertissement, erreur, question), boîtes de dialogue de fichiers (ouverture, enregistrement, dossier) et fenêtres de dialogue personnalisées avec l’apparence et le comportement natifs de chaque plateforme.

![Une boîte de dialogue de question Wails sous macOS avec les boutons Annuler et Ignorer](/assets/screenshots/dialog-question-macos.png)

La même API respecte les conventions de chaque plateforme prise en charge. Cet exemple sous macOS présente une boîte de dialogue de question attachée à sa fenêtre Wails, avec le bouton par défaut et le bouton d’annulation.

## Démarrage rapide

```go
// Information dialog
app.Dialog.Info().
    SetTitle("Success").
    SetMessage("File saved successfully!").
    Show()

// Question dialog with button callbacks
dialog := app.Dialog.Question().
    SetTitle("Confirm").
    SetMessage("Delete this file?")

deleteBtn := dialog.AddButton("Delete")
deleteBtn.OnClick(func() {
    deleteFile()
})

cancelBtn := dialog.AddButton("Cancel")
dialog.SetDefaultButton(cancelBtn)
dialog.SetCancelButton(cancelBtn)
dialog.Show()

// File open dialog
path, _ := app.Dialog.OpenFile().
    SetTitle("Select Image").
    AddFilter("Images", "*.png;*.jpg").
    PromptForSingleSelection()
```

**C’est tout !** Des boîtes de dialogue natives avec un minimum de code.

## Accéder aux boîtes de dialogue

Accédez aux boîtes de dialogue par l’intermédiaire du gestionnaire `app.Dialog` :

```go
app.Dialog.Info()
app.Dialog.Question()
app.Dialog.Warning()
app.Dialog.Error()
app.Dialog.OpenFile()
app.Dialog.SaveFile()
```

## Types de boîtes de dialogue

### Boîte de dialogue d’information

Affichez des messages simples :

```go
app.Dialog.Info().
    SetTitle("Welcome").
    SetMessage("Welcome to our application!").
    Show()
```

**Cas d’utilisation :**

- Messages de réussite
- Avis d’information
- Confirmations d’achèvement

### Boîte de dialogue d’avertissement

Affichez des avertissements :

```go
app.Dialog.Warning().
    SetTitle("Warning").
    SetMessage("This action cannot be undone.").
    Show()
```

**Cas d’utilisation :**

- Avertissements non critiques
- Avis d’obsolescence
- Messages de mise en garde

### Boîte de dialogue d’erreur

Affichez des erreurs :

```go
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Failed to save file: " + err.Error()).
    Show()
```

**Cas d’utilisation :**

- Messages d’erreur
- Notifications d’échec
- Gestion des exceptions

### Boîte de dialogue de question

Posez des questions aux utilisateurs et traitez leurs réponses au moyen des fonctions de rappel des boutons :

```go
dialog := app.Dialog.Question().
    SetTitle("Confirm Delete").
    SetMessage("Are you sure you want to delete this file?")

deleteBtn := dialog.AddButton("Delete")
deleteBtn.OnClick(func() {
    deleteFile()
})

cancelBtn := dialog.AddButton("Cancel")
dialog.SetDefaultButton(cancelBtn)
dialog.SetCancelButton(cancelBtn)
dialog.Show()
```

**Cas d’utilisation :**

- Confirmation d’actions
- Questions Oui/Non
- Choix multiple

## Boîtes de dialogue de fichiers

### Boîte de dialogue d’ouverture de fichier

Sélectionnez les fichiers à ouvrir :

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Image").
    AddFilter("Images", "*.png;*.jpg;*.gif").
    AddFilter("All Files", "*.*").
    PromptForSingleSelection()

if err == nil && path != "" {
    openFile(path)
}
```

**Sélection multiple :**

```go
paths, err := app.Dialog.OpenFile().
    SetTitle("Select Images").
    AddFilter("Images", "*.png;*.jpg").
    PromptForMultipleSelection()

if err == nil {
    for _, path := range paths {
        processFile(path)
    }
}
```

### Boîte de dialogue d’enregistrement de fichier

Choisissez l’emplacement d’enregistrement :

```go
path, err := app.Dialog.SaveFile().
    SetFilename("document.txt").
    AddFilter("Text Files", "*.txt").
    AddFilter("All Files", "*.*").
    PromptForSingleSelection()

if err == nil && path != "" {
    saveFile(path)
}
```

### Boîte de dialogue de sélection de dossier

Choisissez un répertoire à l’aide de la boîte de dialogue d’ouverture de fichier, avec la sélection de répertoires activée :

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Output Folder").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()

if err == nil && path != "" {
    exportToFolder(path)
}
```

## Options des boîtes de dialogue

### Titre et message

```go
dialog := app.Dialog.Info().
    SetTitle("Success").
    SetMessage("Operation completed successfully!")
```

### Boutons

**Bouton par défaut des boîtes de dialogue simples :**

Les boîtes de dialogue d’information, d’avertissement et d’erreur affichent par défaut un bouton « OK » :

```go
app.Dialog.Info().
    SetMessage("Done!").
    Show()
```

**Boutons personnalisés des boîtes de dialogue de question :**

Utilisez `AddButton()` pour ajouter des boutons. Cette méthode renvoie un objet `*Button` que vous pouvez configurer avec des fonctions de rappel :

```go
dialog := app.Dialog.Question().
    SetMessage("Choose action")

save := dialog.AddButton("Save")
save.OnClick(func() {
    saveDocument()
})

dontSave := dialog.AddButton("Don't Save")
dontSave.OnClick(func() {
    discardChanges()
})

cancel := dialog.AddButton("Cancel")
// No callback needed - just dismisses dialog

dialog.SetDefaultButton(save)
dialog.SetCancelButton(cancel)
dialog.Show()
```

**Boutons par défaut et d’annulation :**

Utilisez `SetDefaultButton()` pour indiquer le bouton mis en évidence et activé par Entrée. Utilisez `SetCancelButton()` pour indiquer le bouton activé par Échap.

```go
dialog := app.Dialog.Question().
    SetMessage("Delete file?")

deleteBtn := dialog.AddButton("Delete")
deleteBtn.OnClick(func() {
    performDelete()
})

cancelBtn := dialog.AddButton("Cancel")
dialog.SetDefaultButton(cancelBtn)  // Safe option highlighted by default
dialog.SetCancelButton(cancelBtn)   // Escape triggers Cancel
dialog.Show()
```

### Attachement à une fenêtre

Attachez la boîte de dialogue à une fenêtre précise :

```go
dialog := app.Dialog.Info().
    SetMessage("Window-specific message").
    AttachToWindow(window)

dialog.Show()
```

**Comportement :**

- La boîte de dialogue apparaît centrée sur la fenêtre parente
- La fenêtre parente est désactivée pendant l’affichage de la boîte de dialogue
- La boîte de dialogue se déplace avec la fenêtre parente (macOS)

## Comportement selon la plateforme

@tabs{sync-key="platform"}
[macOS]
**Boîtes de dialogue sous macOS :**

- Apparence native de NSAlert
- Respect du thème système (clair/sombre)
- Prise en charge de la navigation au clavier
- Raccourcis standard (⌘. pour Annuler)
- Fonctionnalités d’accessibilité intégrées
- Présentation sous forme de feuille lorsque la boîte de dialogue est attachée à une fenêtre

**Exemple :**

```go
// Appears as sheet on macOS
dialog := app.Dialog.Question().
    SetMessage("Save changes?").
    AttachToWindow(window)
dialog.AddButton("Yes")
dialog.AddButton("No")
dialog.Show()
```

[Windows]
**Boîtes de dialogue sous Windows :**

- Apparence native de TaskDialog
- Respect du thème système
- Prise en charge de la navigation au clavier
- Raccourcis standard (Échap pour Annuler)
- Fonctionnalités d’accessibilité intégrées
- Modale par rapport à la fenêtre parente

**Exemple :**

```go
// Modal dialog on Windows
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Operation failed").
    Show()
```

[Linux]
**Boîtes de dialogue Linux :**

- Apparence des boîtes de dialogue GTK
- Respectent le thème du bureau
- Prennent en charge la navigation au clavier
- Intégration à l’environnement de bureau
- Varie selon l’environnement de bureau (GNOME, KDE, etc.)

**Exemple :**

```go
// GTK dialog on Linux
app.Dialog.Info().
    SetMessage("Update complete").
    Show()
```

@end

## Modèles courants

### Demander confirmation avant une action destructive

```go
func deleteFile(app *application.App, path string) {
    dialog := app.Dialog.Question().
        SetTitle("Confirm Delete").
        SetMessage(fmt.Sprintf("Delete %s?", filepath.Base(path)))

    deleteBtn := dialog.AddButton("Delete")
    deleteBtn.OnClick(func() {
        if err := os.Remove(path); err != nil {
            app.Dialog.Error().
                SetTitle("Delete Failed").
                SetMessage(err.Error()).
                Show()
        }
    })

    cancelBtn := dialog.AddButton("Cancel")
    dialog.SetDefaultButton(cancelBtn)
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
}
```

### Gestion des erreurs avec une boîte de dialogue

```go
func saveDocument(app *application.App, path string, data []byte) {
    if err := os.WriteFile(path, data, 0644); err != nil {
        app.Dialog.Error().
            SetTitle("Save Failed").
            SetMessage(fmt.Sprintf("Could not save file: %v", err)).
            Show()
        return
    }

    app.Dialog.Info().
        SetTitle("Success").
        SetMessage("File saved successfully!").
        Show()
}
```

### Sélection de fichiers avec validation

```go
func selectImageFile(app *application.App) (string, error) {
    path, err := app.Dialog.OpenFile().
        SetTitle("Select Image").
        AddFilter("Images", "*.png;*.jpg;*.jpeg;*.gif").
        PromptForSingleSelection()

    if err != nil {
        return "", err
    }

    if path == "" {
        return "", errors.New("no file selected")
    }

    // Validate file
    if !isValidImage(path) {
        app.Dialog.Error().
            SetTitle("Invalid File").
            SetMessage("Selected file is not a valid image.").
            Show()
        return "", errors.New("invalid image")
    }

    return path, nil
}
```

### Enchaînement de boîtes de dialogue en plusieurs étapes

```go
func exportData(app *application.App) {
    // Step 1: Confirm export
    dialog := app.Dialog.Question().
        SetTitle("Export Data").
        SetMessage("Export all data to CSV?")

    exportBtn := dialog.AddButton("Export")
    exportBtn.OnClick(func() {
        // Step 2: Select destination
        path, err := app.Dialog.SaveFile().
            SetFilename("export.csv").
            AddFilter("CSV Files", "*.csv").
            PromptForSingleSelection()

        if err != nil || path == "" {
            return
        }

        // Step 3: Perform export
        if err := performExport(path); err != nil {
            app.Dialog.Error().
                SetTitle("Export Failed").
                SetMessage(err.Error()).
                Show()
            return
        }

        // Step 4: Success
        app.Dialog.Info().
            SetTitle("Export Complete").
            SetMessage("Data exported successfully!").
            Show()
    })

    cancelBtn := dialog.AddButton("Cancel")
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
}
```

## Bonnes pratiques

### ✅ À faire

- **Utilisez des boîtes de dialogue natives** — Elles offrent une meilleure expérience utilisateur que les boîtes personnalisées
- **Affichez des messages clairs** — Soyez précis
- **Définissez des titres appropriés** — Le contexte est important
- **Choisissez judicieusement les boutons par défaut** — Utilisez l’option la plus sûre par défaut
- **Gérez l’annulation** — L’utilisateur peut annuler
- **Validez les fichiers sélectionnés** — Vérifiez leur type

### ❌ À ne pas faire

- **N’abusez pas des boîtes de dialogue** — Elles interrompent le flux de travail
- **Ne les utilisez pas pour des messages fréquents** — Utilisez des notifications
- **N’oubliez pas de gérer les erreurs** — L’utilisateur peut annuler
- **Ne bloquez pas inutilement l’application** — Envisagez d’autres solutions
- **N’utilisez pas de messages génériques** — Soyez précis
- **N’ignorez pas les différences entre les plateformes** — Testez sur toutes les plateformes

## Étapes suivantes

@cards{cols="2"}
ℹ Boîtes de dialogue de message
Boîtes de dialogue d’information, d’avertissement et d’erreur.

[En savoir plus →](/features/dialogs/message/)

---
📖 Boîtes de dialogue de fichiers
Ouverture, enregistrement et sélection de dossiers.

[En savoir plus →](/features/dialogs/file/)

---
◆ Boîtes de dialogue personnalisées
Créez des fenêtres de dialogue personnalisées.

[En savoir plus →](/features/dialogs/custom/)

---
▣ Fenêtres
Découvrez la gestion des fenêtres.

[En savoir plus →](/features/windows/basics/)

@end

---

**Des questions ?** Posez-les sur [Discord](https://discord.gg/JDdSxwjhGf) ou consultez les [exemples de boîtes de dialogue](https://github.com/wailsapp/wails/tree/master/v3/examples/dialogs).
