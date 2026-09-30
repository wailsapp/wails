---
title: "Уведомления"
description: "Отображение нативных системных уведомлений с кнопками действий и полем ввода текста"
slug: "features/notifications/overview"
sourcePath: "features/notifications/overview.md"
---

## Введение

Wails предоставляет комплексную кроссплатформенную систему уведомлений для настольных приложений. Этот сервис позволяет отображать нативные системные уведомления и поддерживает:

- Простые уведомления с заголовком, подзаголовком и основным текстом
- Интерактивные уведомления с кнопками действий и текстовыми ответами
- Многократно используемые [категории уведомлений](#--4) для действий
- Пользовательские [звуки](#--7) (стандартные, отключённые или заданные по имени)
- [Вложения](#heading-1) (изображения на всех платформах; аудио и видео в macOS)
- [Группировку связанных уведомлений](#---2) по `ThreadID`
- [Приоритет](#--8) с помощью `InterruptionLevel` (`passive` / `active` / `timeSensitive` / `critical`)
- [Доставку по расписанию](#--9) (нативными средствами в macOS; с помощью внутрипроцессного таймера в Windows и Linux)
- [Обновление уже отправленного уведомления](#--10) по идентификатору

Если платформа не поддерживает новое необязательное поле, связанная с ним функциональность корректно упрощается; матрицу поддержки отдельных возможностей см. в разделе [Особенности платформ](#--11).

## Основы использования

### Создание сервиса

Сначала инициализируйте сервис уведомлений:

```go
import "github.com/wailsapp/wails/v3/pkg/application"
import "github.com/wailsapp/wails/v3/pkg/services/notifications"

// Create a new notification service
notifier := notifications.New()

//Register the service with the application
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(notifier),
    },
})
```

## Разрешение на уведомления

Для показа уведомлений в macOS требуется разрешение пользователя. Запросите разрешение и проверьте его наличие:

```go
authorized, err := notifier.CheckNotificationAuthorization()
if err != nil {
    // Handle authorization error
}
if authorized {
    // Send notifications
} else {
    // Request authorization
    authorized, err = notifier.RequestNotificationAuthorization()
}
```

В Windows и Linux этот вызов всегда возвращает `true`.

## Типы уведомлений

### Простые уведомления

Отправьте пользователям простое уведомление с уникальным идентификатором, заголовком, необязательным подзаголовком (в macOS и Linux) и основным текстом:

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID: "unique-id",
    Title: "New Calendar Invite",
    Subtitle: "From: Jane Doe", // Optional
    Body: "Tap to view the event",
})

```

### Интерактивные уведомления

Отправьте уведомление с кнопками действий и полями ввода текста. Для таких уведомлений необходимо сначала зарегистрировать категорию уведомлений:

```go
// Define a unique category id
categoryID := "unique-category-id"

// Define a category with actions
category := notifications.NotificationCategory{
    ID: categoryID,
    Actions: []notifications.NotificationAction{
        {
            ID:    "OPEN", 
            Title: "Open",
        },
        {
            ID:          "ARCHIVE", 
            Title:       "Archive", 
            Destructive: true,  /* macOS-specific */
        },
    },
    HasReplyField:    true,
    ReplyPlaceholder: "message...",
    ReplyButtonTitle: "Reply",
}

// Register the category
notifier.RegisterNotificationCategory(category)

// Send an interactive notification with the actions registered in the provided category
notifier.SendNotificationWithActions(notifications.NotificationOptions{
    ID:         "unique-id",
    Title:      "New Message",
    Subtitle:   "From: Jane Doe",
    Body:       "Are you able to make it?",
    CategoryID: categoryID,
})
```

## Ответы на уведомления

Обработайте взаимодействия пользователя с уведомлениями:

```go
notifier.OnNotificationResponse(func(result notifications.NotificationResult) {
    response := result.Response
    fmt.Printf("Notification %s was actioned with: %s\n", response.ID, response.ActionIdentifier)

    if response.ActionIdentifier == "TEXT_REPLY" {
        fmt.Printf("User replied: %s\n", response.UserText)
    }

    if data, ok := response.UserInfo["sender"].(string); ok {
        fmt.Printf("Original sender: %s\n", data)
    }

    // Emit an event to the frontend
    app.Event.Emit("notification", result.Response)
})
```

## Настройка уведомлений

### Пользовательские метаданные

В простые и интерактивные уведомления можно включать пользовательские данные:

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID: "unique-id",
    Title: "New Calendar Invite",
    Subtitle: "From: Jane Doe", // Optional
    Body: "Tap to view the event",
    Data: map[string]interface{}{
        "sender": "jane.doe@example.com",
        "timestamp": "2025-03-10T15:30:00Z",
    }
})
```

### Пользовательский звук

Используйте `Sound`, чтобы управлять звуком при доставке уведомления. Если оставить значение `nil`, будет воспроизведён стандартный звук платформы; задайте `Silent: true`, чтобы отключить звук; задайте `Name`, чтобы воспроизвести звук с указанным именем или из комплекта приложения.

```go
// Silent
notifier.SendNotification(notifications.NotificationOptions{
    ID:    "silent-id",
    Title: "Background sync complete",
    Sound: &notifications.NotificationSound{Silent: true},
})

// Named sound
notifier.SendNotification(notifications.NotificationOptions{
    ID:    "named-id",
    Title: "New message",
    Body:  "Tap to read",
    Sound: &notifications.NotificationSound{Name: "Ping"},
})
```

Разрешение значения `Name` на каждой платформе:

- **macOS** — значение `Name` передаётся в `[UNNotificationSound soundNamed:]`; аудиофайл должен находиться в каталоге `Library/Sounds` комплекта приложения.
- **Windows** — если значение `Name` уже начинается с `ms-winsoundevent:` или `ms-appx:`, оно используется без изменений; в противном случае оно помещается в `ms-winsoundevent:` и используется как имя встроенного события всплывающего уведомления (см. документацию Microsoft по схеме всплывающих уведомлений `<audio>`).
- **Linux** — значение передаётся как подсказка freedesktop `sound-name`; воспроизведение зависит от активной службы уведомлений и звуковой темы.

### Вложения

`Attachments` добавляет к уведомлению медиафайлы. macOS поддерживает несколько вложений с медиафайлами любого типа; Windows и Linux обрабатывают первое вложение типа «изображение».

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID:    "image-id",
    Title: "Photo uploaded",
    Body:  "Tap to view",
    Attachments: []notifications.NotificationAttachment{
        {
            ID:   "preview",
            Path: "/absolute/path/to/image.png",
            // Type hint:
            //   macOS:   UTI like "public.png" / "public.audio" (often inferred)
            //   Windows: "hero" | "appLogoOverride" | "inline" (default "inline")
            //   Linux:   ignored
            Type: "inline",
        },
    },
})
```

`Path` должен содержать абсолютный путь в файловой системе. macOS также принимает URL-адреса `file://`.

#### Прикрепление файла из комплекта приложения

ОС считывает вложение с диска при доставке уведомления, поэтому значение `Path` должно разрешаться в реальный файл на компьютере конечного пользователя. Для ресурса из комплекта приложения (значка или изображения, встроенного с помощью `go:embed`) нельзя жёстко задать постоянный абсолютный путь: файл находится внутри исполняемого файла, а не в известном расположении на диске. Один раз при запуске запишите его в каталог, доступный для записи, и передайте полученный путь:

```go
import (
    _ "embed"
    "os"
    "path/filepath"
)

//go:embed assets/preview.png
var previewPNG []byte

// Materialise the embedded asset to a stable path the OS can read.
previewPath := filepath.Join(os.TempDir(), "myapp-preview.png")
if err := os.WriteFile(previewPath, previewPNG, 0o644); err != nil {
    // handle error
}

notifier.SendNotification(notifications.NotificationOptions{
    ID:    "image-id",
    Title: "Photo uploaded",
    Body:  "Tap to view",
    Attachments: []notifications.NotificationAttachment{
        {ID: "preview", Path: previewPath, Type: "inline"},
    },
})
```

У предоставленного пользователем или загруженного файла уже есть реальный путь на диске, поэтому его можно передать непосредственно в `Path`, пропустив этот шаг. Возможность передавать байты вложения из памяти может появиться в одном из будущих выпусков.

### Цепочки и группировка

`ThreadID` объединяет связанные уведомления, позволяя ОС сворачивать их в Центре уведомлений / Центре действий.

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID:       "msg-42",
    Title:    "Jane Doe",
    Body:     "Are you free for lunch?",
    ThreadID: "chat:jane.doe",
})
```

### Уровень прерывания

`InterruptionLevel` управляет приоритетом уведомления. Используйте одну из экспортируемых констант:

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID:                "alert-id",
    Title:             "Server down",
    Body:              "Investigate immediately",
    InterruptionLevel: notifications.InterruptionLevelTimeSensitive,
})
```

| Константа | Значение | Назначение |
| --- | --- | --- |
| `InterruptionLevelPassive` | `"passive"` | Тихая доставка: экран не включается и стандартный звук не воспроизводится |
| `InterruptionLevelActive` | `"active"` | Стандартный уровень |
| `InterruptionLevelTimeSensitive` | `"timeSensitive"` | Обходит режим «Фокусирование» / «Не беспокоить», если это разрешено |
| `InterruptionLevelCritical` | `"critical"` | Обходит режим фокусирования и беззвучный режим; в macOS требуется право Critical Alert (без него уровень незаметно понижается) |

Соответствие платформам:

- **macOS** — задаёт `UNNotificationContent.interruptionLevel`. Для `critical` требуется macOS 12+ и право Critical Alert.
- **Windows** — сопоставляется с атрибутом всплывающего уведомления `<toast scenario="...">`.
- **Linux** — сопоставляется с подсказкой freedesktop `urgency`.

### Отложенная доставка

`Schedule` откладывает доставку. Задайте ровно одно из значений: `DelaySeconds` (количество секунд от текущего момента) или `At` (время Unix в секундах, UTC).

```go
// Deliver in 60 seconds
notifier.SendNotification(notifications.NotificationOptions{
    ID:       "reminder-id",
    Title:    "Stand up and stretch",
    Schedule: &notifications.NotificationSchedule{DelaySeconds: 60},
})

// Deliver at an absolute time
notifier.SendNotification(notifications.NotificationOptions{
    ID:       "alarm-id",
    Title:    "Meeting in 5 minutes",
    Schedule: &notifications.NotificationSchedule{At: time.Now().Add(time.Hour).Unix()},
})
```

@note{type="caution" title="Сохранение после перезапуска"}
В **macOS** для запланированных уведомлений используется нативный триггер, поэтому они сохраняются после перезапуска приложения. В **Windows** и **Linux** планирование выполняется с помощью внутрипроцессного таймера `time.AfterFunc`, и уведомление **будет потеряно, если приложение завершит работу до доставки**: ни `wintoast`, ни спецификация freedesktop не предоставляют примитива отложенной доставки.

@end

### Обновление уведомлений

`UpdateNotification` заменяет активное уведомление с тем же `ID`:

```go
notifier.UpdateNotification(notifications.NotificationOptions{
    ID:    "download-id",
    Title: "Download complete",
    Body:  "report.pdf is ready",
})
```

Поведение на разных платформах:

- **macOS** — `UNUserNotificationCenter` автоматически устраняет дубликаты по идентификатору, поэтому существующее уведомление обновляется на месте.
- **Linux** — для замены предыдущего уведомления используется параметр D-Bus `replaces_id`.
- **Windows** — сейчас уведомление повторно доставляется как новое. Для полноценной замены на месте требуется поддержка `tag` / `group` в вышестоящем проекте `wintoast`.

## Особенности платформ

@tabs
[macOS]
Уведомления в macOS:

- Требуют разрешения пользователя
- Требуют, чтобы приложение было упаковано и подписано (а для распространения — нотариально заверено)
- Используют стандартный системный внешний вид уведомлений
- Поддерживают `Subtitle`
- Поддерживают ввод текста пользователем (ответы)
- Поддерживают параметр действия `Destructive`
- Поддерживают несколько `Attachments` любого типа мультимедиа (изображения, аудио и видео)
- Поддерживают `ThreadID` для группировки в Центре уведомлений
- Поддерживают все значения `InterruptionLevel` (для `critical` требуется право Critical Alert)
- Поддерживают нативную отложенную доставку, сохраняющуюся после перезапуска приложения
- Автоматически устраняют дубликаты вызовов `UpdateNotification` по `ID`
- Автоматически учитывают тёмный или светлый режим

[Windows]
Уведомления в Windows:

- Используют системные стили всплывающих уведомлений Windows через бэкенд `wintoast`
- Адаптируются к настройкам темы Windows
- Поддерживают ввод текста пользователем (ответы)
- Поддерживают дисплеи с высокой плотностью пикселей
- Не поддерживают `Subtitle`
- Поддерживают одно изображение `Attachment` с подсказкой размещения `hero`, `appLogoOverride` или `inline` (по умолчанию — `inline`)
- Поддерживают `ThreadID` для группировки в Центре уведомлений
- Поддерживают `InterruptionLevel` через атрибут всплывающего уведомления `scenario`
- Поддерживают отложенную доставку с помощью внутрипроцессного таймера — **запланированные уведомления будут потеряны, если приложение завершит работу до доставки**
- `UpdateNotification` сейчас повторно доставляет уведомление как новое (полноценная замена на месте ожидает поддержки `tag`/`group` в вышестоящем проекте `wintoast`)

[Linux]
В Linux уведомления используют интерфейс D-Bus `org.freedesktop.Notifications`. Для работы уведомлений совместимый демон уведомлений **должен быть запущен**.

@note{type="caution" title="Системное требование: демон уведомлений"}
Совместимый с freedesktop демон уведомлений должен быть установлен и запущен. Распространённые варианты:

- **dunst** — легковесный и гибко настраиваемый (`apt install dunst` / `dnf install dunst`)
- **mako** — нативный для Wayland (`apt install mako-notifier`)
- **GNOME Shell** — автоматически регистрирует интерфейс в GNOME 43+. В Ubuntu 24.04 (GNOME Shell 46) интерфейс может не зарегистрироваться автоматически при запуске сеанса; если уведомления не появляются, установите `dunst` в качестве резервного варианта.
- **xfce4-notifyd** — входит в состав рабочих окружений XFCE

Если демон не запущен, `SendNotification` вернёт ошибку D-Bus: `The name org.freedesktop.Notifications was not provided by any .service files`. Обработайте эту ошибку в приложении и сообщите пользователю, что необходимо установить демон уведомлений.

@end

Уведомления в Linux:

- Следуют теме рабочего окружения
- Располагаются в соответствии с правилами рабочего окружения
- Поддерживают `Subtitle` (для демонов, которые не отображают его отдельно, он добавляется к тексту уведомления)
- Не поддерживают ввод текста пользователем (он не входит в спецификацию freedesktop)
- Поддерживают одно изображение `Attachment` через подсказку `image-path`
- Поддерживают `ThreadID` (обрабатывается службой уведомлений, если она это поддерживает)
- `Sound.Name` передаётся как подсказка `sound-name`; будет ли воспроизводиться звук, зависит от активной службы уведомлений и звуковой темы
- Сопоставляют `InterruptionLevel` с подсказкой freedesktop `urgency`
- Поддерживают доставку по расписанию с помощью внутрипроцессного таймера — **если приложение завершит работу до момента доставки, запланированные уведомления будут потеряны**
- `UpdateNotification` использует параметр D-Bus `replaces_id`, чтобы заменить предыдущее уведомление на месте

@end

## Рекомендации

1. Проверяйте наличие разрешения и запрашивайте его:
  - В macOS требуется разрешение пользователя


2. Делайте уведомления понятными и лаконичными:
  - Используйте информативные заголовки, подзаголовки, текст и названия действий


3. Правильно обрабатывайте ответы на уведомления:
  - Проверяйте ответы на уведомления на наличие ошибок
  - Предоставляйте обратную связь о действиях пользователя


4. Учитывайте соглашения платформы:
  - Следуйте принятым на платформе шаблонам уведомлений
  - Учитывайте системные настройки


5. В Linux учитывайте зависимость от службы уведомлений:
  - Проверяйте ошибку, возвращаемую `SendNotification`: отсутствие службы уведомлений приводит к ошибке D-Bus
  - В документации пакета или README приложения следует указать, что требуется служба уведомлений, совместимая с freedesktop


## Примеры

Ознакомьтесь с этим примером:

- [Уведомления](https://github.com/wailsapp/wails/tree/master/v3/examples/notifications)

## Справочник API

### Управление службой

| Метод | Описание |
| --- | --- |
| `New()` | Создаёт новую службу уведомлений |

### Разрешение на уведомления

| Метод | Описание |
| --- | --- |
| `RequestNotificationAuthorization()` | Запрашивает разрешение на отображение уведомлений (macOS) |
| `CheckNotificationAuthorization()` | Проверяет текущий статус разрешения на уведомления (macOS) |

### Отправка уведомлений

| Метод | Описание |
| --- | --- |
| `SendNotification(options NotificationOptions)` | Отправляет простое уведомление |
| `SendNotificationWithActions(options NotificationOptions)` | Отправляет интерактивное уведомление с действиями |
| `UpdateNotification(options NotificationOptions)` | Обновляет находящееся в обработке уведомление по `ID` (см. раздел [«Обновление уведомлений»](#--10)) |

### Категории уведомлений

| Метод | Описание |
| --- | --- |
| `RegisterNotificationCategory(category NotificationCategory)` | Регистрирует категорию уведомлений для многократного использования |
| `RemoveNotificationCategory(categoryID string)` | Удаляет ранее зарегистрированную категорию |

### Управление уведомлениями

| Метод | Описание |
| --- | --- |
| `RemoveAllPendingNotifications()` | Удаляет все ожидающие уведомления (только в macOS и Linux) |
| `RemovePendingNotification(identifier string)` | Удаляет указанное ожидающее уведомление (только в macOS и Linux) |
| `RemoveAllDeliveredNotifications()` | Удаляет все доставленные уведомления (только в macOS и Linux) |
| `RemoveDeliveredNotification(identifier string)` | Удаляет указанное доставленное уведомление (только в macOS и Linux) |
| `RemoveNotification(identifier string)` | Удаляет уведомление (только в Linux) |

### Обработка событий

| Метод | Описание |
| --- | --- |
| `OnNotificationResponse(callback func(result NotificationResult))` | Регистрирует функцию обратного вызова для ответов на уведомления |

### Структуры и типы

#### NotificationOptions

```go
type NotificationOptions struct {
    ID         string                 // Unique identifier for the notification (required)
    Title      string                 // Main notification title (required)
    Subtitle   string                 // Subtitle text (macOS and Linux only)
    Body       string                 // Main notification content
    CategoryID string                 // Category identifier for interactive notifications
    Data       map[string]interface{} // Custom data to associate with the notification

    // Sound controls the sound played on delivery. nil = platform default.
    Sound *NotificationSound

    // Attachments are media files shown alongside the notification.
    // macOS supports multiple attachments of any media type;
    // Windows and Linux honour the first image-typed attachment.
    Attachments []NotificationAttachment

    // ThreadID groups related notifications together in
    // Notification Center / Action Center / the Linux notification daemon.
    ThreadID string

    // InterruptionLevel controls priority. One of "passive",
    // "active" (default), "timeSensitive", "critical".
    InterruptionLevel string

    // Schedule defers delivery. macOS uses a native trigger that survives
    // restarts; Windows and Linux use an in-process timer that does NOT.
    Schedule *NotificationSchedule
}
```

#### NotificationSound

```go
type NotificationSound struct {
    Silent bool   // If true, no sound is played
    Name   string // Named/bundled sound (see "Custom Sound" above)
}
```

#### NotificationAttachment

```go
type NotificationAttachment struct {
    ID   string // Optional identifier
    Path string // Absolute filesystem path (macOS also accepts file:// URLs)
    // Type is an optional placement/UTI hint:
    //   macOS:   UTI like "public.png" / "public.audio" (often inferred)
    //   Windows: "hero" | "appLogoOverride" | "inline" (default "inline")
    //   Linux:   ignored (always image-path hint)
    Type string
}
```

#### NotificationSchedule

```go
// Exactly one of DelaySeconds or At must be set. At is Unix seconds (UTC).
type NotificationSchedule struct {
    DelaySeconds int   // Seconds from now until delivery
    At           int64 // Absolute Unix timestamp (seconds, UTC)
}
```

#### Константы InterruptionLevel

```go
const (
    InterruptionLevelPassive       = "passive"
    InterruptionLevelActive        = "active" // default
    InterruptionLevelTimeSensitive = "timeSensitive"
    InterruptionLevelCritical      = "critical"
)
```

#### NotificationCategory

```go
type NotificationCategory struct {
    ID               string                // Unique identifier for the category
    Actions          []NotificationAction  // Button actions for the notification
    HasReplyField    bool                  // Whether to include a text input field
    ReplyPlaceholder string                // Placeholder text for the input field
    ReplyButtonTitle string                // Text for the reply button
}
```

#### NotificationAction

```go
type NotificationAction struct {
    ID          string  // Unique identifier for the action
    Title       string  // Button text
    Destructive bool    // Whether the action is destructive (macOS-specific)
}
```

#### NotificationResponse

```go
type NotificationResponse struct {
    ID               string                  // Notification identifier
    ActionIdentifier string                  // Action that was triggered
    CategoryID       string                  // Category of the notification
    Title            string                  // Title of the notification
    Subtitle         string                  // Subtitle of the notification
    Body             string                  // Body text of the notification
    UserText         string                  // Text entered by the user
    UserInfo         map[string]interface{}  // Custom data from the notification
}
```

#### NotificationResult

```go
type NotificationResult struct {
    Response NotificationResponse  // Response data
    Error    error                 // Any error that occurred
}
```
