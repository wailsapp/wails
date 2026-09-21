---
title: "Referência da CLI"
description: "Referência completa dos comandos da CLI do Wails"
slug: "guides/cli"
sourcePath: "guides/cli.md"
---

A CLI do Wails oferece um conjunto abrangente de comandos para ajudar você a desenvolver, compilar e manter seus aplicativos Wails.

## Comandos principais

Os comandos principais são usados para criar, desenvolver e compilar projetos.

Todos os comandos da CLI têm o seguinte formato: `wails3 <command>`.

### `init`

Inicializa um novo projeto Wails. Durante a inicialização, o comando `go mod tidy` é executado para atualizar os pacotes do projeto. Essa etapa pode ser ignorada usando a opção `-skipgomodtidy` com o comando `init`.

```bash
wails3 init [flags]
```

#### Opções

| Opção | Descrição | Padrão |
| --- | --- | --- |
| `-p` | Nome do pacote Go | `main` |
| `-t` | Nome ou URL do modelo | `vanilla` |
| `-n` | Nome do projeto |  |
| `-d` | Diretório do projeto | `.` |
| `-q` | Suprimir a saída | `false` |
| `-l` | Listar modelos | `false` |
| `-mod` | Caminho do módulo Go (calculado a partir de `-git` se omitido) |  |
| `-git` | URL do repositório Git |  |
| `-s` | Ignorar o aviso ao usar um modelo remoto | `false` |
| `-productname` | Nome do produto | `My Product` |
| `-productdescription` | Descrição do produto | `My Product Description` |
| `-productversion` | Versão do produto | `0.1.0` |
| `-productcompany` | Nome da empresa | `My Company` |
| `-productcopyright` | Aviso de direitos autorais | `© now, My Company` |
| `-productcomments` | Comentários nos arquivos | `This is a comment` |
| `-productidentifier` | Identificador do produto |  |
| `-skipgomodtidy` | Ignorar go mod tidy | `false` |

A opção `-git` aceita vários formatos de URL do Git:

- HTTPS: `https://github.com/username/project`
- SSH: `git@github.com:username/project` ou `ssh://git@github.com/username/project`
- Protocolo Git: `git://github.com/username/project`
- Sistema de arquivos: `file:///path/to/project.git`

Quando fornecida, essa opção:

1. Inicializa um repositório Git no diretório do projeto
2. Define a URL especificada como o repositório remoto origin
3. Atualiza o nome do módulo em `go.mod` para corresponder à URL do repositório
4. Adiciona todos os arquivos

### `dev`

Executa o aplicativo no modo de desenvolvimento. Isso oferece uma visualização em tempo real do código do frontend, permitindo fazer alterações e vê-las refletidas no aplicativo em execução sem precisar recompilar todo o aplicativo. As alterações no código Go também serão detectadas, e o aplicativo será recompilado e reiniciado automaticamente.

```bash
wails3 dev [flags]
```

#### Opções

| Opção | Descrição | Padrão |
| --- | --- | --- |
| `-config` | Caminho do arquivo de configuração | `./build/config.yml` |
| `-port` | Porta do servidor de desenvolvimento do Vite | `9245` |
| `-s` | Ativar HTTPS | `false` |

@note{type="info"}
Isso equivale a executar `wails3 task dev` e executa a tarefa `dev` no Taskfile principal do projeto. Você pode personalizar esse comportamento editando o arquivo `Taskfile.yml`.

@end

### `build`

Compila uma versão de depuração do aplicativo. Por padrão, a compilação é feita para a plataforma e a arquitetura atuais.

```bash
wails3 build [flags] [CLI variables...]
```

#### Opções

| Opção | Descrição | Padrão |
| --- | --- | --- |
| `-tags` | Tags adicionais de compilação do Go (separadas por vírgulas) |  |

Você pode passar variáveis da CLI para personalizar a compilação:

```bash
wails3 build PLATFORM=linux CONFIG=production
```

Use a opção `-tags` para passar tags personalizadas de compilação do Go:

