---
title: "Streams — detalhes internos"
description: "Como funciona o transporte de streams, por que ele foi construído dessa forma, o que significam as constantes de buffer e o que ainda não foi concluído"
slug: "guides/advanced/streams-internals"
sourcePath: "guides/advanced/streams-internals.md"
---

Referência para qualquer pessoa — ou agente — que altere o transporte de streams. A API voltada ao usuário está em [Streams](/guides/streams/); esta página trata dos mecanismos subjacentes e das decisões que os fundamentam, pois várias delas parecem arbitrárias até que se saiba quais problemas evitam.

## Arquivos

| arquivo | função |
| --- | --- |
| `v3/pkg/application/stream.go` | API pública, `StreamConn`, `streamSink`, o gerenciador e seu registro |
| `v3/pkg/application/stream_session.go` | um carregamento de página em uma janela: a fila de saída, os tipos de frame e a tabela de conexões |
| `v3/pkg/application/stream_transport.go` | os dois endpoints HTTP, o enquadramento binário, a remontagem de partes e o preâmbulo do runtime |
| `v3/pkg/application/stream_server.go` | somente `-tags server`: destino WebSocket real |
| `v3/pkg/application/stream_prelude_{server,desktop}.go` | seleciona o transporte do cliente quando o pacote é servido |
| `v3/internal/runtime/desktop/@wailsio/runtime/src/stream.ts` | o cliente com o formato de `WebSocket` |
| `v3/tests/stream-performance/` | ferramenta de testes de carga (`-upload`, `-reloads`, varreduras de cenários) |

## A estrutura

Go→JS e JS→Go usam mecanismos diferentes, e essa assimetria constitui todo o projeto.

```
Go                                  webview
──                                  ───────
Send() ─► per-window queue ─────────► GET  /wails/stream/poll   (held open)
                                      └─ one held request per window,
                                         carrying frames for every connection

Receive() ◄─ per-conn inbox ◄──────── POST /wails/stream/send   (one or more frames)
```

**Go→JS é uma consulta mantida em espera.** A requisição fica suspensa até que haja algo para entregar. Deliberadamente, não há intervalo de consulta nem qualquer mecanismo adaptativo: o servidor mantém a requisição até existir um frame, portanto a latência de entrega já é de ~0, e qualquer intervalo no cliente só poderia aumentá-la. Os frames que chegam enquanto uma resposta está em trânsito se acumulam e seguem na próxima, fazendo da própria viagem de ida e volta a janela de agrupamento — ela se amplia conforme a carga aumenta, sem que nada precise medi-la. Medições: 1.0 frames por resposta a 100/s, ainda 1.0 a 5000/s, 3.4 a 20000/s, e a latência p99 *cai* à medida que a taxa aumenta.

**JS→Go usa POSTs comuns.** Os envios são serializados por conexão com uma cadeia de promises, pois chamadas `fetch` simultâneas não preservam a ordem, e o Go depende de observar os envios na ordem em que foram feitos. Os frames acumulados atrás de uma requisição em trânsito são agrupados no POST seguinte. O Go acrescenta o frame aceito ou o prefixo aceito do lote à caixa de entrada da conexão *antes* de responder, portanto o cliente não pode avançar além dos bytes que o Go já colocou na fila.

**Uma consulta em trânsito por janela, multiplexando todas as conexões.** É isso que garante a ordenação por construção — uma fila, um único consumidor e nenhum segundo caminho de entrega que possa ultrapassar o primeiro. Isso também contorna o limite de seis conexões por host do HTTP/1.1 no Windows, onde essas são requisições de rede reais do Chromium para `http://wails.localhost`.

## Por que essas decisões específicas

Cada uma delas é uma cicatriz do trabalho no transporte de eventos. Remover qualquer uma reintroduz um bug comprovado por medições.

**Nada no caminho Go→JS toca na thread principal.** `Send` acrescenta o item sob a proteção de um mutex e retorna. Antes, os eventos executavam sua avaliação inline quando eram emitidos pela thread principal enquanto uma emissão anterior de uma goroutine ainda estava na fila — 4.4% dos eventos ficavam em ordem invertida nas três plataformas. Uma única fila com um único consumidor não permite que isso aconteça.

