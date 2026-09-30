---
title: "Перетаскивание в HTML"
description: "Перетаскивайте элементы внутри приложения"
slug: "features/drag-and-drop/html"
sourcePath: "features/drag-and-drop/html.md"
---

Механизм перетаскивания HTML5 позволяет пользователям перетаскивать элементы в интерфейсе приложения — например, менять порядок элементов списка или перемещать их между столбцами. Это стандартная веб-функция, которая работает в Wails без какой-либо дополнительной настройки.

## Как сделать элемент перетаскиваемым

По умолчанию большинство элементов нельзя перетаскивать. Чтобы сделать элемент перетаскиваемым, добавьте `draggable="true"`:

```html
<div class="item" draggable="true">Drag me</div>
```

Теперь, когда пользователь нажмёт на элемент и начнёт его перетаскивать, будет отображаться предварительное изображение элемента.

## Определение зоны сброса

По умолчанию элементы не принимают сбрасываемые объекты. Чтобы элемент принимал их, необходимо отменить стандартное поведение для события `dragover`:

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

Необходимо вызвать `preventDefault()` для `dragover`: это указывает, что элемент принимает сбрасываемые объекты. Без этого событие сброса не сработает.

## Оформление при наведении перетаскиваемого элемента

Чтобы показать пользователям, куда можно сбросить объект, добавьте визуальную индикацию при его перемещении над зоной сброса. Событие `dragenter` срабатывает, когда объект входит в зону, а `dragleave` — когда он её покидает:

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

Примечание: `dragleave` также срабатывает при входе в дочерний элемент, из-за чего может возникать мерцание. В полном примере ниже показано, как этого избежать.

## Полный пример

Список задач, элементы которого можно перетаскивать между столбцами приоритетов. Перетаскиваемый элемент отслеживается в переменной — это самый простой подход, когда всё находится на одной странице:

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

## Совместное использование со сбросом файлов

Если приложение использует и перетаскивание HTML, и функцию [сброса файлов](/features/drag-and-drop/files/), зоны сброса HTML также будут получать события, когда пользователи перетаскивают файлы из операционной системы. Чтобы избежать неоднозначности, отфильтровывайте перетаскивание файлов в обработчиках:

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

Массив `dataTransfer.types` содержит `'Files'`, когда пользователь перетаскивает файлы из ОС, а при внутреннем перетаскивании HTML — типы вроде `'text/plain'`. Это позволяет различать эти два случая.

## Передача данных с помощью dataTransfer

В приведённом выше примере перетаскиваемый элемент отслеживается в переменной JavaScript. Этот способ хорошо подходит, когда всё находится на одной странице. Но если требуется перетаскивать объекты между iframe или передавать данные, не связанные с элементом DOM, используйте API `dataTransfer`:

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

Данные хранятся в виде строк, поэтому при необходимости объекты нужно сериализовать с помощью `JSON.stringify()`.

## Дальнейшие действия

- [Сброс файлов](/features/drag-and-drop/files/) — приём файлов из операционной системы
- [API перетаскивания MDN](https://developer.mozilla.org/en-US/docs/Web/API/HTML_Drag_and_Drop_API) — полная справочная документация по API браузера
