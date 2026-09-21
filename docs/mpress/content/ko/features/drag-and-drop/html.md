---
title: "HTML 드래그 앤 드롭"
description: "애플리케이션 내에서 요소 드래그 앤 드롭하기"
slug: "features/drag-and-drop/html"
sourcePath: "features/drag-and-drop/html.md"
---

HTML5 드래그 앤 드롭을 사용하면 목록의 순서를 바꾸거나 열 사이에서 항목을 이동하는 등 앱 UI 내에서 요소를 드래그할 수 있습니다. 이는 별도의 설정 없이 Wails에서 작동하는 표준 웹 기능입니다.

## 요소를 드래그할 수 있도록 설정하기

기본적으로 대부분의 요소는 드래그할 수 없습니다. 요소를 드래그할 수 있도록 `draggable="true"`을 추가하세요:

```html
<div class="item" draggable="true">Drag me</div>
```

이제 사용자가 요소를 클릭하여 드래그하면 드래그 미리 보기가 표시됩니다.

## 드롭 영역 정의하기

기본적으로 요소는 드롭을 허용하지 않습니다. 요소가 드롭을 허용하도록 하려면 `dragover`에서 기본 동작을 취소해야 합니다:

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

`dragover`에서 `preventDefault()`을 호출해야 합니다. 이 호출은 해당 요소가 드롭을 허용한다는 것을 나타냅니다. 호출하지 않으면 drop 이벤트가 발생하지 않습니다.

## 드래그 호버 스타일 지정하기

사용자가 드롭할 수 있는 위치를 알 수 있도록 드롭 영역 위로 항목을 드래그할 때 시각적 피드백을 추가하세요. 항목이 영역에 들어오면 `dragenter` 이벤트가 발생하고, 영역에서 벗어나면 `dragleave` 이벤트가 발생합니다:

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

참고: 자식 요소로 들어갈 때도 `dragleave` 이벤트가 발생하므로 깜박임이 생길 수 있습니다. 아래의 전체 예제에서는 이를 처리하는 방법을 보여 줍니다.

## 전체 예제

항목을 우선순위 열 사이에서 드래그할 수 있는 작업 목록입니다. 드래그 중인 요소를 변수로 추적하며, 모든 항목이 같은 페이지에 있을 때 가장 간단한 방법입니다:

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

## 파일 드롭과 함께 사용하기

앱에서 HTML 드래그 앤 드롭과 [파일 드롭](/features/drag-and-drop/files/)을 모두 사용하면 사용자가 운영 체제에서 파일을 드래그할 때 HTML 드롭 영역에도 이벤트가 전달됩니다. 혼동을 방지하려면 핸들러에서 파일 드래그를 걸러 내세요:

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

사용자가 OS에서 파일을 드래그할 때 `dataTransfer.types` 배열에는 `'Files'`이 포함되지만, 내부 HTML 드래그의 경우에는 `'text/plain'` 같은 유형이 포함됩니다. 이를 통해 두 경우를 구분할 수 있습니다.

## dataTransfer로 데이터 전달하기

위 예제에서는 드래그 중인 요소를 JavaScript 변수로 추적합니다. 모든 항목이 같은 페이지에 있을 때는 이 방법이 적합합니다. 하지만 iframe 사이에서 드래그하거나 DOM 요소와 연결되지 않은 데이터를 전달해야 한다면 `dataTransfer` API를 사용하세요:

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

데이터는 문자열로 저장되므로 필요한 경우 `JSON.stringify()`을 사용해 객체를 직렬화해야 합니다.

## 다음 단계

- [파일 드롭](/features/drag-and-drop/files/) - 운영 체제에서 파일 받기
- [MDN 드래그 앤 드롭 API](https://developer.mozilla.org/en-US/docs/Web/API/HTML_Drag_and_Drop_API) - 전체 브라우저 API 참조
