---
title: "Разрешения"
description: "Управление запросами веб-контента на доступ к камере, микрофону, геолокации и другим возможностям"
slug: "features/windows/permissions"
sourcePath: "features/windows/permissions.md"
---

Веб-контенту, который вызывает `navigator.mediaDevices.getUserMedia()`, API геолокации или API уведомлений, необходимо, чтобы хост-приложение разрешало или отклоняло эти запросы. Wails предоставляет кроссплатформенную карту `Permissions` в `WebviewWindowOptions`, которая позволяет управлять этим декларативно — без платформозависимого кода.

## Быстрый старт

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "My App",
    Permissions: map[application.PermissionType]application.Permission{
        application.PermissionMicrophone: application.PermissionAllow,
        application.PermissionCamera:     application.PermissionAllow,
    },
})
```

Запросы веб-контента этого окна на доступ к камере и микрофону разрешаются без запроса браузера.

## Типы разрешений

`PermissionType` (uint8) определяет возможность, доступ к которой может запросить веб-контент.

| Константа | Возможность |
| --- | --- |
| `PermissionMicrophone` | `getUserMedia({audio: true})` |
| `PermissionCamera` | `getUserMedia({video: true})` |
| `PermissionGeolocation` | `navigator.geolocation` |
| `PermissionNotifications` | `Notification.requestPermission()` |
| `PermissionClipboardRead` | `navigator.clipboard.readText()` |

## Значения разрешений

`Permission` (uint8) — это политика, применяемая к заданному типу.

| Константа | Значение | Назначение |
| --- | --- | --- |
| `PermissionDefault` | 0 | Использовать стандартную обработку платформы (см. ниже) |
| `PermissionAllow` | 1 | Разрешить без запроса |
| `PermissionDeny` | 2 | Отклонить без запроса |

`PermissionDefault` является нулевым значением, поэтому отсутствующие в карте записи обрабатываются согласно политике по умолчанию.

## Поведение на разных платформах

Каждая платформа обрабатывает `PermissionDefault` по-разному, поскольку лежащие в их основе веб-представления имеют разное стандартное поведение.

### Linux (WebKitGTK)

В WebKitGTK **нет собственного запроса разрешений**. Если обработчик не подключён, WebKitGTK без уведомления отклоняет каждый запрос — поэтому до добавления этой возможности `getUserMedia` всегда возвращал `NotAllowedError`.

Сейчас в Linux Wails обрабатывает запросы на доступ к **камере и микрофону**. Геолокация, уведомления и чтение из буфера обмена пока не подключены и остаются запрещёнными независимо от заданной политики.

| Политика | Камера / микрофон | Геолокация, уведомления, буфер обмена |
| --- | --- | --- |
| `PermissionDefault` | **Разрешено** (восстанавливает работу getUserMedia) | Всегда запрещено |
| `PermissionAllow` | Разрешено | Всегда запрещено (пока не реализовано) |
| `PermissionDeny` | Запрещено | Всегда запрещено |

### Windows (WebView2)

В WebView2 есть собственный запрос разрешений и API разрешений для каждого типа. Полностью поддерживаются все пять типов возможностей.

| Политика | Поведение |
| --- | --- |
| `PermissionDefault` | WebView2 показывает собственный запрос разрешения операционной системы или браузера |
| `PermissionAllow` | Разрешено без уведомления |
| `PermissionDeny` | Отклонено без уведомления |

**Важно:** до появления этой возможности Wails безусловно вызывал `SetGlobalPermission(Allow)`, тем самым без уведомления разрешая все возможности. Теперь, если в `Permissions` присутствует хотя бы одна запись, такое разрешение для всех возможностей **не** задаётся. Для возможностей, не указанных в карте, вместо автоматического разрешения отображается стандартный запрос WebView2.

Это означает, что если в Windows вы каким-либо образом настроите `Permissions`, то для каждой возможности, которую вы явно не указали, вместо разрешения без уведомления будет показан запрос. Явно задайте необходимые возможности.

### macOS (WKWebView + TCC)

В macOS разрешение необходимо получить на двух уровнях. Перед началом сеанса захвата WKWebView обращается к приложению, и именно на этот запрос отвечает карта `Permissions`. На нижележащем уровне доступ непосредственно к устройству контролирует системная инфраструктура конфиденциальности (TCC): запрос операционной системы появляется, когда приложение впервые фактически обращается к камере или микрофону, а выбор пользователя запоминается отдельно для каждого приложения в разделе «Системные настройки» → «Конфиденциальность и безопасность».

Wails обрабатывает в macOS 12 и более поздних версиях запросы на доступ к **камере и микрофону**. Для геолокации, уведомлений и чтения из буфера обмена нет эквивалента `WKUIDelegate`, поэтому их обработка не подключена и полностью возложена на TCC. Как и в Linux, заданная для них политика не действует.

| Политика | Камера / микрофон | Геолокация, уведомления, буфер обмена |
| --- | --- | --- |
| `PermissionDefault` | WebKit показывает собственный запрос разрешения | Только TCC |
| `PermissionAllow` | Запрос WebKit пропускается — **TCC продолжает действовать** | Только TCC |
| `PermissionDeny` | Отклонено WebKit; системная авторизация выполняется отдельно | Только TCC |

`PermissionAllow` разрешает запрос веб-представления, а не доступ к устройству. Системная авторизация по-прежнему применяется независимо, и приложение, которому отказано в «Системных настройках», не получает доступ. Пропускается только собственный запрос WebKit.

В версиях ниже macOS 12 метод делегата отсутствует, поэтому карта там игнорируется, а каждый запрос обрабатывается стандартным запросом WebKit.

Убедитесь, что `Info.plist` содержит соответствующие ключи с описаниями причин использования:

```xml
<key>NSMicrophoneUsageDescription</key>
<string>Used for voice input</string>
<key>NSCameraUsageDescription</key>
<string>Used for video calls</string>
```

@note{type="caution" title="Описания причин использования обязательны"}

WebKit может запросить системную авторизацию до вызова делегата разрешений Wails. `PermissionDeny` блокирует запрос WebKit, но не гарантирует отсутствия запроса TCC или запроса авторизации AVFoundation.

Добавьте соответствующие описания причин использования для устройств, к которым обращается приложение, а при работе в песочнице — также соответствующие права доступа к устройствам. Не используйте `PermissionDeny` вместо этих деклараций.

@end

## Распространённые сценарии

### Приложение для захвата мультимедиа

Разрешите доступ к камере и микрофону на всех платформах:

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionAllow,
},
```

