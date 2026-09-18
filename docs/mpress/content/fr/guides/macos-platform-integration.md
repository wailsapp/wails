---
title: "Intégration à la plateforme macOS"
description: "Fenêtres de document, compléments du Dock et des menus, panneaux natifs, éléments de barre d’état, retours, presse-papiers enrichi, glissement sortant, autorisations, alimentation et cycle de vie sur macOS"
slug: "guides/macos-platform-integration"
sourcePath: "guides/macos-platform-integration.md"
---

Plateformes concernées : macOS

Wails v3 donne à votre application les comportements que les utilisateurs de macOS attendent d’une application native : fenêtres de document avec icônes proxy et disposition en cascade, menus avec symboles et badges, menu Dock avec progression, alertes et panneaux natifs, éléments de barre d’état amovibles, retour haptique et synthèse vocale, presse-papiers enrichi avec glissement sortant, informations système sur les autorisations, l’alimentation et les paramètres régionaux, ainsi qu’une présence dans le menu Services, Handoff, AppleScript et Quick Look. Tout est piloté depuis Go au moyen du paquet `application`.

Le même code se compile sous Windows et Linux. Les méthodes de définition conservent leurs valeurs, les requêtes renvoient des valeurs zéro et les opérations nécessitant macOS renvoient une erreur documentée telle que `ErrMacOnly`, `ErrDialogNotSupported` ou `ErrClipboardNotSupported`. Les [notes de plateforme](#platform-notes) ci-dessous décrivent le comportement de chaque domaine hors de macOS.

Pour l’habillage natif des fenêtres (barres d’outils, barres latérales, inspecteurs, accessoires et onglets de fenêtre), consultez le guide [Habillage natif des fenêtres macOS](/guides/macos-native-chrome).

## Fenêtres de document

Une fenêtre de document affiche dans la barre de titre le fichier qu’elle représente, signale les modifications non enregistrées par un point dans le bouton de fermeture et ouvre les nouvelles fenêtres en cascade. Ces fonctionnalités sont toutes des méthodes de `WebviewWindow` et des options de `MacWindow`.

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:           "Report.md",
    URL:             "/editor",
    InitialPosition: application.WindowCascade,
    Mac: application.MacWindow{
        FrameAutosaveName: "editor.main",
        TitleBar: application.MacTitleBar{
            WindowButtonsOffset: &application.Point{X: 12, Y: 8},
        },
    },
})

window.SetRepresentedFile("/Users/me/Documents/Report.md")
window.SetSubtitle("Documents")
window.SetDocumentEdited(true)
```

- `SetRepresentedFile` affiche l’icône proxy du fichier dans la barre de titre. L’utilisateur peut faire glisser l’icône vers une autre application ou cliquer dessus avec la touche Commande pour afficher le chemin. Passez `""` pour la retirer. `RepresentedFile` permet de relire sa valeur.
- `SetDocumentEdited` affiche le point signalant les modifications non enregistrées dans le bouton de fermeture et estompe l’icône proxy. `IsDocumentEdited` permet de relire l’état.
- `SetSubtitle` affiche une deuxième ligne sous le titre sur macOS 11 et les versions ultérieures.
- `InitialPosition: application.WindowCascade` place la fenêtre en bas et à droite de la dernière fenêtre disposée en cascade, comme lors de l’ouverture de nouveaux documents. `CascadeFrom(other)` fait de même pour une fenêtre existante et met à jour le point de départ de la cascade pour les fenêtres suivantes.
- `Mac.FrameAutosaveName` rétablit la position et la taille enregistrées avant le premier affichage de la fenêtre, puis continue de les enregistrer lorsque celle-ci se déplace. Un cadre restauré prime sur `X`, `Y`, `Width`, `Height` et `InitialPosition`. `SetFrameAutosaveName` change le nom sur une fenêtre existante.
- `MacTitleBar.WindowButtonsOffset` déplace les boutons de fermeture, de réduction et de zoom d’un nombre de points donné. `SetWindowButtonsOffset` et `ResetWindowButtonsOffset` modifient ce décalage à l’exécution.

Ces trois méthodes de définition peuvent être appelées avant la création de la fenêtre native ; les valeurs sont appliquées lors de sa création.

### Demandes d’attention

`RequestAttention` fait rebondir l’icône du Dock lorsque l’application est en arrière-plan. Une demande informative la fait rebondir une fois. Une demande critique la fait rebondir jusqu’à ce que l’utilisateur active l’application ou que vous annuliez la demande.

```go
request := window.RequestAttention(true)

// Once the work that needed attention is done:
request.Cancel()
```

`Flash` reste la méthode multiplateforme permettant de demander une fois l’attention de l’utilisateur.

### Impression et exportation

`PrintWithOptions` imprime la WebView avec des paramètres de page explicites. La valeur zéro affiche le panneau d’impression avec les paramètres d’impression partagés. `Print` conserve son comportement historique (orientation paysage, marges de 30 points).

`ExportPDF` génère un document PDF à partir de la page et `Snapshot` en capture une image PNG. Les deux attendent WebKit : appelez-les donc depuis une goroutine et jamais depuis le thread de l’application ; un appel depuis ce thread renvoie `ErrMacExportOnMainThread`.

```go
err := window.PrintWithOptions(application.PrintOptions{
    Orientation: application.PrintOrientationPortrait,
    Margins:     application.PrintMargins{Top: 36, Left: 36, Bottom: 36, Right: 36},
    Silent:      false,
})
if err != nil {
    log.Println("print:", err)
}

go func() {
    pdf, err := window.ExportPDF(application.PDFExportOptions{})
    if err != nil {
        log.Println("export:", err)
        return
    }
    if err := os.WriteFile("report.pdf", pdf, 0o644); err != nil {
        log.Println("write:", err)
    }

    png, err := window.Snapshot(application.SnapshotOptions{Width: 800})
    if err != nil {
        log.Println("snapshot:", err)
        return
    }
    if err := os.WriteFile("preview.png", png, 0o644); err != nil {
        log.Println("write:", err)
    }
}()
```

`PrintOptions` accepte également `PrinterName`, `PaperName` (un nom PostScript tel que `"iso-a4"`) et `Scale`. `PDFExportOptions` et `SnapshotOptions` acceptent un `Rect` facultatif pour limiter la capture et un `Timeout` dont la valeur par défaut est `DefaultMacExportTimeout` (30 secondes).

## Feuilles

Une feuille est une seconde fenêtre attachée en haut de sa fenêtre parente, comme un panneau d’enregistrement. Tout `WebviewWindow` peut être présenté comme feuille d’un autre avec `PresentSheet`, puis fermé avec `EndSheet` et un code de réponse transmis aux rappels `OnSheetEnd` de la feuille. Créez la fenêtre de la feuille avec `Hidden` défini pour qu’elle n’apparaisse pas brièvement à l’écran avant son attachement ; AppKit la masque de nouveau à sa fermeture, ce qui permet de présenter plusieurs fois la même fenêtre.

```go
sheet := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "Rename",
    URL:    "/rename",
    Width:  420,
    Height: 180,
    Hidden: true,
})
sheet.OnSheetEnd(func(code int) {
    if code == application.MacSheetResponseOK {
        log.Println("renamed")
    }
})

// From a menu item or a bound method:
if err := window.PresentSheet(sheet); err != nil {
    log.Println(err)
}

