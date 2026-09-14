// The component palette: a searchable, categorised list of draggable chips.
// Clicking a chip is a keyboard-friendly alternative to dragging — it appends
// the component to the selected container (or the page).

import {categories, components, getDef} from './components';
import {icon} from './icons';
import {ROOT_ID, findParent} from './model';
import {startPaletteDrag, endDrag} from './canvas';
import {store} from './store';

export function initPalette(): void {
    const search = document.getElementById('palette-search') as HTMLInputElement;
    search.addEventListener('input', () => renderPalette(search.value));
    renderPalette('');
}

function renderPalette(filter: string): void {
    const host = document.getElementById('palette')!;
    const q = filter.trim().toLowerCase();
    host.replaceChildren();

    for (const category of categories) {
        const items = components.filter((c) => c.category === category && (!q || c.label.toLowerCase().includes(q) || c.type.includes(q)));
        if (items.length === 0) {
            continue;
        }
        const group = document.createElement('section');
        group.className = 'palette-group';
        const title = document.createElement('h3');
        title.textContent = category;
        group.append(title);
        const grid = document.createElement('div');
        grid.className = 'palette-grid';
        for (const def of items) {
            const chip = document.createElement('button');
            chip.className = 'chip';
            chip.type = 'button';
            chip.draggable = true;
            chip.dataset.type = def.type;
            chip.title = `Drag “${def.label}” onto the canvas, or click to add`;
            chip.innerHTML = `<span class="chip-icon">${icon(def.icon)}</span><span class="chip-label">${def.label}</span>`;
            chip.addEventListener('dragstart', (e) => {
                startPaletteDrag(e, def.type);
                chip.classList.add('is-dragging');
                document.getElementById('canvas')!.classList.add('is-dragging');
            });
            chip.addEventListener('dragend', () => {
                chip.classList.remove('is-dragging');
                endDrag();
            });
            chip.addEventListener('click', () => addToSelection(def.type));
            grid.append(chip);
        }
        group.append(grid);
        host.append(group);
    }

    if (!host.children.length) {
        const empty = document.createElement('p');
        empty.className = 'palette-empty';
        empty.textContent = `No components match “${filter}”.`;
        host.append(empty);
    }
}

/** Append a new node next to / inside the current selection. */
function addToSelection(type: string): void {
    const selected = store.selected;
    if (selected && getDef(selected.type).container) {
        store.insertNew(type, selected.id, selected.children?.length ?? 0);
        return;
    }
    if (selected) {
        const hit = findParent(store.doc.root, selected.id);
        if (hit) {
            store.insertNew(type, hit.parent.id, hit.index + 1);
            return;
        }
    }
    store.insertNew(type, ROOT_ID, store.doc.root.children?.length ?? 0);
}
