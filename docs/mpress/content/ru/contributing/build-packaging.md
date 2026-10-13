---
title: "Конвейер сборки и упаковки"
description: "Что происходит внутри при запуске `wails3 build`, как создаются кроссплатформенные исполняемые файлы и как формируются установочные пакеты для каждой ОС."
slug: "contributing/build-packaging"
sourcePath: "contributing/build-packaging.md"
---

`wails3 build` намеренно сделана **тонкой**: это обёртка над Taskfile, которая передаёт дополнительные теги сборки задаче `build` основного проекта. Основная работа выполняется в собственном файле `build/Taskfile.yml` проекта (созданном командой `wails3 init`), в `internal/commands/build-assets.go` (управляющем ресурсами, внедряемыми при сборке), а также в `internal/packager` (упаковка для Linux с помощью nfpm) и `internal/commands/appimage.go`, `internal/commands/msix.go`, `internal/commands/dmg/dmg.go`, и `internal/commands/dot_desktop.go` (установщики для отдельных платформ).

На этой странице рассматриваются:

1. Фактические точки входа CLI
2. Процесс сборки на основе Taskfile
3. Внедрение ресурсов и сведений о сборке
4. Механизмы упаковки для каждой платформы
5. Настройка конвейера
6. Устранение неполадок

---

## 1. Фактические точки входа CLI

```
wails3 build       → internal/commands.Build       (in task_wrapper.go)
wails3 package     → internal/commands.Package     (in task_wrapper.go)
wails3 generate build-assets → GenerateBuildAssets (in build-assets.go)
wails3 update build-assets   → UpdateBuildAssets   (in build-assets.go)
wails3 tool buildinfo        → BuildInfoOptions    (in tool_buildinfo.go)
wails3 tool package          → internal/packager   (nfpm wrapper)
wails3 generate .desktop     → in dot_desktop.go
```

`internal/commands/task_wrapper.go`:

```go
func Build(buildFlags *flags.Build, otherArgs []string) error {
    // forwards --tags / EXTRA_TAGS, then defers to a Taskfile target
    return wrapTask("build", otherArgs)
}
```

`flags.Build` предоставляет **единственный** флаг — `--tags` (передаваемый как `EXTRA_TAGS=`). У `wails3 build` **нет** флагов `-platform`, `-o`, `-skipbindings`, `-skip-package`, `-package`, `-ldflags`, `-verbose`, `-debug`, `-devbuild`, `-icon` или `-clean`. Кросс-компиляция, пути вывода, значки и прочее настраиваются в **`Taskfile.yml`**, **`build/config.yml`** и вспомогательных командах `wails3 generate icons` / `wails3 generate build-assets`.

`build/build.json` **не** входит в v3 — конфигурация задаётся в `Taskfile.yml` и `build/config.yml`.

---

## 2. Процесс сборки на основе Taskfile

В только что инициализированный проект входит `build/Taskfile.yml` примерно со следующими пространствами имён:

| Пространство имён | Задачи (выборочно) |
| --- | --- |
| `darwin:` | `build`, `build:universal`, `package`, `run`, `dev` |
| `windows:` | `build`, `package`, `run`, `dev` |
| `linux:` | `build`, `package`, `run`, `dev` |
| `common:` | `update:build-assets`, `generate:icons`, `generate:syso` |

По умолчанию `wails3 build` вызывает пространство имён `build` основной ОС; затем Taskfile проекта запускает в оболочке `go build` с флагами для этой ОС. Чтобы выполнить сборку для другой ОС, запустите её задачу напрямую (например, `wails3 task darwin:build:universal`), а не передавайте флаг команде `wails3 build`.

Каталог вывода по умолчанию — **`bin/<APP_NAME>`** (без префикса `build/bin/`).

---

## 3. Ресурсы и сведения о сборке, внедряемые во время сборки

| Назначение | Файл |
| --- | --- |
| Создание и обновление ресурсов сборки | `internal/commands/build-assets.go` |
| Вывод сведений о сборке (CLI: `wails3 tool buildinfo`) | `internal/commands/tool_buildinfo.go` — выводит сведения; это **не** средство внедрения `ldflags` |
| Заглушка для рабочей сборки | `internal/assetserver/build_production.go` — `//go:build production` |
| Пакеты файлов фронтенда | внедряются через `//go:embed` в собственный пакет приложения (например, рядом с `main.go`) |
| Ресурсы Windows (`.syso`) | `internal/commands/syso.go` — создаёт `rsrc_windows_<arch>.syso` |
| MSIX для Windows | `internal/commands/msix.go` + `internal/commands/webview2/` |
| Входные файлы DMG для macOS | `internal/commands/dmg/` |
| `.desktop` для Linux | `internal/commands/dot_desktop.go` |

CLI не генерирует `bundled_assetserver.go` для вашего приложения автоматически — `internal/assetserver/bundled_assetserver.go` **пишется вручную** и служит обёрткой для встроенной среды выполнения JS, расположенной в `bundledassets/`.

---

## 4. Механизмы упаковки

### Linux

Упаковка для Linux выполняется с помощью **nfpm** (а не `fpm`):

