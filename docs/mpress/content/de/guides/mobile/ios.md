---
title: "iOS"
description: "Wails-Anwendungen unter iOS erstellen und ausführen — Einrichtung, Simulator, Geräte-Builds, Konfiguration und native Funktionen"
slug: "guides/mobile/ios"
sourcePath: "guides/mobile/ios.md"
---

@note{type="caution" title="Experimentelle Funktion"}
Die iOS-Unterstützung ist experimentell und kann sich in zukünftigen Versionen ändern.

@end

@note{type="tip"}
Noch keine Erfahrung mit Wails für Mobilgeräte? Beginnen Sie mit [Ihre erste mobile App →](/guides/mobile/first-mobile-app/), um eine schrittweise Anleitung zu erhalten, und kehren Sie anschließend zur vollständigen Referenz hierher zurück.

@end

Wails-v3-Apps werden unter iOS als vollständig native Apps ausgeführt — und das Beste daran: Sie funktionieren *genau* wie die Desktopversion. Dasselbe Go-Backend, dasselbe Frontend und dieselbe `@wailsio/runtime`: Dienstbindungen, Ereignisse, Dialoge und die Zwischenablage verhalten sich alle identisch, mit **null** mobilspezifischen Anpassungen. Es gibt keine separate Codebasis für Mobilgeräte, keine Portierungsschicht und keine spezielle API, die Sie erlernen müssen — Ihre vorhandene Wails-App läuft einfach unter iOS. Die Portierung ist wirklich nahtlos: Übernehmen Sie Ihre App unverändert und veröffentlichen Sie sie.

Dasselbe `main.go` wird sowohl für den Desktop als auch für iOS erstellt; iOS-spezifische Anpassungen werden über `application.Options.IOS` konfiguriert.

## Voraussetzungen

- macOS mit installiertem **vollständigem Xcode** (die Befehlszeilenwerkzeuge allein reichen nicht aus) — `wails3 doctor` zeigt die gefundenen iOS-SDKs an
- Go 1.25+ und npm

## Simulator

Aus Ihrem Projektverzeichnis:

```bash
wails3 task ios:run
```

Dadurch wird Ihre App erstellt, bei Bedarf ein Simulator gestartet und die App darin ausgeführt.

Nützliche Begleitbefehle:

```bash
wails3 task ios:logs:dev    # stream the app's logs from the simulator
wails3 task ios:xcode       # open the generated Xcode project
```

In Debug-Builds kann die WebView über das Menü „Entwickler“ von Safari untersucht werden.

## Paketierung

```bash
wails3 task ios:package             # production .app for the simulator
wails3 task ios:deploy-simulator    # install + launch it
```

Dies sind optimierte Produktions-Builds mit entfernten Debug- und Symbolinformationen.

## Geräte-Builds

```bash
wails3 task ios:package IOS_PLATFORM=device \
    CODESIGN_IDENTITY="Apple Development: You (TEAMID)" \
    PROVISIONING_PROFILE=path/to/profile.mobileprovision

wails3 task ios:deploy-device [DEVICE_ID=<udid>]     # install + launch on a device
wails3 task ios:package:ipa IOS_PLATFORM=device ...  # distribution .ipa
```

`IOS_PLATFORM=device` erstellt einen Build für ein physisches Gerät. Berechtigungen stammen aus `build/ios/entitlements.plist` und gelten nur für Geräte-Builds — fügen Sie die von Ihrer App benötigten Funktionsschlüssel hinzu.

@note{type="tip"}
Öffnen Sie für automatisch verwaltete Signierung, Bereitstellung und App-Store-Archive das generierte Xcode-Projekt mit `wails3 task ios:xcode` und erstellen Sie den Build stattdessen in Xcode.

@end

## Konfiguration

`build/config.yml`:

```yaml
ios:
  bundleID: com.example.myapp
  displayName: My App
  version: 1.0.0
  minIOSVersion: "15.0"
```

Zu den Startoptionen (`application.Options.IOS`) gehören `DisableScroll`, `DisableBounce`, `DisableScrollIndicators`, `DisableInputAccessoryView`, `EnableBackForwardNavigationGestures`, `DisableLinkPreview`, `EnableInlineMediaPlayback`, `EnableAutoplayWithoutUserAction`, `DisableInspectable`, `UserAgent`, `ApplicationNameForUserAgent`, `BackgroundColour` sowie native untere Registerkarten über `EnableNativeTabs` + `NativeTabsItems`.

## Native Funktionen

iOS-spezifische Funktionen sind über `application.IOS` verfügbar. Rufen Sie sie in Go innerhalb einer `//go:build ios`-Datei auf, damit Ihr gemeinsam genutzter Code plattformunabhängig bleibt. Android stellt dieselben Funktionen über `application.Android` bereit.

Einmalige Aktionen kehren sofort zurück:

```go
//go:build ios

application.IOS.Haptic("impact-medium") // impact-light|impact-medium|impact-heavy|success|warning|error|selection
application.IOS.Share(`{"text":"Hi","url":"https://wails.io"}`)
application.IOS.SetKeepAwake(true)
application.IOS.PostNotification(`{"title":"Done","body":"Build finished","delay":2}`)
application.IOS.SecureSet("token", "abc") // stored securely
```

