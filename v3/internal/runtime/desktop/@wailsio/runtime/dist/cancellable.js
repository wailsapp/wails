/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/
var _a;
import isCallable from "./callable.js";
/**
 * Exception class that will be used as rejection reason
 * in case a {@link CancellablePromise} is cancelled successfully.
 *
 * The value of the {@link name} property is the string `"CancelError"`.
 * The value of the {@link cause} property is the cause passed to the cancel method, if any.
 */
export class CancelError extends Error {
    /**
     * Constructs a new `CancelError` instance.
     * @param message - The error message.
     * @param options - Options to be forwarded to the Error constructor.
     */
    constructor(message, options) {
        super(message, options);
        this.name = "CancelError";
    }
}
/**
 * Exception class that will be reported as an unhandled rejection
 * in case a {@link CancellablePromise} rejects after being cancelled,
 * or when the `oncancelled` callback throws or rejects.
 *
 * The value of the {@link name} property is the string `"CancelledRejectionError"`.
 * The value of the {@link cause} property is the reason the promise rejected with.
 *
 * Because the original promise was cancelled,
 * a wrapper promise will be passed to the unhandled rejection listener instead.
 * The {@link promise} property holds a reference to the original promise.
 */
export class CancelledRejectionError extends Error {
    /**
     * Constructs a new `CancelledRejectionError` instance.
     * @param promise - The promise that caused the error originally.
     * @param reason - The rejection reason.
     * @param info - An optional informative message specifying the circumstances in which the error was thrown.
     *               Defaults to the string `"Unhandled rejection in cancelled promise."`.
     */
    constructor(promise, reason, info) {
        super((info !== null && info !== void 0 ? info : "Unhandled rejection in cancelled promise.") + " Reason: " + errorMessage(reason), { cause: reason });
        this.promise = promise;
        this.name = "CancelledRejectionError";
    }
}
// Private field names.
const barrierSym = Symbol("barrier");
const cancelImplSym = Symbol("cancelImpl");
const species = (_a = Symbol.species) !== null && _a !== void 0 ? _a : Symbol("speciesPolyfill");
/**
 * A promise with an attached method for cancelling long-running operations (see {@link CancellablePromise#cancel}).
 * Cancellation can optionally be bound to an {@link AbortSignal}
 * for better composability (see {@link CancellablePromise#cancelOn}).
 *
 * Cancelling a pending promise will result in an immediate rejection
 * with an instance of {@link CancelError} as reason,
 * but whoever started the promise will be responsible
 * for actually aborting the underlying operation.
 * To this purpose, the constructor and all chaining methods
 * accept optional cancellation callbacks.
 *
 * If a `CancellablePromise` still resolves after having been cancelled,
 * the result will be discarded. If it rejects, the reason
 * will be reported as an unhandled rejection,
 * wrapped in a {@link CancelledRejectionError} instance.
 * To facilitate the handling of cancellation requests,
 * cancelled `CancellablePromise`s will _not_ report unhandled `CancelError`s
 * whose `cause` field is the same as the one with which the current promise was cancelled.
 *
 * All usual promise methods are defined and return a `CancellablePromise`
 * whose cancel method will cancel the parent operation as well, propagating the cancellation reason
 * upwards through promise chains.
 * Conversely, cancelling a promise will not automatically cancel dependent promises downstream:
 * ```ts
 * let root = new CancellablePromise((resolve, reject) => { ... });
 * let child1 = root.then(() => { ... });
 * let child2 = child1.then(() => { ... });
 * let child3 = root.catch(() => { ... });
 * child1.cancel(); // Cancels child1 and root, but not child2 or child3
 * ```
 * Cancelling a promise that has already settled is safe and has no consequence.
 *
 * The `cancel` method returns a promise that _always fulfills_
 * after the whole chain has processed the cancel request
 * and all attached callbacks up to that moment have run.
 *
 * All ES2024 promise methods (static and instance) are defined on CancellablePromise,
 * but actual availability may vary with OS/webview version.
 *
 * In line with the proposal at https://github.com/tc39/proposal-rm-builtin-subclassing,
 * `CancellablePromise` does not support transparent subclassing.
 * Extenders should take care to provide their own method implementations.
 * This might be reconsidered in case the proposal is retired.
 *
 * CancellablePromise is a wrapper around the DOM Promise object
 * and is compliant with the [Promises/A+ specification](https://promisesaplus.com/)
 * (it passes the [compliance suite](https://github.com/promises-aplus/promises-tests))
 * if so is the underlying implementation.
 */
