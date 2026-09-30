---
title: "Система шаблонов"
description: "Как Wails v3 создаёт каркас новых проектов, как организованы шаблоны и как создавать собственные."
slug: "contributing/template-system"
sourcePath: "contributing/template-system.md"
---

Wails поставляется с **системой шаблонов**, которая позволяет `wails3 init` создать готовый к запуску проект. Встроенные шаблоны предусмотрены для намеренно небольшого набора фреймворков (Vanilla, React, Vue, Svelte); любой другой фреймворк можно использовать, [подключив собственный фронтенд](/guides/dev/frontend-frameworks/) или опубликовав [пользовательский шаблон](/guides/advanced/custom-templates/).

На этой странице рассматриваются:

1. Структура каталогов шаблонов
2. Как CLI выбирает и обрабатывает шаблоны
3. Пошаговое создание нового шаблона
4. Обновление или переопределение существующих шаблонов
5. Устранение неполадок и рекомендации

---

## 1. Где находятся шаблоны

```
v3/internal/templates/
├── _common/        # Files copied into EVERY project (Taskfile.yml, build/, etc.)
├── base/           # Backend-only "plain Go" base layer (frontend/ + NEXTSTEPS.md)
├── ios/            # iOS bootstrapper
├── vanilla/        vanilla-js/   # TypeScript (default) + JavaScript variant
├── react/          react-js/     # TypeScript (default) + JavaScript variant
├── vue/                          # TypeScript only
├── svelte/                       # TypeScript only
└── templates.go    # Registry + Install/Get APIs (no auto-registration via embed)
```

- **`_common/`** — универсальная основа (Taskfile, каталог `build/` и общая инфраструктура), добавляемая в каждый проект.
- **`base/`** — часть на Go, с которой начинается каждый шаблон. Обратите внимание: сам `base/` **не** содержит `template.json`; этот файл находится внутри каждого шаблона для конкретного фреймворка.
- **Каталоги фреймворков** — содержат фронтенд (`frontend/`), конфигурацию фреймворка и `template.json` с метаданными шаблона.
- Имена каталогов совпадают с **идентификатором шаблона**, который передаётся CLI (`wails3 init -t react`).
- **Соглашение об обозначении языка:** TypeScript используется по умолчанию и получает имя без суффикса (`react`); вариант на JavaScript, если он существует, получает суффикс `-js` (`react-js`). Встроенные шаблоны явно указывают язык с помощью `typescript: true|false` в `template.yaml`. Шаблоны сообщества всё ещё могут использовать устаревший суффикс `-ts`, который поддерживается как резервный вариант.

> Весь каталог `internal/templates/` встраивается в исполняемый файл CLI
>
> с помощью `//go:embed *`, поэтому пользователи могут создавать каркас проектов без подключения к сети.

---

## 2. Как `wails3 init` использует шаблоны

Цепочка вызовов (без `cmd/wails3/init.go` — CLI подключается непосредственно в `cmd/wails3/main.go`):

```
cmd/wails3/main.go             (clir wiring)
       │
       ▼
internal/commands/init.go      Init(options *flags.Init) error
       │
       ▼
internal/templates/templates.go
       │   templates.Install(options)
       │   templates.GetDefaultTemplates()
       ▼
gosod.New(template.FS).Extract(options.ProjectDir, data)   // file extraction
       │
       ▼
go mod tidy (unless --skipgomodtidy / -skipgomodtidy)
```

API `Template.Load()` / `Template.CopyTo()` / `Template.Validate()` не существует — извлечение из встроенного `fs.FS` выполняет `gosod` (`github.com/leaanthony/gosod`).

### Флаги `wails3 init`

Определены в `internal/flags/init.go`:

| Флаг | Назначение | Значение по умолчанию |
| --- | --- | --- |
| `-p` | Имя пакета | `main` |
| `-t` | Имя встроенного шаблона, локальный путь или URL | `vanilla` |
| `-n` | Имя проекта | (пусто) |
| `-d` | Каталог проекта | `.` |
| `-q` | Отключить вывод в консоль | false |
| `-l` | Вывести список шаблонов | false |
| `-skipgomodtidy` | Не запускать `go mod tidy` после извлечения | false |
| `-git` | URL репозитория Git для инициализации | (пусто) |
| `-mod` | Путь модуля Go (если не задан, определяется по `-git`) | (пусто) |
| `-s` | Не показывать предупреждение при использовании удалённых шаблонов | false |
| `-productname` / `-productdescription` / `-productversion` / `-productcompany` / `-productcopyright` / `-productcomments` / `-productidentifier` | Метаданные, встроенные в сгенерированные ресурсы сборки | разумные значения по умолчанию |

Длинного псевдонима `-list` **нет** (есть только `-l`), как нет и отдельного `--help` для каждого шаблона.

### Подстановки

Заполнители представляют собой стандартные директивы шаблонов Go — начальный `.` является частью средства доступа к полю:

| Заполнитель | Пример | Источник |
| --- | --- | --- |
| `{{.ProjectName}}` | `myapp` | Флаг `-n` / имя каталога |
| `{{.ModulePath}}` | `github.com/me/myapp` | Флаг `-mod` или значение, полученное из `-git` |
| `{{.WailsVersion}}` | `v3.0.0-…` | Встроенная при компиляции константа из `internal/version` |
| `{{.ProductName}}`, `{{.ProductDescription}}`, `{{.ProductVersion}}`, `{{.ProductCompany}}`, `{{.ProductCopyright}}`, `{{.ProductComments}}`, `{{.ProductIdentifier}}` | Метаданные, добавляемые при сборке | соответствующие флаги `-product*` |

