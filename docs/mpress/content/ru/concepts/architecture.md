---
title: "Как работает Wails"
description: "Архитектура Wails и способы достижения производительности нативных приложений"
slug: "concepts/architecture"
sourcePath: "concepts/architecture.md"
---

Wails — это фреймворк для создания настольных приложений, использующий **Go для бэкенда** и **веб-технологии для фронтенда**. Но, в отличие от Electron, Wails не включает браузер в комплект поставки, а использует **нативный WebView операционной системы**.

```d2
direction: left

Wails App: {
  shape: sequence_diagram
  label: Приложение Wails

  frontend: Фронтенд
  backend: Бэкенд на Go
  os: Операционная система

  Initialisation: Инициализация {
    shape: sequence_diagram
    backend."Serves Static Web App": Обслуживает статическое веб-приложение
    backend -> frontend: HTML / JS / CSS
    frontend."Render Site via OS-native WebView": Отображает сайт во встроенном WebView операционной системы
  }
  Regular Communication: Обычный обмен данными {
    shape: sequence_diagram
    frontend."Make API-style call": Вызвать API
    frontend -> backend.a: JSON
    backend.a."Service processes request": Сервис обрабатывает запрос
    backend.a -> os: Вызвать системные API
    backend.a."Generate Response": Сформировать ответ
    backend.a -> frontend: JSON
    frontend."Process response": Обработать ответ
  }
  backend.a.label: a
}
```

**Основные отличия от Electron:**

| Аспект | Wails | Electron |
| --- | --- | --- |
| **Браузер** | WebView, предоставляемый ОС | Chromium в комплекте поставки (~100 МБ) |
| **Бэкенд** | Go (компилируемый) | Node.js (интерпретируемый) |
| **Взаимодействие** | Мост в памяти | IPC (межпроцессное взаимодействие) |
| **Размер пакета** | ~15 МБ | ~150 МБ |
| **Память** | ~10 МБ | ~100 МБ+ |
| **Запуск** | &lt;0.5 с | 2-3 с |

## Основные компоненты

### 1. Нативный WebView

Wails использует встроенный в операционную систему движок веб-рендеринга:

@tabs{sync-key="platform"}
[Windows]
**WebView2** (Microsoft Edge WebView2)

- Основан на Chromium (как и браузер Edge)
- Предустановлен в Windows 10/11
- Автоматически обновляется через Windows Update
- Полностью поддерживает современные веб-стандарты

[macOS]
**WebKit** (движок рендеринга Safari)

- Встроен в macOS
- Тот же движок, что и в браузере Safari
- Отличная производительность и длительное время работы от аккумулятора
- Полностью поддерживает современные веб-стандарты

[Linux]
**WebKitGTK** (порт WebKit для GTK)

- Устанавливается через менеджер пакетов
- Тот же движок, что и в GNOME Web (Epiphany)
- Хорошая поддержка стандартов
- Низкое потребление ресурсов и высокая производительность

@end

**Почему это важно:**

- **Нет браузера в комплекте поставки** → Меньший размер приложения
- **Нативность для ОС** → Лучшая интеграция и производительность
- **Автоматические обновления** → Исправления безопасности поступают с обновлениями ОС
- **Привычный рендеринг** → Такой же, как в системном браузере

### 2. Мост Wails

Мост — это сердце Wails: он обеспечивает **прямое взаимодействие** между Go и JavaScript.

```d2
direction: down

Frontend: Фронтенд (JavaScript) {
  shape: rectangle
  style.fill: "#8B5CF6"
}

Bridge: Мост Wails {
  Encoder: Кодировщик JSON {
    shape: rectangle
  }

  Router: Маршрутизатор методов {
    shape: diamond
    style.fill: "#10B981"
  }

  Decoder: Декодировщик JSON {
    shape: rectangle
  }
}

Backend: Бэкенд (Go) {
  Services: Зарегистрированные сервисы {
    shape: rectangle
    style.fill: "#00ADD8"
  }
}

Frontend -> Bridge.Encoder: "1. Вызвать метод Go\nGreet('Alice')"
Bridge.Encoder -> Bridge.Router: "2. Закодировать в JSON\n{method: 'Greet', args: ['Alice']}"
Bridge.Router -> Backend.Services: "3. Направить сервису\nGreetService.Greet('Alice')"
Backend.Services -> Bridge.Decoder: "4. Вернуть результат\n'Hello, Alice!'"
Bridge.Decoder -> Frontend: "5. Декодировать в JS\nPromise выполняется"
```

**Как это работает:**

