---
title: "Mobile-API"
description: "Der plattformübergreifende Manager application.Mobile – ein einzelner, durch Build-Bedingungen abgesicherter Einstiegspunkt für die nativen mobilen Funktionen, die iOS und Android gemeinsam haben"
slug: "guides/mobile/mobile-api"
sourcePath: "guides/mobile/mobile-api.md"
---

@note{type="caution" title="Experimentelle Funktion"}
Die Unterstützung für Mobilgeräte ist experimentell und kann sich in zukünftigen Versionen ändern.

@end

Native mobile Funktionen werden auf zwei Arten bereitgestellt:

- **Plattformspezifische Manager** – `application.IOS` (in `//go:build ios`-Dateien) und `application.Android` (in `//go:build android`-Dateien). Verwende diese für alle plattformspezifischen Funktionen. Die vollständigen plattformspezifischen APIs findest du in den Referenzen zu [iOS](/guides/mobile/ios/) und [Android](/guides/mobile/android/).
- **`application.Mobile`** – ein einzelner, durch Build-Bedingungen abgesicherter Manager für die Teilmenge der Funktionen, die sich auf beiden Plattformen identisch verhalten. Verwende ihn, wenn du einen einzigen Codepfad benötigst, der überall kompiliert und ausgeführt wird.

## `application.Mobile`

`application.Mobile` delegiert unter iOS an `IOS`, unter Android an `Android` und auf dem Desktop an einen wirkungslosen Stub. Da keine Build-Bedingung gilt, kannst du den Manager aus gewöhnlichem, plattformunabhängigem Go-Code aufrufen – ohne eigene `//go:build`-Dateien:

```go
// Works in any file, on any target.
// On desktop this returns "" (no-op); on device it returns the real path.
dbDir := application.Mobile.StoragePath()
if dbDir == "" {
    // Off-device, or the directory could not be created — handle accordingly.
    return
}
db, _ := sql.Open("sqlite", filepath.Join(dbDir, "app.db"))
```

`StoragePath()` gibt den absoluten Pfad zum privaten Dateiverzeichnis der App zurück – unter Android `getFilesDir()`, unter iOS das Verzeichnis „Application Support“. Dieses Verzeichnis wird für Datenbanken und andere persistente Dateien empfohlen. Auf dem Desktop sowie auf dem Gerät, wenn das Verzeichnis nicht verfügbar ist (unter iOS, wenn es nicht erstellt werden kann), wird eine leere Zeichenfolge zurückgegeben. Prüfe daher vor der Verwendung auf `""`.

@note{type="note"}
Außerhalb eines Mobilgeräts (bei Desktop-Builds) sind alle Methoden von `Mobile` wirkungslos, und jede Abfrage gibt ihren Nullwert zurück (`""` bei Zeichenfolgen). Dadurch kann plattformübergreifender Code `application.Mobile.*` bedingungslos aufrufen. Wenn du auch auf dem Desktop einen tatsächlichen Pfad benötigst, verzweige abhängig von der Plattform und greife auf `os.UserConfigDir()` oder eine ähnliche Funktion zurück.

@end

## Funktionen

Der Manager `Mobile` stellt die Funktionen bereit, deren Signaturen unter iOS und Android identisch sind:

| Funktion | API | Hinweise |
| --- | --- | --- |
| Teilen-Dialog | `Mobile.Share(json)` | `{text, url}` |
| URL extern öffnen | `Mobile.OpenURL(url)` | Systembrowser |
| Bildschirm aktiv halten | `Mobile.SetKeepAwake(bool)` |  |
| Taschenlampe | `Mobile.SetTorch(bool)` | → `common:torch` |
| Abstände des sicheren Bereichs | `Mobile.SafeAreaJSON()` | `{top,bottom,left,right}` |
| App-Informationen | `Mobile.AppInfoJSON()` | `{name,version,build,bundleId}` |
| Ausrichtungssperre | `Mobile.SetOrientation(mode)` | `portrait` / `landscape` / `auto` |
| Statusleiste | `Mobile.SetStatusBar(json)` | Stil + Sichtbarkeit |
| Speicherinformationen | `Mobile.StorageJSON()` | `{free,total}` Byte |
| Speicherpfad | `Mobile.StoragePath()` | Privates Dateiverzeichnis der App |
| Stromversorgung / Akku | `Mobile.PowerJSON()` | `{level,charging,lowPower}` |
| Netzwerkstatus | `Mobile.NetworkJSON()` | `{connected,type}` |
| Biometrie | `Mobile.BiometricAuthenticate(reason)` | → `common:biometric` |
| Sicherer Speicher | `Mobile.SecureGet(key)` / `Mobile.SecureDelete(key)` | Schlüsselbund / `EncryptedSharedPreferences` |
| Geolokalisierung | `Mobile.GetLocation()` | einmalig → `common:location` |
| Haptik | `Mobile.Haptic(type)` | Aufprall / Benachrichtigung / Auswahl |
| Beschleunigungssensor | `Mobile.SetMotion(bool)` | → `common:motion` |
| Näherungssensor | `Mobile.SetProximity(bool)` | → `common:proximity` |
| Sprachausgabe | `Mobile.Speak(text)` / `Mobile.StopSpeak()` |  |
| Tastaturabstände | `Mobile.SetKeyboardWatch(bool)` | → `common:keyboard` |
| Bildschirmaufnahme | `Mobile.SetScreenProtect(bool)` | → `common:screenCapture` |
| Kamera | `Mobile.CapturePhoto()` / `Mobile.CaptureVideo()` | → `common:capture` |

Asynchrone Ergebnisse treffen wie bei den plattformspezifischen Managern als `common:*`-Ereignisse ein — die Nutzdaten finden Sie unter [Ereignisse](/guides/mobile/ios/#events).

## Was plattformspezifisch bleibt

Funktionen, deren Struktur sich zwischen iOS und Android unterscheidet, sind **nicht** über `Mobile` verfügbar. Rufen Sie sie aus einer Datei mit Build-Tag über `application.IOS` / `application.Android` auf:

| Aspekt | iOS | Android |
| --- | --- | --- |
| Helligkeit (festlegen) | `IOS.SetBrightness(0.0-1.0)` | `Android.SetBrightness(0-100)` |
| Helligkeit / Ausrichtung (abrufen) | `IOS.GetBrightness()` / `IOS.GetOrientation()` | `Android.BrightnessJSON()` / `Android.OrientationJSON()` |
| Lokale Benachrichtigung | `IOS.PostNotification(json)` | `Android.Notify(json)` |
| Sicherer Speicher (schreiben) | `IOS.SecureSet(key, value)` | `Android.SecureSet(json)` |
| Hintergrundausführung | `IOS.BeginBackgroundTask` / `EndBackgroundTask` | `Android.StartForegroundService` / `StopForegroundService` |

@note{type="tip"}
Die Schnittstelle `MobileManager` definiert den Vertrag, auf dem `application.Mobile` basiert. Da beide Plattform-Manager ihn erfüllen müssen, behält jede oben aufgeführte Methode unter iOS und Android garantiert dieselbe Signatur bei. Sollten die Signaturen jemals voneinander abweichen, schlägt die Kompilierung des Plattform-Builds fehl.

@end
