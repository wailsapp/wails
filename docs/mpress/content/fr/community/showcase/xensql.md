---
title: "XenSQL"
description: "Environnement de travail SQL privilégiant le stockage local, avec prise en charge de plusieurs bases de données et outils avancés pour les requêtes"
slug: "community/showcase/xensql"
sourcePath: "community/showcase/xensql.md"
---

![Éditeur de tables XenSQL avec modification directe](/assets/showcase-images/xensql-1.png) ![Éditeur de requêtes XenSQL avec suggestions de tables et de colonnes](/assets/showcase-images/xensql-2.png) ![Requêtes transactionnelles multiples dans XenSQL](/assets/showcase-images/xensql-3.png)

**[XenSQL](https://github.com/Bare7a/XenSQL)** est un **environnement de travail SQL de bureau rapide, privilégiant le stockage local**, développé avec **Go, Wails et React**. Il réunit PostgreSQL, MySQL/MariaDB et SQLite dans une interface unique, épurée et d’apparence native, sans cloud, sans télémétrie et sans compte.

## Fonctionnalités principales

- **Éditeur SQL puissant** — Basé sur Monaco, avec autocomplétion intelligente tenant compte du schéma, exécution de plusieurs instructions, diffusion progressive des résultats et onglets de résultats propres à chaque instruction
- **Visionneuse de données avancée** — Inspecteur JSON interactif, éditeur de cellules tenant compte de la syntaxe (JSON, XML, HTML, texte), modification directe et inspection complète des enregistrements
- **Modification fluide des données** — Parcourez les tables, préparez les modifications directement, effectuez des opérations groupées et utilisez en toute sécurité `INSERT`/`UPDATE`/`DELETE` avec la prise en charge de `RETURNING`
- **Fonctionnalités de productivité** — Explorateur de schéma, requêtes enregistrées, historique des requêtes, recherche rapide (`Ctrl+P`) et flux de travail centré sur le clavier
- **Options d’exportation** — CSV, JSON, Markdown, instructions SQL INSERT

Entièrement hors ligne et portable. Toutes les données sont stockées localement dans un seul dossier `XenSQL-data/`, qui accompagne l’application.

**Bases de données prises en charge** : PostgreSQL, MySQL, MariaDB et SQLite (avec un mode lecture seule et des options de transport sécurisé).

Conçu pour les développeurs qui recherchent rapidité, clarté et maîtrise, sans la lourdeur des outils SQL traditionnels.
