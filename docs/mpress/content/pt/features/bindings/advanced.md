---
title: "Vinculação avançada"
description: "Técnicas avançadas de vinculação, incluindo diretivas, injeção de código e IDs personalizados"
slug: "features/bindings/advanced"
sourcePath: "features/bindings/advanced.md"
---

Este guia aborda técnicas avançadas para personalizar e otimizar o processo de geração de vinculações no Wails v3.

## Personalização do código gerado com diretivas

### Injeção de código personalizado

A diretiva `//wails:inject` permite injetar código JavaScript/TypeScript personalizado nas vinculações geradas:

```go
//wails:inject console.log("Hello from Wails!");
type MyService struct {}

func (s *MyService) Greet(name string) string {
    return "Hello, " + name
}
```

Isso injetará o código especificado no arquivo JavaScript/TypeScript gerado para o serviço `MyService`.

Você também pode usar a injeção condicional para direcionar o código a formatos de saída específicos:

```go
//wails:inject j*:console.log("Hello JS!");  // JavaScript only
//wails:inject t*:console.log("Hello TS!");  // TypeScript only
```

### Inclusão de arquivos adicionais

A diretiva `//wails:include` permite incluir arquivos adicionais com as vinculações geradas:

```go
//wails:include js/*.js
package mypackage
```

Essa diretiva geralmente é usada nos comentários de documentação do pacote para incluir arquivos JavaScript/TypeScript adicionais com as vinculações geradas.

### Marcação de tipos e métodos internos

A diretiva `//wails:internal` marca um tipo ou método como interno, impedindo que ele seja exportado para o frontend:

```go
//wails:internal
type InternalModel struct {
    Field string
}

//wails:internal
func (s *MyService) InternalMethod() {}
```

Isso é útil para tipos e métodos usados apenas internamente pelo seu código Go e que não devem ser expostos ao frontend.

### Como ignorar métodos

A diretiva `//wails:ignore` faz com que um método seja completamente ignorado durante a geração de vinculações:

```go
//wails:ignore
func (s *MyService) IgnoredMethod() {}
```

Isso é semelhante a `//wails:internal`, mas ignora completamente o método em vez de marcá-lo como interno.

### IDs de método personalizados

A diretiva `//wails:id` especifica um ID personalizado para um método, substituindo o ID padrão baseado em hash:

```go
//wails:id 42
func (s *MyService) CustomIDMethod() {}
```

Isso pode ser útil para manter a compatibilidade ao refatorar o código.

## Como trabalhar com tipos complexos

### Structs aninhadas

O gerador de vinculações processa structs aninhadas automaticamente:

```go
type Address struct {
    Street string
    City   string
    State  string
    Zip    string
}

type Person struct {
    Name    string
    Address Address
}

func (s *MyService) GetPerson() Person {
    return Person{
        Name: "John Doe",
        Address: Address{
            Street: "123 Main St",
            City:   "Anytown",
            State:  "CA",
            Zip:    "12345",
        },
    }
}
```

O código JavaScript/TypeScript gerado incluirá classes para `Person` e `Address`.

### Maps e slices

Maps e slices também são processados automaticamente:

```go
type Person struct {
    Name       string
    Attributes map[string]string
    Friends    []string
}

func (s *MyService) GetPerson() Person {
    return Person{
        Name: "John Doe",
        Attributes: map[string]string{
            "hair": "brown",
            "eyes": "blue",
        },
        Friends: []string{"Jane", "Bob", "Alice"},
    }
}
```

Em JavaScript, maps são representados como objetos e slices, como arrays. Em TypeScript, maps são representados como `Record<K, V>` e slices, como `T[]`.

### Tipos genéricos

O gerador de vinculações oferece suporte a tipos genéricos:

```go
type Result[T any] struct {
    Data  T
    Error string
}

func (s *MyService) GetResult() Result[string] {
    return Result[string]{
        Data:  "Hello, World!",
        Error: "",
    }
}
```

O código TypeScript gerado incluirá uma classe genérica para `Result`:

