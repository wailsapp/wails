---
title: "Fonctionnement interne du runtime"
description: "Présentation détaillée du démarrage et de l’exécution de Wails v3, ainsi que de ses échanges avec le système d’exploitation"
slug: "contributing/runtime-internals"
sourcePath: "contributing/runtime-internals.md"
---

Le **runtime** est la couche qui transforme des fonctions Go ordinaires en application de bureau multiplateforme. Ce document décrit les composants que vous rencontrerez en parcourant le code source.

---

## 1. Cycle de vie de l’application

| Phase | Chemin du code | Déroulement |
| --- | --- | --- |
| **Amorçage** | `pkg/application/application.go:init()` | Enregistre les données générées lors de la compilation et crée un singleton `application` global. |
| **New()** | `application.New(...)` | Valide `Options`, démarre l’**AssetServer** et initialise la journalisation. |
| **Run()** | `application.(*App).Run()` | 1. Appelle la fonction `mainthread.X()` propre à la plateforme pour entrer dans le thread d’interface utilisateur du système d’exploitation.<br />2. Démarre le **runtime** (`internal/runtime`).<br />3. Bloque l’exécution jusqu’à la fermeture de la dernière fenêtre ou jusqu’à l’appel de `Quit()`. |
| **Arrêt** | `application.(*App).Quit()` | Diffuse l’événement `application:shutdown`, écrit les données de journalisation encore en mémoire tampon vers leur destination et détruit les fenêtres et les services. |

Le cycle de vie est strictement à **entrée unique** : vous pouvez créer de nombreuses fenêtres, mais l’objet application lui-même n’est initialisé qu’une seule fois.

---

## 2. Gestion des fenêtres

### API publique

```go
win := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "Dashboard",
    Width:  1280,
    Height: 720,
})
win.Show()
```

> `app.Window.New()` ne prend **aucun argument** ; utilisez `NewWithOptions(...)` lorsque vous
>
> devez transmettre une structure `application.WebviewWindowOptions` (par valeur).

`app.Window.New[WithOptions]()` délègue à `pkg/application/webview_window_*.go`, qui contient les implémentations propres à chaque plateforme :

```
pkg/application/
├── webview_window_darwin.go    // WKWebView
├── webview_window_linux.go     // GTK + WebKitGTK (plus linux_cgo*.go)
└── webview_window_windows.go   // WebView2
```

Chaque fichier :

1. Crée une vue web native (WKWebView, WebKitGTK, WebView2).
2. Enregistre une fonction de rappel du **processeur de messages** (`pkg/application/messageprocessor*.go`).
3. Associe les événements Wails (`WindowDidResize`, `WindowFocus`, `WindowFilesDropped`, …) aux constantes définies dans `pkg/events`.

`internal/runtime/` est réservé au petit code de liaison reposant sur des balises de compilation (`runtime{,_darwin,_linux,_windows,_android,_dev,_prod}.go`) et au runtime JS intégré sous `internal/runtime/desktop/`.

Les fenêtres actives sont suivies par `pkg/application/window_manager.go` / `webview_window.go`. `pkg/application/screenmanager.go` contient les métadonnées d’**affichage** (résolution, échelle, zone de travail) — il ne gère pas les fenêtres.

---

## 3. Pipeline de traitement des messages

Le pont entre JavaScript et Go est implémenté par la famille de **processeurs de messages** dans `pkg/application/messageprocessor_*.go`.

Flux :

1. **JavaScript** appelle `Call.ByID(<fnv-id>, ...args)` depuis `/wails/runtime.js` (implémenté dans `internal/runtime/desktop/@wailsio/runtime/src/calls.ts`), ou `Call.ByName("pkg.Struct.Method", ...args)` pour les builds en mode nom.
2. L’utilitaire du runtime encapsule l’appel et le transmet à Go par l’intermédiaire du pont natif propre à chaque plateforme.
3. **Go** reçoit le message dans `pkg/application/messageprocessor_call.go`.
4. Le processeur recherche la méthode liée dans `pkg/application/bindings.go` (écrit manuellement et basé sur `reflect`), puis l’appelle.
5. Le résultat ou l’erreur est sérialisé et renvoyé à JS, où une `Promise` est résolue ou rejetée.

