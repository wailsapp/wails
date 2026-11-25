---
title: "Recursos experimentais"
description: "Um espaço para experimentos em desenvolvimento no Wails v3 — o que são, por que existem e como enviar feedback."
slug: "experimental"
sourcePath: "experimental/index.md"
---

@note{type="caution" title="Área de experimentos"}
Tudo nesta seção é, por definição, um experimento. Os recursos apresentados aqui são opcionais, ficam desativados por padrão e podem mudar de formato, ser renomeados ou ser totalmente removidos entre versões. Não faça nada essencial depender deles sem estar preparado para mudanças frequentes.

@end

## O que significa "Experimental"

O Wails disponibiliza experimentos publicamente. Um experimento é uma ideia que consideramos promissora o bastante para colocar antecipadamente em suas mãos — mas com a qual ainda não nos comprometemos por completo. Estamos publicando-a *porque* queremos aprender com o uso real antes de decidir se ela se tornará uma parte permanente e oficialmente compatível do Wails.

Isso significa que algumas condições se aplicam a tudo nesta seção:

- **É opcional.** Os experimentos nunca alteram o comportamento padrão de `wails3`. Você os ativa deliberadamente (geralmente com uma variável de ambiente ou flag de compilação) e, quando estão desativados, nada muda no seu fluxo de trabalho existente.
- **Pode não sobreviver.** Alguns experimentos tornam-se recursos estáveis. Outros são reformulados a ponto de ficarem irreconhecíveis ou são descartados. Preferimos experimentar em público e aprender rapidamente a lançar apenas aquilo sobre o qual já temos certeza.
- **A API não está consolidada.** Nomes, flags, valores padrão e comportamentos podem mudar entre versões enquanto o experimento toma forma. As notas de versão destacarão as mudanças, mas não espere as garantias de estabilidade oferecidas pelos recursos estáveis.

## Queremos seu feedback

Esta é a parte importante. Os experimentos prosperam ou desaparecem de acordo com o retorno que recebemos das pessoas que realmente os utilizam. Se você testar um deles, queremos sinceramente saber:

- Funcionou no seu projeto? Em que pontos falhou?
- Foi mais rápido, mais claro ou mais agradável — ou a mudança não valeu a pena?
- Quais condições precisariam ser atendidas para você usá-lo por padrão?

O feedback mais útil é concreto: o que você executou, o que esperava e o que realmente aconteceu. Cada experimento tem sua própria discussão na categoria **Experimentos** do GitHub Discussions:

@container{display="grid" columns="2" gap="1rem"}
@linkcard{title="Discussões sobre experimentos" href="https://github.com/wailsapp/wails/discussions/categories/experiments" description="Encontre a discussão do experimento que você está usando e conte como foi, o que não funcionou ou o que está faltando."}
@end

## Experimentos atuais

@container{display="grid" columns="2" gap="1rem"}
@linkcard{title="Wake" href="/experimental/wake/" description="Experimente o sistema de compilação YAML nativo: crie e configure um projeto em um único assistente, compile e empacote o projeto ou analise uma migração unidirecional de Taskfiles."}
@linkcard{title="Compilações sem CGO" href="/experimental/#cgo-free-builds-purego" description="Compile aplicativos para macOS e Linux com CGO_ENABLED=0: sem toolchain C e com compilação cruzada sem Docker. Ativado com WAILS_PUREGO=1 ou com a tag de compilação purego."}
@linkcard{title="Controle por LLM (MCP)" href="/guides/mcp-service/" description="Um servidor integrado do Model Context Protocol que permite que agentes de LLM inspecionem, testem e controlem um aplicativo Wails em execução — sem exigir código do usuário e ativado com uma tag de compilação."}
@end

## Compilações sem CGO (purego)

@note{type="caution" title="Recurso experimental"}
O backend sem CGO é opcional, fica protegido pela tag de compilação `purego` (ou por `WAILS_PUREGO=1`) e **não** é o padrão. Sem ele, as compilações para macOS e Linux usam o backend cgo exatamente como antes. O comportamento pode mudar entre versões.

@end

### O que é

