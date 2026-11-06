---
title: "Среда выполнения фронтенда"
description: "Пакет среды выполнения JavaScript для интеграции фронтенда с Wails"
slug: "reference/frontend-runtime"
sourcePath: "reference/frontend-runtime.md"
---

Среда выполнения фронтенда Wails — это стандартная библиотека для приложений Wails. Она предоставляет ряд возможностей, которые можно использовать в приложениях, включая:

- Управление окнами
- Диалоговые окна
- Интеграция с браузером
- Буфер обмена
- Меню
- Сведения о системе
- События
- Контекстные меню
- Экраны
- WML (язык разметки Wails)

Среда выполнения необходима для интеграции Go с фронтендом. Интегрировать среду выполнения можно 2 способами:

- С помощью пакета `@wailsio/runtime`
- С помощью готового сборочного файла

## Использование пакета npm

Пакет `@wailsio/runtime` — это пакет JavaScript, предоставляющий фронтенду доступ к среде выполнения Wails. Он используется во всех стандартных шаблонах и является рекомендуемым способом интеграции среды выполнения в приложение. При использовании пакета `@wailsio/runtime` в приложение включаются только используемые части среды выполнения.

Пакет доступен в npm. Установить его можно с помощью следующей команды:

```shell
npm install --save @wailsio/runtime
```

## Использование готового сборочного файла

В некоторых проектах сборщик JavaScript не используется, поэтому для них может быть предпочтительнее готовая собранная версия среды выполнения. Эту версию можно создать локально с помощью следующей команды:

```shell
wails3 generate runtime
```

Команда создаст в текущем каталоге файл `runtime.js` (а также `runtime-debug.js`). Это модуль ES, который сценарии приложения могут импортировать так же, как пакет npm. Кроме того, API экспортируется в глобальный объект window, поэтому в простых приложениях его можно использовать следующим образом:

```html
<html>
    <head>
        <script type="module" src="./runtime.js"></script>
        <script>
            window.onload = function () {
                wails.Window.SetTitle("A new window title");
            }
        </script>
    </head>
    <!--- ... -->
</html>
```

@note{type="caution"}
Важно добавить атрибут `type="module"` в тег `<script>`, загружающий среду выполнения, и дождаться полной загрузки страницы, прежде чем вызывать API, поскольку сценарии с атрибутом `type="module"` выполняются асинхронно.

@end

## Инициализация

Помимо функций API, среда выполнения поддерживает контекстные меню и перетаскивание окон. Эти возможности будут работать должным образом только после инициализации среды выполнения. Даже если вы не используете API, обязательно добавьте в код фронтенда инструкцию импорта исключительно ради побочного эффекта:

```javascript
import "@wailsio/runtime";
```

Сборщик должен обнаружить наличие побочных эффектов и включить в сборку весь необходимый код инициализации.

@note{type="info"}
Если вы предпочитаете готовый сборочный файл, достаточно добавить тег script, как показано выше.

@end

## Плагин Vite для типизированных событий

Среда выполнения включает плагин Vite, который во время разработки обеспечивает поддержку HMR (горячей замены модулей) для типизированных событий.

### Настройка

Добавьте плагин в `vite.config.ts`:

```typescript
import { defineConfig } from 'vite'
import wails from '@wailsio/runtime/plugins/vite'

export default defineConfig({
  plugins: [wails("./bindings")],
})
```

Аргументом служит каталог сгенерированных привязок, указанный относительно проекта Vite. Перед сборкой фронтенда сгенерируйте привязки.

### Преимущества

- **Автоматическая перезагрузка**: привязки событий автоматически создаются заново и перезагружаются при выполнении `wails3 generate bindings`
- **Режим разработки**: без дополнительной настройки работает с `wails3 dev`, мгновенно применяя обновления
- **Безопасность типов**: полная поддержка TypeScript с автодополнением и проверкой типов

### Использование с регистрацией событий

Зарегистрируйте события в Go:

```go
type UserData struct {
    ID   string
    Name string
}

func init() {
    application.RegisterEvent[UserData]("user-updated")
}
```

Сгенерируйте привязки:

```bash
wails3 generate bindings

# Or, to include TypeScript definitions
wails3 generate bindings -ts

# For more options, see:
wails3 generate bindings -help
```

Используйте типизированные события во фронтенде:

```typescript
import { Events } from '@wailsio/runtime'
import { UserUpdated } from './bindings/events'

// Type-safe event with autocomplete
Events.Emit(UserUpdated, {
    ID: "123",
    Name: "John Doe"
})
```

## Справочник API

Среда выполнения разделена на модули, каждый из которых предоставляет определённые возможности. Импортируйте только необходимое:

```javascript
import { Events, Window, Clipboard } from '@wailsio/runtime'
```

### События

Система событий для взаимодействия между Go и JavaScript. Сгенерированные экспорты событий, такие как `UserUpdated`, являются строковыми константами, а не фабриками полезной нагрузки. Генератор дополняет `CustomEvents`, сопоставляя имена зарегистрированных событий с типами полезной нагрузки. Приведённые ниже типы (`WailsEventName`, `WailsEventCallback`, `WailsEventData`) экспортируются модулем Events.

#### On()

Регистрирует функцию обратного вызова для события.

```typescript
function On<E extends WailsEventName = WailsEventName>(eventName: E, callback: WailsEventCallback<E>): () => void
```

**Возвращает:** функцию отмены подписки

**Пример:**

```javascript
import { Events } from '@wailsio/runtime'

// Basic event listening
const unsubscribe = Events.On('user-logged-in', (event) => {
    console.log('User:', event.data.username)
})

// With typed events (TypeScript)
import { UserLogin } from './bindings/events'

Events.On(UserLogin, (event) => {
    // event.data is typed as UserLoginData
    console.log('User:', event.data.username)
})

// Later: unsubscribe()
```

#### Once()

Регистрирует функцию обратного вызова, которая выполняется только один раз.

```typescript
function Once<E extends WailsEventName = WailsEventName>(eventName: E, callback: WailsEventCallback<E>): () => void
```

**Пример:**

```javascript
import { Events } from '@wailsio/runtime'

Events.Once('app-ready', () => {
    console.log('App initialized')
})
```

#### Emit()

Отправляет событие в бэкенд Go или другие окна.

```typescript
function Emit<E extends WailsEventName = WailsEventName>(name: E, data: WailsEventData<E>): Promise<boolean>
function Emit<E extends WailsEventName = WailsEventName>(name: WailsEventData<E> extends null | void ? E : never): Promise<boolean>
```

**Возвращает:** Promise, разрешающийся значением `true`, если событие было отменено, и `false` в противном случае

**Пример:**

```javascript
import { Events } from '@wailsio/runtime'

// Basic event emission
const wasCancelled = await Events.Emit('button-clicked', { buttonId: 'submit' })

// With typed events (TypeScript)
import { UserLogin } from './bindings/events'

const cancelled = await Events.Emit(UserLogin, {
    UserID: "123",
    Username: "john_doe",
    LoginTime: new Date().toISOString()
})

if (cancelled) {
    console.log('Login was cancelled by a hook')
}
```

@note{type="info"}
Возвращаемое значение указывает, было ли событие отменено хуком. Большинство событий отменить нельзя, поэтому они всегда возвращают `false`.

@end

#### Off()

Удаляет обработчики событий.

```typescript
function Off(...eventNames: string[]): void
```

**Пример:**

```javascript
import { Events } from '@wailsio/runtime'

Events.Off('user-logged-in', 'user-logged-out')
```

#### OffAll()

Удаляет все обработчики событий.

```typescript
function OffAll(): void
```

### Окно

Методы управления окнами. Экспортом по умолчанию является текущее окно.

```javascript
import { Window } from '@wailsio/runtime'

// Current window
await Window.SetTitle('New Title')
await Window.Center()

// Get another window
const otherWindow = Window.Get('secondary')
await otherWindow.Show()
```

#### Видимость

**Show()** — показывает окно

```typescript
function Show(): Promise<void>
```

**Hide()** — скрывает окно

```typescript
function Hide(): Promise<void>
```

**Close()** — закрывает окно

```typescript
function Close(): Promise<void>
```

#### Размер и положение

**SetSize(width, height)** — задаёт размер окна

```typescript
function SetSize(width: number, height: number): Promise<void>
```

**Size()** — возвращает размер окна

```typescript
function Size(): Promise<{ width: number, height: number }>
```

**SetPosition(x, y)** — задаёт абсолютное положение

```typescript
function SetPosition(x: number, y: number): Promise<void>
```

