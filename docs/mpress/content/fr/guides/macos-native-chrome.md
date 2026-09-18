---
title: "Habillage natif des fenêtres macOS"
description: "Créez des barres d’outils, barres latérales, listes de contenu, inspecteurs, accessoires et onglets de fenêtre AppKit natifs autour de votre fenêtre Wails"
slug: "guides/macos-native-chrome"
sourcePath: "guides/macos-native-chrome.md"
---

Plateformes concernées : macOS

Wails v3 peut entourer votre WebView de véritables éléments de fenêtre AppKit. La barre d’outils, la barre latérale, la liste de contenu, l’inspecteur et les bandes de la barre de titre sont des contrôles natifs créés depuis Go. Aucun de ces éléments n’est en HTML : ils bénéficient donc directement des matériaux, de la gestion du clavier, des animations et de la persistance d’AppKit, tandis que votre interface reste centrée sur le contenu.

Toutes les API de cette page se trouvent dans le paquet `application` et portent le préfixe `Mac`. Le même code se compile sous Windows et Linux : les constructeurs et les méthodes de définition fonctionnent partout, les appels d’attachement sont sans effet ou renvoient une erreur, et la fenêtre conserve son unique WebView ordinaire.

## Anatomie d’une fenêtre

Une fenêtre complète comprend les éléments suivants, du bord initial au bord final :

| Élément | Type | Classe AppKit |
|------|------|--------------|
| Barre d’outils | `MacToolbar` | `NSToolbar` |
| Barre latérale | `MacSidebar` | liste de sources `NSOutlineView` dans un volet de barre latérale |
| Liste de contenu | `MacContentList` | `NSTableView` dans un volet de liste de contenu |
| Contenu principal | votre WebView ou `MacTextEditor` | `WKWebView` ou `NSTextView` |
| Inspecteur | `MacInspector` | contrôles de propriétés natifs dans un volet d’inspecteur |
| Accessoires | `MacAccessory` | `NSTitlebarAccessoryViewController` ou `NSSplitViewItemAccessoryViewController` |

Les volets sont disposés par un `MacSplitView`, qui est un `NSSplitViewController`. Créez d’abord les éléments, ajoutez-les à la vue fractionnée dans l’ordre du bord initial au bord final, attachez la vue fractionnée à la fenêtre, puis attachez la barre d’outils. Vous pouvez le faire avec `SetSplitView` et `SetToolbar`, ou en une seule étape au moyen des options de fenêtre `Mac.SplitView` et `Mac.Toolbar`.

Une configuration de fenêtre typique associe cet habillage aux options de fenêtre suivantes pour permettre au contenu de défiler sous une barre d’outils unifiée :

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Mac: application.MacWindow{
        ContentLayout: application.MacContentLayoutEdgeToEdge,
        TitleBar: application.MacTitleBar{
            FullSizeContent:      true,
            HideToolbarSeparator: true,
            ToolbarStyle:         application.MacToolbarStyleUnified,
        },
    },
})
```

## Barre d’outils

`NewMacToolbar` crée un `NSToolbar`. Ajoutez des éléments avec les méthodes `Add`, enchaînez les méthodes de définition et les rappels sur les références renvoyées, puis attachez la barre d’outils avec `SetToolbar`. Les identifiants sont générés automatiquement.

```go
toolbar := application.NewMacToolbar().
    SetDisplayMode(application.MacToolbarDisplayModeIconOnly)

// Standard AppKit items. They have no handle because AppKit owns them.
toolbar.AddSidebarToggle()
toolbar.AddSidebarTrackingSeparator()

toolbar.AddButton("New").
    SetSymbol("square.and.pencil").
    SetTooltip("Create a note").
    SetBordered(true).
    OnClick(func(*application.Context) {
        // create a note
    })

toolbar.AddSearch("Search").
    SetSearchPlaceholder("Search notes").
    SetSearchIncremental(true).
    OnSearch(func(_ *application.Context, query string) {
        filterNotes(query)
    })

toolbar.AddFlexibleSpace()

mode := toolbar.AddGroup("Mode", application.ToolbarGroupSelectOne)
mode.AddButton("Write").SetSymbol("pencil").OnClick(func(*application.Context) {
    mode.SetSelectedIndex(0)
})
mode.AddButton("Preview").SetSymbol("doc.richtext").OnClick(func(*application.Context) {
    mode.SetSelectedIndex(1)
})

actions := application.NewMenu()
actions.Add("Export...").OnClick(func(*application.Context) {})
toolbar.AddMenu("Actions", actions).SetSymbol("ellipsis.circle")

