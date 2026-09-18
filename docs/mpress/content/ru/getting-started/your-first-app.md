---
title: "Ваше первое приложение"
description: "Пошаговое создание первого настольного приложения на Wails"
slug: "getting-started/your-first-app"
sourcePath: "getting-started/your-first-app.md"
---

В этом руководстве показано, как создать первое приложение на Wails v3: от настройки проекта до сборки и организации процесса разработки.

<br/>

<br/>

@steps
### Создание нового проекта
Откройте терминал и выполните следующую команду, чтобы создать новый проект Wails:

```bash
wails3 init -n myfirstapp
```

Эта команда создаст новый каталог `myfirstapp` со всеми необходимыми файлами.

   <video src="/assets/wails_init.mp4" controls></video>

### Знакомство со структурой проекта
Перейдите в каталог `myfirstapp`. В нём находятся следующие файлы и папки:

@filetree
- build/           Содержит файлы, используемые в процессе сборки
  - appicon.png  Значок приложения
  - config.yml   Конфигурация сборки
  - Taskfile.yml Build tasks
  - darwin/      Файлы сборки для macOS
    - Info.dev.plist Development configuration
    - Info.plist    Конфигурация для рабочей версии
    - Taskfile.yml  Задачи сборки для macOS
    - icons.icns    Значок приложения для macOS
  - linux/       Файлы сборки для Linux
    - Taskfile.yml  Задачи сборки для Linux
    - appimage/     Упаковка в AppImage
      - build.sh  Сценарий сборки AppImage
    - nfpm/        Упаковка с помощью NFPM
      - nfpm.yaml Package configuration
      - scripts/  Сценарии сборки
  - windows/     Файлы сборки для Windows
    - Taskfile.yml        Задачи сборки для Windows
    - icon.ico           Значок приложения для Windows
    - info.json          Метаданные приложения
    - wails.exe.manifest Windows manifest file
    - nsis/              Файлы установщика NSIS
      - project.nsi                    Файл проекта NSIS
      - wails_tools.nsh               Вспомогательные сценарии NSIS
- frontend/        Файлы клиентской части приложения
  - index.html   Основной HTML-файл
  - main.js      Основной файл JavaScript
  - package.json NPM package configuration
  - public/      Статические ресурсы
  - Inter Font License.txt Font license
- .gitignore      Файл исключений Git
- README.md       Документация проекта
- Taskfile.yml    Задачи проекта
- go.mod          Файл модуля Go
- go.sum          Контрольные суммы модулей Go
- greetservice.go Greeting service
- main.go         Основной код приложения
@end

Уделите немного времени изучению этих файлов и ознакомьтесь со структурой проекта.

