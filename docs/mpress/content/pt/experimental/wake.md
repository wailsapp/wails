---
title: "Wake"
description: "Crie seu primeiro projeto Wake, siga o assistente de configuração, compile e execute o projeto e migre um projeto existente com Taskfile."
slug: "experimental/wake"
sourcePath: "experimental/wake.md"
---

Wake é o sistema de compilação experimental do Wails. Um único `wails.yaml` descreve seu aplicativo, o frontend, as plataformas de destino, os pacotes e as configurações de assinatura. O Wake usa essa configuração para planejar e executar compilações, armazenar em cache o trabalho que não mudou e informar os arquivos resultantes.

Neste tutorial, você criará um pequeno aplicativo para desktop, o configurará no navegador, o compilará para seu próprio computador e o executará no modo de desenvolvimento. Comece com um único destino local e uma compilação não assinada. Adicione a compilação cruzada e a distribuição depois que essa primeira compilação funcionar.

Já tem um aplicativo? Acesse [Migrar um projeto existente](#migrar-um-projeto-existente). Para conhecer os campos individuais e os fluxos de trabalho avançados, consulte o [guia de compilação em YAML](/experimental/yaml-builds/) e a [referência do manifesto](/reference/wails-yaml-fields/).

## O que esperar do experimento

@note{type="caution" title="Use um projeto de teste ou uma ramificação separada"}
O uso do Wake é opcional, e sua configuração, seus comandos e seu comportamento podem mudar entre compilações de prévia. Ele ainda não oferece as garantias de compatibilidade de um sistema de compilação estável. Mantenha uma revisão comprovadamente funcional do projeto e registre as versões da CLI e do framework que você testar. Teste os pacotes nos sistemas de destino antes de distribuí-los.
@end

O Wake e o backend **Pure Go** são escolhas independentes. Ativar o Wake não ativa o Pure Go automaticamente. Esta lição usa o backend padrão; portanto, as compilações de GUI para macOS e Linux precisam das ferramentas nativas de desenvolvimento habituais.

O assistente verifica os requisitos e salva a configuração. Salvar não realiza uma compilação nem comprova que todos os destinos selecionados estão prontos. Ferramentas ou credenciais ausentes podem ser providenciadas posteriormente, por exemplo, em um host de CI. Combinações inválidas de destino e pacote devem ser corrigidas antes de salvar.

A migração é **unidirecional**. O Wake não executa Taskfiles, não expande modelos do Task nem converte YAML de volta em Taskfiles. Projetos não migrados continuam usando o runtime legado separado do Task. Um `wails.yaml` ativo nunca recorre silenciosamente a ele.

## Antes de começar

Você precisa do Go, do Node.js com npm para o frontend inicial e das [dependências da plataforma](/getting-started/installation/#platform-specific-dependencies) do seu computador. Comece em um diretório no qual deseja criar a nova pasta do projeto. Você também precisa de um navegador para usar o assistente.

Use revisões da CLI e do framework que incluam esta prévia. Instalar a CLI `v3.0.0` estável não adiciona as telas de configuração do Wake descritas aqui. Os comandos a seguir fixam a prévia da [PR nº 6531](https://github.com/wailsapp/wails/pull/6531):

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@44e449c91
wails3 version
```

Verifique se o diretório do executável do Go está em `PATH`; o [guia de instalação](/getting-started/installation/) explica como fazer isso. Mais adiante no tutorial, você fixará a dependência Go do aplicativo na mesma revisão. Ao testar uma prévia mais recente, atualize os dois argumentos de revisão em conjunto.

### Ativar o Wake no terminal

No Bash ou Zsh:

```bash
export WAILS_EXP_USE_WAKE=1
```

No PowerShell:

```powershell
$env:WAILS_EXP_USE_WAKE = '1'
```

Mantenha este terminal aberto durante o restante do tutorial. Defina a variável novamente em qualquer novo terminal ou trabalho de CI que execute comandos do Wake. A **presença** dela ativa o Wake: defini-la como `false` ou como uma string vazia não o desativa. O nome antigo `WAILS_USE_WAKE` é ignorado.

## Criar o aplicativo

Execute o assistente de criação de projetos:

```sh
wails3 init -ui
```

Um assistente de projeto é aberto no navegador. Ele cria os arquivos-fonte; o assistente `wails3 setup` separado abaixo configura como compilá-los.

| Tela | Opção para este tutorial | O que acontece em seguida |
| --- | --- | --- |
| **Projeto — Criar um projeto Wails** | Insira `wake-demo` e selecione **Continuar**. | O projeto será criado em uma nova pasta `wake-demo` dentro do diretório exibido. |
| **Framework — Escolher um framework** | Selecione **Vanilla** e depois **Continuar**. | Você obtém um pequeno frontend sem precisar aprender primeiro os conceitos específicos de um framework. |
| **Linguagem — Preferência de linguagem** | Selecione **JavaScript** e depois **Continuar**. | A tela **Bindings**, exclusiva do TypeScript, é ignorada. |
| **Detalhes — Detalhes do projeto** | Use o nome de produto `Wake Demo`, o identificador do pacote `com.example.wakedemo` e a versão `0.1.0`. Clique em **Criar projeto**. | Com o Wake ativado, os detalhes são gravados em `wails.yaml`. |

Não use **Usar padrões globais** nesta lição: essa opção ignora as escolhas de framework e linguagem. Se você escolher TypeScript em outro projeto, a tela **Bindings** permitirá escolher interfaces ou classes para os bindings de dados Go gerados.

Volte ao terminal e aguarde a conclusão da criação do projeto. Em seguida, execute:

```sh
cd wake-demo
go get github.com/wailsapp/wails/v3@44e449c91
```

Agora você deve ter `main.go`, um módulo Go, um diretório `frontend` e um `wails.yaml` nativo. O modelo integrado usa o Wake diretamente; não é necessário migrar esse novo projeto. Um modelo da comunidade pode ter tarefas legadas personalizadas e, nesse caso, talvez exija a revisão de migração descrita abaixo.

## Configurar a compilação, tela por tela

Dentro de `wake-demo`, execute:

```sh
wails3 setup
```

O navegador agora exibe **CONFIGURAÇÃO DO WAKE**. Se ele exibir o fluxo de configuração geral, verifique se `WAILS_EXP_USE_WAKE` está definido no terminal que iniciou o comando e se você está executando a CLI de prévia.

As capturas de tela abaixo mostram este passo a passo no Linux. O destino do host e as ferramentas necessárias serão diferentes no Windows ou no macOS.

### Projeto — somente quando faltarem metadados

Para o aplicativo que você acabou de criar, esta tela é **ignorada**: `init` já gravou os metadados. Esse é o comportamento esperado.

Ao configurar um diretório de código-fonte Go existente sem um manifesto nativo nem uma configuração de compilação legada, a tela **Sobre seu projeto** solicita o nome do produto, o nome do projeto, o identificador do aplicativo, a versão, a empresa e a descrição. Esses valores permanecem no assistente até você selecionar **Salvar**. A configuração não cria arquivos-fonte do aplicativo. Se encontrar Taskfiles ou metadados legados, faça primeiro a migração.

### Aplicativo — O que você está criando?

Escolha **Aplicativo para desktop**. Deixe **Empacotar binários existentes** desmarcado e escolha **Backend padrão**. Clique em **Continuar**.

![Tela Aplicativo com Aplicativo para desktop e Backend padrão selecionados e o empacotamento de binários pré-compilados desativado.](/assets/screenshots/wake/application.png)

- **Aplicativo para desktop** compila uma GUI com um frontend web.
- **Aplicativo de linha de comando** e **Servidor** descrevem programas sem uma janela de desktop. Alterar essa seleção não reescreve o código-fonte do aplicativo; use um projeto de código-fonte que corresponda ao tipo escolhido.
- **Pure Go (experimental)** é uma opção independente. Ela elimina a exigência de um compilador C para as compilações de desktop compatíveis, desde que suas dependências também sejam compiladas sem CGo. Ela não elimina as bibliotecas de runtime do destino nem todas as ferramentas de empacotamento e assinatura.
- **Empacotar binários existentes** prepara um comando usando `--skip-build`.
  Essa opção exige binários já compilados; deixe-a desativada em um projeto novo.

O backend é uma **configuração válida para todo o projeto**. Outros perfis que herdam essas
configurações podem mudar de comportamento quando você a altera aqui.

### Destinos — Onde o aplicativo deve ser executado?

Mantenha apenas o destino correspondente ao computador que você está usando e clique em **Continuar**.
Inicialmente, o assistente seleciona o sistema operacional e a arquitetura da máquina host.

![Tela Destinos com apenas linux/amd64 selecionado para a máquina host do tutorial.](/assets/screenshots/wake/targets.png)

Os nomes dos destinos combinam um sistema operacional e uma arquitetura de CPU: por exemplo,
`windows/amd64`, `darwin/arm64` para um Mac com Apple Silicon ou `linux/amd64`.
`darwin` significa macOS. Selecionar outro destino solicita um executável para essa
máquina; isso não torna o executável compatível com a execução no computador atual.

Deixe as arquiteturas adicionais, as compilações universais para macOS e os destinos móveis
experimentais para uma execução posterior da configuração.

### Pacotes — Escolha os artefatos da versão

Deixe desmarcadas as caixas de formatos de pacote e mantenha **Assinar artefatos da versão** desativado.
Clique em **Continuar**.

![Tela Pacotes sem formatos adicionais de pacote Linux selecionados e com a assinatura desativada.](/assets/screenshots/wake/packages.png)

Sem formatos adicionais selecionados, o Wake produz o artefato padrão do destino:
um executável para desktop ou um pacote de aplicativo macOS para uma GUI do macOS. Isso é
suficiente para verificar o aplicativo localmente. Selecionar formatos posteriormente adiciona
artefatos de distribuição como DMG, MSIX, NSIS, DEB ou AppImage, quando compatíveis.

A assinatura é opcional nesta primeira compilação local. Uma compilação sem assinatura
bem-sucedida não significa que o aplicativo esteja pronto para distribuição pública.

### Ferramentas — Ferramentas para a compilação

Leia os requisitos das suas seleções. A lista é gerada a partir do plano de compilação
do Wake; ela não solicita a instalação de todas as ferramentas em todas as plataformas.

![Tela Ferramentas mostrando os requisitos resolvidos para a compilação local do código-fonte no Linux.](/assets/screenshots/wake/tools.png)

Nesta lição, você precisará do Go e das ferramentas de frontend, além das ferramentas
de desenvolvimento nativas exigidas pelo backend de GUI padrão da máquina host. O Docker não é
um requisito geral. Uma compilação cruzada com CGo pode exigir um compilador e bibliotecas
do destino ou uma cadeia de ferramentas de contêiner compatível.

Resolva as linhas marcadas como **Não está pronto**, **Requer configuração** ou **Atualização necessária**
antes de tentar a compilação local. Use as instruções exibidas, quando disponíveis,
ou o [guia de instalação da plataforma](/getting-started/installation/) e clique em
**Verificar novamente**. Talvez seja necessário reiniciar a configuração se uma instalação tiver alterado o
`PATH` do terminal. Quando os requisitos da máquina host estiverem prontos, clique em **Continuar**.

O assistente pode salvar uma configuração válida mesmo que este computador não tenha as ferramentas necessárias.
Isso é útil para configurar a CI, mas, neste tutorial, você precisa de uma compilação
que possa ser executada nesta máquina.

### Assinatura — somente quando selecionada

Esta tela é **ignorada** porque você deixou a assinatura desativada. Em uma execução posterior, ativar
**Assinar artefatos da versão** exibe os controles de assinatura existentes para as plataformas
de desktop selecionadas. **Configurar a assinatura depois** desativa a assinatura no perfil
de configuração salvo. A assinatura para dispositivos móveis é configurada no YAML, não nesses controles.

Os padrões de assinatura ficam nas configurações do usuário. O Wake pode usar arquivos de certificado
configurados fora do projeto sem copiar chaves privadas nem senhas para
`wails.yaml`. As configurações explícitas de assinatura do projeto têm precedência. Siga as
[orientações de configuração da assinatura](/getting-started/setup/) e o
[guia de empacotamento e assinatura](/experimental/yaml-builds/#package-and-sign-for-distribution)
antes de preparar uma versão pública.

### Salvar — Salve a configuração de compilação

Confira o resumo: aplicativo para desktop, backend padrão, compilação do código-fonte, destino
da máquina host e assinatura desativada. Clique em **Salvar wails.yaml**.

![Tela Salvar resumindo o aplicativo com GUI, o backend padrão, a compilação do código-fonte e o destino linux/amd64 sem assinatura.](/assets/screenshots/wake/save.png)

Você deverá ver **Configuração de compilação salva**, o caminho do manifesto e:

```sh
wails3 build setup
```

Clique em **Concluir** e volte ao terminal. A configuração salvou as opções de destino e
de pacote em `profiles.setup`. Ela também atualiza o tipo do aplicativo,
a ativação do frontend e as configurações do backend. Os comandos de frontend, programas auxiliares,
ganchos e configurações não relacionadas já existentes são preservados.

O nome do perfil é importante: `wails3 build setup` usa essas seleções salvas.
`wails3 build`, sem opções adicionais, não seleciona automaticamente o perfil `setup`.

## Compile e execute seu primeiro aplicativo

Execute o comando exibido pelo assistente:

```sh
wails3 build setup
```

A primeira compilação instala as dependências do frontend, gera os bindings e compila
o aplicativo. As compilações posteriores podem reutilizar o trabalho que não foi alterado. Aguarde por **compilação
bem-sucedida** e leia os caminhos dos artefatos exibidos pela CLI. O diretório padrão
de saída é `bin`; uma GUI do macOS é montada em um pacote `.app`.

Inicie o modo de desenvolvimento no seu próprio computador:

```sh
wails3 dev
```

Você deverá ver a janela do aplicativo inicial. Digite seu nome e teste a
saudação. Altere o texto de placeholder do campo de nome em `frontend/index.html` e salve
o arquivo; o frontend deverá ser atualizado. Edite a saudação em `greetservice.go`
para testar uma recompilação do Go. Ao terminar, pressione **Ctrl+C** no terminal.

O desenvolvimento é um fluxo de trabalho realizado na máquina local. Um perfil de versão pode conter vários
destinos de compilação cruzada sem alterar qual máquina pode executar a janela
de desenvolvimento. Para compilar e iniciar localmente um aplicativo de produção, use
`wails3 run`; consulte a [referência de execução](/reference/wails-yaml-run/).

Agora você tem um projeto com código-fonte, um perfil de compilação revisado e um ciclo local
de desenvolvimento funcional. Confirme o código-fonte e `wails.yaml` no controle de versão.
Os arquivos de trabalho `.wails/` gerados não são a fonte da sua configuração.

## Próximas etapas

| Objetivo | Compatibilidade no Wake | O que verificar primeiro |
| --- | --- | --- |
| Compilar para outro sistema de desktop | Destinos Windows, macOS e Linux, incluindo as duas arquiteturas do macOS e o artefato universal. | Pure Go e CGo têm requisitos de compilador diferentes. Teste no sistema operacional e na arquitetura de destino. |
| Criar instaladores ou pacotes | MSIX, NSIS, aplicativo macOS/DMG, DEB, RPM, Arch, Void e AppImage. | O NSIS requer `makensis`. O AppImage requer um runtime de destino; o empacotamento cruzado de GUI também requer uma árvore de runtime do Linux fornecida. |
| Assinar uma versão | Assinatura portátil com chaves exportadas para Windows, macOS e DEB/RPM/Arch; os repositórios nativos de identidades continuam disponíveis. | Cada identidade, formato e máquina host tem limitações. AppImage e Void não têm um fluxo de assinatura integrado. |
| Submeter um aplicativo macOS à notarização | Envio pela equipe com uma chave de API e anexação do tíquete, ou um perfil do chaveiro do macOS. | Uma assinatura de Developer ID adequada e um carimbo de data e hora ainda são necessários. O fato de assinaturas de teste geradas pelo Go passarem nas verificações nativas não comprova que a Apple aceitará sua versão. |
| Compilar vários programas | Um programa principal com GUI, de CLI ou servidor, além de programas auxiliares nomeados. | Declare os pacotes e tipos dos programas auxiliares no YAML. O assistente não gera o código-fonte deles. |
| Empacotar arquivos compilados em outro lugar | `wails3 build setup --skip-build`. | Forneça nomes de executáveis e arquiteturas correspondentes em `build.output`. Essa opção é uma escolha feita uma única vez no comando, não uma configuração permanente. |
| Adicionar automação de versões | Perfis nomeados, hooks de ciclo de vida limitados, cache e processamento paralelo de pacotes. | Use YAML literal e scripts declarados que pertençam ao projeto. Comportamentos arbitrários do Taskfile não são um mecanismo de extensão do Wake. |

Reabra `wails3 setup` para alterar o perfil `setup`. Outros perfis nomeados e configurações avançadas são editados em `wails.yaml`. Consulte o [guia de compilação com YAML](/experimental/yaml-builds/) para ver exemplos completos e os requisitos dos pacotes. Aplicações com interface gráfica ainda precisam das bibliotecas de sistema do destino; Pure Go não as transforma em aplicações totalmente estáticas.

iOS e Android continuam experimentais, exigem os respectivos SDKs de plataforma e CGo e estão fora do escopo desta primeira lição sobre desktop. O assistente não torna as compilações para dispositivos móveis portáveis entre hosts arbitrários.

## Migrar um projeto existente

Trabalhe em uma ramificação separada ou em uma cópia de uma árvore de trabalho funcional de um projeto **Wails v3**. Esta é uma migração da configuração de compilação, não uma atualização do Wails v2. Preserve o resultado da compilação existente e os testes para poder comparar o comportamento após a ativação.

@note{type="caution" title="A cobertura da migração está incompleta"}
A prévia fixada pode informar bloqueios até mesmo para um projeto gerado com os modelos legados padrão atuais. Exemplos conhecidos incluem tarefas selecionadas dinamicamente e um campo `exit_policy` não reconhecido na configuração de compilação legada. Não remova comportamentos necessários apenas para permitir que a ativação seja concluída. Se a análise não conseguir representar seu projeto, mantenha-o na compilação legada e informe os diagnósticos. O tutorial de novo projeto acima não requer migração.
@end

### 1. Analisar sem alterar a compilação ativa

Ative o Wake e use a CLI de prévia indicada no início desta página. Na raiz do projeto existente, execute:

```sh
wails3 migrate --dry-run
```

Leia as configurações compatíveis e os bloqueios informados. O analisador lê as declarações legadas; ele não executa os comandos delas. A simulação não cria um manifesto ativo. Para obter diagnósticos legíveis por máquina, use `wails3 migrate --json`; use `--dry-run --json` se também quiser evitar a gravação de um rascunho.

### 2. Criar e revisar a proposta inativa

```sh
wails3 migrate
```

Isso cria **`wails.migrated.yaml`**, não `wails.yaml`. Inspecione-o junto com os Taskfiles existentes e `build/config.yml`. Verifique os metadados do projeto, os comandos do frontend, as arquiteturas de destino, o conteúdo dos pacotes, os recursos e as configurações de assinatura. Mantenha os segredos em variáveis de ambiente ou no seu armazenamento de credenciais.

Tarefas internas conhecidas podem ser mapeadas para configurações nativas. Tarefas conservadoras de ciclo de vida, sem argumentos, que chamem um script pertencente ao projeto podem se tornar hooks. Outros comportamentos de tarefas personalizadas acessíveis constituem um bloqueio; o Wake não os descarta silenciosamente nem os executa por meio do Task como alternativa.

Se já houver um rascunho, a migração não substituirá as edições que você revisou. Guarde uma cópia antes de substituí-lo ou gerá-lo novamente de propósito. Não renomeie um rascunho bloqueado como `wails.yaml` para contornar as verificações de ativação. Reduza ou substitua os comportamentos legados incompatíveis e execute novamente a análise antes da ativação. Se o projeto ainda não puder ser representado, mantenha a compilação legada e informe essa lacuna.

### 3. Ativar explicitamente

Depois de revisar o rascunho e resolver os bloqueios, execute:

```sh
wails3 migrate --activate --dry-run
wails3 migrate --activate
```

O primeiro comando verifica a ativação sem publicá-la. O segundo verifica novamente as declarações legadas e valida o rascunho antes de criar `wails.yaml` de forma atômica; uma ativação bem-sucedida remove o `wails.migrated.yaml` inativo. Ele não altera Taskfiles, `build/config.yml`, scripts nem recursos. A partir desse ponto, o manifesto ativo controla as compilações do Wake. Os arquivos antigos não são uma segunda fonte de configuração.

Alguns projetos mais antigos já têm um `wails.yaml` somente de metadados com uma seção `info:` no nível superior. Embora não seja um manifesto nativo do Wake, ele ainda ocupa o destino da ativação. Preserve e mova deliberadamente esse arquivo legado antes da ativação; o migrador não o substituirá nem o renomeará por você.

### 4. Verificar a aplicação migrada

```sh
wails3 config check
wails3 build
wails3 dev
```

Compare a aplicação e o conteúdo dos pacotes com a compilação anterior, incluindo o comportamento personalizado de lançamento. Execute os testes do projeto. Verifique explicitamente cada perfil migrado, por exemplo, `wails3 config check release` e `wails3 build release` quando o rascunho contiver um perfil `release`.

Agora você pode executar `wails3 setup` para adicionar ou substituir o perfil dedicado `setup`. Isso não seleciona nem edita automaticamente um perfil `release` existente.

### Se precisar interromper o experimento

Remover a definição de `WAILS_EXP_USE_WAKE` desativa a adesão voluntária, mas um `wails.yaml` ativo ainda impede o roteamento automático para Taskfiles. Um YAML inválido também interrompe o processo com um erro. Não há comando de migração reversa.

Para voltar ao fluxo de trabalho antigo, restaure no controle de versão a revisão completa do projeto anterior à migração, incluindo a ausência do manifesto ativo, e use a CLI e as dependências compatíveis com ela. Antes disso, preserve todas as novas alterações no código-fonte de que você precise. Manter uma ramificação de teste separada facilita muito essa recuperação deliberada.

## Diagnosticar um problema e enviar feedback

| Sintoma | Primeira verificação |
| --- | --- |
| A configuração mostra o assistente geral | Defina `WAILS_EXP_USE_WAKE` no terminal que inicia o processo e verifique qual executável `wails3` está em `PATH`. |
| A configuração solicita uma migração | Você abriu um projeto legado ou um modelo da comunidade com configuração legada. Revise a migração antes da ativação. |
| Uma compilação salva informa que faltam ferramentas | O salvamento registra a intenção. Volte para **Ferramentas**, corrija os requisitos do host e verifique novamente. |
| Um executável compilado para outra plataforma não inicia neste ambiente | Execute-o no sistema operacional e na arquitetura selecionados; use o destino do host para desenvolvimento local. |
| O modo pré-compilado não encontra um binário | Verifique os [nomes de arquivo e caminhos pré-compilados](/experimental/yaml-builds/#use-prebuilt-desktop-executables), incluindo os subdiretórios de vários destinos. |
| Uma alteração no manifesto falha | Execute `wails3 config check` para o perfil relevante e siga o diagnóstico do campo. O Wake não recorre aos Taskfiles como alternativa. |

Inspecione uma compilação antes de executá-la com `wails3 build setup --plan --json`.
`wails3 config check` valida a configuração, e `wails3 eject` grava um snapshot `wails.ejected.yaml` inativo sem alterar o manifesto ativo.
Revise a saída do comando antes de compartilhá-la: os caminhos do projeto e os locais das credenciais podem identificar seu ambiente local.

Compartilhe seu feedback em [Discussões sobre experimentos](https://github.com/wailsapp/wails/discussions/categories/experiments) ou [registre um bug reproduzível](https://github.com/wailsapp/wails/issues). Inclua as revisões da CLI e do framework, o sistema operacional e a arquitetura do host, o destino selecionado, o comando exato, se você selecionou Pure Go e o diagnóstico relevante. Informe se a falha ocorreu durante a compilação, o empacotamento, a assinatura ou a execução da aplicação. Quando for útil, inclua um manifesto mínimo sem segredos e a saída de `wails3 doctor`.
