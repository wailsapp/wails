---
title: "前端執行階段"
description: "用於前端整合的 Wails JavaScript 執行階段套件"
slug: "reference/frontend-runtime"
sourcePath: "reference/frontend-runtime.md"
---

Wails 前端執行階段是 Wails 應用程式的標準程式庫，提供多項可在應用程式中使用的功能，包括：

- 視窗管理
- 對話方塊
- 瀏覽器整合
- 剪貼簿
- 選單
- 系統資訊
- 事件
- 內容選單
- 螢幕
- WML（Wails 標記語言）

Go 與前端之間的整合需要使用此執行階段。整合執行階段有2種方式：

- 使用`@wailsio/runtime`套件
- 使用預先建置的套件組合

## 使用 npm 套件

`@wailsio/runtime`套件是一個 JavaScript 套件，可讓前端存取 Wails 執行階段。所有標準範本都使用此套件，這也是將執行階段整合至應用程式的建議方式。使用`@wailsio/runtime`套件時，只會納入您實際使用的執行階段部分。

此套件已在 npm 上提供，可使用以下方式安裝：

```shell
npm install --save @wailsio/runtime
```

## 使用預先建置的套件組合

有些專案不使用 JavaScript 打包工具，可能偏好使用預先建置的執行階段套件組合版本。您可以在本機使用以下命令產生此版本：

```shell
wails3 generate runtime
```

此命令會在目前目錄中輸出`runtime.js`檔案（以及`runtime.debug.js`檔案）。這是 ES 模組，應用程式指令碼可像匯入 npm 套件一樣匯入它；此外，其 API 也會匯出至全域 window 物件，因此較簡單的應用程式可以如下使用：

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
請務必在載入執行階段的`<script>`標籤上加入`type="module"`屬性，並等到頁面完全載入後再呼叫 API，因為具有`type="module"`屬性的指令碼會以非同步方式執行。

@end

## 初始化

除了 API 函式外，執行階段也支援內容選單和視窗拖曳。只有在執行階段完成初始化後，這些功能才會如預期運作。即使您不使用 API，也請務必在前端程式碼中的某處加入僅為觸發副作用的匯入陳述式：

```javascript
import "@wailsio/runtime";
```

打包工具應會偵測到副作用，並在組建中納入所有必要的初始化程式碼。

@note{type="info"}
如果您偏好使用預先建置的套件組合，只需依照上述方式加入 script 標籤即可。

@end

## 用於型別化事件的 Vite 外掛程式

執行階段內含 Vite 外掛程式，可在開發期間為型別化事件啟用 HMR（熱模組替換）支援。

### 設定

將此外掛程式加入您的`vite.config.ts`：

```typescript
import { defineConfig } from 'vite'
import wails from '@wailsio/runtime/plugins/vite'

export default defineConfig({
  plugins: [wails()],
})
```

### 優點

- **自動重新載入**：執行`wails3 generate bindings`時，系統會自動重新產生並載入事件繫結
- **開發模式**：可與`wails3 dev`無縫搭配，立即套用更新
- **型別安全**：完整支援 TypeScript，並提供自動完成和型別檢查

### 搭配事件註冊使用

在 Go 中註冊事件：

```go
type UserData struct {
    ID   string
    Name string
}

func init() {
    application.RegisterEvent[UserData]("user-updated")
}
```

產生繫結：

```bash
wails3 generate bindings

# Or, to include TypeScript definitions
wails3 generate bindings -ts

# For more options, see:
wails3 generate bindings -help
```

在前端使用型別化事件：

```typescript
import { Events } from '@wailsio/runtime'
import { UserUpdated } from './bindings/events'

// Type-safe event with autocomplete
Events.Emit(UserUpdated({
    ID: "123",
    Name: "John Doe"
}))
```

## API 參考

執行階段依模組組織，每個模組各自提供特定功能。只匯入您需要的項目：

```javascript
import { Events, Window, Clipboard } from '@wailsio/runtime'
```

