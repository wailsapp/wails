---
title: "Автозапуск"
description: "Настройте запуск приложения при входе пользователя в систему в macOS, Windows и Linux"
slug: "features/autostart/basics"
sourcePath: "features/autostart/basics.md"
---

## Автозапуск

`app.Autostart` регистрирует приложение для автоматического запуска при входе пользователя в систему. Для каждой платформы выбирается подходящий нативный механизм, а пути установки через символические ссылки (Homebrew, Scoop) разрешаются, чтобы регистрация не нарушалась при обновлении исполняемого файла.

Регистрация вступает в силу при **следующем входе в систему**, а не немедленно.

## Быстрый старт

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

Многократный вызов `Enable` безопасен: регистрация каждый раз перезаписывается, поэтому этот метод можно вызывать при каждом запуске, если вы сохраняете выбранную пользователем настройку.

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

Удаляет регистрацию автозапуска. Если приложение не было зарегистрировано, возвращает `nil` — отключение идемпотентно.

```go
func (m *AutostartManager) Disable() error
```

### `IsEnabled`

Сообщает, существует ли регистрация. Выполняется быстро, поскольку не проверяет зарегистрированный путь.

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
| `Path` | `string` | Расположение артефакта регистрации на диске (путь к plist, путь `.desktop` или подраздел реестра). Пусто, если `Enabled` имеет значение false. |
| `Strategy` | `AutostartStrategy` | Механизм, с помощью которого зарегистрировано приложение (см. раздел [«Поведение на разных платформах»](#---)). |

## Поведение на разных платформах

@tabs{sync-key="platform"}
[macOS]
В зависимости от способа упаковки приложения используется один из двух механизмов:

- **macOS 13 и новее, приложение упаковано в `.app`**: `SMAppService.mainAppService`. Работает для приложений в песочнице и сборок для Mac App Store. Запрос TCC на разрешение автоматизации не появляется (исторически использовавшийся подход с AppleScript вызывал такой запрос).
- **macOS до версии 13 или неупакованный исполняемый файл**: plist-файл LaunchAgent записывается в `~/Library/LaunchAgents/<identifier>.plist` с `RunAtLoad=true`.

`Status()` возвращает `AutostartStrategySMAppService` или `AutostartStrategyLaunchAgent`, чтобы вызывающий код мог определить, какой механизм был использован.

Когда приложение обновляется с неупакованной версии до упакованной, `Status()` проверяет оба варианта регистрации, а `Disable()` очищает их, чтобы осиротевший LaunchAgent не продолжал запускать старую сборку.

[Windows]
В разделе реестра `HKCU\Software\Microsoft\Windows\CurrentVersion\Run` добавляется параметр: его именем служит `Identifier` автозапуска, а значением — заключённый в кавычки путь к исполняемому файлу и все аргументы `Arguments`.

Аргументы заключаются в кавычки по правилам `CommandLineToArgvW` (обратные косые черты перед кавычками удваиваются), поэтому пути с пробелами или кавычками корректно восстанавливаются без изменений.

`Status().Strategy` возвращает `AutostartStrategyRegistryRun`.

[Linux]
Запись автозапуска XDG записывается в `$XDG_CONFIG_HOME/autostart/<identifier>.desktop` (по умолчанию — `~/.config/autostart/`) со следующими данными:

```ini
Type=Application
Hidden=false
X-GNOME-Autostart-enabled=true
Exec=<executable> <arguments>
```

Поле `Exec` экранируется согласно [спецификации Desktop Entry сайта freedesktop.org](https://specifications.freedesktop.org/desktop-entry-spec/): перед зарезервированными символами (`"`, `` ` ``, `$`, `\\`) добавляется обратная косая черта, а значение, содержащее пробельные символы, заключается в двойные кавычки.

`Status().Strategy` возвращает `AutostartStrategyXDGAutostart`.

[iOS / Android / сервер]
Не поддерживается. Все методы возвращают `ErrAutostartNotSupported`. Для явного определения отсутствия поддержки используйте `errors.Is(err, application.ErrAutostartNotSupported)`:

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

Если `Options.Identifier` пуст, значение по умолчанию формируется из имени приложения:

| Платформа | Идентификатор по умолчанию |
| --- | --- |
| macOS (упакованное приложение) | Идентификатор пакета приложения, например `com.example.MyApp` |
| macOS (неупакованный исполняемый файл) | `wails.autostart.<slug>`, где `<slug>` формируется из `application.Options.Name` |
| Windows | Slug из `application.Options.Name` (нижний регистр, символы вне `A-Za-z0-9._-` удаляются, пробелы заменяются дефисами) |
| Linux | Тот же slug, что и в Windows |

Идентификаторы должны соответствовать `^[A-Za-z0-9._-]+$` и содержать не более 200 символов. Для macOS рекомендуется формат обратной DNS-записи (он соответствует принятому формату меток launchd).

При переопределении `AutostartOptions.Identifier` этот же идентификатор используется как имя параметра реестра в Windows и имя файла `.desktop` в Linux, поэтому одна строка идентифицирует регистрацию на всех платформах.

## Обнаружение устаревшей регистрации

`Disable()` и `Status()` находят регистрацию, **сопоставляя зарегистрированный путь к исполняемому файлу с `os.Executable()` (после разрешения всех символических ссылок)**, а не выполняя поиск по идентификатору. Это означает следующее:

- **Идентификатор можно безопасно менять между выпусками.** Старую регистрацию по-прежнему можно обнаружить с помощью `Status()` и удалить с помощью `Disable()`, если путь к исполняемому файлу не изменился.
- **Вторая копия приложения, расположенная по другому пути, не перезапишет регистрацию первой.** Каждое расположение исполняемого файла отслеживается независимо.
- **Установки через символические ссылки (Homebrew, Scoop) остаются стабильными.** Перед сопоставлением к `os.Executable()` применяется `filepath.EvalSymlinks`, поэтому при обновлении Homebrew с заменой целевого файла запись не становится недействительной.

Что этот механизм *не* охватывает: если пользователь переместит или переименует исполняемый файл, указав совершенно другой путь, прежняя регистрация станет недействительной (она будет указывать на отсутствующий файл). Приложениям, распространяемым со стабильным путём установки, беспокоиться об этом не нужно; приложения, распространяемые в виде переносимых одиночных исполняемых файлов, должны вызвать `Disable()` перед собственным перемещением либо всегда запускаться через стабильную символическую ссылку.

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

Полный готовый к запуску пример с кнопками проверки состояния, включения и отключения находится в [`examples/autostart/`](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/autostart).
