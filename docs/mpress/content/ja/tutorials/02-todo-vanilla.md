---
title: "TODO リスト"
description: "CRUD 操作を備えた完全な TODO リストアプリケーションを構築する"
slug: "tutorials/02-todo-vanilla"
sourcePath: "tutorials/02-todo-vanilla.md"
---

このチュートリアルでは、完全に機能する TODO リストアプリケーションを構築します。QR コードサービスのチュートリアルから一歩進み、状態の管理、複数の操作の処理、洗練されたユーザーインターフェースの作成方法を学びます。

**構築するもの:**

- 追加、完了、削除の機能を備えた完全な TODO アプリ
- スレッドセーフな状態管理（デスクトップアプリでは重要）
- モダンなグラスモーフィズムの UI デザイン
- すべてバニラ JavaScript で実装 — フレームワークは不要

**学ぶ内容:**

- CRUD 操作（作成、読み取り、更新、削除）
- Go で可変状態を安全に管理する方法
- ユーザー入力とバリデーションの処理
- ネイティブのように感じられるレスポンシブ UI の構築

![TODO リストアプリケーション](/assets/todo-app.png)

**所要時間:** 20 分

## プロジェクトを作成する

@steps
### プロジェクトを生成する
まず、新しい Wails プロジェクトを作成します。ここでは、すっきりとした状態から始められるデフォルトのバニラテンプレートを使用します。

```bash
wails3 init -n todo-app -t vanilla-js
cd todo-app
```

これにより、基本構成を備えた新しいプロジェクトが作成されます。Go バックエンドはルートに、フロントエンドコードは `frontend/` ディレクトリに配置されます。

### TODO サービスを作成する
TODO サービスは、アプリケーションの状態を管理し、CRUD 操作用のメソッドを提供します。各リクエストが分離される Web サーバーとは異なり、デスクトップアプリでは複数の操作が同時に実行される可能性があるため、スレッドセーフな状態管理が必要です。

`greetservice.go` を削除し、新しいファイル `todoservice.go` を作成します。

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

**ここで行っていること:**

**`Todo` 構造体:**

- ID、Title、Completed フィールドでデータの構造を定義します
- `json:` タグは、フロントエンド向けにこの構造体を JSON へ変換する方法を Go に指示します
- バインディングジェネレーターから認識できるように、各フィールドはエクスポートされています（先頭が大文字）

**`TodoService` 構造体:**

- `todos []Todo` — すべての TODO 項目を保持するスライス
- `nextID int` — 次に割り当てる ID を追跡します（自動インクリメントを模倣）
- `mu sync.RWMutex` — スレッドセーフなアクセスを実現する読み書きミューテックス

**`sync.RWMutex` によるスレッドセーフ化:**

- デスクトップアプリでは、UI から複数の操作が同時に実行される可能性があります
- `RLock()` により、複数の読み取り処理を同時に実行できます（複数の `GetAll` 呼び出しなど）
- `Lock()` は、書き込み処理に排他的アクセスを付与します（`Add`、`Toggle`、`Delete` など）
- `defer` により、関数が途中で return した場合やパニックした場合でも、ロックが確実に解放されます

**メソッド:**

- `GetAll()` — すべての TODO を返します（データを変更しないため、読み取りロックを使用）
- `Add(title)` — 入力を検証して新しい TODO を作成し、ID をインクリメントします
- `Toggle(id)` — TODO の完了状態を反転します
- `Delete(id)` — スライスから TODO を削除します

**エラー処理:**

- Go の慣例に従い、最後の値として `error` を返します
- 空のタイトルは拒否されます
- 存在しない TODO に対する操作はエラーを返します
- これらのエラーは、フロントエンドで JavaScript の例外になります

### main.go を更新する
TODO サービスを Wails アプリケーションに登録します。`main.go` の `Services` セクションを見つけ、GreetService を TodoService に置き換えます。

```go {title="main.go" highlight="5"}
Services: []application.Service{
    application.NewService(NewTodoService()),
},
```

**ここで行っていること:**

- デフォルトの GreetService を削除し、代わりに TodoService を追加します
- `application.NewService()` は、Wails がサービスを管理できるようにサービスをラップします
- Wails は、このサービスのすべての公開メソッドに対する JavaScript バインディングを自動的に生成します

### フロントエンド UI を作成する
次に、フロントエンドを構築します。ここで Go のメソッドを呼び出し、UI を表示します。構成をシンプルに保ち、バインディングの動作を直接示すため、バニラ JavaScript を使用します。

