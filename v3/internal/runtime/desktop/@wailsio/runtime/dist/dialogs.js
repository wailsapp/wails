/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/
import { newRuntimeCaller, objectNames } from "./runtime.js";
import { nanoid } from './nanoid.js';
// setup
window._wails = window._wails || {};
window._wails.dialogErrorCallback = dialogErrorCallback;
window._wails.dialogResultCallback = dialogResultCallback;
const call = newRuntimeCaller(objectNames.Dialog);
const dialogResponses = new Map();
// Define constants from the `methods` object in Title Case
const DialogInfo = 0;
const DialogWarning = 1;
const DialogError = 2;
const DialogQuestion = 3;
const DialogOpenFile = 4;
const DialogSaveFile = 5;
/**
 * Handles the result of a dialog request.
 *
 * @param id - The id of the request to handle the result for.
 * @param data - The result data of the request.
 * @param isJSON - Indicates whether the data is JSON or not.
 */
function dialogResultCallback(id, data, isJSON) {
    let resolvers = getAndDeleteResponse(id);
    if (!resolvers) {
        return;
    }
    if (isJSON) {
        try {
            resolvers.resolve(JSON.parse(data));
        }
        catch (err) {
            resolvers.reject(new TypeError("could not parse result: " + err.message, { cause: err }));
        }
    }
    else {
        resolvers.resolve(data);
    }
}
/**
 * Handles the error from a dialog request.
 *
 * @param id - The id of the promise handler.
 * @param message - An error message.
 */
function dialogErrorCallback(id, message) {
    var _a;
    (_a = getAndDeleteResponse(id)) === null || _a === void 0 ? void 0 : _a.reject(new window.Error(message));
}
/**
 * Retrieves and removes the response associated with the given ID from the dialogResponses map.
 *
 * @param id - The ID of the response to be retrieved and removed.
 * @returns The response object associated with the given ID, if any.
 */
function getAndDeleteResponse(id) {
    const response = dialogResponses.get(id);
    dialogResponses.delete(id);
    return response;
}
/**
 * Generates a unique ID using the nanoid library.
 *
 * @returns A unique ID that does not exist in the dialogResponses set.
 */
function generateID() {
    let result;
    do {
        result = nanoid();
    } while (dialogResponses.has(result));
    return result;
}
/**
 * Presents a dialog of specified type with the given options.
 *
 * @param type - Dialog type.
 * @param options - Options for the dialog.
 * @returns A promise that resolves with result of dialog.
 */
function dialog(type, options = {}) {
    const id = generateID();
    return new Promise((resolve, reject) => {
        dialogResponses.set(id, { resolve, reject });
        call(type, Object.assign({ "dialog-id": id }, options)).catch((err) => {
            dialogResponses.delete(id);
            reject(err);
        });
    });
}
/**
 * Presents an info dialog.
 *
 * @param options - Dialog options
 * @returns A promise that resolves with the label of the chosen button.
 */
export function Info(options) { return dialog(DialogInfo, options); }
/**
 * Presents a warning dialog.
 *
 * @param options - Dialog options.
 * @returns A promise that resolves with the label of the chosen button.
 */
export function Warning(options) { return dialog(DialogWarning, options); }
/**
 * Presents an error dialog.
 *
 * @param options - Dialog options.
 * @returns A promise that resolves with the label of the chosen button.
 */
export function Error(options) { return dialog(DialogError, options); }
/**
 * Presents a question dialog.
 *
 * @param options - Dialog options.
 * @returns A promise that resolves with the label of the chosen button.
 */
export function Question(options) { return dialog(DialogQuestion, options); }
export function OpenFile(options) { var _a; return (_a = dialog(DialogOpenFile, options)) !== null && _a !== void 0 ? _a : []; }
/**
 * Presents a file selection dialog to pick a file to save.
 *
 * @param options - Dialog options.
 * @returns Selected file, or a blank string if no file has been selected.
 */
export function SaveFile(options) { return dialog(DialogSaveFile, options); }
