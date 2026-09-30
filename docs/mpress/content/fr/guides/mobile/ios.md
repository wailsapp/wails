---
title: "iOS"
description: "Créez et exécutez des applications Wails sur iOS — installation, simulateur, builds pour appareils, configuration et fonctionnalités natives"
slug: "guides/mobile/ios"
sourcePath: "guides/mobile/ios.md"
---

@note{type="caution" title="Fonctionnalité expérimentale"}
La prise en charge d’iOS est expérimentale et peut évoluer dans de futures versions.

@end

@note{type="tip"}
Vous découvrez le développement mobile avec Wails ? Commencez par [Votre première application mobile →](/guides/mobile/first-mobile-app/) pour suivre un guide pas à pas, puis revenez ici pour consulter la référence complète.

@end

Les applications Wails v3 s’exécutent sur iOS comme des applications entièrement natives — et surtout, elles fonctionnent *exactement* comme leur version de bureau. Le même backend Go, le même frontend et le même `@wailsio/runtime` : les liaisons de services, les événements, les boîtes de dialogue et le presse-papiers se comportent tous de manière identique, sans **aucune** adaptation propre aux appareils mobiles. Il n’existe ni base de code mobile distincte, ni couche de portage, ni API spéciale à apprendre : votre application Wails existante s’exécute simplement sur iOS. Le portage est réellement transparent : reprenez votre application telle quelle et publiez-la.

Le même `main.go` produit les builds pour ordinateur et pour iOS ; les ajustements propres à iOS sont configurés avec `application.Options.IOS`.

## Prérequis

- macOS avec la version **complète de Xcode** installée (les outils en ligne de commande seuls ne suffisent pas) — `wails3 doctor` affiche les SDK iOS qu’il peut trouver
- Go 1.25 ou version ultérieure, et npm

## Simulateur

Depuis le répertoire de votre projet :

```bash
wails3 task ios:run
```

Cette commande compile votre application, démarre un simulateur si aucun ne fonctionne déjà, puis lance l’application.

Commandes complémentaires utiles :

```bash
wails3 task ios:logs:dev    # stream the app's logs from the simulator
wails3 task ios:xcode       # open the generated Xcode project
```

Dans les builds de débogage, la WebView peut être inspectée depuis le menu Développement de Safari.

## Création des paquets

```bash
wails3 task ios:package             # production .app for the simulator
wails3 task ios:deploy-simulator    # install + launch it
```

Il s’agit de builds de production optimisés et dépouillés des informations superflues.

## Builds pour appareils

```bash
wails3 task ios:package IOS_PLATFORM=device \
    CODESIGN_IDENTITY="Apple Development: You (TEAMID)" \
    PROVISIONING_PROFILE=path/to/profile.mobileprovision

wails3 task ios:deploy-device [DEVICE_ID=<udid>]     # install + launch on a device
wails3 task ios:package:ipa IOS_PLATFORM=device ...  # distribution .ipa
```

`IOS_PLATFORM=device` crée un build pour un appareil physique. Les droits proviennent de `build/ios/entitlements.plist` et s’appliquent uniquement aux builds pour appareils — ajoutez les clés de capacités dont votre application a besoin.

@note{type="tip"}
Pour gérer automatiquement la signature, le provisionnement et les archives App Store, ouvrez le projet Xcode généré avec `wails3 task ios:xcode`, puis effectuez plutôt le build depuis Xcode.

@end

## Configuration

`build/config.yml` :

```yaml
ios:
  bundleID: com.example.myapp
  displayName: My App
  version: 1.0.0
  minIOSVersion: "15.0"
```

Les options de démarrage (`application.Options.IOS`) comprennent `DisableScroll`, `DisableBounce`, `DisableScrollIndicators`, `DisableInputAccessoryView`, `EnableBackForwardNavigationGestures`, `DisableLinkPreview`, `EnableInlineMediaPlayback`, `EnableAutoplayWithoutUserAction`, `DisableInspectable`, `UserAgent`, `ApplicationNameForUserAgent`, `BackgroundColour`, ainsi que les onglets inférieurs natifs via `EnableNativeTabs` + `NativeTabsItems`.

## Fonctionnalités natives

