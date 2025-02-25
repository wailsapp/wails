/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

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
    constructor(message?: string, options?: ErrorOptions) {
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
     * Holds a reference to the promise that was cancelled and then rejected.
     */
    promise: CancellablePromise<unknown>;

    /**
     * Constructs a new `CancelledRejectionError` instance.
     * @param promise - The promise that caused the error originally.
     * @param reason - The rejection reason.
     * @param info - An optional informative message specifying the circumstances in which the error was thrown.
     *               Defaults to the string `"Unhandled rejection in cancelled promise."`.
     */
    constructor(promise: CancellablePromise<unknown>, reason?: any, info?: string) {
        super((info ?? "Unhandled rejection in cancelled promise.") + " Reason: " + errorMessage(reason), { cause: reason });
        this.promise = promise;
        this.name = "CancelledRejectionError";
    }
}

type CancellablePromiseResolver<T> = (value: T | PromiseLike<T> | CancellablePromiseLike<T>) => void;
type CancellablePromiseRejector = (reason?: any) => void;
type CancellablePromiseCanceller = (cause?: any) => void | PromiseLike<void>;
type CancellablePromiseExecutor<T> = (resolve: CancellablePromiseResolver<T>, reject: CancellablePromiseRejector) => void;

export interface CancellablePromiseLike<T> {
    then<TResult1 = T, TResult2 = never>(onfulfilled?: ((value: T) => TResult1 | PromiseLike<TResult1> | CancellablePromiseLike<TResult1>) | undefined | null, onrejected?: ((reason: any) => TResult2 | PromiseLike<TResult2> | CancellablePromiseLike<TResult2>) | undefined | null): CancellablePromiseLike<TResult1 | TResult2>;
    cancel(cause?: any): void | PromiseLike<void>;
}

/**
 * Wraps a cancellable promise along with its resolution methods.
 * The `oncancelled` field will be null initially but may be set to provide a custom cancellation function.
 */
export interface CancellablePromiseWithResolvers<T> {
    promise: CancellablePromise<T>;
    resolve: CancellablePromiseResolver<T>;
    reject: CancellablePromiseRejector;
    oncancelled: CancellablePromiseCanceller | null;
}

interface CancellablePromiseState {
    readonly root: CancellablePromiseState;
    resolving: boolean;
    settled: boolean;
    reason?: CancelError;
}