**Nada toca em `evaluateJavaScript`, independentemente do tamanho.** Inserir a carga útil no código-fonte avaliado retém memória do host acima de um limiar específico de cada plataforma: 11.6 GB no macOS e 6.2 GB no WebKitGTK a 100 × 1 MB/s. Os streams não chegam perto desse mecanismo, razão pela qual a varredura com taxa constante de bytes permanece plana para todos os tamanhos de frame.

**Os dados de controle trafegam nos cabeçalhos, nunca no corpo nem na string de consulta.** O 6.0 do WebKitGTK pode entregar corpos de POST como parâmetros de consulta em esquemas de URI personalizados (`transport_http.go` contém uma alternativa exatamente para esse caso), e o WebView2 limita a entrega do corpo a cerca de 2 MB.

**A resposta da consulta é binária, não JSON.** Os frames são `[]byte`; usar base64 dentro de um envelope JSON acrescentaria um custo de 33% a cada frame, além de uma análise sintática na thread da interface.

```
magic "WS1\0" | flags u8 | count u32 | count × ( connID u32 | kind u8 | len u32 | payload )
```

`kind` representa data / open / close / error. Não há número de sequência nem confirmação: um WebSocket não repete dados, e uma conexão interrompida perde o que estava em trânsito. Emular esse comportamento é mais simples e mais honesto do que usar um cursor que o buffer limitado nem sempre poderia atender.

**Manter uma requisição em espera é seguro** porque cada requisição da webview já recebe sua própria goroutine. `dispatchWorkers` em `assetserver_webview.go` está fixado em 0, com um comentário que menciona exatamente esse caso; ativar esse pool exigiria primeiro um limite para a duração da requisição.

## Constantes de buffer

Todas ficam em `stream.go`. **São constantes de tempo de compilação, não opções** — não há `Options.Streams` nem configuração por stream. Para alterá-las, é preciso editar o arquivo.

| constante | valor | o que ela limita |
| --- | ---: | --- |
| `streamOutQueueBytes` | 8 MB | bytes armazenados em buffer por janela aguardando coleta |
| `streamOutQueueDepth` | 256 | frames armazenados em buffer por janela |
| `streamOutQueueBytesGlobal` / `streamOutQueueDepthGlobal` | 256 MB / 8192 | dados de saída armazenados em buffer em todo o aplicativo |
| `streamInQueueBytesGlobal` / `streamInQueueDepthGlobal` | 256 MB / 8192 | dados de entrada aguardando `Receive` em todo o aplicativo |
| `streamMaxConnections` | 256 | conexões ativas mais fechamentos na fila em uma sessão |
| `streamMaxConnectionsGlobal` | 4096 | conexões ativas em todo o aplicativo |
| `streamOutCloseDepthGlobal` | 4096 | notificações de fechamento não entregues em todo o aplicativo |
| `streamMaxSessionsPerWindow` | 16 | sessões que uma janela pode manter antes que uma geração mais recente precise substituir uma anterior |
| `streamMaxSessions` | 1024 | sessões em todo o aplicativo |
| `streamOutControlDepth` / `streamOutControlDepthGlobal` | 256 / 4096 | quadros de controle que não são de fechamento na fila, por sessão e em toda a aplicação |
| `streamMaxChunkSets` / `streamMaxChunkTotal` | 256 / 4096 | uploads incompletos por sessão e partes em um upload |
| `streamMaxChunkBytesGlobal` / `streamMaxChunkPartsGlobal` | 128 MB / 4096 | carga útil de chunks e metadados de partes em toda a aplicação |
| `streamMaxChunkIDLen` | 64 bytes | um identificador de conjunto de chunks fornecido pelo cliente |
| `streamMaxResponseBytes` | 1 MB | uma resposta de polling |
| `streamHoldTimeout` | 20 s | por quanto tempo um polling vazio fica estacionado |
| `streamSessionTTL` | 60 s | nenhum polling por esse período e nenhuma conexão ativa ⇒ a sessão é considerada encerrada |
| `streamSessionGrace` | 10 min | nenhum polling por esse período, com conexões ativas ⇒ a sessão é considerada encerrada |
| `streamSessionSweep` | 20 s | com que frequência o processo de limpeza procura sessões consideradas encerradas |
| `streamMaxFrameBytes` | 64 MB | um quadro em qualquer direção |
| `streamMaxNameLen` | 256 bytes | um nome de stream registrado ou solicitado |
| `streamInQueueDepth` / `streamInQueueBytes` | 256 / 8 MB | quadros recebidos e ainda não consumidos por `Receive` |

