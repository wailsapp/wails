---
title: "Configuração de execução em YAML"
description: "Padrões de execução local, substituições específicas por plataforma e comportamento de inicialização do wails3 run."
---

## Disponibilidade

Esta referência descreve a extensão experimental de execução. Ela não garante a disponibilidade em uma CLI lançada.

`wails3 run` compila e inicia o aplicativo usando os padrões locais do manifesto. A execução orientada pelo manifesto exige a habilitação [`WAILS_EXP_USE_WAKE=1`](/experimental/yaml-builds#enable-yaml-builds). Quando nenhum manifesto é encontrado, o comando executa `go run .` no diretório atual sem exigir essa habilitação. Um manifesto inválido ou ilegível é um erro e nunca aciona o fallback. Se houver um manifesto enquanto o YAML estiver desabilitado, o comando informará como habilitá-lo.

A referência gerada do campo [`wails.yaml`](/reference/wails-yaml-fields) lista os tipos de campo e os valores padrão. A compatibilidade com dispositivos móveis nativos e com bundles do macOS deve ser verificada no host apropriado antes que uma versão declare essas plataformas prontas.

## Contrato do comando

`wails3 run` localiza o `wails.yaml` mais próximo, resolve a configuração para um destino, compila por meio do pipeline existente com cache e inicia o resultado. A descoberta pesquisa a partir do diretório de trabalho em direção aos diretórios superiores, verificando a existência de um manifesto antes de parar em um limite `go.mod`. O diretório do manifesto é a raiz do projeto e o diretório de trabalho do aplicativo para desktop.

O destino padrão é o sistema operacional e a arquitetura do host. Destinos móveis exigem seleção explícita. A compilação cruzada não significa que um executável para desktop possa ser executado no host; destinos de desktop incompatíveis falham antes da compilação. Isso inclui executáveis compilados com `CGO_ENABLED=0`: uma compilação cruzada bem-sucedida somente em Go não fornece um emulador nem altera o contrato de inicialização local. Use [`wails3 build` para compilação cruzada](/experimental/yaml-builds/#cross-compile-desktop-applications) e execute ou instale os arquivos resultantes no sistema de destino. Use `--target android/arm64 --device <adb-serial>` ou `--target android/arm64 --emulator <avd-name>` para Android. A arquitetura selecionada deve corresponder ao dispositivo ou emulador. Use `--target ios/arm64 --device <udid>` para um simulador de iOS e adicione `--destination device` para um dispositivo físico. O planejamento e a execução para iOS exigem macOS e Xcode.

Com YAML, a compilação usa recursos incorporados e o fluxo de compilação para produção, que fornece automaticamente a tag `production`. `wails3 dev` continua responsável pelo recarregamento em tempo real; `dev.tags` não contribui para `wails3 run`.

```shell
wails3 run
wails3 run --tags extra_tag
wails3 run -- --example-mode other
wails3 run --plan
wails3 run --skip-build
```

`--plan` resolve e exibe o destino selecionado, as tags de compilação efetivas, o artefato executável ou instalável, os argumentos e o diretório de trabalho, sem compilar, instalar nem iniciar. O plano lista os nomes das chaves de substituição do ambiente e omite seus valores.

## Opções da CLI

| Opção | Significado |
| --- | --- |
| `--tags <tags>` | Tags adicionais de compilação Go, separadas por vírgulas. |
| `--skip-build` | Importe os executáveis de desktop existentes do diretório de saída do YAML e, em seguida, inicie sem compilar. |
| `--target <os>/<arch>` | Seleciona um destino; o padrão é o sistema operacional e a arquitetura do host. |
| `--plan` | Exibe as configurações resolvidas sem compilar nem iniciar. Planos de desktop multiplataforma são permitidos. |
| `--device <id>` | Número de série adb do Android ou identificador do simulador/dispositivo iOS. |
| `--emulator <name>` | Dispositivo virtual Android a iniciar ou reutilizar; não pode ser combinado com `--device`. |
| `--destination <kind>` | Destino do iOS: `simulator` (padrão) ou `device`. |
| `-- <args...>` | Substitui os argumentos configurados do aplicativo. Um `--` isolado os limpa. |
| `--appargs <string>` | Grafia de compatibilidade: analisa uma string usando regras de aspas no estilo POSIX. Não pode ser repetida nem combinada com argumentos após `--`. |

Seletores móveis exigem um manifesto e um destino móvel correspondente.

`--skip-build` requer um projeto YAML ativo e um executável de desktop compatível. Esse modo ignora o trabalho de frontend, os bindings, a compilação, o UPX e os hooks de build; os argumentos de inicialização e o ambiente configurados ainda se aplicam. Os executáveis importados são validados e copiados para o workspace gerado pelo Wake sem modificar os arquivos de origem. No macOS, os aplicativos com interface gráfica são montados com base nessas entradas e nos assets do pacote atual. Consulte [caminhos e requisitos das entradas pré-compiladas](/experimental/yaml-builds/#use-prebuilt-desktop-executables).
As tags de build não podem alterar um executável pré-compilado; portanto, `--tags` e as tags de build configuradas não têm efeito nesse modo.

## Campos

O mapeamento opcional `run` fornece padrões de inicialização compartilhados. Um mapeamento opcional `run` dentro de `targets["<os>"]` ou `targets["<os>/<arch>"]` fornece padrões específicos da plataforma ou da arquitetura. As configurações específicas da arquitetura têm precedência sobre as configurações gerais do sistema operacional, independentemente da ordem do mapa.

| Campo | Tipo | Padrão | Significado |
| --- | --- | --- | --- |
| `run.tags` | `list(string)` | `[]` | Tags Go adicionais usadas somente por compilações de execução. |
| `run.args` | `list(string)` | `[]` | Argumentos passados ao aplicativo. |
| `run.environment` | `map(string)` | `{}` | Substituições do ambiente de execução, separadas do ambiente do compilador. |
| `target["<os>"].run.tags` | `list(string)` | `[]` | Tags de execução adicionais para a plataforma selecionada. |
| `target["<os>"].run.args` | `list(string)` | herdado | Substitui os argumentos compartilhados quando definido explicitamente; `[]` os limpa. |
| `target["<os>"].run.environment` | `map(string)` | `{}` | Substitui por chave os valores compartilhados do ambiente de execução. |

Os argumentos são passados como um vetor de argumentos, sem expansão pelo shell. Os valores de ambiente são strings literais conforme o contrato existente do YAML, que aceita somente literais.

## Exemplo de plataforma

Este fragmento acompanha a configuração normal do projeto. `example_private_api` é um espaço reservado para a tag real de API privada exigida pelo exemplo.

```yaml
build:
  tags: [devtools]
run:
  args: [--example-mode]
targets:
  darwin:
    tags: [example_private_api]
    run:
      environment:
        EXAMPLE_WINDOW_MODE: native
  linux:
    run:
      environment:
        GDK_BACKEND: wayland
  windows:
    run:
      args: [--example-mode, --windows-option]
```

A configuração de ambiente do Linux é específica deste exemplo e exige uma sessão Wayland; ela não é um padrão do Wails. A tag de API privada é selecionada somente no macOS. O `target.tags` existente se aplica às compilações de produção, incluindo as compilações de execução. `target.run.tags` restringe uma tag adicional às compilações de execução.

Uma substituição específica da arquitetura pode refinar os padrões gerais do sistema operacional:

```yaml
targets:
  linux/arm64:
    run:
      tags: [example_arm64]
      args: []
```

Em conjunto com o mapeamento anterior do Linux, isso mantém `GDK_BACKEND`, adiciona a tag específica da arquitetura e limpa os argumentos compartilhados do aplicativo.

## Regras de resolução

| Configuração | Resolução, da menor para a maior precedência |
| --- | --- |
| Tags | Combina `build.tags`, o `target.tags` selecionado, o `production` automático, `run.tags`, o `target.run.tags` selecionado e o `--tags` da CLI; mantém a primeira ocorrência de cada tag. |
| Ambiente de execução | Herda o shell que inicia o processo, sobrepõe `run.environment` e depois sobrepõe o `target.run.environment` selecionado. Um valor vazio define uma string vazia. |
| Argumentos do aplicativo | O `run.args` compartilhado é substituído pelos argumentos de destino definidos explicitamente e, depois, pelos argumentos explícitos após `--`. Um `--` isolado limpa os argumentos configurados. |

No destino selecionado, as configurações gerais do sistema operacional são resolvidas antes das configurações específicas da arquitetura: as tags são acumuladas, os valores de ambiente são mesclados por chave e os argumentos fornecidos explicitamente substituem a lista menos específica. As compilações de execução também fornecem `android` ou `ios` para o destino móvel correspondente.

A geração de bindings e a compilação Go recebem as mesmas tags resolvidas. As chaves do cache de compilação incluem essas tags. Alterar apenas os argumentos do aplicativo ou o ambiente de execução não invalida a compilação; a próxima inicialização usa os novos valores.

Os mapeamentos de plataforma personalizam uma plataforma; sua presença não torna as demais plataformas incompatíveis. `project.supported_platforms` restringe as plataformas aceitas pelas compilações e por `run`; uma lista omitida ou vazia permite todas as plataformas. Por exemplo, `supported_platforms: [darwin]` dentro de `project` faz com que um exemplo exclusivo para macOS falhe antecipadamente no Linux.

## Comportamento de inicialização por plataforma

| Plataforma | Comportamento de inicialização |
| --- | --- |
| Linux e Windows | Inicia o executável com a entrada e a saída padrão conectadas, tratamento de encerramento e propagação do status de saída do aplicativo. |
| macOS | Monta um pacote de aplicativo e executa seu binário `Contents/MacOS` com a entrada e a saída padrão conectadas. Registra o pacote nos Serviços de Inicialização para esquemas de URL e associações de arquivos. A assinatura e a aceitação do ciclo de vida nativo exigem verificação no macOS. |
| Android | Compila um APK, seleciona um dispositivo ou emulador, instala e inicia o aplicativo. Relata separadamente as falhas de compilação, instalação e inicialização. |
| iOS | Compila um pacote de aplicativo, seleciona um simulador ou dispositivo, atende aos requisitos de assinatura e provisionamento, instala e inicia o aplicativo. As compilações para dispositivos exigem as credenciais adequadas e o conjunto de ferramentas da Apple. |

O Android rejeita argumentos do aplicativo e substituições do ambiente de execução antes da implantação. Quando também há suporte ao Android, as configurações compartilhadas exclusivas de desktop devem ficar nos mapeamentos de destino de desktop. As execuções no Android anexam os logs do aplicativo e interrompem o aplicativo quando o anexo termina. O iOS passa os argumentos para o inicializador selecionado. Os simuladores recebem substituições de ambiente por meio de variáveis `SIMCTL_CHILD_`; os dispositivos físicos recebem um mapa JSON por meio de `devicectl --environment-variables`. Em dispositivos móveis, o status de saída descreve o resultado do inicializador ou do monitoramento de logs, em vez de um código de saída portátil do aplicativo.

## Exemplos somente em Go e modo alternativo

Um exemplo somente em Go com ativos preexistentes ou incorporados desativa as etapas de frontend:

```yaml
frontend:
  disabled: true
```

Isso ignora a instalação das dependências do frontend, a geração de bindings e a compilação do frontend. Como alternativa, exemplos que precisam de classes de binding geradas podem definir `interfaces: false` em `frontend.bindings`.

O modo alternativo sem manifesto executa `go run [-tags <tags>] . [application arguments]` sem adicionar `production` nem `devtools`. Ele segue o comportamento de status de saída do próprio Go. Os argumentos após `--` são passados ao aplicativo, incluindo opções como `--help`. Os valores de tempo de execução exibidos por `--plan` limitam-se aos nomes das chaves de ambiente; os valores não são impressos. O plano do modo alternativo exibe o comando Go e o diretório de trabalho.

## Argumentos do aplicativo em desenvolvimento

Com um projeto YAML ativo, `wails3 dev` aceita argumentos do aplicativo independentemente dos padrões de produção de `run`. `dev.args` fornece padrões compartilhados; um mapeamento `dev` aninhado em `targets["<os>"]` ou `targets["<os>/<arch>"]` os substitui para esse destino. As configurações de arquitetura têm precedência sobre as configurações gerais do sistema operacional, independentemente da ordem do mapa. Um `args: []` explícito remove os argumentos herdados.

```yaml
dev:
  args: [--config-path, config/testing.yaml]
targets:
  windows:
    dev:
      args: [--config-path, config/testing-windows.yaml]
```

`run.args` e `run.tags` não são usados nas inicializações de desenvolvimento. A política de compilação `dev.tags` existente permanece inalterada. Os caminhos relativos do aplicativo são interpretados pelo aplicativo a partir do diretório do manifesto.

```shell
wails3 dev
wails3 dev -- --config-path "profile with spaces.yaml"
wails3 dev --appargs='--config-path "profile with spaces.yaml"'
wails3 dev --
wails3 dev --plan
```

Os argumentos após `--` substituem os argumentos configurados, inclusive quando a lista está vazia. Um `--help` após esse delimitador pertence ao aplicativo. `--appargs=""` também remove os padrões. A opção de compatibilidade aceita uma única string entre aspas; usá-la com `--`, com argumentos posicionais do aplicativo ou com um segundo `--appargs` é um erro. O uso incorreto de aspas é rejeitado antes do início de uma sessão.

O analisador de compatibilidade usa as mesmas regras de citação no estilo POSIX em todos os sistemas host. Ele não executa comandos de shell nem expande variáveis; o shell que fez a chamada pode já ter expandido o texto antes de passá-lo ao Wails. Os argumentos após `--` já são valores separados e não passam por tokenização adicional. Esse é o formato preferencial para caminhos do Windows, por exemplo, no PowerShell:

```powershell
wails3 dev -- --config-path 'C:\Users\Example User\testing.yaml'
```

A sessão preserva o vetor de argumentos resolvido durante as recompilações do código-fonte. Alterar apenas `dev.args` ou os argumentos do destino selecionado reinicia o backend usando o executável armazenado em cache; isso não reinicia o frontend. As substituições da CLI permanecem em vigor após o recarregamento do manifesto. Se uma substituição falhar, o executável anterior com seus argumentos anteriores é mantido ou restaurado. Os argumentos do aplicativo não são opções do compilador e não fazem parte das chaves do cache de compilação. `dev --plan` exibe os argumentos efetivos junto com o plano de compilação finito, sem iniciar o aplicativo.

Linux, Windows e macOS passam o vetor diretamente ao executável. O iOS o passa por meio do inicializador do simulador ou dispositivo. O Android rejeita argumentos de desenvolvimento não vazios do aplicativo antes da descoberta do dispositivo ou da compilação, pois não há um canal equivalente de argumentos compatível. Os padrões exclusivos de desktop podem ser colocados nos mapeamentos de destino de desktop ou removidos com `args: []` em `targets.android.dev`.

Essas opções de desenvolvimento exigem um `wails.yaml` ativo e a ativação explícita no YAML. Elas não reescrevem os comandos de shell do Taskfile nem do observador no fluxo de desenvolvimento legado.
