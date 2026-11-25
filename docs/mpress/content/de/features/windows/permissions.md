---
title: "Berechtigungen"
description: "Anfragen von Webinhalten für Kamera, Mikrofon, Standort und andere Funktionen steuern"
slug: "features/windows/permissions"
sourcePath: "features/windows/permissions.md"
---

Webinhalte, die `navigator.mediaDevices.getUserMedia()`, die Geolocation API oder die Notifications API aufrufen, benötigen eine Genehmigung oder Ablehnung dieser Anfragen durch die Hostanwendung. Wails stellt dafür in `WebviewWindowOptions` eine plattformübergreifende `Permissions`-Map bereit, mit der Sie dies deklarativ steuern können — ohne plattformspezifischen Code.

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

Kamera- und Mikrofonanfragen aus den Webinhalten dieses Fensters werden ohne Browserabfrage genehmigt.

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
| `PermissionDefault` | 0 | Die native Verarbeitung der Plattform verwenden (siehe unten) |
| `PermissionAllow` | 1 | Ohne Nachfrage genehmigen |
| `PermissionDeny` | 2 | Ohne Nachfrage ablehnen |

`PermissionDefault` ist der Nullwert. Nicht gesetzte Map-Einträge verhalten sich daher wie die Standardeinstellung.

## Plattformverhalten

Jede Plattform verarbeitet `PermissionDefault` anders, da sich das native Verhalten der zugrunde liegenden Webviews unterscheidet.

### Linux (WebKitGTK)

WebKitGTK besitzt **keine native Berechtigungsabfrage**. Ist kein Handler registriert, werden alle Anfragen stillschweigend abgelehnt — deshalb gab `getUserMedia` vor Einführung dieser Funktion immer `NotAllowedError` zurück.

Derzeit verarbeitet Wails unter Linux Anfragen für **Kamera und Mikrofon**. Standort, Benachrichtigungen und das Lesen der Zwischenablage sind noch nicht angebunden und bleiben unabhängig von der festgelegten Richtlinie gesperrt.

| Richtlinie | Kamera/Mikrofon | Standort, Benachrichtigungen, Zwischenablage |
| --- | --- | --- |
| `PermissionDefault` | **Zugelassen** (stellt getUserMedia wieder her) | Immer abgelehnt |
| `PermissionAllow` | Zugelassen | Immer abgelehnt (noch nicht implementiert) |
| `PermissionDeny` | Abgelehnt | Immer abgelehnt |

### Windows (WebView2)

WebView2 besitzt eine native Berechtigungsabfrage und eine Berechtigungs-API für die einzelnen Typen. Alle fünf Funktionstypen werden vollständig unterstützt.

| Richtlinie | Verhalten |
| --- | --- |
| `PermissionDefault` | WebView2 zeigt seine native Berechtigungsabfrage des Betriebssystems bzw. Browsers an |
| `PermissionAllow` | Stillschweigend genehmigt |
| `PermissionDeny` | Stillschweigend abgelehnt |

**Wichtig:** Vor Einführung dieser Funktion rief Wails `SetGlobalPermission(Allow)` bedingungslos auf und genehmigte damit stillschweigend alle Funktionen. Sobald `Permissions` einen Eintrag enthält, wird diese pauschale Genehmigung jetzt **nicht** gesetzt. Für nicht konfigurierte Funktionen wird stattdessen die native WebView2-Abfrage angezeigt, anstatt sie automatisch zu genehmigen.

Das bedeutet: Sobald Sie `Permissions` unter Windows konfigurieren, wird für jede nicht ausdrücklich aufgeführte Funktion eine Abfrage angezeigt, anstatt sie stillschweigend zuzulassen. Legen Sie die benötigten Funktionen ausdrücklich fest.

### macOS (WKWebView + TCC)

Unter macOS müssen zwei Ebenen zustimmen. WKWebView fragt die Anwendung, bevor eine Aufnahmesitzung gestartet wird; diese Anfrage wird durch die `Permissions`-Map beantwortet. Darunter kontrolliert das Datenschutzframework des Systems (TCC) den Zugriff auf das Gerät selbst: Die Betriebssystemabfrage erscheint, wenn die App erstmals tatsächlich auf die Kamera oder das Mikrofon zugreift. Die Auswahl des Benutzers wird für jede App unter Systemeinstellungen → Datenschutz & Sicherheit gespeichert.

