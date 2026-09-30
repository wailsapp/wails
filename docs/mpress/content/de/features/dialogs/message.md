---
title: "Meldungsdialoge"
description: "Informationen, Warnungen, Fehler und Fragen anzeigen"
slug: "features/dialogs/message"
sourcePath: "features/dialogs/message.md"
---

## Meldungsdialoge

Wails bietet **native Meldungsdialoge** mit plattformgerechtem Erscheinungsbild: Informations-, Warn-, Fehler- und Fragedialoge mit anpassbaren Titeln, Meldungen und Schaltflächen. Einfache API, natives Verhalten und standardmäßig barrierefrei.

## Dialoge erstellen

Der Zugriff auf Meldungsdialoge erfolgt über den `app.Dialog`-Manager:

```go
app.Dialog.Info()
app.Dialog.Question()
app.Dialog.Warning()
app.Dialog.Error()
```

Alle Methoden geben ein `*MessageDialog` zurück, das sich durch Methodenverkettung konfigurieren lässt.

## Informationsdialog

Informationsmeldungen anzeigen:

```go
app.Dialog.Info().
    SetTitle("Success").
    SetMessage("File saved successfully!").
    Show()
```

**Anwendungsfälle:**

- Erfolgsbestätigungen
- Abschlussbenachrichtigungen
- Informationsmeldungen
- Statusaktualisierungen

**Beispiel – Speicherbestätigung:**

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

## Warndialog

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
- Warnhinweise
- Mögliche Probleme

**Beispiel – Warnung bei geringem Speicherplatz:**

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

## Fehlerdialog

Fehler anzeigen:

```go
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Failed to connect to server.").
    Show()
```

**Anwendungsfälle:**

- Fehlermeldungen
- Benachrichtigungen über fehlgeschlagene Vorgänge
- Ausnahmebehandlung
- Kritische Probleme

**Beispiel – Netzwerkfehler:**

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

## Fragedialog

Benutzern Fragen stellen und Antworten über Rückruffunktionen der Schaltflächen verarbeiten:

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

**Anwendungsfälle:**

- Aktionen bestätigen
- Ja/Nein-Fragen
- Auswahl zwischen mehreren Antworten
- Benutzerentscheidungen

**Beispiel – Nicht gespeicherte Änderungen:**

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

## Dialogoptionen

### Titel und Meldung

```go
dialog := app.Dialog.Info().
    SetTitle("Operation Complete").
    SetMessage("All files have been processed successfully.")
```

**Bewährte Vorgehensweisen:**

- **Titel:** Kurz und aussagekräftig (2-5 Wörter)
- **Meldung:** Klar, konkret und handlungsorientiert
- **Fachjargon vermeiden:** Verständliche Sprache verwenden

### Schaltflächen

**Einzelne Schaltfläche (Information/Warnung/Fehler):**

Informations-, Warn- und Fehlerdialoge zeigen standardmäßig eine „OK“-Schaltfläche an:

```go
app.Dialog.Info().
    SetMessage("Done!").
    Show()
```

Sie können auch benutzerdefinierte Schaltflächen hinzufügen:

```go
dialog := app.Dialog.Info().
    SetMessage("Done!")
ok := dialog.AddButton("Got it!")
dialog.SetDefaultButton(ok)
dialog.Show()
```

**Mehrere Schaltflächen (Frage):**

Fügen Sie mit `AddButton()` Schaltflächen hinzu. Die Methode gibt ein konfigurierbares `*Button` zurück:

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

**Standard- und Abbrechen-Schaltflächen:**

Legen Sie mit `SetDefaultButton()` fest, welche Schaltfläche hervorgehoben und mit der Eingabetaste ausgelöst wird. Legen Sie mit `SetCancelButton()` fest, welche Schaltfläche mit der Esc-Taste ausgelöst wird.

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

Sie können für Schaltflächen auch die Fluent-Methoden `SetAsDefault()` und `SetAsCancel()` verwenden:

```go
dialog := app.Dialog.Question().
    SetMessage("Delete file?")

dialog.AddButton("Delete").OnClick(func() {
    performDelete()
})

dialog.AddButton("Cancel").SetAsDefault().SetAsCancel()

dialog.Show()
```

