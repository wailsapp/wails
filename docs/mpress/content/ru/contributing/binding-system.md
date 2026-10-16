---
title: "Система привязок"
description: "Как Wails v3 позволяет Go и JavaScript вызывать друг друга без шаблонного кода"
slug: "contributing/binding-system"
sourcePath: "contributing/binding-system.md"
---

> «Привязки» — это **типобезопасный контракт**, который позволяет писать:

```go
msg, err := chatService.Send("Hello")
```

на Go *и*

```ts
import { Send } from "../bindings/github.com/you/yourapp/services/chatservice";

const msg = await Send("Hello");
```

на TypeScript **без необходимости вручную писать связующий код IPC**. В этом документе подробно описано, *как* это происходит: от **статического анализа** во время сборки через **генерацию кода** до **моста среды выполнения**, который передаёт байты через WebView.

> Подробное описание см. в [`contributing/architecture/bindings`](/contributing/architecture/bindings/):
>
> это авторитетное подробное руководство по конвейеру генератора, а данная страница —
>
> обзор для участников проекта.

---

## 1. Обзор за 30 секунд

| Этап | Компонент | Результат |
| --- | --- | --- |
| **Сбор и анализ** | `internal/generator/collect/`, `internal/generator/analyse.go` | Модель экспортируемых сервисов Go, методов, параметров, возвращаемых типов и моделей в памяти |
| **Генерация** | `internal/generator/render/templates/*.tmpl` (`service.{js,ts}.tmpl`, `models.{js,ts}.tmpl`, `index.tmpl`, `eventcreate.js.tmpl`, `eventdata.d.ts.tmpl`, `newline.tmpl`) | Модули ES для каждого сервиса в `frontend/bindings/<full Go import path>/...` |
| **Среда выполнения** | `pkg/application/messageprocessor*.go` и встроенная среда выполнения JS в `internal/runtime/desktop/@wailsio/runtime/src/` (`calls.ts`, `events.ts`, …) | Сообщения о вызовах и событиях, передаваемые через нативный мост WebView |

Потоком управляет команда `wails3 generate bindings`, которая запускает `generator.Generate` (определённый в `internal/generator/generate.go`) для набора пакетов Go.

```
wails3 generate bindings
        │
        ▼
internal/generator/generate.go: Generator.Generate(patterns...)
        │
        ├── internal/generator/collect/   // load.go, collector.go, service.go, model.go, …
        ├── internal/generator/analyse.go // semantic checks
        └── internal/generator/render/    // template execution → frontend/bindings/**
```

---

## 2. Статический анализ

### Точка входа

```
internal/generator/generate.go        // Generator + Generate(patterns…)
internal/generator/analyse.go         // semantic validation
internal/generator/collect/load.go    // go/packages loader
internal/generator/collect/collector.go
```

Проход сборщика обходит каждый загруженный пакет и записывает:

- `collect.ServiceInfo` — по одному для каждой экспортируемой структуры Go с привязкой.
- `collect.ServiceMethodInfo` / `collect.MethodInfo` — сведения о сигнатуре каждого метода (имя, параметры, результаты, позиция ошибки, получатель, документация).
- `collect.ModelInfo` / `collect.StructInfo` — генерируются как модели TS / JS.
- Комментарии-директивы, такие как `//wails:inject`, `//wails:include`, `//wails:internal`, `//wails:ignore`, `//wails:id <hex>` (см. `internal/generator/collect/directive.go`).

Неподдерживаемые типы вызывают ошибку генератора, поэтому ошибки обнаруживаются во время сборки, а не во время выполнения.

### Идентификаторы моделей

Конверт вызова среды выполнения идентифицирует метод по **детерминированному хешу FNV-1a** его полного имени (`pkg.Struct.Method`). В сгенерированных привязках он представлен как `$Call.ByID(<numeric-id>, …)`, а при запуске генерации с `-names` — как `$Call.ByName("pkg.Struct.Method", …)`.

---

## 3. Генерация кода

### Шаблоны

`internal/generator/render/templates/`:

| Шаблон | Назначение |
| --- | --- |
| `service.js.tmpl` | Отдельный модуль JS для каждого сервиса с привязкой |
| `service.ts.tmpl` | Сопутствующий код TypeScript (с `-ts`) |
| `models.js.tmpl` | Выходной класс модели (для каждого пакета) |
| `models.ts.tmpl` | Генерируемые объявления моделей в формате `.d.ts` (для каждого пакета) |
| `index.tmpl` | Сводный реэкспорт `index.{js,ts}` для каждого пакета |
| `eventcreate.js.tmpl` / `eventdata.d.ts.tmpl` | Конструктор события / типы полезной нагрузки |
| `newline.tmpl` | Нормализатор завершающего перевода строки |

Результат помещается в `frontend/bindings/<full Go import path>/...`. Например, сервис, определённый в `github.com/you/yourapp/services/chat`, попадает в `frontend/bindings/github.com/you/yourapp/services/chat/`. В v3 каталога `frontend/src/wailsjs/` нет.

### Выходной код JavaScript

Сгенерированные привязки представляют собой модули ES, которые импортируют вспомогательные функции среды выполнения из `/wails/runtime.js`:

```js
import { Call as $Call, Create as $Create } from "/wails/runtime.js";

/**
 * @param {string} msg
 * @returns {Promise<string> & { cancel(): void }}
 */
export function Send(msg) {
    return $Call.ByID(2042131923, msg);
}
```

