---
title: "Fenêtres d’encoche"
description: "Créez des fenêtres macOS natives attachées au logement de la caméra."
slug: "features/windows/notch-windows"
sourcePath: "features/windows/notch-windows.md"
---

`NewNotchWindow` crée un panneau macOS non activant, doté d’une forme personnalisée et attaché au logement de la caméra. Wails gère le positionnement natif, les ailes d’attache noires, le canevas extérieur transparent, le niveau de la fenêtre, le comportement dans les Spaces ainsi que les animations facultatives d’affichage et de masquage. Votre contenu web occupe uniquement le rectangle intérieur demandé. Lorsque le pointeur entre dans la fenêtre, sa webview devient immédiatement la fenêtre clé sans activer l’application.

@note{type="caution" title="La transparence de la webview utilise une API privée"}
`NewNotchWindow` nécessite `-tags private_mac_apis` pour que le contenu web transparent laisse apparaître la forme native de l’encoche. Sans cette balise, le panneau, le positionnement et les animations fonctionnent toujours, mais la webview reste opaque. Consultez [API macOS privées](/guides/build/private-macos-apis/#webview-transparency-and-background).

@end

<figure>
  <img
    src="/images/notch-notification.gif"
    alt="Une notification d’encoche Wails glissant vers le bas depuis le logement de la caméra du MacBook, affichant des métriques système en temps réel, puis se masquant de nouveau sous l’encoche"
    loading="lazy"
    decoding="async"
    style="width: 100%; border-radius: 0.75rem"
  />
  <figcaption>
    Une notification d’encoche animée utilisant une webview persistante avec des transitions
    natives d’affichage et de masquage.
  </figcaption>
</figure>

```go
alert := app.Window.NewNotchWindow(application.NotchWindowOptions{
    Width:    660,
    Height:   92,
    Animated: true,
    WindowOptions: application.WebviewWindowOptions{
        Name: "alert",
        URL:  "/alert",
    },
})

alert.Show()
alert.Hide()
visible := alert.Visibility()
alert.Close()
```

## Options

| Champ | Type | Valeur par défaut | Description |
| --- | --- | --- | --- |
| `Width` | `int` | `660` | Largeur utilisable de la webview à l’intérieur des bords natifs profilés. |
| `Height` | `int` | `92` | Hauteur utilisable de la webview à l’intérieur des bords natifs profilés. |
| `Animated` | `bool` | `false` | Fait glisser la fenêtre vers le bas lors de `Show` et vers le haut lors de `Hide`. |
| `AnimationSpeed` | `time.Duration` | `420ms` | Durée de l’affichage. Le masquage utilise les deux tiers de cette valeur, soit `280ms` par défaut. |
| `Screen` | `*Screen` | écran principal | Cible un écran précis. Sinon, Wails utilise l’écran principal. |
| `WindowOptions` | `WebviewWindowOptions` | valeurs par défaut | Fournit le nom, l’URL ou le HTML, le CSS, le JavaScript, les raccourcis clavier et les autres comportements de la webview. |

`NewNotchWindow` contrôle les dimensions extérieures, la position, le cadre, la transparence, la politique de redimensionnement, la classe du panneau natif, le niveau de la fenêtre et le comportement de collection. Les valeurs de ces champs dans `WindowOptions` sont délibérément remplacées. Les autres champs sont conservés. Le déplacement natif par glissement de l’arrière-plan et les zones de glissement CSS sont désactivés afin que la fenêtre reste attachée au logement de la caméra.

L’objet `NotchWindow` renvoyé n’expose délibérément que `Show`, `Hide`, `Visibility` et `Close` ; la géométrie native ne peut pas être modifiée au moyen de cette interface de haut niveau.

@note{type="note"}
Les fenêtres d’encoche nécessitent macOS. Sur un Mac dépourvu de logement de caméra, Wails place la fenêtre en haut, au centre, sous la barre des menus. Sur les plateformes non prises en charge, `NewNotchWindow` renvoie un objet inerte dont les méthodes de cycle de vie n’effectuent aucune opération et peuvent être appelées sans risque, et dont `Visibility` vaut toujours false.

@end

## Cycle de vie

- `Show` affiche la fenêtre native existante. Lorsque l’animation est activée, elle glisse vers le bas depuis une position située au-dessus de l’écran.
- `Hide` maintient la fenêtre active afin de permettre sa réutilisation. Lorsque l’animation est activée, elle remonte au-dessus de l’écran avant d’être retirée de l’affichage. La webview, l’état JavaScript, les liaisons et les écouteurs d’événements restent chargés, mais la fenêtre masquée ne peut pas recevoir de survol ; l’application doit appeler `Show` pour l’afficher de nouveau.
- `Visibility` indique l’état de visibilité natif actuel.
- `Close` détruit définitivement la fenêtre native. Créez-en une nouvelle avant d’afficher de nouveau cette notification.
- Lorsque le pointeur entre dans la fenêtre d’encoche, celle-ci passe au premier plan et sa webview reçoit le focus pour permettre une interaction immédiate au clavier, tout en conservant le comportement non activant du panneau.

Chaque appel à `NewNotchWindow` crée une fenêtre indépendante, avec son propre contenu ainsi que ses propres états de visibilité et d’animation. Les fenêtres utilisent le même niveau natif ; l’instance affichée ou survolée le plus récemment apparaît donc au premier plan et peut recouvrir les instances antérieures. Lorsque les fenêtres ciblent des écrans différents, chacune est centrée sur le logement de la caméra de l’écran concerné. macOS ne coordonne pas les fenêtres d’encoche appartenant à différentes applications ; si des applications distinctes affichent des fenêtres au même emplacement et au même niveau, la fenêtre placée le plus récemment apparaît au premier plan.

Pour les notifications, les applications réutilisent généralement une fenêtre masquée ou gèrent leur propre file d’attente dans une seule webview. Wails n’impose aucune mise en file d’attente, aucun remplacement, aucun masquage automatique ni aucune politique limitant le nombre de fenêtres à une seule.

@note{type="note"}
Les surfaces réduites persistantes qui se rouvrent au survol sont délibérément distinctes du masquage des notifications et font l’objet du suivi [#6009](https://github.com/wailsapp/wails/issues/6009).

@end

Consultez l’[exemple notch-notification](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/notch-notification), qui présente un moniteur système compact dont l’état JavaScript en temps réel est conservé au fil des cycles répétés d’affichage et de masquage des notifications.