### 事件

用於 Go 與 JavaScript 之間通訊的事件系統。

#### On()

為事件註冊回呼函式。

```typescript
function On(eventName: string, callback: (event: WailsEvent) => void): () => void
function On<T>(eventType: EventType<T>, callback: (event: WailsEvent<T>) => void): () => void
```

<strong>傳回：</strong>取消訂閱函式

**範例：**

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

註冊只執行一次的回呼函式。

```typescript
function Once(eventName: string, callback: (event: WailsEvent) => void): () => void
function Once<T>(eventType: EventType<T>, callback: (event: WailsEvent<T>) => void): () => void
```

**範例：**

```javascript
import { Events } from '@wailsio/runtime'

Events.Once('app-ready', () => {
    console.log('App initialized')
})
```

#### Emit()

將事件發送至 Go 後端或其他視窗。

```typescript
function Emit(name: string, data?: any): Promise<boolean>
function Emit<T>(event: Event<T>): Promise<boolean>
```

<strong>傳回：</strong>一個 Promise；若事件遭取消，會解析為`true`，否則解析為`false`

**範例：**

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
傳回值表示事件是否已被掛鉤函式取消。大多數事件無法取消，因此一律會傳回`false`。

@end

#### Off()

移除事件監聽器。

```typescript
function Off(...eventNames: string[]): void
```

**範例：**

```javascript
import { Events } from '@wailsio/runtime'

Events.Off('user-logged-in', 'user-logged-out')
```

#### OffAll()

移除所有事件監聽器。

```typescript
function OffAll(): void
```

### 視窗

視窗管理方法。預設匯出的是目前視窗。

```javascript
import { Window } from '@wailsio/runtime'

// Current window
await Window.SetTitle('New Title')
await Window.Center()

// Get another window
const otherWindow = Window.Get('secondary')
await otherWindow.Show()
```

#### 可見性

**Show()**－顯示視窗

```typescript
function Show(): Promise<void>
```

**Hide()**－隱藏視窗

```typescript
function Hide(): Promise<void>
```

**Close()** - 關閉視窗

```typescript
function Close(): Promise<void>
```

#### 大小與位置

**SetSize(width, height)** - 設定視窗大小

```typescript
function SetSize(width: number, height: number): Promise<void>
```

**Size()** - 取得視窗大小

```typescript
function Size(): Promise<{ width: number, height: number }>
```

**SetPosition(x, y)** - 設定絕對位置

```typescript
function SetPosition(x: number, y: number): Promise<void>
```

**Position()** - 取得絕對位置

```typescript
function Position(): Promise<{ x: number, y: number }>
```

**Center()** - 將視窗置中

```typescript
function Center(): Promise<void>
```

**範例：**

```javascript
import { Window } from '@wailsio/runtime'

// Resize and center
await Window.SetSize(800, 600)
await Window.Center()

// Get current size
const { width, height } = await Window.Size()
```

#### 視窗狀態

**Minimise()** - 將視窗最小化

```typescript
function Minimise(): Promise<void>
```

**Maximise()** - 將視窗最大化

```typescript
function Maximise(): Promise<void>
```

**Fullscreen()** - 進入全螢幕模式

```typescript
function Fullscreen(): Promise<void>
```

**Restore()** - 從最小化、最大化或全螢幕狀態還原

```typescript
function Restore(): Promise<void>
```

**IsMinimised()** - 檢查視窗是否已最小化

```typescript
function IsMinimised(): Promise<boolean>
```

**IsMaximised()** - 檢查視窗是否已最大化

```typescript
function IsMaximised(): Promise<boolean>
```

**IsFullscreen()** - 檢查視窗是否處於全螢幕模式

```typescript
function IsFullscreen(): Promise<boolean>
```

#### 視窗屬性

**SetTitle(title)** - 設定視窗標題

```typescript
function SetTitle(title: string): Promise<void>
```

**Name()** - 取得視窗名稱

```typescript
function Name(): Promise<string>
```

