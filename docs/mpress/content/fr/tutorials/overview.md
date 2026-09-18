---
title: "Tutoriels"
description: "Apprenez à utiliser Wails en créant des applications"
slug: "tutorials/overview"
sourcePath: "tutorials/overview.md"
---

Des tutoriels pas à pas qui vous enseignent les concepts de Wails en vous faisant créer des applications complètes. Chaque tutoriel comprend du code fonctionnel, des explications et des modèles pratiques.

@note{type="tip" title="Vous débutez avec Go ?"}
Suivez le [Tour de Go](https://go.dev/tour/) avant de commencer les tutoriels.

@end

## Service de codes QR

![Exemple de code QR](/assets/qr1.png)

Découvrez les principes fondamentaux des services Wails en créant un générateur de codes QR. Ce tutoriel présente les concepts essentiels pour organiser la logique de votre application en services réutilisables.

**Ce que vous apprendrez :**

- Créer et structurer un service Wails
- Gérer les dépendances Go externes
- Lier des méthodes Go à votre interface utilisateur
- Transmettre des données entre Go et JavaScript
- Organiser le code pour faciliter sa maintenance

**Public visé :** les personnes qui utilisent Wails pour la première fois et souhaitent comprendre l’architecture des services

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/01-creating-a-service/"><span>Commencer</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### Liste de tâches

![Application de liste de tâches](/assets/todo-app.png)

Créez une application complète de liste de tâches dotée d’une interface moderne et élégante. Ce tutoriel pratique vous enseigne les principaux modèles de Wails au moyen d’une application concrète et réaliste utilisant du JavaScript natif.

**Ce que vous apprendrez :**

- Une architecture fondée sur des services, avec une gestion de l’état sûre pour les accès concurrents
- Les opérations CRUD (création, lecture, mise à jour et suppression)
- Des liaisons à typage sûr entre Go et JavaScript
- Créer des interfaces utilisateur modernes sans la complexité d’un framework
- Des modèles appropriés de gestion des erreurs et de validation

**Temps nécessaire :** environ 20 minutes

**Idéal pour :** créer votre première application Wails complète — parfait pour comprendre les principes fondamentaux avant d’ajouter la complexité d’un framework

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/02-todo-vanilla/"><span>Commencer</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### Notes

![Application de prise de notes](/assets/notes-app.png)

Créez une application inspirée d’Apple Notes, dotée de boîtes de dialogue de fichiers natives et d’une fonction d’enregistrement automatique. Ce tutoriel présente des fonctionnalités propres aux applications de bureau, telles que les opérations sur les fichiers, les boîtes de dialogue natives et les modèles d’interface utilisateur professionnels.

**Ce que vous apprendrez :**

- Les boîtes de dialogue de fichiers natives (Enregistrer, Ouvrir, Informations)
- La persistance des données au format JSON
- Des modèles d’enregistrement automatique avec temporisation anti-rebond
- Des mises en page professionnelles à deux colonnes pour les applications de bureau
- Utiliser les opérations du système de fichiers en Go

**Temps nécessaire :** environ 30 minutes

**Public visé :** les personnes souhaitant découvrir des fonctionnalités propres aux applications de bureau, telles que les opérations sur les fichiers et les boîtes de dialogue natives du système d’exploitation

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/03-notes-vanilla/"><span>Commencer</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### Application Wails à mise à jour automatique

Ajoutez à une application Wails un mécanisme intégré de mise à jour automatique, depuis un nouveau `wails3 init` jusqu’à la vérification des versions signées et au remplacement du binaire en mode auxiliaire. GitHub Releases sert de source pour les mises à jour.

**Ce que vous apprendrez :**

- Comment intégrer `app.Updater` à une application Wails
- Configurer le fournisseur GitHub Releases
- Publier des versions avec `SHA256SUMS` afin d’en vérifier l’empreinte
- Ajouter une signature Ed25519 pour protéger l’application contre les altérations
- Personnaliser la fenêtre par défaut avec du CSS, du HTML personnalisé ou une solution BYO
- Effectuer des vérifications périodiques en arrière-plan avec `CheckInterval`

**Temps nécessaire :** environ 25 minutes

**Idéal pour :** la livraison d’une application de bureau pouvant être mise à jour — couvre l’ensemble du processus de publication, et pas seulement l’API

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/04-self-update-a-wails-app/"><span>Commencer</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>
