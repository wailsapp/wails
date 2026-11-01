---
title: "Builds YAML"
description: "Configurez et exécutez le système expérimental de build natif de Wails v3 avec un seul fichier wails.yaml."
sidebar:
  order: 2
  badge:
    text: "Expérimental"
    variant: caution
---

@note{type="caution" title="Aperçu expérimental"}
La première version acceptée prend uniquement en charge la configuration littérale. Cette expérimentation propose également un prototype de hooks de cycle de vie à durée limitée, afin que leur forme puisse être évaluée avant qu’ils ne soient proposés pour une version publiée. Les étapes personnalisées, les expressions YAML et les appels typés aux outils Wails restent différés. Les projets existants continuent d’utiliser les Taskfiles jusqu’à ce que vous activiez explicitement une migration révisée.
@end

Vous découvrez Wake ? Commencez par le [tutoriel de prise en main](/experimental/wake/).
Il présente les limites de l’expérimentation, tous les écrans de configuration du projet et de la compilation, une première compilation locale ainsi qu’une migration révisée. Cette page décrit plus en détail les différents workflows de compilation, de création de paquet et de configuration.

## Activer les builds YAML

Définissez `WAILS_EXP_USE_WAKE` dans chaque terminal qui exécute les commandes expérimentales. Pour Bash ou Zsh :

```bash
export WAILS_EXP_USE_WAKE=1
```

Pour PowerShell :

```powershell
$env:WAILS_EXP_USE_WAKE = '1'
```

L’activation dépend de la présence de la variable : même une valeur vide ou `false` l’active. Supprimez-la (avec `unset WAILS_EXP_USE_WAKE` dans Bash/Zsh ou `Remove-Item Env:WAILS_EXP_USE_WAKE` dans PowerShell) pour désactiver l’expérimentation. Une fois `wails.yaml` actif, les commandes exigent de réactiver Wake et ne reviennent jamais aux Taskfiles. Sans cette variable, le routage des builds et du développement via YAML ainsi que la génération de projets YAML sont désactivés, et les commandes `migrate`, `config` et `eject` sont indisponibles. Les projets qui n’ont pas migré continuent d’utiliser le runtime Task intégré distinct.

## Configurer un projet

Pour créer et configurer un projet dans le navigateur, exécutez `wails3 init -ui` après avoir activé l’option, puis suivez le [tutoriel Wake](/experimental/wake/). Les choix propres au projet sont enregistrés dans `wails.yaml` ; modifiez ce fichier pour changer un projet existant. `wails3 setup` permet de configurer les préférences partagées, les outils et les identifiants de signature.

Une fois l’activation effectuée, ajoutez `wails.yaml` à la racine du projet. Sa présence sélectionne le système de build natif et conduit Wails à ignorer tous les Taskfiles de ce projet. Maintenez l’activation pour toutes les commandes ci-dessous. La modeline du schéma active la complétion et la validation à la volée dans les éditeurs prenant en charge un serveur de langage YAML.

```yaml
# yaml-language-server: $schema=https://v3.wails.io/schemas/wails.v3.json
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
frontend:
  directory: frontend
  install: [pnpm, install, --frozen-lockfile]
  build: [pnpm, run, build]
  dev: [pnpm, run, dev]
  output: frontend/dist
build:
  output: bin
  trim_path: true
  strip: true
profiles:
  release:
    targets:
      linux/amd64:
        formats: [deb, rpm, appimage]
```

Exécutez le build par défaut pour l’hôte ou la requête de publication nommée complète :

```bash
wails3 build
wails3 build release
```

Pour un build ponctuel, sélectionnez des cibles et des formats séparés par des virgules sans ajouter de profil :

```bash
wails3 build --targets linux/amd64 --formats deb,rpm
```

## Compiler des applications de bureau de façon croisée

Pour le backend de bureau expérimental entièrement écrit en Go, définissez `CGO_ENABLED: "0"` dans l’environnement de build. Wake sélectionne le tag `purego` pour les programmes à interface graphique macOS et Linux, et utilise directement Go sur les hôtes Windows, macOS et Linux. Dans ce mode, la création d’exécutables de bureau ne nécessite ni Docker, ni émulation, ni compilateur C, ni SDK de la plateforme cible. Les dépendances de l’application doivent également prendre en charge `CGO_ENABLED=0`.

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
build:
  output: bin
  environment:
    CGO_ENABLED: "0"
