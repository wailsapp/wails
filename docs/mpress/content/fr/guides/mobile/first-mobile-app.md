---
title: "Votre première application mobile"
description: "Exécutez votre application Wails dans le simulateur iOS ou l’émulateur Android en quelques minutes"
slug: "guides/mobile/first-mobile-app"
sourcePath: "guides/mobile/first-mobile-app.md"
---

Ce guide vous explique comment exécuter une application de bureau Wails standard dans le simulateur iOS ou l’émulateur Android. **Vous n’avez pas besoin de modifier votre code Go.** Le même `main.go` permet de créer l’application pour toutes les cibles.

**Temps nécessaire :** 15–30 minutes (principalement pour l’installation de la chaîne d’outils lors de la première exécution)

## Partir d’un projet de bureau

Si vous n’en avez pas encore, créez un nouveau projet :

```bash
wails3 init -n mymobileapp
cd mymobileapp
```

Vérifiez d’abord que l’application de bureau fonctionne :

```bash
wails3 dev
```

Une fois qu’elle est ouverte, quittez-la et poursuivez. Tout ce qui fonctionne sur ordinateur fonctionne également sur mobile : pour suivre ce guide, vous n’aurez à modifier ni `main.go` ni aucun code Go.

---

## Choisir votre plateforme

@tabs{sync-key="mobile-platform"}
[Simulateur iOS]
### Prérequis

- **macOS** (les builds iOS ne sont possibles que sous macOS)
- **La version complète de Xcode**, et pas seulement les outils en ligne de commande. Installez-la depuis l’App Store, puis exécutez :
  ```bash
  sudo xcode-select -s /Applications/Xcode.app/Contents/Developer
  sudo xcodebuild -license accept
  ```


- **Go 1.25+** et **npm** (déjà installés si vous avez exécuté `wails3 init`)

Exécutez `wails3 doctor` pour effectuer la vérification : cette commande répertorie les SDK iOS qu’elle trouve.

### Exécuter dans le simulateur

@steps
### Lancer l’application
```bash
wails3 task ios:run
```

C’est tout : cette commande compile votre application, démarre un simulateur si aucun ne s’exécute déjà, puis lance l’application.

@note{type="tip"}
La première exécution prend quelques minutes, car elle compile et met en cache le framework Wails pour iOS. Toutes les exécutions suivantes sont beaucoup plus rapides.

@end

Après son lancement, votre application de bureau non modifiée s’exécute dans le simulateur iOS, avec le même `main.go` et le même frontend :

![Une application Wails par défaut s’exécutant dans le simulateur iOS](/assets/ios-simulator-first-app.png)

### Afficher les journaux en continu
Dans un autre terminal :

```bash
wails3 task ios:logs:dev
```

Cette commande affiche en continu le journal du simulateur, filtré pour votre application. La sortie de `fmt.Println` et de `log.Println` apparaît ici.

### Inspecter la WebView
Dans Safari, choisissez **Développement → Simulateur → votre application**. L’inspecteur web complet est disponible : console, débogueur, panneau Réseau, etc.

### Apporter une modification
Modifiez n’importe quel fichier du frontend (`frontend/src/main.js`, `index.html`, etc.), puis réexécutez `wails3 task ios:run`. Wails recompile le frontend et relance l’application.

Après avoir modifié du code Go, réexécutez également `wails3 task ios:run`. La recompilation Go étant incrémentale, seuls les packages modifiés sont recompilés.

@end

### Ouvrir dans Xcode (facultatif)

```bash
wails3 task ios:xcode
```

Cette commande ouvre `build/ios/` dans Xcode. Vous pouvez utiliser Xcode pour déployer l’application sur un appareil, réaliser un profilage avancé ou gérer les profils d’approvisionnement. Wails régénère le projet Xcode à chaque build ; ne modifiez donc pas directement les fichiers générés.

[Émulateur Android]
### Prérequis

Vous avez besoin du **SDK Android**, du **NDK** et d’un **JDK**. Le plus simple est d’utiliser Android Studio ou les outils en ligne de commande :

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
Ajoutez les lignes suivantes à `~/.zshrc` ou à `~/.bashrc` :

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

### Exécuter dans l’émulateur

@steps
### Lancer l’application
```bash
wails3 task android:run
```

Lors de la première exécution, cette commande :

