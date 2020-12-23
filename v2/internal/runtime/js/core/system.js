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

import * as Events from './events';
import * as Store from './store';

// Set up stores
export const LogLevel = Store.New('wails:loglevel');
export const AppConfig = Store.New('wails:appconfig');

// Set up dark mode
export let isDarkMode;

// Register system event listener to keep isDarkMode up to date
Events.On('wails:system:themechange', (darkMode) => {
	isDarkMode = darkMode;
});

