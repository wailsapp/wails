---
title: "Registro de alterações"
description: "Histórico de versões e notas de lançamento do Wails v3"
slug: "changelog"
sourcePath: "changelog.md"
---

Legenda:

-  - macOS
- ⊞ - Windows
- 🐧 - Linux

/_-- Todas as alterações relevantes deste projeto serão documentadas neste arquivo.

O formato é baseado em [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), e este projeto segue o [Versionamento Semântico](https://semver.org/spec/v2.0.0.html).

- `Added` para novos recursos.
- `Changed` para alterações em funcionalidades existentes.
- `Deprecated` para recursos que serão removidos em breve.
- `Removed` para recursos já removidos.
- `Fixed` para quaisquer correções de bugs.
- `Security` em caso de vulnerabilidades.

_/

/_   * NÃO ATUALIZE ESTE ARQUIVO *   As atualizações devem ser adicionadas a `v3/UNRELEASED_CHANGELOG.md`   Obrigado! _/

## [Não lançado]

## v3.0.0-beta.21 - 2026-09-13

## Adicionado

- Disponibilização da documentação do Wails v3 com o M-Press em [PR](https://github.com/wailsapp/wails/pull/6116) por @leaanthony

## Corrigido

- Análise dos valores de slug em JSON no frontmatter MPD para a geração do changelog em [PR](https://github.com/wailsapp/wails/pull/6118) por @leaanthony
- O atualizador limpa as variáveis de ambiente auxiliares e reinicia o destino original após falhas de backup em [PR](https://github.com/wailsapp/wails/pull/6080) por @cnmax
- Inicialização do manipulador de sinais padrão durante App.Run em [PR](https://github.com/wailsapp/wails/pull/6098) por @leaanthony
- O menu do Windows trata menus nil, libera os recursos substituídos e redesenha a barra de menus em [PR](https://github.com/wailsapp/wails/pull/6112) por @taliesin-ai
- Restauração do empacotamento MSIX para novos projetos que usam uma configuração YAML compartilhada em [PR](https://github.com/wailsapp/wails/pull/6115) por @leaanthony
- Correção da falha no carregamento dos bindings JavaScript e TypeScript gerados quando os criadores de modelos genéricos fazem referência a declarações auxiliares posteriores, além da prevenção de estouros de pilha ao criar modelos genéricos mutuamente dependentes (#6062)

## v3.0.0-beta.20 - 2026-09-10

## Alterado

- Atualiza os links da vitrine Clave para o site e o repositório atuais no [PR](https://github.com/wailsapp/wails/pull/6082) por @01xR4in

## Corrigido

- Cancela solicitações de ativos do Windows que foram abortadas, incluindo solicitações de workers, preservando os manipuladores keepalive entre navegações. Encaminha os contextos de solicitações nativas pelo wrapper do aplicativo nas plataformas Apple. (#5963, #5969)
- Preserva as entradas do registro de alterações durante pushes concorrentes por meio de novas tentativas no [PR](https://github.com/wailsapp/wails/pull/6094) por @leaanthony
- Corrige a falha de `go mod vendor` com `pattern arm64/WebView2Loader.dll: no matching files found` em todas as plataformas, removendo os embeds que faziam referência a binários nunca distribuídos no módulo, corrigindo [#5782](https://github.com/wailsapp/wails/issues/5782) e [#5376](https://github.com/wailsapp/wails/issues/5376), no [PR](https://github.com/wailsapp/wails/pull/6031) por @Grantmartin2002

## Removido

- Remove o suporte ao carregador nativo do WebView2, substituído pelo carregador em Go puro. Isso elimina os binários `WebView2Loader.dll` incorporados e a dependência `github.com/jchv/go-winloader`. A tag de compilação `native_webview2loader` continua sendo aceita e não causa mais erros, mas não tem efeito nas compilações da v3, no [PR](https://github.com/wailsapp/wails/pull/6031) por @Grantmartin2002
- Remove tags de compilação não utilizadas e a opção de FPS do guia da API do macOS no [PR](https://github.com/wailsapp/wails/pull/6097) por @leaanthony

## v3.0.0-beta.19 - 2026-09-09

## Adicionado

- Restringe as APIs privadas do macOS por meio de tags de compilação, exigindo adesão explícita para usá-las — consulte a [documentação](https://v3.wails.io/features/browser/integration), a [documentação](https://v3.wails.io/features/environment/info), a [documentação](https://v3.wails.io/features/windows/basics), a [documentação](https://v3.wails.io/features/windows/frameless), a [documentação](https://v3.wails.io/features/windows/notch-windows), a [documentação](https://v3.wails.io/features/windows/options), a [documentação](https://v3.wails.io/guides/build/macos), a [documentação](https://v3.wails.io/guides/build/private-macos-apis) e a [documentação](https://v3.wails.io/reference/overview) no [PR](https://github.com/wailsapp/wails/pull/6087) por @leaanthony

## Corrigido

- Rejeita solicitações ao runtime com mais de 64 MiB usando HTTP 413 no [PR](https://github.com/wailsapp/wails/pull/6091) por @leaanthony

## Segurança

- Reforça a segurança das origens do MCP e do acesso remoto com autenticação por token no [PR](https://github.com/wailsapp/wails/pull/6092) por @leaanthony

## v3.0.0-beta.18 - 2026-09-08

## Corrigido

- Corrige o vazamento de memória de Calloc no Linux e no Darwin usando receptores de ponteiro no [PR](https://github.com/wailsapp/wails/pull/6083) por @4RH1T3CT0R7

## v3.0.0-beta.17 - 2026-09-06

## Corrigido

- Windows: um `GetRequest` com falha ou nil no manipulador WebResourceRequested não encerra mais o processo (`log.Fatal` / pânico por desreferência de nil) — em vez disso, a solicitação é descartada e registrada no log, no [PR](https://github.com/wailsapp/wails/pull/6006) por @midagedev

## v3.0.0-beta.16 - 2026-08-29

## Alterado

- Solicita a senha de notarização em uma nova janela de terminal no [PR](https://github.com/wailsapp/wails/pull/6029) por @leaanthony

## Corrigido

- Trata corretamente os tipos de clique no ícone da bandeja do sistema no macOS no [PR](https://github.com/wailsapp/wails/pull/5919) por @ChewbaccaCookie
- A CI remove os repositórios apt da Microsoft não utilizados antes da atualização no [PR](https://github.com/wailsapp/wails/pull/6041) por @Grantmartin2002

## v3.0.0-beta.15 - 2026-08-27

## Corrigido

- Aumenta para 60 segundos o tempo limite de incorporação do WebView2 no [PR](https://github.com/wailsapp/wails/pull/6043) por @Grantmartin2002

## v3.0.0-beta.14 - 2026-08-26

## Corrigido

- Nomeia corretamente no macOS os pressionamentos de teclas formados por Control e uma letra no [PR](https://github.com/wailsapp/wails/pull/6032) por @taliesin-ai
- Corrige os ícones ICO da bandeja do sistema e faz com que acompanhem o tema da barra de tarefas no Windows no [PR](https://github.com/wailsapp/wails/pull/6016) por @nik9play

## v3.0.0-beta.13 - 2026-08-25

## Corrigido

- Mantém o processamento de tarefas da thread principal no macOS enquanto um loop modal está em execução no [PR](https://github.com/wailsapp/wails/pull/6026) por @leaanthony
- Permite que o armazenamento seguro para dispositivos móveis retorne falhas e bloqueie o acesso em caso de falha no [PR](https://github.com/wailsapp/wails/pull/5923) por @mortenolsrud
- Executa os hooks de eventos do aplicativo mesmo quando nenhum listener está registrado no [PR](https://github.com/wailsapp/wails/pull/5999) por @archy-rock3t-cloud
- Corrige erros de digitação em comentários e na documentação localizada no [PR](https://github.com/wailsapp/wails/pull/6023) por @haoku123
- Remove os binários pré-compilados do macOS enviados ao repositório em `v3/examples` no [PR](https://github.com/wailsapp/wails/pull/6025) por @4RH1T3CT0R7

## v3.0.0-beta.12 - 2026-08-21

## Adicionado

- Adiciona janelas de notificação para o notch do macOS com exemplo de ciclo de vida e telemetria — consulte a [documentação](https://v3.wails.io/features/windows/notch-windows) no [PR](https://github.com/wailsapp/wails/pull/6010) por @leaanthony
- Adiciona suporte a janelas NSPanel no macOS, com novas opções e integração nativa — consulte a [documentação](https://v3.wails.io/features/windows/options) na [PR](https://github.com/wailsapp/wails/pull/6008) de @leaanthony

## Corrigido

- Impede que SQLite Prepare fique bloqueado durante chamadas simultâneas na [PR](https://github.com/wailsapp/wails/pull/5998) de @archy-rock3t-cloud

## v3.0.0-beta.11 - 2026-08-20

## Removido

- Remove da documentação o rastreador de implementação obsoleto na [PR](https://github.com/wailsapp/wails/pull/6005) de @leaanthony

## v3.0.0-beta.10 - 2026-08-19

## Corrigido

- Corrige o descarte, pelo host GTK4 no Linux, dos argumentos de inicialização de protocolo personalizado e associação de arquivos na [PR](https://github.com/wailsapp/wails/pull/6000) de @midagedev
- Processa corretamente linhas excluídas e correções da mesma origem na validação do changelog, na [PR](https://github.com/wailsapp/wails/pull/5993) de @taliesin-ai

## v3.0.0-beta.9 - 2026-08-16

## Adicionado

- Adiciona o servidor MCP seguro wails3 para gerenciamento de projetos assistido por agentes na [PR](https://github.com/wailsapp/wails/pull/5896) de @leaanthony
- Adiciona documentação sobre modelos em bindings — consulte a [documentação](https://v3.wails.io/features/bindings/models) na [PR](https://github.com/wailsapp/wails/pull/5988) de @taliesin-ai
- Adiciona suporte à instalação com rpm-ostree em sistemas Linux atômicos na [PR](https://github.com/wailsapp/wails/pull/5987) de @leaanthony
- Adiciona geração e publicação nativas de um gráfico semanal do histórico de estrelas — consulte a [documentação](https://v3.wails.io/credits), a [documentação](https://v3.wails.io/de/credits), a [documentação](https://v3.wails.io/fr/credits), a [documentação](https://v3.wails.io/id/credits), a [documentação](https://v3.wails.io/ja/credits), a [documentação](https://v3.wails.io/ko/credits), a [documentação](https://v3.wails.io/pt/credits), a [documentação](https://v3.wails.io/ru/credits), a [documentação](https://v3.wails.io/zh-cn/credits) e a [documentação](https://v3.wails.io/zh-tw/credits) na [PR](https://github.com/wailsapp/wails/pull/5986) de @leaanthony
- Adiciona o pacote mac exclusivo do Darwin para resolver recursos do pacote de aplicativos — consulte a [documentação](https://v3.wails.io/guides/build/macos) na [PR](https://github.com/wailsapp/wails/pull/5965) de @leaanthony
- Adiciona a página do showcase do Condui e sua entrada no índice — consulte a [documentação](https://v3.wails.io/community/showcase/condui) e a [documentação](https://v3.wails.io/community/showcase) na [PR](https://github.com/wailsapp/wails/pull/5962) de @mgueregath
- Adiciona a página do showcase do Redis Viewer, com capturas de tela e link para o projeto — consulte a [documentação](https://v3.wails.io/community/showcase) e a [documentação](https://v3.wails.io/community/showcase/redisviewer) na [PR](https://github.com/wailsapp/wails/pull/5984) de @redisviewer

## Alterado

- Atualiza os sinalizadores do aplicativo GTK para G<em>APPLICATION</em>NON_UNIQUE no Linux, na [PR](https://github.com/wailsapp/wails/pull/5971) de @overlordtm
- Registra eventos de janela ausentes no nível de depuração, em vez de aviso, na [PR](https://github.com/wailsapp/wails/pull/5914) de @julianstorer

## Corrigido

- Permite que os atalhos de teclado registrados no macOS tenham precedência sobre a webview, na [PR](https://github.com/wailsapp/wails/pull/5902) de @julianstorer
- Repara links quebrados na barra lateral da documentação, na [PR](https://github.com/wailsapp/wails/pull/5937) de @northes
- Cancela os contextos de requisições de ativos no macOS e iOS quando o WebKit interrompe a tarefa correspondente do esquema personalizado (#5963)
- Processa a mensagem WindowSetFullscreenButtonEnabled na [PR](https://github.com/wailsapp/wails/pull/5976) de @archy-rock3t-cloud
- Importa Fragment no modelo preact-ts para resolver a falha de compilação, na [PR](https://github.com/wailsapp/wails/pull/5979) de @haoku123
- Impede que aplicativos legados GTK3 executados apenas como serviço falhem quando a descoberta de telas é executada antes que uma janela ou tela ativa esteja disponível (#5966)
- Permite que execuções de lançamento com versão explícita prossigam quando o changelog de itens ainda não lançados estiver vazio (#5977)

## Segurança

- Atualiza os lockfiles do nanoid do site para a versão corrigida 3.3.18, resolvendo alertas de segurança na [PR](https://github.com/wailsapp/wails/pull/5985) de @taliesin-ai

## v3.0.0-beta.8 - 2026-08-12

## Adicionado

- Adiciona a geração de URLs de documentação às entradas automáticas do changelog, na [PR](https://github.com/wailsapp/wails/pull/5957) de @taliesin-ai
- Adiciona Streams: fluxos de bytes bidirecionais entre Go e JavaScript, com o modelo de programação WebSocket e sem socket de escuta. Declare um fluxo em Go com `app.HandleStream(name, handler)` e conecte-se pelo frontend com `Stream(name)`, que retorna um objeto com o formato de `WebSocket`. O tráfego Go→JS é transportado por uma única sondagem mantida por janela por meio do servidor de ativos, e o tráfego JS→Go, por um POST normal; nenhum deles vincula uma porta TCP e nenhum passa por `evaluateJavaScript`. Em builds de servidor (`-tags server`), o mesmo manipulador é disponibilizado por meio de um WebSocket real, de modo que o código do aplicativo é idêntico em todos os builds. Por @leaanthony
- Move a entrada do changelog da caixa de correio para a seção de itens ainda não lançados, na [PR](https://github.com/wailsapp/wails/pull/5935) de @leaanthony

## Alterado

- Atualiza a geração automática da barra lateral da documentação e a derivação do tipo de autor do blog, na [PR](https://github.com/wailsapp/wails/pull/5938) de @leaanthony

## Corrigido

- A inicialização do WebView2 passa a usar um prazo limite e uma bomba de mensagens, na [PR](https://github.com/wailsapp/wails/pull/5952) de @leaanthony
- O teste de cookies do WebView2 é ignorado na CI, a menos que seja habilitado explicitamente, e sua execução é vinculada à thread atual do sistema operacional, na [PR](https://github.com/wailsapp/wails/pull/5951) de @leaanthony
- Os construtores de menus do Windows restauram os IDs de comando dos itens pai de submenus, na [PR](https://github.com/wailsapp/wails/pull/5944) de @gilad-ch
- Alinha a imagem oficial de compilação cruzada ao requisito mínimo de GTK 4.14+ para suporte ao Linux (#5928)
- Configura o projeto Xcode para iOS para preservar os sinalizadores herdados do vinculador e adicionar -ObjC, na [PR](https://github.com/wailsapp/wails/pull/5915) de @mortenolsrud
- Corrige a rotatividade excessiva de conexões TCP no proxy de ativos `wails3 dev` em frontends grandes, que poderia esgotar as portas efêmeras do host e fazer com que processos não relacionados falhassem com `EADDRNOTAVAIL`
- Enfileira o JavaScript de eventos de cada janela para despacho ordenado e controle de contrapressão, na [PR](https://github.com/wailsapp/wails/pull/5934) de @leaanthony

## Removido

- Remove o pipeline de lançamento de binários para desktop: os lançamentos da v3 usam apenas tags, e a CLI `wails3` é instalada com `go install`. Exclui `release-v3.yml` e a etapa nightly que o acionava, na [PR](https://github.com/wailsapp/wails/pull/5946) de @leaanthony

## v3.0.0-beta.7 - 2026-08-11

## Adicionado

- Adiciona uma preferência de reprodução automática no macOS para desativar a exigência de ação do usuário para a reprodução de mídia em [PR](https://github.com/wailsapp/wails/pull/5512) por @Eyalm321
- Move a entrada do changelog referente à caixa de mensagens para a seção Não lançado no [PR](https://github.com/wailsapp/wails/pull/5935) por @leaanthony

## Alterado

- A animação de zoom do macOS passa a usar CADisplayLink ou NSTimer para proporcionar um desempenho mais fluido em [PR](https://github.com/wailsapp/wails/pull/5945) por @savely-krasovsky

## Corrigido

- Configura o projeto Xcode para iOS de modo a preservar as flags herdadas do vinculador e adicionar -ObjC em [PR](https://github.com/wailsapp/wails/pull/5915) por @mortenolsrud
- Corrige a rotatividade excessiva de conexões TCP no proxy de ativos `wails3 dev` em frontends grandes, que poderia esgotar as portas efêmeras do host e fazer processos não relacionados falharem com `EADDRNOTAVAIL`
- Enfileira o JavaScript de eventos de cada janela para garantir o despacho ordenado e o controle de contrapressão em [PR](https://github.com/wailsapp/wails/pull/5934) por @leaanthony

### Adicionado

- Implementa uma caixa de mensagens FIFO assíncrona genérica para a entrega ordenada de eventos em [PR](https://github.com/wailsapp/wails/pull/5851) por @savely-krasovsky e @DevLumuz

## v3.0.0-beta.6 - 2026-08-09

## Adicionado

- Implementa armazenamento limitado no lado do host para eventos grandes demais e entrega ordenada ao JavaScript em [PR](https://github.com/wailsapp/wails/pull/5930) por @leaanthony
- Implementa o efeito de salto do ícone no Dock do macOS para sinalizar uma janela em [PR](https://github.com/wailsapp/wails/pull/5921) por @julianstorer

## Corrigido

- O servidor de ativos passa a preservar os erros do detector de tipo de conteúdo e os prefixos ainda não gravados durante o esvaziamento do buffer em [PR](https://github.com/wailsapp/wails/pull/5931) por @leaanthony
- Impede que aplicativos macOS falhem ao substituir o menu do aplicativo a partir de um callback do Wails
- Corrige menus nativos ilegíveis no Windows 10 1809 / Windows Server 2019 (build 17763). As exportações de modo escuro de uxtheme estavam condicionadas ao build 18334, portanto a adesão ao modo escuro no nível do aplicativo nunca era executada nesses hosts: o fundo do menu era pintado de escuro, mas o Windows continuava desenhando o texto do menu com o tema claro, deixando texto escuro sobre fundo escuro. Os ordinais existem desde 17763, portanto a condição agora corresponde a essa versão.
- Corrige `w32.GetStockObject`, que chamava `GetDeviceCaps` em vez de `GetStockObject`, fazendo com que retornasse 0 para todos os objetos predefinidos.
- Melhora o tratamento e o relatório de erros de download do bootstrapper do WebView2 em [PR](https://github.com/wailsapp/wails/pull/5924) por @jannskiee

## v3.0.0-beta.5 - 2026-08-07

## Corrigido

- A ativação de aplicativos no macOS passa a respeitar a política de ativação somente para aplicativos regulares em [PR](https://github.com/wailsapp/wails/pull/5897) por @julianstorer
- Protege contra janelas GTK não inicializadas em builds para Linux em [PR](https://github.com/wailsapp/wails/pull/5898) por @julianstorer
- Define explicitamente uma cor de fundo opaca para janelas WebKit no Linux antes do carregamento da URL em [PR](https://github.com/wailsapp/wails/pull/5899) por @julianstorer

## v3.0.0-beta.4 - 2026-08-05

## Alterado

- As tarefas de build do Android passam a usar arm64 por padrão, e deploy-emulator seleciona a arquitetura do host em [PR](https://github.com/wailsapp/wails/pull/5890) por @mortenolsrud

## Corrigido

- Preserva o estado de zoom das janelas do macOS durante o arraste e reduz o movimento em [PR](https://github.com/wailsapp/wails/pull/5900) por @leaanthony
- Corrige o build do Windows no modo servidor adicionando `!server` à restrição de build do arquivo `webview_window_windows_nonclient.go`

## v3.0.0-beta.3 - 2026-08-03

## Adicionado

- Documenta nos detalhes de implementação a conclusão da verificação beta da Fase 10 em [PR](https://github.com/wailsapp/wails/pull/5881) por @leaanthony

## Corrigido

- Passa o identificador da janela para a API de modo escuro do Windows e valida os argumentos em [PR](https://github.com/wailsapp/wails/pull/5877) por @leaanthony
- Centraliza a resolução do estado dos botões da barra de título para janelas sem moldura no macOS em [PR](https://github.com/wailsapp/wails/pull/5870) por @taliesin-ai
- Impede que o texto dos menus nativos fique ilegível quando um aplicativo Windows solicita o modo escuro enquanto o tema de aplicativos do Windows está claro. O menu agora usa o fundo nativo claro correspondente até que o Windows consiga renderizar o texto do menu para o modo escuro.
- Corrige menus nativos ilegíveis no Windows 10 1809 / Windows Server 2019 (build 17763). As exportações de modo escuro de uxtheme estavam condicionadas ao build 18334, portanto a adesão ao modo escuro no nível do aplicativo nunca era executada nesses hosts: o fundo do menu era pintado de escuro, mas o Windows continuava desenhando o texto do menu com o tema claro, deixando texto escuro sobre fundo escuro. Os ordinais existem desde 17763, portanto a condição agora corresponde a essa versão.

## v3.0.0-beta.2 - 2026-08-02

## Alterado

- Promove a v3 de alfa para beta
- Documenta os padrões inteligentes da bandeja do sistema e o comportamento de ocultação automática do pop-up, com cobertura de regressão para a seleção do manipulador de cliques (#5840).
- O atualizador do GitHub passa a excluir por padrão os ativos de instaladores do Windows em [PR](https://github.com/wailsapp/wails/pull/5861) por @leaanthony
- Adiciona suporte a cantos arredondados, quadrados e com raio personalizado em janelas sem moldura do macOS em [PR](https://github.com/wailsapp/wails/pull/5866) por @leaanthony

## Corrigido

- Informa os tamanhos atuais das janelas GTK4 e emite eventos de redimensionamento, maximização, minimização e estado de tela cheia a partir da superfície configurada (#5830).
- Corrige uma falha do WebKit no Linux ao enviar Blob ou FormData em solicitações fetch em [PR](https://github.com/wailsapp/wails/pull/5854) por @taliesin-ai
- O shim de fetch passa undefined quando os cabeçalhos de Blob/FormData estão ausentes em [PR](https://github.com/wailsapp/wails/pull/5865) por @leaanthony

## v3.0.0-alpha2.122 - 2026-08-01

## Adicionado

## Alterado

- Adiciona suporte a cantos arredondados, quadrados e com raio personalizado em janelas sem moldura do macOS em [PR](https://github.com/wailsapp/wails/pull/5866) por @leaanthony

## Corrigido

- O shim de fetch passa undefined quando os cabeçalhos de Blob/FormData estão ausentes em [PR](https://github.com/wailsapp/wails/pull/5865) por @leaanthony

## v3.0.0-alpha2.121 - 2026-07-31

## Adicionado

- Adiciona suporte ao empacotamento de DMG no macOS, com novas opções e tarefas de compilação, no [PR](https://github.com/wailsapp/wails/pull/5857) de @leaanthony

## Alterado

- O atualizador do GitHub exclui, por padrão, os artefatos do instalador do Windows no [PR](https://github.com/wailsapp/wails/pull/5861) de @leaanthony

## Corrigido

- Corrige uma falha do WebKit no Linux ao enviar Blob ou FormData em requisições fetch no [PR](https://github.com/wailsapp/wails/pull/5854) de @taliesin-ai

## v3.0.0-alpha2.120 - 2026-07-31

## Adicionado

- Implementa ações de maximização ou minimização ao clicar duas vezes na barra de título no macOS, no [PR](https://github.com/wailsapp/wails/pull/5853) de @taliesin-ai

## Alterado

- Atualiza o tutorial do serviço de QR para usar NewServiceWithOptions e adiciona espaçamento no [PR](https://github.com/wailsapp/wails/pull/5849) de @jeongkyu

## Corrigido

- Mantém o WKWebView responsivo durante o zoom no macOS, no [PR](https://github.com/wailsapp/wails/pull/5856) de @leaanthony
- Corrige as consultas de tamanho da janela no GTK4 e emite eventos de redimensionamento e de estado maximizado, minimizado e em tela cheia a partir do `GdkSurface` configurado.

## v3.0.0-alpha2.119 - 2026-07-27

## Corrigido

- Atualiza a documentação em vários idiomas para incluir diagramas de arquitetura no [PR](https://github.com/wailsapp/wails/pull/5833) de @taliesin-ai

## v3.0.0-alpha2.118 - 2026-07-26

## Adicionado

- Fornece caminhos padrão para as entradas e saídas da geração de ícones no [PR](https://github.com/wailsapp/wails/pull/5825) de @taliesin-ai
- Adiciona módulos de entrada do código-fonte a sideEffects do package.json do runtime no [PR](https://github.com/wailsapp/wails/pull/5797) de @savely-krasovsky
- Adiciona uma seção sobre licença e procedência ao guia de contribuição no [PR](https://github.com/wailsapp/wails/pull/5816) de @taliesin-ai

## Corrigido

- Aplica CSS sem moldura com escopo no GTK4 para remover o raio das bordas no [PR](https://github.com/wailsapp/wails/pull/5800) de @savely-krasovsky
- Trata adequadamente as falhas ao obter a posição do cursor no Windows para menus pop-up e enumeração de telas, no [PR](https://github.com/wailsapp/wails/pull/5789) de @wayneforrest
- A caixa de diálogo para abrir arquivos do macOS filtra corretamente as extensões e valida os arquivos permitidos pelo sufixo, no [PR](https://github.com/wailsapp/wails/pull/5678) de @phergul
- Protege a inicialização do modo escuro no Windows contra chamadas de API nil no [PR](https://github.com/wailsapp/wails/pull/5793) de @roachadam
- Corrige uma falha de compilação para 32 bits no atualizador: a constante `maxArchiveTotalSize` (2 GiB) excedia a capacidade do `int` da plataforma quando passada para `fmt.Errorf` em `GOARCH=386`. Agora, ela tem explicitamente o tipo `int64`.
- Corrige um pânico causado por ponteiro nil na inicialização quando uma janela usa uma barra de título escura (ou escura conforme o sistema) em builds do Windows que não carregam as APIs de modo escuro de uxtheme, como Windows 10 1809 / Windows Server 2019 (build 17763). As chamadas a `AllowDarkModeForWindow` na configuração do tema da janela agora são protegidas contra nil, como já ocorre em `w32.SetMenuTheme`.

## v3.0.0-alpha2.117 - 2026-07-08

## Adicionado

- Implementa lógica personalizada de teste de acerto para regiões não clientes no Windows, no [PR](https://github.com/wailsapp/wails/pull/5462) de @savely-krasovsky

## Alterado

- Configura a detecção da escala do monitor pelo WebView2 com base em UseVisualHosting no [PR](https://github.com/wailsapp/wails/pull/5761) de @wayneforrest

## v3.0.0-alpha2.116 - 2026-07-07

## Adicionado

- Atualiza a documentação de perguntas frequentes para se concentrar nos recursos e nas orientações do Wails v3, no [PR](https://github.com/wailsapp/wails/pull/5763) de @taliesin-ai

## v3.0.0-alpha2.115 - 2026-07-06

## Corrigido

- Corrige o problema em que `Menu.Update()` não reconstruía o menu nativo no GTK4 para Linux (#5659, diagnosticado e corrigido de forma independente por @puneetdixit200 em #5539)
- Corrige uma falha ao enumerar as telas do macOS após uma mudança de monitor, copiando as strings de ID/nome das telas e capturando uma cópia instantânea da contagem (#5565, diagnosticado e corrigido de forma independente por @x-haose em #5584)
- Corrige uma falha no Windows quando `WM_ERASEBKGND` desenha um fundo sólido durante uma transição de minimização/restauração na qual `GetClientRect` retorna nil (proteção relatada por @sinspired em #5636)
- Corrige o erro em que a mensagem de erro dos bindings do frontend sempre era analisada como texto, por @mbaklor em #5690
- Corrige uma falha de compilação no Windows ao usar a tag de compilação `server`, causada pela ausência, nos arquivos da GUI do Windows, da restrição de compilação `!server` que os equivalentes para macOS e Linux já possuem (#5680)

## v3.0.0-alpha2.114 - 2026-07-05

## Adicionado

- Implementa o protocolo Update Manifest e o provedor de endpoint no [PR](https://github.com/wailsapp/wails/pull/5720) de @taliesin-ai

## Alterado

- Incorpora o binding `webview2` ao módulo v3 como `v3/internal/webview2`, removendo o módulo independente, seus fluxos de trabalho de lançamento/sincronização nightly e a alternância de versões no go.mod (o v3 é seu único consumidor), no [PR](https://github.com/wailsapp/wails/pull/5711) de @taliesin-ai

## Corrigido

- Move a detecção da escala do monitor pelo WebView2 e a correção da ressincronização do host após mudança de DPI para a seção Não lançado, no [PR](https://github.com/wailsapp/wails/pull/5750) de @taliesin-ai
- Atualiza o marshaling COM do WebView2 para parâmetros float64 e BOOL no [PR](https://github.com/wailsapp/wails/pull/5741) de @wayneforrest
- Evita pânico e desreferência de nil durante a atualização e a destruição do ícone da bandeja do sistema no Windows, no [PR](https://github.com/wailsapp/wails/pull/5703) de @wayneforrest
- Corrige o problema em que janelas ocultas não voltavam a ser ocultadas corretamente no Windows, no [PR](https://github.com/wailsapp/wails/pull/5743) de @wayneforrest
- Sincroniza a visibilidade do controlador WebView2 com a minimização, maximização e restauração da janela no [PR](https://github.com/wailsapp/wails/pull/5742) de @wayneforrest

### Corrigido

- Reativa a detecção da escala do monitor pelo WebView2 e condiciona a ressincronização do host à mudança de DPI, no [PR](https://github.com/wailsapp/wails/pull/5734) de @taliesin-ai, com base na correção validada por @randalmurphal, com verificação da causa raiz por @eleclin e testes de hardware por @qq540491950

## v3.0.0-alpha2.113 - 2026-07-04

## Adicionado

- Adiciona um aviso ao compilar um AAB de lançamento sem `ANDROID_KEYSTORE_FILE` definido (o Google Play rejeita pacotes assinados para depuração) e documenta o empacotamento e a assinatura de App Bundles no [PR](https://github.com/wailsapp/wails/pull/5730) por @taliesin-ai
- Adiciona a documentação Why Wails em vários idiomas no [PR](https://github.com/wailsapp/wails/pull/5739) por @taliesin-ai
- Adiciona suporte ao mapeamento de Go time.Time para JS Date ou string nos bindings no [PR](https://github.com/wailsapp/wails/pull/5398) por @fbbdev
- Adiciona tarefas de empacotamento de Android App Bundle (AAB) (`bundle`, `bundle:fat`, `assemble:aab`, `assemble:aab:release`) para envio à Play Store — as tarefas de APK permanecem disponíveis para testes locais ou em emuladores — no [PR](https://github.com/wailsapp/wails/pull/5728) por @mortenolsrud (corrige [#5726](https://github.com/wailsapp/wails/issues/5726))
- Adiciona destinos de tarefas para dispositivos Android físicos e retoma as permissões de câmera/localização no [PR](https://github.com/wailsapp/wails/pull/5735) por @taliesin-ai

## Alterado

- Atualiza `webview2` para v1.0.28 ([notas de lançamento](https://github.com/wailsapp/wails/releases/tag/webview2%2Fv1.0.28)).
- Atualiza `compileSdk`/`targetSdk` do modelo Android de 34 para 35, conforme exigido pelo Google Play para o envio de novos aplicativos, no [PR](https://github.com/wailsapp/wails/pull/5730) por @taliesin-ai

## Corrigido

- Corrige máscaras de avatar pré-processadas no sponsorkit no [PR](https://github.com/wailsapp/wails/pull/5745) por @leaanthony
- Corrige a seleção da imagem de sistema ou da versão de cmdline-tools incorreta durante a criação automática de AVDs Android, causada pela ordenação lexicográfica das versões, no [PR](https://github.com/wailsapp/wails/pull/5730) por @taliesin-ai
- Corrige a sugestão de uma versão obsoleta do Android NDK pelo assistente de configuração (agora 26.3.11579264, em conformidade com o requisito documentado) no [PR](https://github.com/wailsapp/wails/pull/5730) por @taliesin-ai
- Atualiza a documentação em francês sobre SvelteKit e opções no [PR](https://github.com/wailsapp/wails/pull/5744) por @leaanthony
- Corrige um SIGSEGV durante a enumeração de telas no macOS quando há alterações nos monitores, no [PR](https://github.com/wailsapp/wails/pull/5516) por @flofreud

## v3.0.0-alpha2.112 - 2026-07-03

## Adicionado

- Adiciona um gerador de SVG de colaboradores baseado em Go e atualiza as páginas de créditos da documentação e do site no [PR](https://github.com/wailsapp/wails/pull/5724) por @taliesin-ai
- Adiciona suporte ao mapeamento de Go time.Time para JS Date ou string nos bindings no [PR](https://github.com/wailsapp/wails/pull/5398) por @fbbdev

## Alterado

- Substitui o pipeline de imagens de patrocinadores baseado em Node por um gerador em Go no [PR](https://github.com/wailsapp/wails/pull/5719) por @taliesin-ai

## Corrigido

- Corrige o script de instalação das dependências de assets da compilação para Android no [PR](https://github.com/wailsapp/wails/pull/5729) por @taliesin-ai
- Rejeita o caractere de controle U+0085 (NEXT LINE) em `ValidateAndSanitizeURL`, completando a cobertura de espaços em branco do validador de URLs
- Recalcula a moldura do DWM quando o DPI muda em janelas sem moldura no [PR](https://github.com/wailsapp/wails/pull/4785) por @leaanthony
- Corrige a falha na detecção de zonas de soltura de DnD no Windows quando a escala não é de 100%, no [PR](https://github.com/wailsapp/wails/pull/4632) por @yulesxoxo
- Adiciona gerenciamento explícito de memória Objective-C para objetos Cocoa em diálogos, menus, bandeja do sistema e notificações no Darwin, no [PR](https://github.com/wailsapp/wails/pull/5714) por @taliesin-ai
- Corrige bugs no backend CGO para Linux e problemas na bandeja do sistema no [PR](https://github.com/wailsapp/wails/pull/5718) por @taliesin-ai

## v3.0.0-alpha2.111 - 2026-07-01

## Adicionado

- Adiciona HappyTools à vitrine da comunidade no [PR](https://github.com/wailsapp/wails/pull/5061) por @Aliuyanfeng
- Adiciona suporte à localidade indonésia e documentação abrangente no [PR](https://github.com/wailsapp/wails/pull/5643) por @triadmoko
- Adiciona a opção DisableMenu a WindowsWindow no [PR](https://github.com/wailsapp/wails/pull/4813) por @leaanthony

## Alterado

- Atualiza o modelo de Taskfile e a CLI para encaminhar tarefas de compilação/empacotamento com GOOS e ARCH no [PR](https://github.com/wailsapp/wails/pull/5617) por @leaanthony

## Corrigido

- Corrige um problema com o agrupamento de janelas do macOS em abas no [PR](https://github.com/wailsapp/wails/pull/5708) por @taliesin-ai

## Removido

- Remove os arquivos MDX traduzidos para alemão das seções sobre contribuição, recursos e guias no [PR](https://github.com/wailsapp/wails/pull/5702) por @taliesin-ai

## v3.0.0-alpha2.110 - 2026-06-30

## Adicionado

- Implementa o recarregamento normal e forçado da WebView no macOS e adiciona recuperação após o encerramento do processo WebContent no [PR](https://github.com/wailsapp/wails/pull/5129) por @wayneforrest
- Adiciona documentação abrangente em alemão sobre contribuição, recursos e guias no [PR](https://github.com/wailsapp/wails/pull/5396) por @leaanthony
- Aprimora as notificações com som, anexos, agendamento e API de atualização no [PR](https://github.com/wailsapp/wails/pull/5333) por @popaprozac

## Corrigido

- Recalcula a moldura do DWM quando o DPI muda em janelas sem moldura no [PR](https://github.com/wailsapp/wails/pull/4785) por @leaanthony
- Corrige a falha na detecção de zonas de soltura de DnD no Windows quando a escala não é de 100%, no [PR](https://github.com/wailsapp/wails/pull/4632) por @yulesxoxo

## v3.0.0-alpha2.109 - 2026-06-29

## Adicionado

- Adiciona exemplos de código à documentação de EventsEmit no [PR](https://github.com/wailsapp/wails/pull/5026) por @iamhabbeboy
- Adiciona a opção de hospedagem visual do WebView2 no Windows no [PR](https://github.com/wailsapp/wails/pull/5380) por @MerIijn
- Adiciona Klustr à documentação da vitrine da comunidade no [PR](https://github.com/wailsapp/wails/pull/5536) por @SametKUM
- Adiciona Kira à vitrine da comunidade, com novas páginas e uma entrada no changelog, no [PR](https://github.com/wailsapp/wails/pull/5685) por @thiennguyen93
- Adiciona uma seção de feedback ao guia do serviço MCP no [PR](https://github.com/wailsapp/wails/pull/5694) por @taliesin-ai

## Alterado

- O modo servidor agora tem uma compilação de produção de primeira classe, alinhada às tarefas de compilação para desktop (#5693). `task build:server` compila um binário de produção por padrão (`-tags server,production`, `-trimpath`, sem símbolos) e aceita `DEV=true` (servidor de desenvolvimento), `OBFUSCATED=true` (garble) e `EXTRA_TAGS`. `task run:server` executa um servidor de desenvolvimento. `Dockerfile.server` / `task build:docker` compilam primeiro o servidor de produção (`-tags server,production`) e o frontend de produção; por padrão, a imagem usa uma compilação estática em Go puro sobre distroless/static, com `CGO_ENABLED`, `GO_IMAGE` e `RUNTIME_IMAGE` expostos como argumentos de compilação substituíveis para aplicativos CGO.

## Corrigido

- Evita uma falha ao fechar uma janela com chamadas assíncronas pendentes no [PR](https://github.com/wailsapp/wails/pull/4435) de @leaanthony
- Evita a ativação da janela ao abrir aplicativos ocultos no Windows no [PR](https://github.com/wailsapp/wails/pull/5249) de @leaanthony
- Garante que os metadados das solicitações do WebKit, a conclusão das respostas e o processamento do fluxo do corpo sejam executados na thread principal do GTK no [PR](https://github.com/wailsapp/wails/pull/5668) de @taliesin-ai
- Corrige o problema que fazia `Menu.Update()` não reconstruir o menu nativo no Linux com GTK4 (#5659, diagnosticado e corrigido de forma independente por @puneetdixit200 em #5539)
- Corrige uma falha ao enumerar as telas do macOS após uma alteração de monitor, copiando as strings de ID/nome das telas e criando um instantâneo da contagem (#5565, diagnosticado e corrigido de forma independente por @x-haose em #5584)
- Corrige o conteúdo do WebView2 que encolhia e depois desaparecia após arrastar uma janela entre monitores com DPIs diferentes no Windows, reafirmando os limites do controlador no manipulador `WM_DPICHANGED`, de modo equivalente à ressincronização de DPI ao restaurar uma janela minimizada (#5677)

## v3.0.0-alpha2.108 - 2026-06-28

## Adicionado

- Adiciona atalhos de teclado globais (em todo o sistema) por meio de `app.GlobalShortcut` (`Register`, `Unregister`, `UnregisterAll`, `IsRegistered`, `GetAll`). Os atalhos são acionados mesmo quando o aplicativo não está em foco. Implementação nativa para cada plataforma, sem dependências de terceiros: teclas de atalho do Carbon no macOS, `RegisterHotKey` no Windows, `XGrabKey` no X11 e a interface de atalhos globais do XDG Desktop Portal no Wayland.
- Adiciona um servidor MCP integrado: um servidor Model Context Protocol iniciado automaticamente quando o aplicativo é compilado com a tag `mcp`, permitindo que agentes de LLM testem e controlem um aplicativo Wails em execução — controle de janelas, inspeção do DOM, avaliação de JavaScript, chamadas de métodos vinculados, eventos e entrada simulada de mouse/teclado, exibida com um cursor animado na tela. Não exige código do usuário: a tag `mcp` é adicionada automaticamente por `wails3 build`/`wails3 dev` quando `WAILS_MCP=1` está definido. A configuração é feita inteiramente por variáveis de ambiente (`WAILS_MCP_HOST`, `WAILS_MCP_PORT`, `WAILS_MCP_TIMEOUT`, `WAILS_MCP_HIDE_CURSOR`).

## Corrigido

- Corrige o problema que fazia `Menu.Update()` não reconstruir o menu nativo no Linux com GTK4 (#5659, diagnosticado e corrigido de forma independente por @puneetdixit200 em #5539)
- Corrige uma falha ao enumerar as telas do macOS após uma alteração de monitor, copiando as strings de ID/nome das telas e criando um instantâneo da contagem (#5565, diagnosticado e corrigido de forma independente por @x-haose em #5584)

## v3.0.0-alpha2.107 - 2026-06-27

## Adicionado

- Adiciona documentação experimental do Wake com navegação pela barra lateral no [PR](https://github.com/wailsapp/wails/pull/5613) de @leaanthony

## v3.0.0-alpha2.106 - 2026-06-24

## Alterado

- Atualiza `webview2` para v1.0.27.
  - ci(webview2): corrige a compilação de lançamento (compilação cruzada para Windows + go.sum completo) (#5671)\

  **Diferenças completas:** https://github.com/wailsapp/wails/compare/webview2/v1.0.26...webview2/v1.0.27

- Remove go vet da compilação cruzada no fluxo de trabalho de lançamento do webview2 no [PR](https://github.com/wailsapp/wails/pull/5672) de @taliesin-ai
- Atualiza o modelo OpenRouter do auto-changelog para google/gemini-2.5-flash-lite no [PR](https://github.com/wailsapp/wails/pull/5670) de @taliesin-ai
- Atualiza `webview2` para v1.0.26.

### Correções

- **Recupera-se de erros COM transitórios em tempo de execução em vez de encerrar** (#5658, #5580). Anteriormente, `Chromium.errorCallback` chamava `os.Exit(1)` para *qualquer* erro COM; portanto, uma falha temporária recuperável após a inicialização encerrava todo o aplicativo. Agora, os caminhos de tempo de execução (`Resize`/`GetClientRect`, `Navigate`/`NavigateToString`, `Init`, `MessageReceived`, `PutZoomFactor`, `OpenDevToolsWindow`) registram o erro e se recuperam. Em particular, uma mensagem web malformada/não confiável em `MessageReceived` agora é descartada em vez de encerrar o processo. Isso soluciona a classe de falhas ao atravessar monitores com DPIs diferentes (#5544, #5650). Os caminhos de criação do ambiente/controlador continuam sendo fatais.\

**Diferenças completas:** https://github.com/wailsapp/wails/compare/webview2/v1.0.25...webview2/v1.0.26

## Corrigido

- Corrige o fluxo de trabalho release-webview2 para processar corretamente os arquivos go.sum no [PR](https://github.com/wailsapp/wails/pull/5671) de @taliesin-ai
- Corrige as atualizações de menu no Linux com GTK4, limpando e reconstruindo o menu nativo no [PR](https://github.com/wailsapp/wails/pull/5659) de @taliesin-ai

## v3.0.0-alpha2.105 - 2026-06-21

## Adicionado

- Adiciona `application.System` para detectar a plataforma em tempo de execução a partir de código compartilhado: `System.IsMobile()` (iOS/Android), `System.IsDesktop()` (macOS/Windows/Linux), `System.IsServer()` (a tag de compilação `server`) e `System.IsPlatform(application.PlatformMacOS|PlatformWindows|PlatformLinux|PlatformIOS|PlatformAndroid|PlatformServer)` para testar diretamente um único destino. Ele é compilado para todos os destinos, permitindo criar ramificações sem tags de compilação. Os auxiliares correspondentes para o frontend (`System.IsMobile/IsDesktop/IsIOS/IsAndroid/...`) estão disponíveis em `@wailsio/runtime`
- Adiciona um guia "Como usar outros frameworks de frontend", que mostra como inserir seu próprio projeto Vite em `frontend/` (abrange Solid, Preact, Lit, SvelteKit, Qwik, Angular etc.)
- O assistente `wails3 setup` agora verifica a cadeia de ferramentas para dispositivos móveis (iOS/Android) — Xcode e o runtime do Simulador do iOS, JDK, Android SDK/NDK e emulador — com instalação em um clique e correções copiáveis para a configuração do shell, quando aplicável
- Os projetos gerados incluem uma configuração `frontend/.npmrc` que define um `minimum-release-age` de 7 dias para reduzir a exposição a pacotes recém-publicados (e potencialmente comprometidos); a configuração é respeitada pelo pnpm e pelo bun e ignorada sem efeitos adversos pelo npm

## Alterado

- Redesenha todos os modelos iniciais integrados com um novo visual principal de montanhas neon (web, iOS e Android)
- **O TypeScript agora é o padrão dos modelos iniciais e assume o nome do modelo sem qualificador.** `wails3 init` (sem `-t`) cria a estrutura de um projeto TypeScript; `-t vanilla`, `-t react`, `-t vue` e `-t svelte` usam TypeScript, com variantes JavaScript em `-t vanilla-js`, `-t react-js`, `-t vue-js` e `-t svelte-js`. Os modelos integrados declaram sua linguagem com `typescript:` em `template.yaml`; como alternativa, os modelos da comunidade que usam o sufixo `-ts` continuam funcionando
- Redesenha o assistente `wails3 setup` com o tema neon "Wails digital" (vivacidade de vidro fosco sobre um cenário montanhoso)

## Corrigido

- Corrige uma falha no Windows ao restaurar um aplicativo que ficou minimizado por tempo suficiente para o WebView2 ser suspenso ou para seu processo de renderização/GPU ser reciclado. A resincronização de DPI ao minimizar/restaurar (#5544) agora só altera o controlador do WebView2 quando o DPI da janela realmente mudou, evitando chamadas COM fatais para um controlador suspenso no caso comum de restauração com o mesmo DPI (#5605)
- Corrige falhas nativas repetidas de `SIGABRT`/`SIGSEGV` (normalmente dentro de `g_object_unref` durante o loop principal do GTK) em aplicativos Linux executados por longos períodos e sujeitos a carregamentos frequentes de ativos/mídia. O servidor de ativos concluía `WebKitURISchemeRequest`s a partir de goroutines de trabalho, chamando funções do WebKit2GTK que não são seguras para threads fora da thread principal do GTK; a conclusão (`webkit_uri_scheme_request_finish_with_response`/`finish_error`) agora é executada na thread principal. Conclui a correção parcial de #5566. Afeta tanto as compilações GTK3 quanto GTK4/WebKitGTK 6.0 (#5631, #5557)
- Corrige uma ocorrência intermitente de `fatal error: invalid pointer found on stack` em `setupSignalHandlers` no Linux/GTK3. Os IDs de janela passados como `user_data` de sinal eram armazenados em uma variável local Go do tipo `unsafe.Pointer`; por isso, o coletor de lixo abortava ao examinar o valor (que não era um ponteiro) durante uma cópia da pilha. Agora, o ID é mantido como inteiro (`uintptr_t`) no lado Go, trazendo para o caminho GTK3 legado a mesma correção que #4958 aplicou ao caminho GTK4 (que alterou as funções de sinal C para `uintptr_t` a fim de eliminar erros de `-race`/checkptr) (#5631)

## Removido

- Remove os modelos iniciais `react-swc`, `preact`, `lit`, `solid`, `qwik` e `sveltekit` (e suas variantes `-ts`). O conjunto integrado com suporte agora é composto por `vanilla`, `react`, `vue` e `svelte` — todos usam TypeScript por padrão e têm variantes JavaScript `-js`. Ainda é possível usar qualquer outro framework [fornecendo seu próprio frontend](https://v3.wails.io/guides/dev/frontend-frameworks) ou por meio de um modelo personalizado

## v3.0.0-alpha2.104 - 2026-06-18

## Corrigido

- Corrige uma falha (SIGABRT) no iOS quando um método vinculado de um serviço Go retorna uma string vazia. O gravador de respostas de ativos do iOS verificava o ponteiro do corpo com `buf != nil` em vez de verificar seu comprimento, portanto um corpo de comprimento zero fazia `&buf[0]` entrar em pânico; agora ele verifica o comprimento, como os gravadores para desktop

## v3.0.0-alpha2.103 - 2026-06-15

## Alterado

- Move os recursos nativos de iOS e Android para gerenciadores de plataforma: chame-os por meio de `application.IOS.*` e `application.Android.*` (por exemplo, `application.IOS.Haptic("medium")`, `application.Android.Share(payload)`), em vez das antigas funções independentes `application.IOS*`/`application.Android*` (#5602)
- Renomeia os eventos da ponte móvel: os eventos multiplataforma agora usam o prefixo `common:*` (por exemplo, `common:haptic`, `common:location`), enquanto eventos exclusivos de uma plataforma usam `ios:*` / `android:*` (por exemplo, `ios:backgroundTask`, `android:foregroundService`); o prefixo `native:*` não é mais usado (#5602)

## v3.0.0-alpha.102 - 2026-06-14

## Adicionado

- Adiciona o assistente experimental `wails3 setup` para configuração interativa de projetos e verificação de dependências
- Adiciona o sinalizador `--json` a `wails3 doctor` para gerar uma saída legível por máquina
- Adiciona uma seção de status da assinatura ao comando `wails3 doctor`

## Corrigido

- Corrige a detecção do npm no Linux para verificar o PATH além do gerenciador de pacotes

## v3.0.0-alpha.101 - 2026-06-13

## Adicionado

- iOS: caixas de diálogo nativas de mensagens (UIAlertController) e caixas de diálogo para abrir arquivo, arquivos ou diretório (UIDocumentPickerViewController); caixas de diálogo para salvar retornam um erro explícito
- iOS: suporte à área de transferência por meio de UIPasteboard
- iOS: métricas reais da tela por meio de UIScreen (pontos, pixels, escala e área de trabalho da área segura)
- iOS: compilações para dispositivos (`IOS_PLATFORM=device`), suporte a identidade de assinatura de código/perfil de provisionamento/entitlements, empacotamento `.ipa` e `deploy-device` por meio de devicectl
- iOS: versão mínima configurável do iOS (`ios.minIOSVersion` em build/config.yml)
- iOS: `wails3 doctor` informa a disponibilidade do Xcode e do SDK do iOS no macOS
- iOS: eventos do sistema — bateria, rede, tema, bloqueio de tela e pouca memória são expostos como eventos de aplicativo `events.IOS.*` e como eventos de aplicativo independentes de plataforma `events.Common.*`
- iOS: ponte nativa de recursos móveis (`application.IOS*` exportado) — folha de compartilhamento, abertura de URL, manutenção da tela ativa, lanterna, margens da área segura, brilho, informações do aplicativo, bloqueio de orientação, barra de status, biometria (Face ID/Touch ID), notificações locais e armazenamento seguro no Keychain
- iOS: sensores e hardware — resposta tátil, consulta pontual de geolocalização, acelerômetro, proximidade, conversão de texto em fala, informações de armazenamento, estado de energia/bateria, status da rede, margens do teclado e detecção de captura de tela
- iOS: documentação (IOS.md e um guia no site da documentação)
- Android: caixas de diálogo nativas de mensagens (AlertDialog) e caixas de diálogo para abrir arquivo ou arquivos (Storage Access Framework, importados como cópias no cache); caixas de diálogo para abrir diretório e salvar retornam um erro explícito
- Android: suporte à área de transferência por meio de ClipboardManager
- Android: métricas reais da tela por meio de WindowMetrics/DisplayMetrics (dp, pixels, escala e área útil da tela considerando as barras do sistema)
- Android: métodos de runtime para resposta tátil (`Android.Haptics.Vibrate`), informações do dispositivo (`Android.Device.Info`) e mensagens toast (`Android.Toast.Show`)
- Android: eventos tipados do ciclo de vida (`events.Android.*`, gerados a partir de events.txt), com `ActivityCreated` mapeado para `Common.ApplicationStarted`
- Android: o pipeline de compilação produz APKs instaláveis de depuração e lançamento (`android:run`, `android:package`, `android:package:fat`); por padrão, a assinatura da versão de lançamento usa o keystore de depuração, ou um keystore real por meio de variáveis de ambiente `ANDROID_KEYSTORE_*`
- Android: `wails3 doctor` informa o SDK e o NDK do Android, além do JDK
- Android: eventos do sistema — bateria, rede, tema, bloqueio de tela e pouca memória são expostos como eventos de aplicativo `events.Android.*` e `events.Common.*` independentes de plataforma
- Android: ponte nativa de recursos móveis (`application.Android*` exportado) — compartilhamento, abertura de URL, manutenção da tela ativa, lanterna, margens da área segura, brilho, informações do aplicativo, bloqueio de orientação, barra de status, biometria (BiometricPrompt), notificações locais e armazenamento seguro com EncryptedSharedPreferences
- Android: sensores e hardware — resposta tátil, consulta pontual de geolocalização, acelerômetro, proximidade, conversão de texto em fala, informações de armazenamento, estado de energia/bateria, status da rede, margens do teclado e bloqueio de captura de tela com FLAG_SECURE
- Android: documentação (ANDROID.md e um guia no site da documentação)
- Exemplo: o mostruário `mobile` recebe as abas Dispositivos móveis e Hardware, que demonstram a ponte de recursos nativos no iOS e Android (as abas em formato de pílula quebram em várias linhas)
- Dispositivos móveis: bateria — o acelerômetro, o sensor de proximidade, a lanterna e o relógio periódico do exemplo são pausados quando o aplicativo fica em segundo plano e restaurados quando ele retorna (o Android mantém o processo em execução em segundo plano, e a lanterna é um estado do hardware que persiste no iOS); além disso, os receptores de eventos do sistema Android são registrados somente enquanto o aplicativo está em primeiro plano
- iOS: captura pela câmera — `application.IOSCapturePhoto`/`IOSCaptureVideo` (UIImagePickerController → um evento `native:capture` com uma miniatura em base64)
- iOS: execução em segundo plano — `application.IOSBeginBackgroundTask`/`IOSEndBackgroundTask` (um período de execução de tarefa em segundo plano concedido pela UIApplication) e um `ios.backgroundModes` configurável (build/config.yml), que preenche `UIBackgroundModes` no Info.plist gerado por meio de um template
- Android: captura pela câmera — `application.AndroidCapturePhoto`/`AndroidCaptureVideo` (câmera do sistema via FileProvider → um evento `native:capture`)
- Android: serviço em primeiro plano — `application.AndroidStartForegroundService`/`AndroidStopForegroundService` (um `WailsForegroundService` com uma notificação contínua mantém o processo ativo para trabalhos de longa duração em segundo plano)
- Exemplo: uma guia Câmera que demonstra a captura de fotos e vídeos e a execução em segundo plano (serviço em primeiro plano no Android e período permitido para execução de tarefa em segundo plano no iOS)

## Corrigido

- Corrige `getUserMedia`, que sempre falhava com `NotAllowedError` no Linux: o WebKitGTK nega solicitações de permissão que não são tratadas, e o sinal `permission-request` não estava conectado. Câmera e microfone agora são tratados conforme um novo mapa multiplataforma `WebviewWindowOptions.Permissions` (`map[PermissionType]Permission`), respeitado no Linux (WebKitGTK) e no Windows (WebView2). No Linux, que não tem uma solicitação nativa, câmera e microfone são permitidos por padrão (restaurando `getUserMedia`) e podem ser desativados com `PermissionDeny` (#5552)
- iOS: `GOOS=ios` volta a compilar (`events.IOS` exportado e stubs de nomes de métodos para dispositivos móveis), assim como as compilações com a tag de produção (correções de tags de compilação em pkg/application e em vários serviços)
- iOS: eventos Go→JS e ExecJS agora funcionam — a página não é mais carregada duas vezes na inicialização, e o handshake `wails:runtime:ready` não pode mais ser perdido
- iOS: `ApplicationDidFinishLaunching`/`ApplicationStarted` não entram mais em condição de corrida com a inicialização do aplicativo; removida a espera fixa de 2 segundos na inicialização
- iOS: corrigido um vazamento de string C a cada execução de JavaScript de Go→JS
- iOS: `hasListeners` agora reflete o registro real de listeners
- iOS: o registro de depuração do framework é excluído das compilações de produção
- Android: `GOOS=android` volta a compilar — `events.Android` foi definido, o array de listeners `events_android.go` com acesso fora dos limites foi removido, o stub de nomes de métodos para dispositivos móveis foi adicionado e arquivos do Linux para desktop (`linux_cgo.*`, `events_linux.*`, `environment_linux.go`) deixaram de ser incluídos nas compilações para Android
- Android: os bindings JS→Go agora funcionam — a WebView não consegue entregar corpos POST de `fetch()` a `shouldInterceptRequest`, portanto as chamadas do runtime são encaminhadas por um transporte JavascriptInterface (`nativeHandleRuntimeCall`), em vez de causarem uma falha devido a um corpo de solicitação nil
- Android: as chamadas de runtime `Screens.*` retornam dados reais — o ScreenManager agora é preenchido na inicialização (ele nunca havia sido conectado, portanto `GetAll` retornava nil)
- Android: o registro de depuração do framework é excluído das compilações de produção e, nas compilações de depuração, é encaminhado pelo logcat com a tag `Wails`
- Android: registro `hasListeners` real, tratamento de referências e exceções JNI e ciclo de vida da página com um único carregamento (sem navegação duplicada)
- Corrige a falha de `wails3 generate bindings` com "Acesso negado" no Windows quando o servidor de desenvolvimento do Vite está em execução, sincronizando os arquivos gerados com o diretório de saída em vez de renomeá-los sobre ele (#5515)
- Corrige uma falha fatal intermitente no macOS ao ler informações da tela após uma alteração de monitor: o ID e o nome da tela armazenavam ponteiros para buffers `UTF8String` com autorelease, que podiam ser liberados antes de serem copiados pelo Go (uso após liberação). Agora, aplica-se `strdup` às strings, que são liberadas após a conversão, e a enumeração de telas é executada em um pool de autorelease explícito para não causar mais vazamentos quando chamada por goroutines do Go (#5556)
- Corrige um SIGSEGV intermitente no Linux quando o assetserver fecha um `WebKitURISchemeRequest`: o `g_object_unref` final era executado na goroutine do assetserver, finalizando um GObject do WebKit fora da thread principal do GTK. Agora, o unref é encaminhado ao contexto principal do GTK por meio de `g_main_context_invoke` (#5557)

## v3.0.0-alpha.100 - 2026-06-13

## Adicionado

- Estende `MacWebviewPreferences` com opções adicionais de configuração da WKWebView: `EnableAutoplayWithoutUserAction`, `AllowsAirPlayForMediaPlayback`, `AllowsMagnification`, `JavaScriptCanOpenWindowsAutomatically`, `MinimumFontSize` e `ApplicationNameForUserAgent` (#5549)

## Corrigido

- Corrige a falha de `wails3 generate bindings` com "Acesso negado" no Windows quando o servidor de desenvolvimento do Vite está em execução, sincronizando os arquivos gerados com o diretório de saída em vez de renomeá-los sobre ele (#5561)
- Corrige eventos JS de redimensionamento que não eram disparados em janelas sem moldura no Linux; corrige a detecção das bordas da barra de rolagem em janelas sem moldura (#5368)
- Corrige a falha do atualizador no Windows com "link inválido entre dispositivos" quando o diretório temporário está em um volume diferente do diretório de instalação (#5560)

## v3.0.0-alpha.99 - 2026-06-10

## Corrigido

- Corrige a falha de `wails3 generate bindings` com "Acesso negado" no Windows quando o servidor de desenvolvimento do Vite está em execução, sincronizando os arquivos gerados com o diretório de saída em vez de renomeá-los sobre ele (#5515)

## v3.0.0-alpha.98 - 2026-06-03

## Corrigido

- Corrige o congelamento da interface do WebKit no Linux quando ociosa (por exemplo, com o inspetor aberto), deixando de forçar `SA_ONSTACK` em `SIGUSR1`, o que interrompia a sincronização da thread do coletor de lixo do JavaScriptCore (#5527)

## v3.0.0-alpha.97 - 2026-05-31

## Adicionado

- Adiciona uma página de depuração e instruções para trabalhar com `runtime/trace`

## Alterado

- Remove algumas importações desnecessárias de `_ "embed"`, deixando o código um pouco mais organizado

## Corrigido

- Corrige as restrições de largura e altura mínimas que não eram aplicadas após restaurar uma janela maximizada no Windows (#4593)
- Corrige a passagem de cliques do mouse no modo de tela cheia com as opções de janela Frameless + Transparent (#4408)

## v3.0.0-alpha.96 - 2026-05-25

## Adicionado

- Adiciona suporte à ofuscação com Garble ([#4563](https://github.com/wailsapp/wails/issues/4563)): IDs estáveis de métodos de binding, integração com o sistema de compilação/Taskfile (`build --obfuscated --garbleargs`, `generate bindings -obfuscated`) e tags de struct JSON em todos os payloads voltados ao runtime (`EnvironmentInfo`, `OSInfo`, `Screen`, `Rect`, `Point`, `Size`, `Capabilities`), para que o formato de transmissão sobreviva à renomeação de campos exportados feita pelo Garble.

## v3.0.0-alpha.95 - 2026-05-20

## Adicionado

- Adicionada a página ausente sobre a estrutura do projeto

## Alterado

- Documentação: altera alguns diagramas da página de arquitetura para usar diagramas de sequência e proporcionar uma exibição mais clara
- Documentação: incluir uma observação sobre a instalação do D2 como pré-requisito para a execução

## Corrigido

- Corrige `wails3 generate appimage` no padrão GTK4: o empacotador agora detecta a pilha GTK no binário antes de procurar os arquivos de runtime, selecionando `libwebkitgtkinjectedbundle.so` (em `webkitgtk-6.0/`) para builds GTK4 e `libwebkit2gtkinjectedbundle.so` (em `webkit2gtk-4.1/`) para builds `-tags gtk3`. A verificação de `.relr.dyn` também consulta `libgtk-4.so.1`, portanto a remoção de símbolos é desativada corretamente em toolchains modernas, independentemente da pilha. (#5475)
- Corrige a falha de `wails3 generate appimage` quando invocado com um `-builddir` relativo: o empacotador agora resolve antecipadamente `-binary`, `-icon`, `-desktopfile`, `-builddir` e `-outputdir` como caminhos absolutos, para que o `s.CD` executado no meio do fluxo não interrompa a goroutine de download do AppRun nem a verificação de `ldd` após a cópia.
- Corrige a falha de `wails3 generate appimage` ao mover a AppImage final para `-outputdir` quando o campo `Name=` do arquivo desktop não corresponde ao nome-base do binário: o empacotador agora força o plugin appimage do linuxdeploy (por meio da variável de ambiente `OUTPUT`) a gravar a AppImage em `<binary>-<arch>.AppImage`, em vez de usar o nome derivado do arquivo desktop.
- Corrige `events.Common.ApplicationStarted`, `Common.ThemeChanged`, `Common.SystemWillSleep` e `Common.SystemDidWake`, que não eram disparados no Linux após a pilha GTK4 + WebKitGTK 6.0 se tornar o padrão na alpha.93. O novo `application_linux.go` `run()` padrão não chamava `setupCommonEvents()` (que encaminha eventos `Linux.*` para seus equivalentes `Common.*`) nem `monitorPowerEvents()`. O auxiliar de monitoramento de energia via DBus agora é compartilhado entre os caminhos de build GTK3 e GTK4 por meio de `application_linux_dbus.go`. (#5474)

## v3.0.0-alpha.94 - 2026-05-19

## Corrigido

- Corrige `events.Common.ApplicationStarted`, `Common.ThemeChanged`, `Common.SystemWillSleep` e `Common.SystemDidWake`, que não eram disparados no Linux após a pilha GTK4 + WebKitGTK 6.0 se tornar o padrão na alpha.93. O novo `application_linux.go` `run()` padrão não chamava `setupCommonEvents()` (que encaminha eventos `Linux.*` para seus equivalentes `Common.*`) nem `monitorPowerEvents()`. O auxiliar de monitoramento de energia via DBus agora é compartilhado entre os caminhos de build GTK3 e GTK4 por meio de `application_linux_dbus.go`. (#5474)

## v3.0.0-alpha.93 - 2026-05-17

## Adicionado

- Adiciona `XDG_SESSION_TYPE` à saída de `wails3 doctor` no Linux, por @leaanthony

## Corrigido

- Corrige a falha do menu da janela no Wayland causada pelo acesso do appmenu-gtk-module a uma janela ainda não realizada (#4769), por @leaanthony
- Corrige a falha do aplicativo GTK quando o nome do aplicativo contém caracteres inválidos (espaços, parênteses etc.), por @leaanthony
- Corrige o erro "memória insuficiente" ao inicializar a funcionalidade de arrastar e soltar no Windows (#4701), por @overlordtm
- Corrige uma condição de corrida no armazenamento de callbacks da thread principal, que usava incorretamente RLock para excluir entradas do mapa (Linux, macOS e iOS) (#4424), por @leaanthony
- Corrige o tratamento de variáveis ao passar argumentos de linha de comando para tarefas. As variáveis da CLI especificadas como pares KEY=VALUE agora são inicializadas e propagadas corretamente durante toda a execução da tarefa.
- Corrige o conflito de NSWindowZoomButton no macOS: `MaximiseButtonState` e `FullscreenButtonState` agora aplicam o estado mais restritivo tanto na inicialização quanto em tempo de execução; nenhum dos setters pode mais sobrescrever silenciosamente o outro (#5319)
- Corrige um conjunto de bugs preexistentes no caminho de build GTK3 legado (`-tags gtk3`), revelados pelo CodeRabbit em #5463: inicializações por associação de arquivo não ignoram mais os manipuladores de inicialização; `getTheme` agora verifica limites e tipos com segurança; `appName` não libera mais memória pertencente ao GLib; `clipboardGet` não causa mais vazamento do `gchar*` retornado pelo GTK; `Calloc` agora usa receptores de ponteiro (e `NewCalloc` retorna `*Calloc`), permitindo que o pool realmente rastreie as alocações; `zoomOut` usa o inverso de `zoomInFactor`, em vez de um multiplicador negativo que limitava o valor a 1.0; `execJS` reutiliza o nome de mundo vazio pré-alocado, em vez de causar o vazamento de um `C.CString("")` a cada chamada; um `fmt.Println` de desenvolvimento foi removido de `menuItem.setAccelerator`. Resolve #5465.
- Corrige o mesmo vazamento causado pelo receptor por valor de `Calloc` no caminho de build GTK4 padrão (`linux_cgo.go`): receptores de ponteiro + `NewCalloc() *Calloc`, para que as alocações de `c.String(...)` por janela sejam realmente rastreadas e liberadas.

## v3.0.0-alpha.92 - 2026-05-15

## Adicionado

- Modifica os Taskfiles para permitir o controle do gerenciador de pacotes do frontend usado por meio da opção `PACKAGE_MANAGER`
- Enriquece os dados do template com `{{.Opn}}` e `{{.Cls}}`, tornando mais previsível a criação de templates de Taskfile

## Alterado

- Modifica alguns dos Taskfiles existentes para usar `{{.Opn}} and {{.Cls}}`

## Corrigido

- Corrige um erro fatal em tempo de execução de `concurrent map read and map write` em `linuxSystemTray` quando o menu da bandeja é atualizado enquanto o painel o lê.
- Usa `log` em vez de `fmt` para a saída de erros e rastreamentos de pilha do WebView2, evitando a perda de mensagens quando o aplicativo é executado sem um console conectado no Windows.

## v3.0.0-alpha.91 - 2026-05-12

## Alterado

- Atualiza o SVG dos patrocinadores no [PR](https://github.com/wailsapp/wails/pull/5414) por `@github-actions[bot]`
- **INCOMPATÍVEL (macOS):** normaliza o sistema de coordenadas do macOS para que `GetScreens`, `Position` e `SetPosition` usem o mesmo espaço — pontos lógicos, eixo Y crescente para baixo, com `(0,0)` no canto superior esquerdo da tela principal. Isso corresponde ao Windows, ao GTK e às APIs públicas do Electron e da web. Telas fisicamente acima da principal agora informam `Bounds.Y` negativo (antes era positivo), e os valores de `Position()`/`SetPosition()` agora são expressos em pontos lógicos, em vez de `points × primaryScale`. A conversão de ida e volta `Position()` → `SetPosition()` é preservada; valores absolutos registrados por builds alpha anteriores ou soluções alternativas calculadas manualmente (por exemplo, multiplicar por `primaryScale` ou inverter Y em relação à altura de uma tela) precisarão ser atualizados. Resolve [#5117](https://github.com/wailsapp/wails/issues/5117).

## Corrigido

- Valida defensivamente o nome do sinal DBus e o comprimento do corpo para evitar panics em [PR](https://github.com/wailsapp/wails/pull/5416), por @leaanthony
- Corrige um problema de segurança de memória no tratamento de menus GTK no Linux em [PR](https://github.com/wailsapp/wails/pull/5363), por @leaanthony
- Detecta GPUs NVIDIA e desativa o renderizador DMA-BUF no Linux em [PR](https://github.com/wailsapp/wails/pull/5295), por @leaanthony
- Corrige a conversão de Y entre telas de `SetPosition` no macOS: usa a altura da tela principal como referência global para posicionar as janelas corretamente em monitores com deslocamento vertical em relação à tela principal em [#5117](https://github.com/wailsapp/wails/issues/5117)
- Corrige o template de PR do Git para apontar para a URL de feedback correta em [PR](https://github.com/wailsapp/wails/pull/5109), por @wayneforrest
- Corrige uma série de falhas de `SetMenu` na bandeja do sistema do Windows, causadas por uma chamada de sistema `DestroyMenu` defeituosa que passava quatro argumentos em vez de um, fazendo com que todas as chamadas retornassem FALSE e não liberassem nada. Também libera identificadores HMENU e HBITMAP (inclusive os alocados em tempo de execução por meio de `MenuItem.SetBitmap`) durante as reconstruções do menu, redefine mapas obsoletos de caixas de seleção e botões de opção em `Win32Menu.Update` e remove uma chamada redundante a `Update()` em `systemtray.updateMenu` que duplicava as alocações. Aplicativos de longa duração que usam a bandeja do sistema não apresentam mais vazamentos de objetos GDI/USER a cada reconstrução do menu.

## v3.0.0-alpha.90 - 2026-05-11

## Adicionado

- Adiciona um nome de aplicativo configurável ao User-Agent do WKWebView no macOS no [PR](https://github.com/wailsapp/wails/pull/5261) por @vinhvoit225
- Adiciona a dependência indireta github.com/coder/websocket ao exemplo gin-service no [PR](https://github.com/wailsapp/wails/pull/5400) por @taliesin-ai
- Adiciona suporte à comparação de igualdade profunda aos testes de ativos de compilação no [PR](https://github.com/wailsapp/wails/pull/5402) por @leaanthony

## Alterado

- Consolida a saída da compilação no diretório de ativos no [PR](https://github.com/wailsapp/wails/pull/5401) por @taliesin-ai
- Atualiza o SVG dos patrocinadores no [PR](https://github.com/wailsapp/wails/pull/5399) por `@github-actions[bot]`

## Corrigido

- Usa um objeto de notificação para a mensagem de instância única no macOS no [PR](https://github.com/wailsapp/wails/pull/5289) por @overlordtm
- Agrupa os callbacks do Windows em lotes para evitar a perda de promises sob carga intensa no [PR](https://github.com/wailsapp/wails/pull/5383) por @taliesin-ai

## v3.0.0-alpha.89 - 2026-05-10

## Adicionado

- Adiciona o job go<em>test</em>results para agregar os resultados dos testes do Go no [PR](https://github.com/wailsapp/wails/pull/5316) por @leaanthony

## Alterado

- Divide condicionalmente cargas RPC grandes em requisições POST fragmentadas no [PR](https://github.com/wailsapp/wails/pull/5369) por @leaanthony
- Atualiza o Vite da versão 5.x.x para a 8.0.0 em todos os modelos de frontend no [PR](https://github.com/wailsapp/wails/pull/5386) por @leaanthony
- Migra a configuração da porta do servidor de desenvolvimento do Vite para variáveis de ambiente no [PR](https://github.com/wailsapp/wails/pull/5365) por @leaanthony
- Configura o servidor de desenvolvimento do Vite para escutar em 127.0.0.1 em todos os modelos no [PR](https://github.com/wailsapp/wails/pull/5361) por @leaanthony
- Atualiza o SVG dos patrocinadores no [PR](https://github.com/wailsapp/wails/pull/5384) por `@github-actions[bot]`

## Corrigido

- Sanitiza os stubs do modelo Info.plist durante a atualização de build-assets no [PR](https://github.com/wailsapp/wails/pull/5312) por @leaanthony
- Corrige o estado obsoleto nos menus do macOS aplicando os métodos modificadores dos itens de menu (`setMenuItemChecked()`, `setMenuItemLabel()`, `setMenuItemDisabled()`, `setMenuItemHidden()`, `setMenuItemTooltip()`) de forma síncrona na thread principal, eliminando a condição de corrida em `dispatch_async` que fazia os menus renderizarem o estado anterior quando eram reabertos rapidamente (#5002)
- Ignora arquivos `*_test.go` no modo de desenvolvimento para evitar recompilações desnecessárias no [PR](https://github.com/wailsapp/wails/pull/5203) por @leaanthony
- Evita uma falha de segmentação em Menu.Update() quando o aplicativo não está em execução no [PR](https://github.com/wailsapp/wails/pull/5291) por @wucm667
- Usa lastSizeWParam para controlar o redesenho da barra de menus no Windows no [PR](https://github.com/wailsapp/wails/pull/5382) por @taliesin-ai

## v3.0.0-alpha.88 - 2026-05-09

## Alterado

- Altera HiddenOnTaskbar para usar WS<em>EX</em>TOOLWINDOW no [PR](https://github.com/wailsapp/wails/pull/5371) por @leaanthony
- Reordena as dependências e remove a diretiva replace de webview2 no go.mod no [PR](https://github.com/wailsapp/wails/pull/5370) por @atterpac
- Atualiza o SVG dos patrocinadores no [PR](https://github.com/wailsapp/wails/pull/5358) por `@github-actions[bot]`

## Corrigido

- Remove aliases genéricos de indireção e consolida os tipos de chave de mapas no [PR](https://github.com/wailsapp/wails/pull/5331) por @fbbdev

## Removido

- Exclui o workflow PR-master, removendo a documentação, os testes do Go e a funcionalidade de ignorar testes no [PR](https://github.com/wailsapp/wails/pull/5377) por @leaanthony

## v3.0.0-alpha.87 - 2026-05-07

## Adicionado

- Adiciona documentação em coreano para o Wails v3 no [PR](https://github.com/wailsapp/wails/pull/5352) por @leaanthony
- Adiciona documentação em francês sobre instalação e início rápido no [PR](https://github.com/wailsapp/wails/pull/5354) por @leaanthony
- Adiciona documentação em português sobre início rápido, conceitos e comunidade no [PR](https://github.com/wailsapp/wails/pull/5355) por @leaanthony

## v3.0.0-alpha.86 - 2026-05-06

## Adicionado

- Adiciona a localização da documentação em francês no [PR](https://github.com/wailsapp/wails/pull/5328) por @leaanthony
- Adiciona a localidade alemã ao site de documentação no [PR](https://github.com/wailsapp/wails/pull/5343) por @leaanthony

## Alterado

- Registra todas as 8 localidades traduzidas na configuração da documentação no [PR](https://github.com/wailsapp/wails/pull/5347) por @leaanthony
- Atualiza diversos arquivos relacionados ao Windows para o WebView2 no [PR](https://github.com/wailsapp/wails/pull/5317) por @leaanthony

## Corrigido

- Separa o encaminhamento de diálogos entre GTK3 e GTK4 no Linux no [PR](https://github.com/wailsapp/wails/pull/5340) por @leaanthony
- Garante que os callbacks de diálogos sejam executados na thread do GTK, corrigindo falhas de segmentação no [PR](https://github.com/wailsapp/wails/pull/5339) por @leaanthony

## v3.0.0-alpha.85 - 2026-05-05

## Adicionado

- Adiciona ao repositório a URL do modelo de PR no [PR](https://github.com/wailsapp/wails/pull/5179) por @leaanthony
- Adiciona documentação em alemão para o Wails v3 no [PR](https://github.com/wailsapp/wails/pull/5330) por @leaanthony

## v3.0.0-alpha.84 - 2026-05-03

## Adicionado

- Adiciona uma opção para impedir que a tecla Escape encerre o modo de tela cheia no macOS no [PR](https://github.com/wailsapp/wails/pull/5307) por @leaanthony
- Adiciona uma opção para impedir que a tecla Escape encerre o modo de tela cheia no macOS no [PR](https://github.com/wailsapp/wails/pull/5310) por @leaanthony
- Adiciona a documentação da Pausa à vitrine da comunidade no [PR](https://github.com/wailsapp/wails/pull/5288) por @yuseferi

## Alterado

- Atualiza o SVG dos patrocinadores no [PR](https://github.com/wailsapp/wails/pull/5308) por `@github-actions[bot]`
- Atualiza o comando de geração de ícones para lidar com plataformas não compatíveis no [PR](https://github.com/wailsapp/wails/pull/5309) por @leaanthony
- Substitui a API booleana de tela cheia por ButtonState de três estados e implementa bindings específicos de plataforma no [PR](https://github.com/wailsapp/wails/pull/5224) por @leaanthony

## Corrigido

- Protege as operações de foco do WebView2 contra um estado de controlador nil no [PR](https://github.com/wailsapp/wails/pull/5315) por @leaanthony
- Atualiza o fluxo de trabalho do GitHub Actions para referenciar corretamente o branch base do PR no [PR](https://github.com/wailsapp/wails/pull/5313) por @leaanthony
- Ignora arquivos `*_test.go` no modo de desenvolvimento para evitar recompilações desnecessárias no [PR](https://github.com/wailsapp/wails/pull/5203) por @leaanthony
- Evita falha de segmentação em Menu.Update() quando o aplicativo não está em execução no [PR](https://github.com/wailsapp/wails/pull/5291) por @wucm667

## v3.0.0-alpha.83 - 2026-05-02

## Adicionado

- Adiciona a flag InstallScope e uma opção de compilação para instalação por máquina/usuário no [PR](https://github.com/wailsapp/wails/pull/5094) por @symball
- Adiciona o método SetScreen sem operação a BrowserWindow para satisfazer a interface Window no [PR](https://github.com/wailsapp/wails/pull/5294) por @leaanthony

## Corrigido

- Detecta GPUs NVIDIA e desativa o renderizador DMA-BUF no Linux no [PR](https://github.com/wailsapp/wails/pull/5295) por @leaanthony
- Corrige o modelo de PR do Git para apontar para a URL correta de feedback no [PR](https://github.com/wailsapp/wails/pull/5109) por @wayneforrest
- Corrige uma série de falhas `SetMenu` da bandeja do sistema no Windows causadas por uma chamada de sistema `DestroyMenu` defeituosa, que passava quatro argumentos em vez de um; por isso, todas as chamadas retornavam FALSE e não liberavam nada. Também libera os identificadores HMENU e HBITMAP (inclusive os alocados em tempo de execução por meio de `MenuItem.SetBitmap`) nas reconstruções de menu, redefine mapas obsoletos de caixas de seleção/botões de opção em `Win32Menu.Update` e remove uma chamada redundante de `Update()` em `systemtray.updateMenu` que duplicava as alocações. Aplicativos de bandeja do sistema executados por longos períodos não apresentam mais vazamento de objetos GDI/USER a cada reconstrução do menu.

## v3.0.0-alpha.82 - 2026-05-01

## Corrigido

- Corrige a geração do arquivo desktop para processar corretamente o nome do desktop no [PR](https://github.com/wailsapp/wails/pull/5232) por @leaanthony

## v3.0.0-alpha.81 - 2026-04-30

## Alterado

- Ajusta o agendamento de versões nightly para 15:00 UTC no [PR](https://github.com/wailsapp/wails/pull/5286) por @leaanthony

## Corrigido

- Corrige valores de Screen Bounds, WorkArea e Size reduzidos à metade em Macs com tela Retina —  (#5168)

## v3.0.0-alpha.80 - 2026-04-29

## Alterado

- Atualiza as dependências da documentação e os carregadores de coleções de conteúdo no [PR](https://github.com/wailsapp/wails/pull/5285) por @leaanthony

## v3.0.0-alpha.79 - 2026-04-29

## Adicionado

- Concede permissão actions: write ao trabalho trigger-release no [PR](https://github.com/wailsapp/wails/pull/5270) por @leaanthony

## Alterado

- A tarefa de lançamento passa a usar o branch master como padrão e atualiza a redação do changelog no [PR](https://github.com/wailsapp/wails/pull/5283) por @leaanthony
- Atualiza o fluxo de trabalho de geração automática do changelog para usar a versão mais recente no [PR](https://github.com/wailsapp/wails/pull/5282) por @leaanthony
- Melhora a eficiência dos fluxos de trabalho adicionando filtros de caminho e removendo fluxos obsoletos no [PR](https://github.com/wailsapp/wails/pull/5280) por @leaanthony
- Atualiza a documentação para referenciar o branch master nos links de exemplos no [PR](https://github.com/wailsapp/wails/pull/5274) por @leaanthony
- Atualiza a documentação e os exemplos para a v3 no [PR](https://github.com/wailsapp/wails/pull/5272) por @leaanthony

## Corrigido

- Aprimora o proxy reverso com lógica de repetição e uso forçado de IPv4 no desenvolvimento no [PR](https://github.com/wailsapp/wails/pull/5265) por @AkagiYui
- Reescreve o fluxo de trabalho de acionamento do changelog de itens ainda não lançados no [PR](https://github.com/wailsapp/wails/pull/5281) por @leaanthony

## Removido

- Remove scripts de teste em shell usados para várias finalidades de teste no [PR](https://github.com/wailsapp/wails/pull/5267) por @leaanthony
- Exclui o fluxo de trabalho de implantação da documentação v3-alpha e o registro CNAME no [PR](https://github.com/wailsapp/wails/pull/5266) por @leaanthony

### Adicionado

- Adiciona a entrada Roteamento do frontend à navegação da barra lateral no [PR](https://github.com/wailsapp/wails/pull/5196) por @leaanthony
- Adiciona um guia de roteamento do frontend com recomendações específicas para cada framework no [PR](https://github.com/wailsapp/wails/pull/5185) por @leaanthony
- Adiciona suporte a folhas modais (macOS)
- Atualiza a versão do ghw para oferecer melhor suporte a dispositivos Apple por @leaanthony (#4977)
- Adiciona o método `GetBadge` ao serviço de dock
- Adiciona a flag `-tags` ao comando `wails3 build` para passar tags de compilação personalizadas do Go (por exemplo, `wails3 build -tags gtk4`) (#4957)
- Adiciona documentação sobre a geração automática de enums no gerador de bindings, incluindo uma página dedicada a Enums e sua entrada na navegação da barra lateral (#4972)
- Adiciona a flag `-tags` ao comando `wails3 build` para passar tags de compilação personalizadas do Go (por exemplo, `wails3 build -tags gtk4`) (#4957)
- Adiciona exemplos de APIs Web em `v3/examples/web-apis/` que demonstram 41 APIs de navegador, incluindo Armazenamento (localStorage, sessionStorage, IndexedDB, Cache API), Rede (Fetch, WebSocket, XMLHttpRequest, EventSource, Beacon), Mídia (Canvas, WebGL, Web Audio, MediaDevices, MediaRecorder, Speech Synthesis), Dispositivo (Geolocation, Clipboard, Fullscreen, Device Orientation, Vibration, Gamepad), Desempenho (Performance API, Mutation Observer, Intersection/Resize Observer), Interface do usuário (Web Components, Pointer Events, Selection, Dialog, Drag and Drop) e muito mais
- Adiciona um exemplo de verificador de compatibilidade da API do WebView (`v3/examples/webview-api-check/`), que testa mais de 200 APIs de navegador entre plataformas
- Adiciona o pacote `internal/libpath` para localizar caminhos de bibliotecas nativas no Linux, com pesquisa paralela, cache e suporte a Flatpak/Snap/Nix
- **Em desenvolvimento:** adiciona suporte experimental a WebKitGTK 6.0 / GTK4 no Linux, disponível por meio de `-tags gtk4` (GTK3/WebKit2GTK 4.1 continua sendo o padrão)
- Observação: em gerenciadores de janelas lado a lado (por exemplo, Hyprland e Sway), as operações de minimizar/maximizar podem não funcionar como esperado, pois o gerenciador de janelas controla a geometria da janela
- Adicionada à documentação de **Escuta de eventos em JavaScript** uma explicação sobre como criar **Manipuladores de execução única**, por @AbdelhadiSeddar
- Adiciona a opção `UseApplicationMenu` a `WebviewWindowOptions`, permitindo que as janelas no Windows/Linux herdem o menu do aplicativo definido por meio de `app.Menu.Set()`, por @leaanthony
- Adiciona suporte ao uso de arquivos `.icon` (formato do Apple Icon Composer) para gerar ícones Liquid Glass e catálogos de ativos (macOS) (#4934), por @wimaha
- Adiciona um modo de servidor experimental para implantações sem interface gráfica/na Web (`-tags server`). Permite executar aplicativos Wails como servidores HTTP sem dependências de interface gráfica nativa. Compile com `wails3 task build:server`. Consulte `examples/server` para obter detalhes.
- Adiciona o pacote `internal/libpath` para localizar caminhos de bibliotecas nativas no Linux, com pesquisa paralela, cache e suporte a Flatpak/Snap/Nix
- Adiciona a opção `CollectionBehavior` a `MacWindow` para controlar o comportamento das janelas entre Spaces do macOS e o modo de tela cheia (#4756), por @leaanthony
- Adiciona testes unitários para pkg/application, por @leaanthony
- Adiciona suporte a protocolos personalizados ao empacotamento MSIX, por @leaanthony
- Adiciona detecção do ambiente de desktop no Linux [PR nº 4797](https://github.com/wailsapp/wails/pull/4797)
- Adiciona o método `Window.Print()` ao runtime JavaScript para abrir a caixa de diálogo de impressão pelo frontend (#4290), por @leaanthony
- Adiciona `XDG_SESSION_TYPE` à saída de `wails3 doctor` no Linux, por @leaanthony
- Adiciona outros eventos de alteração de carregamento do WebKit2 no Linux: `WindowLoadStarted`, `WindowLoadRedirected`, `WindowLoadCommitted`, `WindowLoadFinished` (#3896), por @leaanthony
- Adiciona `XDG_SESSION_TYPE` à saída de `wails3 doctor` no Linux, por @leaanthony
- Gera o arquivo `.desktop` durante a compilação para Linux, e não apenas durante o empacotamento (#4575)
- Adiciona documentação sobre as dependências de runtime no Linux, com nomes de pacotes específicos de cada distribuição e exemplos de empacotamento com nfpm (#4339), por @leaanthony
- Adiciona informações sobre a versão do driver NVIDIA à saída de `wails3 doctor` no Linux, por @leaanthony
- Adiciona a origem ao manipulador de mensagens brutas, por @APshenkin, no [PR](https://github.com/wailsapp/wails/pull/4710)
- Adiciona suporte a links universais no macOS, por @APshenkin, no [PR](https://github.com/wailsapp/wails/pull/4712)
- Refatora a camada de transporte de bindings, por @APshenkin, no [PR](https://github.com/wailsapp/wails/pull/4702)
- Adiciona identificadores aria-label aos modelos helloworld para que o aplicativo de exemplo possa ser testado facilmente por clientes de teste Appium, por @chinenual, no [PR](https://github.com/wailsapp/wails/pull/4760)
- Adiciona a origem ao manipulador de mensagens brutas, por @APshenkin, no [PR](https://github.com/wailsapp/wails/pull/4710)
- Adiciona suporte a links universais no macOS, por @APshenkin, no [PR](https://github.com/wailsapp/wails/pull/4712)
- Refatora a camada de transporte de bindings, por @APshenkin, no [PR](https://github.com/wailsapp/wails/pull/4702)
- Eventos tipados, por @fbbdev e @ianvs, em [#4633](https://github.com/wailsapp/wails/pull/4633)
- Adiciona o exemplo `systray-clock`, que demonstra um ícone de bandeja sem interface gráfica com atualizações dinâmicas da dica de ferramenta (#4653).
- Adicionado modelo de protocolo NSIS para Windows, por @Tolfx, em #4510
- Adicionados testes para build-assets, por @Tolfx, em #4510
- macOS: exibe os controles nativos da janela na barra de menus em [#4588](https://github.com/wailsapp/wails/pull/4588), por @nidib
- Adiciona um serviço do Dock do macOS para ocultar/exibir o ícone do aplicativo no Dock, por @popaprozac, no [PR](https://github.com/wailsapp/wails/pull/4451)
- Adiciona um serviço do Dock do macOS para ocultar/exibir o ícone do aplicativo no Dock, por @popaprozac, no [PR](https://github.com/wailsapp/wails/pull/4451)
- Adiciona suporte ao efeito nativo Liquid Glass no macOS com NSGlassEffectView (macOS 15.0+) e fallback para NSVisualEffectView, incluindo opções abrangentes de personalização de materiais, por @leaanthony, em [#4534](https://github.com/wailsapp/wails/pull/4534)
- Sanitização de URLs do navegador, por @leaanthony, em [#4500](https://github.dev/wailsapp/wails/pull/4500). Baseada em [#4484](https://github.com/wailsapp/wails/pull/4484), de @APShenkin.
- Adiciona proteção de conteúdo no Windows/Mac, por [@leaanthony](https://github.com/leaanthony), com base no trabalho original de [@Taiterbase](https://github.com/Taiterbase) neste [PR](https://github.com/wailsapp/wails/pull/4241)
- Adiciona suporte à passagem de variáveis da CLI para comandos Task por meio dos aliases `wails3 build` e `wails3 package` (#4422), por @leaanthony, no [PR](https://github.com/wailsapp/wails/pull/4488)
- Suporte a zonas de soltura, com o evento fornecendo os dados do elemento solto, por [@atterpac](https://github.com/atterpac), em [#4318](https://github.com/wailsapp/wails/pull/4318)
- Adicionada a opção `AdditionalLaunchArgs` às opções de `WindowsWindow` para permitir que argumentos adicionais de linha de comando sejam passados ao navegador WebView2, no [PR](https://github.com/wailsapp/wails/pull/4467)
- Adicionada a execução automática de go mod tidy após wails init, por [@triadmoko](https://github.com/triadmoko), no [PR](https://github.com/wailsapp/wails/pull/4286)
- Recurso Snap Assist do Windows, por @leaanthony, no [PR](https://github.dev/wailsapp/wails/pull/4463)
- Adicionada a opção `AdditionalLaunchArgs` às opções de `WindowsWindow` para permitir que argumentos adicionais de linha de comando sejam passados ao navegador WebView2, no [PR](https://github.com/wailsapp/wails/pull/4467)
- Adicionada a execução automática de go mod tidy após wails init, por [@triadmoko](https://github.com/triadmoko), no [PR](https://github.com/wailsapp/wails/pull/4286)
- Recurso Snap Assist do Windows, por @leaanthony, no [PR](https://github.dev/wailsapp/wails/pull/4463)
- Adiciona a implementação de `getAccentColor` para Windows, por [@almas-x](https://github.com/almas-x), no [PR](https://github.com/wailsapp/wails/pull/4427)
- Adiciona a implementação de `getAccentColor` para Windows, por [@almas-x](https://github.com/almas-x), no [PR](https://github.com/wailsapp/wails/pull/4427)
- Menus e barra de menus com tema escuro no Windows. Por @leaanthony, em [a29b4f0861b1d0a700e9eb213c6f1076ec40efd5](https://github.com/wailsapp/wails/commit/a29b4f0861b1d0a700e9eb213c6f1076ec40efd5)
- Renomeia os serviços integrados para tornar os bindings JS/TS mais claros, por @popaprozac, no [PR](https://github.com/wailsapp/wails/pull/4405)
- `app.Env.GetAccentColor` para obter a cor de destaque do sistema do usuário. Funciona no MacOS. Por [@etesam913](https://github.com/etesam913)
- Adiciona a API `window.ToggleFrameless()`, por [@atterpac](https://github.com/atterpac), em [#4137](https://github.com/wailsapp/wails/pull/4137)
- Adiciona dependências de compilação específicas de cada distribuição para Linux, por @leaanthony, em [PR](https://github.com/wailsapp/wails/pull/4345)
- Adiciona um guia de bindings, por @atterpac, em [PR](https://github.com/wailsapp/wails/pull/4404)
- **Infraestrutura de testes organizada**: os arquivos de teste do Docker foram movidos para o diretório dedicado `test/docker/`, com imagens otimizadas e maior confiabilidade da compilação, por [@leaanthony](https://github.com/leaanthony) em [#4359](https://github.com/wailsapp/wails/pull/4359)
- **Padrões aprimorados de gerenciamento de recursos**: foram adicionados aos exemplos a limpeza adequada de manipuladores de eventos e o gerenciamento de goroutines com reconhecimento de contexto, por [@leaanthony](https://github.com/leaanthony) em [#4359](https://github.com/wailsapp/wails/pull/4359)
- Adiciona suporte à compilação de AppImages para aarch64, por [@AkshayKalose](https://github.com/AkshayKalose) em [#3981](https://github.com/wailsapp/wails/pull/3981)
- Adiciona uma seção de diagnóstico a `wails doctor`, por [@leaanthony](https://github.com/leaanthony)
- Adiciona a janela ao contexto ao chamar um método de serviço, por [@leaanthony](https://github.com/leaanthony)
- Adiciona o exemplo `window-call` para demonstrar como saber qual janela está chamando um serviço, por [@leaanthony](https://github.com/leaanthony)
- Novo guia de menus, por [@leaanthony](https://github.com/leaanthony)
- Melhora o tratamento de panics, por [@leaanthony](https://github.com/leaanthony)
- Novo guia de menus, por [@leaanthony](https://github.com/leaanthony)
- Adiciona comentários de documentação à API de serviços, por [@fbbdev](https://github.com/fbbdev) em [#4024](https://github.com/wailsapp/wails/pull/4024)
- Adiciona a função `application.NewServiceWithOptions` para inicializar serviços com configuração adicional, por [@leaanthony](https://github.com/leaanthony) em [#4024](https://github.com/wailsapp/wails/pull/4024)
- Melhora o controle de menus, por [@FalcoG](https://github.com/FalcoG) e [@leaanthony](https://github.com/leaanthony) em [#4031](https://github.com/wailsapp/wails/pull/4031)
- Mais documentação, por [@leaanthony](https://github.com/leaanthony)
- Adiciona suporte ao cancelamento de eventos em listeners de eventos padrão, por [@leaanthony](https://github.com/leaanthony)
- Adiciona suporte a `Hide`, `Show` e `Destroy` na bandeja do sistema, por [@leaanthony](https://github.com/leaanthony)
- Adiciona suporte a `SetTooltip` na bandeja do sistema, por [@leaanthony](https://github.com/leaanthony). Ideia original de [@lujihong](https://github.com/wailsapp/wails/issues/3487#issuecomment-2633242304)
- Inclui o caminho do pacote nos avisos do gerador de bindings sobre tipos não compatíveis, por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Adiciona ao gerador de bindings suporte a aliases genéricos, por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Adiciona ao gerador de bindings suporte à flag JSON `omitzero`, por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Adiciona a diretiva `//wails:ignore` para impedir a geração de bindings para métodos de serviço selecionados, por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Adiciona a diretiva `//wails:internal` a serviços e modelos para permitir tipos exportados em Go, mas não em JS/TS, por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Adiciona ao gerador de bindings suporte a constantes de tipos alias para permitir enums com tipagem fraca, por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Adiciona testes do gerador de bindings para recursos do Go 1.24, por [@fbbdev](https://github.com/fbbdev) em [#4068](https://github.com/wailsapp/wails/pull/4068)
- Adiciona suporte ao macOS 15 "Sequoia" em `OSInfo.Branding` para melhorar a detecção da versão do sistema operacional, em [#4065](https://github.com/wailsapp/wails/pull/4065)
- Adiciona o hook `PostShutdown` para executar código personalizado após a conclusão do processo de encerramento, por [@fbbdev](https://github.com/fbbdev) em [#4066](https://github.com/wailsapp/wails/pull/4066)
- Adiciona a struct `FatalError` para permitir a detecção de erros fatais em manipuladores de erros personalizados, por [@fbbdev](https://github.com/fbbdev) em [#4066](https://github.com/wailsapp/wails/pull/4066)
- Padroniza e documenta a ordem de inicialização e encerramento dos serviços, por [@fbbdev](https://github.com/fbbdev) em [#4066](https://github.com/wailsapp/wails/pull/4066)
- Adiciona um ambiente de testes para a sequência de inicialização e encerramento do aplicativo e para testes de inicialização e encerramento dos serviços, por [@fbbdev](https://github.com/fbbdev) em [#4066](https://github.com/wailsapp/wails/pull/4066)
- Adiciona o método `RegisterService` para registrar serviços após a criação do aplicativo, por [@fbbdev](https://github.com/fbbdev) em [#4066](https://github.com/wailsapp/wails/pull/4066)
- Adiciona o campo `MarshalError` às opções do aplicativo e dos serviços para o tratamento personalizado de erros em chamadas de bindings, por [@fbbdev](https://github.com/fbbdev) em [#4066](https://github.com/wailsapp/wails/pull/4066)
- Adiciona um wrapper de promise cancelável que propaga solicitações de cancelamento por cadeias de promises, por [@fbbdev](https://github.com/fbbdev) em [#4100](https://github.com/wailsapp/wails/pull/4100)
- Adiciona a capacidade de vincular o cancelamento de chamadas de bindings a um `AbortSignal`, por [@fbbdev](https://github.com/fbbdev) em [#4100](https://github.com/wailsapp/wails/pull/4100)
- Adiciona suporte a atributos `data-wml-*` para WML, juntamente com os atributos `wml-*` usuais, por [@leaanthony](https://github.com/leaanthony)
- Adiciona o método `Configure` a todos os serviços para configuração tardia ou reconfiguração dinâmica, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Quando não está configurado, o serviço `fileserver` envia uma resposta 503 Serviço Indisponível, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Quando não está configurado, o serviço `kvstore` fornece por padrão um armazenamento de chave-valor em memória, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Adiciona o método `Load` ao serviço `kvstore` para recarregar os dados do arquivo após alterações de configuração, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Adiciona o método `Clear` ao serviço `kvstore` para excluir todas as chaves, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Adiciona o tipo `Level` ao serviço `log` para fornecer constantes de nível de log no lado do JS, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Adiciona o método `Log` ao serviço `log` para definir dinamicamente o nível de log, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Quando não está configurado, o serviço `sqlite` fornece por padrão um banco de dados em memória, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Adicionar o método `Close` ao serviço `sqlite` para fechar o banco de dados manualmente, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Adicionar suporte a cancelamento nos métodos de consulta do serviço `sqlite`, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Adicionar suporte a instruções preparadas ao serviço `sqlite`, com bindings JS, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Suporte ao Gin por [Lea Anthony](https://github.com/leaanthony) no [PR](https://github.com/wailsapp/wails/pull/3537), com base no trabalho original de [@AnalogJ](https://github.com/AnalogJ) neste [PR](https://github.com/wailsapp/wails/pull/3537)
- Corrigir o salvamento automático e o salvamento automático de senhas, que estavam sempre habilitados, por [@oSethoum](https://github.com/osethoum) em [#4134](https://github.com/wailsapp/wails/pull/4134)
- Adicionar `SetMenu()` à janela para permitir a definição de um menu nela, por [@leaanthony](https://github.com/leaanthony)
- Adicionar suporte a notificações, por [@popaprozac](https://github.com/popaprozac) em [#4098](https://github.com/wailsapp/wails/pull/4098)
-  Adicionar suporte a associações de arquivos no mac, por [@wimaha](https://github.com/wimaha) em [#4177](https://github.com/wailsapp/wails/pull/4177)
- Adicionar `wails3 tool version` para incrementar versões semânticas, por [@leaanthony](https://github.com/leaanthony)
- Adicionar suporte a badges no macOS e no Windows, por [@popaprozac](https://github.com/popaprozac) em [#](https://github.com/wailsapp/wails/pull/4234)
- Adicionar suporte a eventos registrados e estritamente tipados, por [@fbbdev](https://github.com/fbbdev) e [@IanVS](https://github.com/IanVS) em [#4161](https://github.com/wailsapp/wails/pull/4161)
- Adicionar a capacidade de registrar hooks para eventos personalizados, por [@fbbdev](https://github.com/fbbdev) e [@IanVS](https://github.com/IanVS) em [#4161](https://github.com/wailsapp/wails/pull/4161)
- `app.OpenFileManager(path string, selectFile bool)` para abrir o gerenciador de arquivos do sistema no caminho `path`, com destaque opcional por meio de `selectFile`, por [@Krzysztofz01](https://github.com/Krzysztofz01) [@rcalixte](https://github.com/rcalixte)
- Nova flag `-git` para o comando `wails3 init`, por [@leaanthony](https://github.com/leaanthony)
- Novo comando `wails3 generate webview2bootstrapper`, por [@leaanthony](https://github.com/leaanthony)
- Adicionado o método `init()` ao runtime para permitir sua inicialização manual, por [@leaanthony](https://github.com/leaanthony)
- Adicionada a opção `WindowDidMoveDebounceMS` a WindowOptions da janela, por [@leaanthony](https://github.com/leaanthony)
- Adicionado o recurso de instância única, por [@leaanthony](https://github.com/leaanthony). Baseado no [PR da v2](https://github.com/wailsapp/wails/pull/2951) de @APshenkin.
- Comando `wails3 generate template`, por [@leaanthony](https://github.com/leaanthony)
- Comando `wails3 releasenotes`, por [@leaanthony](https://github.com/leaanthony)
- Comando `wails3 update cli`, por [@leaanthony](https://github.com/leaanthony)
- Opção `-clean` para o comando `wails3 generate bindings`, por [@leaanthony](https://github.com/leaanthony)
- Permitir builds AppImage para Linux aarch64 (arm64), por [@AkshayKalose](https://github.com/AkshayKalose) em [#3981](https://github.com/wailsapp/wails/pull/3981)
- Adicionado hiperlink para o patrocinador por @ansxuman em [#3958](https://github.com/wailsapp/wails/pull/3958)
- Suporte à criação de pacotes deb, rpm e Arch Linux no Linux, por
- Adicionado suporte a builds e pacotes universais para Darwin, por
- Documentação de eventos no site, por
- Templates para sveltekit e sveltekit-ts configurados para desenvolvimento sem SSR
- Atualizar os assets de build usando o novo comando `wails3 update build-assets`, por
- Exemplo para testar a API HTML Drag and Drop, por
- Suporte a associações de arquivos, por [leaanthony](https://github.com/leaanthony) em
- Novo comando `wails3 generate runtime`, por
- Nova opção `InitialPosition` para especificar se a janela deve ser centralizada ou
- Adicionar os métodos `Path` e `Paths` ao pacote `application`, por
- Adicionadas as opções do Windows `GeneralAutofillEnabled` e `PasswordAutosaveEnabled`
- Adicionada a capacidade de recuperar a janela que chamou um método de serviço, por
- Adicionadas as opções `EnabledFeatures` e `DisabledFeatures` para o WebView2, por
- ⊞ Novo sistema DIP para suporte aprimorado a monitores com DPI alto, por
- ⊞ Opção de nome da classe da janela, por [windom](https://github.com/windom/) em
- Os serviços foram ampliados para oferecer funcionalidade de plugins. Por
- 🐧 Eventos WindowDidMove / WindowDidResize em
- ⊞ Evento WindowDidResize em
-  Adicionar o evento ApplicationShouldHandleReopen para permitir o tratamento do Dock
-  Adicionar getPrimaryScreen/getScreens à implementação, por @tmclane em
-  Adicionar opção para exibir a barra de ferramentas no modo de tela cheia no macOS, por
- 🐧 Adicionar lógica onKeyPress para converter o pressionamento de teclas no Linux em um acelerador
- 🐧 Adicionar a tarefa `run:linux`, por
- Exportar o método `SetIcon`, por [@almas-x](https://github.com/almas-x) em
- Aprimorar `OnShutdown`, por [@almas-x](https://github.com/almas-x) em
- Restaurar o método `ToggleMaximise` na interface `Window`, por
- Adicionadas mais informações a `Environment()`. Por @leaanthony em
- Expor o método `WebviewWindow.IsFocused` na interface `Window`, por
- Oferecer suporte a vários eventos de gatilho separados por espaços no sistema WML, por
- Adicionar exportações ESM ao script de runtime JS incluído no bundle, por
- Adicionar ao gerador de bindings uma flag para usar o script de runtime JS incluído no bundle em vez de
- Implementar `setIcon` no Linux, por [@abichinger](https://github.com/abichinger)
- Adicionar a flag `-port` ao comando dev e oferecer suporte à variável de ambiente
- Adicionar testes para chamadas de métodos vinculados por
- ⊞ adicionar `SetIgnoreMouseEvents` para uma janela já criada por
-  Adicionar a capacidade de definir o nível de empilhamento (ordem) de uma janela por

### Corrigido

- Corrigir `Screen.Bounds`, `WorkArea` e `Size` reduzidos pela metade em Macs com tela Retina, convertendo os valores de pontos de NSScreen em pixels do dispositivo nos campos `Physical*`, e preencher `Screen.X`/`Y` no nível superior para que a detecção de monitores adjacentes e o posicionamento na área de trabalho estejam corretos em [PR](https://github.com/wailsapp/wails/pull/5168) por @wayneforrest
- Corrigir a condição de corrida de dados no ScreenManager que causa um deadlock do DisplayLink do WebKit quando a configuração de telas é alterada (por exemplo, conexão a quente de um monitor externo durante a suspensão/retomada)
- Definir CFBundleIconName diretamente como appicon quando Assets.car existir em [PR](https://github.com/wailsapp/wails/pull/5154) por @symball
- Corrigir `wails3 doctor`, que informa pacotes WebKitGTK incorretos no Fedora, openSUSE, Arch e NixOS — as entradas alternativas de 4.0 foram removidas, pois a v3 exige a API 4.1 durante a compilação (#5071)
- Corrigir o nome do pacote webkit2gtk verificado pelo doctor no openSUSE (`webkit2gtk4_1-devel` → `webkit2gtk3-devel`, o nome correto do pacote no openSUSE) (#5071)
- Corrigir o erro `Unexpected token '<'` quando `/wails/custom.js` estiver ausente no modo de desenvolvimento para desktop. Foi adicionado um manipulador explícito de 404 para `/wails/custom.js` e uma validação de `Content-Type` que não diferencia maiúsculas de minúsculas em `loadOptionalScript`, a fim de impedir que fallbacks HTML de SPAs sejam injetados como JavaScript. ([#5068](https://github.com/wailsapp/wails/issues/5068))
- Corrigir o estado de destaque do menu da bandeja do sistema no macOS — o ícone agora exibe o estado selecionado quando o menu está aberto (#4910)
- Corrigir a janela vinculada à bandeja do sistema que aparecia atrás de outras janelas no macOS — agora é usado o nível de janela pop-up apropriado (#4910)
- Corrigir exemplos de importação de `@wailsio/runtime` incorretos em toda a documentação (#4989)
- Corrigir a impossibilidade de minimizar janelas sem moldura no Darwin (#4294)
- Corrige travamentos de 20-30 minutos durante `wails3 build` e `wails3 dev`, excluindo `node_modules/` da verificação do go-task que determina se as tarefas estão atualizadas. Anteriormente, o padrão glob `sources: "**/*"` fazia o go-task enumerar e calcular a soma de verificação de todos os arquivos em `node_modules/` (50000-100000+ arquivos com dependências pesadas, como MUI), o que era especialmente lento no Windows/NTFS (#4939)
- Corrigir a falha de compilação com GTK4 causada pela colisão do typedef C `Screen` com X11 Xlib.h (#4957)
- Corrigir a consistência dos métodos do selo do Dock no macOS
- Corrigir `InvisibleTitleBarHeight`, que era aplicado a todas as janelas do macOS em vez de apenas às janelas sem moldura ou com barra de título transparente (#4960)
- Corrigir a trepidação da janela ao redimensioná-la pelos cantos superiores com `InvisibleTitleBarHeight` habilitado, ignorando o início do arraste perto das bordas da janela (#4960)
- Corrigir a geração de tipos mapeados com chaves de enumerações nos bindings JS/TS (#4437) por @fbbdev
- Corrigir o recurso de arrastar e soltar arquivos no Windows, que não funcionava com escala de exibição diferente de 100%
- Corrigir o recurso interno de arrastar e soltar do HTML5, que não funcionava quando a soltura de arquivos estava habilitada no Windows
- Corrigir as coordenadas de soltura de arquivos, que usavam o espaço de pixels incorreto no Windows (pixels físicos versus pixels CSS)
- Corrigir o funcionamento inconsistente do recurso de arrastar e soltar arquivos com efeitos de foco do ponteiro no Linux
- Corrigir o recurso interno de arrastar e soltar do HTML5, que não funcionava quando a soltura de arquivos estava habilitada no Linux
- Corrigir a exibição/ocultação de janelas no Linux/GTK4, que às vezes as restaurava no estado minimizado, usando `gtk_window_present()` (#4957)
- Corrigir a obtenção/definição da posição da janela no Linux/GTK4, que sempre retornava 0,0, adicionando suporte condicional ao X11 por meio de `XTranslateCoordinates`/`XMoveWindow` (#4957)
- Corrigir a falta de aplicação do tamanho máximo da janela no Linux/GTK4, adicionando uma limitação de tamanho baseada em sinais para substituir o `gtk_window_set_geometry_hints` removido (#4957)
- Corrigir a escala de DPI no Linux/GTK4, implementando o cálculo adequado de PhysicalBounds e o suporte à escala fracionária por meio de `gdk_monitor_get_scale` (GTK 4.14+)
- Corrigir a duplicação de itens de menu ao criar novas janelas no Linux/GTK4
- Corrigir a geração de tipos mapeados com chaves de enumerações nos bindings JS/TS (#4437) por @fbbdev
- Corrigir o recurso de arrastar e soltar arquivos no Windows, que não funcionava com escala de exibição diferente de 100%
- Corrigir o recurso interno de arrastar e soltar do HTML5, que não funcionava quando a soltura de arquivos estava habilitada no Windows
- Corrigir as coordenadas de soltura de arquivos, que usavam o espaço de pixels incorreto no Windows (pixels físicos em vez de pixels CSS)
- Corrigir o funcionamento inconsistente do recurso de arrastar e soltar arquivos com efeitos de foco do ponteiro no Linux
- Corrigir o recurso interno de arrastar e soltar do HTML5, que não funcionava quando a soltura de arquivos estava habilitada no Linux
- Corrigir a escala de DPI no Linux/GTK4, implementando o cálculo adequado de PhysicalBounds e o suporte à escala fracionária por meio de `gdk_monitor_get_scale` (GTK 4.14+)
- Corrigir a duplicação de itens de menu ao criar novas janelas no Linux/GTK4
- Corrigir a geração de tipos mapeados com chaves de enumerações nos bindings JS/TS (#4437) por @fbbdev
- Corrigir o problema de “janelas fantasmas” no macOS, causado pelo acesso às APIs do AppKit fora da thread principal em App.Window.Current() (#4947) por @wimaha
- Corrigir o `<input type="file">` do HTML, que não funcionava no macOS, implementando WKUIDelegate runOpenPanelWithParameters (#4862)
- Corrigir o recurso nativo de arrastar e soltar arquivos, que não funcionava ao usar o módulo npm `@wailsio/runtime` no macOS/Linux (#4953) por @leaanthony
- Corrigir a geração de bindings para aliases de tipos entre pacotes (#4578) por @fbbdev
- Corrigir a falha do OpenFileDialog no Linux causada por uma violação da segurança de threads do GTK (#3683) por @ddmoney420
- Corrigir a falha SIGSEGV ao chamar `Focus()` em uma janela oculta ou destruída (#4890) por @ddmoney420
- Corrigir um possível panic ao definir um ícone ou bitmap vazio no Linux (#4923) por @ddmoney420
- Corrigir a falha do ErrorDialog quando chamado por um binding de serviço no macOS (#3631) por @leaanthony
- Fazer com que os menus sejam exibidos no sistema operacional Windows em `v3\examples\dialogs` por @ndianabasi
- Corrigir a condição de corrida que causava um TypeError durante o recarregamento da página (#4872) por @ddmoney420
- Corrigir a saída incorreta dos testes do gerador de bindings removendo o estado global do método `Collector.IsVoidAlias()` (#4941) por @fbbdev
- Corrigir o seletor de arquivos `<input type="file">`, que não funcionava no macOS (#4862) por @leaanthony
- Corrige o uso de sistemas de coordenadas inconsistentes por `Position()` e `SetPosition()` no macOS, que causava o deslocamento da posição da janela ao salvar/restaurar o estado (#4816), por @leaanthony
- Corrige o erro "Acesso negado" de SetProcessDpiAwarenessContext quando a percepção de DPI já está definida pelo manifesto do aplicativo (#4803)
- Atualiza a página da documentação sobre atalhos de teclado e corrige o tipo do parâmetro de retorno de chamada de `KeyBinding.Add`, por @ndianabasi
- Corrige a documentação sobre a geração de bindings personalizados: é necessário usar `-d String` em vez de `-o String`
- Corrige a falha que impedia o menu de remover os itens filhos em `menu.Update()`
- Corrige referências desatualizadas à API Manager na documentação (31 arquivos foram atualizados para usar o novo padrão, como `app.Window.New()`, `app.Event.Emit()` etc.), por @leaanthony
- Corrige uma falha no Linux quando ocorre um panic em métodos Go vinculados ao JS, devido à substituição dos manipuladores de sinais pelo WebKit (#3965), por @leaanthony
- Corrige SaveFileDialog.SetFilename() sem efeito no Linux (#4841), por @samstanier
- Corrige as coordenadas de soltura exibidas como undefined no exemplo de arrastar e soltar
- Corrige a falha na criação do pacote do aplicativo no macOS quando APP_NAME contém espaços (problema de expansão de chaves)
- Corrige um panic de índice fora dos limites no Windows ao chamar métodos de serviço (reverte goccy/go-json)
- Corrige o recurso de arrastar e soltar arquivos no Windows, que não funcionava com escala de exibição diferente de 100%
- Corrige o recurso interno de arrastar e soltar do HTML5, que deixava de funcionar quando a soltura de arquivos era habilitada no Windows
- Corrige as coordenadas de soltura de arquivos no Windows, que usavam o espaço de pixels incorreto (pixels físicos em vez de pixels CSS)
- Corrige o funcionamento instável do recurso de arrastar e soltar arquivos no Linux com efeitos de foco do ponteiro
- Corrige o recurso interno de arrastar e soltar do HTML5, que deixava de funcionar quando a soltura de arquivos era habilitada no Linux
- Atualiza todos os comandos nos arquivos Taskfile.yml de todos os sistemas operacionais para aceitar espaços em variáveis como `APP_NAME`, por @ndianabasi
- Corrige um erro nos argumentos do comando ao executar a tarefa 'build:universal:lipo:go' no Linux, por @wux1an
- Corrige o erro do Docker "símbolo indefinido: **<em>ubsan</em>handle_xxxxxxx" ao executar 'wails3 build GOOS=darwin GOARCH=arm64' no Linux, por @wux1an
- Consolida a documentação de protocolos personalizados e adiciona seções sobre Universal Links, por @leaanthony
- Corrige uma falha no menu da bandeja do sistema no Windows ao clicar repetidamente no ícone, adicionando uma proteção contra chamadas simultâneas de TrackPopupMenuEx (#4151), por @leaanthony
- Impede uma falha do aplicativo ao chamar systray.Run() antes de app.Run(), por @leaanthony
- Corrige uma falha no macOS ao alternar a visibilidade da janela por meio de Hide()/Show() com ApplicationShouldTerminateAfterLastWindowClosed habilitado (#4389), por @leaanthony
- Corrige um vazamento de memória nos menus de contexto do macOS e do Windows ao abri-los repetidamente (#4012), por @leaanthony
- Corrige a falta de reutilização dos recursos nativos dos menus de contexto no macOS, que fazia com que um novo menu fosse criado a cada exibição (#4012), por @leaanthony
- Corrige o clique no ícone do Dock do macOS, que não exibia as janelas ocultas quando o aplicativo era iniciado com `Hidden: true` (#4583), por @leaanthony
- Corrige a caixa de diálogo de impressão do macOS, que não abria devido a um tipo incorreto de ponteiro de janela na chamada CGO (#4290), por @leaanthony
- Corrige uma falha no menu da janela no Wayland, causada pelo acesso do appmenu-gtk-module a uma janela ainda não realizada (#4769), por @leaanthony
- Corrige uma falha do aplicativo GTK quando o nome do aplicativo contém caracteres inválidos (espaços, parênteses etc.), por @leaanthony
- Corrige o erro "memória insuficiente" ao inicializar o recurso de arrastar e soltar no Windows (#4701), por @overlordtm
- Corrige a abertura do diretório errado pelo explorador de arquivos no Linux devido ao escape incorreto da URI (#4397), por @leaanthony
- Corrige a falha na compilação do AppImage em distribuições Linux modernas (Arch, Fedora 39+, Ubuntu 24.04+) por meio da detecção automática de seções ELF `.relr.dyn` e da desativação da remoção de símbolos (#4642), por @leaanthony
- Corrige `wails doctor`, que indicava incorretamente que os pacotes do webkit estavam instalados em sistemas baseados em Fedora/DNF (#4457), por @leaanthony
- Corrige o comportamento padrão de `config.yml`, que executava `wails3 dev` com uma compilação de produção, por @mbaklor
- Corrige os stubs de serviço do iOS que causavam falhas de compilação devido à importação de um pacote inexistente, por @leaanthony
- Corrige o registro estruturado nos métodos de depuração/informação que causava erros de "nenhuma diretiva de formatação", por @leaanthony
- Remove instruções temporárias de impressão de depuração incluídas acidentalmente durante a mesclagem da plataforma móvel, por @leaanthony
- Corrige uma falha do WebKitGTK no Wayland com GPUs NVIDIA (erro 71 Protocol error) por meio da desativação automática do renderizador DMA-BUF, por @leaanthony
- Corrige o valor alfa ignorado em `application.WebviewWindowOptions.BackgroundColour` no Linux ([#4722](https://github.com/wailsapp/wails/pull/4722), @BradHacker)
- Corrige o ícone da bandeja do sistema no Windows, que não usava por padrão o ícone do aplicativo quando nenhum ícone personalizado era fornecido (#4704)
- Rastreia a propriedade de `HICON` para que apenas os identificadores criados pelo usuário sejam destruídos, evitando falhas quando o Explorer é reiniciado (#4653).
- Libera o listener do tema do sistema do Windows e os ícones da bandeja retidos durante a destruição para impedir o vazamento de goroutines e contextos de dispositivo (#4653).
- Trunca as dicas de ferramenta da bandeja em 127 unidades UTF-16 para evitar a corrupção de pares substitutos e glifos multibyte (#4653).
- Corrige a falha na tarefa de empacotamento para Windows (#4667)
- Corrige a variável appicon do AppImage para Linux no taskfile do Linux [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- Corrige um erro de compilação no Windows causado pela alteração de assinatura no go-webview2 v1.0.22 (#4513, #4645)
- Corrige a variável appicon do AppImage para Linux no taskfile do Linux [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- Corrige a iteração sobre protocolos em desktop.tmpl para Linux, substituindo `<.Info.Protocol>` por `<.Protocol>`, por @Tolfx em #4510
- Corrige o erro de redefinição na demonstração de liquid glass em [#4542](https://github.com/wailsapp/wails/pull/4542), por @Etesam913
- Corrige as atualizações do menu da bandeja do sistema no Linux [#4604](https://github.com/wailsapp/wails/issues/4604), por [@JackDoan](https://github.com/JackDoan)
- Corrige a janela branca que aparecia no Windows ao criar uma janela oculta, por @leaanthony em [#4612](https://github.com/wailsapp/wails/pull/4612)
- Corrige o caminho de importação do pacote de notificações na documentação, por @rxliuli em [#4617](https://github.com/wailsapp/wails/pull/4617)
- Corrige o recurso de arrastar e soltar, que não funcionava ao usar o pacote npm @wailsio/runtime (#4489), por @leaanthony em #4616
- Windows: corrigidos a cintilação da janela na inicialização e o problema que fazia janelas ocultas serem exibidas incorretamente em [PR](https://github.com/wailsapp/wails/pull/4600), por @leaanthony.
- Corrigidos problemas de maximização do tamanho da janela no Wayland (https://github.com/wailsapp/wails/issues/4429), por [@samstanier](https://github.com/samstanier)
- Corrigidos problemas de maximização do tamanho da janela no Wayland (https://github.com/wailsapp/wails/issues/4429), por [@samstanier](https://github.com/samstanier)
- Corrigido o erro de redefinição na demonstração de liquid glass em [#4542](https://github.com/wailsapp/wails/pull/4542), por @Etesam913
- Corrigido o problema que podia causar uma falha no AssetServer no macOS em [#4576](https://github.com/wailsapp/wails/pull/4576), por @jghiloni
- Corrigido o problema de compilação ao criar builds com NextJs. Correção feita em [#4585](https://github.com/wailsapp/wails/pull/4585), por @rev42
- Corrigidos os pipelines da versão nightly em [#4597](https://github.com/wailsapp/wails/pull/4597), por @riadafridishibly
- Corrigido o erro de redefinição na demonstração de liquid glass em [#4542](https://github.com/wailsapp/wails/pull/4542), por @Etesam913
- Corrigido o problema que podia causar uma falha no AssetServer no macOS em [#4576](https://github.com/wailsapp/wails/pull/4576), por @jghiloni
- Corrigido o problema de compilação ao criar builds com NextJs. Correção feita em [#4585](https://github.com/wailsapp/wails/pull/4585), por @rev42
- Corrigidos os pipelines da versão nightly em [#4597](https://github.com/wailsapp/wails/pull/4597), por @riadafridishibly
- Corrigido o erro de redefinição na demonstração de liquid glass em [#4542](https://github.com/wailsapp/wails/pull/4542), por @Etesam913
- Corrige SetBackgroundColour no Windows, por @PPTGamer em [PR](https://github.com/wailsapp/wails/pull/4492)
- Atualiza a documentação para refletir as alterações decorrentes da refatoração da API Manager, por @yulesxoxo em [PR #4476](https://github.com/wailsapp/wails/pull/4476)
- Corrige a variável appicon do arquivo .desktop do Linux no Taskfile do Linux em [PR #4477](https://github.com/wailsapp/wails/pull/4477)
- Atualiza a documentação para refletir as alterações decorrentes da refatoração da API Manager, por @yulesxoxo em [PR #4476](https://github.com/wailsapp/wails/pull/4476)
- Corrige o erro de desreferenciamento de ponteiro nil no Windows relatado em [#4456](https://github.com/wailsapp/wails/issues/4456), por @leaanthony em [#4460](https://github.com/wailsapp/wails/pull/4460)
- Adiciona suporte a `allowsBackForwardNavigationGestures` no WKWebView do macOS para habilitar gestos de navegação por deslizamento com dois dedos (#1857)
- Corrige o problema que impedia onClick de funcionar em itens de menu inicialmente definidos como desabilitados, por @leaanthony em [PR #4469](https://github.com/wailsapp/wails/pull/4469). Agradecimentos a @IanVS pela investigação inicial.
- Corrige o problema que impedia a limpeza do servidor Vite quando o build falhava (#4403)
- Corrigido o panic ao fechar ou cancelar um `SaveFileDialog` no Windows. Correção feita em [PR](https://github.com/wailsapp/wails/pull/4284), por @hkhere
- Corrigido o recurso de arrastar e soltar no nível do HTML no Windows, por [@mbaklor](https://github.com/mbaklor) em [#4259](https://github.com/wailsapp/wails/pull/4259)
- Adiciona suporte a `allowsBackForwardNavigationGestures` no WKWebView do macOS para habilitar gestos de navegação por deslizamento com dois dedos (#1857)
- Corrige o problema que impedia onClick de funcionar em itens de menu inicialmente definidos como desabilitados, por @leaanthony em [PR #4469](https://github.com/wailsapp/wails/pull/4469). Agradecimentos a @IanVS pela investigação inicial.
- Corrige o problema que impedia a limpeza do servidor Vite quando o build falhava (#4403)
- Corrigida a análise de notificações no Windows, por @popaprozac em [PR](https://github.com/wailsapp/wails/pull/4450)
- Corrigido o comando doctor para verificar as dependências do Windows SDK, por [@kodumulo](https://github.com/kodumulo) em [#4390](https://github.com/wailsapp/wails/issues/4390)
- Corrigido o desreferenciamento de ponteiro nil em processURLRequest no Mac, por [@etesam913](https://github.com/etesam913) em [#4366](https://github.com/wailsapp/wails/pull/4366)
- Corrigido um bug no Linux que impedia o uso de caixas de diálogo com filtros, por [@bh90210](https://github.com/bh90210) em [#4287](https://github.com/wailsapp/wails/pull/4287)
- Corrigidos problemas do menu Editar no Windows e no Linux, por [@leaanthony](https://github.com/leaanthony) em [#3f78a3a](https://github.com/wailsapp/wails/commit/3f78a3a8ce7837e8b32242c8edbbed431c68c062)
- Atualizada de 10.13.0 para 10.15.0 a versão mínima do sistema nos arquivos .plist do macOS, por [@AkshayKalose](https://github.com/AkshayKalose) em [#3981](https://github.com/wailsapp/wails/pull/3981)
- Corrigido o problema de salto de IDs de janela, por [@leaanthony](https://github.com/leaanthony)
- Corrige o problema de menu nil ao chamar RegisterContextMenu, por [@leaanthony](https://github.com/leaanthony)
- Corrigidos ciclos de dependência na saída do gerador de bindings, por [@fbbdev](https://github.com/fbbdev) em [#4001](https://github.com/wailsapp/wails/pull/4001)
- Corrigidos erros de uso antes da definição na saída do gerador de bindings, por [@fbbdev](https://github.com/fbbdev) em [#4001](https://github.com/wailsapp/wails/pull/4001)
- Repassa os sinalizadores de build ao gerador de bindings, por [@fbbdev](https://github.com/fbbdev) em [#4023](https://github.com/wailsapp/wails/pull/4023)
- Altera os caminhos no Taskfile do Windows para usar barras normais, garantindo que ele funcione em plataformas diferentes do Windows, por [@leaanthony](https://github.com/leaanthony)
- Corrigidos os eventos do Mac e do JavaScript no Mac, por [@leaanthony](https://github.com/leaanthony)
- Corrigido o deadlock de eventos no macOS, por [@leaanthony](https://github.com/leaanthony)
- Corrigido um erro `Parameter incorrect` na inicialização de Window no Windows quando havia HTML, mas nenhum JS, por [@leaanthony](https://github.com/leaanthony)
- Corrigido o tamanho do prefixo da resposta usado para detectar o tipo de conteúdo no servidor de assets, por [@fbbdev](https://github.com/fbbdev) em [#4049](https://github.com/wailsapp/wails/pull/4049)
- Corrigido o tratamento de respostas que não são 404 no caminho do índice raiz do servidor de assets, por [@fbbdev](https://github.com/fbbdev) em [#4049](https://github.com/wailsapp/wails/pull/4049)
- Corrigido o comportamento indefinido no gerador de bindings ao testar propriedades de tipos genéricos, por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Corrigida a saída do gerador de bindings para modelos quando o tipo subjacente não tem as mesmas propriedades que o wrapper nomeado, por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Corrigida a saída do gerador de bindings para tipos de chave de mapas e para o pré-processamento, por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Corrigida a saída do gerador de bindings para structs que implementam interfaces de marshaling, por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Corrigida a detecção de ciclos de tipos que envolvem tipos genéricos no gerador de bindings, por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Corrigidas referências inválidas a modelos não exportados na saída do gerador de bindings por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Código injetado movido para o final dos arquivos de serviço por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Corrigido o tratamento de erros das operações de fechamento de arquivos no gerador de bindings por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Suprimidos os avisos para serviços que definem métodos de ciclo de vida ou HTTP, mas nenhum outro método vinculado, por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Corrigida a falha dos templates que não usam React ao exibir o rodapé Hello World quando o esquema de cores claro do sistema era usado, por [@marcus-crane](https://github.com/marcus-crane) em [#4056](https://github.com/wailsapp/wails/pull/4056)
- Corrigidos os itens de menu ocultos no macOS por [@leaanthony](https://github.com/leaanthony)
- Corrigidos o tratamento e a formatação de erros nos processadores de mensagens por [@fbbdev](https://github.com/fbbdev) em [#4066](https://github.com/wailsapp/wails/pull/4066)
-  Corrigido o encerramento de serviços que era ignorado ao sair do aplicativo, por [@fbbdev](https://github.com/fbbdev) em [#4066](https://github.com/wailsapp/wails/pull/4066)
-  Garantido que as atualizações de menu ocorram na thread principal por [@leaanthony](https://github.com/leaanthony)
- O mecanismo de arrastar e redimensionar agora é mais robusto e corresponde melhor ao comportamento esperado da plataforma, por [@fbbdev](https://github.com/fbbdev) em [#4100](https://github.com/wailsapp/wails/pull/4100)
- Corrigido [#4097](https://github.com/wailsapp/wails/issues/4097): Webpack/Angular descarta o código de inicialização do runtime, por [@fbbdev](https://github.com/fbbdev) em [#4100](https://github.com/wailsapp/wails/pull/4100)
- Corrigidos os itens de menu inicialmente ocultos por [@IanVS](https://github.com/IanVS) em [#4116](https://github.com/wailsapp/wails/pull/4116)
- Corrigido o assetFileServer, que não servia arquivos `.html` em uma solicitação sem extensão quando `[request]` não existia, mas `[request].html` existia
- Corrigidos os caminhos de geração de ícones por [@robin-samuel](https://github.com/robin-samuel) em [#4125](https://github.com/wailsapp/wails/pull/4125)
- Corrigida a não emissão dos eventos `fullscreen`, `unfullscreen`, `unminimise` e `unmaximise` por [@oSethoum](https://github.com/osethoum) em [#4130](https://github.com/wailsapp/wails/pull/4130)
- Corrigido o erro do NSIS causado por um prefixo incorreto na versão padrão da configuração, por [@robin-samuel](https://github.com/robin-samuel) em [#4126](https://github.com/wailsapp/wails/pull/4126)
- Corrigida a função de runtime Dialogs, que retornava caminhos com caracteres de escape no Windows, por [TheGB0077](https://github.com/TheGB0077) em [#4188](https://github.com/wailsapp/wails/pull/4188)
- Corrigido o caminho de detecção do Webview2 em HKCU por [@leaanthony](https://github.com/leaanthony).
- Corrigido um problema de entrada no macOS por [@leaanthony](https://github.com/leaanthony).
- Corrigido o nome do arquivo da tarefa de geração de ícones do Windows por [@yulesxoxo](https://github.com/yulesxoxo) em [#4219](https://github.com/wailsapp/wails/pull/4219).
- Corrigido um problema de transparência em janelas sem moldura por [@leaanthony](https://github.com/leaanthony), com base no trabalho de @kron.
- Corrigidas as chamadas de foco quando a janela está desabilitada ou minimizada por [@leaanthony](https://github.com/leaanthony), com base no trabalho de @kron.
- Corrigidos os ícones da bandeja do sistema que não reapareciam após a reinicialização da barra de tarefas, por [@leaanthony](https://github.com/leaanthony), com base no trabalho de @kron.
- Corrigido fallbackResponseWriter, que não implementava Flush(), em [#4245](https://github.com/wailsapp/wails/pull/4245)
- Corrigido fallbackResponseWriter, que não implementava Flush(), por [@superDingda] em [#4236](https://github.com/wailsapp/wails/issues/4236)
- Corrigidas falhas ao fechar uma janela no macOS enquanto havia uma chamada assíncrona pendente a uma função Go vinculada, por [@joshhardy](https://github.com/joshhardy) em [#4354](https://github.com/wailsapp/wails/pull/4354)
- Corrigida uma condição de corrida na inicialização do modo de eficiência do Windows por [@leaanthony](https://github.com/leaanthony)
- Corrigida a liberação dos identificadores de ícones do Windows por [@leaanthony](https://github.com/leaanthony).
- Corrigido `OpenFileManager` no Windows por [@PPTGamer](https://github.com/PPTGamer) em [#4375](https://github.com/wailsapp/wails/pull/4375).
- Corrigidas as opções de largura mínima/máxima no Linux por @atterpac em [#3979](https://github.com/wailsapp/wails/pull/3979)
- Corrigidas as definições de tipos dos templates TypeScript por meio de uma atualização da versão no npm, por @atterpac em [#3966](https://github.com/wailsapp/wails/pull/3966)
- Corrigida a referência ao CSS no template SvelteKit por @atterpac em [#3945](https://github.com/wailsapp/wails/pull/3945)
- Garantido que os principais callbacks em run() da janela sejam chamados na thread principal por [@leaanthony](https://github.com/leaanthony)
- Corrigidos os exemplos do seletor de diretórios da caixa de diálogo por [@leaanthony](https://github.com/leaanthony)
- Criada uma nova página de erro em chinês para quando index.html estiver ausente, por [@leaanthony](https://github.com/leaanthony)
-  Garantido que o callback `windowDidBecomeKey` seja executado na thread principal por [@leaanthony](https://github.com/leaanthony)
-  Adicionado suporte a tela cheia para janelas sem moldura por [@leaanthony](https://github.com/leaanthony)
-  Aprimorada a lógica de destruição de janelas por [@leaanthony](https://github.com/leaanthony)
-  Corrigida a lógica de posicionamento da janela quando vinculada a ícones da bandeja do sistema, por [@leaanthony](https://github.com/leaanthony)
-  Adicionado suporte a tela cheia para janelas sem moldura por [@leaanthony](https://github.com/leaanthony)
- Corrigido o tratamento de eventos por [@leaanthony](https://github.com/leaanthony)
- Corrigida a lógica de encerramento de janelas por [@leaanthony](https://github.com/leaanthony)
- O taskfile comum agora gera bindings TypeScript por padrão para templates TypeScript, por [@leaanthony](https://github.com/leaanthony)
- Corrigido o encerramento do aplicativo ao receber a mensagem WM_CLOSE quando nenhuma janela está aberta ou há somente um ícone na bandeja do sistema, por [@mmalcek](https://github.com/mmalcek) em [#3990](https://github.com/wailsapp/wails/pull/3990)
- Corrigida a compilação com garble por @5aaee9 em [#3192](https://github.com/wailsapp/wails/pull/3192)
- Corrigidas as compilações NSIS para Windows por [@leaanthony](https://github.com/leaanthony)
- Corrigido um deadlock na caixa de diálogo do Linux para múltiplas seleções, causado por algo não fechado
- Corrigida a limpeza multiplataforma de arquivos .syso durante a compilação para Windows por
- Corrigida a compilação do AppImage para amd64 por @atterpac em
- Corrigida a atualização dos assets de compilação por @ansxuman em
- Corrigida a implementação de `OnClick` e `OnRightClick` da bandeja do sistema no Linux por @atterpac
- Corrigido o problema de `AlwaysOnTop` não funcionar no Mac por
-  Corrigido `application.NewEditMenu` que incluía uma duplicata
- 🐧 Corrigida a compilação para aarch64
- ⊞ Corrigidos os itens de menu de grupos de opções por
- Corrigido o erro ao compilar um .app executável no macOS quando 'name' e 'outputfilename'
- Corrigido o bug no uso de customEventProcessor no exemplo de arrastar e soltar por
- 🐧 Corrigido o erro de compilação no Linux introduzido pela adição de IgnoreMouseEvents por
- ⊞ Corrigido o bug na geração do arquivo de ícone syso por
- 🐧 Incorporada a correção para execução nativa no Wayland proveniente de
- Não vincular métodos internos dos serviços em
- ⊞ Corrigido o panic ao iniciar a bandeja do sistema em
- Não vincular métodos internos dos serviços em
- ⊞ Corrigido o panic ao iniciar a bandeja do sistema em
- Grande refatoração dos itens de menu e do tratamento de eventos. Por enquanto, melhora principalmente o macOS. Por
- Corrigidos os testes após a refatoração de plugins e eventos em
- ⊞ Corrigido o aviso de `Failed to unregister class Chrome_WidgetWin_0`. Por
- Problemas de módulos
- Corrigido o envio de mensagens de eventos de redimensionamento por [atterpac](https://github.com/atterpac) em
- 🐧 Corrigido o erro no tratamento de temas no NixOS por
- Corrigida a instalação do projeto entre volumes no Windows por
- Corrigido o CSS do modelo React para exibir o rodapé por
- Corrigidos os processos zumbis ao trabalhar no modo de desenvolvimento com a atualização para a versão mais recente do refresh
- Corrigida a obtenção do arquivo do WebKit no AppImage por [Atterpac](https://github.com/atterpac)
- Corrigida a verificação de pacotes apt pelo Doctor por [Atterpac](https://github.com/Atterpac) em
- Corrigido o congelamento do aplicativo ao encerrar (Darwin) por @5aaee9 em
- Corrigidas as cores de fundo dos exemplos no Windows por
- Corrigidos os menus de contexto padrão por [mmghv](https://github.com/mmghv) em
- Corrigidos os valores hexadecimais das teclas de seta no Darwin por
- O recurso de arrastar e soltar no Windows passou a funcionar. Adicionado por
- Corrigido um bug no Doctor para Linux quando o usuário não tem os drivers adequados
- Corrigido o dimensionamento por DPI na inicialização (Windows). Alterado por [@almas-x](https://github.com/almas-x) em
- Corrigida a linha de substituição em `go.mod` para usar caminhos relativos. Corrige caminhos do Windows com
- Corrigido o tratamento de cliques na bandeja do sistema do macOS quando não há uma janela associada por
- Corrigida a falha de compilação no Windows causada por uma opção desconhecida por
- Corrigida a falha no Windows ao clicar com o botão esquerdo no ícone da bandeja do sistema quando não há um
- Corrigida a baseURL incorreta ao abrir a janela duas vezes por @5aaee9 no PR
- Corrigida a ordem das ramificações if no método `WebviewWindow.Restore` por
- Calculado corretamente `startURL` em várias invocações de `GetStartURL` quando
- Corrigido o tipo JS da struct `Screen` para corresponder à sua equivalente em Go por
- Corrigido o método `WML.Reload` para garantir a limpeza adequada dos eventos registrados
- Corrigido o fechamento imediato do menu de contexto personalizado no Linux por
- Corrigidos o caminho de saída e a extensão dos arquivos de modelo produzidos pela vinculação
- Corrigidos os caminhos de importação dos arquivos de modelo no código JS produzido pela vinculação
- Corrigido o recurso de arrastar e soltar em algumas distribuições Linux por
- Corrigida a tarefa ausente no macOS ao usar `wails3 task dev` por
- Corrigido o registro de eventos que causava uma atribuição em um mapa nil por
- Corrigida a desserialização dos parâmetros de métodos vinculados por
- Corrigido o tratamento de vários valores de retorno de métodos vinculados por
- Corrigida a detecção, pelo Doctor, do npm que não foi instalado com o gerenciador de pacotes do sistema
- Corrigida a ausência de MicrosoftEdgeWebview2Setup.exe. Agradecimentos a
- Corrigida uma falha aleatória no Linux causada pelo tratamento do ID da janela por @leaanthony. Baseado em
- Corrigida a falha de systemTray.setIcon no Linux por
- Corrigido para garantir que a moldura da janela seja aplicada na primeira chamada da função `setFrameless` em

### Alterado

- **INCOMPATÍVEL**: as chaves de mapas nas vinculações JS/TS geradas agora são marcadas como opcionais para refletir com precisão a semântica dos mapas em Go. O acesso aos valores de mapas no TypeScript agora retorna `T | undefined` em vez de `T`, exigindo verificações de null ou asserções (#4943) por `@fbbdev`
- Alterado o uso de `Event` para `Events` de acordo com as alterações em `@wailsio/runtime`, bem como as chamadas de função correspondentes na documentação em `Features/Events/Event System`, por @AbdelhadiSeddar
- Movidos `EnabledFeatures`, `DisabledFeatures` e `AdditionalBrowserArgs` das opções específicas de cada janela para `Options.Windows` no nível do aplicativo (#4559) por @leaanthony
- Atualizado o README do exemplo `Drag N Drop`, destacando que `Internal Drag and Drop` é demonstrado no exemplo, por @ndianabasi
- Alterados vários logs de depuração de Info para Debug (por @mbaklor)
- **INCOMPATÍVEL:** renomeado `EnableDragAndDrop` para `EnableFileDrop` nas opções da janela
- **INCOMPATÍVEL:** renomeado `DropZoneDetails` para `DropTargetDetails` no contexto do evento
- **INCOMPATÍVEL:** renomeado o método `DropZoneDetails()` para `DropTargetDetails()` em `WindowEventContext`
- **INCOMPATÍVEL:** removido o evento `WindowDropZoneFilesDropped`; use `WindowFilesDropped` em seu lugar
- **ALTERAÇÃO INCOMPATÍVEL:** Alterar o atributo HTML de `data-wails-dropzone` para `data-file-drop-target`
- **ALTERAÇÃO INCOMPATÍVEL:** Alterar a classe CSS de hover de `wails-dropzone-hover` para `file-drop-target-active`
- **ALTERAÇÃO INCOMPATÍVEL:** Remover as opções `DragEffect`, `OnEnterEffect` e `OnOverEffect` do Windows (faziam parte da IDropTarget removida)
- Migrar para goccy/go-json todo o processamento de JSON em tempo de execução (vinculações de métodos, eventos, solicitações da webview, notificações e kvstore), melhorando o desempenho em 21-63% e reduzindo as alocações de memória em 40-60%
- Otimizar o layout da struct BoundMethod e armazenar em cache o sinalizador isVariadic para reduzir a sobrecarga por chamada
- Usar um buffer de argumentos alocado na pilha para métodos com `<=8` argumentos, evitando alocações no heap
- Otimizar a coleta de resultados em chamadas de métodos para evitar a alocação de uma slice quando houver um único valor de retorno
- Usar sync.Map no cache de tipos MIME para melhorar o desempenho com acessos simultâneos
- Usar um pool de buffers para ler o corpo das solicitações do transporte HTTP
- Alocar de forma adiada o canal CloseNotify no detector de tipo de conteúdo para reduzir as alocações por solicitação
- Remover do servidor de recursos o registro de depuração de CSS
- Expandir o mapa de extensões de tipos MIME para abranger mais de 50 formatos comuns da Web (fontes, áudio, vídeo etc.)
- Atualizar a documentação das opções `X/Y` de Window — @ruhuang2001
- Atualizar a documentação de `Frontend Runtime` adicionando mais opções para gerar vinculações do frontend — @ndianabasi
- Atualizar a página de documentação do Asset Server do Wails v3 — @ndianabasi
- **ALTERAÇÃO INCOMPATÍVEL**: Remover as funções de diálogo no nível do pacote (`application.InfoDialog()`, `application.QuestionDialog()` etc.). Usar o gerenciador `app.Dialog`: `app.Dialog.Info()`, `app.Dialog.Question()`, `app.Dialog.Warning()`, `app.Dialog.Error()`, `app.Dialog.OpenFile()` e `app.Dialog.SaveFile()`
- Atualizar a documentação de diálogos para corresponder à API real: usar `app.Dialog.*`; `AddButton()` com callbacks (não `SetButtons()`); `SetDefaultButton(*Button)` (não uma string); `AddFilter()` (não `SetFilters()`); `SetFilename()` (não `SetDefaultFilename()`); e `app.Dialog.OpenFile().CanChooseDirectories(true)` para selecionar pastas
- **ALTERAÇÃO INCOMPATÍVEL**: As builds de produção agora são o padrão. Para criar builds de desenvolvimento, defina `DEV=true` nos seus Taskfiles. Gere um novo projeto para obter exemplos — @leaanthony
- Ao emitir um evento personalizado com zero ou um argumento de dados, o valor dos dados será atribuído diretamente ao campo Data, sem ser encapsulado em uma slice — [@fbbdev](https://github.com/fbbdev) em [#4633](https://github.com/wailsapp/wails/pull/4633)
- Os ícones da bandeja no Windows agora respeitam `SystemTray.Show()`/`Hide()` alternando `NIS_HIDDEN`, permitindo que os aplicativos desapareçam de fato e reapareçam (#4653).
- O registro do ícone da bandeja reutiliza os ícones resolvidos, define `NOTIFYICON_VERSION_4` uma única vez e habilita `NIF_SHOWTIP` para que as dicas de ferramenta sejam restauradas após a reinicialização do Explorer (#4653).
- macOS: Usar `visibleFrame` em vez de `frame` para centralizar a janela, excluindo as áreas da barra de menus e do Dock
- macOS: Usar `visibleFrame` em vez de `frame` para centralizar a janela, excluindo as áreas da barra de menus e do Dock
- Ao executar `wails3 update build-assets` com o parâmetro `-config`, os valores definidos por meio dos parâmetros `-product*` são
- `window.NativeWindowHandle()` -> `window.NativeWindow()` — @leaanthony em [#4471](https://github.com/wailsapp/wails/pull/4471)
- Refatorar o tratamento interno de janelas — @leaanthony em [#4471](https://github.com/wailsapp/wails/pull/4471)
- Removidos `application.WindowIDKey` e `application.WindowNameKey` (substituídos por `application.WindowKey`) — [@leaanthony](https://github.com/leaanthony)
- ContextMenuData agora retorna uma string em vez de any — [@leaanthony](https://github.com/leaanthony)
- Nas vinculações JS/TS, os campos de classe com tipos de array de tamanho fixo agora são inicializados com o tamanho esperado, em vez de ficarem vazios — [@fbbdev](https://github.com/fbbdev) em [#4001](https://github.com/wailsapp/wails/pull/4001)
- ContextMenuData agora retorna uma string em vez de any — [@leaanthony](https://github.com/leaanthony)
- `application.NewService` não aceita mais opções como parâmetro opcional (use `application.NewServiceWithOptions`) — [@leaanthony](https://github.com/leaanthony) em [#4024](https://github.com/wailsapp/wails/pull/4024)
- Removida a dependência `nanoid` — [@leaanthony](https://github.com/leaanthony)
- Atualizado o exemplo de Window para os estilos de janela mica/acrylic/tabbed — [@leaanthony](https://github.com/leaanthony)
- Nas vinculações JS/TS, os arquivos de modelo `internal.js/ts` foram removidos; agora todos os modelos estão em `models.js/ts` — [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Nas vinculações JS/TS, tipos nomeados nunca são renderizados como aliases de outros tipos nomeados; o comportamento anterior agora se restringe a aliases — [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Nas vinculações JS/TS em modo de classe, os campos de structs cujo tipo é um parâmetro de tipo são marcados como opcionais e nunca são inicializados automaticamente — [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Remover o ESLint dos modelos — [@IanVS](https://github.com/IanVS) em [#4059](https://github.com/wailsapp/wails/pull/4059)
- Atualizar o ano dos direitos autorais para 2025 — [@IanVS](https://github.com/IanVS) em [#4037](https://github.com/wailsapp/wails/pull/4037)
- Adicionar documentação para event.Sender — [@IanVS](https://github.com/IanVS) em [#4075](https://github.com/wailsapp/wails/pull/4075)
- Compatibilidade com Go 1.24 — [@leaanthony](https://github.com/leaanthony)
- Os hooks de `ServiceStartup` agora são invocados quando `App.Run` é chamado, e não em `application.New` — [@fbbdev](https://github.com/fbbdev) em [#4066](https://github.com/wailsapp/wails/pull/4066)
- Os erros de `ServiceStartup` agora são retornados por `App.Run`, em vez de encerrarem o processo — [@fbbdev](https://github.com/fbbdev) em [#4066](https://github.com/wailsapp/wails/pull/4066)
- As chamadas de vinculações e diálogos provenientes de JS agora são rejeitadas com objetos de erro, em vez de strings — [@fbbdev](https://github.com/fbbdev) em [#4066](https://github.com/wailsapp/wails/pull/4066)
- Melhorado o posicionamento do menu da bandeja do sistema no Windows — [@leaanthony](https://github.com/leaanthony)
- O runtime JS foi migrado para TypeScript — [@fbbdev](https://github.com/fbbdev) em [#4100](https://github.com/wailsapp/wails/pull/4100)
- O runtime é inicializado assim que é importado, sem necessidade de aguardar o carregamento da janela, por [@fbbdev](https://github.com/fbbdev) em [#4100](https://github.com/wailsapp/wails/pull/4100)
- O runtime não exporta mais um método init. Uma importação apenas para efeitos colaterais pode ser usada para inicializá-lo, por [@fbbdev](https://github.com/fbbdev) em [#4100](https://github.com/wailsapp/wails/pull/4100)
- Os métodos vinculados agora retornam uma `CancellablePromise` que é rejeitada com um `CancelError` em caso de cancelamento. O resultado efetivo da chamada é descartado, por [@fbbdev](https://github.com/fbbdev) em [#4100](https://github.com/wailsapp/wails/pull/4100)
- Os tipos de serviço integrados agora são chamados consistentemente de `Service`, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- As funções integradas de criação de serviços com opções agora são chamadas consistentemente de `NewWithConfig`, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- O método `Select` do serviço `sqlite` agora se chama `Query`, para manter a consistência com as APIs Go, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Templates: runtime movido para "dependencies" e arquivos package.json organizados, por [@IanVS](https://github.com/IanVS) em [#4133](https://github.com/wailsapp/wails/pull/4133)
- Cria e assina de forma ad hoc os pacotes do aplicativo durante o desenvolvimento para habilitar determinadas APIs do macOS, por [@popaprozac](https://github.com/popaprozac) em [#4171](https://github.com/wailsapp/wails/pull/4171)
- Ativos de compilação movidos para diretórios específicos de cada plataforma, por [@leaanthony](https://github.com/leaanthony)
- Taskfiles movidos e renomeados para diretórios específicos de cada plataforma, por [@leaanthony](https://github.com/leaanthony)
- Experiência significativamente melhor quando `index.html` está ausente, por [@leaanthony](https://github.com/leaanthony)
- [Windows] Melhor desempenho ao minimizar e restaurar, por [@leaanthony](https://github.com/leaanthony). Baseado no [PR](https://github.com/wailsapp/wails/pull/3955) original de [562589540](https://github.com/562589540)
- Opção `ShouldClose` removida (registre um hook para events.Common.WindowClosing em vez disso), por [@leaanthony](https://github.com/leaanthony)
- [Windows] Redução da cintilação ao abrir uma janela, por [@leaanthony](https://github.com/leaanthony)
- `Window.Destroy` removida porque deveria ser uma função interna, por [@leaanthony](https://github.com/leaanthony)
- Eventos `WindowClose` renomeados como `WindowClosing`, por [@leaanthony](https://github.com/leaanthony)
- As compilações do frontend agora usam o ambiente vite "development" ou "production", conforme o tipo de compilação, por [@leaanthony](https://github.com/leaanthony)
- Atualização para go-webview2 v1.19, por [@leaanthony](https://github.com/leaanthony)
- Garantia de uso do fork do taskfile, por @leaanthony
- Atualização do fork do Taskfile para corrigir problemas de versão ao instalar usando
- Uso do fork do Taskfile para corrigir problemas de versão ao instalar usando
- `service.OnStartup` agora encerra o aplicativo em caso de erro e executa
- Refatoração das mensagens de clique na bandeja do sistema para corresponder melhor às interações do usuário, por
- Incorporação de ativos passou a incluir `all:frontend/dist` para oferecer suporte a frameworks que geram
- Refatoração do Taskfile por [leaanthony](https://github.com/leaanthony) em
- Atualização para `go-webview2` v1.0.16 por
- Correção do tipo `Screen` para incluir `ID`, e não `Id`, por
- Atualização da versão do Wails em `go.mod.tmpl` para oferecer suporte a `application.ServiceOptions`, por
- Correção da determinação do nome do serviço por [windom](https://github.com/windom/) em
- mkdocs serve agora usa docker, por [leaanthony](https://github.com/leaanthony)
- Configuração de desenvolvimento consolidada em `config.yml`, por
- A caixa de diálogo da bandeja do sistema agora usa por padrão o ícone do aplicativo, quando disponível (Windows), por
- Melhor relatório de GPU + memória no macOS, por
- `WebviewGpuIsDisabled` e `EnableFraudulentWebsiteWarnings` removidos
- Alteração da API de eventos: `On`/`Emit` -> eventos do usuário, `OnApplicationEvent` ->
- Correção da API de eventos no Linux por [TheGB0077](https://github.com/TheGB0077) em
- [CI] melhorias nas actions e habilitação da execução de actions também em forks e
- `AbsolutePosition()` renomeado como `Position()`, por
- Atualização da dependência do WebKit no Linux para webkit2gtk-4.1 em vez de webkitgtk2-4.0 para
- O script do runtime JS incluído no pacote agora é um módulo ESM: as tags script que o importam
- O pacote `@wailsio/runtime` não publica sua API no `window.wails`
- O módulo `@wailsio/runtime/src/window` da API de janelas agora expõe a janela que contém
- A API de janelas em JS foi atualizada para corresponder à `WebviewWindow` atual do Go
- O gerador de bindings agora usa chamadas por ID por padrão. A opção da CLI `-id`
- Novo layout do código de bindings: anteriormente, os arquivos de saída eram organizados em pastas
- O campo `application.Options.Bind` da struct foi renomeado como
- Nova sintaxe para vincular serviços: agora as instâncias de serviço devem ser encapsuladas em um
- Desativação do indicador giratório em ambientes sem terminal ou de CI, por

### Removido

- **QUEBRA DE COMPATIBILIDADE**: remoção de `EnabledFeatures`, `DisabledFeatures` e `AdditionalLaunchArgs` das opções `WindowsWindow` específicas de cada janela. Em vez disso, use `Options.Windows.EnabledFeatures`, `Options.Windows.DisabledFeatures` e `Options.Windows.AdditionalBrowserArgs` no nível do aplicativo. Essas flags se aplicam globalmente ao ambiente WebView2 compartilhado (#4559), por @leaanthony
- Remoção da implementação nativa de `IDropTarget` no Windows em favor da abordagem baseada em JavaScript (corresponde ao comportamento da v2)
- Remoção da dependência github.com/wailsapp/mimetype em favor de um mapa de extensões ampliado + http.DetectContentType da biblioteca padrão, reduzindo o tamanho do binário em ~1.2MB
- Remoção da dependência gopkg.in/ini.v1 por meio da implementação de um parser mínimo de arquivos .desktop para o explorador de arquivos do Linux, economizando ~45KB
- Remover samber/lo do código de runtime usando o pacote slices da biblioteca padrão do Go 1.21+ e auxiliares internos mínimos, economizando cerca de 310 KB
- Remover instruções printf de depuração do manipulador de esquemas de URL do Darwin (#4834)
- **QUEBRA DE COMPATIBILIDADE**: remover o evento `linux:WindowLoadChanged`; usar `linux:WindowLoadFinished` para detectar quando o WebView terminar de carregar (#3896), por @leaanthony

### Quebras de compatibilidade

- **Refatoração da API de gerenciadores**: a API da aplicação foi reorganizada, passando de uma estrutura plana para gerenciadores organizados, a fim de melhorar a organização e facilitar a descoberta do código, por [@leaanthony](https://github.com/leaanthony) em [#4359](https://github.com/wailsapp/wails/pull/4359)
- `app.NewWebviewWindow()` → `app.Window.New()`
- `app.CurrentWindow()` → `app.Window.Current()`
- `app.GetAllWindows()` → `app.Window.GetAll()`
- `app.WindowByName()` → `app.Window.GetByName()`
- `app.EmitEvent()` → `app.Event.Emit()`
- `app.OnApplicationEvent()` → `app.Event.OnApplicationEvent()`
- `app.OnWindowEvent()` → `app.Event.OnWindowEvent()`
- `app.SetApplicationMenu()` → `app.Menu.SetApplicationMenu()`
- `app.OpenFileDialog()` → `app.Dialog.OpenFile()`
- `app.SaveFileDialog()` → `app.Dialog.SaveFile()`
- `app.MessageDialog()` → `app.Dialog.Message()`
- `app.InfoDialog()` → `app.Dialog.Info()`
- `app.WarningDialog()` → `app.Dialog.Warning()`
- `app.ErrorDialog()` → `app.Dialog.Error()`
- `app.QuestionDialog()` → `app.Dialog.Question()`
- `app.NewSystemTray()` → `app.SystemTray.New()`
- `app.GetSystemTray()` → `app.SystemTray.Get()`
- `app.ShowContextMenu()` → `app.ContextMenu.Show()`
- `app.RegisterKeybinding()` → `app.KeyBinding.Register()`
- `app.UnregisterKeybinding()` → `app.KeyBinding.Unregister()`
- `app.GetPrimaryScreen()` → `app.Screen.GetPrimary()`
- `app.GetAllScreens()` → `app.Screen.GetAll()`
- `app.BrowserOpenURL()` → `app.Browser.OpenURL()`
- `app.Environment()` → `app.Env.GetAll()`
- `app.ClipboardGetText()` → `app.Clipboard.Text()`
- `app.ClipboardSetText()` → `app.Clipboard.SetText()`
- Métodos de Service renomeados: `Name` -> `ServiceName`, `OnStartup` -> `ServiceStartup`, `OnShutdown` -> `ServiceShutdown`, por [@leaanthony](https://github.com/leaanthony)
- Métodos `Path` e `Paths` movidos para o pacote `application`, por [@leaanthony](https://github.com/leaanthony)
- O menu da aplicação agora está disponível somente no macOS, por [@leaanthony](https://github.com/leaanthony)

## v3.0.0-alpha.78 - 2026-04-21

## Adicionado

## Corrigido

## v3.0.0-alpha.77 - 2026-04-18

## Corrigido

## v3.0.0-alpha.76 - 2026-04-17

## Corrigido

## v3.0.0-alpha.75 - 2026-04-16

## Corrigido

## v3.0.0-alpha.74 - 2026-03-01

## Adicionado

## Corrigido

## v3.0.0-alpha.73 - 2026-02-27

## Corrigido

## v3.0.0-alpha.72 - 2026-02-16

## Corrigido

## v3.0.0-alpha.71 - 2026-02-10

## Adicionado

## Corrigido

## v3.0.0-alpha.70 - 2026-02-09

## Adicionado

## Corrigido

## v3.0.0-alpha.69 - 2026-02-08

## Adicionado

## Corrigido

## v3.0.0-alpha.68 - 2026-02-07

## Adicionado

## Alterado

## Corrigido

## v3.0.0-alpha.67 - 2026-02-04

## Adicionado

## Alterado

## Corrigido

## v3.0.0-alpha.66 - 2026-02-03

## Adicionado

## Alterado

## Corrigido

## Removido

## v3.0.0-alpha.65 - 2026-02-01

## Adicionado

## v3.0.0-alpha.64 - 2026-01-26

## Adicionado

## v3.0.0-alpha.63 - 2026-01-25

## Corrigido

## v3.0.0-alpha.62 - 2026-01-22

## Corrigido

## v3.0.0-alpha.61 - 2026-01-20

## Corrigido

## v3.0.0-alpha.60 - 2026-01-14

## Corrigido

## v3.0.0-alpha.59 - 2026-01-11

## Alterado

## v3.0.0-alpha.58 - 2026-01-09

## Corrigido

## v3.0.0-alpha.57 - 2026-01-05

## Alterado

## Corrigido

## v3.0.0-alpha.56 - 2026-01-04

## Adicionado

## Alterado

## Corrigido

## Removido

## v3.0.0-alpha.55 - 2026-01-02

## Alterado

## Corrigido

## Removido

## v3.0.0-alpha.54 - 2025-12-29

## Adicionado

## Corrigido

## Removido

## v3.0.0-alpha.53 - 2025-12-27

## Adicionado

## Corrigido

## v3.0.0-alpha.52 - 2025-12-26

## Corrigido

## v3.0.0-alpha.51 - 2025-12-23

## Corrigido

## v3.0.0-alpha.50 - 2025-12-21

## Alterado

## v3.0.0-alpha.49 - 2025-12-18

## Alterado

## v3.0.0-alpha.48 - 2025-12-16

## Adicionado

## Alterado

## Corrigido

## v3.0.0-alpha.47 - 2025-12-15

## Adicionado

## Corrigido

## v3.0.0-alpha.46 - 2025-12-14

## Adicionado

## Removido

## v3.0.0-alpha.45 - 2025-12-13

## Adicionado

## Corrigido

## v3.0.0-alpha.44 - 2025-12-12

## Adicionado

## Alterado

## Corrigido

## v3.0.0-alpha.43 - 2025-12-11

## Adicionado

## v3.0.0-alpha.42 - 2025-12-10

## Adicionado

## v3.0.0-alpha.41 - 2025-11-23

## Corrigido

## v3.0.0-alpha.40 - 2025-11-13

## Corrigido

## v3.0.0-alpha.39 - 2025-11-12

## Adicionado

## Alterado

## v3.0.0-alpha.38 - 2025-11-04

## Adicionado

## Alterado

## Corrigido

## v3.0.0-alpha.37 - 2025-11-02

## Corrigido

## v3.0.0-alpha.36 - 2025-10-15

## Corrigido

## v3.0.0-alpha.35 - 2025-10-14

## Corrigido

## v3.0.0-alpha.34 - 2025-10-06

## Adicionado

## Corrigido

## v3.0.0-alpha.33 - 2025-10-04

## Corrigido

## v3.0.0-alpha.32 - 2025-10-02

## Corrigido

## v3.0.0-alpha.31 - 2025-09-27

## Corrigido

## v3.0.0-alpha.30 - 2025-09-26

## Corrigido

## v3.0.0-alpha.29 - 2025-09-25

## Adicionado

## Alterado

## Corrigido

## v3.0.0-alpha.29 - 2025-09-25

## Adicionado

## Alterado

## Corrigido

## v3.0.0-alpha.27 - 2025-09-07

## Corrigido

## v3.0.0-alpha.26 - 2025-08-24

## Adicionado

## v3.0.0-alpha.25 - 2025-08-16

## Alterado

deixam de ser ignorados e substituem o valor da configuração.

## v3.0.0-alpha.24 - 2025-08-13

## Adicionado

## v3.0.0-alpha.23 - 2025-08-11

## Corrigido

## v3.0.0-alpha.22 - 2025-08-10

## Adicionado

## Alterado

+ Corrige dependências excessivamente abrangentes de pacotes Linux e dependências de RPM desatualizadas.

## v3.0.0-alpha.21 - 2025-08-07

## Corrigido

## v3.0.0-alpha.20 - 2025-08-06

## Corrigido

## v3.0.0-alpha.19 - 2025-08-05

## Adicionado

## Corrigido

## v3.0.0-alpha.18 - 2025-08-03

## Adicionado

## Corrigido

## v3.0.0-alpha.17 - 2025-07-31

## Corrigido

## v3.0.0-alpha.16 - 2025-07-25

## Adicionado

## v3.0.0-alpha.15 - 2025-07-25

## Adicionado

## v3.0.0-alpha.14 - 2025-07-25

## Adicionado

## v3.0.0-alpha.12 - 2025-07-15

### Adicionado

### Corrigido

## v3.0.0-alpha.11 - 2025-07-12

## Adicionado

## v3.0.0-alpha.10 - 2025-07-06

### Alterações incompatíveis

### Adicionado

### Corrigido

### Alterado

## v3.0.0-alpha.9 - 2025-01-13

### Adicionado

### Corrigido

### Alterado

## v3.0.0-alpha.8.3 - 2024-12-07

### Alterado

## v3.0.0-alpha.8.2 - 2024-12-07

### Alterado

`go install` por @leaanthony

## v3.0.0-alpha.8.1 - 2024-12-07

### Alterado

`go install` por @leaanthony

## v3.0.0-alpha.8 - 2024-12-06

### Adicionado

@atterpac em [#3909](https://github.com/wailsapp/wails/3909)   [ansxuman](https://github.com/ansxuman) em   [#3902](https://github.com/wailsapp/wails/pull/3902)   [atterpac](https://github.com/atterpac) em   [#3867](https://github.com/wailsapp/wails/pull/3867)   por [atterpac](https://github.com/atterpac) em   [#3829](https://github.com/wailsapp/wails/pull/3829)   [leaanthony](https://github.com/leaanthony)   [FerroO2000](https://github.com/FerroO2000) em   [#3856](https://github.com/wailsapp/wails/pull/3856)   [#3873](https://github.com/wailsapp/wails/pull/3873)   [leaanthony](https://github.com/leaanthony)   posicionado nas coordenadas X/Y especificadas, por   [leaanthony](https://github.com/leaanthony) em   [#3885](https://github.com/wailsapp/wails/pull/3885)   [ansxuman](https://github.com/ansxuman) e   [leaanthony](https://github.com/leaanthony) em   [#3823](https://github.com/wailsapp/wails/pull/3823)   por [leaanthony](https://github.com/leaanthony) em   [#3766](https://github.com/wailsapp/wails/pull/3766)   [leaanthony](https://github.com/leaanthony) em   [#3888](https://github.com/wailsapp/wails/pull/3888)   [leaanthony](https://github.com/leaanthony). -

### Alterado

`service.OnShutdown`para quaisquer serviços que já tivessem sido iniciados, por @atterpac em   [#3920](https://github.com/wailsapp/wails/pull/3920)   @atterpac em [#3907](https://github.com/wailsapp/wails/pull/3907)   subpastas por @atterpac em   [#3887](https://github.com/wailsapp/wails/pull/3887)   [#3748](https://github.com/wailsapp/wails/pull/3748)   [leaanthony](https://github.com/leaanthony)   [etesam913](https://github.com/etesam913) em   [#3778](https://github.com/wailsapp/wails/pull/3778)   [northes](https://github.com/northes) em   [#3836](https://github.com/wailsapp/wails/pull/3836)   [#3827](https://github.com/wailsapp/wails/pull/3827)   [leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   (substituído pelas opções `EnabledFeatures` e `DisabledFeatures`) por   [leaanthony](https://github.com/leaanthony)

### Corrigido

variável de canal por @michael-freling em   [#3925](https://github.com/wailsapp/wails/pull/3925)   [ansxuman](https://github.com/ansxuman) em   [#3924](https://github.com/wailsapp/wails/pull/3924)   [#3898](https://github.com/wailsapp/wails/pull/3898)   [#3901](https://github.com/wailsapp/wails/pull/3901)   em [#3886](https://github.com/wailsapp/wails/pull/3886)   [leaanthony](https://github.com/leaanthony) em   [#3841](https://github.com/wailsapp/wails/pull/3841)   função `PasteAndMatchStyle` no menu de edição no Darwin por   [johnmccabe](https://github.com/johnmccabe) em   [#3839](https://github.com/wailsapp/wails/pull/3839)   [#3840](https://github.com/wailsapp/wails/issues/3840) em   [#3854](https://github.com/wailsapp/wails/pull/3854) por   [kodflow](https://github.com/kodflow)   [@leaanthony](https://github.com/leaanthony)   são diferentes. Por @nickisworking em   [#3789](https://github.com/wailsapp/wails/pull/3789)

## v3.0.0-alpha.7 - 2024-09-18

### Adicionado

[mmghv](https://github.com/mmghv) em   [#3665](https://github.com/wailsapp/wails/pull/3665)   [#3682](https://github.com/wailsapp/wails/pull/3682)   [atterpac](https://github.com/atterpac) e   [leaanthony](https://github.com/leaanthony) em   [#3570](https://github.com/wailsapp/wails/pull/3570)

### Alterado

Eventos do aplicativo `OnWindowEvent` -> Eventos da janela, por   [leaanthony](https://github.com/leaanthony)   [#3734](https://github.com/wailsapp/wails/pull/3734)   branches com os prefixos `v3/` ou `v3-` por   [stendler](https://github.com/stendler) em   [#3747](https://github.com/wailsapp/wails/pull/3747)

### Corrigido

[etesam913](https://github.com/etesam913) em   [#3742](https://github.com/wailsapp/wails/pull/3742)   [atterpac](https://github.com/atterpac) em   [#3721](https://github.com/wailsapp/wails/pull/3721)   [atterpac](https://github.com/atterpac) em   [#3675](https://github.com/wailsapp/wails/pull/3675)   [#1811](https://github.com/wailsapp/wails/pull/1811) em   [#3614](https://github.com/wailsapp/wails/pull/3614) por   [@stendler](https://github.com/stendler)   [#3720](https://github.com/wailsapp/wails/pull/3720) por   [leaanthony](https://github.com/leaanthony)   [#3693](https://github.com/wailsapp/wails/issues/3693) por   [@DeltaLaboratory](https://github.com/DeltaLaboratory)   [#3720](https://github.com/wailsapp/wails/pull/3720) por   [leaanthony](https://github.com/leaanthony)   [#3693](https://github.com/wailsapp/wails/issues/3693) por   [@DeltaLaboratory](https://github.com/DeltaLaboratory)   [leaanthony](https://github.com/leaanthony)   [#3746](https://github.com/wailsapp/wails/pull/3746) por   [@stendler](https://github.com/stendler)   [leaanthony](https://github.com/leaanthony)

## v3.0.0-alpha.6 - 2024-07-30

### Corrigido

## v3.0.0-alpha.5 - 2024-07-30

### Adicionado

[#3580](https://github.com/wailsapp/wails/pull/3580)   [#3580](https://github.com/wailsapp/wails/pull/3580)   clique no ícone por @5aaee9 em [#2991](https://github.com/wailsapp/wails/pull/2991)   [#2618](https://github.com/wailsapp/wails/pull/2618)   [@fbbdev](https://github.com/fbbdev) em   [#3282](https://github.com/wailsapp/wails/pull/3282)   @[Atterpac](https://github.com/Atterpac)   em [#3022](https://github.com/wailsapp/wails/pull/3022])   [@marcus-crane](https://github.com/marcus-crane) em   [#3146](https://github.com/wailsapp/wails/pull/3146)   [PR](https://github.com/wailsapp/wails/pull/3147)   [PR](https://github.com/wailsapp/wails/pull/3189)   [@fbbdev](https://github.com/fbbdev) em   [#3281](https://github.com/wailsapp/wails/pull/3281)   [aba82cc](https://github.com/wailsapp/wails/commit/aba82cc52787c97fb99afa58b8b63a0004b7ff6c)   baseado no [PR](https://github.com/wailsapp/wails/pull/2044) de @Mai-Lapyst   [@fbbdev](https://github.com/fbbdev) em   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@fbbdev](https://github.com/fbbdev) em   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@fbbdev](https://github.com/fbbdev) em   [#3295](https://github.com/wailsapp/wails/pull/3295)   o pacote npm por [@fbbdev](https://github.com/fbbdev) em   [#3334](https://github.com/wailsapp/wails/pull/3334)   em [#3354](https://github.com/wailsapp/wails/pull/3354)   `WAILS_VITE_PORT` por [@abichinger](https://github.com/abichinger) em   [#3429](https://github.com/wailsapp/wails/pull/3429)   [@abichinger](https://github.com/abichinger) em   [#3431](https://github.com/wailsapp/wails/pull/3431)   [@bruxaodev](https://github.com/bruxaodev) em   [#3667](https://github.com/wailsapp/wails/pull/3667)   [@OlegGulevskyy](https://github.com/OlegGulevskyy) em   [#3674](https://github.com/wailsapp/wails/pull/3674)

### Corrigido

[#3606](https://github.com/wailsapp/wails/pull/3606)   [tmclane](https://github.com/tmclane) em   [#3515](https://github.com/wailsapp/wails/pull/3515)   [atterpac](https://github.com/atterac) em   [#3512](https://github.com/wailsapp/wails/pull/3512)   [atterpac](https://github.com/atterpac) em   [#3477](https://github.com/wailsapp/wails/pull/3477)   por [Atterpac](https://github.com/atterpac) em   [#3320](https://github.com/wailsapp/wails/pull/3320).   em [#3306](https://github.com/wailsapp/wails/pull/3306).   [#2972](https://github.com/wailsapp/wails/pull/2972).   [#2982](https://github.com/wailsapp/wails/pull/2982)   [mmghv](https://github.com/mmghv) em   [#2750](https://github.com/wailsapp/wails/pull/2750).   [#2753](https://github.com/wailsapp/wails/pull/2753).   [jaybeecave](https://github.com/jaybeecave) em   [#3052](https://github.com/wailsapp/wails/pull/3052).   [@pylotlight](https://github.com/pylotlight) em   [PR](https://github.com/wailsapp/wails/pull/3039)   instalado. Adicionado por [@pylotlight](https://github.com/pylotlight) em   [PR](https://github.com/wailsapp/wails/pull/3032)   [PR](https://github.com/wailsapp/wails/pull/3145)   espaços — @leaanthony.   [thomas-senechal](https://github.com/thomas-senechal) no PR   [#3207](https://github.com/wailsapp/wails/pull/3207)   [thomas-senechal](https://github.com/thomas-senechal) no PR   [#3208](https://github.com/wailsapp/wails/pull/3208)   janela anexada [tw1nk](https://github.com/tw1nk) no PR   [#3271](https://github.com/wailsapp/wails/pull/3271)   [#3273](https://github.com/wailsapp/wails/pull/3273)   [@fbbdev](https://github.com/fbbdev) em   [#3279](https://github.com/wailsapp/wails/pull/3279)   `FRONTEND_DEVSERVER_URL` está presente.   [#3299](https://github.com/wailsapp/wails/pull/3299)   [@fbbdev](https://github.com/fbbdev) em   [#3295](https://github.com/wailsapp/wails/pull/3295)   listeners por [@fbbdev](https://github.com/fbbdev) em   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@abichinger](https://github.com/abichinger) em   [#3330](https://github.com/wailsapp/wails/pull/3330)   gerador por [@fbbdev](https://github.com/fbbdev) em   [#3334](https://github.com/wailsapp/wails/pull/3334)   gerador por [@fbbdev](https://github.com/fbbdev) em   [#3334](https://github.com/wailsapp/wails/pull/3334)   [@abichinger](https://github.com/abichinger) em   [#3346](https://github.com/wailsapp/wails/pull/3346)   [@hfoxy](https://github.com/hfoxy) em   [#3417](https://github.com/wailsapp/wails/pull/3417)   [@hfoxy](https://github.com/hfoxy) em   [#3426](https://github.com/wailsapp/wails/pull/3426)   [@fbbdev](https://github.com/fbbdev) em   [#3431](https://github.com/wailsapp/wails/pull/3431)   [@fbbdev](https://github.com/fbbdev) em   [#3431](https://github.com/wailsapp/wails/pull/3431)   por [@pekim](https://github.com/pekim) em   [#3458](https://github.com/wailsapp/wails/pull/3458)   [@robin-samuel](https://github.com/robin-samuel).   PR [#3466](https://github.com/wailsapp/wails/pull/3622) por   [@5aaee9](https://github.com/5aaee9).   [@windom](https://github.com/windom/) em   [#3636](https://github.com/wailsapp/wails/pull/3636).   Windows por [@bruxaodev](https://github.com/bruxaodev/) em   [#3691](https://github.com/wailsapp/wails/pull/3691).

### Alterado

[mmghv](https://github.com/mmghv) em   [#3611](https://github.com/wailsapp/wails/pull/3611)   compatibilidade com Ubuntu 24.04 LTS por [atterpac](https://github.com/atterpac) em   [#3461](https://github.com/wailsapp/wails/pull/3461)   deve ter o atributo `type="module"`. Por   [@fbbdev](https://github.com/fbbdev) em   [#3295](https://github.com/wailsapp/wails/pull/3295)   objeto e não inicia o sistema WML. Isso foi feito para melhorar   o encapsulamento. Se desejado, o sistema WML pode ser iniciado manualmente chamando   o novo método `WML.Enable`. O script do runtime JS incluído ainda executa ambas   as operações automaticamente. Por [@fbbdev](https://github.com/fbbdev) em   [#3295](https://github.com/wailsapp/wails/pull/3295)   objeto window como exportação padrão. Não é mais possível importar   métodos individuais por meio da sintaxe de importação nomeada ou por namespace do ESM.   API. Alguns métodos tiveram o nome ou o protótipo alterado, especificamente: `Screen`   passa a ser `GetScreen`; `GetZoomLevel`/`SetZoomLevel` passam a ser `GetZoom`/`SetZoom`;   `GetZoom`, `Width` e `Height` agora retornam valores diretamente, em vez de encapsulá-los   em objetos. Por [@fbbdev](https://github.com/fbbdev) em   [#3295](https://github.com/wailsapp/wails/pull/3295)   foi removido. Use a opção `-names` da CLI para voltar às chamadas por nome.   Por [@fbbdev](https://github.com/fbbdev) em   [#3468](https://github.com/wailsapp/wails/pull/3468)   nomeados de acordo com o pacote que os contém; agora são usados os caminhos completos de importação do Go,   incluindo o caminho do módulo. Por [@fbbdev](https://github.com/fbbdev) em   [#3468](https://github.com/wailsapp/wails/pull/3468)   `application.Options.Services`. Por [@fbbdev](https://github.com/fbbdev) em   [#3468](https://github.com/wailsapp/wails/pull/3468)   chamada para `application.NewService`. Por [@fbbdev](https://github.com/fbbdev) em   [#3468](https://github.com/wailsapp/wails/pull/3468)   [@DeltaLaboratory](https://github.com/DeltaLaboratory) em   [#3574](https://github.com/wailsapp/wails/pull/3574)
