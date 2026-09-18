import './style.css';
import './app.css';

// CRITICAL: Register global handlers IMMEDIATELY to prevent file drops from opening new windows
// This must be done as early as possible, before any other code runs
(function() {
    // Helper function to check if drag event contains files
    function isFileDrop(e) {
        return e.dataTransfer && e.dataTransfer.types &&
               (e.dataTransfer.types.includes('Files') ||
                Array.from(e.dataTransfer.types).includes('Files'));
    }

    // Global dragover handler - MUST prevent default for file drops
    window.addEventListener('dragover', function(e) {
        if (isFileDrop(e)) {
            e.preventDefault();
            e.dataTransfer.dropEffect = 'copy';
        }
    }, true); // Use capture phase to handle before any other handlers

    // Global drop handler - MUST prevent default for file drops
    window.addEventListener('drop', function(e) {
        if (isFileDrop(e)) {
            e.preventDefault();
            console.log('Global handler prevented file drop navigation');
        }
    }, true); // Use capture phase to handle before any other handlers

    // Global dragleave handler
    window.addEventListener('dragleave', function(e) {
        if (isFileDrop(e)) {
            e.preventDefault();
        }
    }, true); // Use capture phase

    console.log('Global file drop prevention handlers registered');
})();

document.querySelector('#app').innerHTML = `
    <h1>Wails Drag & Drop Test</h1>

    <div class="compact-container">
        <div class="drag-source">
            <h4>HTML5 Source</h4>
            <div class="draggable" draggable="true" data-item="Item 1">Item 1</div>
            <div class="draggable" draggable="true" data-item="Item 2">Item 2</div>
            <div class="draggable" draggable="true" data-item="Item 3">Item 3</div>
        </div>

        <div class="drop-zone" id="dropZone">
            <h4>HTML5 Drop</h4>
            <p id="dropMessage">Drop here</p>
        </div>

        <div class="file-drop-zone" id="fileDropZone">
            <h4>File Drop</h4>
            <p id="fileDropMessage">Drop files here</p>
        </div>
    </div>

    <div class="status">
        <h4>Event Log</h4>
        <div id="eventLog"></div>
    </div>
`;

// Get all draggable items and drop zones
const draggables = document.querySelectorAll('.draggable');
const dropZone = document.getElementById('dropZone');
const fileDropZone = document.getElementById('fileDropZone');
const eventLog = document.getElementById('eventLog');
const dropMessage = document.getElementById('dropMessage');
const fileDropMessage = document.getElementById('fileDropMessage');

let draggedItem = null;
let eventCounter = 0;

// Function to log events
function logEvent(eventName, details = '') {
    eventCounter++;
    const timestamp = new Date().toLocaleTimeString();
    const logEntry = document.createElement('div');
    logEntry.className = `log-entry ${eventName.replace(' ', '-').toLowerCase()}`;
    logEntry.textContent = `[${timestamp}] ${eventCounter}. ${eventName} ${details}`;
    eventLog.insertBefore(logEntry, eventLog.firstChild);

    // Keep only last 20 events
    while (eventLog.children.length > 20) {
        eventLog.removeChild(eventLog.lastChild);
    }

    console.log(`Event: ${eventName} ${details}`);
}

// Add event listeners to draggable items
draggables.forEach(item => {
    // Drag start
    item.addEventListener('dragstart', (e) => {
        draggedItem = e.target;
        e.target.classList.add('dragging');
        e.dataTransfer.effectAllowed = 'copy';
        e.dataTransfer.setData('text/plain', e.target.dataset.item);
        logEvent('drag-start', `- Started dragging: ${e.target.dataset.item}`);
    });

    // Drag end
    item.addEventListener('dragend', (e) => {
        e.target.classList.remove('dragging');
        logEvent('drag-end', `- Ended dragging: ${e.target.dataset.item}`);
    });
});

// Add event listeners to HTML drop zone
dropZone.addEventListener('dragenter', (e) => {
    e.preventDefault();
    dropZone.classList.add('drag-over');
    logEvent('drag-enter', '- Entered HTML drop zone');
});

dropZone.addEventListener('dragover', (e) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = 'copy';
    // Don't log every dragover to avoid spam
});

dropZone.addEventListener('dragleave', (e) => {
    if (e.target === dropZone) {
        dropZone.classList.remove('drag-over');
        logEvent('drag-leave', '- Left HTML drop zone');
    }
});

dropZone.addEventListener('drop', (e) => {
    e.preventDefault();
    dropZone.classList.remove('drag-over');

    const data = e.dataTransfer.getData('text/plain');
    logEvent('drop', `- Dropped in HTML zone: ${data}`);

    if (draggedItem) {
        // Create a copy of the dragged item
        const droppedElement = document.createElement('div');
        droppedElement.className = 'dropped-item';
        droppedElement.textContent = data;

        // Remove the placeholder message if it exists
        if (dropMessage) {
            dropMessage.style.display = 'none';
        }

        dropZone.appendChild(droppedElement);
    }

    draggedItem = null;
});

// Add event listeners to file drop zone
fileDropZone.addEventListener('dragenter', (e) => {
    e.preventDefault();
    fileDropZone.classList.add('drag-over');
    logEvent('drag-enter', '- Entered file drop zone');
});

fileDropZone.addEventListener('dragover', (e) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = 'copy';
});

fileDropZone.addEventListener('dragleave', (e) => {
    if (e.target === fileDropZone) {
        fileDropZone.classList.remove('drag-over');
        logEvent('drag-leave', '- Left file drop zone');
    }
});

fileDropZone.addEventListener('drop', (e) => {
    e.preventDefault();
    fileDropZone.classList.remove('drag-over');

    const files = [...e.dataTransfer.files];
    if (files.length > 0) {
        logEvent('file-drop', `- Dropped ${files.length} file(s)`);

        // Hide the placeholder message
        if (fileDropMessage) {
            fileDropMessage.style.display = 'none';
        }

        // Display dropped files
        files.forEach(file => {
            const fileElement = document.createElement('div');
            fileElement.className = 'dropped-file';

            // Format file size
            let size = file.size;
            let unit = 'bytes';
            if (size > 1024 * 1024) {
                size = (size / (1024 * 1024)).toFixed(2);
                unit = 'MB';
            } else if (size > 1024) {
                size = (size / 1024).toFixed(2);
                unit = 'KB';
            }

            fileElement.textContent = `📄 ${file.name} (${size} ${unit})`;
            fileDropZone.appendChild(fileElement);
        });
    }
});

// Log when page loads
window.addEventListener('DOMContentLoaded', () => {
    logEvent('page-loaded', '- Wails drag-and-drop test page ready');
    console.log('Wails Drag and Drop test application loaded');

    // Check if Wails drag and drop is enabled
    if (window.wails && window.wails.flags) {
        logEvent('wails-status', `- Wails DnD enabled: ${window.wails.flags.enableWailsDragAndDrop}`);
    }

    // IMPORTANT: Register Wails drag-and-drop handlers to prevent browser navigation
    // This will ensure external files don't open in new windows when dropped anywhere
    if (window.runtime && window.runtime.OnFileDrop) {
        window.runtime.OnFileDrop((x, y, paths) => {
            logEvent('wails-file-drop', `- Wails received ${paths.length} file(s) at (${x}, ${y})`);
            console.log('Wails OnFileDrop:', paths);
        }, false); // false = don't require drop target, handle all file drops
        logEvent('wails-setup', '- Wails OnFileDrop handlers registered');
    }
});