---
title: "Questions fréquentes"
description: "Réponses aux questions les plus fréquentes sur la création d’applications avec Wails v3"
slug: "faq"
sourcePath: "faq.md"
---

## Généralités

### Qu’est-ce que Wails ?

Wails est un framework permettant de créer des applications de bureau avec Go et les technologies web. La logique de votre application est écrite en Go, son interface est construite avec HTML, CSS et JavaScript (ou tout autre framework frontend), puis Wails l’affiche dans la vue web native du système d’exploitation. Vous obtenez ainsi une application compacte, rapide et d’apparence native : aucun navigateur intégré, une faible consommation de mémoire et un seul fichier binaire pesant généralement environ 10 Mo.

### Quelles plateformes Wails prend-il en charge ?

| Plateforme | Configuration requise |
| --- | --- |
| Windows | AMD64 et ARM64. Utilise le runtime [WebView2](https://developer.microsoft.com/microsoft-edge/webview2/). |
| macOS | 10.15 ou version ultérieure sur Intel (les applications peuvent cibler 10.13 ou version ultérieure), 11.0 ou version ultérieure sur Apple Silicon. Les binaires universels sont pris en charge. |
| Linux | AMD64 et ARM64. La pile par défaut repose sur GTK4 avec WebKitGTK 6.0 (Ubuntu 24.04 ou version ultérieure, Debian 13 ou version ultérieure, Fedora 40 ou version ultérieure et distributions similaires). Les distributions qui ne fournissent que WebKit2GTK 4.1, telles qu’Ubuntu 22.04, Debian 12 et RHEL 9, sont prises en charge au moyen de la build héritée `-tags gtk3` (disponible jusqu’à la version v3.1). Les distributions qui ne disposent que de WebKit2GTK 4.0 ne sont pas prises en charge. Consultez le [guide de build pour Linux](/guides/build/linux/). |
| iOS et Android | Expérimental. Consultez les [guides pour appareils mobiles](/guides/mobile/). |

Vous pouvez également servir votre application comme une application web classique à l’aide de la [build serveur](/guides/server-build/).

Exécutez `wails3 doctor` à tout moment pour vérifier votre système et obtenir les instructions d’installation propres à votre plateforme.

### De quoi ai-je besoin pour commencer ?

- Go 1.25 ou version ultérieure
- Node.js et npm (pour la build du frontend)
- Chaîne d’outils de la plateforme : WebView2 sous Windows (préinstallé sur 10/11), les outils en ligne de commande Xcode sous macOS, ainsi que `gcc` et les paquets de développement GTK/WebKit sous Linux

`wails3 doctor` vérifie tous ces éléments pour vous et vous indique précisément ce qui manque. Consultez la section [Installation](/quick-start/installation/) pour obtenir la procédure complète.

### Wails v3 est-il prêt pour la production ?

Wails v3 est un logiciel bêta doté d’une API de bureau stable. Des applications l’utilisent en production, mais vous devriez effectuer des tests approfondis avant tout déploiement pendant que nous apportons les dernières finitions à 3.0. Consultez la [page d’état du projet](/status/) pour connaître la situation actuelle. Wails v2 est la version stable actuelle et continue de recevoir des correctifs.

## Développement

### Dois-je connaître Go ?

Des connaissances de base en Go sont utiles, mais vous n’avez pas besoin d’être un expert. La logique de l’application réside dans de simples méthodes Go, et les [tutoriels](/tutorials/overview/) vous guident pour tout le reste. De nombreux développeurs apprennent Go en créant leur première application Wails.

### Puis-je utiliser mon framework frontend préféré ?

Oui. S’il produit du HTML, du CSS et du JavaScript, il fonctionne avec Wails. Des modèles sont fournis pour React, Vue, Svelte et JavaScript sans framework (chacun avec des variantes TypeScript), et tout autre framework peut être intégré en quelques minutes. Consultez la section [Frameworks frontend](/guides/dev/frontend-frameworks/).

### Comment appeler des fonctions Go depuis JavaScript ?

Enregistrez un service ; Wails génère alors des liaisons typées pour celui-ci :

```go
// Go
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello " + name
}
```

```javascript
// JavaScript
import { GreetService } from "./bindings/changeme";

const message = await GreetService.Greet("World");
```

Les liaisons sont régénérées automatiquement pendant `wails3 dev`, ou à la demande avec `wails3 generate bindings`. Consultez la section [Services](/features/bindings/services/).

### Puis-je utiliser TypeScript ?

Oui. Le générateur de liaisons produit des définitions TypeScript pour vos services et leurs types, de sorte que les appels à Go sont entièrement typés.

### Comment envoyer des événements entre Go et JavaScript ?

```go
// Go
app.Event.Emit("time", time.Now().Format(time.RFC1123))
```

```javascript
// JavaScript
import { Events } from "@wailsio/runtime";

Events.On("time", (event) => {
    console.log(event.data);
});
```

Les noms des événements doivent correspondre exactement. Consultez la [référence des événements](/guides/events-reference/).

### Comment déboguer mon application ?

Exécutez `wails3 dev`, puis faites un clic droit dans la fenêtre pour ouvrir les outils de développement du navigateur, exactement comme vous le feriez sur le web. Le serveur de développement prend également en charge le rechargement à chaud de votre frontend. Consultez la section [Débogage](/guides/dev/debugging/).

## Build et distribution

### Comment créer une build de production ?

```bash
wails3 build
```

Votre fichier binaire est placé dans `bin/`. Les builds de production appliquent déjà des valeurs par défaut adaptées (tags de build, `-trimpath`, suppression des symboles) ; aucun indicateur supplémentaire n’est donc nécessaire pour obtenir un fichier binaire compact.

### Puis-je effectuer une compilation croisée ?

Dans certaines limites. La compilation croisée en Go pur ne s’applique pas, car chaque plateforme utilise des bibliothèques de vue web natives, mais les cas courants sont bien pris en charge :

```bash
# Different architecture, same OS
wails3 build GOOS=windows GOARCH=arm64

# macOS universal binary
wails3 task darwin:build:universal
```

La création d’une build Linux depuis un autre système d’exploitation utilise une chaîne d’outils basée sur Docker. Consultez la section [Builds multiplateformes](/guides/build/cross-platform/) pour connaître la matrice complète.

### Comment créer un programme d’installation ou un paquet ?

```bash
wails3 package
```

Cette opération produit le format natif de la plateforme, et le [guide des programmes d’installation](/guides/installers/) couvre NSIS sous Windows, les bundles `.app` et les DMG sous macOS, ainsi que les paquets Linux.

### Comment signer le code de mon application ?

La signature sous Windows et macOS, y compris la notarisation, est expliquée étape par étape dans le [guide de signature](/guides/build/signing/).

## Fonctionnalités

### Puis-je créer plusieurs fenêtres ?

Oui, la prise en charge de plusieurs fenêtres est native dans la v3 :

```go
window1 := app.Window.New()
window2 := app.Window.New()
```

Consultez la section [Fenêtres multiples](/features/windows/multiple/).

### Wails prend-il en charge la zone de notification système ?

Oui, y compris les menus et les gestionnaires de clics :

```go
systemTray := app.SystemTray.New()
systemTray.SetIcon(iconBytes)
systemTray.SetMenu(myMenu)
```

Consultez la section [Zone de notification système](/features/menus/systray/).

### Puis-je utiliser des boîtes de dialogue natives ?

Oui. Les boîtes de dialogue de fichiers, de messages et de questions utilisent toutes les implémentations natives :

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select File").
    PromptForSingleSelection()
```

Consultez la section [Boîtes de dialogue](/features/dialogs/overview/).

### Wails prend-il en charge les mises à jour automatiques ?

Oui. Wails v3 comprend un mécanisme intégré de mise à jour automatique (`app.Updater`), avec des fournisseurs interchangeables pour GitHub Releases, keygen.sh et Sparkle AppCast, une vérification cryptographique des signatures et une interface utilisateur par défaut que vous pouvez personnaliser ou remplacer. Consultez le guide [Outil de mise à jour intégré à l’application](/guides/updater/) et le tutoriel [Application Wails à mise à jour automatique](/tutorials/04-self-update-a-wails-app/).

## Dépannage

### Quelque chose ne fonctionne pas. Par où commencer ?

```bash
wails3 doctor
```

Il vérifie votre chaîne d’outils, répertorie les dépendances manquantes avec les commandes permettant de les installer et affiche les informations de version à inclure dans tout rapport de bogue.

### La compilation échoue

Essayez les solutions habituelles dans cet ordre :

1. `go mod tidy`
2. `cd frontend && npm install` (l’absence de `node_modules` est la cause la plus fréquente)
3. Mettez à jour la CLI : `go install github.com/wailsapp/wails/v3/cmd/wails3@latest`
4. Sous Linux, recherchez dans `wails3 doctor` les paquets GTK/WebKit manquants

### Mes liaisons sont absentes ou obsolètes

```bash
wails3 generate bindings
```

Les liaisons sont régénérées automatiquement en mode développement. Si vous avez ajouté un service ou modifié des signatures de méthodes en dehors de `wails3 dev`, régénérez-les manuellement.

### Les événements ne se déclenchent pas

Les noms d’événements doivent correspondre exactement entre `app.Event.Emit("name", ...)` dans Go et `Events.On("name", ...)` dans JavaScript. Commencez par rechercher les fautes de frappe et les différences de casse.

### J’ai trouvé un bogue

Veuillez [ouvrir un ticket](https://github.com/wailsapp/wails/issues) et y inclure la sortie de `wails3 doctor`. Le [guide sur les retours](/feedback/) explique comment rendre un rapport facile à traiter.

## Migration depuis la v2

### Dois-je migrer de la v2 vers la v3 ?

La v3 apporte la prise en charge de plusieurs fenêtres, une API plus claire fondée sur des services, un outil de mise à jour intégré, un système de compilation bien plus souple et de meilleures performances. Les nouveaux projets devraient démarrer avec la v3. Pour les projets existants, le [Guide de migration](/migration/v2-to-v3/) présente les différences.

### La v2 continuera-t-elle d’être maintenue ?

Oui. La v2 continue de recevoir des correctifs pendant que la v3 progresse vers sa version stable.

### Puis-je utiliser la v2 et la v3 côte à côte ?

Oui. Les CLI sont des exécutables distincts (`wails` et `wails3`) et les modules utilisent des chemins d’importation différents. Des projets reposant sur des versions majeures différentes peuvent donc coexister sans problème sur une même machine.

## Communauté

### Comment obtenir de l’aide ?

- [Discord](https://discord.gg/JDdSxwjhGf) pour les questions rapides et les échanges
- [GitHub Discussions](https://github.com/wailsapp/wails/discussions) pour les questions nécessitant des échanges plus approfondis
- [GitHub Issues](https://github.com/wailsapp/wails/issues) pour les bogues

### Comment contribuer ?

Consultez le [Guide de contribution](/contributing/). Les correctifs de bogues sont toujours les bienvenus. Les nouvelles fonctionnalités et les modifications du comportement public passent par une PR au statut brouillon proposant une [WEP (proposition d’amélioration de Wails)](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md). Une discussion informelle sur Discord ou GitHub Discussions est facultative.

### Où puis-je trouver des exemples ?

Le dépôt contient plus de 60 exemples exécutables couvrant les fenêtres, les boîtes de dialogue, les événements, la zone de notification système, les services et bien plus encore : [v3/examples](https://github.com/wailsapp/wails/tree/master/v3/examples).

## Vous avez encore des questions ?

Posez-les sur [Discord](https://discord.gg/JDdSxwjhGf) ou [ouvrez une discussion](https://github.com/wailsapp/wails/discussions).
