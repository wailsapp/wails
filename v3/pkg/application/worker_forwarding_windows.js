// WebView2 cannot start a dedicated Worker whose script is requested by another
// Worker from the asset server. The request does reach the server and is
// answered with the correct body and Content-Type, but the nested Worker still
// never starts and fires a silent error event (wails#6128). Scripts requested
// from the document, and blob-backed workers, are served normally.
//
// Injected into every attached Worker before its first script runs. It replaces
// globalThis.Worker so that a same-origin URL is fetched through a blob URL that
// only *imports* the original absolute URL, which WebView2 does serve. Inlining
// the source instead would change import.meta.url to the blob URL and break
// libraries such as Emscripten's pthread pool.
//
// The forward keeps the original script URL as the worker's base, so relative
// Worker(), importScripts(), fetch() and Request() URLs still resolve against
// the real script URL instead of the blob URL the forwarded worker runs from.
(() => {
    'use strict';

    if (globalThis.__wailsWorkerForwarding) return;
    const NativeWorker = globalThis.Worker;
    if (typeof NativeWorker !== 'function') return;
    globalThis.__wailsWorkerForwarding = true;

    const NativeFetch = globalThis.fetch;
    const NativeRequest = globalThis.Request;
    const NativeImportScripts = globalThis.importScripts;

    // The script URL this worker was created from. A forwarded worker runs from a
    // blob URL, so the generated entry records the original URL in
    // __wailsWorkerBase before the forwarded script is imported.
    const workerBase = () => {
        const forwarded = globalThis.__wailsWorkerBase;
        if (typeof forwarded === 'string' && forwarded !== '') return forwarded;
        const location = globalThis.location;
        return location && typeof location.href === 'string' ? location.href : undefined;
    };

    // Resolve a URL the way the worker's own script would, not against the blob
    // URL a forwarded worker actually runs from.
    const resolve = value => {
        const reference = workerBase();
        if (reference === undefined) return value;
        try { return new URL(value, reference).href; } catch (_) { return value; }
    };

    // Only asset-server requests from this worker's own origin are affected by
    // the WebView2 limitation. Cross-origin worker scripts keep their normal
    // CORS-checked handling.
    const sameOrigin = value => {
        const reference = workerBase();
        if (reference === undefined) return false;
        try { return new URL(value).origin === new URL(reference).origin; } catch (_) { return false; }
    };

    // Only the Wails asset-server origin is affected by the WebView2
    // limitation, so only its worker URLs need forwarding. A worker attached on
    // any other HTTP(S) origin keeps the native constructor: its requests go
    // through the normal network stack and a worker-src 'self' policy would
    // reject the replacement blob worker that the original URL satisfies.
    // This mirrors the asset host bound into the Fetch interception patterns.
    const affectedOrigin = value => {
        try { return new URL(value).hostname === 'wails.localhost'; } catch (_) { return false; }
    };

    const forward = (target, module) =>
        'globalThis.__wailsWorkerBase = ' + JSON.stringify(target) + ';\n' +
        (module
            // A static import is hoisted and would evaluate the target before the
            // base is recorded, so import it dynamically once the base is set.
            // The await is top-level on purpose: if the target fails to load, the
            // entry module's evaluation rejects and WebView2 fires the parent's
            // Worker error event. A discarded promise would surface the failure
            // only as a silent unhandled rejection in the child instead.
            ? 'await import(' + JSON.stringify(target) + ');\n'
            // importScripts is synchronous, so a failed load already throws and
            // fires the parent's Worker error event.
            : 'importScripts(' + JSON.stringify(target) + ');\n');

    globalThis.Worker = class extends NativeWorker {
        constructor(url, options) {
            const settings = options || {};
            const script = typeof url === 'string' ? url : url && typeof url.href === 'string' ? url.href : null;
            if (script === null) {
                super(url, settings);
                return;
            }
            const target = resolve(script);
            const forwardable = /^https?:/i.test(target) && sameOrigin(target) && affectedOrigin(target);
            if (!forwardable) {
                super(target, settings);
                return;
            }
            const objectURL = URL.createObjectURL(new Blob([forward(target, settings.type === 'module')], { type: 'text/javascript' }));
            try {
                super(objectURL, settings);
            } catch (error) {
                // The worker never started, so the blob is no longer needed.
                URL.revokeObjectURL(objectURL);
                throw error;
            }
            this.__wailsObjectURL = objectURL;
            // A script that fails to load can report its error after the
            // constructor returns, so the catch above would not run. Release the
            // blob when the child reports that error instead of retaining it for
            // the parent's lifetime. A separate listener never clobbers the
            // caller's own onerror handler.
            this.addEventListener('error', () => {
                const pending = this.__wailsObjectURL;
                if (pending) {
                    this.__wailsObjectURL = undefined;
                    URL.revokeObjectURL(pending);
                }
            }, { once: true });
        }

        terminate() {
            // The child can no longer fetch its entry script once terminated, so
            // release the blob instead of retaining it for the parent's lifetime.
            const objectURL = this.__wailsObjectURL;
            if (objectURL) {
                this.__wailsObjectURL = undefined;
                URL.revokeObjectURL(objectURL);
            }
            return super.terminate();
        }
    };

    // A classic worker's own importScripts("./helper.js") would otherwise resolve
    // against the blob entry URL. Resolve the arguments against the original
    // script URL first; absolute and cross-origin arguments are unchanged.
    if (typeof NativeImportScripts === 'function') {
        globalThis.importScripts = function () {
            return Reflect.apply(NativeImportScripts, this, Array.prototype.map.call(arguments, resolve));
        };
    }

    if (typeof NativeFetch === 'function') {
        globalThis.fetch = function (input, init) {
            if (typeof input === 'string') input = resolve(input);
            return Reflect.apply(NativeFetch, this, [input, init]);
        };
    }

    if (typeof NativeRequest === 'function') {
        globalThis.Request = class extends NativeRequest {
            constructor(input, init) {
                super(typeof input === 'string' ? resolve(input) : input, init);
            }
        };
    }
})();
