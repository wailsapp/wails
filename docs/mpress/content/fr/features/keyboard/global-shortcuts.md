---
title: "Raccourcis globaux"
description: "Enregistrez des raccourcis clavier à l’échelle du système qui se déclenchent même lorsque votre application n’a pas le focus"
slug: "features/keyboard/global-shortcuts"
sourcePath: "features/keyboard/global-shortcuts.md"
---

Les raccourcis globaux sont des raccourcis clavier à l’échelle du système qui se déclenchent quelle que soit l’application ayant actuellement le focus, tant que votre application Wails est en cours d’exécution. Ils sont idéaux pour les raccourcis permettant d’afficher ou de masquer une fenêtre, les outils de capture rapide, les commandes multimédias et les autres fonctionnalités auxquelles les utilisateurs s’attendent à pouvoir accéder depuis n’importe où.

@note{type="info" title="Raccourcis globaux et raccourcis clavier"}
Les [raccourcis clavier](/features/keyboard/shortcuts/) (`app.KeyBinding`) ne se déclenchent que lorsqu’une des fenêtres de votre application a le focus. Les raccourcis globaux (`app.GlobalShortcut`) se déclenchent à l’échelle du système, même lorsque votre application est en arrière-plan. Utilisez le type qui correspond à votre besoin.

@end

Les raccourcis globaux reposent directement sur les fonctionnalités natives de chaque plateforme et n’ajoutent aucune dépendance tierce.

## Accéder au gestionnaire de raccourcis globaux

Le gestionnaire est disponible dans la propriété `GlobalShortcut` de votre instance d’application :

```go
app := application.New(application.Options{
    Name: "Global Shortcuts Demo",
})

globalShortcuts := app.GlobalShortcut
```

## Enregistrer un raccourci

`Register` accepte un accélérateur et une fonction de rappel. La fonction de rappel s’exécute dans sa propre goroutine chaque fois que le raccourci est pressé.

```go
err := app.GlobalShortcut.Register("CmdOrCtrl+Shift+G", func() {
    // Runs even when another application is focused.
    window.Show()
    window.Focus()
})
if err != nil {
    app.Logger.Error("could not register shortcut", "error", err)
}
```

Vous pouvez enregistrer des raccourcis avant d’appeler `app.Run`. La liaison avec le système d’exploitation est alors effectuée automatiquement au démarrage de l’application.

@note{type="tip" title="Interagir avec l’interface utilisateur depuis une fonction de rappel"}
Les fonctions de rappel s’exécutent en dehors du thread principal. Si votre fonction de rappel doit interagir avec des fenêtres ou d’autres éléments de l’interface utilisateur, les méthodes de fenêtre s’en chargent pour vous. Pour une opération personnalisée sur le thread principal, encapsulez-la toutefois avec `application.InvokeSync`.

@end

### Format des accélérateurs

Les raccourcis globaux utilisent le même format d’accélérateur que les accélérateurs de menu et les raccourcis clavier :

```go
"CmdOrCtrl+Shift+G"  // Command on macOS, Control elsewhere
"Ctrl+Alt+K"         // Control + Alt + K
"Cmd+Option+Space"   // Command + Option + Space (macOS)
"Super+D"            // Super / Windows / Logo key + D
"Ctrl+Shift+F5"      // Function keys are supported
```

`CmdOrCtrl` correspond à Command sous macOS et à Control sous Windows et Linux, ce qui facilite la création de raccourcis multiplateformes.

## Gérer les raccourcis

```go
// Check whether a shortcut is registered (modifier order does not matter).
registered := app.GlobalShortcut.IsRegistered("Ctrl+Shift+G")

// List every shortcut this application has registered.
for _, accelerator := range app.GlobalShortcut.GetAll() {
    app.Logger.Info("global shortcut", "accelerator", accelerator)
}

// Release a single shortcut.
app.GlobalShortcut.Unregister("Ctrl+Shift+G")

// Release everything (also done automatically on shutdown).
app.GlobalShortcut.UnregisterAll()
```

Tous les raccourcis enregistrés sont libérés automatiquement lorsque l’application se ferme ; vous n’avez donc pas besoin de les supprimer manuellement.

## Que se passe-t-il lorsque le même raccourci est enregistré deux fois ?

Il existe deux cas distincts, que Wails gère différemment.

### La même application enregistre deux fois un raccourci

Ce cas est géré par Wails lui-même et se comporte de manière identique sur toutes les plateformes. Le second appel à `Register` renvoie une erreur et la liaison d’origine reste en place (« erreur et conservation »). Ce comportement reste ainsi prévisible et signale l’erreur au lieu de remplacer silencieusement un raccourci fonctionnel.

```go
app.GlobalShortcut.Register("Ctrl+Shift+G", showWindow)        // ok
err := app.GlobalShortcut.Register("Shift+Ctrl+G", doSomething) // err: already registered
// showWindow is still the active callback for this shortcut.
```

