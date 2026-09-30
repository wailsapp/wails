---
title: "Notch-Fenster"
description: "Erstellt native macOS-Fenster, die am Kameragehäuse angedockt sind."
slug: "features/windows/notch-windows"
sourcePath: "features/windows/notch-windows.md"
---

`NewNotchWindow` erstellt ein geformtes, nicht aktivierendes macOS-Panel, das am Kameragehäuse angedockt ist. Wails steuert die native Positionierung, die schwarzen seitlichen Ansätze, die transparente äußere Zeichenfläche, die Fensterebene, das Verhalten in Spaces sowie optionale Ein- und Ausblendeanimationen. Ihre Webinhalte belegen nur das angeforderte innere Rechteck. Wenn Sie den Zeiger in das Fenster bewegen, wird dessen Webview sofort zum Key-Window, ohne die Anwendung zu aktivieren.

@note{type="caution" title="Webview-Transparenz verwendet eine private API"}
`NewNotchWindow` benötigt `-tags private_mac_apis`, damit transparente Webinhalte die native Notch-Form sichtbar machen. Ohne dieses Tag funktionieren Panel, Positionierung und Animationen weiterhin, aber die Webview bleibt undurchsichtig. Siehe [Private macOS-APIs](/guides/build/private-macos-apis/#webview-transparency-and-background).

@end

<figure>
  <img
    src="/images/notch-notification.gif"
    alt="Eine Wails-Notch-Benachrichtigung, die vom Kameragehäuse des MacBook nach unten gleitet, Live-Systemmetriken anzeigt und sich wieder unter der Notch verbirgt"
    loading="lazy"
    decoding="async"
    style="width: 100%; border-radius: 0.75rem"
  />
  <figcaption>
    Eine animierte Notch-Benachrichtigung mit einer persistenten Webview und nativen
    Ein- und Ausblendübergängen.
  </figcaption>
</figure>

```go
alert := app.Window.NewNotchWindow(application.NotchWindowOptions{
    Width:    660,
    Height:   92,
    Animated: true,
    WindowOptions: application.WebviewWindowOptions{
        Name: "alert",
        URL:  "/alert",
    },
})

alert.Show()
alert.Hide()
visible := alert.Visibility()
alert.Close()
```

## Optionen

| Feld | Typ | Standardwert | Beschreibung |
| --- | --- | --- | --- |
| `Width` | `int` | `660` | Nutzbare Breite der Webview innerhalb der nativen geformten Ränder. |
| `Height` | `int` | `92` | Nutzbare Höhe der Webview innerhalb der nativen geformten Ränder. |
| `Animated` | `bool` | `false` | Schiebt das Fenster bei `Show` nach unten und bei `Hide` nach oben. |
| `AnimationSpeed` | `time.Duration` | `420ms` | Dauer des Einblendens. Das Ausblenden verwendet zwei Drittel dieses Werts, standardmäßig `280ms`. |
| `Screen` | `*Screen` | primärer Bildschirm | Verwendet einen bestimmten Bildschirm. Andernfalls verwendet Wails den primären Bildschirm. |
| `WindowOptions` | `WebviewWindowOptions` | Standardwerte | Gibt den Namen, die URL oder HTML, CSS, JavaScript, Tastenbelegungen und weiteres Webview-Verhalten an. |

`NewNotchWindow` steuert äußere Größe, Position, Rahmen, Transparenz, Richtlinie zur Größenänderung, native Panel-Klasse, Fensterebene und Sammlungsverhalten. Werte für diese Felder in `WindowOptions` werden absichtlich ersetzt. Andere Felder bleiben erhalten. Das native Ziehen über den Hintergrund und CSS-Ziehbereiche sind deaktiviert, damit das Fenster am Kameragehäuse angedockt bleibt.

Das zurückgegebene `NotchWindow` stellt absichtlich nur `Show`, `Hide`, `Visibility` und `Close` bereit; die native Geometrie kann über das abstrakte Handle nicht geändert werden.

@note{type="note"}
Notch-Fenster erfordern macOS. Auf einem Mac ohne Kameragehäuse positioniert Wails das Fenster oben mittig unterhalb der Menüleiste. Auf nicht unterstützten Plattformen gibt `NewNotchWindow` ein inaktives Handle zurück, dessen Lebenszyklusmethoden sichere No-Ops sind und dessen `Visibility` immer false ist.

@end

## Lebenszyklus

- `Show` blendet das vorhandene native Fenster ein. Bei aktivierter Animation gleitet es von oberhalb des Bildschirms nach unten.
- `Hide` lässt das Fenster zur Wiederverwendung bestehen. Bei aktivierter Animation gleitet es zurück über den Bildschirmrand, bevor es aus der Anzeigereihenfolge entfernt wird. Webview, JavaScript-Zustand, Bindings und Event-Listener bleiben geladen, aber das ausgeblendete Fenster bietet kein Hover-Ziel. Die Anwendung muss `Show` aufrufen, um es wieder einzublenden.
- `Visibility` meldet den aktuellen nativen Sichtbarkeitsstatus.
- `Close` zerstört das native Fenster dauerhaft. Erstellen Sie ein neues, bevor Sie diese Benachrichtigung erneut anzeigen.
- Beim Eintritt des Zeigers wird dieses Notch-Fenster nach vorn gebracht und seine Webview für die sofortige Tastaturinteraktion fokussiert, während das nicht aktivierende Panel-Verhalten erhalten bleibt.

Jeder Aufruf von `NewNotchWindow` erstellt ein unabhängiges Fenster mit eigenem Inhalt, eigener Sichtbarkeit und eigenem Animationszustand. Die Fenster verwenden dieselbe native Ebene, sodass die zuletzt angezeigte oder mit dem Zeiger betretene Instanz im Vordergrund erscheint und frühere Instanzen überlagern kann. Wenn Fenster unterschiedliche Bildschirme verwenden, wird jedes mittig am Kameragehäuse des jeweiligen Bildschirms ausgerichtet. macOS koordiniert keine Notch-Fenster verschiedener Anwendungen. Wenn separate Anwendungen Fenster an derselben Position und auf derselben Ebene anzeigen, erscheint das zuletzt angeordnete Fenster im Vordergrund.

Für Benachrichtigungen verwenden Anwendungen häufig ein ausgeblendetes Fenster erneut oder verwalten eine eigene Warteschlange in einer Webview. Wails erzwingt weder eine Warteschlange noch Ersetzung, automatisches Schließen oder eine Einzelfenster-Richtlinie.

@note{type="note"}
Persistente eingeklappte Oberflächen, die sich beim Darüberbewegen des Zeigers wieder öffnen, sind bewusst vom Ausblenden von Benachrichtigungen getrennt und werden unter [#6009](https://github.com/wailsapp/wails/issues/6009) nachverfolgt.

@end

Siehe das [Beispiel für Notch-Benachrichtigungen](https://github.com/wailsapp/wails/tree/master/v3/examples/notch-notification) für einen kompakten Systemmonitor, dessen aktueller JavaScript-Zustand wiederholte Ein- und Ausblendzyklen der Benachrichtigung überdauert.
