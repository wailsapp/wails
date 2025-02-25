import * as util from "node:util";
import { describe, it, beforeEach, afterEach, assert, expect, vi } from "vitest";
import { CancelError, CancellablePromise, CancelledRejectionError } from "./cancellable";

// TODO: In order of importance:
// TODO: test cancellation of subpromises the main promise resolves to.
// TODO: test cancellation of promise chains built by calling then() and friends:
//   - all promises up the chain should be cancelled;
//   - rejection handlers should be always executed with the CancelError of their parent promise in the chain;
//   - promises returned from rejection handlers should be cancelled too;
//   - if a rejection handler throws or returns a promise that ultimately rejects,
//     it should be reported as an unhandled rejection,
//   - unless it is a CancelError with the same reason given for cancelling the returned promise.
// TODO: test multiple calls to cancel() (second and later should have no effect).

let expectedUnhandled = new Map();

process.on('unhandledRejection', function (error, promise) {
    let reason = error;
    if (reason instanceof CancelledRejectionError) {
        promise = reason.promise;
        reason = reason.cause;
    }

    let reasons = expectedUnhandled.get(promise);
    const callbacks = reasons?.get(reason);
    if (callbacks) {
        for (const cb of callbacks) {
            try {
                cb(reason, promise);
            } catch (e) {
                console.error("Exception in unhandled rejection callback.", e);
            }
        }

        reasons.delete(reason);
        if (reasons.size === 0) {
            expectedUnhandled.delete(promise);
        }
        return;
    }

    console.log(util.format("Unhandled rejection.\nReason: %o\nPromise: %o", reason, promise));
    throw error;
});

function ignoreUnhandled(reason, promise) {
    expectUnhandled(reason, promise, null);
}

function expectUnhandled(reason, promise, cb) {
    let reasons = expectedUnhandled.get(promise);
    if (!reasons) {
        reasons = new Map();
        expectedUnhandled.set(promise, reasons);
    }
    let callbacks = reasons.get(reason);
    if (!callbacks) {
        callbacks = [];
        reasons.set(reason, callbacks);
    }
    if (cb) {
        callbacks.push(cb);
    }
}

afterEach(() => {
    vi.resetAllMocks();
    vi.restoreAllMocks();
});

const dummyValue = { value: "value" };
const dummyCause = { dummy: "dummy" };
const dummyError = new Error("dummy");
const oncancelled = vi.fn().mockName("oncancelled");
const sentinel = vi.fn().mockName("sentinel");
const unhandled = vi.fn().mockName("unhandled");

const resolutionPatterns = [
    ["forever", "pending", (test, value, { cls = CancellablePromise, cb = oncancelled } = {}) => test(
        new cls(() => {}, cb)
    )],
    ["already", "fulfilled", (test, value, { cls = CancellablePromise, cb = oncancelled } = {}) => {
        const prw = cls.withResolvers();
        prw.oncancelled = cb;
        prw.resolve(value ?? dummyValue);
        return test(prw.promise);
    }],
    ["immediately", "fulfilled", (test, value, { cls = CancellablePromise, cb = oncancelled } = {}) => {
        const prw = cls.withResolvers();
        prw.oncancelled = cb;
        const tp = test(prw.promise);
        prw.resolve(value ?? dummyValue);
        return tp;
    }],
    ["eventually", "fulfilled", async (test, value, { cls = CancellablePromise, cb = oncancelled } = {}) => {
        const prw = cls.withResolvers();
        prw.oncancelled = cb;
        const tp = test(prw.promise);
        await new Promise((resolve) => {
            setTimeout(() => {
                prw.resolve(value ?? dummyValue);
                resolve();
            }, 50);
        });
        return tp;
    }],
    ["already", "rejected", (test, reason, { cls = CancellablePromise, cb = oncancelled } = {}) => {
        const prw = cls.withResolvers();
        prw.oncancelled = cb;
        prw.reject(reason ?? dummyError);
        return test(prw.promise);
    }],
    ["immediately", "rejected", (test, reason, { cls = CancellablePromise, cb = oncancelled } = {}) => {
        const prw = cls.withResolvers();
        prw.oncancelled = cb;
        const tp = test(prw.promise);
        prw.reject(reason ?? dummyError);
        return tp;
    }],
    ["eventually", "rejected", async (test, reason, { cls = CancellablePromise, cb = oncancelled } = {}) => {
        const prw = cls.withResolvers();
        prw.oncancelled = cb;
        const tp = test(prw.promise);
        await new Promise((resolve) => {
            setTimeout(() => {
                prw.reject(reason ?? dummyError);
                resolve();
            }, 50);
        });
        return tp;
    }],
];