export class CancellablePromise extends Promise {
    /**
     * Creates a new `CancellablePromise`.
     *
     * @param executor - A callback used to initialize the promise. This callback is passed two arguments:
     *                   a `resolve` callback used to resolve the promise with a value
     *                   or the result of another promise (possibly cancellable),
     *                   and a `reject` callback used to reject the promise with a provided reason or error.
     *                   If the value provided to the `resolve` callback is a thenable _and_ cancellable object
     *                   (it has a `then` _and_ a `cancel` method),
     *                   cancellation requests will be forwarded to that object and the oncancelled will not be invoked anymore.
     *                   If any one of the two callbacks is called _after_ the promise has been cancelled,
     *                   the provided values will be cancelled and resolved as usual,
     *                   but their results will be discarded.
     *                   However, if the resolution process ultimately ends up in a rejection
     *                   that is not due to cancellation, the rejection reason
     *                   will be wrapped in a {@link CancelledRejectionError}
     *                   and bubbled up as an unhandled rejection.
     * @param oncancelled - It is the caller's responsibility to ensure that any operation
     *                      started by the executor is properly halted upon cancellation.
     *                      This optional callback can be used to that purpose.
     *                      It will be called _synchronously_ with a cancellation cause
     *                      when cancellation is requested, _after_ the promise has already rejected
     *                      with a {@link CancelError}, but _before_
     *                      any {@link then}/{@link catch}/{@link finally} callback runs.
     *                      If the callback returns a thenable, the promise returned from {@link cancel}
     *                      will only fulfill after the former has settled.
     *                      Unhandled exceptions or rejections from the callback will be wrapped
     *                      in a {@link CancelledRejectionError} and bubbled up as unhandled rejections.
     *                      If the `resolve` callback is called before cancellation with a cancellable promise,
     *                      cancellation requests on this promise will be diverted to that promise,
     *                      and the original `oncancelled` callback will be discarded.
     */
    constructor(executor, oncancelled) {
        let resolve;
        let reject;
        super((res, rej) => { resolve = res; reject = rej; });
        if (this.constructor[species] !== Promise) {
            throw new TypeError("CancellablePromise does not support transparent subclassing. Please refrain from overriding the [Symbol.species] static property.");
        }
        let promise = {
            promise: this,
            resolve,
            reject,
            get oncancelled() { return oncancelled !== null && oncancelled !== void 0 ? oncancelled : null; },
            set oncancelled(cb) { oncancelled = cb !== null && cb !== void 0 ? cb : undefined; }
        };
        const state = {
            get root() { return state; },
            resolving: false,
            settled: false
        };
        // Setup cancellation system.
        void Object.defineProperties(this, {
            [barrierSym]: {
                configurable: false,
                enumerable: false,
                writable: true,
                value: null
            },
            [cancelImplSym]: {
                configurable: false,
                enumerable: false,
                writable: false,
                value: cancellerFor(promise, state)
            }
        });
        // Run the actual executor.
        const rejector = rejectorFor(promise, state);
        try {
            executor(resolverFor(promise, state), rejector);
        }
        catch (err) {
            if (state.resolving) {
                console.log("Unhandled exception in CancellablePromise executor.", err);
            }
            else {
                rejector(err);
            }
        }
    }
    /**
     * Cancels immediately the execution of the operation associated with this promise.
     * The promise rejects with a {@link CancelError} instance as reason,
     * with the {@link CancelError#cause} property set to the given argument, if any.
     *
     * Has no effect if called after the promise has already settled;
     * repeated calls in particular are safe, but only the first one
     * will set the cancellation cause.
     *
     * The `CancelError` exception _need not_ be handled explicitly _on the promises that are being cancelled:_
     * cancelling a promise with no attached rejection handler does not trigger an unhandled rejection event.
     * Therefore, the following idioms are all equally correct:
     * ```ts
     * new CancellablePromise((resolve, reject) => { ... }).cancel();
     * new CancellablePromise((resolve, reject) => { ... }).then(...).cancel();
     * new CancellablePromise((resolve, reject) => { ... }).then(...).catch(...).cancel();
     * ```
     * Whenever some cancelled promise in a chain rejects with a `CancelError`
     * with the same cancellation cause as itself, the error will be discarded silently.
     * However, the `CancelError` _will still be delivered_ to all attached rejection handlers
     * added by {@link then} and related methods:
     * ```ts
     * let cancellable = new CancellablePromise((resolve, reject) => { ... });
     * cancellable.then(() => { ... }).catch(console.log);
     * cancellable.cancel(); // A CancelError is printed to the console.
     * ```
     * If the `CancelError` is not handled downstream by the time it reaches
     * a _non-cancelled_ promise, it _will_ trigger an unhandled rejection event,
     * just like normal rejections would:
     * ```ts
     * let cancellable = new CancellablePromise((resolve, reject) => { ... });
     * let chained = cancellable.then(() => { ... }).then(() => { ... }); // No catch...
     * cancellable.cancel(); // Unhandled rejection event on chained!
     * ```
     * Therefore, it is important to either cancel whole promise chains from their tail,
     * as shown in the correct idioms above, or take care of handling errors everywhere.
     *
     * @returns A cancellable promise that _fulfills_ after the cancel callback (if any)
     * and all handlers attached up to the call to cancel have run.
     * If the cancel callback returns a thenable, the promise returned by `cancel`
     * will also wait for that thenable to settle.
     * This enables callers to wait for the cancelled operation to terminate
     * without being forced to handle potential errors at the call site.
     * ```ts
     * cancellable.cancel().then(() => {
     *     // Cleanup finished, it's safe to do something else.
     * }, (err) => {
     *     // Unreachable: the promise returned from cancel will never reject.
     * });
     * ```
     * Note that the returned promise will _not_ handle implicitly any rejection
     * that might have occurred already in the cancelled chain.
     * It will just track whether registered handlers have been executed or not.
     * Therefore, unhandled rejections will never be silently handled by calling cancel.
     */
    cancel(cause) {
        return new CancellablePromise((resolve) => {
            // INVARIANT: the result of this[cancelImplSym] and the barrier do not ever reject.
            // Unfortunately macOS High Sierra does not support Promise.allSettled.
            Promise.all([
                this[cancelImplSym](new CancelError("Promise cancelled.", { cause })),
                currentBarrier(this)
            ]).then(() => resolve(), () => resolve());
        });
    }
    /**
     * Binds promise cancellation to the abort event of the given {@link AbortSignal}.
     * If the signal has already aborted, the promise will be cancelled immediately.
     * When either condition is verified, the cancellation cause will be set
     * to the signal's abort reason (see {@link AbortSignal#reason}).
     *
     * Has no effect if called (or if the signal aborts) _after_ the promise has already settled.
     * Only the first signal to abort will set the cancellation cause.
     *
     * For more details about the cancellation process,
     * see {@link cancel} and the `CancellablePromise` constructor.
     *
     * This method enables `await`ing cancellable promises without having
     * to store them for future cancellation, e.g.:
     * ```ts
     * await longRunningOperation().cancelOn(signal);
     * ```
     * instead of:
     * ```ts
     * let promiseToBeCancelled = longRunningOperation();
     * await promiseToBeCancelled;
     * ```
     *
     * @returns This promise, for method chaining.
     */
    cancelOn(signal) {
        if (signal.aborted) {
            void this.cancel(signal.reason);
        }
        else {
            signal.addEventListener('abort', () => void this.cancel(signal.reason), { capture: true });
        }
        return this;
    }
    /**
     * Attaches callbacks for the resolution and/or rejection of the `CancellablePromise`.
     *
     * The optional `oncancelled` argument will be invoked when the returned promise is cancelled,
     * with the same semantics as the `oncancelled` argument of the constructor.
     * When the parent promise rejects or is cancelled, the `onrejected` callback will run,
     * _even after the returned promise has been cancelled:_
     * in that case, should it reject or throw, the reason will be wrapped
     * in a {@link CancelledRejectionError} and bubbled up as an unhandled rejection.
     *
     * @param onfulfilled The callback to execute when the Promise is resolved.
     * @param onrejected The callback to execute when the Promise is rejected.
     * @returns A `CancellablePromise` for the completion of whichever callback is executed.
     * The returned promise is hooked up to propagate cancellation requests up the chain, but not down:
     *
     *   - if the parent promise is cancelled, the `onrejected` handler will be invoked with a `CancelError`
     *     and the returned promise _will resolve regularly_ with its result;
     *   - conversely, if the returned promise is cancelled, _the parent promise is cancelled too;_
     *     the `onrejected` handler will still be invoked with the parent's `CancelError`,
     *     but its result will be discarded
     *     and the returned promise will reject with a `CancelError` as well.
     *
     * The promise returned from {@link cancel} will fulfill only after all attached handlers
     * up the entire promise chain have been run.
     *
     * If either callback returns a cancellable promise,
     * cancellation requests will be diverted to it,
     * and the specified `oncancelled` callback will be discarded.
     */
    then(onfulfilled, onrejected, oncancelled) {
        if (!(this instanceof CancellablePromise)) {
            throw new TypeError("CancellablePromise.prototype.then called on an invalid object.");
        }
        // NOTE: TypeScript's built-in type for then is broken,
        // as it allows specifying an arbitrary TResult1 != T even when onfulfilled is not a function.
        // We cannot fix it if we want to CancellablePromise to implement PromiseLike<T>.
        if (!isCallable(onfulfilled)) {
            onfulfilled = identity;
        }
        if (!isCallable(onrejected)) {
            onrejected = thrower;
        }
        if (onfulfilled === identity && onrejected == thrower) {
            // Shortcut for trivial arguments.
            return new CancellablePromise((resolve) => resolve(this));
        }
        const barrier = {};
        this[barrierSym] = barrier;
        return new CancellablePromise((resolve, reject) => {
            void super.then((value) => {
                var _a;
                if (this[barrierSym] === barrier) {
                    this[barrierSym] = null;
                }
                (_a = barrier.resolve) === null || _a === void 0 ? void 0 : _a.call(barrier);
                try {
                    resolve(onfulfilled(value));
                }
                catch (err) {
                    reject(err);
                }
            }, (reason) => {
                var _a;
                if (this[barrierSym] === barrier) {
                    this[barrierSym] = null;
                }
                (_a = barrier.resolve) === null || _a === void 0 ? void 0 : _a.call(barrier);
                try {
                    resolve(onrejected(reason));
                }
                catch (err) {
                    reject(err);
                }
            });
        }, async (cause) => {
            //cancelled = true;
            try {
                return oncancelled === null || oncancelled === void 0 ? void 0 : oncancelled(cause);
            }
            finally {
                await this.cancel(cause);
            }
        });
    }
    /**
     * Attaches a callback for only the rejection of the Promise.
     *
     * The optional `oncancelled` argument will be invoked when the returned promise is cancelled,
     * with the same semantics as the `oncancelled` argument of the constructor.
     * When the parent promise rejects or is cancelled, the `onrejected` callback will run,
     * _even after the returned promise has been cancelled:_
     * in that case, should it reject or throw, the reason will be wrapped
     * in a {@link CancelledRejectionError} and bubbled up as an unhandled rejection.
     *
     * It is equivalent to
     * ```ts
     * cancellablePromise.then(undefined, onrejected, oncancelled);
     * ```
     * and the same caveats apply.
     *
     * @returns A Promise for the completion of the callback.
     * Cancellation requests on the returned promise
     * will propagate up the chain to the parent promise,
     * but not in the other direction.
     *
     * The promise returned from {@link cancel} will fulfill only after all attached handlers
     * up the entire promise chain have been run.
     *
     * If `onrejected` returns a cancellable promise,
     * cancellation requests will be diverted to it,
     * and the specified `oncancelled` callback will be discarded.
     * See {@link then} for more details.
     */
    catch(onrejected, oncancelled) {
        return this.then(undefined, onrejected, oncancelled);
    }
    /**
     * Attaches a callback that is invoked when the CancellablePromise is settled (fulfilled or rejected). The
     * resolved value cannot be accessed or modified from the callback.
     * The returned promise will settle in the same state as the original one
     * after the provided callback has completed execution,
     * unless the callback throws or returns a rejecting promise,
     * in which case the returned promise will reject as well.
     *
     * The optional `oncancelled` argument will be invoked when the returned promise is cancelled,
     * with the same semantics as the `oncancelled` argument of the constructor.
     * Once the parent promise settles, the `onfinally` callback will run,
     * _even after the returned promise has been cancelled:_
     * in that case, should it reject or throw, the reason will be wrapped
     * in a {@link CancelledRejectionError} and bubbled up as an unhandled rejection.
     *
     * This method is implemented in terms of {@link then} and the same caveats apply.
     * It is polyfilled, hence available in every OS/webview version.
     *
     * @returns A Promise for the completion of the callback.
     * Cancellation requests on the returned promise
     * will propagate up the chain to the parent promise,
     * but not in the other direction.
     *
     * The promise returned from {@link cancel} will fulfill only after all attached handlers
     * up the entire promise chain have been run.
     *
     * If `onfinally` returns a cancellable promise,
     * cancellation requests will be diverted to it,
     * and the specified `oncancelled` callback will be discarded.
     * See {@link then} for more details.
     */
    finally(onfinally, oncancelled) {
        if (!(this instanceof CancellablePromise)) {
            throw new TypeError("CancellablePromise.prototype.finally called on an invalid object.");
        }
        if (!isCallable(onfinally)) {
            return this.then(onfinally, onfinally, oncancelled);
        }
        return this.then((value) => CancellablePromise.resolve(onfinally()).then(() => value), (reason) => CancellablePromise.resolve(onfinally()).then(() => { throw reason; }), oncancelled);
    }
    /**
     * We use the `[Symbol.species]` static property, if available,
     * to disable the built-in automatic subclassing features from {@link Promise}.
     * It is critical for performance reasons that extenders do not override this.
     * Once the proposal at https://github.com/tc39/proposal-rm-builtin-subclassing
     * is either accepted or retired, this implementation will have to be revised accordingly.
     *
     * @ignore
     * @internal
     */
    static get [species]() {
        return Promise;
    }
    static all(values) {
        let collected = Array.from(values);
        const promise = collected.length === 0
            ? CancellablePromise.resolve(collected)
            : new CancellablePromise((resolve, reject) => {
                void Promise.all(collected).then(resolve, reject);
            }, (cause) => cancelAll(promise, collected, cause));
        return promise;
    }
    static allSettled(values) {
        let collected = Array.from(values);
        const promise = collected.length === 0
            ? CancellablePromise.resolve(collected)
            : new CancellablePromise((resolve, reject) => {
                void Promise.allSettled(collected).then(resolve, reject);
            }, (cause) => cancelAll(promise, collected, cause));
        return promise;
    }
    static any(values) {
        let collected = Array.from(values);
        const promise = collected.length === 0
            ? CancellablePromise.resolve(collected)
            : new CancellablePromise((resolve, reject) => {
                void Promise.any(collected).then(resolve, reject);
            }, (cause) => cancelAll(promise, collected, cause));
        return promise;
    }
    static race(values) {
        let collected = Array.from(values);
        const promise = new CancellablePromise((resolve, reject) => {
            void Promise.race(collected).then(resolve, reject);
        }, (cause) => cancelAll(promise, collected, cause));
        return promise;
    }
    /**
     * Creates a new cancelled CancellablePromise for the provided cause.
     *
     * @group Static Methods
     */
    static cancel(cause) {
        const p = new CancellablePromise(() => { });
        p.cancel(cause);
        return p;
    }
    /**
     * Creates a new CancellablePromise that cancels
     * after the specified timeout, with the provided cause.
     *
     * If the {@link AbortSignal.timeout} factory method is available,
     * it is used to base the timeout on _active_ time rather than _elapsed_ time.
     * Otherwise, `timeout` falls back to {@link setTimeout}.
     *
     * @group Static Methods
     */
    static timeout(milliseconds, cause) {
        const promise = new CancellablePromise(() => { });
        if (AbortSignal && typeof AbortSignal === 'function' && AbortSignal.timeout && typeof AbortSignal.timeout === 'function') {
            AbortSignal.timeout(milliseconds).addEventListener('abort', () => void promise.cancel(cause));
        }
        else {
            setTimeout(() => void promise.cancel(cause), milliseconds);
        }
        return promise;
    }
    static sleep(milliseconds, value) {
        return new CancellablePromise((resolve) => {
            setTimeout(() => resolve(value), milliseconds);
        });
    }
    /**
     * Creates a new rejected CancellablePromise for the provided reason.
     *
     * @group Static Methods
     */
    static reject(reason) {
        return new CancellablePromise((_, reject) => reject(reason));
    }
    static resolve(value) {
        if (value instanceof CancellablePromise) {
            // Optimise for cancellable promises.
            return value;
        }
        return new CancellablePromise((resolve) => resolve(value));
    }
    /**
     * Creates a new CancellablePromise and returns it in an object, along with its resolve and reject functions
     * and a getter/setter for the cancellation callback.
     *
     * This method is polyfilled, hence available in every OS/webview version.
     *
     * @group Static Methods
     */
    static withResolvers() {
        let result = { oncancelled: null };
        result.promise = new CancellablePromise((resolve, reject) => {
            result.resolve = resolve;
            result.reject = reject;
        }, (cause) => { var _a; (_a = result.oncancelled) === null || _a === void 0 ? void 0 : _a.call(result, cause); });
        return result;
    }
}
/**
 * Returns a callback that implements the cancellation algorithm for the given cancellable promise.
 * The promise returned from the resulting function does not reject.
 */
