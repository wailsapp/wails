---
title: "Миграция с WebSocket на потоки"
description: "Пошаговый перевод существующей реализации WebSocket на потоки Wails, включая различия, которые приводят к незаметным сбоям"
slug: "guides/streams-from-websockets"
sourcePath: "guides/streams-from-websockets.md"
---

Практическое руководство по переводу приложения, в котором сейчас работает сервер WebSocket, на  
[потоки](/guides/streams/). Инструкции рассчитаны на буквальное выполнение, в том числе агентом.

В результате можно удалить слушатель: больше не будет ни привязанного TCP-порта, ни проверки источника, ни токена — ничего, что могло бы вызвать возражения у брандмауэра или средства защиты конечных точек. API достаточно похож, поэтому большая часть кода фронтенда останется без изменений, однако **три различия приводят к незаметным сбоям**. Они перечислены в первую очередь, поскольку именно из-за них можно потерять вторую половину дня.

## Типичный случай: локальный HTTP-сервер как обходное решение

Обычно WebSocket используется в приложении Wails потому, что раньше не было другого способа передавать непрерывный поток данных во фронтенд. Поэтому приложение запускает собственный `http.Server` на локальном порту, а фронтенд подключается к нему. Если ваша архитектура устроена именно так, при миграции сервер полностью удаляется вместе с несколькими механизмами, созданными *вокруг* него.

**Механизм обнаружения порта больше не нужен.** Фронтенду нужно каким-либо образом сообщать, к какому порту  
подключаться: через привязанный `GetServerPort()`, фиксированный порт с резервным вариантом на случай, если он занят,  
внедрённую глобальную переменную или значение, сохранённое в `localStorage`. Всё это исчезает: к потоку  
обращаются по имени, которое является константой времени компиляции с обеих сторон.

```go
// BEFORE
ln, _ := net.Listen("tcp", "127.0.0.1:0")
go http.Serve(ln, mux)
port := ln.Addr().(*net.TCPAddr).Port     // ...and a binding to hand `port` to the frontend

// AFTER
app.HandleStream("feed", handler)         // that is the entire replacement
```

**Настройка CORS больше не нужна.** В зависимости от платформы источником веб-представления служит `wails://` или  
`http://wails.localhost`, поэтому локальному серверу требуется `CheckOrigin`, заголовок  
`Access-Control-Allow-Origin` либо и то и другое. Потоки используют сервер ресурсов, с которого была загружена страница,  
поэтому разрешать межисточниковый запрос не требуется.

**Любой придуманный вами токен аутентификации больше не нужен.** Порт, привязанный к localhost, доступен любому  
процессу на компьютере, поэтому в тщательно проработанной реализации добавляют токен или одноразовое значение, чтобы другое  
ПО не могло подключиться. Теперь порта, к которому можно подключиться, нет.

**Конечные точки, не относящиеся к WebSocket, переносятся в промежуточное ПО сервера ресурсов.** Такие серверы редко остаются узкоспециализированными: рядом с сокетом обычно постепенно появляются скачивание файлов, конечная точка для изображений и проверка работоспособности. Потоки их не заменяют, но и отдельный сервер для них не требуется. Подключите те же обработчики к серверу ресурсов:

```go
app := application.New(application.Options{
    Assets: application.AssetOptions{
        Handler: yourFrontendAssets,
        Middleware: func(next http.Handler) http.Handler {
            return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                if strings.HasPrefix(r.URL.Path, "/api/") {
                    yourExistingMux.ServeHTTP(w, r)   // the handlers you already wrote
                    return
                }
                next.ServeHTTP(w, r)
            })
        },
    },
})
```

После этого фронтенд обращается к `/api/...` по относительному URL того же источника — без имени хоста, порта и  
CORS. Благодаря потокам и этому изменению у локального сервера больше не остаётся задач.

## Прочтите перед началом

### 1. `ev.data` — это `ArrayBuffer`, а не строка

Это главное различие. WebSocket передаёт текстовые сообщения в виде строк, а поток — все  
сообщения в виде байтов. Такой код **компилируется, запускается и работает неправильно**:

```js
// BEFORE — works with WebSocket
ws.onmessage = (ev) => { const msg = JSON.parse(ev.data); ... };

// AFTER — ev.data is an ArrayBuffer, JSON.parse gets "[object ArrayBuffer]"
```

Исправьте это на границе, а не в каждом месте вызова:

```js
const dec = new TextDecoder();
s.onmessage = (ev) => { const msg = JSON.parse(dec.decode(ev.data)); ... };
```

