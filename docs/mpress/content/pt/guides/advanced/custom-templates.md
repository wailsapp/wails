---
title: "Como criar templates personalizados"
description: "Como gerar, personalizar e hospedar seus próprios templates de projeto Wails v3"
slug: "guides/advanced/custom-templates"
sourcePath: "guides/advanced/custom-templates.md"
---

O Wails inclui um conjunto de templates integrados, mas você pode criar os seus e compartilhá-los com a comunidade. Um template personalizado é simplesmente um repositório Git — depois que ele é hospedado publicamente, qualquer pessoa pode gerar a estrutura inicial de um projeto a partir dele com um único comando.

## Gerar a estrutura básica de um template

O comando `wails3 generate template` gera um diretório de template pronto para personalização:

```bash
wails3 generate template -name MyTemplate
```

Todas as opções:

| Opção | Descrição | Padrão |
| --- | --- | --- |
| `-name` | Nome do template (obrigatório) | — |
| `-author` | Nome do autor | — |
| `-description` | Descrição curta exibida na CLI | — |
| `-helpurl` | URL da documentação deste template | — |
| `-version` | Versão inicial | `v0.0.1` |
| `-frontend` | Copia um diretório de frontend existente para o template | — |
| `-dir` | Local onde o diretório do template será criado | Diretório atual |

Exemplo com todas as opções:

```bash
wails3 generate template \
  -name "My Template" \
  -author "Your Name" \
  -description "React + custom setup" \
  -helpurl "https://github.com/yourname/my-template" \
  -version "v1.0.0" \
  -frontend ./my-existing-frontend
```

O diretório gerado tem esta estrutura:

```
MyTemplate/
├── template.yaml          # Template metadata — edit this
├── NEXTSTEPS.md           # Guidance for you as the template author — delete before publishing
├── README.md              # Shown to users after they create a project
├── main.go.tmpl           # Application entry point
├── greetservice.go        # Example Go service
├── go.mod.tmpl            # Go module file
├── go.sum.tmpl            # Go checksums
├── gitignore.tmpl         # Becomes .gitignore in generated projects
├── Taskfile.tmpl.yml      # Build task definitions
└── frontend/              # Your frontend code
```

@note{type="tip" title="Leia NEXTSTEPS.md"}
O arquivo `NEXTSTEPS.md` gerado contém orientações detalhadas sobre todas as partes do template. Leia-o antes de começar a personalização. Exclua-o antes de publicar — ele não pode aparecer nos projetos criados a partir do seu template.

@end

## Configurar os metadados do template

Abra `template.yaml` para definir os metadados do seu template:

```yaml
# yaml-language-server: $schema=https://v3.wails.io/schemas/template.v3.json
name: "My Template"
shortname: my-template
author: Your Name
description: A template with my preferred setup
helpurl: https://github.com/yourname/my-template
version: v1.0.0
wailsVersion: 3
```

O campo `wailsVersion` é **obrigatório** e deve ser `3`. O comentário `# yaml-language-server` no início do arquivo habilita o preenchimento automático e a validação em linha no VS Code (com a [extensão YAML](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml)) e nas IDEs da JetBrains — você pode mantê-lo ou removê-lo; ele não afeta a execução.

## Personalizar o template

### Frontend

O conteúdo do diretório `frontend/` é copiado sem alterações para todos os projetos criados a partir do seu template. Substitua o conteúdo provisório pelo seu frontend real:

@tabs
[Começar do zero]
```bash
cd MyTemplate/frontend
npm create vite@latest .
```

Siga as instruções exibidas e instale as dependências:

```bash
npm install
```

[Usar um projeto existente]
Ao gerar o template, passe `-frontend` para copiar um frontend existente em uma única etapa:

```bash
wails3 generate template -name MyTemplate -frontend ./my-app/frontend
```

Como alternativa, copie-o manualmente para o diretório `frontend/` depois.

@end

### Tarefas de build

`Taskfile.tmpl.yml` define o fluxo de build. Atualize as tarefas `install:frontend:deps` e `build:frontend` para corresponder ao conjunto de ferramentas do seu frontend:

```yaml
tasks:
  install:frontend:deps:
    dir: frontend
    cmds:
      - npm install       # replace with pnpm install, yarn, etc.

  build:frontend:
    dir: frontend
    deps: [install:frontend:deps, generate:bindings]
    cmds:
      - npm run build     # replace with your build command
```

### Aplicação Go

O arquivo `main.go.tmpl` é o ponto de entrada da aplicação. Ele é processado pelo mecanismo de templates do Wails quando um projeto é criado — variáveis de template como `{{.ProductName}}` são substituídas pelos valores fornecidos pelo usuário.

Para editá-lo como um arquivo Go real (com suporte da IDE), renomeie-o temporariamente como `main.go`, faça as alterações e, antes do commit, restaure o nome `main.go.tmpl`.

