---
title: "待办事项列表"
description: "构建一个支持 CRUD 操作的完整待办事项列表应用"
slug: "tutorials/02-todo-vanilla"
sourcePath: "tutorials/02-todo-vanilla.md"
---

在本教程中，你将构建一个功能完备的待办事项列表应用。本教程比二维码服务教程更进一步——你将学习如何管理状态、处理多项操作，以及创建精致的用户界面。

**你将构建：**

- 一个具备添加、完成和删除功能的完整待办事项应用
- 线程安全的状态管理（对桌面应用很重要）
- 现代的玻璃拟态 UI 设计
- 全部使用原生 JavaScript 实现，无需框架

**你将学到：**

- CRUD 操作（创建、读取、更新、删除）
- 在 Go 中安全地管理可变状态
- 处理和验证用户输入
- 构建具有原生应用体验的响应式 UI

![待办事项列表应用](/assets/todo-app.png)

<strong>完成所需时间：</strong>20分钟

## 创建项目

@steps
### 生成项目
首先，创建一个新的 Wails 项目。我们将使用默认的原生模板，以便从一个简洁的基础开始：

```bash
wails3 init -n todo-app
cd todo-app
```

这会创建一个具有基本结构的新项目：Go 后端位于根目录，前端代码位于`frontend/`目录中。

### 创建待办事项服务
待办事项服务将管理应用状态，并提供执行 CRUD 操作的方法。Web 服务器的每个请求彼此隔离，而桌面应用可能同时执行多项操作，因此我们需要采用线程安全的状态管理方式。

删除`greetservice.go`，然后创建新文件`todoservice.go`：

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

**这里发生了什么：**

**`Todo`结构体：**

- 使用 ID、Title 和 Completed 字段定义数据结构
- `json:`标签告诉 Go 如何将此结构体转换为供前端使用的 JSON
- 每个字段都已导出（首字母大写），因此绑定生成器可以识别它们

**`TodoService`结构体：**

- `todos []Todo`——保存所有待办事项的切片
- `nextID int`——记录下一个要分配的 ID（模拟自增）
- `mu sync.RWMutex`——用于实现线程安全访问的读写互斥锁

**使用`sync.RWMutex`实现线程安全：**

- 桌面应用的 UI 可能同时发起多项操作
- `RLock()`允许多个读取操作同时进行（例如，多次调用`GetAll`）
- `Lock()`为写入操作提供独占访问权限（例如`Add`、`Toggle`和`Delete`）
- `defer`确保即使函数提前返回或发生 panic，也会释放锁

**方法：**

- `GetAll()`——返回所有待办事项（由于不修改数据，因此使用读锁）
- `Add(title)`——创建新的待办事项、验证输入并递增 ID
- `Toggle(id)`——切换待办事项的完成状态
- `Delete(id)`——从切片中移除待办事项

**错误处理：**

- 遵循 Go 的惯例，我们将`error`作为最后一个返回值
- 拒绝空标题
- 对不存在的待办事项执行操作时会返回错误
- 这些错误会在前端转换为 JavaScript 异常

### 更新 main.go
在 Wails 应用中注册待办事项服务。在`main.go`中找到`Services`部分，然后用 TodoService 替换 GreetService：

```go {title="main.go" highlight="5"}
Services: []application.Service{
    application.NewService(NewTodoService()),
},
```

**这里发生了什么：**

- 我们移除默认的 GreetService，改为添加 TodoService
- `application.NewService()`会封装我们的服务，以便 Wails 对其进行管理
- Wails 会自动为此服务的所有公开方法生成 JavaScript 绑定

### 创建前端 UI
现在来构建前端。我们将在这里调用 Go 方法并显示 UI。为了保持简单，并直接展示绑定的工作方式，我们使用原生 JavaScript。

替换`frontend/src/main.js`：

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

**这里发生了什么：**

**导入绑定：**

- `import {TodoService} from "../bindings/changeme"`——导入自动生成的 Go 绑定
- 注意：`changeme`应替换为`go.mod`中的实际模块名称

**`loadTodos()`函数：**

- 调用`TodoService.GetAll()`，从 Go 获取所有待办事项
- 使用模板字面量为每个待办事项构建 HTML
- 动态添加或移除`completed`类以应用样式
- 使用`onclick`属性将按钮连接到我们的函数
- 将所有 HTML 拼接在一起并注入 DOM

**CRUD 函数：**

- `addTodo()`——验证输入、调用 Go 的`Add`方法，然后刷新列表
- `toggleTodo(id)`——调用 Go 的`Toggle`方法，然后刷新列表
- `deleteTodo(id)`——调用 Go 的`Delete`方法，然后刷新列表
- 所有函数都是异步函数，因为 Go 调用会返回 Promise

**为什么要挂载到 window：**

- `window.addTodo = ...`使 HTML 的`onclick`属性能够访问这些函数
- 这是一个适用于原生 JavaScript 的简单模式（框架的处理方式有所不同）
- 在生产环境中，可以改用规范的事件委托

**刷新模式：**

- 每次更改（添加、切换状态或删除）后，我们都会再次调用`loadTodos()`
- 这样可以确保 UI 与 Go 中的状态保持同步
- 另一种方式：让 Go 方法返回新状态，从而避免第二次调用

### 更新 HTML
HTML 为 TODO 应用提供结构。它简洁且语义清晰，真正发挥作用的是 JavaScript 和 CSS。

替换`frontend/index.html`：

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

**这里的工作原理：**

**结构：**

- `container`— 将应用居中并限制其宽度
- `card`— 容纳所有内容的白色主卡片
- `input-box`— 用于输入框和“添加”按钮的弹性容器
- `todo-list`— JavaScript 将各个待办事项插入此处

