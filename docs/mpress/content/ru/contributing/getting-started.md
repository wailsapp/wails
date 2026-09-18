---
title: "Начало работы"
description: "Как начать вносить вклад в Wails v3"
slug: "contributing/getting-started"
sourcePath: "contributing/getting-started.md"
---

## Добро пожаловать в сообщество участников!

Спасибо за интерес к участию в разработке Wails! Это руководство поможет вам внести свой первый вклад.

## Предварительные требования

Прежде чем начать, убедитесь, что у вас есть:

- Установленный **Go 1.25+** ([скачать](https://go.dev/dl/))
- **Node.js 20+** и **npm** ([скачать](https://nodejs.org/))
- **Git**, настроенный для вашей учётной записи GitHub
- Базовое знание Go и JavaScript/TypeScript

### Требования для отдельных платформ

**macOS:**

- Инструменты командной строки Xcode: `xcode-select --install`

**Windows:**

- Рекомендуется MSYS2 или аналогичная Unix-подобная среда
- Среда выполнения WebView2 (обычно уже установлена в Windows 11)

**Linux:**

- `gcc`, `pkg-config`, `libgtk-4-dev`, `libwebkitgtk-6.0-dev` (стек GTK4 по умолчанию)
- Установите с помощью: `sudo apt install build-essential pkg-config libgtk-4-dev libwebkitgtk-6.0-dev` (Debian/Ubuntu)
- Для устаревшего варианта сборки `-tags gtk3` также установите `libgtk-3-dev` и `libwebkit2gtk-4.1-dev`

## Обзор процесса внесения вклада

Обычно процесс внесения вклада состоит из следующих этапов:

1. **Создание форка и клонирование** — создайте собственную копию репозитория Wails
2. **Настройка** — соберите Wails CLI и проверьте свою среду
3. **Ветка** — создайте ветку для своих изменений
4. **Разработка** — внесите изменения в соответствии с нашими стандартами написания кода
5. **Тестирование** — запустите тесты, чтобы убедиться, что всё работает
6. **Коммит** — зафиксируйте изменения с понятными сообщениями в формате Conventional Commits
7. **Отправка** — откройте пул-реквест для проверки
8. **Доработка** — отвечайте на замечания и вносите необходимые изменения
9. **Слияние** — после одобрения ваши изменения станут частью Wails!

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
Создайте ветку для своего исправления:

```bash
git checkout -b fix/issue-123-window-crash
```

### Исправьте ошибку
- Внесите только минимально необходимые для исправления ошибки изменения
- Не выполняйте рефакторинг кода, не связанного с ошибкой
- Добавьте или обновите тесты, чтобы предотвратить регрессию

```bash
# Make your changes
# Add tests in *_test.go files
```

### Протестируйте исправление
Запустите тесты, чтобы убедиться, что исправление работает:

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

### Отправьте пул-реквест
Отправьте ветку и создайте пул-реквест:

```bash
git push origin fix/issue-123-window-crash
```

В описании пул-реквеста:

- Объясните суть и первопричину ошибки
- Опишите своё исправление
- Добавьте ссылку на задачу: «Fixes #123»
- Опишите поведение до и после исправления

### Ответьте на замечания
Учтите замечания проверяющих и при необходимости обновите пул-реквест.

@end

[WEP (предложение по улучшению)]
@steps
### Подготовьте WEP
- Ознакомьтесь с [процессом WEP (предложения по улучшению Wails)](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md)
- Скопируйте шаблон WEP в `v3/wep/proposals/<name>/proposal.md`
- Откройте черновик пул-реквеста с заголовком `[WEP] <title>`, содержащий только WEP
- Прежде чем приступать к реализации, дождитесь решения сопровождающего

### Создайте форк и клонируйте его
Создайте форк репозитория на странице [github.com/wailsapp/wails/fork](https://github.com/wailsapp/wails/fork)

Клонируйте свой форк:

```bash
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails
git remote add upstream https://github.com/wailsapp/wails.git
```

### Настройка среды разработки
Соберите Wails и проверьте среду:

```bash
cd v3
go build -o ../wails3 ./cmd/wails3

# Run tests to ensure everything works
go test ./...
```

### Создание ветки для функции
Создайте ветку с информативным названием:

```bash
git checkout -b feat/window-transparency-support
```

### Реализация функции
- Следуйте нашим [стандартам написания кода](/contributing/standards/)
- Вносите только изменения, относящиеся к функции
- Пишите чистый, документированный код
- Добавьте исчерпывающие тесты

```bash
# Example: Adding a new window method
# 1. Add to window.go interface
# 2. Implement in platform files (darwin, windows, linux)
# 3. Add tests
# 4. Update documentation
```

### Тщательное тестирование
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

### Документирование функции
- Добавьте строки документации ко всем общедоступным API
- Обновите соответствующую документацию в `/docs/mpress/content/`
- Если применимо, добавьте примеры

### Оформление коммитов по соглашению
Используйте соглашение Conventional Commits:

```bash
git commit -m "feat: add window transparency support

- Add SetTransparent() method to Window API
- Implement for macOS, Windows, and Linux
- Add tests and documentation

Closes #456"
```

### Отправка запроса на включение изменений
Отправьте ветку в удалённый репозиторий и создайте PR:

```bash
git push origin feat/window-transparency-support
```

В своём PR:

- Опишите функцию и сценарии её использования
- Приведите примеры или снимки экрана
- Перечислите все несовместимые изменения
- Добавьте ссылку на принятый PR с WEP

### Доработка по результатам проверки
Сопровождающие могут запросить изменения. Проявляйте терпение и будьте готовы к совместной работе.

@end

[Документация]
PR с исправлениями можно отправлять без предварительного создания issue. Для исправлений, затрагивающих только документацию, не требуется тест кода, который завершается с ошибкой. Инструкции по установке M-Press, путям к исходным файлам, предварительному просмотру, проверке и созданию PR приведены в разделе [«Исправление документации»](/contributing/documentation/).

@end

## Поиск задач для работы

- Ищите метки [`good first issue`](https://github.com/wailsapp/wails/labels/good%20first%20issue)
- Проверяйте issue [`help wanted`](https://github.com/wailsapp/wails/labels/help%20wanted)
- Просмотрите [открытые issue](https://github.com/wailsapp/wails/issues) и попросите назначить одну из них вам

## Получение помощи

- **Discord:** присоединитесь к [серверу Wails в Discord](https://discord.gg/JDdSxwjhGf)
- **Обсуждения:** опубликуйте сообщение в [GitHub Discussions](https://github.com/wailsapp/wails/discussions)
- **Issue:** создайте issue для воспроизводимой ошибки; для вопросов используйте Discussions, а для улучшений — PR с WEP

## Кодекс поведения

Проявляйте уважение, конструктивность и доброжелательность. Мы создаём дружелюбное сообщество, участники которого вместе работают над отличным программным обеспечением.

## Следующие шаги

- Настройте [среду разработки](/contributing/setup/)
- Ознакомьтесь с нашими [стандартами написания кода](/contributing/standards/)
- Изучите [техническую документацию](/contributing/overview/)
