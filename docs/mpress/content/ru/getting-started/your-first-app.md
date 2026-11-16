---
title: "Ваше первое приложение"
description: "Пошаговое создание первого настольного приложения Wails"
slug: "getting-started/your-first-app"
sourcePath: "getting-started/your-first-app.md"
---

В этом руководстве показано, как создать первое приложение Wails v3: от настройки проекта до сборки и организации процесса разработки.

<br/>

<br/>

@steps
### Создание нового проекта
Откройте терминал и выполните следующую команду, чтобы создать новый проект Wails:

```bash
wails3 init -n myfirstapp -t vanilla-js
```

В этом уроке явно выбирается шаблон JavaScript Vanilla. Базовый шаблон `vanilla` использует TypeScript. Команда создаёт `myfirstapp` со всеми необходимыми файлами.

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
    - Info.plist    Конфигурация для рабочей среды
    - Taskfile.yml  Задачи сборки для macOS
    - icons.icns    Значок приложения для macOS
  - linux/       Файлы сборки для Linux
    - Taskfile.yml  Задачи сборки для Linux
    - appimage/     Создание пакета AppImage
      - build.sh  Скрипт сборки AppImage
    - nfpm/        Создание пакетов с помощью NFPM
      - nfpm.yaml Package configuration
      - scripts/  Скрипты сборки
  - windows/     Файлы сборки для Windows
    - Taskfile.yml        Задачи сборки для Windows
    - icon.ico           Значок приложения для Windows
    - info.json          Метаданные приложения
    - wails.exe.manifest Windows manifest file
    - nsis/              Файлы установщика NSIS
      - project.nsi                    Файл проекта NSIS
      - wails_tools.nsh               Вспомогательные скрипты NSIS
- frontend/        Файлы клиентского приложения
  - index.html   Основной HTML-файл
  - src/         Исходный код клиентской части
    - main.js    Основной файл JavaScript
  - package.json NPM package configuration
  - public/      Статические ресурсы
  - Inter Font License.txt Font license
- .gitignore      Файл исключений Git
- README.md       Документация проекта
- Taskfile.yml    Задачи проекта
- go.mod          Файл модуля Go
- go.sum          Контрольные суммы модуля Go
- greetservice.go Greeting service
- main.go         Основной код приложения
@end

Уделите немного времени изучению этих файлов и структуры проекта.

@note{type="info"}
Хотя Wails v3 по умолчанию использует [Task](https://taskfile.dev/) в качестве системы сборки, ничто не мешает вам использовать `make` или любую другую альтернативную систему сборки.

@end

### Сборка приложения
Чтобы собрать приложение, выполните:

```bash
wails3 build
```

При использовании сгенерированных Taskfile эта команда создаёт рабочую сборку в `bin`. Для запуска в режиме разработки используйте `wails3 dev`, а для явного создания отладочной сборки — `wails3 build DEV=true`.

@note{type="info"}
`wails3 build` — это сокращение для `wails3 task build`, которое запускает задачу `build` в `Taskfile.yml`.

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

Вы увидите начальный интерфейс. В рабочей сборке не гарантируется ожидаемая при разработке работа инспектора или консоли; для отладки используйте режим разработки.

### Режим разработки
Приложение также можно запустить в режиме разработки. В этом режиме изменения в коде клиентской части сразу отображаются в работающем приложении без необходимости пересобирать его целиком.

1. Откройте новое окно терминала.
2. Выполните `wails3 dev`. Приложение будет скомпилировано и запущено в режиме отладки.
3. Откройте `frontend/index.html` в предпочитаемом редакторе.
4. Измените `placeholder="Your name"` поля ввода имени на `placeholder="Your full name"` и сохраните файл.
5. Сохраните файл.

Это изменение немедленно отобразится в приложении.

Любые изменения в коде серверной части запустят повторную сборку:

1. Откройте `greetservice.go`.
2. Замените строку с `return "Hello " + name + "!"` на `return "Hello there " + name + "!"`.
3. Сохраните файл.

Приложение обновится в течение нескольких секунд.

     <video src="/assets/wails_dev.mp4" controls></video>

### Упаковка приложения
Когда приложение будет готово к распространению, вы сможете создать пакеты для конкретных платформ:

@tabs{sync-key="platform"}
[Mac]
Чтобы создать пакет `.app`:

```bash
wails3 package
```

Будет создана рабочая сборка, упакованная в пакет `.app` в каталоге `bin`.

[Windows]
Чтобы создать установщик NSIS:

```bash
wails3 package
```

Будет создана рабочая сборка, упакованная в установщик NSIS в каталоге `bin`.

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

Подробные сведения о параметрах и конфигурации упаковки см. в нашем [руководстве по сборке и упаковке](/guides/build/building/).

### Настройка системы управления версиями и имени модуля
Проект создаётся с временным именем модуля `changeme`. Рекомендуется заменить его в соответствии с URL-адресом репозитория:

1. Создайте новый репозиторий на GitHub или на предпочитаемом Git-хостинге
2. Инициализируйте Git в каталоге проекта:
  ```bash
  git init
  git add .
  git commit -m "Initial commit"
  ```

3. Задайте удалённый репозиторий, подставив URL-адрес своего репозитория:
  ```bash
  git remote add origin https://github.com/username/myfirstapp.git
  ```

4. Измените имя модуля в `go.mod` в соответствии с URL-адресом репозитория:
  ```bash
  go mod edit -module github.com/username/myfirstapp
  ```

5. Отправьте код в удалённый репозиторий:
  ```bash
  git push -u origin main
  ```


Так имя модуля Go будет соответствовать соглашениям об именовании модулей Go, а делиться кодом станет проще.

@note{type="tip" title="Полезный совет"}
Все этапы начальной настройки можно автоматизировать с помощью флага `-git` при создании проекта:

```bash
wails3 init -n myfirstapp -t vanilla-js -git github.com/username/myfirstapp
```

Поддерживаются различные форматы URL-адресов Git:

- HTTPS: `https://github.com/username/project`
- SSH: `git@github.com:username/project` или `ssh://git@github.com/username/project`
- Протокол Git: `git://github.com/username/project`
- Файловая система: `file:///path/to/project.git`

@end

@end

## Поздравляем!

Вы только что создали, разработали и упаковали своё первое приложение Wails. И это лишь начало того, чего можно достичь с помощью Wails v3.

## Следующие шаги

Если вы только начинаете работать с Wails, рекомендуем далее ознакомиться с нашими руководствами, которые помогут на практике изучить различные возможности Wails. Первое руководство — [«Создание сервиса»](/tutorials/01-creating-a-service/).

Если вы опытный пользователь, ознакомьтесь с [руководством по сборке и упаковке](/guides/build/building/), чтобы получить более подробную информацию об использовании Wails.
