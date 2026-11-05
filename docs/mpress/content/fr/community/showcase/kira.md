---
title: "Kira"
description: "Un client de bureau natif pour macOS qui réunit AWS — ECS, RDS, S3, DynamoDB et bien plus encore — dans une interface unique pilotée au clavier"
slug: "community/showcase/kira"
sourcePath: "community/showcase/kira.md"
---

**[Kira](https://kira.thiennguyen.dev)** est un **client de bureau natif pour macOS destiné à AWS** — *« AWS, sans les complications »*. Développé avec **Go, Wails et React + TypeScript**, il centralise les opérations AWS dans une application unique, pilotée au clavier. Authentifiez-vous une seule fois via AWS SSO et gérez l’infrastructure de nombreux comptes sans jamais avoir à passer de la console à la CLI ni à une multitude d’outils de base de données.

Tout est parti d’une frustration personnelle : devoir passer sans cesse des onglets du navigateur aux commandes `aws`, puis à un client SQL distinct, simplement pour déployer une modification, était laborieux. Kira rassemble ces flux de travail dans une fenêtre native et rapide : connectez-vous, choisissez un compte et accédez à tout ce dont vous avez besoin d’une simple frappe.

![Vue d’ensemble multicomptes de Kira — comptes de production, de préproduction et de développement réunis au même endroit](/assets/showcase-images/kira_screenshot_1.png)

## Fonctionnalités principales

- **AWS SSO multicomptes** - Connectez-vous une seule fois, puis passez d’un compte ou d’une région à l’autre depuis un emplacement unique
- **ECS** - Parcourez les clusters, les services et les tâches ; redéployez, dimensionnez et restaurez une version antérieure des services ; consultez les définitions de tâches ; surveillez les métriques des services ; et ouvrez un shell interactif via ECS Exec
- **Bases de données** - Exécutez des requêtes SQL sur RDS, interrogez et analysez DynamoDB, et connectez-vous à PostgreSQL, MySQL et Redshift, avec une redirection SSH sécurisée et des identifiants stockés dans le trousseau macOS
- **S3** - Parcourez les compartiments et les préfixes ; prévisualisez, téléversez, téléchargez, copiez, renommez et supprimez des objets ; et créez des dossiers
- **Secrets Manager** - Répertoriez les secrets et récupérez leurs valeurs pour le compte actif
- **CloudWatch Logs** - Suivez en temps réel les flux de journaux et effectuez-y des recherches
- **Smart Query** - Génération SQL facultative assistée par l’IA et reposant sur la CLI `claude`
- **Extensions** - Installez des paquets `.kext` personnalisés qui ajoutent des boutons d’action reposant sur de petits scripts Go
- **Navigation rapide** - Un raccourci clavier global pour afficher l’application, une palette de commandes `Cmd+K` et des liens profonds `kira://`

## Présentation détaillée

Surveillez vos services ECS en temps réel — état des tâches, processeur et mémoire, état du déploiement — puis redéployez-les, dimensionnez-les ou restaurez une version antérieure sans quitter la liste.

![Kira parcourant les services ECS avec des métriques en temps réel](/assets/showcase-images/kira_screenshot_17.png)

Parcourez S3 comme avec un gestionnaire de fichiers. Prévisualisez les objets, examinez leurs métadonnées et leurs versions, puis téléversez-les, téléchargez-les, renommez-les ou supprimez-les directement.

![Navigateur d’objets S3 de Kira avec aperçu des objets et métadonnées](/assets/showcase-images/kira_screenshot_5.png)

Distribué sous la forme d’un `.dmg` signé et notarié pour macOS.

[Visiter Kira](https://kira.thiennguyen.dev) | [Lire la documentation](https://docs.kira.thiennguyen.dev)
