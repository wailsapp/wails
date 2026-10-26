---
title: "Politique de publication"
description: "Comment Wails v3 versionne ses publications nocturnes et hebdomadaires, comment les responsables les préparent et pourquoi les utilisateurs doivent épingler les versions."
slug: "contributing/release-policy"
sourcePath: "contributing/release-policy.md"
---

## Le modèle en un tableau

Les publications de Wails v3 suivent un **calendrier, et non les fonctionnalités**. Après la version de référence `v3.0.0` :

| Publication | Tag | Fréquence | Publication GitHub | npm `@wailsio/runtime` |
| --- | --- | --- | --- | --- |
| Nocturne | `v3.<next>.0-nightly.<YYYYMMDD>` (date UTC ; `.N` ajouté en cas de nouvelle exécution le même jour) | automatique, quotidienne, uniquement lorsque `UNRELEASED_CHANGELOG.md` contient des éléments | préversion, **pas** Latest | même version sans `v`, dist-tag `nightly` |
| Hebdomadaire | `v3.<next>.0` | automatique, le dimanche 20:00, heure d’Australia/Sydney **(jour non confirmé)** | publication normale, Latest | même version, dist-tag `latest` |
| Correctif urgent | `v3.<current>.<N>`, `N > 0` | rare, manuelle | version normale, Latest | même version, dist-tag `latest` |
| Majeure | `v4.0.0` | uniquement pour les changements incompatibles | publication normale | `latest` |

`<next>` correspond à la version mineure qui suit la plus élevée des versions **publiées** v3. Une version nightly est un aperçu de la **prochaine version hebdomadaire** ; sa base est donc le tag de version le plus élevé (les préversions sont ignorées) : après un correctif `v3.1.1`, la version nightly suivante reste `v3.2.0-nightly.*`. Elle doit porter un numéro de version supérieur, car une préversion est classée avant la version finale correspondante : `v3.0.0-nightly.<date>` serait classée avant `v3.0.0`.

**L’incrémentation du numéro mineur ne constitue pas un jalon fonctionnel.** `v3.4.0` désigne « la publication hebdomadaire qui suit `v3.3.0` ». Consultez les notes de publication, et non le numéro de version, pour connaître les changements.

La politique est déterminée par la **forme** de la version : un suffixe `-nightly.` indique une version nocturne ; tout autre suffixe (`-alpha`, `-beta`, `-rc`) indique une ancienne build antérieure à la disponibilité générale ; l’absence de suffixe indique une publication.

### Versions synchronisées

La CLI, le module Go et le paquet npm `@wailsio/runtime` portent la **même version**, tirée de `v3/internal/version/version.txt` (`3.2.0-nightly.20261102` sur npm pour le tag `v3.2.0-nightly.20261102`). Les projets créés avec `wails3 init` épinglent `@wailsio/runtime` à la propre version de la CLI ; les builds de développement de la CLI utilisent `latest` comme solution de repli.

## Exemple détaillé

| Date | Événement | Tag | `go install ...@latest` est résolu en |
| --- | --- | --- | --- |
| Dim. 25 oct. | disponibilité générale, publication manuelle | `v3.0.0` | `v3.0.0` |
| Lun. 26 – sam. 31 oct. | versions nocturnes | `v3.1.0-nightly.20261026` ... `v3.1.0-nightly.20261031` | `v3.0.0` |
| Dim. 1 nov. 20:00, heure de Sydney | publication hebdomadaire | `v3.1.0` | `v3.1.0` |
| Lun. 2 nov. | version nocturne (aperçu de la prochaine publication hebdomadaire) | `v3.2.0-nightly.20261102` | `v3.1.0` |
| Mer. 4 nov. | correctif urgent, manuel, depuis une branche de publication | `v3.1.1` | `v3.1.1` |
| Dim. 8 nov. | publication hebdomadaire | `v3.2.0` | `v3.2.0` |

Une version hebdomadaire est toujours classée au-dessus de toutes les versions nocturnes portant son propre numéro, et une version nocturne ne devient jamais `@latest`.