// From the sheet's own page, through a bound method:
sheet.EndSheet(application.MacSheetResponseOK)
```

`PresentCriticalSheet` affiche la feuille devant toute feuille ordinaire déjà attachée, au lieu de la mettre en attente derrière elle. `PresentNativeSheet` fait de même pour un `NativeWindow`. `IsSheet`, `SheetParent`, `AttachedSheet`, `AttachedNativeSheet` et `HasAttachedSheet` décrivent l’état actuel. Fermer la fenêtre d’une feuille au lieu de terminer la feuille transmet `MacSheetResponseStop`.

## Popovers

Un `MacPopover` est un `NSPopover` : un panneau temporaire ancré à un rectangle dans une fenêtre, à un élément de barre d’outils ou à l’élément de barre d’état dans la barre de menus. Son contenu est une bande de contrôles `MacAccessory` natifs, du même type que celle utilisée pour les accessoires de barre de titre dans le guide [Habillage natif des fenêtres macOS](/guides/macos-native-chrome). Ajoutez tous les contrôles avant le premier affichage.

```go
strip := application.NewMacAccessory(application.MacAccessoryLayoutBottom)
strip.AddLabel("Sort by")
strip.AddSegmented([]string{"Date", "Title"}, 0)
strip.AddFlexibleSpace()
strip.AddButton("Apply").OnClick(func(*application.Context) {
    // apply the sort
})

popover := application.NewMacPopover(application.MacPopoverOptions{
    Width:    320,
    Behavior: application.MacPopoverBehaviorTransient,
    Content:  strip,
})
popover.OnClose(func() {
    log.Println("popover closed")
})

// Anchored to a rectangle in the page, in window content coordinates:
err := popover.ShowRelativeTo(application.Rect{X: 20, Y: 60, Width: 120, Height: 28}, window, application.MacRectEdgeMaxY)
if err != nil {
    log.Println(err)
}
```

`MacToolbarItem.ShowPopover` et `SystemTray.ShowPopover` ancrent le même popover à un élément de barre d’outils ou à l’élément de barre d’état. `MacPopoverBehaviorTransient` ferme le popover lors de tout clic à l’extérieur, `Semitransient` uniquement lors de clics dans la fenêtre qui le présente, et le comportement par défaut le garde ouvert jusqu’à l’appel de `Close`. `MacRectEdge` choisit le côté où apparaît le popover. `SetContentSize` et `SetBehavior` ajustent un popover ouvert, et `Destroy` libère le popover natif ainsi que la bande de contenu pour une utilisation ailleurs.

## Restauration de l’état

macOS rouvre les fenêtres d’une application après un plantage, une fermeture forcée ou un redémarrage, ainsi qu’après une fermeture normale lorsque l’option « Fermer les fenêtres à la fermeture d’une application » est désactivée dans les Réglages Système. Attribuez un `Mac.RestorationID` à la fenêtre, enregistrez les informations nécessaires à sa recréation avec `SetRestorationData` et inscrivez `app.Window.OnRestore` pour reconstruire la fenêtre au lancement suivant.

```go
app.Window.OnRestore(func(id string, state application.RestorationState) application.Window {
    if id != "editor" {
        return nil
    }
    path := state.Get("path")
    restored := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: filepath.Base(path),
        URL:   "/editor?path=" + path,
        Mac:   application.MacWindow{RestorationID: id},
    })
    restored.SetRestorationData(state.Data)
    return restored
})

window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Report.md",
    URL:   "/editor?path=/Users/me/Documents/Report.md",
    Mac:   application.MacWindow{RestorationID: "editor"},
})
window.SetRestorationData(map[string]string{"path": "/Users/me/Documents/Report.md"})
```

Seules les fenêtres visibles lorsque l’application se termine sont enregistrées. `RestorationState.Data` est une map de chaînes ; limitez son contenu aux identifiants, chemins et positions. `InteractionState` renvoie la liste de navigation arrière et avant ainsi que les positions de défilement de la WebView sous forme d’un bloc de données opaque (macOS 12+), que `RestoreInteractionState` applique à une fenêtre recréée ; ce bloc est généralement stocké encodé en base64 dans les données de restauration. `SetRestorationID` et `RestorationID` modifient et lisent l’identifiant d’une fenêtre existante.

## Options de présentation

`MacPresentationOptions` reflète `NSApplication.presentationOptions` : un masque de bits qui masque le Dock ou la barre de menus et désactive le changement de processus, la fermeture forcée, la déconnexion ou la commande Masquer pendant que l’application est active. Définissez `Mac.PresentationOptions` dans les options de l’application pour l’appliquer au lancement, ou modifiez-la à l’exécution avec `SetPresentationOptions`. Les combinaisons invalides sont rejetées avant de parvenir à AppKit, avec une erreur enveloppant `ErrMacPresentationOptionsInvalid`.

```go
kiosk := application.MacPresentationHideDock |
    application.MacPresentationHideMenuBar |
    application.MacPresentationDisableProcessSwitching |
    application.MacPresentationDisableForceQuit

if err := app.SetPresentationOptions(kiosk); err != nil {
    log.Println(err)
}
log.Println("presentation:", app.PresentationOptions())

// Restore the standard Dock and menu bar:
_ = app.SetPresentationOptions(application.MacPresentationDefault)
```

Masquer la barre de menus (`HideMenuBar` ou `AutoHideMenuBar`) nécessite l’une des options Dock, et `AutoHideToolbar` nécessite à la fois `FullScreen` et `AutoHideMenuBar`. `Validate` signale la première règle qu’une valeur enfreint, et `Has` teste les indicateurs individuels.

## Menus et Dock

Les éléments de menu bénéficient des SF Symbols, des badges, des en-têtes de section, des palettes de couleurs, de l’état de coche mixte, des variantes et de l’indentation. Toutes ces fonctionnalités sont des méthodes de `MenuItem` et de `Menu` : elles fonctionnent donc dans le menu de l’application, les menus contextuels, les menus de la barre d’état et le menu Dock.

```go
menu := app.NewMenu()
menu.AddRole(application.AppMenu)

fileMenu := menu.AddSubmenu("File")
fileMenu.Add("Open...").SetSymbol("folder").SetAccelerator("CmdOrCtrl+o").
    OnClick(func(*application.Context) {
        // open the file, then note it in the recent list:
        app.Menu.AddRecentDocument("/Users/me/Documents/Report.md")
    })
fileMenu.AddRole(application.OpenRecent)

view := menu.AddSubmenu("View")
view.AddSectionHeader("Mailboxes")
inbox := view.Add("Inbox").SetSymbol("tray.full").SetBadge(3)
inbox.OnClick(func(*application.Context) {
    inbox.ClearBadge()
})
view.Add("Updates").SetBadgeText("New")

view.AddSeparator()
wrap := view.AddCheckbox("Wrap lines", false).SetMixed()
view.Add("Reset").SetIndentationLevel(1).OnClick(func(*application.Context) {
    wrap.SetMixed()
})

view.AddSeparator()
view.Add("Close Tab").SetAccelerator("CmdOrCtrl+w").OnClick(func(*application.Context) {})
view.Add("Close All Tabs").SetAccelerator("CmdOrCtrl+OptionOrAlt+w").SetAlternate(true).
    OnClick(func(*application.Context) {})

