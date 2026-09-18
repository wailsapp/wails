---
title: "Разрешения"
description: "Управление запросами веб-контента на доступ к камере, микрофону, геолокации и другим возможностям"
slug: "features/windows/permissions"
sourcePath: "features/windows/permissions.md"
---

Веб-контенту, который вызывает `navigator.mediaDevices.getUserMedia()`, API геолокации или API уведомлений, требуется, чтобы основное приложение разрешало или отклоняло эти запросы. Wails предоставляет кроссплатформенную карту `Permissions` в `WebviewWindowOptions`, которая позволяет управлять этим декларативно — без платформозависимого кода.

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

| Константа | Значение | Значение политики |
| --- | --- | --- |
| `PermissionDefault` | 0 | Использовать встроенную обработку платформы (см. ниже) |
| `PermissionAllow` | 1 | Разрешить без запроса |
| `PermissionDeny` | 2 | Отклонить без запроса |

`PermissionDefault` — нулевое значение, поэтому отсутствующие в карте записи обрабатываются по умолчанию.

## Поведение на разных платформах

Каждая платформа обрабатывает `PermissionDefault` по-своему, поскольку лежащие в их основе веб-представления имеют разное встроенное поведение.

### Linux (WebKitGTK)

В WebKitGTK **нет встроенного запроса разрешений**. Если обработчик не подключён, WebKitGTK без уведомления отклоняет каждый запрос — поэтому до добавления этой возможности `getUserMedia` всегда возвращал `NotAllowedError`.

В настоящее время Wails обрабатывает в Linux запросы на доступ к **камере и микрофону**. Геолокация, уведомления и чтение из буфера обмена пока не подключены, поэтому доступ к ним остаётся запрещённым независимо от заданной политики.

| Политика | Камера / микрофон | Геолокация, уведомления, буфер обмена |
| --- | --- | --- |
| `PermissionDefault` | **Разрешено** (восстанавливает работу getUserMedia) | Всегда запрещено |
| `PermissionAllow` | Разрешено | Всегда запрещено (пока не реализовано) |
| `PermissionDeny` | Запрещено | Всегда запрещено |

### Windows (WebView2)

В WebView2 есть встроенный запрос разрешений и API разрешений для отдельных типов. Полностью поддерживаются все пять типов возможностей.

| Политика | Поведение |
| --- | --- |
| `PermissionDefault` | WebView2 показывает встроенный запрос разрешений операционной системы или браузера |
| `PermissionAllow` | Разрешается без уведомления |
| `PermissionDeny` | Отклоняется без уведомления |

**Важно:** до появления этой возможности Wails безусловно вызывал `SetGlobalPermission(Allow)`, без уведомления предоставляя доступ ко всем возможностям. Теперь, если в `Permissions` есть хотя бы одна запись, такое общее разрешение **не** устанавливается. Для возможностей, отсутствующих в карте, вместо автоматического разрешения отображается встроенный запрос WebView2.

Это означает, что если в Windows вы вообще настраиваете `Permissions`, для любой возможности, которую вы явно не указали, будет показан запрос, а не предоставлено разрешение без уведомления. Явно задайте необходимые возможности.

### macOS (TCC)

macOS управляет доступом к камере, микрофону, геолокации и уведомлениям с помощью своей системной инфраструктуры конфиденциальности. Запрос операционной системы появляется автоматически, когда веб-контент впервые запрашивает доступ к возможности, а выбор пользователя запоминается отдельно для каждого приложения в разделе «Системные настройки» → «Конфиденциальность и безопасность».

Это корректно работает без настройки `Permissions`. В настоящее время в macOS карта **игнорируется**: все запросы проходят через TCC независимо от заданных параметров. На практике ограничение заключается в том, что `PermissionDeny` не действует в macOS: нельзя запретить веб-представлению использовать возможность, доступ к которой TCC уже предоставила на системном уровне.

Убедитесь, что `Info.plist` содержит соответствующие ключи с описанием назначения:

```xml
<key>NSMicrophoneUsageDescription</key>
<string>Used for voice input</string>
<key>NSCameraUsageDescription</key>
<string>Used for video calls</string>
```

## Типовые сценарии

### Приложение для захвата мультимедиа

Разрешите доступ к камере и микрофону на всех платформах:

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionAllow,
},
```

В **Linux** это явно разрешает доступ к обоим устройствам; остальные возможности остаются запрещёнными. В **Windows** это разрешает доступ к обоим устройствам; для любой другой возможности, не указанной вами, появится системный запрос. В **macOS** это не действует: всеми разрешениями управляет TCC.

### Запрет захвата мультимедиа в Linux

По умолчанию Linux разрешает доступ к камере и микрофону. Чтобы отказаться от этого:

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionDeny,
    application.PermissionCamera:     application.PermissionDeny,
},
```

### Разрешение всех возможностей в Windows

Чтобы разрешить все возможности без запроса:

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

Поле `Windows.Permissions` для отдельного окна (`map[CoreWebView2PermissionKind]CoreWebView2PermissionState`) по-прежнему работает и позволяет переопределять отдельные возможности после применения кроссплатформенной карты. Используйте его, если требуется доступ к типам разрешений WebView2, для которых нет кроссплатформенных эквивалентов (например, `CoreWebView2PermissionKindOtherSensors`).

```go
Windows: application.WindowsWindow{
    Permissions: map[application.CoreWebView2PermissionKind]application.CoreWebView2PermissionState{
        application.CoreWebView2PermissionKindOtherSensors: application.CoreWebView2PermissionStateAllow,
    },
},
```

В Windows политики применяются в следующем порядке:

1. Кроссплатформенная карта `Permissions` (задаёт состояние для каждого типа через `SetPermission`)
2. Карта `Windows.Permissions` (переопределяет отдельные типы)
3. Для любого типа, не охваченного ни одной из карт: системный запрос WebView2 (если политика настроена) или автоматическое разрешение (если политика не настроена — прежнее поведение)

## Матрица поддержки платформ

| Возможность | Linux | Windows | macOS |
| --- | --- | --- | --- |
| Микрофон | ✅ | ✅ | Только TCC |
| Камера | ✅ | ✅ | Только TCC |
| Геолокация | ❌ пока нет | ✅ | Только TCC |
| Уведомления | ❌ пока нет | ✅ | Только TCC |
| Чтение из буфера обмена | ❌ пока нет | ✅ | Только TCC |

## Устранение неполадок

**`getUserMedia` после обновления по-прежнему не работает в Linux**

Убедитесь, что вы явно не задали `PermissionMicrophone: PermissionDeny` или `PermissionCamera: PermissionDeny`. Значение по умолчанию (не задано) разрешает захват мультимедиа в Linux.

**Windows запрашивает разрешения, которые я не настраивал**

Как только в `Permissions` появляется хотя бы одна запись, Wails перестаёт предоставлять общее разрешение `Allow`. Для возможностей, не указанных вами, появится системный запрос WebView2. Добавьте явные записи `PermissionAllow` для каждой возможности, которую использует ваше приложение.

**Разрешения macOS не работают**

Карта `Permissions` не действует в macOS. Убедитесь, что ваш `Info.plist` содержит правильные ключи с описанием целей использования (`NSMicrophoneUsageDescription`, `NSCameraUsageDescription` и т. д.) и что пользователь предоставил доступ в разделе «Системные настройки» → «Конфиденциальность и безопасность».

**Геолокация, уведомления и буфер обмена не работают в Linux**

В настоящее время в Linux обрабатывается только доступ к камере и микрофону. Поддержка других типов возможностей ещё не реализована — они остаются запрещёнными независимо от заданной политики.
