---
title: "Migrando um WebSocket para Streams"
description: "Conversão passo a passo de uma implementação WebSocket existente para streams do Wails, incluindo as diferenças que causam falhas silenciosas"
slug: "guides/streams-from-websockets"
sourcePath: "guides/streams-from-websockets.md"
---

Um guia prático para converter um aplicativo que atualmente executa um servidor WebSocket para [Streams](/guides/streams/). Escrito para ser seguido literalmente, inclusive por um agente.

A vantagem é eliminar o listener: nenhuma porta TCP vinculada, nenhuma verificação de origem, nenhum token, nada a que um firewall ou produto de segurança de endpoints possa se opor. A API é semelhante o suficiente para que a maior parte do código do frontend permaneça intacta — mas **três diferenças causam falhas silenciosas**. Elas são apresentadas primeiro porque podem fazer você perder uma tarde inteira.

## O caso comum: um servidor HTTP local como solução alternativa

O motivo mais comum para um aplicativo Wails ter um WebSocket é não haver outra maneira de enviar um fluxo contínuo ao frontend. Por isso, o aplicativo inicia seu próprio `http.Server` em uma porta local, e o frontend se conecta a ele. Se essa for a estrutura do seu aplicativo, esta migração elimina completamente o servidor — e, com ele, vários recursos que você criou *ao redor* dele.

**O mecanismo de descoberta de porta é eliminado.** Algo precisa informar ao frontend a qual porta se conectar: um `GetServerPort()` vinculado, uma porta fixa com uma alternativa caso esteja ocupada, uma variável global injetada ou um valor armazenado em `localStorage`. Tudo isso desaparece — um stream é identificado por nome, e esse nome é uma constante de tempo de compilação em ambos os lados.

```go
// BEFORE
ln, _ := net.Listen("tcp", "127.0.0.1:0")
go http.Serve(ln, mux)
port := ln.Addr().(*net.TCPAddr).Port     // ...and a binding to hand `port` to the frontend

// AFTER
app.HandleStream("feed", handler)         // that is the entire replacement
```

**A configuração de CORS é eliminada.** A origem da webview é `wails://` ou `http://wails.localhost`, dependendo da plataforma. Portanto, um servidor local precisa de `CheckOrigin`, de um cabeçalho `Access-Control-Allow-Origin` ou de ambos. Os streams usam o mesmo servidor de assets do qual a página foi carregada, portanto não há nenhuma solicitação entre origens a ser permitida.

**Qualquer token de autenticação que você tenha criado é eliminado.** Uma porta vinculada ao localhost pode ser acessada por qualquer processo da máquina; por isso, uma implementação cuidadosa adiciona um token ou nonce para impedir a conexão de outro software. Agora não há mais nenhuma porta para acessar.

**Os endpoints que não são WebSocket passam para o middleware do servidor de assets.** Esses servidores raramente permanecem exclusivos para WebSocket — downloads de arquivos, endpoints de imagens e verificações de integridade tendem a se acumular junto ao socket. Os streams não substituem esses recursos, mas você também não precisa de um segundo servidor para eles. Monte os mesmos handlers no servidor de assets:

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

Em seguida, o frontend chama `/api/...` como uma URL relativa de mesma origem — sem host, sem porta e sem CORS. Com os streams e essa mudança, o servidor local não tem mais nenhuma função.

## Leia isto antes de começar

### 1. `ev.data` é um `ArrayBuffer`, nunca uma string

Esta é a principal diferença. Um WebSocket entrega mensagens de texto como strings; um stream entrega todas as mensagens como bytes. Um código como este **compila, é executado e se comporta incorretamente**:

```js
// BEFORE — works with WebSocket
ws.onmessage = (ev) => { const msg = JSON.parse(ev.data); ... };

// AFTER — ev.data is an ArrayBuffer, JSON.parse gets "[object ArrayBuffer]"
```

Corrija isso no ponto de entrada, em vez de corrigir cada local de chamada:

