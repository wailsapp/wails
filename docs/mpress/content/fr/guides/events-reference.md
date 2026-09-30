---
title: "Guide des événements"
description: "Guide pratique de l’utilisation des événements dans Wails v3 pour la communication au sein de l’application et la gestion de son cycle de vie"
slug: "guides/events-reference"
sourcePath: "guides/events-reference.md"
---

**REMARQUE : ce guide est en cours de rédaction**

## Guide des événements

Les événements sont au cœur de la communication dans les applications Wails. Ils permettent aux différentes parties de votre application de communiquer entre elles sans être étroitement couplées. Ce guide vous présente tout ce que vous devez savoir pour utiliser efficacement les événements dans votre application Wails.

## Comprendre les événements Wails

Considérez les événements comme des messages diffusés dans toute votre application. Chaque partie de votre application peut écouter ces messages et réagir en conséquence. Cela s’avère particulièrement utile pour :

- **Réagir aux changements de la fenêtre** : savoir quand votre fenêtre est réduite, agrandie ou déplacée
- **Gérer les événements système** : réagir aux changements de thème ou aux événements d’alimentation
- **Implémenter une logique propre à l’application** : créer vos propres événements pour des fonctionnalités telles que les mises à jour de données ou les actions utilisateur
- **Faire communiquer les composants** : permettre aux différentes parties de votre application de communiquer sans dépendances directes

## Convention de nommage des événements

Tous les événements Wails suivent une convention d’espace de noms qui indique clairement leur origine :

- `common:` — Événements multiplateformes fonctionnant sous Windows, macOS et Linux
- `windows:` — Événements propres à Windows
- `mac:` — Événements propres à macOS\
- `linux:` — Événements propres à Linux

Par exemple :

- `common:WindowFocus` — La fenêtre a obtenu le focus (fonctionne sur toutes les plateformes)
- `windows:APMSuspend` — Le système entre en veille (Windows uniquement)
- `mac:ApplicationDidBecomeActive` — L’application est devenue active (macOS uniquement)

## Premiers pas avec les événements

### Écouter les événements (frontend)

Le cas d’usage le plus courant consiste à écouter les événements dans le code de votre frontend :

```javascript
import { Events } from '@wailsio/runtime';

// Listen for when the window gains focus
Events.On('common:WindowFocus', () => {
    console.log('Window is now focused!');
    // Maybe refresh some data or resume animations
});

// Listen for theme changes
Events.On('common:ThemeChanged', (event) => {
    console.log('Theme changed:', event.data);
    // Update your app's theme accordingly
});

// Listen for custom events from your Go backend
Events.On('my-app:data-updated', (event) => {
    console.log('Data updated:', event.data);
    // Update your UI with the new data
});
```

### Émettre des événements (backend)

Votre code Go peut émettre des événements que votre frontend peut écouter :

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "time"
)

func (s *Service) UpdateData() {
    // Do some data processing...

    // Notify the frontend
    app := application.Get()
    app.Event.Emit("my-app:data-updated",
        map[string]interface{}{
            "timestamp": time.Now(),
            "count": 42,
        },
	)
}
```

### Émettre des événements (frontend)

Bien que cet usage soit moins courant, votre frontend peut également émettre des événements que votre code Go peut écouter :

```javascript
import { Events } from '@wailsio/runtime';

// Event without data
Events.Emit('myapp:close-window')

// Event with data
Events.Emit('myapp:disconnect-requested', 'id-123')
```

Si vous utilisez TypeScript dans votre frontend et que vous [enregistrez des événements typés](#vnements-typs-avec-sret-des-types) dans votre code Go, vous bénéficierez de l’autocomplétion et de la vérification des noms d’événements, ainsi que de la vérification des types de données.

### Supprimer les écouteurs d’événements

Supprimez toujours vos écouteurs d’événements lorsqu’ils ne sont plus nécessaires :

```javascript
import { Events } from '@wailsio/runtime';

// Store the handler reference
const focusHandler = () => {
    console.log('Window focused');
};

// Add the listener — capture the unsubscribe function it returns
const unsubscribe = Events.On('common:WindowFocus', focusHandler);

