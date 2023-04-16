
/**
 * Opens the given URL in the default browser.
 * @param url {string} - The URL to open.
 * @returns {Promise<void>}
 */
function OpenURL(url) {
    return wails.Plugin("browser", "OpenURL", url);
}

/**
 * Opens the given filename in the default browser.
 * @param filename {string} - The file to open.
 * @returns {Promise<void>}
 */
function OpenFile(filename) {
    return wails.Plugin("browser", "OpenFile", filename);
}

export default {
    Browser: {
        OpenURL,
        OpenFile,
    }
};
