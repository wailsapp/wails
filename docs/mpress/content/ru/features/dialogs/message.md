---
title: "Диалоговые окна сообщений"
description: "Отображение информации, предупреждений, ошибок и вопросов"
slug: "features/dialogs/message"
sourcePath: "features/dialogs/message.md"
---

## Диалоговые окна сообщений

Wails предоставляет **нативные диалоговые окна сообщений**, внешний вид которых соответствует платформе: информационные окна, предупреждения, сообщения об ошибках и вопросы с настраиваемыми заголовками, текстом и кнопками. Простой API, нативное поведение и доступность по умолчанию.

## Создание диалоговых окон

Доступ к диалоговым окнам сообщений осуществляется через менеджер `app.Dialog`:

```go
app.Dialog.Info()
app.Dialog.Question()
app.Dialog.Warning()
app.Dialog.Error()
```

Все методы возвращают объект `*MessageDialog`, который можно настроить с помощью цепочки вызовов методов.

## Информационное диалоговое окно

Отображайте информационные сообщения:

```go
app.Dialog.Info().
    SetTitle("Success").
    SetMessage("File saved successfully!").
    Show()
```

**Варианты использования:**

- Подтверждения успешного выполнения
- Уведомления о завершении
- Информационные сообщения
- Обновления состояния

**Пример — подтверждение сохранения:**

```go
func saveFile(app *application.App, path string, data []byte) error {
    if err := os.WriteFile(path, data, 0644); err != nil {
        return err
    }

    app.Dialog.Info().
        SetTitle("File Saved").
        SetMessage(fmt.Sprintf("Saved to %s", filepath.Base(path))).
        Show()

    return nil
}
```

## Диалоговое окно предупреждения

Отображайте предупреждения:

```go
app.Dialog.Warning().
    SetTitle("Warning").
    SetMessage("This action cannot be undone.").
    Show()
```

**Варианты использования:**

- Некритические предупреждения
- Уведомления об устаревании
- Предупреждающие сообщения
- Потенциальные проблемы

**Пример — предупреждение о свободном месте на диске:**

```go
func checkDiskSpace(app *application.App) {
    available := getDiskSpace()

    if available < 100*1024*1024 { // Less than 100MB
        app.Dialog.Warning().
            SetTitle("Low Disk Space").
            SetMessage(fmt.Sprintf("Only %d MB available.", available/(1024*1024))).
            Show()
    }
}
```

## Диалоговое окно ошибки

Отображайте ошибки:

```go
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Failed to connect to server.").
    Show()
```

**Варианты использования:**

- Сообщения об ошибках
- Уведомления о сбоях
- Обработка исключений
- Критические проблемы

**Пример — сетевая ошибка:**

```go
func fetchData(app *application.App, url string) ([]byte, error) {
    resp, err := http.Get(url)
    if err != nil {
        app.Dialog.Error().
            SetTitle("Network Error").
            SetMessage(fmt.Sprintf("Failed to connect: %v", err)).
            Show()
        return nil, err
    }
    defer resp.Body.Close()

    return io.ReadAll(resp.Body)
}
```

## Диалоговое окно с вопросом

Задавайте пользователям вопросы и обрабатывайте ответы с помощью обратных вызовов кнопок:

```go
dialog := app.Dialog.Question().
    SetTitle("Confirm").
    SetMessage("Save changes before closing?")

save := dialog.AddButton("Save")
save.OnClick(func() {
    saveChanges()
})

dontSave := dialog.AddButton("Don't Save")
dontSave.OnClick(func() {
    // Continue without saving
})

cancel := dialog.AddButton("Cancel")
cancel.OnClick(func() {
    // Don't close
})

dialog.SetDefaultButton(save)
dialog.SetCancelButton(cancel)
dialog.Show()
```

**Варианты использования:**

- Подтверждение действий
- Вопросы с ответами «Да» и «Нет»
- Выбор из нескольких вариантов
- Решения пользователя

**Пример — несохранённые изменения:**

```go
func closeDocument(app *application.App) {
    if !hasUnsavedChanges() {
        doClose()
        return
    }

    dialog := app.Dialog.Question().
        SetTitle("Unsaved Changes").
        SetMessage("Do you want to save your changes?")

    save := dialog.AddButton("Save")
    save.OnClick(func() {
        if saveDocument() {
            doClose()
        }
    })

    dontSave := dialog.AddButton("Don't Save")
    dontSave.OnClick(func() {
        doClose()
    })

    cancel := dialog.AddButton("Cancel")
    // Cancel button has no callback - just closes the dialog

    dialog.SetDefaultButton(save)
    dialog.SetCancelButton(cancel)
    dialog.Show()
}
```

