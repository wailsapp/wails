---
title: "Самообновляемое приложение Wails"
description: "Создайте приложение Wails v3, которое самостоятельно обновляется из GitHub Releases: от `wails3 init` до проверки подписанного выпуска и замены в режиме вспомогательного процесса."
slug: "tutorials/04-self-update-a-wails-app"
sourcePath: "tutorials/04-self-update-a-wails-app.md"
---

В этом руководстве вы добавите встроенный механизм обновления в новое приложение Wails v3. В результате приложение сможет:

- Проверять наличие обновлений в GitHub Releases по запросу (и при необходимости по таймеру).
- Скачивать подходящий ресурс для используемой ОС и архитектуры.
- Сверять дайджест SHA-256 (и при необходимости подпись Ed25519) со скачанными байтами.
- Показывать примечания к выпуску в стандартном окне обновления фреймворка.
- Заменять запущенный исполняемый файл и перезапускаться — без поставки отдельного вспомогательного исполняемого файла.

В качестве источника обновлений мы будем использовать **GitHub Releases**, поскольку этот сервис бесплатен и не требует собственной инфраструктуры. Те же подходы применимы к [keygen.sh](/guides/updater/#keygensh--updaterproviderskeygen) и [Sparkle AppCast](/guides/updater/#sparkle-appcast--updaterprovidersappcast). После завершения ознакомьтесь с [руководством по механизму обновления](/guides/updater/).

@note{type="tip" title="Предварительные требования"}
- Go 1.26 или новее
- Установленный CLI `wails3` (`go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0`)
- Репозиторий GitHub, в котором вы можете публиковать выпуски
- Знакомство с [руководством по сервису QR-кодов](/tutorials/01-creating-a-service/) будет полезно, но не обязательно

@end

<br/>

@steps
### Начните с нового приложения Wails
Создайте заготовку нового проекта на основе стандартного шаблона:

```bash
wails3 init -n updater-tutorial -t vanilla
cd updater-tutorial
```

Теперь у вас должен быть каталог с `main.go`, `frontend/` и `Taskfile.yml`. Убедитесь, что проект собирается и запускается:

```bash
wails3 task dev
```

Должно открыться пустое окно Wails. Закройте приложение и продолжайте.

### Добавьте импорт механизма обновления
Откройте `main.go` и добавьте в импорты два пакета механизма обновления:

```go {title="main.go" ins="6-7"}
package main

import (
    _ "embed"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/updater"
    "github.com/wailsapp/wails/v3/pkg/updater/providers/github"
)
```

Они подключают сам Updater и провайдер GitHub Releases.

### Настройте Updater
`app.Updater` уже подключён к каждому `*application.App` — достаточно вызвать `Init`:

```go {title="main.go"}
const currentVersion = "1.0.0"

gh, err := github.New(github.Config{
    Repository:    "yourorg/your-repo",   // ← change this
    ChecksumAsset: "SHA256SUMS",          // sibling file with sha256 digests
})
if err != nil {
    log.Fatalf("github.New: %v", err)
}

if err := app.Updater.Init(updater.Config{
    CurrentVersion: currentVersion,
    Providers:      []updater.Provider{gh},
}); err != nil {
    log.Fatalf("Updater.Init: %v", err)
}
```

Поместите этот код после `application.New` и перед `app.Run()`.

@note{type="note" title="Формат строки версии"}
Передайте ту же версию, которой вы помечаете выпуски, **без** начального `v`. Со своей стороны провайдер удаляет `v` из имён тегов. Здесь `1.0.0` ↔ в GitHub `v1.0.0`.

@end

### Добавьте пункт меню для запуска обновления
В том же `main.go` добавьте пункт меню «Проверить наличие обновлений…»:

```go {title="main.go"}
menu := app.Menu.New()
app.Menu.SetApplicationMenu(menu)
appMenu := menu.AddSubmenu("App")
appMenu.Add("Check for Updates…").OnClick(func(*application.Context) {
    go func() {
        if err := app.Updater.CheckAndInstall(context.Background()); err != nil {
            app.Logger.Error("update", "error", err)
        }
    }()
})
```

`CheckAndInstall` открывает окно обновления фреймворка, выполняет `Check` и, если выпуск найден, автоматически выполняет `DownloadAndInstall`. Если новых выпусков нет, окно остаётся открытым в состоянии «Обновлений нет» — пользователь закрывает его кнопкой **Закрыть**.

@note{type="caution" title="Запускайте в горутине"}
`CheckAndInstall` блокирует выполнение до завершения проверки и установки. Прямой вызов из обработчика нажатия пункта меню заблокировал бы поток пользовательского интерфейса. Оберните вызов в `go func()`.

@end

### Запустите один раз при отсутствии выпусков
```bash
wails3 task dev
```

Нажмите **Приложение → Проверить наличие обновлений…**. Окно обновления должно ненадолго открыться, обратиться к API GitHub, не найти выпусков новее `1.0.0` и перейти в состояние **Обновлений нет** с зелёной галочкой ✓.

Если на этом этапе возникает ошибка, обычно причина одна из следующих:

| Симптом | Решение |
| --- | --- |
| `404 Not Found` | Поле `Repository` указано неверно — должно быть `owner/repo` |
| `403 rate-limited` | Добавьте `Token: "ghp_…"` в github.Config (используйте PAT с областью действия `public_repo`) |
| Сетевые ошибки | Убедитесь, что запущенное приложение может подключиться к `api.github.com` |

### Опубликуйте тестовый выпуск
Измените `currentVersion` в `main.go` на `1.0.0` (или оставьте как есть). Соберите приложение для одной платформы, чтобы получить исполняемый файл, который можно прикрепить к выпуску:

@tabs
[macOS]
```bash
wails3 task build:darwin
# produces bin/updater-tutorial.app
# zip it for the release asset:
cd bin && zip -r updater-tutorial-darwin-arm64.zip updater-tutorial.app && cd ..
```

[Linux]
```bash
wails3 task build:linux
# produces bin/updater-tutorial
mv bin/updater-tutorial bin/updater-tutorial-linux-amd64
```

[Windows]
```bash
wails3 task build:windows
# produces bin/updater-tutorial.exe
mv bin/updater-tutorial.exe bin/updater-tutorial-windows-amd64.exe
```

@end

Создайте файл `SHA256SUMS` рядом с исполняемым файлом:

```bash
cd bin
shasum -a 256 updater-tutorial-* > SHA256SUMS
cat SHA256SUMS
```

Вы должны увидеть одну или несколько строк следующего вида:

```
abc123…  updater-tutorial-darwin-arm64.zip
```

Теперь опубликуйте это в своём репозитории GitHub как **v2.0.0**:

```bash
gh release create v2.0.0 \
    --title "v2.0.0" \
    --notes "First update for the self-update tutorial.

- **Bold** Markdown renders in the update window
- \`Code spans\` too
- Lists work
- GFM tables work" \
    bin/SHA256SUMS bin/updater-tutorial-*
```

@note{type="note" title="Именование ресурсов"}
Стандартный механизм сопоставления ресурсов ищет в имени файла подстроки `GOOS` и `GOARCH`. Он найдёт ресурс, если его имя содержит `darwin` (либо `linux` или `windows`) и `arm64` (либо `amd64` или `386`). Настраиваемые механизмы сопоставления описаны в [руководстве по механизму обновления](/guides/updater/#github-releases--updaterprovidersgithub).

@end

### Запустите приложение и проверьте обновление
Пока `currentVersion` по-прежнему имеет значение `1.0.0`, снова запустите приложение:

```bash
wails3 task dev
```

Нажмите **Приложение → Проверить наличие обновлений…**. На этот раз вы должны увидеть примерно следующее:

![Стандартное окно обновления в состоянии «Обновление готово» со значком версии, примечаниями к выпуску, отформатированными из Markdown, и основной кнопкой «Перезапустить и применить».](/assets/updater/default-window-ready.png)

- Главный значок меняется с синей стрелки ↓ («Доступно обновление») на зелёную галочку ✓ («Обновление готово»).
- В подзаголовке отображается `v1.0.0 → v2.0.0 · <size>`.
- На панели примечаний к выпуску ваш Markdown отображается с полужирным начертанием, фрагментами кода и таблицей.
- Во время скачивания индикатор выполнения заполняется (это произойдёт быстро — исполняемый файл небольшой).

Updater помещает новый исполняемый файл во временный каталог. Чтобы завершить обновление:

- Нажмите **Перезапустить и применить**.
- Приложение завершает работу, вспомогательный процесс заменяет исполняемый файл, после чего запускается новый исполняемый файл.
- Перезапущенное приложение сообщает версию `currentVersion = "1.0.0"` (поскольку мы жёстко задали её в коде), но содержимое файла на диске соответствует сборке v2.0.0.

В реальном приложении значение `currentVersion` задавалось бы во время сборки через `-ldflags`, чтобы новый исполняемый файл знал, что теперь его версия — v2.0.0, и при следующей проверке обновление не обнаруживалось.

### Свяжите `currentVersion` со сборкой
Замените константу переменной, задаваемой во время сборки:

```go {title="main.go" ins="2,4"}
var (
    currentVersion = "dev" // overridden by -ldflags at release time
)
```

Затем укажите в команде сборки:

```bash
wails3 task build:darwin -- -ldflags "-X main.currentVersion=2.0.0"
```

Либо добавьте `-ldflags` в `Taskfile.yml`, чтобы значение бралось из `git describe --tags`.

### Добавьте криптографическую подпись (рекомендуется для рабочей среды)
Проверка по SHA256SUMS подтверждает *целостность* (байты совпадают с сохранёнными в GitHub), но не *подлинность* (что эти байты созданы вашим конвейером выпуска, а не получены через скомпрометированную учётную запись сопровождающего). Для защиты от подмены подписывайте каждый выпуск ключом Ed25519:

```bash
# One-time: generate the keypair
ssh-keygen -t ed25519 -f updater-key -N "" -C "wails-updater"
#   updater-key      — keep secret (build server, HSM, password manager)
#   updater-key.pub  — bundle in your app
```

Для каждого выпуска подпишите закрытым ключом дайджест SHA-256 каждого артефакта. Небольшая вспомогательная программа на Go:

```go {title="cmd/sign-release/main.go"}
package main

import (
    "crypto/ed25519"
    "crypto/sha256"
    "encoding/base64"
    "fmt"
    "io"
    "os"
)

func main() {
    priv, _ := os.ReadFile("updater-key")
    key := ed25519.PrivateKey(priv) // raw 64-byte private key

    f, _ := os.Open(os.Args[1])
    defer f.Close()
    h := sha256.New()
    _, _ = io.Copy(h, f)
    sig := ed25519.Sign(key, h.Sum(nil))
    fmt.Println(base64.StdEncoding.EncodeToString(sig))
}
```

Стандартный провайдер GitHub сейчас не загружает отдельный файл подписи. Вы можете [написать собственный провайдер](/guides/updater/#writing-your-own-provider), который будет это делать, или перейти на **keygen.sh**: он подписывает каждый артефакт на стороне сервера и предоставляет через свой API как дайджест, так и подпись.

Встройте открытый ключ в приложение:

```go {title="main.go" ins="1,6"}
//go:embed updater-key.pub
var updaterPublicKey []byte

app.Updater.Init(updater.Config{
    CurrentVersion: currentVersion,
    PublicKey:      updaterPublicKey,
    Providers:      []updater.Provider{gh},
})
```

Если задано `PublicKey`, любой выпуск, содержащий `Signature`, должен пройти проверку этим ключом. Источник выпуска не может подменить его собственным ключом — именно для этого ключ закрепляется по независимому каналу во время сборки.

### Настройте окно
Стандартное окно подходит для большинства случаев. Если вам нужен более полный контроль, есть три варианта — выберите подходящий в зависимости от требуемой степени настройки:

@tabs
[Только CSS]
```go
app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{
        CSS: `:root { --accent: #ff6f00; --radius: 16px; }`,
    },
})
```

Полный список переменных приведён в разделе [Настройка темы с помощью переменных CSS](/guides/updater/#theme-via-css-variables).

[Собственный HTML]
```go
//go:embed updater-window.html
var updaterHTML string

app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{HTML: updaterHTML},
})
```

Ваш HTML должен подписываться на события `updater:*` и отправлять действия `updater:user:*` через канал событий Wails. JS-адаптер приведён в разделе [Замена шаблона](/guides/updater/#replace-the-template).

[Собственное окно]
```go
myWin := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:                "My App Updater",
    Width:                520, Height: 460,
    HTML:                 updaterHTML,
    AllowSimpleEventEmit: true,  // required — see security note
})
app.Updater.Init(updater.Config{
    // …
    Window: updater.BYOWindow(myWin.AsUpdaterWindow()),
})
```

Этот вариант полезен, если у вас уже есть собственная оконная инфраструктура и вы хотите, чтобы средство обновления использовало её вместо открытия ещё одного окна. Полностью собственный HTML-шаблон, управляемый теми же событиями средства обновления, что и стандартный шаблон, выглядит так:

![Собственное окно средства обновления с розово-оранжевым градиентным фоном и нестандартной компоновкой карточки со скруглёнными углами, демонстрирующее, что стандартный интерфейс можно полностью заменить.](/assets/updater/byo-custom-window.png)

@note{type="caution" title="Требуется `AllowSimpleEventEmit`"}
Адаптер собственного HTML средства обновления запускает действия Install / Skip / Remind / Restart через сокращённый вызов `wails:event:emit:` postMessage, а в целях безопасности доступ к этому вызову регулируется данным полем. Если его не задать, нажатия кнопок не будут давать никакого результата. Не включайте его для окон, загружающих HTML, который вы не контролируете полностью. Модель угроз описана в разделе руководства [Собственное окно](/guides/updater/#bring-your-own-window).

@end

@end

### Запускайте автоматические проверки в фоновом режиме
Чтобы выполнять проверку по таймеру вместо нажатия пункта меню или в дополнение к нему:

```go {ins="5"}
app.Updater.Init(updater.Config{
    CurrentVersion: currentVersion,
    Providers:      []updater.Provider{gh},
    PublicKey:      updaterPublicKey,
    CheckInterval:  6 * time.Hour,
})
```

При каждом срабатывании таймера выполняется тот же процесс `CheckAndInstall`, что и при нажатии вручную. Задайте `Window: updater.WindowNone`, если периодическая проверка должна выполняться незаметно, пока обновление действительно не будет найдено. Затем самостоятельно подпишитесь на `EventUpdateAvailable` и определите, какой интерфейс следует показать.

@end

## Готово

Теперь ваше приложение Wails:

- Проверяет наличие обновлений в GitHub Releases по запросу и по таймеру.
- Отображает примечания к выпуску в формате Markdown в качественно оформленном стандартном окне.
- Проверяет загружаемые файлы по опубликованному вами дайджесту SHA-256.
- При необходимости проверяет подпись Ed25519 с помощью открытого ключа, встроенного во время сборки.
- Заменяет запущенный исполняемый файл на месте и автоматически перезапускает приложение.

## Дальнейшие действия

- В [руководстве по средству обновления](/guides/updater/) приведены полный справочник API, все события, все параметры конфигурации и описание механизма замены в режиме вспомогательного процесса.
- Полный рабочий пример, который можно клонировать, находится в [`v3/examples/updater`](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples/updater).
- Рекомендуемая структура артефактов выпуска показана в целевом тестовом репозитории [`wailsapp/updater-demo`](https://github.com/wailsapp/updater-demo).

## На что обратить внимание при использовании в рабочей среде

- **Подпись кода в macOS** — Gatekeeper требует, чтобы заменяемый исполняемый файл был подписан и нотариально заверен. Подпишите пакет `.app` *до* его упаковки в ZIP-архив для выпуска. Средство обновления сохраняет байты без изменений и ничего не подписывает заново.
- **Антивирус в Windows** — неподписанные файлы `.exe`, загруженные из интернета, могут вызывать предупреждения SmartScreen. Подпишите исполняемый файл сертификатом Authenticode либо учитывайте, что пользователям компьютеров с жёсткими ограничениями может потребоваться добавить ваше приложение в список разрешённых.
- **Атомарные выпуски** — публикуйте `SHA256SUMS` вместе с исполняемыми файлами, а не отдельными коммитами. Средство обновления загружает сопутствующий файл отдельно от исполняемого; если они рассинхронизированы, проверка дайджеста завершается отказом.
- **Пропущенные версии** — кнопка «Пропустить эту версию» в стандартном окне сохраняет выбор локально. Если вы выпускаете критическое обновление безопасности, присвойте ему новый номер версии, чтобы оно не было автоматически пропущено у пользователей, отклонивших предыдущий выпуск.
