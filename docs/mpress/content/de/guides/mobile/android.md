---
title: "Android"
description: "Wails-Anwendungen für Android erstellen und ausführen — Toolchain-Einrichtung, Emulator, APK-Signierung, Play-Store-Paketierung und API-Referenz"
slug: "guides/mobile/android"
sourcePath: "guides/mobile/android.md"
---

@note{type="caution" title="Experimentelle Funktion"}
Die Android-Unterstützung ist experimentell und kann sich in zukünftigen Versionen ändern.

@end

@note{type="tip"}
Sie entwickeln erstmals mobile Anwendungen mit Wails? Beginnen Sie mit [Ihre erste mobile App →](/guides/mobile/first-mobile-app/) für eine schrittweise Anleitung und kehren Sie anschließend zur vollständigen Referenz hierher zurück.

@end

Wails-v3Anwendungen werden unter Android als native Apps ausgeführt: Eine Android-`WebView` rendert das Frontend, Assets werden **prozessintern** über eine `WebViewAssetLoader` bereitgestellt, die vom Go-Asset-Server unterstützt wird (kein localhost-Server, keine offenen Ports), und die standardmäßige `@wailsio/runtime` funktioniert unverändert — Service-Bindings, Ereignisse, Dialoge und die Zwischenablage werden über den Go-Nachrichtenprozessor geleitet.

Dasselbe `main.go` wird für Desktop und Android erstellt. Der Go-Code wird als gemeinsam genutzte C-Bibliothek kompiliert (`libwails.so`, `GOOS=android` und die NDK-Toolchain) und von einem kleinen Java-Host geladen. Android-spezifisches Verhalten befindet sich in plattformspezifischen Go-Dateien, die durch `//go:build android` geschützt sind.

## Voraussetzungen

- Das **Android SDK** mit Platform-Tools, einer SDK-Plattform (API 35), Build-Tools und dem **NDK** (26.3.x) — `wails3 doctor` zeigt die gefundenen Komponenten an
- Ein **JDK** (z. B. OpenJDK 21) für Gradle; legen Sie `JAVA_HOME` fest, wenn sich `java` nicht in Ihrem `PATH` befindet
- Go 1.26+ und npm
- `ANDROID_HOME` (oder `ANDROID_SDK_ROOT`) mit einem Verweis auf das SDK

Installieren Sie die SDK-Komponenten mit den Befehlszeilentools:

```bash
sdkmanager "platform-tools" "platforms;android-35" "build-tools;35.0.0" \
           "ndk;26.3.11579264" "emulator" \
           "system-images;android-35;google_apis;arm64-v8a"
avdmanager create avd --name wails \
           --package "system-images;android-35;google_apis;arm64-v8a" \
           --device pixel_7
```

## Ausführung im Emulator

Führen Sie im Projektverzeichnis Folgendes aus:

```bash
wails3 task android:run
```

Dadurch wird ein Emulator gestartet, falls noch keiner ausgeführt wird. Anschließend werden die Bindings generiert, das Frontend erstellt, Ihr Go-Code für die ABI des Emulators nach `libwails.so` kompiliert, mit Gradle eine Debug-APK zusammengestellt sowie diese installiert und gestartet.

Nützliche ergänzende Befehle:

```bash
wails3 task android:logs    # stream the app's logcat output
```

In Debug-Builds kann die WebView über Chrome unter `chrome://inspect` untersucht werden.

## Paketierung

```bash
wails3 task android:package             # production release APK
wails3 task android:deploy-emulator     # install + launch it
wails3 task android:bundle              # production release AAB (Android App Bundle)
wails3 task android:bundle:fat          # release AAB containing all ABIs
wails3 task android:run:device          # debug install + launch on a physical device
wails3 task android:deploy-device       # install + launch on a physical device
DEVICE_ID=<serial> wails3 task android:run:device
DEVICE_ID=<serial> wails3 task android:deploy-device
```

Produktions-Builds verwenden `-tags production,android`, werden von nicht benötigten Symbolen bereinigt und kompilieren die internen Diagnosefunktionen des Frameworks nicht mit ein. `wails3 task android:package:fat` erstellt sowohl `arm64-v8a` als auch `x86_64` in einer einzigen APK.

Google Play verlangt für neue App-Einreichungen das Android-App-Bundle-Format (`.aab`). Außerdem müssen neue Einreichungen auf Android 15 (API 35) oder höher ausgerichtet sein. Die Projektvorlage setzt `compileSdk` und `targetSdk` in `build/android/app/build.gradle` auf 35. `wails3 task android:bundle:fat` erzeugt `bin/<AppName>.aab` mit beiden enthaltenen ABIs. Google Play generiert daraus optimierte APKs für die jeweiligen Geräte, weshalb das Fat Bundle das richtige Artefakt für Store-Uploads ist. Für lokale Tests und Emulatortests bleiben APKs der schnellste Weg, da sich ein `.aab` nicht direkt mit `adb` installieren lässt.

