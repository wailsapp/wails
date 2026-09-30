---
title: "Диалоговые окна для работы с файлами"
description: "Диалоговые окна открытия, сохранения и выбора папок"
slug: "features/dialogs/file"
sourcePath: "features/dialogs/file.md"
---

## Диалоговые окна для работы с файлами

Wails предоставляет **нативные диалоговые окна для работы с файлами**, внешний вид которых соответствует платформе, для открытия и сохранения файлов, а также выбора папок. Простой API поддерживает фильтрацию по типам файлов, множественный выбор и начальные расположения.

![Нативное окно выбора файлов macOS, открытое приложением Wails](/assets/screenshots/file-dialog-macos.png)

API Wails делегирует выбор операционной системе, поэтому сохраняются привычные возможности навигации, фильтрации и выбора.

## Создание диалоговых окон для работы с файлами

Доступ к диалоговым окнам для работы с файлами осуществляется через диспетчер `app.Dialog`:

```go
app.Dialog.OpenFile()
app.Dialog.SaveFile()
```

## Диалоговое окно открытия файла

Выберите файлы для открытия:

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Image").
    AddFilter("Images", "*.png;*.jpg;*.gif").
    AddFilter("All Files", "*.*").
    PromptForSingleSelection()

if err != nil || path == "" {
    return
}

openFile(path)
```

**Варианты использования:**

- Открытие документов
- Импорт файлов
- Загрузка изображений
- Выбор файлов конфигурации

### Выбор одного файла

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Open Document").
    AddFilter("Text Files", "*.txt").
    AddFilter("All Files", "*.*").
    PromptForSingleSelection()

if err != nil || path == "" {
    // User cancelled or error occurred
    return
}

// Use selected file
data, _ := os.ReadFile(path)
```

### Выбор нескольких файлов

```go
paths, err := app.Dialog.OpenFile().
    SetTitle("Select Images").
    AddFilter("Images", "*.png;*.jpg;*.jpeg;*.gif").
    PromptForMultipleSelection()

if err != nil {
    return
}

// Process all selected files
for _, path := range paths {
    processFile(path)
}
```

### С начальным каталогом

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Open File").
    SetDirectory("/Users/me/Documents").
    PromptForSingleSelection()
```

## Диалоговое окно сохранения файла

Выберите место сохранения:

```go
path, err := app.Dialog.SaveFile().
    SetFilename("document.txt").
    AddFilter("Text Files", "*.txt").
    AddFilter("All Files", "*.*").
    PromptForSingleSelection()

if err != nil || path == "" {
    return
}

saveFile(path, data)
```

**Варианты использования:**

- Сохранение документов
- Экспорт данных
- Создание новых файлов
- Сохранение под другим именем...

### С именем файла по умолчанию

```go
path, err := app.Dialog.SaveFile().
    SetFilename("export.csv").
    AddFilter("CSV Files", "*.csv").
    PromptForSingleSelection()
```

### С начальным каталогом

```go
path, err := app.Dialog.SaveFile().
    SetDirectory("/Users/me/Documents").
    SetFilename("untitled.txt").
    PromptForSingleSelection()
```

### Подтверждение перезаписи

```go
path, err := app.Dialog.SaveFile().
    SetFilename("document.txt").
    PromptForSingleSelection()

if err != nil || path == "" {
    return
}

// Check if file exists
if _, err := os.Stat(path); err == nil {
    dialog := app.Dialog.Question().
        SetTitle("Confirm Overwrite").
        SetMessage("File already exists. Overwrite?")

    overwriteBtn := dialog.AddButton("Overwrite")
    overwriteBtn.OnClick(func() {
        saveFile(path, data)
    })

    cancelBtn := dialog.AddButton("Cancel")
    dialog.SetDefaultButton(cancelBtn)
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
    return
}

saveFile(path, data)
```

## Диалоговое окно выбора папки

Выберите каталог с помощью диалогового окна открытия файла, включив в нём выбор каталогов:

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Output Folder").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()

if err != nil || path == "" {
    return
}

exportToFolder(path)
```

**Варианты использования:**

- Выбор выходного каталога
- Выбор рабочего пространства
- Выбор места для резервной копии
- Выбор каталога установки

### С начальным каталогом

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Folder").
    SetDirectory("/Users/me/Documents").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()
```

## Фильтры файлов

Добавляйте фильтры типов файлов в диалоговые окна с помощью метода `AddFilter()`. Каждый вызов добавляет новый вариант фильтра.

### Базовые фильтры

```go
path, _ := app.Dialog.OpenFile().
    AddFilter("Text Files", "*.txt").
    AddFilter("All Files", "*.*").
    PromptForSingleSelection()
