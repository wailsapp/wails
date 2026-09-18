---
title: "Упаковка для Windows"
description: "Упакуйте приложение Wails для распространения в Windows"
slug: "guides/build/windows"
sourcePath: "guides/build/windows.md"
---

## Установщик NSIS

Формат упаковки по умолчанию создаёт установщик NSIS:

```bash
wails3 package GOOS=windows
```

При этом запускается `wails3 task windows:package`, который:

1. Собирает приложение
2. Создаёт загрузчик WebView2
3. Создаёт установщик NSIS

Результат: `build/windows/nsis/<AppName>-installer.exe`

### Пакет MSIX

Для распространения через Microsoft Store или современного развёртывания в Windows:

```bash
wails3 package GOOS=windows FORMAT=msix
```

Результат: `bin/<AppName>-<arch>.msix`

@note{type="note"}
Для MSIX требуется либо `makeappx.exe` (Windows SDK), либо отдельный набор инструментов MSIX. В Taskfile для Windows задача установки доступна как `wails3 task install:msix:tools`.

@end

## Настройка установщика

Конфигурация NSIS находится в `build/windows/nsis/project.nsi`. Отредактируйте этот файл, чтобы настроить:

- Интерфейс и фирменное оформление установщика
- Каталог установки
- Ярлыки в меню «Пуск» и на рабочем столе
- Ассоциации файлов
- Лицензионное соглашение

Метаданные приложения берутся из `build/windows/info.json`:

```json
{
  "fixed": {
    "file_version": "1.0.0"
  },
  "info": {
    "0000": {
      "ProductVersion": "1.0.0",
      "CompanyName": "My Company",
      "FileDescription": "My Application",
      "ProductName": "MyApp"
    }
  }
}
```

## Подписание кода

Подпишите исполняемый файл и установщик, чтобы избежать предупреждений SmartScreen:

```bash
# Using the wrapper (auto-detects platform)
wails3 sign GOOS=windows

# Or using tasks directly
wails3 task windows:sign
wails3 task windows:sign:installer
```

Настройте подписание в `build/windows/Taskfile.yml`:

```yaml
vars:
  SIGN_CERTIFICATE: "path/to/certificate.pfx"
  # Or use thumbprint for certificates in Windows store
  SIGN_THUMBPRINT: "certificate-thumbprint"
  TIMESTAMP_SERVER: "http://timestamp.digicert.com"
```

Храните пароль сертификата в безопасном месте:

```bash
wails3 setup signing
```

Подробности см. в разделе [«Подписание приложений»](/guides/build/signing/).

## Сборка для ARM

```bash
wails3 build GOOS=windows GOARCH=arm64
wails3 package GOOS=windows GOARCH=arm64
```

## Устранение неполадок

### Команда makensis не найдена

Установите NSIS:

```bash
# Windows
winget install NSIS.NSIS

# Or download from https://nsis.sourceforge.io/
```

### Предупреждение SmartScreen

Исполняемый файл не подписан. См. раздел [«Подписание кода»](#--1) выше.

### Отсутствует WebView2

Установщик содержит загрузчик WebView2, который при необходимости скачивает среду выполнения. Для автономной установки скачайте Evergreen Standalone Installer с сайта Microsoft.
