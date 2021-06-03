(function (global, factory) {
    typeof exports === 'object' && typeof module !== 'undefined' ? factory(exports) :
    typeof define === 'function' && define.amd ? define(['exports'], factory) :
    (global = typeof globalThis !== 'undefined' ? globalThis : global || self, factory(global.bridge = {}));
}(this, (function (exports) { 'use strict';

    function noop() { }
    const identity = x => x;
    function run(fn) {
        return fn();
    }
    function blank_object() {
        return Object.create(null);
    }
    function run_all(fns) {
        fns.forEach(run);
    }
    function is_function(thing) {
        return typeof thing === 'function';
    }
    function safe_not_equal(a, b) {
        return a != a ? b == b : a !== b || ((a && typeof a === 'object') || typeof a === 'function');
    }
    function is_empty(obj) {
        return Object.keys(obj).length === 0;
    }
    function subscribe(store, ...callbacks) {
        if (store == null) {
            return noop;
        }
        const unsub = store.subscribe(...callbacks);
        return unsub.unsubscribe ? () => unsub.unsubscribe() : unsub;
    }
    function component_subscribe(component, store, callback) {
        component.$$.on_destroy.push(subscribe(store, callback));
    }
    function action_destroyer(action_result) {
        return action_result && is_function(action_result.destroy) ? action_result.destroy : noop;
    }

    const is_client = typeof window !== 'undefined';
    let now = is_client
        ? () => window.performance.now()
        : () => Date.now();
    let raf = is_client ? cb => requestAnimationFrame(cb) : noop;

    const tasks = new Set();
    function run_tasks(now) {
        tasks.forEach(task => {
            if (!task.c(now)) {
                tasks.delete(task);
                task.f();
            }
        });
        if (tasks.size !== 0)
            raf(run_tasks);
    }
    /**
     * Creates a new task that runs on each raf frame
     * until it returns a falsy value or is aborted
     */
    function loop(callback) {
        let task;
        if (tasks.size === 0)
            raf(run_tasks);
        return {
            promise: new Promise(fulfill => {
                tasks.add(task = { c: callback, f: fulfill });
            }),
            abort() {
                tasks.delete(task);
            }
        };
    }

    function append(target, node) {
        target.appendChild(node);
    }
    function insert(target, node, anchor) {
        target.insertBefore(node, anchor || null);
    }
    function detach(node) {
        node.parentNode.removeChild(node);
    }
    function destroy_each(iterations, detaching) {
        for (let i = 0; i < iterations.length; i += 1) {
            if (iterations[i])
                iterations[i].d(detaching);
        }
    }
    function element(name) {
        return document.createElement(name);
    }
    function text(data) {
        return document.createTextNode(data);
    }
    function space() {
        return text(' ');
    }
    function empty() {
        return text('');
    }
    function listen(node, event, handler, options) {
        node.addEventListener(event, handler, options);
        return () => node.removeEventListener(event, handler, options);
    }
    function attr(node, attribute, value) {
        if (value == null)
            node.removeAttribute(attribute);
        else if (node.getAttribute(attribute) !== value)
            node.setAttribute(attribute, value);
    }
    function children(element) {
        return Array.from(element.childNodes);
    }
    function set_data(text, data) {
        data = '' + data;
        if (text.wholeText !== data)
            text.data = data;
    }
    function custom_event(type, detail) {
        const e = document.createEvent('CustomEvent');
        e.initCustomEvent(type, false, false, detail);
        return e;
    }

    const active_docs = new Set();
    let active = 0;
    // https://github.com/darkskyapp/string-hash/blob/master/index.js
    function hash(str) {
        let hash = 5381;
        let i = str.length;
        while (i--)
            hash = ((hash << 5) - hash) ^ str.charCodeAt(i);
        return hash >>> 0;
    }
    function create_rule(node, a, b, duration, delay, ease, fn, uid = 0) {
        const step = 16.666 / duration;
        let keyframes = '{\n';
        for (let p = 0; p <= 1; p += step) {
            const t = a + (b - a) * ease(p);
            keyframes += p * 100 + `%{${fn(t, 1 - t)}}\n`;
        }
        const rule = keyframes + `100% {${fn(b, 1 - b)}}\n}`;
        const name = `__svelte_${hash(rule)}_${uid}`;
        const doc = node.ownerDocument;
        active_docs.add(doc);
        const stylesheet = doc.__svelte_stylesheet || (doc.__svelte_stylesheet = doc.head.appendChild(element('style')).sheet);
        const current_rules = doc.__svelte_rules || (doc.__svelte_rules = {});
        if (!current_rules[name]) {
            current_rules[name] = true;
            stylesheet.insertRule(`@keyframes ${name} ${rule}`, stylesheet.cssRules.length);
        }
        const animation = node.style.animation || '';
        node.style.animation = `${animation ? `${animation}, ` : ''}${name} ${duration}ms linear ${delay}ms 1 both`;
        active += 1;
        return name;
    }
    function delete_rule(node, name) {
        const previous = (node.style.animation || '').split(', ');
        const next = previous.filter(name
            ? anim => anim.indexOf(name) < 0 // remove specific animation
            : anim => anim.indexOf('__svelte') === -1 // remove all Svelte animations
        );
        const deleted = previous.length - next.length;
        if (deleted) {
            node.style.animation = next.join(', ');
            active -= deleted;
            if (!active)
                clear_rules();
        }
    }
    function clear_rules() {
        raf(() => {
            if (active)
                return;
            active_docs.forEach(doc => {
                const stylesheet = doc.__svelte_stylesheet;
                let i = stylesheet.cssRules.length;
                while (i--)
                    stylesheet.deleteRule(i);
                doc.__svelte_rules = {};
            });
            active_docs.clear();
        });
    }

    let current_component;
    function set_current_component(component) {
        current_component = component;
    }
    function get_current_component() {
        if (!current_component)
            throw new Error('Function called outside component initialization');
        return current_component;
    }
    function onMount(fn) {
        get_current_component().$$.on_mount.push(fn);
    }

    const dirty_components = [];
    const binding_callbacks = [];
    const render_callbacks = [];
    const flush_callbacks = [];
    const resolved_promise = Promise.resolve();
    let update_scheduled = false;
    function schedule_update() {
        if (!update_scheduled) {
            update_scheduled = true;
            resolved_promise.then(flush);
        }
    }
    function add_render_callback(fn) {
        render_callbacks.push(fn);
    }
    let flushing = false;
    const seen_callbacks = new Set();
    function flush() {
        if (flushing)
            return;
        flushing = true;
        do {
            // first, call beforeUpdate functions
            // and update components
            for (let i = 0; i < dirty_components.length; i += 1) {
                const component = dirty_components[i];
                set_current_component(component);
                update(component.$$);
            }
            set_current_component(null);
            dirty_components.length = 0;
            while (binding_callbacks.length)
                binding_callbacks.pop()();
            // then, once components are updated, call
            // afterUpdate functions. This may cause
            // subsequent updates...
            for (let i = 0; i < render_callbacks.length; i += 1) {
                const callback = render_callbacks[i];
                if (!seen_callbacks.has(callback)) {
                    // ...so guard against infinite loops
                    seen_callbacks.add(callback);
                    callback();
                }
            }
            render_callbacks.length = 0;
        } while (dirty_components.length);
        while (flush_callbacks.length) {
            flush_callbacks.pop()();
        }
        update_scheduled = false;
        flushing = false;
        seen_callbacks.clear();
    }
    function update($$) {
        if ($$.fragment !== null) {
            $$.update();
            run_all($$.before_update);
            const dirty = $$.dirty;
            $$.dirty = [-1];
            $$.fragment && $$.fragment.p($$.ctx, dirty);
            $$.after_update.forEach(add_render_callback);
        }
    }

    let promise;
    function wait() {
        if (!promise) {
            promise = Promise.resolve();
            promise.then(() => {
                promise = null;
            });
        }
        return promise;
    }
    function dispatch(node, direction, kind) {
        node.dispatchEvent(custom_event(`${direction ? 'intro' : 'outro'}${kind}`));
    }
    const outroing = new Set();
    let outros;
    function group_outros() {
        outros = {
            r: 0,
            c: [],
            p: outros // parent group
        };
    }
    function check_outros() {
        if (!outros.r) {
            run_all(outros.c);
        }
        outros = outros.p;
    }
    function transition_in(block, local) {
        if (block && block.i) {
            outroing.delete(block);
            block.i(local);
        }
    }
    function transition_out(block, local, detach, callback) {
        if (block && block.o) {
            if (outroing.has(block))
                return;
            outroing.add(block);
            outros.c.push(() => {
                outroing.delete(block);
                if (callback) {
                    if (detach)
                        block.d(1);
                    callback();
                }
            });
            block.o(local);
        }
    }
    const null_transition = { duration: 0 };
    function create_bidirectional_transition(node, fn, params, intro) {
        let config = fn(node, params);
        let t = intro ? 0 : 1;
        let running_program = null;
        let pending_program = null;
        let animation_name = null;
        function clear_animation() {
            if (animation_name)
                delete_rule(node, animation_name);
        }
        function init(program, duration) {
            const d = program.b - t;
            duration *= Math.abs(d);
            return {
                a: t,
                b: program.b,
                d,
                duration,
                start: program.start,
                end: program.start + duration,
                group: program.group
            };
        }
        function go(b) {
            const { delay = 0, duration = 300, easing = identity, tick = noop, css } = config || null_transition;
            const program = {
                start: now() + delay,
                b
            };
            if (!b) {
                // @ts-ignore todo: improve typings
                program.group = outros;
                outros.r += 1;
            }
            if (running_program || pending_program) {
                pending_program = program;
            }
            else {
                // if this is an intro, and there's a delay, we need to do
                // an initial tick and/or apply CSS animation immediately
                if (css) {
                    clear_animation();
                    animation_name = create_rule(node, t, b, duration, delay, easing, css);
                }
                if (b)
                    tick(0, 1);
                running_program = init(program, duration);
                add_render_callback(() => dispatch(node, b, 'start'));
                loop(now => {
                    if (pending_program && now > pending_program.start) {
                        running_program = init(pending_program, duration);
                        pending_program = null;
                        dispatch(node, running_program.b, 'start');
                        if (css) {
                            clear_animation();
                            animation_name = create_rule(node, t, running_program.b, running_program.duration, 0, easing, config.css);
                        }
                    }
                    if (running_program) {
                        if (now >= running_program.end) {
                            tick(t = running_program.b, 1 - t);
                            dispatch(node, running_program.b, 'end');
                            if (!pending_program) {
                                // we're done
                                if (running_program.b) {
                                    // intro — we can tidy up immediately
                                    clear_animation();
                                }
                                else {
                                    // outro — needs to be coordinated
                                    if (!--running_program.group.r)
                                        run_all(running_program.group.c);
                                }
                            }
                            running_program = null;
                        }
                        else if (now >= running_program.start) {
                            const p = now - running_program.start;
                            t = running_program.a + running_program.d * easing(p / running_program.duration);
                            tick(t, 1 - t);
                        }
                    }
                    return !!(running_program || pending_program);
                });
            }
        }
        return {
            run(b) {
                if (is_function(config)) {
                    wait().then(() => {
                        // @ts-ignore
                        config = config();
                        go(b);
                    });
                }
                else {
                    go(b);
                }
            },
            end() {
                clear_animation();
                running_program = pending_program = null;
            }
        };
    }

    const globals = (typeof window !== 'undefined'
        ? window
        : typeof globalThis !== 'undefined'
            ? globalThis
            : global);
    function create_component(block) {
        block && block.c();
    }
    function mount_component(component, target, anchor) {
        const { fragment, on_mount, on_destroy, after_update } = component.$$;
        fragment && fragment.m(target, anchor);
        // onMount happens before the initial afterUpdate
        add_render_callback(() => {
            const new_on_destroy = on_mount.map(run).filter(is_function);
            if (on_destroy) {
                on_destroy.push(...new_on_destroy);
            }
            else {
                // Edge case - component was destroyed immediately,
                // most likely as a result of a binding initialising
                run_all(new_on_destroy);
            }
            component.$$.on_mount = [];
        });
        after_update.forEach(add_render_callback);
    }
    function destroy_component(component, detaching) {
        const $$ = component.$$;
        if ($$.fragment !== null) {
            run_all($$.on_destroy);
            $$.fragment && $$.fragment.d(detaching);
            // TODO null out other refs, including component.$$ (but need to
            // preserve final state?)
            $$.on_destroy = $$.fragment = null;
            $$.ctx = [];
        }
    }
    function make_dirty(component, i) {
        if (component.$$.dirty[0] === -1) {
            dirty_components.push(component);
            schedule_update();
            component.$$.dirty.fill(0);
        }
        component.$$.dirty[(i / 31) | 0] |= (1 << (i % 31));
    }
    function init(component, options, instance, create_fragment, not_equal, props, dirty = [-1]) {
        const parent_component = current_component;
        set_current_component(component);
        const $$ = component.$$ = {
            fragment: null,
            ctx: null,
            // state
            props,
            update: noop,
            not_equal,
            bound: blank_object(),
            // lifecycle
            on_mount: [],
            on_destroy: [],
            before_update: [],
            after_update: [],
            context: new Map(parent_component ? parent_component.$$.context : []),
            // everything else
            callbacks: blank_object(),
            dirty,
            skip_bound: false
        };
        let ready = false;
        $$.ctx = instance
            ? instance(component, options.props || {}, (i, ret, ...rest) => {
                const value = rest.length ? rest[0] : ret;
                if ($$.ctx && not_equal($$.ctx[i], $$.ctx[i] = value)) {
                    if (!$$.skip_bound && $$.bound[i])
                        $$.bound[i](value);
                    if (ready)
                        make_dirty(component, i);
                }
                return ret;
            })
            : [];
        $$.update();
        ready = true;
        run_all($$.before_update);
        // `false` as a special case of no DOM component
        $$.fragment = create_fragment ? create_fragment($$.ctx) : false;
        if (options.target) {
            if (options.hydrate) {
                const nodes = children(options.target);
                // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
                $$.fragment && $$.fragment.l(nodes);
                nodes.forEach(detach);
            }
            else {
                // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
                $$.fragment && $$.fragment.c();
            }
            if (options.intro)
                transition_in(component.$$.fragment);
            mount_component(component, options.target, options.anchor);
            flush();
        }
        set_current_component(parent_component);
    }
    /**
     * Base class for Svelte components. Used when dev=false.
     */
    class SvelteComponent {
        $destroy() {
            destroy_component(this, 1);
            this.$destroy = noop;
        }
        $on(type, callback) {
            const callbacks = (this.$$.callbacks[type] || (this.$$.callbacks[type] = []));
            callbacks.push(callback);
            return () => {
                const index = callbacks.indexOf(callback);
                if (index !== -1)
                    callbacks.splice(index, 1);
            };
        }
        $set($$props) {
            if (this.$$set && !is_empty($$props)) {
                this.$$.skip_bound = true;
                this.$$set($$props);
                this.$$.skip_bound = false;
            }
        }
    }

    const subscriber_queue = [];
    /**
     * Create a `Writable` store that allows both updating and reading by subscription.
     * @param {*=}value initial value
     * @param {StartStopNotifier=}start start and stop notifications for subscriptions
     */
    function writable(value, start = noop) {
        let stop;
        const subscribers = [];
        function set(new_value) {
            if (safe_not_equal(value, new_value)) {
                value = new_value;
                if (stop) { // store is ready
                    const run_queue = !subscriber_queue.length;
                    for (let i = 0; i < subscribers.length; i += 1) {
                        const s = subscribers[i];
                        s[1]();
                        subscriber_queue.push(s, value);
                    }
                    if (run_queue) {
                        for (let i = 0; i < subscriber_queue.length; i += 2) {
                            subscriber_queue[i][0](subscriber_queue[i + 1]);
                        }
                        subscriber_queue.length = 0;
                    }
                }
            }
        }
        function update(fn) {
            set(fn(value));
        }
        function subscribe(run, invalidate = noop) {
            const subscriber = [run, invalidate];
            subscribers.push(subscriber);
            if (subscribers.length === 1) {
                stop = start(set) || noop;
            }
            run(value);
            return () => {
                const index = subscribers.indexOf(subscriber);
                if (index !== -1) {
                    subscribers.splice(index, 1);
                }
                if (subscribers.length === 0) {
                    stop();
                    stop = null;
                }
            };
        }
        return { set, update, subscribe };
    }

    function log(message) {
        // eslint-disable-next-line
        console.log(
            '%c wails bridge %c ' + message + ' ',
            'background: #aa0000; color: #fff; border-radius: 3px 0px 0px 3px; padding: 1px; font-size: 0.7rem',
            'background: #009900; color: #fff; border-radius: 0px 3px 3px 0px; padding: 1px; font-size: 0.7rem'
        );
    }

    /** Overlay */
    const overlayVisible = writable(false);

    function showOverlay() {
        overlayVisible.set(true);
    }
    function hideOverlay() {
        overlayVisible.set(false);
    }

    /** Menubar **/
    const menuVisible = writable(false);

    /** Trays **/

    const trays = writable([]);
    function setTray(tray) {
        trays.update((current) => {
            // Remove existing if it exists, else add
            const index = current.findIndex(item => item.ID === tray.ID);
            if ( index === -1 ) {
                current.push(tray);
            } else {
                current[index] = tray;
            }
            return current;
        });
    }
    function updateTrayLabel(tray) {
        trays.update((current) => {
            // Remove existing if it exists, else add
            const index = current.findIndex(item => item.ID === tray.ID);
            if ( index === -1 ) {
                return log("ERROR: Attempted to update tray index ", tray.ID)
            }
            current[index].Label = tray.Label;
            return current;
        });
    }

    function deleteTrayMenu(id) {
        trays.update((current) => {
            // Remove existing if it exists, else add
            const index = current.findIndex(item => item.ID === id);
            if ( index === -1 ) {
                return log("ERROR: Attempted to delete tray index ")
            }
            current.splice(index, 1);
            return current;
        });
    }

    let selectedMenu = writable(null);

    function fade(node, { delay = 0, duration = 400, easing = identity } = {}) {
        const o = +getComputedStyle(node).opacity;
        return {
            delay,
            duration,
            easing,
            css: t => `opacity: ${t * o}`
        };
    }

    /* Overlay.svelte generated by Svelte v3.32.2 */

    function add_css() {
    	var style = element("style");
    	style.id = "svelte-9nqyfr-style";
    	style.textContent = ".wails-reconnect-overlay.svelte-9nqyfr{position:fixed;top:0;left:0;width:100%;height:100%;backdrop-filter:blur(20px) saturate(160%) contrast(45%) brightness(140%);z-index:999999\r\n    }.wails-reconnect-overlay-content.svelte-9nqyfr{position:relative;top:50%;transform:translateY(-50%);margin:0;background-image:url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAC8AAAAuCAMAAACPpbA7AAAAflBMVEUAAAAAAAAAAAAAAAAAAAAAAAAAAAAEBAQAAAAAAAAAAAABAQEEBAQAAAAAAAAEBAQAAAADAwMAAAABAQEAAAAAAAAAAAAAAAAAAAACAgICAgIBAQEAAAAAAAAAAAAAAAAAAAACAgIAAAAAAAAAAAAAAAAAAAAAAAAAAAAFBQWCC3waAAAAKXRSTlMALgUMIBk0+xEqJs70Xhb3lu3EjX2EZTlv5eHXvbarQj3cdmpXSqOeUDwaqNAAAAKCSURBVEjHjZTntqsgEIUPVVCwtxg1vfD+L3hHRe8K6snZf+KKn8OewvzsSSeXLruLnz+KHs0gr6DkT3xsRkU6VVn4Ha/UxLe1Z4y64i847sykPBh/AvQ7ry3eFN70oKrfcBJYvm/tQ1qxP4T3emXPeXAkvodPUvtdjbhk+Ft4c0hslTiXVOzxOJ15NWUblQhRsdu3E1AfCjj3Gdm18zSOsiH8Lk4TB480ksy62fiqNo4OpyU8O21l6+hyRtS6z8r1pHlmle5sR1/WXS6Mq2Nl+YeKt3vr+vdH/q4O68tzXuwkiZmngYb4R8Co1jh0+Ww2UTyWxBvtyxLO7QVjO3YOD/lWZpbXDGellFG2Mws58mMnjVZSn7p+XvZ6IF4nn02OJZV0aTO22arp/DgLPtrgpVoi6TPbZm4XQBjY159w02uO0BDdYsfrOEi0M2ulRXlCIPAOuN1NOVhi+riBR3dgwQplYsZRZJLXq23Mlo5njkbY0rZFu3oiNIYG2kqsbVz67OlNuZZIOlfxHDl0UpyRX86z/OYC/3qf1A1xTrMp/PWWM4ePzf8DDp1nesQRpcFk7BlwdzN08ZIALJpCaciQXO0f6k4dnuT/Ewg4l7qSTNzm2SykdHn6GJ12mWc6aCNj/g1cTXpB8YFfr0uVc96aFkkqiIiX4nO+salKwGtIkvfB+Ja8DxMeD3hIXP5mTOYPB4eVT0+32I5ykvPZjesnkGgIREgYnmLrPb0PdV3hoLup2TjcGBPM4mgsfF5BrawZR4/GpzYQzQfrUZCf0TCWYo2DqhdhTJBQ6j4xqmmLN5LjdRIY8LWExiFUsSrza/nmFBqw3I9tEZB9h0lIQSO9if8DkISDAj8CDawAAAAASUVORK5CYII=);background-repeat:no-repeat;background-position:center\r\n    }.wails-reconnect-overlay-loadingspinner.svelte-9nqyfr{pointer-events:none;width:2.5em;height:2.5em;border:.4em solid transparent;border-color:#f00 #eee0 #f00 #eee0;border-radius:50%;animation:svelte-9nqyfr-loadingspin 1s linear infinite;margin:auto;padding:2.5em\r\n    }@keyframes svelte-9nqyfr-loadingspin{100%{transform:rotate(360deg)}}";
    	append(document.head, style);
    }

    // (8:0) {#if $overlayVisible }
    function create_if_block(ctx) {
    	let div2;
    	let div2_transition;
    	let current;

    	return {
    		c() {
    			div2 = element("div");
    			div2.innerHTML = `<div class="wails-reconnect-overlay-content svelte-9nqyfr"><div class="wails-reconnect-overlay-loadingspinner svelte-9nqyfr"></div></div>`;
    			attr(div2, "class", "wails-reconnect-overlay svelte-9nqyfr");
    		},
    		m(target, anchor) {
    			insert(target, div2, anchor);
    			current = true;
    		},
    		i(local) {
    			if (current) return;

    			add_render_callback(() => {
    				if (!div2_transition) div2_transition = create_bidirectional_transition(div2, fade, { duration: 200 }, true);
    				div2_transition.run(1);
    			});

    			current = true;
    		},
    		o(local) {
    			if (!div2_transition) div2_transition = create_bidirectional_transition(div2, fade, { duration: 200 }, false);
    			div2_transition.run(0);
    			current = false;
    		},
    		d(detaching) {
    			if (detaching) detach(div2);
    			if (detaching && div2_transition) div2_transition.end();
    		}
    	};
    }

    function create_fragment(ctx) {
    	let if_block_anchor;
    	let current;
    	let if_block = /*$overlayVisible*/ ctx[0] && create_if_block();

    	return {
    		c() {
    			if (if_block) if_block.c();
    			if_block_anchor = empty();
    		},
    		m(target, anchor) {
    			if (if_block) if_block.m(target, anchor);
    			insert(target, if_block_anchor, anchor);
    			current = true;
    		},
    		p(ctx, [dirty]) {
    			if (/*$overlayVisible*/ ctx[0]) {
    				if (if_block) {
    					if (dirty & /*$overlayVisible*/ 1) {
    						transition_in(if_block, 1);
    					}
    				} else {
    					if_block = create_if_block();
    					if_block.c();
    					transition_in(if_block, 1);
    					if_block.m(if_block_anchor.parentNode, if_block_anchor);
    				}
    			} else if (if_block) {
    				group_outros();

    				transition_out(if_block, 1, 1, () => {
    					if_block = null;
    				});

    				check_outros();
    			}
    		},
    		i(local) {
    			if (current) return;
    			transition_in(if_block);
    			current = true;
    		},
    		o(local) {
    			transition_out(if_block);
    			current = false;
    		},
    		d(detaching) {
    			if (if_block) if_block.d(detaching);
    			if (detaching) detach(if_block_anchor);
    		}
    	};
    }

    function instance($$self, $$props, $$invalidate) {
    	let $overlayVisible;
    	component_subscribe($$self, overlayVisible, $$value => $$invalidate(0, $overlayVisible = $$value));
    	return [$overlayVisible];
    }

    class Overlay extends SvelteComponent {
    	constructor(options) {
    		super();
    		if (!document.getElementById("svelte-9nqyfr-style")) add_css();
    		init(this, options, instance, create_fragment, safe_not_equal, {});
    	}
    }

    /* Menu.svelte generated by Svelte v3.32.2 */

    function add_css$1() {
    	var style = element("style");
    	style.id = "svelte-1oysp7o-style";
    	style.textContent = ".menu.svelte-1oysp7o.svelte-1oysp7o{padding:3px;background-color:#0008;color:#EEF;border-radius:5px;margin-top:5px;position:absolute;backdrop-filter:blur(3px) saturate(160%) contrast(45%) brightness(140%);border:1px solid rgb(88,88,88);box-shadow:0 0 1px rgb(146,146,148) inset}.menuitem.svelte-1oysp7o.svelte-1oysp7o{display:flex;align-items:center;padding:1px 5px}.menuitem.svelte-1oysp7o.svelte-1oysp7o:hover{display:flex;align-items:center;background-color:rgb(57,131,223);padding:1px 5px;border-radius:5px}.menuitem.svelte-1oysp7o img.svelte-1oysp7o{padding-right:5px}";
    	append(document.head, style);
    }

    function get_each_context(ctx, list, i) {
    	const child_ctx = ctx.slice();
    	child_ctx[2] = list[i];
    	return child_ctx;
    }

    // (8:0) {#if !hidden}
    function create_if_block$1(ctx) {
    	let div;
    	let if_block = /*menu*/ ctx[0].Menu && create_if_block_1(ctx);

    	return {
    		c() {
    			div = element("div");
    			if (if_block) if_block.c();
    			attr(div, "class", "menu svelte-1oysp7o");
    		},
    		m(target, anchor) {
    			insert(target, div, anchor);
    			if (if_block) if_block.m(div, null);
    		},
    		p(ctx, dirty) {
    			if (/*menu*/ ctx[0].Menu) {
    				if (if_block) {
    					if_block.p(ctx, dirty);
    				} else {
    					if_block = create_if_block_1(ctx);
    					if_block.c();
    					if_block.m(div, null);
    				}
    			} else if (if_block) {
    				if_block.d(1);
    				if_block = null;
    			}
    		},
    		d(detaching) {
    			if (detaching) detach(div);
    			if (if_block) if_block.d();
    		}
    	};
    }

    // (10:4) {#if menu.Menu }
    function create_if_block_1(ctx) {
    	let each_1_anchor;
    	let each_value = /*menu*/ ctx[0].Menu.Items;
    	let each_blocks = [];

    	for (let i = 0; i < each_value.length; i += 1) {
    		each_blocks[i] = create_each_block(get_each_context(ctx, each_value, i));
    	}

    	return {
    		c() {
    			for (let i = 0; i < each_blocks.length; i += 1) {
    				each_blocks[i].c();
    			}

    			each_1_anchor = empty();
    		},
    		m(target, anchor) {
    			for (let i = 0; i < each_blocks.length; i += 1) {
    				each_blocks[i].m(target, anchor);
    			}

    			insert(target, each_1_anchor, anchor);
    		},
    		p(ctx, dirty) {
    			if (dirty & /*menu*/ 1) {
    				each_value = /*menu*/ ctx[0].Menu.Items;
    				let i;

    				for (i = 0; i < each_value.length; i += 1) {
    					const child_ctx = get_each_context(ctx, each_value, i);

    					if (each_blocks[i]) {
    						each_blocks[i].p(child_ctx, dirty);
    					} else {
    						each_blocks[i] = create_each_block(child_ctx);
    						each_blocks[i].c();
    						each_blocks[i].m(each_1_anchor.parentNode, each_1_anchor);
    					}
    				}

    				for (; i < each_blocks.length; i += 1) {
    					each_blocks[i].d(1);
    				}

    				each_blocks.length = each_value.length;
    			}
    		},
    		d(detaching) {
    			destroy_each(each_blocks, detaching);
    			if (detaching) detach(each_1_anchor);
    		}
    	};
    }

    // (13:12) {#if menuItem.Image }
    function create_if_block_2(ctx) {
    	let div;
    	let img;
    	let img_src_value;

    	return {
    		c() {
    			div = element("div");
    			img = element("img");
    			attr(img, "alt", "");
    			if (img.src !== (img_src_value = "data:image/png;base64," + /*menuItem*/ ctx[2].Image)) attr(img, "src", img_src_value);
    			attr(img, "class", "svelte-1oysp7o");
    		},
    		m(target, anchor) {
    			insert(target, div, anchor);
    			append(div, img);
    		},
    		p(ctx, dirty) {
    			if (dirty & /*menu*/ 1 && img.src !== (img_src_value = "data:image/png;base64," + /*menuItem*/ ctx[2].Image)) {
    				attr(img, "src", img_src_value);
    			}
    		},
    		d(detaching) {
    			if (detaching) detach(div);
    		}
    	};
    }

    // (11:8) {#each menu.Menu.Items as menuItem}
    function create_each_block(ctx) {
    	let div1;
    	let t0;
    	let div0;
    	let t1_value = /*menuItem*/ ctx[2].Label + "";
    	let t1;
    	let t2;
    	let if_block = /*menuItem*/ ctx[2].Image && create_if_block_2(ctx);

    	return {
    		c() {
    			div1 = element("div");
    			if (if_block) if_block.c();
    			t0 = space();
    			div0 = element("div");
    			t1 = text(t1_value);
    			t2 = space();
    			attr(div0, "class", "menulabel");
    			attr(div1, "class", "menuitem svelte-1oysp7o");
    		},
    		m(target, anchor) {
    			insert(target, div1, anchor);
    			if (if_block) if_block.m(div1, null);
    			append(div1, t0);
    			append(div1, div0);
    			append(div0, t1);
    			append(div1, t2);
    		},
    		p(ctx, dirty) {
    			if (/*menuItem*/ ctx[2].Image) {
    				if (if_block) {
    					if_block.p(ctx, dirty);
    				} else {
    					if_block = create_if_block_2(ctx);
    					if_block.c();
    					if_block.m(div1, t0);
    				}
    			} else if (if_block) {
    				if_block.d(1);
    				if_block = null;
    			}

    			if (dirty & /*menu*/ 1 && t1_value !== (t1_value = /*menuItem*/ ctx[2].Label + "")) set_data(t1, t1_value);
    		},
    		d(detaching) {
    			if (detaching) detach(div1);
    			if (if_block) if_block.d();
    		}
    	};
    }

    function create_fragment$1(ctx) {
    	let if_block_anchor;
    	let if_block = !/*hidden*/ ctx[1] && create_if_block$1(ctx);

    	return {
    		c() {
    			if (if_block) if_block.c();
    			if_block_anchor = empty();
    		},
    		m(target, anchor) {
    			if (if_block) if_block.m(target, anchor);
    			insert(target, if_block_anchor, anchor);
    		},
    		p(ctx, [dirty]) {
    			if (!/*hidden*/ ctx[1]) {
    				if (if_block) {
    					if_block.p(ctx, dirty);
    				} else {
    					if_block = create_if_block$1(ctx);
    					if_block.c();
    					if_block.m(if_block_anchor.parentNode, if_block_anchor);
    				}
    			} else if (if_block) {
    				if_block.d(1);
    				if_block = null;
    			}
    		},
    		i: noop,
    		o: noop,
    		d(detaching) {
    			if (if_block) if_block.d(detaching);
    			if (detaching) detach(if_block_anchor);
    		}
    	};
    }

    function instance$1($$self, $$props, $$invalidate) {
    	let { menu } = $$props;
    	let { hidden = true } = $$props;

    	$$self.$$set = $$props => {
    		if ("menu" in $$props) $$invalidate(0, menu = $$props.menu);
    		if ("hidden" in $$props) $$invalidate(1, hidden = $$props.hidden);
    	};

    	return [menu, hidden];
    }

    class Menu extends SvelteComponent {
    	constructor(options) {
    		super();
    		if (!document.getElementById("svelte-1oysp7o-style")) add_css$1();
    		init(this, options, instance$1, create_fragment$1, safe_not_equal, { menu: 0, hidden: 1 });
    	}
    }

    /* TrayMenu.svelte generated by Svelte v3.32.2 */

    const { document: document_1 } = globals;

    function add_css$2() {
    	var style = element("style");
    	style.id = "svelte-esze1k-style";
    	style.textContent = ".tray-menu.svelte-esze1k{padding-left:0.5rem;padding-right:0.5rem;overflow:visible;font-size:14px}.label.svelte-esze1k{text-align:right;padding-right:10px}";
    	append(document_1.head, style);
    }

    // (47:4) {#if tray.ProcessedMenu }
    function create_if_block$2(ctx) {
    	let menu;
    	let current;

    	menu = new Menu({
    			props: {
    				menu: /*tray*/ ctx[0].ProcessedMenu,
    				hidden: /*hidden*/ ctx[1]
    			}
    		});

    	return {
    		c() {
    			create_component(menu.$$.fragment);
    		},
    		m(target, anchor) {
    			mount_component(menu, target, anchor);
    			current = true;
    		},
    		p(ctx, dirty) {
    			const menu_changes = {};
    			if (dirty & /*tray*/ 1) menu_changes.menu = /*tray*/ ctx[0].ProcessedMenu;
    			if (dirty & /*hidden*/ 2) menu_changes.hidden = /*hidden*/ ctx[1];
    			menu.$set(menu_changes);
    		},
    		i(local) {
    			if (current) return;
    			transition_in(menu.$$.fragment, local);
    			current = true;
    		},
    		o(local) {
    			transition_out(menu.$$.fragment, local);
    			current = false;
    		},
    		d(detaching) {
    			destroy_component(menu, detaching);
    		}
    	};
    }

    function create_fragment$2(ctx) {
    	let span1;
    	let span0;
    	let t0_value = /*tray*/ ctx[0].Label + "";
    	let t0;
    	let t1;
    	let current;
    	let mounted;
    	let dispose;
    	let if_block = /*tray*/ ctx[0].ProcessedMenu && create_if_block$2(ctx);

    	return {
    		c() {
    			span1 = element("span");
    			span0 = element("span");
    			t0 = text(t0_value);
    			t1 = space();
    			if (if_block) if_block.c();
    			attr(span0, "class", "label svelte-esze1k");
    			attr(span1, "class", "tray-menu svelte-esze1k");
    		},
    		m(target, anchor) {
    			insert(target, span1, anchor);
    			append(span1, span0);
    			append(span0, t0);
    			append(span1, t1);
    			if (if_block) if_block.m(span1, null);
    			current = true;

    			if (!mounted) {
    				dispose = [
    					listen(span0, "click", /*trayClicked*/ ctx[3]),
    					action_destroyer(clickOutside.call(null, span1)),
    					listen(span1, "click_outside", /*closeMenu*/ ctx[2])
    				];

    				mounted = true;
    			}
    		},
    		p(ctx, [dirty]) {
    			if ((!current || dirty & /*tray*/ 1) && t0_value !== (t0_value = /*tray*/ ctx[0].Label + "")) set_data(t0, t0_value);

    			if (/*tray*/ ctx[0].ProcessedMenu) {
    				if (if_block) {
    					if_block.p(ctx, dirty);

    					if (dirty & /*tray*/ 1) {
    						transition_in(if_block, 1);
    					}
    				} else {
    					if_block = create_if_block$2(ctx);
    					if_block.c();
    					transition_in(if_block, 1);
    					if_block.m(span1, null);
    				}
    			} else if (if_block) {
    				group_outros();

    				transition_out(if_block, 1, 1, () => {
    					if_block = null;
    				});

    				check_outros();
    			}
    		},
    		i(local) {
    			if (current) return;
    			transition_in(if_block);
    			current = true;
    		},
    		o(local) {
    			transition_out(if_block);
    			current = false;
    		},
    		d(detaching) {
    			if (detaching) detach(span1);
    			if (if_block) if_block.d();
    			mounted = false;
    			run_all(dispose);
    		}
    	};
    }

    function clickOutside(node) {
    	const handleClick = event => {
    		if (node && !node.contains(event.target) && !event.defaultPrevented) {
    			node.dispatchEvent(new CustomEvent("click_outside", node));
    		}
    	};

    	document.addEventListener("click", handleClick, true);

    	return {
    		destroy() {
    			document.removeEventListener("click", handleClick, true);
    		}
    	};
    }

    function instance$2($$self, $$props, $$invalidate) {
    	let hidden;
    	let $selectedMenu;
    	component_subscribe($$self, selectedMenu, $$value => $$invalidate(4, $selectedMenu = $$value));
    	let { tray = null } = $$props;

    	function closeMenu() {
    		selectedMenu.set(null);
    	}

    	function trayClicked() {
    		if ($selectedMenu !== tray) {
    			selectedMenu.set(tray);
    		} else {
    			selectedMenu.set(null);
    		}
    	}

    	$$self.$$set = $$props => {
    		if ("tray" in $$props) $$invalidate(0, tray = $$props.tray);
    	};

    	$$self.$$.update = () => {
    		if ($$self.$$.dirty & /*$selectedMenu, tray*/ 17) {
    			$$invalidate(1, hidden = $selectedMenu !== tray);
    		}
    	};

    	return [tray, hidden, closeMenu, trayClicked, $selectedMenu];
    }

    class TrayMenu extends SvelteComponent {
    	constructor(options) {
    		super();
    		if (!document_1.getElementById("svelte-esze1k-style")) add_css$2();
    		init(this, options, instance$2, create_fragment$2, safe_not_equal, { tray: 0 });
    	}
    }

    /* Menubar.svelte generated by Svelte v3.32.2 */

    function add_css$3() {
    	var style = element("style");
    	style.id = "svelte-1i0zb4n-style";
    	style.textContent = ".tray-menus.svelte-1i0zb4n{display:flex;flex-direction:row;justify-content:flex-end}.wails-menubar.svelte-1i0zb4n{position:relative;display:block;top:0;height:2rem;width:100%;border-bottom:1px solid #b3b3b3;box-shadow:0 0 10px 0 #33333360}.time.svelte-1i0zb4n{padding-left:0.5rem;padding-right:1.5rem;overflow:visible;font-size:14px}";
    	append(document.head, style);
    }

    function get_each_context$1(ctx, list, i) {
    	const child_ctx = ctx.slice();
    	child_ctx[9] = list[i];
    	return child_ctx;
    }

    // (38:0) {#if $menuVisible }
    function create_if_block$3(ctx) {
    	let div;
    	let span1;
    	let t0;
    	let span0;
    	let t1;
    	let div_transition;
    	let current;
    	let each_value = /*$trays*/ ctx[2];
    	let each_blocks = [];

    	for (let i = 0; i < each_value.length; i += 1) {
    		each_blocks[i] = create_each_block$1(get_each_context$1(ctx, each_value, i));
    	}

    	const out = i => transition_out(each_blocks[i], 1, 1, () => {
    		each_blocks[i] = null;
    	});

    	return {
    		c() {
    			div = element("div");
    			span1 = element("span");

    			for (let i = 0; i < each_blocks.length; i += 1) {
    				each_blocks[i].c();
    			}

    			t0 = space();
    			span0 = element("span");
    			t1 = text(/*dateTimeString*/ ctx[0]);
    			attr(span0, "class", "time svelte-1i0zb4n");
    			attr(span1, "class", "tray-menus svelte-1i0zb4n");
    			attr(div, "class", "wails-menubar svelte-1i0zb4n");
    		},
    		m(target, anchor) {
    			insert(target, div, anchor);
    			append(div, span1);

    			for (let i = 0; i < each_blocks.length; i += 1) {
    				each_blocks[i].m(span1, null);
    			}

    			append(span1, t0);
    			append(span1, span0);
    			append(span0, t1);
    			current = true;
    		},
    		p(ctx, dirty) {
    			if (dirty & /*$trays*/ 4) {
    				each_value = /*$trays*/ ctx[2];
    				let i;

    				for (i = 0; i < each_value.length; i += 1) {
    					const child_ctx = get_each_context$1(ctx, each_value, i);

    					if (each_blocks[i]) {
    						each_blocks[i].p(child_ctx, dirty);
    						transition_in(each_blocks[i], 1);
    					} else {
    						each_blocks[i] = create_each_block$1(child_ctx);
    						each_blocks[i].c();
    						transition_in(each_blocks[i], 1);
    						each_blocks[i].m(span1, t0);
    					}
    				}

    				group_outros();

    				for (i = each_value.length; i < each_blocks.length; i += 1) {
    					out(i);
    				}

    				check_outros();
    			}

    			if (!current || dirty & /*dateTimeString*/ 1) set_data(t1, /*dateTimeString*/ ctx[0]);
    		},
    		i(local) {
    			if (current) return;

    			for (let i = 0; i < each_value.length; i += 1) {
    				transition_in(each_blocks[i]);
    			}

    			add_render_callback(() => {
    				if (!div_transition) div_transition = create_bidirectional_transition(div, fade, {}, true);
    				div_transition.run(1);
    			});

    			current = true;
    		},
    		o(local) {
    			each_blocks = each_blocks.filter(Boolean);

    			for (let i = 0; i < each_blocks.length; i += 1) {
    				transition_out(each_blocks[i]);
    			}

    			if (!div_transition) div_transition = create_bidirectional_transition(div, fade, {}, false);
    			div_transition.run(0);
    			current = false;
    		},
    		d(detaching) {
    			if (detaching) detach(div);
    			destroy_each(each_blocks, detaching);
    			if (detaching && div_transition) div_transition.end();
    		}
    	};
    }

    // (41:4) {#each $trays as tray}
    function create_each_block$1(ctx) {
    	let traymenu;
    	let current;
    	traymenu = new TrayMenu({ props: { tray: /*tray*/ ctx[9] } });

    	return {
    		c() {
    			create_component(traymenu.$$.fragment);
    		},
    		m(target, anchor) {
    			mount_component(traymenu, target, anchor);
    			current = true;
    		},
    		p(ctx, dirty) {
    			const traymenu_changes = {};
    			if (dirty & /*$trays*/ 4) traymenu_changes.tray = /*tray*/ ctx[9];
    			traymenu.$set(traymenu_changes);
    		},
    		i(local) {
    			if (current) return;
    			transition_in(traymenu.$$.fragment, local);
    			current = true;
    		},
    		o(local) {
    			transition_out(traymenu.$$.fragment, local);
    			current = false;
    		},
    		d(detaching) {
    			destroy_component(traymenu, detaching);
    		}
    	};
    }

    function create_fragment$3(ctx) {
    	let if_block_anchor;
    	let current;
    	let mounted;
    	let dispose;
    	let if_block = /*$menuVisible*/ ctx[1] && create_if_block$3(ctx);

    	return {
    		c() {
    			if (if_block) if_block.c();
    			if_block_anchor = empty();
    		},
    		m(target, anchor) {
    			if (if_block) if_block.m(target, anchor);
    			insert(target, if_block_anchor, anchor);
    			current = true;

    			if (!mounted) {
    				dispose = listen(window, "keydown", /*handleKeydown*/ ctx[3]);
    				mounted = true;
    			}
    		},
    		p(ctx, [dirty]) {
    			if (/*$menuVisible*/ ctx[1]) {
    				if (if_block) {
    					if_block.p(ctx, dirty);

    					if (dirty & /*$menuVisible*/ 2) {
    						transition_in(if_block, 1);
    					}
    				} else {
    					if_block = create_if_block$3(ctx);
    					if_block.c();
    					transition_in(if_block, 1);
    					if_block.m(if_block_anchor.parentNode, if_block_anchor);
    				}
    			} else if (if_block) {
    				group_outros();

    				transition_out(if_block, 1, 1, () => {
    					if_block = null;
    				});

    				check_outros();
    			}
    		},
    		i(local) {
    			if (current) return;
    			transition_in(if_block);
    			current = true;
    		},
    		o(local) {
    			transition_out(if_block);
    			current = false;
    		},
    		d(detaching) {
    			if (if_block) if_block.d(detaching);
    			if (detaching) detach(if_block_anchor);
    			mounted = false;
    			dispose();
    		}
    	};
    }

    function instance$3($$self, $$props, $$invalidate) {
    	let day;
    	let dom;
    	let mon;
    	let currentTime;
    	let dateTimeString;
    	let $menuVisible;
    	let $trays;
    	component_subscribe($$self, menuVisible, $$value => $$invalidate(1, $menuVisible = $$value));
    	component_subscribe($$self, trays, $$value => $$invalidate(2, $trays = $$value));
    	let time = new Date();

    	onMount(() => {
    		const interval = setInterval(
    			() => {
    				$$invalidate(4, time = new Date());
    			},
    			1000
    		);

    		return () => {
    			clearInterval(interval);
    		};
    	});

    	function handleKeydown(e) {
    		// Backtick toggle
    		if (e.keyCode == 192) {
    			menuVisible.update(current => {
    				return !current;
    			});
    		}
    	}

    	$$self.$$.update = () => {
    		if ($$self.$$.dirty & /*time*/ 16) {
    			$$invalidate(5, day = time.toLocaleString("default", { weekday: "short" }));
    		}

    		if ($$self.$$.dirty & /*time*/ 16) {
    			$$invalidate(6, dom = time.getDate());
    		}

    		if ($$self.$$.dirty & /*time*/ 16) {
    			$$invalidate(7, mon = time.toLocaleString("default", { month: "short" }));
    		}

    		if ($$self.$$.dirty & /*time*/ 16) {
    			$$invalidate(8, currentTime = time.toLocaleString("en-US", {
    				hour: "numeric",
    				minute: "numeric",
    				hour12: true
    			}).toLowerCase());
    		}

    		if ($$self.$$.dirty & /*day, dom, mon, currentTime*/ 480) {
    			$$invalidate(0, dateTimeString = `${day} ${dom} ${mon} ${currentTime}`);
    		}
    	};

    	return [
    		dateTimeString,
    		$menuVisible,
    		$trays,
    		handleKeydown,
    		time,
    		day,
    		dom,
    		mon,
    		currentTime
    	];
    }

    class Menubar extends SvelteComponent {
    	constructor(options) {
    		super();
    		if (!document.getElementById("svelte-1i0zb4n-style")) add_css$3();
    		init(this, options, instance$3, create_fragment$3, safe_not_equal, {});
    	}
    }

    /*
     _       __      _ __
    | |     / /___ _(_) /____
    | | /| / / __ `/ / / ___/
    | |/ |/ / /_/ / / (__  )
    |__/|__/\__,_/_/_/____/
    The lightweight framework for web-like apps
    (c) Lea Anthony 2019-present
    */

    let websocket = null;
    let callback = null;
    let connectTimer;

    function StartWebsocket(userCallback) {

    	callback = userCallback;

    	window.onbeforeunload = function() {
    		if( websocket ) {
    			websocket.onclose = function () { };
    			websocket.close();
    			websocket = null;
    		}
    	};

    	// ...and attempt to connect
    	connect();

    }

    function setupIPCBridge() {
    	window.wailsInvoke = (message) => {
    		websocket.send(message);
    	};
        window.wailsDrag = (message) => {
            websocket.send(message);
        };
        window.wailsContextMenuMessage = (message) => {
            websocket.send(message);
        };
    }

    // Handles incoming websocket connections
    function handleConnect() {
    	log('Connected to backend');
    	setupIPCBridge();
    	hideOverlay();
    	clearInterval(connectTimer);
    	websocket.onclose = handleDisconnect;
    	websocket.onmessage = handleMessage;
    }

    // Handles websocket disconnects
    function handleDisconnect() {
    	log('Disconnected from backend');
    	websocket = null;
    	showOverlay();
    	connect();
    }

    // Try to connect to the backend every 1s (default value).
    function connect() {
    	connectTimer = setInterval(function () {
    		if (websocket == null) {
    			websocket = new WebSocket('ws://' + window.location.hostname + ':34115/bridge');
    			websocket.onopen = handleConnect;
    			websocket.onerror = function (e) {
    				e.stopImmediatePropagation();
    				e.stopPropagation();
    				e.preventDefault();
    				websocket = null;
    				return false;
    			};
    		}
    	}, 1000);
    }

    // Adds a script to the Dom.
    // Removes it if second parameter is true.
    function addScript(script, remove) {
    	const s = document.createElement('script');
    	s.setAttribute('type', 'text/javascript');
    	s.textContent = script;
    	document.head.appendChild(s);

    	// Remove internal messages from the DOM
    	if (remove) {
    		s.parentNode.removeChild(s);
    	}
    }

    function handleMessage(message) {
    	// As a bridge we ignore js and css injections
    	switch (message.data[0]) {
    	// Wails library - inject!
    	case 'b':
    		message = message.data.slice(1);
    		addScript(message);
    		log('Loaded Wails Runtime');

    		// We need to now send a message to the backend telling it
    		// we have loaded (System Start)
    		window.wailsInvoke('SS');
    		
    		// Now wails runtime is loaded, wails for the ready event
    		// and callback to the main app
    		// window.wails.Events.On('wails:loaded', function () {
    		if (callback) {
    			log('Notifying application');
    			callback(window.wails);
    		}
    		// });
    		break;
    		// Notifications
    	case 'n':
    		window.wails._.Notify(message.data.slice(1));
    		break;
    		// 	// Binding
    		// case 'b':
    		// 	const binding = message.data.slice(1);
    		// 	//log("Binding: " + binding)
    		// 	window.wails._.NewBinding(binding);
    		// 	break;
    		// 	// Call back
    	case 'c':
    		const callbackData = message.data.slice(1);
    		window.wails._.Callback(callbackData);
    		break;
    		// Tray
    	case 'T':
    		const trayMessage = message.data.slice(1);
    		switch (trayMessage[0]) {
    		case 'S':
    			// Set tray
    			const trayJSON = trayMessage.slice(1);
    			let tray = JSON.parse(trayJSON);
    			setTray(tray);
    			break;
    		case 'U':
    			// Update label
    			const updateTrayLabelJSON = trayMessage.slice(1);
    			let trayLabelData = JSON.parse(updateTrayLabelJSON);
    			updateTrayLabel(trayLabelData);
    			break;
    		case 'D':
    			// Delete Tray Menu
    			const id = trayMessage.slice(1);
    			deleteTrayMenu(id);
    			break;
    		default:
    			log('Unknown tray message: ' + message.data);
    		}
    		break;

    	default:
    		log('Unknown message: ' + message.data);
    	}
    }

    /*
     _       __      _ __
    | |     / /___ _(_) /____
    | | /| / / __ `/ / / ___/
    | |/ |/ / /_/ / / (__  )
    |__/|__/\__,_/_/_/____/
    The lightweight framework for web-like apps
    (c) Lea Anthony 2019-present
    */

    function setupMenuBar() {
    	new Menubar({
    		target: document.body,
    	});
    }

    // Sets up the overlay
    function setupOverlay() {
    	new Overlay({
    		target: document.body,
    		anchor: document.querySelector('#wails-bridge'),
    	});
    }

    function InitBridge(callback) {

    	setupMenuBar();

    	// Setup the overlay
    	setupOverlay();

    	// Start by showing the overlay...
    	showOverlay();

    	// ...and attempt to connect
    	StartWebsocket(callback);
    }

    exports.InitBridge = InitBridge;

    Object.defineProperty(exports, '__esModule', { value: true });

})));
