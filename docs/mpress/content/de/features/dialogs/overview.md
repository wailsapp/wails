---
title: "Übersicht über Dialogfelder"
description: "Native Systemdialogfelder in Ihrer Anwendung anzeigen"
slug: "features/dialogs/overview"
sourcePath: "features/dialogs/overview.md"
---

## Native Dialogfelder

Wails bietet **native Systemdialogfelder**, die auf allen Plattformen funktionieren: Meldungsdialogfelder (Information, Warnung, Fehler, Frage), Dateidialogfelder (Öffnen, Speichern, Ordner) und benutzerdefinierte Dialogfenster mit plattformeigenem Erscheinungsbild und Verhalten.

![Ein Wails-Fragedialogfeld unter macOS mit den Schaltflächen „Abbrechen“ und „Verwerfen“](/assets/screenshots/dialog-question-macos.png)

Dieselbe API verwendet die Konventionen der jeweiligen unterstützten Plattform. Dieses macOS-Beispiel zeigt ein Fragedialogfeld, das an sein Wails-Fenster angehängt ist und die Standard- sowie die Abbrechen-Schaltfläche enthält.

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

**Das ist alles!** Native Dialogfelder mit minimalem Code.

## Auf Dialogfelder zugreifen

Der Zugriff auf Dialogfelder erfolgt über den `app.Dialog`-Manager:

```go
app.Dialog.Info()
app.Dialog.Question()
app.Dialog.Warning()
app.Dialog.Error()
app.Dialog.OpenFile()
app.Dialog.SaveFile()
```

## Dialogfeldtypen

### Informationsdialogfeld

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

### Warndialogfeld

Warnungen anzeigen:

```go
app.Dialog.Warning().
    SetTitle("Warning").
    SetMessage("This action cannot be undone.").
    Show()
```

**Anwendungsfälle:**

- Unkritische Warnungen
- Hinweise auf veraltete Funktionen
- Vorsichtshinweise

### Fehlerdialogfeld

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

### Fragedialogfeld

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
- Auswahl zwischen mehreren Antwortmöglichkeiten

## Dateidialogfelder

### Dialogfeld „Datei öffnen“

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

### Dialogfeld „Datei speichern“

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

### Dialogfeld „Ordner auswählen“

Ein Verzeichnis über das Dialogfeld „Datei öffnen“ mit aktivierter Verzeichnisauswahl auswählen:

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

## Dialogfeldoptionen

### Titel und Meldung

```go
dialog := app.Dialog.Info().
    SetTitle("Success").
    SetMessage("Operation completed successfully!")
```

### Schaltflächen

**Standardschaltfläche für einfache Dialogfelder:**

Informations-, Warn- und Fehlerdialogfelder zeigen standardmäßig eine „OK“-Schaltfläche an:

```go
app.Dialog.Info().
    SetMessage("Done!").
    Show()
```

**Benutzerdefinierte Schaltflächen für Fragedialogfelder:**

Verwenden Sie `AddButton()`, um Schaltflächen hinzuzufügen. Die Methode gibt ein konfigurierbares `*Button` zurück, für das Sie Callbacks festlegen können:

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

**Standard- und Abbrechen-Schaltflächen:**

Geben Sie mit `SetDefaultButton()` an, welche Schaltfläche hervorgehoben und durch die Eingabetaste ausgelöst wird. Geben Sie mit `SetCancelButton()` an, welche Schaltfläche durch die Escape-Taste ausgelöst wird.

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

### An ein Fenster anhängen

Dialogfeld an ein bestimmtes Fenster anhängen:

```go
dialog := app.Dialog.Info().
    SetMessage("Window-specific message").
    AttachToWindow(window)

dialog.Show()
```

**Verhalten:**

- Das Dialogfeld wird mittig über dem übergeordneten Fenster angezeigt
- Das übergeordnete Fenster ist deaktiviert, solange das Dialogfeld angezeigt wird
- Das Dialogfeld bewegt sich mit dem übergeordneten Fenster (macOS)

## Plattformverhalten

@tabs{sync-key="platform"}
[macOS]
**macOS-Dialogfelder:**

- Natives NSAlert-Erscheinungsbild
- Übernehmen das Systemdesign (hell/dunkel)
- Unterstützen die Tastaturnavigation
- Standardtastenkürzel (⌘. für „Abbrechen“)
- Integrierte Barrierefreiheitsfunktionen
- Bei Anbindung an ein Fenster im Sheet-Stil

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
**Windows-Dialogfelder:**

- Natives TaskDialog-Erscheinungsbild
- Übernehmen das Systemdesign
- Unterstützen die Tastaturnavigation
- Standardtastenkürzel (Esc für „Abbrechen“)
- Integrierte Barrierefreiheitsfunktionen
- Modal zum übergeordneten Fenster

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

- Erscheinungsbild von GTK-Dialogen
- Übernehmen das Desktop-Design
- Unterstützen die Tastaturnavigation
- Integration in die Desktop-Umgebung
- Variiert je nach Desktop-Umgebung (GNOME, KDE usw.)

**Beispiel:**

```go
// GTK dialog on Linux
app.Dialog.Info().
    SetMessage("Update complete").
    Show()
```

@end

## Gängige Muster

### Vor destruktiven Aktionen eine Bestätigung anfordern

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

- **Native Dialoge verwenden** – bessere Benutzerfreundlichkeit als benutzerdefinierte Dialoge
- **Klare Meldungen verwenden** – konkrete Angaben machen
- **Passende Titel festlegen** – der Kontext ist wichtig
- **Standardschaltflächen mit Bedacht verwenden** – die sichere Option als Standard festlegen
- **Abbrüche behandeln** – Benutzer können den Vorgang abbrechen
- **Dateiauswahl validieren** – Dateitypen prüfen

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

**Fragen?** Stelle sie auf [Discord](https://discord.gg/JDdSxwjhGf) oder sieh dir die [Dialogbeispiele](https://github.com/wailsapp/wails/tree/master/v3/examples/dialogs) an.
