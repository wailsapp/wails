---
title: "Установка"
description: "Установите Wails и подготовьте среду для создания приложений"
slug: "quick-start/installation"
sourcePath: "quick-start/installation.md"
---

## Быстрая установка (5 минут)

@note{type="tip" title="Кратко — для опытных разработчиков"}
```bash
# Install Go 1.26+, then:
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0
wails3 setup   # Interactive setup wizard (experimental)
```

Также можно выполнить проверку вручную с помощью `wails3 doctor`. [Перейти к первому приложению →](/quick-start/first-app/)

@end

## Пошаговая установка

@steps
### Установите Go (обязательно)
Для Wails v3 требуется инструментарий сборки Go 1.26 или более поздней версии. Установите актуальный патч-релиз Go 1.26. Более старая программа запуска Go может автоматически загрузить подходящий инструментарий при `GOTOOLCHAIN=auto`; если автоматическая загрузка отключена, установите необходимый инструментарий самостоятельно. Объявление `go.mod` в самом стартовом проекте не снижает требования зависимости Wails.

@tabs{sync-key="os"}
[Windows]
Скачайте установщик для Windows с сайта **[go.dev/dl](https://go.dev/dl/)** и запустите его.

**Проверьте установку:**

```powershell
go version  # Should show 1.26 or later
```

**Проверьте PATH:**

```powershell
$env:PATH -split ';' | Where-Object { $_ -like '*\go\bin' }
```

Если вывод пуст, добавьте `C:\Users\YourName\go\bin` в PATH.

[macOS]
**Вариант 1: официальный установщик**

Скачайте установщик для macOS (файл .pkg) с сайта **[go.dev/dl](https://go.dev/dl/)** и запустите его.

**Вариант 2: Homebrew**

```bash
brew install go
```

**Проверьте установку:**

```bash
go version  # Should show 1.26 or later
echo $PATH | grep go/bin  # Should show ~/go/bin
```

Если `~/go/bin` отсутствует в PATH, добавьте его в `~/.zshrc` или `~/.bash_profile`:

```bash
export PATH=$PATH:~/go/bin
```

[Linux]
**Вариант 1: официальный tar-архив**

Скачайте tar-архив для Linux с сайта **[go.dev/dl](https://go.dev/dl/)**, затем выполните:

```bash
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.26.0.linux-amd64.tar.gz
```

**Вариант 2: менеджер пакетов**

```bash
# Ubuntu/Debian
sudo apt install golang-go

# Fedora
sudo dnf install golang

# Arch
sudo pacman -S go
```

**Добавьте в PATH** (добавьте в `~/.bashrc` или `~/.zshrc`):

```bash
export PATH=$PATH:/usr/local/go/bin:~/go/bin
source ~/.bashrc  # Reload
```

**Проверьте:**

```bash
go version
echo $PATH | grep go/bin
```

@end

### Установите зависимости для платформы
@tabs{sync-key="os"}
[Windows]
**Среда выполнения WebView2** (обычно уже установлена)

В Windows 10/11 WebView2 включён по умолчанию. Если он отсутствует:

- Скачайте с сайта [Microsoft](https://developer.microsoft.com/microsoft-edge/webview2/)
- Или запустите `wails3 doctor` позже — мастер подскажет, что делать

**Готово!** Другие зависимости не требуются.

@note{type="tip" title="Совет по повышению производительности в Windows 11"}
Рассмотрите возможность хранения проектов на [Dev Drive](https://learn.microsoft.com/en-us/windows/dev-drive/). Диски Dev Drive оптимизированы для рабочих нагрузок разработчиков и могут значительно сократить время сборки и повысить скорость доступа к диску — вплоть до 30%.

@end

[macOS]
**Инструменты командной строки Xcode** (обязательно)

```bash
xcode-select --install
```

В появившемся диалоговом окне нажмите «Установить».

**Проверьте:**

```bash
xcode-select -p  # Should show /Library/Developer/CommandLineTools
```

**Готово!** WebKit входит в состав macOS по умолчанию.

[Linux]
**Инструменты сборки и WebKit**

@note{type="caution" title="Минимальные версии дистрибутивов"}
По умолчанию для Wails v3 требуется **WebKitGTK 6.0**. Для дистрибутивов, в которых доступен только WebKit2GTK 4.1 — Ubuntu 22.04 LTS, Debian 12, Fedora ≤ 39, RHEL 9.x, — сборку необходимо выполнять с явным включением устаревшего стека с помощью `-tags gtk3`. Более старые выпуски, в которых доступен только WebKit2GTK 4.0 (Ubuntu 20.04, Debian 11, RHEL 8), не поддерживаются.

@end

@tabs{sync-key="distro"}
[Ubuntu/Debian]
Для использования стека GTK4 по умолчанию требуется Ubuntu 24.04+ или Debian 13+.

```bash
sudo apt update
sudo apt install build-essential pkg-config libgtk-4-dev libwebkitgtk-6.0-dev
```

[Fedora]
```bash
sudo dnf install gcc pkg-config gtk4-devel webkitgtk6.0-devel
```

[Arch]
```bash
sudo pacman -S base-devel gtk4 webkitgtk-6.0
```

[openSUSE]
```bash
sudo zypper install gcc pkg-config gtk4-devel webkitgtk-6_0-devel
```

[Gentoo]
```bash
sudo emerge --ask net-libs/webkit-gtk:6
```

[NixOS]
Добавьте в `shell.nix` или `devShell`:

```nix
buildInputs = with pkgs; [ webkitgtk_6_0 gtk4 pkg-config gcc ];
```

[Другие]
После установки Wails запустите `wails3 doctor` — команда покажет точный список пакетов, необходимых для вашего дистрибутива.

@end

@note{type="info" title="Устаревший стек GTK3"}
Если в целевом дистрибутиве пока нет WebKitGTK 6.0 (например, в Ubuntu 22.04 LTS или Debian 12), вместо него установите библиотеки разработки GTK3 и WebKit2GTK 4.1 (`libgtk-3-dev libwebkit2gtk-4.1-dev` в Debian/Ubuntu; эквивалентные пакеты в других дистрибутивах) и выполняйте сборку с `wails3 build -tags gtk3`. Устаревший вариант поддерживается в линейке v3.0.x и будет удалён в v3.1. Подробности см. в разделе [«Пакетирование для Linux — поддержка устаревшего GTK3»](/guides/build/linux/#legacy-gtk3-support).

@end

@end

### Установите CLI Wails
```bash
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0
```

Команда `wails3` будет установлена в `~/go/bin` (или в `%USERPROFILE%\go\bin` в Windows).

### Запустите мастер настройки (рекомендуется)
```bash
wails3 setup
```

Мастер настройки проверит зависимости, поможет установить отсутствующие и настроит параметры проекта по умолчанию.

@note{type="caution" title="Экспериментальная функция"}
Мастер настройки появился недавно и протестирован преимущественно в Linux. Если возникнут проблемы, [сообщите о них](https://github.com/wailsapp/wails/issues/4904) и вместо мастера используйте `wails3 doctor`.

@end

### Проверьте установку
```bash
wails3 doctor
```

**Ожидаемый вывод (или аналогичный):**

```
Wails (v3.0.0)  Wails Doctor

# System

┌──────────────────────────────────────────────────┐
| Name          | MacOS                            |
| Version       | 26.0                             |
| ID            | 25A354                           |
| Branding      | MacOS 26.0                       |
| Platform      | darwin                           |
| Architecture  | arm64                            |
| Apple Silicon | true                             |
| CPU           | Apple M2 Pro                     |
| CPU 1         | Apple M2 Pro                     |
| CPU 2         | Apple M2 Pro                     |
| GPU           | 16 cores, Metal Support: Metal 4 |
| Memory        | 16 GB                            |
└──────────────────────────────────────────────────┘

# Build Environment

┌─────────────┬─────────────────┐
| Wails CLI   | v3.0.0          |
| Go Version  | go1.26.0        |
└─────────────┴─────────────────┘

# Dependencies

┌─────────────────┬─────────────────────────────────────────────────┐
| npm             | 11.6.2                                          |
| *NSIS           | Not Installed. Install with `brew install...`.  |
| Xcode cli tools | 2412                                            |
└─────────────────┴─────────────────────────────────────────────────┘

# Checking for issues

SUCCESS No issues found

# Diagnosis

SUCCESS Your system is ready for Wails development!
```

@note{type="info" title="Если команда `wails3` не найдена"}
Ваш каталог `~/go/bin` отсутствует в PATH. Исправьте это, как описано выше на шаге 1, затем перезапустите терминал.

@end

### Установите npm (необязательно, но рекомендуется)
В большинстве шаблонов Wails для инструментов фронтенда используется npm.

@tabs{sync-key="os"}
[Windows]
Скачайте установщик с сайта [nodejs.org](https://nodejs.org/) и запустите его.

**Проверьте:**

```powershell
npm --version
```

[macOS]
**Вариант 1: официальный установщик** Скачайте его с сайта [nodejs.org](https://nodejs.org/)

**Вариант 2: Homebrew**

```bash
brew install node
```

**Проверьте:**

```bash
npm --version
```

[Linux]
**Вариант 1: NodeSource**

```bash
curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
sudo apt-get install -y nodejs  # Ubuntu/Debian
```

**Вариант 2: менеджер пакетов**

```bash
sudo dnf install nodejs  # Fedora
sudo pacman -S nodejs npm  # Arch
```

**Проверка:**

```bash
npm --version
```

@end

@note{type="tip" title="Альтернативные менеджеры пакетов"}
Предпочитаете `pnpm`, `yarn` или `bun`? Без проблем! Просто укажите предпочитаемый инструмент в `Taskfile.yml` своего проекта.

@end

@end

## Устранение неполадок

### Команда `wails3` не найдена

**Причина:** `~/go/bin` (или `%USERPROFILE%\go\bin`) отсутствует в PATH.

**Решение:**

@tabs{sync-key="os"}
[Windows]
1. Откройте «Переменные среды» (найдите их через меню «Пуск»)
2. В разделе «Переменные пользователя» найдите `Path`
3. Нажмите «Изменить» → «Создать»
4. Добавьте: `C:\Users\YourName\go\bin` (замените `YourName`)
5. Нажмите «ОК» во всех диалоговых окнах
6. **Перезапустите терминал**

**Проверка:**

```powershell
$env:PATH -split ';' | Where-Object { $_ -like '*\go\bin' }
```

[macOS/Linux]
Добавьте в `~/.zshrc` (macOS) или `~/.bashrc` (Linux):

```bash
export PATH=$PATH:~/go/bin
```

Перезагрузите конфигурацию:

```bash
source ~/.zshrc  # or ~/.bashrc
```

**Проверка:**

```bash
echo $PATH | grep go/bin
wails3 version
```

@end

---

#### `wails3 doctor` сообщает об отсутствующих зависимостях

**Linux:** в выводе точно указано, какие пакеты нужно установить. Пример:

```
❌ webkit2gtk not found
   Install with: sudo apt install libwebkit2gtk-4.1-dev
```

**Windows:** если отсутствует WebView2:

- Скачайте его с сайта [Microsoft](https://developer.microsoft.com/microsoft-edge/webview2/)
- Или он будет установлен автоматически при запуске первого приложения

**macOS:** если отсутствуют инструменты Xcode:

```bash
xcode-select --install
```

---

#### Слишком старая версия Go

Для Wails v3 требуется Go 1.26 или более поздней версии. Если у вас установлена более старая версия:

@tabs{sync-key="os"}
[Windows/macOS]
Скачайте последнюю версию с [go.dev/dl](https://go.dev/dl/) и переустановите её.

[Linux]
Скачайте последний tar-архив с [go.dev/dl](https://go.dev/dl/), затем выполните:

```bash
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.26.0.linux-amd64.tar.gz
```

@end

## Версия для разработки (самая свежая)

Хотите использовать самый свежий код из основной ветки разработки? Вы получите доступ к новым функциям и исправлениям до их выпуска, но рискуете столкнуться с ошибками и несовместимыми изменениями. Рекомендуется только участникам проекта и тем, кому необходимо тестировать готовящиеся функции.

```bash
git clone https://github.com/wailsapp/wails.git
cd wails
git checkout master
cd v3/cmd/wails3
go install
```

@note{type="caution" title="Версия для разработки"}
- Может содержать ошибки или несовместимые изменения
- Созданные проекты будут использовать директиву `replace`, указывающую на локальную версию Wails
- Рекомендуется только участникам проекта или для тестирования новых функций

@end

## Дальнейшие действия

**Установка завершена!** Ваша система готова к разработке с Wails.

@cards{cols="1"}
🚀 Создайте первое приложение
Создайте работающее приложение за 10 минут.

[Руководство по созданию первого приложения →](/quick-start/first-app/)

@end

@cards{cols="1"}
📖 Изучите шаблоны
Посмотрите, что доступно сразу после установки.

```bash
wails3 init -l  # List templates
```

@end

---

**Возникли проблемы?** Задайте вопрос в [Discord](https://discord.gg/JDdSxwjhGf) или [создайте обращение](https://github.com/wailsapp/wails/issues).
