---
title: "Escreva documentação com o M-Press"
description: "Baixe e instale o M-Press, visualize a documentação do Wails e envie uma contribuição para a documentação."
sourcePath: "contributing/documentation.md"
---

O Wails v3 usa o **M-Press** para transformar Markdown neste site de documentação. Você pode corrigir um erro de digitação, melhorar uma explicação, adicionar um exemplo ou guia ou corrigir uma tradução. Correções apenas na documentação não exigem a abertura prévia de uma issue nem um teste de código com falha.

Você precisa do Git para preparar uma contribuição e de uma conta do GitHub para enviar um pull request. O executável disponível para download do M-Press inclui o gerador de documentação e o renderizador de diagramas; a visualização da documentação não requer Go, Node.js, Wails CLI nem dependências de compilação para desktop. O Python 3 é necessário para os scripts adicionais de validação do repositório.

## Comece por uma página da documentação

Abra **Contribuição → Documentação** na barra lateral para acessar este guia. Para trabalhar em uma página específica, selecione o **ícone de edição** no cabeçalho da página e depois **Editar esta página localmente**. Copie o comando correspondente ao seu sistema operacional e execute-o em um terminal.

O comando usa o M-Press se ele estiver instalado; caso contrário, o instalador baixa uma versão oficial e verifica a soma de verificação. Ele prepara uma cópia de trabalho local para a contribuição e abre a página selecionada no M-Press. Antes de executar as verificações do repositório, confira `mpress version` em relação à versão fixada abaixo: o instalador automático pode baixar uma versão mais recente.

Você pode editar o código-fonte no seu próprio editor enquanto a visualização local é recompilada. O assistente também oferece **Revisar alterações**, **Fazer commit das alterações**, **Enviar branch da contribuição** e **Abrir pull request como rascunho** como ações separadas. Antes do envio, siga as orientações abaixo sobre verificações e pull requests. A geração de traduções é opcional e pode exigir um provedor configurado; a edição e a visualização de Markdown não exigem isso.

Para uma configuração manual reproduzível, siga as etapas abaixo.

## Baixe e instale o M-Press

O script de compilação do repositório baixa a **versão estável mais recente do M-Press**. Use a mesma versão localmente: escolha abaixo um método de instalação, execute `mpress version` e compare o resultado com a versão exibida na página da versão.

### Baixe um executável

