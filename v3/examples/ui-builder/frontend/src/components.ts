// The component registry. Each entry describes one thing the user can drop
// into the window: its defaults, the fields shown in the inspector, and how it
// renders. The same `render` function powers the design canvas *and* the
// exported frontend, so there's a single source of truth for the output.
//
// Components are desktop-application building blocks — sidebars, toolbars,
// tables, status bars — and controls can be wired to the Go backend: buttons
// call bound service methods or emit events, and any text-like component can
// display the payload of a Go event through `bind`.

import {type UINode, newId} from './model';

export type FieldKind = 'text' | 'textarea' | 'select' | 'number' | 'color' | 'toggle' | 'range' | 'align';

export interface FieldOption {
    label: string;
    value: string;
}

export interface PropField {
    key: string;
    label: string;
    kind: FieldKind;
    group?: 'Content' | 'Style' | 'Layout' | 'Wails';
    options?: FieldOption[];
    min?: number;
    max?: number;
    step?: number;
    unit?: string;
    placeholder?: string;
    /** Extra explanation shown under the control. */
    hint?: string;
    /** Hide the field unless this predicate holds for the node's props. */
    when?: (props: Record<string, unknown>) => boolean;
    /** Re-render the inspector after a change (for fields that reveal others). */
    rerender?: boolean;
}

export type Category = 'Layout' | 'Content' | 'Controls';

export interface ComponentDef {
    type: string;
    label: string;
    category: Category;
    icon: string;
    /** Containers accept children. `direction` tells the drop logic which axis to split on. */
    container?: boolean;
    direction?: 'row' | 'column';
    defaults: Record<string, unknown>;
    fields: PropField[];
    /** Optional children created alongside the node when dropped from the palette. */
    starter?: () => UINode[];
    render: (node: UINode, children: HTMLElement[]) => HTMLElement;
}

// ---------------------------------------------------------------------------
// Helpers shared by renderers

const el = <K extends keyof HTMLElementTagNameMap>(tag: K, className?: string): HTMLElementTagNameMap[K] => {
    const node = document.createElement(tag);
    if (className) {
        node.className = className;
    }
    return node;
};

export const str = (node: UINode, key: string, fallback = ''): string => {
    const v = node.props[key];
    return v === undefined || v === null ? fallback : String(v);
};
export const num = (node: UINode, key: string, fallback = 0): number => {
    const v = Number(node.props[key]);
    return Number.isFinite(v) ? v : fallback;
};
export const bool = (node: UINode, key: string): boolean => node.props[key] === true || node.props[key] === 'true';
export const lines = (node: UINode, key: string): string[] => str(node, key).split('\n').map((s) => s.trim()).filter(Boolean);

/** Apply the shared box props: padding, gap, background, colour, alignment. */
function applyBox(target: HTMLElement, node: UINode): void {
    const pad = num(node, 'padding', NaN);
    if (!Number.isNaN(pad)) {
        target.style.padding = `${pad}px`;
    }
    const gap = num(node, 'gap', NaN);
    if (!Number.isNaN(gap)) {
        target.style.gap = `${gap}px`;
    }
    const bg = str(node, 'background');
    if (bg) {
        target.style.background = bg;
    }
    const color = str(node, 'color');
    if (color) {
        target.style.color = color;
    }
    const align = str(node, 'align');
    if (align) {
        target.style.textAlign = align;
    }
}

/** Mark an element as a live binding target: the runtime writes the payload of `channel` into it. */
function bindable(target: HTMLElement, node: UINode): void {
    const channel = str(node, 'bind').trim();
    if (channel) {
        target.dataset.bind = channel;
    }
}

const flexAlign = (v: string): string => ({start: 'flex-start', center: 'center', end: 'flex-end', stretch: 'stretch', between: 'space-between'}[v] ?? v);

let nextChildId = 0;
const child = (type: string, props: Record<string, unknown>, children?: UINode[]): UINode => {
    nextChildId += 1;
    return {id: `s${Date.now().toString(36)}${nextChildId}`, type, props, children};
};

// ----- Reusable field definitions -------------------------------------------

