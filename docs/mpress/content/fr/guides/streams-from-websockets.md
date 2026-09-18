---
title: "Migration d’un WebSocket vers les Streams"
description: "Conversion pas à pas d’une implémentation WebSocket existante vers les streams Wails, y compris les différences qui provoquent des dysfonctionnements silencieux"
slug: "guides/streams-from-websockets"
sourcePath: "guides/streams-from-websockets.md"
---

Guide méthodique pour convertir une application qui exécute actuellement un serveur WebSocket afin qu’elle utilise les [Streams](/guides/streams/). Il est conçu pour être suivi à la lettre, y compris par un agent.

Le principal avantage est la suppression de l’écouteur : aucun port TCP lié, aucun contrôle d’origine, aucun jeton, rien qu’un pare-feu ou un produit de sécurité des terminaux puisse bloquer. L’API est suffisamment proche pour que la majeure partie du code frontend reste intacte — mais **trois différences provoquent des dysfonctionnements silencieux**. Elles sont présentées en premier, car ce sont celles qui risquent de vous faire perdre un après-midi.

## Cas courant : un serveur HTTP local comme solution de contournement

Une application Wails utilise généralement un WebSocket parce qu’il n’existait aucun autre moyen d’envoyer un flux continu au frontend. L’application lance donc son propre `http.Server` sur un port local, auquel le frontend se connecte. Si votre architecture correspond à ce cas, cette migration supprime purement et simplement le serveur, ainsi que plusieurs éléments que vous avez construits *autour de* celui-ci.

**Le mécanisme de découverte du port disparaît.** Un mécanisme doit indiquer au frontend le port auquel se connecter : un `GetServerPort()` lié, un port fixe assorti d’un port de secours s’il est occupé, une variable globale injectée ou une valeur stockée dans `localStorage`. Tout cela disparaît : un stream est désigné par un nom, et ce nom est une constante définie à la compilation des deux côtés.

```go
// BEFORE
ln, _ := net.Listen("tcp", "127.0.0.1:0")
go http.Serve(ln, mux)
port := ln.Addr().(*net.TCPAddr).Port     // ...and a binding to hand `port` to the frontend

// AFTER
app.HandleStream("feed", handler)         // that is the entire replacement
```

**La configuration CORS disparaît.** Selon la plateforme, l’origine de la webview est `wails://` ou `http://wails.localhost`. Un serveur local nécessite donc `CheckOrigin`, un en-tête `Access-Control-Allow-Origin`, ou les deux. Les streams utilisent le serveur de ressources depuis lequel la page a été chargée ; aucune requête interorigine n’est donc à autoriser.

**Tout jeton d’authentification que vous avez créé disparaît.** Tous les processus de la machine peuvent accéder à un port lié sur localhost. Une implémentation rigoureuse ajoute donc un jeton ou un nonce pour empêcher d’autres logiciels de s’y connecter. Désormais, il n’existe plus de port auquel accéder.

**Les points de terminaison autres que WebSocket sont déplacés vers le middleware du serveur de ressources.** Ces serveurs restent rarement dédiés au WebSocket : un téléchargement de fichier, un point de terminaison d’image ou une vérification de l’état de fonctionnement du service finissent souvent par s’ajouter au socket. Les streams ne les remplacent pas, mais ils ne nécessitent pas non plus un second serveur. Montez les mêmes gestionnaires sur le serveur de ressources :

```go
app := application.New(application.Options{
    Assets: application.AssetOptions{
        Handler: yourFrontendAssets,
        Middleware: func(next http.Handler) http.Handler {
            return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                if strings.HasPrefix(r.URL.Path, "/api/") {
                    yourExistingMux.ServeHTTP(w, r)   // the handlers you already wrote
                    return
                }
                next.ServeHTTP(w, r)
            })
        },
    },
})
```

Le frontend appelle ensuite `/api/...` au moyen d’une URL relative de même origine : aucun hôte, aucun port, aucune configuration CORS. Avec les streams et ce changement, le serveur local n’a plus rien à faire.

## À lire avant de commencer

### 1. `ev.data` est un `ArrayBuffer`, jamais une chaîne

C’est la différence majeure. Un WebSocket transmet les messages texte sous forme de chaînes ; un stream transmet chaque message sous forme d’octets. Un code comme celui-ci **se compile, s’exécute et se comporte de façon incorrecte** :

```js
// BEFORE — works with WebSocket
ws.onmessage = (ev) => { const msg = JSON.parse(ev.data); ... };

// AFTER — ev.data is an ArrayBuffer, JSON.parse gets "[object ArrayBuffer]"
```

