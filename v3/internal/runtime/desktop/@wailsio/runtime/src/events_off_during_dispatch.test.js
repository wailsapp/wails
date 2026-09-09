import { On, OffAll } from './events';
import { eventListeners } from './listener';
import { expect, describe, it, afterEach } from 'vitest';

const dispatchWailsEvent = window._wails.dispatchWailsEvent;

afterEach(() => {
    OffAll();
});

// Guards #4393: dispatchWailsEvent used to write its pre-dispatch listener
// snapshot back into the map after running callbacks, undoing any
// subscription change made inside a handler.
describe('unsubscribing during dispatch', () => {
    it('keeps a self-removing listener removed (single listener)', () => {
        let calls = 0;
        const off = On('evt', () => { calls++; off(); });

        dispatchWailsEvent({ name: 'evt' });
        dispatchWailsEvent({ name: 'evt' });

        expect(calls).toBe(1);
        expect(eventListeners.has('evt')).toBe(false);
    });

    it('keeps a self-removing listener removed (multiple listeners)', () => {
        let aCalls = 0, bCalls = 0;
        const offA = On('evt', () => { aCalls++; offA(); });
        On('evt', () => { bCalls++; });

        dispatchWailsEvent({ name: 'evt' });
        dispatchWailsEvent({ name: 'evt' });

        expect(aCalls).toBe(1);
        expect(bCalls).toBe(2);
    });

    it('keeps a listener added during dispatch', () => {
        let lateCalls = 0;
        const offSetup = On('evt', () => {
            On('evt', () => { lateCalls++; });
            offSetup();
        });

        dispatchWailsEvent({ name: 'evt' });
        dispatchWailsEvent({ name: 'evt' });

        expect(lateCalls).toBe(1);
    });
});
