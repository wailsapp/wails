---
title: "Fonctionnalités expérimentales"
description: "Un espace consacré aux expérimentations en cours dans Wails v3 : leur nature, leur raison d’être et la manière de nous faire part de vos retours."
slug: "experimental"
sourcePath: "experimental/index.md"
---

@note{type="caution" title="Vous entrez dans la zone des expérimentations"}
Tout ce qui figure dans cette section est, par définition, expérimental. Ces fonctionnalités sont facultatives, désactivées par défaut et peuvent changer de forme, être renommées ou être entièrement supprimées d’une version à l’autre. Ne fondez aucun élément essentiel sur elles sans vous préparer à des changements fréquents.

@end

## Ce que signifie « expérimental »

Wails développe ses expérimentations au grand jour. Une expérimentation est une idée qui nous semble suffisamment prometteuse pour vous être proposée tôt, mais sur laquelle nous ne nous sommes pas encore pleinement engagés. Nous la publions *parce que* nous voulons tirer des enseignements de son utilisation réelle avant de décider si elle doit devenir une composante permanente et prise en charge de Wails.

Cela implique quelques caractéristiques communes à tout ce qui figure dans cette section :

- **Elle est facultative.** Les expérimentations ne modifient jamais le comportement par défaut de `wails3`. Vous les activez délibérément, généralement à l’aide d’une variable d’environnement ou d’une option de compilation. Lorsqu’elles sont désactivées, votre flux de travail existant ne change en rien.
- **Elle peut ne pas perdurer.** Certaines expérimentations deviennent des fonctionnalités stables. D’autres sont remaniées au point de devenir méconnaissables, ou abandonnées. Nous préférons essayer publiquement des idées et apprendre rapidement plutôt que de ne publier que ce dont nous sommes déjà certains.
- **L’API n’est pas figée.** Les noms, les options, les valeurs par défaut et le comportement peuvent évoluer d’une version à l’autre pendant que l’expérimentation prend forme. Les notes de version signaleront ces changements, mais n’attendez pas les garanties de stabilité offertes par les fonctionnalités stables.

## Nous souhaitons recueillir vos retours

C’est le point essentiel. La pérennité des expérimentations dépend des retours des personnes qui les utilisent réellement. Si vous en essayez une, nous souhaitons sincèrement savoir :

- A-t-elle fonctionné pour votre projet ? Où a-t-elle échoué ?
- Était-elle plus rapide, plus claire ou plus agréable, ou le changement n’en valait-il pas la peine ?
- Quelles conditions devraient être remplies pour que vous l’utilisiez par défaut ?

Les retours les plus utiles sont concrets : ce que vous avez exécuté, le résultat attendu et ce qui s’est réellement produit. Chaque expérimentation possède son propre fil dans la catégorie **Expérimentations** des discussions GitHub :

@container{display="grid" columns="2" gap="1rem"}
@linkcard{title="Discussions sur les expérimentations" href="https://github.com/wailsapp/wails/discussions/categories/experiments" description="Trouvez le fil consacré à l’expérimentation que vous utilisez et expliquez-nous comment l’essai s’est déroulé, ce qui n’a pas fonctionné ou ce qui manque."}
@end

## Expérimentations actuelles

@container{display="grid" columns="2" gap="1rem"}
@linkcard{title="Wake" href="/experimental/wake/" description="Essayez le système de compilation YAML natif : créez et configurez un projet dans un seul assistant, compilez-le et empaquetez-le, ou examinez une migration à sens unique depuis des Taskfiles."}
@linkcard{title="Compilations sans CGO" href="/experimental/#cgo-free-builds-purego" description="Compilez des applications macOS et Linux avec CGO_ENABLED=0 : aucune chaîne d’outils C et compilation croisée sans Docker. Activez cette fonctionnalité avec WAILS_PUREGO=1 ou l’étiquette de compilation purego."}
@linkcard{title="Contrôle par LLM (MCP)" href="/guides/mcp-service/" description="Un serveur Model Context Protocol intégré qui permet aux agents LLM d’inspecter, de tester et de piloter une application Wails en cours d’exécution, sans aucun code utilisateur ; il s’active à l’aide d’une étiquette de compilation."}
@end

## Compilations sans CGO (purego)

@note{type="caution" title="Fonctionnalité expérimentale"}
Le backend sans CGO est facultatif et s’active au moyen de l’étiquette de compilation `purego` (ou de `WAILS_PUREGO=1`). Il ne s’agit **pas** du backend par défaut. Sans cette activation, les compilations macOS et Linux utilisent le backend cgo exactement comme auparavant. Le comportement peut changer d’une version à l’autre.

