// Behavioural test for the Windows worker-forwarding shim. Run directly with
// `node worker_forwarding_windows_test.mjs`, or through
// TestWorkerForwardingScriptBehaviour, which skips when node is unavailable.
//
// The shim is a plain script injected into a worker's global scope, so each
// scenario evaluates it in an isolated vm context with a stub Worker/fetch so
// the forwarding decisions can be observed without a browser.

import assert from "node:assert/strict";
import fs from "node:fs";
import vm from "node:vm";
import { fileURLToPath } from "node:url";

const shimPath = fileURLToPath(new URL("./worker_forwarding_windows.js", import.meta.url));
const shimSource = fs.readFileSync(shimPath, "utf8");

// URL whose createObjectURL records the blob without needing a real blob store,
// and whose revokeObjectURL records every released object URL.
function recordingURL(blobs, revoked) {
    const construct = function (...args) {
        return new URL(...args);
    };
    construct.prototype = URL.prototype;
    construct.createObjectURL = blob => {
        blobs.push(blob);
        return `blob:http://wails.localhost/${blobs.length}`;
    };
    construct.revokeObjectURL = url => {
        if (revoked) revoked.push(url);
    };
    return construct;
}

function bootstrap({ href, forwarded, blobs, fetches, revoked, importScripts, createWorker } = {}) {
    const created = [];
    class NativeWorker {
        constructor(url, options) {
            if (createWorker) createWorker(String(url));
            created.push({ url: String(url), options: options || {} });
        }
        terminate() { this.terminated = true; }
        addEventListener(type, listener, options) {
            (this._listeners ??= {})[type] ??= [];
            this._listeners[type].push({ listener, once: !!(options && options.once) });
        }
        dispatchEvent(event) {
            for (const entry of [...((this._listeners && this._listeners[event.type]) || [])]) {
                entry.listener.call(this, event);
                if (entry.once) {
                    this._listeners[event.type] = this._listeners[event.type].filter(candidate => candidate !== entry);
                }
            }
            return true;
        }
    }
    const sandbox = {
        URL: blobs ? recordingURL(blobs, revoked) : URL,
        Blob,
        Request,
        fetch: function (input, init) {
            if (fetches) fetches.push(input);
            return Promise.resolve(new Response("ok"));
        },
        Worker: NativeWorker,
        location: { href },
        console,
    };
    if (importScripts) sandbox.importScripts = importScripts;
    if (forwarded !== undefined) sandbox.__wailsWorkerBase = forwarded;
    vm.createContext(sandbox);
    vm.runInContext(shimSource, sandbox, { filename: shimPath });
    sandbox.__created = created;
    return sandbox;
}

const blobText = blob => blob.text();

