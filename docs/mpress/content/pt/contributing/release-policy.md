---
title: "Política de lançamentos"
description: "Como o Wails v3 versiona seus lançamentos nightly e semanais, como os mantenedores os preparam e por que os usuários devem fixar versões."
slug: "contributing/release-policy"
sourcePath: "contributing/release-policy.md"
---

## O modelo em uma tabela

O Wails v3 faz lançamentos com base em um **cronograma, não em funcionalidades**. Após a versão de referência `v3.0.0`:

| Lançamento | Tag | Periodicidade | Lançamento no GitHub | npm `@wailsio/runtime` |
| --- | --- | --- | --- | --- |
| Nightly | `v3.<next>.0-nightly.<YYYYMMDD>` (data UTC; `.N` é acrescentado se houver uma nova execução no mesmo dia) | automatizado, diariamente, somente quando `UNRELEASED_CHANGELOG.md` contém algo | pré-lançamento, **não** marcado como Latest | mesma versão sem `v`, dist-tag `nightly` |
| Semanal | `v3.<next>.0` | automatizado, domingo às 20:00 no fuso Australia/Sydney **(dia não confirmado)** | lançamento normal, marcado como Latest | mesma versão, dist-tag `latest` |
| Correção emergencial | `v3.<current>.<N>`, `N > 0` | raro, manual | versão normal, Latest | mesma versão, dist-tag `latest` |
| Principal | `v4.0.0` | somente para alterações incompatíveis | lançamento normal | `latest` |

`<next>` é a versão secundária mais alta **lançada** em v3, mais um. Uma versão nightly é uma prévia da **próxima versão semanal**; por isso, sua base é a tag da versão mais alta (pré-lançamentos são ignorados): após um hotfix `v3.1.1`, a próxima nightly ainda é `v3.2.0-nightly.*`. Ela precisa indicar uma versão superior porque um pré-lançamento é ordenado antes de sua própria versão final: `v3.0.0-nightly.<date>` seria ordenada antes de `v3.0.0`.

**O incremento da versão secundária não é um marco de funcionalidade.** `v3.4.0` é "o lançamento semanal posterior a `v3.3.0`". Para saber o que mudou, leia as notas de lançamento, não o número da versão.

A política é escolhida pelo **formato** da versão: um sufixo `-nightly.` indica um nightly; qualquer outro sufixo (`-alpha`, `-beta`, `-rc`) indica uma compilação legada anterior à GA; a ausência de sufixo indica um lançamento.

### Versões sincronizadas

A CLI, o módulo Go e o pacote npm `@wailsio/runtime` usam a **mesma versão**, obtida de `v3/internal/version/version.txt` (`3.2.0-nightly.20261102` no npm para a tag `v3.2.0-nightly.20261102`). Os projetos criados com `wails3 init` fixam `@wailsio/runtime` na própria versão da CLI; as compilações de desenvolvimento da CLI usam `latest` como alternativa.

## Exemplo prático

| Data | Evento | Tag | `go install ...@latest` é resolvido como |
| --- | --- | --- | --- |
| dom., 25 de out. | GA, lançamento manual | `v3.0.0` | `v3.0.0` |
| seg., 26 – sáb., 31 de out. | nightlies | `v3.1.0-nightly.20261026` ... `v3.1.0-nightly.20261031` | `v3.0.0` |
| dom., 1 de nov., às 20:00 no horário de Sydney | semanal | `v3.1.0` | `v3.1.0` |
| seg., 2 de nov. | nightly (prévia do próximo lançamento semanal) | `v3.2.0-nightly.20261102` | `v3.1.0` |
| qua., 4 de nov. | correção emergencial, manual, a partir de um branch de lançamento | `v3.1.1` | `v3.1.1` |
| dom., 8 de nov. | semanal | `v3.2.0` | `v3.2.0` |

Um lançamento semanal sempre fica acima de todos os nightlies com o mesmo número na ordenação, e um nightly nunca se torna `@latest`.

## O que os usuários executam

| Objetivo | Comando |
| --- | --- |
| Lançamento mais recente (semanal ou correção emergencial) | `go install github.com/wailsapp/wails/v3/cmd/wails3@latest` ou `wails3 update cli` |
| **Reproduzível (recomendado)** | fixe uma tag exata: `v3.1.0` em `go.mod`, na CI e na documentação |
| Um nightly | sua tag exata (`@v3.2.0-nightly.20261102`), `wails3 update cli -nightly` |
| Compilação mais recente de qualquer tipo | `wails3 update cli -pre` |
| Pacote npm do runtime | `@wailsio/runtime` (`latest` = semanal/correção emergencial, `nightly` = nightly) |

Por que isso funciona: o Go resolve `@latest` como a tag mais alta **sem um sufixo de pré-lançamento** e só considera pré-lançamentos quando não existe nenhum lançamento. A caixa de seleção "pre-release" do GitHub não é visível para o proxy do Go; portanto, os nightlies usam um sufixo de pré-lançamento semver em vez de depender dela. Sempre fixe a versão em produção.

## Branches

| Branch | Finalidade |
| --- | --- |
| `master` | Desenvolvimento e nightlies. |
| `releases/**` | Branches de curta duração para lançamentos avulsos: a GA `v3.0.0` e correções emergenciais. Exigem a mesma CI que `master`. |

Os lançamentos semanais são preparados a partir de `master`, mas **não da ponta desse branch**: o fluxo de trabalho seleciona o commit mais recente desde o lançamento semanal anterior em que todas as verificações obrigatórias (`v3 Build Results`, `Run Go Tests`, `Bundled runtime matches source`) foram aprovadas e aplica a tag exatamente a esse commit. Se nenhum estiver aprovado, nada será lançado. Os nightlies são preparados a partir de `master` e apenas emitem um aviso quando HEAD não está aprovado.

