---
title: "Обзор диалоговых окон"
description: "Отображение нативных системных диалоговых окон в приложении"
slug: "features/dialogs/overview"
sourcePath: "features/dialogs/overview.md"
---

## Нативные диалоговые окна

Wails предоставляет **нативные системные диалоговые окна**, которые работают на всех платформах: окна сообщений (информация, предупреждение, ошибка, вопрос), диалоговые окна для работы с файлами (открытие, сохранение, выбор папки), а также настраиваемые диалоговые окна с нативными для платформы внешним видом и поведением.

![Диалоговое окно Wails с вопросом в macOS и кнопками «Отмена» и «Не сохранять»](/assets/screenshots/dialog-question-macos.png)

Один и тот же API отображает диалоговые окна в соответствии с соглашениями каждой поддерживаемой платформы. В этом примере для macOS показано диалоговое окно с вопросом, прикреплённое к окну Wails, включая кнопки по умолчанию и отмены.

## Быстрый старт

```go
// Information dialog
app.Dialog.Info().
    SetTitle("Success").
    SetMessage("File saved successfully!").
    Show()

// Question dialog with button callbacks
dialog := app.Dialog.Question().
    SetTitle("Confirm").
    SetMessage("Delete this file?")

deleteBtn := dialog.AddButton("Delete")
deleteBtn.OnClick(func() {
    deleteFile()
})

cancelBtn := dialog.AddButton("Cancel")
dialog.SetDefaultButton(cancelBtn)
dialog.SetCancelButton(cancelBtn)
dialog.Show()

// File open dialog
path, _ := app.Dialog.OpenFile().
    SetTitle("Select Image").
    AddFilter("Images", "*.png;*.jpg").
    PromptForSingleSelection()
```

**Вот и всё!** Нативные диалоговые окна с минимумом кода.

## Доступ к диалоговым окнам

Доступ к диалоговым окнам осуществляется через диспетчер `app.Dialog`:

```go
app.Dialog.Info()
app.Dialog.Question()
app.Dialog.Warning()
app.Dialog.Error()
app.Dialog.OpenFile()
app.Dialog.SaveFile()
```

## Типы диалоговых окон

### Информационное диалоговое окно

Отображайте простые сообщения:

```go
app.Dialog.Info().
    SetTitle("Welcome").
    SetMessage("Welcome to our application!").
    Show()
```

**Примеры использования:**

- Сообщения об успешном выполнении
- Информационные уведомления
- Подтверждения завершения

### Диалоговое окно предупреждения

Показывайте предупреждения:

```go
app.Dialog.Warning().
    SetTitle("Warning").
    SetMessage("This action cannot be undone.").
    Show()
```

**Примеры использования:**

- Некритические предупреждения
- Уведомления об устаревании
- Предостерегающие сообщения

### Диалоговое окно ошибки

Отображайте ошибки:

```go
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Failed to save file: " + err.Error()).
    Show()
```

**Примеры использования:**

- Сообщения об ошибках
- Уведомления о сбоях
- Обработка исключений

### Диалоговое окно с вопросом

Задавайте пользователям вопросы и обрабатывайте ответы с помощью обратных вызовов кнопок:

```go
dialog := app.Dialog.Question().
    SetTitle("Confirm Delete").
    SetMessage("Are you sure you want to delete this file?")

deleteBtn := dialog.AddButton("Delete")
deleteBtn.OnClick(func() {
    deleteFile()
})

cancelBtn := dialog.AddButton("Cancel")
dialog.SetDefaultButton(cancelBtn)
dialog.SetCancelButton(cancelBtn)
dialog.Show()
```

**Примеры использования:**

- Подтверждение действий
- Вопросы с ответами «Да» или «Нет»
- Выбор из нескольких вариантов

## Диалоговые окна для работы с файлами

### Диалоговое окно открытия файла

Выберите файлы для открытия:

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Image").
    AddFilter("Images", "*.png;*.jpg;*.gif").
    AddFilter("All Files", "*.*").
    PromptForSingleSelection()

if err == nil && path != "" {
    openFile(path)
}
```

**Множественный выбор:**

```go
paths, err := app.Dialog.OpenFile().
    SetTitle("Select Images").
    AddFilter("Images", "*.png;*.jpg").
    PromptForMultipleSelection()

if err == nil {
    for _, path := range paths {
        processFile(path)
    }
}
```

### Диалоговое окно сохранения файла

Выберите место сохранения:

```go
path, err := app.Dialog.SaveFile().
    SetFilename("document.txt").
    AddFilter("Text Files", "*.txt").
    AddFilter("All Files", "*.*").
    PromptForSingleSelection()

