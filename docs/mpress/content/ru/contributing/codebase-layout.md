---
title: "Структура кодовой базы"
description: "Как организован репозиторий Wails v3 и как взаимодействуют его компоненты"
slug: "contributing/codebase-layout"
sourcePath: "contributing/codebase-layout.md"
---

Wails v3 размещается в **монорепозитории**, который содержит среду выполнения фреймворка, CLI, примеры, документацию и набор инструментов сборки. На этой странице рассматривается значимая для изучения внутреннего устройства *структура каталогов*.

## Обзор верхнего уровня

```
wails/
├── v3/               # ⬅️ Everything specific to Wails v3 lives here
├── v2/               # Legacy v2 implementation (can be ignored for v3 work)
├── docs/             # M-Press-powered v3 docs site (this page!)
├── website/          # Docusaurus v2 site and marketing pages (main site)
├── scripts/          # Misc helper scripts (e.g. sponsor image generator)
└── *.md              # Project-wide meta files (CHANGELOG, LICENSE, …)
```

Далее мы подробно рассмотрим дерево **`v3/`**.

## Корень `v3/`

```
v3/
├── cmd/          # Compilable commands (currently only the wails3 CLI)
├── internal/     # Framework implementation (not public API)
├── pkg/          # Public Go packages — the API surface
├── tasks/        # Taskfile-based release / generation utilities
├── wep/          # RFC-style proposals (Wails Enhancement Proposals)
├── tests/        # Integration test harness
├── go.mod
└── go.sum
```

> Шаблоны проектов поставляются в каталоге `internal/templates/` (по одной папке на каждый фреймворк
>
> -стек, а также `base/`, `_common/` и `ios/`). Каталога `v3/templates/`
>
> на верхнем уровне нет.

### Ментальная модель

1. **`pkg/`** предоставляет *то, что импортируют разработчики приложений*\
2. **`internal/`** содержит *реализацию внутренней логики*\
3. **`cmd/wails3`** управляет *жизненным циклом проекта и сборками*\

Всё остальное служит опорой для этих трёх ключевых компонентов.

---

## `cmd/` — Команды

| Путь | Примечания |
| --- | --- |
| `v3/cmd/wails3` | Точка входа **CLI**. Минимальный `main.go` передаёт всю логику пакетам в `internal/commands`. |
| `internal/commands/*` | Подкоманды (init, dev, build, doctor, …). Каждая находится в отдельном файле, поэтому её легко найти. |
| `internal/commands/task_wrapper.go` | Связывает флаги CLI с конвейером сборки Taskfile. |

CLI отвечает за:

- **Создание каркаса проекта** (`init`, создание шаблонов)\
- **Оркестрацию сервера разработки** (`dev`, автоматическая перезагрузка)\
- **Сборку и упаковку для рабочей среды** (`build`, `package`, платформенные обёртки)\
- **Диагностику** (`doctor`)\

---

## `internal/` — Внутреннее устройство

```
internal/
├── assetserver/  # Serving & embedding web assets
├── buildinfo/    # Reproducible build metadata
├── commands/     # CLI mechanics (see above)
├── runtime/      # Build-tag glue + embedded JS runtime sources
├── generator/    # Static analysis & binding generator
├── templates/    # Project templates (frontend stacks)
├── packager/     # nfpm wrapper used by `wails3 tool package`
├── capabilities/ # Host OS capability probing
├── dbus/         # Generic D-Bus helper
├── service/      # Service-template scaffolding (`wails3 generate service`)
└── ...           # [other helper sub-packages: flags, hash, term, …]
```

### Ключевые подпакеты

| Пакет | Назначение | Точки взаимодействия |
| --- | --- | --- |
| `runtime` | Содержит небольшой связующий код с тегами сборки `runtime{,_darwin,_linux,_windows,_android,_dev,_prod}.go`, а также встроенную среду выполнения JS в `runtime/desktop/`. Реализации окон, буфера обмена, диалоговых окон и области уведомлений для каждой ОС находятся в `pkg/application/*_{darwin,linux,windows}.go`. | Импортируется косвенно через `pkg/application`. |
| `assetserver` | Файловый сервер с двумя режимами работы:<br />• Разработка: обслуживает файлы с диска и проксирует Vite (`build_dev.go`)<br />• Продакшен: встраивает ресурсы через `go:embed` (`build_production.go`) | Инициализируется `pkg/application` при запуске. |
| `generator` | Анализирует исходный код Go для формирования **метаданных привязок**, на основе которых затем создаются файлы-заглушки TypeScript/JS и константы событий. Точки входа: `generator.Generate` / `generator.Generator` поверх `collect/` + `render/`. | Запускается из `wails3 generate bindings`. |
| `packager` | Обёртка `nfpm`, используемая для создания артефактов Linux в форматах `deb`/`rpm`/`archlinux` (на основе конфигураций nfpm `myapp.DEB`/`.RPM`/`.ARCHLINUX` из `internal/commands/`). | Вызывается из `wails3 tool package`. Средства создания DMG для macOS и MSIX для Windows находятся в `internal/commands/{dmg,msix.go,webview2/}`. |

