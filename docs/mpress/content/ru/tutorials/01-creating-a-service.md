---
title: "Сервис QR-кодов"
description: "Создайте сервис QR-кодов и познакомьтесь с сервисами Wails"
slug: "tutorials/01-creating-a-service"
sourcePath: "tutorials/01-creating-a-service.md"
---

В Wails **сервис** — это структура Go, содержащая бизнес-логику, которую требуется сделать доступной во фронтенде. Сервисы помогают упорядочить код, объединяя связанную функциональность.

Сервис можно представить как набор методов, которые может вызывать код JavaScript. После генерации привязок каждый публичный метод сервиса становится доступен для вызова из фронтенда.

В этом руководстве мы создадим сервис генерации QR-кодов, чтобы продемонстрировать эти концепции. В итоге вы научитесь создавать сервисы, управлять зависимостями и связывать код Go с фронтендом.

<br/>

@steps
### Создание файла сервиса QR-кодов
Создайте в каталоге приложения новый файл с именем `qrservice.go`:

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

**Что здесь происходит:**

- `QRService` — пустая структура, в которой будут находиться наши методы генерации QR-кодов
- `NewQRService()` — функция-конструктор, создающая новый экземпляр нашего сервиса
- `Generate()` — метод, который принимает текст и размер, а затем возвращает QR-код в виде массива байтов PNG
- Метод возвращает `([]byte, error)` в соответствии с принятым в Go соглашением возвращать ошибку последним значением
- Для непосредственной генерации QR-кода используется пакет `github.com/skip2/go-qrcode`

 <br/>

### Регистрация сервиса
Создать сервис недостаточно — его нужно **зарегистрировать** в приложении Wails, чтобы оно узнало о существовании сервиса и смогло сгенерировать для него привязки.

Регистрация выполняется в `main.go` при создании приложения. Передайте экземпляры сервисов в параметр `Services`:

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

**Что здесь происходит:**

- `application.NewService()` оборачивает сервис, чтобы Wails мог им управлять
- Мы вызываем `NewQRService()`, чтобы создать экземпляр сервиса
- Сервис добавляется в срез `Services` в параметрах приложения
- Теперь Wails просканирует этот сервис и сделает его публичные методы доступными во фронтенде

 <br/>

### Установка зависимостей
В коде мы указали пакет `github.com/skip2/go-qrcode`, но пока его не скачали. Go должен получить сведения об этой зависимости и скачать её в проект.

Выполните следующую команду в терминале из каталога проекта:

```bash
go mod tidy
```

**Что здесь происходит:**

- `go mod tidy` сканирует файлы Go на наличие инструкций импорта
- Команда скачивает все отсутствующие пакеты (например, `go-qrcode`) и добавляет их в `go.mod`
- Она также удаляет все зависимости, которые больше не используются
- Благодаря этому в проекте будет весь код, необходимый для успешной компиляции

В выводе должно появиться сообщение о том, что пакет для работы с QR-кодами скачан и добавлен в проект.

 <br/>

### Генерация привязок
Чтобы вызывать эти методы из фронтенда, необходимо сгенерировать привязки. Для этого выполните `wails generate bindings` в корневом каталоге проекта.

@note{type="info"}
При самом первом запуске этой команды в проекте генератор привязок выполняет тщательный анализ кода и зависимостей. Иногда это может занять немного больше времени, чем ожидалось, однако последующие запуски будут гораздо быстрее.

@end

После выполнения команды в терминале должен появиться примерно следующий вывод:

```bash
 % wails3 generate bindings
 INFO  Processed: 337 Packages, 1 Service, 1 Method, 0 Enums, 0 Models in 740.196125ms.
 INFO  Output directory: /Users/leaanthony/myproject/frontend/bindings
```

Обратите внимание: во фронтенд-каталоге появился новый каталог с именем `bindings`:

```bash
frontend/
└── bindings
    └── changeme
        ├── index.js
        └── qrservice.js
```

@note{type="tip" title="Полезный совет"}
При сборке приложения с помощью `wails3 build` привязки генерируются и обновляются автоматически.

@end

 <br/>

