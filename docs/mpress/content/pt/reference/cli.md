---
title: "Visão geral dos comandos da CLI"
description: "Encontre comandos da CLI do Wails por finalidade: projetos, builds, assinatura, ferramentas e serviços"
slug: "reference/cli"
sourcePath: "reference/cli.md"
---

Comece aqui para encontrar um comando por finalidade. Para ver flags específicas de cada comando e exemplos de uso, consulte [Flags e exemplos da CLI](/guides/cli/).

## Visão geral

A CLI do Wails (`wails3`) é o ponto de entrada de linha de comando para criar, desenvolver, compilar, assinar, empacotar e inspecionar aplicações Wails 3. A maior parte da orquestração de builds é delegada aos Taskfiles de cada projeto (em `build/` no seu projeto) — muitos comandos `wails3` são wrappers simples que invocam uma tarefa específica.

Para obter a ajuda mais atualizada sobre qualquer comando, execute:

```bash
wails3 --help
wails3 <command> --help
```

## Ciclo de vida do projeto

| Comando | Descrição |
| --- | --- |
| `wails3 init` | Crie um novo projeto a partir de um modelo. Flags: `-n` (nome do projeto), `-t` (modelo; padrão: `vanilla`), `-p` (nome do pacote Go; padrão: `main`), `-d` (diretório do projeto; padrão: `.`), `-q` (modo silencioso), `-l` (listar modelos), `-mod` (caminho do módulo Go), `--git` (URL do repositório Git), `--skipgomodtidy`, `-s` (ignorar o aviso sobre modelo remoto), `--productname`/`--productdescription`/`--productversion`/`--productcompany`/`--productcopyright`/`--productcomments`/`--productidentifier`. |
| `wails3 dev` | Execute a aplicação no modo de desenvolvimento com recarregamento automático do frontend. Flags: `--config` (padrão: `./build/config.yml`), `--port` (porta de desenvolvimento do Vite), `-s` (habilitar HTTPS). |
| `wails3 build` | Compile o projeto. É um wrapper simples para a tarefa `build` do Taskfile. Flags: `--tags` (repassada como `EXTRA_TAGS=`), `--obfuscated` (compilar com Garble; consulte [Builds ofuscados](/guides/build/obfuscation/)), `--garbleargs` (flags adicionais repassadas a `garble` antes do subcomando `build`). |
| `wails3 package` | Execute a tarefa `package` do Taskfile específica da plataforma. |
| `wails3 task [name]` | Execute qualquer tarefa do Taskfile; quando nenhum nome é informado, `--list` mostra todas as tarefas registradas. |
| `wails3 mcp` | Inicie o servidor MCP do projeto. Usa automaticamente stdio para processos iniciados por agentes ou HTTP Streamable no endereço de loopback para uso interativo no terminal. |
| `wails3 doctor` | Exiba um relatório de diagnóstico do seu ambiente. |
| `wails3 doctor-ng` | Variante TUI mais recente de `doctor`. |
| `wails3 version` | Exiba a versão da CLI. |
| `wails3 releasenotes` | Exiba as notas de versão recentes. |
| `wails3 docs` | Abra o site da documentação no navegador. |
| `wails3 sponsor` | Abra a página de patrocínio. |

## Geração

`wails3 generate <subcommand>`:

| Subcomando | Descrição |
| --- | --- |
| `generate bindings` | Gere bindings do Go para o frontend. Flags: `-d` (diretório de saída), `-models`, `-index`, `-ts`, `-i` (interfaces), `-b` (bundle), `-names` (emitir `Call.ByName`), `-noevents`, `-noindex`, `-dry`, `-silent`, `-v`, `-clean` (padrão: `true`), `-f`, `-obfuscated` (gerar `wails_obfuscated.gen.go` com IDs de binding estáveis para builds com Garble; consulte [Builds ofuscados](/guides/build/obfuscation/)), `-obfuscated-output` (diretório do arquivo gerado; o padrão é o diretório do pacote principal). Aceita padrões de pacotes (por exemplo, `./...`); se nenhum for informado, usa o diretório atual. |
| `generate icons` | Converta um PNG de origem nos formatos de ícone da plataforma. Flags: `-input`, `-windowsfilename`, `-macfilename`, `-iconcomposerinput`, `-macassetdir`. |
| `generate build-assets` | Gere o conteúdo do diretório `build/` (trechos de Taskfile, arquivos NSIS, `Info.plist`, modelo `.desktop` etc.) a partir de `build/config.yml`. |
| `generate runtime` | Gere novamente o `/wails/runtime.js` pré-compilado fornecido ao webview. |
| `generate syso` | Gere o arquivo de recursos `.syso` do Windows (ícone + manifesto + informações de versão). |
| `generate webview2bootstrapper` | Gere um instalador de bootstrap do WebView2 para Windows. |
| `generate constants` | Gere constantes de nomes de eventos em JS a partir dos tipos de eventos do Go. |
| `generate template` | Crie a estrutura inicial de um novo modelo de projeto. |
| `generate .desktop` | Gere um arquivo `.desktop` do Linux (usado por AppImage/DEB/RPM). |
| `generate appimage` | Gere o diretório de build do AppImage. |

