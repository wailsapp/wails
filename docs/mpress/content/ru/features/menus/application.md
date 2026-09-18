---
title: "Меню приложения"
description: "Создавайте нативные строки меню для настольного приложения"
slug: "features/menus/application"
sourcePath: "features/menus/application.md"
---

## Проблема

Профессиональным настольным приложениям нужны строки меню — «Файл», «Правка», «Вид», «Справка». Однако на каждой платформе меню работают по-разному:

- **macOS**: глобальная строка меню в верхней части экрана
- **Windows**: строка меню в заголовке окна
- **Linux**: зависит от среды рабочего стола

Создавать подходящие для каждой платформы меню вручную утомительно и чревато ошибками.

## Решение Wails

Wails предоставляет **унифицированный API**, который автоматически создаёт нативные для платформы меню. Напишите код один раз и получите нативное поведение на всех платформах.

![Меню приложения Wails в macOS со стандартными пунктами, флажками, переключателями и подменю](/assets/screenshots/application-menu-macos.png)

В macOS меню приложения размещается в глобальной строке меню. На этом снимке показано нативное меню, отрисованное API меню Wails, включая отключённые пункты, флажки, переключатели и подменю.

## Быстрый старт

```go
package main

import (
    "runtime"
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "My App",
    })

    // Create menu
    menu := app.NewMenu()

    // Add standard menus (platform-appropriate)
    if runtime.GOOS == "darwin" {
        menu.AddRole(application.AppMenu)  // macOS only
    }
    menu.AddRole(application.FileMenu)
    menu.AddRole(application.EditMenu)
    menu.AddRole(application.WindowMenu)
    menu.AddRole(application.HelpMenu)

    // Set the application menu
    app.Menu.Set(menu)

    // Create window with UseApplicationMenu to inherit the menu on Windows/Linux
    app.Window.NewWithOptions(application.WebviewWindowOptions{
        UseApplicationMenu: true,
    })

    app.Run()
}
```

**Вот и всё!** Теперь у вас есть нативные для платформы меню со стандартными пунктами. Параметр `UseApplicationMenu` обеспечивает отображение меню в окнах Windows и Linux без дополнительного кода.

## Создание меню

### Базовое создание меню

```go
// Create a new menu
menu := app.NewMenu()

// Add a top-level menu
fileMenu := menu.AddSubmenu("File")

// Add menu items
fileMenu.Add("New").OnClick(func(ctx *application.Context) {
    // Handle New
})

fileMenu.Add("Open").OnClick(func(ctx *application.Context) {
    // Handle Open
})

fileMenu.AddSeparator()

fileMenu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})
```

### Назначение меню

**Рекомендуемый подход** — используйте `UseApplicationMenu` для единообразия на всех платформах:

```go
// Set the application menu once
app.Menu.Set(menu)

// Create windows that inherit the menu on Windows/Linux
app.Window.NewWithOptions(application.WebviewWindowOptions{
    UseApplicationMenu: true,  // Window uses the app menu
})
```

При таком подходе:

- В **macOS** меню отображается в верхней части экрана (стандартное поведение).
- В **Windows/Linux** каждое окно с `UseApplicationMenu: true` отображает меню приложения.

**Особенности отдельных платформ:**

@tabs{sync-key="platform"}
[macOS]
**Глобальная строка меню** (одна на приложение):

```go
app.Menu.Set(menu)
```

Меню отображается в верхней части экрана и остаётся доступным даже после закрытия всех окон. Параметр `UseApplicationMenu` не действует в macOS, поскольку все приложения используют глобальное меню.

[Windows]
**Строка меню отдельного окна**:

```go
// Option 1: Use application menu (recommended)
app.Menu.Set(menu)
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    UseApplicationMenu: true,
})

// Option 2: Set menu directly on window
window.SetMenu(menu)
```

У каждого окна может быть собственное меню либо меню, унаследованное от приложения. Меню отображается в заголовке окна.

