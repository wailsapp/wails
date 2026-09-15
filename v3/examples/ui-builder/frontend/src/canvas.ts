// The design canvas: renders the document tree into the artboard, decorates
// every node so it can be selected and dragged, and handles drops from both
// the palette (new nodes) and the canvas itself (moves).

import {type NodeId, type UINode, ROOT_ID, findParent, isWithin} from './model';
import {getDef} from './components';
import {componentStyles} from './theme';
import {store} from './store';
import {type Cleanup, wireRuntime} from './runtime';

interface DragState {
    kind: 'new' | 'move';
    type?: string;
    id?: NodeId;
}

interface DropTarget {
    parentId: NodeId;
    index: number;
    /** Where to draw the indicator. */
    rect: DOMRect;
    axis: 'row' | 'column';
    /** True when dropping *into* an empty container (draw a highlight, not a line). */
    into: boolean;
}

const NEW_MIME = 'application/x-uib-new';
const MOVE_MIME = 'application/x-uib-move';

let drag: DragState | null = null;
let target: DropTarget | null = null;
let hoverEl: HTMLElement | null = null;

export function startPaletteDrag(e: DragEvent, type: string): void {
    drag = {kind: 'new', type};
    e.dataTransfer?.setData(NEW_MIME, type);
    e.dataTransfer!.effectAllowed = 'copy';
}

export function endDrag(): void {
    drag = null;
    target = null;
    hideIndicator();
    canvasEl().classList.remove('is-dragging');
}

const canvasEl = (): HTMLElement => document.getElementById('canvas')!;
const stageEl = (): HTMLElement => document.getElementById('stage')!;
const indicatorEl = (): HTMLElement => document.getElementById('drop-indicator')!;
const emptyEl = (): HTMLElement => document.getElementById('stage-empty')!;

// ---------------------------------------------------------------------------
// Rendering

let styleTag: HTMLStyleElement | null = null;
let unwire: Cleanup | null = null;

export function renderCanvas(): void {
    const canvas = canvasEl();
    if (!styleTag) {
        styleTag = document.createElement('style');
        styleTag.textContent = componentStyles;
        document.head.append(styleTag);
    }
    unwire?.();
    unwire = null;
    canvas.replaceChildren(renderNode(store.doc.root));
    const root = canvas.firstElementChild as HTMLElement;
    root.dataset.theme = store.doc.theme;
    const artboard = document.getElementById('artboard')!;
    artboard.dataset.theme = store.doc.theme;
    artboard.dataset.chrome = store.doc.chrome ?? 'mac';
    document.getElementById('artboard-title')!.textContent = store.doc.name || 'Untitled';
    emptyEl().hidden = (store.doc.root.children?.length ?? 0) > 0;
    applySelection();
    if (store.preview) {
        // Run mode: buttons call Go, bound elements listen for Go events.
        unwire = wireRuntime(root);
    }
}

function renderNode(node: UINode): HTMLElement {
    const def = getDef(node.type);
    const children = (node.children ?? []).map(renderNode);
    const element = def.render(node, children);
    element.dataset.nodeId = node.id;
    element.dataset.nodeType = node.type;
    element.dataset.nodeLabel = def.label;
    element.classList.add('uib-node');
    if (def.container) {
        element.classList.add('uib-container');
        element.dataset.axis = def.direction ?? 'column';
        if (children.length === 0 && node.id !== ROOT_ID) {
            element.classList.add('uib-empty');
        }
    }
    if (node.id !== ROOT_ID) {
        element.draggable = !store.preview;
    }
    return element;
}

export function applySelection(): void {
    const canvas = canvasEl();
    for (const el of canvas.querySelectorAll('.is-selected')) {
        el.classList.remove('is-selected');
    }
    if (store.selectedId) {
        const el = nodeElement(store.selectedId);
        if (el) {
            el.classList.add('is-selected');
        }
    }
}

export function nodeElement(id: NodeId): HTMLElement | null {
    return canvasEl().querySelector<HTMLElement>(`[data-node-id="${CSS.escape(id)}"]`);
}

export function scrollNodeIntoView(id: NodeId): void {
    nodeElement(id)?.scrollIntoView({block: 'nearest', behavior: 'smooth'});
}

// ---------------------------------------------------------------------------
// Drop-target maths

