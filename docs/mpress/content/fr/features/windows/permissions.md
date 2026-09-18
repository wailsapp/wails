---
title: "Autorisations"
description: "Contrôlez les demandes d’accès à la caméra, au microphone, à la géolocalisation et à d’autres fonctionnalités provenant du contenu web"
slug: "features/windows/permissions"
sourcePath: "features/windows/permissions.md"
---

Le contenu web qui appelle `navigator.mediaDevices.getUserMedia()`, l’API Geolocation ou l’API Notifications nécessite que l’application hôte accepte ou refuse ces demandes. Wails expose une table `Permissions` multiplateforme dans `WebviewWindowOptions`, qui permet de contrôler cela de manière déclarative, sans code propre à chaque plateforme.

## Démarrage rapide

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "My App",
    Permissions: map[application.PermissionType]application.Permission{
        application.PermissionMicrophone: application.PermissionAllow,
        application.PermissionCamera:     application.PermissionAllow,
    },
})
```

Les demandes d’accès à la caméra et au microphone provenant du contenu web de cette fenêtre sont acceptées sans invite du navigateur.

## Types d’autorisation

`PermissionType` (uint8) identifie une fonctionnalité que le contenu web peut demander.

| Constante | Fonctionnalité |
| --- | --- |
| `PermissionMicrophone` | `getUserMedia({audio: true})` |
| `PermissionCamera` | `getUserMedia({video: true})` |
| `PermissionGeolocation` | `navigator.geolocation` |
| `PermissionNotifications` | `Notification.requestPermission()` |
| `PermissionClipboardRead` | `navigator.clipboard.readText()` |

## Valeurs des autorisations

`Permission` (uint8) est la stratégie appliquée à un type donné.

| Constante | Valeur | Signification |
| --- | --- | --- |
| `PermissionDefault` | 0 | Utiliser le traitement natif de la plateforme (voir ci-dessous) |
| `PermissionAllow` | 1 | Accepter sans afficher d’invite |
| `PermissionDeny` | 2 | Refuser sans afficher d’invite |

`PermissionDefault` est la valeur zéro ; les entrées absentes de la table adoptent donc le comportement par défaut.

## Comportement selon la plateforme

Chaque plateforme traite `PermissionDefault` différemment, car les vues web sous-jacentes ont des comportements natifs différents.

### Linux (WebKitGTK)

WebKitGTK ne dispose d’**aucune invite d’autorisation native**. Si aucun gestionnaire n’est associé, il refuse silencieusement toutes les demandes. C’est pourquoi `getUserMedia` renvoyait toujours `NotAllowedError` avant l’ajout de cette fonctionnalité.

Actuellement, Wails traite sous Linux les demandes d’accès à la **caméra et au microphone**. La géolocalisation, les notifications et la lecture du presse-papiers ne sont pas encore raccordées et restent refusées quelle que soit la stratégie définie.

| Stratégie | Caméra / Microphone | Géolocalisation, notifications, presse-papiers |
| --- | --- | --- |
| `PermissionDefault` | **Autorisé** (rétablit getUserMedia) | Toujours refusé |
| `PermissionAllow` | Autorisé | Toujours refusé (pas encore implémenté) |
| `PermissionDeny` | Refusé | Toujours refusé |

### Windows (WebView2)

WebView2 dispose d’une invite d’autorisation native et d’une API d’autorisation par type. Les cinq types de fonctionnalités sont entièrement pris en charge.

| Stratégie | Comportement |
| --- | --- |
| `PermissionDefault` | WebView2 affiche l’invite d’autorisation native de son système d’exploitation ou de son navigateur |
| `PermissionAllow` | Accepté silencieusement |
| `PermissionDeny` | Refusé silencieusement |

**Important :** avant l’existence de cette fonctionnalité, Wails appelait systématiquement `SetGlobalPermission(Allow)`, ce qui accordait silencieusement toutes les autorisations. Désormais, dès qu’une entrée est présente dans `Permissions`, cette autorisation globale **n’est pas** définie. Les fonctionnalités non spécifiées déclenchent donc l’invite native de WebView2 au lieu d’être automatiquement autorisées.

Par conséquent, sous Windows, dès que vous configurez `Permissions`, toute fonctionnalité que vous ne répertoriez pas explicitement affiche une invite au lieu d’être autorisée silencieusement. Définissez explicitement les fonctionnalités dont vous avez besoin.

### macOS (TCC)

macOS gère l’accès à la caméra, au microphone, à la géolocalisation et aux notifications au moyen de son infrastructure système de confidentialité. L’invite du système d’exploitation s’affiche automatiquement la première fois que le contenu web demande une fonctionnalité, et le choix de l’utilisateur est mémorisé pour chaque application dans Réglages Système → Confidentialité et sécurité.

Cela fonctionne correctement sans aucune configuration de `Permissions`. La table est **actuellement ignorée sous macOS** : toutes les demandes passent par TCC, quelle que soit la valeur définie. En pratique, `PermissionDeny` est sans effet sous macOS : vous ne pouvez pas empêcher une vue web d’utiliser une fonctionnalité déjà autorisée par TCC au niveau du système.

Vérifiez que votre `Info.plist` comprend les clés de description d’utilisation appropriées :

```xml
<key>NSMicrophoneUsageDescription</key>
<string>Used for voice input</string>
<key>NSCameraUsageDescription</key>
<string>Used for video calls</string>
```

## Scénarios courants

### Application de capture multimédia

Autorisez la caméra et le microphone sur toutes les plateformes :

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionAllow,
},
```

