---
title: "Messages bruts"
description: "Implémenter une communication personnalisée entre le frontend et le backend pour les applications aux performances critiques"
slug: "guides/raw-messages"
sourcePath: "guides/raw-messages.md"
---

Les messages bruts fournissent un canal de communication de bas niveau entre votre frontend et votre backend, en contournant le système de liaisons standard. Vous gagnez ainsi en vitesse au prix d’une moindre simplicité d’utilisation.

## Quand utiliser les messages bruts

Les messages bruts conviennent surtout aux cas limites extrêmes :

- **Mises à jour à très haute fréquence** — Des milliers de messages par seconde, lorsque chaque microseconde compte
- **Protocoles de messages personnalisés** — Lorsque vous avez besoin d’un contrôle total sur le format des données transmises

@note{type="tip"}
Pour presque tous les cas d’utilisation, les [liaisons de services](/features/bindings/services/) standard sont recommandées, car elles garantissent la sûreté des types, assurent la sérialisation automatique et améliorent l’expérience de développement, avec une surcharge négligeable.

@end

## Configuration du backend

Configurez `RawMessageHandler` dans les options de votre application :

```go
package main

import (
    "encoding/json"
    "fmt"

    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "Raw Message Demo",
        Assets: application.AssetOptions{
            Handler: application.BundledAssetFileServer(assets),
        },
        RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
            fmt.Printf("Raw message from window '%s': %s (origin: %+v)\n", window.Name(), message, originInfo.Origin)

            // Process the message and respond via events
            response := processMessage(message)
            window.EmitEvent("raw-response", response)
        },
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: "My App",
        Name:  "main",
    })

    app.Run()
}

func processMessage(message string) map[string]any {
    // Your custom message processing logic
    return map[string]any{
        "received": message,
        "status":   "processed",
    }
}
```

### Signature du gestionnaire

```go
RawMessageHandler func(window Window, message string, originInfo *application.OriginInfo)
```

| Paramètre | Type | Description |
| --- | --- | --- |
| `window` | `Window` | La fenêtre qui a envoyé le message |
| `message` | `string` | Le contenu brut du message |
| `originInfo` | `*application.OriginInfo` | Les informations d’origine de la source du message |

#### Structure OriginInfo

```go
type OriginInfo struct {
	Origin      string
	TopOrigin   string
	IsMainFrame bool
}
```

| Champ | Type | Description |
| --- | --- | --- |
| `Origin` | `string` | L’URL d’origine du document qui a envoyé le message |
| `TopOrigin` | `string` | L’URL d’origine de premier niveau (elle peut différer de Origin dans les iframes) |
| `IsMainFrame` | `bool` | Indique si le message provient du cadre principal |

#### Disponibilité propre à chaque plateforme

- **macOS** : `Origin` et `IsMainFrame` sont fournis
- **Windows** : `Origin` et `TopOrigin` sont fournis
- **Linux** : seul `Origin` est fourni

### Validation de l’origine

@note{type="caution"}
Ne supposez jamais qu’un message est sûr simplement parce qu’il parvient à votre gestionnaire. Vous devez valider les informations d’origine avant de traiter toute opération sensible ou modifiant l’état.

@end

**Vérifiez toujours l’origine des messages entrants avant de les traiter.** Le paramètre `originInfo` fournit des informations de sécurité critiques que vous devez valider afin d’empêcher tout accès non autorisé. Du contenu malveillant ou compromis, ainsi que des scripts non prévus, pourraient envoyer des messages bruts. Sans validation de l’origine, vous risquez de traiter des commandes provenant de sources non fiables. Utilisez `originInfo` pour vous assurer que les messages proviennent des sources attendues.

### Points de validation essentiels

- **Vérifiez toujours `Origin`** — Vérifiez que l’origine correspond aux sources fiables attendues (généralement `wails://wails` ou `http://wails.localhost` pour les ressources locales, ou l’origine propre à votre application)
- **Validez `IsMainFrame`** (macOS) — Vérifiez si le message provient d’une iframe, car cela peut indiquer la présence de contenu intégré relevant de contextes de sécurité différents
- **Utilisez `TopOrigin`** (Windows) — Vérifiez l’origine de premier niveau lorsque vous traitez du contenu dans des cadres
- **Rejetez les origines inattendues** — Adoptez un comportement sûr en cas d’échec en rejetant les messages provenant d’origines que vous n’autorisez pas explicitement

@note{type="info"}
Les messages dont le préfixe est `wails:` sont réservés aux communications internes de Wails et ne sont pas transmis à votre gestionnaire.

@end

## Configuration du frontend

Envoyez des messages bruts à l’aide de `System.invoke()` :

```html
<!DOCTYPE html>
<html>
<head>
    <script type="module">
        import { System, Events } from '@wailsio/runtime'

        // Send raw message
        document.getElementById('send').addEventListener('click', () => {
            const message = document.getElementById('input').value
            System.invoke(message)
        })

        // Listen for response
        Events.On('raw-response', (event) => {
            console.log('Response:', event.data)
        })
    </script>
</head>
<body>
    <input type="text" id="input" placeholder="Enter message" />
    <button id="send">Send</button>
</body>
</html>
```

### Utilisation du bundle précompilé

Si vous n’utilisez pas npm, accédez à `invoke` via l’objet global `wails` :

```html
<script type="module" src="/wails/runtime.js"></script>
<script>
    window.onload = function() {
        document.getElementById('send').onclick = function() {
            wails.System.invoke('my-message')
        }
    }
</script>
```

