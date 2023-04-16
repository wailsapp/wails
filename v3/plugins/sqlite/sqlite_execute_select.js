
/**
 * Execute a SQL statement.
 * @param statement {string} - SQL statement to execute.
 @param args {...any} - Arguments to pass to the statement.
 * @returns {Promise<void>}
 */
function Execute(statement, ...args) {
    return wails.Plugin("sqlite", "Execute", statement, ...args);
}

/**
 * Perform a select query.
 * @param statement {string} - Select SQL statement.
 * @param args {...any} - Arguments to pass to the statement.
 * @returns {Promise<any>}
 */
function Select(statement, ...args) {
    return wails.Plugin("sqlite", "Select", statement, ...args);
}
