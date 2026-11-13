---
title: "Справочник по меню"
description: "Полный справочник по типам, свойствам и методам элементов меню"
slug: "features/menus/reference"
sourcePath: "features/menus/reference.md"
---

## Справочник по меню

Полный справочник по типам, свойствам и динамическому поведению элементов меню. Создавайте профессиональные, отзывчивые меню с флажками, группами переключателей, разделителями и динамическим обновлением.

## Типы элементов меню

### Обычные элементы меню

Самый распространённый тип: отображает текст и запускает действие:

```go
menuItem := menu.Add("Click Me")
menuItem.OnClick(func(ctx *application.Context) {
    fmt.Println("Menu item clicked!")
})
```

**Назначение:** команды, действия, открытие окон

### Флажки

Элементы меню, переключаемые между отмеченным и неотмеченным состояниями:

```go
checkbox := menu.AddCheckbox("Enable Feature", true)  // true = initially checked
checkbox.OnClick(func(ctx *application.Context) {
    isChecked := ctx.ClickedMenuItem().Checked()
    fmt.Printf("Feature is now: %v\n", isChecked)
})
```

**Назначение:** логические параметры, переключатели функций, параметры отображения

**Важно:** состояние флажка автоматически переключается при нажатии.

### Группы переключателей

Взаимоисключающие варианты: выбрать можно только один:

```go
menu.AddRadio("Small", true)   // true = initially selected
menu.AddRadio("Medium", false)
menu.AddRadio("Large", false)
```

**Назначение:** выбор взаимоисключающих вариантов (размер, тема, режим)

**Как работает группировка:**

- Смежные переключатели автоматически образуют группу
- При выборе одного переключателя остальные в группе сбрасываются
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

Вложенные структуры меню для упорядочивания:

```go
submenu := menu.AddSubmenu("More Options")
submenu.Add("Submenu Item 1").OnClick(func(ctx *application.Context) {
    // Handle click
})
submenu.Add("Submenu Item 2")
```

**Назначение:** группировка связанных элементов, уменьшение загромождённости

**Ограничение вложенности:** большинство платформ поддерживает 2-3 уровней. Избегайте более глубокой вложенности.

### Разделители

Визуальные разделители между элементами меню:

```go
menu.Add("Item 1")
menu.AddSeparator()
menu.Add("Item 2")
```

**Назначение:** визуальная группировка связанных элементов

**Рекомендация:** не ставьте разделители в начале или конце меню.

### Заголовки разделов

Нативные неинтерактивные заголовки для групп элементов меню в **macOS 14+**:

```go
menu.AddSectionHeader("Build server")
menu.Add("Start build").OnClick(startBuild)
menu.Add("View logs").OnClick(viewLogs)
```

`AddSectionHeader` использует стандартное оформление заголовков разделов AppKit. Метод не возвращает элемент меню и не назначает действие. В старых версиях macOS, а также в сборках для Windows, Linux, мобильных устройств и серверов вызов ничего не делает и не добавляет строку.

