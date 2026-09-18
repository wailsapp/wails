---
title: "事件指南"
description: "在 Wails v3 中使用事件实现应用通信和生命周期管理的实用指南"
slug: "guides/events-reference"
sourcePath: "guides/events-reference.md"
---

**注意：本指南仍在编写中**

## 事件指南

事件是 Wails 应用中通信的核心。应用的不同部分可以通过事件相互通信，而无需紧密耦合。本指南将全面介绍如何在 Wails 应用中有效地使用事件。

## 了解 Wails 事件

可以将事件看作在整个应用中广播的消息。应用的任何部分都可以监听这些消息并作出相应响应。事件尤其适用于：

- **响应窗口变化**：获知窗口何时最小化、最大化或移动
- **处理系统事件**：响应主题变化或电源事件
- **自定义应用逻辑**：为数据更新或用户操作等功能创建自己的事件
- **跨组件通信**：让应用的不同部分在没有直接依赖关系的情况下进行通信

## 事件命名约定

所有 Wails 事件都遵循命名空间模式，以明确表明其来源：

- `common:` - 可在 Windows、macOS 和 Linux 上使用的跨平台事件
- `windows:` - Windows 专用事件
- `mac:` - macOS 专用事件\
- `linux:` - Linux 专用事件

例如：

- `common:WindowFocus` - 窗口获得焦点（适用于所有平台）
- `windows:APMSuspend` - 系统正在挂起（仅限 Windows）
- `mac:ApplicationDidBecomeActive` - 应用进入活动状态（仅限 macOS）

## 事件入门

### 监听事件（前端）

最常见的用法是在前端代码中监听事件：

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

### 触发事件（后端）

可以从 Go 代码中触发事件，供前端监听：

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

### 触发事件（前端）

虽然不太常用，但也可以从前端触发事件，供 Go 代码监听：

```javascript
import { Events } from '@wailsio/runtime';

// Event without data
Events.Emit('myapp:close-window')

// Event with data
Events.Emit('myapp:disconnect-requested', 'id-123')
```

