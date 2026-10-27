---
title: "Feuille de route"
description: "État du projet Wails v3, fonctionnalités prévues et modalités de contribution"
slug: "status"
sourcePath: "status.md"
---

## Wails v3.0.0 — Version stable

Consultez le [journal des modifications](/changelog/) pour connaître le dernier état du projet.

Cette documentation concerne Wails v3.0.0. Après v3.0.0, les correctifs et les améliorations sont publiés chaque semaine, avec, dans l’intervalle, des versions préliminaires nocturnes datées et des correctifs d’urgence pour les problèmes critiques. Consultez la [politique de publication](/contributing/release-policy/).

## Compatibilité avec les environnements de bureau

Le contrat de compatibilité de la version v3.0.0 couvre les applications de bureau :

| Plateforme | Cibles prises en charge | Configuration requise et remarques |
| --- | --- | --- |
| Windows | amd64 et arm64 | Environnement d’exécution WebView2 |
| macOS | Intel et Apple Silicon | Versions de macOS et de WebKit indiquées dans le guide d’installation |
| Linux | amd64 et arm64 | GTK4 + WebKitGTK 6.0 par défaut ; GTK3 + WebKit2GTK 4.1 reste une option héritée `-tags gtk3` jusqu’à la version v3.0.x incluse et est supprimé dans la version v3.1 |

Pour le développement, toutes les cibles nécessitent Go 1.26 ou une version ultérieure. La prise en charge d’Android et d’iOS reste expérimentale et ne fait pas partie du contrat de compatibilité stable pour les environnements de bureau.

## Comment contribuer

- Testez Wails v3.0.0 et signalez les bogues reproductibles
- Contribuez à la documentation et aux exemples
- Participez aux discussions et donnez votre avis sur les projets de WEP
- Soumettez des pull requests pour les corrections de bogues, la documentation ou les WEP acceptées

Les contributions de la communauté sont les bienvenues. Si vous souhaitez participer à la réalisation de ces objectifs, rejoignez les discussions de la communauté. Les propositions de nouvelles fonctionnalités doivent faire l’objet d’une pull request WEP, et non d’un ticket de demande de fonctionnalité.

## Retours et mises à jour

Cette feuille de route est susceptible d’évoluer en fonction des retours de la communauté et des priorités du projet. Nous la mettrons régulièrement à jour pour refléter l’avancement et les changements d’orientation. Signalez les problèmes reproductibles sous forme de tickets ; proposez les nouvelles fonctionnalités au moyen d’une pull request [WEP (proposition d’amélioration de Wails)](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/README.md).
