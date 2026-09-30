---
title: "Меню"
description: "Руководство по созданию и настройке меню в Wails v3"
slug: "guides/menus"
sourcePath: "guides/menus.md"
---

Wails v3 предоставляет мощную систему меню, которая позволяет создавать как меню приложения, так и контекстные меню. В этом руководстве описаны различные функции и возможности системы меню.

## Создание меню

Чтобы создать новое меню, используйте метод `New()` диспетчера Menus:

```go
menu := app.Menu.New()
```

### Добавление пунктов меню

Wails поддерживает несколько типов пунктов меню, каждый из которых предназначен для определённой задачи:

#### Обычные пункты меню

Обычные пункты меню — это основные элементы меню. Они отображают текст и могут выполнять действия при нажатии:

```go
menuItem := menu.Add("Click Me")
```

#### Флажки

Пункты меню с флажками имеют переключаемое состояние и удобны для включения и отключения функций или параметров:

```go
checkbox := menu.AddCheckbox("My checkbox", true)  // true = initially checked
```

#### Группы переключателей

Группы переключателей позволяют пользователям выбрать один вариант из набора взаимоисключающих вариантов. Они создаются автоматически, когда пункты-переключатели расположены рядом друг с другом:

```go
menu.AddRadio("Option 1", true)   // true = initially selected
menu.AddRadio("Option 2", false)
menu.AddRadio("Option 3", false)
```

#### Разделители

Разделители — это горизонтальные линии, которые помогают объединять пункты меню в логические группы:

```go
menu.AddSeparator()
```

#### Подменю

Подменю — это вложенные меню, которые появляются при наведении указателя на пункт меню или при нажатии на него. Они удобны для организации сложных структур меню:

```go
submenu := menu.AddSubmenu("File")
submenu.Add("Open")
submenu.Add("Save")
```

#### Объединение меню

Одно меню можно добавить в другое, поместив его в конец или в начало.

```go
menu := app.Menu.New()
menu.Add("First Menu")

secondaryMenu := app.Menu.New()
secondaryMenu.Add("Second Menu")

// insert 'secondaryMenu' after 'menu'
menu.Append(secondaryMenu)

// insert 'secondaryMenu' before 'menu'
menu.Prepend(secondaryMenu)

// update the menu
menu.Update()
```

@note{type="info"}
По умолчанию при использовании `prepend` и `append` состояние будет общим с исходным меню. Чтобы создать новое меню с собственным состоянием, вызовите для меню `.Clone()`.

Например: `menu.Append(secondaryMenu.Clone())`

@end

#### Очистка меню

При работе с переменным количеством пунктов меню в некоторых случаях лучше создать совершенно новое меню.

При этом из существующего меню будут удалены все пункты, после чего их можно будет добавить заново.

```go
menu := app.Menu.New()
menu.Add("Waiting for update...")

// after certain logic, the menu has to be updated
menu.Clear()
menu.Add("Update complete!")
menu.Update()
```

@note{type="info"}
При очистке меню удаляются только пункты верхнего уровня. Хотя подменю перестанут отображаться, они по-прежнему будут занимать память, поэтому внимательно управляйте меню.

@end

#### Удаление меню

Чтобы очистить меню и освободить его ресурсы, используйте метод `Destroy()`:

```go
menu := app.Menu.New()
menu.Add("Waiting for update...")

// after certain logic, the menu has to be destroyed
menu.Destroy()
```

### Свойства пунктов меню

Для пунктов меню можно настроить несколько свойств:

| Свойство | Метод | Описание |
| --- | --- | --- |
| Метка | `SetLabel(string)` | Задаёт отображаемый текст |
| Включён | `SetEnabled(bool)` | Включает или отключает пункт |
| Отмечен | `SetChecked(bool)` | Задаёт состояние отметки (для флажков и переключателей) |
| Всплывающая подсказка | `SetTooltip(string)` | Задаёт текст всплывающей подсказки |
| Скрыт | `SetHidden(bool)` | Показывает или скрывает пункт |
| Сочетание клавиш | `SetAccelerator(string)` | Задаёт сочетание клавиш |

### Состояния пунктов меню

Пункты меню могут находиться в различных состояниях, управляющих их видимостью и доступностью для взаимодействия:

#### Видимость

Пункты меню можно динамически показывать и скрывать с помощью метода `SetHidden()`:

```go
menuItem := menu.Add("Dynamic Item")

// Hide the menu item
menuItem.SetHidden(true)

// Show the menu item
menuItem.SetHidden(false)

// Check current visibility
isHidden := menuItem.Hidden()
```

