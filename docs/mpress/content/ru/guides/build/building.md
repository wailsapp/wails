---
title: "Сборка приложений"
description: "Сборка и упаковка приложения Wails"
slug: "guides/build/building"
sourcePath: "guides/build/building.md"
---

Wails v3 использует [Task](https://taskfile.dev) в качестве системы сборки. Команды `wails3 build` и `wails3 package` служат удобными оболочками для Task.

## Сборка

Соберите приложение для текущей платформы:

```bash
wails3 build
```

Соберите приложение для указанной платформы:

```bash
wails3 build GOOS=windows
wails3 build GOOS=darwin
wails3 build GOOS=linux

# With architecture
wails3 build GOOS=darwin GOARCH=arm64

# Environment variable style works too
GOOS=windows wails3 build
```

Результат сборки помещается в каталог `bin/`.

@note{type="tip"}
Для кросс-компиляции под macOS или Linux с другой платформы требуется Docker. Инструкции по настройке см. в разделе [«Кроссплатформенная сборка»](/guides/build/cross-platform/).

@end

## Разработка

Запустите приложение с автоматической перезагрузкой:

```bash
wails3 dev
```

При этом запускается средство отслеживания файлов, которое при изменениях повторно собирает и перезапускает приложение. По умолчанию сервер разработки фронтенда работает на порте 9245.

```bash
# Custom port
wails3 dev -port 3000

# Enable HTTPS
wails3 dev -s
```

## Упаковка

Упакуйте приложение для распространения:

```bash
wails3 package
wails3 package GOOS=windows
wails3 package GOOS=darwin
wails3 package GOOS=linux
```

Будут созданы пакеты для соответствующих платформ:

- **Windows**: установщик NSIS — см. раздел [«Упаковка для Windows»](/guides/build/windows/)
- **macOS**: пакет приложения (`.app`) — см. раздел [«Упаковка для macOS»](/guides/build/macos/)
- **Linux**: AppImage, deb и rpm — см. раздел [«Упаковка для Linux»](/guides/build/linux/)

## Пользовательские теги сборки

Передайте пользовательские теги сборки Go с помощью флага `-tags`:

```bash
# Build with legacy GTK3 + WebKit2GTK 4.1 on Linux (default is GTK4 + WebKitGTK 6.0)
wails3 build -tags gtk3

# Build in server mode (no GUI, CGO-free)
wails3 build -tags server

# Combine multiple tags
wails3 build -tags gtk3,customtag
```

Теги передаются в базовый Taskfile как `EXTRA_TAGS`. Подробности см. в разделах [«Сборка сервера»](/guides/server-build/) и [«Упаковка для Linux — поддержка устаревшей версии GTK3»](/guides/build/linux/#legacy-gtk3-support).

## Непосредственное использование Task

Для более точного управления используйте Task напрямую:

```bash
# List available tasks
wails3 task --list

# Verbose output
wails3 task build -v

# Dry run
wails3 task --dry

# Force rebuild
wails3 task build -f

# Pass variables
wails3 task darwin:build ARCH=amd64
```

Задачи для конкретных платформ, такие как `linux:create:deb` или `darwin:build:universal`, доступны только через Task.

## Создание ресурсов

Повторно создайте значки или обновите конфигурацию сборки:

```bash
wails3 generate icons -input build/appicon.png
wails3 update build-assets -name "MyApp" -config build/config.yml -dir build
```