// Later, remove this specific listener via the returned unsubscribe
unsubscribe();

// Or remove ALL listeners for one (or more) event names — Events.Off takes only event-name strings
Events.Off('common:WindowFocus');
// Events.Off('common:WindowFocus', 'common:WindowLostFocus'); // variadic
// Events.OffAll(); // remove every listener for every event (no args)
```

## Cas d’usage courants

### 1. Suspendre et reprendre selon le focus de la fenêtre

De nombreuses applications doivent suspendre certaines activités lorsque la fenêtre perd le focus :

```javascript
import { Events } from '@wailsio/runtime';

let animationRunning = true;

Events.On('common:WindowLostFocus', () => {
    animationRunning = false;
    pauseBackgroundTasks();
});

Events.On('common:WindowFocus', () => {
    animationRunning = true;
    resumeBackgroundTasks();
});
```

### 2. Réagir aux changements de thème

Synchronisez votre application avec le thème du système :

```javascript
import { Events } from '@wailsio/runtime';

Events.On('common:ThemeChanged', (event) => {
    const isDarkMode = event.data.isDarkMode;

    if (isDarkMode) {
        document.body.classList.add('dark-theme');
        document.body.classList.remove('light-theme');
    } else {
        document.body.classList.add('light-theme');
        document.body.classList.remove('dark-theme');
    }
});
```

### 3. Gérer le dépôt de fichiers

Permettez à votre application d’accepter les fichiers déposés par glisser-déposer :

```javascript
import { Events } from '@wailsio/runtime';

Events.On('common:WindowFilesDropped', (event) => {
    const files = event.data.files;

    files.forEach(file => {
        console.log('File dropped:', file);
        // Process the dropped files
        handleFileUpload(file);
    });
});
```

### 4. Gérer le cycle de vie de la fenêtre

Réagissez aux changements d’état de la fenêtre :

```javascript
import { Events } from '@wailsio/runtime';

Events.On('common:WindowClosing', () => {
    // Save user data before closing
    saveApplicationState();

    // You could also prevent closing by returning false
    // from a registered window close handler
});

Events.On('common:WindowMaximise', () => {
    // Adjust UI for maximized view
    adjustLayoutForMaximized();
});

Events.On('common:WindowRestore', () => {
    // Return UI to normal state
    adjustLayoutForNormal();
});
```

### 5. Fonctionnalités propres à chaque plateforme

Gérez les événements propres à chaque plateforme lorsque nécessaire :

```javascript
import { Events } from '@wailsio/runtime';

// Windows-specific power management
Events.On('windows:APMSuspend', () => {
    console.log('System is going to sleep');
    saveState();
});

Events.On('windows:APMResumeSuspend', () => {
    console.log('System woke up');
    refreshData();
});

// macOS-specific app lifecycle
Events.On('mac:ApplicationWillTerminate', () => {
    console.log('App is about to quit');
    performCleanup();
});
```

## Créer des événements personnalisés

Vous pouvez créer vos propres événements pour répondre aux besoins spécifiques de votre application.

### Backend (Go)

```go
// Emit a custom event when data changes

func (s *Service) ProcessUserData(userData UserData) error {
    // Process the data...

    app := application.Get()
    // Notify all listeners
    app.Event.Emit("user:data-processed",
        map[string]interface{}{
            "userId": userData.ID,
            "status": "completed",
            "timestamp": time.Now(),
        },
    )
    return nil
}

// Emit periodic updates
func (s *Service) StartMonitoring() {
    app := application.Get()
    ticker := time.NewTicker(5 * time.Second)
    go func() {
        for range ticker.C {
            stats := s.collectStats()
            app.Event.Emit("monitor:stats-updated", stats)
        }
    }()
}
```

### Frontend (JavaScript)

```javascript
import { Events } from '@wailsio/runtime';

// Listen for your custom events
Events.On('user:data-processed', (event) => {
    const { userId, status, timestamp } = event.data;

    showNotification(`User ${userId} processing ${status}`);
    updateUIWithNewData();
});

