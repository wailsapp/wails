---
title: "Feuille de route"
description: "État du projet Wails v3, fonctionnalités prévues et modalités de contribution"
slug: "status"
sourcePath: "status.md"
---

## État actuel : bêta

Consultez le [journal des modifications](/changelog/) pour connaître le dernier état du projet.

Notre objectif est de parvenir à une version v3.0 stable. Cette feuille de route présente les principales fonctionnalités et améliorations que nous devons mettre en œuvre avant la version finale. Notez qu’il s’agit d’un document évolutif susceptible d’être mis à jour en fonction de l’évolution des priorités ou de l’apparition de nouveaux enseignements.

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

Cette feuille de route est susceptible d’évoluer en fonction des retours de la communauté et des priorités du projet. Nous la mettrons régulièrement à jour pour refléter les avancées et les changements d’orientation. Signalez les problèmes reproductibles dans des issues ; proposez les nouvelles fonctionnalités au moyen d’une pull request [WEP (Wails Enhancement Proposal)](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md).
