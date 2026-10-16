---
title: "API меню"
description: "Полный справочник по API меню"
slug: "reference/menu"
sourcePath: "reference/menu.md"
---

## Обзор

API меню предоставляет методы для создания и управления меню приложения, контекстными меню и меню в области уведомлений.

**Типы меню:**

- **Меню приложения** — верхняя строка меню («Файл», «Правка» и т. д.)
- **Контекстные меню** — меню, открываемые щелчком правой кнопки мыши
- **Меню в области уведомлений** — меню в системном трее (области уведомлений)

## Создание меню

### NewMenu()

Создаёт новое меню.

```go
func (a *App) NewMenu() *Menu
```

**Пример:**

```go
menu := app.NewMenu()
```

## Методы меню

### Add()

Добавляет пункт в меню.

```go
func (m *Menu) Add(label string) *MenuItem
```

**Параметры:**

- `label` — текст, отображаемый для пункта меню

**Возвращает:** созданный пункт меню

**Пример:**

```go
item := menu.Add("Open File")
item.OnClick(func(ctx *application.Context) {
    // Handle click
})
```

### AddSubmenu()

Добавляет подменю в меню.

```go
func (m *Menu) AddSubmenu(label string) *Menu
```

**Параметры:**

- `label` — метка подменю

**Возвращает:** созданное подменю

**Пример:**

```go
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("New")
fileMenu.Add("Open")
fileMenu.Add("Save")
```

### AddSeparator()

Добавляет визуальную линию-разделитель между пунктами меню.

```go
func (m *Menu) AddSeparator()
```

**Пример:**

```go
menu.Add("Copy")
menu.Add("Paste")
menu.AddSeparator()
menu.Add("Select All")
```

**Рекомендация:** используйте разделители для группировки связанных пунктов меню.

### AddCheckbox()

Добавляет пункт меню с флажком.

```go
func (m *Menu) AddCheckbox(label string, checked bool) *MenuItem
```

**Параметры:**

- `label` — метка флажка
- `checked` — исходное состояние флажка

**Пример:**

```go
darkMode := menu.AddCheckbox("Dark Mode", false)
darkMode.OnClick(func(ctx *application.Context) {
    isChecked := darkMode.Checked()
    // Toggle dark mode
})
```

### AddRadio()

Добавляет пункт меню с переключателем (взаимоисключающая группа).

```go
func (m *Menu) AddRadio(label string, checked bool) *MenuItem
```

**Параметры:**

- `label` — метка переключателя
- `checked` — исходное выбранное состояние

**Пример:**

```go
// Create radio group for view modes
viewMenu := menu.AddSubmenu("View")
listView := viewMenu.AddRadio("List View", true)
gridView := viewMenu.AddRadio("Grid View", false)
treeView := viewMenu.AddRadio("Tree View", false)

listView.OnClick(func(ctx *application.Context) {
    setViewMode("list")
})
gridView.OnClick(func(ctx *application.Context) {
    setViewMode("grid")
})
```

### Update()

Обновляет меню, отражая все изменения, внесённые в его пункты.

```go
func (m *Menu) Update()
```

**Пример:**

```go
item.SetEnabled(false)
menu.Update()  // Must call to apply changes
```

**Важно:** после изменения свойств пунктов меню всегда вызывайте `Update()`.

## Методы пунктов меню

### OnClick()

Регистрирует обработчик щелчка по пункту меню.

```go
func (mi *MenuItem) OnClick(callback func(ctx *application.Context)) *MenuItem
```

**Параметры:**

- `callback` — функция, вызываемая при щелчке по пункту

**Возвращает:** пункт меню (для цепочки вызовов)

**Пример:**

```go
item.OnClick(func(ctx *application.Context) {
    fmt.Println("Menu item clicked")
    app.Logger.Info("User clicked menu item")
})
```

### SetLabel()

Изменяет метку пункта меню.

```go
func (mi *MenuItem) SetLabel(label string) *MenuItem
```

**Пример:**

```go
item.SetLabel("Save As...")
menu.Update()
```

### SetEnabled()

Включает или отключает пункт меню.

```go
func (mi *MenuItem) SetEnabled(enabled bool) *MenuItem
```

**Пример:**

```go
// Disable save when no document is open
saveItem.SetEnabled(hasOpenDocument)
menu.Update()
```

**Распространённый шаблон:**

```go
// Update menu state based on application state
func updateMenuState() {
    saveItem.SetEnabled(hasUnsavedChanges)
    undoItem.SetEnabled(canUndo)
    redoItem.SetEnabled(canRedo)
    menu.Update()
}
```

### SetChecked()

Задаёт состояние выбора для пунктов меню с флажками или переключателями.

```go
func (mi *MenuItem) SetChecked(checked bool) *MenuItem
```

