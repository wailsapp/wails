// Coalesce resize events while keeping native SetBounds requests in order.
export function createLayoutScheduler(layout, report, requestFrame = requestAnimationFrame) {
    let pending = false;
    let scheduled = false;
    let running = false;

    async function flush() {
        scheduled = false;
        running = true;
        do {
            pending = false;
            try { await layout(); } catch (error) { report(error); }
        } while (pending);
        running = false;
    }

    return function schedule() {
        pending = true;
        if (!running && !scheduled) {
            scheduled = true;
            requestFrame(flush);
        }
    };
}