Вспомогательные утилиты (например, `s/`, `hash/` и `flags/`) обеспечивают независимость внутренних компонентов друг от друга.

---

## `pkg/` — публичный API

```
pkg/
├── application/  # Core API: App, windows, menus, dialogs, events, managers
├── events/       # Event constants (Common/Mac/Windows/Linux) + generator
├── services/     # Optional built-in services (notifications, kvstore, …)
├── doctor-ng/    # New-style `wails3 doctor-ng` checks
├── errs/         # Shared error types
├── icons/        # Default platform icons
├── mac/          # macOS-only helpers
└── w32/          # Windows Win32 helpers
```

> Пакетов `pkg/runtime/`, `pkg/options/` и `pkg/menu/` не существует. Параметры окон и меню
>
> располагаются рядом с `pkg/application` (например, `WebviewWindowOptions`, `Menu`,
>
> `MenuItem`), а `assetserver/` находится в `internal/`.

`pkg/application` инициализирует программу Wails:

```go
func main() {
    app := application.New(application.Options{
        Name: "MyApp",
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assetsFS),
        },
    })
    window := app.Window.New()
    window.SetTitle("Hello").SetSize(1024, 768)
    _ = app.Run()
}
```

Внутри пакет выполняет следующие действия:

1. Подключает связующий код для тегов сборки из `internal/runtime` и код для каждой ОС из `pkg/application/`
2. Настраивает экземпляр `internal/assetserver`
3. Регистрирует все обработчики сообщений, используемые привязками
4. Переходит в главный поток ОС

---

## `internal/templates/` — шаблоны каркаса проекта

`internal/templates/` включает **базовые шаблоны** (структура Go в `base/`, `_common/`, `ios/`) и **шаблоны оформления фронтенда** (`vanilla[-ts]`, `react[-ts]`, `react-swc[-ts]`, `lit[-ts]`, `preact[-ts]`, `qwik[-ts]`, `solid[-ts]`, `svelte[-ts]`, `sveltekit[-ts]`, `vue[-ts]`).

При выполнении `wails3 init -t react` интерфейс командной строки:

1. Копирует файлы Go из `_common`
2. Объединяет выбранный пакет фронтенда
3. Запускает `go mod tidy` (можно пропустить с помощью `--skipgomodtidy`)

Изменение шаблонов **не** влияет на существующие приложения — только на будущие запуски `init`. Публичные примеры находятся в `v3/examples/`; они не заменяют автоматизированные наборы тестов, описанные в документации для участников проекта.

---

## `tasks/` — автоматизация выпусков

Taskfile-файлы служат обёртками для сложных операций кросс-компиляции, изменения версии и создания журнала изменений. `internal/commands/task.go` использует их программно, поэтому одна и та же логика применяется в **CLI** и **CI**.

---

## Как взаимодействуют компоненты

```d2
direction: down
CLI: CLI wails3
Generator: internal/generator
AssetDev: assetserver (разработка)
Packager: internal/packager
AppRuntime: {
  label: Среда выполнения приложения
  ApplicationPkg: pkg.application
  InternalRuntime: internal.runtime
  OSAPIs: API ОС
}
CLI -> Generator: сборка / генерация
CLI -> AssetDev: разработка
CLI -> Packager: упаковка
Generator -> ApplicationPkg: привязки
ApplicationPkg -> InternalRuntime
InternalRuntime -> OSAPIs
ApplicationPkg -> AssetDev
ApplicationPkg.label: ApplicationPkg
InternalRuntime.label: InternalRuntime
OSAPIs.label: OSAPIs
```

*CLI → генератор → среда выполнения* образует основной путь от **исходного кода** к **работающему настольному приложению**.

---

## Советы по навигации

| Что нужно понять… | Где искать… |
| --- | --- |
| Платформенные адаптеры | `pkg/application/*_darwin.go`, `*_linux.go`, `*_windows.go` (окна, буфер обмена, диалоговые окна, системный трей, основной поток, events_common). cgo в Linux: `pkg/application/linux_cgo*.go`. |
| Протокол моста | `pkg/application/messageprocessor*.go` |
| Рабочий процесс с ресурсами | `internal/assetserver/` (`build_dev.go` и `build_production.go`) |
| Процесс упаковки | `internal/commands/{appimage,msix,dot_desktop,dmg/}.go`, `internal/packager/` |
| Механизм шаблонов | `internal/templates/` (`templates.Install`, `templates.GetDefaultTemplates`) |
| Статический анализ | `internal/generator/{generate.go,collect/,render/}` |

---

Теперь у вас есть **мысленная карта** репозитория. Используйте её вместе с `ripgrep`, командами «Перейти к файлу/символу» в вашей IDE и примерами приложений, чтобы глубже изучить любую функциональность. Удачного хакинга!