Les fonctionnalités propres à iOS sont disponibles via `application.IOS`, appelé depuis Go dans un fichier `//go:build ios` afin que votre code partagé reste indépendant de la plateforme. Android propose le même ensemble via `application.Android`.

Les actions ponctuelles rendent immédiatement la main :

```go
//go:build ios

application.IOS.Haptic("impact-medium") // impact-light|impact-medium|impact-heavy|success|warning|error|selection
application.IOS.Share(`{"text":"Hi","url":"https://wails.io"}`)
application.IOS.SetKeepAwake(true)
application.IOS.PostNotification(`{"title":"Done","body":"Build finished","delay":2}`)
application.IOS.SecureSet("token", "abc") // stored securely
```

Les fonctions auxiliaires de requête renvoient leurs résultats au format JSON : `SafeAreaJSON()`, `AppInfoJSON()`, `PowerJSON()`, `NetworkJSON()`, `StorageJSON()`, `GetOrientation()` et `GetBrightness()`. `StoragePath()` renvoie le chemin absolu du répertoire Application Support de l’application, un emplacement adapté aux bases de données et autres fichiers persistants (l’équivalent iOS de `getFilesDir()` sous Android). Le répertoire est créé lors du premier accès ; `StoragePath()` renvoie une chaîne vide s’il ne peut pas être créé. Vérifiez donc `""` avant de l’utiliser.

### Événements

Toute opération qui se termine ultérieurement — une demande d’autorisation, un flux de capteur ou une prise de vue avec l’appareil photo — fournit son résultat sous forme d’**événement** plutôt que comme valeur de retour. Vous pouvez écouter cet événement dans Go ou dans le frontend. Les noms portent le préfixe `common:` pour les fonctionnalités partagées avec Android et `ios:` pour celles propres à iOS.

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

| Événement | Déclenché par | Charge utile |
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

L’exemple exhaustif situé dans `v3/examples/mobile` relie de bout en bout toutes les fonctionnalités ci-dessus.

## Contrôles de la WebView

Certains comportements de la WebView peuvent également être modifiés à l’exécution depuis Go :

```go
application.IOS.SetScrollEnabled(false)
application.IOS.SetBounceEnabled(false)
application.IOS.SetScrollIndicatorsEnabled(false)
application.IOS.SetBackForwardGesturesEnabled(true)
application.IOS.SetLinkPreviewEnabled(false)
application.IOS.SetInspectableEnabled(true)
application.IOS.SetCustomUserAgent("MyApp/1.0")
```

Le paquet `@wailsio/runtime` fourni expose également un petit espace de noms iOS côté **frontend** :

```js
import { IOS } from "@wailsio/runtime";
await IOS.Haptics.Impact("medium"); // light|medium|heavy|soft|rigid
const info = await IOS.Device.Info();
```

Les sélections d’onglets inférieurs natifs sont transmises sous la forme d’un événement `nativeTabSelected` sur `window`.

## État de la prise en charge

| Domaine | État |
| --- | --- |
| Rendu du frontend et ressources | ✅ |
| Liaisons de services et événements (dans les deux sens) | ✅ |
| Boîtes de dialogue de message | ✅ |
| Boîtes de dialogue d’ouverture de fichier, de fichiers ou de répertoire | ✅ Importés sous forme de copies dans le bac à sable |
| Boîtes de dialogue d’enregistrement de fichier | ❌ Écrivez plutôt dans le bac à sable de l’application |
| Presse-papiers | ✅ |
| API des écrans | ✅ Inclut la zone de travail délimitée par la zone de sécurité |
| Événements du cycle de vie | ✅ |
| Géométrie des fenêtres, menus et zone de notification système | Sans effet sous iOS |
| Fenêtres multiples | Seule la première fenêtre est affichée |

## Notes de portage

- Le code pour ordinateur de bureau se compile sans modification pour iOS : les appels relatifs aux fenêtres, aux menus et à la zone de notification système ne font tout simplement rien.
- Remplacez les boîtes de dialogue d’enregistrement de fichier par une écriture dans le bac à sable de l’application, suivie d’un partage.
- Concevez un frontend adaptatif ; les zones de sécurité sont gérées automatiquement.
