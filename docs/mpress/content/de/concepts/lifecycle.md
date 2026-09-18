---
title: "Anwendungslebenszyklus"
description: "Den Lebenszyklus einer Wails-Anwendung vom Start bis zum Herunterfahren verstehen"
slug: "concepts/lifecycle"
sourcePath: "concepts/lifecycle.md"
---

## Anwendungslebenszyklus verstehen

Desktopanwendungen durchlaufen einen Lebenszyklus vom Start bis zum Herunterfahren. Wails v3 bietet **Dienste**, **Ereignisse** und **Hooks**, um diesen Lebenszyklus effektiv zu verwalten.

## Phasen des Lebenszyklus

```d2
direction: down

Start: Anwendungsstart {
  shape: oval
  style.fill: "#10B981"
}

Init: Initialisierung {
  Parse: Optionen auswerten {
    shape: rectangle
  }
  Register: Dienste registrieren {
    shape: rectangle
  }
  Setup: Laufzeitumgebung einrichten {
    shape: rectangle
  }
}

AppRun: app.Run() {
  shape: rectangle
  style.fill: "#3B82F6"
}

ServiceStartup: Dienststart {
  shape: rectangle
  style.fill: "#8B5CF6"
}

EventLoop: Ereignisschleife {
  Process: Ereignisse verarbeiten {
    shape: rectangle
  }
  Handle: Nachrichten verarbeiten {
    shape: rectangle
  }
  Update: Benutzeroberfläche aktualisieren {
    shape: rectangle
  }
}

QuitSignal: Beendigungssignal {
  shape: diamond
  style.fill: "#F59E0B"
}

ShouldQuit: ShouldQuit-Prüfung {
  shape: rectangle
  style.fill: "#3B82F6"
}

OnShutdown: OnShutdown-Callbacks {
  shape: rectangle
  style.fill: "#3B82F6"
}

ServiceShutdown: Dienste herunterfahren {
  shape: rectangle
  style.fill: "#8B5CF6"
}

Cleanup: Bereinigung {
  Close: Fenster schließen {
    shape: rectangle
  }
  Release: Ressourcen freigeben {
    shape: rectangle
  }
}

End: Anwendungsende {
  shape: oval
  style.fill: "#EF4444"
}

Start -> Init.Parse
Init.Parse -> Init.Register
Init.Register -> Init.Setup
Init.Setup -> AppRun
AppRun -> ServiceStartup
ServiceStartup -> EventLoop.Process
EventLoop.Process -> EventLoop.Handle
EventLoop.Handle -> EventLoop.Update
EventLoop.Update -> EventLoop.Process: Schleife
EventLoop.Process -> QuitSignal: Benutzer beendet die Anwendung
QuitSignal -> ShouldQuit: Beenden zulässig?
ShouldQuit -> EventLoop.Process: Abgelehnt
ShouldQuit -> OnShutdown: Zulässig
OnShutdown -> ServiceShutdown
ServiceShutdown -> Cleanup.Close
Cleanup.Close -> Cleanup.Release
Cleanup.Release -> End
```

### 1. Anwendung erstellen

Erstellen Sie Ihre Anwendung mit `application.New()`:

```go
app := application.New(application.Options{
    Name:        "My App",
    Description: "An application built with Wails",
    Services: []application.Service{
        application.NewService(&MyService{}),
    },
    Assets: application.AssetOptions{
        Handler: application.BundledAssetFileServer(assets),
    },
})
```

**Was geschieht:**

1. Optionen werden analysiert und validiert
2. Dienste werden registriert (aber noch nicht gestartet)
3. Der Asset-Server wird konfiguriert
4. Die Laufzeitumgebung wird eingerichtet

### 2. Anwendung ausführen

Rufen Sie `app.Run()` auf, um die Anwendung zu starten:

```go
err := app.Run()  // Blocks until quit
if err != nil {
    log.Fatal(err)
}
```

**Was geschieht:**

1. Dienste werden in der Reihenfolge ihrer Registrierung gestartet
2. Ereignis-Listener werden aktiviert
3. Fenster können erstellt werden
4. Die Ereignisschleife beginnt

### 3. Ereignisschleife

Die Anwendung wechselt in die Ereignisschleife, in der sie die meiste Zeit verbringt:

