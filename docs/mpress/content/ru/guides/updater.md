---
title: "Обновление приложения"
description: "Встроенное самообновление приложений Wails v3: подключаемые провайдеры, криптографическая проверка, атомарная замена и стандартный интерфейс, который можно стилизовать или заменить."
slug: "guides/updater"
sourcePath: "guides/updater.md"
---

Модуль обновления доставляет обновления ПО внутри приложения, избавляя вас от необходимости создавать собственный конвейер скачивания, проверки и замены. Он работает поверх `app.Updater`, принимает один или несколько подключаемых `Provider` (GitHub Releases, keygen.sh, Sparkle AppCast, открытый протокол Wails Update Manifest или собственную реализацию), проверяет подлинность скачанных файлов с помощью настроенного открытого ключа, безопасно заменяет работающий исполняемый файл и передаёт сведения о каждом переходе через стандартную шину событий Wails.

![Стандартное окно обновления в состоянии «Обновление готово»: значок, соответствующий состоянию; плашка версии (v1.0.0 → v2.0.1 · 8.8 МБ); примечания к выпуску, отображаемые из Markdown и включающие таблицу GFM; одно основное действие.](/assets/updater/default-window-ready.png)

## Быстрый старт

```go {title="main.go"}
package main

import (
    "context"
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/updater"
    "github.com/wailsapp/wails/v3/pkg/updater/providers/github"
)

func main() {
    app := application.New(application.Options{Name: "Demo"})

    gh, _ := github.New(github.Config{Repository: "myorg/myapp"})
    if err := app.Updater.Init(updater.Config{
        CurrentVersion: "1.0.0",
        Providers:      []updater.Provider{gh},
    }); err != nil {
        log.Fatal(err)
    }

    if err := app.Updater.CheckAndInstall(context.Background()); err != nil {
        log.Printf("update: %v", err)
    }

    _ = app.Run()
}
```

Этот код открывает окно обновления фреймворка, проверяет GitHub, скачивает файл для текущей платформы, проверяет его, заменяет исполняемый файл и ожидает, пока пользователь перезапустит приложение.

## Жизненный цикл

`app.Updater` — это конечный автомат со следующими состояниями (`updater.State`):

| Состояние | Когда возникает |
| --- | --- |
| `unconfigured` | До вызова `Init` |
| `idle` | После `Init`, до первой проверки |
| `checking` | Выполняется `Check` |
| `up-to-date` | Согласно последнему ответу провайдера, у вызывающей стороны актуальная версия |
| `available` | Найден новый выпуск, но скачивание ещё не началось |
| `downloading` | Байты поступают потоком от провайдера |
| `verifying` | Скачивание завершено, проверяется подпись или дайджест |
| `installing` | Проверенные данные распаковываются и переименовываются в каталоге подготовки |
| `ready` | Обновление подготовлено; для применения вызовите `Restart` |
| `error` | На одном из предыдущих этапов произошла ошибка |

