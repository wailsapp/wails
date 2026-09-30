---
title: "前端运行时"
description: "用于前端集成的 Wails JavaScript 运行时包"
slug: "reference/frontend-runtime"
sourcePath: "reference/frontend-runtime.md"
---

Wails 前端运行时是 Wails 应用程序的标准库。它提供了多项可在应用程序中使用的功能，包括：

- 窗口管理
- 对话框
- 浏览器集成
- 剪贴板
- 菜单
- 系统信息
- 事件
- 上下文菜单
- 屏幕
- WML（Wails 标记语言）

Go 与前端之间的集成需要使用该运行时。可通过 2 种方式集成运行时：

- 使用 `@wailsio/runtime` 包
- 使用预构建的捆绑包

## 使用 npm 包

`@wailsio/runtime` 包是一个 JavaScript 包，可供前端访问 Wails 运行时。所有标准模板都使用该包，并且推荐使用这种方式将运行时集成到应用程序中。使用 `@wailsio/runtime` 包时，只会包含实际使用的运行时部分。

该包已发布到 npm，可使用以下命令安装：

```shell
npm install --save @wailsio/runtime
```

## 使用预构建的捆绑包

有些项目不使用 JavaScript 打包工具，可能更适合使用预先构建的运行时捆绑版本。可使用以下命令在本地生成此版本：

```shell
wails3 generate runtime
```

该命令将在当前目录中输出一个 `runtime.js` 文件（以及一个 `runtime.debug.js` 文件）。该文件是 ES 模块，可以像 npm 包一样由应用程序脚本导入；此外，该 API 还会导出到全局 window 对象，因此较简单的应用程序可按如下方式使用：

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
请务必在加载运行时的 `<script>` 标签上添加 `type="module"` 属性，并等待页面完全加载后再调用 API，因为带有 `type="module"` 属性的脚本会异步运行。

@end

## 初始化

除 API 函数外，运行时还支持上下文菜单和窗口拖动。这些功能只有在运行时初始化后才能按预期工作。即使不使用 API，也务必在前端代码的某处添加一条仅用于触发副作用的导入语句：

```javascript
import "@wailsio/runtime";
```

打包工具应能检测到副作用，并在构建中包含所有必需的初始化代码。

@note{type="info"}
如果更倾向于使用预构建的捆绑包，按上文所示添加 script 标签即可。

@end

## 用于类型化事件的 Vite 插件

运行时包含一个 Vite 插件，可在开发期间为类型化事件启用 HMR（热模块替换）支持。

### 设置

将插件添加到 `vite.config.ts`：

```typescript
import { defineConfig } from 'vite'
import wails from '@wailsio/runtime/plugins/vite'

export default defineConfig({
  plugins: [wails()],
})
```

### 优势

- **自动重新加载**：运行 `wails3 generate bindings` 时，会自动重新生成并重新加载事件绑定
- **开发模式**：与 `wails3 dev` 无缝配合，可即时更新
- **类型安全**：提供完整的 TypeScript 支持，包括自动补全和类型检查

### 与事件注册配合使用

在 Go 中注册事件：

```go
type UserData struct {
    ID   string
    Name string
}

func init() {
    application.RegisterEvent[UserData]("user-updated")
}
```

生成绑定：

```bash
wails3 generate bindings

# Or, to include TypeScript definitions
wails3 generate bindings -ts

# For more options, see:
wails3 generate bindings -help
```

在前端中使用类型化事件：

```typescript
import { Events } from '@wailsio/runtime'
import { UserUpdated } from './bindings/events'

// Type-safe event with autocomplete
Events.Emit(UserUpdated({
    ID: "123",
    Name: "John Doe"
}))
```

## API 参考

运行时按模块组织，每个模块提供特定功能。仅导入所需内容：

```javascript
import { Events, Window, Clipboard } from '@wailsio/runtime'
```

### 事件

用于 Go 与 JavaScript 之间通信的事件系统。

#### On()

为事件注册回调函数。

