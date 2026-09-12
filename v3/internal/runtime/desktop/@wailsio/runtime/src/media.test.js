import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { SetSource, ClearSource } from './media';

const sockets = vi.hoisted(() => []);
vi.mock('./stream.js', () => ({
    Stream: vi.fn((name) => {
        const socket = {
            name,
            send: vi.fn(),
            close: vi.fn(),
            binaryType: '',
            onopen: null,
            onmessage: null,
            onclose: null,
            onerror: null,
        };
        sockets.push(socket);
        return socket;
    }),
}));
let player, created, revoked;
function message(socket, bytes) {
    socket.onmessage?.({ data: bytes });
}
function meta(socket, data = {}) {
    message(
        socket,
        new TextEncoder().encode(
            JSON.stringify({ version: 1, size: 3, type: 'audio/wav', ...data }),
        ).buffer,
    );
}
function complete(socket = sockets.at(-1)) {
    meta(socket);
    message(socket, new Uint8Array([1, 2, 3]).buffer);
    socket.onclose?.({ code: 1000 });
}
function load(name = 'a.wav', options) {
    return SetSource(player, 'media', name, options);
}
beforeEach(() => {
    sockets.length = 0;
    player = document.createElement('audio');
    document.body.append(player);
    vi.spyOn(player, 'load').mockImplementation(() => {});
    vi.spyOn(player, 'pause').mockImplementation(() => {});
    created = vi.fn(() => 'blob:test-' + created.mock.calls.length);
    revoked = vi.fn();
    vi.spyOn(URL, 'createObjectURL').mockImplementation(created);
    vi.spyOn(URL, 'revokeObjectURL').mockImplementation(revoked);
});
afterEach(() => {
    ClearSource(player);
    player.remove();
    vi.restoreAllMocks();
});
it('requests a bounded file over its named stream and assigns a typed blob', async () => {
    const pending = load();
    const socket = sockets[0];
    socket.onopen();
    expect(socket.name).toBe('media');
    expect(JSON.parse(socket.send.mock.calls[0][0])).toEqual({
        version: 1,
        name: 'a.wav',
        maxBytes: 32 * 1024 * 1024,
    });
    complete();
    await pending;
    expect(created.mock.calls[0][0].size).toBe(3);
    expect(created.mock.calls[0][0].type).toBe('audio/wav');
    expect(player.src).toBe('blob:test-1');
    expect(socket.close).toHaveBeenCalledTimes(1);
});
it('releases replaced blobs and clears idempotently', async () => {
    let pending = load();
    complete();
    await pending;
    pending = load('b.wav');
    complete();
    await pending;
    expect(revoked.mock.calls).toEqual([['blob:test-1']]);
    ClearSource(player);
    ClearSource(player);
    expect(revoked.mock.calls).toEqual([['blob:test-1'], ['blob:test-2']]);
    expect(player.hasAttribute('src')).toBe(false);
});
it('never lets a stale transfer overwrite a newer selection', async () => {
    const first = load().catch((e) => e);
    const old = sockets[0];
    const lateMessage = old.onmessage;
    const lateClose = old.onclose;
    const second = load('b.wav');
    complete();
    await second;
    lateMessage({ data: new Uint8Array([1]).buffer });
    lateClose({ code: 1000 });
    expect((await first).name).toBe('AbortError');
    expect(old.close).toHaveBeenCalledTimes(1);
    expect(player.src).toBe('blob:test-1');
    expect(created).toHaveBeenCalledTimes(1);
});
it('cancels a pending connection on clear', async () => {
    const pending = load().catch((e) => e);
    ClearSource(player);
    expect((await pending).name).toBe('AbortError');
    expect(sockets[0].close).toHaveBeenCalledTimes(1);
    expect(created).not.toHaveBeenCalled();
});
it('honours caller cancellation while preserving the previous source', async () => {
    let pending = load();
    complete();
    await pending;
    const c = new AbortController();
    pending = load('b.wav', { signal: c.signal }).catch((e) => e);
    c.abort();
    expect((await pending).name).toBe('AbortError');
    expect(player.src).toBe('blob:test-1');
    expect(revoked).not.toHaveBeenCalled();
});
it.each([{ size: 9 }, { error: 'too large', limit: true }])(
    'rejects server byte limits %s',
    async (data) => {
        const pending = load('a.wav', { maxBytes: 8 });
        meta(sockets[0], data);
        await expect(pending).rejects.toBeInstanceOf(RangeError);
        expect(sockets[0].close).toHaveBeenCalledTimes(1);
        expect(created).not.toHaveBeenCalled();
    },
);
it('checks actual bytes against declared size', async () => {
    const pending = load('a.wav', { maxBytes: 8 });
    meta(sockets[0], { size: 3 });
    message(sockets[0], new Uint8Array(4).buffer);
    await expect(pending).rejects.toBeInstanceOf(RangeError);
});
it('rejects oversized chunks', async () => {
    const pending = load();
    meta(sockets[0], { size: 128 * 1024 });
    message(sockets[0], new Uint8Array(64 * 1024 + 1).buffer);
    await expect(pending).rejects.toBeInstanceOf(RangeError);
});
it('accepts exactly the limit', async () => {
    const pending = load('a.wav', { maxBytes: 3 });
    complete();
    await pending;
    expect(created.mock.calls[0][0].size).toBe(3);
});
it.each([0, -1, Infinity, NaN, 1.5])(
    'rejects invalid limit %s before opening a stream',
    async (maxBytes) => {
        await expect(load('a.wav', { maxBytes })).rejects.toBeInstanceOf(
            RangeError,
        );
        expect(sockets).toHaveLength(0);
    },
);
it('rejects already aborted signals', async () => {
    const c = new AbortController();
    c.abort();
    await expect(load('a.wav', { signal: c.signal })).rejects.toHaveProperty(
        'name',
        'AbortError',
    );
    expect(sockets).toHaveLength(0);
});
it('preserves the previous source on server errors and permits retry', async () => {
    let pending = load();
    complete();
    await pending;
    pending = load('b.wav');
    meta(sockets.at(-1), { error: 'not found' });
    await expect(pending).rejects.toThrow('not found');
    expect(player.src).toBe('blob:test-1');
    pending = load('b.wav');
    complete();
    await pending;
    expect(player.src).toBe('blob:test-2');
});
it('rejects truncated files', async () => {
    const pending = load();
    meta(sockets[0]);
    sockets[0].onclose({ code: 1000 });
    await expect(pending).rejects.toThrow('complete file');
});
it('rejects transport failures', async () => {
    const pending = load();
    sockets[0].onerror();
    await expect(pending).rejects.toThrow('stream failed');
});
it('supports detached Audio elements', async () => {
    player.remove();
    const pending = load();
    complete();
    await pending;
    expect(player.src).toBe('blob:test-1');
});
it('shares ownership between runtime module instances', async () => {
    let pending = load();
    complete();
    await pending;
    vi.resetModules();
    const other = await import('./media');
    pending = other.SetSource(player, 'media', 'b.wav');
    complete();
    await pending;
    expect(revoked).toHaveBeenCalledWith('blob:test-1');
    ClearSource(player);
    expect(revoked).toHaveBeenCalledWith('blob:test-2');
});

it('rejects empty data frames instead of retaining an unbounded list', async () => {
    const pending = load();
    meta(sockets[0]);
    message(sockets[0], new ArrayBuffer(0));
    await expect(pending).rejects.toBeInstanceOf(RangeError);
    expect(sockets[0].close).toHaveBeenCalledTimes(1);
});
