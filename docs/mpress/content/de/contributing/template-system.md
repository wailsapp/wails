---
title: "Vorlagensystem"
description: "Wie Wails v3 neue Projekte einrichtet, wie Vorlagen organisiert sind und wie Sie eigene Vorlagen erstellen."
slug: "contributing/template-system"
sourcePath: "contributing/template-system.md"
---

Wails enthält ein **Vorlagensystem**, mit dem `wails3 init` ein sofort ausführbares Projekt erzeugen kann. Für eine bewusst kleine Auswahl an Frameworks sind Vorlagen integriert (Vanilla, React, Vue, Svelte). Jedes andere Framework lässt sich verwenden, indem Sie [Ihr eigenes Frontend mitbringen](/guides/dev/frontend-frameworks/) oder eine [benutzerdefinierte Vorlage](/guides/advanced/custom-templates/) veröffentlichen.

Diese Seite behandelt:

1. Verzeichnisstruktur der Vorlagen
2. Wie die CLI Vorlagen auswählt und rendert
3. Schrittweises Erstellen einer neuen Vorlage
4. Aktualisieren oder Überschreiben vorhandener Vorlagen
5. Fehlerbehebung und bewährte Vorgehensweisen

---

## 1. Speicherort der Vorlagen

```
v3/internal/templates/
├── _common/        # Files copied into EVERY project (Taskfile.yml, build/, etc.)
├── base/           # Minimal frontend skeleton (frontend/ + NEXTSTEPS.md)
├── ios/            # iOS bootstrapper
├── vanilla/        vanilla-js/   # TypeScript (default) + JavaScript variant
├── react/          react-js/     # TypeScript (default) + JavaScript variant
├── vue/                          # TypeScript only
├── svelte/                       # TypeScript only
└── templates.go    # Registry + Install/Get APIs (no auto-registration via embed)
```

- **`_common/`** – universelles Grundgerüst (Taskfile, Verzeichnis `build/` und gemeinsam genutzte Infrastruktur), das in jedes Projekt eingebunden wird.
- **`base/`** – ein minimales Frontend-Grundgerüst für benutzerdefinierte Vorlagen, keine automatische Backend-Schicht. Beim regulären Erstellen der Arbeitskopie wird `_common/` mit der ausgewählten Framework-Vorlage überlagert.
- **Framework-Ordner** – enthalten das Frontend (`frontend/`), die Framework-Konfiguration und eine Datei `template.yaml`, die die Metadaten der Vorlage beschreibt.
- Die Ordnernamen entsprechen der **Vorlagen-ID**, die Sie an die CLI übergeben (`wails3 init -t react`).
- **Sprachkonvention:** TypeScript ist die Standardeinstellung und verwendet den Namen ohne Suffix (`react`). Eine JavaScript-Variante erhält, sofern vorhanden, das Suffix `-js` (`react-js`). Integrierte Vorlagen deklarieren ihre Sprache in `template.yaml` explizit mit `typescript: true|false`. Community-Vorlagen dürfen weiterhin das bisherige Suffix `-ts` verwenden, das als Fallback berücksichtigt wird.

> Das gesamte Verzeichnis `internal/templates/` wird in die CLI-Binärdatei kompiliert
>
> über `//go:embed *`, sodass Benutzer Projekte offline einrichten können.

---

Der Loader akzeptiert außerdem ältere `template.json`-Metadaten. Neue integrierte Vorlagen verwenden YAML einschließlich der expliziten Einstellung `typescript`.

## 2. Wie `wails3 init` Vorlagen verwendet

Aufrufkette (kein `cmd/wails3/init.go` – die CLI wird direkt in `cmd/wails3/main.go` verdrahtet):

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

Es gibt keine API `Template.Load()` / `Template.CopyTo()` / `Template.Validate()`. Das Erstellen der Arbeitskopie übernimmt `gosod` (`github.com/leaanthony/gosod`) anhand von `fs.FS`, das in die Binärdatei eingebettet ist.

### Flags für `wails3 init`

Definiert in `internal/flags/init.go`:

