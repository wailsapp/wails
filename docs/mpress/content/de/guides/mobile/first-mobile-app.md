---
title: "Ihre erste mobile App"
description: "Führen Sie Ihre Wails-App in wenigen Minuten im iOS-Simulator oder Android-Emulator aus"
slug: "guides/mobile/first-mobile-app"
sourcePath: "guides/mobile/first-mobile-app.md"
---

In dieser Anleitung führen Sie eine standardmäßige Wails-Desktop-App im iOS-Simulator oder Android-Emulator aus. **Sie müssen Ihren Go-Code nicht ändern.** Derselbe `main.go` wird für alle Zielplattformen erstellt.

**Benötigte Zeit:** 15–30 Minuten (der Großteil davon entfällt beim ersten Durchlauf auf die Installation der Toolchain)

## Mit einem Desktop-Projekt beginnen

Falls Sie noch keines haben, erstellen Sie ein neues Projekt:

```bash
wails3 init -n mymobileapp
cd mymobileapp
```

Prüfen Sie zunächst, ob die Desktop-App funktioniert:

```bash
wails3 dev
```

Sobald sie geöffnet ist, beenden Sie sie und fahren Sie fort. Alles, was auf dem Desktop läuft, läuft auch auf Mobilgeräten — für diese Anleitung müssen Sie weder `main.go` noch anderen Go-Code ändern.

---

## Plattform auswählen

@tabs{sync-key="mobile-platform"}
[iOS-Simulator]
### Voraussetzungen

- **macOS** (iOS-Builds sind nur unter macOS möglich)
- **Vollständiges Xcode** — nicht nur die Befehlszeilenwerkzeuge. Installieren Sie es aus dem App Store und führen Sie anschließend Folgendes aus:
  ```bash
  sudo xcode-select -s /Applications/Xcode.app/Contents/Developer
  sudo xcodebuild -license accept
  ```


- **Go 1.25+** und **npm** (bereits installiert, wenn Sie `wails3 init` ausgeführt haben)

Führen Sie zur Überprüfung `wails3 doctor` aus — der Befehl listet alle gefundenen iOS-SDKs auf.

### Im Simulator ausführen

@steps
### App starten
```bash
wails3 task ios:run
```

Das ist alles — dadurch wird Ihre App erstellt, bei Bedarf ein Simulator gestartet und die App darin geöffnet.

@note{type="tip"}
Der erste Durchlauf dauert einige Minuten, da das Wails-Framework für iOS kompiliert und zwischengespeichert wird. Alle weiteren Durchläufe sind wesentlich schneller.

@end

Nach dem Start läuft Ihre unveränderte Desktop-App im iOS-Simulator — mit demselben `main.go` und demselben Frontend:

![Eine standardmäßige Wails-App im iOS-Simulator](/assets/ios-simulator-first-app.png)

### Logs fortlaufend anzeigen
In einem separaten Terminal:

```bash
wails3 task ios:logs:dev
```

Damit wird das Simulator-Log fortlaufend und nach Ihrer App gefiltert angezeigt. Die Ausgabe von `fmt.Println` und `log.Println` erscheint hier.

### WebView untersuchen
In Safari: **Entwickler → Simulator → Ihre App**. Der vollständige Web Inspector ist verfügbar — Konsole, Debugger, Netzwerkbereich und alle weiteren Funktionen.

### Änderung vornehmen
Bearbeiten Sie eine beliebige Frontend-Datei (`frontend/src/main.js`, `index.html` usw.) und führen Sie `wails3 task ios:run` erneut aus. Wails erstellt das Frontend neu und startet die App erneut.

Führen Sie nach Änderungen am Go-Code ebenfalls `wails3 task ios:run` erneut aus. Go kompiliert inkrementell neu, sodass nur geänderte Pakete neu erstellt werden.

@end

### In Xcode öffnen (optional)

```bash
wails3 task ios:xcode
```

