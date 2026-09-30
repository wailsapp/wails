---
title: "Справочник по меню"
description: "Полный справочник по типам, свойствам и методам элементов меню"
slug: "features/menus/reference"
sourcePath: "features/menus/reference.md"
---

## Справочник по меню

Полный справочник по типам, свойствам и динамическому поведению элементов меню. Создавайте профессиональные, отзывчивые меню с флажками, группами переключателей, разделителями и динамическими обновлениями.

## Типы элементов меню

### Обычные элементы меню

Самый распространённый тип: отображает текст и запускает действие:

```go
menuItem := menu.Add("Click Me")
menuItem.OnClick(func(ctx *application.Context) {
    fmt.Println("Menu item clicked!")
})
```

**Используйте для:** команд, действий, открытия окон

### Флажки

Элементы меню, которые можно переключать между отмеченным и неотмеченным состояниями:

```go
checkbox := menu.AddCheckbox("Enable Feature", true)  // true = initially checked
checkbox.OnClick(func(ctx *application.Context) {
    isChecked := ctx.ClickedMenuItem().Checked()
    fmt.Printf("Feature is now: %v\n", isChecked)
})
```

**Используйте для:** логических параметров, переключателей функций, параметров представления

**Важно:** при нажатии отмеченное состояние переключается автоматически.

### Группы переключателей

Взаимоисключающие варианты — выбрать можно только один:

```go
menu.AddRadio("Small", true)   // true = initially selected
menu.AddRadio("Medium", false)
menu.AddRadio("Large", false)
```

**Используйте для:** взаимоисключающих вариантов (размер, тема, режим)

**Как работает группировка:**

- Соседние элементы-переключатели автоматически образуют группу
- При выборе одного элемента остальные элементы группы перестают быть выбранными
- Разделяйте группы разделителем или обычным элементом

**Пример с несколькими группами:**

```go
// Group 1: Size
menu.AddRadio("Small", true)
menu.AddRadio("Medium", false)
menu.AddRadio("Large", false)

menu.AddSeparator()

// Group 2: Theme
menu.AddRadio("Light", true)
menu.AddRadio("Dark", false)
```

### Подменю

Вложенные структуры меню для упорядочивания элементов:

```go
submenu := menu.AddSubmenu("More Options")
submenu.Add("Submenu Item 1").OnClick(func(ctx *application.Context) {
    // Handle click
})
submenu.Add("Submenu Item 2")
```

**Используйте для:** группировки связанных элементов и уменьшения загромождённости

**Ограничение вложенности:** большинство платформ поддерживает 2-3 уровней. Избегайте более глубокой вложенности.

### Разделители

Визуальные разделители между элементами меню:

```go
menu.Add("Item 1")
menu.AddSeparator()
menu.Add("Item 2")
```

**Используйте для:** визуальной группировки связанных элементов

**Рекомендация:** не размещайте разделители в начале или конце меню.

## Свойства элементов меню

### Метка

Текст, отображаемый для элемента меню:

```go
menuItem := menu.Add("Initial Label")
menuItem.SetLabel("New Label")

// Get current label
label := menuItem.Label()
```

**Динамические метки:**

```go
updateMenuItem := menu.Add("Check for Updates")
updateMenuItem.OnClick(func(ctx *application.Context) {
    updateMenuItem.SetLabel("Checking...")
    menu.Update()  // Important on Windows!
    
    // Perform update check
    checkForUpdates()
    
    updateMenuItem.SetLabel("Check for Updates")
    menu.Update()
})
```

### Состояние доступности

Управляйте тем, можно ли взаимодействовать с элементом меню:

```go
menuItem := menu.Add("Save")
menuItem.SetEnabled(false)  // Greyed out, can't click

// Enable it later
menuItem.SetEnabled(true)
menu.Update()  // Important: Call this after changing enabled state!

// Check current state
isEnabled := menuItem.Enabled()
```

@note{type="caution" title="Поведение меню в Windows"}
В Windows при изменении состояния меню необходимо создавать заново. **После включения или отключения элементов меню всегда вызывайте `menu.Update()`**, особенно если элемент был создан отключённым.

