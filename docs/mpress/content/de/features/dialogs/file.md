---
title: "Dateidialoge"
description: "Dialoge zum Öffnen, Speichern und Auswählen von Ordnern"
slug: "features/dialogs/file"
sourcePath: "features/dialogs/file.md"
---

## Dateidialoge

Wails bietet **native Dateidialoge** mit plattformgerechtem Erscheinungsbild zum Öffnen und Speichern von Dateien sowie zum Auswählen von Ordnern. Die einfache API unterstützt Dateitypfilter, Mehrfachauswahl und Standardverzeichnisse.

![Eine von einer Wails-Anwendung geöffnete native macOS-Dateiauswahl](/assets/screenshots/file-dialog-macos.png)

Die Wails-API übergibt die Auswahl an das Betriebssystem, sodass die vertraute Navigation sowie das bekannte Filter- und Auswahlverhalten erhalten bleiben.

## Dateidialoge erstellen

Der Zugriff auf Dateidialoge erfolgt über den `app.Dialog`-Manager:

```go
app.Dialog.OpenFile()
app.Dialog.SaveFile()
```

## Dialog zum Öffnen von Dateien

Wählen Sie zu öffnende Dateien aus:

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

**Anwendungsfälle:**

- Dokumente öffnen
- Dateien importieren
- Bilder laden
- Konfigurationsdateien auswählen

### Einzelne Datei auswählen

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

### Mehrere Dateien auswählen

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

### Mit Standardverzeichnis

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Open File").
    SetDirectory("/Users/me/Documents").
    PromptForSingleSelection()
```

## Dialog zum Speichern von Dateien

Wählen Sie den Speicherort aus:

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

**Anwendungsfälle:**

- Dokumente speichern
- Daten exportieren
- Neue Dateien erstellen
- Speichern unter …

### Mit Standarddateinamen

```go
path, err := app.Dialog.SaveFile().
    SetFilename("export.csv").
    AddFilter("CSV Files", "*.csv").
    PromptForSingleSelection()
```

### Mit Standardverzeichnis

```go
path, err := app.Dialog.SaveFile().
    SetDirectory("/Users/me/Documents").
    SetFilename("untitled.txt").
    PromptForSingleSelection()
```

### Überschreiben bestätigen

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

## Dialog zur Ordnerauswahl

Wählen Sie ein Verzeichnis über den Dialog zum Öffnen von Dateien aus, in dem die Verzeichnisauswahl aktiviert ist:

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

**Anwendungsfälle:**

- Ausgabeverzeichnis auswählen
- Arbeitsbereich auswählen
- Speicherort für Sicherungen auswählen
- Installationsverzeichnis auswählen

### Mit Standardverzeichnis

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Folder").
    SetDirectory("/Users/me/Documents").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()
```

## Dateifilter

Fügen Sie Dialogen mit der Methode `AddFilter()` Dateitypfilter hinzu. Jeder Aufruf fügt eine neue Filteroption hinzu.

### Grundlegende Filter

```go
path, _ := app.Dialog.OpenFile().
    AddFilter("Text Files", "*.txt").
    AddFilter("All Files", "*.*").
    PromptForSingleSelection()
```

### Mehrere Dateierweiterungen

Trennen Sie mehrere Dateierweiterungen innerhalb eines Filters durch Semikolons:

```go
dialog := app.Dialog.OpenFile().
    AddFilter("Images", "*.png;*.jpg;*.jpeg;*.gif;*.bmp").
    AddFilter("Documents", "*.txt;*.doc;*.docx;*.pdf").
    AddFilter("All Files", "*.*")
```

### Musterformat

Trennen Sie mehrere Dateierweiterungen innerhalb eines Filters durch **Semikolons**:

```go
// Multiple extensions separated by semicolons
AddFilter("Images", "*.png;*.jpg;*.gif")
```

## Vollständige Beispiele

### Bilddatei öffnen

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

### Dokument mit Validierung speichern

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

### Stapelverarbeitung von Dateien

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

### Mit Ordnerauswahl exportieren

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

### Mit Validierung importieren

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

## Bewährte Vorgehensweisen

### ✅ Empfohlen

- **Dateifilter bereitstellen** – Erleichtert das Auffinden von Dateien
- **Passende Titel festlegen** – Sorgt für einen klaren Kontext
- **Standardverzeichnisse verwenden** – Beginnt an einem sinnvollen Speicherort
- **Auswahl validieren** – Prüft die Dateitypen
- **Abbruch behandeln** – Benutzer können den Vorgang abbrechen
- **Bestätigung anzeigen** – Bei destruktiven Aktionen
- **Rückmeldung geben** – Erfolgs- und Fehlermeldungen

### ❌ Nicht tun

- **Validierung nicht überspringen** – Dateitypen prüfen
- **Fehler nicht ignorieren** – Abbruch behandeln
- **Keine allgemeinen Filter verwenden** – Konkrete Filter angeben
- **„Alle Dateien“ nicht vergessen** – Immer als Option anbieten
- **Pfade nicht fest codieren** – Home-Verzeichnis des Benutzers verwenden
- **Nicht voraussetzen, dass die Datei existiert** – Vor dem Öffnen prüfen

## Plattformunterschiede

### macOS

- Native NSOpenPanel/NSSavePanel
- Bei Anbindung an ein Fenster im Sheet-Stil
- Übernimmt das Systemdesign
- Unterstützt die Quick-Look-Vorschau
- Integration von Tags und Favoriten

### Windows

- Native Dialoge zum Öffnen und Speichern von Dateien
- Übernimmt das Systemdesign
- Integration zuletzt verwendeter Dateien
- Unterstützung für Netzwerkspeicherorte

### Linux

- GTK-Dateiauswahl
- Je nach Desktop-Umgebung unterschiedlich
- Übernimmt das Desktop-Design
- Unterstützung für zuletzt verwendete Dateien

## Nächste Schritte

@cards{cols="2"}
ℹ Meldungsdialoge
Informations-, Warn- und Fehlerdialoge.

[Mehr erfahren →](/features/dialogs/message/)

---
◆ Benutzerdefinierte Dialoge
Erstellen Sie benutzerdefinierte Dialogfenster.

[Mehr erfahren →](/features/dialogs/custom/)

---
🚀 Bindings
Rufen Sie Go-Funktionen aus JavaScript auf.

[Mehr erfahren →](/features/bindings/methods/)

---
★ Ereignisse
Verwenden Sie Ereignisse für Fortschrittsmeldungen.

[Mehr erfahren →](/features/events/system/)

@end

---

**Fragen?** Stellen Sie sie auf [Discord](https://discord.gg/JDdSxwjhGf) oder sehen Sie sich die [Beispiele für Dateidialoge](https://github.com/wailsapp/wails/tree/master/v3/examples/dialogs) an.