- Betriebssystemereignisse werden verarbeitet (Maus-, Tastatur- und Fensterereignisse)
- Go-zu-JS-Nachrichten werden verarbeitet
- JS-zu-Go-Aufrufe werden ausgeführt
- Aktualisierungen der Benutzeroberfläche werden gerendert

### 4. Herunterfahren

Beim Beenden der Anwendung:

1. Der Callback `ShouldQuit` wird geprüft (falls festgelegt)
2. `OnShutdown`-Callbacks werden ausgeführt
3. Dienste werden in umgekehrter Reihenfolge heruntergefahren
4. Fenster werden geschlossen
5. Ressourcen werden freigegeben

## Lebenszyklus von Diensten

Dienste sind in Wails v3 der wichtigste Mechanismus zur Verwaltung des Lebenszyklus. Sie stellen über Schnittstellen Hooks für das Starten und Herunterfahren bereit. Eine vollständige Dokumentation zu Diensten finden Sie im [Leitfaden zu Diensten](/features/bindings/services/).

### Dienst erstellen

```go
type MyService struct {
    db *sql.DB
}

// ServiceStartup is called when the application starts
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    var err error
    s.db, err = sql.Open("sqlite3", "app.db")
    if err != nil {
        return err  // Startup aborts if error returned
    }

    // Run migrations
    if err := s.runMigrations(); err != nil {
        return err
    }

    return nil
}

// ServiceShutdown is called when the application shuts down
func (s *MyService) ServiceShutdown() error {
    if s.db != nil {
        return s.db.Close()
    }
    return nil
}
```

