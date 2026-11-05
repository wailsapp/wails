---
title: "API диалоговых окон"
description: "Полное справочное руководство по API нативных диалоговых окон"
slug: "reference/dialogs"
sourcePath: "reference/dialogs.md"
---

## Обзор

API диалоговых окон предоставляет методы для отображения нативных диалоговых окон выбора файлов и окон сообщений. Доступ к ним осуществляется через менеджер `app.Dialog`.

**Типы диалоговых окон:**

- **Диалоговые окна выбора файлов** — окна открытия и сохранения файлов
- **Диалоговые окна сообщений** — информационные окна, окна ошибок, предупреждений и вопросов

Все диалоговые окна являются **нативными диалоговыми окнами ОС** и соответствуют внешнему виду платформы.

## Доступ к диалоговым окнам

Доступ к диалоговым окнам осуществляется через менеджер `app.Dialog`:

```go
app.Dialog.OpenFile()
app.Dialog.SaveFile()
app.Dialog.Info()
app.Dialog.Question()
app.Dialog.Warning()
app.Dialog.Error()
```

## Диалоговые окна выбора файлов

### OpenFile()

Создаёт диалоговое окно открытия файла.

```go
func (dm *DialogManager) OpenFile() *OpenFileDialogStruct
```

**Пример:**

```go
dialog := app.Dialog.OpenFile()
```

### Методы OpenFileDialogStruct

#### SetTitle()

Задаёт заголовок диалогового окна.

```go
func (d *OpenFileDialogStruct) SetTitle(title string) *OpenFileDialogStruct
```

**Пример:**

```go
dialog.SetTitle("Select Image")
```

#### AddFilter()

Добавляет фильтр по типу файлов.

```go
func (d *OpenFileDialogStruct) AddFilter(displayName, pattern string) *OpenFileDialogStruct
```

**Параметры:**

- `displayName` — описание фильтра, отображаемое пользователю (например, «Изображения» или «Документы»)
- `pattern` — список расширений, разделённых точкой с запятой (например, «*.png;*.jpg»)

**Пример:**

```go
dialog.AddFilter("Images", "*.png;*.jpg;*.gif").
    AddFilter("Documents", "*.pdf;*.docx").
    AddFilter("All Files", "*.*")
```

#### SetDirectory()

Задаёт начальный каталог.

```go
func (d *OpenFileDialogStruct) SetDirectory(directory string) *OpenFileDialogStruct
```

**Пример:**

```go
homeDir, _ := os.UserHomeDir()
dialog.SetDirectory(homeDir)
```

#### CanChooseDirectories()

Разрешает или запрещает выбор каталогов.

```go
func (d *OpenFileDialogStruct) CanChooseDirectories(canChooseDirectories bool) *OpenFileDialogStruct
```

**Пример (выбор папки):**

```go
// Select folders instead of files
path, err := app.Dialog.OpenFile().
    SetTitle("Select Folder").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()
```

#### CanChooseFiles()

Разрешает или запрещает выбор файлов.

```go
func (d *OpenFileDialogStruct) CanChooseFiles(canChooseFiles bool) *OpenFileDialogStruct
```

#### CanCreateDirectories()

Разрешает или запрещает создание новых каталогов.

```go
func (d *OpenFileDialogStruct) CanCreateDirectories(canCreateDirectories bool) *OpenFileDialogStruct
```

#### ShowHiddenFiles()

Показывает или скрывает скрытые файлы.

```go
func (d *OpenFileDialogStruct) ShowHiddenFiles(showHiddenFiles bool) *OpenFileDialogStruct
```

#### AttachToWindow()

Привязывает диалоговое окно к указанному окну приложения.

```go
func (d *OpenFileDialogStruct) AttachToWindow(window Window) *OpenFileDialogStruct
```

#### PromptForSingleSelection()

Отображает диалоговое окно и возвращает выбранный файл.

```go
func (d *OpenFileDialogStruct) PromptForSingleSelection() (string, error)
```

**Возвращаемые значения:**

- `string` — путь к выбранному файлу. Пустая строка означает, что ничего не выбрано (в зависимости от ОС при отмене реализации для разных платформ могут возвращать либо пустую строку, либо ненулевую ошибку).
- `error` — ненулевое значение, если не удалось отобразить само диалоговое окно.