Скрытые пункты полностью удаляются из меню до тех пор, пока не будут показаны снова. Это удобно для контекстно-зависимых пунктов, которые должны отображаться только в определённых состояниях приложения.

#### Состояние доступности

Пункты меню можно включать и отключать с помощью метода `SetEnabled()`:

```go
menuItem := menu.Add("Save")

// Disable the menu item
menuItem.SetEnabled(false)  // Item appears grayed out and cannot be clicked

// Enable the menu item
menuItem.SetEnabled(true)   // Item becomes clickable again

// Check current enabled state
isEnabled := menuItem.Enabled()
```

Отключённые пункты меню остаются видимыми, но отображаются серым цветом, и нажать на них нельзя. Обычно это указывает, что действие сейчас недоступно, например:

- Отключение пункта «Сохранить», когда нет изменений для сохранения
- Отключение пункта «Копировать», когда ничего не выбрано
- Отключение пункта «Отменить», когда нет действия, которое можно отменить

#### Динамическое управление состояниями

Эти состояния можно сочетать с обработчиками событий для создания динамических меню:

```go
saveMenuItem := menu.Add("Save")

// Initially disable the Save menu item
saveMenuItem.SetEnabled(false)

// Enable Save only when there are unsaved changes
documentChanged := func() {
    saveMenuItem.SetEnabled(true)
    menu.Update()  // Remember to update the menu after changing states
}

// Disable Save after saving
documentSaved := func() {
    saveMenuItem.SetEnabled(false)
    menu.Update()
}
```

### Обработка событий

Пункты меню могут обрабатывать события нажатия с помощью метода `OnClick`:

```go
menuItem.OnClick(func(ctx *application.Context) {
    // Handle the click event
    println("Menu item clicked!")
})
```

Контекст содержит сведения о выбранном пункте меню:

```go
menuItem.OnClick(func(ctx *application.Context) {
    // Get the clicked menu item
    clickedItem := ctx.ClickedMenuItem()
    // Get its current state
    isChecked := clickedItem.Checked()
})
```

### Пункты меню на основе ролей

Wails предоставляет набор предопределённых ролей меню, которые автоматически создают пункты меню со стандартной функциональностью. Поддерживаются следующие роли меню:

#### Готовые структуры меню

Эти роли создают целые структуры меню с типовой функциональностью:

| Роль | Описание | Примечания о платформах |
| --- | --- | --- |
| `AppMenu` | Меню приложения с пунктами «О приложении», «Службы», «Скрыть/Показать» и «Завершить работу» | Только macOS |
| `EditMenu` | Стандартное меню «Правка» с командами «Отменить», «Повторить», «Вырезать», «Копировать», «Вставить» и другими | Все платформы |
| `ViewMenu` | Меню «Вид» с командами перезагрузки, масштабирования и перехода в полноэкранный режим | Все платформы |
| `WindowMenu` | Команды управления окном («Свернуть», «Масштабировать» и другие) | Все платформы |
| `HelpMenu` | Меню «Справка» со ссылкой «Подробнее» на сайт Wails | Все платформы |

#### Отдельные пункты меню

Эти роли позволяют добавлять отдельные пункты меню:

| Роль | Описание | Примечания о платформах |
| --- | --- | --- |
| `About` | Показать диалоговое окно «О приложении» | Все платформы |
| `Hide` | Скрыть приложение | Только macOS |
| `HideOthers` | Скрыть другие приложения | Только macOS |
| `UnHide` | Показать скрытое приложение | Только macOS |
| `CloseWindow` | Закрыть текущее окно | Все платформы |
| `Minimise` | Свернуть окно | Все платформы |
| `Zoom` | Масштабировать окно | Только macOS |
| `Front` | Переместить окно на передний план | Только macOS |
| `Quit` | Завершить работу приложения | Все платформы |
| `Undo` | Отменить последнее действие | Все платформы |
| `Redo` | Повторить последнее действие | Все платформы |
| `Cut` | Вырезать выделенное | Все платформы |
| `Copy` | Копировать выделенное | Все платформы |
| `Paste` | Вставить из буфера обмена | Все платформы |
| `PasteAndMatchStyle` | Вставить и согласовать стиль | Только macOS |
| `SelectAll` | Выделить всё | Все платформы |
| `Delete` | Удалить выделенное | Все платформы |
| `Reload` | Перезагрузить текущую страницу | Все платформы |
| `ForceReload` | Принудительно перезагрузить текущую страницу | Все платформы |
| `ToggleFullscreen` | Переключить полноэкранный режим | Все платформы |
| `ResetZoom` | Сбросить масштаб | Все платформы |
| `ZoomIn` | Увеличить масштаб | Все платформы |
| `ZoomOut` | Уменьшить масштаб | Все платформы |