```

```bash
wails3 build --targets windows/amd64,windows/arm64,linux/amd64,linux/arm64,darwin/amd64,darwin/arm64
wails3 build --targets darwin/universal --formats dmg
```

La deuxième commande combine les deux architectures macOS et crée avec Go une image DMG non signée. Elle peut s’exécuter sur tout hôte de bureau pris en charge.

Le mode résolu est sélectionné indépendamment pour chaque exécutable. Une valeur `build.environment.CGO_ENABLED` déclarée, un environnement cible ou l’environnement du binaire remplace la valeur du terminal appelant. Sinon, un tag `purego` explicitement sélectionné, un programme CLI/serveur ou une cible Windows utilise par défaut `0` ; les programmes à interface graphique macOS et Linux utilisent par défaut `1`. `--plan --json` indique pour chaque exécutable le compilateur, le mode CGo et les tags résolus. La compilation croisée crée des fichiers pour la cible ; l’exécution d’une cible de bureau incompatible nécessite toujours l’hôte approprié.

Avec `CGO_ENABLED=1`, Wake utilise d’abord une chaîne d’outils native compatible ou le compilateur croisé `CC` configuré, puis Zig pour une compilation croisée Windows lorsqu’il est disponible, et enfin son mécanisme de repli Docker/Podman existant lorsque celui-ci est pris en charge. Une valeur `CC` explicite doit fournir le compilateur, les bibliothèques et le sysroot de la cible ; Zig seul ne fournit pas de sysroot GTK/WebKit pour Linux. La création et la signature des paquets peuvent nécessiter des outils supplémentaires, comme indiqué ci-dessous. Les cibles mobiles nécessitent CGo et les SDK de leur plateforme.

Une compilation uniquement en Go n’intègre pas les bibliothèques système de la destination. Les applications Linux à interface graphique, y compris les exécutables compagnons dotés d’une interface graphique, nécessitent toujours le runtime GTK/WebKitGTK sélectionné par leurs tags de build. Les dépendances du paquet ou un environnement de déploiement approprié doivent le fournir. Par défaut, les paquets contenant uniquement des programmes CLI/serveurs n’ajoutent pas ces dépendances graphiques ; il reste possible de remplacer explicitement les dépendances.

## Compiler plusieurs exécutables

Le projet définit le programme principal. Définissez `project.package` sur son paquet Go (`.` par défaut) et `project.kind` sur `gui`, `cli` ou `server` (`gui` par défaut). Les programmes supplémentaires utilisent des entrées nommées sous `binaries` :

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
binaries:
  notes-tool:
    package: cmd/notes-tool
    kind: cli
    tags: [tools]
    ldflags: [-X, main.Mode=release]
    environment:
      CGO_ENABLED: "0"
profiles:
  release:
    targets:
      linux/amd64: {}
```

Chaque programme dispose de ses propres entrées de compilation et de cache. Les paquets importés du projet et les fichiers incorporés sont suivis ; la modification d’un programme indépendant ne provoque pas la recompilation des autres. Les exécutables supplémentaires sont pris en charge sur les plateformes de bureau et inclus dans les bundles d’application et programmes d’installation par défaut. Les fichiers de paquet de remplacement appartenant à l’utilisateur déterminent eux-mêmes leur contenu.

```bash
wails3 build --binaries notes,notes-tool --jobs 4
wails3 build --binaries notes-tool --skip-frontend --skip-bindings
```

Un profil nommé constitue une requête complète ; utilisez ces substitutions de sélection avec un build anonyme. `--jobs` limite le nombre de workers d’exécution. La sélection exclusive de programmes CLI/serveurs omet les opérations frontend.

Pour lire une version sans exécuter le code de l’application, configurez à la fois `project.version_file` et `project.version_symbol`. Le fichier doit contenir une déclaration statique de chaîne Go. La valeur littérale `project.version` peut alors être omise. Définissez `commit_symbol: main.Commit` sur le projet ou sur un binaire nommé pour injecter la révision du dépôt dans la variable de chaîne de cet exécutable. L’absence de métadonnées du dépôt constitue une erreur lorsque l’injection du commit est demandée.

Une valeur comprise entre `upx_level: 1` et `9` compresse facultativement les exécutables Linux et Windows de publication ; zéro désactive cette fonction. UPX doit être installé. Les builds de développement et macOS n’utilisent pas la compression. En cas d’échec de la compilation ou de la compression, le programme précédemment publié est conservé.

## Créer et signer des paquets pour la distribution

### Utiliser des exécutables de bureau précompilés

