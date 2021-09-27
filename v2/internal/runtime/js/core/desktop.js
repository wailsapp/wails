/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/
/* jshint esversion: 6 */
import {SetBindings} from './bindings';
import {Init} from './main';

// Initialise the Runtime
Init();

// Load Bindings if they exist
if (window.wailsbindings) {
	SetBindings(window.wailsbindings);
}

// Emit loaded event. Leaving this for now. It will show any errors if runtime fails to load.
window.wails.Events.Emit('wails:loaded');