Ниже приведён пример использования как полных меню, так и отдельных ролей:

```go
menu := app.Menu.New()

// Add complete menu structures
menu.AddRole(application.AppMenu)    // macOS only
menu.AddRole(application.EditMenu)   // Common edit operations
menu.AddRole(application.ViewMenu)   // View controls
menu.AddRole(application.WindowMenu) // Window controls

// Add individual role-based items to a custom menu
fileMenu := menu.AddSubmenu("File")
fileMenu.AddRole(application.CloseWindow)
fileMenu.AddSeparator()
fileMenu.AddRole(application.Quit)
```

## Меню приложения

Меню приложения отображаются в верхней части его окна (Windows/Linux) или в верхней части экрана (macOS).

### Поведение меню приложения

Когда вы задаёте меню приложения с помощью `app.Menu.Set()`, в macOS оно становится главным меню. В Windows/Linux меню задаются отдельно для каждого окна.

```go
app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Custom Menu Window",
    Windows: application.WindowsWindow{
        Menu: customMenu,  // Override application menu for this window
    },
})
```

Ниже приведён полный пример, демонстрирующий эти различия в поведении меню:

```go
func main() {
    app := application.New(application.Options{})

    // Create application menu
    appMenu := app.Menu.New()
    fileMenu := appMenu.AddSubmenu("File")
    fileMenu.Add("New").OnClick(func(ctx *application.Context) {
        // This will be available in all windows unless overridden
        window := app.Window.Current()
        window.SetTitle("New Window")
    })
    
    // Set as application menu - this is for macOS
    app.Menu.Set(appMenu)

    // Window with custom menu on Windows
    customMenu := app.Menu.New()
    customMenu.Add("Custom Action")
    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: "Custom Menu",
        Windows: application.WindowsWindow{
            Menu: customMenu,
        },
    })

    app.Run()
}
```

## Контекстные меню

Контекстные меню — это всплывающие меню, которые появляются при щелчке правой кнопкой мыши по элементам приложения. Они обеспечивают быстрый доступ к действиям, относящимся к выбранному элементу.

### Контекстное меню по умолчанию

Контекстное меню по умолчанию — это встроенное контекстное меню WebView, предоставляющее такие системные операции, как:

- копирование, вырезание и вставка при работе с текстом
- элементы управления выделением текста
- параметры проверки орфографии

#### Управление контекстным меню по умолчанию

Вы можете управлять отображением контекстного меню по умолчанию с помощью CSS-свойства `--default-contextmenu`:

```html
<!-- Always show default context menu -->
<div style="--default-contextmenu: show">
    <input type="text" placeholder="Right-click for text operations"/>
    <textarea>Standard text operations available here</textarea>
</div>

<!-- Hide default context menu -->
<div style="--default-contextmenu: hide">
    <div class="custom-component">Custom context menu only</div>
</div>

<!-- Smart context menu behaviour (default) -->
<div style="--default-contextmenu: auto">
    <!-- Shows default menu when text is selected or in input fields -->
    <p>Select this text to see the default menu</p>
    <input type="text" placeholder="Default menu for input operations"/>
</div>
```

@note{type="info"}
Эта функция будет работать должным образом только после того, как [среда выполнения фронтенда будет готова](/reference/frontend-runtime/).

@end

#### Поведение вложенных контекстных меню

При использовании свойства `--default-contextmenu` для вложенных элементов действуют следующие правила:

1. Дочерние элементы наследуют настройку контекстного меню родительского элемента, если она не переопределена явно
2. Приоритет имеет наиболее конкретная (ближайшая) настройка
3. Значение `auto` можно использовать для возврата к поведению по умолчанию

Пример поведения вложенных контекстных меню:

```html
<!-- Parent sets hide -->
<div style="--default-contextmenu: hide">
    <!-- This inherits hide -->
    <p>No context menu here</p>
    
    <!-- This overrides to show -->
    <div style="--default-contextmenu: show">
        <p>Context menu shown here</p>
        
        <!-- This inherits show -->
        <span>Also has context menu</span>
        
        <!-- This resets to automatic behaviour -->
        <div style="--default-contextmenu: auto">
            <p>Shows menu only when text is selected</p>
        </div>
    </div>
</div>
```

### Пользовательские контекстные меню

Пользовательские контекстные меню позволяют предоставлять действия, относящиеся к приложению и элементу, по которому щёлкнули. Они особенно полезны для следующих задач:

- операции с файлами в диспетчере документов
- инструменты обработки изображений
- пользовательские действия в таблице данных
- Операции, специфичные для компонента

#### Создание пользовательского контекстного меню

При создании пользовательского контекстного меню укажите уникальный идентификатор (имя), связывающий меню с HTML-элементами:

```go
// Create a context menu with identifier "imageMenu"
contextMenu := app.ContextMenu.New()
app.ContextMenu.Add("imageMenu", contextMenu)
```

Параметр name (в этом примере — "imageMenu") служит уникальным идентификатором и используется для следующих целей:

1. Связывать HTML-элементы с этим конкретным контекстным меню
2. Определять, какое меню следует отображать при щелчке правой кнопкой мыши
3. Обновлять и удалять меню

#### Контекстные данные

При обработке событий контекстного меню доступны как выбранный пункт меню, так и связанные с ним контекстные данные:

```go
contextMenu.Add("Process").OnClick(func(ctx *application.Context) {
    // Get the clicked menu item
    menuItem := ctx.ClickedMenuItem()
    
    // Get the context data as a string
    contextData := ctx.ContextMenuData()
    
    // Check if the menu item is checked (for checkbox/radio items)
    isChecked := ctx.IsChecked()
    
    // Use the data
    if contextData != "" {
        processItem(contextData)
    }
})
```

Контекстные данные передаются из свойства `--custom-contextmenu-data` HTML-элемента и доступны в обработчике щелчка через `ctx.ContextMenuData()`. Это особенно полезно в следующих случаях:

- Работа со списками или таблицами, где каждому элементу требуется уникальный идентификатор
- Выполнение операций над конкретными компонентами или элементами
- Передача состояния или метаданных из фронтенда в бэкенд

#### Управление контекстным меню

После внесения изменений в контекстное меню вызовите метод `Update()`, чтобы применить их:

```go
contextMenu.Update()
```

Когда контекстное меню больше не нужно, его можно уничтожить:

```go
contextMenu.Destroy()
```

@note{type="danger" title="Предупреждение"}
После вызова `Destroy()` повторное использование ссылки на контекстное меню приведёт к панике.

@end

### Практический пример: галерея изображений

Ниже приведён полный пример реализации пользовательского контекстного меню для галереи изображений:

```go
// Backend: Create the context menu
imageMenu := app.ContextMenu.New()
app.ContextMenu.Add("imageMenu", imageMenu)

// Add relevant operations
imageMenu.Add("View Full Size").OnClick(func(ctx *application.Context) {
    // Get the image ID from context data
    if imageID := ctx.ContextMenuData(); imageID != "" {
        openFullSizeImage(imageID)
    }
})

imageMenu.Add("Download").OnClick(func(ctx *application.Context) {
    if imageID := ctx.ContextMenuData(); imageID != "" {
        downloadImage(imageID)
    }
})

imageMenu.Add("Share").OnClick(func(ctx *application.Context) {
    if imageID := ctx.ContextMenuData(); imageID != "" {
        showShareDialog(imageID)
    }
})
```

```html
<!-- Frontend: Image gallery implementation -->
<div class="gallery">
    <!-- Each image container with context menu -->
    <div class="image-container" 
         style="--custom-contextmenu: imageMenu; --custom-contextmenu-data: img_123">
        <img src="/images/img_123.jpg" alt="Gallery Image"/>
        <span class="caption">Nature Photo</span>
    </div>
    
    <div class="image-container"
         style="--custom-contextmenu: imageMenu; --custom-contextmenu-data: img_124">
        <img src="/images/img_124.jpg" alt="Gallery Image"/>
        <span class="caption">City Photo</span>
    </div>
</div>
```

В этом примере:

1. Контекстное меню создаётся с идентификатором "imageMenu"
2. Каждый контейнер изображения связывается с меню с помощью `--custom-contextmenu: imageMenu`
3. Каждый контейнер передаёт идентификатор своего изображения в качестве контекстных данных с помощью `--custom-contextmenu-data`
4. Бэкенд получает идентификатор изображения в обработчиках щелчка и может выполнять соответствующие операции
5. Одно и то же меню используется повторно для всех изображений, а контекстные данные указывают, над каким изображением следует выполнить операцию

Этот шаблон особенно эффективен в следующих случаях:

- Таблицы данных, где для отдельных строк требуются определённые операции
- Файловые менеджеры, где для файлов требуются действия, зависящие от контекста
- Инструменты проектирования, где для разных элементов требуются разные операции
- Любые компоненты, в которых одни и те же операции применяются к нескольким экземплярам
