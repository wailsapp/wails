---
title: "Votre première application mobile"
description: "Exécutez votre application Wails sur le simulateur iOS ou l’émulateur Android en quelques minutes"
slug: "guides/mobile/first-mobile-app"
sourcePath: "guides/mobile/first-mobile-app.md"
---

Ce guide part d’une application de bureau Wails standard et l’exécute sur le simulateur iOS ou l’émulateur Android. **Vous n’avez pas besoin de modifier votre code Go.** Le même `main.go` permet de générer l’application pour toutes les cibles.

**Temps nécessaire :** 15–30 minutes (principalement pour installer la chaîne d’outils lors de la première exécution)

## Partir d’un projet pour ordinateur

Si vous n’en avez pas encore, créez un nouveau projet :

```bash
wails3 init -n mymobileapp
cd mymobileapp
```

Vérifiez d’abord que l’application pour ordinateur fonctionne :

```bash
wails3 dev
```

Une fois l’application ouverte, quittez-la et passez à la suite. Tout ce qui fonctionne sur ordinateur fonctionne également sur mobile : pour ce guide, vous n’aurez à modifier ni `main.go` ni aucun code Go.

---

## Choisir votre plateforme

@tabs{sync-key="mobile-platform"}
[Simulateur iOS]
### Prérequis

- **macOS** (les builds iOS ne sont possibles que sous macOS)
- **La version complète de Xcode** — pas seulement les outils en ligne de commande. Installez-la depuis l’App Store, puis :
  ```bash
  sudo xcode-select -s /Applications/Xcode.app/Contents/Developer
  sudo xcodebuild -license accept
  ```


- **Go 1.26+** et **npm** (déjà installés si vous avez exécuté `wails3 init`)

Exécutez `wails3 doctor` pour vérifier l’installation : cette commande répertorie les SDK iOS qu’elle trouve.

### Exécuter l’application sur le simulateur

@steps
### Lancer l’application
```bash
wails3 task ios:run
```

C’est tout : cette commande génère votre application, démarre un simulateur si aucun n’est déjà en cours d’exécution, puis lance l’application.

@note{type="tip"}
La première exécution prend quelques minutes, car elle compile et met en cache le framework Wails pour iOS. Les exécutions suivantes sont beaucoup plus rapides.

@end

Après son lancement, votre application pour ordinateur, sans aucune modification, s’exécute dans le simulateur iOS, avec le même `main.go` et le même frontend :

![Une application Wails par défaut exécutée dans le simulateur iOS](/assets/ios-simulator-first-app.png)

### Suivre les journaux en continu
Dans un autre terminal :

```bash
wails3 task ios:logs:dev
```

Cette commande suit le journal du simulateur en le filtrant pour votre application. Les sorties de `fmt.Println` et `log.Println` s’affichent ici.

### Inspecter la WebView
Dans Safari, sélectionnez **Développement → Simulateur → votre application**. Vous disposez de toutes les fonctionnalités de l’inspecteur web : console, débogueur, panneau Réseau, etc.

### Effectuer une modification
Modifiez n’importe quel fichier du frontend (`frontend/src/main.js`, `index.html`, etc.), puis réexécutez `wails3 task ios:run`. Wails reconstruit le frontend et relance l’application.

Pour les modifications du code Go, réexécutez également `wails3 task ios:run`. La recompilation Go étant incrémentielle, seuls les packages modifiés sont reconstruits.

@end

### Ouvrir dans Xcode (facultatif)

```bash
wails3 task ios:xcode
```

Cette commande ouvre `build/ios/` dans Xcode. Vous pouvez utiliser Xcode pour déployer l’application sur un appareil, effectuer un profilage avancé ou gérer les profils de provisionnement. Wails régénère le projet Xcode à chaque build ; ne modifiez donc pas directement les fichiers générés.

[Émulateur Android]
### Prérequis

Vous avez besoin du **SDK Android**, du **NDK** et d’un **JDK**. La méthode la plus simple consiste à utiliser Android Studio, mais vous pouvez aussi employer les outils en ligne de commande :

