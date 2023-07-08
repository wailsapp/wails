
/**
 * Starts the oauth temporary server.
 * @returns {Promise<void>}
 */
function Start() {
    return wails.Plugin("oauth", "Start");
}

export default {
    OAuth: {
        Start,
    }
};