describe("CancellablePromise.cancel", ()=> {
    it("should suppress its own unhandled cancellation error", async () => {
        const p = new CancellablePromise(() => {});
        p.cancel();

        process.on('unhandledRejection', sentinel);
        await new Promise((resolve) => setTimeout(resolve, 100));
        process.off('unhandledRejection', sentinel);

        expect(sentinel).not.toHaveBeenCalled();
    });

    it.for([
        ["rejections", dummyError],
        ["cancellation errors", new CancelError("dummy", { cause: dummyCause })],
    ])("should not suppress arbitrary unhandled %s", async ([kind, err]) => {
        const p = new CancellablePromise(() => { throw err; });
        p.cancel();

        await new Promise((resolve) => {
            expectUnhandled(err, p, unhandled);
            expectUnhandled(err, p, resolve);
        });

        expect(unhandled).toHaveBeenCalledExactlyOnceWith(err, p);
    });

    describe.for(resolutionPatterns)("when applied to %s %s promises", ([time, state, test]) => {
        if (time === "already") {
            it("should have no effect", () => test(async (promise) => {
                promise.then(sentinel, sentinel);

                let reason;
                try {
                    promise.cancel();
                    await promise;
                    assert(state === "fulfilled", "Promise fulfilled unexpectedly");
                } catch (err) {
                    reason = err;
                    assert(state === "rejected", "Promise rejected unexpectedly");
                }

                expect(sentinel).toHaveBeenCalled();
                expect(oncancelled).not.toHaveBeenCalled();
                expect(reason).not.toBeInstanceOf(CancelError);
            }));
        } else {
            if (state === "rejected") {
                it("should report late rejections as unhandled", () => test(async (promise) => {
                    promise.cancel();

                    await new Promise((resolve) => {
                        expectUnhandled(dummyError, promise, unhandled);
                        expectUnhandled(dummyError, promise, resolve);
                    });

                    expect(unhandled).toHaveBeenCalledExactlyOnceWith(dummyError, promise);
                }));
            }

            it("should reject with a CancelError", () => test(async (promise) => {
                // Ignore the unhandled rejection from the test promise.
                if (state === "rejected") { ignoreUnhandled(dummyError, promise); }

                let reason;
                try {
                    promise.cancel();
                    await promise;
                } catch (err) {
                    reason = err;
                }

                expect(reason).toBeInstanceOf(CancelError);
            }));

            it("should call the oncancelled callback synchronously", () => test(async (promise) => {
                // Ignore the unhandled rejection from the test promise.
                if (state === "rejected") { ignoreUnhandled(dummyError, promise); }

                try {
                    promise.cancel();
                    sentinel();
                    await promise;
                } catch {}

                expect(oncancelled).toHaveBeenCalledBefore(sentinel);
            }));

            it("should propagate the given cause", () => test(async (promise) => {
                // Ignore the unhandled rejection from the test promise.
                if (state === "rejected") { ignoreUnhandled(dummyError, promise); }

                let reason;
                try {
                    promise.cancel(dummyCause);
                    await promise;
                } catch (err) {
                    reason = err;
                }

                expect(reason).toBeInstanceOf(CancelError);
                expect(reason).toHaveProperty('cause', dummyCause);
                expect(oncancelled).toHaveBeenCalledWith(reason.cause);
            }));
        }
    });
});

const onabort = vi.fn().mockName("abort");

const abortPatterns = [
    ["never", "standalone", (test) => {
        const signal = new AbortSignal();
        signal.addEventListener('abort', onabort, { capture: true });
        return test(signal);
    }],
    ["already", "standalone", (test) => {
        const signal = AbortSignal.abort(dummyCause);
        onabort();
        return test(signal);
    }],
    ["eventually", "standalone", (test) => {
        const signal = AbortSignal.timeout(25);
        signal.addEventListener('abort', onabort, { capture: true });
        return test(signal);
    }],
    ["never", "controller-bound", (test) => {
        const signal = new AbortController().signal;
        signal.addEventListener('abort', onabort, { capture: true });
        return test(signal);
    }],
    ["already", " controller-bound", (test) => {
        const ctrl = new AbortController();
        ctrl.signal.addEventListener('abort', onabort, { capture: true });
        ctrl.abort(dummyCause);
        return test(ctrl.signal);
    }],
    ["immediately", "controller-bound", (test) => {
        const ctrl = new AbortController();
        ctrl.signal.addEventListener('abort', onabort, { capture: true });
        const tp = test(ctrl.signal);
        ctrl.abort(dummyCause);
        return tp;
    }],
    ["eventually", "controller-bound", (test) => {
        const ctrl = new AbortController();
        ctrl.signal.addEventListener('abort', onabort, { capture: true });
        const tp = test(ctrl.signal);
        setTimeout(() => ctrl.abort(dummyCause), 25);
        return tp;
    }]
];

