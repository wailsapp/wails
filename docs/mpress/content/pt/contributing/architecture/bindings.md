---
title: "Sistema de bindings"
description: "Como o sistema de bindings coleta e processa informações e gera código JavaScript/TypeScript"
slug: "contributing/architecture/bindings"
sourcePath: "contributing/architecture/bindings.md"
---

Este guia explica o funcionamento interno do sistema de bindings do Wails e ajuda os desenvolvedores a entender os mecanismos por trás da geração automática de código.

## Visão geral da arquitetura

O sistema de bindings do Wails consiste em três componentes principais:

1. **Coleta**: analisa o código Go para extrair informações sobre serviços, modelos e outras declarações
2. **Configuração**: gerencia as configurações e opções do processo de geração de bindings
3. **Renderização**: gera código JavaScript/TypeScript com base nas informações coletadas

@filetree
- internal/generator/
  - collect/     # Análise de pacotes e extração de informações
  - config/      # Estruturas e interfaces de configuração
  - render/      # Geração de código para JS/TS
@end

## Processo de coleta

O processo de coleta analisa pacotes Go e extrai informações sobre serviços, modelos e outras declarações. Essa tarefa é realizada pelo pacote `collect`.

### Componentes principais

- **Coletor**: gerencia as informações dos pacotes e armazena em cache os dados coletados
- **Pacote**: representa um pacote Go que está sendo analisado e armazena os serviços, modelos e diretivas coletados
- **Serviço**: coleta informações sobre os tipos de serviço e seus métodos
- **Modelo**: coleta informações detalhadas sobre os tipos de modelo, incluindo campos, valores e parâmetros de tipo
- **Diretiva**: analisa e interpreta diretivas `//wails:` no código-fonte Go

### Fluxo de coleta

1. O coletor examina os pacotes Go especificados no projeto
2. Ele identifica os tipos de serviço (structs com métodos que serão expostos ao frontend)
3. Para cada serviço, ele coleta informações sobre seus métodos
4. Ele identifica os tipos de modelo (structs usados como parâmetros ou valores de retorno nos métodos de serviço)
5. Para cada modelo, ele coleta informações sobre seus campos e parâmetros de tipo
6. Ele processa todas as diretivas `//wails:` encontradas no código

## Processo de renderização

O processo de renderização gera código JavaScript/TypeScript com base nas informações coletadas. Essa tarefa é realizada pelo pacote `render`.

### Componentes principais

- **Renderizador**: coordena a renderização dos arquivos de serviços, modelos e índices
- **Módulo**: representa um único módulo JavaScript/TypeScript gerado
- **Templates**: templates de texto usados na geração de código

### Fluxo de renderização

1. Para cada serviço, o renderizador gera um arquivo JavaScript/TypeScript com funções que correspondem aos métodos do serviço
2. Para cada modelo, o renderizador gera uma classe JavaScript/TypeScript que corresponde à struct do modelo
3. O renderizador gera arquivos de índice que reexportam todos os serviços e modelos
4. O renderizador aplica todas as injeções de código personalizadas especificadas pelas diretivas `//wails:inject`

## Mapeamento de tipos

Um dos aspectos mais importantes do sistema de bindings é o modo como os tipos Go são mapeados para tipos JavaScript/TypeScript. Veja um resumo desse mapeamento:

| Tipo Go | Tipo JavaScript | Tipo TypeScript |
| --- | --- | --- |
| `bool` | `boolean` | `boolean` |
| `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64` | `number` | `number` |
| `string` | `string` | `string` |
| `[]byte` | `Uint8Array` | `Uint8Array` |
| `[]T` | `Array<T>` | `T[]` |
| `map[string]V` | `Object` | `{ [_: string]: V }` |
| `map[K]V` (`K` que não é string) | `Object` | `{ [_ in K]?: V }` |
| `struct` | `Object` | Classe personalizada |
| `interface{}` | `any` | `any` |
| `*T` | `T \| null` | `T \| null` |
| `func` | Não compatível | Não compatível |
| `chan` | Não compatível | Não compatível |

## Sistema de diretivas

O sistema de bindings oferece várias diretivas que podem ser usadas para personalizar o código gerado. Essas diretivas são adicionadas como comentários no código Go.

### Diretivas disponíveis

- `//wails:inject`: injeta código JavaScript/TypeScript personalizado nos bindings gerados
- `//wails:include`: inclui arquivos adicionais com os bindings gerados
- `//wails:internal`: suprime as exportações no SDK gerado; isso não impede a invocação, em tempo de execução, dos métodos Go exportados
- `//wails:ignore`: ignora completamente um método durante a geração de bindings
- `//wails:id`: especifica um ID personalizado para um método, substituindo o ID padrão baseado em hash

Essas diretivas controlam o código gerado. Em builds normais, o registro em tempo de execução ainda detecta os métodos de serviço exportados; use métodos não exportados ou tipos não vinculados para operações privadas e imponha a autorização nos métodos expostos.

### Processamento de diretivas

1. Durante a fase de coleta, o coletor identifica e analisa as diretivas no código Go
2. As diretivas são armazenadas com as declarações correspondentes (serviços, métodos, modelos etc.)
3. Durante a fase de renderização, o renderizador aplica as diretivas para personalizar o código gerado

## Recursos avançados

### Geração condicional de código

O sistema de vinculação permite gerar código condicionalmente usando um prefixo de condição de dois caracteres nas diretivas `include` e `inject`:

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

## Considerações de desempenho

O gerador de vinculações foi projetado para ser eficiente, mas é preciso considerar alguns pontos:

1. A primeira execução será mais lenta, pois criará um cache dos pacotes que serão examinados
2. As execuções seguintes serão mais rápidas, pois usarão as informações armazenadas em cache
3. O gerador processa todos os pacotes do projeto, o que pode levar tempo em projetos grandes
4. Você pode usar a opção `-clean` para limpar o diretório de saída antes da geração

## Depuração

Se você encontrar problemas na geração das vinculações, use a opção `-v` para ativar a saída de depuração:

```bash
wails3 generate bindings -v
```

Isso fornecerá informações detalhadas sobre o processo de coleta e renderização, o que pode ajudar a identificar a origem do problema.
