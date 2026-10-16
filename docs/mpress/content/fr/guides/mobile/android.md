---
title: "Android"
description: "Créez et exécutez des applications Wails sur Android — configuration de la chaîne d’outils, émulateur, signature des APK, empaquetage pour le Play Store et référence de l’API"
slug: "guides/mobile/android"
sourcePath: "guides/mobile/android.md"
---

@note{type="caution" title="Fonctionnalité expérimentale"}
La prise en charge d’Android est expérimentale et peut évoluer dans les prochaines versions.

@end

@note{type="tip"}
Vous débutez dans le développement mobile avec Wails ? Commencez par [Votre première application mobile →](/guides/mobile/first-mobile-app/) pour suivre un guide pas à pas, puis revenez ici pour consulter la référence complète.

@end

Les applications Wails v3 s’exécutent sur Android comme des applications natives : une `WebView` Android affiche le frontend, les ressources sont servies **dans le processus** par un `WebViewAssetLoader` adossé au serveur de ressources Go (sans serveur localhost ni port ouvert), et le mécanisme `@wailsio/runtime` standard fonctionne sans modification : les liaisons de services, les événements, les boîtes de dialogue et le presse-papiers passent par le processeur de messages Go.

Le même `main.go` permet de compiler pour les environnements de bureau et Android. Le code Go est compilé sous forme de bibliothèque partagée C (`libwails.so`, `GOOS=android` avec la chaîne d’outils du NDK), puis chargé par un petit hôte Java. Le comportement propre à Android se trouve dans des fichiers Go spécifiques à la plateforme, protégés par `//go:build android`.

## Prérequis

- Le **SDK Android** avec les outils de plateforme, une plateforme SDK (API 35), les outils de compilation et le **NDK** (26.3.x) — `wails3 doctor` indique les éléments détectés
- Un **JDK** (par exemple OpenJDK 21) pour Gradle ; définissez `JAVA_HOME` si `java` ne figure pas dans votre `PATH`
- Go 1.25+ et npm
- `ANDROID_HOME` (ou `ANDROID_SDK_ROOT`) pointant vers le SDK

Installez les composants du SDK avec les outils en ligne de commande :

```bash
sdkmanager "platform-tools" "platforms;android-35" "build-tools;35.0.0" \
           "ndk;26.3.11579264" "emulator" \
           "system-images;android-35;google_apis;arm64-v8a"
avdmanager create avd --name wails \
           --package "system-images;android-35;google_apis;arm64-v8a" \
           --device pixel_7
```

## Exécution sur l’émulateur

Depuis le répertoire de votre projet :

```bash
wails3 task android:run
```

Cette commande démarre un émulateur si aucun n’est en cours d’exécution, génère les liaisons, compile le frontend, compile votre code Go en `libwails.so` pour l’ABI de l’émulateur, assemble un APK de débogage avec Gradle, puis l’installe et le lance.

Commandes complémentaires utiles :

```bash
wails3 task android:logs    # stream the app's logcat output
```

Dans les versions de débogage, vous pouvez inspecter la WebView depuis Chrome à l’adresse `chrome://inspect`.

## Empaquetage

```bash
wails3 task android:package             # production release APK
wails3 task android:deploy-emulator     # install + launch it
wails3 task android:bundle              # production release AAB (Android App Bundle)
wails3 task android:bundle:fat          # release AAB containing all ABIs
wails3 task android:run:device          # debug install + launch on a physical device
wails3 task android:deploy-device       # install + launch on a physical device
DEVICE_ID=<serial> wails3 task android:run:device
DEVICE_ID=<serial> wails3 task android:deploy-device
```

Les versions de production utilisent `-tags production,android`, sont débarrassées de leurs symboles et excluent à la compilation les diagnostics internes du framework. `wails3 task android:package:fat` intègre à la fois `arm64-v8a` et `x86_64` dans un seul APK.

Google Play exige le format Android App Bundle (`.aab`) pour toute nouvelle application, qui doit cibler Android 15 (API 35) ou une version ultérieure ; le modèle de projet définit `compileSdk` et `targetSdk` sur 35 dans `build/android/app/build.gradle`. `wails3 task android:bundle:fat` produit `bin/<AppName>.aab` en incluant les deux ABI ; Google Play génère à partir de celui-ci des APK optimisés pour chaque appareil. Ce bundle universel est donc l’artefact approprié à envoyer sur la boutique. Les APK restent la solution la plus rapide pour les tests locaux et sur émulateur, car un `.aab` ne peut pas être installé directement avec `adb`.

