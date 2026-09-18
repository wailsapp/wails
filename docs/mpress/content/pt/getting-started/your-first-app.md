---
title: "Seu primeiro aplicativo"
description: "Crie seu primeiro aplicativo desktop com Wails passo a passo"
slug: "getting-started/your-first-app"
sourcePath: "getting-started/your-first-app.md"
---

Este guia mostra como criar seu primeiro aplicativo com Wails v3, desde a configuração do projeto até a compilação e o fluxo de trabalho de desenvolvimento.

<br/>

<br/>

@steps
### Como criar um novo projeto
Abra o terminal e execute o seguinte comando para criar um novo projeto Wails:

```bash
wails3 init -n myfirstapp
```

Esse comando cria um novo diretório chamado `myfirstapp` com todos os arquivos necessários.

   <video src="/assets/wails_init.mp4" controls></video>

### Como explorar a estrutura do projeto
Acesse o diretório `myfirstapp`. Você encontrará vários arquivos e diretórios:

@filetree
- build/           Contém os arquivos usados pelo processo de compilação
  - appicon.png  Ícone do aplicativo
  - config.yml   Configuração de compilação
  - Taskfile.yml Build tasks
  - darwin/      Arquivos de compilação específicos do macOS
    - Info.dev.plist Development configuration
    - Info.plist    Configuração de produção
    - Taskfile.yml  Tarefas de compilação do macOS
    - icons.icns    Ícone do aplicativo para macOS
  - linux/       Arquivos de compilação específicos do Linux
    - Taskfile.yml  Tarefas de compilação do Linux
    - appimage/     Empacotamento em AppImage
      - build.sh  Script de compilação do AppImage
    - nfpm/        Empacotamento com NFPM
      - nfpm.yaml Package configuration
      - scripts/  Scripts de compilação
  - windows/     Arquivos de compilação específicos do Windows
    - Taskfile.yml        Tarefas de compilação do Windows
    - icon.ico           Ícone do aplicativo para Windows
    - info.json          Metadados do aplicativo
    - wails.exe.manifest Windows manifest file
    - nsis/              Arquivos do instalador NSIS
      - project.nsi                    Arquivo de projeto do NSIS
      - wails_tools.nsh               Scripts auxiliares do NSIS
- frontend/        Arquivos do frontend do aplicativo
  - index.html   Arquivo HTML principal
  - main.js      Arquivo JavaScript principal
  - package.json NPM package configuration
  - public/      Recursos estáticos
  - Inter Font License.txt Font license
- .gitignore      Arquivo de exclusões do Git
- README.md       Documentação do projeto
- Taskfile.yml    Tarefas do projeto
- go.mod          Arquivo de módulo Go
- go.sum          Somas de verificação dos módulos Go
- greetservice.go Greeting service
- main.go         Código principal do aplicativo
@end

Reserve um momento para explorar esses arquivos e se familiarizar com a estrutura.

