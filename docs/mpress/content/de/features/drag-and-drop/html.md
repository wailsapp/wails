---
title: "HTML-Drag-and-Drop"
description: "Elemente innerhalb Ihrer Anwendung ziehen und ablegen"
slug: "features/drag-and-drop/html"
sourcePath: "features/drag-and-drop/html.md"
---

Mit HTML5-Drag-and-Drop können Benutzer Elemente innerhalb der Benutzeroberfläche Ihrer App ziehen, um beispielsweise eine Liste neu anzuordnen oder Elemente zwischen Spalten zu verschieben. Diese Standard-Webfunktion funktioniert in Wails ohne besondere Einrichtung.

## Element ziehbar machen

Standardmäßig lassen sich die meisten Elemente nicht ziehen. Fügen Sie `draggable="true"` hinzu, um ein Element ziehbar zu machen:

```html
<div class="item" draggable="true">Drag me</div>
```

Wenn der Benutzer das Element anklickt und zieht, wird nun eine Ziehvorschau angezeigt.

## Ablagezone definieren

Standardmäßig akzeptieren Elemente keine Ablagevorgänge. Damit ein Element Ablagevorgänge akzeptiert, müssen Sie das Standardverhalten für `dragover` unterbinden:

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

`preventDefault()` muss für `dragover` aufgerufen werden. Dadurch wird signalisiert, dass dieses Element Ablagevorgänge akzeptiert. Andernfalls wird das Drop-Ereignis nicht ausgelöst.

## Darstellung beim Überziehen festlegen

Geben Sie beim Ziehen über eine Ablagezone visuelles Feedback, damit Benutzer erkennen, wo sie etwas ablegen können. Das Ereignis `dragenter` wird ausgelöst, wenn etwas in die Zone gelangt, und `dragleave`, wenn es sie verlässt:

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

Hinweis: `dragleave` wird auch beim Eintritt in ein untergeordnetes Element ausgelöst, was zu Flackern führen kann. Das vollständige Beispiel unten zeigt, wie Sie dies behandeln.

## Vollständiges Beispiel

Eine Aufgabenliste, deren Einträge zwischen Prioritätsspalten verschoben werden können. Dabei wird das gezogene Element in einer Variablen nachverfolgt. Dies ist der einfachste Ansatz, wenn sich alles auf derselben Seite befindet:

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

## Mit Dateiablage kombinieren

Wenn Ihre App sowohl HTML-Drag-and-Drop als auch [Dateiablage](/features/drag-and-drop/files/) verwendet, empfangen Ihre HTML-Ablagezonen auch Ereignisse, wenn Benutzer Dateien aus dem Betriebssystem hineinziehen. Filtern Sie das Ziehen von Dateien in Ihren Handlern heraus, um Verwechslungen zu vermeiden:

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

Das Array `dataTransfer.types` enthält `'Files'`, wenn der Benutzer Dateien aus dem Betriebssystem zieht, bei internen HTML-Ziehvorgängen hingegen Typen wie `'text/plain'`. So können Sie zwischen beiden unterscheiden.

## Daten mit dataTransfer übergeben

Im obigen Beispiel wird das gezogene Element in einer JavaScript-Variablen nachverfolgt. Das funktioniert gut, wenn sich alles auf derselben Seite befindet. Wenn Sie jedoch Elemente zwischen iframes ziehen oder Daten übergeben müssen, die nicht an ein DOM-Element gebunden sind, verwenden Sie die API `dataTransfer`:

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

Die Daten werden als Zeichenfolgen gespeichert. Serialisieren Sie daher Objekte bei Bedarf mit `JSON.stringify()`.

## Nächste Schritte

- [Dateiablage](/features/drag-and-drop/files/) – Dateien aus dem Betriebssystem annehmen
- [MDN-Drag-and-Drop-API](https://developer.mozilla.org/en-US/docs/Web/API/HTML_Drag_and_Drop_API) – Vollständige Referenz der Browser-API
