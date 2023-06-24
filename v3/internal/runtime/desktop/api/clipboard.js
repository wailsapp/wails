/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

/* jshint esversion: 9 */

/**
 * The Clipboard API provides methods to interact with the system clipboard.
 */
export const Clipboard = {
    /**
     * Gets the text from the clipboard
     * @returns {Promise<string>}
     */
    Text: () => {
        return wails.Clipboard.Text();
    },
    /**
     * Sets the text on the clipboard
     * @param {string} text - text to set in the clipboard
     */
    SetText: (text) => {
        return wails.Clipboard.SetText(text);
    },
};
