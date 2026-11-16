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
| Windows EXE/DLL/MSIX/APPX, экспортированный PFX | ✅ | ✅ | ✅ |
| Windows MSI или удостоверение из хранилища сертификатов либо от поставщика | ✅ | ❌ | ❌ |
| macOS .app / Mach-O / DMG, экспортированное удостоверение | ✅ | ✅ | ✅ |
| Нотариальное заверение macOS, учётные данные командного API | ✅ | ✅ | ✅ |
| Удостоверение из связки ключей macOS или профиль нотариального заверения | ❌ | ✅ | ❌ |
| Linux DEB | ✅ | ✅ | ✅ |
| Linux RPM | ✅ | ✅ | ✅ |

@note{type="tip"}
Экспортированные удостоверения Apple и учётные данные командного API позволяют выполнять переносимое подписание и нотариальное заверение. Для удостоверений из связки ключей macOS и профилей нотариального заверения требуется Mac. Ограничения для отдельных форматов описаны ниже.

@end

### Механизмы подписания

Wails автоматически выбирает лучший доступный механизм подписания:

| Платформа | Нативный механизм | Кроссплатформенный механизм |
| --- | --- | --- |
| Windows | `signtool.exe` для удостоверений из хранилища или от поставщика, MSI либо устаревших PFX с RC2 | Встроенный — для экспортированных PFX с RSA и форматов PE/MSIX |
| macOS | `codesign` для удостоверений из связки ключей | Встроенный — для экспортированных удостоверений PKCS#12 или специального подписания (ad hoc) |
| Linux | Не применимо | Встроенный |

Тип учётных данных определяет механизм подписания Apple на любом хосте. Для экспортированных удостоверений используется Go, а для удостоверений из связки ключей macOS — `codesign`. Получение, проверка и прикрепление билетов выполняются средствами Go.

## Быстрый старт

Быстрее всего настроить подписание с помощью мастера настройки:

```bash
wails3 setup
```

