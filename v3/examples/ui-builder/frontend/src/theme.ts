// Styles for the *built* application UI — the components a user drops into the
// window. This CSS is injected into the artboard at design time and written to
// the exported frontend's style.css, so what you build is exactly what ships.
//
// It is deliberately a desktop-app stylesheet, not a web-page one: 13px system
// type, dense controls, panes that fill the window and scroll internally.

export const componentStyles = `
.ub-root {
    --ub-bg: #f5f5f8;
    --ub-sidebar: #ebebf0;
    --ub-surface: #ffffff;
    --ub-surface-2: #f2f2f6;
    --ub-border: #d9d9e1;
    --ub-text: #1c1d24;
    --ub-muted: #6b6f80;
    --ub-accent: #ff4d4d;
    --ub-accent-2: #ff2d72;
    --ub-on-accent: #ffffff;
    --ub-ok: #1f9d5a;
    --ub-warn: #d98a00;
    --ub-danger: #d92d20;
    --ub-radius: 8px;
    --ub-font: -apple-system, BlinkMacSystemFont, "Segoe UI", system-ui, Roboto, "Helvetica Neue", Arial, sans-serif;
    --ub-mono: ui-monospace, "SF Mono", Menlo, Consolas, "Liberation Mono", monospace;
    display: flex;
    flex-direction: column;
    width: 100%;
    height: 100%;
    overflow: hidden;
    background: var(--ub-bg);
    color: var(--ub-text);
    font-family: var(--ub-font);
    font-size: 13px;
    line-height: 1.45;
    -webkit-font-smoothing: antialiased;
    -webkit-user-select: none;
    user-select: none;
}
.ub-root[data-theme="dark"] {
    --ub-bg: #1e1f26;
    --ub-sidebar: #25262f;
    --ub-surface: #2a2b35;
    --ub-surface-2: #31323d;
    --ub-border: #3a3b47;
    --ub-text: #f1f2f7;
    --ub-muted: #9a9eb0;
    --ub-ok: #3ccf7a;
    --ub-warn: #f2b134;
}
.ub-root * { box-sizing: border-box; }
.ub-root input, .ub-root textarea { -webkit-user-select: text; user-select: text; }
.ub-root ::-webkit-scrollbar { width: 10px; height: 10px; }
.ub-root ::-webkit-scrollbar-thumb { background: rgba(127, 127, 140, 0.35); border-radius: 99px; border: 2px solid transparent; background-clip: padding-box; }

/* ----- Layout ----------------------------------------------------------- */
.ub-hstack { display: flex; flex-direction: row; min-height: 0; min-width: 0; }
.ub-vstack { display: flex; flex-direction: column; min-height: 0; min-width: 0; }
.ub-hstack > *, .ub-vstack > * { min-width: 0; }
[data-fill="true"] { flex: 1 1 0; }
[data-scroll="true"] { overflow: auto; }
.ub-sidebar {
    display: flex; flex-direction: column; gap: 4px; flex: 0 0 auto;
    background: var(--ub-sidebar); border-right: 1px solid var(--ub-border);
    overflow: auto;
}
.ub-sidebar[data-side="right"] { border-right: 0; border-left: 1px solid var(--ub-border); }
.ub-toolbar {
    display: flex; align-items: center; gap: 8px; flex: 0 0 auto;
    padding: 0 12px; background: var(--ub-surface); border-bottom: 1px solid var(--ub-border);
}
.ub-toolbar[data-drag="true"] { --wails-draggable: drag; }
.ub-toolbar[data-drag="true"] > * { --wails-draggable: no-drag; }
.ub-content { display: flex; flex-direction: column; flex: 1 1 0; min-height: 0; overflow: auto; }
.ub-group { display: flex; flex-direction: column; gap: 10px; padding: 12px 14px 14px; background: var(--ub-surface); border: 1px solid var(--ub-border); border-radius: var(--ub-radius); }
.ub-group > .ub-group-title { margin: 0 0 2px; font-size: 12px; font-weight: 600; color: var(--ub-muted); text-transform: uppercase; letter-spacing: 0.06em; }
.ub-statusbar {
    display: flex; align-items: center; gap: 14px; flex: 0 0 auto; height: 26px; padding: 0 10px;
    background: var(--ub-surface); border-top: 1px solid var(--ub-border); font-size: 11.5px; color: var(--ub-muted);
}
.ub-statusbar > [data-push="true"] { margin-left: auto; }
.ub-tabs { display: flex; gap: 2px; flex: 0 0 auto; border-bottom: 1px solid var(--ub-border); }
.ub-tabs > button {
    font: inherit; color: var(--ub-muted); background: none; border: none; border-bottom: 2px solid transparent;
    padding: 8px 12px 7px; cursor: pointer; font-weight: 500;
}
.ub-tabs > button[aria-selected="true"] { color: var(--ub-text); border-bottom-color: var(--ub-accent); }
.ub-tabs[data-style="segmented"] { border: none; background: var(--ub-surface-2); padding: 3px; border-radius: 7px; width: fit-content; }
.ub-tabs[data-style="segmented"] > button { border: none; border-radius: 5px; padding: 4px 12px; }
.ub-tabs[data-style="segmented"] > button[aria-selected="true"] { background: var(--ub-surface); box-shadow: 0 1px 2px rgba(0,0,0,0.15); }
.ub-nav { display: flex; flex-direction: column; gap: 1px; padding: 6px; margin: 0; list-style: none; }
.ub-nav > li { display: flex; align-items: center; gap: 9px; padding: 6px 9px; border-radius: 6px; cursor: default; font-weight: 500; color: var(--ub-text); }
.ub-nav > li::before { content: ""; width: 14px; height: 14px; border-radius: 4px; background: currentColor; opacity: 0.35; flex: 0 0 auto; }
.ub-nav > li[aria-current="true"] { background: var(--ub-accent); color: var(--ub-on-accent); }
.ub-nav > li[aria-current="true"]::before { opacity: 0.9; }
.ub-nav[data-style="plain"] > li[aria-current="true"] { background: rgba(127, 127, 140, 0.18); color: var(--ub-text); }
.ub-nav > li.ub-nav-header { font-size: 11px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.06em; color: var(--ub-muted); margin-top: 8px; cursor: default; }
.ub-nav > li.ub-nav-header::before { display: none; }
.ub-divider { border: none; border-top: 1px solid var(--ub-border); margin: 0; width: 100%; flex: 0 0 auto; }
.ub-spacer { flex: 0 0 auto; }
.ub-spacer[data-fill="true"] { flex: 1 1 auto; }

/* ----- Content ---------------------------------------------------------- */
.ub-heading { margin: 0; font-weight: 600; letter-spacing: -0.01em; }
.ub-heading[data-level="window"] { font-size: 13px; font-weight: 600; }
.ub-heading[data-level="title"] { font-size: 20px; line-height: 1.2; }
.ub-heading[data-level="section"] { font-size: 15px; }
.ub-heading[data-level="label"] { font-size: 11px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.06em; color: var(--ub-muted); }
.ub-text { margin: 0; }
.ub-text[data-size="sm"] { font-size: 11.5px; }
.ub-text[data-size="lg"] { font-size: 15px; }
.ub-text[data-tone="muted"] { color: var(--ub-muted); }
.ub-text[data-tone="ok"] { color: var(--ub-ok); }
.ub-text[data-tone="warn"] { color: var(--ub-warn); }
.ub-text[data-tone="danger"] { color: var(--ub-danger); }
.ub-text[data-mono="true"] { font-family: var(--ub-mono); font-size: 12px; }
.ub-badge {
    display: inline-flex; align-items: center; width: fit-content; padding: 1px 8px; border-radius: 999px;
    font-size: 11px; font-weight: 600; letter-spacing: 0.02em;
}
.ub-badge[data-tone="accent"] { color: var(--ub-accent); background: rgba(255, 77, 77, 0.14); }
.ub-badge[data-tone="neutral"] { color: var(--ub-muted); background: rgba(127, 127, 140, 0.16); }
.ub-badge[data-tone="ok"] { color: var(--ub-ok); background: rgba(31, 157, 90, 0.14); }
.ub-badge[data-tone="warn"] { color: var(--ub-warn); background: rgba(217, 138, 0, 0.14); }
.ub-badge[data-tone="danger"] { color: var(--ub-danger); background: rgba(217, 45, 32, 0.14); }
.ub-metric { display: flex; flex-direction: column; gap: 2px; padding: 12px 14px; background: var(--ub-surface); border: 1px solid var(--ub-border); border-radius: var(--ub-radius); }
.ub-metric > .ub-metric-label { font-size: 11px; font-weight: 600; color: var(--ub-muted); text-transform: uppercase; letter-spacing: 0.06em; }
.ub-metric > .ub-metric-value { font-size: 24px; font-weight: 600; letter-spacing: -0.02em; line-height: 1.15; font-variant-numeric: tabular-nums; }
.ub-metric > .ub-metric-note { font-size: 11.5px; color: var(--ub-muted); }
.ub-table { width: 100%; border-collapse: collapse; font-size: 12.5px; background: var(--ub-surface); }
.ub-table th, .ub-table td { padding: 6px 10px; text-align: left; border-bottom: 1px solid var(--ub-border); white-space: nowrap; }
.ub-table th { font-size: 11px; font-weight: 600; color: var(--ub-muted); text-transform: uppercase; letter-spacing: 0.05em; background: var(--ub-surface-2); position: sticky; top: 0; }
.ub-table[data-striped="true"] tbody tr:nth-child(even) { background: rgba(127, 127, 140, 0.07); }
.ub-table tbody tr[aria-selected="true"] { background: rgba(255, 77, 77, 0.14); }
.ub-table-wrap { overflow: auto; border: 1px solid var(--ub-border); border-radius: var(--ub-radius); flex: 0 0 auto; }
.ub-table-wrap[data-fill="true"] { flex: 1 1 0; min-height: 0; }
.ub-list { margin: 0; padding: 0; list-style: none; display: flex; flex-direction: column; background: var(--ub-surface); border: 1px solid var(--ub-border); border-radius: var(--ub-radius); overflow: hidden; }
.ub-list > li { padding: 7px 12px; border-bottom: 1px solid var(--ub-border); }
.ub-list > li:last-child { border-bottom: 0; }
.ub-list[data-style="bullets"] { background: none; border: none; gap: 3px; }
.ub-list[data-style="bullets"] > li { padding: 0 0 0 16px; border: 0; position: relative; }
.ub-list[data-style="bullets"] > li::before { content: ""; position: absolute; left: 4px; top: 8px; width: 5px; height: 5px; border-radius: 50%; background: var(--ub-muted); }
.ub-keyvalue { display: grid; grid-template-columns: max-content 1fr; gap: 6px 16px; margin: 0; font-size: 12.5px; }
.ub-keyvalue dt { color: var(--ub-muted); }
.ub-keyvalue dd { margin: 0; font-variant-numeric: tabular-nums; }
.ub-progress { display: flex; flex-direction: column; gap: 5px; }
.ub-progress > .ub-progress-head { display: flex; justify-content: space-between; font-size: 11.5px; color: var(--ub-muted); }
.ub-progress > .ub-progress-track { height: 6px; border-radius: 99px; background: rgba(127, 127, 140, 0.22); overflow: hidden; }
.ub-progress > .ub-progress-track > i { display: block; height: 100%; border-radius: 99px; background: linear-gradient(90deg, var(--ub-accent), var(--ub-accent-2)); transition: width 0.3s ease; }
.ub-log {
    margin: 0; padding: 10px 12px; flex: 0 0 auto; overflow: auto; white-space: pre-wrap; word-break: break-word;
    font-family: var(--ub-mono); font-size: 12px; line-height: 1.5; color: var(--ub-text);
    background: var(--ub-surface-2); border: 1px solid var(--ub-border); border-radius: var(--ub-radius);
    -webkit-user-select: text; user-select: text;
}
.ub-log[data-fill="true"] { flex: 1 1 0; min-height: 0; }
.ub-image { display: block; max-width: 100%; border-radius: var(--ub-radius); object-fit: cover; }
.ub-image-placeholder {
    display: flex; align-items: center; justify-content: center; border-radius: var(--ub-radius);
    background: repeating-linear-gradient(45deg, transparent 0 8px, rgba(127, 127, 140, 0.1) 8px 16px), var(--ub-surface-2);
    border: 1px dashed var(--ub-border); color: var(--ub-muted); font-size: 12px;
}

/* ----- Controls --------------------------------------------------------- */
.ub-btn {
    display: inline-flex; align-items: center; justify-content: center; gap: 6px; flex: 0 0 auto;
    height: 28px; padding: 0 12px; border-radius: 6px; border: 1px solid transparent; cursor: pointer;
    font: inherit; font-weight: 500; white-space: nowrap; width: fit-content;
    transition: background 0.12s ease, box-shadow 0.12s ease, transform 0.08s ease;
}
.ub-btn[data-size="sm"] { height: 24px; padding: 0 9px; font-size: 12px; }
.ub-btn[data-size="lg"] { height: 34px; padding: 0 18px; font-size: 14px; }
.ub-btn[data-full="true"] { width: 100%; }
.ub-btn[data-variant="primary"] { color: var(--ub-on-accent); background: linear-gradient(180deg, #ff5a4d, #f52d5f); box-shadow: 0 1px 2px rgba(0,0,0,0.2); }
.ub-btn[data-variant="primary"]:hover { filter: brightness(1.06); }
.ub-btn[data-variant="default"] { color: var(--ub-text); background: var(--ub-surface); border-color: var(--ub-border); box-shadow: 0 1px 1px rgba(0,0,0,0.06); }
.ub-btn[data-variant="default"]:hover { background: var(--ub-surface-2); }
.ub-btn[data-variant="ghost"] { color: var(--ub-text); background: transparent; }
.ub-btn[data-variant="ghost"]:hover { background: rgba(127, 127, 140, 0.16); }
.ub-btn[data-variant="danger"] { color: #fff; background: var(--ub-danger); }
.ub-btn:active { transform: translateY(1px); }
.ub-btn[data-busy="true"] { opacity: 0.6; pointer-events: none; }
.ub-field { display: flex; flex-direction: column; gap: 4px; flex: 0 0 auto; }
.ub-field[data-fill="true"] { flex: 1 1 0; }
.ub-field[data-inline="true"] { flex-direction: row; align-items: center; gap: 10px; }
.ub-field[data-inline="true"] > label { min-width: 110px; }
.ub-field > label { font-size: 12px; font-weight: 500; color: var(--ub-text); }
.ub-field > .ub-help { font-size: 11px; color: var(--ub-muted); }
.ub-field input, .ub-field textarea, .ub-field select {
    width: 100%; font: inherit; color: var(--ub-text); padding: 5px 8px; border-radius: 6px;
    background: var(--ub-surface); border: 1px solid var(--ub-border); outline: none;
    transition: border-color 0.12s ease, box-shadow 0.12s ease;
}
.ub-field input { height: 28px; }
.ub-field textarea { resize: vertical; min-height: 64px; }
.ub-field input:focus, .ub-field textarea:focus, .ub-field select:focus { border-color: var(--ub-accent); box-shadow: 0 0 0 3px rgba(255, 77, 77, 0.18); }
.ub-field input[type="search"] { padding-left: 26px; background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%23888c9c' stroke-width='2' stroke-linecap='round'%3E%3Ccircle cx='11' cy='11' r='7'/%3E%3Cpath d='m20 20-3.5-3.5'/%3E%3C/svg%3E"); background-repeat: no-repeat; background-size: 14px; background-position: 8px center; }
.ub-check { display: inline-flex; align-items: center; gap: 8px; cursor: pointer; width: fit-content; flex: 0 0 auto; }
.ub-check > i { width: 15px; height: 15px; border-radius: 4px; border: 1px solid var(--ub-border); background: var(--ub-surface); position: relative; flex: 0 0 auto; }
.ub-check[data-on="true"] > i { background: var(--ub-accent); border-color: var(--ub-accent); }
.ub-check[data-on="true"] > i::after { content: ""; position: absolute; left: 4px; top: 1px; width: 4px; height: 8px; border: solid #fff; border-width: 0 2px 2px 0; transform: rotate(45deg); }
.ub-toggle { display: inline-flex; align-items: center; justify-content: space-between; gap: 12px; cursor: pointer; flex: 0 0 auto; }
.ub-toggle > .ub-switch { position: relative; width: 34px; height: 20px; border-radius: 999px; background: rgba(127, 127, 140, 0.35); transition: background 0.15s ease; flex: 0 0 auto; }
.ub-toggle > .ub-switch::after { content: ""; position: absolute; top: 2px; left: 2px; width: 16px; height: 16px; border-radius: 50%; background: #fff; box-shadow: 0 1px 2px rgba(0,0,0,0.3); transition: transform 0.15s ease; }
.ub-toggle[data-on="true"] > .ub-switch { background: var(--ub-ok); }
.ub-toggle[data-on="true"] > .ub-switch::after { transform: translateX(14px); }
.ub-slider { display: flex; flex-direction: column; gap: 4px; flex: 0 0 auto; }
.ub-slider > .ub-slider-head { display: flex; justify-content: space-between; font-size: 12px; }
.ub-slider > .ub-slider-head > output { color: var(--ub-muted); font-variant-numeric: tabular-nums; }
.ub-slider input[type="range"] { -webkit-appearance: none; appearance: none; width: 100%; height: 4px; border-radius: 2px; background: rgba(127, 127, 140, 0.3); outline: none; }
.ub-slider input[type="range"]::-webkit-slider-thumb { -webkit-appearance: none; width: 16px; height: 16px; border-radius: 50%; background: #fff; border: 1px solid var(--ub-border); box-shadow: 0 1px 3px rgba(0,0,0,0.3); cursor: pointer; }
`;
