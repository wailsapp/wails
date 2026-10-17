---
title: "Übersicht für Mobilgeräte"
description: "Erstellen Sie iOS- und Android-Apps aus derselben Go-Codebasis wie Ihre Desktop-App"
slug: "guides/mobile"
sourcePath: "guides/mobile/index.md"
---

Wails v3 läuft auf **iOS und Android** und verwendet dabei dieselbe `main.go` und dasselbe Frontend, die Sie bereits für Desktop-Anwendungen schreiben. Es gibt weder ein separates Mobilprojekt noch eine Brücke zur gemeinsamen Codenutzung oder eine Neuentwicklung: Das Go-Binärprogramm wird für die mobile Zielplattform kompiliert, und eine native WebView rendert Ihr vorhandenes Frontend.

@cards{cols="2"}
iOS
WKWebView und UIKit-Host. Assets werden über ein benutzerdefiniertes `wails://`-Schema bereitgestellt – ohne offene Ports. Erfordert **macOS** mit vollständiger Xcode-Installation.

[iOS-Anleitung →](/guides/mobile/ios/)

---
Android
Android WebView und `WebViewAssetLoader`. Go wird über das NDK als `libwails.so` kompiliert. Funktioniert unter macOS, Linux und Windows.

[Android-Anleitung →](/guides/mobile/android/)

@end

## In Aktion: das Kitchen-Sink-Beispiel

Am besten erkennen Sie die Möglichkeiten anhand von **Kitchen Sink** – einer einzelnen Wails-App, die aus einer Codebasis identisch auf iOS, Android und Desktop-Plattformen läuft:

@linkcard{title="Mobile Kitchen Sink – GitHub" href="https://github.com/wailsapp/wails/tree/master/v3/examples/mobile" description="Bindings · Ereignisse · Dialoge · Haptik · Geolokalisierung · Biometrie · Benachrichtigungen · sicherer Speicher · und mehr – alles aus einer main.go"}
Das Beispiel demonstriert alle wichtigen mobilen API-Bereiche auf 7 Registerkarten und läuft auch auf Desktop-Plattformen. Die Registerkarten **Mobil** und **Hardware** werden dort durch eine Plattformprüfung im Frontend ausgeblendet; beim Erstellen für Desktop-Plattformen registriert die Go-Seite keine Handler für mobile `common:*`-Ereignisse. Dies ist das empfohlene Muster, um eine Codebasis auf allen Plattformen auszuliefern.

| Registerkarte | Plattformen | Gezeigte Funktionen |
| --- | --- | --- |
| **Bindings** | alle | JS-→-Go-Dienstaufrufe, die Werte, Structs und Fehler zurückgeben |
| **Ereignisse** | alle | Go-→-JS-Uhr, JS-→-Go-→-JS-Ping/Pong, Betriebssystemereignisse (Akku, Netzwerk, Design) |
| **Dialoge** | alle | Native Meldungsdialoge auf jeder Plattform |
| **System** | alle | Zwischenablage, Bildschirmmetriken, Geräteinformationen |
| **Mobil** | iOS und Android | Teilen-Menü, Ruhezustand verhindern, Taschenlampe, Helligkeit, Biometrie, lokale Benachrichtigungen, sicherer Speicher |
| **Hardware** | iOS + Android | Haptik, Geolokalisierung, Beschleunigungssensor, Näherungssensor, Sprachausgabe |
| **Nativ** | iOS + Android | iOS: Haptik + WKWebView-Umschalter · Android: Vibration + Toast-Benachrichtigung |

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

1. **Go-Backend** – Ihre Dienste, Ereignishandler und Anwendungslogik werden unverändert für `GOOS=ios` und `GOOS=android` kompiliert.
2. **Frontend** – exakt dasselbe HTML, JavaScript und CSS. Das Paket `@wailsio/runtime` funktioniert identisch; Dienstbindungen, Ereignisse, Dialoge und die Zwischenablage werden über denselben prozessinternen Transport geleitet.
3. **WebView-Host** – unter iOS eine `WKWebView` innerhalb einer `UIViewController`, unter Android eine `WebView` innerhalb einer `Activity`. Wails richtet die Nachrichtenbrücke automatisch ein.
4. **Prozessinterne Bereitstellung von Assets** – Assets werden direkt aus dem Go-Arbeitsspeicher bereitgestellt, nicht von einem localhost-Server. Keine offenen Ports, kein Loopback, keine zusätzliche Latenz.

