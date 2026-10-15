---
title: "Система сборки"
description: "Как Wails собирает и упаковывает ваше приложение"
slug: "concepts/build-system"
sourcePath: "concepts/build-system.md"
---

## Единая система сборки

Wails предоставляет **единую систему сборки**, которая одной командой компилирует код Go, объединяет ресурсы фронтенда, встраивает всё в один исполняемый файл и выполняет сборку с учётом особенностей платформы.

```bash
wails3 build
```

**Результат:** нативный исполняемый файл со всем встроенным содержимым.

## Обзор процесса сборки

**[Место для диаграммы процесса сборки]**

## Этапы сборки

### 1. Этап анализа

Wails сканирует ваш код Go, чтобы получить сведения о сервисах:

```go
type GreetService struct {
    prefix string
}

func (g *GreetService) Greet(name string) string {
    return g.prefix + name + "!"
}
```

**Что извлекает Wails:**

- Имя сервиса: `GreetService`
- Имя метода: `Greet`
- Типы параметров: `string`
- Возвращаемые типы: `string`

**Используется для:** создания привязок TypeScript

### 2. Этап генерации

#### Привязки TypeScript

Wails создаёт типобезопасные привязки:

```javascript
// Auto-generated: frontend/bindings/<full-go-import-path>/greetservice.js
// (TypeScript is also generated when you pass `-ts`. The shape below is the real
// runtime call format — numeric IDs via $Call.ByID, imported from /wails/runtime.js.)
import { Call as $Call } from "/wails/runtime.js";

export function Greet($0) {
    return $Call.ByID(1234567890, $0);
}
```

**Преимущества:**

- Полная типобезопасность
- Автодополнение в IDE
- Выявление ошибок во время компиляции
- Комментарии JSDoc

#### Сборка фронтенда

Запускается сборщик фронтенда (Vite, webpack и т. п.):

```bash
# Vite example
vite build --outDir dist
```

**Что происходит:**

- Компиляция JavaScript/TypeScript
- Обработка и минификация CSS
- Оптимизация ресурсов
- Создание карт исходного кода (только при разработке)
- Вывод в `frontend/dist/`

### 3. Этап компиляции

#### Компиляция Go

Код Go компилируется с оптимизациями:

```bash
go build -ldflags="-s -w" -o myapp.exe
```

**Флаги:**

- `-s`: удаление таблицы символов
- `-w`: удаление отладочной информации DWARF
- Результат: исполняемый файл меньшего размера (уменьшение примерно на 30 %)

**Особенности платформ:**

- Windows: `.exe` со встроенным значком
- macOS: структура пакета `.app`
- Linux: исполняемый файл ELF

#### Встраивание ресурсов

Ресурсы фронтенда встраиваются в исполняемый файл Go:

```go
//go:embed frontend/dist
var assets embed.FS
```

**Результат:** один исполняемый файл со всем содержимым.

### 4. Результат

**Один нативный исполняемый файл:**

- Windows: `myapp.exe` (~15 МБ)
- macOS: `myapp.app` (~15 МБ)
- Linux: `myapp` (~15 МБ)

**Без зависимостей** (кроме системного WebView).

## Разработка и рабочая сборка

@tabs{sync-key="mode"}
[Разработка (wails3 dev)]
**Оптимизировано для скорости:**

```bash
wails3 dev
```

**Что происходит:**

1. Запускается сервер разработки фронтенда (по умолчанию Vite на порту 9245)
2. Go компилируется без оптимизаций
3. Запускается приложение, подключённое к серверу разработки
4. Включается горячая перезагрузка
5. Включаются карты исходного кода

**Характеристики:**

- **Быстрая повторная сборка** (&lt;1 с при изменениях во фронтенде)
- **Без встраивания ресурсов** (они предоставляются сервером разработки)
- **Отладочные символы** включены
- **Карты исходного кода** включены
- **Подробное журналирование**

**Размер файла:** больше (~50 МБ с отладочными символами)

[Продакшен (wails3 build)]
**Оптимизировано по размеру и производительности:**

```bash
wails3 build
```

**Что происходит:**

1. Фронтенд собирается для продакшена (с минификацией)
2. Код Go компилируется с оптимизациями
3. Отладочные символы удаляются
4. Ресурсы встраиваются
5. Создаётся один исполняемый файл

**Характеристики:**