No macOS e no Linux, o Wails normalmente se comunica com as bibliotecas de interface do sistema (AppKit e WKWebView no macOS; GTK e WebKitGTK no Linux) por meio de cgo. O backend sem CGO faz as mesmas chamadas a partir de Go puro usando [purego](https://github.com/ebitengine/purego), que carrega as bibliotecas do sistema em tempo de execução.

A API pública do aplicativo permanece inalterada. Selecione o backend durante a compilação; não é necessário alterar o código do aplicativo. Teste seu aplicativo no sistema de destino antes de distribuí-lo.

O que você ganha:

- **Nenhuma toolchain C para a compilação Go.** A compilação Go não precisa das ferramentas de linha de comando do Xcode, de `gcc` nem dos cabeçalhos de desenvolvimento do Linux. As compilações do frontend ainda precisam das ferramentas JavaScript habituais.
- **Compilação cruzada sem Docker.** Compile um binário para macOS ou Linux a partir de qualquer sistema operacional usando apenas `go build`.

O Windows já compila sem cgo e não é afetado.

### Antes de ativar

Use versões da CLI e do framework que incluam este experimento. Atualizar apenas a CLI não adiciona o backend a uma dependência mais antiga do framework no `go.mod` do seu aplicativo. Para testar a versão de desenvolvimento, execute estes comandos no diretório do aplicativo:

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@master
go get github.com/wailsapp/wails/v3@master
```

Ambos os comandos acompanham a branch dinâmica `master`. Para obter compilações reproduzíveis, use o mesmo commit testado para a CLI e o framework. Siga o [guia de instalação](/getting-started/installation/) para configurar a toolchain Go e a CLI.

Os backends experimentais para macOS e Linux são compatíveis com `amd64` e `arm64`. A compilação cruzada cria um executável para o destino solicitado; ela não executa esse arquivo na máquina de compilação. Teste-o no sistema operacional e na arquitetura de destino.

### Como usar

Escolha uma destas formas de selecionar o backend.

#### Com a CLI do Wails

Em um projeto gerado com esta versão da CLI, `build`, `package` e `dev` usam o backend experimental quando `WAILS_PUREGO` está ativado. Em um shell POSIX:

```sh
WAILS_PUREGO=1 wails3 dev
WAILS_PUREGO=1 wails3 build
WAILS_PUREGO=1 wails3 package
```

No PowerShell:

```powershell
$env:WAILS_PUREGO = "1"
wails3 build
```

A CLI reconhece `1`, `true`, `on` e `yes`, sem diferenciar maiúsculas de minúsculas e ignorando espaços em branco ao redor. Um valor vazio não ativa este experimento.

Como alternativa, passe a tag para um único comando de compilação ou empacotamento:

```sh
wails3 build -tags purego
wails3 package -tags purego
```

Com a tag definida, os Taskfiles padrão atualizados compilam com `CGO_ENABLED=0` usando a toolchain Go do host. Por exemplo, compile um executável Linux ARM64 no macOS ou no Windows:

```sh
wails3 build -tags purego GOOS=linux ARCH=arm64
```

Projetos mais antigos mantêm os Taskfiles existentes. Compare as tarefas de compilação personalizadas com os modelos atualizados ou use apenas Go para testar o backend sem alterar essas tarefas. A variável de ambiente seleciona a tag nos comandos da CLI do Wails; `go build` usado diretamente exige que a tag seja informada explicitamente.

#### Diretamente com Go

Primeiro compile o frontend para que todos os recursos referenciados por `go:embed` existam; depois, compile o aplicativo:

```sh
CGO_ENABLED=0 go build -tags purego,production .
```

`CGO_ENABLED=1 go build -tags purego` também funciona: o backend é o mesmo; cgo apenas permanece disponível para o restante do programa.

Para uma compilação cruzada, defina também `GOOS` e `GOARCH`:

```sh
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -tags purego,production -o myapp .
```

No Linux, selecione a stack GTK3 com `-tags purego,gtk3,production`.

Para voltar ao backend padrão, remova a tag `purego` e desfaça a definição de `WAILS_PUREGO` (`unset WAILS_PUREGO` em um shell POSIX ou `Remove-Item Env:WAILS_PUREGO` no PowerShell). Ao voltar para o backend cgo do macOS ou Linux, remova também qualquer definição explícita de `CGO_ENABLED=0`.

Se você compilar com `CGO_ENABLED=0`, mas esquecer a tag, a compilação será interrompida com um único erro que informa a correção:

```
undefined: requires_CGO_ENABLED_1_or_build_tag_purego
```

### Requisitos de tempo de execução

O aplicativo carrega as bibliotecas do sistema ao iniciar, portanto elas precisam estar instaladas na máquina de destino, assim como em uma compilação com cgo:

| Plataforma | Bibliotecas | Versão mais antiga compatível |
|---|---|---|
| macOS | AppKit, WebKit (parte do sistema operacional) | Igual ao backend com cgo |
| Linux (padrão) | GTK 4, WebKitGTK 6.0 | GTK 4.8, WebKitGTK 2.40 (Debian 12) |
| Linux (`gtk3`) | GTK 3, WebKit2GTK 4.1 | GTK 3.24, WebKit2GTK 2.36 (Ubuntu 22.04) |

As funções posteriores a essas versões mínimas são localizadas quando o aplicativo é iniciado e usadas somente quando estão disponíveis na biblioteca instalada. Assim, o mesmo binário é executado tanto em sistemas antigos quanto em novos.

O aplicativo ainda usa bibliotecas de interface nativas e pode vincular dinamicamente bibliotecas do sistema, como a libc do Linux. `CGO_ENABLED=0` não significa que o executável seja completamente estático. Quando uma biblioteca de runtime obrigatória estiver ausente, instale no sistema de destino o conjunto GTK/WebKit correspondente; instalar cabeçalhos C na máquina de compilação não corrige a ausência de uma biblioteca de runtime.

### Empacotamento e assinatura

O backend elimina da compilação a exigência de uma cadeia de ferramentas C. Cada formato de pacote mantém seus próprios requisitos. Em particular, selecionar `purego` não torna multiplataforma as ferramentas existentes de implantação de AppImage. Consulte [Empacotamento para Linux](/guides/build/linux/).

Aplicativos macOS compilados de forma cruzada ainda precisam da assinatura e notarização adequadas antes da distribuição. Os instaladores do Windows ainda precisam de suas ferramentas de criação de instaladores. Consulte [Empacotamento para macOS](/guides/build/macos/) e [Empacotamento para Windows](/guides/build/windows/). Para dependências que usam cgo, `CGO_ENABLED=1` também exige um compilador C e todos os cabeçalhos ou bibliotecas necessários para o destino; selecionar o backend do Wails em Go puro não elimina os requisitos dessas dependências.

### Validação

As verificações de bindings gerados e `puregocheck` inspecionam símbolos nativos, proteções para versões opcionais e a cobertura dos backends. O fluxo de trabalho de paridade do Linux testa GTK4 e GTK3 com ambos os backends, incluindo solicitações de recursos, corpos binários de solicitações e despacho nativo; ele também faz compilação cruzada para destinos Linux ARM64. As verificações de compilação cruzada confirmam que o código é compilado, enquanto as verificações de runtime exercitam as bibliotecas nativas instaladas na máquina de teste.

O [guia de testes de backend](https://github.com/wailsapp/wails/blob/master/v3/TESTING.md) descreve as sondagens de paridade e os testes de entrada de ponta a ponta. Ao relatar um problema, diferencie uma falha de compilação de uma falha de runtime ou de empacotamento.

### Diferenças conhecidas

- **Os atalhos de menu no Linux usam registro no nível do aplicativo.** Evite atribuir atalhos conflitantes a menus diferentes. Remover explicitamente o menu de uma janela com `SetMenu(nil)` também impede que os atalhos do menu removido sejam ativados nessa janela; as outras janelas mantêm seus menus e atalhos.
- **Os painéis para abrir e salvar do macOS** são executados em um processo separado do sistema, como ocorre com cgo; os testes verificam se eles abrem e são cancelados corretamente.

Encontrou algo mais? Conte para nós.

### Feedback

Conte como foi sua experiência no tópico **Compilações sem CGO**, em [Experimentos](https://github.com/wailsapp/wails/discussions/categories/experiments): informe seu sistema operacional e a versão, como você fez a compilação e o que funcionou ou não. O relato mais útil que podemos receber é um caso de falha acompanhado da saída de `wails3 doctor`.
