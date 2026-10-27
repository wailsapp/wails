---
title: "Подписание кода"
description: "Руководство по подписанию приложений Wails на всех платформах"
slug: "guides/build/signing"
sourcePath: "guides/build/signing.md"
---

## Подписание кода приложения

В этом руководстве описано, как подписывать приложения Wails для macOS, Windows и Linux. Wails v3 предоставляет встроенные инструменты командной строки для подписания кода, нотариального заверения и управления ключами PGP.

- **macOS** — подписывайте и нотариально заверяйте приложения для macOS
- **Windows** — подписывайте исполняемые файлы и пакеты для Windows
- **Linux** — подписывайте пакеты DEB и RPM ключами PGP

## Матрица кроссплатформенного подписания

Эта матрица показывает, какие форматы можно подписывать на каждой исходной платформе:

| Целевой формат | Из Windows | Из macOS | Из Linux |
| --- | :---: | :---: | :---: |
| Windows EXE/MSI | ✅ | ✅ | ✅ |
| Пакет приложения macOS .app | ❌ | ✅ | ❌ |
| Нотариальное заверение macOS | ❌ | ✅ | ❌ |
| Linux DEB | ✅ | ✅ | ✅ |
| Linux RPM | ✅ | ✅ | ✅ |

@note{type="tip"}
Пакеты Windows и Linux можно подписывать на **любой платформе**. Для подписания macOS требуется компьютер Mac из-за требований к инструментам Apple.

@end

### Механизмы подписания

Wails автоматически выбирает лучший из доступных механизмов подписания:

| Платформа | Нативный механизм | Кроссплатформенный механизм |
| --- | --- | --- |
| Windows | `signtool.exe` (Windows SDK) | Встроенный |
| macOS | `codesign` (Xcode) | Недоступен |
| Linux | Неприменимо | Встроенный |

При запуске на целевой платформе Wails использует нативные инструменты для максимальной совместимости. При кросс-компиляции используется встроенная поддержка подписания.

## Быстрый старт

Самый быстрый способ настроить подписывание — использовать мастер настройки:

```bash
wails3 setup
```

