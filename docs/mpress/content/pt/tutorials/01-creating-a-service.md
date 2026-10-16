---
title: "Serviço de código QR"
description: "Crie um serviço de código QR para aprender sobre os serviços do Wails"
slug: "tutorials/01-creating-a-service"
sourcePath: "tutorials/01-creating-a-service.md"
---

Um **serviço** no Wails é uma struct Go que contém a lógica de negócios que você deseja disponibilizar para o frontend. Os serviços mantêm seu código organizado agrupando funcionalidades relacionadas.

Pense em um serviço como uma coleção de métodos que seu código JavaScript pode chamar. Cada método público do serviço pode ser chamado pelo frontend após a geração dos bindings.

Neste tutorial, criaremos um serviço gerador de códigos QR para demonstrar esses conceitos. Ao final, você saberá como criar serviços, gerenciar dependências e conectar seu código Go ao frontend.

<br/>

@steps
### Crie o arquivo do serviço de QR
Crie um novo arquivo chamado `qrservice.go` no diretório da sua aplicação:

```go {title="qrservice.go"}
package main

import (
    "github.com/skip2/go-qrcode"
)

// QRService handles QR code generation
type QRService struct {
    // We can add state here if needed
}

// NewQRService creates a new QR service
func NewQRService() *QRService {
    return &QRService{}
}

// Generate creates a QR code from the given text
func (s *QRService) Generate(text string, size int) ([]byte, error) {
    // Generate the QR code
    qr, err := qrcode.New(text, qrcode.Medium)
    if err != nil {
        return nil, err
    }

    // Convert to PNG
    png, err := qr.PNG(size)
    if err != nil {
        return nil, err
    }

    return png, nil
}
```

**O que acontece aqui:**

- `QRService` é uma struct vazia que conterá nossos métodos de geração de códigos QR
- `NewQRService()` é uma função construtora que cria uma nova instância do nosso serviço
- `Generate()` é um método que recebe um texto e um tamanho e retorna o código QR como um array de bytes PNG
- O método retorna `([]byte, error)`, seguindo a convenção do Go de retornar erros como o último valor
- Usamos o pacote `github.com/skip2/go-qrcode` para realizar a geração propriamente dita do código QR

 <br/>

### Registre o serviço
Criar um serviço não é suficiente — precisamos **registrá-lo** na aplicação Wails para que ela saiba que o serviço existe e possa gerar bindings para ele.

O registro ocorre em `main.go` quando você cria sua aplicação. Passe as instâncias do serviço para a opção `Services`:

```go {title="main.go" ins="7-9"}
 func main() {

     app := application.New(application.Options{
         Name:        "myproject",
         Description: "A demo of using raw HTML & CSS",
         LogLevel:    slog.LevelDebug,
         Services: []application.Service{
             application.NewService(NewQRService()),
         },
         Assets: application.AssetOptions{
             Handler: application.AssetFileServerFS(assets),
         },
         Mac: application.MacOptions{
             ApplicationShouldTerminateAfterLastWindowClosed: true,
         },
     })

     app.Window.NewWithOptions(application.WebviewWindowOptions{
         Title:  "myproject",
         Width:  600,
         Height: 400,
     })

     // Run the application. This blocks until the application has been exited.
     err := app.Run()

     // If an error occurred while running the application, log it and exit.
     if err != nil {
         log.Fatal(err)
     }
 }
```

**O que acontece aqui:**

- `application.NewService()` encapsula seu serviço para que o Wails possa gerenciá-lo
- Chamamos `NewQRService()` para criar uma instância do nosso serviço
- O serviço é adicionado ao slice `Services` nas opções da aplicação
- Agora, o Wails examinará esse serviço em busca de métodos públicos para disponibilizá-los ao frontend

 <br/>

### Instale as dependências
Referenciamos o pacote `github.com/skip2/go-qrcode` em nosso código, mas ainda não o baixamos. O Go precisa conhecer essa dependência e baixá-la para o seu projeto.

Execute este comando no terminal a partir do diretório do seu projeto:

```bash
go mod tidy
```

**O que acontece aqui:**

- `go mod tidy` examina seus arquivos Go em busca de declarações de importação
- Ele baixa todos os pacotes ausentes (como `go-qrcode`) e os adiciona a `go.mod`
- Ele também remove todas as dependências que não são mais usadas
- Isso garante que seu projeto tenha todo o código necessário para ser compilado corretamente

