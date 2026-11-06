---
title: "Globale Tastenkürzel"
description: "Systemweite Tastenkürzel registrieren, die auch ausgelöst werden, wenn Ihre Anwendung nicht fokussiert ist"
slug: "features/keyboard/global-shortcuts"
sourcePath: "features/keyboard/global-shortcuts.md"
---

Globale Tastenkürzel sind systemweite Tastenkombinationen, die unabhängig davon ausgelöst werden, welche Anwendung gerade fokussiert ist, solange Ihre Wails-Anwendung läuft. Sie eignen sich ideal zum Ein- und Ausblenden, für Schnellerfassungstools, zur Mediensteuerung und für andere Funktionen, auf die Benutzer von überall zugreifen möchten.

@note{type="info" title="Globale Tastenkürzel und Tastenbelegungen"}
[Tastenbelegungen](/features/keyboard/shortcuts/) (`app.KeyBinding`) werden nur ausgelöst, während eines Ihrer Anwendungsfenster fokussiert ist. Globale Tastenkürzel (`app.GlobalShortcut`) werden systemweit ausgelöst, auch wenn Ihre Anwendung im Hintergrund ausgeführt wird. Verwenden Sie die Variante, die Ihren Anforderungen entspricht.

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

`Register` erwartet einen Accelerator und einen Callback. Der Callback wird bei jeder Betätigung des Tastenkürzels in einer eigenen Goroutine ausgeführt.

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
Callbacks werden außerhalb des Haupt-Threads ausgeführt. Wenn Ihr Callback mit Fenstern oder anderen Elementen der Benutzeroberfläche interagieren muss, übernehmen die Fenstermethoden dies für Sie. Benutzerdefinierte Arbeiten im Haupt-Thread müssen Sie jedoch mit `application.InvokeSync` umschließen.

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

`CmdOrCtrl` wird unter macOS zu Command und unter Windows und Linux zu Control aufgelöst. Dadurch lassen sich plattformübergreifende Tastenkürzel bequem definieren.

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

Wails löst diesen Fall selbst auf; das Verhalten ist auf allen Plattformen identisch. Der zweite Aufruf von `Register` gibt einen Fehler zurück, und die ursprüngliche Bindung bleibt bestehen („Fehler melden und beibehalten“). So bleibt das Verhalten vorhersehbar und der Fehler wird sichtbar, statt ein funktionierendes Tastenkürzel stillschweigend zu ersetzen.

```go
app.GlobalShortcut.Register("Ctrl+Shift+G", showWindow)        // ok
err := app.GlobalShortcut.Register("Shift+Ctrl+G", doSomething) // err: already registered
// showWindow is still the active callback for this shortcut.
```

Wenn Sie den Callback eines Tastenkürzels ändern möchten, heben Sie seine Registrierung zunächst mit `Unregister` auf und registrieren Sie es anschließend erneut mit `Register`.

### Eine andere Anwendung besitzt das Tastenkürzel bereits

Dies entscheidet das Betriebssystem, daher ist das Ergebnis plattformspezifisch:

| Plattform | Verhalten, wenn eine andere Anwendung das Tastenkürzel besitzt |
| --- | --- |
| **macOS** | Die Registrierung ist erfolgreich. macOS erlaubt mehreren Anwendungen, denselben Hotkey zu registrieren. Ihr Callback wird daher zusätzlich zum vorhandenen Besitzer hinzugefügt und nicht abgelehnt. |
| **Windows** | Die Registrierung schlägt fehl, und `Register` gibt einen Fehler zurück. Die Anwendung, die das Tastenkürzel zuerst registriert hat, behält es. |
| **Linux (X11)** | Die Registrierung schlägt fehl, und `Register` gibt einen Fehler zurück, da der X-Server ein zweites Grab derselben Kombination ablehnt. |
| **Linux (Wayland)** | Der Compositor entscheidet. Üblicherweise wird der Benutzer im Dialog für globale Tastenkürzel der Desktop-Umgebung aufgefordert, die Bindung zu genehmigen oder auszuwählen. |