[Linux]
**Строка меню отдельного окна** (как правило):

```go
// Option 1: Use application menu (recommended)
app.Menu.Set(menu)
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    UseApplicationMenu: true,
})

// Option 2: Set menu directly on window
window.SetMenu(menu)
```

Поведение зависит от среды рабочего стола. Некоторые среды, например Unity, поддерживают глобальные меню.

@end

@note{type="tip" title="Упрощение кроссплатформенных меню"}
Использование `UseApplicationMenu: true` устраняет необходимость в коде для отдельных платформ, например:

```go
// Old approach - no longer needed
if runtime.GOOS == "darwin" {
    app.Menu.Set(menu)
} else {
    window.SetMenu(menu)
}
```

@end

**Пользовательские меню отдельных окон:**

Если окну требуется меню, отличное от меню приложения, назначьте его непосредственно окну:

```go
window.SetMenu(customMenu)  // Overrides UseApplicationMenu
```

## Роли меню

Wails предоставляет **предопределённые роли меню**, которые автоматически создают подходящие для платформы структуры меню.

### Доступные роли

| Роль | Описание | Примечания по платформам |
| --- | --- | --- |
| `AppMenu` | Меню приложения с пунктами «О приложении», «Настройки» и «Выйти» | **Только macOS** |
| `FileMenu` | Операции с файлами («Создать», «Открыть», «Сохранить» и т. д.) | Все платформы |
| `EditMenu` | Редактирование текста («Отменить», «Повторить», «Вырезать», «Копировать», «Вставить») | Все платформы |
| `WindowMenu` | Управление окнами («Свернуть», «Масштабировать» и т. д.) | Все платформы |
| `HelpMenu` | Справка и информация | Все платформы |

### Использование ролей

```go
menu := app.NewMenu()

// macOS: Add application menu
if runtime.GOOS == "darwin" {
    menu.AddRole(application.AppMenu)
}

// All platforms: Add standard menus
menu.AddRole(application.FileMenu)
menu.AddRole(application.EditMenu)
menu.AddRole(application.WindowMenu)
menu.AddRole(application.HelpMenu)
```

**Что вы получите:**

@tabs{sync-key="platform"}
[macOS]
**AppMenu** (с названием приложения):

- О приложении [Название приложения]
- Настройки... (⌘,)
- ---
- Службы
- ---
- Скрыть [Название приложения] (⌘H)
- Скрыть остальные (⌥⌘H)
- Показать все
- ---
- Завершить [Название приложения] (⌘Q)

**FileMenu**:

- Создать (⌘N)
- Открыть... (⌘O)
- ---
- Закрыть окно (⌘W)

**EditMenu**:

- Отменить (⌘Z)
- Повторить (⇧⌘Z)
- ---
- Вырезать (⌘X)
- Копировать (⌘C)
- Вставить (⌘V)
- Выбрать всё (⌘A)

**WindowMenu**:

- Свернуть (⌘M)
- Масштабировать
- ---
- Переместить все окна на передний план

**HelpMenu**:

- Справка по [Название приложения]

[Windows]
**FileMenu**:

- Создать (Ctrl+N)
- Открыть... (Ctrl+O)
- ---
- Выход (Alt+F4)

**EditMenu**:

- Отменить (Ctrl+Z)
- Повторить (Ctrl+Y)
- ---
- Вырезать (Ctrl+X)
- Копировать (Ctrl+C)
- Вставить (Ctrl+V)
- Выбрать всё (Ctrl+A)

**WindowMenu**:

- Свернуть
- Развернуть

**HelpMenu**:

- О программе [Название приложения]

[Linux]
Аналогично Windows, но сочетания клавиш могут различаться в зависимости от окружения рабочего стола.

@end

### Настройка меню ролей

`Menu.AddRole(role)` возвращает меню-**получатель** (меню верхнего уровня), а **не** подменю роли. Чтобы добавить элементы в подменю роли, найдите добавленный элемент роли с помощью `FindByRole` и вызовите для него `GetSubmenu()`:

```go
menu.AddRole(application.FileMenu)

fileMenu := menu.FindByRole(application.FileMenu).GetSubmenu()
fileMenu.Add("Import...").OnClick(handleImport)
fileMenu.Add("Export...").OnClick(handleExport)
```

## Пользовательские меню

Создавайте собственные меню для функций, специфичных для приложения:

```go
// Add a custom top-level menu
toolsMenu := menu.AddSubmenu("Tools")

// Add items
toolsMenu.Add("Settings").OnClick(func(ctx *application.Context) {
    showSettingsWindow()
})

toolsMenu.AddSeparator()

// Add checkbox
toolsMenu.AddCheckbox("Dark Mode", false).OnClick(func(ctx *application.Context) {
    isDark := ctx.ClickedMenuItem().Checked()
    setTheme(isDark)
})

// Add radio group
toolsMenu.AddRadio("Small", true).OnClick(handleFontSize)
toolsMenu.AddRadio("Medium", false).OnClick(handleFontSize)
toolsMenu.AddRadio("Large", false).OnClick(handleFontSize)

// Add submenu
advancedMenu := toolsMenu.AddSubmenu("Advanced")
advancedMenu.Add("Configure...").OnClick(showAdvancedSettings)
```

**Другие типы элементов меню** описаны в разделе [«Справочник по меню»](/features/menus/reference/).

## Динамические меню

Обновляйте меню в зависимости от состояния приложения:

### Включение и отключение элементов

```go
var saveMenuItem *application.MenuItem

func createMenu() {
    menu := app.NewMenu()
    fileMenu := menu.AddSubmenu("File")
    
    saveMenuItem = fileMenu.Add("Save")
    saveMenuItem.SetEnabled(false)  // Initially disabled
    saveMenuItem.OnClick(handleSave)
    
    app.Menu.Set(menu)
}

func onDocumentChanged() {
    saveMenuItem.SetEnabled(hasUnsavedChanges())
    menu.Update()  // Important!
}
```

@note{type="caution" title="Всегда вызывайте menu.Update()"}
После изменения состояния меню (включения или отключения, подписи либо отметки) **всегда вызывайте `menu.Update()`**. Это особенно важно в Windows, где меню создаются заново.

