---
title: "Foire aux questions"
description: "Réponses aux questions les plus fréquentes sur le développement d’applications avec Wails v3"
slug: "faq"
sourcePath: "faq.md"
---

## Généralités

### Qu’est-ce que Wails ?

Wails est un framework permettant de créer des applications de bureau avec Go et les technologies web. La logique de votre application est écrite en Go, son interface est créée avec HTML, CSS et JavaScript (ou n’importe quel framework frontend), puis Wails l’affiche dans la vue web native du système d’exploitation. Vous obtenez ainsi une application légère, rapide et d’apparence native : aucun navigateur intégré, une faible consommation de mémoire et un seul fichier binaire pesant généralement environ 10 Mo.

### Quelles plateformes Wails prend-il en charge ?

| Plateforme | Configuration requise |
| --- | --- |
| Windows | AMD64 et ARM64. Utilise le runtime [WebView2](https://developer.microsoft.com/microsoft-edge/webview2/). |
| macOS | 10.15 ou version ultérieure sur Intel (les applications peuvent cibler 10.13 ou une version ultérieure), 11.0 ou version ultérieure sur Apple Silicon. Les binaires universels sont pris en charge. |
| Linux | AMD64 et ARM64. La pile par défaut utilise GTK4 avec WebKitGTK 6.0 (Ubuntu 24.04 ou version ultérieure, Debian 13 ou version ultérieure, Fedora 40 ou version ultérieure et distributions similaires). Les distributions qui ne fournissent que WebKit2GTK 4.1, comme Ubuntu 22.04, Debian 12 et RHEL 9, sont prises en charge au moyen de l’ancien build `-tags gtk3` (disponible jusqu’à v3.1). Les distributions qui ne disposent que de WebKit2GTK 4.0 ne sont pas prises en charge. Consultez le [guide de compilation pour Linux](/guides/build/linux/). |
| iOS et Android | Expérimental. Consultez les [guides pour plateformes mobiles](/guides/mobile/). |

Vous pouvez également servir votre application comme une application web classique à l’aide du [build serveur](/guides/server-build/).

Exécutez `wails3 doctor` à tout moment pour vérifier votre système et obtenir les instructions d’installation propres à votre plateforme.

### De quoi ai-je besoin pour commencer ?

- Go 1.26 ou version ultérieure
- Node.js et npm (pour le build frontend)
- Chaîne d’outils de la plateforme : WebView2 sous Windows (préinstallé sur 10/11), les outils en ligne de commande Xcode sous macOS, `gcc` ainsi que les paquets de développement GTK/WebKit sous Linux

`wails3 doctor` vérifie tous ces éléments et vous indique précisément ce qui manque. Consultez la section [Installation](/quick-start/installation/) pour suivre la procédure complète.

### Wails v3 est-il prêt pour la production ?

Wails v3.0.0 est la version stable pour les applications de bureau sous Windows, macOS et Linux. La prise en charge d’Android et d’iOS reste expérimentale et ne relève pas du contrat de compatibilité stable pour les applications de bureau. Consultez la [page d’état du projet](/status/).

## Développement

### Dois-je connaître Go ?

Des connaissances de base en Go sont utiles, mais il n’est pas nécessaire d’être expert. La logique de l’application réside dans de simples méthodes Go, et les [tutoriels](/tutorials/overview/) vous guident pour tout le reste. De nombreux développeurs apprennent Go en créant leur première application Wails.

### Puis-je utiliser mon framework frontend préféré ?

Oui. S’il génère du HTML, du CSS et du JavaScript, il fonctionne avec Wails. Des modèles sont fournis pour React, Vue, Svelte et JavaScript sans framework, chacun avec une variante TypeScript. Tout autre framework peut être intégré en quelques minutes. Consultez la section [Frameworks frontend](/guides/dev/frontend-frameworks/).

### Comment appeler des fonctions Go depuis JavaScript ?

Enregistrez un service ; Wails générera des liaisons typées pour celui-ci :

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

Oui. Le générateur de liaisons produit les définitions TypeScript de vos services et de leurs types, de sorte que les appels vers Go sont entièrement typés.

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

Exécutez `wails3 dev`, puis faites un clic droit dans la fenêtre pour ouvrir les outils de développement du navigateur, exactement comme sur le Web. Le serveur de développement permet également le rechargement à chaud de votre frontend. Consultez la section [Débogage](/guides/dev/debugging/).

## Compilation et distribution

### Comment créer un build de production ?

```bash
wails3 build
```

Votre fichier binaire est placé dans `bin/`. Les builds de production appliquent déjà des valeurs par défaut raisonnables (tags de build, `-trimpath`, symboles supprimés) ; aucun indicateur supplémentaire n’est donc nécessaire pour obtenir un fichier binaire léger.

### Puis-je effectuer une compilation croisée ?

Dans certaines limites. La compilation croisée en Go pur ne s’applique pas, car chaque plateforme utilise des bibliothèques de vue web natives. Les cas courants sont toutefois bien pris en charge :

```bash
# Different architecture, same OS
wails3 build GOOS=windows GOARCH=arm64

# macOS universal binary
wails3 task darwin:build:universal
```

La compilation pour Linux depuis un autre système d’exploitation utilise une chaîne d’outils basée sur Docker. Consultez la section [Builds multiplateformes](/guides/build/cross-platform/) pour connaître la matrice complète.

### Comment créer un programme d’installation ou un paquet ?

```bash
wails3 package
```

Cette opération produit le format natif de la plateforme. Le [guide des programmes d’installation](/guides/installers/) traite de NSIS sous Windows, des bundles `.app` et des DMG sous macOS, ainsi que des paquets Linux.

### Comment signer le code de mon application ?

La signature sous Windows et macOS, y compris la notarisation, est expliquée étape par étape dans le [guide de signature](/guides/build/signing/).

## Fonctionnalités

### Puis-je créer plusieurs fenêtres ?

Oui, la gestion de plusieurs fenêtres est native dans v3 :

```go
window1 := app.Window.New()
window2 := app.Window.New()
```

Consultez la section [Fenêtres multiples](/features/windows/multiple/).

### Wails prend-il en charge la zone de notification ?

Oui, y compris les menus et les gestionnaires de clics :

```go
systemTray := app.SystemTray.New()
systemTray.SetIcon(iconBytes)
systemTray.SetMenu(myMenu)
```

Consultez la section [Zone de notification](/features/menus/systray/).

### Puis-je utiliser des boîtes de dialogue natives ?

Oui. Les boîtes de dialogue de sélection de fichiers, de message et de question utilisent toutes les implémentations natives :

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select File").
    PromptForSingleSelection()
```

Consultez la section [Boîtes de dialogue](/features/dialogs/overview/).

### Wails prend-il en charge les mises à jour automatiques ?

Oui. Wails v3 comprend un mécanisme intégré de mise à jour automatique (`app.Updater`), avec des fournisseurs interchangeables pour GitHub Releases, keygen.sh et Sparkle AppCast, la vérification cryptographique des signatures et une interface utilisateur par défaut que vous pouvez personnaliser ou remplacer. Consultez le guide [Mise à jour intégrée à l’application](/guides/updater/) et le tutoriel [Application Wails à mise à jour automatique](/tutorials/04-self-update-a-wails-app/).

## Dépannage

### Quelque chose ne fonctionne pas. Par où commencer ?

```bash
wails3 doctor
```

Cette commande vérifie votre chaîne d’outils, répertorie les dépendances manquantes avec les commandes permettant de les installer et affiche les informations de version à inclure dans tout rapport de bug.

### Mon build échoue

Essayez les solutions habituelles dans l’ordre suivant :

1. `go mod tidy`
2. `cd frontend && npm install` (l’absence de `node_modules` est la cause la plus fréquente)
3. Mettez à jour la CLI : `go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0`
4. Sous Linux, recherchez les paquets GTK/WebKit manquants avec `wails3 doctor`

### Mes bindings sont absents ou obsolètes

```bash
wails3 generate bindings
```

Les bindings sont régénérés automatiquement en mode développement. Si vous avez ajouté un nouveau service ou modifié des signatures de méthodes en dehors de `wails3 dev`, régénérez-les manuellement.

### Les événements ne se déclenchent pas

Les noms d’événements doivent correspondre exactement entre `app.Event.Emit("name", ...)` en Go et `Events.On("name", ...)` en JavaScript. Vérifiez d’abord les fautes de frappe et les différences de casse.

### J’ai trouvé un bug

Veuillez [ouvrir un ticket](https://github.com/wailsapp/wails/issues) et y inclure la sortie de `wails3 doctor`. Le [guide de retour d’expérience](/feedback/) explique comment rendre un rapport facile à traiter.

## Migration depuis v2

### Dois-je migrer de v2 vers v3 ?

v3 apporte la prise en charge de plusieurs fenêtres, une API basée sur des services plus claire, un système de mise à jour intégré, un système de build beaucoup plus flexible et de meilleures performances. Les nouveaux projets devraient démarrer avec v3. Pour les projets existants, le [guide de migration](/migration/v2-to-v3/) présente les différences pas à pas.

### v2 continuera-t-il à être maintenu ?

Oui. Wails v2 continue de recevoir les mises à jour critiques. Les nouveaux projets devraient utiliser Wails v3.0.0.

### Puis-je utiliser v2 et v3 côte à côte ?

Oui. Les interfaces en ligne de commande sont des exécutables distincts (`wails` et `wails3`), et les modules ont des chemins d’importation différents. Les projets utilisant des versions majeures différentes peuvent donc coexister sans problème sur une même machine.

## Communauté

### Comment obtenir de l’aide ?

- [Discord](https://discord.gg/JDdSxwjhGf) pour les questions rapides et les échanges
- [Discussions GitHub](https://github.com/wailsapp/wails/discussions) pour les questions qui nécessitent des réponses plus développées
- [Tickets GitHub](https://github.com/wailsapp/wails/issues) pour les bugs

### Comment contribuer ?

Consultez le [guide de contribution](/contributing/). Les corrections de bugs sont toujours les bienvenues. Les nouvelles fonctionnalités et les modifications du comportement public nécessitent une pull request contenant un brouillon de [WEP (Wails Enhancement Proposal)](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/README.md). Une discussion informelle sur Discord ou dans les Discussions GitHub est facultative.

### Où trouver des exemples ?

Le dépôt fournit plus de 60 exemples exécutables portant sur les fenêtres, les boîtes de dialogue, les événements, la zone de notification, les services et bien plus encore : [v3/examples](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples).

## Vous avez encore des questions ?

Posez-les sur [Discord](https://discord.gg/JDdSxwjhGf) ou [ouvrez une discussion](https://github.com/wailsapp/wails/discussions).
