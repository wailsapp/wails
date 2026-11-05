---
title: "API mobile"
description: "Le gestionnaire multiplateforme application.Mobile — un point d’entrée unique protégé par des contraintes de compilation pour les fonctionnalités mobiles natives communes à iOS et Android"
slug: "guides/mobile/mobile-api"
sourcePath: "guides/mobile/mobile-api.md"
---

@note{type="caution" title="Fonctionnalité expérimentale"}
La prise en charge des appareils mobiles est expérimentale et peut évoluer dans de futures versions.

@end

Les fonctionnalités mobiles natives sont exposées de deux manières :

- **Gestionnaires propres à chaque plateforme** — `application.IOS` (dans les fichiers `//go:build ios`) et `application.Android` (dans les fichiers `//go:build android`). Utilisez-les pour tout ce qui est propre à une plateforme. Consultez les références [iOS](/guides/mobile/ios/) et [Android](/guides/mobile/android/) pour connaître l’ensemble de l’API propre à chaque plateforme.
- **`application.Mobile`** — un gestionnaire unique, protégé par des contraintes de compilation, qui couvre le sous-ensemble de fonctionnalités au comportement identique sur les deux plateformes. Utilisez-le si vous souhaitez un seul chemin de code qui se compile et s’exécute partout.

## `application.Mobile`

`application.Mobile` délègue à `IOS` sous iOS, à `Android` sous Android et à une implémentation factice sans effet sur ordinateur. Comme il n’est soumis à aucune contrainte de compilation, vous pouvez l’appeler depuis du code Go ordinaire et indépendant de la plateforme, sans devoir créer vos propres fichiers `//go:build` :

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

`StoragePath()` renvoie le chemin absolu du répertoire de fichiers privé de l’application — `getFilesDir()` sous Android et le répertoire Application Support sous iOS — recommandé pour stocker les bases de données et autres fichiers persistants. Il renvoie une chaîne vide sur ordinateur, ainsi que sur l’appareil si le répertoire est indisponible (sous iOS, s’il ne peut pas être créé). Vérifiez donc la valeur `""` avant de l’utiliser.

@note{type="note"}
Hors appareil (dans les versions pour ordinateur), toutes les méthodes `Mobile` sont sans effet et chaque requête renvoie la valeur zéro de son type (`""` pour les chaînes). Le code multiplateforme peut ainsi appeler `application.Mobile.*` sans condition. Si vous avez également besoin d’un chemin réel sur ordinateur, testez la plateforme et utilisez `os.UserConfigDir()` ou une solution similaire comme solution de repli.

@end

## Fonctionnalités

Le gestionnaire `Mobile` expose les fonctionnalités dont les signatures sont identiques sous iOS et Android :

| Fonctionnalité | API | Remarques |
| --- | --- | --- |
| Feuille de partage | `Mobile.Share(json)` | `{text, url}` |
| Ouverture externe d’une URL | `Mobile.OpenURL(url)` | Navigateur système |
| Maintien de l’écran allumé | `Mobile.SetKeepAwake(bool)` |  |
| Torche / lampe de poche | `Mobile.SetTorch(bool)` | → `common:torch` |
| Marges de la zone sûre | `Mobile.SafeAreaJSON()` | `{top,bottom,left,right}` |
| Informations sur l’application | `Mobile.AppInfoJSON()` | `{name,version,build,bundleId}` |
| Verrouillage de l’orientation | `Mobile.SetOrientation(mode)` | `portrait` / `landscape` / `auto` |
| Barre d’état | `Mobile.SetStatusBar(json)` | style et visibilité |
| Informations sur le stockage | `Mobile.StorageJSON()` | `{free,total}` octets |
| Chemin de stockage | `Mobile.StoragePath()` | Répertoire de fichiers privé de l’application |
| Alimentation / batterie | `Mobile.PowerJSON()` | `{level,charging,lowPower}` |
| État du réseau | `Mobile.NetworkJSON()` | `{connected,type}` |
| Authentification biométrique | `Mobile.BiometricAuthenticate(reason)` | → `common:biometric` |
| Stockage sécurisé | `Mobile.SecureGet(key)` / `Mobile.SecureDelete(key)` | Trousseau / `EncryptedSharedPreferences` |
| Géolocalisation | `Mobile.GetLocation()` | ponctuelle → `common:location` |
| Retour haptique | `Mobile.Haptic(type)` | impact / notification / sélection |
| Accéléromètre | `Mobile.SetMotion(bool)` | → `common:motion` |
| Proximité | `Mobile.SetProximity(bool)` | → `common:proximity` |
| Synthèse vocale | `Mobile.Speak(text)` / `Mobile.StopSpeak()` |  |
| Marges du clavier | `Mobile.SetKeyboardWatch(bool)` | → `common:keyboard` |
| Capture d’écran | `Mobile.SetScreenProtect(bool)` | → `common:screenCapture` |
| Appareil photo | `Mobile.CapturePhoto()` / `Mobile.CaptureVideo()` | → `common:capture` |

Les résultats asynchrones arrivent sous forme d’événements `common:*`, exactement comme avec les gestionnaires propres à chaque plateforme — consultez [Événements](/guides/mobile/ios/#events) pour connaître les charges utiles.

## Ce qui reste propre à chaque plateforme

Les fonctionnalités dont la forme diffère entre iOS et Android ne sont **pas** disponibles sur `Mobile`. Appelez-les via `application.IOS` / `application.Android` depuis un fichier doté d’une balise de compilation :

| Fonctionnalité | iOS | Android |
| --- | --- | --- |
| Luminosité (définition) | `IOS.SetBrightness(0.0-1.0)` | `Android.SetBrightness(0-100)` |
| Luminosité / orientation (lecture) | `IOS.GetBrightness()` / `IOS.GetOrientation()` | `Android.BrightnessJSON()` / `Android.OrientationJSON()` |
| Notification locale | `IOS.PostNotification(json)` | `Android.Notify(json)` |
| Stockage sécurisé (écriture) | `IOS.SecureSet(key, value)` | `Android.SecureSet(json)` |
| Exécution en arrière-plan | `IOS.BeginBackgroundTask` / `EndBackgroundTask` | `Android.StartForegroundService` / `StopForegroundService` |

@note{type="tip"}
L’interface `MobileManager` constitue le contrat sur lequel repose `application.Mobile`. Comme les gestionnaires des deux plateformes doivent respecter ce contrat, toute méthode répertoriée ci-dessus est assurée de conserver une signature identique sur iOS et Android. Si leurs signatures divergent un jour, la compilation pour la plateforme échoue.

@end