colours := []application.RGBA{
    application.NewRGB(255, 59, 48),
    application.NewRGB(255, 149, 0),
    application.NewRGB(52, 199, 89),
}
view.AddPalette([]string{"tag.fill"}, colours, 0, func(_ *application.Context, index int) {
    log.Println("tag colour", index)
}).SetLabel("Tag colour")

app.Menu.Set(menu)
```

- `SetSymbol` affiche un SF Symbol à côté du titre (macOS 11+). Il remplace une image définie avec `SetBitmap`.
- `SetBadge` affiche un nombre et `SetBadgeText` une courte chaîne après le titre (macOS 14+). `ClearBadge` le retire ; `BadgeCount` et `BadgeText` permettent de relire sa valeur.
- `AddSectionHeader` ajoute un en-tête non interactif (macOS 14+). Sur les versions antérieures, il devient un élément désactivé portant le même titre.
- `AddPalette` ajoute une rangée d’échantillons de couleurs reposant sur le menu de palette de `NSMenu` (macOS 14+). Passez un symbole pour chaque échantillon, un par couleur, ou une tranche vide pour obtenir des cercles pleins. Sans libellé, la palette apparaît directement dans le menu parent ; `SetLabel` la présente comme un sous-menu titré. `PaletteSelected` renvoie l’index sélectionné.
- `SetMixed` place une case à cocher dans l’état mixte, représenté par un tiret. Un clic la coche entièrement, conformément au comportement d’AppKit.
- `SetAlternate(true)` affiche l’élément à la place de celui qui le précède lorsque la touche de modification qui les distingue est maintenue. Les deux éléments doivent partager une même touche et différer par leurs touches de modification.
- `SetIndentationLevel` indente le titre sur un maximum de 15 niveaux.

### Ouvrir un élément récent

`fileMenu.AddRole(application.OpenRecent)` ajoute le sous-menu standard des éléments récents. Sur macOS, `NSDocumentController` le remplit chaque fois qu’il s’ouvre et y inclut une commande pour effacer le menu. Ajoutez des fichiers avec `app.Menu.AddRecentDocument`, répertoriez-les avec `RecentDocuments` et videz la liste avec `ClearRecentDocuments`. La liste est conservée entre les redémarrages.

Le choix d’un fichier récent produit le même événement que l’ouverture d’un fichier depuis Finder :

```go
app.Event.OnApplicationEvent(events.Common.ApplicationOpenedWithFile, func(event *application.ApplicationEvent) {
    path := event.Context().Filename()
    log.Println("open", path)
})
```

### Menu Dock

`app.Menu.SetDockMenu` installe un menu statique accessible par clic droit sur l’icône du Dock. `OnDockMenu` en construit un à la demande chaque fois qu’il va être affiché ; c’est le bon choix lorsque les éléments reflètent un état changeant. Le constructeur prend le pas sur le menu statique.

```go
dockMenu := application.NewMenu()
dockMenu.Add("New Document").SetSymbol("doc.badge.plus").OnClick(func(*application.Context) {
    // create a document
})
app.Menu.SetDockMenu(dockMenu)

app.Menu.OnDockMenu(func() *application.Menu {
    dynamic := application.NewMenu()
    dynamic.AddSectionHeader("Recent")
    for _, path := range app.Menu.RecentDocuments() {
        path := path
        dynamic.Add(filepath.Base(path)).OnClick(func(*application.Context) {
            app.Menu.OpenRecentDocument(path)
        })
    }
    return dynamic
})
```

### Progression dans le Dock

Le service Dock dessine une barre de progression sur l’icône du Dock, en complément de la prise en charge existante des badges. Enregistrez `dock.New()` comme service et appelez `SetProgress` avec une fraction comprise entre 0 et 1.

```go
import "github.com/wailsapp/wails/v3/pkg/services/dock"

dockService := dock.New()

app := application.New(application.Options{
    Name: "Exporter",
    Services: []application.Service{
        application.NewService(dockService),
    },
})

app.Event.On("export:progress", func(event *application.CustomEvent) {
    if fraction, ok := event.Data.(float64); ok {
        _ = dockService.SetProgress(fraction)
    }
})
app.Event.On("export:done", func(*application.CustomEvent) {
    _ = dockService.ClearProgress()
})
```

`GetProgress` renvoie la fraction actuelle, ou `nil` si aucune barre n’est affichée.

## Boîtes de dialogue

Les boîtes de dialogue de message, d’ouverture et d’enregistrement acceptent des options macOS, et le gestionnaire de boîtes de dialogue propose en plus une invite de saisie ainsi que les panneaux système de couleur et de police.

### Alertes

`SetSuppression` ajoute une case « Ne plus afficher ce message » et `SetHelp` affiche le bouton d’aide. Lisez l’état de la case avec `Suppressed` depuis un rappel de bouton, ou enregistrez `OnSuppression` pour le recevoir en premier.

```go
dialog := app.Dialog.Warning().
    SetTitle("Delete 3 items?").
    SetMessage("The items will be moved to the Bin.").
    SetSuppression("Do not warn me again").
    SetHelp(func() {
        log.Println("help requested")
    }).
    AttachToWindow(window)

dialog.AddButton("Delete").SetAsDefault().OnClick(func() {
    if dialog.Suppressed() {
        // remember not to ask again
    }
})
dialog.AddButton("Cancel").SetAsCancel()
dialog.Show()
```

### Invite de saisie

`Prompt` affiche une alerte avec un champ de texte et bloque jusqu’à sa fermeture : appelez-le donc depuis une goroutine ou une méthode liée. `Secure` transforme le champ en champ de mot de passe et `Window` présente l’alerte comme une feuille.

```go
go func() {
    value, ok, err := app.Dialog.Prompt(application.PromptOptions{
        Title:        "Name this document",
        Message:      "The name is used for the exported file.",
        Placeholder:  "Untitled",
        DefaultValue: "Quarterly report",
        OKLabel:      "Rename",
        Window:       window,
    })
    if err != nil || !ok {
        return
    }
    log.Println("renamed to", value)
}()
```

### Panneaux de fichiers

`AddContentType` filtre par identifiant de type uniforme dans les boîtes de dialogue d’ouverture et d’enregistrement. Il s’utilise avec `AddFilter` : `"public.image"` correspond ainsi à tous les types d’images connus du système, tandis qu’un filtre peut toujours sélectionner les PDF par extension.

```go
paths, err := app.Dialog.OpenFile().
    SetTitle("Choose images").
    AddFilter("PDF", "*.pdf").
    AddContentType("public.image").
    AttachToWindow(window).
    PromptForMultipleSelection()
if err == nil {
    log.Println(paths)
}
```

Les panneaux d’enregistrement proposent en plus un menu déroulant Format, un libellé personnalisé pour le champ du nom et des tags Finder. `SetFormats` change le type autorisé et l’extension du champ du nom lorsque l’utilisateur change l’option du menu, et `SelectedFormat` indique son choix final.

```go
formats := []application.DialogFormat{
    {Label: "PNG image", Extension: "png", UTI: "public.png"},
    {Label: "PDF document", Extension: "pdf", UTI: "com.adobe.pdf"},
}
save := app.Dialog.SaveFile().
    SetFilename("Quarterly report.png").
    SetNameFieldLabel("Export As:").
    SetTags([]string{"Reports", "Draft"}).
    AttachToWindow(window)
