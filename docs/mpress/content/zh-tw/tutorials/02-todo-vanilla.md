---
title: "待辦事項清單"
description: "建置具備 CRUD 操作的完整待辦事項清單應用程式"
slug: "tutorials/02-todo-vanilla"
sourcePath: "tutorials/02-todo-vanilla.md"
---

在本教學中，你將建置一個功能完整的待辦事項清單應用程式。這比 QR 碼服務教學更進一步，你將學會如何管理狀態、處理多項操作，以及打造完善的使用者介面。

**你將建置：**

- 具備新增、完成及刪除功能的完整待辦事項應用程式
- 執行緒安全的狀態管理（對桌面應用程式很重要）
- 現代化的玻璃擬態 UI 設計
- 全部使用原生 JavaScript，無須任何框架

**你將學會：**

- CRUD 操作（建立、讀取、更新、刪除）
- 在 Go 中安全地管理可變狀態
- 處理及驗證使用者輸入
- 建置具備原生應用程式體驗的響應式 UI

![待辦事項清單應用程式](/assets/todo-app.png)

<strong>完成所需時間：</strong>20分鐘

## 建立專案

@steps
### 產生專案
首先，建立新的 Wails 專案。我們將使用預設的原生 JavaScript 範本，作為簡潔的起點：

```bash
wails3 init -n todo-app
cd todo-app
```

這會建立具備基本結構的新專案：Go 後端位於根目錄，前端程式碼則位於`frontend/`目錄。

### 建立待辦事項服務
待辦事項服務將管理應用程式狀態，並提供執行 CRUD 操作的方法。Web 伺服器的每個要求彼此隔離，但桌面應用程式可能同時執行多項操作，因此我們需要以執行緒安全的方式管理狀態。

刪除`greetservice.go`，並建立新檔案`todoservice.go`：

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

**這裡的運作方式：**

**`Todo`結構：**

- 使用 ID、Title 和 Completed 欄位定義資料結構
- `json:`標記會告訴 Go 如何將此結構轉換成供前端使用的 JSON
- 每個欄位都會匯出（以大寫字母開頭），讓繫結產生器能夠識別

**`TodoService`結構：**

- `todos []Todo`—儲存所有待辦事項的切片
- `nextID int`—追蹤下一個要指派的 ID（模擬自動遞增）
- `mu sync.RWMutex`—用於執行緒安全存取的讀寫互斥鎖

**使用`sync.RWMutex`確保執行緒安全：**

- 桌面應用程式可能會從 UI 同時執行多項操作
- `RLock()`允許多個讀取者同時存取（例如同時多次呼叫`GetAll`）
- `Lock()`為寫入操作提供獨佔存取權（例如`Add`、`Toggle`、`Delete`）
- 即使函式提早傳回或發生 panic，`defer`也能確保鎖定解除

**方法：**

- `GetAll()`—傳回所有待辦事項（因為不會修改資料，所以使用讀取鎖定）
- `Add(title)`—建立新的待辦事項、驗證輸入並遞增 ID
- `Toggle(id)`—切換待辦事項的完成狀態
- `Delete(id)`—從切片中移除待辦事項

**錯誤處理：**

- 依照 Go 慣例，我們將`error`作為最後一個值傳回
- 空白標題會遭到拒絕
- 對不存在的待辦事項執行操作時會傳回錯誤
- 這些錯誤會在前端成為 JavaScript 例外狀況

### 更新 main.go
向 Wails 應用程式註冊待辦事項服務。在`main.go`中找到`Services`區段，並以 TodoService 取代 GreetService：

```go {title="main.go" highlight="5"}
Services: []application.Service{
    application.NewService(NewTodoService()),
},
```

**這裡的運作方式：**

- 我們會移除預設的 GreetService，改為加入 TodoService
- `application.NewService()`會包裝我們的服務，讓 Wails 能夠管理它
- Wails 會自動為此服務上的所有公開方法產生 JavaScript 繫結

### 建立前端 UI
現在來建置前端。我們會在這裡呼叫 Go 方法並顯示 UI。為了保持簡單並直接呈現繫結的運作方式，我們將使用原生 JavaScript。