| Flag | Zweck | Standardwert |
| --- | --- | --- |
| `-p` | Paketname | `main` |
| `-t` | Name einer integrierten Vorlage, lokaler Pfad oder URL | `vanilla` |
| `-n` | Projektname | (leer) |
| `-d` | Projektverzeichnis | `.` |
| `-q` | Konsolenausgabe unterdrücken | false |
| `-l` | Vorlagen auflisten | false |
| `-skipgomodtidy` | Ausführen von `go mod tidy` nach dem Erstellen der Arbeitskopie überspringen | false |
| `-git` | URL des zu initialisierenden Git-Repositorys | (leer) |
| `-mod` | Go-Modulpfad (wird aus `-git` abgeleitet, wenn nicht festgelegt) | (leer) |
| `-s` | Warnung bei der Verwendung entfernter Vorlagen überspringen | false |
| `-productname` / `-productdescription` / `-productversion` / `-productcompany` / `-productcopyright` / `-productcomments` / `-productidentifier` | In die generierten Build-Artefakte eingebettete Metadaten | sinnvolle Standardwerte |

Es gibt keinen langen Alias `-list` (nur `-l`) und kein templatespezifisches `--help`.

### Ersetzungen

Platzhalter sind reguläre Go-Vorlagendirektiven – der vorangestellte Teil `.` gehört zum Feldzugriff:

| Platzhalter | Beispiel | Quelle |
| --- | --- | --- |
| `{{.ProjectName}}` | `myapp` | Flag `-n` / Verzeichnisname |
| `{{.ModulePath}}` | `github.com/me/myapp` | Flag `-mod` oder aus `-git` abgeleitet |
| `{{.WailsVersion}}` | `v3.0.0` | Einkompilierte Konstante aus `internal/version` |
| `{{.ProductName}}`, `{{.ProductDescription}}`, `{{.ProductVersion}}`, `{{.ProductCompany}}`, `{{.ProductCopyright}}`, `{{.ProductComments}}`, `{{.ProductIdentifier}}` | Metadaten zum Build-Zeitpunkt | entsprechende `-product*`-Flags |

Wenn Sie einen neuen Platzhalter benötigen, fügen Sie den Vorlagendaten in `internal/templates/templates.go` ein Feld und in `internal/flags/init.go` ein entsprechendes Feld beziehungsweise Flag hinzu (oder legen Sie es über `internal/commands/init.go` fest).

### Hook nach dem Kopieren

Nachdem `gosod` das Erstellen der Arbeitskopie aus der Vorlage abgeschlossen hat, führt die CLI Folgendes aus:

```
go mod tidy
```

sofern Sie nicht `-skipgomodtidy` übergeben. Es gibt keinen Schritt `task deps`.

---

## 3. Eine neue Vorlage erstellen

> Beispiel: Eine **Solid**-Vorlage hinzufügen

### 3.1 Ordner und ID

```
internal/templates/solid/
```

Der Ordnername entspricht der Vorlagen-ID. Verwenden Sie **Kebab-Case**.

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

Kopiere zunächst `react` und entferne nicht benötigte Dateien. Vergiss nicht, eine `template.yaml` zu erstellen — setze für ein TypeScript-Template `typescript: true`; `base/` ist der einzige Ordner ohne eine solche Datei.

### 3.3 Platzhalter aktualisieren

Suche nach den literalen Beispielwerten und ersetze sie durch Go-Template-Direktiven, zum Beispiel:

- `myapp` → `{{.ProjectName}}`
- `github.com/you/myapp` → `{{.ModulePath}}`

### 3.4 Einbinden

Da `templates.go` das eingebettete Dateisystem bei der Initialisierung durchläuft, genügt es normalerweise, unter `internal/templates/<id>/` einen neuen Ordner abzulegen — ein manueller Registrierungsaufruf ist nicht erforderlich. Wenn du zusätzliche Logik benötigst, etwa eine benutzerdefinierte Validierung oder Schritte nach dem Kopieren, füge sie in `templates.Install` in `internal/templates/templates.go` ein.

### 3.5 Testen

```bash
wails3 init -n demo -t solid
cd demo
wails3 dev
```

Stelle Folgendes sicher:

- Der Entwicklungsserver startet auf dem in `WAILS_VITE_PORT` veröffentlichten Port
- Die generierten Bindings erscheinen unter `frontend/bindings/...`
- Hot Reload funktioniert

---