save.SetFormats(formats, 0, nil)

path, err := save.PromptForSingleSelection()
if err == nil && path != "" {
    log.Println("export", path, "as", formats[save.SelectedFormat()].Label)
}
```

### Panneaux de couleur et de police

`PickColor` et `PickFont` ouvrent les panneaux système partagés et bloquent jusqu’à leur fermeture. `OnChange` transmet chaque sélection tant que le panneau est ouvert, ce qui permet à la page d’afficher le choix en aperçu immédiat. Un seul panneau de chaque type peut être ouvert à la fois ; un deuxième appel renvoie `ErrDialogInProgress`.

```go
go func() {
    colour, changed, err := app.Dialog.PickColor(application.ColorPickerOptions{
        Initial:    application.NewRGB(52, 120, 246),
        ShowsAlpha: true,
        Title:      "Accent colour",
        OnChange: func(colour application.RGBA) {
            app.Event.Emit("accent:preview", colour)
        },
    })
    if err == nil && changed {
        log.Println("accent", colour)
    }

    font, changed, err := app.Dialog.PickFont(application.FontPickerOptions{
        Family: "Helvetica Neue",
        Size:   18,
    })
    if err == nil && changed {
        log.Println(font.Family, font.Face, font.PostScriptName, font.Size)
    }
}()
```

## Éléments de barre d’état et retours

### Éléments de barre d’état

Sur macOS, un élément de la barre d’état système est un `NSStatusItem`. Il peut être dessiné à partir d’un SF Symbol, porter une infobulle et être retiré par l’utilisateur comme les éléments intégrés au système.

```go
tray := app.SystemTray.New()
tray.SetSymbol("waveform.circle").SetSymbolConfiguration(0, application.MacSymbolWeightMedium)
tray.SetTooltip("Recorder. Command-drag to remove.")
tray.SetRemovable(true, "com.example.recorder.tray")
tray.OnVisibilityChange(func(visible bool) {
    log.Println("status item visible:", visible)
})

trayMenu := app.Menu.New()
trayMenu.Add("Show").OnClick(func(*application.Context) {
    window.Show().Focus()
})
tray.SetMenu(trayMenu)
tray.Run()
```

- `SetSymbol` affiche le symbole comme image modèle afin qu’il suive l’apparence de la barre de menus (macOS 11+). `SetSymbolConfiguration` définit la taille en points et la graisse.
- `SetTooltip` définit le texte affiché au survol. `Tooltip` permet de le relire.
- `SetRemovable(true, name)` permet à l’utilisateur de retirer l’élément de la barre de menus en le faisant glisser avec la touche Commande. Donnez-lui un nom d’enregistrement automatique stable pour que macOS mémorise son retrait entre les lancements. `Show` ou `SetVisible(true)` le fait réapparaître.
- `IsVisible` lit `NSStatusItem.visible` : sa valeur est donc fausse après le retrait de l’élément par l’utilisateur. `OnVisibilityChange` signale chaque changement.

### Retour haptique

`app.Haptics.Perform` produit un motif sur un trackpad Force Touch ou Magic Trackpad pendant que l’application est active.

```go
app.Haptics.Perform(application.HapticAlignment)
```

Les types sont `HapticGeneric`, `HapticAlignment` (un élément s’aligne à sa place) et `HapticLevelChange` (un cran ou une étape de clic). `IsSupported` indique si la plateforme peut produire un retour haptique.

### Sons

`app.Sound` joue le son d’alerte, un son système nommé ou un fichier audio.

```go
app.Sound.Beep()

if err := app.Sound.Play("Glass"); err != nil {
    log.Println(err)
}

for _, name := range app.Sound.SystemSounds() {
    log.Println(name)
}
```

`Play` accepte un nom provenant de `SystemSounds` ou un chemin absolu vers tout fichier que Core Audio sait décoder. `PlayData` lit un fichier audio complet depuis la mémoire.

### Voix

`app.Speech.Speak` place du texte en attente pour la voix du système et renvoie un `Utterance`. Les énoncés sont lus les uns après les autres ; `Stop` en supprime un et `StopAll` vide la file d’attente. `Voices` répertorie les voix installées avec leurs identifiants et leurs langues.

```go
utterance, err := app.Speech.Speak("Export finished", application.SpeechOptions{
    Voice: "com.apple.voice.compact.en-GB.Daniel",
    Rate:  0.5,
})
if err != nil {
    log.Println(err)
    return
}
utterance.OnFinished(func() {
    log.Println("done, stopped:", utterance.WasStopped())
})
```

`Recognize` transcrit le microphone par défaut avec `SFSpeechRecognizer`. Le premier appel demande les autorisations d’accès au microphone et de reconnaissance vocale, puis bloque jusqu’à la réponse de l’utilisateur : appelez-le donc depuis une goroutine. Les transcriptions partielles arrivent par `OnPartial` ; `Stop` termine la capture et renvoie le texte final.

```go
go func() {
    session, err := app.Speech.Recognize(application.RecognitionOptions{
        Locale: "en-US",
        OnPartial: func(text string) {
            app.Event.Emit("dictation:partial", text)
        },
    })
    if err != nil {
        if errors.Is(err, application.ErrSpeechRecognitionDenied) {
            _ = app.Permissions.OpenSystemSettings(application.PermissionKindMicrophone)
        }
        log.Println(err)
        return
    }
    time.Sleep(5 * time.Second)
    text, err := session.Stop()
    if err != nil {
        log.Println(err)
        return
    }
    app.Event.Emit("dictation:final", text)
}()
```

La reconnaissance nécessite une application empaquetée dont le fichier `Info.plist` déclare `NSSpeechRecognitionUsageDescription` et `NSMicrophoneUsageDescription`. Sans ces clés, macOS refuse l’accès et `Recognize` renvoie `ErrSpeechRecognitionUsageDescription`.

## Presse-papiers et glissement

### Presse-papiers enrichi

`app.Clipboard` lit et écrit des images, des références de fichiers, du HTML, du RTF et des données brutes sous n’importe quel identifiant de type uniforme, en plus du texte brut. `Types` répertorie le contenu du presse-papiers et `OnChange` signale les modifications effectuées par n’importe quelle application.

```go
if err := app.Clipboard.SetHTML("<p>Rich <b>HTML</b></p>", "Rich HTML"); err != nil {
    log.Println(err)
}
_ = app.Clipboard.SetFiles([]string{"/Users/me/Documents/Report.md"})
_ = app.Clipboard.SetData("com.example.record", []byte(`{"id":42}`))

stop := app.Clipboard.OnChange(func() {
    log.Println("clipboard changed", app.Clipboard.ChangeCount(), app.Clipboard.Types())
    if files, err := app.Clipboard.Files(); err == nil && len(files) > 0 {
        log.Println("files:", files)
    }
})
defer stop()
```

`SetImage` et `Image` utilisent des octets PNG ; les images copiées au format TIFF par d’autres applications sont converties automatiquement. Le système ne fournit aucune notification pour les changements du presse-papiers : tant qu’au moins un écouteur est présent, `OnChange` vérifie donc le compteur de modifications toutes les 500 ms.

### Glissement sortant

`StartDrag` lance un glissement système depuis la fenêtre, comme si l’utilisateur avait saisi les éléments dans Finder. Il peut proposer des fichiers existants, des promesses de fichiers dont le contenu n’est produit que lorsqu’une destination accepte le dépôt, ou du texte brut. Lancez-le pendant un geste de souris : liez une méthode Go et appelez-la depuis le gestionnaire `mousedown` ou `pointerdown` de la page sur l’élément déplaçable, en définissant l’attribut HTML `draggable` sur `false` pour que WebKit ne lance pas son propre glissement.

```go
// DragService is bound to the page and called from a mousedown handler.
type DragService struct {
    window *application.WebviewWindow
}