取代`frontend/src/main.js`：

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

**這裡的運作方式：**

**匯入繫結：**

- `import {TodoService} from "../bindings/changeme"`—匯入自動產生的 Go 繫結
- 請注意：`changeme`會是`go.mod`中的實際模組名稱

**`loadTodos()`函式：**

- 呼叫`TodoService.GetAll()`，從 Go 取得所有待辦事項
- 使用範本字面值為每個待辦事項建置 HTML
- 動態新增或移除`completed`類別，以套用樣式
- 使用`onclick`屬性將按鈕連結至我們的函式
- 將所有 HTML 串接起來並注入 DOM

**CRUD 函式：**

- `addTodo()`—驗證輸入、呼叫 Go 的`Add`方法，然後重新整理清單
- `toggleTodo(id)`—呼叫 Go 的`Toggle`方法，然後重新整理清單
- `deleteTodo(id)`—呼叫 Go 的`Delete`方法，然後重新整理清單
- 由於 Go 呼叫會傳回 Promise，因此所有函式都是非同步函式

**為何要附加至 window：**

- `window.addTodo = ...`讓 HTML 的`onclick`屬性能夠存取這些函式
- 這是原生 JavaScript 的簡單模式（框架會以不同方式處理）
- 在正式環境中，你可以改用適當的事件委派

**重新整理模式：**

- 每次變更（新增、切換或刪除）後，我們都會再次呼叫`loadTodos()`
- 這可確保 UI 與 Go 狀態保持同步
- 替代方案：讓 Go 方法傳回新狀態，以免再次呼叫

### 更新 HTML
HTML 為 TODO 應用程式提供結構。它簡潔且具語意；真正發揮作用的是 JavaScript 和 CSS。

取代`frontend/index.html`：

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

**這裡的運作方式：**

**結構：**

- `container`－讓應用程式置中並限制寬度
- `card`－容納所有內容的主要白色卡片
- `input-box`－輸入欄位與「新增」按鈕的彈性容器
- `todo-list`－JavaScript 將個別待辦事項插入此處

**事件處理：**

- `onkeypress="if(event.key==='Enter') addTodo()"`－按下 Enter 時新增待辦事項
- `onclick="addTodo()"`－按一下按鈕時新增待辦事項
- 行內事件處理常式很適合簡單的原生 JavaScript 應用程式

**模組指令碼：**

- `<script type="module">`讓我們能使用 ES6 匯入
- 我們的`main.js`檔案可以匯入繫結並使用現代 JavaScript

### 設定應用程式樣式
CSS 以流暢的轉場效果打造現代玻璃擬態設計。我們希望呈現精緻的質感，讓這個應用程式用起來令人愉快。

取代`frontend/public/style.css`：

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

**這裡的運作方式：**

**玻璃擬態設計：**

- `background: linear-gradient(135deg, #667eea 0%, #764ba2 100%)`－紫色漸層背景
- `backdrop-filter: blur(10px)`－在卡片上建立毛玻璃效果
- `rgba(255, 255, 255, 0.95)`－以半透明白色呈現玻璃效果

**自訂核取方塊樣式：**

- `appearance: none`會移除瀏覽器預設的核取方塊
- 我們使用`::after`建立帶有勾選記號的自訂圓角方塊
- 當`checked`時，會使用 Unicode 字元 ✓ 顯示勾選記號

**游標停留互動：**

- 游標停留時，待辦事項會向右滑動（`transform: translateX(4px)`）
- 刪除按鈕會保持隱藏，直到游標停留時才顯示（`opacity: 0` → `opacity: 1`）
- 游標停留時，按鈕會稍微放大，以提供觸覺般的回饋

**空白狀態：**

- 沒有待辦事項時，`#todo-list:empty::before`會顯示訊息
- 僅使用 CSS 的解決方案，不需要 JavaScript

### 執行應用程式
來看看實際效果！執行開發伺服器：

```bash
wails3 dev
```

應用程式會完成編譯並開啟。試著操作看看：

