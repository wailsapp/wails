---
title: "Как работает Wails"
description: "Архитектура Wails и способы достижения нативной производительности"
slug: "concepts/architecture"
sourcePath: "concepts/architecture.md"
---

Wails — это фреймворк для создания настольных приложений, в котором **Go используется для бэкенда**, а **веб-технологии — для фронтенда**. Но, в отличие от Electron, Wails не включает браузер в состав приложения, а использует **нативный WebView операционной системы**.

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
    frontend."Render Site via OS-native WebView": Отрисовывает сайт в нативном WebView ОС
  }
  Regular Communication: Обычный обмен данными {
    shape: sequence_diagram
    frontend."Make API-style call": Выполнить вызов в стиле API
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

**Ключевые отличия от Electron:**

| Характеристика | Wails | Electron |
| --- | --- | --- |
| **Браузер** | WebView, предоставляемый ОС | Встроенный Chromium (~100 МБ) |
| **Бэкенд** | Go (компилируемый) | Node.js (интерпретируемый) |
| **Обмен данными** | Мост в памяти | IPC (межпроцессное взаимодействие) |
| **Размер пакета приложения** | ~15 МБ | ~150 МБ |
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
- Малый размер и высокая производительность

@end

**Почему это важно:**

- **Нет встроенного браузера** → Меньший размер приложения
- **Нативная технология ОС** → Улучшенная интеграция и производительность
- **Автоматические обновления** → Исправления безопасности через обновления ОС
- **Привычный рендеринг** → Такой же, как в системном браузере

### 2. Мост Wails

Мост — центральный компонент Wails, обеспечивающий **прямой обмен данными** между Go и JavaScript.

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

  Decoder: Декодер JSON {
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
Bridge.Decoder -> Frontend: "5. Декодировать в JS\nPromise завершается"
```

**Как это работает:**

1. **Фронтенд вызывает метод Go** (через автоматически сгенерированную привязку)
2. **Мост кодирует вызов** в JSON (имя метода и аргументы)
3. **Маршрутизатор находит метод Go** в зарегистрированных сервисах
4. **Метод Go выполняется** и возвращает значение
5. **Мост декодирует результат** и отправляет его обратно во фронтенд
6. **Promise завершается** в JavaScript с полученным результатом

**Характеристики производительности:**

- **В памяти**: без сетевых накладных расходов и HTTP
- **Без копирования** там, где это возможно (для больших объёмов данных)
- **Асинхронность по умолчанию**: обе стороны работают без блокировки
- **Безопасность типов**: определения TypeScript генерируются автоматически

### 3. Система сервисов

Сервисы — рекомендуемый способ предоставить фронтенду функциональность Go.

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

- При запуске Wails **анализирует вашу структуру**
- **Экспортируемые методы** становятся доступны для вызова из фронтенда
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
- **Автоматическое обнаружение**: не нужно регистрировать методы вручную
- **Упорядоченность**: связанная функциональность объединена в группы
- **Удобство тестирования**: сервисы — это обычные структуры Go

[Подробнее о сервисах →](/features/bindings/services/)

### 4. Система событий

События обеспечивают **обмен данными между компонентами по модели публикации и подписки**.

```d2
direction: left

Wails Event System: Система событий Wails {
  shape: sequence_diagram

  window1: Окно 1
  window2: Окно 2
  backend: Бэкенд Go

  Event Driver: Диспетчер событий {
    shape: sequence_diagram
    window1."Subscribe to 'data-updated' events": "Подписка на события 'data-updated'"
    window2."Subscribe to 'data-updated' events": "Подписка на события 'data-updated'"
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

**Сценарии использования:**

- **Обмен данными между окнами**: одно окно уведомляет другие
- **Фоновые задачи**: сервис Go уведомляет пользовательский интерфейс о ходе выполнения
- **Синхронизация состояния**: поддержание согласованного состояния нескольких окон
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

Понимание жизненного цикла помогает определить, когда следует инициализировать и освобождать ресурсы.

```d2
direction: down

Start: Запуск приложения {
  shape: oval
  style.fill: "#10B981"
}

Init: Инициализация {
  Create: Создание приложения {
    shape: rectangle
  }

  Register: Регистрация сервисов {
    shape: rectangle
  }

  Setup: Настройка окон и меню {
    shape: rectangle
  }
}

Run: Цикл событий {
  Events: Обработка событий {
    shape: rectangle
  }

  Messages: Обработка сообщений {
    shape: rectangle
  }

  Render: Обновление интерфейса {
    shape: rectangle
  }
}

Shutdown: Завершение работы {
  Cleanup: Освобождение ресурсов {
    shape: rectangle
  }

  Save: Сохранение состояния {
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
Run.Events -> Shutdown.Cleanup: Сигнал завершения
Shutdown.Cleanup -> Shutdown.Save
Shutdown.Save -> End
```

**Перехватчики жизненного цикла:**

```go
app := application.New(application.Options{
    Name: "My App",

    // Cleanly intercept quit requests (e.g. unsaved changes).
    ShouldQuit: func() bool { return true },

    // Called when the app is confirmed to be quitting — save state, close connections, etc.
    OnShutdown: func() {},
})
```

У `application.Options` нет поля `OnStartup`. Действия при запуске следует выполнять в `ServiceStartup(ctx, options)` сервиса, в обратном вызове, зарегистрированном через `app.Event.OnApplicationEvent(events.Common.ApplicationStarted, ...)`, или просто перед `app.Run()`.

[Подробнее о жизненном цикле →](/concepts/lifecycle/)

## Процесс сборки

Как Wails собирает ваше приложение:

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
  AnalyseGo: Анализ кода Go {
    shape: rectangle
  }

  GenerateBindings: Создание привязок {
    shape: rectangle
  }

  BuildFrontend: Сборка фронтенда {
    shape: rectangle
  }

  CompileGo: Компиляция Go {
    shape: rectangle
  }

  Embed: Встраивание ресурсов {
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
Build.AnalyseGo -> Build.GenerateBindings: Извлечение типов
Build.GenerateBindings -> Source.Frontend: Привязки TypeScript
Source.Frontend -> Build.BuildFrontend: Компиляция (Vite/webpack)
Build.BuildFrontend -> Build.Embed: Упакованные ресурсы
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
  - Запуск сборщика модулей (Vite, webpack и т. д.)
  - Минификация и оптимизация
  - Вывод в `frontend/dist/`


4. **Компиляция Go**
  - Компиляция с оптимизациями (`-ldflags="-s -w"`)
  - Добавление метаданных сборки
  - Компиляция для конкретной платформы


5. **Встраивание ресурсов**
  - Встраивание файлов фронтенда в исполняемый файл Go
  - Сжатие ресурсов
  - Создание единого исполняемого файла


**Результат:** единый нативный исполняемый файл со всеми встроенными компонентами.

[Подробнее о сборке →](/guides/build/building/)

## Разработка и промышленная эксплуатация

В режимах разработки и промышленной эксплуатации Wails работает по-разному:

@tabs{sync-key="mode"}
[Разработка (wails3 dev)]
**Характеристики:**

- **Горячая перезагрузка**: изменения во фронтенде применяются мгновенно
- **Карты исходного кода**: отладка по исходному коду
- **DevTools**: доступны инструменты разработчика браузера
- **Ведение журнала**: включено подробное журналирование
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

WailsApp -> DevServer: Проксирование запросов
DevServer -> WebView: Обслуживание с HMR
WebView -> WailsApp: Вызов методов Go
```

**Преимущества:**

- Мгновенная обратная связь при внесении изменений
- Полноценные возможности отладки
- Более быстрые итерации разработки

[Продакшен (wails3 build)]
**Характеристики:**

- **Встроенные ресурсы**: фронтенд встроен в бинарный файл
- **Оптимизация**: минификация и сжатие
- **Без DevTools**: по умолчанию отключены
- **Минимальное журналирование**: только ошибки
- **Один файл**: всё содержится в одном исполняемом файле

**Как это работает:**

```d2
direction: right

Binary: "Единый бинарный файл\n(myapp.exe)" {
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

Binary.Assets -> WebView: Обслуживание из памяти
WebView -> Binary.GoCode: Вызов методов Go
```

**Преимущества:**

- Распространение в виде одного файла
- Меньший размер благодаря минификации
- Более высокая производительность
- Нет внешних зависимостей

@end

## Модель памяти

Понимание принципов использования памяти помогает создавать эффективные приложения.

**Области памяти:**

1. **Куча Go**
  - Ваши сервисы и состояние приложения
  - Управляется сборщиком мусора Go
  - Обычно 5-10 МБ для простых приложений


2. **Память WebView**
  - DOM, куча JavaScript и CSS
  - Управляется движком WebView
  - Обычно 10-20 МБ для простых приложений


3. **Память моста**
  - Буферы сообщений для обмена данными
  - Минимальные накладные расходы (<1 МБ)
  - По возможности — передача больших объёмов данных без копирования


**Советы по оптимизации:**

- **Избегайте передачи больших объёмов данных**: передавайте идентификаторы, а подробные данные загружайте по запросу
- **Используйте события для обновлений**: не выполняйте опрос из фронтенда
- **Передавайте большие файлы потоково**: не загружайте их в память целиком
- **Удаляйте ненужные обработчики**: удаляйте обработчики событий, когда они больше не нужны

[Подробнее о производительности →](/guides/performance/)

## Модель безопасности

Архитектура Wails по умолчанию обеспечивает безопасность:

```d2
direction: down

Frontend: Фронтенд (недоверенный) {
  shape: rectangle
  style.fill: "#EF4444"
}

Bridge: Мост Wails (проверка) {
  shape: diamond
  style.fill: "#F59E0B"
}

Backend: Бэкенд (доверенный) {
  shape: rectangle
  style.fill: "#10B981"
}

Frontend -> Bridge: Вызов метода
Bridge -> Bridge: "Проверка:\n- Метод существует?\n- Типы верны?\n- Доступ разрешён?"
Bridge -> Backend: Выполнение после успешной проверки
Backend -> Bridge: Возврат результата
Bridge -> Frontend: Отправка ответа
```

**Функции безопасности:**

1. **Белый список методов**
  - Можно вызывать только экспортируемые методы
  - Закрытые методы недоступны
  - Требуется явная регистрация сервисов


2. **Проверка типов**
  - Аргументы проверяются на соответствие типам Go
  - Некорректные типы отклоняются
  - Предотвращает атаки с внедрением кода


3. **Без eval()**
  - Фронтенд не может выполнять произвольный код Go
  - Можно вызывать только предопределённые методы
  - Динамическое выполнение кода отсутствует


4. **Изоляция контекстов**
  - У каждого окна есть собственный контекст
  - Сервисы могут проверять контекст вызывающей стороны
  - Можно задавать разрешения для каждого окна


**Рекомендации:**

- **Проверяйте пользовательский ввод** в Go (не доверяйте фронтенду)
- **Используйте контекст** для аутентификации и авторизации
- **Проверяйте и очищайте пути к файлам** перед файловыми операциями
- **Ограничивайте частоту** ресурсоёмких операций

[Подробнее о безопасности →](/guides/security/)

## Дальнейшие шаги

**Жизненный цикл приложения** — узнайте о запуске, завершении работы и обработчиках жизненного цикла [Подробнее →](/concepts/lifecycle/)

**Мост между Go и фронтендом** — подробно изучите работу моста [Подробнее →](/concepts/bridge/)

**Система сборки** — узнайте, как Wails собирает ваше приложение [Подробнее →](/concepts/build-system/)

**Начните разработку** — примените полученные знания на практике, следуя руководству [Руководства →](/tutorials/03-notes-vanilla/)

---

**Есть вопросы об архитектуре?** Задайте их в [Discord](https://discord.gg/JDdSxwjhGf) или обратитесь к [справочнику по API](/reference/overview/).
