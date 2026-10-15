---
title: "事件指南"
description: "在 Wails v3 中使用事件進行應用程式通訊與生命週期管理的實用指南"
slug: "guides/events-reference"
sourcePath: "guides/events-reference.md"
---

**注意：本指南仍在撰寫中**

## 事件指南

事件是 Wails 應用程式通訊的核心。事件讓應用程式的不同部分彼此溝通，而不會形成緊密耦合。本指南將帶您瞭解在 Wails 應用程式中有效使用事件所需的一切知識。

## 瞭解 Wails 事件

您可以將事件視為在整個應用程式中廣播的訊息。應用程式的任何部分都能監聽這些訊息，並據此回應。事件尤其適用於：

- **回應視窗變更**：得知視窗何時最小化、最大化或移動
- **處理系統事件**：回應佈景主題變更或電源事件
- **自訂應用程式邏輯**：針對資料更新或使用者動作等功能建立自己的事件
- **跨元件通訊**：讓應用程式的不同部分在沒有直接相依性的情況下進行通訊

## 事件命名慣例

所有 Wails 事件都遵循命名空間模式，以清楚指出其來源：

- `common:` - 可在 Windows、macOS 與 Linux 上運作的跨平台事件
- `windows:` - Windows 專用事件
- `mac:` - macOS 專用事件\
- `linux:` - Linux 專用事件

例如：

- `common:WindowFocus` - 視窗取得焦點（適用於所有平台）
- `windows:APMSuspend` - 系統正在暫停（僅限 Windows）
- `mac:ApplicationDidBecomeActive` - 應用程式已進入作用中狀態（僅限 macOS）

## 開始使用事件

### 監聽事件（前端）

最常見的使用情境是在前端程式碼中監聽事件：

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

### 發出事件（後端）

您可以從 Go 程式碼發出事件，供前端監聽：

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

### 發出事件（前端）

雖然較不常用，您也可以從前端發出事件，供 Go 程式碼監聽：

```javascript
import { Events } from '@wailsio/runtime';

// Event without data
Events.Emit('myapp:close-window')

// Event with data
Events.Emit('myapp:disconnect-requested', 'id-123')
```