**Причина:** при обновлении меню Windows полностью перестраиваются. Если не вызвать `Update()`, обработчики нажатий не будут срабатывать должным образом.

@end

**Пример: динамическое включение и отключение**

```go
var hasSelection bool

cutMenuItem := menu.Add("Cut")
cutMenuItem.SetEnabled(false)  // Initially disabled

copyMenuItem := menu.Add("Copy")
copyMenuItem.SetEnabled(false)

// When selection changes
func onSelectionChanged(selected bool) {
    hasSelection = selected
    cutMenuItem.SetEnabled(hasSelection)
    copyMenuItem.SetEnabled(hasSelection)
    menu.Update()  // Critical on Windows!
}
```

**Распространённый шаблон: включение при выполнении условия**

```go
saveMenuItem := menu.Add("Save")

func updateSaveMenuItem() {
    canSave := hasUnsavedChanges() && !isSaving()
    saveMenuItem.SetEnabled(canSave)
    menu.Update()
}

// Call whenever state changes
onDocumentChanged(func() {
    updateSaveMenuItem()
})
```

### Отмеченное состояние

Управляйте отмеченным состоянием флажков и переключателей или проверяйте его:

```go
checkbox := menu.AddCheckbox("Feature", false)
checkbox.SetChecked(true)
menu.Update()

// Query state
isChecked := checkbox.Checked()
```

**Автоматическое переключение:** при нажатии флажки переключаются автоматически. Вызывать `SetChecked()` в обработчике нажатия не требуется.

**Ручное управление:**

```go
checkbox := menu.AddCheckbox("Auto-save", false)

// Sync with external state
func syncAutoSave(enabled bool) {
    checkbox.SetChecked(enabled)
    menu.Update()
}
```

### Акселераторы (сочетания клавиш)

Добавляйте сочетания клавиш к элементам меню:

```go
saveMenuItem := menu.Add("Save")
saveMenuItem.SetAccelerator("CmdOrCtrl+S")

quitMenuItem := menu.Add("Quit")
quitMenuItem.SetAccelerator("CmdOrCtrl+Q")
```

**Формат акселератора:**

- `CmdOrCtrl` — Cmd в macOS, Ctrl в Windows/Linux
- `Shift`, `Alt`, `Option` — клавиши-модификаторы
- `A-Z`, `0-9` — буквенные и цифровые клавиши
- `F1-F12` — функциональные клавиши
- `Enter`, `Space`, `Backspace` и т. д. — специальные клавиши

**Примеры:**

```go
"CmdOrCtrl+S"           // Save
"CmdOrCtrl+Shift+S"     // Save As
"CmdOrCtrl+W"           // Close Window
"CmdOrCtrl+Q"           // Quit
"F5"                    // Refresh
"CmdOrCtrl+,"           // Preferences (macOS convention)
"Alt+F4"                // Close (Windows convention)
```

**Акселераторы для конкретных платформ:**

```go
if runtime.GOOS == "darwin" {
    prefsMenuItem.SetAccelerator("Cmd+,")
} else {
    prefsMenuItem.SetAccelerator("Ctrl+P")
}
```

### Всплывающая подсказка

Добавляйте к элементам меню текст, отображаемый при наведении указателя (поддержка зависит от платформы):

```go
menuItem := menu.Add("Advanced Options")
menuItem.SetTooltip("Configure advanced settings")
```

**Поддержка платформами:**

- **Windows:** ✅ поддерживается
- **macOS:** ❌ не поддерживается (всплывающие подсказки не являются стандартными для меню)
- **Linux:** ⚠️ зависит от среды рабочего стола

### Скрытое состояние

Скрывайте элементы меню, не удаляя их:

```go
debugMenuItem := menu.Add("Debug Mode")
debugMenuItem.SetHidden(true)  // Hidden

// Show in debug builds
if isDebugBuild {
    debugMenuItem.SetHidden(false)
    menu.Update()
}
```

