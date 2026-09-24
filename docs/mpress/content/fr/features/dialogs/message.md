---
title: "Boîtes de dialogue de message"
description: "Afficher des informations, des avertissements, des erreurs et des questions"
slug: "features/dialogs/message"
sourcePath: "features/dialogs/message.md"
---

## Boîtes de dialogue de message

Wails fournit des **boîtes de dialogue de message natives** à l’apparence adaptée à chaque plateforme : boîtes de dialogue d’information, d’avertissement, d’erreur et de question, avec des titres, des messages et des boutons personnalisables. API simple, comportement natif et accessibilité par défaut.

## Créer des boîtes de dialogue

Les boîtes de dialogue de message sont accessibles par l’intermédiaire du gestionnaire `app.Dialog` :

```go
app.Dialog.Info()
app.Dialog.Question()
app.Dialog.Warning()
app.Dialog.Error()
```

Toutes les méthodes renvoient un objet `*MessageDialog` qui peut être configuré par chaînage de méthodes.

## Boîte de dialogue d’information

Affichez des messages d’information :

```go
app.Dialog.Info().
    SetTitle("Success").
    SetMessage("File saved successfully!").
    Show()
```

**Cas d’utilisation :**

- Confirmations de réussite
- Notifications d’achèvement
- Messages d’information
- Mises à jour d’état

**Exemple — Confirmation d’enregistrement :**

```go
func saveFile(app *application.App, path string, data []byte) error {
    if err := os.WriteFile(path, data, 0644); err != nil {
        return err
    }

    app.Dialog.Info().
        SetTitle("File Saved").
        SetMessage(fmt.Sprintf("Saved to %s", filepath.Base(path))).
        Show()

    return nil
}
```

## Boîte de dialogue d’avertissement

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
- Problèmes potentiels

**Exemple — Avertissement relatif à l’espace disque :**

```go
func checkDiskSpace(app *application.App) {
    available := getDiskSpace()

    if available < 100*1024*1024 { // Less than 100MB
        app.Dialog.Warning().
            SetTitle("Low Disk Space").
            SetMessage(fmt.Sprintf("Only %d MB available.", available/(1024*1024))).
            Show()
    }
}
```

## Boîte de dialogue d’erreur

Affichez des erreurs :

```go
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Failed to connect to server.").
    Show()
```

**Cas d’utilisation :**

- Messages d’erreur
- Notifications d’échec
- Gestion des exceptions
- Problèmes critiques

**Exemple — Erreur réseau :**

```go
func fetchData(app *application.App, url string) ([]byte, error) {
    resp, err := http.Get(url)
    if err != nil {
        app.Dialog.Error().
            SetTitle("Network Error").
            SetMessage(fmt.Sprintf("Failed to connect: %v", err)).
            Show()
        return nil, err
    }
    defer resp.Body.Close()

    return io.ReadAll(resp.Body)
}
```

## Boîte de dialogue de question

Posez des questions aux utilisateurs et traitez leurs réponses à l’aide de fonctions de rappel associées aux boutons :

```go
dialog := app.Dialog.Question().
    SetTitle("Confirm").
    SetMessage("Save changes before closing?")

save := dialog.AddButton("Save")
save.OnClick(func() {
    saveChanges()
})

dontSave := dialog.AddButton("Don't Save")
dontSave.OnClick(func() {
    // Continue without saving
})

cancel := dialog.AddButton("Cancel")
cancel.OnClick(func() {
    // Don't close
})

dialog.SetDefaultButton(save)
dialog.SetCancelButton(cancel)
dialog.Show()
```

**Cas d’utilisation :**

- Confirmation d’actions
- Questions Oui/Non
- Choix multiple
- Décisions de l’utilisateur

**Exemple — Modifications non enregistrées :**

```go
func closeDocument(app *application.App) {
    if !hasUnsavedChanges() {
        doClose()
        return
    }

    dialog := app.Dialog.Question().
        SetTitle("Unsaved Changes").
        SetMessage("Do you want to save your changes?")

    save := dialog.AddButton("Save")
    save.OnClick(func() {
        if saveDocument() {
            doClose()
        }
    })

    dontSave := dialog.AddButton("Don't Save")
    dontSave.OnClick(func() {
        doClose()
    })

    cancel := dialog.AddButton("Cancel")
    // Cancel button has no callback - just closes the dialog

    dialog.SetDefaultButton(save)
    dialog.SetCancelButton(cancel)
    dialog.Show()
}
```