window.SetToolbar(toolbar)
```

Les types d’éléments sont les suivants :

- `AddButton` ajoute un bouton poussoir. Chaque bouton doit disposer d’un `OnClick` avant l’attachement de la barre d’outils ; sinon, `SetToolbar` signale une erreur et laisse la barre d’outils précédente en place.
- `AddSearch` ajoute un `NSSearchToolbarItem`. `OnSearch` est requis. Utilisez `SetSearchPlaceholder`, `SetSearchIncremental`, `SetSearchRecentsKey` pour un menu persistant des recherches récentes et `SetSearchMenu` pour un menu personnalisé derrière la loupe.
- `AddShare` ajoute l’élément de partage du système (voir ci-dessous).
- `AddGroup` ajoute un `NSToolbarItemGroup` segmenté. Ajoutez ses membres avec la méthode `AddButton` du groupe et choisissez `ToolbarGroupSelectOne`, `ToolbarGroupMomentary` ou `ToolbarGroupSelectAny`.
- `AddMenu` ajoute un `NSMenuToolbarItem` déroulant piloté par un `Menu` ordinaire. `SetShowsIndicator(false)` masque le chevron.
- `AddSpace` et `AddFlexibleSpace` ajoutent les espaces standard.
- `AddSidebarToggle` et `AddSidebarTrackingSeparator` ajoutent les éléments de barre latérale d’AppKit. Le séparateur maintient tous les éléments qui le précèdent alignés au-dessus de la séparation de la barre latérale ; la fenêtre doit donc avoir une vue fractionnée avec un volet de barre latérale.
- `AddInspectorToggle` et `AddInspectorTrackingSeparator` font de même pour un volet d’inspecteur.

`SetDisplayMode` permet de choisir entre `MacToolbarDisplayModeIconAndLabel` (valeur par défaut), `MacToolbarDisplayModeIconOnly`, `MacToolbarDisplayModeLabelOnly` et `MacToolbarDisplayModeDefault`.

### Mises à jour en direct

Chaque référence permet aussi les mises à jour en direct. Les méthodes de définition appelées après l’attachement mettent à jour l’élément natif sur le thread de l’application.

```go
save := toolbar.AddButton("Save").SetSymbol("checkmark.circle")
save.OnClick(func(*application.Context) {
    save.SetBadgeCount(0).SetProminent(false)
})

// Later, when the document changes:
save.SetBadgeCount(1).SetProminent(true)
save.SetVisibilityPriority(application.MacToolbarVisibilityPriorityHigh)

toolbar.Move(save, 0)
toolbar.Remove(save)
```

Parmi les autres méthodes utiles figurent `SetLabel`, `SetTooltip`, `SetEnabled`, `SetHidden`, `SetTintColor` et `SetNavigational`, qui maintient un élément au bord initial, comme Safari le fait avec ses boutons Précédent et Suivant. `SetVisibilityPriority` détermine quels éléments passent en premier dans le menu de débordement lorsque la fenêtre rétrécit.

### Personnalisation par l’utilisateur

`SetCustomizable` active la feuille standard « Personnaliser la barre d’outils… » et enregistre la disposition choisie par l’utilisateur sous la clé que vous fournissez. Attribuez à chaque élément une clé stable avec `SetPersistenceKey` pour que la disposition enregistrée survive à un redémarrage, et utilisez `SetInDefaultSet(false)` pour les éléments qui ne doivent apparaître qu’après leur ajout par l’utilisateur. Appelez `SetCustomizable` avant d’attacher la barre d’outils.

```go
toolbar := application.NewMacToolbar().SetCustomizable("myapp.main-toolbar")

newNote := toolbar.AddButton("New").
    SetPersistenceKey("new").
    OnClick(func(*application.Context) {})

// Offered in the palette but hidden until the user adds it.
toolbar.AddButton("Archive").
    SetPersistenceKey("archive").
    SetInDefaultSet(false).
    OnClick(func(*application.Context) {})

toolbar.SetCenteredItems(newNote)

// From a menu item, for example:
toolbar.RunCustomizationPalette()
```

### Partage

`AddShare` renvoie un `MacToolbarShareItem`. Il reste désactivé jusqu’à ce qu’un `MacShareProvider` annonce au moins une représentation. Wails ne demande les octets au fournisseur que lorsqu’un service de partage les réclame : les exportations volumineuses sont ainsi produites à la demande. `MacShareProviderFunc` adapte deux fonctions en fournisseur ; les applications avec état peuvent implémenter directement l’interface.

```go
share := toolbar.AddShare("Share")
share.SetProvider(application.MacShareProviderFunc{
    Available: []application.MacShareRepresentation{
        {ContentType: application.MacShareTypePDF},
        {ContentType: application.MacShareTypePlainText},
    },
    Load: func(request application.MacShareRequest) ([]byte, error) {
        if request.ContentType == application.MacShareTypePDF {
            return renderPDF()
        }
        return []byte(currentText()), nil
    },
}).SetSubject("Saturday, slowly").SetSuggestedName("Note")

