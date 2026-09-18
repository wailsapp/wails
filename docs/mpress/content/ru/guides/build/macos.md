---
title: "Упаковка для macOS"
description: "Упакуйте приложение Wails для распространения в macOS"
slug: "guides/build/macos"
sourcePath: "guides/build/macos.md"
---

## Закрытые API macOS

По умолчанию Wails v3 использует общедоступные API macOS. Чтобы включить функции, которым требуются недокументированные API Apple, соберите приложение с единственным тегом сборки Go `private_mac_apis`:

```bash
wails3 build -tags private_mac_apis
EXTRA_TAGS=private_mac_apis wails3 dev
wails3 package GOOS=darwin EXTRA_TAGS=private_mac_apis
```

При сборке непосредственно с помощью Go используйте `go build -tags private_mac_apis .` (или `-tags production,private_mac_apis` для рабочей версии). В существующих приложениях, которые зависят от закрытых функций, необходимо добавить этот тег, чтобы сохранить их работу. Тег применяется только к сборкам настольных приложений для macOS.

Полный перечень затрагиваемых функций и значений параметров, точные резервные варианты для сборок с общедоступными API, сопоставления стилей Liquid Glass и комбинации сборок с инспектором см. в разделе [Закрытые API macOS](/guides/build/private-macos-apis/). Без этого тега операции, доступные только через закрытые API, ничего не делают; общедоступный API Go остаётся неизменным.

## Пакет приложения

Упакуйте приложение как стандартный пакет macOS `.app`:

```bash
wails3 package GOOS=darwin
```

Будет создан `bin/<AppName>.app` со следующим содержимым:

- Скомпилированный двоичный файл в `Contents/MacOS/`
- Значок приложения в `Contents/Resources/` (из `icons.icns` или, при наличии, из каталога ресурсов `Assets.car`)
- `Info.plist` с метаданными приложения

## Ресурсы пакета

`Contents/Resources/` — стандартное место для файлов только для чтения, поставляемых с приложением macOS. Используйте его для крупных шаблонов, начальных данных, мультимедийных файлов, языковых пакетов и других данных, которые следует открывать по запросу, а не компилировать в исполняемый файл Go с помощью `embed`.

Wails уже помещает значок приложения в этот каталог. Чтобы добавить собственные файлы, поместите их в исходный каталог, например `build/resources/`, а затем добавьте шаг копирования в задачу `create:app:bundle` в `build/darwin/Taskfile.yml`:

```yaml
tasks:
  create:app:bundle:
    cmds:
      # Existing bundle creation commands...
      - |
          if [ -d build/resources ]; then
            cp -R build/resources/. "{{.BIN_DIR}}/{{.APP_NAME}}.app/Contents/Resources/"
          fi
```

Если вы используете задачу `darwin:run` из Taskfile, добавьте эквивалентную команду в её задачу `run`, указав в качестве целевого расположения `{{.BIN_DIR}}/{{.APP_NAME}}.dev.app/Contents/Resources/`.

### Чтение ресурсов из Go

Импортируйте пакет платформы macOS:

```go
import (
	"io/fs"

	"github.com/wailsapp/wails/v3/pkg/mac"
)
```

Для небольших файлов используйте `LoadResource`:

```go
func loadSplash() ([]byte, error) {
	return mac.LoadResource("images/splash.png")
}
```

Для крупных файлов используйте `ResourceFS`. Эта функция возвращает `io/fs.FS` с корнем в `Contents/Resources`, поэтому вызывающий код может открыть ресурс и читать его в потоковом режиме, не загружая предварительно весь файл в срез байтов Go:

```go
func openCatalogue() (fs.File, error) {
	resources, err := mac.ResourceFS()
	if err != nil {
		return nil, err
	}

	return resources.Open("catalogue/defaults.json")
}
```

Имена ресурсов представляют собой пути с разделителями-косыми чертами относительно `Contents/Resources`. `ResourceFS` и `LoadResource` возвращают `mac.ErrNotInAppBundle`, если исполняемый файл не запущен из `.app/Contents/MacOS`.

Считайте ресурсы пакета неизменяемыми. Изменение файлов внутри подписанного приложения делает его подпись кода недействительной; загруженные, созданные или доступные пользователю для редактирования данные следует хранить в пользовательском каталоге Application Support.

### Универсальный двоичный файл

Соберите приложение как для компьютеров Mac с Apple Silicon, так и для компьютеров Mac с процессорами Intel:

```bash
wails3 task darwin:package:universal
```

Будет создан единый `.app`, который выполняется нативно на обеих архитектурах. Универсальные двоичные файлы можно собирать на любой платформе — в Linux и Windows автоматически используется `wails3 tool lipo`.

## Настройка пакета

Отредактируйте `build/darwin/Info.plist`, чтобы настроить:

- Идентификатор пакета (`CFBundleIdentifier`)
- Название и версию приложения
- Минимальную версию macOS
- Ассоциации файлов
- Схемы URL

Значок приложения создаётся из ресурсов в каталоге `build/`. Используйте задачу `generate:icons`:

```bash
wails3 task common:generate:icons
```

При этом `build/appicon.png` используется для создания `darwin/icons.icns` и `windows/icon.ico`. В macOS также можно предоставить `build/appicon.icon` (формат Icon Composer): задача передаёт `-iconcomposerinput appicon.icon -macassetdir darwin`, в результате чего из файла `.icon` создаются `Assets.car` и `darwin/icons.icns` (на платформах, отличных от macOS, этот шаг пропускается). При наличии `Assets.car` запустите задачу `update:build-assets`, чтобы соответствующим образом обновить `Info.plist` и `CFBundleIconName`:

```bash
wails3 task common:update:build-assets
```

Чтобы запустить команду создания значка вручную из каталога `build/`:

```bash
cd build
wails3 generate icons -input appicon.png -macfilename darwin/icons.icns -windowsfilename windows/icon.ico -iconcomposerinput appicon.icon -macassetdir darwin
```

## Подписание кода

Подпишите приложение для распространения:

```bash
# Using the wrapper (auto-detects platform)
wails3 sign GOOS=darwin

# Or using the task directly
wails3 task darwin:sign
```

Настройте подписание в `build/darwin/Taskfile.yml`:

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  ENTITLEMENTS: "build/darwin/entitlements.plist"
```

### Нотаризация

Apple требует нотариально заверять приложения, распространяемые за пределами Mac App Store:

```bash
wails3 task darwin:sign:notarize
```

Сначала сохраните свои учётные данные. Запустите интерактивный мастер (`wails3 setup signing`) или вызовите `notarytool` напрямую:

```bash
xcrun notarytool store-credentials "my-notarize-profile" \
  --apple-id "you@email.com" \
  --team-id "TEAMID" \
  --password "app-specific-password"
```

Выполните настройку в `build/darwin/Taskfile.yml`:

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
```

Подробности см. в разделе [Подписание приложений](/guides/build/signing/).

## Установщик DMG

Поставляемый с Wails 3 шаблон предоставляет `wails3 task darwin:package:dmg`. Сначала он создаёт `.app`, а затем с помощью библиотеки DMG формирует оформленный образ DMG. По умолчанию в DMG используется градиентный фон в фирменном стиле Wails с красным драконом и словесным логотипом WAILS.

```bash
wails3 task darwin:package:dmg
```

Низкоуровневая задача `darwin:create:dmg` создаёт DMG из существующего пакета `.app`; её можно настроить непосредственно в Taskfile:

```yaml
vars:
  # These are the template defaults; override them when needed.
  DMG_BACKGROUND: build/darwin/dmg-background.png
  DMG_VOLUME_ICON: build/darwin/icons.icns
  DMG_FILE_ICON: build/darwin/dmg-file-icon.icns
  DMG_WINDOW_WIDTH: 540
  DMG_WINDOW_HEIGHT: 380
  DMG_FILES: "Install.command=build/darwin/Install.command,README.txt=README.md"
```

### Макет по умолчанию

Созданный DMG содержит:

- Пакет приложения слева
- Ссылку `Applications` справа
- Окно Finder размером 540 × 380 пикселей
- Значки размером 96 пунктов с подписями под каждым значком
- Фирменный фон Wails из `build/darwin/dmg-background.png`

Значки приложения и `Applications` располагаются относительно заданных размеров окна, поэтому при изменении `DMG_WINDOW_WIDTH` или `DMG_WINDOW_HEIGHT` пропорциональные интервалы в стандартном макете с двумя значками сохраняются. Для наилучшего результата используйте фоновое изображение с теми же размерами в пикселях, что и окно Finder.

### Замена ресурсов DMG

Созданные файлы в каталоге `build/darwin/` являются обычными ресурсами проекта, и их можно заменить:

- `DMG_BACKGROUND` задаёт изображение, отображаемое за содержимым окна Finder.
- `DMG_VOLUME_ICON` задаёт значок подключённого тома.
- `DMG_FILE_ICON` задаёт значок, отображаемый в Finder для итогового файла `.dmg`.

Значок тома и значок файла DMG — это отдельные ресурсы. Замена значка приложения не приводит к автоматической замене ни одного из них.

### Добавление дополнительных файлов

Используйте `DMG_FILES`, чтобы поместить рядом с приложением сценарии установки, примечания к выпуску, лицензии или другие ресурсы. Значение представляет собой разделённый запятыми список пар `name=path`:

```yaml
vars:
  DMG_FILES: "Install.command=build/darwin/Install.command,README.txt=README.md"
```

Имя перед `=` — это имя файла, отображаемое внутри DMG. Путь после `=` указывает на исходный файл в проекте. Начальные и конечные пробелы игнорируются.

Каждое отображаемое имя должно быть уникальным. Дополнительные файлы не могут заменять элементы, уже созданные упаковщиком, в том числе пакет приложения или элемент `Applications`. При конфликте имён упаковка завершается с ошибкой вместо создания повреждённого DMG.

@note{type="note"}
Создание DMG поддерживается только в macOS, поскольку для него используются средства macOS для работы с образами дисков и Finder. Пакеты `.app` можно создавать путём кросс-компиляции на других платформах, но итоговый DMG необходимо создать на Mac.

@end

## Устранение неполадок

### «Приложение повреждено, и его не удаётся открыть»

Приложение не подписано. Подпишите его сертификатом Developer ID, либо пользователи могут обойти Gatekeeper:

```bash
xattr -cr /path/to/YourApp.app
```

### Нотаризация завершается с ошибкой

Распространённые проблемы:

- **Недействительные учётные данные**: повторно выполните `xcrun notarytool store-credentials` (или `wails3 setup signing`)
- **Требуется усиленная среда выполнения**: при необходимости убедитесь, что права включают `com.apple.security.cs.allow-unsigned-executable-memory`
- **Отсутствует метка времени**: при подписании метка времени должна добавляться автоматически

### Приложение, созданное путём кросс-компиляции, не запускается

Исполняемые файлы macOS, созданные путём кросс-компиляции, не подписаны. Перенесите их на Mac и подпишите перед тестированием:

```bash
codesign --force --deep --sign - YourApp.app
```