Corrigez le problème à la frontière plutôt qu’à chaque site d’appel :

```js
const dec = new TextDecoder();
s.onmessage = (ev) => { const msg = JSON.parse(dec.decode(ev.data)); ... };
```

Ou, si vos échanges utilisent JSON — ce qui est généralement le cas — utilisez `JSONStream` au lieu de `Stream` et évitez entièrement le problème :

```js
import { JSONStream } from "@wailsio/runtime";

const s = JSONStream("feed");
s.onmessage = (ev) => dispatch(ev.data);   // already an object
s.send({ hello: true });                   // stringified for you
```

C’est la migration la plus courte pour un WebSocket JSON : remplacez le constructeur, supprimez les appels à `JSON.parse` et `JSON.stringify`, et laissez le reste du gestionnaire inchangé. Pour les échanges non JSON, créez une seule enveloppe et ne modifiez aucun gestionnaire ; consultez [la couche de compatibilité](#couche-de-compatibilit).

### 2. L’envoi est compatible, contrairement à la réception

`send()` accepte une chaîne et l’encode en UTF-8 ; `s.send(JSON.stringify(x))` fonctionne donc sans modification. Seul le chemin de réception doit être modifié. Cette asymétrie passe facilement inaperçue, car la moitié de votre code continue de fonctionner.

### 3. Il n’y a pas d’URL

Un WebSocket transporte les paramètres de connexion dans son URL : chemin, chaîne de requête, sous-protocole et jeton d’authentification. Un stream ne possède qu’un nom. Tout ce que vous transmettiez dans l’URL doit être déplacé dans la première trame ou dans une méthode liée appelée avant la connexion.

```js
// BEFORE
const ws = new WebSocket(`wss://host/feed?topic=${topic}&token=${token}`);

// AFTER — no token needed at all; the app is the only possible caller
const s = Stream("feed");
s.onopen = () => s.send(JSON.stringify({ subscribe: topic }));
```

## Côté Go

Supprimez le serveur HTTP, le composant de mise à niveau et le registre des connexions. Chacun est remplacé par un gestionnaire.

```go
// BEFORE — gorilla/coder websocket
func (a *App) serveWS(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }
    defer conn.Close()

    clients.add(conn)
    defer clients.remove(conn)

    for {
        _, data, err := conn.ReadMessage()
        if err != nil {
            return
        }
        handle(data)
    }
}

// ...plus http.ListenAndServe, a mux entry, an origin checker, and a token check
```

```go
// AFTER
app.HandleStream("feed", func(c *application.StreamConn) {
    defer c.Close()

    for {
        frame, err := c.Receive()
        if err != nil {
            return                  // reload, close, or shutdown
        }
        handle(frame)
    }
})
```

| WebSocket | Stream |
| --- | --- |
| `upgrader.Upgrade` / `websocket.Accept` | *rien — `HandleStream` constitue l’intégralité de l’enregistrement* |
| `conn.ReadMessage()` / `conn.Read(ctx)` | `c.Receive()` |
| `conn.WriteMessage(TextMessage, b)` | `c.Send(b)` |
| `conn.Close()` | `c.Close()`, ou quittez simplement le gestionnaire |
| registre des connexions pour la diffusion | conservez le vôtre — voir [Diffusion](#diffusion) |
| `http.ListenAndServe`, mux, contrôle d’origine, jeton | **supprimer** |
| maintien de connexion par ping/pong | **supprimer** — il n’y a aucun socket inactif à maintenir en vie |
| `r.Context()` | `c.Context()` |

La durée de vie de la goroutine du gestionnaire correspond à celle de la connexion, exactement comme pour un gestionnaire gorilla ; vous pouvez donc conserver telle quelle la structure de votre boucle existante.

## Côté frontend

```js
// BEFORE
const ws = new WebSocket(url);
ws.onopen    = () => ws.send(JSON.stringify(hello));
ws.onmessage = (ev) => dispatch(JSON.parse(ev.data));
ws.onclose   = () => scheduleReconnect();

// AFTER
import { Stream } from "@wailsio/runtime";
const dec = new TextDecoder();

const s = Stream("feed");
s.onopen    = () => s.send(JSON.stringify(hello));   // unchanged
s.onmessage = (ev) => dispatch(JSON.parse(dec.decode(ev.data)));
s.onclose   = () => scheduleReconnect();             // unchanged
```

`readyState`, les quatre constantes d’état, `addEventListener`, `close(code, reason)` et `bufferedAmount` se comportent tous comme sur un `WebSocket`.

### Couche de compatibilité

Si vous préférez ne modifier aucun gestionnaire, encapsulez une seule fois le constructeur. Le code existant qui attend des messages sous forme de chaînes fonctionne alors tel quel :

```js
import { Stream } from "@wailsio/runtime";