```typescript
export class Result<T> {
    "Data": T;
    "Error": string;

    constructor(source: Partial<Result<T>> = {}) {
        if (!("Data" in source)) {
            this["Data"] = null as any;
        }
        if (!("Error" in source)) {
            this["Error"] = "";
        }

        Object.assign(this, source);
    }

    static createFrom<T>(source: string | object = {}): Result<T> {
        let parsedSource = typeof source === "string" ? JSON.parse(source) : source;
        return new Result<T>(parsedSource as Partial<Result<T>>);
    }
}
```

### Interfaces

O gerador de vinculações pode gerar interfaces TypeScript em vez de classes usando o sinalizador `-i`:

```bash
wails3 generate bindings -ts -i
```

Isso gerará interfaces TypeScript para todos os modelos:

```typescript
export interface Person {
    Name: string;
    Attributes: Record<string, string>;
    Friends: string[];
}
```

## Otimização da geração de vinculações

### Uso de nomes em vez de IDs

Por padrão, o gerador de vinculações usa IDs baseados em hash para chamadas de métodos. Você pode usar o sinalizador `-names` para usar nomes em vez deles:

```bash
wails3 generate bindings -names
```

Isso gerará um código que chama o método vinculado pelo nome **totalmente qualificado** (`<package>.<Service>.<Method>`) e usa nomes de parâmetros posicionais `$0`, `$1`, …:

```javascript
export function Greet($0) {
    let $resultPromise = $Call.ByName("main.GreetService.Greet", $0);
    return $resultPromise;
}
```

Isso pode tornar o código gerado mais legível e fácil de depurar, mas talvez seja um pouco menos eficiente.

### Empacotamento do runtime

Por padrão, o código gerado importa o runtime do Wails do pacote npm `@wailsio/runtime`. Você pode usar o sinalizador `-b` para empacotar o runtime com o código gerado:

```bash
wails3 generate bindings -b
```

Isso incluirá o código do runtime diretamente nos arquivos gerados, eliminando a necessidade do pacote npm.

### Desativação de arquivos de índice

Se você não precisar dos arquivos de índice, poderá usar o sinalizador `-noindex` para desativar a geração deles:

```bash
wails3 generate bindings -noindex
```

Isso pode ser útil se você preferir importar serviços e modelos diretamente dos respectivos arquivos.

## Exemplos reais

### Serviço de autenticação

Veja um exemplo de serviço de autenticação com diretivas personalizadas:

```go
package auth

//wails:inject console.log("Auth service initialized");
type AuthService struct {
    // Private fields
    users map[string]User
}

type User struct {
    Username string
    Email    string
    Role     string
}

type LoginRequest struct {
    Username string
    Password string
}

type LoginResponse struct {
    Success bool
    User    User
    Token   string
    Error   string
}

// Login authenticates a user
func (s *AuthService) Login(req LoginRequest) LoginResponse {
    // Implementation...
}

// GetCurrentUser returns the current user
func (s *AuthService) GetCurrentUser() User {
    // Implementation...
}

// Internal helper method
//wails:internal
func (s *AuthService) validateCredentials(username, password string) bool {
    // Implementation...
}
```

### Serviço de processamento de dados

Veja um exemplo de serviço de processamento de dados com tipos genéricos:

```go
package data

type ProcessingResult[T any] struct {
    Data  T
    Error string
}

type DataService struct {}

// Process processes data and returns a result
func (s *DataService) Process(data string) ProcessingResult[map[string]int] {
    // Implementation...
}

// ProcessBatch processes multiple data items
func (s *DataService) ProcessBatch(data []string) ProcessingResult[[]map[string]int] {
    // Implementation...
}

// Internal helper method
//wails:internal
func (s *DataService) parseData(data string) (map[string]int, error) {
    // Implementation...
}
```

### Injeção condicional de código

Veja um exemplo de injeção condicional de código para diferentes formatos de saída:

```go
//wails:inject j*:/**
//wails:inject j*: * @param {string} arg
//wails:inject j*: * @returns {Promise<void>}
//wails:inject j*: */
//wails:inject j*:export async function CustomMethod(arg) {
//wails:inject t*:export async function CustomMethod(arg: string): Promise<void> {
//wails:inject     await InternalMethod("Hello " + arg + "!");
//wails:inject }
type Service struct{}
```

Isso injeta códigos diferentes nas saídas JavaScript e TypeScript, fornecendo anotações de tipo adequadas para cada linguagem.
