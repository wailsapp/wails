// The layers panel: an indented tree of the document. Click selects, drag
// re-parents (it reuses the canvas drop logic by dispatching through the store).

import {getDef} from './components';
import {icon} from './icons';
import {type UINode, ROOT_ID, findParent, isWithin} from './model';
import {store} from './store';
import {scrollNodeIntoView} from './canvas';

const collapsed = new Set<string>();
let dragId: string | null = null;

export function initLayers(): void {
    store.subscribe((change) => {
        if (change.kind === 'doc' || change.kind === 'selection') {
            renderLayers();
        }
    });
    renderLayers();
}

export function renderLayers(): void {
    const host = document.getElementById('layers')!;
    const list = document.createElement('div');
    list.className = 'layer-tree';
    const children = store.doc.root.children ?? [];
    if (children.length === 0) {
        const p = document.createElement('p');
        p.className = 'insp-note';
        p.textContent = 'The page is empty. Drop something onto the canvas to see it here.';
        host.replaceChildren(p);
        return;
    }
    for (const child of children) {
        appendRow(list, child, 0);
    }
    // Drop at the very end of the page.
    const tail = document.createElement('div');
    tail.className = 'layer-tail';
    wireDrop(tail, () => ({parentId: ROOT_ID, index: children.length}));
    list.append(tail);
    host.replaceChildren(list);
}

function appendRow(list: HTMLElement, node: UINode, depth: number): void {
    const def = getDef(node.type);
    const row = document.createElement('div');
    row.className = 'layer-row' + (node.id === store.selectedId ? ' is-selected' : '');
    row.style.setProperty('--depth', String(depth));
    row.draggable = true;
    row.dataset.nodeId = node.id;

    const hasKids = (node.children?.length ?? 0) > 0;
    const twisty = document.createElement('button');
    twisty.type = 'button';
    twisty.className = 'layer-twisty' + (hasKids ? '' : ' is-hidden') + (collapsed.has(node.id) ? ' is-collapsed' : '');
    twisty.innerHTML = icon('down');
    twisty.addEventListener('click', (e) => {
        e.stopPropagation();
        if (collapsed.has(node.id)) {
            collapsed.delete(node.id);
        } else {
            collapsed.add(node.id);
        }
        renderLayers();
    });

    const label = document.createElement('span');
    label.className = 'layer-label';
    label.innerHTML = `${icon(def.icon)}<span>${def.label}</span>`;
    const summary = summarize(node);
    if (summary) {
        const s = document.createElement('em');
        s.textContent = summary;
        label.append(s);
    }
    row.append(twisty, label);

    row.addEventListener('click', () => {
        store.select(node.id);
        scrollNodeIntoView(node.id);
    });
    row.addEventListener('dragstart', (e) => {
        e.stopPropagation();
        dragId = node.id;
        e.dataTransfer!.effectAllowed = 'move';
        row.classList.add('is-dragging');
    });
    row.addEventListener('dragend', () => {
        dragId = null;
        row.classList.remove('is-dragging');
        clearDropClasses();
    });
    // Drop before / after / into this row depending on the pointer's vertical position.
    row.addEventListener('dragover', (e) => {
        if (!dragId || isWithin(store.doc.root, node.id, dragId)) {
            return;
        }
        e.preventDefault();
        clearDropClasses();
        row.classList.add(zoneClass(rowZone(row, node, e.clientY)));
    });
    row.addEventListener('dragleave', () => clearDropClasses());
    row.addEventListener('drop', (e) => {
        e.preventDefault();
        e.stopPropagation();
        if (!dragId) {
            return;
        }
        const zone = rowZone(row, node, e.clientY);
        const id = dragId;
        dragId = null;
        clearDropClasses();
        if (zone === 'into') {
            store.move(id, node.id, node.children?.length ?? 0);
            return;
        }
        const hit = findParent(store.doc.root, node.id);
        if (hit) {
            store.move(id, hit.parent.id, hit.index + (zone === 'after' ? 1 : 0));
        }
    });

    list.append(row);
    if (hasKids && !collapsed.has(node.id)) {
        for (const child of node.children!) {
            appendRow(list, child, depth + 1);
        }
    }
}

type Zone = 'before' | 'after' | 'into';

function rowZone(row: HTMLElement, node: UINode, y: number): Zone {
    const r = row.getBoundingClientRect();
    const rel = (y - r.top) / r.height;
    if (getDef(node.type).container && rel > 0.3 && rel < 0.7) {
        return 'into';
    }
    return rel < 0.5 ? 'before' : 'after';
}

const zoneClass = (z: Zone): string => `is-drop-${z}`;

function clearDropClasses(): void {
    for (const el of document.querySelectorAll('.layer-row.is-drop-before, .layer-row.is-drop-after, .layer-row.is-drop-into, .layer-tail.is-drop-into')) {
        el.classList.remove('is-drop-before', 'is-drop-after', 'is-drop-into');
    }
}

function wireDrop(el: HTMLElement, place: () => { parentId: string; index: number }): void {
    el.addEventListener('dragover', (e) => {
        if (!dragId) {
            return;
        }
        e.preventDefault();
        el.classList.add('is-drop-into');
    });
    el.addEventListener('dragleave', () => el.classList.remove('is-drop-into'));
    el.addEventListener('drop', (e) => {
        e.preventDefault();
        if (!dragId) {
            return;
        }
        const {parentId, index} = place();
        store.move(dragId, parentId, index);
        dragId = null;
        el.classList.remove('is-drop-into');
    });
}

/** A short preview of the node's main text prop, shown dimmed after the label. */
function summarize(node: UINode): string {
    const text = ['text', 'label', 'brand', 'value'].map((k) => node.props[k]).find((v) => typeof v === 'string' && v.trim());
    if (!text) {
        return '';
    }
    const s = String(text).replace(/\s+/g, ' ').trim();
    return s.length > 28 ? s.slice(0, 27) + '…' : s;
}