- Démarre l’émulateur si aucun ne s’exécute
- Génère les liaisons et compile le frontend
- Compile votre code Go pour `libwails.so` à l’aide du compilateur croisé du NDK
- Assemble un APK de débogage avec Gradle
- Installe et lance l’APK dans l’émulateur

@note{type="tip"}
Le premier build télécharge Gradle et compile la chaîne d’outils du NDK ; prévoyez 5–10 minutes. Les builds suivants sont incrémentaux et prennent moins d’une minute.

@end

### Afficher les journaux en continu
Dans un autre terminal :

```bash
wails3 task android:logs
```

Cette commande exécute `adb logcat` en filtrant les résultats pour votre application. La sortie de `fmt.Println` apparaît ici.

### Inspecter la WebView
Ouvrez Chrome et accédez à `chrome://inspect`. La WebView de votre application apparaît sous **Cible distante** ; cliquez sur **inspecter** pour ouvrir les outils de développement.

### Apporter une modification
Modifiez n’importe quel fichier, puis réexécutez `wails3 task android:run`. Grâce au build incrémental de Gradle, seul le code modifié est recompilé.

@end

@end

---

## Comprendre ce qui s’est passé

Votre `main.go` n’a absolument pas changé. Wails s’est chargé de tout :

- **Système de build** — le fichier `Taskfile.yml` de votre projet contient les tâches `ios:*` et `android:*` qui pilotent la chaîne d’outils propre à chaque plateforme.
- **Compilation croisée de Go** — `GOOS=ios` ou `GOOS=android` avec la valeur `GOARCH` et le sysroot appropriés.
- **Hôte natif** — un projet Xcode généré (iOS) ou un projet Gradle (Android) qui intègre votre code Go compilé et héberge la WebView.
- **Mise à disposition des ressources** — votre `frontend/dist/` est intégré au binaire Go et servi dans le processus. Aucun serveur localhost n’est nécessaire.

---

## Adaptez votre application aux appareils mobiles

Votre application fonctionne déjà, mais elle ressemble à une application de bureau sur l’écran d’un téléphone. Quelques petits changements feront toute la différence.

### CSS adaptatif

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

Créez `mobile_ios.go` pour le code réservé à iOS :

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

Créez `mobile_android.go` pour le code réservé à Android :

```go {title="mobile_android.go"}
//go:build android

package main

import "github.com/wailsapp/wails/v3/pkg/application"

func platformOptions() application.AndroidOptions {
    return application.AndroidOptions{}
}
```

Créez `mobile_desktop.go` comme implémentation factice afin que le code partagé soit également compilé pour les plateformes de bureau :

```go {title="mobile_desktop.go"}
//go:build !ios && !android

package main

type mobileOptions struct{}

func platformOptions() mobileOptions { return mobileOptions{} }
```

### Détecter la plateforme en JavaScript et réserver l’interface mobile aux appareils mobiles

Les objets d’exécution `IOS.*` et `Android.*` n’existent que sur leur plateforme respective. Leur appel sur une plateforme de bureau déclenche une exception. La bonne approche, utilisée par Kitchen Sink, consiste à détecter la plateforme une seule fois et à masquer entièrement les contrôles réservés aux appareils mobiles :

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

Puis, dans votre HTML :

```html
<section class="mobile-only">
  <button id="btnHaptic">Haptic feedback</button>
</section>
```

Ainsi, les boutons réservés aux appareils mobiles ne sont jamais affichés sur les plateformes de bureau et vous n’avez pas besoin de protéger chaque appel par une vérification `if (isMobile)`.

Côté Go, associez cette approche à une implémentation factice utilisant une balise de build, afin que les gestionnaires d’événements ne soient enregistrés que sur les plateformes qui en ont besoin :

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