Если же передаются данные JSON — а обычно так и есть, — используйте `JSONStream` вместо `Stream`  
и полностью обойдите эту проблему:

```js
import { JSONStream } from "@wailsio/runtime";

const s = JSONStream("feed");
s.onmessage = (ev) => dispatch(ev.data);   // already an object
s.send({ hello: true });                   // stringified for you
```

Это самый короткий путь миграции для WebSocket с JSON: замените конструктор, удалите вызовы  
`JSON.parse` и `JSON.stringify`, а остальной код обработчика оставьте без изменений. Для  
данных не в формате JSON создайте одну обёртку и не меняйте обработчики — см.  
[адаптер совместимости](#--1).

### 2. Отправка совместима, приём — нет

`send()` принимает строку и кодирует её в UTF-8, поэтому `s.send(JSON.stringify(x))` работает  
без изменений. Редактировать нужно только код приёма. Эту асимметрию легко упустить,  
поскольку половина кода продолжает работать.

### 3. URL отсутствует

WebSocket передаёт параметры подключения в URL: путь, строку запроса, подпротокол и  
токен аутентификации. У потока есть только имя. Всё, что раньше передавалось в URL, необходимо перенести  
в первый фрейм либо в привязанный метод, вызываемый до подключения.

```js
// BEFORE
const ws = new WebSocket(`wss://host/feed?topic=${topic}&token=${token}`);

// AFTER — no token needed at all; the app is the only possible caller
const s = Stream("feed");
s.onopen = () => s.send(JSON.stringify({ subscribe: topic }));
```

## Сторона Go

Удалите HTTP-сервер, механизм повышения протокола и реестр подключений. Каждый из них заменяется обработчиком.

```go
// BEFORE — gorilla/coder websocket
func (a *App) serveWS(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }
    defer conn.Close()

    clients.add(conn)
    defer clients.remove(conn)

    for {
        _, data, err := conn.ReadMessage()
        if err != nil {
            return
        }
        handle(data)
    }
}

// ...plus http.ListenAndServe, a mux entry, an origin checker, and a token check
```

```go
// AFTER
app.HandleStream("feed", func(c *application.StreamConn) {
    defer c.Close()

    for {
        frame, err := c.Receive()
        if err != nil {
            return                  // reload, close, or shutdown
        }
        handle(frame)
    }
})
```

| WebSocket | Поток |
| --- | --- |
| `upgrader.Upgrade` / `websocket.Accept` | *(ничего — `HandleStream` полностью выполняет регистрацию)* |
| `conn.ReadMessage()` / `conn.Read(ctx)` | `c.Receive()` |
| `conn.WriteMessage(TextMessage, b)` | `c.Send(b)` |
| `conn.Close()` | `c.Close()` или просто возврат из обработчика |
| реестр подключений для широковещательной рассылки | используйте собственный — см. [«Широковещательная рассылка»](#--2) |
| `http.ListenAndServe`, мультиплексор, проверка источника, токен | **удалить** |
| поддержание соединения с помощью ping/pong | **удалить** — нет бездействующего сокета, соединение которого нужно поддерживать |
| `r.Context()` | `c.Context()` |

Время жизни горутины обработчика совпадает со временем жизни подключения, как и у обработчика gorilla,  
поэтому структура существующего цикла остаётся без изменений.

## Сторона фронтенда

```js
// BEFORE
const ws = new WebSocket(url);
ws.onopen    = () => ws.send(JSON.stringify(hello));
ws.onmessage = (ev) => dispatch(JSON.parse(ev.data));
ws.onclose   = () => scheduleReconnect();

// AFTER
import { Stream } from "@wailsio/runtime";
const dec = new TextDecoder();

const s = Stream("feed");
s.onopen    = () => s.send(JSON.stringify(hello));   // unchanged
s.onmessage = (ev) => dispatch(JSON.parse(dec.decode(ev.data)));
s.onclose   = () => scheduleReconnect();             // unchanged
```

`readyState`, четыре константы состояния, `addEventListener`, `close(code, reason)` и  
`bufferedAmount` ведут себя так же, как в `WebSocket`.

### Адаптер совместимости

Если вы предпочитаете вообще не менять обработчики, один раз оберните конструктор. Тогда существующий код, рассчитанный на строковые сообщения, будет работать без изменений:

```js
import { Stream } from "@wailsio/runtime";

