---
title: "프런트엔드 런타임"
description: "프런트엔드 통합을 위한 Wails JavaScript 런타임 패키지"
slug: "reference/frontend-runtime"
sourcePath: "reference/frontend-runtime.md"
---

Wails 프런트엔드 런타임은 Wails 애플리케이션용 표준 라이브러리입니다. 애플리케이션에서 사용할 수 있는 다음과 같은 여러 기능을 제공합니다:

- 창 관리
- 대화 상자
- 브라우저 통합
- 클립보드
- 메뉴
- 시스템 정보
- 이벤트
- 상황에 맞는 메뉴
- 화면
- WML(Wails Markup Language)

Go와 프런트엔드를 통합하려면 런타임이 필요합니다. 런타임을 통합하는 방법은 2가지입니다:

- `@wailsio/runtime` 패키지 사용
- 미리 빌드된 번들 사용

## npm 패키지 사용

`@wailsio/runtime` 패키지는 프런트엔드에서 Wails 런타임에 접근할 수 있게 해 주는 JavaScript 패키지입니다. 모든 표준 템플릿에서 사용되며, 애플리케이션에 런타임을 통합할 때 권장되는 방법입니다. `@wailsio/runtime` 패키지를 사용하면 런타임에서 실제로 사용하는 부분만 포함됩니다.

이 패키지는 npm에서 제공되며 다음 명령으로 설치할 수 있습니다:

```shell
npm install --save @wailsio/runtime
```

## 미리 빌드된 번들 사용

일부 프로젝트에서는 JavaScript 번들러를 사용하지 않으므로 미리 빌드된 런타임 번들을 사용하는 편이 나을 수 있습니다. 다음 명령을 사용하여 이 버전을 로컬에서 생성할 수 있습니다:

```shell
wails3 generate runtime
```

이 명령은 현재 디렉터리에 `runtime.js` 파일과 `runtime.debug.js` 파일을 출력합니다. 이 파일은 npm 패키지와 마찬가지로 애플리케이션 스크립트에서 가져올 수 있는 ES 모듈입니다. API도 전역 window 객체로 내보내므로 간단한 애플리케이션에서는 다음과 같이 사용할 수 있습니다:

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
런타임을 로드하는 `<script>` 태그에 `type="module"` 특성을 반드시 포함하고, API를 호출하기 전에 페이지가 완전히 로드될 때까지 기다려야 합니다. `type="module"` 특성이 있는 스크립트는 비동기식으로 실행되기 때문입니다.

@end

## 초기화

런타임은 API 함수 외에도 상황에 맞는 메뉴와 창 끌기를 지원합니다. 이러한 기능은 런타임이 초기화된 후에만 정상적으로 작동합니다. API를 사용하지 않더라도 프런트엔드 코드 어딘가에 부수 효과를 위한 import 문을 반드시 포함하세요:

```javascript
import "@wailsio/runtime";
```

번들러는 부수 효과의 존재를 감지하여 필요한 모든 초기화 코드를 빌드에 포함할 것입니다.

@note{type="info"}
미리 빌드된 번들을 사용하려면 위와 같이 script 태그를 추가하는 것으로 충분합니다.

@end

## 타입 지정 이벤트용 Vite 플러그인

런타임에는 개발 중 타입 지정 이벤트에 대한 HMR(Hot Module Replacement)을 지원하는 Vite 플러그인이 포함되어 있습니다.

### 설정

`vite.config.ts`에 플러그인을 추가하세요:

```typescript
import { defineConfig } from 'vite'
import wails from '@wailsio/runtime/plugins/vite'

export default defineConfig({
  plugins: [wails()],
})
```

### 장점

- **자동 다시 로드**: `wails3 generate bindings`을 실행하면 이벤트 바인딩이 자동으로 다시 생성되고 로드됩니다.
- **개발 모드**: `wails3 dev`와 원활하게 연동되어 변경 사항이 즉시 반영됩니다.
- **타입 안전성**: 자동 완성과 타입 검사를 포함한 완전한 TypeScript 지원

### 이벤트 등록과 함께 사용

Go에서 이벤트를 등록하세요:

```go
type UserData struct {
    ID   string
    Name string
}

func init() {
    application.RegisterEvent[UserData]("user-updated")
}
```

바인딩을 생성하세요:

```bash
wails3 generate bindings

# Or, to include TypeScript definitions
wails3 generate bindings -ts

# For more options, see:
wails3 generate bindings -help
```

프런트엔드에서 타입 지정 이벤트를 사용하세요:

```typescript
import { Events } from '@wailsio/runtime'
import { UserUpdated } from './bindings/events'

// Type-safe event with autocomplete
Events.Emit(UserUpdated({
    ID: "123",
    Name: "John Doe"
}))
```

## API 참조

런타임은 각각 특정 기능을 제공하는 모듈로 구성됩니다. 필요한 항목만 가져오세요:

```javascript
import { Events, Window, Clipboard } from '@wailsio/runtime'
```

### 이벤트

Go와 JavaScript 간 통신을 위한 이벤트 시스템입니다.

#### On()

이벤트에 대한 콜백을 등록합니다.

```typescript
function On(eventName: string, callback: (event: WailsEvent) => void): () => void
function On<T>(eventType: EventType<T>, callback: (event: WailsEvent<T>) => void): () => void
```

**반환값:** 구독 해제 함수

**예:**

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

한 번만 실행되는 콜백을 등록합니다.

```typescript
function Once(eventName: string, callback: (event: WailsEvent) => void): () => void
function Once<T>(eventType: EventType<T>, callback: (event: WailsEvent<T>) => void): () => void
```

**예:**

```javascript
import { Events } from '@wailsio/runtime'

Events.Once('app-ready', () => {
    console.log('App initialized')
})
```

#### Emit()

Go 백엔드 또는 다른 창으로 이벤트를 발생시킵니다.

```typescript
function Emit(name: string, data?: any): Promise<boolean>
function Emit<T>(event: Event<T>): Promise<boolean>
```

**반환값:** 이벤트가 취소되면 `true`로, 그렇지 않으면 `false`로 이행되는 Promise

**예:**

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
반환값은 훅에 의해 이벤트가 취소되었는지를 나타냅니다. 대부분의 이벤트는 취소할 수 없으며 항상 `false`을 반환합니다.

@end

#### Off()

이벤트 리스너를 제거합니다.

```typescript
function Off(...eventNames: string[]): void
```

**예:**

```javascript
import { Events } from '@wailsio/runtime'

Events.Off('user-logged-in', 'user-logged-out')
```

#### OffAll()

모든 이벤트 리스너를 제거합니다.

```typescript
function OffAll(): void
```

### 창

창 관리 메서드입니다. 기본 내보내기는 현재 창입니다.

```javascript
import { Window } from '@wailsio/runtime'

// Current window
await Window.SetTitle('New Title')
await Window.Center()

// Get another window
const otherWindow = Window.Get('secondary')
await otherWindow.Show()
```

#### 표시 여부

**Show()** - 창을 표시합니다.

```typescript
function Show(): Promise<void>
```

**Hide()** - 창을 숨깁니다.

```typescript
function Hide(): Promise<void>
```

**Close()** - 창을 닫습니다

```typescript
function Close(): Promise<void>
```

#### 크기 및 위치

**SetSize(width, height)** - 창 크기를 설정합니다

```typescript
function SetSize(width: number, height: number): Promise<void>
```

**Size()** - 창 크기를 가져옵니다

```typescript
function Size(): Promise<{ width: number, height: number }>
```

**SetPosition(x, y)** - 절대 위치를 설정합니다

```typescript
function SetPosition(x: number, y: number): Promise<void>
```

**Position()** - 절대 위치를 가져옵니다

```typescript
function Position(): Promise<{ x: number, y: number }>
```

**Center()** - 창을 화면 중앙에 배치합니다

```typescript
function Center(): Promise<void>
```

**예:**

```javascript
import { Window } from '@wailsio/runtime'

// Resize and center
await Window.SetSize(800, 600)
await Window.Center()

// Get current size
const { width, height } = await Window.Size()
```

#### 창 상태

**Minimise()** - 창을 최소화합니다

```typescript
function Minimise(): Promise<void>
```

**Maximise()** - 창을 최대화합니다

```typescript
function Maximise(): Promise<void>
```

**Fullscreen()** - 전체 화면 모드로 전환합니다

```typescript
function Fullscreen(): Promise<void>
```

**Restore()** - 최소화, 최대화 또는 전체 화면 상태에서 복원합니다

```typescript
function Restore(): Promise<void>
```

**IsMinimised()** - 최소화 상태인지 확인합니다

```typescript
function IsMinimised(): Promise<boolean>
```

**IsMaximised()** - 최대화 상태인지 확인합니다

```typescript
function IsMaximised(): Promise<boolean>
```

**IsFullscreen()** - 전체 화면 상태인지 확인합니다

```typescript
function IsFullscreen(): Promise<boolean>
```

#### 창 속성

**SetTitle(title)** - 창 제목을 설정합니다

