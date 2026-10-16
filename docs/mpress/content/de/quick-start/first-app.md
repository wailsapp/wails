---
title: "Ihre erste App"
description: "Erstellen Sie in 10 Minuten eine funktionsfähige Wails-Anwendung"
slug: "quick-start/first-app"
sourcePath: "quick-start/first-app.md"
---

Wir erstellen eine einfache Begrüßungsanwendung, die die grundlegenden Konzepte von Wails veranschaulicht:

- Go-Backend zur Verwaltung der Logik
- Frontend, das Go-Funktionen aufruft
- Typsichere Bindings
- Hot Reload während der Entwicklung

**Benötigte Zeit:** 10 Minuten

@note{type="tip" title="Performance-Tipp für Benutzer von Windows 11"}
Erwägen Sie, Ihre Projekte auf einem [Dev Drive](https://learn.microsoft.com/en-us/windows/dev-drive/) zu speichern. Dev Drives sind für Entwicklungs-Workloads optimiert und können die Build-Zeiten und Datenträgerzugriffe gegenüber regulären NTFS-Laufwerken um bis zu 30 % verbessern.

@end

## Projekt erstellen

@steps
### Projekt generieren
```bash
wails3 init -n myapp
cd myapp
```

Dadurch wird ein neues Projekt mit der Standardvorlage Vanilla + Vite erstellt (HTML/CSS/TypeScript mit dem Vite-Bundler).

@note{type="tip" title="Weitere Vorlagen"}
Verwenden Sie `-t react`, `-t vue` oder `-t svelte` für Ihr bevorzugtes Framework. Diese verwenden standardmäßig TypeScript. Für reines JavaScript verwenden Sie `-t vanilla-js` oder `-t react-js`. Führen Sie `wails3 init -l` aus, um alle verfügbaren Vorlagen anzuzeigen, oder [verwenden Sie Ihr eigenes Frontend-Framework](/guides/dev/frontend-frameworks/).

@end

### Projektstruktur verstehen
```
myapp/
├── main.go              # Application entry point
├── greetservice.go      # Greet service
├── frontend/            # Your UI code
│   ├── index.html       # HTML entry point
│   ├── src/
│   │   └── main.ts      # Frontend TypeScript
│   ├── public/
│   │   └── style.css    # Styles
│   ├── package.json     # Frontend dependencies
│   ├── tsconfig.json    # TypeScript configuration
│   └── vite.config.ts   # Vite bundler config
├── build/               # Build configuration
└── Taskfile.yml         # Build tasks
```

### App ausführen
```bash
wails3 dev
```

@note{type="info" title="Erster Start"}
Der erste Start kann länger als erwartet dauern, da dabei Frontend-Abhängigkeiten installiert, Bindings generiert und weitere Schritte ausgeführt werden. Nachfolgende Starts sind deutlich schneller.

@end

Die App öffnet sich mit einer Begrüßungsoberfläche. Geben Sie Ihren Namen ein und klicken Sie auf „Greet“. Das Go-Backend verarbeitet Ihre Eingabe und gibt eine Begrüßung zurück.

@end

## Funktionsweise

Sehen wir uns den Code an, der diese Funktion ermöglicht.

### Das Go-Backend

Öffnen Sie `greetservice.go`:

```go {title="greetservice.go"}
package main

import (
	"fmt"
)

type GreetService struct{}

func (g *GreetService) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
```

**Wichtige Konzepte:**

1. **Service** – eine Go-Struktur mit exportierten Methoden
2. **Exportierte Methode** – `Greet` beginnt mit einem Großbuchstaben und ist dadurch für das Frontend verfügbar
3. **Einfache Logik** – nimmt einen Namen entgegen und gibt eine Begrüßung zurück
4. **Typsicherheit** – Ein- und Ausgabetypen sind definiert

@note{type="tip" title="Services und Bindings verstehen"}
**Services** sind eigenständige Go-Module, die dem Frontend Funktionen bereitstellen. Es handelt sich um reguläre Go-Strukturen mit exportierten Methoden, die Sie im Feld `Services` Ihrer Anwendungskonfiguration registrieren.

**Bindings** sind das automatisch generierte TypeScript-/JavaScript-SDK, über das Ihr Frontend diese Services aufrufen kann. Wenn Sie `wails3 dev` oder `wails3 build` ausführen, analysiert Wails Ihre registrierten Services und generiert typsichere Bindings in `frontend/bindings/`.

Betrachten Sie Services als Ihre Backend-API und Bindings als die Clientbibliothek, die mit ihr kommuniziert.

@end

### Service registrieren

Öffnen Sie `main.go` und suchen Sie die Service-Registrierung:

```go {title="main.go" highlight="4-6"}
err := application.New(application.Options{
    Name: "myapp",
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
    // ... other options
})
```

Dadurch wird `GreetService` bei Wails registriert, sodass alle exportierten Methoden für das Frontend verfügbar sind.

### Das Frontend

Öffnen Sie `frontend/src/main.js`:

```javascript {title="frontend/src/main.js"}
import {GreetService} from "../bindings/changeme";

window.greet = async () => {
    const nameElement = document.getElementById('name');
    const resultElement = document.getElementById('result');

    const name = nameElement.value;
    if (!name) {
        return;
    }

    try {
        const result = await GreetService.Greet(name);
        resultElement.innerText = result;
    } catch (err) {
        console.error(err);
    }
};
```

**Wichtige Konzepte:**

1. **Automatisch generierte Bindings** – `GreetService` wird aus dem generierten Code importiert
2. **Typsichere Aufrufe** – Methodennamen und Signaturen entsprechen Ihrem Go-Code
3. **Standardmäßig asynchron** – alle Go-Aufrufe geben Promises zurück
4. **Fehlerbehandlung** – Fehler aus Go werden mit try/catch abgefangen

@note{type="info" title="Wo befinden sich die Bindings?"}
Die generierten Bindings befinden sich in `frontend/bindings/`. Sie werden automatisch erstellt, wenn Sie `wails3 dev` oder `wails3 build` ausführen.

**Bearbeiten Sie diese Dateien niemals manuell** – sie werden bei jedem Build neu generiert.

@end

## App anpassen

Fügen wir eine neue Funktion hinzu, um den Arbeitsablauf nachzuvollziehen.

### Funktion „Greet Many“ hinzufügen

@steps
### Methode zu GreetService hinzufügen
Fügen Sie Folgendes zu `greetservice.go` hinzu:

```go {title="greetservice.go"}
func (g *GreetService) GreetMany(names []string) []string {
    greetings := make([]string, len(names))
    for i, name := range names {
        greetings[i] = fmt.Sprintf("Hello %s!", name)
    }
    return greetings
}
```

### Die App wird automatisch neu gebaut
Speichern Sie die Datei. `wails3 dev` baut Ihren Go-Code automatisch neu und startet die App neu.

@note{type="info" title="Automatischer Neuaufbau"}
Änderungen am Go-Code lösen automatisch einen Neuaufbau und Neustart aus. Änderungen am Frontend werden ohne Neustart per Hot Reload geladen.

@end

### Im Frontend verwenden
Fügen Sie Folgendes zu `frontend/src/main.js` hinzu:

```javascript {title="frontend/src/main.js"}
window.greetMany = async () => {
    const names = ['Alice', 'Bob', 'Charlie'];
    const greetings = await GreetService.GreetMany(names);
    console.log(greetings);
};
```

Öffnen Sie die Browserkonsole und rufen Sie `greetMany()` auf. Daraufhin wird das Array mit den Begrüßungen angezeigt.

@end

## Für den Produktiveinsatz bauen

Wenn Sie Ihre App verteilen möchten:

```bash
wails3 build
```

**Das geschieht dabei:**

- Kompiliert den Go-Code mit Optimierungen
- Baut das Frontend für den Produktiveinsatz (minifiziert)
- Erstellt eine native ausführbare Datei in `bin/`

@tabs{sync-key="os"}
[Windows]
**Ausgabe:** `bin/myapp.exe`

Zum Ausführen doppelklicken. Es sind keine Abhängigkeiten erforderlich (WebView2 ist Bestandteil von Windows).

[macOS]
**Ausgabe:** `bin/myapp.app`

In den Ordner „Programme“ ziehen oder zum Ausführen doppelklicken.

[Linux]
**Ausgabe:** `bin/myapp`

Mit `./bin/myapp` ausführen oder eine `.desktop`-Datei für den Anwendungsstarter erstellen.

@end

@note{type="tip" title="Plattformübergreifende Builds"}
Möchten Sie für andere Plattformen bauen? Siehe [Plattformübergreifende Builds →](/guides/build/cross-platform/)

@end

## Was wir gelernt haben

**Projektstruktur**

- `main.go` für das Go-Backend
- `frontend/` für den UI-Code
- `Taskfile.yml` für Build-Aufgaben

**Services**

- Go-Structs mit exportierten Methoden erstellen
- Mit `application.NewService()` registrieren
- Methoden sind automatisch im Frontend verfügbar

**Bindings**

- Automatisch generierte TypeScript-Definitionen
- Typsichere Funktionsaufrufe
- Standardmäßig asynchron (Promises)

**Entwicklungsworkflow**

- `wails3 dev` für Hot Reload
- Go-Änderungen lösen automatisch einen neuen Build und Neustart aus
- Frontend-Änderungen werden sofort per Hot Reload übernommen

---

**Fragen?** Treten Sie [Discord](https://discord.gg/JDdSxwjhGf) bei und fragen Sie die Community.
