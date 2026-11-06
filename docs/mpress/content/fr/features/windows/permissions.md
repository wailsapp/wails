---
title: "Autorisations"
description: "Contrôler les demandes d’accès à la caméra, au microphone, à la géolocalisation et à d’autres fonctionnalités provenant du contenu web"
slug: "features/windows/permissions"
sourcePath: "features/windows/permissions.md"
---

Le contenu web qui appelle `navigator.mediaDevices.getUserMedia()`, l’API Geolocation ou l’API Notifications nécessite que l’application hôte accorde ou refuse ces demandes. Wails expose une table `Permissions` multiplateforme dans `WebviewWindowOptions`, qui vous permet de contrôler ce comportement de manière déclarative, sans code propre à chaque plateforme.

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

Les demandes d’accès à la caméra et au microphone provenant du contenu web de cette fenêtre sont accordées sans invite du navigateur.

## Types d’autorisation

`PermissionType` (uint8) identifie une fonctionnalité que le contenu web peut demander.

| Constante | Fonctionnalité |
| --- | --- |
| `PermissionMicrophone` | `getUserMedia({audio: true})` |
| `PermissionCamera` | `getUserMedia({video: true})` |
| `PermissionGeolocation` | `navigator.geolocation` |
| `PermissionNotifications` | `Notification.requestPermission()` |
| `PermissionClipboardRead` | `navigator.clipboard.readText()` |

## Valeurs d’autorisation

`Permission` (uint8) est la politique appliquée à un type donné.

| Constante | Valeur | Signification |
| --- | --- | --- |
| `PermissionDefault` | 0 | Utiliser le mécanisme natif de la plateforme (voir ci-dessous) |
| `PermissionAllow` | 1 | Accorder sans demander confirmation |
| `PermissionDeny` | 2 | Refuser sans demander confirmation |

`PermissionDefault` est la valeur zéro. Les entrées absentes de la table adoptent donc le comportement par défaut.

## Comportement selon la plateforme

Chaque plateforme traite `PermissionDefault` différemment, car les vues web sous-jacentes ont des comportements natifs distincts.

### Linux (WebKitGTK)

WebKitGTK ne dispose d’**aucune invite d’autorisation native**. En l’absence de gestionnaire associé, il refuse silencieusement chaque demande. C’est pourquoi `getUserMedia` renvoyait toujours `NotAllowedError` avant l’ajout de cette fonctionnalité.

Actuellement, Wails gère les demandes d’accès à la **caméra et au microphone** sous Linux. La géolocalisation, les notifications et la lecture du presse-papiers ne sont pas encore prises en charge et restent refusées, quelle que soit la politique définie.

| Politique | Caméra / microphone | Géolocalisation, notifications, presse-papiers |
| --- | --- | --- |
| `PermissionDefault` | **Autorisé** (rétablit getUserMedia) | Toujours refusé |
| `PermissionAllow` | Autorisé | Toujours refusé (pas encore implémenté) |
| `PermissionDeny` | Refusé | Toujours refusé |

### Windows (WebView2)

WebView2 dispose d’une invite d’autorisation native et d’une API d’autorisation par type. Les cinq types de fonctionnalités sont entièrement pris en charge.

| Politique | Comportement |
| --- | --- |
| `PermissionDefault` | WebView2 affiche l’invite d’autorisation native du système d’exploitation ou du navigateur |
| `PermissionAllow` | Accordé silencieusement |
| `PermissionDeny` | Refusé silencieusement |

**Important :** avant l’existence de cette fonctionnalité, Wails appelait systématiquement `SetGlobalPermission(Allow)`, ce qui accordait silencieusement toutes les fonctionnalités. Désormais, dès qu’une entrée figure dans `Permissions`, cette autorisation générale n’est **pas** définie. Les fonctionnalités non configurées déclenchent l’invite native de WebView2 au lieu d’être accordées automatiquement.

Ainsi, si vous configurez `Permissions` sous Windows, toute fonctionnalité que vous n’indiquez pas explicitement affiche une invite au lieu d’être autorisée silencieusement. Définissez explicitement les fonctionnalités dont vous avez besoin.

