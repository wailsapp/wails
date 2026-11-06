---
title: "Streams"
description: "Streams bidirecionais de bytes entre Go e JavaScript, com o modelo de programação WebSocket e sem socket de escuta"
slug: "guides/streams"
sourcePath: "guides/streams.md"
---

Os streams fornecem um canal de bytes nomeado, ordenado e bidirecional entre o Go e o frontend, com o mesmo modelo de programação de um WebSocket — **sem vincular uma porta TCP**.

Não é possível usar o protocolo WebSocket por meio de um esquema de URL personalizado; portanto, a única maneira de disponibilizá-lo dentro de uma webview é executar um servidor HTTP real e escutar em uma porta. Em um aplicativo para desktop, isso significa manter uma porta local aberta, acessível por qualquer outro processo da máquina, o que exige uma verificação de origem e um token para garantir a segurança, além de ficar visível para todos os firewalls e produtos de segurança de endpoint usados pelos seus usuários. Os streams evitam tudo isso: eles trafegam pelo servidor de ativos que seu aplicativo já disponibiliza e que já está vinculado à origem.

Está migrando uma implementação WebSocket existente? Siga [Como migrar um WebSocket para Streams](/guides/streams-from-websockets/) — o guia foi escrito para ser aplicado de forma mecânica e começa pelas três diferenças que causam falhas silenciosas.

## Início rápido

Declare um stream em Go. O manipulador é executado uma vez por conexão, em sua própria goroutine:

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

Conecte-se pelo frontend usando o nome. O objeto implementa a interface `WebSocket`:

```js
import { Stream } from "@wailsio/runtime";

const s = Stream("telemetry");
s.onopen    = () => s.send(new TextEncoder().encode("hello"));
s.onmessage = (ev) => console.log(new Uint8Array(ev.data));
s.onclose   = (ev) => console.log("closed", ev.code);
```

`Stream(name)` retorna **de forma síncrona** com `readyState === CONNECTING`, exatamente como `new WebSocket(url)`, portanto você pode criar um no escopo do módulo:

```js
export const Telemetry = Stream("telemetry");
```

## Os frames são bytes

Cada frame é um `[]byte` em Go e um `ArrayBuffer` em JavaScript. Nenhum esquema ou codificação é imposto — use JSON, protobuf, CBOR ou bytes brutos, como preferir.

Um frame é uma **mensagem, não um stream de bytes**: ele chega inteiro ou não chega, e seu tamanho é transmitido junto com ele. Nenhum dos lados precisa saber o tamanho com antecedência; portanto, uma struct com um campo `[]byte` é serializada no tamanho resultante e enviada como um único frame.

## Envio de objetos

Os frames são bytes, mas raramente é desejável pensar em bytes. Os dois lados oferecem recursos auxiliares de JSON que funcionam em conjunto:

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

`JSONStream` é o mesmo objeto que `Stream`, com a codificação feita no limite — não há um protocolo separado, e um manipulador Go não consegue perceber a diferença. Um frame que não contém JSON válido gera um evento `error` e é descartado, em vez de encerrar a conexão.

Use `Stream` diretamente quando quiser os bytes: protobuf, CBOR, formatos binários ou qualquer caso em que prefira fazer a codificação por conta própria.

## A API Go

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

**A vida útil da goroutine do manipulador é a mesma da conexão.** Retornar do manipulador fecha a conexão; portanto, bloqueie em `Receive` (ou em `c.Context()`) enquanto quiser mantê-la aberta. Essa estrutura é igual à de um manipulador WebSocket `gorilla`/`coder`.

Os erros são `ErrStreamClosed` (o par não está mais disponível) e `ErrStreamFull` (somente de `TrySend`).

## A API JavaScript

`Stream(name)` retorna um objeto que implementa o subconjunto útil de `WebSocket`:

| compatível | observações |
| --- | --- |
| `readyState` + `CONNECTING`/`OPEN`/`CLOSING`/`CLOSED` |  |
| `onopen`, `onmessage`, `onclose`, `onerror` | mais `addEventListener` |
| `send(data)` | string, `ArrayBuffer`, array tipado ou `Blob` — consulte a seção sobre propriedade abaixo |
| `JSONStream(name)` | o mesmo objeto, com objetos na entrada e na saída |
| `close(code, reason)` |  |
| `binaryType` | **usa `"arraybuffer"`** por padrão, não `"blob"` |
| `bufferedAmount` | bytes enfileirados por `send` que ainda não chegaram ao Go |
| `protocol`, `extensions` | sempre `""` — não é negociado |

O valor padrão de `binaryType` é a única divergência intencional em relação ao padrão: os frames são sempre binários, e um `Blob` imporia uma etapa assíncrona adicional para ler cada mensagem. Defina-o como `"blob"` se quiser o comportamento padrão.

São permitidas várias conexões com o mesmo nome de stream, originadas de uma ou várias janelas. Cada conexão recebe seu próprio `StreamConn` e sua própria goroutine de manipulador.

**A propriedade do buffer varia conforme a direção.** O `send()` do JavaScript cria de forma síncrona uma cópia instantânea das entradas binárias mutáveis, de acordo com o comportamento nativo de WebSocket; portanto, os chamadores podem reutilizá-las assim que `send()` retornar. O `Send` do Go transfere a propriedade de seu slice para o transporte e não o copia; não modifique nem reutilize esse slice após uma chamada bem-sucedida. Forneça um novo slice quando o produtor precisar reutilizar o armazenamento.

## Ciclo de vida

Um stream se comporta como um socket, e os eventos que fecham um socket também fecham um stream:

