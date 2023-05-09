/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/
/* jshint esversion: 9 */

/**
 * @typedef {import("./types").MessageDialogOptions} MessageDialogOptions
 * @typedef {import("./types").OpenDialogOptions} OpenDialogOptions
 * @typedef {import("./types").SaveDialogOptions} SaveDialogOptions
 * @typedef {import("./types").Screen} Screen
 * @typedef {import("./types").Size} Size
 * @typedef {import("./types").Position} Position
 *
 */

/**
 * The Clipboard API provides methods to interact with the system clipboard.
 */
export const Clipboard = {
    /**
     * Gets the text from the clipboard
     * @returns {Promise<string>}
     */
    Text: () => {
        return wails.Clipboard.Text();
    },
    /**
     * Sets the text on the clipboard
     * @param {string} text - text to set in the clipboard
     */
    SetText: (text) => {
        return wails.Clipboard.SetText(text);
    },
};

/**
 * The Application API provides methods to interact with the application.
 */
export const Application = {
    /**
     * Hides the application
     */
    Hide: () => {
        return wails.Application.Hide();
    },
    /**
     * Shows the application
     */
    Show: () => {
        return wails.Application.Show();
    },
    /**
     * Quits the application
     */
    Quit: () => {
        return wails.Application.Quit();
    },
};

/**
 * The Screens API provides methods to interact with the system screens/monitors.
 */
export const Screens = {
    /**
     * Get the primary screen
     * @returns {Promise<Screen>}
     */
    GetPrimary: () => {
        return wails.Screens.GetPrimary();
    },
    /**
     * Get all screens
     * @returns {Promise<Screen[]>}
     */
    GetAll: () => {
        return wails.Screens.GetAll();
    },
    /**
     * Get the current screen
     * @returns {Promise<Screen>}
     */
    GetCurrent: () => {
        return wails.Screens.GetCurrent();
    },
};

/**
 * Call a plugin method
 * @param {string} pluginName - name of the plugin
 * @param {string} methodName - name of the method
 * @param {...any} args - arguments to pass to the method
 * @returns {Promise<any>} - promise that resolves with the result
 */
export const Plugin = (pluginName, methodName, ...args) => {
    return wails.Plugin(pluginName, methodName, ...args);
};

/**
 * The Dialog API provides methods to interact with system dialogs.
 */
export const Dialog = {
    /**
     * Shows an info dialog
     * @param {MessageDialogOptions} options - options for the dialog
     * @returns {Promise<string>}
     */
    Info: (options) => {
        return wails.Dialog.Info(options);
    },
    /**
     * Shows a warning dialog
     * @param {MessageDialogOptions} options - options for the dialog
     * @returns {Promise<string>}
     */
    Warning: (options) => {
        return wails.Dialog.Warning(options);
    },
    /**
     * Shows an error dialog
     * @param {MessageDialogOptions} options - options for the dialog
     * @returns {Promise<string>}
     */
    Error: (options) => {
        return wails.Dialog.Error(options);
    },

    /**
     * Shows a question dialog
     * @param {MessageDialogOptions} options - options for the dialog
     * @returns {Promise<string>}
     */
    Question: (options) => {
        return wails.Dialog.Question(options);
    },

    /**
     * Shows a file open dialog and returns the files selected by the user.
     * A blank string indicates that the dialog was cancelled.
     * @param {OpenDialogOptions} options - options for the dialog
     * @returns {Promise<string[]>|Promise<string>}
     */
    OpenFile: (options) => {
        return wails.Dialog.OpenFile(options);
    },

    /**
     * Shows a file save dialog and returns the filename given by the user.
     * A blank string indicates that the dialog was cancelled.
     * @param {SaveDialogOptions} options - options for the dialog
     * @returns {Promise<string>}
     */
    SaveFile: (options) => {
        return wails.Dialog.SaveFile(options);
    },
};

/**
 * The Events API provides methods to interact with the event system.
 */
