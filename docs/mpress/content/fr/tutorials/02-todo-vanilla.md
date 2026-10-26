---
title: "Liste de tâches"
description: "Créez une application complète de liste de tâches avec des opérations CRUD"
slug: "tutorials/02-todo-vanilla"
sourcePath: "tutorials/02-todo-vanilla.md"
---

Dans ce tutoriel, vous allez créer une application de liste de tâches entièrement fonctionnelle. Ce tutoriel va plus loin que celui consacré au service de codes QR : vous apprendrez à gérer l’état, à traiter plusieurs opérations et à créer une interface utilisateur soignée.

**Ce que vous allez créer :**

- Une application complète de liste de tâches permettant d’ajouter, d’achever et de supprimer des tâches
- Une gestion de l’état sûre pour les accès concurrents (importante pour les applications de bureau)
- Une interface moderne au design glassmorphique
- Le tout en JavaScript natif, sans framework

**Ce que vous allez apprendre :**

- Les opérations CRUD (création, lecture, mise à jour et suppression)
- La gestion sûre d’un état mutable en Go
- Le traitement et la validation des saisies utilisateur
- La création d’interfaces utilisateur réactives offrant une expérience native

![Application de liste de tâches](/assets/todo-app.png)

**Temps nécessaire :** 20 minutes

## Créez votre projet

@steps
### Générez le projet
Commencez par créer un projet Wails. Nous utiliserons le modèle JavaScript natif par défaut, qui fournit une base propre :

```bash
wails3 init -n todo-app -t vanilla-js
cd todo-app
```

Cette commande crée un projet avec la structure de base : le backend Go à la racine et le code frontend dans le répertoire `frontend/`.

### Créez le service TODO
Le service TODO gérera l’état de notre application et fournira les méthodes nécessaires aux opérations CRUD. Contrairement à un serveur web, où chaque requête est isolée, une application de bureau peut exécuter plusieurs opérations simultanément. Nous devons donc assurer une gestion de l’état sûre pour les accès concurrents.

Supprimez `greetservice.go` et créez un fichier nommé `todoservice.go` :

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

**Fonctionnement :**

**La structure `Todo` :**

- Définit la forme de nos données avec les champs ID, Title et Completed
- Les balises `json:` indiquent à Go comment convertir cette structure en JSON pour le frontend
- Chaque champ est exporté (son nom commence par une majuscule) afin que le générateur de liaisons puisse le détecter

**La structure `TodoService` :**

- `todos []Todo` : une tranche contenant toutes nos tâches
- `nextID int` : conserve le prochain ID à attribuer (simule une auto-incrémentation)
- `mu sync.RWMutex` : un mutex en lecture-écriture garantissant des accès concurrents sûrs

**Sécurité des accès concurrents avec `sync.RWMutex` :**

- L’interface utilisateur d’une application de bureau peut lancer plusieurs opérations simultanément
- `RLock()` autorise plusieurs lecteurs simultanés (par exemple, plusieurs appels à `GetAll`)
- `Lock()` accorde un accès exclusif aux écritures (par exemple, `Add`, `Toggle` et `Delete`)
- `defer` garantit la libération des verrous même si la fonction se termine prématurément ou déclenche une panique

**Les méthodes :**

- `GetAll()` : renvoie toutes les tâches (utilise un verrou de lecture puisque les données ne sont pas modifiées)
- `Add(title)` : crée une tâche, valide la saisie et incrémente l’ID
- `Toggle(id)` : inverse l’état d’achèvement d’une tâche
- `Delete(id)` : supprime une tâche de la tranche

**Gestion des erreurs :**

- Conformément aux conventions de Go, nous renvoyons `error` en dernière valeur
- Les titres vides sont refusés
- Les opérations portant sur des tâches inexistantes renvoient des erreurs
- Ces erreurs deviennent des exceptions JavaScript dans le frontend

### Mettez à jour main.go
Enregistrez le service TODO auprès de votre application Wails. Dans `main.go`, repérez la section `Services` et remplacez GreetService par notre TodoService :

```go {title="main.go" highlight="5"}
Services: []application.Service{
    application.NewService(NewTodoService()),
},
```

**Fonctionnement :**

- Nous supprimons le GreetService par défaut et ajoutons notre TodoService à la place
- `application.NewService()` encapsule notre service afin que Wails puisse le gérer
- Wails générera automatiquement les liaisons JavaScript de toutes les méthodes publiques de ce service

### Créez l’interface utilisateur frontend
Créons maintenant le frontend. C’est ici que nous appellerons nos méthodes Go et afficherons l’interface utilisateur. Nous utilisons JavaScript natif pour rester simples et vous montrer directement le fonctionnement des liaisons.

