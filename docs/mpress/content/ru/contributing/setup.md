---
title: "Настройка среды разработки"
description: "Настройте среду для разработки Wails v3"
slug: "contributing/setup"
sourcePath: "contributing/setup.md"
---

## Настройка среды разработки

В этом руководстве описана полная настройка среды разработки для работы над Wails v3.

## Необходимые инструменты

### Разработка на Go

1. **Установите Go 1.26 или более поздней версии:**
  ```bash
  # Download from https://go.dev/dl/
  go version  # Verify installation
  ```


2. **Настройте среду Go:**
  ```bash
  # Add to your shell profile (.bashrc, .zshrc, etc.)
  export GOPATH=$HOME/go
  export PATH=$PATH:$GOPATH/bin
  ```


3. **Установите полезные инструменты Go:**
  ```bash
  go install golang.org/x/tools/cmd/goimports@latest
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  ```


### Node.js и npm

Требуются только для примеров интеграции с фронтендом.

```bash
# Install Node.js 20+ and npm
node --version  # Should be 20+
npm --version
```

### Зависимости для отдельных платформ

**macOS:**

```bash
# Install Xcode Command Line Tools
xcode-select --install

# Verify installation
xcode-select -p  # Should output a path
```

**Windows:**

1. Установите [MSYS2](https://www.msys2.org/), чтобы получить Unix-подобную среду
2. Среда выполнения WebView2 (предустановлена в Windows 11; для Windows [скачайте её](https://developer.microsoft.com/en-us/microsoft-edge/webview2/) 10)
3. Необязательно: установите [Git for Windows](https://git-scm.com/download/win)

**Linux (Debian/Ubuntu):**

```bash
sudo apt update
# Default GTK4 + WebKitGTK 6.0 stack (Ubuntu 24.04+ / Debian 13+)
sudo apt install build-essential pkg-config libgtk-4-dev libwebkitgtk-6.0-dev
# For the legacy -tags gtk3 path:
sudo apt install libgtk-3-dev libwebkit2gtk-4.1-dev
```

**Linux (Fedora/RHEL):**

```bash
# Default GTK4 stack
sudo dnf install gcc pkg-config gtk4-devel webkitgtk6.0-devel
# Legacy GTK3 path:
sudo dnf install gtk3-devel webkit2gtk4.1-devel
```

**Linux (Arch):**

```bash
# Default GTK4 stack
sudo pacman -S base-devel gtk4 webkitgtk-6.0
# Legacy GTK3 path:
sudo pacman -S gtk3 webkit2gtk-4.1
```

## Настройка репозитория

### Клонирование и настройка

```bash
# Clone your fork
git clone --branch releases/v3.0.0 https://github.com/YOUR_USERNAME/wails.git
cd wails

# Add upstream remote
git remote add upstream https://github.com/wailsapp/wails.git

# Verify remotes
git remote -v
```

### Сборка Wails CLI

```bash
# Navigate to v3 directory
cd v3

# Build the CLI
go build -o ../wails3 ./cmd/wails3

# Test the build
cd ..
./wails3 version
```

### Добавление в PATH (необязательно)

**Linux/macOS:**

```bash
# Add to ~/.bashrc or ~/.zshrc
export PATH=$PATH:/path/to/wails
```

**Windows:**

Добавьте каталог Wails в переменную среды PATH через свойства системы.

## Настройка IDE

### VS Code (рекомендуется)

1. **Установите VS Code:** [Скачать](https://code.visualstudio.com/)

2. **Установите расширения:**
  - Go (от команды Go в Google)
  - ESLint
  - Prettier
  - MDX (для документации)


3. **Настройте параметры рабочей области** (`.vscode/settings.json`):
  ```json
  {
    "go.useLanguageServer": true,
    "go.lintTool": "golangci-lint",
    "go.lintOnSave": "workspace",
    "editor.formatOnSave": true,
    "go.formatTool": "goimports"
  }
  ```


### GoLand

1. **Установите GoLand:** [Скачать](https://www.jetbrains.com/go/)

2. **Настройте:**
  - Включите поддержку модулей Go
  - Настройте наблюдатели за файлами для `goimports`
  - Настройте стиль кода в соответствии с соглашениями проекта


## Проверка настройки

Выполните следующие команды, чтобы убедиться, что всё работает:

```bash
# Go version check
go version

# Build Wails
cd v3
go build ./cmd/wails3

# Run tests
go test ./pkg/...

# Create a test app
cd ..
./wails3 init -n mytest -t vanilla
cd mytest
../wails3 dev
```

Если тестовое приложение собирается и запускается, среда готова!

## Запуск тестов

### Модульные тесты

```bash
cd v3
go test ./...
cd ..
```

### Тесты отдельных пакетов

```bash
cd v3
go test ./pkg/application
go test ./pkg/events -v  # Verbose output
```

### Запуск с измерением покрытия

```bash
cd v3
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Запуск с детектором гонок данных

```bash
cd v3
go test ./... -race
```

## Работа с документацией

Авторы документации могут использовать M-Press без установки Wails CLI и зависимостей для сборки настольных приложений. Инструкции по скачиванию закреплённых версий исполняемых файлов, установке, интерактивному предварительному просмотру, проверке и процессу создания PR для ветки выпуска см. в разделе [Создание документации с помощью M-Press](/contributing/documentation/).

Англоязычные файлы Markdown находятся в `docs/mpress/content/`, а переводы — в каталогах соответствующих языков. Команда `mpress dev`, запущенная из корня репозитория, открывает предварительный просмотр по адресу `http://127.0.0.1:3000`. Перед выполнением описанных в руководстве команд проверки остановите её с помощью Ctrl+C.

## Отладка

### Отладка кода Go

**VS Code:**

Создайте `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug Wails CLI",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/v3/cmd/wails3",
      "args": ["dev"]
    }
  ]
}
```

**Командная строка:**

```bash
# Use Delve debugger
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug ./cmd/wails3 -- dev
```

### Отладка платформенного кода

Для отладки кода конкретной платформы необходимы соответствующие инструменты:

- **macOS:** Xcode Instruments
- **Windows:** отладчик Visual Studio
- **Linux:** GDB

## Распространённые проблемы

### «command not found: wails3»

Добавьте каталог Wails в PATH или используйте `./wails3` из корня проекта.

### «webkitgtk-6.0 not found» или «webkit2gtk not found» (Linux)

Установите пакеты для разработки, соответствующие стеку, с которым выполняется сборка:

```bash
# Default GTK4 stack (Debian/Ubuntu):
sudo apt install libwebkitgtk-6.0-dev

# Legacy GTK3 path:
sudo apt install libwebkit2gtk-4.1-dev
```

### Сборка завершается ошибками модулей Go

```bash
cd v3
go mod tidy
go mod download
```

### Ошибки «CGO_ENABLED» в Windows

Убедитесь, что компилятор C (MinGW-w64 через MSYS2) доступен в PATH.

## Дальнейшие действия

- Ознакомьтесь со [Стандартами написания кода](/contributing/standards/)
- Изучите [техническую документацию](/contributing/)
- Найдите задачу для работы: [задачи для начинающих](https://github.com/wailsapp/wails/labels/good%20first%20issue)
