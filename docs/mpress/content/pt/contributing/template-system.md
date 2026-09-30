---
title: "Sistema de templates"
description: "Como o Wails v3 gera a estrutura inicial de novos projetos, como os templates são organizados e como criar o seu próprio template."
slug: "contributing/template-system"
sourcePath: "contributing/template-system.md"
---

O Wails inclui um **sistema de templates** que permite ao `wails3 init` gerar um projeto pronto para execução. Um conjunto deliberadamente pequeno de frameworks tem templates integrados (Vanilla, React, Vue, Svelte); qualquer outro framework pode ser usado ao [fornecer seu próprio frontend](/guides/dev/frontend-frameworks/) ou publicar um [template personalizado](/guides/advanced/custom-templates/).

Esta página aborda:

1. Estrutura de diretórios dos templates
2. Como a CLI seleciona e renderiza templates
3. Criação de um novo template passo a passo
4. Atualização ou substituição de templates existentes
5. Solução de problemas e práticas recomendadas

---

## 1. Onde ficam os templates

```
v3/internal/templates/
├── _common/        # Files copied into EVERY project (Taskfile.yml, build/, etc.)
├── base/           # Backend-only "plain Go" base layer (frontend/ + NEXTSTEPS.md)
├── ios/            # iOS bootstrapper
├── vanilla/        vanilla-js/   # TypeScript (default) + JavaScript variant
├── react/          react-js/     # TypeScript (default) + JavaScript variant
├── vue/                          # TypeScript only
├── svelte/                       # TypeScript only
└── templates.go    # Registry + Install/Get APIs (no auto-registration via embed)
```

- **`_common/`** — estrutura básica universal (Taskfile, diretório `build/` e infraestrutura compartilhada), mesclada a todos os projetos.
- **`base/`** — a parte em Go da qual todos os templates partem. Observação: o próprio `base/` **não** contém um `template.json`; esse arquivo fica dentro de cada template específico de framework.
- **Pastas de frameworks** — contêm o frontend (`frontend/`), a configuração do framework e um `template.json` que descreve os metadados do template.
- Os nomes das pastas correspondem ao **ID do template** fornecido à CLI (`wails3 init -t react`).
- **Convenção de linguagem:** TypeScript é o padrão e usa o nome sem sufixo (`react`); quando existe uma variante JavaScript, ela recebe o sufixo `-js` (`react-js`). Os templates integrados declaram explicitamente sua linguagem com `typescript: true|false` em `template.yaml`. Os templates da comunidade ainda podem usar o sufixo legado `-ts`, que é aceito como alternativa.

> Todo o diretório `internal/templates/` é compilado no binário da CLI
>
> por meio de `//go:embed *`, para que os usuários possam gerar a estrutura inicial de projetos sem conexão com a internet.

---

## 2. Como o `wails3 init` usa templates

Cadeia de chamadas (sem `cmd/wails3/init.go` — a CLI é conectada diretamente em `cmd/wails3/main.go`):

```
cmd/wails3/main.go             (clir wiring)
       │
       ▼
internal/commands/init.go      Init(options *flags.Init) error
       │
       ▼
internal/templates/templates.go
       │   templates.Install(options)
       │   templates.GetDefaultTemplates()
       ▼
gosod.New(template.FS).Extract(options.ProjectDir, data)   // file extraction
       │
       ▼
go mod tidy (unless --skipgomodtidy / -skipgomodtidy)
```

Não existe uma API `Template.Load()` / `Template.CopyTo()` / `Template.Validate()` — a extração é realizada por `gosod` (`github.com/leaanthony/gosod`) usando o `fs.FS` incorporado.

### Opções de `wails3 init`

Definidas em `internal/flags/init.go`:

| Opção | Finalidade | Padrão |
| --- | --- | --- |
| `-p` | Nome do pacote | `main` |
| `-t` | Nome de template integrado, caminho local ou URL | `vanilla` |
| `-n` | Nome do projeto | (vazio) |
| `-d` | Diretório do projeto | `.` |
| `-q` | Suprimir a saída do console | false |
| `-l` | Listar templates | false |
| `-skipgomodtidy` | Não executar `go mod tidy` após a extração | false |
| `-git` | URL do repositório Git a ser inicializado | (vazio) |
| `-mod` | Caminho do módulo Go (derivado de `-git` se não for definido) | (vazio) |
| `-s` | Não exibir o aviso ao usar templates remotos | false |
| `-productname` / `-productdescription` / `-productversion` / `-productcompany` / `-productcopyright` / `-productcomments` / `-productidentifier` | Metadados incorporados aos artefatos de compilação gerados | padrões adequados |

**Não** há um alias longo `-list` (somente `-l`), nem um `--help` específico para cada template.

### Substituições

Os placeholders são diretivas padrão de templates Go — o `.` inicial faz parte do acessador do campo:

| Espaço reservado | Exemplo | Origem |
| --- | --- | --- |
| `{{.ProjectName}}` | `myapp` | Opção `-n` / nome do diretório |
| `{{.ModulePath}}` | `github.com/me/myapp` | Opção `-mod` ou derivado de `-git` |
| `{{.WailsVersion}}` | `v3.0.0-…` | Constante incorporada durante a compilação, proveniente de `internal/version` |
| `{{.ProductName}}`, `{{.ProductDescription}}`, `{{.ProductVersion}}`, `{{.ProductCompany}}`, `{{.ProductCopyright}}`, `{{.ProductComments}}`, `{{.ProductIdentifier}}` | Metadados definidos durante a geração do projeto | opções `-product*` correspondentes |

Se precisar de um novo espaço reservado, adicione um campo aos dados do template em `internal/templates/templates.go` e um campo/opção correspondente em `internal/flags/init.go` (ou defina-o a partir de `internal/commands/init.go`).