При этом общая конфигурация подписания записывается в `~/.config/wails/defaults.yaml` и **применяется при подписании на любой платформе** (см. раздел [«Приоритет конфигурации»](#--2)). На этапе настройки подписания мастер:

- Определяет, какие инструменты подписания установлены для выбранной целевой платформы — **на любом хосте** — и показывает команду установки для вашей ОС (например, `brew install gnupg`, `sudo apt install osslsigncode`, `winget install GnuPG.Gpg4win`).
- В macOS выводит список сертификатов Developer ID из связки ключей.
- Для Linux выводит список ключей GPG и позволяет **создать и экспортировать** новый ключ.
- Для Windows позволяет с помощью OpenSSL **создать самоподписанный сертификат** (для тестирования).
- **Безопасно хранит пароли в системном хранилище учётных данных** (а не в Taskfile).

@note{type="tip"}
Пароли хранятся в нативном хранилище учётных данных системы: «Связке ключей» macOS, «Диспетчере учётных данных» Windows или Linux Secret Service. Благодаря этому конфигурация подписания защищена и работает во всех ваших проектах Wails.

@end

## Приоритет конфигурации

При запуске задачи подписания значение каждого параметра определяется в следующем порядке (используется первое найденное):

1. Явный флаг, переданный `wails3 tool sign` (например, `--pgp-key`, `--certificate`, `--identity`).
2. Соответствующая **переменная Taskfile проекта** (`PGP_KEY`, `SIGN_CERTIFICATE`/`SIGN_THUMBPRINT`, `SIGN_IDENTITY`, …).
3. **Глобальная конфигурация** в `~/.config/wails/defaults.yaml` (записывается командой `wails3 setup`).

Это означает, что переменные Taskfile служат **необязательными переопределениями**: если переменная не задана, используется ключ, сертификат или идентификатор из глобальной конфигурации. Если значение отсутствует во всех трёх источниках, команда подписания выводит понятное сообщение об ошибке с указанием способа настройки.

## Конфигурация для отдельного проекта

Чтобы настроить подписание только для одного проекта, а не глобально, запустите мастер конфигурации проекта из каталога этого проекта — он запишет `vars` в файлы `build/<platform>/Taskfile.yml` проекта:

```bash
wails3 signing                                   # all detected platforms
wails3 signing --platform windows --platform linux
```

### Настройка вручную

Кроме того, можно вручную отредактировать Taskfile для каждой платформы. Измените раздел `vars` в начале каждого файла:

@tabs
[macOS]
Отредактируйте `build/darwin/Taskfile.yml`:

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
  # ENTITLEMENTS: "build/darwin/entitlements.plist"
```

Затем выполните:

```bash
wails3 task darwin:sign           # Sign only
wails3 task darwin:sign:notarize  # Sign and notarize
```

[Windows]
Отредактируйте `build/windows/Taskfile.yml`:

```yaml
vars:
  SIGN_CERTIFICATE: "path/to/certificate.pfx"
  # Or use thumbprint instead:
  # SIGN_THUMBPRINT: "certificate-thumbprint"
  # TIMESTAMP_SERVER: "http://timestamp.digicert.com"
```

Пароль извлекается из системного хранилища учётных данных (для настройки выполните `wails3 signing`).

Затем выполните:

```bash
wails3 task windows:sign           # Sign executable
wails3 task windows:sign:installer # Sign NSIS installer
```

[Linux]
Отредактируйте `build/linux/Taskfile.yml`:

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  # SIGN_ROLE: "builder"  # Options: origin, maint, archive, builder
```

Пароль извлекается из системного хранилища учётных данных (для настройки выполните `wails3 signing`).

Затем выполните:

```bash
wails3 task linux:sign:deb       # Sign DEB package
wails3 task linux:sign:rpm       # Sign RPM package
wails3 task linux:sign:packages  # Sign all packages
```

@end

Состояние подписания также можно проверить непосредственно средствами системы:

```bash
# List available macOS code-signing identities
security find-identity -v -p codesigning

# List PGP keys (Linux package signing)
gpg --list-keys
```

Для интерактивной настройки конфигурации подписания на всех платформах используйте мастер:

```bash
wails3 signing
```

## Подписание кода в macOS

### Предварительные требования

- Учётная запись Apple Developer ($99 в год)
- Сертификат Developer ID Application
- Установленные инструменты командной строки Xcode

### Идентификаторы подписания

Проверьте доступные идентификаторы подписания:

```bash
security find-identity -v -p codesigning
```

Вывод:

```
Found 2 signing identities:

  Developer ID Application: Your Company (ABCD1234) [valid]
    Hash: ABC123DEF456...

  Apple Development: your@email.com (XYZ789) [valid]
    Hash: DEF789ABC123...
```

@note{type="tip"}
Для распространения вне App Store необходим сертификат **Developer ID Application**.

@end

### Конфигурация

Отредактируйте `build/darwin/Taskfile.yml` и задайте переменные подписания:

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
  ENTITLEMENTS: "build/darwin/entitlements.plist"
```

| Переменная | Обязательно | Описание |
| --- | --- | --- |
| `SIGN_IDENTITY` | Да | Ваш Developer ID (например, "Developer ID Application: Your Company (TEAMID)") |
| `KEYCHAIN_PROFILE` | Для нотариального заверения | Имя профиля связки ключей с сохранёнными учётными данными |
| `ENTITLEMENTS` | Нет | Путь к файлу разрешений |

Затем выполните:

```bash
wails3 task darwin:sign           # Build, package, and sign
wails3 task darwin:sign:notarize  # Build, package, sign, and notarize
```

### Права доступа

Права доступа определяют, к каким возможностям имеет доступ ваше приложение. Обычно приложениям Wails требуются разные права доступа для разработки и промышленной эксплуатации:

- **Разработка**: требуются права на JIT-компиляцию, неподписанную память и отладку
- **Промышленная эксплуатация**: минимальные права (только доступ к сети)

Создайте оба файла с помощью интерактивного мастера настройки:

```bash
wails3 entitlements
```

Будут созданы:

- `build/darwin/entitlements.dev.plist` — для сборок, используемых при разработке
- `build/darwin/entitlements.plist` — для промышленных/подписанных сборок

**Доступные предустановки:**\

| Предустановка | Описание |
| --- | --- |
| Разработка | JIT-компиляция, неподписанная память, отладка, сеть |
| Промышленная эксплуатация | Только сеть (минимальный и наиболее безопасный вариант) |
| Оба варианта | Создаёт файлы и для разработки, и для промышленной эксплуатации (рекомендуется) |
| App Store | Песочница с доступом к сети и файлам |
| Пользовательская | Выбор отдельных прав доступа |

@note{type="note"}
Задача `run` в Taskfile для darwin автоматически использует `entitlements.dev.plist`. Задачи `sign` используют `entitlements.plist` для промышленных сборок.

@end

Затем задайте для `ENTITLEMENTS` в переменных Taskfile путь к соответствующему файлу.

### Нотаризация

Apple требует нотариально заверять все распространяемые приложения.

@steps
### **Сохраните учётные данные в связке ключей** (однократная настройка). Запустите `wails3 signing` (команда запросит эти значения и самостоятельно вызовет `notarytool`) либо вызовите `notarytool` напрямую:
```bash
xcrun notarytool store-credentials "my-notarize-profile" \
  --apple-id "your@email.com" \
  --team-id "ABCD1234" \
  --password "app-specific-password"
```

### **Задайте KEYCHAIN_PROFILE в Taskfile** в соответствии с указанным выше именем профиля.
### **Подпишите и нотариально заверьте приложение**:
```bash
wails3 task darwin:sign:notarize
```

### **Проверьте нотариальное заверение**:
```bash
spctl --assess --verbose=2 bin/MyApp.app
```

@end

@note{type="note"}
Нотаризация обычно занимает 1-2 минут. Квитанция автоматически прикрепляется к приложению.

@end

## Подписание кода Windows

### Предварительные требования

- Сертификат подписи кода (от DigiCert, Sectigo и т. п.)
- Для нативного подписания в Windows: установленный Windows SDK (для `signtool.exe`)
- Для кроссплатформенного подписания из macOS/Linux: [`osslsigncode`](https://github.com/mtrojnar/osslsigncode) (на этапе подписания `wails3 setup` отображается команда установки для конкретной ОС)

### Создание самоподписанного сертификата (для тестирования)

Если вам нужно лишь проверить процесс подписания, запустите `wails3 setup`, откройте вкладку **Windows** на этапе подписания и выберите **Создать самоподписанный сертификат**. С помощью OpenSSL будет создан предназначенный для подписи кода файл `.pfx`, а путь к нему будет сохранён в глобальной конфигурации.

@note{type="caution"}
Самоподписанные сертификаты предназначены **только для тестирования и внутреннего распространения** — у конечных пользователей они вызывают предупреждения SmartScreen. Для общедоступных выпусков требуется сертификат от доверенного центра сертификации.

@end

### Конфигурация

Откройте `build/windows/Taskfile.yml` и задайте переменные подписания:

```yaml
vars:
  SIGN_CERTIFICATE: "path/to/certificate.pfx"
  # Or use thumbprint instead:
  # SIGN_THUMBPRINT: "certificate-thumbprint"
  # TIMESTAMP_SERVER: "http://timestamp.digicert.com"
```

| Переменная | Обязательность | Описание |
| --- | --- | --- |
| `SIGN_CERTIFICATE` | Переопределение | Путь к файлу сертификата .pfx/.p12 (если не задан, используется глобальная конфигурация) |
| `SIGN_THUMBPRINT` | Переопределение | Отпечаток сертификата в хранилище сертификатов Windows (альтернатива `SIGN_CERTIFICATE`) |
| `TIMESTAMP_SERVER` | Нет | URL-адрес сервера меток времени (по умолчанию: http://timestamp.digicert.com) |

@note{type="note"}
Эти переменные являются необязательными **переопределениями**. Если они не заданы, используется сертификат, настроенный глобально с помощью `wails3 setup`. См. раздел [Приоритет конфигурации](#--2).

@end

@note{type="note"}
Пароль сертификата хранится в системной связке ключей, а не в Taskfile. Для его настройки запустите `wails3 signing` или задайте переменную среды `WAILS_WINDOWS_CERT_PASSWORD` в CI.

@end

Затем выполните:

```bash
wails3 task windows:sign           # Build and sign executable
wails3 task windows:sign:installer # Build and sign NSIS installer
```

### Кроссплатформенное подписание

Исполняемые файлы Windows можно подписывать на любой платформе. Одинаковые конфигурация Taskfile и команды работают в macOS и Linux.

### Поддерживаемые форматы Windows

| Формат | Расширение | Примечания |
| --- | --- | --- |
| Исполняемые файлы | .exe | Стандартное подписание PE |
| Установщики | .msi | Пакеты Windows Installer |
| Пакеты приложений | .msix, .appx | Современные приложения Windows |

## Подписание пакетов Linux

Пакеты Linux (DEB и RPM) подписываются ключами PGP/GPG. В отличие от подписания кода Windows и macOS, подписание пакетов Linux подтверждает, что пакет получен из доверенного источника, а не то, что операционная система доверяет коду.

### Предварительные требования

- Пара ключей PGP (можно создать с помощью Wails)

### Создание ключа PGP

Проще всего воспользоваться мастером настройки: запустите `wails3 setup`, откройте вкладку **Linux** на этапе подписания и выберите **Создать новый ключ GPG**. Мастер:

- создаёт в вашей связке ключей GPG ключ RSA 4096 (оставьте парольную фразу пустой, чтобы ключ подходил для автоматической работы и CI),
- **экспортирует его в `~/.wails/signing/<keyid>.asc`** (файл, которым подписывается сборка) и
- сохраняет как идентификатор ключа, так и путь к экспортированному файлу в `~/.config/wails/defaults.yaml`, поэтому при подписании ключ используется автоматически.

@note{type="note"}
Если ключ был настроен **только по идентификатору** (например, создан до того, как мастер начал автоматически экспортировать ключи), при следующем открытии вкладки Linux он будет экспортирован в файл, а путь к нему — заполнен автоматически; никаких действий вручную не требуется. Сборки подписываются *файлом* ключа, поэтому значение имеет именно путь.

@end

Это также можно сделать вручную с помощью `gpg`:

```bash
# Interactive — the wizard will prompt for name, email, key size and expiry.
gpg --full-generate-key

# Export the key pair to ASCII-armoured files for the Taskfile to consume.
gpg --armor --export-secret-keys "your@email.com" > signing-key.asc
gpg --armor --export "your@email.com" > signing-key.pub.asc
```

Рекомендуемые параметры: RSA 4096 бит, срок действия 1 год, защита надёжным паролем.

@note{type="caution"}
Храните закрытый ключ в безопасности! Храните его в зашифрованном виде и создайте надёжную резервную копию.

@end

### Настройка

Отредактируйте `build/linux/Taskfile.yml` и задайте переменные подписывания:

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  # SIGN_ROLE: "builder"  # Options: origin, maint, archive, builder
```

| Переменная | Обязательно | Описание |
| --- | --- | --- |
| `PGP_KEY` | Переопределение | Путь к экспортированному файлу закрытого ключа PGP (если не задан, используется ключ, настроенный глобально через `wails3 setup`) |
| `SIGN_ROLE` | Нет | Роль подписи DEB (по умолчанию: builder) |

@note{type="note"}
Пароль ключа PGP хранится в системном хранилище ключей, а не в Taskfile. Чтобы настроить его, выполните `wails3 signing`; в CI вместо этого можно задать переменную окружения `WAILS_PGP_PASSWORD`.

@end

Затем выполните:

```bash
wails3 task linux:sign:deb       # Build and sign DEB package
wails3 task linux:sign:rpm       # Build and sign RPM package
wails3 task linux:sign:packages  # Build and sign all packages
```

### Роли подписи DEB

Для пакетов DEB роль подписи можно указать через `SIGN_ROLE`:

- `origin`: подпись создателя пакета
- `maint`: подпись сопровождающего пакета
- `archive`: подпись сопровождающего архива
- `builder`: подпись сборщика пакета (по умолчанию)

### Кроссплатформенное подписывание

Пакеты Linux можно подписывать на любой платформе. Одна и та же конфигурация Taskfile и те же команды работают в Windows и macOS.

### Просмотр сведений о ключе

```bash
gpg --show-keys signing-key.asc
```

Вывод:

```
pub   rsa4096 2024-01-15 [SC] [expires: 2025-01-15]
      1234 5678 90AB CDEF 1234 5678 90AB CDEF 1234 5678
uid                      Your Name <your@email.com>
```

### Проверка пакетов Linux

```bash
# Verify DEB signature
dpkg-sig --verify myapp_1.0.0_amd64.deb

# Verify RPM signature
rpm --checksig myapp-1.0.0.x86_64.rpm
```

### Распространение открытого ключа

Для проверки пакетов пользователям потребуется ваш открытый ключ:

```bash
# Export public key for distribution
gpg --armor --export "your@email.com" > myapp-signing.pub.asc

# Users can import it:
# For DEB (apt):
sudo apt-key add myapp-signing.pub.asc
# Or for modern apt:
sudo cp myapp-signing.pub.asc /etc/apt/trusted.gpg.d/

# For RPM:
sudo rpm --import myapp-signing.pub.asc
```

## Интеграция с GitHub Actions

В средах CI пароли передаются через переменные окружения, а не через системное хранилище ключей:

| Переменная окружения | Описание |
| --- | --- |
| `WAILS_WINDOWS_CERT_PASSWORD` | Пароль сертификата Windows |
| `WAILS_PGP_PASSWORD` | Пароль ключа PGP для пакетов Linux |

Также можно передавать переменные Taskfile напрямую:

```bash
wails3 task darwin:sign SIGN_IDENTITY="$SIGN_IDENTITY" KEYCHAIN_PROFILE="$KEYCHAIN_PROFILE"
```

### Рабочий процесс для macOS

```yaml
name: Build and Sign macOS

on:
  push:
    tags: ['v*']

jobs:
  build:
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Install Wails
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

      - name: Import Certificate
        env:
          CERTIFICATE_BASE64: ${{ secrets.MACOS_CERTIFICATE }}
          CERTIFICATE_PASSWORD: ${{ secrets.MACOS_CERTIFICATE_PASSWORD }}
        run: |
          echo $CERTIFICATE_BASE64 | base64 --decode > certificate.p12
          security create-keychain -p "" build.keychain
          security default-keychain -s build.keychain
          security unlock-keychain -p "" build.keychain
          security import certificate.p12 -k build.keychain -P "$CERTIFICATE_PASSWORD" -T /usr/bin/codesign
          security set-key-partition-list -S apple-tool:,apple:,codesign: -s -k "" build.keychain

      - name: Store Notarization Credentials
        env:
          APPLE_ID: ${{ secrets.APPLE_ID }}
          APPLE_TEAM_ID: ${{ secrets.APPLE_TEAM_ID }}
          APPLE_APP_PASSWORD: ${{ secrets.APPLE_APP_PASSWORD }}
        run: |
          xcrun notarytool store-credentials "notarize-profile" \
            --apple-id "$APPLE_ID" \
            --team-id "$APPLE_TEAM_ID" \
            --password "$APPLE_APP_PASSWORD"

      - name: Build, Sign, and Notarize
        env:
          SIGN_IDENTITY: ${{ secrets.MACOS_SIGN_IDENTITY }}
        run: |
          wails3 task darwin:sign:notarize \
            SIGN_IDENTITY="$SIGN_IDENTITY" \
            KEYCHAIN_PROFILE="notarize-profile"

      - name: Upload Artifact
        uses: actions/upload-artifact@v4
        with:
          name: MyApp-macOS
          path: bin/*.app
```

### Рабочий процесс для Windows

```yaml
name: Build and Sign Windows

on:
  push:
    tags: ['v*']

jobs:
  build:
    runs-on: windows-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Install Wails
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

      - name: Import Certificate
        env:
          CERTIFICATE_BASE64: ${{ secrets.WINDOWS_CERTIFICATE }}
        run: |
          $certBytes = [Convert]::FromBase64String($env:CERTIFICATE_BASE64)
          [IO.File]::WriteAllBytes("certificate.pfx", $certBytes)

      - name: Build and Sign
        env:
          WAILS_WINDOWS_CERT_PASSWORD: ${{ secrets.WINDOWS_CERTIFICATE_PASSWORD }}
        run: |
          wails3 task windows:sign SIGN_CERTIFICATE=certificate.pfx

      - name: Upload Artifact
        uses: actions/upload-artifact@v4
        with:
          name: MyApp-Windows
          path: bin/*.exe
```

### Кроссплатформенный рабочий процесс (исполнитель Linux)

Подписывайте пакеты Windows и Linux с помощью одного исполнителя Linux:

```yaml
name: Build and Sign (Cross-Platform)

on:
  push:
    tags: ['v*']

jobs:
  build-and-sign:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Install Wails
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

      - name: Install Build Dependencies
        run: |
          sudo apt-get update
          sudo apt-get install -y nsis rpm

      # Import certificates
      - name: Import Certificates
        env:
          WINDOWS_CERT_BASE64: ${{ secrets.WINDOWS_CERTIFICATE }}
          PGP_KEY_BASE64: ${{ secrets.PGP_PRIVATE_KEY }}
        run: |
          echo "$WINDOWS_CERT_BASE64" | base64 -d > certificate.pfx
          echo "$PGP_KEY_BASE64" | base64 -d > signing-key.asc

      # Build and sign Windows
      - name: Build and Sign Windows
        env:
          WAILS_WINDOWS_CERT_PASSWORD: ${{ secrets.WINDOWS_CERTIFICATE_PASSWORD }}
        run: |
          wails3 task windows:sign SIGN_CERTIFICATE=certificate.pfx

      # Build and sign Linux packages
      - name: Build and Sign Linux Packages
        env:
          WAILS_PGP_PASSWORD: ${{ secrets.PGP_PASSWORD }}
        run: |
          wails3 task linux:sign:packages PGP_KEY=signing-key.asc

      # Cleanup secrets
      - name: Cleanup
        if: always()
        run: rm -f certificate.pfx signing-key.asc

      - name: Upload Artifacts
        uses: actions/upload-artifact@v4
        with:
          name: signed-binaries
          path: |
            bin/*.exe
            bin/*.deb
            bin/*.rpm
```

@note{type="note"}
Использование исполнителя Linux для кроссплатформенного подписывания упрощает CI/CD, устраняя необходимость в отдельных исполнителях Windows. Для подписывания в macOS по-прежнему требуется исполнитель macOS из-за требований к нативным инструментам Apple.

@end

## Справочник по CLI

### wails3 signing

Интерактивный мастер настройки подписывания для вашего проекта.

```bash
wails3 signing [flags]

Flags:
  --platform    Platform to configure (darwin, windows, linux). Repeatable.
                If omitted, auto-detects which platforms to configure from the build directory.
```

Мастер поможет выполнить следующие действия:

- **macOS**: выбрать сертификат Developer ID и настроить учётные данные для нотариального заверения (вызывает `xcrun notarytool store-credentials`).
- **Windows**: выбрать файл сертификата или отпечаток, задать пароль и сервер меток времени.
- **Linux**: использовать существующий ключ PGP или создать новый (вызывает `gpg`), а также настроить роль подписи.

### wails3 entitlements

Интерактивный мастер настройки разрешений macOS.

```bash
wails3 entitlements [flags]

Flags:
  --output    Output path for entitlements.plist (default: build/darwin/entitlements.plist)
```

**Предустановки:**

- **Разработка**: создаёт `entitlements.dev.plist` с разрешениями для JIT, отладки и доступа к сети
- **Продакшен**: создаёт `entitlements.plist` с минимальным набором разрешений
- **Оба варианта**: создаёт оба файла (рекомендуется)
- **App Store**: создаёт разрешения для песочницы Mac App Store
- **Пользовательская настройка**: позволяет выбрать отдельные разрешения и целевой файл

### wails3 sign

Подписывает двоичные файлы и пакеты для текущей или указанной платформы. Это команда-обёртка, которая вызывает соответствующую платформе задачу подписывания.

```bash
wails3 sign
wails3 sign GOOS=darwin
wails3 sign GOOS=windows
wails3 sign GOOS=linux
```

При этом запускается соответствующая задача `<platform>:sign`, которая использует конфигурацию подписывания из вашего Taskfile.

### wails3 tool sign

Низкоуровневая команда для непосредственного подписывания указанного файла. Используется внутри Taskfile.

```bash
wails3 tool sign [flags]
```

**Общие флаги:**\

| Флаг | Описание |
| --- | --- |
| `--input` | Путь к подписываемому файлу |
| `--output` | Путь для выходного файла (необязательно; по умолчанию исходный файл изменяется на месте) |
| `--verbose` | Включить подробный вывод |

**Флаги Windows/macOS:**\

| Флаг | Описание |
| --- | --- |
| `--certificate` | Путь к сертификату PKCS#12 (.pfx/.p12) |
| `--password` | Пароль сертификата |
| `--timestamp` | URL сервера меток времени |

**Флаги только для macOS:**\

| Флаг | Описание |
| --- | --- |
| `--identity` | Идентификатор подписи (для подписи ad hoc используйте «-») |
| `--entitlements` | Путь к plist-файлу разрешений |
| `--hardened-runtime` | Включить усиленную среду выполнения (по умолчанию: true) |
| `--notarize` | Отправить на нотариальное заверение |
| `--keychain-profile` | Профиль хранилища ключей для нотариального заверения |

**Флаги только для Windows:**\

| Флаг | Описание |
| --- | --- |
| `--thumbprint` | Отпечаток сертификата в хранилище Windows |

**Флаги для Linux:**\

| Флаг | Описание |
| --- | --- |
| `--pgp-key` | Путь к закрытому ключу PGP |
| `--pgp-password` | Пароль ключа PGP |
| `--role` | Роль при подписании DEB (origin/maint/archive/builder) |

### Проверка состояния подписи (штатными инструментами)

`wails3 signing` настраивает учётные данные для подписывания. Чтобы проверить состояние подписи, используйте непосредственно нативные инструменты:

| Задача | Команда |
| --- | --- |
| Вывести список сертификатов для подписания кода macOS | `security find-identity -v -p codesigning` |
| Сохранить учётные данные для нотариального заверения | `xcrun notarytool store-credentials "<profile>" --apple-id … --team-id … --password …` |
| Просмотреть файл ключа PGP | `gpg --show-keys <key.asc>` |
| Создать пару ключей PGP | `gpg --full-generate-key` |
| Экспортировать открытый ключ | `gpg --armor --export <email>` |

## Устранение неполадок

### Проблемы в macOS

**«Сертификат Developer ID не найден»**

- Убедитесь, что сертификат установлен в Связке ключей
- С помощью `security find-identity -v -p codesigning` проверьте, не истёк ли срок его действия
- Убедитесь, что у вас есть сертификат «Developer ID Application», а не только «Apple Development»

**«Не удалось выполнить нотариальное заверение»**

- Проверьте журнал нотариального заверения: `xcrun notarytool log <submission-id> --keychain-profile <profile>`
- Убедитесь, что усиленная среда выполнения включена
- Убедитесь, что приложение не содержит неподписанных двоичных файлов

**«Не удалось выполнить codesign»**

- Убедитесь, что связка ключей разблокирована: `security unlock-keychain`
- Проверьте права доступа к файлам пакета приложения

### Проблемы в Windows

**«Сертификат не найден»**

- Убедитесь, что путь к сертификату указан правильно
- Проверьте пароль сертификата
- Убедитесь, что сертификат действителен: срок его действия не истёк и он не был отозван

**«Ошибка сервера меток времени»**

- Попробуйте другой сервер меток времени:
  - `http://timestamp.digicert.com`
  - `http://timestamp.sectigo.com`
  - `http://timestamp.comodoca.com`


### Проблемы в Linux

**«Недопустимый ключ PGP»**

- Убедитесь, что файл ключа имеет формат ASCII Armor
- С помощью `gpg --show-keys <key.asc>` проверьте, не истёк ли срок действия ключа
- Убедитесь, что пароль указан правильно

**«Не удалось проверить подпись»**

- Убедитесь, что открытый ключ импортирован правильно
- Убедитесь, что пакет не был изменён после подписания

## Дополнительные ресурсы

### Официальная документация

- [Руководство Apple по подписанию кода](https://developer.apple.com/support/code-signing/)
- [Документация Apple по нотариальному заверению](https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution)
- [Подписание кода Microsoft](https://docs.microsoft.com/en-us/windows-hardware/drivers/dashboard/get-a-code-signing-certificate)
- [Подписание пакетов Debian](https://wiki.debian.org/SecureApt)
- [Подписание пакетов RPM](https://rpm-software-management.github.io/rpm/manual/signatures.html)
