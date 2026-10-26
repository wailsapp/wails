---
title: "Система шаблонов"
description: "Как Wails v3 создаёт каркас новых проектов, как организованы шаблоны и как создать собственный."
slug: "contributing/template-system"
sourcePath: "contributing/template-system.md"
---

Wails поставляется с **системой шаблонов**, с помощью которой `wails3 init` создаёт готовый к запуску проект. Встроенные шаблоны предусмотрены для намеренно ограниченного набора фреймворков (Vanilla, React, Vue, Svelte); любой другой фреймворк можно использовать, [добавив собственный фронтенд](/guides/dev/frontend-frameworks/) или опубликовав [пользовательский шаблон](/guides/advanced/custom-templates/).

На этой странице рассматриваются:

1. Структура каталогов шаблонов
2. Как CLI выбирает и формирует шаблоны
3. Пошаговое создание нового шаблона
4. Обновление или переопределение существующих шаблонов
5. Устранение неполадок и рекомендации

---

## 1. Где находятся шаблоны

```
v3/internal/templates/
├── _common/        # Files copied into EVERY project (Taskfile.yml, build/, etc.)
├── base/           # Minimal frontend skeleton (frontend/ + NEXTSTEPS.md)
├── ios/            # iOS bootstrapper
├── vanilla/        vanilla-js/   # TypeScript (default) + JavaScript variant
├── react/          react-js/     # TypeScript (default) + JavaScript variant
├── vue/                          # TypeScript only
├── svelte/                       # TypeScript only
└── templates.go    # Registry + Install/Get APIs (no auto-registration via embed)
```

- **`_common/`** — универсальная заготовка (Taskfile, каталог `build/` и общая инфраструктура), добавляемая в каждый проект.
- **`base/`** — минимальный каркас фронтенда для пользовательских шаблонов, а не автоматически добавляемый слой бэкенда. При обычном создании рабочей копии проекта содержимое шаблона выбранного фреймворка накладывается поверх `_common/`.
- **Каталоги фреймворков** — содержат фронтенд (`frontend/`), конфигурацию фреймворка и файл `template.yaml` с метаданными шаблона.
- Имена каталогов соответствуют **идентификатору шаблона**, который передаётся CLI (`wails3 init -t react`).
- **Соглашение о языках:** по умолчанию используется TypeScript, и соответствующий шаблон получает имя без суффикса (`react`); вариант на JavaScript, если он существует, получает суффикс `-js` (`react-js`). Встроенные шаблоны явно указывают язык с помощью `typescript: true|false` в `template.yaml`. Шаблоны сообщества по-прежнему могут использовать устаревший суффикс `-ts`, который поддерживается как резервный вариант.

> Весь каталог `internal/templates/` встраивается в двоичный файл CLI
>
> с помощью `//go:embed *`, поэтому пользователи могут создавать каркас проектов без подключения к сети.

---

Загрузчик также принимает метаданные в устаревшем формате `template.json`. Новые встроенные шаблоны используют YAML, включая явно заданный параметр `typescript`.

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

API `Template.Load()` / `Template.CopyTo()` / `Template.Validate()` отсутствует — создание рабочей копии из встроенного `fs.FS` выполняет `gosod` (`github.com/leaanthony/gosod`).

### Флаги `wails3 init`

Определены в `internal/flags/init.go`:

| Флаг | Назначение | Значение по умолчанию |
| --- | --- | --- |
| `-p` | Имя пакета | `main` |
| `-t` | Имя встроенного шаблона, локальный путь или URL | `vanilla` |
| `-n` | Имя проекта | (пусто) |
| `-d` | Каталог проекта | `.` |
| `-q` | Не выводить сообщения в консоль | false |
| `-l` | Вывести список шаблонов | false |
| `-skipgomodtidy` | Не запускать `go mod tidy` после создания рабочей копии | false |
| `-git` | URL репозитория Git для инициализации | (пусто) |
| `-mod` | Путь модуля Go (если не задан, определяется на основе `-git`) | (пусто) |
| `-s` | Не показывать предупреждение при использовании удалённых шаблонов | false |
| `-productname` / `-productdescription` / `-productversion` / `-productcompany` / `-productcopyright` / `-productcomments` / `-productidentifier` | Метаданные, встроенные в создаваемые ресурсы сборки | разумные значения по умолчанию |

Длинного псевдонима `-list` нет (поддерживается только `-l`), как и отдельного параметра `--help` для каждого шаблона.

### Подстановки

Заполнители представляют собой стандартные директивы шаблонов Go — начальный `.` является частью выражения доступа к полю:

| Заполнитель | Пример | Источник |
| --- | --- | --- |
| `{{.ProjectName}}` | `myapp` | Флаг `-n` / имя каталога |
| `{{.ModulePath}}` | `github.com/me/myapp` | Флаг `-mod` или значение, полученное из `-git` |
| `{{.WailsVersion}}` | `v3.0.0` | Константа, встраиваемая при компиляции из `internal/version` |
| `{{.ProductName}}`, `{{.ProductDescription}}`, `{{.ProductVersion}}`, `{{.ProductCompany}}`, `{{.ProductCopyright}}`, `{{.ProductComments}}`, `{{.ProductIdentifier}}` | Метаданные, встраиваемые при сборке | соответствующие флаги `-product*` |

