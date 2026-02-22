//@ts-check

//@ts-ignore: Unused imports
import * as $models from "./models.js";

const execSymbol = Symbol("exec"),
      querySymbol = Symbol("query"),
      closeSymbol = Symbol("close");

/**
 * Stmt represents a prepared statement for later queries or executions.
 * Multiple queries or executions may be run concurrently on the same statement.
 *
 * The caller must call the statement's Close method when it is no longer needed.
 * Statements are closed automatically
 * when the connection they are associated with is closed.
 */
export class Stmt {
    /**
     * Constructs a new prepared statement instance.
     * @param {(...args: any[]) => Promise<void>} close
     * @param {(...args: any[]) => Promise<void> & { cancel(): void }} exec
     * @param {(...args: any[]) => Promise<$models.Rows> & { cancel(): void }} query
     */
    constructor(close, exec, query) {
        /**
         * @member
         * @private
         * @type {typeof close}
         */
        this[closeSymbol] = close;

        /**
         * @member
         * @private
         * @type {typeof exec}
         */
        this[execSymbol] = exec;

        /**
         * @member
         * @private
         * @type {typeof query}
         */
        this[querySymbol] = query;
    }

    /**
     * Closes the prepared statement.
     * It has no effect when the statement is already closed.
     * @returns {Promise<void>}
     */
    Close() {
        return this[closeSymbol]();
    }

    /**
     * Executes the prepared statement without returning any rows.
     * It supports early cancellation.
     *
     * @param {any[]} args
     * @returns {Promise<void> & { cancel(): void }}
     */
    Exec(...args) {
        return this[execSymbol](...args);
    }

    /**
     * Executes the prepared statement
     * and returns a slice of key-value records, one per row, with column names as keys.
     * It supports early cancellation, returning the array of results fetched so far.
     *
     * @param {any[]} args
     * @returns {Promise<$models.Rows> & { cancel(): void }}
     */
    Query(...args) {
        return this[querySymbol](...args);
    }
}
