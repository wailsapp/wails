---
title: "Начало работы"
description: "Как начать вносить вклад в Wails v3"
slug: "contributing/getting-started"
sourcePath: "contributing/getting-started.md"
---

## Добро пожаловать!

Благодарим за интерес к участию в разработке Wails! Это руководство поможет вам внести свой первый вклад.

## Предварительные требования

Перед началом работы убедитесь, что у вас есть:

- установленный Go версии **1.26 или новее** ([скачать](https://go.dev/dl/))
- **Node.js версии 20 или новее** и **npm** ([скачать](https://nodejs.org/))
- **Git**, настроенный для вашей учётной записи GitHub
- Базовые знания Go и JavaScript/TypeScript

### Требования для отдельных платформ

**macOS:**

- Инструменты командной строки Xcode: `xcode-select --install`

**Windows:**

- Рекомендуется MSYS2 или аналогичная Unix-подобная среда
- Среда выполнения WebView2 (обычно предустановлена в Windows 11)

**Linux:**

- `gcc`, `pkg-config`, `libgtk-4-dev`, `libwebkitgtk-6.0-dev` (стандартный стек GTK4)
- Установка с помощью: `sudo apt install build-essential pkg-config libgtk-4-dev libwebkitgtk-6.0-dev` (Debian/Ubuntu)
- Для устаревшего варианта сборки `-tags gtk3` также установите `libgtk-3-dev` и `libwebkit2gtk-4.1-dev`

## Обзор процесса внесения изменений

Обычный процесс внесения изменений состоит из следующих этапов:

1. **Создайте форк и клонируйте его** — создайте собственную копию репозитория Wails
2. **Настройте среду** — соберите CLI Wails и проверьте свою среду
3. **Создайте ветку** — создайте функциональную ветку для своих изменений
4. **Разработайте изменения** — внесите изменения в соответствии с нашими стандартами оформления кода
5. **Протестируйте** — запустите тесты и убедитесь, что всё работает
6. **Создайте коммит** — зафиксируйте изменения с понятными сообщениями в формате Conventional Commits
7. **Отправьте изменения** — откройте pull request для проверки
8. **Доработайте изменения** — учтите обратную связь и внесите необходимые корректировки
9. **Объедините изменения** — после одобрения ваши изменения станут частью Wails!

## Пошаговое руководство

Выберите тип вклада:

@tabs
[Исправление ошибки]
@steps
### Найдите ошибку или сообщите о ней
- Проверьте, не зарегистрирована ли уже эта ошибка в [GitHub Issues](https://github.com/wailsapp/wails/issues)
- Если нет, создайте новую задачу и укажите шаги для воспроизведения
- Прежде чем приступать к работе, дождитесь подтверждения

### Создайте форк и клонируйте его
Создайте форк репозитория на странице [github.com/wailsapp/wails/fork](https://github.com/wailsapp/wails/fork)

Клонируйте свой форк:

```bash
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails
git remote add upstream https://github.com/wailsapp/wails.git
```

### Соберите и проверьте
Соберите Wails и убедитесь, что можете воспроизвести ошибку:

```bash
cd v3
go build -o ../wails3 ./cmd/wails3

# Reproduce the bug to understand it
```

### Создайте ветку для исправления ошибки
Создайте ветку для исправления:

```bash
git checkout -b fix/issue-123-window-crash
```

### Исправьте ошибку
- Внесите минимальные изменения, необходимые для исправления ошибки
- Не выполняйте рефакторинг несвязанного кода
- Добавьте или обновите тесты, чтобы предотвратить регрессию

```bash
# Make your changes
# Add tests in *_test.go files
```

### Протестируйте исправление
Запустите тесты и убедитесь, что исправление работает:

```bash
go test ./...

# Test the specific package
go test ./pkg/application -v

# Run with race detector
go test ./... -race
```

### Зафиксируйте исправление
Создайте коммит с понятным сообщением:

```bash
git commit -m "fix: prevent window crash when closing during initialization

Fixes #123"
```

### Отправьте pull request
Отправьте изменения и создайте pull request:

```bash
git push origin fix/issue-123-window-crash
```

В описании pull request:

- Объясните суть и первопричину ошибки
- Опишите своё исправление
- Укажите задачу: «Fixes #123»
- Опишите поведение до и после исправления

### Ответьте на замечания
Учтите замечания рецензентов и при необходимости обновите pull request.

@end

[WEP (улучшение)]
@steps
### Напишите WEP
- Ознакомьтесь с [процессом WEP (Wails Enhancement Proposal)](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/README.md)
- Скопируйте шаблон WEP в `v3/wep/proposals/<name>/proposal.md`
- Откройте черновик pull request с заголовком `[WEP] <title>`, содержащий только WEP
- Прежде чем приступать к реализации, дождитесь решения сопровождающего

### Создайте форк и клонируйте его
Создайте форк репозитория на странице [github.com/wailsapp/wails/fork](https://github.com/wailsapp/wails/fork)

Клонируйте свой форк:

```bash
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails
git remote add upstream https://github.com/wailsapp/wails.git
```

### Настройте среду разработки
Соберите Wails и проверьте свою среду:

```bash
cd v3
go build -o ../wails3 ./cmd/wails3

# Run tests to ensure everything works
go test ./...
```

### Создайте функциональную ветку
Создайте ветку с понятным названием:

```bash
git checkout -b feat/window-transparency-support
```

### Реализуйте функцию
- Следуйте нашим [стандартам оформления кода](/contributing/standards/)
- Ограничьте изменения реализацией этой функции
- Пишите чистый и документированный код
- Добавьте исчерпывающие тесты

```bash
# Example: Adding a new window method
# 1. Add to window.go interface
# 2. Implement in platform files (darwin, windows, linux)
# 3. Add tests
# 4. Update documentation
```

### Проведите тщательное тестирование
Протестируйте функцию:

```bash
# Unit tests
go test ./pkg/application -v

# Integration test - create a test app
cd ..
./wails3 init -n feature-test
cd feature-test
# Add code using your new feature
../wails3 dev
```

### Документируйте функцию
- Добавьте строки документации ко всем публичным API
- Обновите соответствующую документацию в `/docs/mpress/content/`
- Если применимо, добавьте примеры

### Используйте соглашение о коммитах
Используйте формат Conventional Commits:

```bash
git commit -m "feat: add window transparency support

- Add SetTransparent() method to Window API
- Implement for macOS, Windows, and Linux
- Add tests and documentation

Closes #456"
```

### Отправьте pull request
Отправьте изменения и создайте pull request:

```bash
git push origin feat/window-transparency-support
```

В pull request:

- Опишите функцию и сценарии её использования
- Приведите примеры или снимки экрана
- Перечислите все несовместимые изменения
- Укажите принятый pull request с WEP

### Доработайте по результатам проверки
Сопровождающие могут попросить внести изменения. Проявляйте терпение и будьте готовы к совместной работе.

@end

[Документация]
PR с исправлениями документации можно отправлять без предварительного создания issue. Если изменения затрагивают только документацию, тест кода, воспроизводящий ошибку, не требуется. Следуйте руководству [«Работа с документацией в M-Press»](/contributing/documentation/), где описаны установка M-Press, пути к исходным файлам, предварительный просмотр, проверка и отправка PR.

@end

## Поиск задач для работы

- Ищите метки [`good first issue`](https://github.com/wailsapp/wails/labels/good%20first%20issue)
- Просматривайте issue [`help wanted`](https://github.com/wailsapp/wails/labels/help%20wanted)
- Просмотрите [открытые issue](https://github.com/wailsapp/wails/issues) и попросите назначить выбранную задачу вам

## Получение помощи

- **Discord:** присоединитесь к [серверу Wails в Discord](https://discord.gg/JDdSxwjhGf)
- **Обсуждения:** создайте публикацию в [GitHub Discussions](https://github.com/wailsapp/wails/discussions)
- **Issue:** создавайте issue для воспроизводимых ошибок; для вопросов используйте Discussions, а для улучшений — PR с WEP

## Кодекс поведения

Относитесь к другим с уважением, общайтесь конструктивно и будьте доброжелательны. Вместе мы создаём дружелюбное сообщество, нацеленное на разработку отличного программного обеспечения.

## Дальнейшие действия

- Настройте [среду разработки](/contributing/setup/)
- Ознакомьтесь с нашими [стандартами программирования](/contributing/standards/)
- Изучите [техническую документацию](/contributing/overview/)
