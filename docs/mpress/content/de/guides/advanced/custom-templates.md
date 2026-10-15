---
title: "Benutzerdefinierte Vorlagen erstellen"
description: "So erstellen Sie eigene Projektvorlagen für Wails v3, passen sie an und hosten sie"
slug: "guides/advanced/custom-templates"
sourcePath: "guides/advanced/custom-templates.md"
---

Wails wird mit mehreren integrierten Vorlagen ausgeliefert. Sie können jedoch eigene Vorlagen erstellen und mit der Community teilen. Eine benutzerdefinierte Vorlage ist lediglich ein Git-Repository. Sobald es öffentlich gehostet wird, kann jeder mit einem einzigen Befehl daraus ein Projektgerüst erstellen.

## Vorlagengerüst erstellen

Der Befehl `wails3 generate template` erstellt ein Verzeichnis mit einer direkt anpassbaren Vorlage:

```bash
wails3 generate template -name MyTemplate
```

Alle Optionen:

| Option | Beschreibung | Standardwert |
| --- | --- | --- |
| `-name` | Vorlagenname (erforderlich) | — |
| `-author` | Name des Autors | — |
| `-description` | In der CLI angezeigte Kurzbeschreibung | — |
| `-helpurl` | URL zur Dokumentation dieser Vorlage | — |
| `-version` | Ursprüngliche Version | `v0.0.1` |
| `-frontend` | Ein vorhandenes Frontend-Verzeichnis in die Vorlage kopieren | — |
| `-dir` | Speicherort für das Vorlagenverzeichnis | Aktuelles Verzeichnis |

Beispiel mit allen Optionen:

```bash
wails3 generate template \
  -name "My Template" \
  -author "Your Name" \
  -description "React + custom setup" \
  -helpurl "https://github.com/yourname/my-template" \
  -version "v1.0.0" \
  -frontend ./my-existing-frontend
```

Das erstellte Verzeichnis sieht wie folgt aus:

```
MyTemplate/
├── template.yaml          # Template metadata — edit this
├── NEXTSTEPS.md           # Guidance for you as the template author — delete before publishing
├── README.md              # Shown to users after they create a project
├── main.go.tmpl           # Application entry point
├── greetservice.go        # Example Go service
├── go.mod.tmpl            # Go module file
├── go.sum.tmpl            # Go checksums
├── gitignore.tmpl         # Becomes .gitignore in generated projects
├── Taskfile.tmpl.yml      # Build task definitions
└── frontend/              # Your frontend code
```

@note{type="tip" title="NEXTSTEPS.md lesen"}
Die erstellte Datei `NEXTSTEPS.md` enthält ausführliche Anleitungen zu allen Teilen der Vorlage. Lesen Sie sie vor der Anpassung. Löschen Sie sie vor der Veröffentlichung – sie darf nicht in Projekten enthalten sein, die aus Ihrer Vorlage erstellt werden.

@end

## Vorlagenmetadaten konfigurieren

Öffnen Sie `template.yaml`, um die Metadaten Ihrer Vorlage festzulegen:

```yaml
# yaml-language-server: $schema=https://v3.wails.io/schemas/template.v3.json
name: "My Template"
shortname: my-template
author: Your Name
description: A template with my preferred setup
helpurl: https://github.com/yourname/my-template
version: v1.0.0
wailsVersion: 3
```

Das Feld `wailsVersion` ist **erforderlich** und muss `3` sein. Der Kommentar `# yaml-language-server` am Anfang der Datei aktiviert die automatische Vervollständigung und Inline-Validierung in VS Code (mit der [YAML-Erweiterung](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml)) und in JetBrains-IDEs. Sie können ihn beibehalten oder entfernen; zur Laufzeit hat er keine Auswirkungen.

## Vorlage anpassen

### Frontend

Das Verzeichnis `frontend/` wird unverändert in jedes Projekt kopiert, das aus Ihrer Vorlage erstellt wird. Ersetzen Sie den Platzhalterinhalt durch Ihr tatsächliches Frontend:

@tabs
[Neu beginnen]
```bash
cd MyTemplate/frontend
npm create vite@latest .
```

Folgen Sie den Eingabeaufforderungen und installieren Sie anschließend die Abhängigkeiten:

```bash
npm install
```

[Vorhandenes Projekt verwenden]
Übergeben Sie beim Erstellen der Vorlage `-frontend`, um ein vorhandenes Frontend in einem Schritt zu kopieren:

```bash
wails3 generate template -name MyTemplate -frontend ./my-app/frontend
```

Alternativ können Sie es anschließend manuell in das Verzeichnis `frontend/` kopieren.

@end

### Build-Tasks

`Taskfile.tmpl.yml` definiert den Build-Ablauf. Passen Sie die Tasks `install:frontend:deps` und `build:frontend` an Ihre Frontend-Toolchain an:

```yaml
tasks:
  install:frontend:deps:
    dir: frontend
    cmds:
      - npm install       # replace with pnpm install, yarn, etc.

  build:frontend:
    dir: frontend
    deps: [install:frontend:deps, generate:bindings]
    cmds:
      - npm run build     # replace with your build command
```

### Go-Anwendung

Die Datei `main.go.tmpl` ist der Einstiegspunkt der Anwendung. Beim Erstellen eines Projekts wird sie von der Template-Engine von Wails verarbeitet. Vorlagenvariablen wie `{{.ProductName}}` werden dabei durch die vom Benutzer angegebenen Werte ersetzt.