**SetBackgroundColour(r, g, b, a)** - 設定背景顏色

```typescript
function SetBackgroundColour(r: number, g: number, b: number, a: number): Promise<void>
```

**SetAlwaysOnTop(alwaysOnTop)** - 讓視窗保持在最上層

```typescript
function SetAlwaysOnTop(alwaysOnTop: boolean): Promise<void>
```

**SetResizable(resizable)** - 設定視窗可調整大小

```typescript
function SetResizable(resizable: boolean): Promise<void>
```

#### 焦點與螢幕

**Focus()** - 將焦點移至視窗

```typescript
function Focus(): Promise<void>
```

**IsFocused()** - 檢查視窗是否具有焦點

```typescript
function IsFocused(): Promise<boolean>
```

**GetScreen()** - 取得視窗所在的螢幕

```typescript
function GetScreen(): Promise<Screen>
```

#### 內容

**Reload()** - 重新載入頁面

```typescript
function Reload(): Promise<void>
```

**ForceReload()** - 強制重新載入頁面（清除快取）

```typescript
function ForceReload(): Promise<void>
```

#### 縮放

**SetZoom(level)** - 設定縮放層級

```typescript
function SetZoom(level: number): Promise<void>
```

**GetZoom()** - 取得縮放層級

```typescript
function GetZoom(): Promise<number>
```

**ZoomIn()** - 放大

```typescript
function ZoomIn(): Promise<void>
```

**ZoomOut()** - 縮小

```typescript
function ZoomOut(): Promise<void>
```

**ZoomReset()** - 將縮放重設為100%

```typescript
function ZoomReset(): Promise<void>
```

#### 列印

**Print()** - 開啟原生列印對話方塊

```typescript
function Print(): Promise<void>
```

**範例：**

```javascript
import { Window } from '@wailsio/runtime'

// Open print dialog for current window
await Window.Print()
```

<strong>注意：</strong>這會開啟作業系統的原生列印對話方塊，讓使用者選擇印表機設定並列印目前視窗的內容。`window.print()`在 WebView 中可能無法運作，而此方法使用平台的原生列印 API。

### 剪貼簿

剪貼簿操作。

#### SetText()

設定剪貼簿文字。

```typescript
function SetText(text: string): Promise<void>
```

**範例：**

```javascript
import { Clipboard } from '@wailsio/runtime'

await Clipboard.SetText('Hello from Wails!')
```

#### Text()

取得剪貼簿文字。

```typescript
function Text(): Promise<string>
```

**範例：**

```javascript
import { Clipboard } from '@wailsio/runtime'

const clipboardText = await Clipboard.Text()
console.log('Clipboard:', clipboardText)
```

### 系統

用於直接與後端通訊的低階系統方法。

#### invoke()

將原始訊息直接傳送至後端。此操作會略過標準繫結系統，並由應用程式選項中的`RawMessageHandler`處理。

```typescript
function invoke(message: any): void
```

**範例：**

```javascript
import { System } from '@wailsio/runtime'

// Send a raw message to the backend
System.invoke('my-custom-message')

// Send structured data as JSON
System.invoke(JSON.stringify({ action: 'update', value: 42 }))
```

@note{type="caution"}
這是一個不等待回應的函式，沒有傳回值。請使用事件接收後端的回應。

@end

如需更多詳細資訊，請參閱[原始訊息指南](/guides/raw-messages/)。

### 應用程式

應用程式層級的方法。

#### Show()

顯示所有應用程式視窗。

```typescript
function Show(): Promise<void>
```

#### Hide()

隱藏所有應用程式視窗。

```typescript
function Hide(): Promise<void>
```

#### Quit()

退出應用程式。

```typescript
function Quit(): Promise<void>
```

**範例：**

```javascript
import { Application } from '@wailsio/runtime'

// Add quit button
document.getElementById('quit-btn').addEventListener('click', async () => {
    await Application.Quit()
})
```

### 瀏覽器

在預設瀏覽器中開啟 URL。

