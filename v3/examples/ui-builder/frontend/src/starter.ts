// The design that greets a first-time user: a typical desktop application —
// toolbar, sidebar navigation, a content area with a table and metrics, a form
// that calls the Go backend, and a status bar bound to a Go event.

import {type UIDocument, type UINode, newId} from './model';

const n = (type: string, props: Record<string, unknown>, children?: UINode[]): UINode => ({id: newId(), type, props, children});

export function starterDocument(): UIDocument {
    return {
        version: 1,
        name: 'Project Manager',
        theme: 'light',
        chrome: 'mac',
        root: {
            id: 'root', type: 'root', props: {},
            children: [
                n('toolbar', {height: 44, gap: 8, drag: true, background: ''}, [
                    n('heading', {text: 'Project Manager', level: 'window', align: 'left', color: ''}),
                    n('spacer', {size: 0, fill: true}),
                    n('input', {label: '', name: 'query', placeholder: 'Search projects', type: 'search', help: '', fill: false, inline: false}),
                    n('button', {label: 'New project', variant: 'primary', size: 'sm', full: false, action: 'event', event: 'project:new', service: '', method: '', args: '', resultTo: '', url: '', windowAction: 'Minimise'}),
                ]),
                n('hstack', {gap: 0, padding: 0, items: 'stretch', justify: 'start', fill: true, scroll: false, background: ''}, [
                    n('sidebar', {width: 200, side: 'left', padding: 8, gap: 4, background: ''}, [
                        n('nav', {items: 'Overview\nProjects\nBuilds\nReleases\n# Settings\nGeneral\nAccount', active: 1, style: 'accent'}),
                        n('spacer', {size: 0, fill: true}),
                        n('badge', {text: 'Connected', tone: 'ok', bind: ''}),
                    ]),
                    n('content', {padding: 16, gap: 14, scroll: true, background: ''}, [
                        n('hstack', {gap: 12, padding: 0, items: 'center', justify: 'between', fill: false, scroll: false, background: ''}, [
                            n('heading', {text: 'Projects', level: 'title', align: 'left', color: ''}),
                            n('tabs', {items: 'All\nActive\nArchived', active: 0, style: 'segmented'}),
                        ]),
                        n('hstack', {gap: 12, padding: 0, items: 'stretch', justify: 'start', fill: false, scroll: false, background: ''}, [
                            n('metric', {label: 'Projects', value: '14', note: '3 building now', fill: true, bind: ''}),
                            n('metric', {label: 'Builds today', value: '27', note: '2 failed', fill: true, bind: ''}),
                            n('metric', {label: 'Disk used', value: '3.8 GB', note: 'of 20 GB', fill: true, bind: ''}),
                        ]),
                        n('table', {
                            columns: 'Name | Platform | Status | Last build',
                            rows: 'wails-app | darwin/arm64 | Building | 2 min ago\nrelease-notes | windows/amd64 | Ready | 1 h ago\nwebsite | linux/amd64 | Ready | yesterday\ninternal-tools | darwin/arm64 | Failed | 3 days ago',
                            striped: true, selected: 0, fill: false,
                        }),
                        n('group', {title: 'Talk to Go', gap: 10, padding: 14, fill: false, background: ''}, [
                            n('text', {text: 'This form calls GreetService.Greet on the Go side and shows what comes back. Press Run to try it.', size: 'sm', tone: 'muted', mono: false, align: 'left', color: '', bind: ''}),
                            n('hstack', {gap: 8, padding: 0, items: 'end', justify: 'start', fill: false, scroll: false, background: ''}, [
                                n('input', {label: 'Name', name: 'name', placeholder: 'Your name', type: 'text', help: '', fill: true, inline: false}),
                                n('button', {label: 'Greet', variant: 'primary', size: 'md', full: false, action: 'call', service: 'GreetService', method: 'Greet', args: '$name', resultTo: 'greeting', event: '', url: '', windowAction: 'Minimise'}),
                            ]),
                            n('text', {text: 'The reply from Go appears here.', size: 'md', tone: 'default', mono: true, align: 'left', color: '', bind: 'greeting'}),
                        ]),
                    ]),
                ]),
                n('statusbar', {gap: 14, background: ''}, [
                    n('text', {text: 'Ready', size: 'sm', tone: 'muted', mono: false, align: 'left', color: '', bind: 'status'}),
                    n('spacer', {size: 0, fill: true}),
                    n('text', {text: 'Waiting for the time event from Go…', size: 'sm', tone: 'muted', mono: false, align: 'left', color: '', bind: 'time'}),
                ]),
            ],
        },
    };
}
