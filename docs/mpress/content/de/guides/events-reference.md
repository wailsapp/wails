---
title: "Leitfaden zu Ereignissen"
description: "Ein praxisnaher Leitfaden zur Verwendung von Ereignissen in Wails v3 für die Anwendungskommunikation und die Verwaltung des Lebenszyklus"
slug: "guides/events-reference"
sourcePath: "guides/events-reference.md"
---

**HINWEIS: Dieser Leitfaden ist noch in Arbeit**

## Leitfaden zu Ereignissen

Ereignisse sind das Herzstück der Kommunikation in Wails-Anwendungen. Über sie können verschiedene Teile Ihrer Anwendung miteinander kommunizieren, ohne eng gekoppelt zu sein. Dieser Leitfaden vermittelt Ihnen alles, was Sie wissen müssen, um Ereignisse in Ihrer Wails-Anwendung effektiv einzusetzen.

## Wails-Ereignisse verstehen

Stellen Sie sich Ereignisse als Nachrichten vor, die in Ihrer gesamten Anwendung verbreitet werden. Jeder Teil Ihrer Anwendung kann auf diese Nachrichten reagieren. Dies ist besonders nützlich für:

- **Auf Fensteränderungen reagieren**: Erkennen, wenn Ihr Fenster minimiert, maximiert oder verschoben wird
- **Systemereignisse verarbeiten**: Auf Änderungen des Farbschemas oder Energieverwaltungsereignisse reagieren
- **Benutzerdefinierte Anwendungslogik**: Eigene Ereignisse für Funktionen wie Datenaktualisierungen oder Benutzeraktionen erstellen
- **Komponentenübergreifende Kommunikation**: Verschiedene Teile Ihrer Anwendung ohne direkte Abhängigkeiten kommunizieren lassen

## Namenskonvention für Ereignisse

Alle Wails-Ereignisse folgen einem Namensraumschema, das ihren Ursprung eindeutig kennzeichnet:

- `common:` – Plattformübergreifende Ereignisse für Windows, macOS und Linux
- `windows:` – Windows-spezifische Ereignisse
- `mac:` – macOS-spezifische Ereignisse\
- `linux:` – Linux-spezifische Ereignisse

Beispiele:

- `common:WindowFocus` – Das Fenster hat den Fokus erhalten (funktioniert überall)
- `windows:APMSuspend` – Das System wird in den Standbymodus versetzt (nur Windows)
- `mac:ApplicationDidBecomeActive` – Die Anwendung wurde aktiv (nur macOS)

## Erste Schritte mit Ereignissen

### Auf Ereignisse reagieren (Frontend)

Der häufigste Anwendungsfall ist das Reagieren auf Ereignisse in Ihrem Frontend-Code:

```javascript
import { Events } from '@wailsio/runtime';

// Listen for when the window gains focus
Events.On('common:WindowFocus', () => {
    console.log('Window is now focused!');
    // Maybe refresh some data or resume animations
});

// Listen for theme changes
Events.On('common:ThemeChanged', (event) => {
    console.log('Theme changed:', event.data);
    // Update your app's theme accordingly
});

// Listen for custom events from your Go backend
Events.On('my-app:data-updated', (event) => {
    console.log('Data updated:', event.data);
    // Update your UI with the new data
});
```

### Ereignisse auslösen (Backend)

In Ihrem Go-Code können Sie Ereignisse auslösen, auf die Ihr Frontend reagieren kann:

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "time"
)

func (s *Service) UpdateData() {
    // Do some data processing...

    // Notify the frontend
    app := application.Get()
    app.Event.Emit("my-app:data-updated",
        map[string]interface{}{
            "timestamp": time.Now(),
            "count": 42,
        },
	)
}
```

### Ereignisse auslösen (Frontend)

Sie können auch im Frontend Ereignisse auslösen, auf die Ihr Go-Code reagieren kann, auch wenn dies seltener verwendet wird:

```javascript
import { Events } from '@wailsio/runtime';

// Event without data
Events.Emit('myapp:close-window')

// Event with data
Events.Emit('myapp:disconnect-requested', 'id-123')
```

Wenn Sie TypeScript in Ihrem Frontend verwenden und in Ihrem Go-Code [typisierte Ereignisse registrieren](#typisierte-ereignisse-mit-typsicherheit), erhalten Sie eine automatische Vervollständigung und Prüfung von Ereignisnamen sowie eine Datentypprüfung.

### Ereignis-Listener entfernen

Entfernen Sie Ereignis-Listener immer, sobald sie nicht mehr benötigt werden:

```javascript
import { Events } from '@wailsio/runtime';