- **Оптимизированный код** (минифицированный, с удалённым неиспользуемым кодом)
- **Встроенные ресурсы** (без внешних файлов)
- **Отладочные символы удалены**
- **Без карт исходного кода**
- **Минимальное журналирование**

**Размер файла:** меньше (~15 МБ)

@end

## Команды сборки

### Базовая сборка

```bash
wails3 build
```

**Результат:** `bin/<APP_NAME>` (или `bin/<APP_NAME>.exe` в Windows). Каталог `bin/` находится в корне проекта.

`wails3 build` — это тонкая оболочка над `wails3 task build`. Единственный передаваемый ею флаг времени сборки — `--tags`, который становится переменной Taskfile `EXTRA_TAGS`:

```bash
# Build with extra Go build tags
wails3 build --tags "myfeature,gtk4"
```

У `wails3 build` нет флагов `-platform`, `-o`, `-skipbindings`, `-clean`, `-debug`, `-devbuild`, `-icon`, `-ldflags` и `-package`. Кросс-компиляция, пути вывода, значки и создание пакетов настраиваются через Taskfile проекта (`Taskfile.yml` + `build/config.yml`).

### Кроссплатформенные и платформозависимые сборки

Сборки для отдельных платформ доступны как задачи Taskfile в пространствах имён `darwin:` / `windows:` / `linux:` (определённых в `build/Taskfile.<platform>.yml`). Например:

```bash
# macOS — universal binary
wails3 task darwin:build:universal

# macOS — current arch
wails3 task darwin:build

# Windows
wails3 task windows:build

# Linux
wails3 task linux:build
```

Чтобы просмотреть все доступные задачи текущего проекта:

```bash
wails3 task --list
```

### Значки и создание пакетов

Создайте значки для платформ (`build/icons.icns`, `build/icon.ico` и т. д.) из исходного PNG-файла:

```bash
wails3 generate icons -input appicon.png
```

Соберите установщики и пакеты для отдельных платформ:

```bash
wails3 package           # uses the current Go build env
wails3 task linux:create:deb
wails3 task windows:package
wails3 task darwin:package:universal
```

## Конфигурация сборки

### Taskfile.yml

