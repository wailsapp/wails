---
title: "Контекстные меню"
description: "Создание контекстных меню, открываемых правой кнопкой мыши, для вашего приложения"
slug: "features/menus/context"
sourcePath: "features/menus/context.md"
---

## Проблема

Пользователи ожидают, что по щелчку правой кнопкой мыши будут открываться меню с действиями, соответствующими контексту. Для разных элементов нужны разные меню:

- **Текст**: вырезать, копировать, вставить
- **Изображения**: сохранить, копировать, открыть
- **Пользовательские элементы**: действия, относящиеся к приложению

При создании контекстных меню вручную приходится обрабатывать события мыши, позиционирование и различия между платформами.

## Решение Wails

Wails предоставляет **декларативные контекстные меню** на основе свойств CSS. Связывайте меню с HTML-элементами, передавайте данные и обрабатывайте щелчки — всё с использованием нативного поведения платформы.

![Пользовательское контекстное меню Wails, отображаемое поверх webview в macOS](/assets/screenshots/context-menu-macos.png)

Меню является нативным для платформы, а открывший его элемент остаётся частью вашего webview. На этом снимке экрана из macOS используется регистрация пользовательского контекстного меню из приведённого ниже примера.

## Быстрый старт

**Код Go:**

```go
// Create context menu
contextMenu := app.ContextMenu.New()
contextMenu.Add("Cut").OnClick(handleCut)
contextMenu.Add("Copy").OnClick(handleCopy)
contextMenu.Add("Paste").OnClick(handlePaste)

// Register with ID
app.ContextMenu.Add("editor-menu", contextMenu)
```

**HTML:**

```html
<textarea style="--custom-contextmenu: editor-menu">
    Right-click me!
</textarea>
```

**Вот и всё!** Щелчок правой кнопкой мыши по текстовой области открывает ваше пользовательское меню.

## Создание контекстных меню

### Базовое контекстное меню

```go
// Create menu
contextMenu := app.ContextMenu.New()

// Add items
contextMenu.Add("Cut").SetAccelerator("CmdOrCtrl+X").OnClick(func(ctx *application.Context) {
    // Handle cut
})

contextMenu.Add("Copy").SetAccelerator("CmdOrCtrl+C").OnClick(func(ctx *application.Context) {
    // Handle copy
})

contextMenu.Add("Paste").SetAccelerator("CmdOrCtrl+V").OnClick(func(ctx *application.Context) {
    // Handle paste
})

// Register with unique ID
app.ContextMenu.Add("text-menu", contextMenu)
```

**Идентификатор меню:** должен быть уникальным. Он используется для связывания меню с HTML-элементами.

### Меню с подменю

```go
contextMenu := app.ContextMenu.New()

// Add regular items
contextMenu.Add("Open").OnClick(handleOpen)
contextMenu.Add("Delete").OnClick(handleDelete)

contextMenu.AddSeparator()

// Add submenu
exportMenu := contextMenu.AddSubmenu("Export As")
exportMenu.Add("PNG").OnClick(exportPNG)
exportMenu.Add("JPEG").OnClick(exportJPEG)
exportMenu.Add("SVG").OnClick(exportSVG)

app.ContextMenu.Add("image-menu", contextMenu)
```

### Меню с флажками и группами переключателей

```go
contextMenu := app.ContextMenu.New()

// Checkbox
contextMenu.AddCheckbox("Show Grid", true).OnClick(func(ctx *application.Context) {
    showGrid := ctx.ClickedMenuItem().Checked()
    // Toggle grid
})

contextMenu.AddSeparator()

// Radio group
contextMenu.AddRadio("Small", false).OnClick(handleSize)
contextMenu.AddRadio("Medium", true).OnClick(handleSize)
contextMenu.AddRadio("Large", false).OnClick(handleSize)

app.ContextMenu.Add("view-menu", contextMenu)
```

Описание **всех типов пунктов меню** см. в разделе [«Справочник по меню»](/features/menus/reference/).

## Связывание с HTML-элементами

Используйте пользовательские свойства CSS, чтобы привязать контекстные меню:

### Базовое связывание

```html
<div style="--custom-contextmenu: menu-id">
    Right-click me!
</div>
```

**Свойство CSS:** `--custom-contextmenu: <menu-id>`

