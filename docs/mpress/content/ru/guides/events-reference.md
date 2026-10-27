---
title: "Руководство по событиям"
description: "Практическое руководство по использованию событий в Wails v3 для взаимодействия компонентов приложения и управления его жизненным циклом"
slug: "guides/events-reference"
sourcePath: "guides/events-reference.md"
---

**ПРИМЕЧАНИЕ. Это руководство ещё не завершено**

## Руководство по событиям

События лежат в основе взаимодействия в приложениях Wails. Они позволяют разным частям приложения обмениваться данными без жёсткой связи между ними. В этом руководстве описано всё, что нужно знать для эффективного использования событий в приложении Wails.

## Основные сведения о событиях Wails

События можно представить как сообщения, рассылаемые по всему приложению. Любая часть приложения может прослушивать эти сообщения и соответствующим образом реагировать на них. Это особенно полезно в следующих случаях:

- **Реагирование на изменения окна**: определяйте, когда окно свёрнуто, развёрнуто или перемещено
- **Обработка системных событий**: реагируйте на смену темы или события управления питанием
- **Пользовательская логика приложения**: создавайте собственные события для таких функций, как обновление данных или действия пользователя
- **Взаимодействие между компонентами**: позволяйте разным частям приложения взаимодействовать без прямых зависимостей

## Соглашение об именовании событий

Все события Wails именуются по схеме с пространствами имён, которая позволяет однозначно определить их источник:

- `common:` — кроссплатформенные события, работающие в Windows, macOS и Linux
- `windows:` — события только для Windows
- `mac:` — события только для macOS\
- `linux:` — события только для Linux

Например:

- `common:WindowFocus` — окно получило фокус (работает на всех платформах)
- `windows:APMSuspend` — система переходит в спящий режим (только Windows)
- `mac:ApplicationDidBecomeActive` — приложение стало активным (только macOS)

## Начало работы с событиями

### Прослушивание событий (фронтенд)

Наиболее распространённый сценарий — прослушивание событий в коде фронтенда:

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

### Отправка событий (бэкенд)

Из кода Go можно отправлять события, которые будет прослушивать фронтенд:

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

### Отправка событий (фронтенд)

Хотя это используется реже, события также можно отправлять из фронтенда, чтобы их прослушивал код Go:

```javascript
import { Events } from '@wailsio/runtime';

// Event without data
Events.Emit('myapp:close-window')

// Event with data
Events.Emit('myapp:disconnect-requested', 'id-123')
```