```typescript
function SetTitle(title: string): Promise<void>
```

**Name()** - 창 이름을 가져옵니다

```typescript
function Name(): Promise<string>
```

**SetBackgroundColour(r, g, b, a)** - 배경색을 설정합니다

```typescript
function SetBackgroundColour(r: number, g: number, b: number, a: number): Promise<void>
```

**SetAlwaysOnTop(alwaysOnTop)** - 창을 항상 위에 표시합니다

```typescript
function SetAlwaysOnTop(alwaysOnTop: boolean): Promise<void>
```

**SetResizable(resizable)** - 창 크기를 조절할 수 있도록 설정합니다

```typescript
function SetResizable(resizable: boolean): Promise<void>
```

#### 포커스 및 화면

**Focus()** - 창에 포커스를 설정합니다

```typescript
function Focus(): Promise<void>
```

**IsFocused()** - 포커스 상태인지 확인합니다

```typescript
function IsFocused(): Promise<boolean>
```

**GetScreen()** - 창이 표시된 화면을 가져옵니다

```typescript
function GetScreen(): Promise<Screen>
```

#### 콘텐츠

**Reload()** - 페이지를 다시 로드합니다

```typescript
function Reload(): Promise<void>
```

**ForceReload()** - 페이지를 강제로 다시 로드합니다(캐시 삭제)

```typescript
function ForceReload(): Promise<void>
```

#### 확대/축소

**SetZoom(level)** - 확대/축소 수준을 설정합니다

```typescript
function SetZoom(level: number): Promise<void>
```

**GetZoom()** - 확대/축소 수준을 가져옵니다

```typescript
function GetZoom(): Promise<number>
```

**ZoomIn()** - 화면을 확대합니다

```typescript
function ZoomIn(): Promise<void>
```

**ZoomOut()** - 화면을 축소합니다

```typescript
function ZoomOut(): Promise<void>
```

**ZoomReset()** - 확대/축소 수준을 100%로 재설정합니다

```typescript
function ZoomReset(): Promise<void>
```

#### 인쇄

**Print()** - 네이티브 인쇄 대화 상자를 엽니다

```typescript
function Print(): Promise<void>
```

**예:**

```javascript
import { Window } from '@wailsio/runtime'

// Open print dialog for current window
await Window.Print()
```

**참고:** 네이티브 OS 인쇄 대화 상자를 열어 사용자가 프린터 설정을 선택하고 현재 창의 콘텐츠를 인쇄할 수 있게 합니다. WebView에서 작동하지 않을 수 있는 `window.print()`와 달리, 이 기능은 플랫폼의 네이티브 인쇄 API를 사용합니다.

### 클립보드

클립보드 작업입니다.

#### SetText()

클립보드 텍스트를 설정합니다.

```typescript
function SetText(text: string): Promise<void>
```

**예:**

```javascript
import { Clipboard } from '@wailsio/runtime'

await Clipboard.SetText('Hello from Wails!')
```

#### Text()

클립보드 텍스트를 가져옵니다.

```typescript
function Text(): Promise<string>
```

**예:**

```javascript
import { Clipboard } from '@wailsio/runtime'

const clipboardText = await Clipboard.Text()
console.log('Clipboard:', clipboardText)
```

### 시스템

백엔드와 직접 통신하기 위한 저수준 시스템 메서드입니다.

#### invoke()

원시 메시지를 백엔드로 직접 전송합니다. 이 메시지는 표준 바인딩 시스템을 우회하며 애플리케이션 옵션의 `RawMessageHandler`에서 처리됩니다.

```typescript
function invoke(message: any): void
```

**예:**

```javascript
import { System } from '@wailsio/runtime'

// Send a raw message to the backend
System.invoke('my-custom-message')

// Send structured data as JSON
System.invoke(JSON.stringify({ action: 'update', value: 42 }))
```

@note{type="caution"}
이 함수는 반환 값 없이 호출 후 응답을 기다리지 않습니다. 백엔드에서 응답을 받으려면 이벤트를 사용하세요.

@end

자세한 내용은 [원시 메시지 가이드](/guides/raw-messages/)를 참조하세요.

### 애플리케이션

애플리케이션 수준의 메서드입니다.

#### Show()

애플리케이션의 모든 창을 표시합니다.

```typescript
function Show(): Promise<void>
```

#### Hide()

애플리케이션의 모든 창을 숨깁니다.

```typescript
function Hide(): Promise<void>
```

#### Quit()

애플리케이션을 종료합니다.

```typescript
function Quit(): Promise<void>
```

**예:**