func (s *DragService) DragExport() error {
    return s.window.StartDrag(application.DragItems{
        Promises: []application.DragPromise{{
            Filename: "export.csv",
            Data: func() ([]byte, error) {
                return []byte("id,name\n1,Wails\n"), nil
            },
        }},
        Operations: application.DragOperationCopy,
    })
}
```

```go
window.OnDragEnd(func(operation application.DragOperation) {
    log.Println("drag ended:", operation)
})
```

Appeler `StartDrag` hors d’un geste renvoie `ErrDragOutNoGesture`. `DragItems.Image` et `ImageOffset` définissent l’image sous le curseur.

### Dépôts depuis d’autres applications

Les dépôts de fichiers continuent d’utiliser l’événement `WindowFilesDropped`. Pour accepter du texte, des URL ou des images glissés depuis d’autres applications, répertoriez les types dans `DropTypes` et inscrivez `OnDrop`. Ces dépôts sont transmis à Go plutôt qu’aux gestionnaires de dépôt HTML5 propres à la page.

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Inbox",
    URL:   "/",
    DropTypes: []application.DropType{
        application.DropFiles,
        application.DropText,
        application.DropURLs,
        application.DropImages,
    },
})

window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
    log.Println("files:", event.Context().DroppedFiles())
})

window.OnDrop(func(_ *application.Context, data application.DropData) {
    log.Println("text:", data.Text, "urls:", data.URLs, "images:", len(data.Images), "at", data.X, data.Y)
})
```

## Système

### Autorisations

`app.Permissions` indique et demande les autorisations de confidentialité du système (caméra, microphone, enregistrement de l’écran, accessibilité, localisation, notifications, surveillance des entrées et accès complet au disque). `Status` n’affiche jamais de demande. `Request` affiche une demande pour les autorisations dont le statut n’est pas encore déterminé et bloque jusqu’à la réponse de l’utilisateur : appelez-le donc depuis une goroutine. `OpenSystemSettings` ouvre le volet correspondant de Confidentialité et sécurité.

```go
go func() {
    status := app.Permissions.Status(application.PermissionKindCamera)
    if status == application.PermissionStatusNotDetermined {
        status, _ = app.Permissions.Request(application.PermissionKindCamera)
    }
    if status == application.PermissionStatusDenied {
        _ = app.Permissions.OpenSystemSettings(application.PermissionKindCamera)
    }
}()
```

L’accès complet au disque ne peut pas être demandé et renvoie `ErrPermissionNotRequestable` ; dirigez l’utilisateur vers le volet de réglages. Une demande qui ne reçoit jamais de réponse renvoie `ErrPermissionRequestTimeout`, ce qui signifie généralement sur macOS que la clé de description d’utilisation correspondante manque dans `Info.plist`.

L’option de fenêtre `Permissions` est désormais prise en compte sur macOS. Elle détermine le traitement des demandes `getUserMedia` de la page : `PermissionAllow` évite la propre demande de la WebView, `PermissionDeny` refuse sans demander et `PermissionDefault` affiche la demande. La demande TCC au niveau du système apparaît toujours lors de la première utilisation de la caméra ou du microphone.

```go
app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Meeting",
    URL:   "/",
    Permissions: map[application.PermissionType]application.Permission{
        application.PermissionMicrophone: application.PermissionAllow,
        application.PermissionCamera:     application.PermissionDefault,
    },
})
```

### Alimentation

`app.Power.PreventSleep` maintient le système éveillé et, avec `Display`, l’écran également, jusqu’à l’appel de la fonction de libération renvoyée. Les maintiens sont comptabilisés, de sorte que plusieurs parties de l’application peuvent en détenir simultanément. La raison est affichée dans Moniteur d’activité.

```go
release, err := app.Power.PreventSleep("Exporting video", application.PreventSleepOptions{Display: true})
if err != nil {
    log.Println(err)
}
defer release()

state := app.Power.State()
if state.LowPowerMode || state.ThermalState >= application.ThermalStateSerious {
    // trim background work
}
log.Println("battery", state.BatteryLevel, "charging", state.Charging, "on battery", state.OnBattery)
```

Les changements arrivent sous forme de `events.Mac.ApplicationDidChangePowerState` (basculement du mode économie d’énergie) et de `events.Mac.ApplicationDidChangeThermalState`.

### Cycle de vie

macOS peut terminer instantanément une application inactive lors d’une déconnexion ou d’un arrêt si elle active `NSSupportsSuddenTermination`, et peut fermer une application inactive sans fenêtre si elle active `NSSupportsAutomaticTermination`. `app.Lifecycle.HoldTermination` suspend ces deux mécanismes pendant une section critique, comme l’enregistrement d’un fichier.

```go
release := app.Lifecycle.HoldTermination("Saving document")
defer release()
// write the file
```

`SetSuddenTerminationEnabled` active ou désactive la terminaison soudaine à l’exécution ; `SuddenTerminationEnabled` indique l’état actuel, initialisé à partir de la clé `Info.plist`.

### Environnement

`app.Env` propose trois nouvelles requêtes sur les réglages de l’utilisateur.

```go
a11y := app.Env.Accessibility()
if a11y.ReduceMotion || a11y.ReduceTransparency {
    app.Event.Emit("theme:calm", true)
}

layout := app.Env.KeyboardLayout()
log.Println(layout.ID, layout.Name, layout.Languages)

locale := app.Env.Locale()
log.Println(locale.Identifier, locale.Language, locale.Region, locale.Preferred)
```

- `Accessibility` reflète Réduire les animations, Réduire la transparence, Augmenter le contraste, Différencier sans couleur, Inverser les couleurs, VoiceOver et Contrôle de sélection.
- `KeyboardLayout` renvoie la source de saisie active avec son identifiant, son nom localisé et ses langues.
- `Locale` renvoie les paramètres régionaux qu’AppKit a sélectionnés pour l’application ainsi que la liste ordonnée complète `Preferred` de l’utilisateur. `Identifier` ne reflète qu’une langue déclarée par le bundle dans `CFBundleLocalizations` ; utilisez `Preferred` pour choisir vous-même une langue.

### Événements

Ces événements d’application sont nouveaux. Chacun est transmis par `app.Event.OnApplicationEvent` ; interrogez le gestionnaire correspondant pour obtenir la nouvelle valeur.