1. **Фронтенд вызывает метод Go** (через автоматически сгенерированную привязку)
2. **Мост кодирует вызов** в JSON (имя метода + аргументы)
3. **Маршрутизатор находит метод Go** в зарегистрированных сервисах
4. **Метод Go выполняется** и возвращает значение
5. **Мост декодирует результат** и отправляет его обратно во фронтенд
6. **Promise завершается** в JavaScript с полученным результатом

**Характеристики производительности:**

- **В памяти**: без сетевых накладных расходов и без HTTP
- **Без копирования** там, где это возможно (для больших объёмов данных)
- **Асинхронность по умолчанию**: обе стороны работают без блокировки
- **Типобезопасность**: определения TypeScript создаются автоматически

### 3. Система сервисов

Сервисы — рекомендуемый способ предоставления фронтенду функций Go.

```go
// Define a service (just a regular Go struct)
type GreetService struct {
    prefix string
}

// Methods with exported names are automatically available
func (g *GreetService) Greet(name string) string {
    return g.prefix + name + "!"
}

func (g *GreetService) GetTime() time.Time {
    return time.Now()
}

// Register the service
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&GreetService{prefix: "Hello, "}),
    },
})
```

**Обнаружение сервисов:**

- При запуске Wails **сканирует вашу структуру**
- **Экспортируемые методы** становятся доступными для вызова из фронтенда
- **Информация о типах** извлекается для привязок TypeScript
- **Обработка ошибок** выполняется автоматически (ошибки Go → исключения JS)

**Сгенерированная привязка TypeScript:**

```typescript
// Auto-generated in frontend/bindings/GreetService.ts
export function Greet(name: string): Promise<string>
export function GetTime(): Promise<Date>
```

**Зачем нужны сервисы?**

- **Типобезопасность**: полная поддержка TypeScript
- **Автоматическое обнаружение**: регистрировать методы вручную не требуется
- **Упорядоченность**: связанная функциональность сгруппирована
- **Удобство тестирования**: сервисы представляют собой обычные структуры Go

[Подробнее о сервисах →](/features/bindings/services/)

### 4. Система событий

События обеспечивают **обмен данными по модели «издатель — подписчик»** между компонентами.

```d2
direction: left

Wails Event System: Система событий Wails {
  shape: sequence_diagram

  window1: Окно 1
  window2: Окно 2
  backend: Бэкенд на Go

  Event Driver: Диспетчер событий {
    shape: sequence_diagram
    window1."Subscribe to 'data-updated' events": "Подписаться на события 'data-updated'"
    window2."Subscribe to 'data-updated' events": "Подписаться на события 'data-updated'"
    backend.a."App Emit('data-updated', data)": "Приложение вызывает Emit('data-updated', data)"
    backend.a -> window1.a: Шина событий JSON
    backend.a -> window2: Шина событий JSON
    window1.a."Subscriber processes On('data-updated', handler)": "Подписчик обрабатывает On('data-updated', handler)"
    window2."Subscriber processes On('data-updated', handler)": "Подписчик обрабатывает On('data-updated', handler)"
  }
  backend.a.label: a
  window1.a.label: a
}
```

**Варианты использования:**

- **Обмен данными между окнами**: одно окно уведомляет другие
- **Фоновые задачи**: сервис Go сообщает пользовательскому интерфейсу о ходе выполнения
- **Синхронизация состояния**: состояния нескольких окон синхронизируются
- **Слабая связанность**: компонентам не нужны прямые ссылки друг на друга

**Пример:**

```go
// Go: Emit an event
app.Event.Emit("user-logged-in", user)
```

```javascript
// JavaScript: Listen for event
import { Events } from '@wailsio/runtime'

Events.On('user-logged-in', (user) => {
    console.log('User logged in:', user)
})
```

[Подробнее о событиях →](/features/events/system/)

## Жизненный цикл приложения

Понимание жизненного цикла помогает определить, когда инициализировать ресурсы и освобождать их.

```d2
direction: down

Start: Запуск приложения {
  shape: oval
  style.fill: "#10B981"
}

Init: Инициализация {
  Create: Создать приложение {
    shape: rectangle
  }

  Register: Зарегистрировать сервисы {
    shape: rectangle
  }

  Setup: Настроить окна и меню {
    shape: rectangle
  }
}

Run: Цикл обработки событий {
  Events: Обработать события {
    shape: rectangle
  }

  Messages: Обработать сообщения {
    shape: rectangle
  }

  Render: Обновить интерфейс {
    shape: rectangle
  }
}

Shutdown: Завершение работы {
  Cleanup: Освободить ресурсы {
    shape: rectangle
  }

  Save: Сохранить состояние {
    shape: rectangle
  }
}

End: Завершение приложения {
  shape: oval
  style.fill: "#EF4444"
}

Start -> Init.Create
Init.Create -> Init.Register
Init.Register -> Init.Setup
Init.Setup -> Run.Events
Run.Events -> Run.Messages
Run.Messages -> Run.Render
Run.Render -> Run.Events: Цикл
Run.Events -> Shutdown.Cleanup: Сигнал выхода
Shutdown.Cleanup -> Shutdown.Save
Shutdown.Save -> End
```

