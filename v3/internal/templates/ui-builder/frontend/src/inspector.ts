// The inspector: a property form for the selected node, generated from the
// component definition's `fields`. Edits are pushed straight into the store
// (coalesced into single undo steps while typing).

import {type PropField, getDef} from './components';
import {icon} from './icons';
import {pathTo} from './model';
import {store} from './store';

const GROUPS: NonNullable<PropField['group']>[] = ['Content', 'Wails', 'Style', 'Layout'];

export function initInspector(): void {
    store.subscribe((change) => {
        if (change.kind === 'selection') {
            renderInspector();
        } else if (change.kind === 'doc' && change.origin !== 'inspector') {
            // Undo / redo / canvas edits may have changed the selected node's props.
            renderInspector();
        }
    });
    renderInspector();
}

export function renderInspector(): void {
    const host = document.getElementById('inspector')!;
    const node = store.selected;
    if (!node) {
        host.replaceChildren(emptyState());
        return;
    }
    const def = getDef(node.type);
    const frag = document.createDocumentFragment();

    // Header: breadcrumb + type + actions
    const head = document.createElement('div');
    head.className = 'insp-head';
    const crumbs = document.createElement('div');
    crumbs.className = 'insp-crumbs';
    for (const ancestor of pathTo(store.doc.root, node.id).filter((n) => n.id !== 'root')) {
        const b = document.createElement('button');
        b.type = 'button';
        b.textContent = getDef(ancestor.type).label;
        b.addEventListener('click', () => store.select(ancestor.id));
        crumbs.append(b, chevron());
    }
    const current = document.createElement('span');
    current.className = 'insp-current';
    current.innerHTML = `${icon(def.icon)}<strong>${def.label}</strong>`;
    crumbs.append(current);
    head.append(crumbs);

    const actions = document.createElement('div');
    actions.className = 'insp-actions';
    actions.append(
        actionButton('parent', 'Select parent (⇧⇥)', () => store.selectParent()),
        actionButton('up', 'Move up (⌥↑)', () => store.shift(node.id, -1)),
        actionButton('down', 'Move down (⌥↓)', () => store.shift(node.id, 1)),
        actionButton('copy', 'Duplicate (⌘D)', () => store.duplicate(node.id)),
        actionButton('trash', 'Delete (⌫)', () => store.remove(node.id), 'is-danger'),
    );
    head.append(actions);
    frag.append(head);

    if (def.fields.length === 0) {
        const p = document.createElement('p');
        p.className = 'insp-note';
        p.textContent = 'This component has no editable properties.';
        frag.append(p);
    }

    for (const group of GROUPS) {
        const fields = def.fields.filter((f) => (f.group ?? 'Content') === group && (!f.when || f.when(node.props)));
        if (fields.length === 0) {
            continue;
        }
        const section = document.createElement('section');
        section.className = 'insp-group';
        if (group === 'Wails') {
            section.classList.add('insp-group-wails');
        }
        const h = document.createElement('h4');
        h.innerHTML = group === 'Wails' ? `${icon('go')}<span>Go backend</span>` : `<span>${group}</span>`;
        section.append(h);
        for (const field of fields) {
            section.append(renderField(node.id, field, node.props[field.key]));
        }
        frag.append(section);
    }

    const reset = document.createElement('button');
    reset.type = 'button';
    reset.className = 'insp-reset';
    reset.textContent = 'Reset to defaults';
    reset.addEventListener('click', () => store.resetProps(node.id));
    frag.append(reset);

    host.replaceChildren(frag);
}

function emptyState(): HTMLElement {
    const d = document.createElement('div');
    d.className = 'insp-empty';
    d.innerHTML = `${icon('cursor')}<strong>Nothing selected</strong><span>Click a component on the canvas to edit its properties. Double-click to grab its container.</span>`;
    return d;
}

function chevron(): HTMLElement {
    const s = document.createElement('span');
    s.className = 'insp-chev';
    s.innerHTML = icon('chevron');
    return s;
}

function actionButton(name: string, title: string, onClick: () => void, extra = ''): HTMLButtonElement {
    const b = document.createElement('button');
    b.type = 'button';
    b.className = `icon-btn ${extra}`.trim();
    b.title = title;
    b.innerHTML = icon(name);
    b.addEventListener('click', onClick);
    return b;
}