Utilisez `--skip-build` lorsqu’un autre système de build ou une tâche d’intégration continue a déjà produit les exécutables :

```bash
wails3 build --targets linux/amd64 --formats deb,rpm --skip-build
wails3 run --skip-build -- --example-mode
```

Placez chaque exécutable dans `build.output` (par défaut `bin`), en utilisant `project.binary_name` pour l’exécutable principal et le nom de chaque entrée `binaries` pour les exécutables compagnons. Sous Windows, les noms se terminent par `.exe`. Lors d’une invocation multicible, les fichiers de chaque cible sont lus depuis un sous-répertoire `<os>-<arch>`. Pour une interface graphique macOS, Wake lit l’exécutable brut à l’emplacement `bin/<name>` s’il existe, sinon à l’emplacement `bin/<name>.app/Contents/MacOS/<name>`. Une entrée `darwin/universal` doit contenir les tranches d’exécutable amd64 et arm64.

Wake valide le type de fichier et l’architecture cible, copie les exécutables dans son espace de travail généré et suit leur contenu pour la mise en cache des paquets. Les entrées d’exécutables absentes, incompatibles ou constituées de liens symboliques provoquent un échec, au lieu de déclencher une recompilation ou la restauration d’une ancienne entrée. Les exécutables fournis ne sont ni remplacés ni soumis à chmod. Toute modification d’un exécutable invalide les paquets qui l’utilisent.

Ce mode ignore l’installation et les builds frontend, les bindings, la compilation Go, UPX ainsi que les hooks `before_build`/`after_build`. Il génère néanmoins les ressources des paquets, assemble les bundles d’application macOS et exécute les hooks configurés de création et de signature des paquets. Les outils de création et de signature des paquets restent soumis aux exigences du tableau ci-dessous ; Go et les outils de build frontend ne sont pas requis. Les métadonnées du paquet proviennent du manifeste actuel, y compris `version_file` lorsqu’il est configuré ; aucune version n’est déduite du binaire. Utilisez des paramètres de build et de création de paquet concordants, notamment pour le type avec interface graphique ou CLI.

`wails3 package --skip-build` est également pris en charge. Avec `wails3 run`, les arguments de lancement et l’environnement continuent de s’appliquer, et l’exécutable doit pouvoir fonctionner sur l’hôte. Ce mode nécessite un projet YAML actif et prend uniquement en charge les cibles de bureau de production ; il ne s’applique ni aux sessions de développement ni à la création de paquets mobiles.

La création d’un paquet non signé présente des exigences d’hôte et d’outils différentes de celles de la signature :

Les charges utiles IPA et les archives ZIP de notarisation macOS sont écrites en Go, tout en conservant les chemins des bundles, les autorisations d’exécution et les liens symboliques relatifs. Aucune de ces opérations ne nécessite d’exécutable `zip` ou `ditto` externe. Ces archives contiennent les fichiers de l’application, et non les métadonnées du Finder ni les fourches de ressources ; la signature macOS supprime déjà ces métadonnées, car elles invalident les bundles signés.

| Sortie | Hôtes et outils de création | Signature |
| --- | --- | --- |
| MSIX Windows | Windows, macOS ou Linux ; outil d’écriture en Go, aucun SDK Windows requis. | Signature en Go avec un fichier PFX RSA exporté sur tout hôte de bureau ; les identités du magasin de certificats ou d’un fournisseur nécessitent Windows SignTool. |
| Programme d’installation NSIS Windows | Windows, macOS ou Linux ; `makensis` natif de l’hôte requis. | Signature PE en Go avec un fichier PFX RSA exporté sur tout hôte de bureau ; les identités du magasin de certificats ou d’un fournisseur nécessitent Windows SignTool. |
| Application macOS et DMG | Windows, macOS ou Linux ; assemblage en Go, fusion en binaire universel et outil d’écriture DMG portable. | Signature en Go avec une identité RSA PKCS#12 exportée, ou signature ad hoc ; les identités du trousseau continuent d’utiliser les outils macOS. La notarisation par clé d’API et l’agrafage du ticket s’exécutent en Go. |
| DEB, RPM, Arch et Void Linux | Windows, macOS ou Linux ; outils natifs d’écriture de paquets en Go. | Les clés OpenPGP exportées utilisent le signataire natif DEB/RPM/Arch. Les anciens identifiants de clé DEB/RPM nécessitent des outils externes. La sortie Void n’est pas signée. |
| AppImage Linux | Windows, macOS ou Linux ; déploiement en Go et générateur SquashFS, plus un environnement d’exécution AppImage pour la cible. La création multiplateforme d’un paquet d’interface graphique nécessite une arborescence d’environnement d’exécution Linux fournie. | Aucun processus intégré de signature des AppImage. |

