---
title: "Referência da CLI"
description: "Referência completa dos comandos da CLI do Wails"
slug: "reference/cli"
sourcePath: "reference/cli.md"
---

## Visão geral

A CLI do Wails (`wails3`) é o ponto de entrada de linha de comando para criar, desenvolver, compilar, assinar, empacotar e inspecionar aplicações Wails 3. A maior parte da orquestração da compilação é delegada aos Taskfiles de cada projeto (em `build/` no seu projeto) — muitos comandos `wails3` são wrappers simples que invocam uma tarefa específica.

Para obter a ajuda mais atualizada sobre qualquer comando, execute:

```bash
wails3 --help
wails3 <command> --help
```

## Ciclo de vida do projeto

| Comando | Descrição |
| --- | --- |
| `wails3 init` | Cria um novo projeto com base em um modelo. Opções: `-n` (nome do projeto), `-t` (modelo; padrão: `vanilla`), `-p` (nome do pacote Go; padrão: `main`), `-d` (diretório do projeto; padrão: `.`), `-q` (modo silencioso), `-l` (listar modelos), `-mod` (caminho do módulo Go), `--git` (URL do repositório Git), `--skipgomodtidy`, `-s` (ignorar o aviso sobre modelo remoto), `--productname`/`--productdescription`/`--productversion`/`--productcompany`/`--productcopyright`/`--productcomments`/`--productidentifier`. |
| `wails3 dev` | Executa a aplicação no modo de desenvolvimento, com recarregamento automático do frontend. Opções: `--config` (padrão: `./build/config.yml`), `--port` (porta de desenvolvimento do Vite), `-s` (habilitar HTTPS). |
| `wails3 build` | Compila o projeto. É um wrapper simples em torno da tarefa `build` do Taskfile. Opções: `--tags` (repassada como `EXTRA_TAGS=`), `--obfuscated` (compilar com Garble; consulte [Compilações ofuscadas](/guides/build/obfuscation/)), `--garbleargs` (opções adicionais repassadas a `garble` antes do subcomando `build`). |
| `wails3 package` | Executa a tarefa `package` do Taskfile específica da plataforma. |
| `wails3 task [name]` | Executa qualquer tarefa do Taskfile; quando nenhum nome é informado, `--list` mostra todas as tarefas registradas. |
| `wails3 mcp` | Inicia o servidor MCP do projeto. Usa automaticamente stdio para processos iniciados por agentes ou HTTP Streamable no endereço de loopback para uso interativo no terminal. |
| `wails3 doctor` | Exibe um relatório de diagnóstico do seu ambiente. |
| `wails3 doctor-ng` | Variante TUI mais recente de `doctor`. |
| `wails3 version` | Exibe a versão da CLI. |
| `wails3 releasenotes` | Exibe as notas das versões recentes. |
| `wails3 docs` | Abre o site da documentação no navegador. |
| `wails3 sponsor` | Abre a página de patrocínio. |

## Geração

`wails3 generate <subcommand>`:

| Subcomando | Descrição |
| --- | --- |
| `generate bindings` | Gera bindings de Go para o frontend. Opções: `-d` (diretório de saída), `-models`, `-index`, `-ts`, `-i` (interfaces), `-b` (bundle), `-names` (emitir `Call.ByName`), `-noevents`, `-noindex`, `-dry`, `-silent`, `-v`, `-clean` (padrão: `true`), `-f`, `-obfuscated` (gerar `wails_obfuscated.gen.go` com IDs de binding estáveis para compilações com Garble; consulte [Compilações ofuscadas](/guides/build/obfuscation/)), `-obfuscated-output` (diretório do arquivo gerado; o padrão é o diretório do pacote principal). Aceita padrões de pacotes (por exemplo, `./...`); se nenhum for informado, usa o diretório atual. |
| `generate icons` | Converte um PNG de origem nos formatos de ícone da plataforma. Opções: `-input`, `-windowsfilename`, `-macfilename`, `-iconcomposerinput`, `-macassetdir`. |
| `generate build-assets` | Gera o conteúdo do diretório `build/` (trechos de Taskfile, arquivos NSIS, `Info.plist`, modelo `.desktop` etc.) com base em `build/config.yml`. |
| `generate runtime` | Gera novamente o `/wails/runtime.js` pré-compilado fornecido ao webview. |
| `generate syso` | Gera o arquivo de recursos `.syso` do Windows (ícone + manifesto + informações da versão). |
| `generate webview2bootstrapper` | Gera um instalador de inicialização do WebView2 para Windows. |
| `generate constants` | Gera constantes JS de nomes de eventos com base nos tipos de eventos Go. |
| `generate template` | Cria a estrutura inicial de um novo modelo de projeto. |
| `generate .desktop` | Gera um arquivo `.desktop` para Linux (usado por AppImage/DEB/RPM). |
| `generate appimage` | Gera o diretório de compilação do AppImage. |