## Options de la boîte de dialogue

### Titre et message

```go
dialog := app.Dialog.Info().
    SetTitle("Operation Complete").
    SetMessage("All files have been processed successfully.")
```

**Bonnes pratiques :**

- **Titre :** court et descriptif (2-5 mots)
- **Message :** clair, précis et exploitable
- **Évitez le jargon :** employez un langage simple

### Boutons

**Bouton unique (information/avertissement/erreur) :**

Les boîtes de dialogue d’information, d’avertissement et d’erreur affichent par défaut un bouton « OK » :

```go
app.Dialog.Info().
    SetMessage("Done!").
    Show()
```

Vous pouvez également ajouter des boutons personnalisés :

```go
dialog := app.Dialog.Info().
    SetMessage("Done!")
ok := dialog.AddButton("Got it!")
dialog.SetDefaultButton(ok)
dialog.Show()
```

**Plusieurs boutons (question) :**

Utilisez `AddButton()` pour ajouter des boutons. Cette méthode renvoie un objet `*Button` que vous pouvez configurer :

```go
dialog := app.Dialog.Question().
    SetMessage("Choose an action")

option1 := dialog.AddButton("Option 1")
option1.OnClick(func() {
    handleOption1()
})

option2 := dialog.AddButton("Option 2")
option2.OnClick(func() {
    handleOption2()
})

option3 := dialog.AddButton("Option 3")
option3.OnClick(func() {
    handleOption3()
})

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
// No callback needed - just dismisses dialog

dialog.SetDefaultButton(cancelBtn)  // Safe option as default
dialog.SetCancelButton(cancelBtn)   // Escape triggers Cancel
dialog.Show()
```

Vous pouvez également utiliser les méthodes fluides `SetAsDefault()` et `SetAsCancel()` sur les boutons :

```go
dialog := app.Dialog.Question().
    SetMessage("Delete file?")

dialog.AddButton("Delete").OnClick(func() {
    performDelete()
})

dialog.AddButton("Cancel").SetAsDefault().SetAsCancel()

dialog.Show()
```

**Bonnes pratiques :**

- **1-3 boutons :** ne submergez pas les utilisateurs
- **Libellés explicites :** « Enregistrer » plutôt que « OK »
- **Choix par défaut sûr :** une action non destructive
- **L’ordre est important :** placez l’action la plus probable en premier (sauf Annuler)

### Icône personnalisée

Définissez une icône personnalisée pour la boîte de dialogue :

```go
app.Dialog.Info().
    SetTitle("Custom Icon Example").
    SetMessage("Using a custom icon").
    SetIcon(myIconBytes).
    Show()
```

### Association à une fenêtre

Associez la boîte de dialogue à une fenêtre précise :

```go
dialog := app.Dialog.Question().
    SetMessage("Window-specific question").
    AttachToWindow(window)

dialog.AddButton("OK")
dialog.Show()
```

**Avantages :**

- La boîte de dialogue apparaît dans la fenêtre appropriée
- La fenêtre parente est désactivée tant que la boîte de dialogue est affichée
- Meilleure expérience utilisateur avec plusieurs fenêtres

## Exemples complets

### Confirmer une action destructive

```go
func deleteFiles(app *application.App, paths []string) {
    // Confirm deletion
    message := fmt.Sprintf("Delete %d file(s)?", len(paths))
    if len(paths) == 1 {
        message = fmt.Sprintf("Delete %s?", filepath.Base(paths[0]))
    }

    dialog := app.Dialog.Question().
        SetTitle("Confirm Delete").
        SetMessage(message)

    deleteBtn := dialog.AddButton("Delete")
    deleteBtn.OnClick(func() {
        // Perform deletion
        var errs []error
        for _, path := range paths {
            if err := os.Remove(path); err != nil {
                errs = append(errs, err)
            }
        }

        // Show result
        if len(errs) > 0 {
            app.Dialog.Error().
                SetTitle("Delete Failed").
                SetMessage(fmt.Sprintf("Failed to delete %d file(s)", len(errs))).
                Show()
        } else {
            app.Dialog.Info().
                SetTitle("Delete Complete").
                SetMessage(fmt.Sprintf("Deleted %d file(s)", len(paths))).
                Show()
        }
    })

    cancelBtn := dialog.AddButton("Cancel")
    dialog.SetDefaultButton(cancelBtn)
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
}
```

