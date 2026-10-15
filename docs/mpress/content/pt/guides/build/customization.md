---
title: "Personalização da compilação"
description: "Personalize seu processo de compilação usando Task e Taskfile.yml"
slug: "guides/build/customization"
sourcePath: "guides/build/customization.md"
---

## Visão geral

O sistema de compilação do Wails é uma ferramenta flexível e poderosa, projetada para simplificar o processo de compilação dos seus aplicativos Wails. Ele utiliza o [Task](https://taskfile.dev), um executor de tarefas que permite definir e executar tarefas com facilidade. Embora o sistema de compilação da v3 seja o padrão, o Wails incentiva uma abordagem de "use as ferramentas que preferir", permitindo que os desenvolvedores personalizem o processo de compilação conforme necessário.

Saiba mais sobre como usar o Task na [documentação oficial](https://taskfile.dev/usage/).

## Task: o coração do sistema de compilação

O [Task](https://taskfile.dev) é uma alternativa moderna ao Make, escrita em Go. Ele usa um arquivo YAML para definir tarefas e suas dependências. No sistema de compilação do Wails, o [Task](https://taskfile.dev) desempenha um papel central na orquestração do processo de compilação.

O `Taskfile.yml` principal fica na raiz do projeto, enquanto as tarefas específicas de cada plataforma são definidas em arquivos `build/<platform>/Taskfile.yml`. Um arquivo `Taskfile.yml` comum no diretório `build` contém tarefas compartilhadas entre as plataformas.

@filetree

- Project Root
  - Taskfile.yml
  - build
    - windows/Taskfile.yml
    - darwin/Taskfile.yml
    - linux/Taskfile.yml
    - Taskfile.yml
@end

## Taskfile.yml

O arquivo `Taskfile.yml` na raiz do projeto é o principal ponto de entrada do sistema de compilação. Ele define as tarefas e suas dependências. Este é o arquivo `Taskfile.yml` padrão:

```yaml
version: '3'

includes:
  common: ./build/Taskfile.yml
  windows: ./build/windows/Taskfile.yml
  darwin: ./build/darwin/Taskfile.yml
  linux: ./build/linux/Taskfile.yml

vars:
  APP_NAME: "myproject"
  BIN_DIR: "bin"
  VITE_PORT: '{{.WAILS_VITE_PORT | default 9245}}'

tasks:
  build:
    summary: Builds the application
    cmds:
      - task: "{{OS}}:build"

  package:
    summary: Packages a production build of the application
    cmds:
      - task: "{{OS}}:package"

  run:
    summary: Runs the application
    cmds:
      - task: "{{OS}}:run"

  dev:
    summary: Runs the application in development mode
    cmds:
      - wails3 dev -config ./build/config.yml -port {{.VITE_PORT}}


```

## Taskfiles específicos de cada plataforma

Cada plataforma tem seu próprio Taskfile, localizado nos diretórios de plataforma dentro do diretório `build`. Esses arquivos definem as tarefas principais da respectiva plataforma. Cada Taskfile inclui as tarefas comuns do arquivo `build/Taskfile.yml`.

### Windows

Local: `build/windows/Taskfile.yml`

O Taskfile específico do Windows inclui tarefas para compilar, empacotar e executar o aplicativo no Windows. Os principais recursos incluem:

- Compilação com flags opcionais de produção
- Geração do arquivo de ícone `.ico`
- Geração do arquivo `.syso` do Windows
- Criação de um instalador NSIS para empacotamento

### Linux

Local: `build/linux/Taskfile.yml`

O Taskfile específico do Linux inclui tarefas para compilar, empacotar e executar o aplicativo no Linux. Os principais recursos incluem:

- Compilação com flags opcionais de produção
- Criação de pacotes AppImage, deb, rpm e Arch Linux
- Geração do arquivo `.desktop` para aplicativos Linux

### macOS

Local: `build/darwin/Taskfile.yml`

O Taskfile específico do macOS inclui tarefas para compilar, empacotar e executar o aplicativo no macOS. Os principais recursos incluem:

- Compilação de binários para as arquiteturas amd64, arm64 e universal (ambas)
- Geração do arquivo de ícone `.icns`
- Criação de um pacote `.app` para distribuição
- Assinatura ad hoc de pacotes `.app`
- Configuração de flags de compilação e variáveis de ambiente específicas do macOS

## Execução de tarefas e aliases de comandos

O comando `wails3 task` é uma versão incorporada do [Taskfile](https://taskfile.dev) que executa as tarefas definidas no seu `Taskfile.yml`.

Os comandos `wails3 build` e `wails3 package` são aliases de `wails3 task build` e `wails3 task package`, respectivamente. Quando você executa esses comandos, o Wails os converte internamente na execução de tarefa correspondente:

- `wails3 build` → `wails3 task build`
- `wails3 package` → `wails3 task package`

### Como passar parâmetros para tarefas

Você pode passar variáveis da CLI para as tarefas usando o formato `KEY=VALUE`. Essas variáveis são encaminhadas pelos comandos de alias:

```bash
# These are equivalent:
wails3 build PLATFORM=linux CONFIG=production
wails3 task build PLATFORM=linux CONFIG=production

# Package with custom version:
wails3 package VERSION=2.0.0 OUTPUT=myapp.pkg
```

No seu `Taskfile.yml`, você pode acessar essas variáveis usando a sintaxe de templates do Go:

```yaml
tasks:
  build:
    cmds:
      - echo "Building for {{.PLATFORM | default "darwin"}}"
      - go build -tags {{.CONFIG | default "debug"}} -o myapp
```

## Processo de compilação comum

Em todas as plataformas, o processo de compilação normalmente inclui as seguintes etapas:

1. Organização dos módulos Go
2. Compilação do frontend
3. Geração de ícones
4. Compilação do código Go com flags específicas da plataforma
5. Empacotamento do aplicativo (específico da plataforma)

## Personalização do processo de compilação

Embora o sistema de compilação da v3 ofereça uma configuração padrão sólida, você pode personalizá-lo facilmente de acordo com as necessidades do seu projeto. Ao modificar o `Taskfile.yml` e os Taskfiles específicos de cada plataforma, você pode:

- Adicionar novas tarefas
- Modificar tarefas existentes
- Alterar a ordem de execução das tarefas
- Integrar outras ferramentas e scripts

Essa flexibilidade permite adaptar o processo de compilação às suas necessidades específicas e, ao mesmo tempo, aproveitar a estrutura fornecida pelo sistema de compilação do Wails.

@note{type="tip" title="Aprendendo a usar o Taskfile"}
Recomendamos fortemente a leitura da documentação do [Taskfile](https://taskfile.dev) para entender como usá-lo de forma eficaz. Você pode descobrir qual versão do Taskfile está incorporada à CLI do Wails executando `wails3 task --version`.

@end

## Modo de desenvolvimento

O sistema de compilação do Wails inclui um modo de desenvolvimento avançado que melhora a experiência do desenvolvedor ao oferecer recarregamento em tempo real e substituição de módulos a quente. Esse modo é ativado com o comando `wails3 dev`.

### Como funciona

Ao executar `wails3 dev`, ocorre o seguinte processo:

1. O comando verifica se há uma porta disponível e usa 9245 como padrão caso nenhuma seja especificada.
2. Ele configura as variáveis de ambiente para o servidor de desenvolvimento do frontend (Vite).
3. Ele inicia o monitor de arquivos usando a biblioteca [refresh](https://github.com/atterpac/refresh).

A biblioteca [refresh](https://github.com/atterpac/refresh) é responsável por monitorar alterações nos arquivos e acionar recompilações. Ela usa a configuração definida na chave `dev_mode` do arquivo `./build/config.yml`.  
Ela pode ser configurada para ignorar determinados diretórios e arquivos, definir quais arquivos devem ser monitorados e quais ações devem ser executadas quando alterações forem detectadas.  
A configuração padrão funciona muito bem, mas você pode personalizá-la conforme suas necessidades.

### Configuração

Veja um exemplo da estrutura:

```yaml
dev_mode:
  root_path: .
  log_level: warn
  debounce: 1000
  ignore:
    dir:
      - .git
      - node_modules
      - frontend
      - bin
    file:
      - .DS_Store
      - .gitignore
      - .gitkeep
    watched_extension:
      - "*.go"
    git_ignore: true
  executes:
    - cmd: wails3 task common:install:frontend:deps
      type: once
    - cmd: wails3 task common:dev:frontend
      type: background
    - cmd: go mod tidy
      type: blocking
    - cmd: wails3 task build
      type: blocking
    - cmd: wails3 task run
      type: primary
```

Esse arquivo de configuração permite:

- Definir o caminho raiz para o monitoramento de arquivos
- Configurar o nível de log
- Definir um tempo de debounce para eventos de alteração de arquivos
- Ignorar diretórios, arquivos ou extensões de arquivo específicos
- Definir comandos a serem executados quando houver alterações nos arquivos

### Personalização do modo de desenvolvimento

Você pode personalizar a experiência do modo de desenvolvimento modificando esses valores no arquivo `config.yml`.

Algumas formas de personalização incluem:

1. Alterar os diretórios ou arquivos monitorados
2. Ajustar o tempo de debounce para controlar a rapidez com que o sistema responde às alterações
3. Adicionar ou modificar os comandos de execução conforme as necessidades do seu projeto

### Uso de um navegador para desenvolvimento

Embora o Wails v2 oferecesse suporte completo ao uso de um navegador para desenvolvimento, isso causava muita confusão. Aplicativos que funcionavam no navegador não necessariamente funcionavam como aplicativos para desktop, pois nem todas as APIs de navegador estão disponíveis em webviews.

Para trabalhos de desenvolvimento voltados à interface do usuário, você ainda pode usar um navegador na v3 acessando, no modo de desenvolvimento, a URL do Vite em `http://localhost:9245`. Isso permite usar as poderosas ferramentas de desenvolvimento do navegador ao trabalhar no estilo e no layout. Esteja ciente de que os bindings do Go *não funcionarão* nesse modo.  
Quando estiver tudo pronto para testar funcionalidades como bindings e eventos, basta alternar para a visualização do aplicativo para desktop e verificar se tudo funciona perfeitamente no ambiente de produção.
