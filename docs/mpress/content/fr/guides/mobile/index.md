---
title: "Présentation de la prise en charge mobile"
description: "Créez des applications iOS et Android à partir de la même base de code Go que votre application de bureau"
slug: "guides/mobile"
sourcePath: "guides/mobile/index.md"
---

Wails v3 s’exécute sur **iOS et Android** avec le même `main.go` et le même frontend que ceux que vous développez déjà pour les applications de bureau. Aucun projet mobile distinct, aucune passerelle de partage de code et aucune réécriture ne sont nécessaires : le binaire Go est compilé pour la plateforme mobile cible, et une WebView native affiche votre frontend existant.

@cards{cols="2"}
iOS
Hôte WKWebView + UIKit. Les ressources sont servies au moyen d’un schéma `wails://` personnalisé, sans port ouvert. Nécessite **macOS** avec la version complète de Xcode.

[Guide iOS →](/guides/mobile/ios/)

---
Android
Android WebView + `WebViewAssetLoader`. Go est compilé sous forme de `libwails.so` au moyen du NDK. Fonctionne sous macOS, Linux et Windows.

[Guide Android →](/guides/mobile/android/)

@end

## Découvrez-le en fonctionnement : exemple Kitchen Sink

La meilleure façon de comprendre les possibilités offertes consiste à examiner **Kitchen Sink**, une application Wails unique qui, à partir d’une seule base de code, s’exécute de façon identique sur iOS, Android et les plateformes de bureau :

@linkcard{title="Mobile Kitchen Sink — GitHub" href="https://github.com/wailsapp/wails/tree/master/v3/examples/mobile" description="Liaisons · Événements · Boîtes de dialogue · Retour haptique · Géolocalisation · Biométrie · Notifications · Stockage sécurisé · et plus encore — le tout depuis un seul fichier main.go"}
Cet exemple présente toutes les principales surfaces d’API mobiles dans 7 onglets et s’exécute également sur les plateformes de bureau. Les onglets **Mobile** et **Matériel** y sont masqués grâce à une vérification de la plateforme dans le frontend ; côté Go, aucun gestionnaire d’événements mobiles `common:*` n’est enregistré lors d’une compilation pour une plateforme de bureau. Il s’agit du modèle recommandé pour déployer une base de code unique sur toutes les plateformes.

| Onglet | Plateformes | Fonctionnalités présentées |
| --- | --- | --- |
| **Liaisons** | Toutes | Appels de services JS → Go renvoyant des valeurs, des structures et des erreurs |
| **Événements** | Toutes | Horloge Go → JS, ping-pong JS → Go → JS, événements système du système d’exploitation (batterie, réseau, thème) |
| **Boîtes de dialogue** | Toutes | Boîtes de dialogue de message natives sur chaque plateforme |
| **Système** | Toutes | Presse-papiers, caractéristiques de l’écran, informations sur l’appareil |
| **Mobile** | iOS + Android | Feuille de partage, maintien de l’appareil en éveil, lampe torche, luminosité, biométrie, notifications locales, stockage sécurisé |
| **Matériel** | iOS + Android | Retour haptique, géolocalisation, accéléromètre, capteur de proximité, synthèse vocale |
| **Natif** | iOS + Android | iOS : retour haptique + options de WKWebView · Android : vibration + toast |

Pour l’exécuter vous-même :

```bash
git clone https://github.com/wailsapp/wails.git
cd wails/v3/examples/mobile

wails3 task ios:run        # iOS Simulator (macOS + Xcode required)
wails3 task android:run    # Android Emulator
wails3 task run            # Desktop
```

## Fonctionnement

Le même modèle d’application s’applique à toutes les plateformes :

1. **Backend Go** — vos services, gestionnaires d’événements et logique applicative sont compilés sans modification pour `GOOS=ios` et `GOOS=android`.
2. **Frontend** — exactement les mêmes fichiers HTML/JS/CSS. Le paquet `@wailsio/runtime` fonctionne à l’identique ; les liaisons de services, les événements, les boîtes de dialogue et le presse-papiers passent tous par le même transport au sein du processus.
3. **Hôte WebView** — sur iOS, un `WKWebView` dans un `UIViewController` ; sur Android, un `WebView` dans une `Activity`. Wails configure automatiquement la passerelle de messages.
4. **Service des ressources au sein du processus** — les ressources sont servies directement depuis la mémoire de Go, et non depuis un serveur localhost. Aucun port ouvert, aucune interface de bouclage, aucune latence supplémentaire.

