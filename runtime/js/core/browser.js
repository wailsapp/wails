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

import { SystemCall } from './calls';

/**
 * Opens the given URL in the system browser
 *
 * @export
 * @param {string} url
 * @returns
 */
export function OpenURL(url) {
	return SystemCall('Browser.OpenURL', url);
}

/**
 * Opens the given filename using the system's default file handler
 *
 * @export
 * @param {string} filename
 * @returns
 */
export function OpenFile(filename) {
	return SystemCall('Browser.OpenFile', filename);
}
