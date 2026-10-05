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

### macOS (WKWebView + TCC)

Deux couches doivent donner leur accord sous macOS. WKWebView interroge l’application avant de lancer une capture ; la table `Permissions` répond à cette demande. En dessous, l’infrastructure de confidentialité TCC contrôle le périphérique lui-même : l’invite système apparaît au premier accès réel à la caméra ou au microphone, et le choix est mémorisé par application dans Réglages Système → Confidentialité et sécurité.

Wails traite les demandes de **caméra et de microphone** sous macOS 12 et versions ultérieures. La géolocalisation, les notifications et la lecture du presse-papiers n’ont pas d’équivalent `WKUIDelegate` : elles ne sont pas raccordées et restent confiées à TCC. Comme sous Linux, la politique définie pour elles est sans effet.

| Politique | Caméra / Microphone | Géolocalisation, Notifications, Presse-papiers |
| --- | --- | --- |
| `PermissionDefault` | WebKit affiche sa propre invite | TCC uniquement |
| `PermissionAllow` | L’invite WebKit est supprimée — **TCC s’applique toujours** | TCC uniquement |
| `PermissionDeny` | Refus avant tout accès au périphérique | TCC uniquement |

`PermissionAllow` autorise la demande de la vue web, pas le périphérique. La première capture déclenche toujours l’invite TCC et une application refusée dans Réglages Système reste refusée : aucune application ne peut s’accorder l’accès à un périphérique. `PermissionAllow` supprime uniquement l’invite WebKit préalable.

Avant macOS 12, cette méthode de délégué n’existe pas : la table est ignorée et chaque demande revient à l’invite WebKit.

Vérifiez que votre `Info.plist` comprend les clés de description d’utilisation appropriées :

```xml
<key>NSMicrophoneUsageDescription</key>
<string>Used for voice input</string>
<key>NSCameraUsageDescription</key>
<string>Used for video calls</string>
```

@note{type="caution" title="Refusez ce que votre Info.plist ne déclare pas"}

Une demande qui atteint AVFoundation sans la clé de description d’utilisation correspondante ne se contente pas d’échouer : macOS termine l’application.

`PermissionDefault` est la valeur zéro. Ainsi, `{PermissionMicrophone: PermissionAllow}` seul laisse la caméra soumise à l’invite WebKit. Si l’utilisateur accepte et que l’application déclare uniquement `NSMicrophoneUsageDescription`, elle se termine. Définissez explicitement `PermissionDeny` pour chaque fonctionnalité sans description d’utilisation :

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionDeny,
},
```

@end

## Scénarios courants

### Application de capture multimédia

Autorisez la caméra et le microphone sur toutes les plateformes :

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionAllow,
},
```

Sous **Linux**, cela autorise explicitement les deux périphériques ; les autres fonctionnalités restent refusées. Sous **Windows**, cela autorise les deux ; toute autre fonctionnalité non répertoriée déclenchera une invite native. Sous **macOS**, les deux périphériques sont autorisés au niveau WebKit, sans invite du navigateur. TCC demande toujours l’accès aux périphériques lors de la première utilisation, et les deux clés de description d’utilisation doivent figurer dans `Info.plist`.

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
| Microphone | ✅ | ✅ | ✅ (macOS 12+) |
| Caméra | ✅ | ✅ | ✅ (macOS 12+) |
| Géolocalisation | ❌ pas encore | ✅ | ❌ pas encore |
| Notifications | ❌ pas encore | ✅ | ❌ pas encore |
| Lecture du presse-papiers | ❌ pas encore | ✅ | ❌ pas encore |

Pour ✅ sous macOS, la politique répond à WebKit ; TCC contrôle en plus le périphérique. Pour ❌, la fonctionnalité reste confiée à TCC uniquement.

## Dépannage

**`getUserMedia` échoue toujours sous Linux après la mise à niveau**

Vérifiez que vous n’avez pas explicitement défini `PermissionMicrophone: PermissionDeny` ou `PermissionCamera: PermissionDeny`. La valeur par défaut (non définie) autorise la capture multimédia sous Linux.

**Windows demande des autorisations que je n’ai pas configurées**

Dès qu’une entrée figure dans `Permissions`, Wails n’accorde plus l’autorisation globale `Allow`. Les fonctionnalités non répertoriées déclencheront l’invite native de WebView2. Ajoutez des entrées `PermissionAllow` explicites pour chaque fonctionnalité utilisée par votre application.

**Les autorisations macOS ne fonctionnent pas**

`Permissions` couvre la caméra et le microphone sous macOS 12 et versions ultérieures. La géolocalisation, les notifications et la lecture du presse-papiers ne sont pas encore raccordées et ignorent la politique. TCC contrôle toujours le périphérique : `PermissionAllow` supprime l’invite WebKit, pas celle du système. Vérifiez les clés `NSMicrophoneUsageDescription` et `NSCameraUsageDescription` dans `Info.plist` et l’autorisation dans Réglages Système → Confidentialité et sécurité.

**Mon application macOS quitte lorsque le contenu web demande la caméra ou le microphone**

Une demande qui atteint AVFoundation sans la clé de description d’utilisation correspondante ne se contente pas d’échouer : macOS termine l’application. Ajoutez la clé correspondante ou définissez `PermissionDeny` pour cette fonctionnalité afin que la demande n’atteigne pas AVFoundation. `PermissionDefault` conserve l’invite WebKit, que l’utilisateur peut accepter.

**La géolocalisation, les notifications et le presse-papiers n’ont aucun effet sous Linux**

Seuls la caméra et le microphone sont actuellement pris en charge sous Linux. La prise en charge des autres types de fonctionnalités n’est pas encore implémentée : elles restent refusées quelle que soit la politique définie.