**Обработчики жизненного цикла:**

```go
app := application.New(application.Options{
    Name: "My App",

    // Cleanly intercept quit requests (e.g. unsaved changes).
    ShouldQuit: func() bool { return true },

    // Called when the app is confirmed to be quitting — save state, close connections, etc.
    OnShutdown: func() {},
})
```

У `application.Options` нет поля `OnStartup`. Действия при запуске следует выполнять в `ServiceStartup(ctx, options)` сервиса, в обратном вызове, зарегистрированном с помощью `app.Event.OnApplicationEvent(events.Common.ApplicationStarted, ...)`, или просто перед `app.Run()`.

[Подробнее о жизненном цикле →](/concepts/lifecycle/)

## Процесс сборки

Рассмотрим, как Wails собирает приложение:

```d2
direction: down

Source: Исходный код {
  Go: "Код Go\n(main.go, сервисы)" {
    shape: rectangle
    style.fill: "#00ADD8"
  }

  Frontend: "Код фронтенда\n(HTML/CSS/JS)" {
    shape: rectangle
    style.fill: "#8B5CF6"
  }
}

Build: Процесс сборки {
  AnalyseGo: Проанализировать код Go {
    shape: rectangle
  }

  GenerateBindings: Сгенерировать привязки {
    shape: rectangle
  }

  BuildFrontend: Собрать фронтенд {
    shape: rectangle
  }

  CompileGo: Скомпилировать Go {
    shape: rectangle
  }

  Embed: Встроить ресурсы {
    shape: rectangle
  }
}

Output: Результат {
  Binary: "Нативный исполняемый файл\n(myapp.exe/.app)" {
    shape: rectangle
    style.fill: "#10B981"
  }
}

Source.Go -> Build.AnalyseGo
Build.AnalyseGo -> Build.GenerateBindings: Извлечь типы
Build.GenerateBindings -> Source.Frontend: Привязки TypeScript
Source.Frontend -> Build.BuildFrontend: Скомпилировать (Vite/webpack)
Build.BuildFrontend -> Build.Embed: Ресурсы в составе сборки
Source.Go -> Build.CompileGo
Build.CompileGo -> Build.Embed
Build.Embed -> Output.Binary
```

**Этапы сборки:**

1. **Анализ кода Go**
  - Поиск экспортируемых методов в сервисах
  - Извлечение типов параметров и возвращаемых значений
  - Создание сигнатур методов


2. **Создание привязок TypeScript**
  - Создание файлов `.ts` для каждого сервиса
  - Добавление полных определений типов
  - Добавление комментариев JSDoc


3. **Сборка фронтенда**
  - Запуск сборщика (Vite, webpack и т. п.)
  - Минификация и оптимизация
  - Вывод в `frontend/dist/`


4. **Компиляция Go**
  - Компиляция с оптимизациями (`-ldflags="-s -w"`)
  - Добавление метаданных сборки
  - Компиляция для целевой платформы


5. **Встраивание ресурсов**
  - Встраивание файлов фронтенда в исполняемый файл Go
  - Сжатие ресурсов
  - Создание единого исполняемого файла


**Результат:** единый нативный исполняемый файл со всеми встроенными ресурсами.

[Подробнее о сборке →](/guides/build/building/)

## Разработка и эксплуатация

В режимах разработки и эксплуатации Wails работает по-разному:

@tabs{sync-key="mode"}
[Разработка (wails3 dev)]
**Характеристики:**

- **Горячая перезагрузка**: изменения во фронтенде применяются мгновенно
- **Карты исходного кода**: отладка с использованием исходного кода
- **DevTools**: доступны инструменты разработчика браузера
- **Журналирование**: включено подробное журналирование
- **Внешний фронтенд**: обслуживается сервером разработки (Vite)

**Как это работает:**

```d2
direction: right

WailsApp: Приложение Wails {
  shape: rectangle
  style.fill: "#00ADD8"
}

DevServer: "Сервер разработки Vite\n(localhost:5173)" {
  shape: rectangle
  style.fill: "#8B5CF6"
}

WebView: WebView {
  shape: rectangle
  style.fill: "#6B7280"
}

WailsApp -> DevServer: Проксировать запросы
DevServer -> WebView: Обслуживать с HMR
WebView -> WailsApp: Вызывать методы Go
```

**Преимущества:**

- Мгновенная обратная связь при внесении изменений
- Полные возможности отладки
- Более быстрые итерации разработки

