---
title: "Architecture de Wails v3"
description: "Diagrammes détaillés et explications de chaque composant interne de Wails v3"
slug: "contributing/architecture"
sourcePath: "contributing/architecture.md"
---

Wails v3 est un **framework de bureau full-stack** composé d’un environnement d’exécution Go, d’un pont JavaScript, d’une chaîne d’outils pilotée par des tâches et d’un ensemble de modèles qui vous permettent de livrer des applications natives fondées sur des technologies web modernes.

Cette page présente la *vue d’ensemble* en quatre diagrammes :

1. **Architecture globale** – comment tous les sous-systèmes sont interconnectés\
2. **Flux d’exécution** – ce qui se passe lorsque JavaScript appelle Go et inversement\
3. **Développement et production** – les deux modes du serveur de ressources\
4. **Implémentations propres aux plateformes** – où réside le code propre à chaque système d’exploitation\

---

## 1 · Architecture globale

**Wails v3 – architecture générale**

**[Espace réservé au diagramme de l’architecture générale]**

---

## 2 · Flux des appels à l’exécution

**Exécution – parcours des appels JavaScript ⇄ Go**

**[Espace réservé au diagramme du flux des appels à l’exécution]**

Points clés :

- **Aucun HTTP ni IPC** – le pont utilise le canal en mémoire du WebView natif\
- **Identifiants de méthode** – un hachage FNV déterministe permet une recherche en O(1) dans Go\
- **Promesses** – les erreurs se propagent sous forme de rejets avec pile d’appels et code

---

## 3 · Flux des ressources en développement et en production

**Serveur de ressources en développement ↔ production**

**[Espace réservé au diagramme du flux des ressources]**

- En mode **développement**, le serveur transmet les chemins inconnus au serveur de rechargement à chaud du framework et sert les ressources statiques depuis le disque.
- En mode **production**, la même API s’appuie sur `go:embed`, ce qui produit un binaire sans dépendances.

---

## 4 · Séparation de l’environnement d’exécution selon la plateforme

**Fichiers d’exécution propres à chaque système d’exploitation**

**[Espace réservé au diagramme de la séparation par plateforme]**

Chaque fonctionnalité suit ce modèle :

1. **Interface commune** dans `pkg/application`\
2. Entrée du **processeur de messages** dans `pkg/application/messageprocessor_*.go`\
3. **Implémentation propre à chaque système d’exploitation** dans `pkg/application/*_{darwin,linux,windows}.go` (par exemple `webview_window_darwin.go`, `clipboard_linux.go`, `dialogs_windows.go`, `systemtray_*.go`, `mainthread_*.go`), protégée par des contraintes de compilation. Sous Linux, le pont cgo se trouve également dans `linux_cgo.go` / `linux_cgo_gtk4.{go,c,h}`.

`internal/runtime/` ne contient que le petit code de raccordement `runtime{_darwin,_linux,_windows,_android,_dev,_prod}.go` associé aux contraintes de compilation, ainsi que l’environnement d’exécution JS intégré sous `internal/runtime/desktop/`.

`internal/capabilities/` sert à déclarer les ensembles de capacités propres à chaque plateforme, mais il n’existe aucune sentinelle `ErrCapability` : l’activation conditionnelle des fonctionnalités repose sur de simples contraintes de compilation et sur les valeurs de retour des stubs propres aux plateformes (par exemple `nil` ou des erreurs propres aux fonctionnalités).

---

## Résumé

Ces diagrammes indiquent **où réside le code**, **comment les données circulent** et **quelles couches assument quelles responsabilités**. Gardez-les à portée de main lorsque vous explorerez les pages détaillées qui suivent : ils constituent votre carte de l’arborescence des sources de Wails v3.
