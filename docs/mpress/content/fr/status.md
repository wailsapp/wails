---
title: "État du projet"
description: "Compatibilité de la bêta Wails v3, prise en charge de la sécurité et guide de mise à niveau"
slug: "status"
sourcePath: "status.md"
---

## État actuel : bêta

Consultez le [journal des modifications](/changelog/) pour connaître le dernier état du projet.

Notre objectif est une version v3.0 stable. Wails v2 reste la version stable actuelle et continue de recevoir des correctifs. Testez les versions bêta avec votre application avant le déploiement.

## Engagement de compatibilité de la version bêta

Le contrat de la version bêta de v3 couvre les applications de bureau :

| Plateforme | Cibles prises en charge | Prérequis et remarques |
| --- | --- | --- |
| Windows | amd64 et arm64 | Environnement d’exécution WebView2 |
| macOS | Intel et Apple Silicon | Les versions de macOS et de WebKit indiquées dans le guide d’installation |
| Linux | amd64 et arm64 | GTK4 + WebKitGTK 6.0 par défaut ; GTK3 + WebKit2GTK 4.1 reste une option héritée avec `-tags gtk3` jusqu’à la version v3.0.x incluse et est supprimée dans la version v3.1 |

Pour le développement, toutes les cibles nécessitent Go 1.25 ou une version ultérieure. La prise en charge d’Android et d’iOS est expérimentale et ne bloque pas la version bêta pour ordinateurs de bureau. Les API bêta visent la stabilité, mais les défauts des préversions et les changements explicitement annoncés peuvent encore être corrigés avant la version v3.0.0.

## Comment contribuer

- Testez la dernière version bêta et signalez les bogues reproductibles
- Contribuez à la documentation et aux exemples
- Participez aux discussions et donnez votre avis sur les projets de WEP
- Soumettez des pull requests pour les corrections de bogues, la documentation ou les WEP acceptées

Les contributions de la communauté sont les bienvenues. Si vous souhaitez participer à la réalisation de ces objectifs, rejoignez les discussions de la communauté. Les propositions de nouvelles fonctionnalités doivent faire l’objet d’une pull request WEP, et non d’une issue de demande de fonctionnalité.

## Retours et mises à jour

Signalez les problèmes reproductibles dans des issues ; proposez les nouvelles fonctionnalités au moyen d’une pull request [WEP (Wails Enhancement Proposal)](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md).

## Utiliser la bêta

Fixez des versions exactes pour la CLI, le module Go et le runtime frontend plutôt que de suivre `latest`. Pour un projet alpha existant, suivez le [guide de mise à niveau d’alpha vers bêta](/migration/alpha-to-beta/).

La [politique de sécurité](https://github.com/wailsapp/wails/blob/master/SECURITY.md) indique que les versions bêta de v3 sont prises en charge, contrairement aux versions alpha. Signalez les vulnérabilités via le [signalement privé de vulnérabilités](https://github.com/wailsapp/wails/security/advisories/new), et non dans des issues publiques.

## Travaux suivis

- [Bogues ouverts portant le label v3](https://github.com/wailsapp/wails/issues?q=is%3Aissue+is%3Aopen+label%3ABug+label%3Av3)
- [Issues v3 ouvertes portant le label P0 ou P1](https://github.com/wailsapp/wails/issues?q=is%3Aissue+is%3Aopen+label%3Av3+label%3AP0%2CP1)
- [Jalons des versions](https://github.com/wailsapp/wails/milestones)

Ces recherches en direct dépendent des labels des issues ; elles ne constituent ni une liste exhaustive ni un engagement sur une date ou un périmètre de version. Lisez les issues pour évaluer leur incidence sur votre projet.
