---
title: "Benachrichtigungen"
description: "Native Systembenachrichtigungen mit Aktionsschaltflächen und Texteingabe anzeigen"
slug: "features/notifications/overview"
sourcePath: "features/notifications/overview.md"
---

## Einführung

Wails bietet ein umfassendes plattformübergreifendes Benachrichtigungssystem für Desktopanwendungen. Mit diesem Dienst können Sie native Systembenachrichtigungen mit folgenden Funktionen anzeigen:

- Einfache Benachrichtigungen mit Titel, Untertitel und Nachrichtentext
- Interaktive Benachrichtigungen mit Aktionsschaltflächen und Textantworten
- Wiederverwendbare [Benachrichtigungskategorien](#interaktive-benachrichtigungen) für Aktionen
- Benutzerdefinierte [Töne](#benutzerdefinierter-ton) (Standard, stumm oder benannt)
- [Anhänge](#anhnge) (Bilder auf allen Plattformen; Audio und Video unter macOS)
- [Gruppieren zusammengehöriger Benachrichtigungen](#zusammenfassen-und-gruppieren) anhand von `ThreadID`
- [Priorität](#unterbrechungsstufe) über `InterruptionLevel` (`passive` / `active` / `timeSensitive` / `critical`)
- [Geplante Zustellung](#geplante-zustellung) (nativ unter macOS; prozessinterner Timer unter Windows und Linux)
- [Aktualisieren einer bereits versendeten Benachrichtigung](#benachrichtigungen-aktualisieren) anhand ihrer ID

Für jedes neue optionale Feld greift ein geeignetes Rückfallverhalten, wenn eine Plattform es nicht unterstützen kann. Die Unterstützungsmatrix für die einzelnen Funktionen finden Sie unter [Plattformhinweise](#plattformspezifische-aspekte).

## Grundlegende Verwendung

### Dienst erstellen

Initialisieren Sie zuerst den Benachrichtigungsdienst:

```go
import "github.com/wailsapp/wails/v3/pkg/application"
import "github.com/wailsapp/wails/v3/pkg/services/notifications"

// Create a new notification service
notifier := notifications.New()

//Register the service with the application
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(notifier),
    },
})
```

## Autorisierung von Benachrichtigungen

Unter macOS erfordern Benachrichtigungen die Zustimmung des Benutzers. Fordern Sie die Autorisierung an und prüfen Sie sie:

```go
authorized, err := notifier.CheckNotificationAuthorization()
if err != nil {
    // Handle authorization error
}
if authorized {
    // Send notifications
} else {
    // Request authorization
    authorized, err = notifier.RequestNotificationAuthorization()
}
```

Unter Windows und Linux wird immer `true` zurückgegeben.

## Benachrichtigungstypen

### Einfache Benachrichtigungen

Senden Sie Benutzern eine einfache Benachrichtigung mit eindeutiger ID, Titel, optionalem Untertitel (macOS und Linux) und Nachrichtentext:

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID: "unique-id",
    Title: "New Calendar Invite",
    Subtitle: "From: Jane Doe", // Optional
    Body: "Tap to view the event",
})

```

### Interaktive Benachrichtigungen

Senden Sie eine Benachrichtigung mit Aktionsschaltflächen und Texteingabefeldern. Für diese Benachrichtigungen muss zuerst eine Benachrichtigungskategorie registriert werden:

```go
// Define a unique category id
categoryID := "unique-category-id"

// Define a category with actions
category := notifications.NotificationCategory{
    ID: categoryID,
    Actions: []notifications.NotificationAction{
        {
            ID:    "OPEN", 
            Title: "Open",
        },
        {
            ID:          "ARCHIVE", 
            Title:       "Archive", 
            Destructive: true,  /* macOS-specific */
        },
    },
    HasReplyField:    true,
    ReplyPlaceholder: "message...",
    ReplyButtonTitle: "Reply",
}

// Register the category
notifier.RegisterNotificationCategory(category)

// Send an interactive notification with the actions registered in the provided category
notifier.SendNotificationWithActions(notifications.NotificationOptions{
    ID:         "unique-id",
    Title:      "New Message",
    Subtitle:   "From: Jane Doe",
    Body:       "Are you able to make it?",
    CategoryID: categoryID,
})
```

## Reaktionen auf Benachrichtigungen

Verarbeiten Sie Benutzerinteraktionen mit Benachrichtigungen:

```go
notifier.OnNotificationResponse(func(result notifications.NotificationResult) {
    response := result.Response
    fmt.Printf("Notification %s was actioned with: %s\n", response.ID, response.ActionIdentifier)

    if response.ActionIdentifier == "TEXT_REPLY" {
        fmt.Printf("User replied: %s\n", response.UserText)
    }

    if data, ok := response.UserInfo["sender"].(string); ok {
        fmt.Printf("Original sender: %s\n", data)
    }

    // Emit an event to the frontend
    app.Event.Emit("notification", result.Response)
})
```

## Benachrichtigungen anpassen

### Benutzerdefinierte Metadaten

Einfache und interaktive Benachrichtigungen können benutzerdefinierte Daten enthalten:

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID: "unique-id",
    Title: "New Calendar Invite",
    Subtitle: "From: Jane Doe", // Optional
    Body: "Tap to view the event",
    Data: map[string]interface{}{
        "sender": "jane.doe@example.com",
        "timestamp": "2025-03-10T15:30:00Z",
    }
})
```

### Benutzerdefinierter Ton

Steuern Sie mit `Sound`, welcher Ton bei der Zustellung einer Benachrichtigung wiedergegeben wird. Wenn Sie das Feld auf `nil` belassen, wird der Standardton der Plattform wiedergegeben. Mit `Silent: true` unterdrücken Sie den Ton, mit `Name` geben Sie einen benannten oder gebündelten Ton wieder.

```go
// Silent
notifier.SendNotification(notifications.NotificationOptions{
    ID:    "silent-id",
    Title: "Background sync complete",
    Sound: &notifications.NotificationSound{Silent: true},
})

// Named sound
notifier.SendNotification(notifications.NotificationOptions{
    ID:    "named-id",
    Title: "New message",
    Body:  "Tap to read",
    Sound: &notifications.NotificationSound{Name: "Ping"},
})
```

So wird `Name` auf den einzelnen Plattformen aufgelöst:

- **macOS** — `Name` wird an `[UNNotificationSound soundNamed:]` übergeben. Die Audiodatei muss sich im `Library/Sounds`-Verzeichnis Ihres App-Bundles befinden.
- **Windows** — wenn `Name` bereits mit `ms-winsoundevent:` oder `ms-appx:` beginnt, wird der Wert unverändert verwendet. Andernfalls wird er zur Verwendung als Name eines integrierten Toast-Ereignisses in `ms-winsoundevent:` eingeschlossen (siehe Microsofts Dokumentation zum Toast-Schema `<audio>`).
- **Linux** — der Wert wird als freedesktop-Hinweis `sound-name` weitergeleitet. Ob der Ton wiedergegeben wird, hängt vom aktiven Benachrichtigungsdienst und Klangthema ab.

### Anhänge

Mit `Attachments` werden einer Benachrichtigung Mediendateien hinzugefügt. macOS unterstützt mehrere Anhänge beliebiger Medientypen. Windows und Linux berücksichtigen den ersten als Bild typisierten Anhang.

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID:    "image-id",
    Title: "Photo uploaded",
    Body:  "Tap to view",
    Attachments: []notifications.NotificationAttachment{
        {
            ID:   "preview",
            Path: "/absolute/path/to/image.png",
            // Type hint:
            //   macOS:   UTI like "public.png" / "public.audio" (often inferred)
            //   Windows: "hero" | "appLogoOverride" | "inline" (default "inline")
            //   Linux:   ignored
            Type: "inline",
        },
    },
})
```

`Path` muss ein absoluter Dateisystempfad sein. macOS akzeptiert zusätzlich `file://`-URLs.

#### Eine mit Ihrer App ausgelieferte Datei anhängen

Das Betriebssystem liest den Anhang bei der Zustellung der Benachrichtigung vom Datenträger. Daher muss `Path` auf eine tatsächlich vorhandene Datei auf dem Rechner des Endbenutzers verweisen. Für eine Ressource, die Sie mit der Anwendung bündeln, etwa ein mit `go:embed` eingebettetes Symbol oder Bild, können Sie keinen festen absoluten Pfad hartcodieren: Die Datei befindet sich innerhalb der Binärdatei und nicht an einem bekannten Speicherort auf dem Datenträger. Schreiben Sie sie beim Start einmalig in ein beschreibbares Verzeichnis und übergeben Sie diesen Pfad:

```go
import (
    _ "embed"
    "os"
    "path/filepath"
)

//go:embed assets/preview.png
var previewPNG []byte

// Materialise the embedded asset to a stable path the OS can read.
previewPath := filepath.Join(os.TempDir(), "myapp-preview.png")
if err := os.WriteFile(previewPath, previewPNG, 0o644); err != nil {
    // handle error
}

notifier.SendNotification(notifications.NotificationOptions{
    ID:    "image-id",
    Title: "Photo uploaded",
    Body:  "Tap to view",
    Attachments: []notifications.NotificationAttachment{
        {ID: "preview", Path: previewPath, Type: "inline"},
    },
})
```

Eine vom Benutzer bereitgestellte oder heruntergeladene Datei besitzt bereits einen tatsächlichen Pfad auf dem Datenträger. Sie können ihn daher ohne diesen Schritt direkt an `Path` übergeben. Die Übergabe von Anhangsdaten als Bytes im Arbeitsspeicher wird möglicherweise in einer zukünftigen Version ergänzt.

### Zusammenfassen und Gruppieren

`ThreadID` gruppiert zusammengehörige Benachrichtigungen, sodass das Betriebssystem sie in der Mitteilungszentrale bzw. im Info-Center zusammenfassen kann.

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID:       "msg-42",
    Title:    "Jane Doe",
    Body:     "Are you free for lunch?",
    ThreadID: "chat:jane.doe",
})
```

### Unterbrechungsstufe

`InterruptionLevel` steuert die Priorität der Benachrichtigung. Verwenden Sie eine der exportierten Konstanten:

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID:                "alert-id",
    Title:             "Server down",
    Body:              "Investigate immediately",
    InterruptionLevel: notifications.InterruptionLevelTimeSensitive,
})
```

| Konstante | Wert | Bedeutung |
| --- | --- | --- |
| `InterruptionLevelPassive` | `"passive"` | Unauffällige Zustellung; der Bildschirm wird nicht aktiviert und kein Standardton wiedergegeben |
| `InterruptionLevelActive` | `"active"` | Standardstufe |
| `InterruptionLevelTimeSensitive` | `"timeSensitive"` | Umgeht den Fokusmodus bzw. „Nicht stören“, sofern zulässig |
| `InterruptionLevelCritical` | `"critical"` | Umgeht Fokusmodus und Klingelton- bzw. Stummschaltungseinstellungen; unter macOS ist die Berechtigung für kritische Hinweise erforderlich (ohne sie erfolgt stillschweigend eine Herabstufung) |

Plattformspezifische Zuordnung:

- **macOS** — legt `UNNotificationContent.interruptionLevel` fest. `critical` erfordert macOS 12+ und die Berechtigung für kritische Hinweise.
- **Windows** — wird dem Toast-Attribut `<toast scenario="...">` zugeordnet.
- **Linux** — wird dem freedesktop-Hinweis `urgency` zugeordnet.

### Geplante Zustellung

`Schedule` verzögert die Zustellung. Legen Sie genau eine der Optionen `DelaySeconds` (Sekunden ab jetzt) oder `At` (Unix-Sekunden, UTC) fest.

```go
// Deliver in 60 seconds
notifier.SendNotification(notifications.NotificationOptions{
    ID:       "reminder-id",
    Title:    "Stand up and stretch",
    Schedule: &notifications.NotificationSchedule{DelaySeconds: 60},
})

// Deliver at an absolute time
notifier.SendNotification(notifications.NotificationOptions{
    ID:       "alarm-id",
    Title:    "Meeting in 5 minutes",
    Schedule: &notifications.NotificationSchedule{At: time.Now().Add(time.Hour).Unix()},
})
```

@note{type="caution" title="Persistenz"}
Unter **macOS** verwenden geplante Benachrichtigungen einen nativen Trigger und bleiben über App-Neustarts hinweg erhalten. Unter **Windows** und **Linux** greift die Planung auf einen prozessinternen `time.AfterFunc`-Timer zurück und geht **verloren, wenn die App vor der Zustellung beendet wird** — weder `wintoast` noch die freedesktop-Spezifikation stellen ein Primitiv für verzögerte Zustellung bereit.

@end

### Benachrichtigungen aktualisieren

`UpdateNotification` ersetzt eine bereits laufende Benachrichtigung mit demselben `ID`:

```go
notifier.UpdateNotification(notifications.NotificationOptions{
    ID:    "download-id",
    Title: "Download complete",
    Body:  "report.pdf is ready",
})
```

Plattformspezifisches Verhalten:

- **macOS** — `UNUserNotificationCenter` dedupliziert automatisch anhand der Kennung, sodass die vorhandene Benachrichtigung direkt aktualisiert wird.
- **Linux** — verwendet den D-Bus-Parameter `replaces_id`, um die vorherige Benachrichtigung zu ersetzen.
- **Windows** — stellt die Benachrichtigung derzeit erneut als neue Benachrichtigung zu. Ein echtes direktes Ersetzen erfordert vorgelagerte `wintoast`-Unterstützung für `tag`/`group`.

## Plattformspezifische Aspekte

@tabs
[macOS]
Unter macOS gelten für Benachrichtigungen folgende Punkte:

- Erfordern die Autorisierung durch den Benutzer
- Erfordern, dass die App paketiert und signiert ist (für die Verteilung notarisiert)
- Verwenden das systemübliche Erscheinungsbild für Benachrichtigungen
- Unterstützen `Subtitle`
- Unterstützen Texteingaben des Benutzers (Antworten)
- Unterstützen die Aktionsoption `Destructive`
- Unterstützen mehrere `Attachments` beliebiger Medientypen (Bilder, Audio, Video)
- Unterstützen `ThreadID` zur Gruppierung in der Mitteilungszentrale
- Unterstützen alle `InterruptionLevel`-Werte (`critical` erfordert die Berechtigung für kritische Hinweise)
- Unterstützen native geplante Zustellung, die über App-Neustarts hinweg erhalten bleibt
- Deduplizieren `UpdateNotification`-Aufrufe automatisch anhand von `ID`
- Berücksichtigen den dunklen und hellen Modus automatisch

[Windows]
Unter Windows gelten für Benachrichtigungen folgende Punkte:

- Verwenden über das `wintoast`-Backend die Windows-systemeigenen Toast-Stile
- Passen sich an die Windows-Designeinstellungen an
- Unterstützen Texteingaben des Benutzers (Antworten)
- Unterstützen Anzeigen mit hoher DPI-Dichte
- Unterstützen `Subtitle` nicht
- Unterstützen ein einzelnes Bild-`Attachment` mit dem Platzierungshinweis `hero`, `appLogoOverride` oder `inline` (Standard: `inline`)
- Unterstützen `ThreadID` zur Gruppierung im Info-Center
- Unterstützen `InterruptionLevel` über das Toast-Attribut `scenario`
- Unterstützen geplante Zustellung über einen prozessinternen Timer — **geplante Benachrichtigungen gehen verloren, wenn die App vor der Zustellung beendet wird**
- `UpdateNotification` stellt die Benachrichtigung derzeit erneut als neue Benachrichtigung zu (echtes direktes Ersetzen steht aus, bis die vorgelagerte `wintoast`-Unterstützung für `tag`/`group` verfügbar ist)

[Linux]
Unter Linux verwenden Benachrichtigungen die D-Bus-Schnittstelle `org.freedesktop.Notifications`. Damit Benachrichtigungen funktionieren, **muss ein kompatibler Benachrichtigungs-Daemon ausgeführt werden**.

@note{type="caution" title="Systemanforderung: Benachrichtigungs-Daemon"}
Ein freedesktop-kompatibler Benachrichtigungs-Daemon muss installiert sein und ausgeführt werden. Häufig verwendete Optionen:

- **dunst** — ressourcenschonend und umfassend konfigurierbar (`apt install dunst`/`dnf install dunst`)
- **mako** — nativ für Wayland (`apt install mako-notifier`)
- **GNOME Shell** — registriert die Schnittstelle unter GNOME 43+ automatisch. Unter Ubuntu 24.04 (GNOME Shell 46) wird die Schnittstelle beim Sitzungsstart möglicherweise nicht automatisch registriert. Installieren Sie ersatzweise `dunst`, wenn keine Benachrichtigungen angezeigt werden.
- **xfce4-notifyd** — in XFCE-Desktopumgebungen enthalten

Wenn kein Daemon ausgeführt wird, gibt `SendNotification` folgenden D-Bus-Fehler zurück: `The name org.freedesktop.Notifications was not provided by any .service files`. Behandeln Sie diesen Fehler in Ihrer App und weisen Sie den Benutzer darauf hin, einen Benachrichtigungs-Daemon zu installieren.

@end

Unter Linux gelten für Benachrichtigungen folgende Punkte:

- Übernehmen das Design der Desktopumgebung
- Werden gemäß den Regeln der Desktopumgebung positioniert
- Unterstützen `Subtitle` (wird bei Daemons, die es nicht separat darstellen, mit dem Textkörper verkettet)
- Unterstützen keine Texteingaben des Benutzers (nicht Bestandteil der freedesktop-Spezifikation)
- Unterstützen ein einzelnes Bild-`Attachment` über den Hinweis `image-path`
- Unterstützt `ThreadID` (wird, sofern unterstützt, vom Daemon verarbeitet)
- `Sound.Name` wird als Hinweis `sound-name` weitergeleitet; ob der Sound abgespielt wird, hängt vom aktiven Daemon und Soundthema ab
- Ordnet `InterruptionLevel` dem freedesktop-Hinweis `urgency` zu
- Unterstützt die geplante Zustellung über einen prozessinternen Timer – **geplante Benachrichtigungen gehen verloren, wenn die App vor der Zustellung beendet wird**
- `UpdateNotification` verwendet den D-Bus-Parameter `replaces_id`, um die vorherige Benachrichtigung direkt zu ersetzen

@end

## Bewährte Verfahren

1. Berechtigung prüfen und anfordern:
  - macOS erfordert die Berechtigung des Benutzers


2. Klare und prägnante Benachrichtigungen bereitstellen:
  - Aussagekräftige Titel, Untertitel, Texte und Aktionsbezeichnungen verwenden


3. Benachrichtigungsantworten angemessen verarbeiten:
  - Benachrichtigungsantworten auf Fehler prüfen
  - Rückmeldung zu Benutzeraktionen geben


4. Plattformkonventionen berücksichtigen:
  - Plattformspezifische Benachrichtigungsmuster befolgen
  - Systemeinstellungen berücksichtigen


5. Unter Linux die Abhängigkeit vom Daemon berücksichtigen:
  - Den von `SendNotification` zurückgegebenen Fehler prüfen – ein fehlender Daemon verursacht einen D-Bus-Fehler
  - In der Paketdokumentation oder README-Datei der App sollte darauf hingewiesen werden, dass ein freedesktop-Benachrichtigungs-Daemon erforderlich ist


## Beispiele

Dieses Beispiel ansehen:

- [Benachrichtigungen](https://github.com/wailsapp/wails/tree/master/v3/examples/notifications)

## API-Referenz

### Dienstverwaltung

| Methode | Beschreibung |
| --- | --- |
| `New()` | Erstellt einen neuen Benachrichtigungsdienst |

### Benachrichtigungsberechtigung

| Methode | Beschreibung |
| --- | --- |
| `RequestNotificationAuthorization()` | Fordert die Berechtigung zum Anzeigen von Benachrichtigungen an (macOS) |
| `CheckNotificationAuthorization()` | Prüft den aktuellen Status der Benachrichtigungsberechtigung (macOS) |

### Benachrichtigungen senden

| Methode | Beschreibung |
| --- | --- |
| `SendNotification(options NotificationOptions)` | Sendet eine einfache Benachrichtigung |
| `SendNotificationWithActions(options NotificationOptions)` | Sendet eine interaktive Benachrichtigung mit Aktionen |
| `UpdateNotification(options NotificationOptions)` | Aktualisiert eine noch nicht abgeschlossene Benachrichtigung anhand von `ID` (siehe [Benachrichtigungen aktualisieren](#benachrichtigungen-aktualisieren)) |

### Benachrichtigungskategorien

| Methode | Beschreibung |
| --- | --- |
| `RegisterNotificationCategory(category NotificationCategory)` | Registriert eine wiederverwendbare Benachrichtigungskategorie |
| `RemoveNotificationCategory(categoryID string)` | Entfernt eine zuvor registrierte Kategorie |

### Benachrichtigungen verwalten

| Methode | Beschreibung |
| --- | --- |
| `RemoveAllPendingNotifications()` | Entfernt alle ausstehenden Benachrichtigungen (nur macOS und Linux) |
| `RemovePendingNotification(identifier string)` | Entfernt eine bestimmte ausstehende Benachrichtigung (nur macOS und Linux) |
| `RemoveAllDeliveredNotifications()` | Entfernt alle zugestellten Benachrichtigungen (nur macOS und Linux) |
| `RemoveDeliveredNotification(identifier string)` | Entfernt eine bestimmte zugestellte Benachrichtigung (nur macOS und Linux) |
| `RemoveNotification(identifier string)` | Entfernt eine Benachrichtigung (Linux-spezifisch) |

### Ereignisverarbeitung

| Methode | Beschreibung |
| --- | --- |
| `OnNotificationResponse(callback func(result NotificationResult))` | Registriert eine Callback-Funktion für Benachrichtigungsantworten |

### Strukturen und Typen

#### NotificationOptions

```go
type NotificationOptions struct {
    ID         string                 // Unique identifier for the notification (required)
    Title      string                 // Main notification title (required)
    Subtitle   string                 // Subtitle text (macOS and Linux only)
    Body       string                 // Main notification content
    CategoryID string                 // Category identifier for interactive notifications
    Data       map[string]interface{} // Custom data to associate with the notification

    // Sound controls the sound played on delivery. nil = platform default.
    Sound *NotificationSound

    // Attachments are media files shown alongside the notification.
    // macOS supports multiple attachments of any media type;
    // Windows and Linux honour the first image-typed attachment.
    Attachments []NotificationAttachment

    // ThreadID groups related notifications together in
    // Notification Center / Action Center / the Linux notification daemon.
    ThreadID string

    // InterruptionLevel controls priority. One of "passive",
    // "active" (default), "timeSensitive", "critical".
    InterruptionLevel string

    // Schedule defers delivery. macOS uses a native trigger that survives
    // restarts; Windows and Linux use an in-process timer that does NOT.
    Schedule *NotificationSchedule
}
```

#### NotificationSound

```go
type NotificationSound struct {
    Silent bool   // If true, no sound is played
    Name   string // Named/bundled sound (see "Custom Sound" above)
}
```

#### NotificationAttachment

```go
type NotificationAttachment struct {
    ID   string // Optional identifier
    Path string // Absolute filesystem path (macOS also accepts file:// URLs)
    // Type is an optional placement/UTI hint:
    //   macOS:   UTI like "public.png" / "public.audio" (often inferred)
    //   Windows: "hero" | "appLogoOverride" | "inline" (default "inline")
    //   Linux:   ignored (always image-path hint)
    Type string
}
```

#### NotificationSchedule

```go
// Exactly one of DelaySeconds or At must be set. At is Unix seconds (UTC).
type NotificationSchedule struct {
    DelaySeconds int   // Seconds from now until delivery
    At           int64 // Absolute Unix timestamp (seconds, UTC)
}
```

#### InterruptionLevel-Konstanten

```go
const (
    InterruptionLevelPassive       = "passive"
    InterruptionLevelActive        = "active" // default
    InterruptionLevelTimeSensitive = "timeSensitive"
    InterruptionLevelCritical      = "critical"
)
```

#### NotificationCategory

```go
type NotificationCategory struct {
    ID               string                // Unique identifier for the category
    Actions          []NotificationAction  // Button actions for the notification
    HasReplyField    bool                  // Whether to include a text input field
    ReplyPlaceholder string                // Placeholder text for the input field
    ReplyButtonTitle string                // Text for the reply button
}
```

#### NotificationAction

```go
type NotificationAction struct {
    ID          string  // Unique identifier for the action
    Title       string  // Button text
    Destructive bool    // Whether the action is destructive (macOS-specific)
}
```

#### NotificationResponse

```go
type NotificationResponse struct {
    ID               string                  // Notification identifier
    ActionIdentifier string                  // Action that was triggered
    CategoryID       string                  // Category of the notification
    Title            string                  // Title of the notification
    Subtitle         string                  // Subtitle of the notification
    Body             string                  // Body text of the notification
    UserText         string                  // Text entered by the user
    UserInfo         map[string]interface{}  // Custom data from the notification
}
```

#### NotificationResult

```go
type NotificationResult struct {
    Response NotificationResponse  // Response data
    Error    error                 // Any error that occurred
}
```