if err == nil && path != "" {
    saveFile(path)
}
```

### Диалоговое окно выбора папки

Выберите каталог с помощью диалогового окна открытия файла, включив выбор каталогов:

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Output Folder").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()

if err == nil && path != "" {
    exportToFolder(path)
}
```

## Параметры диалоговых окон

### Заголовок и сообщение

```go
dialog := app.Dialog.Info().
    SetTitle("Success").
    SetMessage("Operation completed successfully!")
```

### Кнопки

**Кнопка по умолчанию для простых диалоговых окон:**

В информационных диалоговых окнах, окнах предупреждений и ошибок отображается стандартная кнопка «ОК»:

```go
app.Dialog.Info().
    SetMessage("Done!").
    Show()
```

**Настраиваемые кнопки для диалоговых окон с вопросом:**

Используйте `AddButton()` для добавления кнопок. Метод возвращает `*Button`, который можно настроить с помощью обратных вызовов:

```go
dialog := app.Dialog.Question().
    SetMessage("Choose action")

save := dialog.AddButton("Save")
save.OnClick(func() {
    saveDocument()
})

dontSave := dialog.AddButton("Don't Save")
dontSave.OnClick(func() {
    discardChanges()
})

cancel := dialog.AddButton("Cancel")
// No callback needed - just dismisses dialog

dialog.SetDefaultButton(save)
dialog.SetCancelButton(cancel)
dialog.Show()
```

**Кнопки по умолчанию и отмены:**

Используйте `SetDefaultButton()`, чтобы указать кнопку, которая будет выделена и сработает при нажатии Enter. Используйте `SetCancelButton()`, чтобы указать кнопку, которая сработает при нажатии Escape.

```go
dialog := app.Dialog.Question().
    SetMessage("Delete file?")

deleteBtn := dialog.AddButton("Delete")
deleteBtn.OnClick(func() {
    performDelete()
})

cancelBtn := dialog.AddButton("Cancel")
dialog.SetDefaultButton(cancelBtn)  // Safe option highlighted by default
dialog.SetCancelButton(cancelBtn)   // Escape triggers Cancel
dialog.Show()
```

### Прикрепление к окну

Прикрепите диалоговое окно к определённому окну:

```go
dialog := app.Dialog.Info().
    SetMessage("Window-specific message").
    AttachToWindow(window)

dialog.Show()
```

**Поведение:**

- Диалоговое окно отображается по центру родительского окна
- Пока отображается диалоговое окно, родительское окно недоступно
- Диалоговое окно перемещается вместе с родительским окном (macOS)

## Поведение на разных платформах

@tabs{sync-key="platform"}
[macOS]
**Диалоговые окна macOS:**

- Нативный внешний вид NSAlert
- Соответствие системной теме (светлой или тёмной)
- Поддержка навигации с клавиатуры
- Стандартные сочетания клавиш (⌘. для отмены)
- Встроенные специальные возможности
- Отображение в виде панели при прикреплении к окну

**Пример:**

```go
// Appears as sheet on macOS
dialog := app.Dialog.Question().
    SetMessage("Save changes?").
    AttachToWindow(window)
dialog.AddButton("Yes")
dialog.AddButton("No")
dialog.Show()
```

[Windows]
**Диалоговые окна Windows:**

- Нативный внешний вид TaskDialog
- Соответствие системной теме
- Поддержка навигации с клавиатуры
- Стандартные сочетания клавиш (Esc для отмены)
- Встроенные специальные возможности
- Модальное по отношению к родительскому окну

**Пример:**

```go
// Modal dialog on Windows
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Operation failed").
    Show()
```

[Linux]
**Диалоговые окна Linux:**

- Внешний вид диалоговых окон GTK
- Соответствие теме рабочего стола
- Поддержка навигации с клавиатуры
- Интеграция со средой рабочего стола
- Зависит от среды рабочего стола (GNOME, KDE и т. д.)

**Пример:**

```go
// GTK dialog on Linux
app.Dialog.Info().
    SetMessage("Update complete").
    Show()
```

@end

## Распространённые сценарии

### Подтверждение перед деструктивным действием

