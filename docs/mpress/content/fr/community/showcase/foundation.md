---
title: "Foundation"
description: "Un socle d’application de bureau Wails v3 prêt pour la production"
slug: "community/showcase/foundation"
sourcePath: "community/showcase/foundation.md"
---

![Réglages de disposition de Foundation avec neuf préréglages interchangeables](/assets/showcase-images/foundation-layouts.webp)

![Page d’accueil de Foundation avec le thème sombre](/assets/showcase-images/foundation-home-dark.webp)

![Graphiques X-Pro de Foundation avec le thème Obsidian](/assets/showcase-images/foundation-charts.webp)

**[Foundation](https://github.com/FxRayHughes/foundation)** est un **socle d’application de bureau pour Wails v3** prêt pour la production, construit avec **Go, React 19, TypeScript et Mantine 9**. Il fournit à un nouveau projet les briques qui demandent habituellement des semaines de travail — un backend Go en couches, une fenêtre adaptée à chaque plateforme, la persistance, les thèmes, l’internationalisation et des tests de bout en bout — pour passer directement à la logique métier.

## Points forts

- **Une fenêtre native sur chaque plateforme** — barre de titre sans bordure dessinée par le frontend sous Windows et Linux ; barre de titre immersive sous macOS qui conserve les boutons natifs de la fenêtre.
- **Neuf dispositions interchangeables** — barre d’icônes, barre latérale, navigation supérieure, barre de titre fusionnée, barre à droite, dock inférieur, dock flottant déplaçable, carte en retrait et mode focus, le tout sur une structure CSS Grid ; un préréglage personnalisé s’enregistre en une ligne.
- **Thèmes et internationalisation** — préréglages clair, sombre et Obsidian, plus un éditeur de palette personnalisée ; anglais et chinois inclus ; le thème et la langue restent synchronisés entre la fenêtre principale, les fenêtres enfants et le panneau de la zone de notification.
- **Persistance SQLite** — GORM AutoMigrate sans fichiers de migration, emplacement de la base de données au choix de l’utilisateur et graphiques d’occupation par table.
- **Intégration au bureau** — panneau dans la zone de notification, fenêtres enfants, boîtes de dialogue natives de fichiers et de messages, et export CSV / XLSX via la boîte de dialogue native « Enregistrer sous ».
- **Utilitaires robustes** — client HTTP, sous-processus rattachés à un Job Object Windows, chiffrement AES-GCM, journaux avec rotation et écritures atomiques.
- **Testé de bout en bout** — Playwright pilote une build Wails en mode serveur (vrai backend Go et vrai frontend), et la même suite génère les captures d’écran du README.

Grâce à Wails, le tout tient dans un seul binaire par plateforme : les services Go sont appelés directement via les liaisons générées, tandis que le frontend suit une structure MVVM stricte qui reste maintenable à mesure que les fonctionnalités s’ajoutent.

Disponible pour Windows, macOS et Linux.

[GitHub](https://github.com/FxRayHughes/foundation) |
[Versions](https://github.com/FxRayHughes/foundation/releases/latest)
