---
title: "이벤트 가이드"
description: "애플리케이션 통신 및 수명 주기 관리를 위해 Wails v3에서 이벤트를 사용하는 실용적인 가이드"
slug: "guides/events-reference"
sourcePath: "guides/events-reference.md"
---

**참고: 이 가이드는 작성 중입니다**

## 이벤트 가이드

이벤트는 Wails 애플리케이션 통신의 핵심입니다. 이벤트를 사용하면 애플리케이션의 여러 부분이 서로 긴밀하게 결합되지 않고도 통신할 수 있습니다. 이 가이드에서는 Wails 애플리케이션에서 이벤트를 효과적으로 사용하는 데 필요한 모든 내용을 설명합니다.

## Wails 이벤트 이해하기

이벤트는 애플리케이션 전체에 브로드캐스트되는 메시지라고 생각하면 됩니다. 애플리케이션의 어느 부분에서든 이러한 메시지를 수신하고 그에 따라 동작할 수 있습니다. 이벤트는 특히 다음과 같은 경우에 유용합니다.

- **창 변경에 대응**: 창이 최소화되거나 최대화되거나 이동될 때 이를 감지합니다.
- **시스템 이벤트 처리**: 테마 변경이나 전원 이벤트에 대응합니다.
- **사용자 정의 애플리케이션 로직**: 데이터 업데이트나 사용자 동작 같은 기능을 위한 자체 이벤트를 만듭니다.
- **컴포넌트 간 통신**: 앱의 여러 부분이 직접적인 의존성 없이 통신하도록 합니다.

## 이벤트 명명 규칙

모든 Wails 이벤트는 출처를 명확히 나타내는 네임스페이스 패턴을 따릅니다.

- `common:` - Windows, macOS 및 Linux에서 작동하는 크로스 플랫폼 이벤트
- `windows:` - Windows 전용 이벤트
- `mac:` - macOS 전용 이벤트\
- `linux:` - Linux 전용 이벤트

예를 들면 다음과 같습니다.

- `common:WindowFocus` - 창이 포커스를 얻음(모든 플랫폼에서 작동)
- `windows:APMSuspend` - 시스템이 절전 모드로 전환 중(Windows 전용)
- `mac:ApplicationDidBecomeActive` - 앱이 활성 상태가 됨(macOS 전용)

## 이벤트 시작하기

### 이벤트 수신하기(프런트엔드)

가장 일반적인 사용 사례는 프런트엔드 코드에서 이벤트를 수신하는 것입니다.

```javascript
import { Events } from '@wailsio/runtime';

// Listen for when the window gains focus
Events.On('common:WindowFocus', () => {
    console.log('Window is now focused!');
    // Maybe refresh some data or resume animations
});

// Listen for theme changes
Events.On('common:ThemeChanged', (event) => {
    console.log('Theme changed:', event.data);
    // Update your app's theme accordingly
});

// Listen for custom events from your Go backend
Events.On('my-app:data-updated', (event) => {
    console.log('Data updated:', event.data);
    // Update your UI with the new data
});
```

### 이벤트 발생시키기(백엔드)

Go 코드에서 프런트엔드가 수신할 수 있는 이벤트를 발생시킬 수 있습니다.

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "time"
)

func (s *Service) UpdateData() {
    // Do some data processing...

    // Notify the frontend
    app := application.Get()
    app.Event.Emit("my-app:data-updated",
        map[string]interface{}{
            "timestamp": time.Now(),
            "count": 42,
        },
	)
}
```

### 이벤트 발생시키기(프런트엔드)

자주 사용되지는 않지만, 프런트엔드에서 Go 코드가 수신할 수 있는 이벤트를 발생시킬 수도 있습니다.

```javascript
import { Events } from '@wailsio/runtime';

// Event without data
Events.Emit('myapp:close-window')

// Event with data
Events.Emit('myapp:disconnect-requested', 'id-123')
```

프런트엔드에서 TypeScript를 사용하고 Go 코드에서 [형식이 지정된 이벤트를 등록](#-----)하면 이벤트 이름 자동 완성 및 검사와 데이터 형식 검사를 사용할 수 있습니다.

### 이벤트 리스너 제거하기

이벤트 리스너가 더 이상 필요하지 않으면 항상 정리하세요.

```javascript
import { Events } from '@wailsio/runtime';

// Store the handler reference
const focusHandler = () => {
    console.log('Window focused');
};

