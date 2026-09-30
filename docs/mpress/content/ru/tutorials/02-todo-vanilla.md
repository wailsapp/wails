---
title: "Список задач"
description: "Создайте полноценное приложение со списком задач и CRUD-операциями"
slug: "tutorials/02-todo-vanilla"
sourcePath: "tutorials/02-todo-vanilla.md"
---

В этом руководстве вы создадите полнофункциональное приложение со списком задач. Это следующий шаг после руководства по сервису QR-кодов: вы научитесь управлять состоянием, выполнять различные операции и создавать проработанный пользовательский интерфейс.

**Что вы создадите:**

- Полноценное приложение со списком задач, в котором можно добавлять, отмечать выполненными и удалять задачи
- Потокобезопасное управление состоянием (это важно для настольных приложений)
- Современный пользовательский интерфейс в стиле glassmorphism
- Всё это — на чистом JavaScript, без фреймворков

**Чему вы научитесь:**

- CRUD-операции (создание, чтение, обновление и удаление)
- Безопасное управление изменяемым состоянием в Go
- Обработка и проверка пользовательского ввода
- Создание адаптивных интерфейсов, воспринимаемых как нативные

![Приложение со списком задач](/assets/todo-app.png)

**Время выполнения:** 20 минут

## Создание проекта

@steps
### Создайте проект
Сначала создайте новый проект Wails. Мы воспользуемся стандартным шаблоном на чистом JavaScript, который обеспечит нам удобную отправную точку:

```bash
wails3 init -n todo-app
cd todo-app
```

В результате будет создан новый проект с базовой структурой: серверная часть на Go в корневом каталоге и код клиентской части в каталоге `frontend/`.

### Создайте сервис TODO
Сервис TODO будет управлять состоянием приложения и предоставлять методы для CRUD-операций. В отличие от веб-сервера, где каждый запрос обрабатывается изолированно, в настольном приложении несколько операций могут выполняться одновременно, поэтому необходимо потокобезопасное управление состоянием.

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
- Теги `json:` указывают Go, как преобразовать эту структуру в JSON для клиентской части
- Каждое поле экспортируется (его имя начинается с заглавной буквы), чтобы генератор привязок мог его обнаружить

**Структура `TodoService`:**

- `todos []Todo` — срез, содержащий все наши задачи
- `nextID int` — хранит следующий назначаемый ID (имитирует автоинкремент)
- `mu sync.RWMutex` — мьютекс чтения и записи для потокобезопасного доступа

**Потокобезопасность с помощью `sync.RWMutex`:**

- Пользовательский интерфейс настольного приложения может одновременно выполнять несколько операций
- `RLock()` допускает одновременный доступ нескольких читателей (например, несколько вызовов `GetAll`)
- `Lock()` предоставляет исключительный доступ для записи (например, `Add`, `Toggle`, `Delete`)
- `defer` гарантирует снятие блокировки, даже если функция завершится досрочно или возникнет паника

**Методы:**

- `GetAll()` — возвращает все задачи (использует блокировку чтения, поскольку данные не изменяются)
- `Add(title)` — создаёт новую задачу, проверяет входные данные и увеличивает ID
- `Toggle(id)` — переключает статус выполнения задачи
- `Delete(id)` — удаляет задачу из среза

**Обработка ошибок:**

- Следуя соглашениям Go, последним значением мы возвращаем `error`
- Пустые названия отклоняются
- Операции с несуществующими задачами возвращают ошибки
- В клиентской части эти ошибки становятся исключениями JavaScript

### Обновите main.go
Зарегистрируйте сервис TODO в приложении Wails. Найдите раздел `Services` в файле `main.go` и замените GreetService нашим TodoService:

```go {title="main.go" highlight="5"}
Services: []application.Service{
    application.NewService(NewTodoService()),
},
```

**Что здесь происходит:**

- Мы удаляем стандартный GreetService и вместо него добавляем наш TodoService
- `application.NewService()` оборачивает наш сервис, чтобы Wails мог им управлять
- Wails автоматически создаст привязки JavaScript для всех публичных методов этого сервиса

