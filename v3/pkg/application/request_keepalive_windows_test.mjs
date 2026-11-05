import assert from 'node:assert/strict';
import { webcrypto } from 'node:crypto';
import { readFileSync } from 'node:fs';
import vm from 'node:vm';
import test from 'node:test';

const source = readFileSync(new URL('./request_keepalive_windows.js', import.meta.url), 'utf8');
function fixture(fetch) {
    const events = [];
    const context = { Request, URL, crypto: webcrypto, fetch,
        __wailsAssetKeepalive: payload => events.push(JSON.parse(payload)) };
    vm.runInNewContext(source, context);
    return { fetch: context.fetch, events };
}

test('ordinary fetch forwards the original arguments', async () => {
    const options = { signal: new AbortController().signal };
    const f = fixture((input, init) => {
        assert.equal(input, 'http://wails.localhost/data');
        assert.equal(init, options);
        return Promise.resolve('response');
    });
    assert.equal(await f.fetch('http://wails.localhost/data', options), 'response');
    assert.deepEqual(f.events, []);
});

test('local keepalive preserves POST body, headers and request options', async () => {
    const original = new Request('http://wails.localhost/data#original', {
        method: 'POST', body: 'body', headers: { 'X-Test': 'preserved' },
        credentials: 'include', cache: 'no-store', redirect: 'error', keepalive: true,
    });
    const f = fixture(async (input, init) => {
        const request = new Request(input, init);
        assert.equal(await request.text(), 'body');
        assert.equal(request.headers.get('X-Test'), 'preserved');
        assert.equal(request.headers.has('X-Wails-Keepalive-Id'), false);
        assert.equal(request.credentials, 'include');
        assert.equal(request.cache, 'no-store');
        assert.equal(request.redirect, 'error');
        assert.equal(request.keepalive, true);
        assert.match(request.url, /#original&__wails_keepalive=[\da-f-]{36}$/);
        return 'response';
    });
    assert.equal(await f.fetch(original), 'response');
    assert.equal(original.url, 'http://wails.localhost/data#original');
    assert.deepEqual(f.events.map(e => e.State), ['start', 'end']);
    assert.equal(f.events[0].ID, f.events[1].ID);
});

test('external keepalive requests gain no marker or binding events', async () => {
    const f = fixture(request => {
        assert.equal(request.url, 'https://example.org/data');
        return Promise.resolve('response');
    });
    await f.fetch('https://example.org/data', { keepalive: true });
    assert.deepEqual(f.events, []);
});

test('explicit abort retains AbortSignal behaviour', async () => {
    const controller = new AbortController();
    const f = fixture((input, init) => new Promise((resolve, reject) => {
        init.signal.addEventListener('abort', () => reject(init.signal.reason));
    }));
    const promise = f.fetch('http://wails.localhost/data', { keepalive: true, signal: controller.signal });
    controller.abort('cancelled');
    await assert.rejects(promise, reason => reason === 'cancelled');
    assert.equal(f.events[0].State, 'start');
    assert.ok(f.events.slice(1).every(e => e.State === 'abort'));
});

test('failed delivery is distinct from explicit abort', async () => {
    const f = fixture(() => Promise.reject(new TypeError('Failed to fetch')));
    await assert.rejects(f.fetch('http://wails.localhost/data', { keepalive: true }), /Failed to fetch/);
    assert.deepEqual(f.events.map(e => e.State), ['start', 'failed']);
});

test('Request keepalive can be overridden by init', async () => {
    const f = fixture(request => {
        assert.equal(request.keepalive, false);
        assert.equal(request.url, 'http://wails.localhost/data');
        return Promise.resolve();
    });
    await f.fetch(new Request('http://wails.localhost/data', { keepalive: true }), { keepalive: false });
    assert.deepEqual(f.events, []);
});
