// Inline SVG icons for the palette and layers panel. All share the same
// 24×24 viewBox and inherit `stroke: currentColor` from the stylesheet.

const wrap = (paths: string): string =>
    `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">${paths}</svg>`;

export const icons: Record<string, string> = {
    section: wrap('<rect x="3" y="4" width="18" height="16" rx="2"/><path d="M3 10h18"/>'),
    row: wrap('<rect x="3" y="6" width="5" height="12" rx="1"/><rect x="9.5" y="6" width="5" height="12" rx="1"/><rect x="16" y="6" width="5" height="12" rx="1"/>'),
    card: wrap('<rect x="4" y="4" width="16" height="16" rx="3"/><path d="M8 9h8M8 13h5"/>'),
    navbar: wrap('<rect x="3" y="4" width="18" height="6" rx="1.5"/><path d="M6 7h3M13 7h1.5M17 7h1.5"/><path d="M3 14h18M3 18h12" opacity=".4"/>'),
    heading: wrap('<path d="M5 5v14M13 5v14M5 12h8"/><path d="M17 12h3M18.5 10.5v8.5" />'),
    text: wrap('<path d="M4 6h16M4 10h16M4 14h11M4 18h7"/>'),
    list: wrap('<path d="M9 6h11M9 12h11M9 18h11"/><circle cx="5" cy="6" r="1"/><circle cx="5" cy="12" r="1"/><circle cx="5" cy="18" r="1"/>'),
    badge: wrap('<rect x="3" y="8" width="18" height="8" rx="4"/><path d="M8 12h8"/>'),
    metric: wrap('<path d="M4 19h16"/><path d="M6 15l4-5 3 3 5-7"/>'),
    quote: wrap('<path d="M7 7h3v4c0 2-1 3-3 3"/><path d="M14 7h3v4c0 2-1 3-3 3"/>'),
    button: wrap('<rect x="3" y="8" width="18" height="8" rx="2.5"/><path d="M8 12h8"/>'),
    input: wrap('<rect x="3" y="7" width="18" height="10" rx="2"/><path d="M7 11v2"/>'),
    textarea: wrap('<rect x="3" y="5" width="18" height="14" rx="2"/><path d="M7 9h10M7 13h6"/>'),
    select: wrap('<rect x="3" y="7" width="18" height="10" rx="2"/><path d="m15 11 2 2 2-2"/>'),
    toggle: wrap('<rect x="3" y="8" width="18" height="8" rx="4"/><circle cx="15" cy="12" r="2.5"/>'),
    image: wrap('<rect x="3" y="4" width="18" height="16" rx="2"/><circle cx="9" cy="10" r="1.5"/><path d="m21 16-5-5-8 8"/>'),
    divider: wrap('<path d="M4 12h16"/><path d="M4 7h4M16 7h4M4 17h4M16 17h4" opacity=".35"/>'),
    spacer: wrap('<path d="M4 5h16M4 19h16"/><path d="M12 8v8m-3-3 3 3 3-3m-3-8-3 3m3-3 3 3" />'),
    root: wrap('<rect x="3" y="3" width="18" height="18" rx="3"/>'),
    // Generic UI glyphs
    chevron: wrap('<path d="m9 6 6 6-6 6"/>'),
    trash: wrap('<path d="M4 7h16M10 11v6M14 11v6M6 7l1 13h10l1-13M9 7V4h6v3"/>'),
    copy: wrap('<rect x="9" y="9" width="11" height="11" rx="2"/><path d="M5 15V6a2 2 0 0 1 2-2h9"/>'),
    up: wrap('<path d="m6 14 6-6 6 6"/>'),
    down: wrap('<path d="m6 10 6 6 6-6"/>'),
    parent: wrap('<path d="M12 19V7"/><path d="m6 13 6-6 6 6"/><path d="M4 4h16"/>'),
    cursor: wrap('<path d="m5 4 14 7-6 2-2 6z"/>'),
};

export function icon(name: string): string {
    return icons[name] ?? icons.cursor;
}
