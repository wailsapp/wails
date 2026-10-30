---
title: "Tutoriels"
description: "Apprenez Wails en créant des applications"
slug: "tutorials/overview"
sourcePath: "tutorials/overview.md"
---

Choisissez un tutoriel guidé pour apprendre Wails en réalisant un projet. Chaque tutoriel suit un projet jusqu’à l’obtention d’un résultat concret.

S’il s’agit de votre première visite, commencez par suivre [Installation](/quick-start/installation/) et [Votre première application](/quick-start/first-app/).

## Choisir un tutoriel

Après votre première application, suivez les tutoriels pour ordinateur dans cet ordre :

| Tutoriel | Ce que vous allez pratiquer |
| --- | --- |
| [Service de codes QR](/tutorials/01-creating-a-service/) | Exposer la logique Go par l’intermédiaire d’un service et l’appeler depuis le frontend |
| [Liste de tâches](/tutorials/02-todo-vanilla/) | Créer, mettre à jour et supprimer les données de l’application |
| [Application de notes](/tutorials/03-notes-vanilla/) | Utiliser des boîtes de dialogue natives et enregistrer des données dans des fichiers |
| [Application avec mise à jour automatique](/tutorials/04-self-update-a-wails-app/) | Rechercher et appliquer des mises à jour depuis GitHub Releases |

Pour suivre un tutoriel mobile facultatif, consultez [Votre première application mobile](/guides/mobile/first-mobile-app/). La prise en charge des plateformes mobiles est expérimentale.

Si vous travaillez sur une application existante et devez accomplir une tâche particulière, choisissez plutôt un [guide pratique](/quick-start/next-steps/#features). Pour comprendre les principes de conception, consultez les [Concepts fondamentaux](/concepts/architecture/).

@note{type="tip" title="Vous débutez avec Go ?"}
Suivez le [Tour de Go](https://go.dev/tour/) avant de commencer les tutoriels.

@end

## Service de codes QR

![Exemple de code QR](/assets/qr1.png)

Apprenez les principes fondamentaux des services Wails en créant un générateur de codes QR. Ce tutoriel vous présente les concepts essentiels pour organiser la logique de votre application en services réutilisables.

**Ce que vous apprendrez :**

- Créer et structurer un service Wails
- Gérer des dépendances Go externes
- Lier des méthodes Go à votre frontend
- Transmettre des données entre Go et JavaScript
- Organiser le code pour faciliter sa maintenance

**Idéal pour :** les personnes qui découvrent Wails et souhaitent comprendre l’architecture des services

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/01-creating-a-service/"><span>Commencer</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### Liste de tâches

![Application de liste de tâches](/assets/todo-app.png)

Créez une application complète de liste de tâches dotée d’une interface moderne et soignée. Ce tutoriel pratique vous enseigne les principaux modèles de conception de Wails au moyen d’une application concrète et réaliste utilisant du JavaScript natif.

**Ce que vous apprendrez :**

- Architecture orientée services avec gestion de l’état sûre pour les accès concurrents
- Opérations CRUD (création, lecture, mise à jour et suppression)
- Liaisons à typage sûr entre Go et JavaScript
- Créer des interfaces modernes sans la complexité d’un framework
- Modèles appropriés de gestion des erreurs et de validation

**Durée estimée :** environ 20 minutes

**Idéal pour :** votre première application Wails complète — parfaite pour comprendre les principes fondamentaux avant d’ajouter la complexité d’un framework

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/02-todo-vanilla/"><span>Commencer</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### Notes

![Application de notes](/assets/notes-app.png)

Créez une application dans le style d’Apple Notes, avec des boîtes de dialogue de fichiers natives et une fonction d’enregistrement automatique. Ce tutoriel présente des fonctionnalités propres aux applications de bureau, telles que les opérations sur les fichiers, les boîtes de dialogue natives et les modèles d’interface professionnels.

**Ce que vous apprendrez :**

- Boîtes de dialogue de fichiers natives (Enregistrer, Ouvrir, Informations)
- Persistance des données au format JSON
- Modèles d’enregistrement automatique avec temporisation anti-rebond
- Mises en page professionnelles à deux colonnes pour ordinateur
- Utiliser les opérations du système de fichiers en Go

**Durée estimée :** environ 30 minutes

**Idéal pour :** apprendre les fonctionnalités propres aux applications de bureau, telles que les opérations sur les fichiers et les boîtes de dialogue natives du système d’exploitation

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/03-notes-vanilla/"><span>Commencer</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### Application Wails avec mise à jour automatique

Ajoutez à une application Wails fraîchement créée avec `wails3 init` une fonction intégrée de mise à jour automatique, avec vérification signée des versions et remplacement du binaire en mode auxiliaire. GitHub Releases sert de source de mises à jour.

**Ce que vous apprendrez :**

- Comment intégrer `app.Updater` à une application Wails
- Configurer le fournisseur GitHub Releases
- Publier des versions avec `SHA256SUMS` pour vérifier les empreintes
- Ajouter une signature Ed25519 pour résister aux altérations
- Personnaliser la fenêtre par défaut avec du CSS, du HTML personnalisé ou votre propre interface
- Effectuer des vérifications périodiques en arrière-plan avec `CheckInterval`

**Durée estimée :** environ 25 minutes

**Idéal pour :** distribuer une application de bureau pouvant être mise à jour — couvre l’ensemble du pipeline de publication, et pas seulement l’API

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/04-self-update-a-wails-app/"><span>Commencer</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>
