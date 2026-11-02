---
title: "Zentro"
description: "IDE SQL de bureau centré sur le projet pour PostgreSQL, MySQL, SQL Server et SQLite"
slug: "community/showcase/zentro"
sourcePath: "community/showcase/zentro.md"
---

![Éditeur SQL de Zentro avec saisie semi-automatique selon le schéma, l’explorateur de schéma et un panneau d’enregistrement](/assets/showcase-images/zentro.png)

**[Zentro](https://zentro-ide.com)** est un **IDE SQL de bureau centré sur le projet** pour PostgreSQL, MySQL, SQL Server et SQLite, développé avec **Go, Wails v2 et React**. Chaque projet sait vers quel environnement il pointe, et Zentro demande confirmation avant toute opération irréversible : la requête destinée à la préproduction n’atterrit pas en production.

## Points forts

- **Conscience de l’environnement** - les projets portent leurs environnements (développement, préproduction, production) avec une surcouche de production bien visible, une classification des risques des instructions et un aperçu des lignes affectées avant une écriture destructrice
- **Éditeur SQL** - saisie semi-automatique selon le schéma, diagnostics en ligne et corrections rapides
- **Explorateur de schéma et ERD** - parcourez tables, colonnes, clés et index, et visualisez les relations
- **Suivi basé sur Git** - scripts enregistrés et instantanés de schéma versionnés comme du code, avec historique et différences
- **Tunnels SSH et terminal intégré** - accédez à une base via un tunnel et gardez un terminal local ou SSH à côté de vos requêtes
- **Assistant IA avec votre propre clé** - se connecte au point d’accès du modèle que vous configurez, jamais via un serveur opéré par Zentro, avec en plus un serveur MCP pour les agents externes
- **Priorité au clavier** - navigation par focus dans l’explorateur et la grille de résultats, palette de commandes et raccourcis en séquence
- **Confidentialité par défaut** - les identifiants des bases de données sont conservés dans le trousseau de votre système d’exploitation

Zentro est livré sous la forme d’un binaire natif unique par plateforme pour Windows, macOS et Linux. Le cœur est open source sous GPL-3.0 et une offre Pro payante facultative ajoute des garde-fous de production et davantage, avec un essai gratuit de 30 jours ; le code source de la Community Edition est public.

[Site web](https://zentro-ide.com) |
[GitHub](https://github.com/ZentroIDE/Zentro) |
[Télécharger](https://zentro-ide.com/download)
