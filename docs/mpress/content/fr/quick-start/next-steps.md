---
title: "Étapes suivantes"
description: "Pour aller plus loin après avoir créé votre première application"
slug: "quick-start/next-steps"
sourcePath: "quick-start/next-steps.md"
---

Vous avez créé votre première application Wails. Suivez le parcours d’apprentissage ci-dessous ou accédez directement à un guide si vous savez déjà ce que vous souhaitez créer.

## En savoir plus

Continuez votre apprentissage en réalisant une application complète à la fois :

1. **[Créez un service de génération de codes QR](/tutorials/01-creating-a-service/).** Entraînez-vous à exposer la logique Go à votre frontend par l’intermédiaire d’un service.
2. **[Créez une liste de tâches](/tutorials/02-todo-vanilla/).** Ajoutez, modifiez et supprimez les données de l’application au moyen d’appels de service.
3. **[Créez une application de prise de notes](/tutorials/03-notes-vanilla/).** Manipulez des fichiers, des boîtes de dialogue natives et des données enregistrées.

Lorsque vous serez prêt à explorer la gestion des versions, suivez le [tutoriel sur les applications à mise à jour automatique](/tutorials/04-self-update-a-wails-app/). Consultez [tous les tutoriels](/tutorials/overview/) pour connaître les acquis de chaque leçon.

## Fonctionnalités

**Vous développez déjà votre propre application ?** Consultez le guide pratique correspondant à la tâche qui vous occupe :

| Je souhaite… | Commencez ici |
| --- | --- |
| Appeler du code Go depuis mon frontend | [Appeler du code Go depuis JavaScript](/features/bindings/methods/) |
| Ouvrir ou gérer une fenêtre | [Principes de base des fenêtres](/features/windows/basics/) |
| Ajouter des menus natifs ou une icône dans la zone de notification | [Menus de l’application](/features/menus/application/) ou [menus de la zone de notification](/features/menus/systray/) |
| Demander à l’utilisateur d’ouvrir ou d’enregistrer un fichier | [Boîtes de dialogue de sélection de fichiers](/features/dialogs/file/) |
| Envoyer des événements entre les composants | [Système d’événements](/features/events/system/) |
| Copier du texte ou accepter des fichiers déposés | [Presse-papiers](/features/clipboard/basics/) ou [glisser-déposer de fichiers](/features/drag-and-drop/files/) |
| Ajouter des raccourcis clavier | [Raccourcis clavier](/features/keyboard/shortcuts/) |
| Compiler et empaqueter mon application | [Compilation d’applications](/guides/build/building/), puis empaquetage pour [Windows](/guides/build/windows/), [macOS](/guides/build/macos/) ou [Linux](/guides/build/linux/) |

Ces pages se trouvent sous **Guides pratiques** dans la barre latérale. Vous pouvez ouvrir la rubrique pertinente sans lire toute la section dans l’ordre.

## Comprendre le fonctionnement de Wails

Consultez les **Concepts fondamentaux** pour comprendre les choix de conception qui sous-tendent le code :

- [Fonctionnement de Wails](/concepts/architecture/) explique comment le backend Go, le frontend et la WebView native sont reliés.
- [Pont entre Go et le frontend](/concepts/bridge/) explique comment les appels passent entre Go et JavaScript.
- [Cycle de vie de l’application](/concepts/lifecycle/) explique le démarrage et l’arrêt.
- [Système de compilation](/concepts/build-system/) explique le rôle des tâches du projet et des outils du frontend.

## Consulter une API

Consultez la [Référence de l’API](/reference/overview/) lorsque vous recherchez un type, une méthode, une option ou une commande plutôt qu’une leçon. Pour travailler en ligne de commande, commencez par la [vue d’ensemble des commandes de la CLI](/reference/cli/), puis consultez les [options et exemples de la CLI](/guides/cli/) pour obtenir des informations détaillées sur chaque commande.

## Obtenir de l’aide

Consultez la [foire aux questions](/faq/), parcourez les [exemples fonctionnels](https://github.com/wailsapp/wails/tree/releases/v3-beta/v3/examples) ou posez votre question sur [Discord](https://discord.gg/JDdSxwjhGf). Les problèmes propres à chaque plateforme sont regroupés sous **Dépannage** dans la barre latérale.