Abfragehilfen geben ihre Ergebnisse als JSON zurück — `SafeAreaJSON()`, `AppInfoJSON()`, `PowerJSON()`, `NetworkJSON()`, `StorageJSON()`, `GetOrientation()`, `GetBrightness()`. `StoragePath()` gibt den absoluten Pfad zum Application-Support-Verzeichnis der App zurück — ein geeigneter Speicherort für Datenbanken und andere persistente Dateien (das iOS-Pendant zu `getFilesDir()` unter Android). Das Verzeichnis wird beim ersten Zugriff erstellt; `StoragePath()` gibt eine leere Zeichenfolge zurück, wenn es nicht erstellt werden kann. Prüfen Sie daher vor der Verwendung auf `""`.

### Ereignisse

Alles, was erst später abgeschlossen wird — eine Berechtigungsabfrage, ein Sensorstream oder eine Kameraaufnahme — liefert sein Ergebnis als **Ereignis** statt als Rückgabewert. Sie können es in Go oder im Frontend empfangen. Namen erhalten für mit Android gemeinsam genutzte Funktionen das Präfix `common:` und für ausschließlich unter iOS verfügbare Funktionen das Präfix `ios:`.

```go
// Go
app.Event.On("common:location", func(e *application.CustomEvent) {
    // e.Data -> {"lat":..,"lng":..,"accuracy":..} or {"error":..}
})
```

```js
// frontend
import { Events } from "@wailsio/runtime";
Events.On("common:notification", (e) => { /* {ok, scheduled, presented, tapped, error} */ });
```

| Ereignis | Ausgelöst durch | Nutzdaten |
| --- | --- | --- |
| `common:biometric` | `BiometricAuthenticate(reason)` | `{ok, error}` |
| `common:location` | `GetLocation()` | `{lat, lng, accuracy}` / `{error}` |
| `common:motion` | `SetMotion(true)` | `{x, y, z}` |
| `common:proximity` | `SetProximity(true)` | `{near}` |
| `common:keyboard` | `SetKeyboardWatch(true)` | `{visible, height}` |
| `common:torch` | `SetTorch(bool)` | `{on, available}` |
| `common:notification` | `PostNotification(json)` | `{ok, scheduled, presented, tapped, error}` |
| `common:capture` | `CapturePhoto()` / `CaptureVideo()` | `{type, path, size, thumb}` |
| `common:screenCapture` | `SetScreenProtect(true)` | `{screenshot, recording}` |
| `ios:backgroundTask` | `BeginBackgroundTask(seconds)` | `{message, granted}` |

Das umfassende Beispiel unter `v3/examples/mobile` verbindet jede der oben genannten Funktionen durchgängig.

## WebView-Steuerung

Einige Verhaltensweisen der WebView lassen sich auch zur Laufzeit aus Go ändern:

```go
application.IOS.SetScrollEnabled(false)
application.IOS.SetBounceEnabled(false)
application.IOS.SetScrollIndicatorsEnabled(false)
application.IOS.SetBackForwardGesturesEnabled(true)
application.IOS.SetLinkPreviewEnabled(false)
application.IOS.SetInspectableEnabled(true)
application.IOS.SetCustomUserAgent("MyApp/1.0")
```

Das gebündelte `@wailsio/runtime` stellt außerdem einen kleinen iOS-Namensraum für das **Frontend** bereit:

```js
import { IOS } from "@wailsio/runtime";
await IOS.Haptics.Impact("medium"); // light|medium|heavy|soft|rigid
const info = await IOS.Device.Info();
```

Die Auswahl nativer Tabs am unteren Rand wird als `nativeTabSelected`-Ereignis auf `window` empfangen.

## Unterstützungsstatus

| Bereich | Status |
| --- | --- |
| Frontend-Rendering und Assets | ✅ |
| Service-Bindings, Ereignisse (in beide Richtungen) | ✅ |
| Meldungsdialoge | ✅ |
| Dialoge zum Öffnen einer Datei, mehrerer Dateien oder eines Verzeichnisses | ✅ Als Sandbox-Kopien importiert |
| Dialoge zum Speichern von Dateien | ❌ Schreiben Sie stattdessen in die App-Sandbox |
| Zwischenablage | ✅ |
| Screens-API | ✅ Enthält den Arbeitsbereich innerhalb der sicheren Bereiche |
| Lebenszyklusereignisse | ✅ |
| Fenstergeometrie, Menüs, System-Tray | Unter iOS ohne Wirkung |
| Mehrere Fenster | Nur das erste Fenster wird angezeigt |

## Hinweise zur Portierung

- Desktop-Code lässt sich unverändert für iOS erstellen – Aufrufe für Fenster, Menüs und den System-Tray bleiben dort einfach ohne Wirkung.
- Ersetzen Sie Dialoge zum Speichern von Dateien durch das Schreiben in die App-Sandbox mit anschließendem Teilen.
- Gestalten Sie das Frontend responsiv; sichere Bereiche werden automatisch berücksichtigt.