Você deverá ver uma saída indicando que o pacote de código QR foi baixado e adicionado ao seu projeto.

 <br/>

### Gere os bindings
Para chamar esses métodos pelo frontend, precisamos gerar bindings. Faça isso executando `wails generate bindings` no diretório raiz do seu projeto.

@note{type="info"}
Na primeira vez que você executar esse comando em um projeto, o gerador de bindings fará uma análise minuciosa do código e das dependências. Às vezes, isso pode demorar um pouco mais do que o esperado; no entanto, as execuções seguintes serão muito mais rápidas.

@end

Após executar o comando, você deverá ver algo semelhante ao seguinte no terminal:

```bash
 % wails3 generate bindings
 INFO  Processed: 337 Packages, 1 Service, 1 Method, 0 Enums, 0 Models in 740.196125ms.
 INFO  Output directory: /Users/leaanthony/myproject/frontend/bindings
```

Observe que há um novo diretório chamado `bindings` no diretório do frontend:

```bash
frontend/
└── bindings
    └── changeme
        ├── index.js
        └── qrservice.js
```

@note{type="tip" title="Dica profissional"}
Quando você compilar sua aplicação usando `wails3 build`, os bindings serão gerados e mantidos atualizados automaticamente.

@end

 <br/>

### Entenda os bindings
Vamos examinar os bindings gerados em `bindings/changeme/qrservice.js`:

```js {title="bindings/changeme/qrservice.js"}
 // @ts-check
 // Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
 // This file is automatically generated. DO NOT EDIT

 /**
  * QRService handles QR code generation
  * @module
  */

 // eslint-disable-next-line @typescript-eslint/ban-ts-comment
 // @ts-ignore: Unused imports
 import {Call as $Call, Create as $Create} from "@wailsio/runtime";

 /**
  * Generate creates a QR code from the given text
  * @param {string} text
  * @param {number} size
  * @returns {Promise<string> & { cancel(): void }}
  */
 export function Generate(text, size) {
     let $resultPromise = /** @type {any} */($Call.ByID(3576998831, text, size));
     let $typingPromise = /** @type {any} */($resultPromise.then(($result) => {
         return $Create.ByteSlice($result);
     }));
     $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
     return $typingPromise;
 }
```

Podemos ver que os bindings são gerados para o método `Generate`. Os nomes dos parâmetros foram preservados, assim como os comentários. Também foi gerado JSDoc para o método, fornecendo informações de tipo à sua IDE.

@note{type="info"}
Não é necessário entender completamente os bindings gerados, mas é importante saber como eles funcionam.

@end

Os bindings fornecem:

- Funções equivalentes aos seus métodos Go
- Conversão automática entre tipos Go e JavaScript
- Operações assíncronas baseadas em Promises
- Informações de tipo na forma de comentários JSDoc

@note{type="tip" title="TypeScript"}
O gerador de bindings também permite gerar bindings TypeScript. Faça isso executando `wails3 generate bindings -ts`.

@end

O serviço gerado é reexportado por um arquivo `index.js`:

```js {title="bindings/changeme/index.js"}
 // @ts-check
 // Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
 // This file is automatically generated. DO NOT EDIT

 import * as QRService from "./qrservice.js";
 export {
     QRService
 };
```

Em seguida, você poderá acessá-lo pelo caminho de importação simplificado `./bindings/changeme`, que consiste apenas no caminho do seu pacote Go, sem especificar nenhum nome de arquivo.

@note{type="info"}
Os caminhos de importação simplificados estão disponíveis apenas ao usar empacotadores de frontend. Se você preferir um frontend básico que não use um empacotador, terá que importar `index.js` ou `qrservice.js` manualmente.

@end

 <br/>

### Use os bindings no frontend
Agora podemos chamar nosso serviço Go pelo JavaScript! Os bindings gerados tornam isso fácil e oferecem segurança de tipos.

Atualize `frontend/src/main.js` para usar os novos bindings:

```js {title="frontend/src/main.js"}
 import { QRService } from './bindings/changeme';

 async function generateQR() {
     const text = document.getElementById('text').value;
     if (!text) {
         alert('Please enter some text');
         return;
     }

     try {
         // Generate QR code as base64
         const qrCodeBase64 = await QRService.Generate(text, 256);

         // Display the QR code
         const qrDiv = document.getElementById('qrcode');
         qrDiv.src = `data:image/png;base64,${qrCodeBase64}`;

     } catch (err) {
         console.error('Failed to generate QR code:', err);
         alert('Failed to generate QR code: ' + err);
     }
 }

 export function initializeQRGenerator() {
     const button = document.getElementById('generateButton');
     button.addEventListener('click', generateQR);
 }
```

