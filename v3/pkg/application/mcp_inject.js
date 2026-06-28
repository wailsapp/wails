/*
 * Wails MCP in-page support library.
 *
 * Injected into the webview before every MCP tool evaluation. Provides an
 * animated, visible mouse cursor, realistic input simulation and the callback
 * harness used to return results to the Go side.
 *
 * Everything lives under window.__wailsMCP and installation is idempotent.
 */
(function () {
    'use strict';
    if (window.__wailsMCP) return;

    const MAX_Z = 2147483647;

    // ------------------------------------------------------------------ state

    const state = {
        x: Math.round(window.innerWidth / 2),
        y: Math.round(window.innerHeight / 2),
        cursor: null,
        showCursor: true,
        hovered: null,
        buttons: 0,
    };

    // ----------------------------------------------------------------- cursor

    const CURSOR_SVG =
        '<svg width="26" height="30" viewBox="0 0 26 30" xmlns="http://www.w3.org/2000/svg">' +
        '<path d="M2 1 L2 24 L8 19 L12 28 L16 26 L12 17 L20 17 Z" ' +
        'fill="white" stroke="black" stroke-width="1.6" stroke-linejoin="round"/></svg>';

    function ensureCursor() {
        if (!state.showCursor) return null;
        let cursor = document.getElementById('__wails-mcp-cursor');
        if (cursor) {
            state.cursor = cursor;
            return cursor;
        }
        cursor = document.createElement('div');
        cursor.id = '__wails-mcp-cursor';
        cursor.setAttribute('aria-hidden', 'true');
        cursor.innerHTML = CURSOR_SVG;
        Object.assign(cursor.style, {
            position: 'fixed',
            left: '0px',
            top: '0px',
            width: '26px',
            height: '30px',
            zIndex: String(MAX_Z),
            pointerEvents: 'none',
            transform: `translate(${state.x}px, ${state.y}px)`,
            transformOrigin: '2px 1px',
            filter: 'drop-shadow(1px 2px 3px rgba(0,0,0,0.45))',
            transition: 'opacity 0.2s ease',
            opacity: '0',
        });
        (document.body || document.documentElement).appendChild(cursor);
        // Fade in on first use.
        requestAnimationFrame(() => { cursor.style.opacity = '1'; });
        state.cursor = cursor;
        return cursor;
    }

    function placeCursor(x, y) {
        state.x = x;
        state.y = y;
        const cursor = ensureCursor();
        if (cursor) cursor.style.transform = `translate(${x}px, ${y}px)`;
    }

    function pressCursor(down) {
        const cursor = ensureCursor();
        if (!cursor) return;
        const svg = cursor.firstElementChild;
        if (svg) svg.style.transform = down ? 'scale(0.82)' : 'scale(1)';
    }

    function ripple(x, y, colour) {
        if (!state.showCursor) return;
        const dot = document.createElement('div');
        Object.assign(dot.style, {
            position: 'fixed',
            left: (x - 14) + 'px',
            top: (y - 14) + 'px',
            width: '28px',
            height: '28px',
            borderRadius: '50%',
            border: `2.5px solid ${colour || 'rgba(59,130,246,0.95)'}`,
            background: 'rgba(59,130,246,0.18)',
            zIndex: String(MAX_Z - 1),
            pointerEvents: 'none',
        });
        (document.body || document.documentElement).appendChild(dot);
        const animation = dot.animate(
            [
                { transform: 'scale(0.35)', opacity: 1 },
                { transform: 'scale(1.6)', opacity: 0 },
            ],
            { duration: 420, easing: 'ease-out' },
        );
        animation.onfinish = () => dot.remove();
    }

    // ----------------------------------------------------------------- events

    const MODIFIER_FLAGS = ['ctrl', 'shift', 'alt', 'meta'];

    function modifierInit(modifiers) {
        const list = (modifiers || []).map((m) => String(m).toLowerCase());
        return {
            ctrlKey: list.includes('ctrl') || list.includes('control'),
            shiftKey: list.includes('shift'),
            altKey: list.includes('alt') || list.includes('option'),
            metaKey: list.includes('meta') || list.includes('cmd') || list.includes('command'),
        };
    }

    function elementAt(x, y) {
        return document.elementFromPoint(x, y) || document.documentElement;
    }

    function mouseInit(x, y, extra) {
        return Object.assign({
            bubbles: true,
            cancelable: true,
            composed: true,
            view: window,
            clientX: x,
            clientY: y,
            screenX: x + (window.screenX || 0),
            screenY: y + (window.screenY || 0),
            buttons: state.buttons,
        }, extra || {});
    }

    function firePointer(type, el, x, y, extra) {
        const init = mouseInit(x, y, Object.assign({
            pointerId: 1,
            pointerType: 'mouse',
            isPrimary: true,
            width: 1,
            height: 1,
            pressure: state.buttons ? 0.5 : 0,
        }, extra || {}));
        try {
            el.dispatchEvent(new PointerEvent(type, init));
        } catch (e) { /* PointerEvent unavailable: mouse events still fire */ }
    }

    function fireMouse(type, el, x, y, extra) {
        return el.dispatchEvent(new MouseEvent(type, mouseInit(x, y, extra)));
    }

    // Fire hover transitions (out/leave then over/enter) when the element
    // under the cursor changes.
    function updateHover(x, y, mods) {
        const el = elementAt(x, y);
        const previous = state.hovered;
        if (el === previous) return el;
        if (previous && previous.isConnected) {
            firePointer('pointerout', previous, x, y, mods);
            fireMouse('mouseout', previous, x, y, mods);
            const leaveInit = mouseInit(x, y, Object.assign({ bubbles: false }, mods));
            let node = previous;
            while (node && node !== el && !(el && node.contains(el))) {
                node.dispatchEvent(new MouseEvent('mouseleave', leaveInit));
                node = node.parentElement;
            }
        }
        if (el) {
            firePointer('pointerover', el, x, y, mods);
            fireMouse('mouseover', el, x, y, mods);
            const enterInit = mouseInit(x, y, Object.assign({ bubbles: false }, mods));
            const chain = [];
            let node = el;
            while (node && node !== previous && !(previous && node.contains(previous))) {
                chain.unshift(node);
                node = node.parentElement;
            }
            for (const link of chain) {
                link.dispatchEvent(new MouseEvent('mouseenter', enterInit));
            }
        }
        state.hovered = el;
        return el;
    }

    function moveEvents(x, y, mods) {
        const el = updateHover(x, y, mods);
        firePointer('pointermove', el, x, y, mods);
        fireMouse('mousemove', el, x, y, mods);
        return el;
    }

    // -------------------------------------------------------------- animation

    function easeInOutCubic(t) {
        return t < 0.5 ? 4 * t * t * t : 1 - Math.pow(-2 * t + 2, 3) / 2;
    }

    function clampToViewport(x, y) {
        return {
            x: Math.max(0, Math.min(window.innerWidth - 1, x)),
            y: Math.max(0, Math.min(window.innerHeight - 1, y)),
        };
    }

    function sleep(ms) {
        return new Promise((resolve) => setTimeout(resolve, ms));
    }

    // Wait for the next animation frame, falling back to a timer so that
    // animation continues even when the webview is not rendering (unfocused
    // or occluded windows may suspend requestAnimationFrame).
    function nextFrame() {
        return new Promise((resolve) => {
            let done = false;
            const finish = () => {
                if (!done) {
                    done = true;
                    resolve();
                }
            };
            requestAnimationFrame(finish);
            setTimeout(finish, 32);
        });
    }

    // Animate the cursor from its current position to (x, y), firing move and
    // hover events along the path. Resolves when the cursor arrives.
    async function animateMove(x, y, options) {
        const opts = options || {};
        const clamped = clampToViewport(x, y);
        x = clamped.x;
        y = clamped.y;
        ensureCursor();
        const startX = state.x;
        const startY = state.y;
        const distance = Math.hypot(x - startX, y - startY);
        if (distance < 1) {
            placeCursor(x, y);
            moveEvents(x, y, opts.modifiers);
            return;
        }
        const duration = opts.duration > 0 ? opts.duration : Math.min(900, Math.max(180, distance * 1.4));
        const start = performance.now();
        for (;;) {
            await nextFrame();
            const progress = Math.min(1, (performance.now() - start) / duration);
            const eased = easeInOutCubic(progress);
            const cx = Math.round(startX + (x - startX) * eased);
            const cy = Math.round(startY + (y - startY) * eased);
            placeCursor(cx, cy);
            moveEvents(cx, cy, opts.modifiers);
            if (progress >= 1) return;
        }
    }

    // ---------------------------------------------------------------- targets

    function describe(el) {
        if (!el) return null;
        const rect = el.getBoundingClientRect();
        const description = {
            tag: el.tagName ? el.tagName.toLowerCase() : String(el),
            id: el.id || undefined,
            classes: el.classList && el.classList.length ? Array.from(el.classList) : undefined,
            text: (el.innerText || el.textContent || '').trim().slice(0, 160) || undefined,
            bounds: {
                x: Math.round(rect.x),
                y: Math.round(rect.y),
                width: Math.round(rect.width),
                height: Math.round(rect.height),
            },
        };
        if ('value' in el && typeof el.value === 'string') description.value = el.value.slice(0, 500);
        if (el.disabled) description.disabled = true;
        if (el.href) description.href = el.href;
        return description;
    }

    // Resolve {selector} or {x, y} into concrete viewport coordinates,
    // scrolling selector targets into view first.
    async function resolveTarget(target) {
        if (target && target.selector) {
            const el = document.querySelector(target.selector);
            if (!el) throw new Error('no element matches selector: ' + target.selector);
            el.scrollIntoView({ block: 'center', inline: 'center', behavior: 'instant' });
            await nextFrame();
            const rect = el.getBoundingClientRect();
            if (rect.width === 0 && rect.height === 0) {
                throw new Error('element has zero size (is it hidden?): ' + target.selector);
            }
            const point = clampToViewport(rect.x + rect.width / 2, rect.y + rect.height / 2);
            return { x: Math.round(point.x), y: Math.round(point.y), el };
        }
        if (target && typeof target.x === 'number' && typeof target.y === 'number') {
            const point = clampToViewport(target.x, target.y);
            return { x: Math.round(point.x), y: Math.round(point.y), el: null };
        }
        throw new Error('target requires a selector or x/y coordinates');
    }

    const BUTTONS = { left: 0, middle: 1, right: 2 };
    const BUTTON_MASKS = { 0: 1, 1: 4, 2: 2 };

    function focusTarget(el) {
        const focusable = el && el.closest
            ? el.closest('input, textarea, select, button, a[href], [tabindex], [contenteditable]')
            : null;
        if (focusable && typeof focusable.focus === 'function') {
            focusable.focus();
            return focusable;
        }
        if (document.activeElement && document.activeElement !== document.body) {
            document.activeElement.blur();
        }
        return null;
    }

    // ------------------------------------------------------------ public API

    async function move(target, options) {
        const point = await resolveTarget(target);
        await animateMove(point.x, point.y, options);
        return describe(elementAt(point.x, point.y));
    }

    async function click(target, options) {
        const opts = options || {};
        const mods = modifierInit(opts.modifiers);
        const button = BUTTONS[(opts.button || 'left').toLowerCase()] || 0;
        const count = Math.max(1, Math.min(3, opts.count || 1));
        const point = await resolveTarget(target);
        await animateMove(point.x, point.y, { modifiers: opts.modifiers });

        const x = point.x;
        const y = point.y;
        const colour = button === 2 ? 'rgba(245,158,11,0.95)' : 'rgba(59,130,246,0.95)';

        let el = elementAt(x, y);
        for (let i = 1; i <= count; i++) {
            el = elementAt(x, y);
            const init = Object.assign({ button, detail: i }, mods);
            state.buttons = BUTTON_MASKS[button];
            pressCursor(true);
            ripple(x, y, colour);
            firePointer('pointerdown', el, x, y, init);
            const defaultNotPrevented = fireMouse('mousedown', el, x, y, init);
            if (defaultNotPrevented && button === 0) focusTarget(el);
            await sleep(70);
            state.buttons = 0;
            pressCursor(false);
            firePointer('pointerup', el, x, y, init);
            fireMouse('mouseup', el, x, y, init);
            if (button === 0) {
                fireMouse('click', el, x, y, init);
            } else if (button === 1) {
                fireMouse('auxclick', el, x, y, init);
            } else if (button === 2) {
                fireMouse('contextmenu', el, x, y, init);
            }
            if (i < count) await sleep(90);
        }
        if (count === 2 && button === 0) {
            fireMouse('dblclick', el, x, y, Object.assign({ button, detail: 2 }, mods));
        }
        await sleep(50);
        return describe(el);
    }

    async function drag(from, to, options) {
        const opts = options || {};
        const start = await resolveTarget(from);
        await animateMove(start.x, start.y);

        const source = elementAt(start.x, start.y);
        const html5 = !!(source && source.closest && source.closest('[draggable="true"]'));
        const dragSource = html5 ? source.closest('[draggable="true"]') : source;

        state.buttons = 1;
        pressCursor(true);
        ripple(start.x, start.y, 'rgba(16,185,129,0.95)');
        firePointer('pointerdown', dragSource, start.x, start.y, { button: 0, detail: 1 });
        fireMouse('mousedown', dragSource, start.x, start.y, { button: 0, detail: 1 });

        let dataTransfer = null;
        if (html5) {
            dataTransfer = new DataTransfer();
            dragSource.dispatchEvent(new DragEvent('dragstart', mouseInit(start.x, start.y, { dataTransfer })));
        }

        try {
            // Resolve the destination after pressing, in case the press changed layout.
            const end = await resolveTarget(to);
            const distance = Math.hypot(end.x - start.x, end.y - start.y);
            const duration = opts.duration > 0 ? opts.duration : Math.min(1200, Math.max(300, distance * 1.8));

            const begin = performance.now();
            for (;;) {
                await nextFrame();
                const progress = Math.min(1, (performance.now() - begin) / duration);
                const eased = easeInOutCubic(progress);
                const cx = Math.round(start.x + (end.x - start.x) * eased);
                const cy = Math.round(start.y + (end.y - start.y) * eased);
                placeCursor(cx, cy);
                const over = moveEvents(cx, cy);
                if (html5) {
                    over.dispatchEvent(new DragEvent('dragover', mouseInit(cx, cy, { dataTransfer })));
                }
                if (progress >= 1) break;
            }

            const dropTarget = elementAt(end.x, end.y);
            if (html5) {
                dropTarget.dispatchEvent(new DragEvent('drop', mouseInit(end.x, end.y, { dataTransfer })));
                dragSource.dispatchEvent(new DragEvent('dragend', mouseInit(end.x, end.y, { dataTransfer })));
            }
            state.buttons = 0;
            pressCursor(false);
            firePointer('pointerup', dropTarget, end.x, end.y, { button: 0, detail: 1 });
            fireMouse('mouseup', dropTarget, end.x, end.y, { button: 0, detail: 1 });
            ripple(end.x, end.y, 'rgba(16,185,129,0.95)');
            await sleep(50);
            return { from: describe(dragSource), to: describe(dropTarget) };
        } finally {
            if (state.buttons !== 0) {
                state.buttons = 0;
                pressCursor(false);
            }
        }
    }

    function findScrollable(el, deltaX, deltaY) {
        let node = el;
        while (node && node !== document.documentElement) {
            const style = getComputedStyle(node);
            const canScrollY = deltaY !== 0 && node.scrollHeight > node.clientHeight &&
                ['auto', 'scroll', 'overlay'].includes(style.overflowY);
            const canScrollX = deltaX !== 0 && node.scrollWidth > node.clientWidth &&
                ['auto', 'scroll', 'overlay'].includes(style.overflowX);
            if (canScrollY || canScrollX) return node;
            node = node.parentElement;
        }
        return document.scrollingElement || document.documentElement;
    }

    async function scroll(target, deltaX, deltaY) {
        const point = await resolveTarget(target);
        await animateMove(point.x, point.y);
        const el = elementAt(point.x, point.y);
        const allowed = el.dispatchEvent(new WheelEvent('wheel', mouseInit(point.x, point.y, {
            deltaX: deltaX,
            deltaY: deltaY,
            deltaMode: 0,
        })));
        let scrolled = null;
        if (allowed) {
            scrolled = findScrollable(el, deltaX, deltaY);
            const beforeTop = scrolled.scrollTop;
            const beforeLeft = scrolled.scrollLeft;
            scrolled.scrollBy({ left: deltaX, top: deltaY, behavior: 'smooth' });
            await sleep(400);
            // Smooth scrolling needs rendering frames; fall back to instant if
            // the webview is not rendering (e.g. unfocused window).
            if (scrolled.scrollTop === beforeTop && scrolled.scrollLeft === beforeLeft) {
                scrolled.scrollBy({ left: deltaX, top: deltaY, behavior: 'instant' });
            }
        }
        return {
            target: describe(el),
            scrolled: scrolled ? describe(scrolled) : null,
            scrollTop: scrolled ? Math.round(scrolled.scrollTop) : null,
            scrollLeft: scrolled ? Math.round(scrolled.scrollLeft) : null,
        };
    }

    // --------------------------------------------------------------- keyboard

    const KEY_CODES = {
        Enter: 13, Tab: 9, Backspace: 8, Delete: 46, Escape: 27, ' ': 32,
        ArrowLeft: 37, ArrowUp: 38, ArrowRight: 39, ArrowDown: 40,
        Home: 36, End: 35, PageUp: 33, PageDown: 34, Shift: 16, Control: 17,
        Alt: 18, Meta: 91,
    };

    function keyCode(key) {
        if (KEY_CODES[key] !== undefined) return KEY_CODES[key];
        if (key.length === 1) return key.toUpperCase().charCodeAt(0);
        if (/^F\d{1,2}$/.test(key)) return 111 + Number(key.slice(1));
        return 0;
    }

    function codeFor(key) {
        if (key.length === 1) {
            if (/[a-z]/i.test(key)) return 'Key' + key.toUpperCase();
            if (/[0-9]/.test(key)) return 'Digit' + key;
            if (key === ' ') return 'Space';
            return '';
        }
        return key;
    }

    function keyInit(key, mods) {
        return Object.assign({
            bubbles: true,
            cancelable: true,
            composed: true,
            key: key,
            code: codeFor(key),
            keyCode: keyCode(key),
            which: keyCode(key),
        }, mods || {});
    }

    function editableHost(el) {
        if (!el) return null;
        if (el.isContentEditable) return el;
        const tag = el.tagName ? el.tagName.toLowerCase() : '';
        if (tag === 'textarea') return el;
        if (tag === 'input' && !['checkbox', 'radio', 'button', 'submit', 'reset', 'file', 'image', 'range', 'color'].includes(el.type)) return el;
        return null;
    }

    // Set a field's value through the native setter so frameworks with
    // controlled inputs (React et al.) observe the change.
    function setNativeValue(el, value) {
        const prototype = el instanceof HTMLTextAreaElement
            ? HTMLTextAreaElement.prototype
            : HTMLInputElement.prototype;
        const descriptor = Object.getOwnPropertyDescriptor(prototype, 'value');
        if (descriptor && descriptor.set) {
            descriptor.set.call(el, value);
        } else {
            el.value = value;
        }
    }

    function insertText(el, text) {
        if (el.isContentEditable) {
            const selection = window.getSelection();
            if (selection && selection.rangeCount > 0 && el.contains(selection.anchorNode)) {
                const range = selection.getRangeAt(0);
                range.deleteContents();
                range.insertNode(document.createTextNode(text));
                range.collapse(false);
                selection.removeAllRanges();
                selection.addRange(range);
            } else {
                el.append(document.createTextNode(text));
            }
            return;
        }
        const start = el.selectionStart ?? el.value.length;
        const end = el.selectionEnd ?? el.value.length;
        const value = el.value.slice(0, start) + text + el.value.slice(end);
        setNativeValue(el, value);
        const caret = start + text.length;
        try { el.setSelectionRange(caret, caret); } catch (e) { /* not all inputs support selection */ }
    }

    function deleteBack(el) {
        if (el.isContentEditable) {
            const selection = window.getSelection();
            if (selection && selection.rangeCount > 0) {
                const range = selection.getRangeAt(0);
                if (range.collapsed) selection.modify('extend', 'backward', 'character');
                selection.getRangeAt(0).deleteContents();
            }
            return;
        }
        const start = el.selectionStart ?? el.value.length;
        const end = el.selectionEnd ?? el.value.length;
        if (start === end && start === 0) return;
        const from = start === end ? start - 1 : start;
        setNativeValue(el, el.value.slice(0, from) + el.value.slice(end));
        try { el.setSelectionRange(from, from); } catch (e) { /* ignore */ }
    }

    function fireInput(el, inputType, data) {
        el.dispatchEvent(new InputEvent('input', {
            bubbles: true,
            composed: true,
            inputType: inputType,
            data: data ?? null,
        }));
    }

    async function pressKey(el, key, mods) {
        const init = keyInit(key, mods);
        const keydownNotPrevented = el.dispatchEvent(new KeyboardEvent('keydown', init));
        const editable = editableHost(el);
        if (keydownNotPrevented && editable && !init.ctrlKey && !init.metaKey && !init.altKey) {
            if (key.length === 1) {
                el.dispatchEvent(new KeyboardEvent('keypress', init));
                insertText(editable, key);
                fireInput(editable, 'insertText', key);
            } else if (key === 'Backspace') {
                deleteBack(editable);
                fireInput(editable, 'deleteContentBackward', null);
            } else if (key === 'Enter') {
                if (editable.tagName && editable.tagName.toLowerCase() === 'textarea' || editable.isContentEditable) {
                    insertText(editable, '\n');
                    fireInput(editable, 'insertLineBreak', null);
                } else if (editable.form) {
                    editable.dispatchEvent(new Event('change', { bubbles: true }));
                    editable.form.requestSubmit ? editable.form.requestSubmit() : editable.form.submit();
                }
            }
        }
        if (keydownNotPrevented && key === 'Tab') {
            shiftFocus(!init.shiftKey);
        }
        el.dispatchEvent(new KeyboardEvent('keyup', init));
    }

    function shiftFocus(forward) {
        const focusables = Array.from(document.querySelectorAll(
            'a[href], button, input, textarea, select, [tabindex]:not([tabindex="-1"])',
        )).filter((el) => !el.disabled && el.offsetParent !== null);
        if (focusables.length === 0) return;
        const index = focusables.indexOf(document.activeElement);
        const next = forward
            ? focusables[(index + 1) % focusables.length]
            : focusables[(index - 1 + focusables.length) % focusables.length];
        next.focus();
    }

    async function typeText(text, selector, delay) {
        if (selector) {
            await click({ selector: selector });
        }
        let el = document.activeElement;
        if (!el || el === document.body) {
            throw new Error('no focused element to type into; pass a selector or click a field first');
        }
        const pause = typeof delay === 'number' && delay >= 0 ? delay : 25;
        for (const ch of text) {
            const key = ch === '\n' ? 'Enter' : ch;
            await pressKey(el, key, {});
            if (pause > 0) await sleep(pause);
            el = document.activeElement || el;
        }
        const editable = editableHost(el);
        if (editable) editable.dispatchEvent(new Event('change', { bubbles: true }));
        return describe(el);
    }

    async function press(key, modifiers) {
        const el = document.activeElement || document.body;
        await pressKey(el, key, modifierInit(modifiers));
        return describe(document.activeElement || el);
    }

    // ------------------------------------------------------------- inspection

    function query(selector, limit) {
        const matches = Array.from(document.querySelectorAll(selector));
        const max = limit > 0 ? limit : 25;
        return {
            count: matches.length,
            returned: Math.min(matches.length, max),
            elements: matches.slice(0, max).map((el, index) => {
                const info = describe(el);
                info.index = index;
                const style = getComputedStyle(el);
                info.visible = style.display !== 'none' && style.visibility !== 'hidden' &&
                    el.getBoundingClientRect().width > 0;
                return info;
            }),
        };
    }

    const SKIP_TAGS = new Set(['script', 'style', 'link', 'meta', 'noscript', 'template', 'head']);

    const SNAPSHOT_NODE_BUDGET = 2000;

    function outline(el, depth, maxDepth, budget) {
        if (budget.count >= budget.max) return null;
        if (!el || SKIP_TAGS.has(el.tagName.toLowerCase())) return null;
        const rect = el.getBoundingClientRect();
        const inViewport = rect.bottom > 0 && rect.top < window.innerHeight &&
            rect.right > 0 && rect.left < window.innerWidth;
        if (!inViewport || rect.width === 0 || rect.height === 0) return null;
        budget.count++;
        const node = {
            tag: el.tagName.toLowerCase(),
            bounds: [Math.round(rect.x), Math.round(rect.y), Math.round(rect.width), Math.round(rect.height)],
        };
        if (el.id) node.id = el.id;
        if (el.classList.length) node.classes = Array.from(el.classList).slice(0, 5);
        if ('value' in el && typeof el.value === 'string' && el.value) node.value = el.value.slice(0, 80);
        const ownText = Array.from(el.childNodes)
            .filter((child) => child.nodeType === Node.TEXT_NODE)
            .map((child) => child.textContent.trim())
            .join(' ')
            .trim();
        if (ownText) node.text = ownText.slice(0, 120);
        if (depth < maxDepth) {
            const children = [];
            for (const child of el.children) {
                if (budget.count >= budget.max) break;
                const childNode = outline(child, depth + 1, maxDepth, budget);
                if (childNode) children.push(childNode);
            }
            if (children.length) node.children = children;
        }
        return node;
    }

    function snapshot(maxDepth) {
        const budget = { count: 0, max: SNAPSHOT_NODE_BUDGET };
        return {
            url: location.href,
            title: document.title,
            viewport: { width: window.innerWidth, height: window.innerHeight },
            scroll: { x: Math.round(window.scrollX), y: Math.round(window.scrollY) },
            focused: describe(document.activeElement !== document.body ? document.activeElement : null),
            cursor: { x: state.x, y: state.y },
            tree: outline(document.body, 0, maxDepth > 0 ? maxDepth : 12, budget),
        };
    }

    // ---------------------------------------------------------------- harness

    // run executes fn and posts {id, ok, value|error} to the callback URL.
    // It never throws: all failures are reported through the callback.
    async function run(id, url, showCursor, fn) {
        state.showCursor = !!showCursor;
        let payload;
        try {
            let value = await fn(api);
            if (value === undefined) value = null;
            try {
                JSON.stringify(value);
            } catch (e) {
                value = String(value);
            }
            payload = { id: id, ok: true, value: value };
        } catch (error) {
            payload = {
                id: id,
                ok: false,
                error: String(error && error.stack ? error.stack : error),
            };
        }
        try {
            await fetch(url, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload),
            });
        } catch (error) {
            console.error('[wails-mcp] failed to post result:', error);
        }
    }

    const api = {
        run, move, click, drag, scroll, typeText, press, query, snapshot, describe,
        get cursor() { return { x: state.x, y: state.y }; },
    };

    window.__wailsMCP = api;
})();
