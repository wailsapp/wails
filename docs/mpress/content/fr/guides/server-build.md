---
title: "Build serveur"
description: "Exécutez des applications Wails comme des serveurs HTTP sans fenêtre d’interface graphique native"
slug: "guides/server-build"
sourcePath: "guides/server-build.md"
---

Wails v3 prend en charge le mode serveur, qui permet d’exécuter votre application comme un serveur HTTP pur, sans créer de fenêtres natives ni nécessiter de dépendances graphiques. Vous pouvez ainsi déployer la même application Wails sur des serveurs, dans des conteneurs et dans des navigateurs web.

Le mode serveur est utile pour :

- **Déploiements Docker/conteneur** — Exécutez l’application sans dépendances X11/Wayland
- **Applications côté serveur** — Déployez l’application comme un serveur web accessible depuis un navigateur
- **Accès exclusivement web** — Partagez la même base de code entre les versions de bureau et web
- **Tests CI/CD** — Exécutez des tests d’intégration sans serveur d’affichage
- **Microservices** — Utilisez les liaisons Wails dans des services backend sans interface graphique

## Démarrage rapide

Le mode serveur s’active avec le tag de build `server`. Le code de votre application reste inchangé : il vous suffit d’effectuer le build avec ce tag :

```bash
# Using Taskfile (recommended)
wails3 task build:server
wails3 task run:server

# Or build directly with Go
go build -tags server -o myapp-server .
```

Voici un exemple minimal :

```go
package main

import (
    "embed"
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name: "My App",
        // Server options are used when built with -tags server
        Server: application.ServerOptions{
            Host: "localhost",
            Port: 8080,
        },
        Services: []application.Service{
            application.NewService(&MyService{}),
        },
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
    })

    log.Println("Starting application...")
    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

Le même code peut être compilé en mode bureau (sans le tag) ou en mode serveur (avec `-tags server`).

## Configuration

### ServerOptions

Configurez le serveur HTTP avec `ServerOptions` :

```go
Server: application.ServerOptions{
    // Host to bind to. Default: "localhost"
    // Use "0.0.0.0" to listen on all interfaces
    Host: "localhost",

    // Port to listen on. Default: 8080
    Port: 8080,

    // Request read timeout. Default: 30s
    ReadTimeout: 30 * time.Second,

    // Response write timeout. Default: 30s
    WriteTimeout: 30 * time.Second,

    // Idle connection timeout. Default: 120s
    IdleTimeout: 120 * time.Second,

    // Graceful shutdown timeout. Default: 30s
    ShutdownTimeout: 30 * time.Second,

    // Additional origins allowed to open WebSocket connections and call
    // /wails/runtime. Same-origin requests are always allowed.
    WebSocketOriginPatterns: []string{"app.example.com"},

    // Disable WebSocket and runtime origin checks. Unsafe; default: false.
    WebSocketAllowAllOrigins: false,

    // TLS configuration (optional)
    TLS: &application.TLSOptions{
        CertFile: "/path/to/cert.pem",
        KeyFile:  "/path/to/key.pem",
    },
},
```

## Fonctionnalités

### Point de terminaison de contrôle d’intégrité

Un point de terminaison de contrôle d’intégrité est automatiquement disponible à l’adresse `/health` :

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

Il est utile pour :

- les sondes d’activité et de disponibilité de Kubernetes
- les contrôles d’intégrité des répartiteurs de charge
- les systèmes de supervision

### Liaisons de services

Toutes les liaisons de services fonctionnent exactement comme en mode bureau :

```go
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello, " + name + "!"
}

// Register in options
Services: []application.Service{
    application.NewService(&GreetService{}),
},
```

Le frontend peut appeler ces liaisons à l’aide du runtime Wails standard :

```javascript
const greeting = await wails.Call.ByName('main.GreetService.Greet', 'World');
```

### Événements

En mode serveur, les événements fonctionnent dans les deux sens :

- **Du frontend vers le backend** : les événements émis depuis le navigateur sont envoyés via HTTP et reçus par vos gestionnaires d’événements Go
- **Du backend vers le frontend** : les événements émis depuis Go sont diffusés via WebSocket à tous les navigateurs connectés

Chaque onglet de navigateur est représenté comme une « fenêtre » dotée d’un nom unique (`browser-1`, `browser-2`, etc.), accessible via `event.Sender` :

```go
// Listen for events from browsers
app.Event.On("user-action", func(event *application.CustomEvent) {
    log.Printf("Event from %s: %v", event.Sender, event.Data)
    // event.Sender will be "browser-1", "browser-2", etc.
})

// Emit events to all connected browsers
app.Event.Emit("server-update", data)
```

Depuis le frontend :

```javascript
// Emit event to server (and all other browsers)
await wails.Events.Emit('user-action', { action: 'click' });