### Dienste registrieren

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&MyService{}),
        application.NewService(&AnotherService{}),
    },
})
```

**Wichtige Punkte:**

- Dienste werden in der Reihenfolge ihrer Registrierung gestartet
- Dienste werden in **umgekehrter** Registrierungsreihenfolge heruntergefahren
- Gibt `ServiceStartup` eines Dienstes einen Fehler zurück, bricht die Anwendung ab
- Der an `ServiceStartup` übergebene `ctx` wird abgebrochen, sobald das Herunterfahren beginnt

### Anwendungskontext verwenden

Der an `ServiceStartup` übergebene Kontext bleibt während der gesamten Lebensdauer der Anwendung gültig:

```go
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    // Start a background task that respects shutdown
    go func() {
        ticker := time.NewTicker(5 * time.Minute)
        defer ticker.Stop()

        for {
            select {
            case <-ticker.C:
                s.performBackgroundSync()
            case <-ctx.Done():
                // Application is shutting down
                return
            }
        }
    }()

    return nil
}
```

Sie können auch über die Anwendungsinstanz auf den Kontext zugreifen:

```go
app := application.Get()
ctx := app.Context()
```

## Hooks auf Anwendungsebene

Diese praktischen Callbacks in `application.Options` ermöglichen es, sich in den Anwendungslebenszyklus einzuklinken, ohne einen vollständigen Service zu erstellen. Sie eignen sich für einfache Bereinigungsaufgaben, die Bestätigung des Beendens oder zum Ausführen von Code an bestimmten Stellen der Herunterfahrsequenz.

Verwenden Sie für eine komplexere Verwaltung des Lebenszyklus mit Startlogik, Dependency Injection oder zustandsbehafteten Ressourcen stattdessen [Services](#lebenszyklus-von-diensten).

### ShouldQuit

Der Callback `ShouldQuit` wird bei jeder Anforderung zum Beenden aufgerufen – unabhängig davon, ob der Benutzer das letzte Fenster schließt, Cmd+Q (macOS) bzw. Alt+F4 (Windows) drückt oder `app.Quit()` programmgesteuert aufgerufen wird.

**Rückgabewert:**

- Geben Sie `true` zurück, um das Beenden zuzulassen (die Anwendung wird heruntergefahren).
- Geben Sie `false` zurück, um das Beenden abzubrechen (die Anwendung wird weiter ausgeführt).

Hier können Sie Anforderungen zum Beenden abfangen und bei Bedarf verhindern, beispielsweise um den Benutzer auf ungespeicherte Änderungen hinzuweisen:

```go
app := application.New(application.Options{
    ShouldQuit: func() bool {
        if !hasUnsavedChanges() {
            return true // No unsaved changes, allow quit
        }

        // Prompt the user — MessageDialog.Show() blocks and returns nothing.
        // The button's OnClick callback fires for whichever button the user picks.
        shouldQuit := false

        dlg := application.Get().Dialog.Question().
            SetTitle("Unsaved Changes").
            SetMessage("You have unsaved changes. Quit anyway?")

        quit := dlg.AddButton("Quit")
        cancel := dlg.AddButton("Cancel")
        dlg.SetDefaultButton(cancel)
        dlg.SetCancelButton(cancel)

        quit.OnClick(func() { shouldQuit = true })

        dlg.Show()
        return shouldQuit
    },
})
```

Ist `ShouldQuit` nicht festgelegt, wird die Anwendung bei einer entsprechenden Anforderung sofort beendet.

**Wann ShouldQuit aufgerufen wird:**

- Der Benutzer schließt das letzte Fenster (sofern `DisableQuitOnLastWindowClosed` nicht festgelegt ist).
- Der Benutzer drückt unter macOS Cmd+Q.
- Der Benutzer drückt unter Windows Alt+F4 (während das letzte Fenster den Fokus hat).
- Der Code ruft `app.Quit()` auf.

**Wann ShouldQuit NICHT aufgerufen wird:**

- Der Prozess wird zwangsweise beendet (SIGKILL oder erzwungenes Beenden über den Task-Manager).
- `os.Exit()` wird direkt aufgerufen.

### OnShutdown

Der Callback `OnShutdown` wird aufgerufen, sobald das Beenden der Anwendung bestätigt ist (nachdem `ShouldQuit` den Wert `true` zurückgegeben hat, sofern festgelegt). Verwenden Sie ihn für Bereinigungsaufgaben wie das Speichern des Zustands, das Schließen von Datenbankverbindungen oder das Freigeben von Ressourcen.

```go
app := application.New(application.Options{
    OnShutdown: func() {
        // Save application state
        saveState()

        // Close connections
        cleanup()
    },
})
```

Sie können außerdem jederzeit während des Anwendungslebenszyklus programmgesteuert weitere Callbacks für das Herunterfahren registrieren:

```go
app.OnShutdown(func() {
    log.Println("Application shutting down...")
})
```

Mehrere Callbacks werden in der Reihenfolge ihrer Registrierung ausgeführt. Der Herunterfahrvorgang wird blockiert, bis alle Callbacks abgeschlossen sind.

**Wichtig:** Halten Sie Callbacks für das Herunterfahren kurz (unter 1 Sekunde). Das Betriebssystem kann Anwendungen zwangsweise beenden, wenn der Beendigungsvorgang zu lange dauert. Dadurch könnte die Bereinigung unterbrochen werden und es könnten Daten verloren gehen.

### PostShutdown

Der Callback `PostShutdown` wird unmittelbar vor dem Beenden des Prozesses aufgerufen, nachdem alle Aufgaben zum Herunterfahren abgeschlossen sind. Zu diesem Zeitpunkt kann die Anwendungsinstanz nicht mehr verwendet werden: Die Fenster sind geschlossen, die Services heruntergefahren und die Ressourcen freigegeben.

Dies ist hauptsächlich für Folgendes nützlich:

- Abschließende Protokollierung, die erst nach allen anderen Bereinigungsaufgaben erfolgen darf
- Testen und Debuggen des Verhaltens beim Herunterfahren
- Plattformen, auf denen `app.Run()` nicht zurückkehrt (der Callback stellt sicher, dass Ihr Code ausgeführt wird)

```go
app := application.New(application.Options{
    PostShutdown: func() {
        // Final logging
        log.Println("Application terminated cleanly")

        // Flush any buffered logs
        logger.Sync()
    },
})
```

**Hinweis:** Versuchen Sie in `PostShutdown` nicht, Anwendungsfunktionen wie Fenster oder Dialoge zu verwenden – sie sind nicht mehr verfügbar.

## Ereignisbasierter Lebenszyklus

Wails stellt ein Ereignissystem bereit, das Sie über Vorgänge in Ihrer Anwendung informiert, etwa über das Öffnen von Fenstern, den Start der Anwendung, Änderungen des Farbschemas und weitere Ereignisse. Sie können auf diese Ereignisse reagieren, um Änderungen des Lebenszyklus zu verarbeiten, ohne sie zu blockieren oder abzufangen.

Bei Fensterereignissen können Sie statt `OnWindowEvent` auch `RegisterHook` verwenden, um Aktionen abzufangen und abzubrechen, beispielsweise um das Schließen eines Fensters zu verhindern. Siehe unten [Fenster-Hooks](#fenster-hooks-abbrechbare-ereignisse).

Die vollständige Dokumentation des Ereignissystems finden Sie im [Ereignisleitfaden](/features/events/system/).

### Anwendungsereignisse

Reagieren Sie auf Ereignisse des Anwendungslebenszyklus:

```go
app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(event *application.ApplicationEvent) {
    app.Logger.Info("Application has started!")
})
```

Außerdem sind plattformspezifische Ereignisse verfügbar:

```go
// macOS
app.Event.OnApplicationEvent(events.Mac.ApplicationDidFinishLaunching, func(event *application.ApplicationEvent) {
    // Handle macOS launch
})