share.OnShared(func(_ *application.Context, service string) {
    // service is the localised name of the sharing service
})
share.OnShareError(func(_ *application.Context, service string, err error) {
    window.Error("share via %s failed: %s", service, err)
})
```

Les types de contenu courants sont `MacShareTypePlainText`, `MacShareTypeHTML`, `MacShareTypePDF`, `MacShareTypePNG` et `MacShareTypeJPEG`. Toute autre chaîne UTI est également acceptée.

### Attachement et détachement

Une barre d’outils ne peut appartenir qu’à une fenêtre à la fois. `WebviewWindow.SetToolbar` et `NativeWindow.SetToolbar` acceptent une barre d’outils avant ou après la création de la fenêtre ; passer `nil` la retire et la libère pour une utilisation ailleurs. Sur un `WebviewWindow`, `SetToolbar` signale les problèmes de validation par `Window.Error`, tandis que la version pour `NativeWindow` renvoie l’erreur. L’option de fenêtre `Mac.Toolbar` attache une barre d’outils lors de la création ; elle est appliquée après `Mac.SplitView` afin qu’un séparateur de suivi puisse trouver la barre latérale sur laquelle il s’aligne.

## Vue fractionnée

`MacSplitView` dispose les volets. Ajoutez-les dans l’ordre du bord initial au bord final. `AddSidebar`, `AddContentList` et `AddInspector` prennent le modèle natif qu’ils hébergent et renvoient un `MacSplitPane` pour régler la taille et le repli. `AddPrimaryContent` place la WebView existante de la fenêtre et renvoie un `MacSplitWebviewPane`, qui ajoute `SetContentLayout`.

```go
split := application.NewMacSplitView().SetAutosaveName("myapp.main-window")

sidebarPane := split.AddSidebar(sidebar).
    SetMinimumThickness(200).
    SetMaximumThickness(320).
    SetCollapsible(true)

split.AddContentList(list).
    SetMinimumThickness(240).
    SetCollapsible(true)

split.AddPrimaryContent().
    SetContentLayout(application.MacContentLayoutEdgeToEdge)

split.AddInspector(inspector).
    SetPreferredThicknessFraction(0.25).
    SetHoldingPriority(300).
    SetCollapsible(true).
    SetCanCollapseFromWindowResize(false).
    SetCollapsed(true)

sidebarPane.OnCollapsedChange(func(_ *application.Context, collapsed bool) {
    // fired for toggles, gestures, menu items and SetCollapsed alike
})

window.SetSplitView(split)

// At runtime:
sidebarPane.Toggle()
```

`SetSplitView` fonctionne avant ou après la création de la fenêtre native. Avant sa création, la disposition est mise en attente et installée lors de la création. Après sa création, par exemple depuis un rappel de menu ou de barre d’état dans une application en cours d’exécution, la disposition est installée immédiatement : la WebView existante de la fenêtre devient le volet principal, la barre d’outils actuelle est rattachée pour aligner ses séparateurs de suivi et tous les accessoires en attente sont attachés.

Une fenêtre créée après `app.Run` peut être configurée en un seul appel avec les options `Mac.SplitView` et `Mac.Toolbar` :

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Notes",
    URL:   "/",
    Mac: application.MacWindow{
        SplitView: split,
        Toolbar:   toolbar,
    },
})
```

Les règles d’une disposition sont les suivantes :

- Une disposition nécessite au moins deux volets et exactement un volet principal (`AddPrimaryContent` pour un `WebviewWindow`, `AddTextEditor` pour un `NativeWindow`).
- Elle peut contenir au plus une liste de contenu, placée après la barre latérale et avant le volet principal.
- La structure des volets est figée dès que la vue fractionnée est attachée. Les paramètres des volets, leur état de repli et le contenu de la barre latérale, de la liste et de l’inspecteur peuvent toujours changer à tout moment.
- Une disposition installée ne peut pas être remplacée. Un deuxième appel à `SetSplitView` sur la même fenêtre signale `ErrMacSplitViewAlreadyInstalled` (par `Window.Error` sur un `WebviewWindow`, comme valeur de retour sur un `NativeWindow`) et laisse la fenêtre inchangée. Passer `nil` avant l’installation efface une disposition en attente.
- Une barre latérale, une liste, un inspecteur ou une vue fractionnée ne peut appartenir qu’à une fenêtre à la fois.

Les méthodes de définition des volets sont `SetMinimumThickness`, `SetMaximumThickness`, `SetPreferredThicknessFraction`, `SetHoldingPriority`, `SetCollapsible`, `SetCanCollapseFromWindowResize`, `SetCollapsed`, `Toggle`, `IsCollapsed` et `OnCollapsedChange`. `SetAutosaveName` conserve la position des séparateurs entre les lancements.

Sur le volet principal, `SetContentLayout` permet de choisir entre `MacContentLayoutBelowToolbar` et `MacContentLayoutEdgeToEdge`. `MacContentLayoutAutomatic` hérite de `MacWindow.ContentLayout`, qui suit à son tour `TitleBar.FullSizeContent`. La disposition bord à bord permet à AppKit d’appliquer son effet de bord de défilement sous la barre d’outils sur macOS 26 et les versions ultérieures.

## Barre latérale

`MacSidebar` est une liste de sources native. Elle contient des lignes racines, des sections et des lignes imbriquées à n’importe quelle profondeur. Conservez les références renvoyées pour mettre les lignes à jour plus tard.