## Commandes à exécuter par les utilisateurs

| Objectif | Commande |
| --- | --- |
| Publication la plus récente (hebdomadaire ou correctif urgent) | `go install github.com/wailsapp/wails/v3/cmd/wails3@latest` ou `wails3 update cli` |
| **Reproductible (recommandé)** | épinglez un tag exact : `v3.1.0` dans `go.mod`, l’intégration continue et la documentation |
| Une version nocturne | son tag exact (`@v3.2.0-nightly.20261102`), ou `wails3 update cli -nightly` |
| Build la plus récente, quel que soit son type | `wails3 update cli -pre` |
| Paquet npm du runtime | `@wailsio/runtime` (`latest` = hebdomadaire/correctif urgent, `nightly` = nocturne) |

Pourquoi cela fonctionne : Go résout `@latest` vers le tag le plus élevé **sans suffixe de préversion** et ne tient compte des préversions que lorsqu’aucune publication n’existe. La case « pre-release » de GitHub est invisible pour le proxy Go ; les versions nocturnes portent donc un suffixe de préversion semver au lieu de s’appuyer sur cette case. Épinglez toujours les versions en production.

## Branches

| Branche | Objectif |
| --- | --- |
| `master` | Développement et versions nocturnes. |
| `releases/**` | Branches à courte durée de vie destinées aux publications ponctuelles : la disponibilité générale `v3.0.0` et les correctifs urgents. Elles exigent la même intégration continue que `master`. |

Les versions hebdomadaires sont créées depuis `master`, mais **pas depuis son commit de tête** : le workflow sélectionne le commit le plus récent depuis la précédente publication hebdomadaire pour lequel tous les contrôles obligatoires (`v3 Build Results`, `Run Go Tests`, `Bundled runtime matches source`) ont réussi, puis applique le tag exactement à ce commit. Si aucun commit n’est au vert, aucune publication n’a lieu. Les versions nocturnes sont créées depuis `master` et ne produisent qu’un avertissement lorsque HEAD n’est pas au vert.

## Comment préparer une publication

Les droits de publication sont réservés au responsable du projet (consultez le document de gouvernance). Les contributeurs ne créent jamais de tag et ne publient jamais ; ils ajoutent des entrées à `v3/UNRELEASED_CHANGELOG.md`.

Toutes les opérations s’effectuent depuis **Nightly + Weekly Release v3** (`.github/workflows/nightly-release-v3.yml`) :

- **Nocturne (automatique)** à 15:00 UTC. Elle est ignorée lorsque le journal des modifications non publiées est vide.
- **Version hebdomadaire (automatique)** à 20:00, heure de Sydney. Les tâches cron de GitHub ne prennent en charge que l'heure UTC ; deux tâches cron (09:00 et 10:00 UTC) sont donc planifiées, et une tâche de garde ne poursuit l'exécution que lorsqu'il est 20:00 à Sydney (en tenant compte de l'heure d'été). La version hebdomadaire franchit la porte de validation de la publication, sélectionne le commit réussi le plus récent, incrémente la version mineure, attribue un tag à ce commit sans modifier le journal des modifications, puis met à jour uniquement les champs de version sur `master` (`version.txt` et les versions des packages du runtime), afin de conserver les mises à jour des dépendances et des fichiers de verrouillage fusionnées après le commit sélectionné. npm publie à partir du tag, et non de `master`.
- **Manuelle.** Exécutez le workflow avec `channel`, `dry_run` (activé par défaut) et, facultativement, `release_version`. Une simulation affiche la version, le type, les indicateurs GitHub, le dist-tag npm et les notes sans modifier GitHub ni npm.
- **Version GA ou correctif.** Définissez `release_version` (`v3.0.0`, `v3.1.1`). Pour un correctif destiné à une branche de version antérieure, vous pouvez également définir `ref: releases/...` ; son commit de publication est réappliqué par cherry-pick sur `master`.

Simulation locale : `cd v3/tasks/release && go run release.go --dry-run --channel weekly` ; `go run release.go --print-policy v3.2.0-nightly.20261102` affiche les indicateurs calculés.