### Como escolhê-las

**`streamOutQueueDepth` deliberadamente não é `eventQueueCapacity` (64).** Essa constante foi medida para uma fila drenada uma avaliação por vez, na qual a profundidade só aumentava a latência de cauda. Um polling drena em lotes; portanto, a profundidade aqui precisa comportar a produção durante uma ida e volta — a 5000 quadros/s e com um tempo de ida e volta de 5 ms, cerca de 25 quadros. 256 deixa margem para uma rajada sem bloquear o produtor.

**`streamOutQueueBytes` é o limite que realmente importa**, pois 256 quadros de 1 MB totalizam 256 MB. Ele é a proteção final da memória do host quando um frontend para de coletar.

Duas regras interagem aqui, e é fácil violar a segunda por acidente:

- Os limites de profundidade e de bytes restringem o *acúmulo*.
- **Uma fila vazia sempre aceita um quadro, independentemente do tamanho.** Aplicar incondicionalmente o limite de bytes tornava totalmente impossível enviar um quadro maior que o limite — a condição de espera jamais poderia se tornar verdadeira; assim, `Send` ficava bloqueado para sempre e `TrySend` indicava fila cheia permanentemente. O tamanho do quadro nem sempre fica a critério do chamador; uma struct com um campo `[]byte` é serializada para o tamanho que resultar da serialização.

**`streamMaxResponseBytes` existe por causa do Windows.** O gravador de respostas do WebView2 acumula todo o corpo na memória e só o entrega em `Finish`; portanto, uma resposta sem limite implica uma alocação sem limite nesse ambiente. Aumentá-lo *não* melhora a vazão no Windows — as medições mostram que o gargalo no Windows é por byte, não por resposta: a quantidade de respostas/s varia em um fator de 4 ao longo da varredura de quadros, enquanto MB/s permanece estável em cerca de 90.

**É o limite de entrada que faz o frontend esperar.** No desktop, `deliver` informa que está cheio, o endpoint responde com `429` e o cliente tenta novamente o mesmo quadro ou o sufixo não aceito do lote, com recuo limitado. Isso evita ocupar um slot de solicitação do webview enquanto o manipulador processa os dados pendentes. No modo servidor, o loop de leitura do socket espera e deixa o TCP aplicar contrapressão. Sem o limite, um manipulador que demorasse a chamar `Receive` poderia aumentar o uso de memória do host sem limite.

**Os quadros de controle ignoram os limites de dados, mas têm limites de ciclo de vida independentes.** Perder um quadro de dados sob contrapressão causa uma redução de ritmo; perder uma confirmação de abertura deixa o frontend em `CONNECTING` para sempre, e perder um fechamento faz com que ele considere ativa uma conexão inativa. Portanto, os controles que não são de fechamento têm sua própria fila limitada, separada daquela da qual os fechamentos consomem capacidade, para que uma rajada de aberturas recusadas não consuma a capacidade necessária para uma conexão aceita informar que terminou. Cada sessão também reserva um slot de fechamento por conexão aceita. Quando essa capacidade está ocupada, uma nova abertura recebe contrapressão passível de nova tentativa antes de ser registrada.

**Os limites por sessão também têm equivalentes para toda a aplicação.** Sem eles, cada sessão ou conexão admitida poderia reter toda a sua cota local ao mesmo tempo. Assim, os dados de saída e de entrada compartilham orçamentos separados de 256 MiB / 8192 quadros nos transportes de desktop e servidor. As conexões ativas têm um orçamento de 4096 entradas, e as notificações de fechamento não entregues têm outro de mesmo tamanho. Atingir uma cota compartilhada aplica o mesmo comportamento bloqueante de `Send` ou não bloqueante de `TrySend` aplicado ao atingir uma cota local; além disso, todos os caminhos de drenagem, recebimento, fechamento, falha de gravação e encerramento devolvem suas reservas.

