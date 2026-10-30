---
title: "Ihre erste mobile App"
description: "Führen Sie Ihre Wails-App in wenigen Minuten im iOS-Simulator oder Android-Emulator aus"
slug: "guides/mobile/first-mobile-app"
sourcePath: "guides/mobile/first-mobile-app.md"
---

In diesem Leitfaden führen Sie eine standardmäßige Wails-Desktop-App im iOS-Simulator oder Android-Emulator aus. **Sie müssen Ihren Go-Code nicht ändern.** Dieselbe `main.go` wird für alle Zielplattformen gebaut.

**Benötigte Zeit:** 15–30 Minuten (größtenteils für die Installation der Toolchain beim ersten Durchlauf)

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

Sobald sie geöffnet ist, beenden Sie sie und fahren Sie fort. Alles, was auf dem Desktop läuft, läuft auch auf Mobilgeräten – für diesen Leitfaden müssen Sie weder `main.go` noch anderen Go-Code bearbeiten.

---

## Plattform auswählen

@tabs{sync-key="mobile-platform"}
[iOS-Simulator]
### Voraussetzungen

- **macOS** (iOS-Builds sind nur unter macOS möglich)
- **Vollständiges Xcode** – nicht nur die Befehlszeilentools. Installieren Sie es aus dem App Store und führen Sie anschließend Folgendes aus:
  ```bash
  sudo xcode-select -s /Applications/Xcode.app/Contents/Developer
  sudo xcodebuild -license accept
  ```


- **Go 1.26+** und **npm** (bereits installiert, wenn Sie `wails3 init` ausgeführt haben)

Führen Sie zur Überprüfung `wails3 doctor` aus. Der Befehl listet alle gefundenen iOS-SDKs auf.

### Im Simulator ausführen

@steps
### App starten
```bash
wails3 task ios:run
```

Das ist alles: Ihre App wird gebaut, bei Bedarf ein Simulator gestartet und die App darin geöffnet.

@note{type="tip"}
Der erste Durchlauf dauert einige Minuten, da das Wails-Framework für iOS kompiliert und zwischengespeichert wird. Alle weiteren Durchläufe sind deutlich schneller.

@end

Nach dem Start läuft Ihre unveränderte Desktop-App im iOS-Simulator – mit derselben `main.go` und demselben Frontend:

![Eine standardmäßige Wails-App im iOS-Simulator](/assets/ios-simulator-first-app.png)

### Logs streamen
In einem separaten Terminal:

```bash
wails3 task ios:logs:dev
```

Damit verfolgen Sie das auf Ihre App gefilterte Simulator-Log fortlaufend. Die Ausgaben von `fmt.Println` und `log.Println` erscheinen hier.

### WebView untersuchen
Öffnen Sie in Safari **Entwickler → Simulator → Ihre App**. Der vollständige Web Inspector ist verfügbar – einschließlich Konsole, Debugger, Netzwerkbereich und aller weiteren Funktionen.

### Änderung vornehmen
Bearbeiten Sie eine beliebige Frontend-Datei (`frontend/src/main.js`, `index.html` usw.) und führen Sie `wails3 task ios:run` erneut aus. Wails baut das Frontend neu und startet die App erneut.

Führen Sie nach Änderungen an Go ebenfalls `wails3 task ios:run` erneut aus. Go kompiliert inkrementell, sodass nur geänderte Pakete neu gebaut werden.

@end

### In Xcode öffnen (optional)

```bash
wails3 task ios:xcode
```

Damit wird `build/ios/` in Xcode geöffnet. Sie können Xcode für die Bereitstellung auf Geräten, erweitertes Profiling oder die Verwaltung von Provisioning-Profilen verwenden. Wails generiert das Xcode-Projekt bei jedem Build neu; bearbeiten Sie die generierten Dateien daher nicht direkt.

[Android-Emulator]
### Voraussetzungen

Sie benötigen das **Android SDK**, das **NDK** und ein **JDK**. Am einfachsten geht dies mit Android Studio oder alternativ mit den Befehlszeilentools:

