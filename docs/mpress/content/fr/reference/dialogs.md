---
title: "API des boîtes de dialogue"
description: "Référence complète des API de boîtes de dialogue natives"
slug: "reference/dialogs"
sourcePath: "reference/dialogs.md"
---

## Vue d’ensemble

L’API Dialogs fournit des méthodes permettant d’afficher des boîtes de dialogue natives de sélection de fichiers et de message. Accédez aux boîtes de dialogue via le gestionnaire `app.Dialog`.

**Types de boîtes de dialogue :**

- **Boîtes de dialogue de fichiers** — boîtes de dialogue d’ouverture et d’enregistrement
- **Boîtes de dialogue de message** — boîtes de dialogue d’information, d’erreur, d’avertissement et de question

Toutes les boîtes de dialogue sont des **boîtes de dialogue natives du système d’exploitation** qui respectent l’apparence et le comportement de la plateforme.

## Accès aux boîtes de dialogue

Accédez aux boîtes de dialogue via le gestionnaire `app.Dialog` :

```go
app.Dialog.OpenFile()
app.Dialog.SaveFile()
app.Dialog.Info()
app.Dialog.Question()
app.Dialog.Warning()
app.Dialog.Error()
```

## Boîtes de dialogue de fichiers

### OpenFile()

Crée une boîte de dialogue d’ouverture de fichier.

```go
func (dm *DialogManager) OpenFile() *OpenFileDialogStruct
```

**Exemple :**

```go
dialog := app.Dialog.OpenFile()
```

### Méthodes de OpenFileDialogStruct

#### SetTitle()

Définit le titre de la boîte de dialogue.

```go
func (d *OpenFileDialogStruct) SetTitle(title string) *OpenFileDialogStruct
```

**Exemple :**

```go
dialog.SetTitle("Select Image")
```

#### AddFilter()

Ajoute un filtre de type de fichier.

```go
func (d *OpenFileDialogStruct) AddFilter(displayName, pattern string) *OpenFileDialogStruct
```

**Paramètres :**

- `displayName` — Description du filtre affichée à l’utilisateur (par exemple, « Images » ou « Documents »)
- `pattern` — Liste d’extensions séparées par des points-virgules (par exemple, « *.png;*.jpg »)

**Exemple :**

```go
dialog.AddFilter("Images", "*.png;*.jpg;*.gif").
    AddFilter("Documents", "*.pdf;*.docx").
    AddFilter("All Files", "*.*")
```

#### SetDirectory()

Définit le répertoire initial.

```go
func (d *OpenFileDialogStruct) SetDirectory(directory string) *OpenFileDialogStruct
```

**Exemple :**

```go
homeDir, _ := os.UserHomeDir()
dialog.SetDirectory(homeDir)
```

#### CanChooseDirectories()

Active ou désactive la sélection de répertoires.

```go
func (d *OpenFileDialogStruct) CanChooseDirectories(canChooseDirectories bool) *OpenFileDialogStruct
```

**Exemple (sélection d’un dossier) :**

```go
// Select folders instead of files
path, err := app.Dialog.OpenFile().
    SetTitle("Select Folder").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()
```

#### CanChooseFiles()

Active ou désactive la sélection de fichiers.

```go
func (d *OpenFileDialogStruct) CanChooseFiles(canChooseFiles bool) *OpenFileDialogStruct
```

#### CanCreateDirectories()

Active ou désactive la création de répertoires.

```go
func (d *OpenFileDialogStruct) CanCreateDirectories(canCreateDirectories bool) *OpenFileDialogStruct
```

#### ShowHiddenFiles()

Affiche ou masque les fichiers cachés.

```go
func (d *OpenFileDialogStruct) ShowHiddenFiles(showHiddenFiles bool) *OpenFileDialogStruct
```

#### AttachToWindow()

Associe la boîte de dialogue à une fenêtre précise.

```go
func (d *OpenFileDialogStruct) AttachToWindow(window Window) *OpenFileDialogStruct
```

#### PromptForSingleSelection()

Affiche la boîte de dialogue et renvoie le fichier sélectionné.

```go
func (d *OpenFileDialogStruct) PromptForSingleSelection() (string, error)
```

**Valeurs renvoyées :**

- `string` — Chemin du fichier sélectionné. Considérez une chaîne vide comme une absence de sélection (selon le système d’exploitation, les implémentations propres à la plateforme peuvent renvoyer soit une chaîne vide, soit une erreur non nulle en cas d’annulation).
- `error` — Valeur non nulle si la boîte de dialogue elle-même n’a pas pu être affichée.

