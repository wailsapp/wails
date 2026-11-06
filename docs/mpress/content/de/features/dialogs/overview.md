---
title: "Übersicht über Dialoge"
description: "Native Systemdialoge in Ihrer Anwendung anzeigen"
slug: "features/dialogs/overview"
sourcePath: "features/dialogs/overview.md"
---

## Native Dialoge

Wails bietet **native Systemdialoge**, die auf allen Plattformen funktionieren: Meldungsdialoge (Information, Warnung, Fehler, Frage), Dateidialoge (Öffnen, Speichern, Ordner) und benutzerdefinierte Dialogfenster mit plattformeigenem Erscheinungsbild und Verhalten.

![Ein Wails-Fragedialog unter macOS mit den Schaltflächen „Abbrechen“ und „Verwerfen“](/assets/screenshots/dialog-question-macos.png)

Dieselbe API verwendet die Konventionen jeder unterstützten Plattform. Dieses macOS-Beispiel zeigt einen Fragedialog, der seinem Wails-Fenster zugeordnet ist, einschließlich der Standard- und Abbruchschaltfläche.

## Schnellstart

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

**Das ist alles!** Native Dialoge mit minimalem Code.

## Auf Dialoge zugreifen

Der Zugriff auf Dialoge erfolgt über den `app.Dialog`-Manager:

```go
app.Dialog.Info()
app.Dialog.Question()
app.Dialog.Warning()
app.Dialog.Error()
app.Dialog.OpenFile()
app.Dialog.SaveFile()
```

## Dialogtypen

### Informationsdialog

Einfache Meldungen anzeigen:

```go
app.Dialog.Info().
    SetTitle("Welcome").
    SetMessage("Welcome to our application!").
    Show()
```

**Anwendungsfälle:**

- Erfolgsmeldungen
- Informationshinweise
- Abschlussbestätigungen

### Warndialog

Warnungen anzeigen:

```go
app.Dialog.Warning().
    SetTitle("Warning").
    SetMessage("This action cannot be undone.").
    Show()
```

**Anwendungsfälle:**

- Nicht kritische Warnungen
- Hinweise auf veraltete Funktionen
- Vorsichtshinweise

### Fehlerdialog

Fehler anzeigen:

```go
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Failed to save file: " + err.Error()).
    Show()
```

**Anwendungsfälle:**

- Fehlermeldungen
- Benachrichtigungen über fehlgeschlagene Vorgänge
- Ausnahmebehandlung

### Fragedialog

Benutzern Fragen stellen und Antworten über Schaltflächen-Callbacks verarbeiten:

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

**Anwendungsfälle:**

- Aktionen bestätigen
- Ja/Nein-Fragen
- Mehrfachauswahl

## Dateidialoge

### Dialog zum Öffnen von Dateien

Zu öffnende Dateien auswählen:

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

**Mehrfachauswahl:**

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

### Dialog zum Speichern von Dateien

Speicherort auswählen:

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

### Dialog zur Ordnerauswahl

Ein Verzeichnis über den Dialog zum Öffnen von Dateien auswählen, indem die Verzeichnisauswahl aktiviert wird:

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

## Dialogoptionen

### Titel und Meldung

```go
dialog := app.Dialog.Info().
    SetTitle("Success").
    SetMessage("Operation completed successfully!")
```

### Schaltflächen

**Standardschaltfläche für einfache Dialoge:**

Informations-, Warn- und Fehlerdialoge zeigen standardmäßig eine „OK“-Schaltfläche an:

```go
app.Dialog.Info().
    SetMessage("Done!").
    Show()
```

**Benutzerdefinierte Schaltflächen für Fragedialoge:**

Fügen Sie mit `AddButton()` Schaltflächen hinzu. Die Methode gibt ein `*Button` zurück, das Sie mit Callbacks konfigurieren können:

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

**Standard- und Abbruchschaltflächen:**

Legen Sie mit `SetDefaultButton()` fest, welche Schaltfläche hervorgehoben und durch die Eingabetaste ausgelöst wird. Legen Sie mit `SetCancelButton()` fest, welche Schaltfläche durch Escape ausgelöst wird.

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