Esses dois orçamentos são deliberadamente separados, em vez de constituírem uma única cota que uma conexão repassa ao seu próprio quadro de fechamento. Cada reserva é liberada por exatamente um proprietário: o slot de uma conexão por `shutdown`, que é executado uma vez, e o slot de um quadro de fechamento por aquilo que descarta o quadro — uma drenagem ou o desmonte de sua sessão. A propriedade que migra entre duas partes precisa ser transferida de forma atômica, e uma revisão anterior que permitia que um fechamento herdasse o slot da conexão vazava um slot permanentemente sempre que o desmonte ocorria entre a tentativa de fechamento e a falha dessa tentativa.

**Quadros do Go transferem a propriedade; quadros do JavaScript são copiados como snapshots.** O `Send` do Go retém o slice do chamador até que o transporte o grave; portanto, após uma chamada bem-sucedida, os chamadores não devem modificar nem reutilizar esse armazenamento. O `send()` do JavaScript copia as entradas binárias mutáveis antes de retornar, em conformidade com a semântica de propriedade do WebSocket nativo. A regra assimétrica evita uma segunda cópia do quadro inteiro dentro do Go e, ao mesmo tempo, mantém previsível a API voltada ao navegador.

**O envio no JavaScript segue o contrato de buffer do WebSocket.** `send()` não pode bloquear; portanto, uma aplicação pode enfileirar dados mais rapidamente do que o canal de solicitações do desktop consegue aceitá-los, assim como pode superar a capacidade de um WebSocket nativo. `bufferedAmount` inclui cada byte retido por esse socket e é o sinal de contrapressão para o chamador; as filas no lado do host permanecem limitadas de forma independente pelos limites acima. Uma falha terminal, o fechamento pelo par ou um `close()` local libera as cargas úteis retidas. O fechamento local também cancela uma solicitação de abertura ou de dados que esteja aguardando em `429` antes de publicar o controle de fechamento reservado; assim, a contrapressão de admissão ou do receptor não pode deixar o socket preso em `CLOSING`.

**A remontagem de fragmentos tem uma cota compartilhada de memória do host.** Cada sessão pode montar um único quadro de até 64 MiB, mas essa cota não pode ser multiplicada por todas as sessões admitidas. Portanto, conjuntos de fragmentos incompletos e sujeitos a nova tentativa compartilham um orçamento de 128 MiB para payloads admitidos. A conclusão de um conjunto retém brevemente tanto suas partes quanto o quadro contíguo montado; assim, mesmo duplicando essa cota lógica, o uso permanece dentro do limite efetivo de memória de 256 MiB. As partes retidas também compartilham uma cota de metadados de 4096 entradas, para que fragmentos pequenos ou vazios não aumentem os mapas e a contabilização de slices sem se aproximar do limite de bytes. Uma solicitação que ultrapassaria qualquer uma das cotas recebe contrapressão que permite nova tentativa; a entrega, rejeição, expiração ou finalização da sessão devolve tanto os bytes quanto as entradas de partes aos orçamentos compartilhados.

**O polling repete somente falhas recuperáveis.** Erros de rede, respostas de tempo limite da solicitação (`408`), respostas de dados antecipados (`425`), contrapressão (`429`) e erros do servidor (`5xx`) usam recuo exponencial de 250 ms até 5 segundos. Outras respostas `4xx` são falhas de protocolo ou de propriedade e fecham imediatamente os Streams da página; `410` é o sinal terminal de encerramento normal para uma sessão desativada. Fechar a última conexão cancela um polling em andamento ou o temporizador de recuo, enquanto uma conexão aberta durante esse encerramento inicia um loop de polling substituto.

**`streamSessionTTL` deve permanecer confortavelmente acima de `streamHoldTimeout`**; caso contrário, uma sessão seria removida enquanto seu próprio polling estivesse legitimamente aguardando.