function cancellerFor(promise, state) {
    let cancellationPromise = undefined;
    return (reason) => {
        if (!state.settled) {
            state.settled = true;
            state.reason = reason;
            promise.reject(reason);
            // Attach an error handler that ignores this specific rejection reason and nothing else.
            // In theory, a sane underlying implementation at this point
            // should always reject with our cancellation reason,
            // hence the handler will never throw.
            void Promise.prototype.then.call(promise.promise, undefined, (err) => {
                if (err !== reason) {
                    throw err;
                }
            });
        }
        // If reason is not set, the promise resolved regularly, hence we must not call oncancelled.
        // If oncancelled is unset, no need to go any further.
        if (!state.reason || !promise.oncancelled) {
            return;
        }
        cancellationPromise = new Promise((resolve) => {
            try {
                resolve(promise.oncancelled(state.reason.cause));
            }
            catch (err) {
                Promise.reject(new CancelledRejectionError(promise.promise, err, "Unhandled exception in oncancelled callback."));
            }
        }).catch((reason) => {
            Promise.reject(new CancelledRejectionError(promise.promise, reason, "Unhandled rejection in oncancelled callback."));
        });
        // Unset oncancelled to prevent repeated calls.
        promise.oncancelled = null;
        return cancellationPromise;
    };
}
/**
 * Returns a callback that implements the resolution algorithm for the given cancellable promise.
 */