### macOS (WKWebView + TCC)

Sous macOS, deux couches doivent donner leur accord. WKWebView interroge l’application avant de démarrer une session de capture, et la table `Permissions` répond à cette demande. En dessous, le mécanisme système de protection de la vie privée (TCC) contrôle l’accès à l’appareil lui-même : l’invite du système d’exploitation apparaît la première fois que l’application tente effectivement d’accéder à la caméra ou au microphone, et le choix de l’utilisateur est mémorisé pour chaque application dans Réglages Système → Confidentialité et sécurité.

Wails gère les demandes d’accès à la **caméra et au microphone** sous macOS 12 et versions ultérieures. La géolocalisation, les notifications et la lecture du presse-papiers n’ont pas d’équivalent `WKUIDelegate`. Elles ne sont donc pas prises en charge et relèvent uniquement de TCC ; comme sous Linux, la politique que vous définissez pour ces fonctionnalités n’a aucun effet.

| Politique | Caméra / microphone | Géolocalisation, notifications, presse-papiers |
| --- | --- | --- |
| `PermissionDefault` | WebKit affiche sa propre invite d’autorisation | TCC uniquement |
| `PermissionAllow` | L’invite de WebKit est ignorée — **TCC s’applique toujours** | TCC uniquement |
| `PermissionDeny` | Refusé par WebKit ; l’autorisation système est indépendante | TCC uniquement |

`PermissionAllow` accorde la demande de la vue web, mais pas l’accès à l’appareil. L’autorisation système continue de s’appliquer indépendamment, et une application dont l’accès est refusé dans Réglages Système reste bloquée. Seule l’invite propre à WebKit est ignorée.

Dans les versions antérieures à macOS 12, la méthode du délégué n’existe pas. La table y est donc ignorée et chaque demande utilise l’invite de WebKit.

Assurez-vous que votre `Info.plist` contient les clés de description d’utilisation appropriées :

```xml
<key>NSMicrophoneUsageDescription</key>
<string>Used for voice input</string>
<key>NSCameraUsageDescription</key>
<string>Used for video calls</string>
```

@note{type="caution" title="Les descriptions d’utilisation sont obligatoires"}

WebKit peut demander l’autorisation du système avant d’appeler le délégué d’autorisation de Wails. `PermissionDeny` bloque la demande de WebKit, mais ne garantit pas qu’aucune invite de TCC ni aucune demande d’autorisation d’AVFoundation ne se produira.

Incluez les descriptions d’utilisation appropriées pour les appareils auxquels votre application demande l’accès, ainsi que les droits correspondants lorsque l’application est exécutée dans un bac à sable. N’utilisez pas `PermissionDeny` pour remplacer ces déclarations.

@end

## Modèles courants

### Application de capture multimédia

Accordez l’accès à la caméra et au microphone sur toutes les plateformes :

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionAllow,
},
```

Sous **Linux**, cette configuration autorise explicitement les deux périphériques ; les autres fonctionnalités restent refusées. Sous **Windows**, elle autorise les deux ; toute autre fonctionnalité non répertoriée déclenchera une invite native. Sous **macOS**, elle autorise les deux au niveau de WebKit : aucune invite du navigateur ne s’affiche. TCC demande néanmoins l’accès aux périphériques eux-mêmes lors de leur première utilisation, et les deux clés de description d’utilisation doivent figurer dans `Info.plist`.

### Refuser la capture multimédia sous Linux

Linux autorise par défaut la caméra et le microphone. Pour désactiver cet accès :

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionDeny,
    application.PermissionCamera:     application.PermissionDeny,
},
```

### Autoriser toutes les fonctionnalités sous Windows

Pour autoriser toutes les fonctionnalités sans invite :

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

Chaque fenêtre peut avoir une politique différente :

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

