---
title: "フロントエンドランタイム"
description: "フロントエンド統合用の Wails JavaScript ランタイムパッケージ"
slug: "reference/frontend-runtime"
sourcePath: "reference/frontend-runtime.md"
---

Wails フロントエンドランタイムは、Wails アプリケーションの標準ライブラリです。アプリケーションで使用できる次のような機能を提供します。

- ウィンドウ管理
- ダイアログ
- ブラウザー統合
- クリップボード
- メニュー
- システム情報
- イベント
- コンテキストメニュー
- 画面
- WML（Wails Markup Language）

Go とフロントエンドの統合にはランタイムが必要です。ランタイムの統合方法は 2 つあります。

- `@wailsio/runtime` パッケージを使用する
- ビルド済みバンドルを使用する

## npm パッケージの使用

`@wailsio/runtime` パッケージは、フロントエンドから Wails ランタイムにアクセスできる JavaScript パッケージです。すべての標準テンプレートで使用されており、アプリケーションにランタイムを統合する方法として推奨されています。`@wailsio/runtime` パッケージを使用すると、ランタイムのうち実際に使用する部分だけが含まれます。

このパッケージは npm で提供されており、次のコマンドでインストールできます。

```shell
npm install --save @wailsio/runtime
```

## ビルド済みバンドルの使用

JavaScript バンドラーを使用しないプロジェクトでは、ビルド済みのバンドル版ランタイムを使用する方が適している場合があります。このバージョンは、次のコマンドを使用してローカルで生成できます。

```shell
wails3 generate runtime
```

このコマンドを実行すると、現在のディレクトリに `runtime.js` ファイル（および `runtime.debug.js` ファイル）が出力されます。このファイルは ES モジュールであり、npm パッケージと同様にアプリケーションのスクリプトからインポートできます。また、API はグローバルな window オブジェクトにもエクスポートされるため、単純なアプリケーションでは次のように使用できます。

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
ランタイムを読み込む `<script>` タグには `type="module"` 属性を指定し、API を呼び出す前にページが完全に読み込まれるまで待つことが重要です。これは、`type="module"` 属性を持つスクリプトが非同期で実行されるためです。

@end

## 初期化

ランタイムは API 関数に加えて、コンテキストメニューとウィンドウのドラッグをサポートします。これらの機能は、ランタイムを初期化した後にのみ期待どおりに動作します。API を使用しない場合でも、フロントエンドコード内のいずれかの場所に、副作用のための import 文を必ず含めてください。

```javascript
import "@wailsio/runtime";
```

バンドラーは副作用の存在を検出し、必要な初期化コードをすべてビルドに含めるはずです。

@note{type="info"}
ビルド済みバンドルを使用する場合は、上記のように script タグを追加するだけで十分です。

@end

## 型付きイベント用 Vite プラグイン

ランタイムには、開発中の型付きイベントで HMR（Hot Module Replacement）を利用できるようにする Vite プラグインが含まれています。

### セットアップ

`vite.config.ts` にプラグインを追加します。

```typescript
import { defineConfig } from 'vite'
import wails from '@wailsio/runtime/plugins/vite'

export default defineConfig({
  plugins: [wails()],
})
```

### 利点

- **自動再読み込み**：`wails3 generate bindings` を実行すると、イベントバインディングが自動的に再生成され、再読み込みされます
- **開発モード**：`wails3 dev` とシームレスに連携し、変更が即座に反映されます
- **型安全性**：自動補完と型チェックを備えた完全な TypeScript サポートを提供します

### イベント登録での使用

Go でイベントを登録します。

```go
type UserData struct {
    ID   string
    Name string
}

func init() {
    application.RegisterEvent[UserData]("user-updated")
}
```

バインディングを生成します。

```bash
wails3 generate bindings

# Or, to include TypeScript definitions
wails3 generate bindings -ts

# For more options, see:
wails3 generate bindings -help
```

フロントエンドで型付きイベントを使用します。

```typescript
import { Events } from '@wailsio/runtime'
import { UserUpdated } from './bindings/events'

// Type-safe event with autocomplete
Events.Emit(UserUpdated({
    ID: "123",
    Name: "John Doe"
}))
```

## API リファレンス