```bash
# Build with legacy GTK3 + WebKit2GTK 4.1 on Linux (default is GTK4 + WebKitGTK 6.0)
wails3 build -tags gtk3

# Build in server mode (no GUI)
wails3 build -tags server

# Multiple tags
wails3 build -tags gtk3,customtag
```

As tags são encaminhadas como `EXTRA_TAGS` para o Taskfile subjacente.

@note{type="info"}
Isso equivale a executar `wails3 task build`, que executa a tarefa `build` no Taskfile principal do projeto. Todas as variáveis da CLI passadas para `build` são encaminhadas para a tarefa subjacente. Você pode personalizar o processo de compilação editando o arquivo `Taskfile.yml`.

@end

### `package`

Cria pacotes específicos da plataforma para distribuição.

```bash
wails3 package [CLI variables...]
```

Você pode passar variáveis da CLI para personalizar o empacotamento:

```bash
wails3 package VERSION=2.0.0 OUTPUT=myapp.pkg
```

#### Tipos de pacote

Os seguintes tipos de pacote estão disponíveis para cada plataforma:

| Plataforma | Tipo de pacote |
| --- | --- |
| Windows | `.exe` |
| macOS | `.app`, |
| Linux | `.AppImage`, `.deb`, `.rpm`, `.archlinux` |

@note{type="info"}
Isso equivale a `wails3 task package`, que executa a tarefa `package` no Taskfile principal do projeto. Todas as variáveis da CLI passadas para `package` são encaminhadas para a tarefa subjacente. Você pode personalizar o processo de empacotamento editando o arquivo `Taskfile.yml`.

@end

### `task`