**Пример:**

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Image").
    AddFilter("Images", "*.png;*.jpg;*.gif").
    PromptForSingleSelection()

if err != nil {
    // The dialog failed to present (rare).
    return
}
if path == "" {
    // User cancelled.
    return
}

// Use the selected file
processFile(path)
```

#### PromptForMultipleSelection()

Отображает диалоговое окно и возвращает несколько выбранных файлов.

```go
func (d *OpenFileDialogStruct) PromptForMultipleSelection() ([]string, error)
```

**Возвращаемые значения:**

- `[]string` — массив путей к выбранным файлам
- `error` — ошибка при сбое диалогового окна

**Пример:**

```go
paths, err := app.Dialog.OpenFile().
    SetTitle("Select Images").
    AddFilter("Images", "*.png;*.jpg").
    PromptForMultipleSelection()

if err != nil {
    return
}

for _, path := range paths {
    processFile(path)
}
```

### SaveFile()

Создаёт диалоговое окно сохранения файла.

```go
func (dm *DialogManager) SaveFile() *SaveFileDialogStruct
```

**Пример:**

```go
dialog := app.Dialog.SaveFile()
```

### Методы SaveFileDialogStruct

#### SetTitle()

Задаёт заголовок диалогового окна.

```go
func (d *SaveFileDialogStruct) SetTitle(title string) *SaveFileDialogStruct
```

#### SetFilename()

Задаёт имя файла по умолчанию.

```go
func (d *SaveFileDialogStruct) SetFilename(filename string) *SaveFileDialogStruct
```

**Пример:**

```go
dialog.SetFilename("document.pdf")
```

#### AddFilter()

Добавляет фильтр по типу файлов.

```go
func (d *SaveFileDialogStruct) AddFilter(displayName, pattern string) *SaveFileDialogStruct
```

**Пример:**

```go
dialog.AddFilter("PDF Document", "*.pdf").
    AddFilter("Text Document", "*.txt")
```

#### SetDirectory()

Задаёт начальный каталог.

```go
func (d *SaveFileDialogStruct) SetDirectory(directory string) *SaveFileDialogStruct
```

#### AttachToWindow()

Привязывает диалоговое окно к указанному окну приложения.

```go
func (d *SaveFileDialogStruct) AttachToWindow(window Window) *SaveFileDialogStruct
```

#### PromptForSingleSelection()

Отображает диалоговое окно и возвращает путь сохранения.

```go
func (d *SaveFileDialogStruct) PromptForSingleSelection() (string, error)
```

**Пример:**

```go
path, err := app.Dialog.SaveFile().
    SetTitle("Save Document").
    SetFilename("untitled.pdf").
    AddFilter("PDF Document", "*.pdf").
    PromptForSingleSelection()

if err != nil {
    // User cancelled
    return
}

// Save to the selected path
saveDocument(path)
```

### Выбор папки

Отдельного `SelectFolderDialog` нет. Используйте `OpenFile()` с параметрами выбора каталогов:

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Output Folder").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()

if err != nil {
    // User cancelled
    return
}

// Use the selected folder
outputDir = path
```

## Диалоговые окна сообщений

Все диалоговые окна сообщений возвращают `*MessageDialog` и имеют одинаковые методы.

### Info()

Создаёт информационное диалоговое окно.

```go
func (dm *DialogManager) Info() *MessageDialog
```

**Пример:**

```go
app.Dialog.Info().
    SetTitle("Success").
    SetMessage("File saved successfully!").
    Show()
```

### Error()

Создаёт диалоговое окно с сообщением об ошибке.

```go
func (dm *DialogManager) Error() *MessageDialog
```

**Пример:**

```go
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Failed to save file: " + err.Error()).
    Show()
```

### Warning()

Создаёт диалоговое окно с предупреждением.

```go
func (dm *DialogManager) Warning() *MessageDialog
```

**Пример:**

```go
app.Dialog.Warning().
    SetTitle("Warning").
    SetMessage("This action cannot be undone.").
    Show()
```

### Question()

Создаёт диалоговое окно с вопросом и настраиваемыми кнопками.

```go
func (dm *DialogManager) Question() *MessageDialog
```

**Пример:**