@note{type="info"}
Хотя Wails v3 по умолчанию использует [Task](https://taskfile.dev/) в качестве системы сборки, вы можете использовать `make` или любую другую альтернативную систему сборки.

@end

### Сборка приложения
Чтобы собрать приложение, выполните:

```bash
wails3 build
```

Эта команда скомпилирует отладочную версию приложения и сохранит её в новом каталоге `bin`.

@note{type="info"}
`wails3 build` — это сокращение для `wails3 task build`; команда запустит задачу `build` в `Taskfile.yml`.

@end

     <video src="/assets/wails_build.mp4" controls></video>

После сборки приложение можно запустить так же, как любое обычное приложение:

@tabs{sync-key="platform"}
[Mac]
```sh
./bin/myfirstapp
```

[Windows]
```sh
bin\myfirstapp.exe
```

[Linux]
```sh
./bin/myfirstapp
```

@end

Вы увидите простой пользовательский интерфейс — отправную точку для вашего приложения. Поскольку это отладочная версия, в окне консоли также будут отображаться журналы. Они полезны при отладке.

### Режим разработки
Приложение также можно запустить в режиме разработки. В этом режиме можно изменять код клиентской части и сразу видеть изменения в работающем приложении, не пересобирая его целиком.

1. Откройте новое окно терминала.
2. Выполните `wails3 dev`. Приложение будет скомпилировано и запущено в режиме отладки.
3. Откройте `frontend/index.html` в выбранном редакторе.
4. Отредактируйте код, заменив `Please enter your name below` на `Please enter your name below!!!`.
5. Сохраните файл.

Это изменение немедленно отобразится в приложении.

Любые изменения в коде серверной части запустят повторную сборку:

1. Откройте `greetservice.go`.
2. В строке с `return "Hello " + name + "!"` замените это значение на `return "Hello there " + name + "!"`.
3. Сохраните файл.

Приложение обновится в течение нескольких секунд.

     <video src="/assets/wails_dev.mp4" controls></video>

### Упаковка приложения
Когда приложение будет готово к распространению, можно создать пакеты для соответствующих платформ:

@tabs{sync-key="platform"}
[Mac]
Чтобы создать пакет `.app`:

```bash
wails3 package
```

Будет создана сборка для выпуска, упакованная в пакет `.app` в каталоге `bin`.

[Windows]
Чтобы создать установщик NSIS:

```bash
wails3 package
```

Будет создана сборка для выпуска, упакованная в установщик NSIS в каталоге `bin`.

[Linux]
Wails поддерживает несколько форматов пакетов для распространения в Linux:

```bash
# Create all package types (AppImage, deb, rpm, and Arch Linux)
wails3 package

# Or create specific package types
wails3 task linux:create:appimage  # AppImage format
wails3 task linux:create:deb       # Debian package
wails3 task linux:create:rpm       # Red Hat package
wails3 task linux:create:aur       # Arch Linux package
```

@end

Подробные сведения о вариантах и настройке упаковки см. в нашем [руководстве по сборке и упаковке](/guides/build/building/).

### Настройка системы контроля версий и имени модуля
Проект создаётся с временным именем модуля `changeme`. Рекомендуется заменить его именем, соответствующим URL-адресу репозитория:

1. Создайте новый репозиторий на GitHub (или на предпочитаемом хостинге Git)
2. Инициализируйте Git в каталоге проекта:
  ```bash
  git init
  git add .
  git commit -m "Initial commit"
  ```

3. Укажите удалённый репозиторий (замените значение URL-адресом своего репозитория):
  ```bash
  git remote add origin https://github.com/username/myfirstapp.git
  ```

4. Измените имя модуля в `go.mod`, чтобы оно соответствовало URL-адресу репозитория:
  ```bash
  go mod edit -module github.com/username/myfirstapp
  ```

5. Отправьте код в удалённый репозиторий:
  ```bash
  git push -u origin main
  ```


Это обеспечивает соответствие имени модуля Go правилам именования модулей Go и упрощает совместное использование кода.

@note{type="tip" title="Полезный совет"}
Все этапы инициализации можно автоматизировать с помощью флага `-git` при создании проекта:

```bash
wails3 init -n myfirstapp -git github.com/username/myfirstapp
```

Поддерживаются различные форматы URL-адресов Git:

- HTTPS: `https://github.com/username/project`
- SSH: `git@github.com:username/project` или `ssh://git@github.com/username/project`
- Протокол Git: `git://github.com/username/project`
- Файловая система: `file:///path/to/project.git`

@end

@end

## Поздравляем!

Вы только что создали, разработали и упаковали своё первое приложение на Wails. И это лишь начало того, чего можно достичь с помощью Wails v3.

## Дальнейшие шаги

Если вы только начинаете работать с Wails, рекомендуем далее ознакомиться с нашими учебными руководствами — это практическое знакомство с различными возможностями Wails. Первое руководство — [«Создание сервиса»](/tutorials/01-creating-a-service/).

Если вы уже являетесь опытным пользователем, ознакомьтесь с [руководством по сборке и упаковке](/guides/build/building/), чтобы получить более подробную информацию об использовании Wails.