### Hook pós-cópia

Depois que `gosod` termina de extrair o template, a CLI executa:

```
go mod tidy
```

a menos que você forneça `-skipgomodtidy`. Não há nenhuma etapa `task deps`.

---

## 3. Como criar um novo template

> Exemplo: adicione um template **Solid**

### 3.1 Pasta e ID

```
internal/templates/solid/
```

O nome da pasta é o ID do template. Mantenha-o em **kebab-case**.

### 3.2 Conjunto mínimo de arquivos

```
solid/
├── template.yaml    # name, description, wailsVersion, typescript (required)
├── frontend/        # Your web project (no node_modules/dist)
│   ├── src/
│   ├── package.json
│   └── vite.config.ts
└── ...              # Any extra Go files the template wants to inject
```

Comece copiando `react` e removendo os arquivos desnecessários. Não se esqueça de criar um `template.yaml` — defina `typescript: true` para um template TypeScript — `base/` é a única pasta que não tem um.

### 3.3 Atualize os espaços reservados

Pesquise e substitua os valores literais de exemplo por diretivas de template do Go, por exemplo:

- `myapp` → `{{.ProjectName}}`
- `github.com/you/myapp` → `{{.ModulePath}}`

### 3.4 Integre

Como `templates.go` percorre o sistema de arquivos incorporado durante a inicialização, geralmente basta adicionar uma nova pasta em `internal/templates/<id>/` — não é necessário fazer nenhuma chamada de registro manual. Se precisar de lógica adicional (validação personalizada, etapas pós-cópia), adicione-a a `templates.Install` em `internal/templates/templates.go`.

### 3.5 Teste

```bash
wails3 init -n demo -t solid
cd demo
wails3 dev
```

Verifique se:

- O servidor de desenvolvimento inicia na porta publicada em `WAILS_VITE_PORT`
- Os bindings gerados aparecem em `frontend/bindings/...`
- O recarregamento automático funciona

---

## 4. Como modificar templates existentes

1. Edite os arquivos em `internal/templates/<id>/`.
2. Recompile a CLI (`cd v3 && go build -o ../wails3 ./cmd/wails3`); a diretiva `//go:embed *` incluirá o novo conteúdo.
3. Atualize as **versões das dependências** em `frontend/package.json` e `Taskfile.yml`.
4. Atualize a descrição no arquivo `template.json` do template se o comportamento mudar.

### Ajustes comuns

| Tarefa | Local |
| --- | --- |
| Alterar a porta do servidor de desenvolvimento | `frontend/vite.config.ts` — leia `WAILS_VITE_PORT` |
| Adicionar variáveis de ambiente | `build/Taskfile.yml` ou `frontend/.env` |
| Substituir o gerenciador de pacotes JavaScript | Troque `npm` → `pnpm`/`bun` em `build/Taskfile.yml` |

---

## 5. Dicas para criação de templates

- **Mantenha o frontend genérico** — evite fazer referência a variáveis globais específicas do Wails; `/wails/runtime.js` é servido pelo servidor de ativos em tempo de execução.
- **Não inclua artefatos compilados** — exclua `node_modules`, `dist` e `.DS_Store` do diretório incorporado (ou aplique `.gitignore` a eles para que nunca sejam incluídos em commits).
- **Documente os pré-requisitos** — versão do Node, ferramentas de CLI adicionais etc., em `template.json` ou `NEXTSTEPS.md`.
- **Evite alterações incompatíveis** — se a reformulação for grande, crie um novo ID de template em vez de modificar um existente.

---

## 6. Solução de problemas

| Sintoma | Causa | Correção |
| --- | --- | --- |
| `unknown template name` | Erro de digitação em `-t` ou template não incorporado | Execute `wails3 init -l` para listar os templates disponíveis |
| Placeholders não substituídos | Foi usado `{{ProjectName}}` em vez de `{{.ProjectName}}` | Adicione o `.` inicial (acesso a campo de template do Go) |
| O servidor de desenvolvimento abre uma página em branco | A configuração do Vite não lê `WAILS_VITE_PORT` | Verifique seu `vite.config.ts` |
| Falha ao compilar o frontend para produção | O caminho `base` do Vite foi esquecido | Defina `base: "./"` em `vite.config.ts` |

---

## 7. Mapa dos principais arquivos-fonte

| Arquivo | Responsabilidade |
| --- | --- |
| `internal/templates/templates.go` | Incorpora o sistema de arquivos dos templates e expõe `Install(options *flags.Init) error`, `GetDefaultTemplates()` e `ValidTemplateName(name)` |
| `internal/templates/<id>/**` | Conteúdo propriamente dito do template |
| `internal/commands/init.go` | Integração da CLI: seleciona o template, preenche os metadados e chama `templates.Install` |
| `internal/commands/generate_template.go` | `wails3 generate template` — utilitário para *exportar* um projeto ativo de volta para um template (útil para atualizações) |
| `internal/flags/init.go` | Definições das opções de `wails3 init` |

---

## 8. Recapitulação

- Os templates ficam em **`internal/templates/`** e são incorporados à CLI por meio de `//go:embed *`.
- `wails3 init -t <id>` extrai o template por meio de `gosod` e executa `go mod tidy` (essa etapa pode ser ignorada com `-skipgomodtidy`).
- Criar um template é tão simples quanto **criar uma pasta**, adicionar arquivos e um `template.json` e usar placeholders no estilo de `{{.ProjectName}}`.
- O sistema é **extensível** e **autocontido** — perfeito para compartilhar stacks personalizados com sua equipe ou com a comunidade.

Bom trabalho com os templates!