```go
dialog := app.Dialog.Question().
    SetTitle("Confirm").
    SetMessage("Do you want to save changes?")

save := dialog.AddButton("Save")
save.OnClick(func() {
    saveDocument()
})

dontSave := dialog.AddButton("Don't Save")
dontSave.OnClick(func() {
    // Continue without saving
})

cancel := dialog.AddButton("Cancel")
cancel.OnClick(func() {
    // Do nothing
})

dialog.SetDefaultButton(save)
dialog.SetCancelButton(cancel)
dialog.Show()
```

### Методы MessageDialog

#### SetTitle()

Задаёт заголовок диалогового окна.

```go
func (d *MessageDialog) SetTitle(title string) *MessageDialog
```

#### SetMessage()

Задаёт сообщение диалогового окна.

```go
func (d *MessageDialog) SetMessage(message string) *MessageDialog
```

#### SetIcon()

Задаёт пользовательский значок диалогового окна.

```go
func (d *MessageDialog) SetIcon(icon []byte) *MessageDialog
```

#### AddButton()

Добавляет кнопку в диалоговое окно и возвращает её для дальнейшей настройки.

```go
func (d *MessageDialog) AddButton(label string) *Button
```

**Возвращаемое значение:** `*Button` — экземпляр кнопки для дальнейшей настройки

**Пример:**

```go
button := dialog.AddButton("OK")
button.OnClick(func() {
    // Handle click
})
```

#### SetDefaultButton()

Назначает кнопку по умолчанию, которая активируется при нажатии Enter.

```go
func (d *MessageDialog) SetDefaultButton(button *Button) *MessageDialog
```

**Пример:**

```go
yes := dialog.AddButton("Yes")
no := dialog.AddButton("No")
dialog.SetDefaultButton(yes)
```

#### SetCancelButton()

Назначает кнопку отмены, которая активируется при нажатии Escape.

```go
func (d *MessageDialog) SetCancelButton(button *Button) *MessageDialog
```

**Пример:**

```go
ok := dialog.AddButton("OK")
cancel := dialog.AddButton("Cancel")
dialog.SetCancelButton(cancel)
```

#### AttachToWindow()

Привязывает диалоговое окно к указанному окну.

```go
func (d *MessageDialog) AttachToWindow(window Window) *MessageDialog
```

#### Show()

Отображает диалоговое окно. Ответы пользователя обрабатываются функциями обратного вызова кнопок.

```go
func (d *MessageDialog) Show()
```

**Примечание:** `Show()` не возвращает значение. Обрабатывайте ответы пользователя с помощью функций обратного вызова кнопок.

### Методы кнопок

#### OnClick()

Задаёт функцию обратного вызова, которая выполняется при нажатии кнопки.

```go
func (b *Button) OnClick(callback func()) *Button
```

#### SetAsDefault()

Назначает эту кнопку кнопкой по умолчанию.

```go
func (b *Button) SetAsDefault() *Button
```

#### SetAsCancel()

Назначает эту кнопку кнопкой отмены.

```go
func (b *Button) SetAsCancel() *Button
```

## Полные примеры

### Пример выбора файла

```go
type FileService struct {
    app *application.App
}

func (s *FileService) OpenImage() (string, error) {
    path, err := s.app.Dialog.OpenFile().
        SetTitle("Select Image").
        AddFilter("Images", "*.png;*.jpg;*.jpeg;*.gif").
        AddFilter("All Files", "*.*").
        PromptForSingleSelection()

    if err != nil {
        return "", err
    }

    return path, nil
}

func (s *FileService) SaveDocument(defaultName string) (string, error) {
    path, err := s.app.Dialog.SaveFile().
        SetTitle("Save Document").
        SetFilename(defaultName).
        AddFilter("PDF Document", "*.pdf").
        AddFilter("Text Document", "*.txt").
        PromptForSingleSelection()

    if err != nil {
        return "", err
    }

    return path, nil
}

func (s *FileService) SelectOutputFolder() (string, error) {
    path, err := s.app.Dialog.OpenFile().
        SetTitle("Select Output Folder").
        CanChooseDirectories(true).
        CanChooseFiles(false).
        PromptForSingleSelection()

    if err != nil {
        return "", err
    }

    return path, nil
}
```

### Пример диалогового окна подтверждения