const alignField: PropField = {
    key: 'align', label: 'Text align', kind: 'align', group: 'Layout',
    options: [{label: 'Left', value: 'left'}, {label: 'Center', value: 'center'}, {label: 'Right', value: 'right'}],
};
const colorField: PropField = {key: 'color', label: 'Text colour', kind: 'color', group: 'Style'};
const bgField: PropField = {key: 'background', label: 'Background', kind: 'color', group: 'Style'};
const paddingField = (max = 48): PropField => ({key: 'padding', label: 'Padding', kind: 'range', group: 'Layout', min: 0, max, step: 2, unit: 'px'});
const gapField: PropField = {key: 'gap', label: 'Gap', kind: 'range', group: 'Layout', min: 0, max: 32, step: 2, unit: 'px'};
const fillField: PropField = {key: 'fill', label: 'Fill available space', kind: 'toggle', group: 'Layout', hint: 'Grow to take the remaining room in the parent.'};
const scrollField: PropField = {key: 'scroll', label: 'Scroll overflow', kind: 'toggle', group: 'Layout'};
const itemsField: PropField = {key: 'items', label: 'Align items', kind: 'select', group: 'Layout', options: [
    {label: 'Stretch', value: 'stretch'}, {label: 'Start', value: 'start'}, {label: 'Center', value: 'center'}, {label: 'End', value: 'end'},
]};
const justifyField: PropField = {key: 'justify', label: 'Justify', kind: 'select', group: 'Layout', options: [
    {label: 'Start', value: 'start'}, {label: 'Center', value: 'center'}, {label: 'End', value: 'end'}, {label: 'Space between', value: 'between'},
]};
const bindField = (what: string): PropField => ({
    key: 'bind', label: 'Bind to Go event', kind: 'text', group: 'Wails', placeholder: 'e.g. time',
    hint: `While running, ${what} shows the data of every Go event with this name (app.Event.Emit(name, data)).`,
});
const nameField: PropField = {key: 'name', label: 'Field name', kind: 'text', group: 'Wails', placeholder: 'e.g. name', hint: 'Buttons can pass this field\'s value to Go as $name.'};

const labelledField = (node: UINode, control: HTMLElement): HTMLElement => {
    const field = el('div', 'ub-field');
    field.dataset.fill = String(bool(node, 'fill'));
    field.dataset.inline = String(bool(node, 'inline'));
    const label = str(node, 'label');
    if (label) {
        const l = el('label');
        l.textContent = label;
        field.append(l);
    }
    field.append(control);
    const help = str(node, 'help');
    if (help) {
        const h = el('span', 'ub-help');
        h.textContent = help;
        field.append(h);
    }
    return field;
};

// ---------------------------------------------------------------------------
// The registry

