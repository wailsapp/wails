---
title: "Vorlagensystem"
description: "Wie Wails v3 neue Projekte erzeugt, wie Vorlagen organisiert sind und wie Sie eigene Vorlagen erstellen."
slug: "contributing/template-system"
sourcePath: "contributing/template-system.md"
---

Wails enthält ein **Vorlagensystem**, mit dem `wails3 init` ein sofort ausführbares Projekt erzeugen kann. Für eine bewusst klein gehaltene Auswahl an Frameworks sind Vorlagen integriert (Vanilla, React, Vue, Svelte). Jedes andere Framework lässt sich verwenden, indem Sie [Ihr eigenes Frontend einbinden](/guides/dev/frontend-frameworks/) oder eine [benutzerdefinierte Vorlage](/guides/advanced/custom-templates/) veröffentlichen.

Diese Seite behandelt:

1. Verzeichnisstruktur der Vorlagen
2. Auswahl und Verarbeitung von Vorlagen durch die CLI
3. Schrittweises Erstellen einer neuen Vorlage
4. Aktualisieren oder Überschreiben vorhandener Vorlagen
5. Fehlerbehebung und bewährte Verfahren

---

## 1. Speicherort der Vorlagen

```
v3/internal/templates/
├── _common/        # Files copied into EVERY project (Taskfile.yml, build/, etc.)
├── base/           # Backend-only "plain Go" base layer (frontend/ + NEXTSTEPS.md)
├── ios/            # iOS bootstrapper
├── vanilla/        vanilla-js/   # TypeScript (default) + JavaScript variant
├── react/          react-js/     # TypeScript (default) + JavaScript variant
├── vue/                          # TypeScript only
├── svelte/                       # TypeScript only
└── templates.go    # Registry + Install/Get APIs (no auto-registration via embed)
```

- **`_common/`** — allgemeiner Vorlagencode (Taskfile, Verzeichnis `build/` und gemeinsam genutzte Infrastruktur), der in jedes Projekt eingebunden wird.
- **`base/`** — der Go-Teil, auf dem jede Vorlage basiert. Hinweis: `base/` selbst enthält **keine** `template.json`; diese Datei befindet sich in der jeweiligen Framework-spezifischen Vorlage.
- **Framework-Ordner** — enthalten das Frontend (`frontend/`), die Framework-Konfiguration und eine `template.json` mit den Metadaten der Vorlage.
- Die Ordnernamen entsprechen der **Vorlagen-ID**, die Sie an die CLI übergeben (`wails3 init -t react`).
- **Sprachkonvention:** TypeScript ist die Standardsprache und verwendet den Namen ohne Suffix (`react`). Eine JavaScript-Variante erhält, sofern vorhanden, das Suffix `-js` (`react-js`). Integrierte Vorlagen deklarieren ihre Sprache in `template.yaml` explizit mit `typescript: true|false`. Community-Vorlagen können weiterhin das veraltete Suffix `-ts` verwenden, das als Ausweichlösung berücksichtigt wird.

> Das gesamte Verzeichnis `internal/templates/` wird in die CLI-Binärdatei kompiliert
>
> und zwar über `//go:embed *`, sodass Benutzer Projekte offline erzeugen können.

---

## 2. Verwendung von Vorlagen durch `wails3 init`

Aufrufkette (ohne `cmd/wails3/init.go` — die CLI wird direkt in `cmd/wails3/main.go` verdrahtet):

```
cmd/wails3/main.go             (clir wiring)
       │
       ▼
internal/commands/init.go      Init(options *flags.Init) error
       │
       ▼
internal/templates/templates.go
       │   templates.Install(options)
       │   templates.GetDefaultTemplates()
       ▼
gosod.New(template.FS).Extract(options.ProjectDir, data)   // file extraction
       │
       ▼
go mod tidy (unless --skipgomodtidy / -skipgomodtidy)
```

Es gibt keine API für `Template.Load()` / `Template.CopyTo()` / `Template.Validate()` — `gosod` (`github.com/leaanthony/gosod`) übernimmt die Extraktion aus der eingebetteten `fs.FS`.

### Optionen für `wails3 init`

Definiert in `internal/flags/init.go`:

| Option | Zweck | Standardwert |
| --- | --- | --- |
| `-p` | Paketname | `main` |
| `-t` | Name einer integrierten Vorlage, lokaler Pfad oder URL | `vanilla` |
| `-n` | Projektname | (leer) |
| `-d` | Projektverzeichnis | `.` |
| `-q` | Konsolenausgabe unterdrücken | false |
| `-l` | Vorlagen auflisten | false |
| `-skipgomodtidy` | Ausführung von `go mod tidy` nach der Extraktion überspringen | false |
| `-git` | URL des zu initialisierenden Git-Repositorys | (leer) |
| `-mod` | Go-Modulpfad (wird aus `-git` abgeleitet, falls nicht festgelegt) | (leer) |
| `-s` | Warnung bei Verwendung entfernter Vorlagen überspringen | false |
| `-productname` / `-productdescription` / `-productversion` / `-productcompany` / `-productcopyright` / `-productcomments` / `-productidentifier` | In die generierten Build-Assets eingebettete Metadaten | sinnvolle Standardwerte |

Es gibt **keinen** langen Alias `-list` (nur `-l`) und keine vorlagenspezifische `--help`.

### Ersetzungen

Platzhalter sind standardmäßige Go-Vorlagenanweisungen — das vorangestellte `.` ist Teil des Feldzugriffs:

| Platzhalter | Beispiel | Quelle |
| --- | --- | --- |
| `{{.ProjectName}}` | `myapp` | Flag `-n` / Verzeichnisname |
| `{{.ModulePath}}` | `github.com/me/myapp` | Flag `-mod` oder aus `-git` abgeleitet |
| `{{.WailsVersion}}` | `v3.0.0-…` | Einkompilierte Konstante aus `internal/version` |
| `{{.ProductName}}`, `{{.ProductDescription}}`, `{{.ProductVersion}}`, `{{.ProductCompany}}`, `{{.ProductCopyright}}`, `{{.ProductComments}}`, `{{.ProductIdentifier}}` | Metadaten zum Build-Zeitpunkt | entsprechende `-product*`-Flags |

Wenn Sie einen neuen Platzhalter benötigen, fügen Sie den Vorlagendaten in `internal/templates/templates.go` ein Feld und in `internal/flags/init.go` ein passendes Feld/Flag hinzu (oder setzen Sie es aus `internal/commands/init.go`).

### Hook nach dem Kopieren

Nachdem `gosod` das Entpacken der Vorlage abgeschlossen hat, führt die CLI Folgendes aus:

```
go mod tidy
```

sofern Sie nicht `-skipgomodtidy` übergeben. Einen Schritt `task deps` gibt es nicht.

---

## 3. Eine neue Vorlage erstellen

> Beispiel: Eine **Solid**-Vorlage hinzufügen

### 3.1 Ordner und ID

```
internal/templates/solid/
```

Der Ordnername entspricht der Vorlagen-ID. Verwenden Sie dafür **kebab-case**.

### 3.2 Minimaler Dateisatz

```
solid/
├── template.yaml    # name, description, wailsVersion, typescript (required)
├── frontend/        # Your web project (no node_modules/dist)
│   ├── src/
│   ├── package.json
│   └── vite.config.ts
└── ...              # Any extra Go files the template wants to inject
```

Kopieren Sie zunächst `react` und entfernen Sie nicht benötigte Dateien. Vergessen Sie nicht, eine `template.yaml` zu erstellen – setzen Sie für eine TypeScript-Vorlage `typescript: true`. `base/` ist der einzige Ordner ohne eine solche Datei.

### 3.3 Platzhalter aktualisieren

Suchen und ersetzen Sie wörtliche Beispielwerte durch Go-Vorlagendirektiven, zum Beispiel:

- `myapp` → `{{.ProjectName}}`
- `github.com/you/myapp` → `{{.ModulePath}}`

### 3.4 Einbinden

Da `templates.go` das eingebettete Dateisystem bei der Initialisierung durchläuft, genügt es in der Regel, unter `internal/templates/<id>/` einen neuen Ordner anzulegen – ein manueller Registrierungsaufruf ist nicht erforderlich. Wenn Sie zusätzliche Logik benötigen, etwa eine benutzerdefinierte Validierung oder Schritte nach dem Kopieren, fügen Sie diese zu `templates.Install` in `internal/templates/templates.go` hinzu.

### 3.5 Testen

```bash
wails3 init -n demo -t solid
cd demo
wails3 dev
```

Stellen Sie Folgendes sicher:

- Der Entwicklungsserver startet auf dem in `WAILS_VITE_PORT` veröffentlichten Port
- Die generierten Bindings werden unter `frontend/bindings/...` angelegt
- Hot Reload funktioniert

