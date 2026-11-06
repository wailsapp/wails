---
title: "Dock et barre des tâches"
description: "Gérez la visibilité de l’icône dans le Dock et affichez des badges sous macOS et Windows"
slug: "features/platform/dock"
sourcePath: "features/platform/dock.md"
---

## Introduction

Wails fournit un service Dock multiplateforme pour les applications de bureau. Ce service vous permet d’effectuer les opérations suivantes :

- Masquer et afficher l’icône de l’application dans le Dock de macOS
- Afficher des badges sur la vignette de votre application ou sur son icône dans le Dock ou la barre des tâches (macOS et Windows)

## Utilisation de base

### Création du service

Commencez par initialiser le service Dock :

```go
import "github.com/wailsapp/wails/v3/pkg/application"
import "github.com/wailsapp/wails/v3/pkg/services/dock"

// Create a new Dock service
dockService := dock.New()

// Register the service with the application
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(dockService),
    },
})
```

### Création du service avec des options de badge personnalisées (Windows uniquement)

Sous Windows, vous pouvez personnaliser l’apparence du badge à l’aide de différentes options :

```go
import "github.com/wailsapp/wails/v3/pkg/application"
import "github.com/wailsapp/wails/v3/pkg/services/dock"
import "image/color"

// Create a dock service with custom badge options
options := dock.BadgeOptions{
    TextColour:       color.RGBA{255, 255, 255, 255}, // White text
    BackgroundColour: color.RGBA{0, 0, 255, 255},     // Blue background
    FontName:         "consolab.ttf",                 // Bold Consolas font
    FontSize:         20,                             // Font size for single character
    SmallFontSize:    14,                             // Font size for multiple characters
}

dockService := dock.NewWithOptions(options)

// Register the service with the application
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(dockService),
    },
})
```

## Opérations sur le Dock

### Masquage de l’icône de l’application dans le Dock

Masquez l’icône de l’application dans le Dock de macOS :

```go
// Hide the app icon
dockService.HideAppIcon()
```

### Affichage de l’icône de l’application dans le Dock

Affichez l’icône de l’application dans le Dock de macOS :

```go
// Show the app icon
dockService.ShowAppIcon()
```

## Opérations sur les badges

### Définition d’un badge

Définissez un badge sur la vignette ou l’icône de l’application dans le Dock :

```go
// Set a default badge
dockService.SetBadge("")

// Set a numeric badge
dockService.SetBadge("3")

// Set a text badge
dockService.SetBadge("New")
```

### Définition d’un badge personnalisé (Windows uniquement)

Définissez un badge en lui appliquant des options ponctuelles :

```go
options := dock.BadgeOptions{
    BackgroundColour: color.RGBA{0, 255, 255, 255},
    FontName:         "arialb.ttf", // System font
    FontSize:         16,
    SmallFontSize:    10,
    TextColour:       color.RGBA{0, 0, 0, 255},
}

// Set a default badge
dockService.SetCustomBadge("", options)

// Set a numeric badge
dockService.SetCustomBadge("3", options)

// Set a text badge
dockService.SetCustomBadge("New", options)
```

### Suppression d’un badge

Supprimez le badge de l’icône de l’application :

```go
dockService.RemoveBadge()
```

### Récupération du badge défini

```go
dockService.GetBadge()
```

## Considérations propres aux plateformes

@tabs
[macOS]
Sous macOS :

- L’icône du Dock peut être **masquée** et **affichée**
- Les badges s’affichent directement sur l’icône du Dock
- Les options des badges ne sont **pas personnalisables** (toute option transmise à `NewWithOptions`/`SetCustomBadge` est ignorée)
- Le style standard des badges du Dock de macOS est utilisé et s’adapte automatiquement à l’apparence
- Le système gère le débordement du libellé
- Un libellé vide affiche le badge par défaut « ● »

[Windows]
Sous Windows :