// A top-level worker forwards a relative URL through a module blob.
{
    const blobs = [];
    const box = bootstrap({ href: "http://wails.localhost/parent.js", blobs });
    new box.Worker("./child.js", { type: "module" });
    assert.equal(box.__created.length, 1);
    assert.match(box.__created[0].url, /^blob:/, "relative URL must be forwarded");
    assert.equal(box.__created[0].options.type, "module");
    const source = await blobText(blobs[0]);
    assert.match(source, /__wailsWorkerBase = "http:\/\/wails\.localhost\/child\.js"/);
    assert.ok(
        source.indexOf("__wailsWorkerBase") < source.indexOf("await import("),
        "the base must be recorded before the dynamic import",
    );
    assert.match(source, /await import\("http:\/\/wails\.localhost\/child\.js"\)/);
    assert.doesNotMatch(source, /\(\s*async/, "an async wrapper would swallow load failures");
    assert.doesNotMatch(source, /^\s*import\s/m, "a static import would run before the base is recorded");
}

// A classic worker forwards through importScripts.
{
    const blobs = [];
    const box = bootstrap({ href: "http://wails.localhost/parent.js", blobs });
    new box.Worker("/child.js");
    const source = await blobText(blobs[0]);
    assert.ok(
        source.indexOf("__wailsWorkerBase") < source.indexOf("importScripts("),
        "the base must be recorded before importScripts",
    );
    assert.match(source, /importScripts\("http:\/\/wails\.localhost\/child\.js"\)/);
    assert.equal(box.__created[0].options.type, undefined);
}

// A URL object is forwarded as well, so Emscripten's
// new Worker(new URL(..., import.meta.url)) path is covered.
{
    const blobs = [];
    const box = bootstrap({ href: "http://wails.localhost/parent.js", blobs });
    new box.Worker(new URL("http://wails.localhost/child.js"));
    assert.match(box.__created[0].url, /^blob:/);
}

// blob: and cross-origin URLs are left untouched.
{
    const blobs = [];
    const box = bootstrap({ href: "http://wails.localhost/parent.js", blobs });
    new box.Worker("blob:http://wails.localhost/existing");
    new box.Worker("https://example.com/worker.js");
    assert.equal(box.__created[0].url, "blob:http://wails.localhost/existing");
    assert.equal(box.__created[1].url, "https://example.com/worker.js");
    assert.equal(blobs.length, 0, "unsupported URLs must not be forwarded");
}

// A worker attached on another HTTP(S) origin keeps the native constructor even
// for a same-origin nested URL: only the asset-server origin is affected, and a
// worker-src 'self' policy must never see the blob replacement.
{
    const blobs = [];
    const box = bootstrap({ href: "https://example.com/parent.js", blobs });
    new box.Worker("./child.js");
    assert.equal(box.__created[0].url, "https://example.com/child.js", "a non-asset origin must not be forwarded");
    assert.equal(blobs.length, 0, "a non-asset origin must not create a blob");
}

// A forwarded worker resolves relative URLs against the original script URL.
{
    const blobs = [];
    const box = bootstrap({
        href: "blob:http://wails.localhost/1",
        forwarded: "http://wails.localhost/parent.js",
        blobs,
    });
    new box.Worker("./nested.js");
    assert.match(box.__created[0].url, /^blob:/, "nested workers are forwarded too");
    assert.match(await blobText(blobs[0]), /importScripts\("http:\/\/wails\.localhost\/nested\.js"\)/);
}

// A classic worker's importScripts resolves relative arguments against the
// original script URL while absolute and cross-origin arguments are unchanged.
{
    const calls = [];
    const native = (...urls) => {
        calls.push(urls);
        return undefined;
    };
    const box = bootstrap({
        href: "blob:http://wails.localhost/1",
        forwarded: "http://wails.localhost/workers/parent.js",
        importScripts: native,
    });
    box.importScripts("./helper.js", "https://example.com/lib.js", "/root.js");
    assert.deepEqual(calls, [[
        "http://wails.localhost/workers/helper.js",
        "https://example.com/lib.js",
        "http://wails.localhost/root.js",
    ]]);
}

// A top-level worker resolves against its own script URL, exactly as the
// native fetch would.
{
    const fetches = [];
    const box = bootstrap({ href: "http://wails.localhost/parent.js", fetches });
    await box.fetch("./data.json");
    assert.deepEqual(fetches, ["http://wails.localhost/data.json"]);
}

// Relative fetch() and Request() URLs resolve against the original script URL.
{
    const fetches = [];
    const box = bootstrap({
        href: "blob:http://wails.localhost/1",
        forwarded: "http://wails.localhost/workers/parent.js",
        fetches,
    });
    await box.fetch("./data.json");
    await box.fetch("/root.json");
    await box.fetch("https://example.com/x.json");
    assert.deepEqual(fetches, [
        "http://wails.localhost/workers/data.json",
        "http://wails.localhost/root.json",
        "https://example.com/x.json",
    ]);
    assert.equal(new box.Request("../shared.json").url, "http://wails.localhost/shared.json");
    assert.equal(new box.Request("https://example.com/y.json").url, "https://example.com/y.json");
}

// A forwarded worker's blob URL is released when the child is terminated.
{
    const blobs = [];
    const revoked = [];
    const box = bootstrap({ href: "http://wails.localhost/parent.js", blobs, revoked });
    const worker = new box.Worker("./child.js");
    const objectURL = worker.__wailsObjectURL;
    assert.match(objectURL, /^blob:/);
    worker.terminate();
    assert.deepEqual(revoked, [objectURL], "terminate must release the object URL");
    worker.terminate();
    assert.deepEqual(revoked, [objectURL], "a second terminate must not double release");
    assert.equal(box.__created[0].url, objectURL);
}

// A forwarded worker whose construction fails releases its object URL.
{
    const blobs = [];
    const revoked = [];
    const box = bootstrap({
        href: "http://wails.localhost/parent.js",
        blobs,
        revoked,
        createWorker: () => {
            throw new Error("construction failed");
        },
    });
    assert.throws(() => new box.Worker("./child.js"), /construction failed/);
    assert.equal(blobs.length, 1);
    assert.deepEqual(revoked, ["blob:http://wails.localhost/1"]);
}

// A forwarded worker that reports a startup error after the constructor
// returned releases its object URL, without relying on terminate().
{
    const blobs = [];
    const revoked = [];
    const box = bootstrap({ href: "http://wails.localhost/parent.js", blobs, revoked });
    const worker = new box.Worker("./child.js");
    const objectURL = worker.__wailsObjectURL;
    assert.match(objectURL, /^blob:/);
    worker.dispatchEvent({ type: "error" });
    assert.deepEqual(revoked, [objectURL], "a startup error must release the object URL");
    worker.terminate();
    assert.deepEqual(revoked, [objectURL], "terminate after a startup error must not double release");
}

// Installing twice must not double-wrap (the marker is checked and set).
{
    const box = bootstrap({ href: "http://wails.localhost/parent.js" });
    assert.equal(box.__wailsWorkerForwarding, true);
    const before = box.Worker;
    vm.runInContext(shimSource, box, { filename: shimPath });
    assert.equal(box.Worker, before, "a second injection must be ignored");
}

console.log("worker forwarding shim: ok");