```typescript
function On(eventName: string, callback: (event: WailsEvent) => void): () => void
function On<T>(eventType: EventType<T>, callback: (event: WailsEvent<T>) => void): () => void
```

<strong>返回值：</strong>取消订阅函数

**示例：**

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

注册仅运行一次的回调函数。

```typescript
function Once(eventName: string, callback: (event: WailsEvent) => void): () => void
function Once<T>(eventType: EventType<T>, callback: (event: WailsEvent<T>) => void): () => void
```

**示例：**

```javascript
import { Events } from '@wailsio/runtime'

Events.Once('app-ready', () => {
    console.log('App initialized')
})
```

#### Emit()

向 Go 后端或其他窗口发出事件。

```typescript
function Emit(name: string, data?: any): Promise<boolean>
function Emit<T>(event: Event<T>): Promise<boolean>
```

<strong>返回值：</strong>一个 Promise；如果事件已取消，则解析为 `true`，否则解析为 `false`

**示例：**

```javascript
import { Events } from '@wailsio/runtime'

// Basic event emission
const wasCancelled = await Events.Emit('button-clicked', { buttonId: 'submit' })

// With typed events (TypeScript)
import { UserLogin } from './bindings/events'

const cancelled = await Events.Emit(UserLogin({
    UserID: "123",
    Username: "john_doe",
    LoginTime: new Date().toISOString()
}))

if (cancelled) {
    console.log('Login was cancelled by a hook')
}
```

@note{type="info"}
返回值表示事件是否已被钩子取消。大多数事件无法取消，并且始终返回 `false`。

@end

#### Off()

移除事件监听器。

```typescript
function Off(...eventNames: string[]): void
```

**示例：**

```javascript
import { Events } from '@wailsio/runtime'

Events.Off('user-logged-in', 'user-logged-out')
```

#### OffAll()

移除所有事件监听器。

```typescript
function OffAll(): void
```

### 窗口

窗口管理方法。默认导出为当前窗口。

```javascript
import { Window } from '@wailsio/runtime'

// Current window
await Window.SetTitle('New Title')
await Window.Center()

// Get another window
const otherWindow = Window.Get('secondary')
await otherWindow.Show()
```

#### 可见性

**Show()** - 显示窗口

```typescript
function Show(): Promise<void>
```

**Hide()** - 隐藏窗口

```typescript
function Hide(): Promise<void>
```

**Close()** - 关闭窗口

```typescript
function Close(): Promise<void>
```

#### 大小和位置

**SetSize(width, height)** - 设置窗口大小

```typescript
function SetSize(width: number, height: number): Promise<void>
```

**Size()** - 获取窗口大小

```typescript
function Size(): Promise<{ width: number, height: number }>
```

**SetPosition(x, y)** - 设置绝对位置

```typescript
function SetPosition(x: number, y: number): Promise<void>
```

**Position()** - 获取绝对位置

```typescript
function Position(): Promise<{ x: number, y: number }>
```

**Center()** - 将窗口居中

```typescript
function Center(): Promise<void>
```

**示例：**

```javascript
import { Window } from '@wailsio/runtime'

// Resize and center
await Window.SetSize(800, 600)
await Window.Center()

// Get current size
const { width, height } = await Window.Size()
```

#### 窗口状态

**Minimise()** - 最小化窗口

```typescript
function Minimise(): Promise<void>
```

**Maximise()** - 最大化窗口

```typescript
function Maximise(): Promise<void>
```

**Fullscreen()** - 进入全屏模式

```typescript
function Fullscreen(): Promise<void>
```

**Restore()** - 从最小化、最大化或全屏状态还原

```typescript
function Restore(): Promise<void>
```

**IsMinimised()** - 检查窗口是否已最小化

```typescript
function IsMinimised(): Promise<boolean>
```

**IsMaximised()** - 检查窗口是否已最大化

```typescript
function IsMaximised(): Promise<boolean>
```

**IsFullscreen()** - 检查窗口是否处于全屏模式

