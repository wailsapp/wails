---
title: "Список задач"
description: "Создание полноценного приложения со списком задач и CRUD-операциями"
slug: "tutorials/02-todo-vanilla"
sourcePath: "tutorials/02-todo-vanilla.md"
---

В этом руководстве вы создадите полнофункциональное приложение со списком задач. Это следующий шаг после руководства по сервису QR-кодов: вы научитесь управлять состоянием, обрабатывать несколько операций и создавать проработанный пользовательский интерфейс.

**Что вы создадите:**

- Полноценное приложение со списком задач, в котором можно добавлять, завершать и удалять задачи
- Потокобезопасное управление состоянием (важно для настольных приложений)
- Современный интерфейс в стиле glassmorphism
- Всё это — на чистом JavaScript, без фреймворков

**Чему вы научитесь:**

- CRUD-операции (создание, чтение, обновление и удаление)
- Безопасное управление изменяемым состоянием в Go
- Обработка и проверка пользовательского ввода
- Создание адаптивных интерфейсов, которые выглядят и работают как нативные

![Приложение со списком задач](/assets/todo-app.png)

**Время выполнения:** 20 минут

## Создание проекта

@steps
### Генерация проекта
Сначала создайте новый проект Wails. Мы будем использовать стандартный шаблон на чистом JavaScript, который послужит удобной отправной точкой:

```bash
wails3 init -n todo-app -t vanilla-js
cd todo-app
```

Будет создан новый проект с базовой структурой: бэкенд на Go — в корневом каталоге, код фронтенда — в каталоге `frontend/`.

### Создание сервиса TODO
Сервис TODO будет управлять состоянием приложения и предоставлять методы для CRUD-операций. В отличие от веб-сервера, где каждый запрос обрабатывается изолированно, в настольном приложении несколько операций могут выполняться параллельно, поэтому необходимо потокобезопасное управление состоянием.

Удалите `greetservice.go` и создайте новый файл `todoservice.go`:

```go {title="todoservice.go"}
package main

import (
    "errors"
    "sync"
)

type Todo struct {
    ID        int    `json:"id"`
    Title     string `json:"title"`
    Completed bool   `json:"completed"`
}

type TodoService struct {
    todos  []Todo
    nextID int
    mu     sync.RWMutex
}

func NewTodoService() *TodoService {
    return &TodoService{
        todos:  []Todo{},
        nextID: 1,
    }
}

func (t *TodoService) GetAll() []Todo {
    t.mu.RLock()
    defer t.mu.RUnlock()
    return t.todos
}

func (t *TodoService) Add(title string) (*Todo, error) {
    if title == "" {
        return nil, errors.New("title cannot be empty")
    }

    t.mu.Lock()
    defer t.mu.Unlock()

    todo := Todo{
        ID:        t.nextID,
        Title:     title,
        Completed: false,
    }
    t.todos = append(t.todos, todo)
    t.nextID++

    return &todo, nil
}

func (t *TodoService) Toggle(id int) error {
    t.mu.Lock()
    defer t.mu.Unlock()

    for i := range t.todos {
        if t.todos[i].ID == id {
            t.todos[i].Completed = !t.todos[i].Completed
            return nil
        }
    }
    return errors.New("todo not found")
}

func (t *TodoService) Delete(id int) error {
    t.mu.Lock()
    defer t.mu.Unlock()

    for i, todo := range t.todos {
        if todo.ID == id {
            t.todos = append(t.todos[:i], t.todos[i+1:]...)
            return nil
        }
    }
    return errors.New("todo not found")
}
```

**Что здесь происходит:**

**Структура `Todo`:**

- Определяет структуру данных с полями ID, Title и Completed
- Теги `json:` указывают Go, как преобразовать эту структуру в JSON для фронтенда
- Каждое поле экспортируется (его имя начинается с прописной буквы), чтобы генератор привязок мог его обнаружить

**Структура `TodoService`:**

- `todos []Todo` — срез, содержащий все задачи
- `nextID int` — хранит следующий назначаемый ID (имитирует автоинкремент)
- `mu sync.RWMutex` — мьютекс чтения-записи для потокобезопасного доступа

**Потокобезопасность с помощью `sync.RWMutex`:**