export const components: ComponentDef[] = [
    // ----- Layout ----------------------------------------------------------
    {
        type: 'hstack', label: 'Row', category: 'Layout', icon: 'hstack',
        container: true, direction: 'row',
        defaults: {gap: 8, padding: 0, items: 'stretch', justify: 'start', fill: false, scroll: false, background: ''},
        fields: [gapField, paddingField(), itemsField, justifyField, fillField, scrollField, bgField],
        render(node, children) {
            const row = el('div', 'ub-hstack');
            applyBox(row, node);
            row.style.alignItems = flexAlign(str(node, 'items', 'stretch'));
            row.style.justifyContent = flexAlign(str(node, 'justify', 'start'));
            row.dataset.fill = String(bool(node, 'fill'));
            row.dataset.scroll = String(bool(node, 'scroll'));
            row.append(...children);
            return row;
        },
    },
    {
        type: 'vstack', label: 'Stack', category: 'Layout', icon: 'vstack',
        container: true, direction: 'column',
        defaults: {gap: 8, padding: 0, items: 'stretch', fill: false, scroll: false, background: ''},
        fields: [gapField, paddingField(), itemsField, fillField, scrollField, bgField],
        render(node, children) {
            const stack = el('div', 'ub-vstack');
            applyBox(stack, node);
            stack.style.alignItems = flexAlign(str(node, 'items', 'stretch'));
            stack.dataset.fill = String(bool(node, 'fill'));
            stack.dataset.scroll = String(bool(node, 'scroll'));
            stack.append(...children);
            return stack;
        },
    },
    {
        type: 'sidebar', label: 'Sidebar', category: 'Layout', icon: 'sidebar',
        container: true, direction: 'column',
        defaults: {width: 220, side: 'left', padding: 8, gap: 4, background: ''},
        fields: [
            {key: 'width', label: 'Width', kind: 'range', group: 'Layout', min: 140, max: 400, step: 10, unit: 'px'},
            {key: 'side', label: 'Side', kind: 'select', group: 'Layout', options: [{label: 'Left', value: 'left'}, {label: 'Right', value: 'right'}]},
            paddingField(), gapField, bgField,
        ],
        starter: () => [
            child('heading', {text: 'Workspace', level: 'label', align: 'left', color: ''}),
            child('nav', {items: 'Overview\nProjects\nActivity\n# Settings\nGeneral\nAccount', active: 0, style: 'accent'}),
        ],
        render(node, children) {
            const side = el('aside', 'ub-sidebar');
            applyBox(side, node);
            side.style.width = `${num(node, 'width', 220)}px`;
            side.dataset.side = str(node, 'side', 'left');
            side.append(...children);
            return side;
        },
    },
    {
        type: 'toolbar', label: 'Toolbar', category: 'Layout', icon: 'toolbar',
        container: true, direction: 'row',
        defaults: {height: 44, gap: 8, drag: true, background: ''},
        fields: [
            {key: 'height', label: 'Height', kind: 'range', group: 'Layout', min: 32, max: 72, step: 2, unit: 'px'},
            gapField,
            {key: 'drag', label: 'Window drag region', kind: 'toggle', group: 'Wails', hint: 'Sets --wails-draggable: drag so the window can be moved by this bar (frameless / hidden title bar windows).'},
            bgField,
        ],
        starter: () => [
            child('heading', {text: 'My App', level: 'window', align: 'left', color: ''}),
            child('spacer', {size: 0, fill: true}),
            child('input', {label: '', name: 'query', placeholder: 'Search', type: 'search', help: '', fill: false, inline: false}),
            child('button', {label: 'New', variant: 'primary', size: 'sm', full: false, action: 'none'}),
        ],
        render(node, children) {
            const bar = el('header', 'ub-toolbar');
            applyBox(bar, node);
            bar.style.height = `${num(node, 'height', 44)}px`;
            bar.dataset.drag = String(bool(node, 'drag'));
            bar.append(...children);
            return bar;
        },
    },
    {
        type: 'content', label: 'Content', category: 'Layout', icon: 'content',
        container: true, direction: 'column',
        defaults: {padding: 16, gap: 12, scroll: true, background: ''},
        fields: [paddingField(64), gapField, scrollField, bgField],
        render(node, children) {
            const main = el('main', 'ub-content');
            applyBox(main, node);
            main.style.overflow = bool(node, 'scroll') ? 'auto' : 'hidden';
            main.append(...children);
            return main;
        },
    },
    {
        type: 'group', label: 'Group', category: 'Layout', icon: 'group',
        container: true, direction: 'column',
        defaults: {title: 'General', gap: 10, padding: 14, fill: false, background: ''},
        fields: [
            {key: 'title', label: 'Title', kind: 'text', group: 'Content', placeholder: 'Leave empty for no title'},
            gapField, paddingField(), fillField, bgField,
        ],
        render(node, children) {
            const group = el('section', 'ub-group');
            applyBox(group, node);
            group.dataset.fill = String(bool(node, 'fill'));
            const title = str(node, 'title');
            if (title) {
                const h = el('h3', 'ub-group-title');
                h.textContent = title;
                group.append(h);
            }
            group.append(...children);
            return group;
        },
    },
    {
        type: 'statusbar', label: 'Status bar', category: 'Layout', icon: 'statusbar',
        container: true, direction: 'row',
        defaults: {gap: 14, background: ''},
        fields: [gapField, bgField],
        starter: () => [
            child('text', {text: 'Ready', size: 'sm', tone: 'muted', mono: false, align: 'left', color: '', bind: 'status'}),
            child('spacer', {size: 0, fill: true}),
            child('text', {text: 'Waiting for time event…', size: 'sm', tone: 'muted', mono: false, align: 'left', color: '', bind: 'time'}),
        ],
        render(node, children) {
            const bar = el('footer', 'ub-statusbar');
            applyBox(bar, node);
            bar.append(...children);
            return bar;
        },
    },
    {
        type: 'tabs', label: 'Tabs', category: 'Layout', icon: 'tabs',
        defaults: {items: 'General\nAppearance\nAdvanced', active: 0, style: 'underline'},
        fields: [
            {key: 'items', label: 'Tabs (one per line)', kind: 'textarea', group: 'Content'},
            {key: 'active', label: 'Active tab', kind: 'number', group: 'Content', min: 0, step: 1},
            {key: 'style', label: 'Style', kind: 'select', group: 'Style', options: [{label: 'Underline', value: 'underline'}, {label: 'Segmented', value: 'segmented'}]},
        ],
        render(node) {
            const tabs = el('div', 'ub-tabs');
            tabs.setAttribute('role', 'tablist');
            tabs.dataset.style = str(node, 'style', 'underline');
            const active = num(node, 'active', 0);
            for (const [i, label] of lines(node, 'items').entries()) {
                const b = el('button');
                b.type = 'button';
                b.setAttribute('role', 'tab');
                b.setAttribute('aria-selected', String(i === active));
                b.textContent = label;
                tabs.append(b);
            }
            return tabs;
        },
    },
    {
        type: 'nav', label: 'Nav list', category: 'Layout', icon: 'nav',
        defaults: {items: 'Overview\nProjects\nActivity\n# Settings\nGeneral\nAccount', active: 0, style: 'accent'},
        fields: [
            {key: 'items', label: 'Items (one per line)', kind: 'textarea', group: 'Content', hint: 'Start a line with # to make a section header.'},
            {key: 'active', label: 'Active item', kind: 'number', group: 'Content', min: 0, step: 1},
            {key: 'style', label: 'Selection style', kind: 'select', group: 'Style', options: [{label: 'Accent', value: 'accent'}, {label: 'Plain', value: 'plain'}]},
        ],
        render(node) {
            const list = el('ul', 'ub-nav');
            list.dataset.style = str(node, 'style', 'accent');
            const active = num(node, 'active', 0);
            let index = 0;
            for (const item of lines(node, 'items')) {
                const li = el('li');
                if (item.startsWith('#')) {
                    li.className = 'ub-nav-header';
                    li.textContent = item.replace(/^#\s*/, '');
                } else {
                    li.textContent = item;
                    if (index === active) {
                        li.setAttribute('aria-current', 'true');
                    }
                    index += 1;
                }
                list.append(li);
            }
            return list;
        },
    },
    {
        type: 'divider', label: 'Divider', category: 'Layout', icon: 'divider',
        defaults: {},
        fields: [],
        render() {
            return el('hr', 'ub-divider');
        },
    },
    {
        type: 'spacer', label: 'Spacer', category: 'Layout', icon: 'spacer',
        defaults: {size: 16, fill: false},
        fields: [
            {key: 'size', label: 'Size', kind: 'range', group: 'Layout', min: 0, max: 120, step: 4, unit: 'px'},
            {key: 'fill', label: 'Push apart', kind: 'toggle', group: 'Layout', hint: 'Expands to fill the row or stack, pushing siblings to the edges.'},
        ],
        render(node) {
            const s = el('div', 'ub-spacer');
            const size = num(node, 'size', 16);
            s.style.width = `${size}px`;
            s.style.height = `${size}px`;
            s.dataset.fill = String(bool(node, 'fill'));
            return s;
        },
    },

    // ----- Content ---------------------------------------------------------
    {
        type: 'heading', label: 'Heading', category: 'Content', icon: 'heading',
        defaults: {text: 'Overview', level: 'title', align: 'left', color: ''},
        fields: [
            {key: 'text', label: 'Text', kind: 'text', group: 'Content'},
            {key: 'level', label: 'Level', kind: 'select', group: 'Style', options: [
                {label: 'Window title', value: 'window'}, {label: 'Page title', value: 'title'}, {label: 'Section', value: 'section'}, {label: 'Label', value: 'label'},
            ]},
            alignField, colorField,
        ],
        render(node) {
            const level = str(node, 'level', 'title');
            const h = el(level === 'title' ? 'h1' : level === 'section' ? 'h2' : 'h3', 'ub-heading');
            h.dataset.level = level;
            h.textContent = str(node, 'text');
            applyBox(h, node);
            return h;
        },
    },
    {
        type: 'text', label: 'Text', category: 'Content', icon: 'text',
        defaults: {text: 'Select an item on the left to see its details.', size: 'md', tone: 'default', mono: false, align: 'left', color: '', bind: ''},
        fields: [
            {key: 'text', label: 'Text', kind: 'textarea', group: 'Content'},
            {key: 'size', label: 'Size', kind: 'select', group: 'Style', options: [{label: 'Small', value: 'sm'}, {label: 'Normal', value: 'md'}, {label: 'Large', value: 'lg'}]},
            {key: 'tone', label: 'Tone', kind: 'select', group: 'Style', options: [
                {label: 'Default', value: 'default'}, {label: 'Muted', value: 'muted'}, {label: 'Success', value: 'ok'}, {label: 'Warning', value: 'warn'}, {label: 'Danger', value: 'danger'},
            ]},
            {key: 'mono', label: 'Monospace', kind: 'toggle', group: 'Style'},
            alignField, colorField, bindField('this text'),
        ],
        render(node) {
            const p = el('p', 'ub-text');
            p.textContent = str(node, 'text');
            p.dataset.size = str(node, 'size', 'md');
            p.dataset.tone = str(node, 'tone', 'default');
            p.dataset.mono = String(bool(node, 'mono'));
            applyBox(p, node);
            bindable(p, node);
            return p;
        },
    },
    {
        type: 'badge', label: 'Badge', category: 'Content', icon: 'badge',
        defaults: {text: 'Connected', tone: 'ok', bind: ''},
        fields: [
            {key: 'text', label: 'Text', kind: 'text', group: 'Content'},
            {key: 'tone', label: 'Tone', kind: 'select', group: 'Style', options: [
                {label: 'Accent', value: 'accent'}, {label: 'Neutral', value: 'neutral'}, {label: 'Success', value: 'ok'}, {label: 'Warning', value: 'warn'}, {label: 'Danger', value: 'danger'},
            ]},
            bindField('the badge'),
        ],
        render(node) {
            const badge = el('span', 'ub-badge');
            badge.textContent = str(node, 'text');
            badge.dataset.tone = str(node, 'tone', 'accent');
            bindable(badge, node);
            return badge;
        },
    },
    {
        type: 'metric', label: 'Metric', category: 'Content', icon: 'metric',
        defaults: {label: 'Open issues', value: '128', note: '12 assigned to you', fill: false, bind: ''},
        fields: [
            {key: 'label', label: 'Label', kind: 'text', group: 'Content'},
            {key: 'value', label: 'Value', kind: 'text', group: 'Content'},
            {key: 'note', label: 'Note', kind: 'text', group: 'Content'},
            fillField, bindField('the value'),
        ],
        render(node) {
            const metric = el('div', 'ub-metric');
            metric.dataset.fill = String(bool(node, 'fill'));
            const label = el('span', 'ub-metric-label');
            label.textContent = str(node, 'label');
            const value = el('span', 'ub-metric-value');
            value.textContent = str(node, 'value');
            bindable(value, node);
            metric.append(label, value);
            const note = str(node, 'note');
            if (note) {
                const n = el('span', 'ub-metric-note');
                n.textContent = note;
                metric.append(n);
            }
            return metric;
        },
    },
    {
        type: 'table', label: 'Table', category: 'Content', icon: 'table',
        defaults: {
            columns: 'Name | Status | Updated',
            rows: 'wails-app | Building | 2 min ago\nrelease-notes | Ready | 1 h ago\nwebsite | Ready | yesterday',
            striped: true, selected: -1, fill: false,
        },
        fields: [
            {key: 'columns', label: 'Columns', kind: 'text', group: 'Content', hint: 'Separate columns with |'},
            {key: 'rows', label: 'Rows (one per line)', kind: 'textarea', group: 'Content', hint: 'Cells separated with |'},
            {key: 'striped', label: 'Striped rows', kind: 'toggle', group: 'Style'},
            {key: 'selected', label: 'Selected row', kind: 'number', group: 'Content', min: -1, step: 1, hint: '-1 for none'},
            fillField,
        ],
        render(node) {
            const wrap = el('div', 'ub-table-wrap');
            wrap.dataset.fill = String(bool(node, 'fill'));
            const table = el('table', 'ub-table');
            table.dataset.striped = String(bool(node, 'striped'));
            const split = (s: string): string[] => s.split('|').map((c) => c.trim());
            const thead = el('thead');
            const hr = el('tr');
            for (const c of split(str(node, 'columns'))) {
                const th = el('th');
                th.textContent = c;
                hr.append(th);
            }
            thead.append(hr);
            const tbody = el('tbody');
            const selected = num(node, 'selected', -1);
            for (const [i, line] of lines(node, 'rows').entries()) {
                const tr = el('tr');
                if (i === selected) {
                    tr.setAttribute('aria-selected', 'true');
                }
                for (const c of split(line)) {
                    const td = el('td');
                    td.textContent = c;
                    tr.append(td);
                }
                tbody.append(tr);
            }
            table.append(thead, tbody);
            wrap.append(table);
            return wrap;
        },
    },
    {
        type: 'list', label: 'List', category: 'Content', icon: 'list',
        defaults: {items: 'Build started\nFrontend compiled\nBinary signed', style: 'rows'},
        fields: [
            {key: 'items', label: 'Items (one per line)', kind: 'textarea', group: 'Content'},
            {key: 'style', label: 'Style', kind: 'select', group: 'Style', options: [{label: 'Rows', value: 'rows'}, {label: 'Bullets', value: 'bullets'}]},
        ],
        render(node) {
            const list = el('ul', 'ub-list');
            list.dataset.style = str(node, 'style', 'rows');
            for (const item of lines(node, 'items')) {
                const li = el('li');
                li.textContent = item;
                list.append(li);
            }
            return list;
        },
    },
    {
        type: 'keyvalue', label: 'Key / value', category: 'Content', icon: 'keyvalue',
        defaults: {pairs: 'Version: 1.4.2\nPlatform: darwin/arm64\nGo: 1.25\nWails: v3'},
        fields: [{key: 'pairs', label: 'Pairs (key: value per line)', kind: 'textarea', group: 'Content'}],
        render(node) {
            const dl = el('dl', 'ub-keyvalue');
            for (const line of lines(node, 'pairs')) {
                const [k, ...rest] = line.split(':');
                const dt = el('dt');
                dt.textContent = k.trim();
                const dd = el('dd');
                dd.textContent = rest.join(':').trim();
                dl.append(dt, dd);
            }
            return dl;
        },
    },
    {
        type: 'progress', label: 'Progress', category: 'Content', icon: 'progress',
        defaults: {label: 'Indexing files', value: 62, bind: ''},
        fields: [
            {key: 'label', label: 'Label', kind: 'text', group: 'Content'},
            {key: 'value', label: 'Value', kind: 'range', group: 'Content', min: 0, max: 100, step: 1, unit: '%'},
            {...bindField('the bar'), hint: 'While running, a Go event with this name and a number 0–100 as data moves the bar.'},
        ],
        render(node) {
            const wrap = el('div', 'ub-progress');
            const head = el('div', 'ub-progress-head');
            const label = el('span');
            label.textContent = str(node, 'label');
            const value = el('output');
            value.textContent = `${num(node, 'value', 0)}%`;
            head.append(label, value);
            const track = el('div', 'ub-progress-track');
            const fill = el('i');
            fill.style.width = `${Math.max(0, Math.min(100, num(node, 'value', 0)))}%`;
            track.append(fill);
            wrap.append(head, track);
            bindable(wrap, node);
            return wrap;
        },
    },
    {
        type: 'log', label: 'Console', category: 'Content', icon: 'log',
        defaults: {text: '$ wails3 dev\n  ✓ frontend built in 312ms\n  ✓ bindings generated\n  → window opened', rows: 6, fill: false, bind: ''},
        fields: [
            {key: 'text', label: 'Text', kind: 'textarea', group: 'Content'},
            {key: 'rows', label: 'Rows', kind: 'range', group: 'Layout', min: 3, max: 30, step: 1},
            fillField,
            {...bindField('the console'), hint: 'While running, Go events with this name are appended as new lines.'},
        ],
        render(node) {
            const pre = el('pre', 'ub-log');
            pre.textContent = str(node, 'text');
            pre.style.height = bool(node, 'fill') ? '' : `${num(node, 'rows', 6) * 18 + 20}px`;
            pre.dataset.fill = String(bool(node, 'fill'));
            pre.dataset.bindMode = 'append';
            bindable(pre, node);
            return pre;
        },
    },
    {
        type: 'image', label: 'Image', category: 'Content', icon: 'image',
        defaults: {src: '', alt: 'Preview', height: 160, fit: 'cover'},
        fields: [
            {key: 'src', label: 'Image URL', kind: 'text', group: 'Content', placeholder: 'https://… or /asset.png (empty shows a placeholder)'},
            {key: 'alt', label: 'Alt text', kind: 'text', group: 'Content'},
            {key: 'height', label: 'Height', kind: 'range', group: 'Layout', min: 40, max: 480, step: 8, unit: 'px'},
            {key: 'fit', label: 'Fit', kind: 'select', group: 'Style', options: [{label: 'Cover', value: 'cover'}, {label: 'Contain', value: 'contain'}]},
        ],
        render(node) {
            const src = str(node, 'src').trim();
            const height = num(node, 'height', 160);
            if (!src) {
                const ph = el('div', 'ub-image-placeholder');
                ph.style.height = `${height}px`;
                ph.textContent = str(node, 'alt', 'Image');
                return ph;
            }
            const img = el('img', 'ub-image');
            img.src = src;
            img.alt = str(node, 'alt');
            img.style.width = '100%';
            img.style.height = `${height}px`;
            img.style.objectFit = str(node, 'fit', 'cover');
            return img;
        },
    },

    // ----- Controls --------------------------------------------------------
    {
        type: 'button', label: 'Button', category: 'Controls', icon: 'button',
        defaults: {
            label: 'Greet', variant: 'primary', size: 'md', full: false,
            action: 'call', service: 'GreetService', method: 'Greet', args: '$name', resultTo: 'greeting',
            event: 'button:clicked', url: 'https://wails.io', windowAction: 'Minimise',
        },
        fields: [
            {key: 'label', label: 'Label', kind: 'text', group: 'Content'},
            {key: 'variant', label: 'Variant', kind: 'select', group: 'Style', options: [
                {label: 'Primary', value: 'primary'}, {label: 'Default', value: 'default'}, {label: 'Ghost', value: 'ghost'}, {label: 'Danger', value: 'danger'},
            ]},
            {key: 'size', label: 'Size', kind: 'select', group: 'Style', options: [{label: 'Small', value: 'sm'}, {label: 'Medium', value: 'md'}, {label: 'Large', value: 'lg'}]},
            {key: 'full', label: 'Full width', kind: 'toggle', group: 'Layout'},
            {key: 'action', label: 'On click', kind: 'select', group: 'Wails', rerender: true, options: [
                {label: 'Nothing', value: 'none'},
                {label: 'Call a Go method', value: 'call'},
                {label: 'Emit an event to Go', value: 'event'},
                {label: 'Open a URL in the browser', value: 'url'},
                {label: 'Window action', value: 'window'},
            ]},
            {key: 'service', label: 'Service', kind: 'text', group: 'Wails', placeholder: 'GreetService', when: (p) => p.action === 'call',
                hint: 'A struct registered with application.NewService(&GreetService{}).'},
            {key: 'method', label: 'Method', kind: 'text', group: 'Wails', placeholder: 'Greet', when: (p) => p.action === 'call'},
            {key: 'args', label: 'Arguments (one per line)', kind: 'textarea', group: 'Wails', when: (p) => p.action === 'call',
                hint: '$name reads the input with that field name. Numbers and true/false are sent as such; anything else as a string.'},
            {key: 'resultTo', label: 'Show result in', kind: 'text', group: 'Wails', placeholder: 'greeting', when: (p) => p.action === 'call',
                hint: 'Components bound to this name display the value the Go method returns.'},
            {key: 'event', label: 'Event name', kind: 'text', group: 'Wails', when: (p) => p.action === 'event',
                hint: 'Go receives it with app.Event.On(name, …).'},
            {key: 'url', label: 'URL', kind: 'text', group: 'Wails', when: (p) => p.action === 'url'},
            {key: 'windowAction', label: 'Action', kind: 'select', group: 'Wails', when: (p) => p.action === 'window', options: [
                {label: 'Minimise', value: 'Minimise'}, {label: 'Maximise', value: 'Maximise'}, {label: 'Restore', value: 'Restore'},
                {label: 'Toggle fullscreen', value: 'ToggleFullscreen'}, {label: 'Center', value: 'Center'}, {label: 'Close', value: 'Close'},
            ]},
        ],
        render(node) {
            const b = el('button', 'ub-btn');
            b.type = 'button';
            b.textContent = str(node, 'label');
            b.dataset.variant = str(node, 'variant', 'primary');
            b.dataset.size = str(node, 'size', 'md');
            b.dataset.full = String(bool(node, 'full'));
            switch (str(node, 'action', 'none')) {
                case 'call':
                    b.dataset.call = `${str(node, 'service')}.${str(node, 'method')}`;
                    b.dataset.args = str(node, 'args');
                    if (str(node, 'resultTo').trim()) {
                        b.dataset.resultTo = str(node, 'resultTo').trim();
                    }
                    break;
                case 'event':
                    b.dataset.wmlEvent = str(node, 'event');
                    break;
                case 'url':
                    b.dataset.wmlOpenurl = str(node, 'url');
                    break;
                case 'window':
                    b.dataset.wmlWindow = str(node, 'windowAction', 'Minimise');
                    break;
            }
            return b;
        },
    },
    {
        type: 'input', label: 'Input', category: 'Controls', icon: 'input',
        defaults: {label: 'Name', name: 'name', placeholder: 'Your name', type: 'text', help: '', fill: false, inline: false},
        fields: [
            {key: 'label', label: 'Label', kind: 'text', group: 'Content'},
            {key: 'placeholder', label: 'Placeholder', kind: 'text', group: 'Content'},
            {key: 'help', label: 'Help text', kind: 'text', group: 'Content'},
            {key: 'type', label: 'Type', kind: 'select', group: 'Style', options: [
                {label: 'Text', value: 'text'}, {label: 'Search', value: 'search'}, {label: 'Email', value: 'email'}, {label: 'Password', value: 'password'}, {label: 'Number', value: 'number'},
            ]},
            {key: 'inline', label: 'Label beside field', kind: 'toggle', group: 'Layout'},
            fillField, nameField,
        ],
        render(node) {
            const input = el('input');
            input.type = str(node, 'type', 'text');
            input.placeholder = str(node, 'placeholder');
            if (str(node, 'name').trim()) {
                input.name = str(node, 'name').trim();
            }
            return labelledField(node, input);
        },
    },
    {
        type: 'textarea', label: 'Textarea', category: 'Controls', icon: 'textarea',
        defaults: {label: 'Notes', name: 'notes', placeholder: '', rows: 4, help: '', fill: false, inline: false},
        fields: [
            {key: 'label', label: 'Label', kind: 'text', group: 'Content'},
            {key: 'placeholder', label: 'Placeholder', kind: 'text', group: 'Content'},
            {key: 'rows', label: 'Rows', kind: 'range', group: 'Layout', min: 2, max: 12, step: 1},
            fillField, nameField,
        ],
        render(node) {
            const ta = el('textarea');
            ta.placeholder = str(node, 'placeholder');
            ta.rows = num(node, 'rows', 4);
            if (str(node, 'name').trim()) {
                ta.name = str(node, 'name').trim();
            }
            return labelledField(node, ta);
        },
    },
    {
        type: 'select', label: 'Select', category: 'Controls', icon: 'select',
        defaults: {label: 'Theme', name: 'theme', options: 'System\nLight\nDark', help: '', fill: false, inline: true},
        fields: [
            {key: 'label', label: 'Label', kind: 'text', group: 'Content'},
            {key: 'options', label: 'Options (one per line)', kind: 'textarea', group: 'Content'},
            {key: 'inline', label: 'Label beside field', kind: 'toggle', group: 'Layout'},
            fillField, nameField,
        ],
        render(node) {
            const select = el('select');
            if (str(node, 'name').trim()) {
                select.name = str(node, 'name').trim();
            }
            for (const opt of lines(node, 'options')) {
                const o = el('option');
                o.textContent = opt;
                select.append(o);
            }
            return labelledField(node, select);
        },
    },
    {
        type: 'toggle', label: 'Toggle', category: 'Controls', icon: 'toggle',
        defaults: {label: 'Launch at login', name: 'autostart', on: true},
        fields: [
            {key: 'label', label: 'Label', kind: 'text', group: 'Content'},
            {key: 'on', label: 'Switched on', kind: 'toggle', group: 'Content'},
            nameField,
        ],
        render(node) {
            const t = el('label', 'ub-toggle');
            t.dataset.on = String(bool(node, 'on'));
            const text = el('span');
            text.textContent = str(node, 'label');
            const input = el('input');
            input.type = 'checkbox';
            input.checked = bool(node, 'on');
            input.hidden = true;
            if (str(node, 'name').trim()) {
                input.name = str(node, 'name').trim();
            }
            t.append(text, input, el('span', 'ub-switch'));
            return t;
        },
    },
    {
        type: 'checkbox', label: 'Checkbox', category: 'Controls', icon: 'checkbox',
        defaults: {label: 'Show hidden files', name: 'hidden', on: false},
        fields: [
            {key: 'label', label: 'Label', kind: 'text', group: 'Content'},
            {key: 'on', label: 'Checked', kind: 'toggle', group: 'Content'},
            nameField,
        ],
        render(node) {
            const c = el('label', 'ub-check');
            c.dataset.on = String(bool(node, 'on'));
            const input = el('input');
            input.type = 'checkbox';
            input.checked = bool(node, 'on');
            input.hidden = true;
            if (str(node, 'name').trim()) {
                input.name = str(node, 'name').trim();
            }
            const text = el('span');
            text.textContent = str(node, 'label');
            c.append(input, el('i'), text);
            return c;
        },
    },
    {
        type: 'slider', label: 'Slider', category: 'Controls', icon: 'slider',
        defaults: {label: 'Volume', name: 'volume', value: 60, min: 0, max: 100},
        fields: [
            {key: 'label', label: 'Label', kind: 'text', group: 'Content'},
            {key: 'value', label: 'Value', kind: 'number', group: 'Content'},
            {key: 'min', label: 'Min', kind: 'number', group: 'Content'},
            {key: 'max', label: 'Max', kind: 'number', group: 'Content'},
            nameField,
        ],
        render(node) {
            const wrap = el('div', 'ub-slider');
            const head = el('div', 'ub-slider-head');
            const label = el('span');
            label.textContent = str(node, 'label');
            const out = el('output');
            out.textContent = str(node, 'value', '0');
            head.append(label, out);
            const input = el('input');
            input.type = 'range';
            input.min = str(node, 'min', '0');
            input.max = str(node, 'max', '100');
            input.value = str(node, 'value', '0');
            if (str(node, 'name').trim()) {
                input.name = str(node, 'name').trim();
            }
            wrap.append(head, input);
            return wrap;
        },
    },
];

// The root is a pseudo-component: it's never in the palette, but it renders
// and accepts children like any other container. It is the window's content.
export const rootDef: ComponentDef = {
    type: 'root', label: 'Window', category: 'Layout', icon: 'root',
    container: true, direction: 'column',
    defaults: {},
    fields: [],
    render(_node, children) {
        const root = el('div', 'ub-root');
        root.append(...children);
        return root;
    },
};

const byType = new Map<string, ComponentDef>(components.map((c) => [c.type, c]));
byType.set(rootDef.type, rootDef);

export function getDef(type: string): ComponentDef {
    return byType.get(type) ?? unknownDef(type);
}

export const categories: Category[] = ['Layout', 'Content', 'Controls'];

/** Build a fresh node (with defaults and starter children) for a palette drop. */
export function createNode(type: string): UINode {
    const def = getDef(type);
    return {
        id: newId(),
        type,
        props: {...def.defaults},
        children: def.container ? (def.starter?.() ?? []) : undefined,
    };
}

/** Nodes of a type this build doesn't know (e.g. a layout from a newer version). */
function unknownDef(type: string): ComponentDef {
    return {
        type, label: type, category: 'Content', icon: 'cursor', defaults: {}, fields: [],
        render() {
            const d = el('div');
            d.style.cssText = 'padding:12px;border:1px dashed currentColor;border-radius:8px;opacity:.6;font-size:12px';
            d.textContent = `Unknown component “${type}”`;
            return d;
        },
    };
}