Используйте его в меню системного трея, контекстных меню или меню приложения для обозначения связанных команд. Для разделения групп переключателей на всех платформах используйте `AddSeparator()`, поскольку отсутствующие заголовки не разделяют такие группы. Сигнатуру метода см. в разделе [API меню](/reference/menu#addsectionheader), а пример заголовка с именем хоста — в примере `systray-menu`.

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

Управляйте возможностью взаимодействия с элементом меню:

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
В Windows при изменении состояния меню необходимо пересоздавать. **После включения или отключения элементов меню всегда вызывайте `menu.Update()`**, особенно если элемент был создан отключённым.

**Причина:** при обновлении меню Windows полностью создаются заново. Если не вызвать `Update()`, обработчики нажатий не будут срабатывать должным образом.

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

### Состояние выбора

Управляйте состоянием выбора флажков и переключателей или запрашивайте его:

```go
checkbox := menu.AddCheckbox("Feature", false)
checkbox.SetChecked(true)
menu.Update()

// Query state
isChecked := checkbox.Checked()
```

**Автоматическое переключение:** флажки автоматически переключаются при нажатии. Вызывать `SetChecked()` в обработчике нажатия не требуется.

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

**Платформозависимые акселераторы:**

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

**Поддержка платформ:**

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

**Назначение:** параметры отладки, флаги функций, условно доступные функции

## Обработка событий

### Обработчик OnClick

Выполняйте код при нажатии элемента меню:

```go
menuItem := menu.Add("Click Me")
menuItem.OnClick(func(ctx *application.Context) {
    // Handle click
    fmt.Println("Clicked!")
})
```

**Контекст предоставляет:**

- `ctx.ClickedMenuItem()` — элемент меню, по которому нажали
- Контекст окна (если событие поступило из меню окна)
- Контекст приложения

**Пример: доступ к элементу меню в обработчике**

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

**Рекомендация:** назначьте обработчик один раз и при необходимости используйте внутри него условную логику.

## Динамические меню

### Обновление элементов меню

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
- **macOS/Linux:** менее критично, но всё же рекомендуется
- **Обработчики щелчков:** без Update() не будут срабатывать должным образом

### Пересоздание меню

При значительных изменениях пересоздайте всё меню:

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

**Когда пересоздавать:**

- Изменяется список недавних файлов
- Изменяются меню плагинов
- Происходят значительные переходы между состояниями

**Когда обновлять:**

- Включение и отключение пунктов
- Изменение подписей
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

- Есть меню «Приложение» (с именем приложения)
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
- `Cmd+,` для настроек
- `Cmd+Q` для завершения работы

**Windows:**

- `Ctrl+` для большинства сочетаний клавиш
- `Ctrl+P` или `Ctrl+,` для настроек
- `Alt+F4` для выхода (либо `Ctrl+Q`)

**Linux:**

- Обычно следует соглашениям Windows
- Среда рабочего стола может переопределять эти сочетания

## Рекомендации

### ✅ Рекомендуется

- После изменения состояния меню **вызывайте menu.Update()** (особенно в Windows)
- Для взаимоисключающих вариантов **используйте группы переключателей**
- Для включаемых и отключаемых функций **используйте флажки**
- Для часто используемых действий **добавляйте сочетания клавиш**
- **Группируйте связанные пункты** с помощью разделителей
- **Тестируйте на всех платформах** — поведение различается

### ❌ Не рекомендуется

- **Не забывайте вызывать menu.Update()** — обработчики щелчков не будут работать должным образом
- **Не создавайте слишком глубокую вложенность** — не более 2-3 уровней
- **Не размещайте разделители в начале или конце** — это выглядит непрофессионально
- **Не используйте всплывающие подсказки в macOS** — они не поддерживаются
- **Не задавайте сочетания клавиш для конкретной платформы напрямую** — используйте `CmdOrCtrl`

## Устранение неполадок

### Пункты меню не реагируют

**Симптом:** обработчики щелчков не срабатывают

**Причина:** после включения пункта не был вызван `menu.Update()`

**Решение:**

```go
menuItem.SetEnabled(true)
menu.Update()  // Add this!
```

### Пункты меню отображаются серым цветом

**Симптом:** пункты меню невозможно выбрать

**Причина:** пункты отключены

**Решение:**

```go
menuItem.SetEnabled(true)
menu.Update()
```

### Сочетания клавиш не работают

**Симптом:** сочетания клавиш не активируют пункты меню

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

## Дальнейшие действия

- [Меню приложения](/features/menus/application/) — создание строк меню приложения
- [Контекстные меню](/features/menus/context/) — контекстные меню, открываемые щелчком правой кнопки мыши
- [Меню области уведомлений](/features/menus/systray/) — меню области уведомлений или строки меню
- [Шаблоны меню](/guides/menus/) — распространённые шаблоны меню и рекомендации

---

**Остались вопросы?** Задайте их в [Discord](https://discord.gg/JDdSxwjhGf) или ознакомьтесь с [примерами меню](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/menu).
