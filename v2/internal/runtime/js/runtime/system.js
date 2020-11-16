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

const Events = require('./events');

/**
 * Registers an event listener that will be invoked when the user changes the
 * desktop theme (light mode / dark mode). The callback receives a boolean which
 * indicates if dark mode is enabled.
 *
 * @export
 * @param {function} callback The callback to invoke on theme change
 */
function OnThemeChange(callback) {
	Events.On("wails:system:themechange", callback);
}

/**
 * Checks if dark mode is curently enabled.
 *
 * @export
 * @returns {Promise}
 */
function DarkModeEnabled() {
	return window.wails.System.IsDarkMode.get();
}

/**
 * Mac Application Config
 * @typedef {Object} MacAppConfig
 * @param {MacTitleBar} TitleBar - The window's titlebar configuration
 */

 /**
 * Mac Title Bar Config
 * Check out https://github.com/lukakerr/NSWindowStyles for some examples of these settings
 * @typedef {Object} MacTitleBar
 * @param {bool} TitleBarAppearsTransparent - NSWindow.titleBarAppearsTransparent
 * @param {bool} HideTitle - NSWindow.hideTitle
 * @param {bool} HideTitleBar - NSWindow.hideTitleBar
 * @param {bool} FullSizeContent - Makes the webview portion of the window the full size of the window, even over the titlebar
 * @param {bool} UseToolbar - Set true to add a blank toolbar to the window (makes the title bar larger)
 * @param {bool} HideToolbarSeparator - Set true to remove the separator between the toolbar and the main content area
 * 
 */

/**
 * The application configuration
 * 
 * @typedef {Object} AppConfig
 * @param {string} Title - Application Title
 * @param {number} Width - Window Width
 * @param {number} Height - Window Height
 * @param {boolean} DisableResize - True if resize is disabled
 * @param {boolean} Fullscreen - App started in fullscreen
 * @param {number} MinWidth - Window Minimum Width
 * @param {number} MinHeight - Window Minimum Height
 * @param {number} MaxWidth - Window Maximum Width
 * @param {number} MaxHeight - Window Maximum Height
 * @param {bool} StartHidden - Start with window hidden
 * @param {bool} DevTools - Enables the window devtools
 * @param {number} RBGA - The initial window colour. Convert to hex then it'll mean 0xRRGGBBAA
 * @param {MacAppConfig} [Mac] - Configuration when running on Mac
 * @param {LinuxAppConfig} [Linux] - Configuration when running on Linux
 * @param {WindowsAppConfig} [Windows] - Configuration when running on Windows
 * @param {string} Appearance - The default application appearance. Use the values listed here: https://developer.apple.com/documentation/appkit/nsappearance?language=objc
 * @param {number} WebviewIsTransparent - Makes the background of the webview content transparent. Use this with the Alpha part of the window colour to make parts of your application transparent.
 * @param {number} WindowBackgroundIsTranslucent - Makes the transparent parts of the application window translucent. Example: https://en.wikipedia.org/wiki/MacOS_Big_Sur#/media/File:MacOS_Big_Sur_-_Safari_Extensions_category_in_App_Store.jpg
 * @param {number} LogLevel - The initial log level (lower is more verbose)
 * 
 */

/**
 * Returns the application configuration.
 *
 * @export
 * @returns {Promise<AppConfig>}
 */
function AppConfig() {
	return window.wails.System.AppConfig.get();
}

module.exports = {
	OnThemeChange: OnThemeChange,
	DarkModeEnabled: DarkModeEnabled,
	LogLevel: window.wails.System.LogLevel,
	Platform: window.wails.System.Platform,
	AppType: window.wails.System.AppType,
	AppConfig: AppConfig,
};