```typescript
function IsFullscreen(): Promise<boolean>
```

#### 窗口属性

**SetTitle(title)** - 设置窗口标题

```typescript
function SetTitle(title: string): Promise<void>
```

**Name()** - 获取窗口名称

```typescript
function Name(): Promise<string>
```

**SetBackgroundColour(r, g, b, a)** - 设置背景颜色

```typescript
function SetBackgroundColour(r: number, g: number, b: number, a: number): Promise<void>
```

**SetAlwaysOnTop(alwaysOnTop)** - 使窗口始终置顶

```typescript
function SetAlwaysOnTop(alwaysOnTop: boolean): Promise<void>
```

**SetResizable(resizable)** - 设置窗口是否可调整大小

```typescript
function SetResizable(resizable: boolean): Promise<void>
```

#### 焦点和屏幕

**Focus()** - 使窗口获得焦点

```typescript
function Focus(): Promise<void>
```

**IsFocused()** - 检查窗口是否已获得焦点

```typescript
function IsFocused(): Promise<boolean>
```

**GetScreen()** - 获取窗口所在的屏幕

```typescript
function GetScreen(): Promise<Screen>
```

#### 内容

**Reload()** - 重新加载页面

```typescript
function Reload(): Promise<void>
```

**ForceReload()** - 强制重新加载页面（清除缓存）

```typescript
function ForceReload(): Promise<void>
```

#### 缩放

**SetZoom(level)** - 设置缩放级别

```typescript
function SetZoom(level: number): Promise<void>
```

**GetZoom()** - 获取缩放级别

```typescript
function GetZoom(): Promise<number>
```

**ZoomIn()** - 放大

```typescript
function ZoomIn(): Promise<void>
```

**ZoomOut()** - 缩小

```typescript
function ZoomOut(): Promise<void>
```

**ZoomReset()** - 将缩放重置为100%

```typescript
function ZoomReset(): Promise<void>
```

#### 打印

**Print()** - 打开原生打印对话框

```typescript
function Print(): Promise<void>
```

**示例：**

```javascript
import { Window } from '@wailsio/runtime'

// Open print dialog for current window
await Window.Print()
```

<strong>注意：</strong>此操作会打开操作系统的原生打印对话框，用户可以选择打印机设置并打印当前窗口内容。与可能无法在 WebView 中工作的`window.print()`不同，此方法使用平台原生打印 API。

### 剪贴板

剪贴板操作。

#### SetText()

设置剪贴板文本。

```typescript
function SetText(text: string): Promise<void>
```

**示例：**

```javascript
import { Clipboard } from '@wailsio/runtime'

await Clipboard.SetText('Hello from Wails!')
```

#### Text()

获取剪贴板文本。

```typescript
function Text(): Promise<string>
```

**示例：**

```javascript
import { Clipboard } from '@wailsio/runtime'

const clipboardText = await Clipboard.Text()
console.log('Clipboard:', clipboardText)
```

### 系统

用于与后端直接通信的底层系统方法。

#### invoke()

将原始消息直接发送到后端。此操作会绕过标准绑定系统，并由应用程序选项中的`RawMessageHandler`处理。

```typescript
function invoke(message: any): void
```

**示例：**

```javascript
import { System } from '@wailsio/runtime'

// Send a raw message to the backend
System.invoke('my-custom-message')

// Send structured data as JSON
System.invoke(JSON.stringify({ action: 'update', value: 42 }))
```

@note{type="caution"}
这是一个即发即弃且没有返回值的函数。请使用事件接收后端的响应。

@end

有关更多详细信息，请参阅[原始消息指南](/guides/raw-messages/)。

### 应用程序

应用程序级方法。

#### Show()

显示所有应用程序窗口。

```typescript
function Show(): Promise<void>
```

#### Hide()

隐藏所有应用程序窗口。

```typescript
function Hide(): Promise<void>
```

#### Quit()

退出应用程序。

```typescript
function Quit(): Promise<void>
```

**示例：**