/** A Stream that delivers text messages as strings, like a WebSocket. */
export function TextStream(name) {
    const s = Stream(name);
    const dec = new TextDecoder();
    s.binaryType = "arraybuffer";

    const add = s.addEventListener.bind(s);
    const remove = s.removeEventListener.bind(s);
    const wrappers = new WeakMap();
    const decoded = new WeakMap();

    const decodeEvent = (ev) => {
        if (decoded.has(ev)) return decoded.get(ev);
        const data = typeof ev.data === "string" ? ev.data : dec.decode(ev.data);
        const textEvent = new MessageEvent("message", { data });
        decoded.set(ev, textEvent);
        return textEvent;
    };

    const wrap = (listener) => {
        let wrapper = wrappers.get(listener);
        if (wrapper) return wrapper;
        wrapper = (ev) => {
            const textEvent = decodeEvent(ev);
            if (typeof listener === "function") listener.call(s, textEvent);
            else listener.handleEvent(textEvent);
        };
        wrappers.set(listener, wrapper);
        return wrapper;
    };

    s.addEventListener = (type, listener, options) =>
        add(type, type === "message" && listener ? wrap(listener) : listener, options);
    s.removeEventListener = (type, listener, options) =>
        remove(type, type === "message" && listener ? wrappers.get(listener) ?? listener : listener, options);

    // A WailsSocket implements onmessage through addEventListener, but a native
    // WebSocket uses an internal event-handler slot. Define the property on the
    // instance so both transports pass property handlers through the same
    // decoding wrapper as addEventListener listeners.
    let onmessage = null;
    Object.defineProperty(s, "onmessage", {
        get: () => onmessage,
        set(listener) {
            if (onmessage) s.removeEventListener("message", onmessage);
            onmessage = typeof listener === "function" ? listener : null;
            if (onmessage) s.addEventListener("message", onmessage);
        },
        configurable: true,
        enumerable: true,
    });
    return s;
}
```

После этого `const ws = TextStream("feed")` можно использовать как прямую замену `new WebSocket(url)`.

## Широковещательная рассылка

Сервер WebSocket обычно ведёт реестр подключений, чтобы рассылать сообщения всем получателям. У потоков нет встроенной широковещательной рассылки — сохраните реестр, но храните в нём `*StreamConn` вместо `*websocket.Conn`:

```go
type hub struct {
    mu    sync.Mutex
    conns map[*application.StreamConn]struct{}
}

func (h *hub) add(c *application.StreamConn)    { h.mu.Lock(); h.conns[c] = struct{}{}; h.mu.Unlock() }
func (h *hub) remove(c *application.StreamConn) { h.mu.Lock(); delete(h.conns, c); h.mu.Unlock() }

func (h *hub) broadcast(msg []byte) {
    h.mu.Lock()
    conns := make([]*application.StreamConn, 0, len(h.conns))
    for c := range h.conns {
        conns = append(conns, c)
    }
    h.mu.Unlock()                         // never hold the lock across Send

    for _, c := range conns {
        // TrySend, not Send: one stalled frontend must not block the fan-out.
        _ = c.TrySend(msg)
    }
}

app.HandleStream("feed", func(c *application.StreamConn) {
    h.add(c)
    defer h.remove(c)
    defer c.Close()
    <-c.Context().Done()
})
```

Стоит соблюдать два правила: освобождайте блокировку перед отправкой и отдавайте предпочтение `TrySend` при рассылке нескольким получателям, чтобы один медленный потребитель не задерживал всех остальных клиентов.

## Другой случай: фронтенд подключается непосредственно к брокеру

В приложениях Wails это встречается реже, но об этом стоит знать. Если фронтенд открывает WebSocket **для подключения к брокеру, а не к вашему приложению** — например, использует `nats.ws` для подключения к серверу NATS или MQTT через WebSocket, — поток не будет прямой заменой, поскольку он соединяет фронтенд с *вашим кодом Go*, а не со сторонней системой.

Такая миграция меняет архитектуру — и обычно в лучшую сторону:

```
BEFORE   frontend ──ws──► NATS server            (credentials in the frontend)
AFTER    frontend ──stream──► Go ──►  NATS       (credentials stay in Go)
```

Перенесите клиент брокера в Go, где нативная библиотека лучше браузерной, и предоставьте фронтенду необходимые возможности через поток:

```go
nc, _ := nats.Connect(url, nats.UserCredentials(credsPath))  // creds never reach the frontend