// Listen for events from server
wails.Events.On('server-update', (event) => {
    console.log('Update from server:', event.data);
});
```

### Arrêt progressif

Le serveur gère correctement les signaux `SIGINT` et `SIGTERM` :

1. Cesse d’accepter de nouvelles connexions
2. Attend la fin des requêtes actives (pendant `ShutdownTimeout` au maximum)
3. Exécute les hooks `OnShutdown`
4. Arrête les services dans l’ordre inverse

## Différences par rapport au mode bureau

| Fonctionnalité | Mode bureau | Mode serveur |
| --- | --- | --- |
| Fenêtres natives | Créées | Fenêtres de navigateur (`browser-N`) |
| Zone de notification système | Disponible | Non disponible |
| Boîtes de dialogue natives | Disponibles | Non disponibles |
| Menu de l’application | Disponible | Non disponible |
| Informations sur l’écran | Disponibles | Aucun écran natif ; `GetPrimary()` renvoie `nil` |
| Liaisons de services | Fonctionnent | Fonctionnent |
| Événements | Fonctionnent | Fonctionnent (via WebSocket) |
| Ressources | Via la webview | Via HTTP |
| CGO requis | Oui | Non |

### Comportement de l’API de fenêtre

En mode serveur, les API liées aux fenêtres sont gérées de manière sûre :

- `app.Window.NewWithOptions(options)` — Renvoie un objet de fenêtre Go non nil ; cela ne crée aucune fenêtre de bureau native
- `app.Hide()` / `app.Show()` — Sans effet
- `app.Screen.GetPrimary()` — Renvoie `*application.Screen`, qui vaut `nil` lorsqu’aucun écran principal n’est disponible ; vérifiez sa valeur avant de le déréférencer

Distinguez les objets de fenêtre Go de la disponibilité des fenêtres natives. Protégez le comportement réservé au bureau lors du build d’un serveur : un objet de fenêtre Go non nil ne garantit pas l’existence d’une fenêtre ou d’un écran natif.

### Plateformes prises en charge

Le mode serveur est entièrement écrit en Go et se compile sans CGO pour toutes les plateformes `GOOS` prises en charge par Wails, notamment `android` et `ios`. Sur les cibles mobiles, il compile uniquement un serveur sans interface graphique :

- Les API mobiles natives `application.Android` et `application.IOS` ne sont pas disponibles. Le code qui les utilise doit être exclu des builds serveur à l’aide d’une contrainte de build `!server`.
- Le journaliseur par défaut est celui du mode bureau : dans les builds de développement, il écrit sur stderr plutôt que dans logcat ou dans la console de la WebView.

## Build pour la production

### Avec Task (recommandé)

Les projets créés avec `wails3 init` incluent une tâche `build:server` :

```bash
# Build for server mode
wails3 task build:server

# Build and run
wails3 task run:server
```

### Build manuel

```bash
# Build with server mode
go build -tags server -o myapp-server .
```

### Docker

Les projets Wails incluent une configuration Docker prête à l’emploi. Pour compiler et exécuter votre application dans un conteneur :

```bash
# Build the Docker image
wails3 task build:docker

# Run it
wails3 task run:docker
```

C’est tout ! Votre application sera disponible à l’adresse `http://localhost:8080`.

Vous pouvez personnaliser le build à l’aide de quelques options :

```bash
# Use a custom image tag
wails3 task build:docker TAG=myapp:v1.0.0

# Run on a different port
wails3 task run:docker PORT=3000
```

Le fichier `Dockerfile.server` généré crée une image minimale basée sur distroless. Il configure automatiquement l’écoute réseau afin que votre application soit accessible depuis l’extérieur du conteneur.

### Docker Compose

Pour les déploiements plus complexes, voici une configuration Docker Compose avec des contrôles d’intégrité :

```yaml
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - WAILS_SERVER_HOST=0.0.0.0
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

@note{type="info"}
L’exemple de contrôle d’intégrité utilise `wget`. Si vous utilisez une image de base distroless, vous devrez soit inclure un binaire de contrôle d’intégrité dans votre image, soit utiliser un mécanisme externe de contrôle d’intégrité (par exemple, l’option `curl` de Docker ou un conteneur side-car).

@end

### Dockerfile personnalisé

Si vous avez besoin de davantage de contrôle, vous pouvez créer votre propre Dockerfile. Le point essentiel est de définir `WAILS_SERVER_HOST=0.0.0.0` afin que le serveur accepte les connexions provenant de l’extérieur du conteneur :

```dockerfile
# Build stage
FROM golang:alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY . .
RUN go mod tidy
RUN go build -tags server -ldflags="-s -w" -o server .

