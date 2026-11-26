import { beforeEach, describe, expect, it, vi } from 'vitest';

describe('OpenHarmony runtime transport', () => {
    beforeEach(() => {
        vi.resetModules();
        const host = {
            invoke: () => undefined,
            request: (_session, _id, _payload) => true,
        };
        Object.assign(globalThis.window, {
            crypto: { randomUUID: () => 'session-1' },
            setTimeout: globalThis.setTimeout,
            clearTimeout: globalThis.clearTimeout,
            wailsOpenHarmony: host,
        });
    });

    it('routes runtime calls through the ArkTS proxy and resolves Go replies', async () => {
        const captured = [];
        globalThis.window.wailsOpenHarmony.request = (session, id, payload) => {
            captured.push({ session, id, payload });
            return true;
        };

        const runtime = await import('./runtime.js');
        const pending = runtime.newRuntimeCaller(0, 'main')(7, ['arg']);
        await new Promise((resolve) => queueMicrotask(resolve));

        expect(captured).toHaveLength(1);
        expect(captured[0].session).toBe('session-1');
        expect(JSON.parse(captured[0].payload)).toMatchObject({ object: 0, method: 7, windowName: 'main' });
        expect(globalThis.window.wailsOpenHarmonySession).toBe('session-1');

        globalThis.window.wailsOpenHarmonyResolve(captured[0].id, JSON.stringify({ ok: true, text: 'done' }));
        await expect(pending).resolves.toBe('done');
    });
});
