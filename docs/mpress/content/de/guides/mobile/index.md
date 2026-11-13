---
title: "Mobile-Übersicht"
description: "Erstellen Sie iOS- und Android-Apps aus derselben Go-Codebasis wie Ihre Desktop-App"
slug: "guides/mobile"
sourcePath: "guides/mobile/index.md"
---

Wails v3 läuft auf **iOS und Android** mit demselben `main.go` und Frontend, die Sie bereits für Desktop-Anwendungen schreiben. Es gibt weder ein separates Mobile-Projekt noch eine Brücke zur gemeinsamen Codenutzung oder eine Neuentwicklung: Das Go-Binärprogramm wird für die mobile Zielplattform kompiliert und eine native WebView stellt Ihr bestehendes Frontend dar.

@cards{cols="2"}
iOS
WKWebView mit UIKit-Host. Assets werden über ein benutzerdefiniertes `wails://`-Schema bereitgestellt – ohne offene Ports. Erfordert **macOS** mit der vollständigen Xcode-Installation.

[iOS-Anleitung →](/guides/mobile/ios/)

---
Android
Android WebView mit `WebViewAssetLoader`. Go wird über das NDK als `libwails.so` kompiliert. Funktioniert unter macOS, Linux und Windows.

[Android-Anleitung →](/guides/mobile/android/)

@end

## In Aktion: Kitchen-Sink-Beispiel

Am besten sehen Sie am **Kitchen Sink**, was möglich ist: eine einzige Wails-App, die aus einer Codebasis identisch auf iOS, Android und dem Desktop ausgeführt wird:

@linkcard{title="Mobile Kitchen Sink – GitHub" href="https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/mobile" description="Bindings · Ereignisse · Dialoge · Haptik · Geolokalisierung · Biometrie · Benachrichtigungen · sicherer Speicher · und mehr – alles aus einer einzigen main.go"}
Das Beispiel demonstriert jede wichtige mobile API-Oberfläche auf 7 Registerkarten und läuft auch auf dem Desktop. Die Registerkarten **Mobile** und **Hardware** werden auf dem Desktop durch eine Plattformprüfung im Frontend ausgeblendet. Die Go-Seite registriert bei einem Desktop-Build keine Handler für mobile `common:*`-Ereignisse. Dies ist das empfohlene Muster, um eine einzige Codebasis auf allen Plattformen bereitzustellen.

| Registerkarte | Plattformen | Gezeigte Funktionen |
| --- | --- | --- |
| **Bindings** | alle | JS → Go-Serviceaufrufe, die Werte, Structs und Fehler zurückgeben |
| **Ereignisse** | alle | Go → JS-Uhr, JS → Go → JS-Ping/Pong, OS-Systemereignisse (Akku, Netzwerk, Theme) |
| **Dialoge** | alle | Native Meldungsdialoge auf jeder Plattform |
| **System** | alle | Zwischenablage, Bildschirmmetriken, Geräteinformationen |
| **Mobile** | iOS + Android | Teilen-Menü, Wachhalten, Taschenlampe, Helligkeit, Biometrie, lokale Benachrichtigungen, sicherer Speicher |
| **Hardware** | iOS + Android | Haptik, Geolokalisierung, Beschleunigungssensor, Näherungssensor, Text-to-Speech |
| **Nativ** | iOS + Android | iOS: Haptik + WKWebView-Umschalter · Android: Vibration + Toast |

So führen Sie das Beispiel selbst aus:

```bash
git clone https://github.com/wailsapp/wails.git
cd wails/v3/examples/mobile

wails3 task ios:run        # iOS Simulator (macOS + Xcode required)
wails3 task android:run    # Android Emulator
wails3 task run            # Desktop
```

## Funktionsweise

Auf jeder Plattform gilt dasselbe Anwendungsmodell:

1. **Go-Backend** – Ihre Services, Ereignis-Handler und Anwendungslogik werden unverändert für `GOOS=ios` und `GOOS=android` kompiliert.
2. **Frontend** – exakt dasselbe HTML/JS/CSS. Das Paket `@wailsio/runtime` funktioniert identisch; Service-Bindings, Ereignisse, Dialoge und die Zwischenablage werden alle über denselben prozessinternen Transport geleitet.
3. **WebView-Host** – unter iOS ein `WKWebView` innerhalb eines `UIViewController`, unter Android ein `WebView` innerhalb einer `Activity`. Wails richtet die Nachrichtenbrücke automatisch ein.
4. **Prozessinterne Asset-Bereitstellung** – Assets werden direkt aus dem Go-Speicher bereitgestellt, nicht von einem Localhost-Server. Keine offenen Ports, kein Loopback, keine zusätzliche Latenz.

Plattformspezifisches Verhalten befindet sich in Dateien, die durch `//go:build ios` oder `//go:build android` geschützt sind. Dadurch bleibt Ihr gemeinsam genutzter Code übersichtlich.

## Voraussetzungen im Überblick