```

### Несколько расширений

Чтобы указать несколько расширений в одном фильтре, разделяйте их точками с запятой:

```go
dialog := app.Dialog.OpenFile().
    AddFilter("Images", "*.png;*.jpg;*.jpeg;*.gif;*.bmp").
    AddFilter("Documents", "*.txt;*.doc;*.docx;*.pdf").
    AddFilter("All Files", "*.*")
```

### Формат шаблона

Чтобы разделить несколько расширений в одном фильтре, используйте **точки с запятой**:

```go
// Multiple extensions separated by semicolons
AddFilter("Images", "*.png;*.jpg;*.gif")
```

## Полные примеры

### Открытие файла изображения

```go
func openImage(app *application.App) (image.Image, error) {
    path, err := app.Dialog.OpenFile().
        SetTitle("Select Image").
        AddFilter("Images", "*.png;*.jpg;*.jpeg;*.gif;*.bmp").
        PromptForSingleSelection()

    if err != nil {
        return nil, err
    }

    if path == "" {
        return nil, errors.New("no file selected")
    }

    // Open and decode image
    file, err := os.Open(path)
    if err != nil {
        app.Dialog.Error().
            SetTitle("Open Failed").
            SetMessage(err.Error()).
            Show()
        return nil, err
    }
    defer file.Close()

    img, _, err := image.Decode(file)
    if err != nil {
        app.Dialog.Error().
            SetTitle("Invalid Image").
            SetMessage("Could not decode image file.").
            Show()
        return nil, err
    }

    return img, nil
}
```

### Сохранение документа с проверкой

```go
func saveDocument(app *application.App, content string) {
    path, err := app.Dialog.SaveFile().
        SetFilename("document.txt").
        AddFilter("Text Files", "*.txt").
        AddFilter("Markdown Files", "*.md").
        AddFilter("All Files", "*.*").
        PromptForSingleSelection()

    if err != nil || path == "" {
        return
    }

    // Validate extension
    ext := filepath.Ext(path)
    if ext != ".txt" && ext != ".md" {
        dialog := app.Dialog.Question().
            SetTitle("Confirm Extension").
            SetMessage(fmt.Sprintf("Save as %s file?", ext))

        saveBtn := dialog.AddButton("Save")
        saveBtn.OnClick(func() {
            doSave(app, path, content)
        })

        cancelBtn := dialog.AddButton("Cancel")
        dialog.SetDefaultButton(cancelBtn)
        dialog.SetCancelButton(cancelBtn)
        dialog.Show()
        return
    }

    doSave(app, path, content)
}