- `internal/packager/packager.go` служит обёрткой для `github.com/goreleaser/nfpm/v2` и предоставляет `CreatePackageFromConfig(pkgType, configPath, output)` / `CreatePackageFromConfigWriter(...)`.
- Созданные проекты содержат конфигурации `myapp.DEB`, `myapp.RPM` и `myapp.ARCHLINUX` в формате nfpm в каталоге `internal/commands/` (их использует `wails3 tool package`).
- Создание AppImage реализовано в `internal/commands/appimage.go`, который вызывает `linuxdeploy` + `linuxdeploy-plugin-gtk` (плагин поставляется по адресу `internal/commands/linuxdeploy-plugin-gtk.sh`).

У `wails3 build` **нет** флага `-package deb`/`rpm`. Используйте `wails3 tool package` или платформозависимую цель Taskfile.

### macOS

- `darwin:package` из Taskfile проекта создаёт пакет `.app`.
- Ресурсы DMG находятся в `internal/commands/dmg/`; после завершения `darwin:package` проект может упаковать пакет приложения в DMG с помощью `hdiutil` (Taskfile в новых шаблонах содержит вспомогательную задачу `dmg`).
- Идентификаторы CFBundle, версия и сведения об авторских правах задаются флагами `-product*` во время `wails3 init`, а также берутся из `build/config.yml`.

### Windows

- Для упаковки в Windows используется **MSIX**, а не WiX/MSI. Полный рабочий процесс описан в `internal/commands/msix.go` и `internal/commands/webview2/`.
- **Нет** `internal/commands/packager.go` и **нет** каталога `internal/commands/windows_resources/`.
- Необязательное подписывание исполняемого файла выполняется через `wails3 tool sign` (Authenticode) — см. `internal/commands/sign.go`.

---

## 5. Настройка конвейера

| Задача | Решение |
| --- | --- |
| Дополнительные теги сборки | `wails3 build --tags myFeature,otherTag` |
| Линтер или этап перед сборкой | Добавьте задачу в `build/Taskfile.yml` и укажите её как зависимость специфичной для ОС задачи `build` |
| Кросс-компиляция | Запустите соответствующую задачу для нужной ОС (например, `wails3 task linux:build`) — флага `-platform` нет |
| Пропуск упаковки | Запустите только задачу `build`; `package` выполняется отдельно |
| Собственный инструмент упаковки | Поместите конфигурацию в `internal/commands/myapp.*` и вызовите `wails3 tool package` с параметром `-config <file>` |
| Удаление символов | Измените задачу `darwin:/windows:/linux:` `build`, чтобы передать `-ldflags "-s -w"` непосредственно в `go build`: у самого `wails3 build` нет флага `-ldflags` |

Все цели Taskfile учитывают переменные окружения, публикуемые Wails (`APP_NAME`, `WAILS_VITE_PORT`, `FRONTEND_DEVSERVER_URL`, …), поэтому на них можно полагаться в собственных задачах.

---

## 6. Устранение неполадок

| Симптом | Вероятная причина | Решение |
| --- | --- | --- |
| **`ld: framework not found WebKit` (mac)** | Отсутствуют инструменты командной строки Xcode | `xcode-select --install` |
| **Пустое окно в сборке для промышленной эксплуатации** | Сбой сборки фронтенда или маршрутизация SPA | Убедитесь, что `frontend/dist/index.html` существует и обработчик ресурсов использует его как резервный вариант |
| **Отсутствуют инструменты упаковки MSIX** | Не установлены `WebView2` SDK или инструменты MSIX | Запустите `wails3 task install:msix:tools` |
| **`linuxdeploy` не найден** | Плагин отсутствует в PATH | Установите `linuxdeploy` и запустите `internal/commands/linuxdeploy-plugin-gtk.sh` с помощью предусмотренного в CLI этапа автоматической установки |

У `wails3 build` нет флага `-verbose`. Задайте `TASK_X_VERBOSE=1` (Taskfile) или просмотрите непосредственно целевую задачу, чтобы увидеть выполняемые команды.

---

## 7. Карта основных исходных файлов

| Назначение | Файл |
| --- | --- |
| Обёртка сборки | `internal/commands/task_wrapper.go` (`Build`, `Package`, `SignWrapper`, `wrapTask`) |
| Создание ресурсов сборки | `internal/commands/build-assets.go` (`GenerateBuildAssets`, `UpdateBuildAssets`) |
| Вывод сведений о сборке | `internal/commands/tool_buildinfo.go` |
| Сборщик AppImage | `internal/commands/appimage.go` |
| Упаковка для Linux (nfpm) | `internal/packager/packager.go`, `internal/commands/myapp.{DEB,RPM,ARCHLINUX}` |
| MSIX для Windows | `internal/commands/msix.go`, `internal/commands/webview2/` |
| Генератор ресурсов Windows | `internal/commands/syso.go` |
| Ресурсы DMG для macOS | `internal/commands/dmg/` |
| Генератор `.desktop` | `internal/commands/dot_desktop.go` |
| Константы версии | `internal/version/version.go` |

Держите эту таблицу под рукой при поиске причин сбоя сборки.

---

Теперь у вас есть полная картина процесса: от **исходного кода** до **установщика**. Вкратце: сам `wails3 build` — лишь тонкая обёртка; почти все настройки выполняются в проектных `Taskfile.yml` и `build/config.yml` либо через явные подкоманды `wails3 generate …` и `wails3 tool …`. Успешных выпусков!
