---
title: "Système de bindings"
description: "Comment Wails v3 permet à Go et JavaScript de s’appeler mutuellement sans aucun code répétitif"
slug: "contributing/binding-system"
sourcePath: "contributing/binding-system.md"
---

> Les « bindings » constituent le **contrat à typage sûr** qui vous permet d’écrire :

```go
msg, err := chatService.Send("Hello")
```

en Go *et*

```ts
import { Send } from "../bindings/github.com/you/yourapp/services/chatservice";

const msg = await Send("Hello");
```

en TypeScript **sans écrire manuellement le moindre code de liaison IPC**. Ce document explique en détail *comment* cela fonctionne, depuis l’**analyse statique** lors de la compilation, en passant par la **génération de code**, jusqu’au **pont d’exécution** qui transfère les octets via la WebView.

> Consultez [`contributing/architecture/bindings`](/contributing/architecture/bindings/) pour obtenir
>
> l’analyse de référence détaillée du pipeline du générateur — cette page en donne une vue d’ensemble
>
> destinée aux contributeurs.

---

## 1. Vue d’ensemble en 30 secondes

| Étape | Composant | Sortie |
| --- | --- | --- |
| **Collecte/analyse** | `internal/generator/collect/`, `internal/generator/analyse.go` | Modèle en mémoire des services Go exportés, de leurs méthodes, paramètres et types de retour, ainsi que des modèles |
| **Génération** | `internal/generator/render/templates/*.tmpl` (`service.{js,ts}.tmpl`, `models.{js,ts}.tmpl`, `index.tmpl`, `eventcreate.js.tmpl`, `eventdata.d.ts.tmpl`, `newline.tmpl`) | Modules ES propres à chaque service sous `frontend/bindings/<full Go import path>/...` |
| **Exécution** | `pkg/application/messageprocessor*.go` + l’environnement d’exécution JS embarqué sous `internal/runtime/desktop/@wailsio/runtime/src/` (`calls.ts`, `events.ts`, …) | Messages d’appel et d’événement transmis par le pont natif de la WebView |

La commande `wails3 generate bindings` orchestre le flux et pilote `generator.Generate` (défini dans `internal/generator/generate.go`) sur un ensemble de paquets Go.

```
wails3 generate bindings
        │
        ▼
internal/generator/generate.go: Generator.Generate(patterns...)
        │
        ├── internal/generator/collect/   // load.go, collector.go, service.go, model.go, …
        ├── internal/generator/analyse.go // semantic checks
        └── internal/generator/render/    // template execution → frontend/bindings/**
```

---

## 2. Analyse statique

### Point d’entrée

```
internal/generator/generate.go        // Generator + Generate(patterns…)
internal/generator/analyse.go         // semantic validation
internal/generator/collect/load.go    // go/packages loader
internal/generator/collect/collector.go
```

La passe de collecte parcourt chaque paquet chargé et enregistre :

- `collect.ServiceInfo` — un par structure Go exportée et liée.
- `collect.ServiceMethodInfo` / `collect.MethodInfo` — informations de signature propres à chaque méthode (nom, paramètres, résultats, position de l’erreur, receveur, documentation).
- `collect.ModelInfo` / `collect.StructInfo` — générés sous forme de modèles TS/JS.
- Commentaires de directive tels que `//wails:inject`, `//wails:include`, `//wails:internal`, `//wails:ignore`, `//wails:id <hex>` (voir `internal/generator/collect/directive.go`).

Les types non pris en charge provoquent une erreur du générateur afin que les erreurs apparaissent lors de la compilation, et non lors de l’exécution.

### Identifiants de modèle

L’enveloppe d’appel à l’exécution identifie une méthode à l’aide d’un **hachage FNV-1a déterministe** de son nom pleinement qualifié (`pkg.Struct.Method`). Il apparaît sous la forme `$Call.ByID(<numeric-id>, …)` dans les bindings générés, ou sous la forme `$Call.ByName("pkg.Struct.Method", …)` lorsque la génération s’exécute avec `-names`.

---

## 3. Génération de code

### Modèles de génération

`internal/generator/render/templates/` :

| Modèle de génération | Objectif |
| --- | --- |
| `service.js.tmpl` | Un module JS par service lié |
| `service.ts.tmpl` | Fichier TypeScript complémentaire (avec `-ts`) |
| `models.js.tmpl` | Sortie des classes de modèle (par paquet) |
| `models.ts.tmpl` | Sortie `.d.ts` des modèles (par paquet) |
| `index.tmpl` | Réexportations groupées `index.{js,ts}` propres à chaque paquet |
| `eventcreate.js.tmpl` / `eventdata.d.ts.tmpl` | Constructeur d’événement / typage de la charge utile |
| `newline.tmpl` | Normaliseur de saut de ligne final |

La sortie est placée sous `frontend/bindings/<full Go import path>/...`. Par exemple, un service défini dans `github.com/you/yourapp/services/chat` est placé sous `frontend/bindings/github.com/you/yourapp/services/chat/`. Il n’existe aucun répertoire `frontend/src/wailsjs/` dans la v3.

### Sortie JavaScript

Les bindings générés sont des modules ES qui importent les fonctions utilitaires d’exécution depuis `/wails/runtime.js` :

```js
import { Call as $Call, Create as $Create } from "/wails/runtime.js";

/**
 * @param {string} msg
 * @returns {Promise<string> & { cancel(): void }}
 */
export function Send(msg) {
    return $Call.ByID(2042131923, msg);
}
```

Lorsque la génération s’exécute avec `-names`, `$Call.ByName("pkg.Struct.Method", ...)` est généré à la place — toujours **pleinement qualifié**, jamais seulement `"Method"`.