func doSave(app *application.App, path, content string) {
    if err := os.WriteFile(path, []byte(content), 0644); err != nil {
        app.Dialog.Error().
            SetTitle("Save Failed").
            SetMessage(err.Error()).
            Show()
        return
    }

    app.Dialog.Info().
        SetTitle("Saved").
        SetMessage("Document saved successfully!").
        Show()
}
```

### Пакетная обработка файлов

```go
func processMultipleFiles(app *application.App) {
    paths, err := app.Dialog.OpenFile().
        SetTitle("Select Files to Process").
        AddFilter("Images", "*.png;*.jpg").
        PromptForMultipleSelection()

    if err != nil || len(paths) == 0 {
        return
    }

    // Confirm processing
    dialog := app.Dialog.Question().
        SetTitle("Confirm Processing").
        SetMessage(fmt.Sprintf("Process %d file(s)?", len(paths)))

    processBtn := dialog.AddButton("Process")
    processBtn.OnClick(func() {
        // Process files
        var errs []error
        for i, path := range paths {
            if err := processFile(path); err != nil {
                errs = append(errs, err)
            }

            // Update progress
            // app.Event.Emit("progress", map[string]interface{}{
            //     "current": i + 1,
            //     "total":   len(paths),
            // })
            _ = i // suppress unused variable warning in example
        }

        // Show results
        if len(errs) > 0 {
            app.Dialog.Warning().
                SetTitle("Processing Complete").
                SetMessage(fmt.Sprintf("Processed %d files with %d errors.",
                    len(paths), len(errs))).
                Show()
        } else {
            app.Dialog.Info().
                SetTitle("Success").
                SetMessage(fmt.Sprintf("Processed %d files successfully!", len(paths))).
                Show()
        }
    })

    cancelBtn := dialog.AddButton("Cancel")
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
}
```

### Экспорт с выбором папки

```go
func exportData(app *application.App, data []byte) {
    // Select output folder
    folder, err := app.Dialog.OpenFile().
        SetTitle("Select Export Folder").
        SetDirectory(getDefaultExportFolder()).
        CanChooseDirectories(true).
        CanChooseFiles(false).
        PromptForSingleSelection()

    if err != nil || folder == "" {
        return
    }

    // Generate filename
    filename := fmt.Sprintf("export_%s.csv",
        time.Now().Format("2006-01-02_15-04-05"))
    path := filepath.Join(folder, filename)

    // Save file
    if err := os.WriteFile(path, data, 0644); err != nil {
        app.Dialog.Error().
            SetTitle("Export Failed").
            SetMessage(err.Error()).
            Show()
        return
    }

    // Show success with option to open folder
    dialog := app.Dialog.Question().
        SetTitle("Export Complete").
        SetMessage(fmt.Sprintf("Exported to %s", filename))

    openBtn := dialog.AddButton("Open Folder")
    openBtn.OnClick(func() {
        openFolder(folder)
    })

    dialog.AddButton("OK")
    dialog.Show()
}
```

### Импорт с проверкой

```go
func importConfiguration(app *application.App) {
    path, err := app.Dialog.OpenFile().
        SetTitle("Import Configuration").
        AddFilter("JSON Files", "*.json").
        AddFilter("YAML Files", "*.yaml;*.yml").
        PromptForSingleSelection()

    if err != nil || path == "" {
        return
    }

    // Read file
    data, err := os.ReadFile(path)
    if err != nil {
        app.Dialog.Error().
            SetTitle("Read Failed").
            SetMessage(err.Error()).
            Show()
        return
    }

    // Validate configuration
    config, err := parseConfig(data)
    if err != nil {
        app.Dialog.Error().
            SetTitle("Invalid Configuration").
            SetMessage("File is not a valid configuration.").
            Show()
        return
    }

    // Confirm import
    dialog := app.Dialog.Question().
        SetTitle("Confirm Import").
        SetMessage("Import this configuration?")

    importBtn := dialog.AddButton("Import")
    importBtn.OnClick(func() {
        // Apply configuration
        if err := applyConfig(config); err != nil {
            app.Dialog.Error().
                SetTitle("Import Failed").
                SetMessage(err.Error()).
                Show()
            return
        }

        app.Dialog.Info().
            SetTitle("Success").
            SetMessage("Configuration imported successfully!").
            Show()
    })

    cancelBtn := dialog.AddButton("Cancel")
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
}
```

## Рекомендации

### ✅ Рекомендуется

- **Добавляйте фильтры файлов** — помогайте пользователям находить файлы
- **Задавайте подходящие заголовки** — обеспечивайте понятный контекст
- **Используйте начальные каталоги** — начинайте с логичного расположения
- **Проверяйте выбранные элементы** — проверяйте типы файлов
- **Обрабатывайте отмену** — пользователь может отменить действие
- **Запрашивайте подтверждение** — для разрушительных действий
- **Предоставляйте обратную связь** — сообщайте об успешном выполнении или ошибках

### ❌ Не рекомендуется

- **Не пропускайте проверку** — проверяйте типы файлов
- **Не игнорируйте ошибки** — обрабатывайте отмену
- **Не используйте слишком общие фильтры** — указывайте конкретные типы
- **Не забывайте вариант «Все файлы»** — всегда включайте его в список
- **Не задавайте пути жёстко** — используйте домашний каталог пользователя
- **Не предполагайте, что файл существует** — проверяйте это перед открытием

## Различия между платформами

### macOS

- Нативные панели NSOpenPanel/NSSavePanel
- При привязке к окну отображаются в виде диалогового листа
- Следуют системной теме
- Поддерживают предварительный просмотр Quick Look
- Интегрируются с тегами и избранным

### Windows

- Нативные диалоговые окна открытия и сохранения файлов
- Следуют системной теме
- Интеграция с недавними файлами
- Поддержка сетевых расположений

### Linux

- Диалог выбора файлов GTK
- Зависит от среды рабочего стола
- Соответствует теме рабочего стола
- Поддержка недавних файлов

## Дальнейшие шаги

@cards{cols="2"}
ℹ Диалоговые окна сообщений
Диалоговые окна с информацией, предупреждениями и ошибками.

[Подробнее →](/features/dialogs/message/)

---
◆ Пользовательские диалоговые окна
Создавайте пользовательские диалоговые окна.

[Подробнее →](/features/dialogs/custom/)

---
🚀 Привязки
Вызывайте функции Go из JavaScript.

[Подробнее →](/features/bindings/methods/)

---
★ События
Используйте события для обновления информации о ходе выполнения.

[Подробнее →](/features/events/system/)

@end

---

**Есть вопросы?** Задайте их в [Discord](https://discord.gg/JDdSxwjhGf) или ознакомьтесь с [примерами диалоговых окон выбора файлов](https://github.com/wailsapp/wails/tree/master/v3/examples/dialogs).