### Создайте пользовательский интерфейс клиентской части
Теперь создадим клиентскую часть. Здесь мы будем вызывать методы Go и отображать пользовательский интерфейс. Чтобы упростить пример и наглядно показать работу с привязками, мы используем чистый JavaScript.

Замените `frontend/src/main.js`:

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

- `import {TodoService} from "../bindings/changeme"` — импортирует автоматически созданные привязки Go
- Примечание: `changeme` будет фактическим именем вашего модуля из `go.mod`

**Функция `loadTodos()`:**

- Вызывает `TodoService.GetAll()`, чтобы получить из Go все задачи
- Формирует HTML для каждой задачи с помощью шаблонных строк
- Динамически добавляет или удаляет класс `completed` для применения стилей
- Использует атрибуты `onclick`, чтобы связать кнопки с нашими функциями
- Объединяет весь HTML и вставляет его в DOM

**CRUD-функции:**

- `addTodo()` — проверяет входные данные, вызывает метод Go `Add` и обновляет список
- `toggleTodo(id)` — вызывает метод Go `Toggle` и обновляет список
- `deleteTodo(id)` — вызывает метод Go `Delete` и обновляет список
- Все функции являются асинхронными, поскольку вызовы Go возвращают Promise

**Зачем добавлять функции в window:**

- `window.addTodo = ...` делает функции доступными из атрибутов HTML `onclick`
- Это простой шаблон для обычного JavaScript (во фреймворках это реализуется иначе)
- В рабочей версии вместо этого можно использовать полноценное делегирование событий

**Шаблон обновления:**

- После каждого изменения (добавления, переключения или удаления) мы снова вызываем `loadTodos()`
- Это обеспечивает синхронизацию интерфейса с состоянием в Go
- Альтернатива: возвращайте новое состояние из методов Go, чтобы избежать второго вызова

### Обновите HTML
HTML задаёт структуру нашего приложения для управления задачами. Она минималистична и семантична — основная логика реализована в JavaScript и CSS.

Замените `frontend/index.html`:

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

- `container` — центрирует наше приложение и ограничивает его ширину
- `card` — основная белая карточка, содержащая все элементы
- `input-box` — flex-контейнер для поля ввода и кнопки Add
- `todo-list` — сюда JavaScript будет добавлять отдельные задачи

**Обработка событий:**

- `onkeypress="if(event.key==='Enter') addTodo()"` — добавляет задачу при нажатии Enter
- `onclick="addTodo()"` — добавляет задачу при нажатии кнопки
- Встроенные обработчики событий хорошо подходят для простых приложений на обычном JavaScript

**Скрипт-модуль:**

- `<script type="module">` позволяет использовать импорты ES6
- Наш файл `main.js` может импортировать привязки и использовать современный JavaScript

### Оформите приложение
CSS создаёт современный дизайн в стиле матового стекла с плавными переходами. Мы стремимся придать приложению законченный вид, чтобы им было приятно пользоваться.

Замените `frontend/public/style.css`:

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

**Дизайн в стиле матового стекла:**

- `background: linear-gradient(135deg, #667eea 0%, #764ba2 100%)` — фиолетовый градиентный фон
- `backdrop-filter: blur(10px)` — создаёт на карточке эффект матового стекла
- `rgba(255, 255, 255, 0.95)` — полупрозрачный белый цвет для эффекта стекла

**Настраиваемое оформление флажка:**

- `appearance: none` удаляет стандартное оформление флажка в браузере
- С помощью `::after` мы создаём собственный скруглённый квадрат с галочкой
- При `checked` появляется галочка, созданная с помощью символа Юникода ✓

**Взаимодействие при наведении:**

- При наведении задачи сдвигаются вправо (`transform: translateX(4px)`)
- Кнопка удаления скрыта до наведения (`opacity: 0` → `opacity: 1`)
- При наведении кнопки немного увеличиваются, создавая ощущение тактильного отклика