```go
sidebar := application.NewMacSidebar()

// A root row above the sections.
sidebar.AddItem("All Notes").
    SetSymbol("tray.full").
    SetBadge(12).
    OnClick(func(*application.Context) {})

notes := sidebar.AddSection("Notes")
draft := notes.AddItem("Saturday, slowly").
    SetSymbol("doc.text").
    SetAccessorySymbol("pin.fill").
    SetTooltip("Field notes").
    SetEditable(true).
    OnRename(func(_ *application.Context, label string) {
        // the row label is already updated
    })

// Nested rows with a tinted symbol.
tags := sidebar.AddSection("Tags")
tint := application.NewRGB(0, 122, 255)
work := tags.AddItem("Work").SetSymbol("briefcase").SetTintColor(&tint)
work.AddItem("Meetings").SetSymbol("tag")
work.SetExpanded(true)

sidebar.SetSelectedItem(draft)
```

Les méthodes de définition des lignes sont `SetLabel`, `SetSymbol`, `SetTooltip`, `SetEnabled`, `SetHidden`, `SetBadge`, `SetAccessorySymbol`, `SetTintColor`, `SetEditable` et `SetExpanded`. `OnClick` se déclenche quand AppKit sélectionne la ligne, `OnExpandedChange` quand l’utilisateur ouvre ou ferme ses lignes imbriquées, et `OnRename` après la validation d’un renommage sur place.

### Sélection

La sélection unique est activée par défaut. `SetSelectedItem` sélectionne une ligne sans déclencher son `OnClick`. Lorsque la sélection multiple est activée, `OnClick` se déclenche toujours pour la ligne cliquée et `OnSelectionChange` signale l’ensemble de la sélection.

```go
sidebar.SetAllowsMultipleSelection(true)
sidebar.OnSelectionChange(func(_ *application.Context, items []*application.MacSidebarItem) {
    for _, item := range items {
        _ = item.Section()
    }
})
```

### Menus contextuels

Un clic droit recherche un menu dans cet ordre : le `SetContextMenu` propre à la ligne, puis le rappel `OnContextMenu` de la barre latérale, puis le `SetContextMenu` de secours de la barre latérale. Le rappel s’exécute sur le thread de l’application pendant qu’AppKit attend ; gardez-le rapide.

```go
sidebar.OnContextMenu(func(_ *application.Context, item *application.MacSidebarItem) *application.Menu {
    if item == nil {
        return nil // fall back to the menu set with SetContextMenu
    }
    menu := application.NewMenu()
    menu.Add("Delete").OnClick(func(*application.Context) { item.Remove() })
    return menu
})

empty := application.NewMenu()
empty.Add("New Note").OnClick(func(*application.Context) {})
sidebar.SetContextMenu(empty)
```

### Réorganisation par glissement

`SetReorderable` permet à l’utilisateur de faire glisser des lignes au sein d’une section, entre des sections, ainsi que vers ou depuis la racine. Le modèle Go est mis à jour avant le déclenchement de `OnMove`.

```go
sidebar.SetReorderable(true)
sidebar.OnMove(func(_ *application.Context, item *application.MacSidebarItem, section *application.MacSidebarSection, index int) {
    // section is nil when the row was dropped at the sidebar root
})
```

### Suppression

Supprimer une ligne supprime aussi ses lignes imbriquées. Les références deviennent ensuite inertes.

```go
notes.Remove(draft)
sidebar.RemoveSection(tags)
```

## Liste de contenu

`MacContentList` est la colonne centrale de Finder, de Mail et des navigateurs de documents. Sans colonnes, elle affiche des lignes enrichies : un titre, un sous-titre, un symbole au début, un détail à la fin et un badge de comptage.

```go
list := application.NewMacContentList().
    SetStyle(application.MacContentListStyleInset).
    SetEmptyText("No notes match")

row := list.AddRow("Saturday, slowly").
    SetSubtitle("A slow day is still a day well spent.").
    SetDetail("Yesterday").
    SetSymbol("doc.text").
    SetBadge(2)

list.OnSelectionChange(func(_ *application.Context, rows []*application.MacContentListRow) {})
list.OnActivate(func(_ *application.Context, row *application.MacContentListRow) {
    // double-click or Return
})
list.SetSelectedRow(row)
```

Les styles sont `MacContentListStyleAutomatic`, `MacContentListStyleInset`, `MacContentListStyleSourceList`, `MacContentListStylePlain` et `MacContentListStyleFullWidth`. `SetRowHeight`, `SetAlternatingRowBackgrounds`, `SetHeaderVisible` et `SetAllowsMultipleSelection` complètent les options de présentation. Les lignes prennent en charge `SetHidden`, `SetEnabled`, `SetTooltip`, `InsertRow` à une position donnée, `Remove` et `RemoveAll`.

### Colonnes

`SetColumns` passe en mode tableau. Les lignes affichent alors leurs valeurs `SetCells`, une par colonne. Les colonnes marquées `Sortable` affichent un indicateur de tri ; `SetSortable` active le tri par clic sur les en-têtes. Sans rappel `OnSort`, la liste se trie elle-même selon le texte de la colonne avec `SortBy`.