**O que acontece aqui:**

- Importamos `QRService` dos bindings gerados
- `QRService.Generate()` chama nosso método Go — ele retorna uma Promise, por isso usamos `await`
- O método Go retorna `[]byte`, que o Wails converte automaticamente em uma string base64 para o JavaScript
- Criamos uma URL de dados com a string base64 para exibir a imagem PNG
- O bloco `try/catch` trata todos os erros provenientes do lado Go (como uma entrada inválida)
- Se nosso código Go retornar um erro, a Promise será rejeitada e o capturaremos aqui

Agora, atualize `index.html` para usar os novos bindings na função `initializeQRGenerator`:

```html {title="frontend/src/index.html"}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>QR Code Generator</title>
            <style>
                body {
                font-family: Arial, sans-serif;
                display: flex;
                flex-direction: column;
                align-items: center;
                justify-content: center;
                height: 100vh;
                margin: 0;
            }
                #qrcode {
                margin-bottom: 20px;
                width: 256px;
                height: 256px;
                display: flex;
                align-items: center;
                justify-content: center;
            }
                #controls {
                display: flex;
                gap: 10px;
            }
                #text {
                padding: 5px;
            }
                #generateButton {
                padding: 5px 10px;
                cursor: pointer;
            }
            </style>
</head>
<body>
<img id="qrcode"/>
<div id="controls">
    <input type="text" id="text" placeholder="Enter text">
        <button id="generateButton">Generate QR Code</button>
</div>

<script type="module">
    import { initializeQRGenerator } from './main.js';
    document.addEventListener('DOMContentLoaded', initializeQRGenerator);
</script>
</body>
</html>
```

Execute `wails3 dev` para iniciar o servidor de desenvolvimento. Após alguns segundos, a aplicação deverá abrir.

Digite um texto e clique no botão "Gerar código QR". Você deverá ver um código QR no centro da página:

![Código QR](/assets/qr1.png)

 <br/>

 <br/>

### Abordagem alternativa: handler HTTP
Até agora, abordamos as seguintes áreas:

- Criação de um novo serviço
- Geração de bindings
- Uso dos bindings no código do nosso frontend

**Por que usar um handler HTTP?**

Os bindings de métodos funcionam muito bem para operações com dados, mas há uma abordagem alternativa para servir arquivos, imagens ou outras mídias. Em vez de converter tudo para base64 e enviar pelos bindings, você pode fazer seu serviço funcionar como um minisservidor web.

Isso é útil quando:

- Você serve imagens, vídeos ou arquivos grandes
- Você quer usar tags HTML `<img>` ou `<video>` padrão com atributos `src`
- Você precisa acessar recursos diretamente por URL

Se o seu serviço implementar o método `ServeHTTP(w http.ResponseWriter, r *http.Request)` padrão do Go, o Wails poderá disponibilizá-lo como um endpoint HTTP. Vamos estender nosso serviço de código QR para oferecer esse suporte:

```go {title="qrservice.go" ins="4-5,37-65"}
package main

import (
    "net/http"
    "strconv"

    "github.com/skip2/go-qrcode"
)

// QRService handles QR code generation
type QRService struct {
    // We can add state here if needed
}

// NewQRService creates a new QR service
func NewQRService() *QRService {
    return &QRService{}
}

// Generate creates a QR code from the given text
func (s *QRService) Generate(text string, size int) ([]byte, error) {
    // Generate the QR code
    qr, err := qrcode.New(text, qrcode.Medium)
    if err != nil {
        return nil, err
    }

    // Convert to PNG
    png, err := qr.PNG(size)
    if err != nil {
        return nil, err
    }

    return png, nil
}

func (s *QRService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Extract the text parameter from the request
    text := r.URL.Query().Get("text")
    if text == "" {
        http.Error(w, "Missing 'text' parameter", http.StatusBadRequest)
        return
    }
    // Extract Size parameter from the request
    sizeText := r.URL.Query().Get("size")
    if sizeText == "" {
        sizeText = "256"
    }
    size, err := strconv.Atoi(sizeText)
    if err != nil {
        http.Error(w, "Invalid 'size' parameter", http.StatusBadRequest)
        return
    }

    // Generate the QR code
    qrCodeData, err := s.Generate(text, size)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Write the QR code data to the response
    w.Header().Set("Content-Type", "image/png")
    w.Write(qrCodeData)
}
```

