---
title: "런타임 내부 구조"
description: "Wails v3가 부팅되고 실행되며 OS와 통신하는 방식을 자세히 알아봅니다"
slug: "contributing/runtime-internals"
sourcePath: "contributing/runtime-internals.md"
---

<strong>런타임</strong>은 일반 Go 함수를 크로스 플랫폼 데스크톱 애플리케이션으로 변환하는 계층입니다. 이 문서에서는 소스 코드를 추적할 때 접하게 되는 구성 요소를 설명합니다.

---

## 1. 애플리케이션 수명 주기

| 단계 | 코드 경로 | 수행 작업 |
| --- | --- | --- |
| **부트스트랩** | `pkg/application/application.go:init()` | 빌드 시점 데이터를 등록하고 전역 `application` 싱글턴을 생성합니다. |
| **New()** | `application.New(...)` | `Options`의 유효성을 검사하고 <strong>AssetServer</strong>를 시작하며 로깅을 초기화합니다. |
| **Run()** | `application.(*App).Run()` | 1. 플랫폼 `mainthread.X()`를 호출하여 OS UI 스레드에 진입합니다.<br />2. **런타임**(`internal/runtime`)을 부팅합니다.<br />3. 마지막 창이 닫히거나 `Quit()`이 호출될 때까지 블록됩니다. |
| **종료** | `application.(*App).Quit()` | `application:shutdown` 이벤트를 브로드캐스트하고 로그를 플러시한 뒤 창과 서비스를 해제합니다. |

수명 주기는 엄격한 **단일 진입** 방식입니다. 창은 여러 개 생성할 수 있지만 애플리케이션 객체 자체는 한 번만 초기화됩니다.

---

## 2. 창 관리

### 공개 API

```go
win := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "Dashboard",
    Width:  1280,
    Height: 720,
})
win.Show()
```

> `app.Window.New()`은 **인수를 받지 않습니다**. 다음과 같은 경우 `NewWithOptions(...)`을 사용합니다.
>
> `application.WebviewWindowOptions` 구조체를 값으로 전달해야 하는 경우입니다.

`app.Window.New[WithOptions]()`은 플랫폼별 구현이 있는 `pkg/application/webview_window_*.go`에 처리를 위임합니다.

```
pkg/application/
├── webview_window_darwin.go    // WKWebView
├── webview_window_linux.go     // GTK + WebKitGTK (plus linux_cgo*.go)
└── webview_window_windows.go   // WebView2
```

각 파일은 다음 작업을 수행합니다.

1. 네이티브 웹뷰(WKWebView, WebKitGTK, WebView2)를 생성합니다.
2. **메시지 프로세서** 콜백(`pkg/application/messageprocessor*.go`)을 등록합니다.
3. Wails 이벤트(`WindowDidResize`, `WindowFocus`, `WindowFilesDropped`, …)를 `pkg/events`의 상수에 매핑합니다.

`internal/runtime/`은 작은 빌드 태그 연결 코드(`runtime{,_darwin,_linux,_windows,_android,_dev,_prod}.go`)와 `internal/runtime/desktop/` 아래에 포함된 JS 런타임용으로 예약되어 있습니다.

활성 창은 `pkg/application/window_manager.go` / `webview_window.go`에서 추적합니다. `pkg/application/screenmanager.go`은 **디스플레이** 메타데이터(해상도, 배율, 작업 영역)를 위한 것으로, 창을 관리하지 않습니다.

---

## 3. 메시지 처리 파이프라인

JavaScript와 Go 사이의 브리지는 `pkg/application/messageprocessor_*.go`에 있는 **메시지 프로세서** 계열에서 구현합니다.

흐름:

1. <strong>JavaScript</strong>는 `/wails/runtime.js`에서 `Call.ByID(<fnv-id>, ...args)`(`internal/runtime/desktop/@wailsio/runtime/src/calls.ts`에 구현됨)을 호출합니다. 이름 모드 빌드에서는 `Call.ByName("pkg.Struct.Method", ...args)`를 호출합니다.
2. 런타임 도우미는 호출을 패키징한 후 플랫폼별 네이티브 브리지를 통해 Go로 디스패치합니다.
3. <strong>Go</strong>는 `pkg/application/messageprocessor_call.go`에서 메시지를 수신합니다.
4. 프로세서는 `pkg/application/bindings.go`(직접 작성된 `reflect` 기반 구현)에서 바인딩된 메서드를 찾아 호출합니다.
5. 결과 또는 오류는 JS로 다시 마샬링되며, 여기서 `Promise`가 이행되거나 거부됩니다.

