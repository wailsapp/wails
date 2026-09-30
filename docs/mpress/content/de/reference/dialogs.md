---
title: "Dialogs-API"
description: "Vollständige Referenz für native Dialog-APIs"
slug: "reference/dialogs"
sourcePath: "reference/dialogs.md"
---

## Überblick

Die Dialogs-API bietet Methoden zum Anzeigen nativer Datei- und Meldungsdialoge. Greifen Sie über den Manager `app.Dialog` auf Dialoge zu.

**Dialogtypen:**

- **Dateidialoge** – Dialoge zum Öffnen und Speichern
- **Meldungsdialoge** – Informations-, Fehler-, Warn- und Fragedialoge

Alle Dialoge sind **native Betriebssystemdialoge**, deren Erscheinungsbild der jeweiligen Plattform entspricht.

## Auf Dialoge zugreifen

Der Zugriff auf Dialoge erfolgt über den Manager `app.Dialog`:

```go
app.Dialog.OpenFile()
app.Dialog.SaveFile()
app.Dialog.Info()
app.Dialog.Question()
app.Dialog.Warning()
app.Dialog.Error()
```

## Dateidialoge

### OpenFile()

Erstellt einen Dialog zum Öffnen einer Datei.

```go
func (dm *DialogManager) OpenFile() *OpenFileDialogStruct
```

**Beispiel:**

```go
dialog := app.Dialog.OpenFile()
```

### Methoden von OpenFileDialogStruct

#### SetTitle()

Legt den Titel des Dialogs fest.

```go
func (d *OpenFileDialogStruct) SetTitle(title string) *OpenFileDialogStruct
```

**Beispiel:**

```go
dialog.SetTitle("Select Image")
```

#### AddFilter()

Fügt einen Dateitypfilter hinzu.

```go
func (d *OpenFileDialogStruct) AddFilter(displayName, pattern string) *OpenFileDialogStruct
```

**Parameter:**

- `displayName` – Dem Benutzer angezeigte Filterbeschreibung (z. B. „Bilder“, „Dokumente“)
- `pattern` – Durch Semikolons getrennte Liste von Dateiendungen (z. B. „*.png;*.jpg“)

**Beispiel:**

```go
dialog.AddFilter("Images", "*.png;*.jpg;*.gif").
    AddFilter("Documents", "*.pdf;*.docx").
    AddFilter("All Files", "*.*")
```

#### SetDirectory()

Legt das anfängliche Verzeichnis fest.

```go
func (d *OpenFileDialogStruct) SetDirectory(directory string) *OpenFileDialogStruct
```

**Beispiel:**

```go
homeDir, _ := os.UserHomeDir()
dialog.SetDirectory(homeDir)
```

#### CanChooseDirectories()

Aktiviert oder deaktiviert die Auswahl von Verzeichnissen.

```go
func (d *OpenFileDialogStruct) CanChooseDirectories(canChooseDirectories bool) *OpenFileDialogStruct
```

**Beispiel (Ordnerauswahl):**

```go
// Select folders instead of files
path, err := app.Dialog.OpenFile().
    SetTitle("Select Folder").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()
```

#### CanChooseFiles()

Aktiviert oder deaktiviert die Auswahl von Dateien.

```go
func (d *OpenFileDialogStruct) CanChooseFiles(canChooseFiles bool) *OpenFileDialogStruct
```

#### CanCreateDirectories()

Aktiviert oder deaktiviert das Erstellen neuer Verzeichnisse.

```go
func (d *OpenFileDialogStruct) CanCreateDirectories(canCreateDirectories bool) *OpenFileDialogStruct
```

#### ShowHiddenFiles()

Blendet ausgeblendete Dateien ein oder aus.

```go
func (d *OpenFileDialogStruct) ShowHiddenFiles(showHiddenFiles bool) *OpenFileDialogStruct
```

#### AttachToWindow()

Verknüpft den Dialog mit einem bestimmten Fenster.

