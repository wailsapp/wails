---
title: "CFN Tracker"
description: "Une application de bureau créée avec Wails"
slug: "community/showcase/cfntracker"
sourcePath: "community/showcase/cfntracker.md"
---

![CFN Tracker](/assets/showcase-images/cfntracker.webp)

[CFN Tracker](https://github.com/williamsjokvist/cfn-tracker) — Suivez en direct les matchs de n’importe quel profil CFN Street Fighter 6 ou V. Consultez [le site web](https://cfn.williamsjokvist.se/) pour commencer.

## Fonctionnalités

- Suivi des matchs en temps réel
- Stockage des journaux de matchs et des statistiques
- Prise en charge de l’affichage des statistiques en direct dans OBS via une source Navigateur
- Prise en charge de SF6 et de SFV
- Possibilité pour les utilisateurs de créer leurs propres thèmes de navigateur OBS avec CSS

### Principales technologies utilisées avec Wails

- [Task](https://github.com/go-task/task) — encapsule la CLI de Wails pour faciliter l’utilisation des commandes courantes
- [React](https://github.com/facebook/react) — choisi pour la richesse de son écosystème (radix, framer-motion)
- [Bun](https://github.com/oven-sh/bun) — utilisé pour la rapidité de sa résolution des dépendances et de ses temps de compilation
- [Rod](https://github.com/go-rod/rod) — automatisation d’un navigateur sans interface graphique pour l’authentification et la détection des changements par interrogation périodique
- [SQLite](https://github.com/mattn/go-sqlite3) — utilisé pour stocker les matchs, les sessions et les profils
- [Événements envoyés par le serveur](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events) — flux HTTP servant à envoyer les mises à jour du suivi aux sources Navigateur d’OBS
- [i18next](https://github.com/i18next/) — avec un connecteur backend permettant de fournir les objets de localisation depuis la couche Go
- [xstate](https://github.com/statelyai/xstate) — machines à états pour les processus d’authentification et de suivi
