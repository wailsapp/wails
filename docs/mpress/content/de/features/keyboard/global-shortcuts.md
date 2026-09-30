---
title: "Globale Tastenkürzel"
description: "Systemweite Tastenkürzel registrieren, die auch ausgelöst werden, wenn Ihre Anwendung nicht fokussiert ist"
slug: "features/keyboard/global-shortcuts"
sourcePath: "features/keyboard/global-shortcuts.md"
---

Globale Tastenkürzel sind systemweite Tastenkombinationen, die unabhängig davon ausgelöst werden, welche Anwendung gerade den Fokus hat, solange Ihre Wails-Anwendung ausgeführt wird. Sie eignen sich ideal für Tastenkürzel zum Ein- und Ausblenden, Werkzeuge zur Schnellerfassung, Mediensteuerungen und andere Funktionen, auf die Benutzer von überall aus zugreifen möchten.

@note{type="info" title="Globale Tastenkürzel und Tastenbelegungen"}
[Tastenbelegungen](/features/keyboard/shortcuts/) (`app.KeyBinding`) werden nur ausgelöst, solange eines Ihrer Anwendungsfenster den Fokus hat. Globale Tastenkürzel (`app.GlobalShortcut`) werden systemweit ausgelöst, auch wenn Ihre Anwendung im Hintergrund ausgeführt wird. Verwenden Sie die Variante, die Ihren Anforderungen entspricht.

@end

Globale Tastenkürzel basieren direkt auf den nativen Funktionen der jeweiligen Plattform und fügen keine Abhängigkeiten von Drittanbietern hinzu.

## Auf den Manager für globale Tastenkürzel zugreifen

Der Manager ist über die Eigenschaft `GlobalShortcut` Ihrer Anwendungsinstanz verfügbar:

```go
app := application.New(application.Options{
    Name: "Global Shortcuts Demo",
})

globalShortcuts := app.GlobalShortcut
```

## Ein Tastenkürzel registrieren

`Register` erwartet einen Accelerator und einen Callback. Der Callback wird bei jedem Drücken des Tastenkürzels in einer eigenen Goroutine ausgeführt.

```go
err := app.GlobalShortcut.Register("CmdOrCtrl+Shift+G", func() {
    // Runs even when another application is focused.
    window.Show()
    window.Focus()
})
if err != nil {
    app.Logger.Error("could not register shortcut", "error", err)
}
```

Sie können Tastenkürzel registrieren, bevor Sie `app.Run` aufrufen. Die Bindung an das Betriebssystem erfolgt dann automatisch beim Start der Anwendung.

@note{type="tip" title="Aus einem Callback auf die Benutzeroberfläche zugreifen"}
Callbacks werden außerhalb des Hauptthreads ausgeführt. Wenn Ihr Callback mit Fenstern oder anderen Teilen der Benutzeroberfläche interagieren muss, übernehmen die Fenstermethoden dies für Sie. Benutzerdefinierte Arbeiten im Hauptthread müssen Sie jedoch mit `application.InvokeSync` umschließen.

@end

### Accelerator-Format

Globale Tastenkürzel verwenden dasselbe Accelerator-Format wie Menü-Accelerators und Tastenbelegungen:

```go
"CmdOrCtrl+Shift+G"  // Command on macOS, Control elsewhere
"Ctrl+Alt+K"         // Control + Alt + K
"Cmd+Option+Space"   // Command + Option + Space (macOS)
"Super+D"            // Super / Windows / Logo key + D
"Ctrl+Shift+F5"      // Function keys are supported
```

`CmdOrCtrl` wird unter macOS zu Command und unter Windows und Linux zu Control aufgelöst, was plattformübergreifende Tastenkürzel vereinfacht.

## Tastenkürzel verwalten

```go
// Check whether a shortcut is registered (modifier order does not matter).
registered := app.GlobalShortcut.IsRegistered("Ctrl+Shift+G")

// List every shortcut this application has registered.
for _, accelerator := range app.GlobalShortcut.GetAll() {
    app.Logger.Info("global shortcut", "accelerator", accelerator)
}

// Release a single shortcut.
app.GlobalShortcut.Unregister("Ctrl+Shift+G")

// Release everything (also done automatically on shutdown).
app.GlobalShortcut.UnregisterAll()
```

Alle registrierten Tastenkürzel werden beim Beenden der Anwendung automatisch freigegeben. Sie müssen sie daher nicht manuell bereinigen.

## Was geschieht, wenn dasselbe Tastenkürzel zweimal registriert wird?

Es gibt zwei unterschiedliche Fälle, die Wails verschieden behandelt.

### Dieselbe Anwendung registriert ein Tastenkürzel zweimal

Wails löst diesen Fall selbst auf und verhält sich dabei auf jeder Plattform identisch. Der zweite Aufruf von `Register` gibt einen Fehler zurück, und die ursprüngliche Bindung bleibt bestehen („Fehler melden und beibehalten“). Dadurch bleibt das Verhalten vorhersehbar und der Fehler wird sichtbar, statt ein funktionierendes Tastenkürzel unbemerkt zu ersetzen.

```go
app.GlobalShortcut.Register("Ctrl+Shift+G", showWindow)        // ok
err := app.GlobalShortcut.Register("Shift+Ctrl+G", doSomething) // err: already registered
// showWindow is still the active callback for this shortcut.
```

Wenn Sie den Callback eines Tastenkürzels ändern möchten, führen Sie dafür zuerst `Unregister` und anschließend erneut `Register` aus.