Если нужен новый заполнитель, добавьте поле в данные шаблона в `internal/templates/templates.go` и соответствующее поле или флаг в `internal/flags/init.go` (либо задайте его из `internal/commands/init.go`).

### Хук после копирования

Когда `gosod` завершает создание рабочей копии из шаблона, CLI запускает:

```
go mod tidy
```

если не передан `-skipgomodtidy`. Этап `task deps` отсутствует.

---

## 3. Создание нового шаблона

> Пример: добавление шаблона **Solid**

### 3.1 Каталог и идентификатор

```
internal/templates/solid/
```

Имя каталога совпадает с идентификатором шаблона. Используйте **kebab-case**.

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

Для начала скопируйте `react` и удалите ненужные файлы. Не забудьте создать `template.yaml` — для шаблона TypeScript задайте `typescript: true`. `base/` — единственная папка без этого файла.

### 3.3 Обновление заполнителей

Найдите и замените буквальные значения примеров директивами шаблонов Go, например:

- `myapp` → `{{.ProjectName}}`
- `github.com/you/myapp` → `{{.ModulePath}}`

### 3.4 Подключение

Поскольку `templates.go` при инициализации обходит встроенную файловую систему, обычно достаточно добавить новую папку в `internal/templates/<id>/` — вручную вызывать регистрацию не требуется. Если нужна дополнительная логика, например специальная проверка или действия после копирования, добавьте её в `templates.Install` в файле `internal/templates/templates.go`.

### 3.5 Тестирование

```bash
wails3 init -n demo -t solid
cd demo
wails3 dev
```

Убедитесь, что:

- Сервер разработки запускается на порту, указанном в `WAILS_VITE_PORT`
- Сгенерированные привязки появляются в `frontend/bindings/...`
- Горячая перезагрузка работает

---

## 4. Изменение существующих шаблонов

1. Измените файлы в `internal/templates/<id>/`.
2. Пересоберите CLI (`cd v3 && go build -o ../wails3 ./cmd/wails3`); директива `//go:embed *` включит новое содержимое.
3. Обновите версии **зависимостей** в `frontend/package.json` и `Taskfile.yml`.
4. Если поведение изменилось, обновите описание шаблона в `template.yaml`.

### Распространённые изменения

| Задача | Где изменить |
| --- | --- |
| Изменить порт сервера разработки | `frontend/vite.config.ts` — чтение `WAILS_VITE_PORT` |
| Добавить переменные окружения | `build/Taskfile.yml` или `frontend/.env` |
| Заменить менеджер пакетов JavaScript | В `build/Taskfile.yml` заменить `npm` на `pnpm`/`bun` |

---

## 5. Рекомендации по созданию шаблонов

- **Сохраняйте универсальность фронтенда** — не обращайтесь к глобальным объектам, специфичным для Wails; во время выполнения `/wails/runtime.js` предоставляется сервером ресурсов.
- **Не добавляйте скомпилированные артефакты** — исключите `node_modules`, `dist` и `.DS_Store` из встраиваемого каталога (либо добавьте их в `.gitignore`, чтобы они никогда не попадали в коммиты).
- **Документируйте предварительные требования** — версию Node, дополнительные инструменты CLI и прочее — в `template.yaml` или `NEXTSTEPS.md`.
- **Избегайте несовместимых изменений** — если переработка значительна, создайте новый идентификатор шаблона вместо изменения существующего.

---

## 6. Устранение неполадок

| Симптом | Причина | Решение |
| --- | --- | --- |
| `unknown template name` | Опечатка в `-t` или шаблон не встроен | Выполните `wails3 init -l`, чтобы вывести список доступных шаблонов |
| Заполнители не заменены | Использован `{{ProjectName}}` вместо `{{.ProjectName}}` | Добавьте начальный `.` (доступ к полю шаблона Go) |
| Сервер разработки открывает пустую страницу | Конфигурация Vite не считывает `WAILS_VITE_PORT` | Проверьте `vite.config.ts` |
| Не удаётся собрать фронтенд для рабочей среды | Не указан путь `base` в Vite | Задайте `base: "./"` в `vite.config.ts` |

---

## 7. Карта основных исходных файлов

| Файл | Назначение |
| --- | --- |
| `internal/templates/templates.go` | Встраивает файловую систему шаблонов и предоставляет `Install(options *flags.Init) error`, `GetDefaultTemplates()` и `ValidTemplateName(name)` |
| `internal/templates/<id>/**` | Фактическое содержимое шаблонов |
| `internal/commands/init.go` | Связующий код CLI: выбирает шаблон, заполняет метаданные и вызывает `templates.Install` |
| `internal/commands/generate_template.go` | `wails3 generate template` — утилита для *экспорта* рабочего проекта обратно в шаблон (удобно для обновлений) |
| `internal/flags/init.go` | Определения флагов для `wails3 init` |

---

## 8. Итоги

- Шаблоны находятся в **`internal/templates/`** и встраиваются в CLI с помощью `//go:embed *`.
- `wails3 init -t <id>` извлекает шаблон с помощью `gosod` и выполняет `go mod tidy` (этот шаг можно пропустить с помощью `-skipgomodtidy`).
- Чтобы создать шаблон, достаточно **создать папку**, добавить файлы и `template.yaml`, а затем использовать заполнители в стиле `{{.ProjectName}}`.
- Система **расширяема** и **самодостаточна**, поэтому идеально подходит для распространения пользовательских стеков среди вашей команды или сообщества.

Удачной работы с шаблонами!
