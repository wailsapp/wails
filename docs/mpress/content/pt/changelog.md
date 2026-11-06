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

O formato é baseado no [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), e este projeto segue o [Versionamento Semântico](https://semver.org/spec/v2.0.0.html).

- `Added` para novos recursos.
- `Changed` para alterações em funcionalidades existentes.
- `Deprecated` para recursos que serão removidos em breve.
- `Removed` para recursos já removidos.
- `Fixed` para quaisquer correções de bugs.
- `Security` em caso de vulnerabilidades.

_/

/_   * ATUALIZADO AUTOMATICAMENTE *   As notas de lançamento são geradas a partir de PRs mesclados. Não edite este arquivo em PRs de contribuição comuns nem recrie o `v3/UNRELEASED_CHANGELOG.md` removido. Consulte `AGENTS.md` para saber como manter o registro de alterações. _/

## [Não lançado]

## v3.2.0-nightly.20261102 - 2026-11-02

## Alterado
- Usar DisableOverscroll em vez de DisableBounce no [PR](https://github.com/wailsapp/wails/pull/6077) de @lethusang

## v3.2.0-nightly.20261101 - 2026-11-01

## Adicionado
- adicionar backends experimentais sem CGO para macOS e Linux no [PR](https://github.com/wailsapp/wails/pull/6521) de @taliesin-ai
- compilações YAML nativas, empacotamento portátil e assinatura no [PR](https://github.com/wailsapp/wails/pull/6531) de @taliesin-ai

## Corrigido
- Corrigir diversos problemas das compilações nightly de v3 e limpar as dependências no [PR](https://github.com/wailsapp/wails/pull/6519) de @taliesin-ai
- corrigir a validação de lançamento de v3 e a classificação de fixtures no [PR](https://github.com/wailsapp/wails/pull/6533) de @taliesin-ai
- resolver falhas de validação de lançamento multiplataforma no [PR](https://github.com/wailsapp/wails/pull/6538) de @taliesin-ai

## v3.1.0 - 2026-11-01

## Adicionado
- adicionar backends experimentais sem CGO para macOS e Linux no [PR](https://github.com/wailsapp/wails/pull/6521) de @taliesin-ai
- compilações YAML nativas, empacotamento portátil e assinatura no [PR](https://github.com/wailsapp/wails/pull/6531) de @taliesin-ai

## Alterado
- evitar copiar metadados de desktop descartados no [PR](https://github.com/wailsapp/wails/pull/6391) de @taliesin-ai
- formatar comprimentos de conteúdo sem fmt no [PR](https://github.com/wailsapp/wails/pull/6393) de @taliesin-ai
- selecionar raízes de índice sem coletar candidatos no [PR](https://github.com/wailsapp/wails/pull/6404) de @taliesin-ai
- reutilizar buffers de cópia de arquivos não pesquisáveis no [PR](https://github.com/wailsapp/wails/pull/6408) de @taliesin-ai
- evitar alocar um leitor para detectar prefixos no [PR](https://github.com/wailsapp/wails/pull/6410) de @taliesin-ai
- ignorar leituras de detecção de MIME para extensões conhecidas no [PR](https://github.com/wailsapp/wails/pull/6414) de @taliesin-ai
- impedir o crescimento da fila e a oscilação de redução no [PR](https://github.com/wailsapp/wails/pull/6415) de @taliesin-ai
- evitar alocar contextos de solicitação descartados no [PR](https://github.com/wailsapp/wails/pull/6419) de @taliesin-ai
- ignorar a preparação de logs de depuração desativados no [PR](https://github.com/wailsapp/wails/pull/6420) de @taliesin-ai
- alocar buffers de detecção de MIME sob demanda no [PR](https://github.com/wailsapp/wails/pull/6428) de @taliesin-ai
- reduzir as alocações de formatação de hiperlinks no [PR](https://github.com/wailsapp/wails/pull/6443) de @taliesin-ai
- eliminar alocações temporárias de formatação de marcadores no [PR](https://github.com/wailsapp/wails/pull/6449) de @taliesin-ai
- reduzir a conversão de UUID em string para uma única alocação no [PR](https://github.com/wailsapp/wails/pull/6453) de @taliesin-ai
- reduzir as alocações de formatação de identificadores de instrução no [PR](https://github.com/wailsapp/wails/pull/6454) de @taliesin-ai
- eliminar a alocação temporária do resumo de UUID baseado em nome no [PR](https://github.com/wailsapp/wails/pull/6455) de @taliesin-ai
- evitar copiar linhas incompletas da saída de comandos no [PR](https://github.com/wailsapp/wails/pull/6460) de @taliesin-ai
- transmitir as cópias de executáveis durante a preparação do pacote no [PR](https://github.com/wailsapp/wails/pull/6472) de @taliesin-ai
- evitar examinar novamente linhas de saída incompletas no [PR](https://github.com/wailsapp/wails/pull/6473) de @taliesin-ai
- reutilizar padrões compilados de roteamento de comandos no [PR](https://github.com/wailsapp/wails/pull/6475) de @taliesin-ai
- evitar consultas de metadados durante a descoberta do runtime no [PR](https://github.com/wailsapp/wails/pull/6482) de @taliesin-ai
- reduzir as alocações da página de fallback localizada no [PR](https://github.com/wailsapp/wails/pull/6485) de @taliesin-ai
- evitar recompilar padrões para cada processador no [PR](https://github.com/wailsapp/wails/pull/6490) de @taliesin-ai

## Corrigido
- implementar SetTooltip para a bandeja do sistema no [PR](https://github.com/wailsapp/wails/pull/6017) de @fan711
- ler o modo escuro no portal de configurações no [PR](https://github.com/wailsapp/wails/pull/6072) de @TechHutTV
- permitir que Window.SetZoom fique abaixo de 1.0 no [PR](https://github.com/wailsapp/wails/pull/6166) de @APshenkin
- entregar a cada listener Go na ordem de emissão no [PR](https://github.com/wailsapp/wails/pull/6275) de @DevLumuz
- ler os hooks de eventos sob hookLock em Emit no [PR](https://github.com/wailsapp/wails/pull/6276) de @DevLumuz
- sempre anunciar a propriedade Menu da bandeja para que uma bandeja criada após Run não seja descartada no [PR](https://github.com/wailsapp/wails/pull/6284) de @AlbinoGeek
- reter o trabalho da thread principal enviado antes de Run até que Run seja iniciado no [PR](https://github.com/wailsapp/wails/pull/6287) de @AlbinoGeek
- usar a URL correta do ativo no provedor de atualização do GitHub no [PR](https://github.com/wailsapp/wails/pull/6304) de @gavidroselj
- manter o corpo da solicitação quando o WebView2 omitir Content-Length no [PR](https://github.com/wailsapp/wails/pull/6326) de @strraaa
- relatar as telas do macOS e GTK no espaço DIP (limites com DPI misto) no [PR](https://github.com/wailsapp/wails/pull/6346) de @julianstorer
- não relatar uma notificação dispensada como um clique no [PR](https://github.com/wailsapp/wails/pull/6367) de @bercedev
- interromper os processos de desenvolvimento quando o terminal for fechado no [PR](https://github.com/wailsapp/wails/pull/6371) de @xmasdev
- mapear execuções do registro de alterações para commits de squash merge na verificação de evidências no [PR](https://github.com/wailsapp/wails/pull/6373) de @taliesin-ai
- manter o debounce da reinicialização da barra de tarefas por um segundo no [PR](https://github.com/wailsapp/wails/pull/6386) de @taliesin-ai
- preservar as restrições de tamanho ao entrar em tela cheia quando a janela já começa nesse estado no [PR](https://github.com/wailsapp/wails/pull/6387) de @taliesin-ai
- não redimensionar uma janela em tela cheia ou maximizada para se adequar a uma nova restrição no [PR](https://github.com/wailsapp/wails/pull/6388) de @taliesin-ai
- não exibir uma caixa de diálogo antes de Run no exemplo liquid-glass em [PR](https://github.com/wailsapp/wails/pull/6389) por @taliesin-ai
- uma caixa de diálogo solicitada antes de Run não faz nada em vez de aguardar indefinidamente em [PR](https://github.com/wailsapp/wails/pull/6390) por @taliesin-ai
- usar o fallback quando a entrada de desktop não tiver comando em [PR](https://github.com/wailsapp/wails/pull/6392) por @taliesin-ai
- detectar o conteúdo atual sem reter nomes de arquivos em [PR](https://github.com/wailsapp/wails/pull/6394) por @taliesin-ai
- reconhecer extensões de assets com combinação de maiúsculas e minúsculas em [PR](https://github.com/wailsapp/wails/pull/6395) por @taliesin-ai
- liberar buffers consumidos pelo detector de conteúdo em [PR](https://github.com/wailsapp/wails/pull/6396) por @taliesin-ai
- liberar referências removidas da ringqueue em [PR](https://github.com/wailsapp/wails/pull/6397) por @taliesin-ai
- respeitar exclusões de cabeçalhos antes da seleção do fallback em [PR](https://github.com/wailsapp/wails/pull/6398) por @taliesin-ai
- efetivar os cabeçalhos da resposta de fallback antes do flush em [PR](https://github.com/wailsapp/wails/pull/6399) por @taliesin-ai
- encontrar somente arquivos de índice com nomes exatos em [PR](https://github.com/wailsapp/wails/pull/6400) por @taliesin-ai
- evitar corpos e leituras completas em requisições HEAD não pesquisáveis em [PR](https://github.com/wailsapp/wails/pull/6401) por @taliesin-ai
- preencher o prefixo de detecção de MIME ao longo de leituras curtas em [PR](https://github.com/wailsapp/wails/pull/6402) por @taliesin-ai
- efetivar respostas vazias não pesquisáveis bem-sucedidas em [PR](https://github.com/wailsapp/wails/pull/6403) por @taliesin-ai
- evitar duplicar destinos de URI expandidos em [PR](https://github.com/wailsapp/wails/pull/6405) por @taliesin-ai
- preservar caminhos durante a expansão de campos em [PR](https://github.com/wailsapp/wails/pull/6406) por @taliesin-ai
- respeitar os diretórios de pesquisa de arquivos desktop do XDG em [PR](https://github.com/wailsapp/wails/pull/6407) por @taliesin-ai
- efetivar os cabeçalhos do detector em um flush vazio em [PR](https://github.com/wailsapp/wails/pull/6409) por @taliesin-ai
- ignorar diretórios que ocultam entradas de desktop em [PR](https://github.com/wailsapp/wails/pull/6411) por @taliesin-ai
- preservar destinos encontrados dentro de caminhos de executáveis em [PR](https://github.com/wailsapp/wails/pull/6412) por @taliesin-ai
- liberar capturas de callback após o acionamento do temporizador em [PR](https://github.com/wailsapp/wails/pull/6413) por @taliesin-ai
- evitar correspondências falsas de bibliotecas do ldconfig em [PR](https://github.com/wailsapp/wails/pull/6416) por @taliesin-ai
- restringir as correspondências do fallback de bibliotecas no sistema de arquivos em [PR](https://github.com/wailsapp/wails/pull/6417) por @taliesin-ai
- eliminar a condição de corrida na invalidação do cache em [PR](https://github.com/wailsapp/wails/pull/6418) por @taliesin-ai
- respeitar a configuração opcional do logger em [PR](https://github.com/wailsapp/wails/pull/6421) por @taliesin-ai
- evitar a mutação da URL inicial compartilhada e portas obsoletas em [PR](https://github.com/wailsapp/wails/pull/6422) por @taliesin-ai
- interromper novas tentativas do proxy de desenvolvimento quando a requisição for cancelada em [PR](https://github.com/wailsapp/wails/pull/6423) por @taliesin-ai
- retornar 404 quando faltarem assets incorporados em [PR](https://github.com/wailsapp/wails/pull/6424) por @taliesin-ai
- processar requisições HEAD do runtime incorporado em [PR](https://github.com/wailsapp/wails/pull/6425) por @taliesin-ai
- preservar caminhos com escape ao rotear serviços em [PR](https://github.com/wailsapp/wails/pull/6426) por @taliesin-ai
- ordenar corretamente identificadores numéricos de pré-lançamento excessivamente grandes em [PR](https://github.com/wailsapp/wails/pull/6427) por @taliesin-ai
- rejeitar erros HTTP antes de salvar downloads de ferramentas em [PR](https://github.com/wailsapp/wails/pull/6429) por @taliesin-ai
- respeitar o prazo de inicialização durante a descoberta de desktops no [PR](https://github.com/wailsapp/wails/pull/6430) por @taliesin-ai
- preservar os corpos das requisições durante novas tentativas do proxy de desenvolvimento no [PR](https://github.com/wailsapp/wails/pull/6432) por @taliesin-ai
- fechar os corpos das respostas HTTP em todos os caminhos de retorno no [PR](https://github.com/wailsapp/wails/pull/6434) por @taliesin-ai
- preservar erros de cancelamento em buscas paralelas de bibliotecas no [PR](https://github.com/wailsapp/wails/pull/6435) por @taliesin-ai
- rejeitar respostas malsucedidas de tags do GitHub antes da decodificação no [PR](https://github.com/wailsapp/wails/pull/6436) por @taliesin-ai
- preservar literalmente os caminhos de modelos entre colchetes no [PR](https://github.com/wailsapp/wails/pull/6437) por @taliesin-ai
- instalar o wails3 uma única vez e preservar os erros do instalador no [PR](https://github.com/wailsapp/wails/pull/6438) por @taliesin-ai
- ignorar entradas ocultas e cabeçalhos em blocos de comentários no [PR](https://github.com/wailsapp/wails/pull/6439) por @taliesin-ai
- expor erros de consulta da versão estável no [PR](https://github.com/wailsapp/wails/pull/6440) por @taliesin-ai
- imprimir mensagens simples sem interpretação de formato no [PR](https://github.com/wailsapp/wails/pull/6441) por @taliesin-ai
- evitar panics de escrita após carregar um valor null de JSON no [PR](https://github.com/wailsapp/wails/pull/6442) por @taliesin-ai
- retornar snapshots para leituras do armazenamento inteiro no [PR](https://github.com/wailsapp/wails/pull/6444) por @taliesin-ai
- propagar erros após a iteração parcial das linhas no [PR](https://github.com/wailsapp/wails/pull/6446) por @taliesin-ai
- manter literais os caminhos de destino da extração no [PR](https://github.com/wailsapp/wails/pull/6447) por @taliesin-ai
- liberar identificadores de instruções preparadas ao fechar e reabrir no [PR](https://github.com/wailsapp/wails/pull/6448) por @taliesin-ai
- tratar a condição de corrida entre a preparação e o fechamento ou a reabertura no [PR](https://github.com/wailsapp/wails/pull/6450) por @taliesin-ai
- respeitar requisições HEAD para a página de índice alternativa no [PR](https://github.com/wailsapp/wails/pull/6451) por @taliesin-ai
- exigir identificadores hexadecimais completos de instruções no [PR](https://github.com/wailsapp/wails/pull/6452) por @taliesin-ai
- preservar o banco de dados padrão em memória sob concorrência no [PR](https://github.com/wailsapp/wails/pull/6456) por @taliesin-ai
- separar stderr da saída HEAD analisada no [PR](https://github.com/wailsapp/wails/pull/6457) por @taliesin-ai
- gerar ativos do runtime a partir de bundles incorporados no [PR](https://github.com/wailsapp/wails/pull/6458) por @taliesin-ai
- detectar marcadores TODO, FIXME e TBD no [PR](https://github.com/wailsapp/wails/pull/6459) por @taliesin-ai
- validar respostas da API de notas de versão no [PR](https://github.com/wailsapp/wails/pull/6461) por @taliesin-ai
- respeitar os diretórios das tarefas nas verificações de status no [PR](https://github.com/wailsapp/wails/pull/6462) por @taliesin-ai
- impedir a perda de dados ao copiar um arquivo para ele mesmo no [PR](https://github.com/wailsapp/wails/pull/6463) por @taliesin-ai
- incluir o diretório da tarefa na identidade do cache no [PR](https://github.com/wailsapp/wails/pull/6464) por @taliesin-ai
- sondar as portas uma única vez por verificação no [PR](https://github.com/wailsapp/wails/pull/6465) por @taliesin-ai
- evitar um panic quando entradas do cache forem decodificadas como null no [PR](https://github.com/wailsapp/wails/pull/6466) por @taliesin-ai
- incrementar pré-lançamentos numéricos sem limites de inteiros no [PR](https://github.com/wailsapp/wails/pull/6467) por @taliesin-ai
- impedir que verificações legadas ocultem falhas de cache no [PR](https://github.com/wailsapp/wails/pull/6468) por @taliesin-ai
- executar tarefas padrão no diretório base do executor no [PR](https://github.com/wailsapp/wails/pull/6469) por @taliesin-ai
- preservar sinais de porcentagem literais nos comandos do explorador de arquivos no [PR](https://github.com/wailsapp/wails/pull/6470) por @taliesin-ai
- alinhar os diretórios dos comandos de pré-condição e de tarefa na [PR](https://github.com/wailsapp/wails/pull/6471) por @taliesin-ai
- substituir links simbólicos por arquivos regulares gerados na [PR](https://github.com/wailsapp/wails/pull/6474) por @taliesin-ai
- preservar DMGs existentes quando as compilações falharem na [PR](https://github.com/wailsapp/wails/pull/6476) por @taliesin-ai
- preservar regex entre aspas e barras invertidas em caminhos na [PR](https://github.com/wailsapp/wails/pull/6477) por @taliesin-ai
- validar os tamanhos dos ícones gerados antes da alocação da imagem na [PR](https://github.com/wailsapp/wails/pull/6478) por @taliesin-ai
- encaminhar cadeias de comandos e redirecionamentos pelo shell na [PR](https://github.com/wailsapp/wails/pull/6479) por @taliesin-ai
- preservar pacotes Linux anteriores em caso de falha na compilação na [PR](https://github.com/wailsapp/wails/pull/6480) por @taliesin-ai
- usar execução alternativa de shell específica da plataforma na [PR](https://github.com/wailsapp/wails/pull/6481) por @taliesin-ai
- servir recursos com barras iniciais repetidas no caminho na [PR](https://github.com/wailsapp/wails/pull/6483) por @taliesin-ai
- reconhecer versões estáveis com metadados de compilação na [PR](https://github.com/wailsapp/wails/pull/6484) por @taliesin-ai
- relatar com precisão as falhas na geração de recursos na [PR](https://github.com/wailsapp/wails/pull/6486) por @taliesin-ai
- evitar condições de corrida ao substituir o relator ativo na [PR](https://github.com/wailsapp/wails/pull/6487) por @taliesin-ai
- serializar salvamentos simultâneos de arquivos de cache na [PR](https://github.com/wailsapp/wails/pull/6488) por @taliesin-ai
- oferecer suporte a vários candidatos de igualdade de modelo na [PR](https://github.com/wailsapp/wails/pull/6489) por @taliesin-ai
- normalizar mapas de categorias salvos com valor nulo na [PR](https://github.com/wailsapp/wails/pull/6491) por @taliesin-ai
- mostrar diagnósticos de comandos de shell variáveis que falharam na [PR](https://github.com/wailsapp/wails/pull/6492) por @taliesin-ai
- manter cadeias de comandos de shell fora do cache implícito do Go na [PR](https://github.com/wailsapp/wails/pull/6493) por @taliesin-ai
- serializar a persistência de categorias no Linux na [PR](https://github.com/wailsapp/wails/pull/6494) por @taliesin-ai
- armazenar em cache saídas do Go com espaços entre aspas ou escapados na [PR](https://github.com/wailsapp/wails/pull/6495) por @taliesin-ai
- preservar ícones do Mac quando a codificação falhar na [PR](https://github.com/wailsapp/wails/pull/6496) por @taliesin-ai
- parar de descartar sinalizadores inválidos de contagem de testes do Go na [PR](https://github.com/wailsapp/wails/pull/6497) por @taliesin-ai
- manter as permissões da origem em novas cópias de arquivos na [PR](https://github.com/wailsapp/wails/pull/6498) por @taliesin-ai
- corrigir diversos problemas da compilação nightly v3 e limpar as dependências na [PR](https://github.com/wailsapp/wails/pull/6519) por @taliesin-ai
- corrigir a validação de versão v3 e a classificação de fixtures na [PR](https://github.com/wailsapp/wails/pull/6533) por @taliesin-ai
- resolver falhas de validação de versão entre plataformas na [PR](https://github.com/wailsapp/wails/pull/6538) por @taliesin-ai

## v3.1.0-nightly.20261028 - 2026-10-28

## Alterado
- Melhorar a velocidade de geração do changelog reutilizando padrões regex compilados na [PR](https://github.com/wailsapp/wails/pull/6490) por @taliesin-ai
- Reduzir as alocações ao servir páginas alternativas localizadas na [PR](https://github.com/wailsapp/wails/pull/6485) por @taliesin-ai
- Acelerar a descoberta do runtime do AppImage evitando consultas desnecessárias de metadados na [PR](https://github.com/wailsapp/wails/pull/6482) por @taliesin-ai
- Reutilizar padrões compilados de roteamento de comandos para melhorar o desempenho do wake na [PR](https://github.com/wailsapp/wails/pull/6475) por @taliesin-ai
- Melhorar o desempenho evitando verificar novamente linhas de saída incompletas na [PR](https://github.com/wailsapp/wails/pull/6473) por @taliesin-ai
- Transmitir as cópias de executáveis durante a preparação do pacote MSIX para reduzir o uso de memória e acelerar as compilações na [PR](https://github.com/wailsapp/wails/pull/6472) por @taliesin-ai
- Evitar copiar linhas incompletas da saída de comandos para melhorar o desempenho na [PR](https://github.com/wailsapp/wails/pull/6460) por @taliesin-ai
- Reduzir as alocações ao gerar UUIDs baseados em nomes na [PR](https://github.com/wailsapp/wails/pull/6455) por @taliesin-ai
- Reduzir as alocações ao formatar identificadores de instruções SQLite na [PR](https://github.com/wailsapp/wails/pull/6454) por @taliesin-ai
- Reduzir a conversão de UUID em string a uma única alocação na [PR](https://github.com/wailsapp/wails/pull/6453) por @taliesin-ai
- Reduzir as alocações ao formatar itens com marcadores do changelog na [PR](https://github.com/wailsapp/wails/pull/6449) por @taliesin-ai
- Reduzir as alocações ao formatar hiperlinks de terminal na [PR](https://github.com/wailsapp/wails/pull/6443) por @taliesin-ai
- Reduzir as alocações ao servir recursos inspecionando os buffers MIME de forma tardia na [PR](https://github.com/wailsapp/wails/pull/6428) por @taliesin-ai
- Ignorar a preparação de logs de depuração no servidor de recursos quando o registro de depuração estiver desativado na [PR](https://github.com/wailsapp/wails/pull/6420) por @taliesin-ai
- Reduzir as alocações no processamento de solicitações da webview evitando contextos de solicitação descartados na [PR](https://github.com/wailsapp/wails/pull/6419) por @taliesin-ai
- Acelerar a entrega de recursos ignorando a detecção de MIME para extensões de arquivo conhecidas na [PR](https://github.com/wailsapp/wails/pull/6414) por @taliesin-ai
- Reduzir as alocações ao servir recursos evitando usar um leitor para detectar o conteúdo na [PR](https://github.com/wailsapp/wails/pull/6410) por @taliesin-ai
- Reduzir as alocações ao servir arquivos de recursos que não permitem posicionamento na [PR](https://github.com/wailsapp/wails/pull/6408) por @taliesin-ai
- Acelerar a seleção da raiz do índice do servidor de recursos evitando coletar candidatos na [PR](https://github.com/wailsapp/wails/pull/6404) por @taliesin-ai
- Melhorar o desempenho do servidor de recursos formatando os comprimentos de conteúdo sem fmt na [PR](https://github.com/wailsapp/wails/pull/6393) por @taliesin-ai
- Reduzir a cópia de memória ao ler metadados do explorador de arquivos da área de trabalho na [PR](https://github.com/wailsapp/wails/pull/6391) por @taliesin-ai

## Corrigido
- Manter as permissões do arquivo de origem ao copiá-lo para novos arquivos de destino na [PR](https://github.com/wailsapp/wails/pull/6498) por @taliesin-ai
- Impedir que janelas em tela cheia ou maximizadas sejam redimensionadas quando as restrições de tamanho mínimo ou máximo mudarem na [PR](https://github.com/wailsapp/wails/pull/6388) por @taliesin-ai
- Publicar atualizações programadas do histórico de estrelas do repositório wailsapp/wails na [PR](https://github.com/wailsapp/wails/pull/6500) por @taliesin-ai
- Parar de descartar sinalizadores inválidos de contagem de testes do Go no wake na [PR](https://github.com/wailsapp/wails/pull/6497) por @taliesin-ai
- Preservar ícones do Mac quando a codificação falhar na [PR](https://github.com/wailsapp/wails/pull/6496) por @taliesin-ai
- Corrigir falhas de cache de compilação para saídas do Go com espaços entre aspas ou escapados na [PR](https://github.com/wailsapp/wails/pull/6495) por @taliesin-ai
- Serializar a persistência de categorias de notificação no Linux para evitar condições de corrida em gravações simultâneas na [PR](https://github.com/wailsapp/wails/pull/6494) por @taliesin-ai
- Manter cadeias de comandos de shell fora do cache implícito de compilação do Go na [PR](https://github.com/wailsapp/wails/pull/6493) por @taliesin-ai
- Corrigir a falha no carregamento de mapas de categorias de notificação quando os dados salvos forem nulos na [PR](https://github.com/wailsapp/wails/pull/6491) por @taliesin-ai
- Mostrar diagnósticos de comandos de shell variáveis que falharam na [PR](https://github.com/wailsapp/wails/pull/6492) por @taliesin-ai
- Corrigir a correspondência de modelos do wake quando houver vários candidatos de igualdade na [PR](https://github.com/wailsapp/wails/pull/6489) por @taliesin-ai
- Serializar salvamentos simultâneos de arquivos de cache para evitar gravações corrompidas no cache no [PR](https://github.com/wailsapp/wails/pull/6488) por @taliesin-ai
- Evitar condições de corrida ao substituir o relator ativo no [PR](https://github.com/wailsapp/wails/pull/6487) por @taliesin-ai
- Relatar com precisão as falhas de geração de recursos na CLI no [PR](https://github.com/wailsapp/wails/pull/6486) por @taliesin-ai
- Reconhecer na CLI versões estáveis que incluem metadados de compilação no [PR](https://github.com/wailsapp/wails/pull/6484) por @taliesin-ai
- Servir ativos cujos caminhos de solicitação contêm barras iniciais repetidas no [PR](https://github.com/wailsapp/wails/pull/6483) por @taliesin-ai
- Corrigir a execução alternativa via shell do wake em cada plataforma no [PR](https://github.com/wailsapp/wails/pull/6481) por @taliesin-ai
- Preservar pacotes Linux anteriores quando uma compilação falhar no [PR](https://github.com/wailsapp/wails/pull/6480) por @taliesin-ai
- Encaminhar cadeias de comandos e redirecionamentos pelo shell no [PR](https://github.com/wailsapp/wails/pull/6479) por @taliesin-ai
- Preservar expressões regulares entre aspas e barras invertidas de caminhos ao analisar comandos do wake no [PR](https://github.com/wailsapp/wails/pull/6477) por @taliesin-ai
- Preservar arquivos DMG existentes quando o empacotamento de compilações falhar no [PR](https://github.com/wailsapp/wails/pull/6476) por @taliesin-ai
- Substituir links simbólicos por arquivos regulares gerados na saída de bindings no [PR](https://github.com/wailsapp/wails/pull/6474) por @taliesin-ai
- Preservar sinais de porcentagem literais nos comandos do explorador no Linux no [PR](https://github.com/wailsapp/wails/pull/6470) por @taliesin-ai
- Alinhar os diretórios dos comandos de pré-condição e de tarefa para que sejam resolvidos de modo consistente no [PR](https://github.com/wailsapp/wails/pull/6471) por @taliesin-ai
- Executar tarefas padrão no diretório base do executor no [PR](https://github.com/wailsapp/wails/pull/6469) por @taliesin-ai
- Impedir que verificações legadas ocultem ausências no cache no [PR](https://github.com/wailsapp/wails/pull/6468) por @taliesin-ai
- Corrigir a numeração de pré-lançamentos na CLI para incrementá-la sem limites de estouro de inteiro no [PR](https://github.com/wailsapp/wails/pull/6467) por @taliesin-ai
- Evitar um panic quando entradas do cache forem decodificadas como null no [PR](https://github.com/wailsapp/wails/pull/6466) por @taliesin-ai
- Sondar as portas uma vez por verificação na CLI no [PR](https://github.com/wailsapp/wails/pull/6465) por @taliesin-ai
- Incluir o diretório da tarefa na identidade do cache para evitar resultados obsoletos do wake no [PR](https://github.com/wailsapp/wails/pull/6464) por @taliesin-ai
- Evitar perda de dados quando o empacotamento copiar um arquivo sobre ele mesmo no [PR](https://github.com/wailsapp/wails/pull/6463) por @taliesin-ai
- Fazer com que as verificações de status respeitem os diretórios das tarefas no [PR](https://github.com/wailsapp/wails/pull/6462) por @taliesin-ai
- Detectar marcadores TODO, FIXME e TBD nas entradas do changelog no [PR](https://github.com/wailsapp/wails/pull/6459) por @taliesin-ai
- Gerar ativos de runtime a partir de pacotes incorporados na CLI no [PR](https://github.com/wailsapp/wails/pull/6458) por @taliesin-ai
- Separar stderr da saída HEAD analisada nos comandos git no [PR](https://github.com/wailsapp/wails/pull/6457) por @taliesin-ai
- Preservar o banco de dados SQLite padrão em memória durante acessos simultâneos no [PR](https://github.com/wailsapp/wails/pull/6456) por @taliesin-ai
- Retornar respostas corretas às solicitações HEAD da página de índice alternativa no [PR](https://github.com/wailsapp/wails/pull/6451) por @taliesin-ai
- Tratar a condição de corrida entre a preparação de instruções SQLite e o fechamento ou a reabertura do banco de dados no [PR](https://github.com/wailsapp/wails/pull/6450) por @taliesin-ai
- Liberar identificadores de instruções preparadas quando o banco de dados SQLite for fechado ou reaberto no [PR](https://github.com/wailsapp/wails/pull/6448) por @taliesin-ai
- Manter literais os caminhos de destino da extração ao usar gosod no [PR](https://github.com/wailsapp/wails/pull/6447) por @taliesin-ai
- Propagar erros de consultas SQLite após a iteração parcial das linhas no [PR](https://github.com/wailsapp/wails/pull/6446) por @taliesin-ai
- Retornar snapshots para leituras de todo o armazenamento no kvstore no [PR](https://github.com/wailsapp/wails/pull/6444) por @taliesin-ai
- Evitar panics de gravação no kvstore após carregar um valor JSON null no [PR](https://github.com/wailsapp/wails/pull/6442) por @taliesin-ai
- Imprimir mensagens simples no terminal sem interpretar verbos de formatação no [PR](https://github.com/wailsapp/wails/pull/6441) por @taliesin-ai
- Expor erros ao consultar a versão estável mais recente no atualizador no [PR](https://github.com/wailsapp/wails/pull/6440) por @taliesin-ai
- Ignorar entradas ocultas e cabeçalhos nos blocos de comentários do changelog no [PR](https://github.com/wailsapp/wails/pull/6439) por @taliesin-ai
- Instalar o wails3 uma única vez e expor erros do instalador no atualizador no [PR](https://github.com/wailsapp/wails/pull/6438) por @taliesin-ai
- Preservar caminhos literais de modelos entre colchetes no gosod no [PR](https://github.com/wailsapp/wails/pull/6437) por @taliesin-ai
- Rejeitar respostas de tags do GitHub com falha antes de decodificá-las no atualizador no [PR](https://github.com/wailsapp/wails/pull/6436) por @taliesin-ai
- Preservar erros de cancelamento durante buscas paralelas de bibliotecas no [PR](https://github.com/wailsapp/wails/pull/6435) por @taliesin-ai
- Fechar os corpos das respostas HTTP em todos os caminhos de retorno para evitar vazamentos de recursos no [PR](https://github.com/wailsapp/wails/pull/6434) por @taliesin-ai
- Preservar os corpos das solicitações quando o proxy de desenvolvimento tentar enviá-las novamente no [PR](https://github.com/wailsapp/wails/pull/6432) por @taliesin-ai
- Respeitar o prazo de inicialização durante a descoberta do explorador de arquivos do desktop no [PR](https://github.com/wailsapp/wails/pull/6430) por @taliesin-ai
- Rejeitar erros HTTP antes de salvar downloads das ferramentas do AppImage no [PR](https://github.com/wailsapp/wails/pull/6429) por @taliesin-ai
- Ordenar corretamente identificadores numéricos de pré-lançamento muito grandes nas comparações semver no [PR](https://github.com/wailsapp/wails/pull/6427) por @taliesin-ai
- Preservar caminhos escapados ao rotear serviços no servidor de ativos no [PR](https://github.com/wailsapp/wails/pull/6426) por @taliesin-ai
- Tratar solicitações HEAD do runtime incluído no servidor de ativos no [PR](https://github.com/wailsapp/wails/pull/6425) por @taliesin-ai
- Retornar 404 para ativos incluídos ausentes em vez de uma resposta incorreta no [PR](https://github.com/wailsapp/wails/pull/6424) por @taliesin-ai
- Interromper novas tentativas do proxy de desenvolvimento quando uma solicitação for cancelada no [PR](https://github.com/wailsapp/wails/pull/6423) por @taliesin-ai
- Evitar a mutação compartilhada da URL inicial e portas obsoletas no servidor de ativos no [PR](https://github.com/wailsapp/wails/pull/6422) por @taliesin-ai
- Respeitar a configuração opcional do logger no servidor de ativos no [PR](https://github.com/wailsapp/wails/pull/6421) por @taliesin-ai
- Eliminar a condição de corrida na invalidação do cache de caminhos de bibliotecas no [PR](https://github.com/wailsapp/wails/pull/6418) por @taliesin-ai
- Evitar correspondências falsas de bibliotecas ao pesquisar a saída do ldconfig no [PR](https://github.com/wailsapp/wails/pull/6416) por @taliesin-ai
- Evitar o crescimento e a oscilação de redução da fila do servidor de ativos sob carga no [PR](https://github.com/wailsapp/wails/pull/6415) por @taliesin-ai
- Liberar as capturas do callback após o disparo do temporizador de debounce para evitar a retenção de referências obsoletas no [PR](https://github.com/wailsapp/wails/pull/6413) por @taliesin-ai
- Preservar alvos do explorador de arquivos encontrados em caminhos de executáveis no [PR](https://github.com/wailsapp/wails/pull/6412) por @taliesin-ai
- Ignorar diretórios que ocultam entradas de desktop no explorador de arquivos no [PR](https://github.com/wailsapp/wails/pull/6411) por @taliesin-ai
- Corrigir o não envio dos cabeçalhos do servidor de ativos quando o corpo da resposta está vazio no [PR](https://github.com/wailsapp/wails/pull/6409) por @taliesin-ai
- Respeitar os diretórios XDG de pesquisa de arquivos desktop ao localizar ícones de aplicativos no [PR](https://github.com/wailsapp/wails/pull/6407) por @taliesin-ai
- Preservar os caminhos de arquivo durante a expansão de campos no explorador de arquivos em [PR](https://github.com/wailsapp/wails/pull/6406) por @taliesin-ai
- Evitar a duplicação de destinos de URI expandidos no explorador de arquivos em [PR](https://github.com/wailsapp/wails/pull/6405) por @taliesin-ai
- Confirmar respostas de ativos vazias, bem-sucedidas e não pesquisáveis, em vez de deixá-las pendentes, em [PR](https://github.com/wailsapp/wails/pull/6403) por @taliesin-ai
- Preencher corretamente o prefixo de detecção de MIME ao longo de leituras curtas no servidor de ativos em [PR](https://github.com/wailsapp/wails/pull/6402) por @taliesin-ai
- Evitar a leitura completa dos corpos de resposta em solicitações HEAD não pesquisáveis no servidor de ativos em [PR](https://github.com/wailsapp/wails/pull/6401) por @taliesin-ai
- Servir somente arquivos de índice com correspondência exata ao resolver solicitações de diretórios no servidor de ativos em [PR](https://github.com/wailsapp/wails/pull/6400) por @taliesin-ai
- Corrigir notificações do Linux descartadas que eram relatadas como cliques em servidores sem suporte a ações em [PR](https://github.com/wailsapp/wails/pull/6367) por @bercedev
- Encerrar os processos de desenvolvimento quando o terminal for fechado em [PR](https://github.com/wailsapp/wails/pull/6371) por @xmasdev
- Manter por um segundo o debounce da reinicialização da barra de tarefas do Windows antes de limpar o sinalizador em [PR](https://github.com/wailsapp/wails/pull/6386) por @taliesin-ai
- Corrigir a perda de cabeçalhos de respostas de fallback no servidor de ativos em [PR](https://github.com/wailsapp/wails/pull/6399) por @taliesin-ai
- Respeitar as exclusões de cabeçalhos antes da seleção do fallback no servidor de ativos em [PR](https://github.com/wailsapp/wails/pull/6398) por @taliesin-ai
- Liberar referências removidas da fila circular no servidor de ativos para evitar vazamentos de memória em [PR](https://github.com/wailsapp/wails/pull/6397) por @taliesin-ai
- Liberar os buffers consumidos pelo detector de conteúdo no servidor de ativos para reduzir o uso de memória em [PR](https://github.com/wailsapp/wails/pull/6396) por @taliesin-ai
- Reconhecer extensões de arquivos de ativos com combinação de maiúsculas e minúsculas ao servir ativos em [PR](https://github.com/wailsapp/wails/pull/6395) por @taliesin-ai
- Detectar corretamente o tipo de conteúdo no servidor de ativos sem reter nomes de arquivos em [PR](https://github.com/wailsapp/wails/pull/6394) por @taliesin-ai
- Recorrer a um valor padrão quando uma entrada de desktop não tiver comando em [PR](https://github.com/wailsapp/wails/pull/6392) por @taliesin-ai
- Impedir que o exemplo liquid-glass exiba uma caixa de diálogo antes de o aplicativo iniciar em sistemas que não sejam macOS em [PR](https://github.com/wailsapp/wails/pull/6389) por @taliesin-ai
- Manter as restrições de tamanho da janela quando o modo de tela cheia for iniciado a partir de um tamanho restrito no Windows em [PR](https://github.com/wailsapp/wails/pull/6387) por @taliesin-ai
- Corrigir caixas de diálogo chamadas antes de Run que bloqueavam indefinidamente em vez de retornar em [PR](https://github.com/wailsapp/wails/pull/6390) por @taliesin-ai

## Segurança
- Validar os tamanhos dos ícones gerados antes da alocação de imagens para evitar o uso excessivo de memória em [PR](https://github.com/wailsapp/wails/pull/6478) por @taliesin-ai
- Validar as respostas da API de notas de versão antes de usá-las em [PR](https://github.com/wailsapp/wails/pull/6461) por @taliesin-ai
- Exigir identificadores hexadecimais completos de instruções no plugin SQLite em [PR](https://github.com/wailsapp/wails/pull/6452) por @taliesin-ai
- Restringir as correspondências de fallback de bibliotecas no sistema de arquivos para evitar o carregamento de bibliotecas não pretendidas em [PR](https://github.com/wailsapp/wails/pull/6417) por @taliesin-ai

## v3.1.0-nightly.20261027 - 2026-10-27

## Alterado
- Vincular a página 404 às páginas correspondentes da documentação do Wails v2 em v2.wails.io em [PR](https://github.com/wailsapp/wails/pull/6366) por @taliesin-ai
- Mover o site de documentação v2 para v2.wails.io em [PR](https://github.com/wailsapp/wails/pull/6365) por @taliesin-ai
- Servir a documentação v3 em wails.io e a documentação v2 em v2.wails.io em [PR](https://github.com/wailsapp/wails/pull/6364) por @taliesin-ai
- Atualizar as notas da versão v3.0.0 e corrigir a alegação do aviso de segurança em [PR](https://github.com/wailsapp/wails/pull/6363) por @taliesin-ai

## Corrigido
- Entregar os callbacks dos listeners de eventos do Go em série, na ordem de emissão, em [PR](https://github.com/wailsapp/wails/pull/6275) por @DevLumuz
- Sempre anunciar a propriedade Menu da bandeja para que as bandejas criadas após Run não sejam descartadas no Linux em [PR](https://github.com/wailsapp/wails/pull/6284) por @AlbinoGeek
- Reter o trabalho da thread principal enviado antes de Run até que o aplicativo seja iniciado em [PR](https://github.com/wailsapp/wails/pull/6287) por @AlbinoGeek
- Usar a URL correta da API de ativos ao baixar ativos de versões do atualizador pelo GitHub em [PR](https://github.com/wailsapp/wails/pull/6304) por @gavidroselj
- Corrigir os limites e as posições das telas em configurações de macOS e Linux com DPIs diferentes em [PR](https://github.com/wailsapp/wails/pull/6346) por @julianstorer
- Preservar os corpos das solicitações no Windows quando o WebView2 omitir o cabeçalho Content-Length em [PR](https://github.com/wailsapp/wails/pull/6326) por @strraaa
- Corrigir a condição de corrida de dados na emissão de hooks de eventos criando um snapshot dos hooks sob bloqueio em [PR](https://github.com/wailsapp/wails/pull/6276) por @DevLumuz
- Corrigir `SystemTray.SetTooltip()` que não fazia nada no Linux: a descrição da dica de ferramenta de StatusNotifierItem agora é definida (#6017)
- Executar a verificação de evidências do changelog no fluxo de trabalho de lançamentos semanais em [PR](https://github.com/wailsapp/wails/pull/6374) por @taliesin-ai
- Corrigir a verificação de evidências do changelog de versão para mapear execuções para commits de squash merge em [PR](https://github.com/wailsapp/wails/pull/6373) por @taliesin-ai
- Corrigir a detecção do modo escuro no Linux, que retornava o modo claro em um desktop escuro até que o esquema de cores fosse alterado, e `linux:SystemThemeChanged`, que nunca era acionado no backend GTK4. Agora, ambos os backends compartilham uma única implementação, baseada no namespace do portal `org.freedesktop.appearance` (#6072)
- `Window.SetZoom` não eleva mais fatores abaixo de 1.0 para 1.0 em janelas com moldura (Windows, macOS e Linux): WebView2 `ZoomFactor`, WKWebView `magnification` e o nível de zoom do WebKitGTK aceitam qualquer valor positivo, e a opção de janela `Zoom` já repassava esses valores. Janelas sem moldura mantêm o limite mínimo de 1.0 porque suas regiões de arrastar e redimensionar não são ajustadas ao zoom ([#4590](https://github.com/wailsapp/wails/issues/4590)); fatores não positivos, NaN e Inf agora são ignorados em vez de elevados para 1.0 — por @APshenkin

## v3.0.0 - 2026-10-26

## Alterado
- As pré-versões nightly são publicadas como `vX.Y.0-nightly.YYYYMMDD`, antecipando a próxima versão semanal, e as versões semanais incrementam a versão secundária (`vX.Y.0`), de modo que `@latest` do Go sempre resolva para uma versão semanal ou de correção emergencial; a versão do pacote npm `@wailsio/runtime` sempre corresponde à versão do Wails, com as nightlies na dist-tag `nightly`
- Os modelos de projeto fixam `@wailsio/runtime` na versão da CLI que os gerou
- O Wails v3 requer Go 1.26 ou mais recente (Go 1.26.10 ou 1.27.3 e posteriores incluem a correção de segurança `crypto/x509`)

## Corrigido
- Corrigir `wails3 update cli` que não encontrava versões v3: ele considerava apenas tags v2 e instalava a partir do caminho errado
- Corrigir aplicativos do Windows que falhavam com `panic: SetTextColor failed` ao desenhar a barra de menus do aplicativo em algumas builds do Windows: a barra de menus é desenhada em um contexto de dispositivo válido ou recorre ao desenho padrão, `UAHMENUITEM.Position` agora corresponde ao `int` em C, e o desenho do menu nunca causa panic (#6311)
- Manter a atualização preparada em segurança enquanto uma reinicialização estiver pendente: um download ou uma segunda reinicialização iniciada durante o encerramento não pode excluir nem disputar o payload que o assistente de troca aguarda para instalar por @julianstorer
- Corrigir aplicativos do Windows que falhavam em `applyDarkMode` quando a inicialização de `pkg/w32` descarregava `uxtheme.dll` enquanto seus ponteiros de função por ordinal ainda estavam em uso: agora o módulo permanece carregado durante toda a vida útil do processo (#6356)
- Valide o payload preparado antes da substituição e o aplicativo instalado antes de descartar o backup, para que uma atualização danificada seja revertida em vez de substituir um aplicativo funcional, por @julianstorer

## Removido
- Remova a integração com o Crowdin após a sincronização final das traduções no [PR](https://github.com/wailsapp/wails/pull/5736), por @github-actions[bot]

## Segurança
- Bloqueie mensagens reservadas `wails:` provenientes de frames filhos de outra origem no Linux (GTK3 e GTK4), injetando um script de proteção em todos os frames, para mitigar GHSA-2p5f-2gx8-mh8m
- Atualize `github.com/klauspost/compress` para 1.18.7 (GHSA-259r-337f-4rfw) e `golang.org/x/image` para 0.45.0, e compile com Go 1.26.10 para incluir a correção de `crypto/x509` (GO-2026-6809)
- Atualize `vite` para 8.2.4 ou posterior no pacote de runtime, no assistente de configuração, nos exemplos e nos modelos de projeto (GHSA-9jrq-w75r-8gcw, GHSA-rq7h-c2jc-7f22, GHSA-vfpm-58rq-9qcg, GHSA-r825-w8mm-5pxw, GHSA-587h-7wmj-h5pm, GHSA-698v-4r3x-3vcv)

### Limitações conhecidas
- **Linux:** a proteção de frames filhos não consegue impedir que um frame filho de *mesma origem* (por exemplo, um iframe `srcdoc` ou `about:blank`) acesse o manipulador nativo de mensagens por meio do frame pai, pois o WebKitGTK não fornece ao manipulador informações sobre o frame. Frames de outra origem e frames em sandbox sem `allow-same-origin` são bloqueados. Não incorpore conteúdo não confiável de mesma origem em uma janela do Wails.
- **Android (experimental):** a ponte JavaScript `window.wails` fica exposta a todos os frames da WebView, sem verificação do frame principal nem da origem. O conteúdo de um iframe incorporado, inclusive conteúdo de outra origem, pode chamar métodos vinculados e enviar mensagens reservadas `wails:`. Não incorpore conteúdo não confiável de terceiros na WebView do aplicativo. Esse comportamento difere do macOS e do iOS, que ignoram mensagens `wails:` de frames filhos. Uma correção está planejada para 3.0.1.
- **Todas as plataformas:** chamadas de métodos vinculados provenientes de frames filhos trafegam por HTTP até `/wails/runtime`, e não pelo manipulador nativo de mensagens; portanto, a proteção não as restringe. Considere que qualquer iframe incorporado tem acesso aos seus métodos Go vinculados.

## v3.0.0-beta.34 - 2026-10-20

## Corrigido
- Corrija a WebView do Windows que não recuperava o foco ao ativar a janela e a tecla Tab que ficava presa dentro da WebView em vez de mover o foco para a janela nativa e de volta (#3783)
- Corrija as structs de vtable COM de `internal/webview2/pkg/webview2` (`ICoreWebView2_2`..`_27` e outras interfaces derivadas), que omitiam a cadeia de métodos herdados, fazendo com que as chamadas chegassem ao slot errado da vtable (#6129)
- Corrija as caixas de diálogo de pergunta do Windows que ignoravam rótulos personalizados de botões e nunca acionavam os callbacks dos botões; agora, os botões personalizados são exibidos com uma TaskDialog, recorrendo a uma MessageBox com mapeamento posicional de callbacks quando comctl32 v6 não está disponível (#5794)
- Corrija o travamento de `wails3 dev` por 60s ao usar o Vite: a verificação automática de prontidão do frontend agora usa uma sondagem TCP em vez de uma solicitação HTTP GET que não enviava o cabeçalho `Accept` e recebia um 404 ([#6209](https://github.com/wailsapp/wails/issues/6209))
- Mantenha os menus da área de notificação funcionando quando forem atualizados enquanto estiverem abertos no [PR](https://github.com/wailsapp/wails/pull/6301), por @taliesin-ai
- Trate de forma adequada os erros COM transitórios do WebView2 nas operações de runtime no [PR](https://github.com/wailsapp/wails/pull/6302), por @leaanthony

## v3.0.0-beta.33 - 2026-10-19

## Adicionado
- O gerador de bindings agora emite um módulo `bindings/events` que exporta uma constante de nome tipada para cada evento detectável registrado com `application.RegisterEvent` usando um nome constante em tempo de compilação (por exemplo, `UserLogin` para `user:login`), utilizável com `Events.On`/`Emit` sem o plugin do Vite (#5645)

## Corrigido
- Corrija as caixas de diálogo de pergunta do Windows que ignoravam rótulos personalizados de botões e nunca acionavam os callbacks dos botões; agora, os botões personalizados são exibidos com uma TaskDialog, recorrendo a uma MessageBox com mapeamento posicional de callbacks quando comctl32 v6 não está disponível (#5794)
- Corrija o travamento de `wails3 dev` por 60s ao usar o Vite: a verificação automática de prontidão do frontend agora usa uma sondagem TCP em vez de uma solicitação HTTP GET que não enviava o cabeçalho `Accept` e recebia um 404 ([#6209](https://github.com/wailsapp/wails/issues/6209))
- Mantenha os menus da área de notificação funcionando quando forem atualizados enquanto estiverem abertos no [PR](https://github.com/wailsapp/wails/pull/6301), por @taliesin-ai
- Trate de forma adequada os erros COM transitórios do WebView2 nas operações de runtime no [PR](https://github.com/wailsapp/wails/pull/6302), por @leaanthony

## v3.0.0-beta.32 - 2026-10-18

## Adicionado
- Adicione cabeçalhos nativos de seção de menu no macOS 14+ — consulte o [site da documentação](https://v3.wails.io) (20 páginas atualizadas) no [PR](https://github.com/wailsapp/wails/pull/6297), por @leaanthony

## Corrigido
- Faça SetRelativePosition funcionar em relação à área de trabalho da tela no macOS no [PR](https://github.com/wailsapp/wails/pull/6299), por @leaanthony
- O auxiliar do atualizador reinicia o destino inalterado após o tempo limite do processo pai ou uma falha na substituição no [PR](https://github.com/wailsapp/wails/pull/6298), por @taliesin-ai
- Fixe o pnpm na versão 10 para a documentação versionada no [PR](https://github.com/wailsapp/wails/pull/6293), por @taliesin-ai

## v3.0.0-beta.31 - 2026-10-18

## Adicionado
- Exponha a exportação do subcaminho `@wailsio/runtime/cancellable` para que `CancellablePromise` possa ser importado sem avaliar o ponto de entrada raiz do runtime (#6153)

## Alterado
- Exija Go 1.26 ou posterior: o módulo v3 agora declara `go 1.26.0`, e os projetos gerados por `wails3 init` e `wails3 service init` declaram o mesmo requisito mínimo

## Corrigido
- Corrija a alteração da saída do gerador de bindings ao compilar com Go 1.27, em que tipos `encoding/json` como `Marshaler`, `Unmarshaler` e `RawMessage` passaram a ser aliases dos tipos `encoding/json/v2` e `encoding/json/jsontext`: eles voltam a ser renderizados como `any`, e nenhum modelo é emitido para `encoding/json/v2` (#6292)
- Valide os URLs de modelos remotos passados para `wails3 init -t` antes da clonagem: somente URLs `https://`, `http://`, `ssh://`, `git://` e `user@host:path` são aceitos; os valores são passados ao git após `--`, e o git fica restrito a esses transportes durante a clonagem
- Rejeite valores do cabeçalho `x-wails-window-id` fora do intervalo no transporte HTTP do runtime, em vez de truncá-los silenciosamente para o ID de outra janela
- Proteja a extração de arquivos compactados do atualizador contra zip-slip: agora, os caminhos das entradas são verificados com `filepath.IsLocal` e com uma verificação de contenção no caminho resultante da junção, incluindo testes para entradas com `..` aninhado, barra invertida, caminho absoluto e gravação por meio de link simbólico em arquivos zip e tar.gz
- Corrija AppImages de `wails3 generate appimage` que falhavam ao iniciar fora do Debian e do Ubuntu com "Unable to spawn a new child process: … WebKitNetworkProcess". O WebKit incluído agora executa os processos auxiliares fornecidos na AppImage. O sandbox bubblewrap do WebKit permanece ativado
- Corrija as compilações de `-tags server` que não eram compiladas para Android e iOS, corrigindo [#6271](https://github.com/wailsapp/wails/issues/6271)
- Corrija a falha do atualizador ao substituir arquivos no Linux quando o diretório temporário está em outro sistema de arquivos (como um tmpfs `/tmp`), preparando-os ao lado do binário instalado e copiando-os entre os sistemas de arquivos
- Corrija o atualizador que tentava substituir a montagem somente leitura de uma AppImage em vez do arquivo `.AppImage`
- Corrige a aceitação, em iframes no macOS e iOS, das mensagens reservadas de controle de janela `wails:`. Agora elas são aceitas somente no frame principal da janela; outras mensagens provenientes de iframes continuam chegando a `RawMessageHandler` com suas informações de origem
- Interrompe o envio de `WindowDidUpdate` e `WindowWillUpdate` no macOS: `NSWindowDidUpdateNotification` é disparado a cada passagem pelo loop de execução, portanto a página recebia cerca de 120 chamadas `evaluateJavaScript:` vazias por segundo e, no macOS, cada 27.0 consumia a ativação transitória do clique do usuário — `navigator.clipboard.writeText()` e outras APIs condicionadas à ativação rejeitavam cliques reais com `NotAllowedError` ([#6170](https://github.com/wailsapp/wails/issues/6170)), por @APshenkin
- Corrige o idioma padrão do info.json do Windows para en-US (0409), permitindo que o Explorer exiba os detalhes da versão (#6210)
- Corrige a desativação do renderizador DMA-BUF do Linux em notebooks com uma dGPU NVIDIA junto a uma iGPU AMD ou Intel: no Wayland, agora ele só é desativado quando o WebKit renderiza na NVIDIA (quando o aplicativo foi transferido para a dGPU ou o compositor é executado nela), e não sempre que o módulo nvidia está carregado, o que tornava a renderização extremamente lenta na iGPU; sessões X11 mantêm o comportamento anterior, pois a GPU de renderização do Xorg é definida por sua configuração
- Trata falhas de renomeação entre dispositivos no atualizador do Linux usando uma cópia intermediária como alternativa, no [PR](https://github.com/wailsapp/wails/pull/6135), por @john-okeefe
- Corrige `WebviewWindow.SetMenu` no Windows, que não redesenhava a barra de menus a menos que um tema de barra de menus estivesse configurado, corrigindo [#6103](https://github.com/wailsapp/wails/issues/6103)
- Corrige `WebviewWindow.SetMenu(nil)` no Windows, que causava panic em vez de não realizar nenhuma operação, como no macOS, corrigindo [#6104](https://github.com/wailsapp/wails/issues/6104)
- Preserva entradas do changelog diante de pushes concorrentes usando novas tentativas, no [PR](https://github.com/wailsapp/wails/pull/6094), por @leaanthony
- Corrige a falha de `go mod vendor` com `pattern arm64/WebView2Loader.dll: no matching files found` em todas as plataformas, removendo os embeds que faziam referência a binários nunca distribuídos no módulo, corrigindo [#5782](https://github.com/wailsapp/wails/issues/5782) e [#5376](https://github.com/wailsapp/wails/issues/5376), no [PR](https://github.com/wailsapp/wails/pull/6031), por @Grantmartin2002
- Impede a desmontagem duplicada de janelas no macOS, no [PR](https://github.com/wailsapp/wails/pull/6159), por @charlesbalderson
- Trata WM_CONTEXTMENU para abrir o menu da bandeja do sistema com o botão direito, no [PR](https://github.com/wailsapp/wails/pull/6162), por @wzk2025
- Corrige as cinco variantes de `InvokeSync*` que bloqueavam para sempre seu chamador quando o callback despachado causava panic, corrigindo [#6107](https://github.com/wailsapp/wails/issues/6107), no [PR](https://github.com/wailsapp/wails/pull/6108), por @RALIST
- Reconstrói com segurança submenus GTK3 destruídos, no [PR](https://github.com/wailsapp/wails/pull/6236), por @leaanthony
- Implementa no Windows uma camada de encaminhamento para workers aninhados, no [PR](https://github.com/wailsapp/wails/pull/6261), por @taliesin-ai
- Corrige `OpenFile()` no Windows, que ignorava `PromptForMultipleSelection()` quando `CanChooseDirectories(true)`: agora os usuários podem selecionar várias pastas de uma só vez, como no macOS e no Linux (#6263)
- Corrige argumentos `double` passados por valor para métodos COM do WebView2 em windows/arm64 — `ZoomFactor`, `RasterizationScale`, `Expires` de cookies, bem como a escala, o tamanho da página e as margens das configurações de impressão eram afetados nessa plataforma (antes, os setters das configurações de impressão retornavam `ErrDoubleArgUnsupported`): agora as chamadas afetadas passam por um trampolim `runtime.cgocall` que carrega `d0`, algo que o caminho de syscall do Go não consegue fazer (golang.org/issue/62583), por @APshenkin

## Segurança
- O modo servidor (`-tags server`) agora aceita chamadas `/wails/runtime` somente como solicitações `POST` provenientes da própria origem do servidor ou de uma origem listada em `ServerOptions.WebSocketOriginPatterns`; outros métodos recebem HTTP 405, e outras origens do navegador recebem HTTP 403. Implantações atrás de um proxy reverso que reescreva o cabeçalho `Host` devem adicionar sua origem pública a `WebSocketOriginPatterns`
- O modo servidor (`-tags server`) vinculado a um endereço de loopback (os padrões `localhost`, `127.0.0.0/8` ou `::1`) agora rejeita solicitações com HTTP 403, a menos que o cabeçalho `Host` seja `localhost`, um literal de IP, o host de vinculação ou um host que corresponda a `ServerOptions.WebSocketOriginPatterns`, bloqueando ataques de revinculação de DNS contra `/wails/runtime` e os endpoints WebSocket
- As notificações do Windows não usam mais o PowerShell como alternativa quando a API de notificações toast do WinRT falha; em vez disso, o erro é retornado

## v3.0.0-beta.30 - 2026-10-15

## Corrigido
- Corrige o atualizador, que tentava substituir a montagem somente leitura de uma AppImage em vez do arquivo `.AppImage`
- Interrompe o envio de `WindowDidUpdate` e `WindowWillUpdate` no macOS: `NSWindowDidUpdateNotification` é disparado a cada passagem pelo loop de execução, portanto a página recebia cerca de 120 chamadas `evaluateJavaScript:` vazias por segundo e, no macOS, cada 27.0 consumia a ativação transitória do clique do usuário — `navigator.clipboard.writeText()` e outras APIs condicionadas à ativação rejeitavam cliques reais com `NotAllowedError` ([#6170](https://github.com/wailsapp/wails/issues/6170)), por @APshenkin
- Corrige o idioma padrão do info.json do Windows para en-US (0409), permitindo que o Explorer exiba os detalhes da versão (#6210)
- Trata falhas de renomeação entre dispositivos no atualizador do Linux usando uma cópia intermediária como alternativa, no [PR](https://github.com/wailsapp/wails/pull/6135), por @john-okeefe
- Corrige `WebviewWindow.SetMenu` no Windows, que não redesenhava a barra de menus a menos que um tema de barra de menus estivesse configurado, corrigindo [#6103](https://github.com/wailsapp/wails/issues/6103)
- Corrige `WebviewWindow.SetMenu(nil)` no Windows, que causava panic em vez de não realizar nenhuma operação, como no macOS, corrigindo [#6104](https://github.com/wailsapp/wails/issues/6104)
- Preserva entradas do changelog diante de pushes concorrentes usando novas tentativas, no [PR](https://github.com/wailsapp/wails/pull/6094), por @leaanthony
- Corrige a falha de `go mod vendor` com `pattern arm64/WebView2Loader.dll: no matching files found` em todas as plataformas, removendo os embeds que faziam referência a binários nunca distribuídos no módulo, corrigindo [#5782](https://github.com/wailsapp/wails/issues/5782) e [#5376](https://github.com/wailsapp/wails/issues/5376), no [PR](https://github.com/wailsapp/wails/pull/6031), por @Grantmartin2002
- Impede a desmontagem duplicada de janelas no macOS, no [PR](https://github.com/wailsapp/wails/pull/6159), por @charlesbalderson
- Trata WM_CONTEXTMENU para abrir o menu da bandeja do sistema com o botão direito, no [PR](https://github.com/wailsapp/wails/pull/6162), por @wzk2025
- Corrige as cinco variantes de `InvokeSync*` que bloqueavam para sempre seu chamador quando o callback despachado causava panic, corrigindo [#6107](https://github.com/wailsapp/wails/issues/6107), no [PR](https://github.com/wailsapp/wails/pull/6108), por @RALIST

## v3.0.0-beta.29 - 2026-10-13

## Adicionado
- Agora os atalhos globais no macOS levam em conta o layout do teclado: em layouts que não sejam QWERTY (AZERTY, QWERTZ, Dvorak, ...), os aceleradores de letras são associados à tecla indicada no teclado por meio de `UCKeyTranslate`, em vez de à posição física fixa do QWERTY.

## Corrigido
- Inicializa o menu nativo antes de defini-lo no Linux, no [PR](https://github.com/wailsapp/wails/pull/6233), por @abhicodes-007
- Resolve regressões de prontidão para GA no macOS no tratamento da saída e no assistente de configuração, no [PR](https://github.com/wailsapp/wails/pull/6237), por @taliesin-ai
- Corrige os wrappers COM `AddRef`/`Release` do WebView2 no Windows, que informavam como falha um último erro obsoleto da thread. Agora eles retornam a contagem de referências, portanto uma liberação não pode mais encerrar o aplicativo nem causar falha em uma solicitação de recurso (#3931)
- Corrige `ICoreWebView2Deferral.Release`, que chamava o slot `AddRef` e causava o vazamento do adiamento de todas as solicitações de recursos adiadas no Windows (#3931)

## v3.0.0-beta.28 - 2026-10-05

## Corrigido
- Mantém a saída do aplicativo macOS fora dos provedores de arquivos, no [PR](https://github.com/wailsapp/wails/pull/6225), por @taliesin-ai
- Conclui o assistente de configuração e exibe os detalhes do sistema, no [PR](https://github.com/wailsapp/wails/pull/6224), por @taliesin-ai
- Melhora o tratamento de erros e a limpeza de recursos na recuperação do WebView2 no Windows, no [PR](https://github.com/wailsapp/wails/pull/6222), por @leaanthony
- Aprimora a recuperação do WebView2 no Windows para evitar falhas do aplicativo e controladores abandonados no [PR](https://github.com/wailsapp/wails/pull/6220) de @taliesin-ai
- Corrige as orientações sobre autorização de mídia no macOS no [PR](https://github.com/wailsapp/wails/pull/6221) de @taliesin-ai
- Corrige o clique com o botão direito em um ícone da bandeja do sistema no Linux, que executava o manipulador de clique além de abrir o menu, fazendo com que uma janela vinculada alternasse de estado a cada clique com o botão direito (#6018)
- Mantém os manipuladores de callback do WebView2 ativos até a liberação nativa no Windows no [PR](https://github.com/wailsapp/wails/pull/6184) de @taliesin-ai

## v3.0.0-beta.27 - 2026-10-01

## Corrigido
- Corrige a opção `Permissions`, que era ignorada no macOS 12 e posteriores: o método de captura de mídia `WKUIDelegate` agora está implementado, portanto `PermissionAllow` e `PermissionDeny` se aplicam às solicitações de câmera e microfone da mesma forma que no Linux e no Windows. `NSCameraUsageDescription` / `NSMicrophoneUsageDescription` e, quando executado em sandbox, as permissões de dispositivo correspondentes ainda são necessários

## v3.0.0-beta.26 - 2026-09-25

## Adicionado
- O observador aplica verificações HTTP de prontidão ao servidor de desenvolvimento do frontend no [PR](https://github.com/wailsapp/wails/pull/6150) de @atterpac

## Corrigido
- Recupera-se de falhas do processo do WebView2 em vez de deixar uma janela em branco no [PR](https://github.com/wailsapp/wails/pull/6002) de @taliesin-ai

## v3.0.0-beta.25 - 2026-09-22

## Corrigido
- Encerra `wails3 dev` e seus processos em segundo plano quando o aplicativo principal é encerrado (#6048)

## v3.0.0-beta.24 - 2026-09-20

## Segurança
- Restringe as permissões de GITHUB_TOKEN nos fluxos de trabalho de compilação e publicação para seguir o princípio do menor privilégio
- Corrige travessia de caminho e escapes por links simbólicos no middleware de recursos do exemplo de tela
- Mantém funcional a instalação de dependências do assistente de configuração, restringindo-a aos comandos aprovados de gerenciadores de pacotes

## v3.0.0-beta.23 - 2026-09-16

## Adicionado

- Permite substituir o ID do aplicativo GTK no Linux no [PR](https://github.com/wailsapp/wails/pull/5972) de @overlordtm

## Corrigido

- A instância única no Linux reivindica o nome do barramento usando o próprio ID do aplicativo no [PR](https://github.com/wailsapp/wails/pull/5973) de @overlordtm
- Cancela solicitações de esquema personalizado interrompidas no Linux no [PR](https://github.com/wailsapp/wails/pull/6126) de @leaanthony

## v3.0.0-beta.22 - 2026-09-14

## Adicionado

- Conclui as traduções e migra a mídia da comunidade — consulte o [site da documentação](https://v3.wails.io) (1812 páginas atualizadas) no [PR](https://github.com/wailsapp/wails/pull/6122) de @leaanthony
- Documenta botões incorporáveis do Wails com vários estilos e tamanhos no [PR](https://github.com/wailsapp/wails/pull/6120) de @leaanthony

## v3.0.0-beta.21 - 2026-09-13

## Adicionado

- Disponibiliza a documentação do Wails v3 com o M-Press no [PR](https://github.com/wailsapp/wails/pull/6116) de @leaanthony

## Corrigido

- Analisa valores JSON de slug no frontmatter MPD para gerar o changelog no [PR](https://github.com/wailsapp/wails/pull/6118) de @leaanthony
- O atualizador limpa as variáveis de ambiente auxiliares e reinicia o destino original após falhas de backup no [PR](https://github.com/wailsapp/wails/pull/6080) de @cnmax
- Inicia o manipulador de sinais padrão durante App.Run no [PR](https://github.com/wailsapp/wails/pull/6098) de @leaanthony
- O menu do Windows trata menus nil, libera recursos substituídos e redesenha a barra de menus no [PR](https://github.com/wailsapp/wails/pull/6112) de @taliesin-ai
- Restaura o empacotamento MSIX para projetos novos que usam configuração YAML compartilhada no [PR](https://github.com/wailsapp/wails/pull/6115) de @leaanthony
- Corrige a falha no carregamento de bindings JavaScript e TypeScript gerados quando criadores de modelos genéricos fazem referência a declarações auxiliares posteriores e evita estouros de pilha ao criar modelos genéricos mutuamente dependentes (#6062)

## v3.0.0-beta.20 - 2026-09-10

## Alterado

- Atualiza os links da apresentação do Clave para o site e o repositório atuais no [PR](https://github.com/wailsapp/wails/pull/6082) de @01xR4in

## Corrigido

- Cancela solicitações de recursos interrompidas no Windows, incluindo solicitações de workers, preservando os manipuladores keepalive entre navegações. Encaminha contextos de solicitações nativas pelo wrapper do aplicativo nas plataformas Apple. (#5963, #5969)
- Preserva as entradas do changelog durante pushes concorrentes usando novas tentativas no [PR](https://github.com/wailsapp/wails/pull/6094) de @leaanthony
- Corrige a falha de `go mod vendor` com `pattern arm64/WebView2Loader.dll: no matching files found` em todas as plataformas removendo as incorporações que faziam referência a binários nunca distribuídos no módulo, corrigindo [#5782](https://github.com/wailsapp/wails/issues/5782) e [#5376](https://github.com/wailsapp/wails/issues/5376), no [PR](https://github.com/wailsapp/wails/pull/6031) de @Grantmartin2002

## Removido

- Remove o suporte ao carregador nativo do WebView2, substituído pelo carregador em Go puro. Isso remove os binários `WebView2Loader.dll` incorporados e a dependência `github.com/jchv/go-winloader`. A tag de compilação `native_webview2loader` ainda é aceita e não causa mais erros, mas não tem efeito nas compilações v3, no [PR](https://github.com/wailsapp/wails/pull/6031) de @Grantmartin2002
- Remove tags de compilação não utilizadas e a opção de FPS do guia da API do macOS no [PR](https://github.com/wailsapp/wails/pull/6097) de @leaanthony

## v3.0.0-beta.19 - 2026-09-09

## Adicionado

- Protege APIs privadas do macOS com tags de compilação para uso opcional — consulte a [documentação](https://v3.wails.io/features/browser/integration), a [documentação](https://v3.wails.io/features/environment/info), a [documentação](https://v3.wails.io/features/windows/basics), a [documentação](https://v3.wails.io/features/windows/frameless), a [documentação](https://v3.wails.io/features/windows/notch-windows), a [documentação](https://v3.wails.io/features/windows/options), a [documentação](https://v3.wails.io/guides/build/macos), a [documentação](https://v3.wails.io/guides/build/private-macos-apis) e a [documentação](https://v3.wails.io/reference/overview) no [PR](https://github.com/wailsapp/wails/pull/6087) de @leaanthony

## Corrigido

- Rejeita solicitações de runtime acima de 64 MiB com HTTP 413 no [PR](https://github.com/wailsapp/wails/pull/6091) de @leaanthony

## Segurança

- Reforça a segurança das origens do MCP e do acesso remoto com autenticação por token no [PR](https://github.com/wailsapp/wails/pull/6092) de @leaanthony

## v3.0.0-beta.18 - 2026-09-08

## Corrigido

- Corrige o vazamento de memória de Calloc no Linux e no Darwin usando receptores de ponteiro no [PR](https://github.com/wailsapp/wails/pull/6083) de @4RH1T3CT0R7

## v3.0.0-beta.17 - 2026-09-06

## Corrigido

- Windows: um `GetRequest` com falha ou nil no manipulador WebResourceRequested não encerra mais o processo (`log.Fatal` / panic por desreferência de nil) — em vez disso, a solicitação é descartada e registrada no log, no [PR](https://github.com/wailsapp/wails/pull/6006) de @midagedev

## v3.0.0-beta.16 - 2026-08-29

## Alterado

- Solicita a senha de notarização em uma nova janela do terminal no [PR](https://github.com/wailsapp/wails/pull/6029) por @leaanthony

## Corrigido

- Trata corretamente os tipos de clique no ícone da bandeja do sistema no macOS no [PR](https://github.com/wailsapp/wails/pull/5919) por @ChewbaccaCookie
- A CI remove os repositórios apt não utilizados da Microsoft antes da atualização no [PR](https://github.com/wailsapp/wails/pull/6041) por @Grantmartin2002

## v3.0.0-beta.15 - 2026-08-27

## Corrigido

- Aumenta para 60 segundos o tempo limite de incorporação do WebView2 no [PR](https://github.com/wailsapp/wails/pull/6043) por @Grantmartin2002

## v3.0.0-beta.14 - 2026-08-26

## Corrigido

- Nomeia corretamente os pressionamentos de teclas Control-letra no macOS no [PR](https://github.com/wailsapp/wails/pull/6032) por @taliesin-ai
- Corrige os ícones ICO da bandeja e acompanha o tema da barra de tarefas no Windows no [PR](https://github.com/wailsapp/wails/pull/6016) por @nik9play

## v3.0.0-beta.13 - 2026-08-25

## Corrigido

- Mantém o processamento de tarefas da thread principal no macOS durante a execução de um loop modal no [PR](https://github.com/wailsapp/wails/pull/6026) por @leaanthony
- Faz com que o armazenamento seguro em dispositivos móveis possa falhar e adote o fechamento seguro em caso de falha no [PR](https://github.com/wailsapp/wails/pull/5923) por @mortenolsrud
- Executa os hooks de eventos da aplicação mesmo quando nenhum listener está registrado no [PR](https://github.com/wailsapp/wails/pull/5999) por @archy-rock3t-cloud
- Corrige erros de digitação em comentários e na documentação localizada no [PR](https://github.com/wailsapp/wails/pull/6023) por @haoku123
- Remove os binários pré-compilados do macOS versionados em `v3/examples` no [PR](https://github.com/wailsapp/wails/pull/6025) por @4RH1T3CT0R7

## v3.0.0-beta.12 - 2026-08-21

## Adicionado

- Adiciona janelas de notificação para o entalhe do macOS, com exemplo de ciclo de vida e telemetria — consulte a [documentação](https://v3.wails.io/features/windows/notch-windows) no [PR](https://github.com/wailsapp/wails/pull/6010) por @leaanthony
- Adiciona suporte a janelas NSPanel do macOS, com novas opções e integração nativa — consulte a [documentação](https://v3.wails.io/features/windows/options) no [PR](https://github.com/wailsapp/wails/pull/6008) por @leaanthony

## Corrigido

- Impede que o Prepare do SQLite fique travado durante chamadas simultâneas no [PR](https://github.com/wailsapp/wails/pull/5998) por @archy-rock3t-cloud

## v3.0.0-beta.11 - 2026-08-20

## Removido

- Remove da documentação o rastreador de implementação obsoleto no [PR](https://github.com/wailsapp/wails/pull/6005) por @leaanthony

## v3.0.0-beta.10 - 2026-08-19

## Corrigido

- Corrige o host Linux GTK4, que descartava os argumentos de inicialização de protocolos personalizados e associações de arquivos, no [PR](https://github.com/wailsapp/wails/pull/6000) por @midagedev
- Trata corretamente linhas excluídas e correções na mesma origem durante a validação do changelog no [PR](https://github.com/wailsapp/wails/pull/5993) por @taliesin-ai

## v3.0.0-beta.9 - 2026-08-16

## Adicionado

- Adiciona o servidor seguro wails3 mcp para gerenciamento de projetos assistido por agentes no [PR](https://github.com/wailsapp/wails/pull/5896) por @leaanthony
- Adiciona documentação sobre modelos em bindings — consulte a [documentação](https://v3.wails.io/features/bindings/models) no [PR](https://github.com/wailsapp/wails/pull/5988) por @taliesin-ai
- Adiciona suporte à instalação com rpm-ostree em sistemas Linux atômicos no [PR](https://github.com/wailsapp/wails/pull/5987) por @leaanthony
- Adiciona geração e publicação nativas semanais do gráfico de histórico de estrelas — consulte a [documentação](https://v3.wails.io/credits), a [documentação](https://v3.wails.io/de/credits), a [documentação](https://v3.wails.io/fr/credits), a [documentação](https://v3.wails.io/id/credits), a [documentação](https://v3.wails.io/ja/credits), a [documentação](https://v3.wails.io/ko/credits), a [documentação](https://v3.wails.io/pt/credits), a [documentação](https://v3.wails.io/ru/credits), a [documentação](https://v3.wails.io/zh-cn/credits) e a [documentação](https://v3.wails.io/zh-tw/credits) no [PR](https://github.com/wailsapp/wails/pull/5986) por @leaanthony
- Adiciona um pacote mac exclusivo para Darwin destinado à resolução de recursos do pacote da aplicação — consulte a [documentação](https://v3.wails.io/guides/build/macos) no [PR](https://github.com/wailsapp/wails/pull/5965) por @leaanthony
- Adiciona a página de demonstração do Condui e uma entrada no índice — consulte a [documentação](https://v3.wails.io/community/showcase/condui) e a [documentação](https://v3.wails.io/community/showcase) no [PR](https://github.com/wailsapp/wails/pull/5962) por @mgueregath
- Adiciona a página de demonstração do Redis Viewer, com capturas de tela e link para o projeto — consulte a [documentação](https://v3.wails.io/community/showcase) e a [documentação](https://v3.wails.io/community/showcase/redisviewer) no [PR](https://github.com/wailsapp/wails/pull/5984) por @redisviewer

## Alterado

- Atualiza os sinalizadores da aplicação GTK para G<em>APPLICATION</em>NON_UNIQUE no Linux no [PR](https://github.com/wailsapp/wails/pull/5971) por @overlordtm
- Registra eventos de janela ausentes no nível de depuração, em vez de aviso, no [PR](https://github.com/wailsapp/wails/pull/5914) por @julianstorer

## Corrigido

- Permite que os atalhos de teclado registrados no macOS tenham precedência sobre a webview no [PR](https://github.com/wailsapp/wails/pull/5902) por @julianstorer
- Corrige links quebrados na barra lateral da documentação no [PR](https://github.com/wailsapp/wails/pull/5937) por @northes
- Cancela os contextos de solicitação de recursos no macOS e no iOS quando o WebKit interrompe a tarefa correspondente do esquema personalizado (#5963)
- Trata a mensagem WindowSetFullscreenButtonEnabled no [PR](https://github.com/wailsapp/wails/pull/5976) por @archy-rock3t-cloud
- Importa Fragment no modelo preact-ts para resolver a falha de compilação no [PR](https://github.com/wailsapp/wails/pull/5979) por @haoku123
- Impede que aplicações de serviço legadas baseadas em GTK3 falhem quando a descoberta de telas é executada antes de haver uma janela ativa ou uma tela disponível (#5966)
- Permite que execuções de lançamento com versão explícita prossigam quando o changelog de alterações não lançadas estiver vazio (#5977)

## Segurança

- Atualiza os lockfiles do nanoid do site para a versão corrigida 3.3.18, solucionando os avisos de segurança no [PR](https://github.com/wailsapp/wails/pull/5985) por @taliesin-ai

## v3.0.0-beta.8 - 2026-08-12

## Adicionado

- Adiciona a geração de URLs da documentação às entradas automáticas do changelog no [PR](https://github.com/wailsapp/wails/pull/5957) por @taliesin-ai
- Adiciona Streams: fluxos bidirecionais de bytes entre Go e JavaScript, com o modelo de programação WebSocket e sem socket de escuta. Declare um fluxo em Go com `app.HandleStream(name, handler)` e conecte-se pelo frontend com `Stream(name)`, que retorna um objeto com o formato de `WebSocket`. Go→JS é transportado por uma consulta mantida por janela pelo servidor de recursos, e JS→Go por um POST normal; nenhuma porta TCP é vinculada e nada passa por `evaluateJavaScript`. Em builds de servidor (`-tags server`), o mesmo handler é servido por um WebSocket real, de modo que o código da aplicação é idêntico em todos os builds. por @leaanthony
- Move a entrada do changelog da caixa de correio para Não lançado no [PR](https://github.com/wailsapp/wails/pull/5935) por @leaanthony

## Alterado

- Atualiza a geração automática da barra lateral da documentação e a derivação do tipo de autor do blog no [PR](https://github.com/wailsapp/wails/pull/5938) por @leaanthony

## Corrigido

- A inicialização do WebView2 usa um prazo limite e um processamento de mensagens no [PR](https://github.com/wailsapp/wails/pull/5952) por @leaanthony
- O teste de cookies do WebView2 é ignorado na CI, salvo quando habilitado explicitamente, e fixa a execução à thread atual do sistema operacional em [PR](https://github.com/wailsapp/wails/pull/5951) por @leaanthony
- Os construtores de menus do Windows restauram os IDs de comando dos itens-pai de submenus em [PR](https://github.com/wailsapp/wails/pull/5944) por @gilad-ch
- Alinha a imagem oficial de compilação cruzada com a linha de base de suporte ao GTK 4.14+ no Linux (#5928)
- Configura o projeto Xcode para iOS para preservar os sinalizadores herdados do linker e adicionar -ObjC em [PR](https://github.com/wailsapp/wails/pull/5915) por @mortenolsrud
- Corrige a rotatividade excessiva de conexões TCP no proxy de ativos `wails3 dev` em frontends grandes, que poderia esgotar as portas efêmeras do host e fazer processos não relacionados falharem com `EADDRNOTAVAIL`
- Enfileira o JavaScript de eventos de cada janela para garantir despacho ordenado e contrapressão em [PR](https://github.com/wailsapp/wails/pull/5934) por @leaanthony

## Removido

- Remove o pipeline de lançamento de binários para desktop: as versões v3 contêm apenas tags, e a CLI `wails3` é instalada com `go install`. Exclui `release-v3.yml` e a etapa da versão nightly que o acionava em [PR](https://github.com/wailsapp/wails/pull/5946) por @leaanthony

## v3.0.0-beta.7 - 2026-08-11

## Adicionado

- Adiciona uma preferência de reprodução automática no macOS para desabilitar a exigência de ação do usuário para reproduzir mídia em [PR](https://github.com/wailsapp/wails/pull/5512) por @Eyalm321
- Move a entrada do changelog referente à caixa de correio para Não lançado em [PR](https://github.com/wailsapp/wails/pull/5935) por @leaanthony

## Alterado

- A animação de zoom do macOS usa CADisplayLink ou NSTimer para proporcionar um desempenho mais fluido em [PR](https://github.com/wailsapp/wails/pull/5945) por @savely-krasovsky

## Corrigido

- Configura o projeto Xcode para iOS para preservar os sinalizadores herdados do linker e adicionar -ObjC em [PR](https://github.com/wailsapp/wails/pull/5915) por @mortenolsrud
- Corrige a rotatividade excessiva de conexões TCP no proxy de ativos `wails3 dev` em frontends grandes, que poderia esgotar as portas efêmeras do host e fazer processos não relacionados falharem com `EADDRNOTAVAIL`
- Enfileira o JavaScript de eventos de cada janela para garantir despacho ordenado e contrapressão em [PR](https://github.com/wailsapp/wails/pull/5934) por @leaanthony

### Adicionado

- Implementa uma caixa de correio FIFO assíncrona genérica para a entrega ordenada de eventos em [PR](https://github.com/wailsapp/wails/pull/5851) por @savely-krasovsky e @DevLumuz

## v3.0.0-beta.6 - 2026-08-09

## Adicionado

- Implementa armazenamento limitado no lado do host para eventos grandes demais e entrega ordenada de JavaScript em [PR](https://github.com/wailsapp/wails/pull/5930) por @leaanthony
- Implementa o salto do ícone no Dock do macOS para o alerta visual de janelas em [PR](https://github.com/wailsapp/wails/pull/5921) por @julianstorer

## Corrigido

- O servidor de ativos preserva os erros do detector de tipo de conteúdo e os prefixos não gravados durante o flush em [PR](https://github.com/wailsapp/wails/pull/5931) por @leaanthony
- Impede que aplicativos macOS falhem ao substituir o menu do aplicativo a partir de um callback do Wails
- Corrige menus nativos ilegíveis no Windows 10 1809 / Windows Server 2019 (build 17763). As exportações de modo escuro do uxtheme estavam condicionadas à build 18334, portanto a ativação opcional do modo escuro no nível do aplicativo nunca era executada nesses hosts: o fundo do menu era pintado de escuro, mas o Windows continuava desenhando o texto do menu com o tema claro, deixando texto escuro sobre fundo escuro. Os ordinais existem desde 17763, portanto a condição agora corresponde a essa versão.
- Corrige `w32.GetStockObject`, que chamava `GetDeviceCaps` em vez de `GetStockObject`, fazendo com que retornasse 0 para todos os objetos predefinidos.
- Melhora o tratamento e o relatório de erros de download do bootstrapper do WebView2 em [PR](https://github.com/wailsapp/wails/pull/5924) por @jannskiee

## v3.0.0-beta.5 - 2026-08-07

## Corrigido

- A ativação de aplicativos no macOS agora respeita a política de ativação somente para aplicativos comuns em [PR](https://github.com/wailsapp/wails/pull/5897) por @julianstorer
- Protege contra janelas GTK não inicializadas em builds para Linux em [PR](https://github.com/wailsapp/wails/pull/5898) por @julianstorer
- Define explicitamente uma cor de fundo opaca para janelas WebKit no Linux antes de carregar a URL em [PR](https://github.com/wailsapp/wails/pull/5899) por @julianstorer

## v3.0.0-beta.4 - 2026-08-05

## Alterado

- As tarefas de build do Android usam arm64 por padrão, e deploy-emulator seleciona a arquitetura do host em [PR](https://github.com/wailsapp/wails/pull/5890) por @mortenolsrud

## Corrigido

- Preserva o estado de zoom das janelas do macOS durante o arrasto e reduz o movimento em [PR](https://github.com/wailsapp/wails/pull/5900) por @leaanthony
- Corrige a build do Windows no modo servidor adicionando `!server` à restrição de build `webview_window_windows_nonclient.go`

## v3.0.0-beta.3 - 2026-08-03

## Adicionado

- Documenta, nos detalhes de implementação, a conclusão da verificação beta da Fase 10 em [PR](https://github.com/wailsapp/wails/pull/5881) por @leaanthony

## Corrigido

- Passa o identificador da janela à API de modo escuro do Windows e valida os argumentos em [PR](https://github.com/wailsapp/wails/pull/5877) por @leaanthony
- Centraliza a resolução do estado dos botões da barra de título do macOS para janelas sem moldura em [PR](https://github.com/wailsapp/wails/pull/5870) por @taliesin-ai
- Impede que o texto de menus nativos fique ilegível quando um aplicativo do Windows solicita o modo escuro enquanto o tema de aplicativos do Windows está claro. O menu agora usa o fundo nativo claro correspondente até que o Windows possa renderizar o texto escuro do menu.
- Corrige menus nativos ilegíveis no Windows 10 1809 / Windows Server 2019 (build 17763). As exportações de modo escuro do uxtheme estavam condicionadas à build 18334, portanto a ativação opcional do modo escuro no nível do aplicativo nunca era executada nesses hosts: o fundo do menu era pintado de escuro, mas o Windows continuava desenhando o texto do menu com o tema claro, deixando texto escuro sobre fundo escuro. Os ordinais existem desde 17763, portanto a condição agora corresponde a essa versão.

## v3.0.0-beta.2 - 2026-08-02

## Alterado

- Promove v3 de alfa para beta
- Documenta os padrões inteligentes da bandeja do sistema e o comportamento de ocultação automática de pop-ups, com cobertura de regressão para a seleção do manipulador de cliques (#5840).
- O atualizador do GitHub exclui por padrão os ativos de instaladores do Windows em [PR](https://github.com/wailsapp/wails/pull/5861) por @leaanthony
- Adiciona suporte a cantos arredondados, quadrados e com raio personalizado em janelas sem moldura do macOS em [PR](https://github.com/wailsapp/wails/pull/5866) por @leaanthony

## Corrigido

- Informa os tamanhos atuais das janelas GTK4 e emite eventos de redimensionamento e de estado maximizado, minimizado e em tela cheia a partir da superfície configurada (#5830).
- Corrige uma falha do WebKit no Linux ao enviar Blob ou FormData em requisições fetch em [PR](https://github.com/wailsapp/wails/pull/5854) por @taliesin-ai
- O shim de fetch passa undefined quando os cabeçalhos de Blob/FormData estão ausentes em [PR](https://github.com/wailsapp/wails/pull/5865) por @leaanthony

## v3.0.0-alpha2.122 - 2026-08-01

## Adicionado

## Alterado

- Adiciona suporte a cantos arredondados, quadrados e com raio personalizado em janelas sem moldura do macOS em [PR](https://github.com/wailsapp/wails/pull/5866) por @leaanthony

## Corrigido

- O shim de fetch passa undefined para cabeçalhos ausentes de Blob/FormData no [PR](https://github.com/wailsapp/wails/pull/5865) por @leaanthony

## v3.0.0-alpha2.121 - 2026-07-31

## Adicionado

- Adiciona suporte ao empacotamento DMG no macOS, com novas opções e tarefas de compilação, no [PR](https://github.com/wailsapp/wails/pull/5857) por @leaanthony

## Alterado

- O atualizador do GitHub exclui, por padrão, os artefatos do instalador do Windows no [PR](https://github.com/wailsapp/wails/pull/5861) por @leaanthony

## Corrigido

- Corrige a falha do WebKit no Linux ao enviar Blob ou FormData em requisições fetch no [PR](https://github.com/wailsapp/wails/pull/5854) por @taliesin-ai

## v3.0.0-alpha2.120 - 2026-07-31

## Adicionado

- Implementa ações de clique duplo na barra de título do macOS para maximizar ou minimizar no [PR](https://github.com/wailsapp/wails/pull/5853) por @taliesin-ai

## Alterado

- Atualiza o tutorial do serviço de QR para usar NewServiceWithOptions e adiciona espaçamento no [PR](https://github.com/wailsapp/wails/pull/5849) por @jeongkyu

## Corrigido

- Mantém o WKWebView responsivo durante o zoom no macOS no [PR](https://github.com/wailsapp/wails/pull/5856) por @leaanthony
- Corrige as consultas de tamanho da janela no GTK4 e emite eventos de redimensionamento, maximização, minimização e estado de tela cheia a partir do `GdkSurface` configurado.

## v3.0.0-alpha2.119 - 2026-07-27

## Corrigido

- Atualiza a documentação em vários idiomas para incluir diagramas de arquitetura no [PR](https://github.com/wailsapp/wails/pull/5833) por @taliesin-ai

## v3.0.0-alpha2.118 - 2026-07-26

## Adicionado

- Fornece caminhos padrão para as entradas e saídas da geração de ícones no [PR](https://github.com/wailsapp/wails/pull/5825) por @taliesin-ai
- Adiciona os módulos de entrada do código-fonte ao campo sideEffects do package.json do runtime no [PR](https://github.com/wailsapp/wails/pull/5797) por @savely-krasovsky
- Adiciona uma seção sobre licença e proveniência ao guia de contribuição no [PR](https://github.com/wailsapp/wails/pull/5816) por @taliesin-ai

## Corrigido

- Aplica CSS sem moldura com escopo no GTK4 para remover o raio da borda no [PR](https://github.com/wailsapp/wails/pull/5800) por @savely-krasovsky
- Trata de forma segura as falhas ao obter a posição do cursor no Windows para menus pop-up e enumeração de telas no [PR](https://github.com/wailsapp/wails/pull/5789) por @wayneforrest
- A caixa de diálogo para abrir arquivos do macOS filtra corretamente as extensões e valida os arquivos permitidos pelo sufixo no [PR](https://github.com/wailsapp/wails/pull/5678) por @phergul
- Protege a inicialização do modo escuro no Windows contra chamadas de API nil no [PR](https://github.com/wailsapp/wails/pull/5793) por @roachadam
- Corrige uma falha de compilação de 32 bits no atualizador: a constante `maxArchiveTotalSize` (2 GiB) excedia a capacidade do `int` da plataforma quando passada para `fmt.Errorf` em `GOARCH=386`. Agora, ela tem explicitamente o tipo `int64`.
- Corrige um panic por ponteiro nil na inicialização quando uma janela usa uma barra de título escura (ou escura conforme o sistema) em compilações do Windows que não carregam as APIs de modo escuro do uxtheme, como Windows 10 1809 / Windows Server 2019 (compilação 17763). As chamadas de `AllowDarkModeForWindow` na configuração do tema da janela agora são protegidas contra nil, de modo consistente com a proteção já usada em `w32.SetMenuTheme`.

## v3.0.0-alpha2.117 - 2026-07-08

## Adicionado

- Implementa lógica personalizada de teste de acerto para regiões não cliente no Windows no [PR](https://github.com/wailsapp/wails/pull/5462) por @savely-krasovsky

## Alterado

- Configura a detecção de escala do monitor do WebView2 com base em UseVisualHosting no [PR](https://github.com/wailsapp/wails/pull/5761) por @wayneforrest

## v3.0.0-alpha2.116 - 2026-07-07

## Adicionado

- Atualiza a documentação de perguntas frequentes para se concentrar nos recursos e nas orientações do Wails v3 no [PR](https://github.com/wailsapp/wails/pull/5763) por @taliesin-ai

## v3.0.0-alpha2.115 - 2026-07-06

## Corrigido

- Corrige o problema em que `Menu.Update()` não recompilava o menu nativo no Linux com GTK4 (#5659, diagnosticado e corrigido de forma independente por @puneetdixit200 em #5539)
- Corrige uma falha ao enumerar as telas do macOS após uma alteração de monitor, copiando as strings de ID/nome das telas e capturando a contagem (#5565, diagnosticado e corrigido de forma independente por @x-haose em #5584)
- Corrige uma falha no Windows quando `WM_ERASEBKGND` desenha um fundo sólido durante uma transição de minimizar/restaurar na qual `GetClientRect` retorna nil (proteção relatada por @sinspired em #5636)
- Corrige o erro em que falhas nos bindings do frontend eram sempre analisadas como texto, por @mbaklor em #5690
- Corrige uma falha de compilação no Windows ao usar a tag de compilação `server`, causada pela ausência, nos arquivos da GUI do Windows, da restrição de compilação `!server` que os equivalentes do macOS e do Linux já possuem (#5680)

## v3.0.0-alpha2.114 - 2026-07-05

## Adicionado

- Implementa o protocolo Update Manifest e o provedor de endpoint no [PR](https://github.com/wailsapp/wails/pull/5720) por @taliesin-ai

## Alterado

- Incorpora o binding `webview2` ao módulo v3 como `v3/internal/webview2`, removendo o módulo independente, seus fluxos de trabalho de lançamento/sincronização nightly e a alternância de versões do go.mod (v3 é seu único consumidor), no [PR](https://github.com/wailsapp/wails/pull/5711) por @taliesin-ai

## Corrigido

- Move a detecção de escala do monitor do WebView2 e a correção de ressincronização do host após alteração de DPI para a seção Não lançado no [PR](https://github.com/wailsapp/wails/pull/5750) por @taliesin-ai
- Atualiza o marshaling COM do WebView2 para parâmetros float64 e BOOL no [PR](https://github.com/wailsapp/wails/pull/5741) por @wayneforrest
- Evita panic e desreferência de nil nas atualizações e na destruição do ícone da bandeja do sistema no Windows no [PR](https://github.com/wailsapp/wails/pull/5703) por @wayneforrest
- Corrige o problema em que janelas ocultas não voltavam a ser ocultadas corretamente no Windows no [PR](https://github.com/wailsapp/wails/pull/5743) por @wayneforrest
- Sincroniza a visibilidade do controlador WebView2 com as ações de minimizar, maximizar e restaurar a janela no [PR](https://github.com/wailsapp/wails/pull/5742) por @wayneforrest

### Corrigido

- Reativa a detecção de escala do monitor do WebView2 e condiciona a ressincronização do host à alteração de DPI no [PR](https://github.com/wailsapp/wails/pull/5734) por @taliesin-ai, com base na correção validada por @randalmurphal, com verificação da causa raiz por @eleclin e testes de hardware por @qq540491950

## v3.0.0-alpha2.113 - 2026-07-04

## Adicionado

- Adiciona um aviso ao compilar um AAB de lançamento sem definir `ANDROID_KEYSTORE_FILE` (o Google Play rejeita pacotes assinados para depuração) e documenta o empacotamento e a assinatura de App Bundles no [PR](https://github.com/wailsapp/wails/pull/5730) por @taliesin-ai
- Adiciona a documentação Por que usar o Wails em vários idiomas no [PR](https://github.com/wailsapp/wails/pull/5739) por @taliesin-ai
- Oferece suporte ao mapeamento de Go time.Time para JS Date ou string nos bindings no [PR](https://github.com/wailsapp/wails/pull/5398) por @fbbdev
- Adiciona tarefas de empacotamento de Android App Bundle (AAB) (`bundle`, `bundle:fat`, `assemble:aab`, `assemble:aab:release`) para envio à Play Store — as tarefas de APK continuam disponíveis para testes locais ou em emuladores no [PR](https://github.com/wailsapp/wails/pull/5728) de @mortenolsrud (corrige [#5726](https://github.com/wailsapp/wails/issues/5726))
- Adiciona destinos de tarefas para dispositivos Android físicos e retoma as permissões de câmera/localização no [PR](https://github.com/wailsapp/wails/pull/5735) de @taliesin-ai

## Alterado

- Atualiza `webview2` para v1.0.28 ([notas da versão](https://github.com/wailsapp/wails/releases/tag/webview2%2Fv1.0.28)).
- Atualiza `compileSdk`/`targetSdk` do modelo Android de 34 para 35, conforme exigido pelo Google Play para novos envios de aplicativos, no [PR](https://github.com/wailsapp/wails/pull/5730) de @taliesin-ai

## Corrigido

- Corrige as máscaras de avatar incorporadas no sponsorkit no [PR](https://github.com/wailsapp/wails/pull/5745) de @leaanthony
- Corrige a criação automática de AVD do Android que selecionava a imagem de sistema ou a versão de cmdline-tools incorreta devido à ordenação lexicográfica das versões, no [PR](https://github.com/wailsapp/wails/pull/5730) de @taliesin-ai
- Corrige o assistente de configuração que sugeria uma versão desatualizada do Android NDK (agora 26.3.11579264, correspondente ao requisito documentado) no [PR](https://github.com/wailsapp/wails/pull/5730) de @taliesin-ai
- Atualiza a documentação em francês sobre SvelteKit e opções no [PR](https://github.com/wailsapp/wails/pull/5744) de @leaanthony
- Corrige SIGSEGV durante a enumeração de telas do macOS em alterações de monitor no [PR](https://github.com/wailsapp/wails/pull/5516) de @flofreud

## v3.0.0-alpha2.112 - 2026-07-03

## Adicionado

- Adiciona um gerador de SVG de colaboradores baseado em Go e atualiza as páginas de créditos da documentação e do site no [PR](https://github.com/wailsapp/wails/pull/5724) de @taliesin-ai
- Adiciona suporte ao mapeamento de Go time.Time para JS Date ou string nos bindings no [PR](https://github.com/wailsapp/wails/pull/5398) de @fbbdev

## Alterado

- Substitui o pipeline de imagens de patrocinadores baseado em Node por um gerador em Go no [PR](https://github.com/wailsapp/wails/pull/5719) de @taliesin-ai

## Corrigido

- Corrige o script de instalação das dependências de ativos da compilação Android no [PR](https://github.com/wailsapp/wails/pull/5729) de @taliesin-ai
- Rejeita o caractere de controle U+0085 (NEXT LINE) em `ValidateAndSanitizeURL`, completando a cobertura de espaços em branco do validador de URLs
- Recalcula o quadro DWM ao alterar o DPI de janelas sem moldura no [PR](https://github.com/wailsapp/wails/pull/4785) de @leaanthony
- Corrige a falha na detecção da zona de soltura de DnD no Windows com escalas diferentes de 100% no [PR](https://github.com/wailsapp/wails/pull/4632) de @yulesxoxo
- Adiciona gerenciamento explícito de memória Objective-C para objetos Cocoa em caixas de diálogo, menus, bandeja e notificações do Darwin no [PR](https://github.com/wailsapp/wails/pull/5714) de @taliesin-ai
- Corrige bugs do backend CGO no Linux e problemas da bandeja do sistema no [PR](https://github.com/wailsapp/wails/pull/5718) de @taliesin-ai

## v3.0.0-alpha2.111 - 2026-07-01

## Adicionado

- Adiciona HappyTools à vitrine da comunidade no [PR](https://github.com/wailsapp/wails/pull/5061) de @Aliuyanfeng
- Adiciona suporte à localidade indonésia e documentação abrangente no [PR](https://github.com/wailsapp/wails/pull/5643) de @triadmoko
- Adiciona a opção DisableMenu a WindowsWindow no [PR](https://github.com/wailsapp/wails/pull/4813) de @leaanthony

## Alterado

- Atualiza o modelo do Taskfile e a CLI para despachar tarefas de compilação/empacotamento com GOOS e ARCH no [PR](https://github.com/wailsapp/wails/pull/5617) de @leaanthony

## Corrigido

- Corrige um problema com o agrupamento de janelas em abas no macOS no [PR](https://github.com/wailsapp/wails/pull/5708) de @taliesin-ai

## Removido

- Remove os arquivos MDX traduzidos para o alemão das seções de contribuição, recursos e guias no [PR](https://github.com/wailsapp/wails/pull/5702) de @taliesin-ai

## v3.0.0-alpha2.110 - 2026-06-30

## Adicionado

- Implementa o recarregamento e o recarregamento forçado da WebView no macOS e adiciona recuperação após o encerramento do processo WebContent no [PR](https://github.com/wailsapp/wails/pull/5129) de @wayneforrest
- Adiciona documentação abrangente em alemão para contribuição, recursos e guias no [PR](https://github.com/wailsapp/wails/pull/5396) de @leaanthony
- Aprimora as notificações com som, anexos, agendamento e uma API de atualização no [PR](https://github.com/wailsapp/wails/pull/5333) de @popaprozac

## Corrigido

- Recalcula o quadro DWM ao alterar o DPI de janelas sem moldura no [PR](https://github.com/wailsapp/wails/pull/4785) de @leaanthony
- Corrige a falha na detecção da zona de soltura de DnD no Windows com escalas diferentes de 100% no [PR](https://github.com/wailsapp/wails/pull/4632) de @yulesxoxo

## v3.0.0-alpha2.109 - 2026-06-29

## Adicionado

- Adiciona exemplos de código à documentação de EventsEmit no [PR](https://github.com/wailsapp/wails/pull/5026) de @iamhabbeboy
- Adiciona a opção de hospedagem visual do WebView2 no Windows no [PR](https://github.com/wailsapp/wails/pull/5380) de @MerIijn
- Adiciona Klustr à documentação da vitrine da comunidade no [PR](https://github.com/wailsapp/wails/pull/5536) de @SametKUM
- Adiciona Kira à vitrine da comunidade, com novas páginas e uma entrada no changelog, no [PR](https://github.com/wailsapp/wails/pull/5685) de @thiennguyen93
- Adiciona uma seção de feedback ao guia do serviço MCP no [PR](https://github.com/wailsapp/wails/pull/5694) de @taliesin-ai

## Alterado

- O modo servidor agora conta com uma compilação de produção de primeira classe, consistente com as tarefas de compilação para desktop (#5693). `task build:server` compila um binário de produção por padrão (`-tags server,production`, `-trimpath`, sem símbolos) e aceita `DEV=true` (servidor de desenvolvimento), `OBFUSCATED=true` (garble) e `EXTRA_TAGS`. `task run:server` executa um servidor de desenvolvimento. `Dockerfile.server` / `task build:docker` compilam primeiro o servidor de produção (`-tags server,production`) e o frontend de produção; por padrão, a imagem usa uma compilação estática em Go puro sobre distroless/static, com `CGO_ENABLED`, `GO_IMAGE` e `RUNTIME_IMAGE` expostos como argumentos de compilação substituíveis para aplicativos CGO.

## Corrigido

- Evita uma falha ao fechar uma janela com chamadas assíncronas pendentes no [PR](https://github.com/wailsapp/wails/pull/4435) de @leaanthony
- Evita a ativação da janela ao abrir aplicativos ocultos no Windows no [PR](https://github.com/wailsapp/wails/pull/5249) de @leaanthony
- Garante que os metadados de requisição do WebKit, a conclusão da resposta e o processamento do fluxo do corpo sejam executados na thread principal do GTK no [PR](https://github.com/wailsapp/wails/pull/5668) de @taliesin-ai
- Corrige `Menu.Update()`, que não recompilava o menu nativo no Linux com GTK4 (#5659, diagnosticado e corrigido de forma independente por @puneetdixit200 em #5539)
- Corrige uma falha ao enumerar telas do macOS durante uma alteração de monitor, copiando as strings de ID/nome da tela e capturando a contagem (#5565, diagnosticado e corrigido de forma independente por @x-haose em #5584)
- Corrige o conteúdo do WebView2 que encolhia e depois desaparecia após arrastar uma janela entre monitores com DPIs diferentes no Windows, reafirmando os limites do controlador no manipulador `WM_DPICHANGED`, em correspondência com a ressincronização de DPI ao restaurar a janela (#5677)

## v3.0.0-alpha2.108 - 2026-06-28

## Adicionado

- Adiciona atalhos de teclado globais (em todo o sistema) por meio de `app.GlobalShortcut` (`Register`, `Unregister`, `UnregisterAll`, `IsRegistered`, `GetAll`). Os atalhos são acionados mesmo quando o aplicativo não está em foco. Implementação nativa para cada plataforma, sem dependências de terceiros: teclas de atalho do Carbon no macOS, `RegisterHotKey` no Windows, `XGrabKey` no X11 e a interface de atalhos globais do XDG Desktop Portal no Wayland.
- Adiciona um servidor MCP integrado: um servidor Model Context Protocol que inicia automaticamente quando o aplicativo é compilado com a tag `mcp`, permitindo que agentes de LLM testem e controlem um aplicativo Wails em execução — controle de janelas, inspeção do DOM, avaliação de JavaScript, chamadas a métodos vinculados, eventos e entrada simulada de mouse/teclado, representada por um cursor animado na tela. Não requer código do usuário: a tag `mcp` é adicionada automaticamente por `wails3 build`/`wails3 dev` quando `WAILS_MCP=1` está definido. Configuração feita inteiramente por variáveis de ambiente (`WAILS_MCP_HOST`, `WAILS_MCP_PORT`, `WAILS_MCP_TIMEOUT`, `WAILS_MCP_HIDE_CURSOR`).

## Corrigido

- Corrige `Menu.Update()`, que não reconstruía o menu nativo no Linux com GTK4 (#5659, diagnosticado e corrigido de forma independente por @puneetdixit200 em #5539)
- Corrige uma falha ao enumerar telas do macOS após uma alteração de monitores, copiando as strings de ID/nome das telas e capturando uma cópia da contagem (#5565, diagnosticado e corrigido de forma independente por @x-haose em #5584)

## v3.0.0-alpha2.107 - 2026-06-27

## Adicionado

- Adiciona documentação experimental do Wake com navegação pela barra lateral no [PR](https://github.com/wailsapp/wails/pull/5613), por @leaanthony

## v3.0.0-alpha2.106 - 2026-06-24

## Alterado

- Atualiza `webview2` para v1.0.27.
  - ci(webview2): corrige a compilação de lançamento (compilação cruzada para Windows + go.sum completo) (#5671)\

  **Diff completo:** https://github.com/wailsapp/wails/compare/webview2/v1.0.26...webview2/v1.0.27

- Remove go vet da compilação cruzada no fluxo de trabalho de lançamento do webview2 no [PR](https://github.com/wailsapp/wails/pull/5672), por @taliesin-ai
- Atualiza o modelo OpenRouter do changelog automático para google/gemini-2.5-flash-lite no [PR](https://github.com/wailsapp/wails/pull/5670), por @taliesin-ai
- Atualiza `webview2` para v1.0.26.

### Correções

- **Recupera-se de erros transitórios de COM em tempo de execução em vez de encerrar** (#5658, #5580). Antes, `Chromium.errorCallback` chamava `os.Exit(1)` para *qualquer* erro de COM; portanto, uma falha momentânea recuperável após a inicialização encerrava todo o aplicativo. Agora, os caminhos de tempo de execução (`Resize`/`GetClientRect`, `Navigate`/`NavigateToString`, `Init`, `MessageReceived`, `PutZoomFactor`, `OpenDevToolsWindow`) registram o erro e se recuperam. Em particular, uma mensagem Web malformada/não confiável em `MessageReceived` agora é descartada em vez de encerrar o processo. Isso resolve a classe de falhas ao atravessar monitores com DPIs diferentes (#5544, #5650). Os caminhos de criação do ambiente/controlador continuam sendo fatais.\

**Diff completo:** https://github.com/wailsapp/wails/compare/webview2/v1.0.25...webview2/v1.0.26

## Corrigido

- Corrige o fluxo de trabalho release-webview2 para processar corretamente os arquivos go.sum no [PR](https://github.com/wailsapp/wails/pull/5671), por @taliesin-ai
- Corrige as atualizações de menu do GTK4 no Linux, limpando e reconstruindo o menu nativo no [PR](https://github.com/wailsapp/wails/pull/5659), por @taliesin-ai

## v3.0.0-alpha2.105 - 2026-06-21

## Adicionado

- Adiciona `application.System` para detectar a plataforma em tempo de execução a partir de código compartilhado: `System.IsMobile()` (iOS/Android), `System.IsDesktop()` (macOS/Windows/Linux), `System.IsServer()` (a tag de compilação `server`) e `System.IsPlatform(application.PlatformMacOS|PlatformWindows|PlatformLinux|PlatformIOS|PlatformAndroid|PlatformServer)` para testar diretamente um único destino. Ele é compilado em todos os destinos, permitindo criar ramificações condicionais sem tags de compilação. Os auxiliares correspondentes para o frontend (`System.IsMobile/IsDesktop/IsIOS/IsAndroid/...`) estão disponíveis em `@wailsio/runtime`
- Adiciona um guia "Como usar outros frameworks de frontend" que mostra como colocar seu próprio projeto Vite em `frontend/` (abrange Solid, Preact, Lit, SvelteKit, Qwik, Angular etc.)
- O assistente `wails3 setup` agora verifica o conjunto de ferramentas para dispositivos móveis (iOS/Android) — Xcode e o runtime do Simulador do iOS, JDK, Android SDK/NDK e emulador — com instalação em um clique e correções copiáveis para a configuração do shell, quando aplicável
- Os projetos gerados incluem uma configuração `frontend/.npmrc` que define um `minimum-release-age` de 7 dias para reduzir a exposição a pacotes recém-publicados (e potencialmente comprometidos), respeitada pelo pnpm e pelo bun e ignorada sem efeitos adversos pelo npm

## Alterado

- Redesenha todos os modelos iniciais integrados com um novo visual principal de montanhas em neon (Web, iOS e Android)
- **Agora, TypeScript é o padrão dos modelos iniciais e usa o nome do modelo sem sufixo.** `wails3 init` (sem `-t`) cria a estrutura de um projeto TypeScript; `-t vanilla`, `-t react`, `-t vue` e `-t svelte` usam TypeScript, com variantes JavaScript em `-t vanilla-js`, `-t react-js`, `-t vue-js` e `-t svelte-js`. Os modelos integrados declaram sua linguagem com `typescript:` em `template.yaml`; como alternativa, os modelos da comunidade que usam o sufixo `-ts` continuam funcionando
- Redesenha o assistente `wails3 setup` com o tema neon "digital Wails" (efeito vibrante de vidro fosco sobre um cenário de montanhas)

## Corrigido

- Corrige uma falha no Windows ao restaurar um aplicativo que permaneceu minimizado tempo suficiente para o WebView2 ser suspenso ou para seu processo de renderização/GPU ser reciclado. A nova sincronização de DPI ao minimizar/restaurar (#5544) agora só acessa o controlador do WebView2 quando o DPI da janela realmente mudou, evitando chamadas fatais de COM para um controlador suspenso no caso comum de restauração com o mesmo DPI (#5605)
- Corrige falhas nativas recorrentes de `SIGABRT`/`SIGSEGV` (geralmente dentro de `g_object_unref` durante o loop principal do GTK) em aplicativos Linux executados por longos períodos e submetidos a carregamentos frequentes de ativos/mídia. O servidor de ativos concluía `WebKitURISchemeRequest`s a partir de goroutines de trabalho, chamando funções do WebKit2GTK que não são thread-safe fora da thread principal do GTK; agora, a conclusão (`webkit_uri_scheme_request_finish_with_response`/`finish_error`) é executada na thread principal. Conclui a correção parcial de #5566. Afeta as compilações 6.0 com GTK3 e com GTK4/WebKitGTK (#5631, #5557)
- Corrige ocorrências intermitentes de `fatal error: invalid pointer found on stack` em `setupSignalHandlers` no Linux/GTK3. Os IDs de janela passados como `user_data` de sinal eram mantidos em uma variável local Go `unsafe.Pointer`; por isso, o coletor de lixo abortava ao examinar o valor (que não era um ponteiro) durante uma cópia da pilha. Agora, o ID é mantido como inteiro (`uintptr_t`) no lado Go, retroportando para o caminho legado do GTK3 a mesma correção aplicada por #4958 ao caminho do GTK4 (que alterou as funções de sinal em C para `uintptr_t` a fim de eliminar erros de `-race`/checkptr) (#5631)

## Removido

- Remove os modelos iniciais `react-swc`, `preact`, `lit`, `solid`, `qwik` e `sveltekit` (e suas variantes `-ts`). O conjunto integrado compatível agora é `vanilla`, `react`, `vue` e `svelte` — todos usam TypeScript por padrão e têm variantes JavaScript `-js`. Ainda é possível usar qualquer outro framework ao [trazer seu próprio frontend](https://v3.wails.io/guides/dev/frontend-frameworks) ou utilizar um modelo personalizado

## v3.0.0-alpha2.104 - 2026-06-18

## Corrigido

- Corrige uma falha (SIGABRT) no iOS quando um método vinculado de um serviço Go retorna uma string vazia. O gravador de resposta de ativos do iOS verificava o ponteiro do corpo com `buf != nil` em vez de verificar seu comprimento; por isso, um corpo de comprimento zero fazia `&buf[0]` entrar em pânico. Agora, ele verifica o comprimento, como os gravadores para desktop

## v3.0.0-alpha2.103 - 2026-06-15

## Alterado

- Move os recursos nativos do iOS e Android para gerenciadores de plataforma: chame-os por meio de `application.IOS.*` e `application.Android.*` (por exemplo, `application.IOS.Haptic("medium")`, `application.Android.Share(payload)`), em vez das antigas funções avulsas `application.IOS*`/`application.Android*` (#5602)
- Renomeia os eventos da ponte móvel: os eventos multiplataforma agora usam o prefixo `common:*` (por exemplo, `common:haptic` e `common:location`), e os eventos exclusivos de cada plataforma usam `ios:*` / `android:*` (por exemplo, `ios:backgroundTask` e `android:foregroundService`); o prefixo `native:*` não é mais usado (#5602)

## v3.0.0-alpha.102 - 2026-06-14

## Adicionado

- Adiciona o assistente experimental `wails3 setup` para configuração interativa de projetos e verificação de dependências
- Adiciona a opção `--json` a `wails3 doctor` para gerar uma saída legível por máquina
- Adiciona uma seção de status da assinatura ao comando `wails3 doctor`

## Corrigido

- Corrige a detecção do npm no Linux para verificar o PATH além do gerenciador de pacotes

## v3.0.0-alpha.101 - 2026-06-13

## Adicionado

- iOS: diálogos nativos de mensagem (UIAlertController) e diálogos para abrir arquivo, arquivos ou diretório (UIDocumentPickerViewController); diálogos para salvar retornam um erro explícito
- iOS: suporte à área de transferência por meio de UIPasteboard
- iOS: métricas reais da tela por meio de UIScreen (pontos, pixels, escala e área de trabalho dentro da área segura)
- iOS: builds para dispositivos (`IOS_PLATFORM=device`), suporte a identidade de assinatura de código, perfil de provisionamento e entitlements, empacotamento `.ipa` e `deploy-device` por meio de devicectl
- iOS: versão mínima configurável do iOS (`ios.minIOSVersion` em build/config.yml)
- iOS: `wails3 doctor` informa a disponibilidade do Xcode e do SDK do iOS no macOS
- iOS: eventos do sistema — bateria, rede, tema, bloqueio de tela e pouca memória são expostos como eventos de aplicativo `events.IOS.*` e `events.Common.*`, estes últimos independentes de plataforma
- iOS: ponte nativa de recursos móveis (`application.IOS*` exportado) — folha de compartilhamento, abertura de URL, manutenção da tela ativa, lanterna, margens da área segura, brilho, informações do aplicativo, bloqueio de orientação, barra de status, biometria (Face ID/Touch ID), notificações locais e armazenamento seguro no Keychain
- iOS: sensores e hardware — resposta tátil, geolocalização pontual, acelerômetro, proximidade, conversão de texto em fala, informações de armazenamento, estado de energia/bateria, status da rede, margens do teclado e detecção de captura de tela
- iOS: documentação (IOS.md e um guia no site da documentação)
- Android: diálogos nativos de mensagem (AlertDialog) e diálogos para abrir arquivo ou arquivos (Storage Access Framework, importados como cópias no cache); diálogos para abrir diretório e salvar retornam um erro explícito
- Android: suporte à área de transferência por meio de ClipboardManager
- Android: métricas reais da tela por meio de WindowMetrics/DisplayMetrics (dp, pixels, escala e área de trabalho delimitada pelas barras do sistema)
- Android: métodos de runtime para resposta tátil (`Android.Haptics.Vibrate`), informações do dispositivo (`Android.Device.Info`) e mensagens toast (`Android.Toast.Show`)
- Android: eventos de ciclo de vida tipados (`events.Android.*`, gerados a partir de events.txt), com `ActivityCreated` mapeado para `Common.ApplicationStarted`
- Android: o pipeline de build produz APKs de depuração e de lançamento instaláveis (`android:run`, `android:package` e `android:package:fat`); por padrão, o lançamento é assinado com o keystore de depuração, ou com um keystore real por meio das variáveis de ambiente `ANDROID_KEYSTORE_*`
- Android: `wails3 doctor` informa o Android SDK, o NDK e o JDK
- Android: eventos do sistema — bateria, rede, tema, bloqueio de tela e pouca memória são expostos como eventos de aplicativo `events.Android.*` e `events.Common.*`, estes últimos independentes de plataforma
- Android: ponte nativa de recursos móveis (`application.Android*` exportado) — compartilhamento, abertura de URL, manutenção da tela ativa, lanterna, margens da área segura, brilho, informações do aplicativo, bloqueio de orientação, barra de status, biometria (BiometricPrompt), notificações locais e armazenamento seguro em EncryptedSharedPreferences
- Android: sensores e hardware — resposta tátil, geolocalização pontual, acelerômetro, proximidade, conversão de texto em fala, informações de armazenamento, estado de energia/bateria, status da rede, margens do teclado e bloqueio de captura de tela com FLAG_SECURE
- Android: documentação (ANDROID.md e um guia no site da documentação)
- Exemplo: o kitchen sink `mobile` recebe as abas Mobile e Hardware, que demonstram a ponte nativa de recursos no iOS e no Android (as abas em formato de pílula quebram em várias linhas)
- Dispositivos móveis: bateria — o acelerômetro, o sensor de proximidade, a lanterna e o relógio periódico do exemplo são pausados quando o aplicativo passa para segundo plano e restaurados ao retornar (o Android mantém o processo em execução em segundo plano, e a lanterna é um estado de hardware que persiste no iOS); além disso, os receptores de eventos do sistema Android só são registrados enquanto o aplicativo está em primeiro plano
- iOS: captura da câmera — `application.IOSCapturePhoto`/`IOSCaptureVideo` (UIImagePickerController → um evento `native:capture` com uma miniatura em base64)
- iOS: execução em segundo plano — `application.IOSBeginBackgroundTask`/`IOSEndBackgroundTask` (uma janela de tarefa em segundo plano de UIApplication) e um `ios.backgroundModes` configurável (build/config.yml), que insere `UIBackgroundModes` como modelo no Info.plist gerado
- Android: captura da câmera — `application.AndroidCapturePhoto`/`AndroidCaptureVideo` (câmera do sistema por meio de FileProvider → um evento `native:capture`)
- Android: serviço em primeiro plano — `application.AndroidStartForegroundService`/`AndroidStopForegroundService` (um `WailsForegroundService` com uma notificação permanente mantém o processo ativo para trabalhos prolongados em segundo plano)
- Exemplo: uma aba Camera que demonstra a captura de fotos/vídeos e a execução em segundo plano (serviço em primeiro plano no Android e janela de tarefa em segundo plano no iOS)

## Corrigido

- Corrige `getUserMedia`, que sempre falhava com `NotAllowedError` no Linux: o WebKitGTK nega solicitações de permissão que não são tratadas, e o sinal `permission-request` não estava conectado. Agora, câmera e microfone são tratados por um novo mapa multiplataforma `WebviewWindowOptions.Permissions` (`map[PermissionType]Permission`), respeitado tanto no Linux (WebKitGTK) quanto no Windows (WebView2). No Linux, que não possui um prompt nativo, câmera e microfone são permitidos por padrão (restaurando `getUserMedia`) e podem ser desativados com `PermissionDeny` (#5552)
- iOS: `GOOS=ios` volta a compilar (`events.IOS` exportado e stubs de nomes de métodos móveis), e builds com a tag de produção compilam (correções de tags de build em pkg/application e em vários serviços)
- iOS: eventos Go→JS e ExecJS agora funcionam — a página não é mais carregada duas vezes na inicialização, e o handshake `wails:runtime:ready` não pode mais ser perdido
- iOS: `ApplicationDidFinishLaunching`/`ApplicationStarted` não entram mais em condição de corrida com a inicialização do aplicativo; a espera fixa de 2 segundos na inicialização foi removida
- iOS: corrige um vazamento de string C em cada execução de JavaScript de Go→JS
- iOS: `hasListeners` agora reflete o registro real de listeners
- iOS: o logging de depuração do framework é excluído da compilação dos builds de produção
- Android: `GOOS=android` volta a compilar — `events.Android` foi definido, o array de listeners `events_android.go` com acesso fora dos limites foi removido, o stub de nomes de métodos móveis foi adicionado, e arquivos do Linux para desktop (`linux_cgo.*`, `events_linux.*` e `environment_linux.go`) deixaram de vazar para os builds do Android
- Android: os bindings JS→Go agora funcionam — o WebView não consegue entregar corpos de requisições POST `fetch()` a `shouldInterceptRequest`, portanto as chamadas de runtime são encaminhadas por um transporte JavascriptInterface (`nativeHandleRuntimeCall`), em vez de causarem uma falha por corpo de requisição nil
- Android: as chamadas de runtime `Screens.*` retornam dados reais — o ScreenManager agora é preenchido na inicialização (ele nunca havia sido conectado, portanto `GetAll` retornava nil)
- Android: o logging de depuração do framework é excluído da compilação dos builds de produção e, nos builds de depuração, é encaminhado ao logcat com a tag `Wails`
- Android: registro `hasListeners` real, tratamento de referências/exceções JNI e ciclo de vida da página com um único carregamento (sem navegação duplicada)
- Corrige a falha de `wails3 generate bindings` com "Acesso negado" no Windows quando o servidor de desenvolvimento do Vite está em execução, sincronizando os arquivos gerados com o diretório de saída em vez de renomeá-los sobre ele (#5515)
- Corrige uma falha fatal intermitente no macOS ao ler informações da tela após uma alteração de monitor: o ID e o nome da tela armazenavam ponteiros para buffers `UTF8String` com autorelease, que podiam ser liberados antes de o Go copiá-los (uso após liberação). Agora, as strings passam por `strdup` e são liberadas após a conversão, e a enumeração de telas é executada em um pool de autorelease explícito, portanto não causa mais vazamento quando chamada por goroutines do Go (#5556)
- Corrige uma SIGSEGV intermitente no Linux quando o assetserver fecha um `WebKitURISchemeRequest`: o `g_object_unref` final era executado na goroutine do assetserver, finalizando um GObject do WebKit fora da thread principal do GTK. Agora, o unref é encaminhado ao contexto principal do GTK por meio de `g_main_context_invoke` (#5557)

## v3.0.0-alpha.100 - 2026-06-13

## Adicionado

- Estende `MacWebviewPreferences` com opções adicionais de configuração do WKWebView: `EnableAutoplayWithoutUserAction`, `AllowsAirPlayForMediaPlayback`, `AllowsMagnification`, `JavaScriptCanOpenWindowsAutomatically`, `MinimumFontSize` e `ApplicationNameForUserAgent` (#5549)

## Corrigido

- Corrige a falha de `wails3 generate bindings` com "Acesso negado" no Windows quando o servidor de desenvolvimento do Vite está em execução, sincronizando os arquivos gerados no diretório de saída em vez de renomeá-los sobre os arquivos existentes (#5561)
- Corrige eventos de redimensionamento do JS que não eram disparados em janelas sem moldura no Linux e corrige a detecção das bordas da barra de rolagem nessas janelas (#5368)
- Corrige a falha do atualizador no Windows com "link inválido entre dispositivos" quando o diretório temporário está em um volume diferente do diretório de instalação (#5560)

## v3.0.0-alpha.99 - 2026-06-10

## Corrigido

- Corrige a falha de `wails3 generate bindings` com "Acesso negado" no Windows quando o servidor de desenvolvimento do Vite está em execução, sincronizando os arquivos gerados no diretório de saída em vez de renomeá-los sobre os arquivos existentes (#5515)

## v3.0.0-alpha.98 - 2026-06-03

## Corrigido

- Corrige o congelamento da interface do WebKit no Linux quando ociosa (por exemplo, com o inspetor aberto), deixando de forçar `SA_ONSTACK` em `SIGUSR1`, o que interrompia a sincronização da thread de coleta de lixo do JavaScriptCore (#5527)

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

- Adiciona suporte à ofuscação com Garble ([#4563](https://github.com/wailsapp/wails/issues/4563)): IDs estáveis de métodos de binding, integração com o processo de build e o Taskfile (`build --obfuscated --garbleargs`, `generate bindings -obfuscated`) e tags de struct JSON em todos os payloads expostos ao runtime (`EnvironmentInfo`, `OSInfo`, `Screen`, `Rect`, `Point`, `Size`, `Capabilities`), para que o formato de transmissão continue funcionando após o Garble renomear campos exportados.

## v3.0.0-alpha.95 - 2026-05-20

## Adicionado

- Adiciona a página ausente sobre a estrutura do projeto

## Alterado

- Documentação: altera dois diagramas na página de arquitetura para usar diagramas de sequência e proporcionar uma exibição mais clara
- Documentação: inclui uma observação sobre a instalação do D2 como pré-requisito para a execução

## Corrigido

- Corrige `wails3 generate appimage` na configuração padrão com GTK4: agora, o empacotador detecta a pilha GTK no binário antes de procurar os arquivos de runtime, selecionando `libwebkitgtkinjectedbundle.so` (em `webkitgtk-6.0/`) para builds com GTK4 e `libwebkit2gtkinjectedbundle.so` (em `webkit2gtk-4.1/`) para builds com `-tags gtk3`. A verificação de `.relr.dyn` também consulta `libgtk-4.so.1`, portanto a remoção de símbolos é desativada corretamente em toolchains modernas, independentemente da pilha. (#5475)
- Corrige a falha de `wails3 generate appimage` quando invocado com um `-builddir` relativo: agora, o empacotador resolve antecipadamente `-binary`, `-icon`, `-desktopfile`, `-builddir` e `-outputdir` como caminhos absolutos, para que o `s.CD` executado no meio do fluxo não interrompa a goroutine de download do AppRun nem a verificação de `ldd` após a cópia.
- Corrige a falha de `wails3 generate appimage` ao mover a AppImage final para `-outputdir` quando o campo `Name=` do arquivo desktop não corresponde ao nome-base do binário: agora, por meio da variável de ambiente `OUTPUT`, o empacotador força o plugin appimage do linuxdeploy a gravar a AppImage em `<binary>-<arch>.AppImage`, em vez de usar o nome derivado do arquivo desktop.
- Corrige `events.Common.ApplicationStarted`, `Common.ThemeChanged`, `Common.SystemWillSleep` e `Common.SystemDidWake`, que não eram disparados no Linux após a pilha GTK4 + WebKitGTK 6.0 se tornar o padrão na alpha.93. O novo `application_linux.go` `run()` padrão não chamava `setupCommonEvents()`, que encaminha eventos `Linux.*` aos equivalentes `Common.*`, nem `monitorPowerEvents()`. Agora, o auxiliar de monitoramento de energia via DBus é compartilhado entre os caminhos de build com GTK3 e GTK4 por meio de `application_linux_dbus.go`. (#5474)

## v3.0.0-alpha.94 - 2026-05-19

## Corrigido

- Corrige `events.Common.ApplicationStarted`, `Common.ThemeChanged`, `Common.SystemWillSleep` e `Common.SystemDidWake`, que não eram disparados no Linux após a pilha GTK4 + WebKitGTK 6.0 se tornar o padrão na alpha.93. O novo `application_linux.go` `run()` padrão não chamava `setupCommonEvents()`, que encaminha eventos `Linux.*` aos equivalentes `Common.*`, nem `monitorPowerEvents()`. Agora, o auxiliar de monitoramento de energia via DBus é compartilhado entre os caminhos de build com GTK3 e GTK4 por meio de `application_linux_dbus.go`. (#5474)

## v3.0.0-alpha.93 - 2026-05-17

## Adicionado

- Adiciona `XDG_SESSION_TYPE` à saída de `wails3 doctor` no Linux, por @leaanthony

## Corrigido

- Corrige uma falha no menu da janela no Wayland, causada pelo acesso do appmenu-gtk-module a uma janela ainda não realizada (#4769), por @leaanthony
- Corrige uma falha do aplicativo GTK quando o nome do aplicativo contém caracteres inválidos (espaços, parênteses etc.), por @leaanthony
- Corrige o erro "memória insuficiente" ao inicializar o recurso de arrastar e soltar no Windows (#4701), por @overlordtm
- Corrige uma condição de corrida no armazenamento de callbacks da thread principal, que usava incorretamente RLock para excluir entradas do mapa (Linux, macOS, iOS) (#4424), por @leaanthony
- Corrige o tratamento de variáveis ao passar argumentos de linha de comando às tarefas. Agora, as variáveis da CLI especificadas como pares KEY=VALUE são inicializadas corretamente e propagadas durante toda a execução da tarefa.
- Corrige o conflito de NSWindowZoomButton no macOS: agora, `MaximiseButtonState` e `FullscreenButtonState` aplicam o estado mais restritivo tanto na inicialização quanto em tempo de execução; nenhum dos setters pode sobrescrever silenciosamente o outro (#5319)
- Corrige um conjunto de bugs preexistentes no caminho de build legado com GTK3 (`-tags gtk3`), revelados pelo CodeRabbit em #5463: inicializações por associação de arquivos não ignoram mais os manipuladores de inicialização; `getTheme` agora é seguro quanto aos limites e aos tipos; `appName` não libera mais memória pertencente ao GLib; `clipboardGet` não causa mais vazamento do `gchar*` retornado pelo GTK; agora, `Calloc` usa receivers de ponteiro, e `NewCalloc` retorna `*Calloc`, para que o pool realmente rastreie as alocações; `zoomOut` usa o inverso de `zoomInFactor`, em vez de um multiplicador negativo limitado a 1.0; `execJS` reutiliza o nome de mundo vazio pré-alocado, em vez de causar o vazamento de um `C.CString("")` a cada chamada; e um `fmt.Println` de desenvolvimento foi removido de `menuItem.setAccelerator`. Resolve #5465.
- Corrige o mesmo vazamento do receptor por valor de `Calloc` no caminho de compilação padrão do GTK4 (`linux_cgo.go`): receptores de ponteiro + `NewCalloc() *Calloc`, para que as alocações de `c.String(...)` de cada janela sejam de fato rastreadas e liberadas.

## v3.0.0-alpha.92 - 2026-05-15

## Adicionado

- Modifica os Taskfiles para permitir controlar o gerenciador de pacotes do frontend usado por meio da opção `PACKAGE_MANAGER`
- Enriquece os dados de modelo com `{{.Opn}}` e `{{.Cls}}` para tornar a criação de modelos de Taskfile mais previsível

## Alterado

- Modifica alguns dos Taskfiles existentes para usar `{{.Opn}} and {{.Cls}}`

## Corrigido

- Corrige uma falha fatal em tempo de execução de `concurrent map read and map write` em `linuxSystemTray` quando o menu da bandeja é atualizado enquanto o painel o lê.
- Usa `log` em vez de `fmt` para a saída de erros e rastreamentos de pilha do WebView2, evitando que as mensagens sejam perdidas quando o aplicativo é executado sem um console anexado no Windows.

## v3.0.0-alpha.91 - 2026-05-12

## Alterado

- Atualiza o SVG de patrocinadores no [PR](https://github.com/wailsapp/wails/pull/5414) por `@github-actions[bot]`
- **INCOMPATÍVEL (macOS):** Normaliza o sistema de coordenadas do macOS para que `GetScreens`, `Position` e `SetPosition` usem o mesmo espaço — pontos lógicos, Y crescente para baixo, com `(0,0)` no canto superior esquerdo da tela principal. Isso corresponde ao Windows, ao GTK e às APIs públicas do Electron e da web. Telas fisicamente acima da principal agora informam `Bounds.Y` negativo (antes positivo), e os valores de `Position()`/`SetPosition()` agora são expressos em pontos lógicos em vez de `points × primaryScale`. A conversão de ida e volta `Position()` → `SetPosition()` é preservada; valores absolutos registrados por compilações alfa anteriores ou soluções alternativas calculadas manualmente (por exemplo, multiplicar por `primaryScale` ou inverter Y em relação à altura da tela) precisarão ser atualizados. Resolve [#5117](https://github.com/wailsapp/wails/issues/5117).

## Corrigido

- Valida preventivamente o nome do sinal DBus e o comprimento do corpo para evitar panics no [PR](https://github.com/wailsapp/wails/pull/5416) por @leaanthony
- Corrige um problema de segurança de memória no tratamento de menus do GTK no Linux no [PR](https://github.com/wailsapp/wails/pull/5363) por @leaanthony
- Detecta GPUs NVIDIA e desativa o renderizador DMA-BUF no Linux no [PR](https://github.com/wailsapp/wails/pull/5295) por @leaanthony
- Corrige a conversão de Y entre telas de `SetPosition` no macOS: usa a altura da tela principal como referência global para que as janelas sejam posicionadas corretamente em monitores deslocados verticalmente em relação à tela principal em [#5117](https://github.com/wailsapp/wails/issues/5117)
- Corrige o modelo de PR do Git para apontar para a URL de feedback correta no [PR](https://github.com/wailsapp/wails/pull/5109) por @wayneforrest
- Corrige uma série de falhas de `SetMenu` na bandeja do sistema do Windows causadas por uma chamada de sistema `DestroyMenu` incorreta, que passava quatro argumentos em vez de um, fazendo todas as chamadas retornarem FALSE e não liberarem nada. Também libera identificadores HMENU e HBITMAP (inclusive os alocados em tempo de execução por meio de `MenuItem.SetBitmap`) nas reconstruções do menu, redefine mapas obsoletos de caixas de seleção e botões de opção em `Win32Menu.Update` e remove uma chamada redundante a `Update()` em `systemtray.updateMenu` que duplicava as alocações. Aplicativos de longa duração com bandeja do sistema não vazam mais objetos GDI/USER a cada reconstrução do menu.

## v3.0.0-alpha.90 - 2026-05-11

## Adicionado

- Adiciona um nome de aplicativo configurável ao User-Agent do WKWebView no macOS no [PR](https://github.com/wailsapp/wails/pull/5261) por @vinhvoit225
- Adiciona a dependência indireta github.com/coder/websocket ao exemplo gin-service no [PR](https://github.com/wailsapp/wails/pull/5400) por @taliesin-ai
- Adiciona suporte à comparação de igualdade profunda nos testes de ativos de compilação no [PR](https://github.com/wailsapp/wails/pull/5402) por @leaanthony

## Alterado

- Consolida a saída da compilação no diretório de ativos no [PR](https://github.com/wailsapp/wails/pull/5401) por @taliesin-ai
- Atualiza o SVG de patrocinadores no [PR](https://github.com/wailsapp/wails/pull/5399) por `@github-actions[bot]`

## Corrigido

- Usa um objeto de notificação para a mensagem de instância única no macOS no [PR](https://github.com/wailsapp/wails/pull/5289) por @overlordtm
- Agrupa callbacks do Windows em lotes para evitar a perda de promises sob carga intensa no [PR](https://github.com/wailsapp/wails/pull/5383) por @taliesin-ai

## v3.0.0-alpha.89 - 2026-05-10

## Adicionado

- Adiciona o job go<em>test</em>results para agregar resultados de testes do Go no [PR](https://github.com/wailsapp/wails/pull/5316) por @leaanthony

## Alterado

- Divide condicionalmente payloads RPC grandes em requisições POST fragmentadas no [PR](https://github.com/wailsapp/wails/pull/5369) por @leaanthony
- Atualiza o Vite da versão 5.x.x para 8.0.0 em todos os modelos de frontend no [PR](https://github.com/wailsapp/wails/pull/5386) por @leaanthony
- Migra a configuração da porta do servidor de desenvolvimento do Vite para variáveis de ambiente no [PR](https://github.com/wailsapp/wails/pull/5365) por @leaanthony
- Configura o servidor de desenvolvimento do Vite para escutar em 127.0.0.1 em todos os modelos no [PR](https://github.com/wailsapp/wails/pull/5361) por @leaanthony
- Atualiza o SVG de patrocinadores no [PR](https://github.com/wailsapp/wails/pull/5384) por `@github-actions[bot]`

## Corrigido

- Sanitiza os stubs do modelo Info.plist durante a atualização de build-assets no [PR](https://github.com/wailsapp/wails/pull/5312) por @leaanthony
- Corrige o estado obsoleto nos menus do macOS aplicando os métodos modificadores dos itens de menu (`setMenuItemChecked()`, `setMenuItemLabel()`, `setMenuItemDisabled()`, `setMenuItemHidden()`, `setMenuItemTooltip()`) de forma síncrona na thread principal, eliminando a condição de corrida de `dispatch_async` que fazia os menus exibirem o estado anterior quando eram reabertos rapidamente (#5002)
- Ignora arquivos `*_test.go` no modo de desenvolvimento para evitar recompilações desnecessárias no [PR](https://github.com/wailsapp/wails/pull/5203) por @leaanthony
- Evita falha de segmentação em Menu.Update() quando o aplicativo não está em execução no [PR](https://github.com/wailsapp/wails/pull/5291) por @wucm667
- Usa lastSizeWParam para condicionar o redesenho da barra de menus no Windows no [PR](https://github.com/wailsapp/wails/pull/5382) por @taliesin-ai

## v3.0.0-alpha.88 - 2026-05-09

## Alterado

- Altera HiddenOnTaskbar para usar WS<em>EX</em>TOOLWINDOW no [PR](https://github.com/wailsapp/wails/pull/5371) por @leaanthony
- Reordena as dependências e remove a diretiva replace de webview2 do go.mod no [PR](https://github.com/wailsapp/wails/pull/5370) por @atterpac
- Atualiza o SVG de patrocinadores no [PR](https://github.com/wailsapp/wails/pull/5358) por `@github-actions[bot]`

## Corrigido

- Remove aliases genéricos de indireção e consolida os tipos de chave de mapas no [PR](https://github.com/wailsapp/wails/pull/5331) por @fbbdev

## Removido

- Exclui o fluxo de trabalho PR-master, removendo a documentação, os testes do Go e os testes de omissão no [PR](https://github.com/wailsapp/wails/pull/5377) por @leaanthony

## v3.0.0-alpha.87 - 2026-05-07

## Adicionado

- Adiciona documentação em coreano para o Wails v3 no [PR](https://github.com/wailsapp/wails/pull/5352) por @leaanthony
- Adiciona documentação em francês sobre instalação e início rápido no [PR](https://github.com/wailsapp/wails/pull/5354) por @leaanthony
- Adiciona documentação em português sobre início rápido, conceitos e comunidade no [PR](https://github.com/wailsapp/wails/pull/5355) por @leaanthony

## v3.0.0-alpha.86 - 2026-05-06

## Adicionado

- Adiciona a localização da documentação em francês no [PR](https://github.com/wailsapp/wails/pull/5328) por @leaanthony
- Adiciona o idioma alemão ao site da documentação no [PR](https://github.com/wailsapp/wails/pull/5343) por @leaanthony

## Alterado

- Registra todos os idiomas traduzidos de 8 na configuração da documentação no [PR](https://github.com/wailsapp/wails/pull/5347) por @leaanthony
- Atualiza vários arquivos relacionados ao Windows para o WebView2 no [PR](https://github.com/wailsapp/wails/pull/5317) por @leaanthony

## Corrigido

- Separa o despacho de diálogos entre GTK3 e GTK4 no Linux no [PR](https://github.com/wailsapp/wails/pull/5340) por @leaanthony
- Garante que os callbacks de diálogos sejam executados na thread do GTK, corrigindo falhas de segmentação no [PR](https://github.com/wailsapp/wails/pull/5339) por @leaanthony

## v3.0.0-alpha.85 - 2026-05-05

## Adicionado

- Adiciona ao repositório a URL do modelo de PR no [PR](https://github.com/wailsapp/wails/pull/5179) por @leaanthony
- Adiciona a documentação em alemão do Wails v3 no [PR](https://github.com/wailsapp/wails/pull/5330) por @leaanthony

## v3.0.0-alpha.84 - 2026-05-03

## Adicionado

- Adiciona uma opção para impedir que a tecla Escape saia do modo de tela cheia no macOS no [PR](https://github.com/wailsapp/wails/pull/5307) por @leaanthony
- Adiciona uma opção para impedir que a tecla Escape saia do modo de tela cheia no macOS no [PR](https://github.com/wailsapp/wails/pull/5310) por @leaanthony
- Adiciona a documentação do Pausa à vitrine da comunidade no [PR](https://github.com/wailsapp/wails/pull/5288) por @yuseferi

## Alterado

- Atualiza o SVG dos patrocinadores no [PR](https://github.com/wailsapp/wails/pull/5308) por `@github-actions[bot]`
- Atualiza o comando de geração de ícones para lidar com plataformas sem suporte no [PR](https://github.com/wailsapp/wails/pull/5309) por @leaanthony
- Substitui a API booleana de tela cheia por ButtonState de três estados e implementa os bindings de plataforma no [PR](https://github.com/wailsapp/wails/pull/5224) por @leaanthony

## Corrigido

- Protege as operações de foco do WebView2 contra o estado nulo do controlador no [PR](https://github.com/wailsapp/wails/pull/5315) por @leaanthony
- Atualiza o workflow do GitHub Actions para referenciar corretamente a branch de destino do PR no [PR](https://github.com/wailsapp/wails/pull/5313) por @leaanthony
- Ignora arquivos `*_test.go` no modo de desenvolvimento para evitar recompilações desnecessárias no [PR](https://github.com/wailsapp/wails/pull/5203) por @leaanthony
- Evita uma falha de segmentação em Menu.Update() quando o aplicativo não está em execução no [PR](https://github.com/wailsapp/wails/pull/5291) por @wucm667

## v3.0.0-alpha.83 - 2026-05-02

## Adicionado

- Adiciona a flag InstallScope e uma opção de compilação para instalação por máquina/usuário no [PR](https://github.com/wailsapp/wails/pull/5094) por @symball
- Adiciona ao BrowserWindow um método SetScreen sem operação para satisfazer a interface Window no [PR](https://github.com/wailsapp/wails/pull/5294) por @leaanthony

## Corrigido

- Detecta GPUs NVIDIA e desativa o renderizador DMA-BUF no Linux no [PR](https://github.com/wailsapp/wails/pull/5295) por @leaanthony
- Corrige o modelo de PR do Git para apontar para a URL correta de feedback no [PR](https://github.com/wailsapp/wails/pull/5109) por @wayneforrest
- Corrige uma série de falhas de `SetMenu` na bandeja do sistema do Windows causadas por uma syscall `DestroyMenu` defeituosa que passava quatro argumentos em vez de um, fazendo com que todas as chamadas retornassem FALSE e não liberassem nada. Também libera identificadores HMENU e HBITMAP — incluindo os alocados em tempo de execução por meio de `MenuItem.SetBitmap` — durante a reconstrução de menus, redefine mapas obsoletos de caixas de seleção e botões de opção em `Win32Menu.Update` e remove uma chamada redundante a `Update()` em `systemtray.updateMenu` que duplicava as alocações. Aplicativos de bandeja do sistema executados por longos períodos não vazam mais objetos GDI/USER a cada reconstrução do menu.

## v3.0.0-alpha.82 - 2026-05-01

## Corrigido

- Corrige a geração do arquivo desktop para processar corretamente o nome do desktop no [PR](https://github.com/wailsapp/wails/pull/5232) por @leaanthony

## v3.0.0-alpha.81 - 2026-04-30

## Alterado

- Ajusta a programação das versões nightly para 15:00 UTC no [PR](https://github.com/wailsapp/wails/pull/5286) por @leaanthony

## Corrigido

- Corrige os valores de Screen Bounds, WorkArea e Size reduzidos pela metade em Macs com tela Retina -  (#5168)

## v3.0.0-alpha.80 - 2026-04-29

## Alterado

- Atualiza as dependências da documentação e os carregadores de coleções de conteúdo no [PR](https://github.com/wailsapp/wails/pull/5285) por @leaanthony

## v3.0.0-alpha.79 - 2026-04-29

## Adicionado

- Concede a permissão actions: write ao job trigger-release no [PR](https://github.com/wailsapp/wails/pull/5270) por @leaanthony

## Alterado

- Define a branch master como padrão da tarefa de lançamento e atualiza o texto do changelog no [PR](https://github.com/wailsapp/wails/pull/5283) por @leaanthony
- Atualiza o workflow auto-changelog para usar a versão mais recente no [PR](https://github.com/wailsapp/wails/pull/5282) por @leaanthony
- Melhora a eficiência dos workflows adicionando filtros de caminhos e removendo workflows inativos no [PR](https://github.com/wailsapp/wails/pull/5280) por @leaanthony
- Atualiza a documentação para referenciar a branch master nos links de exemplos no [PR](https://github.com/wailsapp/wails/pull/5274) por @leaanthony
- Atualiza a documentação e os exemplos de v3 no [PR](https://github.com/wailsapp/wails/pull/5272) por @leaanthony

## Corrigido

- Aprimora o proxy reverso com lógica de repetição e uso forçado de IPv4 durante o desenvolvimento no [PR](https://github.com/wailsapp/wails/pull/5265) por @AkagiYui
- Reescreve o workflow que aciona o changelog de alterações ainda não lançadas no [PR](https://github.com/wailsapp/wails/pull/5281) por @leaanthony

## Removido

- Remove scripts de teste em shell usados para várias finalidades de teste no [PR](https://github.com/wailsapp/wails/pull/5267) por @leaanthony
- Exclui o workflow de implantação da documentação v3-alpha e o registro CNAME no [PR](https://github.com/wailsapp/wails/pull/5266) por @leaanthony

### Adicionado

- Adiciona a entrada Roteamento do frontend à navegação da barra lateral no [PR](https://github.com/wailsapp/wails/pull/5196) por @leaanthony
- Adiciona um guia de roteamento do frontend com recomendações específicas para cada framework no [PR](https://github.com/wailsapp/wails/pull/5185) por @leaanthony
- Adiciona suporte a folhas modais (macOS)
- Atualiza a versão do ghw para oferecer melhor suporte a dispositivos Apple por @leaanthony (#4977)
- Adiciona o método `GetBadge` ao serviço de dock
- Adiciona a flag `-tags` ao comando `wails3 build` para passar tags de compilação personalizadas do Go (por exemplo, `wails3 build -tags gtk4`) (#4957)
- Adiciona documentação sobre a geração automática de enums no gerador de bindings, incluindo uma página dedicada a Enums e sua entrada na navegação da barra lateral (#4972)
- Adiciona a flag `-tags` ao comando `wails3 build` para passar tags de compilação personalizadas do Go (por exemplo, `wails3 build -tags gtk4`) (#4957)
- Adicionar exemplos de APIs Web em `v3/examples/web-apis/` que demonstram APIs de navegador de 41, incluindo armazenamento (localStorage, sessionStorage, IndexedDB e Cache API), rede (Fetch, WebSocket, XMLHttpRequest, EventSource e Beacon), mídia (Canvas, WebGL, Web Audio, MediaDevices, MediaRecorder e Speech Synthesis), dispositivo (Geolocation, Clipboard, Fullscreen, Device Orientation, Vibration e Gamepad), desempenho (Performance API, Mutation Observer e Intersection/Resize Observer), interface (Web Components, Pointer Events, Selection, Dialog e Drag and Drop), entre outras
- Adicionar um exemplo de verificador de compatibilidade de APIs do WebView (`v3/examples/webview-api-check/`) que testa mais de 200 APIs de navegador em diferentes plataformas
- Adicionar o pacote `internal/libpath` para localizar caminhos de bibliotecas nativas no Linux, com pesquisa paralela, cache e suporte a Flatpak/Snap/Nix
- **Em desenvolvimento:** adicionar suporte experimental a WebKitGTK 6.0 / GTK4 no Linux, disponível por meio de `-tags gtk4` (GTK3/WebKit2GTK 4.1 continua sendo o padrão)
- Observação: em gerenciadores de janelas lado a lado (por exemplo, Hyprland e Sway), as operações de minimizar/maximizar podem não funcionar conforme esperado, pois o gerenciador de janelas controla a geometria das janelas
- Adicionar à documentação como criar **manipuladores de uso único** em **Escuta de eventos em JavaScript**, por @AbdelhadiSeddar
- Adicionar a opção `UseApplicationMenu` a `WebviewWindowOptions`, permitindo que as janelas no Windows/Linux herdem o menu do aplicativo definido por meio de `app.Menu.Set()`, por @leaanthony
- Adicionar suporte ao uso de arquivos `.icon` (formato Apple Icon Composer) para gerar ícones Liquid Glass e catálogos de recursos (macOS) (#4934), por @wimaha
- Adicionar um modo de servidor experimental para implantações sem interface gráfica/na Web (`-tags server`). Permite executar aplicativos Wails como servidores HTTP sem dependências de interface gráfica nativa. Compile com `wails3 task build:server`. Consulte `examples/server` para obter detalhes.
- Adicionar o pacote `internal/libpath` para localizar caminhos de bibliotecas nativas no Linux, com pesquisa paralela, cache e suporte a Flatpak/Snap/Nix
- Adicionar a opção `CollectionBehavior` a `MacWindow` para controlar o comportamento das janelas entre Spaces do macOS e o modo de tela cheia (#4756), por @leaanthony
- Adicionar testes unitários para pkg/application, por @leaanthony
- Adicionar suporte a protocolos personalizados ao empacotamento MSIX, por @leaanthony
- Adicionar detecção do ambiente de desktop no Linux [PR #4797](https://github.com/wailsapp/wails/pull/4797)
- Adicionar o método `Window.Print()` ao runtime JavaScript para abrir a caixa de diálogo de impressão pelo frontend (#4290), por @leaanthony
- Adicionar `XDG_SESSION_TYPE` à saída de `wails3 doctor` no Linux, por @leaanthony
- Adicionar outros eventos de alteração de carregamento do WebKit2 no Linux: `WindowLoadStarted`, `WindowLoadRedirected`, `WindowLoadCommitted` e `WindowLoadFinished` (#3896), por @leaanthony
- Adicionar `XDG_SESSION_TYPE` à saída de `wails3 doctor` no Linux, por @leaanthony
- Gerar o arquivo `.desktop` durante a compilação para Linux, não apenas durante o empacotamento (#4575)
- Adicionar documentação sobre dependências de runtime no Linux, com nomes de pacotes específicos por distribuição e exemplos de empacotamento com nfpm (#4339), por @leaanthony
- Adicionar informações sobre a versão do driver NVIDIA à saída de `wails3 doctor` no Linux, por @leaanthony
- Adicionar a origem ao manipulador de mensagens brutas, por @APshenkin, no [PR](https://github.com/wailsapp/wails/pull/4710)
- Adicionar suporte a links universais no macOS, por @APshenkin, no [PR](https://github.com/wailsapp/wails/pull/4712)
- Refatorar a camada de transporte dos bindings, por @APshenkin, no [PR](https://github.com/wailsapp/wails/pull/4702)
- Adicionar identificadores aria-label aos modelos helloworld para que o aplicativo de exemplo possa ser testado facilmente por clientes de teste Appium, por @chinenual, no [PR](https://github.com/wailsapp/wails/pull/4760)
- Adicionar a origem ao manipulador de mensagens brutas, por @APshenkin, no [PR](https://github.com/wailsapp/wails/pull/4710)
- Adicionar suporte a links universais no macOS, por @APshenkin, no [PR](https://github.com/wailsapp/wails/pull/4712)
- Refatorar a camada de transporte dos bindings, por @APshenkin, no [PR](https://github.com/wailsapp/wails/pull/4702)
- Eventos tipados, por @fbbdev e @ianvs, em [#4633](https://github.com/wailsapp/wails/pull/4633)
- Adicionar o exemplo `systray-clock`, que mostra um ícone de bandeja sem interface gráfica com atualizações dinâmicas da dica de ferramenta (#4653).
- Adicionar um modelo de protocolo NSIS para Windows, por @Tolfx, em #4510
- Adicionar testes para build-assets, por @Tolfx, em #4510
- macOS: exibir controles nativos de janela na barra de menus em [#4588](https://github.com/wailsapp/wails/pull/4588), por @nidib
- Adicionar um serviço do Dock do macOS para ocultar/exibir o ícone do aplicativo no Dock, por @popaprozac, no [PR](https://github.com/wailsapp/wails/pull/4451)
- Adicionar um serviço do Dock do macOS para ocultar/exibir o ícone do aplicativo no Dock, por @popaprozac, no [PR](https://github.com/wailsapp/wails/pull/4451)
- Adicionar suporte nativo ao efeito Liquid Glass no macOS com NSGlassEffectView (macOS 15.0+) e fallback para NSVisualEffectView, incluindo opções abrangentes de personalização de materiais, por @leaanthony, em [#4534](https://github.com/wailsapp/wails/pull/4534)
- Sanitização de URLs do navegador, por @leaanthony, em [#4500](https://github.dev/wailsapp/wails/pull/4500). Com base em [#4484](https://github.com/wailsapp/wails/pull/4484), por @APShenkin.
- Adicionar proteção de conteúdo no Windows/Mac, por [@leaanthony](https://github.com/leaanthony), com base no trabalho original de [@Taiterbase](https://github.com/Taiterbase) neste [PR](https://github.com/wailsapp/wails/pull/4241)
- Adicionar suporte à passagem de variáveis da CLI para comandos Task por meio dos aliases `wails3 build` e `wails3 package` (#4422), por @leaanthony, no [PR](https://github.com/wailsapp/wails/pull/4488)
- Suporte a zonas de soltura cujo evento fornece os dados do elemento solto, por [@atterpac](https://github.com/atterpac), em [#4318](https://github.com/wailsapp/wails/pull/4318)
- Adicionar `AdditionalLaunchArgs` às opções de `WindowsWindow` para permitir que outros argumentos de linha de comando sejam passados ao navegador WebView2, no [PR](https://github.com/wailsapp/wails/pull/4467)
- Adicionar a execução automática de go mod tidy após wails init, por [@triadmoko](https://github.com/triadmoko), no [PR](https://github.com/wailsapp/wails/pull/4286)
- Recurso Snap Assist do Windows, por @leaanthony, no [PR](https://github.dev/wailsapp/wails/pull/4463)
- Adicionar `AdditionalLaunchArgs` às opções de `WindowsWindow` para permitir que outros argumentos de linha de comando sejam passados ao navegador WebView2, no [PR](https://github.com/wailsapp/wails/pull/4467)
- Adicionar a execução automática de go mod tidy após wails init, por [@triadmoko](https://github.com/triadmoko), no [PR](https://github.com/wailsapp/wails/pull/4286)
- Recurso Snap Assist do Windows, por @leaanthony, no [PR](https://github.dev/wailsapp/wails/pull/4463)
- Adicionar a implementação de `getAccentColor` para Windows, por [@almas-x](https://github.com/almas-x), no [PR](https://github.com/wailsapp/wails/pull/4427)
- Adicionar a implementação de `getAccentColor` para Windows, por [@almas-x](https://github.com/almas-x), no [PR](https://github.com/wailsapp/wails/pull/4427)
- Menus e barra de menus com tema escuro no Windows. Por @leaanthony, em [a29b4f0861b1d0a700e9eb213c6f1076ec40efd5](https://github.com/wailsapp/wails/commit/a29b4f0861b1d0a700e9eb213c6f1076ec40efd5)
- Renomear os serviços integrados para tornar os bindings JS/TS mais claros, por @popaprozac, no [PR](https://github.com/wailsapp/wails/pull/4405)
- `app.Env.GetAccentColor` para obter a cor de destaque do sistema do usuário. Funciona no macOS. Por [@etesam913](https://github.com/etesam913)
- Adicionar a API `window.ToggleFrameless()`, por [@atterpac](https://github.com/atterpac), em [#4137](https://github.com/wailsapp/wails/pull/4137)
- Adicionar dependências de compilação específicas por distribuição para Linux, por @leaanthony, no [PR](https://github.com/wailsapp/wails/pull/4345)
- Adicionado guia de bindings por @atterpac no [PR](https://github.com/wailsapp/wails/pull/4404)
- **Infraestrutura de testes organizada**: arquivos de teste do Docker movidos para o diretório dedicado `test/docker/`, com imagens otimizadas e maior confiabilidade de build, por [@leaanthony](https://github.com/leaanthony) em [#4359](https://github.com/wailsapp/wails/pull/4359)
- **Padrões aprimorados de gerenciamento de recursos**: adicionadas aos exemplos a limpeza adequada de manipuladores de eventos e a gestão de goroutines com reconhecimento de contexto, por [@leaanthony](https://github.com/leaanthony) em [#4359](https://github.com/wailsapp/wails/pull/4359)
- Adicionado suporte a builds de AppImage para aarch64 por [@AkshayKalose](https://github.com/AkshayKalose) em [#3981](https://github.com/wailsapp/wails/pull/3981)
- Adicionada uma seção de diagnóstico a `wails doctor` por [@leaanthony](https://github.com/leaanthony)
- Adicionada a janela ao contexto ao chamar um método de serviço, por [@leaanthony](https://github.com/leaanthony)
- Adicionado o exemplo `window-call` para demonstrar como identificar qual janela está chamando um serviço, por [@leaanthony](https://github.com/leaanthony)
- Novo guia de menus por [@leaanthony](https://github.com/leaanthony)
- Tratamento de pânicos aprimorado por [@leaanthony](https://github.com/leaanthony)
- Novo guia de menus por [@leaanthony](https://github.com/leaanthony)
- Adicionados comentários de documentação à API de serviços por [@fbbdev](https://github.com/fbbdev) em [#4024](https://github.com/wailsapp/wails/pull/4024)
- Adicionada a função `application.NewServiceWithOptions` para inicializar serviços com configuração adicional, por [@leaanthony](https://github.com/leaanthony) em [#4024](https://github.com/wailsapp/wails/pull/4024)
- Controle de menus aprimorado por [@FalcoG](https://github.com/FalcoG) e [@leaanthony](https://github.com/leaanthony) em [#4031](https://github.com/wailsapp/wails/pull/4031)
- Mais documentação por [@leaanthony](https://github.com/leaanthony)
- Adicionado suporte ao cancelamento de eventos em listeners de eventos padrão, por [@leaanthony](https://github.com/leaanthony)
- Adicionado suporte a `Hide`, `Show` e `Destroy` na bandeja do sistema, por [@leaanthony](https://github.com/leaanthony)
- Adicionado suporte a `SetTooltip` na bandeja do sistema, por [@leaanthony](https://github.com/leaanthony). Ideia original de [@lujihong](https://github.com/wailsapp/wails/issues/3487#issuecomment-2633242304)
- Adicionado o caminho do pacote aos avisos do gerador de bindings sobre tipos não compatíveis, por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Adicionado ao gerador de bindings suporte a aliases genéricos, por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Adicionado ao gerador de bindings suporte à flag JSON `omitzero`, por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Adicionada a diretiva `//wails:ignore` para impedir a geração de bindings para métodos de serviço selecionados, por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Adicionada a diretiva `//wails:internal` a serviços e modelos para permitir tipos exportados em Go, mas não em JS/TS, por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Adicionado ao gerador de bindings suporte a constantes de tipos alias, permitindo enums com tipagem fraca, por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Adicionados testes do gerador de bindings para recursos do Go 1.24, por [@fbbdev](https://github.com/fbbdev) em [#4068](https://github.com/wailsapp/wails/pull/4068)
- Adicionado suporte ao macOS 15 "Sequoia" em `OSInfo.Branding` para aprimorar a detecção da versão do sistema operacional em [#4065](https://github.com/wailsapp/wails/pull/4065)
- Adicionado o hook `PostShutdown` para executar código personalizado após a conclusão do processo de encerramento, por [@fbbdev](https://github.com/fbbdev) em [#4066](https://github.com/wailsapp/wails/pull/4066)
- Adicionada a struct `FatalError` para permitir a detecção de erros fatais em manipuladores de erros personalizados, por [@fbbdev](https://github.com/fbbdev) em [#4066](https://github.com/wailsapp/wails/pull/4066)
- Padronizada e documentada a ordem de inicialização e encerramento dos serviços, por [@fbbdev](https://github.com/fbbdev) em [#4066](https://github.com/wailsapp/wails/pull/4066)
- Adicionado um ambiente de testes para a sequência de inicialização e encerramento do aplicativo e para os testes de inicialização e encerramento dos serviços, por [@fbbdev](https://github.com/fbbdev) em [#4066](https://github.com/wailsapp/wails/pull/4066)
- Adicionado o método `RegisterService` para registrar serviços depois que o aplicativo tiver sido criado, por [@fbbdev](https://github.com/fbbdev) em [#4066](https://github.com/wailsapp/wails/pull/4066)
- Adicionado o campo `MarshalError` às opções do aplicativo e dos serviços para tratamento personalizado de erros em chamadas de bindings, por [@fbbdev](https://github.com/fbbdev) em [#4066](https://github.com/wailsapp/wails/pull/4066)
- Adicionado um wrapper de promise cancelável que propaga solicitações de cancelamento por cadeias de promises, por [@fbbdev](https://github.com/fbbdev) em [#4100](https://github.com/wailsapp/wails/pull/4100)
- Adicionada a capacidade de vincular o cancelamento de chamadas de bindings a um `AbortSignal`, por [@fbbdev](https://github.com/fbbdev) em [#4100](https://github.com/wailsapp/wails/pull/4100)
- Adicionado suporte a atributos `data-wml-*` para WML, além dos atributos `wml-*` habituais, por [@leaanthony](https://github.com/leaanthony)
- Adicionado o método `Configure` a todos os serviços para configuração tardia ou reconfiguração dinâmica, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Quando não configurado, o serviço `fileserver` envia uma resposta 503 Serviço indisponível, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Quando não configurado, o serviço `kvstore` fornece por padrão um armazenamento de chave-valor em memória, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Adicionado o método `Load` ao serviço `kvstore` para recarregar os dados do arquivo após alterações na configuração, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Adicionado o método `Clear` ao serviço `kvstore` para excluir todas as chaves, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Adicionado o tipo `Level` ao serviço `log` para fornecer constantes de nível de log no lado do JS, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Adicionado o método `Log` ao serviço `log` para especificar dinamicamente o nível de log, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Quando não configurado, o serviço `sqlite` fornece por padrão um banco de dados em memória, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Adicionado o método `Close` ao serviço `sqlite` para fechar manualmente o banco de dados, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Adicionado suporte a cancelamento nos métodos de consulta do serviço `sqlite`, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Adicionado suporte a instruções preparadas no serviço `sqlite` com bindings para JS, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Suporte ao Gin por [Lea Anthony](https://github.com/leaanthony) no [PR](https://github.com/wailsapp/wails/pull/3537), com base no trabalho original de [@AnalogJ](https://github.com/AnalogJ) neste [PR](https://github.com/wailsapp/wails/pull/3537)
- Corrige o salvamento automático e o salvamento automático de senhas, que estavam sempre ativados, por [@oSethoum](https://github.com/osethoum) em [#4134](https://github.com/wailsapp/wails/pull/4134)
- Adiciona `SetMenu()` à janela para permitir a definição de um menu nela, por [@leaanthony](https://github.com/leaanthony)
- Adiciona suporte a notificações por [@popaprozac](https://github.com/popaprozac) em [#4098](https://github.com/wailsapp/wails/pull/4098)
-  Adiciona suporte a associações de arquivos no Mac por [@wimaha](https://github.com/wimaha) em [#4177](https://github.com/wailsapp/wails/pull/4177)
- Adiciona `wails3 tool version` para incrementar versões semânticas, por [@leaanthony](https://github.com/leaanthony)
- Adiciona suporte a emblemas no macOS e no Windows por [@popaprozac](https://github.com/popaprozac) em [#](https://github.com/wailsapp/wails/pull/4234)
- Adiciona suporte a eventos registrados/estritamente tipados por [@fbbdev](https://github.com/fbbdev) e [@IanVS](https://github.com/IanVS) em [#4161](https://github.com/wailsapp/wails/pull/4161)
- Adiciona a capacidade de registrar hooks para eventos personalizados por [@fbbdev](https://github.com/fbbdev) e [@IanVS](https://github.com/IanVS) em [#4161](https://github.com/wailsapp/wails/pull/4161)
- `app.OpenFileManager(path string, selectFile bool)` para abrir o gerenciador de arquivos do sistema no caminho `path`, com destaque opcional por meio de `selectFile`, por [@Krzysztofz01](https://github.com/Krzysztofz01) [@rcalixte](https://github.com/rcalixte)
- Nova flag `-git` para o comando `wails3 init`, por [@leaanthony](https://github.com/leaanthony)
- Novo comando `wails3 generate webview2bootstrapper` por [@leaanthony](https://github.com/leaanthony)
- Adiciona o método `init()` ao runtime para permitir sua inicialização manual, por [@leaanthony](https://github.com/leaanthony)
- Adiciona a opção `WindowDidMoveDebounceMS` a WindowOptions da janela, por [@leaanthony](https://github.com/leaanthony)
- Adiciona o recurso de instância única por [@leaanthony](https://github.com/leaanthony). Com base no PR [v2](https://github.com/wailsapp/wails/pull/2951) de @APshenkin.
- Comando `wails3 generate template` por [@leaanthony](https://github.com/leaanthony)
- Comando `wails3 releasenotes` por [@leaanthony](https://github.com/leaanthony)
- Comando `wails3 update cli` por [@leaanthony](https://github.com/leaanthony)
- Opção `-clean` para o comando `wails3 generate bindings`, por [@leaanthony](https://github.com/leaanthony)
- Permite builds AppImage para Linux aarch64 (arm64), por [@AkshayKalose](https://github.com/AkshayKalose) em [#3981](https://github.com/wailsapp/wails/pull/3981)
- Adiciona um hiperlink para o patrocinador, por @ansxuman em [#3958](https://github.com/wailsapp/wails/pull/3958)
- Suporte ao empacotamento para Linux com builds de pacotes deb, rpm e Arch Linux, por
- Adiciona suporte a builds e pacotes universais para Darwin, por
- Documentação de eventos no site, por
- Templates para sveltekit e sveltekit-ts configurados para desenvolvimento sem SSR
- Atualiza os recursos de build usando o novo comando `wails3 update build-assets`, por
- Exemplo para testar a API HTML Drag and Drop, por
- Suporte a associações de arquivos por [leaanthony](https://github.com/leaanthony) em
- Novo comando `wails3 generate runtime` por
- Nova opção `InitialPosition` para especificar se a janela deve ser centralizada ou
- Adiciona os métodos `Path` e `Paths` ao pacote `application`, por
- Adiciona as opções do Windows `GeneralAutofillEnabled` e `PasswordAutosaveEnabled`
- Adiciona a capacidade de recuperar a janela que chama um método de serviço, por
- Adiciona as opções `EnabledFeatures` e `DisabledFeatures` para o WebView2, por
- ⊞ Novo sistema de DIP para suporte aprimorado a monitores com alta densidade de pixels, por
- ⊞ Opção de nome da classe da janela por [windom](https://github.com/windom/) em
- Os serviços foram ampliados para oferecer funcionalidade de plugins. Por
- 🐧 Eventos WindowDidMove/WindowDidResize em
- ⊞ Evento WindowDidResize em
-  Adiciona o evento ApplicationShouldHandleReopen para permitir tratar o Dock
-  Adiciona getPrimaryScreen/getScreens à implementação, por @tmclane em
-  Adiciona a opção de exibir a barra de ferramentas no modo de tela cheia no macOS, por
- 🐧 Adiciona a lógica onKeyPress para converter o pressionamento de tecla no Linux em um acelerador
- 🐧 Adiciona a tarefa `run:linux`, por
- Exporta o método `SetIcon`, por [@almas-x](https://github.com/almas-x) em
- Melhora `OnShutdown`, por [@almas-x](https://github.com/almas-x) em
- Restaura o método `ToggleMaximise` na interface `Window`, por
- Adiciona mais informações a `Environment()`. Por @leaanthony em
- Expõe o método `WebviewWindow.IsFocused` na interface `Window`, por
- Suporte a vários eventos de gatilho separados por espaços no sistema WML, por
- Adiciona exportações ESM ao script de runtime JS incluído no pacote, por
- Adiciona uma flag ao gerador de bindings para usar o script de runtime JS incluído no pacote em vez de
- Implementa `setIcon` no Linux, por [@abichinger](https://github.com/abichinger)
- Adiciona a flag `-port` ao comando dev e oferece suporte à variável de ambiente
- Adiciona testes para chamadas de métodos vinculados, por
- ⊞ Adiciona `SetIgnoreMouseEvents` para uma janela já criada, por
-  Adiciona a capacidade de definir o nível de empilhamento (ordem) de uma janela, por

### Corrigido

- Corrige `Screen.Bounds`, `WorkArea` e `Size` reduzidos à metade em Macs com tela Retina, convertendo os valores em pontos de NSScreen em pixels do dispositivo nos campos `Physical*` e preenchendo `Screen.X`/`Y` no nível superior, para que a detecção de contato entre vários monitores e o posicionamento na área de trabalho estejam corretos no [PR](https://github.com/wailsapp/wails/pull/5168) por @wayneforrest
- Corrige a condição de corrida de dados no ScreenManager que causa um deadlock no DisplayLink do WebKit quando a configuração de exibição muda (por exemplo, ao conectar ou desconectar um monitor externo durante a suspensão/retomada)
- Define CFBundleIconName diretamente como appicon quando Assets.car existe, no [PR](https://github.com/wailsapp/wails/pull/5154) por @symball
- Corrige `wails3 doctor`, que informava pacotes WebKitGTK incorretos no Fedora, openSUSE, Arch e NixOS — as entradas alternativas de 4.0 foram removidas, pois v3 exige a API 4.1 em tempo de compilação (#5071)
- Corrige o nome do pacote webkit2gtk do doctor para openSUSE (`webkit2gtk4_1-devel` → `webkit2gtk3-devel`, o nome correto do pacote no openSUSE) (#5071)
- Corrige o erro `Unexpected token '<'` quando `/wails/custom.js` está ausente no modo de desenvolvimento para desktop. Foi adicionado um manipulador 404 explícito para `/wails/custom.js` e uma validação de `Content-Type` que não diferencia maiúsculas de minúsculas em `loadOptionalScript`, a fim de impedir que fallbacks de SPA em HTML sejam injetados como JavaScript. ([#5068](https://github.com/wailsapp/wails/issues/5068))
- Corrige o estado de realce do menu da bandeja do sistema no macOS — agora o ícone exibe o estado selecionado quando o menu está aberto (#4910)
- Corrige a janela anexada à bandeja do sistema que aparecia atrás de outras janelas no macOS — agora ela usa o nível correto de janela pop-up (#4910)
- Corrige exemplos de importação de `@wailsio/runtime` incorretos em toda a documentação (#4989)
- Corrige a impossibilidade de minimizar janelas sem moldura no darwin (#4294)
- Corrige travamentos de 20-30 minuto durante `wails3 build` e `wails3 dev` ao excluir `node_modules/` da verificação de arquivos atualizados do go-task. Anteriormente, o glob `sources: "**/*"` fazia o go-task enumerar e calcular a soma de verificação de todos os arquivos em `node_modules/` (mais de 50000-100000 arquivos, com dependências pesadas como MUI), o que era especialmente lento no Windows/NTFS (#4939)
- Corrige a falha de compilação com GTK4 causada pela colisão do typedef C `Screen` com X11 Xlib.h (#4957)
- Corrige a consistência dos métodos de badge do Dock no macOS
- Corrige a aplicação de `InvisibleTitleBarHeight` a todas as janelas do macOS, em vez de apenas às janelas sem moldura ou com barra de título transparente (#4960)
- Corrige a trepidação da janela ao redimensioná-la pelos cantos superiores com `InvisibleTitleBarHeight` habilitado, ignorando o início do arraste próximo às bordas da janela (#4960)
- Corrige a geração de tipos mapeados com chaves de enum em bindings JS/TS (#4437) por @fbbdev
- Corrige o recurso de arrastar e soltar arquivos no Windows que não funcionava com escala de exibição diferente de 100%
- Corrige o recurso interno de arrastar e soltar do HTML5 que deixava de funcionar quando a soltura de arquivos estava habilitada no Windows
- Corrige as coordenadas de soltura de arquivos que usavam o espaço de pixels incorreto no Windows (pixels físicos em vez de pixels CSS)
- Corrige o funcionamento não confiável do recurso de arrastar e soltar arquivos no Linux com efeitos de foco do ponteiro
- Corrige o recurso interno de arrastar e soltar do HTML5 que deixava de funcionar quando a soltura de arquivos estava habilitada no Linux
- Corrige a exibição ou ocultação de janelas no Linux/GTK4 que, às vezes, as restaurava no estado minimizado, usando `gtk_window_present()` (#4957)
- Corrige a obtenção e definição da posição da janela no Linux/GTK4, que sempre retornava 0,0, adicionando suporte condicional a X11 por meio de `XTranslateCoordinates`/`XMoveWindow` (#4957)
- Corrige a falta de imposição do tamanho máximo da janela no Linux/GTK4, adicionando uma limitação de tamanho baseada em sinais para substituir o `gtk_window_set_geometry_hints` removido (#4957)
- Corrige o dimensionamento de DPI no Linux/GTK4, implementando o cálculo correto de PhysicalBounds e o suporte a dimensionamento fracionário por meio de `gdk_monitor_get_scale` (GTK 4.14+)
- Corrige a duplicação de itens de menu ao criar novas janelas no Linux/GTK4
- Corrige a geração de tipos mapeados com chaves de enum em bindings JS/TS (#4437) por @fbbdev
- Corrige o recurso de arrastar e soltar arquivos no Windows que não funcionava com escala de exibição diferente de 100%
- Corrige o recurso interno de arrastar e soltar do HTML5 que deixava de funcionar quando a soltura de arquivos estava habilitada no Windows
- Corrige as coordenadas de soltura de arquivos que usavam o espaço de pixels incorreto no Windows (pixels físicos em vez de pixels CSS)
- Corrige o funcionamento não confiável do recurso de arrastar e soltar arquivos no Linux com efeitos de foco do ponteiro
- Corrige o recurso interno de arrastar e soltar do HTML5 que deixava de funcionar quando a soltura de arquivos estava habilitada no Linux
- Corrige o dimensionamento de DPI no Linux/GTK4, implementando o cálculo correto de PhysicalBounds e o suporte a dimensionamento fracionário por meio de `gdk_monitor_get_scale` (GTK 4.14+)
- Corrige a duplicação de itens de menu ao criar novas janelas no Linux/GTK4
- Corrige a geração de tipos mapeados com chaves de enum em bindings JS/TS (#4437) por @fbbdev
- Corrige o problema de “janelas fantasmas” no macOS, causado pelo acesso às APIs do AppKit fora da thread principal em App.Window.Current() (#4947) por @wimaha
- Corrige o `<input type="file">` de HTML que não funcionava no macOS, implementando WKUIDelegate runOpenPanelWithParameters (#4862)
- Corrige o recurso nativo de arrastar e soltar arquivos que não funcionava ao usar o módulo npm `@wailsio/runtime` no macOS/Linux (#4953) por @leaanthony
- Corrige a geração de bindings para aliases de tipos entre pacotes (#4578) por @fbbdev
- Corrige uma falha do OpenFileDialog no Linux causada por uma violação de segurança de threads do GTK (#3683) por @ddmoney420
- Corrige uma falha SIGSEGV ao chamar `Focus()` em uma janela oculta ou destruída (#4890) por @ddmoney420
- Corrige um possível panic ao definir um ícone ou bitmap vazio no Linux (#4923) por @ddmoney420
- Corrige uma falha do ErrorDialog ao ser chamado a partir de um binding de serviço no macOS (#3631) por @leaanthony
- Faz com que os menus sejam exibidos no sistema operacional Windows em `v3\examples\dialogs` por @ndianabasi
- Corrige uma condição de corrida que causava TypeError durante o recarregamento da página (#4872) por @ddmoney420
- Corrige a saída incorreta dos testes do gerador de bindings, removendo o estado global do método `Collector.IsVoidAlias()` (#4941) por @fbbdev
- Corrige o seletor de arquivos `<input type="file">` que não funcionava no macOS (#4862) por @leaanthony
- Corrige o uso de sistemas de coordenadas inconsistentes por `Position()` e `SetPosition()` no macOS, que causava o deslocamento da posição da janela ao salvar e restaurar o estado (#4816) por @leaanthony
- Corrige o erro “Access is denied” de SetProcessDpiAwarenessContext quando a percepção de DPI já está definida pelo manifesto do aplicativo (#4803)
- Atualiza a página da documentação sobre atalhos de teclado e corrige o tipo do parâmetro de callback de `KeyBinding.Add` por @ndianabasi
- Corrige a documentação sobre a geração de bindings personalizados: é necessário usar `-d String` em vez de `-o String`
- Corrige a não remoção dos itens filhos do menu em `menu.Update()`
- Corrige referências desatualizadas à API Manager na documentação (os arquivos 31 foram atualizados para usar o novo padrão, como `app.Window.New()`, `app.Event.Emit()` etc.) por @leaanthony
- Corrige uma falha no Linux em caso de panic em métodos Go vinculados ao JS, causada pela substituição dos manipuladores de sinais pelo WebKit (#3965) por @leaanthony
- Corrige SaveFileDialog.SetFilename() que não produzia efeito no Linux (#4841) por @samstanier
- Corrige as coordenadas de soltura exibidas como undefined no exemplo de arrastar e soltar
- Corrige a falha na criação do pacote de aplicativo do macOS quando APP_NAME contém espaços (problema de expansão de chaves)
- Corrige um panic de índice fora dos limites no Windows ao chamar métodos de serviço (reverte goccy/go-json)
- Corrige o recurso de arrastar e soltar arquivos no Windows que não funcionava com escala de exibição diferente de 100%
- Corrige o recurso interno de arrastar e soltar do HTML5 que deixava de funcionar quando a soltura de arquivos estava habilitada no Windows
- Corrige as coordenadas de soltura de arquivos que usavam o espaço de pixels incorreto no Windows (pixels físicos em vez de pixels CSS)
- Corrige o funcionamento não confiável do recurso de arrastar e soltar arquivos no Linux com efeitos de foco do ponteiro
- Corrige o recurso interno de arrastar e soltar do HTML5 que deixava de funcionar quando a soltura de arquivos estava habilitada no Linux
- Atualiza todos os comandos nos arquivos Taskfile.yml de todos os sistemas operacionais para aceitar espaços em variáveis como `APP_NAME` por @ndianabasi
- Corrige o erro de argumento de comando ao executar a tarefa 'build:universal:lipo:go' no Linux por @wux1an
- Corrige o erro do Docker "símbolo indefinido: **<em>ubsan</em>handle_xxxxxxx" ao executar 'wails3 build GOOS=darwin GOARCH=arm64' no Linux, por @wux1an
- Consolida a documentação de protocolos personalizados e adiciona seções sobre Universal Links, por @leaanthony
- Corrige a falha do menu da bandeja do sistema no Windows ao clicar repetidamente no ícone, adicionando uma proteção contra chamadas simultâneas a TrackPopupMenuEx (#4151), por @leaanthony
- Impede que o aplicativo falhe ao chamar systray.Run() antes de app.Run(), por @leaanthony
- Corrige a falha no macOS ao alternar a visibilidade da janela por meio de Hide()/Show() com ApplicationShouldTerminateAfterLastWindowClosed habilitado (#4389), por @leaanthony
- Corrige o vazamento de memória nos menus de contexto no macOS e no Windows ao abri-los repetidamente (#4012), por @leaanthony
- Corrige a falta de reutilização dos recursos nativos dos menus de contexto no macOS, que fazia um novo menu ser criado a cada exibição (#4012), por @leaanthony
- Corrige o clique no ícone do Dock do macOS que não exibia janelas ocultas quando o aplicativo era iniciado com `Hidden: true` (#4583), por @leaanthony
- Corrige a caixa de diálogo de impressão que não abria no macOS devido ao tipo incorreto do ponteiro da janela na chamada CGO (#4290), por @leaanthony
- Corrige a falha do menu da janela no Wayland causada pelo acesso do appmenu-gtk-module a uma janela ainda não realizada (#4769), por @leaanthony
- Corrige a falha do aplicativo GTK quando o nome do aplicativo contém caracteres inválidos (espaços, parênteses etc.), por @leaanthony
- Corrige o erro "memória insuficiente" ao inicializar o recurso de arrastar e soltar no Windows (#4701), por @overlordtm
- Corrige a abertura do diretório errado pelo explorador de arquivos no Linux devido ao escape incorreto do URI (#4397), por @leaanthony
- Corrige a falha de compilação do AppImage em distribuições Linux modernas (Arch, Fedora 39+, Ubuntu 24.04+) detectando automaticamente seções ELF `.relr.dyn` e desabilitando a remoção de símbolos (#4642), por @leaanthony
- Corrige `wails doctor`, que informava incorretamente que os pacotes do WebKit estavam instalados em sistemas baseados no Fedora/DNF (#4457), por @leaanthony
- Corrige o comportamento padrão de `config.yml`, que executava `wails3 dev` com uma compilação de produção, por @mbaklor
- Corrige os stubs de serviços do iOS que causavam falhas de compilação devido à importação de um pacote inexistente, por @leaanthony
- Corrige o registro estruturado nos métodos debug/info que causava erros de "nenhuma diretiva de formatação", por @leaanthony
- Remove instruções temporárias de impressão de depuração incluídas acidentalmente na mesclagem da plataforma móvel, por @leaanthony
- Corrige a falha do WebKitGTK no Wayland com GPUs NVIDIA (erro de protocolo 71) desabilitando automaticamente o renderizador DMA-BUF, por @leaanthony
- Resolve o problema do valor alfa ser ignorado em `application.WebviewWindowOptions.BackgroundColour` no Linux ([#4722](https://github.com/wailsapp/wails/pull/4722), @BradHacker)
- Corrige o ícone da bandeja do sistema no Windows que não usava por padrão o ícone do aplicativo quando nenhum ícone personalizado era fornecido (#4704)
- Rastreia a propriedade de `HICON` para que somente os identificadores criados pelo usuário sejam destruídos, evitando falhas durante a reciclagem do Explorer (#4653).
- Libera o listener do tema do sistema do Windows e os ícones retidos da bandeja durante a destruição para impedir o vazamento de goroutines e contextos de dispositivo (#4653).
- Trunca as dicas de ferramenta da bandeja em 127 unidades UTF-16 para evitar a corrupção de pares substitutos e glifos multibyte (#4653).
- Corrige a falha da tarefa de empacotamento do Windows (#4667)
- Corrige a variável appicon do AppImage do Linux no Taskfile do Linux [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- Corrige o erro de compilação no Windows causado pela alteração da assinatura do go-webview2 v1.0.22 (#4513, #4645)
- Corrige a variável appicon do AppImage do Linux no Taskfile do Linux [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- Corrige o intervalo de protocolos em desktop.tmpl no Linux removendo `<.Info.Protocol>` até `<.Protocol>`, por @Tolfx em #4510
- Corrige o erro de redefinição na demonstração de liquid glass em [#4542](https://github.com/wailsapp/wails/pull/4542), por @Etesam913
- Corrige as atualizações do menu da bandeja do sistema no Linux [#4604](https://github.com/wailsapp/wails/issues/4604), por [@JackDoan](https://github.com/JackDoan)
- Corrige a janela branca que aparecia no Windows ao criar uma janela oculta, por @leaanthony em [#4612](https://github.com/wailsapp/wails/pull/4612)
- Corrige o caminho de importação do pacote de notificações na documentação, por @rxliuli em [#4617](https://github.com/wailsapp/wails/pull/4617)
- Corrige o recurso de arrastar e soltar que não funcionava ao usar o pacote npm @wailsio/runtime (#4489), por @leaanthony em #4616
- Windows: corrige a cintilação da janela na inicialização e a exibição incorreta de janelas ocultas em [PR](https://github.com/wailsapp/wails/pull/4600), por @leaanthony.
- Corrigidos por [@samstanier](https://github.com/samstanier) os problemas de maximização do tamanho da janela no Wayland (https://github.com/wailsapp/wails/issues/4429)
- Corrigidos os problemas de maximização do tamanho da janela no Wayland (https://github.com/wailsapp/wails/issues/4429) por [@samstanier](https://github.com/samstanier)
- Corrige o erro de redefinição na demonstração de liquid glass em [#4542](https://github.com/wailsapp/wails/pull/4542), por @Etesam913
- Corrige o problema que podia causar uma falha no AssetServer no macOS em [#4576](https://github.com/wailsapp/wails/pull/4576), por @jghiloni
- Corrige o problema de compilação ao usar NextJs. Corrigido em [#4585](https://github.com/wailsapp/wails/pull/4585), por @rev42
- Corrige os pipelines da versão nightly em [#4597](https://github.com/wailsapp/wails/pull/4597), por @riadafridishibly
- Corrige o erro de redefinição na demonstração de liquid glass em [#4542](https://github.com/wailsapp/wails/pull/4542), por @Etesam913
- Corrige o problema que podia causar uma falha no AssetServer no macOS em [#4576](https://github.com/wailsapp/wails/pull/4576), por @jghiloni
- Corrige o problema de compilação ao usar NextJs. Corrigido em [#4585](https://github.com/wailsapp/wails/pull/4585), por @rev42
- Corrige os pipelines da versão nightly em [#4597](https://github.com/wailsapp/wails/pull/4597), por @riadafridishibly
- Corrige o erro de redefinição na demonstração de liquid glass em [#4542](https://github.com/wailsapp/wails/pull/4542), por @Etesam913
- Corrige SetBackgroundColour no Windows, por @PPTGamer em [PR](https://github.com/wailsapp/wails/pull/4492)
- Atualiza a documentação para refletir as alterações da refatoração da API Manager, por @yulesxoxo em [PR #4476](https://github.com/wailsapp/wails/pull/4476)
- Corrige a variável appicon do arquivo .desktop do Linux no Taskfile do Linux [PR #4477](https://github.com/wailsapp/wails/pull/4477)
- Atualiza a documentação para refletir as alterações da refatoração da API Manager, por @yulesxoxo em [PR #4476](https://github.com/wailsapp/wails/pull/4476)
- Corrige o bug de desreferenciamento de ponteiro nil no Windows relatado em [#4456](https://github.com/wailsapp/wails/issues/4456), por @leaanthony em [#4460](https://github.com/wailsapp/wails/pull/4460)
- Adiciona suporte a `allowsBackForwardNavigationGestures` no WKWebView do macOS para habilitar gestos de navegação por deslizamento com dois dedos (#1857)
- Corrige o problema em que onClick não funcionava para itens de menu definidos inicialmente como desabilitados, por @leaanthony em [PR #4469](https://github.com/wailsapp/wails/pull/4469). Agradecimentos a @IanVS pela investigação inicial.
- Corrige a falta de encerramento do servidor Vite quando a compilação falha (#4403)
- Corrige o panic ao fechar ou cancelar um `SaveFileDialog` no Windows. Corrigido em [PR](https://github.com/wailsapp/wails/pull/4284), por @hkhere
- Corrigido o recurso de arrastar e soltar no nível do HTML no Windows por [@mbaklor](https://github.com/mbaklor) em [#4259](https://github.com/wailsapp/wails/pull/4259)
- Adicionado suporte a `allowsBackForwardNavigationGestures` no WKWebView do macOS para habilitar gestos de navegação por deslizamento com dois dedos (#1857)
- Corrigido o problema que impedia o funcionamento de onClick em itens de menu inicialmente definidos como desabilitados por @leaanthony em [PR #4469](https://github.com/wailsapp/wails/pull/4469). Agradecemos a @IanVS pela investigação inicial.
- Corrigida a ausência de limpeza do servidor Vite quando a compilação falha (#4403)
- Corrigida a análise de notificações no Windows por @popaprozac em [PR](https://github.com/wailsapp/wails/pull/4450)
- Corrigido o comando doctor para verificar as dependências do Windows SDK por [@kodumulo](https://github.com/kodumulo) em [#4390](https://github.com/wailsapp/wails/issues/4390)
- Corrigida a desreferência de ponteiro nil em processURLRequest no Mac por [@etesam913](https://github.com/etesam913) em [#4366](https://github.com/wailsapp/wails/pull/4366)
- Corrigido um bug no Linux que impedia diálogos com filtros por [@bh90210](https://github.com/bh90210) em [#4287](https://github.com/wailsapp/wails/pull/4287)
- Corrigidos problemas do menu Editar no Windows e no Linux por [@leaanthony](https://github.com/leaanthony) em [#3f78a3a](https://github.com/wailsapp/wails/commit/3f78a3a8ce7837e8b32242c8edbbed431c68c062)
- Atualizada a versão mínima do sistema nos arquivos .plist do macOS de 10.13.0 para 10.15.0 por [@AkshayKalose](https://github.com/AkshayKalose) em [#3981](https://github.com/wailsapp/wails/pull/3981)
- Corrigido o problema de IDs de janela ignorados por [@leaanthony](https://github.com/leaanthony)
- Corrigido o problema de menu nil ao chamar RegisterContextMenu por [@leaanthony](https://github.com/leaanthony)
- Corrigidos ciclos de dependência na saída do gerador de bindings por [@fbbdev](https://github.com/fbbdev) em [#4001](https://github.com/wailsapp/wails/pull/4001)
- Corrigidos erros de uso antes da definição na saída do gerador de bindings por [@fbbdev](https://github.com/fbbdev) em [#4001](https://github.com/wailsapp/wails/pull/4001)
- Sinalizadores de compilação agora são passados ao gerador de bindings por [@fbbdev](https://github.com/fbbdev) em [#4023](https://github.com/wailsapp/wails/pull/4023)
- Alterados os caminhos no Taskfile do Windows para usar barras normais e garantir seu funcionamento em plataformas que não sejam Windows por [@leaanthony](https://github.com/leaanthony)
- Corrigidos os eventos JS de Mac + Mac por [@leaanthony](https://github.com/leaanthony)
- Corrigido o deadlock de eventos no macOS por [@leaanthony](https://github.com/leaanthony)
- Corrigido um erro `Parameter incorrect` na inicialização de Window no Windows quando há HTML, mas não há JS, por [@leaanthony](https://github.com/leaanthony)
- Corrigido o tamanho do prefixo da resposta usado para detectar o tipo de conteúdo no servidor de recursos por [@fbbdev](https://github.com/fbbdev) em [#4049](https://github.com/wailsapp/wails/pull/4049)
- Corrigido o tratamento de respostas que não são 404 no caminho do índice raiz do servidor de recursos por [@fbbdev](https://github.com/fbbdev) em [#4049](https://github.com/wailsapp/wails/pull/4049)
- Corrigido o comportamento indefinido no gerador de bindings ao testar propriedades de tipos genéricos por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Corrigida a saída do gerador de bindings para modelos quando o tipo subjacente não tem as mesmas propriedades que o wrapper nomeado por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Corrigidas a saída do gerador de bindings para tipos de chave de mapas e a etapa de pré-processamento por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Corrigida a saída do gerador de bindings para structs que implementam interfaces de marshaling por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Corrigida a detecção de ciclos de tipos que envolvem tipos genéricos no gerador de bindings por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Corrigidas referências inválidas a modelos não exportados na saída do gerador de bindings por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Movido o código injetado para o final dos arquivos de serviço por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Corrigido o tratamento de erros de operações de fechamento de arquivos no gerador de bindings por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Suprimidos os avisos para serviços que definem métodos de ciclo de vida ou HTTP, mas nenhum outro método vinculado, por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Corrigida a falha dos modelos que não usam React ao exibir o rodapé Hello World quando o esquema de cores claro do sistema está ativo por [@marcus-crane](https://github.com/marcus-crane) em [#4056](https://github.com/wailsapp/wails/pull/4056)
- Corrigidos itens de menu ocultos no macOS por [@leaanthony](https://github.com/leaanthony)
- Corrigidos o tratamento e a formatação de erros nos processadores de mensagens por [@fbbdev](https://github.com/fbbdev) em [#4066](https://github.com/wailsapp/wails/pull/4066)
-  Corrigida a omissão do encerramento de serviços ao sair do aplicativo por [@fbbdev](https://github.com/fbbdev) em [#4066](https://github.com/wailsapp/wails/pull/4066)
-  Garantido que as atualizações de menu ocorram na thread principal por [@leaanthony](https://github.com/leaanthony)
- O mecanismo de arrastar e redimensionar agora é mais robusto e corresponde melhor ao comportamento esperado da plataforma, por [@fbbdev](https://github.com/fbbdev) em [#4100](https://github.com/wailsapp/wails/pull/4100)
- Corrigido [#4097](https://github.com/wailsapp/wails/issues/4097): Webpack/Angular descartava o código de inicialização do runtime, por [@fbbdev](https://github.com/fbbdev) em [#4100](https://github.com/wailsapp/wails/pull/4100)
- Corrigidos itens de menu inicialmente ocultos por [@IanVS](https://github.com/IanVS) em [#4116](https://github.com/wailsapp/wails/pull/4116)
- Corrigido o assetFileServer, que não fornecia arquivos `.html` em uma solicitação sem extensão quando `[request]` não existia, mas `[request].html` existia
- Corrigidos os caminhos de geração de ícones por [@robin-samuel](https://github.com/robin-samuel) em [#4125](https://github.com/wailsapp/wails/pull/4125)
- Corrigida a ausência de emissão dos eventos `fullscreen`, `unfullscreen`, `unminimise` e `unmaximise` por [@oSethoum](https://github.com/osethoum) em [#4130](https://github.com/wailsapp/wails/pull/4130)
- Corrigido o erro do NSIS causado por um prefixo incorreto na versão padrão da configuração por [@robin-samuel](https://github.com/robin-samuel) em [#4126](https://github.com/wailsapp/wails/pull/4126)
- Corrigida a função de runtime Dialogs, que retornava caminhos com caracteres de escape no Windows, por [TheGB0077](https://github.com/TheGB0077) em [#4188](https://github.com/wailsapp/wails/pull/4188)
- Corrigido o caminho de detecção do Webview2 em HKCU por [@leaanthony](https://github.com/leaanthony).
- Corrigido um problema de entrada no macOS por [@leaanthony](https://github.com/leaanthony).
- Corrigido o nome do arquivo de tarefas de geração de ícones do Windows por [@yulesxoxo](https://github.com/yulesxoxo) em [#4219](https://github.com/wailsapp/wails/pull/4219).
- Corrigido o problema de transparência em janelas sem moldura por [@leaanthony](https://github.com/leaanthony), com base no trabalho de @kron.
- Corrigidas as chamadas de foco quando a janela está desabilitada ou minimizada por [@leaanthony](https://github.com/leaanthony), com base no trabalho de @kron.
- Corrigidos os ícones da área de notificação que não reapareciam após a reinicialização da barra de tarefas por [@leaanthony](https://github.com/leaanthony), com base no trabalho de @kron.
- Corrigido o fato de fallbackResponseWriter não implementar Flush() em [#4245](https://github.com/wailsapp/wails/pull/4245)
- Corrigido o fato de fallbackResponseWriter não implementar Flush(), por [@superDingda], em [#4236](https://github.com/wailsapp/wails/issues/4236)
- Corrigidas as falhas ao fechar uma janela do macOS com uma chamada assíncrona pendente de função vinculada ao Go, por [@joshhardy](https://github.com/joshhardy), em [#4354](https://github.com/wailsapp/wails/pull/4354)
- Corrigida a condição de corrida na inicialização do modo de eficiência do Windows, por [@leaanthony](https://github.com/leaanthony)
- Corrigida a limpeza do identificador de ícone do Windows, por [@leaanthony](https://github.com/leaanthony).
- Corrigido `OpenFileManager` no Windows, por [@PPTGamer](https://github.com/PPTGamer), em [#4375](https://github.com/wailsapp/wails/pull/4375).
- Corrigidas as opções de largura mínima/máxima no Linux, por @atterpac, em [#3979](https://github.com/wailsapp/wails/pull/3979)
- Definições de tipos dos modelos TypeScript por meio de atualização da versão do npm, por @atterpac, em [#3966](https://github.com/wailsapp/wails/pull/3966)
- Corrigida a referência ao CSS do modelo SvelteKit, por @atterpac, em [#3945](https://github.com/wailsapp/wails/pull/3945)
- Garantido que os principais retornos de chamada em window run() sejam chamados na thread principal, por [@leaanthony](https://github.com/leaanthony)
- Corrigidos os exemplos do seletor de diretórios da caixa de diálogo, por [@leaanthony](https://github.com/leaanthony)
- Criada uma nova página de erro em chinês para quando index.html estiver ausente, por [@leaanthony](https://github.com/leaanthony)
-  Garantido que o retorno de chamada `windowDidBecomeKey` seja executado na thread principal, por [@leaanthony](https://github.com/leaanthony)
-  Adicionado suporte a tela cheia para janelas sem moldura, por [@leaanthony](https://github.com/leaanthony)
-  Aprimorada a lógica de destruição de janelas, por [@leaanthony](https://github.com/leaanthony)
-  Corrigida a lógica de posicionamento da janela quando anexada a ícones da bandeja do sistema, por [@leaanthony](https://github.com/leaanthony)
-  Adicionado suporte a tela cheia para janelas sem moldura, por [@leaanthony](https://github.com/leaanthony)
- Corrigido o tratamento de eventos, por [@leaanthony](https://github.com/leaanthony)
- Corrigida a lógica de encerramento de janelas, por [@leaanthony](https://github.com/leaanthony)
- O taskfile comum agora gera, por padrão, bindings TypeScript para modelos TypeScript, por [@leaanthony](https://github.com/leaanthony)
- Corrigido o encerramento do aplicativo ao receber a mensagem WM_CLOSE quando não há janelas abertas ou há somente a bandeja do sistema, por [@mmalcek](https://github.com/mmalcek), em [#3990](https://github.com/wailsapp/wails/pull/3990)
- Corrigida a compilação com garble, por @5aaee9, em [#3192](https://github.com/wailsapp/wails/pull/3192)
- Corrigidas as compilações NSIS do Windows, por [@leaanthony](https://github.com/leaanthony)
- Corrigido o deadlock na caixa de diálogo do Linux para múltiplas seleções, causado por um elemento não fechado
- Corrigida a limpeza multiplataforma de arquivos .syso durante a compilação para Windows, por
- Corrigida a compilação do AppImage para amd64, por @atterpac, em
- Corrigida a atualização dos ativos de compilação, por @ansxuman, em
- Corrigida a implementação de `OnClick` e `OnRightClick` da bandeja do sistema no Linux, por @atterpac
- Corrigido o fato de `AlwaysOnTop` não funcionar no Mac, por
-  Corrigido `application.NewEditMenu` que incluía uma duplicata
- 🐧 Corrigida a compilação para aarch64
- ⊞ Corrigidos os itens de menu de grupos de botões de opção, por
- Corrigido o erro ao compilar um .app executável no macOS quando 'name' e 'outputfilename'
- Corrigido um bug no uso de customEventProcessor no exemplo de arrastar e soltar, por
- 🐧 Corrigido o erro de compilação no Linux introduzido pela adição de IgnoreMouseEvents, por
- ⊞ Corrigido o bug na geração do arquivo de ícone syso, por
- 🐧 Incorporida a correção para execução nativa no Wayland, proveniente de
- Não vincular métodos de serviços internos em
- ⊞ Corrigido o panic na inicialização da bandeja do sistema em
- Não vincular métodos de serviços internos em
- ⊞ Corrigido o panic na inicialização da bandeja do sistema em
- Grande refatoração dos itens de menu e do tratamento de eventos. Por enquanto, melhora principalmente o macOS. Por
- Corrigidos os testes após a refatoração dos plugins e eventos em
- ⊞ Corrigido o aviso `Failed to unregister class Chrome_WidgetWin_0`. Por
- Problemas de módulos
- Corrigido o envio de mensagens de eventos de redimensionamento, por [atterpac](https://github.com/atterpac), em
- 🐧 Corrigido o erro no tratamento do tema no NixOS, por
- Corrigida a instalação do projeto entre volumes no Windows, por
- Corrigido o CSS do modelo React para exibir o rodapé, por
- Corrigidos os processos zumbis ao trabalhar no modo de desenvolvimento, com a atualização para a versão mais recente do refresh
- Corrigida a obtenção do arquivo WebKit do AppImage, por [Atterpac](https://github.com/atterpac)
- Corrigida a verificação de pacotes apt pelo Doctor, por [Atterpac](https://github.com/Atterpac), em
- Corrigido o congelamento do aplicativo ao sair (Darwin), por @5aaee9, em
- Corrigidas as cores de fundo dos exemplos no Windows, por
- Corrigidos os menus de contexto padrão, por [mmghv](https://github.com/mmghv), em
- Corrigidos os valores hexadecimais das teclas de seta no Darwin, por
- Implementado o funcionamento de arrastar e soltar no Windows. Adicionado por
- Corrigido um bug do Doctor no Linux quando o usuário não possui os drivers adequados
- Corrigido o dimensionamento por DPI na inicialização (Windows). Alterado por [@almas-x](https://github.com/almas-x), em
- Corrigida a linha de substituição em `go.mod` para usar caminhos relativos. Corrige caminhos do Windows com
- Corrigido o tratamento de cliques na bandeja do sistema do macOS quando não há uma janela anexada, por
- Corrigida a falha da compilação para Windows causada por uma opção desconhecida, por
- Corrigida a falha no Windows ao clicar com o botão esquerdo no ícone da bandeja do sistema sem ter um
- Corrigida a baseURL incorreta ao abrir a janela duas vezes, por @5aaee9, no PR
- Corrigida a ordem das ramificações if no método `WebviewWindow.Restore`, por
- Calculado corretamente `startURL` em várias invocações de `GetStartURL` quando
- Corrigido o tipo JS da struct `Screen` para corresponder à sua equivalente em Go, por
- Corrigido o método `WML.Reload` para garantir a limpeza adequada dos eventos registrados
- Corrigido o fechamento imediato do menu de contexto personalizado no Linux, por
- Corrigidos o caminho de saída e a extensão dos arquivos de modelo produzidos pelo binding
- Corrigidos os caminhos de importação dos arquivos de modelo no código JS produzido pelo binding
- Corrigido o recurso de arrastar e soltar em algumas distribuições Linux, por
- Corrigida a tarefa ausente para macOS ao usar `wails3 task dev`, por
- Corrigido o registro de eventos que causava uma atribuição em um mapa nil, por
- Corrigida a desserialização dos parâmetros de métodos vinculados, por
- Corrigir o tratamento de vários valores de retorno de métodos vinculados por
- Corrigir a detecção, pelo doctor, do npm que não foi instalado com o gerenciador de pacotes do sistema
- Corrigir a ausência de MicrosoftEdgeWebview2Setup.exe. Agradecimentos a
- Corrigir falha aleatória no Linux causada pelo tratamento do ID da janela, por @leaanthony. Com base em
- Corrigir falha de systemTray.setIcon no Linux por
- Garantir que o quadro da janela seja aplicado na primeira chamada da função `setFrameless` em

### Alterado

- **INCOMPATÍVEL**: as chaves de mapas nos bindings JS/TS gerados agora são marcadas como opcionais para refletir corretamente a semântica dos mapas de Go. O acesso aos valores de mapas no TypeScript agora retorna `T | undefined` em vez de `T`, exigindo verificações de valor nulo ou asserções (#4943), por `@fbbdev`
- Alterar o uso de `Event` para `Events` de acordo com as alterações em `@wailsio/runtime` e ajustar as chamadas de função correspondentes na documentação em `Features/Events/Event System`, por @AbdelhadiSeddar
- Mover `EnabledFeatures`, `DisabledFeatures` e `AdditionalBrowserArgs` das opções por janela para `Options.Windows` no nível da aplicação (#4559), por @leaanthony
- Atualizar o README do exemplo `Drag N Drop` e destacar que `Internal Drag and Drop` é demonstrado no exemplo, por @ndianabasi
- Alterar vários logs de depuração de Info para Debug (por @mbaklor)
- **INCOMPATÍVEL:** renomear `EnableDragAndDrop` como `EnableFileDrop` nas opções de janela
- **INCOMPATÍVEL:** renomear `DropZoneDetails` como `DropTargetDetails` no contexto do evento
- **INCOMPATÍVEL:** renomear o método `DropZoneDetails()` como `DropTargetDetails()` em `WindowEventContext`
- **INCOMPATÍVEL:** remover o evento `WindowDropZoneFilesDropped`; usar `WindowFilesDropped` em seu lugar
- **INCOMPATÍVEL:** alterar o atributo HTML de `data-wails-dropzone` para `data-file-drop-target`
- **INCOMPATÍVEL:** alterar a classe CSS de hover de `wails-dropzone-hover` para `file-drop-target-active`
- **INCOMPATÍVEL:** remover as opções `DragEffect`, `OnEnterEffect` e `OnOverEffect` do Windows (faziam parte do IDropTarget removido)
- Adotar goccy/go-json para todo o processamento de JSON em tempo de execução (bindings de métodos, eventos, solicitações da webview, notificações e kvstore), melhorando o desempenho em 21-63% e reduzindo as alocações de memória em 40-60%
- Otimizar o layout da struct BoundMethod e armazenar em cache o sinalizador isVariadic para reduzir a sobrecarga por chamada
- Usar um buffer de argumentos alocado na pilha para métodos com `<=8` argumentos, evitando alocações no heap
- Otimizar a coleta de resultados em chamadas de métodos para evitar a alocação de slices quando houver um único valor de retorno
- Usar sync.Map no cache de tipos MIME para melhorar o desempenho concorrente
- Usar um pool de buffers para ler o corpo de solicitações do transporte HTTP
- Alocar de forma tardia o canal CloseNotify no detector de tipo de conteúdo para reduzir as alocações por solicitação
- Remover do servidor de assets os logs de depuração de CSS
- Expandir o mapa de extensões de tipos MIME para abranger mais de 50 formatos comuns da Web (fontes, áudio, vídeo etc.)
- Atualizar a documentação das opções `X/Y` de Window, por @ruhuang2001
- Atualizar a documentação de `Frontend Runtime`, adicionando mais opções para gerar bindings de frontend, por @ndianabasi
- Atualizar a página de documentação do Asset Server v3 do Wails, por @ndianabasi
- **INCOMPATÍVEL**: remover as funções de diálogo no nível do pacote (`application.InfoDialog()`, `application.QuestionDialog()` etc.). Usar o gerenciador `app.Dialog` em seu lugar: `app.Dialog.Info()`, `app.Dialog.Question()`, `app.Dialog.Warning()`, `app.Dialog.Error()`, `app.Dialog.OpenFile()`, `app.Dialog.SaveFile()`
- Atualizar a documentação de diálogos para corresponder à API real: usar `app.Dialog.*`, `AddButton()` com callbacks (não `SetButtons()`), `SetDefaultButton(*Button)` (não string), `AddFilter()` (não `SetFilters()`), `SetFilename()` (não `SetDefaultFilename()`) e `app.Dialog.OpenFile().CanChooseDirectories(true)` para selecionar pastas
- **INCOMPATÍVEL**: builds de produção agora são o padrão. Para criar builds de desenvolvimento, defina `DEV=true` nos seus Taskfiles. Gerar um novo projeto para os exemplos, por @leaanthony
- Ao emitir um evento personalizado com zero ou um argumento de dados, o valor dos dados será atribuído diretamente ao campo Data, sem ser encapsulado em um slice, por [@fbbdev](https://github.com/fbbdev) em [#4633](https://github.com/wailsapp/wails/pull/4633)
- As bandejas do Windows agora respeitam `SystemTray.Show()`/`Hide()` ao alternar `NIS_HIDDEN`, permitindo que as aplicações realmente desapareçam e reapareçam (#4653).
- O registro da bandeja reutiliza ícones já resolvidos, define `NOTIFYICON_VERSION_4` uma única vez e habilita `NIF_SHOWTIP` para que as dicas de ferramenta sejam restauradas após a reinicialização do Explorer (#4653).
- macOS: usar `visibleFrame` em vez de `frame` para centralizar a janela excluindo as áreas da barra de menus e do Dock
- macOS: usar `visibleFrame` em vez de `frame` para centralizar a janela excluindo as áreas da barra de menus e do Dock
- Ao executar `wails3 update build-assets` com o parâmetro `-config`, os valores definidos por meio dos parâmetros `-product*` são
- `window.NativeWindowHandle()` -> `window.NativeWindow()`, por @leaanthony em [#4471](https://github.com/wailsapp/wails/pull/4471)
- Refatorar o tratamento interno de janelas, por @leaanthony em [#4471](https://github.com/wailsapp/wails/pull/4471)
- Remover `application.WindowIDKey` e `application.WindowNameKey` (substituídos por `application.WindowKey`), por [@leaanthony](https://github.com/leaanthony)
- ContextMenuData agora retorna uma string em vez de any, por [@leaanthony](https://github.com/leaanthony)
- Nos bindings JS/TS, os campos de classe com tipos de array de tamanho fixo agora são inicializados com o tamanho esperado, em vez de ficarem vazios, por [@fbbdev](https://github.com/fbbdev) em [#4001](https://github.com/wailsapp/wails/pull/4001)
- ContextMenuData agora retorna uma string em vez de any, por [@leaanthony](https://github.com/leaanthony)
- `application.NewService` não aceita mais opções como parâmetro opcional (use `application.NewServiceWithOptions` em seu lugar), por [@leaanthony](https://github.com/leaanthony) em [#4024](https://github.com/wailsapp/wails/pull/4024)
- Remover a dependência de `nanoid`, por [@leaanthony](https://github.com/leaanthony)
- Atualizar o exemplo de Window para os estilos de janela mica/acrylic/tabbed, por [@leaanthony](https://github.com/leaanthony)
- Nos bindings JS/TS, os arquivos de modelo `internal.js/ts` foram removidos; agora todos os modelos estão em `models.js/ts`, por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Nos bindings JS/TS, tipos nomeados nunca são renderizados como aliases de outros tipos nomeados; o comportamento anterior agora se restringe a aliases, por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Nos bindings JS/TS, no modo de classe, os campos de struct cujo tipo é um parâmetro de tipo são marcados como opcionais e nunca são inicializados automaticamente, por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Remover o ESLint dos templates, por [@IanVS](https://github.com/IanVS) em [#4059](https://github.com/wailsapp/wails/pull/4059)
- Atualizar o ano dos direitos autorais para 2025, por [@IanVS](https://github.com/IanVS) em [#4037](https://github.com/wailsapp/wails/pull/4037)
- Adicionada documentação para event.Sender por [@IanVS](https://github.com/IanVS) em [#4075](https://github.com/wailsapp/wails/pull/4075)
- Suporte ao Go 1.24 por [@leaanthony](https://github.com/leaanthony)
- Os hooks de `ServiceStartup` agora são invocados quando `App.Run` é chamado, e não em `application.New`, por [@fbbdev](https://github.com/fbbdev) em [#4066](https://github.com/wailsapp/wails/pull/4066)
- Os erros de `ServiceStartup` agora são retornados por `App.Run` em vez de encerrar o processo, por [@fbbdev](https://github.com/fbbdev) em [#4066](https://github.com/wailsapp/wails/pull/4066)
- As chamadas de bindings e caixas de diálogo feitas pelo JS agora são rejeitadas com objetos de erro em vez de strings, por [@fbbdev](https://github.com/fbbdev) em [#4066](https://github.com/wailsapp/wails/pull/4066)
- Melhorado o posicionamento do menu da bandeja do sistema no Windows por [@leaanthony](https://github.com/leaanthony)
- O runtime JS foi portado para TypeScript por [@fbbdev](https://github.com/fbbdev) em [#4100](https://github.com/wailsapp/wails/pull/4100)
- O runtime é inicializado assim que é importado, sem a necessidade de aguardar o carregamento da janela, por [@fbbdev](https://github.com/fbbdev) em [#4100](https://github.com/wailsapp/wails/pull/4100)
- O runtime não exporta mais um método init. Uma importação somente para efeitos colaterais pode ser usada para inicializá-lo, por [@fbbdev](https://github.com/fbbdev) em [#4100](https://github.com/wailsapp/wails/pull/4100)
- Os métodos vinculados agora retornam uma `CancellablePromise` que é rejeitada com um `CancelError` em caso de cancelamento. O resultado efetivo da chamada é descartado, por [@fbbdev](https://github.com/fbbdev) em [#4100](https://github.com/wailsapp/wails/pull/4100)
- Os tipos de serviços integrados agora são chamados consistentemente de `Service`, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- As funções de criação de serviços integrados com opções agora são chamadas consistentemente de `NewWithConfig`, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- O método `Select` do serviço `sqlite` agora se chama `Query` para manter a consistência com as APIs Go, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Templates: runtime movido para "dependencies" e arquivos package.json organizados por [@IanVS](https://github.com/IanVS) em [#4133](https://github.com/wailsapp/wails/pull/4133)
- Cria e assina ad hoc os pacotes de aplicativos durante o desenvolvimento para habilitar determinadas APIs do macOS, por [@popaprozac](https://github.com/popaprozac) em [#4171](https://github.com/wailsapp/wails/pull/4171)
- Recursos de build movidos para diretórios específicos de cada plataforma por [@leaanthony](https://github.com/leaanthony)
- Taskfiles movidos e renomeados para diretórios específicos de cada plataforma por [@leaanthony](https://github.com/leaanthony)
- Criada uma experiência muito melhor quando `index.html` está ausente, por [@leaanthony](https://github.com/leaanthony)
- [Windows] Melhorado o desempenho ao minimizar e restaurar por [@leaanthony](https://github.com/leaanthony). Com base no [PR](https://github.com/wailsapp/wails/pull/3955) original de [562589540](https://github.com/562589540)
- Removida a opção `ShouldClose` (em vez disso, registre um hook para events.Common.WindowClosing) por [@leaanthony](https://github.com/leaanthony)
- [Windows] Reduzida a cintilação ao abrir uma janela por [@leaanthony](https://github.com/leaanthony)
- Removido `Window.Destroy`, pois deveria ser uma função interna, por [@leaanthony](https://github.com/leaanthony)
- Eventos `WindowClose` renomeados para `WindowClosing` por [@leaanthony](https://github.com/leaanthony)
- Os builds do frontend agora usam o ambiente do vite "development" ou "production", conforme o tipo de build, por [@leaanthony](https://github.com/leaanthony)
- Atualização para go-webview2 v1.19 por [@leaanthony](https://github.com/leaanthony)
- Garantido o uso do fork do taskfile por @leaanthony
- Atualizado o fork do Taskfile para corrigir problemas de versão ao instalar usando
- Uso do fork do Taskfile para corrigir problemas de versão ao instalar usando
- `service.OnStartup` agora encerra o aplicativo em caso de erro e executa
- Refatoradas as mensagens de clique na bandeja do sistema para se alinharem melhor às interações do usuário por
- Incorporação de recursos passa a incluir `all:frontend/dist` para oferecer suporte a frameworks que geram
- Refatoração do Taskfile por [leaanthony](https://github.com/leaanthony) em
- Atualização para `go-webview2` v1.0.16 por
- Corrigido o tipo `Screen` para incluir `ID`, e não `Id`, por
- Atualizada a versão do Wails em `go.mod.tmpl` para oferecer suporte a `application.ServiceOptions` por
- Corrigida a determinação do nome do serviço por [windom](https://github.com/windom/) em
- mkdocs serve agora usa Docker por [leaanthony](https://github.com/leaanthony)
- Configuração de desenvolvimento consolidada em `config.yml` por
- A caixa de diálogo da bandeja do sistema agora usa por padrão o ícone do aplicativo, quando disponível (Windows), por
- Melhor relatório de GPU + memória no macOS por
- Removidos `WebviewGpuIsDisabled` e `EnableFraudulentWebsiteWarnings`
- Alteração na API de eventos: `On`/`Emit` -> eventos do usuário, `OnApplicationEvent` ->
- Correção da API de eventos no Linux por [TheGB0077](https://github.com/TheGB0077) em
- [CI] melhorias nas actions e habilitação da execução de actions também em forks e
- `AbsolutePosition()` renomeado para `Position()` por
- Atualizada a dependência webkit do Linux para webkit2gtk-4.1 em vez de webkitgtk2-4.0 para
- O script do runtime JS incluído no pacote agora é um módulo ESM: tags script que o importam
- O pacote `@wailsio/runtime` não publica sua API no `window.wails`
- O módulo `@wailsio/runtime/src/window` da API Window agora expõe o contêiner
- A API de janela do JS foi atualizada para corresponder ao `WebviewWindow` atual do Go
- O gerador de bindings agora usa chamadas por ID por padrão. A opção `-id` da CLI
- Novo layout do código de bindings: anteriormente, os arquivos de saída eram organizados em pastas
- O campo `application.Options.Bind` da struct foi renomeado para
- Nova sintaxe para serviços de binding: agora as instâncias de serviço devem ser encapsuladas em um
- Desabilitado o spinner em ambiente sem terminal ou de CI por

### Removido

- **ALTERAÇÃO INCOMPATÍVEL**: removidos `EnabledFeatures`, `DisabledFeatures` e `AdditionalLaunchArgs` das opções `WindowsWindow` por janela. Em vez disso, use `Options.Windows.EnabledFeatures`, `Options.Windows.DisabledFeatures` e `Options.Windows.AdditionalBrowserArgs` no nível do aplicativo. Esses sinalizadores se aplicam globalmente ao ambiente WebView2 compartilhado (#4559), por @leaanthony
- Removida a implementação nativa de `IDropTarget` no Windows em favor da abordagem baseada em JavaScript (corresponde ao comportamento de v2)
- Removida a dependência github.com/wailsapp/mimetype em favor de um mapa de extensões ampliado + http.DetectContentType da biblioteca padrão, reduzindo o tamanho do binário em cerca de 1.2 MB
- Remove a dependência de gopkg.in/ini.v1 implementando um analisador mínimo de arquivos .desktop para o explorador de arquivos do Linux, economizando cerca de 45 KB
- Remove samber/lo do código de runtime usando o pacote slices da biblioteca padrão do Go 1.21+ e auxiliares internos mínimos, economizando cerca de 310 KB
- Remove as instruções printf de depuração do manipulador de esquemas de URL do Darwin (#4834)
- **ALTERAÇÃO INCOMPATÍVEL**: remove o evento `linux:WindowLoadChanged` — use `linux:WindowLoadFinished` para detectar quando o WebView terminar de carregar (#3896), por @leaanthony

### Alterações incompatíveis

- **Refatoração da API de gerenciadores**: reorganiza a API da aplicação, substituindo a estrutura plana por gerenciadores organizados para melhorar a organização e a facilidade de descoberta do código, por [@leaanthony](https://github.com/leaanthony) em [#4359](https://github.com/wailsapp/wails/pull/4359)
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
- Move os métodos `Path` e `Paths` para o pacote `application`, por [@leaanthony](https://github.com/leaanthony)
- O menu da aplicação agora está disponível apenas no macOS, por [@leaanthony](https://github.com/leaanthony)

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

não são mais ignorados e substituem o valor da configuração.

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

@atterpac in [#3909](https://github.com/wailsapp/wails/3909)   [ansxuman](https://github.com/ansxuman) in   [#3902](https://github.com/wailsapp/wails/pull/3902)   [atterpac](https://github.com/atterpac) in   [#3867](https://github.com/wailsapp/wails/pull/3867)   by [atterpac](https://github.com/atterpac) in   [#3829](https://github.com/wailsapp/wails/pull/3829)   [leaanthony](https://github.com/leaanthony)   [FerroO2000](https://github.com/FerroO2000) in   [#3856](https://github.com/wailsapp/wails/pull/3856)   [#3873](https://github.com/wailsapp/wails/pull/3873)   [leaanthony](https://github.com/leaanthony)   positioned at the given X/Y location by   [leaanthony](https://github.com/leaanthony) in   [#3885](https://github.com/wailsapp/wails/pull/3885)   [ansxuman](https://github.com/ansxuman) and   [leaanthony](https://github.com/leaanthony) in   [#3823](https://github.com/wailsapp/wails/pull/3823)   by [leaanthony](https://github.com/leaanthony) in   [#3766](https://github.com/wailsapp/wails/pull/3766)   [leaanthony](https://github.com/leaanthony) in   [#3888](https://github.com/wailsapp/wails/pull/3888)   [leaanthony](https://github.com/leaanthony). -

### Alterado

`service.OnShutdown`para quaisquer serviços anteriores iniciados por @atterpac em   [#3920](https://github.com/wailsapp/wails/pull/3920)   @atterpac em [#3907](https://github.com/wailsapp/wails/pull/3907)   subpastas por @atterpac em   [#3887](https://github.com/wailsapp/wails/pull/3887)   [#3748](https://github.com/wailsapp/wails/pull/3748)   [leaanthony](https://github.com/leaanthony)   [etesam913](https://github.com/etesam913) em   [#3778](https://github.com/wailsapp/wails/pull/3778)   [northes](https://github.com/northes) em   [#3836](https://github.com/wailsapp/wails/pull/3836)   [#3827](https://github.com/wailsapp/wails/pull/3827)   [leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   (substituído pelas opções `EnabledFeatures` e `DisabledFeatures`) por   [leaanthony](https://github.com/leaanthony)

### Corrigido

variável de canal por @michael-freling em   [#3925](https://github.com/wailsapp/wails/pull/3925)   [ansxuman](https://github.com/ansxuman) em   [#3924](https://github.com/wailsapp/wails/pull/3924)   [#3898](https://github.com/wailsapp/wails/pull/3898)   [#3901](https://github.com/wailsapp/wails/pull/3901)   em [#3886](https://github.com/wailsapp/wails/pull/3886)   [leaanthony](https://github.com/leaanthony) em   [#3841](https://github.com/wailsapp/wails/pull/3841)   função `PasteAndMatchStyle` no menu de edição no Darwin por   [johnmccabe](https://github.com/johnmccabe) em   [#3839](https://github.com/wailsapp/wails/pull/3839)   [#3840](https://github.com/wailsapp/wails/issues/3840) em   [#3854](https://github.com/wailsapp/wails/pull/3854) por   [kodflow](https://github.com/kodflow)   [@leaanthony](https://github.com/leaanthony)   são diferentes. por @nickisworking em   [#3789](https://github.com/wailsapp/wails/pull/3789)

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

[#3580](https://github.com/wailsapp/wails/pull/3580)   [#3580](https://github.com/wailsapp/wails/pull/3580)   clique no ícone por @5aaee9 em [#2991](https://github.com/wailsapp/wails/pull/2991)   [#2618](https://github.com/wailsapp/wails/pull/2618)   [@fbbdev](https://github.com/fbbdev) em   [#3282](https://github.com/wailsapp/wails/pull/3282)   @[Atterpac](https://github.com/Atterpac)   em [#3022](https://github.com/wailsapp/wails/pull/3022])   [@marcus-crane](https://github.com/marcus-crane) em   [#3146](https://github.com/wailsapp/wails/pull/3146)   [PR](https://github.com/wailsapp/wails/pull/3147)   [PR](https://github.com/wailsapp/wails/pull/3189)   [@fbbdev](https://github.com/fbbdev) em   [#3281](https://github.com/wailsapp/wails/pull/3281)   [aba82cc](https://github.com/wailsapp/wails/commit/aba82cc52787c97fb99afa58b8b63a0004b7ff6c)   com base no [PR](https://github.com/wailsapp/wails/pull/2044) por @Mai-Lapyst   [@fbbdev](https://github.com/fbbdev) em   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@fbbdev](https://github.com/fbbdev) em   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@fbbdev](https://github.com/fbbdev) em   [#3295](https://github.com/wailsapp/wails/pull/3295)   o pacote npm por [@fbbdev](https://github.com/fbbdev) em   [#3334](https://github.com/wailsapp/wails/pull/3334)   em [#3354](https://github.com/wailsapp/wails/pull/3354)   `WAILS_VITE_PORT` por [@abichinger](https://github.com/abichinger) em   [#3429](https://github.com/wailsapp/wails/pull/3429)   [@abichinger](https://github.com/abichinger) em   [#3431](https://github.com/wailsapp/wails/pull/3431)   [@bruxaodev](https://github.com/bruxaodev) em   [#3667](https://github.com/wailsapp/wails/pull/3667)   [@OlegGulevskyy](https://github.com/OlegGulevskyy) em   [#3674](https://github.com/wailsapp/wails/pull/3674)

### Corrigido

[#3606](https://github.com/wailsapp/wails/pull/3606)   [tmclane](https://github.com/tmclane) em   [#3515](https://github.com/wailsapp/wails/pull/3515)   [atterpac](https://github.com/atterac) em   [#3512](https://github.com/wailsapp/wails/pull/3512)   [atterpac](https://github.com/atterpac) em   [#3477](https://github.com/wailsapp/wails/pull/3477)   por [Atterpac](https://github.com/atterpac) em   [#3320](https://github.com/wailsapp/wails/pull/3320).   em [#3306](https://github.com/wailsapp/wails/pull/3306).   [#2972](https://github.com/wailsapp/wails/pull/2972).   [#2982](https://github.com/wailsapp/wails/pull/2982)   [mmghv](https://github.com/mmghv) em   [#2750](https://github.com/wailsapp/wails/pull/2750).   [#2753](https://github.com/wailsapp/wails/pull/2753).   [jaybeecave](https://github.com/jaybeecave) em   [#3052](https://github.com/wailsapp/wails/pull/3052).   [@pylotlight](https://github.com/pylotlight) no   [PR](https://github.com/wailsapp/wails/pull/3039)   instalado. Adicionado por [@pylotlight](https://github.com/pylotlight) no   [PR](https://github.com/wailsapp/wails/pull/3032)   [PR](https://github.com/wailsapp/wails/pull/3145)   espaços — @leaanthony.   [thomas-senechal](https://github.com/thomas-senechal) no PR   [#3207](https://github.com/wailsapp/wails/pull/3207)   [thomas-senechal](https://github.com/thomas-senechal) no PR   [#3208](https://github.com/wailsapp/wails/pull/3208)   janela anexada [tw1nk](https://github.com/tw1nk) no PR   [#3271](https://github.com/wailsapp/wails/pull/3271)   [#3273](https://github.com/wailsapp/wails/pull/3273)   [@fbbdev](https://github.com/fbbdev) em   [#3279](https://github.com/wailsapp/wails/pull/3279)   `FRONTEND_DEVSERVER_URL` está presente.   [#3299](https://github.com/wailsapp/wails/pull/3299)   [@fbbdev](https://github.com/fbbdev) em   [#3295](https://github.com/wailsapp/wails/pull/3295)   listeners por [@fbbdev](https://github.com/fbbdev) em   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@abichinger](https://github.com/abichinger) em   [#3330](https://github.com/wailsapp/wails/pull/3330)   gerador por [@fbbdev](https://github.com/fbbdev) em   [#3334](https://github.com/wailsapp/wails/pull/3334)   gerador por [@fbbdev](https://github.com/fbbdev) em   [#3334](https://github.com/wailsapp/wails/pull/3334)   [@abichinger](https://github.com/abichinger) em   [#3346](https://github.com/wailsapp/wails/pull/3346)   [@hfoxy](https://github.com/hfoxy) em   [#3417](https://github.com/wailsapp/wails/pull/3417)   [@hfoxy](https://github.com/hfoxy) em   [#3426](https://github.com/wailsapp/wails/pull/3426)   [@fbbdev](https://github.com/fbbdev) em   [#3431](https://github.com/wailsapp/wails/pull/3431)   [@fbbdev](https://github.com/fbbdev) em   [#3431](https://github.com/wailsapp/wails/pull/3431)   por [@pekim](https://github.com/pekim) em   [#3458](https://github.com/wailsapp/wails/pull/3458)   [@robin-samuel](https://github.com/robin-samuel).   PR [#3466](https://github.com/wailsapp/wails/pull/3622) por   [@5aaee9](https://github.com/5aaee9).   [@windom](https://github.com/windom/) em   [#3636](https://github.com/wailsapp/wails/pull/3636).   Windows por [@bruxaodev](https://github.com/bruxaodev/) em   [#3691](https://github.com/wailsapp/wails/pull/3691).

### Alterado

[mmghv](https://github.com/mmghv) em   [#3611](https://github.com/wailsapp/wails/pull/3611)   compatibilidade com Ubuntu 24.04 LTS por [atterpac](https://github.com/atterpac) em   [#3461](https://github.com/wailsapp/wails/pull/3461)   deve ter o atributo `type="module"`. Por   [@fbbdev](https://github.com/fbbdev) em   [#3295](https://github.com/wailsapp/wails/pull/3295)   objeto e não inicia o sistema WML. Isso foi feito para melhorar o encapsulamento. Se desejado, o sistema WML pode ser iniciado manualmente chamando o novo método `WML.Enable`. O script de runtime JS incluído ainda executa ambas as operações automaticamente. Por [@fbbdev](https://github.com/fbbdev) em   [#3295](https://github.com/wailsapp/wails/pull/3295)   objeto de janela como exportação padrão. Não é mais possível importar métodos individuais por meio da sintaxe ESM de importação nomeada ou de namespace.   API. Alguns métodos tiveram o nome ou o protótipo alterado, especificamente: `Screen` passa a ser `GetScreen`; `GetZoomLevel`/`SetZoomLevel` passam a ser `GetZoom`/`SetZoom`; `GetZoom`, `Width` e `Height` agora retornam valores diretamente, em vez de encapsulá-los em objetos. Por [@fbbdev](https://github.com/fbbdev) em   [#3295](https://github.com/wailsapp/wails/pull/3295)   foi removido. Use a opção `-names` da CLI para voltar às chamadas por nome. Por [@fbbdev](https://github.com/fbbdev) em   [#3468](https://github.com/wailsapp/wails/pull/3468)   nomeados de acordo com o pacote que os contém; agora são usados os caminhos completos de importação do Go, incluindo o caminho do módulo. Por [@fbbdev](https://github.com/fbbdev) em   [#3468](https://github.com/wailsapp/wails/pull/3468)   `application.Options.Services`. Por [@fbbdev](https://github.com/fbbdev) em   [#3468](https://github.com/wailsapp/wails/pull/3468)   chamada a `application.NewService`. Por [@fbbdev](https://github.com/fbbdev) em   [#3468](https://github.com/wailsapp/wails/pull/3468)   [@DeltaLaboratory](https://github.com/DeltaLaboratory) em   [#3574](https://github.com/wailsapp/wails/pull/3574)
