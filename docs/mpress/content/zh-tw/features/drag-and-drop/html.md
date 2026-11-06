---
title: "HTML 拖放"
description: "在應用程式內拖放元素"
slug: "features/drag-and-drop/html"
sourcePath: "features/drag-and-drop/html.md"
---

HTML5 拖放功能可讓使用者在應用程式的使用者介面中拖曳元素，例如重新排序清單，或在欄之間移動項目。這是標準的網頁功能，無須任何特殊設定即可在 Wails 中運作。

## 讓元素可拖曳

預設情況下，大多數元素都無法拖曳。若要讓元素可拖曳，請新增`draggable="true"`：

```html
<div class="item" draggable="true">Drag me</div>
```

現在，使用者按住並拖曳此元素時，元素會顯示拖曳預覽。

## 定義放置區域

預設情況下，元素不接受放置。若要讓元素接受放置，必須取消`dragover`的預設行為：

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

必須在`dragover`上呼叫`preventDefault()`，這表示此元素接受放置。若未呼叫，便不會觸發放置事件。

## 設定拖曳懸停樣式

若要向使用者顯示可放置的位置，請在項目拖曳至放置區域上方時提供視覺回饋。當項目進入區域時會觸發`dragenter`事件，離開時則會觸發`dragleave`：

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

注意：進入子元素時也會觸發`dragleave`，這可能造成閃爍。下方的完整範例示範了處理方式。

## 完整範例

以下是一個可在不同優先順序欄之間拖曳項目的工作清單。此範例使用變數追蹤正在拖曳的元素；當所有內容都位於同一個頁面時，這是最簡單的方法：

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

## 與檔案放置功能搭配使用

如果應用程式同時使用 HTML 拖放和[檔案放置](/features/drag-and-drop/files/)，當使用者從作業系統拖曳檔案時，HTML 放置區域也會收到事件。為避免混淆，請在事件處理常式中排除檔案拖曳：

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

當使用者從作業系統拖曳檔案時，`dataTransfer.types`陣列會包含`'Files'`；若是 HTML 內部拖曳，則會包含`'text/plain'`之類的型別。您可以藉此區分這兩種情況。

## 使用 dataTransfer 傳遞資料

上述範例使用 JavaScript 變數追蹤正在拖曳的元素。當所有內容都位於同一個頁面時，這種方式相當實用。但如果需要在 iframe 之間拖曳，或傳遞未繫結至 DOM 元素的資料，請使用`dataTransfer` API：

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

資料會以字串形式儲存，因此如有需要，必須使用`JSON.stringify()`將物件序列化。

## 後續步驟

- [檔案放置](/features/drag-and-drop/files/)－接受來自作業系統的檔案
- [MDN 拖放 API](https://developer.mozilla.org/en-US/docs/Web/API/HTML_Drag_and_Drop_API)－完整的瀏覽器 API 參考資料