@steps
### Android-Befehlszeilentools installieren
Laden Sie die Tools von [developer.android.com/studio#command-line-tools-only](https://developer.android.com/studio#command-line-tools-only) herunter und entpacken Sie sie nach `~/android-sdk/cmdline-tools/latest/`.

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

Laden Sie die Konfiguration neu: `source ~/.zshrc`

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

Führen Sie `wails3 doctor` aus, um zu prüfen, ob alles gefunden wird.

### Im Emulator ausführen

@steps
### App starten
```bash
wails3 task android:run
```

Beim ersten Durchlauf geschieht Folgendes:

- Der Emulator wird gestartet, falls noch keiner läuft
- Bindings werden generiert und das Frontend wird gebaut
- Ihr Go-Code wird mit dem NDK-Cross-Compiler nach `libwails.so` kompiliert
- Mit Gradle wird eine Debug-APK erstellt
- Die APK wird im Emulator installiert und gestartet

@note{type="tip"}
Beim ersten Build wird Gradle heruntergeladen und die NDK-Toolchain kompiliert – rechnen Sie mit 5–10 Minuten. Nachfolgende Builds erfolgen inkrementell und dauern weniger als eine Minute.

@end

### Logs streamen
In einem separaten Terminal:

```bash
wails3 task android:logs
```

Damit wird `adb logcat` mit einem Filter für Ihre App ausgeführt. Die Ausgabe von `fmt.Println` erscheint hier.

### WebView untersuchen
Öffnen Sie Chrome und rufen Sie `chrome://inspect` auf. Die WebView Ihrer App erscheint unter **Remote Target**. Klicken Sie auf **inspect**, um die DevTools zu öffnen.

### Änderung vornehmen
Bearbeiten Sie eine beliebige Datei und führen Sie `wails3 task android:run` erneut aus. Dank des inkrementellen Gradle-Builds wird nur geänderter Code neu kompiliert.

@end

@end

---

## Was dabei geschehen ist

Ihre `main.go` wurde überhaupt nicht geändert. Wails hat alles übernommen:

- **Build-System** – `Taskfile.yml` in Ihrem Projekt enthält die Tasks `ios:*` und `android:*`, die die plattformspezifische Toolchain steuern.
- **Go-Cross-Kompilierung** – `GOOS=ios` oder `GOOS=android` mit der passenden `GOARCH` und Sysroot.
- **Nativer Host** – ein generiertes Xcode-Projekt (iOS) oder Gradle-Projekt (Android), das Ihren kompilierten Go-Code einbettet und die WebView hostet.
- **Bereitstellung von Assets** – Ihre `frontend/dist/` wird in die Go-Binärdatei eingebettet und innerhalb desselben Prozesses bereitgestellt. Ein localhost-Server ist nicht erforderlich.

---

## App für Mobilgeräte optimieren

Ihre App funktioniert bereits, sieht auf einem Smartphone-Display jedoch noch wie eine Desktop-App aus. Wenige kleine Änderungen bewirken einen großen Unterschied.

### Responsives CSS

Mobile Displays sind schmaler und verwenden andere Eingabemuster. Fügen Sie in `frontend/public/style.css` (oder einer entsprechenden Datei) Folgendes hinzu:

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

Verwenden Sie Build-Tags, um plattformspezifisches Verhalten hinzuzufügen, ohne den gemeinsamen Code zu überladen.

Erstellen Sie `mobile_ios.go` für Code, der nur unter iOS verwendet wird:

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

Erstellen Sie `mobile_android.go` für Code, der nur unter Android verwendet wird:

```go {title="mobile_android.go"}
//go:build android

package main

import "github.com/wailsapp/wails/v3/pkg/application"

func platformOptions() application.AndroidOptions {
    return application.AndroidOptions{}
}
```

Erstellen Sie `mobile_desktop.go` als Stub, damit der gemeinsame Code auch für Desktop-Plattformen kompiliert wird:

```go {title="mobile_desktop.go"}
//go:build !ios && !android

package main

type mobileOptions struct{}

func platformOptions() mobileOptions { return mobileOptions{} }
```

### Plattform in JavaScript erkennen und ausschließlich mobile Benutzeroberfläche bedingt aktivieren

Die Runtime-Objekte `IOS.*` und `Android.*` sind nur auf ihrer jeweiligen Plattform vorhanden. Ein Aufruf auf einer Desktop-Plattform löst einen Fehler aus. Das richtige, auch im Kitchen Sink verwendete Muster besteht darin, die Plattform einmal zu erkennen und ausschließlich mobile Bedienelemente vollständig auszublenden:

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

Dann in deinem HTML:

```html
<section class="mobile-only">
  <button id="btnHaptic">Haptic feedback</button>
</section>
```

So werden Schaltflächen, die nur für Mobilgeräte bestimmt sind, auf dem Desktop nie gerendert, und du musst nicht jeden einzelnen Aufruf mit einer `if (isMobile)`-Prüfung absichern.

Kombiniere dies auf der Go-Seite mit einem Build-Tag-Stub, damit die Event-Handler nur auf den Plattformen registriert werden, die sie benötigen:

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

Genau dieses Muster verwendet das [Kitchen-Sink-Beispiel](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples/mobile) – siehe `native_features_stub.go`, `native_features_ios.go` und `native_features_android.go`.

@note{type="note" title="API für native Funktionen und Eventbenennung"}
Zwei Konventionen solltest du kennen:

- **Native Funktionen auf der Go-Seite verwenden Plattform-Manager.** Rufe sie über die Singletons `application.IOS.*` und `application.Android.*` auf, beispielsweise mit `application.IOS.Haptic("medium")` oder `application.Android.Share(payload)`. Jeder Manager ist nur auf seiner eigenen Plattform vorhanden. Daher befinden sich seine Aufrufe in `//go:build ios`- bzw. `//go:build android`-Dateien.
- **Events erhalten entsprechend ihrer Reichweite einen Namespace.** Alles, was beide Plattformen verstehen, verwendet das Präfix `common:*` (`common:haptic`, `common:location`, …). Events, die nur eine Plattform erzeugen oder verarbeiten kann, verwenden `ios:*` oder `android:*` (beispielsweise `ios:backgroundTask`, `android:foregroundService`). Da fast alle Mobilfunktionen gemeinsam genutzt werden, benötigt dein Frontend unter `common:*` nur einen Listener pro Event.

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

## Für die Produktion bauen

@tabs{sync-key="mobile-platform"}
[iOS]
**Simulator-Build** (zum Testen im Simulator, keine Signierung erforderlich):

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

**Distributions-IPA** (für den App Store oder TestFlight):

```bash
wails3 task ios:package:ipa IOS_PLATFORM=device \
  CODESIGN_IDENTITY="..." \
  PROVISIONING_PROFILE=path/to/distribution.mobileprovision
```

@note{type="tip"}
Verwende für Uploads zu App Store Connect `wails3 task ios:xcode` und überlasse Xcode die Signierung und Archivierung. Xcode handhabt die komplexen Zertifikate, Profile und die Notarisierung automatisch.

@end

[Android]
**Debug-APK** (mit dem Android-Debug-Keystore signiert, direkt installierbar):

```bash
wails3 task android:package
wails3 task android:deploy-emulator
```

**Release-APK** (mit deinem eigenen Keystore signiert):

```bash
ANDROID_KEYSTORE_FILE=/path/to/release.jks \
ANDROID_KEYSTORE_PASSWORD=yourpassword \
ANDROID_KEY_ALIAS=youralias \
ANDROID_KEY_PASSWORD=yourkeypassword \
  wails3 task android:package
```

**Universelles APK** (arm64 und x86_64 in einer Datei):

```bash
wails3 task android:package:fat
```

@note{type="tip"}
Erstelle für Uploads in den Play Store statt eines APK ein `.aab` (Android App Bundle). Öffne `build/android/` in Android Studio und wähle **Build → Generate Signed Bundle / APK**.

@end

@end

---

## Fehlerbehebung

### `wails3 task ios:run` schlägt mit „no iOS SDKs found“ fehl

Die vollständige Xcode-Version muss installiert und ausgewählt sein:

```bash
sudo xcode-select -s /Applications/Xcode.app/Contents/Developer
xcode-select -p  # should print the Xcode path
```

#### `wails3 task android:run` schlägt mit „SDK not found“ fehl

Stelle sicher, dass `ANDROID_HOME` gesetzt und exportiert ist. Prüfe dies mit:

```bash
echo $ANDROID_HOME
ls $ANDROID_HOME/platform-tools/adb
```

#### Der Simulator startet nicht

Liste die verfügbaren Simulatoren auf und starte einen davon manuell:

```bash
xcrun simctl list devices available
xcrun simctl boot "iPhone 16"
```

#### `chrome://inspect` zeigt keine Ziele an

Die WebView muss sich im Debug-Modus befinden (die Standardeinstellung für `android:run`). Stelle sicher, dass du einen Debug-Build und keinen Produktions-Build ausführst. Prüfe außerdem, ob `adb devices` den Emulator als verbunden anzeigt.

#### Safe-Area-Abstände werden nicht angewendet

Stelle sicher, dass dein HTML das Viewport-Meta-Tag enthält:

```html
<meta name="viewport" content="width=device-width, initial-scale=1, viewport-fit=cover">
```

---

## Das Kitchen-Sink-Beispiel erkunden

Sobald deine erste App läuft, ist das **Kitchen-Sink-Beispiel** der schnellste Weg, weitere Möglichkeiten kennenzulernen. Es ist eine vollständige Wails-App, die mit einer einzigen Codebasis unter iOS, Android und auf dem Desktop läuft und unter anderem Haptik, Geolokalisierung, Biometrie, lokale Benachrichtigungen und sicheren Speicher abdeckt:

```bash
git clone https://github.com/wailsapp/wails.git
cd wails/v3/examples/mobile

wails3 task ios:run        # iOS Simulator
wails3 task android:run    # Android Emulator
wails3 task run            # Desktop
```

Durchsuche den Quellcode unter [`v3/examples/mobile`](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples/mobile). Besonders die Dateien `native_features_ios.go` und `native_features_android.go` eignen sich als Ausgangspunkte zum Kopieren und Anpassen plattformspezifischer Funktionen.

## Wie geht es weiter?

@cards{cols="2"}
iOS-Leitfaden
Vollständige Referenz: Konfigurationsoptionen, native Tabs, WKWebView-Umschalter, Geräte-Builds und Signierung.

[iOS-Leitfaden →](/guides/mobile/ios/)

---
Android-Leitfaden
Vollständige Referenz: Konfiguration, Toast-Benachrichtigungen, Paketierung für den Play Store und NDK-Details.

[Android-Leitfaden →](/guides/mobile/android/)

---
📖 Kitchen-Sink-Quellcode
Haptik, Geolokalisierung, Biometrie, Benachrichtigungen und sicherer Speicher – alles in einer ausführbaren App.

[Auf GitHub ansehen →](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples/mobile)

@end
