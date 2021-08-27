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
import * as Log from './log';
import {Emit, Notify, On, Once, OnMultiple} from './events';
import {Callback, SystemCall} from './calls';
import {AddScript, DisableDefaultContextMenu, InjectCSS} from './utils';
import {AddIPCListener, SendMessage} from 'ipc';
import * as Platform from 'platform';

export function Init() {
	// Backend is where the Go struct wrappers get bound to
	window.backend = {};

	// Initialise global if not already
	window.wails = {
		Log,
		Events: {
			On,
			Once,
			OnMultiple,
			Emit,
		},
		_: {
			Callback,
			Notify,
			AddScript,
			InjectCSS,
			DisableDefaultContextMenu,
			// Init,
			AddIPCListener,
			SystemCall,
			SendMessage,
		},
	};

	// Do platform specific Init
	Platform.Init();
}