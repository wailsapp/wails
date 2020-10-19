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
 * Opens the given URL / filename in the system browser
 *
 * @export
 * @param {string} target
 * @returns
 */
export function Open(target) {
	return SendMessage('RBO' + target);
}