**Position()** — возвращает абсолютное положение

```typescript
function Position(): Promise<{ x: number, y: number }>
```

**Center()** — размещает окно по центру

```typescript
function Center(): Promise<void>
```

**Пример:**

```javascript
import { Window } from '@wailsio/runtime'

// Resize and center
await Window.SetSize(800, 600)
await Window.Center()

// Get current size
const { width, height } = await Window.Size()
```

#### Состояние окна

**Minimise()** — сворачивает окно

```typescript
function Minimise(): Promise<void>
```

**Maximise()** — разворачивает окно

```typescript
function Maximise(): Promise<void>
```

**Fullscreen()** — переводит окно в полноэкранный режим

```typescript
function Fullscreen(): Promise<void>
```

**Restore()** — восстанавливает окно из свёрнутого, развёрнутого или полноэкранного состояния

```typescript
function Restore(): Promise<void>
```

**IsMinimised()** — проверяет, свёрнуто ли окно

```typescript
function IsMinimised(): Promise<boolean>
```

**IsMaximised()** — проверяет, развёрнуто ли окно

```typescript
function IsMaximised(): Promise<boolean>
```

**IsFullscreen()** — проверяет, находится ли окно в полноэкранном режиме

```typescript
function IsFullscreen(): Promise<boolean>
```

#### Свойства окна

**SetTitle(title)** — задаёт заголовок окна

```typescript
function SetTitle(title: string): Promise<void>
```

**Name()** — возвращает имя окна

```typescript
function Name(): Promise<string>
```

**SetBackgroundColour(r, g, b, a)** — задаёт цвет фона

```typescript
function SetBackgroundColour(r: number, g: number, b: number, a: number): Promise<void>
```

**SetAlwaysOnTop(alwaysOnTop)** — удерживает окно поверх остальных

```typescript
function SetAlwaysOnTop(alwaysOnTop: boolean): Promise<void>
```

**SetResizable(resizable)** — разрешает изменять размер окна

```typescript
function SetResizable(resizable: boolean): Promise<void>
```

#### Фокус и экран

**Focus()** — переводит фокус на окно

```typescript
function Focus(): Promise<void>
```

**IsFocused()** — проверяет, находится ли окно в фокусе

```typescript
function IsFocused(): Promise<boolean>
```

**GetScreen()** — возвращает экран, на котором находится окно

```typescript
function GetScreen(): Promise<Screen>
```

#### Содержимое

**Reload()** — перезагружает страницу

```typescript
function Reload(): Promise<void>
```

**ForceReload()** — принудительно перезагружает страницу (с очисткой кеша)

```typescript
function ForceReload(): Promise<void>
```

#### Масштаб

**SetZoom(level)** — задаёт уровень масштаба

```typescript
function SetZoom(level: number): Promise<void>
```

**GetZoom()** — возвращает уровень масштаба

```typescript
function GetZoom(): Promise<number>
```

**ZoomIn()** — увеличивает масштаб

```typescript
function ZoomIn(): Promise<void>
```

**ZoomOut()** — уменьшает масштаб

```typescript
function ZoomOut(): Promise<void>
```

**ZoomReset()** — сбрасывает масштаб до 100%

```typescript
function ZoomReset(): Promise<void>
```

#### Печать

**Print()** — открывает нативное диалоговое окно печати

```typescript
function Print(): Promise<void>
```

**Пример:**

```javascript
import { Window } from '@wailsio/runtime'

// Open print dialog for current window
await Window.Print()
```

**Примечание:** Открывает нативное диалоговое окно печати ОС, в котором пользователь может выбрать параметры принтера и распечатать текущее содержимое окна. В отличие от `window.print()`, который может не работать в WebView, этот метод использует нативный API печати платформы.

### Буфер обмена

Операции с буфером обмена.

#### SetText()

Задаёт текст в буфере обмена.

```typescript
function SetText(text: string): Promise<void>
```

**Пример:**

```javascript
import { Clipboard } from '@wailsio/runtime'

await Clipboard.SetText('Hello from Wails!')
```

#### Text()

Возвращает текст из буфера обмена.

```typescript
function Text(): Promise<string>
```

**Пример:**

```javascript
import { Clipboard } from '@wailsio/runtime'

const clipboardText = await Clipboard.Text()
console.log('Clipboard:', clipboardText)
```

