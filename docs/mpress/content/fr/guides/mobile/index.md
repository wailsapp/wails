---
title: "Présentation du mobile"
description: "Créez des applications iOS et Android à partir de la même base de code Go que votre application de bureau"
slug: "guides/mobile"
sourcePath: "guides/mobile/index.md"
---

Wails v3 fonctionne sur **iOS et Android** avec les mêmes `main.go` et frontend que ceux que vous écrivez déjà pour les applications de bureau. Il n’y a ni projet mobile distinct, ni passerelle de partage de code, ni réécriture : le binaire Go est compilé pour la cible mobile et une WebView native affiche votre frontend existant.

@cards{cols="2"}
iOS
WKWebView avec un hôte UIKit. Les ressources sont servies via un schéma `wails://` personnalisé, sans port ouvert. Nécessite **macOS** avec l’installation complète de Xcode.

[Guide iOS →](/guides/mobile/ios/)

---
Android
Android WebView + `WebViewAssetLoader`. Go est compilé sous forme de `libwails.so` via le NDK. Fonctionne sous macOS, Linux et Windows.

[Guide Android →](/guides/mobile/android/)

@end

## Découvrez-le en fonctionnement : exemple Kitchen Sink

La meilleure façon de comprendre les possibilités offertes consiste à consulter **Kitchen Sink**, une application Wails unique qui s’exécute à l’identique sur iOS, Android et les systèmes de bureau à partir d’une seule base de code :

@linkcard{title="Kitchen Sink mobile — GitHub" href="https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/mobile" description="Bindings · Événements · Boîtes de dialogue · Retour haptique · Géolocalisation · Biométrie · Notifications · Stockage sécurisé · et plus encore — le tout depuis un seul fichier main.go"}
Elle présente toutes les principales API mobiles dans 7 onglets et fonctionne également sur les systèmes de bureau. Les onglets **Mobile** et **Matériel** sont masqués sur les systèmes de bureau grâce à une vérification de la plateforme dans le frontend ; lorsqu’il est compilé pour un système de bureau, le code Go n’enregistre aucun gestionnaire pour les événements mobiles `common:*`. C’est le modèle recommandé pour distribuer une seule base de code sur toutes les plateformes.

| Onglet | Plateformes | Fonctionnalités présentées |
| --- | --- | --- |
| **Bindings** | toutes | Appels de services JS → Go renvoyant des valeurs, des structures et des erreurs |
| **Événements** | toutes | Horloge Go → JS, ping-pong JS → Go → JS, événements système du système d’exploitation (batterie, réseau, thème) |
| **Boîtes de dialogue** | toutes | Boîtes de dialogue de message natives sur chaque plateforme |
| **Système** | toutes | Presse-papiers, dimensions de l’écran, informations sur l’appareil |
| **Mobile** | iOS + Android | Feuille de partage, maintien de l’écran actif, lampe torche, luminosité, biométrie, notifications locales, stockage sécurisé |
| **Matériel** | iOS + Android | Retour haptique, géolocalisation, accéléromètre, proximité, synthèse vocale |
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

Le même modèle d’application s’applique sur toutes les plateformes :

1. **Backend Go** — vos services, gestionnaires d’événements et votre logique applicative sont compilés sans modification pour `GOOS=ios` et `GOOS=android`.
2. **Frontend** — exactement les mêmes fichiers HTML/JS/CSS. Le package `@wailsio/runtime` fonctionne à l’identique ; les bindings de services, les événements, les boîtes de dialogue et le presse-papiers passent tous par le même transport en cours de processus.
3. **Hôte WebView** — sur iOS, une `WKWebView` dans un `UIViewController` ; sur Android, une `WebView` dans une `Activity`. Wails configure automatiquement la passerelle de messages.
4. **Service des ressources en cours de processus** — les ressources sont servies directement depuis la mémoire Go et non depuis un serveur localhost. Aucun port ouvert, aucune interface de bouclage, aucune latence supplémentaire.

Le comportement propre à chaque plateforme réside dans des fichiers protégés par `//go:build ios` ou `//go:build android`, ce qui permet de garder votre code partagé propre.

## Prérequis en bref

