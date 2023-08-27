
/**
 * Opens the given URL in the default browser.
 * @param url {string} - The URL to open.
 * @returns {Promise<void>}
 */
function OpenURL(url) {
    return wails.CallByID(3188315539, url);
}

/**
 * Opens the given filename in the default browser.
 * @param filename {string} - The file to open.
 * @returns {Promise<void>}
 */
function OpenFile(filename) {
    return wails.CallByID(3105408620, filename);
}

export default {
    Browser: {
        OpenURL,
        OpenFile,
    }
};