- 輸入待辦事項，然後按 Enter 或按一下「新增」
- 按一下核取方塊，將待辦事項標示為已完成
- 將游標移至待辦事項上，查看刪除按鈕出現
- 請注意 UI 會立即更新，這表示重新整理模式正在運作

**運作方式：**

- Wails 已自動為你的 TodoService 方法產生繫結
- 開發模式包含熱重新載入功能；試著修改 CSS，並觀察畫面更新
- 你的 Go 程式碼以原生方式執行，不需要轉譯或直譯

@end

## 運作方式

### 執行緒安全的狀態管理

`sync.RWMutex`可提供安全的並行存取：

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

**這點為何重要：**

- 多個前端呼叫可能同時發生
- 讀取作業不會互相阻塞
- 寫入作業會取得獨佔存取權
- `defer`可確保鎖定一律會被釋放

### 錯誤處理

服務會針對無效作業傳回錯誤：

```go
func (t *TodoService) Add(title string) (*Todo, error) {
    if title == "" {
        return nil, errors.New("title cannot be empty")
    }
    // ...
}
```

你可以在前端攔截這些錯誤：

```javascript
try {
    await TodoService.Add(title);
} catch (err) {
    alert('Error: ' + err);
}
```

### 狀態同步

每次變更後，我們都會重新載入完整清單：

```javascript
window.addTodo = async () => {
    await TodoService.Add(title);  // Mutation
    await loadTodos();             // Refresh
}
```

<strong>替代作法：</strong>讓每個方法傳回更新後的清單，以免再次呼叫。

## 強化功能

### 新增統計資料

將以下內容新增至`todoservice.go`：

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

在前端顯示：

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

### 新增「清除已完成項目」

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

### 新增持久化功能

對於正式環境應用程式，請使用最適合應用程式的儲存方案新增持久化功能，例如 SQLite、PostgreSQL 或託管服務。

## 建置正式版本

準備好發佈 TODO 應用程式時，請建置正式版本：

```bash
wails3 build
```

這會在`bin/`中建立最佳化的原生可執行檔：

- 以最佳化設定編譯 Go 程式碼
- 建置正式環境使用的前端
- 將所有內容封裝成單一可執行檔
- 產生的應用程式通常為10-20MB（相較之下，Electron 應用程式為150MB 以上）

你可以直接執行該可執行檔，無須執行階段，也不必啟動任何伺服器。這是真正的原生應用程式。

## 你建置的成果

你剛剛建置了一個完整的 TODO 應用程式，其中包含：

**完整的 CRUD 實作：**

- 建立具備建立、讀取、更新及刪除作業的服務
- 加入輸入驗證與錯誤處理
- 了解 Go 錯誤如何轉換為 JavaScript 例外狀況

**執行緒安全的狀態管理：**

- 使用`sync.RWMutex`安全地處理並行存取
- 了解讀取鎖定（RLock）與寫入鎖定（Lock）之間的差異
- 了解`defer`如何藉由保證執行清理來防止死結

**現代且精緻的使用者介面：**

- 運用漸層與模糊效果建置玻璃擬態介面
- 不使用任何框架，自訂核取方塊的樣式
- 加入游標暫留互動與轉場效果，營造原生操作體驗
- 僅使用 CSS 實作空白狀態

**Wails 基礎知識：**

- 服務註冊與自動產生繫結
- 使用 async/await 從 JavaScript 呼叫 Go 方法
- 在 Go 與前端之間同步狀態
- 建置與封裝原生桌面應用程式

## 後續步驟

了解 CRUD 作業與狀態管理後，可以嘗試：

- <strong>加入持久化：</strong>使用 SQLite 或其他儲存方案，讓待辦事項在應用程式重新啟動後仍然保留
- <strong>加入更多功能：</strong>篩選（全部／進行中／已完成）、編輯現有待辦事項、批次作業
- <strong>探索 Notes 教學課程：</strong>了解[Notes](/tutorials/03-notes-vanilla/)中的檔案作業如何運作
- <strong>打造實際的應用程式：</strong>運用這些概念建置你自己的應用程式！
