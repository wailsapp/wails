---
title: "Benutzerdefinierte Vorlagen erstellen"
description: "Eigene Wails-v3-Projektvorlagen erstellen, anpassen und hosten"
slug: "guides/advanced/custom-templates"
sourcePath: "guides/advanced/custom-templates.md"
---

Wails enthält mehrere integrierte Vorlagen. Du kannst jedoch auch eigene Vorlagen erstellen und mit der Community teilen. Eine benutzerdefinierte Vorlage ist lediglich ein Git-Repository. Sobald es öffentlich gehostet wird, kann jeder daraus mit einem einzigen Befehl ein Projekt erzeugen.

## Vorlagengerüst erzeugen

Der Befehl `wails3 generate template` erzeugt ein anpassungsbereites Vorlagenverzeichnis:

```bash
wails3 generate template -name MyTemplate
```

Alle Optionen:

| Option | Beschreibung | Standardwert |
| --- | --- | --- |
| `-name` | Vorlagenname (erforderlich) | — |
| `-author` | Name des Autors | — |
| `-description` | Kurze, in der CLI angezeigte Beschreibung | — |
| `-helpurl` | URL zur Dokumentation dieser Vorlage | — |
| `-version` | Ausgangsversion | `v0.0.1` |
| `-frontend` | Vorhandenes Frontend-Verzeichnis in die Vorlage kopieren | — |
| `-dir` | Speicherort des Vorlagenverzeichnisses | Aktuelles Verzeichnis |

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

Das erzeugte Verzeichnis sieht wie folgt aus:

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
Die erzeugte Datei `NEXTSTEPS.md` enthält ausführliche Anleitungen zu allen Bestandteilen der Vorlage. Lies sie vor der Anpassung. Lösche sie vor der Veröffentlichung – sie darf nicht in Projekten enthalten sein, die aus deiner Vorlage erstellt werden.

@end

## Vorlagenmetadaten konfigurieren

Öffne `template.yaml`, um die Metadaten deiner Vorlage festzulegen:

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

Das Feld `wailsVersion` ist **erforderlich** und muss den Wert `3` haben. Der Kommentar `# yaml-language-server` am Anfang aktiviert die automatische Vervollständigung und Inline-Validierung in VS Code (mit der [YAML-Erweiterung](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml)) sowie in JetBrains-IDEs. Du kannst ihn beibehalten oder entfernen; zur Laufzeit hat er keine Wirkung.

## Vorlage anpassen

### Frontend

Das Verzeichnis `frontend/` wird unverändert in jedes Projekt kopiert, das aus deiner Vorlage erstellt wird. Ersetze den Platzhalterinhalt durch dein tatsächliches Frontend:

@tabs
[Von Grund auf beginnen]
```bash
cd MyTemplate/frontend
npm create vite@latest .
```

Folge den Eingabeaufforderungen und installiere anschließend die Abhängigkeiten:

```bash
npm install
```

[Vorhandenes Projekt verwenden]
Übergib beim Erzeugen der Vorlage `-frontend`, um ein vorhandenes Frontend in einem Schritt zu kopieren:

```bash
wails3 generate template -name MyTemplate -frontend ./my-app/frontend
```

Alternativ kannst du es anschließend manuell in das Verzeichnis `frontend/` kopieren.

@end

### Build-Tasks

`Taskfile.tmpl.yml` definiert den Build-Ablauf. Passe die Tasks `install:frontend:deps` und `build:frontend` an deine Frontend-Toolchain an:

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

Um sie als reguläre Go-Datei mit IDE-Unterstützung zu bearbeiten, benenne sie vorübergehend in `main.go` um, nimm deine Änderungen vor und benenne sie vor dem Commit wieder in `main.go.tmpl` um.

#### Vorlagenvariablen

In jeder `.tmpl`-Datei sind folgende Variablen verfügbar:

| Variable | Beschreibung | Beispiel |
| --- | --- | --- |
| `{{.ProjectName}}` | Vom Benutzer angegebener Projektname | `"MyApp"` |
| `{{.BinaryName}}` | Dateiname der Binärdatei | `"myapp"` |
| `{{.ProductName}}` | Anzeigename des Produkts | `"My Application"` |
| `{{.ProductDescription}}` | Produktbeschreibung | `"An awesome application"` |
| `{{.ProductVersion}}` | Produktversion | `"1.0.0"` |
| `{{.ProductCompany}}` | Name des Unternehmens/Autors | `"My Company Ltd"` |
| `{{.ProductCopyright}}` | Copyright-Zeichenfolge | `"Copyright 2024 My Company Ltd"` |
| `{{.ProductComments}}` | Zusätzliche Produktkommentare | `"Built with Wails"` |
| `{{.ProductIdentifier}}` | Produktkennung im Reverse-DNS-Format | `"com.mycompany.myapp"` |
| `{{.ModulePath}}` | Go-Modulpfad | `"github.com/you/myapp"` |
| `{{.WailsVersion}}` | Zum Erstellen des Projekts verwendete Wails-Version | `"3.0.0"` |
| `{{.Typescript}}` | `true`, wenn der Vorlagenname auf `-ts` endet | `true` |
| `{{.Opn}}` | Literales `{{` – innerhalb von Vorlagen maskieren | `{{` |
| `{{.Cls}}` | Literales `}}` – innerhalb von Vorlagen maskieren | `}}` |

