---
title: "Участие в разработке"
description: "Внесите вклад в Wails"
slug: "contributing"
sourcePath: "contributing/index.md"
---

## Добро пожаловать, участники разработки!

Мы приветствуем вклад в Wails! Исправляете ли вы ошибки, добавляете функции или улучшаете документацию — мы ценим вашу помощь.

## Как внести вклад

### 1. Сообщайте об ошибках

Нашли ошибку? [Создайте задачу](https://github.com/wailsapp/wails/issues/new) и укажите:

- Чёткое описание
- Шаги для воспроизведения
- Ожидаемое и фактическое поведение
- Сведения о системе
- Примеры кода

### 2. Улучшайте документацию

Пул-реквесты с исправлениями приветствуются без предварительно созданной задачи или падающего теста кода.  
Следуйте руководству [«Исправление документации»](/contributing/documentation/), чтобы предварительно просмотреть и проверить изменение с помощью M-Press.

Улучшения документации всегда приветствуются:

- Исправляйте опечатки и ошибки
- Добавляйте примеры
- Делайте пояснения понятнее
- Переводите материалы

### 3. Отправляйте код

Вносите изменения в код с помощью пул-реквестов:

- Исправления ошибок
- Новые функции
- Повышение производительности
- Тесты

### 4. Предлагайте улучшения (WEP)

Новая функциональность и изменения общедоступного поведения проходят процедуру Wails Enhancement Proposal (WEP). Она делает разработку функций прозрачной и гарантирует, что у каждого принятого предложения будет исполнитель. Не создавайте задачу с запросом новой функции.

1. При желании предварительно изложите свою идею в категории [«Идеи»](https://github.com/wailsapp/wails/discussions/categories/ideas) в GitHub Discussions или в [Discord](https://discord.gg/JDdSxwjhGf), чтобы оценить интерес к ней.
2. Скопируйте [`v3/wep/WEP_TEMPLATE.md`](https://github.com/wailsapp/wails/blob/master/v3/wep/WEP_TEMPLATE.md) в `v3/wep/proposals/<proposal name>/proposal.md` и заполните каждый раздел.
3. Создайте черновик пул-реквеста с заголовком `[WEP] <title>`, содержащий только предложение. Этот пул-реквест является официальным местом его обсуждения.
4. Соберите отзывы и поддержку — комментарии и реакции «палец вверх» в пул-реквесте. Отведите на обсуждение не менее двух недель и договоритесь о том, кто реализует предложение.
5. Пометьте пул-реквест как готовый к проверке. Окончательное решение принимают сопровождающие: принятым предложениям присваивают номер WEP, после чего их сливают.

Полная процедура описана в [`v3/wep/README.md`](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md).

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

Тесты — часть изменения, а не заключительная поверхностная проверка. Размещайте модульные тесты рядом с проверяемым ими кодом и отдавайте предпочтение табличным тестам, когда одно поведение проверяется с несколькими входными данными или пограничными случаями. Давайте каждому случаю имя, чтобы при сбое было понятно, какой сценарий не прошёл проверку.

Новую и изменённую логику следует полностью покрывать тестами. Стремитесь обеспечить 100% покрытия операторов Go для кода, который добавляет или изменяет ваш пул-реквест; не рассматривайте процент покрытия всего репозитория как замену тестированию самого изменения. Пробел в покрытии может быть допустим — например, для пути обработки ошибки, возникающей только в определённой ОС, или условия, которое практически невозможно воспроизвести без реального оборудования, — но в описании пул-реквеста объясните этот пробел и почему его нельзя протестировать с разумными затратами.

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

Сведения об интеграционных наборах тестов, обнаружении состояний гонки и полном наборе команд, эквивалентных CI, см. в разделе [«Тестирование и непрерывная интеграция»](/contributing/testing-ci/).

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
4. **Запускайте тесты**, чтобы убедиться, что ничего не нарушено
5. **Фиксируйте изменения** с понятными сообщениями коммитов

### Требования к коммитам

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

### Отправка пул-реквеста

```bash
# Push to your fork
git push origin feature/my-feature

# Open pull request on GitHub
# Provide clear description
# Reference related issues
```

## Требования к пул-реквестам

### Хорошее описание пул-реквеста

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

### Контрольный список пул-реквеста

- [ ] Код соответствует соглашениям Go
- [ ] Тесты добавлены или обновлены
- [ ] Документация обновлена
- [ ] Все тесты проходят
- [ ] Нет нарушающих совместимость изменений либо они задокументированы
- [ ] Сообщения коммитов понятны

## Требования к коду

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

### Написание документации

Для документации используется M-Press. Отредактируйте файлы `.md` в каталоге  
`docs/mpress/content/`, затем из корневого каталога репозитория выполните предварительный просмотр и проверку:

```bash
go install github.com/leaanthony/mpress/cmd/mpress@v1.0.17
mpress dev
mpress build --strict --no-purge-css
mpress check
# Use `python` on Windows if `python3` is unavailable.
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

### Стиль документации

- Используйте орфографию международного варианта английского языка
- Начинайте с описания проблемы
- Приводите рабочие примеры
- Добавляйте рекомендации по устранению неполадок
- Добавляйте перекрёстные ссылки на связанные материалы

## Сообщество

### Получение помощи

- **Discord:** [Присоединиться к нашему сообществу](https://discord.gg/JDdSxwjhGf)
- **Обсуждения на GitHub:** задавайте вопросы
- **Задачи на GitHub:** сообщайте об ошибках

### Кодекс поведения

Будьте уважительны, открыты ко всем и профессиональны. Мы все здесь, чтобы вместе создавать отличное программное обеспечение. Подробности см. в [Кодексе поведения](https://github.com/wailsapp/wails/blob/master/CODE_OF_CONDUCT.md).

## Признание вклада

Участники проекта упоминаются в следующих местах:

- Примечания к выпуску
- Список участников проекта
- Статистика GitHub

Спасибо за ваш вклад в Wails! 🎉

## Дальнейшие действия

@cards{cols="2"}
◆ Репозиторий GitHub
Перейдите в репозиторий Wails.

[Открыть на GitHub →](https://github.com/wailsapp/wails)

---
◆ Сообщество в Discord
Присоединитесь к сообществу.

[Присоединиться в Discord →](https://discord.gg/JDdSxwjhGf)

---
📖 Документация
Ознакомьтесь с документацией.

[Просмотреть документацию →](/quick-start/why-wails/)

@end
