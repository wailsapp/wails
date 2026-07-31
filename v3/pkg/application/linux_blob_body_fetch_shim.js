(function () {
    "use strict";

    if (globalThis.__wailsLinuxBlobBodyFetchShim) return;
    globalThis.__wailsLinuxBlobBodyFetchShim = true;

    var originalFetch = globalThis.fetch;
    if (typeof originalFetch !== "function") return;

    function targetsWailsScheme(input) {
        try {
            var value = input instanceof Request ? input.url : String(input);
            return new URL(value, globalThis.location.href).protocol === "wails:";
        } catch (_) {
            return false;
        }
    }

    globalThis.fetch = async function (input, init) {
        if (!targetsWailsScheme(input)) {
            return originalFetch.apply(this, arguments);
        }

        var request = input instanceof Request ? input : null;
        var hasInitBody = init && ("body" in init);
        var body = hasInitBody ? init.body : null;
        var headers;
        var bytes;

        if (body instanceof Blob) {
            headers = new Headers(init.headers);
            if (body.type && !headers.has("Content-Type")) {
                headers.set("Content-Type", body.type);
            }
            bytes = await body.arrayBuffer();
        } else if (body instanceof FormData) {
            var encoded = new Response(body);
            headers = new Headers(init.headers);
            if (!headers.has("Content-Type")) {
                headers.set("Content-Type", encoded.headers.get("Content-Type"));
            }
            bytes = await encoded.arrayBuffer();
        } else if (request && !hasInitBody && request.body) {
            var requestInit = Object.assign({}, init, {
                body: await request.clone().arrayBuffer()
            });
            var convertedRequest = new Request(request, requestInit);
            return originalFetch.call(this, convertedRequest);
        } else {
            return originalFetch.apply(this, arguments);
        }

        var converted = Object.assign({}, init, {body: bytes, headers: headers});
        return originalFetch.call(this, input, converted);
    };
})();