```go
func (s *Service) DeleteItem(app *application.App, id string) {
    dialog := app.Dialog.Question().
        SetTitle("Confirm Delete").
        SetMessage("Are you sure you want to delete this item?")

    deleteBtn := dialog.AddButton("Delete")
    deleteBtn.OnClick(func() {
        deleteFromDatabase(id)
    })

    cancelBtn := dialog.AddButton("Cancel")
    // Cancel does nothing

    dialog.SetDefaultButton(cancelBtn) // Default to Cancel for safety
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
}
```

### Диалоговое окно сохранения изменений

```go
func (s *Editor) PromptSaveChanges(app *application.App) {
    dialog := app.Dialog.Question().
        SetTitle("Unsaved Changes").
        SetMessage("Do you want to save your changes before closing?")

    save := dialog.AddButton("Save")
    save.OnClick(func() {
        s.Save()
        s.Close()
    })

    dontSave := dialog.AddButton("Don't Save")
    dontSave.OnClick(func() {
        s.Close()
    })

    cancel := dialog.AddButton("Cancel")
    // Cancel does nothing, dialog closes

    dialog.SetDefaultButton(save)
    dialog.SetCancelButton(cancel)
    dialog.Show()
}
```

### Обработка нескольких файлов

```go
func (s *Service) ProcessMultipleFiles(app *application.App) error {
    // Select multiple files
    paths, err := app.Dialog.OpenFile().
        SetTitle("Select Files to Process").
        AddFilter("Images", "*.png;*.jpg").
        PromptForMultipleSelection()

    if err != nil {
        return err
    }

    if len(paths) == 0 {
        app.Dialog.Info().
            SetTitle("No Files Selected").
            SetMessage("Please select at least one file.").
            Show()
        return nil
    }

    // Process files
    for _, path := range paths {
        err := processFile(path)
        if err != nil {
            app.Dialog.Error().
                SetTitle("Processing Error").
                SetMessage(fmt.Sprintf("Failed to process %s: %v", path, err)).
                Show()
            continue
        }
    }

    // Show completion
    app.Dialog.Info().
        SetTitle("Complete").
        SetMessage(fmt.Sprintf("Successfully processed %d files", len(paths))).
        Show()

    return nil
}
```

### Обработка ошибок с помощью диалоговых окон

```go
func (s *Service) SaveFile(app *application.App, data []byte) error {
    // Select save location
    path, err := app.Dialog.SaveFile().
        SetTitle("Save File").
        SetFilename("data.json").
        AddFilter("JSON File", "*.json").
        PromptForSingleSelection()

    if err != nil {
        // User cancelled - not an error
        return nil
    }

    // Attempt to save
    err = os.WriteFile(path, data, 0644)
    if err != nil {
        // Show error dialog
        app.Dialog.Error().
            SetTitle("Save Failed").
            SetMessage(fmt.Sprintf("Could not save file: %v", err)).
            Show()
        return err
    }

    // Show success
    app.Dialog.Info().
        SetTitle("Success").
        SetMessage("File saved successfully!").
        Show()

    return nil
}
```

### Значения по умолчанию для разных платформ

```go
import (
    "os"
    "path/filepath"
    "runtime"
)

func (s *Service) GetDefaultDirectory() string {
    homeDir, _ := os.UserHomeDir()

    switch runtime.GOOS {
    case "windows":
        return filepath.Join(homeDir, "Documents")
    case "darwin":
        return filepath.Join(homeDir, "Documents")
    case "linux":
        return filepath.Join(homeDir, "Documents")
    default:
        return homeDir
    }
}

func (s *Service) OpenWithDefaults(app *application.App) (string, error) {
    return app.Dialog.OpenFile().
        SetTitle("Open File").
        SetDirectory(s.GetDefaultDirectory()).
        AddFilter("All Files", "*.*").
        PromptForSingleSelection()
}
```

## Рекомендации

### Рекомендуется

- **Используйте нативные диалоговые окна** — они соответствуют внешнему виду и поведению платформы
- **Используйте понятные заголовки** — они помогают пользователям понять назначение диалогового окна
- **Задавайте подходящие фильтры** — они помогают пользователям выбирать файлы нужных типов
- **Обрабатывайте отмену** — проверяйте наличие ошибок, поскольку пользователь может отменить действие
- **Запрашивайте подтверждение разрушительных действий** — используйте диалоговые окна Question
- **Предоставляйте обратную связь** — используйте диалоговые окна Info для сообщений об успешном выполнении
- **Задавайте разумные значения по умолчанию** — каталог, имя файла и т. д.
- **Используйте функции обратного вызова для действий кнопок** — корректно обрабатывайте ответы пользователя

