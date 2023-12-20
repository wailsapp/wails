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
 * @typedef {import("./wails").MessageDialogOptions} MessageDialogOptions
 * @typedef {import("./wails").OpenFileDialogOptions} OpenFileDialogOptions
 * @typedef {import("./wails").SaveFileDialogOptions} SaveFileDialogOptions
 */
/**
 * The Dialog API provides methods to interact with system dialogs.
 */
export const Dialog = {
    /**
     * Shows an info dialog
     * @param {MessageDialogOptions} options - options for the dialog
     * @returns {Promise<string>} - returns a promise that resolves to the label of the button clicked by the user
     */
    Info: (options) => wails.Dialog.Info(options),
    /**
     * Shows a warning dialog
     * @param {MessageDialogOptions} options - options for the dialog
     * @returns {Promise<string>} - returns a promise that resolves to the label of the button clicked by the user
     */
    Warning: (options) => wails.Dialog.Warning(options),
    /**
     * Shows an error dialog
     * @param {MessageDialogOptions} options - options for the dialog
     * @returns {Promise<string>} - returns a promise that resolves to the label of the button clicked by the user
     */
    Error: (options) => wails.Dialog.Error(options),

    /**
     * Shows a question dialog
     * @param {MessageDialogOptions} options - options for the dialog
     * @returns {Promise<string>} - returns a promise that resolves to the label of the button clicked by the user
     */
    Question: (options) => wails.Dialog.Question(options),

    /**
     * Shows a file open dialog and returns the files selected by the user.
     * A blank string indicates that the dialog was cancelled.
     * @param {OpenFileDialogOptions} options - options for the dialog
     * @returns {Promise<string[]>|Promise<string>} - returns a promise that resolves to the selected file or an array of selected files if AllowsMultipleSelection is true.
     */
    OpenFile: (options) => wails.Dialog.OpenFile(options),

    /**
     * Shows a file save dialog and returns the filename given by the user.
     * A blank string indicates that the dialog was cancelled.
     * @param {SaveFileDialogOptions} options - options for the dialog
     * @returns {Promise<string>} - returns a promise that resolves to the selected file.
     */
    SaveFile: (options) => wails.Dialog.SaveFile(options),
};
