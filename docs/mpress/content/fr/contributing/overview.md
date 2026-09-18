---
title: "Présentation technique"
description: "Architecture générale et feuille de route du code source de Wails v3"
slug: "contributing/overview"
sourcePath: "contributing/overview.md"
---

## Bienvenue dans la documentation technique de Wails v3

Cette section ne traite **pas** des règles de la communauté ni de la procédure d’ouverture d’une pull request. Elle explique plutôt en détail **comment Wails v3 est construit** afin que vous puissiez rapidement vous repérer dans le code source et commencer à le modifier en toute confiance.

Que vous envisagiez de corriger le runtime, d’étendre la CLI, de créer de nouveaux modèles ou de simplement comprendre le fonctionnement interne, les pages suivantes vous apportent le contexte technique nécessaire.

---

## Architecture générale

@cards{cols="2"}
◇ Backend Go
Au cœur de chaque application Wails se trouve du code Go compilé en exécutable natif. Il prend en charge la logique applicative, l’intégration système et les opérations critiques pour les performances.

---
▤ Frontend web
L’interface utilisateur est écrite avec des technologies web standard (React, Vue, Svelte, Vanilla, etc.) et rendue par une WebView système légère (WebKit sous Linux/macOS, WebView2 sous Windows).

---
◆ Couche de liaison
Une passerelle en mémoire sans copie permet les appels **Go⇄JavaScript** avec conversion automatique des types, propagation des événements et transmission des erreurs.

---
▸ CLI et outils
`wails3` orchestre la création de projets, le serveur de développement avec rechargement à chaud, le regroupement des ressources, la compilation croisée et la création des paquets (deb, rpm, AppImage, msi, dmg…).

@end

---

## Vue d’ensemble de l’architecture

**Wails v3 – Flux de bout en bout**

**[Emplacement réservé au diagramme du flux de bout en bout]**

Le diagramme présente le **flux de bout en bout** :

1. La **CLI** pilote la génération, le serveur de développement, la compilation et la création des paquets.\
2. Le **système de liaison** produit le code d’adaptation qui permet au **frontend web** d’appeler le **backend Go**.\
3. Pendant le développement, le **serveur de ressources** transmet les requêtes au serveur de développement du framework ; en production, il sert les fichiers intégrés.\
4. À l’exécution, le **runtime de bureau** gère les fenêtres et les API du système d’exploitation, tandis que la **passerelle** achemine les messages entre Go et JavaScript.

---

## Contenu de cette documentation

| Sujet | Pourquoi est-ce important ? |
| --- | --- |
| **Organisation du code source** | Cartographie des répertoires `/v3` et des interactions entre les modules. |
| **Fonctionnement interne du runtime** | Gestion des fenêtres, API système, processeur de messages et couches d’adaptation aux plateformes. |
| **Serveur de ressources et de développement** | Mode de distribution des ressources web en développement et de leur intégration en production. |
| **Pipeline de compilation et de création des paquets** | Flux de travail fondé sur Taskfile, compilation multiplateforme et génération des programmes d’installation. |
| **Système de liaison** | Pipeline d’analyse statique qui génère des liaisons Go⇄TS avec sûreté des types. |
| **Système de modèles** | Architecture du générateur sur laquelle repose `wails3 init -t <framework>`. |
| **Tests et CI** | Infrastructure de tests unitaires et d’intégration, GitHub Actions et recommandations concernant le détecteur de situations de compétition. |
| **Extension de Wails** | Ajout de services, de modèles ou de sous-commandes à la CLI. |

Chaque page suivante approfondit ces sujets à l’aide d’exemples de code concrets, de diagrammes et de références aux fichiers sources concernés.

---

@note{type="info"}
Prérequis : vous devriez être à l’aise avec **Go 1.25+**, posséder des notions de base en TypeScript et connaître les outils modernes de compilation frontend. Si vous débutez en Go, envisagez de parcourir d’abord le tutoriel officiel.

@end

Bonne exploration, et bienvenue dans le fonctionnement interne de Wails v3 !
