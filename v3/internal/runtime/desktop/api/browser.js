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
 * The Browser API provides methods to interact with the system browser.
 */
export const Browser = {
    /**
     * Opens a browser window to the given URL
     * @returns {Promise<string>}
     */
    OpenURL: (url) => {
        return wails.Browser.OpenURL(url);
    },
};
