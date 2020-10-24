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
 * Open a dialog with the given parameters
 *
 * @export
 * @param {object} options
 */
export function Open(options) {
	SendMessage('DO'+JSON.stringify(options));
}
