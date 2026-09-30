---
title: "Cycle de vie de l’application"
description: "Comprendre le cycle de vie d’une application Wails, du démarrage à l’arrêt"
slug: "concepts/lifecycle"
sourcePath: "concepts/lifecycle.md"
---

## Comprendre le cycle de vie d’une application

Les applications de bureau suivent un cycle de vie allant du démarrage à l’arrêt. Wails v3 fournit des **services**, des **événements** et des **hooks** pour gérer efficacement ce cycle de vie.

## Les étapes du cycle de vie

```d2
direction: down

Start: Démarrage de l’application {
  shape: oval
  style.fill: "#10B981"
}

Init: Initialisation {
  Parse: Analyser les options {
    shape: rectangle
  }
  Register: Enregistrer les services {
    shape: rectangle
  }
  Setup: Configurer l’environnement d’exécution {
    shape: rectangle
  }
}

AppRun: app.Run() {
  shape: rectangle
  style.fill: "#3B82F6"
}

ServiceStartup: Démarrage des services {
  shape: rectangle
  style.fill: "#8B5CF6"
}

EventLoop: Boucle d’événements {
  Process: Traiter les événements {
    shape: rectangle
  }
  Handle: Gérer les messages {
    shape: rectangle
  }
  Update: Mettre à jour l’interface utilisateur {
    shape: rectangle
  }
}

QuitSignal: Signal de fermeture {
  shape: diamond
  style.fill: "#F59E0B"
}

ShouldQuit: Vérification de ShouldQuit {
  shape: rectangle
  style.fill: "#3B82F6"
}

OnShutdown: Rappels OnShutdown {
  shape: rectangle
  style.fill: "#3B82F6"
}

ServiceShutdown: Arrêt des services {
  shape: rectangle
  style.fill: "#8B5CF6"
}

Cleanup: Nettoyage {
  Close: Fermer les fenêtres {
    shape: rectangle
  }
  Release: Libérer les ressources {
    shape: rectangle
  }
}

End: Fin de l’application {
  shape: oval
  style.fill: "#EF4444"
}

Start -> Init.Parse
Init.Parse -> Init.Register
Init.Register -> Init.Setup
Init.Setup -> AppRun
AppRun -> ServiceStartup
ServiceStartup -> EventLoop.Process
EventLoop.Process -> EventLoop.Handle
EventLoop.Handle -> EventLoop.Update
EventLoop.Update -> EventLoop.Process: Boucle
EventLoop.Process -> QuitSignal: L’utilisateur quitte l’application
QuitSignal -> ShouldQuit: Vérifier si l’opération est autorisée
ShouldQuit -> EventLoop.Process: Refusée
ShouldQuit -> OnShutdown: Autorisée
OnShutdown -> ServiceShutdown
ServiceShutdown -> Cleanup.Close
Cleanup.Close -> Cleanup.Release
Cleanup.Release -> End
```

### 1. Création de l’application

Créez votre application avec `application.New()` :

```go
app := application.New(application.Options{
    Name:        "My App",
    Description: "An application built with Wails",
    Services: []application.Service{
        application.NewService(&MyService{}),
    },
    Assets: application.AssetOptions{
        Handler: application.BundledAssetFileServer(assets),
    },
})
```

**Déroulement :**

1. Les options sont analysées et validées
2. Les services sont enregistrés (mais pas encore démarrés)
3. Le serveur de ressources est configuré
4. L’environnement d’exécution est configuré

### 2. Exécution de l’application

Appelez `app.Run()` pour démarrer l’application :

```go
err := app.Run()  // Blocks until quit
if err != nil {
    log.Fatal(err)
}
```

**Déroulement :**

1. Les services sont démarrés dans leur ordre d’enregistrement
2. Les écouteurs d’événements sont activés
3. Les fenêtres peuvent être créées
4. La boucle d’événements démarre

### 3. Boucle d’événements

L’application entre dans la boucle d’événements, où elle passe la majeure partie de son temps :