// Private field names.
const barrierSym = Symbol("barrier");
const cancelImplSym = Symbol("cancelImpl");
const species = Symbol.species ?? Symbol("speciesPolyfill");

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
export class CancellablePromise<T> extends Promise<T> implements PromiseLike<T>, CancellablePromiseLike<T> {
    // Private fields.
    /** @internal */
    private [barrierSym]!: Partial<PromiseWithResolvers<void>> | null;
    /** @internal */
    private readonly [cancelImplSym]!: (reason: CancelError) => void | PromiseLike<void>;

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
    constructor(executor: CancellablePromiseExecutor<T>, oncancelled?: CancellablePromiseCanceller) {
        let resolve!: (value: T | PromiseLike<T>) => void;
        let reject!: (reason?: any) => void;
        super((res, rej) => { resolve = res; reject = rej; });

        if ((this.constructor as any)[species] !== Promise) {
            throw new TypeError("CancellablePromise does not support transparent subclassing. Please refrain from overriding the [Symbol.species] static property.");
        }

        let promise: CancellablePromiseWithResolvers<T> = {
            promise: this,
            resolve,
            reject,
            get oncancelled() { return oncancelled ?? null; },
            set oncancelled(cb) { oncancelled = cb ?? undefined; }
        };

        const state: CancellablePromiseState = {
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
        } catch (err) {
            if (state.resolving) {
                console.log("Unhandled exception in CancellablePromise executor.", err);
            } else {
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
    cancel(cause?: any): CancellablePromise<void> {
        return new CancellablePromise<void>((resolve) => {
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
    cancelOn(signal: AbortSignal): CancellablePromise<T> {
        if (signal.aborted) {
            void this.cancel(signal.reason)
        } else {
            signal.addEventListener('abort', () => void this.cancel(signal.reason), {capture: true});
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
    then<TResult1 = T, TResult2 = never>(onfulfilled?: ((value: T) => TResult1 | PromiseLike<TResult1> | CancellablePromiseLike<TResult1>) | undefined | null, onrejected?: ((reason: any) => TResult2 | PromiseLike<TResult2> | CancellablePromiseLike<TResult2>) | undefined | null, oncancelled?: CancellablePromiseCanceller): CancellablePromise<TResult1 | TResult2> {
        if (!(this instanceof CancellablePromise)) {
            throw new TypeError("CancellablePromise.prototype.then called on an invalid object.");
        }

        // NOTE: TypeScript's built-in type for then is broken,
        // as it allows specifying an arbitrary TResult1 != T even when onfulfilled is not a function.
        // We cannot fix it if we want to CancellablePromise to implement PromiseLike<T>.

        if (!isCallable(onfulfilled)) { onfulfilled = identity as any; }
        if (!isCallable(onrejected)) { onrejected = thrower; }

        if (onfulfilled === identity && onrejected == thrower) {
            // Shortcut for trivial arguments.
            return new CancellablePromise((resolve) => resolve(this as any));
        }

        const barrier: Partial<PromiseWithResolvers<void>> = {};
        this[barrierSym] = barrier;

        return new CancellablePromise<TResult1 | TResult2>((resolve, reject) => {
            void promiseThen.call(this,
                (value) => {
                    if (this[barrierSym] === barrier) { this[barrierSym] = null; }
                    barrier.resolve?.();

                    try {
                        resolve(onfulfilled!(value));
                    } catch (err) {
                        reject(err);
                    }
                },
                (reason?) => {
                    if (this[barrierSym] === barrier) { this[barrierSym] = null; }
                    barrier.resolve?.();

                    try {
                        resolve(onrejected!(reason));
                    } catch (err) {
                        reject(err);
                    }
                }
            );
        }, async (cause?) => {
            //cancelled = true;
            try {
                return oncancelled?.(cause);
            } finally {
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
    catch<TResult = never>(onrejected?: ((reason: any) => (PromiseLike<TResult> | TResult)) | undefined | null, oncancelled?: CancellablePromiseCanceller): CancellablePromise<T | TResult> {
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
    finally(onfinally?: (() => void) | undefined | null, oncancelled?: CancellablePromiseCanceller): CancellablePromise<T> {
        if (!(this instanceof CancellablePromise)) {
            throw new TypeError("CancellablePromise.prototype.finally called on an invalid object.");
        }

        if (!isCallable(onfinally)) {
            return this.then(onfinally, onfinally, oncancelled);
        }

        return this.then(
            (value) => CancellablePromise.resolve(onfinally()).then(() => value),
            (reason?) => CancellablePromise.resolve(onfinally()).then(() => { throw reason; }),
            oncancelled,
        );
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

    /**
     * Creates a CancellablePromise that is resolved with an array of results
     * when all of the provided Promises resolve, or rejected when any Promise is rejected.
     *
     * Every one of the provided objects that is a thenable _and_ cancellable object
     * will be cancelled when the returned promise is cancelled, with the same cause.
     *
     * @group Static Methods
     */
    static all<T>(values: Iterable<T | PromiseLike<T>>): CancellablePromise<Awaited<T>[]>;
    static all<T extends readonly unknown[] | []>(values: T): CancellablePromise<{ -readonly [P in keyof T]: Awaited<T[P]>; }>;
    static all<T extends Iterable<unknown> | ArrayLike<unknown>>(values: T): CancellablePromise<unknown> {
        let collected = Array.from(values);
        const promise = collected.length === 0
            ? CancellablePromise.resolve(collected)
            : new CancellablePromise<unknown>((resolve, reject) => {
                void Promise.all(collected).then(resolve, reject);
            }, (cause?): Promise<void> => cancelAll(promise, collected, cause));
        return promise;
    }

    /**
     * Creates a CancellablePromise that is resolved with an array of results
     * when all of the provided Promises resolve or reject.
     *
     * Every one of the provided objects that is a thenable _and_ cancellable object
     * will be cancelled when the returned promise is cancelled, with the same cause.
     *
     * @group Static Methods
     */
    static allSettled<T>(values: Iterable<T | PromiseLike<T>>): CancellablePromise<PromiseSettledResult<Awaited<T>>[]>;
    static allSettled<T extends readonly unknown[] | []>(values: T): CancellablePromise<{ -readonly [P in keyof T]: PromiseSettledResult<Awaited<T[P]>>; }>;
    static allSettled<T extends Iterable<unknown> | ArrayLike<unknown>>(values: T): CancellablePromise<unknown> {
        let collected = Array.from(values);
        const promise = collected.length === 0
            ? CancellablePromise.resolve(collected)
            : new CancellablePromise<unknown>((resolve, reject) => {
                void Promise.allSettled(collected).then(resolve, reject);
            }, (cause?): Promise<void> => cancelAll(promise, collected, cause));
        return promise;
    }

    /**
     * The any function returns a promise that is fulfilled by the first given promise to be fulfilled,
     * or rejected with an AggregateError containing an array of rejection reasons
     * if all of the given promises are rejected.
     * It resolves all elements of the passed iterable to promises as it runs this algorithm.
     *
     * Every one of the provided objects that is a thenable _and_ cancellable object
     * will be cancelled when the returned promise is cancelled, with the same cause.
     *
     * @group Static Methods
     */
    static any<T>(values: Iterable<T | PromiseLike<T>>): CancellablePromise<Awaited<T>>;
    static any<T extends readonly unknown[] | []>(values: T): CancellablePromise<Awaited<T[number]>>;
    static any<T extends Iterable<unknown> | ArrayLike<unknown>>(values: T): CancellablePromise<unknown> {
        let collected = Array.from(values);
        const promise = collected.length === 0
            ? CancellablePromise.resolve(collected)
            : new CancellablePromise<unknown>((resolve, reject) => {
                void Promise.any(collected).then(resolve, reject);
            }, (cause?): Promise<void> => cancelAll(promise, collected, cause));
        return promise;
    }

    /**
     * Creates a Promise that is resolved or rejected when any of the provided Promises are resolved or rejected.
     *
     * Every one of the provided objects that is a thenable _and_ cancellable object
     * will be cancelled when the returned promise is cancelled, with the same cause.
     *
     * @group Static Methods
     */
    static race<T>(values: Iterable<T | PromiseLike<T>>): CancellablePromise<Awaited<T>>;
    static race<T extends readonly unknown[] | []>(values: T): CancellablePromise<Awaited<T[number]>>;
    static race<T extends Iterable<unknown> | ArrayLike<unknown>>(values: T): CancellablePromise<unknown> {
        let collected = Array.from(values);
        const promise = new CancellablePromise<unknown>((resolve, reject) => {
            void Promise.race(collected).then(resolve, reject);
        }, (cause?): Promise<void> => cancelAll(promise, collected, cause));
        return promise;
    }

    /**
     * Creates a new cancelled CancellablePromise for the provided cause.
     *
     * @group Static Methods
     */
    static cancel<T = never>(cause?: any): CancellablePromise<T> {
        const p = new CancellablePromise<T>(() => {});
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
    static timeout<T = never>(milliseconds: number, cause?: any): CancellablePromise<T> {
        const promise = new CancellablePromise<T>(() => {});
        if (AbortSignal && typeof AbortSignal === 'function' && AbortSignal.timeout && typeof AbortSignal.timeout === 'function') {
            AbortSignal.timeout(milliseconds).addEventListener('abort', () => void promise.cancel(cause));
        } else {
            setTimeout(() => void promise.cancel(cause), milliseconds);
        }
        return promise;
    }

    /**
     * Creates a new CancellablePromise that resolves after the specified timeout.
     * The returned promise can be cancelled without consequences.
     *
     * @group Static Methods
     */
    static sleep(milliseconds: number): CancellablePromise<void>;
    /**
     * Creates a new CancellablePromise that resolves after
     * the specified timeout, with the provided value.
     * The returned promise can be cancelled without consequences.
     *
     * @group Static Methods
     */
    static sleep<T>(milliseconds: number, value: T): CancellablePromise<T>;
    static sleep<T = void>(milliseconds: number, value?: T): CancellablePromise<T> {
        return new CancellablePromise<T>((resolve) => {
            setTimeout(() => resolve(value!), milliseconds);
        });
    }

    /**
     * Creates a new rejected CancellablePromise for the provided reason.
     *
     * @group Static Methods
     */
    static reject<T = never>(reason?: any): CancellablePromise<T> {
        return new CancellablePromise<T>((_, reject) => reject(reason));
    }

    /**
     * Creates a new resolved CancellablePromise.
     *
     * @group Static Methods
     */
    static resolve(): CancellablePromise<void>;
    /**
     * Creates a new resolved CancellablePromise for the provided value.
     *
     * @group Static Methods
     */
    static resolve<T>(value: T): CancellablePromise<Awaited<T>>;
    /**
     * Creates a new resolved CancellablePromise for the provided value.
     *
     * @group Static Methods
     */
    static resolve<T>(value: T | PromiseLike<T>): CancellablePromise<Awaited<T>>;
    static resolve<T = void>(value?: T | PromiseLike<T>): CancellablePromise<Awaited<T>> {
        if (value instanceof CancellablePromise) {
            // Optimise for cancellable promises.
            return value;
        }
        return new CancellablePromise<any>((resolve) => resolve(value));
    }

    /**
     * Creates a new CancellablePromise and returns it in an object, along with its resolve and reject functions
     * and a getter/setter for the cancellation callback.
     *
     * This method is polyfilled, hence available in every OS/webview version.
     *
     * @group Static Methods
     */
    static withResolvers<T>(): CancellablePromiseWithResolvers<T> {
        let result: CancellablePromiseWithResolvers<T> = { oncancelled: null } as any;
        result.promise = new CancellablePromise<T>((resolve, reject) => {
            result.resolve = resolve;
            result.reject = reject;
        }, (cause?: any) => { result.oncancelled?.(cause); });
        return result;
    }
}

/**
 * Returns a callback that implements the cancellation algorithm for the given cancellable promise.
 * The promise returned from the resulting function does not reject.
 */
function cancellerFor<T>(promise: CancellablePromiseWithResolvers<T>, state: CancellablePromiseState) {
    let cancellationPromise: void | PromiseLike<void> = undefined;

    return (reason: CancelError): void | PromiseLike<void> => {
        if (!state.settled) {
            state.settled = true;
            state.reason = reason;
            promise.reject(reason);

            // Attach an error handler that ignores this specific rejection reason and nothing else.
            // In theory, a sane underlying implementation at this point
            // should always reject with our cancellation reason,
            // hence the handler will never throw.
            void promiseThen.call(promise.promise, undefined, (err) => {
                if (err !== reason) {
                    throw err;
                }
            });
        }

        // If reason is not set, the promise resolved regularly, hence we must not call oncancelled.
        // If oncancelled is unset, no need to go any further.
        if (!state.reason || !promise.oncancelled) { return; }

        cancellationPromise = new Promise<void>((resolve) => {
            try {
                resolve(promise.oncancelled!(state.reason!.cause));
            } catch (err) {
                Promise.reject(new CancelledRejectionError(promise.promise, err, "Unhandled exception in oncancelled callback."));
            }
        }).catch((reason?) => {
            Promise.reject(new CancelledRejectionError(promise.promise, reason, "Unhandled rejection in oncancelled callback."));
        });

        // Unset oncancelled to prevent repeated calls.
        promise.oncancelled = null;

        return cancellationPromise;
    }
}

/**
 * Returns a callback that implements the resolution algorithm for the given cancellable promise.
 */
function resolverFor<T>(promise: CancellablePromiseWithResolvers<T>, state: CancellablePromiseState): CancellablePromiseResolver<T> {
    return (value) => {
        if (state.resolving) { return; }
        state.resolving = true;

        if (value === promise.promise) {
            if (state.settled) { return; }
            state.settled = true;
            promise.reject(new TypeError("A promise cannot be resolved with itself."));
            return;
        }

        if (value != null && (typeof value === 'object' || typeof value === 'function')) {
            let then: any;
            try {
                then = (value as any).then;
            } catch (err) {
                state.settled = true;
                promise.reject(err);
                return;
            }

            if (isCallable(then)) {
                try {
                    let cancel = (value as any).cancel;
                    if (isCallable(cancel)) {
                        const oncancelled = (cause?: any) => {
                            Reflect.apply(cancel, value, [cause]);
                        };
                        if (state.reason) {
                            // If already cancelled, propagate cancellation.
                            // The promise returned from the canceller algorithm does not reject
                            // so it can be discarded safely.
                            void cancellerFor({ ...promise, oncancelled }, state)(state.reason);
                        } else {
                            promise.oncancelled = oncancelled;
                        }
                    }
                } catch {}

                const newState: CancellablePromiseState = {
                    root: state.root,
                    resolving: false,
                    get settled() { return this.root.settled },
                    set settled(value) { this.root.settled = value; },
                    get reason() { return this.root.reason }
                };

                const rejector = rejectorFor(promise, newState);
                try {
                    Reflect.apply(then, value, [resolverFor(promise, newState), rejector]);
                } catch (err) {
                    rejector(err);
                }
                return; // IMPORTANT!
            }
        }

        if (state.settled) { return; }
        state.settled = true;
        promise.resolve(value);
    };
}

/**
 * Returns a callback that implements the rejection algorithm for the given cancellable promise.
 */
function rejectorFor<T>(promise: CancellablePromiseWithResolvers<T>, state: CancellablePromiseState): CancellablePromiseRejector {
    return (reason?) => {
        if (state.resolving) { return; }
        state.resolving = true;

        if (state.settled) {
            try {
                if (reason instanceof CancelError && state.reason instanceof CancelError && Object.is(reason.cause, state.reason.cause)) {
                    // Swallow late rejections that are CancelErrors whose cancellation cause is the same as ours.
                    return;
                }
            } catch {}

            void Promise.reject(new CancelledRejectionError(promise.promise, reason));
        } else {
            state.settled = true;
            promise.reject(reason);
        }
    }
}

/**
 * Cancels all values in an array that look like cancellable thenables.
 * Returns a promise that fulfills once all cancellation procedures for the given values have settled.
 */
function cancelAll(parent: CancellablePromise<unknown>, values: any[], cause?: any): Promise<void> {
    const results = [];

    for (const value of values) {
        let cancel: CancellablePromiseCanceller;
        try {
            if (!isCallable(value.then)) { continue; }
            cancel = value.cancel;
            if (!isCallable(cancel)) { continue; }
        } catch { continue; }

        let result: void | PromiseLike<void>;
        try {
            result = Reflect.apply(cancel, value, [cause]);
        } catch (err) {
            Promise.reject(new CancelledRejectionError(parent, err, "Unhandled exception in cancel method."));
            continue;
        }

        if (!result) { continue; }
        results.push(
            (result instanceof Promise  ? result : Promise.resolve(result)).catch((reason?) => {
                Promise.reject(new CancelledRejectionError(parent, reason, "Unhandled rejection in cancel method."));
            })
        );
    }

    return Promise.all(results) as any;
}

/**
 * Returns its argument.
 */
function identity<T>(x: T): T {
    return x;
}

/**
 * Throws its argument.
 */
function thrower(reason?: any): never {
    throw reason;
}

/**
 * Attempts various strategies to convert an error to a string.
 */
function errorMessage(err: any): string {
    try {
        if (err instanceof Error || typeof err !== 'object' || err.toString !== Object.prototype.toString) {
            return "" + err;
        }
    } catch {}

    try {
        return JSON.stringify(err);
    } catch {}

    try {
        return Object.prototype.toString.call(err);
    } catch {}

    return "<could not convert error to string>";
}

/**
 * Gets the current barrier promise for the given cancellable promise. If necessary, initialises the barrier.
 */
function currentBarrier<T>(promise: CancellablePromise<T>): Promise<void> {
    let pwr: Partial<PromiseWithResolvers<void>> = promise[barrierSym] ?? {};
    if (!('promise' in pwr)) {
        Object.assign(pwr, promiseWithResolvers<void>());
    }
    if (promise[barrierSym] == null) {
        pwr.resolve!();
        promise[barrierSym] = pwr;
    }
    return pwr.promise!;
}

// Stop sneaky people from breaking the barrier mechanism.
const promiseThen = Promise.prototype.then;
Promise.prototype.then = function(...args) {
    if (this instanceof CancellablePromise) {
        return this.then(...args);
    } else {
        return Reflect.apply(promiseThen, this, args);
    }
}

// Polyfill Promise.withResolvers.
let promiseWithResolvers = Promise.withResolvers;
if (promiseWithResolvers && typeof promiseWithResolvers === 'function') {
    promiseWithResolvers = promiseWithResolvers.bind(Promise);
} else {
    promiseWithResolvers = function <T>(): PromiseWithResolvers<T> {
        let resolve!: (value: T | PromiseLike<T>) => void;
        let reject!: (reason?: any) => void;
        const promise = new Promise<T>((res, rej) => { resolve = res; reject = rej; });
        return { promise, resolve, reject };
    }
}