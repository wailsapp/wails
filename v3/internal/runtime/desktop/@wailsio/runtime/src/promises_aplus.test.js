import * as util from "util";
import * as V from "vitest";
import { CancellablePromise } from "./cancellable";

// The Promises/A+ suite handles some errors late.
process.on('rejectionHandled', function () {});

// The Promises/A+ suite leaves some errors unhandled.
process.on('unhandledRejection', function (reason, promise) {
    if (promise instanceof CancellablePromise && reason != null && typeof reason === 'object') {
        for (const key of ['dummy', 'other', 'sentinel']) {
            if (reason[key] === key) {
                return;
            }
        }
    }
    throw new Error(`Unhandled rejection at: ${util.inspect(promise)}; reason: ${util.inspect(reason)}`, { cause: reason });
});

// Emulate a minimal version of the mocha BDD API using vitest primitives.
global.context = global.describe = V.describe;
global.specify = global.it = function it(desc, fn) {
    let viTestFn = fn;
    if (fn && fn.length) {
        viTestFn = () => new Promise((done) => fn(done));
    }
    V.it(desc, viTestFn);
}
global.before = function(desc, fn) { V.beforeAll(typeof desc === 'function' ? desc : fn) };
global.after = function(desc, fn) { V.afterAll(typeof desc === 'function' ? desc : fn) };
global.beforeEach = function(desc, fn) { V.beforeEach(typeof desc === 'function' ? desc : fn) };
global.afterEach = function(desc, fn) { V.afterEach(typeof desc === 'function' ? desc : fn) };

require('promises-aplus-tests').mocha({
    resolved(value) {
        return CancellablePromise.resolve(value);
    },
    rejected(reason) {
        return CancellablePromise.reject(reason);
    },
    deferred() {
        return CancellablePromise.withResolvers();
    }
});
