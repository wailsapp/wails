// A tiny observable store holding the document, selection, and an undo/redo
// history. Panels subscribe and re-render on change; mutations go through the
// methods below so every edit is undoable.

import {
    type ArtboardTheme, type NodeId, type UIDocument, type UINode,
    ROOT_ID, cloneNode, emptyDocument, findNode, findParent, isWithin,
} from './model';
import {createNode, getDef} from './components';

export type Device = 'desktop' | 'tablet' | 'phone';

export type ChangeKind = 'doc' | 'selection' | 'device' | 'preview' | 'file';

export interface Change {
    kind: ChangeKind;
    /** Who made the change — lets the inspector skip re-rendering its own edits. */
    origin?: string;
}

type Listener = (change: Change) => void;

const HISTORY_LIMIT = 200;
const COALESCE_MS = 700;

export class Store {
    doc: UIDocument = emptyDocument();
    selectedId: NodeId | null = null;
    device: Device = 'desktop';
    preview = false;
    /** Path of the file the layout was last saved to / opened from (via Go). */
    filePath = '';
    dirty = false;

    private undoStack: string[] = [];
    private redoStack: string[] = [];
    private lastCoalesceKey = '';
    private lastCoalesceAt = 0;
    private listeners = new Set<Listener>();

    subscribe(fn: Listener): () => void {
        this.listeners.add(fn);
        return () => this.listeners.delete(fn);
    }

    private emit(kind: ChangeKind, origin?: string): void {
        for (const fn of this.listeners) {
            fn({kind, origin});
        }
    }

    // ----- Queries -----------------------------------------------------------

    get selected(): UINode | null {
        return this.selectedId ? findNode(this.doc.root, this.selectedId) : null;
    }

    node(id: NodeId): UINode | null {
        return findNode(this.doc.root, id);
    }

    get canUndo(): boolean {
        return this.undoStack.length > 0;
    }

    get canRedo(): boolean {
        return this.redoStack.length > 0;
    }

    // ----- History -------------------------------------------------------------

    /**
     * Snapshot the document before a mutation. Passing the same coalesceKey
     * in quick succession (e.g. while typing) folds the edits into one undo step.
     */
    private checkpoint(coalesceKey = ''): void {
        const now = Date.now();
        if (coalesceKey && coalesceKey === this.lastCoalesceKey && now - this.lastCoalesceAt < COALESCE_MS) {
            this.lastCoalesceAt = now;
            return;
        }
        this.lastCoalesceKey = coalesceKey;
        this.lastCoalesceAt = now;
        this.undoStack.push(JSON.stringify(this.doc));
        if (this.undoStack.length > HISTORY_LIMIT) {
            this.undoStack.shift();
        }
        this.redoStack = [];
    }

    private commit(origin?: string): void {
        this.dirty = true;
        this.emit('doc', origin);
    }

    undo(): void {
        const prev = this.undoStack.pop();
        if (prev === undefined) {
            return;
        }
        this.redoStack.push(JSON.stringify(this.doc));
        this.doc = JSON.parse(prev) as UIDocument;
        this.lastCoalesceKey = '';
        this.ensureSelectionValid();
        this.commit('history');
    }

    redo(): void {
        const next = this.redoStack.pop();
        if (next === undefined) {
            return;
        }
        this.undoStack.push(JSON.stringify(this.doc));
        this.doc = JSON.parse(next) as UIDocument;
        this.lastCoalesceKey = '';
        this.ensureSelectionValid();
        this.commit('history');
    }

    private ensureSelectionValid(): void {
        if (this.selectedId && !findNode(this.doc.root, this.selectedId)) {
            this.selectedId = null;
            this.emit('selection');
        }
    }

    // ----- Document-level -------------------------------------------------------

    load(doc: UIDocument, filePath = ''): void {
        this.doc = doc;
        this.filePath = filePath;
        this.undoStack = [];
        this.redoStack = [];
        this.selectedId = null;
        this.dirty = false;
        this.emit('doc', 'load');
        this.emit('selection');
        this.emit('file');
    }

    reset(): void {
        this.load(emptyDocument());
    }

    rename(name: string): void {
        if (this.doc.name === name) {
            return;
        }
        this.checkpoint('rename');
        this.doc.name = name;
        this.commit('toolbar');
    }

    setTheme(theme: ArtboardTheme): void {
        this.checkpoint();
        this.doc.theme = theme;
        this.commit('toolbar');
    }

    markSaved(path: string): void {
        this.filePath = path;
        this.dirty = false;
        this.emit('file');
    }

    setDevice(device: Device): void {
        this.device = device;
        this.emit('device');
    }

