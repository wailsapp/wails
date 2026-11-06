---
title: "La voie vers Wails v3"
description: "Notes de version et annonces concernant Wails"
authors: ["leaanthony"]
tags: ["wails","v3"]
date: "2023-01-17"
slug: "blog/the-road-to-wails-v3"
image: "/assets/blog-images/multiwindow.webp"
sourcePath: "blog/the-road-to-wails-v3.md"
---

![capture d’écran du mode multifenêtre](/assets/blog-images/multiwindow.webp)

## Introduction

Wails est un projet qui simplifie le développement d’applications de bureau multiplateformes avec Go. Il utilise des composants de vue web natifs pour l’interface utilisateur — et non des navigateurs intégrés —, offrant ainsi à Go la puissance du système d’interface utilisateur le plus répandu au monde, tout en restant léger.

La version 2 est sortie le 22 septembre 2022 et a apporté de nombreuses améliorations, notamment :

- Développement en direct reposant sur le populaire projet Vite
- Fonctionnalités avancées de gestion des fenêtres et de création de menus
- Composant WebView2 de Microsoft
- Génération de modèles TypeScript reflétant vos structures Go
- Création d’un programme d’installation NSIS
- Builds obscurcis

À l’heure actuelle, Wails v2 fournit des outils puissants pour créer des applications de bureau multiplateformes riches en fonctionnalités.

Cet article de blog présente l’état actuel du projet et les améliorations que nous pouvons lui apporter à l’avenir.

## Où en sommes-nous ?

Il est incroyable de voir la popularité de Wails augmenter depuis la sortie de la v2. La créativité de la communauté et les merveilles qu’elle réalise avec Wails ne cessent de m’étonner. Cette popularité accrue attire davantage de regards sur le projet et, par conséquent, davantage de demandes de fonctionnalités et de rapports de bogues.

Au fil du temps, j’ai pu cerner certains des problèmes les plus urgents auxquels le projet est confronté. J’ai également identifié certains des facteurs qui freinent son évolution.

## Problèmes actuels

J’ai identifié les domaines suivants qui, selon moi, freinent l’évolution du projet :

- L’API
- Génération des liaisons
- Le système de build

### L’API

L’API permettant de créer une application Wails se compose actuellement de 2 parties :

- L’API d’application
- L’API d’exécution

Comme chacun sait, l’API d’application ne comporte que 1 fonction : `Run()`, qui accepte une multitude d’options régissant le fonctionnement de l’application. Cette API est très simple à utiliser, mais aussi très restrictive. Cette approche « déclarative » masque une grande partie de la complexité sous-jacente. Par exemple, aucune référence à la fenêtre principale n’est disponible : vous ne pouvez donc pas interagir directement avec elle. Pour cela, vous devez utiliser l’API d’exécution. Cela pose problème dès que vous souhaitez effectuer des opérations plus complexes, comme créer plusieurs fenêtres.

L’API d’exécution fournit de nombreuses fonctions utilitaires aux développeurs, notamment pour les opérations suivantes :

- Gestion des fenêtres
- Boîtes de dialogue
- Menus
- Événements
- Journaux

Plusieurs aspects de l’API d’exécution ne me satisfont pas. Tout d’abord, elle exige qu’un « contexte » soit transmis de fonction en fonction. C’est à la fois frustrant et déroutant pour les nouveaux développeurs, qui transmettent un contexte avant d’obtenir une erreur d’exécution.

Le principal problème de l’API d’exécution est qu’elle a été conçue pour les applications n’utilisant qu’une seule fenêtre. Au fil du temps, la demande de prise en charge de plusieurs fenêtres a augmenté, et l’API s’y prête mal.

### Réflexions sur l’API de la v3

Ne serait-il pas formidable de pouvoir procéder ainsi ?

```go
func main() {
    app := wails.NewApplication(options.App{})
    myWindow := app.NewWindow(options.Window{})
    myWindow.SetTitle("My Window")
    myWindow.On(events.Window.Close, func() {
        app.Quit()
    })
    app.Run()
}
```

