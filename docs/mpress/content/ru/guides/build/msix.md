---
title: "Упаковка в MSIX"
description: "Упаковка приложения Wails v3 в пакет MSIX"
slug: "guides/build/msix"
sourcePath: "guides/build/msix.md"
---

MSIX — это современный формат упаковки приложений Windows. Wails может создавать пакет MSIX в рамках сборки для Windows.

Инструкции по упаковке в MSIX приведены в руководстве [«Упаковка для Windows»](/guides/build/windows/#msix-package).

## Создание пакетов напрямую через CLI

Запускайте инструмент MSIX в Windows, в том числе в CI. По умолчанию используется `MakeAppx.exe`; для подписи также нужен `signtool.exe`. Оба входят в Windows SDK. Помощник установки открывает Microsoft Store и при необходимости страницу загрузки SDK; завершите установку перед созданием пакета.

Задайте идентификационные данные в `build/config.yml`, затем соберите исполняемый файл и создайте пакет:

```yaml
info:
  companyName: "Example Corp"
  productName: "MyApp"
  productIdentifier: "com.example.myapp"
  description: "MyApp"
  version: "1.0.0"
```

```powershell
wails3 tool msix-install-tools
wails3 build GOOS=windows
wails3 tool msix --executable bin/myapp.exe --name myapp.exe
```

Прямая команда записывает `MyApp.msix` в текущий каталог. [Задача упаковки для Windows](/guides/build/windows/#msix-package) задаёт собственный путь вывода.

## Параметры CLI

| Параметр | Назначение |
| --- | --- |
| `--config` | Файл конфигурации; по умолчанию `build/config.yml`. |
| `--executable`, `--name` | Существующий исполняемый файл и его имя внутри пакета; оба параметра обязательны. |
| `--out` | Выходной файл; по умолчанию `<ProductName>.msix`. |
| `--arch` | Архитектура пакета: `x64` (по умолчанию), `x86`, `arm`, `arm64`, `x86a64` или `neutral`. Допускаются псевдонимы Go `amd64` и `386`. Архитектура должна соответствовать исполняемому файлу. |
| `--publisher` | Идентификатор издателя; по умолчанию `CN=<companyName>`. |
| `--cert`, `--cert-password` | Путь к сертификату PFX и пароль для подписи. |
| `--use-makeappx` | Использовать стандартный упаковщик Windows SDK. |
| `--use-msix-tool` | Явно выбрать `MsixPackagingTool.exe`, который должен находиться в `PATH`. |

## Подпись и CI

Для распространения вне Store подпишите пакет сертификатом, которому доверяет целевая машина. Поле Subject сертификата должно точно совпадать с `--publisher`. При передаче `--cert` механизм MakeAppx вызывает SignTool с SHA256. См. [руководство Microsoft по подписи](https://learn.microsoft.com/en-us/windows/msix/package/sign-msix-package-guide).

Этот шаг рабочего процесса Windows предполагает, что Wails и SDK установлены, а предыдущий шаг безопасно разместил файл PFX по пути `CERT_PATH`. Секрет, содержащий только путь, сам по себе не загружает сертификат:

```yaml
- name: MSIX
  if: runner.os == 'Windows'
  shell: pwsh
  run: |
    wails3 build GOOS=windows
    wails3 tool msix --executable bin/myapp.exe --name myapp.exe --publisher "$env:MSIX_PUBLISHER" --cert "$env:CERT_PATH" --cert-password "$env:CERT_PASSWORD"
  env:
    MSIX_PUBLISHER: ${{ vars.MSIX_PUBLISHER }}
    CERT_PATH: ${{ secrets.WINDOWS_CERT_PATH }}
    CERT_PASSWORD: ${{ secrets.WINDOWS_CERT_PASSWORD }}
```

## Ассоциации файлов и ресурсы

Добавьте расширения без начальной точки в `build/config.yml`; генерируемый манифест добавит точку:

```yaml
fileAssociations:
  - ext: myext
    name: MyApp Document
    description: MyApp Document
    iconName: fileicon
```

Обрабатывайте открытие файлов во время выполнения, как описано в разделе [Ассоциации файлов](/guides/file-associations/). Механизм MakeAppx сейчас копирует только исполняемый файл и создаёт прозрачные изображения-заглушки. Он не импортирует файлы `Assets/` проекта и не преобразует значки `iconName`. Для фирменных изображений или дополнительных DLL используйте собственный процесс упаковки.

| Генерируемый ресурс | Размер (пиксели) |
| --- | --- |
| `Square150x150Logo.png` | 150×150 |
| `Square44x44Logo.png` | 44×44 |
| `Wide310x150Logo.png` | 310×150 |
| `StoreLogo.png` | 50×50 |
| `SplashScreen.png` | 620×300 |
| `FileIcon.png` | 44×44 |

`FileIcon.png` создаётся только при наличии настроенных ассоциаций файлов. Эти файлы находятся в каталоге `Assets/` пакета.

## Отправка в Store и устранение неполадок

Зарезервируйте приложение в [портале Partner Center](https://partner.microsoft.com/dashboard) и при подготовке отправки используйте указанные там идентификатор пакета и издателя. Store подписывает пакеты MSIX при отправке; для этого способа распространения покупать сертификат подписи не нужно. См. [требования Microsoft к пакетам](https://learn.microsoft.com/en-us/windows/apps/publish/publish-your-app/msix/app-package-requirements).

Если `MakeAppx.exe` или `signtool.exe` не найден, установите или восстановите Windows SDK. Wails ищет в `PATH` и стандартных каталогах SDK. При ошибках подписи проверьте поле Subject сертификата, срок его действия и доверие на целевой машине; см. [устранение неполадок MSIX](https://learn.microsoft.com/en-us/windows/msix/msix-troubleshooting-guide).
