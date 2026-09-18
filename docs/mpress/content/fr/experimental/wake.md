---
title: "Wake"
description: "Un exécuteur de builds expérimental adapté à Wails, qui exécute vos Taskfiles existants avec des builds incrémentaux plus rapides, une sortie structurée et une exécution parallèle par défaut."
slug: "experimental/wake"
sourcePath: "experimental/wake.md"
---

@note{type="caution" title="Fonctionnalité expérimentale"}
Wake s’active explicitement au moyen de `WAILS_USE_WAKE=true` et n’est **pas** l’exécuteur par défaut. Lorsque la variable n’est pas définie, `wails3 build / package / sign / task` se comporte exactement comme auparavant. La couverture fonctionnelle et le comportement peuvent changer d’une version à l’autre.

@end

Wake est un **exécuteur de builds alternatif expérimental** pour `wails3`. Il lit les mêmes `Taskfile.yml` que ceux déjà présents dans votre projet — avec la même syntaxe pour les tâches, les dépendances, les variables, les modèles, les inclusions et les espaces de noms de plateformes — puis les exécute à l’aide d’un exécuteur adapté à Wails plutôt qu’avec l’environnement d’exécution généraliste de [Task](https://taskfile.dev).

L’objectif n’est pas de remplacer Task, mais de fournir un exécuteur spécialement conçu pour la manière dont les projets Wails sont réellement construits, avec une sémantique, une sortie et des valeurs par défaut cohérentes avec le reste de la CLI `wails3`. **Si vous utilisez exclusivement Wake, vos Taskfiles ne changent pas.**

## Pourquoi Wake existe

Wake et l’environnement d’exécution de Task sont tous deux intégrés à `wails3` — aucun ne nécessite l’installation d’un binaire distinct. La différence est que Wake **comprend le domaine**. Un exécuteur généraliste exécute les étapes répertoriées dans un Taskfile, dans l’ordre indiqué. Wake sait ce qu’est *réellement* un build Wails : le bundle frontend est incorporé au binaire, celui-ci est empaqueté dans des artefacts propres à chaque plateforme, et les icônes et les bindings sont également générés au cours du build. Wake exploite ces connaissances pour optimiser le build d’une façon inaccessible à un exécuteur générique.

- **Wake n’effectue que le travail réellement nécessaire au build.** Il suit lui-même les véritables entrées et sorties de chaque étape. Pour un build Go, il s’agit du graphe des modules et des sorties des étapes dont il dépend. Ainsi, lorsque rien de pertinent n’a changé, Wake ignore entièrement le compilateur et l’éditeur de liens au lieu de les réexécuter. Un exécuteur généraliste ne peut ignorer une étape que si le Taskfile a indiqué à l’avance et précisément les fichiers à surveiller ; Wake le détermine grâce à ce qu’il sait déjà du build. Pour un nouveau build sans aucune modification, cela représente environ **~20 ms (Wake), contre ~316 ms (Task)**. Les builds à froid prennent le même temps réel : leur durée dépend principalement de `npm install`, de Vite et du compilateur Go.

- **Wake sait quelles opérations peuvent s’exécuter simultanément.** Comme il comprend quelles étapes sont indépendantes, il les exécute en parallèle par défaut, et la ligne de résultat indique le gain de temps obtenu. Désactivez ce comportement avec `WAKE_SERIAL=true` lorsque l’entrelacement des sorties d’étapes parallèles risquerait de compliquer une investigation.

- **Une sortie structurée contrôlée par wails3.** Wake utilise le propre système de rapports de wails3 : une ligne par étape planifiée, un état actualisé en direct, une répartition finale des phases par couleur et des liens `file:line` cliquables dans les panneaux d’échec. `NO_COLOR` et les environnements sans TTY, tels que les journaux de CI, sont pris en charge proprement avec un affichage simplifié.

- **Intégré pour évoluer avec Wails.** Comme Wake fait partie de `wails3` au lieu d’être un outil tiers, de nouvelles fonctionnalités de build peuvent y être ajoutées directement, sans attendre qu’un projet distinct les implémente. Cela ouvre également la voie à l’exécution native de scripts et d’outils multiplateformes là où un Taskfile lance aujourd’hui le binaire `wails3` dans un shell, en créant un processus à chaque appel. L’intégration de ce travail au processus réduit la surcharge et permettra d’autres accélérations.

## Activation de Wake

Wake est entièrement conditionné par la variable d’environnement `WAILS_USE_WAKE=true`. Lorsqu’elle n’est pas définie, ou qu’elle possède une valeur autre que `true`, chaque commande `wails3` utilise l’environnement d’exécution Task intégré exactement comme auparavant.

```bash
# Default: Task runtime, no Wake involvement
wails3 build

# Opt in: Wake drives build / package / sign / task <name>
WAILS_USE_WAKE=true wails3 build
WAILS_USE_WAKE=true wails3 package
WAILS_USE_WAKE=true wails3 task <some-task-name>
```

L’indicateur s’applique à `wails3 build`, `wails3 package`, `wails3 sign` et `wails3 task <name>`. `wails3 dev` n’est **pas** encore concerné : le processus de surveillance en mode développement utilise toujours son propre pipeline.

@note{type="tip" title="L’activation de Wake est sans risque"}
Si Wake rencontre une fonctionnalité de Taskfile qu’il n’implémente pas, il confie l’intégralité de l’exécution à l’environnement Task intégré, au sein du même processus : il n’est pas nécessaire d’installer un binaire `task` externe. Dans le pire des cas, vous obtenez exactement le comportement que vous auriez eu sans l’indicateur.

@end

## Surcharges locales en couches

Wake prend en charge un **Taskfile de base complété par des surcharges locales**. Placez un fichier à côté de votre `Taskfile.yml` ; ses définitions seront prioritaires :

| Fichier | Rôle | Priorité |
| --- | --- | --- |
| `Taskfile.yml` | base, versionnée | la plus basse |
| `Taskfile.override.yml` / `.yaml` | surcharges versionnées communes à l’équipe | intermédiaire |
| `Taskfile.local.yml` / `.yaml` | personnel, généralement ignoré par Git | la plus élevée |

**Sémantique de fusion (la définition locale est prioritaire) :**

- Une tâche portant le **même nom** remplace la tâche de base. Lorsqu’ils sont fournis par la surcharge, les champs de type liste (`cmds`, `deps`, `sources`, `generates`, `platforms`, `status`, `preconditions`, `aliases`) **remplacent** ceux de la base ; les champs omis par la surcharge sont conservés depuis la base.
- `env` et `vars` sont **fusionnés clé par clé**, la surcharge étant prioritaire en cas de conflit.
- Une tâche présente **uniquement** dans un fichier de surcharge est **ajoutée**.

Par exemple, si le fichier `Taskfile.yml` versionné effectue les builds avec des indicateurs de développement, mais que votre machine devrait toujours produire des builds de production :

```yaml
# Taskfile.local.yml (git-ignored, yours)
tasks:
  build:
    cmds:
      - go build -tags production -o bin/app .
  smoke:
    cmds:
      - ./bin/app --selftest
```

Désormais, `build` exécute votre commande de production et `smoke` est disponible, sans modifier le Taskfile versionné.

@note{type="note" title="Modèle de confiance"}
Les fichiers de surcharge sont détectés et appliqués automatiquement, sans confirmation. Cela n’accorde aucune nouvelle capacité : un Taskfile peut déjà exécuter n’importe quelle commande shell ; une surcharge ne peut donc rien faire qu’une modification de `Taskfile.yml` ne permettrait pas déjà. Le fichier `Taskfile.override.*` versionné apparaît dans les différences des PR ; `Taskfile.local.*` est créé sur votre propre machine. Une surcharge incorrecte interrompt l’exécution au lieu d’être ignorée silencieusement. Définissez `WAILS_NO_OVERRIDES=true` pour désactiver entièrement la détection des surcharges et obtenir des builds de CI déterministes.

@end

## Repli automatique

Si Wake rencontre une fonctionnalité de Taskfile qu’il n’implémente pas, il confie l’intégralité de l’exécution à l’environnement Task intégré. Les éléments suivants déclenchent actuellement ce repli :

- `dotenv` au niveau du Taskfile
- les modes `output` autres que `interleaved`
- un bloc `requires`
- `interval` (au niveau du Taskfile ou de la tâche)
- les modes `run` autres que `always`
- `short` dans une tâche
- `defer` dans une tâche

## Variables d’environnement

| Variable | Effet |
| --- | --- |
| `WAILS_USE_WAKE` | `true` active Wake pour les verbes `wails3` pouvant être acheminés ; tout le reste utilise le moteur d’exécution de Task |
| `WAILS_NO_OVERRIDES` | `true` ignore la détection de `Taskfile.local.*` / `.override.*` (compilations déterministes) |
| `WAKE_VERBOSE` | Diffuse en direct les sorties stdout/stderr des sous-processus au lieu de les capturer pour ne les afficher qu’en cas d’échec |
| `WAKE_SILENT` | Supprime entièrement la sortie des tâches |
| `WAKE_SERIAL` | `true` désactive la distribution parallèle de `deps:` (le parallélisme est activé par défaut) |
| `WAKE_FORCE` | `true` contourne tous les caches afin d’effectuer une véritable recompilation complète |
| `WAKE_DEBUG` | Journalise les détails internes du résolveur (DAG, dépendances, références de variables, acheminement de l’exécution) |
| `WAKE_NOTICE` | `off` pour masquer l’avis « wake (expérimental) » à chaque exécution |

Le cache de compilation se trouve dans `.wake/cache.json` (Task utilise `.task/`).

## Commentaires

Wake est une fonctionnalité expérimentale, et vos commentaires détermineront son évolution. Si vous l’essayez, nous aimerions savoir s’il a été plus rapide et plus clair, et si quelque chose a cessé de fonctionner : les rapports les plus utiles indiquent ce que vous avez exécuté, ce que vous attendiez et ce qui s’est réellement produit. Faites-nous-en part dans la [discussion sur Wake](https://github.com/wailsapp/wails/discussions/5679).
