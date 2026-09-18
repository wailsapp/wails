---
title: "Interna der Laufzeitumgebung"
description: "Detaillierter Einblick, wie Wails v3 startet, ausgeführt wird und mit dem Betriebssystem kommuniziert"
slug: "contributing/runtime-internals"
sourcePath: "contributing/runtime-internals.md"
---

Die **Laufzeitumgebung** ist die Schicht, die gewöhnliche Go-Funktionen in eine plattformübergreifende Desktopanwendung umwandelt. Dieses Dokument erläutert die Komponenten, auf die Sie beim Nachverfolgen des Quellcodes stoßen.

---

## 1. Anwendungslebenszyklus

| Phase | Codepfad | Was geschieht |
| --- | --- | --- |
| **Bootstrap** | `pkg/application/application.go:init()` | Registriert Build-Daten und erstellt eine globale `application`-Singleton-Instanz. |
| **New()** | `application.New(...)` | Validiert `Options`, startet den **AssetServer** und initialisiert die Protokollierung. |
| **Run()** | `application.(*App).Run()` | 1. Ruft die plattformspezifische `mainthread.X()` auf, um in den UI-Thread des Betriebssystems einzutreten.<br />2. Startet die **Laufzeitumgebung** (`internal/runtime`).<br />3. Blockiert, bis das letzte Fenster geschlossen oder `Quit()` aufgerufen wird. |
| **Herunterfahren** | `application.(*App).Quit()` | Sendet das Ereignis `application:shutdown` an alle Empfänger, schreibt das Protokoll vollständig und beendet Fenster und Dienste. |

Der Lebenszyklus erlaubt strikt nur einen **einmaligen Einstieg**: Sie können viele Fenster erstellen, das Anwendungsobjekt selbst wird jedoch nur einmal initialisiert.

---

## 2. Fensterverwaltung

### Öffentliche API

```go
win := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "Dashboard",
    Width:  1280,
    Height: 720,
})
win.Show()
```

> `app.Window.New()` akzeptiert **keine Argumente**; verwenden Sie `NewWithOptions(...)`, wenn Sie
>
> eine `application.WebviewWindowOptions`-Struktur als Wert übergeben müssen.

`app.Window.New[WithOptions]()` delegiert an `pkg/application/webview_window_*.go`, wo sich die plattformspezifischen Implementierungen befinden:

```
pkg/application/
├── webview_window_darwin.go    // WKWebView
├── webview_window_linux.go     // GTK + WebKitGTK (plus linux_cgo*.go)
└── webview_window_windows.go   // WebView2
```

Jede Datei:

1. Erstellt eine native Webview (WKWebView, WebKitGTK, WebView2).
2. Registriert einen **Message-Processor**-Callback (`pkg/application/messageprocessor*.go`).
3. Ordnet Wails-Ereignisse (`WindowDidResize`, `WindowFocus`, `WindowFilesDropped`, …) den Konstanten in `pkg/events` zu.

`internal/runtime/` ist für den kleinen Build-Tag-Verbindungscode (`runtime{,_darwin,_linux,_windows,_android,_dev,_prod}.go`) und die eingebettete JS-Laufzeitumgebung unter `internal/runtime/desktop/` vorgesehen.

Aktive Fenster werden von `pkg/application/window_manager.go` / `webview_window.go` nachverfolgt. `pkg/application/screenmanager.go` dient den **Anzeige**-Metadaten (Auflösung, Skalierung, Arbeitsbereich) – es verwaltet keine Fenster.

---

## 3. Nachrichtenverarbeitungspipeline

Die Brücke zwischen JavaScript und Go wird durch die **Message-Processor**-Familie in `pkg/application/messageprocessor_*.go` implementiert.

Ablauf:

1. **JavaScript** ruft `Call.ByID(<fnv-id>, ...args)` aus `/wails/runtime.js` auf (implementiert in `internal/runtime/desktop/@wailsio/runtime/src/calls.ts`) – oder bei Builds im Namensmodus `Call.ByName("pkg.Struct.Method", ...args)`.
2. Die Laufzeithilfsfunktion verpackt den Aufruf und leitet ihn über die native plattformspezifische Brücke an Go weiter.
3. **Go** empfängt die Nachricht in `pkg/application/messageprocessor_call.go`.
4. Der Prozessor sucht die gebundene Methode in `pkg/application/bindings.go` (manuell geschrieben und auf `reflect` basierend) und ruft sie auf.
5. Das Ergebnis oder der Fehler wird an JS zurückserialisiert, wo ein `Promise` erfüllt oder abgelehnt wird.