`frontend/src/main.js` を置き換えます。

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

**ここで行っていること:**

**バインディングのインポート:**

- `import {TodoService} from "../bindings/changeme"` — 自動生成された Go バインディングをインポートします
- 注: `changeme` は、`go.mod` に記載された実際のモジュール名になります

**`loadTodos()` 関数:**

- `TodoService.GetAll()` を呼び出し、Go からすべての TODO を取得します
- テンプレートリテラルを使用して、各 TODO の HTML を生成します
- スタイル設定のため、`completed` クラスを動的に追加または削除します
- `onclick` 属性を使用して、ボタンを関数に関連付けます
- すべての HTML を連結し、DOM に挿入します

**CRUD 関数:**

- `addTodo()` — 入力を検証し、Go の `Add` メソッドを呼び出して、リストを更新します
- `toggleTodo(id)` — Go の `Toggle` メソッドを呼び出して、リストを更新します
- `deleteTodo(id)` — Go の `Delete` メソッドを呼び出して、リストを更新します
- Go の呼び出しは Promise を返すため、すべての関数は async です

**window に関連付ける理由:**

- `window.addTodo = ...` により、HTML の `onclick` 属性から関数にアクセスできるようになります
- これはバニラ JavaScript で使えるシンプルなパターンです（フレームワークでは異なる方法で処理します）
- 本番環境では、代わりに適切なイベントデリゲーションを使用することもできます

**再読み込みのパターン:**

- 各変更操作（追加、完了状態の切り替え、削除）の後に、`loadTodos()` を再度呼び出します
- これにより、UI と Go の状態が常に同期されます
- 別の方法として、2 回目の呼び出しを避けるため、Go のメソッドから新しい状態を返すこともできます

### HTML を更新する
HTML は TODO アプリの構造を定義します。最小限でセマンティックな構成になっており、実際の動作や表現は JavaScript と CSS が担います。

`frontend/index.html` を次の内容に置き換えます。

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

**ここで行っていること：**

**構造：**

- `container` - アプリを中央に配置し、幅を制限します
- `card` - すべての要素を収めるメインの白いカードです
- `input-box` - 入力欄と追加ボタンを配置するフレックスコンテナです
- `todo-list` - JavaScript によって個々の TODO が挿入される場所です

**イベント処理：**

- `onkeypress="if(event.key==='Enter') addTodo()"` - Enter キーが押されたときに TODO を追加します
- `onclick="addTodo()"` - ボタンがクリックされたときに TODO を追加します
- 単純な素の JavaScript アプリでは、インラインイベントハンドラーが適しています

**モジュールスクリプト：**

- `<script type="module">` により、ES6 の import を使用できます
- `main.js` ファイルでは、バインディングを import して最新の JavaScript を使用できます

### アプリのスタイルを整える
CSS を使って、滑らかなトランジションを備えたモダンなグラスモーフィズムデザインを作成します。アプリを快適に使えるよう、洗練された操作感を目指します。

`frontend/public/style.css` を次の内容に置き換えます。

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

**ここで行っていること：**

**グラスモーフィズムデザイン：**

- `background: linear-gradient(135deg, #667eea 0%, #764ba2 100%)` - 紫色のグラデーション背景です
- `backdrop-filter: blur(10px)` - カードに曇りガラスのような効果を加えます
- `rgba(255, 255, 255, 0.95)` - ガラス効果を表現する半透明の白色です

**チェックボックスのカスタムスタイル：**

- `appearance: none` はブラウザー既定のチェックボックスを非表示にします
- `::after` を使って、チェックマーク付きの角丸正方形を独自に作成します
- `checked` のとき、Unicode 文字の ✓ を使ってチェックマークを表示します

**ホバー時のインタラクション：**

- TODO はホバー時に右へスライドします（`transform: translateX(4px)`）
- 削除ボタンはホバーするまで非表示です（`opacity: 0` → `opacity: 1`）
- 触覚的なフィードバックを表現するため、ボタンはホバー時にわずかに拡大します

**空の状態：**

- TODO がない場合、`#todo-list:empty::before` によってメッセージが表示されます
- CSS のみで実現するため、JavaScript は不要です

### アプリを実行する
実際の動作を確認しましょう。開発サーバーを実行します。

```bash
wails3 dev
```

アプリがコンパイルされ、開きます。次の操作を試してください。

