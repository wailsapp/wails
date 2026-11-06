---
title: "Compilações ofuscadas"
description: "Compile seu aplicativo Wails com o Garble para proteger o código-fonte contra engenharia reversa"
slug: "guides/build/obfuscation"
sourcePath: "guides/build/obfuscation.md"
---

O [Garble](https://github.com/burrowers/garble) é uma ferramenta de compilação para Go que substitui `go build` para renomear símbolos, ofuscar constantes e remover informações de depuração do binário resultante. O Wails v3 oferece suporte nativo ao Garble por meio de dois novos comandos.

## Pré-requisitos

- **Go 1.26.2 ou posterior** — exigido pelo Garble v0.16.0
- **Garble v0.16.0**

```bash
go install mvdan.cc/garble@v0.16.0
```

@note{type="tip"}
A versão mínima do Go exigida pelo Garble muda entre as versões. Se você usa uma cadeia de ferramentas Go mais antiga, antes de instalar consulte a [página de versões do Garble](https://github.com/burrowers/garble/releases) para encontrar uma versão compatível com sua cadeia de ferramentas.

@end

## Obrigatório: adicione tags JSON aos tipos dos seus serviços

Toda struct retornada ou aceita por um método de serviço vinculado deve ter tags JSON explícitas em todos os campos exportados:

```go
// Without tags — breaks under Garble
type OrderSummary struct {
    ID        int
    Total     float64
    LineItems []LineItem
}

// With tags — safe under Garble
type OrderSummary struct {
    ID        int       `json:"id"`
    Total     float64   `json:"total"`
    LineItems []LineItem `json:"lineItems"`
}
```

@note{type="caution"}
O Garble renomeia os campos exportados das structs, e o Wails passa essas structs para `json.Marshal` por meio de um parâmetro `interface{}` que o Garble não consegue rastrear estaticamente. Uma compilação ofuscada sem tags JSON será concluída com sucesso, mas o frontend receberá nomes de campos adulterados ou vazios em tempo de execução. Adicione as tags antes de executar uma compilação ofuscada.

@end

Os próprios tipos do Wails — `Screen`, `Rect`, `Point`, `Size`, `EnvironmentInfo`, `OSInfo`, `Capabilities` — já têm tags. Você só precisa adicionar tags aos seus próprios tipos.

## Compilação com ofuscação

@steps
### Gere o arquivo de IDs estáveis
Execute este comando sempre que adicionar, renomear ou remover um método de serviço vinculado:

```bash
wails3 generate bindings -obfuscated
```

Isso cria `wails_obfuscated.gen.go` no diretório do pacote principal — faça commit desse arquivo.

### Compile com o Garble
```bash
wails3 build --obfuscated
```

Compila o aplicativo usando os bindings ofuscados.

@end

## Como passar flags adicionais ao Garble

Use `--garbleargs` para encaminhar opções diretamente a `garble`:

```bash
# Obfuscate string literals and reduce binary size
wails3 build --obfuscated --garbleargs "-literals -tiny"

# Reproducible output — same seed produces the same binary
wails3 build --obfuscated --garbleargs "-seed=deadbeef"
```

Consulte a [documentação do Garble](https://github.com/burrowers/garble#flags) para ver a lista completa de flags compatíveis.

## Avançado: como gravar o arquivo de IDs em outro pacote

Por padrão, `wails_obfuscated.gen.go` é gravado ao lado do seu pacote `main`. Se o projeto mantém os serviços em um subpacote importado por `main`, você pode gravar o arquivo nesse local usando `-obfuscated-output`:

```bash
wails3 generate bindings -obfuscated -obfuscated-output ./internal/services
```

@note{type="caution"}
O pacote de destino deve ser importado — direta ou transitivamente — pelo seu pacote `main`, para que seu `init()` seja executado na inicialização. Se o pacote não estiver acessível, os IDs estáveis nunca serão registrados, e as chamadas de binding falharão (por exemplo, com erros `binding not found` em tempo de execução).

@end

## Solução de problemas

### `garble: command not found`

O Garble não está instalado ou `$(go env GOPATH)/bin` não está no seu `PATH`.

```bash
go install mvdan.cc/garble@v0.16.0
export PATH="$PATH:$(go env GOPATH)/bin"
```

### O frontend recebe valores de campos incorretos ou vazios

Os tipos retornados pelos seus serviços não têm tags `json:"..."`. Verifique todas as structs retornadas pelos métodos vinculados e adicione tags explícitas a cada campo exportado.

### Erro `binding not found` no console do navegador

O arquivo de IDs estáveis está ausente ou não foi incluído na compilação. Verifique se:

- `wails_obfuscated.gen.go` existe no diretório do pacote principal (ou no diretório informado a `-obfuscated-output`)
- Você executou `wails3 build --obfuscated`, que adiciona a tag de compilação `wails_obfuscated`
- Se você usou `-obfuscated-output`, o pacote de destino é importado por `main`

### O Windows Defender identifica a compilação como vírus

Os binários Go ofuscados pelo Garble são sinalizados heuristicamente pelo Windows Defender durante a compilação porque não contêm símbolos de depuração e se assemelham a executáveis empacotados. A compilação falha com:

```
open C:\Users\...\AppData\Local\Temp\go-build...\a.out.exe: The file contains a virus or potentially unwanted software.
```

Adicione seu diretório temporário (onde o Go grava os artefatos intermediários da compilação) e o diretório do projeto à lista de exclusões do Defender:

```powershell
Add-MpPreference -ExclusionPath "$env:TEMP"
Add-MpPreference -ExclusionPath "C:\path\to\your\project"
```

Essas exclusões se aplicam apenas aos caminhos especificados e não desativam o Defender globalmente.

### A compilação falha com `unsupported Go version`

O Garble v0.16.0 exige o Go 1.26.2 ou posterior. Atualize o Go ou consulte a [página de versões do Garble](https://github.com/burrowers/garble/releases) para encontrar uma versão compatível com sua cadeia de ferramentas.