Wails verarbeitet unter macOS 12 und neuer Anfragen für **Kamera und Mikrofon**. Für Standort, Benachrichtigungen und das Lesen der Zwischenablage gibt es keine Entsprechung zu `WKUIDelegate`. Sie sind daher nicht angebunden und werden allein TCC überlassen — wie unter Linux hat die dafür festgelegte Richtlinie keine Wirkung.

| Richtlinie | Kamera/Mikrofon | Standort, Benachrichtigungen, Zwischenablage |
| --- | --- | --- |
| `PermissionDefault` | WebKit zeigt eine eigene Berechtigungsabfrage an | Nur TCC |
| `PermissionAllow` | Die WebKit-Abfrage wird übersprungen — **TCC gilt weiterhin** | Nur TCC |
| `PermissionDeny` | Von WebKit abgelehnt; die Systemautorisierung erfolgt separat | Nur TCC |

`PermissionAllow` genehmigt die Anfrage des Webviews, nicht den Gerätezugriff. Die Systemautorisierung gilt weiterhin unabhängig davon; eine in den Systemeinstellungen abgelehnte App bleibt gesperrt. Nur die WebKit-eigene Abfrage wird übersprungen.

Unterhalb von macOS 12 ist die Delegate-Methode nicht vorhanden. Dort wird die Map daher ignoriert und jede Anfrage auf die WebKit-Abfrage zurückgeführt.

Stellen Sie sicher, dass Ihre `Info.plist` die passenden Schlüssel für Nutzungsbeschreibungen enthält:

```xml
<key>NSMicrophoneUsageDescription</key>
<string>Used for voice input</string>
<key>NSCameraUsageDescription</key>
<string>Used for video calls</string>
```

@note{type="caution" title="Nutzungsbeschreibungen sind erforderlich"}

WebKit kann die Systemautorisierung anfordern, bevor es den Wails-Berechtigungs-Delegate aufruft. `PermissionDeny` blockiert die WebKit-Anfrage, garantiert jedoch nicht, dass keine TCC-Abfrage oder AVFoundation-Autorisierungsanfrage erfolgt.

Fügen Sie die passenden Nutzungsbeschreibungen für die von Ihrer App angeforderten Geräte sowie bei Verwendung der Sandbox die entsprechenden Geräte-Entitlements hinzu. Verwenden Sie `PermissionDeny` nicht als Ersatz für diese Deklarationen.

@end

## Gängige Muster

### App zur Medienaufnahme

Kamera und Mikrofon auf allen Plattformen genehmigen:

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionAllow,
},
```

Unter **Linux** werden beide Geräte dadurch ausdrücklich zugelassen; andere Funktionen bleiben verweigert. Unter **Windows** werden beide zugelassen; für jede andere, nicht aufgeführte Funktion erscheint eine native Abfrage. Unter **macOS** werden beide auf der WebKit-Ebene zugelassen, sodass keine Browserabfrage erscheint. TCC fragt beim ersten Zugriff weiterhin nach der Berechtigung für die Geräte selbst, und beide Schlüssel für die Nutzungsbeschreibung müssen in `Info.plist` vorhanden sein.

### Medienaufnahme unter Linux verweigern

Linux erlaubt den Zugriff auf Kamera und Mikrofon standardmäßig. So deaktivierst du ihn:

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionDeny,
    application.PermissionCamera:     application.PermissionDeny,
},
```

### Alle Funktionen unter Windows zulassen

So gewährst du ohne Abfrage Zugriff auf alle Funktionen:

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

Das fensterspezifische Feld `Windows.Permissions` (`map[CoreWebView2PermissionKind]CoreWebView2PermissionState`) funktioniert weiterhin und kann einzelne Funktionen überschreiben, nachdem die plattformübergreifende Zuordnung angewendet wurde. Verwende es, wenn du Zugriff auf WebView2-Berechtigungstypen benötigst, für die es keine plattformübergreifende Entsprechung gibt (z. B. `CoreWebView2PermissionKindOtherSensors`).

