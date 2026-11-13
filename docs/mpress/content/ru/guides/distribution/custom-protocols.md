---
title: "Пользовательские URL-протоколы"
description: "Регистрация пользовательских URL-схем для запуска приложения по ссылкам"
slug: "guides/distribution/custom-protocols"
sourcePath: "guides/distribution/custom-protocols.md"
---

Пользовательские URL-протоколы (также называемые URL-схемами) позволяют запускать приложение, когда пользователи переходят по ссылкам с вашим протоколом, например `myapp://action` или `myapp://open/document`.

## Обзор

Пользовательские протоколы обеспечивают:

- **Глубокие ссылки**: запуск приложения с определёнными данными
- **Интеграцию с браузером**: обработку ссылок с веб-страниц
- **Ссылки в электронных письмах**: открытие приложения из почтовых клиентов
- **Взаимодействие между приложениями**: запуск из других приложений

**Пример**: `myapp://open/document?id=123` запускает приложение и открывает документ 123.

## Настройка

Определите пользовательские протоколы в параметрах приложения:

Пользовательские протоколы объявляются в `build/config.yml` (эти данные используются во время упаковки платформенными средствами создания пакетов: макросами NSIS в Windows, манифестом MSIX, `CFBundleURLTypes` в macOS и `.desktop`/`xdg-mime` в Linux). Типа `application.Protocol` и поля `Protocols` в `application.Options` не существует.

```yaml
# build/config.yml
protocols:
  - scheme: myapp
    description: "My Application Protocol"
```

В коде Go отслеживайте запуск с URL через событие `ApplicationLaunchedWithUrl`:

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

func main() {
    app := application.New(application.Options{
        Name:        "My Application",
        Description: "My awesome application",
    })

    app.Event.OnApplicationEvent(events.Common.ApplicationLaunchedWithUrl, func(e *application.ApplicationEvent) {
        handleCustomURL(e.Context().URL())
    })

    app.Run()
}

func handleCustomURL(url string) {
    // Parse and handle the custom URL
    // Example: myapp://open/document?id=123
    println("Received URL:", url)
}
```

## Обработчик протокола

Отслеживайте события протокола для обработки входящих URL:

```go
app.Event.OnApplicationEvent(events.Common.ApplicationLaunchedWithUrl, func(e *application.ApplicationEvent) {
    url := e.Context().URL()

    // Parse the URL
    parsedURL, err := parseCustomURL(url)
    if err != nil {
        app.Logger.Error("Failed to parse URL:", err)
        return
    }

    // Handle different actions
    switch parsedURL.Action {
    case "open":
        openDocument(parsedURL.DocumentID)
    case "settings":
        showSettings()
    case "user":
        showUser Profile(parsedURL.UserID)
    default:
        app.Logger.Warn("Unknown action:", parsedURL.Action)
    }
})
```

## Структура URL

Создавайте понятные иерархические структуры URL:

```
myapp://action/resource?param=value

Examples:
myapp://open/document?id=123
myapp://settings/theme?mode=dark
myapp://user/profile?username=john
```

**Рекомендации:**

- Используйте имена схем в нижнем регистре
- Выбирайте короткие и запоминающиеся схемы
- Используйте иерархические пути к ресурсам
- Передавайте необязательные данные в параметрах запроса
- Кодируйте специальные символы в формате URL

## Регистрация на платформах

Пользовательские протоколы регистрируются по-разному на каждой платформе.

@tabs{sync-key="platform"}
[Windows]
### Установщик NSIS для Windows

**При использовании установщиков NSIS Wails v3 автоматически регистрирует пользовательские протоколы**.

#### Автоматическая регистрация

Когда вы собираете приложение с помощью `wails3 build`, установщик NSIS:

1. Автоматически регистрирует все протоколы, объявленные в `build/config.yml`, в разделе `protocols:`
2. Связывает протоколы с исполняемым файлом приложения
3. Создаёт необходимые записи реестра
4. Удаляет связи с протоколами при удалении приложения

**Дополнительная настройка не требуется!**

#### Как это работает

Шаблон NSIS содержит встроенные макросы:

- `wails.associateCustomProtocols` — регистрирует протоколы при установке
- `wails.unassociateCustomProtocols` — удаляет протоколы при удалении приложения

Эти макросы вызываются автоматически в соответствии с вашей конфигурацией `Protocols`.

#### Ручная настройка реестра (для опытных пользователей)

Если вам требуется ручная регистрация (без NSIS):

```batch
@echo off
REM Register custom protocol
REG ADD "HKEY_CURRENT_USER\SOFTWARE\Classes\myapp" /ve /d "URL:My Application Protocol" /f
REG ADD "HKEY_CURRENT_USER\SOFTWARE\Classes\myapp" /v "URL Protocol" /t REG_SZ /d "" /f
REG ADD "HKEY_CURRENT_USER\SOFTWARE\Classes\myapp\shell\open\command" /ve /d "\"%1\"" /f
```

#### Тестирование

Проверьте регистрацию протокола:

```powershell
# Open protocol URL from PowerShell
Start-Process "myapp://test/action"

