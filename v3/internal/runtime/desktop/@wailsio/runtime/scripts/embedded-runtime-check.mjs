import assert from 'node:assert/strict';
import { readFile } from 'node:fs/promises';
import { test } from 'node:test';
import { SourceMap } from 'node:module';
import { createContext, SourceTextModule } from 'node:vm';

const assets = new URL('../../../../../assetserver/bundledassets/', import.meta.url);
const exports = [
    'Android', 'Application', 'Browser', 'Call', 'CancelError', 'CancellablePromise',
    'CancelledRejectionError', 'Clipboard', 'Create', 'Dialogs', 'Events', 'Flags',
    'IOS', 'JSONStream', 'Screens', 'Stream', 'System', 'Updater', 'WML', 'WailsSocket',
    'Window', 'clientId', 'getTransport', 'loadOptionalScript', 'objectNames', 'setTransport',
].sort();

async function loadRuntime(filename, readyState, t, platform) {
    const timers = [];
    t.after(() => timers.forEach(clearInterval));
    const queries = [];
    const messages = [];
    const logs = [];
    const requests = [];
    const nativeCalls = [];
    const document = Object.assign(new EventTarget(), {
        readyState,
        documentElement: new EventTarget(),
        body: { querySelectorAll: selector => { queries.push(selector); return []; } },
    });
    const window = Object.assign(new EventTarget(), {
        _wails: { environment: { OS: platform === 'android' ? 'android' : 'linux' } },
        location: new URL('http://localhost'),
        setInterval: (...args) => {
            const timer = setInterval(...args);
            timers.push(timer);
            return timer;
        },
        clearInterval,
        nativeBridge: { invokeAsync: (id, payload) => {
            nativeCalls.push(JSON.parse(payload));
            window._wailsAndroidCallback(id, JSON.stringify({ ok: true, data: null }), null);
        }, invoke(message) {
            messages.push(message);
            if (message === 'wails:runtime:ready') {
                assert.equal(typeof window.wails?.Call.ByID, 'function');
                assert.equal(typeof window._wails.handlePlatformFileDrop, 'function');
                assert.equal(typeof window._wails.invoke, 'function');
            }
        } },
    });
    if (platform === 'webview2') window.chrome = { webview: { postMessage: window.nativeBridge.invoke } };
    else if (platform === 'webkit') window.webkit = { messageHandlers: { external: { postMessage: window.nativeBridge.invoke } } };
    else window.wails = window.nativeBridge;
    const context = createContext({
        window, document, navigator: { userAgent: 'Wails test' },
        EventTarget, AbortController, AbortSignal, CustomEvent, URL, URLSearchParams,
        MouseEvent: class { buttons = 0; },
        console: Object.fromEntries(['log', 'warn', 'error', 'debug', 'info'].map(
            name => [name, (...args) => logs.push([name, ...args])],
        )),
        fetch: async (...args) => { requests.push(args); return { ok: false }; },
        setTimeout, clearTimeout, TextEncoder, TextDecoder,
    });
    const code = await readFile(new URL(filename, assets), 'utf8');
    const module = new SourceTextModule(code, { context });
    await module.link(specifier => assert.fail(`Unexpected external import: ${specifier}`));
    await module.evaluate();
    await Promise.resolve();
    return { module, window, document, queries, messages, logs, requests, nativeCalls, code };
}

for (const debug of [false, true]) {
    const filename = debug ? 'runtime.debug.js' : 'runtime.js';
    for (const platform of ['webview2', 'webkit', 'android']) {
        for (const readyState of ['loading', 'complete']) {
            test(`${filename}: ${platform}, document ${readyState}`, async t => {
                const runtime = await loadRuntime(filename, readyState, t, platform);
                assert.deepEqual(Object.keys(runtime.module.namespace).sort(), exports);
                assert.deepEqual(Object.keys(runtime.window.wails).sort(), exports);
                for (const key of exports) {
                    assert.equal(runtime.module.namespace[key], runtime.window.wails[key]);
                }
                assert.deepEqual(runtime.messages, ['wails:runtime:ready']);
                runtime.window.wails.System.invoke({ kind: 'test' });
                assert.deepEqual(runtime.messages[1], platform === 'android' ? '{"kind":"test"}' : { kind: 'test' });
                if (platform === 'android') {
                    await runtime.window.wails.Android.Toast.Show('hello');
                    const call = runtime.nativeCalls[0];
                    assert.equal(call.object, runtime.window.wails.objectNames.Android);
                    assert.equal(call.method, 2);
                    assert.deepEqual(call.args, { message: 'hello' });
                    assert.equal(call.clientId, runtime.window.wails.clientId);
                }
                assert.equal(runtime.requests[0][0], '/wails/custom.js');
                assert.deepEqual(runtime.requests[0][1].method, 'HEAD');
                if (readyState === 'loading') {
                    assert.equal(runtime.queries.length, 0);
                    runtime.document.dispatchEvent(new Event('DOMContentLoaded'));
                }
                assert.equal(runtime.queries.length, 1);
                assert.match(runtime.queries[0], /\[data-wml-openurl\]/);
                const loadedLog = runtime.logs.some(entry => entry.some(
                    value => typeof value === 'string' && value.includes('Wails Runtime Loaded'),
                ));
                assert.equal(loadedLog, debug);
                if (!debug) assert.equal(runtime.logs.length, 0);
                const promise = new runtime.module.namespace.CancellablePromise(resolve => resolve(42));
                assert.equal(await promise, 42);
                if (debug) {
                    const match = runtime.code.match(/sourceMappingURL=data:application\/json[^,]*;base64,([^\s]+)/);
                    assert.ok(match, 'debug bundle must contain an inline source map');
                    const map = JSON.parse(Buffer.from(match[1], 'base64'));
                    assert.ok(map.sources.some(name => name.endsWith('/src/index.ts')));
                    assert.ok(map.sourcesContent.some(source => source.includes('export class CancellablePromise')));
                    const source = await readFile(new URL('../src/index.ts', import.meta.url), 'utf8');
                    assert.ok(map.sourcesContent.includes(source));
                    const generatedLine = runtime.code.split('\n').findIndex(line => line.startsWith('function loadOptionalScript('));
                    const originalLine = source.split('\n').findIndex(line => line.startsWith('export function loadOptionalScript('));
                    const position = new SourceMap(map).findEntry(generatedLine, 0);
                    assert.equal(position.originalLine, originalLine);
                    assert.ok(position.originalSource.endsWith('/src/index.ts'));
                }
            });
        }
    }
}