В **Linux** это явно разрешает доступ к обоим устройствам; остальные возможности остаются запрещёнными. В **Windows** это разрешает доступ к обоим устройствам; для любой другой возможности, которую вы не укажете, появится системный запрос. В **macOS** это разрешает доступ к обоим устройствам на уровне WebKit, поэтому запрос браузера не появляется; при первом использовании TCC всё равно запрашивает доступ непосредственно к устройствам, а в `Info.plist` должны присутствовать оба ключа с описанием цели использования.

### Запрет захвата мультимедиа в Linux

По умолчанию Linux разрешает доступ к камере и микрофону. Чтобы отказаться от этого поведения:

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionDeny,
    application.PermissionCamera:     application.PermissionDeny,
},
```

### Разрешение всех возможностей в Windows

Чтобы предоставить доступ ко всем возможностям без запросов:

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone:    application.PermissionAllow,
    application.PermissionCamera:        application.PermissionAllow,
    application.PermissionGeolocation:   application.PermissionAllow,
    application.PermissionNotifications: application.PermissionAllow,
    application.PermissionClipboardRead: application.PermissionAllow,
},
```

### Политики для отдельных окон

Для разных окон можно задать разные политики:

```go
// Main app window — full media access
mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "App",
    Permissions: map[application.PermissionType]application.Permission{
        application.PermissionMicrophone: application.PermissionAllow,
        application.PermissionCamera:     application.PermissionAllow,
    },
})

// Settings window — no special capabilities
settingsWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Settings",
    // No Permissions entry — uses platform defaults
})

// Embedded content window — deny media capture
embeddedWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Embedded",
    Permissions: map[application.PermissionType]application.Permission{
        application.PermissionMicrophone: application.PermissionDeny,
        application.PermissionCamera:     application.PermissionDeny,
    },
})
```

