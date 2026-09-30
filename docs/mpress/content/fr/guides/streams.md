---
title: "Flux"
description: "Flux d’octets bidirectionnels entre Go et JavaScript, avec le modèle de programmation WebSocket et sans socket d’écoute"
slug: "guides/streams"
sourcePath: "guides/streams.md"
---

Les flux fournissent un canal d’octets nommé, ordonné et bidirectionnel entre Go et votre interface, avec le même modèle de programmation qu’un WebSocket — **sans lier de port TCP**.

Le protocole WebSocket ne peut pas transiter par un schéma d’URL personnalisé. Le seul moyen de l’utiliser dans une webview consiste donc à exécuter un véritable serveur HTTP et à écouter sur un port. Dans une application de bureau, cela implique un port local ouvert auquel tout autre processus de la machine peut accéder, qui nécessite une vérification de l’origine et un jeton pour être sécurisé, et que détectent tous les pare-feu et produits de sécurité des terminaux utilisés par vos utilisateurs. Les flux évitent tout cela : ils transitent par le serveur de ressources que votre application utilise déjà et qui est déjà lié à l’origine.

Vous migrez une implémentation WebSocket existante ? Suivez [Migration d’un WebSocket vers les flux](/guides/streams-from-websockets/). Ce guide est conçu pour être appliqué mécaniquement et présente d’abord les trois différences qui provoquent des dysfonctionnements silencieux.

## Démarrage rapide

Déclarez un flux en Go. Le gestionnaire s’exécute une fois par connexion, dans sa propre goroutine :

```go
app.HandleStream("telemetry", func(c *application.StreamConn) {
    defer c.Close()

    for {
        frame, err := c.Receive()   // blocks until a frame arrives
        if err != nil {
            return                  // page reloaded, window closed, or app shutting down
        }
        _ = c.Send(process(frame))  // blocks like a socket write
    }
})
```

Depuis l’interface, connectez-vous au flux par son nom. L’objet implémente l’interface `WebSocket` :

```js
import { Stream } from "@wailsio/runtime";

const s = Stream("telemetry");
s.onopen    = () => s.send(new TextEncoder().encode("hello"));
s.onmessage = (ev) => console.log(new Uint8Array(ev.data));
s.onclose   = (ev) => console.log("closed", ev.code);
```

`Stream(name)` renvoie **de manière synchrone** avec `readyState === CONNECTING`, exactement comme `new WebSocket(url)`. Vous pouvez donc en créer un dans la portée du module :

```js
export const Telemetry = Stream("telemetry");
```

## Les trames sont des octets

Chaque trame est un `[]byte` en Go et un `ArrayBuffer` en JavaScript. Aucun schéma ni encodage ne vous est imposé : utilisez JSON, protobuf, CBOR ou des octets bruts, selon vos préférences.

Une trame est un **message, pas un flux d’octets** : elle arrive en entier ou pas du tout, et sa longueur l’accompagne. Aucun des deux côtés n’a besoin de connaître sa taille à l’avance. Ainsi, une structure dotée d’un champ `[]byte` est sérialisée selon son format habituel et envoyée dans une seule trame.

## Envoi d’objets

Les trames sont des octets, mais il est rarement souhaitable de raisonner en octets. Les deux côtés proposent des fonctions pratiques pour JSON qui fonctionnent de concert :

```go
type Reading struct {
    Sensor string  `json:"sensor"`
    Value  float64 `json:"value"`
}

app.HandleStream("telemetry", func(c *application.StreamConn) {
    defer c.Close()

    var cmd map[string]any
    if err := c.ReceiveJSON(&cmd); err != nil {
        return
    }

    _ = c.SendJSON(Reading{Sensor: "cpu", Value: 42.5})
})
```

```js
import { JSONStream } from "@wailsio/runtime";

const s = JSONStream("telemetry");
s.onopen    = () => s.send({ subscribe: "cpu" });   // stringified for you
s.onmessage = (ev) => console.log(ev.data.value);   // already an object
```

`JSONStream` est le même objet que `Stream`, l’encodage étant effectué à la frontière. Il n’existe aucun protocole distinct et un gestionnaire Go ne peut pas faire la différence. Une trame qui ne contient pas de JSON valide déclenche un événement `error` et est abandonnée au lieu d’interrompre la connexion.

Utilisez directement `Stream` lorsque vous souhaitez manipuler les octets : protobuf, CBOR, formats binaires ou tout autre cas où vous préférez effectuer vous-même l’encodage.

## L’API Go