## Como preparar um lançamento

Os direitos de lançamento são restritos ao líder do projeto (consulte o documento de governança). Os colaboradores nunca criam tags nem publicam; eles adicionam entradas a `v3/UNRELEASED_CHANGELOG.md`.

Tudo é executado por meio de **Nightly + Weekly Release v3** (`.github/workflows/nightly-release-v3.yml`):

- **Nightly (automático)** às 15:00 UTC. É ignorado quando o changelog ainda não lançado está vazio.
- **Semanal (automática)** às 20:00, no horário de Sydney. O cron do GitHub aceita apenas UTC; por isso, dois crons (09:00 e 10:00 UTC) são agendados, e um job de verificação só prossegue quando são 20:00 em Sydney (considerando o horário de verão). A versão semanal passa pelo gate de lançamento, seleciona o commit aprovado mais recente, incrementa a versão secundária, cria uma tag nesse commit sem alterar o changelog e, em seguida, atualiza apenas os campos de versão em `master` (`version.txt` e as versões dos pacotes de runtime), preservando assim as atualizações de dependências e do arquivo de lock mescladas após o commit selecionado. O npm publica a partir da tag, não de `master`.
- **Manual.** Execute o fluxo de trabalho com `channel`, `dry_run` (ativado por padrão) e, opcionalmente, `release_version`. Uma simulação imprime a versão, o tipo, os sinalizadores do GitHub, a dist-tag do npm e as notas sem alterar o GitHub nem o npm.
- **GA ou hotfix.** Defina `release_version` (`v3.0.0`, `v3.1.1`). Um hotfix para uma linha mais antiga também pode definir `ref: releases/...`; o commit de lançamento correspondente é trazido de volta para `master` por cherry-pick.

Simulação local: `cd v3/tasks/release && go run release.go --dry-run --channel weekly`; `go run release.go --print-policy v3.2.0-nightly.20261102` mostra os sinalizadores derivados.

### Mecanismos de proteção

- Uma versão automática deve ser estritamente maior (na ordem semver, pré-lançamento < lançamento) que todas as tags existentes; caso contrário, a ferramenta se recusa a prosseguir.
- Apenas um trabalho de lançamento por vez (grupo de simultaneidade compartilhado).
- A publicação no npm é acionada pela **tag** do lançamento e compila o commit marcado.

## Controle de lançamento

O fluxo **semanal** (nunca o noturno) se recusa a criar uma tag enquanto existir qualquer um dos seguintes itens:

1. issues ou pull requests abertos com o rótulo `P0` ou `release-blocker`;
2. issues abertas com o rótulo `regression` criadas desde o último lançamento semanal;
3. alertas abertos do Dependabot de severidade crítica ou alta, alertas abertos de análise de código de severidade crítica ou alta e qualquer alerta aberto de verificação de segredos, **somente para código distribuído**. Alertas localizados exclusivamente em `v3/examples/`, `v2/examples/`, `website/`, `docs/`, `.github/` ou em qualquer diretório `examples/` são ignorados (`SHIPPED_ALERT_IGNORE` em `v3/tasks/releasegate/main.go`, **o proprietário deve revisar**);
4. avisos de segurança do repositório que ainda estejam em `triage`.

**Quem pode definir um bloqueador.** Um rótulo só é considerado se o evento `labeled` mais recente correspondente (sem um `unlabeled` posterior) tiver sido realizado por um membro da equipe `wailsapp/maintainers` ou por um colaborador com permissão `maintain`/`admin`. O mesmo rótulo aplicado por qualquer outra pessoa é ignorado e registrado, de modo que pessoas externas não possam bloquear nem desbloquear um lançamento. O rótulo `high-priority` aplicado automaticamente **não** é usado deliberadamente. Contas de bot, como `taliesin-ai`, passam pela lista de permissões se tiverem `maintain`; isso é intencional, e o agente de triagem pode aplicar rótulos de bloqueio com evidências.

**Comportamento em caso de falha.** O controle falha de forma **restritiva**: se não for possível ler uma API de alertas ou a lista de permissões, o fluxo semanal é interrompido. Em caso de falha, ele não cria nenhuma tag, grava a lista (com links) no resumo do trabalho e atualiza um comentário na única issue "Status do lançamento semanal" quando a variável de repositório `WEEKLY_STATUS_ISSUE` está definida.

**Configuração que o proprietário deve realizar no GitHub** (nada é criado automaticamente): crie os rótulos `P0`, `release-blocker` e `regression`; crie a issue "Status do lançamento semanal" e defina `WEEKLY_STATUS_ISSUE` com o número dela; conceda a `WAILS_REPO_TOKEN` acesso de leitura a Issues, Pull requests, Administration (permissão de colaborador), alertas do Dependabot, alertas de análise de código, alertas de verificação de segredos, avisos de segurança do repositório e membros da organização, além de acesso de gravação a Issues para o comentário de status. O `GITHUB_TOKEN` padrão não pode ler esses dados.

**Substituição.** `override_reason` (entrada de workflow-dispatch) ignora o controle. O motivo, o autor da ação e todos os bloqueios ignorados são registrados no resumo e na issue de status. Somente o líder de lançamento deve usá-la. O GitHub não pode restringir uma entrada a uma única pessoa; portanto, execute o trabalho de lançamento em um ambiente com revisores obrigatórios (recomendado), para que uma substituição exija aprovação explícita.

## Conteúdo relacionado

- [Testes e integração contínua](/contributing/testing-ci/)
- [Triagem para mantenedores](/contributing/maintainer-triage/)
