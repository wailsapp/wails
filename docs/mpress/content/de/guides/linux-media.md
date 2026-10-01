---
title: "Lokale Audio- und Videodateien abspielen"
description: "Eingebundene Medien unter Linux über größenbegrenzte Blob-URLs abspielen und anschließend freigeben."
slug: "guides/linux-media"
sourcePath: "guides/linux-media.md"
---

Verwenden Sie `Media.SetSource`, um kurze lokale Clips in einer Wails-v3-Anwendung abzuspielen. Unter Linux übergibt WebKitGTK die Medienwiedergabe an GStreamer, das `wails://`-URLs nicht direkt laden kann. Die Hilfsfunktion empfängt den Clip über einen Wails-Stream und weist dem Player eine Blob-URL zu.

Desktop-Übertragungen nutzen den vorhandenen Wails-Asset-Transport, ohne einen lauschenden Socket zu öffnen. Die API funktioniert mit Audio- und Videoelementen. Gewöhnliche HTTP/HTTPS-Medien können direkt das `src` des Players verwenden.

## Mediendateien registrieren

Stellen Sie ein Dateisystem mit den Clips bereit, die Ihr Frontend abspielen darf, und registrieren Sie den Medienhandler auf einem benannten Stream:

```go
import (
    "embed"
    "io/fs"
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/services/media"
)

//go:embed clips
var clips embed.FS

func registerMedia(app *application.App) {
    mediaFiles, err := fs.Sub(clips, "clips")
    if err != nil {
        log.Fatal(err)
    }
    handler, err := media.NewHandler(mediaFiles, 32 << 20) // 32 MiB per file
    if err != nil {
        log.Fatal(err)
    }
    app.HandleStream("media", handler)
}
```

Rufen Sie `registerMedia(app)` nach dem Erstellen der Anwendung und vor `app.Run()` auf.

Verwenden Sie für Dateien auf dem Datenträger `os.OpenRoot(directory)` und übergeben Sie `root.FS()` an `media.NewHandler`. Lassen Sie das Stammverzeichnis bis zur Rückkehr von `app.Run()` geöffnet und schließen Sie es anschließend. Wählen Sie ein Verzeichnis, das nur vom Frontend lesbare Dateien enthält; `os.Root` begrenzt den Zugriff auch dann, wenn ein symbolischer Link aus diesem Verzeichnis herausführt.

## Einen Clip laden

Erstellen Sie einen Player:

```html
<video id="player" controls></video>
```

Für ein Frontend mit der npm-Runtime:

```javascript
import { Media } from '@wailsio/runtime';

const player = document.getElementById('player');

try {
    await Media.SetSource(player, 'media', 'welcome.mp4');
} catch (error) {
    if (error.name !== 'AbortError') {
        console.error('Could not load the clip:', error);
    }
}
```

Ändern Sie bei einer Anwendung mit der gebündelten Runtime den Import in:

```javascript
import { Media } from '/wails/runtime.js';
```

Das zweite Argument ist der registrierte Streamname. Das dritte ist ein durch Schrägstriche getrennter Dateipfad relativ zum Stamm seines Dateisystems, etwa `welcome.mp4` oder `tutorials/intro.mp4`. Es ist weder eine URL noch ein Betriebssystempfad.

Das Promise wird erfüllt, sobald die Quelle zugewiesen ist. Anschließend dekodiert der Player sie. Behandeln Sie das `error`-Ereignis des Players, um nicht unterstützte Codecs zu erkennen, und rufen Sie `player.play()` aus einer Benutzerinteraktion auf, wenn Sie die Wiedergabe selbst starten müssen.

## Einen Clip ersetzen oder freigeben

Rufen Sie `Media.SetSource` erneut auf, um den Clip zu wechseln. Dadurch wird ein vorheriger ausstehender Ladevorgang dieses Players abgebrochen, sodass eine langsame Antwort nicht die neueste Auswahl überschreiben kann. Der bisherige Clip bleibt verfügbar, bis der Ersatz erfolgreich geladen wurde. Danach wird seine Blob-URL widerrufen.

Rufen Sie `Media.ClearSource(player)` beim Schließen des Players oder beim Entfernen seiner Komponente auf:

```javascript
Media.ClearSource(player);
```

Dies bricht einen ausstehenden Ladevorgang ab, setzt den Player zurück und gibt die Blob-URL frei. Rufen Sie die Funktion auf, bevor Sie das Element aus dem Dokument entfernen. Eine Framework-Komponente sollte sie in ihrem Unmount- oder Disposal-Hook aufrufen. Das Entfernen des Elements allein gibt die Blob-URL nicht frei. Verwenden Sie diese Hilfsfunktionen durchgängig für die Quelle des Players.

