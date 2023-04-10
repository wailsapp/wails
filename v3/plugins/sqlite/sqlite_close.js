
/**
 * Close a sqlite DB.
 * @returns {Promise<void>}
 */
function Close() {
    return wails.Plugin("sqlite", "Close");
}
