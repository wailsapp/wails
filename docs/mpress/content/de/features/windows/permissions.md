---
title: "Berechtigungen"
description: "Anfragen von Webinhalten auf Kamera, Mikrofon, Standort und andere Funktionen steuern"
slug: "features/windows/permissions"
sourcePath: "features/windows/permissions.md"
---

Webinhalte, die `navigator.mediaDevices.getUserMedia()`, die Geolocation API oder die Notifications API aufrufen, benötigen von der Hostanwendung eine Genehmigung oder Ablehnung dieser Anfragen. Wails stellt dafür plattformübergreifend eine `Permissions`-Map auf `WebviewWindowOptions` bereit, mit der sich dies deklarativ steuern lässt – plattformspezifischer Code ist nicht erforderlich.

## Schnellstart

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "My App",
    Permissions: map[application.PermissionType]application.Permission{
        application.PermissionMicrophone: application.PermissionAllow,
        application.PermissionCamera:     application.PermissionAllow,
    },
})
```

Anfragen der Webinhalte dieses Fensters auf Kamera und Mikrofon werden ohne Browserabfrage genehmigt.

## Berechtigungstypen

`PermissionType` (uint8) bezeichnet eine Funktion, die Webinhalte anfordern können.

| Konstante | Funktion |
| --- | --- |
| `PermissionMicrophone` | `getUserMedia({audio: true})` |
| `PermissionCamera` | `getUserMedia({video: true})` |
| `PermissionGeolocation` | `navigator.geolocation` |
| `PermissionNotifications` | `Notification.requestPermission()` |
| `PermissionClipboardRead` | `navigator.clipboard.readText()` |

## Berechtigungswerte

`Permission` (uint8) ist die Richtlinie, die auf einen bestimmten Typ angewendet wird.

| Konstante | Wert | Bedeutung |
| --- | --- | --- |
| `PermissionDefault` | 0 | Native Verarbeitung der Plattform verwenden (siehe unten) |
| `PermissionAllow` | 1 | Ohne Nachfrage genehmigen |
| `PermissionDeny` | 2 | Ohne Nachfrage ablehnen |

`PermissionDefault` ist der Nullwert. Nicht gesetzte Map-Einträge verhalten sich daher gemäß der Standardrichtlinie.

## Plattformverhalten

Jede Plattform verarbeitet `PermissionDefault` anders, da sich das native Verhalten der zugrunde liegenden Webviews unterscheidet.

### Linux (WebKitGTK)

WebKitGTK verfügt über **keine native Berechtigungsabfrage**. Ohne registrierten Handler lehnt es jede Anfrage stillschweigend ab. Deshalb gab `getUserMedia` vor Einführung dieser Funktion immer `NotAllowedError` zurück.

Derzeit verarbeitet Wails unter Linux Anfragen auf **Kamera und Mikrofon**. Standortzugriff, Benachrichtigungen und das Lesen der Zwischenablage sind noch nicht angebunden und bleiben unabhängig von der festgelegten Richtlinie gesperrt.

| Richtlinie | Kamera/Mikrofon | Standort, Benachrichtigungen, Zwischenablage |
| --- | --- | --- |
| `PermissionDefault` | **Zugelassen** (stellt getUserMedia wieder her) | Immer abgelehnt |
| `PermissionAllow` | Zugelassen | Immer abgelehnt (noch nicht implementiert) |
| `PermissionDeny` | Abgelehnt | Immer abgelehnt |

### Windows (WebView2)

WebView2 verfügt über eine native Berechtigungsabfrage und eine Berechtigungs-API für die einzelnen Typen. Alle fünf Funktionstypen werden vollständig unterstützt.

| Richtlinie | Verhalten |
| --- | --- |
| `PermissionDefault` | WebView2 zeigt die native Berechtigungsabfrage des Betriebssystems bzw. Browsers an |
| `PermissionAllow` | Stillschweigend genehmigt |
| `PermissionDeny` | Stillschweigend abgelehnt |

**Wichtig:** Vor Einführung dieser Funktion rief Wails `SetGlobalPermission(Allow)` bedingungslos auf und genehmigte damit stillschweigend alle Funktionen. Sobald `Permissions` einen Eintrag enthält, wird diese pauschale Genehmigung nun **nicht** festgelegt. Nicht gesetzte Funktionen werden nicht automatisch genehmigt, sondern an die native Berechtigungsabfrage von WebView2 weitergeleitet.

Wenn du `Permissions` unter Windows überhaupt konfigurierst, wird daher für jede nicht ausdrücklich aufgeführte Funktion eine Abfrage angezeigt, statt sie stillschweigend zuzulassen. Lege die benötigten Funktionen ausdrücklich fest.

### macOS (TCC)

macOS verwaltet den Zugriff auf Kamera, Mikrofon, Standort und Benachrichtigungen über sein systemweites Datenschutz-Framework. Wenn Webinhalte erstmals eine Funktion anfordern, erscheint automatisch die Abfrage des Betriebssystems. Die Auswahl des Benutzers wird in den Systemeinstellungen unter „Datenschutz & Sicherheit“ anwendungsspezifisch gespeichert.

Dies funktioniert auch ohne `Permissions`-Konfiguration ordnungsgemäß. Die Map wird unter macOS **derzeit ignoriert**: Unabhängig von den festgelegten Werten laufen alle Anfragen über TCC. In der Praxis bedeutet dies, dass `PermissionDeny` unter macOS wirkungslos ist: Du kannst ein Webview nicht daran hindern, eine Funktion zu verwenden, die TCC bereits auf Systemebene genehmigt hat.

Stelle sicher, dass deine `Info.plist` die entsprechenden Schlüssel für Verwendungsbeschreibungen enthält:

```xml
<key>NSMicrophoneUsageDescription</key>
<string>Used for voice input</string>
<key>NSCameraUsageDescription</key>
<string>Used for video calls</string>
```

## Gängige Muster

### Anwendung zur Medienaufnahme

Kamera und Mikrofon auf allen Plattformen genehmigen:

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionAllow,
},
```

