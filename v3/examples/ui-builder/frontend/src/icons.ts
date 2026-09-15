// Inline SVG icons for the palette, inspector and layers panel. All share the
// same 24×24 viewBox and inherit `stroke: currentColor` from the stylesheet.

const wrap = (paths: string): string =>
    `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">${paths}</svg>`;

export const icons: Record<string, string> = {
    // Layout
    hstack: wrap('<rect x="3" y="5" width="5.5" height="14" rx="1"/><rect x="9.5" y="5" width="5.5" height="14" rx="1"/><rect x="16" y="5" width="5" height="14" rx="1"/>'),
    vstack: wrap('<rect x="4" y="3" width="16" height="5" rx="1"/><rect x="4" y="9.5" width="16" height="5" rx="1"/><rect x="4" y="16" width="16" height="5" rx="1"/>'),
    sidebar: wrap('<rect x="3" y="4" width="18" height="16" rx="2"/><path d="M9 4v16"/><path d="M5 8h2M5 11h2M5 14h2"/>'),
    toolbar: wrap('<rect x="3" y="4" width="18" height="16" rx="2"/><path d="M3 9h18"/><path d="M6 6.5h3M12 6.5h1.5M16 6.5h2"/>'),
    content: wrap('<rect x="3" y="4" width="18" height="16" rx="2"/><path d="M7 9h10M7 12h10M7 15h6" opacity=".5"/>'),
    group: wrap('<rect x="3" y="6" width="18" height="14" rx="2"/><path d="M6 6V4h6v2"/><path d="M7 11h10M7 15h6" opacity=".5"/>'),
    statusbar: wrap('<rect x="3" y="4" width="18" height="16" rx="2"/><path d="M3 16h18"/><path d="M6 18h4M15 18h3"/>'),
    tabs: wrap('<path d="M3 9h18v10a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1z"/><path d="M3 9V6a1 1 0 0 1 1-1h5l2 4M11 9l2-4h5a1 1 0 0 1 1 1v3"/>'),
    nav: wrap('<path d="M8 6h12M8 12h12M8 18h12"/><rect x="3" y="4.5" width="3" height="3" rx="1"/><rect x="3" y="10.5" width="3" height="3" rx="1" fill="currentColor"/><rect x="3" y="16.5" width="3" height="3" rx="1"/>'),
    divider: wrap('<path d="M4 12h16"/><path d="M4 7h4M16 7h4M4 17h4M16 17h4" opacity=".35"/>'),
    spacer: wrap('<path d="M4 5h16M4 19h16"/><path d="M12 8v8m-3-3 3 3 3-3m-3-8-3 3m3-3 3 3"/>'),
    // Content
    heading: wrap('<path d="M5 5v14M13 5v14M5 12h8"/><path d="M17 12h3M18.5 10.5v8.5"/>'),
    text: wrap('<path d="M4 6h16M4 10h16M4 14h11M4 18h7"/>'),
    badge: wrap('<rect x="3" y="8" width="18" height="8" rx="4"/><path d="M8 12h8"/>'),
    metric: wrap('<path d="M4 19h16"/><path d="M6 15l4-5 3 3 5-7"/>'),
    table: wrap('<rect x="3" y="4" width="18" height="16" rx="2"/><path d="M3 9h18M3 14h18M9 9v11M15 9v11"/>'),
    list: wrap('<path d="M9 6h11M9 12h11M9 18h11"/><circle cx="5" cy="6" r="1"/><circle cx="5" cy="12" r="1"/><circle cx="5" cy="18" r="1"/>'),
    keyvalue: wrap('<path d="M4 7h6M4 12h6M4 17h6"/><path d="M14 7h6M14 12h4M14 17h5" opacity=".5"/>'),
    progress: wrap('<rect x="3" y="9" width="18" height="6" rx="3"/><rect x="3" y="9" width="11" height="6" rx="3" fill="currentColor"/>'),
    log: wrap('<rect x="3" y="4" width="18" height="16" rx="2"/><path d="m7 9 3 3-3 3M13 15h4"/>'),
    image: wrap('<rect x="3" y="4" width="18" height="16" rx="2"/><circle cx="9" cy="10" r="1.5"/><path d="m21 16-5-5-8 8"/>'),
    // Controls
    button: wrap('<rect x="3" y="8" width="18" height="8" rx="2.5"/><path d="M8 12h8"/>'),
    input: wrap('<rect x="3" y="7" width="18" height="10" rx="2"/><path d="M7 11v2"/>'),
    textarea: wrap('<rect x="3" y="5" width="18" height="14" rx="2"/><path d="M7 9h10M7 13h6"/>'),
    select: wrap('<rect x="3" y="7" width="18" height="10" rx="2"/><path d="m15 11 2 2 2-2"/>'),
    toggle: wrap('<rect x="3" y="8" width="18" height="8" rx="4"/><circle cx="15" cy="12" r="2.5"/>'),
    checkbox: wrap('<rect x="4" y="4" width="16" height="16" rx="3"/><path d="m8 12 3 3 5-6"/>'),
    slider: wrap('<path d="M4 12h16"/><circle cx="14" cy="12" r="3" fill="currentColor"/>'),
    root: wrap('<rect x="3" y="3" width="18" height="18" rx="3"/>'),
    // Generic UI glyphs
    chevron: wrap('<path d="m9 6 6 6-6 6"/>'),
    trash: wrap('<path d="M4 7h16M10 11v6M14 11v6M6 7l1 13h10l1-13M9 7V4h6v3"/>'),
    copy: wrap('<rect x="9" y="9" width="11" height="11" rx="2"/><path d="M5 15V6a2 2 0 0 1 2-2h9"/>'),
    up: wrap('<path d="m6 14 6-6 6 6"/>'),
    down: wrap('<path d="m6 10 6 6 6-6"/>'),
    parent: wrap('<path d="M12 19V7"/><path d="m6 13 6-6 6 6"/><path d="M4 4h16"/>'),
    cursor: wrap('<path d="m5 4 14 7-6 2-2 6z"/>'),
    go: wrap('<circle cx="12" cy="12" r="9"/><path d="M8 12h8m-3-3 3 3-3 3"/>'),
};

export function icon(name: string): string {
    return icons[name] ?? icons.cursor;
}