ランタイムはモジュール単位で構成されており、各モジュールが特定の機能を提供します。必要なものだけをインポートしてください。

```javascript
import { Events, Window, Clipboard } from '@wailsio/runtime'
```

### イベント

Go と JavaScript の間で通信するためのイベントシステムです。

#### On()

イベントのコールバックを登録します。

```typescript
function On(eventName: string, callback: (event: WailsEvent) => void): () => void
function On<T>(eventType: EventType<T>, callback: (event: WailsEvent<T>) => void): () => void
```

<strong>戻り値：</strong>購読解除関数

**例：**

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

一度だけ実行されるコールバックを登録します。

```typescript
function Once(eventName: string, callback: (event: WailsEvent) => void): () => void
function Once<T>(eventType: EventType<T>, callback: (event: WailsEvent<T>) => void): () => void
```

**例：**

```javascript
import { Events } from '@wailsio/runtime'

Events.Once('app-ready', () => {
    console.log('App initialized')
})
```

#### Emit()

Go バックエンドまたは他のウィンドウにイベントを送出します。

```typescript
function Emit(name: string, data?: any): Promise<boolean>
function Emit<T>(event: Event<T>): Promise<boolean>
```

<strong>戻り値：</strong>イベントがキャンセルされた場合は `true`、それ以外の場合は `false` に解決される Promise

**例：**

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
戻り値は、イベントがフックによってキャンセルされたかどうかを示します。ほとんどのイベントはキャンセルできず、常に `false` を返します。

@end

#### Off()

イベントリスナーを削除します。

```typescript
function Off(...eventNames: string[]): void
```

**例：**

```javascript
import { Events } from '@wailsio/runtime'

Events.Off('user-logged-in', 'user-logged-out')
```

#### OffAll()

すべてのイベントリスナーを削除します。

```typescript
function OffAll(): void
```

### ウィンドウ

ウィンドウ管理用のメソッドです。デフォルトエクスポートは現在のウィンドウです。

```javascript
import { Window } from '@wailsio/runtime'

// Current window
await Window.SetTitle('New Title')
await Window.Center()

// Get another window
const otherWindow = Window.Get('secondary')
await otherWindow.Show()
```

#### 表示と非表示

**Show()** — ウィンドウを表示します

```typescript
function Show(): Promise<void>
```

**Hide()** — ウィンドウを非表示にします

```typescript
function Hide(): Promise<void>
```

**Close()** - ウィンドウを閉じます

```typescript
function Close(): Promise<void>
```

#### サイズと位置

**SetSize(width, height)** - ウィンドウのサイズを設定します

```typescript
function SetSize(width: number, height: number): Promise<void>
```

**Size()** - ウィンドウのサイズを取得します

```typescript
function Size(): Promise<{ width: number, height: number }>
```

**SetPosition(x, y)** - 絶対位置を設定します

```typescript
function SetPosition(x: number, y: number): Promise<void>
```

**Position()** - 絶対位置を取得します

```typescript
function Position(): Promise<{ x: number, y: number }>
```

**Center()** - ウィンドウを中央に配置します

```typescript
function Center(): Promise<void>
```

**例：**

```javascript
import { Window } from '@wailsio/runtime'

// Resize and center
await Window.SetSize(800, 600)
await Window.Center()

// Get current size
const { width, height } = await Window.Size()
```

#### ウィンドウの状態

**Minimise()** - ウィンドウを最小化します

```typescript
function Minimise(): Promise<void>
```

**Maximise()** - ウィンドウを最大化します

```typescript
function Maximise(): Promise<void>
```

**Fullscreen()** - 全画面表示に切り替えます

```typescript
function Fullscreen(): Promise<void>
```

**Restore()** - 最小化、最大化、または全画面表示から元の状態に戻します

```typescript
function Restore(): Promise<void>
```

**IsMinimised()** - 最小化されているか確認します

```typescript
function IsMinimised(): Promise<boolean>
```

**IsMaximised()** - 最大化されているか確認します

```typescript
function IsMaximised(): Promise<boolean>
```

**IsFullscreen()** - 全画面表示になっているか確認します

```typescript
function IsFullscreen(): Promise<boolean>
```

#### ウィンドウのプロパティ

**SetTitle(title)** - ウィンドウのタイトルを設定します