```go
Windows: application.WindowsWindow{
    Permissions: map[application.CoreWebView2PermissionKind]application.CoreWebView2PermissionState{
        application.CoreWebView2PermissionKindOtherSensors: application.CoreWebView2PermissionStateAllow,
    },
},
```

Unter Windows gilt folgende Auswertungsreihenfolge:

1. Plattformübergreifende Zuordnung `Permissions` (legt den Status jedes Typs über `SetPermission` fest)
2. Zuordnung `Windows.Permissions` (überschreibt einzelne Typen)
3. Für jeden Typ, der von keiner der beiden Zuordnungen abgedeckt wird: native Abfrage von WebView2 (wenn eine Richtlinie konfiguriert ist) oder automatische Zulassung (wenn keine Richtlinie konfiguriert ist – das bisherige Verhalten)

## Matrix der Plattformunterstützung

| Funktion | Linux | Windows | macOS |
| --- | --- | --- | --- |
| Mikrofon | ✅ | ✅ | ✅ (macOS 12+) |
| Kamera | ✅ | ✅ | ✅ (macOS 12+) |
| Geolokalisierung | ❌ noch nicht | ✅ | ❌ noch nicht |
| Benachrichtigungen | ❌ noch nicht | ✅ | ❌ noch nicht |
| Lesen aus der Zwischenablage | ❌ noch nicht | ✅ | ❌ noch nicht |

Wo macOS mit ✅ gekennzeichnet ist, beantwortet die Richtlinie die Anfrage von WebKit; zusätzlich steuert TCC weiterhin den Zugriff auf das Gerät. Wo es mit ❌ gekennzeichnet ist, wird die Funktion ausschließlich TCC überlassen.

## Fehlerbehebung

**`getUserMedia` schlägt unter Linux nach dem Upgrade weiterhin fehl**

Prüfe, ob du `PermissionMicrophone: PermissionDeny` oder `PermissionCamera: PermissionDeny` ausdrücklich festgelegt hast. Standardmäßig (nicht festgelegt) ist die Medienaufnahme unter Linux zulässig.

**Windows fragt nach Berechtigungen, die ich nicht konfiguriert habe**

Sobald `Permissions` einen Eintrag enthält, legt Wails nicht mehr die pauschale Gewährung `Allow` fest. Für nicht aufgeführte Funktionen erscheint die native Abfrage von WebView2. Füge für jede von deiner App verwendete Funktion einen ausdrücklichen `PermissionAllow`-Eintrag hinzu.

**macOS-Berechtigungen funktionieren nicht**

`Permissions` unterstützt Kamera und Mikrofon unter macOS 12 und neuer. Geolokalisierung, Benachrichtigungen und das Lesen aus der Zwischenablage sind noch nicht angebunden und ignorieren die Richtlinie. Wo sie berücksichtigt wird, steuert TCC zusätzlich weiterhin den Zugriff auf das Gerät: `PermissionAllow` entfernt die Abfrage von WebKit, nicht jedoch die Systemabfrage. Stelle sicher, dass `Info.plist` die richtigen Schlüssel für die Nutzungsbeschreibung enthält (`NSMicrophoneUsageDescription`, `NSCameraUsageDescription`) und dass der Benutzer den Zugriff unter Systemeinstellungen → Datenschutz & Sicherheit gewährt hat.

**Meine macOS-App wird beendet, wenn Webinhalte auf die Kamera oder das Mikrofon zugreifen möchten**

Prüfe, ob `Info.plist` für jedes von deiner App angeforderte Gerät `NSCameraUsageDescription` oder `NSMicrophoneUsageDescription` enthält. Fehlende Beschreibungen können dazu führen, dass die Anfrage abgelehnt wird; ein Zugriff auf AVFoundation ohne erforderliche Beschreibung kann die App beenden. `PermissionDeny` steuert die Entscheidung von WebKit und ersetzt diese Schlüssel nicht.

**Geolokalisierung, Benachrichtigungen und Zwischenablage haben unter Linux keine Wirkung**

Unter Linux werden derzeit nur Kamera und Mikrofon verarbeitet. Die Unterstützung anderer Funktionstypen ist noch nicht implementiert – sie bleiben unabhängig von der festgelegten Richtlinie verweigert.