**Exemple :**

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Image").
    AddFilter("Images", "*.png;*.jpg;*.gif").
    PromptForSingleSelection()

if err != nil {
    // The dialog failed to present (rare).
    return
}
if path == "" {
    // User cancelled.
    return
}

// Use the selected file
processFile(path)
```

#### PromptForMultipleSelection()

Affiche la boîte de dialogue et renvoie les fichiers sélectionnés.

```go
func (d *OpenFileDialogStruct) PromptForMultipleSelection() ([]string, error)
```

**Valeurs renvoyées :**

- `[]string` — Tableau des chemins des fichiers sélectionnés
- `error` — Erreur si l’affichage de la boîte de dialogue a échoué

**Exemple :**

```go
paths, err := app.Dialog.OpenFile().
    SetTitle("Select Images").
    AddFilter("Images", "*.png;*.jpg").
    PromptForMultipleSelection()

if err != nil {
    return
}

for _, path := range paths {
    processFile(path)
}
```

### SaveFile()

Crée une boîte de dialogue d’enregistrement de fichier.

```go
func (dm *DialogManager) SaveFile() *SaveFileDialogStruct
```

**Exemple :**

```go
dialog := app.Dialog.SaveFile()
```

### Méthodes de SaveFileDialogStruct

#### SetTitle()

Définit le titre de la boîte de dialogue.

```go
func (d *SaveFileDialogStruct) SetTitle(title string) *SaveFileDialogStruct
```

#### SetFilename()

Définit le nom de fichier par défaut.

```go
func (d *SaveFileDialogStruct) SetFilename(filename string) *SaveFileDialogStruct
```

**Exemple :**

```go
dialog.SetFilename("document.pdf")
```

#### AddFilter()

Ajoute un filtre de type de fichier.

```go
func (d *SaveFileDialogStruct) AddFilter(displayName, pattern string) *SaveFileDialogStruct
```

**Exemple :**

```go
dialog.AddFilter("PDF Document", "*.pdf").
    AddFilter("Text Document", "*.txt")
```

#### SetDirectory()

Définit le répertoire initial.

```go
func (d *SaveFileDialogStruct) SetDirectory(directory string) *SaveFileDialogStruct
```

#### AttachToWindow()

Associe la boîte de dialogue à une fenêtre précise.

```go
func (d *SaveFileDialogStruct) AttachToWindow(window Window) *SaveFileDialogStruct
```

#### PromptForSingleSelection()

Affiche la boîte de dialogue et renvoie le chemin d’enregistrement.

```go
func (d *SaveFileDialogStruct) PromptForSingleSelection() (string, error)
```

**Exemple :**

```go
path, err := app.Dialog.SaveFile().
    SetTitle("Save Document").
    SetFilename("untitled.pdf").
    AddFilter("PDF Document", "*.pdf").
    PromptForSingleSelection()

if err != nil {
    // User cancelled
    return
}

// Save to the selected path
saveDocument(path)
```

### Sélection d’un dossier

Il n’existe pas de `SelectFolderDialog` distinct. Utilisez `OpenFile()` avec les options de sélection de répertoires :

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Output Folder").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()

if err != nil {
    // User cancelled
    return
}

// Use the selected folder
outputDir = path
```

## Boîtes de dialogue de message

Toutes les boîtes de dialogue de message renvoient `*MessageDialog` et partagent les mêmes méthodes.

### Info()

Crée une boîte de dialogue d’information.

```go
func (dm *DialogManager) Info() *MessageDialog
```

**Exemple :**

```go
app.Dialog.Info().
    SetTitle("Success").
    SetMessage("File saved successfully!").
    Show()
```

### Error()

Crée une boîte de dialogue d’erreur.

```go
func (dm *DialogManager) Error() *MessageDialog
```

**Exemple :**

```go
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Failed to save file: " + err.Error()).
    Show()
```

### Warning()

Crée une boîte de dialogue d’avertissement.

```go
func (dm *DialogManager) Warning() *MessageDialog
```

**Exemple :**

```go
app.Dialog.Warning().
    SetTitle("Warning").
    SetMessage("This action cannot be undone.").
    Show()
```

### Question()

Crée une boîte de dialogue de question avec des boutons personnalisés.

```go
func (dm *DialogManager) Question() *MessageDialog
```

**Exemple :**