> L’enveloppe JSON exacte est définie par l’utilitaire du runtime côté JS et
>
> par `messageprocessor_call.go` côté Go — les anciennes versions de cette page
>
> indiquaient une structure `{"t":"c","id":"123","m":"Greet","p":[…]}`, mais celle-ci ne
>
> correspond pas à l’implémentation actuelle. Lisez les deux fichiers conjointement lorsque vous recherchez un
>
> bogue lié au format des données échangées.

Processeurs spécialisés :

| Fichier | Rôle |
| --- | --- |
| `messageprocessor_window.go` | Actions sur les fenêtres (masquer, agrandir, …) |
| `messageprocessor_dialog.go` | Boîtes de dialogue natives (`OpenFile`, `MessageBox`, …) |
| `messageprocessor_clipboard.go` | Lecture/écriture du presse-papiers |
| `messageprocessor_events.go` | Abonnement aux événements / émission d’événements |
| `messageprocessor_browser.go` | Navigation dans le navigateur, outils de développement |

Les processeurs sont **sans état** : ils obtiennent tout ce dont ils ont besoin depuis le `ApplicationContext` transmis avec chaque message.

---

## 4. Système d’événements

Les événements sont des chaînes réparties dans des espaces de noms et distribuées sur trois couches :

1. **Événements d’application** : cycle de vie global (`application:ready`, `application:shutdown`).
2. **Événements de fenêtre** : propres à chaque fenêtre (`window:focus`, `window:resize`).
3. **Événements personnalisés** : définis par l’utilisateur (`chat:new-message`).

Détails d’implémentation :

- Les constantes d’événement se trouvent dans `pkg/events/` (`defaults.go`, `known_events.go`, `events.txt`). Elles sont générées par `v3/tasks/events/generate.go` et exposées sous les noms `events.Common.*`, `events.Mac.*`, `events.Windows.*` et `events.Linux.*`. Vous pouvez les régénérer avec `wails3 generate constants`.
- Côté Go (événements d’application) :
  ```go
  app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {})
  ```

- Côté Go (événements de fenêtre) :
  ```go
  window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {})
  ```

- Côté Go (événements personnalisés) :
  ```go
  app.Event.On("chat:new-message", func(e *application.CustomEvent) {})
  ```

- Côté JS :
  ```js
  import { Events } from "/wails/runtime.js";
  Events.On("chat:new-message", (e) => { /* … */ });
  ```


Les événements d’application, de fenêtre et personnalisés transitent tous par `pkg/application/event_manager.go`. Les abonnements aux événements d’une fenêtre sont limités à celle-ci ; sa fermeture désinscrit donc automatiquement les gestionnaires correspondants.

---

## 5. Implémentations propres aux plateformes

La compilation conditionnelle conserve une API publique identique tout en masquant les particularités des systèmes d’exploitation.

| Aspect | Darwin | Linux | Windows |
| --- | --- | --- | --- |
| Thread principal | `mainthread_darwin.go` (Cgo vers Foundation) | `mainthread_linux.go` (GTK) | `mainthread_windows.go` (Win32 `AttachThreadInput`) |
| Boîtes de dialogue | `dialogs_darwin.*` (NSAlert) | `dialogs_linux.go` (GtkFileChooser) | `dialogs_windows.go` (IFileOpenDialog) |
| Presse-papiers | `clipboard_darwin.go` | `clipboard_linux.go` | `clipboard_windows.go` |
| Icônes de zone de notification | `systemtray_darwin.*` | `systemtray_linux.go` (DBus) | `systemtray_windows.go` (Shell_NotifyIcon) |