// Store the handler reference
const focusHandler = () => {
    console.log('Window focused');
};

// Add the listener — capture the unsubscribe function it returns
const unsubscribe = Events.On('common:WindowFocus', focusHandler);

// Later, remove this specific listener via the returned unsubscribe
unsubscribe();

// Or remove ALL listeners for one (or more) event names — Events.Off takes only event-name strings
Events.Off('common:WindowFocus');
// Events.Off('common:WindowFocus', 'common:WindowLostFocus'); // variadic
// Events.OffAll(); // remove every listener for every event (no args)
```

## Häufige Anwendungsfälle

### 1. Bei Änderung des Fensterfokus pausieren/fortsetzen

Viele Anwendungen müssen bestimmte Aktivitäten pausieren, wenn das Fenster den Fokus verliert:

```javascript
import { Events } from '@wailsio/runtime';

let animationRunning = true;

Events.On('common:WindowLostFocus', () => {
    animationRunning = false;
    pauseBackgroundTasks();
});

Events.On('common:WindowFocus', () => {
    animationRunning = true;
    resumeBackgroundTasks();
});
```

### 2. Auf Änderungen des Farbschemas reagieren

Halten Sie Ihre Anwendung mit dem Farbschema des Systems synchron:

```javascript
import { Events } from '@wailsio/runtime';

Events.On('common:ThemeChanged', (event) => {
    const isDarkMode = event.data.isDarkMode;

    if (isDarkMode) {
        document.body.classList.add('dark-theme');
        document.body.classList.remove('light-theme');
    } else {
        document.body.classList.add('light-theme');
        document.body.classList.remove('dark-theme');
    }
});
```

### 3. Abgelegte Dateien verarbeiten

Ermöglichen Sie Ihrer Anwendung, per Drag-and-drop abgelegte Dateien anzunehmen:

```javascript
import { Events } from '@wailsio/runtime';

Events.On('common:WindowFilesDropped', (event) => {
    const files = event.data.files;

    files.forEach(file => {
        console.log('File dropped:', file);
        // Process the dropped files
        handleFileUpload(file);
    });
});
```

### 4. Fensterlebenszyklus verwalten

Reagieren Sie auf Änderungen des Fensterzustands:

```javascript
import { Events } from '@wailsio/runtime';

Events.On('common:WindowClosing', () => {
    // Save user data before closing
    saveApplicationState();

    // You could also prevent closing by returning false
    // from a registered window close handler
});

Events.On('common:WindowMaximise', () => {
    // Adjust UI for maximized view
    adjustLayoutForMaximized();
});

Events.On('common:WindowRestore', () => {
    // Return UI to normal state
    adjustLayoutForNormal();
});
```

### 5. Plattformspezifische Funktionen

Verarbeiten Sie bei Bedarf plattformspezifische Ereignisse:

```javascript
import { Events } from '@wailsio/runtime';

// Windows-specific power management
Events.On('windows:APMSuspend', () => {
    console.log('System is going to sleep');
    saveState();
});

Events.On('windows:APMResumeSuspend', () => {
    console.log('System woke up');
    refreshData();
});

// macOS-specific app lifecycle
Events.On('mac:ApplicationWillTerminate', () => {
    console.log('App is about to quit');
    performCleanup();
});
```

## Benutzerdefinierte Ereignisse erstellen

Sie können eigene Ereignisse für anwendungsspezifische Anforderungen erstellen.

### Backend (Go)

```go
// Emit a custom event when data changes

func (s *Service) ProcessUserData(userData UserData) error {
    // Process the data...

    app := application.Get()
    // Notify all listeners
    app.Event.Emit("user:data-processed",
        map[string]interface{}{
            "userId": userData.ID,
            "status": "completed",
            "timestamp": time.Now(),
        },
    )
    return nil
}

// Emit periodic updates
func (s *Service) StartMonitoring() {
    app := application.Get()
    ticker := time.NewTicker(5 * time.Second)
    go func() {
        for range ticker.C {
            stats := s.collectStats()
            app.Event.Emit("monitor:stats-updated", stats)
        }
    }()
}
```

### Frontend (JavaScript)

```javascript
import { Events } from '@wailsio/runtime';

// Listen for your custom events
Events.On('user:data-processed', (event) => {
    const { userId, status, timestamp } = event.data;

    showNotification(`User ${userId} processing ${status}`);
    updateUIWithNewData();
});

