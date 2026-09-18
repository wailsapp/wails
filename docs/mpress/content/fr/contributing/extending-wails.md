---
title: "Étendre Wails"
description: "Guide pratique pour ajouter de nouvelles fonctionnalités et plateformes à Wails v3"
slug: "contributing/extending-wails"
sourcePath: "contributing/extending-wails.md"
---

> Wails est conçu pour être **facilement modifiable**.
>
> Chaque sous-système majeur réside dans du code Go que vous pouvez lire, modifier et distribuer.
>
> Cette page indique *par où* commencer et *comment* préserver la compatibilité multiplateforme lorsque vous :

- Ajoutez un **service** (notifications, stockage clé-valeur, IPC personnalisé, etc.)
- Créez une **nouvelle commande de CLI** (`wails3 <foo>`)
- Étendez le **runtime** (API de fenêtre, boîtes de dialogue, événements)
- Introduisez une **capacité de plateforme** (Wayland, etc.)
- Préservez la **compatibilité multiplateforme** sans vous noyer dans les balises `//go:build`

---

## 1. Ajouter un service

Dans v3, un « service » est un type Go fourni par l’utilisateur, enregistré au moyen de `application.Options.Services` et exposé à JS par des liaisons générées. La base de code de v3 fournit :

- `internal/service/` — structure de départ pour `wails3 generate service` :
  ```
  internal/service/
  ├── service.go              # Install(options *flags.ServiceInit)
  └── template/
      ├── README.tmpl.md
      ├── go.mod.tmpl
      ├── service.go.tmpl
      └── service.tmpl.yml
  ```

- `pkg/services/` — services prêts à l’emploi que vous pouvez enregistrer dès maintenant (notifications, kvstore, sqlite, log, fileserver, dock, etc.).

Les fichiers du générateur et de la CLI que d’anciennes versions désignaient comme `internal/service/template/template.go` et `internal/generator/collect/services.go` n’existent pas : l’outil de génération de structure est `internal/service/service.go` (point d’entrée : `service.Install`), et les métadonnées de liaison des services sont collectées dans `internal/generator/collect/service.go`.

### 1.1 Définir le service

```go
package chat

type Service struct {
    messages []string
}

func New() *Service { return &Service{} }

func (s *Service) Send(msg string) string {
    s.messages = append(s.messages, msg)
    return "ok"
}
```

### 1.2 Implémenter les interfaces de cycle de vie (facultatif)

Un service peut, de manière facultative, satisfaire aux interfaces suivantes (définies dans `pkg/application`) :

```go
func (s *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error { return nil }
func (s *Service) ServiceShutdown() error                                                       { return nil }
```

> **Important :** `ServiceShutdown` ne prend **aucun argument**. Une méthode dont la
>
> signature est `ServiceShutdown(ctx context.Context) error` ne satisfait **pas**
>
> à l’interface et ne sera donc jamais appelée, sans qu’aucune erreur soit signalée.

### 1.3 Enregistrer le service auprès de l’application

Il n’existe aucun appel global `services.Register(...)`. Les services sont enregistrés à l’exécution au moyen de `application.Options.Services` :

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(chat.New()),
    },
})
```

Une fois le service enregistré, `wails3 generate bindings` génère sous `frontend/bindings/<your import path>/...` des modules ES qui encapsulent les méthodes exportées.

### 1.4 Appeler depuis JS

```js
import { Send } from "../bindings/github.com/you/yourapp/chat";

await Send("hi");
```

Il n’existe aucune variable globale `window.backend.*` dans v3 : les appels passent par les modules ES générés, qui appellent à leur tour `Call.ByID(...)` depuis `/wails/runtime.js`.

---

## 2. Écrire une nouvelle commande de CLI

La CLI de v3 utilise **`github.com/leaanthony/clir`** (et non cobra). Son câblage se trouve dans `v3/cmd/wails3/main.go` :

```go
import "github.com/leaanthony/clir"