Se você estiver ajustando para uma carga de trabalho com muitas mensagens pequenas, o limite de profundidade será atingido primeiro; para payloads grandes, será o limite de bytes. Nenhum deles precisa ser alterado em um aplicativo típico — no macOS, os valores padrão sustentam 634000 quadros/s e 2100 MB/s.

## Ciclo de vida de conexões e sessões

Uma **sessão** corresponde a um carregamento de página em uma janela e é identificada por um id gerado pelo cliente (como o `clientId` do runtime). As sessões são criadas de forma tardia pela primeira solicitação que chegar. Quando uma plataforma não consegue identificar a janela solicitante (`windowID == 0`), os ids das sessões continuam sujeitos a um limite global, mas suas gerações deliberadamente não são comparadas: elas podem pertencer a clientes de navegador independentes, com contadores de geração não relacionados. Essas sessões expiram por fechamento ou TTL, em vez de substituírem umas às outras.

Três mecanismos fecham os recursos, na ordem da rapidez com que detectam a situação:

1. **Um polling de sessão mais recente para uma janela substitui gerações anteriores.** Um recarregamento fornece à página um novo id de sessão e incrementa uma geração armazenada no `sessionStorage` dessa janela. O mesmo valor é espelhado em `window.name`, que sobrevive a recarregamentos quando o armazenamento está desabilitado, e é ancorado em `performance.timeOrigin` (ou `Date.now()` em mecanismos mais antigos), para que a limpeza de ambos os armazenamentos não reinicie a contagem em um. Cada solicitação carrega o id e a geração da sessão. Um polling desativa apenas gerações anteriores da página, para que o agendamento do servidor não faça uma solicitação atrasada da página anterior parecer mais recente que sua substituta. Se uma política bloquear tanto o armazenamento quanto `window.name`, a ordenação recorrerá ao relógio da página e, portanto, dependerá de as páginas posteriores receberem uma origem de tempo posterior. O gerenciador mantém, por janela, uma marca d'água da geração desativada, para que uma solicitação já em andamento não possa recriar a página antiga sem que seja necessário reter todos os ids históricos de sessão. As conexões da sessão anterior são fechadas imediatamente.
2. **A destruição da janela** remove todas as sessões dessa janela, como `eventPayloadStore.dropWindow`.
3. **A remoção por TTL** trata todo o restante — um renderizador que falhou ou uma máquina suspensa.

Somente os dois primeiros mecanismos desativam a geração da página. A limpeza por TTL remove a sessão ociosa sem avançar a marca d'água da geração desativada: uma página para de fazer polling quando sua última conexão é fechada, mas essa mesma página, ainda carregada, deve poder abrir outro stream posteriormente. Uma geração que realmente foi substituída continua bloqueada porque o polling da página mais recente avança a marca d'água antes de a sessão antiga ser removida.

**As WebViews da Apple informam solicitações canceladas.** No macOS e no iOS, o callback `stopURLSchemeTask` do WebKit cancela o contexto da solicitação correspondente; assim, um polling pertencente a uma página da qual se saiu por navegação é desbloqueado imediatamente. O registro usa como chave a identidade retida da tarefa nativa e remove as entradas quando o processamento das solicitações as fecha. O Linux e o Windows ainda não expõem um callback equivalente para cancelamento antecipado na ponte atual; nessas plataformas, uma solicitação em espera permanece assim até o término do período de retenção. A regra 1 faz com que a *conexão* seja fechada prontamente de qualquer forma. Nos demais casos, o cancelamento aparece como `EPIPE` no Linux e somente em `Finish` no Windows.

## Seleção do transporte

`Stream(name)` consulta `window._wails.streamFactory`. Builds de servidor instalam uma implementação que retorna um `WebSocket` real; builds de webview a deixam indefinida e recebem o cliente de polling.

A fábrica **deve** ser instalada antes da execução do corpo de qualquer módulo, pois os bindings gerados criarão streams no escopo do módulo. `custom.js` não pode fazer isso — `loadOptionalScript` realiza uma solicitação HEAD e depois adiciona uma tag `<script>`, portanto ocorre tarde demais. Em vez disso, a fábrica é adicionada ao início do bundle do runtime quando ele é servido (`stream_prelude_server.go`), o que é síncrono por definição: as dependências de módulos ES são avaliadas antes dos módulos que as importam.

