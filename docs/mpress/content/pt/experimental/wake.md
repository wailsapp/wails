---
title: "Wake"
description: "Crie seu primeiro projeto Wake, siga o assistente de criação de projetos, compile e execute o projeto e migre um projeto existente que usa Taskfile."
slug: "experimental/wake"
sourcePath: "experimental/wake.md"
---

Wake é o sistema de compilação experimental do Wails. Um `wails.yaml` descreve seu aplicativo, frontend, plataformas de destino, pacotes e configurações de assinatura. O Wake usa essa configuração para planejar e executar compilações, armazenar em cache trabalhos que não foram alterados e informar os arquivos resultantes.

Neste tutorial, você criará um pequeno aplicativo para desktop, fará sua configuração no navegador, o compilará para seu próprio computador e o executará no modo de desenvolvimento. Comece com um único destino local e uma compilação não assinada. Adicione a compilação cruzada e a distribuição depois que essa primeira compilação funcionar.

Já tem um aplicativo? Acesse [Migrar um projeto existente](#migrar-um-projeto-existente). Para consultar campos específicos e fluxos de trabalho avançados, use o [guia de compilação com YAML](/experimental/yaml-builds/) e a [referência do manifesto](/reference/wails-yaml-fields/).

## O que esperar do experimento

@note{type="caution" title="Use um projeto de teste ou uma branch separada"}
O Wake é opcional, e sua configuração, seus comandos e seu comportamento podem mudar entre versões de prévia. Ele ainda não oferece as garantias de compatibilidade de um sistema de compilação estável. Mantenha uma revisão do projeto que você saiba que funciona e registre as versões da CLI e do framework que testar. Teste os pacotes nos sistemas de destino antes de distribuí-los.
@end

O Wake e o backend **Pure Go** são opções independentes. Habilitar o Wake não habilita automaticamente o Pure Go. Esta lição usa o backend padrão; portanto, as compilações de GUI para macOS e Linux precisam das ferramentas de desenvolvimento nativas habituais.

O assistente verifica os requisitos e salva a configuração. Salvar não executa uma compilação nem comprova que todos os destinos selecionados estão prontos. Ferramentas ou credenciais ausentes podem ser providenciadas posteriormente, por exemplo, em um host de CI. Combinações inválidas de destino e pacote devem ser corrigidas antes de salvar.

A migração é **unidirecional**. O Wake não executa Taskfiles, não expande modelos do Task nem converte YAML de volta em Taskfiles. Projetos não migrados continuam usando o runtime legado separado do Task. Um `wails.yaml` ativo nunca recorre silenciosamente a esse runtime.

## Antes de começar

Você precisa do Go, do Node.js com npm para o frontend inicial e das [dependências da plataforma](/getting-started/installation/#platform-specific-dependencies) para seu computador. Comece em um diretório no qual deseja criar a nova pasta do projeto. Você também precisa de um navegador para usar o assistente.

Use revisões da CLI e do framework que incluam esta prévia. Instalar a CLI `v3.0.0` estável não adiciona as telas de criação de projetos do Wake descritas aqui. Os comandos a seguir fixam a prévia da [PR nº 6531](https://github.com/wailsapp/wails/pull/6531):

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@b1daf95f3
wails3 version
```

Verifique se o diretório do executável do Go está em `PATH`; o [guia de instalação](/getting-started/installation/) explica como fazer isso. Mais adiante no tutorial, você fixará a dependência Go do aplicativo na mesma revisão. Ao testar uma prévia mais recente, atualize os dois argumentos de revisão juntos.

### Habilite o Wake no terminal

No Bash ou Zsh:

```bash
export WAILS_EXP_USE_WAKE=1
```

No PowerShell:

```powershell
$env:WAILS_EXP_USE_WAKE = '1'
```

Mantenha este terminal aberto durante o restante do tutorial. Defina a variável novamente em qualquer terminal novo ou tarefa de CI que execute comandos do Wake. A **presença** dela habilita o Wake: defini-la como `false` ou como uma string vazia não o desabilita. O nome antigo `WAILS_USE_WAKE` é ignorado.

## Crie o aplicativo

Execute o assistente de criação de projetos:

```sh
wails3 init -ui
```

Um assistente de projeto é aberto no navegador. Com o Wake habilitado, ele reúne as opções de código-fonte e compilação em um único fluxo. Nada é criado até você clicar no botão final **Criar projeto**. `wails3 setup` é um assistente separado para preferências globais, ferramentas e credenciais de assinatura; ele não configura o perfil Wake deste projeto.

| Tela | Opção para este tutorial | O que acontece em seguida |
| --- | --- | --- |
| **Projeto — Criar um projeto Wails** | Insira `wake-demo` e clique em **Continuar**. | O projeto será criado em uma nova pasta `wake-demo` dentro do diretório exibido. |
| **Framework — Escolher um framework** | Selecione **Vanilla** e clique em **Continuar**. | Você obtém um pequeno frontend sem precisar passar por uma etapa de aprendizado específica de um framework. |
| **Linguagem — Preferência de linguagem** | Selecione **JavaScript** e clique em **Continuar**. | A tela **Bindings**, exclusiva do TypeScript, é ignorada. |
| **Detalhes — Detalhes do projeto** | Use `Wake Demo` como nome do produto, `com.example.wakedemo` como identificador do pacote e `0.1.0` como versão. Clique em **Continuar**. | Seus metadados permanecem no assistente enquanto você escolhe as configurações de compilação. |

Não use **Usar padrões globais** nesta lição: essa opção ignora as escolhas de framework e linguagem. Se você escolher TypeScript em outro projeto, a tela **Bindings** permitirá escolher interfaces ou classes para os bindings de dados Go gerados.

## Escolha as configurações de compilação antes da criação

Continue na mesma janela do navegador. As capturas de tela abaixo mostram este passo a passo no Linux. O destino do host e as ferramentas necessárias serão diferentes no Windows ou macOS.

### Aplicativo — O que você está criando?

Mantenha **Aplicativo para desktop** e escolha **Backend padrão**. Clique em **Continuar**.

![Tela Aplicativo do assistente de criação de projetos, com Aplicativo para desktop e Backend padrão selecionados.](/assets/screenshots/wake/application.png)

Esta lição cria uma GUI com um frontend web. **Pure Go (experimental)** é uma opção separada. Ela elimina a exigência de um compilador C para a compilação dos destinos de desktop compatíveis, desde que suas dependências também possam ser compiladas sem CGo. Ela não elimina as bibliotecas de runtime do destino nem todas as ferramentas de empacotamento e assinatura.

O backend é uma **configuração que abrange todo o projeto**. Outros perfis que herdam essas configurações também usam o backend selecionado.

### Destinos — Onde o aplicativo deve ser executado?

Mantenha apenas o destino correspondente ao computador que você está usando e clique em **Continuar**. Inicialmente, o assistente seleciona o sistema operacional e a arquitetura do host.

![Tela Destinos com apenas linux/amd64 selecionado para o host usado no tutorial.](/assets/screenshots/wake/targets.png)

Os nomes dos destinos combinam um sistema operacional e uma arquitetura de CPU: por exemplo, `windows/amd64`, `darwin/arm64` para um Mac com Apple Silicon ou `linux/amd64`. `darwin` significa macOS. Selecionar outro destino solicita um executável para essa máquina; isso não torna o executável compatível com a execução no computador atual.

Deixe arquiteturas adicionais, compilações universais para macOS e destinos móveis experimentais para depois. Você poderá adicioná-los a `wails.yaml` após a primeira compilação local.

### Pacotes — Escolha os artefatos da versão

Deixe as caixas de formatos de pacote desmarcadas e mantenha **Assinar artefatos da versão** desativado. Clique em **Continuar**.

![Tela Pacotes sem formatos adicionais de pacote para Linux selecionados e com a assinatura desativada.](/assets/screenshots/wake/packages.png)

Sem formatos adicionais selecionados, o Wake gera a saída padrão do destino:
um executável para desktop ou, no caso de uma GUI para macOS, um pacote de aplicativo macOS. Isso é
suficiente para verificar o aplicativo localmente. A seleção de formatos adiciona
saídas de distribuição como DMG, MSIX, NSIS, DEB ou AppImage, quando compatíveis.

A assinatura é opcional para esta primeira compilação local. Uma compilação sem assinatura
bem-sucedida não significa que o aplicativo esteja pronto para distribuição pública. Para um projeto
assinado, configure credenciais reutilizáveis com `wails3 setup` antes da criação ou
configure a assinatura posteriormente em `wails.yaml`. O assistente de criação usa as
credenciais configuradas; ele não altera suas preferências globais de assinatura.

O Wake pode usar arquivos de certificado configurados fora do projeto sem copiar
chaves privadas nem senhas para `wails.yaml`. As configurações explícitas de assinatura do projeto
têm precedência. Siga as [orientações para configurar a assinatura](/getting-started/setup/)
e o [guia de empacotamento e assinatura](/experimental/yaml-builds/#package-and-sign-for-distribution)
antes de preparar uma versão pública. A assinatura para dispositivos móveis é configurada em YAML.

### Ferramentas — Ferramentas para sua compilação

Leia os requisitos das suas seleções. A lista vem do plano de compilação do Wake;
não é uma solicitação para instalar todas as ferramentas em todas as plataformas.

![Tela Ferramentas mostrando os requisitos determinados para a compilação local do código-fonte no Linux.](/assets/screenshots/wake/tools.png)

Para esta lição, você precisará do Go e das ferramentas de frontend, além das ferramentas
de desenvolvimento nativas exigidas pelo backend de GUI padrão do seu host. O Docker não é
um requisito geral. Uma compilação cruzada com CGo pode exigir um compilador e bibliotecas
do destino ou uma cadeia de ferramentas compatível em contêiner.

Resolva as linhas marcadas como **Não está pronto**, **Requer configuração** ou **Atualização necessária**
antes de tentar a compilação local. Use as instruções exibidas, quando disponíveis,
ou o [guia de instalação da plataforma](/getting-started/installation/) e clique em
**Verificar requisitos novamente**. Se uma instalação tiver alterado o `PATH` do terminal,
reinicie `wails3 init -ui` em um terminal atualizado e repita suas seleções.

O assistente pode criar um projeto válido mesmo que este computador não tenha as ferramentas
ou credenciais de assinatura necessárias. Isso é útil para configurar a CI, mas, neste tutorial,
você precisa de uma compilação que possa ser executada aqui. Configurações inválidas devem ser corrigidas
antes da criação.

### Criar o projeto

Confira suas escolhas: aplicativo para desktop, backend padrão, destino correspondente ao seu sistema host e assinatura desativada. Em **Ferramentas**, clique em **Criar projeto** e aguarde a mensagem **Escolhas do projeto aceitas**. Volte ao terminal e aguarde até que a CLI termine de criar o projeto; a confirmação no navegador apenas aceita suas escolhas.

![Tela de escolhas do projeto aceitas, mostrando o caminho planejado do aplicativo e os comandos a serem usados após a criação ser concluída com sucesso.](/assets/screenshots/wake/save.png)

O assistente cria o código-fonte e salva as escolhas de destino, pacote e assinatura
em `profiles.setup`, dentro de `wails.yaml`. Ele também define o tipo de aplicativo,
a ativação do frontend e o backend de todo o projeto. Você pode fechar o navegador e
voltar ao terminal. Aguarde a conclusão da criação do projeto e execute:

```sh
cd wake-demo
go get github.com/wailsapp/wails/v3@b1daf95f3
```

Agora você deve ter `main.go`, um módulo Go, um diretório `frontend` e um
`wails.yaml` nativo. O modelo integrado usa o Wake diretamente; você não precisa
migrar este novo projeto. Um modelo da comunidade pode ter tarefas legadas personalizadas
e, nesse caso, exigir a revisão de migração descrita abaixo.

A tela de conclusão mostra `wails3 build setup` e `wails3 dev`. O nome do perfil
é importante: `wails3 build setup` usa as seleções que você acabou de salvar. O comando simples
`wails3 build` não seleciona automaticamente o perfil `setup`. O nome do perfil
não significa que você precise executar `wails3 setup`.

## Compilar e executar seu primeiro aplicativo

Execute o comando mostrado pelo assistente:

```sh
wails3 build setup
```

A primeira compilação instala as dependências do frontend, gera os bindings e compila
o aplicativo. Compilações posteriores podem reutilizar o trabalho que não mudou. Aguarde **compilação
bem-sucedida** e consulte os caminhos de saída impressos pela CLI. O diretório de saída
padrão é `bin`; uma GUI para macOS é montada em um pacote `.app`.

Inicie o modo de desenvolvimento no seu próprio computador:

```sh
wails3 dev
```

Você deverá ver a janela do aplicativo inicial. Digite seu nome e teste a
saudação. Altere o placeholder do campo de nome em `frontend/index.html` e salve
o arquivo; o frontend deverá ser atualizado. Edite a saudação em `greetservice.go`
para acionar uma nova compilação do Go. Quando terminar, pressione **Ctrl+C** no terminal.

O desenvolvimento é um fluxo de trabalho executado no host local. Um perfil de versão pode conter vários
destinos de compilação cruzada sem alterar qual máquina pode executar a janela de
desenvolvimento. Para compilar e iniciar localmente um aplicativo de produção, use
`wails3 run`; consulte a [referência de execução](/reference/wails-yaml-run/).

Agora você tem um projeto com código-fonte, um perfil de compilação revisado e um ciclo local
de desenvolvimento funcional. Adicione o código-fonte e `wails.yaml` ao controle de versão.
Os arquivos de trabalho `.wails/` gerados não são a fonte da sua configuração.

## O que você pode fazer a seguir

| Objetivo | Compatibilidade do Wake | O que verificar primeiro |
| --- | --- | --- |
| Compilar para outro sistema desktop | Destinos Windows, macOS e Linux, incluindo ambas as arquiteturas do macOS e saída universal. | Go puro e CGo têm requisitos de compilador diferentes. Teste no sistema operacional e na arquitetura de destino. |
| Criar instaladores ou pacotes | MSIX, NSIS, aplicativo/DMG para macOS, DEB, RPM, Arch, Void e AppImage. | O NSIS requer `makensis`. O AppImage requer um runtime de destino; o empacotamento cruzado de GUI também requer uma árvore de runtime do Linux fornecida. |
| Assinar uma versão | Assinatura portátil com chave exportada para Windows, macOS e DEB/RPM/Arch; os repositórios nativos de identidades continuam disponíveis. | Cada identidade, formato e host tem limitações. AppImage e Void não têm um fluxo de assinatura integrado. |
| Notarizar um aplicativo macOS | Envio com chave de API da equipe e anexação do tíquete, ou um perfil das Chaves do macOS. | Uma assinatura de Developer ID adequada e um carimbo de data e hora ainda são obrigatórios. O fato de assinaturas de teste geradas pelo Go passarem nas verificações nativas não prova que a Apple aceitará sua versão. |
| Compilar vários programas | Um programa principal de GUI, CLI ou servidor, além de programas auxiliares nomeados. | Declare os pacotes e tipos dos programas auxiliares no YAML. O assistente não gera o código-fonte deles. |
| Empacotar arquivos compilados em outro lugar | `wails3 build setup --skip-build`. | Forneça nomes de executáveis e arquiteturas correspondentes em `build.output`. Essa opção é uma escolha feita uma única vez no comando, não uma configuração permanente do perfil. |
| Adicionar automação de versões | Perfis nomeados, hooks de ciclo de vida limitados, cache e trabalho paralelo de empacotamento. | Use YAML literal e scripts declarados que pertençam ao projeto. Comportamentos arbitrários do Taskfile não são um mecanismo de extensão do Wake. |

Edite `profiles.setup`, outros perfis nomeados e configurações avançadas em
`wails.yaml`. Use `wails3 setup` para alterar preferências e credenciais compartilhadas.
Consulte o
[guia de compilação com YAML](/experimental/yaml-builds/) para ver exemplos completos e requisitos
de pacotes. Aplicativos com GUI ainda precisam das bibliotecas de sistema do destino;
o Go puro não os transforma em aplicativos completamente estáticos.

iOS e Android continuam experimentais, exigem seus respectivos SDKs de plataforma e CGo e
não fazem parte desta primeira lição sobre desktop. O assistente não torna as compilações para dispositivos móveis
portáteis entre hosts arbitrários.

## Migrar um projeto existente

Trabalhe em um branch separado ou em uma cópia da árvore de trabalho de um projeto **Wails v3**. Esta é uma
migração da configuração de compilação, não uma atualização do Wails v2. Preserve o resultado
da compilação existente e os testes para poder comparar o comportamento após a ativação.

@note{type="note" title="Comece com uma migração revisada"}
Projetos novos gerados a partir dos modelos legados integrados atuais podem ser migrados antes da compilação do frontend. O Wake reconhece os seletores de tarefas padrão e o ciclo de vida de desenvolvimento desses projetos sem executar Taskfiles nem avaliar expressões personalizadas.

Tarefas personalizadas, substituições de variáveis e políticas de desenvolvimento ainda exigem revisão. Se a análise não conseguir representar seu projeto, mantenha-o no sistema de compilação legado e informe os diagnósticos. Não remova comportamentos necessários apenas para que a ativação seja aprovada. O tutorial de criação de projetos acima não exige migração.
@end

### 1. Analise sem alterar o sistema de compilação ativo

Ative o Wake e use a CLI de pré-visualização indicada no início desta página. Na raiz do projeto existente, execute:

```sh
wails3 migrate --dry-run
```

Leia as configurações compatíveis e os impedimentos informados. O analisador lê as declarações legadas, mas não executa seus comandos. A simulação não cria um manifesto ativo. Para obter diagnósticos legíveis por máquina, use `wails3 migrate --json`; use `--dry-run --json` se também quiser evitar a gravação de um rascunho.

### 2. Crie e revise a proposta inativa

```sh
wails3 migrate
```

Isso cria **`wails.migrated.yaml`**, e não `wails.yaml`. Examine-o junto com os Taskfiles existentes e o `build/config.yml`. Verifique os metadados do projeto, os comandos do frontend, as arquiteturas de destino, o conteúdo dos pacotes, os ativos e as configurações de assinatura. Mantenha os segredos em variáveis de ambiente ou no seu repositório de credenciais.

Tarefas integradas conhecidas podem ser mapeadas para configurações nativas. Tarefas conservadoras do ciclo de vida, sem argumentos, que chamam um script pertencente ao projeto podem se tornar hooks. Outros comportamentos de tarefas personalizadas alcançáveis constituem um impedimento; o Wake não os descarta silenciosamente nem os executa por meio do Task como alternativa.

Se já existir um rascunho, a migração não sobrescreverá suas edições revisadas. Guarde uma cópia antes de substituí-lo ou gerá-lo novamente de forma deliberada. Não renomeie um rascunho bloqueado como `wails.yaml` para contornar as verificações de ativação. Reduza ou substitua o comportamento legado sem suporte e execute novamente a análise antes da ativação. Se ainda não for possível representar o projeto, mantenha seu sistema de compilação legado e informe essa lacuna.

### 3. Ative explicitamente

Depois de revisar o rascunho e resolver os impedimentos, execute:

```sh
wails3 migrate --activate --dry-run
wails3 migrate --activate
```

O primeiro comando verifica a ativação sem publicá-la. O segundo verifica novamente as declarações legadas e valida o rascunho antes de criar `wails.yaml` de forma atômica; uma ativação bem-sucedida remove o `wails.migrated.yaml` inativo. Ele não altera Taskfiles, `build/config.yml`, scripts nem ativos. A partir desse ponto, o manifesto ativo controla as compilações do Wake. Os arquivos antigos não constituem uma segunda fonte de configuração.

Alguns projetos mais antigos já têm um `wails.yaml` somente de metadados, com uma seção `info:` no nível superior. Embora não seja um manifesto nativo do Wake, ele ainda ocupa o destino da ativação. Preserve e mova deliberadamente esse arquivo legado antes da ativação; o migrador não o sobrescreverá nem o renomeará para você.

### 4. Verifique o aplicativo migrado

```sh
wails3 config check
wails3 build
wails3 dev
```

Compare o aplicativo e o conteúdo dos pacotes com a compilação anterior, incluindo seu comportamento personalizado de lançamento. Execute os testes do projeto. Verifique explicitamente cada perfil migrado, por exemplo, `wails3 config check release` e `wails3 build release` quando o rascunho contiver um perfil `release`.

Edite os perfis migrados diretamente em `wails.yaml`. `wails3 init -ui` cria projetos novos; `wails3 setup` gerencia preferências globais. Nenhum deles edita o perfil de lançamento de um projeto existente.

### Se precisar interromper o experimento

Remover a definição de `WAILS_EXP_USE_WAKE` desativa a adesão ao experimento, mas um `wails.yaml` ativo ainda impede o roteamento automático para o Taskfile. Um YAML inválido também interrompe o processo com um erro. Não há comando de migração reversa.

Para voltar ao fluxo de trabalho antigo, restaure no controle de versão a revisão completa do projeto anterior à migração, inclusive a ausência do manifesto ativo, e use a CLI e as dependências compatíveis com ela. Antes disso, preserve todas as novas alterações de código-fonte de que precisar. Manter uma branch de teste separada torna essa recuperação deliberada muito mais fácil.

## Diagnostique um problema e envie feedback

| Sintoma | Primeira verificação |
| --- | --- |
| A criação do projeto ignora Aplicativo, Destinos, Pacotes e Ferramentas | Defina `WAILS_EXP_USE_WAKE` no terminal que inicia o processo e verifique qual executável `wails3` está no `PATH`. O `wails3 setup` global sempre mantém seu fluxo de preferências. |
| Um modelo da comunidade precisa de migração | Revise a configuração legada dele com os comandos de migração antes da ativação. |
| Uma compilação salva informa que faltam ferramentas | A criação de um projeto registra a intenção. Instale os requisitos de host informados antes de compilar; a criação não instala todas as ferramentas selecionadas. |
| Um executável gerado por compilação cruzada não inicia neste ambiente | Execute-o no sistema operacional e na arquitetura selecionados; use o destino do host para o desenvolvimento local. |
| O modo pré-compilado não encontra um binário | Verifique os [nomes de arquivo e caminhos dos binários pré-compilados](/experimental/yaml-builds/#use-prebuilt-desktop-executables), inclusive os subdiretórios para múltiplos destinos. |
| Uma alteração no manifesto falha | Execute `wails3 config check` para o perfil relevante e siga o diagnóstico do campo. O Wake não recorre aos Taskfiles como alternativa. |

Inspecione uma compilação antes de executá-la com `wails3 build setup --plan --json`. `wails3 config check` valida a configuração, e `wails3 eject` grava um snapshot `wails.ejected.yaml` inativo sem alterar o manifesto ativo. Revise a saída do comando antes de compartilhá-la: os caminhos do projeto e os locais das credenciais podem identificar seu ambiente local.

Compartilhe seu feedback nas [Discussões sobre experimentos](https://github.com/wailsapp/wails/discussions/categories/experiments) ou [registre um bug reproduzível](https://github.com/wailsapp/wails/issues). Inclua as revisões da CLI e do framework, o sistema operacional e a arquitetura do host, o destino selecionado, o comando exato, se você selecionou Pure Go e o diagnóstico relevante. Informe se a falha ocorreu durante a compilação, o empacotamento, a assinatura ou a execução do aplicativo. Quando for útil, inclua um manifesto mínimo sem os segredos e a saída de `wails3 doctor`.