如果您在前端使用 TypeScript，並在 Go 程式碼中[註冊具型別事件](#heading-9)，即可取得事件名稱的自動完成與檢查，以及資料型別檢查。

### 移除事件監聽器

不再需要事件監聽器時，務必將其清除：

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

## 常見使用情境

### 1. 視窗焦點變更時暫停／繼續

許多應用程式需要在視窗失去焦點時暫停某些活動：

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

### 2. 回應佈景主題變更

讓應用程式與系統佈景主題保持同步：

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

### 3. 處理檔案拖放

讓應用程式接受拖入的檔案：

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

### 4. 視窗生命週期管理

回應視窗狀態變更：

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

### 5. 平台專用功能

視需要處理平台專用事件：

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

## 建立自訂事件

您可以針對應用程式的特定需求建立自己的事件。

### 後端（Go）

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

### 前端（JavaScript）

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

## 具型別安全性的事件

Wails v3 可透過事件註冊與自動產生繫結，支援具備完整 TypeScript 型別安全性的具型別事件。

### 註冊自訂事件

在初始化時呼叫`application.RegisterEvent`，以註冊自訂事件名稱及其資料型別：

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
`RegisterEvent`應在初始化時呼叫，並會在下列情況下引發 panic：

- 引數無效
- 以不同的資料型別重複註冊相同的事件名稱

@end

@note{type="info"}
只要資料型別始終相同，就能安全地多次註冊同一事件。當多個套件中的任何一個載入時，這有助於確保事件已註冊。

@end

### 註冊事件的優點

註冊後，系統會根據指定型別，對傳遞給`Event.Emit`的資料引數進行型別檢查。如果型別不相符：

- 系統會發出並記錄錯誤（或將錯誤傳遞給已註冊的錯誤處理常式）
- 不會傳播有問題的事件
- 這可確保已註冊事件的資料欄位一律可指派給宣告的型別

### 嚴格模式

使用`strictevents`建置標籤，在開發期間為未註冊的事件啟用警告：

```bash
go build -tags strictevents
```

啟用嚴格模式後，執行階段對每個未註冊的事件名稱最多發出一次警告，以免記錄檔充斥重複訊息。

### 產生 TypeScript 繫結

繫結產生器會輸出 TypeScript 定義和黏合程式碼，讓前端能以透明方式支援型別化事件。

#### 1. 設定 Vite 外掛程式

在您的`vite.config.ts`中：

```typescript
import { defineConfig } from 'vite'
import wails from '@wailsio/runtime/plugins/vite'

export default defineConfig({
  plugins: [wails()],
})
```

#### 2. 產生繫結

執行繫結產生器：

```bash
wails3 generate bindings
```

這會在前端目錄中建立包含型別化事件建立函式和資料介面的 TypeScript 檔案。

#### 3. 在前端使用型別化事件

```typescript
import { Events } from '@wailsio/runtime'
import { UserLogin, MonitorStats } from './bindings/events'

// Type-safe event emission with autocomplete
Events.Emit(UserLogin({
    UserID: "123",
    Username: "john_doe",
    LoginTime: new Date().toISOString()
}))

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

型別化事件提供：

- 事件名稱的<strong>自動完成</strong>
- 事件資料的<strong>型別檢查</strong>
- 資料型別不相符時的<strong>編譯期錯誤</strong>
- <strong>IntelliSense</strong>文件

## 事件參考

### 通用事件（跨平台）

這些事件適用於所有平台：

| 事件 | 說明 | 使用時機 |
| --- | --- | --- |
| `common:ApplicationStarted` | 應用程式已完全啟動 | 初始化應用程式、載入已儲存的狀態 |
| `common:WindowRuntimeReady` | Wails 執行階段已就緒 | 開始呼叫 Wails API |
| `common:ThemeChanged` | 系統佈景主題已變更 | 更新應用程式外觀 |
| `common:SystemWillSleep` | 系統即將暫停 | 將狀態寫入儲存裝置、關閉通訊端 |
| `common:SystemDidWake` | 系統已從暫停狀態恢復 | 重新連線、重新整理過期資料 |
| `common:WindowFocus` | 視窗取得焦點 | 繼續執行活動、重新整理資料 |
| `common:WindowLostFocus` | 視窗失去焦點 | 暫停活動、儲存狀態 |
| `common:WindowMinimise` | 視窗已最小化 | 暫停轉譯、降低資源使用量 |
| `common:WindowMaximise` | 視窗已最大化 | 調整全螢幕版面配置 |
| `common:WindowRestore` | 視窗已從最小化或最大化狀態還原 | 恢復一般版面配置 |
| `common:WindowClosing` | 視窗即將關閉 | 儲存資料、清理資源 |
| `common:WindowFilesDropped` | 檔案已拖放至視窗 | 處理檔案匯入 |
| `common:WindowDidResize` | 視窗大小已調整 | 調整版面配置、重新轉譯圖表 |
| `common:WindowDidMove` | 視窗已移動 | 更新取決於位置的功能 |

### 平台特定事件

#### Windows 事件

Windows 應用程式的重要事件：

| 事件 | 說明 | 使用情境 |
| --- | --- | --- |
| `windows:SystemThemeChanged` | Windows 佈景主題已變更 | 更新應用程式色彩 |
| `windows:APMSuspend` | 系統正在暫停 | 儲存狀態、暫停作業 |
| `windows:APMResumeAutomatic` | 系統已恢復（恢復時一律觸發） | 還原狀態、重新整理資料 |
| `windows:APMResumeSuspend` | 系統透過使用者輸入恢復（在`APMResumeAutomatic`之後） | 辨別由使用者啟動的喚醒 |
| `windows:APMPowerStatusChange` | 電源狀態已變更 | 調整效能設定 |

#### macOS 事件

重要的 macOS 應用程式事件：

| 事件 | 說明 | 使用情境 |
| --- | --- | --- |
| `mac:ApplicationDidBecomeActive` | 應用程式已變為作用中 | 恢復作業 |
| `mac:ApplicationDidResignActive` | 應用程式已變為非作用中 | 暫停作業 |
| `mac:ApplicationWillTerminate` | 應用程式即將結束 | 執行最後清理 |
| `mac:ApplicationWillSleep` | 系統即將暫停 | 儲存狀態、關閉通訊端 |
| `mac:ApplicationDidWake` | 系統已恢復 | 重新連線、重新整理 |
| `mac:ApplicationScreensDidSleep` | 顯示器已進入睡眠 | 暫停轉譯（有別於系統睡眠） |
| `mac:ApplicationScreensDidWake` | 顯示器已喚醒 | 恢復轉譯 |
| `mac:WindowDidEnterFullScreen` | 已進入全螢幕模式 | 針對全螢幕模式調整使用者介面 |
| `mac:WindowDidExitFullScreen` | 已離開全螢幕模式 | 還原一般使用者介面 |

#### Linux 事件

Linux 的核心視窗事件：

| 事件 | 說明 | 使用情境 |
| --- | --- | --- |
| `linux:SystemThemeChanged` | 桌面佈景主題已變更 | 更新應用程式佈景主題 |
| `linux:SystemWillSleep` | 系統即將暫停（logind） | 儲存狀態 |
| `linux:SystemDidWake` | 系統已恢復（logind） | 重新連線、重新整理 |
| `linux:WindowFocusIn` | 視窗已取得焦點 | 恢復活動 |
| `linux:WindowFocusOut` | 視窗失去焦點 | 暫停活動 |
| `linux:WindowLoadStarted` | WebView 開始載入 | 顯示載入指示器 |
| `linux:WindowLoadRedirected` | WebView 已重新導向 | 追蹤導覽重新導向 |
| `linux:WindowLoadCommitted` | WebView 已確認載入 | 正在接收內容 |
| `linux:WindowLoadFinished` | WebView 已完成載入 | 隱藏載入指示器，注入 JS/CSS |

## 最佳實務

### 1. 使用事件命名空間

建立自訂事件時，請使用命名空間以避免衝突：

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

### 2. 清理監聽器

元件卸載時，務必移除事件監聽器：

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

### 3. 處理平台差異

使用平台特定事件時，請檢查該平台是否支援：

```javascript
import { Events } from '@wailsio/runtime';

// Platform-specific events can be registered unconditionally;
// they will simply never fire on unsupported platforms.
Events.On('windows:APMSuspend', handleSuspend);
Events.On('mac:ApplicationWillTerminate', handleTerminate);
```

### 4. 請勿濫用事件

事件雖然功能強大，但不要事事都使用事件：

- ✅ 適合使用事件的情況：系統通知、生命週期變更、廣播更新
- ❌ 應避免使用事件的情況：直接傳回函式結果、更新單一元件、同步操作

## 偵錯事件

若要偵錯事件問題：

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

## 權威來源

可用事件的完整清單可在 Wails 原始碼中找到：

- 前端事件：[`v3/internal/runtime/desktop/@wailsio/runtime/src/event_types.ts`](https://github.com/wailsapp/wails/blob/main/v3/internal/runtime/desktop/@wailsio/runtime/src/event_types.ts)
- 後端事件：[`v3/pkg/events/events.go`](https://github.com/wailsapp/wails/blob/main/v3/pkg/events/events.go)

請一律參考這些檔案，以取得最新的事件名稱及可用性資訊。

## 總結

Wails 中的事件為應用程式內的通訊提供了功能強大且鬆散耦合的處理方式。遵循本指南中的模式與實務，即可建置回應迅速、能感知平台差異，並順暢回應系統變更與使用者互動的應用程式。

請記住：先使用通用事件以確保跨平台相容性，視需要加入平台特定事件，並務必清理事件監聽器以防止記憶體洩漏。
