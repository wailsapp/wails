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

/_   * ATUALIZADO AUTOMATICAMENTE *   As notas de lançamento são geradas a partir de PRs mescladas. Não edite este arquivo em PRs de contribuição comuns nem recrie o `v3/UNRELEASED_CHANGELOG.md` removido. Consulte `AGENTS.md` para saber como manter o registro de alterações. _/

## [Não lançado]

## v3.2.0-nightly.20261103 - 2026-11-03

## Alterado
- [v3] atualizar os menus da bandeja do sistema em Menu.Update na [PR](https://github.com/wailsapp/wails/pull/6314) de @johnmaguire
- [v3/Windows] Implementar SetIcon para janelas na [PR](https://github.com/wailsapp/wails/pull/6502) de @hypherionmc

## v3.2.0-nightly.20261102 - 2026-11-02

## Alterado
- DisableOverscroll em vez de DisableBounce na [PR](https://github.com/wailsapp/wails/pull/6077) de @lethusang

## v3.2.0-nightly.20261101 - 2026-11-01

## Adicionado
- adicionar back-ends experimentais sem CGO para macOS e Linux na [PR](https://github.com/wailsapp/wails/pull/6521) de @taliesin-ai
- compilações YAML nativas, empacotamento portátil e assinatura na [PR](https://github.com/wailsapp/wails/pull/6531) de @taliesin-ai

## Corrigido
- Corrigir diversos problemas das compilações nightly de v3 e limpar as dependências na [PR](https://github.com/wailsapp/wails/pull/6519) de @taliesin-ai
- corrigir a validação de lançamento de v3 e a classificação de fixtures na [PR](https://github.com/wailsapp/wails/pull/6533) de @taliesin-ai
- resolver falhas de validação de lançamento entre plataformas na [PR](https://github.com/wailsapp/wails/pull/6538) de @taliesin-ai

## v3.1.0 - 2026-11-01

## Adicionado
- adicionar back-ends experimentais sem CGO para macOS e Linux na [PR](https://github.com/wailsapp/wails/pull/6521) de @taliesin-ai
- compilações YAML nativas, empacotamento portátil e assinatura na [PR](https://github.com/wailsapp/wails/pull/6531) de @taliesin-ai

## Alterado
- evitar copiar metadados de desktop descartados na [PR](https://github.com/wailsapp/wails/pull/6391) de @taliesin-ai
- formatar tamanhos de conteúdo sem fmt na [PR](https://github.com/wailsapp/wails/pull/6393) de @taliesin-ai
- selecionar raízes de índice sem coletar candidatos na [PR](https://github.com/wailsapp/wails/pull/6404) de @taliesin-ai
- reutilizar buffers de cópia de arquivos sem suporte a seek na [PR](https://github.com/wailsapp/wails/pull/6408) de @taliesin-ai
- evitar alocar um leitor para detectar prefixos na [PR](https://github.com/wailsapp/wails/pull/6410) de @taliesin-ai
- ignorar leituras de detecção de MIME para extensões conhecidas na [PR](https://github.com/wailsapp/wails/pull/6414) de @taliesin-ai
- impedir o crescimento da fila e a oscilação de redução na [PR](https://github.com/wailsapp/wails/pull/6415) de @taliesin-ai
- evitar alocar contextos de requisição descartados na [PR](https://github.com/wailsapp/wails/pull/6419) de @taliesin-ai
- ignorar a preparação de logs de depuração desativados na [PR](https://github.com/wailsapp/wails/pull/6420) de @taliesin-ai
- alocar buffers de detecção de MIME de forma tardia na [PR](https://github.com/wailsapp/wails/pull/6428) de @taliesin-ai
- reduzir alocações na formatação de hiperlinks na [PR](https://github.com/wailsapp/wails/pull/6443) de @taliesin-ai
- remover alocações temporárias da formatação de marcadores na [PR](https://github.com/wailsapp/wails/pull/6449) de @taliesin-ai
- reduzir a conversão de UUID em string a uma única alocação na [PR](https://github.com/wailsapp/wails/pull/6453) de @taliesin-ai
- reduzir alocações na formatação de identificadores de instrução na [PR](https://github.com/wailsapp/wails/pull/6454) de @taliesin-ai
- eliminar a alocação temporária do resumo de UUID baseado em nome na [PR](https://github.com/wailsapp/wails/pull/6455) de @taliesin-ai
- evitar copiar linhas incompletas da saída de comandos na [PR](https://github.com/wailsapp/wails/pull/6460) de @taliesin-ai
- transmitir as cópias de executáveis durante a preparação do pacote na [PR](https://github.com/wailsapp/wails/pull/6472) de @taliesin-ai
- evitar verificar novamente linhas de saída incompletas na [PR](https://github.com/wailsapp/wails/pull/6473) de @taliesin-ai
- reutilizar padrões compilados de roteamento de comandos na [PR](https://github.com/wailsapp/wails/pull/6475) de @taliesin-ai
- evitar consultas de metadados durante a descoberta do runtime na [PR](https://github.com/wailsapp/wails/pull/6482) de @taliesin-ai
- reduzir alocações da página de fallback localizada na [PR](https://github.com/wailsapp/wails/pull/6485) de @taliesin-ai
- evitar recompilar padrões para cada processador na [PR](https://github.com/wailsapp/wails/pull/6490) de @taliesin-ai

## Corrigido
- implementar SetTooltip para a bandeja do sistema na [PR](https://github.com/wailsapp/wails/pull/6017) de @fan711
- ler o modo escuro no portal de configurações na [PR](https://github.com/wailsapp/wails/pull/6072) de @TechHutTV
- permitir que Window.SetZoom fique abaixo de 1.0 na [PR](https://github.com/wailsapp/wails/pull/6166) de @APshenkin
- entregar a cada listener Go na ordem de emissão na [PR](https://github.com/wailsapp/wails/pull/6275) de @DevLumuz
- ler os hooks de eventos sob hookLock em Emit na [PR](https://github.com/wailsapp/wails/pull/6276) de @DevLumuz
- sempre anunciar a propriedade Menu da bandeja para que uma bandeja criada após Run não seja descartada na [PR](https://github.com/wailsapp/wails/pull/6284) de @AlbinoGeek
- reter o trabalho da thread principal enviado antes de Run até que Run seja iniciado na [PR](https://github.com/wailsapp/wails/pull/6287) de @AlbinoGeek
- usar a URL correta do ativo no provedor de atualização do GitHub na [PR](https://github.com/wailsapp/wails/pull/6304) de @gavidroselj
- manter o corpo da requisição quando o WebView2 omitir Content-Length na [PR](https://github.com/wailsapp/wails/pull/6326) de @strraaa
- relatar as telas do macOS e GTK no espaço DIP (limites com DPI misto) na [PR](https://github.com/wailsapp/wails/pull/6346) de @julianstorer
- não relatar uma notificação dispensada como um clique na [PR](https://github.com/wailsapp/wails/pull/6367) de @bercedev
- encerrar os processos de desenvolvimento quando o terminal fechar na [PR](https://github.com/wailsapp/wails/pull/6371) de @xmasdev
- associar execuções do registro de alterações aos commits de squash merge na verificação de evidências na [PR](https://github.com/wailsapp/wails/pull/6373) de @taliesin-ai
- mantenha por um segundo o debounce da reinicialização da barra de tarefas no [PR](https://github.com/wailsapp/wails/pull/6386) por @taliesin-ai
- mantenha as restrições de tamanho durante uma tela cheia iniciada nesse estado no [PR](https://github.com/wailsapp/wails/pull/6387) por @taliesin-ai
- não redimensione uma janela em tela cheia ou maximizada para atender a uma nova restrição no [PR](https://github.com/wailsapp/wails/pull/6388) por @taliesin-ai
- não exiba uma caixa de diálogo antes de Run no exemplo liquid-glass no [PR](https://github.com/wailsapp/wails/pull/6389) por @taliesin-ai
- faça com que uma caixa de diálogo solicitada antes de Run não faça nada, em vez de esperar para sempre, no [PR](https://github.com/wailsapp/wails/pull/6390) por @taliesin-ai
- use o fallback quando a entrada de desktop não tiver um comando no [PR](https://github.com/wailsapp/wails/pull/6392) por @taliesin-ai
- detecte o conteúdo atual sem reter nomes de arquivos no [PR](https://github.com/wailsapp/wails/pull/6394) por @taliesin-ai
- reconheça extensões de assets com combinação de maiúsculas e minúsculas no [PR](https://github.com/wailsapp/wails/pull/6395) por @taliesin-ai
- libere os buffers consumidos do detector de conteúdo no [PR](https://github.com/wailsapp/wails/pull/6396) por @taliesin-ai
- libere as referências removidas da ringqueue no [PR](https://github.com/wailsapp/wails/pull/6397) por @taliesin-ai
- respeite as exclusões de cabeçalhos antes da seleção do fallback no [PR](https://github.com/wailsapp/wails/pull/6398) por @taliesin-ai
- confirme os cabeçalhos da resposta de fallback antes do flush no [PR](https://github.com/wailsapp/wails/pull/6399) por @taliesin-ai
- descubra apenas arquivos de índice com correspondência exata no [PR](https://github.com/wailsapp/wails/pull/6400) por @taliesin-ai
- evite corpos e leituras completas em solicitações HEAD não posicionáveis no [PR](https://github.com/wailsapp/wails/pull/6401) por @taliesin-ai
- preencha o prefixo de detecção MIME ao longo de leituras curtas no [PR](https://github.com/wailsapp/wails/pull/6402) por @taliesin-ai
- confirme respostas vazias não posicionáveis bem-sucedidas no [PR](https://github.com/wailsapp/wails/pull/6403) por @taliesin-ai
- evite duplicar destinos de URI expandidos no [PR](https://github.com/wailsapp/wails/pull/6405) por @taliesin-ai
- preserve os caminhos durante a expansão de campos no [PR](https://github.com/wailsapp/wails/pull/6406) por @taliesin-ai
- respeite os diretórios de pesquisa de arquivos desktop do XDG no [PR](https://github.com/wailsapp/wails/pull/6407) por @taliesin-ai
- confirme os cabeçalhos do detector em um flush vazio no [PR](https://github.com/wailsapp/wails/pull/6409) por @taliesin-ai
- ignore diretórios que ocultam entradas de desktop no [PR](https://github.com/wailsapp/wails/pull/6411) por @taliesin-ai
- preserve os destinos encontrados dentro de caminhos de executáveis no [PR](https://github.com/wailsapp/wails/pull/6412) por @taliesin-ai
- libere as capturas do callback após o disparo do temporizador no [PR](https://github.com/wailsapp/wails/pull/6413) por @taliesin-ai
- evite correspondências falsas de bibliotecas do ldconfig no [PR](https://github.com/wailsapp/wails/pull/6416) por @taliesin-ai
- restrinja as correspondências de fallback de bibliotecas no sistema de arquivos no [PR](https://github.com/wailsapp/wails/pull/6417) por @taliesin-ai
- elimine a condição de corrida de dados na invalidação do cache no [PR](https://github.com/wailsapp/wails/pull/6418) por @taliesin-ai
- respeite a configuração opcional do logger no [PR](https://github.com/wailsapp/wails/pull/6421) por @taliesin-ai
- evite a mutação compartilhada da URL inicial e portas obsoletas no [PR](https://github.com/wailsapp/wails/pull/6422) por @taliesin-ai
- interrompa as novas tentativas do proxy de desenvolvimento quando a solicitação for cancelada no [PR](https://github.com/wailsapp/wails/pull/6423) por @taliesin-ai
- retorne 404 para assets empacotados ausentes no [PR](https://github.com/wailsapp/wails/pull/6424) por @taliesin-ai
- processe solicitações HEAD do runtime empacotado no [PR](https://github.com/wailsapp/wails/pull/6425) por @taliesin-ai
- preserve os caminhos com escape ao rotear serviços no [PR](https://github.com/wailsapp/wails/pull/6426) por @taliesin-ai
- ordene corretamente identificadores numéricos de pré-lançamento grandes demais no [PR](https://github.com/wailsapp/wails/pull/6427) por @taliesin-ai
- rejeite erros HTTP antes de salvar downloads de ferramentas no [PR](https://github.com/wailsapp/wails/pull/6429) por @taliesin-ai
- respeite o prazo limite de inicialização durante a descoberta do desktop no [PR](https://github.com/wailsapp/wails/pull/6430) por @taliesin-ai
- preserve os corpos das solicitações entre novas tentativas do proxy de desenvolvimento no [PR](https://github.com/wailsapp/wails/pull/6432) por @taliesin-ai
- feche os corpos das respostas HTTP em todos os caminhos de retorno no [PR](https://github.com/wailsapp/wails/pull/6434) por @taliesin-ai
- preserve os erros de cancelamento em pesquisas paralelas de bibliotecas no [PR](https://github.com/wailsapp/wails/pull/6435) por @taliesin-ai
- rejeite respostas malsucedidas de tags do GitHub antes da decodificação no [PR](https://github.com/wailsapp/wails/pull/6436) por @taliesin-ai
- preserve literalmente caminhos de modelos entre colchetes no [PR](https://github.com/wailsapp/wails/pull/6437) por @taliesin-ai
- instale o wails3 uma única vez e preserve os erros do instalador no [PR](https://github.com/wailsapp/wails/pull/6438) por @taliesin-ai
- ignore entradas ocultas e cabeçalhos em blocos de comentários no [PR](https://github.com/wailsapp/wails/pull/6439) por @taliesin-ai
- exponha erros de consulta da versão estável no [PR](https://github.com/wailsapp/wails/pull/6440) por @taliesin-ai
- imprima mensagens simples sem interpretação de formato no [PR](https://github.com/wailsapp/wails/pull/6441) por @taliesin-ai
- evite panics de gravação após carregar null de JSON no [PR](https://github.com/wailsapp/wails/pull/6442) por @taliesin-ai
- retorne snapshots para leituras do armazenamento inteiro no [PR](https://github.com/wailsapp/wails/pull/6444) por @taliesin-ai
- propague erros após a iteração parcial de linhas no [PR](https://github.com/wailsapp/wails/pull/6446) por @taliesin-ai
- mantenha literais os caminhos de destino da extração no [PR](https://github.com/wailsapp/wails/pull/6447) por @taliesin-ai
- libere os handles de instruções preparadas ao fechar e reabrir no [PR](https://github.com/wailsapp/wails/pull/6448) por @taliesin-ai
- processe a preparação que concorre com o fechamento ou a reabertura no [PR](https://github.com/wailsapp/wails/pull/6450) por @taliesin-ai
- respeite solicitações HEAD para a página de índice de fallback no [PR](https://github.com/wailsapp/wails/pull/6451) por @taliesin-ai
- exija handles hexadecimais completos de instruções no [PR](https://github.com/wailsapp/wails/pull/6452) por @taliesin-ai
- preserve o banco de dados padrão em memória sob concorrência no [PR](https://github.com/wailsapp/wails/pull/6456) por @taliesin-ai
- separe stderr da saída HEAD analisada no [PR](https://github.com/wailsapp/wails/pull/6457) por @taliesin-ai
- gere assets de runtime a partir de pacotes incorporados no [PR](https://github.com/wailsapp/wails/pull/6458) por @taliesin-ai
- detecte placeholders TODO, FIXME e TBD no [PR](https://github.com/wailsapp/wails/pull/6459) por @taliesin-ai
- valide as respostas da API de notas de versão no [PR](https://github.com/wailsapp/wails/pull/6461) por @taliesin-ai
- respeite os diretórios das tarefas nas verificações de status no [PR](https://github.com/wailsapp/wails/pull/6462) por @taliesin-ai
- evite a perda de dados ao copiar um arquivo para ele mesmo no [PR](https://github.com/wailsapp/wails/pull/6463) por @taliesin-ai
- inclua o diretório da tarefa na identidade do cache no [PR](https://github.com/wailsapp/wails/pull/6464) por @taliesin-ai
- sonde as portas uma única vez por verificação no [PR](https://github.com/wailsapp/wails/pull/6465) por @taliesin-ai
- evite um panic quando as entradas do cache forem decodificadas como null no [PR](https://github.com/wailsapp/wails/pull/6466) por @taliesin-ai
- incremente pré-lançamentos numéricos sem limites de inteiros no [PR](https://github.com/wailsapp/wails/pull/6467) por @taliesin-ai
- impedir que verificações legadas ocultem falhas de cache no [PR](https://github.com/wailsapp/wails/pull/6468) por @taliesin-ai
- executar tarefas padrão no diretório base do executor no [PR](https://github.com/wailsapp/wails/pull/6469) por @taliesin-ai
- preservar sinais de porcentagem literais em comandos do explorador no [PR](https://github.com/wailsapp/wails/pull/6470) por @taliesin-ai
- alinhar os diretórios dos comandos de pré-condição e de tarefa no [PR](https://github.com/wailsapp/wails/pull/6471) por @taliesin-ai
- substituir links simbólicos por arquivos regulares gerados no [PR](https://github.com/wailsapp/wails/pull/6474) por @taliesin-ai
- preservar DMGs existentes quando as compilações falharem no [PR](https://github.com/wailsapp/wails/pull/6476) por @taliesin-ai
- preservar regex entre aspas e barras invertidas de caminhos no [PR](https://github.com/wailsapp/wails/pull/6477) por @taliesin-ai
- validar os tamanhos dos ícones gerados antes da alocação da imagem no [PR](https://github.com/wailsapp/wails/pull/6478) por @taliesin-ai
- encaminhar cadeias de comandos e redirecionamentos pelo shell no [PR](https://github.com/wailsapp/wails/pull/6479) por @taliesin-ai
- preservar pacotes Linux anteriores em caso de falha na compilação no [PR](https://github.com/wailsapp/wails/pull/6480) por @taliesin-ai
- usar execução alternativa do shell específica para cada plataforma no [PR](https://github.com/wailsapp/wails/pull/6481) por @taliesin-ai
- servir recursos com barras iniciais repetidas no caminho no [PR](https://github.com/wailsapp/wails/pull/6483) por @taliesin-ai
- reconhecer versões estáveis com metadados de compilação no [PR](https://github.com/wailsapp/wails/pull/6484) por @taliesin-ai
- relatar com precisão falhas na geração de recursos no [PR](https://github.com/wailsapp/wails/pull/6486) por @taliesin-ai
- evitar condições de corrida ao substituir o relator ativo no [PR](https://github.com/wailsapp/wails/pull/6487) por @taliesin-ai
- serializar salvamentos simultâneos do arquivo de cache no [PR](https://github.com/wailsapp/wails/pull/6488) por @taliesin-ai
- aceitar vários candidatos de igualdade em templates no [PR](https://github.com/wailsapp/wails/pull/6489) por @taliesin-ai
- normalizar mapas de categorias salvos com valor nulo no [PR](https://github.com/wailsapp/wails/pull/6491) por @taliesin-ai
- mostrar diagnósticos de comandos de shell com variáveis que falharam no [PR](https://github.com/wailsapp/wails/pull/6492) por @taliesin-ai
- manter cadeias de comandos do shell fora do cache implícito do Go no [PR](https://github.com/wailsapp/wails/pull/6493) por @taliesin-ai
- serializar a persistência de categorias no Linux no [PR](https://github.com/wailsapp/wails/pull/6494) por @taliesin-ai
- armazenar em cache saídas do Go com espaços entre aspas ou com escape no [PR](https://github.com/wailsapp/wails/pull/6495) por @taliesin-ai
- preservar ícones do Mac quando a codificação falhar no [PR](https://github.com/wailsapp/wails/pull/6496) por @taliesin-ai
- parar de descartar flags inválidas de contagem de testes do Go no [PR](https://github.com/wailsapp/wails/pull/6497) por @taliesin-ai
- manter as permissões de origem em novas cópias de arquivos no [PR](https://github.com/wailsapp/wails/pull/6498) por @taliesin-ai
- Corrigir diversos problemas da compilação nightly v3 e limpar dependências no [PR](https://github.com/wailsapp/wails/pull/6519) por @taliesin-ai
- corrigir a validação de versões v3 e a classificação de fixtures no [PR](https://github.com/wailsapp/wails/pull/6533) por @taliesin-ai
- resolver falhas de validação de versões entre plataformas no [PR](https://github.com/wailsapp/wails/pull/6538) por @taliesin-ai

## v3.1.0-nightly.20261028 - 2026-10-28

## Alterado
- Melhorar a velocidade de geração do changelog reutilizando padrões regex compilados no [PR](https://github.com/wailsapp/wails/pull/6490) por @taliesin-ai
- Reduzir alocações ao servir páginas alternativas localizadas no [PR](https://github.com/wailsapp/wails/pull/6485) por @taliesin-ai
- Acelerar a descoberta do runtime do AppImage evitando consultas desnecessárias de metadados no [PR](https://github.com/wailsapp/wails/pull/6482) por @taliesin-ai
- Reutilizar padrões compilados de roteamento de comandos para melhorar o desempenho do wake no [PR](https://github.com/wailsapp/wails/pull/6475) por @taliesin-ai
- Melhorar o desempenho evitando verificar novamente linhas de saída incompletas no [PR](https://github.com/wailsapp/wails/pull/6473) por @taliesin-ai
- Transmitir cópias de executáveis durante a preparação de pacotes MSIX para reduzir o uso de memória e acelerar as compilações no [PR](https://github.com/wailsapp/wails/pull/6472) por @taliesin-ai
- Evitar copiar linhas incompletas da saída de comandos para melhorar o desempenho no [PR](https://github.com/wailsapp/wails/pull/6460) por @taliesin-ai
- Reduzir alocações ao gerar UUIDs baseados em nomes no [PR](https://github.com/wailsapp/wails/pull/6455) por @taliesin-ai
- Reduzir alocações ao formatar identificadores de instruções SQLite no [PR](https://github.com/wailsapp/wails/pull/6454) por @taliesin-ai
- Reduzir a conversão de UUIDs em strings a uma única alocação no [PR](https://github.com/wailsapp/wails/pull/6453) por @taliesin-ai
- Reduzir alocações ao formatar itens com marcadores do changelog no [PR](https://github.com/wailsapp/wails/pull/6449) por @taliesin-ai
- Reduzir alocações ao formatar hiperlinks de terminal no [PR](https://github.com/wailsapp/wails/pull/6443) por @taliesin-ai
- Reduzir alocações ao servir recursos detectando buffers MIME somente quando necessário no [PR](https://github.com/wailsapp/wails/pull/6428) por @taliesin-ai
- Ignorar a preparação de logs de depuração no servidor de recursos quando esse tipo de log estiver desativado no [PR](https://github.com/wailsapp/wails/pull/6420) por @taliesin-ai
- Reduzir alocações no processamento de solicitações da webview evitando contextos de solicitação descartados no [PR](https://github.com/wailsapp/wails/pull/6419) por @taliesin-ai
- Acelerar o fornecimento de recursos ignorando a detecção de MIME para extensões de arquivo conhecidas no [PR](https://github.com/wailsapp/wails/pull/6414) por @taliesin-ai
- Reduzir alocações ao servir recursos evitando usar um leitor para detectar o tipo do conteúdo no [PR](https://github.com/wailsapp/wails/pull/6410) por @taliesin-ai
- Reduzir alocações ao servir arquivos de recursos que não permitem posicionamento no [PR](https://github.com/wailsapp/wails/pull/6408) por @taliesin-ai
- Acelerar a seleção da raiz de índice do servidor de recursos evitando coletar candidatos no [PR](https://github.com/wailsapp/wails/pull/6404) por @taliesin-ai
- Melhorar o desempenho do servidor de recursos formatando comprimentos de conteúdo sem fmt no [PR](https://github.com/wailsapp/wails/pull/6393) por @taliesin-ai
- Reduzir a cópia de memória ao ler metadados do explorador de arquivos do desktop no [PR](https://github.com/wailsapp/wails/pull/6391) por @taliesin-ai

## Corrigido
- Manter as permissões do arquivo de origem ao copiá-lo para novos arquivos de destino no [PR](https://github.com/wailsapp/wails/pull/6498) por @taliesin-ai
- Impedir que janelas em tela cheia e maximizadas sejam redimensionadas quando as restrições de tamanho mínimo ou máximo mudarem no [PR](https://github.com/wailsapp/wails/pull/6388) por @taliesin-ai
- Publicar atualizações agendadas do histórico de estrelas do repositório wailsapp/wails no [PR](https://github.com/wailsapp/wails/pull/6500) por @taliesin-ai
- Parar de descartar flags inválidas de contagem de testes do Go no wake no [PR](https://github.com/wailsapp/wails/pull/6497) por @taliesin-ai
- Preservar ícones do Mac quando a codificação falhar no [PR](https://github.com/wailsapp/wails/pull/6496) por @taliesin-ai
- Corrigir falhas do cache de compilação para saídas do Go com espaços entre aspas ou com escape no [PR](https://github.com/wailsapp/wails/pull/6495) por @taliesin-ai
- Serializar a persistência de categorias de notificações no Linux para evitar condições de corrida em gravações simultâneas no [PR](https://github.com/wailsapp/wails/pull/6494) por @taliesin-ai
- Manter cadeias de comandos do shell fora do cache implícito de compilação do Go no [PR](https://github.com/wailsapp/wails/pull/6493) por @taliesin-ai
- Corrigir a falha ao carregar mapas de categorias de notificação quando os dados salvos são nulos no [PR](https://github.com/wailsapp/wails/pull/6491) de @taliesin-ai
- Exibir diagnósticos de comandos de shell com variáveis que falharam no [PR](https://github.com/wailsapp/wails/pull/6492) de @taliesin-ai
- Corrigir a correspondência de modelos do wake quando existem vários candidatos de igualdade no [PR](https://github.com/wailsapp/wails/pull/6489) de @taliesin-ai
- Serializar salvamentos simultâneos de arquivos de cache para evitar gravações corrompidas no cache no [PR](https://github.com/wailsapp/wails/pull/6488) de @taliesin-ai
- Evitar condições de corrida ao substituir o relator ativo no [PR](https://github.com/wailsapp/wails/pull/6487) de @taliesin-ai
- Relatar com precisão as falhas de geração de recursos na CLI no [PR](https://github.com/wailsapp/wails/pull/6486) de @taliesin-ai
- Reconhecer na CLI versões estáveis que incluem metadados de build no [PR](https://github.com/wailsapp/wails/pull/6484) de @taliesin-ai
- Servir ativos cujos caminhos de solicitação contêm barras iniciais repetidas no [PR](https://github.com/wailsapp/wails/pull/6483) de @taliesin-ai
- Corrigir a execução alternativa via shell do wake em cada plataforma no [PR](https://github.com/wailsapp/wails/pull/6481) de @taliesin-ai
- Preservar pacotes Linux anteriores quando um build falhar no [PR](https://github.com/wailsapp/wails/pull/6480) de @taliesin-ai
- Encaminhar cadeias de comandos e redirecionamentos pelo shell no [PR](https://github.com/wailsapp/wails/pull/6479) de @taliesin-ai
- Preservar expressões regulares entre aspas e barras invertidas de caminhos ao analisar comandos do wake no [PR](https://github.com/wailsapp/wails/pull/6477) de @taliesin-ai
- Preservar arquivos DMG existentes quando o empacotamento de builds falhar no [PR](https://github.com/wailsapp/wails/pull/6476) de @taliesin-ai
- Substituir links simbólicos por arquivos regulares gerados na saída de bindings no [PR](https://github.com/wailsapp/wails/pull/6474) de @taliesin-ai
- Preservar sinais de porcentagem literais em comandos do explorador de arquivos no Linux no [PR](https://github.com/wailsapp/wails/pull/6470) de @taliesin-ai
- Alinhar os diretórios dos comandos de pré-condição e de tarefa para que sejam resolvidos de forma consistente no [PR](https://github.com/wailsapp/wails/pull/6471) de @taliesin-ai
- Executar tarefas padrão no diretório-base do executor no [PR](https://github.com/wailsapp/wails/pull/6469) de @taliesin-ai
- Impedir que verificações legadas ocultem falhas de cache no [PR](https://github.com/wailsapp/wails/pull/6468) de @taliesin-ai
- Corrigir a numeração de pré-lançamentos da CLI para que seja incrementada sem limites de estouro de inteiros no [PR](https://github.com/wailsapp/wails/pull/6467) de @taliesin-ai
- Evitar um panic quando entradas do cache forem decodificadas como nulas no [PR](https://github.com/wailsapp/wails/pull/6466) de @taliesin-ai
- Sondar portas uma vez por verificação na CLI no [PR](https://github.com/wailsapp/wails/pull/6465) de @taliesin-ai
- Incluir o diretório da tarefa na identidade do cache para evitar resultados obsoletos do wake no [PR](https://github.com/wailsapp/wails/pull/6464) de @taliesin-ai
- Evitar perda de dados quando o empacotamento copiar um arquivo sobre ele mesmo no [PR](https://github.com/wailsapp/wails/pull/6463) de @taliesin-ai
- Fazer as verificações de status respeitarem os diretórios das tarefas no [PR](https://github.com/wailsapp/wails/pull/6462) de @taliesin-ai
- Detectar placeholders TODO, FIXME e TBD em entradas do changelog no [PR](https://github.com/wailsapp/wails/pull/6459) de @taliesin-ai
- Gerar ativos de runtime a partir de bundles incorporados na CLI no [PR](https://github.com/wailsapp/wails/pull/6458) de @taliesin-ai
- Separar stderr da saída de HEAD analisada nos comandos git no [PR](https://github.com/wailsapp/wails/pull/6457) de @taliesin-ai
- Preservar o banco de dados SQLite padrão em memória durante acessos simultâneos no [PR](https://github.com/wailsapp/wails/pull/6456) de @taliesin-ai
- Retornar respostas corretas a solicitações HEAD da página de índice alternativa no [PR](https://github.com/wailsapp/wails/pull/6451) de @taliesin-ai
- Tratar a condição de corrida entre a preparação de instruções SQLite e o fechamento ou a reabertura do banco de dados no [PR](https://github.com/wailsapp/wails/pull/6450) de @taliesin-ai
- Liberar identificadores de instruções preparadas quando o banco de dados SQLite for fechado ou reaberto no [PR](https://github.com/wailsapp/wails/pull/6448) de @taliesin-ai
- Manter literais os caminhos de destino da extração ao usar gosod no [PR](https://github.com/wailsapp/wails/pull/6447) de @taliesin-ai
- Propagar erros de consultas SQLite após uma iteração parcial das linhas no [PR](https://github.com/wailsapp/wails/pull/6446) de @taliesin-ai
- Retornar snapshots em leituras de todo o armazenamento no kvstore no [PR](https://github.com/wailsapp/wails/pull/6444) de @taliesin-ai
- Evitar panics de gravação no kvstore após o carregamento de um valor JSON nulo no [PR](https://github.com/wailsapp/wails/pull/6442) de @taliesin-ai
- Imprimir mensagens simples no terminal sem interpretar verbos de formatação no [PR](https://github.com/wailsapp/wails/pull/6441) de @taliesin-ai
- Expor erros ao consultar a versão estável mais recente no atualizador no [PR](https://github.com/wailsapp/wails/pull/6440) de @taliesin-ai
- Ignorar entradas ocultas e cabeçalhos em blocos de comentários do changelog no [PR](https://github.com/wailsapp/wails/pull/6439) de @taliesin-ai
- Instalar wails3 uma única vez e expor erros do instalador no atualizador no [PR](https://github.com/wailsapp/wails/pull/6438) de @taliesin-ai
- Preservar caminhos literais de modelos entre colchetes no gosod no [PR](https://github.com/wailsapp/wails/pull/6437) de @taliesin-ai
- Rejeitar respostas malsucedidas de tags do GitHub antes de decodificá-las no atualizador no [PR](https://github.com/wailsapp/wails/pull/6436) de @taliesin-ai
- Preservar erros de cancelamento durante pesquisas paralelas de bibliotecas no [PR](https://github.com/wailsapp/wails/pull/6435) de @taliesin-ai
- Fechar os corpos das respostas HTTP em todos os caminhos de retorno para evitar vazamentos de recursos no [PR](https://github.com/wailsapp/wails/pull/6434) de @taliesin-ai
- Preservar os corpos das solicitações quando o proxy de desenvolvimento repetir as solicitações no [PR](https://github.com/wailsapp/wails/pull/6432) de @taliesin-ai
- Respeitar o prazo de inicialização durante a descoberta do explorador de arquivos da área de trabalho no [PR](https://github.com/wailsapp/wails/pull/6430) de @taliesin-ai
- Rejeitar erros HTTP antes de salvar downloads das ferramentas do AppImage no [PR](https://github.com/wailsapp/wails/pull/6429) de @taliesin-ai
- Ordenar corretamente identificadores numéricos de pré-lançamento muito grandes em comparações semver no [PR](https://github.com/wailsapp/wails/pull/6427) de @taliesin-ai
- Preservar caminhos com escape ao encaminhar serviços no servidor de ativos no [PR](https://github.com/wailsapp/wails/pull/6426) de @taliesin-ai
- Tratar solicitações HEAD do runtime incluído no bundle no servidor de ativos no [PR](https://github.com/wailsapp/wails/pull/6425) de @taliesin-ai
- Retornar 404 para ativos ausentes do bundle, em vez de uma resposta incorreta, no [PR](https://github.com/wailsapp/wails/pull/6424) de @taliesin-ai
- Interromper novas tentativas do proxy de desenvolvimento quando uma solicitação for cancelada no [PR](https://github.com/wailsapp/wails/pull/6423) de @taliesin-ai
- Evitar a mutação compartilhada da URL inicial e portas obsoletas no servidor de ativos no [PR](https://github.com/wailsapp/wails/pull/6422) de @taliesin-ai
- Respeitar a configuração opcional do logger no servidor de ativos no [PR](https://github.com/wailsapp/wails/pull/6421) de @taliesin-ai
- Eliminar a condição de corrida na invalidação do cache de caminhos de bibliotecas no [PR](https://github.com/wailsapp/wails/pull/6418) de @taliesin-ai
- Evitar correspondências falsas de bibliotecas ao pesquisar a saída do ldconfig no [PR](https://github.com/wailsapp/wails/pull/6416) de @taliesin-ai
- Evitar o crescimento e a oscilação de redução da fila do servidor de ativos sob carga no [PR](https://github.com/wailsapp/wails/pull/6415) de @taliesin-ai
- Liberar as capturas do callback após o disparo do temporizador de debounce para evitar a retenção de referências obsoletas no [PR](https://github.com/wailsapp/wails/pull/6413) de @taliesin-ai
- Preservar alvos do explorador de arquivos encontrados dentro de caminhos de executáveis no [PR](https://github.com/wailsapp/wails/pull/6412) de @taliesin-ai
- Ignorar diretórios que ocultam entradas de desktop no explorador de arquivos em [PR](https://github.com/wailsapp/wails/pull/6411) por @taliesin-ai
- Corrigir o não envio dos cabeçalhos do servidor de ativos quando o corpo da resposta está vazio em [PR](https://github.com/wailsapp/wails/pull/6409) por @taliesin-ai
- Respeitar os diretórios de pesquisa de arquivos desktop do XDG ao localizar ícones de aplicativos em [PR](https://github.com/wailsapp/wails/pull/6407) por @taliesin-ai
- Preservar caminhos de arquivos durante a expansão de campos no explorador de arquivos em [PR](https://github.com/wailsapp/wails/pull/6406) por @taliesin-ai
- Evitar a duplicação de destinos de URI expandidos no explorador de arquivos em [PR](https://github.com/wailsapp/wails/pull/6405) por @taliesin-ai
- Confirmar respostas de ativos vazias, bem-sucedidas e sem suporte a busca, em vez de deixá-las pendentes, em [PR](https://github.com/wailsapp/wails/pull/6403) por @taliesin-ai
- Preencher corretamente o prefixo de detecção de MIME ao longo de leituras curtas no servidor de ativos em [PR](https://github.com/wailsapp/wails/pull/6402) por @taliesin-ai
- Evitar a leitura integral dos corpos de resposta em solicitações HEAD sem suporte a busca no servidor de ativos em [PR](https://github.com/wailsapp/wails/pull/6401) por @taliesin-ai
- Servir somente arquivos de índice exatos ao resolver solicitações de diretórios no servidor de ativos em [PR](https://github.com/wailsapp/wails/pull/6400) por @taliesin-ai
- Corrigir notificações do Linux descartadas que eram relatadas como cliques em servidores sem suporte a ações em [PR](https://github.com/wailsapp/wails/pull/6367) por @bercedev
- Encerrar processos de desenvolvimento quando o terminal for fechado em [PR](https://github.com/wailsapp/wails/pull/6371) por @xmasdev
- Manter por um segundo o debounce da reinicialização da barra de tarefas do Windows antes de limpar o sinalizador em [PR](https://github.com/wailsapp/wails/pull/6386) por @taliesin-ai
- Corrigir a perda dos cabeçalhos da resposta de fallback no servidor de ativos em [PR](https://github.com/wailsapp/wails/pull/6399) por @taliesin-ai
- Respeitar as exclusões de cabeçalhos antes da seleção do fallback no servidor de ativos em [PR](https://github.com/wailsapp/wails/pull/6398) por @taliesin-ai
- Liberar referências removidas da fila circular no servidor de ativos para evitar vazamentos de memória em [PR](https://github.com/wailsapp/wails/pull/6397) por @taliesin-ai
- Liberar buffers consumidos pelo detector de conteúdo no servidor de ativos para reduzir o uso de memória em [PR](https://github.com/wailsapp/wails/pull/6396) por @taliesin-ai
- Reconhecer extensões de arquivos de ativos com letras maiúsculas e minúsculas ao servir ativos em [PR](https://github.com/wailsapp/wails/pull/6395) por @taliesin-ai
- Detectar corretamente o tipo de conteúdo no servidor de ativos sem reter nomes de arquivos em [PR](https://github.com/wailsapp/wails/pull/6394) por @taliesin-ai
- Recorrer a um padrão quando uma entrada de desktop não tiver comando em [PR](https://github.com/wailsapp/wails/pull/6392) por @taliesin-ai
- Impedir que o exemplo liquid-glass exiba uma caixa de diálogo antes da inicialização do aplicativo em sistemas que não sejam macOS em [PR](https://github.com/wailsapp/wails/pull/6389) por @taliesin-ai
- Manter as restrições de tamanho da janela quando o modo de tela cheia for iniciado a partir de um tamanho restrito no Windows em [PR](https://github.com/wailsapp/wails/pull/6387) por @taliesin-ai
- Corrigir caixas de diálogo chamadas antes de Run que bloqueavam indefinidamente em vez de retornar em [PR](https://github.com/wailsapp/wails/pull/6390) por @taliesin-ai

## Segurança
- Validar os tamanhos dos ícones gerados antes da alocação da imagem para evitar o uso excessivo de memória em [PR](https://github.com/wailsapp/wails/pull/6478) por @taliesin-ai
- Validar as respostas da API de notas de versão antes do uso em [PR](https://github.com/wailsapp/wails/pull/6461) por @taliesin-ai
- Exigir identificadores de instrução hexadecimais completos no plugin SQLite em [PR](https://github.com/wailsapp/wails/pull/6452) por @taliesin-ai
- Restringir as correspondências de fallback da biblioteca do sistema de arquivos para evitar o carregamento de bibliotecas não pretendidas em [PR](https://github.com/wailsapp/wails/pull/6417) por @taliesin-ai

## v3.1.0-nightly.20261027 - 2026-10-27

## Alterado
- Vincular a página 404 às páginas correspondentes da documentação do Wails v2 em v2.wails.io em [PR](https://github.com/wailsapp/wails/pull/6366) por @taliesin-ai
- Mover o site de documentação do v2 para v2.wails.io em [PR](https://github.com/wailsapp/wails/pull/6365) por @taliesin-ai
- Servir a documentação do v3 em wails.io e a documentação do v2 em v2.wails.io em [PR](https://github.com/wailsapp/wails/pull/6364) por @taliesin-ai
- Atualizar as notas de versão do v3.0.0 e corrigir a alegação do aviso de segurança em [PR](https://github.com/wailsapp/wails/pull/6363) por @taliesin-ai

## Corrigido
- Entregar serialmente os callbacks de ouvintes de eventos do Go na ordem de emissão em [PR](https://github.com/wailsapp/wails/pull/6275) por @DevLumuz
- Sempre anunciar a propriedade Menu da bandeja para que as bandejas criadas após Run não sejam descartadas no Linux em [PR](https://github.com/wailsapp/wails/pull/6284) por @AlbinoGeek
- Reter o trabalho da thread principal enviado antes de Run até o aplicativo iniciar em [PR](https://github.com/wailsapp/wails/pull/6287) por @AlbinoGeek
- Usar a URL correta da API de ativos ao baixar ativos de versões do atualizador no GitHub em [PR](https://github.com/wailsapp/wails/pull/6304) por @gavidroselj
- Corrigir os limites e as posições das telas em configurações de macOS e Linux com DPIs diferentes em [PR](https://github.com/wailsapp/wails/pull/6346) por @julianstorer
- Preservar os corpos das solicitações no Windows quando o WebView2 omitir o cabeçalho Content-Length em [PR](https://github.com/wailsapp/wails/pull/6326) por @strraaa
- Corrigir uma condição de corrida de dados na emissão de hooks de eventos criando um snapshot dos hooks sob bloqueio em [PR](https://github.com/wailsapp/wails/pull/6276) por @DevLumuz
- Corrigir `SystemTray.SetTooltip()` que não fazia nada no Linux: a descrição da dica de ferramenta do StatusNotifierItem agora é definida (#6017)
- Executar a verificação de evidências do changelog no fluxo de trabalho semanal de lançamento em [PR](https://github.com/wailsapp/wails/pull/6374) por @taliesin-ai
- Corrigir a verificação de evidências do changelog de lançamento para mapear execuções para commits de squash merge em [PR](https://github.com/wailsapp/wails/pull/6373) por @taliesin-ai
- Corrigir a detecção do modo escuro no Linux, que retornava o modo claro em um desktop escuro até que o esquema de cores fosse alterado, e `linux:SystemThemeChanged`, que nunca era disparado no backend GTK4. Agora, ambos os backends compartilham uma única implementação, baseada no namespace do portal `org.freedesktop.appearance` (#6072)
- `Window.SetZoom` não eleva mais fatores abaixo de 1.0 para 1.0 em janelas com moldura (Windows, macOS, Linux): WebView2 `ZoomFactor`, WKWebView `magnification` e o nível de zoom do WebKitGTK aceitam qualquer valor positivo, e a opção de janela `Zoom` já repassava esses valores. Janelas sem moldura mantêm o limite mínimo de 1.0 porque suas regiões de arrastar e redimensionar não são ajustadas para o zoom ([#4590](https://github.com/wailsapp/wails/issues/4590)); fatores não positivos, NaN e Inf agora são ignorados em vez de elevados para 1.0 — por @APshenkin

## v3.0.0 - 2026-10-26

## Alterado
- As pré-versões noturnas são publicadas como `vX.Y.0-nightly.YYYYMMDD`, antecipando a próxima versão semanal, e as versões semanais incrementam a versão secundária (`vX.Y.0`), de modo que `@latest` do Go sempre resolva para uma versão semanal ou de correção emergencial; a versão do pacote npm `@wailsio/runtime` sempre corresponde à versão do Wails, com as versões noturnas no dist-tag `nightly`
- Os modelos de projeto fixam `@wailsio/runtime` na versão da CLI que os gerou
- O Wails v3 requer Go 1.26 ou mais recente (Go 1.26.10 ou 1.27.3 e posteriores incluem a correção de segurança `crypto/x509`)

## Corrigido
- Corrigir `wails3 update cli` que não encontrava versões v3: ele considerava apenas tags v2 e instalava a partir do caminho incorreto
- Corrige falhas em aplicativos Windows com `panic: SetTextColor failed` ao desenhar a barra de menus do aplicativo em algumas compilações do Windows: a barra de menus agora é desenhada em um contexto de dispositivo válido ou usa o desenho padrão como alternativa, `UAHMENUITEM.Position` agora corresponde ao `int` de C e o desenho do menu nunca causa panic (#6311)
- Mantém a atualização preparada em segurança enquanto uma reinicialização está pendente: um download ou uma segunda reinicialização iniciados durante o encerramento não podem excluir nem disputar o payload que o auxiliar de substituição aguarda para instalar, por @julianstorer
- Corrige falhas em aplicativos Windows em `applyDarkMode` quando a inicialização de `pkg/w32` descarregava `uxtheme.dll` enquanto seus ponteiros para funções ordinais ainda estavam em uso: o módulo agora permanece carregado durante toda a vida útil do processo (#6356)
- Valida o payload preparado antes da substituição e o aplicativo instalado antes de descartar o backup, para que uma atualização danificada seja revertida em vez de substituir um aplicativo funcional, por @julianstorer

## Removido
- Remove a integração com o Crowdin após a sincronização final das traduções no [PR](https://github.com/wailsapp/wails/pull/5736), por @github-actions[bot]

## Segurança
- Bloqueia mensagens `wails:` reservadas provenientes de frames filhos de origem cruzada no Linux (GTK3 e GTK4), injetando um script de proteção em todos os frames e mitigando GHSA-2p5f-2gx8-mh8m
- Atualiza `github.com/klauspost/compress` para 1.18.7 (GHSA-259r-337f-4rfw) e `golang.org/x/image` para 0.45.0, e compila com Go 1.26.10 para incorporar a correção de `crypto/x509` (GO-2026-6809)
- Atualiza `vite` para 8.2.4 ou posterior no pacote de runtime, no assistente de configuração, nos exemplos e nos modelos de projeto (GHSA-9jrq-w75r-8gcw, GHSA-rq7h-c2jc-7f22, GHSA-vfpm-58rq-9qcg, GHSA-r825-w8mm-5pxw, GHSA-587h-7wmj-h5pm, GHSA-698v-4r3x-3vcv)

### Limitações conhecidas
- **Linux:** a proteção de frames filhos não consegue impedir que um frame filho de *mesma origem* (por exemplo, um iframe `srcdoc` ou `about:blank`) alcance o manipulador de mensagens nativo por meio do frame pai, pois o WebKitGTK não fornece informações de frame ao manipulador. Frames de origem cruzada e frames em sandbox sem `allow-same-origin` são bloqueados. Não incorpore conteúdo não confiável de mesma origem em uma janela do Wails.
- **Android (experimental):** a ponte JavaScript `window.wails` fica exposta a todos os frames da WebView, sem verificação do frame principal nem da origem. O conteúdo de um iframe incorporado, inclusive conteúdo de origem cruzada, pode chamar métodos vinculados e enviar mensagens `wails:` reservadas. Não incorpore conteúdo de terceiros não confiável na WebView do aplicativo. Isso difere do macOS e do iOS, que ignoram mensagens `wails:` provenientes de frames filhos. Uma correção está planejada para 3.0.1.
- **Todas as plataformas:** as chamadas de métodos vinculados provenientes de frames filhos trafegam por HTTP até `/wails/runtime`, e não pelo manipulador de mensagens nativo; portanto, a proteção não as restringe. Considere que qualquer iframe incorporado tem acesso aos seus métodos Go vinculados.

## v3.0.0-beta.34 - 2026-10-20

## Corrigido
- Corrige a WebView no Windows, que não recuperava o foco ao ativar a janela, e a tecla Tab, que ficava presa dentro da WebView em vez de mover o foco para a janela nativa e de volta (#3783)
- Corrige as structs de vtable COM de `internal/webview2/pkg/webview2` (`ICoreWebView2_2`..`_27` e outras interfaces derivadas), que omitiam a cadeia de métodos herdados e faziam as chamadas chegarem ao slot incorreto da vtable (#6129)
- Corrige diálogos de pergunta no Windows que ignoravam rótulos personalizados dos botões e nunca acionavam os callbacks dos botões; os botões personalizados agora são exibidos com uma TaskDialog (com fallback para MessageBox e mapeamento posicional dos callbacks quando comctl32 v6 não está disponível) (#5794)
- Corrige o travamento de `wails3 dev` durante 60s com o Vite: a verificação automática de prontidão do frontend agora usa uma sondagem TCP em vez de uma solicitação HTTP GET que não enviava o cabeçalho `Accept` e recebia um 404 ([#6209](https://github.com/wailsapp/wails/issues/6209))
- Mantém os menus da bandeja do sistema funcionando quando são atualizados enquanto estão abertos no [PR](https://github.com/wailsapp/wails/pull/6301), por @taliesin-ai
- Trata adequadamente erros COM transitórios do WebView2 durante operações do runtime no [PR](https://github.com/wailsapp/wails/pull/6302), por @leaanthony

## v3.0.0-beta.33 - 2026-10-19

## Adicionado
- O gerador de bindings agora emite um módulo `bindings/events` que exporta uma constante de nome tipada para cada evento detectável registrado com `application.RegisterEvent` usando um nome constante em tempo de compilação (por exemplo, `UserLogin` para `user:login`), que pode ser usada com `Events.On`/`Emit` sem o plugin do Vite (#5645)

## Corrigido
- Corrige diálogos de pergunta no Windows que ignoravam rótulos personalizados dos botões e nunca acionavam os callbacks dos botões; os botões personalizados agora são exibidos com uma TaskDialog (com fallback para MessageBox e mapeamento posicional dos callbacks quando comctl32 v6 não está disponível) (#5794)
- Corrige o travamento de `wails3 dev` durante 60s com o Vite: a verificação automática de prontidão do frontend agora usa uma sondagem TCP em vez de uma solicitação HTTP GET que não enviava o cabeçalho `Accept` e recebia um 404 ([#6209](https://github.com/wailsapp/wails/issues/6209))
- Mantém os menus da bandeja do sistema funcionando quando são atualizados enquanto estão abertos no [PR](https://github.com/wailsapp/wails/pull/6301), por @taliesin-ai
- Trata adequadamente erros COM transitórios do WebView2 durante operações do runtime no [PR](https://github.com/wailsapp/wails/pull/6302), por @leaanthony

## v3.0.0-beta.32 - 2026-10-18

## Adicionado
- Adiciona cabeçalhos de seção nativos aos menus do macOS 14+ — consulte o [site da documentação](https://v3.wails.io) (20 páginas atualizadas) no [PR](https://github.com/wailsapp/wails/pull/6297), por @leaanthony

## Corrigido
- Faz SetRelativePosition funcionar em relação à área de trabalho da tela no macOS no [PR](https://github.com/wailsapp/wails/pull/6299), por @leaanthony
- O auxiliar de atualização reinicia o destino inalterado após o tempo limite do processo pai ou uma falha na substituição no [PR](https://github.com/wailsapp/wails/pull/6298), por @taliesin-ai
- Fixa o pnpm na versão 10 para a documentação versionada no [PR](https://github.com/wailsapp/wails/pull/6293), por @taliesin-ai

## v3.0.0-beta.31 - 2026-10-18

## Adicionado
- Expõe a exportação do subcaminho `@wailsio/runtime/cancellable` para que `CancellablePromise` possa ser importado sem avaliar o ponto de entrada raiz do runtime (#6153)

## Alterado
- Passa a exigir Go 1.26 ou posterior: o módulo v3 agora declara `go 1.26.0`, e os projetos gerados por `wails3 init` e `wails3 service init` declaram a mesma versão mínima

## Corrigido
- Corrige a alteração da saída do gerador de bindings ao compilar com Go 1.27, em que tipos `encoding/json` como `Marshaler`, `Unmarshaler` e `RawMessage` se tornavam aliases dos tipos `encoding/json/v2` e `encoding/json/jsontext`: eles voltam a ser renderizados como `any`, e nenhum modelo é emitido para `encoding/json/v2` (#6292)
- Valida os URLs de modelos remotos passados para `wails3 init -t` antes da clonagem: somente URLs `https://`, `http://`, `ssh://`, `git://` e `user@host:path` são aceitos; os valores são passados ao git após `--`; e o git fica restrito a esses transportes durante a clonagem
- Rejeita valores do cabeçalho `x-wails-window-id` fora do intervalo no transporte HTTP do runtime, em vez de truncá-los silenciosamente para o ID de outra janela
- Reforça a extração de arquivos do atualizador contra zip-slip: os caminhos das entradas agora são verificados com `filepath.IsLocal`, além de uma verificação de contenção no caminho combinado, com testes para entradas `..` aninhadas, com barra invertida, com caminho absoluto e que gravam por meio de links simbólicos em arquivos zip e tar.gz
- Corrige a falha na inicialização de AppImages de `wails3 generate appimage` fora do Debian e do Ubuntu com "Unable to spawn a new child process: … WebKitNetworkProcess". O WebKit incluído agora executa os processos auxiliares fornecidos na AppImage. O sandbox bubblewrap do WebKit permanece ativado
- Corrige a falha de compilação de builds de `-tags server` para Android e iOS, resolvendo [#6271](https://github.com/wailsapp/wails/issues/6271)
- Corrige a falha do atualizador ao fazer a substituição no Linux quando o diretório temporário está em outro sistema de arquivos (como um tmpfs em `/tmp`), preparando os arquivos ao lado do binário instalado e copiando-os entre sistemas de arquivos
- Corrige o atualizador, que tentava substituir a montagem somente leitura de uma AppImage em vez do arquivo `.AppImage`
- Corrige a aceitação, no macOS e no iOS, de mensagens reservadas de controle de janela `wails:` provenientes de iframes. Agora, elas são aceitas somente do quadro principal da janela; outras mensagens de iframes ainda chegam a `RawMessageHandler` com suas informações de origem
- Interrompe o envio de `WindowDidUpdate` e `WindowWillUpdate` no macOS: `NSWindowDidUpdateNotification` é acionado a cada passagem do loop de eventos, de modo que a página recebia cerca de 120 chamadas vazias de `evaluateJavaScript:` por segundo; além disso, no macOS, 27.0 fazia com que cada uma delas consumisse a ativação transitória do clique do usuário — `navigator.clipboard.writeText()` e outras APIs que exigem ativação eram rejeitadas com `NotAllowedError` em cliques reais ([#6170](https://github.com/wailsapp/wails/issues/6170)), por @APshenkin
- Corrige o idioma padrão do info.json do Windows para en-US (0409), permitindo que o Explorer exiba os detalhes da versão (#6210)
- Corrige a desativação do renderizador DMA-BUF do Linux em notebooks com uma dGPU NVIDIA junto a uma iGPU AMD ou Intel: no Wayland, agora ele só é desativado quando o WebKit renderiza na NVIDIA (porque o aplicativo foi transferido para a dGPU ou o compositor é executado nela), e não sempre que o módulo nvidia está carregado, o que tornava a renderização extremamente lenta na iGPU; as sessões X11 mantêm o comportamento anterior, pois a GPU de renderização do Xorg é definida por sua configuração
- Trata falhas de renomeação entre dispositivos no atualizador do Linux usando uma cópia intermediária como alternativa, em [PR](https://github.com/wailsapp/wails/pull/6135), por @john-okeefe
- Corrige `WebviewWindow.SetMenu` no Windows, que não redesenhava a barra de menus quando nenhum tema de barra de menus estava configurado, resolvendo [#6103](https://github.com/wailsapp/wails/issues/6103)
- Corrige `WebviewWindow.SetMenu(nil)` no Windows, que causava panic em vez de não realizar nenhuma operação, como no macOS, resolvendo [#6104](https://github.com/wailsapp/wails/issues/6104)
- Preserva as entradas do changelog durante pushes concorrentes usando novas tentativas, em [PR](https://github.com/wailsapp/wails/pull/6094), por @leaanthony
- Corrige a falha de `go mod vendor` com `pattern arm64/WebView2Loader.dll: no matching files found` em todas as plataformas, removendo os embeds que faziam referência a binários nunca distribuídos no módulo; resolve [#5782](https://github.com/wailsapp/wails/issues/5782) e [#5376](https://github.com/wailsapp/wails/issues/5376), em [PR](https://github.com/wailsapp/wails/pull/6031), por @Grantmartin2002
- Evita o encerramento duplicado de janelas no macOS, em [PR](https://github.com/wailsapp/wails/pull/6159), por @charlesbalderson
- Trata WM_CONTEXTMENU para abrir o menu da bandeja do sistema com o botão direito, em [PR](https://github.com/wailsapp/wails/pull/6162), por @wzk2025
- Corrige as cinco variantes de `InvokeSync*`, que bloqueavam para sempre o chamador quando o callback despachado causava panic; resolve [#6107](https://github.com/wailsapp/wails/issues/6107), em [PR](https://github.com/wailsapp/wails/pull/6108), por @RALIST
- Reconstrói com segurança submenus GTK3 destruídos, em [PR](https://github.com/wailsapp/wails/pull/6236), por @leaanthony
- Implementa um shim de encaminhamento de workers no Windows para workers aninhados, em [PR](https://github.com/wailsapp/wails/pull/6261), por @taliesin-ai
- Corrige `OpenFile()` no Windows, que ignorava `PromptForMultipleSelection()` quando `CanChooseDirectories(true)`: agora os usuários podem selecionar várias pastas de uma só vez, como no macOS e no Linux (#6263)
- Corrige argumentos `double` passados por valor para métodos COM do WebView2 no windows/arm64 — `ZoomFactor`, `RasterizationScale`, o `Expires` de cookies, além da escala, do tamanho da página e das margens das configurações de impressão eram afetados nessa plataforma (os setters das configurações de impressão retornavam anteriormente `ErrDoubleArgUnsupported`): agora, as chamadas afetadas passam por um trampoline `runtime.cgocall` que carrega `d0`, algo que o caminho de syscall do Go não consegue fazer (golang.org/issue/62583), por @APshenkin

## Segurança
- O modo servidor (`-tags server`) agora aceita chamadas de `/wails/runtime` somente como solicitações `POST` provenientes da própria origem do servidor ou de uma origem listada em `ServerOptions.WebSocketOriginPatterns`; outros métodos recebem HTTP 405, e outras origens do navegador recebem HTTP 403. Implantações atrás de um proxy reverso que reescreve o cabeçalho `Host` devem adicionar sua origem pública a `WebSocketOriginPatterns`
- O modo servidor (`-tags server`) vinculado a um endereço de loopback (os padrões `localhost`, `127.0.0.0/8` ou `::1`) agora rejeita solicitações com HTTP 403, a menos que o cabeçalho `Host` seja `localhost`, um literal de IP, o host de vinculação ou um host correspondente a `ServerOptions.WebSocketOriginPatterns`, bloqueando ataques de DNS rebinding contra `/wails/runtime` e os endpoints WebSocket
- As notificações do Windows não recorrem mais ao PowerShell quando a API de notificações toast do WinRT falha; em vez disso, o erro é retornado

## v3.0.0-beta.30 - 2026-10-15

## Corrigido
- Corrige o atualizador, que tentava substituir a montagem somente leitura de uma AppImage em vez do arquivo `.AppImage`
- Interrompe o envio de `WindowDidUpdate` e `WindowWillUpdate` no macOS: `NSWindowDidUpdateNotification` é acionado a cada passagem do loop de eventos, de modo que a página recebia cerca de 120 chamadas vazias de `evaluateJavaScript:` por segundo; além disso, no macOS, 27.0 fazia com que cada uma delas consumisse a ativação transitória do clique do usuário — `navigator.clipboard.writeText()` e outras APIs que exigem ativação eram rejeitadas com `NotAllowedError` em cliques reais ([#6170](https://github.com/wailsapp/wails/issues/6170)), por @APshenkin
- Corrige o idioma padrão do info.json do Windows para en-US (0409), permitindo que o Explorer exiba os detalhes da versão (#6210)
- Trata falhas de renomeação entre dispositivos no atualizador do Linux usando uma cópia intermediária como alternativa, em [PR](https://github.com/wailsapp/wails/pull/6135), por @john-okeefe
- Corrige `WebviewWindow.SetMenu` no Windows, que não redesenhava a barra de menus quando nenhum tema de barra de menus estava configurado, resolvendo [#6103](https://github.com/wailsapp/wails/issues/6103)
- Corrige `WebviewWindow.SetMenu(nil)` no Windows, que causava panic em vez de não realizar nenhuma operação, como no macOS, resolvendo [#6104](https://github.com/wailsapp/wails/issues/6104)
- Preserva as entradas do changelog durante pushes concorrentes usando novas tentativas, em [PR](https://github.com/wailsapp/wails/pull/6094), por @leaanthony
- Corrige a falha de `go mod vendor` com `pattern arm64/WebView2Loader.dll: no matching files found` em todas as plataformas, removendo os embeds que faziam referência a binários nunca distribuídos no módulo; resolve [#5782](https://github.com/wailsapp/wails/issues/5782) e [#5376](https://github.com/wailsapp/wails/issues/5376), em [PR](https://github.com/wailsapp/wails/pull/6031), por @Grantmartin2002
- Evita o encerramento duplicado de janelas no macOS, em [PR](https://github.com/wailsapp/wails/pull/6159), por @charlesbalderson
- Trata WM_CONTEXTMENU para abrir o menu da bandeja do sistema com o botão direito, em [PR](https://github.com/wailsapp/wails/pull/6162), por @wzk2025
- Corrige as cinco variantes de `InvokeSync*`, que bloqueavam para sempre o chamador quando o callback despachado causava panic; resolve [#6107](https://github.com/wailsapp/wails/issues/6107), em [PR](https://github.com/wailsapp/wails/pull/6108), por @RALIST

## v3.0.0-beta.29 - 2026-10-13

## Adicionado
- Os atalhos globais no macOS agora consideram o layout do teclado: em layouts que não são QWERTY (AZERTY, QWERTZ, Dvorak, ...), os aceleradores de letras são associados à tecla indicada no teclado por meio de `UCKeyTranslate`, em vez de usarem a posição física fixa do QWERTY.

## Corrigido
- Inicializa o menu nativo antes de defini-lo no Linux, em [PR](https://github.com/wailsapp/wails/pull/6233), por @abhicodes-007
- Resolve regressões de preparação para a disponibilidade geral no macOS no tratamento da saída e no assistente de configuração, em [PR](https://github.com/wailsapp/wails/pull/6237), por @taliesin-ai
- Corrige os wrappers COM `AddRef`/`Release` do WebView2 no Windows, que relatavam como falha um último erro obsoleto da thread. Agora, eles retornam a contagem de referências, portanto uma liberação não pode mais encerrar o aplicativo nem causar falha em uma solicitação de recurso (#3931)
- Corrige `ICoreWebView2Deferral.Release` chamando o slot `AddRef`, que vazava o adiamento de todas as solicitações de ativos adiadas no Windows (#3931)

## v3.0.0-beta.28 - 2026-10-05

## Corrigido
- Mantém a saída do aplicativo macOS fora dos provedores de arquivos no [PR](https://github.com/wailsapp/wails/pull/6225) por @taliesin-ai
- Conclui o assistente de configuração e exibe detalhes do sistema no [PR](https://github.com/wailsapp/wails/pull/6224) por @taliesin-ai
- Melhora o tratamento de erros de recuperação e a limpeza de recursos do WebView2 no Windows no [PR](https://github.com/wailsapp/wails/pull/6222) por @leaanthony
- Melhora a recuperação do WebView2 no Windows para evitar falhas do aplicativo e controladores órfãos no [PR](https://github.com/wailsapp/wails/pull/6220) por @taliesin-ai
- Corrige as orientações sobre autorização de mídia no macOS no [PR](https://github.com/wailsapp/wails/pull/6221) por @taliesin-ai
- Corrige o clique com o botão direito em um ícone da bandeja do sistema no Linux que, além de abrir o menu, executava o manipulador de clique, fazendo uma janela associada alternar de estado a cada clique com o botão direito (#6018)
- Mantém os manipuladores de retorno de chamada do WebView2 ativos até a liberação nativa no Windows no [PR](https://github.com/wailsapp/wails/pull/6184) por @taliesin-ai

## v3.0.0-beta.27 - 2026-10-01

## Corrigido
- Corrige a opção `Permissions`, que era ignorada no macOS 12 e posteriores: o método de captura de mídia `WKUIDelegate` agora está implementado, portanto `PermissionAllow` e `PermissionDeny` se aplicam às solicitações de câmera e microfone da mesma forma que no Linux e no Windows. `NSCameraUsageDescription` / `NSMicrophoneUsageDescription` e, quando o aplicativo estiver em sandbox, os direitos de acesso correspondentes aos dispositivos ainda são necessários

## v3.0.0-beta.26 - 2026-09-25

## Adicionado
- O observador aplica verificações de prontidão HTTP ao servidor de desenvolvimento do frontend no [PR](https://github.com/wailsapp/wails/pull/6150) por @atterpac

## Corrigido
- Recupera-se de falhas do processo do WebView2 em vez de deixar uma janela em branco no [PR](https://github.com/wailsapp/wails/pull/6002) por @taliesin-ai

## v3.0.0-beta.25 - 2026-09-22

## Corrigido
- Encerra `wails3 dev` e seus processos em segundo plano quando o aplicativo principal é encerrado (#6048)

## v3.0.0-beta.24 - 2026-09-20

## Segurança
- Restringe as permissões de GITHUB_TOKEN nos fluxos de trabalho de compilação e publicação para seguir o princípio do privilégio mínimo
- Corrige travessia de caminho e escapes por links simbólicos no middleware de ativos do exemplo de tela
- Mantém funcional a instalação de dependências pelo assistente de configuração, restringindo-a aos comandos de gerenciadores de pacotes aprovados

## v3.0.0-beta.23 - 2026-09-16

## Adicionado

- Permite substituir o ID do aplicativo GTK no Linux no [PR](https://github.com/wailsapp/wails/pull/5972) por @overlordtm

## Corrigido

- A instância única no Linux reivindica o nome do barramento usando o próprio ID do aplicativo no [PR](https://github.com/wailsapp/wails/pull/5973) por @overlordtm
- Cancela solicitações abortadas de esquemas personalizados no Linux no [PR](https://github.com/wailsapp/wails/pull/6126) por @leaanthony

## v3.0.0-beta.22 - 2026-09-14

## Adicionado

- Conclui as traduções e migra a mídia da comunidade — consulte o [site da documentação](https://v3.wails.io) (1812 páginas atualizadas) no [PR](https://github.com/wailsapp/wails/pull/6122) por @leaanthony
- Documenta botões incorporáveis do Wails com vários estilos e tamanhos no [PR](https://github.com/wailsapp/wails/pull/6120) por @leaanthony

## v3.0.0-beta.21 - 2026-09-13

## Adicionado

- Disponibiliza a documentação do Wails v3 com o M-Press no [PR](https://github.com/wailsapp/wails/pull/6116) por @leaanthony

## Corrigido

- Analisa valores JSON de slug no frontmatter MPD para gerar o changelog no [PR](https://github.com/wailsapp/wails/pull/6118) por @leaanthony
- O atualizador limpa as variáveis de ambiente auxiliares e reinicia o destino original após falhas de backup no [PR](https://github.com/wailsapp/wails/pull/6080) por @cnmax
- Inicia o manipulador de sinais padrão durante App.Run no [PR](https://github.com/wailsapp/wails/pull/6098) por @leaanthony
- O menu do Windows trata menus nil, libera recursos substituídos e redesenha a barra de menus no [PR](https://github.com/wailsapp/wails/pull/6112) por @taliesin-ai
- Restaura o empacotamento MSIX para projetos novos que usam configuração YAML compartilhada no [PR](https://github.com/wailsapp/wails/pull/6115) por @leaanthony
- Corrige falhas no carregamento dos bindings JavaScript e TypeScript gerados quando criadores de modelos genéricos fazem referência a declarações auxiliares posteriores e evita estouros de pilha ao criar modelos genéricos mutuamente dependentes (#6062)

## v3.0.0-beta.20 - 2026-09-10

## Alterado

- Atualiza os links da vitrine Clave para o site e o repositório atuais no [PR](https://github.com/wailsapp/wails/pull/6082) por @01xR4in

## Corrigido

- Cancela solicitações de ativos abortadas no Windows, incluindo solicitações de workers, preservando os manipuladores keepalive entre navegações. Encaminha os contextos de solicitação nativos pelo wrapper do aplicativo nas plataformas Apple. (#5963, #5969)
- Preserva as entradas do changelog entre pushes concorrentes com novas tentativas no [PR](https://github.com/wailsapp/wails/pull/6094) por @leaanthony
- Corrige a falha de `go mod vendor` com `pattern arm64/WebView2Loader.dll: no matching files found` em todas as plataformas, removendo os embeds que faziam referência a binários nunca distribuídos no módulo e corrigindo [#5782](https://github.com/wailsapp/wails/issues/5782) e [#5376](https://github.com/wailsapp/wails/issues/5376), no [PR](https://github.com/wailsapp/wails/pull/6031) por @Grantmartin2002

## Removido

- Remove o suporte ao carregador nativo do WebView2, substituído pelo carregador em Go puro. Isso elimina os binários `WebView2Loader.dll` incorporados e a dependência `github.com/jchv/go-winloader`. A tag de compilação `native_webview2loader` continua sendo aceita e não causa mais erros, mas não tem efeito em compilações v3, no [PR](https://github.com/wailsapp/wails/pull/6031) por @Grantmartin2002
- Remove tags de compilação não utilizadas e a opção de FPS do guia da API do macOS no [PR](https://github.com/wailsapp/wails/pull/6097) por @leaanthony

## v3.0.0-beta.19 - 2026-09-09

## Adicionado

- Restringe APIs privadas do macOS por meio de tags de compilação para uso opcional — consulte a [documentação](https://v3.wails.io/features/browser/integration), a [documentação](https://v3.wails.io/features/environment/info), a [documentação](https://v3.wails.io/features/windows/basics), a [documentação](https://v3.wails.io/features/windows/frameless), a [documentação](https://v3.wails.io/features/windows/notch-windows), a [documentação](https://v3.wails.io/features/windows/options), a [documentação](https://v3.wails.io/guides/build/macos), a [documentação](https://v3.wails.io/guides/build/private-macos-apis) e a [documentação](https://v3.wails.io/reference/overview) no [PR](https://github.com/wailsapp/wails/pull/6087) por @leaanthony

## Corrigido

- Rejeita solicitações de runtime maiores que 64 MiB com HTTP 413 no [PR](https://github.com/wailsapp/wails/pull/6091) por @leaanthony

## Segurança

- Reforça as origens do MCP e o acesso remoto com autenticação por token no [PR](https://github.com/wailsapp/wails/pull/6092) por @leaanthony

## v3.0.0-beta.18 - 2026-09-08

## Corrigido

- Corrige o vazamento de memória de Calloc no Linux e no Darwin usando receptores de ponteiro no [PR](https://github.com/wailsapp/wails/pull/6083) por @4RH1T3CT0R7

## v3.0.0-beta.17 - 2026-09-06

## Corrigido

- Windows: um `GetRequest` com falha ou nil no manipulador WebResourceRequested não encerra mais o processo (`log.Fatal` / pânico por desreferência de nil) — em vez disso, a solicitação é descartada e registrada no log no [PR](https://github.com/wailsapp/wails/pull/6006) por @midagedev

## v3.0.0-beta.16 - 2026-08-29

## Alterado

- Solicita a senha de notarização em uma nova janela de terminal no [PR](https://github.com/wailsapp/wails/pull/6029) por @leaanthony

## Corrigido

- Processa corretamente os tipos de clique na bandeja do sistema no macOS no [PR](https://github.com/wailsapp/wails/pull/5919) por @ChewbaccaCookie
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
- Permite que o armazenamento seguro em dispositivos móveis falhe e adota uma postura segura em caso de falha no [PR](https://github.com/wailsapp/wails/pull/5923) por @mortenolsrud
- Executa os hooks de eventos da aplicação mesmo quando nenhum listener está registrado no [PR](https://github.com/wailsapp/wails/pull/5999) por @archy-rock3t-cloud
- Corrige erros de digitação em comentários e na documentação localizada no [PR](https://github.com/wailsapp/wails/pull/6023) por @haoku123
- Remove os binários pré-compilados do macOS adicionados ao repositório em `v3/examples` no [PR](https://github.com/wailsapp/wails/pull/6025) por @4RH1T3CT0R7

## v3.0.0-beta.12 - 2026-08-21

## Adicionado

- Adiciona janelas de notificação para o notch do macOS, com exemplo de ciclo de vida e telemetria — consulte a [documentação](https://v3.wails.io/features/windows/notch-windows) no [PR](https://github.com/wailsapp/wails/pull/6010) por @leaanthony
- Adiciona suporte a janelas NSPanel do macOS, com novas opções e integração nativa — consulte a [documentação](https://v3.wails.io/features/windows/options) no [PR](https://github.com/wailsapp/wails/pull/6008) por @leaanthony

## Corrigido

- Impede que Prepare do SQLite fique bloqueado durante chamadas simultâneas no [PR](https://github.com/wailsapp/wails/pull/5998) por @archy-rock3t-cloud

## v3.0.0-beta.11 - 2026-08-20

## Removido

- Remove da documentação o rastreador de implementação obsoleto no [PR](https://github.com/wailsapp/wails/pull/6005) por @leaanthony

## v3.0.0-beta.10 - 2026-08-19

## Corrigido

- Corrige o host GTK4 no Linux que descartava argumentos de inicialização de protocolo personalizado e associação de arquivos no [PR](https://github.com/wailsapp/wails/pull/6000) por @midagedev
- Processa corretamente linhas excluídas e correções na mesma fonte durante a validação do changelog no [PR](https://github.com/wailsapp/wails/pull/5993) por @taliesin-ai

## v3.0.0-beta.9 - 2026-08-16

## Adicionado

- Adiciona o servidor seguro wails3 mcp para gerenciamento de projetos assistido por agentes no [PR](https://github.com/wailsapp/wails/pull/5896) por @leaanthony
- Adiciona documentação sobre modelos em bindings — consulte a [documentação](https://v3.wails.io/features/bindings/models) no [PR](https://github.com/wailsapp/wails/pull/5988) por @taliesin-ai
- Adiciona suporte à instalação via rpm-ostree em sistemas Linux atômicos no [PR](https://github.com/wailsapp/wails/pull/5987) por @leaanthony
- Adiciona geração e publicação nativas de um gráfico semanal do histórico de estrelas — consulte a [documentação](https://v3.wails.io/credits), a [documentação](https://v3.wails.io/de/credits), a [documentação](https://v3.wails.io/fr/credits), a [documentação](https://v3.wails.io/id/credits), a [documentação](https://v3.wails.io/ja/credits), a [documentação](https://v3.wails.io/ko/credits), a [documentação](https://v3.wails.io/pt/credits), a [documentação](https://v3.wails.io/ru/credits), a [documentação](https://v3.wails.io/zh-cn/credits) e a [documentação](https://v3.wails.io/zh-tw/credits) no [PR](https://github.com/wailsapp/wails/pull/5986) por @leaanthony
- Adiciona o pacote mac exclusivo para Darwin para resolver recursos do bundle da aplicação — consulte a [documentação](https://v3.wails.io/guides/build/macos) no [PR](https://github.com/wailsapp/wails/pull/5965) por @leaanthony
- Adiciona a página de demonstração do Condui e uma entrada no índice — consulte a [documentação](https://v3.wails.io/community/showcase/condui) e a [documentação](https://v3.wails.io/community/showcase) no [PR](https://github.com/wailsapp/wails/pull/5962) por @mgueregath
- Adiciona a página de demonstração do Redis Viewer com capturas de tela e link para o projeto — consulte a [documentação](https://v3.wails.io/community/showcase) e a [documentação](https://v3.wails.io/community/showcase/redisviewer) no [PR](https://github.com/wailsapp/wails/pull/5984) por @redisviewer

## Alterado

- Atualiza as flags da aplicação GTK para G<em>APPLICATION</em>NON_UNIQUE no Linux no [PR](https://github.com/wailsapp/wails/pull/5971) por @overlordtm
- Registra eventos de janela ausentes no nível de depuração em vez de aviso no [PR](https://github.com/wailsapp/wails/pull/5914) por @julianstorer

## Corrigido

- Permite que atalhos de teclado registrados no macOS tenham precedência sobre a webview no [PR](https://github.com/wailsapp/wails/pull/5902) por @julianstorer
- Repara links quebrados da barra lateral da documentação no [PR](https://github.com/wailsapp/wails/pull/5937) por @northes
- Cancela os contextos de solicitação de recursos do macOS e do iOS quando o WebKit interrompe a tarefa correspondente do esquema personalizado (#5963)
- Processa a mensagem WindowSetFullscreenButtonEnabled no [PR](https://github.com/wailsapp/wails/pull/5976) por @archy-rock3t-cloud
- Importa Fragment no modelo preact-ts para corrigir a falha de compilação no [PR](https://github.com/wailsapp/wails/pull/5979) por @haoku123
- Impede que aplicações legadas somente de serviço em GTK3 falhem quando a descoberta de telas é executada antes que uma janela ou tela ativa esteja disponível (#5966)
- Permite que execuções de lançamento com versão explícita prossigam quando o changelog de itens não lançados estiver vazio (#5977)

## Segurança

- Atualiza os lockfiles de nanoid do site para a versão corrigida 3.3.18 a fim de resolver alertas de segurança no [PR](https://github.com/wailsapp/wails/pull/5985) por @taliesin-ai

## v3.0.0-beta.8 - 2026-08-12

## Adicionado

- Adiciona a geração de URLs de documentação às entradas automáticas do changelog no [PR](https://github.com/wailsapp/wails/pull/5957) por @taliesin-ai
- Adiciona Streams: fluxos bidirecionais de bytes entre Go e JavaScript com o modelo de programação WebSocket e sem socket de escuta. Declare um fluxo em Go com `app.HandleStream(name, handler)` e conecte-se pelo frontend com `Stream(name)`, que retorna um objeto no formato de `WebSocket`. O tráfego Go→JS é transportado por uma consulta mantida por janela pelo servidor de assets, e o tráfego JS→Go por um POST normal; nada vincula uma porta TCP nem passa por `evaluateJavaScript`. Em builds de servidor (`-tags server`), o mesmo handler é servido por um WebSocket real, de modo que o código da aplicação é idêntico em todos os builds. por @leaanthony
- Move a entrada do changelog sobre a caixa de mensagens para Não lançado no [PR](https://github.com/wailsapp/wails/pull/5935) por @leaanthony

## Alterado

- Atualiza a geração automática da barra lateral da documentação e a derivação do tipo de autor do blog no [PR](https://github.com/wailsapp/wails/pull/5938) por @leaanthony

## Corrigido

- A inicialização do WebView2 usa um prazo-limite e um loop de mensagens no [PR](https://github.com/wailsapp/wails/pull/5952) por @leaanthony
- O teste de cookies do WebView2 é ignorado na CI, salvo quando habilitado explicitamente, e fixa a execução à thread atual do SO no [PR](https://github.com/wailsapp/wails/pull/5951) por @leaanthony
- Os construtores de menus do Windows restauram os IDs de comando dos itens-pai de submenus no [PR](https://github.com/wailsapp/wails/pull/5944) por @gilad-ch
- Alinha a imagem oficial de compilação cruzada à linha de base de suporte do GTK 4.14+ no Linux (#5928)
- Configura o projeto Xcode para iOS para preservar as flags herdadas do linker e adicionar -ObjC no [PR](https://github.com/wailsapp/wails/pull/5915) por @mortenolsrud
- Corrige a rotatividade excessiva de conexões TCP no proxy de assets `wails3 dev` em frontends grandes, que poderia esgotar as portas efêmeras do host e fazer processos não relacionados falharem com `EADDRNOTAVAIL`
- Enfileira o JavaScript de eventos de cada janela para despacho ordenado e controle de contrapressão no [PR](https://github.com/wailsapp/wails/pull/5934) por @leaanthony

## Removido

- Remove o pipeline de lançamento de binários para desktop: as versões v3 contêm apenas tags, e a CLI `wails3` é instalada com `go install`. Exclui `release-v3.yml` e a etapa nightly que o acionava no [PR](https://github.com/wailsapp/wails/pull/5946) por @leaanthony

## v3.0.0-beta.7 - 2026-08-11

## Adicionado

- Adiciona uma preferência de reprodução automática no macOS para desabilitar a exigência de ação do usuário para reproduzir mídia no [PR](https://github.com/wailsapp/wails/pull/5512) por @Eyalm321
- Move a entrada do changelog sobre a caixa de mensagens para Não lançado no [PR](https://github.com/wailsapp/wails/pull/5935) por @leaanthony

## Alterado

- A animação de zoom do macOS usa CADisplayLink ou NSTimer para obter um desempenho mais suave no [PR](https://github.com/wailsapp/wails/pull/5945) por @savely-krasovsky

## Corrigido

- Configura o projeto Xcode para iOS para preservar as flags herdadas do linker e adicionar -ObjC no [PR](https://github.com/wailsapp/wails/pull/5915) por @mortenolsrud
- Corrige a rotatividade excessiva de conexões TCP no proxy de assets `wails3 dev` em frontends grandes, que poderia esgotar as portas efêmeras do host e fazer processos não relacionados falharem com `EADDRNOTAVAIL`
- Enfileira o JavaScript de eventos de cada janela para despacho ordenado e controle de contrapressão no [PR](https://github.com/wailsapp/wails/pull/5934) por @leaanthony

### Adicionado

- Implementa uma caixa de mensagens FIFO assíncrona genérica para entrega ordenada de eventos no [PR](https://github.com/wailsapp/wails/pull/5851) por @savely-krasovsky e @DevLumuz

## v3.0.0-beta.6 - 2026-08-09

## Adicionado

- Implementa armazenamento limitado no host para eventos grandes demais e entrega ordenada ao JavaScript no [PR](https://github.com/wailsapp/wails/pull/5930) por @leaanthony
- Implementa a animação do ícone no Dock do macOS ao chamar a atenção para uma janela no [PR](https://github.com/wailsapp/wails/pull/5921) por @julianstorer

## Corrigido

- O servidor de assets preserva os erros do detector de tipo de conteúdo e os prefixos não gravados durante o flush no [PR](https://github.com/wailsapp/wails/pull/5931) por @leaanthony
- Impede que aplicações macOS falhem ao substituir o menu da aplicação a partir de um callback do Wails
- Corrige menus nativos ilegíveis no Windows 10 1809 / Windows Server 2019 (build 17763). As exportações de modo escuro do uxtheme estavam condicionadas ao build 18334, portanto a habilitação do modo escuro no nível da aplicação nunca era executada nesses hosts: o fundo do menu era pintado de escuro, mas o Windows continuava desenhando o texto do menu no tema claro, deixando texto escuro sobre fundo escuro. Os ordinais existem desde 17763, portanto a condição agora corresponde a essa versão.
- Corrige `w32.GetStockObject`, que chamava `GetDeviceCaps` em vez de `GetStockObject`, fazendo com que retornasse 0 para todos os objetos predefinidos.
- Melhora o tratamento e o relatório de erros no download do bootstrapper do WebView2 no [PR](https://github.com/wailsapp/wails/pull/5924) por @jannskiee

## v3.0.0-beta.5 - 2026-08-07

## Corrigido

- A ativação de aplicações no macOS agora respeita a política de ativação somente para aplicações regulares no [PR](https://github.com/wailsapp/wails/pull/5897) por @julianstorer
- Protege contra janelas GTK não inicializadas em builds para Linux no [PR](https://github.com/wailsapp/wails/pull/5898) por @julianstorer
- Define uma cor de fundo opaca explícita para janelas WebKit no Linux antes de carregar a URL no [PR](https://github.com/wailsapp/wails/pull/5899) por @julianstorer

## v3.0.0-beta.4 - 2026-08-05

## Alterado

- As tarefas de build do Android agora usam arm64 por padrão, e deploy-emulator seleciona a arquitetura do host no [PR](https://github.com/wailsapp/wails/pull/5890) por @mortenolsrud

## Corrigido

- Preserva o estado de zoom da janela no macOS durante o arrasto e reduz o movimento no [PR](https://github.com/wailsapp/wails/pull/5900) por @leaanthony
- Corrige o build do modo de servidor no Windows adicionando `!server` à restrição de build `webview_window_windows_nonclient.go`

## v3.0.0-beta.3 - 2026-08-03

## Adicionado

- Documenta a conclusão da verificação beta da Fase 10 nos detalhes de implementação no [PR](https://github.com/wailsapp/wails/pull/5881) por @leaanthony

## Corrigido

- Passa o identificador da janela para a API de modo escuro do Windows e valida os argumentos no [PR](https://github.com/wailsapp/wails/pull/5877) por @leaanthony
- Centraliza a resolução do estado dos botões da barra de título do macOS para janelas sem moldura no [PR](https://github.com/wailsapp/wails/pull/5870) por @taliesin-ai
- Impede que o texto do menu nativo fique ilegível quando uma aplicação Windows solicita o modo escuro enquanto o tema de aplicações do Windows está claro. O menu agora usa o fundo nativo claro correspondente até que o Windows possa renderizar texto de menu no modo escuro.
- Corrige menus nativos ilegíveis no Windows 10 1809 / Windows Server 2019 (build 17763). As exportações de modo escuro do uxtheme estavam condicionadas ao build 18334, portanto a habilitação do modo escuro no nível da aplicação nunca era executada nesses hosts: o fundo do menu era pintado de escuro, mas o Windows continuava desenhando o texto do menu no tema claro, deixando texto escuro sobre fundo escuro. Os ordinais existem desde 17763, portanto a condição agora corresponde a essa versão.

## v3.0.0-beta.2 - 2026-08-02

## Alterado

- Promove v3 de alfa para beta
- Documenta os padrões inteligentes da bandeja do sistema e o comportamento de ocultação automática do menu pop-up, com cobertura de regressão para a seleção do handler de clique (#5840).
- O atualizador do GitHub exclui por padrão os assets de instalação do Windows no [PR](https://github.com/wailsapp/wails/pull/5861) por @leaanthony
- Adiciona suporte a cantos arredondados, quadrados e com raio personalizado em janelas sem moldura no macOS no [PR](https://github.com/wailsapp/wails/pull/5866) de @leaanthony

## Corrigido

- Informa os tamanhos atuais das janelas GTK4 e emite eventos de redimensionamento, maximização, minimização e estado de tela cheia a partir da superfície configurada (#5830).
- Corrige uma falha do WebKit no Linux ao enviar Blob ou FormData em solicitações fetch no [PR](https://github.com/wailsapp/wails/pull/5854) de @taliesin-ai
- O shim de fetch passa undefined quando os cabeçalhos de Blob/FormData estão ausentes no [PR](https://github.com/wailsapp/wails/pull/5865) de @leaanthony

## v3.0.0-alpha2.122 - 2026-08-01

## Adicionado

## Alterado

- Adiciona suporte a cantos arredondados, quadrados e com raio personalizado em janelas sem moldura no macOS no [PR](https://github.com/wailsapp/wails/pull/5866) de @leaanthony

## Corrigido

- O shim de fetch passa undefined quando os cabeçalhos de Blob/FormData estão ausentes no [PR](https://github.com/wailsapp/wails/pull/5865) de @leaanthony

## v3.0.0-alpha2.121 - 2026-07-31

## Adicionado

- Adiciona suporte ao empacotamento de DMG no macOS, com novas opções e tarefas de compilação, no [PR](https://github.com/wailsapp/wails/pull/5857) de @leaanthony

## Alterado

- Por padrão, o atualizador do GitHub exclui os artefatos do instalador do Windows no [PR](https://github.com/wailsapp/wails/pull/5861) de @leaanthony

## Corrigido

- Corrige uma falha do WebKit no Linux ao enviar Blob ou FormData em solicitações fetch no [PR](https://github.com/wailsapp/wails/pull/5854) de @taliesin-ai

## v3.0.0-alpha2.120 - 2026-07-31

## Adicionado

- Implementa ações de maximização ou minimização ao clicar duas vezes na barra de título do macOS no [PR](https://github.com/wailsapp/wails/pull/5853) de @taliesin-ai

## Alterado

- Atualiza o tutorial do serviço de QR para usar NewServiceWithOptions e adiciona espaçamento no [PR](https://github.com/wailsapp/wails/pull/5849) de @jeongkyu

## Corrigido

- Mantém a responsividade do WKWebView durante o zoom no macOS no [PR](https://github.com/wailsapp/wails/pull/5856) de @leaanthony
- Corrige as consultas de tamanho de janela do GTK4 e emite eventos de redimensionamento, maximização, minimização e estado de tela cheia a partir do `GdkSurface` configurado.

## v3.0.0-alpha2.119 - 2026-07-27

## Corrigido

- Atualiza a documentação em vários idiomas para incluir diagramas de arquitetura no [PR](https://github.com/wailsapp/wails/pull/5833) de @taliesin-ai

## v3.0.0-alpha2.118 - 2026-07-26

## Adicionado

- Fornece caminhos padrão para as entradas e saídas da geração de ícones no [PR](https://github.com/wailsapp/wails/pull/5825) de @taliesin-ai
- Adiciona módulos de entrada de origem a sideEffects do package.json do runtime no [PR](https://github.com/wailsapp/wails/pull/5797) de @savely-krasovsky
- Adiciona uma seção sobre licença e procedência ao guia de contribuição no [PR](https://github.com/wailsapp/wails/pull/5816) de @taliesin-ai

## Corrigido

- Aplica CSS sem moldura com escopo no GTK4 para remover o raio da borda no [PR](https://github.com/wailsapp/wails/pull/5800) de @savely-krasovsky
- Trata adequadamente as falhas ao obter a posição do cursor no Windows para menus pop-up e enumeração de telas no [PR](https://github.com/wailsapp/wails/pull/5789) de @wayneforrest
- A caixa de diálogo para abrir arquivos do macOS filtra corretamente as extensões e valida os arquivos permitidos pelo sufixo no [PR](https://github.com/wailsapp/wails/pull/5678) de @phergul
- Protege a inicialização do modo escuro do Windows contra chamadas de API nil no [PR](https://github.com/wailsapp/wails/pull/5793) de @roachadam
- Corrige uma falha de compilação de 32 bits no atualizador: a constante `maxArchiveTotalSize` (2 GiB) excedia a capacidade do `int` da plataforma quando passada para `fmt.Errorf` em `GOARCH=386`. Agora, ela tem explicitamente o tipo `int64`.
- Corrige um panic causado por ponteiro nil durante a inicialização quando uma janela usa uma barra de título Dark (ou escura conforme o sistema) em compilações para Windows que não carregam as APIs de modo escuro do uxtheme, como Windows 10 1809 / Windows Server 2019 (compilação 17763). As chamadas a `AllowDarkModeForWindow` na configuração do tema da janela agora verificam nil, como já ocorre em `w32.SetMenuTheme`.

## v3.0.0-alpha2.117 - 2026-07-08

## Adicionado

- Implementa lógica personalizada de teste de acerto para regiões não clientes no Windows no [PR](https://github.com/wailsapp/wails/pull/5462) de @savely-krasovsky

## Alterado

- Configura a detecção de escala do monitor do WebView2 com base em UseVisualHosting no [PR](https://github.com/wailsapp/wails/pull/5761) de @wayneforrest

## v3.0.0-alpha2.116 - 2026-07-07

## Adicionado

- Atualiza a documentação de perguntas frequentes para se concentrar nos recursos e nas orientações do Wails v3 no [PR](https://github.com/wailsapp/wails/pull/5763) de @taliesin-ai

## v3.0.0-alpha2.115 - 2026-07-06

## Corrigido

- Corrige o problema em que `Menu.Update()` não recompilava o menu nativo no GTK4 para Linux (#5659, diagnosticado e corrigido de forma independente por @puneetdixit200 em #5539)
- Corrige uma falha ao enumerar as telas do macOS após uma alteração de monitor, copiando as strings de ID/nome da tela e capturando uma cópia instantânea da contagem (#5565, diagnosticado e corrigido de forma independente por @x-haose em #5584)
- Corrige uma falha no Windows quando `WM_ERASEBKGND` desenha um fundo sólido durante uma transição de minimização/restauração em que `GetClientRect` retorna nil (verificação relatada por @sinspired em #5636)
- Corrige o problema em que o erro dos bindings do frontend sempre era analisado como texto, por @mbaklor em #5690
- Corrige uma falha de compilação no Windows ao usar a tag de compilação `server`, causada pela ausência, nos arquivos da GUI do Windows, da restrição de compilação `!server` que os equivalentes para macOS e Linux já possuem (#5680)

## v3.0.0-alpha2.114 - 2026-07-05

## Adicionado

- Implementa o protocolo Update Manifest e o provedor de endpoint no [PR](https://github.com/wailsapp/wails/pull/5720) de @taliesin-ai

## Alterado

- Incorpora o binding `webview2` ao módulo v3 como `v3/internal/webview2`, removendo o módulo independente, seus fluxos de trabalho de lançamento/sincronização nightly e o malabarismo de versões do go.mod (v3 é seu único consumidor), no [PR](https://github.com/wailsapp/wails/pull/5711) de @taliesin-ai

## Corrigido

- Move a correção da detecção de escala do monitor e da ressincronização do host após alterações de DPI do WebView2 para a seção Não lançado no [PR](https://github.com/wailsapp/wails/pull/5750) de @taliesin-ai
- Atualiza o marshaling COM do WebView2 para parâmetros float64 e BOOL no [PR](https://github.com/wailsapp/wails/pull/5741) de @wayneforrest
- Evita panic e desreferenciamento de nil nas atualizações e na destruição do ícone da bandeja do sistema no Windows no [PR](https://github.com/wailsapp/wails/pull/5703) de @wayneforrest
- Corrige o problema em que janelas ocultas não voltavam a ser ocultadas corretamente no Windows no [PR](https://github.com/wailsapp/wails/pull/5743) de @wayneforrest
- Sincroniza a visibilidade do controlador do WebView2 com a minimização, maximização e restauração da janela no [PR](https://github.com/wailsapp/wails/pull/5742) de @wayneforrest

### Corrigido

- Reativar a detecção da escala do monitor pelo WebView2 e condicionar a ressincronização do host à alteração de DPI em [PR](https://github.com/wailsapp/wails/pull/5734) por @taliesin-ai, com base na correção validada por @randalmurphal, com verificação da causa raiz por @eleclin e testes de hardware por @qq540491950

## v3.0.0-alpha2.113 - 2026-07-04

## Adicionado

- Adicionar um aviso ao compilar um AAB de lançamento sem `ANDROID_KEYSTORE_FILE` definido (o Google Play rejeita pacotes assinados para depuração) e documentar o empacotamento e a assinatura de App Bundles em [PR](https://github.com/wailsapp/wails/pull/5730) por @taliesin-ai
- Adicionar a documentação Por que usar o Wails em vários idiomas em [PR](https://github.com/wailsapp/wails/pull/5739) por @taliesin-ai
- Permitir o mapeamento de Go time.Time para JS Date ou string nos bindings em [PR](https://github.com/wailsapp/wails/pull/5398) por @fbbdev
- Adicionar tarefas de empacotamento de Android App Bundle (AAB) (`bundle`, `bundle:fat`, `assemble:aab`, `assemble:aab:release`) para envio à Play Store — as tarefas de APK permanecem disponíveis para testes locais ou em emuladores em [PR](https://github.com/wailsapp/wails/pull/5728) por @mortenolsrud (corrige [#5726](https://github.com/wailsapp/wails/issues/5726))
- Adicionar destinos de tarefas para dispositivos Android físicos e retomar as permissões de câmera/localização em [PR](https://github.com/wailsapp/wails/pull/5735) por @taliesin-ai

## Alterado

- Atualizar `webview2` para v1.0.28 ([notas de lançamento](https://github.com/wailsapp/wails/releases/tag/webview2%2Fv1.0.28)).
- Atualizar `compileSdk`/`targetSdk` do modelo Android de 34 para 35, exigido pelo Google Play para o envio de novos aplicativos, em [PR](https://github.com/wailsapp/wails/pull/5730) por @taliesin-ai

## Corrigido

- Corrigir máscaras de avatar incorporadas ao sponsorkit em [PR](https://github.com/wailsapp/wails/pull/5745) por @leaanthony
- Corrigir a criação automática de AVDs Android que selecionava a imagem do sistema ou a versão de cmdline-tools incorreta devido à ordenação lexicográfica das versões em [PR](https://github.com/wailsapp/wails/pull/5730) por @taliesin-ai
- Corrigir a sugestão de uma versão desatualizada do Android NDK no assistente de configuração (agora 26.3.11579264, de acordo com o requisito documentado) em [PR](https://github.com/wailsapp/wails/pull/5730) por @taliesin-ai
- Atualizar a documentação em francês sobre SvelteKit e opções em [PR](https://github.com/wailsapp/wails/pull/5744) por @leaanthony
- Corrigir SIGSEGV durante a enumeração de telas no macOS em alterações de monitor em [PR](https://github.com/wailsapp/wails/pull/5516) por @flofreud

## v3.0.0-alpha2.112 - 2026-07-03

## Adicionado

- Adicionar um gerador de SVG de colaboradores baseado em Go e atualizar as páginas de créditos da documentação e do site em [PR](https://github.com/wailsapp/wails/pull/5724) por @taliesin-ai
- Permitir o mapeamento de Go time.Time para JS Date ou string nos bindings em [PR](https://github.com/wailsapp/wails/pull/5398) por @fbbdev

## Alterado

- Substituir o pipeline de imagens de patrocinadores baseado em Node por um gerador em Go em [PR](https://github.com/wailsapp/wails/pull/5719) por @taliesin-ai

## Corrigido

- Corrigir o script de instalação de dependências dos ativos de compilação do Android em [PR](https://github.com/wailsapp/wails/pull/5729) por @taliesin-ai
- Rejeitar o caractere de controle U+0085 (NEXT LINE) em `ValidateAndSanitizeURL`, completando a cobertura de espaços em branco do validador de URLs
- Recalcular o quadro do DWM ao alterar o DPI de janelas sem moldura em [PR](https://github.com/wailsapp/wails/pull/4785) por @leaanthony
- Corrigir a falha na detecção da zona de soltura de DnD em escalas diferentes de 100% no Windows em [PR](https://github.com/wailsapp/wails/pull/4632) por @yulesxoxo
- Adicionar gerenciamento explícito de memória do Objective-C para objetos Cocoa em caixas de diálogo, menus, bandeja e notificações no Darwin em [PR](https://github.com/wailsapp/wails/pull/5714) por @taliesin-ai
- Corrigir bugs do backend CGO no Linux e problemas na bandeja do sistema em [PR](https://github.com/wailsapp/wails/pull/5718) por @taliesin-ai

## v3.0.0-alpha2.111 - 2026-07-01

## Adicionado

- Adicionar o HappyTools à vitrine da comunidade em [PR](https://github.com/wailsapp/wails/pull/5061) por @Aliuyanfeng
- Adicionar suporte à localidade indonésia e documentação abrangente em [PR](https://github.com/wailsapp/wails/pull/5643) por @triadmoko
- Adicionar a opção DisableMenu a WindowsWindow em [PR](https://github.com/wailsapp/wails/pull/4813) por @leaanthony

## Alterado

- Atualizar o modelo do Taskfile e a CLI para despachar tarefas de compilação/empacotamento com GOOS e ARCH em [PR](https://github.com/wailsapp/wails/pull/5617) por @leaanthony

## Corrigido

- Corrigir o problema de agrupamento de janelas em abas no macOS em [PR](https://github.com/wailsapp/wails/pull/5708) por @taliesin-ai

## Removido

- Remover os arquivos MDX traduzidos para alemão das seções de contribuição, recursos e guias em [PR](https://github.com/wailsapp/wails/pull/5702) por @taliesin-ai

## v3.0.0-alpha2.110 - 2026-06-30

## Adicionado

- Implementar o recarregamento e o recarregamento forçado da WebView no macOS e adicionar recuperação após o encerramento do processo WebContent em [PR](https://github.com/wailsapp/wails/pull/5129) por @wayneforrest
- Adicionar documentação abrangente em alemão sobre contribuição, recursos e guias em [PR](https://github.com/wailsapp/wails/pull/5396) por @leaanthony
- Aprimorar as notificações com som, anexos, agendamento e API de atualização em [PR](https://github.com/wailsapp/wails/pull/5333) por @popaprozac

## Corrigido

- Recalcular o quadro do DWM ao alterar o DPI de janelas sem moldura em [PR](https://github.com/wailsapp/wails/pull/4785) por @leaanthony
- Corrigir a falha na detecção da zona de soltura de DnD em escalas diferentes de 100% no Windows em [PR](https://github.com/wailsapp/wails/pull/4632) por @yulesxoxo

## v3.0.0-alpha2.109 - 2026-06-29

## Adicionado

- Adicionar exemplos de código à documentação de EventsEmit em [PR](https://github.com/wailsapp/wails/pull/5026) por @iamhabbeboy
- Adicionar a opção de hospedagem visual do WebView2 no Windows em [PR](https://github.com/wailsapp/wails/pull/5380) por @MerIijn
- Adicionar o Klustr à documentação da vitrine da comunidade em [PR](https://github.com/wailsapp/wails/pull/5536) por @SametKUM
- Adicionar o Kira à vitrine da comunidade, com novas páginas e uma entrada no changelog, em [PR](https://github.com/wailsapp/wails/pull/5685) por @thiennguyen93
- Adicionar uma seção de feedback ao guia do serviço MCP em [PR](https://github.com/wailsapp/wails/pull/5694) por @taliesin-ai

## Alterado

- O modo servidor agora conta com uma compilação de produção de primeira classe, consistente com as tarefas de compilação para desktop (#5693). Por padrão, `task build:server` compila um binário de produção (`-tags server,production`, `-trimpath`, sem símbolos) e aceita `DEV=true` (servidor de desenvolvimento), `OBFUSCATED=true` (garble) e `EXTRA_TAGS`. `task run:server` executa um servidor de desenvolvimento. `Dockerfile.server` / `task build:docker` compilam primeiro o servidor de produção (`-tags server,production`) e o frontend de produção; por padrão, a imagem usa uma compilação estática somente em Go sobre distroless/static, com `CGO_ENABLED`, `GO_IMAGE` e `RUNTIME_IMAGE` expostos como argumentos de compilação substituíveis para aplicativos CGO.

## Corrigido

- Evitar falha ao fechar uma janela com chamadas assíncronas pendentes em [PR](https://github.com/wailsapp/wails/pull/4435) por @leaanthony
- Impedir a ativação da janela ao abrir aplicativos ocultos no Windows em [PR](https://github.com/wailsapp/wails/pull/5249) por @leaanthony
- Garantir que o processamento dos metadados de requisição, da conclusão da resposta e do fluxo do corpo do WebKit seja executado na thread principal do GTK em [PR](https://github.com/wailsapp/wails/pull/5668) por @taliesin-ai
- Corrige o problema em que `Menu.Update()` não reconstruía o menu nativo no Linux com GTK4 (#5659, diagnosticado e corrigido de forma independente por @puneetdixit200 em #5539)
- Corrige uma falha ao enumerar telas do macOS após uma alteração de monitor, copiando as strings de ID/nome das telas e registrando uma captura da contagem (#5565, diagnosticado e corrigido de forma independente por @x-haose em #5584)
- Corrige o conteúdo do WebView2 que encolhia e depois desaparecia ao arrastar uma janela entre monitores com DPIs diferentes no Windows, redefinindo os limites do controlador no manipulador `WM_DPICHANGED`, de forma equivalente à nova sincronização do DPI ao restaurar uma janela minimizada (#5677)

## v3.0.0-alpha2.108 - 2026-06-28

## Adicionado

- Adiciona atalhos de teclado globais (em todo o sistema) por meio de `app.GlobalShortcut` (`Register`, `Unregister`, `UnregisterAll`, `IsRegistered`, `GetAll`). Os atalhos são acionados mesmo quando o aplicativo não está em foco. Implementação nativa para cada plataforma, sem dependências de terceiros: teclas de atalho Carbon no macOS, `RegisterHotKey` no Windows, `XGrabKey` no X11 e a interface de atalhos globais do XDG Desktop Portal no Wayland.
- Adiciona um servidor MCP integrado: um servidor do Model Context Protocol que inicia automaticamente quando o aplicativo é compilado com a tag `mcp`, permitindo que agentes de LLM testem e controlem um aplicativo Wails em execução — controle de janelas, inspeção do DOM, avaliação de JavaScript, chamadas de métodos vinculados, eventos e simulação de entradas de mouse/teclado, exibidas com um cursor animado na tela. Não exige código do usuário: a tag `mcp` é adicionada automaticamente por `wails3 build`/`wails3 dev` quando `WAILS_MCP=1` está definido. Configuração feita inteiramente por variáveis de ambiente (`WAILS_MCP_HOST`, `WAILS_MCP_PORT`, `WAILS_MCP_TIMEOUT`, `WAILS_MCP_HIDE_CURSOR`).

## Corrigido

- Corrige o problema em que `Menu.Update()` não reconstruía o menu nativo no Linux com GTK4 (#5659, diagnosticado e corrigido de forma independente por @puneetdixit200 em #5539)
- Corrige uma falha ao enumerar telas do macOS após uma alteração de monitor, copiando as strings de ID/nome das telas e registrando uma captura da contagem (#5565, diagnosticado e corrigido de forma independente por @x-haose em #5584)

## v3.0.0-alpha2.107 - 2026-06-27

## Adicionado

- Adiciona documentação experimental do Wake com navegação pela barra lateral no [PR](https://github.com/wailsapp/wails/pull/5613) por @leaanthony

## v3.0.0-alpha2.106 - 2026-06-24

## Alterado

- Atualiza `webview2` para v1.0.27.
  - ci(webview2): corrige a compilação de lançamento (compilação cruzada para Windows + go.sum completo) (#5671)\

  **Diferenças completas:** https://github.com/wailsapp/wails/compare/webview2/v1.0.26...webview2/v1.0.27

- Remove go vet da compilação cruzada no fluxo de trabalho de lançamento do webview2 no [PR](https://github.com/wailsapp/wails/pull/5672) por @taliesin-ai
- Atualiza o modelo do OpenRouter usado pelo changelog automático para google/gemini-2.5-flash-lite no [PR](https://github.com/wailsapp/wails/pull/5670) por @taliesin-ai
- Atualiza `webview2` para v1.0.26.

### Correções

- **Recupera-se de erros COM transitórios em tempo de execução em vez de encerrar** (#5658, #5580). Anteriormente, `Chromium.errorCallback` chamava `os.Exit(1)` para *qualquer* erro COM, portanto uma falha temporária recuperável após a inicialização encerrava todo o aplicativo. Agora, os caminhos de tempo de execução (`Resize`/`GetClientRect`, `Navigate`/`NavigateToString`, `Init`, `MessageReceived`, `PutZoomFactor`, `OpenDevToolsWindow`) registram o erro e se recuperam. Em particular, uma mensagem da Web malformada ou não confiável em `MessageReceived` agora é descartada em vez de encerrar o processo. Isso corrige a classe de falhas ao atravessar monitores com DPIs diferentes (#5544, #5650). Os caminhos de criação do ambiente/controlador continuam sendo fatais.\

**Diferenças completas:** https://github.com/wailsapp/wails/compare/webview2/v1.0.25...webview2/v1.0.26

## Corrigido

- Corrige o fluxo de trabalho release-webview2 para processar corretamente arquivos go.sum no [PR](https://github.com/wailsapp/wails/pull/5671) por @taliesin-ai
- Corrige as atualizações de menu do GTK4 no Linux limpando e reconstruindo o menu nativo no [PR](https://github.com/wailsapp/wails/pull/5659) por @taliesin-ai

## v3.0.0-alpha2.105 - 2026-06-21

## Adicionado

- Adiciona `application.System` para detectar a plataforma em tempo de execução a partir de código compartilhado: `System.IsMobile()` (iOS/Android), `System.IsDesktop()` (macOS/Windows/Linux), `System.IsServer()` (a tag de compilação `server`) e `System.IsPlatform(application.PlatformMacOS|PlatformWindows|PlatformLinux|PlatformIOS|PlatformAndroid|PlatformServer)` para testar diretamente um único destino. Ele é compilado para todos os destinos, permitindo criar ramificações sem tags de compilação. Os auxiliares correspondentes para o frontend (`System.IsMobile/IsDesktop/IsIOS/IsAndroid/...`) estão disponíveis em `@wailsio/runtime`
- Adiciona um guia "Como usar outros frameworks de frontend" que mostra como inserir seu próprio projeto Vite em `frontend/` (abrange Solid, Preact, Lit, SvelteKit, Qwik, Angular etc.)
- O assistente `wails3 setup` agora verifica a cadeia de ferramentas para dispositivos móveis (iOS/Android) — Xcode e o runtime do Simulador do iOS, JDK, Android SDK/NDK e emulador — com instalação em um clique e correções copiáveis para a configuração do shell, quando aplicável
- Os projetos gerados incluem uma configuração `frontend/.npmrc` que define um `minimum-release-age` de 7 dias para reduzir a exposição a pacotes recém-publicados (e potencialmente comprometidos), respeitada pelo pnpm e pelo bun e ignorada pelo npm sem causar problemas

## Alterado

- Redesenha todos os modelos iniciais integrados com um novo visual principal de montanhas em neon (Web, iOS e Android)
- **O TypeScript agora é o padrão dos modelos iniciais e assume o nome do modelo sem sufixo.** `wails3 init` (sem `-t`) cria a estrutura de um projeto TypeScript; `-t vanilla`, `-t react`, `-t vue` e `-t svelte` usam TypeScript, com variantes JavaScript em `-t vanilla-js`, `-t react-js`, `-t vue-js` e `-t svelte-js`. Os modelos integrados declaram a linguagem com `typescript:` em `template.yaml`; como alternativa, modelos da comunidade que usam o sufixo `-ts` continuam funcionando
- Redesenha o assistente `wails3 setup` com o tema neon "Wails digital" (efeito vibrante de vidro fosco sobre um cenário de montanhas)

## Corrigido

- Corrige uma falha no Windows ao restaurar um aplicativo que permaneceu minimizado por tempo suficiente para o WebView2 ser suspenso ou para seu processo de renderização/GPU ser reciclado. A nova sincronização do DPI ao minimizar/restaurar (#5544) agora só altera o controlador do WebView2 quando o DPI da janela realmente mudou, evitando chamadas COM fatais a um controlador suspenso no caso comum de restauração com o mesmo DPI (#5605)
- Corrige falhas nativas recorrentes de `SIGABRT`/`SIGSEGV` (normalmente dentro de `g_object_unref` durante o loop principal do GTK) em aplicativos Linux executados por longos períodos e sujeitos a carregamentos frequentes de ativos/mídia. O servidor de ativos concluía `WebKitURISchemeRequest`s a partir de goroutines de trabalho, chamando funções do WebKit2GTK que não são seguras para threads fora da thread principal do GTK; agora, a conclusão (`webkit_uri_scheme_request_finish_with_response`/`finish_error`) é executada na thread principal. Conclui a correção parcial de #5566. Afeta as compilações 6.0 com GTK3 e GTK4/WebKitGTK (#5631, #5557)
- Corrige um `fatal error: invalid pointer found on stack` intermitente em `setupSignalHandlers` no Linux/GTK3. Os IDs de janela passados como `user_data` de sinal eram armazenados em uma variável local Go do tipo `unsafe.Pointer`, portanto o coletor de lixo abortava quando verificava o valor (que não era um ponteiro) durante uma cópia da pilha. Agora, o ID é mantido como inteiro (`uintptr_t`) no lado Go, retroportando para o caminho legado do GTK3 a mesma correção que #4958 aplicou ao caminho do GTK4 (que mudou as funções de sinal C para `uintptr_t` a fim de eliminar erros de `-race`/checkptr) (#5631)

## Removido

- Remove os modelos iniciais `react-swc`, `preact`, `lit`, `solid`, `qwik` e `sveltekit` (e suas variantes `-ts`). O conjunto integrado compatível agora é `vanilla`, `react`, `vue` e `svelte` — todos usam TypeScript por padrão, com variantes JavaScript `-js`. Qualquer outro framework ainda pode ser usado ao [trazer seu próprio frontend](https://v3.wails.io/guides/dev/frontend-frameworks) ou por meio de um modelo personalizado

## v3.0.0-alpha2.104 - 2026-06-18

## Corrigido

- Corrige uma falha do iOS (SIGABRT) quando um método de serviço Go vinculado retorna uma string vazia. O gravador de respostas de ativos do iOS verificava o ponteiro do corpo com `buf != nil` em vez de verificar seu comprimento; por isso, um corpo de comprimento zero fazia `&buf[0]` entrar em pânico. Agora, a verificação usa o comprimento, como nos gravadores para desktop

## v3.0.0-alpha2.103 - 2026-06-15

## Alterado

- Transfere os recursos nativos do iOS e Android para gerenciadores de plataforma: chame-os por meio de `application.IOS.*` e `application.Android.*` (por exemplo, `application.IOS.Haptic("medium")` e `application.Android.Share(payload)`), em vez das antigas funções livres `application.IOS*`/`application.Android*` (#5602)
- Renomeia os eventos da ponte móvel: os eventos multiplataforma agora usam o prefixo `common:*` (por exemplo, `common:haptic` e `common:location`), enquanto os eventos exclusivos de cada plataforma usam `ios:*` / `android:*` (por exemplo, `ios:backgroundTask` e `android:foregroundService`); o prefixo `native:*` não é mais usado (#5602)

## v3.0.0-alpha.102 - 2026-06-14

## Adicionado

- Adiciona o assistente experimental `wails3 setup` para configuração interativa de projetos e verificação de dependências
- Adiciona a opção `--json` a `wails3 doctor` para produzir uma saída legível por máquina
- Adiciona uma seção de status da assinatura ao comando `wails3 doctor`

## Corrigido

- Corrige a detecção do npm no Linux para verificar o PATH além do gerenciador de pacotes

## v3.0.0-alpha.101 - 2026-06-13

## Adicionado

- iOS: diálogos nativos de mensagem (UIAlertController) e diálogos para abrir arquivo, arquivos ou diretório (UIDocumentPickerViewController); os diálogos para salvar retornam um erro explícito
- iOS: suporte à área de transferência por meio de UIPasteboard
- iOS: métricas reais da tela por meio de UIScreen (pontos, pixels, escala e área de trabalho dentro da área segura)
- iOS: compilações para dispositivos (`IOS_PLATFORM=device`), suporte a identidade de assinatura de código, perfil de provisionamento e direitos, empacotamento `.ipa` e `deploy-device` por meio de devicectl
- iOS: versão mínima configurável do iOS (`ios.minIOSVersion` em build/config.yml)
- iOS: `wails3 doctor` informa a disponibilidade do Xcode e do SDK do iOS no macOS
- iOS: eventos do sistema — bateria, rede, tema, bloqueio de tela e pouca memória são expostos como eventos de aplicação `events.IOS.*` e `events.Common.*` independentes de plataforma
- iOS: ponte nativa de recursos móveis (`application.IOS*` exportado) — folha de compartilhamento, abertura de URL, manutenção da tela ativa, lanterna, margens da área segura, brilho, informações do aplicativo, bloqueio de orientação, barra de status, biometria (Face ID/Touch ID), notificações locais e armazenamento seguro no Keychain
- iOS: sensores e hardware — resposta tátil, geolocalização pontual, acelerômetro, proximidade, conversão de texto em fala, informações de armazenamento, estado de energia/bateria, status da rede, margens do teclado e detecção de captura de tela
- iOS: documentação (IOS.md e um guia no site de documentação)
- Android: diálogos nativos de mensagem (AlertDialog) e diálogos para abrir arquivo ou arquivos (Storage Access Framework, importados como cópias no cache); os diálogos para abrir diretório e salvar retornam um erro explícito
- Android: suporte à área de transferência por meio de ClipboardManager
- Android: métricas reais da tela por meio de WindowMetrics/DisplayMetrics (dp, pixels, escala e área de trabalho delimitada pelas barras do sistema)
- Android: métodos de runtime para resposta tátil (`Android.Haptics.Vibrate`), informações do dispositivo (`Android.Device.Info`) e mensagens toast (`Android.Toast.Show`)
- Android: eventos tipados do ciclo de vida (`events.Android.*`, gerados de events.txt), com `ActivityCreated` mapeado para `Common.ApplicationStarted`
- Android: o pipeline de compilação produz APKs de depuração e lançamento instaláveis (`android:run`, `android:package` e `android:package:fat`); por padrão, a versão de lançamento é assinada com o keystore de depuração, ou com um keystore real por meio das variáveis de ambiente `ANDROID_KEYSTORE_*`
- Android: `wails3 doctor` informa o Android SDK, o NDK e o JDK
- Android: eventos do sistema — bateria, rede, tema, bloqueio de tela e pouca memória são expostos como eventos de aplicação `events.Android.*` e `events.Common.*` independentes de plataforma
- Android: ponte nativa de recursos móveis (`application.Android*` exportado) — compartilhamento, abertura de URL, manutenção da tela ativa, lanterna, margens da área segura, brilho, informações do aplicativo, bloqueio de orientação, barra de status, biometria (BiometricPrompt), notificações locais e armazenamento seguro com EncryptedSharedPreferences
- Android: sensores e hardware — resposta tátil, geolocalização pontual, acelerômetro, proximidade, conversão de texto em fala, informações de armazenamento, estado de energia/bateria, status da rede, margens do teclado e bloqueio de captura de tela com FLAG_SECURE
- Android: documentação (ANDROID.md e um guia no site de documentação)
- Exemplo: o exemplo abrangente `mobile` recebe as abas Mobile e Hardware, que demonstram a ponte de recursos nativos no iOS e Android (as abas em formato de pílula quebram em várias linhas)
- Dispositivos móveis: bateria — o acelerômetro, o sensor de proximidade, a lanterna e o relógio periódico do exemplo são pausados quando o aplicativo vai para segundo plano e restaurados quando ele retorna (o Android mantém o processo em execução em segundo plano, e a lanterna é um estado de hardware que persiste no iOS); além disso, os receptores de eventos do sistema Android só são registrados enquanto o aplicativo está em primeiro plano
- iOS: captura com a câmera — `application.IOSCapturePhoto`/`IOSCaptureVideo` (UIImagePickerController → um evento `native:capture` com uma miniatura em base64)
- iOS: execução em segundo plano — `application.IOSBeginBackgroundTask`/`IOSEndBackgroundTask` (uma janela de tarefa em segundo plano do UIApplication) e um `ios.backgroundModes` configurável (build/config.yml), que insere `UIBackgroundModes` como modelo no Info.plist gerado
- Android: captura com a câmera — `application.AndroidCapturePhoto`/`AndroidCaptureVideo` (câmera do sistema por meio de FileProvider → um evento `native:capture`)
- Android: serviço em primeiro plano — `application.AndroidStartForegroundService`/`AndroidStopForegroundService` (um `WailsForegroundService` com uma notificação contínua mantém o processo ativo para trabalhos prolongados em segundo plano)
- Exemplo: uma aba Camera que demonstra a captura de fotos e vídeos e a execução em segundo plano (serviço em primeiro plano no Android e janela de tarefa em segundo plano no iOS)

## Corrigido

- Corrige `getUserMedia`, que sempre falhava com `NotAllowedError` no Linux: o WebKitGTK nega solicitações de permissão que ninguém trata, e o sinal `permission-request` não estava conectado. Agora, a câmera e o microfone são tratados de acordo com um novo mapa multiplataforma `WebviewWindowOptions.Permissions` (`map[PermissionType]Permission`), respeitado tanto no Linux (WebKitGTK) quanto no Windows (WebView2). No Linux, que não tem um prompt nativo, a câmera e o microfone são permitidos por padrão (restaurando `getUserMedia`) e podem ser desativados com `PermissionDeny` (#5552)
- iOS: `GOOS=ios` volta a compilar (`events.IOS` exportado e stubs de nomes de métodos móveis), e as compilações com a tag de produção também compilam (correções de tags de compilação em pkg/application e em vários serviços)
- iOS: os eventos Go→JS e ExecJS agora funcionam — a página não é mais carregada duas vezes na inicialização, e o handshake `wails:runtime:ready` não pode mais ser perdido
- iOS: `ApplicationDidFinishLaunching`/`ApplicationStarted` não apresentam mais uma condição de corrida com a inicialização do aplicativo; a espera fixa de 2 segundos durante a inicialização foi removida
- iOS: corrige um vazamento de string C em cada execução de JavaScript de Go→JS
- iOS: `hasListeners` agora reflete o registro real de ouvintes
- iOS: o registro de depuração do framework é excluído das compilações de produção
- Android: `GOOS=android` volta a compilar — `events.Android` foi definido, o array de ouvintes `events_android.go` com acesso fora dos limites foi removido, o stub de nomes de métodos móveis foi adicionado e os arquivos do Linux para desktop (`linux_cgo.*`, `events_linux.*` e `environment_linux.go`) deixaram de ser incluídos nas compilações para Android
- Android: as vinculações JS→Go agora funcionam — o WebView não consegue entregar corpos POST de `fetch()` a `shouldInterceptRequest`; portanto, as chamadas de runtime são encaminhadas por um transporte JavascriptInterface (`nativeHandleRuntimeCall`), em vez de causarem uma falha devido a um corpo de solicitação nil
- Android: as chamadas de runtime `Screens.*` retornam dados reais — o ScreenManager agora é preenchido na inicialização (ele nunca havia sido conectado, portanto `GetAll` retornava nil)
- Android: o registro de depuração do framework é excluído das compilações de produção e, nas compilações de depuração, é encaminhado pelo logcat com a tag `Wails`
- Android: registro `hasListeners` real, tratamento de referências/exceções JNI e ciclo de vida da página com um único carregamento (sem navegação duplicada)
- Corrige a falha de `wails3 generate bindings` com "Acesso negado" no Windows quando o servidor de desenvolvimento do Vite está em execução, sincronizando os arquivos gerados com o diretório de saída em vez de renomeá-los sobre ele (#5515)
- Corrige uma falha fatal intermitente no macOS ao ler informações da tela após uma mudança de monitor: o ID e o nome da tela armazenavam ponteiros para buffers `UTF8String` com autorelease que podiam ser liberados antes que o Go os copiasse (uso após liberação). Agora, as strings passam por `strdup` e são liberadas após a conversão, e a enumeração de telas é executada em um pool de autorelease explícito, eliminando o vazamento quando chamada por goroutines do Go (#5556)
- Corrige um SIGSEGV intermitente no Linux quando o servidor de ativos fecha um `WebKitURISchemeRequest`: o `g_object_unref` final era executado na goroutine do servidor de ativos, finalizando um GObject do WebKit fora da thread principal do GTK. Agora, o unref é encaminhado ao contexto principal do GTK por meio de `g_main_context_invoke` (#5557)

## v3.0.0-alpha.100 - 2026-06-13

## Adicionado

- Estende `MacWebviewPreferences` com opções adicionais de configuração do WKWebView: `EnableAutoplayWithoutUserAction`, `AllowsAirPlayForMediaPlayback`, `AllowsMagnification`, `JavaScriptCanOpenWindowsAutomatically`, `MinimumFontSize` e `ApplicationNameForUserAgent` (#5549)

## Corrigido

- Corrige a falha de `wails3 generate bindings` com "Acesso negado" no Windows quando o servidor de desenvolvimento do Vite está em execução, sincronizando os arquivos gerados com o diretório de saída em vez de renomeá-los sobre ele (#5561)
- Corrige eventos de redimensionamento do JS que não eram disparados em janelas sem moldura no Linux e a detecção das bordas da barra de rolagem em janelas sem moldura (#5368)
- Corrige a falha do atualizador no Windows com "link inválido entre dispositivos" quando o diretório temporário está em um volume diferente do diretório de instalação (#5560)

## v3.0.0-alpha.99 - 2026-06-10

## Corrigido

- Corrige a falha de `wails3 generate bindings` com "Acesso negado" no Windows quando o servidor de desenvolvimento do Vite está em execução, sincronizando os arquivos gerados com o diretório de saída em vez de renomeá-los sobre ele (#5515)

## v3.0.0-alpha.98 - 2026-06-03

## Corrigido

- Corrige o congelamento da interface do WebKit no Linux quando ociosa (por exemplo, com o inspetor aberto), deixando de forçar `SA_ONSTACK` em `SIGUSR1`, o que interrompia a sincronização da thread de GC do JavaScriptCore (#5527)

## v3.0.0-alpha.97 - 2026-05-31

## Adicionado

- Adiciona uma página de depuração e instruções para trabalhar com `runtime/trace`

## Alterado

- Remove algumas importações desnecessárias de `_ "embed"`, deixando o código um pouco mais organizado

## Corrigido

- Corrige as restrições de largura/altura mínimas que não eram aplicadas após restaurar uma janela maximizada no Windows (#4593)
- Corrige a passagem de cliques do mouse no modo de tela cheia com as opções de janela Frameless + Transparent (#4408)

## v3.0.0-alpha.96 - 2026-05-25

## Adicionado

- Adiciona suporte à ofuscação com Garble ([#4563](https://github.com/wailsapp/wails/issues/4563)): IDs estáveis para métodos de binding, integração com o sistema de compilação/Taskfile (`build --obfuscated --garbleargs`, `generate bindings -obfuscated`) e tags de struct JSON em todos os payloads expostos ao runtime (`EnvironmentInfo`, `OSInfo`, `Screen`, `Rect`, `Point`, `Size`, `Capabilities`), para que o formato de comunicação sobreviva à renomeação de campos exportados feita pelo Garble.

## v3.0.0-alpha.95 - 2026-05-20

## Adicionado

- Adicionada a página ausente sobre a estrutura do projeto

## Alterado

- Documentação: altera alguns diagramas da página de arquitetura para usar diagramas de sequência e proporcionar uma exibição mais clara
- Documentação: inclui uma observação sobre a instalação do D2 como pré-requisito para a execução

## Corrigido

- Corrige `wails3 generate appimage` no GTK4 padrão: agora, o empacotador detecta a pilha GTK no binário antes de procurar os arquivos de runtime, escolhendo `libwebkitgtkinjectedbundle.so` (em `webkitgtk-6.0/`) para compilações GTK4 e `libwebkit2gtkinjectedbundle.so` (em `webkit2gtk-4.1/`) para compilações `-tags gtk3`. A verificação de `.relr.dyn` também consulta `libgtk-4.so.1`, de modo que a remoção de símbolos seja desativada corretamente em toolchains modernas, independentemente da pilha. (#5475)
- Corrige a falha de `wails3 generate appimage` quando invocado com um `-builddir` relativo: agora, o empacotador resolve antecipadamente `-binary`, `-icon`, `-desktopfile`, `-builddir` e `-outputdir` como caminhos absolutos, para que o `s.CD` executado durante o fluxo não interrompa a goroutine de download do AppRun nem a verificação de `ldd` após a cópia.
- Corrige a falha de `wails3 generate appimage` ao mover o AppImage final para `-outputdir` quando o campo `Name=` do arquivo desktop não corresponde ao nome-base do binário: agora, o empacotador força o plugin appimage do linuxdeploy (por meio da variável de ambiente `OUTPUT`) a gravar o AppImage em `<binary>-<arch>.AppImage`, em vez de usar o nome derivado do arquivo desktop.
- Corrige `events.Common.ApplicationStarted`, `Common.ThemeChanged`, `Common.SystemWillSleep` e `Common.SystemDidWake`, que não eram disparados no Linux depois que a pilha GTK4 + WebKitGTK 6.0 passou a ser o padrão na alpha.93. O novo `application_linux.go` `run()` padrão não chamava `setupCommonEvents()` (que encaminha eventos `Linux.*` para seus equivalentes `Common.*`) nem `monitorPowerEvents()`. Agora, o auxiliar de monitoramento de energia via DBus é compartilhado entre os caminhos de compilação GTK3 e GTK4 por meio de `application_linux_dbus.go`. (#5474)

## v3.0.0-alpha.94 - 2026-05-19

## Corrigido

- Corrige `events.Common.ApplicationStarted`, `Common.ThemeChanged`, `Common.SystemWillSleep` e `Common.SystemDidWake`, que não eram disparados no Linux depois que a pilha GTK4 + WebKitGTK 6.0 passou a ser o padrão na alpha.93. O novo `application_linux.go` `run()` padrão não chamava `setupCommonEvents()` (que encaminha eventos `Linux.*` para seus equivalentes `Common.*`) nem `monitorPowerEvents()`. Agora, o auxiliar de monitoramento de energia via DBus é compartilhado entre os caminhos de compilação GTK3 e GTK4 por meio de `application_linux_dbus.go`. (#5474)

## v3.0.0-alpha.93 - 2026-05-17

## Adicionado

- Adiciona `XDG_SESSION_TYPE` à saída de `wails3 doctor` no Linux, por @leaanthony

## Corrigido

- Corrige a falha do menu da janela no Wayland causada pelo acesso do appmenu-gtk-module a uma janela ainda não realizada (#4769), por @leaanthony
- Corrige a falha do aplicativo GTK quando o nome do aplicativo contém caracteres inválidos (espaços, parênteses etc.), por @leaanthony
- Corrige o erro "memória insuficiente" ao inicializar o recurso de arrastar e soltar no Windows (#4701), por @overlordtm
- Corrige uma condição de corrida no armazenamento de callbacks da thread principal causada pelo uso incorreto de RLock para excluir itens do mapa (Linux, macOS, iOS) (#4424), por @leaanthony
- Corrige o tratamento de variáveis ao passar argumentos de linha de comando para tarefas. Agora, as variáveis da CLI especificadas como pares KEY=VALUE são inicializadas corretamente e propagadas durante toda a execução da tarefa.
- Corrige o conflito de NSWindowZoomButton no macOS: agora, `MaximiseButtonState` e `FullscreenButtonState` aplicam o estado mais restritivo tanto na inicialização quanto em tempo de execução; nenhum dos setters pode mais sobrescrever silenciosamente o outro (#5319)
- Corrige um conjunto de bugs preexistentes no caminho de compilação legado do GTK3 (`-tags gtk3`), revelados pelo CodeRabbit no nº 5463: inicializações por associação de arquivo não ignoram mais os manipuladores de inicialização; `getTheme` agora é seguro quanto aos limites e tipos; `appName` não libera mais memória pertencente à GLib; `clipboardGet` não causa mais vazamento do `gchar*` retornado pelo GTK; `Calloc` agora usa receptores de ponteiro (e `NewCalloc` retorna `*Calloc`), para que o pool realmente rastreie as alocações; `zoomOut` usa o recíproco de `zoomInFactor`, em vez de um multiplicador negativo que era limitado a 1.0; `execJS` reutiliza o nome de mundo vazio pré-alocado, em vez de causar o vazamento de um `C.CString("")` por chamada; um `fmt.Println` de desenvolvimento foi removido de `menuItem.setAccelerator`. Resolve o nº 5465.
- Corrige o mesmo vazamento causado pelo receptor de valor de `Calloc` no caminho de compilação padrão do GTK4 (`linux_cgo.go`): receptores de ponteiro + `NewCalloc() *Calloc`, para que as alocações de `c.String(...)` por janela sejam realmente rastreadas e liberadas.

## v3.0.0-alpha.92 - 2026-05-15

## Adicionado

- Modifica os Taskfiles para permitir controlar o gerenciador de pacotes do frontend por meio da opção `PACKAGE_MANAGER`
- Enriquece os dados de modelo com `{{.Opn}}` e `{{.Cls}}` para tornar a criação de modelos de Taskfile mais previsível

## Alterado

- Modifica alguns dos Taskfiles existentes para usar `{{.Opn}} and {{.Cls}}`

## Corrigido

- Corrige uma falha fatal de `concurrent map read and map write` em `linuxSystemTray` quando o menu da bandeja é atualizado enquanto o painel o lê.
- Usa `log` em vez de `fmt` para a saída de erros e rastreamentos de pilha do WebView2, evitando que as mensagens sejam perdidas quando o aplicativo é executado sem um console conectado no Windows.

## v3.0.0-alpha.91 - 2026-05-12

## Alterado

- Atualiza o SVG de patrocinadores em [PR](https://github.com/wailsapp/wails/pull/5414) por `@github-actions[bot]`
- **INCOMPATÍVEL (macOS):** normaliza o sistema de coordenadas do macOS para que `GetScreens`, `Position` e `SetPosition` usem o mesmo espaço — pontos lógicos, eixo Y para baixo, com `(0,0)` no canto superior esquerdo da tela principal. Isso corresponde ao Windows, ao GTK e às APIs públicas do Electron e da Web. Telas fisicamente acima da principal agora relatam `Bounds.Y` negativo (antes positivo), e os valores de `Position()`/`SetPosition()` agora estão em pontos lógicos, em vez de `points × primaryScale`. A conversão de ida e volta `Position()` → `SetPosition()` é preservada; valores absolutos registrados por compilações alfa anteriores ou soluções alternativas calculadas manualmente (por exemplo, multiplicar por `primaryScale` ou inverter Y em relação à altura de uma tela) precisarão ser atualizados. Resolve [nº 5117](https://github.com/wailsapp/wails/issues/5117).

## Corrigido

- Valida defensivamente o nome do sinal DBus e o tamanho do corpo para evitar panics em [PR](https://github.com/wailsapp/wails/pull/5416) por @leaanthony
- Corrige um problema de segurança de memória no processamento de menus GTK no Linux em [PR](https://github.com/wailsapp/wails/pull/5363) por @leaanthony
- Detecta GPUs NVIDIA e desativa o renderizador DMA-BUF no Linux em [PR](https://github.com/wailsapp/wails/pull/5295) por @leaanthony
- Corrige a conversão de Y entre telas de `SetPosition` no macOS: usa a altura da tela principal como referência global para que as janelas sejam posicionadas corretamente em monitores com deslocamento vertical em relação à tela principal em [nº 5117](https://github.com/wailsapp/wails/issues/5117)
- Corrige o modelo de PR do Git para apontar para a URL correta de feedback em [PR](https://github.com/wailsapp/wails/pull/5109) por @wayneforrest
- Corrige uma série de falhas de `SetMenu` na bandeja do sistema do Windows, causadas por uma chamada de sistema `DestroyMenu` incorreta que passava quatro argumentos em vez de um, fazendo com que todas as chamadas retornassem FALSE e não liberassem nada. Também libera identificadores HMENU e HBITMAP (inclusive os alocados em tempo de execução por `MenuItem.SetBitmap`) nas reconstruções do menu, redefine mapas obsoletos de caixas de seleção/botões de opção em `Win32Menu.Update` e remove uma chamada redundante a `Update()` em `systemtray.updateMenu` que duplicava as alocações. Aplicativos de longa duração na bandeja do sistema não causam mais vazamento de objetos GDI/USER a cada reconstrução do menu.

## v3.0.0-alpha.90 - 2026-05-11

## Adicionado

- Adiciona um nome de aplicativo configurável ao User-Agent do WKWebView no macOS em [PR](https://github.com/wailsapp/wails/pull/5261) por @vinhvoit225
- Adiciona a dependência indireta github.com/coder/websocket ao exemplo gin-service em [PR](https://github.com/wailsapp/wails/pull/5400) por @taliesin-ai
- Adiciona suporte à comparação de igualdade profunda aos testes de ativos de compilação em [PR](https://github.com/wailsapp/wails/pull/5402) por @leaanthony

## Alterado

- Consolida a saída da compilação no diretório de ativos em [PR](https://github.com/wailsapp/wails/pull/5401) por @taliesin-ai
- Atualiza o SVG de patrocinadores em [PR](https://github.com/wailsapp/wails/pull/5399) por `@github-actions[bot]`

## Corrigido

- Usa um objeto de notificação para a mensagem de instância única no macOS em [PR](https://github.com/wailsapp/wails/pull/5289) por @overlordtm
- Agrupa callbacks do Windows em lotes para evitar a perda de promises sob carga intensa em [PR](https://github.com/wailsapp/wails/pull/5383) por @taliesin-ai

## v3.0.0-alpha.89 - 2026-05-10

## Adicionado

- Adiciona o job go<em>test</em>results para agregar os resultados dos testes Go em [PR](https://github.com/wailsapp/wails/pull/5316) por @leaanthony

## Alterado

- Divide condicionalmente payloads RPC grandes em requisições POST fragmentadas em [PR](https://github.com/wailsapp/wails/pull/5369) por @leaanthony
- Atualiza o Vite da versão 5.x.x para a 8.0.0 em todos os modelos de frontend em [PR](https://github.com/wailsapp/wails/pull/5386) por @leaanthony
- Migra a configuração da porta do servidor de desenvolvimento do Vite para variáveis de ambiente em [PR](https://github.com/wailsapp/wails/pull/5365) por @leaanthony
- Configura o servidor de desenvolvimento do Vite para escutar em 127.0.0.1 em todos os modelos em [PR](https://github.com/wailsapp/wails/pull/5361) por @leaanthony
- Atualiza o SVG de patrocinadores em [PR](https://github.com/wailsapp/wails/pull/5384) por `@github-actions[bot]`

## Corrigido

- Sanitiza os stubs de modelo de Info.plist durante a atualização dos ativos de compilação em [PR](https://github.com/wailsapp/wails/pull/5312) por @leaanthony
- Corrige o estado obsoleto nos menus do macOS aplicando os modificadores dos itens de menu (`setMenuItemChecked()`, `setMenuItemLabel()`, `setMenuItemDisabled()`, `setMenuItemHidden()`, `setMenuItemTooltip()`) de forma síncrona na thread principal, eliminando a condição de corrida de `dispatch_async` que fazia os menus renderizarem o estado anterior quando eram reabertos rapidamente (nº 5002)
- Ignora arquivos `*_test.go` no modo de desenvolvimento para evitar recompilações desnecessárias em [PR](https://github.com/wailsapp/wails/pull/5203) por @leaanthony
- Impede uma falha de segmentação em Menu.Update() quando o aplicativo não está em execução em [PR](https://github.com/wailsapp/wails/pull/5291) por @wucm667
- Usa lastSizeWParam para controlar os redesenhos da barra de menus no Windows em [PR](https://github.com/wailsapp/wails/pull/5382) por @taliesin-ai

## v3.0.0-alpha.88 - 2026-05-09

## Alterado

- Altera HiddenOnTaskbar para usar WS<em>EX</em>TOOLWINDOW em [PR](https://github.com/wailsapp/wails/pull/5371) por @leaanthony
- Reordena as dependências e remove a diretiva replace de webview2 do go.mod em [PR](https://github.com/wailsapp/wails/pull/5370) por @atterpac
- Atualiza o SVG de patrocinadores em [PR](https://github.com/wailsapp/wails/pull/5358) por `@github-actions[bot]`

## Corrigido

- Remove aliases genéricos de indireção e consolida os tipos de chave de mapas em [PR](https://github.com/wailsapp/wails/pull/5331) por @fbbdev

## Removido

- Exclui o fluxo de trabalho PR-master, removendo a documentação, os testes de Go e os testes ignorados no [PR](https://github.com/wailsapp/wails/pull/5377), por @leaanthony

## v3.0.0-alpha.87 - 2026-05-07

## Adicionado

- Adiciona a documentação em coreano do Wails v3 no [PR](https://github.com/wailsapp/wails/pull/5352), por @leaanthony
- Adiciona a documentação em francês sobre instalação e início rápido no [PR](https://github.com/wailsapp/wails/pull/5354), por @leaanthony
- Adiciona a documentação em português sobre início rápido, conceitos e comunidade no [PR](https://github.com/wailsapp/wails/pull/5355), por @leaanthony

## v3.0.0-alpha.86 - 2026-05-06

## Adicionado

- Adiciona a localização da documentação em francês no [PR](https://github.com/wailsapp/wails/pull/5328), por @leaanthony
- Adiciona o idioma alemão ao site da documentação no [PR](https://github.com/wailsapp/wails/pull/5343), por @leaanthony

## Alterado

- Registra todos os 8 idiomas traduzidos na configuração da documentação no [PR](https://github.com/wailsapp/wails/pull/5347), por @leaanthony
- Atualiza vários arquivos relacionados ao Windows para o WebView2 no [PR](https://github.com/wailsapp/wails/pull/5317), por @leaanthony

## Corrigido

- Separa o despacho de caixas de diálogo entre GTK3 e GTK4 no Linux no [PR](https://github.com/wailsapp/wails/pull/5340), por @leaanthony
- Garante que os callbacks das caixas de diálogo sejam executados na thread do GTK, corrigindo falhas de segmentação no [PR](https://github.com/wailsapp/wails/pull/5339), por @leaanthony

## v3.0.0-alpha.85 - 2026-05-05

## Adicionado

- Adiciona ao repositório a URL do modelo de PR no [PR](https://github.com/wailsapp/wails/pull/5179), por @leaanthony
- Adiciona a documentação em alemão do Wails v3 no [PR](https://github.com/wailsapp/wails/pull/5330), por @leaanthony

## v3.0.0-alpha.84 - 2026-05-03

## Adicionado

- Adiciona uma opção para impedir que a tecla Escape encerre o modo de tela cheia no macOS no [PR](https://github.com/wailsapp/wails/pull/5307), por @leaanthony
- Adiciona uma opção para impedir que a tecla Escape encerre o modo de tela cheia no macOS no [PR](https://github.com/wailsapp/wails/pull/5310), por @leaanthony
- Adiciona a documentação da Pausa à vitrine da comunidade no [PR](https://github.com/wailsapp/wails/pull/5288), por @yuseferi

## Alterado

- Atualiza o SVG dos patrocinadores no [PR](https://github.com/wailsapp/wails/pull/5308), por `@github-actions[bot]`
- Atualiza o comando de geração de ícones para lidar com plataformas não compatíveis no [PR](https://github.com/wailsapp/wails/pull/5309), por @leaanthony
- Substitui a API booleana de tela cheia pelo ButtonState de três estados e implementa os bindings de plataforma no [PR](https://github.com/wailsapp/wails/pull/5224), por @leaanthony

## Corrigido

- Protege as operações de foco do WebView2 contra o estado nulo do controlador no [PR](https://github.com/wailsapp/wails/pull/5315), por @leaanthony
- Atualiza o fluxo de trabalho do GitHub Actions para referenciar corretamente a branch de base do PR no [PR](https://github.com/wailsapp/wails/pull/5313), por @leaanthony
- Ignora arquivos `*_test.go` no modo de desenvolvimento para evitar recompilações desnecessárias no [PR](https://github.com/wailsapp/wails/pull/5203), por @leaanthony
- Evita uma falha de segmentação em Menu.Update() quando o aplicativo não está em execução no [PR](https://github.com/wailsapp/wails/pull/5291), por @wucm667

## v3.0.0-alpha.83 - 2026-05-02

## Adicionado

- Adiciona a flag InstallScope e uma opção de compilação para instalação por máquina/usuário no [PR](https://github.com/wailsapp/wails/pull/5094), por @symball
- Adiciona ao BrowserWindow um método SetScreen sem operação para satisfazer a interface Window no [PR](https://github.com/wailsapp/wails/pull/5294), por @leaanthony

## Corrigido

- Detecta GPUs NVIDIA e desativa o renderizador DMA-BUF no Linux no [PR](https://github.com/wailsapp/wails/pull/5295), por @leaanthony
- Corrige o modelo de PR do Git para apontar para a URL correta de feedback no [PR](https://github.com/wailsapp/wails/pull/5109), por @wayneforrest
- Corrige um conjunto de falhas `SetMenu` na bandeja do sistema do Windows causadas por uma chamada de sistema `DestroyMenu` incorreta, que passava quatro argumentos em vez de um; por isso, todas as chamadas retornavam FALSE e não liberavam nada. Também libera os identificadores HMENU e HBITMAP (inclusive os alocados em tempo de execução por meio de `MenuItem.SetBitmap`) nas reconstruções do menu, redefine mapas obsoletos de caixas de seleção/botões de opção em `Win32Menu.Update` e remove uma chamada redundante de `Update()` em `systemtray.updateMenu` que duplicava as alocações. Aplicativos de bandeja do sistema executados por longos períodos não vazam mais objetos GDI/USER a cada reconstrução do menu.

## v3.0.0-alpha.82 - 2026-05-01

## Corrigido

- Corrige a geração do arquivo desktop para processar corretamente o nome do desktop no [PR](https://github.com/wailsapp/wails/pull/5232), por @leaanthony

## v3.0.0-alpha.81 - 2026-04-30

## Alterado

- Ajusta o agendamento da versão nightly para 15:00 UTC no [PR](https://github.com/wailsapp/wails/pull/5286), por @leaanthony

## Corrigido

- Corrige valores de Screen Bounds, WorkArea e Size reduzidos pela metade em Macs com tela Retina -  (#5168)

## v3.0.0-alpha.80 - 2026-04-29

## Alterado

- Atualiza as dependências da documentação e os carregadores de coleções de conteúdo no [PR](https://github.com/wailsapp/wails/pull/5285), por @leaanthony

## v3.0.0-alpha.79 - 2026-04-29

## Adicionado

- Concede a permissão actions: write ao job trigger-release no [PR](https://github.com/wailsapp/wails/pull/5270), por @leaanthony

## Alterado

- Define a branch master como padrão da tarefa de lançamento e atualiza a redação do changelog no [PR](https://github.com/wailsapp/wails/pull/5283), por @leaanthony
- Atualiza o fluxo de trabalho auto-changelog para usar a versão mais recente no [PR](https://github.com/wailsapp/wails/pull/5282), por @leaanthony
- Melhora a eficiência dos fluxos de trabalho adicionando filtros de caminho e removendo fluxos de trabalho obsoletos no [PR](https://github.com/wailsapp/wails/pull/5280), por @leaanthony
- Atualiza a documentação para referenciar a branch master nos links de exemplos no [PR](https://github.com/wailsapp/wails/pull/5274), por @leaanthony
- Atualiza a documentação e os exemplos de v3 no [PR](https://github.com/wailsapp/wails/pull/5272), por @leaanthony

## Corrigido

- Aprimora o proxy reverso com lógica de repetição e uso forçado de IPv4 durante o desenvolvimento no [PR](https://github.com/wailsapp/wails/pull/5265), por @AkagiYui
- Reescreve o fluxo de trabalho que aciona o changelog de alterações ainda não lançadas no [PR](https://github.com/wailsapp/wails/pull/5281), por @leaanthony

## Removido

- Remove scripts de teste do shell usados para várias finalidades de teste no [PR](https://github.com/wailsapp/wails/pull/5267), por @leaanthony
- Exclui o fluxo de trabalho de implantação da documentação v3-alpha e o registro CNAME no [PR](https://github.com/wailsapp/wails/pull/5266), por @leaanthony

### Adicionado

- Adiciona a entrada Roteamento do frontend à navegação da barra lateral no [PR](https://github.com/wailsapp/wails/pull/5196), por @leaanthony
- Adiciona um guia de roteamento do frontend com recomendações específicas para cada framework no [PR](https://github.com/wailsapp/wails/pull/5185), por @leaanthony
- Adicionar suporte a folhas modais (macOS)
- Atualizar a versão do ghw para oferecer melhor suporte a dispositivos Apple, por @leaanthony (#4977)
- Adicionar o método `GetBadge` ao serviço do Dock
- Adicionar a opção `-tags` ao comando `wails3 build` para passar tags de compilação personalizadas do Go (por exemplo, `wails3 build -tags gtk4`) (#4957)
- Adicionar documentação sobre a geração automática de enums no gerador de bindings, incluindo uma página dedicada a Enums e a navegação na barra lateral (#4972)
- Adicionar a opção `-tags` ao comando `wails3 build` para passar tags de compilação personalizadas do Go (por exemplo, `wails3 build -tags gtk4`) (#4957)
- Adicionar exemplos de APIs Web em `v3/examples/web-apis/` que demonstram APIs de navegador 41, incluindo armazenamento (localStorage, sessionStorage, IndexedDB e Cache API), rede (Fetch, WebSocket, XMLHttpRequest, EventSource e Beacon), mídia (Canvas, WebGL, Web Audio, MediaDevices, MediaRecorder e Speech Synthesis), dispositivo (Geolocation, Clipboard, Fullscreen, Device Orientation, Vibration e Gamepad), desempenho (Performance API, Mutation Observer e Intersection/Resize Observer), interface (Web Components, Pointer Events, Selection, Dialog e Drag and Drop), entre outras
- Adicionar um exemplo de verificador de compatibilidade das APIs do WebView (`v3/examples/webview-api-check/`), que testa mais de 200 APIs de navegador em diferentes plataformas
- Adicionar o pacote `internal/libpath` para localizar caminhos de bibliotecas nativas no Linux, com busca paralela, cache e suporte a Flatpak/Snap/Nix
- **WIP:** adicionar suporte experimental a WebKitGTK 6.0 / GTK4 no Linux, disponível por meio de `-tags gtk4` (GTK3/WebKit2GTK 4.1 continua sendo o padrão)
- Observação: em gerenciadores de janelas em mosaico (por exemplo, Hyprland e Sway), as operações de minimizar/maximizar podem não funcionar como esperado, pois o gerenciador de janelas controla a geometria da janela
- Adicionar à documentação como usar **manipuladores de uso único** em **Escuta de eventos em JavaScript**, por @AbdelhadiSeddar
- Adicionar a opção `UseApplicationMenu` a `WebviewWindowOptions`, permitindo que janelas no Windows/Linux herdem o menu do aplicativo definido por meio de `app.Menu.Set()`, por @leaanthony
- Adicionar suporte ao uso de arquivos `.icon` (formato do Apple Icon Composer) para gerar ícones Liquid Glass e catálogos de recursos (macOS) (#4934), por @wimaha
- Adicionar um modo de servidor experimental para implantações headless/Web (`-tags server`). Permite executar aplicativos Wails como servidores HTTP sem dependências de interface gráfica nativa. Compile com `wails3 task build:server`. Consulte `examples/server` para obter detalhes.
- Adicionar o pacote `internal/libpath` para localizar caminhos de bibliotecas nativas no Linux, com busca paralela, cache e suporte a Flatpak/Snap/Nix
- Adicionar a opção `CollectionBehavior` a `MacWindow` para controlar o comportamento das janelas entre Spaces do macOS e o modo de tela cheia (#4756), por @leaanthony
- Adicionar testes unitários para pkg/application, por @leaanthony
- Adicionar suporte a protocolos personalizados ao empacotamento MSIX, por @leaanthony
- Adicionar detecção do ambiente de desktop no Linux [PR #4797](https://github.com/wailsapp/wails/pull/4797)
- Adicionar o método `Window.Print()` ao runtime JavaScript para abrir a caixa de diálogo de impressão pelo frontend (#4290), por @leaanthony
- Adicionar `XDG_SESSION_TYPE` à saída de `wails3 doctor` no Linux, por @leaanthony
- Adicionar mais eventos de alteração de carregamento do WebKit2 no Linux: `WindowLoadStarted`, `WindowLoadRedirected`, `WindowLoadCommitted` e `WindowLoadFinished` (#3896), por @leaanthony
- Adicionar `XDG_SESSION_TYPE` à saída de `wails3 doctor` no Linux, por @leaanthony
- Gerar o arquivo `.desktop` durante a compilação no Linux, não apenas no empacotamento (#4575)
- Adicionar documentação sobre as dependências de runtime no Linux, com nomes de pacotes específicos por distribuição e exemplos de empacotamento com nfpm (#4339), por @leaanthony
- Adicionar informações sobre a versão do driver NVIDIA à saída de `wails3 doctor` no Linux, por @leaanthony
- Adicionar a origem ao manipulador de mensagens brutas, por @APshenkin em [PR](https://github.com/wailsapp/wails/pull/4710)
- Adicionar suporte a links universais no macOS, por @APshenkin em [PR](https://github.com/wailsapp/wails/pull/4712)
- Refatorar a camada de transporte de bindings, por @APshenkin em [PR](https://github.com/wailsapp/wails/pull/4702)
- Adicionar identificadores aria-label aos modelos helloworld para facilitar os testes do aplicativo de exemplo por clientes de teste Appium, por @chinenual em [PR](https://github.com/wailsapp/wails/pull/4760)
- Adicionar a origem ao manipulador de mensagens brutas, por @APshenkin em [PR](https://github.com/wailsapp/wails/pull/4710)
- Adicionar suporte a links universais no macOS, por @APshenkin em [PR](https://github.com/wailsapp/wails/pull/4712)
- Refatorar a camada de transporte de bindings, por @APshenkin em [PR](https://github.com/wailsapp/wails/pull/4702)
- Eventos tipados, por @fbbdev e @ianvs em [#4633](https://github.com/wailsapp/wails/pull/4633)
- Adicionar o exemplo `systray-clock`, que mostra um ícone da área de notificação sem janela, com atualizações dinâmicas da dica de ferramenta (#4653).
- Adicionar modelo de protocolo NSIS para Windows, por @Tolfx em #4510
- Adicionar testes para build-assets, por @Tolfx em #4510
- macOS: exibir os controles nativos da janela na barra de menus em [#4588](https://github.com/wailsapp/wails/pull/4588), por @nidib
- Adicionar um serviço do Dock do macOS para ocultar/exibir o ícone do aplicativo no Dock, por @popaprozac em [PR](https://github.com/wailsapp/wails/pull/4451)
- Adicionar um serviço do Dock do macOS para ocultar/exibir o ícone do aplicativo no Dock, por @popaprozac em [PR](https://github.com/wailsapp/wails/pull/4451)
- Adicionar suporte nativo ao efeito Liquid Glass no macOS com NSGlassEffectView (macOS 15.0+) e fallback para NSVisualEffectView, incluindo opções abrangentes de personalização de materiais, por @leaanthony em [#4534](https://github.com/wailsapp/wails/pull/4534)
- Sanitização de URLs do navegador, por @leaanthony em [#4500](https://github.dev/wailsapp/wails/pull/4500). Baseada em [#4484](https://github.com/wailsapp/wails/pull/4484), por @APShenkin.
- Adicionar proteção de conteúdo no Windows/Mac por [@leaanthony](https://github.com/leaanthony), com base no trabalho original de [@Taiterbase](https://github.com/Taiterbase) neste [PR](https://github.com/wailsapp/wails/pull/4241)
- Adicionar suporte para passar variáveis da CLI aos comandos do Task por meio dos aliases `wails3 build` e `wails3 package` (#4422), por @leaanthony em [PR](https://github.com/wailsapp/wails/pull/4488)
- Suporte a zonas de soltura, com o evento fornecendo os dados do elemento solto, por [@atterpac](https://github.com/atterpac) em [#4318](https://github.com/wailsapp/wails/pull/4318)
- Adicionar `AdditionalLaunchArgs` às opções de `WindowsWindow` para permitir que argumentos adicionais de linha de comando sejam passados ao navegador WebView2, em [PR](https://github.com/wailsapp/wails/pull/4467)
- Adicionar a execução automática de go mod tidy após wails init, por [@triadmoko](https://github.com/triadmoko) em [PR](https://github.com/wailsapp/wails/pull/4286)
- Recurso Snap Assist do Windows, por @leaanthony em [PR](https://github.dev/wailsapp/wails/pull/4463)
- Adicionar `AdditionalLaunchArgs` às opções de `WindowsWindow` para permitir que argumentos adicionais de linha de comando sejam passados ao navegador WebView2, em [PR](https://github.com/wailsapp/wails/pull/4467)
- Adicionar a execução automática de go mod tidy após wails init, por [@triadmoko](https://github.com/triadmoko) em [PR](https://github.com/wailsapp/wails/pull/4286)
- Recurso Snap Assist do Windows, por @leaanthony em [PR](https://github.dev/wailsapp/wails/pull/4463)
- Adicionar a implementação de `getAccentColor` para Windows, por [@almas-x](https://github.com/almas-x) em [PR](https://github.com/wailsapp/wails/pull/4427)
- Adicionar a implementação de `getAccentColor` para Windows, por [@almas-x](https://github.com/almas-x) em [PR](https://github.com/wailsapp/wails/pull/4427)
- Menus e barra de menus com tema escuro no Windows. Por @leaanthony em [a29b4f0861b1d0a700e9eb213c6f1076ec40efd5](https://github.com/wailsapp/wails/commit/a29b4f0861b1d0a700e9eb213c6f1076ec40efd5)
- Renomeia os serviços integrados para tornar os bindings JS/TS mais claros, por @popaprozac em [PR](https://github.com/wailsapp/wails/pull/4405)
- `app.Env.GetAccentColor` para obter a cor de destaque do sistema do usuário. Funciona no macOS. Por [@etesam913](https://github.com/etesam913)
- Adiciona a API `window.ToggleFrameless()`, por [@atterpac](https://github.com/atterpac) em [#4137](https://github.com/wailsapp/wails/pull/4137)
- Adiciona dependências de compilação específicas de cada distribuição Linux, por @leaanthony em [PR](https://github.com/wailsapp/wails/pull/4345)
- Adiciona um guia de bindings, por @atterpac em [PR](https://github.com/wailsapp/wails/pull/4404)
- **Infraestrutura de testes organizada**: move os arquivos de teste do Docker para o diretório dedicado `test/docker/`, com imagens otimizadas e maior confiabilidade de compilação, por [@leaanthony](https://github.com/leaanthony) em [#4359](https://github.com/wailsapp/wails/pull/4359)
- **Padrões aprimorados de gerenciamento de recursos**: adiciona aos exemplos a limpeza adequada de manipuladores de eventos e o gerenciamento de goroutines com reconhecimento de contexto, por [@leaanthony](https://github.com/leaanthony) em [#4359](https://github.com/wailsapp/wails/pull/4359)
- Adiciona suporte à compilação de AppImages para aarch64, por [@AkshayKalose](https://github.com/AkshayKalose) em [#3981](https://github.com/wailsapp/wails/pull/3981)
- Adiciona uma seção de diagnóstico a `wails doctor`, por [@leaanthony](https://github.com/leaanthony)
- Adiciona a janela ao contexto ao chamar um método de serviço, por [@leaanthony](https://github.com/leaanthony)
- Adiciona o exemplo `window-call` para demonstrar como identificar qual janela está chamando um serviço, por [@leaanthony](https://github.com/leaanthony)
- Novo guia de menus, por [@leaanthony](https://github.com/leaanthony)
- Melhora o tratamento de panics, por [@leaanthony](https://github.com/leaanthony)
- Novo guia de menus, por [@leaanthony](https://github.com/leaanthony)
- Adiciona comentários de documentação à API de serviços, por [@fbbdev](https://github.com/fbbdev) em [#4024](https://github.com/wailsapp/wails/pull/4024)
- Adiciona a função `application.NewServiceWithOptions` para inicializar serviços com configurações adicionais, por [@leaanthony](https://github.com/leaanthony) em [#4024](https://github.com/wailsapp/wails/pull/4024)
- Aprimora o controle de menus, por [@FalcoG](https://github.com/FalcoG) e [@leaanthony](https://github.com/leaanthony) em [#4031](https://github.com/wailsapp/wails/pull/4031)
- Mais documentação, por [@leaanthony](https://github.com/leaanthony)
- Adiciona suporte ao cancelamento de eventos em listeners de eventos padrão, por [@leaanthony](https://github.com/leaanthony)
- Adiciona suporte a `Hide`, `Show` e `Destroy` na bandeja do sistema, por [@leaanthony](https://github.com/leaanthony)
- Adiciona suporte a `SetTooltip` na bandeja do sistema, por [@leaanthony](https://github.com/leaanthony). Ideia original de [@lujihong](https://github.com/wailsapp/wails/issues/3487#issuecomment-2633242304)
- Informa o caminho do pacote nos avisos do gerador de bindings sobre tipos sem suporte, por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Adiciona ao gerador de bindings suporte a aliases genéricos, por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Adiciona ao gerador de bindings suporte à flag JSON `omitzero`, por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Adiciona a diretiva `//wails:ignore` para impedir a geração de bindings para métodos de serviço selecionados, por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Adiciona a diretiva `//wails:internal` a serviços e modelos para permitir tipos exportados em Go, mas não em JS/TS, por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Adiciona ao gerador de bindings suporte a constantes de tipo alias para permitir enums com tipagem fraca, por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Adiciona testes do gerador de bindings para recursos do Go 1.24, por [@fbbdev](https://github.com/fbbdev) em [#4068](https://github.com/wailsapp/wails/pull/4068)
- Adiciona suporte ao macOS 15 "Sequoia" em `OSInfo.Branding` para melhorar a detecção da versão do sistema operacional em [#4065](https://github.com/wailsapp/wails/pull/4065)
- Adiciona o hook `PostShutdown` para executar código personalizado após a conclusão do processo de encerramento, por [@fbbdev](https://github.com/fbbdev) em [#4066](https://github.com/wailsapp/wails/pull/4066)
- Adiciona a struct `FatalError` para permitir a detecção de erros fatais em manipuladores de erros personalizados, por [@fbbdev](https://github.com/fbbdev) em [#4066](https://github.com/wailsapp/wails/pull/4066)
- Padroniza e documenta a ordem de inicialização e encerramento dos serviços, por [@fbbdev](https://github.com/fbbdev) em [#4066](https://github.com/wailsapp/wails/pull/4066)
- Adiciona uma infraestrutura de testes para a sequência de inicialização e encerramento da aplicação e para testes de inicialização e encerramento dos serviços, por [@fbbdev](https://github.com/fbbdev) em [#4066](https://github.com/wailsapp/wails/pull/4066)
- Adiciona o método `RegisterService` para registrar serviços depois que a aplicação tiver sido criada, por [@fbbdev](https://github.com/fbbdev) em [#4066](https://github.com/wailsapp/wails/pull/4066)
- Adiciona o campo `MarshalError` às opções da aplicação e dos serviços para o tratamento personalizado de erros em chamadas de bindings, por [@fbbdev](https://github.com/fbbdev) em [#4066](https://github.com/wailsapp/wails/pull/4066)
- Adiciona um wrapper de promise cancelável que propaga solicitações de cancelamento por cadeias de promises, por [@fbbdev](https://github.com/fbbdev) em [#4100](https://github.com/wailsapp/wails/pull/4100)
- Adiciona a possibilidade de vincular o cancelamento de chamadas de bindings a um `AbortSignal`, por [@fbbdev](https://github.com/fbbdev) em [#4100](https://github.com/wailsapp/wails/pull/4100)
- Adiciona suporte a atributos `data-wml-*` para WML, além dos atributos `wml-*` habituais, por [@leaanthony](https://github.com/leaanthony)
- Adiciona o método `Configure` a todos os serviços para configuração tardia ou reconfiguração dinâmica, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Quando não está configurado, o serviço `fileserver` envia uma resposta 503 Service Unavailable, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Quando não está configurado, o serviço `kvstore` fornece por padrão um armazenamento de chave-valor em memória, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Adiciona o método `Load` ao serviço `kvstore` para recarregar os dados do arquivo após alterações na configuração, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Adiciona o método `Clear` ao serviço `kvstore` para excluir todas as chaves, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Adiciona o tipo `Level` ao serviço `log` para fornecer constantes de nível de log no lado do JS, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Adiciona o método `Log` ao serviço `log` para definir dinamicamente o nível de log, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- O serviço `sqlite` fornece, por padrão, um banco de dados em memória quando não está configurado, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Adiciona o método `Close` ao serviço `sqlite` para fechar o banco de dados manualmente, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Adiciona suporte a cancelamento nos métodos de consulta do serviço `sqlite`, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Adiciona suporte a instruções preparadas ao serviço `sqlite`, com bindings JS, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Suporte ao Gin por [Lea Anthony](https://github.com/leaanthony) no [PR](https://github.com/wailsapp/wails/pull/3537), com base no trabalho original de [@AnalogJ](https://github.com/AnalogJ) neste [PR](https://github.com/wailsapp/wails/pull/3537)
- Corrige o salvamento automático e o salvamento automático de senhas, que estavam sempre habilitados, por [@oSethoum](https://github.com/osethoum) em [#4134](https://github.com/wailsapp/wails/pull/4134)
- Adiciona `SetMenu()` à janela para permitir a definição de um menu nela, por [@leaanthony](https://github.com/leaanthony)
- Adiciona suporte a notificações, por [@popaprozac](https://github.com/popaprozac) em [#4098](https://github.com/wailsapp/wails/pull/4098)
-  Adiciona suporte a associações de arquivos no macOS, por [@wimaha](https://github.com/wimaha) em [#4177](https://github.com/wailsapp/wails/pull/4177)
- Adiciona `wails3 tool version` para incrementar versões semânticas, por [@leaanthony](https://github.com/leaanthony)
- Adiciona suporte a emblemas no macOS e no Windows, por [@popaprozac](https://github.com/popaprozac) em [#](https://github.com/wailsapp/wails/pull/4234)
- Adiciona suporte a eventos registrados e com tipagem estrita, por [@fbbdev](https://github.com/fbbdev) e [@IanVS](https://github.com/IanVS) em [#4161](https://github.com/wailsapp/wails/pull/4161)
- Adiciona a capacidade de registrar hooks para eventos personalizados, por [@fbbdev](https://github.com/fbbdev) e [@IanVS](https://github.com/IanVS) em [#4161](https://github.com/wailsapp/wails/pull/4161)
- `app.OpenFileManager(path string, selectFile bool)` para abrir o gerenciador de arquivos do sistema no caminho `path`, com destaque opcional por meio de `selectFile`, por [@Krzysztofz01](https://github.com/Krzysztofz01) [@rcalixte](https://github.com/rcalixte)
- Nova flag `-git` para o comando `wails3 init`, por [@leaanthony](https://github.com/leaanthony)
- Novo comando `wails3 generate webview2bootstrapper`, por [@leaanthony](https://github.com/leaanthony)
- Adiciona o método `init()` ao runtime para permitir sua inicialização manual, por [@leaanthony](https://github.com/leaanthony)
- Adiciona a opção `WindowDidMoveDebounceMS` a WindowOptions da janela, por [@leaanthony](https://github.com/leaanthony)
- Adiciona o recurso de instância única, por [@leaanthony](https://github.com/leaanthony). Baseado no PR [v2](https://github.com/wailsapp/wails/pull/2951) de @APshenkin.
- Comando `wails3 generate template`, por [@leaanthony](https://github.com/leaanthony)
- Comando `wails3 releasenotes`, por [@leaanthony](https://github.com/leaanthony)
- Comando `wails3 update cli`, por [@leaanthony](https://github.com/leaanthony)
- Opção `-clean` para o comando `wails3 generate bindings`, por [@leaanthony](https://github.com/leaanthony)
- Permite compilações AppImage para Linux em aarch64 (arm64), por [@AkshayKalose](https://github.com/AkshayKalose) em [#3981](https://github.com/wailsapp/wails/pull/3981)
- Adiciona um hiperlink para o patrocinador, por @ansxuman em [#3958](https://github.com/wailsapp/wails/pull/3958)
- Suporte ao empacotamento para Linux nos formatos deb e rpm e a compilações de pacotes para Arch Linux, por
- Adiciona suporte a compilações e pacotes universais para Darwin, por
- Documentação sobre eventos no site, por
- Templates para sveltekit e sveltekit-ts configurados para desenvolvimento sem SSR
- Atualiza os assets de compilação usando o novo comando `wails3 update build-assets`, por
- Exemplo para testar a API HTML Drag and Drop, por
- Suporte a associações de arquivos, por [leaanthony](https://github.com/leaanthony) em
- Novo comando `wails3 generate runtime`, por
- Nova opção `InitialPosition` para especificar se a janela deve ser centralizada ou
- Adiciona os métodos `Path` e `Paths` ao pacote `application`, por
- Adiciona as opções `GeneralAutofillEnabled` e `PasswordAutosaveEnabled` do Windows
- Adiciona a capacidade de obter a janela que chamou um método de serviço, por
- Adiciona as opções `EnabledFeatures` e `DisabledFeatures` para WebView2, por
- ⊞ Novo sistema DIP para aprimorar o suporte a monitores com DPI alto, por
- ⊞ Opção de nome da classe da janela, por [windom](https://github.com/windom/) em
- Os serviços foram expandidos para fornecer funcionalidade de plugins. Por
- 🐧 Eventos WindowDidMove / WindowDidResize em
- ⊞ Evento WindowDidResize em
-  Adiciona o evento ApplicationShouldHandleReopen para permitir o tratamento do Dock
-  Adiciona getPrimaryScreen/getScreens à implementação, por @tmclane em
-  Adiciona uma opção para exibir a barra de ferramentas no modo de tela cheia no macOS, por
- 🐧 Adiciona a lógica onKeyPress para converter o pressionamento de teclas no Linux em um acelerador
- 🐧 Adiciona a tarefa `run:linux`, por
- Exporta o método `SetIcon`, por [@almas-x](https://github.com/almas-x) em
- Aprimora `OnShutdown`, por [@almas-x](https://github.com/almas-x) em
- Restaura o método `ToggleMaximise` na interface `Window`, por
- Adiciona mais informações a `Environment()`. Por @leaanthony em
- Expõe o método `WebviewWindow.IsFocused` na interface `Window`, por
- Oferece suporte a vários eventos de acionamento separados por espaços no sistema WML, por
- Adiciona exportações ESM ao script agrupado do runtime JS, por
- Adiciona uma flag ao gerador de bindings para usar o script agrupado do runtime JS em vez de
- Implementa `setIcon` no Linux, por [@abichinger](https://github.com/abichinger)
- Adiciona a flag `-port` ao comando dev e oferece suporte à variável de ambiente
- Adiciona testes para chamadas de métodos vinculados, por
- ⊞ Adiciona `SetIgnoreMouseEvents` para uma janela já criada, por
-  Adiciona a capacidade de definir o nível de empilhamento (ordem) de uma janela, por

### Corrigido

- Corrige `Screen.Bounds`, `WorkArea` e `Size` reduzidos pela metade em Macs com tela Retina, convertendo os valores em pontos de NSScreen em pixels do dispositivo nos campos `Physical*`, e preenche `Screen.X`/`Y` no nível superior para que a detecção de contato entre vários monitores e o posicionamento na área de trabalho estejam corretos no [PR](https://github.com/wailsapp/wails/pull/5168), por @wayneforrest
- Corrige uma condição de corrida de dados no ScreenManager que causa um deadlock no DisplayLink do WebKit quando a configuração de exibição muda (por exemplo, ao conectar um monitor externo durante a suspensão/retomada)
- Define diretamente CFBundleIconName como appicon quando Assets.car existe no [PR](https://github.com/wailsapp/wails/pull/5154), por @symball
- Corrige o `wails3 doctor` que informava pacotes WebKitGTK incorretos no Fedora, openSUSE, Arch e NixOS — as entradas alternativas de 4.0 foram removidas, pois v3 requer a API 4.1 em tempo de compilação (#5071)
- Corrige o nome do pacote webkit2gtk do doctor para openSUSE (`webkit2gtk4_1-devel` → `webkit2gtk3-devel`, o nome correto do pacote no openSUSE) (#5071)
- Corrige o erro de `Unexpected token '<'` quando `/wails/custom.js` está ausente no modo de desenvolvimento para desktop. Foi adicionado um manipulador explícito de 404 para `/wails/custom.js` e uma validação de `Content-Type` que não diferencia maiúsculas de minúsculas em `loadOptionalScript`, evitando que fallbacks de SPA em HTML sejam injetados como JavaScript. ([#5068](https://github.com/wailsapp/wails/issues/5068))
- Corrige o estado de destaque do menu da bandeja do sistema no macOS — agora o ícone exibe o estado selecionado quando o menu está aberto (#4910)
- Corrige a janela anexada à bandeja do sistema que aparecia atrás de outras janelas no macOS — agora usa o nível de janela pop-up apropriado (#4910)
- Corrige exemplos incorretos de importação de `@wailsio/runtime` em toda a documentação (#4989)
- Corrige a impossibilidade de minimizar janelas sem moldura no darwin (#4294)
- Corrige travamentos de 20-30 minuto durante `wails3 build` e `wails3 dev`, excluindo `node_modules/` da verificação de atualização do go-task. Anteriormente, o glob `sources: "**/*"` fazia o go-task enumerar e calcular a soma de verificação de todos os arquivos em `node_modules/` (mais de 50000-100000 arquivos com dependências pesadas, como MUI), o que era especialmente lento no Windows/NTFS (#4939)
- Corrige a falha de compilação do GTK4 causada pela colisão do typedef `Screen` de C com o X11 Xlib.h (#4957)
- Corrige a consistência dos métodos de selo do Dock no macOS
- Corrige a aplicação de `InvisibleTitleBarHeight` a todas as janelas do macOS, em vez de apenas às janelas sem moldura ou com barra de título transparente (#4960)
- Corrige a trepidação da janela ao redimensioná-la pelos cantos superiores com `InvisibleTitleBarHeight` habilitado, ignorando o início do arraste próximo às bordas da janela (#4960)
- Corrige a geração de tipos mapeados com chaves de enumeração nos bindings JS/TS (#4437), por @fbbdev
- Corrige o recurso de arrastar e soltar arquivos no Windows, que não funcionava com escalas de exibição diferentes de 100%
- Corrige o recurso interno de arrastar e soltar do HTML5, que deixava de funcionar quando a soltura de arquivos estava habilitada no Windows
- Corrige as coordenadas de soltura de arquivos no Windows, que usavam o espaço de pixels incorreto (pixels físicos em vez de pixels CSS)
- Corrige o recurso de arrastar e soltar arquivos no Linux, que não funcionava de forma confiável com efeitos de passagem do cursor
- Corrige o recurso interno de arrastar e soltar do HTML5, que deixava de funcionar quando a soltura de arquivos estava habilitada no Linux
- Corrige a exibição e ocultação de janelas no Linux/GTK4, que às vezes restaurava a janela no estado minimizado, usando `gtk_window_present()` (#4957)
- Corrige a obtenção e definição da posição da janela no Linux/GTK4, que sempre retornava 0,0, adicionando suporte condicional ao X11 por meio de `XTranslateCoordinates`/`XMoveWindow` (#4957)
- Corrige a falta de aplicação do tamanho máximo da janela no Linux/GTK4, adicionando limitação de tamanho baseada em sinais para substituir o `gtk_window_set_geometry_hints` removido (#4957)
- Corrige o dimensionamento por DPI no Linux/GTK4, implementando o cálculo correto de PhysicalBounds e o suporte a escala fracionária por meio de `gdk_monitor_get_scale` (GTK 4.14+)
- Corrige a duplicação de itens de menu ao criar novas janelas no Linux/GTK4
- Corrige a geração de tipos mapeados com chaves de enumeração nos bindings JS/TS (#4437), por @fbbdev
- Corrige o recurso de arrastar e soltar arquivos no Windows, que não funcionava com escalas de exibição diferentes de 100%
- Corrige o recurso interno de arrastar e soltar do HTML5, que deixava de funcionar quando a soltura de arquivos estava habilitada no Windows
- Corrige as coordenadas de soltura de arquivos no Windows, que usavam o espaço de pixels incorreto (pixels físicos em vez de pixels CSS)
- Corrige o recurso de arrastar e soltar arquivos no Linux, que não funcionava de forma confiável com efeitos de passagem do cursor
- Corrige o recurso interno de arrastar e soltar do HTML5, que deixava de funcionar quando a soltura de arquivos estava habilitada no Linux
- Corrige o dimensionamento por DPI no Linux/GTK4, implementando o cálculo correto de PhysicalBounds e o suporte a escala fracionária por meio de `gdk_monitor_get_scale` (GTK 4.14+)
- Corrige a duplicação de itens de menu ao criar novas janelas no Linux/GTK4
- Corrige a geração de tipos mapeados com chaves de enumeração nos bindings JS/TS (#4437), por @fbbdev
- Corrige o problema de "janelas fantasmas" no macOS, causado pelo acesso às APIs do AppKit fora da thread principal em App.Window.Current() (#4947), por @wimaha
- Corrige o `<input type="file">` do HTML, que não funcionava no macOS, implementando WKUIDelegate runOpenPanelWithParameters (#4862)
- Corrige o recurso nativo de arrastar e soltar arquivos, que não funcionava ao usar o módulo npm `@wailsio/runtime` no macOS/Linux (#4953), por @leaanthony
- Corrige a geração de bindings para aliases de tipos entre pacotes (#4578), por @fbbdev
- Corrige a falha do OpenFileDialog no Linux causada por uma violação da segurança de threads do GTK (#3683), por @ddmoney420
- Corrige a falha SIGSEGV ao chamar `Focus()` em uma janela oculta ou destruída (#4890), por @ddmoney420
- Corrige um possível panic ao definir um ícone ou bitmap vazio no Linux (#4923), por @ddmoney420
- Corrige a falha do ErrorDialog quando chamado por um binding de serviço no macOS (#3631), por @leaanthony
- Faz com que os menus sejam exibidos no Windows em `v3\examples\dialogs`, por @ndianabasi
- Corrige uma condição de corrida que causava TypeError durante o recarregamento da página (#4872), por @ddmoney420
- Corrige a saída incorreta dos testes do gerador de bindings, removendo o estado global do método `Collector.IsVoidAlias()` (#4941), por @fbbdev
- Corrige o seletor de arquivos de `<input type="file">`, que não funcionava no macOS (#4862), por @leaanthony
- Corrige o uso de sistemas de coordenadas inconsistentes por `Position()` e `SetPosition()` no macOS, que causava desvio da posição da janela ao salvar e restaurar o estado (#4816), por @leaanthony
- Corrige o erro "Access is denied" de SetProcessDpiAwarenessContext quando o reconhecimento de DPI já está definido pelo manifesto do aplicativo (#4803)
- Atualiza a página da documentação sobre atalhos de teclado e corrige o tipo do parâmetro de callback de `KeyBinding.Add`, por @ndianabasi
- Corrige a documentação sobre a geração de bindings personalizados: é necessário usar `-d String` em vez de `-o String`
- Corrige o menu que não removia os itens filhos em `menu.Update()`
- Corrige referências desatualizadas à API Manager na documentação (os arquivos 31 foram atualizados para usar o novo padrão, como `app.Window.New()`, `app.Event.Emit()` etc.), por @leaanthony
- Corrige a falha no Linux causada por um panic em métodos Go vinculados ao JS, devido à substituição dos manipuladores de sinais pelo WebKit (#3965), por @leaanthony
- Corrige SaveFileDialog.SetFilename(), que não produzia efeito no Linux (#4841), por @samstanier
- Corrige as coordenadas de soltura exibidas como undefined no exemplo de arrastar e soltar
- Corrige a falha na criação do pacote de aplicativo do macOS quando APP_NAME contém espaços (problema de expansão de chaves)
- Corrige um panic por índice fora dos limites no Windows ao chamar métodos de serviço (reverte goccy/go-json)
- Corrige o recurso de arrastar e soltar arquivos no Windows, que não funcionava com escalas de exibição diferentes de 100%
- Corrige o recurso interno de arrastar e soltar do HTML5, que deixava de funcionar quando a soltura de arquivos estava habilitada no Windows
- Corrige as coordenadas de soltura de arquivos no Windows, que usavam o espaço de pixels incorreto (pixels físicos em vez de pixels CSS)
- Corrige o recurso de arrastar e soltar arquivos no Linux, que não funcionava de forma confiável com efeitos de foco
- Corrige a interrupção do recurso interno de arrastar e soltar do HTML5 quando a soltura de arquivos estava habilitada no Linux
- Atualiza todos os comandos nos arquivos Taskfile.yml de todos os sistemas operacionais para aceitar espaços em variáveis como `APP_NAME`, por @ndianabasi
- Corrige o erro de argumento de comando ao executar a tarefa 'build:universal:lipo:go' no Linux, por @wux1an
- Corrige o erro do Docker "símbolo indefinido: **<em>ubsan</em>handle_xxxxxxx" ao executar 'wails3 build GOOS=darwin GOARCH=arm64' no Linux, por @wux1an
- Consolida a documentação de protocolos personalizados e adiciona seções sobre Universal Links, por @leaanthony
- Corrige a falha do menu da bandeja do sistema no Windows ao clicar repetidamente no ícone, adicionando uma proteção contra chamadas simultâneas a TrackPopupMenuEx (#4151), por @leaanthony
- Impede que o aplicativo falhe ao chamar systray.Run() antes de app.Run(), por @leaanthony
- Corrige uma falha no macOS ao alternar a visibilidade da janela por meio de Hide()/Show() com ApplicationShouldTerminateAfterLastWindowClosed habilitado (#4389), por @leaanthony
- Corrige um vazamento de memória nos menus de contexto do macOS e do Windows ao abri-los repetidamente (#4012), por @leaanthony
- Corrige a falta de reutilização dos recursos nativos dos menus de contexto no macOS, que fazia um novo menu ser criado a cada exibição (#4012), por @leaanthony
- Corrige o clique no ícone do Dock do macOS, que não exibia janelas ocultas quando o aplicativo era iniciado com `Hidden: true` (#4583), por @leaanthony
- Corrige a caixa de diálogo de impressão do macOS, que não abria devido ao tipo incorreto do ponteiro da janela na chamada CGO (#4290), por @leaanthony
- Corrige uma falha do menu da janela no Wayland causada pelo acesso do appmenu-gtk-module a uma janela ainda não realizada (#4769), por @leaanthony
- Corrige uma falha do aplicativo GTK quando o nome do aplicativo contém caracteres inválidos (espaços, parênteses etc.), por @leaanthony
- Corrige o erro "memória insuficiente" ao inicializar o recurso de arrastar e soltar no Windows (#4701), por @overlordtm
- Corrige o explorador de arquivos que abria o diretório errado no Linux devido ao escape incorreto do URI (#4397), por @leaanthony
- Corrige a falha na compilação do AppImage em distribuições Linux modernas (Arch, Fedora 39+, Ubuntu 24.04+) detectando automaticamente as seções ELF `.relr.dyn` e desabilitando a remoção de símbolos (#4642), por @leaanthony
- Corrige `wails doctor`, que informava incorretamente que os pacotes do webkit estavam instalados em sistemas baseados em Fedora/DNF (#4457), por @leaanthony
- Corrige o comportamento padrão de `config.yml`, que executava `wails3 dev` com uma compilação de produção, por @mbaklor
- Corrige stubs de serviço do iOS que causavam falhas de compilação devido à importação de um pacote inexistente, por @leaanthony
- Corrige o registro estruturado nos métodos de depuração/informação, que causava erros de "nenhuma diretiva de formatação", por @leaanthony
- Remove instruções temporárias de impressão de depuração incluídas acidentalmente durante a mesclagem da plataforma móvel, por @leaanthony
- Corrige uma falha do WebKitGTK no Wayland com GPUs NVIDIA (erro de protocolo 71) desabilitando automaticamente o renderizador DMA-BUF, por @leaanthony
- Corrige o valor alfa ignorado em `application.WebviewWindowOptions.BackgroundColour` no Linux ([#4722](https://github.com/wailsapp/wails/pull/4722), @BradHacker)
- Corrige o ícone da bandeja do sistema no Windows, que não usava por padrão o ícone do aplicativo quando nenhum ícone personalizado era fornecido (#4704)
- Rastreia a propriedade de `HICON` para destruir somente os identificadores criados pelo usuário, evitando falhas durante a reciclagem do Explorer (#4653).
- Libera o listener de tema do sistema do Windows e os ícones retidos da bandeja durante a destruição, eliminando o vazamento de goroutines e contextos de dispositivo (#4653).
- Trunca as dicas de ferramenta da bandeja em 127 unidades UTF-16 para evitar corromper pares substitutos e glifos multibyte (#4653).
- Corrige a falha da tarefa de empacotamento do Windows (#4667)
- Corrige a variável appicon do AppImage para Linux no taskfile do Linux [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- Corrige o erro de compilação no Windows causado pela alteração da assinatura do go-webview2 v1.0.22 (#4513, #4645)
- Corrige a variável appicon do AppImage para Linux no taskfile do Linux [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- Corrige o intervalo de protocolos em desktop.tmpl no Linux, removendo `<.Info.Protocol>` de `<.Protocol>`, por @Tolfx em #4510
- Corrige o erro de redefinição na demonstração de liquid glass em [#4542](https://github.com/wailsapp/wails/pull/4542), por @Etesam913
- Corrige as atualizações do menu da bandeja do sistema no Linux em [#4604](https://github.com/wailsapp/wails/issues/4604), por [@JackDoan](https://github.com/JackDoan)
- Corrige a janela branca que aparecia no Windows ao criar uma janela oculta, por @leaanthony em [#4612](https://github.com/wailsapp/wails/pull/4612)
- Corrige o caminho de importação do pacote de notificações na documentação, por @rxliuli em [#4617](https://github.com/wailsapp/wails/pull/4617)
- Corrige o recurso de arrastar e soltar, que não funcionava ao usar o pacote npm @wailsio/runtime (#4489), por @leaanthony em #4616
- Windows: corrige o piscar da janela ao iniciar e a exibição incorreta de janelas ocultas em [PR](https://github.com/wailsapp/wails/pull/4600), por @leaanthony.
- Corrigidos problemas de maximização do tamanho da janela no Wayland (https://github.com/wailsapp/wails/issues/4429) por [@samstanier](https://github.com/samstanier)
- Corrigidos os problemas de maximização do tamanho da janela no Wayland (https://github.com/wailsapp/wails/issues/4429) por [@samstanier](https://github.com/samstanier)
- Corrige o erro de redefinição na demonstração de liquid glass em [#4542](https://github.com/wailsapp/wails/pull/4542), por @Etesam913
- Corrige um problema que podia causar uma falha no AssetServer no macOS em [#4576](https://github.com/wailsapp/wails/pull/4576), por @jghiloni
- Corrige um problema de compilação ao compilar com NextJs. Corrigido em [#4585](https://github.com/wailsapp/wails/pull/4585) por @rev42
- Corrige os pipelines da versão nightly em [#4597](https://github.com/wailsapp/wails/pull/4597), por @riadafridishibly
- Corrige o erro de redefinição na demonstração de liquid glass em [#4542](https://github.com/wailsapp/wails/pull/4542), por @Etesam913
- Corrige um problema que podia causar uma falha no AssetServer no macOS em [#4576](https://github.com/wailsapp/wails/pull/4576), por @jghiloni
- Corrige um problema de compilação ao compilar com NextJs. Corrigido em [#4585](https://github.com/wailsapp/wails/pull/4585) por @rev42
- Corrige os pipelines da versão nightly em [#4597](https://github.com/wailsapp/wails/pull/4597), por @riadafridishibly
- Corrige o erro de redefinição na demonstração de liquid glass em [#4542](https://github.com/wailsapp/wails/pull/4542), por @Etesam913
- Corrige SetBackgroundColour no Windows, por @PPTGamer em [PR](https://github.com/wailsapp/wails/pull/4492)
- Atualiza a documentação para refletir as alterações da refatoração da API Manager, por @yulesxoxo em [PR #4476](https://github.com/wailsapp/wails/pull/4476)
- Corrige a variável appicon do arquivo .desktop do Linux no taskfile do Linux [PR #4477](https://github.com/wailsapp/wails/pull/4477)
- Atualiza a documentação para refletir as alterações da refatoração da API Manager, por @yulesxoxo em [PR #4476](https://github.com/wailsapp/wails/pull/4476)
- Corrige um bug de desreferência de ponteiro nulo no Windows, relatado em [#4456](https://github.com/wailsapp/wails/issues/4456), por @leaanthony em [#4460](https://github.com/wailsapp/wails/pull/4460)
- Adiciona suporte a `allowsBackForwardNavigationGestures` no WKWebView do macOS para habilitar gestos de navegação com deslize de dois dedos (#1857)
- Corrige o problema em que onClick não funcionava para itens de menu inicialmente definidos como desabilitados, por @leaanthony em [PR nº 4469](https://github.com/wailsapp/wails/pull/4469). Agradecimentos a @IanVS pela investigação inicial.
- Corrige a falta de encerramento do servidor Vite quando a compilação falha (nº 4403)
- Corrigido o panic ao fechar ou cancelar um `SaveFileDialog` no Windows. Corrigido por @hkhere em [PR](https://github.com/wailsapp/wails/pull/4284)
- Corrigido o recurso de arrastar e soltar no nível do HTML no Windows por [@mbaklor](https://github.com/mbaklor) em [nº 4259](https://github.com/wailsapp/wails/pull/4259)
- Adiciona suporte a `allowsBackForwardNavigationGestures` no WKWebView do macOS para habilitar gestos de navegação com deslize de dois dedos (nº 1857)
- Corrige o problema em que onClick não funcionava para itens de menu inicialmente definidos como desabilitados, por @leaanthony em [PR nº 4469](https://github.com/wailsapp/wails/pull/4469). Agradecimentos a @IanVS pela investigação inicial.
- Corrige a falta de encerramento do servidor Vite quando a compilação falha (nº 4403)
- Corrigida a análise de notificações no Windows por @popaprozac em [PR](https://github.com/wailsapp/wails/pull/4450)
- Corrigido o comando doctor para verificar as dependências do SDK do Windows por [@kodumulo](https://github.com/kodumulo) em [nº 4390](https://github.com/wailsapp/wails/issues/4390)
- Corrigida a desreferenciação de ponteiro nil em processURLRequest no Mac por [@etesam913](https://github.com/etesam913) em [nº 4366](https://github.com/wailsapp/wails/pull/4366)
- Corrigido um bug no Linux que impedia diálogos filtrados por [@bh90210](https://github.com/bh90210) em [nº 4287](https://github.com/wailsapp/wails/pull/4287)
- Corrigidos problemas do menu Edit no Windows e no Linux por [@leaanthony](https://github.com/leaanthony) em [nº 3f78a3a](https://github.com/wailsapp/wails/commit/3f78a3a8ce7837e8b32242c8edbbed431c68c062)
- Atualizada a versão mínima do sistema nos arquivos .plist do macOS de 10.13.0 para 10.15.0 por [@AkshayKalose](https://github.com/AkshayKalose) em [nº 3981](https://github.com/wailsapp/wails/pull/3981)
- Problema de salto no ID da janela corrigido por [@leaanthony](https://github.com/leaanthony)
- Corrige o problema de menu nil ao chamar RegisterContextMenu por [@leaanthony](https://github.com/leaanthony)
- Corrigidos ciclos de dependência na saída do gerador de bindings por [@fbbdev](https://github.com/fbbdev) em [nº 4001](https://github.com/wailsapp/wails/pull/4001)
- Corrigidos erros de uso antes da definição na saída do gerador de bindings por [@fbbdev](https://github.com/fbbdev) em [nº 4001](https://github.com/wailsapp/wails/pull/4001)
- Repassa as flags de compilação ao gerador de bindings por [@fbbdev](https://github.com/fbbdev) em [nº 4023](https://github.com/wailsapp/wails/pull/4023)
- Altera os caminhos no Taskfile do Windows para usar barras normais e garantir seu funcionamento em plataformas que não sejam Windows, por [@leaanthony](https://github.com/leaanthony)
- Eventos JS do Mac + Mac agora corrigidos por [@leaanthony](https://github.com/leaanthony)
- Corrigido um deadlock de eventos no macOS por [@leaanthony](https://github.com/leaanthony)
- Corrigido um erro `Parameter incorrect` na inicialização de Window no Windows quando há HTML, mas não há JS, por [@leaanthony](https://github.com/leaanthony)
- Corrigido o tamanho do prefixo da resposta usado para detectar o tipo de conteúdo no servidor de assets por [@fbbdev](https://github.com/fbbdev) em [nº 4049](https://github.com/wailsapp/wails/pull/4049)
- Corrigido o tratamento de respostas que não são 404 no caminho do índice raiz do servidor de assets por [@fbbdev](https://github.com/fbbdev) em [nº 4049](https://github.com/wailsapp/wails/pull/4049)
- Corrigido o comportamento indefinido no gerador de bindings ao testar propriedades de tipos genéricos por [@fbbdev](https://github.com/fbbdev) em [nº 4045](https://github.com/wailsapp/wails/pull/4045)
- Corrigida a saída do gerador de bindings para modelos quando o tipo subjacente não tem as mesmas propriedades que o wrapper nomeado, por [@fbbdev](https://github.com/fbbdev) em [nº 4045](https://github.com/wailsapp/wails/pull/4045)
- Corrigida a saída do gerador de bindings para tipos de chave de mapas e para o pré-processamento por [@fbbdev](https://github.com/fbbdev) em [nº 4045](https://github.com/wailsapp/wails/pull/4045)
- Corrigida a saída do gerador de bindings para structs que implementam interfaces de serialização por [@fbbdev](https://github.com/fbbdev) em [nº 4045](https://github.com/wailsapp/wails/pull/4045)
- Corrigida a detecção de ciclos de tipos que envolvem tipos genéricos no gerador de bindings por [@fbbdev](https://github.com/fbbdev) em [nº 4045](https://github.com/wailsapp/wails/pull/4045)
- Corrigidas referências inválidas a modelos não exportados na saída do gerador de bindings por [@fbbdev](https://github.com/fbbdev) em [nº 4045](https://github.com/wailsapp/wails/pull/4045)
- Movido o código injetado para o final dos arquivos de serviço por [@fbbdev](https://github.com/fbbdev) em [nº 4045](https://github.com/wailsapp/wails/pull/4045)
- Corrigido o tratamento de erros das operações de fechamento de arquivos no gerador de bindings por [@fbbdev](https://github.com/fbbdev) em [nº 4045](https://github.com/wailsapp/wails/pull/4045)
- Suprimidos os avisos para serviços que definem métodos de ciclo de vida ou HTTP, mas nenhum outro método vinculado, por [@fbbdev](https://github.com/fbbdev) em [nº 4045](https://github.com/wailsapp/wails/pull/4045)
- Corrigida a falha dos templates que não usam React ao exibir o rodapé Hello World quando o esquema de cores claro do sistema está em uso, por [@marcus-crane](https://github.com/marcus-crane) em [nº 4056](https://github.com/wailsapp/wails/pull/4056)
- Corrigidos itens de menu ocultos no macOS por [@leaanthony](https://github.com/leaanthony)
- Corrigidos o tratamento e a formatação de erros nos processadores de mensagens por [@fbbdev](https://github.com/fbbdev) em [nº 4066](https://github.com/wailsapp/wails/pull/4066)
-  Corrigida a omissão do encerramento de serviços ao sair do aplicativo por [@fbbdev](https://github.com/fbbdev) em [nº 4066](https://github.com/wailsapp/wails/pull/4066)
-  Garante que as atualizações de menu ocorram na thread principal por [@leaanthony](https://github.com/leaanthony)
- O mecanismo de arrastar e redimensionar agora é mais robusto e corresponde melhor ao comportamento esperado da plataforma, por [@fbbdev](https://github.com/fbbdev) em [nº 4100](https://github.com/wailsapp/wails/pull/4100)
- Corrigido [nº 4097](https://github.com/wailsapp/wails/issues/4097): Webpack/Angular descartava o código de inicialização do runtime, por [@fbbdev](https://github.com/fbbdev) em [nº 4100](https://github.com/wailsapp/wails/pull/4100)
- Corrigidos itens de menu inicialmente ocultos por [@IanVS](https://github.com/IanVS) em [nº 4116](https://github.com/wailsapp/wails/pull/4116)
- Corrigida a falha de assetFileServer ao servir arquivos `.html` em uma solicitação sem extensão quando `[request]` não existe, mas `[request].html` existe
- Corrigidos os caminhos de geração de ícones por [@robin-samuel](https://github.com/robin-samuel) em [nº 4125](https://github.com/wailsapp/wails/pull/4125)
- Corrigida a ausência de emissão dos eventos `fullscreen`, `unfullscreen`, `unminimise` e `unmaximise` por [@oSethoum](https://github.com/osethoum) em [nº 4130](https://github.com/wailsapp/wails/pull/4130)
- Corrigido o erro do NSIS causado por um prefixo incorreto na versão padrão da configuração por [@robin-samuel](https://github.com/robin-samuel) em [nº 4126](https://github.com/wailsapp/wails/pull/4126)
- Corrigido o retorno de caminhos com caracteres de escape pela função Dialogs do runtime no Windows, por [TheGB0077](https://github.com/TheGB0077) em [nº 4188](https://github.com/wailsapp/wails/pull/4188)
- Corrigido o caminho de detecção do WebView2 em HKCU por [@leaanthony](https://github.com/leaanthony).
- Corrigido um problema de entrada no macOS por [@leaanthony](https://github.com/leaanthony).
- Corrigido o nome do arquivo da tarefa de geração de ícones do Windows por [@yulesxoxo](https://github.com/yulesxoxo) em [#4219](https://github.com/wailsapp/wails/pull/4219).
- Corrigido o problema de transparência em janelas sem moldura por [@leaanthony](https://github.com/leaanthony), com base no trabalho de @kron.
- Corrigidas as chamadas de foco quando a janela está desabilitada ou minimizada por [@leaanthony](https://github.com/leaanthony), com base no trabalho de @kron.
- Corrigido o problema que impedia as bandejas do sistema de reaparecer após a reinicialização da barra de tarefas por [@leaanthony](https://github.com/leaanthony), com base no trabalho de @kron.
- Corrigido fallbackResponseWriter, que não implementava Flush(), em [#4245](https://github.com/wailsapp/wails/pull/4245)
- Corrigido fallbackResponseWriter, que não implementava Flush(), por [@superDingda] em [#4236](https://github.com/wailsapp/wails/issues/4236)
- Corrigidas as falhas ao fechar uma janela no macOS com uma chamada assíncrona pendente de uma função vinculada ao Go por [@joshhardy](https://github.com/joshhardy) em [#4354](https://github.com/wailsapp/wails/pull/4354)
- Corrigida uma condição de corrida na inicialização do modo de eficiência do Windows por [@leaanthony](https://github.com/leaanthony)
- Corrigida a limpeza de identificadores de ícone do Windows por [@leaanthony](https://github.com/leaanthony).
- Corrigido `OpenFileManager` no Windows por [@PPTGamer](https://github.com/PPTGamer) em [#4375](https://github.com/wailsapp/wails/pull/4375).
- Corrigidas as opções de largura mínima/máxima no Linux por @atterpac em [#3979](https://github.com/wailsapp/wails/pull/3979)
- Definições de tipos dos modelos TypeScript por meio da atualização da versão npm por @atterpac em [#3966](https://github.com/wailsapp/wails/pull/3966)
- Corrigida a referência CSS do modelo SvelteKit por @atterpac em [#3945](https://github.com/wailsapp/wails/pull/3945)
- Garantido que os principais callbacks em window run() sejam chamados na thread principal por [@leaanthony](https://github.com/leaanthony)
- Corrigidos os exemplos do seletor de diretórios da caixa de diálogo por [@leaanthony](https://github.com/leaanthony)
- Criada uma nova página de erro em chinês para quando index.html estiver ausente por [@leaanthony](https://github.com/leaanthony)
-  Garantido que o callback `windowDidBecomeKey` seja executado na thread principal por [@leaanthony](https://github.com/leaanthony)
-  Adicionado suporte ao modo de tela cheia em janelas sem moldura por [@leaanthony](https://github.com/leaanthony)
-  Aprimorada a lógica de destruição de janelas por [@leaanthony](https://github.com/leaanthony)
-  Corrigida a lógica de posicionamento de janelas quando anexadas às bandejas do sistema por [@leaanthony](https://github.com/leaanthony)
-  Adicionado suporte ao modo de tela cheia em janelas sem moldura por [@leaanthony](https://github.com/leaanthony)
- Corrigido o tratamento de eventos por [@leaanthony](https://github.com/leaanthony)
- Corrigida a lógica de encerramento de janelas por [@leaanthony](https://github.com/leaanthony)
- O taskfile comum agora gera, por padrão, bindings TypeScript para modelos TypeScript por [@leaanthony](https://github.com/leaanthony)
- Corrigido o fechamento do aplicativo ao receber a mensagem WM_CLOSE quando não há janelas abertas ou há somente a bandeja do sistema por [@mmalcek](https://github.com/mmalcek) em [#3990](https://github.com/wailsapp/wails/pull/3990)
- Corrigida a compilação com garble por @5aaee9 em [#3192](https://github.com/wailsapp/wails/pull/3192)
- Corrigidas as compilações NSIS do Windows por [@leaanthony](https://github.com/leaanthony)
- Corrigido o deadlock na caixa de diálogo do Linux para múltiplas seleções, causado por um elemento não fechado
- Corrigida a limpeza multiplataforma de arquivos .syso durante a compilação para Windows por
- Corrigida a compilação do AppImage para amd64 por @atterpac em
- Corrigida a atualização dos ativos de compilação por @ansxuman em
- Corrigida a implementação de `OnClick` e `OnRightClick` da bandeja do sistema no Linux por @atterpac
- Corrigido `AlwaysOnTop`, que não funcionava no Mac, por
-  Corrigido `application.NewEditMenu`, que incluía uma duplicata
- 🐧 Corrigida a compilação para aarch64
- ⊞ Corrigidos os itens de menu de grupos de botões de opção por
- Corrigido o erro ao compilar um .app executável no macOS quando 'name' e 'outputfilename'
- Corrigido o bug no uso de customEventProcessor no exemplo de arrastar e soltar por
- 🐧 Corrigido o erro de compilação no Linux introduzido pela adição de IgnoreMouseEvents por
- ⊞ Corrigido o bug na geração do arquivo de ícone syso por
- 🐧 Incorporida a correção para execução nativa no Wayland de
- Não vincular métodos de serviços internos em
- ⊞ Corrigido o pânico na inicialização da bandeja do sistema em
- Não vincular métodos de serviços internos em
- ⊞ Corrigido o pânico na inicialização da bandeja do sistema em
- Grande refatoração de itens de menu e do tratamento de eventos. Por enquanto, melhora principalmente o macOS. Por
- Corrigidos os testes após a refatoração de plugins e eventos em
- ⊞ Corrigido o aviso `Failed to unregister class Chrome_WidgetWin_0`. Por
- Problemas de módulos
- Corrigido o envio de mensagens de eventos de redimensionamento por [atterpac](https://github.com/atterpac) em
- 🐧 Corrigido o erro no tratamento de temas no NixOS por
- Corrigida a instalação de projetos entre volumes no Windows por
- Corrigido o CSS do modelo React para exibir o rodapé por
- Corrigidos os processos zumbis ao trabalhar no modo de desenvolvimento por meio da atualização para a versão mais recente de refresh
- Corrigida a obtenção do arquivo WebKit do AppImage por [Atterpac](https://github.com/atterpac)
- Corrigida a verificação de pacotes apt do Doctor por [Atterpac](https://github.com/Atterpac) em
- Corrigido o congelamento do aplicativo ao sair (Darwin) por @5aaee9 em
- Corrigidas as cores de fundo dos exemplos no Windows por
- Corrigidos os menus de contexto padrão por [mmghv](https://github.com/mmghv) em
- Corrigidos os valores hexadecimais das teclas de seta no Darwin por
- Implementado o funcionamento de arrastar e soltar no Windows. Adicionado por
- Corrigido um bug no Doctor no Linux quando o usuário não tem os drivers adequados
- Corrigido o dimensionamento de DPI na inicialização (Windows). Alterado por [@almas-x](https://github.com/almas-x) em
- Corrigida a linha de substituição em `go.mod` para usar caminhos relativos. Corrige caminhos do Windows com
- Corrigido o tratamento de cliques na bandeja do sistema no macOS quando não há uma janela anexada por
- Corrigida a falha de compilação no Windows causada por uma opção desconhecida por
- Corrigida a falha no Windows ao clicar com o botão esquerdo no ícone da bandeja do sistema quando não há um
- Corrigida a baseURL incorreta ao abrir a janela duas vezes por @5aaee9 no PR
- Corrigida a ordem das ramificações if no método `WebviewWindow.Restore` por
- Calculado corretamente `startURL` em várias invocações de `GetStartURL` quando
- Corrigido o tipo JS da struct `Screen` para corresponder à sua equivalente em Go por
- Corrigido o método `WML.Reload` para garantir a limpeza adequada dos eventos registrados
- Corrigido o fechamento imediato do menu de contexto personalizado no Linux por
- Corrigir o caminho de saída e a extensão dos arquivos de modelo produzidos pelo binding
- Corrigir os caminhos de importação dos arquivos de modelo no código JS produzido pelo binding
- Corrigir o recurso de arrastar e soltar em algumas distribuições Linux por
- Corrigir a tarefa ausente no macOS ao usar `wails3 task dev` por
- Corrigir a atribuição a um mapa nil causada pelo registro de eventos por
- Corrigir a desserialização dos parâmetros de métodos vinculados por
- Corrigir o tratamento de vários valores de retorno de métodos vinculados por
- Corrigir a detecção, pelo doctor, do npm que não foi instalado com o gerenciador de pacotes do sistema
- Corrigir a ausência de MicrosoftEdgeWebview2Setup.exe. Agradecimentos a
- Corrigir uma falha aleatória no Linux causada pelo tratamento do ID da janela por @leaanthony. Baseado em
- Corrigir a falha de systemTray.setIcon no Linux por
- Garantir que a moldura da janela seja aplicada na primeira chamada da função `setFrameless` em

### Alterado

- **INCOMPATÍVEL**: As chaves de mapas nos bindings JS/TS gerados agora são marcadas como opcionais para refletir com precisão a semântica dos mapas do Go. O acesso aos valores de mapas no TypeScript agora retorna `T | undefined` em vez de `T`, exigindo verificações de nulo ou asserções (#4943) por `@fbbdev`
- Alterar o uso de `Event` para `Events` de acordo com as alterações em `@wailsio/runtime` e atualizar as chamadas de função correspondentes na documentação em `Features/Events/Event System` por @AbdelhadiSeddar
- Mover `EnabledFeatures`, `DisabledFeatures` e `AdditionalBrowserArgs` das opções específicas de cada janela para `Options.Windows` no nível da aplicação (#4559) por @leaanthony
- Atualizar o README do exemplo `Drag N Drop` e destacar que o exemplo demonstra `Internal Drag and Drop` por @ndianabasi
- Alterar diversos logs de depuração de Info para Debug (por @mbaklor)
- **INCOMPATÍVEL:** Renomear `EnableDragAndDrop` como `EnableFileDrop` nas opções de janela
- **INCOMPATÍVEL:** Renomear `DropZoneDetails` como `DropTargetDetails` no contexto do evento
- **INCOMPATÍVEL:** Renomear o método `DropZoneDetails()` como `DropTargetDetails()` em `WindowEventContext`
- **INCOMPATÍVEL:** Remover o evento `WindowDropZoneFilesDropped`; usar `WindowFilesDropped` em seu lugar
- **INCOMPATÍVEL:** Alterar o atributo HTML de `data-wails-dropzone` para `data-file-drop-target`
- **INCOMPATÍVEL:** Alterar a classe CSS de hover de `wails-dropzone-hover` para `file-drop-target-active`
- **INCOMPATÍVEL:** Remover as opções `DragEffect`, `OnEnterEffect` e `OnOverEffect` do Windows (faziam parte do IDropTarget removido)
- Adotar goccy/go-json para todo o processamento de JSON em tempo de execução (bindings de métodos, eventos, solicitações da webview, notificações e kvstore), melhorando o desempenho em 21-63% e reduzindo as alocações de memória em 40-60%
- Otimizar o layout da struct BoundMethod e armazenar em cache o sinalizador isVariadic para reduzir a sobrecarga por chamada
- Usar um buffer de argumentos alocado na pilha para métodos com `<=8` argumentos, evitando alocações no heap
- Otimizar a coleta de resultados em chamadas de métodos para evitar a alocação de slice quando houver um único valor de retorno
- Usar sync.Map no cache de tipos MIME para melhorar o desempenho em operações concorrentes
- Usar um pool de buffers para ler o corpo das solicitações do transporte HTTP
- Alocar o canal CloseNotify sob demanda no detector de tipo de conteúdo para reduzir as alocações por solicitação
- Remover o registro de depuração de CSS do servidor de assets
- Expandir o mapa de extensões de tipos MIME para abranger mais de 50 formatos comuns da Web (fontes, áudio, vídeo etc.)
- Atualizar a documentação das opções `X/Y` de Window por @ruhuang2001
- Atualizar a documentação de `Frontend Runtime`, adicionando mais opções para gerar bindings de frontend, por @ndianabasi
- Atualizar a página de documentação do servidor de assets v3 do Wails por @ndianabasi
- **INCOMPATÍVEL**: Remover as funções de diálogo no nível do pacote (`application.InfoDialog()`, `application.QuestionDialog()` etc.). Usar o gerenciador `app.Dialog` em seu lugar: `app.Dialog.Info()`, `app.Dialog.Question()`, `app.Dialog.Warning()`, `app.Dialog.Error()`, `app.Dialog.OpenFile()`, `app.Dialog.SaveFile()`
- Atualizar a documentação de diálogos para corresponder à API real: usar `app.Dialog.*` e `AddButton()` com callbacks (não `SetButtons()`), `SetDefaultButton(*Button)` (não uma string), `AddFilter()` (não `SetFilters()`), `SetFilename()` (não `SetDefaultFilename()`) e `app.Dialog.OpenFile().CanChooseDirectories(true)` para selecionar pastas
- **INCOMPATÍVEL**: Os builds de produção agora são o padrão. Para criar builds de desenvolvimento, defina `DEV=true` nos seus Taskfiles. Gerar um novo projeto para os exemplos por @leaanthony
- Ao emitir um evento personalizado com zero ou um argumento de dados, o valor dos dados será atribuído diretamente ao campo Data, sem ser encapsulado em um slice, por [@fbbdev](https://github.com/fbbdev) em [#4633](https://github.com/wailsapp/wails/pull/4633)
- As bandejas do sistema no Windows agora respeitam `SystemTray.Show()`/`Hide()` alternando `NIS_HIDDEN`, permitindo que os aplicativos realmente desapareçam e retornem (#4653).
- O registro da bandeja do sistema reutiliza ícones já resolvidos, define `NOTIFYICON_VERSION_4` uma única vez e habilita `NIF_SHOWTIP` para que as dicas de ferramenta sejam restauradas após a reinicialização do Explorer (#4653).
- macOS: Usar `visibleFrame` em vez de `frame` para centralizar a janela excluindo as áreas da barra de menus e do Dock
- macOS: Usar `visibleFrame` em vez de `frame` para centralizar a janela excluindo as áreas da barra de menus e do Dock
- Ao executar `wails3 update build-assets` com o parâmetro `-config`, os valores definidos por meio dos parâmetros `-product*` são
- `window.NativeWindowHandle()` -> `window.NativeWindow()` por @leaanthony em [#4471](https://github.com/wailsapp/wails/pull/4471)
- Refatorar o tratamento interno de janelas por @leaanthony em [#4471](https://github.com/wailsapp/wails/pull/4471)
- Remover `application.WindowIDKey` e `application.WindowNameKey` (substituídos por `application.WindowKey`) por [@leaanthony](https://github.com/leaanthony)
- ContextMenuData agora retorna uma string em vez de any por [@leaanthony](https://github.com/leaanthony)
- Nos bindings JS/TS, os campos de classe de tipos de array de tamanho fixo agora são inicializados com o tamanho esperado, em vez de ficarem vazios, por [@fbbdev](https://github.com/fbbdev) em [#4001](https://github.com/wailsapp/wails/pull/4001)
- ContextMenuData agora retorna uma string em vez de any por [@leaanthony](https://github.com/leaanthony)
- `application.NewService` não aceita mais options como parâmetro opcional (use `application.NewServiceWithOptions` em seu lugar) por [@leaanthony](https://github.com/leaanthony) em [#4024](https://github.com/wailsapp/wails/pull/4024)
- Remover a dependência de `nanoid` por [@leaanthony](https://github.com/leaanthony)
- Atualizar o exemplo Window para os estilos de janela mica/acrylic/tabbed por [@leaanthony](https://github.com/leaanthony)
- Nos bindings JS/TS, os arquivos de modelo `internal.js/ts` foram removidos; agora todos os modelos podem ser encontrados em `models.js/ts` por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Nos bindings JS/TS, tipos nomeados nunca são renderizados como aliases de outros tipos nomeados; o comportamento anterior agora se restringe a aliases por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Nos bindings JS/TS, no modo de classe, os campos de structs cujo tipo é um parâmetro de tipo são marcados como opcionais e nunca são inicializados automaticamente, por [@fbbdev](https://github.com/fbbdev) em [#4045](https://github.com/wailsapp/wails/pull/4045)
- Remoção do ESLint dos templates, por [@IanVS](https://github.com/IanVS) em [#4059](https://github.com/wailsapp/wails/pull/4059)
- Atualização da data de copyright para 2025, por [@IanVS](https://github.com/IanVS) em [#4037](https://github.com/wailsapp/wails/pull/4037)
- Adição de documentação para event.Sender, por [@IanVS](https://github.com/IanVS) em [#4075](https://github.com/wailsapp/wails/pull/4075)
- Suporte ao Go 1.24, por [@leaanthony](https://github.com/leaanthony)
- Os hooks de `ServiceStartup` agora são invocados quando `App.Run` é chamado, e não em `application.New`, por [@fbbdev](https://github.com/fbbdev) em [#4066](https://github.com/wailsapp/wails/pull/4066)
- Os erros de `ServiceStartup` agora são retornados por `App.Run` em vez de encerrarem o processo, por [@fbbdev](https://github.com/fbbdev) em [#4066](https://github.com/wailsapp/wails/pull/4066)
- As chamadas de bindings e diálogos feitas pelo JS agora são rejeitadas com objetos de erro em vez de strings, por [@fbbdev](https://github.com/fbbdev) em [#4066](https://github.com/wailsapp/wails/pull/4066)
- Melhoria no posicionamento do menu da bandeja do sistema no Windows, por [@leaanthony](https://github.com/leaanthony)
- O runtime JS foi portado para TypeScript por [@fbbdev](https://github.com/fbbdev) em [#4100](https://github.com/wailsapp/wails/pull/4100)
- O runtime é inicializado assim que é importado, sem a necessidade de aguardar o carregamento da janela, por [@fbbdev](https://github.com/fbbdev) em [#4100](https://github.com/wailsapp/wails/pull/4100)
- O runtime não exporta mais um método init. Uma importação somente para efeitos colaterais pode ser usada para inicializá-lo, por [@fbbdev](https://github.com/fbbdev) em [#4100](https://github.com/wailsapp/wails/pull/4100)
- Os métodos vinculados agora retornam uma `CancellablePromise` que é rejeitada com um `CancelError` em caso de cancelamento. O resultado efetivo da chamada é descartado, por [@fbbdev](https://github.com/fbbdev) em [#4100](https://github.com/wailsapp/wails/pull/4100)
- Os tipos de serviços integrados agora são chamados consistentemente de `Service`, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- As funções de criação de serviços integrados com opções agora são chamadas consistentemente de `NewWithConfig`, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- O método `Select` do serviço `sqlite` agora se chama `Query`, para manter a consistência com as APIs Go, por [@fbbdev](https://github.com/fbbdev) em [#4067](https://github.com/wailsapp/wails/pull/4067)
- Templates: runtime movido para "dependencies" e arquivos package.json organizados, por [@IanVS](https://github.com/IanVS) em [#4133](https://github.com/wailsapp/wails/pull/4133)
- Criação e assinatura ad hoc de pacotes de aplicativos durante o desenvolvimento para habilitar determinadas APIs do macOS, por [@popaprozac](https://github.com/popaprozac) em [#4171](https://github.com/wailsapp/wails/pull/4171)
- Recursos de compilação movidos para diretórios específicos de cada plataforma, por [@leaanthony](https://github.com/leaanthony)
- Taskfiles movidos e renomeados para diretórios específicos de cada plataforma, por [@leaanthony](https://github.com/leaanthony)
- Experiência significativamente melhor quando `index.html` está ausente, por [@leaanthony](https://github.com/leaanthony)
- [Windows] Melhoria no desempenho ao minimizar e restaurar, por [@leaanthony](https://github.com/leaanthony). Com base no [PR](https://github.com/wailsapp/wails/pull/3955) original de [562589540](https://github.com/562589540)
- Remoção da opção `ShouldClose` (registre um hook para events.Common.WindowClosing em vez disso), por [@leaanthony](https://github.com/leaanthony)
- [Windows] Redução da oscilação visual ao abrir uma janela, por [@leaanthony](https://github.com/leaanthony)
- Remoção de `Window.Destroy`, pois deveria ser uma função interna, por [@leaanthony](https://github.com/leaanthony)
- Eventos `WindowClose` renomeados para `WindowClosing`, por [@leaanthony](https://github.com/leaanthony)
- As compilações do frontend agora usam o ambiente "development" ou "production" do vite, dependendo do tipo de compilação, por [@leaanthony](https://github.com/leaanthony)
- Atualização para go-webview2 v1.19, por [@leaanthony](https://github.com/leaanthony)
- Garantia de uso do fork do taskfile, por @leaanthony
- Atualização do fork do Taskfile para corrigir problemas de versão durante a instalação usando
- Uso do fork do Taskfile para corrigir problemas de versão durante a instalação usando
- `service.OnStartup` agora encerra o aplicativo em caso de erro e executa
- Refatoração das mensagens de clique na bandeja do sistema para melhor correspondência com as interações do usuário, por
- Incorporação de recursos passa a incluir `all:frontend/dist` para oferecer suporte a frameworks que geram
- Refatoração do Taskfile por [leaanthony](https://github.com/leaanthony) em
- Atualização para `go-webview2` v1.0.16 por
- Correção do tipo `Screen` para incluir `ID`, e não `Id`, por
- Atualização da versão do Wails em `go.mod.tmpl` para oferecer suporte a `application.ServiceOptions`, por
- Correção da determinação do nome do serviço por [windom](https://github.com/windom/) em
- mkdocs serve agora usa Docker, por [leaanthony](https://github.com/leaanthony)
- Configuração de desenvolvimento consolidada em `config.yml` por
- A caixa de diálogo da bandeja do sistema agora usa por padrão o ícone do aplicativo, se disponível (Windows), por
- Melhoria nos relatórios de GPU e memória no macOS, por
- Remoção de `WebviewGpuIsDisabled` e `EnableFraudulentWebsiteWarnings`
- Alteração na API de eventos: `On`/`Emit` -> eventos do usuário, `OnApplicationEvent` ->
- Correção da API de eventos no Linux por [TheGB0077](https://github.com/TheGB0077) em
- [CI] melhorias nas actions e habilitação da execução de actions também em forks e
- `AbsolutePosition()` renomeado para `Position()` por
- Atualização da dependência WebKit do Linux para webkit2gtk-4.1 em vez de webkitgtk2-4.0 para
- O script do runtime JS incluído no pacote agora é um módulo ESM: tags de script que o importam
- O pacote `@wailsio/runtime` não publica sua API no `window.wails`
- O módulo `@wailsio/runtime/src/window` da API de janelas agora expõe o contêiner
- A API de janelas do JS foi atualizada para corresponder à API Go `WebviewWindow` atual
- O gerador de bindings agora usa chamadas por ID como padrão. A opção `-id` da CLI
- Novo layout do código de bindings: anteriormente, os arquivos de saída eram organizados em pastas
- O campo `application.Options.Bind` da struct foi renomeado para
- Nova sintaxe para vincular serviços: agora, as instâncias de serviços devem ser encapsuladas em um
- Desativação do indicador giratório em ambientes sem terminal ou de CI, por

### Removido

- **ALTERAÇÃO INCOMPATÍVEL**: remove `EnabledFeatures`, `DisabledFeatures` e `AdditionalLaunchArgs` das opções `WindowsWindow` por janela. Em vez disso, use `Options.Windows.EnabledFeatures`, `Options.Windows.DisabledFeatures` e `Options.Windows.AdditionalBrowserArgs` no nível do aplicativo. Esses sinalizadores se aplicam globalmente ao ambiente WebView2 compartilhado (#4559), por @leaanthony
- Remove a implementação nativa de `IDropTarget` no Windows em favor da abordagem baseada em JavaScript (corresponde ao comportamento de v2)
- Remove a dependência github.com/wailsapp/mimetype em favor de um mapa de extensões ampliado e de http.DetectContentType da biblioteca padrão, reduzindo o tamanho do binário em cerca de 1.2 MB
- Remove a dependência gopkg.in/ini.v1 por meio da implementação de um analisador mínimo de arquivos .desktop para o explorador de arquivos do Linux, economizando cerca de 45 KB
- Remove samber/lo do código de runtime usando o pacote slices da biblioteca padrão do Go 1.21+ e auxiliares internos mínimos, economizando cerca de 310 KB
- Remove instruções printf de depuração do manipulador de esquemas de URL do Darwin (#4834)
- **ALTERAÇÃO INCOMPATÍVEL**: remove o evento `linux:WindowLoadChanged`; use `linux:WindowLoadFinished` para detectar quando o WebView terminar de carregar (#3896), por @leaanthony

### Alterações incompatíveis

- **Refatoração da API de gerenciadores**: reorganizou a API do aplicativo, substituindo a estrutura plana por gerenciadores organizados para melhorar a organização e facilitar a descoberta do código, por [@leaanthony](https://github.com/leaanthony) em [#4359](https://github.com/wailsapp/wails/pull/4359)
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
- Moveu os métodos `Path` e `Paths` para o pacote `application`, por [@leaanthony](https://github.com/leaanthony)
- O menu do aplicativo agora está disponível somente no macOS, por [@leaanthony](https://github.com/leaanthony)

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

+ Corrige dependências de pacotes Linux excessivamente abrangentes e dependências RPM desatualizadas.

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

variável de canal por @michael-freling em   [#3925](https://github.com/wailsapp/wails/pull/3925)   [ansxuman](https://github.com/ansxuman) em   [#3924](https://github.com/wailsapp/wails/pull/3924)   [#3898](https://github.com/wailsapp/wails/pull/3898)   [#3901](https://github.com/wailsapp/wails/pull/3901)   em [#3886](https://github.com/wailsapp/wails/pull/3886)   [leaanthony](https://github.com/leaanthony) em   [#3841](https://github.com/wailsapp/wails/pull/3841)   `PasteAndMatchStyle` função no menu de edição no Darwin por   [johnmccabe](https://github.com/johnmccabe) em   [#3839](https://github.com/wailsapp/wails/pull/3839)   [#3840](https://github.com/wailsapp/wails/issues/3840) em   [#3854](https://github.com/wailsapp/wails/pull/3854) por   [kodflow](https://github.com/kodflow)   [@leaanthony](https://github.com/leaanthony)   são diferentes. por @nickisworking em   [#3789](https://github.com/wailsapp/wails/pull/3789)

## v3.0.0-alpha.7 - 2024-09-18

### Adicionado

[mmghv](https://github.com/mmghv) em   [#3665](https://github.com/wailsapp/wails/pull/3665)   [#3682](https://github.com/wailsapp/wails/pull/3682)   [atterpac](https://github.com/atterpac) e   [leaanthony](https://github.com/leaanthony) em   [#3570](https://github.com/wailsapp/wails/pull/3570)

### Alterado

Eventos do aplicativo `OnWindowEvent` -> Eventos da janela, por   [leaanthony](https://github.com/leaanthony)   [#3734](https://github.com/wailsapp/wails/pull/3734)   branches prefixadas com `v3/` ou `v3-` por   [stendler](https://github.com/stendler) em   [#3747](https://github.com/wailsapp/wails/pull/3747)

### Corrigido

[etesam913](https://github.com/etesam913) em   [#3742](https://github.com/wailsapp/wails/pull/3742)   [atterpac](https://github.com/atterpac) em   [#3721](https://github.com/wailsapp/wails/pull/3721)   [atterpac](https://github.com/atterpac) em   [#3675](https://github.com/wailsapp/wails/pull/3675)   [#1811](https://github.com/wailsapp/wails/pull/1811) em   [#3614](https://github.com/wailsapp/wails/pull/3614) por   [@stendler](https://github.com/stendler)   [#3720](https://github.com/wailsapp/wails/pull/3720) por   [leaanthony](https://github.com/leaanthony)   [#3693](https://github.com/wailsapp/wails/issues/3693) por   [@DeltaLaboratory](https://github.com/DeltaLaboratory)   [#3720](https://github.com/wailsapp/wails/pull/3720) por   [leaanthony](https://github.com/leaanthony)   [#3693](https://github.com/wailsapp/wails/issues/3693) por   [@DeltaLaboratory](https://github.com/DeltaLaboratory)   [leaanthony](https://github.com/leaanthony)   [#3746](https://github.com/wailsapp/wails/pull/3746) por   [@stendler](https://github.com/stendler)   [leaanthony](https://github.com/leaanthony)

## v3.0.0-alpha.6 - 2024-07-30

### Corrigido

## v3.0.0-alpha.5 - 2024-07-30

### Adicionado

[#3580](https://github.com/wailsapp/wails/pull/3580)   [#3580](https://github.com/wailsapp/wails/pull/3580)   clique no ícone por @5aaee9 em [#2991](https://github.com/wailsapp/wails/pull/2991)   [#2618](https://github.com/wailsapp/wails/pull/2618)   [@fbbdev](https://github.com/fbbdev) em   [#3282](https://github.com/wailsapp/wails/pull/3282)   @[Atterpac](https://github.com/Atterpac)   em [#3022](https://github.com/wailsapp/wails/pull/3022])   [@marcus-crane](https://github.com/marcus-crane) em   [#3146](https://github.com/wailsapp/wails/pull/3146)   [PR](https://github.com/wailsapp/wails/pull/3147)   [PR](https://github.com/wailsapp/wails/pull/3189)   [@fbbdev](https://github.com/fbbdev) em   [#3281](https://github.com/wailsapp/wails/pull/3281)   [aba82cc](https://github.com/wailsapp/wails/commit/aba82cc52787c97fb99afa58b8b63a0004b7ff6c)   com base no [PR](https://github.com/wailsapp/wails/pull/2044) de @Mai-Lapyst   [@fbbdev](https://github.com/fbbdev) em   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@fbbdev](https://github.com/fbbdev) em   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@fbbdev](https://github.com/fbbdev) em   [#3295](https://github.com/wailsapp/wails/pull/3295)   o pacote npm por [@fbbdev](https://github.com/fbbdev) em   [#3334](https://github.com/wailsapp/wails/pull/3334)   em [#3354](https://github.com/wailsapp/wails/pull/3354)   `WAILS_VITE_PORT` por [@abichinger](https://github.com/abichinger) em   [#3429](https://github.com/wailsapp/wails/pull/3429)   [@abichinger](https://github.com/abichinger) em   [#3431](https://github.com/wailsapp/wails/pull/3431)   [@bruxaodev](https://github.com/bruxaodev) em   [#3667](https://github.com/wailsapp/wails/pull/3667)   [@OlegGulevskyy](https://github.com/OlegGulevskyy) em   [#3674](https://github.com/wailsapp/wails/pull/3674)

### Corrigido

[#3606](https://github.com/wailsapp/wails/pull/3606)   [tmclane](https://github.com/tmclane) em   [#3515](https://github.com/wailsapp/wails/pull/3515)   [atterpac](https://github.com/atterac) em   [#3512](https://github.com/wailsapp/wails/pull/3512)   [atterpac](https://github.com/atterpac) em   [#3477](https://github.com/wailsapp/wails/pull/3477)   por [Atterpac](https://github.com/atterpac) em   [#3320](https://github.com/wailsapp/wails/pull/3320).   em [#3306](https://github.com/wailsapp/wails/pull/3306).   [#2972](https://github.com/wailsapp/wails/pull/2972).   [#2982](https://github.com/wailsapp/wails/pull/2982)   [mmghv](https://github.com/mmghv) em   [#2750](https://github.com/wailsapp/wails/pull/2750).   [#2753](https://github.com/wailsapp/wails/pull/2753).   [jaybeecave](https://github.com/jaybeecave) em   [#3052](https://github.com/wailsapp/wails/pull/3052).   [@pylotlight](https://github.com/pylotlight) no [PR](https://github.com/wailsapp/wails/pull/3039)   instalado. Adicionado por [@pylotlight](https://github.com/pylotlight) no   [PR](https://github.com/wailsapp/wails/pull/3032)   [PR](https://github.com/wailsapp/wails/pull/3145)   espaços — @leaanthony.   [thomas-senechal](https://github.com/thomas-senechal) no PR   [#3207](https://github.com/wailsapp/wails/pull/3207)   [thomas-senechal](https://github.com/thomas-senechal) no PR   [#3208](https://github.com/wailsapp/wails/pull/3208)   janela anexada [tw1nk](https://github.com/tw1nk) no PR   [#3271](https://github.com/wailsapp/wails/pull/3271)   [#3273](https://github.com/wailsapp/wails/pull/3273)   [@fbbdev](https://github.com/fbbdev) em   [#3279](https://github.com/wailsapp/wails/pull/3279)   `FRONTEND_DEVSERVER_URL` está presente.   [#3299](https://github.com/wailsapp/wails/pull/3299)   [@fbbdev](https://github.com/fbbdev) em   [#3295](https://github.com/wailsapp/wails/pull/3295)   listeners por [@fbbdev](https://github.com/fbbdev) em   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@abichinger](https://github.com/abichinger) em   [#3330](https://github.com/wailsapp/wails/pull/3330)   gerador por [@fbbdev](https://github.com/fbbdev) em   [#3334](https://github.com/wailsapp/wails/pull/3334)   gerador por [@fbbdev](https://github.com/fbbdev) em   [#3334](https://github.com/wailsapp/wails/pull/3334)   [@abichinger](https://github.com/abichinger) em   [#3346](https://github.com/wailsapp/wails/pull/3346)   [@hfoxy](https://github.com/hfoxy) em   [#3417](https://github.com/wailsapp/wails/pull/3417)   [@hfoxy](https://github.com/hfoxy) em   [#3426](https://github.com/wailsapp/wails/pull/3426)   [@fbbdev](https://github.com/fbbdev) em   [#3431](https://github.com/wailsapp/wails/pull/3431)   [@fbbdev](https://github.com/fbbdev) em   [#3431](https://github.com/wailsapp/wails/pull/3431)   por [@pekim](https://github.com/pekim) em   [#3458](https://github.com/wailsapp/wails/pull/3458)   [@robin-samuel](https://github.com/robin-samuel).   PR [#3466](https://github.com/wailsapp/wails/pull/3622) por   [@5aaee9](https://github.com/5aaee9).   [@windom](https://github.com/windom/) em   [#3636](https://github.com/wailsapp/wails/pull/3636).   Windows por [@bruxaodev](https://github.com/bruxaodev/) em   [#3691](https://github.com/wailsapp/wails/pull/3691).

### Alterado

[mmghv](https://github.com/mmghv) em   [#3611](https://github.com/wailsapp/wails/pull/3611)   oferece suporte ao Ubuntu 24.04 LTS, por [atterpac](https://github.com/atterpac) em   [#3461](https://github.com/wailsapp/wails/pull/3461)   deve ter o atributo `type="module"`. Por   [@fbbdev](https://github.com/fbbdev) em   [#3295](https://github.com/wailsapp/wails/pull/3295)   objeto e não inicia o sistema WML. Isso foi feito para melhorar o encapsulamento. Se desejar, você pode iniciar o sistema WML manualmente chamando o novo método `WML.Enable`. O script de runtime JS incluído ainda executa ambas as operações automaticamente. Por [@fbbdev](https://github.com/fbbdev) em   [#3295](https://github.com/wailsapp/wails/pull/3295)   objeto de janela como exportação padrão. Não é mais possível importar métodos individuais por meio da sintaxe de importação nomeada ou de namespace do ESM.   API. Alguns métodos mudaram de nome ou de protótipo, especificamente: `Screen`   passa a ser `GetScreen`; `GetZoomLevel`/`SetZoomLevel` passam a ser `GetZoom`/`SetZoom`;   `GetZoom`, `Width` e `Height` agora retornam valores diretamente, em vez de encapsulá-los em objetos. Por [@fbbdev](https://github.com/fbbdev) em   [#3295](https://github.com/wailsapp/wails/pull/3295)   foi removido. Use a opção de CLI `-names` para voltar a usar chamadas por nome.   Por [@fbbdev](https://github.com/fbbdev) em   [#3468](https://github.com/wailsapp/wails/pull/3468)   nomeados de acordo com o pacote que os contém; agora são usados os caminhos completos de importação do Go, incluindo o caminho do módulo. Por [@fbbdev](https://github.com/fbbdev) em   [#3468](https://github.com/wailsapp/wails/pull/3468)   `application.Options.Services`. Por [@fbbdev](https://github.com/fbbdev) em   [#3468](https://github.com/wailsapp/wails/pull/3468)   chamada a `application.NewService`. Por [@fbbdev](https://github.com/fbbdev) em   [#3468](https://github.com/wailsapp/wails/pull/3468)   [@DeltaLaboratory](https://github.com/DeltaLaboratory) em   [#3574](https://github.com/wailsapp/wails/pull/3574)