- Les événements du système d’exploitation sont traités (souris, clavier et fenêtres)
- Les messages de Go vers JS sont traités
- Les appels de JS vers Go sont exécutés
- Les mises à jour de l’interface utilisateur sont affichées

### 4. Arrêt

Lorsque l’application se ferme :

1. Le rappel `ShouldQuit` est vérifié (s’il est défini)
2. Les rappels `OnShutdown` sont exécutés
3. Les services sont arrêtés dans l’ordre inverse
4. Les fenêtres sont fermées
5. Les ressources sont libérées

## Cycle de vie des services

Les services constituent le principal moyen de gérer le cycle de vie dans Wails v3. Ils fournissent, par l’intermédiaire d’interfaces, des hooks de démarrage et d’arrêt. Pour consulter la documentation complète sur les services, reportez-vous au [guide des services](/features/bindings/services/).

### Création d’un service

```go
type MyService struct {
    db *sql.DB
}

// ServiceStartup is called when the application starts
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    var err error
    s.db, err = sql.Open("sqlite3", "app.db")
    if err != nil {
        return err  // Startup aborts if error returned
    }

    // Run migrations
    if err := s.runMigrations(); err != nil {
        return err
    }

    return nil
}

// ServiceShutdown is called when the application shuts down
func (s *MyService) ServiceShutdown() error {
    if s.db != nil {
        return s.db.Close()
    }
    return nil
}
```