`android:run` und `android:deploy-emulator` sind auf Emulatoren ausgerichtete Tasks. Verwenden Sie für ein physisches Android-Gerät `android:run:device` für eine Debug-APK oder `android:deploy-device` für eine Release-APK. Beide erstellen den Build für `arm64`, wählen den ersten verbundenen Eintrag aus `adb devices` aus, der kein Emulator ist, installieren ihn und starten `com.wails.app.MainActivity`. Übergeben Sie `DEVICE_ID=<serial>`, um ein bestimmtes Gerät auszuwählen.

## Signierung und Release-Builds

Ohne einen Keystore werden Release-Builds mit dem Android-**Debug**-Keystore signiert, damit sie zu Testzwecken installiert werden können. Legen Sie Folgendes fest, um mit Ihrem eigenen Keystore zu signieren:

```bash
ANDROID_KEYSTORE_FILE=/path/to/release.jks \
ANDROID_KEYSTORE_PASSWORD=... \
ANDROID_KEY_ALIAS=... \
ANDROID_KEY_PASSWORD=... \
  wails3 task android:package
```

Dieselben Variablen signieren App Bundles: Führen Sie `wails3 task android:bundle:fat` mit gesetzten Variablen aus, um ein für Google Play geeignetes `.aab` zu erzeugen. Ohne diese Variablen wird das Bundle mit dem Debug-Keystore signiert und von Google Play abgelehnt; deshalb gibt der Task eine Warnung aus.

@note{type="tip"}
Bei [Play App Signing](https://support.google.com/googleplay/android-developer/answer/9842756) ist der Keystore, mit dem Sie lokal signieren, Ihr **Upload-Schlüssel**: Google überprüft damit Ihren Upload und signiert die App anschließend erneut mit dem von Google verwalteten App-Signaturschlüssel. Beachten Sie außerdem, dass Google Play bei jedem Upload einen höheren `versionCode` verlangt; erhöhen Sie ihn in `build/android/app/build.gradle`.

@end

## Konfiguration

Das Frontend steuert Android-Funktionen zur Laufzeit über das Laufzeitobjekt `Android`: `Android.Haptics.Vibrate(durationMs)`, `Android.Device.Info()`, `Android.Toast.Show(message)`. Der Paketname wird in den Build-Tasks durch `APP_ID` festgelegt.

## Was funktioniert und was nicht

| Bereich | Status |
| --- | --- |
| WebView und prozessinterne Assets (`WebViewAssetLoader`) | ✅ |
| Service-Bindings und Ereignisse (in beide Richtungen) | ✅ |
| Nachrichtendialoge | ✅ AlertDialog mit Button-Callbacks |
| Dialoge zum Öffnen einer oder mehrerer Dateien | ✅ Storage Access Framework (Dateien werden als Kopien in den Cache importiert) |
| Dialoge zum Öffnen eines Verzeichnisses oder Speichern einer Datei | ❌ Geben einen Fehler zurück — schreiben Sie stattdessen in die App-Sandbox |
| Zwischenablage | ✅ ClipboardManager |
| Screens API | ✅ WindowMetrics einschließlich des Arbeitsbereichs ohne Systemleisten |
| Lebenszyklusereignisse (`events.Android.*`) | ✅ |
| Haptik, Geräteinformationen und Toast-Benachrichtigungen | ✅ `Android.*`-Laufzeit-API |
| Builds für Emulatoren und physische Geräte | ✅ `android:run`, `android:run:device`, `android:deploy-emulator`, `android:deploy-device` |
| Fenstergeometrie, Menüs und Infobereich | Bewusst wirkungslose Operationen |
| Mehrere Fenster | Nur das erste Fenster wird angezeigt |

## Hinweise zur Portierung

- Desktop-Code wird unter `GOOS=android` unverändert kompiliert; Aufrufe für Geometrie, Menüs und den Infobereich werden zu wirkungslosen Operationen, da Android-Apps im Vollbildmodus ausgeführt werden.
- `android` **impliziert das `linux`Build-Tag** (Android verwendet einen Linux-Kernel): Dateien ausschließlich für Desktop-Linux benötigen `//go:build linux && !android`, und zur Laufzeit ist `runtime.GOOS` gleich `"android"`.
- Ersetzen Sie Dialoge zum Speichern von Dateien und Auswählen von Verzeichnissen durch Schreibvorgänge in die App-Sandbox sowie einen Freigabeablauf per Intent. Dialoge zum Öffnen von Dateien funktionieren und importieren die ausgewählten Dokumente als Kopien in das Cache-Verzeichnis, sodass Sie echte Dateisystempfade erhalten.
- Eine echte App wird immer mit `CGO_ENABLED=1` und dem NDK erstellt; der Pfad ohne cgo ist nur vorhanden, damit Werkzeuge wie `wails3 generate bindings` das Paket laden können.
- Gestalten Sie das Frontend responsiv; der Arbeitsbereich von `Screens` schließt die Status- und Navigationsleisten aus.