**O que acontece aqui:**

- `ServeHTTP` é a interface padrão do Go para tratar requisições HTTP
- Analisamos os parâmetros de consulta da URL (`?text=hello&size=256`)
- Chamamos nosso método `Generate()` existente para criar o código QR
- Definimos o tipo de conteúdo como `image/png` para que os navegadores saibam que se trata de uma imagem
- Gravamos os bytes brutos do PNG diretamente na resposta — sem precisar de base64!

Agora, atualize `main.go` para especificar a rota pela qual o serviço de código QR deve ficar acessível:

```go {title="main.go" ins="8-10"}
 func main() {

     app := application.New(application.Options{
         Name:        "myproject",
         Description: "A demo of using raw HTML & CSS",
         LogLevel:    slog.LevelDebug,
         Services: []application.Service{
             application.NewServiceWithOptions(NewQRService(), application.ServiceOptions{
                 Route: "/qrservice",
             }),
         },
         Assets: application.AssetOptions{
             Handler: application.AssetFileServerFS(assets),
         },
         Mac: application.MacOptions{
             ApplicationShouldTerminateAfterLastWindowClosed: true,
         },
     })

     app.Window.NewWithOptions(application.WebviewWindowOptions{
         Title:  "myproject",
         Width:  600,
         Height: 400,
     })

     // Run the application. This blocks until the application has been exited.
     err := app.Run()

     // If an error occurred while running the application, log it and exit.
     if err != nil {
         log.Fatal(err)
     }
 }
```

**O que acontece aqui:**

- Adicionamos `application.ServiceOptions` para configurar como o serviço é exposto
- `Route: "/qrservice"` disponibiliza o manipulador HTTP em `/qrservice`
- Agora, qualquer requisição para `/qrservice?text=hello` chamará nosso método `ServeHTTP`
- Sem definir `Route`, a funcionalidade do manipulador HTTP fica desativada

@note{type="info"}
Se você não definir explicitamente a opção `Route`, o manipulador HTTP não ficará acessível pelo frontend.

@end

Por fim, atualize `main.js` para usar um simples atributo `src` de imagem em vez da codificação base64:

```js {title="frontend/src/main.js"}
async function generateQR() {
    const text = document.getElementById('text').value;
    if (!text) {
        alert('Please enter some text');
        return;
    }

    const img = document.getElementById('qrcode');
    // Make the image source the path to the QR code service, passing the text
    img.src = `/qrservice?text=${encodeURIComponent(text)}`
}

export function initializeQRGenerator() {
    const button = document.getElementById('generateButton');
    if (button) {
        button.addEventListener('click', generateQR);
    } else {
        console.error('Generate button not found');
    }
}
```

**O que acontece aqui:**

- Removemos a importação e a chamada a `await QRService.Generate()`
- Em vez disso, simplesmente definimos `img.src` para apontar para nosso endpoint HTTP
- `encodeURIComponent()` escapa com segurança os caracteres especiais na URL
- O navegador faz automaticamente uma requisição HTTP GET quando definimos `src`
- Isso é mais simples e eficiente para imagens — sem precisar converter para base64!

Ao executar o aplicativo novamente, você deverá obter o mesmo código QR:

![Código QR](/assets/qr1.png)

 <br/>

 <br/>

### Suporte a configurações dinâmicas
**O problema das rotas fixas no código:**

No exemplo acima, usamos uma rota `/qrservice` fixa no código JavaScript. Isso cria um forte acoplamento entre a configuração do Go e o código do frontend.

Se você editar `main.go` e alterar a opção `Route` sem atualizar `main.js`, o aplicativo deixará de funcionar:

```go {title="main.go" ins="3"}
        // ...
            application.NewServiceWithOptions(NewQRService(), application.ServiceOptions{
                Route: "/services/qr",
            }),
        // ...
```

Rotas fixas no código funcionam em aplicativos simples, mas tornam o código frágil e mais difícil de manter.

**A solução: configuração dinâmica**

Bindings de métodos e manipuladores HTTP podem trabalhar juntos! Podemos usar bindings para informar ao frontend qual rota deve ser usada, tornando a configuração dinâmica e eliminando o caminho fixo no código.

Veja como funciona:

1. O método de ciclo de vida `ServiceStartup` é executado quando o aplicativo é iniciado
2. Salvamos a rota configurada nas opções
3. Adicionamos um método `URL()` que o frontend pode chamar para obter a rota correta
4. Agora, o frontend solicita a rota ao serviço Go em vez de tentar adivinhá-la

Primeiro, implemente a interface `ServiceStartup` e adicione um novo método `URL`:

```go {title="qrservice.go" ins="4,6,10,15,23-27,46-55"}
package main

import (
    "context"
    "net/http"
    "net/url"
    "strconv"

    "github.com/skip2/go-qrcode"
    "github.com/wailsapp/wails/v3/pkg/application"
)

// QRService handles QR code generation
type QRService struct {
    route string
}

// NewQRService creates a new QR service
func NewQRService() *QRService {
    return &QRService{}
}

// ServiceStartup runs at application startup.
func (s *QRService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    s.route = options.Route
    return nil
}

// Generate creates a QR code from the given text
func (s *QRService) Generate(text string, size int) ([]byte, error) {
    // Generate the QR code
    qr, err := qrcode.New(text, qrcode.Medium)
    if err != nil {
        return nil, err
    }

    // Convert to PNG
    png, err := qr.PNG(size)
    if err != nil {
        return nil, err
    }

    return png, nil
}

// URL returns an URL that may be used to fetch
// a QR code with the given text and size.
// It returns an error if the HTTP handler is not available.
func (s *QRService) URL(text string, size int) (string, error) {
    if s.route == "" {
        return "", errors.New("http handler unavailable")
    }

    return fmt.Sprintf("%s?text=%s&size=%d", s.route, url.QueryEscape(text), size), nil
}

func (s *QRService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Extract the text parameter from the request
    text := r.URL.Query().Get("text")
    if text == "" {
        http.Error(w, "Missing 'text' parameter", http.StatusBadRequest)
        return
    }
    // Extract Size parameter from the request
    sizeText := r.URL.Query().Get("size")
    if sizeText == "" {
        sizeText = "256"
    }
    size, err := strconv.Atoi(sizeText)
    if err != nil {
        http.Error(w, "Invalid 'size' parameter", http.StatusBadRequest)
        return
    }

    // Generate the QR code
    qrCodeData, err := s.Generate(text, size)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Write the QR code data to the response
    w.Header().Set("Content-Type", "image/png")
    w.Write(qrCodeData)
}
```

**O que acontece aqui:**

- Adicionamos um campo `route` para armazenar a rota configurada em `ServiceStartup`
- `ServiceStartup(ctx, options)` é chamado quando o aplicativo é iniciado — salvamos a rota aqui
- O método `URL()` cria a URL completa com os parâmetros de consulta
- Se nenhuma rota estiver configurada (a rota estiver vazia), retornaremos um erro
- `url.QueryEscape()` codifica o texto com segurança para uso em uma URL
- Esse método ficará disponível para o frontend por meio dos bindings

Agora, atualize `main.js` para usar o método `URL` no lugar de um caminho fixo no código:

```js {title="frontend/src/main.js" ins="1,11-12"}
import { QRService } from "./bindings/changeme";

async function generateQR() {
    const text = document.getElementById('text').value;
    if (!text) {
        alert('Please enter some text');
        return;
    }

    const img = document.getElementById('qrcode');
    // Invoke the URL method to obtain an URL for the given text.
    img.src = await QRService.URL(text, 256);
}

export function initializeQRGenerator() {
    const button = document.getElementById('generateButton');
    if (button) {
        button.addEventListener('click', generateQR);
    } else {
        console.error('Generate button not found');
    }
}
```

**O que acontece aqui:**

- Importamos `QRService` para voltar a usar os bindings
- Em vez de fixar `/qrservice` no código, chamamos `await QRService.URL(text, 256)`
- O serviço Go cria a URL com a rota e os parâmetros corretos
- Agora, se você alterar a rota em `main.go`, o frontend usará automaticamente a nova rota
- Não é mais necessário sincronizar manualmente a configuração do Go com o código do frontend!

O funcionamento deve ser igual ao do exemplo anterior, mas alterar a rota do serviço em `main.go` não fará mais o frontend deixar de funcionar.

@note{type="info"}
Se um método Go retornar um erro diferente de nil, a promise no lado do JS será rejeitada e as instruções await lançarão uma exceção.

@end

 <br/>

@end
