---
title: "Seret \u0026 Lepas HTML"
description: "Seret dan lepas elemen di dalam aplikasi Anda"
slug: "features/drag-and-drop/html"
sourcePath: "features/drag-and-drop/html.md"
---

Fitur seret dan lepas HTML5 memungkinkan pengguna menyeret elemen di dalam UI aplikasi Anda—misalnya, untuk mengurutkan ulang daftar atau memindahkan item antarkolom. Ini adalah fungsionalitas web standar yang dapat digunakan di Wails tanpa penyiapan khusus.

## Membuat Elemen Dapat Diseret

Secara default, sebagian besar elemen tidak dapat diseret. Agar elemen dapat diseret, tambahkan `draggable="true"`:

```html
<div class="item" draggable="true">Drag me</div>
```

Elemen tersebut kini akan menampilkan pratinjau penyeretan saat pengguna mengeklik lalu menyeretnya.

## Menentukan Zona Pelepasan

Secara default, elemen tidak menerima pelepasan. Agar elemen menerima pelepasan, Anda perlu membatalkan perilaku default pada `dragover`:

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

Anda wajib memanggil `preventDefault()` pada `dragover`—pemanggilan ini menandakan bahwa elemen tersebut menerima pelepasan. Tanpanya, peristiwa pelepasan tidak akan dipicu.

## Menata Tampilan Saat Seretan Melintas

Untuk menunjukkan tempat pengguna dapat melepaskan item, tambahkan umpan balik visual saat item diseret di atas zona pelepasan. Peristiwa `dragenter` dipicu saat sesuatu memasuki zona tersebut, sedangkan `dragleave` dipicu saat sesuatu meninggalkannya:

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

Catatan: `dragleave` juga dipicu saat memasuki elemen turunan, yang dapat menyebabkan tampilan berkedip. Contoh lengkap di bawah menunjukkan cara menanganinya.

## Contoh Lengkap

Daftar tugas yang itemnya dapat diseret antarkolom prioritas. Contoh ini melacak elemen yang sedang diseret dalam sebuah variabel, yang merupakan pendekatan paling sederhana jika semuanya berada di halaman yang sama:

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

## Menggabungkan dengan Pelepasan File

Jika aplikasi Anda menggunakan fitur seret dan lepas HTML sekaligus [Pelepasan File](/features/drag-and-drop/files/), zona pelepasan HTML Anda juga akan menerima peristiwa saat pengguna menyeret file dari sistem operasi. Untuk menghindari kebingungan, saring penyeretan file dalam handler Anda:

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

Array `dataTransfer.types` berisi `'Files'` saat pengguna menyeret file dari OS, tetapi berisi tipe seperti `'text/plain'` untuk penyeretan HTML internal. Dengan demikian, Anda dapat membedakan keduanya.

## Meneruskan Data dengan dataTransfer

Contoh di atas melacak elemen yang sedang diseret dalam sebuah variabel JavaScript. Cara ini efektif jika semuanya berada di halaman yang sama. Namun, jika Anda perlu menyeret antar-iframe atau meneruskan data yang tidak terikat pada elemen DOM, gunakan API `dataTransfer`:

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

Data disimpan sebagai string, jadi jika diperlukan, Anda harus melakukan serialisasi objek dengan `JSON.stringify()`.

## Langkah Berikutnya

- [Pelepasan File](/features/drag-and-drop/files/)—Terima file dari sistem operasi
- [API Seret dan Lepas MDN](https://developer.mozilla.org/en-US/docs/Web/API/HTML_Drag_and_Drop_API)—Referensi lengkap API browser