func main() {
    app := clir.NewCli("wails", "The Wails3 CLI", "v3")
    app.NewSubCommand("hello", "Prints Hello Wails").Action(func() error {
        fmt.Println("Hello Wails!")
        return nil
    })
    // ... other subcommands explicitly wired here
    _ = app.Run()
}
```

Il n’existe aucun enregistrement automatique fondé sur `init()`. Ajoutez votre nouvelle sous-commande à `cmd/wails3/main.go`, ainsi qu’une fonction sous-jacente dans `internal/commands/` (et une structure d’options sous `internal/flags/` si elle accepte des options). Recompilez la CLI :

```
cd v3
go install ./cmd/wails3
wails3 hello
```

Si votre commande nécessite l’intégration au Taskfile, réutilisez les fonctions auxiliaires de `internal/commands/task_wrapper.go` (`wrapTask("yourtask", args)`).

---

## 3. Modifier le runtime

Motifs courants :

- Nouvelle fonctionnalité de fenêtre (`SetOpacity`, `Shake`, etc.)
- Boîte de dialogue supplémentaire (`ColorPicker`)
- API au niveau du système (luminosité de l’écran)

### 3.1 API publique

Ajoutez la méthode à `pkg/application/webview_window.go` (l’interface se trouve dans `window.go`) :

```go
func (w *WebviewWindow) SetOpacity(o float32) Window {
    InvokeSync(func() { w.impl.setOpacity(o) })
    return w
}
```

Utilisez les fonctions auxiliaires `InvokeSync`/`InvokeAsync` existantes pour garantir que l’appel s’exécute sur le thread principal.

### 3.2 Processeur de messages

Si JS doit appeler la nouvelle méthode, étendez le fichier `pkg/application/messageprocessor_*.go` approprié. Le processeur de messages utilise des méthodes basées sur switch dans `MessageProcessor`, et non un appel global `register(...)` :

```go
// inside messageprocessor_window.go
case "setOpacity":
    var args struct {
        WindowID uint    `json:"windowID"`
        Opacity  float32 `json:"opacity"`
    }
    if err := json.Unmarshal(payload, &args); err != nil { ... }
    window, _ := m.app.Window.GetByID(args.WindowID)
    window.SetOpacity(args.Opacity)
