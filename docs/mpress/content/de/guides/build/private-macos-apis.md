---
title: "Private macOS-APIs"
description: "Alle Wails-Funktionen und -Optionen, die von privaten macOS-APIs abhängen, einschließlich Befehlen zur expliziten Aktivierung und Alternativen für öffentliche Builds."
slug: "guides/build/private-macos-apis"
sourcePath: "guides/build/private-macos-apis.md"
---

Wails verwendet standardmäßig öffentliche macOS-APIs. Das einzelne Go-Build-Tag `private_mac_apis` aktiviert die auf dieser Seite aufgeführten privaten WebKit- und AppKit-Aufrufe. Alle öffentlichen Go-Optionen und -Methoden bleiben in beiden Builds verfügbar. Ohne das Tag haben ausschließlich private Operationen keine Wirkung; Funktionen mit öffentlichen Alternativen verwenden diese Alternativen.

@note{type="caution" title="Privates macOS-Verhalten aktivieren"}
Das Festlegen einer Fensteroption aktiviert keine privaten APIs. Fügen Sie dem Build-Befehl `private_mac_apis` hinzu, um sie zu aktivieren. Das Tag gilt nur für macOS-Desktop-Builds, nicht für iOS, Android, Windows, Linux oder Server-Builds.

@end

## Private APIs aktivieren

```bash
# Build an application
wails3 build -tags private_mac_apis

# Run with live reload
EXTRA_TAGS=private_mac_apis wails3 dev

# Package an application
wails3 package GOOS=darwin EXTRA_TAGS=private_mac_apis

# Run a Go-only example from its directory
go run -tags private_mac_apis .
```

Verwenden Sie für einen direkten Produktions-Build `go build -tags production,private_mac_apis .`. Befolgen Sie bei Frontend-Beispielen deren README, um Bindings und Assets vor der Ausführung zu erstellen. Benutzerdefinierte oder ältere Taskfiles müssen `EXTRA_TAGS` an den Go-Compiler weiterleiten.

## Funktionsübersicht

| Funktion oder Wert | Was `private_mac_apis` aktiviert | Ohne das Tag |
| --- | --- | --- |
| `Mac.Backdrop: MacBackdropTransparent` | Transparentes WKWebView über dem nativen Fenster | Das native Fenster wird konfiguriert, aber die Webview bleibt undurchsichtig |
| `Mac.Backdrop: MacBackdropTranslucent` | Transparentes WKWebView, durch das die native Unschärfe sichtbar ist | Die Unschärfe wird hinter einer undurchsichtigen Webview konfiguriert |
| `Mac.Backdrop: MacBackdropLiquidGlass` | Transparentes WKWebView über der Glasschicht mit privater Steuerung des Webview-Hintergrunds | Die Glasschicht wird hinter einer undurchsichtigen Webview konfiguriert; für die Gestaltung werden öffentliche Alternativen verwendet |
| Leeren des Webview-Hintergrunds bei der Einrichtung von Liquid Glass | Private WebKit-Steuerung über `backgroundColor` | Öffentliches `underPageBackgroundColor` unter macOS 12+ oder eine Ebenenfarbe unter älteren macOS-Versionen; macht die Webview nicht transparent |
| `app.Window.NewNotchWindow(...)` | Transparente Webview innerhalb des geformten Notch-Panels | Das Panel funktioniert weiterhin einschließlich Positionierung und Animationen, seine Webview bleibt jedoch undurchsichtig |
| `Mac.LiquidGlass.Style` | Die bestehende Zuordnung nativer Stile von Wails einschließlich eines undokumentierten Werts für einen dunklen Stil | Verwendet öffentliche reguläre/transparente Stile und helle/dunkle Erscheinungsbilder; siehe die Wertetabelle unten |
| `Mac.LiquidGlass.GroupID` | Fordert für eine nicht leere Kennung eine private Glasgruppierung an | Wird ignoriert; es wird keine Gruppierung angefordert |
| `Mac.LiquidGlass.GroupSpacing` | Fordert für einen Wert größer als null einen privaten Gruppenabstand an | Wird ignoriert |
| `window.OpenDevTools()` und JavaScript-`Window.OpenDevTools()` | Öffnet den WebKit-Inspector unter macOS 12+ programmgesteuert | Keine Wirkung |
| `WebviewWindowOptions.OpenInspectorOnStartup: true` | Fordert unter macOS 12+ das programmgesteuerte Öffnen des Inspectors an, wenn das Fenster erstmals angezeigt wird | Keine Wirkung |
| Veraltete Aktivierung des Inspectors vor macOS 13.3 | Aktiviert die WebKit-Entwicklerextras, wenn die Inspector-Unterstützung einkompiliert ist | Keine Wirkung; die öffentliche Überprüfung mit Safari erfordert macOS 13.3+ |

## Transparenz und Hintergrund der Webview

**Erfordert private APIs:** die von `MacBackdropTransparent`, `MacBackdropTranslucent`, `MacBackdropLiquidGlass` und Notch-Fenstern verwendete Webview-Transparenz. Intern setzt Wails den privaten WebKit-Schlüssel `drawsBackground`. Ein transparenter HTML- oder CSS-Hintergrund allein kann ein undurchsichtiges natives WKWebView nicht transparent machen.

```go
Mac: application.MacWindow{
    // Requires -tags private_mac_apis for the blur to show through the webview.
    Backdrop: application.MacBackdropTranslucent,
},
```

