---
title: "Сборки с обфускацией"
description: "Соберите приложение Wails с помощью Garble, чтобы защитить исходный код от обратной разработки"
slug: "guides/build/obfuscation"
sourcePath: "guides/build/obfuscation.md"
---

[Garble](https://github.com/burrowers/garble) — это инструмент сборки Go, который подменяет `go build`, чтобы переименовывать символы, обфусцировать константы и удалять отладочную информацию из полученного исполняемого файла. В Wails v3 полноценная поддержка Garble реализована с помощью двух новых команд.

## Предварительные требования

- **Go 1.26.2 или более поздней версии** — требуется для Garble v0.16.0
- **Garble v0.16.0**

```bash
go install mvdan.cc/garble@v0.16.0
```

@note{type="tip"}
Минимальная версия Go, необходимая для Garble, меняется от выпуска к выпуску. Если вы используете более старую цепочку инструментов Go, перед установкой найдите на [странице выпусков Garble](https://github.com/burrowers/garble/releases) версию, совместимую с вашей цепочкой инструментов.

@end

## Обязательно: добавьте теги JSON к типам сервисов

У каждой структуры, которую возвращает или принимает метод привязанного сервиса, должны быть явно заданы теги JSON для всех экспортируемых полей:

```go
// Without tags — breaks under Garble
type OrderSummary struct {
    ID        int
    Total     float64
    LineItems []LineItem
}

// With tags — safe under Garble
type OrderSummary struct {
    ID        int       `json:"id"`
    Total     float64   `json:"total"`
    LineItems []LineItem `json:"lineItems"`
}
```

@note{type="caution"}
Garble переименовывает экспортируемые поля структур, а Wails передаёт эти структуры в `json.Marshal` через параметр `interface{}`, который Garble не может отследить с помощью статического анализа. Сборка с обфускацией без тегов JSON успешно скомпилируется, но во время выполнения фронтенд получит искажённые или пустые имена полей. Добавьте теги перед запуском сборки с обфускацией.

@end

Собственные типы Wails — `Screen`, `Rect`, `Point`, `Size`, `EnvironmentInfo`, `OSInfo`, `Capabilities` — уже содержат теги. Добавить теги нужно только к собственным типам.

## Сборка с обфускацией

@steps
### Создайте файл стабильных идентификаторов
Выполняйте эту команду каждый раз, когда добавляете, переименовываете или удаляете метод привязанного сервиса:

```bash
wails3 generate bindings -obfuscated
```

Команда создаёт `wails_obfuscated.gen.go` в каталоге вашего пакета main — добавьте этот файл в репозиторий.

### Соберите с помощью Garble
```bash
wails3 build --obfuscated
```

Собирает приложение с использованием обфусцированных привязок.

@end

## Передача дополнительных флагов в Garble

Чтобы передать параметры непосредственно в `garble`, используйте `--garbleargs`:

```bash
# Obfuscate string literals and reduce binary size
wails3 build --obfuscated --garbleargs "-literals -tiny"

# Reproducible output — same seed produces the same binary
wails3 build --obfuscated --garbleargs "-seed=deadbeef"
```

Полный список поддерживаемых флагов приведён в [документации Garble](https://github.com/burrowers/garble#flags).

## Для опытных пользователей: запись файла идентификаторов в другой пакет

По умолчанию `wails_obfuscated.gen.go` записывается рядом с вашим пакетом `main`. Если в вашем проекте сервисы находятся в подпакете, который импортируется `main`, файл можно записать туда с помощью `-obfuscated-output`:

```bash
wails3 generate bindings -obfuscated -obfuscated-output ./internal/services
```

@note{type="caution"}
Пакет назначения должен импортироваться вашим пакетом `main` напрямую или транзитивно, чтобы его `init()` выполнялся при запуске. Если пакет недоступен по цепочке импортов, стабильные идентификаторы не регистрируются и вызовы привязок завершаются ошибкой (например, во время выполнения возникают ошибки `binding not found`).

@end

## Устранение неполадок

### `garble: command not found`

Garble не установлен или `$(go env GOPATH)/bin` отсутствует в переменной `PATH`.

```bash
go install mvdan.cc/garble@v0.16.0
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Фронтенд получает неверные или пустые значения полей

В типах, возвращаемых вашими сервисами, отсутствуют теги `json:"..."`. Проверьте все структуры, возвращаемые привязанными методами, и явно добавьте теги к каждому экспортируемому полю.

### Ошибка `binding not found` в консоли браузера

Файл стабильных идентификаторов отсутствует или не был включён в компиляцию. Проверьте следующее:

- `wails_obfuscated.gen.go` существует в каталоге вашего пакета main (или в каталоге, переданном в `-obfuscated-output`)
- Вы выполнили команду `wails3 build --obfuscated`, которая добавляет тег сборки `wails_obfuscated`
- Если вы использовали `-obfuscated-output`, пакет назначения импортируется `main`

### Windows Defender определяет сборку как вирус

Во время сборки Windows Defender эвристически определяет обфусцированные с помощью Garble исполняемые файлы Go как вредоносные, поскольку в них отсутствуют отладочные символы и они похожи на упакованные исполняемые файлы. Сборка завершается следующей ошибкой:

```
open C:\Users\...\AppData\Local\Temp\go-build...\a.out.exe: The file contains a virus or potentially unwanted software.
```

Добавьте каталог временных файлов, в который Go записывает промежуточные артефакты сборки, и каталог проекта в список исключений Defender:

```powershell
Add-MpPreference -ExclusionPath "$env:TEMP"
Add-MpPreference -ExclusionPath "C:\path\to\your\project"
```

Эти исключения действуют только для указанных путей и не отключают Defender во всей системе.

### Сборка завершается ошибкой `unsupported Go version`

Для Garble v0.16.0 требуется Go 1.26.2 или более поздней версии. Обновите Go либо найдите на [странице выпусков Garble](https://github.com/burrowers/garble/releases) версию, совместимую с вашей цепочкой инструментов.
