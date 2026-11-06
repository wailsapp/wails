---
title: "Foundation"
description: "Socle d’application de bureau pour Wails v3.1 avec React 19 et Mantine 9"
slug: "community/showcase/foundation"
sourcePath: "community/showcase/foundation.md"
---

![Réglages de disposition de Foundation avec neuf préréglages interchangeables](/assets/showcase-images/foundation-layouts.webp)

![Page d’accueil de Foundation avec le thème sombre](/assets/showcase-images/foundation-home-dark.webp)

![Graphiques X-Pro de Foundation avec le thème Obsidian](/assets/showcase-images/foundation-charts.webp)

**[Foundation](https://github.com/FxRayHughes/foundation)** est un **socle d’application de bureau pour Wails v3.1**, construit avec **Go, React 19, TypeScript et Mantine 9**. Il fournit un backend Go en couches, une fenêtre adaptée à chaque plateforme, la persistance SQLite, des thèmes, l’internationalisation et des tests de bout en bout comme point de départ pour de nouveaux projets.

## Fonctionnalités

- **Fenêtres adaptées à chaque plateforme** — fenêtre sans bordure avec barre de titre dessinée par le frontend sous Windows et Linux ; sous macOS, barre de titre immersive qui conserve les boutons natifs de la fenêtre.
- **Neuf préréglages de disposition** — barre d’icônes, barre latérale, navigation supérieure, barre de titre fusionnée, barre à droite, dock inférieur, dock flottant déplaçable, carte en retrait et mode focus, sur une structure CSS Grid ; des préréglages personnalisés peuvent être enregistrés.
- **Thèmes et internationalisation** — préréglages clair, sombre et Obsidian, et un éditeur de palette personnalisée ; interface en anglais et en chinois ; le thème et la langue restent synchronisés entre la fenêtre principale, les fenêtres enfants et le panneau de la zone de notification.
- **Persistance SQLite** — GORM AutoMigrate sans fichiers de migration, emplacement de la base de données au choix de l’utilisateur et graphiques d’occupation par table.
- **Intégration au bureau** — panneau dans la zone de notification, fenêtres enfants, boîtes de dialogue natives de fichiers et de messages, et export CSV / XLSX via la boîte de dialogue native « Enregistrer sous ».
- **Utilitaires** — client HTTP, sous-processus rattachés à un Job Object Windows, chiffrement AES-GCM, journaux avec rotation et écritures atomiques.
- **Tests de bout en bout** — Playwright s’exécute sur une build Wails en mode serveur (vrai backend Go et vrai frontend) ; la même suite génère les captures d’écran du README.

Le frontend appelle les services Go via les liaisons générées par Wails, et chaque plateforme produit un seul binaire. Le frontend suit une structure MVVM, et le README contient un guide pas à pas pour renommer le socle en votre propre projet.

Fonctionne sous Windows, macOS et Linux.

[GitHub](https://github.com/FxRayHughes/foundation) |
[Versions](https://github.com/FxRayHughes/foundation/releases/latest)
