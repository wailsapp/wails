---
title: "Autostart"
description: "Registrieren Sie Ihre Anwendung so, dass sie unter macOS, Windows und Linux bei der Benutzeranmeldung gestartet wird"
slug: "features/autostart/basics"
sourcePath: "features/autostart/basics.md"
---

## Autostart

`app.Autostart` registriert Ihre Anwendung so, dass sie automatisch gestartet wird, wenn sich der Benutzer anmeldet. Dabei wird für jede Plattform der passende native Mechanismus ausgewählt und über symbolische Links referenzierte Installationspfade (Homebrew, Scoop) werden aufgelöst, damit Registrierungen bei einem Upgrade der Binärdatei nicht ungültig werden.

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

`Enable` kann gefahrlos wiederholt aufgerufen werden – die Registrierung wird jedes Mal überschrieben. Daher können Sie die Funktion bei jedem Start aufrufen, wenn Sie die Einstellung des Benutzers dauerhaft gespeichert haben.

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

Gibt an, ob eine Registrierung vorhanden ist. Schnell – der registrierte Pfad wird nicht validiert.

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
| `Strategy` | `AutostartStrategy` | Gibt an, mit welchem Mechanismus die Anwendung registriert wurde (siehe [Plattformverhalten](#plattformverhalten)). |

## Plattformverhalten

@tabs{sync-key="platform"}
[macOS]
Je nach Paketierung der Anwendung wird einer von zwei Mechanismen verwendet:

- **macOS 13+, gebündeltes `.app`**: `SMAppService.mainAppService`. Funktioniert für Sandbox-Anwendungen und Builds für den Mac App Store. Keine TCC-Automatisierungsabfrage (der frühere AppleScript-Ansatz löste eine solche Abfrage aus).
- **macOS vor 13 oder nicht gebündelte Binärdatei**: Eine LaunchAgent-plist wird mit `RunAtLoad=true` nach `~/Library/LaunchAgents/<identifier>.plist` geschrieben.

`Status()` gibt `AutostartStrategySMAppService` oder `AutostartStrategyLaunchAgent` zurück, sodass Aufrufer erkennen können, welcher Pfad verwendet wurde.

Wenn die Anwendung von einer nicht gebündelten auf eine gebündelte Version aktualisiert wird, prüft `Status()` beide Pfade und `Disable()` bereinigt sie, damit kein verwaister LaunchAgent weiterhin den alten Build startet.

[Windows]
Unter `HKCU\Software\Microsoft\Windows\CurrentVersion\Run` wird ein Registrierungswert hinzugefügt. Als Wertname wird die Autostart-`Identifier` und als Daten werden der in Anführungszeichen gesetzte Pfad der ausführbaren Datei sowie alle `Arguments` verwendet.

Die Argumente werden gemäß den Regeln von `CommandLineToArgvW` in Anführungszeichen gesetzt (umgekehrte Schrägstriche vor Anführungszeichen werden verdoppelt), sodass Pfade mit Leerzeichen oder Anführungszeichen korrekt hin und zurück konvertiert werden.

`Status().Strategy` gibt `AutostartStrategyRegistryRun` zurück.

[Linux]
Ein XDG-Autostarteintrag wird nach `$XDG_CONFIG_HOME/autostart/<identifier>.desktop` geschrieben (standardmäßig `~/.config/autostart/`) und enthält:

```ini
Type=Application
Hidden=false
X-GNOME-Autostart-enabled=true
Exec=<executable> <arguments>
```

Das Feld `Exec` wird gemäß der [Desktop-Entry-Spezifikation von freedesktop.org](https://specifications.freedesktop.org/desktop-entry-spec/) maskiert – reservierten Zeichen (`"`, `` ` ``, `$`, `\\`) wird ein umgekehrter Schrägstrich vorangestellt, und der Wert wird in doppelte Anführungszeichen gesetzt, wenn er Leerraum enthält.

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
| macOS (gebündelt) | Bundle-Bezeichner der Anwendung, z. B. `com.example.MyApp` |
| macOS (nicht gebündelt) | `wails.autostart.<slug>`, wobei `<slug>` aus `application.Options.Name` abgeleitet wird |
| Windows | Slug von `application.Options.Name` (Kleinbuchstaben, alle Zeichen außer `A-Za-z0-9._-` werden entfernt, Leerzeichen werden zu Bindestrichen) |
| Linux | Derselbe Slug wie unter Windows |

Bezeichner müssen `^[A-Za-z0-9._-]+$` entsprechen und dürfen höchstens 200 Zeichen lang sein. Für macOS wird die umgekehrte DNS-Notation empfohlen, da sie der üblichen Schreibweise von launchd-Labels entspricht.

Wenn `AutostartOptions.Identifier` überschrieben wird, wird derselbe Bezeichner unter Windows als Name des Registrierungswerts und unter Linux als Dateiname für `.desktop` wiederverwendet, sodass eine einzige Zeichenfolge die Registrierung plattformübergreifend identifiziert.

## Erkennung veralteter Registrierungen

`Disable()` und `Status()` finden die Registrierung, indem sie **den Pfad der registrierten ausführbaren Datei mit `os.Executable()` (nach Auflösung aller symbolischen Links)** abgleichen, nicht durch Nachschlagen des Bezeichners. Das bedeutet:

- **Der Bezeichner kann zwischen Releases gefahrlos geändert werden.** Die alte Registrierung kann weiterhin von `Status()` gefunden und von `Disable()` bereinigt werden – solange der Pfad der ausführbaren Datei unverändert ist.
- **Eine zweite Kopie der App unter einem anderen Pfad überschreibt die Registrierung der ersten nicht.** Jeder Speicherort der Binärdatei wird unabhängig verwaltet.
- **Über symbolische Links installierte Apps (Homebrew, Scoop) bleiben stabil.** `filepath.EvalSymlinks` wird vor dem Abgleich auf `os.Executable()` angewendet. Daher hinterlässt ein Homebrew-Upgrade, das das Ziel austauscht, keinen verwaisten Eintrag.

Was dadurch *nicht* abgedeckt ist: Wenn der Benutzer die Binärdatei an einen nicht verwandten Pfad verschiebt oder sie umbenennt, wird die alte Registrierung verwaist (sie verweist auf die nun fehlende Datei). Apps, die von einem stabilen Installationspfad aus bereitgestellt werden, müssen dies nicht berücksichtigen. Apps, die als portable Einzeldatei-Binärdateien bereitgestellt werden, sollten vor dem Verschieben ihrer eigenen Binärdatei `Disable()` aufrufen oder immer über einen stabilen symbolischen Link gestartet werden.

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

Ein vollständiges ausführbares Beispiel mit Schaltflächen für Status, Aktivieren und Deaktivieren befindet sich in [`examples/autostart/`](https://github.com/wailsapp/wails/tree/master/v3/examples/autostart).