```go
// Register a handler. Runs once per connection, on its own goroutine.
func (a *App) HandleStream(name string, handler StreamHandler)

// The connection.
func (c *StreamConn) Send(data []byte) error        // blocks when the buffer is full
func (c *StreamConn) TrySend(data []byte) error     // ErrStreamFull instead of blocking
func (c *StreamConn) Receive() ([]byte, error)      // blocks until a frame or close
func (c *StreamConn) SendJSON(v any) error          // marshal and send as one frame
func (c *StreamConn) ReceiveJSON(v any) error       // receive one frame and unmarshal
func (c *StreamConn) Context() context.Context      // cancelled on disconnect
func (c *StreamConn) Window() Window                // nil in server mode
func (c *StreamConn) Name() string
func (c *StreamConn) Close() error
```

**La goroutine du gestionnaire existe aussi longtemps que la connexion.** Le retour du gestionnaire ferme la connexion. Bloquez donc sur `Receive` (ou sur `c.Context()`) aussi longtemps que vous souhaitez la maintenir ouverte. Cette structure est identique à celle d’un gestionnaire WebSocket `gorilla`/`coder`.

Les erreurs sont `ErrStreamClosed` (le pair n’est plus disponible) et `ErrStreamFull` (uniquement depuis `TrySend`).

## L’API JavaScript

`Stream(name)` renvoie un objet qui implémente le sous-ensemble utile de `WebSocket` :

| pris en charge | remarques |
| --- | --- |
| `readyState` + `CONNECTING`/`OPEN`/`CLOSING`/`CLOSED` |  |
| `onopen`, `onmessage`, `onclose`, `onerror` | plus `addEventListener` |
| `send(data)` | chaîne, `ArrayBuffer`, tableau typé ou `Blob` — consultez les règles de propriété ci-dessous |
| `JSONStream(name)` | même objet, avec des objets en entrée et en sortie |
| `close(code, reason)` |  |
| `binaryType` | **vaut `"arraybuffer"`** par défaut, et non `"blob"` |
| `bufferedAmount` | octets mis en file d’attente par `send` qui ne sont pas encore parvenus à Go |
| `protocol`, `extensions` | toujours `""` — aucune négociation |

La valeur par défaut de `binaryType` constitue le seul écart volontaire par rapport à la norme : les trames sont toujours binaires et un `Blob` imposerait une étape asynchrone supplémentaire pour lire chaque message. Définissez-la sur `"blob"` si vous souhaitez obtenir le comportement standard.

Plusieurs connexions portant le même nom de flux sont autorisées, depuis une ou plusieurs fenêtres. Chacune possède son propre `StreamConn` et sa propre goroutine de gestionnaire.

**La propriété des tampons diffère selon le sens du transfert.** La fonction JavaScript `send()` crée de manière synchrone un instantané des entrées binaires modifiables, conformément au comportement du WebSocket natif. Le code appelant peut donc les réutiliser dès que `send()` rend la main. La fonction Go `Send` transfère au transport la propriété de sa tranche sans la copier ; ne modifiez ni ne réutilisez cette tranche après la réussite de l’appel. Fournissez une nouvelle tranche lorsque le producteur doit réutiliser son espace de stockage.

## Cycle de vie

Un flux se comporte comme un socket : les événements qui ferment un socket ferment également un flux :

| événement | conséquence |
| --- | --- |
| Rechargement de la page ou navigation | la connexion se ferme, l’appel `Receive` du gestionnaire renvoie une erreur et la nouvelle page établit une nouvelle connexion |
| `window.close()` / destruction de la fenêtre | toutes les connexions de cette fenêtre se ferment |
| `s.close()` en JS | l’appel `Receive` du gestionnaire renvoie `ErrStreamClosed` |
| Retour du gestionnaire | l’interface reçoit `onclose` |
| Arrêt de l’application | le contexte de chaque connexion est annulé |

Il n’existe **aucune reconnexion automatique**, conformément à `WebSocket`. Si votre application en a besoin, la logique de reconnexion que vous utilisez déjà pour un WebSocket fonctionnera sans modification : recréez le flux dans `onclose`.

## Contre-pression

`Send` se bloque lorsque l’interface ne suit pas le rythme, comme l’écriture sur un socket se bloque lorsque le tampon d’envoi est plein. Utilisez `TrySend` si vous préférez abandonner les données plutôt qu’attendre :

```go
if err := c.TrySend(sample); errors.Is(err, application.ErrStreamFull) {
    // frontend is behind — skip this sample rather than stalling the producer
}
```

Une interface en pause — en raison d’un point d’arrêt dans les outils de développement, d’une fenêtre masquée ou d’App Nap — cesse de récupérer les données, et la limite du tampon bloque alors le producteur. Ce comportement est intentionnel : il borne l’utilisation de la mémoire au lieu de laisser croître sans limite un flux qui n’est pas lu.

## Mode serveur