### С контекстными данными

Передавайте данные из HTML в Go:

```html
<div style="--custom-contextmenu: file-menu; --custom-contextmenu-data: file-123">
    Right-click this file
</div>
```

**Обработчик Go:**

```go
contextMenu := app.ContextMenu.New()
contextMenu.Add("Open").OnClick(func(ctx *application.Context) {
    fileID := ctx.ContextMenuData()  // "file-123"
    openFile(fileID)
})

app.ContextMenu.Add("file-menu", contextMenu)
```

**Свойства CSS:**

- `--custom-contextmenu: <menu-id>` — какое меню показывать
- `--custom-contextmenu-data: <data>` — данные, передаваемые обработчикам

### Динамические данные

Формируйте данные динамически в JavaScript:

```html
<div id="file-item" style="--custom-contextmenu: file-menu">
    File.txt
</div>

<script>
// Set data dynamically
const fileItem = document.getElementById('file-item')
fileItem.style.setProperty('--custom-contextmenu-data', 'file-' + fileId)
</script>
```

### Одно меню для нескольких элементов

```html
<div class="file-item" style="--custom-contextmenu: file-menu; --custom-contextmenu-data: file-1">
    Document.pdf
</div>

<div class="file-item" style="--custom-contextmenu: file-menu; --custom-contextmenu-data: file-2">
    Image.png
</div>

<div class="file-item" style="--custom-contextmenu: file-menu; --custom-contextmenu-data: file-3">
    Video.mp4
</div>
```

**Одно меню, но разные данные для каждого элемента.**

## Контекстные данные

### Доступ к контекстным данным

```go
contextMenu.Add("Process").OnClick(func(ctx *application.Context) {
    data := ctx.ContextMenuData()  // Get data from HTML
    
    // Use the data
    processItem(data)
})
```

**Тип данных:** всегда `string`. Выполняйте синтаксический разбор по мере необходимости.

### Передача сложных данных

Для сложных данных используйте JSON:

```html
<div style="--custom-contextmenu: item-menu; --custom-contextmenu-data: {&quot;id&quot;:123,&quot;type&quot;:&quot;image&quot;}">
    Image.png
</div>
```

**Обработчик Go:**

```go
import "encoding/json"

type ItemData struct {
    ID   int    `json:"id"`
    Type string `json:"type"`
}

contextMenu.Add("Process").OnClick(func(ctx *application.Context) {
    dataStr := ctx.ContextMenuData()
    
    var data ItemData
    if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
        log.Printf("Invalid data: %v", err)
        return
    }
    
    processItem(data.ID, data.Type)
})
```

@note{type="caution" title="Безопасность"}
**Всегда проверяйте контекстные данные**, поступающие из фронтенда. Пользователи могут изменять свойства CSS, поэтому рассматривайте эти данные как недоверенные входные данные.

@end

### Пример проверки

```go
contextMenu.Add("Delete").OnClick(func(ctx *application.Context) {
    fileID := ctx.ContextMenuData()
    
    // Validate
    if !isValidFileID(fileID) {
        log.Printf("Invalid file ID: %s", fileID)
        return
    }
    
    // Check permissions
    if !canDeleteFile(fileID) {
        showError("Permission denied")
        return
    }
    
    // Safe to proceed
    deleteFile(fileID)
})
```

## Контекстное меню по умолчанию

WebView предоставляет встроенное контекстное меню для стандартных операций (копирования, вставки и инспектирования). Управляйте им с помощью `--default-contextmenu`:

### Скрытие меню по умолчанию

```html
<div style="--default-contextmenu: hide">
    No default menu here
</div>
```

**Пример использования:** пользовательские элементы интерфейса, для которых меню по умолчанию не имеет смысла.

### Отображение меню по умолчанию

```html
<div style="--default-contextmenu: show">
    Default menu always shown
</div>
```

**Пример использования:** текстовые области, поля ввода и редактируемое содержимое.

### Автоматический (интеллектуальный) режим

```html
<div style="--default-contextmenu: auto">
    Smart context menu
</div>
```

**Поведение по умолчанию.** Меню по умолчанию отображается, когда:

- выделен текст
- указатель находится в текстовом поле ввода
- указатель находится в редактируемом содержимом (`contenteditable`)