Pour modifier la fonction de rappel d’un raccourci, appelez d’abord `Unregister` pour ce raccourci, puis appelez à nouveau `Register`.

### Une autre application détient déjà le raccourci

Ce cas est tranché par le système d’exploitation ; le résultat dépend donc de la plateforme :

| Plateforme | Comportement lorsqu’une autre application détient le raccourci |
| --- | --- |
| **macOS** | L’enregistrement réussit. macOS autorise plusieurs applications à enregistrer le même raccourci global ; votre fonction de rappel est donc ajoutée en plus de celle du détenteur existant, au lieu d’être rejetée. |
| **Windows** | L’enregistrement échoue et `Register` renvoie une erreur. L’application qui a enregistré le raccourci en premier le conserve. |
| **Linux (X11)** | L’enregistrement échoue et `Register` renvoie une erreur, car le serveur X refuse une seconde capture de la même combinaison. |
| **Linux (Wayland)** | Le compositeur arbitre. L’utilisateur est généralement invité à approuver ou à choisir la liaison dans la boîte de dialogue des raccourcis globaux de l’environnement de bureau. |

Compte tenu de ces différences, vérifiez toujours l’erreur renvoyée par `Register` et prévoyez un raccourci de secours ou informez l’utilisateur lorsqu’un raccourci ne peut pas être réservé.

## Considérations propres aux plateformes

@tabs
[macOS]
Les raccourcis globaux utilisent l’API de raccourcis clavier du Carbon Event Manager. Il s’agit du mécanisme standard de macOS pour les raccourcis clavier à l’échelle du système, qui ne nécessite pas l’autorisation Accessibilité.

Les raccourcis clavier sont liés à la position physique des touches : sur les dispositions autres que QWERTY, un raccourci correspond donc à la touche occupant la même position dans la disposition ANSI/QWERTY standard.

@note{type="caution" title="Masquage par raccourci et `ApplicationShouldTerminateAfterLastWindowClosed`"}
Sous macOS, `window.Hide()` utilise `orderOut:`, ce qui rend la fenêtre non visible. AppKit considère la dernière fenêtre non visible comme fermée. Par conséquent, si vous définissez `Mac.ApplicationShouldTerminateAfterLastWindowClosed: true` et utilisez un raccourci global pour masquer votre unique fenêtre, l’application se fermera au lieu de rester en arrière-plan. Lorsque vous dépendez d’un raccourci d’affichage et de masquage, laissez cette option non définie — comme elle l’est par défaut — afin de pouvoir masquer la fenêtre, puis la rappeler ultérieurement.

@end

[Windows]
Les raccourcis globaux utilisent l’API Win32 `RegisterHotKey`. La répétition automatique est désactivée : maintenir les touches enfoncées déclenche donc la fonction de rappel une seule fois, et non de manière répétée.

L’enregistrement échoue si une autre application détient déjà la combinaison ; privilégiez donc des combinaisons moins courantes comme valeurs par défaut.

[Linux]
Dans les sessions **X11**, Wails capture le raccourci directement auprès du serveur X. L’accélérateur demandé est donc associé exactement comme indiqué.

Dans les sessions **Wayland**, une application ne peut, par conception, pas capturer directement les touches. Wails utilise à la place l’interface `org.freedesktop.portal.GlobalShortcuts` du portail de bureau XDG. Avec ce portail, l’accélérateur que vous transmettez est un déclencheur *préféré* ; le compositeur, puis en dernier ressort l’utilisateur, décide de la combinaison de touches définitive. Votre fonction de rappel est toujours exécutée lorsque le raccourci est activé, mais rien ne garantit que les touches exactes correspondent à votre demande. De plus, `IsRegistered`/`GetAll` indiquent la combinaison que vous avez demandée, et non celle associée par le compositeur.

Le portail nécessite un environnement de bureau qui implémente le portail des raccourcis globaux, par exemple une version récente de GNOME ou de KDE Plasma.

@end

## Exemple complet

```go
package main

import (
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "Global Shortcuts Demo",
    })

    window := app.Window.New()

    // Bring the window to the front from anywhere.
    if err := app.GlobalShortcut.Register("CmdOrCtrl+Shift+G", func() {
        window.Show()
        window.Focus()
    }); err != nil {
        log.Printf("could not register show shortcut: %v", err)
    }

    // Hide the window from anywhere.
    if err := app.GlobalShortcut.Register("CmdOrCtrl+Shift+H", func() {
        window.Hide()
    }); err != nil {
        log.Printf("could not register hide shortcut: %v", err)
    }

    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

@note{type="danger" title="Éviter les raccourcis système critiques"}
Certaines combinaisons sont réservées par le système d’exploitation ou l’environnement de bureau et ne peuvent pas être utilisées par les applications. Choisissez des valeurs par défaut peu susceptibles d’entrer en conflit et gérez toujours l’erreur renvoyée par `Register`.

@end