**Пустое состояние:**

- `#todo-list:empty::before` показывает сообщение, когда задач нет
- Решение реализовано только на CSS — JavaScript не требуется

### Запустите приложение
Посмотрим, как оно работает! Запустите сервер разработки:

```bash
wails3 dev
```

Приложение скомпилируется и откроется. Попробуйте его в работе:

- Введите задачу и нажмите Enter или кнопку Add
- Установите флажок, чтобы отметить задачу как выполненную
- Наведите указатель на задачу, чтобы появилась кнопка удаления
- Обратите внимание, что интерфейс обновляется мгновенно — так работает наш шаблон обновления

**Что происходит:**

- Wails автоматически создал привязки для методов TodoService
- Режим разработки поддерживает горячую перезагрузку — попробуйте изменить CSS и посмотрите, как обновится приложение
- Код Go выполняется нативно — его не нужно транслировать или интерпретировать

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

- Несколько вызовов из фронтенда могут выполняться одновременно
- Операции чтения не блокируют друг друга
- Операции записи получают эксклюзивный доступ
- `defer` гарантирует, что блокировки всегда снимаются

### Обработка ошибок

Сервис возвращает ошибки при недопустимых операциях:

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

### Добавление статистики

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

### Добавление функции «Очистить выполненные»

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

### Добавление постоянного хранения данных

Для рабочих приложений добавьте постоянное хранение данных с помощью решения, которое лучше всего подходит вашему приложению, например SQLite, PostgreSQL или размещённого сервиса.

## Сборка для рабочей среды

Когда приложение TODO будет готово к распространению, соберите его для рабочей среды:

```bash
wails3 build
```

В результате в `bin/` будет создан оптимизированный нативный исполняемый файл:

- Компилирует код Go с оптимизациями
- Собирает фронтенд для рабочей среды
- Объединяет все компоненты в один исполняемый файл
- Размер готового приложения обычно составляет 10-20 МБ (для сравнения: у Electron — 150 МБ и более)

Исполняемый файл можно запускать напрямую: среда выполнения не требуется, серверы запускать не нужно. Это настоящее нативное приложение.

## Что вы создали

Вы только что создали полноценное приложение TODO со следующими возможностями:

**Полная реализация CRUD:**

- Создан сервис с операциями создания, чтения, обновления и удаления
- Добавлены проверка входных данных и обработка ошибок
- Вы узнали, как ошибки Go преобразуются в исключения JavaScript

**Потокобезопасное управление состоянием:**

- Для безопасной обработки конкурентного доступа использован `sync.RWMutex`
- Вы разобрались в различиях между блокировками чтения (RLock) и записи (Lock)
- Вы увидели, как `defer` предотвращает взаимные блокировки, гарантируя освобождение ресурсов

**Современный, тщательно проработанный интерфейс:**

- Создан интерфейс в стиле глассморфизма с градиентами и эффектами размытия
- Созданы флажки с собственным оформлением без использования фреймворков
- Добавлены эффекты при наведении и переходы, благодаря которым интерфейс воспринимается как нативный
- Состояние отсутствия элементов реализовано исключительно средствами CSS

**Основы Wails:**

- Регистрация сервисов и автоматическое создание привязок
- Вызов методов Go из JavaScript с помощью async/await
- Синхронизация состояния между Go и фронтендом
- Сборка и упаковка нативного настольного приложения

## Следующие шаги

Теперь, когда вы разобрались с операциями CRUD и управлением состоянием, попробуйте следующее:

- **Добавьте постоянное хранение данных:** сохраняйте задачи между перезапусками приложения с помощью SQLite или другого решения для хранения данных
- **Добавьте другие возможности:** фильтрацию (все/активные/выполненные), редактирование существующих задач и массовые операции
- **Изучите руководство по приложению «Заметки»:** узнайте, как работают файловые операции, в руководстве [«Заметки»](/tutorials/03-notes-vanilla/)
- **Создайте настоящее приложение:** примените эти концепции и разработайте собственное приложение!