Sous **Linux**, cela autorise explicitement les deux périphériques ; les autres fonctionnalités restent refusées. Sous **Windows**, cela autorise les deux ; toute autre fonctionnalité non répertoriée déclenchera une invite native. Sous **macOS**, cela n’a aucun effet ; TCC gère tout.

### Refuser la capture multimédia sous Linux

Linux autorise la caméra et le microphone par défaut. Pour les désactiver :

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionDeny,
    application.PermissionCamera:     application.PermissionDeny,
},
```

### Autoriser toutes les fonctionnalités sous Windows

Pour accorder toutes les autorisations sans afficher d’invite :

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone:    application.PermissionAllow,
    application.PermissionCamera:        application.PermissionAllow,
    application.PermissionGeolocation:   application.PermissionAllow,
    application.PermissionNotifications: application.PermissionAllow,
    application.PermissionClipboardRead: application.PermissionAllow,
},
```

### Politiques propres à chaque fenêtre

Chaque fenêtre peut avoir une politique différente :

```go
// Main app window — full media access
mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "App",
    Permissions: map[application.PermissionType]application.Permission{
        application.PermissionMicrophone: application.PermissionAllow,
        application.PermissionCamera:     application.PermissionAllow,
    },
})

// Settings window — no special capabilities
settingsWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Settings",
    // No Permissions entry — uses platform defaults
})

// Embedded content window — deny media capture
embeddedWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Embedded",
    Permissions: map[application.PermissionType]application.Permission{
        application.PermissionMicrophone: application.PermissionDeny,
        application.PermissionCamera:     application.PermissionDeny,
    },
})
```

## Remplacement propre à Windows

Le champ `Windows.Permissions` propre à chaque fenêtre (`map[CoreWebView2PermissionKind]CoreWebView2PermissionState`) fonctionne toujours et peut remplacer les paramètres de fonctionnalités individuelles après l’application de la table multiplateforme. Utilisez-le lorsque vous devez accéder à des types d’autorisations WebView2 sans équivalent multiplateforme (par exemple, `CoreWebView2PermissionKindOtherSensors`).

```go
Windows: application.WindowsWindow{
    Permissions: map[application.CoreWebView2PermissionKind]application.CoreWebView2PermissionState{
        application.CoreWebView2PermissionKindOtherSensors: application.CoreWebView2PermissionStateAllow,
    },
},
```

Sous Windows, l’ordre d’évaluation est le suivant :

1. Table `Permissions` multiplateforme (définit l’état de chaque type via `SetPermission`)
2. Table `Windows.Permissions` (remplace les paramètres de types individuels)
3. Pour tout type non couvert par l’une ou l’autre table : invite native de WebView2 (lorsqu’une politique est configurée) ou autorisation automatique (lorsqu’aucune politique n’est configurée — comportement historique)

## Matrice de prise en charge des plateformes

| Fonctionnalité | Linux | Windows | macOS |
| --- | --- | --- | --- |
| Microphone | ✅ | ✅ | TCC uniquement |
| Caméra | ✅ | ✅ | TCC uniquement |
| Géolocalisation | ❌ pas encore | ✅ | TCC uniquement |
| Notifications | ❌ pas encore | ✅ | TCC uniquement |
| Lecture du presse-papiers | ❌ pas encore | ✅ | TCC uniquement |

## Dépannage

**`getUserMedia` échoue toujours sous Linux après la mise à niveau**

Vérifiez que vous n’avez pas explicitement défini `PermissionMicrophone: PermissionDeny` ou `PermissionCamera: PermissionDeny`. La valeur par défaut (non définie) autorise la capture multimédia sous Linux.

**Windows demande des autorisations que je n’ai pas configurées**

Dès qu’une entrée figure dans `Permissions`, Wails n’accorde plus l’autorisation globale `Allow`. Les fonctionnalités non répertoriées déclencheront l’invite native de WebView2. Ajoutez des entrées `PermissionAllow` explicites pour chaque fonctionnalité utilisée par votre application.

**Les autorisations macOS ne fonctionnent pas**

La table `Permissions` n’a aucun effet sous macOS. Vérifiez que votre `Info.plist` contient les clés de description d’utilisation appropriées (`NSMicrophoneUsageDescription`, `NSCameraUsageDescription`, etc.) et que l’utilisateur a accordé l’accès dans Réglages Système → Confidentialité et sécurité.

**La géolocalisation, les notifications et le presse-papiers n’ont aucun effet sous Linux**

Seuls la caméra et le microphone sont actuellement pris en charge sous Linux. La prise en charge des autres types de fonctionnalités n’est pas encore implémentée : elles restent refusées quelle que soit la politique définie.
