---
title: "Verschleierte Builds"
description: "Erstellen Sie Ihre Wails-Anwendung mit Garble, um Ihren Quellcode vor Reverse Engineering zu schützen"
slug: "guides/build/obfuscation"
sourcePath: "guides/build/obfuscation.md"
---

[Garble](https://github.com/burrowers/garble) ist ein Go-Build-Tool, das `go build` ersetzt, um Symbole umzubenennen, Konstanten zu verschleiern und Debug-Informationen aus der resultierenden Binärdatei zu entfernen. Wails v3 bietet über zwei neue Befehle direkte Unterstützung für Garble.

## Voraussetzungen

- **Go 1.26.2 oder höher** – erforderlich für Garble v0.16.0
- **Garble v0.16.0**

```bash
go install mvdan.cc/garble@v0.16.0
```

@note{type="tip"}
Die Go-Mindestversion von Garble ändert sich zwischen den Releases. Wenn Sie eine ältere Go-Toolchain verwenden, suchen Sie vor der Installation auf der [Garble-Releaseseite](https://github.com/burrowers/garble/releases) nach einer Version, die zu Ihrer Toolchain passt.

@end

## Erforderlich: JSON-Tags zu Ihren Servicetypen hinzufügen

Jede Struktur, die eine gebundene Servicemethode zurückgibt oder akzeptiert, muss für jedes exportierte Feld ein explizites JSON-Tag enthalten:

```go
// Without tags — breaks under Garble
type OrderSummary struct {
    ID        int
    Total     float64
    LineItems []LineItem
}

// With tags — safe under Garble
type OrderSummary struct {
    ID        int       `json:"id"`
    Total     float64   `json:"total"`
    LineItems []LineItem `json:"lineItems"`
}
```

@note{type="caution"}
Garble benennt exportierte Strukturfelder um, und Wails übergibt diese Strukturen über einen `interface{}`-Parameter an `json.Marshal`, den Garble nicht statisch nachverfolgen kann. Ein verschleierter Build ohne JSON-Tags wird erfolgreich kompiliert, aber Ihr Frontend erhält zur Laufzeit verstümmelte oder leere Feldnamen. Fügen Sie die Tags hinzu, bevor Sie einen verschleierten Build ausführen.

@end

Die Wails-eigenen Typen – `Screen`, `Rect`, `Point`, `Size`, `EnvironmentInfo`, `OSInfo`, `Capabilities` – sind bereits mit Tags versehen. Sie müssen nur Ihre eigenen Typen mit Tags versehen.

## Mit Verschleierung erstellen

@steps
### Stabile ID-Datei generieren
Führen Sie diesen Befehl immer aus, wenn Sie eine gebundene Servicemethode hinzufügen, umbenennen oder entfernen:

```bash
wails3 generate bindings -obfuscated
```

Dadurch wird `wails_obfuscated.gen.go` im Verzeichnis Ihres Hauptpakets erstellt – committen Sie diese Datei.

### Mit Garble erstellen
```bash
wails3 build --obfuscated
```

Erstellt die Anwendung mit den verschleierten Bindings.

@end

## Zusätzliche Flags an Garble übergeben

Verwenden Sie `--garbleargs`, um Optionen direkt an `garble` weiterzuleiten:

```bash
# Obfuscate string literals and reduce binary size
wails3 build --obfuscated --garbleargs "-literals -tiny"

# Reproducible output — same seed produces the same binary
wails3 build --obfuscated --garbleargs "-seed=deadbeef"
```

Die vollständige Liste der unterstützten Flags finden Sie in der [Garble-Dokumentation](https://github.com/burrowers/garble#flags).

## Für Fortgeschrittene: ID-Datei in ein anderes Paket schreiben

Standardmäßig wird `wails_obfuscated.gen.go` neben Ihr `main`-Paket geschrieben. Wenn Ihr Projekt Services in einem Unterpaket verwaltet, das von `main` importiert wird, können Sie die Datei stattdessen mit `-obfuscated-output` dorthin schreiben:

```bash
wails3 generate bindings -obfuscated -obfuscated-output ./internal/services
```

@note{type="caution"}
Das Zielpaket muss direkt oder transitiv von Ihrem `main`-Paket importiert werden, damit dessen `init()` beim Start ausgeführt wird. Wenn das Paket nicht erreichbar ist, werden die stabilen IDs nie registriert und Binding-Aufrufe schlagen fehl (beispielsweise mit `binding not found`-Fehlern zur Laufzeit).

@end

## Fehlerbehebung

### `garble: command not found`

Garble ist nicht installiert oder `$(go env GOPATH)/bin` befindet sich nicht in Ihrem `PATH`.

```bash
go install mvdan.cc/garble@v0.16.0
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Frontend erhält falsche oder leere Feldwerte

In den Rückgabetypen Ihres Services fehlen `json:"..."`-Tags. Prüfen Sie jede von Ihren gebundenen Methoden zurückgegebene Struktur und fügen Sie jedem exportierten Feld ein explizites Tag hinzu.

### `binding not found`-Fehler in der Browserkonsole

Die stabile ID-Datei fehlt oder wurde nicht mitkompiliert. Prüfen Sie Folgendes:

- `wails_obfuscated.gen.go` ist im Verzeichnis Ihres Hauptpakets vorhanden (oder in dem Verzeichnis, das Sie an `-obfuscated-output` übergeben haben)
- Sie haben `wails3 build --obfuscated` ausgeführt; dieser Befehl fügt das Build-Tag `wails_obfuscated` hinzu
- Wenn Sie `-obfuscated-output` verwendet haben, wird das Zielpaket von `main` importiert

### Windows Defender stuft den Build als Virus ein

Mit Garble verschleierte Go-Binärdateien werden während des Builds von Windows Defender aufgrund heuristischer Verfahren als schädlich eingestuft, weil ihnen Debug-Symbole fehlen und sie gepackten ausführbaren Dateien ähneln. Der Build schlägt mit folgender Meldung fehl:

```
open C:\Users\...\AppData\Local\Temp\go-build...\a.out.exe: The file contains a virus or potentially unwanted software.
```

Fügen Sie Ihr temporäres Verzeichnis (in das Go temporäre Build-Artefakte schreibt) und Ihr Projektverzeichnis zur Ausschlussliste von Defender hinzu:

```powershell
Add-MpPreference -ExclusionPath "$env:TEMP"
Add-MpPreference -ExclusionPath "C:\path\to\your\project"
```

Diese Ausschlüsse gelten nur für die angegebenen Pfade und deaktivieren Defender nicht global.

### Build schlägt mit `unsupported Go version` fehl

Garble v0.16.0 erfordert Go 1.26.2 oder höher. Aktualisieren Sie entweder Go oder suchen Sie auf der [Garble-Releaseseite](https://github.com/burrowers/garble/releases) nach einer mit Ihrer Toolchain kompatiblen Version.