```go
table := application.NewMacContentList().
    SetColumns(
        application.MacContentListColumn{Title: "Name", Width: 220, Sortable: true},
        application.MacContentListColumn{Title: "Size", Width: 80, Alignment: application.MacContentListAlignTrailing},
        application.MacContentListColumn{Title: "Modified", Sortable: true},
    ).
    SetSortable(true).
    SetAlternatingRowBackgrounds(true)

table.AddRow("notes.txt").SetCells("notes.txt", "4 KB", "Today")
table.AddRow("ideas.txt").SetCells("ideas.txt", "1 KB", "Yesterday")

// Without OnSort the list sorts itself by the column text.
table.OnSort(func(_ *application.Context, column int, ascending bool) {
    table.SortRows(func(a, b *application.MacContentListRow) bool {
        return (a.Cells()[column] < b.Cells()[column]) == ascending
    })
})
```

### Menus contextuels

Les menus contextuels sont choisis dans le même ordre que pour la barre latérale : le `SetContextMenu` de la ligne, puis `OnContextMenu`, puis le menu de secours de la liste.

```go
list.OnContextMenu(func(_ *application.Context, row *application.MacContentListRow) *application.Menu {
    if row == nil {
        return nil
    }
    menu := application.NewMenu()
    menu.Add("Delete").OnClick(func(*application.Context) { row.Remove() })
    return menu
})
```

## Inspecteur

`MacInspector` est un panneau de propriétés situé au bord final, constitué de contrôles natifs regroupés en sections. Chaque méthode `Add` renvoie une référence `MacInspectorControl` avec des méthodes de définition et des rappels propres au type de contrôle.

```go
inspector := application.NewMacInspector()

document := inspector.AddSection("Document")
document.AddTextField("Title", "Saturday, slowly").
    OnTextChange(func(_ *application.Context, value string) {})
document.AddPopup("Category", []string{"Personal", "Work"}, 0).
    OnSelectionChange(func(_ *application.Context, index int, value string) {})
document.AddCheckbox("Pinned", false).
    OnToggle(func(_ *application.Context, checked bool) {})

appearance := inspector.AddSection("Appearance").SetCollapsible(true)
priority := appearance.AddSlider("Priority", 0, 5, 0)
priority.OnValueChange(func(_ *application.Context, value float64) {})
appearance.AddStepper("Indent", 0, 8, 1, 2)
appearance.AddSegmented("Align", []string{"Left", "Centre", "Right"}, 0)
appearance.AddColorWell("Tint", application.NewRGB(0, 122, 255)).
    OnColorChange(func(_ *application.Context, colour application.RGBA) {})
appearance.AddDatePicker("Due", time.Time{}).
    OnDateChange(func(_ *application.Context, t time.Time) {})
appearance.AddButton("Reset").OnClick(func(*application.Context) {
    priority.SetFloatValue(0)
})

statistics := inspector.AddSection("Statistics")
words := statistics.AddLabel("Words", "0")

// Programmatic setters never fire the callbacks.
words.SetValue("128")
```

Types de contrôles et méthodes de définition correspondantes :

| Contrôle | Méthodes de définition | Rappel |
|---------|---------|----------|
| `AddLabel` | `SetValue` | aucun |
| `AddTextField` | `SetValue` | `OnTextChange` |
| `AddCheckbox` | `SetChecked` | `OnToggle` |
| `AddPopup`, `AddSegmented` | `SetOptions`, `SetSelectedIndex` | `OnSelectionChange` |
| `AddSlider`, `AddStepper` | `SetFloatValue`, `SetRange`, `SetStep` | `OnValueChange` |
| `AddColorWell` | `SetColor` | `OnColorChange` |
| `AddDatePicker` | `SetDate` | `OnDateChange` |
| `AddButton` | `SetLabel` | `OnClick` |

`SetLabel`, `SetTooltip`, `SetEnabled` et `SetHidden` s’appliquent à tous les types. Les sections peuvent être repliables ; les sections comme les contrôles peuvent être déplacés ou supprimés à tout moment.

```go
appearance.SetCollapsed(true)
inspector.MoveSection(statistics, 0)
appearance.Remove(priority)
inspector.RemoveSection(statistics)
```

## Accessoires

`MacAccessory` est une bande de contrôles natifs. Choisissez une disposition lors de sa création : `MacAccessoryLayoutLeading` se place à côté des boutons de fenêtre, `MacAccessoryLayoutTrailing` au bord final de la barre de titre, `MacAccessoryLayoutBottom` s’étend sur toute la largeur sous la barre de titre et la barre d’outils, et `MacAccessoryLayoutTop` se place en haut d’un volet de vue fractionnée. Ajoutez tous les contrôles avant d’attacher l’accessoire.

