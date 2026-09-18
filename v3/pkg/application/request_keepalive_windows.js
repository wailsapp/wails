// DevTools reports page teardown as ERR_ABORTED even for keepalive fetches.
// Preserve those requests until completion or explicit AbortSignal cancellation.
// This script runs before page/worker scripts and only marks local keepalive calls.
(() => {
    const notify = globalThis.__wailsAssetKeepalive;
    const nativeFetch = globalThis.fetch;
    if (typeof notify !== 'function' || typeof nativeFetch !== 'function') return;
    globalThis.fetch = function(input, init) {
        // Leave the usual fetch path untouched, including external requests.
        if (!(init && init.keepalive) && !(input instanceof Request && input.keepalive)) {
            return Reflect.apply(nativeFetch, this, arguments);
        }
        let request;
        try { request = new Request(input, init); }
        catch (error) { return Promise.reject(error); }
        const url = new URL(request.url);
        if (!request.keepalive || request.signal.aborted || url.protocol !== 'http:' || url.hostname !== 'wails.localhost') {
            return nativeFetch.call(this, request);
        }
        const id = crypto.randomUUID();
        const send = state => notify(JSON.stringify({ ID: id, State: state }));
        const abort = () => send('abort');
        url.hash += '&__wails_keepalive=' + id;
        send('start');
        request.signal.addEventListener('abort', abort, { once: true });
        const dispatch = async () => {
            const options = {
                method: request.method, headers: request.headers, mode: request.mode,
                credentials: request.credentials, cache: request.cache, redirect: request.redirect,
                referrer: request.referrer, referrerPolicy: request.referrerPolicy,
                integrity: request.integrity, keepalive: true, signal: request.signal,
            };
            // A Request body is exposed as a stream, which keepalive forbids as
            // input. Materialise it before rebuilding the URL; native fetch
            // still enforces its keepalive body-size limit.
            if (request.body !== null) options.body = await request.arrayBuffer();
            return nativeFetch.call(this, url, options);
        };
        return dispatch().then(response => {
            request.signal.removeEventListener('abort', abort);
            send('end');
            return response;
        }, error => {
            request.signal.removeEventListener('abort', abort);
            // WebView2 rejects the promise on navigation even though keepalive
            // permits an already dispatched Go handler to finish its work.
            send(request.signal.aborted ? 'abort' : 'failed');
            throw error;
        });
    };
})();
