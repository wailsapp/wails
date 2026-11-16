---
title: "Notch-Fenster"
description: "Native macOS-Fenster erstellen, die am Kameragehäuse angedockt sind."
slug: "features/windows/notch-windows"
sourcePath: "features/windows/notch-windows.md"
---

`NewNotchWindow` erstellt ein geformtes, nicht aktivierendes macOS-Panel, das am Kameragehäuse angedockt ist. Wails verwaltet die native Positionierung, die schwarzen seitlichen Verbindungen, die transparente äußere Zeichenfläche, die Fensterebene, das Verhalten in Spaces sowie optionale Ein- und Ausblendanimationen. Ihre Webinhalte belegen nur das angeforderte innere Rechteck. Sobald der Mauszeiger in das Fenster bewegt wird, erhält dessen Webview unmittelbar den Tastaturfokus, ohne die Anwendung zu aktivieren.

@note{type="caution" title="Die Webview-Transparenz verwendet eine private API"}
`NewNotchWindow` benötigt `-tags private_mac_apis`, damit transparente Webinhalte die native Notch-Form sichtbar machen. Ohne dieses Tag funktionieren das Panel, die Positionierung und die Animationen weiterhin, der Webview bleibt jedoch undurchsichtig. Siehe [Private macOS-APIs](/guides/build/private-macos-apis/#webview-transparency-and-background).

@end

<figure>
  <img
    src="/images/notch-notification.gif"
    alt="Eine Wails-Notch-Benachrichtigung, die aus dem Kameragehäuse des MacBook nach unten gleitet, Live-Systemmetriken anzeigt und anschließend wieder unter der Notch verschwindet"
    loading="lazy"
    decoding="async"
    style="width: 100%; border-radius: 0.75rem"
  />
  <figcaption>
    Eine animierte Notch-Benachrichtigung mit einem persistenten Webview und nativen
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
| `Width` | `int` | `660` | Nutzbare Breite des Webviews innerhalb der nativen geformten Ränder. |
| `Height` | `int` | `92` | Nutzbare Höhe des Webviews innerhalb der nativen geformten Ränder. |
| `Animated` | `bool` | `false` | Schiebt das Fenster bei `Show` nach unten und bei `Hide` nach oben. |
| `AnimationSpeed` | `time.Duration` | `420ms` | Dauer des Einblendens. Das Ausblenden dauert zwei Drittel dieses Werts, standardmäßig `280ms`. |
| `Screen` | `*Screen` | primärer Bildschirm | Legt einen bestimmten Bildschirm als Ziel fest. Andernfalls verwendet Wails den primären Bildschirm. |
| `WindowOptions` | `WebviewWindowOptions` | Standardwerte | Gibt den Namen, die URL oder das HTML, CSS, JavaScript, die Tastenkürzel und weiteres Webview-Verhalten an. |

`NewNotchWindow` bestimmt die äußere Größe, Position, den Rahmen, die Transparenz, die Größenänderungsrichtlinie, die native Panelklasse, die Fensterebene und das Sammlungsverhalten. Werte für diese Felder in `WindowOptions` werden absichtlich ersetzt. Andere Felder bleiben erhalten. Das native Ziehen über den Hintergrund und CSS-Ziehbereiche sind deaktiviert, damit das Fenster am Kameragehäuse angedockt bleibt.

Das zurückgegebene `NotchWindow` stellt absichtlich nur `Show`, `Hide`, `Visibility` und `Close` bereit; die native Geometrie kann über dieses abstrakte Handle nicht geändert werden.

@note{type="note"}
Notch-Fenster erfordern macOS. Auf einem Mac ohne Kameragehäuse positioniert Wails das Fenster oben mittig unterhalb der Menüleiste. Auf nicht unterstützten Plattformen gibt `NewNotchWindow` ein inaktives Handle zurück, dessen Lebenszyklusmethoden sichere No-ops sind und dessen `Visibility` immer false ist.

@end

## Lebenszyklus

- `Show` blendet das vorhandene native Fenster ein. Bei aktivierter Animation gleitet es von oberhalb des Bildschirms nach unten.
- `Hide` hält das Fenster zur Wiederverwendung am Leben. Bei aktivierter Animation gleitet es zurück über den Bildschirmrand, bevor es aus der Fensterreihenfolge entfernt wird. Webview, JavaScript-Zustand, Bindings und Ereignis-Listener bleiben geladen, das ausgeblendete Fenster besitzt jedoch kein Hover-Ziel; die Anwendung muss `Show` aufrufen, um es erneut einzublenden.
- `Visibility` meldet den aktuellen nativen Sichtbarkeitsstatus.
- `Close` zerstört das native Fenster dauerhaft. Erstellen Sie ein neues, bevor Sie diese Benachrichtigung erneut anzeigen.
- Beim Eintritt des Mauszeigers wird dieses Notch-Fenster in den Vordergrund gebracht und sein Webview für die sofortige Tastaturinteraktion fokussiert, während das nicht aktivierende Panelverhalten erhalten bleibt.

Jeder Aufruf von `NewNotchWindow` erstellt ein unabhängiges Fenster mit eigenen Inhalten sowie einem eigenen Sichtbarkeits- und Animationszustand. Die Fenster verwenden dieselbe native Ebene, sodass die zuletzt eingeblendete oder vom Mauszeiger betretene Instanz im Vordergrund erscheint und frühere Instanzen überlagern kann. Wenn Fenster auf unterschiedliche Bildschirme ausgerichtet sind, wird jedes am Kameragehäuse des jeweiligen Bildschirms zentriert. macOS koordiniert keine Notch-Fenster verschiedener Anwendungen; zeigen separate Anwendungen Fenster an derselben Position und auf derselben Ebene an, erscheint das zuletzt in der Fensterreihenfolge platzierte Fenster im Vordergrund.

Für Benachrichtigungen verwenden Anwendungen üblicherweise ein ausgeblendetes Fenster erneut oder verwalten eine eigene Warteschlange in einem Webview. Wails erzwingt weder Warteschlangen, Ersetzungen, automatisches Ausblenden noch eine Richtlinie, die nur ein einzelnes Fenster zulässt.

@note{type="note"}
Persistente eingeklappte Oberflächen, die sich beim Darüberfahren mit dem Mauszeiger wieder öffnen, sind bewusst vom Ausblenden von Benachrichtigungen getrennt und werden unter [#6009](https://github.com/wailsapp/wails/issues/6009) nachverfolgt.

@end

Ein kompaktes Systemmonitor-Beispiel, dessen aktueller JavaScript-Zustand wiederholte Ein- und Ausblendzyklen von Benachrichtigungen überdauert, finden Sie im [Beispiel notch-notification](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/notch-notification).