В проектах Wails 3 для оркестрации сборки используется [Taskfile](https://taskfile.dev/). Корневой файл `Taskfile.yml` включает файлы задач для отдельных платформ из `build/`:

```yaml
# Taskfile.yml (excerpt — the real templates are richer)
version: '3'

includes:
  common: ./build/Taskfile.yml
  darwin: ./build/Taskfile.darwin.yml
  windows: ./build/Taskfile.windows.yml
  linux: ./build/Taskfile.linux.yml

tasks:
  build:
    desc: Build the application
    cmds:
      - task: "{{OS}}:build"
```

Запускайте задачи с помощью `wails3 task <name>` или `task <name>`:

```bash
wails3 task windows:build
wails3 task darwin:package:universal
wails3 task linux:create:appimage
```

### Конфигурация проекта: `build/config.yml`

Метаданные проекта (имя, идентификатор, версия, значения info-plist, параметры NSIS, поля `.desktop`, пользовательские протоколы и т. д.) находятся в `build/config.yml`. Taskfile считывает этот файл при создании значков, манифестов, установщиков и других подобных компонентов. В Wails 3 файла `build/build.json` **нет**.

```yaml
# build/config.yml (illustrative)
info:
  productName: "My App"
  productIdentifier: "com.example.myapp"
  productVersion: "1.0.0"
  companyName: "Example Ltd."
  productDescription: "An application built with Wails"
```

Выполните `wails3 generate build-assets` (или `wails3 update build-assets`), чтобы обновить из этой конфигурации платформозависимые ресурсы сборки.

## Встраивание ресурсов

### Как это работает

Wails использует пакет Go `embed`:

```go
package main

import (
    "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name:   "My App",
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
    })
    
    app.Window.New()
    app.Run()
}
```

**Во время сборки:**

1. Фронтенд собирается в `frontend/dist/`
2. Директива `//go:embed` включает файлы
3. Файлы компилируются в исполняемый файл
4. Исполняемый файл содержит всё необходимое

**Во время выполнения:**

1. Приложение запускается
2. Ресурсы загружаются из памяти
3. Для ресурсов не требуется дисковый ввод-вывод
4. Быстрая загрузка

### Пользовательские ресурсы

Встройте дополнительные файлы:

```go
//go:embed frontend/dist
var frontendAssets embed.FS

//go:embed data/*.json
var dataAssets embed.FS

//go:embed templates/*.html
var templateAssets embed.FS
```

## Оптимизация сборки

### Оптимизация фронтенда

**Vite (по умолчанию):**

```javascript
// vite.config.js
export default {
  build: {
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: true,  // Remove console.log
        drop_debugger: true,
      },
    },
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['react', 'react-dom'],  // Separate vendor bundle
        },
      },
    },
  },
}
```

**Результаты:**

- JavaScript минифицирован (размер уменьшен примерно на 70 %)
- CSS минифицирован (размер уменьшен примерно на 60 %)
- Изображения оптимизированы
- Выполнено удаление неиспользуемого кода

### Оптимизация Go

**Флаги компилятора:**

```bash
-ldflags="-s -w"
```

- `-s`: удалить таблицу символов (уменьшение примерно на 10 %)
- `-w`: удалить отладочную информацию DWARF (уменьшение примерно на 20 %)

**Дополнительные оптимизации:**

```bash
-ldflags="-s -w -X main.version=1.0.0"
```

- `-X`: задать значения переменных во время сборки
- Полезно для номеров версий и дат сборки

### Сжатие исполняемого файла

**UPX (необязательно):**

```bash
# After building
upx --best bin/myapp.exe
```

**Результаты:**

- Уменьшение размера примерно на 50 %
- Незначительное замедление запуска (примерно на 100 мс)
- Не рекомендуется для macOS из-за проблем с подписыванием кода

## Сборки для отдельных платформ

### Windows

**Результат:** `myapp.exe`

**Включает:**

- Значок приложения
- Сведения о версии
- Манифест (настройки UAC)

**Значок:**

```bash
# Generate platform icons from a source PNG
wails3 generate icons -input appicon.png -windowsfilename build/icon.ico
```

Затем на этапе `tool package` для Windows созданный файл `.ico` встраивается в исполняемый файл.

**Манифест:**

```xml
<!-- build/windows/manifest.xml -->
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">
  <assemblyIdentity version="1.0.0.0" name="MyApp"/>
  <trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
      <requestedPrivileges>
        <requestedExecutionLevel level="asInvoker" uiAccess="false"/>
      </requestedPrivileges>
    </security>
  </trustInfo>
</assembly>
```

### macOS

**Результат:** `myapp.app` (пакет приложения)

**Структура:**

```
myapp.app/
├── Contents/
│   ├── Info.plist          # App metadata
│   ├── MacOS/
│   │   └── myapp           # Binary
│   ├── Resources/
│   │   └── icon.icns       # Icon
│   └── _CodeSignature/     # Code signature (if signed)
```

**Info.plist:**

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleName</key>
    <string>My App</string>
    <key>CFBundleIdentifier</key>
    <string>com.example.myapp</string>
    <key>CFBundleVersion</key>
    <string>1.0.0</string>
</dict>
</plist>
```

**Универсальный исполняемый файл:**

Taskfile для macOS содержит задачу `darwin:build:universal` (а также `darwin:package:universal`), которая выполняет сборку для обеих архитектур и объединяет результаты с помощью `wails3 tool lipo`:

```bash
wails3 task darwin:build:universal
```

### Linux

**Результат:** `myapp` (исполняемый файл ELF)

**Зависимости:**

- GTK3
- WebKitGTK

**Файл рабочего стола:**

```ini
# myapp.desktop
[Desktop Entry]
Name=My App
Exec=/usr/bin/myapp
Icon=myapp
Type=Application
Categories=Utility;
```

**Установка:**

```bash
# Copy binary
sudo cp myapp /usr/bin/

# Copy desktop file
sudo cp myapp.desktop /usr/share/applications/

# Copy icon
sudo cp icon.png /usr/share/icons/hicolor/256x256/apps/myapp.png
```

## Производительность сборки

### Типичное время сборки

| Этап | Время | Примечания |
| --- | --- | --- |
| Анализ | &lt;1 с | Сканирование кода Go |
| Создание привязок | &lt;1 с | Генерация TypeScript |
| Сборка фронтенда | 5-30 с | Зависит от размера проекта |
| Компиляция Go | 2-10 с | Зависит от объёма кода |
| Встраивание ресурсов | &lt;1 с | Встраивание фронтенда |
| **Итого** | **10-45 с** | Первая сборка |
| **Инкрементная сборка** | **5-15 с** | Последующие сборки |

### Ускорение сборки

**1. Используйте кеш сборки:**

```bash
# Go build cache is automatic
# Frontend cache (Vite)
npm run build  # Uses cache by default
```

**2. Запускайте только необходимые задачи:**

```bash
# Pick the specific Taskfile target you actually need
wails3 task common:build:frontend   # rebuild only the frontend
wails3 task windows:build           # rebuild only the Windows binary
```

**3. Используйте параллельные сборки (на нескольких машинах или в CI):**

В v3 кросс-компиляция между Linux, Windows и macOS обычно выполняется в контейнере Docker `wails-cross` или на выделенных исполнителях для каждой платформы — сама команда `wails3 build` создаёт сборку для ОС хоста. Поддерживаемые рабочие процессы описаны в разделе [«Кроссплатформенные сборки»](/guides/build/cross-platform/).

**4. Используйте более быстрые инструменты:**

```bash
# Use esbuild instead of webpack
# (Vite uses esbuild by default)
```

## Устранение неполадок

### Сбой сборки

**Признак:** `wails3 build` завершается с ошибкой

**Распространённые причины:**

1. **Ошибка компиляции Go**
  ```bash
  # Check Go code compiles
  go build
  ```


2. **Ошибка сборки фронтенда**
  ```bash
  # Check frontend builds
  cd frontend
  npm run build
  ```


3. **Отсутствующие зависимости**
  ```bash
  # Install dependencies
  npm install
  go mod download
  ```


### Слишком большой исполняемый файл

**Признак:** размер исполняемого файла превышает 50 МБ

**Решения:**

1. **Удалите отладочные символы** (поставляемый Taskfile уже передаёт `-ldflags="-s -w"` в `go build`).

2. **Проверьте встроенные ресурсы**
  ```bash
  # Remove unnecessary files from frontend/dist/
  # Check for large images, videos, etc.
  ```


3. **Используйте сжатие UPX**
  ```bash
  upx --best bin/myapp.exe
  ```


### Медленная сборка

**Признак:** сборка занимает более 1 минуты

**Решения:**

1. **Используйте кеш сборки**
  - Кеширование Go выполняется автоматически
  - Кеширование фронтенда (Vite) выполняется автоматически


2. **Запускайте только нужную задачу**
  ```bash
  wails3 task common:build:frontend
  wails3 task windows:build
  ```


3. **Оптимизируйте сборку фронтенда**
  ```javascript
  // vite.config.js
  export default {
    build: {
      minify: 'esbuild',  // Faster than terser
    },
  }
  ```


## Рекомендации

### ✅ Делайте

- **Используйте `wails3 dev` во время разработки** — это ускоряет итерации
- **Используйте `wails3 build` для выпусков** — это обеспечивает оптимизированный результат
- **Версионируйте сборки** — используйте `-ldflags` для встраивания версии
- **Тестируйте сборки на целевых платформах** — кросс-компиляция не идеальна
- **Обеспечьте быструю сборку фронтенда** — оптимизируйте конфигурацию сборщика
- **Используйте кеш сборки** — это ускоряет последующие сборки

### ❌ Не делайте

- **Не добавляйте каталог `build/` в коммиты** — добавьте его в `.gitignore`
- **Не пропускайте тестирование сборок** — всегда тестируйте их перед выпуском
- **Не встраивайте ненужные ресурсы** — сохраняйте небольшой размер исполняемых файлов
- **Не используйте отладочные сборки в рабочей среде** — используйте оптимизированные сборки
- **Не забывайте подписывать код** — это необходимо для распространения

## Следующие шаги

**Сборка приложений** — подробное руководство по сборке и упаковке [Подробнее →](/guides/build/building/)

**Кроссплатформенные сборки** — сборка для всех платформ на одном компьютере [Подробнее →](/guides/build/cross-platform/)

**Создание установщиков** — создание установщиков для конечных пользователей [Подробнее →](/guides/installers/)

---

**Есть вопросы о сборке?** Задайте их в [Discord](https://discord.gg/JDdSxwjhGf) или ознакомьтесь с [примерами сборки](https://github.com/wailsapp/wails/tree/master/v3/examples/build).
