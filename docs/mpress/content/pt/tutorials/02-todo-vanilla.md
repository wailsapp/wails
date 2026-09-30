---
title: "Lista de tarefas"
description: "Crie um aplicativo completo de lista de tarefas com operações CRUD"
slug: "tutorials/02-todo-vanilla"
sourcePath: "tutorials/02-todo-vanilla.md"
---

Neste tutorial, você criará um aplicativo de lista de tarefas totalmente funcional. Este tutorial é um avanço em relação ao tutorial do Serviço de Código QR: você aprenderá a gerenciar o estado, lidar com várias operações e criar uma interface do usuário refinada.

**O que você criará:**

- Um aplicativo completo de lista de tarefas com funcionalidades para adicionar, concluir e excluir tarefas
- Gerenciamento de estado seguro para threads (importante para aplicativos desktop)
- Design de interface moderno com efeito de vidro
- Tudo usando JavaScript puro, sem precisar de frameworks

**O que você aprenderá:**

- Operações CRUD (criar, ler, atualizar e excluir)
- Gerenciamento seguro de estado mutável em Go
- Tratamento e validação da entrada do usuário
- Criação de interfaces responsivas com aparência nativa

![Aplicativo de lista de tarefas](/assets/todo-app.png)

**Tempo para concluir:** 20 minutos

## Crie seu projeto

@steps
### Gere o projeto
Primeiro, crie um novo projeto Wails. Usaremos o modelo padrão com JavaScript puro, que oferece um ponto de partida sem elementos desnecessários:

```bash
wails3 init -n todo-app
cd todo-app
```

Isso cria um novo projeto com a estrutura básica: o backend em Go na raiz e o código do frontend no diretório `frontend/`.

### Crie o serviço de tarefas
O serviço de tarefas gerenciará o estado do nosso aplicativo e fornecerá métodos para operações CRUD. Diferentemente de um servidor web, no qual cada solicitação é isolada, os aplicativos desktop podem executar várias operações simultaneamente. Portanto, precisamos de um gerenciamento de estado seguro para threads.

Exclua `greetservice.go` e crie um novo arquivo `todoservice.go`:

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

**O que está acontecendo aqui:**

**A struct `Todo`:**

- Define o formato dos nossos dados com os campos ID, Title e Completed
- As tags `json:` informam ao Go como converter essa struct em JSON para o frontend
- Cada campo é exportado (começa com letra maiúscula) para que o gerador de bindings possa detectá-lo

**A struct `TodoService`:**

- `todos []Todo` — um slice que armazena todas as nossas tarefas
- `nextID int` — controla o próximo ID a ser atribuído (simula o incremento automático)
- `mu sync.RWMutex` — um mutex de leitura e escrita para acesso seguro entre threads

**Segurança entre threads com `sync.RWMutex`:**

- Aplicativos desktop podem receber várias operações simultâneas da interface do usuário
- `RLock()` permite vários leitores ao mesmo tempo (por exemplo, várias chamadas a `GetAll`)
- `Lock()` concede acesso exclusivo para gravações (por exemplo, `Add`, `Toggle` e `Delete`)
- `defer` garante que os bloqueios sejam liberados mesmo que a função retorne antecipadamente ou entre em pânico

**Os métodos:**

- `GetAll()` — Retorna todas as tarefas (usa um bloqueio de leitura, pois não estamos modificando os dados)
- `Add(title)` — Cria uma nova tarefa, valida a entrada e incrementa o ID
- `Toggle(id)` — Alterna o status de conclusão de uma tarefa
- `Delete(id)` — Remove uma tarefa do slice

**Tratamento de erros:**

- Retornamos `error` como o último valor, seguindo as convenções do Go
- Títulos vazios são rejeitados
- Operações em tarefas inexistentes retornam erros
- Esses erros se tornam exceções de JavaScript no frontend

### Atualize main.go
Registre o serviço de tarefas no seu aplicativo Wails. Localize a seção `Services` em `main.go` e substitua GreetService pelo nosso TodoService:

```go {title="main.go" highlight="5"}
Services: []application.Service{
    application.NewService(NewTodoService()),
},
```

**O que está acontecendo aqui:**

- Estamos removendo o GreetService padrão e adicionando nosso TodoService no lugar dele
- `application.NewService()` encapsula nosso serviço para que o Wails possa gerenciá-lo
- O Wails gerará automaticamente bindings JavaScript para todos os métodos públicos desse serviço