> 정확한 JSON 엔벌로프는 JS 측 런타임 도우미와
>
> Go 측 `messageprocessor_call.go`에서 정의됩니다. 이 페이지의 이전 초안에는
>
> `{"t":"c","id":"123","m":"Greet","p":[…]}` 형식이 제시되어 있었지만, 이는
>
> 현재 구현과 일치하지 않습니다. 와이어 형식 버그를 추적할 때는 두 파일을 함께
>
> 확인하십시오.

특수 프로세서:

| 파일 | 용도 |
| --- | --- |
| `messageprocessor_window.go` | 창 작업(숨기기, 최대화, …) |
| `messageprocessor_dialog.go` | 네이티브 대화 상자(`OpenFile`, `MessageBox`, …) |
| `messageprocessor_clipboard.go` | 클립보드 읽기/쓰기 |
| `messageprocessor_events.go` | 이벤트 구독/발행 |
| `messageprocessor_browser.go` | 브라우저 탐색, 개발자 도구 |

프로세서는 **상태 비저장** 방식입니다. 필요한 모든 정보는 각 메시지와 함께 전달되는 `ApplicationContext`에서 가져옵니다.

---

## 4. 이벤트 시스템

이벤트는 네임스페이스가 지정된 문자열이며 세 계층에 걸쳐 디스패치됩니다.

1. **애플리케이션 이벤트**: 전역 수명 주기(`application:ready`, `application:shutdown`).
2. **창 이벤트**: 창별 이벤트(`window:focus`, `window:resize`).
3. **사용자 지정 이벤트**: 사용자가 정의한 이벤트(`chat:new-message`).

구현 세부 정보:

- 이벤트 상수는 `pkg/events/`에 있습니다(`defaults.go`, `known_events.go`, `events.txt`). 이 상수는 `v3/tasks/events/generate.go`에서 생성되며 `events.Common.*`, `events.Mac.*`, `events.Windows.*`, `events.Linux.*`로 노출됩니다. `wails3 generate constants`을 사용해 다시 빌드할 수 있습니다.
- Go 측(애플리케이션 이벤트):
  ```go
  app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {})
  ```

- Go 측(창 이벤트):
  ```go
  window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {})
  ```

- Go 측(사용자 지정 이벤트):
  ```go
  app.Event.On("chat:new-message", func(e *application.CustomEvent) {})
  ```

- JS 측:
  ```js
  import { Events } from "/wails/runtime.js";
  Events.On("chat:new-message", (e) => { /* … */ });
  ```


애플리케이션/창/사용자 지정 이벤트는 모두 `pkg/application/event_manager.go`을 통해 전달됩니다. 창 이벤트 구독의 범위는 해당 창으로 한정되므로 창을 닫으면 그 핸들러가 자동으로 등록 해제됩니다.

---

## 5. 플랫폼별 구현

조건부 컴파일을 사용하면 OS별 차이를 숨기면서 공개 API를 동일하게 유지할 수 있습니다.

| 고려 사항 | Darwin | Linux | Windows |
| --- | --- | --- | --- |
| 메인 스레드 | `mainthread_darwin.go`(Foundation에 Cgo로 연결) | `mainthread_linux.go`(GTK) | `mainthread_windows.go`(Win32 `AttachThreadInput`) |
| 대화 상자 | `dialogs_darwin.*`(NSAlert) | `dialogs_linux.go`(GtkFileChooser) | `dialogs_windows.go`(IFileOpenDialog) |
| 클립보드 | `clipboard_darwin.go` | `clipboard_linux.go` | `clipboard_windows.go` |
| 트레이 아이콘 | `systemtray_darwin.*` | `systemtray_linux.go`(DBus) | `systemtray_windows.go`(Shell_NotifyIcon) |

핵심 원칙:

- <strong>macOS와 Windows</strong>에서는 Cgo를 제한적으로 사용합니다(주로 `pkg/mac/` 및 `pkg/w32`의 `w32` Win32 래퍼를 통해 사용).
- **Linux에서는 필요상 Cgo를 광범위하게 사용합니다**. `pkg/application/linux_cgo.go`(약 69KB)와 `linux_cgo_gtk4.{c,go,h}`(약 50KB 이상)가 GTK/WebKitGTK를 직접 구동합니다.
- OS별 파일을 읽기 쉽게 유지하려면 **빌드 태그**(`//go:build darwin`, `//go:build linux`, …)를 사용하십시오.
- `internal/capabilities/`은 플랫폼별 기능 플래그를 위해 존재하지만, 프레임워크는 `ErrCapability` 센티널을 내보내지 **않습니다**. 기능 게이팅은 플랫폼별 스텁 반환 값으로 처리합니다.

---

## 6. 파일 안내

| 파일 | 수정하는 경우 |
| --- | --- |
| `internal/runtime/runtime_*.go` | 소규모 빌드 태그 스텁 계층(개발 환경과 프로덕션 환경의 차이, OS별 연결 코드)을 변경합니다. |
| `pkg/application/webview_window_*.go` | 새 창 힌트 또는 동작을 구현합니다. |
| `pkg/application/messageprocessor*.go` | JS에서 호출할 수 있는 새 브리지 명령을 추가합니다. |
| `pkg/events/*.go` | 기본 제공 이벤트 정의를 확장합니다(그런 다음 `wails3 generate constants`을 다시 실행합니다). |
| `internal/assetserver/*` | 개발/프로덕션 환경의 애셋 처리를 조정합니다. |
| `internal/runtime/desktop/@wailsio/runtime/src/*` | 내장 JS 런타임(호출/이벤트 디스패치, 대화 상자, 드래그 등)을 수정합니다. |

---

## 7. 디버깅 팁

- `Options.LogLevel`을 구성하고(예: `slog.LevelDebug`) `Options.Logger` 출력을 확인하십시오. `WAILS_LOG_LEVEL` 환경 변수는 없습니다.
- `wails3 dev` 플래그는 `--config`, `--port`, `-s`(HTTPS 활성화)이며, 전역 플래그로 `--no-colour`가 있습니다. `-verbose` 플래그는 없습니다.
- macOS에서는 Objective-C 예외를 조기에 포착할 수 있도록 `lldb --`에서 실행하십시오.
- Windows의 Chromium 문제를 조사하려면 WebView2 디버그 로그를 활성화하십시오. `set WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS=--remote-debugging-port=9222`

---

## 8. 런타임 확장

1. 선택 사항: 새 기능 플래그가 있으면 `internal/capabilities/`에 선언하십시오.
2. 빌드 태그를 사용하여 각 `pkg/application/*_{darwin,linux,windows}.go` 변형에 기능을 구현하십시오. 기능을 지원하지 않는 플랫폼에는 스텁을 제공하십시오.
3. `pkg/application`에 공개 API를 추가하십시오(인터페이스와 구체적인 `WebviewWindow` 메서드, 옵션 구조체 등).
4. JS에서 호출해야 하는 경우 새 메시지 프로세서 메서드(`pkg/application/messageprocessor*.go`)를 등록하고, JS 런타임에 이에 대응하는 헬퍼를 등록하십시오.
5. 이벤트를 추가하는 경우 `pkg/events/`에 해당 상수를 선언하고 `wails3 generate constants`을 실행하여 생성된 파일을 갱신하세요.

이 체크리스트를 따르면 크로스 플랫폼 계약을 온전히 유지할 수 있습니다.

---

## 9. 드래그 앤 드롭

파일 드래그 앤 드롭은 모든 플랫폼에서 <strong>JavaScript 우선 접근 방식</strong>을 사용합니다. 네이티브 계층이 OS 드래그 이벤트를 가로채지만, 실제 드롭 처리와 DOM 상호 작용은 JavaScript에서 이루어집니다.

### 흐름