### Система

Низкоуровневые системные методы для прямого взаимодействия с бэкендом.

#### invoke()

Отправляет необработанное сообщение непосредственно бэкенду. При этом стандартная система привязок обходится, а сообщение обрабатывается обработчиком `RawMessageHandler`, заданным в параметрах приложения.

```typescript
function invoke(message: any): void
```

**Пример:**

```javascript
import { System } from '@wailsio/runtime'

// Send a raw message to the backend
System.invoke('my-custom-message')

// Send structured data as JSON
System.invoke(JSON.stringify({ action: 'update', value: 42 }))
```

@note{type="caution"}
Эта функция работает по принципу «отправил и забыл» и ничего не возвращает. Для получения ответов от бэкенда используйте события.

@end

Подробнее см. в [руководстве по необработанным сообщениям](/guides/raw-messages/).

### Приложение

Методы уровня приложения.

#### Show()

Показывает все окна приложения.

```typescript
function Show(): Promise<void>
```

#### Hide()

Скрывает все окна приложения.

```typescript
function Hide(): Promise<void>
```

#### Quit()

Завершает работу приложения.

```typescript
function Quit(): Promise<void>
```

**Пример:**

```javascript
import { Application } from '@wailsio/runtime'

// Add quit button
document.getElementById('quit-btn').addEventListener('click', async () => {
    await Application.Quit()
})
```

### Браузер

Открытие URL-адресов в браузере по умолчанию.

#### OpenURL()

Открывает URL-адрес в системном браузере.

```typescript
function OpenURL(url: string | URL): Promise<void>
```

**Пример:**

```javascript
import { Browser } from '@wailsio/runtime'

await Browser.OpenURL('https://wails.io')
```

### Экраны

Получение сведений об экранах и управление ими.

#### GetAll()

Возвращает все экраны.

```typescript
function GetAll(): Promise<Screen[]>
```

#### GetPrimary()

Возвращает основной экран.

```typescript
function GetPrimary(): Promise<Screen>
```

#### GetCurrent()

Возвращает текущий активный экран.

```typescript
function GetCurrent(): Promise<Screen>
```

**Интерфейс Screen:**

```typescript
interface Screen {
    ID: string
    Name: string
    ScaleFactor: number
    X: number
    Y: number
    Size: { Width: number, Height: number }
    Bounds: { X: number, Y: number, Width: number, Height: number }
    WorkArea: { X: number, Y: number, Width: number, Height: number }
    IsPrimary: boolean
    Rotation: number
}
```

**Пример:**

```javascript
import { Screens } from '@wailsio/runtime'

// List all screens
const screens = await Screens.GetAll()
screens.forEach(screen => {
    console.log(`${screen.Name}: ${screen.Size.Width}x${screen.Size.Height}`)
})

// Get primary screen
const primary = await Screens.GetPrimary()
console.log('Primary screen:', primary.Name)
```

### Диалоговые окна

Нативные диалоговые окна ОС, вызываемые из JavaScript.

#### Info()

Показывает информационное диалоговое окно.

```typescript
function Info(options: MessageDialogOptions): Promise<string>
```

**Пример:**

```javascript
import { Dialogs } from '@wailsio/runtime'

await Dialogs.Info({
    Title: 'Success',
    Message: 'Operation completed successfully!'
})
```

#### Error()

Показывает диалоговое окно с сообщением об ошибке.

```typescript
function Error(options: MessageDialogOptions): Promise<string>
```

#### Warning()

Показывает диалоговое окно с предупреждением.

```typescript
function Warning(options: MessageDialogOptions): Promise<string>
```

#### Question()

Показывает диалоговое окно с вопросом и настраиваемыми кнопками.

```typescript
function Question(options: MessageDialogOptions): Promise<string>
```

**Пример:**

```javascript
import { Dialogs } from '@wailsio/runtime'

const result = await Dialogs.Question({
    Title: 'Confirm Delete',
    Message: 'Are you sure you want to delete this file?',
    Buttons: [
        { Label: 'Delete', IsDefault: false },
        { Label: 'Cancel', IsDefault: true }
    ]
})

if (result === 'Delete') {
    // Delete the file
}
```

#### OpenFile()

Показывает диалоговое окно открытия файла.

```typescript
function OpenFile(options: OpenFileDialogOptions): Promise<string | string[]>
```