В остальных случаях меню по умолчанию скрывается.

### Совместное использование пользовательского меню и меню по умолчанию

```html
<!-- Custom menu + default menu -->
<textarea style="--custom-contextmenu: editor-menu; --default-contextmenu: show">
    Both menus available
</textarea>
```

**Поведение:**

1. Сначала отображается пользовательское меню
2. Если пользовательское меню пусто или не найдено, отображается меню по умолчанию
3. Оба меню могут использоваться совместно (зависит от платформы)

## Динамические контекстные меню

Обновляйте меню в соответствии с состоянием приложения:

### Включение и отключение пунктов

```go
var cutMenuItem *application.MenuItem
var copyMenuItem *application.MenuItem

func createContextMenu() {
    contextMenu := app.ContextMenu.New()
    
    cutMenuItem = contextMenu.Add("Cut")
    cutMenuItem.SetEnabled(false)  // Initially disabled
    cutMenuItem.OnClick(handleCut)
    
    copyMenuItem = contextMenu.Add("Copy")
    copyMenuItem.SetEnabled(false)
    copyMenuItem.OnClick(handleCopy)
    
    app.ContextMenu.Add("editor-menu", contextMenu)
}

func onSelectionChanged(hasSelection bool) {
    cutMenuItem.SetEnabled(hasSelection)
    copyMenuItem.SetEnabled(hasSelection)
    contextMenu.Update()  // Important!
}
```

@note{type="caution" title="Всегда вызывайте Update()"}
После изменения состояния меню **вызовите `contextMenu.Update()`**. Это критически важно в Windows.