- Ce service ne prend actuellement pas en charge le masquage et l’affichage de l’icône dans la barre des tâches
- Les badges s’affichent sous forme d’icône superposée dans la barre des tâches
- Les badges prennent en charge les valeurs textuelles
- L’apparence du badge peut être personnalisée à l’aide de `BadgeOptions`
- L’application doit disposer d’une fenêtre pour que les badges puissent s’afficher
- Une taille de police plus petite est automatiquement utilisée pour les libellés comportant plusieurs caractères
- Le débordement du libellé n’est pas géré
- Options de personnalisation :
  - **TextColour** : couleur du texte (par défaut : blanc)
  - **BackgroundColour** : couleur d’arrière-plan du badge (par défaut : rouge)
  - **FontName** : nom du fichier de police (par défaut : « segoeuib.ttf »)
  - **FontSize** : taille de police pour un seul caractère (par défaut : 18)
  - **SmallFontSize** : taille de police pour plusieurs caractères (par défaut : 14)


[Linux]
Sous Linux :

- La gestion de la visibilité de l’icône dans le Dock et la fonctionnalité de badge ne sont pas disponibles

@end

## Bonnes pratiques

1. **Lorsque vous masquez l’icône du Dock (macOS) :**
  - Veillez à ce que les utilisateurs puissent toujours accéder à votre application (par exemple via la [zone de notification](/features/menus/systray/))
  - Ajoutez une option « Quitter » à votre interface utilisateur alternative
  - L’application n’apparaîtra pas dans le sélecteur Command+Tab
  - Les fenêtres ouvertes restent visibles et fonctionnelles
  - La fermeture de toutes les fenêtres peut ne pas quitter l’application (le comportement varie sous macOS)
  - Les utilisateurs ne peuvent plus quitter l’application de la manière habituelle en effectuant un clic droit sur son icône dans le Dock


2. **Utilisez les badges avec parcimonie :**
  - Des mises à jour trop fréquentes des badges peuvent distraire les utilisateurs
  - Réservez les badges aux notifications importantes


3. **Utilisez un texte court dans les badges :**
  - Les badges numériques sont les plus efficaces
  - Sous macOS, le texte des badges devrait être bref


4. **Pour personnaliser les badges sous Windows :**
  - Assurez un contraste élevé entre les couleurs du texte et de l’arrière-plan
  - Testez différentes longueurs de texte, car la taille de la police diminue lorsque la longueur augmente
  - Utilisez des polices système courantes pour garantir leur disponibilité


## Référence de l’API

### Gestion du service

| Méthode | Description |
| --- | --- |
| `New()` | Crée un nouveau service de Dock |
| `NewWithOptions(options BadgeOptions)` | Crée un nouveau service de Dock avec des options de badge personnalisées (Windows uniquement ; les options sont ignorées sous macOS et Linux) |

### Opérations sur le Dock

| Méthode | Description |
| --- | --- |
| `HideAppIcon()` | Masque l’icône de l’application dans le Dock de macOS (macOS uniquement) |
| `ShowAppIcon()` | Affiche l’icône de l’application dans le Dock de macOS (macOS uniquement) |

### Opérations sur les badges

| Méthode | Description |
| --- | --- |
| `SetBadge(label string) error` | Définit un badge avec le libellé spécifié |
| `SetCustomBadge(label string, options BadgeOptions) error` | Définit un badge avec le libellé spécifié et des options de style personnalisées (Windows uniquement) |
| `RemoveBadge() error` | Supprime le badge de l’icône de l’application |
| `GetBadge() *string` | Récupère le badge actuel |

### Structures et types

```go
// Options for customizing badge appearance (Windows only)
type BadgeOptions struct {
    TextColour       color.RGBA  // Color of the badge text
    BackgroundColour color.RGBA  // Color of the badge background
    FontName         string      // Font file name (e.g., "segoeuib.ttf")
    FontSize         int         // Font size for single character
    SmallFontSize    int         // Font size for multiple characters
}
```