1. 사용자가 OS에서 Wails 창 위로 파일을 드래그합니다.
2. 네이티브 계층이 드래그를 감지하고 호버 효과를 위해 JavaScript에 알립니다.
3. 사용자가 파일을 드롭합니다.
4. 네이티브 계층이 파일 경로와 좌표를 JavaScript로 전송합니다.
5. JavaScript가 드롭 대상 요소(`data-file-drop-target`)를 찾습니다.
6. JavaScript가 파일 경로와 요소 세부 정보를 Go 백엔드로 전송합니다.
7. Go가 전체 컨텍스트와 함께 `WindowFilesDropped` 이벤트를 발생시킵니다.

### 플랫폼별 구현

| 플랫폼 | 네이티브 계층 | 핵심 과제 |
| --- | --- | --- |
| **Windows** | WebView2의 기본 제공 드래그 지원 | 좌표가 CSS 픽셀 단위이므로 변환할 필요 없음 |
| **macOS** | NSWindow 드래그 델리게이트 | 창 기준 좌표를 웹뷰 기준 좌표로 변환 |
| **Linux** | GTK3 드래그 시그널 | 파일 드래그와 내부 HTML5 드래그를 반드시 구분해야 함 |

### Linux: 드래그 유형 구분

GTK와 WebKit 모두 드래그 이벤트를 처리하려고 합니다. 핵심은 드래그 대상 유형을 확인하는 것입니다.

```c
static gboolean is_file_drag(GdkDragContext *context) {
    GList *targets = gdk_drag_context_list_targets(context);
    for (GList *l = targets; l != NULL; l = l->next) {
        GdkAtom atom = GDK_POINTER_TO_ATOM(l->data);
        gchar *name = gdk_atom_name(atom);
        if (name && g_strcmp0(name, "text/uri-list") == 0) {
            g_free(name);
            return TRUE;  // External file drag
        }
        g_free(name);
    }
    return FALSE;  // Internal HTML5 drag
}
```

시그널 핸들러는 내부 드래그의 경우 WebKit이 처리하도록 `FALSE`을 반환하고, 파일 드래그의 경우 직접 처리하기 위해 `TRUE`을 반환합니다.

### 파일 드롭 차단

`EnableFileDrop`이 `false`인 경우에도 브라우저가 드롭된 파일로 이동하지 못하도록 해야 합니다. 각 플랫폼에서는 이를 서로 다르게 처리합니다.

- **Windows**: JavaScript가 드래그 이벤트에서 `preventDefault()`을 호출합니다.
- **macOS**: JavaScript가 드래그 이벤트에서 `preventDefault()`을 호출합니다.\
- **Linux**: GTK 시그널 핸들러가 네이티브 계층에서 파일 드래그를 가로채 거부합니다.

### 주요 파일

| 파일 | 용도 |
| --- | --- |
| `pkg/application/linux_cgo.go` | GTK 드래그 시그널 핸들러(cgo 프리앰블의 C 코드) |
| `pkg/application/webview_window_darwin.go` | macOS 드래그 델리게이트 |
| `pkg/application/webview_window_windows.go` | WebView2 메시지 처리 |
| `internal/runtime/desktop/@wailsio/runtime/src/window.ts` | JavaScript 드롭 처리 |

### 디버깅

- **Linux**: C 코드에 `printf`을 추가하세요(`fflush(stdout)`을 잊지 마세요).
- **Windows**: `globalApplication.debug()`을 사용하세요.
- **JavaScript**: 브라우저 콘솔을 확인하고 디버그 모드를 활성화하세요.

일반적인 문제:

1. **내부 HTML5 드래그가 작동하지 않음**: 네이티브 핸들러가 이를 가로채고 있습니다(파일이 아닌 드래그에는 `FALSE`을 반환하세요).
2. **호버 효과가 표시되지 않음**: JavaScript 핸들러가 호출되지 않고 있습니다.
3. **좌표가 잘못됨**: 좌표 공간 변환을 확인하세요.

---

이제 런타임 내부 구조를 단계별로 살펴보았습니다. 이 지식을 **코드베이스 레이아웃** 맵 및 **애셋 서버** 문서와 함께 활용하면 구조를 자신 있게 파악하고 의미 있는 기여를 할 수 있습니다. 즐거운 코딩 되세요!