**事件处理：**

- `onkeypress="if(event.key==='Enter') addTodo()"`— 按下 Enter 时添加待办事项
- `onclick="addTodo()"`— 单击按钮时添加待办事项
- 内联事件处理程序很适合简单的原生 JavaScript 应用

**模块脚本：**

- `<script type="module">`允许我们使用 ES6 导入
- 我们的`main.js`文件可以导入绑定并使用现代 JavaScript

### 设置应用样式
CSS 通过平滑过渡营造出具有现代感的玻璃拟态设计。我们的目标是打造精致的使用体验，让应用用起来更加愉悦。

替换`frontend/public/style.css`：

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

**这里的工作原理：**

**玻璃拟态设计：**

- `background: linear-gradient(135deg, #667eea 0%, #764ba2 100%)`— 紫色渐变背景
- `backdrop-filter: blur(10px)`— 在卡片上营造磨砂玻璃效果
- `rgba(255, 255, 255, 0.95)`— 使用半透明白色营造玻璃效果

**自定义复选框样式：**

- `appearance: none`会移除浏览器的默认复选框
- 我们使用`::after`创建带有对勾的自定义圆角方框
- 当`checked`时，会使用 Unicode 字符 ✓ 显示对勾

**悬停交互：**

- 悬停时，待办事项向右滑动（`transform: translateX(4px)`）
- 删除按钮会一直隐藏，直到悬停时才显示（`opacity: 0` → `opacity: 1`）
- 悬停时按钮会略微放大，以提供触感反馈

**空状态：**

- 没有待办事项时，`#todo-list:empty::before`会显示一条消息
- 纯 CSS 解决方案，无需 JavaScript

### 运行应用
来看看实际效果！运行开发服务器：

```bash
wails3 dev
```

应用将完成编译并打开。请尝试以下操作：

- 输入待办事项，然后按 Enter 或单击“添加”
- 单击复选框，将其标记为已完成
- 将鼠标悬停在待办事项上，查看删除按钮出现
- 注意 UI 如何即时更新——这正是刷新模式在发挥作用

**工作原理：**

- Wails 自动为 TodoService 方法生成了绑定
- 开发模式包含热重载功能——尝试修改 CSS，观察它自动更新
- Go 代码以原生方式运行，无需转换或解释

@end

## 工作原理

### 线程安全的状态管理

`sync.RWMutex`提供安全的并发访问：

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

**为什么这很重要：**

- 多个前端调用可能并发发生
- 读取操作不会相互阻塞
- 写入操作会获得独占访问权
- `defer`可确保锁始终得到释放

### 错误处理

对于无效操作，服务会返回错误：

```go
func (t *TodoService) Add(title string) (*Todo, error) {
    if title == "" {
        return nil, errors.New("title cannot be empty")
    }
    // ...
}
```

在前端中，可以捕获这些错误：

```javascript
try {
    await TodoService.Add(title);
} catch (err) {
    alert('Error: ' + err);
}
```

### 状态同步

每次更改后，我们都会重新加载完整列表：

```javascript
window.addTodo = async () => {
    await TodoService.Add(title);  // Mutation
    await loadTodos();             // Refresh
}
```

<strong>另一种方式：</strong>让每个方法返回更新后的列表，从而避免第二次调用。

## 增强功能

### 添加统计信息

将以下内容添加到`todoservice.go`：

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

在前端中显示：

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

### 添加“清除已完成项”

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

### 添加持久化

对于生产应用，请使用最适合应用的存储解决方案添加持久化，例如 SQLite、PostgreSQL 或托管服务。

## 构建生产版本

准备分发 TODO 应用时，请构建生产版本：

```bash
wails3 build
```

这会在`bin/`中创建经过优化的原生可执行文件：

- 启用优化来编译 Go 代码
- 为生产环境构建前端
- 将所有内容打包到一个可执行文件中
- 生成的应用通常为10-20MB（相比之下，Electron 应用为150MB 以上）

你可以直接运行该可执行文件——无需运行时，也无需启动服务器。它是真正的原生应用。

## 你构建的内容

你刚刚构建了一个完整的 TODO 应用，其中包括：

**完整的 CRUD 实现：**

- 创建了一个提供创建、读取、更新和删除操作的服务
- 添加了输入验证和错误处理
- 了解了 Go 错误如何转换为 JavaScript 异常

**线程安全的状态管理：**

- 使用`sync.RWMutex`安全处理并发访问
- 理解了读锁（RLock）与写锁（Lock）之间的区别
- 了解了`defer`如何通过保证执行清理操作来防止死锁

**现代而精致的 UI：**

- 构建了带有渐变和模糊效果的玻璃拟态界面
- 不使用任何框架，创建了自定义样式的复选框
- 添加了悬停交互和过渡效果，营造原生应用体验
- 使用纯 CSS 实现了空状态

**Wails 基础知识：**

- 服务注册和自动生成绑定
- 使用 async/await 从 JavaScript 调用 Go 方法
- 在 Go 与前端之间同步状态
- 构建和打包原生桌面应用

## 后续步骤

现在你已经了解 CRUD 操作和状态管理，可以尝试：

- <strong>添加持久化：</strong>使用 SQLite 或其他存储解决方案，让待办事项在应用重启后仍然保留
- <strong>添加更多功能：</strong>筛选（全部/未完成/已完成）、编辑现有待办事项、批量操作
- <strong>探索 Notes 教程：</strong>了解[Notes](/tutorials/03-notes-vanilla/)中的文件操作如何工作
- <strong>构建真正实用的应用：</strong>运用这些概念构建你自己的应用！