#### Variáveis de template

Estas variáveis estão disponíveis em qualquer arquivo `.tmpl`:

| Variável | Descrição | Exemplo |
| --- | --- | --- |
| `{{.ProjectName}}` | Nome do projeto fornecido pelo usuário | `"MyApp"` |
| `{{.BinaryName}}` | Nome do arquivo binário | `"myapp"` |
| `{{.ProductName}}` | Nome de exibição do produto | `"My Application"` |
| `{{.ProductDescription}}` | Descrição do produto | `"An awesome application"` |
| `{{.ProductVersion}}` | Versão do produto | `"1.0.0"` |
| `{{.ProductCompany}}` | Nome da empresa/do autor | `"My Company Ltd"` |
| `{{.ProductCopyright}}` | Texto de direitos autorais | `"Copyright 2024 My Company Ltd"` |
| `{{.ProductComments}}` | Comentários adicionais sobre o produto | `"Built with Wails"` |
| `{{.ProductIdentifier}}` | Identificador do produto no formato DNS reverso | `"com.mycompany.myapp"` |
| `{{.ModulePath}}` | Caminho do módulo Go | `"github.com/you/myapp"` |
| `{{.WailsVersion}}` | Versão do Wails usada para criar o projeto | `"3.0.0"` |
| `{{.Typescript}}` | `true` se o nome do template terminar em `-ts` | `true` |
| `{{.Opn}}` | `{{` literal — use escape dentro dos templates | `{{` |
| `{{.Cls}}` | `}}` literal — use escape dentro dos templates | `}}` |

@note{type="tip"}
Qualquer arquivo do seu template pode ser um arquivo `.tmpl` — inclusive arquivos HTML, JSON e YAML. Os arquivos sem o sufixo `.tmpl` são copiados sem alterações.

@end

## Testar seu template localmente

Antes de publicar, teste o template criando um projeto a partir de um caminho local:

```bash
wails3 init -n testproject -t /path/to/MyTemplate
```

Em seguida, verifique se o projeto funciona:

```bash
cd testproject
wails3 dev    # development mode with hot reload
wails3 build  # production binary
```

Verifique se:

- O recarregamento automático do frontend funciona
- As alterações no código Go recompilam e reiniciam a aplicação
- O binário de produção em `bin/` é executado corretamente

## Publicar no GitHub

@steps
### **Crie um repositório público no GitHub** para seu template. A raiz do repositório deve conter `template.yaml`.
### **Exclua `NEXTSTEPS.md`** — esse arquivo contém orientações para autores de templates e não deve aparecer nos projetos que os usuários criarem com seu template.
### **Faça o commit e o push** do conteúdo do diretório do template como a raiz do repositório:
```bash
git init
git add .
git commit -m "Initial template"
git remote add origin https://github.com/yourname/my-template.git
git push -u origin main
```

### **Crie uma tag de versão** usando versionamento semântico:
```bash
git tag v1.0.0
git push origin v1.0.0
```

@end

Agora os usuários podem criar projetos com seu template:

```bash
# Latest commit on the default branch
wails3 init -n myapp -t https://github.com/yourname/my-template

# Pinned to a specific release tag
wails3 init -n myapp -t https://github.com/yourname/my-template@v1.0.0
```

@note{type="caution" title="Aviso sobre template de terceiros"}
Quando um usuário instala um template remoto, o Wails exibe um aviso explicando que o template é código de terceiros e que o projeto Wails não se responsabiliza pelo conteúdo dele. Os usuários devem confirmar explicitamente antes que o projeto seja criado.

Como autor do template, você é responsável pela segurança e pela correção de todo o código contido nele.

Templates são código. As tarefas `Taskfile.yml` são comandos de shell arbitrários, e `wails3 dev`, `wails3 build` e `wails3 package` os executam com as permissões do usuário, assim como quaisquer scripts `package.json` que eles invoquem. Antes de executar qualquer um desses comandos, revise `Taskfile.yml`, o diretório `build/` e `frontend/package.json` do template remoto.

@end

## Práticas recomendadas

- **Escreva um `README.md` claro** — ele é exibido aos usuários depois que criam um projeto. Explique como executar, compilar e personalizar o projeto.
- **Preencha `helpurl`** — inclua um link para seu repositório ou para uma documentação específica. Os usuários verão esse link na lista de templates da CLI do Wails.
- **Fixe as versões das dependências do frontend** em `package.json` para evitar falhas de instalação causadas por atualizações upstream.
- **Teste antes de criar a tag** — crie um novo projeto a partir da versão com a tag antes de anunciá-la à comunidade.
- **Mantenha `wailsVersion: 3`** — esse campo informa ao Wails qual é a versão principal de destino do template. Não o altere.
- **Atualize regularmente** — mantenha as dependências atualizadas e teste com as novas versões do Wails.