- Пользовательский интерфейс настольного приложения может параллельно выполнять несколько операций
- `RLock()` разрешает одновременный доступ нескольким читателям (например, при нескольких вызовах `GetAll`)
- `Lock()` предоставляет эксклюзивный доступ для записи (например, для `Add`, `Toggle` и `Delete`)
- `defer` гарантирует снятие блокировки, даже если функция завершится досрочно или возникнет паника

**Методы:**

- `GetAll()` — возвращает все задачи (использует блокировку чтения, поскольку данные не изменяются)
- `Add(title)` — создаёт новую задачу, проверяет входные данные и увеличивает ID
- `Toggle(id)` — переключает статус завершения задачи
- `Delete(id)` — удаляет задачу из среза

**Обработка ошибок:**

- В соответствии с соглашениями Go мы возвращаем `error` последним значением
- Пустые названия отклоняются
- Операции с несуществующими задачами возвращают ошибки
- На фронтенде эти ошибки преобразуются в исключения JavaScript

### Обновление main.go
Зарегистрируйте сервис TODO в приложении Wails. Найдите раздел `Services` в файле `main.go` и замените GreetService на TodoService:

```go {title="main.go" highlight="5"}
Services: []application.Service{
    application.NewService(NewTodoService()),
},
```

**Что здесь происходит:**

- Мы удаляем стандартный GreetService и вместо него добавляем TodoService
- `application.NewService()` создаёт обёртку для сервиса, чтобы Wails мог им управлять
- Wails автоматически сгенерирует привязки JavaScript для всех публичных методов этого сервиса

### Создание пользовательского интерфейса фронтенда
Теперь создадим фронтенд. Здесь мы будем вызывать методы Go и отображать пользовательский интерфейс. Чтобы упростить пример и наглядно показать непосредственную работу привязок, мы используем чистый JavaScript.

Замените содержимое `frontend/src/main.js`:

```javascript {title="frontend/src/main.js"}
import {TodoService} from "../bindings/changeme";

async function loadTodos() {
    const todos = await TodoService.GetAll();
    const list = document.getElementById('todo-list');

    list.innerHTML = todos.map(todo => `
        <div class="todo ${todo.completed ? 'completed' : ''}">
            <input type="checkbox"
                   ${todo.completed ? 'checked' : ''}
                   onchange="toggleTodo(${todo.id})">
            <span>${todo.title}</span>
            <button onclick="deleteTodo(${todo.id})">Delete</button>
        </div>
    `).join('');
}

window.addTodo = async () => {
    const input = document.getElementById('todo-input');
    const title = input.value.trim();

    if (title) {
        await TodoService.Add(title);
        input.value = '';
        await loadTodos();
    }
}

window.toggleTodo = async (id) => {
    await TodoService.Toggle(id);
    await loadTodos();
}

window.deleteTodo = async (id) => {
    await TodoService.Delete(id);
    await loadTodos();
}

// Load todos on startup
loadTodos();
```

**Что здесь происходит:**

**Импорт привязок:**

- `import {TodoService} from "../bindings/changeme"` — импортирует автоматически сгенерированные привязки Go
- Примечание: вместо `changeme` будет фактическое имя вашего модуля из `go.mod`

**Функция `loadTodos()`:**

- Вызывает `TodoService.GetAll()`, чтобы получить из Go все задачи
- Формирует HTML для каждой задачи с помощью шаблонных литералов
- Динамически добавляет и удаляет класс `completed` для применения стилей
- Использует атрибуты `onclick`, чтобы связать кнопки с нашими функциями
- Объединяет весь HTML и вставляет его в DOM

**Функции CRUD:**

- `addTodo()` — проверяет входные данные, вызывает метод Go `Add` и обновляет список
- `toggleTodo(id)` — вызывает метод Go `Toggle` и обновляет список
- `deleteTodo(id)` — вызывает метод Go `Delete` и обновляет список
- Все функции асинхронные, поскольку вызовы Go возвращают объекты Promise

**Зачем добавлять функции в window:**

- `window.addTodo = ...` делает функции доступными из атрибутов HTML `onclick`
- Это простой подход для чистого JavaScript (во фреймворках это реализуется иначе)
- В рабочей версии приложения вместо этого можно использовать полноценное делегирование событий

