/**
 * DND NPM Runtime Test
 *
 * This file tests drag-and-drop functionality using the @wailsio/runtime npm module
 * instead of the bundled /wails/runtime.js.
 *
 * The key difference:
 * - Bundled runtime: import { Events } from '/wails/runtime.js'
 * - NPM module: import { Events } from '@wailsio/runtime'
 */

import { Events } from '@wailsio/runtime';

const documentsEl = document.getElementById('documents-list');
const imagesEl = document.getElementById('images-list');
const otherEl = document.getElementById('other-list');
const dropDetails = document.getElementById('drop-details');

// ===== External File Drop =====
const imageExtensions = ['.png', '.jpg', '.jpeg', '.gif', '.bmp', '.svg', '.webp', '.ico', '.tiff', '.tif'];
const documentExtensions = ['.pdf', '.doc', '.docx', '.txt', '.rtf', '.odt', '.xls', '.xlsx', '.ppt', '.pptx', '.md', '.csv', '.json', '.xml', '.html', '.htm'];

function getFileName(path) {
    return path.split(/[/\\]/).pop();
}

function getExtension(path) {
    const name = getFileName(path);
    const idx = name.lastIndexOf('.');
    return idx > 0 ? name.substring(idx).toLowerCase() : '';
}

function categoriseFile(path) {
    const ext = getExtension(path);
    if (imageExtensions.includes(ext)) return 'images';
    if (documentExtensions.includes(ext)) return 'documents';
    return 'other';
}

function addFileToList(listEl, fileName) {
    const empty = listEl.querySelector('.empty');
    if (empty) empty.remove();

    const li = document.createElement('li');
    li.textContent = fileName;
    listEl.appendChild(li);
}

// Listen for files-dropped event from Go backend
Events.On('files-dropped', (event) => {
    const { files, details } = event.data;

    files.forEach(filePath => {
        const fileName = getFileName(filePath);
        const category = categoriseFile(filePath);

        switch (category) {
            case 'documents':
                addFileToList(documentsEl, fileName);
                break;
            case 'images':
                addFileToList(imagesEl, fileName);
                break;
            default:
                addFileToList(otherEl, fileName);
        }
    });

    let info = `External: ${files.length} file(s) dropped`;
    if (details) {
        info += ` at (${details.x}, ${details.y})`;
    }
    dropDetails.textContent = info;
});

console.log('[DND NPM Test] Initialized with @wailsio/runtime');