`android:run` et `android:deploy-emulator` sont des tâches destinées à l’émulateur. Pour un appareil Android physique, utilisez `android:run:device` pour un APK de débogage ou `android:deploy-device` pour un APK de production. Ces deux tâches compilent pour `arm64`, sélectionnent la première entrée connectée qui n’est pas un émulateur dans `adb devices`, installent l’APK et lancent `com.wails.app.MainActivity`. Passez `DEVICE_ID=<serial>` pour cibler un appareil précis.

## Signature et versions de production

En l’absence de magasin de clés, les versions de production sont signées avec le magasin de clés Android de **débogage** afin de permettre leur installation à des fins de test. Pour les signer avec votre propre magasin de clés, définissez :

```bash
ANDROID_KEYSTORE_FILE=/path/to/release.jks \
ANDROID_KEYSTORE_PASSWORD=... \
ANDROID_KEY_ALIAS=... \
ANDROID_KEY_PASSWORD=... \
  wails3 task android:package
```

Les mêmes variables permettent de signer les App Bundles : exécutez `wails3 task android:bundle:fat` après les avoir définies afin de produire un `.aab` prêt pour le Play Store. Sans ces variables, le bundle est signé avec le magasin de clés de débogage et Google Play le rejettera ; la tâche affiche donc un avertissement.

@note{type="tip"}
Avec [Play App Signing](https://support.google.com/googleplay/android-developer/answer/9842756), le magasin de clés utilisé pour la signature locale contient votre **clé d’importation** : Google s’en sert pour vérifier votre envoi, puis signe de nouveau l’application avec la clé de signature de l’application qu’il gère. Notez également que Google Play exige une valeur `versionCode` plus élevée à chaque envoi ; augmentez-la dans `build/android/app/build.gradle`.

@end

## Configuration

Le frontend contrôle les fonctionnalités Android à l’exécution par l’intermédiaire de l’objet d’exécution `Android` : `Android.Haptics.Vibrate(durationMs)`, `Android.Device.Info()`, `Android.Toast.Show(message)`. Le nom du paquet est défini par `APP_ID` dans les tâches de compilation.

## Ce qui fonctionne et ce qui ne fonctionne pas

| Domaine | État |
| --- | --- |
| WebView et ressources servies dans le processus (`WebViewAssetLoader`) | ✅ |
| Liaisons de services et événements (dans les deux sens) | ✅ |
| Boîtes de dialogue de message | ✅ AlertDialog avec fonctions de rappel pour les boutons |
| Boîtes de dialogue d’ouverture d’un ou de plusieurs fichiers | ✅ Storage Access Framework (fichiers importés sous forme de copies dans le cache) |
| Boîtes de dialogue de sélection d’un répertoire ou d’enregistrement d’un fichier | ❌ Renvoient une erreur — écrivez plutôt dans le bac à sable de l’application |
| Presse-papiers | ✅ ClipboardManager |
| API Screens | ✅ WindowMetrics, y compris la zone de travail délimitée par les barres système |
| Événements du cycle de vie (`events.Android.*`) | ✅ |
| Retour haptique, informations sur l’appareil et notifications toast | ✅ API d’exécution `Android.*` |
| Compilations pour émulateur et appareil physique | ✅ `android:run`, `android:run:device`, `android:deploy-emulator`, `android:deploy-device` |
| Géométrie des fenêtres, menus et zone de notification système | Opérations intentionnellement sans effet |
| Fenêtres multiples | Seule la première fenêtre est affichée |

## Notes de portage

- Le code destiné aux environnements de bureau se compile sans modification sous `GOOS=android` ; les appels relatifs à la géométrie, aux menus et à la zone de notification deviennent des opérations sans effet, car les applications Android s’exécutent en plein écran.
- `android` **implique la balise de compilation `linux`** (Android repose sur un noyau Linux) : les fichiers réservés à Linux pour environnement de bureau nécessitent `//go:build linux && !android` et, à l’exécution, `runtime.GOOS` vaut `"android"`.
- Remplacez les boîtes de dialogue d’enregistrement de fichier et de sélection de répertoire par des écritures dans le bac à sable de l’application, associées à un flux de partage par intent. Les boîtes de dialogue d’ouverture de fichier fonctionnent et importent les documents sélectionnés sous forme de copies dans le répertoire de cache ; vous obtenez ainsi de véritables chemins de système de fichiers.
- Une application réelle est toujours compilée avec `CGO_ENABLED=1` et le NDK ; le chemin sans cgo existe uniquement pour que des outils tels que `wails3 generate bindings` puissent charger le paquet.
- Concevez une interface adaptative ; la zone de travail `Screens` exclut les barres d’état et de navigation.
