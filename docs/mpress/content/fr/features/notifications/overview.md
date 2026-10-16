---
title: "Notifications"
description: "Afficher des notifications système natives avec des boutons d’action et des champs de saisie"
slug: "features/notifications/overview"
sourcePath: "features/notifications/overview.md"
---

## Introduction

Wails fournit un système de notifications multiplateforme complet pour les applications de bureau. Ce service vous permet d’afficher des notifications système natives et prend en charge :

- Notifications simples avec titre, sous-titre et corps de texte
- Notifications interactives avec boutons d’action et réponses textuelles
- [Catégories de notifications](#notifications-interactives) réutilisables pour les actions
- [Sons](#son-personnalis) personnalisés (par défaut, silencieux ou nommés)
- [Pièces jointes](#pices-jointes) (images sur toutes les plateformes ; contenu audio et vidéo sous macOS)
- [Regroupement des notifications associées](#fils-de-discussion-et-regroupement) par `ThreadID`
- [Priorité](#niveau-dinterruption) définie avec `InterruptionLevel` (`passive` / `active` / `timeSensitive` / `critical`)
- [Envoi planifié](#envoi-planifi) (natif sous macOS ; minuteur intégré au processus sous Windows et Linux)
- [Mise à jour d’une notification en cours](#mise--jour-des-notifications) à l’aide de son ID

La prise en charge de chaque nouveau champ facultatif s’adapte aux capacités de la plateforme lorsque celle-ci ne peut pas le gérer pleinement ; consultez [Considérations propres aux plateformes](#considrations-propres-aux-plateformes) pour connaître la matrice de prise en charge de chaque fonctionnalité.

## Utilisation de base

### Création du service

Commencez par initialiser le service de notifications :

```go
import "github.com/wailsapp/wails/v3/pkg/application"
import "github.com/wailsapp/wails/v3/pkg/services/notifications"

// Create a new notification service
notifier := notifications.New()

//Register the service with the application
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(notifier),
    },
})
```

## Autorisation des notifications

Sous macOS, les notifications nécessitent l’autorisation de l’utilisateur. Demandez cette autorisation et vérifiez-la :

```go
authorized, err := notifier.CheckNotificationAuthorization()
if err != nil {
    // Handle authorization error
}
if authorized {
    // Send notifications
} else {
    // Request authorization
    authorized, err = notifier.RequestNotificationAuthorization()
}
```

Sous Windows et Linux, cette opération renvoie toujours `true`.

## Types de notifications

### Notifications simples

Envoyez aux utilisateurs une notification simple avec un identifiant unique, un titre, un sous-titre facultatif (sous macOS et Linux) et un corps de texte :

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID: "unique-id",
    Title: "New Calendar Invite",
    Subtitle: "From: Jane Doe", // Optional
    Body: "Tap to view the event",
})

```

### Notifications interactives

Envoyez une notification avec des boutons d’action et des champs de saisie. Ces notifications nécessitent l’enregistrement préalable d’une catégorie de notifications :

```go
// Define a unique category id
categoryID := "unique-category-id"

// Define a category with actions
category := notifications.NotificationCategory{
    ID: categoryID,
    Actions: []notifications.NotificationAction{
        {
            ID:    "OPEN", 
            Title: "Open",
        },
        {
            ID:          "ARCHIVE", 
            Title:       "Archive", 
            Destructive: true,  /* macOS-specific */
        },
    },
    HasReplyField:    true,
    ReplyPlaceholder: "message...",
    ReplyButtonTitle: "Reply",
}

// Register the category
notifier.RegisterNotificationCategory(category)

// Send an interactive notification with the actions registered in the provided category
notifier.SendNotificationWithActions(notifications.NotificationOptions{
    ID:         "unique-id",
    Title:      "New Message",
    Subtitle:   "From: Jane Doe",
    Body:       "Are you able to make it?",
    CategoryID: categoryID,
})
```

## Réponses aux notifications

Traitez les interactions de l’utilisateur avec les notifications :

```go
notifier.OnNotificationResponse(func(result notifications.NotificationResult) {
    response := result.Response
    fmt.Printf("Notification %s was actioned with: %s\n", response.ID, response.ActionIdentifier)

    if response.ActionIdentifier == "TEXT_REPLY" {
        fmt.Printf("User replied: %s\n", response.UserText)
    }

    if data, ok := response.UserInfo["sender"].(string); ok {
        fmt.Printf("Original sender: %s\n", data)
    }

    // Emit an event to the frontend
    app.Event.Emit("notification", result.Response)
})
```

## Personnalisation des notifications

### Métadonnées personnalisées

Les notifications simples et interactives peuvent inclure des données personnalisées :

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID: "unique-id",
    Title: "New Calendar Invite",
    Subtitle: "From: Jane Doe", // Optional
    Body: "Tap to view the event",
    Data: map[string]interface{}{
        "sender": "jane.doe@example.com",
        "timestamp": "2025-03-10T15:30:00Z",
    }
})
```

### Son personnalisé

Utilisez `Sound` pour contrôler le son émis lors de la remise d’une notification. Laissez-le défini sur `nil` pour utiliser le son par défaut de la plateforme ; utilisez le réglage `Silent: true` pour désactiver le son ; renseignez le champ `Name` pour lire un son nommé ou fourni avec l’application.

```go
// Silent
notifier.SendNotification(notifications.NotificationOptions{
    ID:    "silent-id",
    Title: "Background sync complete",
    Sound: &notifications.NotificationSound{Silent: true},
})

// Named sound
notifier.SendNotification(notifications.NotificationOptions{
    ID:    "named-id",
    Title: "New message",
    Body:  "Tap to read",
    Sound: &notifications.NotificationSound{Name: "Ping"},
})
```

Résolution de `Name` selon la plateforme :

- **macOS** — `Name` est transmis à `[UNNotificationSound soundNamed:]` ; le fichier audio doit se trouver dans le répertoire `Library/Sounds` du paquet de votre application.
- **Windows** — si `Name` commence déjà par `ms-winsoundevent:` ou `ms-appx:`, il est utilisé tel quel ; sinon, il est encapsulé dans `ms-winsoundevent:` afin d’être utilisé comme nom d’événement toast intégré (consultez la documentation Microsoft relative au schéma `<audio>` des notifications toast).
- **Linux** — transmis comme indication freedesktop `sound-name` ; sa lecture dépend du démon de notifications et du thème sonore actifs.

### Pièces jointes

`Attachments` ajoute des fichiers multimédias à une notification. macOS prend en charge plusieurs pièces jointes de tout type de média ; Windows et Linux prennent en charge la première pièce jointe de type image.

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID:    "image-id",
    Title: "Photo uploaded",
    Body:  "Tap to view",
    Attachments: []notifications.NotificationAttachment{
        {
            ID:   "preview",
            Path: "/absolute/path/to/image.png",
            // Type hint:
            //   macOS:   UTI like "public.png" / "public.audio" (often inferred)
            //   Windows: "hero" | "appLogoOverride" | "inline" (default "inline")
            //   Linux:   ignored
            Type: "inline",
        },
    },
})
```

`Path` doit être un chemin absolu dans le système de fichiers. macOS accepte également les URL `file://`.

#### Joindre un fichier fourni avec votre application

Le système d’exploitation lit la pièce jointe sur le disque lors de la remise de la notification ; `Path` doit donc correspondre à un fichier réel sur la machine de l’utilisateur final. Pour une ressource fournie avec l’application (une icône ou une image incorporée avec `go:embed`), vous ne pouvez pas coder en dur un chemin absolu fixe, car le fichier se trouve dans le binaire et non à un emplacement connu sur le disque. Écrivez-le une fois dans un répertoire accessible en écriture au démarrage, puis transmettez ce chemin :

```go
import (
    _ "embed"
    "os"
    "path/filepath"
)

//go:embed assets/preview.png
var previewPNG []byte

// Materialise the embedded asset to a stable path the OS can read.
previewPath := filepath.Join(os.TempDir(), "myapp-preview.png")
if err := os.WriteFile(previewPath, previewPNG, 0o644); err != nil {
    // handle error
}

notifier.SendNotification(notifications.NotificationOptions{
    ID:    "image-id",
    Title: "Photo uploaded",
    Body:  "Tap to view",
    Attachments: []notifications.NotificationAttachment{
        {ID: "preview", Path: previewPath, Type: "inline"},
    },
})
```

Un fichier fourni par l’utilisateur ou téléchargé possède déjà un chemin réel sur le disque ; vous pouvez donc le transmettre directement à `Path` sans effectuer cette étape. La possibilité de transmettre en mémoire les octets d’une pièce jointe pourra être ajoutée dans une version ultérieure.

### Fils de discussion et regroupement

`ThreadID` regroupe les notifications associées afin que le système d’exploitation puisse les réduire dans le centre de notifications.

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID:       "msg-42",
    Title:    "Jane Doe",
    Body:     "Are you free for lunch?",
    ThreadID: "chat:jane.doe",
})
```

### Niveau d’interruption

`InterruptionLevel` contrôle la priorité des notifications. Utilisez l’une des constantes exportées :

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID:                "alert-id",
    Title:             "Server down",
    Body:              "Investigate immediately",
    InterruptionLevel: notifications.InterruptionLevelTimeSensitive,
})
```

| Constante | Valeur | Signification |
| --- | --- | --- |
| `InterruptionLevelPassive` | `"passive"` | Remise discrète ; n’allume pas l’écran et n’émet pas le son par défaut |
| `InterruptionLevelActive` | `"active"` | Niveau par défaut |
| `InterruptionLevelTimeSensitive` | `"timeSensitive"` | Contourne les modes Concentration / Ne pas déranger lorsque cela est autorisé |
| `InterruptionLevelCritical` | `"critical"` | Contourne les modes Concentration et silencieux ; macOS requiert l’autorisation Critical Alert (sans elle, le comportement est silencieusement dégradé) |

Correspondance selon la plateforme :

- **macOS** — définit `UNNotificationContent.interruptionLevel`. `critical` requiert macOS 12+ et l’autorisation Critical Alert.
- **Windows** — correspond à l’attribut `<toast scenario="...">` de la notification toast.
- **Linux** — correspond à l’indication freedesktop `urgency`.

### Envoi planifié

`Schedule` diffère l’envoi. Définissez exactement l’une des valeurs suivantes : `DelaySeconds` (nombre de secondes à partir de maintenant) ou `At` (secondes Unix, UTC).

```go
// Deliver in 60 seconds
notifier.SendNotification(notifications.NotificationOptions{
    ID:       "reminder-id",
    Title:    "Stand up and stretch",
    Schedule: &notifications.NotificationSchedule{DelaySeconds: 60},
})

// Deliver at an absolute time
notifier.SendNotification(notifications.NotificationOptions{
    ID:       "alarm-id",
    Title:    "Meeting in 5 minutes",
    Schedule: &notifications.NotificationSchedule{At: time.Now().Add(time.Hour).Unix()},
})
```

@note{type="caution" title="Persistance"}
Sous **macOS**, les notifications planifiées utilisent un déclencheur natif et persistent après le redémarrage de l’application. Sous **Windows** et **Linux**, la planification utilise à défaut un minuteur `time.AfterFunc` intégré au processus et la notification est **perdue si l’application se ferme avant son envoi** : ni `wintoast` ni la spécification freedesktop ne fournissent de primitive d’envoi différé.

@end

### Mise à jour des notifications

`UpdateNotification` remplace une notification en cours portant le même `ID` :

```go
notifier.UpdateNotification(notifications.NotificationOptions{
    ID:    "download-id",
    Title: "Download complete",
    Body:  "report.pdf is ready",
})
```

Comportement selon la plateforme :

- **macOS** — `UNUserNotificationCenter` déduplique automatiquement les notifications selon leur identifiant ; la notification existante est donc mise à jour sur place.
- **Linux** — utilise le paramètre D-Bus `replaces_id` pour remplacer la notification précédente.
- **Windows** — renvoie actuellement la notification comme une nouvelle notification. Un véritable remplacement sur place nécessite la prise en charge en amont de `tag` / `group` par `wintoast`.

## Considérations propres aux plateformes

@tabs
[macOS]
Sous macOS, les notifications :

- Nécessitent l’autorisation de l’utilisateur
- Nécessitent que l’application soit empaquetée et signée (et notariée pour la distribution)
- Utilisent l’apparence standard des notifications système
- Prennent en charge `Subtitle`
- Prennent en charge la saisie de texte par l’utilisateur (réponses)
- Prennent en charge l’option d’action `Destructive`
- Prennent en charge plusieurs `Attachments` de tout type de média (images, audio et vidéo)
- Prennent en charge `ThreadID` pour le regroupement dans le Centre de notifications
- Prennent en charge toutes les valeurs de `InterruptionLevel` (`critical` requiert l’autorisation Critical Alert)
- Prennent en charge l’envoi planifié natif, qui persiste après le redémarrage de l’application
- Dédupliquent automatiquement les appels à `UpdateNotification` selon `ID`
- Gèrent automatiquement les modes sombre et clair

[Windows]
Sous Windows, les notifications :

- Utilisent les styles de notifications toast du système Windows au moyen du backend `wintoast`
- S’adaptent aux paramètres de thème de Windows
- Prennent en charge la saisie de texte par l’utilisateur (réponses)
- Prennent en charge les écrans à haute résolution (DPI)
- Ne prennent pas en charge `Subtitle`
- Prennent en charge une seule image `Attachment` avec l’indication de placement `hero`, `appLogoOverride` ou `inline` (`inline` par défaut)
- Prennent en charge `ThreadID` pour le regroupement dans le centre de notifications
- Prennent en charge `InterruptionLevel` au moyen de l’attribut `scenario` de la notification toast
- Prennent en charge l’envoi planifié au moyen d’un minuteur intégré au processus : **les notifications planifiées sont perdues si l’application se ferme avant leur envoi**
- `UpdateNotification` renvoie actuellement la notification comme une nouvelle notification (un véritable remplacement sur place reste en attente de la prise en charge en amont de `tag`/`group` par `wintoast`)

[Linux]
Sous Linux, les notifications utilisent l’interface D-Bus `org.freedesktop.Notifications`. Pour que les notifications fonctionnent, un démon de notification compatible **doit être en cours d’exécution**.

@note{type="caution" title="Configuration système requise : démon de notification"}
Un démon de notification compatible avec freedesktop doit être installé et en cours d’exécution. Les choix courants sont :

- **dunst** — léger et hautement configurable (`apt install dunst` / `dnf install dunst`)
- **mako** — natif pour Wayland (`apt install mako-notifier`)
- **GNOME Shell** — enregistre automatiquement l’interface sous GNOME 43+. Sous Ubuntu 24.04 (GNOME Shell 46), il est possible que l’interface ne s’enregistre pas automatiquement au démarrage de la session ; si les notifications ne s’affichent pas, installez `dunst` comme solution de secours.
- **xfce4-notifyd** — fourni avec les environnements de bureau XFCE

Si aucun démon n’est en cours d’exécution, `SendNotification` renvoie une erreur D-Bus : `The name org.freedesktop.Notifications was not provided by any .service files`. Gérez cette erreur dans votre application et indiquez à l’utilisateur d’installer un démon de notification.

@end

Sous Linux, les notifications :

- Suivent le thème de l’environnement de bureau
- Sont positionnées selon les règles de l’environnement de bureau
- Prennent en charge `Subtitle` (concaténé au corps pour les démons qui ne l’affichent pas séparément)
- Ne prennent pas en charge la saisie de texte par l’utilisateur (elle ne fait pas partie de la spécification freedesktop)
- Prennent en charge une seule image `Attachment` au moyen de l’indication `image-path`
- Prennent en charge `ThreadID` (géré par le démon lorsqu’il le prend en charge)
- `Sound.Name` est transmis comme indication `sound-name` ; sa lecture dépend du démon actif et du thème sonore
- Associe `InterruptionLevel` à l’indication freedesktop `urgency`
- Prend en charge la remise planifiée au moyen d’un minuteur interne au processus — **les notifications planifiées sont perdues si l’application se ferme avant leur remise**
- `UpdateNotification` utilise le paramètre D-Bus `replaces_id` pour remplacer sur place la notification précédente

@end

## Bonnes pratiques

1. Vérifiez et demandez l’autorisation :
  - macOS exige l’autorisation de l’utilisateur


2. Proposez des notifications claires et concises :
  - Utilisez des titres, des sous-titres, des textes et des intitulés d’action explicites


3. Traitez correctement les réponses aux notifications :
  - Vérifiez si les réponses aux notifications contiennent des erreurs
  - Fournissez un retour pour les actions de l’utilisateur


4. Tenez compte des conventions de la plateforme :
  - Respectez les modèles de notification propres à la plateforme
  - Respectez les réglages du système


5. Sous Linux, gérez la dépendance au démon :
  - Vérifiez l’erreur renvoyée par `SendNotification` — l’absence de démon produit une erreur D-Bus
  - La documentation du paquet ou le fichier README de l’application devrait préciser qu’un démon de notification freedesktop est requis


## Exemples

Consultez cet exemple :

- [Notifications](https://github.com/wailsapp/wails/tree/master/v3/examples/notifications)

## Référence de l’API

### Gestion du service

| Méthode | Description |
| --- | --- |
| `New()` | Crée un nouveau service de notifications |

### Autorisation des notifications

| Méthode | Description |
| --- | --- |
| `RequestNotificationAuthorization()` | Demande l’autorisation d’afficher des notifications (macOS) |
| `CheckNotificationAuthorization()` | Vérifie l’état actuel de l’autorisation des notifications (macOS) |

### Envoi de notifications

| Méthode | Description |
| --- | --- |
| `SendNotification(options NotificationOptions)` | Envoie une notification simple |
| `SendNotificationWithActions(options NotificationOptions)` | Envoie une notification interactive comportant des actions |
| `UpdateNotification(options NotificationOptions)` | Met à jour une notification en cours à partir de `ID` (voir [Mise à jour des notifications](#mise--jour-des-notifications)) |

### Catégories de notifications

| Méthode | Description |
| --- | --- |
| `RegisterNotificationCategory(category NotificationCategory)` | Enregistre une catégorie de notifications réutilisable |
| `RemoveNotificationCategory(categoryID string)` | Supprime une catégorie précédemment enregistrée |

### Gestion des notifications

| Méthode | Description |
| --- | --- |
| `RemoveAllPendingNotifications()` | Supprime toutes les notifications en attente (macOS et Linux uniquement) |
| `RemovePendingNotification(identifier string)` | Supprime une notification en attente spécifique (macOS et Linux uniquement) |
| `RemoveAllDeliveredNotifications()` | Supprime toutes les notifications remises (macOS et Linux uniquement) |
| `RemoveDeliveredNotification(identifier string)` | Supprime une notification remise spécifique (macOS et Linux uniquement) |
| `RemoveNotification(identifier string)` | Supprime une notification (propre à Linux) |

### Gestion des événements

| Méthode | Description |
| --- | --- |
| `OnNotificationResponse(callback func(result NotificationResult))` | Enregistre une fonction de rappel pour les réponses aux notifications |

### Structures et types

#### NotificationOptions

```go
type NotificationOptions struct {
    ID         string                 // Unique identifier for the notification (required)
    Title      string                 // Main notification title (required)
    Subtitle   string                 // Subtitle text (macOS and Linux only)
    Body       string                 // Main notification content
    CategoryID string                 // Category identifier for interactive notifications
    Data       map[string]interface{} // Custom data to associate with the notification

    // Sound controls the sound played on delivery. nil = platform default.
    Sound *NotificationSound

    // Attachments are media files shown alongside the notification.
    // macOS supports multiple attachments of any media type;
    // Windows and Linux honour the first image-typed attachment.
    Attachments []NotificationAttachment

    // ThreadID groups related notifications together in
    // Notification Center / Action Center / the Linux notification daemon.
    ThreadID string

    // InterruptionLevel controls priority. One of "passive",
    // "active" (default), "timeSensitive", "critical".
    InterruptionLevel string

    // Schedule defers delivery. macOS uses a native trigger that survives
    // restarts; Windows and Linux use an in-process timer that does NOT.
    Schedule *NotificationSchedule
}
```

#### NotificationSound

```go
type NotificationSound struct {
    Silent bool   // If true, no sound is played
    Name   string // Named/bundled sound (see "Custom Sound" above)
}
```

#### NotificationAttachment

```go
type NotificationAttachment struct {
    ID   string // Optional identifier
    Path string // Absolute filesystem path (macOS also accepts file:// URLs)
    // Type is an optional placement/UTI hint:
    //   macOS:   UTI like "public.png" / "public.audio" (often inferred)
    //   Windows: "hero" | "appLogoOverride" | "inline" (default "inline")
    //   Linux:   ignored (always image-path hint)
    Type string
}
```

#### NotificationSchedule

```go
// Exactly one of DelaySeconds or At must be set. At is Unix seconds (UTC).
type NotificationSchedule struct {
    DelaySeconds int   // Seconds from now until delivery
    At           int64 // Absolute Unix timestamp (seconds, UTC)
}
```

#### Constantes InterruptionLevel

```go
const (
    InterruptionLevelPassive       = "passive"
    InterruptionLevelActive        = "active" // default
    InterruptionLevelTimeSensitive = "timeSensitive"
    InterruptionLevelCritical      = "critical"
)
```

#### NotificationCategory

```go
type NotificationCategory struct {
    ID               string                // Unique identifier for the category
    Actions          []NotificationAction  // Button actions for the notification
    HasReplyField    bool                  // Whether to include a text input field
    ReplyPlaceholder string                // Placeholder text for the input field
    ReplyButtonTitle string                // Text for the reply button
}
```

#### NotificationAction

```go
type NotificationAction struct {
    ID          string  // Unique identifier for the action
    Title       string  // Button text
    Destructive bool    // Whether the action is destructive (macOS-specific)
}
```

#### NotificationResponse

```go
type NotificationResponse struct {
    ID               string                  // Notification identifier
    ActionIdentifier string                  // Action that was triggered
    CategoryID       string                  // Category of the notification
    Title            string                  // Title of the notification
    Subtitle         string                  // Subtitle of the notification
    Body             string                  // Body text of the notification
    UserText         string                  // Text entered by the user
    UserInfo         map[string]interface{}  // Custom data from the notification
}
```

#### NotificationResult

```go
type NotificationResult struct {
    Response NotificationResponse  // Response data
    Error    error                 // Any error that occurred
}
```
