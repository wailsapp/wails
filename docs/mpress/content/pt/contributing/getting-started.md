---
title: "Primeiros passos"
description: "Como começar a contribuir com o Wails v3"
slug: "contributing/getting-started"
sourcePath: "contributing/getting-started.md"
---

## Boas-vindas, colaborador!

Agradecemos seu interesse em contribuir com o Wails! Este guia ajudará você a fazer sua primeira contribuição.

## Pré-requisitos

Antes de começar, verifique se você tem:

- **Go 1.26+** instalado ([baixar](https://go.dev/dl/))
- **Node.js 20+** e **npm** ([baixar](https://nodejs.org/))
- **Git** configurado com sua conta do GitHub
- Conhecimentos básicos de Go e JavaScript/TypeScript

### Requisitos específicos da plataforma

**macOS:**

- Ferramentas de linha de comando do Xcode: `xcode-select --install`

**Windows:**

- Recomenda-se o MSYS2 ou um ambiente semelhante ao Unix
- Runtime do WebView2 (normalmente pré-instalado no Windows 11)

**Linux:**

- `gcc`, `pkg-config`, `libgtk-4-dev`, `libwebkitgtk-6.0-dev` (pilha GTK4 padrão)
- Instale com: `sudo apt install build-essential pkg-config libgtk-4-dev libwebkitgtk-6.0-dev` (Debian/Ubuntu)
- Para o fluxo de compilação legado `-tags gtk3`, instale também `libgtk-3-dev` e `libwebkit2gtk-4.1-dev`

## Visão geral do processo de contribuição

O fluxo de trabalho típico para contribuições segue estas etapas:

1. **Faça um fork e clone** — Crie sua própria cópia do repositório do Wails
2. **Configure** — Compile a CLI do Wails e verifique seu ambiente
3. **Crie uma branch** — Crie uma branch de funcionalidade para suas alterações
4. **Desenvolva** — Faça suas alterações de acordo com nossos padrões de código
5. **Teste** — Execute os testes para garantir que tudo funcione
6. **Faça o commit** — Faça o commit com mensagens claras e convencionais
7. **Envie** — Abra uma pull request para revisão
8. **Faça os ajustes** — Responda ao feedback e faça os ajustes necessários
9. **Faça o merge** — Após a aprovação, suas alterações passarão a fazer parte do Wails!

## Guia passo a passo

Escolha o tipo de contribuição:

@tabs
[Correção de bug]
@steps
### Encontre ou relate o bug
- Verifique se o bug já foi relatado nas [Issues do GitHub](https://github.com/wailsapp/wails/issues)
- Caso contrário, crie uma nova issue com as etapas para reproduzi-lo
- Aguarde a confirmação antes de começar o trabalho

### Faça um fork e clone
Faça um fork do repositório em [github.com/wailsapp/wails/fork](https://github.com/wailsapp/wails/fork)

Clone seu fork:

```bash
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails
git remote add upstream https://github.com/wailsapp/wails.git
```

### Compile e verifique
Compile o Wails e verifique se você consegue reproduzir o bug:

```bash
cd v3
go build -o ../wails3 ./cmd/wails3

# Reproduce the bug to understand it
```

### Crie uma branch para a correção do bug
Crie uma branch para sua correção:

```bash
git checkout -b fix/issue-123-window-crash
```

### Corrija o bug
- Faça apenas as alterações mínimas necessárias para corrigir o bug
- Não refatore código não relacionado
- Adicione ou atualize testes para evitar regressões

```bash
# Make your changes
# Add tests in *_test.go files
```

### Teste sua correção
Execute os testes para garantir que a correção funcione:

```bash
go test ./...

# Test the specific package
go test ./pkg/application -v

# Run with race detector
go test ./... -race
```

### Faça o commit da correção
Faça o commit com uma mensagem clara:

```bash
git commit -m "fix: prevent window crash when closing during initialization

Fixes #123"
```

### Envie a pull request
Faça o push e crie a PR:

```bash
git push origin fix/issue-123-window-crash
```

Na descrição da sua PR:

- Explique o bug e sua causa raiz
- Descreva sua correção
- Faça referência à issue: "Fixes #123"
- Inclua o comportamento antes e depois da correção

### Responda ao feedback
Atenda aos comentários da revisão e atualize sua PR conforme necessário.

@end

[WEP (Aprimoramento)]
@steps
### Escreva uma WEP
- Leia o [processo de WEP (Wails Enhancement Proposal)](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/README.md)
- Copie o modelo de WEP para `v3/wep/proposals/<name>/proposal.md`
- Abra uma PR de rascunho intitulada `[WEP] <title>` contendo apenas a WEP
- Aguarde a decisão de um mantenedor antes de implementar

### Faça um fork e clone
Faça um fork do repositório em [github.com/wailsapp/wails/fork](https://github.com/wailsapp/wails/fork)

Clone seu fork:

```bash
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails
git remote add upstream https://github.com/wailsapp/wails.git
```

### Configure o ambiente de desenvolvimento
Compile o Wails e verifique seu ambiente:

```bash
cd v3
go build -o ../wails3 ./cmd/wails3

# Run tests to ensure everything works
go test ./...
```

### Crie uma branch de funcionalidade
Crie uma branch com um nome descritivo:

```bash
git checkout -b feat/window-transparency-support
```

### Implemente a funcionalidade
- Siga nossos [padrões de código](/contributing/standards/)
- Mantenha as alterações focadas na funcionalidade
- Escreva código limpo e documentado
- Adicione testes abrangentes

```bash
# Example: Adding a new window method
# 1. Add to window.go interface
# 2. Implement in platform files (darwin, windows, linux)
# 3. Add tests
# 4. Update documentation
```

### Teste minuciosamente
Teste sua funcionalidade:

```bash
# Unit tests
go test ./pkg/application -v

# Integration test - create a test app
cd ..
./wails3 init -n feature-test
cd feature-test
# Add code using your new feature
../wails3 dev
```

### Documente sua funcionalidade
- Adicione comentários de documentação a todas as APIs públicas
- Atualize a documentação pertinente em `/docs/mpress/content/`
- Adicione exemplos, se aplicável

### Siga a convenção nos commits
Use commits convencionais:

```bash
git commit -m "feat: add window transparency support

- Add SetTransparent() method to Window API
- Implement for macOS, Windows, and Linux
- Add tests and documentation

Closes #456"
```

### Envie a pull request
Faça o push e crie a PR:

```bash
git push origin feat/window-transparency-support
```

Na sua PR:

- Descreva a funcionalidade e os casos de uso
- Mostre exemplos ou capturas de tela
- Liste todas as alterações incompatíveis
- Faça referência à PR da WEP aceita

### Faça os ajustes com base na revisão
Os mantenedores podem solicitar alterações. Tenha paciência e colabore.

@end

[Documentação]
PRs de correção são bem-vindos sem a necessidade de abrir uma issue primeiro. Correções apenas na documentação não precisam de um teste de código que falhe. Siga [Escrever documentação com o M-Press](/contributing/documentation/) para conhecer as etapas de instalação do M-Press, os caminhos do código-fonte, a pré-visualização, a validação e o PR.

@end

## Como encontrar issues para trabalhar

- Procure as etiquetas [`good first issue`](https://github.com/wailsapp/wails/labels/good%20first%20issue)
- Confira as issues [`help wanted`](https://github.com/wailsapp/wails/labels/help%20wanted)
- Consulte as [issues abertas](https://github.com/wailsapp/wails/issues) e peça para ser designado

## Como obter ajuda

- **Discord:** entre no [Discord do Wails](https://discord.gg/JDdSxwjhGf)
- **Discussões:** publique no [GitHub Discussions](https://github.com/wailsapp/wails/discussions)
- **Issues:** abra uma issue para relatar um bug reproduzível; use o Discussions para perguntas e um PR de WEP para melhorias

## Código de Conduta

Seja respeitoso, construtivo e acolhedor. Estamos construindo uma comunidade amigável, dedicada a criar ótimos softwares em conjunto.

## Próximas etapas

- Configure seu [ambiente de desenvolvimento](/contributing/setup/)
- Consulte nossos [padrões de codificação](/contributing/standards/)
- Explore a [documentação técnica](/contributing/overview/)