// Add the listener — capture the unsubscribe function it returns
const unsubscribe = Events.On('common:WindowFocus', focusHandler);

// Later, remove this specific listener via the returned unsubscribe
unsubscribe();

// Or remove ALL listeners for one (or more) event names — Events.Off takes only event-name strings
Events.Off('common:WindowFocus');
// Events.Off('common:WindowFocus', 'common:WindowLostFocus'); // variadic
// Events.OffAll(); // remove every listener for every event (no args)
```

## 일반적인 사용 사례

### 1. 창 포커스에 따라 일시 중지/재개하기

많은 애플리케이션은 창이 포커스를 잃으면 특정 작업을 일시 중지해야 합니다.

```javascript
import { Events } from '@wailsio/runtime';

let animationRunning = true;

Events.On('common:WindowLostFocus', () => {
    animationRunning = false;
    pauseBackgroundTasks();
});

Events.On('common:WindowFocus', () => {
    animationRunning = true;
    resumeBackgroundTasks();
});
```

### 2. 테마 변경에 대응하기

앱의 테마를 시스템 테마와 동기화하세요.

```javascript
import { Events } from '@wailsio/runtime';

Events.On('common:ThemeChanged', (event) => {
    const isDarkMode = event.data.isDarkMode;

    if (isDarkMode) {
        document.body.classList.add('dark-theme');
        document.body.classList.remove('light-theme');
    } else {
        document.body.classList.add('light-theme');
        document.body.classList.remove('dark-theme');
    }
});
```

### 3. 파일 드롭 처리하기

앱에서 드래그한 파일을 받을 수 있도록 하세요.

```javascript
import { Events } from '@wailsio/runtime';

Events.On('common:WindowFilesDropped', (event) => {
    const files = event.data.files;

    files.forEach(file => {
        console.log('File dropped:', file);
        // Process the dropped files
        handleFileUpload(file);
    });
});
```

### 4. 창 수명 주기 관리

창 상태 변경에 대응하세요.

```javascript
import { Events } from '@wailsio/runtime';

Events.On('common:WindowClosing', () => {
    // Save user data before closing
    saveApplicationState();

    // You could also prevent closing by returning false
    // from a registered window close handler
});

Events.On('common:WindowMaximise', () => {
    // Adjust UI for maximized view
    adjustLayoutForMaximized();
});

Events.On('common:WindowRestore', () => {
    // Return UI to normal state
    adjustLayoutForNormal();
});
```

### 5. 플랫폼별 기능

필요할 때 플랫폼별 이벤트를 처리하세요.

```javascript
import { Events } from '@wailsio/runtime';

// Windows-specific power management
Events.On('windows:APMSuspend', () => {
    console.log('System is going to sleep');
    saveState();
});

Events.On('windows:APMResumeSuspend', () => {
    console.log('System woke up');
    refreshData();
});

// macOS-specific app lifecycle
Events.On('mac:ApplicationWillTerminate', () => {
    console.log('App is about to quit');
    performCleanup();
});
```

## 사용자 정의 이벤트 만들기

애플리케이션별 요구 사항에 맞는 자체 이벤트를 만들 수 있습니다.

### 백엔드(Go)

```go
// Emit a custom event when data changes

func (s *Service) ProcessUserData(userData UserData) error {
    // Process the data...

    app := application.Get()
    // Notify all listeners
    app.Event.Emit("user:data-processed",
        map[string]interface{}{
            "userId": userData.ID,
            "status": "completed",
            "timestamp": time.Now(),
        },
    )
    return nil
}

// Emit periodic updates
func (s *Service) StartMonitoring() {
    app := application.Get()
    ticker := time.NewTicker(5 * time.Second)
    go func() {
        for range ticker.C {
            stats := s.collectStats()
            app.Event.Emit("monitor:stats-updated", stats)
        }
    }()
}
```

### 프런트엔드(JavaScript)

```javascript
import { Events } from '@wailsio/runtime';

// Listen for your custom events
Events.On('user:data-processed', (event) => {
    const { userId, status, timestamp } = event.data;

    showNotification(`User ${userId} processing ${status}`);
    updateUIWithNewData();
});

Events.On('monitor:stats-updated', (event) => {
    updateDashboard(event.data);
});
```

## 형식 안전성을 갖춘 형식 지정 이벤트

Wails v3는 이벤트 등록 및 자동 바인딩 생성을 통해 TypeScript의 완전한 형식 안전성을 제공하는 형식 지정 이벤트를 지원합니다.

### 사용자 정의 이벤트 등록하기

초기화 시 `application.RegisterEvent`을 호출하여 사용자 정의 이벤트 이름과 해당 데이터 형식을 등록하세요.

```go
package main