**Используйте для:** параметров отладки, флагов функций, условно доступных функций

## Обработка событий

### Обработчик OnClick

Выполните код при выборе пункта меню:

```go
menuItem := menu.Add("Click Me")
menuItem.OnClick(func(ctx *application.Context) {
    // Handle click
    fmt.Println("Clicked!")
})
```

**Контекст предоставляет:**

- `ctx.ClickedMenuItem()` — выбранный пункт меню
- Контекст окна (если событие поступило из меню окна)
- Контекст приложения

**Пример: доступ к пункту меню в обработчике**

```go
checkbox := menu.AddCheckbox("Feature", false)
checkbox.OnClick(func(ctx *application.Context) {
    item := ctx.ClickedMenuItem()
    isChecked := item.Checked()
    fmt.Printf("Feature is now: %v\n", isChecked)
})
```

### Несколько обработчиков

Можно задать несколько обработчиков (используется последний):

```go
menuItem := menu.Add("Action")
menuItem.OnClick(func(ctx *application.Context) {
    fmt.Println("First handler")
})

// This replaces the first handler
menuItem.OnClick(func(ctx *application.Context) {
    fmt.Println("Second handler - this one runs")
})
```

**Рекомендация:** задайте обработчик один раз и при необходимости используйте в нём условную логику.

## Динамические меню

### Обновление пунктов меню

**Главное правило:** после изменения состояния меню всегда вызывайте `menu.Update()`.

```go
// ✅ Correct
menuItem.SetEnabled(true)
menu.Update()

// ❌ Wrong (especially on Windows)
menuItem.SetEnabled(true)
// Forgot to call Update() - click handlers may not work!
```

**Почему это важно:**

- **Windows:** при обновлении меню создаются заново
- **macOS/Linux:** это менее критично, но всё же рекомендуется
- **Обработчики нажатий:** без Update() не будут срабатывать должным образом

### Пересоздание меню

При существенных изменениях пересоздайте всё меню:

```go
func rebuildFileMenu() {
    menu := app.Menu.New()
    
    menu.Add("New").OnClick(handleNew)
    menu.Add("Open").OnClick(handleOpen)
    
    if hasRecentFiles() {
        recentMenu := menu.AddSubmenu("Open Recent")
        for _, file := range getRecentFiles() {
            recentMenu.Add(file).OnClick(func(ctx *application.Context) {
                openFile(file)
            })
        }
    }
    
    menu.AddSeparator()
    menu.Add("Quit").OnClick(handleQuit)
    
    // Set the new menu
    window.SetMenu(menu)
}
```

**Когда следует пересоздавать меню:**

- Изменяется список последних файлов
- Изменяются меню плагинов
- Происходят существенные переходы между состояниями

**Когда следует обновлять меню:**

- Включение и отключение пунктов
- Изменение названий
- Переключение флажков

### Контекстно-зависимые меню

Настраивайте меню в зависимости от состояния приложения:

```go
func updateEditMenu() {
    cutMenuItem.SetEnabled(hasSelection())
    copyMenuItem.SetEnabled(hasSelection())
    pasteMenuItem.SetEnabled(hasClipboardContent())
    undoMenuItem.SetEnabled(canUndo())
    redoMenuItem.SetEnabled(canRedo())
    menu.Update()
}

// Call whenever state changes
onSelectionChanged(updateEditMenu)
onClipboardChanged(updateEditMenu)
onUndoStackChanged(updateEditMenu)
```

## Различия между платформами

### Расположение строки меню

| Платформа | Расположение | Примечания |
| --- | --- | --- |
| **macOS** | В верхней части экрана | Глобальная строка меню |
| **Windows** | В верхней части окна | Отдельное меню для каждого окна |
| **Linux** | В верхней части окна | Отдельное для каждого окна (обычно) |

### Стандартные меню

**macOS:**

- Есть меню «Приложение» (с названием приложения)
- Пункт «Настройки» находится в меню приложения
- Пункт «Завершить» находится в меню приложения