Le comportement propre à chaque plateforme réside dans des fichiers protégés par `//go:build ios` ou `//go:build android`, ce qui préserve la propreté de votre code partagé.

## Prérequis en un coup d’œil

| Prérequis | iOS | Android |
| --- | --- | --- |
| Système d’exploitation | macOS uniquement | macOS, Linux, Windows |
| Chaîne d’outils | Version complète de Xcode (pas seulement les outils en ligne de commande) | Android SDK + NDK 26.3.x + JDK |
| Go | 1.25+ | 1.25+ |
| npm | ✅ | ✅ |
| Vérification avec | `wails3 doctor` | `wails3 doctor` |

@note{type="tip"}
Après avoir configuré votre chaîne d’outils, exécutez `wails3 doctor` : cette commande indique précisément les éléments détectés et manquants pour chaque plateforme.

@end

## Fonctionnalités prises en charge

Les deux plateformes partagent le même ensemble de fonctionnalités de base :

| Fonctionnalité | iOS | Android |
| --- | --- | --- |
| Liaisons de services (JS → Go) | ✅ | ✅ |
| Événements (dans les deux sens) | ✅ | ✅ |
| Boîtes de dialogue de message | ✅ UIAlertController | ✅ AlertDialog |
| Boîtes de dialogue d’ouverture de fichier | ✅ UIDocumentPicker | ✅ Storage Access Framework |
| Boîtes de dialogue d’enregistrement de fichier | ❌ écrire plutôt dans le bac à sable | ❌ écrire plutôt dans le bac à sable |
| Presse-papiers | ✅ UIPasteboard | ✅ ClipboardManager |
| Écrans / métriques de la zone sûre | ✅ | ✅ |
| Événements du cycle de vie | ✅ `events.IOS.*` | ✅ `events.Android.*` |
| Retour haptique | ✅ `IOS.Haptics.*` | ✅ `Android.Haptics.Vibrate` |
| Informations sur l’appareil | ✅ `IOS.Device.Info()` | ✅ `Android.Device.Info()` |
| Onglets natifs (iOS) | ✅ UITabBar | — |
| Messages toast (Android) | — | ✅ `Android.Toast.Show` |
| Fenêtres multiples | ❌ première fenêtre uniquement | ❌ première fenêtre uniquement |
| Géométrie des fenêtres / menus / zone de notification | opérations volontairement sans effet | opérations volontairement sans effet |

## Règles des balises de compilation

Deux règles importantes à connaître lorsque vous écrivez du code conditionnel selon la plateforme :

- **`ios` implique `darwin`** : un fichier portant la balise `//go:build darwin` sera également compilé pour iOS. Pour cibler uniquement macOS, utilisez `//go:build darwin && !ios`.
- **`android` implique `linux`** : un fichier portant la balise `//go:build linux` sera également compilé pour Android. Pour cibler uniquement Linux sur ordinateur, utilisez `//go:build linux && !android`.

À l’exécution, `runtime.GOOS` renvoie respectivement `"ios"` et `"android"`.

## Détection de la plateforme à l’exécution

Les balises de compilation sont destinées au code qui peut uniquement être *compilé* sur une plateforme donnée. Pour les branchements ordinaires dans du code partagé, utilisez `application.System`, disponible dans chaque build (aucune balise de compilation nécessaire), de sorte que le même fichier fonctionne partout :

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

Disponibles : `IsMobile()`, `IsDesktop()`, `IsServer()` (la balise de compilation `server`) et `IsPlatform(application.PlatformMacOS | PlatformWindows | PlatformLinux | PlatformIOS | PlatformAndroid | PlatformServer)`.

Le frontend dispose des fonctions auxiliaires correspondantes dans `@wailsio/runtime` :

```js
import { System } from "@wailsio/runtime";

if (System.IsMobile()) { /* iOS or Android */ }
if (System.IsIOS()) { /* … */ }      // also IsAndroid, IsMac, IsWindows, IsLinux, IsDesktop
```

## Étapes suivantes

@cards{cols="2"}
🚀 Votre première application mobile
Prenez une application Wails de bureau et exécutez-la en quelques minutes dans le simulateur iOS ou l’émulateur Android.

[Commencer →](/guides/mobile/first-mobile-app/)

---
Guide iOS
Configuration complète de la chaîne d’outils iOS, simulateur, builds pour appareil, signature, configuration et référence de l’API.

[Guide iOS →](/guides/mobile/ios/)

---
Guide Android
Configuration complète des SDK et NDK Android, émulateur, signature des APK, création des packages pour le Play Store et référence de l’API.

[Guide Android →](/guides/mobile/android/)

@end