Events.On('monitor:stats-updated', (event) => {
    updateDashboard(event.data);
});
```

## Événements typés avec sûreté des types

Wails v3 prend en charge les événements typés avec une sûreté complète des types TypeScript grâce à l’enregistrement des événements et à la génération automatique des liaisons.

### Enregistrer des événements personnalisés

Appelez `application.RegisterEvent` lors de l’initialisation pour enregistrer les noms d’événements personnalisés avec leurs types de données :

```go
package main

import "github.com/wailsapp/wails/v3/pkg/application"

type UserLoginData struct {
    UserID   string
    Username string
    LoginTime string
}

type MonitorStats struct {
    CPUUsage    float64
    MemoryUsage float64
}

func init() {
    // Register events with their data types
    application.RegisterEvent[UserLoginData]("user:login")
    application.RegisterEvent[MonitorStats]("monitor:stats")

    // Register events without data (void events)
    application.RegisterEvent[application.Void]("app:ready")
}
```

@note{type="caution"}
`RegisterEvent` est destiné à être appelé lors de l’initialisation et provoquera une panique si :

- Les arguments ne sont pas valides
- Le même nom d’événement est enregistré deux fois avec des types de données différents

@end

@note{type="info"}
Vous pouvez enregistrer sans risque le même événement plusieurs fois, à condition que son type de données soit toujours identique. Cela peut être utile pour garantir l’enregistrement d’un événement lorsque l’un quelconque de plusieurs paquets est chargé.

@end

### Avantages de l’enregistrement des événements

Après l’enregistrement, le type des arguments de données transmis à `Event.Emit` est vérifié par rapport au type spécifié. En cas d’incompatibilité :

- Une erreur est émise et consignée dans les journaux (ou transmise au gestionnaire d’erreurs enregistré)
- L’événement concerné ne sera pas propagé
- Cela garantit que le champ de données des événements enregistrés est toujours assignable au type déclaré

### Mode strict

Utilisez le tag de build `strictevents` pour activer, pendant le développement, les avertissements concernant les événements non enregistrés :

```bash
go build -tags strictevents
```

Lorsque le mode strict est activé, l’environnement d’exécution émet au maximum un avertissement par nom d’événement non enregistré afin de ne pas saturer les journaux.

### Génération des liaisons TypeScript

Le générateur de liaisons produit des définitions TypeScript et du code d’intégration assurant une prise en charge transparente des événements typés dans le frontend.

#### 1. Configurez le plugin Vite

Dans votre fichier `vite.config.ts` :

```typescript
import { defineConfig } from 'vite'
import wails from '@wailsio/runtime/plugins/vite'

export default defineConfig({
  plugins: [wails()],
})
```

#### 2. Générez les liaisons

Exécutez le générateur de liaisons :

```bash
wails3 generate bindings
```

Cette commande crée, dans le répertoire de votre frontend, des fichiers TypeScript contenant des fonctions de création d’événements typés et des interfaces de données.

#### 3. Utilisez les événements typés dans le frontend

```typescript
import { Events } from '@wailsio/runtime'
import { UserLogin, MonitorStats } from './bindings/events'

// Type-safe event emission with autocomplete
Events.Emit(UserLogin({
    UserID: "123",
    Username: "john_doe",
    LoginTime: new Date().toISOString()
}))

// Type-safe event listening
Events.On(UserLogin, (event) => {
    // event.data is typed as UserLoginData
    console.log(`User ${event.data.Username} logged in`)
})

