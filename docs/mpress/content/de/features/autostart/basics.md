---
title: "Autostart"
description: "Registrieren Sie Ihre Anwendung für den Start bei der Benutzeranmeldung unter macOS, Windows und Linux"
slug: "features/autostart/basics"
sourcePath: "features/autostart/basics.md"
---

## Autostart

`app.Autostart` registriert Ihre Anwendung so, dass sie beim Anmelden des Benutzers automatisch gestartet wird. Dabei wird für jede Plattform der passende native Mechanismus gewählt und werden symbolische Links in Installationspfaden (Homebrew, Scoop) aufgelöst, damit Registrierungen nach einem Upgrade der Binärdatei weiterhin funktionieren.

Die Registrierung wird erst bei der **nächsten Anmeldung** wirksam, nicht sofort.

## Schnellstart

```go
import "github.com/wailsapp/wails/v3/pkg/application"

// Register to launch at login
if err := app.Autostart.Enable(); err != nil {
    app.Logger.Error("autostart enable failed", "error", err)
}

// Stop launching at login
if err := app.Autostart.Disable(); err != nil {
    app.Logger.Error("autostart disable failed", "error", err)
}

// Check status
enabled, err := app.Autostart.IsEnabled()
```

## API

### `Enable`

Registriert die Anwendung mit den Standardoptionen für den Start bei der Anmeldung.

```go
func (m *AutostartManager) Enable() error
```

`Enable` kann gefahrlos wiederholt aufgerufen werden – die Registrierung wird jedes Mal überschrieben. Wenn Sie die Einstellung des Benutzers dauerhaft gespeichert haben, können Sie die Funktion daher bei jedem Start aufrufen.

### `EnableWithOptions`

Registriert die Anwendung mit benutzerdefinierten Optionen.

```go
func (m *AutostartManager) EnableWithOptions(opts AutostartOptions) error
```

**`AutostartOptions`:**

| Feld | Typ | Beschreibung |
| --- | --- | --- |
| `Identifier` | `string` | Überschreibt die automatisch abgeleitete Registrierungs-ID. Siehe „Bezeichner“ weiter unten. |
| `Arguments` | `[]string` | Zusätzliche Argumente, die beim Start während der Anmeldung an den Pfad der ausführbaren Datei angehängt werden (z. B. `--hidden`). |

### `Disable`

Entfernt die Autostart-Registrierung. Gibt `nil` zurück, wenn die Anwendung nicht registriert war – das Deaktivieren ist idempotent.

```go
func (m *AutostartManager) Disable() error
```

### `IsEnabled`

Gibt an, ob eine Registrierung vorhanden ist. Die Abfrage ist schnell, da der registrierte Pfad nicht validiert wird.

```go
func (m *AutostartManager) IsEnabled() (bool, error)
```

### `Status`

Gibt den vollständigen Registrierungsstatus zurück.

```go
func (m *AutostartManager) Status() (AutostartStatus, error)
```

**`AutostartStatus`:**

