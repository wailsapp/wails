---
title: "Wake"
description: "Créez votre premier projet Wake, suivez l’assistant de création de projet, compilez et exécutez le projet, puis migrez un projet Taskfile existant."
slug: "experimental/wake"
sourcePath: "experimental/wake.md"
---

Wake est le système de compilation expérimental de Wails. Un fichier `wails.yaml` décrit votre application, son frontend, les plateformes cibles, les paquets et les paramètres de signature. Wake utilise cette configuration pour planifier et exécuter les compilations, mettre en cache les tâches inchangées et indiquer les fichiers produits.

Dans ce tutoriel, vous allez créer une petite application de bureau, la configurer dans le navigateur, la compiler pour votre propre ordinateur et l’exécuter en mode développement. Commencez avec une seule cible locale et une compilation non signée. Ajoutez la compilation croisée et les formats de distribution une fois cette première compilation fonctionnelle.

Vous avez déjà une application ? Consultez la section [Migrer un projet existant](#migrer-un-projet-existant). Pour obtenir des informations sur chaque champ et les workflows avancés, consultez le [guide de compilation YAML](/experimental/yaml-builds/) et la [référence du manifeste](/reference/wails-yaml-fields/).

## À quoi s’attendre avec cette fonctionnalité expérimentale

@note{type="caution" title="Utilisez un projet de test ou une branche distincte"}
Wake est facultatif, et sa configuration, ses commandes et son comportement peuvent changer d’une préversion à l’autre. Il n’offre pas encore les garanties de compatibilité d’un système de compilation stable. Conservez une révision fonctionnelle connue de votre projet et consignez les versions du CLI et du framework que vous testez. Testez les paquets sur leurs systèmes de destination avant de les distribuer.
@end

Wake et le backend **Pure Go** sont deux choix indépendants. Activer Wake n’active pas automatiquement Pure Go. Cette leçon utilise le backend standard ; les compilations d’interfaces graphiques pour macOS et Linux nécessitent donc les outils de développement natifs habituels.

L’assistant vérifie les prérequis et enregistre la configuration. Cet enregistrement ne lance aucune compilation et ne prouve pas que toutes les cibles sélectionnées sont prêtes. Vous pourrez installer ultérieurement les outils manquants ou configurer les identifiants requis, par exemple sur un hôte d’intégration continue. Vous devez corriger toute combinaison non valide de cible et de paquet avant l’enregistrement.

La migration est **à sens unique**. Wake n’exécute pas les Taskfiles, ne développe pas les modèles Task et ne reconvertit pas le YAML en Taskfiles. Les projets non migrés continuent d’utiliser l’ancien environnement d’exécution Task distinct. Un fichier `wails.yaml` actif ne se rabat jamais silencieusement sur celui-ci.

## Avant de commencer

Vous devez disposer de Go, de Node.js avec npm pour le frontend de démarrage et des [dépendances propres à la plateforme](/getting-started/installation/#platform-specific-dependencies) de votre ordinateur. Placez-vous dans le répertoire où vous souhaitez créer le dossier du nouveau projet. L’assistant nécessite également un navigateur.

Utilisez une révision du CLI et du framework qui inclut cette préversion. L’installation du CLI `v3.0.0` stable n’ajoute pas les écrans de création de projet Wake décrits ici. Les commandes suivantes épinglent la préversion issue de la [PR nº 6531](https://github.com/wailsapp/wails/pull/6531) :

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@b1daf95f3
wails3 version
```

Vérifiez que le répertoire des exécutables Go figure dans `PATH` ; le [guide d’installation](/getting-started/installation/) explique comment procéder. Plus loin dans ce tutoriel, vous épinglerez la dépendance Go de l’application à la même révision. Lorsque vous testez une préversion plus récente, mettez à jour simultanément les deux arguments de révision.

### Activer Wake dans votre terminal

Dans Bash ou Zsh :

```bash
export WAILS_EXP_USE_WAKE=1
```

Dans PowerShell :

```powershell
$env:WAILS_EXP_USE_WAKE = '1'
```

Gardez ce terminal ouvert pendant le reste du tutoriel. Définissez à nouveau la variable dans tout nouveau terminal ou toute nouvelle tâche d’intégration continue qui exécute des commandes Wake. Sa **présence** active Wake : lui attribuer `false` ou une chaîne vide ne le désactive pas. L’ancien nom `WAILS_USE_WAKE` est ignoré.

## Créer l’application

Exécutez l’assistant de création de projet :

```sh
wails3 init -ui
```

L’assistant de projet s’ouvre dans le navigateur. Lorsque Wake est activé, il regroupe les choix relatifs au code source et à la compilation dans un même parcours. Rien n’est créé avant l’utilisation du dernier bouton **Créer le projet**. `wails3 setup` est un assistant distinct consacré aux préférences globales, aux outils et aux identifiants de signature ; il ne configure pas le profil Wake de ce projet.

| Écran | Choix pour ce tutoriel | Étape suivante |
| --- | --- | --- |
| **Projet — Créer un projet Wails** | Saisissez `wake-demo`, puis cliquez sur **Continuer**. | Le projet sera créé dans un nouveau dossier `wake-demo` sous le répertoire affiché. |
| **Framework — Choisir un framework** | Sélectionnez **Vanilla**, puis cliquez sur **Continuer**. | Vous obtenez un petit frontend sans avoir à apprendre les particularités d’un framework. |
| **Langage — Préférence de langage** | Sélectionnez **JavaScript**, puis cliquez sur **Continuer**. | L’écran **Liaisons**, réservé à TypeScript, est ignoré. |
| **Détails — Détails du projet** | Utilisez `Wake Demo` comme nom de produit, `com.example.wakedemo` comme identifiant de bundle et `0.1.0` comme version. Cliquez sur **Continuer**. | Vos métadonnées restent dans l’assistant pendant que vous choisissez les paramètres de compilation. |

Pour cette leçon, n’utilisez pas **Utiliser les valeurs globales par défaut** : cette option ignore les choix du framework et du langage. Si vous choisissez TypeScript dans un autre projet, l’écran **Liaisons** permet de choisir des interfaces ou des classes pour les liaisons de données Go générées.

## Choisir les paramètres de compilation avant la création

Poursuivez dans la même fenêtre de navigateur. Les captures d’écran ci-dessous illustrent cette procédure sous Linux. La cible correspondant à votre hôte et les outils requis seront différents sous Windows ou macOS.

### Application — Que créez-vous ?

Conservez **Application de bureau** et choisissez **Backend standard**. Cliquez sur **Continuer**.

![Écran Application de l’assistant de création de projet, avec Application de bureau et Backend standard sélectionnés.](/assets/screenshots/wake/application.png)

Cette leçon crée une interface graphique dotée d’un frontend web. **Pure Go (expérimental)** est une option distincte. Elle supprime l’exigence d’un compilateur C pour les compilations de bureau prises en charge, à condition que vos dépendances puissent également être compilées sans CGo. Elle ne supprime ni les bibliothèques d’exécution de la cible, ni tous les outils de création de paquets et de signature.

Le backend est un **paramètre qui s’applique à l’ensemble du projet**. Les autres profils qui héritent de ces paramètres utilisent eux aussi le backend sélectionné.

### Cibles — Où l’application doit-elle s’exécuter ?

Ne conservez que la cible correspondant à l’ordinateur que vous utilisez, puis cliquez sur **Continuer**. L’assistant sélectionne initialement le système d’exploitation et l’architecture de l’hôte.

![Écran Cibles avec uniquement linux/amd64 sélectionné pour l’hôte utilisé dans le tutoriel.](/assets/screenshots/wake/targets.png)

Le nom d’une cible associe un système d’exploitation et une architecture de processeur : par exemple, `windows/amd64`, `darwin/arm64` pour un Mac Apple Silicon, ou `linux/amd64`. `darwin` désigne macOS. Sélectionner une autre cible demande la création d’un exécutable pour cette machine ; cela ne rend pas cet exécutable utilisable sur votre ordinateur actuel.

Réservez pour plus tard les architectures supplémentaires, les compilations universelles macOS et les cibles mobiles expérimentales. Vous pourrez les ajouter à `wails.yaml` après la première compilation locale.

### Paquets — Choisir les formats de publication

Laissez les cases des formats de paquet décochées et désactivez **Signer les fichiers de publication**. Cliquez sur **Continuer**.

![Écran Paquets sans format de paquet Linux supplémentaire sélectionné et avec la signature désactivée.](/assets/screenshots/wake/packages.png)

Sans format supplémentaire sélectionné, Wake produit la sortie par défaut de la cible : un exécutable de bureau ou, pour une interface graphique macOS, un paquet d’application macOS. Cela suffit pour vérifier l’application localement. La sélection de formats ajoute, lorsqu’ils sont pris en charge, des sorties de distribution telles que DMG, MSIX, NSIS, DEB ou AppImage.

La signature est facultative pour cette première compilation locale. La réussite d’une compilation non signée ne signifie pas que l’application est prête à être distribuée publiquement. Pour un projet signé, configurez des identifiants réutilisables avec `wails3 setup` avant la création, ou configurez ensuite la signature dans `wails.yaml`. L’assistant de création utilise les identifiants que vous avez configurés ; il ne modifie pas vos préférences globales de signature.

Wake peut utiliser des fichiers de certificat configurés hors du projet sans copier de clés privées ni de mots de passe dans `wails.yaml`. Les paramètres de signature explicites du projet sont prioritaires. Avant de préparer une version publique, suivez les [instructions de configuration de la signature](/getting-started/setup/) et le [guide de création de paquets et de signature](/experimental/yaml-builds/#package-and-sign-for-distribution). La signature mobile se configure dans le fichier YAML.

### Outils — Outils pour votre compilation

Lisez les prérequis correspondant à vos choix. La liste provient du plan de compilation de Wake ; elle ne vous demande pas d’installer tous les outils sur chaque plateforme.

![Écran Outils présentant les prérequis déterminés pour la compilation locale des sources sous Linux.](/assets/screenshots/wake/tools.png)

Pour cette leçon, prévoyez Go et les outils frontend, ainsi que les outils de développement natifs requis par le backend d’interface graphique standard de votre système hôte. Docker n’est pas systématiquement requis. Une compilation croisée avec CGo peut nécessiter un compilateur et les bibliothèques de la cible, ou une chaîne d’outils prise en charge dans un conteneur.

Avant de tenter la compilation locale, corrigez les lignes marquées **Non prêt**, **Configuration requise** ou **Mise à jour requise**. Suivez les instructions affichées lorsqu’elles sont disponibles, ou consultez le [guide d’installation de la plateforme](/getting-started/installation/), puis cliquez sur **Revérifier les prérequis**. Si une installation a modifié la variable `PATH` du terminal, redémarrez `wails3 init -ui` depuis un terminal mis à jour et répétez vos sélections.

L’assistant peut créer un projet valide même si cet ordinateur ne dispose pas des outils ou des identifiants de signature nécessaires. Cela est utile pour configurer l’intégration continue, mais, dans ce tutoriel, vous devez obtenir une compilation exécutable sur cet ordinateur. Toute configuration non valide doit être corrigée avant la création.

### Créer le projet

Vérifiez vos choix : application de bureau, backend standard, cible correspondant à votre système hôte et signature désactivée. Dans **Outils**, cliquez sur **Créer le projet** et attendez que le message **Choix du projet acceptés** s’affiche. Revenez dans le terminal et attendez que la CLI termine la création du projet ; la confirmation dans le navigateur ne fait qu’accepter vos choix.

![Écran « Choix du projet acceptés » affichant le chemin prévu de l’application et les commandes à utiliser une fois la création terminée.](/assets/screenshots/wake/save.png)

L’assistant crée les sources et enregistre les choix de cible, de paquet et de signature sous `profiles.setup` dans `wails.yaml`. Il définit également le type d’application, l’activation du frontend et le backend commun à l’ensemble du projet. Vous pouvez fermer le navigateur et revenir au terminal. Attendez la fin de la création du projet, puis exécutez :

```sh
cd wake-demo
go get github.com/wailsapp/wails/v3@b1daf95f3
```

Vous devez maintenant disposer de `main.go`, d’un module Go, d’un répertoire `frontend` et d’un fichier `wails.yaml` natif. Le modèle intégré utilise directement Wake ; vous n’avez pas besoin de migrer ce nouveau projet. Un modèle communautaire peut contenir des tâches héritées personnalisées et nécessiter à la place l’examen de migration décrit ci-dessous.

L’écran de fin affiche `wails3 build setup` et `wails3 dev`. Le nom du profil est important : `wails3 build setup` utilise les sélections que vous venez d’enregistrer. La commande `wails3 build` seule ne sélectionne pas automatiquement le profil `setup`. Le nom du profil ne signifie pas que vous devez exécuter `wails3 setup`.

## Compiler et exécuter votre première application

Exécutez la commande affichée par l’assistant :

```sh
wails3 build setup
```

La première compilation installe les dépendances frontend, génère les liaisons et compile l’application. Les compilations suivantes peuvent réutiliser les résultats inchangés. Attendez le message **compilation réussie** et consultez les chemins de sortie affichés par l’interface en ligne de commande. Le répertoire de sortie par défaut est `bin` ; une interface graphique macOS est assemblée dans un paquet `.app`.

Démarrez le mode développement sur votre ordinateur :

```sh
wails3 dev
```

La fenêtre de l’application de démarrage doit s’afficher. Saisissez votre nom et essayez le message d’accueil. Modifiez le texte indicatif du champ de saisie du nom dans `frontend/index.html`, puis enregistrez le fichier ; le frontend doit s’actualiser. Modifiez le message d’accueil dans `greetservice.go` pour déclencher une recompilation Go. Lorsque vous avez terminé, appuyez sur **Ctrl+C** dans le terminal.

Le développement s’effectue localement sur la machine hôte. Un profil de publication peut contenir plusieurs cibles de compilation croisée sans changer la machine capable d’exécuter la fenêtre de développement. Pour compiler et lancer localement une application de production, utilisez `wails3 run` ; consultez la [référence de la commande run](/reference/wails-yaml-run/).

Vous disposez maintenant d’un projet source, d’un profil de compilation vérifié et d’une boucle de développement locale fonctionnelle. Validez les sources et `wails.yaml` dans le système de gestion de versions. Les fichiers de travail `.wails/` générés ne constituent pas la source de votre configuration.

## Étapes suivantes

| Objectif | Prise en charge par Wake | Points à vérifier en premier |
| --- | --- | --- |
| Compiler pour un autre système de bureau | Cibles Windows, macOS et Linux, y compris les deux architectures macOS et une sortie universelle. | Go pur et CGo ont des exigences de compilation différentes. Testez sur le système d’exploitation et l’architecture de destination. |
| Créer des programmes d’installation ou des paquets | MSIX, NSIS, application macOS/DMG, DEB, RPM, Arch, Void et AppImage. | NSIS nécessite `makensis`. AppImage nécessite un environnement d’exécution cible ; la création croisée de paquets d’interface graphique nécessite également une arborescence d’environnement d’exécution Linux fournie. |
| Signer une version | Signature portable par clé exportée pour Windows, macOS et DEB/RPM/Arch ; les magasins d’identités natifs restent disponibles. | Chaque identité, format et système hôte présente des limites. AppImage et Void ne disposent d’aucun processus de signature intégré. |
| Notarier une application macOS | Envoi par clé d’API d’équipe et agrafage du ticket, ou profil de trousseau macOS. | Une signature Developer ID appropriée et un horodatage restent obligatoires. La réussite des vérifications natives par des signatures de test générées avec Go ne prouve pas qu’Apple acceptera votre version. |
| Compiler plusieurs programmes | Un programme principal avec interface graphique, en ligne de commande ou serveur, accompagné de programmes complémentaires nommés. | Déclarez les paquets et les types des programmes complémentaires dans le fichier YAML. L’assistant ne génère pas leur code source. |
| Créer des paquets à partir de fichiers compilés ailleurs | `wails3 build setup --skip-build`. | Indiquez dans `build.output` des noms d’exécutables et des architectures correspondants. Cette option est un choix ponctuel de la commande, et non un paramètre permanent du profil. |
| Ajouter l’automatisation des publications | Profils nommés, hooks de cycle de vie limités, mise en cache et création parallèle de paquets. | Utilisez du YAML littéral et des scripts déclarés appartenant au projet. Le comportement arbitraire d’un Taskfile n’est pas un mécanisme d’extension de Wake. |

Modifiez `profiles.setup`, les autres profils nommés et les paramètres avancés dans `wails.yaml`. Utilisez `wails3 setup` pour modifier les préférences et identifiants partagés. Consultez le [guide de compilation YAML](/experimental/yaml-builds/) pour obtenir des exemples complets et connaître les prérequis des paquets. Les applications avec interface graphique ont toujours besoin des bibliothèques système de la destination ; Go pur n’en fait pas des applications entièrement statiques.

iOS et Android restent expérimentaux, nécessitent leurs SDK de plateforme et CGo, et sortent du cadre de cette première leçon consacrée aux applications de bureau. L’assistant ne rend pas les compilations mobiles portables entre des systèmes hôtes arbitraires.

## Migrer un projet existant

Travaillez sur une branche distincte ou une copie de l’arbre de travail d’un projet **Wails v3**. Il s’agit d’une migration de la configuration de compilation, et non d’une mise à niveau depuis Wails v2. Conservez le résultat de compilation et les tests existants afin de pouvoir comparer le comportement après l’activation.

@note{type="note" title="Commencer par une migration vérifiée"}
Les nouveaux projets générés à partir des modèles historiques intégrés actuels peuvent être migrés avant la compilation de leur frontend. Wake reconnaît leurs sélecteurs de tâches standard et leur cycle de développement sans exécuter de Taskfiles ni évaluer d’expressions personnalisées.

Les tâches personnalisées, les substitutions de variables et les politiques de développement doivent toujours être examinées. Si l’analyse ne peut pas représenter votre projet, conservez l’ancien système de compilation et signalez les diagnostics. Ne supprimez pas un comportement dont vous avez besoin dans le seul but de réussir l’activation. Le tutoriel ci-dessus consacré aux nouveaux projets ne nécessite aucune migration.
@end

### 1. Analyser sans modifier le système de compilation actif

Activez Wake et utilisez la CLI de préversion indiquée au début de cette page. Depuis la racine du projet existant, exécutez :

```sh
wails3 migrate --dry-run
```

Consultez les paramètres pris en charge et les blocages signalés. L’analyseur lit les déclarations historiques, mais n’exécute pas leurs commandes. Une simulation ne crée pas de manifeste actif. Pour obtenir des diagnostics lisibles par une machine, utilisez `wails3 migrate --json` ; utilisez `--dry-run --json` si vous souhaitez également éviter l’écriture d’un brouillon.

### 2. Créer et examiner la proposition inactive

```sh
wails3 migrate
```

Cette opération crée **`wails.migrated.yaml`**, et non `wails.yaml`. Examinez ce fichier en parallèle de vos Taskfiles existants et de `build/config.yml`. Vérifiez les métadonnées du projet, les commandes du frontend, les architectures cibles, le contenu des paquets, les ressources et les paramètres de signature. Conservez les secrets dans des variables d’environnement ou dans votre gestionnaire d’identifiants.

Les tâches intégrées connues peuvent être converties en paramètres natifs. Les tâches de cycle de vie prudentes, sans argument, qui appellent un script appartenant au projet peuvent devenir des hooks. Tout autre comportement accessible provenant d’une tâche personnalisée constitue un blocage : Wake ne l’ignore pas silencieusement et ne l’exécute pas via Task comme solution de repli.

S’il existe déjà un brouillon, la migration n’écrase pas les modifications que vous avez examinées. Conservez-en une copie avant de le remplacer ou de le régénérer délibérément. Ne renommez pas un brouillon bloqué en `wails.yaml` pour contourner les contrôles d’activation. Réduisez ou remplacez le comportement historique non pris en charge, puis relancez l’analyse avant l’activation. Si le projet ne peut pas encore être représenté, conservez son ancien système de compilation et signalez cette lacune.

### 3. Activer explicitement

Après avoir examiné le brouillon et résolu les blocages, exécutez :

```sh
wails3 migrate --activate --dry-run
wails3 migrate --activate
```

La première commande vérifie l’activation sans la publier. La seconde vérifie à nouveau les déclarations historiques et valide le brouillon avant de créer `wails.yaml` de manière atomique ; si l’activation réussit, elle supprime le fichier inactif `wails.migrated.yaml`. Elle ne modifie ni les Taskfiles, ni `build/config.yml`, ni les scripts, ni les ressources. À partir de ce moment, le manifeste actif contrôle les compilations Wake. Les anciens fichiers ne constituent pas une seconde source de configuration.

Certains projets plus anciens possèdent déjà un fichier `wails.yaml` contenant uniquement des métadonnées, avec une section `info:` au niveau supérieur. Bien qu’il ne s’agisse pas d’un manifeste Wake natif, il occupe tout de même la destination d’activation. Conservez ce fichier historique et déplacez-le délibérément avant l’activation ; l’outil de migration ne l’écrasera ni ne le renommera à votre place.

### 4. Vérifier l’application migrée

```sh
wails3 config check
wails3 build
wails3 dev
```

Comparez l’application et le contenu des paquets avec la compilation précédente, y compris votre comportement de publication personnalisé. Exécutez les tests de votre projet. Vérifiez explicitement chaque profil migré, par exemple `wails3 config check release` et `wails3 build release` lorsque le brouillon contient un profil `release`.

Modifiez directement les profils migrés dans `wails.yaml`. `wails3 init -ui` crée de nouveaux projets ; `wails3 setup` gère les préférences globales. Aucun des deux ne permet de modifier le profil de publication d’un projet existant.

### Si vous devez interrompre l’expérimentation

La suppression de `WAILS_EXP_USE_WAKE` désactive l’activation volontaire, mais la présence d’un fichier `wails.yaml` actif bloque toujours le routage automatique vers les Taskfiles. Un YAML non valide provoque également un arrêt avec une erreur. Il n’existe aucune commande de migration inverse.

Pour revenir à l’ancien flux de travail, restaurez dans le contrôle de version la révision complète du projet antérieure à la migration, y compris l’absence du manifeste actif, puis utilisez la CLI et les dépendances compatibles avec cette révision. Au préalable, conservez toute nouvelle modification du code source dont vous avez besoin. L’utilisation d’une branche d’essai distincte facilite considérablement cette restauration délibérée.

## Diagnostiquer un problème et transmettre vos commentaires

| Symptôme | Première vérification |
| --- | --- |
| La création du projet ignore les étapes Application, Targets, Packages et Tools | Définissez `WAILS_EXP_USE_WAKE` dans le terminal de lancement et vérifiez quel exécutable `wails3` se trouve dans `PATH`. La commande globale `wails3 setup` conserve toujours son flux de préférences. |
| Un modèle communautaire doit être migré | Examinez sa configuration historique à l’aide des commandes de migration avant l’activation. |
| Une compilation enregistrée signale des outils manquants | La création d’un projet enregistre l’intention. Installez les prérequis de l’hôte signalés avant la compilation ; la création n’installe pas tous les outils sélectionnés. |
| Un exécutable produit par compilation croisée ne démarre pas ici | Exécutez-le sur le système d’exploitation et l’architecture sélectionnés ; utilisez la cible de l’hôte pour le développement local. |
| Le mode précompilé ne trouve pas de binaire | Vérifiez les [noms de fichiers et chemins précompilés](/experimental/yaml-builds/#use-prebuilt-desktop-executables), y compris les sous-répertoires pour plusieurs cibles. |
| Une modification du manifeste échoue | Exécutez `wails3 config check` pour le profil concerné et suivez le diagnostic du champ. Wake ne se rabat pas sur les Taskfiles. |

Inspectez une compilation avant de l’exécuter avec `wails3 build setup --plan --json`. `wails3 config check` valide la configuration, tandis que `wails3 eject` écrit un instantané `wails.ejected.yaml` inactif sans modifier le manifeste actif. Examinez la sortie de la commande avant de la partager : les chemins du projet et les emplacements des identifiants peuvent révéler votre environnement local.

Faites part de vos commentaires dans les [discussions sur les expérimentations](https://github.com/wailsapp/wails/discussions/categories/experiments) ou [signalez un bug reproductible](https://github.com/wailsapp/wails/issues). Indiquez les révisions de votre CLI et du framework, le système d’exploitation et l’architecture de l’hôte, la cible sélectionnée, la commande exacte, si vous avez sélectionné Pure Go, ainsi que le diagnostic pertinent. Précisez si l’échec s’est produit lors de la compilation, de la création du paquet, de la signature ou de l’exécution de l’application. Si cela est utile, joignez un manifeste minimal dont les secrets ont été supprimés ainsi que la sortie de `wails3 doctor`.