import "github.com/wailsapp/wails/v3/pkg/application"

type UserLoginData struct {
    UserID   string
    Username string
    LoginTime string
}

type MonitorStats struct {
    CPUUsage    float64
    MemoryUsage float64
}

func init() {
    // Register events with their data types
    application.RegisterEvent[UserLoginData]("user:login")
    application.RegisterEvent[MonitorStats]("monitor:stats")

    // Register events without data (void events)
    application.RegisterEvent[application.Void]("app:ready")
}
```

@note{type="caution"}
`RegisterEvent`은 초기화 시 호출하도록 설계되었으며, 다음과 같은 경우 패닉이 발생합니다.

- 인수가 유효하지 않은 경우
- 같은 이벤트 이름을 서로 다른 데이터 형식으로 두 번 등록한 경우

@end

@note{type="info"}
데이터 형식이 항상 같다면 동일한 이벤트를 여러 번 등록해도 안전합니다. 여러 패키지 중 어느 패키지가 로드되더라도 이벤트가 등록되도록 보장할 때 유용합니다.

@end

### 이벤트 등록의 이점

등록이 완료되면 `Event.Emit`에 전달된 데이터 인수가 지정된 형식과 일치하는지 검사합니다. 일치하지 않으면 다음과 같이 처리됩니다.

- 오류가 발생하여 로그에 기록됩니다(또는 등록된 오류 처리기에 전달됩니다).
- 문제가 있는 이벤트는 전파되지 않습니다.
- 이를 통해 등록된 이벤트의 데이터 필드에 항상 선언된 형식에 할당 가능한 값만 지정되도록 보장합니다.

### 엄격 모드

개발 중 등록되지 않은 이벤트에 대한 경고를 활성화하려면 `strictevents` 빌드 태그를 사용하세요.

```bash
go build -tags strictevents
```

엄격 모드를 활성화하면 로그가 과도하게 쌓이지 않도록 런타임이 등록되지 않은 이벤트 이름마다 최대 한 번만 경고를 발생시킵니다.

### TypeScript 바인딩 생성

바인딩 생성기는 프런트엔드에서 형식 지정 이벤트를 별도의 변환 없이 지원할 수 있도록 TypeScript 정의와 연결 코드를 출력합니다.

#### 1. Vite 플러그인 설정하기

`vite.config.ts`에서 다음과 같이 설정하세요.

```typescript
import { defineConfig } from 'vite'
import wails from '@wailsio/runtime/plugins/vite'

export default defineConfig({
  plugins: [wails()],
})
```

#### 2. 바인딩 생성하기

바인딩 생성기를 실행하세요.

```bash
wails3 generate bindings
```

그러면 형식이 지정된 이벤트 데이터 정의가 포함된 TypeScript 파일이 프런트엔드 디렉터리에 생성됩니다. 또한 생성기가 검색할 수 있는 컴파일 타임 상수인 등록된 이벤트 이름별 상수를 하나씩 내보내는 `bindings/events` 모듈도 생성됩니다(예: `user:login`은 `UserLogin`이 됩니다).

#### 3. 프런트엔드에서 형식 지정 이벤트 사용하기

```typescript
import { Events } from '@wailsio/runtime'
import { UserLogin, MonitorStats } from './bindings/events'

// Type-safe event emission with autocomplete
Events.Emit(UserLogin, {
    UserID: "123",
    Username: "john_doe",
    LoginTime: new Date().toISOString()
})

// Type-safe event listening
Events.On(UserLogin, (event) => {
    // event.data is typed as UserLoginData
    console.log(`User ${event.data.Username} logged in`)
})