| Feld | Typ | Beschreibung |
| --- | --- | --- |
| `Enabled` | `bool` | Gibt an, ob eine Registrierung vorhanden ist. |
| `Path` | `string` | Speicherort des Registrierungsartefakts auf dem Datenträger (plist-Pfad, `.desktop`-Pfad oder Registrierungsunterschlüssel). Leer, wenn `Enabled` false ist. |
| `Strategy` | `AutostartStrategy` | Der Mechanismus, mit dem die Anwendung registriert wurde (siehe [Plattformverhalten](#plattformverhalten)). |

## Plattformverhalten

@tabs{sync-key="platform"}
[macOS]
Je nach Paketierung der Anwendung wird einer von zwei Mechanismen verwendet:

- **macOS 13+, gebündelte `.app`**: `SMAppService.mainAppService`. Funktioniert für Sandbox-Anwendungen und Builds für den Mac App Store. Es erscheint keine TCC-Automatisierungsabfrage, die beim früher verwendeten AppleScript-Ansatz ausgelöst wurde.
- **macOS vor 13 oder nicht gebündelte Binärdatei**: Eine LaunchAgent-plist wird mit `RunAtLoad=true` nach `~/Library/LaunchAgents/<identifier>.plist` geschrieben.

`Status()` gibt `AutostartStrategySMAppService` oder `AutostartStrategyLaunchAgent` zurück, sodass Aufrufer erkennen können, welcher Pfad verwendet wurde.

Wenn die Anwendung von einer nicht gebündelten zu einer gebündelten Version aktualisiert wird, prüft `Status()` beide Pfade und `Disable()` bereinigt sie, damit ein verwaister LaunchAgent nicht weiterhin den alten Build startet.

[Windows]
Unter `HKCU\Software\Microsoft\Windows\CurrentVersion\Run` wird ein Registrierungswert hinzugefügt. Als Wertname wird die Autostart-`Identifier` und als Wert der in Anführungszeichen gesetzte Pfad der ausführbaren Datei zusammen mit etwaigen `Arguments` festgelegt.

Die Argumente werden gemäß den `CommandLineToArgvW`-Regeln in Anführungszeichen gesetzt (umgekehrte Schrägstriche vor Anführungszeichen werden verdoppelt), sodass Pfade mit Leerzeichen oder Anführungszeichen unverändert verarbeitet werden.

`Status().Strategy` gibt `AutostartStrategyRegistryRun` zurück.

[Linux]
Unter `$XDG_CONFIG_HOME/autostart/<identifier>.desktop` (standardmäßig `~/.config/autostart/`) wird ein XDG-Autostart-Eintrag mit folgendem Inhalt geschrieben:

```ini
Type=Application
Hidden=false
X-GNOME-Autostart-enabled=true
Exec=<executable> <arguments>
```

Das Feld `Exec` wird gemäß der [Desktop-Entry-Spezifikation von freedesktop.org](https://specifications.freedesktop.org/desktop-entry-spec/) maskiert: Reservierten Zeichen (`"`, `` ` ``, `$`, `\\`) wird ein umgekehrter Schrägstrich vorangestellt; enthält der Wert Leerraum, wird er in doppelte Anführungszeichen gesetzt.

`Status().Strategy` gibt `AutostartStrategyXDGAutostart` zurück.

[iOS / Android / Server]
Nicht unterstützt. Alle Methoden geben `ErrAutostartNotSupported` zurück. Verwenden Sie `errors.Is(err, application.ErrAutostartNotSupported)`, um dies zuverlässig zu erkennen:

```go
if err := app.Autostart.Enable(); err != nil {
    if errors.Is(err, application.ErrAutostartNotSupported) {
        // hide the toggle in the UI
        return
    }
    // real failure — surface it
}
```

@end

## Bezeichner

Wenn `Options.Identifier` leer ist, wird aus dem Namen Ihrer Anwendung ein Standardwert abgeleitet:

| Plattform | Standardbezeichner |
| --- | --- |
| macOS (gebündelt) | Der Bundle-Bezeichner der Anwendung, z. B. `com.example.MyApp` |
| macOS (nicht gebündelt) | `wails.autostart.<slug>`, wobei `<slug>` aus `application.Options.Name` abgeleitet wird |
| Windows | Slug von `application.Options.Name` (Kleinbuchstaben, alle Zeichen außer `A-Za-z0-9._-` entfernt, Leerzeichen werden zu Bindestrichen) |
| Linux | Derselbe Slug wie unter Windows |

Bezeichner müssen `^[A-Za-z0-9._-]+$` entsprechen und dürfen höchstens 200 Zeichen lang sein. Für macOS wird die umgekehrte DNS-Notation empfohlen, da sie der üblichen Schreibweise von launchd-Labels entspricht.

Wenn `AutostartOptions.Identifier` überschrieben wird, dient derselbe Bezeichner unter Windows als Name des Registrierungswerts und unter Linux als Name der `.desktop`-Datei. So identifiziert eine einzige Zeichenfolge die Registrierung plattformübergreifend.

## Erkennung veralteter Registrierungen

`Disable()` und `Status()` finden die Registrierung, indem sie **den registrierten Pfad der ausführbaren Datei mit `os.Executable()` vergleichen, nachdem alle symbolischen Links aufgelöst wurden**, und nicht durch Nachschlagen des Bezeichners. Das bedeutet:

- **Der Bezeichner kann zwischen Releases gefahrlos geändert werden.** `Status()` kann die alte Registrierung weiterhin finden und `Disable()` kann sie bereinigen, solange der Pfad der ausführbaren Datei unverändert bleibt.
- **Eine zweite Kopie der Anwendung unter einem anderen Pfad überschreibt nicht die Registrierung der ersten Kopie.** Jeder Speicherort der Binärdatei wird unabhängig verwaltet.
- **Über symbolische Links installierte Anwendungen (Homebrew, Scoop) bleiben stabil.** Vor dem Vergleich wird `filepath.EvalSymlinks` auf `os.Executable()` angewendet. Daher hinterlässt ein Homebrew-Upgrade, das das Ziel austauscht, keinen verwaisten Eintrag.

Was dies *nicht* abdeckt: Wenn der Benutzer die Binärdatei an einen anderen, nicht zugehörigen Pfad verschiebt oder sie umbenennt, bleibt die alte Registrierung verwaist zurück (sie verweist auf die nun nicht mehr vorhandene Datei). Anwendungen, die von einem stabilen Installationspfad aus ausgeliefert werden, müssen sich darum nicht kümmern. Als portable Einzeldatei-Binärdateien ausgelieferte Anwendungen sollten vor dem Verschieben ihrer eigenen Binärdatei `Disable()` aufrufen oder immer über einen stabilen symbolischen Link gestartet werden.

## Beispiel

```go
package main

import (
    "errors"

    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "My App",
    })

    // Restore the user's preference on startup
    if userPrefersAutostart() {
        if err := app.Autostart.Enable(); err != nil {
            if !errors.Is(err, application.ErrAutostartNotSupported) {
                app.Logger.Error("autostart", "error", err)
            }
        }
    }

    app.Run()
}
```

Ein vollständiges ausführbares Beispiel mit Schaltflächen zum Anzeigen des Status sowie zum Aktivieren und Deaktivieren befindet sich unter [`examples/autostart/`](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/autostart).