Prüfen Sie aufgrund dieser Unterschiede stets den von `Register` zurückgegebenen Fehler. Stellen Sie außerdem ein alternatives Tastenkürzel bereit oder informieren Sie den Benutzer, wenn ein Tastenkürzel nicht beansprucht werden kann.

## Plattformspezifische Hinweise

@tabs
[macOS]
Globale Tastenkürzel verwenden die Hotkey-API des Carbon Event Manager. Dies ist unter macOS der Standardmechanismus für systemweite Hotkeys und erfordert keine Bedienungshilfen-Berechtigung.

Tastenkürzel mit Buchstaben berücksichtigen das Tastaturlayout. Die Carbon-API bindet an physische Tastenpositionen. Daher übersetzt Wails einen Buchstaben in die Taste, die ihn im aktiven Tastaturlayout erzeugt (über `UCKeyTranslate`). `Cmd+A` liegt deshalb bei AZERTY-, QWERTZ-, Dvorak- und ähnlichen Layouts auf der mit **A** beschrifteten Taste und nicht auf der QWERTY-Position. Ziffern, Satzzeichen und benannte Tasten (Pfeiltasten, Funktionstasten usw.) verwenden ihre feste physische Position. Die Bindung wird bei der Registrierung des Tastenkürzels aufgelöst. Wenn Sie das Tastaturlayout bei laufender Anwendung wechseln, behalten vorhandene Tastenkürzel mit Buchstaben die physische Taste, an die sie ursprünglich gebunden wurden.

@note{type="caution" title="Tastenkürzel zum Ausblenden und `ApplicationShouldTerminateAfterLastWindowClosed`"}
`window.Hide()` verwendet unter macOS `orderOut:`, wodurch das Fenster unsichtbar wird. AppKit behandelt das letzte unsichtbare Fenster als geschlossen. Wenn Sie daher `Mac.ApplicationShouldTerminateAfterLastWindowClosed: true` festlegen und Ihr einziges Fenster mit einem globalen Tastenkürzel ausblenden, wird die Anwendung beendet, statt im Hintergrund weiterzulaufen. Lassen Sie diese Option unverändert (Standardeinstellung), wenn Sie einen Hotkey zum Ein- und Ausblenden verwenden. Dann kann das Fenster ausgeblendet und später wieder aufgerufen werden.

@end

[Windows]
Globale Tastenkürzel verwenden die Win32-API `RegisterHotKey`. Die automatische Tastenwiederholung wird unterdrückt. Wenn Sie die Tasten gedrückt halten, wird der Callback daher nur einmal und nicht wiederholt ausgelöst.

Die Registrierung schlägt fehl, wenn eine andere Anwendung die Kombination bereits besitzt. Bevorzugen Sie daher weniger gebräuchliche Kombinationen als Standardwerte.

[Linux]
In **X11**-Sitzungen greift Wails das Tastenkürzel direkt vom X-Server ab. Der angeforderte Accelerator wird daher genau wie angegeben gebunden.

In **Wayland**-Sitzungen kann eine Anwendung konzeptbedingt nicht direkt auf Tasten zugreifen. Wails verwendet stattdessen die Schnittstelle `org.freedesktop.portal.GlobalShortcuts` des XDG Desktop Portal. Beim Portal ist der übergebene Accelerator ein *bevorzugter* Auslöser; der Compositor und letztlich der Benutzer bestimmen die endgültige Tastenkombination. Ihr Callback wird weiterhin ausgelöst, wenn das Tastenkürzel aktiviert wird. Es ist jedoch nicht garantiert, dass die tatsächlichen Tasten Ihrer Anforderung entsprechen, und `IsRegistered`/`GetAll` geben Ihre Anforderung zurück, nicht die vom Compositor eingerichtete Bindung.

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
Einige Kombinationen sind vom Betriebssystem oder von der Desktop-Umgebung reserviert und können von Anwendungen nicht beansprucht werden. Wählen Sie Standardwerte, bei denen Konflikte unwahrscheinlich sind, und behandeln Sie stets den Fehler von `Register`.

@end