/** A Stream that delivers text messages as strings, like a WebSocket. */
export function TextStream(name) {
    const s = Stream(name);
    const dec = new TextDecoder();
    s.binaryType = "arraybuffer";

    const add = s.addEventListener.bind(s);
    const remove = s.removeEventListener.bind(s);
    const wrappers = new WeakMap();
    const decoded = new WeakMap();

    const decodeEvent = (ev) => {
        if (decoded.has(ev)) return decoded.get(ev);
        const data = typeof ev.data === "string" ? ev.data : dec.decode(ev.data);
        const textEvent = new MessageEvent("message", { data });
        decoded.set(ev, textEvent);
        return textEvent;
    };

    const wrap = (listener) => {
        let wrapper = wrappers.get(listener);
        if (wrapper) return wrapper;
        wrapper = (ev) => {
            const textEvent = decodeEvent(ev);
            if (typeof listener === "function") listener.call(s, textEvent);
            else listener.handleEvent(textEvent);
        };
        wrappers.set(listener, wrapper);
        return wrapper;
    };

    s.addEventListener = (type, listener, options) =>
        add(type, type === "message" && listener ? wrap(listener) : listener, options);
    s.removeEventListener = (type, listener, options) =>
        remove(type, type === "message" && listener ? wrappers.get(listener) ?? listener : listener, options);

    // A WailsSocket implements onmessage through addEventListener, but a native
    // WebSocket uses an internal event-handler slot. Define the property on the
    // instance so both transports pass property handlers through the same
    // decoding wrapper as addEventListener listeners.
    let onmessage = null;
    Object.defineProperty(s, "onmessage", {
        get: () => onmessage,
        set(listener) {
            if (onmessage) s.removeEventListener("message", onmessage);
            onmessage = typeof listener === "function" ? listener : null;
            if (onmessage) s.addEventListener("message", onmessage);
        },
        configurable: true,
        enumerable: true,
    });
    return s;
}
```

`const ws = TextStream("feed")` remplace alors directement `new WebSocket(url)`.

## Diffusion

Un serveur WebSocket tient généralement un registre pour pouvoir diffuser les messages. Les streams n’intègrent aucun mécanisme de diffusion : conservez le registre, mais stockez `*StreamConn` au lieu de `*websocket.Conn` :

```go
type hub struct {
    mu    sync.Mutex
    conns map[*application.StreamConn]struct{}
}

func (h *hub) add(c *application.StreamConn)    { h.mu.Lock(); h.conns[c] = struct{}{}; h.mu.Unlock() }
func (h *hub) remove(c *application.StreamConn) { h.mu.Lock(); delete(h.conns, c); h.mu.Unlock() }

func (h *hub) broadcast(msg []byte) {
    h.mu.Lock()
    conns := make([]*application.StreamConn, 0, len(h.conns))
    for c := range h.conns {
        conns = append(conns, c)
    }
    h.mu.Unlock()                         // never hold the lock across Send

    for _, c := range conns {
        // TrySend, not Send: one stalled frontend must not block the fan-out.
        _ = c.TrySend(msg)
    }
}

app.HandleStream("feed", func(c *application.StreamConn) {
    h.add(c)
    defer h.remove(c)
    defer c.Close()
    <-c.Context().Done()
})
```

Deux règles sont à conserver : libérez le verrou avant l’envoi et, lors d’une diffusion, privilégiez `TrySend` afin qu’un seul consommateur lent ne puisse pas bloquer tous les autres clients.

## L’autre cas : le frontend communique directement avec un broker

Ce cas est moins courant dans une application Wails, mais il est utile de le connaître. Si le frontend ouvre un WebSocket **vers un broker plutôt que vers votre application** — `nats.ws` vers un serveur NATS, ou MQTT sur WebSocket —, un stream ne peut pas le remplacer directement, car il connecte le frontend à *votre code Go* et non à un tiers.

Cette migration constitue une modification de l’architecture, généralement bénéfique :

```
BEFORE   frontend ──ws──► NATS server            (credentials in the frontend)
AFTER    frontend ──stream──► Go ──►  NATS       (credentials stay in Go)
```

Déplacez le client du broker dans Go, où la bibliothèque native est meilleure que celle du navigateur, puis exposez par l’intermédiaire d’un stream les éléments nécessaires au frontend :

```go
nc, _ := nats.Connect(url, nats.UserCredentials(credsPath))  // creds never reach the frontend