### Enregistrement des services

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&MyService{}),
        application.NewService(&AnotherService{}),
    },
})
```

**Points clés :**

- Les services démarrent dans leur ordre d’enregistrement
- Les services s’arrêtent dans l’ordre **inverse** de leur enregistrement
- Si la méthode `ServiceStartup` d’un service renvoie une erreur, l’application abandonne son démarrage
- Le contexte `ctx` transmis à `ServiceStartup` est annulé lorsque l’arrêt commence

### Utilisation du contexte de l’application

Le contexte transmis à `ServiceStartup` reste valide pendant toute la durée de vie de l’application :

```go
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    // Start a background task that respects shutdown
    go func() {
        ticker := time.NewTicker(5 * time.Minute)
        defer ticker.Stop()

        for {
            select {
            case <-ticker.C:
                s.performBackgroundSync()
            case <-ctx.Done():
                // Application is shutting down
                return
            }
        }
    }()

    return nil
}
```

Vous pouvez également accéder au contexte depuis l’instance de l’application :

```go
app := application.Get()
ctx := app.Context()
```

## Hooks au niveau de l’application

Il s’agit de fonctions de rappel pratiques dans `application.Options` qui vous permettent d’intervenir dans le cycle de vie de l’application sans créer un service complet. Elles sont utiles pour les tâches de nettoyage simples, la confirmation de la fermeture de l’application ou l’exécution de code à des étapes précises de la séquence d’arrêt.

Pour une gestion plus complexe du cycle de vie, avec une logique de démarrage, l’injection de dépendances ou des ressources avec état, utilisez plutôt les [services](#cycle-de-vie-des-services).

### ShouldQuit

La fonction de rappel `ShouldQuit` est appelée chaque fois que la fermeture de l’application est demandée, que l’utilisateur ferme la dernière fenêtre, appuie sur Cmd+Q (macOS) ou Alt+F4 (Windows), ou qu’un appel à `app.Quit()` est effectué par programmation.

**Valeur de retour :**

- Renvoyez `true` pour autoriser la fermeture (l’application s’arrêtera).
- Renvoyez `false` pour annuler la fermeture (l’application continuera de s’exécuter).

Vous pouvez ainsi intercepter les demandes de fermeture et éventuellement les empêcher, par exemple pour avertir l’utilisateur de modifications non enregistrées :

```go
app := application.New(application.Options{
    ShouldQuit: func() bool {
        if !hasUnsavedChanges() {
            return true // No unsaved changes, allow quit
        }

        // Prompt the user — MessageDialog.Show() blocks and returns nothing.
        // The button's OnClick callback fires for whichever button the user picks.
        shouldQuit := false

        dlg := application.Get().Dialog.Question().
            SetTitle("Unsaved Changes").
            SetMessage("You have unsaved changes. Quit anyway?")

        quit := dlg.AddButton("Quit")
        cancel := dlg.AddButton("Cancel")
        dlg.SetDefaultButton(cancel)
        dlg.SetCancelButton(cancel)

        quit.OnClick(func() { shouldQuit = true })

        dlg.Show()
        return shouldQuit
    },
})
```

Si `ShouldQuit` n’est pas défini, l’application se ferme immédiatement dès que cela lui est demandé.

**Cas dans lesquels ShouldQuit est appelé :**

- L’utilisateur ferme la dernière fenêtre (sauf si `DisableQuitOnLastWindowClosed` est défini).
- L’utilisateur appuie sur Cmd+Q sous macOS.
- L’utilisateur appuie sur Alt+F4 sous Windows (lorsque la dernière fenêtre a le focus).
- Le code appelle `app.Quit()`.

**Cas dans lesquels ShouldQuit n’est PAS appelé :**

- Le processus est arrêté de force (SIGKILL ou arrêt forcé depuis le Gestionnaire des tâches).
- `os.Exit()` est appelé directement.

### OnShutdown

La fonction de rappel `OnShutdown` est appelée lorsque la fermeture de l’application est confirmée (après que `ShouldQuit` a renvoyé `true`, si cette fonction est définie). Utilisez-la pour effectuer des tâches de nettoyage, comme enregistrer l’état, fermer les connexions aux bases de données ou libérer les ressources.

```go
app := application.New(application.Options{
    OnShutdown: func() {
        // Save application state
        saveState()

        // Close connections
        cleanup()
    },
})
```

Vous pouvez également enregistrer par programmation des fonctions de rappel d’arrêt supplémentaires à tout moment du cycle de vie de l’application :

```go
app.OnShutdown(func() {
    log.Println("Application shutting down...")
})
```

Les fonctions de rappel sont exécutées dans leur ordre d’enregistrement. Le processus d’arrêt reste bloqué jusqu’à ce qu’elles soient toutes terminées.

**Important :** veillez à ce que les fonctions de rappel d’arrêt s’exécutent rapidement (en moins de 1 seconde). Le système d’exploitation peut forcer l’arrêt des applications qui mettent trop de temps à se fermer, ce qui pourrait interrompre vos tâches de nettoyage et entraîner une perte de données.

### PostShutdown

La fonction de rappel `PostShutdown` est appelée une fois toutes les tâches d’arrêt terminées, juste avant la fin du processus. À ce stade, l’instance de l’application n’est plus utilisable : les fenêtres sont fermées, les services sont arrêtés et les ressources sont libérées.

Elle est principalement utile pour :

- La journalisation finale qui doit avoir lieu après toutes les autres opérations de nettoyage
- Le test et le débogage du comportement à l’arrêt
- Les plateformes sur lesquelles `app.Run()` ne rend pas la main (la fonction de rappel garantit l’exécution de votre code)

```go
app := application.New(application.Options{
    PostShutdown: func() {
        // Final logging
        log.Println("Application terminated cleanly")

        // Flush any buffered logs
        logger.Sync()
    },
})
```

**Remarque :** n’essayez pas d’utiliser les fonctionnalités de l’application (fenêtres, boîtes de dialogue, etc.) dans `PostShutdown`, car elles ne sont plus disponibles.

## Cycle de vie basé sur les événements

Wails fournit un système d’événements qui vous avertit lorsque des événements surviennent dans votre application : ouverture de fenêtres, démarrage de l’application, changements de thème, etc. Vous pouvez écouter ces événements pour réagir aux changements du cycle de vie sans les bloquer ni les intercepter.

Pour les événements de fenêtre, vous pouvez également utiliser `RegisterHook` plutôt que `OnWindowEvent` afin d’intercepter et d’annuler des actions, par exemple pour empêcher la fermeture d’une fenêtre. Consultez la section [Hooks de fenêtre](#hooks-de-fentre-vnements-annulables) ci-dessous.

Pour obtenir la documentation complète du système d’événements, consultez le [guide des événements](/features/events/system/).

### Événements de l’application

Écoutez les événements du cycle de vie de l’application :

```go
app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(event *application.ApplicationEvent) {
    app.Logger.Info("Application has started!")
})
```

Des événements propres à chaque plateforme sont également disponibles :

```go
// macOS
app.Event.OnApplicationEvent(events.Mac.ApplicationDidFinishLaunching, func(event *application.ApplicationEvent) {
    // Handle macOS launch
})