Remplacez `frontend/src/main.js` :

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

**Fonctionnement :**

**Importation des liaisons :**

- `import {TodoService} from "../bindings/changeme"` : importe les liaisons Go générées automatiquement
- Remarque : `changeme` correspondra au véritable nom de votre module défini dans `go.mod`

**La fonction `loadTodos()` :**

- Appelle `TodoService.GetAll()` pour récupérer toutes les tâches depuis Go
- Génère le HTML de chaque tâche à l’aide de littéraux de gabarit
- Ajoute ou supprime dynamiquement la classe `completed` pour appliquer le style
- Utilise les attributs `onclick` pour relier les boutons à nos fonctions
- Concatène tout le HTML et l’insère dans le DOM

**Les fonctions CRUD :**

- `addTodo()` : valide la saisie, appelle la méthode Go `Add` et actualise la liste
- `toggleTodo(id)` : appelle la méthode Go `Toggle` et actualise la liste
- `deleteTodo(id)` : appelle la méthode Go `Delete` et actualise la liste
- Toutes les fonctions sont asynchrones, car les appels à Go renvoient des promesses

**Pourquoi les associer à window :**

- `window.addTodo = ...` rend les fonctions accessibles depuis les attributs HTML `onclick`
- Il s’agit d’un modèle simple pour JavaScript natif (les frameworks gèrent cela différemment)
- En production, vous pourriez utiliser à la place une véritable délégation d’événements

**Le modèle d’actualisation :**

- Après chaque modification (ajout, basculement ou suppression), nous appelons de nouveau `loadTodos()`
- Cela garantit que l’interface utilisateur reste synchronisée avec l’état Go
- Autre possibilité : faire renvoyer le nouvel état par les méthodes Go afin d’éviter le second appel

### Mettez à jour le HTML
Le HTML fournit la structure de notre application de tâches. Il est minimaliste et sémantique : toute la magie opère dans le JavaScript et le CSS.

Remplacez `frontend/index.html` :

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

**Voici ce qui se passe :**

**La structure :**

- `container` : centre notre application et en limite la largeur
- `card` : la carte blanche principale qui contient tous les éléments
- `input-box` : conteneur flexible pour le champ de saisie et le bouton Ajouter
- `todo-list` : emplacement où JavaScript insérera chaque tâche

**Gestion des événements :**

- `onkeypress="if(event.key==='Enter') addTodo()"` : ajoute une tâche lorsque vous appuyez sur Entrée
- `onclick="addTodo()"` : ajoute une tâche lorsque vous cliquez sur le bouton
- Les gestionnaires d’événements en ligne conviennent bien aux applications simples en JavaScript natif

**Script de module :**

- `<script type="module">` nous permet d’utiliser les imports ES6
- Notre fichier `main.js` peut importer les bindings et utiliser du JavaScript moderne

### Mettre en forme l’application
Le CSS crée un design moderne inspiré du verre dépoli, avec des transitions fluides. L’objectif est d’obtenir une finition soignée qui rende l’application agréable à utiliser.

Remplacez `frontend/public/style.css` :

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

**Voici ce qui se passe :**

**Design inspiré du verre dépoli :**

- `background: linear-gradient(135deg, #667eea 0%, #764ba2 100%)` : arrière-plan avec un dégradé violet
- `backdrop-filter: blur(10px)` : crée l’effet de verre dépoli sur la carte
- `rgba(255, 255, 255, 0.95)` : blanc semi-transparent pour produire l’effet de verre

**Style personnalisé des cases à cocher :**

- `appearance: none` supprime l’apparence par défaut de la case à cocher du navigateur
- Nous créons un carré arrondi personnalisé avec une coche à l’aide de `::after`
- La coche apparaît lorsque `checked`, grâce au caractère Unicode ✓

**Interactions au survol :**

- Les tâches se décalent vers la droite au survol (`transform: translateX(4px)`)
- Le bouton de suppression reste masqué jusqu’au survol (`opacity: 0` → `opacity: 1`)
- Les boutons s’agrandissent légèrement au survol pour fournir un retour tactile

**État vide :**

- `#todo-list:empty::before` affiche un message lorsqu’il n’y a aucune tâche
- Solution entièrement en CSS : aucun JavaScript nécessaire

### Exécuter l’application
Voyons-la fonctionner ! Lancez le serveur de développement :

```bash
wails3 dev
```

L’application sera compilée et s’ouvrira. Essayez-la :