```go
dialog := app.Dialog.Question().
    SetTitle("Confirm").
    SetMessage("Do you want to save changes?")

save := dialog.AddButton("Save")
save.OnClick(func() {
    saveDocument()
})

dontSave := dialog.AddButton("Don't Save")
dontSave.OnClick(func() {
    // Continue without saving
})

cancel := dialog.AddButton("Cancel")
cancel.OnClick(func() {
    // Do nothing
})

dialog.SetDefaultButton(save)
dialog.SetCancelButton(cancel)
dialog.Show()
```

### Méthodes de MessageDialog

#### SetTitle()

Définit le titre de la boîte de dialogue.

```go
func (d *MessageDialog) SetTitle(title string) *MessageDialog
```

#### SetMessage()

Définit le message de la boîte de dialogue.

```go
func (d *MessageDialog) SetMessage(message string) *MessageDialog
```

#### SetIcon()

Définit une icône personnalisée pour la boîte de dialogue.

```go
func (d *MessageDialog) SetIcon(icon []byte) *MessageDialog
```

#### AddButton()

Ajoute un bouton à la boîte de dialogue et renvoie ce bouton afin de permettre sa configuration.

```go
func (d *MessageDialog) AddButton(label string) *Button
```

**Renvoie :** `*Button` — L’instance du bouton permettant de poursuivre sa configuration

**Exemple :**

```go
button := dialog.AddButton("OK")
button.OnClick(func() {
    // Handle click
})
```

#### SetDefaultButton()

Définit le bouton par défaut (activé en appuyant sur Entrée).

```go
func (d *MessageDialog) SetDefaultButton(button *Button) *MessageDialog
```

**Exemple :**

```go
yes := dialog.AddButton("Yes")
no := dialog.AddButton("No")
dialog.SetDefaultButton(yes)
```

#### SetCancelButton()

Définit le bouton d’annulation (activé en appuyant sur Échap).

```go
func (d *MessageDialog) SetCancelButton(button *Button) *MessageDialog
```

**Exemple :**

```go
ok := dialog.AddButton("OK")
cancel := dialog.AddButton("Cancel")
dialog.SetCancelButton(cancel)
```

#### AttachToWindow()

Associe la boîte de dialogue à une fenêtre donnée.

```go
func (d *MessageDialog) AttachToWindow(window Window) *MessageDialog
```

#### Show()

Affiche la boîte de dialogue. Les fonctions de rappel des boutons traitent les réponses de l’utilisateur.

```go
func (d *MessageDialog) Show()
```

**Remarque :** `Show()` ne renvoie aucune valeur. Utilisez les fonctions de rappel des boutons pour traiter les réponses de l’utilisateur.

### Méthodes des boutons

#### OnClick()

Définit la fonction de rappel exécutée lorsque l’utilisateur clique sur le bouton.

```go
func (b *Button) OnClick(callback func()) *Button
```

#### SetAsDefault()

Définit ce bouton comme bouton par défaut.

```go
func (b *Button) SetAsDefault() *Button
```

#### SetAsCancel()

Définit ce bouton comme bouton d’annulation.

```go
func (b *Button) SetAsCancel() *Button
```

## Exemples complets

### Exemple de sélection de fichiers

```go
type FileService struct {
    app *application.App
}

func (s *FileService) OpenImage() (string, error) {
    path, err := s.app.Dialog.OpenFile().
        SetTitle("Select Image").
        AddFilter("Images", "*.png;*.jpg;*.jpeg;*.gif").
        AddFilter("All Files", "*.*").
        PromptForSingleSelection()

    if err != nil {
        return "", err
    }

    return path, nil
}

func (s *FileService) SaveDocument(defaultName string) (string, error) {
    path, err := s.app.Dialog.SaveFile().
        SetTitle("Save Document").
        SetFilename(defaultName).
        AddFilter("PDF Document", "*.pdf").
        AddFilter("Text Document", "*.txt").
        PromptForSingleSelection()

    if err != nil {
        return "", err
    }

    return path, nil
}

func (s *FileService) SelectOutputFolder() (string, error) {
    path, err := s.app.Dialog.OpenFile().
        SetTitle("Select Output Folder").
        CanChooseDirectories(true).
        CanChooseFiles(false).
        PromptForSingleSelection()

    if err != nil {
        return "", err
    }

    return path, nil
}
```

### Exemple de boîte de dialogue de confirmation

