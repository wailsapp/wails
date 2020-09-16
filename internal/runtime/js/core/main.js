/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The lightweight framework for web-like apps
(c) Lea Anthony 2019-present
*/
/* jshint esversion: 6 */
import * as Log from './log';
import * as Browser from './browser';
import * as Window from './window';
import { On, OnMultiple, Emit, Notify, Heartbeat, Acknowledge } from './events';
import { Callback } from './calls';
import { AddScript, InjectCSS } from './utils';
import { AddIPCListener } from 'ipc';
import * as Platform from 'platform';

export function Init() {
	// Backend is where the Go struct wrappers get bound to
	window.backend = {};

	// Initialise global if not already
	window.wails = {
		System: Platform.System,
		Log,
		Browser,
		Window,
		Events: {
			On,
			OnMultiple,
			Emit,
			Heartbeat,
			Acknowledge,
		},
		_: {
			Callback,
			Notify,
			AddScript,
			InjectCSS,
			Init,
			AddIPCListener
		}
	};

	// Do platform specific Init
	Platform.Init();
}