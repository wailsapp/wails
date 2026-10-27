---
title: "Упаковка для Linux"
description: "Упаковка приложения Wails для распространения в Linux"
slug: "guides/build/linux"
sourcePath: "guides/build/linux.md"
---

## Форматы пакетов

Упакуйте приложение для распространения в Linux:

```bash
wails3 package GOOS=linux
```

При этом в каталоге `bin/` создаются пакеты нескольких форматов:

- **AppImage**: переносимый формат, работающий в любом дистрибутиве Linux
- **DEB**: для Debian, Ubuntu и производных дистрибутивов
- **RPM**: для Fedora, RHEL и производных дистрибутивов
- **Arch**: для Arch Linux и производных дистрибутивов

### Отдельные форматы

Соберите пакеты нужных форматов:

```bash
wails3 task linux:create:appimage
wails3 task linux:create:deb
wails3 task linux:create:rpm
wails3 task linux:create:aur
```

## Настройка пакетов

### Запись рабочего стола

Файл `.desktop` определяет, как приложение отображается в меню приложений. Он создаётся на основе значений из `build/linux/Taskfile.yml`:

```yaml
vars:
  APP_NAME: 'MyApp'
  EXEC: 'MyApp'
  ICON: 'MyApp'
  CATEGORIES: 'Development;'
```

### Метаданные пакета

Чтобы настроить пакеты DEB и RPM, отредактируйте `build/linux/nfpm/nfpm.yaml`:

```yaml
name: myapp
version: 1.0.0
maintainer: Your Name <you@example.com>
description: My awesome Wails application
homepage: https://example.com
license: MIT
```

### AppImage

Конфигурация AppImage находится в `build/linux/appimage/`. Значок приложения берётся из `build/appicon.png`.

В AppImage включены WebKit и его вспомогательные процессы, а встроенная библиотека исправлена так, чтобы запускать их из AppImage. Сценарий `AppRun` запускает приложение из `usr/lib` внутри AppImage. Исправленный путь разрешается как за пределами песочницы bubblewrap WebKit, так и внутри неё, поэтому песочница остаётся включённой. Используйте собственный `AppRun`, только если он сохраняет этот рабочий каталог.

## Подписание пакетов

Подпишите пакеты DEB и RPM ключом PGP:

```bash
# Using the wrapper (auto-detects platform)
wails3 sign GOOS=linux

# Or using tasks directly
wails3 task linux:sign:deb
wails3 task linux:sign:rpm
wails3 task linux:sign:packages  # Both
```

Настройте подписание в `build/linux/Taskfile.yml`:

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  SIGN_ROLE: "builder"  # origin, maint, archive, or builder
```

Сохраните пароль ключа:

```bash
wails3 signing
```

Подробности см. в разделе [Подписание приложений](/guides/build/signing/).

## Сборка для ARM

```bash
wails3 build GOOS=linux GOARCH=arm64
wails3 package GOOS=linux GOARCH=arm64
```

@note{type="note"}
При сборке для ARM64 на хостах x86_64 для кросс-компиляции CGO используется Docker.

@end

## Поддержка устаревшего GTK3

По умолчанию Wails v3 собирается на основе **GTK4 с WebKitGTK 6.0**. Для дистрибутивов, в которых WebKitGTK 6.0 пока отсутствует (Ubuntu 22.04 LTS, Debian 12, Fedora ≤ 39, RHEL 9.x), по-прежнему доступен устаревший вариант с GTK3 / WebKit2GTK 4.1. Этот вариант включается явно с помощью тега сборки; его планируется удалить в v3.1.

@note{type="caution" title="Устаревший вариант"}
Вариант с GTK3 / WebKit2GTK 4.1 поддерживается на протяжении ветки v3.0.x. Планируйте переход на GTK4 с учётом доступности GTK4 / WebKitGTK 6.0 в целевом дистрибутиве: `-tags gtk3` будет удалён в v3.1.

@end

### Зависимости

Установите библиотеки для разработки GTK3 и WebKit2GTK 4.1:

```bash
# Ubuntu/Debian
sudo apt install libgtk-3-dev libwebkit2gtk-4.1-dev

# Fedora
sudo dnf install gtk3-devel webkit2gtk4.1-devel