## Параметры диалогового окна

### Заголовок и сообщение

```go
dialog := app.Dialog.Info().
    SetTitle("Operation Complete").
    SetMessage("All files have been processed successfully.")
```

**Рекомендации:**

- **Заголовок:** краткий и информативный (2-5 слов)
- **Сообщение:** понятное, конкретное и содержащее указание к действию
- **Избегайте жаргона:** используйте простой язык

### Кнопки

**Одна кнопка (информация, предупреждение или ошибка):**

В информационных диалоговых окнах, предупреждениях и сообщениях об ошибках по умолчанию отображается кнопка «ОК»:

```go
app.Dialog.Info().
    SetMessage("Done!").
    Show()
```

Также можно добавлять собственные кнопки:

```go
dialog := app.Dialog.Info().
    SetMessage("Done!")
ok := dialog.AddButton("Got it!")
dialog.SetDefaultButton(ok)
dialog.Show()
```

**Несколько кнопок (вопрос):**

Добавляйте кнопки с помощью `AddButton()`. Этот метод возвращает объект `*Button`, который можно настроить:

```go
dialog := app.Dialog.Question().
    SetMessage("Choose an action")

option1 := dialog.AddButton("Option 1")
option1.OnClick(func() {
    handleOption1()
})

option2 := dialog.AddButton("Option 2")
option2.OnClick(func() {
    handleOption2()
})

option3 := dialog.AddButton("Option 3")
option3.OnClick(func() {
    handleOption3()
})

dialog.Show()
```

**Кнопки по умолчанию и отмены:**

Используйте `SetDefaultButton()`, чтобы указать кнопку, которая будет выделена и активирована клавишей Enter. Используйте `SetCancelButton()`, чтобы указать кнопку, которая будет активирована клавишей Escape.

```go
dialog := app.Dialog.Question().
    SetMessage("Delete file?")

deleteBtn := dialog.AddButton("Delete")
deleteBtn.OnClick(func() {
    performDelete()
})

cancelBtn := dialog.AddButton("Cancel")
// No callback needed - just dismisses dialog

dialog.SetDefaultButton(cancelBtn)  // Safe option as default
dialog.SetCancelButton(cancelBtn)   // Escape triggers Cancel
dialog.Show()
```

Для кнопок также можно использовать методы `SetAsDefault()` и `SetAsCancel()` с цепочкой вызовов:

```go
dialog := app.Dialog.Question().
    SetMessage("Delete file?")

dialog.AddButton("Delete").OnClick(func() {
    performDelete()
})

dialog.AddButton("Cancel").SetAsDefault().SetAsCancel()

dialog.Show()
```

**Рекомендации:**

- **1-3 кнопки:** не перегружайте пользователей выбором
- **Понятные подписи:** «Сохранить» вместо «ОК»
- **Безопасное действие по умолчанию:** действие без необратимых последствий
- **Порядок важен:** наиболее вероятное действие размещайте первым (кроме отмены)

### Пользовательский значок

Задайте пользовательский значок для диалогового окна:

```go
app.Dialog.Info().
    SetTitle("Custom Icon Example").
    SetMessage("Using a custom icon").
    SetIcon(myIconBytes).
    Show()
```

### Привязка к окну

Привяжите диалоговое окно к определённому окну:

```go
dialog := app.Dialog.Question().
    SetMessage("Window-specific question").
    AttachToWindow(window)

dialog.AddButton("OK")
dialog.Show()
```

**Преимущества:**

- Диалоговое окно отображается поверх нужного окна
- Родительское окно недоступно, пока отображается диалоговое окно
- Более удобная работа с несколькими окнами

## Полные примеры

### Подтверждение необратимого действия

```go
func deleteFiles(app *application.App, paths []string) {
    // Confirm deletion
    message := fmt.Sprintf("Delete %d file(s)?", len(paths))
    if len(paths) == 1 {
        message = fmt.Sprintf("Delete %s?", filepath.Base(paths[0]))
    }

    dialog := app.Dialog.Question().
        SetTitle("Confirm Delete").
        SetMessage(message)

    deleteBtn := dialog.AddButton("Delete")
    deleteBtn.OnClick(func() {
        // Perform deletion
        var errs []error
        for _, path := range paths {
            if err := os.Remove(path); err != nil {
                errs = append(errs, err)
            }
        }

        // Show result
        if len(errs) > 0 {
            app.Dialog.Error().
                SetTitle("Delete Failed").
                SetMessage(fmt.Sprintf("Failed to delete %d file(s)", len(errs))).
                Show()
        } else {
            app.Dialog.Info().
                SetTitle("Delete Complete").
                SetMessage(fmt.Sprintf("Deleted %d file(s)", len(paths))).
                Show()
        }
    })

    cancelBtn := dialog.AddButton("Cancel")
    dialog.SetDefaultButton(cancelBtn)
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
}
```