```go
// Next to the window buttons.
folders := application.NewMacAccessory(application.MacAccessoryLayoutLeading)
folders.AddSegmented([]string{"Inbox", "Starred"}, 0).
    SetSegmentSymbols("tray", "star").
    OnSelectionChange(func(_ *application.Context, index int, label string) {})

// At the trailing edge of the titlebar.
tools := application.NewMacAccessory(application.MacAccessoryLayoutTrailing)
tools.AddSearch("Search mail").
    SetIncremental(true).
    SetWidth(200).
    OnSearch(func(_ *application.Context, query string) {})
tools.AddSymbolButton("square.and.pencil").
    SetTooltip("Compose").
    OnClick(func(*application.Context) {})

// A full-width strip beneath the titlebar and toolbar.
status := application.NewMacAccessory(application.MacAccessoryLayoutBottom).
    SetHeight(30).
    SetPreferredScrollEdgeEffectStyle(application.MacScrollEdgeEffectStyleSoft)
label := status.AddLabel("Saved").SetSymbol("checkmark.circle")
status.AddFlexibleSpace()
status.AddButton("Mark All Read").OnClick(func(*application.Context) {})

for _, accessory := range []*application.MacAccessory{folders, tools, status} {
    if err := window.AddTitlebarAccessory(accessory); err != nil {
        window.Error("titlebar accessory: %s", err)
    }
}

// Live updates and lifecycle.
label.SetText("Edited").SetSymbol("pencil.circle")
status.SetHidden(true)
tools.Remove()
```

Les contrôles sont `AddSearch`, `AddSegmented`, `AddButton`, `AddSymbolButton`, `AddMenuButton`, `AddLabel`, `AddFlexibleSpace` et, pour les intégrations natives, `AddNativeView`. `AddTitlebarAccessory` existe sur `WebviewWindow` comme sur `NativeWindow` ; les appels effectués avant la création de la fenêtre sont mis en attente et appliqués lors de sa création. `Remove` détache un accessoire pour permettre de l’attacher ailleurs, et `SetHidden` le replie sur place.

### Accessoires de volet

Sur macOS 26 et les versions ultérieures, un accessoire peut se placer en haut ou en bas d’un volet fractionné, comme le champ de filtre de la barre latérale de Finder. Créez l’accessoire avec la disposition `Top` ou `Bottom`, puis attachez-le au volet avec `AddTopAccessory` ou `AddBottomAccessory`.

```go
filter := application.NewMacAccessory(application.MacAccessoryLayoutTop)
filter.AddSearch("Filter").
    SetIncremental(true).
    OnSearch(func(_ *application.Context, query string) { filterNotes(query) })

if err := sidebarPane.AddTopAccessory(filter); err != nil {
    window.Error("sidebar filter: %s", err)
}
```

`SetPreferredScrollEdgeEffectStyle` choisit le rendu `Automatic`, `Soft` ou `Hard` d’AppKit pour le contenu qui défile derrière l’accessoire. Les styles explicites nécessitent macOS 26.1 ; sur les versions antérieures, la demande est signalée par le gestionnaire d’erreurs de la fenêtre et le style automatique reste appliqué.

## Assemblage

Ce programme crée une fenêtre à trois volets avec une barre latérale, la WebView et un inspecteur, ainsi qu’une barre d’outils qui suit les deux séparateurs.

```go title="main.go"
package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:   "Notes",
		Assets: application.AssetOptions{Handler: application.AssetFileServerFS(assets)},
	})

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Notes",
		Width:  1100,
		Height: 700,
		URL:    "/",
		Mac: application.MacWindow{
			ContentLayout: application.MacContentLayoutEdgeToEdge,
			TitleBar: application.MacTitleBar{
				FullSizeContent:      true,
				HideToolbarSeparator: true,
				ToolbarStyle:         application.MacToolbarStyleUnified,
			},
		},
	})

	sidebar := application.NewMacSidebar()
	notes := sidebar.AddSection("Notes")
	notes.AddItem("Saturday, slowly").SetSymbol("doc.text").OnClick(func(*application.Context) {
		app.Event.Emit("note:selected", "saturday")
	})

	inspector := application.NewMacInspector()
	inspector.AddSection("Document").AddTextField("Title", "Saturday, slowly").
		OnTextChange(func(_ *application.Context, value string) {
			app.Event.Emit("note:title", value)
		})

	split := application.NewMacSplitView().SetAutosaveName("notes.main")
	split.AddSidebar(sidebar).SetMinimumThickness(200).SetCollapsible(true)
	split.AddPrimaryContent()
	split.AddInspector(inspector).SetMinimumThickness(240).SetCollapsible(true)
	window.SetSplitView(split)

	toolbar := application.NewMacToolbar().SetDisplayMode(application.MacToolbarDisplayModeIconOnly)
	toolbar.AddSidebarToggle()
	toolbar.AddSidebarTrackingSeparator()
	toolbar.AddButton("New").SetSymbol("square.and.pencil").SetBordered(true).
		OnClick(func(*application.Context) {
			notes.AddItem("Untitled").SetSymbol("doc.text")
		})
	toolbar.AddFlexibleSpace()
	toolbar.AddInspectorTrackingSeparator()
	toolbar.AddInspectorToggle()
	window.SetToolbar(toolbar)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
```

L’habillage natif et l’interface communiquent au moyen des événements et services Wails habituels. La barre latérale émet `note:selected`, l’inspecteur émet `note:title` et la page écoute avec `Events.On` du runtime.

## Fenêtres natives

@note{type="caution" title="Expérimental"}
`NativeWindow`, `NativeWindowManager` et `MacTextEditor` sont expérimentaux dans v3. L’API est volontairement limitée et peut changer lorsque l’API commune des fenêtres sera repensée pour v4.
@end

