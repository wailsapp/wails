// The component registry. Each entry describes one thing the user can drop
// onto the canvas: its defaults, the fields shown in the inspector, and how it
// renders. The same `render` function powers the design canvas *and* the
// exported HTML, so there's a single source of truth for the output.

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
    group?: 'Content' | 'Style' | 'Layout';
    options?: FieldOption[];
    min?: number;
    max?: number;
    step?: number;
    unit?: string;
    placeholder?: string;
}

export type Category = 'Layout' | 'Content' | 'Form' | 'Media';

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

const str = (node: UINode, key: string, fallback = ''): string => {
    const v = node.props[key];
    return v === undefined || v === null ? fallback : String(v);
};
const num = (node: UINode, key: string, fallback = 0): number => {
    const v = Number(node.props[key]);
    return Number.isFinite(v) ? v : fallback;
};
const bool = (node: UINode, key: string): boolean => node.props[key] === true || node.props[key] === 'true';

/** Apply the "box" style props (padding, background, radius, max width) that most components share. */
function applyBox(target: HTMLElement, node: UINode): void {
    const pad = num(node, 'padding', NaN);
    if (!Number.isNaN(pad)) {
        target.style.padding = `${pad}px`;
    }
    const bg = str(node, 'background');
    if (bg && bg !== 'transparent') {
        target.style.background = bg;
    }
    const radius = num(node, 'radius', NaN);
    if (!Number.isNaN(radius)) {
        target.style.borderRadius = `${radius}px`;
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

const alignField: PropField = {
    key: 'align', label: 'Align', kind: 'align', group: 'Layout',
    options: [{label: 'Left', value: 'left'}, {label: 'Center', value: 'center'}, {label: 'Right', value: 'right'}],
};
const colorField: PropField = {key: 'color', label: 'Text colour', kind: 'color', group: 'Style'};
const paddingField = (max = 96): PropField => ({key: 'padding', label: 'Padding', kind: 'range', group: 'Layout', min: 0, max, step: 4, unit: 'px'});
const gapField: PropField = {key: 'gap', label: 'Gap', kind: 'range', group: 'Layout', min: 0, max: 64, step: 4, unit: 'px'};
const radiusField: PropField = {key: 'radius', label: 'Corner radius', kind: 'range', group: 'Style', min: 0, max: 32, step: 2, unit: 'px'};
const bgField: PropField = {key: 'background', label: 'Background', kind: 'color', group: 'Style'};

const flexAlign = (v: string): string => ({start: 'flex-start', center: 'center', end: 'flex-end', stretch: 'stretch', between: 'space-between'}[v] ?? v);

let nextChildId = 0;
const child = (type: string, props: Record<string, unknown>, children?: UINode[]): UINode => {
    nextChildId += 1;
    return {id: `s${Date.now().toString(36)}${nextChildId}`, type, props, children};
};

// ---------------------------------------------------------------------------
// The registry

export const components: ComponentDef[] = [
    // ----- Layout ----------------------------------------------------------
    {
        type: 'section', label: 'Section', category: 'Layout', icon: 'section',
        container: true, direction: 'column',
        defaults: {padding: 48, gap: 16, maxWidth: 960, background: 'transparent', align: 'left', items: 'stretch'},
        fields: [
            paddingField(160), gapField,
            {key: 'maxWidth', label: 'Content width', kind: 'range', group: 'Layout', min: 320, max: 1400, step: 20, unit: 'px'},
            {key: 'items', label: 'Items', kind: 'select', group: 'Layout', options: [
                {label: 'Stretch', value: 'stretch'}, {label: 'Start', value: 'start'}, {label: 'Center', value: 'center'}, {label: 'End', value: 'end'},
            ]},
            alignField, bgField,
        ],
        render(node, children) {
            const section = el('section', 'ub-section');
            const inner = el('div', 'ub-inner');
            applyBox(section, node);
            inner.style.gap = `${num(node, 'gap', 16)}px`;
            inner.style.maxWidth = `${num(node, 'maxWidth', 960)}px`;
            inner.style.alignItems = flexAlign(str(node, 'items', 'stretch'));
            inner.append(...children);
            section.append(inner);
            return section;
        },
    },
    {
        type: 'row', label: 'Row', category: 'Layout', icon: 'row',
        container: true, direction: 'row',
        defaults: {gap: 16, justify: 'start', items: 'stretch', equal: true, stack: true, padding: 0},
        fields: [
            gapField,
            {key: 'justify', label: 'Justify', kind: 'select', group: 'Layout', options: [
                {label: 'Start', value: 'start'}, {label: 'Center', value: 'center'}, {label: 'End', value: 'end'}, {label: 'Space between', value: 'between'},
            ]},
            {key: 'items', label: 'Align items', kind: 'select', group: 'Layout', options: [
                {label: 'Stretch', value: 'stretch'}, {label: 'Start', value: 'start'}, {label: 'Center', value: 'center'}, {label: 'End', value: 'end'},
            ]},
            {key: 'equal', label: 'Equal widths', kind: 'toggle', group: 'Layout'},
            {key: 'stack', label: 'Stack on phones', kind: 'toggle', group: 'Layout'},
            paddingField(), bgField,
        ],
        starter: () => [child('card', {padding: 24, radius: 12, elevated: true, gap: 8}, [
            child('heading', {text: 'Column one', level: 'h4'}),
            child('text', {text: 'Rows split into equal columns. Drop anything inside.', size: 'sm', muted: true}),
        ]), child('card', {padding: 24, radius: 12, elevated: true, gap: 8}, [
            child('heading', {text: 'Column two', level: 'h4'}),
            child('text', {text: 'Turn off “Equal widths” to let content size the columns.', size: 'sm', muted: true}),
        ])],
        render(node, children) {
            const row = el('div', 'ub-row');
            applyBox(row, node);
            row.style.gap = `${num(node, 'gap', 16)}px`;
            row.style.justifyContent = flexAlign(str(node, 'justify', 'start'));
            row.style.alignItems = flexAlign(str(node, 'items', 'stretch'));
            row.dataset.equal = String(bool(node, 'equal'));
            row.dataset.stack = String(bool(node, 'stack'));
            row.append(...children);
            return row;
        },
    },
    {
        type: 'card', label: 'Card', category: 'Layout', icon: 'card',
        container: true, direction: 'column',
        defaults: {padding: 24, gap: 12, radius: 14, elevated: true, background: '', align: 'left'},
        fields: [
            paddingField(), gapField, radiusField,
            {key: 'elevated', label: 'Shadow', kind: 'toggle', group: 'Style'},
            alignField, bgField,
        ],
        render(node, children) {
            const card = el('div', 'ub-card');
            applyBox(card, node);
            card.style.gap = `${num(node, 'gap', 12)}px`;
            card.dataset.elevated = String(bool(node, 'elevated'));
            card.append(...children);
            return card;
        },
    },
    {
        type: 'navbar', label: 'Navbar', category: 'Layout', icon: 'navbar',
        defaults: {brand: 'Acme', links: 'Product\nPricing\nDocs\nBlog', cta: 'Get started', background: ''},
        fields: [
            {key: 'brand', label: 'Brand', kind: 'text', group: 'Content'},
            {key: 'links', label: 'Links (one per line)', kind: 'textarea', group: 'Content'},
            {key: 'cta', label: 'Button label', kind: 'text', group: 'Content', placeholder: 'Leave empty to hide'},
            bgField,
        ],
        render(node) {
            const bar = el('header', 'ub-navbar');
            applyBox(bar, node);
            const brand = el('span', 'ub-navbar-brand');
            brand.append(el('i'), document.createTextNode(str(node, 'brand', 'Brand')));
            const nav = el('nav');
            for (const link of str(node, 'links').split('\n').map((s) => s.trim()).filter(Boolean)) {
                const a = el('a');
                a.href = '#';
                a.textContent = link;
                nav.append(a);
            }
            bar.append(brand, nav);
            const cta = str(node, 'cta');
            if (cta) {
                const btn = el('a', 'ub-btn');
                btn.href = '#';
                btn.dataset.variant = 'primary';
                btn.dataset.size = 'sm';
                btn.textContent = cta;
                bar.append(btn);
            }
            return bar;
        },
    },

    // ----- Content ---------------------------------------------------------
    {
        type: 'heading', label: 'Heading', category: 'Content', icon: 'heading',
        defaults: {text: 'Build interfaces by dragging', level: 'h1', align: 'left', gradient: false, color: ''},
        fields: [
            {key: 'text', label: 'Text', kind: 'textarea', group: 'Content'},
            {key: 'level', label: 'Level', kind: 'select', group: 'Style', options: [
                {label: 'H1 · Display', value: 'h1'}, {label: 'H2 · Title', value: 'h2'}, {label: 'H3 · Section', value: 'h3'}, {label: 'H4 · Label', value: 'h4'},
            ]},
            {key: 'gradient', label: 'Gradient text', kind: 'toggle', group: 'Style'},
            alignField, colorField,
        ],
        render(node) {
            const level = str(node, 'level', 'h1') as 'h1' | 'h2' | 'h3' | 'h4';
            const h = el(level, `ub-heading ub-${level}`);
            h.textContent = str(node, 'text');
            h.dataset.gradient = String(bool(node, 'gradient'));
            applyBox(h, node);
            if (bool(node, 'gradient') && str(node, 'align') === 'center') {
                h.style.marginInline = 'auto';
            }
            return h;
        },
    },
    {
        type: 'text', label: 'Text', category: 'Content', icon: 'text',
        defaults: {text: 'Compose pages from ready-made blocks, tune every property in the inspector, then export clean HTML.', size: 'md', muted: false, align: 'left', color: ''},
        fields: [
            {key: 'text', label: 'Text', kind: 'textarea', group: 'Content'},
            {key: 'size', label: 'Size', kind: 'select', group: 'Style', options: [
                {label: 'Small', value: 'sm'}, {label: 'Medium', value: 'md'}, {label: 'Large', value: 'lg'},
            ]},
            {key: 'muted', label: 'Muted', kind: 'toggle', group: 'Style'},
            alignField, colorField,
        ],
        render(node) {
            const p = el('p', 'ub-text');
            p.textContent = str(node, 'text');
            p.dataset.size = str(node, 'size', 'md');
            p.dataset.muted = String(bool(node, 'muted'));
            applyBox(p, node);
            return p;
        },
    },
    {
        type: 'list', label: 'List', category: 'Content', icon: 'list',
        defaults: {items: 'Native desktop window\nType-safe Go bindings\nOne binary, every platform', style: 'check'},
        fields: [
            {key: 'items', label: 'Items (one per line)', kind: 'textarea', group: 'Content'},
            {key: 'style', label: 'Style', kind: 'select', group: 'Style', options: [
                {label: 'Checkmarks', value: 'check'}, {label: 'Bullets', value: 'disc'}, {label: 'Numbered', value: 'decimal'},
            ]},
            colorField,
        ],
        render(node) {
            const style = str(node, 'style', 'check');
            const list = el(style === 'decimal' ? 'ol' : 'ul', 'ub-list');
            list.dataset.style = style;
            if (style === 'disc') {
                list.style.listStyle = 'disc';
            }
            for (const item of str(node, 'items').split('\n').map((s) => s.trim()).filter(Boolean)) {
                const li = el('li');
                li.textContent = item;
                list.append(li);
            }
            applyBox(list, node);
            return list;
        },
    },
    {
        type: 'badge', label: 'Badge', category: 'Content', icon: 'badge',
        defaults: {text: 'New in v3', tone: 'accent'},
        fields: [
            {key: 'text', label: 'Text', kind: 'text', group: 'Content'},
            {key: 'tone', label: 'Tone', kind: 'select', group: 'Style', options: [
                {label: 'Accent', value: 'accent'}, {label: 'Neutral', value: 'neutral'}, {label: 'Success', value: 'success'}, {label: 'Info', value: 'info'},
            ]},
        ],
        render(node) {
            const badge = el('span', 'ub-badge');
            badge.textContent = str(node, 'text');
            badge.dataset.tone = str(node, 'tone', 'accent');
            return badge;
        },
    },
    {
        type: 'metric', label: 'Metric', category: 'Content', icon: 'metric',
        defaults: {label: 'Active users', value: '12,480', delta: '+8.2% this week', trend: 'up', align: 'left'},
        fields: [
            {key: 'label', label: 'Label', kind: 'text', group: 'Content'},
            {key: 'value', label: 'Value', kind: 'text', group: 'Content'},
            {key: 'delta', label: 'Change', kind: 'text', group: 'Content'},
            {key: 'trend', label: 'Trend', kind: 'select', group: 'Style', options: [
                {label: 'Up', value: 'up'}, {label: 'Down', value: 'down'}, {label: 'Flat', value: 'flat'},
            ]},
            alignField,
        ],
        render(node) {
            const metric = el('div', 'ub-metric');
            const label = el('span', 'ub-metric-label');
            label.textContent = str(node, 'label');
            const value = el('span', 'ub-metric-value');
            value.textContent = str(node, 'value');
            metric.append(label, value);
            const delta = str(node, 'delta');
            if (delta) {
                const d = el('span', 'ub-metric-delta');
                d.textContent = delta;
                d.dataset.trend = str(node, 'trend', 'up');
                metric.append(d);
            }
            applyBox(metric, node);
            return metric;
        },
    },
    {
        type: 'quote', label: 'Quote', category: 'Content', icon: 'quote',
        defaults: {text: 'Wails let us ship a native desktop app with the web stack we already knew.', cite: 'A happy developer'},
        fields: [
            {key: 'text', label: 'Quote', kind: 'textarea', group: 'Content'},
            {key: 'cite', label: 'Attribution', kind: 'text', group: 'Content'},
            colorField,
        ],
        render(node) {
            const q = el('blockquote', 'ub-quote');
            q.append(document.createTextNode('“' + str(node, 'text') + '”'));
            const cite = str(node, 'cite');
            if (cite) {
                const footer = el('footer');
                footer.textContent = '— ' + cite;
                q.append(footer);
            }
            applyBox(q, node);
            return q;
        },
    },

    // ----- Form ------------------------------------------------------------
    {
        type: 'button', label: 'Button', category: 'Form', icon: 'button',
        defaults: {label: 'Get started', variant: 'primary', size: 'md', full: false, href: '#'},
        fields: [
            {key: 'label', label: 'Label', kind: 'text', group: 'Content'},
            {key: 'href', label: 'Link', kind: 'text', group: 'Content', placeholder: 'https://'},
            {key: 'variant', label: 'Variant', kind: 'select', group: 'Style', options: [
                {label: 'Primary', value: 'primary'}, {label: 'Secondary', value: 'secondary'}, {label: 'Ghost', value: 'ghost'}, {label: 'Danger', value: 'danger'},
            ]},
            {key: 'size', label: 'Size', kind: 'select', group: 'Style', options: [
                {label: 'Small', value: 'sm'}, {label: 'Medium', value: 'md'}, {label: 'Large', value: 'lg'},
            ]},
            {key: 'full', label: 'Full width', kind: 'toggle', group: 'Layout'},
        ],
        render(node) {
            const a = el('a', 'ub-btn');
            a.href = str(node, 'href', '#') || '#';
            a.textContent = str(node, 'label');
            a.dataset.variant = str(node, 'variant', 'primary');
            a.dataset.size = str(node, 'size', 'md');
            a.dataset.full = String(bool(node, 'full'));
            return a;
        },
    },
    {
        type: 'input', label: 'Input', category: 'Form', icon: 'input',
        defaults: {label: 'Email address', placeholder: 'you@example.com', type: 'email', help: ''},
        fields: [
            {key: 'label', label: 'Label', kind: 'text', group: 'Content'},
            {key: 'placeholder', label: 'Placeholder', kind: 'text', group: 'Content'},
            {key: 'help', label: 'Help text', kind: 'text', group: 'Content'},
            {key: 'type', label: 'Type', kind: 'select', group: 'Style', options: [
                {label: 'Text', value: 'text'}, {label: 'Email', value: 'email'}, {label: 'Password', value: 'password'}, {label: 'Number', value: 'number'}, {label: 'Search', value: 'search'},
            ]},
        ],
        render(node) {
            const field = el('div', 'ub-field');
            const label = str(node, 'label');
            if (label) {
                const l = el('label');
                l.textContent = label;
                field.append(l);
            }
            const input = el('input');
            input.type = str(node, 'type', 'text');
            input.placeholder = str(node, 'placeholder');
            field.append(input);
            const help = str(node, 'help');
            if (help) {
                const h = el('span', 'ub-help');
                h.textContent = help;
                field.append(h);
            }
            return field;
        },
    },
    {
        type: 'textarea', label: 'Textarea', category: 'Form', icon: 'textarea',
        defaults: {label: 'Message', placeholder: 'Tell us more…', rows: 4},
        fields: [
            {key: 'label', label: 'Label', kind: 'text', group: 'Content'},
            {key: 'placeholder', label: 'Placeholder', kind: 'text', group: 'Content'},
            {key: 'rows', label: 'Rows', kind: 'range', group: 'Layout', min: 2, max: 12, step: 1},
        ],
        render(node) {
            const field = el('div', 'ub-field');
            const label = str(node, 'label');
            if (label) {
                const l = el('label');
                l.textContent = label;
                field.append(l);
            }
            const ta = el('textarea');
            ta.placeholder = str(node, 'placeholder');
            ta.rows = num(node, 'rows', 4);
            field.append(ta);
            return field;
        },
    },
    {
        type: 'select', label: 'Select', category: 'Form', icon: 'select',
        defaults: {label: 'Plan', options: 'Hobby\nPro\nTeam'},
        fields: [
            {key: 'label', label: 'Label', kind: 'text', group: 'Content'},
            {key: 'options', label: 'Options (one per line)', kind: 'textarea', group: 'Content'},
        ],
        render(node) {
            const field = el('div', 'ub-field');
            const label = str(node, 'label');
            if (label) {
                const l = el('label');
                l.textContent = label;
                field.append(l);
            }
            const select = el('select');
            for (const opt of str(node, 'options').split('\n').map((s) => s.trim()).filter(Boolean)) {
                const o = el('option');
                o.textContent = opt;
                select.append(o);
            }
            field.append(select);
            return field;
        },
    },
    {
        type: 'toggle', label: 'Toggle', category: 'Form', icon: 'toggle',
        defaults: {label: 'Email me product updates', on: true},
        fields: [
            {key: 'label', label: 'Label', kind: 'text', group: 'Content'},
            {key: 'on', label: 'Switched on', kind: 'toggle', group: 'Content'},
        ],
        render(node) {
            const t = el('label', 'ub-toggle');
            t.dataset.on = String(bool(node, 'on'));
            const text = el('span');
            text.textContent = str(node, 'label');
            t.append(text, el('span', 'ub-switch'));
            return t;
        },
    },

    // ----- Media -----------------------------------------------------------
    {
        type: 'image', label: 'Image', category: 'Media', icon: 'image',
        defaults: {src: '', alt: 'Illustration', height: 240, radius: 14, fit: 'cover'},
        fields: [
            {key: 'src', label: 'Image URL', kind: 'text', group: 'Content', placeholder: 'https://… (empty shows a placeholder)'},
            {key: 'alt', label: 'Alt text', kind: 'text', group: 'Content'},
            {key: 'height', label: 'Height', kind: 'range', group: 'Layout', min: 80, max: 640, step: 8, unit: 'px'},
            radiusField,
            {key: 'fit', label: 'Fit', kind: 'select', group: 'Style', options: [
                {label: 'Cover', value: 'cover'}, {label: 'Contain', value: 'contain'},
            ]},
        ],
        render(node) {
            const src = str(node, 'src').trim();
            const height = num(node, 'height', 240);
            if (!src) {
                const ph = el('div', 'ub-image-placeholder');
                ph.style.height = `${height}px`;
                ph.style.borderRadius = `${num(node, 'radius', 14)}px`;
                ph.textContent = str(node, 'alt', 'Image');
                return ph;
            }
            const img = el('img', 'ub-image');
            img.src = src;
            img.alt = str(node, 'alt');
            img.style.width = '100%';
            img.style.height = `${height}px`;
            img.style.objectFit = str(node, 'fit', 'cover');
            img.style.borderRadius = `${num(node, 'radius', 14)}px`;
            return img;
        },
    },
    {
        type: 'divider', label: 'Divider', category: 'Media', icon: 'divider',
        defaults: {},
        fields: [],
        render() {
            return el('hr', 'ub-divider');
        },
    },
    {
        type: 'spacer', label: 'Spacer', category: 'Media', icon: 'spacer',
        defaults: {height: 32},
        fields: [{key: 'height', label: 'Height', kind: 'range', group: 'Layout', min: 4, max: 200, step: 4, unit: 'px'}],
        render(node) {
            const s = el('div', 'ub-spacer');
            s.style.height = `${num(node, 'height', 32)}px`;
            return s;
        },
    },
];

// The root is a pseudo-component: it's never in the palette, but it renders
// and accepts children like any other container.
export const rootDef: ComponentDef = {
    type: 'root', label: 'Page', category: 'Layout', icon: 'root',
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

export const categories: Category[] = ['Layout', 'Content', 'Form', 'Media'];

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
            d.style.cssText = 'padding:12px;border:1px dashed currentColor;border-radius:8px;opacity:.6;font-size:13px';
            d.textContent = `Unknown component “${type}”`;
            return d;
        },
    };
}