Cette approche programmatique est beaucoup plus intuitive et permet au développeur d’interagir directement avec les éléments de l’application. Toutes les méthodes d’exécution actuelles relatives aux fenêtres deviendraient simplement des méthodes de l’objet fenêtre. Quant aux autres méthodes d’exécution, nous pourrions les déplacer vers l’objet application comme suit :

```go
app := wails.NewApplication(options.App{})
app.NewInfoDialog(options.InfoDialog{})
app.Log.Info("Hello World")
```

Cette API beaucoup plus puissante permettra de créer des applications plus complexes. Elle permettra également de créer plusieurs fenêtres, ce qui constitue [la fonctionnalité ayant recueilli le plus de votes sur GitHub](https://github.com/wailsapp/wails/issues/1480) :

```go
func main() {
    app := wails.NewApplication(options.App{})
    myWindow := app.NewWindow(options.Window{})
    myWindow.SetTitle("My Window")
    myWindow.On(events.Window.Close, func() {
        app.Quit()
    })
    myWindow2 := app.NewWindow(options.Window{})
    myWindow2.SetTitle("My Window 2")
    myWindow2.On(events.Window.Close, func() {
        app.Quit()
    })
    app.Run()
}
```

### Génération des liaisons

L’une des fonctionnalités essentielles de Wails consiste à générer des liaisons pour vos méthodes Go afin qu’elles puissent être appelées depuis JavaScript. La méthode actuelle relève quelque peu du bricolage. Elle consiste à compiler l’application avec un indicateur spécial, puis à exécuter le binaire obtenu, qui utilise la réflexion pour déterminer les éléments liés. Cela crée une sorte de problème de la poule et de l’œuf : vous ne pouvez pas compiler l’application sans les liaisons, ni générer les liaisons sans compiler l’application. Il existe de nombreuses solutions de contournement, mais la meilleure serait de renoncer entièrement à cette approche.

Plusieurs tentatives ont été menées pour écrire un analyseur statique destiné aux projets Wails, mais elles n’ont pas beaucoup progressé. Plus récemment, l’abondance accrue de ressources sur le sujet a légèrement facilité cette tâche.

Par rapport à la réflexion, l’approche fondée sur l’AST est beaucoup plus rapide, mais aussi nettement plus complexe. Dans un premier temps, nous devrons peut-être imposer certaines contraintes sur la manière de définir les liaisons dans le code. L’objectif est de prendre en charge les cas d’utilisation les plus courants, puis d’élargir cette prise en charge ultérieurement.

### Le système de build

À l’instar de l’approche déclarative de l’API, le système de build a été créé pour masquer la complexité de la création d’une application de bureau. Lorsque vous exécutez `wails build`, il effectue de nombreuses opérations en arrière-plan :

- Compile le binaire du backend nécessaire aux liaisons et génère celles-ci
- Installe les dépendances du frontend
- Génère les ressources du frontend
- Détermine si l’icône de l’application est présente et, le cas échéant, l’incorpore
- Compile le binaire final
- Si le build cible `darwin/universal`, compile 2 binaires, l’un pour `darwin/amd64` et l’autre pour `darwin/arm64`, puis crée un binaire universel à l’aide de `lipo`
- Si une compression est requise, compresse le binaire avec UPX
- Détermine si ce binaire doit être empaqueté et, le cas échéant :
  - Vérifie que l’icône et le manifeste de l’application sont compilés dans le binaire (Windows)
  - Crée le paquet applicatif, génère le jeu d’icônes, puis copie celui-ci, le binaire et Info.plist dans le paquet applicatif (Mac)

- Si un programme d’installation NSIS est requis, le génère

L’ensemble de ce processus, bien que très puissant, est également très opaque. Il est très difficile à personnaliser et à déboguer.

Pour remédier à cela dans la v3, je souhaiterais adopter un système de build externe à Wails. Après avoir utilisé [Task](https://taskfile.dev/) pendant quelque temps, j’en suis devenu un fervent adepte. C’est un excellent outil pour configurer des systèmes de build, qui devrait être assez familier à quiconque a déjà utilisé des Makefiles.

Le système de build serait configuré à l’aide d’un fichier `Taskfile.yml`, généré par défaut avec chacun des modèles pris en charge. Il contiendrait toutes les étapes nécessaires à l’exécution des tâches actuelles, telles que la compilation ou l’empaquetage de l’application, ce qui faciliterait sa personnalisation.

Cet outillage ne nécessitera aucune dépendance externe, puisqu’il fera partie de la CLI de Wails. Vous pourrez donc continuer à utiliser `wails build`, qui effectuera toujours tout ce qu’il fait actuellement. Toutefois, si vous souhaitez personnaliser le processus de build, vous pourrez le faire en modifiant le fichier `Taskfile.yml`. Vous pourrez également comprendre facilement les étapes du build et utiliser votre propre système de build si vous le souhaitez.

La pièce manquante du système de build réside dans les opérations élémentaires du processus, telles que la génération des icônes, la compression et l’empaquetage. Imposer de nombreux outils externes ne constituerait pas une bonne expérience pour les développeurs. Pour y remédier, la CLI de Wails intégrera toutes ces fonctionnalités. Les builds continueront ainsi de fonctionner comme prévu sans outil externe supplémentaire, mais vous pourrez remplacer n’importe quelle étape par l’outil de votre choix.

Ce système de build sera beaucoup plus transparent, facilitera la personnalisation et résoudra de nombreux problèmes signalés à son sujet.

## Les bénéfices

Ces changements positifs apporteront d’importants avantages au projet :

- La nouvelle API sera beaucoup plus intuitive et permettra de créer des applications plus complexes.
- Le recours à l’analyse statique pour générer les liaisons sera beaucoup plus rapide et réduira considérablement la complexité du processus actuel.
- L’utilisation d’un système de build externe et éprouvé rendra le processus entièrement transparent et permettra une personnalisation poussée.

Les responsables de la maintenance du projet bénéficieront des avantages suivants :

- La nouvelle API sera beaucoup plus facile à maintenir et à adapter aux nouvelles fonctionnalités et plateformes.
- Le nouveau système de build sera beaucoup plus facile à maintenir et à étendre. J’espère que cela favorisera l’émergence d’un nouvel écosystème de pipelines de build créés par la communauté.
- Une meilleure séparation des responsabilités au sein du projet facilitera l’ajout de nouvelles fonctionnalités et plateformes.

## Le plan

Une grande partie des expérimentations nécessaires a déjà été réalisée et les résultats sont prometteurs. Aucun calendrier n’est actuellement défini pour ces travaux, mais j’espère qu’une version alpha pour Mac sera disponible d’ici la fin du premier trimestre 2023, afin que la communauté puisse la tester, l’expérimenter et faire part de ses retours.

## Résumé

- L’API de la v2 est déclarative, masque beaucoup de choses aux développeurs et ne convient pas à des fonctionnalités telles que la gestion de plusieurs fenêtres. Une nouvelle API, plus simple, plus intuitive et plus puissante, sera créée.
- Le système de build est opaque et difficile à personnaliser. Nous adopterons donc un système externe qui rendra l’ensemble du processus transparent.
- La génération des liaisons est lente et complexe. Nous adopterons donc l’analyse statique, qui éliminera une grande partie de la complexité de la méthode actuelle.

Un travail considérable a été consacré aux mécanismes internes de la v2, qui sont solides. Il est maintenant temps de nous attaquer à la couche qui les recouvre afin d’améliorer nettement l’expérience des développeurs.

J’espère que cette perspective vous enthousiasme autant que moi. J’ai hâte de connaître votre avis et de recevoir vos retours.

Cordialement,

&dash; Lea

P.-S. : Si vous-même ou votre entreprise trouvez Wails utile, pensez à [soutenir financièrement le projet](https://github.com/sponsors/leaanthony). Merci !

P.-P.-S. : Oui, il s’agit bien d’une véritable capture d’écran d’une application multifenêtre créée avec Wails. Ce n’est pas une maquette. Elle est réelle. Elle est géniale. Elle arrive bientôt.
