---
title: "HTML ドラッグ＆ドロップ"
description: "アプリケーション内で要素をドラッグ＆ドロップする"
slug: "features/drag-and-drop/html"
sourcePath: "features/drag-and-drop/html.md"
---

HTML5 のドラッグ＆ドロップを使用すると、リストの並べ替えや列間での項目の移動など、アプリの UI 内で要素をドラッグできます。これは標準の Web 機能であり、特別な設定をしなくても Wails で動作します。

## 要素をドラッグ可能にする

デフォルトでは、ほとんどの要素をドラッグできません。要素をドラッグ可能にするには、`draggable="true"` を追加します。

```html
<div class="item" draggable="true">Drag me</div>
```

これで、ユーザーが要素をクリックしてドラッグすると、ドラッグプレビューが表示されます。

## ドロップゾーンを定義する

デフォルトでは、要素はドロップを受け付けません。要素がドロップを受け付けるようにするには、`dragover` でデフォルトの動作をキャンセルする必要があります。

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

`dragover` で `preventDefault()` を呼び出す必要があります。これにより、この要素がドロップを受け付けることが示されます。呼び出さない場合、drop イベントは発生しません。

## ドラッグ中のホバー表示を設定する

ドロップできる場所をユーザーに示すには、ドロップゾーン上にドラッグされたときに視覚的なフィードバックを追加します。対象がゾーンに入ると `dragenter` イベントが発生し、ゾーンから出ると `dragleave` が発生します。

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

注: `dragleave` は子要素に入ったときにも発生するため、表示がちらつくことがあります。以下の完全な例では、この問題への対処方法を示します。

## 完全な例

項目を優先度ごとの列間でドラッグできるタスクリストです。この例では、ドラッグ中の要素を変数で追跡します。すべてが同じページ上にある場合は、これが最も簡単な方法です。

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

## ファイルドロップと組み合わせる

アプリで HTML ドラッグ＆ドロップと [ファイルドロップ](/features/drag-and-drop/files/) の両方を使用している場合、ユーザーがオペレーティングシステムからファイルをドラッグしたときにも、HTML のドロップゾーンがイベントを受け取ります。混同を防ぐには、ハンドラーでファイルのドラッグを除外します。

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

ユーザーが OS からファイルをドラッグしている場合、`dataTransfer.types` 配列には `'Files'` が含まれます。一方、HTML 内部のドラッグでは、`'text/plain'` などの型が含まれます。これにより、両者を区別できます。

## dataTransfer でデータを渡す

上記の例では、ドラッグ中の要素を JavaScript 変数で追跡しています。これは、すべてが同じページ上にある場合に適しています。ただし、iframe 間でドラッグする必要がある場合や、DOM 要素に紐づかないデータを渡す場合は、`dataTransfer` API を使用します。

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

データは文字列として保存されるため、必要に応じて `JSON.stringify()` でオブジェクトをシリアライズする必要があります。

## 次のステップ

- [ファイルドロップ](/features/drag-and-drop/files/) - オペレーティングシステムからファイルを受け付ける
- [MDN Drag and Drop API](https://developer.mozilla.org/en-US/docs/Web/API/HTML_Drag_and_Drop_API) - ブラウザー API の完全なリファレンス