Se você adicionar um terceiro transporte, coloque-o também no preâmbulo. Não ceda à tentação de voltar para `custom.js`.

## O que ainda não está concluído

|  | status |
| --- | --- |
| Cancelamento de solicitações pela camada da plataforma | **Concluído para Apple; pendente para Linux/Windows** — veja acima |
| Constantes de buffer como opções | não concluído; apenas em tempo de compilação |
| Streams tipados | deliberadamente não concluído — por decisão, os quadros são `[]byte` |
| Pipelining (um segundo polling em andamento) | não concluído; exigiria remontagem ordenada em JS |
| Equidade por conexão | não concluído — as conexões de uma janela compartilham uma fila; portanto, uma conexão que a inunda torna as demais mais lentas |
| Coalescência de quadros de JS→Go | **concluído** — os quadros acumulados enquanto uma solicitação está em andamento são enviados em lotes limitados; conexões com pouca carga ainda enviam um quadro por POST |
| Taxa de transferência no Windows | ~100 MB/s, limitada pelo marshaling de `WebResourceRequested`. Buffers compartilhados (`PostSharedBufferToScript`) são a possível solução; há bindings em `internal/webview2/pkg/webview2/`, mas eles não estão integrados a `pkg/edge` |
| `wails3 dev` / Vite | **funciona** — verificado com um projeto `vanilla-js` gerado: o servidor de desenvolvimento do Vite faz proxy em `/`, e `/wails/stream/*` é correspondido pelo middleware do servidor de assets antes do proxy; portanto, os streams não são afetados |
| Múltiplas janelas | não testado sob carga, embora, por definição, as sessões tenham escopo de janela |

## O pacote de frontend no modo de desenvolvimento

Um projeto gerado importa `@wailsio/runtime` do **npm**, e não do `/wails/runtime.js` incluído no bundle servido pelo servidor de assets. Com `wails3 dev`, o Vite o resolve a partir de `node_modules`; portanto, um aplicativo compilado com um runtime publicado não verá adições no lado do cliente feitas em uma branch.

Enquanto os streams ainda não tiverem sido lançados, aponte um aplicativo de teste para o pacote desta cópia de trabalho:

```bash
task v3:install-runtime -- ./path/to/your-app/frontend
```

Isso primeiro recompila `dist/`, portanto sempre instala o código-fonte atual. Desfaça essa alteração com  
`npm install @wailsio/runtime@latest` no mesmo diretório.

Observe que há **duas** saídas de build do cliente e é fácil recompilar uma, mas não a  
outra: `task v3:runtime:build:package` gera o `dist/` do pacote npm (que o frontend de um aplicativo  
importa), enquanto `task v3:runtime:build:assets` gera  
`bundledassets/runtime.js` (que o webview carrega do servidor de recursos). Uma alteração em  
`stream.ts` exige ambas.

## Testes

```bash
go test ./pkg/application/ -run TestStream -race        # protocol, ordering, backpressure
go test -tags server ./pkg/application/ -run TestServerMode
pnpm --dir v3/internal/runtime/desktop/@wailsio/runtime test
```

O teste de ordenação é o que importa: oito goroutines enviam simultaneamente usando um  
contador distribuído enquanto o bloqueio da fila está mantido, e a ordem de esvaziamento deve corresponder  
exatamente à ordem de aceitação. Se esse teste falhar, significa que a invariante de um único consumidor da fila foi violada.

Harness de teste de carga:

```bash
go run ./tests/stream-performance -duration 20s              # full sweep
go run ./tests/stream-performance -upload -duration 10s      # JS→Go matrix
go run ./tests/stream-performance -reloads 6                 # connection lifecycle
```

No Windows, ele deve ser executado na sessão interativa do console — uma simples invocação por SSH é encerrada na  
sessão 0 com saída de comprimento zero —, e o binário deve ser colocado em um local que tanto a conta  
SSH quanto a conta do console possam ler, pois `C:\Users\<user>` tem sua ACL restrita ao proprietário.