Pour les DMG signés, Wake signe normalement une copie privée de l’application avant que le générateur Go ne compresse l’image, puis signe l’image. Cela évite de monter, d’agrandir et de recompresser le DMG avec `hdiutil`. La signature s’exécute à chaque demande ; le paquet signé n’est pas réutilisé depuis le cache de build. Les projets dotés de hooks `after_package` ou `before_sign` conservent le processus d’origine, qui crée d’abord le paquet puis le signe, y compris `hdiutil`, afin de préserver les modifications que ces hooks apportent à l’image non signée. Ce processus de modification de l’image compatible avec les hooks nécessite macOS et `hdiutil`.

Pour une signature Apple portable, définissez `darwin.signing.identity: "-"` afin d’utiliser une signature ad hoc, ou utilisez `darwin.signing.certificate` avec une identité RSA PKCS#12 exportée. `darwin.signing.credential` indique le nom de la variable d’environnement contenant son mot de passe. Les signatures par certificat utilisent l’implémentation Go partagée de PFX/CMS ; aucune clé privée n’est exportée automatiquement depuis un trousseau. Les exécutables monolithiques et universels, les bundles d’application imbriqués, les frameworks et les DMG ont réussi les contrôles natifs stricts d’Apple sur les signatures, notamment le rejet après altération, avec des identités de test générées.

Pour la notarisation, une identité Developer ID exportée nécessite également un service `darwin.signing.timestamp_server` RFC 3161 compatible. Le vérificateur partagé accepte les profils d’horodatage SHA-2 qu’il prend en charge et rejette les anciens algorithmes non pris en charge. Une signature ad hoc ne peut pas être notariée. La validation native du format du certificat est distincte de la confiance réelle accordée au Developer ID et de l’acceptation d’une soumission par Apple.

L’outil d’agrafage en Go authentifie la signature du ticket d’Apple et vérifie qu’il couvre chaque architecture avant de le joindre à une copie privée. Il utilise le protocole de remise des tickets observé, qu’Apple ne publie pas comme API prise en charge. Il ne remplace pas les contrôles de révocation actuels de Gatekeeper. Les sorties d’application et de DMG signé ont réussi les contrôles natifs `stapler validate` et les contrôles stricts des signatures.

Les archives MSIX Windows sont créées directement en Go ; la création d’un paquet non signé ne nécessite ni le SDK Windows ni Docker. Les identités RSA PFX exportées permettent de signer les exécutables Windows et chaque paquet MSIX directement en Go sur tout hôte de bureau. Définissez `windows.signing.certificate` sur le chemin du fichier PFX relatif au projet et `windows.signing.credential` sur la variable d’environnement contenant son mot de passe. Les exportations AES/PBKDF2 et celles utilisant l’ancien algorithme 3DES sont prises en charge ; les exportations RC2 doivent être réexportées ou nécessitent le recours à l’outil Windows SignTool natif. Les identités provenant du magasin de certificats ou d’un fournisseur, ainsi que la signature MSI, nécessitent Windows SignTool. Le service d’horodatage par défaut est DigiCert ; les horodatages sont vérifiés avant la publication de la sortie signée. Consultez le [guide de signature](/guides/build/signing/) pour connaître les limites de format.
NSIS utilise le compilateur NSIS externe et accepte les définitions littérales via `packages.nsis.defines`. Sous macOS, lorsque les paramètres régionaux hérités ne sont pas en UTF-8, Wails fournit des paramètres régionaux UTF-8 uniquement au processus enfant du compilateur NSIS ; les environnements du terminal et de l’application restent inchangés. Le programme d’installation NSIS standard omet le programme d’amorçage WebView2 lorsque tous les exécutables inclus dans le paquet sont des outils en ligne de commande ou des serveurs. WebView2 est conservé si l’exécutable principal ou un exécutable compagnon possède une interface graphique. Les scripts NSIS personnalisés conservent le programme d’amorçage pour rester compatibles avec leur propre logique d’installation de l’environnement d’exécution.