```typescript
function SetTitle(title: string): Promise<void>
```

**Name()** - ウィンドウ名を取得します

```typescript
function Name(): Promise<string>
```

**SetBackgroundColour(r, g, b, a)** - 背景色を設定します

```typescript
function SetBackgroundColour(r: number, g: number, b: number, a: number): Promise<void>
```

**SetAlwaysOnTop(alwaysOnTop)** - ウィンドウを常に最前面に表示します

```typescript
function SetAlwaysOnTop(alwaysOnTop: boolean): Promise<void>
```

**SetResizable(resizable)** - ウィンドウをサイズ変更可能にします

```typescript
function SetResizable(resizable: boolean): Promise<void>
```

#### フォーカスと画面

**Focus()** - ウィンドウにフォーカスを移します

```typescript
function Focus(): Promise<void>
```

**IsFocused()** - フォーカスされているか確認します

```typescript
function IsFocused(): Promise<boolean>
```

**GetScreen()** - ウィンドウが表示されている画面を取得します

```typescript
function GetScreen(): Promise<Screen>
```

#### コンテンツ

**Reload()** - ページを再読み込みします

```typescript
function Reload(): Promise<void>
```

**ForceReload()** - キャッシュを消去してページを強制的に再読み込みします

```typescript
function ForceReload(): Promise<void>
```

#### ズーム

**SetZoom(level)** - ズームレベルを設定します

```typescript
function SetZoom(level: number): Promise<void>
```

**GetZoom()** - ズームレベルを取得します

```typescript
function GetZoom(): Promise<number>
```

**ZoomIn()** - ズーム倍率を上げます

```typescript
function ZoomIn(): Promise<void>
```

**ZoomOut()** - ズーム倍率を下げます

```typescript
function ZoomOut(): Promise<void>
```

**ZoomReset()** - ズームを100%にリセットします

```typescript
function ZoomReset(): Promise<void>
```

#### 印刷

**Print()** - OSネイティブの印刷ダイアログを開きます

```typescript
function Print(): Promise<void>
```

**例：**

```javascript
import { Window } from '@wailsio/runtime'

// Open print dialog for current window
await Window.Print()
```

**注：** OSネイティブの印刷ダイアログが開き、ユーザーはプリンター設定を選択して、現在のウィンドウの内容を印刷できます。WebViewでは動作しない場合がある`window.print()`とは異なり、これは各プラットフォームのネイティブ印刷APIを使用します。

### クリップボード

クリップボード操作。

#### SetText()

クリップボードのテキストを設定します。

```typescript
function SetText(text: string): Promise<void>
```

**例：**

```javascript
import { Clipboard } from '@wailsio/runtime'

await Clipboard.SetText('Hello from Wails!')
```

#### Text()

クリップボードのテキストを取得します。

```typescript
function Text(): Promise<string>
```

**例：**

```javascript
import { Clipboard } from '@wailsio/runtime'

const clipboardText = await Clipboard.Text()
console.log('Clipboard:', clipboardText)
```

### システム

バックエンドと直接通信するための低レベルのシステムメソッド。

#### invoke()

生のメッセージをバックエンドへ直接送信します。標準のバインディングシステムを経由せず、アプリケーションオプションの`RawMessageHandler`によって処理されます。

```typescript
function invoke(message: any): void
```

**例：**

```javascript
import { System } from '@wailsio/runtime'

// Send a raw message to the backend
System.invoke('my-custom-message')

// Send structured data as JSON
System.invoke(JSON.stringify({ action: 'update', value: 42 }))
```

@note{type="caution"}
これは戻り値のない、送信後に完了を待たない関数です。バックエンドから応答を受信するにはイベントを使用します。

@end

詳細については、[生メッセージガイド](/guides/raw-messages/)を参照してください。

### アプリケーション

アプリケーションレベルのメソッド。

#### Show()

アプリケーションのすべてのウィンドウを表示します。

```typescript
function Show(): Promise<void>
```

#### Hide()

アプリケーションのすべてのウィンドウを非表示にします。

```typescript
function Hide(): Promise<void>
```

#### Quit()

アプリケーションを終了します。

```typescript
function Quit(): Promise<void>
```

**例：**