| evento | o que acontece |
| --- | --- |
| Recarregamento ou navegação da página | a conexão é fechada, o `Receive` do manipulador retorna um erro e a nova página estabelece uma nova conexão |
| `window.close()` / janela destruída | todas as conexões dessa janela são fechadas |
| `s.close()` no JS | o `Receive` do manipulador retorna `ErrStreamClosed` |
| O manipulador retorna | o frontend recebe `onclose` |
| Encerramento do aplicativo | o contexto de cada conexão é cancelado |

Não há **reconexão automática**, o que corresponde ao comportamento de `WebSocket`. Se seu aplicativo precisar desse recurso, a lógica de reconexão que você já usa para um WebSocket funcionará sem alterações — recrie o stream em `onclose`.

## Contrapressão

`Send` bloqueia quando o frontend não consegue acompanhar, assim como uma gravação em socket bloqueia quando o buffer de envio está cheio. Use `TrySend` se preferir descartar em vez de esperar:

```go
if err := c.TrySend(sample); errors.Is(err, application.ErrStreamFull) {
    // frontend is behind — skip this sample rather than stalling the producer
}
```

Um frontend pausado — por um ponto de interrupção nas ferramentas de desenvolvimento, uma janela oculta ou o App Nap — deixa de coletar os dados, e o limite do buffer passa a bloquear o produtor. Isso é intencional: limita o uso de memória, em vez de permitir que um stream não lido cresça sem limite.

## Modo servidor

Compilar com `-tags server` substitui o transporte por um **WebSocket real** em `/wails/stream/ws`, pois o modo servidor já tem um listener no qual fazer o upgrade. O handler Go e o código do frontend são idênticos — nada muda no seu aplicativo. O runtime escolhe o transporte para você antes da execução de qualquer código de módulo. Por padrão, as conexões WebSocket têm a mesma origem. Um servidor que hospede deliberadamente seu frontend em outra origem confiável pode adicionar esse host com `ServerOptions.WebSocketOriginPatterns`.

## Desempenho

Medições feitas com `v3/tests/stream-performance`, sem limitação, com 0 perdas e 0 reordenações ao longo de ~41 milhões de frames:

|  | Pico de Go→JS | Pico de JS→Go |
| --- | ---: | ---: |
| macOS / WebKit-Cocoa | **3117 MB/s** | 2793 MB/s |
| Linux / WebKitGTK | 226 MB/s | 727 MB/s |
| Windows / WebView2 | 100 MB/s | 99 MB/s |

O padrão importa mais do que os picos:

- **Go→JS é muito mais rápido para frames pequenos** — 634000 frames/s no macOS, contra ~6200/s no sentido inverso. Uma única resposta agrega até 256 frames; JS→Go também agrupa os frames acumulados enquanto há uma solicitação em andamento, mas cada conexão ainda serializa sua própria cadeia de POSTs. Se você envia muitas mensagens pequenas, prefira Go→JS ou agrupe-as no nível do aplicativo antes de enviá-las.
- **No Windows, 512 KB é o tamanho ideal para uploads.** Frames maiores são divididos em várias solicitações, e um frame de 4 MB apresenta desempenho *inferior* ao de um frame de 512 KB.
- **A latência é baixa e permanece baixa**: p99 de ~1–2 ms no macOS, sem piora sob carga — com 20000 frames/s, mediu-se um p99 *menor* do que com 100 frames/s.

As tabelas completas por plataforma e o método estão no registro de medição que acompanha este recurso.

## Limites

|  | limite | o que acontece ao atingi-lo |
| --- | --- | --- |
| Armazenado em buffer por janela, aguardando coleta | 8 MB ou 256 frames, o que ocorrer primeiro | `Send` bloqueia; `TrySend` retorna `ErrStreamFull` |
| Armazenado em buffer em todo o aplicativo, aguardando coleta/gravação | 256 MB ou 8192 frames de dados | o mesmo |
| Recebido por conexão, aguardando `Receive` | 8 MB ou 256 frames | o `send()` do frontend é repetido automaticamente até que o handler consiga acompanhar |
| Recebido em todo o aplicativo, aguardando `Receive` | 256 MB ou 8192 frames | o mesmo |
| Conexões por janela | 256 | a abertura é repetida automaticamente até que uma vaga seja liberada |
| Conexões ativas em todo o aplicativo | 4096 | o mesmo |
| Sessões por janela | 16 | uma recarga substitui sua própria sessão anterior; caso contrário, a abertura é repetida |
| Um único frame, em qualquer sentido | 64 MB | Go retorna `ErrStreamTooLarge`; no JS, o stream gera `error` e é fechado |
| Um nome de stream | 256 bytes UTF-8 | a abertura é rejeitada e o stream gera `error` |
| Retenção de polling ocioso | 20 s | o polling retorna vazio e o runtime o emite novamente de imediato |

Nenhuma das situações acima descarta dados silenciosamente. As duas linhas que dizem *repetido automaticamente* representam contrapressão normal — o runtime retém o frame e tenta novamente após um breve intervalo de espera, de modo que seu código percebe um stream mais lento, não um erro. As linhas que geram `error` representam erros de programação, não carga, e são expostas em vez de ocultadas.

Atualmente, esses valores são constantes definidas em tempo de compilação, não opções. Consulte o guia de funcionamento interno se precisar alterá-los.

## Quando não usar um stream

- **Para solicitação/resposta, use bindings.** Streams destinam-se a dados contínuos ou não solicitados; uma chamada que retorna um valor é mais simples como método vinculado.
- **Para eventos do aplicativo, use `Emit`/`On`.** Os eventos são distribuídos a todos os listeners e constituem um sistema separado e consolidado. Streams são ponto a ponto.
- **Não disponível para janelas `InitialHTML`.** Elas são carregadas com `origin === "null"`, portanto não conseguem acessar o servidor de assets.
