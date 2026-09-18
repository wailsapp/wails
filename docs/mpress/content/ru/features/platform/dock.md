---
title: "Dock и панель задач"
description: "Управление видимостью значка в Dock и отображение индикаторов в macOS и Windows"
slug: "features/platform/dock"
sourcePath: "features/platform/dock.md"
---

## Введение

Wails предоставляет кроссплатформенную службу Dock для настольных приложений. Эта служба позволяет:

- Скрывать и показывать значок приложения в Dock macOS
- Отображать индикаторы на плитке приложения или значке в Dock либо на панели задач (macOS и Windows)

## Основы использования

### Создание службы

Сначала инициализируйте службу Dock:

```go
import "github.com/wailsapp/wails/v3/pkg/application"
import "github.com/wailsapp/wails/v3/pkg/services/dock"

// Create a new Dock service
dockService := dock.New()

// Register the service with the application
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(dockService),
    },
})
```

### Создание службы с пользовательскими параметрами индикатора (только для Windows)

В Windows внешний вид индикатора можно настроить с помощью различных параметров:

```go
import "github.com/wailsapp/wails/v3/pkg/application"
import "github.com/wailsapp/wails/v3/pkg/services/dock"
import "image/color"

// Create a dock service with custom badge options
options := dock.BadgeOptions{
    TextColour:       color.RGBA{255, 255, 255, 255}, // White text
    BackgroundColour: color.RGBA{0, 0, 255, 255},     // Blue background
    FontName:         "consolab.ttf",                 // Bold Consolas font
    FontSize:         20,                             // Font size for single character
    SmallFontSize:    14,                             // Font size for multiple characters
}

dockService := dock.NewWithOptions(options)

// Register the service with the application
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(dockService),
    },
})
```

## Операции с Dock

### Скрытие значка приложения в Dock

Скройте значок приложения из Dock macOS:

```go
// Hide the app icon
dockService.HideAppIcon()
```

### Отображение значка приложения в Dock

Покажите значок приложения в Dock macOS:

```go
// Show the app icon
dockService.ShowAppIcon()
```

## Операции с индикаторами

### Установка индикатора

Установите индикатор на плитке приложения или значке в Dock:

```go
// Set a default badge
dockService.SetBadge("")

// Set a numeric badge
dockService.SetBadge("3")

// Set a text badge
dockService.SetBadge("New")
```

### Установка пользовательского индикатора (только для Windows)

Установите индикатор с однократно применяемыми параметрами:

```go
options := dock.BadgeOptions{
    BackgroundColour: color.RGBA{0, 255, 255, 255},
    FontName:         "arialb.ttf", // System font
    FontSize:         16,
    SmallFontSize:    10,
    TextColour:       color.RGBA{0, 0, 0, 255},
}

// Set a default badge
dockService.SetCustomBadge("", options)

// Set a numeric badge
dockService.SetCustomBadge("3", options)

// Set a text badge
dockService.SetCustomBadge("New", options)
```

### Удаление индикатора

Удалите индикатор со значка приложения:

```go
dockService.RemoveBadge()
```

### Получение установленного индикатора

```go
dockService.GetBadge()
```

## Особенности платформ

@tabs
[macOS]
В macOS:

- Значок в Dock можно **скрывать** и **показывать**
- Индикаторы отображаются непосредственно на значке в Dock
- Параметры индикатора **нельзя настраивать** (все параметры, переданные в `NewWithOptions`/`SetCustomBadge`, игнорируются)
- Используется стандартное оформление индикатора Dock в macOS, которое автоматически адаптируется к внешнему виду системы
- Переполнение текста обрабатывается системой
- Если указать пустой текст, отображается индикатор по умолчанию «●»

[Windows]
В Windows:

- Эта служба пока не поддерживает скрытие и отображение значка на панели задач
- Индикаторы отображаются на панели задач в виде значка-наложения
- Индикаторы поддерживают текстовые значения
- Внешний вид индикатора можно настроить с помощью `BadgeOptions`
- Для отображения индикаторов у приложения должно быть окно
- Для текста из нескольких символов автоматически используется меньший размер шрифта
- Переполнение текста не обрабатывается
- Параметры настройки:
  - **TextColour**: цвет текста (по умолчанию: белый)
  - **BackgroundColour**: цвет фона индикатора (по умолчанию: красный)
  - **FontName**: имя файла шрифта (по умолчанию: «segoeuib.ttf»)
  - **FontSize**: размер шрифта для одного символа (по умолчанию: 18)
  - **SmallFontSize**: размер шрифта для нескольких символов (по умолчанию: 14)


[Linux]
В Linux:

- Управление видимостью значка в Dock и функция индикаторов недоступны

@end

## Рекомендации

1. **При скрытии значка в Dock (macOS):**
  - Убедитесь, что пользователи по-прежнему могут открыть приложение (например, через [область уведомлений](/features/menus/systray/))
  - Добавьте пункт «Выйти» в альтернативный пользовательский интерфейс
  - Приложение не будет отображаться в переключателе Command+Tab
  - Открытые окна останутся видимыми и работоспособными
  - Закрытие всех окон может не завершить работу приложения (поведение macOS различается)
  - Пользователи лишатся стандартного способа завершить работу приложения через контекстное меню значка в Dock


2. **Используйте индикаторы умеренно:**
  - Слишком частое обновление индикатора может отвлекать пользователей
  - Используйте индикаторы только для важных уведомлений


3. **Используйте короткий текст индикатора:**
  - Числовые индикаторы наиболее эффективны
  - В macOS текстовые индикаторы рекомендуется делать краткими


4. **При настройке индикатора в Windows:**
  - Обеспечьте высокую контрастность между цветами текста и фона
  - Проверяйте отображение текста разной длины, поскольку с увеличением длины текста размер шрифта уменьшается
  - Используйте распространённые системные шрифты, чтобы гарантировать их наличие


## Справочник API

### Управление службой

| Метод | Описание |
| --- | --- |
| `New()` | Создаёт новую службу Dock |
| `NewWithOptions(options BadgeOptions)` | Создаёт новую службу Dock с настраиваемыми параметрами значка (только для Windows; в macOS и Linux параметры игнорируются) |

### Операции с Dock

| Метод | Описание |
| --- | --- |
| `HideAppIcon()` | Скрывает значок приложения из Dock в macOS (только для macOS) |
| `ShowAppIcon()` | Показывает значок приложения в Dock в macOS (только для macOS) |

### Операции со значком

| Метод | Описание |
| --- | --- |
| `SetBadge(label string) error` | Устанавливает значок с указанной меткой |
| `SetCustomBadge(label string, options BadgeOptions) error` | Устанавливает значок с указанной меткой и настраиваемыми параметрами оформления (только для Windows) |
| `RemoveBadge() error` | Удаляет значок с иконки приложения |
| `GetBadge() *string` | Возвращает текущий значок |

### Структуры и типы

```go
// Options for customizing badge appearance (Windows only)
type BadgeOptions struct {
    TextColour       color.RGBA  // Color of the badge text
    BackgroundColour color.RGBA  // Color of the badge background
    FontName         string      // Font file name (e.g., "segoeuib.ttf")
    FontSize         int         // Font size for single character
    SmallFontSize    int         // Font size for multiple characters
}
```
