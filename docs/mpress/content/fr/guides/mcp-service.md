---
title: "Contrôle par LLM (MCP)"
description: "Permettez aux agents LLM de tester et de contrôler votre application via le Model Context Protocol"
slug: "guides/mcp-service"
sourcePath: "guides/mcp-service.md"
---

@note{type="caution" title="Fonctionnalité expérimentale"}
Le serveur MCP intégré est expérimental et son API pourra évoluer dans de futures versions.

@end

Wails v3 intègre un serveur [Model Context Protocol](https://modelcontextprotocol.io) (MCP) qui permet aux agents LLM — Claude Code, assistants d’IDE ou tout client MCP — d’inspecter, de tester et de piloter une application Wails **en cours d’exécution**.

Pour automatiser le cycle de vie des projets, utilisez le serveur CLI `wails3 mcp` distinct. Il permet aux agents d’inspecter et d’initialiser des projets, d’exécuter des diagnostics, de lancer des compilations et des tâches de développement, de générer des liaisons, d’exécuter des tâches Taskfile nommées et de récupérer une quantité limitée de sortie des tâches. Par défaut, le serveur CLI est limité au répertoire courant et ne permet pas l’exécution arbitraire de commandes shell. Consultez la [documentation du MCP CLI](/guides/cli/#mcp) pour en savoir plus sur le transport, l’authentification et les outils.

Lorsqu’il est activé, un agent connecté à votre application peut :

- **Répertorier et contrôler les fenêtres** — taille, position, focus, plein écran, outils de développement, rechargement, etc.
- **Inspecter le DOM** — rechercher des éléments, obtenir le HTML, prendre un instantané structurel
- **Évaluer du JavaScript** — exécuter du code arbitraire dans n’importe quelle fenêtre et obtenir le résultat
- **Simuler les actions de l’utilisateur** — mouvements de souris, clics, glissements et défilements affichés avec un **curseur animé à l’écran** afin que vous puissiez observer le travail de l’agent
- **Saisir du texte et appuyer sur des touches** — événements réalistes générés caractère par caractère et compatibles avec les champs contrôlés de React
- **Appeler des méthodes Go liées** et émettre ou attendre des événements de l’application

## Fonctionnement

Le serveur MCP n’est compilé dans votre application que si l’étiquette de compilation **`mcp`** est présente. Sans cette étiquette, le code du serveur est totalement absent du fichier binaire : aucune surcharge à l’exécution, aucun port ouvert et aucune surface d’attaque.

Lorsque l’étiquette est présente, le serveur démarre automatiquement dans `App.Run()`, se lie par défaut à `127.0.0.1:9099` et consigne son point de terminaison. Aucun code utilisateur n’est requis.

## Tutoriel

### Étape 1 — écrivez une application Wails normale

MCP ne nécessite ni importation ni enregistrement. Créez votre application comme vous le feriez habituellement :

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

### Étape 2 — compilez ou exécutez avec l’étiquette `mcp`

@tabs
[CLI Wails (recommandée)]
Définissez `WAILS_MCP=1` ; la CLI Wails ajoutera automatiquement l’étiquette `mcp` :

```shell
# Development
WAILS_MCP=1 wails3 dev

# Production build
WAILS_MCP=1 wails3 build
```

[Go directement]
Transmettez directement l’étiquette à `go run` ou à `go build` :

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

Le serveur utilise le **transport HTTP diffusé en continu de MCP**. Connectez-vous avec n’importe quel client compatible avec MCP.

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
Indiquez l’adresse suivante à tout client MCP prenant en charge le transport HTTP diffusé en continu :

```
http://127.0.0.1:9099/mcp
```

@end

### Étape 4 — exécutez une session de test

Demandez à l’agent de tester votre application. Voici quelques exemples d’invites :

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

Toute la configuration s’effectue au moyen de variables d’environnement ; aucune modification du code n’est requise.

| Variable d’environnement | Valeur par défaut | Description |
| --- | --- | --- |
| `WAILS_MCP` | (non définie) | Définissez-la sur `1`, `true`, `on` ou `yes` pour ajouter automatiquement l’étiquette de compilation `mcp` lors de l’utilisation de la CLI Wails. |
| `WAILS_MCP_HOST` | `127.0.0.1` | Interface à laquelle se lier. Les liaisons à une interface autre que l’interface de bouclage nécessitent `WAILS_MCP_TOKEN`. |
| `WAILS_MCP_TOKEN` | non définie | Jeton de porteur facultatif sur l’interface de bouclage, mais obligatoire sur les autres adresses de liaison. Les clients envoient `Authorization: Bearer <token>`. |
| `WAILS_MCP_PORT` | `9099` | Port d’écoute. Définissez-le sur `0` pour utiliser un port libre attribué aléatoirement (indiqué dans le journal). |
| `WAILS_MCP_TIMEOUT` | `30000` | Délai d’expiration par défaut de l’évaluation JavaScript, en **millisecondes**. |
| `WAILS_MCP_HIDE_CURSOR` | (non définie) | Définissez-la sur `1` ou `true` pour désactiver la superposition du curseur animé. |

Exemple — port personnalisé et délai d’expiration de 60 secondes :

```shell
WAILS_MCP=1 WAILS_MCP_PORT=9200 WAILS_MCP_TIMEOUT=60000 wails3 dev
```

## Outils disponibles

| Outil | Fonction |
| --- | --- |
| `app_info` | Informations sur l’application : plateforme, architecture, toutes les fenêtres et point de terminaison MCP |
| `windows_list` | Répertorie toutes les fenêtres avec leur géométrie et leur état |
| `window_control` | Donner le focus, redimensionner, déplacer, passer en plein écran, ouvrir les outils de développement, recharger, définir l’URL, etc. (22 actions) |
| `js_eval` | Évalue du JavaScript dans une fenêtre (corps asynchrone, `return` pour la valeur) |
| `dom_html` | Obtenir le HTML de la page ou d’un élément précis |
| `dom_query` | Rechercher des éléments à l’aide d’un sélecteur CSS — balise, texte, limites, visibilité |
| `screenshot_dom` | Instantané structurel de la page visible (fondé sur le DOM, sans pixels) |
| `mouse_move` | Animer le curseur jusqu’à un point ou un sélecteur CSS |
| `mouse_click` | Cliquer avec le curseur animé (bouton gauche, droit ou central, double-clic, touches de modification) |
| `mouse_drag` | Effectuer un glisser-déposer avec le curseur animé (prend en charge les éléments de glisser-déposer HTML5) |
| `mouse_scroll` | Faire défiler à un point ou sur un élément |
| `keyboard_type` | Saisir du texte caractère par caractère avec des événements réalistes |
| `keyboard_press` | Appuyer sur une seule touche (Enter, Tab, Escape, ArrowDown, …), avec des touches de modification facultatives |
| `call_bound_method` | Appeler une méthode liée d’un service Go, par exemple `main.GreetService.Greet` |
| `emit_event` | Émettre un événement d’application Wails |
| `wait_for_event` | Attendre un événement d’application Wails et renvoyer ses données |

### Prise en charge de plusieurs fenêtres

Tous les outils qui agissent sur une fenêtre acceptent un argument facultatif `window` contenant le **nom** de la fenêtre (défini via `WebviewWindowOptions.Name`). S’il est omis, l’outil cible la fenêtre actuellement active ou, si aucune ne l’est, la première fenêtre.

```
List all windows, then click the "New" button in the window named "editor".
```

### Sélection des éléments

Les outils de souris et de clavier acceptent au choix :

- **Sélecteur CSS** — `selector: "#submit-btn"` (l’élément défile automatiquement jusqu’à devenir visible)
- **Coordonnées** — `x: 400, y: 300` (pixels CSS par rapport à la fenêtre d’affichage)

Pour les opérations de glisser-déposer, utilisez les préfixes `from_` et `to_` :

```
Drag from selector: ".card" to selector: ".dropzone"
```

## Sécurité

@note{type="caution"}
Le serveur MCP offre un contrôle programmatique complet de votre application. Toute personne ayant accès à ses outils peut lire le DOM, évaluer du JavaScript, cliquer sur des boutons et appeler des méthodes Go.

@end

- Par défaut, le serveur écoute sur `127.0.0.1`. Les origines des navigateurs doivent être des origines de bouclage HTTP(S) ; les origines opaques (`null`), mal formées ou externes sont rejetées.
- Les clients MCP natifs sans en-têtes restent pris en charge. Sans `WAILS_MCP_TOKEN`, les processus locaux et les origines locales de navigateur autorisées sont considérés comme fiables ; les vérifications d’origine ne constituent pas une authentification. Définissez un jeton à forte entropie afin d’exiger une authentification par jeton au porteur pour tous les appels `/mcp`. Configurez le même jeton dans l’en-tête `Authorization: Bearer <token>` du client. Les requêtes de contrôle préliminaire ne nécessitent pas le jeton.
- Le rappel `/eval-result` utilise des identifiants imprévisibles propres à chaque évaluation plutôt que le jeton au porteur du client, afin que la transmission des résultats provenant de la vue web reste compatible.
- Les versions de production ne devraient **pas** inclure le tag de compilation `mcp`. La CLI Wails ne l’ajoute que lorsque `WAILS_MCP=1` est explicitement défini, et la compilation par défaut avec `wails3 build` ne contient aucun code du serveur.
- Si vous devez exposer le serveur sur une interface autre que celle de bouclage (par exemple pour des tests sur un réseau local), définissez `WAILS_MCP_HOST=0.0.0.0` et un `WAILS_MCP_TOKEN` à forte entropie ; sans jeton, le démarrage échoue. Sur les réseaux non fiables, utilisez un tunnel chiffré ou un proxy assurant la terminaison TLS, car le programme d’écoute intégré utilise HTTP.

## Exemple d’application

Une application bac à sable complète présentant tous les outils est disponible à l’emplacement [`v3/examples/mcp`](https://github.com/wailsapp/wails/tree/releases/v3-beta/v3/examples/mcp). Elle comprend :

- Compteur avec des boutons d’incrémentation et de réinitialisation
- Champ de saisie du nom avec les méthodes liées Greet, Add et Shout
- Source et cible de glisser-déposer HTML5
- Liste à défilement (50 éléments)
- Journal des événements

Exécutez-la avec :

```shell
cd v3/examples/mcp
go run -tags mcp .
```

Connectez ensuite Claude Code ou n’importe quel client MCP à `http://127.0.0.1:9099/mcp`, puis demandez-lui d’utiliser l’interface utilisateur.

## Retours

Le serveur MCP intégré est expérimental, et vos retours détermineront son évolution. Si vous l’essayez, nous aimerions savoir quels client et outils vous avez utilisés, ce que vous attendiez et ce qui s’est réellement produit, ainsi que s’il vous a été utile de laisser un agent piloter votre application — les rapports les plus utiles indiquent exactement ce que vous avez exécuté. Faites-nous-en part dans la [discussion consacrée aux retours sur le serveur MCP](https://github.com/wailsapp/wails/discussions/5692).