### Eine andere Anwendung besitzt das Tastenkürzel bereits

Darüber entscheidet das Betriebssystem, weshalb das Ergebnis plattformspezifisch ist:

| Plattform | Verhalten, wenn eine andere Anwendung das Tastenkürzel besitzt |
| --- | --- |
| **macOS** | Die Registrierung ist erfolgreich. macOS erlaubt mehreren Anwendungen, dasselbe Tastenkürzel zu registrieren. Ihr Callback wird daher zusätzlich zum bereits vorhandenen Besitzer registriert und nicht abgelehnt. |
| **Windows** | Die Registrierung schlägt fehl und `Register` gibt einen Fehler zurück. Die Anwendung, die das Tastenkürzel zuerst registriert hat, behält es. |
| **Linux (X11)** | Die Registrierung schlägt fehl und `Register` gibt einen Fehler zurück, weil der X-Server einen zweiten Grab derselben Tastenkombination ablehnt. |
| **Linux (Wayland)** | Der Compositor entscheidet. Üblicherweise wird der Benutzer aufgefordert, die Bindung im Dialog für globale Tastenkürzel der Desktop-Umgebung zu genehmigen oder auszuwählen. |

Prüfen Sie wegen dieser Unterschiede stets den von `Register` zurückgegebenen Fehler. Stellen Sie außerdem ein alternatives Tastenkürzel bereit oder informieren Sie den Benutzer, wenn ein Tastenkürzel nicht beansprucht werden kann.

## Plattformspezifische Hinweise

@tabs
[macOS]
Globale Tastenkürzel verwenden die Hotkey-API des Carbon Event Managers. Dies ist der Standardmechanismus für systemweite Tastenkürzel unter macOS und erfordert keine Bedienungshilfen-Berechtigung.

Hotkeys werden an physische Tastenpositionen gebunden. Bei Tastaturlayouts, die nicht QWERTY entsprechen, wird ein Tastenkürzel daher der Taste an der entsprechenden Position im standardmäßigen ANSI/QWERTY-Layout zugeordnet.

@note{type="caution" title="Tastenkürzel zum Ausblenden und `ApplicationShouldTerminateAfterLastWindowClosed`"}
`window.Hide()` verwendet unter macOS `orderOut:`, wodurch das Fenster unsichtbar wird. AppKit behandelt das letzte unsichtbare Fenster als geschlossen. Wenn Sie daher `Mac.ApplicationShouldTerminateAfterLastWindowClosed: true` festlegen und Ihr einziges Fenster mit einem globalen Tastenkürzel ausblenden, wird die Anwendung beendet, statt im Hintergrund weiterzulaufen. Lassen Sie diese Option deaktiviert (Standardeinstellung), wenn Sie ein Tastenkürzel zum Ein- und Ausblenden verwenden, damit das Fenster ausgeblendet und später wieder aufgerufen werden kann.

@end

[Windows]
Globale Tastenkürzel verwenden die Win32-API `RegisterHotKey`. Die automatische Tastenwiederholung wird unterdrückt, sodass der Callback beim Gedrückthalten der Tasten nur einmal und nicht wiederholt ausgelöst wird.

Die Registrierung schlägt fehl, wenn eine andere Anwendung die Tastenkombination bereits besitzt. Bevorzugen Sie daher für die Standardbelegung weniger gebräuchliche Kombinationen.

[Linux]
In **X11**-Sitzungen registriert Wails das Tastenkürzel direkt beim X-Server, sodass der angeforderte Tastaturbeschleuniger genau wie angegeben gebunden wird.

In **Wayland**-Sitzungen können Anwendungen grundsätzlich keine Tasten direkt registrieren. Wails verwendet stattdessen die `org.freedesktop.portal.GlobalShortcuts`-Schnittstelle des XDG Desktop Portal. Beim Portal ist der übergebene Tastaturbeschleuniger ein *bevorzugter* Auslöser; der Compositor und letztlich der Benutzer legen die endgültige Tastenkombination fest. Ihr Callback wird weiterhin ausgeführt, wenn das Tastenkürzel aktiviert wird. Es ist jedoch nicht garantiert, dass die tatsächlichen Tasten Ihrer Anforderung entsprechen, und `IsRegistered`/`GetAll` melden die von Ihnen angeforderte statt der vom Compositor gebundenen Kombination.

Das Portal erfordert eine Desktop-Umgebung, die das Portal für globale Tastenkürzel implementiert, beispielsweise eine aktuelle Version von GNOME oder KDE Plasma.

@end

## Vollständiges Beispiel

```go
package main

import (
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "Global Shortcuts Demo",
    })

    window := app.Window.New()

    // Bring the window to the front from anywhere.
    if err := app.GlobalShortcut.Register("CmdOrCtrl+Shift+G", func() {
        window.Show()
        window.Focus()
    }); err != nil {
        log.Printf("could not register show shortcut: %v", err)
    }

    // Hide the window from anywhere.
    if err := app.GlobalShortcut.Register("CmdOrCtrl+Shift+H", func() {
        window.Hide()
    }); err != nil {
        log.Printf("could not register hide shortcut: %v", err)
    }

    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

@note{type="danger" title="Kritische Systemtastenkürzel vermeiden"}
Einige Tastenkombinationen sind dem Betriebssystem oder der Desktop-Umgebung vorbehalten und können nicht von Anwendungen registriert werden. Wählen Sie Standardeinstellungen, bei denen Konflikte unwahrscheinlich sind, und behandeln Sie immer den von `Register` zurückgegebenen Fehler.

@end