## Atualização

`wails3 update <subcommand>`:

| Subcomando | Descrição |
| --- | --- |
| `update build-assets` | Atualize o diretório `build/` a partir de `build/config.yml` (preserva as edições do usuário quando possível). |
| `update cli` | Atualize o próprio binário `wails3`. |

## Assinatura de código e empacotamento

| Comando | Descrição |
| --- | --- |
| `wails3 signing` | Assistente interativo que configura a assinatura para as plataformas detectadas em `build/`. Flags: `--platform` (pode ser repetida; por padrão, detecta automaticamente com base no diretório de build). |
| `wails3 entitlements` | Assistente interativo para entitlements do macOS. Flag: `--output` (caminho; padrão: `build/darwin/entitlements.plist`). |
| `wails3 sign [GOOS=…]` | Wrapper que executa a tarefa `*:sign` do Taskfile específica da plataforma para o sistema operacional atual (ou para aquele especificado por `GOOS`). |
| `wails3 tool sign` | Ponto de entrada de baixo nível para assinatura direta. Flags: `--input`, `--output`, `--verbose`, `--certificate`, `--password`, `--thumbprint`, `--timestamp`, `--identity`, `--entitlements`, `--hardened-runtime`, `--notarize`, `--keychain-profile`, `--pgp-key`, `--pgp-password`, `--role`. |

Nesta versão, o analisador da CLI expõe `signing` e `entitlements` no nível superior. Consulte `wails3 signing --help` e `wails3 entitlements --help`; `wails3 setup signing` não seleciona o assistente de assinatura. Para gerenciar credenciais diretamente, use `xcrun notarytool store-credentials` ou `gpg`.

## Ferramentas

`wails3 tool <subcommand>`:

| Subcomando | Descrição |
| --- | --- |
| `tool checkport` | Verifica se uma porta TCP está aberta (útil para aguardar o Vite). |
| `tool watcher` | Executa um comando sempre que os arquivos monitorados são alterados. |
| `tool cp` | Copia arquivos entre plataformas. |
| `tool buildinfo` | Exibe as informações de compilação do Go incorporadas em um binário. |
| `tool package` | Cria um pacote Linux (`deb`, `rpm`, `archlinux`) a partir de `build/linux/nfpm`. |
| `tool version` | Incrementa a versão semântica de um projeto. |
| `tool lipo` | Combina binários de várias arquiteturas do macOS em um binário universal. |
| `tool capabilities` | Verifica no sistema a disponibilidade de GTK3/GTK4 e WebKit. |
| `tool sign` | (Consulte [Assinatura de código e empacotamento](#assinatura-de-cdigo-e-empacotamento).) |

## Serviços

`wails3 service <subcommand>`:

| Subcomando | Descrição |
| --- | --- |
| `service init` | Gera a estrutura inicial de um novo pacote de serviço. |

## iOS

`wails3 ios <subcommand>`:

| Subcomando | Descrição |
| --- | --- |
| `ios overlay:gen` | Gera o overlay do Go para o shim da ponte do iOS. |
| `ios xcode:gen` | Gera um projeto do Xcode no diretório de saída. |

## Caminhos de saída da compilação

- Os binários nativos são gravados em `bin/<APP_NAME>` (ou em `bin/<APP_NAME>.exe` no Windows). Não há `build/bin/`.
- As saídas empacotadas (`.app`, `.dmg`, instalador NSIS, MSIX, DEB/RPM/AppImage) também são gravadas em `bin/` (ou em subdiretórios específicos da plataforma criados pela tarefa correspondente do Taskfile).

## Flags compartilhadas

| Flag | Aplicável a | Descrição |
| --- | --- | --- |
| `-nocolour` | Comandos que a expõem em `--help` (por exemplo, `init` e `update cli`) | Desativa as cores ANSI na saída da CLI. Nem todos os comandos aceitam essa flag, incluindo `version`. |

---

**Dúvidas?** Pergunte no [Discord](https://discord.gg/JDdSxwjhGf) ou consulte os [exemplos](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples).