function resolverFor(promise, state) {
    return (value) => {
        if (state.resolving) {
            return;
        }
        state.resolving = true;
        if (value === promise.promise) {
            if (state.settled) {
                return;
            }
            state.settled = true;
            promise.reject(new TypeError("A promise cannot be resolved with itself."));
            return;
        }
        if (value != null && (typeof value === 'object' || typeof value === 'function')) {
            let then;
            try {
                then = value.then;
            }
            catch (err) {
                state.settled = true;
                promise.reject(err);
                return;
            }
            if (isCallable(then)) {
                try {
                    let cancel = value.cancel;
                    if (isCallable(cancel)) {
                        const oncancelled = (cause) => {
                            Reflect.apply(cancel, value, [cause]);
                        };
                        if (state.reason) {
                            // If already cancelled, propagate cancellation.
                            // The promise returned from the canceller algorithm does not reject
                            // so it can be discarded safely.
                            void cancellerFor(Object.assign(Object.assign({}, promise), { oncancelled }), state)(state.reason);
                        }
                        else {
                            promise.oncancelled = oncancelled;
                        }
                    }
                }
                catch (_a) { }
                const newState = {
                    root: state.root,
                    resolving: false,
                    get settled() { return this.root.settled; },
                    set settled(value) { this.root.settled = value; },
                    get reason() { return this.root.reason; }
                };
                const rejector = rejectorFor(promise, newState);
                try {
                    Reflect.apply(then, value, [resolverFor(promise, newState), rejector]);
                }
                catch (err) {
                    rejector(err);
                }
                return; // IMPORTANT!
            }
        }
        if (state.settled) {
            return;
        }
        state.settled = true;
        promise.resolve(value);
    };
}
/**
 * Returns a callback that implements the rejection algorithm for the given cancellable promise.
 */