### Подтверждение выхода

```go
func confirmQuit(app *application.App) {
    dialog := app.Dialog.Question().
        SetTitle("Quit").
        SetMessage("You have unsaved work. Are you sure you want to quit?")

    yes := dialog.AddButton("Yes")
    yes.OnClick(func() {
        app.Quit()
    })

    no := dialog.AddButton("No")
    dialog.SetDefaultButton(no)
    dialog.Show()
}
```

### Диалоговое окно обновления с возможностью скачивания

```go
func showUpdateDialog(app *application.App) {
    dialog := app.Dialog.Question().
        SetTitle("Update").
        SetMessage("A new version is available. The cancel button is selected when pressing escape.")

    download := dialog.AddButton("📥 Download")
    download.OnClick(func() {
        app.Dialog.Info().SetMessage("Downloading...").Show()
    })

    cancel := dialog.AddButton("Cancel")

    dialog.SetDefaultButton(download)
    dialog.SetCancelButton(cancel)
    dialog.Show()
}
```

### Вопрос с пользовательским значком

```go
func showCustomIconQuestion(app *application.App, iconBytes []byte) {
    dialog := app.Dialog.Question().
        SetTitle("Custom Icon Example").
        SetMessage("Using a custom icon").
        SetIcon(iconBytes)

    likeIt := dialog.AddButton("I like it!")
    likeIt.OnClick(func() {
        app.Dialog.Info().SetMessage("Thanks!").Show()
    })

    notKeen := dialog.AddButton("Not so keen...")
    notKeen.OnClick(func() {
        app.Dialog.Info().SetMessage("Too bad!").Show()
    })

    dialog.SetDefaultButton(likeIt)
    dialog.Show()
}
```

## Рекомендации

### ✅ Рекомендуется

- **Будьте конкретны** — «Файл сохранён в папке „Документы“» вместо «Успешно»
- **Выбирайте подходящий тип** — Error для ошибок, Warning для предупреждений
- **Предоставляйте контекст** — указывайте существенные подробности
- **Используйте понятные подписи кнопок** — «Удалить» вместо «ОК»
- **Задавайте безопасные значения по умолчанию** — выбирайте действие, не приводящее к потере данных
- **Обрабатывайте отмену** — пользователь может закрыть диалоговое окно

### ❌ Не делайте так

- **Не злоупотребляйте диалоговыми окнами** — они прерывают рабочий процесс
- **Не используйте их для частых обновлений** — вместо этого используйте уведомления
- **Не используйте общие сообщения** — сообщение «Ошибка» ничего не объясняет
- **Не игнорируйте ошибки** — обрабатывайте ошибки dialog.Show()
- **Не блокируйте выполнение без необходимости** — рассмотрите асинхронные альтернативы
- **Не используйте технический жаргон** — пишите простым языком

## Различия между платформами

### macOS

- При привязке к окну отображается в виде листа
- Стандартные сочетания клавиш (⌘. для отмены)
- Автоматически следует системной теме
- Встроенная поддержка специальных возможностей

### Windows

- Модальные диалоговые окна
- Оформление TaskDialog
- Esc для отмены
- Следует системной теме

### Linux

- Диалоговые окна GTK
- Внешний вид зависит от среды рабочего стола
- Следует теме рабочего стола
- Стандартная навигация с клавиатуры

## Что дальше

@cards{cols="2"}
📖 Диалоговые окна для работы с файлами
Открытие и сохранение файлов, а также выбор папок.

[Подробнее →](/features/dialogs/file/)

---
◆ Пользовательские диалоговые окна
Создавайте пользовательские диалоговые окна.

[Подробнее →](/features/dialogs/custom/)

---
● Уведомления
Ненавязчивые уведомления.

[Подробнее →](/features/notifications/overview/)

---
★ События
Используйте события для неблокирующего обмена данными.

[Подробнее →](/features/events/system/)

@end

---

**Есть вопросы?** Задайте их в [Discord](https://discord.gg/JDdSxwjhGf) или посмотрите [примеры диалоговых окон](https://github.com/wailsapp/wails/tree/master/v3/examples/dialogs).