```javascript
import { Application } from '@wailsio/runtime'

// Add quit button
document.getElementById('quit-btn').addEventListener('click', async () => {
    await Application.Quit()
})
```

### 브라우저

기본 브라우저에서 URL을 엽니다.

#### OpenURL()

시스템 브라우저에서 URL을 엽니다.

```typescript
function OpenURL(url: string | URL): Promise<void>
```

**예:**

```javascript
import { Browser } from '@wailsio/runtime'

await Browser.OpenURL('https://wails.io')
```

### 화면

화면 정보 및 관리 기능입니다.

#### GetAll()

모든 화면을 가져옵니다.

```typescript
function GetAll(): Promise<Screen[]>
```

#### GetPrimary()

기본 화면을 가져옵니다.

```typescript
function GetPrimary(): Promise<Screen>
```

#### GetCurrent()

현재 활성 화면을 가져옵니다.

```typescript
function GetCurrent(): Promise<Screen>
```

**Screen 인터페이스:**

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

**예:**

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

### 대화 상자

JavaScript에서 네이티브 OS 대화 상자를 사용합니다.

#### Info()

정보 대화 상자를 표시합니다.

```typescript
function Info(options: MessageDialogOptions): Promise<string>
```

**예:**

```javascript
import { Dialogs } from '@wailsio/runtime'

await Dialogs.Info({
    Title: 'Success',
    Message: 'Operation completed successfully!'
})
```

#### Error()

오류 대화 상자를 표시합니다.

```typescript
function Error(options: MessageDialogOptions): Promise<string>
```

#### Warning()

경고 대화 상자를 표시합니다.

```typescript
function Warning(options: MessageDialogOptions): Promise<string>
```

#### Question()

사용자 지정 버튼이 있는 질문 대화 상자를 표시합니다.

```typescript
function Question(options: MessageDialogOptions): Promise<string>
```

**예:**

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

파일 열기 대화 상자를 표시합니다.

```typescript
function OpenFile(options: OpenFileDialogOptions): Promise<string | string[]>
```

**예:**

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

파일 저장 대화 상자를 표시합니다.

```typescript
function SaveFile(options: SaveFileDialogOptions): Promise<string>
```

### WML(Wails Markup Language)

WML은 일반적인 작업을 위한 선언적 속성을 제공합니다. HTML 요소에 속성을 추가하세요:

#### 속성

**wml-event** - 클릭할 때 이벤트를 발생시킵니다

```html
<button wml-event="save-clicked">Save</button>
```

**wml-window** - 창 메서드를 호출합니다

```html
<button wml-window="Close">Close Window</button>
<button wml-window="Minimise">Minimize</button>
```

**wml-target-window** - wml-window의 대상 창을 지정합니다

```html
<button wml-window="Show" wml-target-window="settings">
    Show Settings
</button>
```

**wml-openurl** - 브라우저에서 URL을 엽니다

```html
<a href="#" wml-openurl="https://wails.io">Visit Wails</a>
```

**wml-confirm** - 작업 전에 확인 대화 상자를 표시합니다

```html
<button wml-window="Close" wml-confirm="Are you sure you want to close?">
    Close
</button>
```

**예:**

```html
<div>
    <button wml-event="save-clicked">Save</button>
    <button wml-window="Minimise">Minimize</button>
    <button wml-window="Close" wml-confirm="Close window?">Close</button>
    <a href="#" wml-openurl="https://github.com/wailsapp/wails">GitHub</a>
</div>
```

## 전체 예제

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

## 모범 사례

### ✅ 권장 사항

- **선택적으로 가져오기** - 필요한 항목만 가져오세요
- **프로미스 처리하기** - 모든 메서드는 프로미스를 반환합니다
- **간단한 작업에는 WML 사용하기** - JavaScript보다 간결합니다
- **반환값 확인하기** - 특히 대화 상자에서 확인하세요
- **이벤트 구독 해제하기** - 작업이 끝나면 정리하세요

### ❌ 금지 사항

- **await를 빠뜨리지 않기** - 대부분의 메서드는 비동기입니다
- **UI를 차단하지 않기** - async/await를 올바르게 사용하세요
- **오류를 무시하지 않기** - 거부는 항상 처리하세요

## TypeScript 지원

런타임에는 완전한 TypeScript 정의가 포함되어 있습니다:

```typescript
import { Events, Window } from '@wailsio/runtime'

Events.On('custom-event', (event) => {
    // TypeScript knows event.data, event.name, event.sender
    console.log(event.data)
})

// All methods are fully typed
const size: { width: number, height: number } = await Window.Size()
```