app.Event.OnApplicationEvent(events.Mac.ApplicationWillTerminate, func(event *application.ApplicationEvent) {
    // Handle macOS termination
})

// Windows
app.Event.OnApplicationEvent(events.Windows.ApplicationStarted, func(event *application.ApplicationEvent) {
    // Handle Windows start
})
```

### Fensterereignisse

Reagieren Sie auf Ereignisse des Fensterlebenszyklus:

```go
window := app.Window.New()

window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window gained focus")
})

window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
    app.Logger.Info("Window is closing")
})
```

### Fenster-Hooks (abbrechbare Ereignisse)

Verwenden Sie `RegisterHook` statt `OnWindowEvent`, wenn Sie ein Ereignis **abbrechen** müssen:

```go
window := app.Window.New()

var countdown = 3

window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    countdown--
    if countdown > 0 {
        app.Logger.Info("Not closing yet!", "remaining", countdown)
        e.Cancel()  // Prevent the window from closing
        return
    }
    app.Logger.Info("Window closing now")
})
```

**Unterschied zwischen OnWindowEvent und RegisterHook:**

- `OnWindowEvent`: Benachrichtigt Sie, wenn ein Ereignis eintritt (es kann nicht abgebrochen werden)
- `RegisterHook`: Ermöglicht es Ihnen, das Ereignis abzufangen und gegebenenfalls abzubrechen

## Fensterlebenszyklus

Fenster haben einen eigenen Lebenszyklus von der Erstellung bis zur Zerstörung. Jedes Fenster lädt seine Frontend-Inhalte unabhängig und kann jederzeit angezeigt, ausgeblendet oder geschlossen werden. Wenn ein Benutzer versucht, ein Fenster zu schließen, können Sie dies mit einem `RegisterHook` abfangen, um eine Bestätigung anzufordern oder das Fenster auszublenden, statt es zu zerstören.

Die vollständige Dokumentation zu Fenstern finden Sie im [Leitfaden zu Fenstern](/features/windows/basics/).

```d2
direction: down

Create: Fenster erstellen {
  shape: oval
  style.fill: "#10B981"
}

Load: Frontend laden {
  shape: rectangle
}

Show: Fenster anzeigen {
  shape: rectangle
}

Active: Fenster aktiv {
  Events: Ereignisse verarbeiten {
    shape: rectangle
  }
}

CloseRequest: Schließanforderung {
  shape: diamond
  style.fill: "#F59E0B"
}

Hook: WindowClosing-Hook {
  shape: rectangle
  style.fill: "#3B82F6"
}

Destroy: Fenster zerstören {
  shape: rectangle
}

End: Fenster geschlossen {
  shape: oval
  style.fill: "#EF4444"
}

Create -> Load
Load -> Show
Show -> Active.Events
Active.Events -> Active.Events: Schleife
Active.Events -> CloseRequest: Benutzer schließt das Fenster
CloseRequest -> Hook
Hook -> Active.Events: Abgebrochen
Hook -> Destroy: Zulässig
Destroy -> End
```

### Fenster erstellen

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "My Window",
    Width:  800,
    Height: 600,
})
```

### Schließen eines Fensters verhindern

```go
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    if !hasUnsavedChanges() {
        return
    }

    // MessageDialog.Show() returns nothing; per-button OnClick handlers fire.
    dlg := application.Get().Dialog.Question().
        SetTitle("Unsaved Changes").
        SetMessage("Save before closing?")

    save := dlg.AddButton("Save")
    discard := dlg.AddButton("Discard")
    cancel := dlg.AddButton("Cancel")
    dlg.SetDefaultButton(save)
    dlg.SetCancelButton(cancel)

    save.OnClick(func() { saveChanges() })
    cancel.OnClick(func() { e.Cancel() }) // Prevent close
    _ = discard                            // "Discard" falls through and allows close

    dlg.Show()
})
```