```go
func (d *OpenFileDialogStruct) AttachToWindow(window Window) *OpenFileDialogStruct
```

#### PromptForSingleSelection()

Zeigt den Dialog an und gibt die ausgewählte Datei zurück.

```go
func (d *OpenFileDialogStruct) PromptForSingleSelection() (string, error)
```

**Rückgabewerte:**

- `string` – Pfad der ausgewählten Datei. Behandeln Sie eine leere Zeichenfolge als „keine Auswahl“ (je nach Betriebssystem geben die Plattformimplementierungen beim Abbrechen entweder eine leere Zeichenfolge oder einen Fehler ungleich nil zurück).
- `error` – Ungleich nil, wenn der Dialog selbst nicht angezeigt werden konnte.

**Beispiel:**

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

Zeigt den Dialog an und gibt mehrere ausgewählte Dateien zurück.

```go
func (d *OpenFileDialogStruct) PromptForMultipleSelection() ([]string, error)
```

**Rückgabewerte:**

- `[]string` – Array mit den Pfaden der ausgewählten Dateien
- `error` – Fehler, wenn der Dialog fehlgeschlagen ist

**Beispiel:**

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

Erstellt einen Dialog zum Speichern einer Datei.

```go
func (dm *DialogManager) SaveFile() *SaveFileDialogStruct
```

**Beispiel:**

```go
dialog := app.Dialog.SaveFile()
```

### Methoden von SaveFileDialogStruct

#### SetTitle()

Legt den Titel des Dialogs fest.

```go
func (d *SaveFileDialogStruct) SetTitle(title string) *SaveFileDialogStruct
```

#### SetFilename()

Legt den standardmäßigen Dateinamen fest.

```go
func (d *SaveFileDialogStruct) SetFilename(filename string) *SaveFileDialogStruct
```

**Beispiel:**

```go
dialog.SetFilename("document.pdf")
```

#### AddFilter()

Fügt einen Dateitypfilter hinzu.

```go
func (d *SaveFileDialogStruct) AddFilter(displayName, pattern string) *SaveFileDialogStruct
```

**Beispiel:**

```go
dialog.AddFilter("PDF Document", "*.pdf").
    AddFilter("Text Document", "*.txt")
```

#### SetDirectory()

Legt das anfängliche Verzeichnis fest.

```go
func (d *SaveFileDialogStruct) SetDirectory(directory string) *SaveFileDialogStruct
```

#### AttachToWindow()

Verknüpft den Dialog mit einem bestimmten Fenster.

```go
func (d *SaveFileDialogStruct) AttachToWindow(window Window) *SaveFileDialogStruct
```

#### PromptForSingleSelection()

Zeigt den Dialog an und gibt den Speicherpfad zurück.

```go
func (d *SaveFileDialogStruct) PromptForSingleSelection() (string, error)
```

**Beispiel:**

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

### Ordnerauswahl

Es gibt keine separate `SelectFolderDialog`. Verwenden Sie `OpenFile()` mit Verzeichnisoptionen:

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

## Meldungsdialoge

Alle Meldungsdialoge geben `*MessageDialog` zurück und stellen dieselben Methoden bereit.

### Info()

Erstellt einen Informationsdialog.

```go
func (dm *DialogManager) Info() *MessageDialog
```

**Beispiel:**

```go
app.Dialog.Info().
    SetTitle("Success").
    SetMessage("File saved successfully!").
    Show()
```

### Error()

Erstellt einen Fehlerdialog.

```go
func (dm *DialogManager) Error() *MessageDialog
```

**Beispiel:**

```go
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Failed to save file: " + err.Error()).
    Show()
```

### Warning()

Erstellt einen Warnungsdialog.

```go
func (dm *DialogManager) Warning() *MessageDialog
```

**Beispiel:**

```go
app.Dialog.Warning().
    SetTitle("Warning").
    SetMessage("This action cannot be undone.").
    Show()
```

