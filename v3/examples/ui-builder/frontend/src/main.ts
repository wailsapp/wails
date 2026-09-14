import {System, WML} from "@wailsio/runtime";
import {LayoutService} from "../bindings/github.com/wailsapp/wails/v3/examples/ui-builder";

import {initCanvas, renderCanvas, applySelection, scrollNodeIntoView} from './canvas';
import {initPalette} from './palette';
import {initInspector} from './inspector';
import {initLayers} from './layers';
import {renderToHTML} from './exporter';
import {starterDocument} from './starter';
import {type UIDocument, countNodes} from './model';
import {type Device, store} from './store';

const AUTOSAVE_KEY = 'uib:autosave';
const DEVICE_WIDTHS: Record<Device, number> = {desktop: 1200, tablet: 820, phone: 390};

// Wire up data-wml-openURL links (logo + footer "Docs" link).
WML.Enable();

// Show the actual Wails version this project was generated against.
document.getElementById('version')!.innerText = "v3.0.0-dev";

// The top bar doubles as the window's drag region on macOS (hidden-inset
// title bar), so leave room for the traffic lights.
if (System.IsMac()) {
    document.body.classList.add('is-mac');
}

const $ = <T extends HTMLElement = HTMLElement>(id: string): T => document.getElementById(id) as T;

// ---------------------------------------------------------------------------
// Toast + status

let toastTimer: ReturnType<typeof setTimeout>;
function toast(message: string, label = 'From Go', tone: 'ok' | 'error' = 'ok'): void {
    $('toast-label').textContent = label;
    $('toast-msg').textContent = message;
    const t = $('toast');
    t.classList.toggle('is-error', tone === 'error');
    t.classList.add('is-visible');
    clearTimeout(toastTimer);
    toastTimer = setTimeout(() => t.classList.remove('is-visible'), 3600);
}

function refreshStatus(): void {
    const n = countNodes(store.doc.root);
    const file = store.filePath ? store.filePath.split(/[\\/]/).pop() : 'not saved to disk';
    $('status-text').textContent = `${n} component${n === 1 ? '' : 's'} · ${file}${store.dirty && store.filePath ? ' · edited' : ''}`;
    $('doc-status').classList.toggle('is-dirty', store.dirty);
    $<HTMLButtonElement>('undo').disabled = !store.canUndo;
    $<HTMLButtonElement>('redo').disabled = !store.canRedo;
    document.title = `${store.doc.name || 'Untitled'} — UI Builder`;
}

// ---------------------------------------------------------------------------
// Persistence: autosave to this browser, save/open/export through Go

let autosaveTimer: ReturnType<typeof setTimeout>;
function scheduleAutosave(): void {
    clearTimeout(autosaveTimer);
    autosaveTimer = setTimeout(() => {
        try {
            localStorage.setItem(AUTOSAVE_KEY, JSON.stringify({doc: store.doc, filePath: store.filePath}));
        } catch {
            // Storage may be unavailable (private mode); the design stays in memory.
        }
    }, 400);
}

function restoreAutosave(): boolean {
    try {
        const raw = localStorage.getItem(AUTOSAVE_KEY);
        if (!raw) {
            return false;
        }
        const saved = JSON.parse(raw) as { doc: UIDocument; filePath?: string };
        if (saved?.doc?.version === 1 && saved.doc.root) {
            store.load(saved.doc, saved.filePath ?? '');
            return true;
        }
    } catch {
        // Corrupt autosave: fall through to the starter document.
    }
    return false;
}

async function save(saveAs = false): Promise<void> {
    try {
        const path = await LayoutService.SaveLayout(saveAs ? '' : store.filePath, store.doc.name, JSON.stringify(store.doc, null, 2));
        if (!path) {
            return; // cancelled
        }
        store.markSaved(path);
        toast(`Saved ${path.split(/[\\/]/).pop()}`);
    } catch (err) {
        toast(String(err), 'Save failed', 'error');
    }
}

async function open(): Promise<void> {
    try {
        const file = await LayoutService.OpenLayout();
        if (!file) {
            return; // cancelled
        }
        const doc = JSON.parse(file.data) as UIDocument;
        if (doc?.version !== 1 || !doc.root) {
            throw new Error('That file is not a UI Builder layout.');
        }
        store.load(doc, file.path);
        toast(`Opened ${file.path.split(/[\\/]/).pop()}`);
    } catch (err) {
        toast(String(err), 'Open failed', 'error');
    }
}

async function exportHTML(): Promise<void> {
    const html = renderToHTML(store.doc);
    try {
        const path = await LayoutService.ExportHTML(store.doc.name, html);
        if (path) {
            toast(`Exported ${path.split(/[\\/]/).pop()}`);
        }
    } catch (err) {
        // Outside the desktop app (plain `vite dev` in a browser) the Go side
        // isn't there — fall back to the clipboard so the export is still usable.
        try {
            await navigator.clipboard.writeText(html);
            toast('Go backend unavailable — HTML copied to the clipboard instead.', 'Export', 'error');
        } catch {
            toast(String(err), 'Export failed', 'error');
        }
    }
}

// ---------------------------------------------------------------------------
// Toolbar

