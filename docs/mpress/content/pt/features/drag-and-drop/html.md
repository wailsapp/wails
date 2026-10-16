---
title: "Arrastar e soltar com HTML"
description: "Arraste e solte elementos dentro do seu aplicativo"
slug: "features/drag-and-drop/html"
sourcePath: "features/drag-and-drop/html.md"
---

O recurso de arrastar e soltar do HTML5 permite que os usuários arrastem elementos dentro da interface do seu aplicativo — por exemplo, para reordenar uma lista ou mover itens entre colunas. Essa é uma funcionalidade padrão da Web que funciona no Wails sem nenhuma configuração especial.

## Tornar um elemento arrastável

Por padrão, a maioria dos elementos não pode ser arrastada. Para tornar um elemento arrastável, adicione `draggable="true"`:

```html
<div class="item" draggable="true">Drag me</div>
```

Agora, o elemento exibirá uma prévia do arraste quando o usuário clicar nele e o arrastar.

## Definir uma zona de soltura

Por padrão, os elementos não aceitam solturas. Para fazer com que um elemento as aceite, você precisa cancelar o comportamento padrão em `dragover`:

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

É obrigatório chamar `preventDefault()` em `dragover`, pois isso sinaliza que o elemento aceita solturas. Sem essa chamada, o evento de soltura não será disparado.

## Estilizar a zona durante o arraste

Para mostrar aos usuários onde eles podem soltar um item, adicione um feedback visual quando ele for arrastado sobre uma zona de soltura. O evento `dragenter` é disparado quando algo entra na zona, e `dragleave` é disparado quando sai:

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

Observação: `dragleave` também é disparado ao entrar em um elemento filho, o que pode causar oscilação visual. O exemplo completo abaixo mostra como lidar com isso.

## Exemplo completo

Uma lista de tarefas cujos itens podem ser arrastados entre colunas de prioridade. Este exemplo mantém o elemento arrastado em uma variável, que é a abordagem mais simples quando tudo está na mesma página:

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

## Combinar com a soltura de arquivos

Se o seu aplicativo usa tanto o recurso de arrastar e soltar do HTML quanto [Soltura de arquivos](/features/drag-and-drop/files/), suas zonas de soltura do HTML também receberão eventos quando os usuários arrastarem arquivos do sistema operacional. Para evitar confusão, filtre os arrastes de arquivos nos seus manipuladores:

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

O array `dataTransfer.types` contém `'Files'` quando o usuário está arrastando arquivos do sistema operacional, mas contém tipos como `'text/plain'` nos arrastes internos do HTML. Isso permite distinguir os dois casos.

## Passar dados com dataTransfer

O exemplo acima mantém o elemento arrastado em uma variável JavaScript. Isso funciona bem quando tudo está na mesma página. Porém, se você precisar arrastar entre iframes ou passar dados que não estejam associados a um elemento DOM, use a API `dataTransfer`:

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

Os dados são armazenados como strings; portanto, se necessário, você precisará serializar objetos com `JSON.stringify()`.

## Próximas etapas

- [Soltura de arquivos](/features/drag-and-drop/files/) — Aceite arquivos do sistema operacional
- [API de arrastar e soltar da MDN](https://developer.mozilla.org/en-US/docs/Web/API/HTML_Drag_and_Drop_API) — Referência completa da API do navegador