**Bewährte Vorgehensweisen:**

- **1-3 Schaltflächen:** Benutzer nicht überfordern
- **Eindeutige Beschriftungen:** „Speichern“ statt „OK“
- **Sichere Standardeinstellung:** Nicht destruktive Aktion
- **Die Reihenfolge ist wichtig:** Die wahrscheinlichste Aktion zuerst (außer „Abbrechen“)

### Benutzerdefiniertes Symbol

Ein benutzerdefiniertes Symbol für den Dialog festlegen:

```go
app.Dialog.Info().
    SetTitle("Custom Icon Example").
    SetMessage("Using a custom icon").
    SetIcon(myIconBytes).
    Show()
```

### Fensterzuordnung

Den Dialog einem bestimmten Fenster zuordnen:

```go
dialog := app.Dialog.Question().
    SetMessage("Window-specific question").
    AttachToWindow(window)

dialog.AddButton("OK")
dialog.Show()
```

**Vorteile:**

- Der Dialog erscheint im richtigen Fenster
- Das übergeordnete Fenster ist während der Anzeige deaktiviert
- Bessere Benutzerfreundlichkeit bei mehreren Fenstern

## Vollständige Beispiele

### Destruktive Aktion bestätigen

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

### Beenden bestätigen

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

### Aktualisierungsdialog mit Downloadoption

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

### Fragedialog mit benutzerdefiniertem Symbol

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

## Bewährte Vorgehensweisen

### ✅ Empfohlen

- **Konkrete Angaben machen** – „Datei unter Dokumente gespeichert“ statt „Erfolgreich“
- **Passenden Typ verwenden** – „Error“ für Fehler, „Warning“ für Warnungen
- **Kontext angeben** – Relevante Details aufnehmen
- **Eindeutige Schaltflächenbeschriftungen verwenden** – „Löschen“ statt „OK“
- **Sichere Standardaktion festlegen** – Eine nicht destruktive Aktion
- **Abbruch behandeln** – Der Benutzer könnte den Dialog schließen

### ❌ Nicht empfohlen

- **Nicht übermäßig verwenden** – Unterbricht den Arbeitsablauf
- **Nicht für häufige Aktualisierungen verwenden** – Stattdessen Benachrichtigungen verwenden
- **Keine allgemeinen Meldungen verwenden** – „Fehler“ enthält keine hilfreichen Informationen
- **Fehler nicht ignorieren** – Fehler von dialog.Show() behandeln
- **Nicht unnötig blockieren** – Asynchrone Alternativen in Betracht ziehen
- **Keine Fachsprache verwenden** – Einfache Sprache verwenden

## Plattformunterschiede

### macOS

- Bei Verknüpfung mit einem Fenster im Sheet-Stil
- Standardtastenkürzel (⌘. für Abbrechen)
- Übernimmt automatisch das Systemdesign
- Integrierte Barrierefreiheit

### Windows

- Modale Dialoge
- TaskDialog-Darstellung
- Esc zum Abbrechen
- Übernimmt das Systemdesign

### Linux

- GTK-Dialoge
- Je nach Desktop-Umgebung unterschiedlich
- Übernimmt das Desktop-Design
- Standardmäßige Tastaturnavigation

## Nächste Schritte

@cards{cols="2"}
📖 Dateidialoge
Dateien öffnen und speichern sowie Ordner auswählen.

[Mehr erfahren →](/features/dialogs/file/)

---
◆ Benutzerdefinierte Dialoge
Benutzerdefinierte Dialogfenster erstellen.

[Mehr erfahren →](/features/dialogs/custom/)

---
● Benachrichtigungen
Nicht störende Benachrichtigungen.

[Mehr erfahren →](/features/notifications/overview/)

---
★ Ereignisse
Ereignisse für nicht blockierende Kommunikation verwenden.

[Mehr erfahren →](/features/events/system/)

@end

---

**Fragen?** Stelle sie auf [Discord](https://discord.gg/JDdSxwjhGf) oder sieh dir die [Dialogbeispiele](https://github.com/wailsapp/wails/tree/master/v3/examples/dialogs) an.