---

## 4. Vorhandene Vorlagen ändern

1. Bearbeiten Sie die Dateien unter `internal/templates/<id>/`.
2. Erstellen Sie die CLI neu (`cd v3 && go build -o ../wails3 ./cmd/wails3`); die Direktive `//go:embed *` übernimmt die neuen Inhalte.
3. Aktualisieren Sie alle **Abhängigkeitsversionen** in `frontend/package.json` und `Taskfile.yml`.
4. Aktualisieren Sie die Beschreibung der Vorlage in `template.json`, wenn sich ihr Verhalten ändert.

### Häufige Anpassungen

| Aufgabe | Ort |
| --- | --- |
| Port des Entwicklungsservers ändern | `frontend/vite.config.ts` – `WAILS_VITE_PORT` auslesen |
| Umgebungsvariablen hinzufügen | `build/Taskfile.yml` oder `frontend/.env` |
| JS-Paketmanager ersetzen | In `build/Taskfile.yml` `npm` durch `pnpm`/`bun` ersetzen |

---

## 5. Tipps zum Erstellen von Vorlagen

- **Frontend generisch halten** – vermeiden Sie Verweise auf Wails-spezifische globale Variablen; `/wails/runtime.js` wird zur Laufzeit vom Asset-Server bereitgestellt.
- **Keine kompilierten Artefakte** – schließen Sie `node_modules`, `dist` und `.DS_Store` aus dem eingebetteten Verzeichnis aus (oder tragen Sie sie in `.gitignore` ein, damit sie nie committet werden).
- **Voraussetzungen dokumentieren** – dokumentieren Sie die Node-Version, zusätzliche CLI-Werkzeuge usw. in `template.json` oder `NEXTSTEPS.md`.
- **Breaking Changes vermeiden** – wenn die Überarbeitung umfangreich ist, erstellen Sie eine neue Vorlagen-ID, statt eine vorhandene Vorlage zu verändern.

---

## 6. Fehlerbehebung

| Symptom | Ursache | Behebung |
| --- | --- | --- |
| `unknown template name` | Tippfehler in `-t` oder Vorlage nicht eingebettet | Führe `wails3 init -l` aus, um die verfügbaren Vorlagen aufzulisten |
| Platzhalter werden nicht ersetzt | `{{ProjectName}}` statt `{{.ProjectName}}` verwendet | Füge das führende `.` hinzu (Feldzugriff in Go-Vorlagen) |
| Entwicklungsserver öffnet eine leere Seite | Vite-Konfiguration liest `WAILS_VITE_PORT` nicht ein | Überprüfe deine `vite.config.ts` |
| Frontend-Build für die Produktionsumgebung schlägt fehl | Vite-Pfad `base` vergessen | Setze `base: "./"` in `vite.config.ts` |

---

## 7. Übersicht der wichtigsten Quelldateien

| Datei | Aufgabe |
| --- | --- |
| `internal/templates/templates.go` | Bettet das Vorlagen-Dateisystem ein und stellt `Install(options *flags.Init) error`, `GetDefaultTemplates()` und `ValidTemplateName(name)` bereit |
| `internal/templates/<id>/**` | Eigentlicher Vorlageninhalt |
| `internal/commands/init.go` | CLI-Anbindung: Wählt die Vorlage aus, füllt die Metadaten aus und ruft `templates.Install` auf |
| `internal/commands/generate_template.go` | `wails3 generate template` — Dienstprogramm, um ein aktives Projekt wieder als Vorlage zu *exportieren* (praktisch für Aktualisierungen) |
| `internal/flags/init.go` | Flag-Definitionen für `wails3 init` |

---

## 8. Zusammenfassung

- Vorlagen befinden sich in **`internal/templates/`** und werden über `//go:embed *` in die CLI eingebettet.
- `wails3 init -t <id>` extrahiert die Vorlage mithilfe von `gosod` und führt `go mod tidy` aus (mit `-skipgomodtidy` überspringbar).
- Zum Erstellen einer Vorlage musst du lediglich **einen Ordner anlegen**, Dateien sowie eine `template.json` hinzufügen und Platzhalter im Stil von `{{.ProjectName}}` verwenden.
- Das System ist **erweiterbar** und **eigenständig** — ideal, um benutzerdefinierte Stacks mit deinem Team oder der Community zu teilen.

Viel Spaß beim Erstellen von Vorlagen!