При запуске генерации с `-names` вместо этого генерируется `$Call.ByName("pkg.Struct.Method", ...)` — всегда с **полностью квалифицированным** именем метода, а не просто `"Method"`.

Сгенерированные классы моделей используют шаблон конструктора `$$source` со значениями по умолчанию `if (!("X" in $$source))` для каждого поля, заключёнными в кавычки именами полей и `static createFrom(...)`, который применяет `JSON.parse` к строковым входным данным.

### Основные особенности сопоставления типов

Проверено по `internal/generator/render/`:

| Go | TypeScript |
| --- | --- |
| `map[string]V` | `{ [_: string]: V }` |
| `map[K]V` (нестроковый `K`) | `{ [_ in K]?: V }` (не `Map<K, V>` и не `Record<K, V>`) |
| `[]byte` | `Uint8Array` |
| `time.Time` | `string` (ISO 8601 в JSON) |
| `error` (в позиции возвращаемого значения) | отклонённый промис |

### Примечание о рефлексии

`pkg/application/bindings.go` написан **вручную** и использует `reflect` для диспетчеризации методов на основе реестра `BoundMethod`. Не воспринимайте слишком буквально прежние утверждения о «полном отсутствии рефлексии во время выполнения»: генератор избегает рефлексии, но диспетчер среды выполнения её использует.

---

## 4. Протокол вызовов во время выполнения

### Сторона JavaScript

```ts
import { Call } from "/wails/runtime.js";

await Call.ByID(0x7a1201d3 /* ChatService.Send */, "Hello");
// or, with -names:
await Call.ByName("chatservice.ChatService.Send", "Hello");
```

Вспомогательные средства среды выполнения находятся в `internal/runtime/desktop/@wailsio/runtime/src/calls.ts` (диспетчеризация вызовов), `events.ts` (события) и связанных файлах — в этом дереве нет `invoke.ts` или `errors.ts`. Точный формат передаваемого конверта кодируется `calls.ts` на стороне JS и декодируется `pkg/application/messageprocessor_call.go` на стороне Go; при отладке моста рассматривайте эти два файла вместе.

### Сторона Go

1. `pkg/application/messageprocessor_call.go` получает сообщение о вызове.
2. Находит привязанный метод по идентификатору или имени в `pkg/application/bindings.go` (на основе `reflect`).
3. Вызывает привязанный метод и сериализует `{result, error}` для отправки обратно в JS.

### Сопоставление ошибок

| Go | JavaScript |
| --- | --- |
| `error == nil` | `Promise` завершается успешно с результатом |
| `error != nil` | `Promise` отклоняется с объектом `Error`, у которого `message` содержит строку ошибки Go |

---

## 5. Вызов JavaScript из Go

Генератор привязок работает только в одном направлении (предоставляет методы Go для JS). Для взаимодействия из Go с JS используйте шину событий или выполняйте JS в окне:

```go
app.Event.Emit("chat:new-message", msg)
window.ExecJS(`window.dispatchEvent(new CustomEvent("ping"))`)
```

На стороне JS подпишитесь с помощью `Events.On(name, cb)` из `/wails/runtime.js`.

---

## 6. Расширение и устранение неполадок

### Ошибка неподдерживаемого типа

```
error: field "Client" uses unsupported type: chan struct{}
```

→ скройте канал за API метода или пометьте поле с помощью `//wails:internal`, чтобы генератор его пропустил.

### Устаревшие привязки

Сгенерированные файлы перезаписываются при каждом выполнении `wails3 generate bindings`, `wails3 dev` или `wails3 build`. Если автодополнение IDE показывает устаревшие заглушки, удалите `frontend/bindings/` и повторно запустите генератор. Флаг `-clean` (в текущих сборках по умолчанию `true`) очищает каталог привязок перед каждым запуском.

### Советы по производительности

- Не передавайте большие срезы байтов потоково через мост — вместо этого обслуживайте их через сервер ресурсов.
- Если важна задержка, объединяйте несколько быстрых вызовов в один метод.
- Для небольших структур параметров предпочитайте получатели-значения, чтобы сократить количество выделений памяти.

---

## 7. Карта ключевых файлов

| Назначение | Файл |
| --- | --- |
| Оркестрация генератора | `internal/generator/generate.go` |
| Семантические проверки | `internal/generator/analyse.go` |
| Сбор данных (службы, методы, модели) | `internal/generator/collect/{service,method,model,struct,package}.go` |
| Шаблоны рендеринга | `internal/generator/render/templates/*.tmpl` |
| Расположение сгенерированных привязок | `frontend/bindings/<full Go import path>/...` |
| Диспетчер на стороне Go | `pkg/application/bindings.go`, `messageprocessor_call.go` |
| Среда выполнения JS | `internal/runtime/desktop/@wailsio/runtime/src/{calls,events,index}.ts` |

Держите эту памятку под рукой при поиске ошибок в мосте.

---

## 8. Итоги

1. **Сборщик** сканирует ваш код Go → семантическая модель в памяти.
2. **Шаблоны** создают ES-модули для каждой службы и файлы моделей и индексов для каждого пакета.
3. **Обработчик сообщений** диспетчеризует вызовы на стороне Go через реестр привязок.
4. **Среда выполнения JS** предоставляет для всего этого идиоматичные промисы с возможностью отмены.

И всё это — без единой написанной вами строки шаблонного кода IPC. Такова система привязок Wails v3. Вперёд — создавайте привязки!