### Confirmation avant de quitter

```go
func confirmQuit(app *application.App) {
    dialog := app.Dialog.Question().
        SetTitle("Quit").
        SetMessage("You have unsaved work. Are you sure you want to quit?")

    yes := dialog.AddButton("Yes")
    yes.OnClick(func() {
        app.Quit()
    })

    no := dialog.AddButton("No")
    dialog.SetDefaultButton(no)
    dialog.Show()
}
```

### Boîte de dialogue de mise à jour avec option de téléchargement

```go
func showUpdateDialog(app *application.App) {
    dialog := app.Dialog.Question().
        SetTitle("Update").
        SetMessage("A new version is available. The cancel button is selected when pressing escape.")

    download := dialog.AddButton("📥 Download")
    download.OnClick(func() {
        app.Dialog.Info().SetMessage("Downloading...").Show()
    })

    cancel := dialog.AddButton("Cancel")

    dialog.SetDefaultButton(download)
    dialog.SetCancelButton(cancel)
    dialog.Show()
}
```

### Question avec une icône personnalisée

```go
func showCustomIconQuestion(app *application.App, iconBytes []byte) {
    dialog := app.Dialog.Question().
        SetTitle("Custom Icon Example").
        SetMessage("Using a custom icon").
        SetIcon(iconBytes)

    likeIt := dialog.AddButton("I like it!")
    likeIt.OnClick(func() {
        app.Dialog.Info().SetMessage("Thanks!").Show()
    })

    notKeen := dialog.AddButton("Not so keen...")
    notKeen.OnClick(func() {
        app.Dialog.Info().SetMessage("Too bad!").Show()
    })

    dialog.SetDefaultButton(likeIt)
    dialog.Show()
}
```

## Bonnes pratiques

### ✅ À faire

- **Soyez précis** — « Fichier enregistré dans Documents » plutôt que « Réussite »
- **Utilisez le type approprié** — Error pour les erreurs, Warning pour les avertissements
- **Fournissez le contexte** — Incluez les détails pertinents
- **Utilisez des libellés de boutons explicites** — « Supprimer » plutôt que « OK »
- **Définissez des choix par défaut sûrs** — Choisissez une action non destructive
- **Gérez l’annulation** — L’utilisateur peut fermer la boîte de dialogue

### ❌ À éviter

- **N’en abusez pas** — Cela interrompt le flux de travail
- **Ne les utilisez pas pour les mises à jour fréquentes** — Utilisez plutôt des notifications
- **N’utilisez pas de messages génériques** — « Erreur » n’apporte aucune information
- **N’ignorez pas les erreurs** — Gérez les erreurs de dialog.Show()
- **Ne bloquez pas inutilement** — Envisagez des solutions asynchrones
- **N’utilisez pas de jargon technique** — Employez un langage clair

## Différences entre les plateformes

### macOS

- Présentation sous forme de feuille lorsque la boîte de dialogue est attachée à une fenêtre
- Raccourcis clavier standard (⌘. pour annuler)
- Adopte automatiquement le thème du système
- Fonctions d’accessibilité intégrées

### Windows

- Boîtes de dialogue modales
- Apparence de TaskDialog
- Touche Échap pour annuler
- Adopte le thème du système

### Linux

- Boîtes de dialogue GTK
- Varie selon l’environnement de bureau
- Adopte le thème de l’environnement de bureau
- Navigation standard au clavier

## Étapes suivantes

@cards{cols="2"}
📖 Boîtes de dialogue de fichiers
Ouverture, enregistrement et sélection de dossiers.

[En savoir plus →](/features/dialogs/file/)

---
◆ Boîtes de dialogue personnalisées
Créez des fenêtres de dialogue personnalisées.

[En savoir plus →](/features/dialogs/custom/)

---
● Notifications
Notifications non intrusives.

[En savoir plus →](/features/notifications/overview/)

---
★ Événements
Utilisez des événements pour une communication non bloquante.

[En savoir plus →](/features/events/system/)

@end

---

**Des questions ?** Posez-les sur [Discord](https://discord.gg/JDdSxwjhGf) ou consultez les [exemples de boîtes de dialogue](https://github.com/wailsapp/wails/tree/master/v3/examples/dialogs).
