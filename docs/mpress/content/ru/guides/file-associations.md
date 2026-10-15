---
title: "Ассоциации файлов"
description: "Настройка ассоциаций файлов для приложения Wails"
slug: "guides/file-associations"
sourcePath: "guides/file-associations.md"
---

Поддерживаемые платформы: <span class="mpress-badge mpress-badge-note">Windows</span> <span class="mpress-badge mpress-badge-success">macOS</span>

<br/>

Ассоциации файлов позволяют приложению обрабатывать файлы определённых типов, когда пользователи открывают их. Это особенно полезно для текстовых редакторов, программ просмотра изображений и любых других приложений, работающих с определёнными форматами файлов. В этом руководстве объясняется, как реализовать ассоциации файлов в приложении Wails v3.

## Обзор

В настоящее время Wails v3 поддерживает ассоциации файлов для следующих платформ:

- Windows (установочные пакеты NSIS)
- macOS (пакеты приложений)

## Настройка

Ассоциации файлов настраиваются в файле `config.yml`, расположенном в каталоге `build` вашего проекта.

### Базовая настройка

Чтобы настроить ассоциации файлов:

1. Откройте `build/config.yml`
2. Добавьте ассоциации файлов в раздел `fileAssociations`
3. Выполните `wails3 update build-assets`, чтобы обновить ресурсы сборки
4. Задайте поле `FileAssociations` в параметрах приложения
5. Упакуйте приложение с помощью `wails3 package`

Пример конфигурации:

```yaml
fileAssociations:
  - ext: myapp
    name: MyApp Document
    description: MyApp Document File
    iconName: myappFileIcon
    role: Editor
  - ext: custom
    name: Custom Format
    description: Custom File Format
    iconName: customFileIcon
    role: Editor
```

### Свойства конфигурации

| Свойство | Описание | Платформа |
| --- | --- | --- |
| ext | Расширение файла без начальной точки (например, `txt`) | Все |
| name | Отображаемое имя типа файла | Все |
| description | Описание, отображаемое в свойствах файла | Windows |
| iconName | Имя файла значка без расширения в папке сборки | Все |
| role | Роль приложения для этого типа файлов (например, `Editor`, `Viewer`) | macOS |
| mimeType | MIME-тип файла (например, `image/jpeg`) | macOS |

## Обработка событий открытия файлов

Чтобы обрабатывать события открытия файлов в приложении, можно отслеживать событие `events.Common.ApplicationOpenedWithFile`:

```go {title="main.go"}
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

func main() {
    app := application.New(application.Options{
        Name: "MyApp",
        FileAssociations: []string{".txt", ".md"}, // Specify supported extensions
    })

    // Listen for files being used to open the application
    app.Event.OnApplicationEvent(events.Common.ApplicationOpenedWithFile, func(event *application.ApplicationEvent) {
        associatedFile := event.Context().Filename()
        app.Dialog.Info().SetMessage("Application opened with file: " + associatedFile).Show()
    })

    // Create your window and run the app...
}

```

## Пошаговое руководство

Рассмотрим настройку ассоциаций файлов для простого текстового редактора:

@steps
### Создание значков
- Создайте значки для своего типа файлов (рекомендуемые размеры: 16x16, 32x32, 48x48, 256x256)
- Сохраните значки в папке `build` вашего проекта
- Назовите их в соответствии с конфигурацией `iconName` (например, `textFileIcon.png`)

@note{type="tip"}
Для создания необходимых значков можно использовать `wails3 generate icons`. Чтобы получить дополнительную информацию, выполните `wails3 generate icons --help`.

@end

- Для macOS добавьте в задачу `create:app:bundle:` команду копирования, например `cp build/darwin/documenticon.icns {{.BIN_DIR}}/{{.APP_NAME}}.app/Contents/Resources`.

### Настройка ассоциаций файлов
Измените файл `build/config.yml`, добавив в него ассоциации файлов:

```yaml
# build/config.yml
fileAssociations:
  - ext: txt
    name: Text Document
    description: Plain Text Document
    iconName: textFileIcon
    role: Editor
```

### Обновление ресурсов сборки
Чтобы обновить ресурсы сборки, выполните следующую команду:

```bash
wails3 update build-assets
```

### Задание ассоциаций файлов в параметрах приложения
В файле `main.go` задайте поле `FileAssociations` в параметрах приложения:

```go
app := application.New(application.Options{
  Name: "MyApp",
  FileAssociations: []string{".txt", ".md"}, // Specify supported extensions
})
```

@note{type="tip" title="Почему расширения файлов необходимо указывать и в конфигурации приложения, и в config.yml?"}
Когда в Windows файл открывают с помощью ассоциации, приложение запускается с именем файла в качестве первого аргумента. Приложение не может определить, является ли первый аргумент файлом или аргументом командной строки, поэтому оно использует поле `FileAssociations` в параметрах приложения, чтобы определить, является ли первый аргумент ассоциированным файлом.

@end

### Упаковка приложения
Упакуйте приложение с помощью следующей команды:

```bash
wails3 package
```

Упакованное приложение будет создано в каталоге `bin`. После этого приложение можно установить и протестировать.

## Дополнительные примечания

- Значки следует размещать в папке сборки в формате PNG
- Для тестирования ассоциаций файлов необходимо установить упакованное приложение

@end