Les classes de modèle générées utilisent un patron de constructeur `$$source` avec des valeurs par défaut `if (!("X" in $$source))` propres à chaque champ, des noms de champs entre guillemets et une `static createFrom(...)` qui exécute `JSON.parse` sur les entrées de type chaîne.

### Principales correspondances de types

Vérifié par rapport à `internal/generator/render/` :

| Go | TypeScript |
| --- | --- |
| `map[string]V` | `{ [_: string]: V }` |
| `map[K]V` (`K` non chaîne) | `{ [_ in K]?: V }` (ni `Map<K, V>` ni `Record<K, V>`) |
| `[]byte` | `Uint8Array` |
| `time.Time` | `string` (chaîne JSON au format ISO 8601) |
| `error` (en position de retour) | promesse rejetée |

### Remarque sur la réflexion

`pkg/application/bindings.go` est **écrit à la main** et utilise `reflect` pour piloter la distribution des méthodes à partir d’un registre `BoundMethod`. Ne prenez pas trop au pied de la lettre les anciennes affirmations selon lesquelles il n’y aurait « aucune réflexion à l’exécution » : le générateur évite la réflexion, mais le distributeur d’exécution l’utilise.

---

## 4. Protocole d’invocation à l’exécution

### Côté JavaScript

```ts
import { Call } from "/wails/runtime.js";

await Call.ByID(0x7a1201d3 /* ChatService.Send */, "Hello");
// or, with -names:
await Call.ByName("chatservice.ChatService.Send", "Hello");
```

Les utilitaires d’exécution se trouvent dans `internal/runtime/desktop/@wailsio/runtime/src/calls.ts` (distribution des appels), `events.ts` (événements) et les fichiers associés : cet arbre ne contient ni `invoke.ts` ni `errors.ts`. L’enveloppe exacte transmise sur le canal est encodée par `calls.ts` côté JS et décodée par `pkg/application/messageprocessor_call.go` côté Go ; consultez ces deux fichiers ensemble lorsque vous déboguez le pont.

### Côté Go

1. `pkg/application/messageprocessor_call.go` reçoit le message d’appel.
2. Recherche la méthode liée par son ID ou son nom dans `pkg/application/bindings.go` (piloté par `reflect`).
3. Appelle la méthode liée et sérialise `{result, error}` pour le renvoyer à JS.

### Correspondance des erreurs

| Go | JavaScript |
| --- | --- |
| `error == nil` | `Promise` est résolue avec le résultat |
| `error != nil` | `Promise` est rejetée avec une `Error` dont `message` contient la chaîne de l’erreur Go |

---

## 5. Appel de JavaScript depuis Go

Le générateur de liaisons fonctionne dans un seul sens (les méthodes Go sont exposées à JS). Pour communiquer de Go vers JS, utilisez le bus d’événements ou exécutez du JS dans une fenêtre :

```go
app.Event.Emit("chat:new-message", msg)
window.ExecJS(`window.dispatchEvent(new CustomEvent("ping"))`)
```

Côté JS, abonnez-vous avec `Events.On(name, cb)` depuis `/wails/runtime.js`.

---

## 6. Extension et dépannage

### Erreur de type non pris en charge

```
error: field "Client" uses unsupported type: chan struct{}
```

→ encapsulez le canal derrière une API de méthodes ou marquez le champ avec `//wails:internal` afin que le générateur l’ignore.

### Liaisons obsolètes

La sortie générée est remplacée à chaque `wails3 generate bindings`, `wails3 dev` ou `wails3 build`. Si la saisie semi-automatique de l’IDE affiche des stubs obsolètes, supprimez `frontend/bindings/` et relancez le générateur. L’option `-clean` (définie par défaut sur `true` dans les versions actuelles) vide le répertoire des liaisons avant chaque exécution.

### Conseils de performance

- Évitez de transmettre en continu de grandes tranches d’octets via le pont : servez-les plutôt au moyen du serveur de ressources.
- Lorsque la latence est un critère important, regroupez plusieurs appels rapides dans une seule méthode.
- Pour les petites structures de paramètres, préférez des récepteurs par valeur afin de réduire les allocations.

---

## 7. Carte des fichiers clés

| Fonction | Fichier |
| --- | --- |
| Orchestration du générateur | `internal/generator/generate.go` |
| Vérifications sémantiques | `internal/generator/analyse.go` |
| Collecte (services, méthodes, modèles) | `internal/generator/collect/{service,method,model,struct,package}.go` |
| Modèles de rendu | `internal/generator/render/templates/*.tmpl` |
| Emplacement des liaisons générées | `frontend/bindings/<full Go import path>/...` |
| Distributeur côté Go | `pkg/application/bindings.go`, `messageprocessor_call.go` |
| Environnement d’exécution JS | `internal/runtime/desktop/@wailsio/runtime/src/{calls,events,index}.ts` |

Gardez cet aide-mémoire à portée de main lorsque vous recherchez un bogue dans le pont.

---

## 8. Récapitulatif

1. **Collecteur** : analyse votre code Go → modèle sémantique en mémoire.
2. **Modèles** : produisent des modules ES pour chaque service ainsi que des fichiers de modèles et d’index pour chaque paquet.
3. **Processeur de messages** : distribue les appels côté Go via le registre des liaisons.
4. **Environnement d’exécution JS** : encapsule le tout dans des promesses idiomatiques avec annulation.

Le tout sans que vous ayez à écrire la moindre ligne de code répétitif pour l’IPC. Voilà le système de liaisons de Wails v3. À vous de jouer et de créer vos liaisons !