| Prérequis | iOS | Android |
| --- | --- | --- |
| Système d’exploitation | macOS uniquement | macOS, Linux, Windows |
| Chaîne d’outils | Installation complète de Xcode (pas seulement les outils en ligne de commande) | Android SDK + NDK 26.3.x + JDK |
| Go | 1.25+ | 1.25+ |
| npm | ✅ | ✅ |
| Vérification avec | `wails3 doctor` | `wails3 doctor` |

@note{type="tip"}
Exécutez `wails3 doctor` après avoir configuré votre chaîne d’outils : la commande indique précisément les éléments détectés et ceux qui manquent pour chaque plateforme.

@end

## Fonctionnalités prises en charge

Les deux plateformes partagent le même ensemble de fonctionnalités de base :

| Fonctionnalité | iOS | Android |
| --- | --- | --- |
| Bindings de services (JS → Go) | ✅ | ✅ |
| Événements (dans les deux sens) | ✅ | ✅ |
| Boîtes de dialogue de message | ✅ UIAlertController | ✅ AlertDialog |
| Boîtes de dialogue d’ouverture de fichier | ✅ UIDocumentPicker | ✅ Storage Access Framework |
| Boîtes de dialogue d’enregistrement de fichier | ❌ écrivez plutôt dans le bac à sable | ❌ écrivez plutôt dans le bac à sable |
| Presse-papiers | ✅ UIPasteboard | ✅ ClipboardManager |
| Écrans / dimensions des zones sûres | ✅ | ✅ |
| Événements du cycle de vie | ✅ `events.IOS.*` | ✅ `events.Android.*` |
| Retour haptique | ✅ `IOS.Haptics.*` | ✅ `Android.Haptics.Vibrate` |
| Informations sur l’appareil | ✅ `IOS.Device.Info()` | ✅ `Android.Device.Info()` |
| Onglets natifs (iOS) | ✅ UITabBar | — |
| Messages toast (Android) | — | ✅ `Android.Toast.Show` |
| Fenêtres multiples | ❌ première fenêtre uniquement | ❌ première fenêtre uniquement |
| Géométrie des fenêtres / menus / zone de notification | sans effet, volontairement | sans effet, volontairement |

## Règles des balises de compilation

Deux règles importantes à connaître lorsque vous écrivez du code conditionnel selon la plateforme :

- **`ios` implique `darwin`** : un fichier portant la balise `//go:build darwin` sera également compilé pour iOS. Pour cibler uniquement macOS, utilisez `//go:build darwin && !ios`.
- **`android` implique `linux`** : un fichier portant la balise `//go:build linux` sera également compilé pour Android. Pour cibler uniquement Linux sur ordinateur, utilisez `//go:build linux && !android`.

À l’exécution, `runtime.GOOS` renvoie respectivement `"ios"` et `"android"`.

## Détection de la plateforme à l’exécution

Les balises de compilation servent au code qui ne peut être *compilé* que sur une plateforme donnée. Pour les branchements ordinaires dans du code partagé, utilisez `application.System` : cette fonction est disponible dans chaque build, sans balise de compilation, de sorte que le même fichier fonctionne partout :

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

Valeurs disponibles : `IsMobile()`, `IsDesktop()`, `IsServer()` (la balise de compilation `server`) et `IsPlatform(application.PlatformMacOS | PlatformWindows | PlatformLinux | PlatformIOS | PlatformAndroid | PlatformServer)`.

Le frontend dispose des fonctions utilitaires correspondantes dans `@wailsio/runtime` :

```js
import { System } from "@wailsio/runtime";

if (System.IsMobile()) { /* iOS or Android */ }
if (System.IsIOS()) { /* … */ }      // also IsAndroid, IsMac, IsWindows, IsLinux, IsDesktop
```

## Étapes suivantes

@cards{cols="2"}
🚀 Votre première application mobile
Prenez une application Wails pour ordinateur et exécutez-la en quelques minutes dans le simulateur iOS ou l’émulateur Android.

[Commencer →](/guides/mobile/first-mobile-app/)

---
Guide iOS
Configuration complète de la chaîne d’outils iOS, simulateur, builds pour appareils, signature, configuration et référence de l’API.

[Guide iOS →](/guides/mobile/ios/)

---
Guide Android
Configuration complète du SDK et du NDK Android, émulateur, signature des APK, création de paquets pour le Play Store et référence de l’API.

[Guide Android →](/guides/mobile/android/)

@end