Die private Operation für die Hintergrundfarbe der Webview verwendet den WebKit-Schlüssel `backgroundColor`. Die Einrichtung von Liquid Glass verwendet ihn, um den Webview-Hintergrund zu leeren. Ohne das Tag verwendet diese interne Operation unter macOS 12+ das öffentliche `underPageBackgroundColor` oder unter älteren macOS-Versionen die Ebene der Ansicht. Diese Alternativen machen die Webview nicht transparent.

`WebviewWindowOptions.BackgroundColour` und `window.SetBackgroundColour()` legen unter macOS die Farbe des **nativen Fensters** fest und erfordern selbst keine privaten APIs. Ebenso verwenden `Frameless` und `Mac.TitleBar.AppearsTransparent` öffentliche AppKit-APIs; die private Abhängigkeit ist die Webview-Transparenz, nicht die Transparenz der Titelleiste. Konfigurieren Sie für macOS-Hintergrundeffekte `Mac.Backdrop`, statt sich allein auf `BackgroundType` zu verlassen.

Siehe [Fensteroptionen](/features/windows/options/#mac-options), [rahmenlose Fenster](/features/windows/frameless/#with-transparent-background) und [Notch-Fenster](/features/windows/notch-windows/).

## Liquid-Glass-Werte

Wenn ein natives `NSGlassEffectView` verfügbar ist (macOS 26+), gelten die folgenden Zuordnungen. Dokumentiert sind nur die nativen Stilwerte `0` (regulär) und `1` (transparent). Die Go-Konstanten behalten in beiden Builds ihre bestehenden Werte.

| Wert von `MacLiquidGlassStyle` | Mit `private_mac_apis` | Ohne das Tag |
| --- | --- | --- |
| `LiquidGlassStyleAutomatic` (`0`) | Nativer regulärer Stil (`0`) | Nativer regulärer Stil (`0`) |
| `LiquidGlassStyleLight` (`1`) | Bestehende Zuordnung nativer Stile (`1`, transparent) | Nativer regulärer Stil (`0`) mit Aqua-Erscheinungsbild |
| `LiquidGlassStyleDark` (`2`) | **Undokumentierter nativer Stilwert `2`** | Nativer regulärer Stil (`0`) mit Dark-Aqua-Erscheinungsbild |
| `LiquidGlassStyleVibrant` (`3`) | Entspricht dem vorhandenen hellen/nativen transparenten Stil (`1`) | Nativer transparenter Stil (`1`) |

Automatic und Vibrant verwenden dokumentierte native Stilwerte, für einen fensterfüllenden Liquid-Glass-Hintergrund ist das Tag jedoch weiterhin für die **Transparenz der Webview** erforderlich. Das Erscheinungsbild von Light unterscheidet sich zwischen den beiden Builds. Die Zuordnung des privaten Stils garantiert nicht, dass eine zukünftige macOS-Version denselben Effekt rendert.

**Immer privat:** `GroupID` und `GroupSpacing`. Wails prüft vor dem Anfordern der Gruppierung, ob die privaten Selektoren `setGroupIdentifier:`, `setGroupName:` und `setGroupSpacing:` vorhanden sind. Ohne das Tag bewirken diese Operationen nichts. Das Aktivieren des Tags garantiert nicht, dass die ausgeführte macOS-Version diese Selektoren unterstützt.

`MacLiquidGlass.Material`, `CornerRadius` und `TintColor` erfordern selbst keine privaten APIs. Auf macOS-Versionen ohne native Liquid-Glass-Unterstützung konfiguriert Wails eine transluzente Ausweichdarstellung; damit sie durch die Webview sichtbar ist, ist das Tag weiterhin erforderlich.

## Web-Inspector

**Erfordert private APIs:** der Aufruf von `OpenDevTools()` oder das Setzen von `OpenInspectorOnStartup: true`, um den Inspector von WebKit aus der Anwendung heraus zu öffnen. Wails verwendet den privaten Selektor `_inspector`. Ohne `private_mac_apis` bewirken diese Operationen stillschweigend nichts.

Die Inspector-Unterstützung muss außerdem einkompiliert sein. Die vorhandenen Tags `production` und `devtools` behalten ihre Bedeutung:

| Build-Tags | Programmgesteuertes Öffnen des Inspectors | Öffentliche Safari-Inspektion unter macOS 13.3+ |
| --- | --- | --- |
| Keine | Keine Wirkung | Aktiviert |
| `private_mac_apis` | Unter macOS 12+ aktiviert | Aktiviert |
| `production` | Keine Wirkung | Deaktiviert |
| `production,private_mac_apis` | Keine Wirkung | Deaktiviert |
| `production,devtools` | Keine Wirkung | Aktiviert |
| `production,devtools,private_mac_apis` | Unter macOS 12+ aktiviert | Aktiviert |

Unter macOS 13.3+ aktiviert Wails die Safari-Inspektion über die öffentliche API `WKWebView.inspectable`; dafür sind keine privaten APIs erforderlich. Vor macOS 13.3 verwendet die Ausweichlösung die private Einstellung `developerExtrasEnabled` und erfordert daher neben einkompilierter Inspector-Unterstützung auch `private_mac_apis`.

```bash
# Production build with programmatic inspector support
go build -tags production,devtools,private_mac_apis .
```

## Diese Liste pflegen

Alle privaten nativen Aufrufe sind in `v3/pkg/application/mac_private_api_darwin.go` isoliert; Standard-Builds wählen `mac_public_api_darwin.go` aus. Diese Übersicht umfasst Transparenz, die Hintergrundfarbe der Webview, Glass-Stile, Glass-Gruppierung, das Öffnen des Inspectors und die Aktivierung des Legacy-Inspectors. Bei Änderungen an diesen Implementierungen sollten diese Seite und die Dokumentation der betroffenen Optionen gemeinsam aktualisiert werden.