Events.On(MonitorStats, (event) => {
    // event.data is typed as MonitorStats
    updateDashboard({
        cpu: event.data.CPUUsage,
        memory: event.data.MemoryUsage
    })
})
```

형식 지정 이벤트는 다음 기능을 제공합니다.

- 이벤트 이름 **자동 완성**
- 이벤트 데이터 **형식 검사**
- 데이터 형식 불일치에 대한 **컴파일 타임 오류**
- **IntelliSense** 문서

## 이벤트 참조

### 공통 이벤트(크로스 플랫폼)

다음 이벤트는 모든 플랫폼에서 작동합니다.

| 이벤트 | 설명 | 사용 시점 |
| --- | --- | --- |
| `common:ApplicationStarted` | 애플리케이션이 완전히 시작됨 | 앱을 초기화하고 저장된 상태 로드 |
| `common:WindowRuntimeReady` | Wails 런타임 준비 완료 | Wails API 호출 시작 |
| `common:ThemeChanged` | 시스템 테마 변경됨 | 앱 모양 업데이트 |
| `common:SystemWillSleep` | 시스템이 곧 절전 모드로 전환됨 | 상태를 저장하고 소켓 닫기 |
| `common:SystemDidWake` | 시스템이 절전 모드에서 재개됨 | 다시 연결하고 오래된 데이터 새로 고침 |
| `common:WindowFocus` | 창이 포커스를 얻음 | 작업을 재개하고 데이터 새로 고침 |
| `common:WindowLostFocus` | 창이 포커스를 잃음 | 작업을 일시 중지하고 상태 저장 |
| `common:WindowMinimise` | 창이 최소화됨 | 렌더링을 일시 중지하고 리소스 사용량 줄이기 |
| `common:WindowMaximise` | 창이 최대화됨 | 전체 화면에 맞게 레이아웃 조정 |
| `common:WindowRestore` | 최소화 또는 최대화된 창이 복원됨 | 일반 레이아웃으로 복귀 |
| `common:WindowClosing` | 창이 곧 닫힘 | 데이터를 저장하고 리소스 정리 |
| `common:WindowFilesDropped` | 창에 파일이 드롭됨 | 파일 가져오기 처리 |
| `common:WindowDidResize` | 창 크기가 조정됨 | 레이아웃을 조정하고 차트를 다시 렌더링 |
| `common:WindowDidMove` | 창이 이동됨 | 위치 종속 기능 업데이트 |

### 플랫폼별 이벤트

#### Windows 이벤트

Windows 애플리케이션의 주요 이벤트:

| 이벤트 | 설명 | 사용 사례 |
| --- | --- | --- |
| `windows:SystemThemeChanged` | Windows 테마 변경됨 | 앱 색상 업데이트 |
| `windows:APMSuspend` | 시스템이 절전 모드로 전환 중 | 상태를 저장하고 작업 일시 중지 |
| `windows:APMResumeAutomatic` | 시스템이 재개됨(재개 시 항상 발생) | 상태를 복원하고 데이터 새로 고침 |
| `windows:APMResumeSuspend` | 사용자 입력으로 시스템이 재개됨(`APMResumeAutomatic` 이후) | 사용자가 시작한 절전 해제 구분 |
| `windows:APMPowerStatusChange` | 전원 상태 변경됨 | 성능 설정 조정 |

#### macOS 이벤트

주요 macOS 애플리케이션 이벤트:

| 이벤트 | 설명 | 사용 사례 |
| --- | --- | --- |
| `mac:ApplicationDidBecomeActive` | 앱이 활성 상태가 됨 | 작업 재개 |
| `mac:ApplicationDidResignActive` | 앱이 비활성 상태가 됨 | 작업 일시 중지 |
| `mac:ApplicationWillTerminate` | 앱이 곧 종료됨 | 최종 정리 |
| `mac:ApplicationWillSleep` | 시스템이 곧 절전 모드로 전환됨 | 상태를 저장하고 소켓 닫기 |
| `mac:ApplicationDidWake` | 시스템이 재개됨 | 다시 연결하고 새로 고침 |
| `mac:ApplicationScreensDidSleep` | 디스플레이가 절전 모드로 전환됨 | 렌더링 일시 중지(시스템 절전 모드와 별개) |
| `mac:ApplicationScreensDidWake` | 디스플레이가 절전 모드에서 해제됨 | 렌더링 재개 |
| `mac:WindowDidEnterFullScreen` | 전체 화면 모드로 전환됨 | 전체 화면에 맞게 UI 조정 |
| `mac:WindowDidExitFullScreen` | 전체 화면 모드에서 나옴 | 일반 UI 복원 |

#### Linux 이벤트

주요 Linux 창 이벤트:

| 이벤트 | 설명 | 사용 사례 |
| --- | --- | --- |
| `linux:SystemThemeChanged` | 데스크톱 테마 변경됨 | 앱 테마 업데이트 |
| `linux:SystemWillSleep` | 시스템이 곧 절전 모드로 전환됨(logind) | 상태 저장 |
| `linux:SystemDidWake` | 시스템이 재개됨(logind) | 다시 연결하고 새로 고침 |
| `linux:WindowFocusIn` | 창이 포커스를 얻음 | 작업 재개 |
| `linux:WindowFocusOut` | 창이 포커스를 잃음 | 작업 일시 중지 |
| `linux:WindowLoadStarted` | WebView가 로드를 시작함 | 로딩 표시기 표시 |
| `linux:WindowLoadRedirected` | WebView가 리디렉션됨 | 탐색 리디렉션 추적 |
| `linux:WindowLoadCommitted` | WebView가 로드를 커밋함 | 콘텐츠 수신 중 |
| `linux:WindowLoadFinished` | WebView 로드 완료 | 로딩 표시기 숨기기, JS/CSS 삽입 |

## 모범 사례

### 1. 이벤트 네임스페이스 사용

사용자 지정 이벤트를 만들 때는 충돌을 방지하도록 네임스페이스를 사용하세요:

```javascript
import { Events } from '@wailsio/runtime';

