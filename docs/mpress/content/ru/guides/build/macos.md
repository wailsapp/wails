---
title: "Упаковка для macOS"
description: "Упакуйте приложение Wails для распространения в macOS"
slug: "guides/build/macos"
sourcePath: "guides/build/macos.md"
---

## Закрытые API macOS

По умолчанию Wails v3 использует общедоступные API macOS. Чтобы включить функции, для которых требуются недокументированные API Apple, соберите приложение с единственным тегом сборки Go `private_mac_apis`:

```bash
wails3 build -tags private_mac_apis
EXTRA_TAGS=private_mac_apis wails3 dev
wails3 package GOOS=darwin EXTRA_TAGS=private_mac_apis
```

При непосредственной сборке средствами Go используйте `go build -tags private_mac_apis .` (или `-tags production,private_mac_apis` для рабочей версии). В существующие приложения, зависящие от закрытых функций, необходимо добавить этот тег, чтобы сохранить их поведение. Тег применяется только к сборкам настольных приложений для macOS.

Полный перечень затронутых функций и значений параметров, точные резервные варианты поведения общедоступной сборки, сопоставления стилей Liquid Glass и комбинации сборок инспектора приведены в разделе [Закрытые API macOS](/guides/build/private-macos-apis/). Без этого тега операции, доступные только через закрытые API, ничего не делают; общедоступный API Go остаётся неизменным.

## Пакет приложения

Упакуйте приложение как стандартный пакет macOS `.app`:

```bash
wails3 package GOOS=darwin
```

Будет создан `bin/<AppName>.app` со следующим содержимым:

- Скомпилированный двоичный файл в `Contents/MacOS/`
- Значок приложения в `Contents/Resources/` (из `icons.icns` или, при наличии, из каталога ресурсов `Assets.car`)
- `Info.plist` с метаданными приложения

### Проекты в синхронизируемых папках

Некоторые поставщики файлов в macOS добавляют метаданные Finder в каталоги `.app`. Из-за этого `codesign` отклоняет пакет либо его подпись становится недействительной после подписания. Простого удаления метаданных недостаточно, если поставщик добавляет их повторно.

Сгенерированные задачи сборки для macOS проверяют наличие такого поведения. При его обнаружении Wails перемещает сгенерированный выходной каталог проекта в `~/Library/Application Support/Wails/build-output/<project>-<hash>` и заменяет исходный выходной каталог символической ссылкой. Существующие выходные файлы сохраняются, а `bin/<AppName>.app` остаётся путём, используемым для упаковки, подписания и режима разработки. При перемещении CLI выводит расположение локального хранилища.

Чтобы выбрать другое локальное несинхронизируемое хранилище, задайте `WAILS_MACOS_OUTPUT_ROOT`. Корневой каталог для перемещения должен находиться в той же файловой системе, что и проект. Явно заданный `BIN_DIR` за пределами проекта используется в соответствии с конфигурацией. Чтобы задействовать этот механизм в существующих проектах, требуется обновлённый `build/darwin/Taskfile.yml`.

Для распространения скопируйте `bin/<AppName>.app` или создайте DMG. При копировании только символической ссылки `bin` артефакты сборки не включаются. Команда непосредственного запуска `wails3 tool sign` сообщает об ошибке, если к её входной папке повторно прикрепляются метаданные; перед подписанием переместите этот пакет в несинхронизируемую папку.

## Ресурсы пакета

`Contents/Resources/` — стандартное расположение для доступных только для чтения файлов, поставляемых с приложением macOS. Используйте его для крупных шаблонов, начальных данных, медиафайлов, языковых пакетов и других данных, которые следует открывать по запросу, а не компилировать в исполняемый файл Go с помощью `embed`.

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

Если вы используете задачу `darwin:run` из Taskfile, добавьте эквивалентную команду в её задачу `run`, указав в качестве целевого пути `{{.BIN_DIR}}/{{.APP_NAME}}.dev.app/Contents/Resources/`.

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

Для крупных файлов используйте `ResourceFS`. Эта функция возвращает `io/fs.FS` с корнем в `Contents/Resources`, поэтому вызывающий код может открыть ресурс и читать его как поток, не загружая предварительно целиком в срез байтов Go:

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

Считайте ресурсы пакета неизменяемыми. Изменение файлов внутри подписанного приложения делает его цифровую подпись недействительной; загруженные, сгенерированные или редактируемые пользователем данные храните вместо этого в пользовательском каталоге Application Support.

### Универсальный двоичный файл

Соберите приложение как для компьютеров Mac с Apple Silicon, так и для компьютеров Mac с процессорами Intel:

```bash
wails3 task darwin:package:universal
```

Будет создан единый `.app`, который выполняется нативно на обеих архитектурах. Универсальные двоичные файлы можно собирать на любой платформе — в Linux и Windows `wails3 tool lipo` используется автоматически.

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