## Messages structurés

Pour les données complexes, sérialisez-les au format JSON :

### Frontend

```javascript
import { System } from '@wailsio/runtime'

const command = {
    action: 'update',
    payload: {
        id: 123,
        value: 'new value'
    }
}

System.invoke(JSON.stringify(command))
```

### Backend

```go
RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
    var cmd struct {
        Action  string `json:"action"`
        Payload struct {
            ID    int    `json:"id"`
            Value string `json:"value"`
        } `json:"payload"`
    }

    if err := json.Unmarshal([]byte(message), &cmd); err != nil {
        window.EmitEvent("error", err.Error())
        return
    }

    switch cmd.Action {
    case "update":
        // Handle update
        result := handleUpdate(cmd.Payload.ID, cmd.Payload.Value)
        window.EmitEvent("update-complete", result)
    default:
        window.EmitEvent("error", "unknown action")
    }
}
```

## Comparaison des performances

| Approche | Surcharge | Sécurité des types | Cas d’utilisation |
| --- | --- | --- | --- |
| Liaisons de services | Plus élevée | Complète | Usage général |
| Messages bruts | Minime | Manuelle | Haute fréquence, performances critiques |

### Exemple de benchmark

Pour les charges utiles simples, les messages bruts peuvent traiter nettement plus de messages par seconde que les liaisons de services :

```go
// Raw message handler - minimal overhead
RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
    // Direct string processing, no reflection or marshaling
    counter++
}
```

## Exemple complet

Voici un exemple complet qui implémente un protocole de commande simple :

### main.go

```go
package main

import (
    "embed"
    "encoding/json"
    "fmt"
    "time"

    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

type Command struct {
    Type string          `json:"type"`
    Data json.RawMessage `json:"data"`
}

func main() {
    app := application.New(application.Options{
        Name: "Raw Message Demo",
        Assets: application.AssetOptions{
            Handler: application.BundledAssetFileServer(assets),
        },
        Mac: application.MacOptions{
            ApplicationShouldTerminateAfterLastWindowClosed: true,
        },
        RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
            var cmd Command
            if err := json.Unmarshal([]byte(message), &cmd); err != nil {
                window.EmitEvent("error", map[string]string{"error": err.Error()})
                return
            }

            switch cmd.Type {
            case "ping":
                window.EmitEvent("pong", map[string]any{
                    "time":   time.Now().UnixMilli(),
                    "window": window.Name(),
                })
            case "echo":
                var text string
                json.Unmarshal(cmd.Data, &text)
                window.EmitEvent("echo", text)
            default:
                window.EmitEvent("error", map[string]string{
                    "error": fmt.Sprintf("unknown command: %s", cmd.Type),
                })
            }
        },
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: "Raw Message Demo",
        Name:  "main",
        Width: 400,
        Height: 300,
    })

    app.Run()
}
```

### assets/index.html

```html
<!DOCTYPE html>
<html>
<head>
    <title>Raw Message Demo</title>
    <style>
        body { font-family: sans-serif; padding: 20px; }
        button { margin: 5px; padding: 10px 20px; }
        #output { margin-top: 20px; padding: 10px; background: #f0f0f0; }
    </style>
</head>
<body>
    <h1>Raw Message Demo</h1>

    <button id="ping">Ping</button>
    <button id="echo">Echo "Hello"</button>

    <div id="output">Waiting for response...</div>

    <script type="module">
        import { System, Events } from '@wailsio/runtime'

        const output = document.getElementById('output')

        function send(type, data) {
            System.invoke(JSON.stringify({ type, data }))
        }

        document.getElementById('ping').onclick = () => send('ping')
        document.getElementById('echo').onclick = () => send('echo', 'Hello')

        Events.On('pong', (e) => {
            output.textContent = `Pong from ${e.data.window} at ${e.data.time}`
        })

        Events.On('echo', (e) => {
            output.textContent = `Echo: ${e.data}`
        })

        Events.On('error', (e) => {
            output.textContent = `Error: ${e.data.error}`
        })
    </script>
</body>
</html>
```

## Bonnes pratiques

### À faire

- Utilisez les messages bruts pour les chemins dont les performances sont réellement critiques
- Implémentez une gestion appropriée des erreurs dans votre gestionnaire
- Utilisez des événements pour renvoyer les réponses au frontend
- Envisagez d’utiliser JSON pour les données structurées
- Traitez rapidement les messages pour éviter tout blocage

### À ne pas faire

- Utiliser des messages bruts lorsque les liaisons de services suffisent
- Oublier de valider les messages entrants
- Bloquer le gestionnaire avec des opérations longues (utilisez des goroutines)
- Ignorer le paramètre de fenêtre lorsque les réponses doivent cibler des fenêtres précises

## Considérations relatives aux applications multifenêtres

Le paramètre `window` identifie la fenêtre qui a envoyé le message, ce qui vous permet :

- d’envoyer les réponses à la bonne fenêtre
- d’implémenter un comportement propre à chaque fenêtre
- de suivre la source des messages à des fins de débogage

```go
RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
    // Respond only to the sending window
    window.EmitEvent("response", result)

    // Or broadcast to all windows
    app.Event.Emit("broadcast", result)
}
```

## Étapes suivantes

- [Liaisons de services](/features/bindings/services/) — Approche standard pour la plupart des applications
- [Événements](/guides/events-reference/) — Système d’événements pour la communication du backend vers le frontend
- [Performances](/guides/performance/) — Optimisation générale des performances
