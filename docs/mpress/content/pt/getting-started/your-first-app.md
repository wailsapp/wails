---
title: "Seu primeiro aplicativo"
description: "Crie seu primeiro aplicativo desktop com Wails passo a passo"
slug: "getting-started/your-first-app"
sourcePath: "getting-started/your-first-app.md"
---

Este guia mostra como criar seu primeiro aplicativo Wails v3, abordando a configuração do projeto, a compilação e o fluxo de trabalho de desenvolvimento.

<br/>

<br/>

@steps
### Como criar um novo projeto
Abra o terminal e execute o seguinte comando para criar um novo projeto Wails:

```bash
wails3 init -n myfirstapp -t vanilla-js
```

Esta lição seleciona explicitamente o modelo JavaScript Vanilla. O modelo básico `vanilla` usa TypeScript. O comando cria `myfirstapp` com todos os arquivos necessários.

   <video src="/assets/wails_init.mp4" controls></video>

### Como explorar a estrutura do projeto
Acesse o diretório `myfirstapp`. Você encontrará vários arquivos e diretórios:

@filetree
- build/           Contém arquivos usados pelo processo de compilação
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
  - src/         Código-fonte do frontend
    - main.js    Arquivo JavaScript principal
  - package.json NPM package configuration
  - public/      Recursos estáticos
  - Inter Font License.txt Font license
- .gitignore      Arquivo de exclusões do Git
- README.md       Documentação do projeto
- Taskfile.yml    Tarefas do projeto
- go.mod          Arquivo de módulo Go
- go.sum          Somas de verificação do módulo Go
- greetservice.go Greeting service
- main.go         Código principal do aplicativo
@end

Reserve um momento para explorar esses arquivos e se familiarizar com a estrutura.

@note{type="info"}
Embora o Wails v3 use o [Task](https://taskfile.dev/) como sistema de compilação padrão, nada impede que você use o `make` ou qualquer outro sistema de compilação alternativo.

@end

### Como compilar seu aplicativo
Para compilar seu aplicativo, execute:

```bash
wails3 build
```

Com os Taskfiles gerados, esse comando cria uma compilação de produção em `bin`. Use `wails3 dev` para o modo de desenvolvimento ou `wails3 build DEV=true` para criar explicitamente uma compilação de desenvolvimento.

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

Você verá a interface inicial. Uma compilação de produção não garante o comportamento do inspetor nem do console de desenvolvimento; use o modo de desenvolvimento ao depurar.

### Modo de desenvolvimento
Também podemos executar o aplicativo no modo de desenvolvimento. Esse modo permite alterar o código do frontend e ver as alterações refletidas no aplicativo em execução sem precisar recompilar o aplicativo inteiro.

1. Abra uma nova janela do terminal.
2. Execute `wails3 dev`. O aplicativo será compilado e executado no modo de depuração.
3. Abra `frontend/index.html` no editor de sua preferência.
4. Altere o `placeholder="Your name"` do campo de entrada de nome para `placeholder="Your full name"` e salve o arquivo.
5. Salve o arquivo.

Essa alteração será refletida imediatamente no aplicativo.

Qualquer alteração no código do backend acionará uma nova compilação:

1. Abra `greetservice.go`.
2. Substitua a linha que contém `return "Hello " + name + "!"` por `return "Hello there " + name + "!"`.
3. Salve o arquivo.

O aplicativo será atualizado em poucos segundos.

     <video src="/assets/wails_dev.mp4" controls></video>

### Como empacotar seu aplicativo
Quando seu aplicativo estiver pronto para distribuição, você poderá criar pacotes específicos para cada plataforma:

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

Para obter informações mais detalhadas sobre as opções e configurações de empacotamento, consulte nosso [Guia de compilação e empacotamento](/guides/build/building/).

### Como configurar o controle de versão e o nome do módulo
Seu projeto é criado com o nome de módulo provisório `changeme`. Recomendamos atualizá-lo para que corresponda à URL do seu repositório:

1. Crie um novo repositório no GitHub (ou no serviço de hospedagem Git de sua preferência)
2. Inicialize o Git no diretório do projeto:
  ```bash
  git init
  git add .
  git commit -m "Initial commit"
  ```

3. Configure o repositório remoto (substitua pela URL do seu repositório):
  ```bash
  git remote add origin https://github.com/username/myfirstapp.git
  ```

4. Atualize o nome do módulo em `go.mod` para que corresponda à URL do seu repositório:
  ```bash
  go mod edit -module github.com/username/myfirstapp
  ```

5. Envie seu código:
  ```bash
  git push -u origin main
  ```


Isso garante que o nome do módulo Go esteja de acordo com as convenções de nomenclatura de módulos do Go e facilita o compartilhamento do código.

@note{type="tip" title="Dica profissional"}
Você pode automatizar todas as etapas de inicialização usando a opção `-git` ao criar o projeto:

```bash
wails3 init -n myfirstapp -t vanilla-js -git github.com/username/myfirstapp
```

Há suporte para vários formatos de URL do Git:

- HTTPS: `https://github.com/username/project`
- SSH: `git@github.com:username/project` ou `ssh://git@github.com/username/project`
- Protocolo Git: `git://github.com/username/project`
- Sistema de arquivos: `file:///path/to/project.git`

@end

@end

## Parabéns!

Você acabou de criar, desenvolver e empacotar seu primeiro aplicativo Wails. Este é apenas o começo do que você pode alcançar com o Wails v3.

## Próximos passos

Se você está começando a usar o Wails, recomendamos que leia nossos Tutoriais em seguida. Eles oferecem um guia prático sobre os diversos recursos do Wails. O primeiro tutorial é [Como criar um serviço](/tutorials/01-creating-a-service/).

Se você é um usuário mais avançado, consulte o [Guia de compilação e empacotamento](/guides/build/building/) para obter informações mais detalhadas sobre como usar o Wails.
