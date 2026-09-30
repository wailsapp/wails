---
title: "TODO 목록"
description: "CRUD 작업을 지원하는 완전한 TODO 목록 애플리케이션 만들기"
slug: "tutorials/02-todo-vanilla"
sourcePath: "tutorials/02-todo-vanilla.md"
---

이 튜토리얼에서는 모든 기능을 갖춘 TODO 목록 애플리케이션을 만듭니다. QR Code Service 튜토리얼보다 한 단계 발전한 내용으로, 상태를 관리하고 여러 작업을 처리하며 완성도 높은 사용자 인터페이스를 만드는 방법을 배웁니다.

**만들 결과물:**

- 추가, 완료 및 삭제 기능을 갖춘 완전한 TODO 앱
- 스레드로부터 안전한 상태 관리(데스크톱 앱에서 중요)
- 현대적인 글래스모피즘 UI 디자인
- 프레임워크 없이 모두 순수 JavaScript로 구현

**배울 내용:**

- CRUD 작업(Create, Read, Update, Delete)
- Go에서 변경 가능한 상태를 안전하게 관리하는 방법
- 사용자 입력 처리 및 유효성 검사
- 네이티브 앱처럼 느껴지는 반응형 UI 만들기

![TODO 목록 애플리케이션](/assets/todo-app.png)

**완료 예상 시간:** 20분

## 프로젝트 만들기

@steps
### 프로젝트 생성하기
먼저 새 Wails 프로젝트를 만듭니다. 깔끔하게 시작할 수 있도록 기본 순수 JavaScript 템플릿을 사용하겠습니다.

```bash
wails3 init -n todo-app
cd todo-app
```

이렇게 하면 루트에 Go 백엔드가 있고 `frontend/` 디렉터리에 프런트엔드 코드가 있는 기본 구조의 새 프로젝트가 생성됩니다.

### TODO 서비스 만들기
TODO 서비스는 애플리케이션 상태를 관리하고 CRUD 작업을 위한 메서드를 제공합니다. 각 요청이 격리되는 웹 서버와 달리 데스크톱 앱에서는 여러 작업이 동시에 실행될 수 있으므로, 스레드로부터 안전한 상태 관리가 필요합니다.

`greetservice.go`을 삭제하고 새 파일 `todoservice.go`을 만드세요.

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

**코드의 동작:**

**`Todo` 구조체:**

- ID, Title 및 Completed 필드로 데이터의 형태를 정의합니다.
- `json:` 태그는 프런트엔드에서 사용할 수 있도록 이 구조체를 JSON으로 변환하는 방법을 Go에 알려 줍니다.
- 바인딩 생성기가 볼 수 있도록 각 필드를 내보냅니다(대문자로 시작).

**`TodoService` 구조체:**

- `todos []Todo` - 모든 TODO 항목을 저장하는 슬라이스입니다.
- `nextID int` - 다음에 할당할 ID를 추적합니다(자동 증가를 모방합니다).
- `mu sync.RWMutex` - 스레드로부터 안전한 접근을 위한 읽기/쓰기 뮤텍스입니다.

**`sync.RWMutex`을 사용한 스레드 안전성:**

- 데스크톱 앱에서는 UI에서 여러 작업을 동시에 실행할 수 있습니다.
- `RLock()`을 사용하면 여러 읽기 작업을 동시에 수행할 수 있습니다(예: 여러 `GetAll` 호출).
- `Lock()`은 쓰기 작업에 독점 접근 권한을 부여합니다(예: `Add`, `Toggle`, `Delete`).
- `defer`을 사용하면 함수가 일찍 반환되거나 패닉이 발생하더라도 잠금이 해제됩니다.

**메서드:**

- `GetAll()` - 모든 TODO를 반환합니다(데이터를 수정하지 않으므로 읽기 잠금을 사용합니다).
- `Add(title)` - 새 TODO를 만들고 입력을 검증한 후 ID를 증가시킵니다.
- `Toggle(id)` - TODO의 완료 상태를 반전합니다.
- `Delete(id)` - 슬라이스에서 TODO를 제거합니다.