export const Events = {
    /**
     * Emit an event
     * @param {string} name
     * @param {any=} data
     */
    Emit: (name, data) => {
        return wails.Events.Emit(name, data);
    },
    /**
     * Subscribe to an event
     * @param {string} name - name of the event
     * @param {(any) => void} callback - callback to call when the event is emitted
     @returns {function()} unsubscribeMethod - method to unsubscribe from the event
     */
    On: (name, callback) => {
        return wails.Events.On(name, callback);
    },
    /**
     * Subscribe to an event once
     * @param {string} name - name of the event
     * @param {(any) => void} callback - callback to call when the event is emitted
     * @returns {function()} unsubscribeMethod - method to unsubscribe from the event
     */
    Once: (name, callback) => {
        return wails.Events.Once(name, callback);
    },
    /**
     * Subscribe to an event multiple times
     * @param {string} name - name of the event
     * @param {(any) => void} callback - callback to call when the event is emitted
     * @param {number} count - number of times to call the callback
     * @returns {Promise<void>} unsubscribeMethod - method to unsubscribe from the event
     */
    OnMultiple: (name, callback, count) => {
        return wails.Events.OnMultiple(name, callback, count);
    },
    /**
     * Unsubscribe from an event
     * @param {string} name - name of the event to unsubscribe from
     * @param {...string} additionalNames - additional names of events to unsubscribe from
     */
    Off: (name, ...additionalNames) => {
        wails.Events.Off(name, additionalNames);
    },
    /**
     * Unsubscribe all listeners from all events
     */
    OffAll: () => {
        wails.Events.OffAll();
    },
};

/**
 * The Window API provides methods to interact with the window.
 */
export const Window = {
    /**
     * Center the window.
     */
    Center: () => void wails.Window.Center(),
    /**
     * Set the window title.
     * @param title
     */
    SetTitle: (title) => void wails.Window.SetTitle(title),

    /**
     * Makes the window fullscreen.
     */
    Fullscreen: () => void wails.Window.Fullscreen(),

    /**
     * Unfullscreen the window.
     */
    UnFullscreen: () => void wails.Window.UnFullscreen(),

    /**
     * Set the window size.
     * @param {number} width The window width
     * @param {number} height The window height
     */
    SetSize: (width, height) => void wails.Window.SetSize(width, height),

    /**
     * Get the window size.
     * @returns {Promise<Size>} The window size
     */
    Size: () => {
        return wails.Window.Size();
    },

    /**
     * Set the window maximum size.
     * @param {number} width
     * @param {number} height
     */
    SetMaxSize: (width, height) => void wails.Window.SetMaxSize(width, height),

    /**
     * Set the window minimum size.
     * @param {number} width
     * @param {number} height
     */
    SetMinSize: (width, height) => void wails.Window.SetMinSize(width, height),

    /**
     * Set window to be always on top.
     * @param {boolean} onTop Whether the window should be always on top
     */
    SetAlwaysOnTop: (onTop) => void wails.Window.SetAlwaysOnTop(onTop),

    /**
     * Set the window position.
     * @param {number} x
     * @param {number} y
     */
    SetPosition: (x, y) => void wails.Window.SetPosition(x, y),

    /**
     * Get the window position.
     * @returns {Promise<Position>} The window position
     */
    Position: () => {
        return wails.Window.Position();
    },

    /**
     * Get the screen the window is on.
     * @returns {Promise<Screen>}
     */
    Screen: () => {
        return wails.Window.Screen();
    },

    /**
     * Hide the window
     */
    Hide: () => void wails.Window.Hide(),

    /**
     * Maximise the window
     */
    Maximise: () => void wails.Window.Maximise(),

    /**
     * Show the window
     */
    Show: () => void wails.Window.Show(),

    /**
     * Close the window
     */
    Close: () => void wails.Window.Close(),

    /**
     * Toggle the window maximise state
     */
    ToggleMaximise: () => void wails.Window.ToggleMaximise(),

    /**
     * Unmaximise the window
     */
    UnMaximise: () => void wails.Window.UnMaximise(),

    /**
     * Minimise the window
     */
    Minimise: () => void wails.Window.Minimise(),

    /**
     * Unminimise the window
     */
    UnMinimise: () => void wails.Window.UnMinimise(),

    /**
     * Set the background colour of the window.
     * @param {number} r - The red value between 0 and 255
     * @param {number} g - The green value between 0 and 255
     * @param {number} b - The blue value between 0 and 255
     * @param {number} a - The alpha value between 0 and 255
     */
    SetBackgroundColour: (r, g, b, a) => void wails.Window.SetBackgroundColour(r, g, b, a),
};