### Не рекомендуется

- **Не игнорируйте ошибки** — при отмене пользователем возвращается ошибка
- **Не используйте неоднозначные подписи кнопок** — указывайте действие явно: «Сохранить»/«Отмена»
- **Не злоупотребляйте диалоговыми окнами** — они прерывают рабочий процесс
- **Не показывайте ошибки при отмене** — это обычное действие
- **Не забывайте о фильтрах файлов** — помогите пользователям найти нужные файлы
- **Не задавайте пути жёстко в коде** — используйте os.UserHomeDir() или аналогичную функцию

## Типы диалоговых окон на разных платформах

### macOS

- Диалоговые окна выдвигаются из строки заголовка
- Диалоговое окно в стиле «листа», прикреплённое к родительскому окну
- Нативный внешний вид macOS

### Windows

- Стандартные диалоговые окна Windows
- Соответствуют рекомендациям по дизайну Windows
- Современный внешний вид Windows 10/11

### Linux

- Диалоговые окна GTK в системах на базе GTK
- Диалоговые окна Qt в системах на базе Qt
- Соответствуют окружению рабочего стола

#### Поведение диалоговых окон в Linux

В Linux сборка GTK4 по умолчанию использует **xdg-desktop-portal** для диалоговых окон выбора файлов, что обеспечивает нативную интеграцию с рабочим столом, но приводит к тому, что некоторые параметры не действуют. Устаревший вариант на GTK3 (`-tags gtk3`) сохраняет полный программный контроль над этими параметрами:

| Параметр | GTK3 (`-tags gtk3`) | GTK4 (по умолчанию) | Примечания |
| --- | --- | --- | --- |
| `ShowHiddenFiles()` | ✅ Работает | ❌ Не действует | Пользователь управляет этим с помощью переключателя в интерфейсе диалогового окна (Ctrl+H или меню) |
| `CanCreateDirectories()` | ✅ Работает | ❌ Не действует | Всегда включено в портале |
| `ResolvesAliases()` | ✅ Работает | ❌ Не действует | Портал обрабатывает разрешение символических ссылок |
| `SetButtonText()` | ✅ Работает | ✅ Работает | Пользовательский текст кнопки подтверждения работает |

**Почему существуют эти ограничения:** Диалоговые окна GTK4 на основе портала передают управление интерфейсом окружению рабочего стола (GNOME, KDE и т. д.). Так и задумано: портал обеспечивает единообразное взаимодействие с пользователем во всех приложениях и учитывает пользовательские настройки.

@note{type="info"}
Сборка GTK4 по умолчанию использует диалоговые окна на основе портала. Если вашему приложению требуется полный программный контроль над перечисленными выше параметрами диалоговых окон, используйте сборку с устаревшим вариантом `-tags gtk3` (поддерживается до версии v3.0.x включительно; удалён в v3.1) — см. [«Пакетирование для Linux — поддержка устаревшего GTK3»](/guides/build/linux/#legacy-gtk3-support).

@end

## Распространённые шаблоны

### Шаблон «Сохранить как»

```go
func (s *Service) SaveAs(app *application.App, currentPath string) (string, error) {
    // Extract filename from current path
    filename := filepath.Base(currentPath)

    // Show save dialog
    path, err := app.Dialog.SaveFile().
        SetTitle("Save As").
        SetFilename(filename).
        PromptForSingleSelection()

    if err != nil {
        return "", err
    }

    return path, nil
}
```

### Шаблон «Открыть недавний файл»

```go
func (s *Service) OpenRecent(app *application.App, recentPath string) error {
    // Check if file still exists
    if _, err := os.Stat(recentPath); os.IsNotExist(err) {
        dialog := app.Dialog.Question().
            SetTitle("File Not Found").
            SetMessage("The file no longer exists. Remove from recent files?")

        remove := dialog.AddButton("Remove")
        remove.OnClick(func() {
            s.removeFromRecent(recentPath)
        })

        cancel := dialog.AddButton("Cancel")
        dialog.SetCancelButton(cancel)
        dialog.Show()

        return err
    }

    return s.openFile(recentPath)
}
```