**Пример:**

```go
darkModeItem.SetChecked(isDarkModeEnabled)
menu.Update()
```

### Checked()

Возвращает текущее состояние выбора.

```go
func (mi *MenuItem) Checked() bool
```

**Пример:**

```go
if darkModeItem.Checked() {
    // Dark mode is enabled
}
```

### SetAccelerator()

Задаёт сочетание клавиш для пункта меню.

```go
func (mi *MenuItem) SetAccelerator(accelerator string) *MenuItem
```

**Параметры:**

- `accelerator` — сочетание клавиш (например, «Ctrl+S», «Cmd+Q»)

**Формат сочетания клавиш:**

- **Клавиши-модификаторы:** `Ctrl`, `Cmd`, `Alt`, `Shift`
- **Клавиши:** `A-Z`, `0-9`, `F1-F12`, `Enter`, `Backspace` и т. д.
- **Платформа:** используйте `Cmd` в macOS и `Ctrl` в Windows/Linux

**Пример:**

```go
saveItem.SetAccelerator("Ctrl+S")
quitItem.SetAccelerator("Ctrl+Q")
newItem.SetAccelerator("Ctrl+N")
```

**Пример с учётом платформы:**

```go
import "runtime"

var quitShortcut string
if runtime.GOOS == "darwin" {
    quitShortcut = "Cmd+Q"
} else {
    quitShortcut = "Ctrl+Q"
}
quitItem.SetAccelerator(quitShortcut)
```

### SetTooltip()

Задаёт всплывающую подсказку, которая появляется при наведении указателя на пункт меню.

```go
func (mi *MenuItem) SetTooltip(tooltip string) *MenuItem
```

**Пример:**

```go
item.SetTooltip("Opens a file from disk")
```

### SetHidden()

Показывает или скрывает пункт меню.

```go
func (mi *MenuItem) SetHidden(hidden bool) *MenuItem
```

**Пример:**

```go
// Hide debug menu in production
debugItem.SetHidden(!isDevelopment)
menu.Update()
```

## Меню приложения

### app.Menu.Set()

Задаёт главную строку меню приложения.

```go
func (mm *MenuManager) Set(menu *Menu)
```

**Пример:**

```go
menu := app.NewMenu()

// File menu
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("New").SetAccelerator("Ctrl+N").OnClick(newFile)
fileMenu.Add("Open").SetAccelerator("Ctrl+O").OnClick(openFile)
fileMenu.Add("Save").SetAccelerator("Ctrl+S").OnClick(saveFile)
fileMenu.AddSeparator()
fileMenu.Add("Exit").SetAccelerator("Ctrl+Q").OnClick(func(ctx *application.Context) {
    app.Quit()
})

// Edit menu
editMenu := menu.AddSubmenu("Edit")
editMenu.Add("Undo").SetAccelerator("Ctrl+Z").OnClick(undo)
editMenu.Add("Redo").SetAccelerator("Ctrl+Y").OnClick(redo)
editMenu.AddSeparator()
editMenu.Add("Cut").SetAccelerator("Ctrl+X").OnClick(cut)
editMenu.Add("Copy").SetAccelerator("Ctrl+C").OnClick(copy)
editMenu.Add("Paste").SetAccelerator("Ctrl+V").OnClick(paste)

app.Menu.Set(menu)
```

**Примечания для разных платформ:**

- **macOS:** меню отображается в верхней строке меню
- **Windows/Linux:** меню отображается в строке заголовка окна
- **macOS:** меню приложения с его названием добавляется автоматически

## Контекстные меню

### app.ContextMenu.New() / app.ContextMenu.Add()

Создайте `*ContextMenu` с помощью менеджера и зарегистрируйте его под определённым именем. `ContextMenuManager.Add` принимает `*ContextMenu`, а **не** `*Menu`; метода **`app.RegisterContextMenu`** также нет.

```go
func (cm *ContextMenuManager) New() *ContextMenu
func (cm *ContextMenuManager) Add(name string, menu *ContextMenu)
func (cm *ContextMenuManager) Get(name string) (*ContextMenu, bool)
func (cm *ContextMenuManager) Remove(name string)
```

Также функция `application.NewContextMenu(name string) *ContextMenu` уровня пакета позволяет создать и зарегистрировать контекстное меню за один шаг.

**Go:**

```go
// Build the context menu via the manager
contextMenu := app.ContextMenu.New()
contextMenu.Add("Cut").OnClick(cut)
contextMenu.Add("Copy").OnClick(copy)
contextMenu.Add("Paste").OnClick(paste)
contextMenu.AddSeparator()
contextMenu.Add("Select All").OnClick(selectAll)

// Register under a name; HTML opts into it via the CSS custom property below.
app.ContextMenu.Add("editor", contextMenu)
```