```go
func (s *Service) DeleteItem(app *application.App, id string) {
    dialog := app.Dialog.Question().
        SetTitle("Confirm Delete").
        SetMessage("Are you sure you want to delete this item?")

    deleteBtn := dialog.AddButton("Delete")
    deleteBtn.OnClick(func() {
        deleteFromDatabase(id)
    })

    cancelBtn := dialog.AddButton("Cancel")
    // Cancel does nothing

    dialog.SetDefaultButton(cancelBtn) // Default to Cancel for safety
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
}
```

### Boîte de dialogue d’enregistrement des modifications

```go
func (s *Editor) PromptSaveChanges(app *application.App) {
    dialog := app.Dialog.Question().
        SetTitle("Unsaved Changes").
        SetMessage("Do you want to save your changes before closing?")

    save := dialog.AddButton("Save")
    save.OnClick(func() {
        s.Save()
        s.Close()
    })

    dontSave := dialog.AddButton("Don't Save")
    dontSave.OnClick(func() {
        s.Close()
    })

    cancel := dialog.AddButton("Cancel")
    // Cancel does nothing, dialog closes

    dialog.SetDefaultButton(save)
    dialog.SetCancelButton(cancel)
    dialog.Show()
}
```

### Traitement de plusieurs fichiers

```go
func (s *Service) ProcessMultipleFiles(app *application.App) error {
    // Select multiple files
    paths, err := app.Dialog.OpenFile().
        SetTitle("Select Files to Process").
        AddFilter("Images", "*.png;*.jpg").
        PromptForMultipleSelection()

    if err != nil {
        return err
    }

    if len(paths) == 0 {
        app.Dialog.Info().
            SetTitle("No Files Selected").
            SetMessage("Please select at least one file.").
            Show()
        return nil
    }

    // Process files
    for _, path := range paths {
        err := processFile(path)
        if err != nil {
            app.Dialog.Error().
                SetTitle("Processing Error").
                SetMessage(fmt.Sprintf("Failed to process %s: %v", path, err)).
                Show()
            continue
        }
    }

    // Show completion
    app.Dialog.Info().
        SetTitle("Complete").
        SetMessage(fmt.Sprintf("Successfully processed %d files", len(paths))).
        Show()

    return nil
}
```

### Gestion des erreurs avec les boîtes de dialogue

```go
func (s *Service) SaveFile(app *application.App, data []byte) error {
    // Select save location
    path, err := app.Dialog.SaveFile().
        SetTitle("Save File").
        SetFilename("data.json").
        AddFilter("JSON File", "*.json").
        PromptForSingleSelection()

    if err != nil {
        // User cancelled - not an error
        return nil
    }

    // Attempt to save
    err = os.WriteFile(path, data, 0644)
    if err != nil {
        // Show error dialog
        app.Dialog.Error().
            SetTitle("Save Failed").
            SetMessage(fmt.Sprintf("Could not save file: %v", err)).
            Show()
        return err
    }

    // Show success
    app.Dialog.Info().
        SetTitle("Success").
        SetMessage("File saved successfully!").
        Show()

    return nil
}
```

### Valeurs par défaut propres à chaque plateforme

```go
import (
    "os"
    "path/filepath"
    "runtime"
)

func (s *Service) GetDefaultDirectory() string {
    homeDir, _ := os.UserHomeDir()

    switch runtime.GOOS {
    case "windows":
        return filepath.Join(homeDir, "Documents")
    case "darwin":
        return filepath.Join(homeDir, "Documents")
    case "linux":
        return filepath.Join(homeDir, "Documents")
    default:
        return homeDir
    }
}

func (s *Service) OpenWithDefaults(app *application.App) (string, error) {
    return app.Dialog.OpenFile().
        SetTitle("Open File").
        SetDirectory(s.GetDefaultDirectory()).
        AddFilter("All Files", "*.*").
        PromptForSingleSelection()
}
```

## Bonnes pratiques

### À faire

- **Utilisez des boîtes de dialogue natives** — Elles respectent l’apparence et le comportement de la plateforme
- **Donnez des titres explicites** — Aidez les utilisateurs à comprendre l’objectif de la boîte de dialogue
- **Définissez des filtres appropriés** — Guidez les utilisateurs vers les types de fichiers adéquats
- **Gérez l’annulation** — Vérifiez les erreurs, car l’utilisateur peut annuler l’opération
- **Demandez confirmation avant toute action destructive** — Utilisez des boîtes de dialogue Question
- **Fournissez un retour d’information** — Utilisez des boîtes de dialogue Info pour les messages de réussite
- **Définissez des valeurs par défaut pertinentes** — Répertoire et nom de fichier par défaut, etc.
- **Utilisez des fonctions de rappel pour les actions des boutons** — Traitez correctement les réponses de l’utilisateur