**오류 처리:**

- Go 관례에 따라 마지막 값으로 `error`을 반환합니다.
- 빈 제목은 거부합니다.
- 존재하지 않는 TODO에 대한 작업은 오류를 반환합니다.
- 이러한 오류는 프런트엔드에서 JavaScript 예외가 됩니다.

### main.go 업데이트하기
Wails 애플리케이션에 TODO 서비스를 등록하세요. `main.go`에서 `Services` 섹션을 찾아 GreetService를 TodoService로 바꾸세요.

```go {title="main.go" highlight="5"}
Services: []application.Service{
    application.NewService(NewTodoService()),
},
```

**코드의 동작:**

- 기본 GreetService를 제거하고 대신 TodoService를 추가합니다.
- Wails가 서비스를 관리할 수 있도록 `application.NewService()`이 서비스를 래핑합니다.
- Wails는 이 서비스의 모든 공개 메서드에 대한 JavaScript 바인딩을 자동으로 생성합니다.

### 프런트엔드 UI 만들기
이제 프런트엔드를 만들어 보겠습니다. 여기에서 Go 메서드를 호출하고 UI를 표시합니다. 구성을 단순하게 유지하면서 바인딩이 직접 작동하는 방식을 보여 주기 위해 순수 JavaScript를 사용합니다.

`frontend/src/main.js`을 다음과 같이 바꾸세요.

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

**코드의 동작:**

**바인딩 가져오기:**

- `import {TodoService} from "../bindings/changeme"` - 자동 생성된 Go 바인딩을 가져옵니다.
- 참고: `changeme`에는 `go.mod`에 정의된 실제 모듈 이름이 들어갑니다.

**`loadTodos()` 함수:**

- Go에서 모든 TODO를 가져오기 위해 `TodoService.GetAll()`을 호출합니다.
- 템플릿 리터럴을 사용해 각 TODO의 HTML을 만듭니다.
- 스타일을 적용하기 위해 `completed` 클래스를 동적으로 추가하거나 제거합니다.
- `onclick` 속성을 사용해 버튼을 함수에 연결합니다.
- 모든 HTML을 하나로 결합해 DOM에 삽입합니다.

**CRUD 함수:**

- `addTodo()` - 입력을 검증하고 Go의 `Add` 메서드를 호출한 후 목록을 새로 고칩니다.
- `toggleTodo(id)` - Go의 `Toggle` 메서드를 호출한 후 목록을 새로 고칩니다.
- `deleteTodo(id)` - Go의 `Delete` 메서드를 호출한 후 목록을 새로 고칩니다.
- Go 호출은 Promise를 반환하므로 모든 함수가 비동기 함수입니다.

**window에 연결하는 이유:**

- `window.addTodo = ...`을 사용하면 HTML의 `onclick` 속성에서 함수에 접근할 수 있습니다.
- 이는 순수 JavaScript를 위한 간단한 패턴입니다(프레임워크에서는 다르게 처리합니다).
- 프로덕션 환경에서는 적절한 이벤트 위임을 대신 사용할 수 있습니다.

**새로 고침 패턴:**

- 각 변경 작업(추가/전환/삭제) 후 `loadTodos()`을 다시 호출합니다.
- 이렇게 하면 UI가 Go 상태와 계속 동기화됩니다.
- 대안: 두 번째 호출을 피하려면 Go 메서드가 새 상태를 반환하도록 하세요.

### HTML 업데이트
HTML은 TODO 앱의 구조를 제공합니다. 최소한의 시맨틱 구조로 구성되어 있으며, 핵심 동작은 JavaScript와 CSS에서 이루어집니다.

`frontend/index.html`을 다음과 같이 바꾸세요:

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

**여기서 일어나는 작업:**

**구조:**

