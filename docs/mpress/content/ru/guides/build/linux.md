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

Чтобы собрать пакеты определённых форматов:

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
wails3 setup signing
```

Подробности см. в разделе [«Подписание приложений»](/guides/build/signing/).

## Сборка для ARM

```bash
wails3 build GOOS=linux GOARCH=arm64
wails3 package GOOS=linux GOARCH=arm64
```

@note{type="note"}
При сборке для ARM64 на компьютерах x86_64 для кросс-компиляции CGO используется Docker.

@end

## Поддержка устаревшей версии GTK3

По умолчанию Wails v3 использует для сборки **GTK4 с WebKitGTK 6.0**. Для дистрибутивов, в которых WebKitGTK 6.0 ещё не поставляется (Ubuntu 22.04 LTS, Debian 12, Fedora ≤ 39, RHEL 9.x), по-прежнему доступен устаревший вариант с GTK3 / WebKit2GTK 4.1. Этот вариант включается явно с помощью тега сборки, а его удаление запланировано на v3.1.

@note{type="caution" title="Устаревший вариант"}
Вариант с GTK3 / WebKit2GTK 4.1 поддерживается в линейке v3.0.x. Планируйте переход на GTK4 с учётом доступности GTK4 / WebKitGTK 6.0 в целевом дистрибутиве — `-tags gtk3` будет удалён в v3.1.

@end

### Зависимости

Установите библиотеки разработки GTK3 и WebKit2GTK 4.1:

```bash
# Ubuntu/Debian
sudo apt install libgtk-3-dev libwebkit2gtk-4.1-dev

# Fedora
sudo dnf install gtk3-devel webkit2gtk4.1-devel

# Arch
sudo pacman -S gtk3 webkit2gtk-4.1
```

Необходимые пакеты pkg-config: `gtk+-3.0` и `webkit2gtk-4.1`.

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

- **Диалоговые окна выбора файлов**: GTK4 по умолчанию использует для диалоговых окон выбора файлов `xdg-desktop-portal`, поэтому некоторые параметры этих окон (например, начальный каталог и отображение пользовательских фильтров) работают иначе, чем в GTK3. Подробности см. в разделе [«Справочник по диалоговым окнам — поведение диалоговых окон в Linux»](/reference/dialogs/#linux-dialog-behavior).
- **Стиль меню**: GTK4 поддерживает параметр `LinuxMenuStylePrimaryMenu`, который в соответствии с GNOME HIG отображает в панели заголовка кнопку меню-гамбургера (☰). Этот параметр не влияет на сборки `-tags gtk3`. См. [«API окон — Linux MenuStyle»](/reference/window/#linux).
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

Если приложение не запускается, проверьте наличие отсутствующих зависимостей WebKit:

```bash
# Debian/Ubuntu
sudo apt install libwebkit2gtk-4.1-0

# Fedora
sudo dnf install webkit2gtk4.1

# Arch
sudo pacman -S webkit2gtk-4.1
```

### Компилятор C не найден

Для CGO системе сборки необходим GCC или Clang:

```bash
# Debian/Ubuntu
sudo apt install build-essential

# Fedora
sudo dnf install gcc

# Arch
sudo pacman -S base-devel
```

Также можно выполнить `wails3 task setup:docker`, и система сборки автоматически воспользуется Docker.

### Пустое или белое окно при использовании графического процессора NVIDIA

В Linux при использовании проприетарных драйверов NVIDIA приложения Wails могут при запуске отображать пустое или белое окно. Причина — ошибка WebKitGTK, из-за которой средство визуализации DMA-BUF не работает с `gbm_bo_map()` и проприетарным драйвером NVIDIA (затрагиваются X11 и Wayland, версии драйвера 377–580+, графические процессоры серии 10 и более старые GT 710).

**При обнаружении модуля ядра NVIDIA (`/sys/module/nvidia`) Wails автоматически применяет `WEBKIT_DISABLE_DMABUF_RENDERER=1`**, поэтому большинству пользователей ничего делать не потребуется.

Если окно по-прежнему остаётся пустым (например, в контейнере, где путь к модулю недоступен), задайте переменную окружения вручную перед запуском приложения:

```bash
WEBKIT_DISABLE_DMABUF_RENDERER=1 ./myapp
```

Связанные ошибки в вышестоящем проекте: [WebKit #262607](https://bugs.webkit.org/show_bug.cgi?id=262607), [WebKit #180739](https://bugs.webkit.org/show_bug.cgi?id=180739).

### Совместимость AppImage с удалением отладочных символов

В современных дистрибутивах Linux (Arch Linux, Fedora 39+, Ubuntu 24.04+) системные библиотеки компилируются с разделами ELF `.relr.dyn` для более эффективного выполнения перемещений. Инструмент `linuxdeploy`, используемый для создания AppImage, включает более старый исполняемый файл `strip`, который не может обрабатывать эти современные разделы.

Wails автоматически обнаруживает эту ситуацию, проверяя системные библиотеки GTK перед сборкой AppImage. При обнаружении удаление отладочной информации отключается (`NO_STRIP=1`) для обеспечения совместимости.

**Что это означает:**

- В затронутых системах файлы AppImage будут немного больше (~20-40 %)
- Это не влияет на функциональность приложения
- Это обрабатывается автоматически — никаких действий не требуется

Если в современных системах вам нужны файлы AppImage меньшего размера, можно установить более новую версию исполняемого файла `strip` и настроить `linuxdeploy` на её использование вместо версии из комплекта поставки.