@note{type="info"}
Embora o Wails v3 use [Task](https://taskfile.dev/) como sistema de compilação padrão, nada impede que você use `make` ou qualquer outro sistema de compilação alternativo.

@end

### Como compilar seu aplicativo
Para compilar seu aplicativo, execute:

```bash
wails3 build
```

Esse comando compila uma versão de depuração do aplicativo e a salva em um novo diretório `bin`.

@note{type="info"}
`wails3 build` é uma forma abreviada de `wails3 task build` e executará a tarefa `build` em `Taskfile.yml`.

@end

     <video src="/assets/wails_build.mp4" controls></video>

Após a compilação, você poderá executá-lo como qualquer aplicativo comum:

@tabs{sync-key="platform"}
[Mac]
```sh
./bin/myfirstapp
```

[Windows]
```sh
bin\myfirstapp.exe
```

[Linux]
```sh
./bin/myfirstapp
```

@end

Você verá uma interface simples, que será o ponto de partida do seu aplicativo. Como essa é a versão de depuração, você também verá logs na janela do console. Isso é útil para fins de depuração.

### Modo de desenvolvimento
Também podemos executar o aplicativo no modo de desenvolvimento. Esse modo permite alterar o código do frontend e ver as mudanças refletidas no aplicativo em execução sem precisar recompilar o aplicativo inteiro.

1. Abra uma nova janela do terminal.
2. Execute `wails3 dev`. O aplicativo será compilado e executado no modo de depuração.
3. Abra `frontend/index.html` no editor de sua preferência.
4. Edite o código e substitua `Please enter your name below` por `Please enter your name below!!!`.
5. Salve o arquivo.

Essa alteração será refletida imediatamente no aplicativo.

Qualquer alteração no código do backend acionará uma nova compilação:

1. Abra `greetservice.go`.
2. Na linha que contém `return "Hello " + name + "!"`, substitua-o por `return "Hello there " + name + "!"`.
3. Salve o arquivo.

O aplicativo será atualizado em questão de segundos.

     <video src="/assets/wails_dev.mp4" controls></video>

### Como empacotar seu aplicativo
Quando o aplicativo estiver pronto para distribuição, você poderá criar pacotes específicos para cada plataforma:

@tabs{sync-key="platform"}
[Mac]
Para criar um pacote `.app`:

```bash
wails3 package
```

Isso criará uma compilação de produção e a empacotará em um pacote `.app` no diretório `bin`.

[Windows]
Para criar um instalador NSIS:

```bash
wails3 package
```

Isso criará uma compilação de produção e a empacotará em um instalador NSIS no diretório `bin`.

[Linux]
O Wails oferece suporte a vários formatos de pacote para distribuição no Linux:

```bash
# Create all package types (AppImage, deb, rpm, and Arch Linux)
wails3 package

# Or create specific package types
wails3 task linux:create:appimage  # AppImage format
wails3 task linux:create:deb       # Debian package
wails3 task linux:create:rpm       # Red Hat package
wails3 task linux:create:aur       # Arch Linux package
```

@end

Para obter informações mais detalhadas sobre as opções e a configuração de empacotamento, consulte nosso [Guia de compilação e empacotamento](/guides/build/building/).

### Como configurar o controle de versão e o nome do módulo
Seu projeto é criado com o nome de módulo provisório `changeme`. Recomendamos atualizá-lo para corresponder à URL do seu repositório:

1. Crie um novo repositório no GitHub (ou no serviço de hospedagem Git de sua preferência)
2. Inicialize o Git no diretório do projeto:
  ```bash
  git init
  git add .
  git commit -m "Initial commit"
  ```

3. Defina seu repositório remoto (substitua pela URL do seu repositório):
  ```bash
  git remote add origin https://github.com/username/myfirstapp.git
  ```

4. Atualize o nome do módulo em `go.mod` para corresponder à URL do seu repositório:
  ```bash
  go mod edit -module github.com/username/myfirstapp
  ```

5. Envie seu código:
  ```bash
  git push -u origin main
  ```


Isso garante que o nome do seu módulo Go siga as convenções de nomenclatura de módulos do Go e facilita o compartilhamento do código.

@note{type="tip" title="Dica profissional"}
Você pode automatizar todas as etapas de inicialização usando a opção `-git` ao criar seu projeto:

```bash
wails3 init -n myfirstapp -git github.com/username/myfirstapp
```

Há suporte para vários formatos de URL do Git:

- HTTPS: `https://github.com/username/project`
- SSH: `git@github.com:username/project` ou `ssh://git@github.com/username/project`
- Protocolo Git: `git://github.com/username/project`
- Sistema de arquivos: `file:///path/to/project.git`

@end

@end

## Parabéns!

Você acabou de criar, desenvolver e empacotar seu primeiro aplicativo Wails. Isso é apenas o começo do que você pode fazer com o Wails v3.

## Próximas etapas

Se você está começando a usar o Wails, recomendamos que leia a seguir nossos Tutoriais, que oferecem um guia prático sobre os diversos recursos do Wails. O primeiro tutorial é [Como criar um serviço](/tutorials/01-creating-a-service/).

Se você tem mais experiência, consulte o [Guia de compilação e empacotamento](/guides/build/building/) para obter informações mais detalhadas sobre como usar o Wails.