**HTML/CSS:**

Среда выполнения вызывает зарегистрированное контекстное меню, если у элемента, по которому щёлкнули правой кнопкой мыши, или у любого из его предков пользовательскому CSS-свойству `--custom-contextmenu` присвоено имя этого меню. Значение необязательного свойства `--custom-contextmenu-data` передаётся в функцию обратного вызова Go через `ctx.ContextMenuData()`. Чтобы отключить стандартное контекстное меню браузера, задайте `--default-contextmenu: hide` (либо `auto`/`show`).

```html
<!-- Trigger context menu on right-click -->
<div style="--custom-contextmenu: editor; --default-contextmenu: hide">
    Right-click here for context menu
</div>
```

Атрибута `data-wails-context-menu="..."` нет — он никогда не был подключён в среде выполнения.

**Динамические контекстные меню:**

```go
// Update context menu based on selection
func updateContextMenu() {
    contextMenu := app.ContextMenu.New()

    if hasSelection {
        contextMenu.Add("Cut").OnClick(cut)
        contextMenu.Add("Copy").OnClick(copy)
    }

    contextMenu.Add("Paste").SetEnabled(hasClipboardContent).OnClick(paste)

    app.ContextMenu.Add("editor", contextMenu)
}
```

## Меню в области уведомлений

### app.SystemTray.New()

Создаёт новый значок в области уведомлений.

```go
func (sm *SystemTrayManager) New() *SystemTray
```

**Пример:**

```go
tray := app.SystemTray.New()
```

### SetIcon()

Задаёт значок в области уведомлений.

```go
func (st *SystemTray) SetIcon(icon []byte) *SystemTray
```

**Пример:**

```go
iconData, _ := os.ReadFile("icon.png")
tray.SetIcon(iconData)
```

### SetMenu()

Задаёт меню для области уведомлений.

```go
func (st *SystemTray) SetMenu(menu *Menu) *SystemTray
```

**Пример:**

```go
trayMenu := app.NewMenu()
trayMenu.Add("Show Window").OnClick(func(ctx *application.Context) {
    window.Show()
    window.Focus()
})
trayMenu.Add("Settings").OnClick(openSettings)
trayMenu.AddSeparator()
trayMenu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})

tray.SetMenu(trayMenu)
```

### SetTooltip()

Задаёт всплывающую подсказку, отображаемую при наведении указателя на значок в области уведомлений. Ничего не возвращает.

```go
func (st *SystemTray) SetTooltip(tooltip string)
```

**Пример:**

```go
tray.SetTooltip("My Application - Running")
```

### OnClick()

Обрабатывает щелчок левой кнопкой мыши по значку в области уведомлений.

```go
func (st *SystemTray) OnClick(callback func()) *SystemTray
```

**Пример:**

```go
tray.OnClick(func() {
    if window.IsVisible() {
        window.Hide()
    } else {
        window.Show()
        window.Focus()
    }
})
```

## Полные примеры

### Стандартное меню приложения

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
)