Compiler avec `-tags server` remplace le transport par un **véritable WebSocket** à l’adresse `/wails/stream/ws`, puisque le mode serveur dispose déjà d’un écouteur sur lequel effectuer la mise à niveau. Le gestionnaire Go et le code frontend sont identiques : rien ne change dans votre application. Le runtime choisit le transport pour vous avant l’exécution du moindre code de module. Par défaut, les connexions WebSocket sont de même origine. Un serveur qui héberge délibérément son frontend sur une autre origine de confiance peut ajouter cet hôte avec `ServerOptions.WebSocketOriginPatterns`.

## Performances

Mesures effectuées avec `v3/tests/stream-performance`, sans limitation, avec 0 pertes et 0 réordonnancements sur environ 41 millions de trames :

|  | Pic Go→JS | Pic JS→Go |
| --- | ---: | ---: |
| macOS / WebKit-Cocoa | **3117 Mo/s** | 2793 Mo/s |
| Linux / WebKitGTK | 226 Mo/s | 727 Mo/s |
| Windows / WebView2 | 100 Mo/s | 99 Mo/s |

Le profil des performances compte davantage que les pics :

- **Go→JS est beaucoup plus rapide pour les petites trames** : 634000 trames/s sous macOS, contre environ 6200/s dans l’autre sens. Une seule réponse regroupe jusqu’à 256 trames ; JS→Go regroupe également les trames qui s’accumulent derrière une requête en cours, mais chaque connexion sérialise toujours sa propre chaîne de requêtes POST. Si vous envoyez beaucoup de petits messages, privilégiez Go→JS ou regroupez-les au niveau de l’application avant de les envoyer vers Go.
- **Sous Windows, 512 Ko est la taille optimale pour les téléversements.** Les trames plus grandes sont réparties entre plusieurs requêtes, et les mesures montrent qu’une trame de 4 Mo est *plus lente* qu’une trame de 512 Ko.
- **La latence est faible et le reste** : environ 1 à 2 ms au 99e centile sous macOS, sans dégradation sous la charge ; à 20000 trames/s, le 99e centile mesuré était *inférieur* à celui mesuré à 100 trames/s.

Les tableaux complets par plateforme et la méthode figurent dans le relevé des mesures qui accompagne cette fonctionnalité.

## Limites

|  | Limite | Comportement lorsque vous l’atteignez |
| --- | --- | --- |
| Données mises en mémoire tampon par fenêtre, en attente de collecte | 8 Mo ou 256 trames, selon la première limite atteinte | `Send` bloque ; `TrySend` renvoie `ErrStreamFull` |
| Données mises en mémoire tampon dans toute l’application, en attente de collecte ou d’écriture | 256 Mo ou 8192 trames de données | Même comportement |
| Données reçues par connexion, en attente de `Receive` | 8 Mo ou 256 trames | Le `send()` du frontend est automatiquement retenté jusqu’à ce que le gestionnaire rattrape son retard |
| Données reçues dans toute l’application, en attente de `Receive` | 256 Mo ou 8192 trames | Même comportement |
| Connexions par fenêtre | 256 | L’ouverture est automatiquement retentée jusqu’à ce qu’une place se libère |
| Connexions actives dans toute l’application | 4096 | Même comportement |
| Sessions par fenêtre | 16 | Un rechargement remplace sa propre session antérieure ; sinon, l’ouverture est retentée |
| Une seule trame, dans les deux sens | 64 Mo | Go renvoie `ErrStreamTooLarge` ; depuis JS, le flux déclenche `error` et se ferme |
| Nom d’un flux | 256 octets UTF-8 | L’ouverture est rejetée et le flux déclenche `error` |
| Durée de maintien d’une interrogation inactive | 20 s | L’interrogation renvoie un résultat vide et le runtime la relance immédiatement |

Aucune des limites ci-dessus n’entraîne de perte silencieuse de données. Les deux lignes qui indiquent *automatiquement retenté* décrivent une contre-pression ordinaire : le runtime conserve la trame et réessaie après un bref délai progressif. Votre code voit donc un flux plus lent plutôt qu’une erreur. Les lignes qui déclenchent `error` correspondent à des erreurs de programmation, et non à la charge ; elles sont donc signalées au lieu d’être masquées.

Il s’agit actuellement de constantes définies à la compilation, et non d’options. Consultez le guide sur le fonctionnement interne si vous devez les modifier.

## Quand ne pas utiliser de flux

- **Pour les échanges requête-réponse, utilisez les bindings.** Les flux servent aux données continues ou non sollicitées ; pour un appel qui renvoie une valeur, une méthode liée est plus simple.
- **Pour les événements d’application, utilisez `Emit`/`On`.** Les événements sont diffusés à tous les écouteurs et reposent sur un système distinct et éprouvé. Les flux établissent une communication point à point.
- **Non disponible pour les fenêtres `InitialHTML`.** Celles-ci sont chargées avec `origin === "null"` et ne peuvent donc pas du tout accéder au serveur de ressources.
