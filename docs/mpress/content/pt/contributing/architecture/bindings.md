---
title: "Sistema de vinculação"
description: "Como o sistema de vinculação coleta, processa e gera código JavaScript/TypeScript"
slug: "contributing/architecture/bindings"
sourcePath: "contributing/architecture/bindings.md"
---

Este guia explica o funcionamento interno do sistema de vinculação do Wails e oferece informações para desenvolvedores que desejam compreender os mecanismos por trás da geração automática de código.

## Visão geral da arquitetura

O sistema de vinculação do Wails consiste em três componentes principais:

1. **Coleta**: analisa o código Go para extrair informações sobre serviços, modelos e outras declarações
2. **Configuração**: gerencia as configurações e opções do processo de geração de vinculações
3. **Renderização**: gera código JavaScript/TypeScript com base nas informações coletadas

@filetree
- internal/generator/
  - collect/     # Análise de pacotes e extração de informações
  - config/      # Estruturas e interfaces de configuração
  - render/      # Geração de código para JS/TS
@end

## Processo de coleta

O processo de coleta é responsável por analisar pacotes Go e extrair informações sobre serviços, modelos e outras declarações. Essa tarefa é realizada pelo pacote `collect`.

### Componentes principais

- **Collector**: gerencia as informações dos pacotes e armazena em cache os dados coletados
- **Package**: representa um pacote Go em análise e armazena os serviços, modelos e diretivas coletados
- **Service**: coleta informações sobre os tipos de serviço e seus métodos
- **Model**: coleta informações detalhadas sobre os tipos de modelo, incluindo campos, valores e parâmetros de tipo
- **Directive**: analisa e interpreta as diretivas `//wails:` no código-fonte Go

### Fluxo de coleta

1. O coletor examina os pacotes Go especificados no projeto
2. Ele identifica os tipos de serviço (structs com métodos que serão expostos ao frontend)
3. Para cada serviço, ele coleta informações sobre seus métodos
4. Ele identifica os tipos de modelo (structs usados como parâmetros ou valores de retorno nos métodos de serviço)
5. Para cada modelo, ele coleta informações sobre seus campos e parâmetros de tipo
6. Ele processa todas as diretivas `//wails:` encontradas no código

## Processo de renderização

O processo de renderização é responsável por gerar código JavaScript/TypeScript com base nas informações coletadas. Essa tarefa é realizada pelo pacote `render`.

### Componentes principais

- **Renderer**: coordena a renderização dos arquivos de serviços, modelos e índices
- **Module**: representa um único módulo JavaScript/TypeScript gerado
- **Templates**: modelos de texto usados para gerar código

### Fluxo de renderização

1. Para cada serviço, o renderizador gera um arquivo JavaScript/TypeScript com funções que refletem os métodos do serviço
2. Para cada modelo, o renderizador gera uma classe JavaScript/TypeScript que reflete a struct do modelo
3. O renderizador gera arquivos de índice que reexportam todos os serviços e modelos
4. O renderizador aplica todas as injeções de código personalizadas especificadas pelas diretivas `//wails:inject`

## Mapeamento de tipos

Um dos aspectos mais importantes do sistema de vinculação é a forma como os tipos Go são mapeados para tipos JavaScript/TypeScript. Veja um resumo do mapeamento:

| Tipo Go | Tipo JavaScript | Tipo TypeScript |
| --- | --- | --- |
| `bool` | `boolean` | `boolean` |
| `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64` | `number` | `number` |
| `string` | `string` | `string` |
| `[]byte` | `Uint8Array` | `Uint8Array` |
| `[]T` | `Array<T>` | `T[]` |
| `map[string]V` | `Object` | `{ [_: string]: V }` |
| `map[K]V` (`K` que não seja string) | `Object` | `{ [_ in K]?: V }` |
| `struct` | `Object` | Classe personalizada |
| `interface{}` | `any` | `any` |
| `*T` | `T \| null` | `T \| null` |
| `func` | Não compatível | Não compatível |
| `chan` | Não compatível | Não compatível |

## Sistema de diretivas

O sistema de bindings oferece várias diretivas que podem ser usadas para personalizar o código gerado. Essas diretivas são adicionadas como comentários no código Go.

### Diretivas disponíveis

- `//wails:inject`: injeta código JavaScript/TypeScript personalizado nos bindings gerados
- `//wails:include`: inclui arquivos adicionais nos bindings gerados
- `//wails:internal`: marca um tipo ou método como interno, impedindo que ele seja exportado para o frontend
- `//wails:ignore`: ignora completamente um método durante a geração dos bindings
- `//wails:id`: especifica um ID personalizado para um método, substituindo o ID padrão baseado em hash

### Processamento de diretivas

1. Durante a fase de coleta, o coletor identifica e analisa as diretivas no código Go
2. As diretivas são armazenadas com as declarações correspondentes (serviços, métodos, modelos etc.)
3. Durante a fase de renderização, o renderizador aplica as diretivas para personalizar o código gerado

## Recursos avançados

### Geração condicional de código

O sistema de bindings permite a geração condicional de código usando um prefixo de condição de dois caracteres nas diretivas `include` e `inject`:

```
<language><style>:<content>
```

Em que:

- `<language>` pode ser:
  - `*` - JavaScript e TypeScript
  - `j` - Somente JavaScript
  - `t` - Somente TypeScript


- `<style>` pode ser:
  - `*` - Classes e interfaces
  - `c` - Somente classes
  - `i` - Somente interfaces


Por exemplo:

```go
//wails:inject j*:console.log("JavaScript only");
//wails:inject t*:console.log("TypeScript only");
```

### IDs de método personalizados

Por padrão, os métodos são identificados por um ID baseado em hash. No entanto, você pode especificar um ID personalizado usando a diretiva `//wails:id`:

```go
//wails:id 42
func (s *Service) CustomIDMethod() {}
```

Isso pode ser útil para manter a compatibilidade durante a refatoração do código.

## Considerações sobre desempenho

O gerador de bindings foi projetado para ser eficiente, mas há alguns aspectos que devem ser considerados:

1. A primeira execução será mais lenta enquanto o gerador cria um cache dos pacotes a serem verificados
2. As execuções seguintes serão mais rápidas porque usarão as informações armazenadas em cache
3. O gerador processa todos os pacotes do projeto, o que pode ser demorado em projetos grandes
4. Você pode usar a opção `-clean` para limpar o diretório de saída antes da geração

## Depuração

Se você encontrar problemas na geração dos bindings, poderá usar a opção `-v` para ativar a saída de depuração:

```bash
wails3 generate bindings -v
```

Isso fornecerá informações detalhadas sobre o processo de coleta e renderização, o que pode ajudar a identificar a origem do problema.