app.Event.OnApplicationEvent(events.Mac.ApplicationWillTerminate, func(event *application.ApplicationEvent) {
    // Handle macOS termination
})

// Windows
app.Event.OnApplicationEvent(events.Windows.ApplicationStarted, func(event *application.ApplicationEvent) {
    // Handle Windows start
})
```

### Événements de fenêtre

Écoutez les événements du cycle de vie des fenêtres :

```go
window := app.Window.New()

window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window gained focus")
})

window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
    app.Logger.Info("Window is closing")
})
```

### Hooks de fenêtre (événements annulables)

Utilisez `RegisterHook` plutôt que `OnWindowEvent` lorsque vous devez **annuler** un événement :

```go
window := app.Window.New()

var countdown = 3

window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    countdown--
    if countdown > 0 {
        app.Logger.Info("Not closing yet!", "remaining", countdown)
        e.Cancel()  // Prevent the window from closing
        return
    }
    app.Logger.Info("Window closing now")
})
```

**Différence entre OnWindowEvent et RegisterHook :**

- `OnWindowEvent` : vous avertit lorsqu’un événement se produit (sans possibilité de l’annuler)
- `RegisterHook` : vous permet d’intercepter l’événement et, éventuellement, de l’annuler

## Cycle de vie des fenêtres

Les fenêtres ont leur propre cycle de vie, de leur création à leur destruction. Chaque fenêtre charge son contenu frontend indépendamment et peut être affichée, masquée ou fermée à tout moment. Lorsqu’un utilisateur tente de fermer une fenêtre, vous pouvez intercepter cette action avec un `RegisterHook` afin de demander confirmation ou de masquer la fenêtre au lieu de la détruire.

Pour consulter la documentation complète sur les fenêtres, reportez-vous au [guide des fenêtres](/features/windows/basics/).

```d2
direction: down

Create: Créer la fenêtre {
  shape: oval
  style.fill: "#10B981"
}

Load: Charger l’interface {
  shape: rectangle
}

Show: Afficher la fenêtre {
  shape: rectangle
}

Active: Fenêtre active {
  Events: Gérer les événements {
    shape: rectangle
  }
}

CloseRequest: Demande de fermeture {
  shape: diamond
  style.fill: "#F59E0B"
}

Hook: Hook WindowClosing {
  shape: rectangle
  style.fill: "#3B82F6"
}

Destroy: Détruire la fenêtre {
  shape: rectangle
}

End: Fenêtre fermée {
  shape: oval
  style.fill: "#EF4444"
}

Create -> Load
Load -> Show
Show -> Active.Events
Active.Events -> Active.Events: Boucle
Active.Events -> CloseRequest: L’utilisateur ferme la fenêtre
CloseRequest -> Hook
Hook -> Active.Events: Annulée
Hook -> Destroy: Autorisée
Destroy -> End
```

### Création de fenêtres

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "My Window",
    Width:  800,
    Height: 600,
})
```

### Empêcher la fermeture d’une fenêtre

```go
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    if !hasUnsavedChanges() {
        return
    }

    // MessageDialog.Show() returns nothing; per-button OnClick handlers fire.
    dlg := application.Get().Dialog.Question().
        SetTitle("Unsaved Changes").
        SetMessage("Save before closing?")

    save := dlg.AddButton("Save")
    discard := dlg.AddButton("Discard")
    cancel := dlg.AddButton("Cancel")
    dlg.SetDefaultButton(save)
    dlg.SetCancelButton(cancel)

    save.OnClick(func() { saveChanges() })
    cancel.OnClick(func() { e.Cancel() }) // Prevent close
    _ = discard                            // "Discard" falls through and allows close

    dlg.Show()
})
```