| Événement | Déclenchement | Lecture avec |
|-------|------------|-----------|
| `events.Mac.ApplicationDidChangePowerState` | basculement du mode économie d’énergie | `app.Power.State()` |
| `events.Mac.ApplicationDidChangeThermalState` | changement de la pression thermique | `app.Power.State()` |
| `events.Common.AccessibilitySettingsChanged` | changement d’un réglage d’affichage lié à l’accessibilité | `app.Env.Accessibility()` |
| `events.Mac.ApplicationDidChangeAccessibilitySettings` | forme macOS du même changement | `app.Env.Accessibility()` |
| `events.Mac.ApplicationDidChangeKeyboardLayout` | changement de la source de saisie | `app.Env.KeyboardLayout()` |
| `events.Mac.ApplicationDidChangeLocale` | changement des paramètres régionaux | `app.Env.Locale()` |

```go
app.Event.OnApplicationEvent(events.Mac.ApplicationDidChangeThermalState, func(*application.ApplicationEvent) {
    app.Event.Emit("system:power", app.Power.State())
})
app.Event.OnApplicationEvent(events.Common.AccessibilitySettingsChanged, func(*application.ApplicationEvent) {
    app.Event.Emit("system:accessibility", app.Env.Accessibility())
})
```

## Intégration

### Menu Services

`app.ServicesProvider.Register` ajoute une entrée au sous-menu Services que chaque application macOS affiche pour du texte ou des fichiers sélectionnés. Le gestionnaire reçoit le presse-papiers sous forme de `ServiceRequest` et renvoie un `ServiceResponse` pour y réécrire des données ; une réponse vide laisse la sélection intacte. Gardez le gestionnaire rapide, car AppKit l’attend sur le thread principal.

```go
err := app.ServicesProvider.Register(application.ServiceDefinition{
    Name:          "summarise",
    MenuTitle:     "Summarise with Notes",
    SendTypes:     []string{"public.utf8-plain-text"},
    ReturnTypes:   []string{"public.utf8-plain-text"},
    KeyEquivalent: "S",
    Handler: func(_ *application.Context, request application.ServiceRequest) (application.ServiceResponse, error) {
        return application.ServiceResponse{Text: "Summary: " + request.Text}, nil
    },
})
if err != nil {
    log.Println(err)
}

// Paste this inside the top-level <dict> of Info.plist.
log.Println(app.ServicesProvider.InfoPlistXML())
```

`Name` est le message envoyé par AppKit et doit être un simple identifiant. `SendTypes` et `ReturnTypes` sont des types de presse-papiers ; un service doit en avoir au moins un. L’enregistrement seul ne rend pas le service visible : le fichier `Info.plist` du bundle doit le déclarer sous `NSServices`. `InfoPlistXML` renvoie ce bloc prêt à être collé et `InfoPlistEntries` renvoie les mêmes données sous forme de maps pour un sérialiseur plist. Dans ces entrées, `NSPortName` correspond au `Name` de l’application, qui doit être identique à `CFBundleName`.

Les projets compilés avec la CLI Wails peuvent déclarer ces mêmes services une seule fois dans `build/config.yml` et laisser le processus d’empaquetage produire le bloc `NSServices` :

```yaml
services:
  - name: SummariseText
    menuTitle: Summarise with My Product
    sendTypes:
      - public.utf8-plain-text
    returnTypes:
      - public.utf8-plain-text
    keyEquivalent: S
```

Chaque entrée doit correspondre à une `ServiceDefinition` enregistrée en Go avec le même `Name`. Exécutez `pbs -update` après l’installation d’une nouvelle version pour que le menu Services prenne en compte le changement sans déconnexion.

### Handoff et activités utilisateur

`app.Activity.Publish` rend une `NSUserActivity` active pour que l’utilisateur puisse la poursuivre sur un autre appareil, la retrouver dans Spotlight ou recevoir une suggestion de Siri. Le `PublishedActivity` renvoyé peut être mis à jour lorsque l’état change et invalidé à la fermeture du document. Publier une nouvelle activité remplace la précédente.

```go
activity, err := app.Activity.Publish(application.UserActivity{
    Type:               "com.example.notes.editing",
    Title:              "Editing Quarterly report",
    UserInfo:           map[string]any{"note": "quarterly-report"},
    WebpageURL:         "https://example.com/notes/quarterly-report",
    EligibleForHandoff: true,
    EligibleForSearch:  true,
    Keywords:           []string{"report", "quarterly"},
})
if err != nil {
    log.Println(err)
    return
}
if err := activity.Update(map[string]any{"note": "quarterly-report", "cursor": 120}); err != nil {
    log.Println(err)
}
// When the document closes:
activity.Invalidate()
```

Les activités entrantes arrivent par `OnContinue`. Les liens universels portent le type `UserActivityTypeBrowsingWeb`, avec la page dans `WebpageURL`, et sont aussi transmis sous forme de `events.Common.ApplicationLaunchedWithUrl` afin qu’une application puisse utiliser un seul chemin de traitement des URL. `OnWillContinue`, `OnFailed` et `OnUpdated` couvrent le reste du délégué.

```go
app.Activity.OnContinue(func(_ *application.Context, incoming application.UserActivity) bool {
    if incoming.Type == application.UserActivityTypeBrowsingWeb {
        log.Println("universal link", incoming.WebpageURL)
        return true
    }
    note, _ := incoming.UserInfo["note"].(string)
    log.Println("continue editing", note)
    return true
})
```

Chaque type d’activité doit figurer sous `NSUserActivityTypes` dans `Info.plist`. Les liens universels nécessitent également l’autorisation `com.apple.developer.associated-domains` avec une entrée `applinks:example.com` et le fichier `apple-app-site-association` correspondant sur ce domaine.

### Apple Events

`app.AppleEvents.Handle` enregistre un gestionnaire pour une classe et un ID d’événement, afin qu’AppleScript, Raccourcis et d’autres applications puissent piloter l’application. Les codes sont des chaînes de quatre caractères. Le paramètre direct est décodé en une valeur Go (`string`, `[]string` de chemins de fichiers, `int64`, `float64`, `bool`, `[]any` ou `AppleEventRawData`) et le `Result` de la réponse accepte les mêmes types. Les gestionnaires s’exécutent chacun dans leur propre goroutine pendant que l’événement est suspendu.

```go
err := app.AppleEvents.Handle("WAIL", "note", func(_ *application.Context, event application.AppleEvent) (application.AppleEventReply, error) {
    text, _ := event.DirectObject.(string)
    return application.AppleEventReply{Result: "noted: " + text}, nil
})
if err != nil {
    log.Println(err)
}

go func() {
    result, err := app.AppleEvents.Send("com.apple.finder", "misc", "actv", nil)
    log.Println(result, err)
}()

sdef := app.AppleEvents.ScriptingDefinition()
if err := os.WriteFile("Notes.sdef", []byte(sdef), 0o644); err != nil {
    log.Println(err)
}
```

Un script peut appeler immédiatement le gestionnaire avec la syntaxe d’événement brut :

```applescript
tell application id "com.example.notes" to «event WAILnote» "hello"
```

`ScriptingDefinition` génère un `.sdef` minimal qui donne un nom de commande à chaque gestionnaire. Incluez-le dans `Contents/Resources` et référencez-le dans `Info.plist` avec `NSAppleScriptEnabled` et `OSAScriptingDefinition` ; l’Éditeur de scripts l’affiche alors sous Fichier > Ouvrir le dictionnaire. Wails gère déjà l’événement Get URL pour les schémas d’URL personnalisés ; un gestionnaire pour `"GURL"`/`"GURL"` s’ajoute à ce traitement, tandis qu’un gestionnaire pour `"aevt"`/`"odoc"` remplace la transmission intégrée de l’événement d’ouverture de documents. `Send` cible une application en cours d’exécution par son identifiant de bundle, nécessite `NSAppleEventsUsageDescription` dans une application empaquetée et bloque la goroutine appelante jusqu’à l’arrivée de la réponse.

