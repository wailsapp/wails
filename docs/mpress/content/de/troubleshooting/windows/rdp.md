---
title: "WebView2 hängt über Remotedesktop (RDP)"
description: "Behebt mehrere Sekunden dauernde Stillstände der WebView2-Benutzeroberfläche, wenn eine Wails-App über eine RDP-Sitzung ausgeführt wird, die während der Sitzung die Monitor-DPI ändert."
slug: "troubleshooting/windows/rdp"
sourcePath: "troubleshooting/windows/rdp.md"
---

## Problem

Wenn eine Wails-Anwendung über eine Remotedesktop-Sitzung (RDP) verwendet wird, kann die Benutzeroberfläche bei häufigen Interaktionen mehrere Sekunden lang stillstehen:

- Vom Klick bis zum sichtbaren Inhalt dauert das Öffnen eines Popup-Fensters etwa 4 bis 8 Sekunden.
- Das Schließen eines Fensters blockiert das übergeordnete Fenster etwa 2 Sekunden lang.
- Der langsame Zustand bleibt auch nach erneuten Verbindungen bestehen und verschwindet erst nach einem Neustart des Hostrechners.

Dies tritt am häufigsten mit dem Microsoft-Remote-Desktop-Client unter iOS auf, der im Laufe der Sitzung einen für Retina optimierten virtuellen Monitor bereitstellt. Jeder RDP-Client, der während der Sitzung einen Monitor mit einem anderen DPI-Kontext hinzufügt, kann dasselbe Verhalten auslösen.

## Ursache

Standardmäßig verwendet WebView2 fensterbasiertes Hosting, bei dem sich seine Compositor-Oberfläche in einem untergeordneten Fenster befindet. Wenn ein RDP-Client einen Monitor hinzufügt, dessen DPI-Kontext von dem der Sitzung abweicht, erzwingt jeder Aufruf des WebView2-Controllers (`PutIsVisible`, `MoveFocus`, erstes Rendern und Freigeben der Oberfläche) ein synchrones erneutes Marshaling durch DirectComposition. Jedes erneute Marshaling blockiert den UI-Thread etwa 2 Sekunden lang. Deshalb summieren sich die Stillstände in Apps, die viele Popups verwenden.

Eine native Anwendung auf Basis von Win32 und WebView2 ist auf demselben Rechner nicht betroffen, da sie visuelles Hosting verwendet. Dies weist darauf hin, dass der Hostingmodus die Ursache ist und nicht ein allgemeines Problem mit WebView2 oder dem Windows-Compositor.

## Lösung

Aktivieren Sie visuelles Hosting, indem Sie in Ihren Windows-Optionen `UseVisualHosting` festlegen. Beim visuellen Hosting wird die Compositor-Oberfläche von WebView2 über ein DirectComposition-Visual verwaltet, das dem Host gehört. Änderungen des DPI-Kontexts lösen daher kein synchrones erneutes Marshaling mehr aus.

```go
package main

import "github.com/wailsapp/wails/v3/pkg/application"

func main() {
    app := application.New(application.Options{
        Windows: application.WindowsOptions{
            UseVisualHosting: true,
        },
    })

    // ... create your windows, then:
    app.Run()
}
```

Nach der Aktivierung öffnen sich Popups innerhalb der normalen Navigationszeit (etwa 150 bis 500 ms), und das Schließen eines Fensters blockiert das übergeordnete Fenster nicht mehr.

@note{type="caution"}
`UseVisualHosting` muss vor `app.Run()` festgelegt werden. Wails liest die Option beim Start der Anwendung und setzt die Umgebungsvariable `COREWEBVIEW2_FORCED_HOSTING_MODE` auf `COREWEBVIEW2_HOSTING_MODE_WINDOW_TO_VISUAL`, bevor die WebView2-Umgebung initialisiert wird. Eine spätere Festlegung hat keine Wirkung.

@end

Der Standardwert der Option ist `false`, sodass fensterbasiertes Hosting die Standardeinstellung bleibt. Bei bestehenden Apps ändert sich das Verhalten nur, wenn sie die Option ausdrücklich aktivieren.

## Wann die Option aktiviert werden sollte

Legen Sie `UseVisualHosting: true` fest, wenn Ihre App regelmäßig über RDP verwendet wird – insbesondere mit dem Microsoft-Remote-Desktop-Client für iOS – und beim Öffnen oder Schließen von Fenstern mehrere Sekunden dauernde Stillstände auftreten. Wenn Ihre App nicht über RDP ausgeführt wird, benötigen Sie diese Option nicht und können die Standardeinstellung beibehalten.

## Referenzen

- [WebView2: fensterbasiertes und visuelles Hosting im Vergleich](https://learn.microsoft.com/en-us/microsoft-edge/webview2/concepts/windowed-vs-visual-hosting)
- [WebView2Feedback-Issue Nr. 5248](https://github.com/MicrosoftEdge/WebView2Feedback/issues/5248)
- [WebView2Feedback-Issue Nr. 4485](https://github.com/MicrosoftEdge/WebView2Feedback/issues/4485)