- Saisissez une tâche, puis appuyez sur Entrée ou cliquez sur Ajouter
- Cliquez sur la case à cocher pour marquer la tâche comme terminée
- Survolez une tâche pour faire apparaître le bouton de suppression
- Remarquez que l’interface se met à jour instantanément : notre modèle d’actualisation fonctionne

**Voici ce qui se passe :**

- Wails a généré automatiquement les bindings des méthodes de votre TodoService
- Le mode développement inclut le rechargement à chaud : modifiez le CSS et observez la mise à jour
- Votre code Go s’exécute nativement : aucune traduction ni interprétation n’est nécessaire

@end

## Fonctionnement

### Gestion de l’état sûre pour les accès concurrents

Le `sync.RWMutex` permet un accès concurrent sécurisé :

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

**Pourquoi est-ce important ?**

- Plusieurs appels du frontend peuvent avoir lieu simultanément
- Les opérations de lecture ne se bloquent pas mutuellement
- Les opérations d’écriture bénéficient d’un accès exclusif
- `defer` garantit que les verrous sont toujours libérés

### Gestion des erreurs

Le service renvoie des erreurs pour les opérations non valides :

```go
func (t *TodoService) Add(title string) (*Todo, error) {
    if title == "" {
        return nil, errors.New("title cannot be empty")
    }
    // ...
}
```

Dans le frontend, vous pouvez intercepter ces erreurs :

```javascript
try {
    await TodoService.Add(title);
} catch (err) {
    alert('Error: ' + err);
}
```

### Synchronisation de l’état

Après chaque modification, nous rechargeons la liste complète :

```javascript
window.addTodo = async () => {
    await TodoService.Add(title);  // Mutation
    await loadTodos();             // Refresh
}
```

**Autre approche :** renvoyez la liste mise à jour depuis chaque méthode afin d’éviter le second appel.

## Améliorations

### Ajouter des statistiques

Ajoutez ceci à `todoservice.go` :

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

Affichez-les dans le frontend :

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

### Ajouter « Effacer les tâches terminées »

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

### Ajouter la persistance

Pour les applications en production, ajoutez la persistance à l’aide de la solution de stockage la mieux adaptée à votre application, telle que SQLite, PostgreSQL ou un service hébergé.

## Compiler pour la production

Lorsque vous êtes prêt à distribuer votre application de tâches, compilez-la pour la production :

```bash
wails3 build
```

Cette commande crée un exécutable natif optimisé dans `bin/` :

- Compile votre code Go avec des optimisations
- Compile votre frontend pour la production
- Regroupe tous les éléments dans un seul exécutable
- L’application obtenue occupe généralement 10-20 Mo, contre plus de 150 Mo pour Electron

Vous pouvez exécuter directement le fichier : aucun environnement d’exécution n’est nécessaire et aucun serveur ne doit être démarré. Il s’agit d’une véritable application native.

## Ce que vous avez créé

Vous venez de créer une application de tâches complète avec :

**Implémentation CRUD complète :**

- Création d’un service avec des opérations de création, de lecture, de mise à jour et de suppression
- Ajout de la validation des entrées et de la gestion des erreurs
- Découverte de la conversion des erreurs Go en exceptions JavaScript

**Gestion de l’état sûre pour les accès concurrents :**

- Utilisation de `sync.RWMutex` pour gérer les accès concurrents en toute sécurité
- Compréhension de la différence entre les verrous en lecture (RLock) et les verrous en écriture (Lock)
- Observation de la manière dont `defer` évite les interblocages en garantissant le nettoyage

**Interface moderne et soignée :**

- Création d’une interface inspirée du verre dépoli, avec des dégradés et des effets de flou
- Création de cases à cocher au style personnalisé sans aucun framework
- Ajout d’interactions au survol et de transitions pour une sensation d’application native
- Implémentation d’un état vide entièrement en CSS

**Principes fondamentaux de Wails :**

- Enregistrement des services et génération automatique des liaisons
- Appel de méthodes Go depuis JavaScript avec async/await
- Synchronisation de l’état entre Go et le frontend
- Compilation et empaquetage d’une application de bureau native

## Étapes suivantes

Maintenant que vous maîtrisez les opérations CRUD et la gestion de l’état, essayez ceci :

- **Ajoutez la persistance :** conservez les tâches après le redémarrage de l’application avec SQLite ou une autre solution de stockage
- **Ajoutez d’autres fonctionnalités :** filtrage (toutes/actives/terminées), modification des tâches existantes et opérations par lots
- **Explorez le tutoriel Notes :** découvrez le fonctionnement des opérations sur les fichiers dans [Notes](/tutorials/03-notes-vanilla/)
- **Créez une véritable application :** mettez ces concepts en pratique pour créer votre propre application !