Linux prend en charge les paquets AppImage, DEB, RPM, Arch Linux et Void. Wake assemble les AppImage directement en Go, y compris le déploiement de GTK/WebKit, la compression SquashFS et la génération zsync. Il n’appelle ni linuxdeploy, ni le plug-in GTK, ni appimagetool, ni mksquashfs. Un environnement d’exécution AppImage de type 2, épinglé et dont la somme de contrôle est vérifiée, est téléchargé pour l’architecture cible ; `packages.appimage.runtime_path` permet plutôt de sélectionner un fichier ELF d’environnement d’exécution relatif au projet.

Sous Linux, la source de bibliothèques par défaut est le système de fichiers racine de l’hôte. Pour empaqueter une interface graphique destinée à une autre architecture ou depuis Windows/macOS, définissez `packages.appimage.runtime_root` sur une arborescence racine Linux cible relative au projet, contenant ses bibliothèques, ses schémas compilés et ses caches de modules. Il s’agit de l’ensemble de l’arborescence racine, y compris `usr/`, et non du répertoire `usr/` lui-même. Fournissez un environnement d’exécution provenant de la plus ancienne distribution Linux que vous prévoyez de prendre en charge. Wake lit les importations ELF et les profils GTK/WebKit explicites sans exécuter les programmes cibles. Il inclut la fermeture de dépendances sélectionnée, mais ni glibc ni le chargeur dynamique. Les hachages du contenu des bibliothèques et des ressources invalident le cache des paquets après la mise à jour de l’environnement d’exécution.

Les builds d’interfaces graphiques entièrement en Go sont pris en charge. Wake sélectionne GTK 3 pour les exécutables avec interface graphique construits avec `gtk3`, et GTK 4 dans les autres cas, y compris pour les exécutables compagnons avec interface graphique d’une CLI ou d’un serveur. Tous les exécutables avec interface graphique d’une image doivent utiliser la même version de GTK. Pour les binaires précompilés, `packages.appimage.gtk_version` permet de sélectionner explicitement `3` ou `4` ; les bibliothèques CGo importées doivent correspondre à ce choix.

Définissez `packages.appimage.update_info` pour demander un fichier zsync et `packages.appimage.checksum: true` pour demander un fichier SHA256. Les fichiers annexes sont restaurés depuis le cache des paquets et répertoriés dans le reçu des artefacts avec l’AppImage. L’empaquetage sous Windows et macOS a été vérifié avec une arborescence GTK/WebKit fournie ; les images d’interface graphique entièrement en Go ainsi produites ont réussi les contrôles d’exécution sous Linux. Pour un même build d’interface graphique CGo, le processus natif a pris 8.77 secondes, contre 26.19 secondes pour linuxdeploy. L’image gzip native pesait 133 Mio, contre 114 Mio avec l’ancien processus zstd ; ces mesures portent donc sur l’ensemble du processus d’empaquetage, et non sur la seule compression.

Les paquets Arch préservent les préversions dans les métadonnées natives du paquet. Par exemple, `1.2.3-beta.2` devient `1.2.3beta.2-1`, qui est classé avant la version finale. Les dépendances Void acceptent des noms de paquets seuls, tels que `gtk4`, normalisés en `gtk4>=0` ; les comparaisons de versions explicites et les valeurs `name-version_revision` complètes sont préservées. Utilisez les noms de paquets Void plutôt que les noms de fichiers des bibliothèques.

La signature DEB, RPM et Arch peut utiliser des clés privées OpenPGP exportées, au format ASCII armor ou binaire, sur tout hôte de bureau pris en charge, sans GnuPG, `dpkg-sig` ni `rpmsign`. Configurez `linux.signing.certificate` avec le fichier de clé et `credential` avec le nom de la variable d’environnement contenant sa phrase secrète. DEB et RPM incorporent les signatures ; Arch publie une signature détachée. Les clés privées sont lues uniquement lors de l’étape de signature et ne sont copiées ni dans les ressources du paquet ni dans les artefacts de cache réutilisables.

Les configurations DEB/RPM existantes contenant un simple identifiant de clé GnuPG continuent d’utiliser leur signataire externe. Pour supprimer cette dépendance, exportez la clé secrète et configurez plutôt son chemin. Utilisez un chemin tel que `./signing/release.asc` ; si le fichier de clé est absent, illisible ou non valide, l’opération échoue au lieu de se rabattre sur un trousseau de clés. Les signatures DEB conservent le format de signature de paquet `dpkg-sig` ; cela ne remplace pas la signature des métadonnées de publication d’un dépôt apt.