```javascript
import { Application } from '@wailsio/runtime'

// Add quit button
document.getElementById('quit-btn').addEventListener('click', async () => {
    await Application.Quit()
})
```

### 浏览器

在默认浏览器中打开 URL。

#### OpenURL()

在系统浏览器中打开 URL。

```typescript
function OpenURL(url: string | URL): Promise<void>
```

**示例：**

```javascript
import { Browser } from '@wailsio/runtime'

await Browser.OpenURL('https://wails.io')
```

### 屏幕

屏幕信息和管理。

#### GetAll()

获取所有屏幕。

```typescript
function GetAll(): Promise<Screen[]>
```

#### GetPrimary()

获取主屏幕。

```typescript
function GetPrimary(): Promise<Screen>
```

#### GetCurrent()

获取当前活动屏幕。

```typescript
function GetCurrent(): Promise<Screen>
```

**Screen 接口：**

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

**示例：**

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

### 对话框

通过 JavaScript 显示操作系统原生对话框。

#### Info()

显示信息对话框。

```typescript
function Info(options: MessageDialogOptions): Promise<string>
```

**示例：**

```javascript
import { Dialogs } from '@wailsio/runtime'

await Dialogs.Info({
    Title: 'Success',
    Message: 'Operation completed successfully!'
})
```

#### Error()

显示错误对话框。

```typescript
function Error(options: MessageDialogOptions): Promise<string>
```

#### Warning()

显示警告对话框。

```typescript
function Warning(options: MessageDialogOptions): Promise<string>
```

#### Question()

显示带有自定义按钮的询问对话框。

```typescript
function Question(options: MessageDialogOptions): Promise<string>
```

**示例：**

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

显示文件打开对话框。

```typescript
function OpenFile(options: OpenFileDialogOptions): Promise<string | string[]>
```

**示例：**

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

显示文件保存对话框。

```typescript
function SaveFile(options: SaveFileDialogOptions): Promise<string>
```

### WML（Wails 标记语言）

WML 为常用操作提供声明式属性。请将这些属性添加到 HTML 元素：

#### 属性

**wml-event** - 点击时发出事件

```html
<button wml-event="save-clicked">Save</button>
```

**wml-window** - 调用窗口方法

```html
<button wml-window="Close">Close Window</button>
<button wml-window="Minimise">Minimize</button>
```

**wml-target-window** - 指定 wml-window 的目标窗口

```html
<button wml-window="Show" wml-target-window="settings">
    Show Settings
</button>
```

**wml-openurl** - 在浏览器中打开 URL

```html
<a href="#" wml-openurl="https://wails.io">Visit Wails</a>
```

**wml-confirm** - 执行操作前显示确认对话框

```html
<button wml-window="Close" wml-confirm="Are you sure you want to close?">
    Close
</button>
```

**示例：**

```html
<div>
    <button wml-event="save-clicked">Save</button>
    <button wml-window="Minimise">Minimize</button>
    <button wml-window="Close" wml-confirm="Close window?">Close</button>
    <a href="#" wml-openurl="https://github.com/wailsapp/wails">GitHub</a>
</div>
```

## 完整示例

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

## 最佳实践

### ✅ 应该做

- **按需导入** - 只导入所需内容
- **处理 Promise** - 所有方法都返回 Promise
- **使用 WML 执行简单操作** - 比 JavaScript 更简洁
- **检查返回值** - 对话框尤其如此
- **取消订阅事件** - 使用完毕后进行清理

### ❌ 不应该做

- **不要忘记 await** - 大多数方法都是异步的
- **不要阻塞 UI** - 正确使用 async/await
- **不要忽略错误** - 始终处理 Promise 拒绝

## TypeScript 支持

运行时包含完整的 TypeScript 类型定义：

```typescript
import { Events, Window } from '@wailsio/runtime'

Events.On('custom-event', (event) => {
    // TypeScript knows event.data, event.name, event.sender
    console.log(event.data)
})

// All methods are fully typed
const size: { width: number, height: number } = await Window.Size()
```