**Пример:**

```javascript
import { Dialogs } from '@wailsio/runtime'

const file = await Dialogs.OpenFile({
    Title: 'Select Image',
    Filters: [
        { DisplayName: 'Images', Pattern: '*.png;*.jpg;*.jpeg' },
        { DisplayName: 'All Files', Pattern: '*.*' }
    ]
})

if (file) {
    console.log('Selected:', file)
}
```

#### SaveFile()

Показывает диалоговое окно сохранения файла.

```typescript
function SaveFile(options: SaveFileDialogOptions): Promise<string>
```

### WML (язык разметки Wails)

WML предоставляет декларативные атрибуты для распространённых действий. Добавьте атрибуты к HTML-элементам:

#### Атрибуты

**wml-event** — генерирует событие при щелчке

```html
<button wml-event="save-clicked">Save</button>
```

**wml-window** — вызывает метод окна

```html
<button wml-window="Close">Close Window</button>
<button wml-window="Minimise">Minimize</button>
```

**wml-target-window** — задаёт целевое окно для wml-window

```html
<button wml-window="Show" wml-target-window="settings">
    Show Settings
</button>
```

**wml-openurl** — открывает URL-адрес в браузере

```html
<a href="#" wml-openurl="https://wails.io">Visit Wails</a>
```

**wml-confirm** — перед выполнением действия показывает диалоговое окно подтверждения

```html
<button wml-window="Close" wml-confirm="Are you sure you want to close?">
    Close
</button>
```

**Пример:**

```html
<div>
    <button wml-event="save-clicked">Save</button>
    <button wml-window="Minimise">Minimize</button>
    <button wml-window="Close" wml-confirm="Close window?">Close</button>
    <a href="#" wml-openurl="https://github.com/wailsapp/wails">GitHub</a>
</div>
```

## Полный пример

```javascript
import { Events, Window, Clipboard, Dialogs, Screens } from '@wailsio/runtime'

// Listen for events from Go
Events.On('data-updated', (event) => {
    console.log('Data:', event.data)
    updateUI(event.data)
})

// Window management
document.getElementById('center-btn').addEventListener('click', async () => {
    await Window.Center()
})

document.getElementById('fullscreen-btn').addEventListener('click', async () => {
    const isFullscreen = await Window.IsFullscreen()
    if (isFullscreen) {
        await Window.UnFullscreen()
    } else {
        await Window.Fullscreen()
    }
})

// Clipboard operations
document.getElementById('copy-btn').addEventListener('click', async () => {
    await Clipboard.SetText('Copied from Wails!')
})

// Dialog with confirmation
document.getElementById('delete-btn').addEventListener('click', async () => {
    const result = await Dialogs.Question({
        Title: 'Confirm',
        Message: 'Delete this item?',
        Buttons: [
            { Label: 'Delete' },
            { Label: 'Cancel', IsDefault: true }
        ]
    })

    if (result === 'Delete') {
        await Events.Emit('delete-item', { id: currentItemId })
    }
})

// Screen information
const screens = await Screens.GetAll()
console.log(`Detected ${screens.length} screen(s)`)
screens.forEach(screen => {
    console.log(`- ${screen.Name}: ${screen.Size.Width}x${screen.Size.Height}`)
})
```

## Рекомендации

### ✅ Делайте так

- **Импортируйте выборочно** — импортируйте только необходимое
- **Обрабатывайте промисы** — все методы возвращают промисы
- **Используйте WML для простых действий** — это лаконичнее, чем JavaScript
- **Проверяйте возвращаемые значения** — особенно при работе с диалоговыми окнами
- **Отписывайтесь от событий** — удаляйте подписки, когда они больше не нужны

### ❌ Не делайте так

- **Не забывайте об await** — большинство методов асинхронны
- **Не блокируйте пользовательский интерфейс** — правильно используйте async/await
- **Не игнорируйте ошибки** — всегда обрабатывайте отклонённые промисы

## Поддержка TypeScript

Среда выполнения включает полные определения типов TypeScript:

```typescript
import { Events, Window } from '@wailsio/runtime'

Events.On('custom-event', (event) => {
    // TypeScript knows event.data, event.name, event.sender
    console.log(event.data)
})

// All methods are fully typed
const size: { width: number, height: number } = await Window.Size()
```