Pour créer un paquet Arch et une signature détachée, utilisez un fichier exporté contenant une clé privée PGP au format ASCII armor et, si nécessaire, indiquez le nom d’une variable d’environnement contenant la phrase secrète :

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
build:
  environment:
    CGO_ENABLED: "0"
linux:
  signing:
    certificate: signing/private-key.asc
    credential: ARCH_SIGNING_PASSPHRASE
profiles:
  release:
    targets:
      linux/amd64:
        formats: [archlinux]
        sign: true
```

Définissez la variable d’environnement indiquée dans votre environnement local ou votre gestionnaire de secrets CI, puis exécutez `wails3 build release`. Omettez `credential` pour une clé non chiffrée. Le module de signature Go fonctionne sur tous les hôtes de bureau pris en charge et ne nécessite pas GPG. Il publie le fichier `.pkg.tar.zst` d’origine et un fichier `.pkg.tar.zst.sig` adjacent ; tous deux sont consignés séparément dans le reçu des artefacts. Le contenu de la clé n’est ni placé dans la zone intermédiaire ni stocké dans le cache de compilation. Les signatures et leur publication ne sont jamais réutilisées depuis le cache. En cas d’échec de la signature, la paire précédemment publiée est préservée ; si l’installation de l’un des fichiers échoue, la publication rétablit les deux fichiers. Cela ne rend pas deux renommages de fichiers atomiques pour les lecteurs concurrents.

Les bundles d’application macOS incluent les programmes compagnons. Définissez `darwin.microphone_usage_description` pour ajouter le texte décrivant l’utilisation ainsi que, lorsqu’aucun fichier de droits personnalisé n’est fourni, le droit d’accès à l’entrée audio. Les sources plist et de droits personnalisées restent inchangées. Une identité RSA PKCS#12 exportée utilise la signature en Go sur tout hôte de bureau ; `identity: "-"` sélectionne la signature ad hoc portable. Les identités du trousseau continuent d’utiliser `codesign` sous macOS. L’agrafage du ticket s’exécute en Go. La soumission pour notarisation peut utiliser un profil de trousseau notarial macOS existant ou les identifiants d’API d’équipe App Store Connect :

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
darwin:
  signing:
    certificate: signing/developer-id.p12
    credential: APPLE_CERT_PASSWORD
    timestamp_server: https://timestamp.digicert.com
  notarization:
    api_key: signing/AuthKey.p8
    key_id: ABCDEFGHIJ
    issuer_id: 00000000-0000-0000-0000-000000000000
```

Le mode API effectue la soumission, le téléversement et l’interrogation depuis Go sans `notarytool`. Définissez les trois champs et omettez `credential` ; ne placez pas la clé privée sous gestion de versions. La clé est lue uniquement au moment de la signature et n’est copiée ni dans les ressources générées ni dans le cache. Ce mode prend en charge les clés d’API d’équipe et les archives allant jusqu’à 5 Gio ; les téléversements multiparties et les clés d’API individuelles ne sont pas implémentés. En cas de rejet, l’identifiant de soumission est conservé et le journal de diagnostic d’Apple est enregistré dans un fichier temporaire privé dont le chemin est indiqué dans l’erreur. Une soumission annulée ou rejetée ne remplace pas la sortie publiée précédemment. Après acceptation, Wake récupère et vérifie le ticket d’Apple, puis l’agrafe en Go. Avec une identité exportée et des identifiants d’API d’équipe, la signature, la soumission et l’agrafage peuvent s’exécuter sous Linux, Windows ou macOS. Les profils de trousseau nécessitent toujours macOS et `notarytool`.

La signature portable par certificat prend en charge les identités RSA dans des fichiers PKCS#12 utilisant AES/PBKDF2 ou l’ancien chiffrement 3DES. Les exportations chiffrées avec RC2 doivent être réexportées avec un algorithme de chiffrement pris en charge. Les signatures portables notariées nécessitent un service d’horodatage RFC 3161 SHA-2 explicite ; le vérificateur Go rejette les anciennes signatures d’horodatage SHA-1. Le service de tickets suit le protocole observé d’Apple et épingle Apple Root G3 ; Gatekeeper reste l’autorité en matière de confiance pour la distribution et de révocation.

Les sorties signées conservent leur extension propre à la plateforme, par exemple `notes.signed.app` et `notes.signed.dmg`.

## Commandes de développement

Pendant `wails3 dev`, saisissez une commande, puis appuyez sur Entrée :

