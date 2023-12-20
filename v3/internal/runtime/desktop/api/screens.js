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
 * The Screens API provides methods to interact with the system screens/monitors.
 */
export const Screens = {
    /**
     * Get the primary screen
     * @returns {Promise<Screen>}
     */
    GetPrimary: () => wails.Screens.GetPrimary(),
    /**
     * Get all screens
     * @returns {Promise<Screen[]>}
     */
    GetAll: () => wails.Screens.GetAll(),
    /**
     * Get the current screen
     * @returns {Promise<Screen>}
     */
    GetCurrent: () => wails.Screens.GetCurrent(),
};