### À ne pas faire

- **N’ignorez pas les erreurs** — L’annulation par l’utilisateur renvoie une erreur
- **N’utilisez pas de libellés de boutons ambigus** — Soyez précis : « Enregistrer »/« Annuler »
- **N’abusez pas des boîtes de dialogue** — Elles interrompent le flux de travail
- **N’affichez pas d’erreur en cas d’annulation** — Il s’agit d’une action normale
- **N’oubliez pas les filtres de fichiers** — Aidez les utilisateurs à trouver les fichiers appropriés
- **Ne codez pas les chemins en dur** — Utilisez os.UserHomeDir() ou une fonction similaire

## Types de boîtes de dialogue selon la plateforme

### macOS

- Les boîtes de dialogue se déroulent depuis la barre de titre
- Style « feuille » attaché à la fenêtre parente
- Apparence native de macOS

### Windows

- Boîtes de dialogue Windows standard
- Respecte les directives de conception de Windows
- Apparence moderne de Windows 10/11

### Linux

- Boîtes de dialogue GTK sur les systèmes basés sur GTK
- Boîtes de dialogue Qt sur les systèmes basés sur Qt
- S’adapte à l’environnement de bureau

#### Comportement des boîtes de dialogue sous Linux

Sous Linux, la compilation GTK4 par défaut utilise **xdg-desktop-portal** pour les boîtes de dialogue de fichiers, ce qui assure une intégration native au bureau, mais signifie que certaines options sont sans effet. L’ancienne voie GTK3 (`-tags gtk3`) conserve un contrôle programmatique total sur ces options :

| Option | GTK3 (`-tags gtk3`) | GTK4 (par défaut) | Remarques |
| --- | --- | --- | --- |
| `ShowHiddenFiles()` | ✅ Fonctionne | ❌ Sans effet | L’utilisateur contrôle cette option au moyen du bouton bascule de la boîte de dialogue (Ctrl+H ou menu) |
| `CanCreateDirectories()` | ✅ Fonctionne | ❌ Sans effet | Toujours activé dans le portail |
| `ResolvesAliases()` | ✅ Fonctionne | ❌ Sans effet | Le portail gère la résolution des liens symboliques |
| `SetButtonText()` | ✅ Fonctionne | ✅ Fonctionne | Le texte personnalisé du bouton d’acceptation fonctionne |

**Pourquoi ces limitations existent :** les boîtes de dialogue de GTK4 fondées sur le portail délèguent le contrôle de l’interface utilisateur à l’environnement de bureau (GNOME, KDE, etc.). Ce comportement est intentionnel : le portail offre une expérience utilisateur cohérente entre les applications et respecte les préférences de l’utilisateur.

@note{type="info"}
La compilation GTK4 par défaut utilise les boîtes de dialogue adossées au portail. Si votre application nécessite un contrôle programmatique total sur les options de boîte de dialogue ci-dessus, utilisez l’ancienne voie `-tags gtk3` lors de la compilation (prise en charge jusqu’à la version v3.0.x ; supprimée dans la version v3.1) — consultez [Empaquetage sous Linux — Prise en charge de l’ancienne version GTK3](/guides/build/linux/#legacy-gtk3-support).

@end

## Modèles courants

### Modèle « Enregistrer sous »

```go
func (s *Service) SaveAs(app *application.App, currentPath string) (string, error) {
    // Extract filename from current path
    filename := filepath.Base(currentPath)

    // Show save dialog
    path, err := app.Dialog.SaveFile().
        SetTitle("Save As").
        SetFilename(filename).
        PromptForSingleSelection()

    if err != nil {
        return "", err
    }

    return path, nil
}
```

### Modèle « Ouvrir un élément récent »

```go
func (s *Service) OpenRecent(app *application.App, recentPath string) error {
    // Check if file still exists
    if _, err := os.Stat(recentPath); os.IsNotExist(err) {
        dialog := app.Dialog.Question().
            SetTitle("File Not Found").
            SetMessage("The file no longer exists. Remove from recent files?")

        remove := dialog.AddButton("Remove")
        remove.OnClick(func() {
            s.removeFromRecent(recentPath)
        })

        cancel := dialog.AddButton("Cancel")
        dialog.SetCancelButton(cancel)
        dialog.Show()

        return err
    }

    return s.openFile(recentPath)
}
```