// Good - namespaced events
Events.Emit('myapp:user:login');
Events.Emit('myapp:data:updated');
Events.Emit('myapp:network:connected');

// Avoid - generic names that might conflict
Events.Emit('login');
Events.Emit('update');
```

### 2. 리스너 정리

컴포넌트가 언마운트될 때는 항상 이벤트 리스너를 제거하세요:

```javascript
import { Events } from '@wailsio/runtime';

// React example
useEffect(() => {
    const handler = (event) => {
        // Handle event
    };

    const off = Events.On('common:WindowDidResize', handler);

    // Cleanup — call the unsubscribe returned by Events.On
    return () => {
        off();
    };
}, []);
```

### 3. 플랫폼 차이 처리

플랫폼별 이벤트를 사용할 때는 해당 플랫폼에서 이벤트를 사용할 수 있는지 확인하세요:

```javascript
import { Events } from '@wailsio/runtime';

// Platform-specific events can be registered unconditionally;
// they will simply never fire on unsupported platforms.
Events.On('windows:APMSuspend', handleSuspend);
Events.On('mac:ApplicationWillTerminate', handleTerminate);
```

### 4. 이벤트를 과도하게 사용하지 않기

이벤트는 강력하지만 모든 작업에 사용해서는 안 됩니다:

- ✅ 이벤트를 사용하기 적합한 경우: 시스템 알림, 수명 주기 변경, 브로드캐스트 업데이트
- ❌ 이벤트 사용을 피해야 하는 경우: 직접적인 함수 반환, 단일 컴포넌트 업데이트, 동기 작업

## 이벤트 디버깅

이벤트 문제를 디버깅하려면 다음을 수행하세요:

```javascript
import { Events } from '@wailsio/runtime';

// Log all events (development only)
if (isDevelopment) {
    const originalOn = Events.On;
    Events.On = function(eventName, handler) {
        console.log(`[Event Registered] ${eventName}`);
        return originalOn.call(this, eventName, function(event) {
            console.log(`[Event Fired] ${eventName}`, event);
            return handler(event);
        });
    };
}
```

## 공식 기준 소스

사용 가능한 이벤트의 전체 목록은 Wails 소스 코드에서 확인할 수 있습니다:

- 프런트엔드 이벤트: [`v3/internal/runtime/desktop/@wailsio/runtime/src/event_types.ts`](https://github.com/wailsapp/wails/blob/main/v3/internal/runtime/desktop/@wailsio/runtime/src/event_types.ts)
- 백엔드 이벤트: [`v3/pkg/events/events.go`](https://github.com/wailsapp/wails/blob/main/v3/pkg/events/events.go)

최신 이벤트 이름과 사용 가능 여부는 항상 이 파일들을 참조하세요.

## 요약

Wails의 이벤트는 애플리케이션 내 통신을 처리하는 강력하고 결합도가 낮은 방법을 제공합니다. 이 가이드의 패턴과 모범 사례를 따르면 시스템 변경과 사용자 상호 작용에 원활하게 반응하는 응답성 높은 플랫폼 인식 애플리케이션을 구축할 수 있습니다.

기억하세요. 크로스 플랫폼 호환성을 위해 공통 이벤트부터 사용하고, 필요할 때 플랫폼별 이벤트를 추가하며, 메모리 누수를 방지하도록 이벤트 리스너를 항상 정리하세요.