Um sie als echte Go-Datei mit IDE-Unterstützung zu bearbeiten, benennen Sie sie vorübergehend in `main.go` um, nehmen Sie Ihre Änderungen vor und benennen Sie sie vor dem Commit wieder in `main.go.tmpl` um.

#### Vorlagenvariablen

Diese Variablen sind in jeder `.tmpl`-Datei verfügbar:

| Variable | Beschreibung | Beispiel |
| --- | --- | --- |
| `{{.ProjectName}}` | Vom Benutzer angegebener Projektname | `"MyApp"` |
| `{{.BinaryName}}` | Dateiname der Binärdatei | `"myapp"` |
| `{{.ProductName}}` | Anzeigename des Produkts | `"My Application"` |
| `{{.ProductDescription}}` | Produktbeschreibung | `"An awesome application"` |
| `{{.ProductVersion}}` | Produktversion | `"1.0.0"` |
| `{{.ProductCompany}}` | Unternehmens- oder Autorenname | `"My Company Ltd"` |
| `{{.ProductCopyright}}` | Copyright-Vermerk | `"Copyright 2024 My Company Ltd"` |
| `{{.ProductComments}}` | Zusätzliche Produktkommentare | `"Built with Wails"` |
| `{{.ProductIdentifier}}` | Produktkennung im Reverse-DNS-Format | `"com.mycompany.myapp"` |
| `{{.ModulePath}}` | Go-Modulpfad | `"github.com/you/myapp"` |
| `{{.WailsVersion}}` | Zum Erstellen des Projekts verwendete Wails-Version | `"3.0.0"` |
| `{{.Typescript}}` | `true`, wenn der Vorlagenname auf `-ts` endet | `true` |
| `{{.Opn}}` | Literales `{{` – innerhalb von Vorlagen maskieren | `{{` |
| `{{.Cls}}` | Literales `}}` – innerhalb von Vorlagen maskieren | `}}` |

@note{type="tip"}
Jede Datei in deiner Vorlage kann eine `.tmpl`-Datei sein – einschließlich HTML-, JSON- und YAML-Dateien. Dateien ohne das Suffix `.tmpl` werden unverändert kopiert.

@end

## Vorlage lokal testen

Teste die Vorlage vor der Veröffentlichung, indem du ein Projekt aus einem lokalen Pfad erstellst:

```bash
wails3 init -n testproject -t /path/to/MyTemplate
```

Prüfe anschließend, ob das Projekt funktioniert:

```bash
cd testproject
wails3 dev    # development mode with hot reload
wails3 build  # production binary
```

Prüfe Folgendes:

- Hot Reload des Frontends funktioniert
- Änderungen am Go-Code erstellen die App neu und starten sie erneut
- Die Produktionsbinärdatei in `bin/` wird ordnungsgemäß ausgeführt

## Auf GitHub veröffentlichen

@steps
### **Erstelle ein öffentliches GitHub-Repository** für deine Vorlage. Das Stammverzeichnis des Repositorys muss `template.yaml` enthalten.
### **Lösche `NEXTSTEPS.md`** – diese Datei enthält Hinweise für Vorlagenautoren und darf nicht in Projekten enthalten sein, die Benutzer aus deiner Vorlage erstellen.
### **Committe und pushe** den Inhalt des Vorlagenverzeichnisses als Stammverzeichnis des Repositorys:
```bash
git init
git add .
git commit -m "Initial template"
git remote add origin https://github.com/yourname/my-template.git
git push -u origin main
```

### **Erstelle ein Release-Tag** gemäß der semantischen Versionierung:
```bash
git tag v1.0.0
git push origin v1.0.0
```

@end

Benutzer können nun Projekte aus deiner Vorlage erstellen:

```bash
# Latest commit on the default branch
wails3 init -n myapp -t https://github.com/yourname/my-template

# Pinned to a specific release tag
wails3 init -n myapp -t https://github.com/yourname/my-template@v1.0.0
```

@note{type="caution" title="Warnung vor Drittanbieter-Vorlagen"}
Wenn ein Benutzer eine Remote-Vorlage installiert, zeigt Wails eine Warnung an, die darauf hinweist, dass die Vorlage Code eines Drittanbieters enthält und das Wails-Projekt keine Verantwortung für deren Inhalt übernimmt. Benutzer müssen ausdrücklich bestätigen, bevor das Projekt erstellt wird.

Als Vorlagenautor bist du für die Sicherheit und Korrektheit des gesamten Codes in deiner Vorlage verantwortlich.

@end

## Bewährte Verfahren

- **Verfasse eine verständliche `README.md`** – sie wird Benutzern angezeigt, nachdem sie ein Projekt erstellt haben. Erkläre, wie das Projekt ausgeführt, gebaut und angepasst wird.
- **Fülle `helpurl`** aus – verlinke dein Repository oder eine eigene Dokumentation. Benutzer sehen den Link in der Vorlagenliste der Wails CLI.
- **Fixiere die Versionen der Frontend-Abhängigkeiten** in `package.json`, damit Aktualisierungen vorgelagerter Projekte Installationen nicht beeinträchtigen.
- **Teste vor dem Tagging** – erstelle ein neues Projekt aus dem getaggten Release, bevor du es der Community ankündigst.
- **Behalte `wailsVersion: 3`** bei – dieses Feld teilt Wails mit, auf welche Hauptversion die Vorlage ausgerichtet ist. Ändere es nicht.
- **Aktualisiere regelmäßig** – halte die Abhängigkeiten aktuell und teste mit neuen Wails-Releases.