# Arch
sudo pacman -S gtk3 webkit2gtk-4.1
```

Требуются пакеты pkg-config `gtk+-3.0` и `webkit2gtk-4.1`.

### Сборка с GTK3

Используйте флаг `-tags gtk3`:

```bash
wails3 build -tags gtk3
```

Или выполните сборку непосредственно с помощью Go:

```bash
go build -tags gtk3 -o myapp .
```

### Известные отличия от GTK4

- **Диалоги выбора файлов**: по умолчанию GTK4 использует для диалогов выбора файлов `xdg-desktop-portal`, поэтому некоторые параметры диалогов, например начальный каталог и отображение пользовательских фильтров, работают иначе, чем в GTK3. Подробности см. в разделе [Справочник по диалогам — поведение диалогов в Linux](/reference/dialogs/#linux-dialog-behavior).
- **Стиль меню**: GTK4 поддерживает параметр `LinuxMenuStylePrimaryMenu`, который в соответствии с GNOME HIG отображает кнопку меню-гамбургера (☰) в панели заголовка. Этот параметр не влияет на сборки `-tags gtk3`. См. [API окна — Linux MenuStyle](/reference/window/#linux).
- **Масштабирование DPI**: для поддержки дробного масштабирования GTK4 использует `gdk_monitor_get_scale` (GTK 4.14+).

### Проверка сборки

Чтобы проверить конфигурацию, выполните `wails3 doctor`. Без флагов команда проверяет наличие GTK4 / WebKitGTK 6.0, используемых по умолчанию. Устаревшие пакеты GTK3 / WebKit2GTK 4.1 указаны как необязательные.

## Устранение неполадок

### AppImage не запускается

Сделайте файл исполняемым:

```bash
chmod +x MyApp-x86_64.AppImage
```

### Отсутствующие зависимости

Установите библиотеки, соответствующие варианту сборки исполняемого файла. Для стандартной сборки с GTK4/WebKitGTK 6.0:

```bash
# Debian/Ubuntu versions that provide WebKitGTK 6.0
sudo apt install libgtk-4-1 libwebkitgtk-6.0-4
# Fedora
sudo dnf install gtk4 webkitgtk6.0
# Arch
sudo pacman -S gtk4 webkitgtk-6.0
```

Для исполняемого файла, явно собранного с `-tags gtk3`, вместо них установите библиотеки устаревшего стека GTK3/WebKit2GTK 4.1:

```bash
# Debian/Ubuntu
sudo apt install libgtk-3-0 libwebkit2gtk-4.1-0
# Fedora
sudo dnf install gtk3 webkit2gtk4.1
# Arch
sudo pacman -S gtk3 webkit2gtk-4.1
```

Доступность и названия пакетов зависят от выпуска дистрибутива. Используйте поддерживаемый выпуск с необходимыми библиотеками либо документированный способ кросс-сборки или сборки в контейнере. Установка устаревших библиотек не удовлетворяет зависимостям стандартной сборки. Переносимость AppImage зависит от того, что фактически включено в пакет и какие библиотеки хоста ему по-прежнему нужны.

### Не найден компилятор C

Для работы CGO системе сборки требуется GCC или Clang:

```bash
# Debian/Ubuntu
sudo apt install build-essential

# Fedora
sudo dnf install gcc

# Arch
sudo pacman -S base-devel
```

Можно также выполнить `wails3 task setup:docker`, и система сборки автоматически воспользуется Docker.

### Пустое или белое окно на GPU NVIDIA

В Linux с проприетарными драйверами NVIDIA приложения Wails при запуске могут отображать пустое или белое окно. Причина — ошибка WebKitGTK: средство визуализации DMA-BUF завершается с ошибкой `gbm_bo_map()` при работе с проприетарным драйвером NVIDIA. Ошибка затрагивает X11 и Wayland, версии драйвера 377–580+, GPU серии 10 и более старые GT 710.

**Wails автоматически применяет `WEBKIT_DISABLE_DMABUF_RENDERER=1`** при обнаружении модуля ядра NVIDIA (`/sys/module/nvidia`), поэтому большинству пользователей ничего делать не потребуется.

Если окно по-прежнему остаётся пустым, например в контейнере, где путь к модулю недоступен, задайте переменную окружения вручную перед запуском приложения:

```bash
WEBKIT_DISABLE_DMABUF_RENDERER=1 ./myapp
```

Связанные ошибки в вышестоящем проекте: [WebKit #262607](https://bugs.webkit.org/show_bug.cgi?id=262607), [WebKit #180739](https://bugs.webkit.org/show_bug.cgi?id=180739).

### Совместимость удаления символов из AppImage

В современных дистрибутивах Linux (Arch Linux, Fedora 39+, Ubuntu 24.04+) системные библиотеки компилируются с разделами ELF `.relr.dyn` для более эффективного выполнения перемещений. Инструмент `linuxdeploy`, используемый для создания AppImage, включает более старый исполняемый файл `strip`, который не может обрабатывать эти современные разделы.

Перед сборкой AppImage Wails автоматически выявляет эту ситуацию, проверяя системные библиотеки GTK. Если она обнаружена, удаление символов отключается (`NO_STRIP=1`) для обеспечения совместимости.

**Что это означает:**

- В затронутых системах файлы AppImage будут немного больше (примерно на 20-40%)
- Функциональность приложения не изменяется
- Всё обрабатывается автоматически — никаких действий не требуется

Если в современных системах вам нужны файлы AppImage меньшего размера, установите более новый исполняемый файл `strip` и настройте `linuxdeploy` на его использование вместо встроенной версии.
