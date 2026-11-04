---
title: "Création d’une couche de transport personnalisée"
description: "Découvrez comment créer et personnaliser votre propre couche de transport IPC personnalisée pour Wails v3"
slug: "guides/custom-transport"
sourcePath: "guides/custom-transport.md"
---

Wails v3 vous permet de fournir une couche de transport IPC personnalisée tout en conservant toutes les liaisons générées et la communication des événements. Vous pouvez ainsi remplacer le transport par défaut fondé sur les requêtes HTTP Fetch par des WebSockets, des protocoles personnalisés ou tout autre mécanisme de transport.

## Présentation

Par défaut, Wails utilise les requêtes HTTP Fetch du frontend pour communiquer avec le backend via `/wails/runtime`. L’API de transport personnalisé vous permet de :

- Remplacer le transport HTTP par des WebSockets, gRPC ou tout protocole personnalisé
- Conserver une compatibilité totale avec la génération de code de Wails
- Conserver toutes les liaisons, tous les événements, toutes les boîtes de dialogue et les autres fonctionnalités existantes de Wails
- Implémenter votre propre gestion des connexions, de l’authentification et des erreurs

## Architecture

```text
┌─────────────────────────────────────────────────┐
│  Frontend (TypeScript)                          │
│  - Generated bindings still work                │
│  - Your custom client transport                 │
└──────────────────┬──────────────────────────────┘
                   │
                   │ Your Protocol (WebSocket/etc)
                   │
┌──────────────────▼──────────────────────────────┐
│  Backend (Go)                                   │
│  - Your Transport implementation                │
│  - Wails MessageProcessor                       │
│  - All existing Wails infrastructure            │
└─────────────────────────────────────────────────┘
```

## Exécuter l’exemple WebSocket complet

L’[exemple de transport WebSocket](https://github.com/wailsapp/wails/tree/master/v3/examples/websocket-transport) comprend les deux parties du protocole : les identifiants de requête et les réponses, la gestion du cycle de vie, le client JavaScript et la transmission des événements. Partez de ces fichiers plutôt que d’un transport qui implémente uniquement `Start` et `Stop`.

Depuis une copie de travail du dépôt Wails dans laquelle les dépendances natives de la plateforme sont installées :

```bash
cd v3/examples/websocket-transport
go run .
```

L’exemple ouvre une fenêtre native et démarre un point de terminaison WebSocket sur le port 9099. Utilisez-le dans un environnement local de développement fiable. Avant de déployer une variante accessible sur le réseau, restreignez l’écouteur, validez les origines et authentifiez les clients. La politique d’origine de l’exemple est permissive et ne constitue pas une politique de sécurité adaptée à la production.

### Contrat du backend

Un transport doit implémenter les trois méthodes requises. Dans l’exemple, `WebSocketTransport` implémente `Start`, `JSClient` et `Stop` ; `JSClient` fournit le shim JavaScript que Wails injecte dans la page. Ajoutez à votre implémentation une vérification de l’interface à la compilation :

```go
var _ application.Transport = (*WebSocketTransport)(nil)
```

Configurez l’application à l’aide du constructeur de l’exemple :

```go
Transport: NewWebSocketTransport("127.0.0.1:9099"),
```

Ce champ doit figurer dans `application.Options`, aux côtés des ressources et services habituels. Copiez les fichiers d’implémentation ainsi que l’appel du constructeur. `Start` reçoit le `MessageProcessor` ; `Stop` doit fermer les écouteurs, les connexions actives et les opérations en attente.

### Contrat du frontend

Dans l’exemple, `websocket-transport.js` implémente le transport du runtime, notamment la corrélation entre requêtes et réponses ainsi que le rejet des appels ayant échoué. Après avoir servi ce fichier depuis vos ressources, initialisez-le avant d’appeler les liaisons de services générées :

```javascript
import { setTransport } from '/wails/runtime.js'
import { createWebSocketTransport } from '/websocket-transport.js'

setTransport(createWebSocketTransport('ws://127.0.0.1:9099/wails/ws'))
```

Les liaisons de services générées restent inchangées. Si votre transport remplace le client, il doit prendre en charge les méthodes utilisées par votre application et préserver leur sémantique d’erreur. La transmission des événements nécessite le contrat facultatif de transport des événements ; une implémentation `Transport` arbitraire ne le fournit pas automatiquement. L’exemple WebSocket complet en propose une démonstration.

Pour les déploiements distants, utilisez TLS ainsi qu’une politique explicite d’authentification et de contrôle des origines. Un hébergement exclusivement dans un navigateur nécessite également de servir les ressources de l’application, comme décrit ci-dessous.

## Référence de l’API

### Interface Transport

```go
type Transport interface {
    Start(ctx context.Context, messageProcessor *application.MessageProcessor) error
    // JSClient returns the JavaScript shim that the runtime injects into the
    // window so that frontend code can call into the transport.
    JSClient() []byte
    Stop() error
}
```

### Interface AssetServerTransport (facultative)

Pour les déploiements dans un navigateur, ou si vous souhaitez servir à la fois les ressources et l’IPC au moyen de votre transport personnalisé, implémentez l’interface `AssetServerTransport` :

```go
type AssetServerTransport interface {
    Transport

    // ServeAssets configures the transport to serve assets alongside IPC.
    // The assetHandler is Wails' internal asset server that handles all assets,
    // runtime.js, capabilities, flags, etc.
    ServeAssets(assetHandler http.Handler) error
}
```

**Quand implémenter cette interface :**

- Exécution de l’application dans un navigateur plutôt que dans une webview
- Distribution des ressources via HTTP parallèlement à votre transport IPC personnalisé
- Création d’applications accessibles sur le réseau

**Exemple d’implémentation :**

```go
func (t *MyTransport) ServeAssets(assetHandler http.Handler) error {
    mux := http.NewServeMux()

    // Mount your IPC endpoint
    mux.HandleFunc("/my/ipc/endpoint", t.handleIPC)

    // Mount Wails asset server for everything else
    mux.Handle("/", assetHandler)

    // Start HTTP server
    t.httpServer.Handler = mux
    go t.httpServer.ListenAndServe()

    return nil
}
```

Lorsque `ServeAssets()` est appelé, assetHandler fournit :

- Toutes les ressources statiques (HTML, CSS, JS, images, etc.)
- `/wails/runtime.js` - La bibliothèque du runtime Wails

## Voir aussi

- `transport.go` - Interfaces et types fondamentaux du transport
- `messageprocessor.go` - Le processeur de messages sous-jacent qui gère toutes les communications IPC de Wails