[Эксплуатация (wails3 build)]
**Характеристики:**

- **Встроенные ресурсы**: фронтенд включён в исполняемый файл
- **Оптимизация**: минификация и сжатие
- **Без DevTools**: по умолчанию отключены
- **Минимальное журналирование**: регистрируются только ошибки
- **Один файл**: всё содержится в одном исполняемом файле

**Как это работает:**

```d2
direction: right

Binary: "Единый исполняемый файл\n(myapp.exe)" {
  GoCode: Скомпилированный код Go {
    shape: rectangle
    style.fill: "#00ADD8"
  }

  Assets: "Встроенные ресурсы\n(HTML/CSS/JS)" {
    shape: rectangle
    style.fill: "#8B5CF6"
  }
}

WebView: WebView {
  shape: rectangle
  style.fill: "#6B7280"
}

Binary.Assets -> WebView: Обслуживать из памяти
WebView -> Binary.GoCode: Вызывать методы Go
```

**Преимущества:**

- Распространение в виде одного файла
- Меньший размер (благодаря минификации)
- Более высокая производительность
- Нет внешних зависимостей

@end

## Модель памяти

Понимание использования памяти помогает создавать эффективные приложения.

**Области памяти:**

1. **Куча Go**
  - Ваши сервисы и состояние приложения
  - Управляется сборщиком мусора Go
  - Обычно 5-10 МБ для простых приложений


2. **Память WebView**
  - DOM, куча JavaScript, CSS
  - Управляется движком WebView
  - Обычно 10-20 МБ для простых приложений


3. **Память моста**
  - Буферы сообщений для обмена данными
  - Минимальные накладные расходы (<1 МБ)
  - По возможности передача больших объёмов данных без копирования


**Советы по оптимизации:**

- **Избегайте передачи больших объёмов данных**: передавайте идентификаторы, а подробные данные загружайте по запросу
- **Используйте события для обновлений**: не выполняйте опрос из фронтенда
- **Передавайте большие файлы потоково**: не загружайте их в память целиком
- **Удаляйте ненужные обработчики**: удаляйте обработчики событий, когда они больше не нужны

[Подробнее о производительности →](/guides/performance/)

## Модель безопасности

Архитектура Wails безопасна по умолчанию:

```d2
direction: down

Frontend: Фронтенд (недоверенная среда) {
  shape: rectangle
  style.fill: "#EF4444"
}

Bridge: Мост Wails (проверка) {
  shape: diamond
  style.fill: "#F59E0B"
}

Backend: Бэкенд (доверенная среда) {
  shape: rectangle
  style.fill: "#10B981"
}

Frontend -> Bridge: Вызвать метод
Bridge -> Bridge: "Проверка:\n- Метод существует?\n- Типы указаны верно?\n- Доступ разрешён?"
Bridge -> Backend: Выполнить, если проверка пройдена
Backend -> Bridge: Вернуть результат
Bridge -> Frontend: Отправить ответ
```

**Механизмы безопасности:**

1. **Белый список методов**
  - Можно вызывать только экспортируемые методы
  - Закрытые методы недоступны
  - Требуется явная регистрация сервисов


2. **Проверка типов**
  - Аргументы проверяются на соответствие типам Go
  - Значения недопустимых типов отклоняются
  - Предотвращает атаки с внедрением кода


3. **Без eval()**
  - Фронтенд не может выполнять произвольный код Go
  - Можно вызывать только предопределённые методы
  - Нет динамического выполнения кода


4. **Изоляция контекста**
  - У каждого окна есть собственный контекст
  - Сервисы могут проверять контекст вызывающей стороны
  - Для каждого окна можно задать отдельные разрешения


**Рекомендации:**

- **Проверяйте пользовательский ввод** в Go (не доверяйте фронтенду)
- **Используйте контекст** для аутентификации и авторизации
- **Очищайте пути к файлам** перед операциями с файлами
- **Ограничивайте частоту** ресурсоёмких операций

[Подробнее о безопасности →](/guides/security/)

## Следующие шаги

**Жизненный цикл приложения** — изучите запуск, завершение работы и перехватчики жизненного цикла [Подробнее →](/concepts/lifecycle/)

**Мост между Go и фронтендом** — подробно изучите работу моста [Подробнее →](/concepts/bridge/)

**Система сборки** — узнайте, как Wails собирает ваше приложение [Подробнее →](/concepts/build-system/)

**Приступайте к разработке** — примените полученные знания на практике в учебном руководстве [Учебные руководства →](/tutorials/03-notes-vanilla/)

---

**Есть вопросы об архитектуре?** Задайте их в [Discord](https://discord.gg/JDdSxwjhGf) или обратитесь к [справочнику по API](/reference/overview/).
