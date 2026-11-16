---
title: "Contrôle par LLM (MCP)"
description: "Permettez aux agents LLM de tester et de contrôler votre application via le Model Context Protocol"
slug: "guides/mcp-service"
sourcePath: "guides/mcp-service.md"
---

@note{type="caution" title="Fonctionnalité expérimentale"}
Le serveur MCP intégré est expérimental et son API pourra évoluer dans de prochaines versions.

@end

Wails v3 intègre un serveur [Model Context Protocol](https://modelcontextprotocol.io) (MCP) qui permet aux agents LLM — Claude Code, assistants d’IDE ou tout client MCP — d’inspecter, de tester et de piloter une application Wails **en cours d’exécution**.

Pour automatiser le cycle de vie des projets, utilisez le serveur CLI `wails3 mcp` distinct. Il permet aux agents d’inspecter et d’initialiser des projets, d’exécuter des diagnostics, de lancer des builds et des tâches de développement, de générer des bindings, d’exécuter des tâches Taskfile nommées et de récupérer une sortie de tâche limitée. Par défaut, le serveur CLI est confiné au répertoire courant et ne permet pas d’exécuter des commandes shell arbitraires. Consultez la [documentation du MCP de la CLI](/guides/cli/#mcp) pour en savoir plus sur le transport, l’authentification et les outils.

Lorsqu’il est activé, un agent connecté à votre application peut :

- **Répertorier et contrôler les fenêtres** — taille, position, focus, plein écran, outils de développement, rechargement, etc.
- **Inspecter le DOM** — rechercher des éléments, obtenir le HTML, prendre un instantané structurel
- **Évaluer du JavaScript** — exécuter du code arbitraire dans n’importe quelle fenêtre et en obtenir le résultat
- **Simuler les actions de l’utilisateur** — mouvements de souris, clics, glisser-déposer et défilements affichés avec un **curseur animé à l’écran** pour vous permettre de suivre le travail de l’agent
- **Saisir du texte et appuyer sur des touches** — événements réalistes générés caractère par caractère et compatibles avec les champs contrôlés de React
- **Appeler des méthodes Go liées** et émettre ou attendre des événements de l’application

## Fonctionnement

Le serveur MCP n’est compilé dans votre application que si le tag de build **`mcp`** est présent. Sans ce tag, le code du serveur est totalement absent du binaire : aucune surcharge à l’exécution, aucun port ouvert et aucune surface d’attaque.

Lorsque le tag est présent, le serveur démarre automatiquement dans `App.Run()`, écoute par défaut sur `127.0.0.1:9099` et consigne son point de terminaison. Aucun code utilisateur n’est nécessaire.

## Tutoriel

### Étape 1 — écrivez une application Wails normale

MCP ne nécessite ni importation ni enregistrement. Créez votre application comme vous le feriez normalement :

```go {title="main.go"}
package main

import (
    "embed"
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name: "My App",
        Assets: application.AssetOptions{
            Handler: application.BundledAssetFileServer(assets),
        },
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: "My App",
        Width: 1024, Height: 768,
    })

    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

### Étape 2 — compilez ou exécutez avec le tag `mcp`

@tabs
[CLI Wails (recommandé)]
Définissez `WAILS_MCP=1` pour que la CLI Wails ajoute automatiquement le tag `mcp` :

```shell
# Development
WAILS_MCP=1 wails3 dev

# Production build
WAILS_MCP=1 wails3 build
```

[Go directement]
Transmettez directement le tag à `go run` ou à `go build` :

```shell
go run -tags mcp .
go build -tags mcp -o myapp .
```

[Windows (PowerShell)]
```powershell
$env:WAILS_MCP = "1"
wails3 dev
# or
wails3 build
```

@end

Au démarrage, l’application consigne le point de terminaison MCP :

```
INFO MCP server started. Connect MCP clients using the streamable HTTP transport.
     url=http://127.0.0.1:9099/mcp
```

### Étape 3 — connectez un client

Le serveur utilise le **transport HTTP streamable de MCP**. Connectez-vous avec n’importe quel client compatible avec MCP.

@tabs
[Claude Code]
```shell
claude mcp add --transport http my-app http://127.0.0.1:9099/mcp
```

Demandez ensuite à Claude d’interagir avec votre application :

```
Click the "Submit" button, then verify a success toast appears.
```

[VS Code (GitHub Copilot)]
Ajoutez ceci à `.vscode/settings.json` :

```json
{
  "github.copilot.chat.mcp.enabled": true,
  "mcp": {
    "servers": {
      "my-wails-app": {
        "type": "http",
        "url": "http://127.0.0.1:9099/mcp"
      }
    }
  }
}
```

[Autres clients]
Configurez tout client MCP prenant en charge le transport HTTP streamable pour qu’il utilise l’adresse suivante :

```
http://127.0.0.1:9099/mcp
```

@end

### Étape 4 — exécutez une session de test

Demandez à l’agent de tester votre application. Voici quelques exemples de prompts :

```
Take a DOM snapshot of the main window.
```

```
Click the "Add item" button, type "Hello world" in the input field,
then press Enter and verify the item appears in the list.
```

```
Call the bound method main.GreetService.Greet with argument ["World"]
and return the result.
```

```
Wait for the event "save:complete" while clicking the Save button.
```

## Configuration

Toute la configuration s’effectue au moyen de variables d’environnement : aucune modification du code n’est nécessaire.

| Variable d’environnement | Valeur par défaut | Description |
| --- | --- | --- |
| `WAILS_MCP` | (non définie) | Définissez cette variable sur `1`, `true`, `on` ou `yes` pour ajouter automatiquement le tag de build `mcp` lors de l’utilisation de la CLI Wails. |
| `WAILS_MCP_HOST` | `127.0.0.1` | Interface d’écoute. Une écoute sur une interface autre que l’interface de bouclage nécessite `WAILS_MCP_TOKEN`. |
| `WAILS_MCP_TOKEN` | non définie | Jeton Bearer facultatif sur l’interface de bouclage, mais obligatoire sur les autres adresses d’écoute. Les clients envoient `Authorization: Bearer <token>`. |
| `WAILS_MCP_PORT` | `9099` | Port d’écoute. Définissez-le sur `0` pour utiliser un port libre attribué aléatoirement et affiché dans le journal. |
| `WAILS_MCP_TIMEOUT` | `30000` | Délai d’expiration par défaut de l’évaluation JavaScript, en **millisecondes**. |
| `WAILS_MCP_HIDE_CURSOR` | (non définie) | Définissez cette variable sur `1` ou `true` pour désactiver la superposition du curseur animé. |

Exemple — port personnalisé et délai d’expiration de 60 secondes :

```shell
WAILS_MCP=1 WAILS_MCP_PORT=9200 WAILS_MCP_TIMEOUT=60000 wails3 dev
```

## Outils disponibles

| Outil | Fonction |
| --- | --- |
| `app_info` | Informations sur l’application : plateforme, architecture, toutes les fenêtres et point de terminaison MCP |
| `windows_list` | Répertorie toutes les fenêtres avec leur géométrie et leur état |
| `window_control` | Met au premier plan, redimensionne, déplace, passe en plein écran, ouvre les outils de développement, recharge, définit l’URL, etc. (22 actions) |
| `js_eval` | Évalue du JavaScript dans une fenêtre (corps asynchrone, `return` pour la valeur) |
| `dom_html` | Obtient le HTML de la page ou d’un élément précis |
| `dom_query` | Recherche des éléments à l’aide d’un sélecteur CSS : balise, texte, limites et visibilité |
| `screenshot_dom` | Instantané structurel de la page visible (fondé sur le DOM, sans pixels) |
| `mouse_move` | Anime le curseur jusqu’à un point ou un sélecteur CSS |
| `mouse_click` | Clique avec le curseur animé (bouton gauche, droit ou central, double-clic, touches de modification) |
| `mouse_drag` | Effectue un glisser-déposer avec le curseur animé (prend en charge les éléments HTML5 glissables) |
| `mouse_scroll` | Fait défiler la page à un point ou sur un élément |
| `keyboard_type` | Saisit du texte caractère par caractère avec des événements réalistes |
| `keyboard_press` | Appuyer sur une seule touche (Entrée, Tabulation, Échap, Flèche vers le bas, …), avec des touches de modification facultatives |
| `call_bound_method` | Appeler une méthode d’un service Go lié, par exemple `main.GreetService.Greet` |
| `emit_event` | Émettre un événement d’application Wails |
| `wait_for_event` | Attendre un événement d’application Wails et renvoyer ses données |

### Prise en charge de plusieurs fenêtres

Tous les outils qui agissent sur une fenêtre acceptent un argument `window` facultatif contenant le **nom** de la fenêtre (défini via `WebviewWindowOptions.Name`). S’il est omis, l’outil cible la fenêtre actuellement active ou, si aucune ne l’est, la première fenêtre.

```
List all windows, then click the "New" button in the window named "editor".
```

### Sélection des éléments

Les outils de souris et de clavier acceptent au choix :

- **Sélecteur CSS** — `selector: "#submit-btn"` (l’élément défile automatiquement jusqu’à devenir visible)
- **Coordonnées** — `x: 400, y: 300` (pixels CSS relatifs à la zone d’affichage)

Pour les opérations de glisser-déposer, utilisez les préfixes `from_` et `to_` :

```
Drag from selector: ".card" to selector: ".dropzone"
```

## Sécurité

@note{type="caution"}
Le serveur MCP permet de contrôler entièrement votre application par programmation. Toute personne ayant accès à ses outils peut lire le DOM, évaluer du JavaScript, cliquer sur des boutons et appeler des méthodes Go.

@end

- Par défaut, le serveur écoute sur `127.0.0.1`. Les origines de navigateur doivent être des origines de bouclage HTTP(S) ; les origines opaques (`null`), mal formées et externes sont rejetées.
- Les clients MCP natifs sans en-têtes restent pris en charge. Sans `WAILS_MCP_TOKEN`, les processus locaux et les origines de navigateur locales autorisées sont approuvés ; les contrôles d’origine ne constituent pas une authentification. Définissez un jeton à forte entropie pour exiger une authentification par jeton porteur pour tous les appels `/mcp`. Configurez le même jeton dans l’en-tête `Authorization: Bearer <token>` du client. Les requêtes de contrôle préalable ne nécessitent pas ce jeton.
- Le rappel `/eval-result` utilise des identifiants imprévisibles propres à chaque évaluation plutôt que le jeton porteur du client, ce qui préserve la compatibilité de la transmission des résultats de la vue web.
- Les builds de production ne doivent **pas** inclure le tag `mcp`. La CLI Wails ne l’ajoute que lorsque `WAILS_MCP=1` est explicitement défini, et le `wails3 build` par défaut ne contient aucun code du serveur.
- Si vous devez exposer le serveur sur une interface autre que l’interface de bouclage (par exemple pour effectuer des tests sur un réseau local), définissez `WAILS_MCP_HOST=0.0.0.0` ainsi qu’un `WAILS_MCP_TOKEN` à forte entropie ; le démarrage échoue en l’absence de jeton. Sur les réseaux non fiables, utilisez un tunnel chiffré ou un proxy assurant la terminaison TLS, car le service d’écoute intégré utilise HTTP.

## Application d’exemple

Une application complète de démonstration de tous les outils est disponible à l’emplacement [`v3/examples/mcp`](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/mcp). Elle comprend :

- Compteur avec boutons d’incrémentation et de réinitialisation
- Champ de saisie du nom avec les méthodes liées Greet, Add et Shout
- Source et cible de glisser-déposer HTML5
- Liste déroulante (50 éléments)
- Journal des événements

Exécutez-la avec :

```shell
cd v3/examples/mcp
go run -tags mcp .
```

Connectez ensuite Claude Code ou tout autre client MCP à `http://127.0.0.1:9099/mcp` et demandez-lui d’interagir avec l’interface utilisateur.

## Vos retours

Le serveur MCP intégré est expérimental, et vos retours détermineront son évolution. Si vous l’essayez, nous aimerions savoir quels client et outils vous avez utilisés, ce à quoi vous vous attendiez et ce qui s’est réellement produit, ainsi que l’utilité de laisser un agent piloter votre application ; les rapports les plus utiles indiquent précisément ce que vous avez exécuté. Faites-nous-en part dans la [discussion sur les retours concernant le serveur MCP](https://github.com/wailsapp/wails/discussions/5692).
