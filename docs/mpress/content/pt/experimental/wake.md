---
title: "Wake"
description: "Um executor de builds experimental e integrado ao Wails que executa seus Taskfiles existentes com builds incrementais mais rápidos, saída estruturada e execução paralela por padrão."
slug: "experimental/wake"
sourcePath: "experimental/wake.md"
---

@note{type="caution" title="Recurso experimental"}
O Wake precisa ser habilitado por meio de `WAILS_USE_WAKE=true` e **não** é o executor padrão. Quando a variável não está definida, os comandos `wails3 build / package / sign / task` se comportam exatamente como antes. A cobertura de recursos e o comportamento podem mudar entre versões.

@end

O Wake é um **executor de builds alternativo e experimental** para `wails3`. Ele lê o mesmo `Taskfile.yml` que seu projeto já tem — com a mesma sintaxe de tarefas, dependências, variáveis, modelos, inclusões e namespaces de plataforma — e o executa por meio de um executor integrado ao Wails, em vez de usar o ambiente de execução de uso geral do [Task](https://taskfile.dev).

O objetivo não é substituir o Task, mas oferecer um executor criado especificamente para a forma como os projetos Wails são realmente compilados, com semântica, saída e padrões compatíveis com o restante da CLI `wails3`. **Se você usar apenas o Wake, seus Taskfiles não precisarão ser alterados.**

## Por que ele existe

Tanto o Wake quanto o ambiente de execução do Task são compilados no `wails3` — nenhum deles exige a instalação de um binário separado. A diferença é que o Wake **entende o domínio**. Um executor de uso geral executa as etapas listadas por um Taskfile na ordem determinada. O Wake sabe o que um build do Wails realmente *é* — o pacote do frontend é incorporado ao binário, o binário é empacotado em artefatos específicos de cada plataforma, e ícones e bindings também são gerados nesse processo — e usa esse conhecimento para otimizar o build de maneiras que um executor genérico não consegue.

- **Ele executa apenas o trabalho realmente necessário para o build.** O Wake rastreia por conta própria as entradas e saídas reais de cada etapa. Para um build do Go, isso inclui o grafo de módulos e as saídas das etapas das quais ele depende. Assim, quando nada relevante mudou, ele ignora completamente o compilador e o vinculador, em vez de executá-los novamente. Um executor de uso geral só consegue ignorar uma etapa quando o Taskfile especifica antecipadamente e com exatidão quais arquivos devem ser monitorados; o Wake deduz isso com base no que já sabe sobre o build. Em uma recompilação sem alterações, são aproximadamente **~20 ms (Wake) contra ~316 ms (Task)**. Builds a frio levam o mesmo tempo total — são dominados por `npm install`, pelo Vite e pelo compilador Go.

- **Ele sabe o que pode ser executado simultaneamente.** Como o Wake entende quais etapas são independentes, ele as executa em paralelo por padrão, e a linha de resultado informa a aceleração obtida. Desative esse comportamento com `WAKE_SERIAL=true` quando a saída intercalada de etapas irmãs puder atrapalhar uma investigação.

- **Saída estruturada controlada pelo wails3.** O Wake renderiza a saída usando o próprio gerador de relatórios do wails3: uma linha por etapa planejada, status em tempo real, um detalhamento das fases por cores ao final e links `file:line` clicáveis nos painéis de falha. `NO_COLOR` e ambientes sem TTY (logs de CI) são tratados adequadamente com uma apresentação simplificada.

- **Integrado, para que possa evoluir com o Wails.** Como o Wake faz parte do `wails3`, em vez de ser uma ferramenta de terceiros, novos recursos de build podem ser adicionados diretamente — sem esperar que outro projeto os implemente. Isso também possibilita executar scripts e ferramentas multiplataforma de forma nativa nos casos em que, atualmente, um Taskfile invoca o binário `wails3` por meio do shell, criando um processo a cada chamada. Incorporar esse trabalho ao processo reduz a sobrecarga e possibilitará novas melhorias de desempenho.

## Como habilitar o Wake

O Wake é controlado exclusivamente pela variável de ambiente `WAILS_USE_WAKE=true`. Quando ela não está definida (ou contém qualquer valor diferente de `true`), todos os comandos `wails3` usam o ambiente de execução integrado do Task exatamente como antes.

```bash
# Default: Task runtime, no Wake involvement
wails3 build

# Opt in: Wake drives build / package / sign / task <name>
WAILS_USE_WAKE=true wails3 build
WAILS_USE_WAKE=true wails3 package
WAILS_USE_WAKE=true wails3 task <some-task-name>
```

A opção abrange `wails3 build`, `wails3 package`, `wails3 sign` e `wails3 task <name>`. `wails3 dev` **ainda não** é afetado — o monitor de desenvolvimento continua usando seu próprio pipeline.

@note{type="tip" title="É seguro habilitar o Wake"}
Se o Wake encontrar um recurso do Taskfile que não implementa, ele delegará toda a execução ao ambiente de execução integrado do Task, dentro do mesmo processo — não há nenhum binário `task` externo para instalar. Na pior das hipóteses, você obterá exatamente o mesmo comportamento que teria sem a opção.

@end

## Substituições locais em camadas

O Wake permite usar um **Taskfile base com substituições locais**. Coloque um arquivo ao lado do seu `Taskfile.yml`, e as definições desse arquivo terão precedência:

| Arquivo | Finalidade | Precedência |
| --- | --- | --- |
| `Taskfile.yml` | base, versionado | mais baixa |
| `Taskfile.override.yml` / `.yaml` | substituições versionadas para toda a equipe | intermediária |
| `Taskfile.local.yml` / `.yaml` | pessoal, normalmente ignorado pelo Git | mais alta |

**Semântica da mesclagem (o arquivo local prevalece):**

- Uma tarefa com o **mesmo nome** substitui a tarefa base. Quando fornecidos pela substituição, os campos de lista (`cmds`, `deps`, `sources`, `generates`, `platforms`, `status`, `preconditions`, `aliases`) **substituem** os campos da base; os campos omitidos pela substituição são mantidos da base.
- `env` e `vars` são **mesclados por chave**, e a substituição prevalece em caso de conflito.
- Uma tarefa que existe **somente** em um arquivo de substituição é **adicionada**.

Por exemplo, se o `Taskfile.yml` versionado compila com opções de desenvolvimento, mas sua máquina deve sempre gerar builds de produção:

```yaml
# Taskfile.local.yml (git-ignored, yours)
tasks:
  build:
    cmds:
      - go build -tags production -o bin/app .
  smoke:
    cmds:
      - ./bin/app --selftest
```

Agora, `build` executa seu comando de produção, e `smoke` fica disponível, sem alterar o Taskfile versionado.

@note{type="note" title="Modelo de confiança"}
Os arquivos de substituição são descobertos e aplicados automaticamente, sem confirmação. Isso não concede nenhum recurso novo — um Taskfile já executa comandos de shell arbitrários, portanto uma substituição não pode fazer nada que não pudesse ser feito editando `Taskfile.yml`. O arquivo versionado `Taskfile.override.*` aparece nos diffs dos PRs; `Taskfile.local.*` é criado na sua própria máquina. Uma substituição malformada interrompe a execução, em vez de ser ignorada silenciosamente. Defina `WAILS_NO_OVERRIDES=true` para desativar completamente a descoberta de substituições em builds de CI determinísticos.

@end

## Fallback automático

Se o Wake encontrar um recurso do Taskfile que não implementa, ele delegará toda a execução ao ambiente de execução integrado do Task. Atualmente, os seguintes recursos acionam esse fallback:

- `dotenv` no nível do Taskfile
- modos `output` diferentes de `interleaved`
- um bloco `requires`
- `interval` (no nível do Taskfile ou da tarefa)
- modos `run` diferentes de `always`
- `short` em uma tarefa
- `defer` em uma tarefa

## Variáveis de ambiente

| Variável | Efeito |
| --- | --- |
| `WAILS_USE_WAKE` | `true` habilita o Wake para os verbos `wails3` que podem ser roteados; qualquer outro valor usa o runtime do Task |
| `WAILS_NO_OVERRIDES` | `true` ignora a descoberta de `Taskfile.local.*` / `.override.*` (builds determinísticos) |
| `WAKE_VERBOSE` | Transmite stdout/stderr dos subprocessos em tempo real, em vez de capturá-los para exibição somente em caso de falha |
| `WAKE_SILENT` | Suprime completamente a saída das tarefas |
| `WAKE_SERIAL` | `true` desabilita a distribuição paralela de `deps:` (a execução paralela é o padrão) |
| `WAKE_FORCE` | `true` ignora todos os caches para executar uma recompilação realmente limpa |
| `WAKE_DEBUG` | Registra detalhes internos do resolvedor (DAG, dependências, referências de variáveis, roteamento de execução) |
| `WAKE_NOTICE` | `off` para silenciar o aviso "wake (experimental)" a cada execução |

O cache de build fica em `.wake/cache.json` (o Task usa `.task/`).

## Feedback

O Wake é um experimento, e seu feedback determina o rumo dele. Se você testá-lo, gostaríamos de saber se ele foi mais rápido e mais claro e se algo deixou de funcionar — os relatos mais úteis informam o que você executou, o que esperava e o que realmente aconteceu. Conte para nós na [discussão sobre feedback do Wake](https://github.com/wailsapp/wails/discussions/5679).
