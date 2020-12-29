/*
 _       __      _ __
| |     / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The lightweight framework for web-like apps
(c) Lea Anthony 2019-present
*/

/* jshint esversion: 6 */

import { SendMessage } from 'ipc';

/**
 * Sets the tray icon to the icon referenced by the given ID.
 * Tray icons must follow this convention:
 *   - They must be PNG files
 *   - They must reside in a "trayicons" directory in the project root
 *   - They must have a ".png" extension
 *
 * The icon ID is the name of the file, without the ".png"
 *
 * @param {string} trayIconID - The tray icon ID
 */
export function SetIcon(trayIconID) {
	SendMessage('TI' + trayIconID);
}