C’est exactement l’approche utilisée par [Kitchen Sink](https://github.com/wailsapp/wails/tree/master/v3/examples/mobile) : consultez `native_features_stub.go`, `native_features_ios.go` et `native_features_android.go`.

@note{type="note" title="API des fonctionnalités natives et nommage des événements"}
Deux conventions sont à connaître :

- **Les fonctionnalités natives côté Go utilisent des gestionnaires propres à chaque plateforme.** Appelez-les au moyen des singletons `application.IOS.*` et `application.Android.*`, par exemple `application.IOS.Haptic("medium")` ou `application.Android.Share(payload)`. Chaque gestionnaire n’existe que sur sa propre plateforme ; ses appels se trouvent donc dans des fichiers `//go:build ios` ou `//go:build android`.
- **Les espaces de noms des événements dépendent de leur portée.** Tout événement compris par les deux plateformes utilise le préfixe `common:*` (`common:haptic`, `common:location`, etc.) ; les événements qu’une seule plateforme peut produire ou traiter utilisent `ios:*` ou `android:*` (par exemple `ios:backgroundTask`, `android:foregroundService`). Comme presque toutes les fonctionnalités mobiles sont partagées, votre frontend conserve un seul écouteur par événement sous `common:*`.

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

### Ajouter des vibrations (Android)

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
**Build pour simulateur** (pour les tests sur simulateur, aucune signature nécessaire) :

```bash
wails3 task ios:package
wails3 task ios:deploy-simulator
```

**Build pour appareil** (nécessite une identité de signature et un profil d’approvisionnement) :

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
Pour les téléversements vers App Store Connect, utilisez `wails3 task ios:xcode` et laissez Xcode gérer la signature et l’archivage : il prend automatiquement en charge la complexité des certificats, des profils et de la notarisation.

@end

[Android]
**APK de débogage** (signé avec le keystore de débogage Android, s’installe directement) :

```bash
wails3 task android:package
wails3 task android:deploy-emulator
```

**APK de publication** (signé avec votre propre keystore) :

```bash
ANDROID_KEYSTORE_FILE=/path/to/release.jks \
ANDROID_KEYSTORE_PASSWORD=yourpassword \
ANDROID_KEY_ALIAS=youralias \
ANDROID_KEY_PASSWORD=yourkeypassword \
  wails3 task android:package
```

**APK universel** (arm64 et x86_64 dans un même fichier) :

```bash
wails3 task android:package:fat
```

@note{type="tip"}
Pour les téléversements vers le Play Store, produisez un `.aab` (Android App Bundle) plutôt qu’un APK : ouvrez `build/android/` dans Android Studio et sélectionnez **Build → Generate Signed Bundle / APK**.

@end

@end

---

## Résolution des problèmes

### Échec de `wails3 task ios:run` avec « no iOS SDKs found »

La version complète de Xcode doit être installée et sélectionnée :

```bash
sudo xcode-select -s /Applications/Xcode.app/Contents/Developer
xcode-select -p  # should print the Xcode path
```

#### Échec de `wails3 task android:run` avec « SDK not found »

Vérifiez que `ANDROID_HOME` est défini et exporté. Effectuez la vérification avec :

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

La WebView doit être en mode débogage (le mode par défaut pour `android:run`). Vérifiez que vous exécutez un build de débogage et non un build de production. Vérifiez également que `adb devices` indique que l’émulateur est connecté.

#### Les marges de la zone sûre ne sont pas appliquées

Vérifiez que votre HTML inclut la balise meta viewport :

```html
<meta name="viewport" content="width=device-width, initial-scale=1, viewport-fit=cover">
```

---

## Explorer le Kitchen Sink

Une fois votre première application en cours d’exécution, l’exemple **Kitchen Sink** constitue le moyen le plus rapide de découvrir les autres possibilités. Il s’agit d’une application Wails complète qui s’exécute sur iOS, Android et les systèmes de bureau à partir d’une base de code unique. Elle couvre le retour haptique, la géolocalisation, la biométrie, les notifications locales, le stockage sécurisé et bien plus encore :

```bash
git clone https://github.com/wailsapp/wails.git
cd wails/v3/examples/mobile

wails3 task ios:run        # iOS Simulator
wails3 task android:run    # Android Emulator
wails3 task run            # Desktop
```

Parcourez le code source à l’adresse [`v3/examples/mobile`](https://github.com/wailsapp/wails/tree/master/v3/examples/mobile). Les fichiers `native_features_ios.go` et `native_features_android.go` sont particulièrement utiles comme points de départ à copier-coller pour les fonctionnalités propres à chaque plateforme.

## Et ensuite ?

@cards{cols="2"}
Guide iOS
Référence complète : options de configuration, onglets natifs, options de WKWebView, builds pour appareils et signature.

[Guide iOS →](/guides/mobile/ios/)

---
Guide Android
Référence complète : configuration, messages toast, création de paquets pour le Play Store et détails sur le NDK.

[Guide Android →](/guides/mobile/android/)

---
📖 Code source du Kitchen Sink
Retour haptique, géolocalisation, biométrie, notifications et stockage sécurisé — le tout dans une seule application exécutable.

[Voir sur GitHub →](https://github.com/wailsapp/wails/tree/master/v3/examples/mobile)

@end
