// Turns the document into a standalone HTML page. The page carries the same
// component stylesheet the canvas uses, so the export looks identical to the
// artboard (minus the builder chrome).

import {getDef} from './components';
import {type UIDocument, type UINode} from './model';
import {componentStyles} from './theme';

export function renderToHTML(doc: UIDocument): string {
    const body = renderPure(doc.root);
    body.dataset.theme = doc.theme;
    const title = escapeHTML(doc.name || 'Untitled');
    return `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>${title}</title>
<meta name="generator" content="Wails UI Builder">
<style>
html, body { margin: 0; min-height: 100%; }
body { display: flex; flex-direction: column; }
${componentStyles.trim()}
</style>
</head>
<body>
${indent(body.outerHTML)}
</body>
</html>
`;
}

/** Render a node without any of the builder's data attributes or classes. */
function renderPure(node: UINode): HTMLElement {
    const def = getDef(node.type);
    const children = (node.children ?? []).map(renderPure);
    return def.render(node, children);
}

function escapeHTML(s: string): string {
    return s.replace(/[&<>"']/g, (c) => ({'&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;'}[c] ?? c));
}

/** Light pretty-printing: one element per line, indented by depth. */
function indent(html: string): string {
    const out: string[] = [];
    let depth = 0;
    // Split into tags and text runs without regex lookbehind (older WebKit lacks it).
    const tokens: string[] = [];
    let buf = '';
    for (const ch of html.replace(/>\s*</g, '><')) {
        if (ch === '<' && buf) {
            tokens.push(buf);
            buf = '';
        }
        buf += ch;
        if (ch === '>') {
            tokens.push(buf);
            buf = '';
        }
    }
    if (buf) {
        tokens.push(buf);
    }
    for (const token of tokens) {
        const isClose = token.startsWith('</');
        const isVoid = /^<(img|hr|br|input|meta|link)\b/i.test(token) || token.endsWith('/>');
        const isOpen = /^<[a-z]/i.test(token) && !isVoid;
        if (isClose) {
            depth = Math.max(0, depth - 1);
        }
        out.push('  '.repeat(depth) + token);
        if (isOpen) {
            depth += 1;
        }
    }
    return out.join('\n');
}
