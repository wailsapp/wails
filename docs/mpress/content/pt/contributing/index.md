---
title: "Como contribuir"
description: "Contribua com o Wails"
slug: "contributing"
sourcePath: "contributing/index.md"
---

## Boas-vindas, colaboradores!

Contribuições para o Wails são bem-vindas! Seja corrigindo bugs, adicionando recursos ou aprimorando a documentação, agradecemos sua ajuda.

## Formas de contribuir

### 1. Relate problemas

Encontrou um bug? [Abra uma issue](https://github.com/wailsapp/wails/issues/new) com:

- Descrição clara
- Etapas para reproduzir
- Comportamento esperado e comportamento observado
- Informações do sistema
- Exemplos de código

### 2. Aprimore a documentação

PRs de correção são bem-vindos sem a necessidade de uma issue prévia ou de um teste de código com falha.  
Siga [Corrigir a documentação](/contributing/documentation/) para visualizar  
e validar uma alteração com o M-Press.

Melhorias na documentação são sempre bem-vindas:

- Corrija erros de digitação e outros erros
- Adicione exemplos
- Esclareça explicações
- Traduza o conteúdo

### 3. Envie código

Contribua com código por meio de pull requests:

- Correções de bugs
- Novos recursos
- Melhorias de desempenho
- Testes

### 4. Proponha uma melhoria (WEP)

Novas funcionalidades e alterações no comportamento público seguem o processo Wails Enhancement Proposal (WEP). Ele mantém o desenvolvimento de recursos transparente e garante que cada proposta aceita tenha alguém responsável por implementá-la. Não abra uma issue de solicitação de recurso.

1. Opcionalmente, apresente sua ideia na categoria [Ideias](https://github.com/wailsapp/wails/discussions/categories/ideas) do GitHub Discussions ou no [Discord](https://discord.gg/JDdSxwjhGf) para avaliar o interesse.
2. Copie [`v3/wep/WEP_TEMPLATE.md`](https://github.com/wailsapp/wails/blob/master/v3/wep/WEP_TEMPLATE.md) para `v3/wep/proposals/<proposal name>/proposal.md` e preencha todas as seções.
3. Abra um pull request em rascunho com o título `[WEP] <title>` contendo apenas a proposta. O PR é o local oficial para discuti-la.
4. Reúna feedback e apoio (comentários e reações de positivo no PR). Reserve pelo menos duas semanas para a discussão e chegue a um acordo sobre quem implementará a proposta.
5. Marque o PR como pronto para revisão. Os mantenedores tomam a decisão final: as propostas aceitas recebem um número WEP e passam por merge.

O processo completo está documentado em [`v3/wep/README.md`](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md).

## Primeiros passos

### Faça um fork e clone

```bash
# Fork the repository on GitHub
# Then clone your fork
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails

# Add upstream remote
git remote add upstream https://github.com/wailsapp/wails.git
```

### Compile a partir do código-fonte

```bash
# Build the v3 CLI. Go downloads any required modules automatically.
cd v3
go build -o ../wails3 ./cmd/wails3

# Confirm the built CLI runs and reports its version.
../wails3 version
```

### Execute os testes

Os testes fazem parte da alteração, não são uma verificação rápida final. Mantenha os testes unitários junto ao código que eles exercitam e prefira testes orientados por tabelas sempre que um comportamento for verificado com várias entradas ou casos extremos. Dê um nome a cada caso para que uma falha explique o cenário.

A lógica nova e alterada deve ter cobertura completa. Procure alcançar 100% de cobertura de instruções Go no código que seu PR adiciona ou altera; não use uma porcentagem referente ao repositório inteiro como substituto para testar a alteração. Uma lacuna pode ser válida — por exemplo, um caminho de erro exclusivo de um sistema operacional ou uma condição cuja reprodução seja inviável sem hardware real —, mas explique a lacuna e por que não é razoavelmente possível testá-la na descrição do PR.

```bash
# Run all v3 tests
cd v3
go test ./...

# Run specific package tests
go test ./pkg/application

# Inspect coverage for the packages you changed
go test ./pkg/application -coverprofile=coverage.out
go tool cover -func=coverage.out
cd ..
```

Para suítes de integração, detecção de condições de corrida e todos os comandos equivalentes aos da CI, consulte [Testes e integração contínua](/contributing/testing-ci/).

## Como fazer alterações

### Crie uma branch

```bash
# Update master
git checkout master
git pull upstream master

# Create feature branch
git checkout -b feature/my-feature
```

### Faça suas alterações

1. **Escreva código** seguindo as convenções do Go
2. **Adicione testes** para novas funcionalidades
3. **Atualize a documentação** se necessário
4. **Execute os testes** para garantir que nada deixe de funcionar
5. **Faça o commit das alterações** com mensagens claras

### Diretrizes para commits

```bash
# Good commit messages
git commit -m "fix: resolve window focus issue on macOS"
git commit -m "feat: add support for custom window chrome"
git commit -m "docs: improve bindings documentation"

# Use conventional commits:
# - feat: New feature
# - fix: Bug fix
# - docs: Documentation
# - test: Tests
# - refactor: Code refactoring
# - chore: Maintenance
```

### Envie um pull request

```bash
# Push to your fork
git push origin feature/my-feature

# Open pull request on GitHub
# Provide clear description
# Reference related issues
```

## Diretrizes para pull requests

### Boa descrição de PR

```markdown
## Description
Brief description of changes

## Changes
- Added feature X
- Fixed bug Y
- Updated documentation

## Testing
- Tested on macOS 14
- Tested on Windows 11
- All tests passing

## Related Issues
Fixes #123
```

### Checklist do PR

- [ ] O código segue as convenções do Go
- [ ] Testes adicionados/atualizados
- [ ] Documentação atualizada
- [ ] Todos os testes estão passando
- [ ] Nenhuma alteração incompatível (ou está documentada)
- [ ] Mensagens de commit claras

## Diretrizes de código

### Estilo de código Go

```go
// ✅ Good: Clear, documented, tested
// ProcessData processes the input data and returns the result.
// It returns an error if the data is invalid.
func ProcessData(data string) (string, error) {
    if data == "" {
        return "", errors.New("data cannot be empty")
    }
    
    result := process(data)
    return result, nil
}

// ❌ Bad: No docs, no error handling
func ProcessData(data string) string {
    return process(data)
}
```

### Testes

```go
func TestProcessData(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid input", "test", "processed", false},
        {"empty input", "", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ProcessData(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ProcessData() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("ProcessData() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Documentação

### Como escrever a documentação

A documentação usa o M-Press. Edite os arquivos `.md` em  
`docs/mpress/content/` e, em seguida, visualize e valide as alterações a partir da raiz do repositório:

```bash
go install github.com/leaanthony/mpress/cmd/mpress@v1.0.17
mpress dev
mpress build --strict --no-purge-css
mpress check
# Use `python` on Windows if `python3` is unavailable.
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

### Estilo da documentação

- Use a ortografia do inglês internacional
- Comece pelo problema
- Forneça exemplos funcionais
- Inclua orientações para solução de problemas
- Inclua referências cruzadas a conteúdos relacionados

## Comunidade

### Obtenha ajuda

- **Discord:** [Participe da nossa comunidade](https://discord.gg/JDdSxwjhGf)
- **Discussões do GitHub:** Faça perguntas
- **Issues do GitHub:** Relate bugs

### Código de Conduta

Seja respeitoso, inclusivo e profissional. Estamos todos aqui para criar ótimos softwares juntos. Consulte o [Código de Conduta](https://github.com/wailsapp/wails/blob/master/CODE_OF_CONDUCT.md) para obter detalhes.

## Reconhecimento

As pessoas que contribuem são reconhecidas em:

- Notas de versão
- Lista de colaboradores
- Insights do GitHub

Agradecemos por contribuir com o Wails! 🎉

## Próximos passos

@cards{cols="2"}
◆ Repositório do GitHub
Acesse o repositório do Wails.

[Ver no GitHub →](https://github.com/wailsapp/wails)

---
◆ Comunidade no Discord
Participe da comunidade.

[Entrar no Discord →](https://discord.gg/JDdSxwjhGf)

---
📖 Documentação
Leia a documentação.

[Explorar a documentação →](/quick-start/why-wails/)

@end
