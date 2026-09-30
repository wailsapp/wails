---
title: "Dateizuordnungen"
description: "Dateizuordnungen für Ihre Wails-Anwendung konfigurieren"
slug: "guides/file-associations"
sourcePath: "guides/file-associations.md"
---

Relevante Plattformen: <span class="mpress-badge mpress-badge-note">Windows</span> <span class="mpress-badge mpress-badge-success">macOS</span>

<br/>

Mit Dateizuordnungen kann Ihre Anwendung bestimmte Dateitypen verarbeiten, wenn Benutzer sie öffnen. Dies ist besonders für Texteditoren, Bildbetrachter und andere Anwendungen nützlich, die mit bestimmten Dateiformaten arbeiten. In dieser Anleitung erfahren Sie, wie Sie Dateizuordnungen in Ihrer Wails-v3-Anwendung implementieren.

## Überblick

Dateizuordnungen werden in Wails v3 derzeit auf folgenden Plattformen unterstützt:

- Windows (NSIS-Installationspakete)
- macOS (Anwendungspakete)

## Konfiguration

Dateizuordnungen werden in der Datei `config.yml` im Verzeichnis `build` Ihres Projekts konfiguriert.

### Grundkonfiguration

So richten Sie Dateizuordnungen ein:

1. Öffnen Sie `build/config.yml`
2. Fügen Sie Ihre Dateizuordnungen im Abschnitt `fileAssociations` hinzu
3. Führen Sie `wails3 update build-assets` aus, um die Build-Assets zu aktualisieren
4. Legen Sie in den Anwendungsoptionen das Feld `FileAssociations` fest
5. Paketieren Sie Ihre Anwendung mit `wails3 package`

Hier sehen Sie eine Beispielkonfiguration:

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

### Konfigurationseigenschaften

| Eigenschaft | Beschreibung | Plattform |
| --- | --- | --- |
| ext | Dateierweiterung ohne führenden Punkt (z. B. `txt`) | Alle |
| name | Anzeigename des Dateityps | Alle |
| description | In den Dateieigenschaften angezeigte Beschreibung | Windows |
| iconName | Name der Symboldatei im Build-Ordner (ohne Dateierweiterung) | Alle |
| role | Rolle der Anwendung für diesen Dateityp (z. B. `Editor`, `Viewer`) | macOS |
| mimeType | MIME-Typ der Datei (z. B. `image/jpeg`) | macOS |

## Auf Dateiöffnungsereignisse reagieren

Um Dateiöffnungsereignisse in Ihrer Anwendung zu verarbeiten, können Sie auf das Ereignis `events.Common.ApplicationOpenedWithFile` reagieren:

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

## Schritt-für-Schritt-Anleitung

Richten Sie nun Schritt für Schritt Dateizuordnungen für einen einfachen Texteditor ein:

@steps
### Symbole erstellen
- Erstellen Sie Symbole für Ihren Dateityp (empfohlene Größen: 16x16, 32x32, 48x48, 256x256)
- Speichern Sie die Symbole im Ordner `build` Ihres Projekts
- Benennen Sie sie entsprechend Ihrer `iconName`-Konfiguration (z. B. `textFileIcon.png`)

@note{type="tip"}
Mit `wails3 generate icons` können Sie die benötigten Symbole generieren. Führen Sie `wails3 generate icons --help` aus, um weitere Informationen zu erhalten.

@end

- Fügen Sie für macOS in der Aufgabe `create:app:bundle:` eine Kopieranweisung wie `cp build/darwin/documenticon.icns {{.BIN_DIR}}/{{.APP_NAME}}.app/Contents/Resources` hinzu.

### Dateizuordnungen konfigurieren
Bearbeiten Sie die Datei `build/config.yml`, um Ihre Dateizuordnungen hinzuzufügen:

```yaml
# build/config.yml
fileAssociations:
  - ext: txt
    name: Text Document
    description: Plain Text Document
    iconName: textFileIcon
    role: Editor
```

### Build-Assets aktualisieren
Führen Sie den folgenden Befehl aus, um die Build-Assets zu aktualisieren:

```bash
wails3 update build-assets
```

### Dateizuordnungen in den Anwendungsoptionen festlegen
Legen Sie in Ihrer Datei `main.go` das Feld `FileAssociations` in den Anwendungsoptionen fest:

```go
app := application.New(application.Options{
  Name: "MyApp",
  FileAssociations: []string{".txt", ".md"}, // Specify supported extensions
})
```

@note{type="tip" title="Warum müssen Dateierweiterungen sowohl in der Anwendungskonfiguration als auch in config.yml angegeben werden?"}
Wenn unter Windows eine Datei über eine Dateizuordnung geöffnet wird, wird die Anwendung mit dem Dateinamen als erstem Argument gestartet. Die Anwendung kann nicht erkennen, ob das erste Argument eine Datei oder ein Befehlszeilenargument ist. Daher ermittelt sie anhand des Felds `FileAssociations` in den Anwendungsoptionen, ob das erste Argument eine zugeordnete Datei ist.

@end

### Anwendung paketieren
Paketieren Sie Ihre Anwendung mit dem folgenden Befehl:

```bash
wails3 package
```

Die paketierte Anwendung wird im Verzeichnis `bin` erstellt. Anschließend können Sie die Anwendung installieren und testen.

## Zusätzliche Hinweise

- Symbole sollten im Build-Ordner im PNG-Format bereitgestellt werden
- Zum Testen von Dateizuordnungen muss die paketierte Anwendung installiert werden

@end