**Схема обновления:**

- После каждого изменения (добавления, переключения статуса или удаления) мы снова вызываем `loadTodos()`
- Это гарантирует, что пользовательский интерфейс остаётся синхронизированным с состоянием Go
- Альтернативный вариант: методы Go могут возвращать новое состояние, чтобы избежать второго вызова

### Обновление HTML
HTML задаёт структуру нашего приложения TODO. Она минималистична и семантична — вся магия происходит в JavaScript и CSS.

Замените содержимое `frontend/index.html`:

```html {title="frontend/index.html"}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
    <title>TODO App</title>
    <link rel="stylesheet" href="./style.css"/>
</head>
<body>
    <div class="container">
        <h1>My TODOs</h1>
        <div class="card">
            <div class="input-box">
                <input type="text"
                       id="todo-input"
                       class="input"
                       placeholder="Add a new todo..."
                       onkeypress="if(event.key==='Enter') addTodo()">
                <button class="btn" onclick="addTodo()">Add</button>
            </div>
            <div id="todo-list"></div>
        </div>
    </div>
    <script type="module" src="./src/main.js"></script>
</body>
</html>
```

**Что здесь происходит:**

**Структура:**

- `container` — центрирует приложение и ограничивает его ширину
- `card` — основная белая карточка, содержащая все элементы
- `input-box` — flex-контейнер для поля ввода и кнопки «Добавить»
- `todo-list` — сюда JavaScript будет добавлять отдельные задачи

**Обработка событий:**

- `onkeypress="if(event.key==='Enter') addTodo()"` — добавляет задачу при нажатии Enter
- `onclick="addTodo()"` — добавляет задачу при нажатии кнопки
- Встроенные обработчики событий хорошо подходят для простых приложений на чистом JavaScript

**Модульный скрипт:**

- `<script type="module">` позволяет использовать импорты ES6
- Наш файл `main.js` может импортировать привязки и использовать современный JavaScript

### Оформите приложение
CSS создаёт современный дизайн в стиле стекломорфизма с плавными переходами. Мы стремимся придать приложению завершённый вид, чтобы им было приятно пользоваться.

Замените содержимое `frontend/public/style.css`:

```css {title="frontend/public/style.css"}
:root {
    font-family: "Inter", -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto",
    "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue",
    sans-serif;
    font-size: 16px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: rgba(255, 255, 255, 0.87);
}

body {
    margin: 0;
    display: flex;
    place-items: center;
    justify-content: center;
    min-height: 100vh;
}

.container {
    width: 100%;
    max-width: 600px;
    padding: 20px;
}

h1 {
    text-align: center;
    color: white;
    font-size: 2.5em;
    font-weight: 300;
    margin: 0 0 30px 0;
    text-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
}

.card {
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(10px);
    border-radius: 16px;
    padding: 30px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.input-box {
    display: flex;
    gap: 10px;
    margin-bottom: 25px;
}

.input {
    flex: 1;
    border: 2px solid #e0e0e0;
    border-radius: 12px;
    height: 50px;
    padding: 0 20px;
    font-size: 16px;
    transition: all 0.3s ease;
}

.input:focus {
    border-color: #667eea;
    outline: none;
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.btn {
    height: 50px;
    padding: 0 30px;
    border: none;
    border-radius: 12px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s ease;
    box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.btn:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(102, 126, 234, 0.6);
}

#todo-list {
    display: flex;
    flex-direction: column;
    gap: 10px;
}

.todo {
    display: flex;
    align-items: center;
    padding: 18px 20px;
    background: white;
    border: 2px solid #f0f0f0;
    border-radius: 12px;
    transition: all 0.3s ease;
    gap: 15px;
}

.todo:hover {
    border-color: #667eea;
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.15);
    transform: translateX(4px);
}

.todo.completed {
    opacity: 0.6;
}

.todo.completed span {
    text-decoration: line-through;
    color: #999;
}

.todo input[type="checkbox"] {
    width: 24px;
    height: 24px;
    cursor: pointer;
    appearance: none;
    -webkit-appearance: none;
    border: 2px solid #667eea;
    border-radius: 6px;
    position: relative;
    transition: all 0.3s ease;
    flex-shrink: 0;
}

.todo input[type="checkbox"]:hover {
    background: rgba(102, 126, 234, 0.1);
}

.todo input[type="checkbox"]:checked {
    background: #667eea;
    border-color: #667eea;
}

.todo input[type="checkbox"]:checked::after {
    content: '✓';
    position: absolute;
    color: white;
    font-size: 16px;
    font-weight: bold;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
}

.todo span {
    flex: 1;
    font-size: 16px;
    color: #333;
}

.todo button {
    padding: 8px 16px;
    background: #ff4757;
    color: white;
    border: none;
    border-radius: 8px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s ease;
    opacity: 0;
    flex-shrink: 0;
}

.todo:hover button {
    opacity: 1;
}

.todo button:hover {
    background: #ee5a6f;
    transform: scale(1.05);
}

#todo-list:empty::before {
    content: "No todos yet. Add one above!";
    display: block;
    text-align: center;
    padding: 40px 20px;
    color: #999;
}
```