Events.On('monitor:stats-updated', (event) => {
    updateDashboard(event.data);
});
```

## Typisierte Ereignisse mit Typsicherheit

Wails v3 unterstützt durch die Registrierung von Ereignissen und die automatische Generierung von Bindings typisierte Ereignisse mit vollständiger TypeScript-Typsicherheit.

### Benutzerdefinierte Ereignisse registrieren

Rufen Sie während der Initialisierung `application.RegisterEvent` auf, um benutzerdefinierte Ereignisnamen mit ihren Datentypen zu registrieren:

```go
package main

import "github.com/wailsapp/wails/v3/pkg/application"

type UserLoginData struct {
    UserID   string
    Username string
    LoginTime string
}

type MonitorStats struct {
    CPUUsage    float64
    MemoryUsage float64
}

func init() {
    // Register events with their data types
    application.RegisterEvent[UserLoginData]("user:login")
    application.RegisterEvent[MonitorStats]("monitor:stats")

    // Register events without data (void events)
    application.RegisterEvent[application.Void]("app:ready")
}
```

@note{type="caution"}
`RegisterEvent` ist für den Aufruf während der Initialisierung vorgesehen und löst in folgenden Fällen eine Panic aus:

- Die Argumente sind ungültig
- Derselbe Ereignisname wird zweimal mit unterschiedlichen Datentypen registriert

@end

@note{type="info"}
Dasselbe Ereignis kann mehrmals sicher registriert werden, solange der Datentyp immer identisch ist. So lässt sich sicherstellen, dass ein Ereignis registriert wird, sobald eines von mehreren Paketen geladen wird.

@end

### Vorteile der Ereignisregistrierung

Nach der Registrierung werden an `Event.Emit` übergebene Datenargumente anhand des angegebenen Typs geprüft. Bei einer Abweichung gilt:

- Ein Fehler wird ausgegeben und protokolliert (oder an den registrierten Fehler-Handler übergeben)
- Das betreffende Ereignis wird nicht weitergeleitet
- Dadurch ist sichergestellt, dass das Datenfeld registrierter Ereignisse stets dem deklarierten Typ zugewiesen werden kann

### Strikter Modus

Verwenden Sie das Build-Tag `strictevents`, um während der Entwicklung Warnungen für nicht registrierte Ereignisse zu aktivieren:

```bash
go build -tags strictevents
```

Bei aktiviertem striktem Modus gibt die Runtime pro nicht registriertem Ereignisnamen höchstens eine Warnung aus, damit die Protokolle nicht überflutet werden.

### TypeScript-Bindings generieren

Der Binding-Generator erzeugt TypeScript-Definitionen und Glue-Code zur transparenten Unterstützung typisierter Ereignisse im Frontend.

#### 1. Vite-Plug-in einrichten

In Ihrer `vite.config.ts`:

```typescript
import { defineConfig } from 'vite'
import wails from '@wailsio/runtime/plugins/vite'

export default defineConfig({
  plugins: [wails()],
})
```

#### 2. Bindings generieren

Führen Sie den Binding-Generator aus:

```bash
wails3 generate bindings
```

Dadurch werden in Ihrem Frontend-Verzeichnis TypeScript-Dateien mit typisierten Definitionen der Ereignisdaten sowie ein Modul `bindings/events` erstellt. Dieses Modul exportiert für jeden registrierten Ereignisnamen, der eine vom Generator erkennbare Kompilierzeitkonstante ist, jeweils eine Konstante (beispielsweise wird `user:login` zu `UserLogin`).

#### 3. Typisierte Ereignisse im Frontend verwenden

```typescript
import { Events } from '@wailsio/runtime'
import { UserLogin, MonitorStats } from './bindings/events'

// Type-safe event emission with autocomplete
Events.Emit(UserLogin, {
    UserID: "123",
    Username: "john_doe",
    LoginTime: new Date().toISOString()
})

// Type-safe event listening
Events.On(UserLogin, (event) => {
    // event.data is typed as UserLoginData
    console.log(`User ${event.data.Username} logged in`)
})

