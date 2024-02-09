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
 * @typedef {Object} OpenFileDialogOptions
 * @property {boolean} [CanChooseDirectories] - Indicates if directories can be chosen.
 * @property {boolean} [CanChooseFiles] - Indicates if files can be chosen.
 * @property {boolean} [CanCreateDirectories] - Indicates if directories can be created.
 * @property {boolean} [ShowHiddenFiles] - Indicates if hidden files should be shown.
 * @property {boolean} [ResolvesAliases] - Indicates if aliases should be resolved.
 * @property {boolean} [AllowsMultipleSelection] - Indicates if multiple selection is allowed.
 * @property {boolean} [HideExtension] - Indicates if the extension should be hidden.
 * @property {boolean} [CanSelectHiddenExtension] - Indicates if hidden extensions can be selected.
 * @property {boolean} [TreatsFilePackagesAsDirectories] - Indicates if file packages should be treated as directories.
 * @property {boolean} [AllowsOtherFiletypes] - Indicates if other file types are allowed.
 * @property {FileFilter[]} [Filters] - Array of file filters.
 * @property {string} [Title] - Title of the dialog.
 * @property {string} [Message] - Message to show in the dialog.
 * @property {string} [ButtonText] - Text to display on the button.
 * @property {string} [Directory] - Directory to open in the dialog.
 * @property {boolean} [Detached] - Indicates if the dialog should appear detached from the main window.
 */


/**
 * @typedef {Object} SaveFileDialogOptions
 * @property {string} [Filename] - Default filename to use in the dialog.
 * @property {boolean} [CanChooseDirectories] - Indicates if directories can be chosen.
 * @property {boolean} [CanChooseFiles] - Indicates if files can be chosen.
 * @property {boolean} [CanCreateDirectories] - Indicates if directories can be created.
 * @property {boolean} [ShowHiddenFiles] - Indicates if hidden files should be shown.
 * @property {boolean} [ResolvesAliases] - Indicates if aliases should be resolved.
 * @property {boolean} [AllowsMultipleSelection] - Indicates if multiple selection is allowed.
 * @property {boolean} [HideExtension] - Indicates if the extension should be hidden.
 * @property {boolean} [CanSelectHiddenExtension] - Indicates if hidden extensions can be selected.
 * @property {boolean} [TreatsFilePackagesAsDirectories] - Indicates if file packages should be treated as directories.
 * @property {boolean} [AllowsOtherFiletypes] - Indicates if other file types are allowed.
 * @property {FileFilter[]} [Filters] - Array of file filters.
 * @property {string} [Title] - Title of the dialog.
 * @property {string} [Message] - Message to show in the dialog.
 * @property {string} [ButtonText] - Text to display on the button.
 * @property {string} [Directory] - Directory to open in the dialog.
 * @property {boolean} [Detached] - Indicates if the dialog should appear detached from the main window.
 */

/**
 * @typedef {Object} MessageDialogOptions
 * @property {string} [Title] - The title of the dialog window.
 * @property {string} [Message] - The main message to show in the dialog.
 * @property {Button[]} [Buttons] - Array of button options to show in the dialog.
 * @property {boolean} [Detached] - True if the dialog should appear detached from the main window (if applicable).
 */

/**
 * @typedef {Object} Button
 * @property {string} [Label] - Text that appears within the button.
 * @property {boolean} [IsCancel] - True if the button should cancel an operation when clicked.
 * @property {boolean} [IsDefault] - True if the button should be the default action when the user presses enter.
 */

/**
 * @typedef {Object} FileFilter
 * @property {string} [DisplayName] - Display name for the filter, it could be "Text Files", "Images" etc.
 * @property {string} [Pattern] - Pattern to match for the filter, e.g. "*.txt;*.md" for text markdown files.
 */

// setup
window._wails = window._wails || {};
window._wails.dialogErrorCallback = dialogErrorCallback;
window._wails.dialogResultCallback = dialogResultCallback;

import {newRuntimeCallerWithID, objectNames} from "./runtime";

import { nanoid } from 'nanoid/non-secure';

// Define constants from the `methods` object in Title Case
const DialogInfo = 0;
const DialogWarning = 1;
const DialogError = 2;
const DialogQuestion = 3;
const DialogOpenFile = 4;
const DialogSaveFile = 5;

const call = newRuntimeCallerWithID(objectNames.Dialog, '');
const dialogResponses = new Map();

/**
 * Generates a unique id that is not present in dialogResponses.
 * @returns {string} unique id
 */
function generateID() {
    let result;
    do {
        result = nanoid();
    } while (dialogResponses.has(result));
    return result;
}

/**
 * Shows a dialog of specified type with the given options.
 * @param {number} type - type of dialog
 * @param {MessageDialogOptions|OpenFileDialogOptions|SaveFileDialogOptions} options - options for the dialog
 * @returns {Promise} promise that resolves with result of dialog
 */
function dialog(type, options = {}) {
    const id = generateID();
    options["dialog-id"] = id;
    return new Promise((resolve, reject) => {
        dialogResponses.set(id, {resolve, reject});
        call(type, options).catch((error) => {
            reject(error);
            dialogResponses.delete(id);
        });
    });
}

/**
 * Handles the callback from a dialog.
 *
 * @param {string} id - The ID of the dialog response.
 * @param {string} data - The data received from the dialog.
 * @param {boolean} isJSON - Flag indicating whether the data is in JSON format.
 *
 * @return {undefined}
 */
function dialogResultCallback(id, data, isJSON) {
    let p = dialogResponses.get(id);
    if (p) {
        if (isJSON) {
            p.resolve(JSON.parse(data));
        } else {
            p.resolve(data);
        }
        dialogResponses.delete(id);
    }
}

/**
 * Callback function for handling errors in dialog.
 *
 * @param {string} id - The id of the dialog response.
 * @param {string} message - The error message.
 *
 * @return {void}
 */
function dialogErrorCallback(id, message) {
    let p = dialogResponses.get(id);
    if (p) {
        p.reject(message);
        dialogResponses.delete(id);
    }
}


// Replace `methods` with constants in Title Case

/**
 * @param {MessageDialogOptions} options - Dialog options
 * @returns {Promise<string>} - The label of the button pressed
 */
export const Info = (options) => dialog(DialogInfo, options);

/**
 * @param {MessageDialogOptions} options - Dialog options
 * @returns {Promise<string>} - The label of the button pressed
 */
export const Warning = (options) => dialog(DialogWarning, options);

/**
 * @param {MessageDialogOptions} options - Dialog options
 * @returns {Promise<string>} - The label of the button pressed
 */
export const Error = (options) => dialog(DialogError, options);

/**
 * @param {MessageDialogOptions} options - Dialog options
 * @returns {Promise<string>} - The label of the button pressed
 */
export const Question = (options) => dialog(DialogQuestion, options);

/**
 * @param {OpenFileDialogOptions} options - Dialog options
 * @returns {Promise<string[]|string>} Returns selected file or list of files. Returns blank string if no file is selected.
 */
export const OpenFile = (options) => dialog(DialogOpenFile, options);

/**
 * @param {SaveFileDialogOptions} options - Dialog options
 * @returns {Promise<string>} Returns the selected file. Returns blank string if no file is selected.
 */
export const SaveFile = (options) => dialog(DialogSaveFile, options);