Подробности см. в разделе [«Справочник по меню»](/features/menus/reference/#enabled-state).

@end

### Изменение подписей

```go
updateMenuItem := menu.Add("Check for Updates")

updateMenuItem.OnClick(func(ctx *application.Context) {
    updateMenuItem.SetLabel("Checking...")
    menu.Update()
    
    checkForUpdates()
    
    updateMenuItem.SetLabel("Check for Updates")
    menu.Update()
})
```

### Пересоздание меню

При значительных изменениях пересоздайте всё меню:

```go
func rebuildFileMenu() {
    menu := app.NewMenu()
    fileMenu := menu.AddSubmenu("File")
    
    fileMenu.Add("New").OnClick(handleNew)
    fileMenu.Add("Open").OnClick(handleOpen)
    
    // Add recent files dynamically
    if hasRecentFiles() {
        recentMenu := fileMenu.AddSubmenu("Open Recent")
        for _, file := range getRecentFiles() {
            filePath := file  // Capture for closure
            recentMenu.Add(filepath.Base(file)).OnClick(func(ctx *application.Context) {
                openFile(filePath)
            })
        }
        recentMenu.AddSeparator()
        recentMenu.Add("Clear Recent").OnClick(clearRecentFiles)
    }
    
    fileMenu.AddSeparator()
    fileMenu.Add("Quit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })
    
    app.Menu.Set(menu)
}
```

## Управление окнами из меню

Элементы меню могут управлять окнами:

```go
viewMenu := menu.AddSubmenu("View")

// Toggle fullscreen
viewMenu.Add("Toggle Fullscreen").OnClick(func(ctx *application.Context) {
    if window, ok := app.Window.GetByName("main"); ok {
        window.ToggleFullscreen()
    }
})

// Zoom controls
viewMenu.Add("Zoom In").SetAccelerator("CmdOrCtrl++").OnClick(func(ctx *application.Context) {
    // Increase zoom
})

viewMenu.Add("Zoom Out").SetAccelerator("CmdOrCtrl+-").OnClick(func(ctx *application.Context) {
    // Decrease zoom
})

viewMenu.Add("Reset Zoom").SetAccelerator("CmdOrCtrl+0").OnClick(func(ctx *application.Context) {
    // Reset zoom
})
```

**Получение активного окна:**

```go
menuItem.OnClick(func(ctx *application.Context) {
    window := application.Get().Window.Current() // the window the menu was invoked from
    // Use window
})
```

## Особенности отдельных платформ

### macOS

**Поведение строки меню:**

- Отображается **в верхней части экрана** (глобально)
- Остаётся после закрытия всех окон
- Первым **всегда идёт меню приложения**
- Для стандартных элементов используйте `menu.AddRole(application.AppMenu)`

**Стандартное расположение:**

- **О программе**: меню приложения
- **Настройки**: меню приложения (⌘,)
- **Завершить**: меню приложения (⌘Q)
- **Справка**: меню «Справка»

**Пример:**

```go
if runtime.GOOS == "darwin" {
    menu.AddRole(application.AppMenu)  // Adds About, Preferences, Quit
    
    // Don't add Quit to File menu on macOS
    // Don't add About to Help menu on macOS
}
```

### Windows

**Поведение строки меню:**

- Отображается в **строке заголовка окна**
- У каждого окна есть собственное меню
- Меню приложения отсутствует

**Стандартное расположение:**

- **Выход**: меню «Файл» (Alt+F4)
- **Настройки**: меню «Инструменты» или «Правка»
- **О программе**: меню «Справка»

**Пример:**

```go
if runtime.GOOS == "windows" {
    menu.AddRole(application.FileMenu) // Exit is added automatically
    menu.AddRole(application.HelpMenu) // About is added automatically
}
```

### Linux

**Поведение строки меню:**

- Обычно отдельная для каждого окна (как в Windows)
- Некоторые среды рабочего стола поддерживают глобальные меню (Unity, GNOME с расширением)
- Внешний вид зависит от среды рабочего стола

**Рекомендация:** следуйте соглашениям Windows и тестируйте в целевых средах рабочего стола.

## Полный пример

Ниже приведена готовая к использованию в рабочей среде структура меню:

```go
package main

import (
    "runtime"
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "My Application",
    })

    // Create and set menu
    createMenu(app)

    // Create main window with UseApplicationMenu for cross-platform menu support
    app.Window.NewWithOptions(application.WebviewWindowOptions{
        UseApplicationMenu: true,
    })

    app.Run()
}

func createMenu(app *application.App) {
    menu := app.NewMenu()

    // Platform-specific application menu (macOS only)
    if runtime.GOOS == "darwin" {
        menu.AddRole(application.AppMenu)
    }

    // File menu — AddRole returns the receiver menu, not the role submenu.
    // To add items into the File submenu, look it up via FindByRole + GetSubmenu.
    menu.AddRole(application.FileMenu)
    fileMenu := menu.FindByRole(application.FileMenu).GetSubmenu()
    fileMenu.Add("Import...").SetAccelerator("CmdOrCtrl+I").OnClick(handleImport)
    fileMenu.Add("Export...").SetAccelerator("CmdOrCtrl+E").OnClick(handleExport)

    // Edit menu
    menu.AddRole(application.EditMenu)

    // View menu
    viewMenu := menu.AddSubmenu("View")
    viewMenu.Add("Toggle Fullscreen").SetAccelerator("F11").OnClick(toggleFullscreen)
    viewMenu.AddSeparator()
    viewMenu.AddCheckbox("Show Sidebar", true).OnClick(toggleSidebar)
    viewMenu.AddCheckbox("Show Toolbar", true).OnClick(toggleToolbar)

    // Tools menu
    toolsMenu := menu.AddSubmenu("Tools")
    
    // Settings location varies by platform
    if runtime.GOOS == "darwin" {
        // On macOS, Preferences is in Application menu (added by AppMenu role)
    } else {
        toolsMenu.Add("Settings").SetAccelerator("CmdOrCtrl+,").OnClick(showSettings)
    }
    
    toolsMenu.AddSeparator()
    toolsMenu.AddCheckbox("Dark Mode", false).OnClick(toggleDarkMode)

    // Window menu
    menu.AddRole(application.WindowMenu)

    // Help menu
    helpMenu := menu.AddRole(application.HelpMenu)
    helpMenu.Add("Documentation").OnClick(openDocumentation)
    
    // About location varies by platform
    if runtime.GOOS == "darwin" {
        // On macOS, About is in Application menu (added by AppMenu role)
    } else {
        helpMenu.AddSeparator()
        helpMenu.Add("About").OnClick(showAbout)
    }

    // Set the application menu
    app.Menu.Set(menu)
}

func handleImport(ctx *application.Context) {
    // Implementation
}

func handleExport(ctx *application.Context) {
    // Implementation
}

func toggleFullscreen(ctx *application.Context) {
    window := application.Get().Window.Current()
    window.ToggleFullscreen()
}

func toggleSidebar(ctx *application.Context) {
    // Implementation
}

func toggleToolbar(ctx *application.Context) {
    // Implementation
}

func showSettings(ctx *application.Context) {
    // Implementation
}

func toggleDarkMode(ctx *application.Context) {
    isDark := ctx.ClickedMenuItem().Checked()
    // Apply theme
}

func openDocumentation(ctx *application.Context) {
    // Open browser
}

func showAbout(ctx *application.Context) {
    // Show about dialog
}
```

## Рекомендации

### ✅ Делайте так

- **Используйте роли меню** для стандартных меню («Файл», «Правка» и т. д.)
- **Следуйте соглашениям платформы** при формировании структуры меню
- **Добавляйте сочетания клавиш** для часто используемых действий
- **Вызывайте menu.Update()** после изменения состояния меню
- **Тестируйте на всех платформах** — поведение различается
- **Не создавайте глубокую иерархию меню** — не более 2-3 уровней
- **Используйте понятные названия** — «Сохранить проект», а не «Сохранить»

### ❌ Не делайте так

- **Не задавайте сочетания клавиш для платформы жёстко** — используйте `CmdOrCtrl`
- **Не помещайте команду «Выход» в меню «Файл» в macOS** — она находится в меню приложения
- **Не помещайте пункт «О программе» в меню «Справка» в macOS** — он находится в меню приложения
- **Не забывайте вызывать menu.Update()** — иначе меню не будут работать должным образом
- **Не создавайте слишком глубокую вложенность** — пользователи путаются
- **Не используйте жаргон** — названия должны быть понятны пользователям

## Дальнейшие шаги

@cards{cols="2"}
📖 Справочник по меню
Полный справочник по типам и свойствам элементов меню.

[Подробнее →](/features/menus/reference/)

---
◆ Контекстные меню
Создавайте контекстные меню, открываемые щелчком правой кнопки мыши.

[Подробнее →](/features/menus/context/)

---
★ Меню в области уведомлений
Добавьте интеграцию с областью уведомлений или строкой меню.

[Подробнее →](/features/menus/systray/)

---
📖 Шаблоны меню
Распространённые шаблоны меню и рекомендации.

[Подробнее →](/guides/menus/)

@end

---

**Есть вопросы?** Задайте их в [Discord](https://discord.gg/JDdSxwjhGf) или ознакомьтесь с [примером меню](https://github.com/wailsapp/wails/tree/master/v3/examples/menu).