Подробнее см. в разделе [«Справочник по меню»](/features/menus/reference/#enabled-state).

@end

### Изменение подписей

```go
playMenuItem := contextMenu.Add("Play")

playMenuItem.OnClick(func(ctx *application.Context) {
    if isPlaying {
        playMenuItem.SetLabel("Pause")
    } else {
        playMenuItem.SetLabel("Play")
    }
    contextMenu.Update()
})
```

### Пересоздание меню

При значительных изменениях пересоздайте всё меню:

```go
func rebuildContextMenu(fileType string) {
    contextMenu := app.ContextMenu.New()
    
    // Common items
    contextMenu.Add("Open").OnClick(handleOpen)
    contextMenu.Add("Delete").OnClick(handleDelete)
    
    contextMenu.AddSeparator()
    
    // Type-specific items
    switch fileType {
    case "image":
        contextMenu.Add("Edit Image").OnClick(editImage)
        contextMenu.Add("Set as Wallpaper").OnClick(setWallpaper)
    case "video":
        contextMenu.Add("Play").OnClick(playVideo)
        contextMenu.Add("Extract Audio").OnClick(extractAudio)
    case "document":
        contextMenu.Add("Print").OnClick(printDocument)
        contextMenu.Add("Export PDF").OnClick(exportPDF)
    }
    
    app.ContextMenu.Add("file-menu", contextMenu)
}
```

## Поведение на разных платформах

Контекстные меню используют **нативный интерфейс платформы**:

@tabs{sync-key="platform"}
[macOS]
**Нативные контекстные меню macOS:**

- Системные анимации и переходы
- Щелчок правой кнопкой мыши = Control+Click (автоматически)
- Адаптация к системному оформлению (светлому или тёмному)
- Стандартные операции с текстом в меню по умолчанию
- Нативная прокрутка длинных меню

**Соглашения macOS:**

- Используйте регистр предложений для пунктов меню
- Добавляйте многоточие (...) к пунктам, открывающим диалоговые окна
- Распространённые сочетания клавиш: ⌘C (Копировать), ⌘V (Вставить)

[Windows]
**Нативные контекстные меню Windows:**

- Нативный стиль Windows
- Соответствие теме Windows
- Стандартные операции Windows в меню по умолчанию
- Поддержка сенсорного и перьевого ввода

**Соглашения Windows:**

- Начинайте каждое значимое слово в пунктах меню с прописной буквы
- Добавляйте многоточие (...) к пунктам, открывающим диалоговые окна
- Распространённые сочетания клавиш: Ctrl+C (Копировать), Ctrl+V (Вставить)

[Linux]
**Интеграция с окружением рабочего стола:**

- Адаптация к теме рабочего стола (GTK, Qt и т. д.)
- Поведение при щелчке правой кнопкой мыши определяется системными настройками
- Содержимое меню по умолчанию зависит от окружения
- Расположение соответствует соглашениям окружения рабочего стола

**Особенности Linux:**

- Тестируйте в целевых окружениях рабочего стола
- Поведение GTK и Qt различается
- Некоторые окружения рабочего стола настраивают контекстные меню по-своему

@end

## Полный пример

**Код Go:**

```go
package main

import (
    "encoding/json"
    "log"
    "github.com/wailsapp/wails/v3/pkg/application"
)

type FileData struct {
    ID   string `json:"id"`
    Type string `json:"type"`
    Name string `json:"name"`
}

func main() {
    app := application.New(application.Options{
        Name: "Context Menu Demo",
    })

    // Create file context menu
    fileMenu := createFileMenu(app)
    app.ContextMenu.Add("file-menu", fileMenu)

    // Create image context menu
    imageMenu := createImageMenu(app)
    app.ContextMenu.Add("image-menu", imageMenu)

    // Create text context menu
    textMenu := createTextMenu(app)
    app.ContextMenu.Add("text-menu", textMenu)

    app.Window.New()
    app.Run()
}

func createFileMenu(app *application.App) *application.ContextMenu {
    menu := app.ContextMenu.New()
    
    menu.Add("Open").OnClick(func(ctx *application.Context) {
        data := parseFileData(ctx.ContextMenuData())
        openFile(data.ID)
    })
    
    menu.Add("Rename").OnClick(func(ctx *application.Context) {
        data := parseFileData(ctx.ContextMenuData())
        renameFile(data.ID)
    })
    
    menu.AddSeparator()
    
    menu.Add("Delete").OnClick(func(ctx *application.Context) {
        data := parseFileData(ctx.ContextMenuData())
        deleteFile(data.ID)
    })
    
    return menu
}

func createImageMenu(app *application.App) *application.ContextMenu {
    menu := app.ContextMenu.New()
    
    menu.Add("View").OnClick(func(ctx *application.Context) {
        data := parseFileData(ctx.ContextMenuData())
        viewImage(data.ID)
    })
    
    menu.Add("Edit").OnClick(func(ctx *application.Context) {
        data := parseFileData(ctx.ContextMenuData())
        editImage(data.ID)
    })
    
    menu.AddSeparator()
    
    exportMenu := menu.AddSubmenu("Export As")
    exportMenu.Add("PNG").OnClick(exportPNG)
    exportMenu.Add("JPEG").OnClick(exportJPEG)
    exportMenu.Add("WebP").OnClick(exportWebP)
    
    return menu
}

func createTextMenu(app *application.App) *application.ContextMenu {
    menu := app.ContextMenu.New()
    
    menu.Add("Cut").SetAccelerator("CmdOrCtrl+X").OnClick(handleCut)
    menu.Add("Copy").SetAccelerator("CmdOrCtrl+C").OnClick(handleCopy)
    menu.Add("Paste").SetAccelerator("CmdOrCtrl+V").OnClick(handlePaste)
    
    return menu
}

func parseFileData(dataStr string) FileData {
    var data FileData
    if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
        log.Printf("Invalid file data: %v", err)
    }
    return data
}

// Handler implementations...
func openFile(id string) { /* ... */ }
func renameFile(id string) { /* ... */ }
func deleteFile(id string) { /* ... */ }
func viewImage(id string) { /* ... */ }
func editImage(id string) { /* ... */ }
func exportPNG(ctx *application.Context) { /* ... */ }
func exportJPEG(ctx *application.Context) { /* ... */ }
func exportWebP(ctx *application.Context) { /* ... */ }
func handleCut(ctx *application.Context) { /* ... */ }
func handleCopy(ctx *application.Context) { /* ... */ }
func handlePaste(ctx *application.Context) { /* ... */ }
```

**HTML:**

```html
<!DOCTYPE html>
<html>
<head>
    <style>
        .file-item {
            padding: 10px;
            margin: 5px;
            border: 1px solid #ccc;
            cursor: pointer;
        }
        
        .file-item:hover {
            background: #f0f0f0;
        }
        
        textarea {
            width: 100%;
            height: 200px;
        }
    </style>
</head>
<body>
    <h2>Files</h2>
    
    <!-- Regular file -->
    <div class="file-item" 
         style="--custom-contextmenu: file-menu; 
                --custom-contextmenu-data: {&quot;id&quot;:&quot;file-1&quot;,&quot;type&quot;:&quot;document&quot;,&quot;name&quot;:&quot;Report.pdf&quot;}">
        📄 Report.pdf
    </div>
    
    <!-- Image file -->
    <div class="file-item" 
         style="--custom-contextmenu: image-menu; 
                --custom-contextmenu-data: {&quot;id&quot;:&quot;file-2&quot;,&quot;type&quot;:&quot;image&quot;,&quot;name&quot;:&quot;Photo.jpg&quot;}">
        🖼️ Photo.jpg
    </div>
    
    <h2>Text Editor</h2>
    
    <!-- Text area with custom menu + default menu -->
    <textarea 
        style="--custom-contextmenu: text-menu; --default-contextmenu: show"
        placeholder="Type here, then right-click...">
    </textarea>
    
    <h2>No Context Menu</h2>
    
    <!-- Disable default menu -->
    <div style="--default-contextmenu: hide; padding: 20px; border: 1px solid #ccc;">
        Right-click here - no menu appears
    </div>
</body>
</html>
```

## Рекомендации

### ✅ Что следует делать

- **Не перегружайте меню** — включайте только действия, относящиеся к элементу
- **Проверяйте контекстные данные** — считайте их недоверенным вводом
- **Используйте понятные подписи** — «Удалить файл» вместо «Удалить»
- **Вызывайте menu.Update()** — после изменения состояния меню
- **Тестируйте на всех платформах** — поведение различается
- **Добавляйте сочетания клавиш** — для распространённых действий
- **Группируйте связанные пункты** — используйте разделители

### ❌ Чего не следует делать

- **Не доверяйте контекстным данным** — всегда проверяйте их
- **Не делайте меню слишком длинными** — не более 7-10 пунктов
- **Не забывайте вызывать menu.Update()** — иначе меню не будет работать должным образом
- **Не создавайте слишком глубокую вложенность** — не более 2 уровней
- **Не используйте жаргон** — подписи должны быть понятны пользователям
- **Не блокируйте обработчики** — они должны выполняться быстро

## Устранение неполадок

### Контекстное меню не появляется

**Возможные причины:**

1. Несоответствие идентификатора меню
2. Опечатка в CSS-свойстве
3. Среда выполнения не инициализирована

**Решение:**

```go
// Check menu is registered
app.ContextMenu.Add("my-menu", contextMenu)
```

```html
<!-- Check ID matches -->
<div style="--custom-contextmenu: my-menu">
```

### Контекстные данные не получены

**Возможные причины:**

1. CSS-свойство не задано
2. Данные содержат специальные символы

**Решение:**

```html
<!-- Escape quotes in JSON -->
<div style="--custom-contextmenu-data: {&quot;id&quot;:123}">
```

Или используйте JavaScript:

```javascript
element.style.setProperty('--custom-contextmenu-data', JSON.stringify(data))
```

### Пункты меню не реагируют

**Причина:** после включения не был вызван `menu.Update()`

**Решение:**

```go
menuItem.SetEnabled(true)
contextMenu.Update()  // Add this!
```

## Следующие шаги

@cards{cols="2"}
📖 Справочник по меню
Полный справочник по типам и свойствам элементов меню.

[Подробнее →](/features/menus/reference/)

---
☰ Меню приложения
Создавайте строки меню приложения.

[Подробнее →](/features/menus/application/)

---
★ Меню в области уведомлений
Добавьте интеграцию с областью уведомлений или строкой меню.

[Подробнее →](/features/menus/systray/)

---
📖 Шаблоны меню
Распространённые шаблоны меню и рекомендации по их использованию.

[Подробнее →](/guides/menus/)

@end

---

**Есть вопросы?** Задайте их в [Discord](https://discord.gg/JDdSxwjhGf) или ознакомьтесь с [примером контекстного меню](https://github.com/wailsapp/wails/tree/master/v3/examples/contextmenus).
