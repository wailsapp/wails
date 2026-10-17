---
title: "Glisser-déposer HTML"
description: "Faites glisser et déposez des éléments dans votre application"
slug: "features/drag-and-drop/html"
sourcePath: "features/drag-and-drop/html.md"
---

Le glisser-déposer HTML5 permet aux utilisateurs de déplacer des éléments dans l’interface de votre application, par exemple pour réorganiser une liste ou déplacer des éléments entre des colonnes. Cette fonctionnalité web standard fonctionne dans Wails sans configuration particulière.

## Rendre un élément déplaçable

Par défaut, la plupart des éléments ne peuvent pas être déplacés. Pour rendre un élément déplaçable, ajoutez `draggable="true"` :

```html
<div class="item" draggable="true">Drag me</div>
```

L’élément affiche désormais un aperçu de déplacement lorsque l’utilisateur clique dessus et le fait glisser.

## Définir une zone de dépôt

Par défaut, les éléments n’acceptent pas les dépôts. Pour qu’un élément les accepte, vous devez annuler le comportement par défaut sur `dragover` :

```html
<div class="drop-zone" id="target">Drop here</div>

<script>
const target = document.getElementById('target');

target.addEventListener('dragover', (e) => {
    e.preventDefault(); // Allow the drop
});

target.addEventListener('drop', (e) => {
    e.preventDefault();
    // Handle the drop
});
</script>
```

Vous devez appeler `preventDefault()` sur `dragover` : cela indique que cet élément accepte les dépôts. Sans cet appel, l’événement de dépôt ne se déclenchera pas.

## Mettre en évidence la zone survolée

Pour indiquer aux utilisateurs où ils peuvent effectuer le dépôt, ajoutez un retour visuel lors du survol d’une zone de dépôt. L’événement `dragenter` se déclenche lorsqu’un élément entre dans la zone, et `dragleave` lorsqu’il en sort :

```css
.drop-zone {
    border: 2px dashed #ccc;
    padding: 40px;
    transition: all 0.2s ease;
}

.drop-zone.drag-over {
    border-color: #007bff;
    background-color: rgba(0, 123, 255, 0.1);
}
```

```javascript
const target = document.getElementById('target');

target.addEventListener('dragenter', () => {
    target.classList.add('drag-over');
});

target.addEventListener('dragleave', () => {
    target.classList.remove('drag-over');
});

target.addEventListener('drop', (e) => {
    e.preventDefault();
    target.classList.remove('drag-over');
    // Handle the drop
});
```

Remarque : `dragleave` se déclenche également lors de l’entrée dans un élément enfant, ce qui peut provoquer un scintillement. L’exemple complet ci-dessous montre comment gérer ce comportement.

## Exemple complet

Voici une liste de tâches dont les éléments peuvent être déplacés entre des colonnes de priorité. Cet exemple conserve l’élément déplacé dans une variable, ce qui constitue l’approche la plus simple lorsque tout se trouve sur la même page :

```html
<div class="tasks">
    <div class="item" draggable="true">Fix login bug</div>
    <div class="item" draggable="true">Update docs</div>
    <div class="item" draggable="true">Add dark mode</div>
</div>

<div class="columns">
    <div class="drop-zone" data-priority="high">
        <h3>High Priority</h3>
        <ul></ul>
    </div>
    <div class="drop-zone" data-priority="low">
        <h3>Low Priority</h3>
        <ul></ul>
    </div>
</div>

<script>
let draggedItem = null;

// Track which item is being dragged
document.querySelectorAll('.item').forEach(item => {
    item.addEventListener('dragstart', () => {
        draggedItem = item;
        item.classList.add('dragging');
    });
    
    item.addEventListener('dragend', () => {
        item.classList.remove('dragging');
    });
});

// Handle drops on each zone
document.querySelectorAll('.drop-zone').forEach(zone => {
    zone.addEventListener('dragover', (e) => {
        e.preventDefault();
    });
    
    zone.addEventListener('dragenter', () => {
        zone.classList.add('drag-over');
    });
    
    zone.addEventListener('dragleave', (e) => {
        // Only remove the class if we're leaving the zone entirely,
        // not just entering a child element
        if (!zone.contains(e.relatedTarget)) {
            zone.classList.remove('drag-over');
        }
    });
    
    zone.addEventListener('drop', (e) => {
        e.preventDefault();
        zone.classList.remove('drag-over');
        
        if (draggedItem) {
            const li = document.createElement('li');
            li.textContent = draggedItem.textContent;
            zone.querySelector('ul').appendChild(li);
            draggedItem.remove();
        }
    });
});
</script>

<style>
.item {
    padding: 12px 16px;
    background: #f0f0f0;
    margin: 8px 0;
    border-radius: 8px;
    cursor: grab;
}

.item.dragging {
    opacity: 0.5;
}

.drop-zone {
    min-height: 150px;
    border: 2px dashed #ccc;
    border-radius: 8px;
    padding: 15px;
    transition: all 0.2s ease;
}

.drop-zone.drag-over {
    border-color: #007bff;
    background: rgba(0, 123, 255, 0.1);
}
</style>
```

## Combiner avec le dépôt de fichiers

Si votre application utilise à la fois le glisser-déposer HTML et le [dépôt de fichiers](/features/drag-and-drop/files/), vos zones de dépôt HTML recevront aussi des événements lorsque les utilisateurs y feront glisser des fichiers depuis le système d’exploitation. Pour éviter toute confusion, excluez les glissements de fichiers dans vos gestionnaires :

```javascript
zone.addEventListener('dragenter', (e) => {
    // Ignore external file drags
    if (e.dataTransfer?.types.includes('Files')) return;
    
    zone.classList.add('drag-over');
});

zone.addEventListener('dragover', (e) => {
    // Ignore external file drags
    if (e.dataTransfer?.types.includes('Files')) return;
    
    e.preventDefault();
});

zone.addEventListener('drop', (e) => {
    // Ignore external file drags
    if (e.dataTransfer?.types.includes('Files')) return;
    
    e.preventDefault();
    zone.classList.remove('drag-over');
    // Handle the internal drop
});
```

Le tableau `dataTransfer.types` contient `'Files'` lorsque l’utilisateur fait glisser des fichiers depuis le système d’exploitation, mais des types tels que `'text/plain'` pour les glissements HTML internes. Vous pouvez ainsi distinguer les deux.

## Transmettre des données avec dataTransfer

L’exemple ci-dessus conserve l’élément déplacé dans une variable JavaScript. Cette méthode fonctionne bien lorsque tout se trouve sur la même page. Toutefois, si vous devez effectuer un glisser-déposer entre des iframes ou transmettre des données qui ne sont pas liées à un élément DOM, utilisez l’API `dataTransfer` :

```javascript
// When drag starts, store data
item.addEventListener('dragstart', (e) => {
    e.dataTransfer.setData('text/plain', item.id);
});

// When dropped, retrieve the data
target.addEventListener('drop', (e) => {
    e.preventDefault();
    const itemId = e.dataTransfer.getData('text/plain');
    const item = document.getElementById(itemId);
    // Move or copy the item
});
```

Les données sont stockées sous forme de chaînes de caractères. Si nécessaire, vous devrez donc sérialiser les objets avec `JSON.stringify()`.

## Étapes suivantes

- [Dépôt de fichiers](/features/drag-and-drop/files/) – Accepter les fichiers provenant du système d’exploitation
- [API de glisser-déposer MDN](https://developer.mozilla.org/en-US/docs/Web/API/HTML_Drag_and_Drop_API) – Référence complète de l’API du navigateur
