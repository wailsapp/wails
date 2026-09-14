// Styles for the *built* UI — the components a user drops onto the canvas.
// This CSS is injected into the artboard at design time and inlined into the
// exported HTML page, so what you build is exactly what you ship.

export const componentStyles = `
.ub-root {
    --ub-bg: #ffffff;
    --ub-surface: #f6f7fb;
    --ub-surface-2: #ffffff;
    --ub-border: #e4e7ef;
    --ub-text: #0f1222;
    --ub-muted: #5f6b85;
    --ub-accent: #ff4d4d;
    --ub-accent-2: #ff2d72;
    --ub-accent-text: #ffffff;
    --ub-radius: 12px;
    --ub-shadow: 0 1px 2px rgba(15, 18, 34, 0.06), 0 8px 24px rgba(15, 18, 34, 0.06);
    --ub-font: "Inter", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
    background: var(--ub-bg);
    color: var(--ub-text);
    font-family: var(--ub-font);
    font-size: 16px;
    line-height: 1.5;
    -webkit-font-smoothing: antialiased;
    display: flex;
    flex-direction: column;
    min-height: 100%;
}
.ub-root[data-theme="dark"] {
    --ub-bg: #0b0d16;
    --ub-surface: #131627;
    --ub-surface-2: #181c30;
    --ub-border: rgba(255, 255, 255, 0.1);
    --ub-text: #f4f6fb;
    --ub-muted: #9aa6c0;
    --ub-shadow: 0 1px 2px rgba(0, 0, 0, 0.4), 0 12px 32px rgba(0, 0, 0, 0.35);
}
.ub-root * { box-sizing: border-box; }
.ub-root img { max-width: 100%; display: block; }

/* Layout */
.ub-section { display: flex; flex-direction: column; width: 100%; }
.ub-section > .ub-inner { width: 100%; margin: 0 auto; display: flex; flex-direction: column; }
.ub-row { display: flex; width: 100%; }
.ub-row > * { min-width: 0; }
.ub-row[data-equal="true"] > * { flex: 1 1 0; }
.ub-card {
    display: flex; flex-direction: column;
    background: var(--ub-surface-2);
    border: 1px solid var(--ub-border);
    border-radius: var(--ub-radius);
}
.ub-card[data-elevated="true"] { box-shadow: var(--ub-shadow); }
.ub-navbar {
    display: flex; align-items: center; justify-content: space-between; gap: 24px;
    width: 100%; padding: 14px 24px;
    border-bottom: 1px solid var(--ub-border);
    background: var(--ub-bg);
}
.ub-navbar .ub-navbar-brand { font-weight: 700; letter-spacing: -0.01em; display: inline-flex; align-items: center; gap: 10px; }
.ub-navbar .ub-navbar-brand i { width: 22px; height: 22px; border-radius: 7px; background: linear-gradient(135deg, var(--ub-accent), var(--ub-accent-2)); }
.ub-navbar nav { display: flex; gap: 22px; }
.ub-navbar nav a { color: var(--ub-muted); text-decoration: none; font-size: 14.5px; font-weight: 500; }
.ub-navbar nav a:hover { color: var(--ub-text); }

/* Content */
.ub-heading { margin: 0; font-weight: 800; letter-spacing: -0.025em; line-height: 1.1; }
.ub-heading.ub-h1 { font-size: 52px; }
.ub-heading.ub-h2 { font-size: 36px; }
.ub-heading.ub-h3 { font-size: 26px; letter-spacing: -0.015em; }
.ub-heading.ub-h4 { font-size: 19px; letter-spacing: -0.01em; }
.ub-heading[data-gradient="true"] {
    background: linear-gradient(135deg, var(--ub-accent) 0%, var(--ub-accent-2) 100%);
    -webkit-background-clip: text; background-clip: text;
    -webkit-text-fill-color: transparent; color: transparent;
    width: fit-content; max-width: 100%;
}
.ub-text { margin: 0; }
.ub-text[data-size="sm"] { font-size: 14px; }
.ub-text[data-size="md"] { font-size: 16px; }
.ub-text[data-size="lg"] { font-size: 19px; line-height: 1.55; }
.ub-text[data-muted="true"] { color: var(--ub-muted); }
.ub-list { margin: 0; padding-left: 1.3em; display: flex; flex-direction: column; gap: 6px; }
.ub-list[data-style="check"] { list-style: none; padding-left: 0; }
.ub-list[data-style="check"] li { display: flex; gap: 10px; align-items: flex-start; }
.ub-list[data-style="check"] li::before {
    content: ""; flex: 0 0 auto; width: 18px; height: 18px; margin-top: 3px; border-radius: 50%;
    background: linear-gradient(135deg, var(--ub-accent), var(--ub-accent-2));
    -webkit-mask: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%23000' stroke-width='3' stroke-linecap='round' stroke-linejoin='round'%3E%3Ccircle cx='12' cy='12' r='10' fill='%23000' stroke='none'/%3E%3Cpath d='m8 12 3 3 5-6' stroke='%23fff'/%3E%3C/svg%3E") center / contain no-repeat;
    mask: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%23000' stroke-width='3' stroke-linecap='round' stroke-linejoin='round'%3E%3Ccircle cx='12' cy='12' r='10' fill='%23000' stroke='none'/%3E%3Cpath d='m8 12 3 3 5-6' stroke='%23fff'/%3E%3C/svg%3E") center / contain no-repeat;
}
.ub-badge {
    display: inline-flex; align-items: center; width: fit-content;
    padding: 4px 11px; border-radius: 999px;
    font-size: 12px; font-weight: 700; letter-spacing: 0.04em; text-transform: uppercase;
}
.ub-badge[data-tone="accent"] { color: var(--ub-accent); background: rgba(255, 77, 77, 0.12); }
.ub-badge[data-tone="neutral"] { color: var(--ub-muted); background: var(--ub-surface); border: 1px solid var(--ub-border); }
.ub-badge[data-tone="success"] { color: #12855a; background: rgba(18, 133, 90, 0.12); }
.ub-badge[data-tone="info"] { color: #2563eb; background: rgba(37, 99, 235, 0.12); }
.ub-root[data-theme="dark"] .ub-badge[data-tone="success"] { color: #4ade80; }
.ub-root[data-theme="dark"] .ub-badge[data-tone="info"] { color: #7db1ff; }
.ub-metric { display: flex; flex-direction: column; gap: 4px; }
.ub-metric .ub-metric-label { font-size: 13px; font-weight: 600; color: var(--ub-muted); text-transform: uppercase; letter-spacing: 0.06em; }
.ub-metric .ub-metric-value { font-size: 34px; font-weight: 800; letter-spacing: -0.03em; line-height: 1.1; }
.ub-metric .ub-metric-delta { font-size: 13px; font-weight: 600; }
.ub-metric .ub-metric-delta[data-trend="up"] { color: #12855a; }
.ub-metric .ub-metric-delta[data-trend="down"] { color: var(--ub-accent); }
.ub-root[data-theme="dark"] .ub-metric .ub-metric-delta[data-trend="up"] { color: #4ade80; }
.ub-quote { margin: 0; padding-left: 18px; border-left: 3px solid var(--ub-accent); font-size: 19px; line-height: 1.5; }
.ub-quote footer { margin-top: 10px; font-size: 14px; color: var(--ub-muted); }
.ub-divider { border: none; border-top: 1px solid var(--ub-border); margin: 0; width: 100%; }
.ub-spacer { width: 100%; flex: 0 0 auto; }
.ub-image { border-radius: var(--ub-radius); object-fit: cover; background: var(--ub-surface); }
.ub-image-placeholder {
    display: flex; align-items: center; justify-content: center; width: 100%;
    border-radius: var(--ub-radius); color: var(--ub-muted);
    background:
        linear-gradient(135deg, rgba(255, 77, 77, 0.14), rgba(255, 45, 114, 0.06)),
        repeating-linear-gradient(45deg, transparent 0 10px, rgba(127, 127, 127, 0.06) 10px 20px);
    border: 1px dashed var(--ub-border);
    font-size: 13px; font-weight: 600;
}

/* Form */
.ub-btn {
    display: inline-flex; align-items: center; justify-content: center; gap: 8px;
    border: 1px solid transparent; border-radius: 10px; cursor: pointer;
    font-family: inherit; font-weight: 600; text-decoration: none;
    transition: transform 0.1s ease, box-shadow 0.2s ease, background 0.2s ease;
    width: fit-content; white-space: nowrap;
}
.ub-btn[data-size="sm"] { padding: 7px 14px; font-size: 13.5px; }
.ub-btn[data-size="md"] { padding: 10px 20px; font-size: 15px; }
.ub-btn[data-size="lg"] { padding: 14px 28px; font-size: 17px; border-radius: 12px; }
.ub-btn[data-full="true"] { width: 100%; }
.ub-btn[data-variant="primary"] { color: var(--ub-accent-text); background: linear-gradient(135deg, var(--ub-accent) 0%, var(--ub-accent-2) 100%); box-shadow: 0 6px 18px rgba(255, 45, 114, 0.3); }
.ub-btn[data-variant="primary"]:hover { box-shadow: 0 8px 24px rgba(255, 45, 114, 0.45); }
.ub-btn[data-variant="secondary"] { color: var(--ub-text); background: var(--ub-surface); border-color: var(--ub-border); }
.ub-btn[data-variant="secondary"]:hover { background: var(--ub-surface-2); }
.ub-btn[data-variant="ghost"] { color: var(--ub-text); background: transparent; }
.ub-btn[data-variant="ghost"]:hover { background: var(--ub-surface); }
.ub-btn[data-variant="danger"] { color: #fff; background: #d92d20; }
.ub-btn:active { transform: scale(0.98); }
.ub-field { display: flex; flex-direction: column; gap: 6px; width: 100%; }
.ub-field label { font-size: 13.5px; font-weight: 600; color: var(--ub-text); }
.ub-field .ub-help { font-size: 12.5px; color: var(--ub-muted); }
.ub-field input, .ub-field textarea, .ub-field select {
    width: 100%; font: inherit; font-size: 15px; color: var(--ub-text);
    padding: 10px 12px; border-radius: 10px;
    background: var(--ub-surface-2); border: 1px solid var(--ub-border);
    outline: none; transition: border-color 0.15s ease, box-shadow 0.15s ease;
}
.ub-field textarea { resize: vertical; min-height: 80px; }
.ub-field input:focus, .ub-field textarea:focus, .ub-field select:focus { border-color: var(--ub-accent); box-shadow: 0 0 0 3px rgba(255, 77, 77, 0.18); }
.ub-toggle { display: inline-flex; align-items: center; gap: 12px; cursor: pointer; width: fit-content; }
.ub-toggle span:first-child { font-size: 15px; font-weight: 500; }
.ub-toggle .ub-switch { position: relative; width: 42px; height: 24px; border-radius: 999px; background: var(--ub-border); transition: background 0.2s ease; flex: 0 0 auto; }
.ub-toggle .ub-switch::after { content: ""; position: absolute; top: 3px; left: 3px; width: 18px; height: 18px; border-radius: 50%; background: #fff; box-shadow: 0 1px 3px rgba(0,0,0,0.3); transition: transform 0.2s ease; }
.ub-toggle[data-on="true"] .ub-switch { background: linear-gradient(135deg, var(--ub-accent), var(--ub-accent-2)); }
.ub-toggle[data-on="true"] .ub-switch::after { transform: translateX(18px); }

@media (max-width: 640px) {
    .ub-heading.ub-h1 { font-size: 38px; }
    .ub-heading.ub-h2 { font-size: 30px; }
    .ub-row[data-stack="true"] { flex-direction: column; }
    .ub-navbar nav { display: none; }
}
`;