> Die genaue JSON-Hülle wird durch die Laufzeithilfsfunktion auf der JS-Seite und
>
> durch `messageprocessor_call.go` auf der Go-Seite definiert – in älteren Entwürfen dieser Seite
>
> wurde eine `{"t":"c","id":"123","m":"Greet","p":[…]}`-Struktur angegeben, die jedoch nicht
>
> der aktuellen Implementierung entspricht. Lesen Sie beide Dateien gemeinsam, wenn Sie einen
>
> Fehler im Übertragungsformat untersuchen.

Spezialisierte Prozessoren:

| Datei | Zweck |
| --- | --- |
| `messageprocessor_window.go` | Fensteraktionen (ausblenden, maximieren, …) |
| `messageprocessor_dialog.go` | Native Dialogfelder (`OpenFile`, `MessageBox`, …) |
| `messageprocessor_clipboard.go` | Zwischenablage lesen/schreiben |
| `messageprocessor_events.go` | Ereignisse abonnieren/auslösen |
| `messageprocessor_browser.go` | Browsernavigation, Entwicklertools |

Prozessoren sind **zustandslos** – sie beziehen alles Benötigte aus dem `ApplicationContext`, der mit jeder Nachricht übergeben wird.

---

## 4. Ereignissystem

Ereignisse sind Zeichenketten mit Namensräumen, die über drei Ebenen verteilt werden:

1. **Anwendungsereignisse**: globaler Lebenszyklus (`application:ready`, `application:shutdown`).
2. **Fensterereignisse**: fensterspezifisch (`window:focus`, `window:resize`).
3. **Benutzerdefinierte Ereignisse**: vom Benutzer definiert (`chat:new-message`).

Implementierungsdetails:

- Ereigniskonstanten befinden sich in `pkg/events/` (`defaults.go`, `known_events.go`, `events.txt`). Sie werden von `v3/tasks/events/generate.go` generiert und als `events.Common.*`, `events.Mac.*`, `events.Windows.*` und `events.Linux.*` bereitgestellt. Mit `wails3 generate constants` können Sie sie neu erzeugen.
- Go-Seite (Anwendungsereignisse):
  ```go
  app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {})
  ```

- Go-Seite (Fensterereignisse):
  ```go
  window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {})
  ```

- Go-Seite (benutzerdefinierte Ereignisse):
  ```go
  app.Event.On("chat:new-message", func(e *application.CustomEvent) {})
  ```

- JS-Seite:
  ```js
  import { Events } from "/wails/runtime.js";
  Events.On("chat:new-message", (e) => { /* … */ });
  ```


Anwendungs-, Fenster- und benutzerdefinierte Ereignisse durchlaufen alle `pkg/application/event_manager.go`. Abonnements für Fensterereignisse gelten nur für das jeweilige Fenster. Beim Schließen des Fensters werden seine Handler daher automatisch abgemeldet.

---

## 5. Plattformspezifische Implementierungen

Bedingte Kompilierung hält die öffentliche API identisch und verbirgt zugleich die Besonderheiten der Betriebssysteme.

| Aspekt | Darwin | Linux | Windows |
| --- | --- | --- | --- |
| Hauptthread | `mainthread_darwin.go` (Cgo zu Foundation) | `mainthread_linux.go` (GTK) | `mainthread_windows.go` (Win32 `AttachThreadInput`) |
| Dialoge | `dialogs_darwin.*` (NSAlert) | `dialogs_linux.go` (GtkFileChooser) | `dialogs_windows.go` (IFileOpenDialog) |
| Zwischenablage | `clipboard_darwin.go` | `clipboard_linux.go` | `clipboard_windows.go` |
| Systemtray-Symbole | `systemtray_darwin.*` | `systemtray_linux.go` (DBus) | `systemtray_windows.go` (Shell_NotifyIcon) |

Grundprinzipien:

- **macOS und Windows** verwenden Cgo sparsam (hauptsächlich über `pkg/mac/` und die `w32`-Win32-Wrapper in `pkg/w32`).
- **Linux verwendet zwangsläufig intensiv Cgo**: `pkg/application/linux_cgo.go` (ca. 69 KB) und `linux_cgo_gtk4.{c,go,h}` (ca. 50 KB oder mehr) steuern GTK/WebKitGTK direkt an.
- Verwenden Sie **Build-Tags** (`//go:build darwin`, `//go:build linux`, …), damit die betriebssystemspezifischen Dateien übersichtlich bleiben.
- `internal/capabilities/` ist für plattformspezifische Capability-Flags vorgesehen, das Framework exportiert jedoch **keinen** `ErrCapability`-Sentinelwert. Die Verfügbarkeitsprüfung erfolgt über Rückgaben plattformspezifischer Stubs.

---

## 6. Dateiübersicht

| Datei | Wann Sie sie bearbeiten |
| --- | --- |
| `internal/runtime/runtime_*.go` | Die kleine Stub-Schicht mit Build-Tags ändern (Entwicklung gegenüber Produktion, betriebssystemspezifischer Verbindungscode). |
| `pkg/application/webview_window_*.go` | Einen neuen Fensterhinweis oder ein neues Fensterverhalten implementieren. |
| `pkg/application/messageprocessor*.go` | Einen neuen, von JS aufrufbaren Bridge-Befehl hinzufügen. |
| `pkg/events/*.go` | Die integrierten Ereignisdefinitionen erweitern (anschließend `wails3 generate constants` erneut ausführen). |
| `internal/assetserver/*` | Die Asset-Verarbeitung für Entwicklung und Produktion anpassen. |
| `internal/runtime/desktop/@wailsio/runtime/src/*` | Die eingebettete JS-Laufzeit bearbeiten (Aufruf-/Ereignisweiterleitung, Dialoge, Ziehen, …). |

---

## 7. Tipps zur Fehlersuche

- Konfigurieren Sie `Options.LogLevel` (z. B. `slog.LevelDebug`) und prüfen Sie die Ausgabe von `Options.Logger`. Eine Umgebungsvariable namens `WAILS_LOG_LEVEL` gibt es nicht.
- Die `wails3 dev`-Flags sind `--config`, `--port`, `-s` (aktiviert HTTPS) sowie das globale `--no-colour`. Ein `-verbose`-Flag gibt es nicht.
- Führen Sie die Anwendung unter macOS mit `lldb --` aus, um Objective-C-Ausnahmen frühzeitig zu erkennen.
- Aktivieren Sie bei Chromium-Problemen unter Windows die WebView2-Debugprotokolle: `set WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS=--remote-debugging-port=9222`

---

## 8. Laufzeit erweitern

1. Optional: Deklarieren Sie neue Capability-Flags in `internal/capabilities/`.
2. Implementieren Sie die Funktion mithilfe von Build-Tags in jeder `pkg/application/*_{darwin,linux,windows}.go`-Variante. Stellen Sie auf Plattformen, die sie nicht unterstützen, einen Stub bereit.
3. Fügen Sie die öffentliche API in `pkg/application` hinzu (Schnittstelle und konkrete `WebviewWindow`-Methoden, Optionsstruktur usw.).
4. Registrieren Sie eine neue Methode des Nachrichtenprozessors (`pkg/application/messageprocessor*.go`) und eine passende Hilfsfunktion in der JS-Laufzeit, falls JS sie aufrufen muss.
5. Wenn du ein Ereignis hinzufügst, deklariere dessen Konstante in `pkg/events/` und führe `wails3 generate constants` aus, um die generierten Dateien zu aktualisieren.

Halte dich an diese Checkliste, damit der plattformübergreifende Vertrag intakt bleibt.

---

## 9. Drag-and-drop

Drag-and-drop für Dateien verwendet auf allen Plattformen einen **JavaScript-zentrierten Ansatz**. Die native Schicht fängt Drag-Ereignisse des Betriebssystems ab, die eigentliche Verarbeitung des Ablegens und die DOM-Interaktion erfolgen jedoch in JavaScript.

### Ablauf