@steps
### Installer les outils Android en ligne de commande
Téléchargez-les depuis [developer.android.com/studio#command-line-tools-only](https://developer.android.com/studio#command-line-tools-only), puis décompressez-les dans `~/android-sdk/cmdline-tools/latest/`.

### Installer les composants du SDK
```bash
sdkmanager "platform-tools" \
           "platforms;android-35" \
           "build-tools;35.0.0" \
           "ndk;26.3.11579264" \
           "emulator" \
           "system-images;android-35;google_apis;arm64-v8a"
```

### Créer un émulateur
```bash
avdmanager create avd \
  --name wails \
  --package "system-images;android-35;google_apis;arm64-v8a" \
  --device pixel_7
```

### Définir les variables d’environnement
Ajoutez ceci à `~/.zshrc` ou à `~/.bashrc` :

```bash
export ANDROID_HOME=~/android-sdk
export ANDROID_SDK_ROOT=~/android-sdk
export PATH=$PATH:$ANDROID_HOME/platform-tools:$ANDROID_HOME/cmdline-tools/latest/bin
```

Rechargez la configuration : `source ~/.zshrc`

### Installer un JDK
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

Exécutez `wails3 doctor` pour vérifier que tous les composants sont détectés.

### Exécuter l’application sur l’émulateur

@steps
### Lancer l’application
```bash
wails3 task android:run
```

Lors de la première exécution, cette commande :

- Démarre l’émulateur si aucun n’est en cours d’exécution
- Génère les liaisons et construit le frontend
- Compile votre code Go pour `libwails.so` à l’aide du compilateur croisé du NDK
- Assemble un APK de débogage avec Gradle
- L’installe et le lance sur l’émulateur

@note{type="tip"}
Le premier build télécharge Gradle et compile la chaîne d’outils du NDK ; prévoyez 5–10 minutes. Les builds suivants sont incrémentiels et prennent moins d’une minute.

@end

### Suivre les journaux en continu
Dans un autre terminal :

```bash
wails3 task android:logs
```

Cette commande exécute `adb logcat` en filtrant la sortie pour votre application. La sortie de `fmt.Println` s’affiche ici.

### Inspecter la WebView
Ouvrez Chrome et accédez à `chrome://inspect`. La WebView de votre application apparaît sous **Cible distante** ; cliquez sur **inspecter** pour ouvrir DevTools.

### Effectuer une modification
Modifiez n’importe quel fichier, puis réexécutez `wails3 task android:run`. Grâce au build incrémentiel de Gradle, seul le code modifié est recompilé.

@end

@end

---

## Comprendre ce qui s’est passé

Votre `main.go` n’a absolument pas changé. Wails s’est occupé de tout :

- **Système de build** — le fichier `Taskfile.yml` de votre projet contient des tâches `ios:*` et `android:*` qui pilotent la chaîne d’outils propre à chaque plateforme.
- **Compilation croisée de Go** — `GOOS=ios` ou `GOOS=android` avec le `GOARCH` et le sysroot appropriés.
- **Hôte natif** — un projet Xcode généré pour iOS ou un projet Gradle généré pour Android, qui incorpore votre code Go compilé et héberge la WebView.
- **Service des ressources** — votre `frontend/dist/` est incorporé au binaire Go et servi depuis le même processus. Aucun serveur localhost n’est nécessaire.

---

## Adapter votre application aux appareils mobiles

Votre application fonctionne déjà, mais sur l’écran d’un téléphone, elle conserve l’apparence d’une application pour ordinateur. Quelques petites modifications font toute la différence.

### CSS responsive

Les écrans mobiles sont plus étroits et utilisent des modes de saisie différents. Dans `frontend/public/style.css` (ou son équivalent) :

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

### Détecter la plateforme en Go

Utilisez des balises de build pour ajouter un comportement propre à chaque plateforme sans encombrer le code partagé.

Créez `mobile_ios.go` pour le code propre à iOS :

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

Créez `mobile_android.go` pour le code propre à Android :

```go {title="mobile_android.go"}
//go:build android

package main

import "github.com/wailsapp/wails/v3/pkg/application"

func platformOptions() application.AndroidOptions {
    return application.AndroidOptions{}
}
```

Créez `mobile_desktop.go` comme implémentation factice afin que le code partagé soit également compilable pour ordinateur :

```go {title="mobile_desktop.go"}
//go:build !ios && !android

package main

type mobileOptions struct{}

func platformOptions() mobileOptions { return mobileOptions{} }
```

### Détecter la plateforme en JavaScript et réserver l’interface mobile aux appareils mobiles

Les objets d’exécution `IOS.*` et `Android.*` n’existent que sur leurs plateformes respectives. Leur appel sur ordinateur déclenche une exception. La bonne approche, utilisée par Kitchen Sink, consiste à détecter la plateforme une seule fois et à masquer entièrement les contrôles propres aux appareils mobiles :

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

Puis, dans votre code HTML :

```html
<section class="mobile-only">
  <button id="btnHaptic">Haptic feedback</button>
</section>
```

Ainsi, les boutons réservés aux plateformes mobiles ne sont jamais affichés sur ordinateur et vous n’avez pas besoin de protéger chaque appel par une vérification `if (isMobile)`.

Côté Go, associez cette approche à un stub utilisant une balise de compilation afin que les gestionnaires d’événements ne soient enregistrés que sur les plateformes qui en ont besoin :

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

C’est exactement l’approche utilisée par le [Kitchen Sink](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples/mobile) : consultez `native_features_stub.go`, `native_features_ios.go` et `native_features_android.go`.

@note{type="note" title="API des fonctionnalités natives et nommage des événements"}
Deux conventions sont à connaître :

- **Les fonctionnalités natives côté Go utilisent des gestionnaires propres à chaque plateforme.** Appelez-les par l’intermédiaire des singletons `application.IOS.*` et `application.Android.*`, par exemple avec `application.IOS.Haptic("medium")` ou `application.Android.Share(payload)`. Chaque gestionnaire n’existe que sur sa propre plateforme ; ses appels doivent donc se trouver dans des fichiers `//go:build ios` ou `//go:build android`.
- **Les événements sont placés dans des espaces de noms selon leur portée.** Tout événement compris par les deux plateformes utilise le préfixe `common:*` (`common:haptic`, `common:location`, etc.) ; les événements qu’une seule plateforme peut produire ou gérer utilisent `ios:*` ou `android:*` (par exemple `ios:backgroundTask`, `android:foregroundService`). Comme presque toutes les fonctionnalités mobiles sont communes, votre frontend conserve un seul écouteur par événement sous `common:*`.

@end

### Ajouter un retour haptique (iOS)

```javascript
import { IOS } from '@wailsio/runtime';

async function onButtonTap() {
  if (isIOS) {
    await IOS.Haptics.Impact({ style: 'medium' });
  }
  // ... rest of your handler
}
```

### Ajouter une vibration (Android)

```javascript
import { Android } from '@wailsio/runtime';

async function onButtonTap() {
  if (isAndroid) {
    await Android.Haptics.Vibrate(50); // 50ms
  }
}
```

---

## Compiler pour la production

@tabs{sync-key="mobile-platform"}
[iOS]
**Build pour simulateur** (pour les tests dans le simulateur, aucune signature requise) :

```bash
wails3 task ios:package
wails3 task ios:deploy-simulator
```

**Build pour appareil** (nécessite une identité de signature et un profil de provisionnement) :

```bash
wails3 task ios:package \
  IOS_PLATFORM=device \
  CODESIGN_IDENTITY="Apple Development: You (TEAMID)" \
  PROVISIONING_PROFILE=path/to/profile.mobileprovision

wails3 task ios:deploy-device   # installs via xcrun devicectl
```

**IPA de distribution** (pour l’App Store ou TestFlight) :

```bash
wails3 task ios:package:ipa IOS_PLATFORM=device \
  CODESIGN_IDENTITY="..." \
  PROVISIONING_PROFILE=path/to/distribution.mobileprovision
```

@note{type="tip"}
Pour les envois vers App Store Connect, utilisez `wails3 task ios:xcode` et laissez Xcode gérer la signature et l’archivage : il prend automatiquement en charge la complexité des certificats, des profils et de la notarisation.

@end

[Android]
**APK de débogage** (signé avec le keystore de débogage Android et directement installable) :

```bash
wails3 task android:package
wails3 task android:deploy-emulator
```

**APK de version** (signé avec votre propre keystore) :

```bash
ANDROID_KEYSTORE_FILE=/path/to/release.jks \
ANDROID_KEYSTORE_PASSWORD=yourpassword \
ANDROID_KEY_ALIAS=youralias \
ANDROID_KEY_PASSWORD=yourkeypassword \
  wails3 task android:package
```

**APK universel** (arm64 et x86_64 dans un seul fichier) :

```bash
wails3 task android:package:fat
```

@note{type="tip"}
Pour les envois vers le Play Store, produisez un `.aab` (Android App Bundle) plutôt qu’un APK : ouvrez `build/android/` dans Android Studio, puis sélectionnez **Build → Generate Signed Bundle / APK**.

@end

@end

---

## Dépannage

### Échec de `wails3 task ios:run` avec le message « no iOS SDKs found »

La version complète de Xcode doit être installée et sélectionnée :

```bash
sudo xcode-select -s /Applications/Xcode.app/Contents/Developer
xcode-select -p  # should print the Xcode path
```

#### Échec de `wails3 task android:run` avec le message « SDK not found »

Vérifiez que `ANDROID_HOME` est défini et exporté. Pour le vérifier :

```bash
echo $ANDROID_HOME
ls $ANDROID_HOME/platform-tools/adb
```

#### Le simulateur ne démarre pas

Répertoriez les simulateurs disponibles et démarrez-en un manuellement :

```bash
xcrun simctl list devices available
xcrun simctl boot "iPhone 16"
```

#### `chrome://inspect` n’affiche aucune cible

La WebView doit être en mode débogage, ce qui est le réglage par défaut pour `android:run`. Assurez-vous d’exécuter un build de débogage et non un build de production. Vérifiez également que `adb devices` indique que l’émulateur est connecté.

#### Les marges de la zone sûre ne sont pas appliquées

Assurez-vous que votre code HTML inclut la balise meta viewport :

```html
<meta name="viewport" content="width=device-width, initial-scale=1, viewport-fit=cover">
```

---

## Explorer le Kitchen Sink

Une fois votre première application en cours d’exécution, l’exemple **Kitchen Sink** est le moyen le plus rapide de découvrir les autres possibilités. Il s’agit d’une application Wails complète qui s’exécute sur iOS, Android et les environnements de bureau à partir d’une base de code unique. Elle couvre notamment le retour haptique, la géolocalisation, la biométrie, les notifications locales et le stockage sécurisé :

```bash
git clone https://github.com/wailsapp/wails.git
cd wails/v3/examples/mobile

wails3 task ios:run        # iOS Simulator
wails3 task android:run    # Android Emulator
wails3 task run            # Desktop
```

Parcourez le code source à l’adresse [`v3/examples/mobile`](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples/mobile). Les fichiers `native_features_ios.go` et `native_features_android.go` constituent des points de départ particulièrement utiles à copier-coller pour les fonctionnalités propres à chaque plateforme.

## Étapes suivantes

@cards{cols="2"}
Guide iOS
Référence complète : options de configuration, onglets natifs, options de WKWebView, builds pour appareils et signature.

[Guide iOS →](/guides/mobile/ios/)

---
Guide Android
Référence complète : configuration, toasts, empaquetage pour le Play Store et détails du NDK.

[Guide Android →](/guides/mobile/android/)

---
📖 Code source du Kitchen Sink
Retour haptique, géolocalisation, biométrie, notifications et stockage sécurisé : le tout dans une seule application exécutable.

[Voir sur GitHub →](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples/mobile)

@end