# Or from command prompt
start myapp://test/action
```

### Пакет MSIX для Windows

При использовании пакетов MSIX пользовательские протоколы также регистрируются автоматически.

#### Автоматическая регистрация

Когда вы собираете приложение в формате MSIX, в манифест автоматически добавляются регистрации протоколов из конфигурации протоколов `build/config.yml`.

Созданный манифест содержит:

```xml
<uap:Extension Category="windows.protocol">
  <uap:Protocol Name="myapp">
    <uap:DisplayName>My Application Protocol</uap:DisplayName>
  </uap:Protocol>
</uap:Extension>
```

#### Универсальные ссылки (ссылки из веб-сайта в приложение)

Windows поддерживает **ссылки из веб-сайта в приложение**, которые работают аналогично универсальным ссылкам в macOS. При развёртывании приложения в виде пакета MSIX можно разрешить запуск приложения напрямую по HTTPS-ссылкам.

@note{type="note"}
Для ссылок из веб-сайта в приложение необходимо вручную настроить манифест. Пользовательские схемы протоколов настраиваются автоматически из `build/config.yml`, но связанные домены необходимо добавить в манифест MSIX вручную.

@end

Чтобы включить ссылки из веб-сайта в приложение, следуйте [руководству Microsoft по таким ссылкам](https://learn.microsoft.com/en-us/windows/apps/develop/launch/web-to-app-linking). Вам потребуется:

1. **Вручную добавить App URI Handler в манифест MSIX** (`build/windows/msix/app_manifest.xml`):
  ```xml
  <uap3:Extension Category="windows.appUriHandler">
    <uap3:AppUriHandler>
      <uap3:Host Name="myawesomeapp.com"/>
    </uap3:AppUriHandler>
  </uap3:Extension>
  ```


2. **Настроить `windows-app-web-link` на своём веб-сайте:** разместите файл `windows-app-web-link` по адресу `https://myawesomeapp.com/.well-known/windows-app-web-link`. Этот файл должен содержать сведения о пакете приложения и обрабатываемых им путях.

Когда ссылка из веб-сайта в приложение запускает ваше приложение, вы получаете то же событие `ApplicationLaunchedWithUrl`, что и при использовании пользовательских схем протоколов.

[macOS]
### Настройка Info.plist

В macOS протоколы регистрируются через файл `Info.plist`.

#### Автоматическая настройка

При сборке с помощью `wails3 build` Wails автоматически создаёт `Info.plist` с вашими протоколами.

Протоколы, объявленные в `build/config.yml`, добавляются в:

```xml
<key>CFBundleURLTypes</key>
<array>
    <dict>
        <key>CFBundleURLName</key>
        <string>My Application Protocol</string>
        <key>CFBundleURLSchemes</key>
        <array>
            <string>myapp</string>
        </array>
        <key>CFBundleTypeRole</key>
        <string>Editor</string>
    </dict>
</array>
```

#### Тестирование