```javascript
import { Application } from '@wailsio/runtime'

// Add quit button
document.getElementById('quit-btn').addEventListener('click', async () => {
    await Application.Quit()
})
```

### ブラウザー

デフォルトのブラウザーで URL を開きます。

#### OpenURL()

システムのブラウザーで URL を開きます。

```typescript
function OpenURL(url: string | URL): Promise<void>
```

**例：**

```javascript
import { Browser } from '@wailsio/runtime'

await Browser.OpenURL('https://wails.io')
```

### 画面

画面情報の取得と管理を行います。

#### GetAll()

すべての画面を取得します。

```typescript
function GetAll(): Promise<Screen[]>
```

#### GetPrimary()

プライマリー画面を取得します。

```typescript
function GetPrimary(): Promise<Screen>
```

#### GetCurrent()

現在アクティブな画面を取得します。

```typescript
function GetCurrent(): Promise<Screen>
```

**Screen インターフェース：**

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

**例：**

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

### ダイアログ

JavaScript から OS ネイティブのダイアログを使用します。

#### Info()

情報ダイアログを表示します。

```typescript
function Info(options: MessageDialogOptions): Promise<string>
```

**例：**

```javascript
import { Dialogs } from '@wailsio/runtime'

await Dialogs.Info({
    Title: 'Success',
    Message: 'Operation completed successfully!'
})
```

#### Error()

エラーダイアログを表示します。

```typescript
function Error(options: MessageDialogOptions): Promise<string>
```

#### Warning()

警告ダイアログを表示します。

```typescript
function Warning(options: MessageDialogOptions): Promise<string>
```

#### Question()

カスタムボタンを備えた質問ダイアログを表示します。

```typescript
function Question(options: MessageDialogOptions): Promise<string>
```

**例：**

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

ファイルを開くダイアログを表示します。

```typescript
function OpenFile(options: OpenFileDialogOptions): Promise<string | string[]>
```

**例：**

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

ファイルを保存するダイアログを表示します。

```typescript
function SaveFile(options: SaveFileDialogOptions): Promise<string>
```

### WML（Wails Markup Language）

WML は、一般的なアクションを宣言的に指定するための属性を提供します。HTML 要素に属性を追加します：

#### 属性

**wml-event** - クリック時にイベントを発行します

```html
<button wml-event="save-clicked">Save</button>
```

**wml-window** - ウィンドウのメソッドを呼び出します

```html
<button wml-window="Close">Close Window</button>
<button wml-window="Minimise">Minimize</button>
```

**wml-target-window** - wml-window の対象ウィンドウを指定します

```html
<button wml-window="Show" wml-target-window="settings">
    Show Settings
</button>
```

**wml-openurl** - ブラウザーで URL を開きます

```html
<a href="#" wml-openurl="https://wails.io">Visit Wails</a>
```

**wml-confirm** - アクションを実行する前に確認ダイアログを表示します

```html
<button wml-window="Close" wml-confirm="Are you sure you want to close?">
    Close
</button>
```

**例：**

```html
<div>
    <button wml-event="save-clicked">Save</button>
    <button wml-window="Minimise">Minimize</button>
    <button wml-window="Close" wml-confirm="Close window?">Close</button>
    <a href="#" wml-openurl="https://github.com/wailsapp/wails">GitHub</a>
</div>
```

## 完全な例

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

## ベストプラクティス

### ✅ 推奨事項

- **選択的にインポートする** - 必要なものだけをインポートします
- **Promise を処理する** - すべてのメソッドは Promise を返します
- **単純なアクションには WML を使用する** - JavaScript より簡潔に記述できます
- **戻り値を確認する** - 特にダイアログでは重要です
- **イベントの購読を解除する** - 完了したらクリーンアップします

### ❌ 禁止事項

- **await を忘れない** - ほとんどのメソッドは非同期です
- **UI をブロックしない** - async/await を適切に使用します
- **エラーを無視しない** - Promise の拒否を必ず処理します

## TypeScript のサポート

ランタイムには完全な TypeScript 型定義が含まれています：

```typescript
import { Events, Window } from '@wailsio/runtime'

Events.On('custom-event', (event) => {
    // TypeScript knows event.data, event.name, event.sender
    console.log(event.data)
})

// All methods are fully typed
const size: { width: number, height: number } = await Window.Size()
```