| Commande | Action |
| --- | --- |
| `r` | Forcer une recompilation et un redémarrage, même si le back-end n’a pas changé. |
| `b` | Régénérer les liaisons et recompiler. |
| `f` | Redémarrer le front-end. |
| `c` | Effacer le terminal. |
| `q` | Arrêter la session et ses processus enfants. |
| `h` | Afficher les commandes disponibles. |

En cas d’échec d’une compilation, l’application actuelle continue de s’exécuter. Les commandes nécessitent un terminal interactif ; une entrée redirigée ne modifie pas le comportement de la session.

## Examiner les opérations avant de les exécuter

L’inspection du plan résout le même graphe immuable que celui utilisé pour l’exécution. Le texte est destiné aux utilisateurs ; le JSON constitue un contrat versionné lisible par machine.

```bash
wails3 build release --plan
wails3 build release --plan --json
```

## Valider la configuration sans compiler

Utilisez la commande de vérification préalable dédiée dans la CI, dans une tâche d’éditeur ou avant de valider la configuration dans le dépôt. Sans profil, elle valide le manifeste de base et chaque profil nommé ; avec un profil, elle vérifie uniquement cette requête complète :

```bash
wails3 config check
wails3 config check release
```

Les erreurs de syntaxe, de schéma, de sémantique, de chemin, de cible, de format, de signature, de planificateur et de capacité de l’hôte liées à la configuration utilisent le même moteur d’affichage des diagnostics que Build et Dev. Les diagnostics identifient le champ et la plage correspondante dans la source, affichent la ligne source et un curseur, et fournissent une suggestion de correction lorsqu’elle est utile. Les valeurs d’environnement sont masquées dans les extraits de source des diagnostics. Les erreurs de l’hôte dues uniquement aux options de ligne de commande restent des erreurs de commande ordinaires au lieu d’être signalées à tort comme du YAML non valide.

## Utiliser les ressources appartenant au projet

Les chemins du manifeste désignent des fichiers qui vous appartiennent, tels que des icônes, des manifestes de plateforme, des modèles de paquets et des scripts de maintenance. Wails ne modifie jamais ces sources. Il les copie dans un espace de travail `.wails/` temporaire avant leur utilisation par un outil. Les artefacts finaux sont placés sous `build.output` et un reçu vérifié est écrit dans `.wails/artifacts/receipt.json`.

## Déployer sur Android

Les profils Android de production génèrent des AAB. Pour un déploiement local, Wails compile en interne un APK de débogage propre à l’ABI, l’installe avec ADB, puis lance l’application :

```bash
wails3 android devices
wails3 android run --device SERIAL
wails3 android run --emulator AVD_NAME
```

Sans sélecteur explicite, un terminal interactif propose les appareils connectés et les appareils virtuels Android configurés. Lorsque plusieurs cibles sont disponibles, l’automatisation doit fournir `--device` ou `--emulator`.

Attachez un flux de journaux limité à l’application et surveillez la cible jusqu’à la fermeture de l’application ou jusqu’à ce que vous appuyiez sur Ctrl+C :

```bash
wails3 android run --emulator AVD_NAME --logs
wails3 android run --emulator AVD_NAME --logs --stop-emulator
```

Les exécutions attachées suivent les redémarrages de l’application et signalent les déconnexions de la cible ou son passage hors ligne. Par défaut, les émulateurs restent en cours d’exécution afin d’accélérer les itérations. `--stop-emulator` est explicite, nécessite `--logs` et n’arrête une cible que s’il s’agit d’un émulateur, à la fin de la session attachée. Si Wails démarre un nouvel émulateur, mais que l’attente de sa disponibilité échoue ou est annulée, il arrête automatiquement ce processus.

## Exécuter des hooks de cycle de vie délimités

L’interface expérimentale des hooks appelle un script appartenant au projet à chaque phase déclarée. Un hook n’est pas une commande shell : Wails exécute directement le fichier depuis la racine du projet, sauf si `directory` est défini.

Sur les hôtes Unix, le script doit être exécutable et commencer par un shebang. Les hôtes Windows acceptent les scripts `.cmd`, `.bat` et `.ps1` au moyen d’adaptateurs d’interpréteur fixes ; Wails n’analyse de commande shell en ligne sur aucune des deux plateformes.

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
hooks:
  before_build:
    script: scripts/generate-version.sh
  after_package:
    script: scripts/check-packages.sh
profiles:
  release:
    targets:
      linux/amd64: {}