function rejectorFor(promise, state) {
    return (reason) => {
        if (state.resolving) {
            return;
        }
        state.resolving = true;
        if (state.settled) {
            try {
                if (reason instanceof CancelError && state.reason instanceof CancelError && Object.is(reason.cause, state.reason.cause)) {
                    // Swallow late rejections that are CancelErrors whose cancellation cause is the same as ours.
                    return;
                }
            }
            catch (_a) { }
            void Promise.reject(new CancelledRejectionError(promise.promise, reason));
        }
        else {
            state.settled = true;
            promise.reject(reason);
        }
    };
}
/**
 * Cancels all values in an array that look like cancellable thenables.
 * Returns a promise that fulfills once all cancellation procedures for the given values have settled.
 */
function cancelAll(parent, values, cause) {
    const results = [];
    for (const value of values) {
        let cancel;
        try {
            if (!isCallable(value.then)) {
                continue;
            }
            cancel = value.cancel;
            if (!isCallable(cancel)) {
                continue;
            }
        }
        catch (_a) {
            continue;
        }
        let result;
        try {
            result = Reflect.apply(cancel, value, [cause]);
        }
        catch (err) {
            Promise.reject(new CancelledRejectionError(parent, err, "Unhandled exception in cancel method."));
            continue;
        }
        if (!result) {
            continue;
        }
        results.push((result instanceof Promise ? result : Promise.resolve(result)).catch((reason) => {
            Promise.reject(new CancelledRejectionError(parent, reason, "Unhandled rejection in cancel method."));
        }));
    }
    return Promise.all(results);
}
/**
 * Returns its argument.
 */