Un `NativeWindow` ne contient pas de WebView. Son contenu principal est un `MacTextEditor`, c’est-à-dire un `NSTextView` dans un `NSScrollView`, et il accepte les mêmes types de barres d’outils, de vues fractionnées et d’accessoires qu’un `WebviewWindow`. Créez-en un avec `app.NativeWindow.New` ou `app.NativeWindow.NewWithOptions`, puis retrouvez-le avec `Get` ou `GetByID`.

Une fenêtre native n’est créée que lorsqu’elle dispose d’un contenu. Fournissez une vue fractionnée dont le volet principal a été ajouté avec `AddTextEditor`, soit par `NativeWindowOptions.SplitView`, soit en appelant `SetSplitView`. Avant `app.Run`, la disposition est mise en attente ; dans une application en cours d’exécution, `SetSplitView` crée et affiche immédiatement la fenêtre (sauf si `Hidden` est défini) et renvoie toute erreur de création. Une fenêtre sans disposition reste en attente et `Run` renvoie `ErrNativeWindowContentRequired` ; une disposition sans éditeur de texte est rejetée avec `ErrNativeWindowEditorRequired`. Utilisez `NativeWindowOptions.Toolbar` et `NativeWindowOptions.SplitView` plutôt que les champs `Mac.Toolbar` et `Mac.SplitView`, qu’une fenêtre native ignore.

```go title="main.go"
package main

import (
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	app := application.New(application.Options{
		Name:       "Native Notes",
		NativeOnly: true,
	})

	editor := application.NewMacTextEditor()
	editor.OnChange(func(*application.Context) {
		// mark the document dirty; call editor.Text() only when saving
	})

	sidebar := application.NewMacSidebar()
	sidebar.AddSection("Files").AddItem("README.txt").
		SetSymbol("doc.plaintext").
		OnClick(func(*application.Context) {
			editor.SetText("Hello from AppKit")
		})

	split := application.NewMacSplitView().SetAutosaveName("native-notes.main")
	split.AddSidebar(sidebar).SetMinimumThickness(200).SetCollapsible(true)
	split.AddTextEditor(editor).SetMinimumThickness(400)

	window := app.NativeWindow.NewWithOptions(application.NativeWindowOptions{
		Title:  "Native Notes",
		Width:  900,
		Height: 600,
		Mac: application.MacWindow{
			TitleBar: application.MacTitleBar{
				FullSizeContent: true,
				ToolbarStyle:    application.MacToolbarStyleUnified,
			},
		},
	})
	if err := window.SetSplitView(split); err != nil {
		log.Fatal(err)
	}

	toolbar := application.NewMacToolbar()
	toolbar.AddSidebarToggle()
	toolbar.AddSidebarTrackingSeparator()
	toolbar.AddFlexibleSpace()
	toolbar.AddButton("Save").SetSymbol("square.and.arrow.down").SetBordered(true).
		OnClick(func(*application.Context) {
			log.Printf("%d bytes", len(editor.Text()))
		})
	if err := window.SetToolbar(toolbar); err != nil {
		log.Fatal(err)
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
```

La même fenêtre peut être créée en un seul appel depuis une application en cours d’exécution en passant l’habillage comme options :

```go
window := app.NativeWindow.NewWithOptions(application.NativeWindowOptions{
    Title:     "Native Notes",
    SplitView: split,
    Toolbar:   toolbar,
})
```

`MacTextEditor` propose `SetText`, `Text`, `SetEditable`, `OnChange` et `Focus`. Les appels programmatiques à `SetText` ne déclenchent pas `OnChange` : charger un fichier ne le marque donc jamais comme modifié. `Text` lit l’intégralité du document depuis AppKit ; appelez-le lorsque vous avez besoin du contenu, plutôt qu’à chaque modification.

Deux options permettent d’alléger une application exclusivement native :

- `NativeOnly: true` dans `application.Options` désactive le transport de l’interface et le serveur de ressources à l’exécution. Ne créez pas de `WebviewWindow` lorsque cette option est définie.
- Le tag de compilation `wails_native` exclut entièrement du binaire le code de la WebView, de l’interface et du système de mise à jour, et définit automatiquement `NativeOnly` :

```sh
go build -tags wails_native .
```

La prise en charge d’une instance unique est également exclue d’une compilation `wails_native` ; ajoutez le tag `wails_single_instance` si vous en avez besoin.

## Onglets de fenêtre

macOS peut regrouper les fenêtres en onglets. Le mode d’onglets est fixé à la création de la fenêtre : définissez donc `Mac.TabbingMode` sur `MacWindowTabbingModePreferred` ou `MacWindowTabbingModeAutomatic` pour chaque fenêtre qui doit participer. `MacWindowTabbingModeDisallowed` (ainsi que la valeur par défaut non définie) exclut une fenêtre des groupes d’onglets.

```go
first := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Document 1",
    URL:   "/",
    Mac: application.MacWindow{
        TabbingMode: application.MacWindowTabbingModePreferred,
    },
})
```

