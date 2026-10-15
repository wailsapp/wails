---
title: "HTML 拖放"
description: "在应用程序内拖放元素"
slug: "features/drag-and-drop/html"
sourcePath: "features/drag-and-drop/html.md"
---

HTML5 拖放功能允许用户在应用的 UI 中拖动元素，例如对列表重新排序或在列之间移动项目。这是标准的 Web 功能，无需任何特殊设置即可在 Wails 中使用。

## 使元素可拖动

默认情况下，大多数元素都无法拖动。要使元素可拖动，请添加`draggable="true"`：

```html
<div class="item" draggable="true">Drag me</div>
```

现在，当用户点击并拖动该元素时，会显示拖动预览。

## 定义放置区域

默认情况下，元素不接受放置。要使元素接受放置，需要取消`dragover`的默认行为：

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

必须对`dragover`调用`preventDefault()`，这表示该元素接受放置。否则，不会触发放置事件。

## 设置拖动悬停样式

为了向用户显示可放置的位置，请在拖动到放置区域上方时提供视觉反馈。当某个对象进入该区域时，会触发`dragenter`事件；当它离开时，会触发`dragleave`事件：

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

注意：进入子元素时也会触发`dragleave`，这可能导致闪烁。下面的完整示例展示了如何处理此问题。

## 完整示例

这是一个可在优先级列之间拖动项目的任务列表。该示例使用变量跟踪正在拖动的元素；当所有内容都位于同一页面时，这是最简单的方法：

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

## 与文件放置结合使用

如果应用同时使用 HTML 拖放和[文件放置](/features/drag-and-drop/files/)，当用户从操作系统拖入文件时，HTML 放置区域也会收到事件。为避免混淆，请在处理程序中过滤掉文件拖动：

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

当用户从操作系统拖动文件时，`dataTransfer.types`数组包含`'Files'`；对于 HTML 内部拖动，该数组则包含`'text/plain'`之类的类型。这样便可区分两者。

## 使用 dataTransfer 传递数据

上面的示例使用 JavaScript 变量跟踪正在拖动的元素。当所有内容都位于同一页面时，这种方法效果很好。但如果需要在 iframe 之间拖动，或传递不与 DOM 元素关联的数据，请使用`dataTransfer` API：

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

数据以字符串形式存储，因此如有需要，必须使用`JSON.stringify()`序列化对象。

## 后续步骤

- [文件放置](/features/drag-and-drop/files/) - 接收来自操作系统的文件
- [MDN 拖放 API](https://developer.mozilla.org/en-US/docs/Web/API/HTML_Drag_and_Drop_API) - 完整的浏览器 API 参考