function identity(x) {
    return x;
}
/**
 * Throws its argument.
 */
function thrower(reason) {
    throw reason;
}
/**
 * Attempts various strategies to convert an error to a string.
 */
function errorMessage(err) {
    try {
        if (err instanceof Error || typeof err !== 'object' || err.toString !== Object.prototype.toString) {
            return "" + err;
        }
    }
    catch (_a) { }
    try {
        return JSON.stringify(err);
    }
    catch (_b) { }
    try {
        return Object.prototype.toString.call(err);
    }
    catch (_c) { }
    return "<could not convert error to string>";
}
/**
 * Gets the current barrier promise for the given cancellable promise. If necessary, initialises the barrier.
 */
function currentBarrier(promise) {
    var _a;
    let pwr = (_a = promise[barrierSym]) !== null && _a !== void 0 ? _a : {};
    if (!('promise' in pwr)) {
        Object.assign(pwr, promiseWithResolvers());
    }
    if (promise[barrierSym] == null) {
        pwr.resolve();
        promise[barrierSym] = pwr;
    }
    return pwr.promise;
}
// Polyfill Promise.withResolvers.
let promiseWithResolvers = Promise.withResolvers;
if (promiseWithResolvers && typeof promiseWithResolvers === 'function') {
    promiseWithResolvers = promiseWithResolvers.bind(Promise);
}
else {
    promiseWithResolvers = function () {
        let resolve;
        let reject;
        const promise = new Promise((res, rej) => { resolve = res; reject = rej; });
        return { promise, resolve, reject };
    };
}