### Masquer au lieu de fermer

Une approche courante pour les applications de la zone de notification :

```go
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    window.Hide()  // Hide instead of destroy
    e.Cancel()     // Prevent actual close
})
```

## Cycle de vie avec plusieurs fenêtres

Avec plusieurs fenêtres :

```go
mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Main Window",
})

settingsWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "Settings",
    Width:  400,
    Height: 600,
    Hidden: true,  // Start hidden
})
```

**Le comportement par défaut varie selon la plateforme :**

| Plateforme | Comportement par défaut à la fermeture de la dernière fenêtre |
| --- | --- |
| macOS | L’application continue de s’exécuter (la barre des menus reste affichée) |
| Windows | L’application se ferme |
| Linux | L’application se ferme |

macOS suit les conventions natives de la plateforme, selon lesquelles les applications restent généralement actives dans la barre des menus même lorsqu’aucune fenêtre n’est ouverte. Sous Windows et Linux, elles se ferment par défaut.

**Fermer l’application sur toutes les plateformes lorsque la dernière fenêtre se ferme :**

```go
app := application.New(application.Options{
    Mac: application.MacOptions{
        ApplicationShouldTerminateAfterLastWindowClosed: true,
    },
})
```

**Maintenir l’application en cours d’exécution sur toutes les plateformes lorsque la dernière fenêtre se ferme :**

Ce comportement est utile pour les applications de la zone de notification ou celles qui doivent continuer de s’exécuter en arrière-plan.

```go
app := application.New(application.Options{
    Windows: application.WindowsOptions{
        DisableQuitOnLastWindowClosed: true,
    },
    Linux: application.LinuxOptions{
        DisableQuitOnLastWindowClosed: true,
    },
})
```

## Approches courantes

### Approche 1 : service de base de données

```go
type DatabaseService struct {
    db *sql.DB
}

func (s *DatabaseService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    var err error
    s.db, err = sql.Open("sqlite3", "app.db")
    if err != nil {
        return fmt.Errorf("failed to open database: %w", err)
    }

    if err := s.db.PingContext(ctx); err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }

    return nil
}

func (s *DatabaseService) ServiceShutdown() error {
    if s.db != nil {
        return s.db.Close()
    }
    return nil
}

// Exported methods are available to the frontend
func (s *DatabaseService) GetUsers() ([]User, error) {
    // Query implementation
}
```

### Approche 2 : service de configuration

```go
type ConfigService struct {
    config *Config
    path   string
}

func (s *ConfigService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    s.path = "config.json"

    data, err := os.ReadFile(s.path)
    if err != nil {
        if os.IsNotExist(err) {
            s.config = &Config{} // Default config
            return nil
        }
        return err
    }

    return json.Unmarshal(data, &s.config)
}

func (s *ConfigService) ServiceShutdown() error {
    data, err := json.MarshalIndent(s.config, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(s.path, data, 0644)
}
```

### Approche 3 : tâche d’arrière-plan

```go
type WorkerService struct {
    cancel context.CancelFunc
}

func (s *WorkerService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    workerCtx, cancel := context.WithCancel(ctx)
    s.cancel = cancel

    go s.runWorker(workerCtx)

    return nil
}

func (s *WorkerService) ServiceShutdown() error {
    if s.cancel != nil {
        s.cancel()
    }
    return nil
}

func (s *WorkerService) runWorker(ctx context.Context) {
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            s.doWork()
        case <-ctx.Done():
            return
        }
    }
}
```

## Référence du cycle de vie