### Mécanismes de sécurité

- Une version générée automatiquement doit être strictement supérieure, selon l’ordre SemVer où une préversion est antérieure à la version finale, à toutes les balises existantes, sinon l’outil la refuse.
- Une seule tâche de publication à la fois grâce à un groupe de concurrence partagé.
- La publication npm est déclenchée par la **balise** de publication et compile le commit balisé.

## Point de contrôle de publication

La publication **hebdomadaire**, jamais la nightly, refuse de créer une balise tant que l’un des éléments suivants existe :

1. des issues ou des pull requests ouvertes portant le libellé `P0` ou `release-blocker` ;
2. des issues ouvertes portant le libellé `regression` et créées depuis la dernière publication hebdomadaire ;
3. des alertes Dependabot ouvertes de gravité critique ou élevée, des alertes d’analyse du code ouvertes de gravité critique ou élevée, ainsi que toute alerte ouverte d’analyse des secrets, **uniquement pour le code distribué**. Les alertes situées exclusivement sous `v3/examples/`, `v2/examples/`, `website/`, `docs/`, `.github/` ou dans tout répertoire `examples/` sont ignorées (`SHIPPED_ALERT_IGNORE` dans `v3/tasks/releasegate/main.go`, **à examiner par le responsable**) ;
4. des avis de sécurité du dépôt toujours à l’état `triage`.

**Qui peut définir un blocage.** Un label n'est pris en compte que si son événement `labeled` le plus récent (sans événement `unlabeled` ultérieur) a été déclenché par un membre de l'équipe `wailsapp/maintainers` ou par un collaborateur disposant de l'autorisation `maintain`/`admin`. Le même label appliqué par toute autre personne est ignoré et consigné dans le journal, afin qu'aucune personne extérieure ne puisse bloquer ou débloquer une publication. Le label `high-priority` appliqué automatiquement n'est délibérément **pas** utilisé. Les comptes de bots tels que `taliesin-ai` passent la liste d'autorisation s'ils disposent de `maintain` ; c'est intentionnel, et l'agent de triage peut appliquer des labels de blocage en fournissant des éléments justificatifs.

**Comportement en cas d’échec.** Le point de contrôle échoue en mode **fermé** : si une API d’alertes ou la liste d’autorisation ne peut pas être lue, la publication hebdomadaire s’arrête. En cas d’échec, il ne crée aucune balise, inscrit la liste, liens compris, dans le résumé de la tâche et met à jour un commentaire dans l’unique issue « État de la publication hebdomadaire » lorsque la variable de dépôt `WEEKLY_STATUS_ISSUE` est définie.

**Configuration que le responsable doit effectuer sur GitHub** (rien n’est créé automatiquement) : créez les libellés `P0`, `release-blocker` et `regression` ; créez l’issue « État de la publication hebdomadaire » et définissez `WEEKLY_STATUS_ISSUE` sur son numéro ; accordez à `WAILS_REPO_TOKEN` un accès en lecture aux issues, aux pull requests, à l’administration (autorisation des collaborateurs), aux alertes Dependabot, aux alertes d’analyse du code, aux alertes d’analyse des secrets, aux avis de sécurité du dépôt et aux membres de l’organisation, ainsi qu’un accès en écriture aux issues pour le commentaire d’état. La valeur `GITHUB_TOKEN` par défaut ne peut pas lire ces éléments.

**Contournement.** `override_reason`, une entrée de déclenchement manuel du workflow, ignore le point de contrôle. Le motif, l’auteur de l’action et chaque blocage ignoré sont consignés dans le résumé et dans l’issue d’état. Seul le responsable des publications doit l’utiliser. GitHub ne peut pas restreindre une entrée à une seule personne ; exécutez donc la tâche de publication dans un environnement imposant des approbateurs (recommandé), afin qu’un contournement nécessite une approbation explicite.

## Voir aussi

- [Tests et intégration continue](/contributing/testing-ci/)
- [Triage par les mainteneurs](/contributing/maintainer-triage/)