**Что здесь происходит:**

**Дизайн в стиле стекломорфизма:**

- `background: linear-gradient(135deg, #667eea 0%, #764ba2 100%)` — фиолетовый градиентный фон
- `backdrop-filter: blur(10px)` — создаёт на карточке эффект матового стекла
- `rgba(255, 255, 255, 0.95)` — полупрозрачный белый цвет для эффекта стекла

**Пользовательское оформление флажка:**

- `appearance: none` удаляет стандартное оформление флажка в браузере
- С помощью `::after` мы создаём собственный скруглённый квадрат с галочкой
- Когда выполняется условие `checked`, появляется галочка в виде символа Unicode ✓

**Взаимодействие при наведении:**

- При наведении задачи сдвигаются вправо (`transform: translateX(4px)`)
- Кнопка удаления скрыта до наведения (`opacity: 0` → `opacity: 1`)
- При наведении кнопки немного увеличиваются, создавая тактильный отклик

**Пустое состояние:**

- `#todo-list:empty::before` показывает сообщение, когда задач нет
- Решение реализовано только на CSS — JavaScript не требуется

### Запустите приложение
Посмотрим на приложение в действии! Запустите сервер разработки:

```bash
wails3 dev
```

Приложение скомпилируется и откроется. Попробуйте его в работе:

- Введите задачу и нажмите Enter или кнопку «Добавить»
- Установите флажок, чтобы отметить задачу как выполненную
- Наведите указатель на задачу, чтобы появилась кнопка удаления
- Обратите внимание, что интерфейс обновляется мгновенно — так работает наш шаблон обновления

**Что происходит:**

- Wails автоматически сгенерировал привязки для методов вашего TodoService
- Режим разработки поддерживает горячую перезагрузку — попробуйте изменить CSS и посмотрите, как обновится приложение
- Ваш код Go выполняется нативно — преобразование или интерпретация не требуются

@end

## Как это работает

### Потокобезопасное управление состоянием

`sync.RWMutex` обеспечивает безопасный конкурентный доступ:

```go
func (t *TodoService) GetAll() []Todo {
    t.mu.RLock()  // Read lock - multiple readers allowed
    defer t.mu.RUnlock()
    return t.todos
}

func (t *TodoService) Add(title string) (*Todo, error) {
    t.mu.Lock()  // Write lock - exclusive access
    defer t.mu.Unlock()
    // ... mutations
}
```

**Почему это важно:**

- Несколько вызовов из фронтенда могут выполняться конкурентно
- Операции чтения не блокируют друг друга
- Операции записи получают исключительный доступ
- `defer` гарантирует, что блокировки всегда будут сняты

### Обработка ошибок

При недопустимых операциях сервис возвращает ошибки:

```go
func (t *TodoService) Add(title string) (*Todo, error) {
    if title == "" {
        return nil, errors.New("title cannot be empty")
    }
    // ...
}
```

Во фронтенде эти ошибки можно перехватывать:

```javascript
try {
    await TodoService.Add(title);
} catch (err) {
    alert('Error: ' + err);
}
```

### Синхронизация состояния

После каждого изменения мы заново загружаем весь список:

