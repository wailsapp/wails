---
title: "iOS"
description: "Créer et exécuter des applications Wails sur iOS — installation, simulateur, builds pour appareils, configuration et fonctionnalités natives"
slug: "guides/mobile/ios"
sourcePath: "guides/mobile/ios.md"
---

@note{type="caution" title="Fonctionnalité expérimentale"}
La prise en charge d’iOS est expérimentale et peut évoluer dans les prochaines versions.

@end

@note{type="tip"}
Vous débutez dans le développement mobile avec Wails ? Commencez par [Votre première application mobile →](/guides/mobile/first-mobile-app/) pour suivre un guide pas à pas, puis revenez ici pour consulter la documentation de référence complète.

@end

Les applications Wails v3 s’exécutent sur iOS comme des applications entièrement natives — et surtout, elles fonctionnent *exactement* comme leur version de bureau. Le même backend Go, le même frontend et les mêmes `@wailsio/runtime` : les liaisons de services, les événements, les boîtes de dialogue et le presse-papiers se comportent tous de manière identique, sans **aucune** adaptation propre aux appareils mobiles. Il n’existe ni base de code mobile distincte, ni couche de portage, ni API spéciale à apprendre : votre application Wails existante s’exécute simplement sur iOS. Le portage est réellement transparent : reprenez votre application telle quelle et publiez-la.

Le même `main.go` permet de créer des builds pour les environnements de bureau et iOS ; les adaptations propres à iOS se configurent via `application.Options.IOS`.

## Prérequis

- macOS avec la version **complète de Xcode** installée (les seuls outils en ligne de commande ne suffisent pas) — `wails3 doctor` affiche les SDK iOS détectés
- Go 1.26+ et npm

## Simulateur

Depuis le répertoire de votre projet :

```bash
wails3 task ios:run
```

Cette commande crée votre application, démarre un simulateur si aucun n’est déjà en cours d’exécution, puis y lance l’application.

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

Il s’agit de builds de production optimisés et allégés.

## Builds pour appareils

```bash
wails3 task ios:package IOS_PLATFORM=device \
    CODESIGN_IDENTITY="Apple Development: You (TEAMID)" \
    PROVISIONING_PROFILE=path/to/profile.mobileprovision

wails3 task ios:deploy-device [DEVICE_ID=<udid>]     # install + launch on a device
wails3 task ios:package:ipa IOS_PLATFORM=device ...  # distribution .ipa
```

`IOS_PLATFORM=device` crée un build pour un appareil physique. Les droits proviennent de `build/ios/entitlements.plist` et ne s’appliquent qu’aux builds pour appareils : ajoutez les clés de capacité dont votre application a besoin.

@note{type="tip"}
Pour gérer automatiquement la signature et le provisionnement, ainsi que créer des archives pour l’App Store, ouvrez plutôt le projet Xcode généré avec `wails3 task ios:xcode`, puis lancez le build depuis Xcode.

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

Les fonctionnalités propres à iOS sont accessibles via `application.IOS`. Appelez-les depuis Go dans un fichier `//go:build ios` afin que votre code partagé reste indépendant de la plateforme. Android propose le même ensemble via `application.Android`.

Les actions ponctuelles rendent immédiatement la main :

```go
//go:build ios

application.IOS.Haptic("impact-medium") // impact-light|impact-medium|impact-heavy|success|warning|error|selection
application.IOS.Share(`{"text":"Hi","url":"https://wails.io"}`)
application.IOS.SetKeepAwake(true)
application.IOS.PostNotification(`{"title":"Done","body":"Build finished","delay":2}`)
application.IOS.SecureSet("token", "abc") // stored securely
```

Les fonctions auxiliaires d’interrogation renvoient leurs résultats au format JSON : `SafeAreaJSON()`, `AppInfoJSON()`, `PowerJSON()`, `NetworkJSON()`, `StorageJSON()`, `GetOrientation()`, `GetBrightness()`. `StoragePath()` renvoie le chemin absolu du répertoire Application Support de l’application, un emplacement adapté aux bases de données et aux autres fichiers persistants (l’équivalent sous iOS de `getFilesDir()` sous Android). Le répertoire est créé lors du premier accès ; `StoragePath()` renvoie une chaîne vide s’il ne peut pas être créé. Vérifiez donc la présence de `""` avant de l’utiliser.

### Événements

Tout ce qui se termine ultérieurement — une demande d’autorisation, un flux de capteur ou une capture par l’appareil photo — transmet son résultat sous forme d’**événement** plutôt que comme valeur de retour. Vous pouvez écouter cet événement dans Go ou dans le frontend. Les noms portent le préfixe `common:` pour les fonctionnalités partagées avec Android et `ios:` pour celles qui sont propres à iOS.

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

L’exemple exhaustif situé sous `v3/examples/mobile` met en œuvre de bout en bout chacune des fonctionnalités ci-dessus.

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

Le `@wailsio/runtime` fourni expose également un petit espace de noms iOS côté **frontend** :

```js
import { IOS } from "@wailsio/runtime";
await IOS.Haptics.Impact("medium"); // light|medium|heavy|soft|rigid
const info = await IOS.Device.Info();
```

Les sélections d’onglets inférieurs natifs arrivent sous forme d’événement `nativeTabSelected` sur `window`.

## État de la prise en charge

| Domaine | État |
| --- | --- |
| Rendu du frontend et ressources | ✅ |
| Liaisons de services et événements (dans les deux sens) | ✅ |
| Boîtes de dialogue de message | ✅ |
| Boîtes de dialogue d’ouverture de fichier, de fichiers ou de répertoire | ✅ Importés sous forme de copies dans le bac à sable |
| Boîtes de dialogue d’enregistrement de fichier | ❌ Écrivez plutôt dans le bac à sable de l’application |
| Presse-papiers | ✅ |
| API des écrans | ✅ Inclut la zone de travail respectant les zones de sécurité |
| Événements du cycle de vie | ✅ |
| Géométrie des fenêtres, menus et zone de notification | Sans effet sous iOS |
| Fenêtres multiples | Seule la première fenêtre est affichée |

## Notes de portage

- Le code pour ordinateur se compile sans modification pour iOS — les appels relatifs aux fenêtres, aux menus et à la zone de notification n’ont tout simplement aucun effet.
- Remplacez les boîtes de dialogue d’enregistrement de fichiers par une écriture dans le bac à sable de l’application, suivie d’un partage.
- Concevez une interface frontend adaptative ; les zones de sécurité sont gérées automatiquement.