Events.On(MonitorStats, (event) => {
    // event.data is typed as MonitorStats
    updateDashboard({
        cpu: event.data.CPUUsage,
        memory: event.data.MemoryUsage
    })
})
```

Les événements typés offrent :

- **Autocomplétion** des noms d’événements
- **Vérification des types** des données d’événement
- **Erreurs à la compilation** en cas d’incompatibilité des types de données
- Documentation **IntelliSense**

## Référence des événements

### Événements communs (multiplateformes)

Ces événements fonctionnent sur toutes les plateformes :

| Événement | Description | Quand l’utiliser |
| --- | --- | --- |
| `common:ApplicationStarted` | L’application a complètement démarré | Initialiser votre application et charger l’état enregistré |
| `common:WindowRuntimeReady` | L’environnement d’exécution Wails est prêt | Commencer à appeler l’API Wails |
| `common:ThemeChanged` | Le thème du système a changé | Mettre à jour l’apparence de l’application |
| `common:SystemWillSleep` | Le système est sur le point d’entrer en veille | Enregistrer l’état sur le stockage et fermer les sockets |
| `common:SystemDidWake` | Le système est sorti de veille | Se reconnecter et actualiser les données obsolètes |
| `common:WindowFocus` | La fenêtre a obtenu le focus | Reprendre les activités et actualiser les données |
| `common:WindowLostFocus` | La fenêtre a perdu le focus | Suspendre les activités et enregistrer l’état |
| `common:WindowMinimise` | La fenêtre a été réduite | Suspendre le rendu et réduire l’utilisation des ressources |
| `common:WindowMaximise` | La fenêtre a été agrandie | Adapter la mise en page au plein écran |
| `common:WindowRestore` | La fenêtre a été restaurée après avoir été réduite ou agrandie | Rétablir la mise en page normale |
| `common:WindowClosing` | La fenêtre est sur le point de se fermer | Enregistrer les données et libérer les ressources |
| `common:WindowFilesDropped` | Des fichiers ont été déposés sur la fenêtre | Traiter les importations de fichiers |
| `common:WindowDidResize` | La fenêtre a été redimensionnée | Adapter la mise en page et restituer à nouveau les graphiques |
| `common:WindowDidMove` | La fenêtre a été déplacée | Mettre à jour les fonctionnalités dépendant de la position |

### Événements propres à chaque plateforme

#### Événements Windows

Principaux événements pour les applications Windows :

| Événement | Description | Cas d’utilisation |
| --- | --- | --- |
| `windows:SystemThemeChanged` | Le thème Windows a changé | Mettre à jour les couleurs de l’application |
| `windows:APMSuspend` | Le système passe en veille | Enregistrer l’état et suspendre les opérations |
| `windows:APMResumeAutomatic` | Le système est sorti de veille (toujours déclenché à la sortie de veille) | Restaurer l’état et actualiser les données |
| `windows:APMResumeSuspend` | Le système est sorti de veille à la suite d’une saisie utilisateur (après `APMResumeAutomatic`) | Distinguer une sortie de veille déclenchée par l’utilisateur |
| `windows:APMPowerStatusChange` | L’état de l’alimentation a changé | Adapter les paramètres de performances |

#### Événements macOS

Événements importants des applications macOS :

| Événement | Description | Cas d’utilisation |
| --- | --- | --- |
| `mac:ApplicationDidBecomeActive` | L’application est devenue active | Reprendre les opérations |
| `mac:ApplicationDidResignActive` | L’application est devenue inactive | Suspendre les opérations |
| `mac:ApplicationWillTerminate` | L’application va se fermer | Effectuer le nettoyage final |
| `mac:ApplicationWillSleep` | Le système est sur le point de passer en veille | Enregistrer l’état et fermer les sockets |
| `mac:ApplicationDidWake` | Le système est sorti de veille | Se reconnecter et actualiser |
| `mac:ApplicationScreensDidSleep` | Les écrans se sont mis en veille | Suspendre le rendu (indépendamment de la mise en veille du système) |
| `mac:ApplicationScreensDidWake` | Les écrans sont sortis de veille | Reprendre le rendu |
| `mac:WindowDidEnterFullScreen` | Passage en plein écran | Adapter l’interface utilisateur au mode plein écran |
| `mac:WindowDidExitFullScreen` | Sortie du mode plein écran | Rétablir l’interface utilisateur normale |

#### Événements Linux

Principaux événements de fenêtre sous Linux :

| Événement | Description | Cas d’utilisation |
| --- | --- | --- |
| `linux:SystemThemeChanged` | Le thème du bureau a changé | Mettre à jour le thème de l’application |
| `linux:SystemWillSleep` | Le système est sur le point de passer en veille (logind) | Enregistrer l’état |
| `linux:SystemDidWake` | Le système est sorti de veille (logind) | Se reconnecter et actualiser |
| `linux:WindowFocusIn` | La fenêtre a obtenu le focus | Reprendre les activités |
| `linux:WindowFocusOut` | La fenêtre a perdu le focus | Suspendre les activités |
| `linux:WindowLoadStarted` | La WebView a commencé le chargement | Afficher l’indicateur de chargement |
| `linux:WindowLoadRedirected` | La WebView a été redirigée | Suivre les redirections de navigation |
| `linux:WindowLoadCommitted` | La WebView a validé le chargement | Le contenu est en cours de réception |
| `linux:WindowLoadFinished` | La WebView a terminé le chargement | Masquer l’indicateur de chargement, injecter du JS/CSS |

## Bonnes pratiques

### 1. Utiliser des espaces de noms pour les événements

Lors de la création d’événements personnalisés, utilisez des espaces de noms pour éviter les conflits :

```javascript
import { Events } from '@wailsio/runtime';

