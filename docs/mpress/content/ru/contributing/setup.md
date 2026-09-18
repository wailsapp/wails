---
title: "Настройка среды разработки"
description: "Настройте среду для разработки Wails v3"
slug: "contributing/setup"
sourcePath: "contributing/setup.md"
---

## Настройка среды разработки

В этом руководстве описана полная настройка среды для работы над Wails v3.

## Необходимые инструменты

### Разработка на Go

1. **Установите Go 1.25 или более поздней версии:**
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

Требуется только для примеров интеграции с фронтендом.

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
2. Среда выполнения WebView2 (предустановлена в Windows 11; для Windows 10 её можно [скачать](https://developer.microsoft.com/en-us/microsoft-edge/webview2/))
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
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails

# Add upstream remote
git remote add upstream https://github.com/wailsapp/wails.git

# Verify remotes
git remote -v
```

### Сборка CLI Wails

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
  - Настройте наблюдение за файлами для `goimports`
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

### Тесты определённого пакета

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

Документация Wails v3 написана с использованием M-Press. Исходные файлы на английском языке находятся в  
`docs/mpress/content/`, а переводы — в каталогах соответствующих языков, например  
`fr/` и `id/`.

Предварительно просмотрите и проверьте изменения в документации из корневого каталога репозитория:

```bash
go install github.com/leaanthony/mpress/cmd/mpress@v1.0.17
mpress dev
mpress build --strict --no-purge-css
mpress check
# Use `python` on Windows if `python3` is unavailable.
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

Рабочий сайт является статическим. Для локальной работы с документацией не требуются Node.js, поставщик услуг перевода и учетные данные Cloudflare.

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

Для отладки кода конкретной платформы необходимы платформенные инструменты:

- **macOS:** Xcode Instruments
- **Windows:** Visual Studio Debugger
- **Linux:** GDB

## Распространённые проблемы

### «command not found: wails3»

Добавьте каталог Wails в PATH или используйте `./wails3` из корневого каталога проекта.

### «webkitgtk-6.0 не найден» или «webkit2gtk не найден» (Linux)

Установите пакеты для разработки для того стека, под который выполняется сборка:

```bash
# Default GTK4 stack (Debian/Ubuntu):
sudo apt install libwebkitgtk-6.0-dev

# Legacy GTK3 path:
sudo apt install libwebkit2gtk-4.1-dev
```

### Сбой сборки из-за ошибок модулей Go

```bash
cd v3
go mod tidy
go mod download
```

### Ошибки «CGO_ENABLED» в Windows

Убедитесь, что компилятор C (MinGW-w64 из MSYS2) доступен через PATH.

## Дальнейшие действия

- Ознакомьтесь со [стандартами написания кода](/contributing/standards/)
- Изучите [техническую документацию](/contributing/)
- Найдите задачу, над которой можно поработать: [задачи для начинающих](https://github.com/wailsapp/wails/labels/good%20first%20issue)