Если нужен новый заполнитель, добавьте поле в данные шаблона в `internal/templates/templates.go` и соответствующее поле или флаг в `internal/flags/init.go` (либо задайте его из `internal/commands/init.go`).

### Обработчик после копирования

После того как `gosod` завершит извлечение шаблона, CLI выполнит:

```
go mod tidy
```

если только не передан `-skipgomodtidy`. Этапа `task deps` нет.

---

## 3. Создание нового шаблона

> Пример: добавление шаблона **Solid**

### 3.1 Папка и идентификатор

```
internal/templates/solid/
```

Имя папки совпадает с идентификатором шаблона. Используйте формат **kebab-case**.

### 3.2 Минимальный набор файлов

```
solid/
├── template.yaml    # name, description, wailsVersion, typescript (required)
├── frontend/        # Your web project (no node_modules/dist)
│   ├── src/
│   ├── package.json
│   └── vite.config.ts
└── ...              # Any extra Go files the template wants to inject
```

Сначала скопируйте `react` и удалите ненужные файлы. Не забудьте создать `template.yaml` — для шаблона TypeScript задайте `typescript: true`. `base/` — единственная папка без такого файла.

### 3.3 Обновление заполнителей

Найдите и замените буквальные демонстрационные значения директивами шаблонов Go, например:

- `myapp` → `{{.ProjectName}}`
- `github.com/you/myapp` → `{{.ModulePath}}`

### 3.4 Подключение

Поскольку `templates.go` обходит встроенную файловую систему при инициализации, обычно достаточно добавить новую папку в `internal/templates/<id>/` — вручную регистрировать её не требуется. Если нужна дополнительная логика (пользовательская проверка или действия после копирования), добавьте её в `templates.Install` в `internal/templates/templates.go`.

### 3.5 Тестирование

```bash
wails3 init -n demo -t solid
cd demo
wails3 dev
```

Убедитесь, что:

- Сервер разработки запускается на порте, указанном в `WAILS_VITE_PORT`
- Сгенерированные привязки появляются в `frontend/bindings/...`
- Горячая перезагрузка работает

---

## 4. Изменение существующих шаблонов

1. Измените файлы в `internal/templates/<id>/`.
2. Пересоберите CLI (`cd v3 && go build -o ../wails3 ./cmd/wails3`); директива `//go:embed *` подхватит новое содержимое.
3. Обновите **версии зависимостей** в `frontend/package.json` и `Taskfile.yml`.
4. Если поведение изменилось, обновите описание шаблона в `template.json`.

### Типичные изменения

| Задача | Где |
| --- | --- |
| Изменить порт сервера разработки | `frontend/vite.config.ts` — считывайте `WAILS_VITE_PORT` |
| Добавить переменные окружения | `build/Taskfile.yml` или `frontend/.env` |
| Заменить менеджер пакетов JS | Замените `npm` → `pnpm`/`bun` в `build/Taskfile.yml` |

---

## 5. Рекомендации по созданию шаблонов

- **Сохраняйте фронтенд универсальным** — не обращайтесь к глобальным объектам, специфичным для Wails; `/wails/runtime.js` предоставляется сервером ресурсов во время выполнения.
- **Не добавляйте скомпилированные артефакты** — исключите `node_modules`, `dist` и `.DS_Store` из встроенного каталога (либо добавьте их в `.gitignore`, чтобы они никогда не попадали в коммиты).
- **Документируйте предварительные требования** — версию Node, дополнительные инструменты CLI и прочее — в `template.json` или `NEXTSTEPS.md`.
- **Избегайте несовместимых изменений** — если переработка значительна, создайте новый идентификатор шаблона вместо изменения существующего.

---

## 6. Устранение неполадок

| Симптом | Причина | Решение |
| --- | --- | --- |
| `unknown template name` | Опечатка в `-t` или шаблон не встроен | Выполните `wails3 init -l`, чтобы вывести список доступных шаблонов |
| Заполнители не заменяются | Использован `{{ProjectName}}` вместо `{{.ProjectName}}` | Добавьте начальный `.` (доступ к полю шаблона Go) |
| Сервер разработки открывает пустую страницу | Конфигурация Vite не считывает `WAILS_VITE_PORT` | Проверьте файл `vite.config.ts` |
| Не удаётся собрать фронтенд для рабочей среды | Не указан путь `base` в Vite | Задайте `base: "./"` в `vite.config.ts` |

---

## 7. Карта основных файлов исходного кода

| Файл | Назначение |
| --- | --- |
| `internal/templates/templates.go` | Встраивает файловую систему шаблонов, предоставляет `Install(options *flags.Init) error`, `GetDefaultTemplates()` и `ValidTemplateName(name)` |
| `internal/templates/<id>/**` | Фактическое содержимое шаблона |
| `internal/commands/init.go` | Связующий код CLI: выбирает шаблон, заполняет метаданные и вызывает `templates.Install` |
| `internal/commands/generate_template.go` | `wails3 generate template` — утилита для *экспорта* рабочего проекта обратно в шаблон (удобно для обновлений) |
| `internal/flags/init.go` | Определения флагов для `wails3 init` |

---

## 8. Итоги

- Шаблоны находятся в **`internal/templates/`** и встраиваются в CLI посредством `//go:embed *`.
- `wails3 init -t <id>` извлекает шаблон посредством `gosod` и запускает `go mod tidy` (этот шаг можно пропустить с помощью `-skipgomodtidy`).
- Чтобы создать шаблон, достаточно **создать папку**, добавить файлы и `template.json`, а затем использовать заполнители в стиле `{{.ProjectName}}`.
- Система **расширяема** и **самодостаточна** — она идеально подходит для предоставления собственных наборов технологий вашей команде или сообществу.

Успешной работы с шаблонами!