// Good - namespaced events
Events.Emit('myapp:user:login');
Events.Emit('myapp:data:updated');
Events.Emit('myapp:network:connected');

// Avoid - generic names that might conflict
Events.Emit('login');
Events.Emit('update');
```

### 2. Supprimer les écouteurs inutilisés

Supprimez toujours les écouteurs d’événements lorsque les composants sont démontés :

```javascript
import { Events } from '@wailsio/runtime';

// React example
useEffect(() => {
    const handler = (event) => {
        // Handle event
    };

    const off = Events.On('common:WindowDidResize', handler);

    // Cleanup — call the unsubscribe returned by Events.On
    return () => {
        off();
    };
}, []);
```

### 3. Gérer les différences entre les plateformes

Lorsque vous utilisez des événements propres à une plateforme, vérifiez leur disponibilité sur celle-ci :

```javascript
import { Events } from '@wailsio/runtime';

// Platform-specific events can be registered unconditionally;
// they will simply never fire on unsupported platforms.
Events.On('windows:APMSuspend', handleSuspend);
Events.On('mac:ApplicationWillTerminate', handleTerminate);
```

### 4. Ne pas abuser des événements

Bien que les événements soient puissants, ne les utilisez pas pour tout :

- ✅ Utilisez les événements pour : les notifications système, les changements de cycle de vie et la diffusion de mises à jour
- ❌ Évitez les événements pour : les valeurs de retour directes des fonctions, les mises à jour d’un seul composant et les opérations synchrones

## Débogage des événements

Pour déboguer les problèmes liés aux événements :

```javascript
import { Events } from '@wailsio/runtime';

// Log all events (development only)
if (isDevelopment) {
    const originalOn = Events.On;
    Events.On = function(eventName, handler) {
        console.log(`[Event Registered] ${eventName}`);
        return originalOn.call(this, eventName, function(event) {
            console.log(`[Event Fired] ${eventName}`, event);
            return handler(event);
        });
    };
}
```

## Source de référence

La liste complète des événements disponibles se trouve dans le code source de Wails :

- Événements frontend : [`v3/internal/runtime/desktop/@wailsio/runtime/src/event_types.ts`](https://github.com/wailsapp/wails/blob/main/v3/internal/runtime/desktop/@wailsio/runtime/src/event_types.ts)
- Événements backend : [`v3/pkg/events/events.go`](https://github.com/wailsapp/wails/blob/main/v3/pkg/events/events.go)

Consultez toujours ces fichiers pour connaître les noms d’événements et leur disponibilité les plus récents.

## Résumé

Dans Wails, les événements offrent un moyen puissant et découplé de gérer la communication au sein de votre application. En suivant les modèles et les pratiques présentés dans ce guide, vous pouvez créer des applications réactives, adaptées aux différentes plateformes et capables de réagir avec fluidité aux changements du système et aux interactions de l’utilisateur.

À retenir : commencez par les événements communs pour assurer la compatibilité multiplateforme, ajoutez des événements propres à chaque plateforme lorsque cela est nécessaire et supprimez toujours vos écouteurs d’événements afin d’éviter les fuites de mémoire.