| Hook/interface | Moment de l’appel | Annulation possible ? | Utilisation |
| --- | --- | --- | --- |
| `ServiceStartup` | Pendant `app.Run()`, avant la boucle d’événements | Non (renvoyez une erreur pour interrompre) | Initialisation |
| `ServiceShutdown` | Pendant l’arrêt, après `OnShutdown` | Non | Nettoyage |
| `OnShutdown` | Lorsque la fermeture de l’application est confirmée | Non | Nettoyage de l’application |
| `ShouldQuit` | Lorsqu’une fermeture de l’application est demandée | Oui (renvoyez false) | Confirmer la fermeture de l’application |
| `RegisterHook(WindowClosing)` | Lorsqu’une fermeture de fenêtre est demandée | Oui (`e.Cancel()`) | Empêcher la fermeture de la fenêtre |
| `OnWindowEvent` | Lorsqu’un événement se produit | Non | Réagir aux événements |
| `OnApplicationEvent` | Lorsqu’un événement se produit | Non | Réagir aux événements |

## Différences entre les plateformes

### macOS

- Le **menu de l’application** reste disponible même lorsqu’aucune fenêtre n’est ouverte
- **Cmd+Q** déclenche la fermeture de l’application (en passant par `ShouldQuit`)
- L’**icône du Dock** reste affichée, sauf si elle est masquée
- Utilisez `ApplicationShouldTerminateAfterLastWindowClosed` pour contrôler le comportement lors de la fermeture de l’application

### Windows

- Aucun **menu d’application** sans fenêtre
- **Alt+F4** ferme la fenêtre (cette fermeture peut être empêchée avec `RegisterHook`)
- La **zone de notification** peut maintenir l’application en cours d’exécution

### Linux

- Le **comportement varie** selon l’environnement de bureau
- **Généralement similaire à Windows**

## Débogage des problèmes de cycle de vie

### Problème : l’application ne se ferme pas

**Causes :**

1. `ShouldQuit` renvoie `false`
2. `OnShutdown` prend trop de temps
3. Les goroutines en arrière-plan ne s’arrêtent pas

**Solution :**

```go
// 1. Check ShouldQuit logic
ShouldQuit: func() bool {
    log.Println("ShouldQuit called")
    return true
}

// 2. Keep OnShutdown fast
OnShutdown: func() {
    log.Println("OnShutdown started")
    // Fast cleanup only
    log.Println("OnShutdown finished")
}

// 3. Use context for background tasks
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    go func() {
        <-ctx.Done()
        log.Println("Context cancelled, stopping background work")
    }()
    return nil
}
```

### Problème : le démarrage du service échoue

**Solution :** renvoyez des erreurs descriptives :

```go
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    if err := s.init(); err != nil {
        return fmt.Errorf("failed to initialise: %w", err)
    }
    return nil
}
```

L’erreur sera consignée et l’application ne démarrera pas.

## Bonnes pratiques

### À faire

- **Utilisez des services pour gérer le cycle de vie** : ils fournissent les hooks nécessaires au démarrage et à l’arrêt
- **Assurez un arrêt rapide** : visez moins de 1 seconde pour l’ensemble du nettoyage
- **Utilisez le contexte pour l’annulation** : arrêtez correctement les tâches en arrière-plan
- **Gérez les erreurs au démarrage** : renvoyez les erreurs pour interrompre proprement le démarrage
- **Consignez les événements du cycle de vie** : cela facilite le débogage

### À ne pas faire

- **Ne bloquez pas le démarrage d’un service** : veillez à ce que l’initialisation soit rapide (moins de 2 secondes)
- **N’affichez pas de boîtes de dialogue pendant l’arrêt** : l’application est en cours de fermeture et l’interface utilisateur risque de ne pas fonctionner
- **N’ignorez pas le contexte** : vérifiez toujours `ctx.Done()` dans les goroutines
- **Ne laissez pas de ressources non libérées** : implémentez toujours `ServiceShutdown`

## Étapes suivantes

**Services** : approfondissez le système de services [En savoir plus →](/features/bindings/services/)

**Système d’événements** : utilisez les événements pour communiquer [En savoir plus →](/features/events/system/)

**Gestion des fenêtres** : créez et gérez des fenêtres [En savoir plus →](/features/windows/basics/)

---

**Des questions sur le cycle de vie ?** Posez-les sur [Discord](https://discord.gg/JDdSxwjhGf) ou consultez les [exemples](https://github.com/wailsapp/wails/tree/master/v3/examples).