- TODO を入力し、Enter キーを押すか、追加ボタンをクリックします
- チェックボックスをクリックして完了済みにします
- TODO にカーソルを合わせ、削除ボタンが表示されることを確認します
- UI が即座に更新されることに注目してください。これは、ここまでに実装した更新パターンによるものです

**内部で行われていること：**

- Wails によって TodoService のメソッド用バインディングが自動生成されました
- 開発モードではホットリロードを利用できます。CSS を変更し、自動的に反映されることを確認してください
- Go コードはネイティブに実行されるため、変換やインタープリターによる解釈は不要です

@end

## 仕組み

### スレッドセーフな状態管理

`sync.RWMutex` により、安全な並行アクセスを実現します。

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

**これが重要な理由：**

- 複数のフロントエンド呼び出しが同時に発生する可能性があります
- 読み取り操作同士は互いにブロックしません
- 書き込み操作は排他的にアクセスします
- `defer` により、ロックが常に解放されます

### エラー処理

無効な操作に対して、サービスはエラーを返します。

```go
func (t *TodoService) Add(title string) (*Todo, error) {
    if title == "" {
        return nil, errors.New("title cannot be empty")
    }
    // ...
}
```

フロントエンドでは、これらのエラーを捕捉できます。

```javascript
try {
    await TodoService.Add(title);
} catch (err) {
    alert('Error: ' + err);
}
```

### 状態の同期

状態を変更するたびに、リスト全体を再読み込みします。

```javascript
window.addTodo = async () => {
    await TodoService.Add(title);  // Mutation
    await loadTodos();             // Refresh
}
```

**別の方法：** 2 回目の呼び出しを避けるには、各メソッドから更新後のリストを返します。

## 機能拡張

### 統計を追加する

次の内容を `todoservice.go` に追加します。

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

フロントエンドに表示します。

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

### 「完了済みをクリア」を追加する

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

### 永続化を追加する

本番環境向けアプリでは、SQLite、PostgreSQL、ホスティングサービスなど、アプリケーションに最適なストレージソリューションを使って永続化を追加します。

## 本番環境向けにビルドする

TODO アプリを配布する準備が整ったら、本番環境向けにビルドします。

```bash
wails3 build
```

これにより、最適化されたネイティブ実行可能ファイルが `bin/` に生成されます。

- Go コードを最適化してコンパイルします
- フロントエンドを本番環境向けにビルドします
- すべてを単一の実行可能ファイルにまとめます
- 生成されるアプリのサイズは通常 10-20MB です（Electron は 150MB 以上）

実行可能ファイルは直接実行できます。ランタイムは不要で、サーバーを起動する必要もありません。完全なネイティブアプリケーションです。

## 作成したもの

これで、次の機能を備えた完全な TODO アプリケーションを作成できました。

**完全な CRUD 実装：**

- Create、Read、Update、Delete の各操作を備えたサービスを作成しました
- 入力検証とエラー処理を追加しました
- Go のエラーが JavaScript の例外になる仕組みを学びました

**スレッドセーフな状態管理：**

- `sync.RWMutex` を使用して、並行アクセスを安全に処理しました
- 読み取りロック（RLock）と書き込みロック（Lock）の違いを理解しました
- `defer` が確実にクリーンアップを行い、デッドロックを防ぐ仕組みを確認しました

**モダンで洗練された UI：**

- グラデーションとぼかし効果を使ったグラスモーフィズムのインターフェースを作成しました
- フレームワークを使わずに、独自のスタイルを適用したチェックボックスを作成しました
- ネイティブらしい操作感を実現するため、ホバー時のインタラクションとトランジションを追加しました
- 純粋な CSS で空の状態を実装しました

**Wails の基礎：**

- サービスの登録とバインディングの自動生成
- async/await を使用した JavaScript からの Go メソッドの呼び出し
- Go とフロントエンド間の状態同期
- ネイティブデスクトップアプリのビルドとパッケージ化

## 次のステップ

CRUD 操作と状態管理を理解したら、次のことを試してみましょう：

- **永続化を追加：** SQLite または別のストレージソリューションを使用し、アプリを再起動しても TODO が保持されるようにします
- **機能をさらに追加：** フィルタリング（すべて／未完了／完了）、既存の TODO の編集、一括操作
- **Notes チュートリアルを確認：** [Notes](/tutorials/03-notes-vanilla/) でファイル操作の仕組みを確認します
- **実用的なものを構築：** ここで学んだ概念を活用して、独自のアプリを構築しましょう！