Abra a página da [versão estável mais recente do M-Press](https://github.com/leaanthony/mpress/releases/latest) ou baixe abaixo o arquivo correspondente. ARM64 inclui Apple Silicon; AMD64 é a arquitetura comum de PCs Intel/AMD.

| Computador | AMD64 / Intel | ARM64 / Apple Silicon |
| --- | --- | --- |
| macOS | [mpress-darwin-amd64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-darwin-amd64.tar.gz) | [mpress-darwin-arm64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-darwin-arm64.tar.gz) |
| Linux | [mpress-linux-amd64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-linux-amd64.tar.gz) | [mpress-linux-arm64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-linux-arm64.tar.gz) |
| Windows | [mpress-windows-amd64.zip](https://github.com/leaanthony/mpress/releases/latest/download/mpress-windows-amd64.zip) | [mpress-windows-arm64.zip](https://github.com/leaanthony/mpress/releases/latest/download/mpress-windows-arm64.zip) |

Baixe [checksums.txt](https://github.com/leaanthony/mpress/releases/latest/download/checksums.txt) da mesma versão. No diretório de downloads, calcule o hash SHA-256 do arquivo e compare-o com a linha referente a esse nome de arquivo em `checksums.txt`:

@tabs
[macOS]
```sh
shasum -a 256 mpress-darwin-arm64.tar.gz
```

[Linux]
```sh
sha256sum mpress-linux-amd64.tar.gz
```

[Windows PowerShell]
```powershell
Get-FileHash .\mpress-windows-amd64.zip -Algorithm SHA256
```
@end

Use o nome real do arquivo baixado se você escolheu outra arquitetura. Depois de confirmar que os hashes coincidem, extraia e instale o executável:

@tabs
[macOS]
No diretório de downloads, usando o arquivo compactado para Apple Silicon:

```sh
mkdir -p mpress-download
tar -xzf mpress-darwin-arm64.tar.gz -C mpress-download
mkdir -p "$HOME/.local/bin"
install -m 755 mpress-download/mpress "$HOME/.local/bin/mpress"
export PATH="$HOME/.local/bin:$PATH"
mpress version
```

Para Macs Intel, substitua `mpress-darwin-amd64.tar.gz`. Adicione a linha `export PATH` a `~/.zshrc` para manter o comando disponível em novos terminais.

[Linux]
No diretório de downloads:

```sh
mkdir -p mpress-download
tar -xzf mpress-linux-amd64.tar.gz -C mpress-download
mkdir -p "$HOME/.local/bin"
install -m 755 mpress-download/mpress "$HOME/.local/bin/mpress"
export PATH="$HOME/.local/bin:$PATH"
mpress version
```

Para ARM64, substitua `mpress-linux-arm64.tar.gz`. Adicione a linha `export PATH` ao arquivo de inicialização do seu shell, como `~/.bashrc`, para manter o comando disponível em novos terminais.

[Windows PowerShell]
No diretório de downloads:

```powershell
Expand-Archive .\mpress-windows-amd64.zip -DestinationPath .\mpress-download
New-Item -ItemType Directory -Force "$HOME\bin"
Copy-Item .\mpress-download\mpress.exe "$HOME\bin\mpress.exe"
$env:Path = "$HOME\bin;$env:Path"
mpress version
```

Para ARM64, substitua `mpress-windows-arm64.zip`. Adicione o diretório `bin` do seu usuário à variável **Path** do usuário nas **Variáveis de Ambiente** do Windows e reabra o terminal para manter o comando disponível em sessões futuras.
@end

### Instale com o Go

Se você já usa Go, pode compilar a versão estável mais recente. Consulte o `go.mod` dessa versão para verificar o conjunto de ferramentas necessário; a seleção automática do conjunto de ferramentas do Go poderá baixá-lo. A versão usada nesta revisão da documentação requer Go 1.26.10 ou posterior.

```sh
go install github.com/leaanthony/mpress/cmd/mpress@latest
mpress version
```

O Go instala o comando em `GOBIN` ou no diretório `bin` dentro de `go env GOPATH` quando `GOBIN` não está definido. Adicione esse diretório ao seu `PATH` e reabra o terminal se `mpress` não for encontrado.

## Obtenha o código-fonte da documentação

[Crie um fork de wailsapp/wails](https://github.com/wailsapp/wails/fork) e execute estes comandos, substituindo `YOUR_USERNAME` pelo seu nome de usuário do GitHub:

```sh
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails
git remote add upstream https://github.com/wailsapp/wails.git
git fetch upstream master
git switch -c docs/my-change upstream/master
```

Isso cria sua branch de documentação a partir da branch **master** do repositório upstream. Use um nome de branch descritivo no lugar de `docs/my-change`.

## Execute a visualização ao vivo

Na raiz do repositório do Wails, onde fica `mpress.yaml`:

```sh
mpress dev
```

Abra [http://127.0.0.1:3000](http://127.0.0.1:3000) no navegador. Para este guia, acesse [http://127.0.0.1:3000/contributing/documentation/](http://127.0.0.1:3000/contributing/documentation/). Mantenha o comando em execução enquanto edita; as alterações salvas no código-fonte recompilam a visualização. Pressione **Ctrl+C** para interrompê-lo antes de executar os comandos de validação abaixo.

O repositório já contém a configuração do M-Press. Não é necessário inicializar um novo site. Para usar outra porta, execute `mpress dev --port 3001` e abra o endereço localhost correspondente.

## Edite ou adicione uma página

As páginas-fonte em inglês ficam em `docs/mpress/content/`. As traduções ficam em pastas de idiomas, como `fr/` e `id/`. As instruções de contribuição do ícone de edição identificam o arquivo-fonte da página atual.

Para uma página existente, edite o Markdown e preserve os metadados YAML no início. Use títulos, parágrafos, listas, links e blocos de código delimitados. Preserve os componentes pareados do M-Press, como `@note` / `@end` e `@tabs` / `@end`. A [documentação de autoria do M-Press](https://m-press.me/authoring/) descreve a sintaxe compatível.

Para um novo guia, crie um arquivo como `docs/mpress/content/guides/my-guide.md`:

```markdown
---
title: "My guide"
description: "The task this guide helps someone complete."
slug: "guides/my-guide"
sourcePath: "guides/my-guide.md"
---

Explain the task, give working steps, and show the expected result.
```

Adicione a página à seção apropriada de `docs/mpress/content/_nav.yaml`, a menos que essa seção já use `autogenerate`. Verifique se o leitor consegue acessá-la pela navegação e pelos links de páginas relacionadas. Novas páginas em inglês também precisam de cópias completas para todas as localidades antes que as verificações de cobertura possam passar; coordene esse trabalho com os mantenedores ou colaboradores de tradução.

A saída gerada fica em `docs/mpress/site/`. Inclua no commit as páginas-fonte e todas as alterações intencionais de navegação ou recursos estáticos; os arquivos gerados do site e os caches locais devem ficar fora do pull request.

## Verifique sua contribuição

Na raiz do repositório:

```sh
python3 docs/mpress/scripts/check_translations.py
mpress build --strict --no-purge-css
mpress check
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

No Windows, use `python` se foi assim que o Python 3 foi instalado. Se o [Task](https://taskfile.dev/) já estiver instalado, `task docs:build` executará as mesmas verificações com o `mpress` disponível no seu PATH.

A verificação das traduções falha quando há arquivos de tradução ausentes, vazios ou órfãos. Por padrão, trechos traduzidos desatualizados ou corrompidos são relatados como avisos; `--strict` faz com que essas ocorrências também provoquem falha. Uma edição de um texto existente em inglês pode ser enviada enquanto a sincronização das traduções é atualizada. Uma correção em uma tradução pode alterar apenas esse idioma. Mantenha intactos os comandos, nomes de API, links, exemplos de código e conexões dos diagramas; não use marcadores de posição em inglês como conteúdo traduzido.

Leia a página alterada no navegador e execute todos os exemplos de código que você modificar. Se uma verificação falhar em uma página que não foi alterada, inclua a saída e o contexto relevante no seu PR, em vez de tornar a verificação menos rigorosa.

## Envie um pull request

Revise o diff, faça commit das alterações no código-fonte e envie o branch da sua contribuição para o seu fork. Abra um PR para **master** em `wailsapp/wails`. Descreva o que mudou, inclua um link para a página afetada e liste as verificações executadas. Para alterações de layout, inclua capturas de tela de antes e depois; para exemplos de código alterados, informe os detalhes da plataforma e da versão.

O site de produção é estático; a pré-visualização e a validação locais não exigem credenciais de implantação nem um provedor de tradução. Os mantenedores cuidam da publicação após a revisão.

Para contribuições de código e propostas de funcionalidades, consulte [Como contribuir com o Wails](/contributing/). Para conhecer os componentes internos do framework, consulte a [Visão geral técnica](/contributing/overview/).
