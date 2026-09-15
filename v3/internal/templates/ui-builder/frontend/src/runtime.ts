// Run mode: wires the rendered window to the Go backend, exactly as the
// exported frontend's main.js does (see exporter.ts — keep the two in step).
//
//  - Buttons with data-call="Service.Method" call that bound Go method through
//    Call.ByName, passing the arguments listed in data-args ($field reads the
//    value of the control with that field name) and publish the result to the
//    channel named in data-result-to.
//  - Elements with data-bind="channel" display whatever is published to that
//    channel: the data of Go events with that name, or a call's result.
//  - data-wml-* attributes (events, window actions, URLs) are handled by the
//    runtime's WML support.

import {Call, Events, WML} from '@wailsio/runtime';

export type Cleanup = () => void;

export function wireRuntime(root: HTMLElement): Cleanup {
    const cleanups: Cleanup[] = [];

    // --- Bindings: one Go event subscription per distinct channel.
    const bound = new Map<string, HTMLElement[]>();
    for (const el of root.querySelectorAll<HTMLElement>('[data-bind]')) {
        const channel = el.dataset.bind!;
        bound.set(channel, [...(bound.get(channel) ?? []), el]);
    }
    const publish = (channel: string, data: unknown): void => {
        for (const el of bound.get(channel) ?? []) {
            display(el, data);
        }
    };
    for (const channel of bound.keys()) {
        cleanups.push(Events.On(channel, (ev) => publish(channel, ev.data)));
    }

    // --- Go method calls.
    for (const button of root.querySelectorAll<HTMLButtonElement>('[data-call]')) {
        const handler = async (): Promise<void> => {
            const name = qualify(button.dataset.call!);
            const args = parseArgs(button.dataset.args ?? '', root);
            const target = button.dataset.resultTo;
            button.dataset.busy = 'true';
            try {
                const result = await Call.ByName(name, ...args);
                if (target) {
                    publish(target, result);
                }
            } catch (err) {
                console.error(`${name} failed:`, err);
                if (target) {
                    publish(target, `Error: ${describe(err)}`);
                }
            } finally {
                delete button.dataset.busy;
            }
        };
        button.addEventListener('click', handler);
        cleanups.push(() => button.removeEventListener('click', handler));
    }

    // --- Declarative event / window / URL actions.
    WML.Reload();

    return () => {
        for (const fn of cleanups) {
            fn();
        }
        WML.Reload();
    };
}

/** "GreetService.Greet" → "main.GreetService.Greet" (a package-qualified name is kept as is). */
export function qualify(call: string): string {
    return call.split('.').length >= 3 ? call : `main.${call}`;
}

/** One argument per line; $field reads a control's value; numbers and booleans are typed. */
export function parseArgs(spec: string, root: ParentNode): unknown[] {
    return spec.split('\n').map((s) => s.trim()).filter(Boolean).map((arg) => {
        if (arg.startsWith('$')) {
            const control = root.querySelector<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>(`[name="${CSS.escape(arg.slice(1))}"]`);
            if (!control) {
                return '';
            }
            if (control instanceof HTMLInputElement && control.type === 'checkbox') {
                return control.checked;
            }
            if (control instanceof HTMLInputElement && (control.type === 'number' || control.type === 'range')) {
                return Number(control.value);
            }
            return control.value;
        }
        if (arg === 'true' || arg === 'false') {
            return arg === 'true';
        }
        if (arg !== '' && !Number.isNaN(Number(arg))) {
            return Number(arg);
        }
        return arg;
    });
}

/** Write a published value into a bound element according to what it is. */
export function display(el: HTMLElement, data: unknown): void {
    const text = typeof data === 'string' ? data : JSON.stringify(data);
    if (el.classList.contains('ub-progress')) {
        const pct = Math.max(0, Math.min(100, Number(data) || 0));
        el.querySelector<HTMLElement>('.ub-progress-track > i')!.style.width = `${pct}%`;
        el.querySelector('output')!.textContent = `${Math.round(pct)}%`;
        return;
    }
    if (el.dataset.bindMode === 'append') {
        el.textContent = (el.textContent ? el.textContent + '\n' : '') + text;
        el.scrollTop = el.scrollHeight;
        return;
    }
    el.textContent = text;
}

function describe(err: unknown): string {
    if (err instanceof Error) {
        return err.message || 'the call failed (is the Go backend running?)';
    }
    return String(err) || 'the call failed';
}