Dadurch wird `build/ios/` in Xcode geöffnet. Sie können Xcode für die Bereitstellung auf Geräten, erweitertes Profiling oder die Verwaltung von Bereitstellungsprofilen verwenden. Wails generiert das Xcode-Projekt bei jedem Build neu; ändern Sie die generierten Dateien daher nicht direkt.

[Android-Emulator]
### Voraussetzungen

Sie benötigen das **Android SDK**, das **NDK** und ein **JDK**. Am einfachsten geht dies mit Android Studio oder alternativ mit den Befehlszeilenwerkzeugen:

@steps
### Android-Befehlszeilenwerkzeuge installieren
Laden Sie sie von [developer.android.com/studio#command-line-tools-only](https://developer.android.com/studio#command-line-tools-only) herunter und entpacken Sie sie nach `~/android-sdk/cmdline-tools/latest/`.

### SDK-Komponenten installieren
```bash
sdkmanager "platform-tools" \
           "platforms;android-35" \
           "build-tools;35.0.0" \
           "ndk;26.3.11579264" \
           "emulator" \
           "system-images;android-35;google_apis;arm64-v8a"
```

### Emulator erstellen
```bash
avdmanager create avd \
  --name wails \
  --package "system-images;android-35;google_apis;arm64-v8a" \
  --device pixel_7
```

### Umgebungsvariablen festlegen
Fügen Sie Folgendes zu `~/.zshrc` oder `~/.bashrc` hinzu:

```bash
export ANDROID_HOME=~/android-sdk
export ANDROID_SDK_ROOT=~/android-sdk
export PATH=$PATH:$ANDROID_HOME/platform-tools:$ANDROID_HOME/cmdline-tools/latest/bin
```

Neu laden: `source ~/.zshrc`

### JDK installieren
```bash
# macOS
brew install openjdk@21
export JAVA_HOME=$(brew --prefix openjdk@21)

# Ubuntu/Debian
sudo apt install openjdk-21-jdk
export JAVA_HOME=/usr/lib/jvm/java-21-openjdk-amd64

# Windows (scoop)
scoop install openjdk21
```

@end

Führen Sie `wails3 doctor` aus, um zu bestätigen, dass alle Komponenten gefunden werden.

### Im Emulator ausführen

@steps
### App starten
```bash
wails3 task android:run
```

Beim ersten Durchlauf geschieht Folgendes:

- Der Emulator wird gestartet, falls noch keiner läuft
- Bindings werden generiert und das Frontend wird erstellt
- Ihr Go-Code wird mit dem NDK-Cross-Compiler zu `libwails.so` kompiliert
- Mit Gradle wird eine Debug-APK erstellt
- Sie wird im Emulator installiert und gestartet

@note{type="tip"}
Beim ersten Build wird Gradle heruntergeladen und die NDK-Toolchain kompiliert — rechnen Sie mit 5–10 Minuten. Nachfolgende Builds sind inkrementell und dauern weniger als eine Minute.

@end

### Logs fortlaufend anzeigen
In einem separaten Terminal:

```bash
wails3 task android:logs
```

Damit wird `adb logcat` ausgeführt und nach Ihrer App gefiltert. Die Ausgabe von `fmt.Println` erscheint hier.

### WebView untersuchen
Öffnen Sie Chrome und rufen Sie `chrome://inspect` auf. Die WebView Ihrer App wird unter **Remote-Ziel** angezeigt — klicken Sie auf **Untersuchen**, um die DevTools zu öffnen.

### Änderung vornehmen
Bearbeiten Sie eine beliebige Datei und führen Sie `wails3 task android:run` erneut aus. Dank des inkrementellen Builds von Gradle wird nur geänderter Code neu kompiliert.

@end

@end

---

## Was ist geschehen?

Ihr `main.go` wurde überhaupt nicht geändert. Wails hat alles übernommen:

- **Build-System** — `Taskfile.yml` in Ihrem Projekt enthält die Tasks `ios:*` und `android:*`, die die plattformspezifische Toolchain steuern.
- **Go-Cross-Kompilierung** — `GOOS=ios` oder `GOOS=android` mit dem passenden `GOARCH` und Sysroot.
- **Nativer Host** — ein generiertes Xcode-Projekt (iOS) oder Gradle-Projekt (Android), das Ihren kompilierten Go-Code einbettet und die WebView hostet.
- **Bereitstellung von Assets** — Ihr `frontend/dist/` wird in die Go-Binärdatei eingebettet und innerhalb des Prozesses bereitgestellt. Ein localhost-Server ist nicht erforderlich.

---

## Ihre App für Mobilgeräte optimieren

Ihre App funktioniert bereits, sieht auf einem Smartphone-Bildschirm aber wie eine Desktop-App aus. Einige kleine Änderungen bewirken viel.

### Responsives CSS

Bildschirme von Mobilgeräten sind schmaler und verwenden andere Eingabemuster. In `frontend/public/style.css` (oder einer entsprechenden Datei):

```css
/* Prevent horizontal scrolling */
body {
  overflow-x: hidden;
}

/* Touch-friendly tap targets */
button {
  min-height: 44px;
  min-width: 44px;
}

/* Respect the iOS safe area (notch, home indicator) */
body {
  padding-top: env(safe-area-inset-top);
  padding-bottom: env(safe-area-inset-bottom);
  padding-left: env(safe-area-inset-left);
  padding-right: env(safe-area-inset-right);
}
```

### Plattform in Go erkennen

Verwenden Sie Build-Tags, um plattformspezifisches Verhalten hinzuzufügen, ohne den gemeinsam genutzten Code unübersichtlich zu machen.

Erstellen Sie `mobile_ios.go` für Code, der ausschließlich unter iOS verwendet wird:

```go {title="mobile_ios.go"}
//go:build ios

package main

import "github.com/wailsapp/wails/v3/pkg/application"

func platformOptions() application.IOSOptions {
    return application.IOSOptions{
        DisableBounce: true,
    }
}
```

Erstellen Sie `mobile_android.go` für Code, der ausschließlich unter Android verwendet wird:

```go {title="mobile_android.go"}
//go:build android

package main

import "github.com/wailsapp/wails/v3/pkg/application"

func platformOptions() application.AndroidOptions {
    return application.AndroidOptions{}
}
```

Erstellen Sie `mobile_desktop.go` als Stub, damit der gemeinsam genutzte Code auch für Desktop-Plattformen kompiliert wird:

```go {title="mobile_desktop.go"}
//go:build !ios && !android

package main

type mobileOptions struct{}

func platformOptions() mobileOptions { return mobileOptions{} }
```

### Plattform in JavaScript erkennen und ausschließlich mobile Benutzeroberflächen gezielt aktivieren

Die Runtime-Objekte `IOS.*` und `Android.*` sind jeweils nur auf der entsprechenden Plattform vorhanden. Wenn Sie sie auf einer Desktop-Plattform aufrufen, wird eine Exception ausgelöst. Das richtige, auch vom Kitchen Sink verwendete Muster besteht darin, die Plattform einmal zu erkennen und ausschließlich für Mobilgeräte bestimmte Steuerelemente vollständig auszublenden:

```javascript
// Detect platform from the bridge the host injects into the WebView
const platform = (() => {
  if (typeof window.wails?.platform === 'function') return window.wails.platform(); // Android
  if (window.webkit?.messageHandlers?.external) return 'ios';
  return 'desktop';
})();

const isIOS     = platform === 'ios';
const isAndroid = platform === 'android';
const isMobile  = isIOS || isAndroid;

// Hide any element marked as mobile-only
document.querySelectorAll('.mobile-only').forEach(el => {
  el.style.display = isMobile ? '' : 'none';
});
```

Anschließend in Ihrem HTML:

```html
<section class="mobile-only">
  <button id="btnHaptic">Haptic feedback</button>
</section>
```

So werden ausschließlich für Mobilgeräte bestimmte Schaltflächen auf Desktop-Plattformen nie gerendert, und Sie müssen nicht jeden einzelnen Aufruf durch eine `if (isMobile)`-Prüfung absichern.

Kombinieren Sie dies auf der Go-Seite mit einem Build-Tag-Stub, damit die Event-Handler nur auf den Plattformen registriert werden, die sie benötigen:

```go {title="native_desktop.go"}
//go:build !ios && !android

package main

import "github.com/wailsapp/wails/v3/pkg/application"

// No-op on desktop — mobile tabs are hidden in the frontend so these
// events are never emitted.
func registerNativeFeatures(app *application.App) {}
```

```go {title="native_ios.go"}
//go:build ios

package main

import "github.com/wailsapp/wails/v3/pkg/application"

func registerNativeFeatures(app *application.App) {
    app.Event.On("common:haptic", func(e *application.CustomEvent) {
        // only compiled and called on iOS
        application.IOS.Haptic("medium")
    })
    // ... other handlers
}
```

Genau dieses Muster verwendet der [Kitchen Sink](https://github.com/wailsapp/wails/tree/master/v3/examples/mobile) — siehe `native_features_stub.go`, `native_features_ios.go` und `native_features_android.go`.

@note{type="note" title="API für native Funktionen und Benennung von Events"}
Zwei Konventionen sollten Sie kennen:

- **Native Funktionen auf der Go-Seite verwenden Plattform-Manager.** Rufen Sie sie über die Singletons `application.IOS.*` und `application.Android.*` auf — beispielsweise `application.IOS.Haptic("medium")` oder `application.Android.Share(payload)`. Jeder Manager ist nur auf seiner eigenen Plattform vorhanden; deshalb befinden sich seine Aufrufe in `//go:build ios`- bzw. `//go:build android`-Dateien.
- **Events verwenden Namensräume entsprechend ihrer Reichweite.** Alles, was beide Plattformen verstehen, verwendet das Präfix `common:*` (`common:haptic`, `common:location`, …); Events, die nur eine Plattform erzeugen oder verarbeiten kann, verwenden `ios:*` oder `android:*` (beispielsweise `ios:backgroundTask`, `android:foregroundService`). Da fast alle mobilen Funktionen gemeinsam genutzt werden, benötigt Ihr Frontend unter `common:*` nur einen Listener pro Event.

@end

### Haptisches Feedback hinzufügen (iOS)

```javascript
import { IOS } from '@wailsio/runtime';

async function onButtonTap() {
  if (isIOS) {
    await IOS.Haptics.Impact({ style: 'medium' });
  }
  // ... rest of your handler
}
```

### Vibration hinzufügen (Android)

```javascript
import { Android } from '@wailsio/runtime';

async function onButtonTap() {
  if (isAndroid) {
    await Android.Haptics.Vibrate(50); // 50ms
  }
}
```

---

## Für den Produktivbetrieb bauen

@tabs{sync-key="mobile-platform"}
[iOS]
**Simulator-Build** (zum Testen im Simulator; keine Signierung erforderlich):

```bash
wails3 task ios:package
wails3 task ios:deploy-simulator
```

**Geräte-Build** (erfordert eine Signierungsidentität und ein Bereitstellungsprofil):

```bash
wails3 task ios:package \
  IOS_PLATFORM=device \
  CODESIGN_IDENTITY="Apple Development: You (TEAMID)" \
  PROVISIONING_PROFILE=path/to/profile.mobileprovision

wails3 task ios:deploy-device   # installs via xcrun devicectl
```

**IPA für die Verteilung** (für den App Store oder TestFlight):

```bash
wails3 task ios:package:ipa IOS_PLATFORM=device \
  CODESIGN_IDENTITY="..." \
  PROVISIONING_PROFILE=path/to/distribution.mobileprovision
```

@note{type="tip"}
Verwenden Sie für Uploads zu App Store Connect `wails3 task ios:xcode` und lassen Sie Xcode die Signierung und Archivierung verwalten — die komplexe Handhabung von Zertifikaten, Profilen und Notarisierung übernimmt Xcode automatisch.

@end

[Android]
**Debug-APK** (mit dem Android-Debug-Keystore signiert und direkt installierbar):

```bash
wails3 task android:package
wails3 task android:deploy-emulator
```

**Release-APK** (mit Ihrem eigenen Keystore signiert):

```bash
ANDROID_KEYSTORE_FILE=/path/to/release.jks \
ANDROID_KEYSTORE_PASSWORD=yourpassword \
ANDROID_KEY_ALIAS=youralias \
ANDROID_KEY_PASSWORD=yourkeypassword \
  wails3 task android:package
```

**Universelle APK** (arm64 und x86_64 in einer Datei):

```bash
wails3 task android:package:fat
```

@note{type="tip"}
Erstellen Sie für Uploads zum Play Store anstelle einer APK ein `.aab` (Android App Bundle) — öffnen Sie `build/android/` in Android Studio und wählen Sie **Build → Generate Signed Bundle / APK**.

@end

@end

---

## Fehlerbehebung

### `wails3 task ios:run` schlägt mit "no iOS SDKs found" fehl

Die vollständige Xcode-Version muss installiert und ausgewählt sein:

```bash
sudo xcode-select -s /Applications/Xcode.app/Contents/Developer
xcode-select -p  # should print the Xcode path
```

#### `wails3 task android:run` schlägt mit "SDK not found" fehl

Stellen Sie sicher, dass `ANDROID_HOME` festgelegt und exportiert ist. Überprüfen Sie dies mit:

```bash
echo $ANDROID_HOME
ls $ANDROID_HOME/platform-tools/adb
```

#### Simulator startet nicht

Listen Sie die verfügbaren Simulatoren auf und starten Sie einen davon manuell:

```bash
xcrun simctl list devices available
xcrun simctl boot "iPhone 16"
```

#### `chrome://inspect` zeigt keine Ziele an

Die WebView muss sich im Debug-Modus befinden (der Standard für `android:run`). Stellen Sie sicher, dass Sie einen Debug-Build und keinen Produktions-Build ausführen. Prüfen Sie außerdem, ob `adb devices` den Emulator als verbunden anzeigt.

#### Safe-Area-Abstände werden nicht angewendet

Stellen Sie sicher, dass Ihr HTML das Viewport-Meta-Tag enthält:

```html
<meta name="viewport" content="width=device-width, initial-scale=1, viewport-fit=cover">
```

---

## Kitchen Sink erkunden

Sobald Ihre erste App läuft, ist das Beispiel **Kitchen Sink** der schnellste Weg, weitere Möglichkeiten kennenzulernen. Es ist eine vollständige Wails-App, die mit einer einzigen Codebasis unter iOS, Android und auf Desktop-Systemen läuft und unter anderem Haptik, Geolokalisierung, Biometrie, lokale Benachrichtigungen und sichere Speicherung abdeckt:

```bash
git clone https://github.com/wailsapp/wails.git
cd wails/v3/examples/mobile

wails3 task ios:run        # iOS Simulator
wails3 task android:run    # Android Emulator
wails3 task run            # Desktop
```

Sehen Sie sich den Quellcode unter [`v3/examples/mobile`](https://github.com/wailsapp/wails/tree/master/v3/examples/mobile) an. Insbesondere die Dateien `native_features_ios.go` und `native_features_android.go` eignen sich als Ausgangspunkte zum Kopieren und Einfügen für plattformspezifische Funktionen.

## Wie geht es weiter?

@cards{cols="2"}
iOS-Leitfaden
Vollständige Referenz: Konfigurationsoptionen, native Tabs, WKWebView-Umschalter, Builds für Geräte und Signierung.

[iOS-Leitfaden →](/guides/mobile/ios/)

---
Android-Leitfaden
Vollständige Referenz: Konfiguration, Toast-Benachrichtigungen, Paketierung für den Play Store und NDK-Details.

[Android-Leitfaden →](/guides/mobile/android/)

---
📖 Kitchen-Sink-Quellcode
Haptik, Geolokalisierung, Biometrie, Benachrichtigungen und sichere Speicherung – alles in einer einzigen ausführbaren App.

[Auf GitHub ansehen →](https://github.com/wailsapp/wails/tree/master/v3/examples/mobile)

@end