### Ausblenden statt schließen

Ein gängiges Muster für Anwendungen im Infobereich:

```go
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    window.Hide()  // Hide instead of destroy
    e.Cancel()     // Prevent actual close
})
```

## Lebenszyklus bei mehreren Fenstern

Bei mehreren Fenstern:

```go
mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Main Window",
})

settingsWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "Settings",
    Width:  400,
    Height: 600,
    Hidden: true,  // Start hidden
})
```

**Das Standardverhalten variiert je nach Plattform:**

| Plattform | Standardverhalten beim Schließen des letzten Fensters |
| --- | --- |
| macOS | Anwendung wird weiter ausgeführt (Menüleiste bleibt sichtbar) |
| Windows | Anwendung wird beendet |
| Linux | Anwendung wird beendet |

macOS folgt den nativen Plattformkonventionen, nach denen Anwendungen normalerweise auch ohne geöffnete Fenster in der Menüleiste aktiv bleiben. Unter Windows und Linux werden sie standardmäßig beendet.

**Alle Plattformen so konfigurieren, dass die Anwendung beim Schließen des letzten Fensters beendet wird:**

```go
app := application.New(application.Options{
    Mac: application.MacOptions{
        ApplicationShouldTerminateAfterLastWindowClosed: true,
    },
})
```

**Alle Plattformen so konfigurieren, dass die Anwendung beim Schließen des letzten Fensters weiter ausgeführt wird:**

Dies ist für Anwendungen im Infobereich oder Anwendungen nützlich, die im Hintergrund weiter ausgeführt werden sollen.

```go
app := application.New(application.Options{
    Windows: application.WindowsOptions{
        DisableQuitOnLastWindowClosed: true,
    },
    Linux: application.LinuxOptions{
        DisableQuitOnLastWindowClosed: true,
    },
})
```

## Gängige Muster

### Muster 1: Datenbankdienst

```go
type DatabaseService struct {
    db *sql.DB
}

func (s *DatabaseService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    var err error
    s.db, err = sql.Open("sqlite3", "app.db")
    if err != nil {
        return fmt.Errorf("failed to open database: %w", err)
    }

    if err := s.db.PingContext(ctx); err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }

    return nil
}

func (s *DatabaseService) ServiceShutdown() error {
    if s.db != nil {
        return s.db.Close()
    }
    return nil
}

// Exported methods are available to the frontend
func (s *DatabaseService) GetUsers() ([]User, error) {
    // Query implementation
}
```

### Muster 2: Konfigurationsdienst

```go
type ConfigService struct {
    config *Config
    path   string
}

func (s *ConfigService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    s.path = "config.json"

    data, err := os.ReadFile(s.path)
    if err != nil {
        if os.IsNotExist(err) {
            s.config = &Config{} // Default config
            return nil
        }
        return err
    }

    return json.Unmarshal(data, &s.config)
}

func (s *ConfigService) ServiceShutdown() error {
    data, err := json.MarshalIndent(s.config, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(s.path, data, 0644)
}
```

### Muster 3: Hintergrund-Worker

```go
type WorkerService struct {
    cancel context.CancelFunc
}

func (s *WorkerService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    workerCtx, cancel := context.WithCancel(ctx)
    s.cancel = cancel

    go s.runWorker(workerCtx)

    return nil
}

func (s *WorkerService) ServiceShutdown() error {
    if s.cancel != nil {
        s.cancel()
    }
    return nil
}

func (s *WorkerService) runWorker(ctx context.Context) {
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            s.doWork()
        case <-ctx.Done():
            return
        }
    }
}
```

## Lebenszyklusreferenz

| Hook/Schnittstelle | Aufrufzeitpunkt | Abbruch möglich? | Verwendungszweck |
| --- | --- | --- | --- |
| `ServiceStartup` | Während `app.Run()`, vor der Ereignisschleife | Nein (zum Abbrechen Fehler zurückgeben) | Initialisierung |
| `ServiceShutdown` | Während des Herunterfahrens, nach `OnShutdown` | Nein | Bereinigung |
| `OnShutdown` | Nach Bestätigung des Beendens | Nein | Anwendungsbereinigung |
| `ShouldQuit` | Wenn das Beenden angefordert wird | Ja (`false` zurückgeben) | Beenden bestätigen |
| `RegisterHook(WindowClosing)` | Wenn das Schließen des Fensters angefordert wird | Ja (`e.Cancel()`) | Schließen des Fensters verhindern |
| `OnWindowEvent` | Wenn das Ereignis eintritt | Nein | Auf Ereignisse reagieren |
| `OnApplicationEvent` | Wenn das Ereignis eintritt | Nein | Auf Ereignisse reagieren |