Übergeben Sie bei `<video><source ...></video>` den gewählten Dateinamen an `Media.SetSource(video, "media", name)`. Die Hilfsfunktion setzt das `src` des übergeordneten Players, das Vorrang vor seinen `<source>`-Kindern hat. Nach dem Entfernen dieses Attributs kann der Browser die Kinder wieder berücksichtigen; verwenden Sie einen leeren Player, wenn Sie alle Quellen über diese API verwalten.

Auch nicht in das Dokument eingebundene Audioobjekte funktionieren:

```javascript
const sound = new Audio();
await Media.SetSource(sound, 'media', 'notification.mp3');
sound.addEventListener('ended', () => Media.ClearSource(sound), { once: true });
// Call sound.play() from an appropriate user interaction.
// Also clear it if playback is cancelled or the owning component is disposed.
```

## Downloads begrenzen und Laden abbrechen

Die Standardgrenze beträgt **32 MiB pro Quelle**. Sie können eine kleinere oder größere positive Ganzzahl in Bytes wählen, innerhalb der in Go konfigurierten Grenze:

```javascript
const controller = new AbortController();
const loading = Media.SetSource(player, 'media', 'welcome.mp4', {
    maxBytes: 8 * 1024 * 1024,
    signal: controller.signal,
});

// Call controller.abort() to cancel this load.
await loading;
```

Eine zu große Datei führt zur Ablehnung mit `RangeError`. Der Go-Handler prüft die Dateigröße vor dem Lesen und überträgt höchstens die kleinere der konfigurierten Grenze und der Frontend-Grenze. Auch das Frontend prüft die empfangene Größe und weist unvollständige Übertragungen zurück. Ein Abbruch führt zur Ablehnung mit `AbortError` oder dem an `AbortController.abort(reason)` übergebenen Grund.

Dateien werden in Frames von 64 KiB übertragen. Der Wails-Desktop-Stream-Transport begrenzt Poll-Antworten auf 1 MiB, sodass WebView2 beim Puffern vollständiger Antworten nicht die gesamte Mediendatei in einer einzigen Antwort puffert. Stream-Warteschlangen haben eine eigene begrenzte Rückstausteuerung. Diese Transportgrenzen ändern nichts daran, dass die Player-Hilfsfunktion einen vollständigen Blob verwendet.

**Die gesamte Datei wird vor der Wiedergabe heruntergeladen.** Die Byte-Grenze gilt pro Datei und begrenzt nicht den gesamten Anwendungsspeicher. Mehrere Player, der bisherige Clip während des Ersetzens, der Blob-Aufbau und dekodierte Medien können zusätzlichen Speicher benötigen. Die Hilfsfunktion überträgt beim Aufruf unabhängig von der `preload`-Einstellung des Players. Rufen Sie sie auf, wenn der Benutzer einen Clip zum Laden auswählt.

Für große lokale Dateien ist diese Hilfsfunktion keine Streaming-Lösung. Eine höhere Grenze erhöht auch den Speicherbedarf. Bereits über HTTP/HTTPS bereitgestellte Medien sollten das native Medienladen verwenden, damit der Browser über Bereichsanfragen streamen und die Wiedergabeposition ändern kann.

## Wiedergabeprobleme unter Linux beheben

- Meldet die direkte lokale Wiedergabe **No URI handler implemented for "wails"**, laden Sie den Clip mit `Media.SetSource`. Dies betrifft sowohl den standardmäßigen GTK4-Stack als auch den älteren Stack mit `-tags gtk3`.
- Gelingt das Laden, aber nicht das Dekodieren, prüfen Sie die auf dem Zielsystem installierten GStreamer-Codecs. MP4 benötigt meist Unterstützung für H.264-Video und AAC-Audio; MP3 benötigt einen MP3-Decoder. Testen Sie die ausgelieferten Formate auf den unterstützten Distributionen.
- Schlägt die Übertragung fehl, prüfen Sie den registrierten Streamnamen, den relativen Dateinamen und die Dateisystemberechtigungen. Änderungen während der Übertragung können zu einer unvollständigen Übertragung führen; versuchen Sie es nach Abschluss des Schreibens erneut.
- Verwendet Ihre Anwendung eine Content Security Policy, erlauben Sie den Wails-Asset-Ursprung in `connect-src` und `blob:` in `media-src`. Für eine rein lokale Richtlinie können diese Direktiven `connect-src 'self'; media-src 'self' blob:` lauten. Behalten Sie die übrigen Direktiven bei.
- Überschreitet eine Datei die Grenze, wählen Sie einen kürzeren oder kleineren Clip oder setzen Sie eine explizite Grenze passend zum Speicherbudget Ihrer Anwendung.

Starten Sie das [audio-video-Beispiel](https://github.com/wailsapp/wails/tree/master/v3/examples/audio-video) mit `go run .`, um die enthaltenen MP3- und MP4-Beispiele auf Ihrem Rechner zu prüfen.