```go
func deleteFile(app *application.App, path string) {
    dialog := app.Dialog.Question().
        SetTitle("Confirm Delete").
        SetMessage(fmt.Sprintf("Delete %s?", filepath.Base(path)))

    deleteBtn := dialog.AddButton("Delete")
    deleteBtn.OnClick(func() {
        if err := os.Remove(path); err != nil {
            app.Dialog.Error().
                SetTitle("Delete Failed").
                SetMessage(err.Error()).
                Show()
        }
    })

    cancelBtn := dialog.AddButton("Cancel")
    dialog.SetDefaultButton(cancelBtn)
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
}
```

### Обработка ошибок с помощью диалоговых окон

```go
func saveDocument(app *application.App, path string, data []byte) {
    if err := os.WriteFile(path, data, 0644); err != nil {
        app.Dialog.Error().
            SetTitle("Save Failed").
            SetMessage(fmt.Sprintf("Could not save file: %v", err)).
            Show()
        return
    }

    app.Dialog.Info().
        SetTitle("Success").
        SetMessage("File saved successfully!").
        Show()
}
```

### Выбор файла с проверкой

```go
func selectImageFile(app *application.App) (string, error) {
    path, err := app.Dialog.OpenFile().
        SetTitle("Select Image").
        AddFilter("Images", "*.png;*.jpg;*.jpeg;*.gif").
        PromptForSingleSelection()

    if err != nil {
        return "", err
    }

    if path == "" {
        return "", errors.New("no file selected")
    }

    // Validate file
    if !isValidImage(path) {
        app.Dialog.Error().
            SetTitle("Invalid File").
            SetMessage("Selected file is not a valid image.").
            Show()
        return "", errors.New("invalid image")
    }

    return path, nil
}
```

### Многоэтапная последовательность диалоговых окон

```go
func exportData(app *application.App) {
    // Step 1: Confirm export
    dialog := app.Dialog.Question().
        SetTitle("Export Data").
        SetMessage("Export all data to CSV?")

    exportBtn := dialog.AddButton("Export")
    exportBtn.OnClick(func() {
        // Step 2: Select destination
        path, err := app.Dialog.SaveFile().
            SetFilename("export.csv").
            AddFilter("CSV Files", "*.csv").
            PromptForSingleSelection()

        if err != nil || path == "" {
            return
        }

        // Step 3: Perform export
        if err := performExport(path); err != nil {
            app.Dialog.Error().
                SetTitle("Export Failed").
                SetMessage(err.Error()).
                Show()
            return
        }

        // Step 4: Success
        app.Dialog.Info().
            SetTitle("Export Complete").
            SetMessage("Data exported successfully!").
            Show()
    })

    cancelBtn := dialog.AddButton("Cancel")
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
}
```

## Рекомендации

### ✅ Рекомендуется

- **Используйте нативные диалоговые окна** — они удобнее для пользователей, чем пользовательские
- **Выводите понятные сообщения** — указывайте конкретную информацию
- **Задавайте подходящие заголовки** — контекст имеет значение
- **Разумно выбирайте кнопки по умолчанию** — по умолчанию выбирайте безопасный вариант
- **Обрабатывайте отмену** — пользователь может отменить действие
- **Проверяйте выбранные файлы** — проверяйте их типы

### ❌ Не делайте

- **Не злоупотребляйте диалоговыми окнами** — они прерывают рабочий процесс
- **Не используйте их для часто выводимых сообщений** — используйте уведомления
- **Не забывайте об обработке ошибок** — пользователь может отменить действие
- **Не блокируйте работу без необходимости** — рассмотрите альтернативы
- **Не используйте общие формулировки** — указывайте конкретную информацию
- **Не игнорируйте различия между платформами** — тестируйте на всех платформах

## Дальнейшие шаги

@cards{cols="2"}
ℹ Диалоговые окна сообщений
Диалоговые окна с информацией, предупреждениями и сообщениями об ошибках.

[Подробнее →](/features/dialogs/message/)

---
📖 Диалоговые окна выбора файлов
Открытие и сохранение файлов, а также выбор папок.

[Подробнее →](/features/dialogs/file/)

---
◆ Пользовательские диалоговые окна
Создание пользовательских диалоговых окон.

[Подробнее →](/features/dialogs/custom/)

---
▣ Окна
Узнайте об управлении окнами.

[Подробнее →](/features/windows/basics/)

@end

---

**Есть вопросы?** Задайте их в [Discord](https://discord.gg/JDdSxwjhGf) или ознакомьтесь с [примерами диалоговых окон](https://github.com/wailsapp/wails/tree/master/v3/examples/dialogs).
