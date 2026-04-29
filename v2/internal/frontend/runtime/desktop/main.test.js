import { expect, describe, it, vi, beforeEach, afterEach } from 'vitest'

function createMouseEvent(type, opts = {}) {
    return new MouseEvent(type, {
        bubbles: true,
        cancelable: true,
        clientX: opts.clientX || 0,
        clientY: opts.clientY || 0,
        buttons: opts.buttons !== undefined ? opts.buttons : 1,
        detail: opts.detail !== undefined ? opts.detail : 1,
    });
}

function createDragController(threshold = 3) {
    const state = {
        shouldDrag: false,
        dragStartX: 0,
        dragStartY: 0,
        dragThreshold: threshold,
        deferDragToMouseMove: true,
    };

    const WailsInvoke = vi.fn();

    function dragTest(e, isDraggable) {
        if (!isDraggable) return false;
        if (e.buttons !== 1) return false;
        if (e.detail !== 1) return false;
        return true;
    }

    function handleMouseDown(e, isDraggable) {
        if (dragTest(e, isDraggable)) {
            if (state.deferDragToMouseMove) {
                state.shouldDrag = true;
                state.dragStartX = e.clientX;
                state.dragStartY = e.clientY;
            } else {
                e.preventDefault();
                WailsInvoke("drag");
            }
        } else {
            state.shouldDrag = false;
        }
    }

    function handleMouseMove(e) {
        if (state.shouldDrag) {
            const dx = Math.abs(e.clientX - state.dragStartX);
            const dy = Math.abs(e.clientY - state.dragStartY);
            if (dx < state.dragThreshold && dy < state.dragThreshold) {
                return;
            }
            state.shouldDrag = false;
            const mousePressed = e.buttons !== undefined ? e.buttons : e.which;
            if (mousePressed > 0) {
                WailsInvoke("drag");
                return;
            }
        }
    }

    function handleMouseUp() {
        state.shouldDrag = false;
    }

    return { state, WailsInvoke, handleMouseDown, handleMouseMove, handleMouseUp };
}

describe('Drag threshold', () => {
    let controller;

    beforeEach(() => {
        controller = createDragController(3);
    });

    it('should not start drag on mouse movement below threshold (issue #4285)', () => {
        const downEvent = createMouseEvent('mousedown', {
            clientX: 100,
            clientY: 100,
            buttons: 1,
            detail: 1
        });

        controller.handleMouseDown(downEvent, true);
        expect(controller.state.shouldDrag).toBe(true);

        const moveEvent = createMouseEvent('mousemove', {
            clientX: 101,
            clientY: 101,
            buttons: 1
        });
        controller.handleMouseMove(moveEvent);
        expect(controller.WailsInvoke).not.toHaveBeenCalled();
        expect(controller.state.shouldDrag).toBe(true);

        controller.handleMouseUp();
        expect(controller.state.shouldDrag).toBe(false);
    });

    it('should start drag when movement exceeds threshold', () => {
        const downEvent = createMouseEvent('mousedown', {
            clientX: 100,
            clientY: 100,
            buttons: 1,
            detail: 1
        });

        controller.handleMouseDown(downEvent, true);
        expect(controller.state.shouldDrag).toBe(true);

        const moveEvent = createMouseEvent('mousemove', {
            clientX: 105,
            clientY: 105,
            buttons: 1
        });
        controller.handleMouseMove(moveEvent);
        expect(controller.WailsInvoke).toHaveBeenCalledWith("drag");
        expect(controller.state.shouldDrag).toBe(false);
    });

    it('should not invoke drag when click has no movement', () => {
        const downEvent = createMouseEvent('mousedown', {
            clientX: 100,
            clientY: 100,
            buttons: 1,
            detail: 1
        });

        controller.handleMouseDown(downEvent, true);

        const moveEvent = createMouseEvent('mousemove', {
            clientX: 100,
            clientY: 100,
            buttons: 1
        });
        controller.handleMouseMove(moveEvent);
        expect(controller.WailsInvoke).not.toHaveBeenCalled();

        controller.handleMouseUp();
        expect(controller.state.shouldDrag).toBe(false);
    });

    it('should not set shouldDrag on non-draggable elements', () => {
        const downEvent = createMouseEvent('mousedown', {
            clientX: 100,
            clientY: 100,
            buttons: 1,
            detail: 1
        });

        controller.handleMouseDown(downEvent, false);
        expect(controller.state.shouldDrag).toBe(false);
    });

    it('should handle diagonal movement within threshold', () => {
        const downEvent = createMouseEvent('mousedown', {
            clientX: 100,
            clientY: 100,
            buttons: 1,
            detail: 1
        });

        controller.handleMouseDown(downEvent, true);

        const moveEvent = createMouseEvent('mousemove', {
            clientX: 102,
            clientY: 102,
            buttons: 1
        });
        controller.handleMouseMove(moveEvent);
        expect(controller.WailsInvoke).not.toHaveBeenCalled();
        expect(controller.state.shouldDrag).toBe(true);

        const moveEvent2 = createMouseEvent('mousemove', {
            clientX: 104,
            clientY: 103,
            buttons: 1
        });
        controller.handleMouseMove(moveEvent2);
        expect(controller.WailsInvoke).toHaveBeenCalledWith("drag");
    });
});