@end

### Présentation

Sous macOS et Linux, Wails communique normalement avec les bibliothèques d’interface utilisateur du système par l’intermédiaire de cgo : AppKit et WKWebView sous macOS, GTK et WebKitGTK sous Linux. Le backend sans CGO effectue les mêmes appels depuis du Go pur au moyen de [purego](https://github.com/ebitengine/purego), qui charge les bibliothèques système lors de l’exécution.

L’API publique de l’application reste inchangée. Sélectionnez le backend lors de la compilation ; le code de l’application n’a pas besoin d’être modifié. Testez votre application sur le système de destination avant de la distribuer.

Avantages :

- **Aucune chaîne d’outils C pour la compilation Go.** La compilation Go ne nécessite ni les outils en ligne de commande Xcode, ni `gcc`, ni les en-têtes de développement Linux. Les compilations du frontend nécessitent toujours leurs outils JavaScript habituels.
- **Compilation croisée sans Docker.** Compilez un exécutable macOS ou Linux depuis n’importe quel système d’exploitation avec le simple `go build`.

Windows se compile déjà sans cgo et n’est pas concerné.

### Avant l’activation

Utilisez des versions de la CLI et du framework qui incluent toutes deux cette expérimentation. La seule mise à jour de la CLI n’ajoute pas le backend à une ancienne dépendance du framework dans le fichier `go.mod` de votre application. Pour essayer la version de développement, exécutez ces commandes depuis le répertoire de votre application :

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@master
go get github.com/wailsapp/wails/v3@master
```

Les deux commandes suivent la branche mouvante `master`. Pour des compilations reproductibles, utilisez le même commit testé pour la CLI et le framework. Consultez le [guide d’installation](/getting-started/installation/) pour configurer la chaîne d’outils Go et la CLI.

Les backends macOS et Linux expérimentaux prennent en charge `amd64` et `arm64`. La compilation croisée crée un exécutable pour la cible demandée ; elle ne l’exécute pas sur la machine de compilation. Testez-le sur le système d’exploitation et l’architecture cibles.

### Utilisation

Choisissez l’une des méthodes suivantes pour sélectionner le backend.

#### Avec la CLI Wails

Dans un projet généré avec cette version de la CLI, `build`, `package` et `dev` utilisent le backend expérimental lorsque `WAILS_PUREGO` est activé. Dans un shell POSIX :

```sh
WAILS_PUREGO=1 wails3 dev
WAILS_PUREGO=1 wails3 build
WAILS_PUREGO=1 wails3 package
```

Dans PowerShell :

```powershell
$env:WAILS_PUREGO = "1"
wails3 build
```

La CLI reconnaît `1`, `true`, `on` et `yes` sans tenir compte de la casse ni des espaces environnants. Une valeur vide n’active pas cette expérimentation.

Vous pouvez également transmettre l’étiquette à une seule commande de compilation ou d’empaquetage :

```sh
wails3 build -tags purego
wails3 package -tags purego
```

Lorsque l’étiquette est définie, les Taskfiles standard mis à jour compilent avec `CGO_ENABLED=0` en utilisant la chaîne d’outils Go de l’hôte. Par exemple, compilez un exécutable Linux ARM64 depuis macOS ou Windows :

```sh
wails3 build -tags purego GOOS=linux ARCH=arm64
```

Les anciens projets conservent leurs Taskfiles existants. Comparez vos tâches de compilation personnalisées aux modèles mis à jour, ou utilisez directement Go pour tester le backend sans modifier ces tâches. La variable d’environnement sélectionne l’étiquette dans les commandes de la CLI Wails ; avec le simple `go build`, vous devez spécifier explicitement l’étiquette.

#### Directement avec Go

Compilez d’abord votre frontend afin que toutes les ressources référencées par `go:embed` existent, puis compilez l’application :

```sh
CGO_ENABLED=0 go build -tags purego,production .
```

`CGO_ENABLED=1 go build -tags purego` fonctionne également : le backend reste le même, mais cgo demeure disponible pour le reste de votre programme.

Pour une compilation croisée, définissez également `GOOS` et `GOARCH` :

```sh
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -tags purego,production -o myapp .
```

Sous Linux, sélectionnez la pile GTK3 avec `-tags purego,gtk3,production`.

Pour revenir au backend par défaut, supprimez l’étiquette `purego` et désactivez `WAILS_PUREGO` (avec `unset WAILS_PUREGO` dans un shell POSIX ou `Remove-Item Env:WAILS_PUREGO` dans PowerShell). Lorsque vous revenez au backend cgo sous macOS ou Linux, supprimez également toute définition explicite de `CGO_ENABLED=0`.

Si vous compilez avec `CGO_ENABLED=0` mais oubliez l’étiquette, la compilation s’arrête avec une seule erreur qui indique la correction à appliquer :

```
undefined: requires_CGO_ENABLED_1_or_build_tag_purego
```

### Prérequis d’exécution

L’application charge les bibliothèques système au démarrage. Elles doivent donc être installées sur la machine cible, comme pour une compilation avec cgo :

| Plateforme | Bibliothèques | Version minimale prise en charge |
|---|---|---|
| macOS | AppKit, WebKit (inclus dans le système d’exploitation) | Identique au backend cgo |
| Linux (par défaut) | GTK 4, WebKitGTK 6.0 | GTK 4.8, WebKitGTK 2.40 (Debian 12) |
| Linux (`gtk3`) | GTK 3, WebKit2GTK 4.1 | GTK 3.24, WebKit2GTK 2.36 (Ubuntu 22.04) |

Les fonctions plus récentes que ces versions minimales sont recherchées au démarrage de l’application et utilisées uniquement si la bibliothèque installée les fournit. Un même binaire fonctionne donc aussi bien sur les anciens systèmes que sur les nouveaux.

L’application utilise toujours les bibliothèques d’interface utilisateur natives et peut effectuer une liaison dynamique avec des bibliothèques système telles que la libc sous Linux. `CGO_ENABLED=0` ne désigne pas un exécutable entièrement statique. Si une bibliothèque d’exécution requise est absente, installez la pile GTK/WebKit correspondante sur le système cible ; installer les en-têtes C sur la machine de compilation ne résout pas l’absence d’une bibliothèque d’exécution.

### Création des paquets et signature

Le backend supprime l’exigence d’une chaîne d’outils C pour la compilation. Chaque format de paquet conserve ses propres exigences. En particulier, sélectionner `purego` ne rend pas multiplateformes les outils de déploiement AppImage existants. Consultez la page [Création de paquets Linux](/guides/build/linux/).

Les applications macOS compilées pour une autre plateforme nécessitent toujours la signature et la notarisation appropriées avant leur distribution. Les programmes d’installation Windows nécessitent toujours leurs propres outils de création. Consultez les pages [Création de paquets macOS](/guides/build/macos/) et [Création de paquets Windows](/guides/build/windows/). Pour les dépendances qui utilisent cgo, `CGO_ENABLED=1` nécessite également un compilateur C ainsi que tous les en-têtes ou bibliothèques de la cible ; sélectionner le backend Wails en Go pur ne supprime pas les exigences de ces dépendances.

### Validation

Les vérifications des liaisons générées et `puregocheck` inspectent les symboles natifs, les protections conditionnelles liées aux versions facultatives et la couverture des backends. Le workflow de parité Linux teste GTK4 et GTK3 avec les deux backends, notamment les requêtes de ressources, les corps de requête binaires et la répartition native ; il effectue également une compilation croisée pour les cibles Linux ARM64. Les vérifications de compilation croisée établissent que le code peut être compilé, tandis que les vérifications d’exécution testent les bibliothèques natives installées sur la machine de test.

Le [guide de test des backends](https://github.com/wailsapp/wails/blob/master/v3/TESTING.md) décrit les tests de parité et les tests d’entrée de bout en bout. Lorsque vous signalez un problème, distinguez un échec de compilation d’un échec à l’exécution ou lors de la création du paquet.

### Différences connues

- **Sous Linux, les raccourcis clavier des menus sont enregistrés au niveau de l’application.** Évitez d’attribuer des raccourcis identiques à des menus différents. La suppression explicite du menu d’une fenêtre avec `SetMenu(nil)` empêche également l’activation, dans cette fenêtre, des raccourcis du menu supprimé ; les autres fenêtres conservent leurs menus et leurs raccourcis.
- **Les boîtes de dialogue d’ouverture et d’enregistrement de macOS** s’exécutent dans un processus système distinct, comme avec cgo ; les tests vérifient qu’elles s’ouvrent et peuvent être annulées correctement.

Vous avez trouvé autre chose ? Faites-le-nous savoir.

### Retours

Racontez-nous votre expérience dans la discussion **Compilations sans CGO**, sous [Expérimentations](https://github.com/wailsapp/wails/discussions/categories/experiments) : indiquez votre système d’exploitation et sa version, votre méthode de compilation, ainsi que ce qui a fonctionné ou non. Le signalement le plus utile que vous puissiez nous fournir est un cas d’échec accompagné de la sortie de `wails3 doctor`.