function initToolbar(): void {
    const name = $<HTMLInputElement>('doc-name');
    name.addEventListener('change', () => store.rename(name.value.trim() || 'Untitled'));
    name.addEventListener('keydown', (e) => {
        if (e.key === 'Enter') {
            name.blur();
        }
    });

    for (const btn of $('device-switch').querySelectorAll<HTMLButtonElement>('[data-device]')) {
        btn.addEventListener('click', () => store.setDevice(btn.dataset.device as Device));
    }
    $('theme-toggle').addEventListener('click', () => store.setTheme(store.doc.theme === 'light' ? 'dark' : 'light'));
    $('undo').addEventListener('click', () => store.undo());
    $('redo').addEventListener('click', () => store.redo());
    $('open').addEventListener('click', () => void open());
    $('save').addEventListener('click', () => void save());
    $('export').addEventListener('click', () => void exportHTML());
    $('preview').addEventListener('click', () => store.setPreview(!store.preview));

    for (const tab of document.querySelectorAll<HTMLButtonElement>('.panel-tabs .tab')) {
        tab.addEventListener('click', () => {
            for (const t of document.querySelectorAll('.panel-tabs .tab')) {
                t.classList.toggle('is-active', t === tab);
                t.setAttribute('aria-selected', String(t === tab));
            }
            for (const page of document.querySelectorAll<HTMLElement>('.tab-page')) {
                page.classList.toggle('is-active', page.dataset.page === tab.dataset.tab);
            }
        });
    }

    $('shortcuts-close').addEventListener('click', () => ($('shortcuts').hidden = true));
    $('shortcuts').addEventListener('click', (e) => {
        if (e.target === e.currentTarget) {
            $('shortcuts').hidden = true;
        }
    });
}

function applyDevice(): void {
    const artboard = $('artboard');
    artboard.dataset.device = store.device;
    $('artboard-size').textContent = `${DEVICE_WIDTHS[store.device]} × auto`;
    for (const btn of $('device-switch').querySelectorAll<HTMLButtonElement>('[data-device]')) {
        const active = btn.dataset.device === store.device;
        btn.classList.toggle('is-active', active);
        btn.setAttribute('aria-checked', String(active));
    }
}

function applyPreview(): void {
    document.body.classList.toggle('is-preview', store.preview);
    $('preview').classList.toggle('is-active', store.preview);
    renderCanvas();
}

// ---------------------------------------------------------------------------
// Keyboard shortcuts

function initKeyboard(): void {
    document.addEventListener('keydown', (e) => {
        const mod = e.metaKey || e.ctrlKey;
        const inField = (e.target as HTMLElement).matches('input, textarea, select, [contenteditable]');

        if (mod && e.key.toLowerCase() === 's') { e.preventDefault(); void save(e.shiftKey); return; }
        if (mod && e.key.toLowerCase() === 'o') { e.preventDefault(); void open(); return; }
        if (mod && e.key.toLowerCase() === 'e') { e.preventDefault(); void exportHTML(); return; }
        if (mod && e.key.toLowerCase() === 'p') { e.preventDefault(); store.setPreview(!store.preview); return; }
        if (e.key === 'Escape') {
            if (!$('shortcuts').hidden) { $('shortcuts').hidden = true; return; }
            if (store.preview) { store.setPreview(false); return; }
            if (inField) { (e.target as HTMLElement).blur(); return; }
            store.select(null);
            return;
        }
        if (inField) {
            // Let text fields keep their native editing keys (including ⌘Z).
            return;
        }
        if (mod && e.key.toLowerCase() === 'z') { e.preventDefault(); e.shiftKey ? store.redo() : store.undo(); return; }
        if (mod && e.key.toLowerCase() === 'y') { e.preventDefault(); store.redo(); return; }
        if (e.key === '?') { $('shortcuts').hidden = !$('shortcuts').hidden; return; }

        const id = store.selectedId;
        if (!id) {
            if (e.key === 'ArrowDown') { store.selectSibling(1); }
            return;
        }
        if (mod && e.key.toLowerCase() === 'd') { e.preventDefault(); store.duplicate(id); return; }
        if (e.key === 'Backspace' || e.key === 'Delete') { e.preventDefault(); store.remove(id); return; }
        if (e.altKey && e.key === 'ArrowUp') { e.preventDefault(); store.shift(id, -1); return; }
        if (e.altKey && e.key === 'ArrowDown') { e.preventDefault(); store.shift(id, 1); return; }
        if (e.key === 'ArrowUp') { e.preventDefault(); store.selectSibling(-1); return; }
        if (e.key === 'ArrowDown') { e.preventDefault(); store.selectSibling(1); return; }
        if (e.shiftKey && e.key === 'Tab') { e.preventDefault(); store.selectParent(); }
    });
}

// ---------------------------------------------------------------------------
// Boot

initToolbar();
initKeyboard();
initCanvas();
initPalette();
initInspector();
initLayers();

store.subscribe((change) => {
    switch (change.kind) {
        case 'doc':
            renderCanvas();
            $<HTMLInputElement>('doc-name').value = store.doc.name;
            scheduleAutosave();
            break;
        case 'selection':
            applySelection();
            if (store.selectedId) {
                scrollNodeIntoView(store.selectedId);
            }
            break;
        case 'device':
            applyDevice();
            break;
        case 'preview':
            applyPreview();
            break;
        case 'file':
            scheduleAutosave();
            break;
    }
    refreshStatus();
});

if (!restoreAutosave()) {
    store.load(starterDocument());
}
applyDevice();
refreshStatus();