Les opérations sur les onglets nécessitent des fenêtres natives existantes : appelez-les donc depuis un gestionnaire de menu, une méthode de service ou tout autre code exécuté après le démarrage de `app.Run`. `TabGroup` renvoie une référence `MacWindowTabGroup` qui retrouve le groupe natif à chaque appel ; une référence nil peut être utilisée sans risque et toutes ses méthodes renvoient leur valeur zéro.

```go
second := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Document 2",
    URL:   "/",
    Mac: application.MacWindow{
        TabbingMode: application.MacWindowTabbingModeAutomatic,
    },
})
if err := first.AddTab(second, application.MacTabOrderAbove); err != nil {
    app.Logger.Error("add tab", "error", err)
}
second.SetTabTitle("Draft")

if group := first.TabGroup(); group != nil {
    group.SelectNext()
    group.ToggleTabBar()
    for _, member := range group.Windows() {
        app.Logger.Info("tab", "name", member.Name())
    }
}

first.MoveTabToNewWindow()
first.MergeAllWindows()
```

`MacWindowTabGroup` fournit `Identifier`, `Count`, `Windows`, `NativeWindows`, `SelectedWindow`, `Select`, `SelectNative`, `SelectNext`, `SelectPrevious`, `IsTabBarVisible`, `ToggleTabBar`, `IsOverviewVisible` et `ToggleTabOverview`. `app.Window.TabGroups` répertorie tous les groupes contenant un `WebviewWindow`. `AddNativeTab` ajoute un `NativeWindow` au groupe d’une fenêtre WebView, et `SetTabTooltip` définit le texte affiché au survol d’un onglet. Hors de macOS, les méthodes d’onglets renvoient `ErrMacWindowTabsUnsupported`.

## Versions requises

Tout ce qui figure sur cette page est propre à macOS. Sur les versions antérieures, les fonctionnalités se dégradent progressivement comme indiqué ci-dessous ; l’API Go est identique partout.

| Fonctionnalité | Version minimale de macOS | Comportement sur les versions antérieures |
|---------|---------------|-------------------------------|
| Onglets de fenêtre | 10.12 | indisponibles |
| Groupes de barre d’outils, éléments avec bordure, `AddMenu` | 10.15 | les éléments de menu sont omis ; les groupes utilisent la présentation classique |
| SF Symbols (`SetSymbol` sur les éléments de barre d’outils, de barre latérale, de liste et d’accessoire) | 11 | aucune image n’est affichée |
| `AddSearch` comme `NSSearchToolbarItem`, `SetNavigational`, séparateur de suivi de barre latérale | 11 | la recherche utilise à la place un champ de recherche simple ; le séparateur est omis |
| Rôles de volet fractionné pour l’inspecteur et la liste de contenu | 11 | les mêmes volets sont hébergés dans des éléments fractionnés ordinaires |
| Styles de liste de contenu | 11 | ignorés |
| `SetCenteredItems` | 13 | ignoré |
| Bouton de bascule et séparateur de suivi de l’inspecteur | 14 | Wails fournit un bouton de bascule natif ; le séparateur est omis |
| `SetBadgeCount`, `SetProminent`, `SetTintColor` sur les éléments de barre d’outils | 26 | conservés et appliqués lorsqu’ils sont disponibles |
| Accessoires de volet (`AddTopAccessory`, `AddBottomAccessory`) | 26 | `ErrMacSplitItemAccessoryUnavailable` |
| Styles explicites d’effet de bord de défilement | 26.1 | signalés par le gestionnaire d’erreurs de la fenêtre ; le style automatique reste appliqué |

Les accessoires de barre de titre, les vues fractionnées, les barres latérales, les inspecteurs et l’éditeur de texte n’ont aucune exigence de version au-delà du minimum requis par Wails.

## Exemples

Chaque exemple est une application complète et exécutable :

- [`v3/examples/mac-toolbar`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-toolbar) : un éditeur de notes avec barre d’outils, barre latérale, liste de contenu, inspecteur, fournisseur de partage, personnalisation de la barre d’outils et accessoire de volet.
- [`v3/examples/mac-titlebar-accessory`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-titlebar-accessory) : des accessoires de barre de titre placés au début, à la fin et en bas, qui pilotent une vue de boîte aux lettres.
- [`v3/examples/mac-window-tabs`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-window-tabs) : les modes d’onglets, `AddTab`, `AddNativeTab` et l’API des groupes d’onglets depuis un menu.
- [`v3/examples/mac-native-editor`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-native-editor) : un éditeur de texte `NativeWindow` sans WebView, compilé avec le tag `wails_native`.

## Pour aller plus loin

L’habillage décrit sur cette page constitue une moitié d’une application macOS native. L’autre moitié concerne le comportement de l’application : fenêtres de document avec icônes proxy et disposition en cascade, menus avec symboles et badges, menu Dock avec progression, alertes et panneaux natifs, éléments de barre d’état amovibles, retour haptique et synthèse vocale, presse-papiers enrichi et glissement sortant, ainsi que les informations système sur les autorisations, l’alimentation et les paramètres régionaux. Ces API sont présentées dans le guide [Intégration à la plateforme macOS](/guides/macos-platform-integration).