**Windows/Linux:**

- Меню приложения отсутствует
- Пункт «Настройки» находится в меню «Правка» или «Инструменты»
- Пункт «Выход» находится в меню «Файл»

**Пример: структура, соответствующая платформе**

```go
menu := app.Menu.New()

// macOS gets Application menu
if runtime.GOOS == "darwin" {
    menu.AddRole(application.AppMenu)
}

// File menu
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("New")
fileMenu.Add("Open")

// Preferences location varies
if runtime.GOOS == "darwin" {
    // On macOS, preferences are in Application menu (added by AppMenu role)
} else {
    // On Windows/Linux, add to Edit or Tools menu
    editMenu := menu.AddSubmenu("Edit")
    editMenu.Add("Preferences")
}
```

### Соглашения о сочетаниях клавиш

**macOS:**

- `Cmd+` для большинства сочетаний клавиш
- `Cmd+,` для пункта «Настройки»
- `Cmd+Q` для пункта «Завершить»

**Windows:**

- `Ctrl+` для большинства сочетаний клавиш
- `Ctrl+P` или `Ctrl+,` для пункта «Настройки»
- `Alt+F4` для пункта «Выход» (или `Ctrl+Q`)

**Linux:**

- Обычно используются соглашения Windows
- Среда рабочего стола может переопределить сочетания клавиш

## Рекомендации

### ✅ Рекомендуется

- После изменения состояния меню **вызывайте menu.Update()** (особенно в Windows)
- Для взаимоисключающих вариантов **используйте группы переключателей**
- Для включаемых и отключаемых функций **используйте флажки**
- **Добавляйте сочетания клавиш** для часто используемых действий
- **Группируйте связанные пункты** с помощью разделителей
- **Тестируйте на всех платформах** — поведение различается

### ❌ Не делайте так

- **Не забывайте вызывать menu.Update()** — обработчики щелчков не будут работать должным образом
- **Не создавайте слишком глубокую вложенность** — не более 2-3 уровней
- **Не ставьте разделители в начале или конце** — это выглядит непрофессионально
- **Не используйте всплывающие подсказки в macOS** — они не поддерживаются
- **Не задавайте платформенные сочетания клавиш напрямую** — используйте `CmdOrCtrl`

## Устранение неполадок

### Пункты меню не реагируют

**Признак:** обработчики щелчков не срабатывают

**Причина:** после включения пункта не был вызван `menu.Update()`

**Решение:**

```go
menuItem.SetEnabled(true)
menu.Update()  // Add this!
```

### Пункты меню отображаются серым цветом

**Признак:** пункты меню невозможно выбрать

**Причина:** пункты отключены

**Решение:**

```go
menuItem.SetEnabled(true)
menu.Update()
```

### Сочетания клавиш не работают

**Признак:** сочетания клавиш не активируют пункты меню

**Причины:**

1. Неверный формат сочетания клавиш
2. Конфликт с системными сочетаниями клавиш
3. Окно не находится в фокусе

**Решение:**

```go
// Check format
menuItem.SetAccelerator("CmdOrCtrl+S")  // ✅ Correct
menuItem.SetAccelerator("Ctrl+S")       // ❌ Wrong (macOS uses Cmd)

// Avoid conflicts
// ❌ Cmd+H (Hide Window on macOS - system shortcut)
// ✅ Cmd+Shift+H (Custom shortcut)
```

## Дальнейшие шаги

- [Меню приложения](/features/menus/application/) — создание строк меню приложения
- [Контекстные меню](/features/menus/context/) — контекстные меню, открываемые щелчком правой кнопки мыши
- [Меню в области уведомлений](/features/menus/systray/) — меню в области уведомлений или строке меню
- [Шаблоны меню](/guides/menus/) — распространённые шаблоны меню и рекомендации

---

**Есть вопросы?** Задайте их в [Discord](https://discord.gg/JDdSxwjhGf) или ознакомьтесь с [примерами меню](https://github.com/wailsapp/wails/tree/master/v3/examples/menu).