### Question()

Erstellt einen Fragedialog mit benutzerdefinierten Schaltflächen.

```go
func (dm *DialogManager) Question() *MessageDialog
```

**Beispiel:**

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

### MessageDialog-Methoden

#### SetTitle()

Legt den Titel des Dialogs fest.

```go
func (d *MessageDialog) SetTitle(title string) *MessageDialog
```

#### SetMessage()

Legt die Meldung des Dialogs fest.

```go
func (d *MessageDialog) SetMessage(message string) *MessageDialog
```

#### SetIcon()

Legt ein benutzerdefiniertes Symbol für den Dialog fest.

```go
func (d *MessageDialog) SetIcon(icon []byte) *MessageDialog
```

#### AddButton()

Fügt dem Dialog eine Schaltfläche hinzu und gibt sie zur Konfiguration zurück.

```go
func (d *MessageDialog) AddButton(label string) *Button
```

**Rückgabewert:** `*Button` – Die Schaltflächeninstanz zur weiteren Konfiguration

**Beispiel:**

```go
button := dialog.AddButton("OK")
button.OnClick(func() {
    // Handle click
})
```

#### SetDefaultButton()

Legt fest, welche Schaltfläche die Standardschaltfläche ist (wird durch Drücken der Eingabetaste aktiviert).

```go
func (d *MessageDialog) SetDefaultButton(button *Button) *MessageDialog
```

**Beispiel:**

```go
yes := dialog.AddButton("Yes")
no := dialog.AddButton("No")
dialog.SetDefaultButton(yes)
```

#### SetCancelButton()

Legt fest, welche Schaltfläche die Abbrechen-Schaltfläche ist (wird durch Drücken der Esc-Taste aktiviert).

```go
func (d *MessageDialog) SetCancelButton(button *Button) *MessageDialog
```

**Beispiel:**

```go
ok := dialog.AddButton("OK")
cancel := dialog.AddButton("Cancel")
dialog.SetCancelButton(cancel)
```

#### AttachToWindow()

Ordnet den Dialog einem bestimmten Fenster zu.

```go
func (d *MessageDialog) AttachToWindow(window Window) *MessageDialog
```

#### Show()

Zeigt den Dialog an. Schaltflächen-Callbacks verarbeiten die Benutzerreaktionen.

```go
func (d *MessageDialog) Show()
```

**Hinweis:** `Show()` gibt keinen Wert zurück. Verwenden Sie Schaltflächen-Callbacks, um Benutzerreaktionen zu verarbeiten.

### Schaltflächenmethoden

#### OnClick()

Legt die Callback-Funktion fest, die beim Klicken auf die Schaltfläche aufgerufen wird.

```go
func (b *Button) OnClick(callback func()) *Button
```

#### SetAsDefault()

Kennzeichnet diese Schaltfläche als Standardschaltfläche.

```go
func (b *Button) SetAsDefault() *Button
```

#### SetAsCancel()

Kennzeichnet diese Schaltfläche als Abbrechen-Schaltfläche.

```go
func (b *Button) SetAsCancel() *Button
```

## Vollständige Beispiele

### Beispiel für die Dateiauswahl

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

### Beispiel für einen Bestätigungsdialog

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

### Dialog zum Speichern von Änderungen

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

### Verarbeitung mehrerer Dateien

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

### Fehlerbehandlung mit Dialogen

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

### Plattformspezifische Standardwerte

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

## Bewährte Vorgehensweisen

### Empfohlen

- **Native Dialoge verwenden** – Sie entsprechen dem Erscheinungsbild der Plattform
- **Eindeutige Titel verwenden** – Sie helfen Benutzern, den Zweck zu verstehen
- **Geeignete Filter festlegen** – Sie führen Benutzer zu den richtigen Dateitypen
- **Abbrüche behandeln** – Auf Fehler prüfen (der Benutzer kann abbrechen)
- **Bei destruktiven Aktionen eine Bestätigung anzeigen** – Question-Dialoge verwenden
- **Rückmeldung geben** – Info-Dialoge für Erfolgsmeldungen verwenden
- **Sinnvolle Standardwerte festlegen** – Standardverzeichnis, Dateiname usw.
- **Callbacks für Schaltflächenaktionen verwenden** – Benutzerreaktionen korrekt verarbeiten