@note{type="tip"}
Jede Datei in deiner Vorlage kann eine `.tmpl`-Datei sein, einschließlich HTML-, JSON- und YAML-Dateien. Dateien ohne das Suffix `.tmpl` werden unverändert kopiert.

@end

## Vorlage lokal testen

Teste die Vorlage vor der Veröffentlichung, indem du ein Projekt aus einem lokalen Pfad erstellst:

```bash
wails3 init -n testproject -t /path/to/MyTemplate
```

Überprüfe anschließend, ob das Projekt funktioniert:

```bash
cd testproject
wails3 dev    # development mode with hot reload
wails3 build  # production binary
```

Prüfe Folgendes:

- Hot Reload des Frontends funktioniert
- Änderungen am Go-Code erstellen die Anwendung neu und starten sie erneut
- Die Produktionsbinärdatei in `bin/` wird ordnungsgemäß ausgeführt

## Auf GitHub veröffentlichen

@steps
### **Erstellen Sie ein öffentliches GitHub-Repository** für Ihre Vorlage. Das Stammverzeichnis des Repositorys muss `template.yaml` enthalten.
### **Löschen Sie `NEXTSTEPS.md`** – diese Datei dient Vorlagenautoren als Anleitung und darf nicht in Projekten enthalten sein, die Benutzer aus Ihrer Vorlage erstellen.
### **Committen und pushen Sie** den Inhalt des Vorlagenverzeichnisses als Stammverzeichnis des Repositorys:
```bash
git init
git add .
git commit -m "Initial template"
git remote add origin https://github.com/yourname/my-template.git
git push -u origin main
```

### **Kennzeichnen Sie ein Release mit einem Tag** und verwenden Sie dabei semantische Versionierung:
```bash
git tag v1.0.0
git push origin v1.0.0
```

@end

Benutzer können nun Projekte aus Ihrer Vorlage erstellen:

```bash
# Latest commit on the default branch
wails3 init -n myapp -t https://github.com/yourname/my-template

# Pinned to a specific release tag
wails3 init -n myapp -t https://github.com/yourname/my-template@v1.0.0
```

@note{type="caution" title="Warnung zu Drittanbieter-Vorlagen"}
Wenn ein Benutzer eine Remote-Vorlage installiert, zeigt Wails eine Warnung an. Sie weist darauf hin, dass die Vorlage Code eines Drittanbieters enthält und das Wails-Projekt keine Verantwortung für deren Inhalt übernimmt. Benutzer müssen ausdrücklich zustimmen, bevor das Projekt erstellt wird.

Als Vorlagenautor sind Sie für die Sicherheit und Korrektheit des gesamten Codes in Ihrer Vorlage verantwortlich.

Vorlagen sind Code. `Taskfile.yml`-Tasks sind beliebige Shell-Befehle. `wails3 dev`, `wails3 build` und `wails3 package` führen sie sowie alle von ihnen aufgerufenen `package.json`-Skripte mit den Berechtigungen des Benutzers aus. Prüfen Sie `Taskfile.yml`, das Verzeichnis `build/` und `frontend/package.json` einer Remote-Vorlage, bevor Sie einen dieser Befehle ausführen.

@end

## Bewährte Vorgehensweisen

- **Verfassen Sie eine verständliche `README.md`** – sie wird Benutzern angezeigt, nachdem sie ein Projekt erstellt haben. Erklären Sie, wie das Projekt ausgeführt, gebaut und angepasst wird.
- **Füllen Sie `helpurl` aus** – verlinken Sie auf Ihr Repository oder eine eigene Dokumentation. Benutzer sehen den Link in der Vorlagenliste der Wails CLI.
- **Fixieren Sie die Versionen der Frontend-Abhängigkeiten** in `package.json`, um fehlerhafte Installationen aufgrund von Upstream-Aktualisierungen zu vermeiden.
- **Testen Sie vor dem Tagging** – erstellen Sie ein neues Projekt aus dem getaggten Release, bevor Sie es der Community ankündigen.
- **Behalten Sie `wailsVersion: 3` bei** – dieses Feld teilt Wails mit, für welche Hauptversion die Vorlage vorgesehen ist. Ändern Sie es nicht.
- **Aktualisieren Sie regelmäßig** – halten Sie die Abhängigkeiten aktuell und testen Sie mit neuen Wails-Releases.
