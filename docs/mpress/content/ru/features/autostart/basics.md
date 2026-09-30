---
title: "Автозапуск"
description: "Регистрация приложения для запуска при входе пользователя в систему на macOS, Windows и Linux"
slug: "features/autostart/basics"
sourcePath: "features/autostart/basics.md"
---

## Автозапуск

`app.Autostart` регистрирует приложение для автоматического запуска при входе пользователя в систему. Для каждой платформы выбирается соответствующий нативный механизм, а пути установки с символическими ссылками (Homebrew, Scoop) разрешаются, чтобы регистрация не нарушалась после обновления исполняемого файла.

Регистрация вступает в силу при **следующем входе в систему**, а не сразу.

## Краткое руководство

```go
import "github.com/wailsapp/wails/v3/pkg/application"

// Register to launch at login
if err := app.Autostart.Enable(); err != nil {
    app.Logger.Error("autostart enable failed", "error", err)
}

// Stop launching at login
if err := app.Autostart.Disable(); err != nil {
    app.Logger.Error("autostart disable failed", "error", err)
}

// Check status
enabled, err := app.Autostart.IsEnabled()
```

## API

### `Enable`

Регистрирует приложение для запуска при входе в систему с параметрами по умолчанию.

```go
func (m *AutostartManager) Enable() error
```

Многократный вызов `Enable` безопасен: регистрация каждый раз перезаписывается, поэтому этот метод можно вызывать при каждом запуске, если сохранена выбранная пользователем настройка.

### `EnableWithOptions`

Регистрирует приложение с пользовательскими параметрами.

```go
func (m *AutostartManager) EnableWithOptions(opts AutostartOptions) error
```

**`AutostartOptions`:**

| Поле | Тип | Описание |
| --- | --- | --- |
| `Identifier` | `string` | Переопределяет автоматически сформированный идентификатор регистрации. См. раздел «Идентификатор» ниже. |
| `Arguments` | `[]string` | Дополнительные аргументы, добавляемые после пути к исполняемому файлу при запуске во время входа в систему (например, `--hidden`). |

### `Disable`

Удаляет регистрацию автозапуска. Возвращает `nil`, если приложение не было зарегистрировано: отключение идемпотентно.

```go
func (m *AutostartManager) Disable() error
```

### `IsEnabled`

Сообщает, существует ли регистрация. Выполняется быстро, поскольку зарегистрированный путь не проверяется.

```go
func (m *AutostartManager) IsEnabled() (bool, error)
```

### `Status`

Возвращает полное состояние регистрации.

```go
func (m *AutostartManager) Status() (AutostartStatus, error)
```

**`AutostartStatus`:**