Executa tarefas definidas no arquivo Taskfile.yml do seu projeto. Esta é uma versão integrada do [Taskfile](https://taskfile.dev) que permite definir e executar tarefas personalizadas de compilação, teste e implantação.

```bash
wails3 task [taskname] [CLI variables...] [flags]
```

#### Variáveis da CLI

Você pode passar variáveis para as tarefas no formato `KEY=VALUE`:

```bash
wails3 task build PLATFORM=linux CONFIG=production
wails3 task deploy ENV=staging VERSION=1.2.3
```

Essas variáveis podem ser acessadas no arquivo Taskfile.yml usando a sintaxe de templates do Go:

```yaml
tasks:
  build:
    cmds:
      - echo "Building for {{.PLATFORM | default "darwin"}}"
      - echo "Config: {{.CONFIG | default "debug"}}"
```

#### Opções

| Opção | Descrição | Padrão |
| --- | --- | --- |
| `-h` | Exibe as instruções de uso do Task | `false` |
| `-i` | Cria um novo Taskfile.yml | `false` |
| `-list` | Lista as tarefas com descrições | `false` |
| `-list-all` | Lista todas as tarefas (com ou sem descrições) | `false` |
| `-json` | Formata a lista de tarefas como JSON | `false` |
| `-status` | Encerra com um código diferente de zero se a tarefa não estiver atualizada | `false` |
| `-f` | Força a execução mesmo quando a tarefa está atualizada | `false` |
| `-w` | Ativa o modo de monitoramento para a tarefa especificada | `false` |
| `-v` | Ativa o modo detalhado | `false` |
| `-version` | Exibe a versão do Task | `false` |
| `-s` | Desativa a exibição dos comandos | `false` |
| `-p` | Executa tarefas em paralelo | `false` |
| `-dry` | Compila e exibe as tarefas sem executá-las | `false` |
| `-summary` | Exibe um resumo sobre uma tarefa | `false` |
| `-x` | Propaga o código de saída da tarefa | `false` |
| `-dir` | Define o diretório de execução |  |
| `-taskfile` | Escolhe qual Taskfile executar |  |
| `-output` | Define o estilo de saída: [interleaved|group|prefixed] |  |
| `-c` | Saída colorida (habilitada por padrão) | `true` |
| `-C` | Limita o número de tarefas executadas simultaneamente |  |
| `-interval` | Intervalo de monitoramento de alterações (em segundos) |  |

#### Exemplos

```bash
# Run the default task
wails3 task

# Run a specific task
wails3 task test

# Run a task with variables
wails3 task build PLATFORM=windows ARCH=amd64

# List all available tasks
wails3 task --list

# Run multiple tasks in parallel
wails3 task -p task1 task2 task3

# Watch for changes and re-run task
wails3 task -w dev
```

### `mcp`

Inicia o servidor MCP do projeto Wails para o gerenciamento de projetos assistido por agentes. Ele é separado do servidor MCP compilado em um aplicativo em execução: `wails3 mcp` gerencia os arquivos do projeto e os comandos de ciclo de vida, enquanto o servidor MCP do aplicativo controla o WebView em execução.

```bash
wails3 mcp [flags]
```

O transporte é selecionado automaticamente:

- Quando um host MCP inicia o Wails com entrada e saída padrão redirecionadas por pipes, o servidor usa **stdio**.
- Quando executado interativamente em um terminal, o servidor usa **HTTP com streaming** em `127.0.0.1` e solicita ao sistema operacional uma porta livre.

Use `--stdio` ou `--http` para selecionar explicitamente um transporte. Use `--port 0` para escolher uma porta de loopback livre no modo HTTP.

#### Opções do MCP

| Opção | Descrição | Padrão |
| --- | --- | --- |
| `--root` | Raiz permitida do projeto. Caminhos e links simbólicos fora dela são rejeitados. | Diretório atual |
| `--token` | Token de sessão/bearer para ferramentas de modificação e controle de processos. Se não for definido, usa `WAILS_MCP_TOKEN`. | Gerado com segurança |
| `--stdio` | Força o transporte stdio. | Automático |
| `--http` | Força o transporte HTTP com streaming. | Automático |
| `--port` | Porta HTTP; `0` seleciona uma porta de loopback livre. | `0` |

No modo HTTP, o Wails imprime o endpoint e o token bearer na saída de erro padrão. No modo stdio, o token é incluído nas instruções de inicialização do MCP. O servidor não disponibiliza a execução arbitrária de comandos do shell. Modelos remotos e repositórios remotos do Git exigem aprovação explícita por meio da entrada `allowExternal` da ferramenta.

### `doctor`

Executa uma verificação do sistema e exibe um relatório de status.

```bash
wails3 doctor
```

## Comandos de geração

Os comandos de geração ajudam a criar vários recursos do projeto, como bindings, ícones e arquivos de build. Todos os comandos de geração usam o comando base: `wails3 generate <command>`.

### `generate bindings`

Gera bindings e modelos para o seu código Go.

```bash
wails3 generate bindings [flags] [patterns...]
```

#### Opções

| Opção | Descrição | Padrão |
| --- | --- | --- |
| `-f` | Opções adicionais de build do Go |  |
| `-d` | Diretório de saída | `frontend/bindings` |
| `-models` | Nome do arquivo de modelos | `models` |
| `-index` | Nome do arquivo de índice | `index` |
| `-ts` | Gera TypeScript | `false` |
| `-i` | Usa interfaces do TypeScript | `false` |
| `-b` | Usa o runtime incluído no pacote | `false` |
| `-names` | Usa nomes em vez de IDs | `false` |
| `-noindex` | Ignora os arquivos de índice | `false` |
| `-noevents` | Não gerar associações relacionadas a eventos | `false` |
| `-dry` | Simulação | `false` |
| `-silent` | Modo silencioso | `false` |
| `-v` | Saída de depuração | `false` |
| `-clean` | Limpar o diretório de saída antes da geração | `true` |

### `generate build-assets`

Gera os recursos de compilação do seu aplicativo.

```bash
wails3 generate build-assets [flags]
```

#### Opções

| Opção | Descrição | Padrão |
| --- | --- | --- |
| `-name` | Nome do projeto |  |
| `-dir` | Diretório de saída | `build` |
| `-silent` | Suprimir a saída | `false` |
| `-company` | Nome da empresa |  |
| `-productname` | Nome do produto |  |
| `-description` | Descrição do produto |  |
| `-version` | Versão do produto |  |
| `-identifier` | Identificador do produto | `com.wails.[name]` |
| `-copyright` | Aviso de direitos autorais |  |
| `-comments` | Comentários do arquivo |  |

### `generate icons`

Gera os ícones do aplicativo.

```bash
wails3 generate icons [flags]
```

#### Opções

| Opção | Descrição | Padrão |
| --- | --- | --- |
| `-input` | Arquivo PNG de entrada | Obrigatório |
| `-windowsfilename` | Nome do arquivo de saída para Windows |  |
| `-macfilename` | Nome do arquivo de saída para macOS |  |
| `-sizes` | Tamanhos dos ícones (separados por vírgulas) | `256,128,64,48,32,16` |
| `-example` | Gerar ícone de exemplo | `false` |
| `-iconcomposerinput` | Arquivo do Icon Composer de entrada (`.icon`) |  |
| `-macassetdir` | Diretório de saída dos recursos para Mac (Assets.car + icns) |  |

#### Icon Composer (macOS)

No macOS 26 ou posterior, você pode usar arquivos `.icon` do Icon Composer para gerar `Assets.car` e `icons.icns`:

```bash
wails3 generate icons -iconcomposerinput build/appicon.icon -macassetdir build
```

Isso compila o arquivo `.icon` usando o comando `actool` da Apple. Requer o Xcode com o `actool` na versão 26 ou posterior.

Ao usar o Icon Composer, defina `cfBundleIconName` no seu `build/config.yml` para que corresponda ao nome do arquivo `.icon` (sem a extensão):

```yaml
info:
  cfBundleIconName: "appicon"
```

Se não estiver definido e `Assets.car` existir, o valor padrão será `"appicon"`.

### `generate syso`

Gera um arquivo .syso do Windows.

```bash
wails3 generate syso [flags]
```

#### Opções

| Opção | Descrição | Padrão |
| --- | --- | --- |
| `-manifest` | Caminho para o arquivo de manifesto | Obrigatório |
| `-icon` | Caminho para o arquivo de ícone | Obrigatório |
| `-info` | Caminho para o arquivo de informações da versão |  |
| `-arch` | Arquitetura de destino | GOARCH atual |
| `-out` | Nome do arquivo de saída | `rsrc_windows_[arch].syso` |

### `generate .desktop`

Gera um arquivo .desktop do Linux.

```bash
wails3 generate .desktop [flags]
```

#### Opções

| Opção | Descrição | Padrão |
| --- | --- | --- |
| `-name` | Nome do aplicativo | Obrigatório |
| `-exec` | Caminho para o executável | Obrigatório |
| `-icon` | Caminho para o ícone |  |
| `-categories` | Categorias do aplicativo | `Utility` |
| `-comment` | Comentário do aplicativo |  |
| `-terminal` | Executar no terminal | `false` |
| `-keywords` | Palavras-chave de pesquisa |  |
| `-version` | Versão do aplicativo |  |
| `-genericname` | Nome genérico |  |
| `-startupnotify` | Mostrar notificação de inicialização | `false` |
| `-mimetype` | Tipos MIME compatíveis |  |
| `-output` | Nome do arquivo de saída | `[name].desktop` |

### `generate runtime`

Gera a versão pré-compilada do runtime.

```bash
wails3 generate runtime
```

### `generate constants`

Gera constantes JavaScript a partir de código Go.

```bash
wails3 generate constants
```

### `generate webview2bootstrapper`

Gera um instalador bootstrap do WebView2 para distribuição no Windows.

```bash
wails3 generate webview2bootstrapper [flags]
```

### `generate template`

Cria a estrutura de um novo diretório de modelo de projeto.

```bash
wails3 generate template [flags]
```

### `generate appimage`

Gera uma AppImage para Linux.

```bash
wails3 generate appimage [flags]
```

#### Opções

| Opção | Descrição | Padrão |
| --- | --- | --- |
| `-binary` | Caminho para o arquivo binário | Obrigatório |
| `-icon` | Caminho para o arquivo de ícone | Obrigatório |
| `-desktop` | Caminho para o arquivo .desktop | Obrigatório |
| `-builddir` | Diretório de compilação | Diretório temporário |
| `-output` | Diretório de saída | `.` |

## Comandos de serviço

Os comandos de serviço ajudam a gerenciar os serviços do Wails. Todos os comandos de serviço usam o comando-base: `wails3 service <command>`.

### `service init`

Inicializa um novo serviço.

```bash
wails3 service init [flags]
```

#### Opções

| Opção | Descrição | Padrão |
| --- | --- | --- |
| `-n` | Nome do serviço | `example_service` |
| `-d` | Descrição do serviço | `Example service` |
| `-p` | Nome do pacote |  |
| `-o` | Diretório de saída | `.` |
| `-q` | Suprimir a saída | `false` |
| `-a` | Nome do autor |  |
| `-v` | Versão |  |
| `-w` | URL do site |  |
| `-r` | URL do repositório |  |
| `-l` | Licença |  |

## Comandos de ferramenta

Os comandos de ferramenta fornecem utilitários para desenvolvimento e depuração. Todos os comandos de ferramenta usam o comando-base: `wails3 tool <command>`.

### `tool checkport`

Verifica se uma porta está aberta. É útil para testar se o vite está em execução.

```bash
wails3 tool checkport [flags]
```

#### Opções

| Opção | Descrição | Padrão |
| --- | --- | --- |
| `-port` | Porta a ser verificada | `9245` |
| `-host` | Host a ser verificado | `localhost` |

### `tool watcher`

Monitora arquivos e executa um comando quando eles são alterados.

```bash
wails3 tool watcher [flags]
```

#### Opções

| Opção | Descrição | Padrão |
| --- | --- | --- |
| `-config` | Caminho do arquivo de configuração | `./build/config.yml` |
| `-ignore` | Padrões a serem ignorados |  |
| `-include` | Padrões a serem incluídos |  |

### `tool cp`

Copia arquivos.

```bash
wails3 tool cp
```

### `tool buildinfo`

Exibe informações de compilação sobre o aplicativo.

```bash
wails3 tool buildinfo
```

### `tool version`

Incrementa uma versão semântica com base nas opções fornecidas.

```bash
wails3 tool version [flags]
```

#### Opções

| Opção | Descrição | Padrão |
| --- | --- | --- |
| `-v` | Versão atual a ser incrementada |  |
| `-major` | Incrementar a versão principal | `false` |
| `-minor` | Incrementar a versão secundária | `false` |
| `-patch` | Incrementar a versão de correção | `false` |
| `-prerelease` | Incrementar a versão de pré-lançamento (por exemplo, de alpha.5 para alpha.6) | `false` |

O comando segue esta ordem de precedência: principal > secundária > correção > pré-lançamento. Ele preserva o prefixo "v", se estiver presente na versão de entrada, bem como quaisquer componentes de pré-lançamento e metadados.

Exemplo de uso:

```bash
wails3 tool version -v 1.2.3 -major      # Output: 2.0.0
wails3 tool version -v v1.2.3 -minor     # Output: v1.3.0
wails3 tool version -v 1.2.3-alpha -patch # Output: 1.2.4-alpha
wails3 tool version -v v3.0.0-alpha.5 -prerelease # Output: v3.0.0-alpha.6
```

### `tool package`

Gera pacotes Linux (deb, rpm, archlinux).

```bash
wails3 tool package [flags]
```

#### Opções

| Opção | Descrição | Padrão |
| --- | --- | --- |
| `-format` | Formato do pacote (deb, rpm, archlinux) | `deb` |
| `-name` | Nome do executável | `myapp` |
| `-config` | Caminho do arquivo de configuração |  |
| `-out` | Diretório de saída | `.` |

### `tool lipo`

Cria um binário universal do macOS combinando binários específicos de cada arquitetura.

```bash
wails3 tool lipo [flags]
```

#### Opções

| Opção | Descrição | Padrão |
| --- | --- | --- |
| `-output` | Caminho do binário de saída |  |

### `tool capabilities`

Verifica os recursos de compilação do sistema (disponibilidade do GTK4/GTK3 no Linux).

```bash
wails3 tool capabilities
```

### `tool docker-mounts`

Gera opções de montagem de volume do Docker para compilação cruzada. Produz opções `-v` para o cache de módulos Go e para todas as diretivas locais `replace` em `go.mod`, para uso nos comandos `docker run` do Taskfile.

```bash
wails3 tool docker-mounts
```

### `tool has`

Verifica se uma ferramenta ou um recurso está disponível, imprimindo `true` ou `false` na saída padrão. Foi projetado para uso em variáveis `sh:` do Taskfile como alternativa multiplataforma a `command -v`.

Use `|` para verificar se qualquer uma entre várias alternativas está disponível.

```bash
wails3 tool has <tool>
```

#### Exemplos

```bash
# Check for a C compiler (gcc or clang)
wails3 tool has gcc|clang

# Check for a specific tool
wails3 tool has git
wails3 tool has node
```

#### Uso em um Taskfile

```yaml
vars:
  HAS_CC:
    sh: 'wails3 tool has gcc|clang'
```

### `tool has-cc`

@note{type="caution" title="Obsoleto"}
`wails3 tool has-cc` está obsoleto. Atualize seu Taskfile para usar `wails3 tool has gcc|clang` em seu lugar.

@end

Um alias compatível com versões anteriores para `wails3 tool has gcc|clang`. Verifica se `gcc` ou `clang` está disponível no PATH e imprime `true` ou `false`.

```bash
wails3 tool has-cc
```

## Comandos de atualização

Os comandos de atualização ajudam a gerenciar e atualizar os recursos do projeto. Todos eles usam o comando-base: `wails3 update <command>`.

### `update cli`

Atualiza a CLI do Wails para uma nova versão.

```bash
wails3 update cli [flags]
```

#### Opções

| Opção | Descrição | Padrão |
| --- | --- | --- |
| `-pre` | Atualizar para a versão de pré-lançamento mais recente | `false` |
| `-version` | Atualizar para uma versão específica |  |
| `-nocolour` | Desativar a saída colorida | `false` |

O comando update cli permite atualizar sua instalação da CLI do Wails. Por padrão, ele atualiza para a versão estável mais recente. Você pode usar a opção `-pre` para atualizar para a versão de pré-lançamento mais recente ou especificar uma versão usando a opção `-version`.

Após a atualização, lembre-se de atualizar o arquivo go.mod do seu projeto para usar a mesma versão:

```bash
require github.com/wailsapp/wails/v3 v3.x.x
```

### `update build-assets`

Atualiza os recursos de compilação usando o arquivo de configuração fornecido.

```bash
wails3 update build-assets [flags]
```

#### Opções

| Opção | Descrição | Padrão |
| --- | --- | --- |
| `-config` | Caminho do arquivo de configuração |  |
| `-dir` | Diretório de saída | `build` |
| `-silent` | Suprimir a saída | `false` |
| `-company` | Nome da empresa |  |
| `-productname` | Nome do produto |  |
| `-description` | Descrição do produto |  |
| `-version` | Versão do produto |  |
| `-identifier` | Identificador do produto |  |
| `-copyright` | Aviso de direitos autorais |  |
| `-comments` | Comentários do arquivo |  |

## Comandos utilitários

Os comandos utilitários oferecem atalhos úteis para tarefas comuns. Use esses comandos diretamente com o comando base: `wails3 <command>`.

### `docs`

Abre a documentação do Wails no navegador padrão.

```bash
wails3 docs
```

### `releasenotes`

Exibe as notas de versão da versão atual ou da versão especificada.

```bash
wails3 releasenotes [flags]
```

#### Opções

| Opção | Descrição | Padrão |
| --- | --- | --- |
| `-v` | Versão cujas notas de versão serão exibidas |  |
| `-n` | Desativar a saída colorida | `false` |

### `version`

Exibe a versão atual do Wails.

```bash
wails3 version
```

### `sponsor`

Abre a página de patrocínio do Wails no navegador padrão.

```bash
wails3 sponsor

```