/** Work out where a drop at (x, y) over `over` would land. */
function computeTarget(over: HTMLElement, x: number, y: number): DropTarget | null {
    const id = over.dataset.nodeId!;
    if (drag?.kind === 'move' && drag.id && isWithin(store.doc.root, id, drag.id)) {
        return null; // can't drop a node into itself
    }
    const def = getDef(over.dataset.nodeType!);
    const rect = over.getBoundingClientRect();

    if (def.container) {
        const axis = def.direction ?? 'column';
        const kids = Array.from(over.querySelectorAll<HTMLElement>(':scope > .uib-node, :scope > .ub-inner > .uib-node'));
        if (kids.length === 0) {
            return {parentId: id, index: 0, rect, axis, into: true};
        }
        // Inside a populated container, land between children based on the pointer.
        // The root and sections take the whole pointer band; nested containers only
        // claim the inner 60% so their edges still act as "before / after me" zones
        // for the parent container.
        const inset = id === ROOT_ID ? 0 : (axis === 'column' ? rect.height * 0.2 : rect.width * 0.2);
        const inside = axis === 'column'
            ? y > rect.top + inset && y < rect.bottom - inset
            : x > rect.left + inset && x < rect.right - inset;
        if (inside || id === ROOT_ID) {
            let index = kids.length;
            for (const [i, kid] of kids.entries()) {
                const r = kid.getBoundingClientRect();
                const mid = axis === 'column' ? r.top + r.height / 2 : r.left + r.width / 2;
                if ((axis === 'column' ? y : x) < mid) {
                    index = i;
                    break;
                }
            }
            const lineRect = edgeRect(kids, index, axis, rect);
            return {parentId: id, index, rect: lineRect, axis, into: false};
        }
    }

    // A leaf (or the edge band of a container): insert before / after it in its parent.
    const hit = findParent(store.doc.root, id);
    if (!hit) {
        return null;
    }
    const parentDef = getDef(hit.parent.type);
    const axis = parentDef.direction ?? 'column';
    const after = axis === 'column' ? y > rect.top + rect.height / 2 : x > rect.left + rect.width / 2;
    const index = hit.index + (after ? 1 : 0);
    const line = new DOMRect(
        axis === 'column' ? rect.left : (after ? rect.right : rect.left),
        axis === 'column' ? (after ? rect.bottom : rect.top) : rect.top,
        axis === 'column' ? rect.width : 0,
        axis === 'column' ? 0 : rect.height,
    );
    return {parentId: hit.parent.id, index, rect: line, axis, into: false};
}

/** A zero-thickness rect along the edge where a new child at `index` would go. */
function edgeRect(kids: HTMLElement[], index: number, axis: 'row' | 'column', container: DOMRect): DOMRect {
    const ref = index < kids.length ? kids[index].getBoundingClientRect() : kids[kids.length - 1].getBoundingClientRect();
    const before = index < kids.length;
    if (axis === 'column') {
        return new DOMRect(ref.left, before ? ref.top : ref.bottom, ref.width || container.width, 0);
    }
    return new DOMRect(before ? ref.left : ref.right, ref.top, 0, ref.height || container.height);
}

function showIndicator(t: DropTarget): void {
    const ind = indicatorEl();
    const stage = stageEl().getBoundingClientRect();
    ind.classList.add('is-visible');
    ind.classList.toggle('is-into', t.into);
    ind.classList.toggle('is-vertical', !t.into && t.axis === 'row');
    if (t.into) {
        ind.style.left = `${t.rect.left - stage.left}px`;
        ind.style.top = `${t.rect.top - stage.top}px`;
        ind.style.width = `${t.rect.width}px`;
        ind.style.height = `${t.rect.height}px`;
        return;
    }
    if (t.axis === 'column') {
        ind.style.left = `${t.rect.left - stage.left}px`;
        ind.style.top = `${t.rect.top - stage.top - 2}px`;
        ind.style.width = `${t.rect.width}px`;
        ind.style.height = '4px';
    } else {
        ind.style.left = `${t.rect.left - stage.left - 2}px`;
        ind.style.top = `${t.rect.top - stage.top}px`;
        ind.style.width = '4px';
        ind.style.height = `${t.rect.height}px`;
    }
}

function hideIndicator(): void {
    indicatorEl().classList.remove('is-visible', 'is-into', 'is-vertical');
}

// ---------------------------------------------------------------------------
// Event wiring (delegated, one listener per event on the stage)