### Как устроены привязки
Рассмотрим сгенерированные привязки в `bindings/changeme/qrservice.js`:

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

Мы видим, что привязки сгенерированы для метода `Generate`. Имена параметров сохранены, как и комментарии. Для метода также сгенерирован JSDoc, предоставляющий среде разработки информацию о типах.

@note{type="info"}
Полностью разбираться в сгенерированных привязках необязательно, но важно понимать принцип их работы.

@end

Привязки предоставляют:

- Функции, эквивалентные методам Go
- Автоматическое преобразование между типами Go и JavaScript
- Асинхронные операции на основе Promise
- Информацию о типах в виде комментариев JSDoc

@note{type="tip" title="TypeScript"}
Генератор привязок также поддерживает создание привязок TypeScript. Для этого выполните `wails3 generate bindings -ts`.

@end

Сгенерированный сервис повторно экспортируется файлом `index.js`:

```js {title="bindings/changeme/index.js"}
 // @ts-check
 // Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
 // This file is automatically generated. DO NOT EDIT

 import * as QRService from "./qrservice.js";
 export {
     QRService
 };
```

После этого к нему можно обращаться по упрощённому пути импорта `./bindings/changeme`, содержащему только путь к пакету Go, без указания имени файла.

@note{type="info"}
Упрощённые пути импорта доступны только при использовании сборщиков фронтенда. Если вы предпочитаете обычный фронтенд без сборщика, потребуется вручную импортировать `index.js` или `qrservice.js`.

@end

 <br/>

### Использование привязок во фронтенде
Теперь наш сервис Go можно вызывать из JavaScript! Благодаря сгенерированным привязкам это просто и типобезопасно.

Обновите `frontend/src/main.js`, чтобы использовать новые привязки:

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

**Что здесь происходит:**

- Мы импортируем `QRService` из сгенерированных привязок
- `QRService.Generate()` вызывает наш метод Go. Он возвращает Promise, поэтому мы используем `await`
- Метод Go возвращает `[]byte`, который Wails автоматически преобразует в строку base64 для JavaScript
- Чтобы отобразить изображение PNG, мы создаём URL данных со строкой base64
- Блок `try/catch` обрабатывает любые ошибки на стороне Go (например, недопустимые входные данные)
- Если код Go возвращает ошибку, Promise отклоняется, и здесь мы перехватываем эту ошибку

Теперь обновите `index.html`, чтобы использовать новые привязки в функции `initializeQRGenerator`:

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

Выполните `wails3 dev`, чтобы запустить сервер разработки. Через несколько секунд приложение должно открыться.

Введите текст и нажмите кнопку «Сгенерировать QR-код». В центре страницы должен появиться QR-код:

![QR-код](/assets/qr1.png)

 <br/>

 <br/>

### Альтернативный подход: обработчик HTTP
К этому моменту мы рассмотрели следующие темы:

- Создание нового сервиса
- Генерация привязок
- Использование привязок в коде фронтенда

**Зачем использовать обработчик HTTP?**

Привязки методов отлично подходят для операций с данными, но для передачи файлов, изображений и других медиаданных можно использовать другой подход. Вместо того чтобы преобразовывать всё в base64 и передавать через привязки, можно превратить сервис в небольшой веб-сервер.

Это полезно, когда:

- Нужно передавать изображения, видео или большие файлы
- Нужно использовать стандартные HTML-теги `<img>` или `<video>` с атрибутами `src`
- Нужен прямой доступ к ресурсам по URL

Если ваш сервис реализует стандартный метод Go `ServeHTTP(w http.ResponseWriter, r *http.Request)`, Wails может сделать его доступным как конечную точку HTTP. Расширим наш сервис QR-кодов, добавив такую поддержку:

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

**Что здесь происходит:**

- `ServeHTTP` — стандартный интерфейс Go для обработки HTTP-запросов
- Мы разбираем параметры запроса из URL (`?text=hello&size=256`)
- Мы вызываем существующий метод `Generate()`, чтобы создать QR-код
- Мы задаём тип содержимого `image/png`, чтобы браузеры распознали изображение
- Мы записываем необработанные байты PNG непосредственно в ответ — base64 не требуется!