Он записывает общую конфигурацию подписания в `~/.config/wails/defaults.yaml`, которая **учитывается при подписании на всех платформах** (см. раздел [«Приоритет конфигурации»](#--2)). На этапе настройки подписания мастер:

- Показывает инструменты настройки для создания или экспорта сертификатов и ключей, например OpenSSL и GnuPG. Для экспортированных удостоверений Apple и Windows используется встроенное средство подписания; для удостоверений из хранилища сертификатов Windows или от поставщика, а также для пакетов MSI требуется Windows SignTool.
- В macOS выводит список сертификатов Developer ID из связки ключей.
- В Linux выводит список ключей GPG и позволяет **создать и экспортировать** новый ключ.
- В Windows позволяет через OpenSSL **создать самоподписанный сертификат** (для тестирования).
- **Безопасно сохраняет пароли в системном хранилище учётных данных** (а не в Taskfile-файлах).

@note{type="tip"}
Пароли сохраняются в нативном хранилище учётных данных системы: «Связке ключей» macOS, «Диспетчере учётных данных» Windows или службе секретов Linux. Благодаря этому конфигурация подписания защищена и работает во всех проектах Wails.

@end

## Приоритет конфигурации

При запуске задачи подписания значение каждого параметра определяется в следующем порядке (используется первое найденное):

1. Явно заданный флаг, переданный `wails3 tool sign` (например, `--pgp-key`, `--certificate` или `--identity`).
2. Соответствующая **переменная Taskfile проекта** (`PGP_KEY`, `SIGN_CERTIFICATE`/`SIGN_THUMBPRINT`, `SIGN_IDENTITY` и т. д.).
3. **Глобальная конфигурация** в `~/.config/wails/defaults.yaml` (записывается командой `wails3 setup`).

Таким образом, переменные Taskfile служат **необязательными переопределениями**: если переменная не задана, используется ключ, сертификат или удостоверение из глобальной конфигурации. Если значение отсутствует во всех трёх источниках, команда подписания выдаёт понятную ошибку с указанием способа настройки.

## Конфигурация для отдельного проекта

Чтобы настроить подписание только для одного проекта, а не глобально, запустите мастер конфигурации проекта из каталога этого проекта — он запишет `vars` в файлы `build/<platform>/Taskfile.yml` проекта:

```bash
wails3 signing                                   # all detected platforms
wails3 signing --platform windows --platform linux
```

### Настройка вручную

Также можно вручную изменить платформенные Taskfile-файлы. Отредактируйте раздел `vars` в начале каждого файла:

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

Для интерактивной настройки конфигурации подписания всех платформ используйте мастер:

```bash
wails3 signing
```

## Подписание кода в macOS

### Предварительные требования

- Учётная запись Apple Developer ($99 в год)
- Сертификат Developer ID Application
- Для удостоверения из связки ключей macOS: установленные инструменты командной строки Xcode
- Для переносимого подписания: экспортированное удостоверение RSA PKCS#12; инструменты подписания Apple не требуются

### Переносимое подписание с помощью Wake

Wake считывает `wails.yaml` напрямую и не выполняет Taskfile-файлы. Настройте экспортированное удостоверение Developer ID и ключ командного API App Store Connect:

```yaml
darwin:
  signing:
    certificate: signing/developer-id.p12
    credential: APPLE_CERT_PASSWORD
    timestamp_server: https://timestamp.digicert.com
  notarization:
    api_key: signing/AuthKey.p8
    key_id: ABCDEFGHIJ
    issuer_id: 00000000-0000-0000-0000-000000000000
```

`credential` задаёт имя переменной окружения, содержащей пароль сертификата. Храните оба файла закрытых ключей вне системы контроля версий. Включите подписание и нотариальное заверение в профиле выпуска, описанном в [руководстве по сборке с помощью YAML](/experimental/yaml-builds/).

Переносимое средство подписи поддерживает Mach-O-файлы для 64-разрядных процессоров Intel и Apple Silicon, универсальные двоичные файлы, стандартные пакеты приложений с вложенным кодом, фреймворки и образы UDIF DMG размером менее 4 ГиБ. Неподдерживаемые структуры отклоняются до публикации результата. Для подписи ad hoc используется `identity: "-"`, и такую подпись нельзя нотариально заверить. Удостоверения PKCS#12 могут использовать шифрование AES/PBKDF2 или устаревшее шифрование 3DES; экспортированные удостоверения с RC2 необходимо повторно экспортировать с поддерживаемым шифром.

Для переносимых нотариально заверенных подписей требуется явно заданная служба временных меток SHA-2 RFC 3161. Устаревшие подписи временных меток SHA-1 отклоняются. Отправка через Team API поддерживает одиночные загрузки размером до 5 ГиБ; индивидуальные ключи API и составные загрузки не поддерживаются. Проверка и прикрепление билета используют наблюдаемый протокол службы Apple с закреплённым сертификатом Apple Root G3, без `xcrun stapler`. Решение о доверии при распространении и об отзыве по-прежнему принимает Gatekeeper.

В следующих разделах описан существующий рабочий процесс с использованием связки ключей macOS, применяемый в устаревших проектах и конфигурациях Wake, в которых выбрано удостоверение из связки ключей.

### Удостоверения для подписи

Проверьте доступные удостоверения для подписи:

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
Для распространения вне App Store требуется сертификат **Developer ID Application**.

@end

### Конфигурация

Отредактируйте `build/darwin/Taskfile.yml` и задайте переменные подписи:

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
  ENTITLEMENTS: "build/darwin/entitlements.plist"
```

| Переменная | Обязательность | Описание |
| --- | --- | --- |
| `SIGN_IDENTITY` | Да | Ваш Developer ID (например, «Developer ID Application: Your Company (TEAMID)») |
| `KEYCHAIN_PROFILE` | Для нотариального заверения | Имя профиля связки ключей с сохранёнными учётными данными |
| `ENTITLEMENTS` | Нет | Путь к файлу разрешений |

Затем выполните:

```bash
wails3 task darwin:sign           # Build, package, and sign
wails3 task darwin:sign:notarize  # Build, package, sign, and notarize
```

### Разрешения

Разрешения определяют, к каким возможностям имеет доступ приложение. Для приложений Wails обычно требуются разные разрешения при разработке и в рабочей сборке:

- **Разработка**: требуются разрешения на JIT, неподписанную память и отладку
- **Рабочая сборка**: минимальные разрешения (только доступ к сети)

Создайте оба файла с помощью интерактивного мастера настройки:

```bash
wails3 entitlements
```

Будут созданы:

- `build/darwin/entitlements.dev.plist` — для сборок разработки
- `build/darwin/entitlements.plist` — для рабочих/подписанных сборок

**Доступные предустановки:**\

| Предустановка | Описание |
| --- | --- |
| Разработка | JIT, неподписанная память, отладка, сеть |
| Рабочая сборка | Только сеть (минимальный и наиболее безопасный набор) |
| Оба варианта | Создаёт файлы как для разработки, так и для рабочей сборки (рекомендуется) |
| App Store | Включена песочница с доступом к сети и файлам |
| Пользовательская | Выбор отдельных разрешений |

@note{type="note"}
Задача `run` в Taskfile для darwin автоматически использует `entitlements.dev.plist`. Задачи `sign` используют `entitlements.plist` для рабочих сборок.

@end

Затем задайте для `ENTITLEMENTS` в переменных Taskfile путь к соответствующему файлу.

### Нотариальное заверение

Apple требует нотариального заверения всех распространяемых приложений.

@steps
### **Сохраните учётные данные в связке ключей** (однократная настройка). Выполните `wails3 signing` (команда запросит эти значения и внутри вызовет `notarytool`) либо вызовите `notarytool` напрямую:
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
Нотариальное заверение обычно занимает 1-2 минут. Билет автоматически прикрепляется к приложению.

@end

## Подписание кода в Windows

### Предварительные требования

- Сертификат подписи кода (от DigiCert, Sectigo и т. п.)
- Для переносимой подписи: экспортированное удостоверение RSA PFX, зашифрованное с помощью AES/PBKDF2 или устаревшего алгоритма 3DES; внешние средства подписи не требуются
- Для удостоверений из хранилища сертификатов или от поставщика, пакетов MSI либо устаревших PFX с RC2: Windows со средствами подписи из Windows SDK в `PATH`

### Создание самоподписанного сертификата (для тестирования)

Если требуется только проверить конвейер подписи, выполните `wails3 setup`, откройте вкладку **Windows** на шаге подписи и выберите **Создать самоподписанный сертификат**. При этом OpenSSL создаст предназначенный для подписи кода файл `.pfx`, а путь к нему будет записан в глобальную конфигурацию.

@note{type="caution"}
Самоподписанные сертификаты предназначены **только для тестирования и внутреннего распространения** — у конечных пользователей они вызывают предупреждения SmartScreen. Для публичных выпусков требуется сертификат от доверенного центра сертификации.

@end

### Конфигурация

Отредактируйте `build/windows/Taskfile.yml` и задайте переменные подписи:

```yaml
vars:
  SIGN_CERTIFICATE: "path/to/certificate.pfx"
  # Or use thumbprint instead:
  # SIGN_THUMBPRINT: "certificate-thumbprint"
  # TIMESTAMP_SERVER: "http://timestamp.digicert.com"
```

| Переменная | Обязательность | Описание |
| --- | --- | --- |
| `SIGN_CERTIFICATE` | Переопределение | Путь к файлу сертификата .pfx/.p12 (если не задан, используется значение из глобальной конфигурации) |
| `SIGN_THUMBPRINT` | Переопределение | Отпечаток сертификата в хранилище сертификатов Windows (альтернатива `SIGN_CERTIFICATE`) |
| `TIMESTAMP_SERVER` | Нет | URL-адрес сервера временных меток (по умолчанию: http://timestamp.digicert.com) |

@note{type="note"}
Эти переменные являются необязательными **переопределениями**. Если они не заданы, используется сертификат, глобально настроенный с помощью `wails3 setup`. См. раздел [Приоритет конфигурации](#--2).

@end

@note{type="note"}
Пароль сертификата хранится в системной связке ключей, а не в Taskfile. Чтобы настроить его, выполните `wails3 signing` либо задайте переменную среды `WAILS_WINDOWS_CERT_PASSWORD` в CI.

@end

Затем выполните:

```bash
wails3 task windows:sign           # Build and sign executable
wails3 task windows:sign:installer # Build and sign NSIS installer
```

### Кроссплатформенная подпись

Экспортированные удостоверения RSA в формате PFX используются для подписания файлов EXE/DLL, исполняемых файлов установщиков NSIS и отдельных пакетов MSIX/APPX средствами Go в Windows, macOS и Linux. Существующие подписи можно заменить. При подписании выходной файл сначала помещается во временную область и публикуется только после успешной проверки подписи и метки времени; в случае сбоя исходный файл сохраняется. Используйте `--output`, чтобы записать подписанный файл отдельно.

Wake считывает эти параметры непосредственно из `wails.yaml`:

```yaml
windows:
  signing:
    certificate: signing/windows.pfx
    credential: WINDOWS_SIGNING_PASSWORD
    timestamp_server: http://timestamp.digicert.com
```

`credential` задаёт имя переменной окружения с паролем, но не сам пароль. Включите `sign: true` для цели выпуска; см. руководство по сборке [YAML](/experimental/yaml-builds/). Wake не запускает Taskfile. В устаревших проектах можно продолжать использовать приведённую выше конфигурацию Taskfile.

Переносимый профиль использует RSA/SHA-256 и проверяет метки времени RFC 3161, включая доверенную цепочку сертификатов подписанта метки времени. По умолчанию используется служба DigiCert. Значение MSIX Publisher должно совпадать с субъектом сертификата подписи; наборы пакетов и архивы размером не менее 4 ГиБ не поддерживаются. Размер файлов PFX ограничен 16 МиБ. Другие алгоритмы сертификатов и шифрование RC2 отклоняются. Для переносимого подписания повторно экспортируйте PFX с RC2, используя AES или 3DES; в Windows для прошедшего проверку подлинности PFX с RC2 вместо этого используется SignTool. Неверный пароль или повреждённый PFX никогда не приводит к такому переключению.

Для ключей из хранилища сертификатов, аппаратных токенов и ключей, предоставляемых провайдерами, сохраняется нативный рабочий процесс Windows SignTool с использованием `thumbprint`. Для подписания MSI также требуется SignTool в Windows. Настройка экспортированного сертификата не делает эти рабочие процессы переносимыми.

### Поддерживаемые форматы Windows

| Формат | Расширение | Примечания |
| --- | --- | --- |
| Исполняемые файлы и библиотеки | .exe, .dll | Переносимое подписание PE, включая установщики NSIS |
| Установщики | .msi | Требуется Windows SignTool |
| Пакеты приложений | .msix, .appx | Переносимое подписание отдельных пакетов; наборы не поддерживаются |

## Подписание пакетов Linux

Пакеты Linux (DEB и RPM) подписываются ключами PGP/GPG. В отличие от подписания кода в Windows и macOS, подписание пакетов Linux подтверждает, что пакет получен из доверенного источника, а не то, что операционная система считает код доверенным.

Для экспортированного файла закрытого ключа в текстовом формате ASCII Armor или двоичном формате используется нативный подписант Go в Windows, macOS и Linux. При подписании существующего пакета DEB или RPM его содержимое сохраняется без распаковки и повторного сжатия; устанавливать GnuPG, `dpkg-sig` или `rpmsign` не требуется. Для пакетов Arch тот же загрузчик ключей создаёт отсоединённую подпись. GnuPG по-прежнему полезен для создания и экспорта ключей, а также независимой проверки подписей. В устаревших конфигурациях DEB/RPM, где вместо файла ключа указано удостоверение из связки ключей, по-прежнему используются внешние средства подписания.

### Предварительные требования

- Пара ключей PGP (можно создать с помощью Wails)

### Создание ключа PGP

Проще всего воспользоваться мастером настройки: запустите `wails3 setup`, на этапе настройки подписи откройте вкладку **Linux** и выберите **Создать новый ключ GPG**. Мастер:

- создаёт ключ RSA 4096 в вашей связке ключей GPG (оставьте парольную фразу пустой, чтобы ключ можно было использовать без участия пользователя и в CI),
- **экспортирует его в `~/.wails/signing/<keyid>.asc`** (этот файл используется для подписи при сборке) и
- сохраняет идентификатор ключа и путь к экспортированному файлу в `~/.config/wails/defaults.yaml`, чтобы ключ автоматически использовался при подписании.

@note{type="note"}
Если ключ был настроен **только по идентификатору** (например, создан до того, как мастер начал автоматически экспортировать ключи), при следующем открытии вкладки Linux ключ будет экспортирован в файл, а путь — заполнен автоматически; никаких действий вручную не требуется. При сборке для подписания используется *файл* ключа, поэтому важен именно путь.

@end

Это также можно сделать вручную с помощью `gpg`:

```bash
# Interactive — the wizard will prompt for name, email, key size and expiry.
gpg --full-generate-key

# Export the key pair to ASCII-armoured files for the Taskfile to consume.
gpg --armor --export-secret-keys "your@email.com" > signing-key.asc
gpg --armor --export "your@email.com" > signing-key.pub.asc
```

Рекомендуемые параметры: RSA, 4096 бит, срок действия 1 год, защита надёжным паролем.

@note{type="caution"}
Обеспечьте безопасность закрытого ключа! Храните его в зашифрованном виде и создайте надёжную резервную копию.

@end

### Конфигурация

Измените `build/linux/Taskfile.yml` и задайте переменные подписания:

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  # SIGN_ROLE: "builder"  # Options: origin, maint, archive, builder
```

| Переменная | Обязательна | Описание |
| --- | --- | --- |
| `PGP_KEY` | Переопределение | Путь к экспортированному файлу закрытого ключа PGP (если не задан, используется ключ, глобально настроенный через `wails3 setup`) |
| `SIGN_ROLE` | Нет | Роль при подписании DEB (по умолчанию: builder) |

@note{type="note"}
Пароль ключа PGP хранится в системной связке ключей, а не в Taskfile. Чтобы настроить его, запустите `wails3 signing`; в CI можно вместо этого задать переменную окружения `WAILS_PGP_PASSWORD`.

@end

Затем выполните:

```bash
wails3 task linux:sign:deb       # Build and sign DEB package
wails3 task linux:sign:rpm       # Build and sign RPM package
wails3 task linux:sign:packages  # Build and sign all packages
```

### Роли при подписании DEB

Для пакетов DEB роль подписанта можно указать через `SIGN_ROLE`:

- `origin`: подпись источника пакета
- `maint`: подпись сопровождающего пакета
- `archive`: подпись сопровождающего архива
- `builder`: подпись сборщика пакета (по умолчанию)

### Кроссплатформенное подписание

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

В средах CI пароли передаются через переменные окружения, а не через системную связку ключей:

| Переменная окружения | Описание |
| --- | --- |
| `WAILS_WINDOWS_CERT_PASSWORD` | Пароль сертификата Windows |
| `WAILS_PGP_PASSWORD` | Пароль ключа PGP для пакетов Linux |

Переменные Taskfile также можно передавать напрямую:

```bash
wails3 task darwin:sign SIGN_IDENTITY="$SIGN_IDENTITY" KEYCHAIN_PROFILE="$KEYCHAIN_PROFILE"
```

### Рабочий процесс macOS

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

### Рабочий процесс Windows

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
Использование исполнителя Linux для кроссплатформенного подписания упрощает CI/CD, устраняя потребность в отдельных исполнителях Windows. Для удостоверений в связке ключей Apple и профилей нотариального заверения требуется исполнитель macOS. Экспортированные удостоверения и учётные данные API команды поддерживают описанный выше переносимый рабочий процесс Wake.

@end

## Справочник CLI

### wails3 signing

Интерактивный мастер настройки подписания для проекта.

```bash
wails3 signing [flags]

Flags:
  --platform    Platform to configure (darwin, windows, linux). Repeatable.
                If omitted, auto-detects which platforms to configure from the build directory.
```

Мастер поможет выполнить следующие действия:

- **macOS**: выбор сертификата Developer ID и настройка учётных данных для нотариального заверения (вызывает `xcrun notarytool store-credentials`).
- **Windows**: выбор между файлом сертификата и отпечатком, задание пароля и сервера меток времени.
- **Linux**: использование существующего ключа PGP или создание нового (вызывает `gpg`), настройка роли для подписания.

### wails3 entitlements

Интерактивный мастер настройки разрешений macOS.

```bash
wails3 entitlements [flags]

Flags:
  --output    Output path for entitlements.plist (default: build/darwin/entitlements.plist)
```

**Предустановки:**

- **Разработка**: создаёт `entitlements.dev.plist` с разрешениями для JIT, отладки и сети
- **Рабочая среда**: создаёт `entitlements.plist` с минимальным набором разрешений
- **Оба варианта**: создаёт оба файла (рекомендуется)
- **App Store**: создаёт разрешения с песочницей для Mac App Store
- **Пользовательская настройка**: позволяет выбрать отдельные разрешения и целевой файл

### wails3 sign

Подписывает исполняемые файлы и пакеты для текущей или указанной платформы. Это команда-обёртка, которая вызывает соответствующую задачу подписания для конкретной платформы.

```bash
wails3 sign
wails3 sign GOOS=darwin
wails3 sign GOOS=windows
wails3 sign GOOS=linux
```

При этом запускается соответствующая задача `<platform>:sign`, использующая конфигурацию подписания из вашего Taskfile.

### wails3 tool sign

Низкоуровневая команда для непосредственного подписания указанного файла. Используется внутри Taskfile.

```bash
wails3 tool sign [flags]
```

**Общие флаги:**\

| Флаг | Описание |
| --- | --- |
| `--input` | Путь к подписываемому файлу |
| `--output` | Выходной путь (необязательно; по умолчанию файл изменяется на месте) |
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
| `--identity` | Идентификатор подписанта (используйте «-» для произвольной подписи) |
| `--entitlements` | Путь к plist-файлу разрешений |
| `--hardened-runtime` | Включить усиленную среду выполнения (по умолчанию: true) |
| `--notarize` | Отправить на нотариальное заверение |
| `--keychain-profile` | Профиль Связки ключей для нотариального заверения |
| `--notary-key` | Файл закрытого `.p8` ключа API команды App Store Connect; используется клиент Go Notary API |
| `--notary-key-id` | Идентификатор ключа API команды; обязателен при использовании `--notary-key` |
| `--notary-issuer-id` | Идентификатор издателя команды; обязателен при использовании `--notary-key` |

**Флаги только для Windows:**\

| Флаг | Описание |
| --- | --- |
| `--thumbprint` | Отпечаток сертификата в хранилище Windows |

**Флаги только для Linux:**\

| Флаг | Описание |
| --- | --- |
| `--pgp-key` | Путь к закрытому ключу PGP |
| `--pgp-password` | Пароль ключа PGP |
| `--role` | Роль подписания DEB (origin/maint/archive/builder) |

### Проверка состояния подписи (системными инструментами)

`wails3 signing` настраивает учётные данные для подписания. Чтобы проверить состояние подписи, используйте непосредственно системные инструменты:

| Задача | Команда |
| --- | --- |
| Вывести список идентификаторов для подписания кода macOS | `security find-identity -v -p codesigning` |
| Сохранить учётные данные для нотариального заверения | `xcrun notarytool store-credentials "<profile>" --apple-id … --team-id … --password …` |
| Проверить файл ключа PGP | `gpg --show-keys <key.asc>` |
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
- Убедитесь, что приложение не содержит неподписанных исполняемых файлов

**«Не удалось выполнить codesign»**

- Убедитесь, что Связка ключей разблокирована: `security unlock-keychain`
- Проверьте права доступа к файлам в пакете приложения

### Проблемы в Windows

**«Сертификат не найден»**

- Убедитесь, что путь к сертификату указан правильно
- Проверьте пароль сертификата
- Убедитесь, что сертификат действителен (срок его действия не истёк и он не был отозван)

**«Ошибка сервера меток времени»**

- Попробуйте другой сервер меток времени:
  - `http://timestamp.digicert.com`
  - `http://timestamp.sectigo.com`
  - `http://timestamp.comodoca.com`


### Проблемы в Linux

**«Недействительный ключ PGP»**

- Убедитесь, что файл ключа имеет формат ASCII Armor
- С помощью `gpg --show-keys <key.asc>` проверьте, не истёк ли срок действия ключа
- Убедитесь, что пароль указан верно

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