Unter **Linux** werden damit beide Geräte ausdrücklich zugelassen; andere Funktionen bleiben verweigert. Unter **Windows** werden damit beide zugelassen; für jede andere nicht aufgeführte Funktion wird eine native Abfrage angezeigt. Unter **macOS** hat dies keine Wirkung; TCC übernimmt alles.

### Medienaufnahme unter Linux verweigern

Linux lässt Kamera und Mikrofon standardmäßig zu. So deaktivieren Sie dies:

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionDeny,
    application.PermissionCamera:     application.PermissionDeny,
},
```

### Alle Funktionen unter Windows zulassen

So gewähren Sie alle Funktionen ohne Abfrage:

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone:    application.PermissionAllow,
    application.PermissionCamera:        application.PermissionAllow,
    application.PermissionGeolocation:   application.PermissionAllow,
    application.PermissionNotifications: application.PermissionAllow,
    application.PermissionClipboardRead: application.PermissionAllow,
},
```

### Richtlinien pro Fenster

Für verschiedene Fenster können unterschiedliche Richtlinien gelten:

```go
// Main app window — full media access
mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "App",
    Permissions: map[application.PermissionType]application.Permission{
        application.PermissionMicrophone: application.PermissionAllow,
        application.PermissionCamera:     application.PermissionAllow,
    },
})

// Settings window — no special capabilities
settingsWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Settings",
    // No Permissions entry — uses platform defaults
})

// Embedded content window — deny media capture
embeddedWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Embedded",
    Permissions: map[application.PermissionType]application.Permission{
        application.PermissionMicrophone: application.PermissionDeny,
        application.PermissionCamera:     application.PermissionDeny,
    },
})
```

## Windows-spezifische Überschreibung

Das fensterspezifische Feld `Windows.Permissions` (`map[CoreWebView2PermissionKind]CoreWebView2PermissionState`) funktioniert weiterhin und kann einzelne Funktionen überschreiben, nachdem die plattformübergreifende Zuordnung angewendet wurde. Verwenden Sie es, wenn Sie Zugriff auf WebView2-Berechtigungstypen benötigen, für die es keine plattformübergreifende Entsprechung gibt (z. B. `CoreWebView2PermissionKindOtherSensors`).

```go
Windows: application.WindowsWindow{
    Permissions: map[application.CoreWebView2PermissionKind]application.CoreWebView2PermissionState{
        application.CoreWebView2PermissionKindOtherSensors: application.CoreWebView2PermissionStateAllow,
    },
},
```

Die Auswertungsreihenfolge unter Windows lautet:

1. Plattformübergreifende `Permissions`-Zuordnung (legt den Status pro Typ über `SetPermission` fest)
2. `Windows.Permissions`-Zuordnung (überschreibt einzelne Typen)
3. Für jeden Typ, den keine der beiden Zuordnungen abdeckt: native Abfrage von WebView2 (wenn eine Richtlinie konfiguriert ist) oder automatische Zulassung (wenn keine Richtlinie konfiguriert ist – das bisherige Verhalten)

## Matrix der Plattformunterstützung

| Funktion | Linux | Windows | macOS |
| --- | --- | --- | --- |
| Mikrofon | ✅ | ✅ | Nur TCC |
| Kamera | ✅ | ✅ | Nur TCC |
| Geolokalisierung | ❌ noch nicht | ✅ | Nur TCC |
| Benachrichtigungen | ❌ noch nicht | ✅ | Nur TCC |
| Lesen aus der Zwischenablage | ❌ noch nicht | ✅ | Nur TCC |

## Fehlerbehebung

**`getUserMedia` schlägt unter Linux nach dem Upgrade weiterhin fehl**

Prüfen Sie, dass Sie weder `PermissionMicrophone: PermissionDeny` noch `PermissionCamera: PermissionDeny` ausdrücklich festgelegt haben. Standardmäßig (nicht festgelegt) ist die Medienaufnahme unter Linux zugelassen.

**Windows fragt nach nicht konfigurierten Berechtigungen**

Sobald `Permissions` einen Eintrag enthält, erteilt Wails nicht mehr die pauschale Freigabe `Allow`. Für nicht aufgeführte Funktionen wird die native Abfrage von WebView2 angezeigt. Fügen Sie für jede von Ihrer Anwendung verwendete Funktion einen ausdrücklichen `PermissionAllow`-Eintrag hinzu.

**macOS-Berechtigungen funktionieren nicht**

Die `Permissions`-Zuordnung hat unter macOS keine Wirkung. Stellen Sie sicher, dass `Info.plist` die richtigen Schlüssel für Nutzungsbeschreibungen enthält (`NSMicrophoneUsageDescription`, `NSCameraUsageDescription` usw.) und dass der Benutzer den Zugriff unter Systemeinstellungen → Datenschutz & Sicherheit gewährt hat.

**Geolokalisierung, Benachrichtigungen und Zwischenablage haben unter Linux keine Wirkung**

Unter Linux werden derzeit nur Kamera und Mikrofon berücksichtigt. Die Unterstützung anderer Funktionstypen ist noch nicht implementiert – sie bleiben unabhängig von der festgelegten Richtlinie verweigert.