app.HandleStream("nats", func(c *application.StreamConn) {
    defer c.Close()

    var subs []*nats.Subscription
    defer func() {
        for _, s := range subs {
            _ = s.Unsubscribe()
        }
    }()

    for {
        frame, err := c.Receive()
        if err != nil {
            return
        }

        var cmd struct {
            Op      string          `json:"op"`       // "sub" | "pub"
            Subject string          `json:"subject"`
            Data    json.RawMessage `json:"data"`
        }
        if json.Unmarshal(frame, &cmd) != nil {
            continue
        }

        switch cmd.Op {
        case "sub":
            sub, err := nc.Subscribe(cmd.Subject, func(m *nats.Msg) {
                out, _ := json.Marshal(map[string]any{"subject": m.Subject, "data": m.Data})
                // TrySend: a slow frontend must not block the NATS callback.
                _ = c.TrySend(out)
            })
            if err == nil {
                subs = append(subs, sub)
            }
        case "pub":
            _ = nc.Publish(cmd.Subject, cmd.Data)
        }
    }
})
```

Notez `TrySend` dans le callback de l’abonnement : ce callback s’exécute dans la goroutine du client du broker. Le bloquer interromprait la livraison pour tous les abonnements de la connexion.

Vous y gagnez les avantages suivants : les identifiants du broker n’atteignent jamais le frontend, aucun port WebSocket n’est exposé sur la machine et la reconnexion avec temporisation progressive est gérée par le client Go éprouvé plutôt que par celui du navigateur.

## Liste de contrôle de la migration

- [ ] `HandleStream` enregistré pour chaque endpoint WebSocket existant
- [ ] Boucle de lecture convertie : `ReadMessage`/`Read` → `c.Receive()`
- [ ] Écritures converties : `WriteMessage` → `c.Send()`, ou `TrySend` dans chaque diffusion ou callback du broker
- [ ] Serveur HTTP, entrée du mux, mécanisme de mise à niveau, contrôle de l’origine et jeton d’authentification **supprimés**
- [ ] Maintien de connexion par ping/pong **supprimé**
- [ ] Paramètres d’URL déplacés dans une première frame ou une méthode liée
- [ ] **Chaque lecture de `ev.data` décodée** — `new TextDecoder().decode(ev.data)` — ou couche de compatibilité adoptée
- [ ] Logique de reconnexion conservée telle quelle (par conception, aucune reconnexion n’est intégrée)
- [ ] Clients de broker déplacés dans Go si le frontend communiquait directement avec un broker
- [ ] Vérification de la présence de fenêtres `InitialHTML` : elles ne peuvent pas du tout utiliser les streams

Éléments qu’il est utile de rechercher avec grep pendant la conversion :

```
new WebSocket(     ev.data            .onmessage
websocket.Accept   upgrader.Upgrade   ReadMessage
WriteMessage       ListenAndServe     CheckOrigin
```

## Différences de comportement à prévoir

|  | WebSocket | Stream |
| --- | --- | --- |
| Type de message | texte ou binaire | octets uniquement |
| `ev.data` | chaîne ou `Blob`/`ArrayBuffer` | toujours `ArrayBuffer` (sauf avec `binaryType = "blob"`) |
| Valeur par défaut de `binaryType` | `"blob"` | `"arraybuffer"` |
| Sous-protocoles, `extensions` | négociés | non pris en charge ; toujours `""` |
| Paramètres de connexion | URL et chaîne de requête | première frame ou appel lié |
| Authentification | jeton ou cookie | inutile : l’application est le seul appelant |
| Maintien de connexion | ping/pong | aucun requis |
| Reconnexion automatique | aucune | aucune (identique) |
| Codes de fermeture | plage complète | `1000` fermeture normale, `1001` session fermée, `1002` incompatibilité de framing, `1006` erreur |
| Contre-pression | tampon de socket du noyau | 8 Mo / 256 trames par fenêtre, puis `Send` se bloque |
| Plusieurs connexions, un seul point de terminaison | oui | oui |

## Après la migration

Vérifications rapides permettant de détecter les erreurs courantes :

1. Rechargez la page à plusieurs reprises : le gestionnaire devrait s’arrêter et un nouveau devrait démarrer à chaque fois, sans que les gestionnaires ne s’accumulent jamais.
2. Envoyez dans chaque direction un message de plus de 512 Ko.
3. Laissez l’application inactive pendant plus d’une minute ; le trafic devrait reprendre sans nouvelle connexion.
4. Ouvrez les outils de développement et suspendez l’exécution sur un point d’arrêt lorsque le système est sous charge, puis reprenez-la : le producteur devrait se bloquer puis reprendre son fonctionnement, plutôt que de perdre des données ou de croître sans limite.