```js
const dec = new TextDecoder();
s.onmessage = (ev) => { const msg = JSON.parse(dec.decode(ev.data)); ... };
```

Ou, se o tráfego for JSON — como geralmente é — use `JSONStream` em vez de `Stream` e evite completamente o problema:

```js
import { JSONStream } from "@wailsio/runtime";

const s = JSONStream("feed");
s.onmessage = (ev) => dispatch(ev.data);   // already an object
s.send({ hello: true });                   // stringified for you
```

Essa é a migração mais curta para um WebSocket JSON: troque o construtor, remova as chamadas `JSON.parse` e `JSON.stringify`, e o restante do handler permanece inalterado. Para tráfego que não seja JSON, encapsule uma vez e mantenha todos os handlers intactos — consulte [a camada de compatibilidade](#camada-de-compatibilidade).

### 2. O envio é compatível; o recebimento não

`send()` aceita uma string e a codifica como UTF-8, portanto `s.send(JSON.stringify(x))` funciona sem alterações. Apenas o caminho de recebimento precisa ser editado. É fácil não perceber essa assimetria, pois metade do código continua funcionando.

### 3. Não há URL

Um WebSocket transporta parâmetros de conexão em sua URL — caminho, string de consulta, subprotocolo e token de autenticação. Um stream tem apenas um nome. Tudo o que era passado na URL precisa ir para o primeiro frame ou para um método vinculado chamado antes da conexão.

```js
// BEFORE
const ws = new WebSocket(`wss://host/feed?topic=${topic}&token=${token}`);

// AFTER — no token needed at all; the app is the only possible caller
const s = Stream("feed");
s.onopen = () => s.send(JSON.stringify({ subscribe: topic }));
```

## Lado do Go

Exclua o servidor HTTP, o upgrader e o registro de conexões. Cada um deles é substituído por um handler.

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
| `upgrader.Upgrade` / `websocket.Accept` | *(nada — `HandleStream` é todo o registro)* |
| `conn.ReadMessage()` / `conn.Read(ctx)` | `c.Receive()` |
| `conn.WriteMessage(TextMessage, b)` | `c.Send(b)` |
| `conn.Close()` | `c.Close()` ou simplesmente retorne do handler |
| registro de conexões para broadcast | mantenha o seu próprio — consulte [Broadcast](#difuso) |
| `http.ListenAndServe`, mux, verificação de origem, token | **excluir** |
| keepalive com ping/pong | **excluir** — não há socket ocioso cuja conexão precise ser mantida ativa |
| `r.Context()` | `c.Context()` |

O tempo de vida da goroutine do handler é o tempo de vida da conexão, exatamente como em um handler do gorilla; portanto, a estrutura do loop existente permanece inalterada.

## Lado do frontend

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

`readyState`, as quatro constantes de estado, `addEventListener`, `close(code, reason)` e `bufferedAmount` se comportam da mesma forma que em um `WebSocket`.

### Camada de compatibilidade

Se preferir não alterar nenhum handler, encapsule o construtor uma única vez. Assim, o código existente que espera mensagens de texto funciona sem nenhuma alteração:

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

Assim, `const ws = TextStream("feed")` pode substituir diretamente `new WebSocket(url)`.

## Difusão

Um servidor WebSocket geralmente mantém um registro para poder distribuir mensagens a vários destinatários. Os streams não têm difusão integrada — mantenha o registro, mas armazene `*StreamConn` em vez de `*websocket.Conn`:

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

Vale a pena manter duas regras: libere o bloqueio antes de enviar e, em uma difusão, prefira `TrySend` para que um único consumidor lento não bloqueie todos os outros clientes.

## O outro caso: o frontend se comunica diretamente com um broker

Isso é menos comum em um aplicativo Wails, mas é importante saber. Se o frontend abrir um WebSocket **com um broker, e não com seu aplicativo** — `nats.ws` com um servidor NATS ou MQTT sobre WebSocket —, um stream não será um substituto direto, pois ele conecta o frontend *ao seu código Go*, e não a terceiros.

A migração é uma mudança de arquitetura e geralmente é uma boa mudança:

```
BEFORE   frontend ──ws──► NATS server            (credentials in the frontend)
AFTER    frontend ──stream──► Go ──►  NATS       (credentials stay in Go)
```

Mova o cliente do broker para o Go, onde a biblioteca nativa é melhor do que a biblioteca para navegador, e exponha por meio de um stream as partes necessárias ao frontend:

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

Observe `TrySend` dentro do callback da assinatura: esse callback é executado na goroutine do cliente do broker, e bloqueá-lo interromperia a entrega de todas as assinaturas da conexão.

As vantagens são: as credenciais do broker nunca chegam ao frontend, nenhuma porta WebSocket fica exposta na máquina e a reconexão e o recuo ficam a cargo do cliente Go consolidado, em vez do cliente para navegador.

## Lista de verificação da migração

- [ ] `HandleStream` registrado para cada endpoint WebSocket existente
- [ ] Loop de leitura convertido: `ReadMessage`/`Read` → `c.Receive()`
- [ ] Gravações convertidas: `WriteMessage` → `c.Send()`, ou `TrySend` em qualquer difusão ou callback do broker
- [ ] Servidor HTTP, entrada do mux, upgrader, verificação de origem e token de autenticação **excluídos**
- [ ] Keepalive com ping/pong **excluído**
- [ ] Parâmetros da URL movidos para o primeiro frame ou para um método vinculado
- [ ] **Todas as leituras de `ev.data` decodificadas** — `new TextDecoder().decode(ev.data)` — ou shim adotado
- [ ] Lógica de reconexão mantida como está (por design, não há reconexão integrada)
- [ ] Clientes de broker movidos para o Go se o frontend se comunicava diretamente com um deles
- [ ] Verificação da existência de janelas `InitialHTML` — elas não podem usar streams de forma alguma

Itens que vale a pena procurar com grep durante a conversão:

```
new WebSocket(     ev.data            .onmessage
websocket.Accept   upgrader.Upgrade   ReadMessage
WriteMessage       ListenAndServe     CheckOrigin
```

## Diferenças de comportamento esperadas

|  | WebSocket | Stream |
| --- | --- | --- |
| Tipo de mensagem | texto ou binário | somente bytes |
| `ev.data` | string ou `Blob`/`ArrayBuffer` | sempre `ArrayBuffer` (a menos que `binaryType = "blob"`) |
| padrão de `binaryType` | `"blob"` | `"arraybuffer"` |
| Subprotocolos, `extensions` | negociados | não compatíveis; sempre `""` |
| Parâmetros de conexão | URL e consulta | primeiro frame ou uma chamada vinculada |
| Autenticação | token ou cookie | desnecessária — o aplicativo é o único chamador |
| Keepalive | ping/pong | não é necessário |
| Reconexão automática | nenhuma | nenhuma (igual) |
| Códigos de fechamento | faixa completa | `1000` normal, `1001` sessão encerrada, `1002` incompatibilidade de enquadramento, `1006` erro |
| Contrapressão | buffer de socket do kernel | 8 MB/256 quadros por janela; depois, `Send` é bloqueado |
| Várias conexões, um único endpoint | sim | sim |

## Após a migração

Verificações básicas que detectam os erros comuns:

1. Recarregue a página várias vezes — o manipulador deve encerrar e um novo deve iniciar a cada vez, sem nunca se acumularem.
2. Envie uma mensagem maior que 512 KB em cada direção.
3. Deixe-o ocioso por mais de um minuto; o tráfego deve ser retomado sem uma nova conexão.
4. Abra as ferramentas de desenvolvimento e pause em um ponto de interrupção sob carga; depois, retome a execução — o produtor deve bloquear e se recuperar, em vez de perder dados ou crescer sem limites.