Le champ `Windows.Permissions` propre à chaque fenêtre (`map[CoreWebView2PermissionKind]CoreWebView2PermissionState`) continue de fonctionner et peut remplacer les réglages de fonctionnalités individuelles après l’application de la table multiplateforme. Utilisez-le lorsque vous devez accéder à des types d’autorisations WebView2 sans équivalent multiplateforme (par exemple `CoreWebView2PermissionKindOtherSensors`).

```go
Windows: application.WindowsWindow{
    Permissions: map[application.CoreWebView2PermissionKind]application.CoreWebView2PermissionState{
        application.CoreWebView2PermissionKindOtherSensors: application.CoreWebView2PermissionStateAllow,
    },
},
```

Sous Windows, l’ordre d’évaluation est le suivant :

1. Table multiplateforme `Permissions` (définit l’état de chaque type via `SetPermission`)
2. Table `Windows.Permissions` (remplace les réglages de types individuels)
3. Pour tout type absent des deux tables : invite native de WebView2 (lorsqu’une politique est configurée) ou autorisation automatique (lorsqu’aucune politique n’est configurée — comportement historique)

## Matrice de prise en charge des plateformes

| Fonctionnalité | Linux | Windows | macOS |
| --- | --- | --- | --- |
| Microphone | ✅ | ✅ | ✅ (macOS 12+) |
| Caméra | ✅ | ✅ | ✅ (macOS 12+) |
| Géolocalisation | ❌ pas encore | ✅ | ❌ pas encore |
| Notifications | ❌ pas encore | ✅ | ❌ pas encore |
| Lecture du presse-papiers | ❌ pas encore | ✅ | ❌ pas encore |

Lorsque macOS affiche ✅, la politique répond à la demande de WebKit ; TCC contrôle néanmoins l’accès au périphérique par-dessus cette politique. Lorsque macOS affiche ❌, la fonctionnalité est laissée au seul contrôle de TCC.

## Résolution des problèmes

**`getUserMedia` échoue toujours sous Linux après la mise à niveau**

Vérifiez que vous n’avez pas défini explicitement `PermissionMicrophone: PermissionDeny` ou `PermissionCamera: PermissionDeny`. Le réglage par défaut (non défini) autorise la capture multimédia sous Linux.

**Windows demande des autorisations que je n’ai pas configurées**

Dès qu’une entrée figure dans `Permissions`, Wails n’accorde plus l’autorisation globale `Allow`. Les fonctionnalités non répertoriées déclencheront l’invite native de WebView2. Ajoutez des entrées `PermissionAllow` explicites pour chaque fonctionnalité utilisée par votre application.

**Les autorisations macOS ne fonctionnent pas**

`Permissions` prend en charge la caméra et le microphone sous macOS 12 et versions ultérieures ; la géolocalisation, les notifications et la lecture du presse-papiers ne sont pas encore raccordées et ignorent la politique. Lorsqu’elle est respectée, TCC contrôle néanmoins l’accès au périphérique par-dessus cette politique : `PermissionAllow` supprime l’invite de WebKit, mais pas celle du système. Vérifiez que votre `Info.plist` contient les clés de description d’utilisation appropriées (`NSMicrophoneUsageDescription`, `NSCameraUsageDescription`) et que l’utilisateur a accordé l’accès dans Réglages Système → Confidentialité et sécurité.

**Mon application macOS se ferme lorsque du contenu web demande l’accès à la caméra ou au microphone**

Vérifiez que `Info.plist` contient `NSCameraUsageDescription` ou `NSMicrophoneUsageDescription` pour chaque périphérique demandé par votre application. L’absence de descriptions peut entraîner le rejet de la demande ; accéder à AVFoundation sans une description requise peut provoquer l’arrêt de l’application. `PermissionDeny` contrôle la décision de WebKit et ne remplace pas ces clés.

**La géolocalisation, les notifications et le presse-papiers sont sans effet sous Linux**

Seuls la caméra et le microphone sont actuellement pris en charge sous Linux. La prise en charge des autres types de fonctionnalités n’est pas encore implémentée : ils restent refusés quelle que soit la politique définie.
