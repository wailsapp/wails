---
title: "Участие в разработке"
description: "Внесите свой вклад в Wails"
slug: "contributing"
sourcePath: "contributing/index.md"
---

## Добро пожаловать, участники проекта!

Мы рады любому вкладу в Wails! Исправляете ли вы ошибки, добавляете функции или улучшаете документацию — мы ценим вашу помощь.

## Как внести вклад

### 1. Сообщайте о проблемах

Нашли ошибку? [Создайте issue](https://github.com/wailsapp/wails/issues/new) и укажите:

- Чёткое описание
- Шаги для воспроизведения
- Ожидаемое и фактическое поведение
- Сведения о системе
- Примеры кода

### 2. Улучшайте документацию

PR с исправлениями приветствуются без предварительного создания issue или падающего теста кода.  
Инструкции по скачиванию, установке в macOS/Linux/Windows, предварительному просмотру в реальном времени и отправке изменений приведены в разделе [«Создание документации с помощью M-Press»](/contributing/documentation/).

Улучшения документации всегда приветствуются:

- Исправляйте опечатки и ошибки
- Добавляйте примеры
- Делайте объяснения понятнее
- Переводите материалы

### 3. Отправляйте код

Вносите изменения в код с помощью pull request:

- Исправления ошибок
- Новые функции
- Улучшения производительности
- Тесты

### 4. Предлагайте улучшения (WEP)

Для новой функциональности и изменений публичного поведения используется процесс Wails Enhancement Proposal (WEP). Он обеспечивает прозрачность разработки функций и наличие исполнителя у каждого принятого предложения. Не создавайте issue с запросом функции.

1. При желании сначала предложите идею в категории [Ideas](https://github.com/wailsapp/wails/discussions/categories/ideas) в GitHub Discussions или в [Discord](https://discord.gg/JDdSxwjhGf), чтобы оценить интерес сообщества.
2. Скопируйте [`v3/wep/WEP_TEMPLATE.md`](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/WEP_TEMPLATE.md) в `v3/wep/proposals/<proposal name>/proposal.md` и заполните все разделы.
3. Создайте черновой pull request с заголовком `[WEP] <title>`, содержащий только предложение. PR служит официальным местом для его обсуждения.
4. Соберите отзывы и поддержку — комментарии и реакции «палец вверх» в PR. Выделите на обсуждение не менее двух недель и согласуйте, кто реализует предложение.
5. Отметьте PR как готовый к проверке. Окончательное решение принимают сопровождающие проекта: принятым предложениям присваивается номер WEP, после чего они сливаются в основную ветку.

Полное описание процесса приведено в [`v3/wep/README.md`](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/README.md).

## Начало работы

### Создание форка и клонирование

```bash
# Fork the repository on GitHub
# Then clone your fork
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails

# Add upstream remote
git remote add upstream https://github.com/wailsapp/wails.git
```

### Сборка из исходного кода

```bash
# Build the v3 CLI. Go downloads any required modules automatically.
cd v3
go build -o ../wails3 ./cmd/wails3

# Confirm the built CLI runs and reports its version.
../wails3 version
```

### Запуск тестов

Тесты — часть изменения, а не заключительная поверхностная проверка. Размещайте модульные тесты рядом с проверяемым ими кодом и отдавайте предпочтение табличным тестам, когда одно поведение проверяется с несколькими входными данными или пограничными случаями. Давайте каждому случаю такое имя, чтобы при сбое было понятно, какой сценарий не прошёл проверку.

Новая и изменённая логика должна быть полностью покрыта тестами. Стремитесь к покрытию 100% операторов Go в коде, который добавляет или изменяет ваш PR; общий процент покрытия репозитория не заменяет тестирование самого изменения. Пробел в покрытии может быть оправдан — например, для пути обработки ошибки, возникающей только в определённой ОС, или условия, которое практически невозможно воспроизвести без реального оборудования, — но в описании PR объясните этот пробел и почему его нельзя протестировать разумным способом.

```bash
# Run all v3 tests
cd v3
go test ./...

# Run specific package tests
go test ./pkg/application

# Inspect coverage for the packages you changed
go test ./pkg/application -coverprofile=coverage.out
go tool cover -func=coverage.out
cd ..
```

Интеграционные наборы тестов, обнаружение состояний гонки и полный набор команд, эквивалентный CI, описаны в разделе [«Тестирование и непрерывная интеграция»](/contributing/testing-ci/).

## Внесение изменений

### Создание ветки

```bash
# Update master
git checkout master
git pull upstream master

# Create feature branch
git checkout -b feature/my-feature
```

### Внесение изменений

1. **Пишите код** в соответствии с соглашениями Go
2. **Добавляйте тесты** для новой функциональности
3. При необходимости **обновляйте документацию**
4. **Запускайте тесты**, чтобы убедиться, что ничего не сломалось
5. **Фиксируйте изменения в коммитах** с понятными сообщениями

### Рекомендации по коммитам

```bash
# Good commit messages
git commit -m "fix: resolve window focus issue on macOS"
git commit -m "feat: add support for custom window chrome"
git commit -m "docs: improve bindings documentation"

# Use conventional commits:
# - feat: New feature
# - fix: Bug fix
# - docs: Documentation
# - test: Tests
# - refactor: Code refactoring
# - chore: Maintenance
```

### Отправка pull request

```bash
# Push to your fork
git push origin feature/my-feature

# Open pull request on GitHub
# Provide clear description
# Reference related issues
```

## Рекомендации по pull request

### Хорошее описание PR

```markdown
## Description
Brief description of changes

## Changes
- Added feature X
- Fixed bug Y
- Updated documentation

## Testing
- Tested on macOS 14
- Tested on Windows 11
- All tests passing

## Related Issues
Fixes #123
```

### Контрольный список PR

- [ ] Код соответствует соглашениям Go
- [ ] Тесты добавлены или обновлены
- [ ] Документация обновлена
- [ ] Все тесты проходят
- [ ] Нет несовместимых изменений либо они документированы
- [ ] Сообщения коммитов понятны

## Рекомендации по коду

### Стиль кода Go

```go
// ✅ Good: Clear, documented, tested
// ProcessData processes the input data and returns the result.
// It returns an error if the data is invalid.
func ProcessData(data string) (string, error) {
    if data == "" {
        return "", errors.New("data cannot be empty")
    }
    
    result := process(data)
    return result, nil
}

// ❌ Bad: No docs, no error handling
func ProcessData(data string) string {
    return process(data)
}
```

### Тестирование

```go
func TestProcessData(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid input", "test", "processed", false},
        {"empty input", "", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ProcessData(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ProcessData() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("ProcessData() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Документация

### Создание документации

Полный процесс работы с документацией описан в разделе [«Создание документации с помощью M-Press»](/contributing/documentation/). В нём рассматриваются закреплённая версия, варианты установки исполняемого файла и Go, редактируемые исходные файлы, навигация, проверка переводов и отправка PR для `master`.

Для работы только с документацией не требуется собирать приложение Wails или устанавливать платформенные зависимости. M-Press запускает предварительный просмотр существующего сайта с помощью `mpress dev` из корня репозитория.

### Стиль документации

- Используйте британский вариант английского языка
- Начинайте с описания проблемы
- Приводите работающие примеры
- Добавляйте рекомендации по устранению неполадок
- Добавляйте перекрёстные ссылки на связанные материалы

## Сообщество

### Получение помощи

- **Discord:** [присоединяйтесь к нашему сообществу](https://discord.gg/JDdSxwjhGf)
- **GitHub Discussions:** задавайте вопросы
- **GitHub Issues:** сообщайте об ошибках

### Кодекс поведения

Проявляйте уважение, поддерживайте инклюзивность и ведите себя профессионально. Мы все здесь, чтобы вместе создавать отличное программное обеспечение. Подробности см. в [Кодексе поведения](https://github.com/wailsapp/wails/blob/master/CODE_OF_CONDUCT.md).

## Признание вклада

Участники проекта упоминаются в следующих разделах:

- Примечания к выпуску
- Список участников проекта
- Статистика GitHub

Спасибо за ваш вклад в Wails! 🎉

## Следующие шаги

@cards{cols="2"}
◆ Репозиторий GitHub
Посетите репозиторий Wails.

[Открыть на GitHub →](https://github.com/wailsapp/wails)

---
◆ Сообщество в Discord
Присоединяйтесь к сообществу.

[Присоединиться в Discord →](https://discord.gg/JDdSxwjhGf)

---
📖 Документация
Прочитайте документацию.

[Просмотреть документацию →](/quick-start/why-wails/)

@end