## 4. Vorhandene Templates ändern

1. Bearbeite die Dateien unter `internal/templates/<id>/`.
2. Erstelle die CLI neu (`cd v3 && go build -o ../wails3 ./cmd/wails3`); die Direktive `//go:embed *` übernimmt die neuen Inhalte.
3. Aktualisiere alle **Abhängigkeitsversionen** in `frontend/package.json` und `Taskfile.yml`.
4. Aktualisiere die Beschreibung in `template.yaml` des Templates, wenn sich dessen Verhalten ändert.

### Häufige Anpassungen

| Aufgabe | Ort |
| --- | --- |
| Port des Entwicklungsservers ändern | `frontend/vite.config.ts` — `WAILS_VITE_PORT` auslesen |
| Umgebungsvariablen hinzufügen | `build/Taskfile.yml` oder `frontend/.env` |
| JS-Paketmanager ersetzen | In `build/Taskfile.yml` `npm` durch `pnpm`/`bun` ersetzen |

---

## 5. Tipps zum Erstellen von Templates

- **Frontend generisch halten** — vermeide Verweise auf Wails-spezifische globale Variablen; `/wails/runtime.js` wird zur Laufzeit vom Asset-Server bereitgestellt.
- **Keine kompilierten Artefakte** — schließe `node_modules`, `dist` und `.DS_Store` aus dem eingebetteten Verzeichnis aus (oder trage sie in `.gitignore` ein, damit sie nie committet werden).
- **Voraussetzungen dokumentieren** — dokumentiere die Node-Version, zusätzliche CLI-Werkzeuge usw. in `template.yaml` oder `NEXTSTEPS.md`.
- **Inkompatible Änderungen vermeiden** — wenn die Überarbeitung umfangreich ist, erstelle eine neue Template-ID, statt eine vorhandene zu verändern.

---

## 6. Fehlerbehebung

| Symptom | Ursache | Lösung |
| --- | --- | --- |
| `unknown template name` | Tippfehler in `-t` oder Template nicht eingebettet | Führe `wails3 init -l` aus, um die verfügbaren Templates aufzulisten |
| Platzhalter werden nicht ersetzt | `{{ProjectName}}` statt `{{.ProjectName}}` verwendet | Füge das vorangestellte `.` hinzu (Feldzugriff in Go-Templates) |
| Der Entwicklungsserver öffnet eine leere Seite | Die Vite-Konfiguration liest `WAILS_VITE_PORT` nicht aus | Überprüfe deine `vite.config.ts` |
| Das Frontend kann für die Produktion nicht gebaut werden | Vite-Pfad `base` vergessen | Setze `base: "./"` in `vite.config.ts` |

---

## 7. Übersicht der wichtigsten Quelldateien

| Datei | Zuständigkeit |
| --- | --- |
| `internal/templates/templates.go` | Bettet das Template-Dateisystem ein und stellt `Install(options *flags.Init) error`, `GetDefaultTemplates()` und `ValidTemplateName(name)` bereit |
| `internal/templates/<id>/**` | Eigentliche Template-Inhalte |
| `internal/commands/init.go` | CLI-Verknüpfung: wählt das Template aus, trägt Metadaten ein und ruft `templates.Install` auf |
| `internal/commands/generate_template.go` | `wails3 generate template` — Hilfsprogramm, um ein aktives Projekt wieder als Template zu *exportieren* (praktisch für Aktualisierungen) |
| `internal/flags/init.go` | Flag-Definitionen für `wails3 init` |

---

## 8. Zusammenfassung

- Templates befinden sich in **`internal/templates/`** und werden über `//go:embed *` in die CLI eingebettet.
- `wails3 init -t <id>` extrahiert das Template über `gosod` und führt `go mod tidy` aus; mit `-skipgomodtidy` lässt sich dieser Schritt überspringen.
- Um ein Template zu erstellen, musst du lediglich **einen Ordner anlegen**, Dateien sowie eine `template.yaml` hinzufügen und Platzhalter im Stil von `{{.ProjectName}}` verwenden.
- Das System ist **erweiterbar** und **eigenständig** — ideal, um benutzerdefinierte Stacks mit deinem Team oder der Community zu teilen.

Viel Freude beim Erstellen von Templates!