Если во фронтенде используется TypeScript, а в коде Go [регистрируются типизированные события](#----), будут доступны автодополнение и проверка имён событий, а также проверка типов данных.

### Удаление обработчиков событий

Всегда удаляйте обработчики событий, когда они больше не нужны:

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

## Распространённые сценарии использования

### 1. Приостановка и возобновление при изменении фокуса окна

Многим приложениям необходимо приостанавливать определённые действия, когда окно теряет фокус:

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

### 2. Реагирование на смену темы

Синхронизируйте тему приложения с системной темой:

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

### 3. Обработка перетаскивания файлов

Добавьте в приложение возможность принимать перетаскиваемые файлы:

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

### 4. Управление жизненным циклом окна

Реагируйте на изменения состояния окна:

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

### 5. Возможности для отдельных платформ

При необходимости обрабатывайте события, относящиеся к конкретной платформе:

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

## Создание собственных событий

Для специфических потребностей приложения можно создавать собственные события.

### Бэкенд (Go)

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

### Фронтенд (JavaScript)

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

## Типизированные события с контролем типов

Wails v3 поддерживает типизированные события с полным контролем типов TypeScript благодаря регистрации событий и автоматической генерации привязок.

### Регистрация собственных событий

Вызовите `application.RegisterEvent` во время инициализации, чтобы зарегистрировать имена собственных событий и типы их данных:

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
`RegisterEvent` предназначена для вызова во время инициализации и вызовет панику, если:

- Аргументы недопустимы
- Одно имя события дважды зарегистрировано с разными типами данных

@end

@note{type="info"}
Одно и то же событие можно безопасно регистрировать несколько раз, если тип данных всегда одинаков. Это позволяет гарантировать регистрацию события при загрузке любого из нескольких пакетов.

@end

### Преимущества регистрации событий

После регистрации тип аргументов данных, переданных в `Event.Emit`, проверяется на соответствие указанному типу. При несовпадении:

- Формируется и записывается в журнал ошибка (либо передаётся зарегистрированному обработчику ошибок)
- Событие, вызвавшее ошибку, не распространяется
- Таким образом, поле данных зарегистрированных событий всегда совместимо с объявленным типом

### Строгий режим

Используйте тег сборки `strictevents`, чтобы во время разработки включить предупреждения о незарегистрированных событиях:

```bash
go build -tags strictevents
```

Когда строгий режим включён, среда выполнения выводит не более одного предупреждения для каждого имени незарегистрированного события, чтобы не засорять журналы.

### Генерация привязок TypeScript

Генератор привязок создаёт определения TypeScript и связующий код для прозрачной поддержки типизированных событий во фронтенде.

#### 1. Настройка плагина Vite

В файле `vite.config.ts`:

```typescript
import { defineConfig } from 'vite'
import wails from '@wailsio/runtime/plugins/vite'

export default defineConfig({
  plugins: [wails()],
})
```

#### 2. Генерация привязок

Запустите генератор привязок:

```bash
wails3 generate bindings
```

В результате в каталоге фронтенда создаются файлы TypeScript с типизированными определениями данных событий, а также модуль `bindings/events`. Он экспортирует по одной константе для каждого зарегистрированного имени события, если это константа времени компиляции, которую может обнаружить генератор (например, `user:login` преобразуется в `UserLogin`).

#### 3. Использование типизированных событий во фронтенде

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

Типизированные события предоставляют:

- **Автодополнение** имён событий
- **Проверку типов** данных событий
- **Ошибки времени компиляции** при несовпадении типов данных
- **Документацию IntelliSense**

## Справочник событий

### Общие события (кроссплатформенные)

Эти события работают на всех платформах:

| Событие | Описание | Когда использовать |
| --- | --- | --- |
| `common:ApplicationStarted` | Приложение полностью запущено | Инициализировать приложение, загрузить сохранённое состояние |
| `common:WindowRuntimeReady` | Среда выполнения Wails готова | Начать вызывать API Wails |
| `common:ThemeChanged` | Системная тема изменилась | Обновить внешний вид приложения |
| `common:SystemWillSleep` | Система готовится перейти в спящий режим | Сохранить состояние, закрыть сокеты |
| `common:SystemDidWake` | Система вышла из спящего режима | Восстановить подключения, обновить устаревшие данные |
| `common:WindowFocus` | Окно получило фокус | Возобновить работу, обновить данные |
| `common:WindowLostFocus` | Окно потеряло фокус | Приостановить работу, сохранить состояние |
| `common:WindowMinimise` | Окно было свёрнуто | Приостановить отрисовку, сократить потребление ресурсов |
| `common:WindowMaximise` | Окно было развёрнуто | Адаптировать макет к полноэкранному режиму |
| `common:WindowRestore` | Окно восстановлено после сворачивания или разворачивания | Вернуться к обычному макету |
| `common:WindowClosing` | Окно готовится закрыться | Сохранить данные, освободить ресурсы |
| `common:WindowFilesDropped` | Файлы перетащили в окно | Обработать импорт файлов |
| `common:WindowDidResize` | Размер окна изменился | Адаптировать макет, перерисовать диаграммы |
| `common:WindowDidMove` | Окно было перемещено | Обновить функции, зависящие от положения окна |

### События для отдельных платформ

#### События Windows

Основные события приложений Windows:

| Событие | Описание | Сценарий использования |
| --- | --- | --- |
| `windows:SystemThemeChanged` | Тема Windows изменилась | Обновить цвета приложения |
| `windows:APMSuspend` | Система переходит в спящий режим | Сохранить состояние, приостановить операции |
| `windows:APMResumeAutomatic` | Система вышла из спящего режима (событие возникает всегда при возобновлении работы) | Восстановить состояние, обновить данные |
| `windows:APMResumeSuspend` | Система возобновила работу после ввода пользователя (после `APMResumeAutomatic`) | Определить пробуждение по инициативе пользователя |
| `windows:APMPowerStatusChange` | Состояние питания изменилось | Скорректировать параметры производительности |

#### События macOS

Важные события приложений macOS:

| Событие | Описание | Сценарий использования |
| --- | --- | --- |
| `mac:ApplicationDidBecomeActive` | Приложение стало активным | Возобновить операции |
| `mac:ApplicationDidResignActive` | Приложение стало неактивным | Приостановить операции |
| `mac:ApplicationWillTerminate` | Приложение готовится завершить работу | Выполнить окончательную очистку |
| `mac:ApplicationWillSleep` | Система готовится перейти в спящий режим | Сохранить состояние, закрыть сокеты |
| `mac:ApplicationDidWake` | Система возобновила работу | Восстановить подключения, обновить данные |
| `mac:ApplicationScreensDidSleep` | Дисплеи перешли в спящий режим | Приостановить отрисовку (в отличие от перехода системы в спящий режим) |
| `mac:ApplicationScreensDidWake` | Дисплеи вышли из спящего режима | Возобновить отрисовку |
| `mac:WindowDidEnterFullScreen` | Выполнен переход в полноэкранный режим | Адаптировать интерфейс к полноэкранному режиму |
| `mac:WindowDidExitFullScreen` | Выполнен выход из полноэкранного режима | Восстановить обычный интерфейс |

#### События Linux

Основные события окон Linux:

| Событие | Описание | Сценарий использования |
| --- | --- | --- |
| `linux:SystemThemeChanged` | Тема рабочего стола изменилась | Обновить тему приложения |
| `linux:SystemWillSleep` | Система готовится перейти в спящий режим (logind) | Сохранить состояние |
| `linux:SystemDidWake` | Система возобновила работу (logind) | Восстановить подключения, обновить данные |
| `linux:WindowFocusIn` | Окно получило фокус | Возобновить работу |
| `linux:WindowFocusOut` | Окно потеряло фокус | Приостановить работу |
| `linux:WindowLoadStarted` | WebView начал загрузку | Показать индикатор загрузки |
| `linux:WindowLoadRedirected` | WebView перенаправлен | Отслеживание перенаправлений при навигации |
| `linux:WindowLoadCommitted` | WebView подтвердил загрузку | Получение содержимого |
| `linux:WindowLoadFinished` | WebView завершил загрузку | Скрытие индикатора загрузки, внедрение JS/CSS |

## Рекомендации

### 1. Используйте пространства имён событий

При создании пользовательских событий используйте пространства имён, чтобы избежать конфликтов:

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

### 2. Удаляйте обработчики событий

Всегда удаляйте обработчики событий при размонтировании компонентов:

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

### 3. Учитывайте различия между платформами

При использовании событий, специфичных для отдельных платформ, проверяйте их доступность на текущей платформе:

```javascript
import { Events } from '@wailsio/runtime';

// Platform-specific events can be registered unconditionally;
// they will simply never fire on unsupported platforms.
Events.On('windows:APMSuspend', handleSuspend);
Events.On('mac:ApplicationWillTerminate', handleTerminate);
```

### 4. Не злоупотребляйте событиями

Несмотря на широкие возможности событий, не используйте их для решения любых задач:

- ✅ Используйте события для системных уведомлений, изменений жизненного цикла и широковещательных обновлений
- ❌ Не используйте события для прямого возврата значений из функций, обновления одного компонента и синхронных операций

## Отладка событий

Чтобы устранить проблемы с событиями:

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

## Первоисточник

Полный список доступных событий можно найти в исходном коде Wails:

- События фронтенда: [`v3/internal/runtime/desktop/@wailsio/runtime/src/event_types.ts`](https://github.com/wailsapp/wails/blob/main/v3/internal/runtime/desktop/@wailsio/runtime/src/event_types.ts)
- События бэкенда: [`v3/pkg/events/events.go`](https://github.com/wailsapp/wails/blob/main/v3/pkg/events/events.go)

Всегда сверяйтесь с этими файлами, чтобы получить актуальные сведения о названиях и доступности событий.

## Итоги

События в Wails обеспечивают мощный и слабосвязанный механизм взаимодействия в приложении. Следуя шаблонам и рекомендациям из этого руководства, вы сможете создавать отзывчивые приложения, учитывающие особенности платформ и плавно реагирующие на системные изменения и действия пользователя.

Помните: для кроссплатформенной совместимости начинайте с общих событий, при необходимости добавляйте события для конкретных платформ и всегда удаляйте обработчики событий, чтобы предотвратить утечки памяти.
