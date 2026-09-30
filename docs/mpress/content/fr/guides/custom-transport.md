---
title: "Création d’une couche de transport personnalisée"
description: "Découvrez comment créer et personnaliser votre propre couche de transport IPC pour Wails v3"
slug: "guides/custom-transport"
sourcePath: "guides/custom-transport.md"
---

Wails v3 vous permet de fournir une couche de transport IPC personnalisée tout en conservant l’ensemble des liaisons générées et des communications par événements. Vous pouvez ainsi remplacer le transport par défaut fondé sur les requêtes HTTP Fetch par des WebSockets, des protocoles personnalisés ou tout autre mécanisme de transport.

## Vue d’ensemble

Par défaut, Wails utilise des requêtes HTTP Fetch envoyées par le frontend pour communiquer avec le backend via `/wails/runtime`. L’API de transport personnalisé vous permet de :

- Remplacer le transport HTTP par des WebSockets, gRPC ou tout protocole personnalisé
- Conserver une compatibilité totale avec la génération de code de Wails
- Conserver toutes les liaisons, tous les événements et toutes les boîtes de dialogue existants, ainsi que les autres fonctionnalités de Wails
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

## Utilisation

### 1. Implémenter l’interface de transport

Créez un transport personnalisé en implémentant l’interface `Transport` :

```go
package main

import (
    "context"
    "github.com/wailsapp/wails/v3/pkg/application"
)

type MyCustomTransport struct {
    // Your fields
}

func (t *MyCustomTransport) Start(ctx context.Context, processor *application.MessageProcessor) error {
    // Initialize your transport (WebSocket server, gRPC server, etc.)
    // When you receive requests, call processor.HandleRuntimeCallWithIDs()
    return nil
}

func (t *MyCustomTransport) Stop() error {
    // Clean up your transport
    return nil
}
```

### 2. Configurer votre application

Transmettez votre transport personnalisé dans les options de l’application :

```go
func main() {
    app := application.New(application.Options{
        Name: "My App",
        Transport: &MyCustomTransport{},
        // ... other options
    })

    err := app.Run()
    if err != nil {
        log.Fatal(err)
    }
}
```

### 3. Modifier le runtime du frontend

Si vous utilisez un transport personnalisé, vous devrez modifier le runtime du frontend afin qu’il utilise votre transport au lieu de HTTP Fetch. Implémentez l’interface `RuntimeTransport` qui servira à traiter les requêtes :

```typescript
const { setTransport } = await import('/wails/runtime.js');

class MyRuntimeTransport {
  call(objectID: number, method: number, windowName: string, args: any): Promise<any> {
    // TODO: implement IPC call with your transport protocol

    return resp;
  }
}

const myTransport = new MyRuntimeTransport();
setTransport(myTransport);
```

## Remarques

- Le transport HTTP par défaut continue de fonctionner si aucun transport personnalisé n’est spécifié
- Les liaisons générées restent inchangées : seule la couche de transport change
- Les événements, les boîtes de dialogue, le presse-papiers et toutes les autres fonctionnalités de Wails fonctionnent de manière transparente
- Vous êtes responsable de la gestion des erreurs, de la logique de reconnexion et de la sécurité de votre transport personnalisé
- L’exemple WebSocket fourni est destiné à la démonstration et peut nécessiter un renforcement avant toute utilisation en production

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

Pour les déploiements dans un navigateur, ou si vous souhaitez fournir à la fois les ressources et l’IPC par l’intermédiaire de votre transport personnalisé, implémentez l’interface `AssetServerTransport` :

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
- Mise à disposition des ressources via HTTP parallèlement à votre transport IPC personnalisé
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

Lorsque `ServeAssets()` est appelée, assetHandler fournit :

- Toutes les ressources statiques (HTML, CSS, JS, images, etc.)
- `/wails/runtime.js` — La bibliothèque du runtime de Wails

## Voir aussi

- `transport.go` — Interfaces et types fondamentaux du transport
- `messageprocessor.go` — Le processeur de messages sous-jacent qui gère toutes les communications IPC de Wails