func createMenu(app *application.App) *application.Menu {
    menu := app.NewMenu()

    // File menu
    fileMenu := menu.AddSubmenu("File")
    fileMenu.Add("New").
        SetAccelerator("Ctrl+N").
        OnClick(func(ctx *application.Context) {
            // Create new document
        })
    fileMenu.Add("Open").
        SetAccelerator("Ctrl+O").
        OnClick(func(ctx *application.Context) {
            // Open file dialog
        })
    fileMenu.Add("Save").
        SetAccelerator("Ctrl+S").
        OnClick(func(ctx *application.Context) {
            // Save document
        })
    fileMenu.AddSeparator()
    fileMenu.Add("Exit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })

    // Edit menu
    editMenu := menu.AddSubmenu("Edit")
    editMenu.Add("Undo").SetAccelerator("Ctrl+Z")
    editMenu.Add("Redo").SetAccelerator("Ctrl+Y")
    editMenu.AddSeparator()
    editMenu.Add("Cut").SetAccelerator("Ctrl+X")
    editMenu.Add("Copy").SetAccelerator("Ctrl+C")
    editMenu.Add("Paste").SetAccelerator("Ctrl+V")

    // View menu
    viewMenu := menu.AddSubmenu("View")
    darkMode := viewMenu.AddCheckbox("Dark Mode", false)
    darkMode.OnClick(func(ctx *application.Context) {
        // Toggle dark mode
        isChecked := darkMode.Checked()
        app.Logger.Info("Dark mode", "enabled", isChecked)
    })
    viewMenu.AddSeparator()
    viewMenu.AddRadio("List View", true)
    viewMenu.AddRadio("Grid View", false)
    viewMenu.AddRadio("Detail View", false)

    // Help menu
    helpMenu := menu.AddSubmenu("Help")
    helpMenu.Add("Documentation").OnClick(func(ctx *application.Context) {
        // Open docs
    })
    helpMenu.Add("About").OnClick(func(ctx *application.Context) {
        // Show about dialog
    })

    return menu
}

func main() {
    app := application.New(application.Options{
        Name: "Menu Demo",
    })

    menu := createMenu(app)
    app.Menu.Set(menu)

    window := app.Window.New()
    window.Show()

    app.Run()
}
```

### Приложение со значком в области уведомлений

```go
func setupSystemTray(app *application.App, window application.Window) {
    // Create system tray
    tray := app.SystemTray.New()

    // Set icon
    iconData, _ := os.ReadFile("icon.png")
    tray.SetIcon(iconData)
    tray.SetTooltip("My App - Running")

    // Handle left-click on tray icon
    tray.OnClick(func() {
        if window.IsVisible() {
            window.Hide()
        } else {
            window.Show()
            window.Focus()
        }
    })

    // Create tray menu
    trayMenu := app.NewMenu()

    showItem := trayMenu.Add("Show Window")
    showItem.OnClick(func(ctx *application.Context) {
        window.Show()
        window.Focus()
    })

    trayMenu.AddSeparator()

    trayMenu.Add("Settings").OnClick(func(ctx *application.Context) {
        // Open settings window
    })

    trayMenu.Add("About").OnClick(func(ctx *application.Context) {
        // Show about dialog
    })

    trayMenu.AddSeparator()

    trayMenu.Add("Quit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })

    tray.SetMenu(trayMenu)
}
```

### Динамическое обновление меню

```go
type Editor struct {
    app         *application.App
    menu        *application.Menu
    undoItem    *application.MenuItem
    redoItem    *application.MenuItem
    saveItem    *application.MenuItem
    undoStack   []string
    redoStack   []string
    hasChanges  bool
}

func (e *Editor) createMenu() {
    e.menu = e.app.NewMenu()

    fileMenu := e.menu.AddSubmenu("File")
    e.saveItem = fileMenu.Add("Save").SetAccelerator("Ctrl+S")
    e.saveItem.OnClick(func(ctx *application.Context) {
        e.save()
    })

    editMenu := e.menu.AddSubmenu("Edit")
    e.undoItem = editMenu.Add("Undo").SetAccelerator("Ctrl+Z")
    e.undoItem.OnClick(func(ctx *application.Context) {
        e.undo()
    })

    e.redoItem = editMenu.Add("Redo").SetAccelerator("Ctrl+Y")
    e.redoItem.OnClick(func(ctx *application.Context) {
        e.redo()
    })

    e.updateMenuState()
    e.app.Menu.Set(e.menu)
}

func (e *Editor) updateMenuState() {
    // Update menu items based on current state
    e.saveItem.SetEnabled(e.hasChanges)
    e.undoItem.SetEnabled(len(e.undoStack) > 0)
    e.redoItem.SetEnabled(len(e.redoStack) > 0)
    e.menu.Update()
}

func (e *Editor) onChange() {
    e.hasChanges = true
    e.updateMenuState()
}

func (e *Editor) save() {
    // Save logic
    e.hasChanges = false
    e.updateMenuState()
}
```

## Рекомендации

### ✅ Рекомендуется

- **Используйте стандартные сочетания клавиш** — следуйте соглашениям платформы (например, Ctrl+C для копирования)
- **Вызывайте Update() после изменений** — иначе изменения не отобразятся в меню
- **Группируйте связанные пункты** — используйте разделители для упорядочивания пунктов меню
- **Отключайте недоступные действия** — не скрывайте их, а отключайте с помощью SetEnabled(false)
- **Используйте понятные названия** — формулируйте их кратко и информативно
- **Следуйте соглашениям платформы** — учитывайте различия между структурой меню в macOS и Windows/Linux

### ❌ Не рекомендуется

- **Не забывайте вызывать Update()** — это самая распространённая ошибка
- **Не создавайте слишком глубокую вложенность** — используйте не более 2-3 уровней меню
- **Не используйте неоднозначные названия** — «Обработать» вместо «Обработать документ»
- **Не усложняйте** — меню должны быть простыми и целенаправленными
- **Не смешивайте метафоры** — используйте единообразные названия и структуру

## Примечания для отдельных платформ

### macOS

- Меню приложения с его названием добавляется автоматически
- Для сочетаний клавиш используйте `Cmd` вместо `Ctrl`
- Пункты «О приложении», «Настройки» и «Выход» по умолчанию находятся в меню приложения

### Windows/Linux

- Меню приложения не создаётся автоматически
- Используйте `Ctrl` для сочетаний клавиш
- Пункт «Выход» обычно находится в меню «Файл»
