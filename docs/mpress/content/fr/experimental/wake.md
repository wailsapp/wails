---
title: "Wake"
description: "Créez votre premier projet Wake, suivez l’assistant de configuration, compilez et exécutez le projet, puis migrez un projet Taskfile existant."
slug: "experimental/wake"
sourcePath: "experimental/wake.md"
---

Wake est le système de compilation expérimental de Wails. Un fichier `wails.yaml` décrit votre application, son frontend, les plateformes cibles, les paquets et les paramètres de signature. Wake utilise cette configuration pour planifier et exécuter les compilations, mettre en cache les tâches inchangées et indiquer les fichiers produits.

Dans ce tutoriel, vous allez créer une petite application de bureau, la configurer dans le navigateur, la compiler pour votre ordinateur et l’exécuter en mode développement. Commencez avec une seule cible locale et une compilation non signée. Ajoutez la compilation croisée et la distribution une fois cette première compilation fonctionnelle.

Vous avez déjà une application ? Accédez à [Migrer un projet existant](#migrer-un-projet-existant). Pour en savoir plus sur les différents champs et les processus avancés, consultez le [guide de compilation YAML](/experimental/yaml-builds/) et la [référence du manifeste](/reference/wails-yaml-fields/).

## À quoi vous attendre avec cette expérimentation

@note{type="caution" title="Utilisez un projet de test ou une branche distincte"}
Wake est facultatif, et sa configuration, ses commandes et son comportement peuvent changer d’une préversion à l’autre. Il n’offre pas encore les garanties de compatibilité d’un système de compilation stable. Conservez une révision fonctionnelle connue de votre projet et notez les versions de l’interface en ligne de commande et du framework que vous testez. Testez les paquets sur leurs systèmes de destination avant de les distribuer.
@end

Wake et le backend **Pure Go** sont deux choix indépendants. L’activation de Wake n’active pas automatiquement Pure Go. Cette leçon utilise le backend standard ; les compilations d’interfaces graphiques pour macOS et Linux nécessitent donc leurs outils de développement natifs habituels.

L’assistant vérifie les prérequis et enregistre la configuration. L’enregistrement ne lance pas de compilation et ne prouve pas que chaque cible sélectionnée est prête. Vous pourrez fournir ultérieurement les outils ou identifiants manquants, par exemple sur un hôte d’intégration continue. Vous devez corriger les combinaisons de cibles et de paquets non valides avant l’enregistrement.

La migration est **irréversible**. Wake n’exécute pas les Taskfiles, ne développe pas les modèles Task et ne reconvertit pas le YAML en Taskfiles. Les projets non migrés continuent d’utiliser l’ancien environnement d’exécution Task distinct. Un fichier `wails.yaml` actif ne s’y rabat jamais silencieusement.

## Avant de commencer

Vous avez besoin de Go, de Node.js avec npm pour le frontend de démarrage, ainsi que des [dépendances propres à la plateforme](/getting-started/installation/#platform-specific-dependencies) de votre ordinateur. Placez-vous dans le répertoire où vous souhaitez créer le dossier du nouveau projet. L’assistant nécessite également un navigateur.

Utilisez des révisions de l’interface en ligne de commande et du framework qui incluent cette préversion. L’installation de l’interface en ligne de commande `v3.0.0` stable n’ajoute pas les écrans de configuration de Wake décrits ici. Les commandes suivantes épinglent la préversion provenant de la [PR nº 6531](https://github.com/wailsapp/wails/pull/6531) :

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@0243da963
wails3 version
```

Vérifiez que le répertoire contenant l’exécutable Go figure dans `PATH` ; le [guide d’installation](/getting-started/installation/) explique comment procéder. Plus loin dans ce tutoriel, vous épinglerez la dépendance Go de l’application à la même révision. Lorsque vous testez une préversion plus récente, mettez à jour simultanément les deux arguments de révision.

### Activer Wake dans votre terminal

Dans Bash ou Zsh :

```bash
export WAILS_EXP_USE_WAKE=1
```

Dans PowerShell :

```powershell
$env:WAILS_EXP_USE_WAKE = '1'
```

Gardez ce terminal ouvert pendant tout le reste du tutoriel. Redéfinissez la variable dans tout nouveau terminal ou toute nouvelle tâche d’intégration continue exécutant des commandes Wake. Sa **présence** active Wake : lui attribuer `false` ou une chaîne vide ne le désactive pas. L’ancien nom `WAILS_USE_WAKE` est ignoré.

## Créer l’application

Exécutez l’assistant de création de projet :

```sh
wails3 init -ui
```

L’assistant de projet s’ouvre dans le navigateur. Il crée les fichiers sources ; l’assistant `wails3 setup` distinct présenté ci-dessous configure leur compilation.

| Écran | Choix pour ce tutoriel | Étape suivante |
| --- | --- | --- |
| **Projet — Créer un projet Wails** | Saisissez `wake-demo`, puis sélectionnez **Continuer**. | Le projet sera créé dans un nouveau dossier `wake-demo` sous le répertoire affiché. |
| **Framework — Choisir un framework** | Sélectionnez **Vanilla**, puis **Continuer**. | Vous obtenez un petit frontend sans étape d’apprentissage propre à un framework. |
| **Langage — Préférence de langage** | Sélectionnez **JavaScript**, puis **Continuer**. | L’écran **Liaisons**, réservé à TypeScript, est ignoré. |
| **Détails — Détails du projet** | Utilisez `Wake Demo` comme nom de produit, `com.example.wakedemo` comme identifiant de bundle et `0.1.0` comme version. Cliquez sur **Créer le projet**. | Lorsque Wake est activé, les détails sont écrits dans `wails.yaml`. |

N’utilisez pas **Utiliser les valeurs par défaut globales** pour cette leçon : cette option omet les choix du framework et du langage. Si vous choisissez TypeScript dans un autre projet, l’écran **Liaisons** vous permet de choisir entre des interfaces et des classes pour les liaisons de données Go générées.

Revenez au terminal et attendez la fin de la création du projet. Exécutez ensuite :

```sh
cd wake-demo
go get github.com/wailsapp/wails/v3@0243da963
```

Vous devriez maintenant disposer de `main.go`, d’un module Go, d’un répertoire `frontend` et d’un fichier `wails.yaml` natif. Le modèle intégré utilise directement Wake ; vous n’avez pas besoin de migrer ce nouveau projet. Un modèle communautaire peut contenir d’anciennes tâches personnalisées et nécessiter à la place l’examen de migration décrit ci-dessous.

## Configurer la compilation, écran par écran

Depuis le répertoire `wake-demo`, exécutez :

```sh
wails3 setup
```

Le navigateur affiche maintenant **CONFIGURATION DE WAKE**. S’il affiche le processus de configuration général, vérifiez que `WAILS_EXP_USE_WAKE` est défini dans le terminal qui a lancé la commande et que vous utilisez l’interface en ligne de commande de la préversion.

Les captures d’écran ci-dessous illustrent cette procédure sous Linux. La cible correspondant à votre système hôte et les outils requis seront différents sous Windows ou macOS.

### Projet — uniquement si les métadonnées sont absentes

Pour l’application que vous venez de créer, cet écran est **ignoré** : `init` a déjà écrit les métadonnées. C’est le comportement attendu.

Lorsque vous configurez un répertoire de sources Go existant dépourvu de manifeste natif ou d’ancienne configuration de compilation, l’écran **À propos de votre projet** demande le nom du produit, le nom du projet, l’identifiant de l’application, la version, l’entreprise et la description. Ces valeurs restent dans l’assistant jusqu’à ce que vous sélectionniez **Enregistrer**. La configuration ne crée aucun fichier source de l’application. Si elle détecte des Taskfiles ou d’anciennes métadonnées, commencez par suivre la procédure de migration.

### Application — Que créez-vous ?

Choisissez **Application de bureau**. Laissez **Empaqueter des binaires existants** décoché et choisissez **Backend standard**. Cliquez sur **Continuer**.

![Écran Application avec Application de bureau et Backend standard sélectionnés, et l’empaquetage de binaires précompilés désactivé.](/assets/screenshots/wake/application.png)

- **Application de bureau** compile une interface graphique dotée d’un frontend web.
- **Application en ligne de commande** et **Serveur** désignent des programmes sans fenêtre de bureau. La modification de cette sélection ne réécrit pas le code source de votre application ; utilisez un projet source correspondant au type choisi.
- **Pure Go (expérimental)** est une option distincte. Pour les compilations de bureau prises en charge, elle supprime la nécessité d’un compilateur C, à condition que vos dépendances puissent également être compilées sans CGo. Elle ne supprime ni les bibliothèques d’exécution de la cible, ni tous les outils d’empaquetage et de signature.
- **Empaqueter des binaires existants** prépare une commande utilisant `--skip-build`.
  Cette option nécessite des binaires déjà compilés ; laissez-la désactivée pour un nouveau projet.

Le backend est un **paramètre commun à tout le projet**. Si vous le modifiez ici, le comportement des autres profils qui héritent de ces paramètres peut changer.

### Cibles — Où l’application doit-elle s’exécuter ?

Conservez uniquement la cible correspondant à l’ordinateur que vous utilisez, puis cliquez sur **Continuer**.
L’assistant sélectionne initialement le système d’exploitation et l’architecture de la machine hôte.

![Écran Cibles où seule la cible linux/amd64 est sélectionnée pour la machine hôte utilisée dans ce tutoriel.](/assets/screenshots/wake/targets.png)

Les noms des cibles combinent un système d’exploitation et une architecture de processeur : par exemple, `windows/amd64`, `darwin/arm64` pour un Mac Apple Silicon, ou `linux/amd64`.
`darwin` désigne macOS. Sélectionner une autre cible demande la création d’un exécutable pour la machine correspondante ; cela ne rend pas cet exécutable utilisable sur votre ordinateur actuel.

Vous pourrez sélectionner des architectures supplémentaires, les builds universels macOS et les cibles mobiles expérimentales lors d’un passage ultérieur dans la configuration.

### Paquets — Choisir les livrables de publication

Laissez les cases des formats de paquet décochées et l’option **Signer les livrables de publication** désactivée.
Cliquez sur **Continuer**.

![Écran Paquets sans format de paquet Linux supplémentaire sélectionné et avec la signature désactivée.](/assets/screenshots/wake/packages.png)

Si vous ne sélectionnez aucun format supplémentaire, Wake produit le livrable par défaut de la cible : un exécutable de bureau, ou un paquet d’application macOS pour une interface graphique macOS. Cela suffit pour tester l’application localement. Vous pourrez sélectionner ultérieurement des formats afin d’ajouter, lorsqu’ils sont pris en charge, des livrables de distribution tels que DMG, MSIX, NSIS, DEB ou AppImage.

La signature est facultative pour ce premier build local. La réussite d’un build non signé ne signifie pas que l’application est prête à être distribuée publiquement.

### Outils — Outils nécessaires à votre build

Lisez les prérequis correspondant à vos sélections. Cette liste provient du plan de build de Wake ; elle ne vous demande pas d’installer tous les outils sur chaque plateforme.

![Écran Outils affichant les prérequis déterminés pour le build local depuis les sources sous Linux.](/assets/screenshots/wake/tools.png)

Pour cette leçon, vous aurez besoin de Go et des outils frontend, ainsi que des outils de développement natifs requis par le backend d’interface graphique standard de votre machine hôte. Docker n’est pas systématiquement requis. Une compilation croisée avec CGo peut nécessiter un compilateur et des bibliothèques pour la cible, ou une chaîne d’outils de conteneur prise en charge.

Corrigez les lignes marquées **Non prêt**, **Configuration requise** ou **Mise à jour requise** avant de tenter le build local. Suivez les instructions affichées lorsqu’elles sont disponibles, ou consultez le [guide d’installation propre à la plateforme](/getting-started/installation/), puis cliquez sur **Revérifier**. Vous devrez peut-être redémarrer la configuration si une installation a modifié le `PATH` du terminal. Lorsque les prérequis de la machine hôte sont satisfaits, cliquez sur **Continuer**.

L’assistant peut enregistrer une configuration valide même si cet ordinateur ne dispose pas de certains outils. Cela est utile pour configurer la CI, mais, dans ce tutoriel, vous devez obtenir un build qui puisse s’exécuter ici.

### Signature — Uniquement si elle est sélectionnée

Cet écran est **ignoré**, car vous avez laissé la signature désactivée. Lors d’une exécution ultérieure, l’activation de l’option **Signer les livrables de publication** affiche les contrôles de signature existants pour les plateformes de bureau sélectionnées. L’option **Configurer la signature ultérieurement** désactive la signature dans le profil de configuration enregistré. La signature mobile se configure dans le fichier YAML et non au moyen de ces contrôles.

Les paramètres de signature par défaut se trouvent dans vos paramètres utilisateur. Wake peut utiliser des fichiers de certificat configurés en dehors du projet sans copier de clés privées ni de mots de passe dans `wails.yaml`. Les paramètres de signature explicites du projet sont prioritaires. Avant de préparer une publication publique, suivez les [instructions de configuration de la signature](/getting-started/setup/) et le [guide d’empaquetage et de signature](/experimental/yaml-builds/#package-and-sign-for-distribution).

### Enregistrement — Enregistrer votre configuration de build

Vérifiez le récapitulatif : application de bureau, backend standard, build depuis les sources, cible de votre machine hôte et signature désactivée. Cliquez sur **Enregistrer wails.yaml**.

![Écran Enregistrement récapitulant l’application avec interface graphique, le backend standard, le build depuis les sources et la cible linux/amd64 non signée.](/assets/screenshots/wake/save.png)

Vous devriez voir **Configuration de build enregistrée**, le chemin d’accès à votre manifeste, ainsi que :

```sh
wails3 build setup
```

Cliquez sur **Terminer**, puis revenez au terminal. La configuration a enregistré les choix de cible et de paquet sous `profiles.setup`. Elle met également à jour le type d’application, l’activation du frontend et les paramètres du backend. Les commandes frontend, programmes compagnons, hooks et paramètres sans rapport déjà présents sont conservés.

Le nom du profil est important : `wails3 build setup` utilise ces sélections enregistrées.
La commande `wails3 build` seule ne sélectionne pas automatiquement le profil `setup`.

## Compiler et exécuter votre première application

Exécutez la commande affichée par l’assistant :

```sh
wails3 build setup
```

Le premier build installe les dépendances frontend, génère les bindings et compile l’application. Les builds ultérieurs peuvent réutiliser les éléments inchangés. Attendez le message **build réussi** et consultez les chemins de sortie affichés par la CLI. Le répertoire de sortie par défaut est `bin` ; une interface graphique macOS est assemblée dans un paquet `.app`.

Lancez le mode développement sur votre propre ordinateur :

```sh
wails3 dev
```

La fenêtre de l’application de démarrage devrait apparaître. Saisissez votre nom et essayez son message de bienvenue. Modifiez le texte indicatif du champ de saisie du nom dans `frontend/index.html`, puis enregistrez le fichier ; le frontend devrait s’actualiser. Modifiez le message de bienvenue dans `greetservice.go` pour déclencher un nouveau build Go. Lorsque vous avez terminé, appuyez sur **Ctrl+C** dans le terminal.

Le développement est un workflow local sur la machine hôte. Un profil de publication peut contenir plusieurs cibles de compilation croisée sans changer la machine sur laquelle la fenêtre de développement peut s’exécuter. Pour compiler et lancer localement une application de production, utilisez `wails3 run` ; consultez la [référence de la commande run](/reference/wails-yaml-run/).

Vous disposez maintenant d’un projet source, d’un profil de build vérifié et d’une boucle de développement locale fonctionnelle. Validez le code source et `wails.yaml` dans le système de gestion de versions.
Les fichiers de travail `.wails/` générés ne constituent pas la source de votre configuration.

## Étapes suivantes

| Objectif | Prise en charge par Wake | Premiers points à vérifier |
| --- | --- | --- |
| Compiler pour un autre système de bureau | Cibles Windows, macOS et Linux, y compris les deux architectures macOS et les livrables universels. | Pure Go et CGo ont des exigences différentes en matière de compilateur. Testez sur le système d’exploitation et l’architecture de destination. |
| Créer des installateurs ou des paquets | MSIX, NSIS, application macOS/DMG, DEB, RPM, Arch, Void et AppImage. | NSIS nécessite `makensis`. AppImage nécessite un environnement d’exécution cible ; l’empaquetage croisé d’une interface graphique nécessite également de fournir une arborescence d’environnement d’exécution Linux. |
| Signer une publication | Signature portable avec une clé exportée pour Windows, macOS et DEB/RPM/Arch ; les magasins d’identités natifs restent disponibles. | Chaque identité, format et machine hôte présente des limites. AppImage et Void ne disposent d’aucun workflow de signature intégré. |
| Notariser une application macOS | Soumission à l’aide d’une clé d’API d’équipe et agrafage du ticket, ou utilisation d’un profil du trousseau macOS. | Une signature Developer ID appropriée et un horodatage restent obligatoires. Le fait que des signatures de test générées par Go passent les vérifications natives ne prouve pas qu’Apple acceptera votre publication. |
| Compiler plusieurs programmes | Un programme principal avec interface graphique, en ligne de commande ou serveur, accompagné de programmes compagnons nommés. | Déclarez les paquets et les types des programmes compagnons dans le fichier YAML. L’assistant ne génère pas leur code source. |
| Empaqueter des fichiers compilés ailleurs | `wails3 build setup --skip-build`. | Fournissez les noms d’exécutables et les architectures correspondants dans `build.output`. Ce drapeau constitue un choix ponctuel pour la commande, et non un paramètre de configuration permanent. |
| Ajouter l’automatisation des versions | Profils nommés, hooks de cycle de vie délimités, mise en cache et traitement parallèle des paquets. | Utilisez du YAML littéral et des scripts déclarés appartenant au projet. Les comportements arbitraires des Taskfiles ne constituent pas un mécanisme d’extension de Wake. |

Rouvrez `wails3 setup` pour modifier le profil `setup`. Modifiez les autres profils nommés et les paramètres avancés dans `wails.yaml`. Consultez le [guide de compilation YAML](/experimental/yaml-builds/) pour obtenir des exemples complets et connaître les exigences relatives aux paquets. Les applications avec interface graphique ont toujours besoin des bibliothèques système de la destination ; Pure Go n’en fait pas des applications entièrement statiques.

iOS et Android restent expérimentaux, nécessitent leurs SDK de plateforme ainsi que CGo, et sortent du cadre de cette première leçon consacrée aux applications de bureau. L’assistant ne rend pas les compilations mobiles portables entre des hôtes arbitraires.

## Migrer un projet existant

Travaillez dans une branche distincte ou sur une copie d’un projet **Wails v3** fonctionnel. Il s’agit d’une migration de la configuration de compilation, et non d’une mise à niveau depuis Wails v2. Conservez le résultat de compilation et les tests existants afin de pouvoir comparer le comportement après l’activation.

@note{type="note" title="Commencer par une migration révisée"}
Les nouveaux projets générés à partir des modèles hérités intégrés actuels peuvent être migrés avant la compilation de leur frontend. Wake reconnaît leurs sélecteurs de tâches standard et leur cycle de développement sans exécuter de fichiers Taskfile ni évaluer d’expressions personnalisées.

Les tâches personnalisées, les substitutions de variables et les stratégies de développement doivent toujours être révisées. Si l’analyse ne peut pas représenter votre projet, conservez le système de compilation hérité et transmettez les diagnostics. Ne supprimez pas un comportement dont vous avez besoin uniquement pour permettre l’activation. Le tutoriel ci-dessus consacré aux nouveaux projets ne nécessite aucune migration.
@end

### 1. Analyser sans modifier la compilation active

Activez Wake et utilisez la CLI de préversion indiquée au début de cette page. Exécutez la commande suivante depuis la racine du projet existant :

```sh
wails3 migrate --dry-run
```

Examinez les paramètres pris en charge et les blocages signalés. L’analyseur lit les déclarations historiques ; il n’exécute pas leurs commandes. Une simulation ne crée pas de manifeste actif. Pour obtenir des diagnostics lisibles par une machine, utilisez `wails3 migrate --json` ; utilisez `--dry-run --json` si vous souhaitez également éviter l’écriture d’un brouillon.

### 2. Créer et examiner la proposition inactive

```sh
wails3 migrate
```

Cette opération crée **`wails.migrated.yaml`**, et non `wails.yaml`. Examinez ce fichier avec vos Taskfiles existants et `build/config.yml`. Vérifiez les métadonnées du projet, les commandes du frontend, les architectures cibles, le contenu des paquets, les ressources et les paramètres de signature. Conservez les secrets dans des variables d’environnement ou dans votre gestionnaire d’identifiants.

Les tâches intégrées connues peuvent être associées à des paramètres natifs. Les tâches de cycle de vie prudentes et sans argument qui appellent un script appartenant au projet peuvent devenir des hooks. Tout autre comportement de tâche personnalisée accessible constitue un blocage : Wake ne l’abandonne pas silencieusement et ne l’exécute pas non plus via Task comme solution de repli.

Si un brouillon existe déjà, la migration n’écrase pas les modifications que vous avez examinées. Conservez-en une copie avant de le remplacer ou de le régénérer volontairement. Ne renommez pas un brouillon bloqué en `wails.yaml` pour contourner les contrôles d’activation. Réduisez ou remplacez les comportements historiques non pris en charge, puis relancez l’analyse avant l’activation. Si le projet ne peut pas encore être représenté, conservez son système de compilation historique et signalez cette lacune.

### 3. Activer explicitement

Après avoir examiné le brouillon et résolu les blocages, exécutez :

```sh
wails3 migrate --activate --dry-run
wails3 migrate --activate
```

La première commande vérifie l’activation sans la publier. La seconde contrôle de nouveau les déclarations historiques et valide le brouillon avant de créer `wails.yaml` de façon atomique ; une activation réussie supprime le fichier inactif `wails.migrated.yaml`. Elle ne modifie ni les Taskfiles, ni `build/config.yml`, ni les scripts, ni les ressources. À partir de ce moment, le manifeste actif régit les compilations Wake. Les anciens fichiers ne constituent pas une deuxième source de configuration.

Certains projets plus anciens possèdent déjà un fichier `wails.yaml` contenant uniquement des métadonnées et une section `info:` de premier niveau. Bien qu’il ne s’agisse pas d’un manifeste Wake natif, ce fichier occupe tout de même la destination d’activation. Conservez ce fichier historique et déplacez-le délibérément avant l’activation ; l’outil de migration ne l’écrasera pas et ne le renommera pas à votre place.

### 4. Vérifier l’application migrée

```sh
wails3 config check
wails3 build
wails3 dev
```

Comparez l’application et le contenu des paquets avec la compilation précédente, y compris vos comportements de publication personnalisés. Exécutez les tests de votre projet. Vérifiez explicitement chaque profil migré, par exemple `wails3 config check release` et `wails3 build release` lorsque le brouillon contient un profil `release`.

Vous pouvez maintenant exécuter `wails3 setup` pour ajouter ou remplacer le profil dédié `setup`. Cette commande ne sélectionne ni ne modifie automatiquement un profil `release` existant.

### Si vous devez interrompre l’expérimentation

La suppression de `WAILS_EXP_USE_WAKE` désactive l’activation facultative, mais un fichier `wails.yaml` actif bloque toujours le routage automatique vers les Taskfiles. Un YAML non valide provoque également un arrêt avec une erreur. Il n’existe aucune commande de migration inverse.

Pour revenir à l’ancien workflow, restaurez dans le système de gestion de versions la révision complète du projet antérieure à la migration, y compris l’absence du manifeste actif, puis utilisez la CLI et les dépendances compatibles avec cette révision. Avant cela, conservez toutes les nouvelles modifications du code source dont vous avez besoin. L’utilisation d’une branche d’essai distincte facilite considérablement cette restauration volontaire.

## Diagnostiquer un problème et donner votre avis

| Symptôme | Première vérification |
| --- | --- |
| La configuration affiche l’assistant général | Définissez `WAILS_EXP_USE_WAKE` dans le terminal de lancement et vérifiez quel exécutable `wails3` se trouve dans `PATH`. |
| La configuration demande une migration | Vous avez ouvert un projet historique ou un modèle communautaire doté d’une configuration historique. Examinez la migration avant l’activation. |
| Une compilation enregistrée signale des outils manquants | L’enregistrement consigne l’intention. Revenez à **Outils**, corrigez les exigences de l’hôte, puis relancez la vérification. |
| Un exécutable compilé pour une autre cible ne démarre pas sur cet hôte | Exécutez-le sur le système d’exploitation et l’architecture sélectionnés ; utilisez la cible correspondant à l’hôte pour le développement local. |
| Le mode précompilé ne trouve pas de binaire | Vérifiez les [noms et chemins des fichiers précompilés](/experimental/yaml-builds/#use-prebuilt-desktop-executables), y compris les sous-répertoires destinés aux différentes cibles. |
| La modification d’un manifeste échoue | Exécutez `wails3 config check` pour le profil concerné et suivez le diagnostic du champ. Wake ne se rabat pas sur les Taskfiles. |

Inspectez un build avant de l’exécuter avec `wails3 build setup --plan --json`.
`wails3 config check` valide la configuration, tandis que `wails3 eject` écrit un instantané `wails.ejected.yaml` inactif sans modifier le manifeste actif.
Examinez la sortie de la commande avant de la partager : les chemins du projet et les emplacements des identifiants peuvent permettre d’identifier votre environnement local.

Faites part de vos commentaires dans les [discussions sur les expérimentations](https://github.com/wailsapp/wails/discussions/categories/experiments) ou [signalez un bug reproductible](https://github.com/wailsapp/wails/issues). Indiquez les révisions de votre CLI et du framework, le système d’exploitation et l’architecture de l’hôte, la cible sélectionnée, la commande exacte, si vous avez sélectionné Pure Go, ainsi que le diagnostic pertinent. Précisez si l’échec s’est produit lors de la compilation, de la création du paquet, de la signature ou de l’exécution de l’application. Si cela s’avère utile, joignez un manifeste minimal dont les secrets ont été supprimés ainsi que la sortie de `wails3 doctor`.
