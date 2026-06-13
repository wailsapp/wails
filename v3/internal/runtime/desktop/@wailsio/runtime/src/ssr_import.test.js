// @vitest-environment node
import { describe, it, expect } from 'vitest';
describe('SSR import safety', () => {
    it('imports the full runtime without a DOM', async () => {
        expect(typeof window).toBe('undefined');
        await import('./index');
    });
});