```bash
# Open protocol URL from terminal
open "myapp://test/action"

# Check registered handlers
/System/Library/Frameworks/CoreServices.framework/Versions/A/Frameworks/LaunchServices.framework/Versions/A/Support/lsregister -dump | grep myapp
```

#### Универсальные ссылки

Помимо пользовательских схем протоколов, macOS также поддерживает **универсальные ссылки**, которые позволяют запускать приложение по обычным HTTPS-ссылкам (например, `https://myawesomeapp.com/path`). Универсальные ссылки обеспечивают удобный и целостный переход между веб-приложением и настольным приложением.

@note{type="caution"}
Чтобы универсальные ссылки работали, приложение для macOS должно быть **подписано** действительным сертификатом Apple Developer с использованием профиля подготовки. Сборки без подписи или со специальной подписью не смогут открывать универсальные ссылки. Перед тестированием убедитесь, что приложение правильно подписано.

@end

Чтобы включить универсальные ссылки, следуйте [руководству Apple по поддержке универсальных ссылок в приложении](https://developer.apple.com/documentation/xcode/supporting-universal-links-in-your-app). Вам потребуется:

1. **Добавить разрешения** в `entitlements.plist`:
  ```xml
  <key>com.apple.developer.associated-domains</key>
  <array>
    <string>applinks:myawesomeapp.com</string>
  </array>
  ```


2. **Добавить NSUserActivityTypes в Info.plist**:
  ```xml
  <key>NSUserActivityTypes</key>
  <array>
    <string>NSUserActivityTypeBrowsingWeb</string>
  </array>
  ```


3. **Настроить `apple-app-site-association` на своём веб-сайте:** разместите файл `apple-app-site-association` по адресу `https://myawesomeapp.com/.well-known/apple-app-site-association`.

Когда универсальная ссылка запускает приложение, вы получаете то же событие `ApplicationLaunchedWithUrl`, поэтому код обработки остаётся таким же, как для пользовательских схем протоколов.

[Linux]
### Файл desktop entry

В Linux протоколы регистрируются с помощью файлов `.desktop`.

#### Автоматическая настройка

При сборке с помощью `wails3 build` Wails создаёт файл desktop entry с обработчиками протоколов.

**Исправлено в v3**: шаблон рабочего стола Linux теперь корректно поддерживает обработку протоколов.

Созданный файл рабочего стола содержит:

```ini
[Desktop Entry]
Type=Application
Name=My Application
Exec=/usr/bin/myapp %u
MimeType=x-scheme-handler/myapp;
```

#### Регистрация вручную

При необходимости установите файл рабочего стола вручную:

```bash
# Copy desktop file
cp myapp.desktop ~/.local/share/applications/

# Update desktop database
update-desktop-database ~/.local/share/applications/

# Register protocol handler
xdg-mime default myapp.desktop x-scheme-handler/myapp
```

#### Тестирование

```bash
# Open protocol URL
xdg-open "myapp://test/action"

# Check registered handler
xdg-mime query default x-scheme-handler/myapp
```

@end

## Полный пример

Ниже приведён полный пример обработки нескольких действий протокола:

```go
package main

import (
    "fmt"
    "net/url"
    "strings"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

type App struct {
    app    *application.App
    window *application.WebviewWindow
}

func main() {
    // Protocol registration lives in build/config.yml (the platform packagers
    // consume it); the application code just listens for the launch event.
    app := application.New(application.Options{
        Name:        "DeepLink Demo",
        Description: "Custom protocol demonstration",
    })

    myApp := &App{app: app}
    myApp.setup()

    app.Run()
}

func (a *App) setup() {
    // Create window
    a.window = a.app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:  "DeepLink Demo",
        Width:  800,
        Height: 600,
        URL:    "http://wails.localhost/",
    })

    // Handle custom protocol URLs
    a.app.Event.OnApplicationEvent(events.Common.ApplicationLaunchedWithUrl, func(e *application.ApplicationEvent) {
        a.handleDeepLink(e.Context().URL())
    })
}

func (a *App) handleDeepLink(rawURL string) {
    // Parse URL
    parsedURL, err := url.Parse(rawURL)
    if err != nil {
        a.app.Logger.Error("Failed to parse URL:", err)
        return
    }

    // Bring window to front
    a.window.Show()
    a.window.Focus()

    // Extract path and query
    path := strings.Trim(parsedURL.Path, "/")
    query := parsedURL.Query()

    // Handle different actions
    parts := strings.Split(path, "/")
    if len(parts) == 0 {
        return
    }

    action := parts[0]

    switch action {
    case "open":
        if len(parts) >= 2 {
            resource := parts[1]
            id := query.Get("id")
            a.openResource(resource, id)
        }

    case "settings":
        section := ""
        if len(parts) >= 2 {
            section = parts[1]
        }
        a.openSettings(section)

    case "user":
        if len(parts) >= 2 {
            username := parts[1]
            a.openUserProfile(username)
        }

    default:
        a.app.Logger.Warn("Unknown action:", action)
    }
}

func (a *App) openResource(resourceType, id string) {
    fmt.Printf("Opening %s with ID: %s\n", resourceType, id)
    // Emit event to frontend
    a.app.Event.Emit("navigate", map[string]string{
        "type": resourceType,
        "id":   id,
    })
}

func (a *App) openSettings(section string) {
    fmt.Printf("Opening settings section: %s\n", section)
    a.app.Event.Emit("navigate", map[string]string{
        "page":    "settings",
        "section": section,
    })
}

func (a *App) openUserProfile(username string) {
    fmt.Printf("Opening user profile: %s\n", username)
    a.app.Event.Emit("navigate", map[string]string{
        "page": "user",
        "user": username,
    })
}
```

## Интеграция с фронтендом

Обрабатывайте события навигации во фронтенде:

```javascript
import { Events } from '@wailsio/runtime'

// Listen for navigation events from protocol handler
Events.On('navigate', (event) => {
    const { type, id, page, section, user } = event.data

    if (type === 'document') {
        // Open document with ID
        router.push(`/document/${id}`)
    } else if (page === 'settings') {
        // Open settings
        router.push(`/settings/${section}`)
    } else if (page === 'user') {
        // Open user profile
        router.push(`/user/${user}`)
    }
})
```

## Рекомендации по безопасности

### Проверяйте все входные данные

Всегда проверяйте и очищайте URL-адреса из внешних источников:

```go
func (a *App) handleDeepLink(rawURL string) {
    // Parse URL
    parsedURL, err := url.Parse(rawURL)
    if err != nil {
        a.app.Logger.Error("Invalid URL:", err)
        return
    }

    // Validate scheme
    if parsedURL.Scheme != "myapp" {
        a.app.Logger.Warn("Invalid scheme:", parsedURL.Scheme)
        return
    }

    // Validate path
    path := strings.Trim(parsedURL.Path, "/")
    if !isValidPath(path) {
        a.app.Logger.Warn("Invalid path:", path)
        return
    }

    // Sanitize parameters
    params := sanitizeQueryParams(parsedURL.Query())

    // Process validated URL
    a.processDeepLink(path, params)
}

func isValidPath(path string) bool {
    // Only allow alphanumeric and forward slashes
    validPath := regexp.MustCompile(`^[a-zA-Z0-9/]+$`)
    return validPath.MatchString(path)
}

func sanitizeQueryParams(query url.Values) map[string]string {
    sanitized := make(map[string]string)
    for key, values := range query {
        if len(values) > 0 {
            // Take first value and sanitize
            sanitized[key] = sanitizeString(values[0])
        }
    }
    return sanitized
}
```

### Предотвращайте атаки с внедрением кода

Никогда не выполняйте URL-адреса напрямую как код или SQL-запросы:

```go
// ❌ DON'T: Execute URL content
func badHandler(url string) {
    exec.Command("sh", "-c", url).Run() // DANGEROUS!
}

// ✅ DO: Parse and validate
func goodHandler(url string) {
    parsed, _ := url.Parse(url)
    action := parsed.Query().Get("action")

    // Whitelist allowed actions
    allowed := map[string]bool{
        "open":     true,
        "settings": true,
        "help":     true,
    }

    if allowed[action] {
        handleAction(action)
    }
}
```

## Тестирование

### Ручное тестирование

Тестируйте обработчики протоколов во время разработки:

**Windows:**

```powershell
Start-Process "myapp://test/action?id=123"
```

**macOS:**

```bash
open "myapp://test/action?id=123"
```

**Linux:**

```bash
xdg-open "myapp://test/action?id=123"
```

### Тестирование с помощью HTML

Создайте тестовую HTML-страницу:

```html
<!DOCTYPE html>
<html>
<head>
    <title>Protocol Test</title>
</head>
<body>
    <h1>Custom Protocol Test Links</h1>

    <ul>
        <li><a href="myapp://open/document?id=123">Open Document 123</a></li>
        <li><a href="myapp://settings/theme?mode=dark">Dark Mode Settings</a></li>
        <li><a href="myapp://user/profile?username=john">User Profile</a></li>
    </ul>
</body>
</html>
```

## Устранение неполадок

### Протокол не зарегистрирован

**Windows:**

- Проверьте реестр: `HKEY_CURRENT_USER\SOFTWARE\Classes\<scheme>`
- Переустановите приложение с помощью установщика NSIS
- Убедитесь, что установщик был запущен с необходимыми разрешениями

**macOS:**

- Пересоберите приложение с помощью `wails3 build`
- Проверьте `Info.plist` в пакете приложения: `MyApp.app/Contents/Info.plist`
- Сбросьте настройки Launch Services: `/System/Library/Frameworks/CoreServices.framework/Frameworks/LaunchServices.framework/Support/lsregister -kill`

**Linux:**

- Проверьте файл рабочего стола: `~/.local/share/applications/myapp.desktop`
- Обновите базу данных: `update-desktop-database ~/.local/share/applications/`
- Проверьте обработчик: `xdg-mime query default x-scheme-handler/myapp`

### Приложение не запускается

**Проверьте журналы:**

```go
app := application.New(application.Options{
    LogLevel: slog.LevelDebug, // requires `import "log/slog"`
    // ...
})
```

**Распространённые проблемы:**

- Приложение установлено не в ожидаемую папку
- Путь к исполняемому файлу в регистрационных данных не соответствует его фактическому расположению
- Проблемы с разрешениями

## Рекомендации

### ✅ Рекомендуется

- **Используйте понятные имена схем** — `mycompany-myapp` вместо `mca`
- **Проверяйте все входные данные** — никогда не доверяйте URL-адресам из внешних источников
- **Корректно обрабатывайте ошибки** — записывайте недопустимые URL-адреса в журнал и не допускайте аварийного завершения приложения
- **Информируйте пользователя** — показывайте, какое действие было запущено
- **Тестируйте на всех платформах** — обработка протоколов различается
- **Документируйте структуру URL-адресов** — помогите пользователям и специалистам по интеграции

### ❌ Не рекомендуется

- **Не используйте распространённые имена схем** — избегайте `http`, `file`, `app` и подобных имён
- **Не выполняйте URL-адреса как код** — это создаёт огромный риск для безопасности
- **Не открывайте доступ к критически важным операциям** — запрашивайте подтверждение для деструктивных действий
- **Не полагайтесь на то, что протоколы работают везде** — предусмотрите резервные механизмы
- **Не забывайте о кодировании URL-адресов** — правильно обрабатывайте специальные символы

## Дальнейшие шаги

- [Упаковка для Windows](/guides/build/windows/) — узнайте о параметрах установщика NSIS
- [Ассоциации файлов](/guides/file-associations/) — открывайте файлы в своём приложении
- [Один экземпляр](/guides/single-instance/) — предотвращайте запуск нескольких экземпляров приложения

---

**Остались вопросы?** Задайте их в [Discord](https://discord.gg/JDdSxwjhGf) или ознакомьтесь с [примерами](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples).