```javascript
window.addTodo = async () => {
    await TodoService.Add(title);  // Mutation
    await loadTodos();             // Refresh
}
```

**Альтернативный подход:** возвращайте обновлённый список из каждого метода, чтобы избежать второго вызова.

## Улучшения

### Добавьте статистику

Добавьте следующий код в `todoservice.go`:

```go
type TodoStats struct {
    Total     int `json:"total"`
    Completed int `json:"completed"`
    Active    int `json:"active"`
}

func (t *TodoService) GetStats() TodoStats {
    t.mu.RLock()
    defer t.mu.RUnlock()

    stats := TodoStats{
        Total: len(t.todos),
    }

    for _, todo := range t.todos {
        if todo.Completed {
            stats.Completed++
        } else {
            stats.Active++
        }
    }

    return stats
}
```

Отобразите статистику во фронтенде:

```javascript
async function loadTodos() {
    const [todos, stats] = await Promise.all([
        TodoService.GetAll(),
        TodoService.GetStats()
    ]);

    // Display stats
    document.getElementById('stats').textContent =
        `${stats.active} active, ${stats.completed} completed`;

    // ... render todos
}
```

### Добавьте «Очистить выполненные»

```go
func (t *TodoService) ClearCompleted() int {
    t.mu.Lock()
    defer t.mu.Unlock()

    removed := 0
    newTodos := []Todo{}

    for _, todo := range t.todos {
        if !todo.Completed {
            newTodos = append(newTodos, todo)
        } else {
            removed++
        }
    }

    t.todos = newTodos
    return removed
}
```

### Добавьте постоянное хранение данных

В рабочем приложении добавьте постоянное хранение данных с помощью решения, которое лучше всего подходит для вашего приложения, например SQLite, PostgreSQL или размещённого сервиса.

## Соберите приложение для рабочей среды

Когда будете готовы распространять приложение TODO, соберите его для рабочей среды:

```bash
wails3 build
```

При этом в `bin/` создаётся оптимизированный нативный исполняемый файл:

- Код Go компилируется с оптимизациями
- Фронтенд собирается для рабочей среды
- Все компоненты упаковываются в единый исполняемый файл
- Размер полученного приложения обычно составляет 10-20 МБ (для сравнения: у Electron — 150 МБ и более)

Исполняемый файл можно запускать напрямую — среда выполнения не нужна и запускать серверы не требуется. Это настоящее нативное приложение.

## Что вы создали

Вы только что создали полноценное приложение TODO со следующими возможностями:

**Полная реализация CRUD:**

- Создали сервис с операциями создания, чтения, обновления и удаления
- Добавили проверку входных данных и обработку ошибок
- Узнали, как ошибки Go преобразуются в исключения JavaScript

**Потокобезопасное управление состоянием:**

- Использовали `sync.RWMutex` для безопасной обработки конкурентного доступа
- Разобрались в различиях между блокировками чтения (RLock) и записи (Lock)
- Увидели, как `defer` предотвращает взаимные блокировки, гарантируя освобождение ресурсов

**Современный и качественно оформленный интерфейс:**

- Создали интерфейс в стиле стекломорфизма с градиентами и эффектами размытия
- Создали флажки с пользовательским оформлением без использования фреймворков
- Добавили эффекты наведения и переходы, чтобы интерфейс воспринимался как нативный
- Реализовано пустое состояние с помощью чистого CSS

**Основы Wails:**

- Регистрация сервисов и автоматическое создание привязок
- Вызов методов Go из JavaScript с помощью async/await
- Синхронизация состояния между Go и фронтендом
- Сборка и упаковка нативного настольного приложения

## Следующие шаги

Теперь, когда вы разобрались с CRUD-операциями и управлением состоянием, попробуйте следующее:

- **Добавьте постоянное хранение данных:** сохраняйте задачи между перезапусками приложения с помощью SQLite или другого решения для хранения данных
- **Добавьте новые возможности:** фильтрацию (все/активные/завершённые), редактирование существующих задач и массовые операции
- **Изучите руководство по приложению Notes:** узнайте, как работают файловые операции в приложении [Notes](/tutorials/03-notes-vanilla/)
- **Создайте настоящее приложение:** используйте эти концепции, чтобы создать собственное приложение!