## Переопределение для Windows

Поле `Windows.Permissions` для отдельного окна (`map[CoreWebView2PermissionKind]CoreWebView2PermissionState`) по-прежнему работает и позволяет переопределять отдельные возможности после применения кроссплатформенной карты. Используйте его, когда требуется доступ к типам разрешений WebView2, у которых нет кроссплатформенного эквивалента (например, `CoreWebView2PermissionKindOtherSensors`).

```go
Windows: application.WindowsWindow{
    Permissions: map[application.CoreWebView2PermissionKind]application.CoreWebView2PermissionState{
        application.CoreWebView2PermissionKindOtherSensors: application.CoreWebView2PermissionStateAllow,
    },
},
```

В Windows используется следующий порядок обработки:

1. Кроссплатформенная карта `Permissions` (задаёт состояние для каждого типа через `SetPermission`)
2. Карта `Windows.Permissions` (переопределяет отдельные типы)
3. Для любого типа, не охваченного ни одной из карт: системный запрос WebView2 (если политика настроена) или автоматическое разрешение (если политика не настроена — прежнее поведение)

## Матрица поддержки платформ

| Возможность | Linux | Windows | macOS |
| --- | --- | --- | --- |
| Микрофон | ✅ | ✅ | ✅ (macOS 12 и новее) |
| Камера | ✅ | ✅ | ✅ (macOS 12 и новее) |
| Геолокация | ❌ пока нет | ✅ | ❌ пока нет |
| Уведомления | ❌ пока нет | ✅ | ❌ пока нет |
| Чтение буфера обмена | ❌ пока нет | ✅ | ❌ пока нет |

Если для macOS указано ✅, политика отвечает на запрос WebKit, но поверх неё доступ к устройству по-прежнему контролирует TCC. Если указано ❌, возможность контролируется только TCC.

## Устранение неполадок

**`getUserMedia` по-прежнему не работает в Linux после обновления**

Убедитесь, что вы не задали явно `PermissionMicrophone: PermissionDeny` или `PermissionCamera: PermissionDeny`. По умолчанию, когда эти параметры не заданы, захват мультимедиа в Linux разрешён.

**Windows запрашивает разрешения, которые я не настраивал**

Как только в `Permissions` появляется хотя бы одна запись, Wails больше не предоставляет общее разрешение `Allow`. Для возможностей, которые вы не укажете, появится системный запрос WebView2. Добавьте явные записи `PermissionAllow` для каждой возможности, используемой вашим приложением.

**Разрешения macOS не работают**

`Permissions` распространяется на камеру и микрофон в macOS 12 и новее; геолокация, уведомления и чтение буфера обмена пока не подключены и игнорируют политику. Там, где политика учитывается, поверх неё доступ к устройству по-прежнему контролирует TCC: `PermissionAllow` устраняет запрос WebKit, но не системный запрос. Убедитесь, что ваш `Info.plist` содержит правильные ключи с описанием цели использования (`NSMicrophoneUsageDescription`, `NSCameraUsageDescription`), а пользователь предоставил доступ в разделе «Системные настройки» → «Конфиденциальность и безопасность».

**Моё приложение macOS завершает работу, когда веб-контент запрашивает доступ к камере или микрофону**

Убедитесь, что `Info.plist` содержит `NSCameraUsageDescription` или `NSMicrophoneUsageDescription` для каждого устройства, доступ к которому запрашивает ваше приложение. Отсутствие описаний может привести к отклонению запроса; обращение к AVFoundation без обязательного описания может привести к завершению работы приложения. `PermissionDeny` управляет решением WebKit и не заменяет эти ключи.

**Геолокация, уведомления и буфер обмена не работают в Linux**

В настоящее время в Linux обрабатываются только камера и микрофон. Поддержка остальных типов возможностей пока не реализована — они остаются запрещёнными независимо от заданной политики.