如果前端使用 TypeScript，并在 Go 代码中[注册类型化事件](#heading-9)，即可获得事件名称自动补全和检查以及数据类型检查。

### 移除事件监听器

不再需要事件监听器时，务必将其清理：

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

## 常见用例

### 1. 窗口焦点变化时暂停/恢复

许多应用需要在窗口失去焦点时暂停某些活动：

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

### 2. 响应主题变化

让应用与系统主题保持同步：

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

### 3. 处理文件拖放

让应用接受拖入的文件：

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

### 4. 窗口生命周期管理

响应窗口状态变化：

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

### 5. 平台专用功能

按需处理平台专用事件：

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

## 创建自定义事件

可以根据应用的特定需求创建自己的事件。

### 后端（Go）

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

## 具有类型安全保障的类型化事件

Wails v3 通过事件注册和自动生成绑定，支持具有完整 TypeScript 类型安全保障的类型化事件。

### 注册自定义事件

在初始化时调用`application.RegisterEvent`，注册自定义事件名称及其数据类型：

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
`RegisterEvent`应在初始化时调用，并会在以下情况下触发 panic：

- 参数无效
- 使用不同的数据类型重复注册同一事件名称

@end

@note{type="info"}
只要数据类型始终相同，多次注册同一事件就是安全的。当加载多个包中的任意一个时都需要确保某个事件已注册，这种方式会很有用。

@end

### 事件注册的优势

注册后，系统会根据指定类型对传给`Event.Emit`的数据参数进行类型检查。如果类型不匹配：

- 系统会产生并记录错误（或将错误传给已注册的错误处理程序）
- 不会传播引发错误的事件
- 这可确保已注册事件的 data 字段始终可赋值给声明的类型

### 严格模式

使用 `strictevents` 构建标签，在开发环境中为未注册的事件启用警告：

```bash
go build -tags strictevents
```

启用严格模式后，为避免日志中出现大量重复信息，运行时针对每个未注册的事件名称最多发出一次警告。

### 生成 TypeScript 绑定

绑定生成器会输出 TypeScript 定义和粘合代码，为前端提供透明的类型化事件支持。

#### 1. 配置 Vite 插件

在你的 `vite.config.ts` 中：

```typescript
import { defineConfig } from 'vite'
import wails from '@wailsio/runtime/plugins/vite'

export default defineConfig({
  plugins: [wails()],
})
```

#### 2. 生成绑定

运行绑定生成器：

```bash
wails3 generate bindings
```

这会在前端目录中创建包含类型化事件创建器和数据接口的 TypeScript 文件。

#### 3. 在前端使用类型化事件

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

类型化事件提供：

- 事件名称的<strong>自动补全</strong>
- 事件数据的<strong>类型检查</strong>
- 数据类型不匹配时的<strong>编译时错误</strong>
- **IntelliSense** 文档

## 事件参考

### 通用事件（跨平台）

以下事件适用于所有平台：

| 事件 | 说明 | 使用场景 |
| --- | --- | --- |
| `common:ApplicationStarted` | 应用程序已完全启动 | 初始化应用程序并加载已保存的状态 |
| `common:WindowRuntimeReady` | Wails 运行时已就绪 | 开始调用 Wails API |
| `common:ThemeChanged` | 系统主题已更改 | 更新应用程序外观 |
| `common:SystemWillSleep` | 系统即将挂起 | 将状态写入持久存储并关闭套接字 |
| `common:SystemDidWake` | 系统已从挂起状态恢复 | 重新连接并刷新过期数据 |
| `common:WindowFocus` | 窗口已获得焦点 | 恢复活动并刷新数据 |
| `common:WindowLostFocus` | 窗口已失去焦点 | 暂停活动并保存状态 |
| `common:WindowMinimise` | 窗口已最小化 | 暂停渲染并减少资源使用 |
| `common:WindowMaximise` | 窗口已最大化 | 针对全屏调整布局 |
| `common:WindowRestore` | 窗口已从最小化或最大化状态还原 | 恢复正常布局 |
| `common:WindowClosing` | 窗口即将关闭 | 保存数据并清理资源 |
| `common:WindowFilesDropped` | 文件已拖放到窗口中 | 处理文件导入 |
| `common:WindowDidResize` | 窗口大小已调整 | 调整布局并重新渲染图表 |
| `common:WindowDidMove` | 窗口已移动 | 更新依赖位置的功能 |

### 平台特定事件

#### Windows 事件

Windows 应用程序的关键事件：

| 事件 | 说明 | 用例 |
| --- | --- | --- |
| `windows:SystemThemeChanged` | Windows 主题已更改 | 更新应用颜色 |
| `windows:APMSuspend` | 系统正在挂起 | 保存状态并暂停操作 |
| `windows:APMResumeAutomatic` | 系统已恢复（恢复时始终触发） | 恢复状态并刷新数据 |
| `windows:APMResumeSuspend` | 系统通过用户输入恢复（在`APMResumeAutomatic`之后） | 识别由用户发起的唤醒 |
| `windows:APMPowerStatusChange` | 电源状态已更改 | 调整性能设置 |

#### macOS 事件

重要的 macOS 应用程序事件：

| 事件 | 说明 | 用例 |
| --- | --- | --- |
| `mac:ApplicationDidBecomeActive` | 应用已变为活跃状态 | 恢复操作 |
| `mac:ApplicationDidResignActive` | 应用已变为非活跃状态 | 暂停操作 |
| `mac:ApplicationWillTerminate` | 应用即将退出 | 执行最终清理 |
| `mac:ApplicationWillSleep` | 系统即将挂起 | 保存状态并关闭套接字 |
| `mac:ApplicationDidWake` | 系统已恢复 | 重新连接并刷新 |
| `mac:ApplicationScreensDidSleep` | 显示器已进入睡眠状态 | 暂停渲染（不同于系统睡眠） |
| `mac:ApplicationScreensDidWake` | 显示器已唤醒 | 恢复渲染 |
| `mac:WindowDidEnterFullScreen` | 已进入全屏模式 | 针对全屏模式调整 UI |
| `mac:WindowDidExitFullScreen` | 已退出全屏模式 | 恢复常规 UI |

#### Linux 事件

Linux 的核心窗口事件：

| 事件 | 说明 | 用例 |
| --- | --- | --- |
| `linux:SystemThemeChanged` | 桌面主题已更改 | 更新应用主题 |
| `linux:SystemWillSleep` | 系统即将挂起（logind） | 保存状态 |
| `linux:SystemDidWake` | 系统已恢复（logind） | 重新连接并刷新 |
| `linux:WindowFocusIn` | 窗口已获得焦点 | 恢复活动 |
| `linux:WindowFocusOut` | 窗口失去焦点 | 暂停活动 |
| `linux:WindowLoadStarted` | WebView 开始加载 | 显示加载指示器 |
| `linux:WindowLoadRedirected` | WebView 已重定向 | 跟踪导航重定向 |
| `linux:WindowLoadCommitted` | WebView 已提交加载 | 正在接收内容 |
| `linux:WindowLoadFinished` | WebView 已完成加载 | 隐藏加载指示器，注入 JS/CSS |

## 最佳实践

### 1. 使用事件命名空间

创建自定义事件时，请使用命名空间以避免冲突：

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

### 2. 清理监听器

组件卸载时，务必移除事件监听器：

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

### 3. 处理平台差异

使用平台特定事件时，请检查其在当前平台上是否可用：

```javascript
import { Events } from '@wailsio/runtime';

// Platform-specific events can be registered unconditionally;
// they will simply never fire on unsupported platforms.
Events.On('windows:APMSuspend', handleSuspend);
Events.On('mac:ApplicationWillTerminate', handleTerminate);
```

### 4. 不要过度使用事件

事件虽然功能强大，但不要将其用于所有场景：

- ✅ 适合使用事件的场景：系统通知、生命周期变化、广播更新
- ❌ 避免使用事件的场景：直接返回函数结果、更新单个组件、同步操作

## 调试事件

若要调试事件问题：

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

## 权威来源

可用事件的完整列表可以在 Wails 源代码中找到：

- 前端事件：[`v3/internal/runtime/desktop/@wailsio/runtime/src/event_types.ts`](https://github.com/wailsapp/wails/blob/main/v3/internal/runtime/desktop/@wailsio/runtime/src/event_types.ts)
- 后端事件：[`v3/pkg/events/events.go`](https://github.com/wailsapp/wails/blob/main/v3/pkg/events/events.go)

请始终以这些文件为准，获取最新的事件名称及其可用性信息。

## 总结

Wails 中的事件提供了一种强大且解耦的应用内通信方式。遵循本指南中的模式和实践，即可构建响应迅速、感知平台差异的应用，使其能够顺畅响应系统变化和用户交互。

请记住：首先使用通用事件以确保跨平台兼容性，在需要时添加平台特定事件，并始终清理事件监听器以防止内存泄漏。