| Поле | Тип | Описание |
| --- | --- | --- |
| `Enabled` | `bool` | Существует ли регистрация. |
| `Path` | `string` | Расположение артефакта регистрации на диске (путь к plist, путь к `.desktop` или подраздел реестра). Пустое, если `Enabled` имеет значение false. |
| `Strategy` | `AutostartStrategy` | Механизм, с помощью которого зарегистрировано приложение (см. раздел [«Поведение на разных платформах»](#---)). |

## Поведение на разных платформах

@tabs{sync-key="platform"}
[macOS]
В зависимости от способа упаковки приложения используется один из двух механизмов:

- **macOS 13+ и приложение в составе пакета `.app`**: `SMAppService.mainAppService`. Работает для приложений в песочнице и сборок для Mac App Store. Запрос TCC на разрешение автоматизации не появляется (прежний подход с AppleScript вызывал такой запрос).
- **Версии macOS до 13 либо исполняемый файл вне пакета**: plist-файл LaunchAgent записывается в `~/Library/LaunchAgents/<identifier>.plist` с `RunAtLoad=true`.

`Status()` возвращает `AutostartStrategySMAppService` или `AutostartStrategyLaunchAgent`, позволяя вызывающему коду определить, какой механизм был использован.

Когда приложение обновляется с версии вне пакета до версии в составе пакета, `Status()` проверяет оба механизма, а `Disable()` очищает их, чтобы оставшийся LaunchAgent не продолжал запускать старую сборку.

[Windows]
В раздел `HKCU\Software\Microsoft\Windows\CurrentVersion\Run` добавляется параметр реестра: его именем служит `Identifier` автозапуска, а данными — заключённый в кавычки путь к исполняемому файлу и все `Arguments`.

Аргументы заключаются в кавычки по правилам `CommandLineToArgvW` (обратные косые черты перед кавычками удваиваются), поэтому пути с пробелами или кавычками корректно преобразуются в обоих направлениях.

`Status().Strategy` возвращает `AutostartStrategyRegistryRun`.

[Linux]
Запись автозапуска XDG записывается в `$XDG_CONFIG_HOME/autostart/<identifier>.desktop` (по умолчанию — в `~/.config/autostart/`) со следующими данными:

```ini
Type=Application
Hidden=false
X-GNOME-Autostart-enabled=true
Exec=<executable> <arguments>
```

Поле `Exec` экранируется согласно [спецификации Desktop Entry от freedesktop.org](https://specifications.freedesktop.org/desktop-entry-spec/): перед зарезервированными символами (`"`, `` ` ``, `$`, `\\`) добавляется обратная косая черта, а значение с пробельными символами заключается в двойные кавычки.

`Status().Strategy` возвращает `AutostartStrategyXDGAutostart`.

[iOS / Android / сервер]
Не поддерживается. Все методы возвращают `ErrAutostartNotSupported`. Используйте `errors.Is(err, application.ErrAutostartNotSupported)` для корректного определения этой ситуации:

```go
if err := app.Autostart.Enable(); err != nil {
    if errors.Is(err, application.ErrAutostartNotSupported) {
        // hide the toggle in the UI
        return
    }
    // real failure — surface it
}
```

@end

## Идентификатор

Если `Options.Identifier` пуст, идентификатор по умолчанию формируется из имени приложения:

| Платформа | Идентификатор по умолчанию |
| --- | --- |
| macOS (в составе пакета) | Идентификатор пакета приложения, например `com.example.MyApp` |
| macOS (вне пакета) | `wails.autostart.<slug>`, где `<slug>` формируется из `application.Options.Name` |
| Windows | Slug, сформированный из `application.Options.Name` (нижний регистр, символы, не относящиеся к `A-Za-z0-9._-`, удаляются, а пробелы заменяются дефисами) |
| Linux | Тот же slug, что и в Windows |

Идентификаторы должны соответствовать `^[A-Za-z0-9._-]+$` и содержать не более 200 символов. Для macOS рекомендуется формат обратного DNS (он соответствует принятому способу записи меток launchd).

Если значение `AutostartOptions.Identifier` переопределено, один и тот же идентификатор используется как имя параметра реестра в Windows и имя файла `.desktop` в Linux, поэтому одна строка идентифицирует регистрацию на всех платформах.

## Обнаружение устаревшей регистрации

`Disable()` и `Status()` находят регистрацию, **сопоставляя путь зарегистрированного исполняемого файла с `os.Executable()` (после разрешения всех символических ссылок)**, а не выполняя поиск по идентификатору. Это означает следующее:

- **Идентификатор можно безопасно изменять между выпусками.** Старую регистрацию по-прежнему можно обнаружить с помощью `Status()` и удалить с помощью `Disable()` — при условии, что путь к исполняемому файлу не изменился.
- **Вторая копия приложения, расположенная по другому пути, не перезапишет регистрацию первой.** Каждое расположение исполняемого файла отслеживается независимо.
- **Установки через символические ссылки (Homebrew, Scoop) работают стабильно.** Перед сопоставлением к `os.Executable()` применяется `filepath.EvalSymlinks`, поэтому обновление Homebrew, при котором меняется цель ссылки, не оставляет запись потерянной.

Что при этом *не* учитывается: если пользователь переместит или переименует исполняемый файл, указав не связанный с прежним путь, старая регистрация останется бесхозной (она будет указывать на уже отсутствующий файл). Приложениям, распространяемым с неизменным путём установки, об этом беспокоиться не нужно; приложениям, распространяемым как переносимые одиночные исполняемые файлы, следует вызывать `Disable()` перед собственным перемещением либо всегда запускаться через неизменную символическую ссылку.

## Пример

```go
package main

import (
    "errors"

    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "My App",
    })

    // Restore the user's preference on startup
    if userPrefersAutostart() {
        if err := app.Autostart.Enable(); err != nil {
            if !errors.Is(err, application.ErrAutostartNotSupported) {
                app.Logger.Error("autostart", "error", err)
            }
        }
    }

    app.Run()
}
```

Полный готовый к запуску пример с кнопками проверки состояния, включения и отключения находится в [`examples/autostart/`](https://github.com/wailsapp/wails/tree/master/v3/examples/autostart).