### Quick Look

`app.QuickLook.Preview` ouvre le panneau Quick Look partagé pour un ou plusieurs fichiers ; avec plusieurs chemins, le panneau affiche des flèches permettant de passer de l’un à l’autre. `Thumbnail` produit une miniature de fichier avec les fournisseurs de miniatures du système et renvoie un PNG ; les documents, images, PDF et vidéos sont donc tous pris en charge.

```go
if err := app.QuickLook.Preview([]string{"/Users/me/Documents/Report.pdf"}); err != nil {
    log.Println(err)
}

go func() {
    png, err := app.QuickLook.Thumbnail("/Users/me/Documents/Report.pdf", application.ThumbnailOptions{
        Width: 256,
        Scale: 2,
    })
    if err != nil {
        log.Println(err)
        return
    }
    if err := os.WriteFile("thumbnail.png", png, 0o644); err != nil {
        log.Println(err)
    }
}()
```

Les chemins doivent être absolus et exister. `ClosePreview` et `IsPreviewOpen` gèrent le panneau. `ThumbnailOptions.IconMode` dessine la bordure de document dans le style de Finder, et `Scale: 2` produit une image Retina. `Thumbnail` bloque la goroutine appelante : appelez-le donc depuis une goroutine ou une méthode liée.

### Utilitaires d’espace de travail

`app.Browser` propose trois utilitaires autour de `NSWorkspace`. `OpenWith` ouvre un fichier avec une application précise désignée par son identifiant de bundle ou par le chemin de son bundle. `ApplicationsForFile` répertorie les applications installées capables d’ouvrir un fichier, avec le gestionnaire par défaut en premier. `ActivateApplication` place une application en cours d’exécution au premier plan.

```go
apps := app.Browser.ApplicationsForFile("/Users/me/Documents/Report.md")
for _, info := range apps {
    log.Println(info.Name, info.BundleID, info.Path)
}

if err := app.Browser.OpenWith("/Users/me/Documents/Report.md", "com.apple.TextEdit"); err != nil {
    log.Println(err)
}

if err := app.Browser.ActivateApplication("com.apple.TextEdit"); err != nil {
    log.Println(err)
}
```

### Spotlight

`app.Spotlight.Index` ajoute du contenu de l’application à l’index de recherche du système au moyen de Core Spotlight. Chaque `SearchableItem` possède un `ID`, un `Title` et, facultativement, un `Domain` pour la suppression groupée, une `Description`, des `Keywords`, un `ContentType`, une miniature PNG, une `URL` de lien profond et une date d’expiration. `OnOpen` est appelé lorsque l’utilisateur choisit un de ces éléments dans Spotlight ou sélectionne « Rechercher dans l’app » avec une requête.

```go
err := app.Spotlight.Index([]application.SearchableItem{{
    ID:          "note:quarterly-report",
    Domain:      "notes",
    Title:       "Quarterly report",
    Description: "Draft for the board meeting",
    Keywords:    []string{"finance", "q3"},
    ContentType: "public.plain-text",
    URL:         "notes://open/quarterly-report",
}})
if err != nil {
    log.Println(err)
}

stop := app.Spotlight.OnOpen(func(_ *application.Context, id string, query string) {
    if query != "" {
        log.Println("search in app:", query)
        return
    }
    log.Println("open item", id)
})
defer stop()

// Later:
_ = app.Spotlight.Delete([]string{"note:quarterly-report"})
_ = app.Spotlight.DeleteDomain("notes")
```

`IsAvailable` indique si l’index accepte des éléments. L’indexation nécessite une application empaquetée : les éléments indexés par un binaire `go run` non empaqueté n’apparaissent jamais dans Spotlight. `DeleteAll` supprime tous les éléments indexés par l’application.

## À venir

Un exemple `mac-windows-extra` couvrant les feuilles, les popovers, les options de présentation et la restauration de l’état est en cours d’ajout ; il sera lié ici lorsqu’il sera disponible.

## Versions requises

Tout ce qui figure sur cette page est propre à macOS, et l’API Go est identique partout. Wails cible macOS 10.13 et les versions ultérieures ; les fonctionnalités nécessitant une version plus récente se dégradent comme indiqué.

| Fonctionnalité | Version minimale de macOS | Comportement sur les versions antérieures |
|---------|---------------|-------------------------------|
| Statut des autorisations de caméra et de microphone | 10.14 | signalé comme autorisé (les versions antérieures ne contrôlent pas l’accès aux périphériques de capture) |
| `Speech.Speak` et `Voices` | 10.14 | `ErrSpeechNotSupported` |
| Autorisations d’enregistrement de l’écran et de surveillance des entrées | 10.15 | signalées comme autorisées |
| `Speech.Recognize` | 10.15 | `ErrSpeechRecognitionNotSupported` |
| `QuickLook.Thumbnail` | 10.15 | `ErrQuickLookNotSupported` |
| `SetSubtitle` | 11 | ignoré avec une entrée dans le journal de débogage |
| `ExportPDF` | 11 | `ErrMacExportUnsupported` |
| `SetSymbol` sur les éléments de menu et de barre d’état | 11 | aucune image n’est affichée |
| `AddContentType`, `SetFormats` par UTI | 11 | les mêmes identifiants sont appliqués au moyen de l’ancienne API des types de fichiers autorisés |
| Icônes des promesses de fichiers dans `StartDrag` | 11 | une icône de document générique |
| `PowerState.LowPowerMode` et son événement | 12 | toujours faux ; l’événement ne se déclenche jamais |
| `InteractionState` et `RestoreInteractionState` | 12 | `ErrMacInteractionStateUnsupported` |
| Volet Notifications dans `OpenSystemSettings` | 13 | l’ancien volet des préférences Notifications s’ouvre |
| Badges de menu, en-têtes de section, palettes | 14 | les badges ne s’affichent pas ; les en-têtes sont des éléments désactivés ; les palettes sont masquées |
| `MacToolbarItem.ShowPopover` sur les éléments sans vue personnalisée | 14 | `ErrMacPopoverAnchorUnavailable` |

Tout le reste ne nécessite aucune version au-delà du minimum requis par Wails.

## Clés Info.plist

Plusieurs fonctionnalités dépendent de clés dans le fichier `Info.plist` de l’application. Les descriptions d’utilisation sont présentées à l’utilisateur dans la demande d’autorisation ; sans la clé, macOS n’affiche jamais la demande et celle-ci expire.