Plattformspezifisches Verhalten befindet sich in Dateien, die durch `//go:build ios` oder `//go:build android` geschützt sind. So bleibt Ihr gemeinsam genutzter Code übersichtlich.

## Voraussetzungen auf einen Blick

| Anforderung | iOS | Android |
| --- | --- | --- |
| Betriebssystem | Nur macOS | macOS, Linux, Windows |
| Toolchain | Vollständiges Xcode (nicht nur die CLI-Tools) | Android SDK + NDK 26.3.x + JDK |
| Go | 1.25+ | 1.25+ |
| npm | ✅ | ✅ |
| Überprüfen mit | `wails3 doctor` | `wails3 doctor` |

@note{type="tip"}
Führen Sie nach dem Einrichten Ihrer Toolchain `wails3 doctor` aus. Der Befehl zeigt für jede Plattform genau an, was gefunden wurde und was fehlt.

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
| Bildschirm- und Safe-Area-Metriken | ✅ | ✅ |
| Lebenszyklusereignisse | ✅ `events.IOS.*` | ✅ `events.Android.*` |
| Haptisches Feedback | ✅ `IOS.Haptics.*` | ✅ `Android.Haptics.Vibrate` |
| Geräteinformationen | ✅ `IOS.Device.Info()` | ✅ `Android.Device.Info()` |
| Native Tabs (iOS) | ✅ UITabBar | — |
| Toast-Meldungen (Android) | — | ✅ `Android.Toast.Show` |
| Mehrere Fenster | ❌ nur das erste Fenster | ❌ nur das erste Fenster |
| Fenstergeometrie, Menüs und Tray | bewusste No-Op-Implementierungen | bewusste No-Op-Implementierungen |

## Regeln für Build-Tags

Beachten Sie beim Schreiben von plattformabhängigem Code zwei wichtige Regeln:

- **`ios` impliziert `darwin`** – eine mit `//go:build darwin` gekennzeichnete Datei wird auch für iOS kompiliert. Verwenden Sie `//go:build darwin && !ios`, um ausschließlich macOS anzusprechen.
- **`android` impliziert `linux`** – eine mit `//go:build linux` gekennzeichnete Datei wird auch für Android kompiliert. Verwenden Sie `//go:build linux && !android`, um ausschließlich Desktop-Linux anzusprechen.

Zur Laufzeit gibt `runtime.GOOS` jeweils `"ios"` beziehungsweise `"android"` zurück.

## Plattformerkennung zur Laufzeit

Build-Tags sind für Code vorgesehen, der nur auf einer bestimmten Plattform *kompiliert* werden kann. Verwenden Sie für gewöhnliche Verzweigungen in gemeinsam genutztem Code `application.System`. Dieses Objekt ist in jedem Build verfügbar (keine Build-Tags erforderlich), sodass dieselbe Datei überall funktioniert:

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

Verfügbar sind `IsMobile()`, `IsDesktop()`, `IsServer()` (das Build-Tag `server`) und `IsPlatform(application.PlatformMacOS | PlatformWindows | PlatformLinux | PlatformIOS | PlatformAndroid | PlatformServer)`.

Das Frontend enthält die entsprechenden Hilfsfunktionen in `@wailsio/runtime`:

```js
import { System } from "@wailsio/runtime";

if (System.IsMobile()) { /* iOS or Android */ }
if (System.IsIOS()) { /* … */ }      // also IsAndroid, IsMac, IsWindows, IsLinux, IsDesktop
```

## Nächste Schritte

@cards{cols="2"}
🚀 Ihre erste mobile App
Bringen Sie eine Wails-Desktop-App innerhalb weniger Minuten im iOS-Simulator oder Android-Emulator zum Laufen.

[Erste Schritte →](/guides/mobile/first-mobile-app/)

---
iOS-Leitfaden
Vollständige Einrichtung der iOS-Toolchain sowie Simulator, Geräte-Builds, Signierung, Konfiguration und API-Referenz.

[iOS-Leitfaden →](/guides/mobile/ios/)

---
Android-Leitfaden
Vollständige Einrichtung von Android SDK und NDK sowie Emulator, APK-Signierung, Play-Store-Paketierung und API-Referenz.

[Android-Leitfaden →](/guides/mobile/android/)

@end