Events.On(MonitorStats, (event) => {
    // event.data is typed as MonitorStats
    updateDashboard({
        cpu: event.data.CPUUsage,
        memory: event.data.MemoryUsage
    })
})
```

Die typisierten Ereignisse bieten:

- **Automatische Vervollständigung** für Ereignisnamen
- **Typprüfung** für Ereignisdaten
- **Kompilierzeitfehler** bei nicht übereinstimmenden Datentypen
- **IntelliSense**-Dokumentation

## Ereignisreferenz

### Allgemeine Ereignisse (plattformübergreifend)

Diese Ereignisse funktionieren auf allen Plattformen:

| Ereignis | Beschreibung | Verwendungszweck |
| --- | --- | --- |
| `common:ApplicationStarted` | Die Anwendung wurde vollständig gestartet | App initialisieren, gespeicherten Zustand laden |
| `common:WindowRuntimeReady` | Wails-Laufzeit ist bereit | Wails-API-Aufrufe starten |
| `common:ThemeChanged` | Systemdesign wurde geändert | Erscheinungsbild der App aktualisieren |
| `common:SystemWillSleep` | System wird gleich in den Ruhezustand versetzt | Zustand schreiben, Sockets schließen |
| `common:SystemDidWake` | System wurde aus dem Ruhezustand fortgesetzt | Verbindung wiederherstellen, veraltete Daten aktualisieren |
| `common:WindowFocus` | Fenster hat den Fokus erhalten | Aktivitäten fortsetzen, Daten aktualisieren |
| `common:WindowLostFocus` | Fenster hat den Fokus verloren | Aktivitäten pausieren, Zustand speichern |
| `common:WindowMinimise` | Fenster wurde minimiert | Rendering pausieren, Ressourcennutzung reduzieren |
| `common:WindowMaximise` | Fenster wurde maximiert | Layout für den Vollbildmodus anpassen |
| `common:WindowRestore` | Fenster wurde aus dem minimierten oder maximierten Zustand wiederhergestellt | Zum normalen Layout zurückkehren |
| `common:WindowClosing` | Fenster wird gleich geschlossen | Daten speichern, Ressourcen bereinigen |
| `common:WindowFilesDropped` | Dateien wurden auf dem Fenster abgelegt | Dateiimporte verarbeiten |
| `common:WindowDidResize` | Fenstergröße wurde geändert | Layout anpassen, Diagramme neu rendern |
| `common:WindowDidMove` | Fenster wurde verschoben | Positionsabhängige Funktionen aktualisieren |

### Plattformspezifische Ereignisse

#### Windows-Ereignisse

Wichtige Ereignisse für Windows-Anwendungen:

| Ereignis | Beschreibung | Anwendungsfall |
| --- | --- | --- |
| `windows:SystemThemeChanged` | Windows-Design wurde geändert | App-Farben aktualisieren |
| `windows:APMSuspend` | System wird in den Ruhezustand versetzt | Zustand speichern, Vorgänge pausieren |
| `windows:APMResumeAutomatic` | System wurde fortgesetzt (wird beim Fortsetzen immer ausgelöst) | Zustand wiederherstellen, Daten aktualisieren |
| `windows:APMResumeSuspend` | System wurde durch eine Benutzereingabe fortgesetzt (nach `APMResumeAutomatic`) | Vom Benutzer initiiertes Aufwecken erkennen |
| `windows:APMPowerStatusChange` | Energiestatus wurde geändert | Leistungseinstellungen anpassen |

#### macOS-Ereignisse

Wichtige macOS-Anwendungsereignisse:

| Ereignis | Beschreibung | Anwendungsfall |
| --- | --- | --- |
| `mac:ApplicationDidBecomeActive` | App wurde aktiv | Vorgänge fortsetzen |
| `mac:ApplicationDidResignActive` | App wurde inaktiv | Vorgänge pausieren |
| `mac:ApplicationWillTerminate` | App wird beendet | Abschließende Bereinigung |
| `mac:ApplicationWillSleep` | System wird gleich in den Ruhezustand versetzt | Zustand speichern, Sockets schließen |
| `mac:ApplicationDidWake` | System wurde fortgesetzt | Verbindung wiederherstellen, aktualisieren |
| `mac:ApplicationScreensDidSleep` | Bildschirme wurden in den Ruhezustand versetzt | Rendering pausieren (unabhängig vom Systemruhezustand) |
| `mac:ApplicationScreensDidWake` | Bildschirme wurden aufgeweckt | Rendering fortsetzen |
| `mac:WindowDidEnterFullScreen` | Vollbildmodus wurde aktiviert | Benutzeroberfläche für den Vollbildmodus anpassen |
| `mac:WindowDidExitFullScreen` | Vollbildmodus wurde beendet | Normale Benutzeroberfläche wiederherstellen |

#### Linux-Ereignisse

Zentrale Linux-Fensterereignisse:

| Ereignis | Beschreibung | Anwendungsfall |
| --- | --- | --- |
| `linux:SystemThemeChanged` | Desktop-Design wurde geändert | App-Design aktualisieren |
| `linux:SystemWillSleep` | System wird gleich in den Ruhezustand versetzt (logind) | Zustand speichern |
| `linux:SystemDidWake` | System wurde fortgesetzt (logind) | Verbindung wiederherstellen, aktualisieren |
| `linux:WindowFocusIn` | Fenster hat den Fokus erhalten | Aktivitäten fortsetzen |
| `linux:WindowFocusOut` | Fenster hat den Fokus verloren | Aktivitäten pausieren |
| `linux:WindowLoadStarted` | WebView hat den Ladevorgang gestartet | Ladeanzeige anzeigen |
| `linux:WindowLoadRedirected` | WebView wurde umgeleitet | Navigationsweiterleitungen verfolgen |
| `linux:WindowLoadCommitted` | WebView hat den Ladevorgang bestätigt | Inhalte werden empfangen |
| `linux:WindowLoadFinished` | WebView hat den Ladevorgang abgeschlossen | Ladeanzeige ausblenden, JS/CSS einschleusen |

## Bewährte Vorgehensweisen

### 1. Ereignisnamensräume verwenden

Verwenden Sie beim Erstellen benutzerdefinierter Ereignisse Namensräume, um Konflikte zu vermeiden:

```javascript
import { Events } from '@wailsio/runtime';

