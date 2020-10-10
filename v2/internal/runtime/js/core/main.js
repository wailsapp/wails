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
import { Callback, SystemCall } from './calls';
import { AddScript, InjectCSS } from './utils';
import { AddIPCListener } from 'ipc';
import * as Platform from 'platform';
import * as Store from './store';

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
			AddIPCListener,
			SystemCall,
		},
		Store,
	};

	// Setup system
	window.wails.System.IsDarkMode = Store.New('isdarkmode');
	window.wails.System.LogLevel = Store.New('loglevel');

	// Do platform specific Init
	Platform.Init();

	window.wailsloader.runtime = true;
}