## Atualização

`wails3 update <subcommand>`:

| Subcomando | Descrição |
| --- | --- |
| `update build-assets` | Atualiza o diretório `build/` com base em `build/config.yml`, preservando as alterações do usuário quando possível. |
| `update cli` | Atualiza o próprio binário `wails3`. |

## Assinatura de código e empacotamento

| Comando | Descrição |
| --- | --- |
| `wails3 setup signing` | Assistente interativo que configura a assinatura para as plataformas detectadas em `build/`. Opções: `--platform` (pode ser repetida; por padrão, detecta automaticamente com base no diretório de compilação). |
| `wails3 setup entitlements` | Assistente interativo para configurar os direitos do macOS. Opção: `--output` (caminho; padrão: `build/darwin/entitlements.plist`). |
| `wails3 sign [GOOS=…]` | Wrapper que executa a tarefa `*:sign` do Taskfile específica da plataforma para o sistema operacional atual (ou para o especificado por meio de `GOOS`). |
| `wails3 tool sign` | Ponto de entrada de baixo nível para assinatura direta. Opções: `--input`, `--output`, `--verbose`, `--certificate`, `--password`, `--thumbprint`, `--timestamp`, `--identity`, `--entitlements`, `--hardened-runtime`, `--notarize`, `--keychain-profile`, `--pgp-key`, `--pgp-password`, `--role`. |

**Não existe** um subcomando `wails3 signing` — para credenciais do chaveiro, use `xcrun notarytool store-credentials`; para chaves PGP, use `gpg` diretamente (o assistente `wails3 setup signing` automatiza ambos).

## Ferramentas

`wails3 tool <subcommand>`:

| Subcomando | Descrição |
| --- | --- |
| `tool checkport` | Verifica se uma porta TCP está aberta (útil para aguardar o Vite). |
| `tool watcher` | Executa um comando sempre que os arquivos monitorados são alterados. |
| `tool cp` | Copia arquivos de maneira multiplataforma. |
| `tool buildinfo` | Exibe as informações de compilação Go incorporadas a um binário. |
| `tool package` | Cria um pacote Linux (`deb`, `rpm`, `archlinux`) a partir de `build/linux/nfpm`. |
| `tool version` | Incrementa a versão semântica de um projeto. |
| `tool lipo` | Combina binários de várias arquiteturas do macOS em um binário universal. |
| `tool capabilities` | Verifica no sistema a disponibilidade do GTK3/GTK4 e do WebKit. |
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
| `ios overlay:gen` | Gera o overlay Go para o shim da ponte do iOS. |
| `ios xcode:gen` | Gera um projeto do Xcode no diretório de saída. |

## Caminhos de saída da compilação

- Os binários nativos são gerados em `bin/<APP_NAME>` (ou em `bin/<APP_NAME>.exe` no Windows). Não há `build/bin/`.
- Da mesma forma, as saídas empacotadas (`.app`, `.dmg`, instalador NSIS, MSIX, DEB/RPM/AppImage) são geradas em `bin/` (ou em subdiretórios específicos da plataforma criados pela tarefa correspondente do Taskfile).

## Opções globais

| Opção | Aplicável a | Descrição |
| --- | --- | --- |
| `--no-colour` | Todos os comandos | Desativa as cores ANSI na saída da CLI. |

---

**Dúvidas?** Pergunte no [Discord](https://discord.gg/JDdSxwjhGf) ou consulte os [exemplos](https://github.com/wailsapp/wails/tree/master/v3/examples).
