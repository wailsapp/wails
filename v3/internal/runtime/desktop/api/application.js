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
 * The Application API provides methods to interact with the application.
 */
export const Application = {
    /**
     * Hides the application
     * @returns {Promise<void>}
     */
    Hide: () => wails.Application.Hide(),
    /**
     * Shows the application
     * @returns {Promise<void>}
     */
    Show: () => wails.Application.Show(),
    /**
     * Quits the application
     * @returns {Promise<void>}
     */
    Quit: () => wails.Application.Quit(),
};