| Clé | Nécessaire pour |
|-----|-----------|
| `NSCameraUsageDescription` | `Permissions.Request(PermissionKindCamera)`, accès à la caméra depuis la page |
| `NSMicrophoneUsageDescription` | `Permissions.Request(PermissionKindMicrophone)`, accès au microphone depuis la page, `Speech.Recognize` |
| `NSSpeechRecognitionUsageDescription` | `Speech.Recognize` |
| `NSLocationUsageDescription` | `Permissions.Request(PermissionKindLocation)` |
| `NSSupportsSuddenTermination` | l’état initial de `Lifecycle.SuddenTerminationEnabled` ; `HoldTermination` le suspend |
| `NSSupportsAutomaticTermination` | permet à macOS de fermer l’application inactive ; `HoldTermination` suspend ce mécanisme |
| `CFBundleLocalizations` | les langues que `Env.Locale().Identifier` peut indiquer |
| `NSServices` | une entrée par service enregistré avec `app.ServicesProvider` ; générez le bloc avec `InfoPlistXML` |
| `NSUserActivityTypes` | chaque `UserActivity.Type` publié ou poursuivi par `app.Activity` |
| `NSAppleScriptEnabled` et `OSAScriptingDefinition` | indiquent que l’application peut être pilotée par script et désignent le `.sdef` produit par `app.AppleEvents.ScriptingDefinition` |
| `NSAppleEventsUsageDescription` | `app.AppleEvents.Send` vers d’autres applications |

L’autorisation des notifications et la reconnaissance vocale nécessitent aussi que l’application s’exécute comme un bundle doté d’un identifiant de bundle ; un binaire `go run` non empaqueté indique `PermissionStatusUnsupported` pour les notifications.

<a id="platform-notes"></a>

## Notes de plateforme

Toutes les API de cette page se compilent sous Windows et Linux. Hors de macOS :

- Compléments de fenêtre : `SetRepresentedFile`, `SetDocumentEdited`, `SetSubtitle`, `CascadeFrom`, `SetFrameAutosaveName` et `SetWindowButtonsOffset` sont sans effet. `WindowCascade` se comporte comme `WindowCentered`. `RequestAttention` renvoie une référence dont `Cancel` est sans effet. `PrintWithOptions` appelle `Print`. `ExportPDF` et `Snapshot` renvoient `ErrMacOnly`.
- Menus : les symboles, badges, états mixtes, variantes et indentations sont conservés, mais pas dessinés. Les en-têtes de section sont des éléments désactivés. Les palettes sont masquées. Le sous-menu des éléments récents est rempli à partir de la liste Go des éléments récents lors de la création du menu. Les menus Dock ne s’affichent jamais.
- Boîtes de dialogue : `SetSuppression`, `SetHelp`, `SetNameFieldLabel` et `SetTags` sont ignorés. `AddContentType` et `SetFormats` deviennent des filtres d’extension lorsqu’une correspondance est connue. `Prompt`, `PickColor` et `PickFont` renvoient `ErrDialogNotSupported`.
- Éléments de barre d’état : `SetSymbol`, `SetRemovable` et `OnVisibilityChange` sont sans effet. `IsVisible` reflète le dernier appel à `Show` ou `Hide`.
- Retours : le retour haptique fonctionne sur iOS et Android ; sous Windows et Linux, les appels sont sans effet. `Sound.Beep` et `Sound.Play` fonctionnent sous Windows avec les fichiers WAV et les alias du registre ; ailleurs, `Play` renvoie `ErrSoundNotSupported`. Les fonctions vocales renvoient `ErrSpeechNotSupported` et `ErrSpeechRecognitionNotSupported`.
- Presse-papiers : les méthodes enrichies renvoient `ErrClipboardNotSupported`, `Types` est vide, `ChangeCount` vaut 0 et `OnChange` ne se déclenche jamais.
- Glissement : `StartDrag` renvoie `ErrDragOutUnsupported`. Les types de dépôt autres que les fichiers ne sont pas transmis ; les dépôts de fichiers continuent de fonctionner par `WindowFilesDropped`.
- Système : `Permissions.Status` indique `PermissionStatusUnsupported` et `Request` renvoie `ErrPermissionsUnsupported`. `PreventSleep` renvoie `ErrPreventSleepUnsupported` avec une fonction de libération sans effet. `HoldTermination` renvoie une fonction de libération sans effet. Tous les champs de `Accessibility` sont faux, `KeyboardLayout` prend la valeur zéro et `Locale` est déduit de `LC_ALL`, `LC_MESSAGES` et `LANG`. L’option de fenêtre `Permissions` est multiplateforme.
- Intégration : `ServicesProvider.Register` renvoie `ErrServicesUnsupported` et `Activity.Publish` renvoie `ErrActivityUnsupported` ; `InfoPlistXML`, `InfoPlistEntries` et les gestionnaires d’activité fonctionnent toujours. `Handle` et `Send` de `AppleEvents` renvoient `ErrAppleEventsNotSupported` ; `ScriptingDefinition` est généré partout. `QuickLook.Preview` et `Thumbnail` renvoient `ErrQuickLookNotSupported`. Les méthodes d’indexation de `Spotlight` renvoient `ErrSpotlightNotSupported` et `OnOpen` ne se déclenche jamais. `Browser.OpenWith` lance l’exécutable désigné avec le chemin comme argument, `ApplicationsForFile` est vide et `ActivateApplication` renvoie `ErrApplicationNotRunning`.
- Options de présentation : `SetPresentationOptions` renvoie `ErrMacOnly` et `PresentationOptions` vaut `MacPresentationDefault`.
- Feuilles : `PresentSheet`, `PresentCriticalSheet` et `PresentNativeSheet` renvoient `ErrMacSheetUnsupported` ; `EndSheet` est sans effet et les méthodes de requête n’indiquent aucune feuille.
- Popovers : `NewMacPopover` fonctionne, les méthodes d’affichage renvoient `ErrMacPopoverUnsupported` et `IsShown` vaut faux.
- Restauration de l’état : `SetRestorationID` et `SetRestorationData` sont sans effet, `OnRestore` n’est jamais appelé, et `InteractionState` ainsi que `RestoreInteractionState` renvoient `ErrMacOnly`.

## Exemples

Chaque exemple est une application complète et exécutable :

- [`v3/examples/mac-toolbar`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-toolbar) : `windowextras.go` intègre à l’éditeur de notes le point indiquant les modifications, le sous-titre, l’exportation PDF, les options d’impression et les fenêtres en cascade.
- [`v3/examples/mac-menus-dock`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-menus-dock) : symboles, badges, en-têtes de section, état mixte, variantes, palette, éléments récents, menu Dock dynamique et progression dans le Dock.
- [`v3/examples/mac-dialogs`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-dialogs) : case de suppression et bouton d’aide, invites de saisie, types de contenu, menu Format, tags Finder et panneaux de couleur et de police.
- [`v3/examples/mac-feedback`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-feedback) : élément de barre d’état SF Symbol amovible, retour haptique, sons système, synthèse vocale et reconnaissance vocale.
- [`v3/examples/mac-clipboard-drag`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-clipboard-drag) : presse-papiers enrichi avec suivi des changements, glissement sortant avec promesses de fichiers et dépôts de texte, d’URL et d’images.
- [`v3/examples/mac-system`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-system) : autorisations, prévention de la mise en veille, suspension de la terminaison, état d’alimentation, accessibilité, disposition du clavier et paramètres régionaux avec mises à jour en direct.
- [`v3/examples/mac-integration`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-integration) : entrée du menu Services, activité Handoff avec gestionnaires de continuation et utilitaires d’espace de travail sur `app.Browser`.
- [`v3/examples/mac-search-preview`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-search-preview) : indexation Spotlight avec `OnOpen`, aperçus et miniatures Quick Look, et Apple Event personnalisé avec sa définition de script.
