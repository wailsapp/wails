---
title: "Configuration d’exécution YAML"
description: "Valeurs d’exécution locales par défaut, substitutions propres aux plateformes et comportement de lancement de wails3 run."
---

## Disponibilité

Cette référence décrit l’extension d’exécution expérimentale. Elle ne garantit pas sa disponibilité dans une version publiée de la CLI.

`wails3 run` compile et lance l’application en utilisant les valeurs locales par défaut du manifeste. L’exécution pilotée par le manifeste nécessite l’activation explicite de [`WAILS_EXP_USE_WAKE=1`](/experimental/yaml-builds#enable-yaml-builds). Lorsqu’aucun manifeste n’est trouvé, la commande exécute `go run .` dans le répertoire courant sans nécessiter cette activation. Un manifeste non valide ou illisible constitue une erreur ; il ne déclenche jamais la solution de repli. Si un manifeste existe alors que YAML est désactivé, la commande indique comment l’activer.

La référence générée du champ [`wails.yaml`](/reference/wails-yaml-fields) répertorie les types de champs et leurs valeurs par défaut. Avant qu’une version ne déclare ces plateformes prêtes, la prise en charge native des appareils mobiles et des bundles macOS doit être vérifiée sur l’hôte approprié.

## Contrat de la commande

`wails3 run` détecte le fichier `wails.yaml` le plus proche, résout la configuration d’une cible, effectue la compilation au moyen du pipeline existant avec mise en cache, puis lance le résultat. La détection remonte depuis le répertoire de travail, en recherchant un manifeste avant de s’arrêter à une limite `go.mod`. Le répertoire du manifeste constitue la racine du projet et le répertoire de travail de l’application de bureau.

Par défaut, la cible correspond au système d’exploitation et à l’architecture de l’hôte. Les cibles mobiles doivent être sélectionnées explicitement. La compilation croisée ne signifie pas qu’un exécutable de bureau peut s’exécuter sur l’hôte ; les cibles de bureau incompatibles échouent avant la compilation. Cela inclut les exécutables compilés avec `CGO_ENABLED=0` : la réussite d’une compilation croisée Go uniquement ne fournit aucun émulateur et ne modifie pas le contrat de lancement local. Utilisez [`wails3 build` pour la compilation croisée](/experimental/yaml-builds/#cross-compile-desktop-applications), puis exécutez ou installez les fichiers obtenus sur leur système cible. Pour Android, utilisez `--target android/arm64 --device <adb-serial>` ou `--target android/arm64 --emulator <avd-name>`. L’architecture sélectionnée doit correspondre à l’appareil ou à l’émulateur. Pour un simulateur iOS, utilisez `--target ios/arm64 --device <udid>` ; ajoutez `--destination device` pour un appareil physique. La planification et l’exécution iOS nécessitent macOS et Xcode.

Avec YAML, la compilation utilise les ressources intégrées et le chemin de compilation de production, qui fournit automatiquement le tag `production`. `wails3 dev` reste responsable du rechargement à chaud ; `dev.tags` ne contribue pas à `wails3 run`.

```shell
wails3 run
wails3 run --tags extra_tag
wails3 run -- --example-mode other
wails3 run --plan
wails3 run --skip-build
```

`--plan` résout et affiche la cible sélectionnée, les tags de compilation effectifs, l’artefact exécutable ou installable, les arguments et le répertoire de travail, sans compiler, installer ni lancer quoi que ce soit. Le plan répertorie les noms des clés de substitution d’environnement et omet leurs valeurs.

## Options de la CLI

| Option | Signification |
| --- | --- |
| `--tags <tags>` | Tags de compilation Go supplémentaires, séparés par des virgules. |
| `--skip-build` | Importez les exécutables de bureau existants depuis le répertoire de sortie YAML, puis lancez-les sans compilation. |
| `--target <os>/<arch>` | Sélectionne une cible ; par défaut, utilise le système d’exploitation et l’architecture de l’hôte. |
| `--plan` | Affiche les paramètres résolus sans compiler ni lancer. Les plans de bureau multiplateformes sont autorisés. |
| `--device <id>` | Numéro de série adb Android ou identifiant d’un simulateur ou appareil iOS. |
| `--emulator <name>` | Appareil virtuel Android à démarrer ou à réutiliser ; ne peut pas être combiné avec `--device`. |
| `--destination <kind>` | Destination iOS : `simulator` (par défaut) ou `device`. |
| `-- <args...>` | Remplace les arguments d’application configurés. Un `--` seul les efface. |
| `--appargs <string>` | Graphie de compatibilité : analyse une chaîne unique selon les règles de citation POSIX. Cette option ne peut pas être répétée ni combinée avec des arguments placés après `--`. |

Les sélecteurs mobiles nécessitent un manifeste et une cible mobile correspondante.

`--skip-build` nécessite un projet YAML actif et un exécutable de bureau compatible. Cette option ignore les opérations sur le frontend, les liaisons, la compilation, UPX et les hooks de build ; les arguments de lancement et l’environnement configurés continuent toutefois de s’appliquer. Les exécutables importés sont validés et copiés dans l’espace de travail généré par Wake sans modification de leurs fichiers sources. Les applications macOS avec interface graphique sont assemblées à partir de ces entrées et des ressources du paquet actuel. Consultez [les chemins et les exigences des entrées précompilées](/experimental/yaml-builds/#use-prebuilt-desktop-executables). Les tags de build ne peuvent pas modifier un exécutable précompilé ; `--tags` et les tags de build configurés n’ont donc aucun effet dans ce mode.

## Champs

Le mappage facultatif `run` fournit les valeurs de lancement par défaut communes. Un mappage facultatif `run` dans `targets["<os>"]` ou `targets["<os>/<arch>"]` fournit les valeurs par défaut propres à une plateforme ou à une architecture. Les paramètres propres à l’architecture ont priorité sur ceux définis pour l’ensemble du système d’exploitation, quel que soit l’ordre du mappage.

| Champ | Type | Valeur par défaut | Signification |
| --- | --- | --- | --- |
| `run.tags` | `list(string)` | `[]` | Tags Go supplémentaires utilisés uniquement par les compilations d’exécution. |
| `run.args` | `list(string)` | `[]` | Arguments transmis à l’application. |
| `run.environment` | `map(string)` | `{}` | Substitutions de l’environnement d’exécution, distinctes de l’environnement du compilateur. |
| `target["<os>"].run.tags` | `list(string)` | `[]` | Tags d’exécution supplémentaires pour la plateforme sélectionnée. |
| `target["<os>"].run.args` | `list(string)` | héritée | Remplace les arguments communs lorsqu’il est défini explicitement ; `[]` les efface. |
| `target["<os>"].run.environment` | `map(string)` | `{}` | Remplace par clé les valeurs communes de l’environnement d’exécution. |

Les arguments sont transmis sous forme de vecteur d’arguments, sans expansion par l’interpréteur de commandes. Les valeurs d’environnement sont des chaînes littérales conformément au contrat existant de YAML, limité aux valeurs littérales.

## Exemple propre à une plateforme

Ce fragment accompagne la configuration normale du projet. `example_private_api` est un substitut du véritable tag d’API privée requis par l’exemple.

```yaml
build:
  tags: [devtools]
run:
  args: [--example-mode]
targets:
  darwin:
    tags: [example_private_api]
    run:
      environment:
        EXAMPLE_WINDOW_MODE: native
  linux:
    run:
      environment:
        GDK_BACKEND: wayland
  windows:
    run:
      args: [--example-mode, --windows-option]
```

Le paramètre d’environnement Linux est propre à cet exemple et nécessite une session Wayland ; il ne s’agit pas d’une valeur par défaut de Wails. Le tag d’API privée est sélectionné uniquement sous macOS. Le paramètre `target.tags` existant s’applique aux compilations de production, y compris aux compilations d’exécution. `target.run.tags` réserve un tag supplémentaire aux compilations d’exécution.

Une substitution propre à une architecture peut affiner les valeurs par défaut définies pour l’ensemble du système d’exploitation :

```yaml
targets:
  linux/arm64:
    run:
      tags: [example_arm64]
      args: []
```

Combinée au mappage Linux précédent, cette configuration conserve `GDK_BACKEND`, ajoute le tag propre à l’architecture et efface les arguments d’application communs.

## Règles de résolution

| Paramètre | Résolution, de la priorité la plus faible à la plus élevée |
| --- | --- |
| Tags | Combine `build.tags`, le paramètre `target.tags` sélectionné, le paramètre automatique `production`, `run.tags`, le paramètre `target.run.tags` sélectionné et le paramètre CLI `--tags` ; conserve la première occurrence de chaque tag. |
| Environnement d’exécution | Hérite de l’interpréteur de commandes qui effectue le lancement, applique ensuite `run.environment`, puis le paramètre `target.run.environment` sélectionné. Une valeur vide définit une chaîne vide. |
| Arguments de l’application | Utilise le paramètre commun `run.args`, remplacé par les arguments de cible définis explicitement, eux-mêmes remplacés par les arguments explicites placés après `--`. Un `--` seul efface les arguments configurés. |

Dans la cible sélectionnée, les paramètres définis pour l’ensemble du système d’exploitation sont résolus avant ceux propres à l’architecture : les tags s’accumulent, les valeurs d’environnement sont fusionnées par clé et les arguments fournis explicitement remplacent la liste moins spécifique. Les compilations d’exécution fournissent également `android` ou `ios` pour la cible mobile correspondante.

La génération des liaisons et la compilation Go reçoivent les mêmes tags résolus. Les clés du cache de compilation incluent ces tags. La modification des seuls arguments de l’application ou de l’environnement d’exécution n’invalide pas la compilation ; le lancement suivant utilise les nouvelles valeurs.

Les mappages de plateforme personnalisent une plateforme ; leur présence ne rend pas les autres plateformes non prises en charge. `project.supported_platforms` limite les plateformes acceptées par les compilations et `run` ; une liste omise ou vide autorise toutes les plateformes. Par exemple, `supported_platforms: [darwin]` dans `project` provoque l’échec anticipé sous Linux d’un exemple réservé à macOS.

## Comportement de lancement selon la plateforme

| Plateforme | Comportement de lancement |
| --- | --- |
| Linux et Windows | Lancez l’exécutable avec les entrées/sorties standard connectées, la gestion de l’arrêt et la propagation du code de sortie de l’application. |
| macOS | Assemblez un paquet d’application et exécutez son binaire `Contents/MacOS` avec les entrées/sorties standard connectées. Enregistrez le paquet auprès de Launch Services pour les schémas d’URL et les associations de fichiers. La signature et la validation du cycle de vie natif nécessitent une vérification sous macOS. |
| Android | Générez un APK, sélectionnez un appareil ou un émulateur, puis installez et lancez l’application. Signalez séparément les échecs de génération, d’installation et de lancement. |
| iOS | Générez un paquet d’application, sélectionnez un simulateur ou un appareil, remplissez les exigences de signature et de provisionnement, puis procédez à l’installation et au lancement. Les builds destinés aux appareils nécessitent les identifiants appropriés et la chaîne d’outils Apple. |

Android rejette les arguments d’application et les substitutions de l’environnement d’exécution avant le déploiement. Lorsque Android est également pris en charge, placez les paramètres partagés réservés aux environnements de bureau dans les mappages des cibles de bureau. Lors d’une exécution sous Android, les journaux de l’application sont suivis et l’application est arrêtée à la fin de ce suivi. iOS transmet les arguments au lanceur sélectionné. Les simulateurs reçoivent les substitutions d’environnement par l’intermédiaire de variables `SIMCTL_CHILD_`, tandis que les appareils physiques reçoivent une table JSON par l’intermédiaire de `devicectl --environment-variables`. Sur les plateformes mobiles, le code de sortie décrit le résultat du lanceur ou de la surveillance des journaux, et non un code de sortie d’application portable.

## Exemples en Go uniquement et solution de repli

Un exemple en Go uniquement, avec des ressources préexistantes ou incorporées, désactive les étapes du frontend :

```yaml
frontend:
  disabled: true
```

Cela ignore l’installation des dépendances du frontend, la génération des bindings et la génération du frontend. Les exemples nécessitant des classes de binding générées peuvent à la place définir `interfaces: false` sous `frontend.bindings`.

La solution de repli sans manifeste exécute `go run [-tags <tags>] . [application arguments]` sans ajouter `production` ni `devtools`. Elle suit le comportement propre à Go en matière de code de sortie. Les arguments placés après `--` sont transmis à l’application, y compris les options telles que `--help`. Les valeurs d’exécution affichées par `--plan` se limitent aux noms des clés d’environnement ; les valeurs ne sont pas affichées. Le plan de repli affiche la commande Go et le répertoire de travail.

## Arguments de l’application en développement

Avec un projet YAML actif, `wails3 dev` accepte des arguments d’application indépendamment des valeurs par défaut de production de `run`. `dev.args` fournit les valeurs par défaut partagées ; un mappage `dev` imbriqué dans `targets["<os>"]` ou `targets["<os>/<arch>"]` les remplace pour cette cible. Les paramètres propres à l’architecture prévalent sur ceux applicables à l’ensemble du système d’exploitation, quel que soit l’ordre du mappage. Une valeur `args: []` explicite efface les arguments hérités.

```yaml
dev:
  args: [--config-path, config/testing.yaml]
targets:
  windows:
    dev:
      args: [--config-path, config/testing-windows.yaml]
```

`run.args` et `run.tags` ne contribuent pas aux lancements en développement. La politique de build `dev.tags` existante reste inchangée. L’application interprète les chemins relatifs à partir du répertoire du manifeste.

```shell
wails3 dev
wails3 dev -- --config-path "profile with spaces.yaml"
wails3 dev --appargs='--config-path "profile with spaces.yaml"'
wails3 dev --
wails3 dev --plan
```

Les arguments placés après `--` remplacent les arguments configurés, même lorsque la liste est vide. `--help` placé après ce délimiteur appartient à l’application. `--appargs=""` efface également les valeurs par défaut. L’option de compatibilité accepte une seule chaîne entre guillemets ; son utilisation avec `--`, avec des arguments d’application positionnels ou avec un second `--appargs` constitue une erreur. Une mise entre guillemets incorrecte est rejetée avant le démarrage d’une session.

L’analyseur de compatibilité applique les mêmes règles de mise entre guillemets de style POSIX sur tous les hôtes. Il n’exécute aucune commande shell et ne développe aucune variable ; le shell appelant peut toutefois avoir déjà développé le texte avant de le transmettre à Wails. Les arguments placés après `--` sont déjà des valeurs distinctes et ne font l’objet d’aucune tokenisation supplémentaire. Il s’agit de la forme recommandée pour les chemins Windows, par exemple dans PowerShell :

```powershell
wails3 dev -- --config-path 'C:\Users\Example User\testing.yaml'
```

La session conserve le vecteur d’arguments résolu lors des reconstructions du code source. Si vous modifiez uniquement `dev.args` ou les arguments de la cible sélectionnée, le backend redémarre avec l’exécutable mis en cache, sans redémarrer le frontend. Les substitutions de la CLI restent actives après le rechargement du manifeste. En cas d’échec d’un remplacement, l’exécutable précédent et ses arguments précédents sont conservés ou restaurés. Les arguments de l’application ne sont pas des options du compilateur et ne sont pas inclus dans les clés du cache de compilation. `dev --plan` affiche les arguments effectifs avec le plan de build fini, sans effectuer de lancement.

Linux, Windows et macOS transmettent directement le vecteur à l’exécutable. iOS le transmet par l’intermédiaire du lanceur du simulateur ou de l’appareil. Android rejette les arguments d’application de développement non vides avant la détection des appareils ou le build, car aucun canal équivalent de transmission des arguments n’est pris en charge. Vous pouvez placer les valeurs par défaut réservées aux environnements de bureau dans les mappages des cibles de bureau, ou les effacer avec `args: []` sous `targets.android.dev`.

Ces options de développement nécessitent un `wails.yaml` actif ainsi que l’activation explicite dans le fichier YAML. Elles ne réécrivent pas les commandes shell du Taskfile ou de l’observateur dans l’ancien chemin de développement.
