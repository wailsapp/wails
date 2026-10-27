---
title: "Blocages de WebView2 via le Bureau à distance (RDP)"
description: "Corrigez les blocages de plusieurs secondes de l’interface WebView2 lorsqu’une application Wails s’exécute via une session RDP qui modifie la résolution en PPP du moniteur en cours de session."
slug: "troubleshooting/windows/rdp"
sourcePath: "troubleshooting/windows/rdp.md"
---

## Problème

Lorsqu’une application Wails est utilisée via une session Bureau à distance (RDP), l’interface peut se bloquer pendant plusieurs secondes lors d’interactions courantes :

- Après un clic, l’ouverture d’une fenêtre contextuelle et l’affichage de son contenu prennent environ 4 à 8 secondes.
- La fermeture d’une fenêtre bloque la fenêtre parente pendant environ 2 secondes.
- Le ralentissement persiste après les reconnexions et ne disparaît qu’après le redémarrage de la machine hôte.

Ce problème se produit le plus souvent avec le client Microsoft Remote Desktop sur iOS, qui crée en cours de session un moniteur virtuel optimisé pour les écrans Retina. Tout client RDP qui introduit en cours de session un moniteur doté d’un contexte PPP différent peut provoquer le même comportement.

## Origine du problème

Par défaut, WebView2 utilise l’hébergement fenêtré, dans lequel sa surface de composition réside dans une fenêtre enfant. Lorsqu’un client RDP introduit un moniteur dont le contexte PPP diffère de celui de la session, chaque appel au contrôleur WebView2 (`PutIsVisible`, `MoveFocus`, premier rendu et libération de la surface) impose une nouvelle opération de marshaling synchrone de DirectComposition. Chaque nouvelle opération de marshaling bloque le thread d’interface utilisateur pendant environ 2 secondes, ce qui explique l’accumulation des blocages dans les applications qui utilisent beaucoup de fenêtres contextuelles.

Une application native Win32 utilisant WebView2 sur la même machine n’est pas affectée, car elle repose sur l’hébergement visuel. Cela désigne le mode d’hébergement comme la cause, plutôt qu’un problème général de WebView2 ou du compositeur Windows.

## Solution

Activez l’hébergement visuel en définissant `UseVisualHosting` dans vos options Windows. Avec l’hébergement visuel, la surface de composition de WebView2 est gérée par l’intermédiaire d’un visuel DirectComposition appartenant à l’hôte ; les changements de contexte PPP ne déclenchent donc plus de nouvelle opération de marshaling synchrone.

```go
package main

import "github.com/wailsapp/wails/v3/pkg/application"

func main() {
    app := application.New(application.Options{
        Windows: application.WindowsOptions{
            UseVisualHosting: true,
        },
    })

    // ... create your windows, then:
    app.Run()
}
```

Une fois cette option activée, les fenêtres contextuelles s’ouvrent dans le délai normal de navigation (environ 150 à 500 ms) et la fermeture d’une fenêtre ne bloque plus la fenêtre parente.

@note{type="caution"}
`UseVisualHosting` doit être défini avant `app.Run()`. Wails lit ce paramètre au démarrage de l’application et définit la variable d’environnement `COREWEBVIEW2_FORCED_HOSTING_MODE` sur `COREWEBVIEW2_HOSTING_MODE_WINDOW_TO_VISUAL` avant l’initialisation de l’environnement WebView2. Le définir ultérieurement n’a aucun effet.

@end

La valeur par défaut de l’option est `false` ; l’hébergement fenêtré reste donc le mode par défaut. Le comportement des applications existantes ne change pas, sauf si elles activent explicitement cette option.

## Quand l’activer

Définissez `UseVisualHosting: true` si votre application est régulièrement utilisée via RDP, en particulier avec le client Microsoft Remote Desktop pour iOS, et si vous constatez des blocages de plusieurs secondes à l’ouverture ou à la fermeture des fenêtres. Si votre application ne s’exécute pas via RDP, cette option n’est pas nécessaire et vous pouvez conserver sa valeur par défaut.

## Références

- [WebView2 : hébergement fenêtré ou visuel](https://learn.microsoft.com/en-us/microsoft-edge/webview2/concepts/windowed-vs-visual-hosting)
- [Ticket WebView2Feedback nº 5248](https://github.com/MicrosoftEdge/WebView2Feedback/issues/5248)
- [Ticket WebView2Feedback nº 4485](https://github.com/MicrosoftEdge/WebView2Feedback/issues/4485)
