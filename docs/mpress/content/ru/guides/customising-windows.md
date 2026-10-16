---
title: "Настройка окон в Wails"
description: "Настройка внешнего вида и поведения окон в приложениях Wails"
slug: "guides/customising-windows"
sourcePath: "guides/customising-windows.md"
---

Поддерживаемые платформы: <span class="mpress-badge mpress-badge-note">Windows</span> <span class="mpress-badge mpress-badge-success">macOS</span>

<br/>

Wails предоставляет API для управления внешним видом и функциональностью элементов управления окна. Эта возможность доступна в Windows и macOS, но не в Linux.

## Настройка состояний кнопок окна

Состояния кнопок определяются перечислением `ButtonState`:

```go
type ButtonState int

const (
    ButtonEnabled   ButtonState = 0
    ButtonDisabled  ButtonState = 1
    ButtonHidden    ButtonState = 2
)
```

- `ButtonEnabled`: кнопка включена и видима.
- `ButtonDisabled`: кнопка видима, но отключена (отображается серым цветом).
- `ButtonHidden`: кнопка скрыта с панели заголовка.

Состояния кнопок можно задать при создании окна или изменить во время выполнения.

### Настройка состояний кнопок при создании окна

При создании нового окна начальные состояния кнопок можно задать с помощью структуры `WebviewWindowOptions`:

```go {title="main.go"}
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name:        "My Application",
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        MinimiseButtonState:   application.ButtonHidden,
        MaximiseButtonState:   application.ButtonDisabled,
        CloseButtonState:      application.ButtonEnabled,
        FullscreenButtonState: application.ButtonEnabled,
    })

    app.Run()
}
```

В приведённом выше примере кнопка сворачивания скрыта, кнопка разворачивания неактивна (отображается серым цветом), а кнопка закрытия активна.

### Настройка состояний кнопок во время выполнения

Состояния кнопок также можно изменять во время выполнения с помощью следующих методов интерфейса `Window`:

```go
window.SetMinimiseButtonState(wails.ButtonHidden)
window.SetMaximiseButtonState(wails.ButtonEnabled)
window.SetCloseButtonState(wails.ButtonDisabled)
window.SetFullscreenButtonState(wails.ButtonEnabled)
```

### macOS: MaximiseButtonState и FullscreenButtonState относятся к одной кнопке

В macOS зелёная кнопка «светофора» (`NSWindowZoomButton`) — это один и тот же физический элемент управления как для разворачивания окна, так и для перехода в полноэкранный режим. Если при создании окна задать разные значения для `MaximiseButtonState` и `FullscreenButtonState`, одно из них без предупреждения переопределило бы другое по принципу «последняя запись имеет приоритет».

Чтобы избежать этого, при инициализации Wails применяет **более строгое** из двух состояний в следующем порядке: `ButtonEnabled` < `ButtonDisabled` < `ButtonHidden`.

| `MaximiseButtonState` | `FullscreenButtonState` | Фактическое состояние в macOS |
| --- | --- | --- |
| `ButtonEnabled` | `ButtonEnabled` | `ButtonEnabled` |
| `ButtonDisabled` | `ButtonEnabled` | `ButtonDisabled` |
| `ButtonEnabled` | `ButtonHidden` | `ButtonHidden` |
| `ButtonDisabled` | `ButtonHidden` | `ButtonHidden` |

Во время выполнения `SetMaximiseButtonState` и `SetFullscreenButtonState` в macOS воздействуют на `NSWindowZoomButton`, поэтому приоритет имеет последний вызов.

### Различия между платформами

Управление состояниями кнопок работает немного по-разному в Windows и macOS:

|  | Windows | Mac |
| --- | --- | --- |
| Отключить кнопки сворачивания, разворачивания и закрытия | Отключает кнопки сворачивания, разворачивания и закрытия | Отключает кнопки сворачивания, разворачивания и закрытия |
| Скрыть кнопку сворачивания | Отключает кнопку сворачивания | Скрывает кнопку сворачивания |
| Скрыть кнопку разворачивания | Отключает кнопку разворачивания | Скрывает кнопку разворачивания |
| Скрыть кнопку закрытия | Скрывает все элементы управления | Скрывает кнопку закрытия |
| `FullscreenButtonState` | Не выполняет никаких действий | Управляет кнопкой масштабирования (зелёной) |

Примечание: в Windows невозможно скрыть кнопки сворачивания и разворачивания по отдельности. Однако если отключить обе кнопки, они обе будут скрыты и останется только кнопка закрытия. В стандартной строке заголовка Windows нет отдельной кнопки полноэкранного режима, поэтому `FullscreenButtonState` в Windows не выполняет никаких действий.

### Управление стилем окна (Windows)

Чтобы управлять стилем строки заголовка в Windows, используйте поле `ExStyle` структуры `WebviewWindowOptions`:

Пример:

```go {title="main.go"}
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/w32"
)

func main() {
    app := application.New(application.Options{
        Name: "My Application",
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Windows: application.WindowsWindow{
            ExStyle: w32.WS_EX_TOOLWINDOW | w32.WS_EX_NOREDIRECTIONBITMAP | w32.WS_EX_TOPMOST,
        },
    })

    app.Run()
}
```

Эта настройка переопределяет другие параметры, влияющие на расширенный стиль окна:

- HiddenOnTaskbar
- AlwaysOnTop
- IgnoreMouseEvents
- BackgroundType