    setPreview(on: boolean): void {
        this.preview = on;
        this.emit('preview');
    }

    // ----- Selection -----------------------------------------------------------

    select(id: NodeId | null): void {
        if (id === ROOT_ID) {
            id = null;
        }
        if (this.selectedId === id) {
            return;
        }
        this.selectedId = id;
        this.emit('selection');
    }

    selectParent(): void {
        if (!this.selectedId) {
            return;
        }
        const hit = findParent(this.doc.root, this.selectedId);
        this.select(hit && hit.parent.id !== ROOT_ID ? hit.parent.id : null);
    }

    selectSibling(delta: 1 | -1): void {
        if (!this.selectedId) {
            const first = this.doc.root.children?.[0];
            this.select(first ? first.id : null);
            return;
        }
        const hit = findParent(this.doc.root, this.selectedId);
        if (!hit) {
            return;
        }
        const siblings = hit.parent.children ?? [];
        const next = siblings[hit.index + delta];
        if (next) {
            this.select(next.id);
        }
    }

    // ----- Tree mutations --------------------------------------------------------

    /** Insert a brand-new node of `type` into `parentId` at `index`. Returns the new node. */
    insertNew(type: string, parentId: NodeId, index: number): UINode {
        const node = createNode(type);
        this.checkpoint();
        this.attach(node, parentId, index);
        this.commit('canvas');
        this.select(node.id);
        return node;
    }

    /** Move an existing node to a new parent / index. */
    move(id: NodeId, parentId: NodeId, index: number): void {
        if (isWithin(this.doc.root, parentId, id)) {
            return; // never drop a node into itself
        }
        const hit = findParent(this.doc.root, id);
        if (!hit) {
            return;
        }
        // Dropping onto its own slot (or the slot right after it) is a no-op.
        if (hit.parent.id === parentId && (index === hit.index || index === hit.index + 1)) {
            return;
        }
        this.checkpoint();
        const [node] = hit.parent.children!.splice(hit.index, 1);
        if (hit.parent.id === parentId && index > hit.index) {
            index -= 1;
        }
        this.attach(node, parentId, index);
        this.commit('canvas');
        this.select(id);
    }

    private attach(node: UINode, parentId: NodeId, index: number): void {
        const parent = findNode(this.doc.root, parentId);
        if (!parent) {
            return;
        }
        if (!getDef(parent.type).container) {
            return;
        }
        parent.children ??= [];
        const at = Math.max(0, Math.min(index, parent.children.length));
        parent.children.splice(at, 0, node);
    }

    remove(id: NodeId): void {
        const hit = findParent(this.doc.root, id);
        if (!hit) {
            return;
        }
        this.checkpoint();
        hit.parent.children!.splice(hit.index, 1);
        if (this.selectedId && isWithin(this.doc.root, this.selectedId, id)) {
            // Selection was inside the removed subtree: fall back to a neighbour.
            const siblings = hit.parent.children ?? [];
            const neighbour = siblings[hit.index] ?? siblings[hit.index - 1];
            this.selectedId = neighbour ? neighbour.id : (hit.parent.id === ROOT_ID ? null : hit.parent.id);
            this.emit('selection');
        }
        this.commit('canvas');
    }

    duplicate(id: NodeId): void {
        const hit = findParent(this.doc.root, id);
        if (!hit) {
            return;
        }
        this.checkpoint();
        const copy = cloneNode(hit.parent.children![hit.index]);
        hit.parent.children!.splice(hit.index + 1, 0, copy);
        this.commit('canvas');
        this.select(copy.id);
    }

    /** Nudge a node one slot up or down among its siblings. */
    shift(id: NodeId, delta: 1 | -1): void {
        const hit = findParent(this.doc.root, id);
        if (!hit) {
            return;
        }
        const siblings = hit.parent.children!;
        const to = hit.index + delta;
        if (to < 0 || to >= siblings.length) {
            return;
        }
        this.checkpoint();
        const [node] = siblings.splice(hit.index, 1);
        siblings.splice(to, 0, node);
        this.commit('canvas');
    }

    setProp(id: NodeId, key: string, value: unknown, origin = 'inspector'): void {
        const node = findNode(this.doc.root, id);
        if (!node || node.props[key] === value) {
            return;
        }
        this.checkpoint(`prop:${id}:${key}`);
        node.props[key] = value;
        this.commit(origin);
    }

    resetProps(id: NodeId): void {
        const node = findNode(this.doc.root, id);
        if (!node) {
            return;
        }
        this.checkpoint();
        node.props = {...getDef(node.type).defaults};
        this.commit('inspector-reset');
    }
}

export const store = new Store();