// Good - namespaced events
Events.Emit('myapp:user:login');
Events.Emit('myapp:data:updated');
Events.Emit('myapp:network:connected');

// Avoid - generic names that might conflict
Events.Emit('login');
Events.Emit('update');
```

### 2. Listener bereinigen

Entfernen Sie Ereignis-Listener immer, wenn Komponenten ausgehängt werden:

```javascript
import { Events } from '@wailsio/runtime';

// React example
useEffect(() => {
    const handler = (event) => {
        // Handle event
    };

    const off = Events.On('common:WindowDidResize', handler);

    // Cleanup — call the unsubscribe returned by Events.On
    return () => {
        off();
    };
}, []);
```

### 3. Plattformunterschiede berücksichtigen

Prüfen Sie bei plattformspezifischen Ereignissen deren Verfügbarkeit auf der jeweiligen Plattform:

```javascript
import { Events } from '@wailsio/runtime';

// Platform-specific events can be registered unconditionally;
// they will simply never fire on unsupported platforms.
Events.On('windows:APMSuspend', handleSuspend);
Events.On('mac:ApplicationWillTerminate', handleTerminate);
```

### 4. Ereignisse nicht übermäßig verwenden

Ereignisse sind zwar leistungsfähig, sollten aber nicht für alles verwendet werden:

- ✅ Ereignisse verwenden für: Systembenachrichtigungen, Lebenszyklusänderungen, übertragene Aktualisierungen
- ❌ Ereignisse vermeiden für: direkte Funktionsrückgaben, Aktualisierungen einzelner Komponenten, synchrone Vorgänge

## Ereignisse debuggen

So debuggen Sie Probleme mit Ereignissen:

```javascript
import { Events } from '@wailsio/runtime';

// Log all events (development only)
if (isDevelopment) {
    const originalOn = Events.On;
    Events.On = function(eventName, handler) {
        console.log(`[Event Registered] ${eventName}`);
        return originalOn.call(this, eventName, function(event) {
            console.log(`[Event Fired] ${eventName}`, event);
            return handler(event);
        });
    };
}
```

## Maßgebliche Quelle

Die vollständige Liste der verfügbaren Ereignisse finden Sie im Wails-Quellcode:

- Frontend-Ereignisse: [`v3/internal/runtime/desktop/@wailsio/runtime/src/event_types.ts`](https://github.com/wailsapp/wails/blob/main/v3/internal/runtime/desktop/@wailsio/runtime/src/event_types.ts)
- Backend-Ereignisse: [`v3/pkg/events/events.go`](https://github.com/wailsapp/wails/blob/main/v3/pkg/events/events.go)

Die aktuellsten Ereignisnamen und Angaben zu ihrer Verfügbarkeit finden Sie stets in diesen Dateien.

## Zusammenfassung

Ereignisse bieten in Wails eine leistungsfähige, entkoppelte Möglichkeit, die Kommunikation in Ihrer Anwendung zu handhaben. Wenn Sie die Muster und bewährten Vorgehensweisen in diesem Leitfaden befolgen, können Sie reaktionsschnelle, plattformgerechte Anwendungen entwickeln, die flüssig auf Systemänderungen und Benutzerinteraktionen reagieren.

Beachten Sie: Beginnen Sie für plattformübergreifende Kompatibilität mit allgemeinen Ereignissen, ergänzen Sie bei Bedarf plattformspezifische Ereignisse und bereinigen Sie stets Ihre Ereignis-Listener, um Speicherlecks zu vermeiden.