Она использует `build/appicon.png` для создания `darwin/icons.icns` и `windows/icon.ico`. В macOS также можно предоставить `build/appicon.icon` (формат Icon Composer): задача передаёт `-iconcomposerinput appicon.icon -macassetdir darwin`, в результате чего из файла `.icon` создаются `Assets.car` и `darwin/icons.icns` (на платформах, отличных от macOS, этот шаг пропускается). Если имеется `Assets.car`, запустите задачу `update:build-assets`, чтобы соответствующим образом обновить `Info.plist` и `CFBundleIconName`:

```bash
wails3 task common:update:build-assets
```

Чтобы выполнить команду создания значка вручную из каталога `build/`:

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

Apple требует нотариализовать приложения, распространяемые вне Mac App Store:

```bash
wails3 task darwin:sign:notarize
```

Сначала сохраните учётные данные. Запустите интерактивный мастер (`wails3 signing`) либо вызовите `notarytool` напрямую:

```bash
xcrun notarytool store-credentials "my-notarize-profile" \
  --apple-id "you@email.com" \
  --team-id "TEAMID" \
  --password "app-specific-password"
```

Настройте в `build/darwin/Taskfile.yml`:

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
```

Подробности приведены в разделе [Подписание приложений](/guides/build/signing/).

## Установщик DMG

Поставляемый шаблон Wails 3 предоставляет `wails3 task darwin:package:dmg`. Сначала он создаёт `.app`, а затем с помощью библиотеки DMG собирает оформленный DMG. По умолчанию в DMG используется фирменный градиентный фон Wails с красным символом дракона и логотипом WAILS.

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
- Окно Finder размером 540×380 пикселей
- Значки размером 96 пунктов с подписями под каждым значком
- Фирменный фон Wails из `build/darwin/dmg-background.png`

Значки приложения и `Applications` располагаются относительно настроенных размеров окна, поэтому при изменении `DMG_WINDOW_WIDTH` или `DMG_WINDOW_HEIGHT` в стандартном макете с двумя значками сохраняются пропорциональные интервалы. Для наилучшего результата используйте фоновое изображение с теми же размерами в пикселях, что и окно Finder.

### Замена ресурсов DMG

Созданные файлы в `build/darwin/` являются обычными ресурсами проекта, и их можно заменить:

- `DMG_BACKGROUND` определяет изображение, отображаемое за содержимым окна Finder.
- `DMG_VOLUME_ICON` определяет значок подключённого тома.
- `DMG_FILE_ICON` определяет значок, отображаемый в Finder для созданного файла `.dmg`.

Значок тома и значок файла DMG — это отдельные ресурсы. Замена значка приложения не приводит к автоматической замене ни одного из них.

### Добавление дополнительных файлов

Используйте `DMG_FILES`, чтобы добавить рядом с приложением сценарии установки, примечания к выпуску, лицензии или другие ресурсы. Значение представляет собой разделённый запятыми список пар `name=path`:

```yaml
vars:
  DMG_FILES: "Install.command=build/darwin/Install.command,README.txt=README.md"
```

Имя перед `=` — это имя файла, отображаемое внутри DMG. Путь после `=` указывает на исходный файл в проекте. Начальные и конечные пробелы игнорируются.

Каждое отображаемое имя должно быть уникальным. Дополнительные файлы не могут заменять элементы, уже созданные упаковщиком, включая пакет приложения или элемент `Applications`. При конфликте имён упаковка завершается ошибкой, а повреждённый DMG не создаётся.

@note{type="note"}
Создание DMG поддерживается только в macOS, поскольку для него используются средства macOS для работы с образами дисков и Finder. Пакеты `.app` можно создавать путём кросс-компиляции в других системах, но итоговый DMG необходимо создать на компьютере Mac.

@end

## Устранение неполадок

### «Приложение повреждено, и его не удаётся открыть»

Приложение не подписано. Подпишите его сертификатом Developer ID либо предложите пользователям обойти Gatekeeper:

```bash
xattr -cr /path/to/YourApp.app
```

### Нотаризация завершается ошибкой

Распространённые проблемы:

- **Недействительные учётные данные**: повторно выполните `xcrun notarytool store-credentials` (или `wails3 signing`)
- **Требуется усиленная среда выполнения**: при необходимости убедитесь, что файл разрешений содержит `com.apple.security.cs.allow-unsigned-executable-memory`
- **Отсутствует метка времени**: процесс подписания должен автоматически добавлять метку времени

### Кросс-компилированное приложение не запускается

Кросс-компилированные двоичные файлы macOS не подписаны. Перед тестированием перенесите их на компьютер Mac и подпишите:

```bash
codesign --force --deep --sign - YourApp.app
```