# Runtime stage
FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/server /server
COPY --from=builder /app/frontend/dist /frontend/dist
EXPOSE 8080
ENV WAILS_SERVER_HOST=0.0.0.0
ENTRYPOINT ["/server"]
```

## Considérations de sécurité

Lors du déploiement d’applications en mode serveur :

1. **Liez l’application à localhost par défaut** — N’utilisez `0.0.0.0` que lorsque cela est nécessaire
2. **Utilisez TLS en production** — Configurez `ServerOptions.TLS`
3. **Placez l’application derrière un proxy inverse** — Utilisez nginx/traefik pour renforcer la sécurité
4. **Conservez la même origine pour les WebSockets et les appels au runtime** — Ajoutez uniquement des origines de confiance avec `WebSocketOriginPatterns` ; évitez `WebSocketAllowAllOrigins`
5. **Validez toutes les entrées** — Appliquez les mêmes pratiques de sécurité que pour toute application web

## Exemple

Un exemple complet est disponible dans `v3/examples/server/` :

```bash
cd v3/examples/server

# Using Taskfile
task dev

# Or run directly
go run -tags server .

# Open http://localhost:8080 in browser
```

## Variables d’environnement

Pour les scénarios de déploiement dans lesquels vous devez remplacer la configuration du serveur sans modifier le code, Wails reconnaît les variables d’environnement suivantes :

| Variable | Description | Valeur par défaut |
| --- | --- | --- |
| `WAILS_SERVER_HOST` | Interface réseau à laquelle se lier | `localhost` |
| `WAILS_SERVER_PORT` | Port d’écoute | `8080` |

Ces variables ont priorité sur `ServerOptions` dans votre code. C’est pourquoi les exemples Docker définissent `WAILS_SERVER_HOST=0.0.0.0` : cela permet au conteneur d’accepter les connexions externes sans qu’aucune modification de votre application soit nécessaire.

## Voir aussi

- [Transport personnalisé](/guides/custom-transport/) — Pour la personnalisation avancée de l’IPC
- [Services](/features/bindings/services/) — Documentation sur la liaison de services
- [Événements](/guides/events-reference/) — Documentation du système d’événements

### Taille des requêtes du runtime

La taille des requêtes adressées à `/wails/runtime` est limitée à 64 Mio avant le traitement JSON. Une requête ordinaire dépassant cette limite reçoit le statut HTTP 413, y compris si elle ne contient pas `Content-Length`. Les téléversements segmentés du runtime conservent leurs limites distinctes de 1 Mio par segment et de 64 Mio pour la charge utile assemblée. Utilisez un middleware applicatif ou un proxy inverse pour imposer une limite inférieure lorsque cela est approprié.

### Origine des requêtes du runtime

`/wails/runtime` n’accepte que les requêtes `POST` ; les autres méthodes reçoivent le statut HTTP 405. Une requête de navigateur dont l’`Origin` ne correspond ni à l’hôte de la requête ni à `WebSocketOriginPatterns` reçoit le statut HTTP 403. Il en va de même pour une requête sans `Origin` dont `Sec-Fetch-Site` vaut `cross-site` ou `same-site`. Les requêtes ne comportant aucun de ces deux en-têtes, telles que celles provenant de `curl` ou d’autres clients non basés sur un navigateur, sont acceptées. Si un proxy inverse réécrit l’en-tête `Host`, ajoutez l’origine publique à `WebSocketOriginPatterns`.

`WebSocketOriginPatterns` est prévu pour les déploiements derrière un proxy inverse, dans lesquels l’origine publique diffère de l’`Host` vue par l’application. Les frontends servis depuis une autre origine et appelant directement `/wails/runtime` ne sont pas pris en charge : le point de terminaison n’émet aucun en-tête CORS et répond aux requêtes préliminaires `OPTIONS` avec le statut HTTP 405. Servez le frontend depuis le serveur Wails, ou placez les deux derrière la même origine de proxy.

### Vérification de l’en-tête Host pour les adresses de bouclage

Lorsque le serveur se lie à une adresse de bouclage (la valeur par défaut `localhost`, toute adresse `127.0.0.0/8` ou `::1`, qu’elle soit définie dans `ServerOptions.Host` ou `WAILS_SERVER_HOST`), chaque requête doit comporter un en-tête `Host` indiquant `localhost`, une adresse IP littérale, l’hôte de liaison ou un hôte correspondant à la partie hôte d’une entrée de `WebSocketOriginPatterns` ; toute autre valeur de `Host` reçoit le statut HTTP 403. Cette vérification empêche le rebind DNS, où une page web malveillante fait pointer son propre nom d’hôte vers `127.0.0.1` afin que ses requêtes passent les contrôles de même origine décrits ci-dessus. Le port n’est pas comparé, et les graphies telles que `localhost.` (avec un point final) ou `127.1` sont rejetées. La vérification s’applique toujours lorsque `WebSocketAllowAllOrigins` est défini. Si un proxy inverse local transmet son `Host` public, ajoutez cet hôte à `WebSocketOriginPatterns`. Les liaisons à d’autres adresses, telles que `0.0.0.0` dans les conteneurs, ne sont pas concernées.