Текущее состояние можно в любой момент получить с помощью `app.Updater.State()`. При каждом переходе также отправляется событие Wails (см. раздел [«События»](#heading-1)).

`Restart` ожидает, пока вспомогательный процесс достигнет состояния `application.New`, и только затем предлагает работающему приложению завершить работу. Стандартный тайм-аут запуска составляет 30 секунд. Если до `application.New` ваше приложение выполняет длительную инициализацию, задайте для `Config.HelperReadyTimeout` более продолжительный интервал, например `time.Minute`. Нулевое значение выбирает стандартный интервал; отрицательные значения отклоняются. Если время запуска истечёт, `Restart` вернёт `updater.ErrHelperNotReady` и оставит работающее приложение открытым.

Обновление подготавливается в каталоге `wails-update-*` рядом с установленным исполняемым файлом (в macOS — рядом с пакетом `.app`), поэтому этот каталог должен быть доступен для записи.

Внутри AppImage модуль обновления заменяет и повторно запускает файл `.AppImage`, указанный переменной `$APPIMAGE`, а не файловую систему, смонтированную только для чтения.

Стандартное окно автоматически отражает текущее состояние. Например, если `Check` сообщает, что обновление не требуется, пользователь видит следующее сообщение и закрывает окно кнопкой **«Закрыть»**:

![Стандартное окно обновления в состоянии «Установлена актуальная версия»: зелёная галочка, заголовок «У вас актуальная версия» и единственная кнопка «Закрыть».](/assets/updater/default-window-up-to-date.png)

## Провайдеры

Провайдером `Provider` может быть любая реализация следующего интерфейса:

```go
type Provider interface {
    Name() string
    Check(ctx context.Context, req CheckRequest) (*Release, error)
    Download(ctx context.Context, r *Release, dst io.Writer, onProgress func(written, total int64)) error
}
```

В исходный код входят четыре реализации.

### GitHub Releases — `updater/providers/github`

```go {title="github provider"}
gh, err := github.New(github.Config{
    Repository:    "myorg/myapp",     // your owner/repo (required)
    Token:         "ghp_…",           // optional; raises rate limit + private repos
    Prerelease:    false,             // include pre-releases in latest lookup
    ChecksumAsset: "SHA256SUMS",      // optional sibling asset for digest verification
    BaseURL:       "",                // optional override (e.g. GitHub Enterprise)
    AssetMatcher:  nil,               // optional custom asset-picker; nil uses DefaultAssetMatcher
    HTTPClient:    nil,               // optional client override
})
```

Стандартный алгоритм сопоставления выбирает ресурс по подстрокам `GOOS` и `GOARCH` в имени файла, распознавая распространённые псевдонимы (`amd64` / `x86_64` / `x64`, `arm64` / `aarch64`, `386` / `i386` / `x86` / `ia32`). Для собственных схем именования:

```go
gh, _ := github.New(github.Config{
    Repository: "myorg/myapp",
    AssetMatcher: func(req updater.CheckRequest, assets []github.ReleaseAsset) int {
        for i, a := range assets {
            if strings.Contains(a.Name, "my-naming-convention") &&
               strings.Contains(a.Name, req.Platform) {
                return i
            }
        }
        return -1 // no match
    },
})
```

`ChecksumAsset` — имя соседнего ресурса выпуска, содержимое которого состоит из строк `<sha256>  <filename>` (в таком формате выводят данные `sha256sum` и `shasum -a 256`). Провайдер получает этот ресурс во время `Check`, находит строку, соответствующую выбранному артефакту, и заполняет `Release.Verification.Digest`, чтобы фреймворк проверил скачанный файл.

### keygen.sh — `updater/providers/keygen`

```go {title="keygen provider"}
kg, err := keygen.New(keygen.Config{
    Account:    "your-account-slug", // required
    Product:    "product-uuid",      // optional but recommended when account has multiple products
    Package:    "",                  // optional further narrowing
    Channel:    "stable",            // "stable" / "rc" / "beta" / "alpha" / "dev"
    Filetype:   "",                  // optional artifact filetype filter ("dmg", "exe", …)
    Token:      "prod-…",            // product / environment / user / admin token; wins over LicenseKey
    LicenseKey: "",                  // license key auth (used only when Token is empty)
    BaseURL:    "",                  // optional API base override
    HTTPClient: nil,                 // optional client override; redirect-strip wrapper still applied
})
```

Провайдер автоматически преобразует контрольную сумму SHA-512 и подпись Ed25519ph каждого артефакта keygen.sh в блок `Release.Verification` фреймворка — дополнительная настройка не требуется.

**Форматы токенов:** токены keygen.sh содержат префикс роли (`admi-` / `prod-` / `envi-` / `user-`). Необработанный UUID, отображаемый на панели управления, — это *идентификатор* токена, а не его секретное значение. Секрет отображается только при создании токена. Подробнее см. в [документации keygen.sh по аутентификации](https://keygen.sh/docs/api/authentication/).

### Sparkle AppCast — `updater/providers/appcast`

```go {title="appcast provider"}
ac, err := appcast.New(appcast.Config{
    URL:        "https://your.app/appcast.xml", // required
    Channel:    "stable",                       // optional sparkle:channel filter
    HTTPClient: nil,                            // optional client override
})
```

Подключается к существующей инфраструктуре Sparkle / WinSparkle без каких-либо изменений. Считывает из ленты `sparkle:shortVersionString`, `<enclosure url type length sparkle:os sparkle:edSignature>` и `sparkle:channel`.

Подписи DSA (`sparkle:dsaSignature`) из Sparkle 1 не поддерживаются. Проектам, использующим эту схему подписи, следует перейти на EdDSA (Sparkle 2).

### Wails Update Manifest — `updater/providers/endpoint`

```go {title="endpoint provider"}
ep, err := endpoint.New(endpoint.Config{
    URL:        "https://updates.example.com/check", // required; supports {{platform}} / {{arch}} / {{version}} / {{channel}} placeholders
    Channel:    "stable",                            // optional channel filter
    Headers:    nil,                                 // optional headers, e.g. {"Authorization": "License <key>"}
    HTTPClient: nil,                                 // optional client override; redirect-strip wrapper still applied
})
```

Использует открытый [протокол Wails Update Manifest](/reference/update-manifest/): один документ JSON описывает последний выпуск и его артефакты для каждой платформы, включая встроенные контрольные суммы и подписи. Один и тот же документ можно разместить на статическом файловом хостинге (S3, GitHub Pages или любом CDN — публикуйте для каждого канала по одному манифесту со всеми платформами) либо на динамическом сервере обновлений (при каждой проверке провайдер отправляет `platform`, `arch`, `version` и `channel`, поэтому сервер может вернуть ровно один артефакт или предоставить его только при наличии лицензии).

Благодаря заполнителям URL статическую структуру можно настроить одной строкой:

```go
ep, _ := endpoint.New(endpoint.Config{
    URL: "https://cdn.example.com/updates/{{platform}}/{{arch}}/stable.json",
})
```

Настроенные заголовки отправляются с каждым запросом манифеста. При скачивании артефактов они повторно используются только на том же хосте, где расположен манифест, и лишь при отсутствии перехода с `https` на `http`; заголовок `Authorization` удаляется при любом перенаправлении на другой источник или протокол с более низким уровнем безопасности.

За публикацию отвечает CLI: команда `wails3 updater manifest` за один вызов вычисляет дайджесты, подписывает и описывает файлы выпуска, а `wails3 updater verify` повторно проверяет результат перед загрузкой. См. раздел [«Публикация с помощью CLI wails3»](/reference/update-manifest/#publishing-with-the-wails3-cli).

### Цепочка резервных провайдеров

Список `Config.Providers` упорядочен. Модуль обновления последовательно обходит его: используется первый провайдер, вернувший выпуск; первый ответ «установлена актуальная версия» немедленно останавливает обход цепочки (резервирование предназначено для ситуации «основной провайдер недоступен», а не «провайдеры расходятся во мнениях»). При ошибке модуль переходит к следующему провайдеру.

```go
app.Updater.Init(updater.Config{
    CurrentVersion: "1.0.0",
    Providers: []updater.Provider{
        kg, // primary: licensed customers
        gh, // fallback: public mirror
    },
})
```

### Создание собственного провайдера

Для типичной реализации достаточно трёх методов и примерно 150 строк. Updater отвечает за проверку, атомарную подготовку, замену и окно, а код провайдера определяет следующий выпуск и передаёт байты потоком:

```go
type CustomProvider struct { /* config */ }

func (p *CustomProvider) Name() string { return "custom" }

func (p *CustomProvider) Check(ctx context.Context, req updater.CheckRequest) (*updater.Release, error) {
    // Hit your update endpoint, decide whether req.CurrentVersion is current,
    // and return either nil (no upgrade) or a *Release with Artifact + optional
    // Verification populated.
    // Errors here drop through to the next provider in Config.Providers.
}

func (p *CustomProvider) Download(ctx context.Context, r *updater.Release, dst io.Writer, onProgress func(int64, int64)) error {
    // Stream the artifact's bytes to dst. Call onProgress(written, total) as
    // bytes flow past; the Updater debounces emits to ~10 Hz on the event bus.
}
```

Ориентируйтесь на встроенные провайдеры — каждый из них реализован в одном файле Go.

## Криптографическая проверка

Средство проверки фреймворка подтверждает подлинность выпусков, используя `Config.PublicKey` как корень доверия:

```go
//go:embed publickey.pem
var publicKey []byte

app.Updater.Init(updater.Config{
    CurrentVersion: "1.0.0",
    Providers:      []updater.Provider{...},
    PublicKey:      publicKey,
})
```

Поддерживаемые алгоритмы (`Release.Verification.SignatureAlgo`):

| Алгоритм | Что подписывается | Примечания |
| --- | --- | --- |
| `ed25519` | Дайджест SHA-256 артефакта | Используется Sparkle EdDSA |
| `ed25519ph` | Весь артефакт через предварительное хеширование Ed25519ph (внутри используется SHA-512) | Используется keygen.sh |
| `ecdsa-p256` | Дайджест SHA-256 артефакта | Принимаются как необработанные подписи `r∥s`, так и подписи DER |

Также поддерживается проверка только по дайджесту (`DigestAlgo`: `sha256` / `sha512`), если релиз содержит хеш, но не содержит подписи.

`Config.PublicKey` — ЕДИНСТВЕННЫЙ корень доверия для проверки подписи: источник релиза никак не может подставить собственный ключ. Для релизов, содержащих `Signature`, но не имеющих настроенного `Config.PublicKey`, проверка завершается безопасным отказом. Верификатор вычисляет дайджест в потоковом режиме во время скачивания, поэтому даже при обновлениях размером в несколько гигабайт проверка не требует дополнительного прохода по диску.

@note{type="caution" title="Проверка только по дайджесту ≠ криптографическая проверка"}
Подлинность релиза, содержащего только `Digest`, подтверждается TLS реестра и теми гарантиями целостности, которые предоставляет сам реестр, а не контролируемым вами криптографическим корнем. Используйте проверку только по дайджесту для обнаружения порчи данных, а подписи — для защиты от подмены при компрометации конвейера выпуска релизов.

@end

### Создание ключа подписи

```bash
wails3 updater genkey
# updater.key       — keep secret, use to sign releases
# updater.key.pub   — bundle in your app via go:embed
```

Закрытый ключ имеет формат PKCS#8 PEM, а открытый — PKIX PEM; `Config.PublicKey` принимает файл `.pub` без изменений (также принимается необработанный ключ длиной 32 байт или его представление в base64, которое `genkey` выводит для встраивания). Подписывайте релизы с помощью `wails3 updater manifest -key updater.key ...` или `wails3 updater sign`; см. раздел [Публикация с помощью CLI wails3](/reference/update-manifest/#publishing-with-the-wails3-cli).

Или на Go:

```go
import "crypto/ed25519"
import "crypto/rand"

pub, priv, _ := ed25519.GenerateKey(rand.Reader)
// Persist `priv` securely (HSM, signing CI, etc.); embed `pub` in your binary.
```

## Форматы артефактов

Провайдеры передают поток байтов опубликованного вами файла; затем перед заменой фреймворк распаковывает его:

- **Один исполняемый файл** (например, `myapp-linux-amd64`) — используется без изменений. Распространённый вариант для Linux.
- **`.zip`** — извлекается на место назначения. Архив должен содержать ровно один элемент верхнего уровня (обычно пакет macOS `.app` или один исполняемый файл). Рекомендуемый формат упаковки для macOS.
- **`.tar.gz`** / **`.tgz`** — извлекается на место назначения по тому же правилу об одном элементе верхнего уровня. Полезно для дистрибутивов Linux, которые поставляют дерево среды выполнения вместе с исполняемым файлом.

Архивы с несколькими элементами верхнего уровня отклоняются: фреймворк заменяет одну цель на диске, поэтому команда «установить этот архив на место назначения» неоднозначна, если архив содержит несколько объектов. `.dmg` и `.pkg` (macOS), а также `.msi` (Windows) не поддерживаются в v1 — вместо них распространяйте `.zip` пакета. При извлечении обеспечивается защита от zip-slip, отклоняются символические ссылки, выходящие за пределы корня архива, а также ограничиваются общий размер распакованных данных (2 ГиБ) и количество элементов (50 000).

## Окно по умолчанию

`app.Updater.CheckAndInstall(ctx)` открывает принадлежащее фреймворку окно размером 520×540 со следующими элементами:

- Основной значок, зависящий от состояния (синяя ↓ — обновление доступно или скачивается, зелёная ✓ — готово или установлена актуальная версия, красный ! — ошибка)
- Метка версии: `v1.0.0 → v2.0.1 · 8.8 MB`
- Прокручиваемая панель примечаний к релизу с **отрендеренным Markdown** (абзацы, полужирный текст и курсив, списки, таблицы GFM, встроенный код, блоки кода с ограждением, заголовки h1–h3, ссылки)
- Одно основное действие для каждого состояния (Установить / Перезапустить и применить / Повторить)
- Второстепенные действия в стиле ghost-кнопок (Пропустить эту версию / Напомнить позже)
- Тёмный и светлый режимы через `prefers-color-scheme`
- Мерцающий индикатор неопределённого прогресса, когда общий размер неизвестен

Окно прослушивает события `updater:*` в шине событий Wails и отправляет действия `updater:user:*` обратно в Go.

### Настройка темы с помощью переменных CSS

```go
app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{
        CSS: `:root { --accent: #ff6f00; --bg: #1a1a1a; --fg: #fafafa; }`,
    },
})
```

Таблица стилей по умолчанию предоставляет следующие переменные — можно переопределить любую из них:

| Переменная | Значение по умолчанию (светлая тема) | Значение по умолчанию (тёмная тема) |
| --- | --- | --- |
| `--bg` | `#f8f8fa` | `#1a1a1c` |
| `--surface` | `#ffffff` | `#232326` |
| `--surface-2` | `#f0f0f3` | `#2c2c30` |
| `--fg` | `#1d1d1f` | `#f5f5f7` |
| `--fg-dim` | `#6b6b73` | `#b0b0b8` |
| `--fg-faint` | `#99999f` | `#7a7a82` |
| `--border` | `#d6d6dc` | `#3a3a3e` |
| `--accent` | `#0a84ff` | `#0a84ff` |
| `--accent-fg` | `#ffffff` | — |
| `--success` | `#34c759` | — |
| `--error` | `#ff3b30` | — |
| `--radius` | `10px` | — |
| `--font` | системный стек | — |

### Замена шаблона

Предоставьте собственный HTML; он должен лишь прослушивать события `wails:updater:*` и отправлять действия `wails:updater:user:*`:

```go
app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{
        HTML: myCustomTemplate,
    },
})
```

Окна InitialHTML загружаются без источника сервера ресурсов, поэтому не могут динамически получать `/wails/runtime.js`. В таком окне есть два способа взаимодействовать с хостом:

1. **Просто напишите HTML.** Фреймворк автоматически внедряет минимальную прослойку `window.wails.Events` в любое окно, открытое с заданными `WebviewWindowOptions.AllowSimpleEventEmit = true` и `HTML`, — именно так работают встроенный сценарий средства обновления и сценарий с собственным окном. Этап сборки не требуется. Этот вариант используется в примере ниже.
2. **Соберите `@wailsio/runtime` предпочитаемым сборщиком** (Vite, esbuild, Rollup) и во время сборки импортируйте его в собственный HTML. `Events.On` работает без дополнительной настройки, поскольку выполняется исключительно на стороне клиента; `Events.Emit` использует fetch-транспорт среды выполнения, который не работает при нулевом origin. Поэтому установите небольшой postMessage-транспорт через перехватчик [`setTransport`](/guides/custom-transport/) среды выполнения, направляющий запросы через `window._wails.invoke("wails:event:emit:<name>")`. Если `window.wails.Events` уже доступен в области видимости, внедрение фреймворка ничего не делает, поэтому эти два подхода не конфликтуют.

В обоих случаях JavaScript в вашем HTML использует один и тот же API `Events.On` / `Events.Emit`:

```html
<script>
const { On, Emit } = window.wails.Events;

On("wails:updater:update-available", (e) => {
    const rel = e.data ?? e;
    document.getElementById("ver").textContent = rel.version;
});

document.getElementById("install").addEventListener("click",
    () => Emit("wails:updater:user:install"));

// Ask the host to replay the current state so we paint correctly on (re)open.
Emit("wails:updater:window:ready");
</script>
```

Адаптер предоставляет ту часть современной среды выполнения, которая нужна событиям с простыми именами: `Events.On(name, cb)` возвращает функцию отписки, а `Events.Emit(nameOrEventObject)` передаёт событие хосту через защищённый путь postMessage `wails:event:emit:`. Адаптер устанавливается один раз при загрузке страницы, до выполнения любых ваших встроенных скриптов.

Если вы *хотите* переопределить адаптер (или загружаете полную среду выполнения другим способом), задайте `window.wails.Events` до выполнения первого тега `<script>` на странице — тогда внедрение будет пропущено.

### Оформление окна

Переопределите параметры окна (размер, отсутствие рамки, отображение поверх остальных окон), не изменяя HTML:

```go
Window: &updater.BuiltinWindow{
    Options: updater.WindowOptions{
        Title:         "My App Updater",
        Width:         640,
        Height:        480,
        Frameless:     true,
        AlwaysOnTop:   true,
        DisableResize: false,
    },
},
```

### Собственное окно

Управляйте процессом обновления через созданный вами `*application.WebviewWindow`. Средство обновления вызывает для вашего окна `Show()` / `Close()` / `EmitEvent()`, а ваш HTML определяет, что отображать:

![Собственное окно средства обновления с розово-оранжевым градиентным фоном, одной белой скруглённой карточкой, пользовательской типографикой и теми же событиями средства обновления, управляющими видимым состоянием. Показывает, насколько полно можно заменить стандартный интерфейс.](/assets/updater/byo-custom-window.png)

```go
myWin := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:                "My Updater",
    Width:                520, Height: 460,
    HTML:                 myCustomHTML,
    AllowSimpleEventEmit: true,           // see security note below
})
app.Updater.Init(updater.Config{
    // …
    Window: updater.BYOWindow(myWin.AsUpdaterWindow()),
})
```

Ваш HTML использует `window.wails.Events.On` / `Events.Emit` так же, как встроенный шаблон: автоматически внедряемый фреймворком адаптер добавляется в любое окно с `AllowSimpleEventEmit: true` независимо от того, принадлежит ли окно фреймворку или вам. Описание API см. в разделе [Замена шаблона](#--4).

@note{type="caution" title="Для собственных окон средства обновления требуется `AllowSimpleEventEmit`"}
В целях безопасности фреймворк разрешает сокращённый путь postMessage `wails:event:emit:` только при включённом значении этого поля: окно, для которого оно не задано, не может создавать пользовательские события на стороне хоста. Адаптер средства обновления для собственного HTML отправляет события `updater:user:*` через этот сокращённый путь. Поэтому, если в собственном окне забыть задать это поле, все нажатия кнопок будут молча игнорироваться: пользователь нажмёт «Установить», но ничего не произойдёт.

Оставляйте `AllowSimpleEventEmit` **выключенным** для любого окна, загружающего HTML, который вы контролируете не полностью (удалённые URL-адреса, содержимое от пользователей). Если параметр включён, любой JavaScript на странице, включая уязвимые для XSS участки, может вызвать любой обработчик `app.Event.On(name, …)`. Сокращённый путь передаёт только простые имена без полезной нагрузки и не предоставляет доступ к пути binding/Call, но всё же может запускать привилегированные обработчики пользовательских событий в вашем коде Go, если они действуют, основываясь лишь на имени события.

Во *встроенном* окне средства обновления фреймворк задаёт этот параметр внутренне — помнить о нём нужно только при использовании собственного окна.

@end

### Без интерфейса

```go
app.Updater.Init(updater.Config{
    // …
    Window: updater.WindowNone,
})
```

Окно вообще не открывается. Подпишитесь на события `updater:*` в собственном интерфейсе (или в существующем главном окне) и вызывайте `app.Updater.CheckAndInstall(ctx)` из обработчика кнопки. Это удобно для периодических фоновых проверок, которые должны становиться видимыми только при обнаружении обновления, а также для приложений, интегрирующих процесс обновления в собственную панель настроек.

## События

И Go, и JavaScript подписываются через стандартную шину событий Wails. **Не вводите передаваемые строки вручную** — используйте константы, экспортируемые пакетом updater (Go) или пакетом среды выполнения (JS). Оба уровня используют один и тот же набор имён, синхронность которого контролируется регрессионным тестом.

### Из Go

Константы находятся в `github.com/wailsapp/wails/v3/pkg/updater`. Подпишитесь через `app.Event.On(name, fn)`; функция обратного вызова получает `*application.CustomEvent`, поле `Data` которого содержит типизированную полезную нагрузку, указанную в [справочнике событий](#--8). Используйте утверждение типа, а не декодирование JSON:

```go {title="main.go"}
import (
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/updater"
)

// …

app.Event.On(updater.EventUpdateAvailable, func(e *application.CustomEvent) {
    rel, ok := e.Data.(*updater.Release)
    if !ok { return }
    log.Printf("update found: %s", rel.Version)
})

app.Event.On(updater.EventDownloadProgress, func(e *application.CustomEvent) {
    p, ok := e.Data.(updater.Progress)
    if !ok { return }
    log.Printf("%d / %d bytes (%.0f KB/s)", p.Written, p.Total, p.Rate/1024)
})

app.Event.On(updater.EventError, func(e *application.CustomEvent) {
    info, ok := e.Data.(updater.ErrorInfo)
    if !ok { return }
    log.Printf("update failed during %s: %s", info.Stage, info.Message)
})
```

Все доступные константы Go:

| Константа | Передаваемая строка |
| --- | --- |
| `updater.EventCheckStarted` | `wails:updater:check-started` |
| `updater.EventUpdateAvailable` | `wails:updater:update-available` |
| `updater.EventNoUpdate` | `wails:updater:no-update` |
| `updater.EventDownloadStarted` | `wails:updater:download-started` |
| `updater.EventDownloadProgress` | `wails:updater:download-progress` |
| `updater.EventDownloadComplete` | `wails:updater:download-complete` |
| `updater.EventVerifying` | `wails:updater:verifying` |
| `updater.EventInstalling` | `wails:updater:installing` |
| `updater.EventUpdateReady` | `wails:updater:update-ready` |
| `updater.EventError` | `wails:updater:error` |
| `updater.EventMeta` | `wails:updater:meta` |
| `updater.EventWindowReady` | `wails:updater:window:ready` |
| `updater.EventUserInstall` | `wails:updater:user:install` |
| `updater.EventUserSkip` | `wails:updater:user:skip` |
| `updater.EventUserRemind` | `wails:updater:user:remind` |
| `updater.EventUserCancel` | `wails:updater:user:cancel` |
| `updater.EventUserRestart` | `wails:updater:user:restart` |

### Из JavaScript

Константы находятся в `Updater.Events` внутри `@wailsio/runtime`. Они имеют те же имена, что и в Go, и сгруппированы по подпространствам имён (`User.*`, `Window.*`), чтобы их было легко находить с помощью автодополнения:

```js
import { Events, Updater } from "@wailsio/runtime";

Events.On(Updater.Events.UpdateAvailable, (e) => {
    console.log("update found:", e.data.version);
});

Events.On(Updater.Events.DownloadProgress, (e) => {
    const p = e.data;
    console.log(`${p.written} / ${p.total} bytes (${(p.rate/1024).toFixed(1)} KB/s)`);
});

Events.On(Updater.Events.Error, (e) => {
    const info = e.data;
    console.error(`update failed during ${info.stage}: ${info.message}`);
});
```

События действий пользователя, которые ваш HTML отправляет *обратно* хосту, находятся в `Updater.Events.User`:

```js
import { Updater } from "@wailsio/runtime";

document.getElementById("install-btn").addEventListener("click", () => {
    // The framework window does this internally via the postMessage shim;
    // shown here for BYO templates that need to drive the flow themselves.
    window._wails.invoke("wails:event:emit:" + Updater.Events.User.Install);
});
```

### Справочник событий

Сторона подписки (хост → страница):

| Константа (Go) | Константа (JS) | Полезная нагрузка | Когда |
| --- | --- | --- | --- |
| `updater.EventCheckStarted` | `Updater.Events.CheckStarted` | нет | Перед каждым циклом обмена `Check` |
| `updater.EventUpdateAvailable` | `Updater.Events.UpdateAvailable` | `*Release` | `Check` обнаружил более новый выпуск |
| `updater.EventNoUpdate` | `Updater.Events.NoUpdate` | нет | `Check` подтвердил, что установлена актуальная версия |
| `updater.EventDownloadStarted` | `Updater.Events.DownloadStarted` | `*Release` | Начинается потоковая передача байтов |
| `updater.EventDownloadProgress` | `Updater.Events.DownloadProgress` | `Progress` | ~10 Гц во время загрузки |
| `updater.EventDownloadComplete` | `Updater.Events.DownloadComplete` | `*Release` | Все байты записаны, проверка ещё не началась |
| `updater.EventVerifying` | `Updater.Events.Verifying` | `*Release` | Начинается проверка подписи или дайджеста |
| `updater.EventInstalling` | `Updater.Events.Installing` | `*Release` | Начинаются распаковка и подготовка |
| `updater.EventUpdateReady` | `Updater.Events.UpdateReady` | `*Release` | Ожидается перезапуск |
| `updater.EventError` | `Updater.Events.Error` | `ErrorInfo` | Сбой на любом этапе |
| `updater.EventMeta` | `Updater.Events.Meta` | `Meta` | Один раз за сеанс перед воспроизведением снимка состояния |

На стороне страницы (страница → хост) — ваш код подписывается на эти события, если вы создаёте собственный шаблон:

| Константа (Go) | Константа (JS) | Когда |
| --- | --- | --- |
| `updater.EventWindowReady` | `Updater.Events.Window.Ready` | Окно завершило загрузку; хост воспроизводит текущее состояние |
| `updater.EventUserInstall` | `Updater.Events.User.Install` | Основное действие в состоянии `available` |
| `updater.EventUserRestart` | `Updater.Events.User.Restart` | Основное действие в состоянии `ready` |
| `updater.EventUserSkip` | `Updater.Events.User.Skip` | «Пропустить эту версию» |
| `updater.EventUserRemind` | `Updater.Events.User.Remind` | «Напомнить позже» |
| `updater.EventUserCancel` | `Updater.Events.User.Cancel` | Кнопка закрытия |

## Справочник API

### `updater.Config`

| Поле | Тип | Примечания |
| --- | --- | --- |
| `CurrentVersion` | `string` | **Обязательно.** Та же строка, которой вы помечаете релизы тегами (без префикса `v`) |
| `Providers` | `[]updater.Provider` | **Обязательно.** Упорядоченная цепочка резервных вариантов |
| `PublicKey` | `[]byte` | PEM или необработанные байты. Необязательно, но без этого подписанные релизы отклоняются |
| `CheckInterval` | `time.Duration` | Ненулевое значение запускает фоновый цикл опроса, вызывающий `CheckAndInstall` |
| `Platform` | `string` | Переопределяет `runtime.GOOS` для выбора ресурса |
| `Arch` | `string` | Переопределяет `runtime.GOARCH` для выбора ресурса |
| `Channel` | `string` | Сейчас используется только для информации; фильтрация каналов зависит от провайдера |
| `Window` | `updater.WindowOption` | `nil` (встроенные значения по умолчанию), `&BuiltinWindow{…}`, `BYOWindow(handle)` или `WindowNone` |

### Методы `*updater.Updater`

| Сигнатура | Назначение |
| --- | --- |
| `Init(cfg Config) error` | Выполняет настройку. При втором вызове возвращает `ErrAlreadyConfigured` |
| `State() State` | Текущая фаза жизненного цикла |
| `CurrentVersion() string` | Версия, переданная в `Init` |
| `Check(ctx) (*Release, error)` | Проходит по цепочке провайдеров. `(rel, nil)` = обновление найдено, `(nil, nil)` = установлена актуальная версия, `(nil, err)` = все попытки завершились сбоем |
| `DownloadAndInstall(ctx) error` | Выполняет потоковую загрузку, проверку, извлечение (если это архив) и подготовку. Требует предварительного вызова `Check` |
| `CheckAndInstall(ctx) error` | Вспомогательный метод: открывает окно, вызывает `Check`, а при обнаружении обновления — `DownloadAndInstall` |
| `Restart(ctx) error` | Запускает вспомогательный процесс, вызывает `Host.Quit` и завершает работу; вспомогательный процесс заменяет приложение и перезапускает его |
| `DownloadedPath() string` | Расположение подготовленного обновления на диске или `""`, если его нет |
| `SkipVersion(v string)` | Помечает `v` как пропущенную версию; последующие вызовы `Check` считают её актуальной |
| `SkippedVersion() string` | Возвращает текущую пропущенную версию |
| `StopPeriodicCheck()` | Отменяет таймер, запущенный `Config.CheckInterval`, и ожидает завершения цикла |

### Ошибки

| Сигнальная ошибка | Возвращается методом |
| --- | --- |
| `ErrAlreadyConfigured` | `Init` после первого успешного выполнения |
| `ErrNotConfigured` | Любая операция до `Init` |
| `ErrNoPendingRelease` | `DownloadAndInstall` без предварительного вызова `Check` |
| `ErrDownloadInProgress` | Вызов `DownloadAndInstall` во время выполнения другого экземпляра |
| `ErrNotReady` | `Restart` при отсутствии подготовленного обновления |

## Как происходит замена

`Restart` повторно запускает текущий исполняемый файл с установленными сигнальными переменными окружения. `application.New` обнаруживает их при запуске и переключается в режим вспомогательного процесса:

1. Вспомогательный процесс ждёт до 30 с, пока завершится родительский процесс с указанным PID (`platformIsAlive` опрашивает его состояние через `syscall.OpenProcess` + `GetExitCodeProcess` в Windows и `os.FindProcess` + `proc.Signal(syscall.Signal(0))` в Unix).
2. Вспомогательный процесс создаёт резервную копию целевого объекта (копирует файл либо рекурсивно копирует каталог пакета macOS `.app`).
3. Вспомогательный процесс заменяет целевой объект подготовленным артефактом, выполняя до 20 попыток с задержкой 500 мс между ними:
  - **Unix** — `os.RemoveAll(target)` + `os.Rename(newPath, target)`. Открытые файловые дескрипторы, указывающие на старый inode, остаются действительными.
  - **Windows** — `os.Rename(target, target.old.<nanos>)` + `os.Rename(newPath, target)`. Windows позволяет переименовывать файлы, образ которых всё ещё отображён в память, но не позволяет удалять их; при следующем обновлении вспомогательный процесс удаляет все оставшиеся соседние файлы `.old.*`, отображение которых в памяти ядра уже освобождено.

4. Вспомогательный процесс восстанавливает у нового исполняемого файла исходный режим доступа (загруженный файл был создан с маской umask по умолчанию, из-за чего в Unix сбрасывается `+x`; в Windows эта операция ничего не делает).
5. Вспомогательный процесс удаляет из окружения переменные режима вспомогательного процесса и повторно запускает уже заменённый исполняемый файл.
6. Вспомогательный процесс завершается.

Если запуск завершается с ошибкой, вспомогательный процесс восстанавливает резервную копию. Если родительский процесс не завершается в течение 30 с, вспомогательный процесс прерывает операцию, не изменяя целевой объект (поэтому у пользователя остаётся работающее приложение, даже если диалог завершения работы блокирует `Quit`).

Для пакетов macOS `.app`, распространяемых в виде `.zip` (рекомендуемый формат упаковки), архив распаковывается между проверкой и переходом в состояние готовности, чтобы вспомогательный процесс мог подставить на место реальный каталог.

## Периодическая проверка

```go
app.Updater.Init(updater.Config{
    // …
    CheckInterval: 6 * time.Hour,
})
```

Если `CheckInterval > 0`, фоновая горутина вызывает `CheckAndInstall` через заданные интервалы. Срабатывания таймера, поступившие во время уже выполняющегося процесса (проверки, загрузки, верификации или установки), отбрасываются: параллельные конечные автоматы не поддерживаются.

Чтобы фоновый опрос выполнялся без уведомлений, пока не будет что-либо найдено, задайте `Window: updater.WindowNone` и обрабатывайте `EventUpdateAvailable` в собственном интерфейсе.

## Пропуск и напоминание

Кнопка «Пропустить эту версию» в стандартном окне сохраняет доступную версию через `SkipVersion(rel.Version)`. Последующие вызовы `Check` обнаруживают ту же версию и считают приложение актуальным (пока пользователь не обновит `CurrentVersion`, что происходит автоматически после успешного вызова `Restart`). Кнопка «Напомнить позже» просто закрывает окно, ничего не сохраняя.

```go
// Reading what the user skipped (e.g. to surface in app settings)
if v := app.Updater.SkippedVersion(); v != "" {
    log.Printf("user skipped %s", v)
}

// Programmatically clearing the skip:
app.Updater.SkipVersion("")
```

## Контрольный список перед распространением

Перед публикацией выпуска, который будет устанавливаться средством обновления:

1. **Выберите правильный формат архива.** macOS: `.zip` пакета `.app`. Linux: один исполняемый файл или `.tar.gz`. Windows: один файл `.exe` или `.zip`. Форматы `.dmg`, `.msi` и `.pkg` не поддерживаются.
2. **Подпишите артефакт** закрытым ключом, соответствующим `Config.PublicKey`. Для лент, публикуемых провайдерами (keygen.sh, AppCast), следуйте процедуре подписания соответствующего провайдера. Для GitHub Releases с `ChecksumAsset` создайте файл `SHA256SUMS` с помощью `sha256sum` / `shasum -a 256`.
3. **Обеспечьте совпадение строк версий.** Значение `Config.CurrentVersion` должно точно совпадать с тегом версии выпуска (например, `1.0.0` ↔ тег `v1.0.0`; начальный символ `v` удаляется на стороне провайдера).
4. **Протестируйте замену на целевой платформе** хотя бы один раз перед выпуском: подписание кода, нотариальное заверение и обработка Gatekeeper зависят от платформы и не выполняются самим средством обновления.

## Устранение неполадок

**«signature requires a public key but none configured»** — у выпуска есть поле `Signature`, но `Config.PublicKey` пуст. Укажите открытый ключ или измените процесс выпуска так, чтобы подпись не добавлялась.

**«digest mismatch»** — загруженные байты не соответствуют данным, заявленным провайдером. Обычно это вызвано неполной загрузкой из-за сбоя сети или повреждением артефакта. Повторный запуск часто устраняет проблему.

**Окно открывается, но сразу исчезает; нет ни Markdown, ни индикатора хода выполнения** — пользовательский HTML не вызвал `wails:runtime:ready`. См. адаптер [«Замена шаблона»](#--4).

**Обновление в Windows не завершается; в журнале вспомогательного процесса указано «remove old (attempt N): Access is denied»** — это происходит только в версиях этого PR до `de764fb`; текущая реализация переименовывает старый файл в стороне и не подвержена этой проблеме. Обновитесь.

**Gatekeeper в macOS блокирует заменённый исполняемый файл** — подпись кода должна сохраняться на всех этапах. Подпишите исходный пакет `.app` *и* повторно подпишите перезапускаемый исполняемый файл, если процесс сборки изменяет права доступа к ресурсам приложения во время обновления.

## См. также

- Готовый к запуску пример: [`v3/examples/updater`](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/updater)
- Репозиторий с тестовой демонстрацией: [`wailsapp/updater-demo`](https://github.com/wailsapp/updater-demo)
- Руководство: [Добавление автоматического обновления в приложение Wails](/tutorials/04-self-update-a-wails-app/)