## Plattformunterschiede

### macOS

- Das **Anwendungsmenü** bleibt auch ohne geöffnete Fenster verfügbar
- **Cmd+Q** löst den Beendigungsvorgang aus (über `ShouldQuit`)
- Das **Dock-Symbol** bleibt sichtbar, sofern es nicht ausgeblendet wird
- Verwenden Sie `ApplicationShouldTerminateAfterLastWindowClosed`, um das Beendigungsverhalten zu steuern

### Windows

- Ohne Fenster gibt es **kein Anwendungsmenü**
- **Alt+F4** schließt das Fenster (kann mit `RegisterHook` verhindert werden)
- Der **Infobereich der Taskleiste** kann die Anwendung weiter ausführen lassen

### Linux

- Das **Verhalten variiert** je nach Desktop-Umgebung
- **Im Allgemeinen ähnlich wie unter Windows**

## Probleme mit dem Lebenszyklus beheben

### Problem: Anwendung lässt sich nicht beenden

**Ursachen:**

1. `ShouldQuit` gibt `false` zurück
2. `OnShutdown` dauert zu lange
3. Hintergrund-Goroutinen werden nicht beendet

**Lösung:**

```go
// 1. Check ShouldQuit logic
ShouldQuit: func() bool {
    log.Println("ShouldQuit called")
    return true
}

// 2. Keep OnShutdown fast
OnShutdown: func() {
    log.Println("OnShutdown started")
    // Fast cleanup only
    log.Println("OnShutdown finished")
}

// 3. Use context for background tasks
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    go func() {
        <-ctx.Done()
        log.Println("Context cancelled, stopping background work")
    }()
    return nil
}
```

### Problem: Dienst kann nicht gestartet werden

**Lösung:** Geben Sie aussagekräftige Fehler zurück:

```go
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    if err := s.init(); err != nil {
        return fmt.Errorf("failed to initialise: %w", err)
    }
    return nil
}
```

Der Fehler wird protokolliert und die Anwendung wird nicht gestartet.

## Bewährte Methoden

### Empfohlen

- **Verwenden Sie Dienste für die Lebenszyklusverwaltung** – Sie stellen geeignete Hooks für Start und Herunterfahren bereit
- **Sorgen Sie für schnelles Herunterfahren** – Streben Sie für sämtliche Bereinigungen insgesamt weniger als 1 Sekunde an
- **Verwenden Sie den Kontext für Abbrüche** – Beenden Sie Hintergrundaufgaben ordnungsgemäß
- **Behandeln Sie Fehler beim Start** – Geben Sie Fehler zurück, um den Start ordnungsgemäß abzubrechen
- **Protokollieren Sie Lebenszyklusereignisse** – Das erleichtert die Fehlersuche

### Nicht empfohlen

- **Führen Sie beim Dienststart keine blockierenden Operationen aus** – Halten Sie die Initialisierung kurz (unter 2 Sekunden)
- **Zeigen Sie beim Herunterfahren keine Dialogfelder an** – Die Anwendung wird beendet; die Benutzeroberfläche funktioniert möglicherweise nicht mehr
- **Ignorieren Sie den Kontext nicht** – Prüfen Sie in Goroutinen immer `ctx.Done()`
- **Vermeiden Sie Ressourcenlecks** – Implementieren Sie immer `ServiceShutdown`

## Nächste Schritte

**Dienste** – Erfahren Sie mehr über das Dienstsystem [Mehr erfahren →](/features/bindings/services/)

**Ereignissystem** – Verwenden Sie Ereignisse zur Kommunikation [Mehr erfahren →](/features/events/system/)

**Fensterverwaltung** – Erstellen und verwalten Sie Fenster [Mehr erfahren →](/features/windows/basics/)

---

**Fragen zum Lebenszyklus?** Fragen Sie auf [Discord](https://discord.gg/JDdSxwjhGf) nach oder sehen Sie sich die [Beispiele](https://github.com/wailsapp/wails/tree/master/v3/examples) an.
