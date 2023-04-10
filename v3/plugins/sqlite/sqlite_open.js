/**
 * Open a sqlite DB.
 * @param filename {string} - file to open.
 * @returns {Promise<void>}
 */
function Open(filename) {
    return wails.Plugin("sqlite", "Open", filename);
}