export function initCanvas(): void {
    const canvas = canvasEl();
    const stage = stageEl();

    canvas.addEventListener('click', (e) => {
        if (store.preview) {
            return;
        }
        const el = (e.target as HTMLElement).closest<HTMLElement>('.uib-node');
        store.select(el && el.dataset.nodeId !== ROOT_ID ? el.dataset.nodeId! : null);
        // Links inside the design shouldn't navigate while editing.
        if ((e.target as HTMLElement).closest('a')) {
            e.preventDefault();
        }
    });
    canvas.addEventListener('dblclick', (e) => {
        // Double-click a container to select it directly, even when a child was hit.
        const el = (e.target as HTMLElement).closest<HTMLElement>('.uib-container');
        if (el && el.dataset.nodeId !== ROOT_ID) {
            store.select(el.dataset.nodeId!);
        }
    });

    canvas.addEventListener('mouseover', (e) => {
        if (store.preview) {
            return;
        }
        const el = (e.target as HTMLElement).closest<HTMLElement>('.uib-node');
        if (el === hoverEl) {
            return;
        }
        hoverEl?.classList.remove('is-hover');
        hoverEl = el && el.dataset.nodeId !== ROOT_ID ? el : null;
        hoverEl?.classList.add('is-hover');
    });
    canvas.addEventListener('mouseleave', () => {
        hoverEl?.classList.remove('is-hover');
        hoverEl = null;
    });

    canvas.addEventListener('dragstart', (e) => {
        const el = (e.target as HTMLElement).closest<HTMLElement>('.uib-node');
        if (!el || el.dataset.nodeId === ROOT_ID || store.preview) {
            e.preventDefault();
            return;
        }
        e.stopPropagation();
        drag = {kind: 'move', id: el.dataset.nodeId};
        e.dataTransfer?.setData(MOVE_MIME, el.dataset.nodeId!);
        e.dataTransfer!.effectAllowed = 'move';
        store.select(el.dataset.nodeId!);
        // Defer so the browser captures the drag image before we fade the source.
        requestAnimationFrame(() => {
            el.classList.add('is-drag-source');
            canvas.classList.add('is-dragging');
        });
    });

    stage.addEventListener('dragover', (e) => {
        if (!drag) {
            // A drag that started outside the app (e.g. a file) — ignore.
            const types = e.dataTransfer?.types ?? [];
            if (!types.includes(NEW_MIME) && !types.includes(MOVE_MIME)) {
                return;
            }
        }
        e.preventDefault();
        e.dataTransfer!.dropEffect = drag?.kind === 'move' ? 'move' : 'copy';
        const over = (e.target as HTMLElement).closest<HTMLElement>('.uib-node') ?? canvas.querySelector<HTMLElement>('.uib-node');
        if (!over) {
            target = null;
            hideIndicator();
            return;
        }
        target = computeTarget(over, e.clientX, e.clientY);
        if (target) {
            showIndicator(target);
        } else {
            hideIndicator();
        }
    });

    stage.addEventListener('dragleave', (e) => {
        if (!stage.contains(e.relatedTarget as Node | null)) {
            hideIndicator();
            target = null;
        }
    });

    stage.addEventListener('drop', (e) => {
        e.preventDefault();
        const t = target;
        const d = drag ?? recoverDrag(e);
        endDrag();
        if (!t || !d) {
            return;
        }
        if (d.kind === 'new' && d.type) {
            store.insertNew(d.type, t.parentId, t.index);
        } else if (d.kind === 'move' && d.id) {
            store.move(d.id, t.parentId, t.index);
        }
    });

    canvas.addEventListener('dragend', () => {
        for (const el of canvas.querySelectorAll('.is-drag-source')) {
            el.classList.remove('is-drag-source');
        }
        endDrag();
    });

    // Keep the drop indicator honest while the artboard scrolls under a drag.
    document.getElementById('stage-scroll')!.addEventListener('scroll', () => {
        if (target) {
            hideIndicator();
        }
    }, {passive: true});
}

/** Rebuild drag state from the dataTransfer if the module-level state was lost. */
function recoverDrag(e: DragEvent): DragState | null {
    const type = e.dataTransfer?.getData(NEW_MIME);
    if (type) {
        return {kind: 'new', type};
    }
    const id = e.dataTransfer?.getData(MOVE_MIME);
    if (id) {
        return {kind: 'move', id};
    }
    return null;
}