Теперь обновите `main.go`, указав маршрут, по которому должен быть доступен сервис QR-кодов:

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

**Что здесь происходит:**

- Мы добавляем `application.ServiceOptions`, чтобы настроить способ публикации сервиса
- `Route: "/qrservice"` делает обработчик HTTP доступным по адресу `/qrservice`
- Теперь любой запрос к `/qrservice?text=hello` будет вызывать наш метод `ServeHTTP`
- Если не задать `Route`, функциональность обработчика HTTP будет отключена

@note{type="info"}
Если явно не задать параметр `Route`, обработчик HTTP будет недоступен из фронтенда.

@end

Наконец, обновите `main.js`, чтобы вместо кодирования в base64 использовать обычный атрибут `src` изображения:

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

**Что здесь происходит:**

- Мы удалили импорт и вызов `await QRService.Generate()`
- Вместо этого мы просто задаём для `img.src` адрес нашей конечной точки HTTP
- `encodeURIComponent()` безопасно экранирует специальные символы в URL
- Когда мы задаём `src`, браузер автоматически отправляет HTTP-запрос GET
- Для изображений это проще и эффективнее — преобразование в base64 не требуется!

После повторного запуска приложения должен появиться тот же QR-код:

![QR-код](/assets/qr1.png)

 <br/>

 <br/>

### Поддержка динамической конфигурации
**Проблема жёстко заданных маршрутов:**

В приведённом выше примере мы жёстко задали маршрут `/qrservice` в коде JavaScript. Это создаёт сильную зависимость между конфигурацией Go и кодом фронтенда.

Если изменить параметр `Route` в `main.go`, но не обновить `main.js`, приложение перестанет работать:

```go {title="main.go" ins="3"}
        // ...
            application.NewServiceWithOptions(NewQRService(), application.ServiceOptions{
                Route: "/services/qr",
            }),
        // ...
```

Жёстко заданные маршруты подходят для простых приложений, но делают код хрупким и усложняют его сопровождение.

**Решение: динамическая конфигурация**

Привязки методов и обработчики HTTP могут работать вместе! С помощью привязок можно сообщить фронтенду, какой маршрут использовать. Это сделает конфигурацию динамической и позволит отказаться от жёстко заданного пути.

Вот как это работает:

1. Метод жизненного цикла `ServiceStartup` выполняется при запуске приложения
2. Мы сохраняем настроенный в параметрах маршрут
3. Мы добавляем метод `URL()`, который фронтенд может вызвать, чтобы получить правильный маршрут
4. Теперь фронтенд запрашивает маршрут у сервиса Go, а не пытается его угадать

Сначала реализуйте интерфейс `ServiceStartup` и добавьте новый метод `URL`:

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

**Что здесь происходит:**

- Мы добавили поле `route` для хранения маршрута, настроенного в `ServiceStartup`
- `ServiceStartup(ctx, options)` вызывается при запуске приложения — здесь мы сохраняем маршрут
- Метод `URL()` формирует полный URL с параметрами запроса
- Если маршрут не настроен (то есть пуст), мы возвращаем ошибку
- `url.QueryEscape()` безопасно кодирует текст для использования в URL
- Этот метод будет доступен фронтенду через привязки

Теперь обновите `main.js`, чтобы вместо жёстко заданного пути использовать метод `URL`:

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

**Что здесь происходит:**

- Мы импортируем `QRService`, чтобы снова использовать привязки
- Вместо того чтобы жёстко задавать `/qrservice`, мы вызываем `await QRService.URL(text, 256)`
- Сервис Go формирует URL с правильным маршрутом и параметрами
- Теперь при изменении маршрута в `main.go` фронтенд автоматически использует новый маршрут
- Больше не нужно вручную синхронизировать конфигурацию Go с кодом фронтенда!

Всё должно работать так же, как в предыдущем примере, но теперь изменение маршрута сервиса в `main.go` не нарушит работу фронтенда.

@note{type="info"}
Если метод Go возвращает ошибку, отличную от nil, промис на стороне JS будет отклонён, а инструкции await сгенерируют исключение.

@end

 <br/>

@end