Principes clés :

- **macOS et Windows** utilisent Cgo avec parcimonie (principalement par l’intermédiaire de `pkg/mac/` et des wrappers Win32 `w32` dans `pkg/w32`).
- **Linux fait largement appel à Cgo** par nécessité : `pkg/application/linux_cgo.go` (environ 69 Ko), ainsi que `linux_cgo_gtk4.{c,go,h}` (environ 50 Ko ou plus), pilotent directement GTK/WebKitGTK.
- Utilisez les **contraintes de compilation** (`//go:build darwin`, `//go:build linux`, …) pour préserver la lisibilité des fichiers propres à chaque système d’exploitation.
- `internal/capabilities/` fournit des indicateurs de capacité propres à chaque plateforme, mais le framework n’exporte **pas** de sentinelle `ErrCapability` : la disponibilité des fonctionnalités est gérée par les valeurs de retour des stubs propres aux plateformes.

---

## 6. Guide des fichiers

| Fichier | Pourquoi le modifier |
| --- | --- |
| `internal/runtime/runtime_*.go` | Modifier la petite couche de stubs fondée sur les contraintes de compilation (développement ou production, code d’intégration propre au système d’exploitation). |
| `pkg/application/webview_window_*.go` | Implémenter une nouvelle indication ou un nouveau comportement de fenêtre. |
| `pkg/application/messageprocessor*.go` | Ajouter une nouvelle commande de pont appelable depuis JS. |
| `pkg/events/*.go` | Étendre les définitions d’événements intégrées (puis réexécuter `wails3 generate constants`). |
| `internal/assetserver/*` | Ajuster la gestion des ressources en développement et en production. |
| `internal/runtime/desktop/@wailsio/runtime/src/*` | Modifier le runtime JS intégré (distribution des appels et des événements, boîtes de dialogue, glisser-déposer, etc.). |

---

## 7. Conseils de débogage

- Configurez `Options.LogLevel` (par exemple `slog.LevelDebug`) et examinez la sortie de `Options.Logger` : il n’existe aucune variable d’environnement `WAILS_LOG_LEVEL`.
- Les indicateurs de `wails3 dev` sont `--config`, `--port` et `-s` (qui active HTTPS), auxquels s’ajoute l’indicateur global `--no-colour`. Il n’existe aucun indicateur `-verbose`.
- Sur macOS, exécutez l’application sous `lldb --` afin de détecter rapidement les exceptions Objective-C.
- Pour les problèmes liés à Chromium sous Windows, activez les journaux de débogage de WebView2 : `set WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS=--remote-debugging-port=9222`

---

## 8. Extension du runtime

1. Facultatif : déclarez tout nouvel indicateur de capacité dans `internal/capabilities/`.
2. Implémentez la fonctionnalité dans chaque variante de `pkg/application/*_{darwin,linux,windows}.go` à l’aide de contraintes de compilation. Fournissez un stub sur les plateformes qui ne la prennent pas en charge.
3. Ajoutez l’API publique dans `pkg/application` (interface et méthodes concrètes de `WebviewWindow`, structure d’options, etc.).
4. Si JS doit l’appeler, enregistrez une nouvelle méthode du processeur de messages (`pkg/application/messageprocessor*.go`) et ajoutez une fonction auxiliaire correspondante au runtime JS.
5. Si vous ajoutez un événement, déclarez sa constante dans `pkg/events/` et exécutez `wails3 generate constants` pour actualiser les fichiers générés.

Suivez cette liste de contrôle pour préserver le contrat multiplateforme.

---

## 9. Glisser-déposer

Sur toutes les plateformes, le glisser-déposer de fichiers suit une approche **centrée sur JavaScript**. La couche native intercepte les événements de glissement du système d’exploitation, mais le traitement effectif du dépôt et l’interaction avec le DOM s’effectuent en JavaScript.

### Déroulement