### Zuordnung zu einem Fenster

Dialog einem bestimmten Fenster zuordnen:

```go
dialog := app.Dialog.Info().
    SetMessage("Window-specific message").
    AttachToWindow(window)

dialog.Show()
```

**Verhalten:**

- Der Dialog erscheint mittig über dem übergeordneten Fenster
- Das übergeordnete Fenster ist deaktiviert, solange der Dialog angezeigt wird
- Der Dialog bewegt sich mit dem übergeordneten Fenster (macOS)

## Plattformverhalten

@tabs{sync-key="platform"}
[macOS]
**macOS-Dialoge:**

- Natives NSAlert-Erscheinungsbild
- Verwenden das Systemdesign (hell/dunkel)
- Unterstützen die Tastaturnavigation
- Standardtastenkürzel (⌘. für Abbrechen)
- Integrierte Barrierefreiheitsfunktionen
- Bei Zuordnung zu einem Fenster im Sheet-Stil

**Beispiel:**

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
**Windows-Dialoge:**

- Natives TaskDialog-Erscheinungsbild
- Verwenden das Systemdesign
- Unterstützen die Tastaturnavigation
- Standardtastenkürzel (Esc für Abbrechen)
- Integrierte Barrierefreiheitsfunktionen
- Modal gegenüber dem übergeordneten Fenster

**Beispiel:**

```go
// Modal dialog on Windows
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Operation failed").
    Show()
```

[Linux]
**Linux-Dialoge:**

- GTK-Dialogerscheinungsbild
- Verwenden das Desktopdesign
- Unterstützen die Tastaturnavigation
- Integration in die Desktopumgebung
- Variiert je nach Desktopumgebung (GNOME, KDE usw.)

**Beispiel:**

```go
// GTK dialog on Linux
app.Dialog.Info().
    SetMessage("Update complete").
    Show()
```

@end

## Gängige Muster

### Vor destruktiven Aktionen bestätigen

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

### Fehlerbehandlung mit Dialogen

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

### Dateiauswahl mit Validierung

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

### Mehrstufiger Dialogablauf

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

## Bewährte Vorgehensweisen

### ✅ Empfohlen

- **Native Dialoge verwenden** – bessere UX als bei benutzerdefinierten Dialogen
- **Klare Meldungen formulieren** – präzise sein
- **Passende Titel festlegen** – der Kontext ist wichtig
- **Standardschaltflächen mit Bedacht verwenden** – die sichere Option als Standard festlegen
- **Abbrüche behandeln** – Benutzer könnten den Vorgang abbrechen
- **Dateiauswahlen validieren** – Dateitypen prüfen

### ❌ Nicht empfohlen

- **Dialoge nicht übermäßig verwenden** – sie unterbrechen den Arbeitsablauf
- **Nicht für häufige Meldungen verwenden** – stattdessen Benachrichtigungen verwenden
- **Fehlerbehandlung nicht vergessen** – Benutzer können den Vorgang abbrechen
- **Nicht unnötig blockieren** – Alternativen in Betracht ziehen
- **Keine allgemeinen Meldungen verwenden** – konkrete Angaben machen
- **Plattformunterschiede nicht ignorieren** – auf allen Plattformen testen

## Nächste Schritte

@cards{cols="2"}
ℹ Meldungsdialoge
Informations-, Warn- und Fehlerdialoge.

[Mehr erfahren →](/features/dialogs/message/)

---
📖 Dateidialoge
Dateien öffnen und speichern sowie Ordner auswählen.

[Mehr erfahren →](/features/dialogs/file/)

---
◆ Benutzerdefinierte Dialoge
Benutzerdefinierte Dialogfenster erstellen.

[Mehr erfahren →](/features/dialogs/custom/)

---
▣ Fenster
Mehr über die Fensterverwaltung erfahren.

[Mehr erfahren →](/features/windows/basics/)

@end

---

**Fragen?** Stelle sie auf [Discord](https://discord.gg/JDdSxwjhGf) oder sieh dir die [Dialogbeispiele](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/dialogs) an.