function renderField(id: string, field: PropField, value: unknown): HTMLElement {
    const row = document.createElement('label');
    row.className = `field field-${field.kind}`;
    const label = document.createElement('span');
    label.className = 'field-label';
    label.textContent = field.label;
    row.append(label);

    const set = (v: unknown): void => {
        store.setProp(id, field.key, v);
        if (field.rerender) {
            renderInspector();
        }
    };

    switch (field.kind) {
        case 'text': {
            const input = document.createElement('input');
            input.type = 'text';
            input.value = value == null ? '' : String(value);
            input.placeholder = field.placeholder ?? '';
            input.addEventListener('input', () => set(input.value));
            row.append(input);
            break;
        }
        case 'textarea': {
            const ta = document.createElement('textarea');
            ta.value = value == null ? '' : String(value);
            ta.rows = 3;
            ta.placeholder = field.placeholder ?? '';
            ta.addEventListener('input', () => set(ta.value));
            row.append(ta);
            break;
        }
        case 'select': {
            const select = document.createElement('select');
            for (const opt of field.options ?? []) {
                const o = document.createElement('option');
                o.value = opt.value;
                o.textContent = opt.label;
                o.selected = opt.value === String(value ?? '');
                select.append(o);
            }
            select.addEventListener('change', () => set(select.value));
            row.append(wrapSelect(select));
            break;
        }
        case 'align': {
            const seg = document.createElement('div');
            seg.className = 'segmented';
            for (const opt of field.options ?? []) {
                const b = document.createElement('button');
                b.type = 'button';
                b.className = 'seg' + (opt.value === String(value ?? '') ? ' is-active' : '');
                b.textContent = opt.label;
                b.addEventListener('click', () => {
                    for (const sib of seg.children) {
                        sib.classList.toggle('is-active', sib === b);
                    }
                    set(opt.value);
                });
                seg.append(b);
            }
            row.append(seg);
            break;
        }
        case 'number': {
            const input = document.createElement('input');
            input.type = 'number';
            input.value = value == null ? '' : String(value);
            if (field.min !== undefined) input.min = String(field.min);
            if (field.max !== undefined) input.max = String(field.max);
            if (field.step !== undefined) input.step = String(field.step);
            input.addEventListener('input', () => set(Number(input.value)));
            row.append(input);
            break;
        }
        case 'range': {
            const wrap = document.createElement('div');
            wrap.className = 'range-wrap';
            const input = document.createElement('input');
            input.type = 'range';
            input.min = String(field.min ?? 0);
            input.max = String(field.max ?? 100);
            input.step = String(field.step ?? 1);
            input.value = String(Number(value ?? field.min ?? 0));
            const out = document.createElement('output');
            out.textContent = `${input.value}${field.unit ?? ''}`;
            input.addEventListener('input', () => {
                out.textContent = `${input.value}${field.unit ?? ''}`;
                set(Number(input.value));
            });
            wrap.append(input, out);
            row.append(wrap);
            break;
        }
        case 'color': {
            const wrap = document.createElement('div');
            wrap.className = 'color-wrap';
            const swatch = document.createElement('input');
            swatch.type = 'color';
            swatch.value = toHex(value);
            const text = document.createElement('input');
            text.type = 'text';
            text.placeholder = 'inherit';
            text.value = value == null ? '' : String(value);
            const clear = document.createElement('button');
            clear.type = 'button';
            clear.className = 'color-clear';
            clear.title = 'Clear';
            clear.textContent = '×';
            swatch.addEventListener('input', () => {
                text.value = swatch.value;
                set(swatch.value);
            });
            text.addEventListener('input', () => {
                swatch.value = toHex(text.value);
                set(text.value);
            });
            clear.addEventListener('click', () => {
                text.value = '';
                set('');
            });
            wrap.append(swatch, text, clear);
            row.append(wrap);
            break;
        }
        case 'toggle': {
            row.classList.add('field-inline');
            const input = document.createElement('input');
            input.type = 'checkbox';
            input.checked = value === true || value === 'true';
            const track = document.createElement('span');
            track.className = 'switch';
            input.addEventListener('change', () => set(input.checked));
            row.append(input, track);
            break;
        }
    }
    if (field.hint) {
        const hint = document.createElement('span');
        hint.className = 'field-hint';
        hint.textContent = field.hint;
        row.append(hint);
    }
    return row;
}

function wrapSelect(select: HTMLSelectElement): HTMLElement {
    const wrap = document.createElement('div');
    wrap.className = 'select-wrap';
    wrap.append(select);
    wrap.insertAdjacentHTML('beforeend', icon('down'));
    return wrap;
}

/** Best-effort conversion of a CSS colour string to #rrggbb for the colour input. */
function toHex(value: unknown): string {
    const s = String(value ?? '').trim();
    if (/^#[0-9a-f]{6}$/i.test(s)) {
        return s;
    }
    if (/^#[0-9a-f]{3}$/i.test(s)) {
        return '#' + s.slice(1).split('').map((c) => c + c).join('');
    }
    return '#8a8f9f';
}