- `container` - 앱을 가운데에 배치하고 너비를 제한합니다.
- `card` - 모든 요소를 담는 기본 흰색 카드입니다.
- `input-box` - 입력 필드와 Add 버튼을 위한 플렉스 컨테이너입니다.
- `todo-list` - JavaScript가 개별 TODO 항목을 삽입하는 곳입니다.

**이벤트 처리:**

- `onkeypress="if(event.key==='Enter') addTodo()"` - Enter를 누르면 TODO 항목을 추가합니다.
- `onclick="addTodo()"` - 버튼을 클릭하면 TODO 항목을 추가합니다.
- 인라인 이벤트 핸들러는 간단한 순수 JavaScript 앱에 적합합니다.

**모듈 스크립트:**

- `<script type="module">`을 사용하면 ES6 import를 사용할 수 있습니다.
- `main.js` 파일에서 바인딩을 가져와 최신 JavaScript를 사용할 수 있습니다.

### 앱 스타일 지정
CSS는 부드러운 전환 효과가 적용된 현대적인 글래스모피즘 디자인을 구현합니다. 앱을 즐겁게 사용할 수 있도록 세련된 느낌을 추구합니다.

`frontend/public/style.css`을 다음과 같이 바꾸세요:

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

**여기서 일어나는 작업:**

**글래스모피즘 디자인:**

- `background: linear-gradient(135deg, #667eea 0%, #764ba2 100%)` - 보라색 그라데이션 배경입니다.
- `backdrop-filter: blur(10px)` - 카드에 반투명 유리 효과를 만듭니다.
- `rgba(255, 255, 255, 0.95)` - 유리 효과를 위한 반투명 흰색입니다.

**사용자 지정 체크박스 스타일:**

- `appearance: none`은 브라우저의 기본 체크박스를 제거합니다.
- `::after`을 사용해 체크 표시가 있는 사용자 지정 둥근 사각형을 만듭니다.
- `checked` 상태가 되면 유니코드 문자 ✓를 사용한 체크 표시가 나타납니다.

**마우스 오버 상호작용:**

- TODO 항목에 마우스를 올리면 오른쪽으로 이동합니다(`transform: translateX(4px)`).
- 삭제 버튼은 마우스를 올릴 때까지 숨겨집니다(`opacity: 0` → `opacity: 1`).
- 촉각적인 피드백을 주기 위해 마우스를 올리면 버튼이 약간 커집니다.

**빈 상태:**

- TODO 항목이 없으면 `#todo-list:empty::before`이 메시지를 표시합니다.
- CSS만 사용하는 방식이므로 JavaScript가 필요하지 않습니다.

### 앱 실행
이제 실제로 작동하는 모습을 확인해 보겠습니다! 개발 서버를 실행하세요:

```bash
wails3 dev
```

앱이 컴파일된 후 열립니다. 다음 기능을 사용해 보세요:

- TODO 항목을 입력하고 Enter를 누르거나 Add를 클릭하세요.
- 완료로 표시하려면 체크박스를 클릭하세요.
- 삭제 버튼이 나타나는 것을 확인하려면 TODO 항목에 마우스를 올리세요.
- UI가 즉시 업데이트되는 것을 확인하세요. 이것이 새로 고침 패턴의 동작입니다.

**일어나는 작업:**

- Wails가 TodoService 메서드의 바인딩을 자동으로 생성했습니다.
- 개발 모드에는 핫 리로드가 포함됩니다. CSS를 변경하고 업데이트되는 모습을 확인해 보세요.
- Go 코드는 네이티브로 실행되므로 변환이나 인터프리트가 필요하지 않습니다.

@end

## 작동 방식

### 스레드 안전 상태 관리

`sync.RWMutex`은 안전한 동시 접근을 제공합니다:

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

**이것이 중요한 이유:**

- 여러 프런트엔드 호출이 동시에 발생할 수 있습니다.
- 읽기 작업은 서로를 차단하지 않습니다.
- 쓰기 작업은 배타적 접근 권한을 얻습니다.
- `defer`을 사용하면 잠금이 항상 해제됩니다.

### 오류 처리

서비스는 유효하지 않은 작업에 대해 오류를 반환합니다:

```go
func (t *TodoService) Add(title string) (*Todo, error) {
    if title == "" {
        return nil, errors.New("title cannot be empty")
    }
    // ...
}
```

프런트엔드에서 이러한 오류를 포착할 수 있습니다:

```javascript
try {
    await TodoService.Add(title);
} catch (err) {
    alert('Error: ' + err);
}
```

### 상태 동기화

각 변경 작업 후 전체 목록을 다시 불러옵니다:

```javascript
window.addTodo = async () => {
    await TodoService.Add(title);  // Mutation
    await loadTodos();             // Refresh
}
```

**대안:** 두 번째 호출을 피하려면 각 메서드에서 업데이트된 목록을 반환하세요.

## 기능 개선

### 통계 추가

`todoservice.go`에 다음을 추가하세요:

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

프런트엔드에 표시하세요:

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

### "완료 항목 지우기" 추가

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

### 영구 저장 기능 추가

프로덕션 앱에서는 SQLite, PostgreSQL 또는 호스팅 서비스 등 애플리케이션에 가장 적합한 스토리지 솔루션을 사용해 영구 저장 기능을 추가하세요.

## 프로덕션용 빌드

TODO 앱을 배포할 준비가 되면 프로덕션용으로 빌드하세요:

```bash
wails3 build
```

그러면 `bin/`에 최적화된 네이티브 실행 파일이 생성됩니다:

- Go 코드를 최적화하여 컴파일합니다
- 프런트엔드를 프로덕션용으로 빌드합니다
- 모든 항목을 하나의 실행 파일로 묶습니다
- 생성된 앱의 크기는 일반적으로 10-20MB입니다(Electron의 150MB 이상과 비교)

런타임이나 시작할 서버 없이 실행 파일을 직접 실행할 수 있습니다. 진정한 네이티브 애플리케이션입니다.

## 완성한 내용

이제 다음 기능을 갖춘 완전한 TODO 애플리케이션을 만들었습니다:

**완전한 CRUD 구현:**

- Create, Read, Update, Delete 작업을 제공하는 서비스를 만들었습니다
- 입력 유효성 검사와 오류 처리를 추가했습니다
- Go 오류가 JavaScript 예외로 변환되는 방식을 배웠습니다

**스레드 안전 상태 관리:**

- 동시 접근을 안전하게 처리하기 위해 `sync.RWMutex`을 사용했습니다
- 읽기 잠금(RLock)과 쓰기 잠금(Lock)의 차이를 이해했습니다
- `defer`이 정리를 보장하여 교착 상태를 방지하는 방식을 확인했습니다

**현대적이고 세련된 UI:**

- 그라데이션과 블러 효과를 적용한 글래스모피즘 인터페이스를 만들었습니다
- 프레임워크 없이 사용자 정의 스타일의 체크박스를 만들었습니다
- 네이티브와 같은 사용감을 위해 호버 상호작용과 전환 효과를 추가했습니다
- 순수 CSS로 빈 상태를 구현했습니다

**Wails 기본 사항:**

- 서비스 등록 및 바인딩 자동 생성
- async/await를 사용하여 JavaScript에서 Go 메서드 호출
- Go와 프런트엔드 간 상태 동기화
- 네이티브 데스크톱 앱 빌드 및 패키징

## 다음 단계

이제 CRUD 작업과 상태 관리를 이해했으므로 다음을 시도해 보세요:

- **영속성 추가:** SQLite 또는 다른 스토리지 솔루션을 사용하여 앱을 다시 시작해도 TODO 항목이 유지되도록 하세요
- **기능 추가:** 필터링(전체/진행 중/완료), 기존 TODO 항목 편집, 일괄 작업
- **Notes 튜토리얼 살펴보기:** [Notes](/tutorials/03-notes-vanilla/)에서 파일 작업이 어떻게 이루어지는지 확인하세요
- **실제 앱 만들기:** 이러한 개념을 활용하여 자신만의 앱을 만들어 보세요!