1. Der Benutzer zieht Dateien aus dem Betriebssystem über das Wails-Fenster
2. Die native Schicht erkennt den Drag-Vorgang und benachrichtigt JavaScript für Hover-Effekte
3. Der Benutzer legt die Dateien ab
4. Die native Schicht sendet Dateipfade und Koordinaten an JavaScript
5. JavaScript ermittelt das Zielelement für das Ablegen (`data-file-drop-target`)
6. JavaScript sendet Dateipfade und Elementdetails an das Go-Backend
7. Go löst das Ereignis `WindowFilesDropped` mit dem vollständigen Kontext aus

### Plattformspezifische Implementierungen

| Plattform | Native Schicht | Zentrale Herausforderung |
| --- | --- | --- |
| **Windows** | Integrierte Drag-Unterstützung von WebView2 | Koordinaten in CSS-Pixeln, keine Umrechnung erforderlich |
| **macOS** | Drag-Delegates von NSWindow | Fensterrelative in Webview-relative Koordinaten umrechnen |
| **Linux** | Drag-Signale von GTK3 | Datei-Drag-Vorgänge müssen von internen HTML5-Drag-Vorgängen unterschieden werden |

### Linux: Drag-Typen unterscheiden

Sowohl GTK als auch WebKit wollen Drag-Ereignisse verarbeiten. Entscheidend ist, den Typ des Drag-Ziels zu prüfen:

```c
static gboolean is_file_drag(GdkDragContext *context) {
    GList *targets = gdk_drag_context_list_targets(context);
    for (GList *l = targets; l != NULL; l = l->next) {
        GdkAtom atom = GDK_POINTER_TO_ATOM(l->data);
        gchar *name = gdk_atom_name(atom);
        if (name && g_strcmp0(name, "text/uri-list") == 0) {
            g_free(name);
            return TRUE;  // External file drag
        }
        g_free(name);
    }
    return FALSE;  // Internal HTML5 drag
}
```

Signal-Handler geben bei internen Drag-Vorgängen `FALSE` zurück, damit WebKit sie verarbeitet, und bei Datei-Drag-Vorgängen `TRUE`, damit wir sie selbst verarbeiten.

### Ablegen von Dateien blockieren

Wenn `EnableFileDrop` den Wert `false` hat, müssen wir weiterhin verhindern, dass der Browser zu abgelegten Dateien navigiert. Jede Plattform handhabt dies anders:

- **Windows**: JavaScript ruft bei Drag-Ereignissen `preventDefault()` auf
- **macOS**: JavaScript ruft bei Drag-Ereignissen `preventDefault()` auf\
- **Linux**: GTK-Signal-Handler fangen Datei-Drag-Vorgänge auf nativer Ebene ab und weisen sie zurück

### Wichtige Dateien

| Datei | Zweck |
| --- | --- |
| `pkg/application/linux_cgo.go` | GTK-Signal-Handler für Drag-Vorgänge (C-Code in der cgo-Präambel) |
| `pkg/application/webview_window_darwin.go` | Drag-Delegates für macOS |
| `pkg/application/webview_window_windows.go` | Nachrichtenverarbeitung von WebView2 |
| `internal/runtime/desktop/@wailsio/runtime/src/window.ts` | Drop-Verarbeitung in JavaScript |

### Fehlersuche

- **Linux**: Füge `printf` in den C-Code ein (denke an `fflush(stdout)`)
- **Windows**: Verwende `globalApplication.debug()`
- **JavaScript**: Prüfe die Browserkonsole und aktiviere den Debug-Modus

Häufige Probleme:

1. **Interner HTML5-Drag-Vorgang funktioniert nicht**: Der native Handler fängt ihn ab (gib bei Drag-Vorgängen ohne Dateien `FALSE` zurück)
2. **Hover-Effekte werden nicht angezeigt**: Die JavaScript-Handler werden nicht aufgerufen
3. **Falsche Koordinaten**: Prüfe die Umrechnungen zwischen den Koordinatensystemen

---

Du hast nun einen geführten Rundgang durch die Interna der Runtime absolviert. Nutze dieses Wissen zusammen mit der Übersicht **Codebasisstruktur** und der Dokumentation zum **Asset-Server**, um dich sicher zurechtzufinden, wirkungsvoll beizutragen und erfolgreich zu programmieren!