```

Les phases prises en charge sont `before_build`, `after_build`, `before_package`, `after_package`, `before_sign` et `after_sign`. Selon le cas, les hooks de compilation s’exécutent une fois pour le projet partagé ou une fois par cible. Les hooks d’empaquetage et de signature forment des barrières autour de tous les formats demandés pour une cible, tandis que les opérations sur les formats à l’intérieur de la barrière restent concurrentes.

Les hooks reçoivent une variable d’environnement appartenant à Wails : `WAILS_HOOK_CONTEXT_FILE`. Elle désigne un fichier JSON protégé contenant le contexte d’appel versionné :

```json
{
  "version": 1,
  "phase": "after_build",
  "command": "build",
  "scope": "target",
  "project_dir": "/work/notes",
  "working_directory": "/work/notes",
  "profile": "release",
  "target": { "os": "linux", "arch": "amd64" },
  "output": "/work/notes/bin/notes",
  "declared_outputs": ["/work/notes/generated/version.go"]
}
```

Le contexte contient les chemins résolus, mais ne copie jamais l’environnement hérité du processus ni les secrets. Wails le crée de manière atomique sous `.wails/hooks/context/`, avec des autorisations de fichier réservées au propriétaire. Il est supprimé après une réussite ou une annulation. En cas d’échec, il est conservé à des fins de diagnostic et son chemin figure dans l’erreur. Un code de sortie différent de zéro fait échouer l’étape Plan normale, et une annulation met fin au groupe de processus du hook.

Par défaut, les hooks s’exécutent à chaque appel. N’activez la mise en cache que si le script déclare l’intégralité de ses entrées et sorties :

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
hooks:
  before_build:
    script: scripts/generate-version.sh
    cache: true
    inputs: [version.txt]
    outputs: [generated/version.go]
profiles:
  release:
    targets:
      linux/amd64: {}
```

Le script et son mode exécutable sont des entrées implicites. Plusieurs sorties doivent partager un même répertoire autre que la racine, qui constitue l’artefact délimité que Wails enregistre et restaure. La racine des sorties ne peut contenir ni le script ni une entrée déclarée, et les arborescences de sorties mises en cache ne peuvent contenir que des fichiers ordinaires et des répertoires, jamais de liens symboliques ni de fichiers spéciaux. Un hook qui s’exécute systématiquement empêche la mise en cache de ses opérations en aval, car les effets de bord non déclarés ne peuvent pas être identifiés de manière fiable par une empreinte. Un hook pouvant être mis en cache préserve la mise en cache en aval grâce à l’identité déclarée de son artefact.

## Migrer un projet existant

Créez une ébauche inactive et examinez-la :

```bash
wails3 migrate
wails3 migrate --dry-run
wails3 migrate --json
```

Les tâches de cycle de vie prudentes, sans argument et qui appellent un script appartenant au projet sont converties en hooks. Tout autre comportement personnalisé accessible du Taskfile est signalé comme bloquant. Après avoir examiné `wails.migrated.yaml`, activez-le explicitement :

```bash
wails3 migrate --activate
```

L’activation crée `wails.yaml` de manière atomique. Elle ne supprime, ne renomme ni ne modifie les Taskfiles, `build/config.yml`, les scripts ou les ressources. Un YAML non valide échoue en tant que YAML ; il ne revient jamais à l’ancien mécanisme.

## Matérialiser toutes les valeurs par défaut

Pour examiner ou copier l’intégralité de la configuration résolue, écrivez un instantané inactif :

```bash
wails3 eject
```

Cette opération crée `wails.ejected.yaml`. Elle ne réécrit jamais la configuration active.

## Fenêtre de compatibilité

`wails3 package` et `wails3 sign` sont des alias obsolètes du même planificateur ; les nouvelles automatisations doivent utiliser `wails3 build` avec un profil ou des formats. Le déploiement est conditionné par les versions : la première version fait de YAML le mécanisme utilisé pour les nouveaux projets, la version suivante émet un avertissement lors du routage par l’ancien Taskfile, puis la version d’après supprime l’exécution du Taskfile ainsi que les alias. Pour chaque transition, la matrice de tests d’acceptation sur les hôtes natifs de cette version doit être entièrement validée.

Consultez la référence générée des champs [`wails.yaml`](/reference/wails-yaml-fields/) pour connaître chaque champ, son type, sa valeur par défaut et ses règles d’applicabilité.

La référence de configuration [`wails3 run`](/reference/wails-yaml-run/) décrit les valeurs par défaut du lancement local, les substitutions propres aux plateformes et le mécanisme de repli `go run .` en l’absence de manifeste.
