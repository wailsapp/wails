---
title: "Présentation technique"
description: "Architecture générale et feuille de route du code source de Wails v3"
slug: "contributing/overview"
sourcePath: "contributing/overview.md"
---

## Bienvenue dans la documentation technique de Wails v3

Cette section ne traite **pas** des règles de la communauté ni de la manière d’ouvrir une pull request. Elle explique plutôt **comment Wails v3 est construit**, afin que vous puissiez rapidement vous repérer dans le code source et commencer à le modifier en toute confiance.

Que vous envisagiez de corriger le runtime, d’étendre la CLI, de créer de nouveaux modèles ou simplement de comprendre le fonctionnement interne, les pages suivantes vous fourniront le contexte technique nécessaire.

---

## Architecture générale

@cards{cols="2"}
◇ Backend Go
Le cœur de chaque application Wails est constitué de code Go compilé en exécutable natif. Il prend en charge la logique applicative, l’intégration au système et les opérations critiques pour les performances.

---
▤ Frontend web
L’interface utilisateur est écrite avec des technologies web standard (React, Vue, Svelte, JavaScript natif, etc.) et affichée par une WebView système légère (WebKit sous Linux/macOS, WebView2 sous Windows).

---
◆ Couche de liaison
Un pont en mémoire sans copie permet les appels **Go⇄JavaScript** avec conversion automatique des types, propagation des événements et transmission des erreurs.

---
▸ CLI et outils
`wails3` orchestre la création de projets, le serveur de développement avec rechargement à chaud, le regroupement des ressources, la compilation croisée et la création de paquets (deb, rpm, AppImage, msi, dmg, etc.).

@end

---

## Vue d’ensemble de l’architecture

**Wails v3 – Flux de bout en bout**

**[Emplacement réservé au diagramme du flux de bout en bout]**

Le diagramme illustre le **flux de bout en bout** :

1. La **CLI** pilote la génération, le serveur de développement, la compilation et la création de paquets.\
2. Le **système de liaisons** produit le code d’adaptation qui permet au **frontend web** d’appeler le **backend Go**.\
3. Pendant le développement, le **serveur de ressources** fait office de proxy vers le serveur de développement du framework ; en production, il sert les fichiers intégrés.\
4. Lors de l’exécution, le **runtime de bureau** gère les fenêtres et les API du système d’exploitation, tandis que le **pont** assure le transfert des messages entre Go et JavaScript.

---

## Contenu de cette documentation

| Sujet | Pourquoi est-ce important ? |
| --- | --- |
| **Organisation du code source** | Cartographie des répertoires `/v3` et des interactions entre les modules. |
| **Fonctionnement interne du runtime** | Gestion des fenêtres, API système, processeur de messages et couches d’adaptation aux plateformes. |
| **Serveur de ressources et de développement** | Mode de mise à disposition des ressources web en développement et de leur intégration en production. |
| **Pipeline de compilation et de création de paquets** | Workflow fondé sur Taskfile, compilation multiplateforme et génération de programmes d’installation. |
| **Système de liaisons** | Pipeline d’analyse statique qui génère des liaisons Go⇄TS avec sûreté des types. |
| **Système de modèles** | Architecture du générateur qui alimente `wails3 init -t <framework>`. |
| **Tests et CI** | Infrastructure de tests unitaires et d’intégration, GitHub Actions et recommandations sur le détecteur de situations de compétition. |
| **Extension de Wails** | Ajout de services, de modèles ou de sous-commandes à la CLI. |

Chaque page suivante approfondit ces domaines à l’aide d’exemples de code concrets, de diagrammes et de références aux fichiers sources concernés.

---

@note{type="info"}
Prérequis : vous devez être à l’aise avec **Go 1.26+**, posséder des connaissances de base en TypeScript et maîtriser les outils modernes de compilation frontend. Si vous débutez en Go, envisagez de parcourir d’abord le tutoriel officiel.

@end

Bonne exploration — et bienvenue dans le fonctionnement interne de Wails v3 !