1. L’utilisateur fait glisser des fichiers depuis le système d’exploitation au-dessus de la fenêtre Wails
2. La couche native détecte le glissement et avertit JavaScript pour les effets de survol
3. L’utilisateur dépose les fichiers
4. La couche native envoie à JavaScript les chemins des fichiers et les coordonnées
5. JavaScript trouve l’élément cible du dépôt (`data-file-drop-target`)
6. JavaScript envoie au backend Go les chemins des fichiers et les informations sur l’élément
7. Go émet l’événement `WindowFilesDropped` avec le contexte complet

### Implémentations propres aux plateformes

| Plateforme | Couche native | Difficulté principale |
| --- | --- | --- |
| **Windows** | Prise en charge intégrée du glissement par WebView2 | Coordonnées en pixels CSS, aucune conversion nécessaire |
| **macOS** | Délégués de glissement de NSWindow | Conversion des coordonnées relatives à la fenêtre en coordonnées relatives à la vue web |
| **Linux** | Signaux de glissement GTK3 | Il faut distinguer les glissements de fichiers des glissements HTML5 internes |

### Linux : distinguer les types de glissement

GTK et WebKit veulent tous deux traiter les événements de glissement. L’essentiel est de vérifier le type de cible du glissement :

```c
static gboolean is_file_drag(GdkDragContext *context) {
    GList *targets = gdk_drag_context_list_targets(context);
    for (GList *l = targets; l != NULL; l = l->next) {
        GdkAtom atom = GDK_POINTER_TO_ATOM(l->data);
        gchar *name = gdk_atom_name(atom);
        if (name && g_strcmp0(name, "text/uri-list") == 0) {
            g_free(name);
            return TRUE;  // External file drag
        }
        g_free(name);
    }
    return FALSE;  // Internal HTML5 drag
}
```

Les gestionnaires de signaux renvoient `FALSE` pour les glissements internes, afin de laisser WebKit les traiter, et `TRUE` pour les glissements de fichiers, afin de les traiter nous-mêmes.

### Blocage du dépôt de fichiers

Lorsque `EnableFileDrop` vaut `false`, nous devons tout de même empêcher le navigateur de naviguer vers les fichiers déposés. Chaque plateforme procède différemment :

- **Windows** : JavaScript appelle `preventDefault()` lors des événements de glissement
- **macOS** : JavaScript appelle `preventDefault()` lors des événements de glissement\
- **Linux** : les gestionnaires de signaux GTK interceptent et rejettent les glissements de fichiers au niveau natif

### Fichiers principaux

| Fichier | Rôle |
| --- | --- |
| `pkg/application/linux_cgo.go` | Gestionnaires des signaux de glissement GTK (code C dans le préambule cgo) |
| `pkg/application/webview_window_darwin.go` | Délégués de glissement de macOS |
| `pkg/application/webview_window_windows.go` | Traitement des messages WebView2 |
| `internal/runtime/desktop/@wailsio/runtime/src/window.ts` | Traitement du dépôt en JavaScript |

### Débogage

- **Linux** : ajoutez `printf` dans le code C (n’oubliez pas `fflush(stdout)`)
- **Windows** : utilisez `globalApplication.debug()`
- **JavaScript** : consultez la console du navigateur et activez le mode de débogage

Problèmes courants :

1. **Le glissement HTML5 interne ne fonctionne pas** : le gestionnaire natif l’intercepte (renvoyez `FALSE` pour les glissements qui ne concernent pas des fichiers)
2. **Les effets de survol ne s’affichent pas** : les gestionnaires JavaScript ne sont pas appelés
3. **Coordonnées incorrectes** : vérifiez les conversions entre espaces de coordonnées

---

Vous disposez maintenant d’une visite guidée des mécanismes internes du runtime. Associez ces connaissances à la carte **Structure du code source** et à la documentation du **serveur de ressources** pour vous repérer avec assurance et apporter des contributions significatives. Bon développement !