```

La base de code ne contient **aucun** fichier `messageprocessor_window_opacity.go` ni modèle `register(MsgSetOpacity, ...)` fondé sur `init()`.

### 3.3 Implémentation propre à la plateforme

Ajoutez l’implémentation à chaque fichier propre à un système d’exploitation sous `pkg/application/` :

```
pkg/application/
├── webview_window_darwin.go   //go:build darwin
├── webview_window_linux.go    //go:build linux
└── webview_window_windows.go  //go:build windows
```

Si une plateforme ne peut pas prendre en charge la fonctionnalité, écrivez un stub sans effet. Le framework ne possède aucune sentinelle `ErrCapability` : indiquez la prise en charge dans la documentation et, si nécessaire, au moyen du champ booléen approprié de `Options` ou de la structure d’options propre à la plateforme.

### 3.4 Indicateur de capacité (facultatif)

Le paquet `internal/capabilities/` sert à déclarer des ensembles de capacités propres à chaque plateforme. Il n’existe aucune API publique `application.HasCapability` / `application.CapOpacity`. Si vous souhaitez pouvoir vérifier une capacité à l’exécution, ajoutez-la sous `internal/capabilities/` et exposez un accesseur typé depuis `pkg/application`.

---

## 4. Ajouter de nouvelles capacités de plateforme

Exemple : prise en charge facultative de Wayland sous Linux.

1. Scindez le fichier `pkg/application/*_linux.go` concerné en `*_linux_x11.go` (`//go:build linux && !wayland`) et `*_linux_wayland.go` (`//go:build linux && wayland`).
2. Permettez aux utilisateurs d’activer explicitement la fonctionnalité avec `wails3 build --tags wayland`. Transmettez les balises supplémentaires par le câblage `EXTRA_TAGS` existant dans `internal/commands/task_wrapper.go`. Il n’existe aucun indicateur `--tags wayland` au niveau de `dev` : `wails3 dev` accepte uniquement `--config`, `--port` et `-s`.
3. Mettez à jour la documentation et tout fichier README propre à une plateforme sous `pkg/application/`.

> Réduisez au minimum les balises de compilation par défaut ; réservez les balises activables aux fonctionnalités spécialisées.

---

## 5. Liste de contrôle de compatibilité multiplateforme

| ✅ Étape | Pourquoi |
| --- | --- |
| Fournissez **toutes** les méthodes publiques dans les fichiers de chaque plateforme (même sous forme de stubs) | Garantit la réussite de la compilation sur chaque système d’exploitation |
| Documentez la dégradation contrôlée pour chaque système d’exploitation | Les applications peuvent effectuer un branchement selon `runtime.GOOS` sans erreurs masquées |
| Utilisez d’abord du **Go pur** et Cgo uniquement lorsque c’est nécessaire | Simplifie la compilation croisée (Linux supporte déjà le coût de Cgo) |
| Exécutez `task test:cli`, `task test:generator` et `task test:templates` | Reproduit localement l’intégration continue |
| Documentez les nouvelles balises de compilation dans la documentation destinée aux contributeurs ou dans le README du modèle | Les utilisateurs doivent connaître les fonctionnalités optionnelles à activer |

---

## 6. Compilations de débogage et vitesse d’itération

- Utilisez `Options.LogLevel = slog.LevelDebug` (`Options.Logger = slog.Default()`) pour afficher des informations détaillées sur l’activité du runtime. Il n’existe aucune variable d’environnement `WAILS_LOG_LEVEL`.
- Les options de `wails3 dev` sont `--config`, `--port` et `-s`. Il n’existe aucune option `-race` ni `-verbose` : utilisez le détecteur de situations de concurrence avec `go test -race ./...`, ou en exécutant `go build -race` sur votre application puis en lançant celle-ci directement.
- Le guide de test des conditions de concurrence et de Cgo se trouve dans `v3/TESTING.md` (d’anciennes versions préliminaires renvoyaient vers `pkg/application/RACE.md`, qui n’existe pas).

---

## 7. Contributions au projet en amont

1. Pour une nouvelle fonctionnalité ou une modification du comportement public, ouvrez une PR préliminaire de **WEP (Wails Enhancement Proposal)** afin de discuter de l’idée et de la conception. N’utilisez une issue que pour un bogue reproductible ou un problème de documentation.
2. Pour l’implémentation, suivez les techniques décrites ci-dessus.
3. Ajoutez :
  - Des tests unitaires (`*_test.go`)
  - De la documentation (ce fichier ou la page `docs/...` appropriée)
  - Un test de non-régression sous `internal/generator/testcases/` si vous avez modifié le générateur de bindings

4. Avant de pousser vos modifications, exécutez localement `task precommit` ainsi que les cibles `task test:*` appropriées.

---

### Liens rapides

| Domaine | Emplacement |
| --- | --- |
| Services intégrés | `pkg/services/` |
| Générateur de structure de service | `internal/service/` |
| Câblage de la CLI | `v3/cmd/wails3/main.go` |
| Corps des commandes de la CLI | `internal/commands/` |
| Runtime propre à chaque système d’exploitation | `pkg/application/*_{darwin,linux,windows}.go` |
| Déclarations de capacités | `internal/capabilities/` |
| DSL du Taskfile | `v3/Taskfile.yaml` |
| Générateur de constantes d’événements | `v3/tasks/events/generate.go` |

---

Vous disposez désormais d’une **feuille de route** pour façonner Wails à votre guise : ajoutez des services, agrémentez la CLI d’un peu de magie, modifiez le runtime ou apportez des fonctionnalités entièrement nouvelles aux systèmes d’exploitation. Bonne extension !
