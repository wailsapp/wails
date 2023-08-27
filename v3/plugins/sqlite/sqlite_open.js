/**
 * Open a sqlite DB.
 * @param filename {string} - file to open.
 * @returns {Promise<void>}
 */
function Open(filename) {
    return wails.CallByID(147348976, filename);
}