describe("CancellablePromise.cancelOn", ()=> {
    it("should return the target promise for chaining", () => {
        const p = new CancellablePromise(() => {});
        expect(p.cancelOn(AbortSignal.abort())).toBe(p);
    });

    function tests(abortTime, mode, testSignal, resolveTime, state, testPromise) {
        if (abortTime !== "never") {
            it(`should call CancellablePromise.cancel ${abortTime === "already" ? "immediately" : "on abort"} with the abort reason as cause`, () => testSignal((signal) => testPromise(async (promise) => {
                // Ignore the unhandled rejection from the test promise.
                if (state === "rejected") { ignoreUnhandled(dummyError, promise); }

                const cancelSpy = vi.spyOn(promise, 'cancel');

                promise.catch(() => {});
                promise.cancelOn(signal);

                if (signal.aborted) {
                    sentinel();
                } else {
                    await new Promise((resolve) => {
                        signal.onabort = () => {
                            sentinel();
                            resolve();
                        };
                    });
                }

                expect(cancelSpy).toHaveBeenCalledAfter(onabort);
                expect(cancelSpy).toHaveBeenCalledBefore(sentinel);
                expect(cancelSpy).toHaveBeenCalledExactlyOnceWith(signal.reason);
            })));
        }

        if (
            resolveTime === "already"
            || abortTime === "never"
            || (
                ["immediately", "eventually"].includes(abortTime)
                && ["already", "immediately"].includes(resolveTime)
            )
        ) {
            it("should have no effect", () => testSignal((signal) => testPromise(async (promise) => {
                promise.then(sentinel, sentinel);

                let reason;
                try {
                    if (resolveTime !== "forever") {
                        await promise.cancelOn(signal);
                        assert(state === "fulfilled", "Promise fulfilled unexpectedly");
                    } else {
                        await Promise.race([promise, new Promise((resolve) => setTimeout(resolve, 100))]).then(sentinel);
                    }
                } catch (err) {
                    reason = err;
                    assert(state === "rejected", "Promise rejected unexpectedly");
                }

                if (abortTime !== "never" && !signal.aborted) {
                    // Wait for the AbortSignal to have actually aborted.
                    await new Promise((resolve) => signal.onabort = resolve);
                }

                expect(sentinel).toHaveBeenCalled();
                expect(oncancelled).not.toHaveBeenCalled();
                expect(reason).not.toBeInstanceOf(CancelError);
            })));
        } else {
            if (state === "rejected") {
                it("should report late rejections as unhandled", () => testSignal((signal) => testPromise(async (promise) => {
                    promise.cancelOn(signal);

                    await new Promise((resolve) => {
                        expectUnhandled(dummyError, promise, unhandled);
                        expectUnhandled(dummyError, promise, resolve);
                    });

                    expect(unhandled).toHaveBeenCalledExactlyOnceWith(dummyError, promise);
                })));
            }

            it("should reject with a CancelError", () => testSignal((signal) => testPromise(async (promise)=> {
                // Ignore the unhandled rejection from the test promise.
                if (state === "rejected") { ignoreUnhandled(dummyError, promise); }

                let reason;
                try {
                    await promise.cancelOn(signal);
                } catch (err) {
                    reason = err;
                }

                expect(reason).toBeInstanceOf(CancelError);
            })));

            it(`should call the oncancelled callback ${abortTime === "already" ? "" : "a"}synchronously`, () => testSignal((signal) => testPromise(async (promise) => {
                // Ignore the unhandled rejection from the test promise.
                if (state === "rejected") { ignoreUnhandled(dummyError, promise); }

                try {
                    promise.cancelOn(signal);
                    sentinel();
                    await promise;
                } catch {}

                expect(oncancelled).toHaveBeenCalledAfter(onabort);
                if (abortTime === "already") {
                    expect(oncancelled).toHaveBeenCalledBefore(sentinel);
                } else {
                    expect(oncancelled).toHaveBeenCalledAfter(sentinel);
                }
            })));

            it("should propagate the abort reason as cause", () => testSignal((signal) => testPromise(async (promise) => {
                // Ignore the unhandled rejection from the test promise.
                if (state === "rejected") { ignoreUnhandled(dummyError, promise); }

                let reason;
                try {
                    await promise.cancelOn(signal);
                } catch (err) {
                    reason = err;
                }

                expect(reason).toBeInstanceOf(CancelError);
                expect(reason).toHaveProperty('cause', signal.reason);
                expect(oncancelled).toHaveBeenCalledWith(signal.reason);
            })));
        }
    }

    describe.for(abortPatterns)("when called with %s aborted %s signals", ([abortTime, mode, testSignal]) => {
        describe.for(resolutionPatterns)("when applied to %s %s promises", ([resolveTime, state, testPromise]) => {
            tests(abortTime, mode, testSignal, resolveTime, state, testPromise);
        });
    });

    describe.for(resolutionPatterns)("when applied to %s %s promises", ([resolveTime, state, testPromise]) => {
        describe.for(abortPatterns)("when called with %s aborted %s signals", ([abortTime, mode, testSignal]) => {
            tests(abortTime, mode, testSignal, resolveTime, state, testPromise);
        });
    });
});