### Crie a interface do frontend
Agora, vamos criar o frontend. É nele que chamaremos nossos métodos Go e exibiremos a interface do usuário. Usaremos JavaScript puro para simplificar e mostrar diretamente como os bindings funcionam.

Substitua `frontend/src/main.js`:

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

**O que está acontecendo aqui:**

**Importação dos bindings:**

- `import {TodoService} from "../bindings/changeme"` — importa os bindings Go gerados automaticamente
- Observação: `changeme` será o nome real do seu módulo, definido em `go.mod`

**A função `loadTodos()`:**

- Chama `TodoService.GetAll()` para buscar todas as tarefas no Go
- Cria o HTML de cada tarefa usando template literals
- Adiciona ou remove dinamicamente a classe `completed` para aplicar a estilização
- Usa atributos `onclick` para conectar os botões às nossas funções
- Concatena todo o HTML e o insere no DOM

**As funções CRUD:**

- `addTodo()` — Valida a entrada, chama o método Go `Add` e atualiza a lista
- `toggleTodo(id)` — Chama o método `Toggle` do Go e atualiza a lista
- `deleteTodo(id)` — Chama o método `Delete` do Go e atualiza a lista
- Todas as funções são assíncronas porque as chamadas ao Go retornam Promises

**Por que anexar a window:**

- `window.addTodo = ...` torna as funções acessíveis por meio dos atributos HTML `onclick`
- Este é um padrão simples para JavaScript puro (os frameworks lidam com isso de outra forma)
- Em produção, você pode optar por uma delegação de eventos adequada

**O padrão de atualização:**

- Após cada alteração (adicionar/alternar/excluir), chamamos `loadTodos()` novamente
- Isso garante que a interface permaneça sincronizada com o estado no Go
- Alternativa: faça os métodos Go retornarem o novo estado para evitar a segunda chamada

### Atualize o HTML
O HTML fornece a estrutura do nosso aplicativo de tarefas. Ele é mínimo e semântico — a mágica acontece no JavaScript e no CSS.

Substitua `frontend/index.html`:

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

**O que acontece aqui:**

**A estrutura:**

- `container` — centraliza nosso aplicativo e limita sua largura
- `card` — o cartão branco principal que contém tudo
- `input-box` — contêiner flexível para o campo de entrada e o botão Adicionar
- `todo-list` — onde o JavaScript inserirá cada tarefa

**Tratamento de eventos:**

- `onkeypress="if(event.key==='Enter') addTodo()"` — adiciona uma tarefa quando Enter é pressionado
- `onclick="addTodo()"` — adiciona uma tarefa quando o botão é clicado
- Manipuladores de eventos inline funcionam bem em aplicativos simples feitos com JavaScript puro

**Script de módulo:**

- `<script type="module">` permite usar importações ES6
- Nosso arquivo `main.js` pode importar os bindings e usar JavaScript moderno

### Estilize o aplicativo
O CSS cria um design moderno de glassmorphism com transições suaves. Buscamos um acabamento refinado que torne agradável usar o aplicativo.

Substitua `frontend/public/style.css`:

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

**O que acontece aqui:**

**Design com glassmorphism:**

- `background: linear-gradient(135deg, #667eea 0%, #764ba2 100%)` — fundo com gradiente roxo
- `backdrop-filter: blur(10px)` — cria o efeito de vidro fosco no cartão
- `rgba(255, 255, 255, 0.95)` — branco semitransparente para produzir o efeito de vidro

**Estilização personalizada da caixa de seleção:**

- `appearance: none` remove a caixa de seleção padrão do navegador
- Criamos um quadrado arredondado personalizado com uma marca de seleção usando `::after`
- A marca de seleção aparece quando `checked`, usando o caractere Unicode ✓

**Interações ao passar o cursor:**

- As tarefas deslizam para a direita quando o cursor passa sobre elas (`transform: translateX(4px)`)
- O botão de exclusão fica oculto até que o cursor passe sobre a tarefa (`opacity: 0` → `opacity: 1`)
- Os botões aumentam ligeiramente quando o cursor passa sobre eles, proporcionando uma resposta tátil

**Estado vazio:**

- `#todo-list:empty::before` exibe uma mensagem quando não há tarefas
- Solução apenas com CSS — não requer JavaScript

### Execute o aplicativo
Vamos vê-lo em ação! Execute o servidor de desenvolvimento:

```bash
wails3 dev
```

O aplicativo será compilado e aberto. Experimente:

- Digite uma tarefa e pressione Enter ou clique em Adicionar
- Clique na caixa de seleção para marcar a tarefa como concluída
- Passe o cursor sobre uma tarefa para exibir o botão de exclusão
- Observe como a interface é atualizada instantaneamente — é o nosso padrão de atualização em ação

**O que acontece:**

- O Wails gerou automaticamente bindings para os métodos do seu TodoService
- O modo de desenvolvimento inclui recarregamento automático — experimente alterar o CSS e veja a atualização
- Seu código Go está sendo executado nativamente — não é necessário traduzi-lo nem interpretá-lo

@end

## Como funciona

### Gerenciamento de estado seguro para concorrência

O `sync.RWMutex` fornece acesso concorrente seguro:

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

**Por que isso é importante:**

- Várias chamadas do frontend podem ocorrer simultaneamente
- As operações de leitura não bloqueiam umas às outras
- As operações de escrita obtêm acesso exclusivo
- `defer` garante que os bloqueios sejam sempre liberados

### Tratamento de erros

O serviço retorna erros para operações inválidas:

```go
func (t *TodoService) Add(title string) (*Todo, error) {
    if title == "" {
        return nil, errors.New("title cannot be empty")
    }
    // ...
}
```

No frontend, você pode capturar esses erros:

```javascript
try {
    await TodoService.Add(title);
} catch (err) {
    alert('Error: ' + err);
}
```

### Sincronização de estado

Após cada alteração, recarregamos a lista completa:

```javascript
window.addTodo = async () => {
    await TodoService.Add(title);  // Mutation
    await loadTodos();             // Refresh
}
```

**Abordagem alternativa:** retorne a lista atualizada de cada método para evitar a segunda chamada.

## Aprimoramentos

### Adicione estatísticas

Adicione isto a `todoservice.go`:

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

Exiba no frontend:

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

### Adicione “Limpar concluídas”

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

### Adicione persistência

Para aplicativos de produção, adicione persistência usando a solução de armazenamento mais adequada ao seu aplicativo, como SQLite, PostgreSQL ou um serviço hospedado.

## Compile para produção

Quando estiver tudo pronto para distribuir seu aplicativo de tarefas, compile-o para produção:

```bash
wails3 build
```

Isso cria um executável nativo otimizado em `bin/`:

- Compila seu código Go com otimizações
- Compila seu frontend para produção
- Empacota tudo em um único executável
- O aplicativo resultante normalmente ocupa 10-20 MB (em comparação com mais de 150 MB do Electron)

Você pode executar o executável diretamente: não é necessário nenhum runtime nem iniciar servidores. É um aplicativo verdadeiramente nativo.

## O que você criou

Você acaba de criar um aplicativo de tarefas completo com:

**Implementação completa de CRUD:**

- Criação de um serviço com operações de criação, leitura, atualização e exclusão
- Adição de validação de entrada e tratamento de erros
- Compreensão de como os erros do Go se tornam exceções do JavaScript

**Gerenciamento de estado seguro para uso concorrente:**

- Uso de `sync.RWMutex` para gerenciar o acesso concorrente com segurança
- Compreensão da diferença entre bloqueios de leitura (RLock) e de escrita (Lock)
- Compreensão de como `defer` evita deadlocks ao garantir a liberação de recursos

**Interface moderna e refinada:**

- Criação de uma interface com efeito de vidro, gradientes e efeitos de desfoque
- Criação de caixas de seleção com estilo personalizado sem usar nenhum framework
- Adição de interações ao passar o cursor e transições para proporcionar a sensação de um aplicativo nativo
- Implementação de um estado vazio usando apenas CSS

**Fundamentos do Wails:**

- Registro de serviços e geração automática de bindings
- Chamada de métodos Go pelo JavaScript com async/await
- Sincronização de estado entre o Go e o frontend
- Compilação e empacotamento de um aplicativo desktop nativo

## Próximas etapas

Agora que você entende as operações CRUD e o gerenciamento de estado, experimente:

- **Adicione persistência:** Faça com que as tarefas persistam após a reinicialização do aplicativo usando SQLite ou outra solução de armazenamento
- **Adicione mais recursos:** filtragem (todas/ativas/concluídas), edição de tarefas existentes e operações em massa
- **Explore o tutorial do Notes:** Veja como as operações de arquivo funcionam em [Notes](/tutorials/03-notes-vanilla/)
- **Crie algo de verdade:** Use esses conceitos para criar seu próprio aplicativo!
