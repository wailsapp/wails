---
title: "Pourquoi Wails ?"
description: "Découvrez pourquoi Wails est le choix idéal pour votre application de bureau"
slug: "quick-start/why-wails"
sourcePath: "quick-start/why-wails.md"
---

Wails associe **les performances et la simplicité de Go** à **la flexibilité des interfaces web modernes**, ce qui vous permet de créer de belles applications de bureau natives avec les outils que vous connaissez déjà.

## Des performances perceptibles par les utilisateurs

**Applications Wails :**

- **Fichiers binaires d’environ 15 Mo** (contre 150 Mo pour Electron)
- **Environ 10 Mo de mémoire au repos** (contre plus de 100 Mo pour Electron)
- **&lt;0.5 s au démarrage** (contre 2-3 s pour Electron)
- **Rendu natif** à l’aide de la WebView fournie par le système d’exploitation

Les utilisateurs perçoivent votre application comme rapide, légère et professionnelle.

## Expérience de développement

**Écrivez une fois, exécutez partout :**

- Une seule base de code Go pour Windows, macOS et Linux
- Utilisez n’importe quel framework web (React, Vue, Svelte ou JavaScript natif)
- Rechargement à chaud pendant le développement
- Liaisons TypeScript générées automatiquement à partir du code Go

Publiez plus vite avec moins de code à maintenir.

## Fonctionnalités prêtes pour la production

**Tout ce dont vous avez besoin :**

- Plusieurs fenêtres avec des cycles de vie indépendants
- Menus natifs (application, menus contextuels et zone de notification)
- Boîtes de dialogue de fichiers avec une interface native à chaque plateforme
- Intégration au système (notifications, presse-papiers et raccourcis clavier)
- Signature du code et empaquetage pour toutes les plateformes

Créez des applications professionnelles, pas des prototypes.

## Développement plus rapide

- **Une seule base de code, trois plateformes** — Écrivez une fois, puis compilez pour Windows, macOS et Linux
- **Mettez à profit vos compétences actuelles** — Go pour le backend, HTML/CSS/JS pour l’interface utilisateur
- **Retour instantané** – Rechargement à chaud pendant le développement, avec des temps de compilation de quelques secondes
- **Binaires compacts** – Des applications de 15 Mo pour des compilations, des téléchargements et des itérations plus rapides

## Quand choisir Wails

**Wails convient parfaitement aux cas suivants :**

- **Applications métier** (CRM, gestion des stocks, tableaux de bord, outils d’administration)
- **Outils de développement** (clients de base de données, outils de test d’API, outils de déploiement)
- **Applications de productivité** (prise de notes, gestionnaires de tâches, suivi du temps)
- **Outils de création** (éditeurs d’images, outils de traitement vidéo, utilitaires de conception)
- **Outils internes** (applications propres à l’entreprise, outils d’automatisation)

## Exemples de réussite concrets

@note{type="tip" title="Applications en production"}
Wails fait fonctionner de véritables applications utilisées par des milliers de personnes :

- **Outils de gestion de bases de données** dotés d’interfaces utilisateur complexes
- **Tableaux de bord financiers** traitant des données en temps réel
- **Outils de montage vidéo** offrant des performances natives
- **Utilitaires de développement** utilisés par des équipes d’ingénierie

[Voir la galerie →](/community/showcase/)

@end

## Fonctionnement de Wails

Contrairement à Electron, qui embarque un navigateur complet et l’environnement d’exécution Node.js, Wails adopte une approche fondamentalement différente : votre code Go est compilé en un binaire natif, tandis que votre interface utilisateur s’exécute dans la WebView intégrée au système d’exploitation. Cette architecture produit des binaires compacts, assure un démarrage rapide et réduit l’utilisation de la mémoire, donnant ainsi aux applications Wails le comportement d’applications natives.

### Architecture

Les applications Wails se composent de deux parties principales qui communiquent de manière transparente : un backend Go chargé de la logique métier et des opérations système, et un frontend web constituant votre interface utilisateur. La WebView fournie par le système d’exploitation affiche votre interface sans embarquer de navigateur, tandis que la couche de liaisons assure une communication avec typage sûr entre Go et JavaScript.

<div style="display: flex; justify-content: center; align-items: center; margin: 2rem 0;">
  <img src="/img/architecture.svg" alt="Architecture de Wails : un backend Go et votre interface utilisateur web compilés en un seul binaire natif, reliés par des bindings générés et rendus par la WebView du système d’exploitation" style="max-width: 640px; width: 100%;" />
</div>

Cette architecture simple permet au code JavaScript d’appeler directement des fonctions Go au moyen de liaisons générées automatiquement, tandis que Go peut renvoyer des événements et des données au frontend. Les deux couches communiquent par l’intermédiaire d’un pont en mémoire efficace dont la surcharge est inférieure à une milliseconde.

**Comment Wails atteint de telles performances :**

1. **Aucun environnement d’exécution embarqué** – Utilise le binaire compilé de Go
2. **WebView native** – Moteur de rendu fourni par le système d’exploitation
3. **Pont direct Go ↔ JS** – Communication en mémoire, sans surcharge réseau
4. **Binaire compilé** – Démarrage instantané, sans compilation JIT

## Étapes suivantes

Maintenant que vous connaissez les possibilités offertes par Wails, configurons votre environnement :

1. **Installez Wails** — Configurez votre environnement de développement en 5 minutes [Guide d’installation →](/quick-start/installation/)

2. **Créez votre première application** — Créez une application fonctionnelle et découvrez les principes de base [Tutoriel sur la première application →](/quick-start/first-app/)

3. **Explorez les fonctionnalités** — Découvrez ce que Wails peut apporter à votre application [Vue d’ensemble des fonctionnalités →](/quick-start/next-steps/)

---

**Vous avez encore des questions ?** Rejoignez notre [communauté Discord](https://discord.gg/JDdSxwjhGf) et posez-les directement à l’équipe.
