---
title: "Тестирование и непрерывная интеграция"
description: "Как Wails v3 обеспечивает качество с помощью модульных и интеграционных тестов, обнаружения гонок данных и CI на базе GitHub Actions."
slug: "contributing/testing-ci"
sourcePath: "contributing/testing-ci.md"
---

Надёжным фреймворкам для настольных приложений требуется безотказное тестирование. В Wails v3 применяется **многоуровневая стратегия**:

| Уровень | Цель | Инструменты |
| --- | --- | --- |
| Модульные тесты | Быстрая обратная связь для изолированных функций | `go test ./...` |
| Тесты генератора и CLI | Проверка `wails3 generate bindings` и внутренней инфраструктуры CLI | `task test:generator`, `task test:cli` |
| Тесты шаблонов | Проверка того, что каждый поставляемый шаблон по-прежнему собирается | `task test:templates` |
| Обнаружение гонок данных | Выявление гонок данных в среде выполнения и мосте | `go test -race ./...` |
| Матрица CI | Уверенность в работе на разных ОС для каждого PR | GitHub Actions |

В этом документе объясняется, **где находятся тесты**, **как их запускать** и **какими операциями управляет Taskfile**.

> Руководство по гонкам данных, на которое в более ранних черновиках ссылались как на `pkg/application/RACE.md`,
>
> сейчас находится в `v3/TESTING.md`.

---

## 1. Соглашения о структуре каталогов

```
v3/
├── internal/.../_test.go     # Unit tests for internal packages
├── pkg/.../_test.go          # Public API tests
├── tasks/events/generate.go  # Code generator for event constants (NOT a test harness)
├── tests/                    # Top-level integration test harness
└── TESTING.md                # Race / Cgo testing guidance
```

Рекомендации:

- **Размещайте модульные тесты рядом с кодом** (`foo.go` ↔ `foo_test.go`).
- Для пакетов `pkg/` используйте **тестирование по принципу чёрного ящика** (`package application_test`), если это помогает поддерживать чистоту API.
- Общие фикстуры размещайте там, где они используются (в этом рабочем дереве нет централизованного пакета `internal/testutil/` — вместо этого подключайте вспомогательные средства отдельно для каждого пакета).

---

## 2. Модульные тесты

### Написание тестов

```go
func TestEventConstants(t *testing.T) {
    assert.NotEmpty(t, events.Common.WindowFocus)
}
```

Рекомендации:

- Используйте [`stretchr/testify`](https://github.com/stretchr/testify) — он уже есть в `go.mod`.
- Отдавайте предпочтение **табличным** тестам, когда несколько входных данных, граничных случаев или ожидаемых результатов проверяют одно и то же поведение. Дайте каждому случаю описательное имя.
- При необходимости подменяйте платформозависимое поведение заглушками с помощью тегов сборки (`foo_windows_test.go`, `foo_darwin_test.go`, …).

### Требования к покрытию

Для новой и изменённой логики ожидается покрытие операторов Go на уровне 100%. Измеряйте покрытие изменённого пакета, а не полагайтесь на процент покрытия всего репозитория:

```bash
go test ./pkg/application -coverprofile=coverage.out
go tool cover -func=coverage.out
```

Некоторые пути выполнения практически не поддаются тестированию в обычной тестовой среде — например, сбои, возникающие только на определённой платформе, поведение, зависящее от оборудования, или защитный резервный путь, который нельзя безопасно активировать. Сводите такие исключения к минимуму и описывайте каждый непокрытый путь в описании PR.

### Локальный запуск

```bash
cd v3
go test ./... -cover
```

Также можно использовать Taskfile (указывайте реальные цели — сокращения `task test` не существует):

```
task test:cli            # CLI plumbing tests
task test:generator      # bindings generator round-trip tests
task test:templates      # build every shipped template
task test:infrastructure # supporting helpers
task test:examples       # exercise the example matrix (downloads as needed)
task test:all            # everything above
task sanity              # quick smoke check (also: sanity:gtk4)
task precommit           # what you should run before pushing
```

---

## 3. Интеграционные тесты

В `v3/tests/` размещена инфраструктура межпакетных интеграционных тестов. Цели Taskfile `test:example:*` и `test:examples:*` запускают проверки сборки и запуска в darwin / windows / linux (включая матрицы GTK3 / GTK4 на базе Docker в Linux).

> Готовые к запуску примеры находятся в `v3/examples/`. Цели тестирования выбирают и собирают
>
> примеры, соответствующие платформе хоста или матрице CI.

Чтобы запустить набор дымовых тестов для платформы хоста, выполните:

```
task test:examples       # host
task test:examples:all   # full matrix (slow)
```

---

## 4. Обнаружение гонок данных

Гонки данных приводят к фатальным сбоям в средах выполнения графического интерфейса.

### Руководство по гонкам данных

В `v3/TESTING.md` описаны:

- Известные безопасные гонки данных и обоснование их подавления
- Порядок интерпретации трассировок стека, пересекающих границы Cgo (Linux GTK + WebKit2GTK)

### Локальный набор тестов на гонки данных

```
go test -race ./...
```

> У `wails3 dev` нет флага `-race` — поддерживаются следующие флаги CLI: `--config`, `--port`
>
> и `-s` (включает HTTPS). Чтобы проверить среду выполнения с помощью детектора гонок данных,
>
> соберите тестовое приложение с `go build -race` и запустите его напрямую.

---

## 5. Рабочие процессы GitHub Actions

Фактические файлы рабочих процессов в `.github/workflows/` (сверено с рабочим деревом):

| Файл | Назначение |
| --- | --- |
| `build-and-test-v3.yml` | Основная матрица сборки и тестирования v3. Использует `actions/setup-go@v5` с `go-version: 1.25`. Запускает `task runtime:check`, `task runtime:test`, `task runtime:build`, `task test:examples` (а для пути GTK4 также `BUILD_TAGS=gtk4 task test:examples`), `task generator:test:check`, `task install`, а затем `wails3 build` в качестве дымовой проверки. Задания Linux устанавливают `libgtk-3-dev libwebkit2gtk-4.1-dev libwayland-dev build-essential pkg-config xvfb x11-xserver-utils at-spi2-core xdg-desktop-portal-gtk` и запускают набор тестов в `dbus-run-session -- xvfb-run`. |
| `cross-compile-test-v3.yml` | Базовые проверки кросс-компиляции |
| `auto-changelog-v3.yml`, `changelog-v3.yml` | Автоматизация журнала изменений |
| `nightly-release-v3.yml` | Артефакты ночных выпусков v3 |
| `bump-webview2-v3.yml`, `release-webview2.yml` | Управление зависимостью WebView2 и её выпусками |
| `build-cross-image.yml` | Собирает образ контейнера с кросс-компилятором |
| `publish-npm.yml` | Публикует встроенную среду выполнения JS `@wailsio/runtime` в npm |
| `pr-master.yml` | Проверки в PR относительно ветки `master` |
| `semgrep.yml` | Статический анализ с помощью Semgrep |
| `stale-issues.yml`, `issue-labeler.yml`, `file-labeler.yml`, `claude.yml`, `generate-sponsor-image.yml`, `sync-translated-documents.yml`, `upload-source-documents.yml`, `build-and-test.yml`, `weekly-release-v2.yml` | Обслуживание репозитория и процессы для v2 |

В этом рабочем дереве **нет** `qodana.yaml` и **нет** `runtime.yml` — в более ранних черновиках этой страницы упоминались оба, однако статический анализ выполняет только `semgrep.yml`, а пакет среды выполнения JS публикуется через `publish-npm.yml`.

Шаги CI соответствуют указанным выше целям Taskfile (`task test:cli`, `task test:generator`, `task test:templates`, `task test:examples`, …), поэтому вы можете в точности воспроизвести CI локально. Шаг быстрой проверки `wails3 build` в `build-and-test-v3.yml` вызывается **без дополнительных флагов** — у `wails3 build` нет флага `-skip-package`.

---

## 6. Локальное воспроизведение CI

Единой общей цели `task ci` нет. Чтобы воспроизвести CI, последовательно запустите реальные цели:

```
task precommit
task test:cli
task test:generator
task test:templates
task test:examples
```

---

## 7. Устранение неполадок при сбоях тестов

| Симптом | Вероятная причина | Решение |
| --- | --- | --- |
| **Состояние гонки в `webview_window_darwin.go`** | Изменение состояния окна вне главного потока | Передайте вызов через `application.InvokeAsync` / `Invoke`, чтобы он выполнялся в главном потоке |
| **Тест для Linux зависает в CI без графического окружения** | Для GTK требуется дисплей | Запустите через `xvfb-run`, например: `xvfb-run task test:examples:linux` |
| **Сбой сборки шаблона** | Файл блокировки зависимостей фронтенда устарел | Повторно запустите `wails3 init` для чистого каталога, чтобы обновить шаблон |
| **Ошибки Coverpkg** | Интеграционный тест импортирует `main` | Перейдите на тег сборки `//go:build integration` и сделайте импорт условным |

---

## 8. Добавление новых тестов

1. **Модульные тесты** — создайте `*_test.go` и запустите `go test ./...`
2. **Генератор / CLI** — дополните тестовые случаи в `internal/generator/testcases/` или `internal/commands/*_test.go` и повторно запустите `task test:generator` / `task test:cli`
3. **Шаблоны / примеры** — убедитесь, что поставляемые шаблоны по-прежнему собираются с помощью `task test:templates`

---

## 9. Карта ключевых файлов

| Назначение | Путь |
| --- | --- |
| Тест полного цикла генератора | `internal/generator/generate_test.go` |
| Тест сборки ресурсов | `internal/commands/build-assets_test.go` |
| Руководство по состояниям гонки / Cgo | `v3/TESTING.md` |
| Цели тестирования в Taskfile | `v3/Taskfile.yaml` |
| Генератор констант событий | `v3/tasks/events/generate.go` |
| Рабочий процесс CI | `.github/workflows/build-and-test-v3.yml` (Go 1.25 через `actions/setup-go@v5`) |
| Статический анализ | `.github/workflows/semgrep.yml` |
| Публикация среды выполнения в npm | `.github/workflows/publish-npm.yml` |

---

В Wails v3 качество — не второстепенная задача. Благодаря модульным тестам, наборам тестов генератора и шаблонов, обнаружению состояний гонки и кроссплатформенной матрице CI вы можете уверенно вносить изменения, зная, что они успешно проходят проверку во всех поддерживаемых нами ОС. Успешного тестирования!