app.HandleStream("nats", func(c *application.StreamConn) {
    defer c.Close()

    var subs []*nats.Subscription
    defer func() {
        for _, s := range subs {
            _ = s.Unsubscribe()
        }
    }()

    for {
        frame, err := c.Receive()
        if err != nil {
            return
        }

        var cmd struct {
            Op      string          `json:"op"`       // "sub" | "pub"
            Subject string          `json:"subject"`
            Data    json.RawMessage `json:"data"`
        }
        if json.Unmarshal(frame, &cmd) != nil {
            continue
        }

        switch cmd.Op {
        case "sub":
            sub, err := nc.Subscribe(cmd.Subject, func(m *nats.Msg) {
                out, _ := json.Marshal(map[string]any{"subject": m.Subject, "data": m.Data})
                // TrySend: a slow frontend must not block the NATS callback.
                _ = c.TrySend(out)
            })
            if err == nil {
                subs = append(subs, sub)
            }
        case "pub":
            _ = nc.Publish(cmd.Subject, cmd.Data)
        }
    }
})
```

Обратите внимание на `TrySend` внутри функции обратного вызова подписки: эта функция выполняется в горутине клиента брокера, и её блокировка приостановит доставку по всем подпискам данного соединения.

Преимущества: учётные данные брокера не попадают во фронтенд, на компьютере не открывается порт WebSocket, а переподключение и увеличение задержки между попытками выполняет зрелый клиент Go, а не браузерный клиент.

## Контрольный список миграции

- [ ] `HandleStream` зарегистрирован для каждой имевшейся конечной точки WebSocket
- [ ] Цикл чтения преобразован: `ReadMessage`/`Read` → `c.Receive()`
- [ ] Операции записи преобразованы: `WriteMessage` → `c.Send()` либо в `TrySend` во всех рассылках нескольким получателям и функциях обратного вызова брокера
- [ ] HTTP-сервер, запись мультиплексора, модуль повышения протокола, проверка источника и токен аутентификации **удалены**
- [ ] Механизм поддержания соединения ping/pong **удалён**
- [ ] Параметры URL перенесены в первый фрейм или связанный метод
- [ ] **Каждое прочитанное значение `ev.data` декодируется** — `new TextDecoder().decode(ev.data)` — либо используется адаптер совместимости
- [ ] Логика переподключения сохранена без изменений (встроенного переподключения намеренно нет)
- [ ] Клиенты брокеров перенесены в Go, если фронтенд подключался к брокеру напрямую
- [ ] Проверено наличие окон `InitialHTML` — они вообще не могут использовать потоки

Что стоит искать с помощью grep при преобразовании:

```
new WebSocket(     ev.data            .onmessage
websocket.Accept   upgrader.Upgrade   ReadMessage
WriteMessage       ListenAndServe     CheckOrigin
```

## Ожидаемые различия в поведении

|  | WebSocket | Поток |
| --- | --- | --- |
| Тип сообщения | текст или двоичные данные | только байты |
| `ev.data` | строка или `Blob`/`ArrayBuffer` | всегда `ArrayBuffer` (если не используется `binaryType = "blob"`) |
| `binaryType` по умолчанию | `"blob"` | `"arraybuffer"` |
| Подпротоколы, `extensions` | согласовываются | не поддерживаются; всегда `""` |
| Параметры подключения | URL и строка запроса | первый фрейм или вызов Go-метода через привязки |
| Аутентификация | токен или файл cookie | не требуется — приложение является единственным вызывающим клиентом |
| Поддержание соединения | ping/pong | не требуется |
| Автоматическое переподключение | нет | нет (так же) |
| Коды закрытия | полный диапазон | `1000` — нормальное закрытие, `1001` — сеанс закрыт, `1002` — несоответствие формата фреймов, `1006` — ошибка |
| Противодавление | буфер сокета ядра | 8 МБ / 256 кадров на окно, после чего `Send` блокируется |
| Несколько соединений, одна конечная точка | да | да |

## После миграции

Проверки, помогающие выявить распространённые ошибки:

1. Несколько раз перезагрузите страницу — при каждой перезагрузке обработчик должен завершаться, а новый — запускаться; обработчики не должны накапливаться.
2. Отправьте в каждом направлении сообщение размером более 512 КБ.
3. Оставьте соединение бездействовать более минуты; передача данных должна возобновиться без повторного подключения.
4. Откройте инструменты разработчика и под нагрузкой приостановите выполнение на точке останова, затем продолжите его — производитель должен заблокироваться, а затем восстановить работу, не теряя данные и не увеличивая потребление ресурсов без ограничений.