### Nicht empfohlen

- **Fehler nicht ignorieren** – Ein Abbruch durch den Benutzer gibt einen Fehler zurück
- **Keine mehrdeutigen Schaltflächenbeschriftungen verwenden** – Eindeutige Bezeichnungen verwenden: „Speichern“/„Abbrechen“
- **Dialoge nicht übermäßig verwenden** – sie unterbrechen den Arbeitsablauf
- **Bei einem Abbruch keine Fehler anzeigen** – er ist eine normale Aktion
- **Dateifilter nicht vergessen** – sie helfen Benutzern, die richtigen Dateien zu finden
- **Pfade nicht fest codieren** – os.UserHomeDir() oder Ähnliches verwenden

## Dialogtypen nach Plattform

### macOS

- Dialoge werden von der Titelleiste nach unten eingeblendet
- Als „Sheet“ an das übergeordnete Fenster angeheftet
- Natives macOS-Erscheinungsbild

### Windows

- Windows-Standarddialoge
- Entspricht den Windows-Designrichtlinien
- Modernes Windows-10/11-Erscheinungsbild

### Linux

- GTK-Dialoge auf GTK-basierten Systemen
- Qt-Dialoge auf Qt-basierten Systemen
- Passt sich an die Desktop-Umgebung an

#### Verhalten von Dialogen unter Linux

Unter Linux verwendet der standardmäßige GTK4-Build für Dateidialoge **xdg-desktop-portal**. Dies ermöglicht eine native Desktop-Integration, bedeutet jedoch, dass einige Optionen wirkungslos sind. Der ältere GTK3-Pfad (`-tags gtk3`) behält die vollständige programmgesteuerte Kontrolle über diese Optionen:

| Option | GTK3 (`-tags gtk3`) | GTK4 (Standard) | Hinweise |
| --- | --- | --- | --- |
| `ShowHiddenFiles()` | ✅ Funktioniert | ❌ Keine Wirkung | Benutzer steuern dies über den Schalter in der Benutzeroberfläche des Dialogs (Ctrl+H oder Menü) |
| `CanCreateDirectories()` | ✅ Funktioniert | ❌ Keine Wirkung | Im Portal immer aktiviert |
| `ResolvesAliases()` | ✅ Funktioniert | ❌ Keine Wirkung | Das Portal übernimmt die Auflösung symbolischer Links |
| `SetButtonText()` | ✅ Funktioniert | ✅ Funktioniert | Benutzerdefinierter Text für die Bestätigungsschaltfläche funktioniert |

**Warum diese Einschränkungen bestehen:** Die portalbasierten Dialoge von GTK4 übertragen die Kontrolle über die Benutzeroberfläche an die Desktop-Umgebung (GNOME, KDE usw.). Dies ist beabsichtigt: Das Portal sorgt anwendungsübergreifend für eine einheitliche Benutzererfahrung und berücksichtigt die Benutzereinstellungen.

@note{type="info"}
Der standardmäßige GTK4-Build verwendet portalbasierte Dialoge. Wenn Ihre Anwendung die vollständige programmgesteuerte Kontrolle über die oben genannten Dialogoptionen benötigt, erstellen Sie den Build mit dem älteren `-tags gtk3`-Pfad (unterstützt bis einschließlich v3.0.x; entfernt in v3.1) – siehe [Linux-Paketierung – Unterstützung für das ältere GTK3](/guides/build/linux/#legacy-gtk3-support).

@end

## Gängige Muster

### Muster „Speichern unter“

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

### Muster „Zuletzt verwendet öffnen“

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