#### OpenURL()

在系統瀏覽器中開啟 URL。

```typescript
function OpenURL(url: string | URL): Promise<void>
```

**範例：**

```javascript
import { Browser } from '@wailsio/runtime'

await Browser.OpenURL('https://wails.io')
```

### 螢幕

螢幕資訊與管理。

#### GetAll()

取得所有螢幕。

```typescript
function GetAll(): Promise<Screen[]>
```

#### GetPrimary()

取得主要螢幕。

```typescript
function GetPrimary(): Promise<Screen>
```

#### GetCurrent()

取得目前作用中的螢幕。

```typescript
function GetCurrent(): Promise<Screen>
```

**Screen 介面：**

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

**範例：**

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

### 對話方塊

從 JavaScript 顯示作業系統原生對話方塊。

#### Info()

顯示資訊對話方塊。

```typescript
function Info(options: MessageDialogOptions): Promise<string>
```

**範例：**

```javascript
import { Dialogs } from '@wailsio/runtime'

await Dialogs.Info({
    Title: 'Success',
    Message: 'Operation completed successfully!'
})
```

#### Error()

顯示錯誤對話方塊。

```typescript
function Error(options: MessageDialogOptions): Promise<string>
```

#### Warning()

顯示警告對話方塊。

```typescript
function Warning(options: MessageDialogOptions): Promise<string>
```

#### Question()

顯示含自訂按鈕的詢問對話方塊。

```typescript
function Question(options: MessageDialogOptions): Promise<string>
```

**範例：**

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

顯示開啟檔案對話方塊。

```typescript
function OpenFile(options: OpenFileDialogOptions): Promise<string | string[]>
```

**範例：**

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

顯示儲存檔案對話方塊。

```typescript
function SaveFile(options: SaveFileDialogOptions): Promise<string>
```

### WML（Wails 標記語言）

WML 為常見操作提供宣告式屬性。請將屬性新增至 HTML 元素：

#### 屬性

**wml-event** - 點擊時發出事件

```html
<button wml-event="save-clicked">Save</button>
```

**wml-window** - 呼叫視窗方法

```html
<button wml-window="Close">Close Window</button>
<button wml-window="Minimise">Minimize</button>
```

**wml-target-window** - 指定 wml-window 的目標視窗

```html
<button wml-window="Show" wml-target-window="settings">
    Show Settings
</button>
```

**wml-openurl** - 在瀏覽器中開啟 URL

```html
<a href="#" wml-openurl="https://wails.io">Visit Wails</a>
```

**wml-confirm** - 執行操作前顯示確認對話方塊

```html
<button wml-window="Close" wml-confirm="Are you sure you want to close?">
    Close
</button>
```

**範例：**

```html
<div>
    <button wml-event="save-clicked">Save</button>
    <button wml-window="Minimise">Minimize</button>
    <button wml-window="Close" wml-confirm="Close window?">Close</button>
    <a href="#" wml-openurl="https://github.com/wailsapp/wails">GitHub</a>
</div>
```

## 完整範例

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

## 最佳實務

### ✅ 建議做法

- **選擇性匯入** - 僅匯入所需項目
- **處理 Promise** - 所有方法都會傳回 Promise
- **對簡單操作使用 WML** - 比 JavaScript 更簡潔
- **檢查傳回值** - 對話方塊尤其如此
- **取消訂閱事件** - 完成後進行清理

### ❌ 請勿這樣做

- **不要忘記 await** - 大多數方法都是非同步的
- **不要阻塞 UI** - 正確使用 async/await
- **不要忽略錯誤** - 一律處理 Promise 拒絕

## TypeScript 支援

執行階段包含完整的 TypeScript 型別定義：

```typescript
import { Events, Window } from '@wailsio/runtime'

Events.On('custom-event', (event) => {
    // TypeScript knows event.data, event.name, event.sender
    console.log(event.data)
})

// All methods are fully typed
const size: { width: number, height: number } = await Window.Size()
```
