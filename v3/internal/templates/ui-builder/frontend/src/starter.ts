// The design that greets a first-time user: a small landing page that shows
// off nesting, rows, cards and the form controls. It's just a document — edit
// or delete anything.

import {type UIDocument, type UINode, newId} from './model';

const n = (type: string, props: Record<string, unknown>, children?: UINode[]): UINode => ({id: newId(), type, props, children});

export function starterDocument(): UIDocument {
    return {
        version: 1,
        name: 'Launch page',
        theme: 'light',
        root: {
            id: 'root', type: 'root', props: {},
            children: [
                n('navbar', {brand: 'Nimbus', links: 'Features\nPricing\nDocs\nChangelog', cta: 'Download'}),
                n('section', {padding: 88, gap: 20, maxWidth: 820, background: 'transparent', align: 'center', items: 'center'}, [
                    n('badge', {text: 'Now with Wails v3', tone: 'accent'}),
                    n('heading', {text: 'Ship native desktop apps with the web stack you love', level: 'h1', align: 'center', gradient: false, color: ''}),
                    n('text', {text: 'Nimbus pairs a Go backend with a modern frontend, wraps it in a real native window, and builds a single binary for macOS, Windows and Linux.', size: 'lg', muted: true, align: 'center', color: ''}),
                    n('row', {gap: 12, justify: 'center', items: 'center', equal: false, stack: true, padding: 0}, [
                        n('button', {label: 'Get started', variant: 'primary', size: 'lg', full: false, href: '#'}),
                        n('button', {label: 'Read the docs', variant: 'secondary', size: 'lg', full: false, href: '#'}),
                    ]),
                ]),
                n('section', {padding: 24, gap: 16, maxWidth: 1040, background: 'transparent', align: 'left', items: 'stretch'}, [
                    n('row', {gap: 16, justify: 'start', items: 'stretch', equal: true, stack: true, padding: 0}, [
                        n('card', {padding: 24, gap: 10, radius: 16, elevated: true, background: '', align: 'left'}, [
                            n('metric', {label: 'Binary size', value: '9.4 MB', delta: '−38% vs Electron', trend: 'up', align: 'left'}),
                        ]),
                        n('card', {padding: 24, gap: 10, radius: 16, elevated: true, background: '', align: 'left'}, [
                            n('metric', {label: 'Cold start', value: '120 ms', delta: '+ instant window', trend: 'up', align: 'left'}),
                        ]),
                        n('card', {padding: 24, gap: 10, radius: 16, elevated: true, background: '', align: 'left'}, [
                            n('metric', {label: 'Platforms', value: '3', delta: 'one codebase', trend: 'flat', align: 'left'}),
                        ]),
                    ]),
                ]),
                n('section', {padding: 64, gap: 24, maxWidth: 1040, background: '', align: 'left', items: 'stretch'}, [
                    n('row', {gap: 48, justify: 'start', items: 'center', equal: true, stack: true, padding: 0}, [
                        n('section', {padding: 0, gap: 14, maxWidth: 520, background: 'transparent', align: 'left', items: 'start'}, [
                            n('heading', {text: 'Everything you need, nothing you don’t', level: 'h2', align: 'left', gradient: true, color: ''}),
                            n('text', {text: 'Bind Go methods, call them from TypeScript with full type safety, and let the runtime handle windows, menus, dialogs and events.', size: 'md', muted: true, align: 'left', color: ''}),
                            n('list', {items: 'Type-safe Go ↔ TS bindings\nNative menus, dialogs and tray\nHot reload in development', style: 'check'}),
                        ]),
                        n('image', {src: '', alt: 'Product screenshot', height: 300, radius: 18, fit: 'cover'}),
                    ]),
                ]),
                n('section', {padding: 64, gap: 16, maxWidth: 520, background: 'transparent', align: 'center', items: 'stretch'}, [
                    n('card', {padding: 28, gap: 14, radius: 18, elevated: true, background: '', align: 'left'}, [
                        n('heading', {text: 'Join the beta', level: 'h3', align: 'left', gradient: false, color: ''}),
                        n('input', {label: 'Email address', placeholder: 'you@example.com', type: 'email', help: 'We send one email a month. No spam.'}),
                        n('toggle', {label: 'Notify me about releases', on: true}),
                        n('button', {label: 'Request access', variant: 'primary', size: 'md', full: true, href: '#'}),
                    ]),
                ]),
                n('divider', {}),
                n('section', {padding: 28, gap: 8, maxWidth: 1040, background: 'transparent', align: 'center', items: 'center'}, [
                    n('text', {text: '© 2026 Nimbus Labs · Built with Wails', size: 'sm', muted: true, align: 'center', color: ''}),
                ]),
            ],
        },
    };
}