| Anforderung | iOS | Android |
| --- | --- | --- |
| Betriebssystem | nur macOS | macOS, Linux, Windows |
| Toolchain | Vollständiges Xcode (nicht nur die CLI-Tools) | Android SDK + NDK 26.3.x + JDK |
| Go | 1.25+ | 1.25+ |
| npm | ✅ | ✅ |
| Überprüfen mit | `wails3 doctor` | `wails3 doctor` |

@note{type="tip"}
Führen Sie nach der Einrichtung Ihrer Toolchain `wails3 doctor` aus. Der Befehl zeigt für jede Plattform genau an, was gefunden wurde und was fehlt.

@end

## Unterstützte Funktionen

Beide Plattformen unterstützen dieselben Kernfunktionen:

| Funktion | iOS | Android |
| --- | --- | --- |
| Service-Bindings (JS → Go) | ✅ | ✅ |
| Ereignisse (in beide Richtungen) | ✅ | ✅ |
| Meldungsdialoge | ✅ UIAlertController | ✅ AlertDialog |
| Dialoge zum Öffnen von Dateien | ✅ UIDocumentPicker | ✅ Storage Access Framework |
| Dialoge zum Speichern von Dateien | ❌ stattdessen in die Sandbox schreiben | ❌ stattdessen in die Sandbox schreiben |
| Zwischenablage | ✅ UIPasteboard | ✅ ClipboardManager |
| Bildschirm-/Safe-Area-Metriken | ✅ | ✅ |
| Lebenszyklusereignisse | ✅ `events.IOS.*` | ✅ `events.Android.*` |
| Haptik | ✅ `IOS.Haptics.*` | ✅ `Android.Haptics.Vibrate` |
| Geräteinformationen | ✅ `IOS.Device.Info()` | ✅ `Android.Device.Info()` |
| Native Registerkarten (iOS) | ✅ UITabBar | — |
| Toast-Meldungen (Android) | — | ✅ `Android.Toast.Show` |
| Mehrere Fenster | ❌ nur das erste Fenster | ❌ nur das erste Fenster |
| Fenstergeometrie / Menüs / Taskleistensymbol | bewusst wirkungslose Operationen | bewusst wirkungslose Operationen |

## Regeln für Build-Tags

Beim Schreiben von plattformabhängigem Code sind zwei wichtige Regeln zu beachten:

- **`ios` impliziert `darwin`** – eine mit `//go:build darwin` gekennzeichnete Datei wird daher auch für iOS kompiliert. Verwenden Sie `//go:build darwin && !ios`, um ausschließlich macOS anzusprechen.
- **`android` impliziert `linux`** – eine mit `//go:build linux` gekennzeichnete Datei wird daher auch für Android kompiliert. Verwenden Sie `//go:build linux && !android`, um ausschließlich Desktop-Linux anzusprechen.

Zur Laufzeit gibt `runtime.GOOS` entsprechend `"ios"` beziehungsweise `"android"` zurück.

## Plattformerkennung zur Laufzeit

Build-Tags sind für Code vorgesehen, der nur auf einer bestimmten Plattform *kompiliert* werden kann. Verwenden Sie für gewöhnliche Verzweigungen in gemeinsam genutztem Code `application.System`. Diese Funktion ist in jedem Build verfügbar (Build-Tags sind nicht erforderlich), sodass dieselbe Datei überall funktioniert:

```go
import "github.com/wailsapp/wails/v3/pkg/application"

if application.System.IsMobile() {
    // iOS or Android
} else if application.System.IsDesktop() {
    // macOS, Windows or Linux
}

// Or test a single target directly:
if application.System.IsPlatform(application.PlatformIOS) {
    // iOS only
}
```

Verfügbar sind: `IsMobile()`, `IsDesktop()`, `IsServer()` (der Build-Tag `server`) und `IsPlatform(application.PlatformMacOS | PlatformWindows | PlatformLinux | PlatformIOS | PlatformAndroid | PlatformServer)`.

Im Frontend stehen die entsprechenden Hilfsfunktionen in `@wailsio/runtime` zur Verfügung:

```js
import { System } from "@wailsio/runtime";

if (System.IsMobile()) { /* iOS or Android */ }
if (System.IsIOS()) { /* … */ }      // also IsAndroid, IsMac, IsWindows, IsLinux, IsDesktop
```

## Nächste Schritte

@cards{cols="2"}
🚀 Ihre erste mobile App
Übertragen Sie eine Wails-Desktop-App und führen Sie sie innerhalb weniger Minuten im iOS-Simulator oder Android-Emulator aus.

[Erste Schritte →](/guides/mobile/first-mobile-app/)

---
iOS-Leitfaden
Vollständige Einrichtung der iOS-Toolchain sowie Informationen zu Simulator, Geräte-Builds, Signierung, Konfiguration und API-Referenz.

[iOS-Leitfaden →](/guides/mobile/ios/)

---
Android-Leitfaden
Vollständige Einrichtung von Android SDK und NDK sowie Informationen zu Emulator, APK-Signierung, Play-Store-Paketierung und API-Referenz.

[Android-Leitfaden →](/guides/mobile/android/)

@end
